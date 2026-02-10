# Module 7: Space Station Defense Buildings - Implementation Plan

**Date:** 2026-02-07
**Architect:** architect
**Module:** Space Station Defense Buildings
**Estimated Total:** 15-22 hours (3-4 days)

---

## 1. EXECUTIVE SUMMARY

### 1.1 Module Overview

The Space Station Defense Buildings module adds **5 orbital defense structures** that automatically attack incoming PvP fleets, significantly strengthening planetary defense. These buildings integrate with the **Planetary Defense Science tech tree** (12 techs from Module 1) and the **Module 6 PvP combat system**.

**Complexity Level:** MEDIUM (★★★☆☆)
- Simpler than combat/commanders (reuses existing combat engine)
- Database schema already 90% complete
- Main work: Combat integration + 3D models
- Lower risk: Defense buildings are optional enhancement to PvP

### 1.2 Scope (From final-scope.md)

**IN SCOPE:**
- ✅ 5 defense buildings (Meteor Star, Particle Cannon, Anti-Aircraft Gun, Thor's Cannon, Celestial Base)
- ✅ Auto-attack incoming fleets during PvP combat
- ✅ Integration with Planetary Defense Science techs (12 techs, already seeded)
- ✅ Defense stats contribute to planet defense value
- ✅ Building placement on 20×20 planet grid
- ✅ 3D procedural models for all 5 buildings
- ✅ Tech bonuses apply (attack, build speed, cost reduction, max count)

**OUT OF SCOPE:**
- ❌ Defense building repairs (buildings are destroyed, must rebuild)
- ❌ Defense building movement/repositioning (fixed after placement)
- ❌ Automated rebuild after destruction (manual only)
- ❌ Defense building blueprints (buildings unlocked by tech or default)

### 1.3 Current Implementation Status

**Database:** 90% Complete
- ✅ 5 building types seeded in `building_types`
- ✅ 4 `*_levels` lookup tables (meteor_star, particle_cannon, anti_aircraft_gun, thors_cannon)
- ✅ 12 Planetary Defense Science techs seeded
- ✅ Tech bonuses defined (defense_cost_reduction, defense_build_speed, emplacement_attack, etc.)
- ❌ NO `celestial_base_levels` lookup table (building exists but no level data)
- ❌ NO defense stats in level tables (attack, HP, range, cooldown)

**Backend:** 10% Complete
- ✅ Tech effects service already handles defense bonuses (tech_effects.go)
- ❌ NO defense building combat integration
- ❌ NO defense attack logic in PvP combat
- ❌ NO defense building placement validation

**Frontend:** 5% Complete
- ✅ Building placement system exists (from Phase 1)
- ❌ NO 3D models for 5 defense buildings
- ❌ NO defense-specific UI (attack stats display, placement limits)

---

## 2. DATABASE ANALYSIS

### 2.1 Existing Schema Review

**Table: `building_types` (existing - Phase 1)**
```sql
-- Defense buildings already seeded:
('meteor_star',        'Meteor Star',        'defense',  'space', ...)
('particle_cannon',    'Particle Cannon',    'defense',  'space', ...)
('anti_aircraft_gun',  'Anti-Aircraft Gun',  'defense',  'space', ...)
('thors_cannon',       'Thor''s Cannon',     'defense',  'space', ...)
('celestial_base',     'Celestial Base',     'space',    'space', ...)  -- Category 'space' not 'defense'
```

**Issue:** Celestial Base has category 'space' instead of 'defense'. Need to verify intent.

**Table: `meteor_star_levels` (existing - Phase 1)**
```sql
CREATE TABLE meteor_star_levels (
    level INTEGER PRIMARY KEY,
    space_station_req INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL
);
-- Seeded with 12 levels
```

**Similar tables:** `particle_cannon_levels`, `anti_aircraft_gun_levels`, `thors_cannon_levels`

**Issue:** Level tables have no combat stats (attack, HP, range, cooldown).

**Table: `tech_types` (existing - Module 1)**
```sql
-- Planetary Defense Science tree (12 techs):
1. energy_control               - defense_cost_reduction (-1-10%)
2. rapid_defense_buildup        - defense_build_speed (+1-10%)
3. defense_enhancement          - defense_value (+1-10%)
4. emplacement_mastery          - emplacement_attack (+1-10%)
5. utmost_defense_buildup       - max_defense_structures (+1-10%)
6. range_extension              - defense_range (Particle Cannon + AA Gun)
7. thor_buildup                 - max_thor_cannon (+1 flat)
8. augment_propulsion           - defense_movement (+1-2 ship speed)
9-12. (Additional techs not listed in grep output)
```

**Tech Bonuses Already Implemented:**
- `TechBonuses.DefenseCostReduction` ✅
- `TechBonuses.DefenseBuildSpeed` ✅
- `TechBonuses.DefenseValue` ✅
- `TechBonuses.EmplacementAttack` ✅
- `TechBonuses.MaxDefenseStructures` ✅
- `TechBonuses.MaxThorCannon` ✅

### 2.2 Required Schema Changes

**ALTER `*_levels` tables** - Add combat stats:

```sql
-- Meteor Star
ALTER TABLE meteor_star_levels
ADD COLUMN attack INTEGER NOT NULL DEFAULT 50,
ADD COLUMN hp BIGINT NOT NULL DEFAULT 200,
ADD COLUMN defense NUMERIC(5,2) NOT NULL DEFAULT 0,
ADD COLUMN range_squares INTEGER NOT NULL DEFAULT 0,  -- Meteor Star is melee (range 0)
ADD COLUMN cooldown_rounds INTEGER NOT NULL DEFAULT 0; -- No cooldown

-- Update seed data with GO2 values
UPDATE meteor_star_levels SET
    attack = 50,
    hp = CASE level
        WHEN 1 THEN 200
        WHEN 2 THEN 1000
        WHEN 3 THEN 5000
        WHEN 4 THEN 25000
        WHEN 5 THEN 125000
        WHEN 6 THEN 625000
        WHEN 7 THEN 3125000
        WHEN 8 THEN 15625000
        WHEN 9 THEN 78125000
        WHEN 10 THEN 390625000
        WHEN 11 THEN 1953125000
        WHEN 12 THEN 9765625000
    END,
    defense = 0,
    range_squares = 0,
    cooldown_rounds = 0;

-- Particle Cannon
ALTER TABLE particle_cannon_levels
ADD COLUMN attack INTEGER NOT NULL DEFAULT 5000,
ADD COLUMN hp BIGINT NOT NULL DEFAULT 10000,
ADD COLUMN defense NUMERIC(5,2) NOT NULL DEFAULT 8,
ADD COLUMN range_squares INTEGER NOT NULL DEFAULT 5,
ADD COLUMN cooldown_rounds INTEGER NOT NULL DEFAULT 0; -- No cooldown (fires every round)

UPDATE particle_cannon_levels SET
    attack = CASE level
        WHEN 1 THEN 5000
        WHEN 2 THEN 20000
        WHEN 3 THEN 80000
        WHEN 4 THEN 140000
        WHEN 5 THEN 240000
        WHEN 6 THEN 400000
        WHEN 7 THEN 650000
        WHEN 8 THEN 1000000
        WHEN 9 THEN 1400000
        WHEN 10 THEN 1800000
        WHEN 11 THEN 1900000
        WHEN 12 THEN 2000000
    END,
    hp = CASE level
        WHEN 1 THEN 10000
        WHEN 2 THEN 50000
        WHEN 3 THEN 250000
        WHEN 4 THEN 1000000
        WHEN 5 THEN 3000000
        WHEN 6 THEN 7000000
        WHEN 7 THEN 12000000
        WHEN 8 THEN 15000000
        WHEN 9 THEN 17000000
        WHEN 10 THEN 18500000
        WHEN 11 THEN 19250000
        WHEN 12 THEN 20000000
    END,
    defense = 8,
    range_squares = CASE level
        WHEN 1 THEN 5
        WHEN 2 THEN 6
        WHEN 3 THEN 7
        WHEN 4 THEN 8
        WHEN 5 THEN 9
        WHEN 6 THEN 10
        WHEN 7 THEN 11
        WHEN 8 THEN 12
        WHEN 9 THEN 13
        WHEN 10 THEN 14
        WHEN 11 THEN 15
        WHEN 12 THEN 16
    END,
    cooldown_rounds = 0;

-- Anti-Aircraft Gun
ALTER TABLE anti_aircraft_gun_levels
ADD COLUMN attack INTEGER NOT NULL DEFAULT 25000,
ADD COLUMN hp BIGINT NOT NULL DEFAULT 20000,
ADD COLUMN defense NUMERIC(5,2) NOT NULL DEFAULT 10,
ADD COLUMN range_squares INTEGER NOT NULL DEFAULT 3,
ADD COLUMN cooldown_rounds INTEGER NOT NULL DEFAULT 1;  -- 1 round cooldown

UPDATE anti_aircraft_gun_levels SET
    attack = CASE level
        WHEN 1 THEN 25000
        WHEN 2 THEN 100000
        WHEN 3 THEN 200000
        WHEN 4 THEN 400000
        WHEN 5 THEN 700000
        WHEN 6 THEN 1000000
        WHEN 7 THEN 1500000
        WHEN 8 THEN 2000000
        WHEN 9 THEN 2400000
        WHEN 10 THEN 2700000
        WHEN 11 THEN 2850000
        WHEN 12 THEN 3000000
    END,
    hp = CASE level
        WHEN 1 THEN 20000
        WHEN 2 THEN 100000
        WHEN 3 THEN 500000
        WHEN 4 THEN 2000000
        WHEN 5 THEN 7000000
        WHEN 6 THEN 15000000
        WHEN 7 THEN 25000000
        WHEN 8 THEN 35000000
        WHEN 9 THEN 42000000
        WHEN 10 THEN 46000000
        WHEN 11 THEN 48000000
        WHEN 12 THEN 50000000
    END,
    defense = 10,
    range_squares = CASE level
        WHEN 1 THEN 3
        WHEN 2 THEN 4
        WHEN 3 THEN 5
        WHEN 4 THEN 6
        WHEN 5 THEN 7
        WHEN 6 THEN 8
        WHEN 7 THEN 9
        WHEN 8 THEN 10
        WHEN 9 THEN 11
        WHEN 10 THEN 12
        WHEN 11 THEN 13
        WHEN 12 THEN 14
    END,
    cooldown_rounds = 1;

-- Thor's Cannon
ALTER TABLE thors_cannon_levels
ADD COLUMN attack INTEGER NOT NULL DEFAULT 10000,
ADD COLUMN hp BIGINT NOT NULL DEFAULT 200000,
ADD COLUMN defense NUMERIC(5,2) NOT NULL DEFAULT 40,
ADD COLUMN range_squares INTEGER NOT NULL DEFAULT 10,
ADD COLUMN cooldown_rounds INTEGER NOT NULL DEFAULT 2;  -- 2 round cooldown (1 round at Lv10+)

UPDATE thors_cannon_levels SET
    attack = CASE level
        WHEN 1 THEN 10000
        WHEN 2 THEN 50000
        WHEN 3 THEN 250000
        WHEN 4 THEN 600000
        WHEN 5 THEN 1200000
        WHEN 6 THEN 2000000
        WHEN 7 THEN 3000000
        WHEN 8 THEN 4000000
        WHEN 9 THEN 5000000
        WHEN 10 THEN 5500000
        WHEN 11 THEN 5750000
        WHEN 12 THEN 6000000
    END,
    hp = CASE level
        WHEN 1 THEN 200000
        WHEN 2 THEN 1000000
        WHEN 3 THEN 5000000
        WHEN 4 THEN 15000000
        WHEN 5 THEN 30000000
        WHEN 6 THEN 45000000
        WHEN 7 THEN 60000000
        WHEN 8 THEN 72000000
        WHEN 9 THEN 81000000
        WHEN 10 THEN 85500000
        WHEN 11 THEN 87750000
        WHEN 12 THEN 90000000
    END,
    defense = 40,
    range_squares = CASE level
        WHEN 1 THEN 10
        WHEN 2 THEN 11
        WHEN 3 THEN 12
        WHEN 4 THEN 13
        WHEN 5 THEN 14
        WHEN 6 THEN 15
        WHEN 7 THEN 16
        WHEN 8 THEN 18
        WHEN 9 THEN 20
        WHEN 10 THEN 21
        WHEN 11 THEN 22
        WHEN 12 THEN 23
    END,
    cooldown_rounds = CASE
        WHEN level >= 10 THEN 1  -- Cooldown drops to 1 at Lv10
        ELSE 2
    END;
```

**NEW TABLE: `celestial_base_levels`**
```sql
-- Celestial Base is listed in scope but not documented in GO2 wiki
-- Assuming it's a support/utility building rather than direct attack

CREATE TABLE celestial_base_levels (
    level INTEGER PRIMARY KEY CHECK (level BETWEEN 1 AND 12),
    space_station_req INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL,
    hp BIGINT NOT NULL DEFAULT 1000000,        -- High HP, support building
    defense NUMERIC(5,2) NOT NULL DEFAULT 20,  -- Moderate defense
    support_bonus_pct INTEGER NOT NULL DEFAULT 5  -- +5-60% defense value boost to other defenses
);

-- Seed data (12 levels)
INSERT INTO celestial_base_levels (level, space_station_req, metal_cost, he3_cost, gold_cost, build_time_seconds, hp, defense, support_bonus_pct) VALUES
(1,  6,  600,   500,   600,   200,  1000000,   20,  5),
(2,  6,  1818,  1515,  1818,  606,  5000000,   25,  10),
(3,  7,  5509,  4590,  5509,  1836, 25000000,  30,  15),
(4,  7,  16691, 13905, 16691, 5562, 100000000, 35,  20),
(5,  8,  50564, 42136, 50564, 16857,300000000, 40,  25),
(6,  8,  153207,127672,153207,51060,700000000, 45,  30),
(7,  9,  464419,386849,464419,154665,1500000000,50, 35),
(8,  9,  1407446,1172871,1407446,468332,3000000000,55,40),
(9,  10, 4264891,3555226,4264891,1419480,6000000000,60,45),
(10, 10, 12921053,10774965,12921053,4301822,12000000000,65,50),
(11, 11, 39161513,32647277,39161513,13031301,24000000000,70,55),
(12, 11, 118646926,98919212,118646926,39458671,48000000000,75,60);

CREATE INDEX idx_celestial_base_levels ON celestial_base_levels(level);
ALTER TABLE celestial_base_levels ENABLE ROW LEVEL SECURITY;
CREATE POLICY celestial_base_levels_select_all ON celestial_base_levels FOR SELECT USING (true);
```

**Purpose:** Celestial Base provides passive defense boost to all other defense buildings.

**ALTER `building_types`** - Fix Celestial Base category:
```sql
UPDATE building_types
SET category = 'defense'
WHERE id = 'celestial_base';
```

**Time Estimate:** 4-5 hours (migration creation + testing)

---

## 3. GALAXY ONLINE 2 DEFENSE MECHANICS RESEARCH

### 3.1 Defense Building Stats Summary

| Building | Attack | HP (Lv1/Lv12) | Range | Cooldown | Special |
|----------|--------|---------------|-------|----------|---------|
| **Meteor Star** | 50 (constant) | 200 / 90M | 0 (melee) | None | Blocks movement, high HP scaling |
| **Particle Cannon** | 5K / 2M | 10K / 20M | 5-16 | None | Fires every round, longest range |
| **Anti-Aircraft Gun** | 25K / 3M | 20K / 50M | 3-14 | 1 round | AoE damage, moderate power |
| **Thor's Cannon** | 10K / 6M | 200K / 90M | 10-23 | 2 rounds (1 at Lv10+) | Highest attack, hits all in range |
| **Celestial Base** | N/A | 1M / 48B | N/A | N/A | Support building, +5-60% defense boost |

### 3.2 Combat Integration

**From GO2 wiki:** "All enemy ships MUST defeat all Orbital Defenses before they can start attacking your Space Station."

**Defense Combat Phase:**
1. Attacker fleet arrives
2. **Defense buildings attack first** (before fleet vs fleet)
3. Defense buildings attack in order (Particle Cannon → AA Gun → Thor's Cannon → Meteor Star)
4. Each building attacks based on range + cooldown
5. If all defenses destroyed, proceed to fleet vs defense fleet (Module 6)
6. If attacker destroyed by defenses, defender wins

**Attack Priority:**
- Particle Cannon: No cooldown, fires every round
- Anti-Aircraft Gun: 1 round cooldown (fires every other round)
- Thor's Cannon: 2 round cooldown (1 at Lv10+), AOE damage to all ships in range
- Meteor Star: Melee (range 0), forces engagement, very high HP

**Cooldown Mechanic:**
- Round 1: All defenses attack
- Round 2: Particle Cannon attacks, AA Gun on cooldown, Thor on cooldown
- Round 3: Particle Cannon + AA Gun attack, Thor on cooldown
- Round 4: All defenses attack again

### 3.3 Tech Bonuses

**Planetary Defense Science Techs:**
1. **Energy Control** (Lv1-10): -1-10% defense construction costs
2. **Rapid Defense Buildup** (Lv1-10): +1-10% defense construction speed
3. **Defense Enhancement** (Lv1-10): +1-10% defensive value (HP)
4. **Emplacement Mastery** (Lv1-10): +1-10% emplacement attack power
5. **Utmost Defense Buildup** (Lv1-10): +1-10% max defensive structures
6. **Range Extension** (Lv1-2): Increases range of Particle Cannon + AA Gun
7. **Thor Buildup** (Lv1): +1 max Thor's Cannon allowed (base max = 3, can go to 4)
8. **Augment Propulsion** (Lv1-2): +1-2 ship movement speed when defending

**Applied in:**
- Construction: Cost reduction, build speed
- Combat: Attack bonus, defense value (HP), range
- Placement: Max count limits

---

## 4. COMBAT INTEGRATION DESIGN

### 4.1 Defense Phase in PvP Combat

**Modified PvP Combat Flow (Module 6 + Module 7):**

```go
func ResolvePvPCombat(attackFleetID, defenderPlayerID string) {
    // 1. Load attacker fleet
    attackerFleet := LoadFleet(attackFleetID)

    // 2. Load defender defenses
    defenses := LoadDefenseBuildings(defenderPlayerID)

    // 3. Defense Combat Phase (NEW)
    if len(defenses) > 0 {
        defenseResult := ResolveDefenseCombat(attackerFleet, defenses)

        if defenseResult.AttackerDestroyed {
            // Attacker defeated by defenses
            SaveCombatReport(defenseResult, "defender_win")
            return
        }

        // Update attacker fleet casualties
        attackerFleet = defenseResult.AttackerFinalState
    }

    // 4. Fleet vs Fleet Combat (Module 5 engine)
    defenseFleetID := GetDefenseFleet(defenderPlayerID)

    if defenseFleetID == "" {
        // No defense fleet, attacker wins
        loot := CalculateLoot(defenderPlayerID)
        GrantLoot(attackerFleet.PlayerID, loot)
        SaveCombatReport(attackerFleet, nil, "attacker_win", loot)
        return
    }

    // 5. Standard fleet combat
    result := services.ResolveCombat(attackerFleet.ID, defenseFleetID, "pvp")
    // ... rest of Module 6 combat logic
}
```

### 4.2 Defense Combat Resolution

```go
type DefenseBuilding struct {
    ID               string
    BuildingType     string
    Level            int
    Attack           int
    HP               int64
    Defense          float64
    Range            int
    Cooldown         int
    LastFireRound    int  // Track cooldown
    Position         Position
}

type DefenseCombatResult struct {
    TotalRounds        int
    AttackerDestroyed  bool
    DefensesDestroyed  bool
    AttackerFinalState *FleetState
    DefensesFinalState []DefenseBuilding
    RoundsLog          []DefenseRoundData
}

func ResolveDefenseCombat(attackerFleet *FleetState, defenses []DefenseBuilding) *DefenseCombatResult {
    result := &DefenseCombatResult{
        AttackerFinalState: attackerFleet,
        DefensesFinalState: defenses,
        RoundsLog: []DefenseRoundData{},
    }

    // Apply tech bonuses to defenses
    techBonus := GetPlayerTechBonuses(defenses[0].PlayerID)
    for i := range defenses {
        defenses[i].Attack = int(float64(defenses[i].Attack) * (1.0 + techBonus.EmplacementAttack/100.0))
        defenses[i].HP = int64(float64(defenses[i].HP) * (1.0 + techBonus.DefenseValue/100.0))
    }

    // Apply Celestial Base support bonus
    celestialBonus := getCelestialBaseBonus(defenses)
    if celestialBonus > 0 {
        for i := range defenses {
            if defenses[i].BuildingType != "celestial_base" {
                defenses[i].Attack = int(float64(defenses[i].Attack) * (1.0 + float64(celestialBonus)/100.0))
                defenses[i].HP = int64(float64(defenses[i].HP) * (1.0 + float64(celestialBonus)/100.0))
            }
        }
    }

    // Combat loop (max 99 rounds)
    for round := 1; round <= 99; round++ {
        roundData := DefenseRoundData{Round: round}

        // Defense buildings attack (in order: Particle → AA → Thor → Meteor)
        attackOrder := sortDefensesByAttackPriority(defenses)

        for _, defense := range attackOrder {
            if defense.HP <= 0 {
                continue // Building destroyed
            }

            // Check cooldown
            if round - defense.LastFireRound < defense.Cooldown {
                continue // On cooldown
            }

            // Attack attacker fleet
            damage := calculateDefenseDamage(defense, attackerFleet)
            casualties := applyFleetDamage(attackerFleet, damage)

            defense.LastFireRound = round

            roundData.DefenseAttacks = append(roundData.DefenseAttacks, DefenseAttack{
                BuildingType: defense.BuildingType,
                Damage: damage,
                Casualties: casualties,
            })
        }

        // Attacker fleet attacks back (simplified - just target random defense)
        if attackerFleet.TotalShips() > 0 {
            target := selectRandomAliveDefense(defenses)
            if target != nil {
                damage := calculateFleetVsDefenseDamage(attackerFleet, target)
                target.HP -= damage

                if target.HP <= 0 {
                    target.HP = 0
                    roundData.DefensesDestroyed = append(roundData.DefensesDestroyed, target.BuildingType)
                }

                roundData.AttackerDamage = damage
            }
        }

        roundData.AttackerShipsRemaining = attackerFleet.TotalShips()
        roundData.DefensesRemaining = countAliveDefenses(defenses)

        result.RoundsLog = append(result.RoundsLog, roundData)

        // Check victory conditions
        if attackerFleet.TotalShips() == 0 {
            result.AttackerDestroyed = true
            result.TotalRounds = round
            return result
        }

        if countAliveDefenses(defenses) == 0 {
            result.DefensesDestroyed = true
            result.TotalRounds = round
            return result
        }
    }

    // Round 99 reached, attacker survives = attacker wins
    result.DefensesDestroyed = true
    result.TotalRounds = 99
    return result
}

func calculateDefenseDamage(defense DefenseBuilding, fleet *FleetState) int {
    baseDamage := defense.Attack

    // Range check (simplified - assume all ships in range for now)
    // Future: Calculate actual range based on grid positions

    // Thor's Cannon hits all ships
    if defense.BuildingType == "thors_cannon" {
        return baseDamage // AOE damage distributed across all ships
    }

    // Other buildings hit single stack
    return baseDamage
}

func sortDefensesByAttackPriority(defenses []DefenseBuilding) []DefenseBuilding {
    // Order: Particle Cannon → AA Gun → Thor → Meteor → Celestial Base (no attack)
    order := map[string]int{
        "particle_cannon":     1,
        "anti_aircraft_gun":   2,
        "thors_cannon":        3,
        "meteor_star":         4,
        "celestial_base":      99, // Support building, doesn't attack
    }

    sorted := make([]DefenseBuilding, len(defenses))
    copy(sorted, defenses)

    sort.Slice(sorted, func(i, j int) bool {
        return order[sorted[i].BuildingType] < order[sorted[j].BuildingType]
    })

    return sorted
}

func getCelestialBaseBonus(defenses []DefenseBuilding) int {
    for _, d := range defenses {
        if d.BuildingType == "celestial_base" && d.HP > 0 {
            // Get support bonus from level (5% per level, up to 60% at Lv12)
            return d.Level * 5
        }
    }
    return 0
}
```

**Time Estimate:** 8-10 hours

---

## 5. BACKEND IMPLEMENTATION

### 5.1 Service Layer (`backend/internal/services/defense_combat.go`)

**Functions:**
```go
// LoadDefenseBuildings gets all defense buildings for a player
func LoadDefenseBuildings(playerID string) ([]DefenseBuilding, error)

// ResolveDefenseCombat executes defense buildings vs attacking fleet
func ResolveDefenseCombat(attackerFleet *FleetState, defenses []DefenseBuilding) (*DefenseCombatResult, error)

// calculateDefenseDamage computes damage from a defense building
func calculateDefenseDamage(defense DefenseBuilding, fleet *FleetState) int

// applyDefenseDamage applies damage to attacking fleet
func applyDefenseDamage(fleet *FleetState, damage int) int

// sortDefensesByAttackPriority orders defenses for attack sequence
func sortDefensesByAttackPriority(defenses []DefenseBuilding) []DefenseBuilding

// getCelestialBaseBonus calculates support bonus from Celestial Base
func getCelestialBaseBonus(defenses []DefenseBuilding) int

// UpdateDefenseBuildingHP saves defense HP after combat
func UpdateDefenseBuildingHP(defenses []DefenseBuilding) error
```

**Time Estimate:** 6-8 hours

---

### 5.2 Update PvP Combat Handler

**Modify `backend/internal/handlers/pvp.go`:**

```go
func ResolvePvPCombat(attackFleetID, defenderPlayerID, movementID string) {
    // 1. Load attacker fleet
    attackerFleet := LoadFleet(attackFleetID)

    // 2. Load defender defenses (NEW)
    defenses, err := services.LoadDefenseBuildings(defenderPlayerID)
    if err != nil {
        log.Printf("Failed to load defenses: %v", err)
    }

    // 3. Defense Combat Phase (NEW)
    var defenseResult *services.DefenseCombatResult
    if len(defenses) > 0 {
        defenseResult, err = services.ResolveDefenseCombat(attackerFleet, defenses)

        if err != nil {
            log.Printf("Defense combat failed: %v", err)
        } else {
            // Save defense combat results
            SaveDefenseCombatLog(defenseResult)

            // Update defense building HP
            services.UpdateDefenseBuildingHP(defenseResult.DefensesFinalState)

            if defenseResult.AttackerDestroyed {
                // Attacker defeated by defenses alone
                SaveCombatReport(attackFleetID, defenderPlayerID, "defender_win", defenseResult)
                NotifyPlayer(attackerFleet.PlayerID, "Your fleet was destroyed by orbital defenses!")
                NotifyPlayer(defenderPlayerID, "Your orbital defenses defeated the attacker!")
                InitiateReturnJourney(attackFleetID, movementID) // Return survivors (if any)
                return
            }

            // Attacker survived defenses, update fleet state
            attackerFleet = defenseResult.AttackerFinalState
        }
    }

    // 4. Continue with fleet vs fleet combat (Module 6 existing code)
    defenseFleetID := GetDefenseFleet(defenderPlayerID)
    // ... rest of Module 6 combat logic
}
```

**Time Estimate:** 2-3 hours

---

## 6. FRONTEND IMPLEMENTATION

### 6.1 3D Models for Defense Buildings

**Files:** `frontend/src/components/three/buildings/*.tsx`

Create 5 procedural 3D models:
1. **MeteorStar.tsx** - Spherical structure with spikes, metallic texture
2. **ParticleCannon.tsx** - Cylindrical barrel with energy rings, long range
3. **AntiAircraftGun.tsx** - Dual-barrel turret, rotating base
4. **ThorsCannon.tsx** - Massive quad-barrel cannon, imposing structure
5. **CelestialBase.tsx** - Orbital platform with support beams, glowing core

**Example: MeteorStar.tsx**
```typescript
import React, { useRef } from 'react';
import { Mesh } from 'three';
import { useFrame } from '@react-three/fiber';

interface Props {
  position: [number, number, number];
  level: number;
}

export function MeteorStar({ position, level }: Props) {
  const meshRef = useRef<Mesh>(null);

  // Rotate slowly
  useFrame(() => {
    if (meshRef.current) {
      meshRef.current.rotation.y += 0.005;
    }
  });

  const size = 0.8 + (level * 0.05); // Scale with level

  return (
    <group position={position}>
      {/* Main sphere */}
      <mesh ref={meshRef}>
        <sphereGeometry args={[size, 16, 16]} />
        <meshStandardMaterial color="#8b7355" metalness={0.7} roughness={0.3} />
      </mesh>

      {/* Spikes (8 total) */}
      {[...Array(8)].map((_, i) => {
        const angle = (i / 8) * Math.PI * 2;
        const x = Math.cos(angle) * size * 0.8;
        const z = Math.sin(angle) * size * 0.8;

        return (
          <mesh key={i} position={[x, 0, z]}>
            <coneGeometry args={[0.2, 0.6, 8]} />
            <meshStandardMaterial color="#6b5345" metalness={0.8} roughness={0.2} />
          </mesh>
        );
      })}

      {/* Base platform */}
      <mesh position={[0, -size * 0.6, 0]}>
        <cylinderGeometry args={[size * 1.2, size * 1.4, 0.3, 16]} />
        <meshStandardMaterial color="#4a4a4a" metalness={0.5} roughness={0.5} />
      </mesh>
    </group>
  );
}
```

**Similar procedural generation for other 4 buildings.**

**Time Estimate:** 6-8 hours (1-1.5h per building)

---

### 6.2 Building Info Display

**Update `BuildingModel.tsx` factory:**
```typescript
// Add defense building dispatch
case 'meteor_star':
  return <MeteorStar position={worldPos} level={building.level} />;
case 'particle_cannon':
  return <ParticleCannon position={worldPos} level={building.level} />;
case 'anti_aircraft_gun':
  return <AntiAircraftGun position={worldPos} level={building.level} />;
case 'thors_cannon':
  return <ThorsCannon position={worldPos} level={building.level} />;
case 'celestial_base':
  return <CelestialBase position={worldPos} level={building.level} />;
```

**Update building tooltips to show defense stats:**
```typescript
// In BuildingTooltip.tsx
if (building.category === 'defense') {
  const stats = await api.getDefenseStats(building.type, building.level);

  return (
    <div className="defense-stats">
      <div>Attack: {stats.attack.toLocaleString()}</div>
      <div>HP: {stats.hp.toLocaleString()}</div>
      <div>Range: {stats.range} squares</div>
      <div>Cooldown: {stats.cooldown} rounds</div>
      {stats.support_bonus && (
        <div>Support Bonus: +{stats.support_bonus}% to other defenses</div>
      )}
    </div>
  );
}
```

**Time Estimate:** 2-3 hours

---

### 6.3 Defense Combat Report Display

**Update `CombatReportPanel.tsx`:**

```typescript
// Add defense phase section
function renderDefensePhase(report: CombatReport) {
  if (!report.defense_phase) return null;

  return (
    <div className="defense-phase-section">
      <h3>Orbital Defense Phase</h3>
      <p>Rounds: {report.defense_phase.total_rounds}</p>

      <table className="defense-rounds-table">
        <thead>
          <tr>
            <th>Round</th>
            <th>Defense Attacks</th>
            <th>Attacker Ships Lost</th>
            <th>Defenses Destroyed</th>
          </tr>
        </thead>
        <tbody>
          {report.defense_phase.rounds.map((round) => (
            <tr key={round.round}>
              <td>{round.round}</td>
              <td>
                {round.defense_attacks.map((atk) => (
                  <div key={atk.building_type}>
                    {atk.building_type}: {atk.damage.toLocaleString()} dmg
                  </div>
                ))}
              </td>
              <td>{round.casualties}</td>
              <td>
                {round.defenses_destroyed.map((def) => (
                  <div key={def}>{def} destroyed</div>
                ))}
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {report.defense_phase.attacker_destroyed && (
        <div className="result defeat">
          Fleet destroyed by orbital defenses!
        </div>
      )}
    </div>
  );
}
```

**Time Estimate:** 2 hours

---

## 7. PLACEMENT VALIDATION

### 7.1 Defense Building Limits

**Backend validation:**
```go
func ValidateDefensePlacement(playerID, buildingType string) error {
    // Get tech bonuses
    techBonus := GetPlayerTechBonuses(playerID)

    // Base max counts
    maxCounts := map[string]int{
        "meteor_star":       40,  // Base max from building_types
        "particle_cannon":   20,
        "anti_aircraft_gun": 15,
        "thors_cannon":      3,   // Special limit
        "celestial_base":    1,   // Only 1 allowed
    }

    // Apply tech bonuses
    if buildingType != "thors_cannon" && buildingType != "celestial_base" {
        // Utmost Defense Buildup increases max structures
        maxCounts[buildingType] += int(float64(maxCounts[buildingType]) * techBonus.MaxDefenseStructures / 100.0)
    }

    // Thor Buildup adds +1 to Thor's Cannon
    if buildingType == "thors_cannon" {
        maxCounts["thors_cannon"] += techBonus.MaxThorCannon
    }

    // Check current count
    var currentCount int
    db.QueryRow(`
        SELECT COUNT(*) FROM buildings
        WHERE player_id = $1 AND building_type = $2
    `, playerID, buildingType).Scan(&currentCount)

    if currentCount >= maxCounts[buildingType] {
        return fmt.Errorf("max %s limit reached (%d)", buildingType, maxCounts[buildingType])
    }

    return nil
}
```

**Time Estimate:** 2 hours

---

## 8. QA CHECKLIST

### 8.1 Database Tests
- [ ] All 4 `*_levels` tables have combat stats (attack, HP, range, cooldown)
- [ ] `celestial_base_levels` table created with 12 levels
- [ ] Celestial Base category changed to 'defense'
- [ ] Stats match GO2 values (scaled appropriately)
- [ ] 12 Planetary Defense Science techs seeded

### 8.2 Combat Integration Tests
- [ ] Defense buildings load correctly for defender
- [ ] Defense combat phase executes before fleet combat
- [ ] Defense buildings attack in order (Particle → AA → Thor → Meteor)
- [ ] Cooldown mechanics work (AA Gun, Thor's Cannon)
- [ ] Thor's Cannon AOE damage hits all ships
- [ ] Celestial Base support bonus applies (+5-60%)
- [ ] Tech bonuses apply (emplacement attack, defense value)
- [ ] Attacker defeated by defenses = defender win
- [ ] Defenses destroyed = proceed to fleet combat
- [ ] Defense building HP updated after combat
- [ ] Destroyed defenses removed/marked

### 8.3 Placement Tests
- [ ] Max defense structure limits enforced
- [ ] Utmost Defense Buildup tech increases max count
- [ ] Thor Buildup tech adds +1 Thor's Cannon max
- [ ] Celestial Base limit = 1 (only one allowed)
- [ ] Building placement on 20×20 grid works
- [ ] Defense buildings visible on planet view

### 8.4 Frontend Tests
- [ ] 5 3D models render correctly
- [ ] Models scale with building level
- [ ] Defense stats display in tooltips
- [ ] Combat reports show defense phase
- [ ] Defense rounds table displays correctly
- [ ] "Destroyed by defenses" message shown

### 8.5 Tech Bonus Tests
- [ ] Energy Control reduces defense costs
- [ ] Rapid Defense Buildup increases build speed
- [ ] Defense Enhancement increases HP
- [ ] Emplacement Mastery increases attack
- [ ] Range Extension increases Particle/AA range
- [ ] Thor Buildup allows 4th Thor's Cannon

### 8.6 Edge Cases
- [ ] No defense buildings = skip defense phase
- [ ] All defenses destroyed by attacker
- [ ] Attacker fleet destroyed by defenses
- [ ] Cooldown persists across rounds correctly
- [ ] Round 99 reached in defense combat
- [ ] Thor's Cannon cooldown drops to 1 at Lv10
- [ ] Celestial Base destroyed = bonus removed mid-combat

**Total QA Checks:** 37

---

## 9. RISK ASSESSMENT

### 9.1 Technical Risks

**RISK 1: Combat Balance (Defenses Too Strong/Weak)**
- **Severity:** MEDIUM
- **Impact:** Defenses too strong = no successful attacks, too weak = useless
- **Mitigation:** Use GO2 stat values, playtesting for balance
- **Adjustment:** Can scale stats by 0.5× or 2× if needed

**RISK 2: Cooldown Logic Complexity**
- **Severity:** LOW
- **Impact:** Cooldown tracking across rounds may have bugs
- **Mitigation:** Simple LastFireRound tracking, well-tested in GO2
- **Test:** Explicit test cases for AA Gun (1 round) and Thor (2 rounds)

**RISK 3: AOE Damage Distribution**
- **Severity:** LOW
- **Impact:** Thor's Cannon AOE may be too powerful or too weak
- **Mitigation:** Distribute damage across all ships, not per-ship
- **Formula:** Total AOE damage = Thor attack / total attacker ships

---

### 9.2 Balance Risks

**RISK 4: Celestial Base Support Bonus Too High**
- **Severity:** MEDIUM
- **Impact:** +60% at Lv12 may make defenses unbeatable
- **Mitigation:** Test with/without Celestial Base
- **Adjustment:** Can reduce to +30% max if too strong

**RISK 5: Tech Bonuses Stack Too High**
- **Severity:** MEDIUM
- **Impact:** +10% emplacement + +60% Celestial Base = +70% total
- **Mitigation:** Bonuses are multiplicative, not additive
- **Formula:** Final Attack = Base × (1 + Tech/100) × (1 + Celestial/100)

---

### 9.3 UX Risks

**RISK 6: Defense Combat Too Long**
- **Severity:** LOW
- **Impact:** Defense phase adds 10-20 rounds to combat
- **Mitigation:** This is expected, matches GO2
- **Alternative:** Show "Defense Phase" summary, skip detailed rounds

---

## 10. IMPLEMENTATION TIMELINE

### 10.1 Phase Breakdown

**Phase 1: Database Migration (Day 1)**
- Add combat stats to 4 `*_levels` tables
- Create `celestial_base_levels` table
- Seed stats with GO2 values
- Update Celestial Base category
- **Time:** 4-5 hours

**Phase 2: Defense Combat Engine (Days 1-2)**
- Create `defense_combat.go` service
- Implement combat resolution
- Cooldown tracking
- AOE damage logic
- Tech bonus integration
- **Time:** 8-10 hours

**Phase 3: PvP Integration (Day 2)**
- Update PvP combat handler
- Load defense buildings
- Execute defense phase
- Update defense HP after combat
- **Time:** 2-3 hours

**Phase 4: 3D Models (Day 3)**
- Create 5 procedural models
- Integrate with BuildingModel factory
- Test rendering on planet grid
- **Time:** 6-8 hours

**Phase 5: Frontend UI (Day 3)**
- Defense stats tooltips
- Combat report defense phase display
- Placement validation UI
- **Time:** 3-4 hours

**Phase 6: QA & Balance (Day 4)**
- Run QA checklist (37 checks)
- Balance testing (attack multiple times, test defenses)
- Bug fixes
- **Time:** 4-6 hours

---

### 10.2 Critical Path

```
Day 1: Database Migration → Defense Combat Engine Core
       ↓
Day 2: Combat Engine Complete → PvP Integration
       ↓
Day 3: 3D Models → Frontend UI
       ↓
Day 4: QA Testing → Balance → Bug Fixes
```

**Total Estimated Time:** 15-22 hours (3-4 working days)

---

## 11. SUCCESS CRITERIA

Module 7 is considered **COMPLETE** when:

✅ **Database:**
- All defense building level tables have combat stats
- Celestial Base level table created
- Stats match GO2 values (scaled)

✅ **Combat:**
- Defense buildings attack before fleet combat
- Cooldown mechanics work correctly
- Thor's Cannon AOE damage functional
- Celestial Base support bonus applies
- Tech bonuses integrate correctly
- Attacker can be defeated by defenses alone

✅ **Frontend:**
- 5 3D models render on planet grid
- Defense stats display in tooltips
- Combat reports show defense phase
- Placement validation enforces limits

✅ **QA:**
- All 37 QA checks pass
- Combat balance feels fair (defenses help but not unbeatable)
- No critical bugs

---

## 12. NEXT STEPS AFTER MODULE 7

After Module 7 completion, defense buildings strengthen PvP defense. Next modules:

1. **Module 8: Recycling Plant** - Scrap unwanted ships for resources
2. **Module 9: Inventory System** - Full item management (resource packs, truce cards, commander cards)
3. **Module 10: World Chat** - Basic global chat channel
4. **Module 11: Production Polish** - Improve UX (animations, sound, tooltips, tutorial)

**Recommended Next Module:** Module 8 (Recycling Plant) - Simple utility feature.

---

## APPENDIX A: DEFENSE BUILDING STATS TABLE

| Building | Lv1 Attack | Lv12 Attack | Lv1 HP | Lv12 HP | Range Lv1 | Range Lv12 | Cooldown |
|----------|------------|-------------|--------|---------|-----------|------------|----------|
| Meteor Star | 50 | 50 | 200 | 9.7B | 0 | 0 | 0 |
| Particle Cannon | 5K | 2M | 10K | 20M | 5 | 16 | 0 |
| Anti-Aircraft Gun | 25K | 3M | 20K | 50M | 3 | 14 | 1 |
| Thor's Cannon | 10K | 6M | 200K | 90M | 10 | 23 | 2 (1 at Lv10+) |
| Celestial Base | N/A | N/A | 1M | 48B | N/A | N/A | N/A |

---

## APPENDIX B: References

**Galaxy Online 2 Wiki Sources:**
- [Composite Orbital Defenses Table](https://galaxyonlineii.fandom.com/wiki/Composite_Orbital_Defenses_Table) - Defense stats
- [Meteor Star](https://galaxyonlineii.fandom.com/wiki/Meteor_Star) - Blocking mechanics
- [Particle Cannon](https://galaxyonlineii.fandom.com/wiki/Particle_Cannon) - No cooldown, longest range
- [Thor's Cannon](https://galaxyonlineii.fandom.com/wiki/Thor's_Cannon) - AOE damage, 2-round cooldown

**Project Documents:**
- `/docs/planning/final-scope.md` - Module 7 scope
- `/docs/planning/module-5-combat-system-plan.md` - Combat engine (reused)
- `/docs/planning/module-6-pvp-combat-plan.md` - PvP integration point
- `/supabase/migrations/20260206005232_phase1_mvp.sql` - Existing defense building tables

---

**END OF MODULE 7 IMPLEMENTATION PLAN**
