# Module 1: Research System - Implementation Plan

**Date:** 2026-02-07
**Architect:** Claude (Sonnet)
**Status:** READY FOR IMPLEMENTATION
**Estimated Duration:** 4-5 days

---

## Executive Summary

The Research System is **90% implemented** in Phase 1. This plan covers the remaining 10%:

1. **Tech Effects Application** - Apply bonuses when research completes
2. **Frontend Integration** - Wire up existing ResearchPanel
3. **QA Testing** - Verify all 111 techs work correctly

**Good News:** Database schema, migrations, seed data (111 techs), backend handlers (5 endpoints), frontend hooks, and ResearchPanel already exist and are functional.

---

## 1. Database Analysis

### ✅ Already Complete

**Tables:**
- `tech_types` - 111 techs seeded across 7 trees
- `technologies` - Player research progress tracking

**Migration:** `/supabase/migrations/20260206040000_research.sql`

**Schema:**
```sql
tech_types:
  id, name, display_name, tree, prerequisites_json,
  base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier,
  base_time_seconds, time_multiplier, max_level,
  effects_json, description

technologies:
  id, player_id, tech_type, level, is_researching,
  research_finish_at, created_at, updated_at
```

**Constraints:**
- UNIQUE (player_id, tech_type)
- CHECK is_researching consistency
- 7 tree types validated

### ❌ Missing

**NONE** - Database is complete.

---

## 2. Backend Analysis

### ✅ Already Complete

**Files:**
- `/backend/internal/handlers/research.go` (692 lines)
- `/backend/internal/models/research.go` (40 lines)

**Endpoints (5 total):**
1. `GET /api/research` - List all trees with player progress ✅
2. `GET /api/research/trees/{tree}` - Get single tree ✅
3. `POST /api/research/start` - Start research (validate prereqs, resources, Tech Center level) ✅
4. `POST /api/research/cancel` - Cancel active research ✅
5. `POST /api/research/speedup` - Spend vouchers (3 per 30min) ✅

**Auto-complete Worker:**
- `applyCompletedResearch()` function exists ✅
- Called on every research endpoint ✅
- Updates quest progress via `services.UpdateQuestProgress()` ✅

**Tech Center Bonus:**
- 3% time reduction per Tech Center level applied ✅

**Cost/Time Calculations:**
- `calcLevelCost()` - baseCost × multiplier^(level-1) ✅
- `calcLevelTime()` - baseTime × multiplier^(level-1) ✅

### ❌ Missing: Tech Effects Application

**Problem:** Research completes but bonuses don't apply.

**Current Code (research.go:544-558):**
```go
_, err = database.DB.Exec(`
    UPDATE technologies
    SET level = level + 1, is_researching = false, research_finish_at = NULL, updated_at = now()
    WHERE player_id = $1 AND is_researching = true AND research_finish_at <= $2
`, playerID, now)
```

**Issue:** No logic to parse `effects_json` and apply bonuses.

**Effects Types (from seed data):**

**Production Bonuses:**
- `type: "metal_output"` - +1-10% Metal production
- `type: "he3_output"` - +1-10% He3 production
- `type: "gold_output"` - +1-10% Gold production
- `type: "warehouse_capacity"` - +50k-350k storage

**Construction Bonuses:**
- `type: "construction_slots"` - +1 building slot
- `type: "build_speed"` - +1-15% building speed
- `type: "build_cost_reduction"` - -1-15% building costs
- `type: "ship_production_slots"` - +1 shipbuilding slot
- `type: "ship_build_speed"` - +1-15% ship build speed
- `type: "ship_build_cost_reduction"` - -1-15% ship resource costs
- `type: "ship_repair_percent"` - +1-10% repair percentage

**Module Unlocks:**
- Implied by tech prerequisites (e.g., Ballistics Lv3 required for module X)

**Combat Bonuses:**
- Applied via tech level lookup during combat resolution
- Examples: `ballistic_damage`, `directional_accuracy`, `missile_hit_rate`, etc.

---

## 3. Tech Effects Implementation Strategy

### Option A: Cached Bonuses (RECOMMENDED)

**Approach:** Calculate total bonuses on-demand when needed.

**Pros:**
- No database changes
- Simple implementation
- Easy to debug
- No data sync issues

**Cons:**
- Recalculate on every query (minimal cost ~20ms)

**Implementation:**

**File:** `/backend/internal/services/tech_effects.go` (NEW)

```go
package services

import (
    "database/sql"
    "encoding/json"
    "log"
    "github.com/cryptomines-online/backend/internal/database"
)

type TechEffect struct {
    Type     string  `json:"type"`
    PerLevel float64 `json:"per_level,omitempty"`
    Flat     int     `json:"flat,omitempty"`
    Unit     string  `json:"unit,omitempty"`
}

type TechBonus struct {
    MetalOutput           float64 // Percentage (0.10 = 10%)
    He3Output             float64
    GoldOutput            float64
    WarehouseCapacity     int64
    ConstructionSlots     int
    BuildSpeed            float64
    BuildCostReduction    float64
    ShipProductionSlots   int
    ShipBuildSpeed        float64
    ShipBuildCostReduction float64
    ShipRepairPercent     float64
}

// GetPlayerTechBonuses calculates all active tech bonuses for a player
func GetPlayerTechBonuses(playerID string) TechBonus {
    var bonus TechBonus

    rows, err := database.DB.Query(`
        SELECT t.level, tt.effects_json
        FROM technologies t
        JOIN tech_types tt ON t.tech_type = tt.id
        WHERE t.player_id = $1 AND t.level > 0
    `, playerID)
    if err != nil {
        log.Printf("Failed to query tech bonuses: %v", err)
        return bonus
    }
    defer rows.Close()

    for rows.Next() {
        var level int
        var effectsJSON []byte
        if err := rows.Scan(&level, &effectsJSON); err != nil {
            continue
        }

        var effect TechEffect
        if err := json.Unmarshal(effectsJSON, &effect); err != nil {
            continue
        }

        // Apply effect based on type
        switch effect.Type {
        case "metal_output":
            bonus.MetalOutput += effect.PerLevel * float64(level) / 100.0
        case "he3_output":
            bonus.He3Output += effect.PerLevel * float64(level) / 100.0
        case "gold_output":
            bonus.GoldOutput += effect.PerLevel * float64(level) / 100.0
        case "warehouse_capacity":
            bonus.WarehouseCapacity += int64(effect.PerLevel) * int64(level)
        case "construction_slots":
            bonus.ConstructionSlots += effect.PerLevel * level
        case "build_speed":
            bonus.BuildSpeed += effect.PerLevel * float64(level) / 100.0
        case "build_cost_reduction":
            bonus.BuildCostReduction += effect.PerLevel * float64(level) / 100.0
        case "ship_production_slots":
            bonus.ShipProductionSlots += effect.Flat
        case "ship_build_speed":
            bonus.ShipBuildSpeed += effect.PerLevel * float64(level) / 100.0
        case "ship_build_cost_reduction":
            bonus.ShipBuildCostReduction += effect.PerLevel * float64(level) / 100.0
        case "ship_repair_percent":
            bonus.ShipRepairPercent += effect.PerLevel * float64(level) / 100.0
        }
    }

    return bonus
}

// GetTechLevel returns the current level of a tech for a player
func GetTechLevel(playerID, techName string) int {
    var level int
    err := database.DB.QueryRow(`
        SELECT COALESCE(t.level, 0)
        FROM tech_types tt
        LEFT JOIN technologies t ON t.tech_type = tt.id AND t.player_id = $1
        WHERE tt.name = $2
    `, playerID, techName).Scan(&level)
    if err != nil {
        return 0
    }
    return level
}
```

**Integration Points:**

**1. Resource Production** (`handlers/buildings.go:GetResources`)

```go
// BEFORE (line ~50):
var res models.Resources
err := database.DB.QueryRow(`SELECT ...`, planetID).Scan(...)

// AFTER:
var res models.Resources
err := database.DB.QueryRow(`SELECT ...`, planetID).Scan(...)

// Apply tech bonuses
bonus := services.GetPlayerTechBonuses(playerID)
res.MetalPerHour = int64(float64(res.MetalPerHour) * (1.0 + bonus.MetalOutput))
res.He3PerHour = int64(float64(res.He3PerHour) * (1.0 + bonus.He3Output))
res.GoldPerHour = int64(float64(res.GoldPerHour) * (1.0 + bonus.GoldOutput))
res.StorageCapacity += bonus.WarehouseCapacity
```

**2. Building Construction** (`handlers/buildings.go:ConstructBuilding`)

```go
// BEFORE (line ~150):
effectiveTime := calcConstructionTime(bt.BuildTimeSeconds, targetLevel)

// AFTER:
baseTime := calcConstructionTime(bt.BuildTimeSeconds, targetLevel)
bonus := services.GetPlayerTechBonuses(playerID)
effectiveTime := int(float64(baseTime) * (1.0 - bonus.BuildSpeed))

// BEFORE (line ~180):
metalCost := calcLevelCost(bt.BaseCostMetal, bt.CostMultiplier, targetLevel)
he3Cost := calcLevelCost(bt.BaseCostHe3, bt.CostMultiplier, targetLevel)
goldCost := calcLevelCost(bt.BaseCostGold, bt.CostMultiplier, targetLevel)

// AFTER:
baseMetal := calcLevelCost(bt.BaseCostMetal, bt.CostMultiplier, targetLevel)
baseHe3 := calcLevelCost(bt.BaseCostHe3, bt.CostMultiplier, targetLevel)
baseGold := calcLevelCost(bt.BaseCostGold, bt.CostMultiplier, targetLevel)
bonus := services.GetPlayerTechBonuses(playerID)
metalCost := int64(float64(baseMetal) * (1.0 - bonus.BuildCostReduction))
he3Cost := int64(float64(baseHe3) * (1.0 - bonus.BuildCostReduction))
goldCost := int64(float64(baseGold) * (1.0 - bonus.BuildCostReduction))
```

**3. Ship Production** (`handlers/ship_factory.go:BuildShips`)

```go
// BEFORE (line ~120):
baseTime := calcShipBuildTime(...)

// AFTER:
baseTime := calcShipBuildTime(...)
bonus := services.GetPlayerTechBonuses(playerID)
effectiveTime := int(float64(baseTime) * (1.0 - bonus.ShipBuildSpeed))

// Resource costs (similar pattern)
baseMetal := design.MetalCost * quantity
baseHe3 := design.He3Cost * quantity
baseGold := design.GoldCost * quantity
bonus := services.GetPlayerTechBonuses(playerID)
metalCost := int64(float64(baseMetal) * (1.0 - bonus.ShipBuildCostReduction))
he3Cost := int64(float64(baseHe3) * (1.0 - bonus.ShipBuildCostReduction))
goldCost := int64(float64(baseGold) * (1.0 - bonus.ShipBuildCostReduction))
```

**4. Construction/Production Slots**

**File:** `/backend/internal/handlers/buildings.go:ConstructBuilding`

```go
// BEFORE (line ~100):
var activeCount int
err = tx.QueryRow(`SELECT COUNT(*) FROM buildings WHERE planet_id = $1 AND is_constructing = true`, planetID).Scan(&activeCount)
maxSlots := 2 // Base slots
if activeCount >= maxSlots {
    return // error: all slots busy
}

// AFTER:
var activeCount int
err = tx.QueryRow(`SELECT COUNT(*) FROM buildings WHERE planet_id = $1 AND is_constructing = true`, planetID).Scan(&activeCount)
bonus := services.GetPlayerTechBonuses(playerID)
maxSlots := 2 + bonus.ConstructionSlots // Base 2 + tech bonus
if activeCount >= maxSlots {
    return // error: all slots busy
}
```

**File:** `/backend/internal/handlers/ship_factory.go:BuildShips`

```go
// BEFORE (line ~80):
var activeCount int
err = tx.QueryRow(`SELECT COUNT(*) FROM ship_factory_slots WHERE planet_id = $1 AND is_active = true`, planetID).Scan(&activeCount)
factoryLevel := getShipFactoryLevel(tx, planetID)
maxSlots := min(factoryLevel / 6 + 1, 5) // Max 5 slots at Lv24
if activeCount >= maxSlots {
    return // error: all slots busy
}

// AFTER:
var activeCount int
err = tx.QueryRow(`SELECT COUNT(*) FROM ship_factory_slots WHERE planet_id = $1 AND is_active = true`, planetID).Scan(&activeCount)
factoryLevel := getShipFactoryLevel(tx, playerID)
baseSlots := min(factoryLevel / 6 + 1, 4) // Max 4 base slots
bonus := services.GetPlayerTechBonuses(playerID)
maxSlots := baseSlots + bonus.ShipProductionSlots // +1 from Sync Shipbuilding tech
if activeCount >= maxSlots {
    return // error: all slots busy
}
```

**5. Combat Bonuses** (Module 5: Combat System)

Combat bonuses are applied during combat resolution by looking up tech levels:

```go
// Example: Ballistics damage bonus
ballisticsLevel := services.GetTechLevel(playerID, "ballistics_base")
damageBonus := 1.0 + (float64(ballisticsLevel) * 0.05) // +5% per level
finalDamage := baseDamage * damageBonus
```

**Note:** Combat integration deferred to Module 5 implementation.

---

### Option B: Materialized Table (NOT RECOMMENDED)

**Approach:** Store computed bonuses in `player_tech_bonuses` table.

**Pros:**
- Faster queries (no recalculation)

**Cons:**
- Requires migration
- Complex sync logic
- Potential data inconsistency
- Overkill for ~111 techs

**Verdict:** Use Option A (cached). Performance difference negligible (~20ms vs ~2ms).

---

## 4. Frontend Analysis

### ✅ Already Complete

**Components:**
- `ResearchPanel.tsx` (300+ lines) - Full UI with tabs, tech tree, timers ✅
- `useResearch.ts` hook - API integration ✅

**Features:**
- 7 tree tabs (Logistics, Ballistics, Directional, Missile, Ship-Based, Ship Defense, Planetary) ✅
- Tech display with current level, max level, costs ✅
- Prerequisites display ✅
- Start/cancel/speedup buttons ✅
- Active research countdown timer ✅
- Confirmation modals ✅
- Toast notifications ✅

**API Integration:**
- `getResearch()` ✅
- `getResearchTree(tree)` ✅
- `startResearch(techTypeId)` ✅
- `cancelResearch(techTypeId)` ✅
- `speedupResearch(techTypeId, minutes)` ✅

### ❌ Missing

**1. Tech Effects Tooltips**

**Current:** Tooltips show description only.

**Needed:** Show active bonuses (e.g., "Current: +5% Metal Output (Level 5)")

**File:** `ResearchPanel.tsx:TechCard`

```tsx
// BEFORE:
<div className="tooltip">
  <p>{tech.description}</p>
  <p>Cost: {formatNumber(tech.cost_next_level.gold)} Gold</p>
</div>

// AFTER:
<div className="tooltip">
  <p>{tech.description}</p>
  {tech.current_level > 0 && (
    <p className="active-bonus">
      Active Bonus (Lv{tech.current_level}): {formatEffect(tech.effects, tech.current_level)}
    </p>
  )}
  {tech.current_level < tech.max_level && (
    <p>Next Level: {formatEffect(tech.effects, tech.current_level + 1)}</p>
  )}
  <p>Cost: {formatNumber(tech.cost_next_level.gold)} Gold</p>
</div>

// Helper function:
function formatEffect(effects: any, level: number): string {
  const type = effects.type;
  const perLevel = effects.per_level || 0;
  const value = perLevel * level;

  switch (type) {
    case 'metal_output': return `+${value}% Metal Production`;
    case 'he3_output': return `+${value}% He3 Production`;
    case 'gold_output': return `+${value}% Gold Production`;
    case 'build_speed': return `+${value}% Build Speed`;
    // ... etc
    default: return effects.description || '';
  }
}
```

**2. Visual Feedback for Active Bonuses**

**File:** `ResearchPanel.tsx`

Add green highlight or icon to techs with active bonuses:

```tsx
<div className={`tech-card ${tech.current_level > 0 ? 'active' : ''}`}>
  {tech.current_level > 0 && <span className="active-badge">Lv{tech.current_level}</span>}
  {/* ... */}
</div>
```

**CSS:**
```css
.tech-card.active {
  border-color: #4ade80;
  box-shadow: 0 0 8px rgba(74, 222, 128, 0.3);
}

.active-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  background: #4ade80;
  color: #000;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: bold;
}
```

---

## 5. Integration Points

**Files to Modify:**

1. **NEW:** `/backend/internal/services/tech_effects.go` - Tech bonus calculation
2. **MODIFY:** `/backend/internal/handlers/buildings.go` - Apply bonuses to construction (3 locations)
3. **MODIFY:** `/backend/internal/handlers/ship_factory.go` - Apply bonuses to ship production (2 locations)
4. **MODIFY:** `/frontend/src/components/panels/ResearchPanel.tsx` - Tooltips + active badges
5. **MODIFY:** `/frontend/src/types/index.ts` - Add `TechEffect` type

**No database changes required.**

---

## 6. QA Checklist

### Tech Trees (7 trees, 111 techs)

**Logistics Construction Science (11 techs):**
- [ ] Concurrent Construction (+1 construction slot) - CRITICAL
- [ ] Construction Boost (+1-15% build speed)
- [ ] Quality Materials (-1-15% build cost)
- [ ] Ship Building Boost (+1-15% ship build speed)
- [ ] Ship Building Logistics (-1-15% ship cost)
- [ ] Sync Shipbuilding (+1 ship production slot) - CRITICAL
- [ ] Repair Technology (+1-10% repair %)
- [ ] High Yield Mining (+1-10% Metal)
- [ ] High Yield Chemistry (+1-10% He3)
- [ ] High Yield Investing (+1-10% Gold)
- [ ] Expand Capacity (+50k-350k storage)

**Ballistics Science (14 techs):**
- [ ] All 14 techs researchable
- [ ] Prerequisites enforced correctly
- [ ] Combat bonuses apply (defer to Module 5)

**Ship Defense Science (20 techs):**
- [ ] All 20 techs researchable
- [ ] Shield/Structure branches work
- [ ] Combat bonuses apply (defer to Module 5)

**Directional Science (15 techs):**
- [ ] All 15 techs researchable
- [ ] Optics tree complete
- [ ] Combat bonuses apply (defer to Module 5)

**Missile Science (14 techs):**
- [ ] All 14 techs researchable
- [ ] Missile tree complete
- [ ] Combat bonuses apply (defer to Module 5)

**Ship-Based Science (10 techs):**
- [ ] All 10 techs researchable
- [ ] Fighter weapons tree complete
- [ ] Combat bonuses apply (defer to Module 5)

**Planetary Defense Science (8 techs):**
- [ ] All 8 techs researchable
- [ ] Defense enhancement techs work
- [ ] Combat bonuses apply (defer to Module 5)

### Tech Effects Application

**Production Bonuses:**
- [ ] High Yield Mining: Metal production increases correctly
- [ ] High Yield Chemistry: He3 production increases correctly
- [ ] High Yield Investing: Gold production increases correctly
- [ ] Expand Capacity: Storage capacity increases correctly

**Construction Bonuses:**
- [ ] Concurrent Construction: +1 building slot appears
- [ ] Construction Boost: Building time reduces correctly
- [ ] Quality Materials: Building costs reduce correctly

**Ship Production Bonuses:**
- [ ] Ship Building Boost: Ship build time reduces correctly
- [ ] Ship Building Logistics: Ship costs reduce correctly
- [ ] Sync Shipbuilding: 5th ship production slot appears

**Prerequisites:**
- [ ] Cannot research tech without prerequisites
- [ ] Tech Center level requirement enforced
- [ ] Max level cap enforced (10 for most, varies by tech)

**Auto-Complete:**
- [ ] Research completes automatically when timer expires
- [ ] Level increments correctly
- [ ] Quest progress updates (`research_tech` type)
- [ ] Active research clears

**Cancel/Speedup:**
- [ ] Cancel works (no refund currently - intended?)
- [ ] Speedup costs 3 vouchers per 30 minutes
- [ ] Speedup reduces timer correctly

**UI:**
- [ ] All 7 tree tabs display correctly
- [ ] Tech cards show correct state (locked/available/researching/completed)
- [ ] Countdown timer updates in real-time
- [ ] Tooltips show costs, prereqs, effects
- [ ] Active bonuses highlighted
- [ ] Start/cancel/speedup buttons work
- [ ] Confirmation modals appear
- [ ] Toast notifications display
- [ ] Error messages helpful

---

## 7. Task Breakdown

### Task 1: Create Tech Effects Service (2-3 hours)

**Owner:** backend-dev

**Files:**
- NEW: `/backend/internal/services/tech_effects.go`

**Deliverables:**
- `GetPlayerTechBonuses(playerID)` function
- `GetTechLevel(playerID, techName)` function
- Unit tests (optional, time permitting)

**Acceptance:**
- Function returns correct bonuses for test player with High Yield Mining Lv5
- Function returns zero bonuses for new player

---

### Task 2: Integrate Tech Bonuses - Buildings (2-3 hours)

**Owner:** backend-dev

**Files:**
- MODIFY: `/backend/internal/handlers/buildings.go`

**Locations:**
1. `GetResources()` - Apply production bonuses (line ~50)
2. `ConstructBuilding()` - Apply build speed/cost bonuses (line ~150-180)
3. `ConstructBuilding()` - Check construction slots (line ~100)

**Deliverables:**
- Production rates reflect tech bonuses
- Build time/cost reflect tech bonuses
- Construction slots respect tech bonuses

**Acceptance:**
- Build Lv1 Metal Mine with/without High Yield Mining Lv5 → production differs by 5%
- Build Lv1 Civic Center with/without Construction Boost Lv3 → time differs by 4.5%
- Concurrent Construction Lv1 → can build 3 buildings simultaneously (was 2)

---

### Task 3: Integrate Tech Bonuses - Ships (2-3 hours)

**Owner:** backend-dev

**Files:**
- MODIFY: `/backend/internal/handlers/ship_factory.go`

**Locations:**
1. `BuildShips()` - Apply ship build speed/cost bonuses (line ~120-150)
2. `BuildShips()` - Check ship production slots (line ~80)

**Deliverables:**
- Ship build time/cost reflect tech bonuses
- Ship production slots respect tech bonuses

**Acceptance:**
- Build 10 Weikes with/without Ship Building Boost Lv5 → time differs by 7.5%
- Sync Shipbuilding Lv1 → 5th production slot available (was 4 max)

---

### Task 4: Frontend - Tech Effects Tooltips (1-2 hours)

**Owner:** frontend-dev

**Files:**
- MODIFY: `/frontend/src/components/panels/ResearchPanel.tsx`
- MODIFY: `/frontend/src/types/index.ts`

**Deliverables:**
- `formatEffect()` helper function
- Tooltips show active bonuses
- Tooltips show next level bonuses
- Green highlight/badge for active techs

**Acceptance:**
- Hover over High Yield Mining Lv5 → tooltip shows "+5% Metal Production"
- Hover over locked tech → shows prerequisites
- Active techs have green border

---

### Task 5: QA Testing (4-6 hours)

**Owner:** qa-agent

**Test Cases:**

**Production Bonuses (30 min):**
1. Research High Yield Mining to Lv5
2. Verify Metal production increased by 5%
3. Repeat for He3 and Gold
4. Test Expand Capacity (warehouse size)

**Construction Bonuses (30 min):**
1. Research Concurrent Construction to Lv1
2. Verify 3rd construction slot available
3. Research Construction Boost to Lv3
4. Build Civic Center, verify 4.5% time reduction
5. Research Quality Materials to Lv3
6. Build Metal Mine, verify 4.5% cost reduction

**Ship Production Bonuses (30 min):**
1. Research Ship Building Boost to Lv5
2. Build 10 Weikes, verify 7.5% time reduction
3. Research Ship Building Logistics to Lv3
4. Build 10 Weikes, verify 4.5% cost reduction
5. Research Sync Shipbuilding to Lv1
6. Verify 5th ship production slot available

**Prerequisites (1 hour):**
1. Try to research Construction Boost without Concurrent Construction → expect error
2. Research Concurrent Construction, then Construction Boost → success
3. Test 5 random techs with multi-level prerequisites

**Auto-Complete (30 min):**
1. Start research with 10-second timer (modify seed data temporarily)
2. Wait 10 seconds
3. Verify research completes automatically
4. Verify quest progress increments

**Cancel/Speedup (30 min):**
1. Start research
2. Cancel → verify research stops
3. Start research
4. Speedup 30 minutes → verify cost 3 vouchers, timer reduced

**UI (1 hour):**
1. Click through all 7 tree tabs
2. Verify all techs display correctly
3. Test start/cancel/speedup buttons
4. Verify tooltips show correct info
5. Verify active research countdown updates

**Edge Cases (1 hour):**
1. Try to research with insufficient Gold → expect error
2. Try to research while already researching in same tree → expect error
3. Try to research past max level → expect error
4. Start research in Ballistics tree, then start in Logistics tree → should succeed (different trees)

---

## 8. Risk Assessment

### LOW RISK
- ✅ Database schema complete
- ✅ Backend handlers functional
- ✅ Frontend UI complete
- ✅ Auto-complete worker exists

### MEDIUM RISK
- ⚠️ **Tech effects integration** - Touching 3 handlers (buildings, ship_factory, resources)
  - **Mitigation:** Small, isolated changes. Test each bonus type separately.
- ⚠️ **Combat bonuses deferred** - Won't see combat tech effects until Module 5
  - **Mitigation:** Document clearly. Implement lookup functions now, integrate later.

### POTENTIAL BLOCKERS
- **NONE IDENTIFIED**

---

## 9. Estimated Complexity

| Task | Complexity | Hours |
|------|-----------|-------|
| Tech Effects Service | Low | 2-3 |
| Buildings Integration | Medium | 2-3 |
| Ships Integration | Medium | 2-3 |
| Frontend Tooltips | Low | 1-2 |
| QA Testing | Medium | 4-6 |
| **Total** | **Medium** | **11-17 hours** |

**Estimated Duration:** 2-3 days (assuming 6-8 hours/day)

**Critical Path:**
1. Tech Effects Service (backend-dev)
2. Buildings Integration (backend-dev)
3. Ships Integration (backend-dev)
4. Frontend Tooltips (frontend-dev) - can run in parallel
5. QA Testing (qa-agent)

---

## 10. Success Criteria

**Module 1 is COMPLETE when:**

1. ✅ All 111 techs researchable across 7 trees
2. ✅ Prerequisites validated correctly
3. ✅ Research costs deducted (Gold only)
4. ✅ Tech Center level reduces research time (3% per level)
5. ✅ Auto-complete works (timer expires → level increases)
6. ✅ Production bonuses apply (Metal/He3/Gold output +%)
7. ✅ Construction bonuses apply (slots, speed, cost)
8. ✅ Ship production bonuses apply (slots, speed, cost)
9. ✅ Warehouse capacity bonuses apply
10. ✅ Frontend displays active bonuses in tooltips
11. ✅ Quest progress updates on research complete
12. ✅ Cancel/speedup work correctly
13. ✅ All 30 QA test cases pass
14. ✅ No critical bugs

**Combat bonuses:** Deferred to Module 5 (Combat System). Functions prepared but not integrated.

---

## 11. Dependencies

**Upstream (Required Before This Module):**
- ✅ Phase 1 complete (buildings, resources, Tech Center)
- ✅ Research migration seeded (111 techs)
- ✅ Quest system (auto-progress integration)

**Downstream (Modules That Depend On This):**
- Module 2: Blueprint Research (uses similar tech effect pattern)
- Module 5: Combat System (reads tech levels for combat bonuses)

---

## 12. Code Snippets Reference

### Tech Effects JSON Examples (from seed data)

```json
// Simple percentage bonus
{"type":"metal_output","per_level":1,"unit":"percent"}
// Lv5 → +5% Metal production

// Flat bonus
{"type":"construction_slots","per_level":1}
// Lv1 → +1 construction slot

// Flat amount per level
{"type":"warehouse_capacity","per_level":50000,"unit":"flat"}
// Lv7 → +350,000 storage

// Combat bonuses (for Module 5)
{"type":"ballistic_damage","per_level":5,"unit":"percent"}
// Lv10 → +50% ballistic damage
```

### Tech Prerequisites JSON Examples

```json
// No prerequisites
[]

// Single prerequisite
[{"tech":"concurrent_construction","level":1}]

// Multiple prerequisites (AND logic)
[
  {"tech":"ballistics_base","level":6},
  {"tech":"ballistic_malice","level":3}
]
```

---

## 13. Notes

**Tech Center Bonus:**
- Already implemented in `StartResearch()` handler
- 3% time reduction per Tech Center level
- Applied before research starts

**Gold-Only Costs:**
- All 111 techs cost Gold only (metal=0, he3=0)
- This is GO2 canonical behavior
- Schema retains metal/he3 columns for flexibility

**Max Levels:**
- Most techs: 10 levels
- Some techs: 1-7 levels (varies)
- Enforced by `max_level` column

**Tree-Level Concurrency:**
- Cannot research 2 techs in same tree simultaneously
- CAN research 1 tech per tree (max 7 concurrent researches)
- Current implementation: 1 active research total (simpler for MVP)
- **Future enhancement:** Allow 1 per tree

**Cancel Refund:**
- Current implementation: No refund
- **Future enhancement:** 50% refund (requires migration to add `cancel_count` tracking)

---

## 14. Future Enhancements (Out of Scope)

**Not Included in Module 1:**

1. **Multi-tree concurrent research** - Allow 1 active research per tree (currently 1 total)
2. **Cancel refund** - Refund 50% of resources on cancel
3. **Research queue** - Queue next tech to start after current completes
4. **Tech respec** - Reset tech tree (costs Gold/vouchers)
5. **Tech presets** - Save/load tech build templates
6. **Combat tech effects** - Full integration with combat system (Module 5)
7. **Tech unlocks for modules** - Gate modules behind tech prerequisites (Module 2)

---

**Document Status:** FINAL
**Ready for Implementation:** YES 🚀
**Next Step:** Assign tasks to backend-dev, frontend-dev, qa-agent
