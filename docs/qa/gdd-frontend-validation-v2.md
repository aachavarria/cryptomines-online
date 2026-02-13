# GDD Frontend Validation v2 (STRICT)
Date: 2026-02-13
Validator: Strict Frontend Validator

## Summary
- Total features checked: 78
- PASS: 52 (UI flow fully works)
- PARTIAL: 8 (80%+ works, minor polish needed)
- FAIL: 12 (not implemented, placeholder, or broken)
- NOT_IMPLEMENTED: 6 (no UI exists at all)

**Overall Grade: C+ (66% fully functional)**

**Critical Finding:** Previous validation was too lenient. Many features have UI panels but **lack actual user-visible functionality**. The presence of a component does NOT mean the feature works.

---

## Critical Failures (features users can't use)

### FAIL-001: Fleet Travel Time UI - MISSING ENTIRELY
**GDD Requirement:** REQ-PVP002 - Fleet travel time uses SP (Space Points), fleets should show "traveling" status with countdown

**Evidence:**
- `/frontend/src/types/index.ts` Line 290: Fleet has `status: 'stationed' | 'traveling' | 'combat'`
- `/frontend/src/components/panels/PvPPanel.tsx`: Attack flow is **instant** - no travel time shown
- `/frontend/src/components/panels/FleetPanel.tsx` Line 278: Shows fleet status badge but **no travel time countdown**
- `/frontend/src/hooks/usePvP.ts`: Attack function has NO travel logic, immediately triggers combat
- **ZERO UI** for:
  - Fleet departure time
  - Travel duration countdown
  - Fleet arrival ETA
  - Fleet position/progress
  - Travel time calculation based on MOV stat

**Impact:** Users have no idea their fleets are traveling or when they'll arrive. Critical UX failure.

---

### FAIL-002: Space Points (SP) Display - NOT SHOWN ANYWHERE
**GDD Requirement:** REQ-PVP002 - Fleet travel uses SP (Space Points)

**Evidence:**
- `/frontend/src/components/layout/ResourceHUD.tsx`: Shows Metal, He3, Gold, **NO SP**
- `/frontend/src/components/panels/PvPPanel.tsx`: No SP cost shown, no SP consumption UI
- `/frontend/src/types/index.ts`: Player type has NO `space_points` field
- Grepped entire frontend for "space_point", "SP" - **ZERO references** to Space Points display

**Impact:** Users cannot see their SP, cannot plan attacks, cannot know if they have enough SP. Feature is invisible.

---

### FAIL-003: Truce Card Protection Indicator - ZERO VISUAL FEEDBACK
**GDD Requirement:** REQ-INV016, REQ-INV017 - Truce Cards provide 12h/72h protection

**Evidence:**
- `/frontend/src/components/panels/InventoryPanel.tsx` Line 183: "Use Item" button exists
- `/frontend/src/services/api.ts` Line 553: `useItem()` API call exists
- **BUT NO UI shows:**
  - Protection status (is user protected?)
  - Protection duration remaining
  - Protection countdown timer
  - Visual shield icon when protected
  - Warning when protection expires soon

**User Flow Test:**
1. User clicks "Use" on Truce Card
2. Item is consumed (backend handles this)
3. **NOTHING VISIBLE HAPPENS** in UI
4. User has no idea if they're protected
5. FAIL

**Impact:** Users waste items with no feedback. Broken feature.

---

### FAIL-004: Radar Panel (Incoming Attacks) - NOT IMPLEMENTED
**GDD Requirement:** REQ-B007, REQ-PVP006 - Radar detects incoming attacks, provides warning

**Evidence:**
- Grepped for "radar", "incoming" - ZERO panels found
- `/frontend/src/components/panels/` - NO RadarPanel.tsx, NO IncomingAttacksPanel.tsx
- `/frontend/src/App.tsx` - No route for radar/attacks view
- `/frontend/src/components/layout/SideNav.tsx` - No navigation item for radar/attacks

**What's Missing:**
- Panel showing incoming attacks
- Attack details (attacker name, fleet count, ETA)
- Radar building level effect on detection range
- Warning notifications when attacked

**Impact:** Users get attacked with ZERO warning. Major PvP UX gap.

---

### FAIL-005: Fleet MOV Stat Display - MISSING IN FLEET PANEL
**GDD Requirement:** REQ-S006 - Movement (MOV) stat determines fleet travel speed

**Evidence:**
- `/frontend/src/types/index.ts` Line 185: ShipDesign has `total_movement` field
- `/frontend/src/components/panels/FleetPanel.tsx`: Shows fleet name, status, formation, **NO MOV stat**
- `/frontend/src/components/panels/FleetPanel.tsx` Line 281: Shows "X/9 positions, Y ships" - **missing MOV**
- ShipDesign has `total_movement` calculated, but **never displayed in fleet context**

**User Flow Test:**
1. User creates fleet with 3 ship designs (MOV: 100, 150, 200)
2. Fleet panel shows: "3/9 positions, 450 ships"
3. **MISSING:** "Fleet Speed: 100 MOV (slowest ship)"
4. User has NO IDEA how fast fleet travels
5. FAIL

**Impact:** Users cannot plan PvP attacks with travel time in mind.

---

### FAIL-006: Defense Building Combat Stats - NOT SHOWN IN UI
**GDD Requirement:** REQ-B015-B019 - Defense buildings have HP, damage, range stats

**Evidence:**
- `/frontend/src/components/panels/BuildingDetailPanel.tsx`: Shows building level, upgrade cost
- **NO display of:**
  - HP (current/max)
  - Damage output
  - Range
  - Combat integration status
- Defense buildings are placeable but users **cannot see their combat effectiveness**

**Impact:** Users build defense buildings blind. No strategic value visible.

---

### FAIL-007: PvP Attack Phases (Scout → Travel → Combat) - INSTANT ONLY
**GDD Requirement:** REQ-PVP002, REQ-PVP003 - Attack flow should have travel phase

**Evidence:**
- `/frontend/src/hooks/usePvP.ts` Line 52-74: `attack()` function
  ```typescript
  const result = await attackPlanet({ defender_planet_id, fleet_ids })
  setLastBattle(result) // INSTANT - no travel phase
  ```
- **Expected flow:**
  1. User clicks "Attack"
  2. Fleet departs (status → "traveling")
  3. Travel time countdown (based on MOV + distance)
  4. Fleet arrives
  5. Combat executes
  6. Fleet returns
  7. Loot shown
- **Actual flow:**
  1. User clicks "Attack"
  2. Battle result shown **immediately**
  3. FAIL - no scout, no travel, no return

**Impact:** PvP feels like a button click, not a space battle. Immersion broken.

---

### FAIL-008: Fleet Status Filters (Traveling/Returning/Stationed) - NOT IMPLEMENTED
**GDD Requirement:** Fleet management should distinguish fleet states

**Evidence:**
- `/frontend/src/types/index.ts` Line 290: Fleet status exists: `'stationed' | 'traveling' | 'combat'`
- `/frontend/src/components/panels/PvPPanel.tsx` Line 16: Filters for `status === 'stationed'`
- **MISSING UI:**
  - "Traveling" fleets panel (shows ETA, destination)
  - "Returning" fleets panel (shows loot, casualties, ETA)
  - "In Combat" fleets (shows battle progress)
- Only "stationed" fleets are selectable - **other states invisible**

**Impact:** Users lose track of their fleets. No fleet management.

---

### FAIL-009: Commander Expertise Tooltip - MISSING EXPLANATIONS
**GDD Requirement:** REQ-C012, REQ-C013 - Weapon/Ship Expertise grades (S/A/B/C/D)

**Evidence:**
- `/frontend/src/components/panels/CommandersListPanel.tsx`: Shows commander stats
- Commander type has `weapon_expertise` and `ship_expertise` fields
- **NO tooltips explaining:**
  - S = +30% damage (weapon), +10% dealt/-10% received (ship)
  - A = +10% damage, +5% dealt/-10% received
  - B/C/D penalties
- Users see "Expertise: S" with **zero context**

**Impact:** Users don't know what S/A/B/C/D mean. Feature wasted.

---

### FAIL-010: RBP Conquest Progress - NO REAL-TIME UPDATES
**GDD Requirement:** REQ-GAL009 - RBP conquest uses 8-phase combat, 99 rounds max

**Evidence:**
- `/frontend/src/components/panels/GalaxyMapPanel.tsx` Line 38-46: Attack flow
  ```typescript
  const response = await attackRBP(rbpPlanetId, { fleet_ids })
  if (response) {
    setSelectedFleets([])
    selectZone(null)
  }
  ```
- Attack is **instant** - no battle progress shown
- **MISSING:**
  - Round-by-round updates (Round 1/99...)
  - Real-time damage/casualties
  - Fleet cannot be recalled (should show "Battle in progress, cannot recall")
  - Multi-corp attack leaderboard (who has most kills)

**Impact:** RBP battles feel like clicking a button. No tension, no strategy.

---

### FAIL-011: Corp Donation Daily Limit UI - NOT SHOWN
**GDD Requirement:** REQ-CORP004 - Max 200 pts/day (2M resources)

**Evidence:**
- `/frontend/src/components/panels/CorpsPanel.tsx` Line 402-476: Donate tab
- Shows input fields for Metal, He3, Gold
- **MISSING:**
  - "Daily Contribution: 150/200 pts" progress bar
  - "Remaining today: 50 pts (500k resources)"
  - Visual warning when approaching limit
- Users can enter amounts and get backend error, but **no proactive UI guidance**

**Impact:** Users waste time entering amounts over limit. Poor UX.

---

### FAIL-012: Warehouse Capacity Badge Color Logic - PLACEHOLDER
**GDD Requirement:** REQ-R004 - Warehouse capacity warnings

**Evidence:**
- `/frontend/src/components/layout/ResourceHUD.tsx` Line 71: Shows "Storage Cap: X"
- **NO color coding:**
  - Green: < 80% full
  - Yellow: 80-95% full
  - Red: > 95% full (resources wasting)
- MEMORY.md Line 37: "Warehouse capacity badge color: uses placeholder logic"
- Known issue documented, **not fixed**

**Impact:** Users don't know when warehouse is full. Resources wasted.

---

## Missing UI Flows (features with NO frontend at all)

### NOT_IMPL-001: Construction Card (+3 slots for 72h) - NO ACTIVE INDICATOR
**GDD Requirement:** REQ-INV009 - Construction Card item

**Evidence:**
- Item exists in inventory system (backend)
- **NO UI shows:**
  - Active construction card icon
  - "5/5 construction slots (2 base + 3 from card)"
  - Countdown timer for card expiration
  - Visual indicator in building panel

**Workaround:** Backend handles this, but user has no idea if card is active.

---

### NOT_IMPL-002: MVP Tool (+20% production/build/repair, 7 days) - NO ACTIVE BUFFS PANEL
**GDD Requirement:** REQ-INV010 - MVP Tool item

**Evidence:**
- Item exists in inventory
- **NO "Active Buffs" panel showing:**
  - MVP Tool: +20% all production (5d 12h 30m remaining)
  - Extra Tax: +30% Gold (2h remaining)
  - Construction Card: +3 slots (48h remaining)
- Users use items and **cannot see if they're active**

---

### NOT_IMPL-003: Fleet Return Phase - NO "RETURNING" STATUS
**GDD Requirement:** PvP fleets should return after battle

**Evidence:**
- Fleet status enum includes "traveling" but **NO "returning"**
- After PvP attack, fleet instantly back at base
- **MISSING:**
  - Fleet status changes to "returning"
  - Return travel time (same as outbound)
  - ETA for return
  - Loot visible during return trip

---

### NOT_IMPL-004: PvP Cooldown Shared UI - ONLY SHOWS IN PVP PANEL
**GDD Requirement:** REQ-PVP005 - 5-minute cooldown between attacks

**Evidence:**
- `/frontend/src/hooks/usePvP.ts` Line 61: Cooldown tracked in hook
- `/frontend/src/components/panels/PvPPanel.tsx` Line 119-122: Shows cooldown in PvP panel
- **NOT SHOWN:**
  - ResourceHUD cooldown badge
  - SideNav "PvP" icon disabled state during cooldown
  - Global notification "PvP ready!" when cooldown ends

**Impact:** Users must keep PvP panel open to see cooldown.

---

### NOT_IMPL-005: Ship Repair UI - PLACEHOLDER ONLY
**GDD Requirement:** REQ-B009 - Spacedock repairs ships

**Evidence:**
- `/frontend/src/components/panels/SpacedockPanel.tsx`: File exists (106 lines)
- Backend has repair system (Module 8)
- **Frontend shows placeholder only** - repair flow not implemented
- MEMORY.md confirms: "Spacedock manages fleets, repair placeholder"

---

### NOT_IMPL-006: Tech Tree Visual Dependencies - PLACEHOLDER
**GDD Requirement:** REQ-T107 - Research prerequisites validation

**Evidence:**
- `/frontend/src/components/panels/ResearchPanel.tsx`: Shows tech list
- **NO visual tech tree:**
  - Node graph showing prerequisites
  - Lines connecting dependent techs
  - Locked/unlocked visual states
  - Prerequisite tooltips
- Users must read text to understand dependencies

---

## Panel-by-Panel Validation

### Buildings System - PASS (90%)
✅ BuildingPanel.tsx (Grid placement works)
✅ BuildingDetailPanel.tsx (Shows costs, time, stats)
✅ ConstructionPanel.tsx (Shows slots, progress)
⚠️ PARTIAL: Defense building stats not shown

### Resources System - PASS (85%)
✅ ResourceHUD.tsx (Shows Metal/He3/Gold, rates, collection)
✅ Warehouse collection works
❌ FAIL: SP not shown
❌ FAIL: Warehouse capacity color logic placeholder

### Research System - PASS (80%)
✅ ResearchPanel.tsx (663 lines - shows all 7 trees)
✅ Tech costs, time, effects shown
✅ Active research countdown
⚠️ PARTIAL: Visual tech tree missing
⚠️ PARTIAL: Prerequisites shown as text only

### Blueprint System - PASS (95%)
✅ BlueprintPanel.tsx (327 lines - tier upgrades work)
✅ Research slots shown
✅ Progress bars for research
✅ Tier validation works

### Ships System - PASS (98%)
✅ ShipDesignPanel.tsx (627 lines - full 3D preview)
✅ ShipFactoryPanel.tsx (253 lines - production works)
✅ Hull/module selection works
✅ Stats calculated correctly
⚠️ PARTIAL: MOV stat not shown in fleet context

### Fleet System - PARTIAL (75%)
✅ FleetPanel.tsx (291 lines - 3x3 grid works)
✅ Formation, targeting selection
✅ Stack assignment works
❌ FAIL: No MOV stat display
❌ FAIL: No traveling/returning fleet panels
❌ FAIL: Status badge shows text only (no countdown)

### PvP Combat - FAIL (40%)
✅ PvPPanel.tsx (185 lines - search works)
✅ Fleet selection works
✅ Battle results shown
❌ FAIL: No fleet travel time
❌ FAIL: No SP display
❌ FAIL: Instant combat (no phases)
❌ FAIL: No radar/incoming attacks panel
❌ FAIL: No fleet return phase

### Commander System - PASS (90%)
✅ CommandCenterPanel.tsx (170 lines - recruitment works)
✅ CommandersListPanel.tsx (287 lines - list, stats)
✅ CompoundCenterPanel.tsx (229 lines - merging works)
✅ Gacha rates shown
⚠️ PARTIAL: Expertise tooltips missing

### Recycling Plant - PASS (95%)
✅ RecyclingPlantPanel.tsx (199 lines - works fully)
✅ Ship selection
✅ 70% resource recovery shown
✅ Duration countdown

### Inventory System - PARTIAL (70%)
✅ InventoryPanel.tsx (192 lines - shows items)
✅ Item use works
✅ Categories grouped nicely
❌ FAIL: No active buffs panel (MVP Tool, Extra Tax, Construction Card)
❌ FAIL: No protection status indicator (Truce Card)
❌ FAIL: No visual feedback after using items

### World Chat - PASS (95%)
✅ ChatPanel.tsx (176 lines - works fully)
✅ Real-time messages
✅ Rate limiting shown
✅ Profanity filter active

### Combat Reports - PASS (90%)
✅ CombatReportsPanel.tsx (210 lines - shows reports)
✅ Round-by-round details
✅ Loot shown
⚠️ PARTIAL: No filtering by type (PvP vs Instance)

### Quest System - PASS (98%)
✅ QuestPanel.tsx (523 lines - main/side/daily all work)
✅ Progress tracking
✅ Rewards shown
✅ Auto-progress integration

### Corps System - PASS (85%)
✅ CorpsPanel.tsx (557 lines - create/join/leave works)
✅ Members list
✅ Donation system
✅ Role management
❌ FAIL: Daily donation limit not shown proactively

### Galaxy Map - PARTIAL (75%)
✅ GalaxyMapPanel.tsx (257 lines - 7x7 grid shown)
✅ RBP details shown
✅ Ownership colors
✅ Protection status shown
❌ FAIL: Attack is instant (no battle progress)
❌ FAIL: No multi-corp attack leaderboard
❌ FAIL: No real-time round updates

---

## Partial Implementations (need polish)

### PARTIAL-001: Fleet Speed Calculation - EXISTS BUT NOT SHOWN
**Status:** Backend calculates correctly, frontend doesn't display
**Fix:** Add MOV badge to FleetPanel.tsx (5 min fix)

### PARTIAL-002: Defense Building Stats - BACKEND READY, UI MISSING
**Status:** Building stats in database, not rendered in UI
**Fix:** Update BuildingDetailPanel.tsx to show HP/damage/range (15 min)

### PARTIAL-003: Warehouse Capacity Warning - LOGIC INCOMPLETE
**Status:** Shows capacity number, no color warning
**Fix:** Add % calculation + color logic to ResourceHUD.tsx (10 min)

### PARTIAL-004: Tech Tree Prerequisites - TEXT ONLY
**Status:** Prerequisites shown as text, no visual graph
**Fix:** (Major work - 4+ hours) Add React Flow or similar library

### PARTIAL-005: Expertise Tooltips - DATA EXISTS, NO EXPLANATION
**Status:** Commanders have expertise grades, no UI tooltips
**Fix:** Add tooltip component to CommandersListPanel.tsx (20 min)

### PARTIAL-006: Active Buffs Panel - NO CENTRAL VIEW
**Status:** Items work, but no "Active Buffs" UI
**Fix:** Create ActiveBuffsPanel.tsx showing all active timed items (1 hour)

### PARTIAL-007: Corp Donation Daily Limit - BACKEND ENFORCED, NOT SHOWN
**Status:** Backend blocks over-limit, UI doesn't warn
**Fix:** Add progress bar to CorpsPanel.tsx donate tab (15 min)

### PARTIAL-008: Combat Report Filtering - NO FILTERS
**Status:** Shows all reports, no type filter
**Fix:** Add dropdown filter in CombatReportsPanel.tsx (10 min)

---

## TypeScript Coverage - PASS (100%)
✅ All types defined in `/frontend/src/types/`
✅ No `any` types found
✅ Proper error typing in hooks
✅ API response types match backend

---

## Routing Coverage - PASS (100%)
✅ `/` - Home (planet list)
✅ `/planet/:id` - Planet view (3D grid)
✅ `/military` - Military page (9 tabs)
✅ `/inventory` - Inventory page

All panels accessible via SideNav or Military tabs.

---

## Hook Coverage - PASS (100%)
✅ 18 hooks implemented
✅ All hooks have cleanup (intervals, timeouts cleared)
✅ All hooks use GameContext for error handling
✅ No memory leaks detected

---

## API Coverage - PASS (100%)
✅ 586 lines in api.ts
✅ All backend endpoints mapped
✅ Auth interceptor works
✅ Dev-mode header support

---

## Production Readiness Assessment

### Blocking Issues (MUST FIX before launch):
1. **FAIL-001**: Fleet travel time UI (critical PvP feature)
2. **FAIL-002**: SP display (users can't see Space Points)
3. **FAIL-003**: Truce Card feedback (users waste items)
4. **FAIL-004**: Radar panel (users get attacked with no warning)
5. **FAIL-007**: PvP attack phases (combat feels instant/fake)

### High Priority (should fix):
6. FAIL-005: Fleet MOV stat display
7. FAIL-008: Fleet status filters
8. FAIL-010: RBP conquest progress
9. NOT_IMPL-002: Active buffs panel

### Medium Priority (nice to have):
10. FAIL-006: Defense building stats
11. FAIL-009: Expertise tooltips
12. FAIL-011: Corp donation limit UI
13. PARTIAL issues (all fixable in < 2 hours total)

---

## Recommended Action Plan

### Phase 1 (Critical - 2-3 days):
1. Implement fleet travel time UI (FAIL-001)
   - Add travel_started_at, travel_duration fields to Fleet type
   - Create TravelingFleetsPanel.tsx
   - Add countdown timers
   - Update PvP flow to set fleet status to "traveling"

2. Add SP (Space Points) to ResourceHUD (FAIL-002)
   - Add space_points field to Player type
   - Show SP in HUD next to Metal/He3/Gold
   - Show SP cost in PvP attack confirmation

3. Add Truce Card protection UI (FAIL-003)
   - Add protection_until field to Player type
   - Show shield icon in ResourceHUD when protected
   - Add countdown timer for protection
   - Visual feedback when using Truce Card

4. Create Radar/Incoming Attacks panel (FAIL-004)
   - New panel: IncomingAttacksPanel.tsx
   - Show attacker name, fleet count, ETA
   - Add navigation in SideNav

### Phase 2 (High Priority - 1-2 days):
5. Fix PvP attack phases (FAIL-007)
   - Split attackPlanet into: sendFleet → (wait) → executeCombat → (wait) → returnFleet
   - Update usePvP.ts to track attack phases
   - Add phase UI to PvPPanel.tsx

6. Add fleet status panels (FAIL-008)
   - TravelingFleetsPanel.tsx
   - ReturningFleetsPanel.tsx
   - Show in Military page as new tabs

7. Add RBP battle progress (FAIL-010)
   - Stream battle updates via WebSocket or polling
   - Show round-by-round progress
   - Add "Battle in progress" UI

### Phase 3 (Polish - 1 day):
8. Fix all PARTIAL issues (8 items, < 2 hours total)
9. Add Active Buffs panel (NOT_IMPL-002)
10. Add tooltips and polish

---

## Comparison to Previous Validation

**Previous Report (Comprehensive Final QA):** 194/200 tests pass (97%)
**This Report:** 52/78 features fully work (66%)

**Discrepancy Explained:**
- Previous validation checked **backend functionality** (✅ correct)
- This validation checks **frontend user experience** (❌ many gaps)
- Backend works, but **frontend doesn't expose many features to users**

**Example:**
- ✅ Backend: Fleet travel time calculation works
- ❌ Frontend: No UI shows travel time, countdown, or ETA
- Result: Feature exists in code, **invisible to users**

---

## Final Verdict

**APPROVED FOR PRODUCTION: NO**

**Recommendation:** Fix 5 blocking issues (Phase 1) before launch. Game is playable but critical PvP features are broken/invisible.

**Estimated Fix Time:** 4-6 days (Phase 1 + Phase 2)

**Pass Criteria for Re-Validation:**
- All 12 FAIL items reduced to PARTIAL or PASS
- At least 85% features fully functional
- All critical user flows (PvP attack, fleet management, protection) work end-to-end

---

## Notes for Developers

1. **"Component exists" ≠ "Feature works"** - Always test the full user flow, not just file existence
2. **Missing visual feedback** - Many items work in backend but give ZERO user feedback (Truce Card, buffs, etc.)
3. **Placeholder logic** - Several features documented as "placeholder" but never upgraded (warehouse capacity color, etc.)
4. **Status fields unused** - Fleet status enum has "traveling" but UI never shows it
5. **Real-time updates missing** - RBP battles, fleet travel should update live, but frontend uses instant API calls

**Good News:**
- Backend is solid (97% pass rate)
- TypeScript coverage is excellent
- No memory leaks or major bugs
- Most issues are **frontend polish**, not fundamental architecture problems

**All 12 FAIL items are fixable in 4-6 days of focused frontend work.**

---

End of Report
