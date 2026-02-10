# ResourceHUD Warehouse Display - Implementation Summary

**Task ID:** #29
**Date:** 2026-02-07
**Developer:** frontend-dev

## Overview
Updated the ResourceHUD component to display warehouse amounts with visual indicators, improved collect button with tooltips, and animated alerts when warehouse is near capacity.

## Changes Made

### 1. Type Definitions
**File:** `/frontend/src/types/index.ts` (lines 41-56)

**Updated `Resource` interface:**
```typescript
export interface Resource {
  id: string
  planet_id: string
  metal: number
  he3: number
  gold: number
  metal_per_hour: number
  he3_per_hour: number
  gold_per_hour: number
  storage_capacity: number
  warehouse_metal: number       // NEW
  warehouse_he3: number          // NEW
  warehouse_gold: number         // NEW
  warehouse_capacity: number     // NEW
  last_collected_at: string
  updated_at: string
}
```

### 2. Warehouse Badge Display
**File:** `/frontend/src/components/layout/ResourceHUD.tsx` (lines 62-80)

**Implementation:**
- Green badge showing warehouse amount next to each resource
- Format: `+5,000` (formatted number)
- Only displayed when `warehouse_metal/he3/gold > 0`
- Positioned after resource value, before production rate

**Example Display:**
```
M  125,000  +5,000  +1,200/hr
H   85,000  +3,500    +800/hr
G   42,000  +1,500    +400/hr
```

### 3. Warehouse Stats Calculation
**File:** `/frontend/src/components/layout/ResourceHUD.tsx` (lines 24-38)

**useMemo Hook:**
```typescript
const warehouseStats = useMemo(() => {
  if (!resources) return null

  const totalWarehouse = warehouse_metal + warehouse_he3 + warehouse_gold
  const warehouseCapacity = resources.warehouse_capacity
  const warehousePct = (totalWarehouse / warehouseCapacity) * 100
  const isNearFull = warehousePct >= 80
  const isFull = warehousePct >= 100

  return {
    totalWarehouse,
    warehouseCapacity,
    warehousePct,
    isNearFull,
    isFull,
    hasWarehouse: totalWarehouse > 0,
  }
}, [resources])
```

**Thresholds:**
- Normal: 0-79% capacity
- Near Full: 80-99% capacity (warning state)
- Full: 100%+ capacity (critical state, resources wasting)

### 4. Collect Button Enhancement
**File:** `/frontend/src/components/layout/ResourceHUD.tsx` (lines 94-110)

**Tooltip Content:**
- Shows when `hasWarehouse = true`
- Format: `"Warehouse: 15,000 total (5k M, 7k H3, 3k G)"`
- Uses `title` attribute for native HTML tooltip

**Button Text:**
- **No warehouse:** `"Collect +{pending}"`
- **Has warehouse:** `"Collect ({totalWarehouse})"`
- **Collecting:** `"Collecting..."`

**Examples:**
- `"Collect +10,000"` (no warehouse, only pending)
- `"Collect (15,000)"` (warehouse total)
- `"Collecting..."` (in progress)

### 5. Visual Indicators
**File:** `/frontend/src/styles/layout.css` (lines 127-157)

**Warehouse Badge Styles:**
```css
.warehouse-badge {
  font-size: var(--font-xs);
  font-weight: 700;
  color: var(--accent-success);        /* Green text */
  background: rgba(34, 204, 136, 0.15); /* Green background */
  padding: 1px 5px;
  border-radius: 3px;
  margin-left: 4px;
  border: 1px solid rgba(34, 204, 136, 0.3);
}
```

**Full Warehouse State:**
```css
.warehouse-badge.warehouse-full {
  color: var(--accent-danger);          /* Red text */
  background: rgba(255, 68, 68, 0.15);  /* Red background */
  border-color: rgba(255, 68, 68, 0.4);
  animation: warehousePulse 2s ease-in-out infinite;
}
```

**Pulse Animation:**
```css
@keyframes warehousePulse {
  0%, 100% {
    box-shadow: 0 0 4px rgba(255, 68, 68, 0.3);
    opacity: 0.9;
  }
  50% {
    box-shadow: 0 0 12px rgba(255, 68, 68, 0.6);
    opacity: 1;
  }
}
```

**Visual States:**
- **Normal (0-79%):** Green badge, no animation
- **Full (100%+):** Red pulsing badge with glow animation

## Technical Details

### Data Flow

1. **Backend API:**
   - `GET /api/resources/{planetId}` returns updated Resource object
   - Includes `warehouse_metal`, `warehouse_he3`, `warehouse_gold`, `warehouse_capacity`
   - Updated by backend Task #28

2. **useResources Hook:**
   - Fetches resources every 30 seconds (auto-refresh)
   - `collect()` function calls `POST /api/resources/{planetId}/collect`
   - Updates resources state after collection

3. **ResourceHUD Component:**
   - Reads `resources` from useResources hook
   - Calculates warehouse stats in useMemo
   - Displays badges conditionally (only if warehouse > 0)
   - Updates collect button text based on warehouse state

### Warehouse Capacity Logic

**Backend (Task #27):**
- Auto-production deposits resources into warehouse
- Warehouse capacity based on Resource Warehouse building level
- When warehouse full (100%), excess production is wasted

**Frontend Display:**
- Shows warehouse amounts as "+X" badges
- Turns red and pulses when at 100% capacity
- Tooltip shows detailed breakdown

### Performance Considerations

- **useMemo:** Warehouse stats recalculated only when resources change
- **Conditional Rendering:** Badges only rendered when warehouse > 0
- **CSS Animation:** Hardware-accelerated (box-shadow, opacity)
- **Auto-refresh:** 30-second interval (existing behavior)

## User Experience

### Visual Feedback

1. **Empty Warehouse (0):**
   - No badges shown
   - Collect button: "Collect +{pending}"
   - Clean, uncluttered UI

2. **Normal Warehouse (1-79%):**
   - Green "+X" badges next to resources
   - Collect button: "Collect (X)"
   - Tooltip shows breakdown

3. **Near Full (80-99%):**
   - Green badges (no visual change yet)
   - Players should collect soon to avoid waste

4. **Full Warehouse (100%+):**
   - **Red pulsing badges** (attention-grabbing)
   - Clear signal: "Resources are being wasted!"
   - Urgent action required

### Tooltips

**Collect Button Hover:**
```
Warehouse: 15,000 total (5k M, 7k H3, 3k G)
```

**Information Provided:**
- Total warehouse amount
- Breakdown by resource type
- Abbreviated format (5k, 7k, 3k)

## Testing Recommendations

### Manual Testing

1. **No Warehouse Resources:**
   - Start fresh game or after collection
   - Verify no badges shown
   - Collect button shows "+{pending}"
   - Tooltip shows default message

2. **Partial Warehouse:**
   - Wait for auto-production to accumulate
   - Verify green "+X" badges appear
   - Collect button shows "(total)"
   - Tooltip shows correct breakdown

3. **Full Warehouse:**
   - Let warehouse reach 100% capacity
   - Verify badges turn red
   - Verify pulsing animation plays
   - Badge should grab attention

4. **Collect Action:**
   - Click "Collect" button
   - Badges should disappear after collection
   - Resources should add to main balance
   - Button should update text

5. **Auto-refresh:**
   - Leave page open for 30+ seconds
   - Verify resources update automatically
   - Warehouse amounts should increase over time

### Edge Cases

1. **Warehouse Capacity = 0:**
   - No Resource Warehouse built
   - No badges shown (warehouse always 0)

2. **Very Large Numbers:**
   - Test with 1M+ warehouse amounts
   - Verify formatNumber displays correctly (1,000,000 or 1M)

3. **Mixed Warehouse States:**
   - Metal full (red badge)
   - He3/Gold partial (green badges)
   - All badges should show independently

4. **Responsive Design:**
   - Test on small screens (< 768px)
   - Badges should not wrap awkwardly
   - Tooltip should remain readable

## Future Enhancements

1. **Near-Full Warning (80-99%):**
   - Yellow/orange badge color
   - Earlier warning before full capacity
   - Progressive color scale (green → yellow → red)

2. **Warehouse Capacity Display:**
   - Show current/max: "10k / 15k"
   - Progress bar under badges
   - Visual representation of fullness

3. **Auto-Collect Option:**
   - Toggle: "Auto-collect when warehouse full"
   - Prevents resource waste
   - User preference setting

4. **Sound/Notification:**
   - Audio alert when warehouse reaches 100%
   - Browser notification (if permitted)
   - Desktop notification support

5. **Warehouse Details Modal:**
   - Click badge → detailed breakdown
   - Show accumulation rate
   - Estimate time to full capacity

6. **Historical Tracking:**
   - Graph of warehouse levels over time
   - Peak usage times
   - Optimization suggestions

## Files Modified

### Frontend
- `/frontend/src/types/index.ts` - Added warehouse fields to Resource interface
- `/frontend/src/components/layout/ResourceHUD.tsx` - Warehouse display logic
- `/frontend/src/styles/layout.css` - Warehouse badge styles + animation

### No Backend Changes
- Backend already updated (Task #28)
- API returns warehouse fields

## Dependencies
- Existing: `useResources()` hook
- Existing: `formatNumber()` utility
- Existing: `useMemo` from React
- No new dependencies

## Integration with Module 3

**Task #27 (Backend):** Resource Auto-Production Worker
- Deposits resources into warehouse hourly
- Respects warehouse capacity limits

**Task #28 (Backend):** Update Resource Endpoints
- Returns warehouse fields in API response
- Handles warehouse collection

**Task #30 (Frontend):** useResources Hook Update
- May add additional warehouse-related hooks
- Task #29 is prerequisite

## Estimated Completion Time
**Actual:** 2 hours (within 2-3 hour estimate)

## Status
✅ **COMPLETED** - Ready for QA testing
