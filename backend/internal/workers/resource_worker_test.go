package workers

import (
	"testing"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/google/uuid"
)

func setupResourceWorkerTest(t *testing.T) (playerID, planetID, resourceID string, cleanup func()) {
	if err := database.InitSupabase(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	playerID = uuid.New().String()
	planetID = uuid.New().String()
	resourceID = uuid.New().String()
	anonID := "anon-" + uuid.New().String()
	username := "test-res-worker-" + uuid.New().String()

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
		INSERT INTO planets (id, player_id, name, is_homeworld, position_x, position_y)
		VALUES ($1, $2, 'Test Planet', true, $3, $4)
	`, planetID, playerID, 1000+time.Now().Unix()%1000, 1000+time.Now().Unix()%1000)
	if err != nil {
		t.Fatalf("Failed to create test planet: %v", err)
	}

	// Create resources manually (not auto-created by trigger in test environment)
	_, err = database.DB.Exec(`
		INSERT INTO resources (id, planet_id)
		VALUES ($1, $2)
	`, resourceID, planetID)
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	cleanup = func() {
		database.DB.Exec("DELETE FROM resources WHERE planet_id = $1", planetID)
		database.DB.Exec("DELETE FROM planets WHERE id = $1", planetID)
		database.DB.Exec("DELETE FROM players WHERE id = $1", playerID)
	}

	return playerID, planetID, resourceID, cleanup
}

func TestUpdateWarehouseResources_Basic(t *testing.T) {
	_, _, resourceID, cleanup := setupResourceWorkerTest(t)
	defer cleanup()

	// Set production rates and simulate 1 hour has passed
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	_, err := database.DB.Exec(`
		UPDATE resources
		SET metal_per_hour = 1000,
		    he3_per_hour = 500,
		    gold_per_hour = 100,
		    warehouse_metal = 0,
		    warehouse_he3 = 0,
		    warehouse_gold = 0,
		    last_warehouse_update = $1
		WHERE id = $2
	`, oneHourAgo, resourceID)
	if err != nil {
		t.Fatalf("Failed to setup resource: %v", err)
	}

	// Run the worker
	UpdateWarehouseResources()

	// Verify warehouse was updated
	var warehouseMetal, warehouseHe3, warehouseGold int64
	err = database.DB.QueryRow(`
		SELECT warehouse_metal, warehouse_he3, warehouse_gold
		FROM resources WHERE id = $1
	`, resourceID).Scan(&warehouseMetal, &warehouseHe3, &warehouseGold)
	if err != nil {
		t.Fatalf("Failed to query updated resource: %v", err)
	}

	// Allow for small timing differences (should be close to 1 hour of production)
	if warehouseMetal < 950 || warehouseMetal > 1050 {
		t.Errorf("Expected warehouse_metal ~1000, got %d", warehouseMetal)
	}
	if warehouseHe3 < 475 || warehouseHe3 > 525 {
		t.Errorf("Expected warehouse_he3 ~500, got %d", warehouseHe3)
	}
	if warehouseGold < 95 || warehouseGold > 105 {
		t.Errorf("Expected warehouse_gold ~100, got %d", warehouseGold)
	}
}

func TestUpdateWarehouseResources_CapacityCap(t *testing.T) {
	_, _, resourceID, cleanup := setupResourceWorkerTest(t)
	defer cleanup()

	// Set high production rates and low capacity to test capping
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	_, err := database.DB.Exec(`
		UPDATE resources
		SET metal_per_hour = 100000,
		    he3_per_hour = 100000,
		    gold_per_hour = 100000,
		    metal = 50000,
		    he3 = 60000,
		    gold = 70000,
		    storage_capacity = 100000,
		    warehouse_metal = 0,
		    warehouse_he3 = 0,
		    warehouse_gold = 0,
		    last_warehouse_update = $1
		WHERE id = $2
	`, oneHourAgo, resourceID)
	if err != nil {
		t.Fatalf("Failed to setup resource: %v", err)
	}

	// Run the worker
	UpdateWarehouseResources()

	// Verify warehouse was capped at storage_capacity - current_resources
	var warehouseMetal, warehouseHe3, warehouseGold int64
	err = database.DB.QueryRow(`
		SELECT warehouse_metal, warehouse_he3, warehouse_gold
		FROM resources WHERE id = $1
	`, resourceID).Scan(&warehouseMetal, &warehouseHe3, &warehouseGold)
	if err != nil {
		t.Fatalf("Failed to query updated resource: %v", err)
	}

	// Warehouse should be capped at (capacity - current)
	// Metal: 100000 - 50000 = 50000
	// He3: 100000 - 60000 = 40000
	// Gold: 100000 - 70000 = 30000
	if warehouseMetal != 50000 {
		t.Errorf("Expected warehouse_metal = 50000 (capacity cap), got %d", warehouseMetal)
	}
	if warehouseHe3 != 40000 {
		t.Errorf("Expected warehouse_he3 = 40000 (capacity cap), got %d", warehouseHe3)
	}
	if warehouseGold != 30000 {
		t.Errorf("Expected warehouse_gold = 30000 (capacity cap), got %d", warehouseGold)
	}
}

func TestUpdateWarehouseResources_Accumulation(t *testing.T) {
	_, _, resourceID, cleanup := setupResourceWorkerTest(t)
	defer cleanup()

	// Set production rates and existing warehouse amounts
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	_, err := database.DB.Exec(`
		UPDATE resources
		SET metal_per_hour = 1000,
		    he3_per_hour = 500,
		    gold_per_hour = 100,
		    warehouse_metal = 500,
		    warehouse_he3 = 250,
		    warehouse_gold = 50,
		    last_warehouse_update = $1
		WHERE id = $2
	`, oneHourAgo, resourceID)
	if err != nil {
		t.Fatalf("Failed to setup resource: %v", err)
	}

	// Run the worker
	UpdateWarehouseResources()

	// Verify warehouse accumulated (existing + new production)
	var warehouseMetal, warehouseHe3, warehouseGold int64
	err = database.DB.QueryRow(`
		SELECT warehouse_metal, warehouse_he3, warehouse_gold
		FROM resources WHERE id = $1
	`, resourceID).Scan(&warehouseMetal, &warehouseHe3, &warehouseGold)
	if err != nil {
		t.Fatalf("Failed to query updated resource: %v", err)
	}

	// Should be: existing + ~1 hour of production
	// Metal: 500 + 1000 = 1500
	// He3: 250 + 500 = 750
	// Gold: 50 + 100 = 150
	if warehouseMetal < 1450 || warehouseMetal > 1550 {
		t.Errorf("Expected warehouse_metal ~1500, got %d", warehouseMetal)
	}
	if warehouseHe3 < 725 || warehouseHe3 > 775 {
		t.Errorf("Expected warehouse_he3 ~750, got %d", warehouseHe3)
	}
	if warehouseGold < 145 || warehouseGold > 155 {
		t.Errorf("Expected warehouse_gold ~150, got %d", warehouseGold)
	}
}

func TestUpdateWarehouseResources_ZeroProduction(t *testing.T) {
	_, _, resourceID, cleanup := setupResourceWorkerTest(t)
	defer cleanup()

	// Set zero production rates
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	_, err := database.DB.Exec(`
		UPDATE resources
		SET metal_per_hour = 0,
		    he3_per_hour = 0,
		    gold_per_hour = 0,
		    warehouse_metal = 100,
		    warehouse_he3 = 100,
		    warehouse_gold = 100,
		    last_warehouse_update = $1
		WHERE id = $2
	`, oneHourAgo, resourceID)
	if err != nil {
		t.Fatalf("Failed to setup resource: %v", err)
	}

	// Run the worker
	UpdateWarehouseResources()

	// Verify warehouse unchanged (no production)
	var warehouseMetal, warehouseHe3, warehouseGold int64
	err = database.DB.QueryRow(`
		SELECT warehouse_metal, warehouse_he3, warehouse_gold
		FROM resources WHERE id = $1
	`, resourceID).Scan(&warehouseMetal, &warehouseHe3, &warehouseGold)
	if err != nil {
		t.Fatalf("Failed to query updated resource: %v", err)
	}

	if warehouseMetal != 100 {
		t.Errorf("Expected warehouse_metal = 100 (unchanged), got %d", warehouseMetal)
	}
	if warehouseHe3 != 100 {
		t.Errorf("Expected warehouse_he3 = 100 (unchanged), got %d", warehouseHe3)
	}
	if warehouseGold != 100 {
		t.Errorf("Expected warehouse_gold = 100 (unchanged), got %d", warehouseGold)
	}
}

func TestUpdateWarehouseResources_RecentUpdate(t *testing.T) {
	_, _, resourceID, cleanup := setupResourceWorkerTest(t)
	defer cleanup()

	// Set last_warehouse_update to now (no elapsed time)
	_, err := database.DB.Exec(`
		UPDATE resources
		SET metal_per_hour = 1000,
		    he3_per_hour = 500,
		    gold_per_hour = 100,
		    warehouse_metal = 0,
		    warehouse_he3 = 0,
		    warehouse_gold = 0,
		    last_warehouse_update = now()
		WHERE id = $1
	`, resourceID)
	if err != nil {
		t.Fatalf("Failed to setup resource: %v", err)
	}

	// Run the worker
	UpdateWarehouseResources()

	// Verify warehouse is still ~0 (minimal elapsed time)
	var warehouseMetal, warehouseHe3, warehouseGold int64
	err = database.DB.QueryRow(`
		SELECT warehouse_metal, warehouse_he3, warehouse_gold
		FROM resources WHERE id = $1
	`, resourceID).Scan(&warehouseMetal, &warehouseHe3, &warehouseGold)
	if err != nil {
		t.Fatalf("Failed to query updated resource: %v", err)
	}

	// Should be very close to 0 (only a few milliseconds elapsed)
	if warehouseMetal > 10 {
		t.Errorf("Expected warehouse_metal ~0 (no time elapsed), got %d", warehouseMetal)
	}
	if warehouseHe3 > 10 {
		t.Errorf("Expected warehouse_he3 ~0 (no time elapsed), got %d", warehouseHe3)
	}
	if warehouseGold > 10 {
		t.Errorf("Expected warehouse_gold ~0 (no time elapsed), got %d", warehouseGold)
	}
}

func TestUpdateWarehouseResources_TimestampUpdate(t *testing.T) {
	_, _, resourceID, cleanup := setupResourceWorkerTest(t)
	defer cleanup()

	// Set last_warehouse_update to 1 hour ago
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	_, err := database.DB.Exec(`
		UPDATE resources
		SET metal_per_hour = 1000,
		    last_warehouse_update = $1
		WHERE id = $2
	`, oneHourAgo, resourceID)
	if err != nil {
		t.Fatalf("Failed to setup resource: %v", err)
	}

	beforeUpdate := time.Now()

	// Run the worker
	UpdateWarehouseResources()

	// Verify last_warehouse_update was updated to ~now
	var lastWarehouseUpdate time.Time
	err = database.DB.QueryRow(`
		SELECT last_warehouse_update FROM resources WHERE id = $1
	`, resourceID).Scan(&lastWarehouseUpdate)
	if err != nil {
		t.Fatalf("Failed to query updated resource: %v", err)
	}

	// Should be within a few seconds of now
	if lastWarehouseUpdate.Before(beforeUpdate) {
		t.Errorf("last_warehouse_update was not updated (still %v, should be after %v)", lastWarehouseUpdate, beforeUpdate)
	}
	if time.Since(lastWarehouseUpdate) > 5*time.Second {
		t.Errorf("last_warehouse_update is too old (%v), should be recent", lastWarehouseUpdate)
	}
}
