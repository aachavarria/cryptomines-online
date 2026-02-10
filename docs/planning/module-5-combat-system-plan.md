# Module 5: Combat System (8-Phase) - Implementation Plan

**Date:** 2026-02-07
**Architect:** architect
**Module:** Combat System - 8-Phase Combat Resolution Engine
**Estimated Total:** 40-55 hours (6-9 days)

---

## 1. EXECUTIVE SUMMARY

### 1.1 Module Overview

The Combat System is the **MOST CRITICAL module** in Cryptomines Online. This is the core game engine that makes all military features functional - without combat, ships/commanders/instances are meaningless. This module implements the complete Galaxy Online 2 combat mechanics with **8-phase round-by-round resolution**.

**Complexity Level:** HIGHEST (★★★★★)
- Most complex algorithms (damage calc, armor effectiveness, targeting)
- Highest integration surface (ships, commanders, techs, instances, PvP)
- Largest testing requirement (combat balance is game-critical)

### 1.2 Scope (From final-scope.md)

**IN SCOPE:**
- ✅ 8-phase combat resolution (Positioning → Targeting → Damage → Shields → Armor → HP → Casualties → Loot)
- ✅ Ship type advantage (Frigate > Cruiser > Battleship > Frigate rock-paper-scissors)
- ✅ Armor types vs Damage types effectiveness matrix (4x4 matrix)
- ✅ Commander bonuses (accuracy, dodge, speed, electron, effective stack)
- ✅ Tech bonuses (ballistic damage, missile hit rate, etc.)
- ✅ Combat reports (round-by-round detailed log)
- ✅ Instance combat integration (replace placeholder)
- ✅ Min 20 rounds, max 99 rounds
- ✅ Formation bonuses (7 formations)
- ✅ Targeting commands (6 strategies)
- ✅ Critical hits & successive strikes
- ✅ Shield penetration mechanics
- ✅ He3 fuel consumption

**OUT OF SCOPE (deferred to Module 6: PvP Combat):**
- ❌ PvP attack mechanics (fleet travel, SP consumption)
- ❌ Defense fleet integration
- ❌ Planetary defense buildings (Module 7)
- ❌ Loot 20% resources from defender
- ❌ Radar warnings
- ❌ Truce cards

**Combat Types in Module 5:**
- `instance_normal` - 30 Normal Instances (PRIMARY FOCUS)
- Foundation for future: `pvp`, `instance_restricted`, `league`, `rbp_attack`

### 1.3 Current Implementation Status

**Database:** 60% Complete
- ✅ `combat_reports` table exists (Phase 2 migration)
- ✅ `instances` table exists with placeholder combat
- ✅ `hull_types` with armor types (nano/chrome/neutralizing/regen/light)
- ✅ `module_types` with damage types (heat/kinetic/magnetic/explosive)
- ✅ `fleets` with formation + targeting_command
- ❌ NO armor effectiveness matrix defined
- ❌ NO combat formulas implemented

**Backend:** 0% Complete
- ❌ NO combat engine (combat_engine.go)
- ❌ NO 8-phase resolution logic
- ❌ NO damage calculation formulas
- ❌ NO armor effectiveness system
- ❌ NO hit chance calculations
- ❌ NO critical hit/successive strike logic
- ❌ Instance combat is placeholder (mock victory)

**Frontend:** 10% Complete
- ✅ Combat reports table placeholder (shows mock data)
- ❌ NO round-by-round visualization
- ❌ NO combat animation/effects
- ❌ NO detailed damage breakdown

**Integration Points:**
- Ships (attack power, defense, agility)
- Commanders (accuracy, dodge, speed, electron, effective stack)
- Tech bonuses (weapon damage, hit rate, defense)
- Formations (7 types with bonuses)
- Targeting commands (6 strategies)

---

## 2. DATABASE ANALYSIS

### 2.1 Existing Schema Review

**Table: `combat_reports` (existing)**
```sql
CREATE TABLE combat_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attacker_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    defender_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    attacker_fleet_id UUID,
    defender_fleet_id UUID,
    combat_type TEXT NOT NULL DEFAULT 'pvp'
        CHECK (combat_type IN ('pvp', 'instance_normal', 'instance_restricted',
               'instance_trial', 'instance_constellation', 'instance_humaroid',
               'league', 'championship', 'rbp_attack')),
    result TEXT NOT NULL CHECK (result IN ('attacker_win', 'defender_win', 'draw')),
    total_rounds INTEGER NOT NULL DEFAULT 0,
    rounds_json JSONB NOT NULL DEFAULT '[]',
    loot_json JSONB NOT NULL DEFAULT '{}',
    attacker_losses_json JSONB NOT NULL DEFAULT '{}',
    defender_losses_json JSONB NOT NULL DEFAULT '{}',
    he3_consumed BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

**✅ Perfect** - Already supports all required fields.

**`rounds_json` Structure:**
```json
[
  {
    "round": 1,
    "attacker_damage": 45000,
    "defender_damage": 38000,
    "attacker_ships_lost": 120,
    "defender_ships_lost": 200,
    "attacker_remaining": 2680,
    "defender_remaining": 2300,
    "events": [
      "Attacker Fleet Alpha fires: 450 hits",
      "Defender interceptors block 80 hits",
      "Critical hit! +50% damage",
      "Defender loses 200 ships"
    ]
  }
]
```

**Table: `hull_types` (existing - Phase 2)**
```sql
hull_types:
  armor_type TEXT NOT NULL -- 'nano', 'chrome', 'neutralizing', 'regen', 'light'
  base_shield INTEGER
  base_structure INTEGER
  base_defense NUMERIC(5,2)
  base_agility INTEGER
```

**Table: `module_types` (existing - Phase 2)**
```sql
module_types:
  damage_type TEXT -- 'heat', 'kinetic', 'magnetic', 'explosive', NULL (for non-weapons)
  min_damage INTEGER
  max_damage INTEGER
  weapon_range_min INTEGER
  weapon_range_max INTEGER
  cooldown INTEGER
  he3_per_round INTEGER
```

**Table: `fleets` (existing - Phase 2)**
```sql
fleets:
  formation TEXT NOT NULL DEFAULT 'phalanx'
      CHECK (formation IN ('phalanx', 'diamond', 'battle_line', 'skirmish',
             'tee_forward', 'enfilade', 'tee_reverse'))
  commander_id UUID REFERENCES commanders(id) ON DELETE SET NULL
  targeting_command TEXT NOT NULL DEFAULT 'max_attack'
      CHECK (targeting_command IN ('max_attack', 'min_attack', 'max_durability',
             'min_durability', 'closest', 'by_commander_rank'))
```

**Table: `commanders` (existing - Phase 2)**
```sql
commanders:
  accuracy INTEGER
  dodge INTEGER
  speed INTEGER
  electron INTEGER
  effective_stack INTEGER (300 + star_rank * 50)
  weapon_expertise JSONB -- {"ballistic":"A","directional":"B","missile":"B","ship_based":"B"}
  ship_expertise JSONB -- {"frigate":"B","cruiser":"B","battleship":"B"}
```

### 2.2 Required Schema Changes

**NONE REQUIRED** ✅

The existing schema is sufficient. We only need to implement the combat engine logic and populate `combat_reports.rounds_json` correctly.

**Optional Enhancement (future):**
```sql
-- Table for combat simulation debugging (NOT for Module 5)
CREATE TABLE combat_debug_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    combat_report_id UUID REFERENCES combat_reports(id) ON DELETE CASCADE,
    round_number INTEGER NOT NULL,
    phase TEXT NOT NULL, -- 'positioning', 'targeting', 'damage', etc.
    log_data JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

---

## 3. GALAXY ONLINE 2 COMBAT MECHANICS RESEARCH

### 3.1 Core Combat Structure (7-8 Phases)

Based on Galaxy Online 2 wiki research, combat follows these steps:

**Phase 1: Fleet Positioning**
- Determine turn order based on Speed stat (commander + ship agility)
- Formation bonuses applied
- Effective stack calculated per grid slot (base 1100/1000/900 for F/C/B)

**Phase 2: Targeting Selection**
- Apply targeting command (max_attack, min_attack, max_durability, etc.)
- Select target stack from enemy formation
- Consider weapon range constraints

**Phase 3: Attack Calculation**
- Calculate total attacks = effective stack × weapons per ship
- Commander accuracy bonus applied
- Tech bonuses (ballistic damage, hit rate, etc.)

**Phase 4: Interception** (optional)
- Defender's interceptor modules fire
- 55% intercept rate per PPC (Particle Pulse Cannon) against single targets
- Reduces incoming hits

**Phase 5: Hit Chance & Damage Roll**
- Hit chance = attacker hit rate / defender agility
- Roll for each attack: random(min_damage, max_damage)
- Critical hits (electron stat)
- Successive strikes (speed stat)

**Phase 6: Shield Absorption**
- Shields absorb damage before hull
- Shield penetration for ballistic/directional weapons
- EOS shields: 30% chance for double absorption

**Phase 7: Armor Effectiveness**
- Apply armor type vs damage type matrix
- Reduce/amplify damage based on matchup

**Phase 8: HP Reduction & Casualties**
- Damage applied to Structure (HP)
- Ships destroyed when structure reaches 0
- Scatter damage to other stacks (unmitigated bonus damage)

**Loot Calculation** (after all rounds):
- Instance loot from loot tables
- PvP: 20% of defender resources (Module 6)

### 3.2 Armor Effectiveness Matrix

From GO2 wiki research:

| Armor Type     | vs Heat | vs Kinetic | vs Magnetic | vs Explosive |
|----------------|---------|------------|-------------|--------------|
| **Chrome**     | 0.5×    | 0.5×       | 1.5×        | 2.0×         |
| **Nano**       | 0.5×    | 1.5×       | 2.0×        | 0.5×         |
| **Neutralizing** | 1.5×  | 2.0×       | 0.5×        | 0.5×         |
| **Regen**      | 2.0×    | 0.5×       | 0.5×        | 1.5×         |
| **Light**      | 1.0×    | 1.0×       | 1.0×        | 1.0×         |

**Interpretation:**
- 0.5× = Strong resistance (50% damage)
- 1.0× = Neutral (100% damage)
- 1.5× = Weak resistance (150% damage)
- 2.0× = Very weak (200% damage)

**Pattern:**
- Each armor (except Light) is strong against 2 damage types and weak against 2
- Chrome: Strong vs Heat/Kinetic, Weak vs Magnetic/Explosive
- Nano: Strong vs Heat/Explosive, Weak vs Kinetic/Magnetic
- Neutralizing: Strong vs Magnetic/Explosive, Weak vs Heat/Kinetic
- Regen: Strong vs Kinetic/Magnetic, Weak vs Heat/Explosive
- Light: Balanced (no bonuses/penalties)

### 3.3 Ship Type Advantages

**Base Effective Stack:**
- Frigate: 1,100 ships per stack
- Cruiser: 1,000 ships per stack
- Battleship: 900 ships per stack

**Type Advantage (Rock-Paper-Scissors):**
- Frigate > Cruiser (+20% damage)
- Cruiser > Battleship (+20% damage)
- Battleship > Frigate (+20% damage)

**Why Frigates have higher stack:**
- Smaller ships = more units per formation slot
- More attacks per round
- Less HP per ship (balanced)

### 3.4 Hit Chance Formula

From GO2 wiki:

```
Hit Chance = (Attacker Hit Rate + Accuracy Bonus) / Defender Agility
```

Where:
- Attacker Hit Rate = Base weapon hit rate (e.g., 85%)
- Accuracy Bonus = Commander accuracy + Tech accuracy bonuses
- Defender Agility = Ship agility + Commander dodge + Tech bonuses

**Hit chance can exceed 100%** (not capped, but displayed as 100% in UI).

**Agility Penalty:** Roughly 4% hit reduction per agility point.

### 3.5 Critical Hits & Successive Strikes

**Critical Hits:**
- Trigger chance based on Electron stat
- Effect: +50% damage on that hit
- Formula: `Crit Chance = Base Crit Rate + (Electron / 200)`
- Base Crit Rate: 5%

**Successive Strikes:**
- Trigger chance based on Speed stat
- Effect: Attack twice in same round
- Formula: `Successive Strike Chance = Speed / 500`
- Max: 20% chance at 100 speed

### 3.6 Commander Bonuses

**Accuracy:** Increases hit chance directly (+1% per point)
**Dodge:** Increases effective agility (+1 agility per point)
**Speed:** Increases turn order + successive strike chance
**Electron:** Increases crit rate + crit damage

**Effective Stack Multiplier:**
```
Effective Stack = Base Stack × (Commander Effective Stack / 300)
```

Example:
- Frigate base: 1,100
- Commander: Star Rank 5 = 550 effective stack
- Result: 1,100 × (550/300) = 2,017 ships per stack

**This is the biggest combat upgrade** - doubling ship count per slot.

### 3.7 Formation Bonuses

| Formation     | Attack Bonus | Defense Bonus | Special Effect |
|---------------|--------------|---------------|----------------|
| Phalanx       | 0%           | +10%          | Balanced defensive |
| Diamond       | +5%          | +5%           | Balanced |
| Battle Line   | +10%         | 0%            | Aggressive |
| Skirmish      | +15%         | -10%          | High risk/reward |
| Tee Forward   | +8%          | +2%           | Offensive focus |
| Enfilade      | +12%         | -5%           | Flanking bonus |
| Tee Reverse   | -5%          | +15%          | Defensive focus |

### 3.8 Targeting Commands

1. **max_attack:** Target enemy stack with highest attack power
2. **min_attack:** Target enemy stack with lowest attack power (pick off weak)
3. **max_durability:** Target enemy stack with highest HP (focus fire strongest)
4. **min_durability:** Target enemy stack with lowest HP (finish off wounded)
5. **closest:** Target enemy front line (shortest range)
6. **by_commander_rank:** Target enemy commander's stack first

### 3.9 Round Limits & Victory Conditions

**Min Rounds:** 1 (instant victory if one side has 0 ships)
**Max Rounds:** 99 (draw if neither side wins by round 99)
**Typical Combat:** 20-40 rounds for balanced fleets

**Victory Conditions:**
- Attacker Win: Defender has 0 ships remaining
- Defender Win: Attacker has 0 ships remaining OR attacker retreats
- Draw: Round 99 reached, both sides have ships

**He3 Consumption:**
- Each ship consumes `he3_per_round` per round
- Total: `he3_consumed = ships × he3_per_round × rounds`

---

## 4. COMBAT ENGINE DESIGN

### 4.1 Core Combat Loop

**High-Level Pseudocode:**
```go
func ResolveCombat(attacker Fleet, defender Fleet) CombatReport {
    // Initialize combat state
    attackerState := BuildFleetState(attacker)
    defenderState := BuildFleetState(defender)

    report := CombatReport{
        rounds: []
    }

    // Combat loop
    for round := 1; round <= 99; round++ {
        // Phase 1: Positioning (turn order by speed)
        turnOrder := DetermineTurnOrder(attackerState, defenderState)

        roundData := RoundData{round: round}

        // Each side attacks
        for _, fleet := range turnOrder {
            if fleet == attacker {
                // Phase 2-8 for attacker
                damage := ExecuteAttack(attackerState, defenderState)
                casualties := ApplyDamage(defenderState, damage)

                roundData.attacker_damage += damage.total
                roundData.defender_ships_lost += casualties
            } else {
                // Phase 2-8 for defender
                damage := ExecuteAttack(defenderState, attackerState)
                casualties := ApplyDamage(attackerState, damage)

                roundData.defender_damage += damage.total
                roundData.attacker_ships_lost += casualties
            }
        }

        roundData.attacker_remaining = attackerState.TotalShips()
        roundData.defender_remaining = defenderState.TotalShips()

        report.rounds = append(report.rounds, roundData)

        // Check victory conditions
        if defenderState.TotalShips() == 0 {
            report.result = "attacker_win"
            break
        }
        if attackerState.TotalShips() == 0 {
            report.result = "defender_win"
            break
        }
    }

    if report.result == "" {
        report.result = "draw" // Round 99 reached
    }

    report.total_rounds = len(report.rounds)

    return report
}
```

### 4.2 Fleet State Structure

```go
type FleetState struct {
    PlayerID       string
    FleetID        string
    CommanderBonus CommanderBonus
    TechBonus      TechBonus
    Formation      string // 'phalanx', 'diamond', etc.
    TargetCommand  string // 'max_attack', 'min_attack', etc.
    Stacks         []StackState // 3x3 grid = 9 stacks max
}

type StackState struct {
    Position       int    // 0-8 (3x3 grid position)
    ShipDesignID   string
    ShipCount      int    // Current ships in this stack
    MaxShipCount   int    // Original ship count

    // Ship stats (from design)
    ArmorType      string // 'nano', 'chrome', etc.
    Shield         int
    Structure      int    // HP per ship
    Defense        float64
    Agility        int
    AttackPower    int

    // Weapons
    Weapons        []Weapon

    // State per round
    CurrentShield  int
    CurrentHP      int // Total HP of all ships in stack
}

type Weapon struct {
    DamageType     string // 'heat', 'kinetic', etc.
    MinDamage      int
    MaxDamage      int
    HitRate        int    // Base hit rate (e.g., 85%)
    He3PerRound    int
    Quantity       int    // Number of this weapon on ship
}

type CommanderBonus struct {
    Accuracy       int
    Dodge          int
    Speed          int
    Electron       int
    EffectiveStack int   // 300 + star_rank * 50
}

type TechBonus struct {
    BallisticDamage      float64 // Percent bonus
    BallisticHitRate     float64
    DirectionalDamage    float64
    MissileDamage        float64
    MissileHitRate       float64
    FighterDamage        float64
    BaseShield           float64 // Percent bonus
    BaseStructure        float64
    BaseAgility          float64
    // ... (from tech_effects.go TechBonuses)
}
```

### 4.3 Phase-by-Phase Implementation

**Phase 1: Positioning (Turn Order)**
```go
func DetermineTurnOrder(attacker, defender *FleetState) []*FleetState {
    attackerSpeed := attacker.CommanderBonus.Speed + attacker.TechBonus.SpeedBonus
    defenderSpeed := defender.CommanderBonus.Speed + defender.TechBonus.SpeedBonus

    if attackerSpeed > defenderSpeed {
        return []*FleetState{attacker, defender}
    } else if defenderSpeed > attackerSpeed {
        return []*FleetState{defender, attacker}
    } else {
        // Tie: attacker goes first
        return []*FleetState{attacker, defender}
    }
}
```

**Phase 2: Targeting**
```go
func SelectTarget(attacker, defender *FleetState) *StackState {
    switch attacker.TargetCommand {
    case "max_attack":
        return findStackWithMaxAttack(defender.Stacks)
    case "min_attack":
        return findStackWithMinAttack(defender.Stacks)
    case "max_durability":
        return findStackWithMaxHP(defender.Stacks)
    case "min_durability":
        return findStackWithMinHP(defender.Stacks)
    case "closest":
        return findClosestStack(defender.Stacks)
    case "by_commander_rank":
        return findCommanderStack(defender.Stacks)
    default:
        return findStackWithMaxAttack(defender.Stacks)
    }
}
```

**Phase 3: Attack Calculation**
```go
func CalculateAttacks(attackerStack *StackState, attacker *FleetState) int {
    // Base attacks = ship count × weapons per ship
    baseAttacks := attackerStack.ShipCount

    // Effective stack multiplier (commander bonus)
    effectiveStackMultiplier := float64(attacker.CommanderBonus.EffectiveStack) / 300.0
    totalAttacks := int(float64(baseAttacks) * effectiveStackMultiplier)

    return totalAttacks
}
```

**Phase 4: Interception** (simplified for Module 5)
```go
// For Module 5, we'll skip interceptors (no PPC modules yet)
// Module 6+ can add this
func ApplyInterception(attacks int, defenderStack *StackState) int {
    // TODO: Calculate interceptor count from defender modules
    // interceptRate := 0.55 per PPC
    // blockedAttacks := int(float64(attacks) * interceptRate)
    return attacks // No interception in Module 5
}
```

**Phase 5: Hit Chance & Damage Roll**
```go
func CalculateHits(attacks int, attackerStack *StackState, defenderStack *StackState,
                    attacker, defender *FleetState) []int {
    damages := []int{}

    for i := 0; i < attacks; i++ {
        // Hit chance formula
        attackerAccuracy := 100.0 + float64(attacker.CommanderBonus.Accuracy)
        defenderAgility := float64(defenderStack.Agility) + float64(defender.CommanderBonus.Dodge)

        hitChance := attackerAccuracy / defenderAgility
        if hitChance > 1.0 {
            hitChance = 1.0 // Cap at 100%
        }

        // Roll hit
        if rand.Float64() < hitChance {
            // Hit! Roll damage
            weapon := selectRandomWeapon(attackerStack.Weapons)
            damage := randInt(weapon.MinDamage, weapon.MaxDamage)

            // Apply formation bonus
            formationBonus := getFormationAttackBonus(attacker.Formation)
            damage = int(float64(damage) * (1.0 + formationBonus))

            // Apply tech bonus
            techBonus := getTechDamageBonus(weapon.DamageType, attacker.TechBonus)
            damage = int(float64(damage) * (1.0 + techBonus/100.0))

            // Critical hit?
            critChance := 0.05 + (float64(attacker.CommanderBonus.Electron) / 200.0)
            if rand.Float64() < critChance {
                damage = int(float64(damage) * 1.5) // +50% crit damage
            }

            // Ship type advantage?
            typeBonus := getShipTypeAdvantage(attackerStack.HullType, defenderStack.HullType)
            damage = int(float64(damage) * typeBonus)

            damages = append(damages, damage)

            // Successive strike?
            successiveChance := float64(attacker.CommanderBonus.Speed) / 500.0
            if rand.Float64() < successiveChance {
                // Attack twice!
                damage2 := randInt(weapon.MinDamage, weapon.MaxDamage)
                damages = append(damages, damage2)
            }
        }
    }

    return damages
}
```

**Phase 6: Shield Absorption**
```go
func ApplyShieldDamage(damages []int, stack *StackState) []int {
    remainingDamages := []int{}

    for _, dmg := range damages {
        if stack.CurrentShield > 0 {
            absorbed := min(dmg, stack.CurrentShield)
            stack.CurrentShield -= absorbed
            overflow := dmg - absorbed
            if overflow > 0 {
                remainingDamages = append(remainingDamages, overflow)
            }
        } else {
            remainingDamages = append(remainingDamages, dmg)
        }
    }

    return remainingDamages
}
```

**Phase 7: Armor Effectiveness**
```go
func ApplyArmorReduction(damages []int, damageTypes []string, armorType string) []int {
    modifiedDamages := []int{}

    for i, dmg := range damages {
        multiplier := getArmorEffectiveness(armorType, damageTypes[i])
        finalDamage := int(float64(dmg) * multiplier)
        modifiedDamages = append(modifiedDamages, finalDamage)
    }

    return modifiedDamages
}

func getArmorEffectiveness(armor, damageType string) float64 {
    matrix := map[string]map[string]float64{
        "chrome": {
            "heat":       0.5,
            "kinetic":    0.5,
            "magnetic":   1.5,
            "explosive":  2.0,
        },
        "nano": {
            "heat":       0.5,
            "kinetic":    1.5,
            "magnetic":   2.0,
            "explosive":  0.5,
        },
        "neutralizing": {
            "heat":       1.5,
            "kinetic":    2.0,
            "magnetic":   0.5,
            "explosive":  0.5,
        },
        "regen": {
            "heat":       2.0,
            "kinetic":    0.5,
            "magnetic":   0.5,
            "explosive":  1.5,
        },
        "light": {
            "heat":       1.0,
            "kinetic":    1.0,
            "magnetic":   1.0,
            "explosive":  1.0,
        },
    }

    if armorMap, ok := matrix[armor]; ok {
        if multiplier, ok := armorMap[damageType]; ok {
            return multiplier
        }
    }

    return 1.0 // Default neutral
}
```

**Phase 8: HP Reduction & Casualties**
```go
func ApplyCasualties(stack *StackState, damages []int) int {
    totalDamage := sum(damages)

    // Apply defense reduction
    effectiveDamage := int(float64(totalDamage) / (1.0 + stack.Defense/100.0))

    stack.CurrentHP -= effectiveDamage

    if stack.CurrentHP <= 0 {
        // Stack destroyed
        casualties := stack.ShipCount
        stack.ShipCount = 0
        stack.CurrentHP = 0
        return casualties
    }

    // Calculate ships lost
    hpPerShip := stack.Structure
    shipsLost := (stack.MaxShipCount * hpPerShip - stack.CurrentHP) / hpPerShip

    if shipsLost > stack.ShipCount {
        shipsLost = stack.ShipCount
    }

    stack.ShipCount -= shipsLost

    return shipsLost
}
```

---

## 5. BACKEND IMPLEMENTATION

### 5.1 Service Layer (`backend/internal/services/combat_engine.go`)

**File Structure:**
```go
package services

// Core types
type FleetState struct { ... }
type StackState struct { ... }
type CombatResult struct { ... }
type RoundData struct { ... }

// Main entry point
func ResolveCombat(attackerFleetID, defenderFleetID string, combatType string) (*CombatResult, error)

// Fleet building
func BuildFleetState(fleetID string) (*FleetState, error)
func LoadCommander(commanderID string) (*CommanderBonus, error)
func LoadTechBonuses(playerID string) (*TechBonus, error)

// Combat phases
func DetermineTurnOrder(attacker, defender *FleetState) []*FleetState
func SelectTarget(attacker, defender *FleetState) *StackState
func CalculateAttacks(stack *StackState, fleet *FleetState) int
func CalculateHits(attacks int, attackerStack, defenderStack *StackState, attacker, defender *FleetState) []DamageRoll
func ApplyShieldDamage(damages []DamageRoll, stack *StackState) []DamageRoll
func ApplyArmorReduction(damages []DamageRoll, stack *StackState) []DamageRoll
func ApplyCasualties(stack *StackState, damages []DamageRoll) int

// Helpers
func getArmorEffectiveness(armor, damageType string) float64
func getFormationAttackBonus(formation string) float64
func getFormationDefenseBonus(formation string) float64
func getShipTypeAdvantage(attackerType, defenderType string) float64
func getTechDamageBonus(damageType string, techBonus *TechBonus) float64

// Loot
func CalculateInstanceLoot(instanceID string, victory bool) (map[string]int64, error)

// Persistence
func SaveCombatReport(result *CombatResult) error
func UpdateFleetShips(fleetID string, stacks []StackState) error
```

**Time Estimate:** 18-24 hours

---

### 5.2 Handler Layer (`backend/internal/handlers/combat.go`)

**Endpoints:**

1. **POST /api/instances/:id/attack** - Start instance combat (replace placeholder)
2. **GET /api/combat-reports/:id** - View detailed combat report
3. **GET /api/combat-reports** - List player's combat reports
4. **POST /api/combat/simulate** - Debug endpoint for testing combat (optional)

**Handler Code Skeleton:**

```go
package handlers

// AttackInstance initiates combat with an instance
func AttackInstance(w http.ResponseWriter, r *http.Request) {
    playerID := getPlayerIDFromContext(r)
    vars := mux.Vars(r)
    instanceID := vars["id"]

    var req struct {
        FleetID string `json:"fleet_id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
        return
    }

    // Validate fleet belongs to player
    var fleetOwner string
    err := database.DB.QueryRow(`SELECT player_id FROM fleets WHERE id = $1`, req.FleetID).Scan(&fleetOwner)
    if err != nil || fleetOwner != playerID {
        http.Error(w, `{"error":"Fleet not found"}`, http.StatusBadRequest)
        return
    }

    // Check fleet has ships
    var shipCount int
    database.DB.QueryRow(`SELECT COUNT(*) FROM fleet_ships WHERE fleet_id = $1`, req.FleetID).Scan(&shipCount)
    if shipCount == 0 {
        http.Error(w, `{"error":"Fleet has no ships"}`, http.StatusBadRequest)
        return
    }

    // Get instance defender fleet
    var defenderFleetID string
    err = database.DB.QueryRow(`SELECT defender_fleet_json->>'fleet_id' FROM instances WHERE id = $1`, instanceID).Scan(&defenderFleetID)
    if err != nil {
        http.Error(w, `{"error":"Instance not found"}`, http.StatusNotFound)
        return
    }

    // Resolve combat
    result, err := services.ResolveCombat(req.FleetID, defenderFleetID, "instance_normal")
    if err != nil {
        log.Printf("Combat resolution failed: %v", err)
        http.Error(w, `{"error":"Combat failed"}`, http.StatusInternalServerError)
        return
    }

    // Save combat report
    err = services.SaveCombatReport(result)
    if err != nil {
        log.Printf("Failed to save combat report: %v", err)
    }

    // Update fleet ships (apply casualties)
    err = services.UpdateFleetShips(req.FleetID, result.AttackerFinalState.Stacks)
    if err != nil {
        log.Printf("Failed to update fleet ships: %v", err)
    }

    // Grant loot if victory
    if result.Result == "attacker_win" {
        loot, err := services.CalculateInstanceLoot(instanceID, true)
        if err == nil {
            // Add loot to player resources
            database.DB.Exec(`
                UPDATE resources
                SET metal = metal + $1, he3 = he3 + $2, gold = gold + $3
                WHERE player_id = $4
            `, loot["metal"], loot["he3"], loot["gold"], playerID)
        }
    }

    // Return combat report
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

// GetCombatReport retrieves a detailed combat report
func GetCombatReport(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    reportID := vars["id"]

    var report models.CombatReport
    var roundsJSON []byte
    var lootJSON []byte
    var attackerLossesJSON []byte
    var defenderLossesJSON []byte

    err := database.DB.QueryRow(`
        SELECT id, attacker_id, defender_id, attacker_fleet_id, defender_fleet_id,
               combat_type, result, total_rounds, rounds_json, loot_json,
               attacker_losses_json, defender_losses_json, he3_consumed, created_at
        FROM combat_reports
        WHERE id = $1
    `, reportID).Scan(
        &report.ID, &report.AttackerID, &report.DefenderID,
        &report.AttackerFleetID, &report.DefenderFleetID,
        &report.CombatType, &report.Result, &report.TotalRounds,
        &roundsJSON, &lootJSON, &attackerLossesJSON, &defenderLossesJSON,
        &report.He3Consumed, &report.CreatedAt,
    )

    if err != nil {
        http.Error(w, `{"error":"Combat report not found"}`, http.StatusNotFound)
        return
    }

    json.Unmarshal(roundsJSON, &report.Rounds)
    json.Unmarshal(lootJSON, &report.Loot)
    json.Unmarshal(attackerLossesJSON, &report.AttackerLosses)
    json.Unmarshal(defenderLossesJSON, &report.DefenderLosses)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(report)
}

// ListCombatReports lists all combat reports for a player
func ListCombatReports(w http.ResponseWriter, r *http.Request) {
    playerID := getPlayerIDFromContext(r)

    rows, err := database.DB.Query(`
        SELECT id, attacker_id, defender_id, combat_type, result, total_rounds, created_at
        FROM combat_reports
        WHERE attacker_id = $1 OR defender_id = $1
        ORDER BY created_at DESC
        LIMIT 50
    `, playerID)

    if err != nil {
        http.Error(w, `{"error":"Failed to fetch reports"}`, http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    reports := []map[string]interface{}{}
    for rows.Next() {
        var report map[string]interface{}
        var id, attackerID, defenderID, combatType, result string
        var totalRounds int
        var createdAt time.Time

        rows.Scan(&id, &attackerID, &defenderID, &combatType, &result, &totalRounds, &createdAt)

        report = map[string]interface{}{
            "id":           id,
            "attacker_id":  attackerID,
            "defender_id":  defenderID,
            "combat_type":  combatType,
            "result":       result,
            "total_rounds": totalRounds,
            "created_at":   createdAt,
        }

        reports = append(reports, report)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reports)
}
```

**Routes:**
```go
// Combat
api.HandleFunc("/instances/{id}/attack", handlers.AttackInstance).Methods("POST")
api.HandleFunc("/combat-reports/{id}", handlers.GetCombatReport).Methods("GET")
api.HandleFunc("/combat-reports", handlers.ListCombatReports).Methods("GET")
```

**Time Estimate:** 6-8 hours

---

## 6. FRONTEND IMPLEMENTATION

### 6.1 API Client (`frontend/src/api/combat.ts`)

```typescript
import { apiClient } from './client';

export interface CombatReport {
  id: string;
  attacker_id: string;
  defender_id: string;
  attacker_fleet_id: string;
  defender_fleet_id: string;
  combat_type: string;
  result: 'attacker_win' | 'defender_win' | 'draw';
  total_rounds: number;
  rounds: RoundData[];
  loot: { [key: string]: number };
  attacker_losses: { [key: string]: number };
  defender_losses: { [key: string]: number };
  he3_consumed: number;
  created_at: string;
}

export interface RoundData {
  round: number;
  attacker_damage: number;
  defender_damage: number;
  attacker_ships_lost: number;
  defender_ships_lost: number;
  attacker_remaining: number;
  defender_remaining: number;
  events: string[];
}

export const combatApi = {
  attackInstance: (instanceId: string, fleetId: string) =>
    apiClient.post<CombatReport>(`/instances/${instanceId}/attack`, { fleet_id: fleetId }),

  getCombatReport: (reportId: string) =>
    apiClient.get<CombatReport>(`/combat-reports/${reportId}`),

  listCombatReports: () =>
    apiClient.get<CombatReport[]>('/combat-reports'),
};
```

**Time Estimate:** 1 hour

---

### 6.2 React Hook (`frontend/src/hooks/useCombat.ts`)

```typescript
import { useState } from 'react';
import { combatApi, CombatReport } from '../api/combat';

export function useCombat() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentReport, setCurrentReport] = useState<CombatReport | null>(null);

  const attackInstance = async (instanceId: string, fleetId: string): Promise<CombatReport | null> => {
    try {
      setLoading(true);
      const report = await combatApi.attackInstance(instanceId, fleetId);
      setCurrentReport(report);
      setError(null);
      return report;
    } catch (err: any) {
      setError(err.message || 'Attack failed');
      return null;
    } finally {
      setLoading(false);
    }
  };

  const viewReport = async (reportId: string): Promise<CombatReport | null> => {
    try {
      setLoading(true);
      const report = await combatApi.getCombatReport(reportId);
      setCurrentReport(report);
      setError(null);
      return report;
    } catch (err: any) {
      setError(err.message || 'Failed to load report');
      return null;
    } finally {
      setLoading(false);
    }
  };

  return {
    loading,
    error,
    currentReport,
    attackInstance,
    viewReport,
  };
}
```

**Time Estimate:** 1 hour

---

### 6.3 Combat Report Panel (`frontend/src/components/panels/CombatReportPanel.tsx`)

**Features:**
- Round-by-round summary table
- Damage charts (attacker vs defender per round)
- Casualties timeline
- Final results (victory, loot, losses)
- Event log

```typescript
import React from 'react';
import { CombatReport, RoundData } from '../../api/combat';

interface Props {
  report: CombatReport;
}

export function CombatReportPanel({ report }: Props) {
  const renderResult = () => {
    switch (report.result) {
      case 'attacker_win':
        return <div className="result victory">Victory!</div>;
      case 'defender_win':
        return <div className="result defeat">Defeat</div>;
      case 'draw':
        return <div className="result draw">Draw</div>;
    }
  };

  const renderRoundsSummary = () => {
    return (
      <table className="rounds-table">
        <thead>
          <tr>
            <th>Round</th>
            <th>Attacker Damage</th>
            <th>Defender Damage</th>
            <th>Attacker Lost</th>
            <th>Defender Lost</th>
            <th>Attacker Ships</th>
            <th>Defender Ships</th>
          </tr>
        </thead>
        <tbody>
          {report.rounds.map((round) => (
            <tr key={round.round}>
              <td>{round.round}</td>
              <td className="damage">{round.attacker_damage.toLocaleString()}</td>
              <td className="damage">{round.defender_damage.toLocaleString()}</td>
              <td className="losses">{round.attacker_ships_lost}</td>
              <td className="losses">{round.defender_ships_lost}</td>
              <td className="ships">{round.attacker_remaining}</td>
              <td className="ships">{round.defender_remaining}</td>
            </tr>
          ))}
        </tbody>
      </table>
    );
  };

  const renderLoot = () => {
    if (report.result !== 'attacker_win' || !report.loot) return null;

    return (
      <div className="loot-section">
        <h3>Loot</h3>
        <div className="resources">
          {Object.entries(report.loot).map(([resource, amount]) => (
            <div key={resource} className="resource-item">
              <span className="label">{resource}:</span>
              <span className="value">+{amount.toLocaleString()}</span>
            </div>
          ))}
        </div>
      </div>
    );
  };

  return (
    <div className="combat-report-panel">
      <h2>Combat Report</h2>

      {renderResult()}

      <div className="summary">
        <div className="stat">
          <span className="label">Total Rounds:</span>
          <span className="value">{report.total_rounds}</span>
        </div>
        <div className="stat">
          <span className="label">He3 Consumed:</span>
          <span className="value">{report.he3_consumed.toLocaleString()}</span>
        </div>
      </div>

      {renderRoundsSummary()}

      {renderLoot()}

      <div className="losses-section">
        <h3>Casualties</h3>
        <div className="losses-grid">
          <div className="attacker-losses">
            <h4>Attacker Losses</h4>
            {Object.entries(report.attacker_losses).map(([shipType, count]) => (
              <div key={shipType}>
                {shipType}: {count}
              </div>
            ))}
          </div>
          <div className="defender-losses">
            <h4>Defender Losses</h4>
            {Object.entries(report.defender_losses).map(([shipType, count]) => (
              <div key={shipType}>
                {shipType}: {count}
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
```

**Time Estimate:** 4-5 hours

---

### 6.4 Instance Panel - Attack Button

Update `InstancePanel.tsx` to add attack functionality:

```typescript
import { useCombat } from '../../hooks/useCombat';
import { useFleets } from '../../hooks/useFleets';

// Inside InstancePanel component
const { attackInstance, currentReport } = useCombat();
const { fleets } = useFleets();

const handleAttack = async (instanceId: string) => {
  // Show fleet selection modal
  const selectedFleetId = await showFleetSelectionModal(fleets);

  if (selectedFleetId) {
    const report = await attackInstance(instanceId, selectedFleetId);

    if (report) {
      // Show combat report
      setShowCombatReport(true);
    }
  }
};

// In JSX for each instance:
<button onClick={() => handleAttack(instance.id)}>
  Attack
</button>
```

**Time Estimate:** 2-3 hours

---

## 7. INTEGRATION & TESTING

### 7.1 Quest Integration

Update `quest_service.go` to track combat victories:

```go
case "defeat_instance":
    // Triggered when instance combat won
    // targetValue = instance name (e.g., "Ancestral Recall")
    progress = 1

case "win_pvp_battle":
    // Triggered when PvP combat won (Module 6)
    progress = 1
```

Add quest trigger in `SaveCombatReport`:
```go
if result.Result == "attacker_win" && result.CombatType == "instance_normal" {
    UpdateQuestProgress(result.AttackerID, "defeat_instance", instanceName, 1)
}
```

**Time Estimate:** 1 hour

---

### 7.2 Tech Bonuses Integration

Combat engine must read tech bonuses from `services.GetPlayerTechBonuses()`:

```go
func LoadTechBonuses(playerID string) (*TechBonus, error) {
    bonuses, err := GetPlayerTechBonuses(playerID)
    if err != nil {
        return nil, err
    }

    return &TechBonus{
        BallisticDamage:      bonuses.BallisticDamage,
        BallisticHitRate:     bonuses.BallisticHitRate,
        DirectionalDamage:    bonuses.DirectionalDamage,
        MissileDamage:        bonuses.MissileDamage,
        MissileHitRate:       bonuses.MissileHitRate,
        FighterDamage:        bonuses.FighterDamage,
        BaseShield:           bonuses.BaseShield,
        BaseStructure:        bonuses.BaseStructure,
        BaseAgility:          bonuses.BaseAgility,
        BaseDefense:          bonuses.BaseDefense,
        // ... all tech bonuses
    }, nil
}
```

Apply in damage calculation:
```go
damage := int(float64(baseDamage) * (1.0 + techBonus.BallisticDamage/100.0))
```

**Time Estimate:** 2 hours

---

### 7.3 Commander Bonuses Integration

Load commander from fleet:

```go
func LoadCommander(commanderID string) (*CommanderBonus, error) {
    if commanderID == "" {
        // No commander assigned
        return &CommanderBonus{
            Accuracy:       0,
            Dodge:          0,
            Speed:          0,
            Electron:       0,
            EffectiveStack: 300, // Base effective stack
        }, nil
    }

    var cmd CommanderBonus
    err := database.DB.QueryRow(`
        SELECT accuracy, dodge, speed, electron, effective_stack
        FROM commanders
        WHERE id = $1
    `, commanderID).Scan(&cmd.Accuracy, &cmd.Dodge, &cmd.Speed, &cmd.Electron, &cmd.EffectiveStack)

    if err != nil {
        return nil, err
    }

    return &cmd, nil
}
```

Apply in hit chance:
```go
attackerAccuracy := 100.0 + float64(attacker.CommanderBonus.Accuracy)
defenderAgility := float64(defenderStack.Agility) + float64(defender.CommanderBonus.Dodge)
hitChance := attackerAccuracy / defenderAgility
```

**Time Estimate:** 2 hours

---

## 8. QA CHECKLIST

### 8.1 Combat Engine Tests
- [ ] ResolveCombat() completes without errors
- [ ] Turn order determined by speed (higher speed goes first)
- [ ] Targeting commands work (max_attack, min_attack, etc.)
- [ ] Hit chance formula correct (accuracy / agility)
- [ ] Damage rolls within weapon range (min_damage to max_damage)
- [ ] Critical hits trigger (~5% + electron/200)
- [ ] Successive strikes trigger (speed/500)
- [ ] Formation bonuses apply (attack/defense modifiers)
- [ ] Ship type advantage applies (+20% damage)
- [ ] Armor effectiveness matrix correct (chrome vs heat = 0.5×)
- [ ] Shield absorbs damage before hull
- [ ] HP reduction calculates casualties correctly
- [ ] Combat ends when one side has 0 ships
- [ ] Draw declared at round 99 if both sides alive
- [ ] He3 consumption calculated correctly

### 8.2 Instance Combat Tests
- [ ] POST /api/instances/:id/attack works
- [ ] Fleet validation (fleet belongs to player)
- [ ] Instance defender fleet loaded correctly
- [ ] Combat resolves with expected winner
- [ ] Combat report saved to database
- [ ] Fleet ships updated after combat (casualties)
- [ ] Loot granted on victory
- [ ] Quest progress triggered (defeat_instance)
- [ ] 30 instances all attackable

### 8.3 Integration Tests
- [ ] Commander bonuses apply in combat
- [ ] Tech bonuses apply in combat
- [ ] No commander = base 300 effective stack
- [ ] Commander star rank 5 = 550 effective stack (~2x ships per stack)
- [ ] Ballistic Damage Tech +10% = +10% ballistic weapon damage
- [ ] Ship Defense Tech +5% structure = +5% HP
- [ ] Formation Phalanx = +10% defense
- [ ] Formation Battle Line = +10% attack

### 8.4 Frontend Tests
- [ ] Instance panel shows "Attack" button
- [ ] Fleet selection modal works
- [ ] Combat report displays after attack
- [ ] Round summary table shows all rounds
- [ ] Victory/defeat/draw displayed correctly
- [ ] Loot shown on victory
- [ ] Casualties displayed
- [ ] Combat reports list shows history
- [ ] Report detail view works

### 8.5 Edge Cases
- [ ] Attack with empty fleet = error
- [ ] Attack with insufficient He3 = error (future)
- [ ] Attack with fleet already in combat = error (future)
- [ ] 99 round draw handled correctly
- [ ] 1 round instant win handled correctly
- [ ] All ships destroyed in round 50 = victory
- [ ] Shield regeneration per round (future enhancement)
- [ ] Armor type NULL = light armor (1.0× all)
- [ ] Damage type NULL = neutral damage (1.0× all)

**Total QA Checks:** 45

---

## 9. RISK ASSESSMENT

### 9.1 Technical Risks

**RISK 1: Combat Balance (Too Easy/Hard)**
- **Severity:** HIGH
- **Impact:** Instances too easy = boring, too hard = frustrating
- **Mitigation:** Conservative instance defender stats, extensive testing
- **Resolution:** Balancing pass after Module 5 complete

**RISK 2: Performance (Long Combat Rounds)**
- **Severity:** MEDIUM
- **Impact:** 99-round combat with 9 stacks each = thousands of calculations
- **Mitigation:** Optimize critical path (hit chance, damage roll), profile performance
- **Target:** <500ms for average 30-round combat

**RISK 3: Armor Effectiveness Matrix Accuracy**
- **Severity:** MEDIUM
- **Impact:** Incorrect multipliers break game balance
- **Mitigation:** Values derived from GO2 wiki, verify with testing
- **Fallback:** Adjust multipliers if too extreme

**RISK 4: Commander Effective Stack Too Powerful**
- **Severity:** MEDIUM
- **Impact:** Star rank 9 commanders (750 effective stack) = 2.5× ships per stack
- **Mitigation:** This matches GO2, but we can adjust formula if needed
- **Alternative:** Cap effective stack at 600 (2× multiplier)

**RISK 5: Random Number Generator Bias**
- **Severity:** LOW
- **Impact:** Poor RNG = unfair combat outcomes
- **Mitigation:** Use crypto/rand for better distribution, seed properly
- **Test:** 10,000 combat simulations should show expected win rates

**RISK 6: Float Precision in Damage Calculations**
- **Severity:** LOW
- **Impact:** Rounding errors accumulate over 99 rounds
- **Mitigation:** Use integers where possible, round at final step
- **Test:** Verify total damage adds up correctly

---

### 9.2 Integration Risks

**RISK 7: Tech Bonuses Not Applied**
- **Severity:** HIGH
- **Impact:** Techs feel useless if not affecting combat
- **Mitigation:** Explicit integration test for each tech bonus type
- **Test:** Research Ballistics 1 (+5% damage) → verify +5% in combat

**RISK 8: Commander Bonuses Not Applied**
- **Severity:** HIGH
- **Impact:** Commanders feel useless if not affecting combat
- **Mitigation:** Explicit test for each commander stat
- **Test:** Commander accuracy 30 → verify hit rate increase

**RISK 9: Fleet Ships Not Updated After Combat**
- **Severity:** CRITICAL
- **Impact:** Casualties not applied = infinite ships
- **Mitigation:** Transaction-based fleet update, rollback on error
- **Test:** Attack instance → lose 500 ships → verify fleet count reduced

---

### 9.3 UX Risks

**RISK 10: Combat Report Too Complex**
- **Severity:** MEDIUM
- **Impact:** Users overwhelmed by 99-round detailed log
- **Mitigation:** Collapsible rounds, summary view by default
- **Alternative:** Show only key rounds (first 3, last 3, critical hits)

**RISK 11: Combat Takes Too Long**
- **Severity:** LOW
- **Impact:** 500ms × 30 rounds = 15 seconds feels slow
- **Mitigation:** Show "Combat in progress..." spinner, async processing
- **Alternative:** Pre-calculate combat, show instant result

**RISK 12: No Combat Animation**
- **Severity:** LOW
- **Impact:** Text-only combat feels dry
- **Mitigation:** This is Phase A foundation, animation in Module 11 (Polish)
- **Acceptable:** Text reports are standard in GO2 clones

---

### 9.4 Timeline Risks

**RISK 13: Underestimated Complexity**
- **Severity:** HIGH
- **Impact:** 8-phase combat engine is the most complex code in the project
- **Buffer:** Estimated 40-55 hours with +25% contingency
- **Contingency:** If >60 hours, simplify (remove interception, simplify armor)

**RISK 14: Testing Takes Longer Than Expected**
- **Severity:** MEDIUM
- **Impact:** 45 QA checks + balancing = significant testing time
- **Mitigation:** Allocate 15-20 hours for QA/balancing phase
- **Alternative:** Deploy with basic balance, iterate in production

---

## 10. IMPLEMENTATION TIMELINE

### 10.1 Phase Breakdown

**Phase 1: Combat Engine Core (Days 1-3)**
- FleetState/StackState structures
- BuildFleetState() (load ships, commander, techs)
- ResolveCombat() main loop
- Phases 1-2: Positioning, Targeting
- **Time:** 12-15 hours

**Phase 2: Damage Calculation (Days 3-4)**
- Phase 3-5: Attack calculation, hit chance, damage rolls
- Critical hits, successive strikes
- Tech bonuses integration
- Commander bonuses integration
- **Time:** 8-10 hours

**Phase 3: Defense & Casualties (Day 5)**
- Phase 6-8: Shields, armor, HP reduction
- Armor effectiveness matrix
- Formation bonuses
- Ship type advantage
- **Time:** 8-10 hours

**Phase 4: Persistence & Loot (Day 5)**
- SaveCombatReport()
- UpdateFleetShips()
- CalculateInstanceLoot()
- Quest integration
- **Time:** 4-6 hours

**Phase 5: Backend Handlers (Day 6)**
- AttackInstance endpoint
- GetCombatReport endpoint
- ListCombatReports endpoint
- Error handling
- **Time:** 6-8 hours

**Phase 6: Frontend UI (Days 7-8)**
- API client + hook
- CombatReportPanel component
- Instance attack button
- Fleet selection modal
- **Time:** 8-10 hours

**Phase 7: QA & Balancing (Day 9)**
- Run full QA checklist (45 checks)
- Combat simulation tests (10,000 battles)
- Balance instance defender stats
- Fix bugs
- **Time:** 12-16 hours

---

### 10.2 Critical Path

```
Day 1-2: FleetState → BuildFleetState → ResolveCombat loop
         ↓
Day 3-4: Damage calculation → Tech/Commander integration
         ↓
Day 5: Defense mechanics → Casualties → Loot
         ↓
Day 6: Backend handlers → Instance attack endpoint
         ↓
Day 7-8: Frontend UI → CombatReportPanel
         ↓
Day 9: QA testing → Balancing → Bug fixes
```

**Total Estimated Time:** 40-55 hours (6-9 working days)

---

## 11. SUCCESS CRITERIA

Module 5 is considered **COMPLETE** when:

✅ **Combat Engine:**
- 8-phase combat resolution implemented
- Hit chance, damage, armor, casualties all functional
- Min 1 round, max 99 rounds enforced
- Victory/defeat/draw conditions correct

✅ **Integration:**
- Commander bonuses apply correctly
- Tech bonuses apply correctly
- Formation bonuses apply correctly
- Ship type advantage works

✅ **Instance Combat:**
- All 30 Normal Instances attackable
- Combat reports generated and saved
- Fleet casualties applied
- Loot granted on victory
- Quest progress triggered

✅ **Frontend:**
- Combat report displays round-by-round data
- Victory/defeat/draw shown clearly
- Loot and casualties displayed
- Combat reports list functional

✅ **QA:**
- All 45 QA checklist items pass
- 10,000 simulation battles complete without errors
- Combat balance feels fair (instances winnable but challenging)
- No critical bugs

---

## 12. NEXT STEPS AFTER MODULE 5

After Module 5 completion, the following features become possible:

1. **Module 6: PvP Combat** - Reuse combat engine for player vs player
2. **Module 7: Space Defense Buildings** - Add defense buildings to defender fleet
3. **Balancing Pass** - Adjust instance difficulty, tech bonuses, commander impact
4. **Combat Animations** - Add visual effects (Module 11: Polish)

**Recommended Next Module:** Module 6 (PvP Combat) to enable player interaction.

---

## APPENDIX A: FORMULAS REFERENCE

### A.1 Hit Chance
```
Hit Chance = (Base Hit Rate + Commander Accuracy) / (Defender Agility + Commander Dodge)
Capped at 100% for display, but can exceed internally
```

### A.2 Critical Hits
```
Crit Chance = 5% + (Commander Electron / 200)
Crit Damage = Base Damage × 1.5
```

### A.3 Successive Strikes
```
Successive Strike Chance = Commander Speed / 500
Max: 20% at 100 speed
```

### A.4 Effective Stack
```
Effective Stack = Base Stack × (Commander Effective Stack / 300)
Base Stack: Frigate 1100, Cruiser 1000, Battleship 900
Commander Effective Stack: 300 + (Star Rank × 50)
```

### A.5 Formation Bonuses
```
Attack Multiplier = 1.0 + Formation Attack Bonus
Defense Multiplier = 1.0 + Formation Defense Bonus
```

### A.6 Ship Type Advantage
```
Frigate vs Cruiser: 1.2× damage
Cruiser vs Battleship: 1.2× damage
Battleship vs Frigate: 1.2× damage
Otherwise: 1.0× damage
```

### A.7 Armor Effectiveness
```
Final Damage = Base Damage × Armor Multiplier
Armor Multiplier: See 3.2 Armor Effectiveness Matrix
```

### A.8 Defense Reduction
```
Effective Damage = Base Damage / (1 + Defense / 100)
Defense from hull + modules + tech bonuses
```

### A.9 He3 Consumption
```
Total He3 = Ships × He3 Per Round × Rounds
```

---

## APPENDIX B: References

**Galaxy Online 2 Wiki Sources:**
- [Combat Mechanics](https://galaxyonlineii.fandom.com/wiki/Combat_Mechanics) - Core combat phases
- [All Ship Armor Types](https://galaxyonlineii.fandom.com/wiki/All_Ship_Armor_Types) - Armor effectiveness
- [Commander Cards](https://galaxyonlineii.fandom.com/wiki/Commander_Cards) - Commander bonuses

**Project Documents:**
- `/docs/planning/final-scope.md` - Module 5 scope definition
- `/docs/planning/module-4-commander-system-plan.md` - Commander integration
- `/backend/internal/services/ship_formulas.go` - Existing ship stats
- `/backend/internal/services/tech_effects.go` - Tech bonuses

**External Research:**
- [DevilsMMO GO2 Guide](https://www.devilsmmo.com/forum/galaxy-online-ii-attributes-analysis) - Attribute analysis
- [GuideScroll Ships List](https://guidescroll.com/2011/05/facebook-galaxy-online-ii-ships-and-armor-list/) - Ship armor types

---

**END OF MODULE 5 IMPLEMENTATION PLAN**
