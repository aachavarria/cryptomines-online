# Blueprint Panel Research UI - Implementation Summary

**Task ID:** #23
**Date:** 2026-02-07
**Developer:** frontend-dev

## Overview
Added comprehensive blueprint research functionality to the BlueprintPanel, including research level stars, research button, active research timer, and visual state indicators.

## Changes Made

### 1. Backend: Active Blueprint Research Endpoint
**File:** `/backend/internal/handlers/blueprints.go` (lines 272-318)

- **New Endpoint:** `GET /api/blueprint-research/active`
- **Function:** `GetActiveBlueprintResearch`
- **Returns:** Active blueprint research with countdown timer, or null if none active
- **Response:**
  ```json
  {
    "id": "uuid",
    "player_blueprint_id": "uuid",
    "target_level": 2,
    "is_researching": true,
    "research_finish_at": "2026-02-07T15:30:00Z",
    "blueprint_id": 123,
    "blueprint_name": "frigate_hull_t1",
    "blueprint_type": "hull"
  }
  ```

**File:** `/backend/cmd/server/main.go` (line 70)
- Registered route: `GET /api/blueprint-research/active`

### 2. Frontend: Type Definitions
**File:** `/frontend/src/types/index.ts` (lines 237-275)

- **Updated `Blueprint` interface:**
  - Added `display_name: string`
  - Added `hull_class?: string`
  - Added `module_category?: string`

- **Updated `PlayerBlueprint` interface:**
  - Added `is_researching?: boolean`
  - Added `research_finish_at?: string`

- **New `ActiveBlueprintResearch` interface:**
  ```typescript
  {
    id: string
    player_blueprint_id: string
    target_level: number
    is_researching: boolean
    research_finish_at: string
    blueprint_id: number
    blueprint_name: string
    blueprint_type: 'hull' | 'module'
  }
  ```

### 3. Frontend: API Functions
**File:** `/frontend/src/services/api.ts` (lines 19, 224-227)

- Added `ActiveBlueprintResearch` import
- Added `getActiveBlueprintResearch()` function
  - Calls `GET /api/blueprint-research/active`
  - Returns active research or null

### 4. Frontend: useBlueprints Hook
**File:** `/frontend/src/hooks/useBlueprints.ts`

**New State:**
- `activeResearch: ActiveBlueprintResearch | null`

**Updated Functions:**
- `refresh()` - Now fetches active research in parallel with blueprints
- `research(id)` - New function to start blueprint research

**Updated Return:**
- Added `activeResearch` to hook return
- Added `research` function to hook return

### 5. Frontend: BlueprintPanel Component
**File:** `/frontend/src/components/panels/BlueprintPanel.tsx` (completely rewritten)

**New Features:**

1. **Research Level Stars Display (★★★)**
   - Shows 3 stars for each blueprint
   - Filled stars (gold, glowing) = researched levels
   - Empty stars (gray, dim) = locked levels
   - Example: Lv2 blueprint shows ★★☆

2. **Research Button**
   - Shows "Research Lv2" or "Research Lv3" button
   - Only visible for activated blueprints with research_level < 3
   - Disabled when WRC is busy with another research
   - Opens confirmation modal on click

3. **Active Research Timer**
   - Shows active blueprint research at top of panel
   - Real-time countdown timer
   - Progress bar (1 hour per level)
   - Auto-refreshes when research completes (polls every 2 seconds)

4. **Visual States:**
   - **Locked:** Gray card + lock icon + "🔒 Not Owned"
   - **Owned:** Yellow border + "Activate" button
   - **Activated:** Green border + research stars + research button
   - **Researching:** Blue pulsing border + "Researching..." status
   - **Max Research (Lv3):** Gold border + glow + "Max Research" status
   - **WRC Busy:** "WRC Busy" status when another research is active

5. **Research Confirmation Modal:**
   - Shows current level → target level
   - Displays costs: Metal, He3, Gold
   - Shows research time (1 hour per level)
   - Confirmation buttons: CANCEL / START RESEARCH

**Components:**
- `BlueprintPanel` - Main component with grid layout
- `ActiveResearchBar` - Active research timer bar (at top)
- `ResearchConfirmModal` - Research confirmation dialog

### 6. Frontend: CSS Styles
**File:** `/frontend/src/styles/phase2.css` (lines 1370-1622)

**New Styles:**
- `.bp-status-maxed` - Gold text for max research
- `.bp-status-researching` - Blue text for researching
- `.bp-status-busy` - Warning text for busy WRC
- `.bp-card--researching` - Blue pulsing border animation
- `.bp-card--maxed` - Gold border + glow
- `.bp-research-stars` - Star container
- `.bp-star.filled` - Gold glowing star
- `.bp-star.empty` - Gray dim star
- `.p2-btn-research` - Blue gradient research button
- `.bp-active-research` - Active research bar container
- `.bp-active-info` - Research info row
- `.bp-active-name` - Blueprint name
- `.bp-active-level` - Level indicator (→ Lv2)
- `.bp-active-timer` - Countdown timer
- `.bp-progress-bar` - Progress bar container
- `.bp-progress-fill` - Progress bar fill
- `.bp-confirm-backdrop` - Modal backdrop
- `.bp-confirm-modal` - Modal container
- `.bp-confirm-header` - Modal header
- `.bp-confirm-body` - Modal body
- `.bp-confirm-row` - Cost row
- `.bp-confirm-label` - Cost label
- `.bp-confirm-value` - Cost value
- `.bp-confirm-sep` - Separator line
- `.bp-confirm-note` - Info note
- `.bp-confirm-actions` - Button container
- `.bp-confirm-cancel` - Cancel button
- `.bp-confirm-ok` - Confirm button

**Animations:**
- `@keyframes bpResearchPulse` - Blue pulsing glow for researching state

## Technical Details

### Research Costs (from backend)
```go
targetLevel := playerBp.ResearchLevel + 1
baseCost := 10000 * targetLevel
metalCost := baseCost
he3Cost := baseCost * 3 / 4
goldCost := baseCost / 2
```

**Example:**
- Lv1 → Lv2: 10k Metal, 7.5k He3, 5k Gold, 2 hours
- Lv2 → Lv3: 20k Metal, 15k He3, 10k Gold, 3 hours

### Research Time
- 1 hour per level (3600 seconds * target_level)
- Lv1 → Lv2: 2 hours
- Lv2 → Lv3: 3 hours

### Backend Requirements
- Weapon Research Center (WRC) level ≥ 1
- Only 1 active research at a time per WRC
- Backend auto-complete worker (Task #21) handles completion

### Frontend Polling
- Polls active research status every 2 seconds when research is active
- Auto-refreshes blueprint list when research completes
- No manual refresh button needed

## Data Flow

1. **Page Load:**
   - `useBlueprints()` fetches blueprints + active research
   - Sets `allBlueprints`, `myBlueprints`, `activeResearch`

2. **Start Research:**
   - User clicks "Research Lv2" button
   - Confirmation modal shows costs + time
   - User confirms → `POST /api/blueprints/{id}/research`
   - Backend creates `blueprint_research` record
   - Frontend refreshes → shows active research bar

3. **Research Progress:**
   - ActiveResearchBar shows countdown timer
   - Progress bar animates based on elapsed time
   - Polls every 2 seconds for completion

4. **Research Complete:**
   - Backend worker auto-completes research
   - Updates `player_blueprints.research_level`
   - Deletes `blueprint_research` record
   - Frontend detects completion (poll)
   - Refreshes blueprint list
   - Active research bar disappears
   - Stars update to show new level

## Testing Recommendations

### Manual Testing

1. **Research Level Stars:**
   - View owned blueprints → verify 3 stars displayed
   - Check star states: Lv0 = ☆☆☆, Lv1 = ★☆☆, Lv2 = ★★☆, Lv3 = ★★★
   - Verify gold glow on filled stars

2. **Research Button:**
   - Owned but not activated → no research button
   - Activated Lv0/1/2 → "Research Lv1/2/3" button shows
   - Lv3 → "Max Research" status, no button
   - WRC busy → "WRC Busy" status, no button

3. **Research Confirmation Modal:**
   - Click research button → modal opens
   - Verify costs: Metal, He3, Gold
   - Verify time: 2h for Lv2, 3h for Lv3
   - Cancel → modal closes
   - Confirm → research starts

4. **Active Research Bar:**
   - Start research → bar appears at top
   - Verify blueprint name + target level
   - Verify countdown timer updates
   - Verify progress bar animates
   - Wait for completion → bar disappears

5. **Visual States:**
   - Locked: Gray card + lock icon
   - Owned: Yellow border
   - Activated: Green border + stars
   - Researching: Blue pulsing border
   - Max Research: Gold border + glow

6. **Edge Cases:**
   - Start research with insufficient resources → error message
   - Try to research while WRC busy → button disabled
   - Try to research Lv3 blueprint → "Max Research" status
   - Close modal while pending → no action taken

### Browser Compatibility
- Chrome/Edge: Primary target
- Firefox: Test pulsing animations
- Safari: Test modal backdrop blur

## Future Enhancements

1. **Cancel Research Button:** Allow canceling active research (partial refund)
2. **Speedup Research Button:** Use vouchers to accelerate research
3. **Research History:** Log of completed research
4. **Bulk Research Queue:** Research multiple blueprints in sequence
5. **WRC Level Indicators:** Show WRC level + research slots
6. **Research Notifications:** Browser notifications when research completes

## Files Modified

### Backend
- `/backend/internal/handlers/blueprints.go` - New endpoint
- `/backend/cmd/server/main.go` - Route registration

### Frontend
- `/frontend/src/types/index.ts` - Type definitions
- `/frontend/src/services/api.ts` - API function
- `/frontend/src/hooks/useBlueprints.ts` - Hook logic
- `/frontend/src/components/panels/BlueprintPanel.tsx` - UI component
- `/frontend/src/styles/phase2.css` - Styles + animations

## Dependencies
- No new dependencies added
- Uses existing hooks: `useBlueprints`, `useCountdown`, `formatDuration`, `formatNumber`
- Uses existing API client: `axios`

## Estimated Completion Time
**Actual:** 3 hours (within 3-4 hour estimate)

## Status
✅ **COMPLETED** - Ready for QA testing
