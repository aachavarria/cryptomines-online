package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
	"github.com/google/uuid"
)

func setupTierTest(t *testing.T) (playerID string, cleanup func()) {
	if err := database.InitSupabase(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	playerID = uuid.New().String()
	anonID := "anon-" + uuid.New().String()

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username)
		VALUES ($1, $2, 'tier-test-user')
	`, playerID, anonID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}

	cleanup = func() {
		database.DB.Exec("DELETE FROM ship_designs WHERE player_id = $1", playerID)
		database.DB.Exec("DELETE FROM player_blueprints WHERE player_id = $1", playerID)
		database.DB.Exec("DELETE FROM players WHERE id = $1", playerID)
	}

	return playerID, cleanup
}

func TestCreateShipDesign_HullTierNotUnlocked(t *testing.T) {
	playerID, cleanup := setupTierTest(t)
	defer cleanup()

	// Find a tier 2 hull (weikes_ii)
	var tier2HullID int
	var blueprintID int
	err := database.DB.QueryRow(`
		SELECT ht.id, b.id
		FROM hull_types ht
		JOIN blueprints b ON b.hull_type_id = (
			SELECT id FROM hull_types WHERE name = 'weikes_i'
		)
		WHERE ht.name = 'weikes_ii' AND b.blueprint_type = 'hull'
	`).Scan(&tier2HullID, &blueprintID)
	if err != nil {
		t.Fatalf("Failed to find tier 2 hull: %v", err)
	}

	// Give player the blueprint but only at research level 1 (tier 1 only)
	_, err = database.DB.Exec(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 1)
	`, playerID, blueprintID)
	if err != nil {
		t.Fatalf("Failed to create player blueprint: %v", err)
	}

	// Get a module
	var moduleID int
	err = database.DB.QueryRow(`
		SELECT id FROM module_types WHERE tier = 1 LIMIT 1
	`).Scan(&moduleID)
	if err != nil {
		t.Fatalf("Failed to find module: %v", err)
	}

	// Try to create design with tier 2 hull (should fail)
	reqBody := createDesignRequest{
		Name:       "test-design",
		HullTypeID: tier2HullID,
		Modules: []models.DesignModule{
			{ModuleTypeID: moduleID, Quantity: 1, PlacementOrder: 0},
		},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/ship-designs", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	CreateShipDesign(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	if response["error"] != "hull tier not unlocked - research blueprint to unlock" {
		t.Errorf("Expected tier error, got: %s", response["error"])
	}
}

func TestCreateShipDesign_HullTierUnlocked(t *testing.T) {
	playerID, cleanup := setupTierTest(t)
	defer cleanup()

	// Find a tier 2 hull (weikes_ii)
	var tier2HullID int
	var blueprintID int
	err := database.DB.QueryRow(`
		SELECT ht.id, b.id
		FROM hull_types ht
		JOIN blueprints b ON b.hull_type_id = (
			SELECT id FROM hull_types WHERE name = 'weikes_i'
		)
		WHERE ht.name = 'weikes_ii' AND b.blueprint_type = 'hull'
	`).Scan(&tier2HullID, &blueprintID)
	if err != nil {
		t.Fatalf("Failed to find tier 2 hull: %v", err)
	}

	// Give player the blueprint at research level 2 (tier 2 unlocked)
	_, err = database.DB.Exec(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 2)
	`, playerID, blueprintID)
	if err != nil {
		t.Fatalf("Failed to create player blueprint: %v", err)
	}

	// Get a tier 1 module and activate its blueprint
	var moduleID int
	var moduleBlueprintID int
	err = database.DB.QueryRow(`
		SELECT mt.id, b.id
		FROM module_types mt
		JOIN blueprints b ON b.module_type_id = mt.id
		WHERE mt.tier = 1 AND b.blueprint_type = 'module'
		LIMIT 1
	`).Scan(&moduleID, &moduleBlueprintID)
	if err != nil {
		t.Fatalf("Failed to find module: %v", err)
	}

	_, err = database.DB.Exec(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 1)
	`, playerID, moduleBlueprintID)
	if err != nil {
		t.Fatalf("Failed to create module blueprint: %v", err)
	}

	// Try to create design with tier 2 hull (should succeed)
	reqBody := createDesignRequest{
		Name:       "test-design-tier2",
		HullTypeID: tier2HullID,
		Modules: []models.DesignModule{
			{ModuleTypeID: moduleID, Quantity: 1, PlacementOrder: 0},
		},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/ship-designs", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	CreateShipDesign(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateShipDesign_ModuleTierNotUnlocked(t *testing.T) {
	playerID, cleanup := setupTierTest(t)
	defer cleanup()

	// Find a tier 1 hull
	var tier1HullID int
	var hullBlueprintID int
	err := database.DB.QueryRow(`
		SELECT ht.id, b.id
		FROM hull_types ht
		JOIN blueprints b ON b.hull_type_id = ht.id
		WHERE ht.name = 'weikes_i' AND b.blueprint_type = 'hull'
	`).Scan(&tier1HullID, &hullBlueprintID)
	if err != nil {
		t.Fatalf("Failed to find tier 1 hull: %v", err)
	}

	// Activate hull blueprint at level 1
	_, err = database.DB.Exec(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 1)
	`, playerID, hullBlueprintID)
	if err != nil {
		t.Fatalf("Failed to create hull blueprint: %v", err)
	}

	// Find a tier 2 module by name suffix (_ii)
	var tier2ModuleID int
	var tier2ModuleName string
	err = database.DB.QueryRow(`
		SELECT id, name
		FROM module_types
		WHERE name LIKE '%_ii' AND tier = 2
		LIMIT 1
	`).Scan(&tier2ModuleID, &tier2ModuleName)
	if err != nil {
		t.Fatalf("Failed to find tier 2 module: %v", err)
	}

	// Find the tier 1 version of this module's blueprint
	baseName := tier2ModuleName[:len(tier2ModuleName)-3] // Remove "_ii"
	tier1Name := baseName + "_i"

	var moduleBlueprintID int
	err = database.DB.QueryRow(`
		SELECT b.id
		FROM blueprints b
		JOIN module_types mt ON b.module_type_id = mt.id
		WHERE b.blueprint_type = 'module' AND mt.name = $1
	`, tier1Name).Scan(&moduleBlueprintID)
	if err != nil {
		t.Fatalf("Failed to find module blueprint for %s: %v", tier1Name, err)
	}

	// Give player the module blueprint but only at research level 1
	_, err = database.DB.Exec(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 1)
	`, playerID, moduleBlueprintID)
	if err != nil {
		t.Fatalf("Failed to create module blueprint: %v", err)
	}

	// Try to create design with tier 2 module (should fail)
	reqBody := createDesignRequest{
		Name:       "test-module-tier",
		HullTypeID: tier1HullID,
		Modules: []models.DesignModule{
			{ModuleTypeID: tier2ModuleID, Quantity: 1, PlacementOrder: 0},
		},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/ship-designs", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	CreateShipDesign(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	if response["error"] != "module tier not unlocked - research blueprint to unlock" {
		t.Errorf("Expected module tier error, got: %s", response["error"])
	}
}

func TestUpdateShipDesign_HullTierNotUnlocked(t *testing.T) {
	playerID, cleanup := setupTierTest(t)
	defer cleanup()

	// Create a design with tier 1 hull first
	var tier1HullID int
	var hullBlueprintID int
	err := database.DB.QueryRow(`
		SELECT ht.id, b.id
		FROM hull_types ht
		JOIN blueprints b ON b.hull_type_id = ht.id
		WHERE ht.name = 'weikes_i' AND b.blueprint_type = 'hull'
	`).Scan(&tier1HullID, &hullBlueprintID)
	if err != nil {
		t.Fatalf("Failed to find tier 1 hull: %v", err)
	}

	// Activate hull blueprint at level 1
	_, err = database.DB.Exec(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 1)
	`, playerID, hullBlueprintID)
	if err != nil {
		t.Fatalf("Failed to create hull blueprint: %v", err)
	}

	// Get a module and activate its blueprint
	var moduleID int
	var moduleBlueprintID int
	err = database.DB.QueryRow(`
		SELECT mt.id, b.id
		FROM module_types mt
		JOIN blueprints b ON b.module_type_id = mt.id
		WHERE mt.tier = 1 AND b.blueprint_type = 'module'
		LIMIT 1
	`).Scan(&moduleID, &moduleBlueprintID)
	if err != nil {
		t.Fatalf("Failed to find module: %v", err)
	}

	_, err = database.DB.Exec(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 1)
	`, playerID, moduleBlueprintID)
	if err != nil {
		t.Fatalf("Failed to create module blueprint: %v", err)
	}

	// Create the design
	var designID string
	err = database.DB.QueryRow(`
		INSERT INTO ship_designs (player_id, name, hull_type_id, modules_json, total_shield, total_structure, total_defense, total_agility, total_movement, total_storage, attack_power, weapon_range_min, weapon_range_max, volume_used, he3_per_round, metal_cost, he3_cost, gold_cost, build_time_seconds)
		VALUES ($1, 'original-design', $2, '[]', 100, 100, 10, 10, 5, 100, 50, 1, 5, 10, 5, 1000, 500, 50, 3600)
		RETURNING id
	`, playerID, tier1HullID).Scan(&designID)
	if err != nil {
		t.Fatalf("Failed to create design: %v", err)
	}

	// Find tier 2 hull
	var tier2HullID int
	err = database.DB.QueryRow(`
		SELECT id FROM hull_types WHERE name = 'weikes_ii'
	`).Scan(&tier2HullID)
	if err != nil {
		t.Fatalf("Failed to find tier 2 hull: %v", err)
	}

	// Try to update design to use tier 2 hull (should fail - only level 1 unlocked)
	reqBody := createDesignRequest{
		Name:       "updated-design",
		HullTypeID: tier2HullID,
		Modules: []models.DesignModule{
			{ModuleTypeID: moduleID, Quantity: 1, PlacementOrder: 0},
		},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/ship-designs/"+designID, bytes.NewReader(body))
	req.SetPathValue("id", designID)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	UpdateShipDesign(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	if response["error"] != "hull tier not unlocked - research blueprint to unlock" {
		t.Errorf("Expected tier error, got: %s", response["error"])
	}
}
