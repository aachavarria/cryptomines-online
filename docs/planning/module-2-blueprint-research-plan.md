# Module 2: Blueprint Research System - Implementation Plan

**Date:** 2026-02-07
**Architect:** Claude (Sonnet)
**Status:** READY FOR IMPLEMENTATION
**Estimated Duration:** 3-4 days

---

## Executive Summary

The Blueprint Research System allows players to upgrade hulls/modules from **tier 1 → tier 2 → tier 3** at the **Weapon Research Center**.

**Current State:** ~60% implemented
- ✅ Database schema (`blueprint_research` table exists)
- ✅ Endpoint exists: `POST /api/blueprints/{id}/research`
- ✅ Basic validation (WRC level, research slots, costs)
- ❌ Auto-complete worker (MISSING)
- ❌ Tier validation in ship design (MISSING)
- ❌ Frontend UI (MISSING)

**Work Remaining:** 40%
1. Auto-complete worker (apply research_level when timer expires)
2. Tier validation in ship design (enforce unlocked tiers)
3. Frontend: Research button + stars + timer in BlueprintPanel
4. Frontend: Filter hulls/modules by unlocked tiers in ShipDesignPanel
5. QA testing

---

## 1. Database Analysis

### ✅ Already Complete

**Tables:**
```sql
blueprint_research:
  id, player_id, player_blueprint_id, target_level,
  is_researching, research_finish_at,
  metal_cost, he3_cost, gold_cost, created_at

player_blueprints:
  id, player_id, blueprint_id, is_activated,
  research_level, -- CRITICAL FIELD (0-3)
  acquired_at

blueprints:
  id, name, blueprint_type, hull_type_id, module_type_id,
  source, research_level, description

hull_types:
  id, name, tier (1-3), ... -- tier IN hull name suffix
  -- Example: weikes_i (tier 1), weikes_ii (tier 2), weikes_iii (tier 3)

module_types:
  id, name, tier (0-3), ... -- tier IN separate column
  -- Example: "Plasma Cannon" tier 1, "Plasma Cannon" tier 2
```

**Constraints:**
- `research_level` in `player_blueprints`: 0-3 (0 = unresearched, 1-3 = researched tiers)
- `target_level` in `blueprint_research`: 2-3 (cannot research to level 1, that's base)
- Research consistency: `(is_researching = true AND research_finish_at IS NOT NULL) OR (is_researching = false AND research_finish_at IS NULL)`

### ❌ Missing

**NONE** - Database schema is complete.

---

## 2. Backend Analysis

### ✅ Already Complete

**File:** `/backend/internal/handlers/blueprints.go`

**Endpoint:** `POST /api/blueprints/{id}/research` (lines 140-270)

**Implemented:**
- Load player blueprint
- Validate: is_activated, research_level < 3
- Check WRC level (must have WRC built)
- Check research slots (1 active max per player)
- Calculate costs:
  - `baseCost = 10000 × targetLevel`
  - `metalCost = baseCost`
  - `he3Cost = baseCost × 3/4 = 7500 × targetLevel`
  - `goldCost = baseCost / 2 = 5000 × targetLevel`
- Research time: `3600 seconds × targetLevel` (1hr per level)
- Deduct resources from homeworld
- Insert into `blueprint_research` table
- Return response with finish time

**Example Costs:**
- Research to Lv2: 10k Metal, 7.5k He3, 5k Gold, 2 hours
- Research to Lv3: 20k Metal, 15k He3, 10k Gold, 3 hours

### ❌ Missing: Auto-Complete Worker

**Problem:** Research completes but `player_blueprints.research_level` doesn't increment.

**Current Flow:**
1. Player calls `/api/blueprints/{id}/research`
2. `blueprint_research` row created with `research_finish_at`
3. Timer expires... **NOTHING HAPPENS**
4. `research_level` stays at 0/1/2, never increments

**Needed:** Worker to check for completed research and apply level.

---

### ❌ Missing: Tier Validation in Ship Design

**Problem:** Player can use tier 2/3 hulls/modules without researching.

**Current Flow (`ship_designs.go`):**
- Line 126-150: Checks player owns activated hull blueprint
- **MISSING:** Check `player_blueprints.research_level >= hull.tier`

**Example Exploit:**
1. Player owns "Weikes Blueprint" (activated, research_level=0)
2. Player creates design with `weikes_iii` (tier 3) hull
3. **Should fail** but currently succeeds

**Same Issue for Modules:**
- Line 150+: Module validation exists but no tier check

---

### ❌ Missing: Blueprint Service Functions

**Needed for tier logic:**

**File:** `/backend/internal/services/blueprint_service.go` (NEW)

```go
package services

import (
    "database/sql"
    "fmt"
    "strings"
    "github.com/cryptomines-online/backend/internal/database"
)

// GetHullBaseName extracts base name from tier-suffixed hull
// "weikes_i" → "weikes"
// "weikes_ii" → "weikes"
// "weikes_iii" → "weikes"
func GetHullBaseName(hullName string) string {
    name := strings.TrimSuffix(hullName, "_iii")
    name = strings.TrimSuffix(name, "_ii")
    name = strings.TrimSuffix(name, "_i")
    return name
}

// GetHullTierSuffix returns suffix for tier
// 1 → "_i", 2 → "_ii", 3 → "_iii"
func GetHullTierSuffix(tier int) string {
    switch tier {
    case 1:
        return "_i"
    case 2:
        return "_ii"
    case 3:
        return "_iii"
    default:
        return "_i"
    }
}

// GetPlayerBlueprintResearchLevel returns research level for a hull/module blueprint
// Returns 0 if not owned, 1-3 if researched
func GetPlayerBlueprintResearchLevel(playerID string, blueprintType string, refID int) (int, error) {
    var researchLevel int
    var query string

    if blueprintType == "hull" {
        // Hull blueprints reference tier 1 hull, research level applies to all tiers
        query = `
            SELECT COALESCE(MAX(pb.research_level), 0)
            FROM player_blueprints pb
            JOIN blueprints bp ON pb.blueprint_id = bp.id
            JOIN hull_types ht ON bp.hull_type_id = ht.id
            JOIN hull_types ht2 ON (
                ht2.id = $2
                AND ht.hull_class = ht2.hull_class
                AND ht.tier = 1
            )
            WHERE pb.player_id = $1 AND pb.is_activated = true
        `
    } else {
        // Module blueprints reference all tiers of same name
        query = `
            SELECT COALESCE(MAX(pb.research_level), 0)
            FROM player_blueprints pb
            JOIN blueprints bp ON pb.blueprint_id = bp.id
            JOIN module_types mt ON bp.module_type_id = mt.id
            JOIN module_types mt2 ON (
                mt2.id = $2
                AND mt.name = mt2.name
            )
            WHERE pb.player_id = $1 AND pb.is_activated = true
        `
    }

    err := database.DB.QueryRow(query, playerID, refID).Scan(&researchLevel)
    if err != nil && err != sql.ErrNoRows {
        return 0, fmt.Errorf("failed to get research level: %w", err)
    }

    return researchLevel, nil
}

// CanUseHullTier checks if player can use a hull of given tier
func CanUseHullTier(playerID string, hullTypeID int, tier int) (bool, error) {
    researchLevel, err := GetPlayerBlueprintResearchLevel(playerID, "hull", hullTypeID)
    if err != nil {
        return false, err
    }

    // research_level 0 = can use tier 1 only
    // research_level 1 = can use tier 1 only
    // research_level 2 = can use tier 1-2
    // research_level 3 = can use tier 1-3
    return tier <= researchLevel+1 || (researchLevel == 0 && tier == 1), nil
}

// CanUseModuleTier checks if player can use a module of given tier
func CanUseModuleTier(playerID string, moduleTypeID int, tier int) (bool, error) {
    researchLevel, err := GetPlayerBlueprintResearchLevel(playerID, "module", moduleTypeID)
    if err != nil {
        return false, err
    }

    // Tier 0 modules (e.g., Orbital Shield) always available
    if tier == 0 {
        return true, nil
    }

    // Same logic as hulls
    return tier <= researchLevel+1 || (researchLevel == 0 && tier == 1), nil
}
```

**Research Level Logic:**
- `research_level = 0`: Blueprint owned but not researched → Can use tier 1 only
- `research_level = 1`: Blueprint owned (base level) → Can use tier 1 only
- `research_level = 2`: Researched to level 2 → Can use tier 1-2
- `research_level = 3`: Researched to level 3 (max) → Can use tier 1-3

---

## 3. Auto-Complete Worker Implementation

**File:** `/backend/internal/workers/blueprint_worker.go` (NEW)

```go
package workers

import (
    "log"
    "time"
    "github.com/cryptomines-online/backend/internal/database"
    "github.com/cryptomines-online/backend/internal/services"
)

// ApplyCompletedBlueprintResearch checks for finished blueprint research and applies levels
func ApplyCompletedBlueprintResearch() {
    now := time.Now()

    // Find completed research
    rows, err := database.DB.Query(`
        SELECT br.id, br.player_id, br.player_blueprint_id, br.target_level, pb.blueprint_id
        FROM blueprint_research br
        JOIN player_blueprints pb ON br.player_blueprint_id = pb.id
        WHERE br.is_researching = true AND br.research_finish_at <= $1
    `, now)
    if err != nil {
        log.Printf("[BlueprintWorker] Failed to query completed research: %v", err)
        return
    }

    type completedResearch struct {
        researchID        string
        playerID          string
        playerBlueprintID string
        targetLevel       int
        blueprintID       int
    }
    var completed []completedResearch

    for rows.Next() {
        var cr completedResearch
        rows.Scan(&cr.researchID, &cr.playerID, &cr.playerBlueprintID, &cr.targetLevel, &cr.blueprintID)
        completed = append(completed, cr)
    }
    rows.Close()

    if len(completed) == 0 {
        return
    }

    log.Printf("[BlueprintWorker] Found %d completed blueprint research to apply", len(completed))

    for _, cr := range completed {
        tx, err := database.DB.Begin()
        if err != nil {
            log.Printf("[BlueprintWorker] Failed to begin transaction: %v", err)
            continue
        }

        // Update player_blueprints.research_level
        _, err = tx.Exec(`
            UPDATE player_blueprints
            SET research_level = $1
            WHERE id = $2
        `, cr.targetLevel, cr.playerBlueprintID)
        if err != nil {
            tx.Rollback()
            log.Printf("[BlueprintWorker] Failed to update research level: %v", err)
            continue
        }

        // Mark research as complete
        _, err = tx.Exec(`
            UPDATE blueprint_research
            SET is_researching = false, research_finish_at = NULL
            WHERE id = $1
        `, cr.researchID)
        if err != nil {
            tx.Rollback()
            log.Printf("[BlueprintWorker] Failed to mark research complete: %v", err)
            continue
        }

        if err := tx.Commit(); err != nil {
            log.Printf("[BlueprintWorker] Failed to commit: %v", err)
            continue
        }

        // Update quest progress (optional, if blueprint research quests exist)
        // services.UpdateQuestProgress(cr.playerID, "research_blueprint", fmt.Sprintf("%d", cr.blueprintID), 1)

        log.Printf("[BlueprintWorker] Applied blueprint research: player=%s, target_level=%d", cr.playerID, cr.targetLevel)
    }
}

// StartBlueprintWorker runs the blueprint research worker every 30 seconds
func StartBlueprintWorker() {
    ticker := time.NewTicker(30 * time.Second)
    go func() {
        for range ticker.C {
            ApplyCompletedBlueprintResearch()
        }
    }()
    log.Println("[BlueprintWorker] Started (checking every 30s)")
}
```

**Integration:** Add to `cmd/server/main.go`:

```go
import "github.com/cryptomines-online/backend/internal/workers"

func main() {
    // ... existing setup ...

    // Start workers
    workers.StartBlueprintWorker()

    // ... start server ...
}
```

---

## 4. Tier Validation Integration

### Location 1: Ship Design Creation

**File:** `/backend/internal/handlers/ship_designs.go:CreateShipDesign`

**Insert after line 150 (hull blueprint check):**

```go
// Check hull tier
canUse, err := services.CanUseHullTier(playerID, hull.ID, hull.Tier)
if err != nil {
    log.Printf("Failed to check hull tier: %v", err)
    http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
    return
}
if !canUse {
    http.Error(w, `{"error":"hull tier not unlocked - research required"}`, http.StatusConflict)
    return
}
```

### Location 2: Module Validation

**File:** `/backend/internal/handlers/ship_designs.go:CreateShipDesign`

**Current module loop (line ~170+):**

Find the loop where modules are validated. Add after blueprint check:

```go
// Check module tier
canUseModule, err := services.CanUseModuleTier(playerID, mt.ID, mt.Tier)
if err != nil {
    log.Printf("Failed to check module tier: %v", err)
    http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
    return
}
if !canUseModule {
    http.Error(w, `{"error":"module tier not unlocked - research required"}`, http.StatusConflict)
    return
}
```

### Location 3: Ship Design Update

**Same validations in `UpdateShipDesign()` function.**

---

## 5. Frontend Implementation

### Task 1: BlueprintPanel - Research Button + Stars

**File:** `/frontend/src/components/panels/BlueprintPanel.tsx`

**Current State (lines 59-100):**
- Shows blueprint cards with Activate button
- No research level stars
- No research button

**Changes:**

**1. Add research level stars display:**

```tsx
// After line 85 (inside bp-card):
{playerBp && playerBp.research_level > 0 && (
  <div className="bp-research-stars">
    {[1, 2, 3].map(star => (
      <span key={star} className={star <= playerBp.research_level ? 'star-filled' : 'star-empty'}>
        ★
      </span>
    ))}
  </div>
)}
```

**2. Add research button:**

```tsx
// Replace Activate button section (lines 84-96) with:
{activated ? (
  <>
    <span className="bp-status-active">Activated</span>
    {playerBp && playerBp.research_level < 3 && (
      <button
        className="p2-btn p2-btn-primary p2-btn-sm"
        onClick={() => handleStartResearch(bp.id, playerBp.research_level + 1)}
      >
        Research Lv{playerBp.research_level + 1}
      </button>
    )}
    {playerBp && playerBp.research_level === 3 && (
      <span className="bp-status-maxed">Max Level</span>
    )}
  </>
) : owned ? (
  <button
    className="p2-btn p2-btn-success p2-btn-sm"
    onClick={() => handleActivate(bp.id)}
    disabled={activating === bp.id}
  >
    {activating === bp.id ? '...' : 'Activate'}
  </button>
) : (
  <span className="bp-status-locked">Not Owned</span>
)}
```

**3. Add handleStartResearch function:**

```tsx
// After line 30:
async function handleStartResearch(bpId: number, targetLevel: number) {
  setActivating(bpId) // Reuse activating state
  try {
    await startBlueprintResearch(bpId)
    // Refresh blueprints after research starts
    // (hook should handle this)
  } catch (error: any) {
    alert(error.message || 'Failed to start research')
  } finally {
    setActivating(null)
  }
}
```

**4. Add API function:**

**File:** `/frontend/src/services/api.ts`

```ts
export async function startBlueprintResearch(blueprintId: number): Promise<any> {
  const response = await fetch(`/api/blueprints/${blueprintId}/research`, {
    method: 'POST',
    headers: getHeaders(),
  })
  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Failed to start blueprint research')
  }
  return response.json()
}
```

**5. Update useBlueprints hook:**

**File:** `/frontend/src/hooks/useBlueprints.ts`

Add `research` state and polling:

```ts
const [activeResearch, setActiveResearch] = useState<any | null>(null)

// Poll for active research
useEffect(() => {
  async function fetchActiveResearch() {
    try {
      const response = await fetch('/api/blueprint-research/active', {
        headers: getHeaders(),
      })
      if (response.ok) {
        const data = await response.json()
        setActiveResearch(data)
      }
    } catch {
      // Ignore
    }
  }

  fetchActiveResearch()
  const interval = setInterval(fetchActiveResearch, 5000)
  return () => clearInterval(interval)
}, [])
```

**6. CSS for stars:**

```css
.bp-research-stars {
  display: flex;
  gap: 2px;
  margin-top: 4px;
}

.star-filled {
  color: #fbbf24;
  font-size: 16px;
}

.star-empty {
  color: #4b5563;
  font-size: 16px;
}
```

---

### Task 2: ShipDesignPanel - Filter by Unlocked Tiers

**File:** `/frontend/src/components/panels/ShipDesignPanel.tsx`

**Current State:**
- Shows all hulls/modules regardless of tier
- No visual indicator for locked tiers

**Changes:**

**1. Add tier filter logic (inside DesignEditor component):**

```tsx
// After line 100:
const availableHulls = useMemo(() => {
  return hullTypes.filter(hull => {
    if (!hasHullBlueprint(hull.id)) return false

    // Get player blueprint research level
    const playerBp = myBlueprints.find(pb =>
      pb.hull_type_id === hull.id && pb.is_activated
    )

    if (!playerBp) return false

    // Tier 1 always available if blueprint owned
    if (hull.tier === 1) return true

    // Tier 2 requires research_level >= 2
    if (hull.tier === 2) return playerBp.research_level >= 2

    // Tier 3 requires research_level >= 3
    if (hull.tier === 3) return playerBp.research_level >= 3

    return false
  })
}, [hullTypes, hasHullBlueprint, myBlueprints])

const availableModules = useMemo(() => {
  return moduleTypes.filter(mod => {
    if (!hasModuleBlueprint(mod.id)) return false

    // Tier 0 always available (e.g., Orbital Shield)
    if (mod.tier === 0) return true

    const playerBp = myBlueprints.find(pb =>
      pb.module_type_id === mod.id && pb.is_activated
    )

    if (!playerBp) return false

    // Same logic as hulls
    if (mod.tier === 1) return true
    if (mod.tier === 2) return playerBp.research_level >= 2
    if (mod.tier === 3) return playerBp.research_level >= 3

    return false
  })
}, [moduleTypes, hasModuleBlueprint, myBlueprints])
```

**2. Update dropdown to use filtered lists:**

```tsx
// Replace hullTypes with availableHulls in hull dropdown
// Replace moduleTypes with availableModules in module lists
```

**3. Add visual indicator for locked tiers:**

```tsx
// In hull dropdown:
{hullTypes.map(hull => {
  const isAvailable = availableHulls.includes(hull)
  return (
    <option
      key={hull.id}
      value={hull.id}
      disabled={!isAvailable}
      className={!isAvailable ? 'locked-tier' : ''}
    >
      {hull.display_name} {hull.tier > 1 ? `(Tier ${hull.tier})` : ''}
      {!isAvailable && ' 🔒 Research Required'}
    </option>
  )
})}
```

---

### Task 3: Active Research Timer Display (Optional)

**Create:** `/frontend/src/components/panels/BlueprintResearchPanel.tsx` (NEW)

```tsx
import { useEffect, useState } from 'react'
import { formatDuration } from '../../hooks/useCountdown'

export default function BlueprintResearchPanel() {
  const [activeResearch, setActiveResearch] = useState<any | null>(null)

  useEffect(() => {
    async function fetch() {
      try {
        const response = await fetch('/api/blueprint-research/active', {
          headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
        })
        if (response.ok) {
          const data = await response.json()
          setActiveResearch(data)
        }
      } catch {
        // Ignore
      }
    }

    fetch()
    const interval = setInterval(fetch, 2000)
    return () => clearInterval(interval)
  }, [])

  if (!activeResearch) {
    return (
      <div className="research-status">
        <span>No active blueprint research</span>
      </div>
    )
  }

  const finishTime = new Date(activeResearch.research_finish_at).getTime()
  const now = Date.now()
  const remaining = Math.max(0, finishTime - now)

  return (
    <div className="research-status active">
      <div className="research-icon">🔬</div>
      <div>
        <div className="research-name">{activeResearch.blueprint_name} → Lv{activeResearch.target_level}</div>
        <div className="research-timer">
          {remaining > 0 ? formatDuration(remaining) : 'Completing...'}
        </div>
      </div>
    </div>
  )
}
```

**Add endpoint:** `GET /api/blueprint-research/active`

```go
// File: /backend/internal/handlers/blueprints.go

func GetActiveBlueprintResearch(w http.ResponseWriter, r *http.Request) {
    playerID := middleware.GetPlayerID(r)

    rows, err := database.DB.Query(`
        SELECT br.id, br.target_level, br.research_finish_at, bp.name
        FROM blueprint_research br
        JOIN player_blueprints pb ON br.player_blueprint_id = pb.id
        JOIN blueprints bp ON pb.blueprint_id = bp.id
        WHERE br.player_id = $1 AND br.is_researching = true
        LIMIT 1
    `, playerID)
    if err != nil {
        http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    if !rows.Next() {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("null"))
        return
    }

    var data struct {
        ID               string `json:"id"`
        TargetLevel      int    `json:"target_level"`
        ResearchFinishAt string `json:"research_finish_at"`
        BlueprintName    string `json:"blueprint_name"`
    }

    rows.Scan(&data.ID, &data.TargetLevel, &data.ResearchFinishAt, &data.BlueprintName)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}
```

---

## 6. Integration Points

**Files to Create:**
1. **NEW:** `/backend/internal/services/blueprint_service.go` - Tier logic functions
2. **NEW:** `/backend/internal/workers/blueprint_worker.go` - Auto-complete worker
3. **NEW (optional):** `/frontend/src/components/panels/BlueprintResearchPanel.tsx` - Active research display

**Files to Modify:**
1. `/backend/cmd/server/main.go` - Start blueprint worker
2. `/backend/internal/handlers/ship_designs.go` - Add tier validation (2 locations)
3. `/backend/internal/handlers/blueprints.go` - Add active research endpoint (optional)
4. `/frontend/src/components/panels/BlueprintPanel.tsx` - Add research button + stars
5. `/frontend/src/components/panels/ShipDesignPanel.tsx` - Filter by unlocked tiers
6. `/frontend/src/hooks/useBlueprints.ts` - Add active research polling
7. `/frontend/src/services/api.ts` - Add startBlueprintResearch()

**No database changes required.**

---

## 7. QA Checklist

### Backend - Auto-Complete Worker

- [ ] Start research to Lv2 with 10-second timer (modify temporarily)
- [ ] Wait 10 seconds
- [ ] Verify `player_blueprints.research_level` increments to 2
- [ ] Verify `blueprint_research.is_researching` = false
- [ ] Repeat for Lv3

### Backend - Tier Validation

**Hull Tier Validation:**
- [ ] Own "Weikes Blueprint", research_level=0
- [ ] Try to create design with `weikes_ii` (tier 2) → expect error "tier not unlocked"
- [ ] Research blueprint to Lv2
- [ ] Create design with `weikes_ii` → success
- [ ] Try `weikes_iii` (tier 3) → expect error
- [ ] Research to Lv3
- [ ] Create design with `weikes_iii` → success

**Module Tier Validation:**
- [ ] Own "Plasma Cannon" blueprint, research_level=0
- [ ] Try to add tier 2 Plasma Cannon to design → expect error
- [ ] Research to Lv2
- [ ] Add tier 2 Plasma Cannon → success
- [ ] Try tier 3 → expect error
- [ ] Research to Lv3
- [ ] Add tier 3 → success

### Frontend - BlueprintPanel

- [ ] Activate blueprint → research button appears
- [ ] Research_level=0 → 3 empty stars (☆☆☆)
- [ ] Click "Research Lv1" → costs deducted, timer starts
- [ ] After completion → 1 filled star (★☆☆)
- [ ] Click "Research Lv2" → 2 filled stars (★★☆)
- [ ] Click "Research Lv3" → 3 filled stars (★★★)
- [ ] At research_level=3 → button says "Max Level"

### Frontend - ShipDesignPanel

- [ ] Hull dropdown shows only unlocked tiers
- [ ] Tier 2/3 hulls appear after research
- [ ] Locked hulls show "🔒 Research Required"
- [ ] Same for module lists (filter by tier)

### Edge Cases

- [ ] Try to start research while already researching → expect error "research slot in use"
- [ ] Try to research without WRC → expect error "weapon research center required"
- [ ] Try to research unactivated blueprint → expect error "must be activated first"
- [ ] Research to Lv2, then to Lv3 → both should work sequentially
- [ ] Insufficient resources → expect error "insufficient resources"

---

## 8. Estimated Complexity

| Task | Complexity | Hours |
|------|-----------|-------|
| Blueprint Service Functions | Low | 2-3 |
| Auto-Complete Worker | Low-Medium | 2-3 |
| Tier Validation (ship_designs.go) | Medium | 2-3 |
| BlueprintPanel UI (stars + button) | Medium | 3-4 |
| ShipDesignPanel Filtering | Medium | 3-4 |
| Active Research Timer (optional) | Low | 1-2 |
| QA Testing | Medium | 4-6 |
| **Total** | **Medium** | **17-25 hours** |

**Estimated Duration:** 3-4 days (assuming 6-8 hours/day)

**Critical Path:**
1. blueprint_service.go (backend-dev) - BLOCKS tier validation
2. blueprint_worker.go (backend-dev) - Can run in parallel
3. Tier validation integration (backend-dev) - DEPENDS ON #1
4. BlueprintPanel UI (frontend-dev) - Can run in parallel
5. ShipDesignPanel filtering (frontend-dev) - DEPENDS ON backend tier logic
6. QA testing (qa-agent) - After all backend + frontend complete

---

## 9. Risk Assessment

### LOW RISK
- ✅ Database schema complete
- ✅ Research endpoint functional (60% done)
- ✅ Simple tier logic (suffix parsing for hulls, column for modules)

### MEDIUM RISK
- ⚠️ **Tier validation complexity** - Hulls use `_i/_ii/_iii` suffix, modules use tier column
  - **Mitigation:** blueprint_service.go centralizes logic, well-tested
- ⚠️ **Frontend filtering** - Need to pass `myBlueprints` with research_level to ShipDesignPanel
  - **Mitigation:** useBlueprints hook already fetches this data

### POTENTIAL BLOCKERS
- **NONE IDENTIFIED**

---

## 10. Success Criteria

**Module 2 is COMPLETE when:**

1. ✅ Auto-complete worker applies research_level when timer expires
2. ✅ Ship design creation validates hull tier (blocks tier 2/3 without research)
3. ✅ Ship design creation validates module tier
4. ✅ BlueprintPanel shows research level stars (★★★)
5. ✅ BlueprintPanel shows "Research Lv2/Lv3" button for activated blueprints
6. ✅ ShipDesignPanel filters hulls/modules by unlocked tiers
7. ✅ Locked tiers show "🔒 Research Required" in UI
8. ✅ Active research timer displays (optional enhancement)
9. ✅ All 20 QA test cases pass
10. ✅ No critical bugs

---

## 11. Dependencies

**Upstream (Required Before This Module):**
- ✅ Phase 2 complete (blueprints, ship designs, WRC building)
- ✅ Endpoint exists: `POST /api/blueprints/{id}/research`

**Downstream (Modules That Depend On This):**
- Module 5: Combat System (tier 2/3 hulls/modules in combat)
- Module 9: Inventory System (blueprint items → activate → research flow)

**Synergy with Module 1:**
- Same pattern: research timer → auto-complete worker → apply effect
- Can reuse worker pattern from tech research

---

## 12. Code Snippets Reference

### Hull Tier Pattern (from seed data)

```sql
-- Tier 1 (base)
weikes_i, air_wanderer_i, valkyrie_i, ...

-- Tier 2 (researched to Lv2)
weikes_ii, air_wanderer_ii, valkyrie_ii, ...

-- Tier 3 (researched to Lv3)
weikes_iii, air_wanderer_iii, valkyrie_iii, ...
```

**Blueprint Naming:**
- Blueprint name: "Weikes Blueprint" (no tier suffix)
- Applies to: `weikes_i`, `weikes_ii`, `weikes_iii`

### Module Tier Pattern

```sql
-- Plasma Cannon (name = "plasma_cannon", tier varies)
('plasma_cannon', 'Plasma Cannon I', 'ballistic', 1, ...) -- tier 1
('plasma_cannon', 'Plasma Cannon II', 'ballistic', 2, ...) -- tier 2
('plasma_cannon', 'Plasma Cannon III', 'ballistic', 3, ...) -- tier 3
```

**Blueprint Naming:**
- Blueprint name: "Plasma Cannon Blueprint"
- Applies to: ALL tiers of `plasma_cannon`

### Research Costs (from existing code)

```go
targetLevel := research_level + 1
baseCost := 10000 * targetLevel

metalCost := baseCost                 // 10k (Lv1), 20k (Lv2), 30k (Lv3)
he3Cost := baseCost * 3 / 4           // 7.5k (Lv1), 15k (Lv2), 22.5k (Lv3)
goldCost := baseCost / 2              // 5k (Lv1), 10k (Lv2), 15k (Lv3)
researchTime := 3600 * targetLevel    // 1hr (Lv1), 2hr (Lv2), 3hr (Lv3)
```

**Note:** Current code treats research_level=0 as "unresearched" and increments from there. The target_level=2 means "researching to level 2" (from 0 or 1).

**Clarification Needed:**
- Does `research_level=0` mean "not researched" or "tier 1 available"?
- Current implementation: research_level=0 → can use tier 1, research_level=2 → can use tier 2

**Recommended:**
- `research_level=0`: Blueprint owned, tier 1 available
- `research_level=1`: (not used, or same as 0)
- `research_level=2`: Researched to Lv2, tier 2 unlocked
- `research_level=3`: Researched to Lv3, tier 3 unlocked

---

## 13. Notes

**Weapon Research Center:**
- Controls research slots (currently 1 slot max)
- Does NOT reduce research time (unlike Tech Center for tech research)
- **Future enhancement:** WRC level could reduce time/cost

**Blueprint Types:**
- Hull blueprints: 1 blueprint per hull line (e.g., "Weikes Blueprint" covers weikes_i/ii/iii)
- Module blueprints: 1 blueprint per module line (e.g., "Plasma Cannon Blueprint" covers all 3 tiers)

**Research Progression:**
- Must research sequentially: Lv1 → Lv2 → Lv3
- Cannot skip levels
- Max level is 3 (all 3 tiers unlocked)

**Inventory Flow (Module 9):**
- Blueprint drops as item → Inventory
- Player "uses" item → Unlocks in player_blueprints (is_activated=false)
- Player "activates" blueprint → is_activated=true, research_level=0
- Player researches → research_level increments to 2, then 3

**Current Implementation Gap:**
- Blueprints currently bypass inventory (quest rewards directly activate)
- Module 9 will fix this flow

---

## 14. Future Enhancements (Out of Scope)

**Not Included in Module 2:**

1. **Multiple WRC research slots** - Currently 1 slot max, could scale with WRC level
2. **WRC time reduction bonus** - Tech Center reduces tech research time, WRC could do same
3. **Cancel/Speedup blueprint research** - No cancel/speedup endpoints yet
4. **Research queue** - Queue next blueprint after current completes
5. **Blueprint respec** - Reset research level (refund resources?)
6. **Mass research** - Research multiple blueprints simultaneously
7. **Research cost balancing** - Current costs may need adjustment after playtesting

---

**Document Status:** FINAL
**Ready for Implementation:** YES 🚀
**Next Step:** Assign tasks to backend-dev, frontend-dev, qa-agent
