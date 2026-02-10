package services

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/cryptomines-online/backend/internal/database"
)

// GetHullBaseName extracts the base hull name without tier suffix.
// Examples:
//   - "weikes_i" → "weikes"
//   - "weikes_ii" → "weikes"
//   - "weikes_iii" → "weikes"
//   - "darak_i" → "darak"
//   - "scarab" → "scarab" (no tier suffix)
func GetHullBaseName(hullName string) string {
	// Remove tier suffixes: _i, _ii, _iii
	name := strings.TrimSuffix(hullName, "_iii")
	name = strings.TrimSuffix(name, "_ii")
	name = strings.TrimSuffix(name, "_i")
	return name
}

// GetHullTier extracts the tier number from a hull name.
// Returns 1 for _i, 2 for _ii, 3 for _iii, or 0 if no tier suffix.
func GetHullTier(hullName string) int {
	if strings.HasSuffix(hullName, "_iii") {
		return 3
	}
	if strings.HasSuffix(hullName, "_ii") {
		return 2
	}
	if strings.HasSuffix(hullName, "_i") {
		return 1
	}
	return 0 // No tier suffix
}

// CanUseHullTier checks if a player can use a specific hull based on blueprint research.
// Players must have the blueprint activated and researched to the appropriate tier level.
//
// Tier requirements:
//   - Tier 1 (_i): Blueprint activated (research_level >= 1)
//   - Tier 2 (_ii): Blueprint researched to level 2
//   - Tier 3 (_iii): Blueprint researched to level 3
func CanUseHullTier(playerID string, hullTypeID int) (bool, error) {
	// Get hull info
	var hullName string
	var tier int
	err := database.DB.QueryRow(`
		SELECT name, tier
		FROM hull_types
		WHERE id = $1
	`, hullTypeID).Scan(&hullName, &tier)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("hull type not found")
		}
		log.Printf("Failed to get hull type: %v", err)
		return false, fmt.Errorf("failed to get hull type: %w", err)
	}

	// Get base hull name (without tier suffix)
	baseName := GetHullBaseName(hullName)

	// Find the blueprint for this hull's base name
	// The blueprint references the tier 1 hull (e.g., weikes_i)
	var blueprintID int
	err = database.DB.QueryRow(`
		SELECT b.id
		FROM blueprints b
		JOIN hull_types ht ON b.hull_type_id = ht.id
		WHERE b.blueprint_type = 'hull'
		  AND ht.name = $1
	`, baseName+"_i").Scan(&blueprintID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No blueprint exists for this hull - allow usage (base game hulls)
			return true, nil
		}
		log.Printf("Failed to find blueprint: %v", err)
		return false, fmt.Errorf("failed to find blueprint: %w", err)
	}

	// Check if player has this blueprint activated and researched to required level
	var researchLevel int
	var isActivated bool
	err = database.DB.QueryRow(`
		SELECT research_level, is_activated
		FROM player_blueprints
		WHERE player_id = $1 AND blueprint_id = $2
	`, playerID, blueprintID).Scan(&researchLevel, &isActivated)
	if err != nil {
		if err == sql.ErrNoRows {
			// Player doesn't have this blueprint
			return false, nil
		}
		log.Printf("Failed to check player blueprint: %v", err)
		return false, fmt.Errorf("failed to check player blueprint: %w", err)
	}

	// Blueprint must be activated
	if !isActivated {
		return false, nil
	}

	// Check if research level meets tier requirement
	// Tier 1 requires level >= 1 (just activated)
	// Tier 2 requires level >= 2
	// Tier 3 requires level >= 3
	return researchLevel >= tier, nil
}

// CanUseModuleTier checks if a player can use a specific module based on blueprint research.
// Similar to CanUseHullTier but for modules.
//
// Module blueprints work the same way:
//   - Tier 1: Blueprint activated (research_level >= 1)
//   - Tier 2: Blueprint researched to level 2
//   - Tier 3: Blueprint researched to level 3
func CanUseModuleTier(playerID string, moduleTypeID int) (bool, error) {
	// Get module info
	var moduleName string
	var tier int
	err := database.DB.QueryRow(`
		SELECT name, tier
		FROM module_types
		WHERE id = $1
	`, moduleTypeID).Scan(&moduleName, &tier)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("module type not found")
		}
		log.Printf("Failed to get module type: %v", err)
		return false, fmt.Errorf("failed to get module type: %w", err)
	}

	// Get base module name (modules may have tier suffixes too)
	baseName := GetHullBaseName(moduleName) // Reuse same logic

	// Find the blueprint for this module's base name
	// The blueprint references the tier 1 module
	var blueprintID int
	err = database.DB.QueryRow(`
		SELECT b.id
		FROM blueprints b
		JOIN module_types mt ON b.module_type_id = mt.id
		WHERE b.blueprint_type = 'module'
		  AND mt.name = $1
	`, baseName+"_i").Scan(&blueprintID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No blueprint exists for this module - allow usage (base game modules)
			return true, nil
		}
		log.Printf("Failed to find module blueprint: %v", err)
		return false, fmt.Errorf("failed to find module blueprint: %w", err)
	}

	// Check if player has this blueprint activated and researched to required level
	var researchLevel int
	var isActivated bool
	err = database.DB.QueryRow(`
		SELECT research_level, is_activated
		FROM player_blueprints
		WHERE player_id = $1 AND blueprint_id = $2
	`, playerID, blueprintID).Scan(&researchLevel, &isActivated)
	if err != nil {
		if err == sql.ErrNoRows {
			// Player doesn't have this blueprint
			return false, nil
		}
		log.Printf("Failed to check player module blueprint: %v", err)
		return false, fmt.Errorf("failed to check player module blueprint: %w", err)
	}

	// Blueprint must be activated
	if !isActivated {
		return false, nil
	}

	// Check if research level meets tier requirement
	return researchLevel >= tier, nil
}

// GetPlayerBlueprintLevel returns the research level of a blueprint for a player.
// Returns 0 if player doesn't have the blueprint or it's not activated.
// This is useful for determining what tier of hulls/modules a player can use.
func GetPlayerBlueprintLevel(playerID string, blueprintID int) (int, error) {
	var researchLevel int
	var isActivated bool
	err := database.DB.QueryRow(`
		SELECT research_level, is_activated
		FROM player_blueprints
		WHERE player_id = $1 AND blueprint_id = $2
	`, playerID, blueprintID).Scan(&researchLevel, &isActivated)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // Player doesn't have this blueprint
		}
		log.Printf("Failed to get player blueprint level: %v", err)
		return 0, fmt.Errorf("failed to get player blueprint level: %w", err)
	}

	if !isActivated {
		return 0, nil // Blueprint not activated
	}

	return researchLevel, nil
}
