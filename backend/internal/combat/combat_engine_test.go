package combat

import (
	"testing"
)

func TestNewCombatEngine(t *testing.T) {
	engine := NewCombatEngine(12345)
	if engine.Seed != 12345 {
		t.Errorf("Expected seed 12345, got %d", engine.Seed)
	}
}

func TestGetShipTypeAdvantage(t *testing.T) {
	engine := NewCombatEngine(0)

	tests := []struct {
		attacker ShipType
		defender ShipType
		expected float64
	}{
		{ShipTypeFrigate, ShipTypeCruiser, 1.05},      // Frigate beats Cruiser
		{ShipTypeCruiser, ShipTypeBattleship, 1.05},  // Cruiser beats Battleship
		{ShipTypeBattleship, ShipTypeFrigate, 1.05},  // Battleship beats Frigate
		{ShipTypeCruiser, ShipTypeFrigate, 0.95},      // Cruiser weak to Frigate
		{ShipTypeBattleship, ShipTypeCruiser, 0.95},  // Battleship weak to Cruiser
		{ShipTypeFrigate, ShipTypeBattleship, 0.95},  // Frigate weak to Battleship
		{ShipTypeFrigate, ShipTypeFrigate, 1.0},      // Same type, no advantage
	}

	for _, tt := range tests {
		attacker := &FleetStack{ShipType: tt.attacker}
		defender := &FleetStack{ShipType: tt.defender}
		result := engine.getShipTypeAdvantage(attacker, defender)
		if result != tt.expected {
			t.Errorf("%s vs %s: expected %.2f, got %.2f", tt.attacker, tt.defender, tt.expected, result)
		}
	}
}

func TestGetArmorEffectiveness(t *testing.T) {
	engine := NewCombatEngine(0)

	tests := []struct {
		damageType DamageType
		armorType  ArmorType
		expected   float64
	}{
		// Chrome armor
		{DamageExplosive, ArmorChrome, 1.25}, // Weak to Explosive
		{DamageKinetic, ArmorChrome, 0.75},   // Resists Kinetic
		{DamageHeat, ArmorChrome, 1.0},       // Neutral
		// Regen armor
		{DamageHeat, ArmorRegen, 1.25},       // Weak to Heat
		{DamageExplosive, ArmorRegen, 0.75},  // Resists Explosive
		{DamageKinetic, ArmorRegen, 1.0},     // Neutral
		// Nano armor
		{DamageMagnetic, ArmorNano, 1.25},    // Weak to Magnetic
		{DamageHeat, ArmorNano, 0.75},        // Resists Heat
		{DamageKinetic, ArmorNano, 1.0},      // Neutral
		// Neutralizing armor
		{DamageKinetic, ArmorNeutralizing, 1.25},   // Weak to Kinetic
		{DamageMagnetic, ArmorNeutralizing, 0.75},  // Resists Magnetic
		{DamageHeat, ArmorNeutralizing, 1.0},       // Neutral
	}

	for _, tt := range tests {
		result := engine.getArmorEffectiveness(tt.damageType, tt.armorType)
		if result != tt.expected {
			t.Errorf("%s vs %s: expected %.2f, got %.2f", tt.damageType, tt.armorType, tt.expected, result)
		}
	}
}

func TestCalculateHitChance(t *testing.T) {
	engine := NewCombatEngine(0)

	tests := []struct {
		accuracy int
		dodge    int
		expected float64
	}{
		{100, 100, 0.5},   // Equal stats, 50% hit chance
		{200, 100, 0.6},   // +100 accuracy, +10% hit
		{100, 200, 0.4},   // +100 dodge, -10% hit
		{1000, 0, 0.95},   // Max hit chance (clamped)
		{0, 1000, 0.05},   // Min hit chance (clamped)
	}

	for _, tt := range tests {
		attacker := &FleetStack{EffectiveAccuracy: tt.accuracy}
		defender := &FleetStack{EffectiveDodge: tt.dodge}
		result := engine.phase4CalculateHitChance(attacker, defender)
		if result != tt.expected {
			t.Errorf("Accuracy %d vs Dodge %d: expected %.2f, got %.2f", tt.accuracy, tt.dodge, tt.expected, result)
		}
	}
}

func TestCalculateDamage(t *testing.T) {
	engine := NewCombatEngine(0)

	attacker := &FleetStack{
		EffectiveAttack: 100,
		EffectiveStacks: 10,
		ShipType:        ShipTypeFrigate,
		DamageType:      DamageKinetic,
	}

	defender := &FleetStack{
		ShipType:  ShipTypeCruiser,
		ArmorType: ArmorChrome,
	}

	// Base damage = 100 * 10 = 1000
	// Frigate vs Cruiser = 1.05x = 1050
	// Kinetic vs Chrome = 0.75x = 787.5 → 788
	expected := 788
	result := engine.phase5CalculateDamage(attacker, defender)

	if result != expected {
		t.Errorf("Expected damage %d, got %d", expected, result)
	}
}

func TestApplyDamage(t *testing.T) {
	engine := NewCombatEngine(0)

	target := &FleetStack{
		CurrentShield:    500,
		CurrentStructure: 1000,
	}

	// Apply 300 damage (all to shield)
	shieldDmg, structDmg := engine.phase6ApplyDamage(target, 300)
	if shieldDmg != 300 || structDmg != 0 {
		t.Errorf("Expected shield=300, struct=0, got shield=%d, struct=%d", shieldDmg, structDmg)
	}
	if target.CurrentShield != 200 || target.CurrentStructure != 1000 {
		t.Errorf("Expected shield=200, struct=1000, got shield=%d, struct=%d", target.CurrentShield, target.CurrentStructure)
	}

	// Apply 300 more damage (200 to shield, 100 to structure)
	shieldDmg, structDmg = engine.phase6ApplyDamage(target, 300)
	if shieldDmg != 200 || structDmg != 100 {
		t.Errorf("Expected shield=200, struct=100, got shield=%d, struct=%d", shieldDmg, structDmg)
	}
	if target.CurrentShield != 0 || target.CurrentStructure != 900 {
		t.Errorf("Expected shield=0, struct=900, got shield=%d, struct=%d", target.CurrentShield, target.CurrentStructure)
	}
}

func TestCalculateCasualties(t *testing.T) {
	engine := NewCombatEngine(0)

	// Stack with 10 ships, structure destroyed
	target := &FleetStack{
		CurrentShips:       10,
		CurrentShield:      0,
		CurrentStructure:   0,
		EffectiveStacks:    10,
		EffectiveStructure: 100,
	}

	destroyed := engine.phase7CalculateCasualties(target)
	if destroyed != 10 {
		t.Errorf("Expected 10 ships destroyed, got %d", destroyed)
	}
	if target.CurrentShips != 0 {
		t.Errorf("Expected 0 ships remaining, got %d", target.CurrentShips)
	}
}

func TestExecuteCombat_BasicScenario(t *testing.T) {
	engine := NewCombatEngine(12345) // Deterministic seed

	// Create simple attacker fleet
	attackerStack := &FleetStack{
		ID:                "a1",
		ShipCount:         100,
		ShipType:          ShipTypeFrigate,
		DamageType:        DamageKinetic,
		ArmorType:         ArmorChrome,
		BaseAttack:        50,
		BaseDefense:       30,
		BaseSpeed:         80,
		BaseAccuracy:      100,
		BaseDodge:         50,
		BaseShield:        200,
		BaseStructure:     300,
		GridRow:           0,
		GridCol:           0,
	}

	attackerFleet := &Fleet{
		PlayerID:       "player1",
		FleetID:        "fleet1",
		CommanderBonus: nil,
		TechBonuses:    &TechBonuses{},
		Stacks:         []*FleetStack{attackerStack},
		Formation:      "phalanx",
		Targeting:      "max_attack",
		Side:           "attacker",
	}

	// Create simple defender fleet
	defenderStack := &FleetStack{
		ID:                "d1",
		ShipCount:         50,
		ShipType:          ShipTypeCruiser,
		DamageType:        DamageExplosive,
		ArmorType:         ArmorRegen,
		BaseAttack:        40,
		BaseDefense:       25,
		BaseSpeed:         60,
		BaseAccuracy:      90,
		BaseDodge:         40,
		BaseShield:        150,
		BaseStructure:     250,
		GridRow:           0,
		GridCol:           0,
	}

	defenderFleet := &Fleet{
		PlayerID:       "player2",
		FleetID:        "fleet2",
		CommanderBonus: nil,
		TechBonuses:    &TechBonuses{},
		Stacks:         []*FleetStack{defenderStack},
		Formation:      "phalanx",
		Targeting:      "max_attack",
		Side:           "defender",
	}

	// Execute combat
	result, err := engine.ExecuteCombat(attackerFleet, defenderFleet)

	if err != nil {
		t.Fatalf("Combat execution failed: %v", err)
	}

	// Verify result structure
	if result.Winner == "" {
		t.Error("Expected a winner, got empty string")
	}

	if result.TotalRounds < 1 {
		t.Error("Expected at least 1 round of combat")
	}

	if result.Winner != "attacker" && result.Winner != "defender" && result.Winner != "draw" {
		t.Errorf("Invalid winner: %s", result.Winner)
	}

	t.Logf("Combat result: %s won in %d rounds", result.Winner, result.TotalRounds)
	t.Logf("Attacker casualties: %d, Defender casualties: %d", result.AttackerCasualties, result.DefenderCasualties)
}

func TestFleetDestroyed(t *testing.T) {
	engine := NewCombatEngine(0)

	aliveStack := &FleetStack{CurrentShips: 10}
	deadStack := &FleetStack{CurrentShips: 0}

	aliveFleet := &Fleet{Stacks: []*FleetStack{aliveStack}}
	deadFleet := &Fleet{Stacks: []*FleetStack{deadStack}}
	mixedFleet := &Fleet{Stacks: []*FleetStack{aliveStack, deadStack}}

	if engine.isFleetDestroyed(aliveFleet) {
		t.Error("Alive fleet should not be destroyed")
	}

	if !engine.isFleetDestroyed(deadFleet) {
		t.Error("Dead fleet should be destroyed")
	}

	if engine.isFleetDestroyed(mixedFleet) {
		t.Error("Mixed fleet should not be destroyed")
	}
}
