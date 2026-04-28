package combat

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/cryptomines-online/backend/internal/services"
)

// CombatEngine handles 8-phase combat resolution for Galaxy Online 2
type CombatEngine struct {
	Seed int64
	rng  *rand.Rand
}

// NewCombatEngine creates a new combat engine with optional seed for deterministic testing
func NewCombatEngine(seed int64) *CombatEngine {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	return &CombatEngine{
		Seed: seed,
		rng:  rand.New(rand.NewSource(seed)),
	}
}

// ShipType represents hull classifications
type ShipType string

const (
	ShipTypeFrigate    ShipType = "frigate"
	ShipTypeCruiser    ShipType = "cruiser"
	ShipTypeBattleship ShipType = "battleship"
)

// DamageType represents weapon damage types
type DamageType string

const (
	DamageKinetic   DamageType = "kinetic"
	DamageExplosive DamageType = "explosive"
	DamageHeat      DamageType = "heat"
	DamageMagnetic  DamageType = "magnetic"
)

// WeaponCategory represents the tech-tree weapon category for bonus application
type WeaponCategory string

const (
	WeaponBallistic   WeaponCategory = "ballistic"
	WeaponDirectional WeaponCategory = "directional"
	WeaponMissile     WeaponCategory = "missile"
	WeaponFighter     WeaponCategory = "fighter"
)

// ArmorType represents armor classifications
type ArmorType string

const (
	ArmorChrome       ArmorType = "chrome"
	ArmorRegen        ArmorType = "regen"
	ArmorNano         ArmorType = "nano"
	ArmorNeutralizing ArmorType = "neutralizing"
)

// FleetStack represents a stack of ships in the 3x3 grid
type FleetStack struct {
	ID             string
	ShipDesignID   string
	GridRow        int
	GridCol        int
	ShipCount      int
	ShipType       ShipType
	DamageType     DamageType
	ArmorType      ArmorType
	WeaponCategory WeaponCategory
	BaseAttack     int
	BaseDefense    int
	BaseSpeed      int
	BaseAccuracy   int
	BaseDodge      int
	BaseShield     int
	BaseStructure  int
	BaseAgility    int
	// Effective stats (after bonuses)
	EffectiveStacks    int
	EffectiveAttack    int
	EffectiveDefense   int
	EffectiveSpeed     int
	EffectiveAccuracy  int
	EffectiveDodge     int
	EffectiveShield    int
	EffectiveStructure int
	// Current combat state
	CurrentShield    int
	CurrentStructure int
	CurrentShips     int

	// --- Advanced combat properties ---

	// Scatter/AoE: damage spreads to adjacent or all enemy ships
	ScatterDamage          float64 // % of damage dealt to adjacent stacks (by grid position)
	ScatterAll             float64 // % of damage spread across ALL enemy stacks
	ScatterRate            float64 // chance (0-1) that scatter triggers on a hit
	ScatterBonus           float64 // flat bonus added to scatter damage
	ScatterVsLowStructure  float64 // extra scatter % when target structure < 50%
	ScatterVsHighStructure float64 // extra scatter % when target structure >= 50%

	// Piercing: damage passes through the primary target to stacks behind it
	PiercingDamage      float64 // % of damage that pierces to the stack behind
	PiercingCritical    bool    // if true, piercing damage can crit
	PiercingDamageBonus float64 // multiplier applied to piercing damage (1.0 = no bonus)

	// Restoration: per-round healing
	ShieldRestore    float64 // % of max shield restored each round
	StructureRestore float64 // % of max structure restored each round
	AbsorbDouble     float64 // chance (0-1) to absorb 2x damage (halves effective damage)

	// Reflection: reflect damage back to attacker before HP drops to 0
	ReflectDamage          float64 // % of incoming damage reflected (from shield damage)
	ReflectStructureDamage float64 // % of incoming structure damage reflected

	// Positional: knockback and range-dependent damage
	Knockback   int              // push target back N grid rows on hit
	RangeDamage map[int]float64  // grid-row distance -> damage multiplier

	// Debuffs applied on hit
	DamageTakenIncrease  float64 // target takes X% more damage for rest of combat
	EnemyAttackReduction float64 // reduce target's attack by X% for EnemyAttackReductionRounds
	EnemyAttackReductionRounds int // number of rounds the attack reduction lasts
}

// Commander bonuses
type CommanderBonus struct {
	Accuracy        int
	Dodge           int
	Speed           int
	Electron        int
	StarRank        int
	WeaponExpertise string // ballistic, directional, missile, fighter (or empty)
	ShipExpertise   string // frigate, cruiser, battleship (or empty)
	EffectiveStack  float64 // Percentage bonus to effective stacks
}

// Fleet represents a complete fleet in combat
type Fleet struct {
	PlayerID       string
	FleetID        string
	CommanderBonus *CommanderBonus
	TechBonuses    *services.TechBonuses
	Stacks         []*FleetStack
	Formation      string
	Targeting      string
	Side           string // "attacker" or "defender"
}

// StackDebuff tracks a temporary debuff applied to a stack
type StackDebuff struct {
	AttackReduction float64 // % attack reduction
	RoundsLeft      int     // rounds remaining
}

// CombatState tracks the current state of combat
type CombatState struct {
	Attacker    *Fleet
	Defender    *Fleet
	Round       int
	MaxRounds   int
	Logs        []string
	RoundLogs   []*CombatRound
	IsComplete  bool
	Winner      string // "attacker", "defender", or "draw"

	// Advanced combat state
	DamageTakenModifiers map[string]float64      // stack ID -> cumulative % extra damage taken
	AttackDebuffs        map[string]*StackDebuff  // stack ID -> active attack reduction debuff
}

// CombatRound stores detailed information about a combat round
type CombatRound struct {
	RoundNumber int
	Attacks     []*Attack
	Casualties  map[string]int // stack_id -> ships destroyed
}

// Attack represents a single attack action
type Attack struct {
	AttackerStackID   string
	DefenderStackID   string
	AttackerSide      string
	DefenderSide      string
	Hit               bool
	CriticalHit       bool
	SuccessiveStrike  bool
	Damage            int
	ShieldDamage      int
	StructureDamage   int
	ShipsDestroyed    int

	// Advanced combat effects
	ScatterDamage     int            // total AoE damage dealt to other stacks
	ScatterTargets    map[string]int // stack ID -> scatter damage dealt
	PiercingDamage    int            // damage dealt to stack behind target
	PiercingTargetID  string         // stack that received piercing damage
	ReflectedDamage   int            // damage reflected back to attacker
	AbsorbedDouble    bool           // true if target absorbed 2x (halved incoming)
	KnockbackApplied  int            // rows the target was pushed back
}

// CombatResult is the final outcome of combat
type CombatResult struct {
	Winner              string
	AttackerCasualties  int
	DefenderCasualties  int
	TotalRounds         int
	Loot                *Loot
	CombatLog           []string
	DetailedRounds      []*CombatRound
}

// Loot represents resources looted from PvP
type Loot struct {
	Metal int64
	He3   int64
	Gold  int64
}

// ExecuteCombat runs the 8-phase combat engine
func (ce *CombatEngine) ExecuteCombat(attacker, defender *Fleet) (*CombatResult, error) {
	state := &CombatState{
		Attacker:             attacker,
		Defender:             defender,
		Round:                0,
		MaxRounds:            99,
		Logs:                 []string{},
		RoundLogs:            []*CombatRound{},
		DamageTakenModifiers: make(map[string]float64),
		AttackDebuffs:        make(map[string]*StackDebuff),
	}

	state.log(fmt.Sprintf("Combat started: %s vs %s", attacker.PlayerID, defender.PlayerID))
	state.log(fmt.Sprintf("Attacker formation: %s | Defender formation: %s", attacker.Formation, defender.Formation))

	// Phase 1: Calculate effective stacks (pre-combat, once)
	ce.phase1CalculateEffectiveStacks(state)

	// Combat loop
	for state.Round < state.MaxRounds {
		state.Round++
		state.log(fmt.Sprintf("=== Round %d ===", state.Round))

		round := &CombatRound{
			RoundNumber: state.Round,
			Attacks:     []*Attack{},
			Casualties:  make(map[string]int),
		}

		// Restoration phase: heal shields/structure at the start of each round
		ce.applyRestoration(state)

		// Tick down attack debuffs
		ce.tickDebuffs(state)

		// Phase 2: Ship type advantage calculation (embedded in damage calculation)
		// Phase 3: Determine attack order
		attackOrder := ce.phase3DetermineAttackOrder(state)

		// Execute attacks for each stack in order
		for _, stack := range attackOrder {
			if stack.CurrentShips <= 0 {
				continue
			}

			// Find target
			target := ce.selectTarget(state, stack)
			if target == nil || target.CurrentShips <= 0 {
				continue
			}

			// Get fleet to access commander bonuses
			attackerFleet := ce.getFleetForStack(state, stack)

			// Check for successive strike (Speed stat: Speed/500 chance to attack twice)
			attackCount := 1
			successiveStrike := false
			if attackerFleet != nil && attackerFleet.CommanderBonus != nil {
				successiveChance := float64(attackerFleet.CommanderBonus.Speed) / 500.0
				if ce.rng.Float64() < successiveChance {
					attackCount = 2
					successiveStrike = true
				}
			}

			// Execute attack(s) - can be 1 or 2 if successive strike triggers
			for attackNum := 0; attackNum < attackCount; attackNum++ {
				// Phase 4: Calculate hit chance (with tech bonuses)
				hitChance := ce.phase4CalculateHitChanceWithTech(stack, target, attackerFleet)
				hit := ce.rng.Float64() < hitChance

				attack := &Attack{
					AttackerStackID:  stack.ID,
					DefenderStackID:  target.ID,
					AttackerSide:     ce.getSide(state, stack),
					DefenderSide:     ce.getSide(state, target),
					Hit:              hit,
					SuccessiveStrike: successiveStrike && attackNum > 0,
				}

				if !hit {
					state.log(fmt.Sprintf("Stack %s missed target %s", stack.ID, target.ID))
					round.Attacks = append(round.Attacks, attack)
					continue
				}

				// Phase 5: Calculate damage with formation bonuses
				defenderFleet := ce.getFleetForStack(state, target)
				damage := ce.phase5CalculateDamageWithFormation(stack, target, attackerFleet, defenderFleet)

				// Apply attack debuff reduction to attacker
				if debuff, ok := state.AttackDebuffs[stack.ID]; ok && debuff.RoundsLeft > 0 {
					damage = int(math.Round(float64(damage) * (1.0 - debuff.AttackReduction/100.0)))
				}

				// Apply range-based damage modifier
				if stack.RangeDamage != nil {
					distance := abs(stack.GridRow - target.GridRow)
					if mult, ok := stack.RangeDamage[distance]; ok {
						damage = int(math.Round(float64(damage) * mult))
					}
				}

				// Apply DamageTakenIncrease modifier on the target
				if extraDmg, ok := state.DamageTakenModifiers[target.ID]; ok && extraDmg > 0 {
					damage = int(math.Round(float64(damage) * (1.0 + extraDmg/100.0)))
				}

				// Check for critical hit (Electron stat: 5% + Electron/200 + tech crit rate)
				criticalHit := false
				critMultiplier := 1.5
				{
					critChance := 0.05

					// Commander Electron bonus
					if attackerFleet != nil && attackerFleet.CommanderBonus != nil {
						critChance += float64(attackerFleet.CommanderBonus.Electron) / 200.0
					}

					// Tech crit rate bonus (weapon-specific)
					if attackerFleet != nil && attackerFleet.TechBonuses != nil {
						tb := attackerFleet.TechBonuses
						switch stack.WeaponCategory {
						case WeaponBallistic:
							critChance += tb.BallisticCritRate / 100.0
							critMultiplier += tb.BallisticCritDamage / 100.0
						case WeaponDirectional:
							critChance += tb.DirectionalCritRate / 100.0
						}
					}

					if ce.rng.Float64() < critChance {
						damage = int(math.Round(float64(damage) * critMultiplier))
						criticalHit = true
					}
				}

				// AbsorbDouble: chance to halve incoming damage
				if target.AbsorbDouble > 0 && ce.rng.Float64() < target.AbsorbDouble {
					damage = damage / 2
					attack.AbsorbedDouble = true
					state.log(fmt.Sprintf("Stack %s absorbed double (halved damage to %d)", target.ID, damage))
				}

				attack.Damage = damage
				attack.CriticalHit = criticalHit

				// Reflection: reflect damage back to attacker before applying
				if target.ReflectDamage > 0 || target.ReflectStructureDamage > 0 {
					reflectedDamage := ce.applyReflection(state, stack, target, damage, round)
					attack.ReflectedDamage = reflectedDamage
				}

				// Phase 6: Apply damage
				shieldDamage, structureDamage := ce.phase6ApplyDamage(target, damage)
				attack.ShieldDamage = shieldDamage
				attack.StructureDamage = structureDamage

				// Phase 7: Calculate casualties
				destroyed := ce.phase7CalculateCasualties(target)
				attack.ShipsDestroyed = destroyed
				round.Casualties[target.ID] += destroyed

				// Apply debuffs from attacker to target
				ce.applyDebuffs(state, stack, target)

				// Knockback: push target back N rows
				if stack.Knockback > 0 && target.CurrentShips > 0 {
					oldRow := target.GridRow
					target.GridRow += stack.Knockback
					if target.GridRow > 2 {
						target.GridRow = 2
					}
					if target.GridRow != oldRow {
						attack.KnockbackApplied = target.GridRow - oldRow
						state.log(fmt.Sprintf("Stack %s knocked %s back from row %d to row %d",
							stack.ID, target.ID, oldRow, target.GridRow))
					}
				}

				// Scatter/AoE damage
				if stack.ScatterDamage > 0 || stack.ScatterAll > 0 {
					scatterTotal, scatterTargets := ce.applyScatter(state, stack, target, damage, round)
					attack.ScatterDamage = scatterTotal
					attack.ScatterTargets = scatterTargets
				}

				// Piercing damage
				if stack.PiercingDamage > 0 {
					pierceDmg, pierceTargetID := ce.applyPiercing(state, stack, target, damage, criticalHit, critMultiplier, round)
					attack.PiercingDamage = pierceDmg
					attack.PiercingTargetID = pierceTargetID
				}

				critMsg := ""
				if criticalHit {
					critMsg = " [CRIT]"
				}
				strikeMsg := ""
				if successiveStrike && attackNum > 0 {
					strikeMsg = " [SUCCESSIVE STRIKE]"
				}
				absorbMsg := ""
				if attack.AbsorbedDouble {
					absorbMsg = " [ABSORB 2x]"
				}
				state.log(fmt.Sprintf("Stack %s hit %s for %d dmg%s%s%s (%d shield, %d structure), %d ships destroyed",
					stack.ID, target.ID, damage, critMsg, strikeMsg, absorbMsg, shieldDamage, structureDamage, destroyed))

				round.Attacks = append(round.Attacks, attack)

				// Check if target destroyed before second attack
				if target.CurrentShips <= 0 {
					break
				}
			}
		}

		state.RoundLogs = append(state.RoundLogs, round)

		// Check victory conditions
		if ce.checkVictory(state) {
			break
		}

		// Min 20 rounds
		if state.Round >= 20 && ce.isFleetDestroyed(state.Attacker) || ce.isFleetDestroyed(state.Defender) {
			break
		}
	}

	// Determine winner
	winner := ce.determineWinner(state)
	state.Winner = winner
	state.log(fmt.Sprintf("Combat ended. Winner: %s", winner))

	// Calculate result
	result := &CombatResult{
		Winner:             winner,
		AttackerCasualties: ce.countCasualties(state.Attacker),
		DefenderCasualties: ce.countCasualties(state.Defender),
		TotalRounds:        state.Round,
		CombatLog:          state.Logs,
		DetailedRounds:     state.RoundLogs,
	}

	// Phase 8: Loot (PvP only, if attacker wins)
	if winner == "attacker" {
		result.Loot = ce.phase8CalculateLoot(state)
	}

	return result, nil
}

// Phase 1: Calculate effective stacks (base ships + commander bonus + tech bonuses)
func (ce *CombatEngine) phase1CalculateEffectiveStacks(state *CombatState) {
	for _, fleet := range []*Fleet{state.Attacker, state.Defender} {
		for _, stack := range fleet.Stacks {
			// Base effective stacks = ship count
			effectiveStacks := float64(stack.ShipCount)

			// Apply commander effective stack bonus (percentage)
			if fleet.CommanderBonus != nil && fleet.CommanderBonus.EffectiveStack > 0 {
				effectiveStacks *= (1.0 + fleet.CommanderBonus.EffectiveStack/100.0)
			}

			stack.EffectiveStacks = int(math.Round(effectiveStacks))
			stack.CurrentShips = stack.ShipCount

			// Apply commander stat bonuses
			stack.EffectiveAccuracy = stack.BaseAccuracy
			stack.EffectiveDodge = stack.BaseDodge
			stack.EffectiveSpeed = stack.BaseSpeed
			stack.EffectiveAttack = stack.BaseAttack
			stack.EffectiveDefense = stack.BaseDefense
			stack.EffectiveShield = stack.BaseShield
			stack.EffectiveStructure = stack.BaseStructure

			if fleet.CommanderBonus != nil {
				stack.EffectiveAccuracy += fleet.CommanderBonus.Accuracy
				stack.EffectiveDodge += fleet.CommanderBonus.Dodge
				stack.EffectiveSpeed += fleet.CommanderBonus.Speed
				// Electron affects shield/structure
			}

			// Apply tech bonuses
			if fleet.TechBonuses != nil {
				stack.EffectiveShield = int(float64(stack.EffectiveShield) * (1.0 + fleet.TechBonuses.BaseShield/100.0))
				stack.EffectiveStructure = int(float64(stack.EffectiveStructure) * (1.0 + fleet.TechBonuses.BaseStructure/100.0))
				stack.EffectiveDefense = int(float64(stack.EffectiveDefense) * (1.0 + fleet.TechBonuses.BaseDefense/100.0))
				stack.EffectiveDodge = int(float64(stack.EffectiveDodge) * (1.0 + fleet.TechBonuses.BaseAgility/100.0))
			}

			// Set current combat values
			stack.CurrentShield = stack.EffectiveShield * stack.EffectiveStacks
			stack.CurrentStructure = stack.EffectiveStructure * stack.EffectiveStacks
		}
	}
}

// Phase 2: Ship type advantage (embedded in damage calculation)
// Frigate > Cruiser > Battleship > Frigate (+5%/-5%)
func (ce *CombatEngine) getShipTypeAdvantage(attacker, defender *FleetStack) float64 {
	if attacker.ShipType == ShipTypeFrigate && defender.ShipType == ShipTypeCruiser {
		return 1.05
	}
	if attacker.ShipType == ShipTypeCruiser && defender.ShipType == ShipTypeBattleship {
		return 1.05
	}
	if attacker.ShipType == ShipTypeBattleship && defender.ShipType == ShipTypeFrigate {
		return 1.05
	}
	if attacker.ShipType == ShipTypeCruiser && defender.ShipType == ShipTypeFrigate {
		return 0.95
	}
	if attacker.ShipType == ShipTypeBattleship && defender.ShipType == ShipTypeCruiser {
		return 0.95
	}
	if attacker.ShipType == ShipTypeFrigate && defender.ShipType == ShipTypeBattleship {
		return 0.95
	}
	return 1.0 // No advantage
}

// Phase 3: Determine attack order based on speed (highest first, random tiebreaker)
func (ce *CombatEngine) phase3DetermineAttackOrder(state *CombatState) []*FleetStack {
	allStacks := append([]*FleetStack{}, state.Attacker.Stacks...)
	allStacks = append(allStacks, state.Defender.Stacks...)

	// Simple bubble sort by speed (descending)
	for i := 0; i < len(allStacks)-1; i++ {
		for j := 0; j < len(allStacks)-i-1; j++ {
			// If speeds equal, randomize
			if allStacks[j].EffectiveSpeed == allStacks[j+1].EffectiveSpeed {
				if ce.rng.Intn(2) == 0 {
					allStacks[j], allStacks[j+1] = allStacks[j+1], allStacks[j]
				}
			} else if allStacks[j].EffectiveSpeed < allStacks[j+1].EffectiveSpeed {
				allStacks[j], allStacks[j+1] = allStacks[j+1], allStacks[j]
			}
		}
	}

	return allStacks
}

// Phase 4: Calculate hit chance (accuracy vs dodge, clamp 5%-95%)
func (ce *CombatEngine) phase4CalculateHitChance(attacker, defender *FleetStack) float64 {
	// Base hit chance = 50% + (accuracy - dodge) / 10
	hitChance := 0.5 + float64(attacker.EffectiveAccuracy-defender.EffectiveDodge)/1000.0

	// Clamp to 5%-95%
	if hitChance < 0.05 {
		hitChance = 0.05
	}
	if hitChance > 0.95 {
		hitChance = 0.95
	}

	return hitChance
}

// phase4CalculateHitChanceWithTech applies weapon-specific tech hit rate bonuses
func (ce *CombatEngine) phase4CalculateHitChanceWithTech(attacker, defender *FleetStack, attackerFleet *Fleet) float64 {
	hitChance := ce.phase4CalculateHitChance(attacker, defender)

	if attackerFleet != nil && attackerFleet.TechBonuses != nil {
		tb := attackerFleet.TechBonuses
		switch attacker.WeaponCategory {
		case WeaponBallistic:
			hitChance += tb.BallisticHitRate / 100.0
		case WeaponDirectional:
			hitChance += tb.DirectionalAccuracy / 100.0
		case WeaponMissile:
			hitChance += tb.MissileHitRate / 100.0
		case WeaponFighter:
			hitChance += tb.FighterHitRate / 100.0
		}
	}

	// Re-clamp after tech bonuses
	if hitChance < 0.05 {
		hitChance = 0.05
	}
	if hitChance > 0.95 {
		hitChance = 0.95
	}

	return hitChance
}

// Phase 5: Calculate damage (weapon damage * type advantage * armor effectiveness * position modifier * formation bonus)
func (ce *CombatEngine) phase5CalculateDamage(attacker, defender *FleetStack) int {
	baseDamage := float64(attacker.EffectiveAttack * attacker.EffectiveStacks)

	// Apply ship type advantage
	typeAdvantage := ce.getShipTypeAdvantage(attacker, defender)
	baseDamage *= typeAdvantage

	// Apply armor effectiveness
	armorEffectiveness := ce.getArmorEffectiveness(attacker.DamageType, defender.ArmorType)
	baseDamage *= armorEffectiveness

	// Apply position-based attack modifier (GO2 formula)
	positionModifier := ce.getPositionAttackModifier(attacker)
	baseDamage *= positionModifier

	return int(math.Round(baseDamage))
}

// Phase 5 with formation bonuses (called from combat loop with fleet context)
func (ce *CombatEngine) phase5CalculateDamageWithFormation(attacker, defender *FleetStack, attackerFleet, defenderFleet *Fleet) int {
	baseDamage := ce.phase5CalculateDamage(attacker, defender)

	// Apply weapon-category tech damage bonus
	if attackerFleet != nil && attackerFleet.TechBonuses != nil {
		tb := attackerFleet.TechBonuses
		switch attacker.WeaponCategory {
		case WeaponBallistic:
			baseDamage = int(math.Round(float64(baseDamage) * (1.0 + tb.BallisticDamage/100.0)))
		case WeaponDirectional:
			baseDamage = int(math.Round(float64(baseDamage) * (1.0 + tb.DirectionalDamage/100.0)))
		case WeaponMissile:
			baseDamage = int(math.Round(float64(baseDamage) * (1.0 + tb.MissileDamage/100.0)))
		case WeaponFighter:
			baseDamage = int(math.Round(float64(baseDamage) * (1.0 + tb.FighterDamage/100.0)))
		}
	}

	// Apply attacker formation attack bonus
	if attackerFleet != nil {
		attackBonus, _ := ce.getFormationBonuses(attackerFleet.Formation)
		baseDamage = int(math.Round(float64(baseDamage) * attackBonus))
	}

	// Apply defender formation defense bonus (reduces incoming damage)
	if defenderFleet != nil {
		_, defenseBonus := ce.getFormationBonuses(defenderFleet.Formation)
		// Defense bonus reduces damage: if defense is 1.10 (110%), damage is reduced by ~9%
		baseDamage = int(math.Round(float64(baseDamage) / defenseBonus))
	}

	return baseDamage
}

// Position-based attack modifier (GO2 formula: Front 100%, Middle 90%, Back 75%)
// Based on grid row position in 3x3 formation
func (ce *CombatEngine) getPositionAttackModifier(attacker *FleetStack) float64 {
	// GridRow: 0 = Front (top), 1 = Middle, 2 = Back (bottom)
	switch attacker.GridRow {
	case 0:
		return 1.00 // Front rank: 100% attack power
	case 1:
		return 0.90 // Middle rank: 90% attack power
	case 2:
		return 0.75 // Back rank: 75% attack power
	default:
		return 1.00 // Default to 100% if invalid position
	}
}

// Formation bonuses (GO2 formula: 7 formations with attack/defense modifiers)
// Returns (attackBonus, defenseBonus) as multipliers
func (ce *CombatEngine) getFormationBonuses(formation string) (attackBonus float64, defenseBonus float64) {
	switch formation {
	case "phalanx":
		return 1.00, 1.10 // 0% attack, +10% defense (balanced defensive)
	case "diamond":
		return 1.05, 1.05 // +5% attack, +5% defense (balanced)
	case "battle_line":
		return 1.10, 1.00 // +10% attack, 0% defense (aggressive)
	case "skirmish":
		return 1.15, 0.90 // +15% attack, -10% defense (high risk/reward)
	case "tee_forward":
		return 1.08, 1.02 // +8% attack, +2% defense (offensive focus)
	case "enfilade":
		return 1.12, 0.95 // +12% attack, -5% defense (flanking bonus)
	case "tee_reverse":
		return 0.95, 1.15 // -5% attack, +15% defense (defensive focus)
	default:
		return 1.00, 1.00 // No formation bonuses
	}
}

// Armor effectiveness matrix
func (ce *CombatEngine) getArmorEffectiveness(damageType DamageType, armorType ArmorType) float64 {
	// Chrome: weak to Explosive (1.25x), resists Kinetic (0.75x)
	if armorType == ArmorChrome {
		if damageType == DamageExplosive {
			return 1.25
		}
		if damageType == DamageKinetic {
			return 0.75
		}
	}

	// Regen: weak to Heat (1.25x), resists Explosive (0.75x)
	if armorType == ArmorRegen {
		if damageType == DamageHeat {
			return 1.25
		}
		if damageType == DamageExplosive {
			return 0.75
		}
	}

	// Nano: weak to Magnetic (1.25x), resists Heat (0.75x)
	if armorType == ArmorNano {
		if damageType == DamageMagnetic {
			return 1.25
		}
		if damageType == DamageHeat {
			return 0.75
		}
	}

	// Neutralizing: weak to Kinetic (1.25x), resists Magnetic (0.75x)
	if armorType == ArmorNeutralizing {
		if damageType == DamageKinetic {
			return 1.25
		}
		if damageType == DamageMagnetic {
			return 0.75
		}
	}

	return 1.0 // Neutral
}

// Phase 6: Apply damage (shields first, then structure)
func (ce *CombatEngine) phase6ApplyDamage(target *FleetStack, damage int) (shieldDamage, structureDamage int) {
	remainingDamage := damage

	// Apply to shield first
	if target.CurrentShield > 0 {
		shieldDamage = remainingDamage
		if shieldDamage > target.CurrentShield {
			shieldDamage = target.CurrentShield
		}
		target.CurrentShield -= shieldDamage
		remainingDamage -= shieldDamage
	}

	// Apply remainder to structure
	if remainingDamage > 0 && target.CurrentStructure > 0 {
		structureDamage = remainingDamage
		if structureDamage > target.CurrentStructure {
			structureDamage = target.CurrentStructure
		}
		target.CurrentStructure -= structureDamage
	}

	return shieldDamage, structureDamage
}

// Phase 7: Calculate casualties (ships destroyed when structure reaches 0)
func (ce *CombatEngine) phase7CalculateCasualties(target *FleetStack) int {
	if target.CurrentStructure <= 0 && target.CurrentShips > 0 {
		// Entire stack destroyed
		destroyed := target.CurrentShips
		target.CurrentShips = 0
		target.CurrentShield = 0
		target.CurrentStructure = 0
		return destroyed
	}

	// Partial casualties based on structure damage
	if target.EffectiveStructure > 0 {
		totalStructure := target.EffectiveStructure * target.EffectiveStacks
		damagePercent := 1.0 - (float64(target.CurrentStructure) / float64(totalStructure))
		casualtyCount := int(math.Floor(float64(target.CurrentShips) * damagePercent))

		if casualtyCount > target.CurrentShips {
			casualtyCount = target.CurrentShips
		}

		target.CurrentShips -= casualtyCount
		return casualtyCount
	}

	return 0
}

// Phase 8: Calculate loot (PvP only, 20% of resources)
func (ce *CombatEngine) phase8CalculateLoot(state *CombatState) *Loot {
	// TODO: Query defender planet resources and calculate 20% loot
	// For now, return empty loot
	return &Loot{
		Metal: 0,
		He3:   0,
		Gold:  0,
	}
}

// --- Advanced combat subsystem methods ---

// applyRestoration heals shields and structure for all alive stacks at the start of each round.
func (ce *CombatEngine) applyRestoration(state *CombatState) {
	for _, fleet := range []*Fleet{state.Attacker, state.Defender} {
		for _, stack := range fleet.Stacks {
			if stack.CurrentShips <= 0 {
				continue
			}

			if stack.ShieldRestore > 0 {
				maxShield := stack.EffectiveShield * stack.EffectiveStacks
				restore := int(math.Round(float64(maxShield) * stack.ShieldRestore / 100.0))
				stack.CurrentShield += restore
				if stack.CurrentShield > maxShield {
					stack.CurrentShield = maxShield
				}
				if restore > 0 {
					state.log(fmt.Sprintf("Stack %s restored %d shield (%.1f%%)", stack.ID, restore, stack.ShieldRestore))
				}
			}

			if stack.StructureRestore > 0 {
				maxStructure := stack.EffectiveStructure * stack.EffectiveStacks
				restore := int(math.Round(float64(maxStructure) * stack.StructureRestore / 100.0))
				stack.CurrentStructure += restore
				if stack.CurrentStructure > maxStructure {
					stack.CurrentStructure = maxStructure
				}
				if restore > 0 {
					state.log(fmt.Sprintf("Stack %s restored %d structure (%.1f%%)", stack.ID, restore, stack.StructureRestore))
				}
			}
		}
	}
}

// tickDebuffs decrements round counters on attack reduction debuffs and removes expired ones.
func (ce *CombatEngine) tickDebuffs(state *CombatState) {
	for id, debuff := range state.AttackDebuffs {
		debuff.RoundsLeft--
		if debuff.RoundsLeft <= 0 {
			delete(state.AttackDebuffs, id)
			state.log(fmt.Sprintf("Stack %s attack debuff expired", id))
		}
	}
}

// applyDebuffs applies DamageTakenIncrease and EnemyAttackReduction from attacker to target.
func (ce *CombatEngine) applyDebuffs(state *CombatState, attacker, target *FleetStack) {
	if attacker.DamageTakenIncrease > 0 {
		state.DamageTakenModifiers[target.ID] += attacker.DamageTakenIncrease
		state.log(fmt.Sprintf("Stack %s debuffed %s: +%.1f%% damage taken (total: +%.1f%%)",
			attacker.ID, target.ID, attacker.DamageTakenIncrease, state.DamageTakenModifiers[target.ID]))
	}

	if attacker.EnemyAttackReduction > 0 && attacker.EnemyAttackReductionRounds > 0 {
		state.AttackDebuffs[target.ID] = &StackDebuff{
			AttackReduction: attacker.EnemyAttackReduction,
			RoundsLeft:      attacker.EnemyAttackReductionRounds,
		}
		state.log(fmt.Sprintf("Stack %s reduced %s attack by %.1f%% for %d rounds",
			attacker.ID, target.ID, attacker.EnemyAttackReduction, attacker.EnemyAttackReductionRounds))
	}
}

// applyReflection reflects a portion of incoming damage back to the attacker.
// Returns the total amount of reflected damage.
func (ce *CombatEngine) applyReflection(state *CombatState, attacker, target *FleetStack, incomingDamage int, round *CombatRound) int {
	reflectedTotal := 0

	// Reflect shield damage portion
	if target.ReflectDamage > 0 {
		reflected := int(math.Round(float64(incomingDamage) * target.ReflectDamage / 100.0))
		if reflected > 0 {
			ce.phase6ApplyDamage(attacker, reflected)
			ce.phase7CalculateCasualties(attacker)
			reflectedTotal += reflected
			state.log(fmt.Sprintf("Stack %s reflected %d damage back to %s", target.ID, reflected, attacker.ID))
		}
	}

	// Reflect structure damage portion
	if target.ReflectStructureDamage > 0 {
		reflected := int(math.Round(float64(incomingDamage) * target.ReflectStructureDamage / 100.0))
		if reflected > 0 {
			// Structure reflection bypasses shields, applies directly to structure
			if attacker.CurrentStructure > 0 {
				dmg := reflected
				if dmg > attacker.CurrentStructure {
					dmg = attacker.CurrentStructure
				}
				attacker.CurrentStructure -= dmg
				ce.phase7CalculateCasualties(attacker)
				reflectedTotal += dmg
				state.log(fmt.Sprintf("Stack %s reflected %d structure damage back to %s", target.ID, dmg, attacker.ID))
			}
		}
	}

	return reflectedTotal
}

// applyScatter deals AoE damage to adjacent or all enemy stacks.
// Returns total scatter damage dealt and a map of target stack IDs to damage dealt.
func (ce *CombatEngine) applyScatter(state *CombatState, attacker, primaryTarget *FleetStack, primaryDamage int, round *CombatRound) (int, map[string]int) {
	// Check scatter rate (chance to trigger)
	if attacker.ScatterRate > 0 && ce.rng.Float64() >= attacker.ScatterRate {
		return 0, nil
	}

	var enemyStacks []*FleetStack
	if ce.isAttackerStack(state, attacker) {
		enemyStacks = state.Defender.Stacks
	} else {
		enemyStacks = state.Attacker.Stacks
	}

	scatterTargets := make(map[string]int)
	totalScatter := 0

	// Determine conditional scatter bonus based on target structure
	conditionalBonus := 0.0
	if primaryTarget.EffectiveStructure > 0 && primaryTarget.EffectiveStacks > 0 {
		maxStruct := float64(primaryTarget.EffectiveStructure * primaryTarget.EffectiveStacks)
		structPercent := float64(primaryTarget.CurrentStructure) / maxStruct
		if structPercent < 0.5 && attacker.ScatterVsLowStructure > 0 {
			conditionalBonus = attacker.ScatterVsLowStructure
		} else if structPercent >= 0.5 && attacker.ScatterVsHighStructure > 0 {
			conditionalBonus = attacker.ScatterVsHighStructure
		}
	}

	for _, enemyStack := range enemyStacks {
		if enemyStack.ID == primaryTarget.ID || enemyStack.CurrentShips <= 0 {
			continue
		}

		var scatterPercent float64

		if attacker.ScatterAll > 0 {
			// ScatterAll: spread to ALL enemy stacks
			scatterPercent = attacker.ScatterAll + conditionalBonus
		} else if attacker.ScatterDamage > 0 {
			// ScatterDamage: only adjacent stacks (within 1 row/col distance)
			rowDist := abs(primaryTarget.GridRow - enemyStack.GridRow)
			colDist := abs(primaryTarget.GridCol - enemyStack.GridCol)
			if rowDist <= 1 && colDist <= 1 {
				scatterPercent = attacker.ScatterDamage + conditionalBonus
			}
		}

		if scatterPercent <= 0 {
			continue
		}

		scatterDmg := int(math.Round(float64(primaryDamage)*scatterPercent/100.0 + attacker.ScatterBonus))
		if scatterDmg <= 0 {
			continue
		}

		ce.phase6ApplyDamage(enemyStack, scatterDmg)
		destroyed := ce.phase7CalculateCasualties(enemyStack)
		round.Casualties[enemyStack.ID] += destroyed
		scatterTargets[enemyStack.ID] = scatterDmg
		totalScatter += scatterDmg

		state.log(fmt.Sprintf("Scatter: %s dealt %d AoE damage to %s, %d ships destroyed",
			attacker.ID, scatterDmg, enemyStack.ID, destroyed))
	}

	return totalScatter, scatterTargets
}

// applyPiercing deals damage through the primary target to the stack behind it (higher grid row).
// Returns piercing damage dealt and the target stack ID.
func (ce *CombatEngine) applyPiercing(state *CombatState, attacker, primaryTarget *FleetStack, primaryDamage int, wasCrit bool, critMultiplier float64, round *CombatRound) (int, string) {
	var enemyStacks []*FleetStack
	if ce.isAttackerStack(state, attacker) {
		enemyStacks = state.Defender.Stacks
	} else {
		enemyStacks = state.Attacker.Stacks
	}

	// Find the stack "behind" the primary target (next higher grid row)
	var behindStack *FleetStack
	bestRow := -1
	for _, s := range enemyStacks {
		if s.ID == primaryTarget.ID || s.CurrentShips <= 0 {
			continue
		}
		if s.GridRow > primaryTarget.GridRow {
			if bestRow == -1 || s.GridRow < bestRow {
				bestRow = s.GridRow
				behindStack = s
			}
		}
	}

	if behindStack == nil {
		return 0, ""
	}

	pierceDmg := int(math.Round(float64(primaryDamage) * attacker.PiercingDamage / 100.0))

	// Apply piercing damage bonus multiplier
	if attacker.PiercingDamageBonus > 0 {
		pierceDmg = int(math.Round(float64(pierceDmg) * attacker.PiercingDamageBonus))
	}

	// Apply crit to piercing if enabled and the primary attack was a crit
	if attacker.PiercingCritical && wasCrit {
		pierceDmg = int(math.Round(float64(pierceDmg) * critMultiplier))
	}

	if pierceDmg <= 0 {
		return 0, ""
	}

	ce.phase6ApplyDamage(behindStack, pierceDmg)
	destroyed := ce.phase7CalculateCasualties(behindStack)
	round.Casualties[behindStack.ID] += destroyed

	state.log(fmt.Sprintf("Piercing: %s dealt %d piercing damage to %s (behind %s), %d ships destroyed",
		attacker.ID, pierceDmg, behindStack.ID, primaryTarget.ID, destroyed))

	return pierceDmg, behindStack.ID
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Helper methods

func (ce *CombatEngine) selectTarget(state *CombatState, attacker *FleetStack) *FleetStack {
	var targets []*FleetStack

	// Determine enemy fleet
	if ce.isAttackerStack(state, attacker) {
		targets = state.Defender.Stacks
	} else {
		targets = state.Attacker.Stacks
	}

	// Filter alive targets
	var aliveTargets []*FleetStack
	for _, t := range targets {
		if t.CurrentShips > 0 {
			aliveTargets = append(aliveTargets, t)
		}
	}

	if len(aliveTargets) == 0 {
		return nil
	}

	// Honour the attacking fleet's targeting_command. The valid set comes from
	// handlers/fleets.go (max_attack, min_attack, max_durability, min_durability,
	// closest, by_commander_rank).
	fleet := ce.getFleetForStack(state, attacker)
	cmd := ""
	if fleet != nil {
		cmd = fleet.Targeting
	}

	score := func(t *FleetStack) int {
		switch cmd {
		case "max_attack":
			return t.EffectiveAttack
		case "min_attack":
			return -t.EffectiveAttack
		case "max_durability":
			return t.CurrentShield + t.CurrentStructure
		case "min_durability":
			return -(t.CurrentShield + t.CurrentStructure)
		case "closest":
			return -abs(attacker.GridRow - t.GridRow)
		default:
			return 0
		}
	}

	best := aliveTargets[0]
	bestScore := score(best)
	for _, t := range aliveTargets[1:] {
		if s := score(t); s > bestScore {
			best = t
			bestScore = s
		}
	}
	return best
}

func (ce *CombatEngine) isAttackerStack(state *CombatState, stack *FleetStack) bool {
	for _, s := range state.Attacker.Stacks {
		if s.ID == stack.ID {
			return true
		}
	}
	return false
}

func (ce *CombatEngine) getFleetForStack(state *CombatState, stack *FleetStack) *Fleet {
	if ce.isAttackerStack(state, stack) {
		return state.Attacker
	}
	return state.Defender
}

func (ce *CombatEngine) getSide(state *CombatState, stack *FleetStack) string {
	if ce.isAttackerStack(state, stack) {
		return "attacker"
	}
	return "defender"
}

func (ce *CombatEngine) checkVictory(state *CombatState) bool {
	attackerAlive := !ce.isFleetDestroyed(state.Attacker)
	defenderAlive := !ce.isFleetDestroyed(state.Defender)

	return !attackerAlive || !defenderAlive
}

func (ce *CombatEngine) isFleetDestroyed(fleet *Fleet) bool {
	for _, stack := range fleet.Stacks {
		if stack.CurrentShips > 0 {
			return false
		}
	}
	return true
}

func (ce *CombatEngine) determineWinner(state *CombatState) string {
	attackerAlive := !ce.isFleetDestroyed(state.Attacker)
	defenderAlive := !ce.isFleetDestroyed(state.Defender)

	if attackerAlive && !defenderAlive {
		return "attacker"
	}
	if defenderAlive && !attackerAlive {
		return "defender"
	}

	// If both alive or both dead, compare total remaining HP
	attackerHP := ce.getTotalHP(state.Attacker)
	defenderHP := ce.getTotalHP(state.Defender)

	if attackerHP > defenderHP {
		return "attacker"
	}
	if defenderHP > attackerHP {
		return "defender"
	}

	return "draw"
}

func (ce *CombatEngine) getTotalHP(fleet *Fleet) int {
	total := 0
	for _, stack := range fleet.Stacks {
		total += stack.CurrentShield + stack.CurrentStructure
	}
	return total
}

func (ce *CombatEngine) countCasualties(fleet *Fleet) int {
	casualties := 0
	for _, stack := range fleet.Stacks {
		casualties += (stack.ShipCount - stack.CurrentShips)
	}
	return casualties
}

func (state *CombatState) log(message string) {
	state.Logs = append(state.Logs, message)
	log.Println(message)
}
