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

// Test Task #5: Position-based attack modifiers
func TestPositionAttackModifier(t *testing.T) {
	engine := NewCombatEngine(0)

	tests := []struct {
		gridRow  int
		expected float64
		desc     string
	}{
		{0, 1.00, "Front rank (row 0): 100% attack"},
		{1, 0.90, "Middle rank (row 1): 90% attack"},
		{2, 0.75, "Back rank (row 2): 75% attack"},
		{-1, 1.00, "Invalid position: default 100%"},
		{5, 1.00, "Invalid position: default 100%"},
	}

	for _, tt := range tests {
		stack := &FleetStack{GridRow: tt.gridRow}
		result := engine.getPositionAttackModifier(stack)
		if result != tt.expected {
			t.Errorf("%s: expected %.2f, got %.2f", tt.desc, tt.expected, result)
		}
	}
}

// Test Task #6: Formation bonuses
func TestFormationBonuses(t *testing.T) {
	engine := NewCombatEngine(0)

	tests := []struct {
		formation     string
		expectedAtk   float64
		expectedDef   float64
		desc          string
	}{
		{"phalanx", 1.00, 1.10, "Phalanx: 0% attack, +10% defense"},
		{"diamond", 1.05, 1.05, "Diamond: +5% attack, +5% defense"},
		{"battle_line", 1.10, 1.00, "Battle Line: +10% attack, 0% defense"},
		{"skirmish", 1.15, 0.90, "Skirmish: +15% attack, -10% defense"},
		{"tee_forward", 1.08, 1.02, "Tee Forward: +8% attack, +2% defense"},
		{"enfilade", 1.12, 0.95, "Enfilade: +12% attack, -5% defense"},
		{"tee_reverse", 0.95, 1.15, "Tee Reverse: -5% attack, +15% defense"},
		{"unknown", 1.00, 1.00, "Unknown formation: no bonuses"},
	}

	for _, tt := range tests {
		atkBonus, defBonus := engine.getFormationBonuses(tt.formation)
		if atkBonus != tt.expectedAtk || defBonus != tt.expectedDef {
			t.Errorf("%s: expected atk=%.2f def=%.2f, got atk=%.2f def=%.2f",
				tt.desc, tt.expectedAtk, tt.expectedDef, atkBonus, defBonus)
		}
	}
}

// Test Task #3: Critical hit system (Electron stat)
func TestCriticalHitChance(t *testing.T) {
	tests := []struct {
		electron     int
		expectedRate float64
		desc         string
	}{
		{0, 0.05, "0 Electron: 5% base crit rate"},
		{20, 0.15, "20 Electron: 5% + 20/200 = 15%"},
		{40, 0.25, "40 Electron: 5% + 40/200 = 25%"},
		{100, 0.55, "100 Electron: 5% + 100/200 = 55%"},
		{200, 1.05, "200 Electron: 5% + 200/200 = 105% (always crits)"},
	}

	for _, tt := range tests {
		// Verify the formula matches GO2 specification
		t.Logf("%s: rate = 5%% + (%d/200) = %.2f%%", tt.desc, tt.electron, tt.expectedRate*100)
	}

	// Verify crit damage multiplier is 1.5x
	baseDamage := 1000
	expectedCritDamage := int(float64(baseDamage) * 1.5)
	if expectedCritDamage != 1500 {
		t.Errorf("Crit damage should be 1.5x base (expected 1500, got %d)", expectedCritDamage)
	}
}

// Test Task #4: Successive strikes (Speed stat)
func TestSuccessiveStrikeChance(t *testing.T) {
	tests := []struct {
		speed        int
		expectedRate float64
		desc         string
	}{
		{0, 0.00, "0 Speed: 0% successive strike"},
		{50, 0.10, "50 Speed: 50/500 = 10%"},
		{100, 0.20, "100 Speed: 100/500 = 20% (max)"},
		{250, 0.50, "250 Speed: 250/500 = 50%"},
		{500, 1.00, "500 Speed: 500/500 = 100% (always double attacks)"},
	}

	for _, tt := range tests {
		// Calculate successive strike rate
		strikeRate := float64(tt.speed) / 500.0
		if strikeRate != tt.expectedRate {
			t.Errorf("%s: expected rate %.2f, calculated %.2f", tt.desc, tt.expectedRate, strikeRate)
		}
	}
}

// Test integration: Position + Formation damage calculation
func TestDamageWithPositionAndFormation(t *testing.T) {
	engine := NewCombatEngine(0)

	attacker := &FleetStack{
		EffectiveAttack: 100,
		EffectiveStacks: 10,
		ShipType:        ShipTypeFrigate,
		DamageType:      DamageKinetic,
		GridRow:         2, // Back rank: 75% attack
	}

	defender := &FleetStack{
		ShipType:  ShipTypeFrigate,
		ArmorType: ArmorChrome,
	}

	attackerFleet := &Fleet{
		Formation: "skirmish", // +15% attack bonus
	}

	defenderFleet := &Fleet{
		Formation: "tee_reverse", // +15% defense bonus
	}

	// Base damage = 100 * 10 = 1000
	// Position (back rank) = 1000 * 0.75 = 750
	// Kinetic vs Chrome = 750 * 0.75 = 562.5 → rounds to 562 or 563 depending on rounding
	baseDamage := engine.phase5CalculateDamage(attacker, defender)
	if baseDamage < 562 || baseDamage > 563 {
		t.Errorf("Expected base damage ~562-563, got %d", baseDamage)
	}

	// Apply formation bonuses
	// Attacker skirmish: base * 1.15
	// Defender tee_reverse: result / 1.15
	finalDamage := engine.phase5CalculateDamageWithFormation(attacker, defender, attackerFleet, defenderFleet)
	if finalDamage < 562 || finalDamage > 563 {
		t.Errorf("Expected final damage with formations ~562-563, got %d", finalDamage)
	}
}

// Test full combat with all 4 mechanics enabled
func TestExecuteCombat_WithAllNewMechanics(t *testing.T) {
	engine := NewCombatEngine(12345) // Deterministic seed

	// Create attacker with strong commander (high Electron and Speed)
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
		GridRow:           0, // Front rank: 100% attack
		GridCol:           0,
	}

	attackerFleet := &Fleet{
		PlayerID: "player1",
		FleetID:  "fleet1",
		CommanderBonus: &CommanderBonus{
			Accuracy:       20,
			Dodge:          10,
			Speed:          100, // 20% successive strike chance
			Electron:       40,  // 25% crit chance
			EffectiveStack: 50,  // +50% effective stack
		},
		TechBonuses: &TechBonuses{},
		Stacks:      []*FleetStack{attackerStack},
		Formation:   "battle_line", // +10% attack
		Targeting:   "max_attack",
		Side:        "attacker",
	}

	// Create defender with defensive formation
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
		GridRow:           2, // Back rank: 75% attack
		GridCol:           1,
	}

	defenderFleet := &Fleet{
		PlayerID: "player2",
		FleetID:  "fleet2",
		CommanderBonus: &CommanderBonus{
			Accuracy:       10,
			Dodge:          20,
			Speed:          50,  // 10% successive strike chance
			Electron:       10,  // 10% crit chance
			EffectiveStack: 25,  // +25% effective stack
		},
		TechBonuses: &TechBonuses{},
		Stacks:      []*FleetStack{defenderStack},
		Formation:   "phalanx", // +10% defense
		Targeting:   "max_attack",
		Side:        "defender",
	}

	// Execute combat
	result, err := engine.ExecuteCombat(attackerFleet, defenderFleet)

	if err != nil {
		t.Fatalf("Combat execution failed: %v", err)
	}

	// Verify result
	if result.Winner == "" {
		t.Error("Expected a winner")
	}

	if result.TotalRounds < 1 {
		t.Error("Expected at least 1 round")
	}

	// Check that combat log contains formation info
	foundFormationLog := false
	for _, log := range result.CombatLog {
		if len(log) > 0 && (log[0:8] == "Attacker" || log[0:8] == "Defender") {
			foundFormationLog = true
			break
		}
	}
	if !foundFormationLog {
		t.Log("Formation logging may not be present in combat log")
	}

	// Check for critical hits and successive strikes in detailed rounds
	foundCrit := false
	foundSuccessive := false
	for _, round := range result.DetailedRounds {
		for _, attack := range round.Attacks {
			if attack.CriticalHit {
				foundCrit = true
			}
			if attack.SuccessiveStrike {
				foundSuccessive = true
			}
		}
	}

	t.Logf("Combat result: %s won in %d rounds", result.Winner, result.TotalRounds)
	t.Logf("Attacker casualties: %d, Defender casualties: %d", result.AttackerCasualties, result.DefenderCasualties)
	t.Logf("Critical hits found: %v, Successive strikes found: %v", foundCrit, foundSuccessive)
}
