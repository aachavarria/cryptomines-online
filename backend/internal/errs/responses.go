package errs

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a detailed error response with context
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Details string                 `json:"details,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
}

// ResourceContext provides resource requirement details
type ResourceContext struct {
	Required map[string]int64 `json:"required"`
	Current  map[string]int64 `json:"current"`
}

// PrerequisiteContext provides prerequisite details
type PrerequisiteContext struct {
	Required string `json:"required"`
	Current  int    `json:"current_level"`
	Needed   int    `json:"needed_level"`
}

// CooldownContext provides cooldown information
type CooldownContext struct {
	Action       string `json:"action"`
	AvailableAt  string `json:"available_at"`
	RemainingSeconds int `json:"remaining_seconds"`
}

// LevelContext provides level requirement details
type LevelContext struct {
	Required int `json:"required"`
	Current  int `json:"current"`
}

// WriteJSON writes an ErrorResponse as JSON with the given HTTP status
func (e *ErrorResponse) WriteJSON(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(e)
}

// InsufficientResources creates a resource shortage error
func InsufficientResources(details string, required, current map[string]int64) *ErrorResponse {
	return &ErrorResponse{
		Error:   "insufficient_resources",
		Details: details,
		Context: map[string]interface{}{
			"required": required,
			"current":  current,
		},
	}
}

// PrerequisiteNotMet creates a prerequisite error
func PrerequisiteNotMet(details, requiredBuilding string, currentLevel, neededLevel int) *ErrorResponse {
	return &ErrorResponse{
		Error:   "prerequisite_not_met",
		Details: details,
		Context: map[string]interface{}{
			"required_building": requiredBuilding,
			"current_level":     currentLevel,
			"needed_level":      neededLevel,
		},
	}
}

// MaxLevelReached creates a max level error
func MaxLevelReached(details string, maxLevel int) *ErrorResponse {
	return &ErrorResponse{
		Error:   "max_level_reached",
		Details: details,
		Context: map[string]interface{}{
			"max_level": maxLevel,
		},
	}
}

// MaxCountReached creates a max count error
func MaxCountReached(details string, maxCount, currentCount int) *ErrorResponse {
	return &ErrorResponse{
		Error:   "max_count_reached",
		Details: details,
		Context: map[string]interface{}{
			"max_count":     maxCount,
			"current_count": currentCount,
		},
	}
}

// OnCooldown creates a cooldown error
func OnCooldown(details, action, availableAt string, remainingSeconds int) *ErrorResponse {
	return &ErrorResponse{
		Error:   "on_cooldown",
		Details: details,
		Context: map[string]interface{}{
			"action":            action,
			"available_at":      availableAt,
			"remaining_seconds": remainingSeconds,
		},
	}
}

// LevelTooLow creates a level requirement error
func LevelTooLow(details string, required, current int) *ErrorResponse {
	return &ErrorResponse{
		Error:   "level_too_low",
		Details: details,
		Context: map[string]interface{}{
			"required": required,
			"current":  current,
		},
	}
}

// NotFound creates a not found error
func NotFound(details string) *ErrorResponse {
	return &ErrorResponse{
		Error:   "not_found",
		Details: details,
	}
}

// AlreadyExists creates an already exists error
func AlreadyExists(details string) *ErrorResponse {
	return &ErrorResponse{
		Error:   "already_exists",
		Details: details,
	}
}

// InvalidRequest creates an invalid request error
func InvalidRequest(details string) *ErrorResponse {
	return &ErrorResponse{
		Error:   "invalid_request",
		Details: details,
	}
}

// InternalError creates an internal server error
func InternalError(details string) *ErrorResponse {
	return &ErrorResponse{
		Error:   "internal_error",
		Details: details,
	}
}

// Conflict creates a conflict error
func Conflict(details string, context map[string]interface{}) *ErrorResponse {
	return &ErrorResponse{
		Error:   "conflict",
		Details: details,
		Context: context,
	}
}
