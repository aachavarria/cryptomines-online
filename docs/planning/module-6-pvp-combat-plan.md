# Module 6: PvP Combat - Implementation Plan

**Date:** 2026-02-07
**Architect:** architect
**Module:** PvP Combat (Player vs Player Attacks)
**Estimated Total:** 25-35 hours (4-6 days)

---

## 1. EXECUTIVE SUMMARY

### 1.1 Module Overview

The PvP Combat System enables **player-vs-player warfare**, allowing players to attack neighbor planets, loot resources, and compete for dominance. This module builds directly on the **Module 5 Combat Engine** foundation but adds travel mechanics, neighbor selection, defense fleets, Radar warnings, and protection systems.

**Complexity Level:** HIGH (★★★★☆)
- Reuses Module 5 combat engine (80% shared code)
- New: Fleet travel system, neighbor discovery, defense mechanics
- Integration: Radar building, Truce cards, resource looting

### 1.2 Scope (From final-scope.md)

**IN SCOPE:**
- ✅ Attack neighbor planets (select nearby players)
- ✅ Fleet travel time (uses SP - Space Points for fuel/travel)
- ✅ Combat resolution (reuse Module 5 engine with 'pvp' combat_type)
- ✅ Loot 20% of defender's **harvested resources** (excluding warehouse)
- ✅ Radar building integration (detect incoming attacks, warning system)
- ✅ Defense fleets (assign fleet to defend planet)
- ✅ PvP combat reports (attacker + defender receive report)
- ✅ Truce cards (12h/72h protection shields)

**OUT OF SCOPE (deferred to future modules):**
- ❌ Space Defense Buildings combat integration (Module 7)
- ❌ Corps/Alliance warfare mechanics (cut from scope)
- ❌ Galaxy map visualization (basic neighbor list instead)
- ❌ Revenge attacks (future enhancement)
- ❌ Attack cooldowns per target (unlimited attacks for now)
- ❌ Honor/Ranking system (future enhancement)

### 1.3 Current Implementation Status

**Database:** 40% Complete
- ✅ `combat_reports` table supports 'pvp' combat_type
- ✅ `radar` building exists with 9 levels
- ✅ `radar_levels` lookup table exists
- ✅ `fleets` table exists
- ❌ NO `fleet_movements` table for tracking traveling fleets
- ❌ NO `defense_assignments` table for defense fleets
- ❌ NO `active_protections` table for truce cards
- ❌ NO neighbor discovery system

**Backend:** 0% Complete
- ❌ NO PvP attack endpoint
- ❌ NO neighbor discovery logic
- ❌ NO fleet travel system
- ❌ NO defense fleet mechanics
- ❌ NO Radar warning system
- ❌ NO truce card validation

**Frontend:** 0% Complete
- ❌ NO galaxy/neighbors UI
- ❌ NO attack interface
- ❌ NO incoming attack warnings
- ❌ NO defense fleet assignment UI

---

## 2. DATABASE ANALYSIS

### 2.1 Existing Schema Review

**Table: `combat_reports` (existing - Phase 2)**
```sql
combat_reports:
  combat_type TEXT CHECK (... 'pvp' ...)  ✅ Already supports PvP
  attacker_id UUID REFERENCES players(id)
  defender_id UUID REFERENCES players(id)
  loot_json JSONB
```

**Table: `radar` building (existing - Phase 1)**
```sql
building_types:
  'radar' → Civic Center req, 9 levels
```

**Table: `radar_levels` (existing - Phase 1)**
```sql
radar_levels:
  level (1-9)
  civic_center_req
  detection_time_minutes  -- 30min (Lv1) to 270min (Lv9)
  metal_cost, he3_cost, gold_cost
  build_time_seconds
```

**Table: `fleets` (existing - Phase 2)**
```sql
fleets:
  player_id UUID
  status TEXT CHECK (status IN ('idle', 'moving', 'in_combat', 'defending'))
  commander_id UUID
  formation TEXT
  targeting_command TEXT
```

**Table: `players` (existing - Phase 1)**
```sql
players:
  id UUID
  username TEXT
  level INTEGER
  is_online BOOLEAN
  -- NO galaxy coordinates yet
```

**Table: `resources` (existing - Phase 1)**
```sql
resources:
  player_id UUID
  metal BIGINT
  he3 BIGINT
  gold BIGINT
  -- NO warehouse fields (Module 3 adds these)
```

### 2.2 Required Schema Changes

**NEW TABLE: `fleet_movements`**
```sql
CREATE TABLE fleet_movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id UUID NOT NULL REFERENCES fleets(id) ON DELETE CASCADE,
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    origin_player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    target_player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    movement_type TEXT NOT NULL CHECK (movement_type IN ('attack', 'return')),
    departure_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    arrival_time TIMESTAMPTZ NOT NULL,
    sp_consumed INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'traveling'
        CHECK (status IN ('traveling', 'arrived', 'cancelled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_fleet_movements_fleet ON fleet_movements(fleet_id);
CREATE INDEX idx_fleet_movements_arrival ON fleet_movements(arrival_time) WHERE status = 'traveling';
CREATE INDEX idx_fleet_movements_target ON fleet_movements(target_player_id) WHERE status = 'traveling';
```

**Purpose:** Track fleets in transit (outbound attacks + return trips).

**NEW TABLE: `defense_assignments`**
```sql
CREATE TABLE defense_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    fleet_id UUID UNIQUE NOT NULL REFERENCES fleets(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(player_id, fleet_id)
);

CREATE INDEX idx_defense_assignments_player ON defense_assignments(player_id);
```

**Purpose:** Mark which fleet(s) are assigned to defend the planet.

**NEW TABLE: `active_protections`**
```sql
CREATE TABLE active_protections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    protection_type TEXT NOT NULL CHECK (protection_type IN ('truce_card', 'newbie_protection', 'corp_protection')),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(player_id, protection_type)
);

CREATE INDEX idx_active_protections_player ON active_protections(player_id);
CREATE INDEX idx_active_protections_expiry ON active_protections(expires_at);
```

**Purpose:** Track active protection shields (Truce cards, newbie protection).

**ALTER TABLE `players`** - Add galaxy coordinates:
```sql
ALTER TABLE players
ADD COLUMN galaxy_sector INTEGER NOT NULL DEFAULT 1,
ADD COLUMN galaxy_x INTEGER NOT NULL DEFAULT 0,
ADD COLUMN galaxy_y INTEGER NOT NULL DEFAULT 0;

CREATE INDEX idx_players_coordinates ON players(galaxy_sector, galaxy_x, galaxy_y);
```

**Purpose:** Position players in galaxy for neighbor discovery.

**Coordinate System:**
- Galaxy divided into sectors (1-10)
- Each sector is 1000×1000 grid
- Players spawn randomly in their sector
- Neighbors = players within distance threshold

**ALTER TABLE `resources`** - Add SP (Space Points):
```sql
ALTER TABLE resources
ADD COLUMN space_points INTEGER NOT NULL DEFAULT 100,
ADD COLUMN sp_regen_at TIMESTAMPTZ NOT NULL DEFAULT now();
```

**Purpose:** SP used for fleet travel, regenerates over time.

---

## 3. GALAXY ONLINE 2 PVP MECHANICS RESEARCH

### 3.1 Attack Flow

From GO2 wiki research:

1. **Neighbor Discovery**: Players see nearby planets in galaxy view
2. **Target Selection**: Click planet to view player card (level, corps, defenses visible/hidden)
3. **Scout Attack** (optional): Send cheap ships to reveal defenses
4. **Main Attack**: Send attack fleet with SP cost
5. **Fleet Travel**: Time based on distance, consumes SP
6. **Combat Resolution**: Reuse Module 5 combat engine
7. **Loot Distribution**: Winner gets 20% of loser's **harvested resources**
8. **Return Travel**: Takes half the time of outbound trip
9. **Battle Report**: Both players receive mail with combat report

### 3.2 Space Points (SP) System

**SP Mechanics** (inferred from GO2 patterns):
- **Max SP:** 100 (base), +10 per level (max ~200 at high levels)
- **Regen Rate:** 1 SP per 6 minutes (10 SP/hour)
- **Attack Cost:** Distance-based (1 SP per 10 distance units)
- **Usage:** Consumed on fleet departure, NOT refunded

**Example:**
- Player A at (100, 100)
- Player B at (150, 200)
- Distance = √((150-100)² + (200-100)²) = √(2500 + 10000) = √12500 ≈ 112
- SP Cost = 112 / 10 = 11.2 → 12 SP (rounded up)

### 3.3 Fleet Travel Time

**Travel Time Formula:**
```
Outbound Time = Base Time × Distance Multiplier
Return Time = Outbound Time / 2
```

**Base Time:** 1 minute per 10 distance units (configurable)

**Example:**
- Distance = 112 units
- Outbound = 112 / 10 = 11.2 → 12 minutes
- Return = 6 minutes

**Speed Bonuses** (future):
- Tech bonuses (Ship Movement Speed research)
- Commander speed stat
- Formation bonuses

### 3.4 Resource Looting Rules

From GO2 wiki: "20% of the losing planet's **harvested resources**"

**Key Detail:** "Harvested resources" = resources NOT in warehouse.

**In Cryptomines Online:**
- Total Resources = `resources.metal + resources.he3 + resources.gold`
- Warehouse Resources = `resources.warehouse_metal + warehouse_he3 + warehouse_gold` (Module 3)
- **Lootable Resources = Total - Warehouse**
- **Loot Amount = Lootable × 20%**

**Example:**
- Defender has: 500k Metal, 300k He3, 200k Gold (total)
- Warehouse: 100k Metal, 50k He3, 50k Gold
- Lootable: 400k Metal, 250k He3, 150k Gold
- Attacker wins: 80k Metal, 50k He3, 30k Gold

**Winner Takes All:** If attacker wins, they get loot. If defender wins, no loot.

### 3.5 Radar Building Mechanics

From GO2 wiki:

**Radar Levels (1-9):**
| Level | Detection Time | Info Revealed |
|-------|----------------|---------------|
| 1-2   | 30-60 min      | Arrival time only |
| 3-4   | 90-120 min     | Arrival time + Origin coords |
| 5-6   | 150-180 min    | Above + Fleet size |
| 7-8   | 210-240 min    | Above + Fleet types |
| 9     | 270 min (4.5h) | Above + Commander |

**Warning System:**
1. Attacker launches fleet at time T
2. Defender's Radar checks detection time D
3. If (T + D) < arrival_time, warning shown
4. Warning displays info based on Radar level

**Example:**
- Attacker launches at 10:00 AM
- Travel time: 30 minutes (arrival 10:30 AM)
- Defender Radar Level 5 (detection 150 min)
- Warning time: 10:00 + 150 min = 12:30 PM > 10:30 AM arrival
- **Result:** No warning (arrival before detection)

**For effective warning:**
- Radar detection time < fleet travel time
- Higher level = earlier warning
- Level 9 Radar (4.5h detection) warns about long-distance attacks

### 3.6 Defense Fleet Mechanics

**From GO2 wiki:**
- Players can assign fleet(s) to "Defend" status
- Defense fleets automatically engage attackers
- Multiple defense fleets fight sequentially (not simultaneously)
- If no defense fleet, attacker fights empty planet (instant win)

**In Cryptomines Online:**
- Player assigns 1 primary defense fleet (simple version)
- Defense fleet must be idle (not traveling, not in other combat)
- Defense fleet uses same combat engine (Module 5)
- Defense fleet casualties apply after combat
- If defense fleet destroyed, can't defend until repaired/rebuilt

### 3.7 Truce Card Protection

**From GO2 wiki:**
- **Truce Card:** 12-hour protection (cannot be attacked)
- **Advanced Truce Card:** 72-hour protection (3 days)
- **Restrictions:**
  - Cannot be used during incoming attack (must use before)
  - Player cannot attack others while protected
  - Protection expires naturally after duration

**In Cryptomines Online:**
- Truce Card (12h) - From inventory (Module 9)
- Advanced Truce Card (72h) - From inventory
- Check `active_protections` table before allowing attacks
- Validation: Attacker not protected, Defender not protected

---

## 4. PVP COMBAT FLOW DESIGN

### 4.1 Attack Initiation

**Player Action:** Select neighbor → Click "Attack" → Choose fleet

**Backend Validation:**
1. Check attacker has fleet with ships
2. Check attacker SP ≥ travel cost
3. Check attacker not protected by truce card
4. Check defender not protected by truce card
5. Check fleet not already in combat/traveling
6. Calculate travel time and arrival time
7. Deduct SP from attacker
8. Set fleet status = 'moving'
9. Create `fleet_movements` row
10. Return success + arrival time

**Pseudocode:**
```go
func InitiatePvPAttack(attackerID, targetID, fleetID string) (*FleetMovement, error) {
    // 1. Validate fleet
    var fleet Fleet
    err := db.Get(&fleet, `SELECT * FROM fleets WHERE id = $1 AND player_id = $2`, fleetID, attackerID)
    if err != nil || fleet.Status != "idle" {
        return nil, errors.New("fleet not available")
    }

    // 2. Check ship count
    var shipCount int
    db.Get(&shipCount, `SELECT COUNT(*) FROM fleet_ships WHERE fleet_id = $1`, fleetID)
    if shipCount == 0 {
        return nil, errors.New("fleet has no ships")
    }

    // 3. Check protections
    attackerProtected := HasActiveProtection(attackerID)
    defenderProtected := HasActiveProtection(targetID)
    if attackerProtected {
        return nil, errors.New("cannot attack while protected")
    }
    if defenderProtected {
        return nil, errors.New("target is protected by truce card")
    }

    // 4. Calculate distance and SP cost
    distance := CalculateDistance(attackerID, targetID)
    spCost := int(math.Ceil(float64(distance) / 10.0))

    // 5. Check SP
    var sp int
    db.Get(&sp, `SELECT space_points FROM resources WHERE player_id = $1`, attackerID)
    if sp < spCost {
        return nil, errors.New("insufficient Space Points")
    }

    // 6. Calculate travel time
    travelMinutes := int(math.Ceil(float64(distance) / 10.0))
    arrivalTime := time.Now().Add(time.Duration(travelMinutes) * time.Minute)

    // 7. Begin transaction
    tx, _ := db.Begin()
    defer tx.Rollback()

    // 8. Deduct SP
    tx.Exec(`UPDATE resources SET space_points = space_points - $1 WHERE player_id = $2`, spCost, attackerID)

    // 9. Update fleet status
    tx.Exec(`UPDATE fleets SET status = 'moving' WHERE id = $1`, fleetID)

    // 10. Create fleet movement
    var movementID string
    tx.QueryRow(`
        INSERT INTO fleet_movements (fleet_id, player_id, origin_player_id, target_player_id,
                                      movement_type, arrival_time, sp_consumed, status)
        VALUES ($1, $2, $3, $4, 'attack', $5, $6, 'traveling')
        RETURNING id
    `, fleetID, attackerID, attackerID, targetID, arrivalTime, spCost).Scan(&movementID)

    tx.Commit()

    return &FleetMovement{
        ID: movementID,
        ArrivalTime: arrivalTime,
        SPConsumed: spCost,
    }, nil
}
```

**Time Estimate:** 4-5 hours

---

### 4.2 Fleet Arrival & Combat Resolution

**Worker Process:** Runs every 30 seconds, checks for arrived fleets.

```go
func ProcessArrivedFleets() {
    rows, _ := db.Query(`
        SELECT id, fleet_id, origin_player_id, target_player_id
        FROM fleet_movements
        WHERE status = 'traveling' AND arrival_time <= now()
    `)

    for rows.Next() {
        var movement FleetMovement
        rows.Scan(&movement.ID, &movement.FleetID, &movement.OriginPlayerID, &movement.TargetPlayerID)

        // Execute combat
        ResolvePvPCombat(movement.FleetID, movement.TargetPlayerID, movement.ID)
    }
}

func ResolvePvPCombat(attackFleetID, defenderPlayerID, movementID string) {
    // 1. Get defender's defense fleet
    defenseFleetID := GetDefenseFleet(defenderPlayerID)

    // 2. If no defense fleet, attacker wins instantly
    if defenseFleetID == "" {
        // Instant victory, loot resources
        loot := CalculateLoot(defenderPlayerID)
        GrantLoot(attackFleetID.PlayerID, loot)

        // Create combat report (0 rounds, instant win)
        SaveCombatReport(attackFleetID, "", "attacker_win", 0, loot)

        // Start return journey
        InitiateReturnJourney(attackFleetID, movementID)
        return
    }

    // 3. Execute combat (reuse Module 5 engine)
    result := services.ResolveCombat(attackFleetID, defenseFleetID, "pvp")

    // 4. Apply casualties
    UpdateFleetShips(attackFleetID, result.AttackerFinalState)
    UpdateFleetShips(defenseFleetID, result.DefenderFinalState)

    // 5. Handle loot
    if result.Result == "attacker_win" {
        loot := CalculateLoot(defenderPlayerID)
        GrantLoot(attackFleetID.PlayerID, loot)
        result.Loot = loot
    }

    // 6. Save combat report
    SaveCombatReport(result)

    // 7. Send notifications to both players
    NotifyPlayer(attackFleetID.PlayerID, "Combat report available")
    NotifyPlayer(defenderPlayerID, "Your planet was attacked!")

    // 8. Update fleet movement status
    db.Exec(`UPDATE fleet_movements SET status = 'arrived' WHERE id = $1`, movementID)

    // 9. Start return journey
    InitiateReturnJourney(attackFleetID, movementID)
}
```

**Time Estimate:** 6-8 hours

---

### 4.3 Resource Loot Calculation

```go
func CalculateLoot(defenderPlayerID string) map[string]int64 {
    // Get defender resources
    var total struct {
        Metal BIGINT
        He3   BIGINT
        Gold  BIGINT
        WarehouseMetal BIGINT
        WarehouseHe3   BIGINT
        WarehouseGold  BIGINT
    }

    db.QueryRow(`
        SELECT metal, he3, gold, warehouse_metal, warehouse_he3, warehouse_gold
        FROM resources
        WHERE player_id = $1
    `, defenderPlayerID).Scan(
        &total.Metal, &total.He3, &total.Gold,
        &total.WarehouseMetal, &total.WarehouseHe3, &total.WarehouseGold,
    )

    // Calculate lootable (exclude warehouse)
    lootableMetal := total.Metal - total.WarehouseMetal
    lootableHe3 := total.He3 - total.WarehouseHe3
    lootableGold := total.Gold - total.WarehouseGold

    // Clamp to zero (can't loot negative)
    if lootableMetal < 0 { lootableMetal = 0 }
    if lootableHe3 < 0 { lootableHe3 = 0 }
    if lootableGold < 0 { lootableGold = 0 }

    // Take 20%
    loot := map[string]int64{
        "metal": int64(float64(lootableMetal) * 0.20),
        "he3":   int64(float64(lootableHe3) * 0.20),
        "gold":  int64(float64(lootableGold) * 0.20),
    }

    return loot
}

func GrantLoot(attackerPlayerID string, loot map[string]int64) {
    db.Exec(`
        UPDATE resources
        SET metal = metal + $1, he3 = he3 + $2, gold = gold + $3
        WHERE player_id = $4
    `, loot["metal"], loot["he3"], loot["gold"], attackerPlayerID)

    // Deduct from defender
    db.Exec(`
        UPDATE resources
        SET metal = GREATEST(0, metal - $1),
            he3 = GREATEST(0, he3 - $2),
            gold = GREATEST(0, gold - $3)
        WHERE player_id = $4
    `, loot["metal"], loot["he3"], loot["gold"], defenderPlayerID)
}
```

**Time Estimate:** 2 hours

---

### 4.4 Return Journey

```go
func InitiateReturnJourney(fleetID, originalMovementID string) {
    // Get original movement
    var original FleetMovement
    db.QueryRow(`SELECT origin_player_id, target_player_id, sp_consumed FROM fleet_movements WHERE id = $1`,
        originalMovementID).Scan(&original.OriginPlayerID, &original.TargetPlayerID, &original.SPConsumed)

    // Calculate return time (half of outbound)
    distance := CalculateDistance(original.TargetPlayerID, original.OriginPlayerID)
    returnMinutes := int(math.Ceil(float64(distance) / 20.0)) // Double speed
    arrivalTime := time.Now().Add(time.Duration(returnMinutes) * time.Minute)

    // Create return movement
    db.Exec(`
        INSERT INTO fleet_movements (fleet_id, player_id, origin_player_id, target_player_id,
                                      movement_type, arrival_time, sp_consumed, status)
        VALUES ($1, $2, $3, $4, 'return', $5, 0, 'traveling')
    `, fleetID, original.OriginPlayerID, original.TargetPlayerID, original.OriginPlayerID, arrivalTime)
}

func ProcessReturningFleets() {
    rows, _ := db.Query(`
        SELECT id, fleet_id
        FROM fleet_movements
        WHERE status = 'traveling' AND movement_type = 'return' AND arrival_time <= now()
    `)

    for rows.Next() {
        var movementID, fleetID string
        rows.Scan(&movementID, &fleetID)

        // Fleet arrives home
        db.Exec(`UPDATE fleets SET status = 'idle' WHERE id = $1`, fleetID)
        db.Exec(`UPDATE fleet_movements SET status = 'arrived' WHERE id = $1`, movementID)

        // Notify player
        var playerID string
        db.QueryRow(`SELECT player_id FROM fleets WHERE id = $1`, fleetID).Scan(&playerID)
        NotifyPlayer(playerID, "Your fleet has returned home")
    }
}
```

**Time Estimate:** 2 hours

---

### 4.5 Radar Warning System

```go
func CheckRadarWarnings() {
    // Run every minute

    // Get all active incoming attacks
    rows, _ := db.Query(`
        SELECT fm.id, fm.target_player_id, fm.origin_player_id, fm.arrival_time,
               fm.fleet_id, fm.departure_time
        FROM fleet_movements fm
        WHERE fm.status = 'traveling' AND fm.movement_type = 'attack'
    `)

    for rows.Next() {
        var movement struct {
            ID            string
            TargetPlayerID string
            OriginPlayerID string
            ArrivalTime    time.Time
            FleetID        string
            DepartureTime  time.Time
        }

        rows.Scan(&movement.ID, &movement.TargetPlayerID, &movement.OriginPlayerID,
                  &movement.ArrivalTime, &movement.FleetID, &movement.DepartureTime)

        // Check if warning already sent
        var warningSent bool
        db.QueryRow(`SELECT EXISTS(SELECT 1 FROM radar_warnings WHERE movement_id = $1)`,
            movement.ID).Scan(&warningSent)

        if warningSent {
            continue
        }

        // Get defender's Radar level
        radarLevel := GetRadarLevel(movement.TargetPlayerID)

        if radarLevel == 0 {
            continue // No Radar, no warning
        }

        // Get detection time for Radar level
        var detectionMinutes int
        db.QueryRow(`SELECT detection_time_minutes FROM radar_levels WHERE level = $1`,
            radarLevel).Scan(&detectionMinutes)

        // Check if current time is within detection window
        detectionTime := movement.DepartureTime.Add(time.Duration(detectionMinutes) * time.Minute)

        if time.Now().After(detectionTime) || time.Now().After(movement.ArrivalTime) {
            continue // Too late to warn
        }

        // Generate warning info based on Radar level
        warningInfo := GenerateRadarWarning(movement, radarLevel)

        // Save warning
        db.Exec(`
            INSERT INTO radar_warnings (movement_id, player_id, warning_json, created_at)
            VALUES ($1, $2, $3, now())
        `, movement.ID, movement.TargetPlayerID, warningInfo)

        // Notify player
        NotifyPlayer(movement.TargetPlayerID, "Incoming attack detected!")
    }
}

func GenerateRadarWarning(movement Movement, radarLevel int) string {
    info := map[string]interface{}{
        "arrival_time": movement.ArrivalTime,
    }

    if radarLevel >= 3 {
        // Show origin coords
        coords := GetPlayerCoordinates(movement.OriginPlayerID)
        info["origin_coords"] = coords
    }

    if radarLevel >= 5 {
        // Show fleet size
        var shipCount int
        db.QueryRow(`SELECT COUNT(*) FROM fleet_ships WHERE fleet_id = $1`,
            movement.FleetID).Scan(&shipCount)
        info["fleet_size"] = shipCount
    }

    if radarLevel >= 7 {
        // Show fleet types
        fleetTypes := GetFleetComposition(movement.FleetID)
        info["fleet_types"] = fleetTypes
    }

    if radarLevel >= 9 {
        // Show commander
        var commanderName string
        db.QueryRow(`
            SELECT c.name FROM fleets f
            JOIN commanders c ON f.commander_id = c.id
            WHERE f.id = $1
        `, movement.FleetID).Scan(&commanderName)
        info["commander"] = commanderName
    }

    jsonData, _ := json.Marshal(info)
    return string(jsonData)
}
```

**NEW TABLE: `radar_warnings`**
```sql
CREATE TABLE radar_warnings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movement_id UUID NOT NULL REFERENCES fleet_movements(id) ON DELETE CASCADE,
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    warning_json JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(movement_id, player_id)
);

CREATE INDEX idx_radar_warnings_player ON radar_warnings(player_id);
```

**Time Estimate:** 4-5 hours

---

## 5. NEIGHBOR DISCOVERY SYSTEM

### 5.1 Galaxy Coordinate System

**Coordinate Generation (on player creation):**
```go
func AssignGalaxyCoordinates(playerID string) {
    // Assign random sector (1-10)
    sector := rand.Intn(10) + 1

    // Assign random coords within sector (0-999)
    x := rand.Intn(1000)
    y := rand.Intn(1000)

    db.Exec(`
        UPDATE players
        SET galaxy_sector = $1, galaxy_x = $2, galaxy_y = $3
        WHERE id = $4
    `, sector, x, y, playerID)
}
```

### 5.2 Neighbor Discovery

**Find neighbors within distance threshold:**
```go
func GetNeighbors(playerID string, maxDistance int) ([]Player, error) {
    // Get player coords
    var player struct {
        Sector int
        X      int
        Y      int
    }

    db.QueryRow(`SELECT galaxy_sector, galaxy_x, galaxy_y FROM players WHERE id = $1`,
        playerID).Scan(&player.Sector, &player.X, &player.Y)

    // Find neighbors in same sector within distance
    rows, err := db.Query(`
        SELECT id, username, level, galaxy_x, galaxy_y,
               SQRT(POWER($2 - galaxy_x, 2) + POWER($3 - galaxy_y, 2)) AS distance
        FROM players
        WHERE id != $1
          AND galaxy_sector = $4
          AND SQRT(POWER($2 - galaxy_x, 2) + POWER($3 - galaxy_y, 2)) <= $5
        ORDER BY distance ASC
        LIMIT 50
    `, playerID, player.X, player.Y, player.Sector, maxDistance)

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    neighbors := []Player{}
    for rows.Next() {
        var n Player
        rows.Scan(&n.ID, &n.Username, &n.Level, &n.X, &n.Y, &n.Distance)
        neighbors = append(neighbors, n)
    }

    return neighbors, nil
}
```

**Default Distance Threshold:** 500 units (half of sector size).

**Time Estimate:** 3-4 hours

---

## 6. DEFENSE FLEET SYSTEM

### 6.1 Assign Defense Fleet

```go
func AssignDefenseFleet(playerID, fleetID string) error {
    // Validate fleet
    var fleet Fleet
    err := db.QueryRow(`SELECT id, status FROM fleets WHERE id = $1 AND player_id = $2`,
        fleetID, playerID).Scan(&fleet.ID, &fleet.Status)

    if err != nil {
        return errors.New("fleet not found")
    }

    if fleet.Status != "idle" {
        return errors.New("fleet must be idle to assign as defense")
    }

    // Check ship count
    var shipCount int
    db.QueryRow(`SELECT COUNT(*) FROM fleet_ships WHERE fleet_id = $1`, fleetID).Scan(&shipCount)

    if shipCount == 0 {
        return errors.New("fleet has no ships")
    }

    // Assign defense
    tx, _ := db.Begin()
    defer tx.Rollback()

    // Remove existing defense assignment
    tx.Exec(`DELETE FROM defense_assignments WHERE player_id = $1`, playerID)

    // Assign new defense fleet
    tx.Exec(`INSERT INTO defense_assignments (player_id, fleet_id) VALUES ($1, $2)`,
        playerID, fleetID)

    // Update fleet status
    tx.Exec(`UPDATE fleets SET status = 'defending' WHERE id = $1`, fleetID)

    tx.Commit()

    return nil
}

func GetDefenseFleet(playerID string) string {
    var fleetID string
    db.QueryRow(`SELECT fleet_id FROM defense_assignments WHERE player_id = $1`,
        playerID).Scan(&fleetID)

    return fleetID // Empty string if no defense fleet
}
```

**Time Estimate:** 2 hours

---

## 7. PROTECTION SYSTEM (TRUCE CARDS)

### 7.1 Activate Truce Card

```go
func ActivateTruceCard(playerID string, duration time.Duration) error {
    // Check no incoming attacks
    var incomingCount int
    db.QueryRow(`
        SELECT COUNT(*)
        FROM fleet_movements
        WHERE target_player_id = $1 AND status = 'traveling' AND movement_type = 'attack'
    `, playerID).Scan(&incomingCount)

    if incomingCount > 0 {
        return errors.New("cannot use truce card during incoming attack")
    }

    // Activate protection
    expiresAt := time.Now().Add(duration)

    _, err := db.Exec(`
        INSERT INTO active_protections (player_id, protection_type, expires_at)
        VALUES ($1, 'truce_card', $2)
        ON CONFLICT (player_id, protection_type)
        DO UPDATE SET expires_at = $2
    `, playerID, expiresAt)

    return err
}

func HasActiveProtection(playerID string) bool {
    var count int
    db.QueryRow(`
        SELECT COUNT(*)
        FROM active_protections
        WHERE player_id = $1 AND expires_at > now()
    `, playerID).Scan(&count)

    return count > 0
}
```

**Truce Card Durations:**
- Truce Card: 12 hours
- Advanced Truce Card: 72 hours

**Time Estimate:** 2 hours

---

## 8. BACKEND IMPLEMENTATION

### 8.1 Handler Layer (`backend/internal/handlers/pvp.go`)

**Endpoints:**

1. **GET /api/pvp/neighbors** - List nearby players
2. **POST /api/pvp/attack** - Initiate attack
3. **GET /api/pvp/incoming** - List incoming attacks (Radar warnings)
4. **POST /api/pvp/defense/assign** - Assign defense fleet
5. **DELETE /api/pvp/defense** - Unassign defense fleet
6. **POST /api/pvp/protection/activate** - Use truce card
7. **GET /api/pvp/protection** - Check active protections

**Time Estimate:** 8-10 hours

---

### 8.2 Worker Layer (`backend/internal/workers/pvp_worker.go`)

**Workers:**
1. **ProcessArrivedFleets()** - Check fleet arrivals every 30s
2. **ProcessReturningFleets()** - Check return arrivals every 30s
3. **CheckRadarWarnings()** - Generate warnings every 1 min
4. **RegenerateSP()** - Regen SP every 6 minutes (1 SP)
5. **ExpireProtections()** - Remove expired truce cards every 5 min

**Time Estimate:** 4-5 hours

---

## 9. FRONTEND IMPLEMENTATION

### 9.1 API Client (`frontend/src/api/pvp.ts`)

```typescript
export interface Neighbor {
  id: string;
  username: string;
  level: number;
  distance: number;
  coords: { x: number; y: number };
  is_protected: boolean;
  has_defense: boolean;
}

export interface IncomingAttack {
  id: string;
  origin_coords?: { x: number; y: number };
  arrival_time: string;
  fleet_size?: number;
  fleet_types?: string[];
  commander?: string;
}

export const pvpApi = {
  getNeighbors: () =>
    apiClient.get<Neighbor[]>('/pvp/neighbors'),

  attack: (targetPlayerId: string, fleetId: string) =>
    apiClient.post('/pvp/attack', { target_player_id: targetPlayerId, fleet_id: fleetId }),

  getIncomingAttacks: () =>
    apiClient.get<IncomingAttack[]>('/pvp/incoming'),

  assignDefenseFleet: (fleetId: string) =>
    apiClient.post('/pvp/defense/assign', { fleet_id: fleetId }),

  unassignDefenseFleet: () =>
    apiClient.delete('/pvp/defense'),

  activateTruceCard: (duration: number) =>
    apiClient.post('/pvp/protection/activate', { duration_hours: duration }),

  getProtection: () =>
    apiClient.get('/pvp/protection'),
};
```

**Time Estimate:** 1 hour

---

### 9.2 Galaxy/Neighbors Panel (`frontend/src/components/panels/GalaxyPanel.tsx`)

**Features:**
- List of neighbors (sorted by distance)
- Player info (username, level, coords, distance)
- "Attack" button per neighbor
- Fleet selection modal
- SP cost display
- Travel time estimate
- Protection shield indicator

**Time Estimate:** 4-5 hours

---

### 9.3 Incoming Attacks Panel (`frontend/src/components/panels/IncomingAttacksPanel.tsx`)

**Features:**
- List of incoming attacks
- Countdown timers to arrival
- Info based on Radar level
- "View Details" button
- Defense fleet status
- "Assign Defense" button

**Time Estimate:** 3-4 hours

---

### 9.4 Defense Fleet Panel

Update `FleetPanel.tsx` to add:
- "Assign as Defense" button
- "Unassign Defense" button
- Defense status indicator

**Time Estimate:** 2 hours

---

## 10. QA CHECKLIST

### 10.1 Attack Flow Tests
- [ ] POST /api/pvp/attack initiates attack successfully
- [ ] Fleet status changes to 'moving'
- [ ] SP deducted correctly based on distance
- [ ] Fleet movement row created with correct arrival time
- [ ] Cannot attack without sufficient SP
- [ ] Cannot attack with fleet already in motion
- [ ] Cannot attack while protected by truce card
- [ ] Cannot attack protected target
- [ ] Cannot attack with empty fleet

### 10.2 Combat Resolution Tests
- [ ] Arrived fleets trigger combat automatically
- [ ] Combat uses Module 5 engine correctly
- [ ] Defense fleet participates in combat
- [ ] No defense fleet = instant attacker win
- [ ] Casualties applied to both fleets
- [ ] Loot calculated correctly (20% of lootable)
- [ ] Loot granted to attacker on win
- [ ] Combat report saved with 'pvp' type
- [ ] Both players receive combat report notification

### 10.3 Return Journey Tests
- [ ] Return journey initiated after combat
- [ ] Return time is half of outbound time
- [ ] Fleet arrives home and status = 'idle'
- [ ] No SP cost for return
- [ ] Player notified on fleet return

### 10.4 Neighbor Discovery Tests
- [ ] GET /api/pvp/neighbors returns nearby players
- [ ] Neighbors sorted by distance
- [ ] Distance calculated correctly (Euclidean)
- [ ] Only players in same sector shown
- [ ] Max 50 neighbors returned
- [ ] Player not included in own neighbor list

### 10.5 Radar Warning Tests
- [ ] Radar Level 1-2 shows arrival time only
- [ ] Radar Level 3-4 shows arrival + origin coords
- [ ] Radar Level 5-6 shows fleet size
- [ ] Radar Level 7-8 shows fleet types
- [ ] Radar Level 9 shows commander
- [ ] Warning generated when detection time < arrival
- [ ] No warning if Radar level 0
- [ ] Warning sent only once per attack
- [ ] Player notified of incoming attack

### 10.6 Defense Fleet Tests
- [ ] POST /api/pvp/defense/assign sets defense fleet
- [ ] Fleet status = 'defending'
- [ ] Only idle fleets can be assigned
- [ ] Only 1 defense fleet per player
- [ ] Reassigning replaces old defense fleet
- [ ] Defense fleet engages attackers automatically
- [ ] DELETE /api/pvp/defense unassigns fleet
- [ ] Fleet status returns to 'idle' after unassign

### 10.7 Truce Card Tests
- [ ] POST /api/pvp/protection/activate activates truce
- [ ] 12h truce card = 12 hours protection
- [ ] 72h truce card = 72 hours protection
- [ ] Cannot activate during incoming attack
- [ ] Protected player cannot be attacked
- [ ] Protected player cannot attack others
- [ ] Protection expires after duration
- [ ] Expired protections removed automatically

### 10.8 SP System Tests
- [ ] SP regenerates 1 per 6 minutes
- [ ] Max SP = 100 + (10 × level)
- [ ] SP does not exceed max
- [ ] Attack deducts SP correctly
- [ ] Insufficient SP prevents attack

### 10.9 Resource Looting Tests
- [ ] Lootable = Total - Warehouse
- [ ] Loot = 20% of lootable
- [ ] Loot granted to attacker on win
- [ ] Loot deducted from defender
- [ ] No loot if defender wins
- [ ] Negative lootable = 0 loot

### 10.10 Integration Tests
- [ ] Module 5 combat engine reused correctly
- [ ] Commander bonuses apply in PvP
- [ ] Tech bonuses apply in PvP
- [ ] Quest progress triggered (win_pvp_battle)
- [ ] Warehouse resources excluded from loot
- [ ] Fleet casualties permanent (ships lost)

**Total QA Checks:** 52

---

## 11. RISK ASSESSMENT

### 11.1 Technical Risks

**RISK 1: Worker Performance (Fleet Arrivals)**
- **Severity:** MEDIUM
- **Impact:** 100+ players attacking = 100+ fleet arrivals to process
- **Mitigation:** Optimize worker query (index on arrival_time), process in batches
- **Target:** <1s for 100 simultaneous arrivals

**RISK 2: Coordinate System Collisions**
- **Severity:** LOW
- **Impact:** Multiple players at same coords
- **Mitigation:** Random spawn in 1000×1000 grid = 1M positions, very low collision rate
- **Alternative:** Add unique constraint, regenerate on conflict

**RISK 3: Radar Warning Spam**
- **Severity:** LOW
- **Impact:** Player attacked 10 times = 10 warnings
- **Mitigation:** One warning per attack (unique constraint on radar_warnings)
- **Enhancement:** Limit warnings to latest 20

**RISK 4: SP Regeneration Abuse**
- **Severity:** LOW
- **Impact:** Player manipulates SP regen timing
- **Mitigation:** Server-side regen only, no client control
- **Formula:** SP = min(max_sp, current_sp + elapsed_minutes/6)

---

### 11.2 Balance Risks

**RISK 5: Loot Too High/Low**
- **Severity:** MEDIUM
- **Impact:** 20% loot feels too generous or too stingy
- **Mitigation:** Follows GO2 standard, playtesting will validate
- **Adjustment:** Can change 20% to 10-30% if needed

**RISK 6: Travel Time Too Fast/Slow**
- **Severity:** LOW
- **Impact:** 1 min per 10 units may feel wrong
- **Mitigation:** Configurable formula, easy to adjust
- **Alternative:** 2 min per 10 units, or tech-based speed bonuses

**RISK 7: Defense Fleet Balance**
- **Severity:** MEDIUM
- **Impact:** 1 defense fleet may be too weak vs multiple attacks
- **Mitigation:** This matches GO2 (sequential defense fleets)
- **Enhancement:** Module 7 adds defense buildings for extra power

---

### 11.3 UX Risks

**RISK 8: No Galaxy Map Visualization**
- **Severity:** LOW
- **Impact:** Text-only neighbor list feels basic
- **Mitigation:** Acceptable for Phase A, visual map in Module 11 (Polish)
- **Workaround:** Show distance + coords for spatial awareness

**RISK 9: Incoming Attack Confusion**
- **Severity:** MEDIUM
- **Impact:** Players miss warnings, lose resources
- **Mitigation:** Clear UI notification, countdown timers, defense assignment prompts
- **Enhancement:** Browser notifications for incoming attacks

---

## 12. IMPLEMENTATION TIMELINE

### 12.1 Phase Breakdown

**Phase 1: Database & Coordinates (Day 1)**
- Migration: fleet_movements, defense_assignments, active_protections, radar_warnings
- ALTER players (galaxy coords)
- ALTER resources (SP field)
- Seed player coordinates
- **Time:** 4-5 hours

**Phase 2: Neighbor Discovery & Attack (Days 1-2)**
- Neighbor discovery logic
- Attack initiation endpoint
- SP validation
- Protection checks
- Distance calculation
- **Time:** 6-8 hours

**Phase 3: Combat Resolution (Day 2)**
- Fleet arrival worker
- Combat integration (Module 5)
- Loot calculation
- Return journey
- **Time:** 6-8 hours

**Phase 4: Defense & Radar (Day 3)**
- Defense fleet assignment
- Radar warning worker
- Warning generation by level
- Radar info display
- **Time:** 6-8 hours

**Phase 5: Backend Handlers (Day 4)**
- PvP endpoints (7 total)
- Error handling
- Validation logic
- **Time:** 4-5 hours

**Phase 6: Frontend UI (Days 4-5)**
- GalaxyPanel (neighbors list)
- IncomingAttacksPanel
- Defense fleet UI updates
- API client + hooks
- **Time:** 8-10 hours

**Phase 7: QA & Polish (Day 6)**
- Run full QA checklist (52 checks)
- Balance testing (loot %, travel time)
- Fix bugs
- **Time:** 6-8 hours

---

### 12.2 Critical Path

```
Day 1: Database → Coordinates → Neighbor Discovery → Attack Initiation
       ↓
Day 2: Fleet Arrival Worker → Combat Resolution → Return Journey
       ↓
Day 3: Defense Fleet → Radar Warnings → Worker Integration
       ↓
Day 4: Backend Handlers → Frontend Galaxy Panel
       ↓
Day 5: Frontend Incoming Attacks → Defense UI
       ↓
Day 6: QA Testing → Balancing → Bug Fixes
```

**Total Estimated Time:** 25-35 hours (4-6 working days)

---

## 13. SUCCESS CRITERIA

Module 6 is considered **COMPLETE** when:

✅ **Attack Flow:**
- Players can select neighbors and initiate attacks
- Fleet travel system works (outbound + return)
- SP deducted correctly
- Protection system enforced

✅ **Combat:**
- Module 5 engine reused for PvP
- Defense fleets engage attackers
- Casualties applied correctly
- Loot distributed (20% of lootable resources)

✅ **Radar System:**
- Warnings generated based on Radar level
- Info revealed progressively (arrival → coords → size → types → commander)
- Players notified of incoming attacks

✅ **Defense:**
- Players can assign defense fleet
- Defense fleet automatically defends
- Fleet status managed correctly

✅ **Protection:**
- Truce cards prevent attacks (12h/72h)
- Cannot use during incoming attack
- Protected players cannot attack

✅ **Frontend:**
- Galaxy/neighbors list displays
- Incoming attacks panel shows warnings
- Defense fleet assignment works
- Combat reports accessible

✅ **QA:**
- All 52 QA checks pass
- No critical bugs
- Balance feels fair (loot, travel time, SP costs)

---

## 14. NEXT STEPS AFTER MODULE 6

After Module 6 completion, PvP combat is fully functional. Next modules:

1. **Module 7: Space Defense Buildings** - Add 5 defensive structures to boost defense
2. **Module 8: Recycling Plant** - Scrap unwanted ships for resources
3. **Module 9: Inventory System** - Full item management (truce cards, resource packs)
4. **Module 11: Production Polish** - Improve PvP UX (galaxy map, notifications, animations)

**Recommended Next Module:** Module 7 (Space Defense Buildings) to strengthen defense mechanics.

---

## APPENDIX A: FORMULAS REFERENCE

### A.1 Distance Calculation
```
Distance = √((x2 - x1)² + (y2 - y1)²)
```

### A.2 SP Cost
```
SP Cost = CEIL(Distance / 10)
```

### A.3 Travel Time
```
Outbound Time (minutes) = CEIL(Distance / 10)
Return Time (minutes) = CEIL(Outbound Time / 2)
```

### A.4 SP Regeneration
```
SP Regen = 1 SP per 6 minutes
Max SP = 100 + (Player Level × 10)
```

### A.5 Resource Loot
```
Lootable = Total Resources - Warehouse Resources
Loot = Lootable × 20%
```

### A.6 Radar Detection Time
| Level | Detection Minutes |
|-------|-------------------|
| 1     | 30                |
| 2     | 60                |
| 3     | 90                |
| 4     | 120               |
| 5     | 150               |
| 6     | 180               |
| 7     | 210               |
| 8     | 240               |
| 9     | 270               |

---

## APPENDIX B: References

**Galaxy Online 2 Wiki Sources:**
- [Attacking Neighbors (PvP)](https://galaxyonlineii.fandom.com/wiki/Attacking_Neighbors_(PvP)) - PvP mechanics
- [Radar](https://galaxyonlineii.fandom.com/wiki/Radar) - Radar building levels
- [Battle Items](https://galaxyonlineii.fandom.com/wiki/Battle_Items) - Truce cards

**Project Documents:**
- `/docs/planning/final-scope.md` - Module 6 scope
- `/docs/planning/module-5-combat-system-plan.md` - Combat engine (reused)
- `/supabase/migrations/20260206020000_phase2_ships.sql` - Existing PvP support

---

**END OF MODULE 6 IMPLEMENTATION PLAN**
