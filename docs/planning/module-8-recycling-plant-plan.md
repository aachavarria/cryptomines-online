# Module 8: Recycling Plant - Implementation Plan

**Date:** 2026-02-07
**Module:** Recycling Plant (Ship Scrapping System)
**Estimated Hours:** 8-12 hours (2-3 days)
**Priority:** MEDIUM (Quality of Life)

---

## 1. Executive Summary

### Overview
The Recycling Plant allows players to scrap unwanted or damaged ships for resource recovery. This is a simple CRUD system with confirmation UI to prevent accidental ship loss. Players select ships from their inventory (unassigned ships) or from fleets, confirm the scrap operation, and receive a percentage of the original build cost back as resources.

### Scope
**IN SCOPE:**
- Recycling Plant building (already exists - 11 levels, 2x2 grid, 3D model complete)
- Ship scrapping from player inventory (`ships` table with `quantity`)
- Ship scrapping from fleet stacks (`fleet_stacks` table)
- Resource recovery based on Recycling Plant level (5%-55% recovery)
- Confirmation modal with resource preview
- Backend endpoint: `POST /api/recycling/scrap`
- Frontend: RecyclingPanel.tsx component

**OUT OF SCOPE:**
- Quality Materials tech bonus (not yet implemented in Phase 1-2)
- Bulk scrap operations (scrap all damaged ships at once)
- Scrap history/logs
- Advanced filtering (by hull class, damage %, etc.)

### Key Design Decisions
1. **Two scrap sources:** Players can scrap ships from inventory (ships.quantity) OR from fleet stacks (fleet_stacks.ship_count)
2. **Recovery formula:** `recoveredResources = designCost × (recoveryPct / 100.0)`
3. **Recovery percentage:** Based on `recycling_plant_levels.recovery_pct` (5% at Lv1 → 55% at Lv11)
4. **Building requirement:** Recycling Plant must be built to unlock scrapping (check `player_buildings` for building_type='recycling_plant')
5. **Atomic operation:** Scrap transaction must be atomic (deduct ships + add resources in single DB transaction)

---

## 2. Current Database Analysis

### Existing Tables

#### `recycling_plant_levels`
```sql
CREATE TABLE recycling_plant_levels (
    level INTEGER PRIMARY KEY,
    civic_center_req INTEGER NOT NULL,
    ship_factory_req INTEGER NOT NULL DEFAULT 0,
    recovery_pct INTEGER NOT NULL,              -- 5% → 55%
    metal_cost BIGINT NOT NULL DEFAULT 0,
    he3_cost BIGINT NOT NULL DEFAULT 0,
    gold_cost BIGINT NOT NULL DEFAULT 0,
    build_time_seconds INTEGER NOT NULL DEFAULT 0
);
```

**Current Data:**
- Level 1: 5% recovery, requires Civic Center Lv2 + Ship Factory Lv1
- Level 5: 25% recovery
- Level 9: 45% recovery
- Level 11: 55% recovery (max)

#### `building_types` (Recycling Plant entry)
```sql
-- Already seeded in phase1_mvp.sql
('recycling_plant', 'Recycling Plant', 'military', 'ground', 500, 400, 550, 200,
 3.0300, 2.8700, 0, 1.0000, 11, 1, 'civic_center', 2,
 'Recovers resources from scrapped ships')
```

**Building Attributes:**
- Type: `military` / `ground`
- Max Level: 11
- Size: 2x2 grid
- Prerequisite: Civic Center Lv2

#### `player_buildings`
```sql
CREATE TABLE player_buildings (
    id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES players(id),
    building_type TEXT NOT NULL,
    level INTEGER NOT NULL DEFAULT 1,
    grid_x INTEGER,
    grid_y INTEGER,
    is_constructing BOOLEAN DEFAULT false,
    construction_finish_at TIMESTAMPTZ,
    ...
);
```

#### `ships` (Player ship inventory)
```sql
CREATE TABLE ships (
    id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES players(id),
    ship_design_id UUID NOT NULL REFERENCES ship_designs(id),
    quantity INTEGER NOT NULL DEFAULT 0,           -- Available ships (not in fleets)
    is_building BOOLEAN NOT NULL DEFAULT false,
    build_quantity INTEGER NOT NULL DEFAULT 0,
    build_finish_at TIMESTAMPTZ,
    production_slot INTEGER NOT NULL DEFAULT 1,
    ...
);
```

#### `fleet_stacks` (Ships in fleets)
```sql
CREATE TABLE fleet_stacks (
    id UUID PRIMARY KEY,
    fleet_id UUID NOT NULL REFERENCES fleets(id),
    ship_design_id UUID NOT NULL REFERENCES ship_designs(id),
    grid_row INTEGER NOT NULL CHECK (grid_row BETWEEN 0 AND 2),
    grid_col INTEGER NOT NULL CHECK (grid_col BETWEEN 0 AND 2),
    ship_count INTEGER NOT NULL DEFAULT 0 CHECK (ship_count BETWEEN 0 AND 3000),
    ...
);
```

#### `ship_designs` (Cost reference)
```sql
CREATE TABLE ship_designs (
    id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES players(id),
    name TEXT NOT NULL,
    hull_type_id INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL DEFAULT 0,          -- Build cost per ship
    he3_cost BIGINT NOT NULL DEFAULT 0,
    gold_cost BIGINT NOT NULL DEFAULT 0,
    build_time_seconds INTEGER NOT NULL DEFAULT 10,
    ...
);
```

### Database Completeness: 100%
**Status:** ✅ ALL TABLES EXIST - No migrations needed!

The database schema is complete:
- ✅ `recycling_plant_levels` table seeded (11 levels)
- ✅ `building_types` has Recycling Plant entry
- ✅ `player_buildings` tracks player's Recycling Plant
- ✅ `ships` table has quantity field for scrapping
- ✅ `fleet_stacks` table for scrapping ships from fleets
- ✅ `ship_designs` table has cost fields for recovery calculation

**Required:** Zero new tables or migrations.

---

## 3. Galaxy Online 2 Research

### Sources
Research conducted on Galaxy Online 2 wiki (game no longer active, archived information):

1. **[Recycling Plant Wiki](https://galaxyonlineii.fandom.com/wiki/Recycling_Plant)**
   - Building allows recycling old ships for resource recovery
   - Recovery percentage increases with building level
   - Formula: `recoveredResources = shipCost × recoveryPct`
   - 11 levels total (5% → 55% recovery)

2. **[Corps Mall Ship Scrap Value Table](https://galaxyonlineii.fandom.com/wiki/Corps_Mall_Ship_Scrap_Value_Table)**
   - Different ship types have same proportional recovery
   - "Ratio of cost to reward remains the same" regardless of ship type
   - Quality Materials tech can affect recovery (not in our Phase 2 scope)

### GO2 Mechanics Summary

| Level | Recovery % | Civic Center Req | Ship Factory Req | Build Time |
|-------|-----------|------------------|------------------|-----------|
| 1 | 5% | 2 | 1 | 3:20 |
| 2 | 10% | 3 | 3 | - |
| 3 | 15% | 4 | 5 | - |
| 4 | 20% | 5 | 7 | - |
| 5 | 25% | 6 | 9 | 4:42:49 |
| 6 | 30% | 7 | 9 | - |
| 7 | 35% | 8 | 11 | - |
| 8 | 40% | 9 | 11 | - |
| 9 | 45% | 10 | 13 | 421:25:27 |
| 10 | 50% | 11 | 13 | - |
| 11 | 55% | 12 | 15 | 4062:57:03 |

**Key Mechanics:**
1. **Building Requirement:** Recycling Plant must exist and be built to level X to get X% recovery
2. **Recovery Formula:** `Metal Recovered = Ship Metal Cost × (Recovery % / 100)` (same for He3, Gold)
3. **Quality Materials Bonus (OUT OF SCOPE):** In GO2, Quality Materials tech reduces ship build cost to 85%, which effectively increases recovery percentage when scrapping (e.g., 45% becomes ~53%). We're NOT implementing this in Phase 2.
4. **No Cooldowns:** Players can scrap ships instantly without cooldown timers
5. **No Undo:** Scrapping is permanent - ships are deleted and cannot be recovered

### Cryptomines Online Adaptations
- **1:1 Copy:** We're using the exact GO2 recovery percentages (5%-55%)
- **No Quality Materials:** We're ignoring the Quality Materials tech bonus for now (future enhancement)
- **Simplified UI:** GO2 had complex Corps Mall mechanics - we're focusing on basic player-owned ship scrapping only

---

## 4. System Architecture

### Component Overview
```
┌─────────────────────────────────────────────────────────────┐
│                     RECYCLING SYSTEM                         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐      ┌──────────────┐     ┌────────────┐│
│  │   Frontend   │◄────►│   Backend    │◄───►│  Database  ││
│  │              │      │              │     │            ││
│  │  Recycling   │ HTTP │    Scrap     │ SQL │   ships    ││
│  │    Panel     │      │   Handler    │     │fleet_stacks││
│  │              │      │              │     │ship_designs││
│  │  • Ship list │      │  • Validate  │     │  players   ││
│  │  • Confirm   │      │  • Calculate │     │            ││
│  │  • Preview   │      │  • Deduct    │     │            ││
│  └──────────────┘      │  • Refund    │     └────────────┘│
│                        └──────────────┘                    │
└─────────────────────────────────────────────────────────────┘
```

### Data Flow: Ship Scrapping

```
1. User clicks Recycling Plant building
   ↓
2. RecyclingPanel opens, fetches:
   - Recycling Plant level → recovery_pct
   - Player's ships (quantity > 0)
   - Player's fleet stacks (ship_count > 0)
   ↓
3. User selects ship design + quantity to scrap
   ↓
4. Frontend calculates preview:
   - Metal recovered = design.metal_cost × qty × (recovery_pct / 100)
   - He3 recovered = design.he3_cost × qty × (recovery_pct / 100)
   - Gold recovered = design.gold_cost × qty × (recovery_pct / 100)
   ↓
5. User confirms in modal
   ↓
6. POST /api/recycling/scrap
   {
     "source": "inventory" | "fleet",
     "ship_id": "uuid",              // ships.id or fleet_stacks.id
     "ship_design_id": "uuid",
     "quantity": 10
   }
   ↓
7. Backend validates:
   - Recycling Plant exists and has level > 0
   - Player owns the ships
   - Quantity available >= requested quantity
   ↓
8. Backend executes transaction:
   BEGIN;
     - Deduct ships.quantity (or fleet_stacks.ship_count)
     - Add resources to players (metal, he3, gold)
   COMMIT;
   ↓
9. Return success + recovered resources
   ↓
10. Frontend updates UI (ship list, resource HUD)
```

---

## 5. Backend Implementation

### 5.1 New Endpoint: POST /api/recycling/scrap

**File:** `/backend/internal/handlers/recycling.go`

```go
package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"

    "github.com/cryptomines-online/backend/internal/database"
    "github.com/cryptomines-online/backend/internal/middleware"
)

type scrapRequest struct {
    Source        string `json:"source"`          // "inventory" or "fleet"
    ShipID        string `json:"ship_id"`         // ships.id or fleet_stacks.id
    ShipDesignID  string `json:"ship_design_id"`  // for validation
    Quantity      int    `json:"quantity"`
}

type scrapResponse struct {
    Success         bool   `json:"success"`
    MetalRecovered  int64  `json:"metal_recovered"`
    He3Recovered    int64  `json:"he3_recovered"`
    GoldRecovered   int64  `json:"gold_recovered"`
    RemainingShips  int    `json:"remaining_ships"`
    Message         string `json:"message"`
}

// ScrapShips handles POST /api/recycling/scrap
func ScrapShips(w http.ResponseWriter, r *http.Request) {
    playerID := r.Context().Value(middleware.PlayerIDKey).(string)

    var req scrapRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
        return
    }

    // Validate input
    if req.Source != "inventory" && req.Source != "fleet" {
        http.Error(w, `{"error":"source must be 'inventory' or 'fleet'"}`, http.StatusBadRequest)
        return
    }
    if req.Quantity <= 0 {
        http.Error(w, `{"error":"quantity must be positive"}`, http.StatusBadRequest)
        return
    }

    // Get Recycling Plant level
    var plantLevel int
    var recoveryPct int
    err := database.DB.QueryRow(`
        SELECT pb.level, rpl.recovery_pct
        FROM player_buildings pb
        JOIN recycling_plant_levels rpl ON pb.level = rpl.level
        WHERE pb.player_id = $1 AND pb.building_type = 'recycling_plant'
    `, playerID).Scan(&plantLevel, &recoveryPct)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, `{"error":"recycling plant not built"}`, http.StatusForbidden)
            return
        }
        log.Printf("Failed to get recycling plant: %v", err)
        http.Error(w, `{"error":"failed to check recycling plant"}`, http.StatusInternalServerError)
        return
    }

    // Get ship design costs
    var metalCost, he3Cost, goldCost int64
    err = database.DB.QueryRow(`
        SELECT metal_cost, he3_cost, gold_cost
        FROM ship_designs
        WHERE id = $1 AND player_id = $2
    `, req.ShipDesignID, playerID).Scan(&metalCost, &he3Cost, &goldCost)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, `{"error":"ship design not found"}`, http.StatusNotFound)
            return
        }
        log.Printf("Failed to get ship design: %v", err)
        http.Error(w, `{"error":"failed to get ship design"}`, http.StatusInternalServerError)
        return
    }

    // Calculate recovery amounts
    recoveryFactor := float64(recoveryPct) / 100.0
    metalRecovered := int64(float64(metalCost) * float64(req.Quantity) * recoveryFactor)
    he3Recovered := int64(float64(he3Cost) * float64(req.Quantity) * recoveryFactor)
    goldRecovered := int64(float64(goldCost) * float64(req.Quantity) * recoveryFactor)

    // Begin transaction
    tx, err := database.DB.Begin()
    if err != nil {
        log.Printf("Failed to begin transaction: %v", err)
        http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    var remainingShips int

    if req.Source == "inventory" {
        // Scrap from ships table (player inventory)
        var currentQty int
        err = tx.QueryRow(`
            SELECT quantity
            FROM ships
            WHERE id = $1 AND player_id = $2 AND ship_design_id = $3
            FOR UPDATE
        `, req.ShipID, playerID, req.ShipDesignID).Scan(&currentQty)
        if err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, `{"error":"ship not found or not owned"}`, http.StatusNotFound)
                return
            }
            log.Printf("Failed to get ship: %v", err)
            http.Error(w, `{"error":"failed to get ship"}`, http.StatusInternalServerError)
            return
        }

        if currentQty < req.Quantity {
            http.Error(w, `{"error":"insufficient ships to scrap"}`, http.StatusBadRequest)
            return
        }

        // Deduct ships
        remainingShips = currentQty - req.Quantity
        _, err = tx.Exec(`
            UPDATE ships
            SET quantity = $1, updated_at = now()
            WHERE id = $2
        `, remainingShips, req.ShipID)
        if err != nil {
            log.Printf("Failed to deduct ships: %v", err)
            http.Error(w, `{"error":"failed to deduct ships"}`, http.StatusInternalServerError)
            return
        }

    } else { // source == "fleet"
        // Scrap from fleet_stacks
        var currentCount int
        var fleetID string
        err = tx.QueryRow(`
            SELECT fs.ship_count, fs.fleet_id
            FROM fleet_stacks fs
            JOIN fleets f ON fs.fleet_id = f.id
            WHERE fs.id = $1 AND f.player_id = $2 AND fs.ship_design_id = $3
            FOR UPDATE
        `, req.ShipID, playerID, req.ShipDesignID).Scan(&currentCount, &fleetID)
        if err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, `{"error":"fleet stack not found or not owned"}`, http.StatusNotFound)
                return
            }
            log.Printf("Failed to get fleet stack: %v", err)
            http.Error(w, `{"error":"failed to get fleet stack"}`, http.StatusInternalServerError)
            return
        }

        if currentCount < req.Quantity {
            http.Error(w, `{"error":"insufficient ships in fleet stack"}`, http.StatusBadRequest)
            return
        }

        // Deduct ships from fleet stack
        remainingShips = currentCount - req.Quantity
        if remainingShips == 0 {
            // Delete empty stack
            _, err = tx.Exec(`DELETE FROM fleet_stacks WHERE id = $1`, req.ShipID)
        } else {
            _, err = tx.Exec(`
                UPDATE fleet_stacks
                SET ship_count = $1
                WHERE id = $2
            `, remainingShips, req.ShipID)
        }
        if err != nil {
            log.Printf("Failed to update fleet stack: %v", err)
            http.Error(w, `{"error":"failed to update fleet stack"}`, http.StatusInternalServerError)
            return
        }
    }

    // Add recovered resources to player
    _, err = tx.Exec(`
        UPDATE players
        SET metal = metal + $1,
            he3 = he3 + $2,
            gold = gold + $3,
            updated_at = now()
        WHERE id = $4
    `, metalRecovered, he3Recovered, goldRecovered, playerID)
    if err != nil {
        log.Printf("Failed to add resources: %v", err)
        http.Error(w, `{"error":"failed to add resources"}`, http.StatusInternalServerError)
        return
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        log.Printf("Failed to commit transaction: %v", err)
        http.Error(w, `{"error":"failed to commit transaction"}`, http.StatusInternalServerError)
        return
    }

    // Success response
    response := scrapResponse{
        Success:        true,
        MetalRecovered: metalRecovered,
        He3Recovered:   he3Recovered,
        GoldRecovered:  goldRecovered,
        RemainingShips: remainingShips,
        Message:        "Ships scrapped successfully",
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

### 5.2 Route Registration

**File:** `/backend/cmd/server/main.go` (update)

```go
// Recycling endpoints
http.HandleFunc("POST /api/recycling/scrap", middleware.AuthMiddleware(handlers.ScrapShips))
```

### 5.3 Backend Tasks Breakdown

| Task | Description | Hours |
|------|-------------|-------|
| Create handlers/recycling.go | Implement ScrapShips handler with validation + transaction logic | 2.5 |
| Route registration | Add POST /api/recycling/scrap route | 0.25 |
| Error handling | Comprehensive error messages for all failure cases | 0.5 |
| Testing (manual) | Test with Postman/curl for all edge cases | 0.75 |
| **Subtotal** | | **4 hours** |

---

## 6. Frontend Implementation

### 6.1 New Component: RecyclingPanel.tsx

**File:** `/frontend/src/components/panels/RecyclingPanel.tsx`

```typescript
import { useState, useEffect } from 'react'
import { useGame } from '../../contexts/GameContext'
import { scrapShips, getRecyclingPlantLevel } from '../../api/recycling'

interface ShipOption {
  id: string                    // ships.id or fleet_stacks.id
  source: 'inventory' | 'fleet'
  designId: string
  designName: string
  quantity: number              // available quantity
  metalCost: number
  he3Cost: number
  goldCost: number
  fleetName?: string            // for fleet ships
}

export default function RecyclingPanel() {
  const { buildings, ships, fleets, resources, refreshResources, refreshShips } = useGame()
  const [shipOptions, setShipOptions] = useState<ShipOption[]>([])
  const [selectedShip, setSelectedShip] = useState<ShipOption | null>(null)
  const [scrapQuantity, setScrapQuantity] = useState(1)
  const [recoveryPct, setRecoveryPct] = useState(0)
  const [showConfirm, setShowConfirm] = useState(false)
  const [loading, setLoading] = useState(false)

  // Fetch recycling plant level and build ship options
  useEffect(() => {
    const recyclingPlant = buildings.find(b => b.building_type === 'recycling_plant')
    if (!recyclingPlant) {
      // No recycling plant built
      return
    }

    // Fetch recovery percentage
    getRecyclingPlantLevel(recyclingPlant.level).then(data => {
      setRecoveryPct(data.recovery_pct)
    })

    // Build ship options from inventory
    const inventoryOptions: ShipOption[] = ships
      .filter(s => s.quantity > 0 && !s.is_building)
      .map(s => ({
        id: s.id,
        source: 'inventory' as const,
        designId: s.ship_design_id,
        designName: s.design_name,
        quantity: s.quantity,
        metalCost: s.metal_cost,
        he3Cost: s.he3_cost,
        goldCost: s.gold_cost,
      }))

    // Build ship options from fleets
    const fleetOptions: ShipOption[] = []
    fleets.forEach(fleet => {
      fleet.stacks.forEach(stack => {
        if (stack.ship_count > 0) {
          fleetOptions.push({
            id: stack.id,
            source: 'fleet' as const,
            designId: stack.ship_design_id,
            designName: stack.design_name,
            quantity: stack.ship_count,
            metalCost: stack.metal_cost,
            he3Cost: stack.he3_cost,
            goldCost: stack.gold_cost,
            fleetName: fleet.name,
          })
        }
      })
    })

    setShipOptions([...inventoryOptions, ...fleetOptions])
  }, [buildings, ships, fleets])

  const handleScrap = async () => {
    if (!selectedShip) return

    setLoading(true)
    try {
      const result = await scrapShips({
        source: selectedShip.source,
        ship_id: selectedShip.id,
        ship_design_id: selectedShip.designId,
        quantity: scrapQuantity,
      })

      // Success - refresh data
      await refreshResources()
      await refreshShips()

      setShowConfirm(false)
      setSelectedShip(null)
      setScrapQuantity(1)

      alert(`Scrapped ${scrapQuantity} ships!\nRecovered:\n` +
            `Metal: ${result.metal_recovered}\n` +
            `He3: ${result.he3_recovered}\n` +
            `Gold: ${result.gold_recovered}`)
    } catch (err: any) {
      alert(`Failed to scrap ships: ${err.message}`)
    } finally {
      setLoading(false)
    }
  }

  const calculatePreview = () => {
    if (!selectedShip) return { metal: 0, he3: 0, gold: 0 }
    const factor = recoveryPct / 100
    return {
      metal: Math.floor(selectedShip.metalCost * scrapQuantity * factor),
      he3: Math.floor(selectedShip.he3Cost * scrapQuantity * factor),
      gold: Math.floor(selectedShip.goldCost * scrapQuantity * factor),
    }
  }

  const recyclingPlant = buildings.find(b => b.building_type === 'recycling_plant')
  if (!recyclingPlant) {
    return (
      <div className="recycling-panel">
        <h2>Recycling Plant</h2>
        <p>Build a Recycling Plant to scrap ships for resources.</p>
      </div>
    )
  }

  const preview = calculatePreview()

  return (
    <div className="recycling-panel">
      <h2>Recycling Plant (Level {recyclingPlant.level})</h2>
      <p>Recovery Rate: {recoveryPct}%</p>

      <div className="ship-list">
        <h3>Select Ships to Scrap</h3>
        {shipOptions.length === 0 && <p>No ships available to scrap.</p>}
        {shipOptions.map(ship => (
          <div
            key={ship.id}
            className={`ship-option ${selectedShip?.id === ship.id ? 'selected' : ''}`}
            onClick={() => {
              setSelectedShip(ship)
              setScrapQuantity(1)
            }}
          >
            <strong>{ship.designName}</strong>
            <span>Quantity: {ship.quantity}</span>
            {ship.fleetName && <span>Fleet: {ship.fleetName}</span>}
            <span>Cost: {ship.metalCost}M / {ship.he3Cost}H / {ship.goldCost}G</span>
          </div>
        ))}
      </div>

      {selectedShip && (
        <div className="scrap-controls">
          <h3>Scrap {selectedShip.designName}</h3>
          <label>
            Quantity:
            <input
              type="number"
              min={1}
              max={selectedShip.quantity}
              value={scrapQuantity}
              onChange={e => setScrapQuantity(Math.max(1, Math.min(selectedShip.quantity, parseInt(e.target.value) || 1)))}
            />
          </label>

          <div className="preview">
            <h4>Resources Recovered ({recoveryPct}%)</h4>
            <p>Metal: {preview.metal}</p>
            <p>He3: {preview.he3}</p>
            <p>Gold: {preview.gold}</p>
          </div>

          <button onClick={() => setShowConfirm(true)} disabled={loading}>
            Scrap Ships
          </button>
        </div>
      )}

      {showConfirm && (
        <div className="confirm-modal">
          <div className="modal-content">
            <h3>⚠️ Confirm Scrap</h3>
            <p>
              Are you sure you want to scrap <strong>{scrapQuantity} {selectedShip?.designName}</strong>?
            </p>
            <p>This action cannot be undone!</p>
            <p>You will receive:</p>
            <ul>
              <li>Metal: {preview.metal}</li>
              <li>He3: {preview.he3}</li>
              <li>Gold: {preview.gold}</li>
            </ul>
            <button onClick={handleScrap} disabled={loading}>
              {loading ? 'Scrapping...' : 'Confirm Scrap'}
            </button>
            <button onClick={() => setShowConfirm(false)} disabled={loading}>
              Cancel
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
```

### 6.2 API Function

**File:** `/frontend/src/api/recycling.ts`

```typescript
import { apiRequest } from './base'

export interface ScrapRequest {
  source: 'inventory' | 'fleet'
  ship_id: string
  ship_design_id: string
  quantity: number
}

export interface ScrapResponse {
  success: boolean
  metal_recovered: number
  he3_recovered: number
  gold_recovered: number
  remaining_ships: number
  message: string
}

export async function scrapShips(request: ScrapRequest): Promise<ScrapResponse> {
  return apiRequest<ScrapResponse>('/api/recycling/scrap', {
    method: 'POST',
    body: JSON.stringify(request),
  })
}

export async function getRecyclingPlantLevel(level: number): Promise<{ recovery_pct: number }> {
  // Fetch from local config or backend
  const recoveryMap: Record<number, number> = {
    1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30,
    7: 35, 8: 40, 9: 45, 10: 50, 11: 55,
  }
  return { recovery_pct: recoveryMap[level] || 5 }
}
```

### 6.3 Open Panel on Building Click

**File:** `/frontend/src/contexts/GameContext.tsx` (update)

```typescript
// Add to handleBuildingClick in GameContext
if (building.building_type === 'recycling_plant') {
  setActivePanel('recycling')
  return
}
```

**File:** `/frontend/src/App.tsx` (update)

```typescript
import RecyclingPanel from './components/panels/RecyclingPanel'

// Add to panel rendering logic
{activePanel === 'recycling' && <RecyclingPanel />}
```

### 6.4 Frontend Tasks Breakdown

| Task | Description | Hours |
|------|-------------|-------|
| Create RecyclingPanel.tsx | Build ship list, quantity selector, preview calculation | 2.0 |
| Create api/recycling.ts | API request functions | 0.25 |
| Confirmation modal | Warning modal with resource preview | 0.5 |
| GameContext integration | Open RecyclingPanel on building click | 0.25 |
| Styling | CSS for RecyclingPanel (list, controls, modal) | 0.5 |
| Testing | Manual UI testing for all flows | 0.5 |
| **Subtotal** | | **4 hours** |

---

## 7. QA Checklist

### 7.1 Backend Tests

**Manual Testing (Go):**

| # | Test Case | Expected Result | Status |
|---|-----------|----------------|--------|
| 1 | POST /api/recycling/scrap with no Recycling Plant | 403 Forbidden: "recycling plant not built" | ⬜ |
| 2 | POST /api/recycling/scrap with Recycling Plant Lv1 (5% recovery) | 200 OK, 5% resources recovered | ⬜ |
| 3 | POST /api/recycling/scrap with Recycling Plant Lv11 (55% recovery) | 200 OK, 55% resources recovered | ⬜ |
| 4 | Scrap 10 ships from inventory (ships.quantity = 50) | ships.quantity = 40, resources added | ⬜ |
| 5 | Scrap all ships from inventory (quantity = 10) | ships.quantity = 0 | ⬜ |
| 6 | Scrap ships with quantity > available | 400 Bad Request: "insufficient ships" | ⬜ |
| 7 | Scrap 500 ships from fleet_stacks (ship_count = 1000) | ship_count = 500 | ⬜ |
| 8 | Scrap all ships from fleet_stacks (ship_count = 100) | fleet_stacks row deleted | ⬜ |
| 9 | Scrap ships not owned by player | 404 Not Found: "ship not found or not owned" | ⬜ |
| 10 | Scrap with invalid source ("invalid") | 400 Bad Request: "source must be 'inventory' or 'fleet'" | ⬜ |
| 11 | Scrap with quantity = 0 | 400 Bad Request: "quantity must be positive" | ⬜ |
| 12 | Scrap with quantity = -5 | 400 Bad Request: "quantity must be positive" | ⬜ |
| 13 | Verify transaction atomicity (rollback on resource update failure) | All changes rolled back on error | ⬜ |
| 14 | Verify Metal/He3/Gold resources increase correctly | players.metal/he3/gold += calculated amounts | ⬜ |
| 15 | Scrap ship design with 0 costs (edge case) | 0 resources recovered (but operation succeeds) | ⬜ |

### 7.2 Frontend Tests

**Manual Testing (UI):**

| # | Test Case | Expected Result | Status |
|---|-----------|----------------|--------|
| 16 | Open Recycling Plant building (no plant built) | Message: "Build a Recycling Plant to scrap ships" | ⬜ |
| 17 | Open Recycling Plant Lv5 | Panel shows "Level 5, Recovery Rate: 25%" | ⬜ |
| 18 | Ship list shows inventory ships (quantity > 0) | All non-building ships displayed | ⬜ |
| 19 | Ship list shows fleet ships (ship_count > 0) | All fleet stacks displayed with fleet names | ⬜ |
| 20 | Ship list empty (no ships available) | Message: "No ships available to scrap" | ⬜ |
| 21 | Select ship, preview shows correct recovery amounts | Metal/He3/Gold = cost × qty × (recovery_pct / 100) | ⬜ |
| 22 | Change quantity slider, preview updates | Preview recalculates dynamically | ⬜ |
| 23 | Click "Scrap Ships", confirmation modal opens | Modal shows ship name, quantity, preview | ⬜ |
| 24 | Confirm scrap in modal | Ships deducted, resources added, panel refreshes | ⬜ |
| 25 | Cancel scrap in modal | Modal closes, no changes made | ⬜ |
| 26 | Scrap all ships (quantity = available) | Ship removed from list after scrap | ⬜ |
| 27 | Resource HUD updates after scrap | Metal/He3/Gold counts increase | ⬜ |
| 28 | Error message on backend failure | Alert shows error message | ⬜ |
| 29 | Loading state during scrap | Button disabled, "Scrapping..." text | ⬜ |
| 30 | Scrap from fleet, fleet panel updates | Fleet stack quantity decreases | ⬜ |

### 7.3 Integration Tests

| # | Test Case | Expected Result | Status |
|---|-----------|----------------|--------|
| 31 | End-to-end: Build Recycling Plant → Scrap ships → Verify resources | Full flow works without errors | ⬜ |
| 32 | Scrap ships from multiple sources (inventory + fleet) | Both sources work independently | ⬜ |
| 33 | Upgrade Recycling Plant, verify recovery % increases | Higher level = higher recovery | ⬜ |
| 34 | Scrap ships while building is in progress | Scrapping allowed (no construction check) | ⬜ |
| 35 | Concurrent scrap requests (race condition test) | Database locks prevent double-scrap | ⬜ |

---

## 8. Edge Cases & Error Handling

### 8.1 Edge Cases

| Case | Handling |
|------|----------|
| No Recycling Plant built | 403 Forbidden: "recycling plant not built" |
| Ship with 0 cost | Allow scrap, return 0 resources (no error) |
| Scrap quantity > available | 400 Bad Request: "insufficient ships" |
| Empty fleet stack after scrap | DELETE fleet_stacks row (ship_count = 0) |
| Concurrent scrap requests | Database `FOR UPDATE` lock prevents race conditions |
| Ship design deleted after selection | 404 Not Found: "ship design not found" |
| Recycling Plant destroyed during scrap | Transaction fails, rolled back |
| Player has max resources (overflow?) | Allow overflow (no cap on resources in Phase 2) |

### 8.2 Validation Rules

**Backend:**
1. ✅ Player must have Recycling Plant built (level >= 1)
2. ✅ Source must be "inventory" or "fleet"
3. ✅ Quantity must be positive integer
4. ✅ Player must own the ships (player_id match)
5. ✅ Available quantity >= requested quantity
6. ✅ Ship design must exist and belong to player

**Frontend:**
1. ✅ Quantity input capped at available quantity
2. ✅ Confirmation modal required for all scrap operations
3. ✅ Disable scrap button while loading
4. ✅ Show clear error messages for all failures

---

## 9. Risk Assessment

| Risk | Impact | Likelihood | Mitigation |
|------|--------|-----------|------------|
| **Transaction failure during scrap** | HIGH | LOW | Use database transactions (BEGIN/COMMIT/ROLLBACK) |
| **Race condition (double scrap)** | MEDIUM | LOW | Use `FOR UPDATE` locks on ships/fleet_stacks |
| **Resource overflow** | LOW | LOW | No resource cap in Phase 2 (allow overflow) |
| **Accidental ship loss** | MEDIUM | MEDIUM | Require confirmation modal with preview |
| **Recovery % calculation error** | HIGH | LOW | Unit tests for recovery formula |
| **Frontend state desync after scrap** | MEDIUM | MEDIUM | Call refreshResources() + refreshShips() after scrap |

**Critical Paths:**
1. Transaction atomicity (deduct ships + add resources must succeed/fail together)
2. Ownership validation (prevent scrapping other players' ships)
3. Quantity validation (prevent scrapping more ships than available)

---

## 10. Timeline & Estimates

### 10.1 Task Breakdown

| Phase | Task | Owner | Hours | Dependencies |
|-------|------|-------|-------|-------------|
| **Backend** | | | | |
| 1 | Create handlers/recycling.go | backend-dev | 2.5 | - |
| 2 | Route registration | backend-dev | 0.25 | Task 1 |
| 3 | Error handling | backend-dev | 0.5 | Task 1 |
| 4 | Manual testing (Postman) | backend-dev | 0.75 | Task 1-3 |
| **Frontend** | | | | |
| 5 | Create RecyclingPanel.tsx | frontend-dev | 2.0 | Backend Task 1 |
| 6 | Create api/recycling.ts | frontend-dev | 0.25 | - |
| 7 | Confirmation modal | frontend-dev | 0.5 | Task 5 |
| 8 | GameContext integration | frontend-dev | 0.25 | Task 5 |
| 9 | Styling (CSS) | frontend-dev | 0.5 | Task 5 |
| 10 | Manual UI testing | frontend-dev | 0.5 | Task 5-9 |
| **QA** | | | | |
| 11 | Backend QA (15 tests) | qa-agent | 1.0 | Backend tasks |
| 12 | Frontend QA (15 tests) | qa-agent | 1.0 | Frontend tasks |
| 13 | Integration QA (5 tests) | qa-agent | 0.5 | All tasks |
| **TOTAL** | | | **10.5 hours** | |

### 10.2 Timeline (Parallel Execution)

**Day 1 (4 hours):**
- Backend: Tasks 1-3 (handlers, routes, error handling) - 3.25 hours
- Frontend: Tasks 6, 8 (API setup, GameContext) - 0.5 hours

**Day 2 (5 hours):**
- Backend: Task 4 (testing) - 0.75 hours
- Frontend: Tasks 5, 7, 9 (RecyclingPanel, modal, styling) - 3.0 hours
- Frontend: Task 10 (UI testing) - 0.5 hours
- QA: Task 11 (backend tests) - 1.0 hours

**Day 3 (2.5 hours):**
- QA: Tasks 12-13 (frontend + integration tests) - 1.5 hours
- Bugfixes + polish - 1.0 hours

**Total:** 10.5 hours (2-3 days with parallel work)

---

## 11. Success Criteria

### 11.1 Must-Have (MVP)
- ✅ Recycling Plant building functional (already complete)
- ✅ Players can scrap ships from inventory (`ships` table)
- ✅ Players can scrap ships from fleets (`fleet_stacks` table)
- ✅ Resource recovery based on Recycling Plant level (5%-55%)
- ✅ Confirmation modal prevents accidental scrap
- ✅ Backend validates ownership + quantity
- ✅ Transaction atomicity (all-or-nothing)
- ✅ Resource HUD updates after scrap

### 11.2 Nice-to-Have (Future)
- ⬜ Quality Materials tech bonus (increases recovery %)
- ⬜ Bulk scrap (scrap all damaged ships at once)
- ⬜ Scrap history log (track scrapped ships)
- ⬜ Advanced filtering (by hull class, damage %, etc.)
- ⬜ Scrap preview before clicking (hover tooltip)

### 11.3 Definition of Done
- All 35 QA tests pass
- Backend endpoint handles all error cases gracefully
- Frontend RecyclingPanel matches GO2 mechanics (1:1 copy)
- No resource loss bugs (atomicity verified)
- User cannot accidentally scrap ships (confirmation required)

---

## 12. Future Enhancements

### Phase 3+ Additions
1. **Quality Materials Tech:** Reduce ship build cost by 15% → effective recovery increases (e.g., 45% becomes ~53%)
2. **Bulk Scrap:** "Scrap All Damaged" button for ships below X% HP
3. **Scrap Analytics:** Show total resources recovered (lifetime stats)
4. **Scrap History:** Log all scrapped ships with timestamps
5. **Advanced Filters:** Filter ships by hull class, tier, damage %, etc.
6. **Scrap Queue:** Queue multiple scrap operations (unlikely needed)

---

## 13. Appendix

### 13.1 GO2 References
- [Recycling Plant Wiki](https://galaxyonlineii.fandom.com/wiki/Recycling_Plant) - Building stats and recovery percentages
- [Corps Mall Ship Scrap Value Table](https://galaxyonlineii.fandom.com/wiki/Corps_Mall_Ship_Scrap_Value_Table) - Ship scrap value mechanics

### 13.2 Database Schema Reference

**Recycling Plant Levels:**
```sql
SELECT level, recovery_pct FROM recycling_plant_levels;
-- Level 1: 5%
-- Level 2: 10%
-- Level 3: 15%
-- Level 4: 20%
-- Level 5: 25%
-- Level 6: 30%
-- Level 7: 35%
-- Level 8: 40%
-- Level 9: 45%
-- Level 10: 50%
-- Level 11: 55%
```

**Ship Cost Example:**
```sql
SELECT name, metal_cost, he3_cost, gold_cost
FROM ship_designs
WHERE player_id = 'xxx';

-- Example: Frigate design
-- metal_cost: 10000
-- he3_cost: 5000
-- gold_cost: 500

-- With Recycling Plant Lv5 (25% recovery), scrapping 10 frigates:
-- Metal recovered: 10000 × 10 × 0.25 = 25000
-- He3 recovered: 5000 × 10 × 0.25 = 12500
-- Gold recovered: 500 × 10 × 0.25 = 1250
```

### 13.3 Formula Reference

**Recovery Calculation:**
```
recoveredMetal = ship.metal_cost × quantity × (recovery_pct / 100.0)
recoveredHe3 = ship.he3_cost × quantity × (recovery_pct / 100.0)
recoveredGold = ship.gold_cost × quantity × (recovery_pct / 100.0)
```

**Example (Lv9 Recycling Plant, 45% recovery):**
- Ship cost: 50000 Metal, 30000 He3, 2000 Gold
- Quantity: 25 ships
- Recovered:
  - Metal: 50000 × 25 × 0.45 = 562,500
  - He3: 30000 × 25 × 0.45 = 337,500
  - Gold: 2000 × 25 × 0.45 = 22,500

---

## 14. Implementation Notes

### Backend Notes
- Use `FOR UPDATE` locks to prevent race conditions on concurrent scrap requests
- Transaction must be atomic: deduct ships + add resources in one transaction
- Delete empty `fleet_stacks` rows (ship_count = 0) to keep database clean
- Recovery percentage is integer (5-55), divide by 100.0 for float calculation
- No need to check if ships are in combat (out of scope for Phase 2)

### Frontend Notes
- Fetch Recycling Plant level on panel open (not on every render)
- Show ships from both inventory (`ships`) and fleets (`fleet_stacks`)
- Quantity input must be capped at available quantity
- Preview calculation must match backend formula exactly
- Refresh resources + ships after successful scrap (avoid stale data)
- Confirmation modal is REQUIRED (prevent accidental ship loss)

### Testing Notes
- Test with Recycling Plant Lv1 (5%), Lv6 (30%), Lv11 (55%)
- Test scrapping partial quantity (10 out of 50)
- Test scrapping all ships (quantity = available)
- Test scrapping from empty fleet stack (should fail)
- Test concurrent scrap requests (database locks should prevent double-deduct)
- Test transaction rollback (simulate resource update failure)

---

**END OF PLAN**
