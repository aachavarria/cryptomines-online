# Phase 4: Corps, RBPs & Galaxy Map

## Context

Phase 4 from the GDD roadmap: "Corps, RBPs, corp donations, Galactic Wars, Alliance Center". All mechanics are 1:1 GO2 faithful. The DB already has forward-thinking RBP fields in `planets` table and `rbp_attack` combat type in `combat_reports`. We need to create the corps/galaxy infrastructure and plug it into existing systems.

## Scope

| Feature | Description |
|---------|-------------|
| **Corps** | Create/join/leave, roles (leader/officer/member), donations, wealth, levels |
| **Galaxy Map** | 7x7 zone grid, each zone has 1 RBP planet |
| **RBPs** | Resource Bonus Planets (5%-280% bonus), NPC defenses, conquest |
| **Galactic Wars** | Corp vs Corp RBP attacks using existing combat engine |
| **Alliance Center** | Building prerequisite for corps features (already in building_types) |

## Team Structure

5 agents (all Sonnet), parallelized where possible:

| Agent | Role | Work |
|-------|------|------|
| **infra-dev** | DB + Backend services | Migration, corp_service.go, corp_worker.go |
| **backend-dev** | Handlers + Routes | corps.go (10 endpoints), dev.go additions, main.go routes |
| **frontend-dev** | Panels + Hooks | CorpsPanel, GalaxyMapPanel, useCorp, useGalaxy, types, api.ts |
| **frontend-nav** | Navigation + Bridge | SideNav unlock, GameActionsBridge, CSS |
| **tools-dev** | WebMCP tools | webmcp-ui.ts corp/galaxy tools, webmcp-dev.ts dev tools |

**Dependency order:**
1. infra-dev (migration + services) - FIRST, blocks backend-dev
2. backend-dev + frontend-dev + frontend-nav - PARALLEL after migration
3. tools-dev - LAST (needs bridge methods from frontend-nav)

## Implementation Details

### 1. Database Migration
**File:** `supabase/migrations/20260213000000_phase4_corps_galaxy.sql`

New tables:
- `corps` (id, name, tag, leader_id, level, wealth, max_members, description)
- `corp_members` (corp_id, player_id, role, contribution_points, daily_contribution)
- `corp_donations` (corp_id, player_id, metal/he3/gold, contribution_points)
- `galaxy_zones` (zone_x 0-6, zone_y 0-6, rbp_planet_id)
- `rbp_attacks` (rbp_planet_id, attacking_corp_id, combat_report_id, total_kills)

Alterations:
- `planets.controlling_corp_id` FK → `corps.id`

Seed data:
- 49 galaxy zones (7x7)
- 49 RBP planets (system-owned, is_rbp=true, 72h initial protection)

### 2. Backend Services
**File:** `backend/internal/services/corp_service.go`

- `CalculateRBPBonus(rbpLevel int) float64` - GO2 scaling: Lv1-10 = 5%+0.5%/lv, Lv11-20 = +1%/lv, etc. up to 280% at Lv100
- `GetCorpBonuses(corpID string) (*CorpBonuses, error)` - Aggregate all corp RBP bonuses
- `GetCorpRBPCount(corpID string) (int, error)` - Count controlled RBPs
- `CanCorpControlMoreRBPs(corpID string) (bool, error)` - Max RBPs = corp level

### 3. Backend Worker
**File:** `backend/internal/workers/corp_worker.go`

- `StartCorpWorker()` - 1hr ticker
- `ResetDailyContributions()` - Reset daily_contribution when day changes

### 4. Backend Handlers (10 endpoints)
**File:** `backend/internal/handlers/corps.go`

| Method | Route | Handler | Notes |
|--------|-------|---------|-------|
| GET | /api/corp | GetCorp | Player's corp or null |
| POST | /api/corp | CreateCorp | Requires Alliance Center Lv1+ |
| POST | /api/corp/join | JoinCorp | Requires Alliance Center Lv1+ |
| POST | /api/corp/leave | LeaveCorp | Leader cannot leave |
| GET | /api/corp/members | ListCorpMembers | With player names |
| POST | /api/corp/donate | DonateResources | Max 200 pts/day = 2M resources |
| PUT | /api/corp/members/{id}/role | UpdateMemberRole | Leader only |
| GET | /api/corp/search | SearchCorps | By name/tag |
| POST | /api/corp/rbp/{id}/attack | AttackRBP | Reuses combat engine |
| GET | /api/galaxy/map | GetGalaxyMap | 7x7 grid with ownership |

Dev endpoints (in dev.go):
- POST /api/dev/create-corp
- POST /api/dev/join-corp
- POST /api/dev/give-corp-wealth

### 5. Frontend Types
**File:** `frontend/src/types/corps.ts`

Interfaces: Corp, CorpBonuses, CorpMember, CreateCorpRequest, DonateRequest, GalaxyZone, GalaxyMapResponse, AttackRBPRequest, AttackRBPResponse

### 6. Frontend API Functions
**File:** `frontend/src/services/api.ts` (add ~100 lines)

10 functions: getCorp, createCorp, joinCorp, leaveCorp, listCorpMembers, donateResources, updateMemberRole, searchCorps, getGalaxyMap, attackRBP

### 7. Frontend Hooks
**Files:** `frontend/src/hooks/useCorp.ts`, `frontend/src/hooks/useGalaxy.ts`

- useCorp: corp state, members, CRUD operations, search, donate
- useGalaxy: zones grid, attack, refresh

### 8. Frontend Panels
**Files:** `frontend/src/components/panels/CorpsPanel.tsx`, `frontend/src/components/panels/GalaxyMapPanel.tsx`

CorpsPanel states:
- Not in corp → Create/Join tabs
- In corp (member) → Overview, Members, Donate tabs
- In corp (leader) → + Role management

GalaxyMapPanel:
- 7x7 HTML/CSS grid (no Three.js needed)
- Zone cards showing RBP name, corp owner, level, protection timer
- Click zone → RBP detail + Attack button

### 9. Navigation Updates
**Files:** SideNav.tsx, GameActionsBridge.tsx, bridge.ts

- Unlock `galaxy` and `corp` nav items (locked: false)
- Add corpOpen/galaxyOpen state + panel renders
- Extend bridge with corp/galaxy methods
- Add 'corp' and 'galaxy' to PanelId type

### 10. WebMCP Tools (12 total)
**Files:** webmcp-ui.ts, webmcp-dev.ts

Corp tools (8): ui_corp_open, ui_corp_get_info, ui_corp_create, ui_corp_join, ui_corp_donate, ui_corp_get_members, ui_corp_change_role, ui_corp_search
Galaxy tools (4): ui_galaxy_open, ui_galaxy_get_map, ui_galaxy_attack_rbp, ui_galaxy_get_rbp_detail
Dev tools (3): dev_create_corp, dev_join_corp, dev_give_corp_wealth

## Key Patterns to Follow

- **Backend handlers**: Direct SQL, `middleware.GetPlayerID(r)`, `database.DB.Query/Exec`, transactions with `defer tx.Rollback()`
- **Error handling**: Use `errs` package (NotFound, InsufficientResources, etc.)
- **Frontend panels**: `createPortal` overlay, `onClose` prop, ESC listener, LoadingButton
- **Hooks**: `useCallback` for functions, `dispatch SET_ERROR` for errors, cleanup in useEffect
- **WebMCP tools**: `ui_<category>_<action>` naming, `getBridge()` + `ok()`/`error()` helpers

## Verification

1. Run migration: `supabase db reset` or apply migration
2. Build backend: `go build ./cmd/server/`
3. Build frontend: `npm run build`
4. Manual test flow:
   - dev_create_corp → ui_corp_get_info → dev_give_resources → ui_corp_donate → ui_galaxy_open → ui_galaxy_get_map
5. RBP attack flow:
   - Build ships → Create fleet → ui_galaxy_attack_rbp → Check combat report

## Estimated Size

~3000 lines total across all files. Largest: corps.go (~800 lines), CorpsPanel.tsx (~500 lines), migration (~300 lines).
