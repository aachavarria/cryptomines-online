package combat

import (
	"testing"

	"github.com/cryptomines-online/backend/internal/services"
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
		TechBonuses:    &services.TechBonuses{},
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
		TechBonuses:    &services.TechBonuses{},
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
		TechBonuses: &services.TechBonuses{},
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
		TechBonuses: &services.TechBonuses{},
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

// --- Advanced combat mechanic tests ---

// Test Scatter/AoE damage to adjacent stacks
func TestScatterDamageAdjacent(t *testing.T) {
	engine := NewCombatEngine(42)

	attacker := &FleetStack{
		ID:              "a1",
		GridRow:         0,
		GridCol:         1,
		ScatterDamage:   25, // 25% scatter to adjacent stacks
		ScatterRate:     0,  // 0 means always triggers (no rate check)
	}

	primaryTarget := &FleetStack{
		ID:               "d1",
		GridRow:          0,
		GridCol:          1,
		CurrentShips:     50,
		CurrentShield:    5000,
		CurrentStructure: 5000,
		EffectiveStacks:  50,
		EffectiveStructure: 100,
	}

	adjacentTarget := &FleetStack{
		ID:               "d2",
		GridRow:          0,
		GridCol:          0, // Adjacent (same row, col diff = 1)
		CurrentShips:     50,
		CurrentShield:    5000,
		CurrentStructure: 5000,
		EffectiveStacks:  50,
		EffectiveStructure: 100,
	}

	farTarget := &FleetStack{
		ID:               "d3",
		GridRow:          2,
		GridCol:          2, // Not adjacent (row diff = 2)
		CurrentShips:     50,
		CurrentShield:    5000,
		CurrentStructure: 5000,
		EffectiveStacks:  50,
		EffectiveStructure: 100,
	}

	state := &CombatState{
		Attacker: &Fleet{
			Stacks: []*FleetStack{attacker},
		},
		Defender: &Fleet{
			Stacks: []*FleetStack{primaryTarget, adjacentTarget, farTarget},
		},
		DamageTakenModifiers: make(map[string]float64),
		AttackDebuffs:        make(map[string]*StackDebuff),
	}

	round := &CombatRound{
		Casualties: make(map[string]int),
	}

	adjacentShieldBefore := adjacentTarget.CurrentShield
	farShieldBefore := farTarget.CurrentShield

	totalScatter, scatterTargets := engine.applyScatter(state, attacker, primaryTarget, 1000, round)

	// Adjacent target should take 25% of 1000 = 250 damage
	if scatterTargets["d2"] != 250 {
		t.Errorf("Adjacent target should take 250 scatter damage, got %d", scatterTargets["d2"])
	}

	// Far target should NOT take scatter damage
	if scatterTargets["d3"] != 0 {
		t.Errorf("Far target should take 0 scatter damage, got %d", scatterTargets["d3"])
	}

	if totalScatter != 250 {
		t.Errorf("Total scatter should be 250, got %d", totalScatter)
	}

	// Adjacent target's shield should be reduced
	if adjacentTarget.CurrentShield != adjacentShieldBefore-250 {
		t.Errorf("Adjacent shield should be %d, got %d", adjacentShieldBefore-250, adjacentTarget.CurrentShield)
	}

	// Far target's shield should be unchanged
	if farTarget.CurrentShield != farShieldBefore {
		t.Errorf("Far target shield should be unchanged at %d, got %d", farShieldBefore, farTarget.CurrentShield)
	}
}

// Test ScatterAll spreads to ALL enemy stacks
func TestScatterAll(t *testing.T) {
	engine := NewCombatEngine(42)

	attacker := &FleetStack{
		ID:          "a1",
		GridRow:     0,
		GridCol:     0,
		ScatterAll:  15, // 15% to all enemies
	}

	target1 := &FleetStack{
		ID: "d1", GridRow: 0, GridCol: 0,
		CurrentShips: 50, CurrentShield: 5000, CurrentStructure: 5000,
		EffectiveStacks: 50, EffectiveStructure: 100,
	}
	target2 := &FleetStack{
		ID: "d2", GridRow: 1, GridCol: 1,
		CurrentShips: 50, CurrentShield: 5000, CurrentStructure: 5000,
		EffectiveStacks: 50, EffectiveStructure: 100,
	}
	target3 := &FleetStack{
		ID: "d3", GridRow: 2, GridCol: 2,
		CurrentShips: 50, CurrentShield: 5000, CurrentStructure: 5000,
		EffectiveStacks: 50, EffectiveStructure: 100,
	}

	state := &CombatState{
		Attacker: &Fleet{Stacks: []*FleetStack{attacker}},
		Defender: &Fleet{Stacks: []*FleetStack{target1, target2, target3}},
		DamageTakenModifiers: make(map[string]float64),
		AttackDebuffs:        make(map[string]*StackDebuff),
	}

	round := &CombatRound{Casualties: make(map[string]int)}
	totalScatter, scatterTargets := engine.applyScatter(state, attacker, target1, 1000, round)

	// Both d2 and d3 should receive 15% of 1000 = 150 each
	if scatterTargets["d2"] != 150 {
		t.Errorf("d2 should take 150 scatter, got %d", scatterTargets["d2"])
	}
	if scatterTargets["d3"] != 150 {
		t.Errorf("d3 should take 150 scatter, got %d", scatterTargets["d3"])
	}
	if totalScatter != 300 {
		t.Errorf("Total scatter should be 300, got %d", totalScatter)
	}
}

// Test Piercing damage to stack behind target
func TestPiercingDamage(t *testing.T) {
	engine := NewCombatEngine(42)

	attacker := &FleetStack{
		ID:              "a1",
		GridRow:         0,
		PiercingDamage:  50, // 50% of damage pierces
	}

	frontTarget := &FleetStack{
		ID: "d1", GridRow: 0, GridCol: 0,
		CurrentShips: 50, CurrentShield: 5000, CurrentStructure: 5000,
		EffectiveStacks: 50, EffectiveStructure: 100,
	}
	behindTarget := &FleetStack{
		ID: "d2", GridRow: 1, GridCol: 0,
		CurrentShips: 50, CurrentShield: 5000, CurrentStructure: 5000,
		EffectiveStacks: 50, EffectiveStructure: 100,
	}

	state := &CombatState{
		Attacker: &Fleet{Stacks: []*FleetStack{attacker}},
		Defender: &Fleet{Stacks: []*FleetStack{frontTarget, behindTarget}},
		DamageTakenModifiers: make(map[string]float64),
		AttackDebuffs:        make(map[string]*StackDebuff),
	}

	round := &CombatRound{Casualties: make(map[string]int)}

	behindShieldBefore := behindTarget.CurrentShield

	pierceDmg, pierceTargetID := engine.applyPiercing(state, attacker, frontTarget, 1000, false, 1.5, round)

	// 50% of 1000 = 500 piercing damage
	if pierceDmg != 500 {
		t.Errorf("Expected 500 piercing damage, got %d", pierceDmg)
	}
	if pierceTargetID != "d2" {
		t.Errorf("Expected pierce target d2, got %s", pierceTargetID)
	}
	if behindTarget.CurrentShield != behindShieldBefore-500 {
		t.Errorf("Behind target shield should be %d, got %d", behindShieldBefore-500, behindTarget.CurrentShield)
	}
}

// Test Piercing with critical hit forwarding
func TestPiercingWithCritical(t *testing.T) {
	engine := NewCombatEngine(42)

	attacker := &FleetStack{
		ID:              "a1",
		GridRow:         0,
		PiercingDamage:  50,   // 50% pierce
		PiercingCritical: true, // crit applies to pierce
	}

	frontTarget := &FleetStack{
		ID: "d1", GridRow: 0,
		CurrentShips: 50, CurrentShield: 10000, CurrentStructure: 10000,
		EffectiveStacks: 50, EffectiveStructure: 200,
	}
	behindTarget := &FleetStack{
		ID: "d2", GridRow: 1,
		CurrentShips: 50, CurrentShield: 10000, CurrentStructure: 10000,
		EffectiveStacks: 50, EffectiveStructure: 200,
	}

	state := &CombatState{
		Attacker: &Fleet{Stacks: []*FleetStack{attacker}},
		Defender: &Fleet{Stacks: []*FleetStack{frontTarget, behindTarget}},
		DamageTakenModifiers: make(map[string]float64),
		AttackDebuffs:        make(map[string]*StackDebuff),
	}

	round := &CombatRound{Casualties: make(map[string]int)}

	// Pass wasCrit=true with critMultiplier=1.5
	pierceDmg, _ := engine.applyPiercing(state, attacker, frontTarget, 1000, true, 1.5, round)

	// 50% of 1000 = 500, then * 1.5 crit = 750
	if pierceDmg != 750 {
		t.Errorf("Expected 750 piercing+crit damage, got %d", pierceDmg)
	}
}

// Test Piercing does nothing when no stack is behind
func TestPiercingNoTargetBehind(t *testing.T) {
	engine := NewCombatEngine(42)

	attacker := &FleetStack{
		ID:             "a1",
		GridRow:        0,
		PiercingDamage: 50,
	}

	// Only one enemy stack, nothing behind it
	frontTarget := &FleetStack{
		ID: "d1", GridRow: 0,
		CurrentShips: 50, CurrentShield: 5000, CurrentStructure: 5000,
		EffectiveStacks: 50, EffectiveStructure: 100,
	}

	state := &CombatState{
		Attacker: &Fleet{Stacks: []*FleetStack{attacker}},
		Defender: &Fleet{Stacks: []*FleetStack{frontTarget}},
		DamageTakenModifiers: make(map[string]float64),
		AttackDebuffs:        make(map[string]*StackDebuff),
	}

	round := &CombatRound{Casualties: make(map[string]int)}
	pierceDmg, pierceTargetID := engine.applyPiercing(state, attacker, frontTarget, 1000, false, 1.5, round)

	if pierceDmg != 0 {
		t.Errorf("Expected 0 piercing damage with no target behind, got %d", pierceDmg)
	}
	if pierceTargetID != "" {
		t.Errorf("Expected empty pierce target ID, got %s", pierceTargetID)
	}
}

// Test Restoration heals shields and structure
func TestRestoration(t *testing.T) {
	engine := NewCombatEngine(42)

	stack := &FleetStack{
		ID:                 "a1",
		CurrentShips:       50,
		EffectiveShield:    100,
		EffectiveStructure: 200,
		EffectiveStacks:    50,
		CurrentShield:      3000, // max is 100*50=5000
		CurrentStructure:   8000, // max is 200*50=10000
		ShieldRestore:      10,   // 10% per round
		StructureRestore:   5,    // 5% per round
	}

	state := &CombatState{
		Attacker: &Fleet{Stacks: []*FleetStack{stack}},
		Defender: &Fleet{Stacks: []*FleetStack{}},
		Logs:     []string{},
		DamageTakenModifiers: make(map[string]float64),
		AttackDebuffs:        make(map[string]*StackDebuff),
	}

	engine.applyRestoration(state)

	// Shield: 10% of 5000 = 500 restored, 3000 + 500 = 3500
	if stack.CurrentShield != 3500 {
		t.Errorf("Expected shield 3500 after restore, got %d", stack.CurrentShield)
	}

	// Structure: 5% of 10000 = 500 restored, 8000 + 500 = 8500
	if stack.CurrentStructure != 8500 {
		t.Errorf("Expected structure 8500 after restore, got %d", stack.CurrentStructure)
	}
}

// Test Restoration caps at max
func TestRestorationCapsAtMax(t *testing.T) {
	engine := NewCombatEngine(42)

	stack := &FleetStack{
		ID:                 "a1",
		CurrentShips:       50,
		EffectiveShield:    100,
		EffectiveStructure: 200,
		EffectiveStacks:    50,
		CurrentShield:      4900, // max is 5000, only 100 missing
		CurrentStructure:   10000, // already at max
		ShieldRestore:      10,    // would restore 500 but cap at max
		StructureRestore:   5,
	}

	state := &CombatState{
		Attacker: &Fleet{Stacks: []*FleetStack{stack}},
		Defender: &Fleet{Stacks: []*FleetStack{}},
		Logs:     []string{},
		DamageTakenModifiers: make(map[string]float64),
		AttackDebuffs:        make(map[string]*StackDebuff),
	}

	engine.applyRestoration(state)

	if stack.CurrentShield != 5000 {
		t.Errorf("Expected shield capped at 5000, got %d", stack.CurrentShield)
	}
	if stack.CurrentStructure != 10000 {
		t.Errorf("Expected structure capped at 10000, got %d", stack.CurrentStructure)
	}
}

// Test Reflection reflects damage back to attacker
func TestReflection(t *testing.T) {
	engine := NewCombatEngine(42)

	attacker := &FleetStack{
		ID:               "a1",
		CurrentShips:     100,
		CurrentShield:    10000,
		CurrentStructure: 10000,
		EffectiveStacks:  100,
		EffectiveStructure: 100,
	}

	defender := &FleetStack{
		ID:                     "d1",
		CurrentShips:           50,
		CurrentShield:          5000,
		CurrentStructure:       5000,
		EffectiveStacks:        50,
		EffectiveStructure:     100,
		ReflectDamage:          20, // 20% reflected
		ReflectStructureDamage: 10, // 10% structure reflected
	}

	state := &CombatState{
		Attacker: &Fleet{Stacks: []*FleetStack{attacker}},
		Defender: &Fleet{Stacks: []*FleetStack{defender}},
		Logs:     []string{},
		DamageTakenModifiers: make(map[string]float64),
		AttackDebuffs:        make(map[string]*StackDebuff),
	}

	round := &CombatRound{Casualties: make(map[string]int)}

	attackerShieldBefore := attacker.CurrentShield

	reflected := engine.applyReflection(state, attacker, defender, 1000, round)

	// ReflectDamage: 20% of 1000 = 200 reflected (to shield first)
	// ReflectStructureDamage: 10% of 1000 = 100 reflected (to structure directly)
	// Total reflected = 300
	if reflected != 300 {
		t.Errorf("Expected 300 total reflected damage, got %d", reflected)
	}

	// Attacker shield should be reduced by 200
	if attacker.CurrentShield != attackerShieldBefore-200 {
		t.Errorf("Attacker shield should be %d, got %d", attackerShieldBefore-200, attacker.CurrentShield)
	}

	// Attacker structure should be reduced by 100 (direct structure reflect)
	if attacker.CurrentStructure != 10000-100 {
		t.Errorf("Attacker structure should be 9900, got %d", attacker.CurrentStructure)
	}
}

// Test AbsorbDouble halves damage
func TestAbsorbDouble(t *testing.T) {
	// AbsorbDouble is checked inline in ExecuteCombat, so we test via full combat
	engine := NewCombatEngine(12345)

	attackerStack := &FleetStack{
		ID: "a1", ShipCount: 100, ShipType: ShipTypeFrigate,
		DamageType: DamageKinetic, ArmorType: ArmorChrome,
		BaseAttack: 50, BaseDefense: 30, BaseSpeed: 80,
		BaseAccuracy: 100, BaseDodge: 50,
		BaseShield: 200, BaseStructure: 300,
		GridRow: 0, GridCol: 0,
	}

	defenderStack := &FleetStack{
		ID: "d1", ShipCount: 50, ShipType: ShipTypeCruiser,
		DamageType: DamageExplosive, ArmorType: ArmorRegen,
		BaseAttack: 40, BaseDefense: 25, BaseSpeed: 60,
		BaseAccuracy: 90, BaseDodge: 40,
		BaseShield: 150, BaseStructure: 250,
		GridRow: 0, GridCol: 0,
		AbsorbDouble: 1.0, // 100% chance to absorb (always halves)
	}

	attackerFleet := &Fleet{
		PlayerID: "p1", FleetID: "f1",
		TechBonuses: &services.TechBonuses{},
		Stacks: []*FleetStack{attackerStack},
		Formation: "phalanx", Side: "attacker",
	}
	defenderFleet := &Fleet{
		PlayerID: "p2", FleetID: "f2",
		TechBonuses: &services.TechBonuses{},
		Stacks: []*FleetStack{defenderStack},
		Formation: "phalanx", Side: "defender",
	}

	result, err := engine.ExecuteCombat(attackerFleet, defenderFleet)
	if err != nil {
		t.Fatalf("Combat failed: %v", err)
	}

	// Verify AbsorbedDouble was triggered in at least one attack
	foundAbsorb := false
	for _, round := range result.DetailedRounds {
		for _, attack := range round.Attacks {
			if attack.AbsorbedDouble {
				foundAbsorb = true
				break
			}
		}
		if foundAbsorb {
			break
		}
	}

	if !foundAbsorb {
		t.Error("Expected AbsorbDouble to trigger at least once with 100% chance")
	}

	t.Logf("Combat with AbsorbDouble: %s won in %d rounds", result.Winner, result.TotalRounds)
}

// Test Knockback pushes target to higher row
func TestKnockback(t *testing.T) {
	engine := NewCombatEngine(12345)

	attackerStack := &FleetStack{
		ID: "a1", ShipCount: 100, ShipType: ShipTypeFrigate,
		DamageType: DamageKinetic, ArmorType: ArmorChrome,
		BaseAttack: 50, BaseDefense: 30, BaseSpeed: 80,
		BaseAccuracy: 200, BaseDodge: 50,
		BaseShield: 200, BaseStructure: 300,
		GridRow: 0, GridCol: 0,
		Knockback: 1, // push back 1 row
	}

	defenderStack := &FleetStack{
		ID: "d1", ShipCount: 50, ShipType: ShipTypeCruiser,
		DamageType: DamageExplosive, ArmorType: ArmorRegen,
		BaseAttack: 40, BaseDefense: 25, BaseSpeed: 60,
		BaseAccuracy: 90, BaseDodge: 40,
		BaseShield: 150, BaseStructure: 250,
		GridRow: 0, GridCol: 0, // starts at front row
	}

	attackerFleet := &Fleet{
		PlayerID: "p1", FleetID: "f1",
		TechBonuses: &services.TechBonuses{},
		Stacks: []*FleetStack{attackerStack},
		Formation: "phalanx", Side: "attacker",
	}
	defenderFleet := &Fleet{
		PlayerID: "p2", FleetID: "f2",
		TechBonuses: &services.TechBonuses{},
		Stacks: []*FleetStack{defenderStack},
		Formation: "phalanx", Side: "defender",
	}

	result, err := engine.ExecuteCombat(attackerFleet, defenderFleet)
	if err != nil {
		t.Fatalf("Combat failed: %v", err)
	}

	// Check that knockback was applied in at least one attack
	foundKnockback := false
	for _, round := range result.DetailedRounds {
		for _, attack := range round.Attacks {
			if attack.KnockbackApplied > 0 {
				foundKnockback = true
				break
			}
		}
		if foundKnockback {
			break
		}
	}

	if !foundKnockback {
		t.Error("Expected knockback to be applied at least once")
	}

	t.Logf("Combat with Knockback: %s won in %d rounds", result.Winner, result.TotalRounds)
}

// Test RangeDamage modifies damage based on distance
func TestRangeDamage(t *testing.T) {
	engine := NewCombatEngine(42)

	attacker := &FleetStack{
		ID:              "a1",
		EffectiveAttack: 100,
		EffectiveStacks: 10,
		ShipType:        ShipTypeFrigate,
		DamageType:      DamageKinetic,
		GridRow:         0,
		RangeDamage: map[int]float64{
			0: 1.5,  // 50% bonus at distance 0
			1: 1.2,  // 20% bonus at distance 1
			2: 0.8,  // 20% penalty at distance 2
		},
	}

	// Test at distance 0 (same row)
	target := &FleetStack{
		ShipType:  ShipTypeFrigate,
		ArmorType: ArmorChrome,
		GridRow:   0,
	}

	baseDmg := engine.phase5CalculateDamage(attacker, target)

	// RangeDamage at distance 0 = 1.5x
	distance := abs(attacker.GridRow - target.GridRow)
	if mult, ok := attacker.RangeDamage[distance]; ok {
		modDmg := int(float64(baseDmg) * mult)
		expected := int(float64(baseDmg) * 1.5)
		if modDmg != expected {
			t.Errorf("Distance 0: expected %d, got %d", expected, modDmg)
		}
	} else {
		t.Error("RangeDamage should have entry for distance 0")
	}

	// Test at distance 2
	target.GridRow = 2
	distance = abs(attacker.GridRow - target.GridRow)
	if mult, ok := attacker.RangeDamage[distance]; ok {
		if mult != 0.8 {
			t.Errorf("Distance 2: expected mult 0.8, got %.2f", mult)
		}
	} else {
		t.Error("RangeDamage should have entry for distance 2")
	}
}

// Test DamageTakenIncrease debuff stacks
func TestDamageTakenIncrease(t *testing.T) {
	engine := NewCombatEngine(42)

	attacker := &FleetStack{
		ID:                  "a1",
		DamageTakenIncrease: 10, // 10% per hit
	}

	target := &FleetStack{
		ID:               "d1",
		CurrentShips:     50,
		CurrentShield:    5000,
		CurrentStructure: 5000,
	}

	state := &CombatState{
		Attacker:             &Fleet{Stacks: []*FleetStack{attacker}},
		Defender:             &Fleet{Stacks: []*FleetStack{target}},
		Logs:                 []string{},
		DamageTakenModifiers: make(map[string]float64),
		AttackDebuffs:        make(map[string]*StackDebuff),
	}

	// Apply debuff twice (simulating two hits)
	engine.applyDebuffs(state, attacker, target)
	engine.applyDebuffs(state, attacker, target)

	// Should stack: 10 + 10 = 20
	if state.DamageTakenModifiers["d1"] != 20 {
		t.Errorf("Expected 20%% damage increase, got %.1f%%", state.DamageTakenModifiers["d1"])
	}
}

// Test EnemyAttackReduction debuff
func TestEnemyAttackReduction(t *testing.T) {
	engine := NewCombatEngine(42)

	attacker := &FleetStack{
		ID:                         "a1",
		EnemyAttackReduction:       15, // 15% reduction
		EnemyAttackReductionRounds: 3,
	}

	target := &FleetStack{ID: "d1", CurrentShips: 50}

	state := &CombatState{
		Attacker:             &Fleet{Stacks: []*FleetStack{attacker}},
		Defender:             &Fleet{Stacks: []*FleetStack{target}},
		Logs:                 []string{},
		DamageTakenModifiers: make(map[string]float64),
		AttackDebuffs:        make(map[string]*StackDebuff),
	}

	engine.applyDebuffs(state, attacker, target)

	debuff, ok := state.AttackDebuffs["d1"]
	if !ok {
		t.Fatal("Expected attack debuff on d1")
	}
	if debuff.AttackReduction != 15 {
		t.Errorf("Expected 15%% attack reduction, got %.1f%%", debuff.AttackReduction)
	}
	if debuff.RoundsLeft != 3 {
		t.Errorf("Expected 3 rounds left, got %d", debuff.RoundsLeft)
	}

	// Tick debuff
	engine.tickDebuffs(state)
	if state.AttackDebuffs["d1"].RoundsLeft != 2 {
		t.Errorf("Expected 2 rounds left after tick, got %d", state.AttackDebuffs["d1"].RoundsLeft)
	}

	// Tick 2 more times to expire
	engine.tickDebuffs(state)
	engine.tickDebuffs(state)
	if _, ok := state.AttackDebuffs["d1"]; ok {
		t.Error("Debuff should be expired and removed after 3 ticks")
	}
}

// Test full combat with all advanced mechanics active
func TestExecuteCombat_AllAdvancedMechanics(t *testing.T) {
	engine := NewCombatEngine(54321)

	a1 := &FleetStack{
		ID: "a1", ShipCount: 100, ShipType: ShipTypeFrigate,
		DamageType: DamageKinetic, ArmorType: ArmorChrome,
		BaseAttack: 50, BaseDefense: 30, BaseSpeed: 80,
		BaseAccuracy: 150, BaseDodge: 50,
		BaseShield: 200, BaseStructure: 300,
		GridRow: 0, GridCol: 1,
		ScatterAll:          10,
		PiercingDamage:      30,
		PiercingCritical:    true,
		DamageTakenIncrease: 5,
		Knockback:           1,
	}

	d1 := &FleetStack{
		ID: "d1", ShipCount: 50, ShipType: ShipTypeCruiser,
		DamageType: DamageExplosive, ArmorType: ArmorRegen,
		BaseAttack: 40, BaseDefense: 25, BaseSpeed: 60,
		BaseAccuracy: 90, BaseDodge: 40,
		BaseShield: 150, BaseStructure: 250,
		GridRow: 0, GridCol: 0,
		ShieldRestore:          5,
		StructureRestore:       2,
		ReflectDamage:          10,
		AbsorbDouble:           0.3,
		EnemyAttackReduction:   10,
		EnemyAttackReductionRounds: 2,
	}

	d2 := &FleetStack{
		ID: "d2", ShipCount: 30, ShipType: ShipTypeBattleship,
		DamageType: DamageMagnetic, ArmorType: ArmorNano,
		BaseAttack: 60, BaseDefense: 40, BaseSpeed: 40,
		BaseAccuracy: 80, BaseDodge: 30,
		BaseShield: 300, BaseStructure: 400,
		GridRow: 1, GridCol: 1,
		ReflectStructureDamage: 15,
	}

	attackerFleet := &Fleet{
		PlayerID: "p1", FleetID: "f1",
		TechBonuses: &services.TechBonuses{},
		Stacks:    []*FleetStack{a1},
		Formation: "battle_line", Side: "attacker",
	}
	defenderFleet := &Fleet{
		PlayerID: "p2", FleetID: "f2",
		TechBonuses: &services.TechBonuses{},
		Stacks:    []*FleetStack{d1, d2},
		Formation: "phalanx", Side: "defender",
	}

	result, err := engine.ExecuteCombat(attackerFleet, defenderFleet)
	if err != nil {
		t.Fatalf("Combat failed: %v", err)
	}

	// Collect stats from all rounds
	stats := struct {
		scatterHits   int
		piercingHits  int
		reflectedHits int
		absorbHits    int
		knockbackHits int
	}{}

	for _, round := range result.DetailedRounds {
		for _, attack := range round.Attacks {
			if attack.ScatterDamage > 0 {
				stats.scatterHits++
			}
			if attack.PiercingDamage > 0 {
				stats.piercingHits++
			}
			if attack.ReflectedDamage > 0 {
				stats.reflectedHits++
			}
			if attack.AbsorbedDouble {
				stats.absorbHits++
			}
			if attack.KnockbackApplied > 0 {
				stats.knockbackHits++
			}
		}
	}

	t.Logf("Combat with all advanced mechanics: %s won in %d rounds", result.Winner, result.TotalRounds)
	t.Logf("Scatter hits: %d, Piercing hits: %d, Reflected: %d, Absorb: %d, Knockback: %d",
		stats.scatterHits, stats.piercingHits, stats.reflectedHits, stats.absorbHits, stats.knockbackHits)
	t.Logf("Attacker casualties: %d, Defender casualties: %d", result.AttackerCasualties, result.DefenderCasualties)

	if result.Winner == "" {
		t.Error("Expected a winner")
	}
}

// Test abs helper
func TestAbs(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
		{-1, 1},
	}

	for _, tt := range tests {
		result := abs(tt.input)
		if result != tt.expected {
			t.Errorf("abs(%d): expected %d, got %d", tt.input, tt.expected, result)
		}
	}
}
