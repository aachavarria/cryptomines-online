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

func setupWarehouseTest(t *testing.T) (playerID, planetID, resourceID string, cleanup func()) {
	if err := database.InitSupabase(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	playerID = uuid.New().String()
	planetID = uuid.New().String()
	resourceID = uuid.New().String()
	anonID := "anon-" + uuid.New().String()
	username := "test-warehouse-" + uuid.New().String()

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
		INSERT INTO planets (id, player_id, name, position_x, position_y)
		VALUES ($1, $2, 'Test Planet', 1, 2)
	`, planetID, playerID)
	if err != nil {
		t.Fatalf("Failed to create test planet: %v", err)
	}

	// last_warehouse_update = now() so auto-accumulation contributes ~0 over
	// the test's elapsed wall-clock; the handler runs UpdateWarehouseResources
	// before reading.
	_, err = database.DB.Exec(`
		INSERT INTO resources (id, planet_id, metal, he3, gold,
			metal_per_hour, he3_per_hour, gold_per_hour, storage_capacity,
			warehouse_metal, warehouse_he3, warehouse_gold,
			last_warehouse_update, last_collected_at)
		VALUES ($1, $2, 1000, 500, 100, 1000, 500, 100, 100000,
			5000, 3000, 500, now(), now())
	`, resourceID, planetID)
	if err != nil {
		t.Fatalf("Failed to create test resources: %v", err)
	}

	cleanup = func() {
		database.DB.Exec("DELETE FROM resources WHERE id = $1", resourceID)
		database.DB.Exec("DELETE FROM planets WHERE id = $1", planetID)
		database.DB.Exec("DELETE FROM players WHERE id = $1", playerID)
	}

	return playerID, planetID, resourceID, cleanup
}

func TestCollectWarehouse_Success(t *testing.T) {
	playerID, _, resourceID, cleanup := setupWarehouseTest(t)
	defer cleanup()

	// Create request
	req := httptest.NewRequest("POST", "/api/resources/collect-warehouse", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	CollectWarehouse(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response collectResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify collected amounts
	if response.Collected.Metal != 5000 {
		t.Errorf("Expected collected metal 5000, got %d", response.Collected.Metal)
	}
	if response.Collected.He3 != 3000 {
		t.Errorf("Expected collected he3 3000, got %d", response.Collected.He3)
	}
	if response.Collected.Gold != 500 {
		t.Errorf("Expected collected gold 500, got %d", response.Collected.Gold)
	}

	// Verify new balances
	if response.Resources.Metal != 6000 {
		t.Errorf("Expected metal 6000, got %d", response.Resources.Metal)
	}
	if response.Resources.He3 != 3500 {
		t.Errorf("Expected he3 3500, got %d", response.Resources.He3)
	}
	if response.Resources.Gold != 600 {
		t.Errorf("Expected gold 600, got %d", response.Resources.Gold)
	}

	// Verify warehouse is now empty
	if response.Resources.WarehouseMetal != 0 {
		t.Errorf("Expected warehouse_metal 0, got %d", response.Resources.WarehouseMetal)
	}
	if response.Resources.WarehouseHe3 != 0 {
		t.Errorf("Expected warehouse_he3 0, got %d", response.Resources.WarehouseHe3)
	}
	if response.Resources.WarehouseGold != 0 {
		t.Errorf("Expected warehouse_gold 0, got %d", response.Resources.WarehouseGold)
	}

	// Verify database state
	var dbWarehouseMetal, dbWarehouseHe3, dbWarehouseGold int64
	err := database.DB.QueryRow(`
		SELECT warehouse_metal, warehouse_he3, warehouse_gold
		FROM resources WHERE id = $1
	`, resourceID).Scan(&dbWarehouseMetal, &dbWarehouseHe3, &dbWarehouseGold)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}

	if dbWarehouseMetal != 0 || dbWarehouseHe3 != 0 || dbWarehouseGold != 0 {
		t.Errorf("Expected warehouse empty in DB, got metal=%d, he3=%d, gold=%d",
			dbWarehouseMetal, dbWarehouseHe3, dbWarehouseGold)
	}
}

func TestCollectWarehouse_PartialCollection_CapacityLimit(t *testing.T) {
	playerID, _, resourceID, cleanup := setupWarehouseTest(t)
	defer cleanup()

	// Set resources near capacity. Pin last_warehouse_update to now() so the
	// handler's auto-accumulation pass adds 0 before the explicit collection.
	_, err := database.DB.Exec(`
		UPDATE resources
		SET metal = 98000, he3 = 99000, gold = 99500,
		    warehouse_metal = 5000, warehouse_he3 = 3000, warehouse_gold = 1000,
		    last_warehouse_update = now()
		WHERE id = $1
	`, resourceID)
	if err != nil {
		t.Fatalf("Failed to update resources: %v", err)
	}

	req := httptest.NewRequest("POST", "/api/resources/collect-warehouse", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	CollectWarehouse(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response collectResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Should only collect what fits
	// Metal: 98000 + ? <= 100000, so max 2000
	// He3: 99000 + ? <= 100000, so max 1000
	// Gold: 99500 + ? <= 100000, so max 500
	if response.Collected.Metal != 2000 {
		t.Errorf("Expected collected metal 2000 (capped), got %d", response.Collected.Metal)
	}
	if response.Collected.He3 != 1000 {
		t.Errorf("Expected collected he3 1000 (capped), got %d", response.Collected.He3)
	}
	if response.Collected.Gold != 500 {
		t.Errorf("Expected collected gold 500 (capped), got %d", response.Collected.Gold)
	}

	// Verify balances at capacity
	if response.Resources.Metal != 100000 {
		t.Errorf("Expected metal at capacity 100000, got %d", response.Resources.Metal)
	}
	if response.Resources.He3 != 100000 {
		t.Errorf("Expected he3 at capacity 100000, got %d", response.Resources.He3)
	}
	if response.Resources.Gold != 100000 {
		t.Errorf("Expected gold at capacity 100000, got %d", response.Resources.Gold)
	}

	// Verify remaining warehouse
	if response.Resources.WarehouseMetal != 3000 { // 5000 - 2000
		t.Errorf("Expected warehouse_metal 3000, got %d", response.Resources.WarehouseMetal)
	}
	if response.Resources.WarehouseHe3 != 2000 { // 3000 - 1000
		t.Errorf("Expected warehouse_he3 2000, got %d", response.Resources.WarehouseHe3)
	}
	if response.Resources.WarehouseGold != 500 { // 1000 - 500
		t.Errorf("Expected warehouse_gold 500, got %d", response.Resources.WarehouseGold)
	}
}

func TestCollectWarehouse_EmptyWarehouse(t *testing.T) {
	playerID, _, resourceID, cleanup := setupWarehouseTest(t)
	defer cleanup()

	// Empty the warehouse and freeze the auto-accumulation timestamp so the
	// handler doesn't refill it from the per_hour rates before collecting.
	_, err := database.DB.Exec(`
		UPDATE resources
		SET warehouse_metal = 0, warehouse_he3 = 0, warehouse_gold = 0,
		    last_warehouse_update = now()
		WHERE id = $1
	`, resourceID)
	if err != nil {
		t.Fatalf("Failed to update resources: %v", err)
	}

	req := httptest.NewRequest("POST", "/api/resources/collect-warehouse", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	CollectWarehouse(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	if response["error"] != "no resources to collect or storage full" {
		t.Errorf("Expected empty warehouse error, got: %s", response["error"])
	}
}

func TestCollectWarehouse_StorageFull(t *testing.T) {
	playerID, _, resourceID, cleanup := setupWarehouseTest(t)
	defer cleanup()

	_, err := database.DB.Exec(`
		UPDATE resources
		SET metal = 100000, he3 = 100000, gold = 100000,
		    warehouse_metal = 5000, warehouse_he3 = 3000, warehouse_gold = 500,
		    last_warehouse_update = now()
		WHERE id = $1
	`, resourceID)
	if err != nil {
		t.Fatalf("Failed to update resources: %v", err)
	}

	req := httptest.NewRequest("POST", "/api/resources/collect-warehouse", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	CollectWarehouse(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	if response["error"] != "no resources to collect or storage full" {
		t.Errorf("Expected storage full error, got: %s", response["error"])
	}
}

func TestGetResources_IncludesWarehouseFields(t *testing.T) {
	playerID, planetID, _, cleanup := setupWarehouseTest(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/planets/"+planetID+"/resources", nil)
	req.SetPathValue("id", planetID)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	GetResources(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify warehouse fields are present
	if _, ok := response["warehouse_metal"]; !ok {
		t.Error("Expected warehouse_metal field in response")
	}
	if _, ok := response["warehouse_he3"]; !ok {
		t.Error("Expected warehouse_he3 field in response")
	}
	if _, ok := response["warehouse_gold"]; !ok {
		t.Error("Expected warehouse_gold field in response")
	}
	if _, ok := response["last_warehouse_update"]; !ok {
		t.Error("Expected last_warehouse_update field in response")
	}

	// Verify warehouse values
	if response["warehouse_metal"].(float64) != 5000 {
		t.Errorf("Expected warehouse_metal 5000, got %v", response["warehouse_metal"])
	}
	if response["warehouse_he3"].(float64) != 3000 {
		t.Errorf("Expected warehouse_he3 3000, got %v", response["warehouse_he3"])
	}
	if response["warehouse_gold"].(float64) != 500 {
		t.Errorf("Expected warehouse_gold 500, got %v", response["warehouse_gold"])
	}
}

func TestCollectWarehouse_TransactionRollback(t *testing.T) {
	playerID, _, resourceID, cleanup := setupWarehouseTest(t)
	defer cleanup()

	// Get initial state
	var initialMetal, initialWarehouse int64
	err := database.DB.QueryRow(`
		SELECT metal, warehouse_metal FROM resources WHERE id = $1
	`, resourceID).Scan(&initialMetal, &initialWarehouse)
	if err != nil {
		t.Fatalf("Failed to get initial state: %v", err)
	}

	// Create request
	req := httptest.NewRequest("POST", "/api/resources/collect-warehouse", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	CollectWarehouse(w, req)

	// Verify state changed (transaction committed)
	var finalMetal, finalWarehouse int64
	err = database.DB.QueryRow(`
		SELECT metal, warehouse_metal FROM resources WHERE id = $1
	`, resourceID).Scan(&finalMetal, &finalWarehouse)
	if err != nil {
		t.Fatalf("Failed to get final state: %v", err)
	}

	// Metal should have increased
	if finalMetal <= initialMetal {
		t.Errorf("Expected metal to increase from %d, got %d", initialMetal, finalMetal)
	}

	// Warehouse should have decreased
	if finalWarehouse >= initialWarehouse {
		t.Errorf("Expected warehouse to decrease from %d, got %d", initialWarehouse, finalWarehouse)
	}
}
