package combat

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
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
	ID            string
	ShipDesignID  string
	GridRow       int
	GridCol       int
	ShipCount     int
	ShipType      ShipType
	DamageType    DamageType
	ArmorType     ArmorType
	BaseAttack    int
	BaseDefense   int
	BaseSpeed     int
	BaseAccuracy  int
	BaseDodge     int
	BaseShield    int
	BaseStructure int
	BaseAgility   int
	// Effective stats (after bonuses)
	EffectiveStacks   int
	EffectiveAttack   int
	EffectiveDefense  int
	EffectiveSpeed    int
	EffectiveAccuracy int
	EffectiveDodge    int
	EffectiveShield   int
	EffectiveStructure int
	// Current combat state
	CurrentShield    int
	CurrentStructure int
	CurrentShips     int
}

// Commander bonuses
type CommanderBonus struct {
	Accuracy      int
	Dodge         int
	Speed         int
	Electron      int
	EffectiveStack float64 // Percentage bonus to effective stacks
}

// Fleet represents a complete fleet in combat
type Fleet struct {
	PlayerID       string
	FleetID        string
	CommanderBonus *CommanderBonus
	TechBonuses    *TechBonuses
	Stacks         []*FleetStack
	Formation      string
	Targeting      string
	Side           string // "attacker" or "defender"
}

// TechBonuses represents technology bonuses
type TechBonuses struct {
	BallisticDamage     float64
	BallisticCritRate   float64
	BallisticCritDamage float64
	BallisticHitRate    float64
	DirectionalDamage   float64
	DirectionalCritRate float64
	DirectionalAccuracy float64
	MissileDamage       float64
	MissileHitRate      float64
	BaseShield          float64
	BaseStructure       float64
	BaseAgility         float64
	BaseDefense         float64
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
}

// CombatRound stores detailed information about a combat round
type CombatRound struct {
	RoundNumber int
	Attacks     []*Attack
	Casualties  map[string]int // stack_id -> ships destroyed
}

// Attack represents a single attack action
type Attack struct {
	AttackerStackID string
	DefenderStackID string
	AttackerSide    string
	DefenderSide    string
	Hit             bool
	Damage          int
	ShieldDamage    int
	StructureDamage int
	ShipsDestroyed  int
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
		Attacker:  attacker,
		Defender:  defender,
		Round:     0,
		MaxRounds: 99,
		Logs:      []string{},
		RoundLogs: []*CombatRound{},
	}

	state.log(fmt.Sprintf("Combat started: %s vs %s", attacker.PlayerID, defender.PlayerID))

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

			// Phase 4: Calculate hit chance
			hitChance := ce.phase4CalculateHitChance(stack, target)
			hit := ce.rng.Float64() < hitChance

			attack := &Attack{
				AttackerStackID: stack.ID,
				DefenderStackID: target.ID,
				AttackerSide:    ce.getSide(state, stack),
				DefenderSide:    ce.getSide(state, target),
				Hit:             hit,
			}

			if !hit {
				state.log(fmt.Sprintf("Stack %s missed target %s", stack.ID, target.ID))
				round.Attacks = append(round.Attacks, attack)
				continue
			}

			// Phase 5: Calculate damage
			damage := ce.phase5CalculateDamage(stack, target)
			attack.Damage = damage

			// Phase 6: Apply damage
			shieldDamage, structureDamage := ce.phase6ApplyDamage(target, damage)
			attack.ShieldDamage = shieldDamage
			attack.StructureDamage = structureDamage

			// Phase 7: Calculate casualties
			destroyed := ce.phase7CalculateCasualties(target)
			attack.ShipsDestroyed = destroyed
			round.Casualties[target.ID] = destroyed

			state.log(fmt.Sprintf("Stack %s hit %s for %d dmg (%d shield, %d structure), %d ships destroyed",
				stack.ID, target.ID, damage, shieldDamage, structureDamage, destroyed))

			round.Attacks = append(round.Attacks, attack)
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

			// Apply tech bonuses (simplified for now)
			if fleet.TechBonuses != nil {
				stack.EffectiveShield = int(float64(stack.EffectiveShield) * (1.0 + fleet.TechBonuses.BaseShield/100.0))
				stack.EffectiveStructure = int(float64(stack.EffectiveStructure) * (1.0 + fleet.TechBonuses.BaseStructure/100.0))
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

// Phase 5: Calculate damage (weapon damage * type advantage * armor effectiveness)
func (ce *CombatEngine) phase5CalculateDamage(attacker, defender *FleetStack) int {
	baseDamage := float64(attacker.EffectiveAttack * attacker.EffectiveStacks)

	// Apply ship type advantage
	typeAdvantage := ce.getShipTypeAdvantage(attacker, defender)
	baseDamage *= typeAdvantage

	// Apply armor effectiveness
	armorEffectiveness := ce.getArmorEffectiveness(attacker.DamageType, defender.ArmorType)
	baseDamage *= armorEffectiveness

	return int(math.Round(baseDamage))
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

	// Simple targeting: pick first alive target (TODO: implement targeting strategies)
	return aliveTargets[0]
}

func (ce *CombatEngine) isAttackerStack(state *CombatState, stack *FleetStack) bool {
	for _, s := range state.Attacker.Stacks {
		if s.ID == stack.ID {
			return true
		}
	}
	return false
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
