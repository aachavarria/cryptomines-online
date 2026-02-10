# Enhanced Error Messages Implementation

## Overview
Replaced generic error responses with detailed, user-friendly messages that include context about requirements, current state, and helpful hints.

## Implementation Date
2026-02-07

## Files Modified

### 1. New Error Response Package
**File**: `/backend/internal/errs/responses.go`

Created standardized error response structure with helper functions:

```go
type ErrorResponse struct {
    Error   string                 `json:"error"`
    Details string                 `json:"details,omitempty"`
    Context map[string]interface{} `json:"context,omitempty"`
}
```

**Helper Functions**:
- `InsufficientResources()` - Resource shortage with required vs current
- `PrerequisiteNotMet()` - Prerequisite details with current vs needed levels
- `MaxLevelReached()` - Max level information
- `MaxCountReached()` - Max count with current count
- `OnCooldown()` - Cooldown info with remaining time
- `LevelTooLow()` - Level requirement with current level
- `NotFound()` - Not found errors
- `Conflict()` - Conflict errors with context
- `InternalError()` - Internal server errors

### 2. Buildings Handler
**File**: `/backend/internal/handlers/buildings.go`

**Enhanced Errors**:
1. **Construction Slots Full**:
   ```json
   {
     "error": "conflict",
     "details": "All 2 construction slots are in use",
     "context": {
       "slots_used": 2,
       "slots_available": 2,
       "hint": "Upgrade Concurrent Construction tech or wait for current construction to finish"
     }
   }
   ```

2. **Max Count Reached**:
   ```json
   {
     "error": "max_count_reached",
     "details": "Maximum Metal Collector count reached (10/10)",
     "context": {
       "max_count": 10,
       "current_count": 10
     }
   }
   ```

3. **Prerequisite Not Met**:
   ```json
   {
     "error": "prerequisite_not_met",
     "details": "Ship Factory requires Civic Center level 5 (current: 3)",
     "context": {
       "required_building": "Civic Center",
       "current_level": 3,
       "needed_level": 5
     }
   }
   ```

4. **Insufficient Resources**:
   ```json
   {
     "error": "insufficient_resources",
     "details": "Not enough resources to construct Metal Collector",
     "context": {
       "required": {"metal": 1000, "he3": 500, "gold": 200},
       "current": {"metal": 800, "he3": 300, "gold": 150}
     }
   }
   ```

5. **Max Level Reached**:
   ```json
   {
     "error": "max_level_reached",
     "details": "Metal Collector is already at max level 12",
     "context": {
       "max_level": 12
     }
   }
   ```

6. **Civic Center Requirement**:
   ```json
   {
     "error": "prerequisite_not_met",
     "details": "Civic Center must be level 5 to upgrade Metal Collector to level 5 (current: 3)",
     "context": {
       "required_building": "Civic Center",
       "current_level": 3,
       "needed_level": 5
     }
   }
   ```

**Helper Function Added**:
- `checkAndDeductResources()` - Checks resource sufficiency before deducting, returns detailed error

### 3. Research Handler
**File**: `/backend/internal/handlers/research.go`

**Enhanced Errors**:
1. **Already Researching**:
   ```json
   {
     "error": "conflict",
     "details": "Already researching a technology in the logistics_construction tree",
     "context": {
       "tree": "logistics_construction",
       "hint": "Wait for current research to complete or cancel it first"
     }
   }
   ```

2. **Max Level Reached**:
   ```json
   {
     "error": "max_level_reached",
     "details": "Building Construction is already at max level 10",
     "context": {
       "max_level": 10
     }
   }
   ```

3. **Prerequisite Not Met**:
   ```json
   {
     "error": "prerequisite_not_met",
     "details": "Advanced Construction requires Building Construction level 5 (current: 2)",
     "context": {
       "required_building": "Building Construction",
       "current_level": 2,
       "needed_level": 5
     }
   }
   ```

4. **Insufficient Resources**:
   ```json
   {
     "error": "insufficient_resources",
     "details": "Not enough resources to research Building Construction level 3",
     "context": {
       "required": {"metal": 5000, "he3": 3000, "gold": 10000},
       "current": {"metal": 3000, "he3": 2000, "gold": 8000}
     }
   }
   ```

### 4. Ship Factory Handler
**File**: `/backend/internal/handlers/ship_factory.go`

**Enhanced Errors**:
1. **Insufficient Resources**:
   ```json
   {
     "error": "insufficient_resources",
     "details": "Not enough resources to build 100 x Viper Frigate",
     "context": {
       "required": {"metal": 100000, "he3": 50000, "gold": 25000},
       "current": {"metal": 80000, "he3": 40000, "gold": 20000}
     }
   }
   ```

### 5. Commanders Handler
**File**: `/backend/internal/handlers/commanders.go`

**Enhanced Errors**:
1. **Insufficient Gold**:
   ```json
   {
     "error": "insufficient_resources",
     "details": "Not enough Gold to recruit commander (need 10000, have 8000)",
     "context": {
       "required": {"gold": 10000},
       "current": {"gold": 8000}
     }
   }
   ```

### 6. Fleets Handler
**File**: `/backend/internal/handlers/fleets.go`

**Enhanced Errors**:
1. **Planet Not Found**:
   ```json
   {
     "error": "not_found",
     "details": "Planet abc-123 not found or does not belong to you"
   }
   ```

### 7. PvP Handler
**File**: `/backend/internal/handlers/pvp.go`

**Enhanced Errors**:
1. **Target Not Found**:
   ```json
   {
     "error": "not_found",
     "details": "Target planet abc-123 not found"
   }
   ```

## Error Response Format

All enhanced errors follow this structure:

```typescript
interface ErrorResponse {
  error: string;              // Error type (snake_case)
  details?: string;           // Human-readable description
  context?: {                 // Additional context
    required?: object;        // Required values
    current?: object;         // Current values
    hint?: string;            // Helpful hint for user
    [key: string]: any;       // Other contextual data
  };
}
```

## Error Types

1. `insufficient_resources` - Not enough resources
2. `prerequisite_not_met` - Prerequisite requirement not satisfied
3. `max_level_reached` - Entity at maximum level
4. `max_count_reached` - Maximum count reached
5. `on_cooldown` - Action on cooldown
6. `level_too_low` - Player/building level too low
7. `not_found` - Resource not found
8. `already_exists` - Resource already exists
9. `invalid_request` - Invalid request data
10. `internal_error` - Internal server error
11. `conflict` - General conflict error

## Benefits

1. **User-Friendly**: Clear, detailed messages explain what went wrong
2. **Actionable**: Users know exactly what's needed to proceed
3. **Context-Rich**: Shows required vs current values for easy comparison
4. **Consistent**: Standardized format across all endpoints
5. **Debuggable**: Detailed context helps with troubleshooting

## Testing

Compile test successful:
```bash
cd /Users/yurei/cryptomines-online/backend && go build ./...
```

## Next Steps

Frontend should be updated to:
1. Parse and display enhanced error messages
2. Show resource requirements vs current
3. Display prerequisite information
4. Provide visual feedback based on error context
5. Show helpful hints to users

## Example Frontend Integration

```typescript
// Frontend error handling
if (response.error === 'insufficient_resources') {
  const { required, current } = response.context;
  showResourceComparison(required, current);
  highlightShortages(required, current);
}

if (response.error === 'prerequisite_not_met') {
  const { required_building, current_level, needed_level } = response.context;
  showPrerequisiteTooltip(required_building, current_level, needed_level);
}
```

## Coverage

- ✅ Buildings: Construction, Upgrade errors
- ✅ Research: Tech research errors
- ✅ Ship Factory: Ship building errors
- ✅ Commanders: Recruitment errors
- ✅ Fleets: Fleet management errors
- ✅ PvP: Combat errors
- 🔄 Additional handlers can be enhanced as needed

## Notes

- All errors now include detailed context
- Resource errors show exact shortages
- Prerequisite errors show current vs required levels
- Helper functions ensure consistency
- Backward compatible with existing error handling
