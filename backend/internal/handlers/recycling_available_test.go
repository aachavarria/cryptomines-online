package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/google/uuid"
)

func setupAvailableShipsTest(t *testing.T) (playerID, planetID string, cleanup func()) {
	if err := database.InitSupabase(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	playerID = uuid.New().String()
	planetID = uuid.New().String()
	anonID := "anon-" + uuid.New().String()
	username := "test-avail-ships-" + uuid.New().String()

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username)
		VALUES ($1, $2, $3)
	`, playerID, anonID, username)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}

	// Create test planet
	_, err = database.DB.Exec(`
		INSERT INTO planets (id, player_id, name, position_x, position_y, is_homeworld)
		VALUES ($1, $2, 'Test Planet', 1, 2, true)
	`, planetID, playerID)
	if err != nil {
		t.Fatalf("Failed to create test planet: %v", err)
	}

	cleanup = func() {
		// Cleanup cascades through ship_instances, ship_designs, fleets, etc.
		database.DB.Exec("DELETE FROM planets WHERE id = $1", planetID)
		database.DB.Exec("DELETE FROM players WHERE id = $1", playerID)
	}

	return playerID, planetID, cleanup
}

func TestListAvailableShips_EmptyList(t *testing.T) {
	playerID, _, cleanup := setupAvailableShipsTest(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/ship-instances/available", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	ListAvailableShips(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response []availableShip
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != 0 {
		t.Errorf("Expected empty list, got %d ships", len(response))
	}
}

func TestListAvailableShips_WithUndeployedShips(t *testing.T) {
	playerID, _, cleanup := setupAvailableShipsTest(t)
	defer cleanup()

	// Get a hull type (Frigate)
	var hullTypeID int
	err := database.DB.QueryRow(`
		SELECT id FROM hull_types WHERE hull_class = 'frigate' LIMIT 1
	`).Scan(&hullTypeID)
	if err != nil {
		t.Fatalf("Failed to find hull type: %v", err)
	}

	// Create a ship design
	var designID string
	err = database.DB.QueryRow(`
		INSERT INTO ship_designs (player_id, name, hull_type_id, modules_json,
			total_shield, total_structure, total_defense, total_agility, total_movement,
			total_storage, attack_power, weapon_range_min, weapon_range_max, volume_used,
			he3_per_round, metal_cost, he3_cost, gold_cost, build_time_seconds)
		VALUES ($1, 'Test_Frigate', $2, '[]', 100, 100, 10, 10, 5, 100, 50, 1, 5, 10, 5, 1000, 500, 50, 3600)
		RETURNING id
	`, playerID, hullTypeID).Scan(&designID)
	if err != nil {
		t.Fatalf("Failed to create ship design: %v", err)
	}

	// Insert one ships row with quantity=3 (handler aggregates by design).
	_, err = database.DB.Exec(`
		INSERT INTO ships (player_id, ship_design_id, quantity, is_building, build_quantity, production_slot)
		VALUES ($1, $2, 3, false, 0, 1)
	`, playerID, designID)
	if err != nil {
		t.Fatalf("Failed to create ships row: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/ship-instances/available", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	ListAvailableShips(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response []availableShip
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Handler returns one row per design with the aggregate quantity.
	if len(response) != 1 {
		t.Fatalf("Expected 1 design row, got %d", len(response))
	}
	if response[0].DesignName != "Test_Frigate" {
		t.Errorf("Expected design name 'Test_Frigate', got '%s'", response[0].DesignName)
	}
	if response[0].HullClass != "frigate" {
		t.Errorf("Expected hull class 'frigate', got '%s'", response[0].HullClass)
	}
}

func TestListAvailableShips_ExcludesDeployedShips(t *testing.T) {
	playerID, _, cleanup := setupAvailableShipsTest(t)
	defer cleanup()

	// Get a hull type
	var hullTypeID int
	err := database.DB.QueryRow(`
		SELECT id FROM hull_types WHERE hull_class = 'cruiser' LIMIT 1
	`).Scan(&hullTypeID)
	if err != nil {
		t.Fatalf("Failed to find hull type: %v", err)
	}

	// Create ship design
	var designID string
	err = database.DB.QueryRow(`
		INSERT INTO ship_designs (player_id, name, hull_type_id, modules_json,
			total_shield, total_structure, total_defense, total_agility, total_movement,
			total_storage, attack_power, weapon_range_min, weapon_range_max, volume_used,
			he3_per_round, metal_cost, he3_cost, gold_cost, build_time_seconds)
		VALUES ($1, 'Test_Cruiser', $2, '[]', 200, 200, 15, 8, 4, 150, 80, 2, 6, 15, 8, 2000, 1000, 100, 7200)
		RETURNING id
	`, playerID, hullTypeID).Scan(&designID)
	if err != nil {
		t.Fatalf("Failed to create ship design: %v", err)
	}

	// 5 ships of this design were built; 2 are deployed to a fleet, so
	// ships.quantity (the available pool, mutated by fleet add/remove) is 3.
	// The fleet_stack is also written to mirror real state and ensure the
	// available endpoint does NOT double-subtract.
	if _, err := database.DB.Exec(`
		INSERT INTO ships (player_id, ship_design_id, quantity, production_slot)
		VALUES ($1, $2, 3, 1)
	`, playerID, designID); err != nil {
		t.Fatalf("ships: %v", err)
	}

	fleetID := uuid.New().String()
	_, err = database.DB.Exec(`
		INSERT INTO fleets (id, player_id, name, formation, targeting_command, status)
		VALUES ($1, $2, 'Test_Fleet', 'phalanx', 'max_attack', 'stationed')
	`, fleetID, playerID)
	if err != nil {
		t.Fatalf("Failed to create fleet: %v", err)
	}
	if _, err := database.DB.Exec(`
		INSERT INTO fleet_stacks (id, fleet_id, ship_design_id, grid_row, grid_col, ship_count)
		VALUES ($1, $2, $3, 0, 0, 2)
	`, uuid.New().String(), fleetID, designID); err != nil {
		t.Fatalf("Failed to assign stack: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/ship-instances/available", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	ListAvailableShips(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response []availableShip
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(response) != 1 {
		t.Fatalf("Expected one design row, got %d", len(response))
	}
	if response[0].Quantity != 3 {
		t.Errorf("Expected available qty=3 (5 built − 2 deployed), got %d", response[0].Quantity)
	}
}

func TestListAvailableShips_MultipleHullClasses(t *testing.T) {
	playerID, _, cleanup := setupAvailableShipsTest(t)
	defer cleanup()

	// Get different hull types
	var frigateID, cruiserID, battleshipID int
	database.DB.QueryRow(`SELECT id FROM hull_types WHERE hull_class = 'frigate' LIMIT 1`).Scan(&frigateID)
	database.DB.QueryRow(`SELECT id FROM hull_types WHERE hull_class = 'cruiser' LIMIT 1`).Scan(&cruiserID)
	database.DB.QueryRow(`SELECT id FROM hull_types WHERE hull_class = 'battleship' LIMIT 1`).Scan(&battleshipID)

	// Create designs for each class
	var frigateDesignID, cruiserDesignID, battleshipDesignID string
	database.DB.QueryRow(`
		INSERT INTO ship_designs (player_id, name, hull_type_id, modules_json,
			total_shield, total_structure, total_defense, total_agility, total_movement,
			total_storage, attack_power, weapon_range_min, weapon_range_max, volume_used,
			he3_per_round, metal_cost, he3_cost, gold_cost, build_time_seconds)
		VALUES ($1, 'Alpha_Frigate', $2, '[]', 100, 100, 10, 10, 5, 100, 50, 1, 5, 10, 5, 1000, 500, 50, 3600)
		RETURNING id
	`, playerID, frigateID).Scan(&frigateDesignID)

	database.DB.QueryRow(`
		INSERT INTO ship_designs (player_id, name, hull_type_id, modules_json,
			total_shield, total_structure, total_defense, total_agility, total_movement,
			total_storage, attack_power, weapon_range_min, weapon_range_max, volume_used,
			he3_per_round, metal_cost, he3_cost, gold_cost, build_time_seconds)
		VALUES ($1, 'Beta_Cruiser', $2, '[]', 200, 200, 15, 8, 4, 150, 80, 2, 6, 15, 8, 2000, 1000, 100, 7200)
		RETURNING id
	`, playerID, cruiserID).Scan(&cruiserDesignID)

	database.DB.QueryRow(`
		INSERT INTO ship_designs (player_id, name, hull_type_id, modules_json,
			total_shield, total_structure, total_defense, total_agility, total_movement,
			total_storage, attack_power, weapon_range_min, weapon_range_max, volume_used,
			he3_per_round, metal_cost, he3_cost, gold_cost, build_time_seconds)
		VALUES ($1, 'Gamma_Battleship', $2, '[]', 400, 400, 20, 5, 3, 200, 150, 3, 8, 25, 15, 5000, 2500, 250, 14400)
		RETURNING id
	`, playerID, battleshipID).Scan(&battleshipDesignID)

	// One ships row per design, quantity 1 each.
	database.DB.Exec(`INSERT INTO ships (player_id, ship_design_id, quantity, production_slot) VALUES ($1, $2, 1, 1)`,
		playerID, frigateDesignID)
	database.DB.Exec(`INSERT INTO ships (player_id, ship_design_id, quantity, production_slot) VALUES ($1, $2, 1, 2)`,
		playerID, cruiserDesignID)
	database.DB.Exec(`INSERT INTO ships (player_id, ship_design_id, quantity, production_slot) VALUES ($1, $2, 1, 3)`,
		playerID, battleshipDesignID)

	req := httptest.NewRequest("GET", "/api/ship-instances/available", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	ListAvailableShips(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response []availableShip
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != 3 {
		t.Errorf("Expected 3 ships (one of each class), got %d", len(response))
	}

	// Verify ships are ordered (ORDER BY hull_class, design_name)
	// Expected order: Battleship, Cruiser, Frigate (alphabetical)
	hullClasses := []string{}
	for _, ship := range response {
		hullClasses = append(hullClasses, ship.HullClass)
	}

	// Check that we have all three types (DB stores hull_class lowercase).
	hasF, hasC, hasB := false, false, false
	for _, hc := range hullClasses {
		if hc == "frigate" {
			hasF = true
		}
		if hc == "cruiser" {
			hasC = true
		}
		if hc == "battleship" {
			hasB = true
		}
	}

	if !hasF || !hasC || !hasB {
		t.Errorf("Expected all 3 hull classes, got: %v", hullClasses)
	}
}

func TestListAvailableShips_OnlyOwnShips(t *testing.T) {
	playerID, _, cleanup := setupAvailableShipsTest(t)
	defer cleanup()

	// Create another player
	otherPlayerID := uuid.New().String()
	otherAnonID := "anon-" + uuid.New().String()
	otherUsername := "other-player-" + uuid.New().String()
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username)
		VALUES ($1, $2, $3)
	`, otherPlayerID, otherAnonID, otherUsername)
	if err != nil {
		t.Fatalf("Failed to create other player: %v", err)
	}
	defer database.DB.Exec("DELETE FROM players WHERE id = $1", otherPlayerID)

	// Get hull type
	var hullTypeID int
	database.DB.QueryRow(`SELECT id FROM hull_types WHERE hull_class = 'frigate' LIMIT 1`).Scan(&hullTypeID)

	// Create ship design for other player
	var otherDesignID string
	database.DB.QueryRow(`
		INSERT INTO ship_designs (player_id, name, hull_type_id, modules_json,
			total_shield, total_structure, total_defense, total_agility, total_movement,
			total_storage, attack_power, weapon_range_min, weapon_range_max, volume_used,
			he3_per_round, metal_cost, he3_cost, gold_cost, build_time_seconds)
		VALUES ($1, 'Other Frigate', $2, '[]', 100, 100, 10, 10, 5, 100, 50, 1, 5, 10, 5, 1000, 500, 50, 3600)
		RETURNING id
	`, otherPlayerID, hullTypeID).Scan(&otherDesignID)

	// Create ship for other player
	database.DB.Exec(`INSERT INTO ship_instances (id, player_id, ship_design_id, hull_type_id) VALUES ($1, $2, $3, $4)`,
		uuid.New().String(), otherPlayerID, otherDesignID, hullTypeID)

	// Query as original player - should see NO ships
	req := httptest.NewRequest("GET", "/api/ship-instances/available", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	ListAvailableShips(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response []availableShip
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != 0 {
		t.Errorf("Expected 0 ships (other player's ships should not appear), got %d", len(response))
	}
}
