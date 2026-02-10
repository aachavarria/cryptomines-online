# Module 3: Resource Auto-Production - Implementation Plan

**Date:** 2026-02-07
**Architect:** Claude (Sonnet)
**Status:** READY FOR IMPLEMENTATION
**Estimated Duration:** 2-3 days

---

## Executive Summary

Module 3 implements **automatic warehouse resource accumulation** to replace the current manual collection system.

**Current System (Manual Collection):**
- Resources tracked: `metal`, `he3`, `gold` (player's direct balance)
- Production rates: `metal_per_hour`, `he3_per_hour`, `gold_per_hour`
- Collection: Player clicks "Collect" → calculates elapsed time × rate → adds to balance
- Problem: Resources only accumulate when player actively collects

**New System (Auto-Production with Warehouse):**
- **Warehouse fields:** `warehouse_metal`, `warehouse_he3`, `warehouse_gold` (auto-accumulate)
- **Worker:** Updates warehouse every 5 minutes, caps at `storage_capacity`
- **Collection:** Transfer warehouse → player balance (clear warehouse)
- **Benefit:** Resources accumulate automatically even when offline

**Current State:** 30% implemented (production rates exist, collect endpoint exists)

**Work Remaining:** 70%
1. Add warehouse fields to `resources` table
2. Create resource worker to auto-accumulate warehouse
3. Update collect endpoint to transfer warehouse → player
4. Update frontend to show warehouse vs player balance
5. QA testing

---

## 1. Current Implementation Analysis

### ✅ What Exists

**Database Schema:**
```sql
resources table:
  id, planet_id,
  metal, he3, gold,  -- Player's direct balance
  metal_per_hour, he3_per_hour, gold_per_hour,  -- Production rates
  storage_capacity,  -- Warehouse cap
  last_collected_at,  -- Last collect timestamp
  updated_at
```

**Backend (`resources.go`):**
- `GET /api/planets/{id}/resources` - Returns resources + pending amounts
- `POST /api/planets/{id}/resources/collect` - Calculates elapsed time, adds to balance
- `recalculateProductionRates()` - Updates production rates when buildings change
- `applyCompletedUpgrades()` - Completes building upgrades, triggers recalculation

**Current Collection Logic (lines 117-131):**
```go
elapsed := time.Since(res.LastCollectedAt).Hours()
collectedMetal := int64(float64(res.MetalPerHour) * elapsed)
collectedHe3 := int64(float64(res.He3PerHour) * elapsed)
collectedGold := int64(float64(res.GoldPerHour) * elapsed)

// Cap resources at storage capacity
newMetal := min64(res.Metal+collectedMetal, res.StorageCapacity)
newHe3 := min64(res.He3+collectedHe3, res.StorageCapacity)
newGold := min64(res.Gold+collectedGold, res.StorageCapacity)
```

**Frontend:**
- `ResourceHUD.tsx` - Displays resources, collect button, pending amounts
- `useResources.ts` - Fetches resources every 30s, handles collect action
- Shows `pending_metal/he3/gold` (calculated from elapsed time)

### ❌ What's Missing

**Problem with Current System:**
1. **No persistent accumulation** - Resources only calculated on-demand when player views
2. **Pending calculation on every GET** - Inefficient, recalculates elapsed time every request
3. **No warehouse separation** - `metal/he3/gold` fields serve dual purpose (confusing)

**Needed Changes:**
1. **Warehouse fields** - Separate storage for auto-accumulated resources
2. **Worker** - Periodically update warehouse (every 5 min)
3. **Collect endpoint** - Transfer warehouse → player balance
4. **Frontend** - Show warehouse amounts separately

---

## 2. Database Changes

### Migration: Add Warehouse Fields

**File:** `/supabase/migrations/YYYYMMDDHHMMSS_add_warehouse_fields.sql` (NEW)

```sql
-- Add warehouse fields to resources table
ALTER TABLE resources
ADD COLUMN warehouse_metal BIGINT NOT NULL DEFAULT 0,
ADD COLUMN warehouse_he3 BIGINT NOT NULL DEFAULT 0,
ADD COLUMN warehouse_gold BIGINT NOT NULL DEFAULT 0,
ADD COLUMN last_warehouse_update TIMESTAMPTZ NOT NULL DEFAULT now();

-- Add constraints for non-negative warehouse values
ALTER TABLE resources
ADD CONSTRAINT chk_warehouse_non_negative CHECK (
    warehouse_metal >= 0 AND warehouse_he3 >= 0 AND warehouse_gold >= 0
);

-- Update existing rows: initialize warehouse with current pending amounts
UPDATE resources
SET
    warehouse_metal = LEAST(
        GREATEST(0, CAST(EXTRACT(EPOCH FROM (now() - last_collected_at)) / 3600.0 * metal_per_hour AS BIGINT)),
        storage_capacity
    ),
    warehouse_he3 = LEAST(
        GREATEST(0, CAST(EXTRACT(EPOCH FROM (now() - last_collected_at)) / 3600.0 * he3_per_hour AS BIGINT)),
        storage_capacity
    ),
    warehouse_gold = LEAST(
        GREATEST(0, CAST(EXTRACT(EPOCH FROM (now() - last_collected_at)) / 3600.0 * gold_per_hour AS BIGINT)),
        storage_capacity
    ),
    last_warehouse_update = now();

-- IMPORTANT: Do NOT alter last_collected_at or metal/he3/gold
-- These remain as player's direct balance
```

**Schema After Migration:**
```sql
resources table:
  id, planet_id,
  -- Player's direct balance (used for spending)
  metal, he3, gold,
  -- Auto-accumulating warehouse (collect to transfer)
  warehouse_metal, warehouse_he3, warehouse_gold,
  -- Production rates
  metal_per_hour, he3_per_hour, gold_per_hour,
  storage_capacity,
  last_collected_at,      -- Last time player clicked "Collect"
  last_warehouse_update,  -- Last time worker updated warehouse
  updated_at
```

---

## 3. Resource Worker Implementation

### File: `/backend/internal/workers/resource_worker.go` (NEW)

```go
package workers

import (
    "log"
    "time"
    "github.com/cryptomines-online/backend/internal/database"
)

// UpdateAllWarehouses updates warehouse resources for all planets
func UpdateAllWarehouses() {
    now := time.Now()

    // Find all planets with production
    rows, err := database.DB.Query(`
        SELECT r.planet_id, r.metal_per_hour, r.he3_per_hour, r.gold_per_hour,
               r.warehouse_metal, r.warehouse_he3, r.warehouse_gold,
               r.storage_capacity, r.last_warehouse_update
        FROM resources r
        WHERE r.metal_per_hour > 0 OR r.he3_per_hour > 0 OR r.gold_per_hour > 0
    `)
    if err != nil {
        log.Printf("[ResourceWorker] Failed to query resources: %v", err)
        return
    }

    type planetResource struct {
        planetID             string
        metalPerHour         int64
        he3PerHour           int64
        goldPerHour          int64
        warehouseMetal       int64
        warehouseHe3         int64
        warehouseGold        int64
        storageCapacity      int64
        lastWarehouseUpdate  time.Time
    }

    var planets []planetResource

    for rows.Next() {
        var pr planetResource
        rows.Scan(
            &pr.planetID, &pr.metalPerHour, &pr.he3PerHour, &pr.goldPerHour,
            &pr.warehouseMetal, &pr.warehouseHe3, &pr.warehouseGold,
            &pr.storageCapacity, &pr.lastWarehouseUpdate,
        )
        planets = append(planets, pr)
    }
    rows.Close()

    if len(planets) == 0 {
        return
    }

    log.Printf("[ResourceWorker] Updating warehouses for %d planets", len(planets))

    for _, pr := range planets {
        // Calculate elapsed time since last warehouse update
        elapsed := now.Sub(pr.lastWarehouseUpdate).Hours()

        // Calculate production in elapsed time
        producedMetal := int64(float64(pr.metalPerHour) * elapsed)
        producedHe3 := int64(float64(pr.he3PerHour) * elapsed)
        producedGold := int64(float64(pr.goldPerHour) * elapsed)

        // Add to warehouse, cap at storage capacity
        newWarehouseMetal := min64(pr.warehouseMetal+producedMetal, pr.storageCapacity)
        newWarehouseHe3 := min64(pr.warehouseHe3+producedHe3, pr.storageCapacity)
        newWarehouseGold := min64(pr.warehouseGold+producedGold, pr.storageCapacity)

        // Update warehouse
        _, err := database.DB.Exec(`
            UPDATE resources
            SET warehouse_metal = $1,
                warehouse_he3 = $2,
                warehouse_gold = $3,
                last_warehouse_update = $4,
                updated_at = $4
            WHERE planet_id = $5
        `, newWarehouseMetal, newWarehouseHe3, newWarehouseGold, now, pr.planetID)

        if err != nil {
            log.Printf("[ResourceWorker] Failed to update warehouse for planet %s: %v", pr.planetID, err)
            continue
        }
    }

    log.Printf("[ResourceWorker] Warehouse update complete")
}

func min64(a, b int64) int64 {
    if a < b {
        return a
    }
    return b
}

// StartResourceWorker runs the warehouse update worker every 5 minutes
func StartResourceWorker() {
    // Run immediately on startup
    UpdateAllWarehouses()

    ticker := time.NewTicker(5 * time.Minute)
    go func() {
        for range ticker.C {
            UpdateAllWarehouses()
        }
    }()
    log.Println("[ResourceWorker] Started (updating warehouses every 5 minutes)")
}
```

**Integration:** Add to `cmd/server/main.go`:

```go
import "github.com/cryptomines-online/backend/internal/workers"

func main() {
    // ... existing setup ...

    // Start workers
    workers.StartBlueprintWorker()  // Existing
    workers.StartResourceWorker()   // NEW

    // ... start server ...
}
```

---

## 4. Backend Endpoint Updates

### Update: `GET /api/planets/{id}/resources`

**File:** `/backend/internal/handlers/resources.go:GetResources`

**Current (lines 28-62):**
```go
var res models.Resource
err := database.DB.QueryRow(
    `SELECT id, planet_id, metal, he3, gold, metal_per_hour, he3_per_hour,
            gold_per_hour, storage_capacity, last_collected_at, updated_at
     FROM resources WHERE planet_id = $1`, planetID,
).Scan(...)

// Calculate pending accumulated resources (not yet collected)
elapsed := time.Since(res.LastCollectedAt).Hours()
pendingMetal := int64(float64(res.MetalPerHour) * elapsed)
pendingHe3 := int64(float64(res.He3PerHour) * elapsed)
pendingGold := int64(float64(res.GoldPerHour) * elapsed)

type resourcesResponse struct {
    models.Resource
    PendingMetal int64 `json:"pending_metal"`
    PendingHe3   int64 `json:"pending_he3"`
    PendingGold  int64 `json:"pending_gold"`
}
```

**NEW (replace):**
```go
var res models.Resource
err := database.DB.QueryRow(
    `SELECT id, planet_id, metal, he3, gold,
            warehouse_metal, warehouse_he3, warehouse_gold,  -- NEW
            metal_per_hour, he3_per_hour, gold_per_hour,
            storage_capacity, last_collected_at, last_warehouse_update, updated_at
     FROM resources WHERE planet_id = $1`, planetID,
).Scan(
    &res.ID, &res.PlanetID, &res.Metal, &res.He3, &res.Gold,
    &res.WarehouseMetal, &res.WarehouseHe3, &res.WarehouseGold,  -- NEW
    &res.MetalPerHour, &res.He3PerHour, &res.GoldPerHour,
    &res.StorageCapacity, &res.LastCollectedAt, &res.LastWarehouseUpdate, &res.UpdatedAt,
)

// No pending calculation needed - warehouse values are already current
type resourcesResponse struct {
    models.Resource  // Already includes warehouse_* fields
}

resp := resourcesResponse{Resource: res}
```

**Update Model:** `/backend/internal/models/resource.go`

```go
type Resource struct {
    ID                   string    `json:"id"`
    PlanetID             string    `json:"planet_id"`
    Metal                int64     `json:"metal"`
    He3                  int64     `json:"he3"`
    Gold                 int64     `json:"gold"`
    WarehouseMetal       int64     `json:"warehouse_metal"`        // NEW
    WarehouseHe3         int64     `json:"warehouse_he3"`          // NEW
    WarehouseGold        int64     `json:"warehouse_gold"`         // NEW
    MetalPerHour         int64     `json:"metal_per_hour"`
    He3PerHour           int64     `json:"he3_per_hour"`
    GoldPerHour          int64     `json:"gold_per_hour"`
    StorageCapacity      int64     `json:"storage_capacity"`
    LastCollectedAt      time.Time `json:"last_collected_at"`
    LastWarehouseUpdate  time.Time `json:"last_warehouse_update"`  // NEW
    UpdatedAt            time.Time `json:"updated_at"`
}
```

---

### Update: `POST /api/planets/{id}/resources/collect`

**File:** `/backend/internal/handlers/resources.go:CollectResources`

**Current (lines 100-166):**
```go
// Get current resources with lock
var res models.Resource
err = tx.QueryRow(`SELECT ... FROM resources WHERE planet_id = $1 FOR UPDATE`, planetID).Scan(...)

// Calculate accumulated resources
elapsed := time.Since(res.LastCollectedAt).Hours()
collectedMetal := int64(float64(res.MetalPerHour) * elapsed)
collectedHe3 := int64(float64(res.He3PerHour) * elapsed)
collectedGold := int64(float64(res.GoldPerHour) * elapsed)

// Cap resources at storage capacity
newMetal := min64(res.Metal+collectedMetal, res.StorageCapacity)
newHe3 := min64(res.He3+collectedHe3, res.StorageCapacity)
newGold := min64(res.Gold+collectedGold, res.StorageCapacity)

actualMetal := newMetal - res.Metal
actualHe3 := newHe3 - res.He3
actualGold := newGold - res.Gold

// Update resources
err = tx.QueryRow(`
    UPDATE resources
    SET metal = $1, he3 = $2, gold = $3, last_collected_at = $4, updated_at = $4
    WHERE planet_id = $5
    RETURNING ...
`, newMetal, newHe3, newGold, now, planetID).Scan(...)
```

**NEW (replace):**
```go
// Get current resources with lock
var res models.Resource
err = tx.QueryRow(`
    SELECT id, planet_id, metal, he3, gold,
           warehouse_metal, warehouse_he3, warehouse_gold,  -- NEW
           metal_per_hour, he3_per_hour, gold_per_hour,
           storage_capacity, last_collected_at, last_warehouse_update, updated_at
    FROM resources WHERE planet_id = $1 FOR UPDATE
`, planetID).Scan(
    &res.ID, &res.PlanetID, &res.Metal, &res.He3, &res.Gold,
    &res.WarehouseMetal, &res.WarehouseHe3, &res.WarehouseGold,  -- NEW
    &res.MetalPerHour, &res.He3PerHour, &res.GoldPerHour,
    &res.StorageCapacity, &res.LastCollectedAt, &res.LastWarehouseUpdate, &res.UpdatedAt,
)

// Transfer warehouse → player balance (no elapsed time calculation)
collectedMetal := res.WarehouseMetal
collectedHe3 := res.WarehouseHe3
collectedGold := res.WarehouseGold

// Add to player balance
newMetal := res.Metal + collectedMetal
newHe3 := res.He3 + collectedHe3
newGold := res.Gold + collectedGold

// Update resources: add warehouse to balance, clear warehouse
now := time.Now()
err = tx.QueryRow(`
    UPDATE resources
    SET metal = $1, he3 = $2, gold = $3,
        warehouse_metal = 0, warehouse_he3 = 0, warehouse_gold = 0,  -- Clear warehouse
        last_collected_at = $4, updated_at = $4
    WHERE planet_id = $5
    RETURNING id, planet_id, metal, he3, gold,
              warehouse_metal, warehouse_he3, warehouse_gold,
              metal_per_hour, he3_per_hour, gold_per_hour,
              storage_capacity, last_collected_at, last_warehouse_update, updated_at
`, newMetal, newHe3, newGold, now, planetID).Scan(
    &res.ID, &res.PlanetID, &res.Metal, &res.He3, &res.Gold,
    &res.WarehouseMetal, &res.WarehouseHe3, &res.WarehouseGold,
    &res.MetalPerHour, &res.He3PerHour, &res.GoldPerHour,
    &res.StorageCapacity, &res.LastCollectedAt, &res.LastWarehouseUpdate, &res.UpdatedAt,
)
```

**Key Change:**
- **Before:** Calculate elapsed time × rate, add to balance
- **After:** Transfer warehouse → balance, clear warehouse (worker already accumulated)

---

## 5. Frontend Updates

### Update: `ResourceHUD.tsx`

**File:** `/frontend/src/components/layout/ResourceHUD.tsx`

**Current (lines 38-52):**
```tsx
<div className="hud-resource">
  <span className="hud-resource-icon metal">M</span>
  <span className="hud-resource-value">{formatNumber(resources.metal)}</span>
  <span className="hud-resource-rate">+{formatNumber(resources.metal_per_hour)}/hr</span>
</div>
```

**NEW (show warehouse separately):**
```tsx
<div className="hud-resource">
  <span className="hud-resource-icon metal">M</span>
  <div className="hud-resource-split">
    <div className="hud-resource-balance">
      {formatNumber(resources.metal)}
      {resources.warehouse_metal > 0 && (
        <span className="hud-warehouse-badge">
          +{formatNumber(resources.warehouse_metal)}
        </span>
      )}
    </div>
    <span className="hud-resource-rate">+{formatNumber(resources.metal_per_hour)}/hr</span>
  </div>
</div>
```

**Update Collect Button (lines 63-69):**
```tsx
const warehouse = resources
  ? (resources.warehouse_metal + resources.warehouse_he3 + resources.warehouse_gold)
  : 0

<button
  className="hud-collect-btn"
  onClick={handleCollect}
  disabled={collecting || warehouse <= 0}
  title={warehouse > 0 ? `Warehouse: ${formatNumber(warehouse)} total` : 'Warehouse empty'}
>
  {collecting ? 'Collecting...' : `Collect +${formatNumber(warehouse)}`}
</button>
```

**CSS Update:**
```css
.hud-resource-split {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.hud-resource-balance {
  display: flex;
  align-items: center;
  gap: 4px;
}

.hud-warehouse-badge {
  font-size: 11px;
  color: #4ade80;
  background: rgba(74, 222, 128, 0.15);
  padding: 1px 4px;
  border-radius: 3px;
}
```

### Update: `useResources.ts`

**File:** `/frontend/src/hooks/useResources.ts`

**Current (lines 27-46):**
```ts
const collect = useCallback(async () => {
  if (!planetId) return
  try {
    const result = await collectResources(planetId)
    dispatch({
      type: 'SET_RESOURCES',
      payload: {
        ...result.resources,
        pending_metal: 0,  // Clear pending
        pending_he3: 0,
        pending_gold: 0,
      },
    })
    // ...
  }
}, [planetId, dispatch])
```

**NEW (remove pending logic):**
```ts
const collect = useCallback(async () => {
  if (!planetId) return
  try {
    const result = await collectResources(planetId)
    dispatch({
      type: 'SET_RESOURCES',
      payload: result.resources,  // Warehouse already cleared by backend
    })
    // Refresh buildings to detect completed upgrades
    const buildings = await listBuildings(planetId)
    dispatch({ type: 'SET_BUILDINGS', payload: buildings })
  } catch {
    dispatch({ type: 'SET_ERROR', payload: 'Failed to collect resources' })
  }
}, [planetId, dispatch])
```

---

## 6. Integration Points

**Files to Create:**
1. **NEW:** `/supabase/migrations/YYYYMMDDHHMMSS_add_warehouse_fields.sql` - Migration
2. **NEW:** `/backend/internal/workers/resource_worker.go` - Warehouse worker

**Files to Modify:**
1. `/backend/cmd/server/main.go` - Start resource worker
2. `/backend/internal/models/resource.go` - Add warehouse fields
3. `/backend/internal/handlers/resources.go` - Update GET/POST endpoints
4. `/frontend/src/components/layout/ResourceHUD.tsx` - Show warehouse + tooltip
5. `/frontend/src/hooks/useResources.ts` - Remove pending logic
6. `/frontend/src/types/index.ts` - Add warehouse fields to Resource type

**No breaking changes to API response structure** - warehouse fields are additive.

---

## 7. QA Checklist

### Backend - Warehouse Worker

- [ ] Start server → worker runs immediately
- [ ] Verify logs: "ResourceWorker Started (updating warehouses every 5 minutes)"
- [ ] Wait 5 minutes → verify logs: "Updating warehouses for N planets"
- [ ] Check DB: `warehouse_metal/he3/gold` should increase
- [ ] Verify cap: Warehouse never exceeds `storage_capacity`
- [ ] Test with 0 production: Warehouse stays at 0

### Backend - Migration

- [ ] Run migration → 4 new columns added
- [ ] Existing resources: warehouse initialized with pending amounts
- [ ] Constraints: Cannot set warehouse to negative values
- [ ] `last_warehouse_update` defaults to now()

### Backend - Collect Endpoint

- [ ] Player has warehouse_metal=10k, metal=5k
- [ ] Call `/api/planets/{id}/resources/collect`
- [ ] Verify: `metal=15k`, `warehouse_metal=0`
- [ ] Repeat collect immediately → collected=0 (warehouse empty)
- [ ] Wait 10 minutes → warehouse fills again → collect works

### Frontend - ResourceHUD

- [ ] Resources display correctly (balance + warehouse badge)
- [ ] Warehouse badge shows "+10k" in green
- [ ] Collect button disabled when warehouse=0
- [ ] Collect button shows tooltip: "Warehouse: 30k total"
- [ ] Click collect → warehouse badge disappears, balance increases
- [ ] Auto-refresh (30s) updates warehouse values

### Edge Cases

- [ ] Storage capacity upgrade → warehouse caps at new capacity
- [ ] Building destroyed → production rate drops → warehouse accumulates slower
- [ ] Player offline 24 hours → warehouse caps at storage_capacity (not overflow)
- [ ] Tech bonus (High Yield Mining Lv5) → warehouse accumulates faster
- [ ] Multiple warehouses → uses highest storage_capacity (existing logic)

### Performance

- [ ] Worker handles 1000+ planets in <1 second
- [ ] GET /resources response time unchanged (<50ms)
- [ ] Collect endpoint response time <100ms

---

## 8. Estimated Complexity

| Task | Complexity | Hours |
|------|-----------|-------|
| Database Migration | Low | 1 |
| Resource Worker | Low-Medium | 2-3 |
| Update GET Endpoint | Low | 1 |
| Update Collect Endpoint | Low-Medium | 1-2 |
| Update Resource Model | Low | 0.5 |
| Frontend ResourceHUD | Medium | 2-3 |
| Frontend useResources | Low | 1 |
| QA Testing | Medium | 3-4 |
| **Total** | **Low-Medium** | **11-16 hours** |

**Estimated Duration:** 2-3 days (assuming 6-8 hours/day)

**Critical Path:**
1. Migration (backend-dev) - BLOCKS all other tasks
2. Resource Worker (backend-dev) - After migration
3. Update endpoints (backend-dev) - After migration
4. Frontend updates (frontend-dev) - After backend complete
5. QA testing (qa-agent) - After all complete

---

## 9. Risk Assessment

### LOW RISK
- ✅ Migration is additive (no data loss risk)
- ✅ Worker pattern already proven (blueprint_worker, research auto-complete)
- ✅ Existing collect endpoint - only logic change, not API contract

### MEDIUM RISK
- ⚠️ **Data migration** - Initializing warehouse with pending amounts
  - **Mitigation:** Test migration on dev DB first, verify warehouse = pending
- ⚠️ **Worker performance** - Updating all planets every 5 minutes
  - **Mitigation:** Query only planets with production > 0, single UPDATE per planet
  - **Scaling:** For 10k planets: 10k × 1ms = 10 seconds (acceptable)

### POTENTIAL BLOCKERS
- **NONE IDENTIFIED**

---

## 10. Success Criteria

**Module 3 is COMPLETE when:**

1. ✅ Migration adds 4 warehouse fields to `resources` table
2. ✅ Resource worker updates warehouses every 5 minutes
3. ✅ Warehouse caps at `storage_capacity` (no overflow)
4. ✅ GET /resources returns warehouse values
5. ✅ Collect endpoint transfers warehouse → player balance
6. ✅ Collect endpoint clears warehouse after transfer
7. ✅ Frontend shows warehouse badge "+X" next to resources
8. ✅ Collect button tooltip shows total warehouse amount
9. ✅ Auto-refresh (30s) updates warehouse display
10. ✅ All 15 QA test cases pass
11. ✅ No critical bugs
12. ✅ Player can collect resources offline (warehouse accumulated)

---

## 11. Dependencies

**Upstream (Required Before This Module):**
- ✅ Phase 1 complete (resources table, production rates)
- ✅ Existing collect endpoint

**Downstream (Modules That Depend On This):**
- Module 1: Tech bonuses (High Yield Mining/Chemistry/Investing affect production rates, which affect warehouse accumulation)
- All modules: Spending resources uses player balance, not warehouse

**Synergy with Module 1:**
- Tech bonuses apply to `metal_per_hour/he3_per_hour/gold_per_hour`
- Worker reads these values → warehouse accumulates faster with tech bonuses
- No additional integration needed

---

## 12. Code Snippets Reference

### Worker Update Logic

```go
// Calculate elapsed time
elapsed := now.Sub(lastWarehouseUpdate).Hours()

// Production in elapsed time
produced := int64(float64(productionPerHour) * elapsed)

// Add to warehouse, cap at storage
newWarehouse := min(currentWarehouse + produced, storageCapacity)

// Update DB
UPDATE resources SET warehouse_metal = newWarehouse, last_warehouse_update = now()
```

### Collect Logic

```go
// Before (manual calculation):
elapsed := time.Since(lastCollectedAt).Hours()
collected := int64(productionPerHour * elapsed)
newBalance := min(balance + collected, storageCapacity)

// After (transfer warehouse):
collected := warehouseMetal
newBalance := balance + collected
newWarehouse := 0
```

### Frontend Display Logic

```tsx
// Show balance with warehouse badge
{formatNumber(resources.metal)}
{resources.warehouse_metal > 0 && (
  <span className="warehouse-badge">
    +{formatNumber(resources.warehouse_metal)}
  </span>
)}
```

---

## 13. Notes

**Worker Frequency:**
- **5 minutes** chosen to balance:
  - Responsiveness (players see warehouse fill up reasonably fast)
  - Performance (10k planets × 5min = acceptable load)
  - Accuracy (5min precision is fine for hourly production rates)

**Alternative: 1 minute** would be more responsive but 5× load on DB.

**Why Not Real-Time?**
- Updating warehouse on every GET request defeats the purpose (same as current system)
- Worker approach allows true offline accumulation + reduced DB load

**Storage Capacity:**
- Warehouse caps at `storage_capacity` (same cap as player balance)
- **Future enhancement:** Separate warehouse cap vs player balance cap
- Current design: One shared cap (simpler, matches GO2 behavior)

**Tech Bonuses:**
- Module 1 bonuses apply to `metal_per_hour/he3_per_hour/gold_per_hour`
- Worker reads these values → automatically applies bonuses
- No additional code needed in worker

**Player Balance vs Warehouse:**
- **Player Balance (`metal/he3/gold`)**: Used for spending (building, ships, research)
- **Warehouse (`warehouse_*`)**: Auto-accumulates, cannot spend directly, must collect first
- **Why separate?** Prevents instant spending of "pending" resources (must consciously collect)

---

## 14. Future Enhancements (Out of Scope)

**Not Included in Module 3:**

1. **Separate warehouse capacity** - warehouse_capacity independent from storage_capacity
2. **Overflow protection alert** - Notify player when warehouse near cap
3. **Auto-collect option** - Toggle to auto-collect when warehouse full
4. **Resource boost items** - +20% production for 24h (Module 9: Inventory)
5. **Offline bonus** - +10% production when offline >8 hours
6. **Visual warehouse fill animation** - Animated progress bar showing accumulation
7. **Resource history graph** - Chart showing production over time

---

**Document Status:** FINAL
**Ready for Implementation:** YES 🚀
**Next Step:** Assign tasks to backend-dev, frontend-dev, qa-agent
