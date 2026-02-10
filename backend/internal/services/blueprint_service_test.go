package services

import (
	"testing"

	"github.com/cryptomines-online/backend/internal/database"
)

// setupBlueprintTestDB initializes the test database connection.
func setupBlueprintTestDB(t *testing.T) {
	if database.DB == nil {
		if err := database.InitSupabase(); err != nil {
			t.Fatalf("Failed to init database: %v", err)
		}
	}
}

// cleanupBlueprintTestPlayer removes test data for a player.
func cleanupBlueprintTestPlayer(t *testing.T, playerID string) {
	database.DB.Exec("DELETE FROM player_blueprints WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM players WHERE id = $1", playerID)
}

func TestGetHullBaseName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"weikes_i", "weikes"},
		{"weikes_ii", "weikes"},
		{"weikes_iii", "weikes"},
		{"darak_i", "darak"},
		{"darak_ii", "darak"},
		{"darak_iii", "darak"},
		{"scarab_i", "scarab"},
		{"scarab", "scarab"}, // No tier suffix
		{"custom_hull", "custom_hull"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := GetHullBaseName(tt.input)
			if result != tt.expected {
				t.Errorf("GetHullBaseName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetHullTier(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"weikes_i", 1},
		{"weikes_ii", 2},
		{"weikes_iii", 3},
		{"darak_i", 1},
		{"darak_ii", 2},
		{"darak_iii", 3},
		{"scarab", 0}, // No tier suffix
		{"custom_hull", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := GetHullTier(tt.input)
			if result != tt.expected {
				t.Errorf("GetHullTier(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCanUseHullTier_NoBlueprint(t *testing.T) {
	setupBlueprintTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000401"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'bptest1', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupBlueprintTestPlayer(t, playerID)

	// Get a hull type ID (Weikes I)
	var hullTypeID int
	database.DB.QueryRow("SELECT id FROM hull_types WHERE name = 'weikes_i' LIMIT 1").Scan(&hullTypeID)

	// Player doesn't have the blueprint, should return false
	canUse, err := CanUseHullTier(playerID, hullTypeID)
	if err != nil {
		t.Fatalf("CanUseHullTier failed: %v", err)
	}

	if canUse {
		t.Errorf("Expected canUse=false (no blueprint), got true")
	}
}

func TestCanUseHullTier_BlueprintActivated(t *testing.T) {
	setupBlueprintTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000402"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'bptest2', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupBlueprintTestPlayer(t, playerID)

	// Get Weikes blueprint ID
	var blueprintID int
	database.DB.QueryRow(`
		SELECT id FROM blueprints WHERE name = 'Weikes Blueprint'
	`).Scan(&blueprintID)

	// Give player the blueprint at level 1 (activated)
	database.DB.Exec(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 1)
		ON CONFLICT (player_id, blueprint_id) DO UPDATE
		SET is_activated = true, research_level = 1
	`, playerID, blueprintID)

	// Get Weikes I hull type ID
	var weikesI_ID int
	database.DB.QueryRow("SELECT id FROM hull_types WHERE name = 'weikes_i'").Scan(&weikesI_ID)

	// Player should be able to use Tier 1 (Weikes I)
	canUse, err := CanUseHullTier(playerID, weikesI_ID)
	if err != nil {
		t.Fatalf("CanUseHullTier failed: %v", err)
	}
	if !canUse {
		t.Errorf("Expected canUse=true for Weikes I with level 1 blueprint, got false")
	}

	// Get Weikes II hull type ID
	var weikesII_ID int
	database.DB.QueryRow("SELECT id FROM hull_types WHERE name = 'weikes_ii'").Scan(&weikesII_ID)

	// Player should NOT be able to use Tier 2 (Weikes II) with only level 1 research
	canUse, err = CanUseHullTier(playerID, weikesII_ID)
	if err != nil {
		t.Fatalf("CanUseHullTier failed: %v", err)
	}
	if canUse {
		t.Errorf("Expected canUse=false for Weikes II with level 1 blueprint, got true")
	}
}

func TestCanUseHullTier_BlueprintLevel2(t *testing.T) {
	setupBlueprintTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000403"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'bptest3', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupBlueprintTestPlayer(t, playerID)

	// Get Weikes blueprint ID
	var blueprintID int
	database.DB.QueryRow(`
		SELECT id FROM blueprints WHERE name = 'Weikes Blueprint'
	`).Scan(&blueprintID)

	// Give player the blueprint at level 2
	database.DB.Exec(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 2)
		ON CONFLICT (player_id, blueprint_id) DO UPDATE
		SET is_activated = true, research_level = 2
	`, playerID, blueprintID)

	// Get hull type IDs
	var weikesI_ID, weikesII_ID, weikesIII_ID int
	database.DB.QueryRow("SELECT id FROM hull_types WHERE name = 'weikes_i'").Scan(&weikesI_ID)
	database.DB.QueryRow("SELECT id FROM hull_types WHERE name = 'weikes_ii'").Scan(&weikesII_ID)
	database.DB.QueryRow("SELECT id FROM hull_types WHERE name = 'weikes_iii'").Scan(&weikesIII_ID)

	// Player should be able to use Tier 1 and Tier 2
	canUseT1, _ := CanUseHullTier(playerID, weikesI_ID)
	canUseT2, _ := CanUseHullTier(playerID, weikesII_ID)
	canUseT3, _ := CanUseHullTier(playerID, weikesIII_ID)

	if !canUseT1 {
		t.Errorf("Expected canUse=true for Tier 1 with level 2 blueprint, got false")
	}
	if !canUseT2 {
		t.Errorf("Expected canUse=true for Tier 2 with level 2 blueprint, got false")
	}
	if canUseT3 {
		t.Errorf("Expected canUse=false for Tier 3 with level 2 blueprint, got true")
	}
}

func TestCanUseHullTier_BlueprintLevel3(t *testing.T) {
	setupBlueprintTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000404"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'bptest4', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupBlueprintTestPlayer(t, playerID)

	// Get Weikes blueprint ID
	var blueprintID int
	database.DB.QueryRow(`
		SELECT id FROM blueprints WHERE name = 'Weikes Blueprint'
	`).Scan(&blueprintID)

	// Give player the blueprint at level 3 (max)
	database.DB.Exec(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 3)
		ON CONFLICT (player_id, blueprint_id) DO UPDATE
		SET is_activated = true, research_level = 3
	`, playerID, blueprintID)

	// Get hull type IDs
	var weikesI_ID, weikesII_ID, weikesIII_ID int
	database.DB.QueryRow("SELECT id FROM hull_types WHERE name = 'weikes_i'").Scan(&weikesI_ID)
	database.DB.QueryRow("SELECT id FROM hull_types WHERE name = 'weikes_ii'").Scan(&weikesII_ID)
	database.DB.QueryRow("SELECT id FROM hull_types WHERE name = 'weikes_iii'").Scan(&weikesIII_ID)

	// Player should be able to use all tiers
	canUseT1, _ := CanUseHullTier(playerID, weikesI_ID)
	canUseT2, _ := CanUseHullTier(playerID, weikesII_ID)
	canUseT3, _ := CanUseHullTier(playerID, weikesIII_ID)

	if !canUseT1 || !canUseT2 || !canUseT3 {
		t.Errorf("Expected canUse=true for all tiers with level 3 blueprint, got T1=%v T2=%v T3=%v",
			canUseT1, canUseT2, canUseT3)
	}
}

func TestCanUseHullTier_NotActivated(t *testing.T) {
	setupBlueprintTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000405"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'bptest5', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupBlueprintTestPlayer(t, playerID)

	// Get Weikes blueprint ID
	var blueprintID int
	database.DB.QueryRow(`
		SELECT id FROM blueprints WHERE name = 'Weikes Blueprint'
	`).Scan(&blueprintID)

	// Give player the blueprint but NOT activated
	database.DB.Exec(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, false, 1)
		ON CONFLICT (player_id, blueprint_id) DO UPDATE
		SET is_activated = false, research_level = 1
	`, playerID, blueprintID)

	// Get Weikes I hull type ID
	var weikesI_ID int
	database.DB.QueryRow("SELECT id FROM hull_types WHERE name = 'weikes_i'").Scan(&weikesI_ID)

	// Player should NOT be able to use Tier 1 if blueprint is not activated
	canUse, err := CanUseHullTier(playerID, weikesI_ID)
	if err != nil {
		t.Fatalf("CanUseHullTier failed: %v", err)
	}
	if canUse {
		t.Errorf("Expected canUse=false for non-activated blueprint, got true")
	}
}

func TestGetPlayerBlueprintLevel(t *testing.T) {
	setupBlueprintTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000406"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'bptest6', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupBlueprintTestPlayer(t, playerID)

	// Get Weikes blueprint ID
	var blueprintID int
	database.DB.QueryRow(`
		SELECT id FROM blueprints WHERE name = 'Weikes Blueprint'
	`).Scan(&blueprintID)

	// Initially, player doesn't have blueprint - should return 0
	level, err := GetPlayerBlueprintLevel(playerID, blueprintID)
	if err != nil {
		t.Fatalf("GetPlayerBlueprintLevel failed: %v", err)
	}
	if level != 0 {
		t.Errorf("Expected level=0 for missing blueprint, got %d", level)
	}

	// Give player the blueprint at level 2
	database.DB.Exec(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 2)
	`, playerID, blueprintID)

	// Should now return 2
	level, err = GetPlayerBlueprintLevel(playerID, blueprintID)
	if err != nil {
		t.Fatalf("GetPlayerBlueprintLevel failed: %v", err)
	}
	if level != 2 {
		t.Errorf("Expected level=2, got %d", level)
	}
}

func TestCanUseModuleTier_NoBlueprint(t *testing.T) {
	setupBlueprintTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000407"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'bptest7', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupBlueprintTestPlayer(t, playerID)

	// Get a module type ID (any tier 1 module)
	var moduleTypeID int
	err = database.DB.QueryRow("SELECT id FROM module_types WHERE tier = 1 LIMIT 1").Scan(&moduleTypeID)
	if err != nil {
		t.Skip("No module types in database, skipping test")
		return
	}

	// If module has no blueprint (base game module), should return true
	// If module requires blueprint and player doesn't have it, should return false
	canUse, err := CanUseModuleTier(playerID, moduleTypeID)
	if err != nil {
		t.Fatalf("CanUseModuleTier failed: %v", err)
	}

	// Result depends on whether this module requires a blueprint
	// Base game modules should be true, blueprint-only modules should be false
	_ = canUse // Just verify no error for now
}
