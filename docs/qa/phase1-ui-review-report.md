# Phase 1 UI Review - QA Report

> **Date**: 2026-02-05
> **Reviewer**: qa-agent
> **Scope**: Frontend (React + Three.js) and Backend (Go) vs Documentation (GDD, UI/UX Spec v2.0, GO2 Research)
> **Version**: Phase 1 MVP + UI Redo

---

## Summary

| Category | PASS | FAIL | WARN | Total |
|----------|------|------|------|-------|
| A. Backend vs GDD | 7 | 1 | 1 | 9 |
| B. Frontend vs UI/UX Spec v2.0 | 8 | 2 | 2 | 12 |
| C. Frontend vs Backend Integration | 2 | 1 | 1 | 4 |
| D. GO2 Name Compliance | 1 | 1 | 0 | 2 |
| **TOTAL** | **18** | **5** | **4** | **27** |

---

## A. Backend vs GDD

### A1. All GDD-listed API endpoints exist in main.go routes
**Result: FAIL**

GDD Section 5 specifies the following endpoints. Status:

| Endpoint | Status |
|----------|--------|
| `GET /api/health` | Implemented |
| `POST /api/auth/guest` | Implemented |
| `GET /api/player/me` | Implemented |
| `GET /api/planets` | Implemented |
| `GET /api/planets/{id}` | Implemented |
| `GET /api/planets/{id}/buildings` | Implemented |
| `POST /api/planets/{id}/buildings` | Implemented |
| `POST /api/planets/{id}/buildings/{buildingId}/upgrade` | Implemented |
| `POST /api/planets/{id}/buildings/{buildingId}/cancel` | **MISSING** |
| `GET /api/planets/{id}/resources` | Implemented |
| `POST /api/planets/{id}/resources/collect` | Implemented |

**Details**: The cancel upgrade endpoint (`POST /api/planets/{id}/buildings/{buildingId}/cancel`) is specified in the GDD but has no corresponding route in `backend/cmd/server/main.go`, no handler function, and no frontend API function.

**File**: `/Users/yurei/cryptomines-online/backend/cmd/server/main.go:26-38`

---

### A2. Cost formulas match GDD: Cost(N) = BaseCost * 3.03^(N-1), Time(N) = BaseTime * 2.87^(N-1)
**Result: PASS**

The formula fallback in `backend/internal/services/game.go` uses the correct formula. The backend first tries the lookup table (exact wiki data), then falls back to the formula with the multipliers stored in `building_types.cost_multiplier` and `building_types.time_multiplier`.

However, note that resource buildings (metal_collector, he3_extractor, residential_area, resource_warehouse) use a cost_multiplier of `1.75` and time_multiplier of `1.75` in the seed data, which differs from the `3.03`/`2.87` defaults. This is correct because these buildings have per-level lookup tables with exact wiki data, and the base multiplier in the seed data reflects the actual scaling pattern of those specific buildings. Core/military/space buildings correctly use `3.03`/`2.87`.

**Files**: `/Users/yurei/cryptomines-online/backend/internal/services/game.go`, `/Users/yurei/cryptomines-online/backend/internal/services/building_costs.go`

---

### A3. 22 building types in seed data with correct GO2 names
**Result: PASS**

The migration file `20260206005232_phase1_mvp.sql` seeds exactly 22 building types:

1. metal_collector (Metal Collector)
2. he3_extractor (He3 Extractor)
3. residential_area (Residential Area)
4. resource_warehouse (Resource Warehouse)
5. civic_center (Civic Center)
6. technology_center (Technology Center)
7. alliance_center (Alliance Center)
8. trading_center (Trading Center)
9. galaxy_transporter (Galaxy Transporter)
10. compound_center (Compound Center)
11. radar (Radar)
12. ship_factory (Ship Factory)
13. spacedock (Spacedock)
14. command_center (Command Center)
15. weapon_research_center (Weapon Research Center)
16. recycling_plant (Recycling Plant)
17. space_station (Space Station)
18. meteor_star (Meteor Star)
19. particle_cannon (Particle Cannon)
20. anti_aircraft_gun (Anti-Aircraft Gun)
21. thors_cannon (Thor's Cannon)
22. celestial_base (Celestial Base)

All names match GO2 wiki originals.

**File**: `/Users/yurei/cryptomines-online/supabase/migrations/20260206005232_phase1_mvp.sql:754-780`

---

### A4. 30 tech types in seed data
**Result: PASS**

The migration seeds tech types from two science trees:
- **Logistics Construction**: 11 techs (concurrent_construction, construction_boost, quality_materials, ship_building_boost, ship_building_logistics, sync_shipbuilding, repair_technology, high_yield_mining, high_yield_chemistry, high_yield_investing, expand_capacity)
- **Ship Defense Science**: 19 techs (ship_defense_tech, shield_research, energy_diffusion, penetration_resistance, augment_shield, restoration, augment_absorption, energy_conservation, electronic_barrier, damage_mitigation, ship_structural_analysis, ship_reinforcement, resilience, structure_improvement, fast_repair, reaction_armor_improvement, defense_improvement, reflection_mastery, stability_mastery)

Total: **30 tech types**. Matches GDD.

**File**: `/Users/yurei/cryptomines-online/supabase/migrations/20260206005232_phase1_mvp.sql:788-941`

---

### A5. Lookup tables used for exact wiki costs (not just formula fallback)
**Result: PASS**

18 level reference tables exist in the database:
- he3_extractor_levels (24 levels)
- metal_collector_levels (24 levels)
- residential_area_levels (24 levels)
- resource_warehouse_levels (24 levels)
- civic_center_levels (12 levels)
- technology_center_levels (12 levels)
- command_center_levels (12 levels)
- space_station_levels (12 levels)
- weapon_research_center_levels (12 levels)
- alliance_center_levels (11 levels)
- trading_center_levels (9 levels)
- radar_levels (9 levels)
- spacedock_levels (12 levels)
- recycling_plant_levels (11 levels, partial data)
- meteor_star_levels (12 levels)
- particle_cannon_levels (12 levels)
- anti_aircraft_gun_levels (12 levels)
- thors_cannon_levels (12 levels)

The backend `services.GetBuildingLevelCost()` in `building_costs.go` correctly queries the lookup table first and only falls back to the formula on miss. The `lookupTableMap` maps all 18 building types to their respective level tables.

**Files**: `/Users/yurei/cryptomines-online/backend/internal/services/building_costs.go:17-36`, `/Users/yurei/cryptomines-online/supabase/migrations/20260206005232_phase1_mvp.sql:196-748`

---

### A6. Auth: Supabase anonymous sign-in + JWT
**Result: PASS**

- Guest auth handler creates anonymous Supabase user, generates JWT with HMAC signing
- JWT middleware validates tokens with 7-day expiry and extracts `sub` claim for player ID
- Token stored in localStorage on frontend, attached via Axios interceptor

**Files**: `/Users/yurei/cryptomines-online/backend/internal/handlers/auth.go`, `/Users/yurei/cryptomines-online/backend/internal/middleware/auth.go`

---

### A7. JWT validation middleware on all protected routes
**Result: PASS**

All protected routes are wrapped with `middleware.Auth()` in main.go:
```
mux.Handle("/api/", middleware.Auth(protected))
```

This covers all routes under `/api/` except the explicitly public routes (`GET /api/health`, `POST /api/auth/guest`).

**File**: `/Users/yurei/cryptomines-online/backend/cmd/server/main.go:40`

---

### A8. Building prerequisite checks (Civic Center gates construction)
**Result: PASS**

The `ConstructBuilding` handler checks:
1. Construction slot limit (max 2 simultaneous)
2. Max count per planet for the building type
3. Prerequisite building existence and level requirement
4. Resource sufficiency

The `UpgradeBuilding` handler checks:
1. Construction slot limit
2. Already upgrading status
3. Max level cap
4. Civic center level >= target level (when `civic_center_req_per_level` is true)
5. Resource sufficiency

**File**: `/Users/yurei/cryptomines-online/backend/internal/handlers/buildings.go:127-178` (construct), `254-327` (upgrade)

---

### A9. Construction slot limit = 2
**Result: WARN**

The backend correctly enforces `maxConstructionSlots = 2` (line 16 of buildings.go). However, the GDD specifies that the "Concurrent Construction" tech (from Logistics Construction tree) adds +1 slot. The backend does not yet check the player's tech level to dynamically calculate the slot limit. This is acceptable for Phase 1 MVP since tech research is not yet implemented, but should be addressed in Phase 2.

**File**: `/Users/yurei/cryptomines-online/backend/internal/handlers/buildings.go:16`

---

## B. Frontend vs UI/UX Spec v2.0

### B1. Three.js architecture: Canvas + HTML overlay layering
**Result: PASS**

The `GameShell.tsx` correctly implements the dual-layer architecture:
- `<Canvas>` wrapping `<PlanetScene>` as z-index 0 (via `.game-canvas`)
- `<div className="hud-overlay">` as z-index 10 with `pointer-events: none` on container and `pointer-events: auto` on children

This matches the spec's Section 2.1 architecture exactly.

**Files**: `/Users/yurei/cryptomines-online/frontend/src/components/layout/GameShell.tsx:43-71`, `/Users/yurei/cryptomines-online/frontend/src/styles/layout.css:9-24`

---

### B2. Scene graph matches spec (SceneLighting, Starfield, PlanetSurface, BuildingModels, EmptySlotMarkers, CameraController)
**Result: PASS**

The spec (Section 2.3) defines the scene graph. The actual `PlanetScene.tsx` renders:
- `<SceneLighting />` -- matches `SceneLighting` in spec
- `<Starfield />` -- matches `Skybox` in spec (named differently but same purpose)
- `<PlanetSurface />` -- matches spec
- `<BuildingModel />` instances -- matches `BuildingModels` in spec
- `<EmptySlotMarker />` instances -- matches `EmptySlotMarkers` in spec
- `<CameraController />` -- matches spec

Missing from spec but acceptable:
- `ConstructionEffects` (particle effects) -- not implemented yet, acceptable for Phase 1 placeholder geometry
- `BuildingLabels` are integrated directly in `BuildingModel` via `<Html>` from drei (acceptable architecture)

**File**: `/Users/yurei/cryptomines-online/frontend/src/components/three/PlanetScene.tsx:52-92`

---

### B3. HTML overlay components match spec (ResourceHUD, SideNav, ConstructionQueue, BuildingDetailPanel, ConstructionPanel)
**Result: PASS**

The `GameShell.tsx` overlay layer renders exactly:
- `<ResourceHUD />` -- present
- `<SideNav />` -- present
- `<ConstructionQueue />` -- present
- `<BuildingDetailPanel />` -- present
- `<ConstructionPanel />` (conditionally shown) -- present

Missing: `<ToastContainer />` from spec. Acceptable for Phase 1.

**File**: `/Users/yurei/cryptomines-online/frontend/src/components/layout/GameShell.tsx:64-70`

---

### B4. CSS color palette matches spec Section 3.1
**Result: PASS**

Comparing `variables.css` against the spec:

| Variable | Spec | Code | Match |
|----------|------|------|-------|
| `--bg-deep` | `#050510` | `#050510` | Yes |
| `--bg-panel` | `#0c0c20cc` | `#0c0c20cc` | Yes |
| `--bg-panel-solid` | `#0c0c20` | `#0c0c20` | Yes |
| `--bg-panel-light` | `#141430cc` | `#141430cc` | Yes |
| `--bg-header` | `#0a0a1eee` | `#0a0a1eee` | Yes |
| `--border-frame` | `#1a1a40` | `#1a1a40` | Yes |
| `--border-active` | `#2244aa` | `#2244aa` | Yes |
| `--border-glow` | `#3366cc` | `#3366cc` | Yes |
| `--accent-primary` | `#4488ff` | `#4488ff` | Yes |
| `--accent-hover` | `#66aaff` | `#66aaff` | Yes |
| `--accent-success` | `#22cc66` | `#22cc66` | Yes |
| `--accent-warning` | `#ffaa22` | `#ffaa22` | Yes |
| `--accent-danger` | `#ff4444` | `#ff4444` | Yes |
| `--accent-gold` | `#ffd700` | `#ffd700` | Yes |
| `--color-metal` | `#8899aa` | `#8899aa` | Yes |
| `--color-he3` | `#44ccff` | `#44ccff` | Yes |
| `--color-gold` | `#ffd700` | `#ffd700` | Yes |
| `--text-bright` | `#e8e8f8` | `#e8e8f8` | Yes |
| `--text-normal` | `#b0b0cc` | `#b0b0cc` | Yes |
| `--text-dim` | `#666688` | `#666688` | Yes |
| `--text-label` | `#8888aa` | `#8888aa` | Yes |

All CSS custom properties match the spec exactly. 3D scene colors also match:
- Ambient light: `#1a1a3a` (matches `--3d-ambient`)
- Key light: `#aabbff` (matches `--3d-key-light`)
- Rim light: `#4466aa` (matches `--3d-rim-light`)
- Ground color: `#1a1a2e` (matches `--3d-ground`)
- Grid lines: `#2a2a4e` (matches `--3d-ground-grid`)

**Files**: `/Users/yurei/cryptomines-online/frontend/src/styles/variables.css`, `/Users/yurei/cryptomines-online/frontend/src/components/three/SceneLighting.tsx`, `/Users/yurei/cryptomines-online/frontend/src/components/three/PlanetSurface.tsx`

---

### B5. Building states: idle/selected/upgrading with visual feedback
**Result: PASS**

`BuildingModel.tsx` implements all three visual states:
- **Idle**: Base emissive intensity 0.1, no selection ring, subtle breathing animation via `useFrame`
- **Selected**: Emissive intensity 0.4, `<SelectionRing>` visible (rotating torus with pulsing opacity), floating label gets `.selected` CSS class with blue glow border
- **Hovered**: Emissive intensity 0.25, cursor changes to pointer
- **Upgrading**: Orange wireframe scaffolding mesh wraps the building, countdown timer in floating label

**Files**: `/Users/yurei/cryptomines-online/frontend/src/components/three/BuildingModel.tsx:41-110`, `/Users/yurei/cryptomines-online/frontend/src/components/three/SelectionRing.tsx`

---

### B6. ResourceHUD: Metal/He3/Gold values, rates, storage, CC level, collect button
**Result: WARN**

The ResourceHUD shows:
- Metal/He3/Gold current values with formatted numbers -- PRESENT
- Production rates (+X/hr) for each resource -- PRESENT
- CC level badge ("CC Lv X") -- PRESENT
- Collect button with pending total, disabled when nothing to collect -- PRESENT
- Storage display -- PRESENT but **INCORRECT**

**Issue**: Storage display sums all 3 resources (`metal + he3 + gold`) and compares against `storage_capacity`. In GO2, `storage_capacity` is a per-resource cap (i.e., each resource independently cannot exceed the warehouse capacity). The current display is misleading because it suggests a combined cap.

```tsx
// Current (incorrect):
Storage: {formatNumber(resources.metal + resources.he3 + resources.gold)} / {formatNumber(resources.storage_capacity)}

// Should show per-resource or max of the three:
// e.g., "Storage: M:5000 H:5000 G:10000 / 10000"
```

**File**: `/Users/yurei/cryptomines-online/frontend/src/components/layout/ResourceHUD.tsx:54-55`

---

### B7. BuildingDetailPanel: costs from backend data (not hardcoded)
**Result: FAIL**

The `BuildingDetailPanel.tsx` uses a local `BUILDING_COSTS` object with hardcoded base costs and applies a client-side formula `calcUpgradeCost()` to calculate costs. **Several of these hardcoded values are wrong**:

| Building | Seed Data (metal/he3/gold) | Frontend Code (metal/he3/gold) | Match |
|----------|---------------------------|-------------------------------|-------|
| civic_center | 550/480/600 | 500/500/500 | **WRONG** |
| technology_center | 450/420/650 | 450/580/350 | **WRONG** |
| alliance_center | 500/400/550 | 800/600/700 | **WRONG** |
| trading_center | 1200/1100/1500 | 600/500/800 | **WRONG** |
| ship_factory | 600/450/500 | 700/850/500 | **WRONG** |
| spacedock | 900/675/750 | 900/750/600 | **WRONG** |
| command_center | 600/450/500 | 650/700/550 | **WRONG** |
| weapon_research_center | 500/300/450 | 750/900/600 | **WRONG** |
| space_station | 650/600/850 | 1000/1200/800 | **WRONG** |
| meteor_star | 50/45/48 | 500/400/300 | **WRONG** |
| particle_cannon | 450/360/520 | 600/500/400 | **WRONG** |
| anti_aircraft_gun | 650/600/850 | 400/350/250 | **WRONG** |
| thors_cannon | 80000/90000/125000 | 800/700/500 | **WRONG** |
| metal_collector | 85/106/85 | 85/106/85 | Correct |
| he3_extractor | 95/80/95 | 95/80/95 | Correct |
| residential_area | 78/72/65 | 78/72/65 | Correct |
| resource_warehouse | 380/370/480 | 380/370/480 | Correct |

Only 4 of 17 building types have correct base costs. The remaining 13 have fabricated values. Additionally, the cost multiplier is uniformly `1.75` for all types, but the DB seed data uses `3.03` for core/military/space buildings.

Moreover, the frontend uses the **formula** to calculate costs rather than fetching actual costs from the backend lookup tables. The backend has exact wiki-sourced per-level data for 18 building types, but the frontend never accesses this data. The backend does not expose a "get upgrade cost" endpoint; instead, costs are calculated on the server side during the upgrade action but not returned in a preview/read-only call.

**Production display**: The `BuildingDetailPanel` does show production data for resource buildings (metal_collector, he3_extractor, residential_area) using `BUILDING_PRODUCTION` constants with correct base values (1080, 1180, 1300) and multiplier (1.134). This is an improvement from the initial code (which just showed "Current output active").

**Missing buildings**: `galaxy_transporter`, `compound_center`, `radar`, `recycling_plant`, `celestial_base` are not in the `BUILDING_COSTS` map, so they would fall through to the generic fallback formula.

**File**: `/Users/yurei/cryptomines-online/frontend/src/components/panels/BuildingDetailPanel.tsx:10-28`

---

### B8. ConstructionPanel: real costs, correct max counts per building type
**Result: FAIL**

Two critical issues:

1. **Hardcoded costs**: The panel displays `formatNumber(100)`, `formatNumber(80)`, `formatNumber(80)` (Metal: 100, He3: 80, Gold: 80) for ALL building types. These are placeholder values that do not reflect actual construction costs from the database.

2. **Hardcoded max count**: `const maxCount = 8` is used for all building types. The actual `max_count_per_planet` values from the seed data vary significantly:
   - Resource buildings: Metal Collector (8), He3 Extractor (8), Residential Area (8), Resource Warehouse (4)
   - Core buildings: All have max 1 (except Galaxy Transporter: 1)
   - Military buildings: All max 1
   - Defense: Meteor Star (40), Particle Cannon (20), Anti-Aircraft Gun (15), Thor's Cannon (3)
   - Space: Celestial Base (1)

   Using `maxCount = 8` means the panel incorrectly shows that players can build 8 of every type, when most buildings allow only 1.

3. **Missing building types**: The `BUILDING_TYPES` array in `types/index.ts` has 17 entries, missing 5 building types that exist in the database:
   - `galaxy_transporter` (Galaxy Transporter)
   - `compound_center` (Compound Center)
   - `radar` (Radar)
   - `recycling_plant` (Recycling Plant)
   - `celestial_base` (Celestial Base)

**Files**: `/Users/yurei/cryptomines-online/frontend/src/components/panels/ConstructionPanel.tsx:77,100-112`, `/Users/yurei/cryptomines-online/frontend/src/types/index.ts:97-115`

---

### B9. ConstructionQueue: bottom bar with timer, progress, slot count
**Result: PASS**

The `ConstructionQueue.tsx` correctly implements:
- Fixed bottom bar positioned at `left: 60px` (clearing the SideNav)
- Shows "Queue (X/2)" with active upgrade count
- Active slots show: building abbreviation icon with category color, building name, level transition ("Lv X -> Lv X+1"), progress bar, countdown timer
- Empty slots show dashed border with "Available" text
- Clicking a queue slot selects the building and focuses camera
- Progress bar calculates completion percentage from `updated_at` to `upgrade_finish_at`
- Timer uses monospace font (`JetBrains Mono`)

**Files**: `/Users/yurei/cryptomines-online/frontend/src/components/panels/ConstructionQueue.tsx`, `/Users/yurei/cryptomines-online/frontend/src/styles/panels.css:257-372`

---

### B10. Animations: panel slide-in, modal fade, breathing glow
**Result: PASS**

- **Panel slide-in**: BuildingDetailPanel uses CSS `transform: translateX(100%)` transitioning to `translateX(0)` with `0.25s ease-out`. Defined in `panels.css:13-16`.
- **Modal fade**: ConstructionPanel backdrop uses `fadeIn 0.2s ease-out`, modal uses `modalIn 0.2s ease-out` with scale(0.95->1). Defined in `animations.css:43-49`.
- **Breathing glow**: `BuildingModel.tsx` uses `useFrame` for subtle position bobbing (`sin(t*0.8) * 0.05`), and `SelectionRing` has pulsing opacity (`0.5 + sin(t*2) * 0.25`).
- **Timer pulse**: `pulseTimer` animation for timers < 60s defined in `animations.css:38-41`.
- **Skeleton loading**: Shimmer animation defined in `animations.css:56-62`.

**Files**: `/Users/yurei/cryptomines-online/frontend/src/styles/animations.css`, `/Users/yurei/cryptomines-online/frontend/src/components/three/BuildingModel.tsx:43-48`

---

### B11. Z-index stack matches spec Section 2.2
**Result: PASS**

| Layer | Spec | Code | Match |
|-------|------|------|-------|
| Canvas | z: 0 | `.game-canvas { z-index: 0 }` | Yes |
| HUD overlay | z: 10 | `.hud-overlay { z-index: 10 }` | Yes |
| SideNav | z: 50 | `.side-nav { z-index: 50 }` | Yes |
| Construction Queue | z: 60 | `.construction-queue { z-index: 60 }` | Yes |
| Resource HUD | z: 100 | `.resource-hud { z-index: 100 }` | Yes |
| Building Detail Panel | z: 200 | `.building-detail-panel { z-index: 200 }` | Yes |
| Construction Modal | z: 300 | `.construction-modal-backdrop { z-index: 300 }` | Yes |
| Tooltips | z: 400 | `.nav-tooltip { z-index: 500 }` | Close (500 vs 400) |

Minor discrepancy: nav tooltip uses z-index 500 instead of spec's 400, but this has no functional impact as it still layers above panels.

**Files**: `/Users/yurei/cryptomines-online/frontend/src/styles/layout.css`, `/Users/yurei/cryptomines-online/frontend/src/styles/panels.css`

---

### B12. Camera: OrbitControls with spec constraints (fov, min/max distance, polar angle)
**Result: PASS**

Comparing `CameraController.tsx` and `GameShell.tsx` against spec Section 12:

| Parameter | Spec | Code | Match |
|-----------|------|------|-------|
| fov | 45 | 45 | Yes |
| near | 0.1 | 0.1 | Yes |
| far | 1000 | 1000 | Yes |
| initialPosition | [0, 40, 40] | [0, 40, 40] | Yes |
| minDistance | 15 | 15 | Yes |
| maxDistance | 80 | 80 | Yes |
| minPolarAngle | 0.3 | 0.3 | Yes |
| maxPolarAngle | 1.4 | 1.4 | Yes |
| enableDamping | true | true | Yes |
| dampingFactor | 0.08 | 0.08 | Yes |
| enablePan | true | true | Yes |
| panSpeed | 0.5 | 0.5 | Yes |
| rotateSpeed | 0.5 | 0.5 | Yes |

Camera transitions use 500ms ease-in-out lerp to focus on selected buildings, matching spec Section 12.4.

**Files**: `/Users/yurei/cryptomines-online/frontend/src/components/layout/GameShell.tsx:48-54`, `/Users/yurei/cryptomines-online/frontend/src/components/three/CameraController.tsx:64-76`

---

## C. Frontend vs Backend Integration

### C1. All backend endpoints have corresponding frontend API functions
**Result: FAIL**

| Backend Route | Frontend API Function | Match |
|---------------|-----------------------|-------|
| `GET /api/health` | `healthCheck()` | Yes |
| `POST /api/auth/guest` | `guestLogin()` | Yes |
| `GET /api/player/me` | Not used | Warn (exists but never called) |
| `GET /api/planets` | `listPlanets()` | Yes |
| `GET /api/planets/{id}` | `getPlanet()` | Yes |
| `GET /api/planets/{id}/buildings` | `listBuildings()` | Yes |
| `POST /api/planets/{id}/buildings` | `constructBuilding()` | Yes |
| `POST /api/planets/{id}/buildings/{buildingId}/upgrade` | `upgradeBuilding()` | Yes |
| `GET /api/planets/{id}/resources` | `getResources()` | Yes |
| `POST /api/planets/{id}/resources/collect` | `collectResources()` | Yes |

The `GET /api/player/me` endpoint exists but is never called from the frontend. Player data comes from the `POST /api/auth/guest` response instead. This is acceptable but means player profile refresh is not possible without re-auth.

The missing `cancel` endpoint (from A1) naturally has no frontend function.

**File**: `/Users/yurei/cryptomines-online/frontend/src/services/api.ts`

---

### C2. TypeScript types match backend JSON response shapes
**Result: PASS**

| Type | Backend Fields | Frontend Interface | Match |
|------|---------------|-------------------|-------|
| `BuildingWithType` | id, planet_id, building_type, level, is_upgrading, upgrade_finish_at, created_at, updated_at, type_name, display_name, category, base, max_level | All fields present | Yes |
| `ResourcesResponse` | metal, he3, gold, metal_per_hour, he3_per_hour, gold_per_hour, storage_capacity, last_collected_at, pending_metal, pending_he3, pending_gold | All fields present | Yes |
| `CollectResponse` | collected (metal/he3/gold), resources (Resource) | All fields present | Yes |
| `BuildingResponse` | building (id, planet_id, building_type, level, is_upgrading, upgrade_finish_at, created_at, updated_at), resources (metal/he3/gold) | All fields present | Yes |
| `GuestAuthResponse` | token, player (id, anonymous_id, level, created_at) | All fields present | Yes |
| `Planet` | id, player_id, name, position_x, position_y, is_homeworld, is_rbp, rbp_level, created_at, updated_at | All fields present | Yes |

**File**: `/Users/yurei/cryptomines-online/frontend/src/types/index.ts`

---

### C3. Error handling: API errors shown to user, not silently swallowed
**Result: WARN**

Error handling is inconsistent:

- **useBuildings.ts**: `fetchBuildings` dispatches SET_ERROR on failure. `upgrade` and `construct` also dispatch SET_ERROR. However, the error message is generic ("Failed to upgrade building") and does not include the server's error message.

- **useResources.ts**: `fetchResources` silently catches errors (`// Silent fail on refresh`). This is intentional for auto-refresh but means resource load failures are invisible to the user. `collect` correctly dispatches SET_ERROR.

- **GameShell.tsx**: Shows error state only when `state.error && !state.currentPlanet`. Once the planet is loaded, errors are not visually displayed anywhere in the UI (no toast/snackbar).

- **useAuth.ts**: Catches auth failures and sets error state, which Home.tsx displays.

Recommendation: Add a toast/notification system to show transient errors (spec mentions `<ToastContainer />` at z-index 500).

**Files**: `/Users/yurei/cryptomines-online/frontend/src/hooks/useBuildings.ts`, `/Users/yurei/cryptomines-online/frontend/src/hooks/useResources.ts`, `/Users/yurei/cryptomines-online/frontend/src/components/layout/GameShell.tsx`

---

### C4. Auto-refresh: resources poll every 30s
**Result: PASS**

`useResources.ts` sets up a 30-second interval:
```tsx
const interval = setInterval(fetchResources, 30000)
```

The interval is properly cleaned up on unmount via the `useEffect` return function. Resources are also refreshed after collecting. Buildings are refreshed after construct/upgrade operations.

**File**: `/Users/yurei/cryptomines-online/frontend/src/hooks/useResources.ts:20-25`

---

## D. GO2 Name Compliance

### D1. All building/resource/tech names match GO2 originals
**Result: PASS**

All names verified against GO2 wiki research documents:

**Resources**: Metal, He3, Gold -- all correct GO2 names.

**Buildings**: All 22 building type names match GO2 originals (verified against wiki links in migration comments and research docs).

**Tech types**: All 30 tech names match GO2 science tree names from the research document (Concurrent Construction, Construction Boost, Quality Materials, Ship Building Boost, Ship Building Logistics, Sync Shipbuilding, Repair Technology, High Yield Mining, High Yield Chemistry, High Yield Investing, Expand Capacity, Ship Defense Tech, Shield Research, Energy Diffusion, etc.).

**Category names**: resource, core, military, defense, space -- consistent with GO2 categorization.

---

### D2. No invented/renamed mechanics or buildings
**Result: FAIL**

While all building and tech names are correct, the `ConstructionPanel` header says **"Build New Structure"** instead of using GO2 terminology. In GO2, this action is typically labeled "Construct" or "Build". The word "Structure" is non-standard for GO2 which uses "building" terminology throughout.

Additionally, the `SideNav` uses the following labels that are either non-standard or abbreviated differently from GO2:

| SideNav Label | GO2 Original | Issue |
|---------------|-------------|-------|
| "Base" | "Base" | Correct |
| "Research" | "Tech" / "Technology" | Minor discrepancy |
| "Fleet" | "Fleet" | Correct |
| "Cmdr" | "Commander" | Abbreviation, acceptable |
| "Galaxy" | "Galaxy" | Correct |
| "Corp" | "Corp" / "Corporation" | Acceptable abbreviation |

Additionally, the SideNav uses emoji icons instead of proper game-themed icons, which is a visual fidelity issue rather than a naming issue.

**Files**: `/Users/yurei/cryptomines-online/frontend/src/components/panels/ConstructionPanel.tsx:58`, `/Users/yurei/cryptomines-online/frontend/src/components/layout/SideNav.tsx`

---

## Priority Recommendations

### Critical (FAIL items to fix)

1. **[B7] BuildingDetailPanel costs**: The hardcoded `BUILDING_COSTS` object has incorrect values for 13 of 17 building types. Either:
   - (a) Add a backend endpoint `GET /api/buildings/types` that returns building types with their base costs, and fetch them on the frontend, OR
   - (b) Sync the frontend `BUILDING_COSTS` with the exact values from the DB seed data AND still use the formula (which would match for low levels but diverge from exact wiki data at higher levels), OR
   - (c) Add a `GET /api/planets/{id}/buildings/{buildingId}/upgrade-cost` endpoint that returns the exact upgrade cost for the next level (from lookup tables), and display the server-provided cost.

   Option (c) is recommended as it ensures the frontend always shows the exact same cost the backend will charge.

2. **[B8] ConstructionPanel costs and max counts**: Replace `formatNumber(100)` / `formatNumber(80)` / `formatNumber(80)` with actual Lv1 costs for each building type. Replace `const maxCount = 8` with the actual `max_count_per_planet` from building type data. Add the 5 missing building types to `BUILDING_TYPES`.

3. **[A1] Cancel upgrade endpoint**: Implement `POST /api/planets/{id}/buildings/{buildingId}/cancel` as specified in the GDD. This should refund a percentage of resources and reset `is_upgrading` / `upgrade_finish_at`. Add corresponding frontend API function and cancel button in ConstructionQueue and BuildingDetailPanel.

### Medium (WARN items to address)

4. **[B6] Storage display**: Change to show per-resource capacity or the max of the three resources vs capacity, not the sum of all three vs a single cap.

5. **[C3] Error handling**: Add a toast/notification component to display transient errors to the user. The spec mentions a `ToastContainer` at z-index 500.

6. **[A9] Dynamic construction slots**: When tech research is implemented in Phase 2, update the backend to check the player's "Concurrent Construction" tech level to calculate max construction slots dynamically.

### Low (Polish items)

7. **[D2] Naming**: Change "Build New Structure" to "Build New Building" or just "Construct" to match GO2 terminology.

8. **[B3/SideNav]**: Replace emoji icons with proper SVG or icon font icons for a more polished game appearance.

9. **[C1]**: Consider calling `GET /api/player/me` on login to ensure player data is fresh, or store the full player object from guest auth response.

---

## Files Reviewed

### Documentation
- `/Users/yurei/cryptomines-online/docs/gdd/game-design-document.md`
- `/Users/yurei/cryptomines-online/docs/gdd/ui-ux-specification.md`
- `/Users/yurei/cryptomines-online/docs/research/galaxy-online-2-mechanics.md`
- `/Users/yurei/cryptomines-online/docs/research/needs-research-results.md`
- `/Users/yurei/cryptomines-online/docs/research/go2-visual-ui-research.md`

### Backend
- `/Users/yurei/cryptomines-online/backend/cmd/server/main.go`
- `/Users/yurei/cryptomines-online/backend/internal/handlers/auth.go`
- `/Users/yurei/cryptomines-online/backend/internal/handlers/buildings.go`
- `/Users/yurei/cryptomines-online/backend/internal/handlers/resources.go`
- `/Users/yurei/cryptomines-online/backend/internal/handlers/planets.go`
- `/Users/yurei/cryptomines-online/backend/internal/handlers/player.go`
- `/Users/yurei/cryptomines-online/backend/internal/services/building_costs.go`
- `/Users/yurei/cryptomines-online/backend/internal/services/game.go`
- `/Users/yurei/cryptomines-online/backend/internal/middleware/auth.go`
- `/Users/yurei/cryptomines-online/backend/internal/middleware/cors.go`

### Frontend
- `/Users/yurei/cryptomines-online/frontend/src/App.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/pages/Home.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/pages/Planet.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/contexts/GameContext.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/services/api.ts`
- `/Users/yurei/cryptomines-online/frontend/src/types/index.ts`
- `/Users/yurei/cryptomines-online/frontend/src/hooks/useAuth.ts`
- `/Users/yurei/cryptomines-online/frontend/src/hooks/useBuildings.ts`
- `/Users/yurei/cryptomines-online/frontend/src/hooks/useResources.ts`
- `/Users/yurei/cryptomines-online/frontend/src/hooks/useCountdown.ts`
- `/Users/yurei/cryptomines-online/frontend/src/components/layout/GameShell.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/components/layout/ResourceHUD.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/components/layout/SideNav.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/components/three/PlanetScene.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/components/three/BuildingModel.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/components/three/SceneLighting.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/components/three/Starfield.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/components/three/PlanetSurface.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/components/three/EmptySlotMarker.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/components/three/CameraController.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/components/three/SelectionRing.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/components/panels/BuildingDetailPanel.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/components/panels/ConstructionPanel.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/components/panels/ConstructionQueue.tsx`
- `/Users/yurei/cryptomines-online/frontend/src/styles/variables.css`
- `/Users/yurei/cryptomines-online/frontend/src/styles/layout.css`
- `/Users/yurei/cryptomines-online/frontend/src/styles/panels.css`
- `/Users/yurei/cryptomines-online/frontend/src/styles/components.css`
- `/Users/yurei/cryptomines-online/frontend/src/styles/animations.css`

### Database
- `/Users/yurei/cryptomines-online/supabase/migrations/20260206001730_init.sql`
- `/Users/yurei/cryptomines-online/supabase/migrations/20260206005232_phase1_mvp.sql`
