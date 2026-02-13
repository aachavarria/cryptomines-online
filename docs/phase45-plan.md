# Phase 4.5: PvP Travel & Protection Systems

## Overview
Refactor PvP from synchronous (instant combat) to asynchronous (fleet travel → combat on arrival → return home). Add SP system, Truce Card activation, and Radar incoming attacks detection.

## Architecture Change

### Current Flow (Instant)
```
POST /api/pvp/attack → validate → combat resolves → response with result
```

### New Flow (Asynchronous)
```
POST /api/pvp/attack → validate → fleets set to "traveling" → returns travel_time
Worker: pending_attacks where arrival_at < now() → resolve combat → fleets "returning"
Worker: returning fleets where arrival_at < now() → fleets "stationed"
GET /api/radar/incoming → shows incoming attacks based on radar level
```

---

## 1. Database Migration

### New table: `pending_attacks`
```sql
CREATE TABLE pending_attacks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attacker_id UUID NOT NULL REFERENCES players(id),
    defender_id UUID NOT NULL REFERENCES players(id),
    defender_planet_id UUID NOT NULL REFERENCES planets(id),
    fleet_ids TEXT[] NOT NULL,        -- array of fleet UUIDs sent
    status TEXT NOT NULL DEFAULT 'traveling'
        CHECK (status IN ('traveling', 'combat', 'resolved', 'cancelled')),
    depart_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    arrival_at TIMESTAMPTZ NOT NULL,  -- calculated travel time
    return_at TIMESTAMPTZ,            -- set after combat (50% of travel time)
    sp_cost INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_pending_attacks_arrival ON pending_attacks (arrival_at) WHERE status = 'traveling';
CREATE INDEX idx_pending_attacks_defender ON pending_attacks (defender_id) WHERE status = 'traveling';
```

### Add to `players` table
```sql
ALTER TABLE players ADD COLUMN space_points INTEGER NOT NULL DEFAULT 20;
ALTER TABLE players ADD COLUMN max_space_points INTEGER NOT NULL DEFAULT 20;
ALTER TABLE players ADD COLUMN sp_last_reset TIMESTAMPTZ NOT NULL DEFAULT now();
```

### planets.protection_until already exists — no change needed

---

## 2. Travel Time Formula

**Fleet speed** = MIN(TotalMovement) across all ship designs in all sent fleets (slowest ship dictates fleet speed). Minimum 1.

**Distance** = Euclidean distance between attacker planet position and defender planet position:
```
distance = sqrt((x2-x1)² + (y2-y1)²)
```

**Travel time (seconds)**:
```
travel_seconds = (distance / fleet_speed) * BASE_TRAVEL_FACTOR
```
Where `BASE_TRAVEL_FACTOR = 300` (5 minutes per distance unit at speed 1).

**Minimum travel time**: 30 seconds
**Maximum travel time**: 3600 seconds (1 hour)

**Return time** = 50% of outbound travel time (per GO2 wiki).

---

## 3. Backend Changes

### 3.1 Refactor AttackPlanet handler (`pvp.go`)

**Before:** Validates → loads fleets → runs combat → returns result
**After:** Validates → calculates travel time → deducts SP → creates pending_attack → sets fleets traveling → returns travel info

```go
type attackPlanetResponse struct {
    PendingAttackID string  `json:"pending_attack_id"`
    TravelSeconds   int     `json:"travel_seconds"`
    ArrivalAt       string  `json:"arrival_at"`
    FleetsDispatched int    `json:"fleets_dispatched"`
}
```

New validations:
- Check attacker has SP >= 1
- Check defender planet `protection_until` is null or past
- Check attacker is not under truce protection (bidirectional)

### 3.2 New PvP Worker (`pvp_worker.go`)

Runs every 5 seconds. Two jobs:

**Job 1: Resolve arrived attacks**
```
SELECT * FROM pending_attacks WHERE status = 'traveling' AND arrival_at <= now()
```
For each:
1. Load attacker fleets + defender fleets + defense buildings
2. Run combat engine
3. Apply casualties, loot, defense reset
4. Create combat_report
5. Set pending_attack status = 'resolved'
6. Set attacker fleets status = 'returning', arrival_at = now() + (travel_time * 0.5)

**Job 2: Return arrived fleets**
```
SELECT * FROM fleets WHERE status = 'returning' AND arrival_at <= now()
```
For each:
1. Set fleet status = 'stationed'
2. Clear destination_x/y, arrival_at

### 3.3 SP System

**Daily reset**: PvP worker checks `sp_last_reset`. If not today, reset to `max_space_points`.

**SP consumption**: 1 SP per PvP attack (deducted in AttackPlanet handler).

**SP Card use** (`inventory.go`):
```go
case "sp_grant":
    _, err = tx.Exec(`UPDATE players SET space_points = LEAST(space_points + $1, max_space_points) WHERE id = $2`, battleValue, playerID)
```

### 3.4 Truce Card Activation (`inventory.go`)

```go
case "protection":
    // Check no outgoing pending attacks
    // Check no incoming pending attacks within radar detection window
    duration := time.Duration(battleValue) * time.Hour
    _, err = tx.Exec(`UPDATE planets SET protection_until = $1 WHERE player_id = $2 AND is_homeworld = true`, time.Now().Add(duration), playerID)
```

Validation in AttackPlanet:
```go
// Check defender protection
var protUntil sql.NullTime
db.QueryRow(`SELECT protection_until FROM planets WHERE id = $1`, defenderPlanetID).Scan(&protUntil)
if protUntil.Valid && protUntil.Time.After(time.Now()) {
    return error("target under truce protection")
}
// Check attacker protection (bidirectional)
// If attacker is under truce, they can't attack either
```

### 3.5 Radar Incoming Attacks (`pvp.go`)

New endpoint: `GET /api/radar/incoming`

```go
func GetIncomingAttacks(w http.ResponseWriter, r *http.Request) {
    // Get player's radar level
    // Query pending_attacks targeting this player where status = 'traveling'
    // Filter info based on radar level:
    //   Lv1-2: arrival_at only
    //   Lv3-4: + attacker coordinates
    //   Lv5-6: + fleet count/strength
    //   Lv7-8: + fleet ship types
    //   Lv9:   + commander info
}
```

### 3.6 Cancel Attack (`pvp.go`)

New endpoint: `POST /api/pvp/cancel/{id}`
- Only works if pending_attack status = 'traveling'
- Sets status = 'cancelled'
- Sets fleets to 'returning' with return_time = time_elapsed * 0.5

### 3.7 Get Player SP (`resources.go` or new endpoint)

Include SP in resource response or add dedicated endpoint.

---

## 4. Frontend Changes

### 4.1 PvP Panel Rework
- Attack button → shows travel time estimate before confirming
- After sending: show fleet travel countdown (arrival_at)
- On arrival: auto-fetch combat result from combat_reports
- Show fleet return countdown

### 4.2 Radar Panel (new)
- Show incoming attacks list
- Info revealed based on radar building level
- Countdown timers per incoming attack

### 4.3 SP Display
- Add SP to ResourceHUD (e.g., "SP: 18/20")
- SP Card use in inventory updates SP display

### 4.4 Truce Protection Indicator
- Show shield icon on ResourceHUD when under truce protection
- Show remaining protection time

### 4.5 Fleet Status Indicators
- Fleet panel shows "Traveling" / "Returning" status with countdown
- Cannot modify traveling/returning fleets

---

## 5. Implementation Order

1. **Migration** — Add pending_attacks table, SP columns
2. **Backend: PvP worker** — Fleet arrival resolution + combat + return
3. **Backend: Refactor AttackPlanet** — Travel time calculation, pending_attack creation
4. **Backend: SP system** — Deduction, reset, SP Card activation
5. **Backend: Truce Card** — Activation logic, protection check in PvP
6. **Backend: Radar endpoint** — Incoming attacks query
7. **Backend: Cancel attack** — Fleet recall during travel
8. **Frontend: PvP panel rework** — Travel countdown, async results
9. **Frontend: Radar panel** — Incoming attacks display
10. **Frontend: SP + Truce UI** — HUD indicators

---

## 6. Scope Boundaries

**IN SCOPE:**
- Fleet travel time (distance-based)
- Asynchronous PvP combat via worker
- SP system (daily reset, consumption, SP Card)
- Truce Card / Adv Truce Card activation
- Radar incoming attacks
- Cancel/recall attack during travel
- Fleet return at 50% travel time

**OUT OF SCOPE (future phases):**
- He3 fuel consumption during travel (complex, needs more research)
- Scout attack mechanic (separate from main attack)
- Spacedock repair rework
- Fleet recall mid-travel (partially in scope as cancel)
- Multiple fleet coordination timing
- Wormhole return positions
