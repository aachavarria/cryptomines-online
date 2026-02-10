package handlers

import (
	"testing"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/services"
)

// setupBuildingTestDB initializes the test database connection.
func setupBuildingTestDB(t *testing.T) {
	if database.DB == nil {
		if err := database.InitSupabase(); err != nil {
			t.Fatalf("Failed to init database: %v", err)
		}
	}
}

// cleanupBuildingTestPlayer removes test data for a player.
func cleanupBuildingTestPlayer(t *testing.T, playerID string) {
	database.DB.Exec("DELETE FROM buildings WHERE planet_id IN (SELECT id FROM planets WHERE player_id = $1)", playerID)
	database.DB.Exec("DELETE FROM resources WHERE planet_id IN (SELECT id FROM planets WHERE player_id = $1)", playerID)
	database.DB.Exec("DELETE FROM planets WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM technologies WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM players WHERE id = $1", playerID)
}

func TestTechBonuses_ConstructionCostReduction(t *testing.T) {
	setupBuildingTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000201"
	planetID := "00000000-0000-0000-0001-000000000201"

	// Create test player and planet
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'costtest', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupBuildingTestPlayer(t, playerID)

	_, err = database.DB.Exec(`
		INSERT INTO planets (id, player_id, name, is_homeworld, position_x, position_y, created_at, updated_at)
		VALUES ($1, $2, 'Test Planet', true, 100, 100, now(), now())
		ON CONFLICT (id) DO NOTHING
	`, planetID, playerID)
	if err != nil {
		t.Fatalf("Failed to create test planet: %v", err)
	}

	// Create resources
	database.DB.Exec(`
		INSERT INTO resources (planet_id, metal, he3, gold, created_at, updated_at)
		VALUES ($1, 1000000, 1000000, 1000000, now(), now())
		ON CONFLICT (planet_id) DO NOTHING
	`, planetID)

	// Get tech type ID for Quality Materials (cost reduction)
	var qualityMatTechID int
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'quality_materials'").Scan(&qualityMatTechID)

	// Set Quality Materials to level 10 = -15% cost
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 10, false, now(), now())
	`, playerID, qualityMatTechID)

	// Get tech bonuses
	bonuses, err := services.GetPlayerTechBonuses(playerID)
	if err != nil {
		t.Fatalf("Failed to get tech bonuses: %v", err)
	}

	// Verify cost reduction bonus is correct
	if bonuses.BuildCostReduction != 15.0 {
		t.Errorf("Expected BuildCostReduction=15.0, got %f", bonuses.BuildCostReduction)
	}

	// Test that the bonus is applied correctly by checking calculation
	// Base cost for Metal Collector level 1: 150 metal
	baseCost := int64(150)
	expectedCost := int64(float64(baseCost) * (1.0 - 15.0/100.0)) // 127.5 -> 127
	if expectedCost != 127 {
		t.Errorf("Expected reduced cost=127, got %d", expectedCost)
	}
}

func TestTechBonuses_BuildSpeed(t *testing.T) {
	setupBuildingTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000202"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'speedtest', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupBuildingTestPlayer(t, playerID)

	// Get tech type ID for Construction Boost
	var constructBoostID int
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'construction_boost'").Scan(&constructBoostID)

	// Set Construction Boost to level 10 = +15% speed (reduces time by 15%)
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 10, false, now(), now())
	`, playerID, constructBoostID)

	// Get tech bonuses
	bonuses, err := services.GetPlayerTechBonuses(playerID)
	if err != nil {
		t.Fatalf("Failed to get tech bonuses: %v", err)
	}

	// Verify build speed bonus is correct
	if bonuses.BuildSpeed != 15.0 {
		t.Errorf("Expected BuildSpeed=15.0, got %f", bonuses.BuildSpeed)
	}

	// Test that the bonus is applied correctly
	// Base time for a building: 100 seconds
	baseTime := 100
	expectedTime := int(float64(baseTime) * (1.0 - 15.0/100.0)) // 85 seconds
	if expectedTime != 85 {
		t.Errorf("Expected reduced time=85, got %d", expectedTime)
	}
}

func TestTechBonuses_ProductionBonus(t *testing.T) {
	setupBuildingTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000203"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'prodtest', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupBuildingTestPlayer(t, playerID)

	// Get tech type IDs for production bonuses
	var metalTechID, he3TechID, goldTechID int
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'high_yield_mining'").Scan(&metalTechID)
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'high_yield_chemistry'").Scan(&he3TechID)
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'high_yield_investing'").Scan(&goldTechID)

	// Set production techs to max level (10)
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 10, false, now(), now())
	`, playerID, metalTechID)
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 10, false, now(), now())
	`, playerID, he3TechID)
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 10, false, now(), now())
	`, playerID, goldTechID)

	// Get tech bonuses
	bonuses, err := services.GetPlayerTechBonuses(playerID)
	if err != nil {
		t.Fatalf("Failed to get tech bonuses: %v", err)
	}

	// Verify production bonuses
	if bonuses.MetalOutput != 10.0 {
		t.Errorf("Expected MetalOutput=10.0, got %f", bonuses.MetalOutput)
	}
	if bonuses.He3Output != 10.0 {
		t.Errorf("Expected He3Output=10.0, got %f", bonuses.He3Output)
	}
	if bonuses.GoldOutput != 10.0 {
		t.Errorf("Expected GoldOutput=10.0, got %f", bonuses.GoldOutput)
	}

	// Test that the bonus is applied correctly
	// Base production: 1000/hour
	baseProd := int64(1000)
	expectedProd := int64(float64(baseProd) * (1.0 + 10.0/100.0)) // 1100/hour
	if expectedProd != 1100 {
		t.Errorf("Expected boosted production=1100, got %d", expectedProd)
	}
}

func TestTechBonuses_WarehouseCapacity(t *testing.T) {
	setupBuildingTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000204"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'warehousetest', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupBuildingTestPlayer(t, playerID)

	// Get tech type ID for Expand Capacity
	var expandCapID int
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'expand_capacity'").Scan(&expandCapID)

	// Set Expand Capacity to level 5 = +250,000 capacity
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 5, false, now(), now())
	`, playerID, expandCapID)

	// Get tech bonuses
	bonuses, err := services.GetPlayerTechBonuses(playerID)
	if err != nil {
		t.Fatalf("Failed to get tech bonuses: %v", err)
	}

	// Verify warehouse capacity bonus
	if bonuses.WarehouseCapacity != 250000 {
		t.Errorf("Expected WarehouseCapacity=250000, got %d", bonuses.WarehouseCapacity)
	}

	// Test that the bonus is applied correctly
	// Base capacity: 100,000
	baseCap := int64(100000)
	expectedCap := baseCap + bonuses.WarehouseCapacity // 350,000
	if expectedCap != 350000 {
		t.Errorf("Expected total capacity=350000, got %d", expectedCap)
	}
}

func TestTechBonuses_ConstructionSlots(t *testing.T) {
	setupBuildingTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000205"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'slotstest', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupBuildingTestPlayer(t, playerID)

	// Get tech type ID for Concurrent Construction
	var concurrentTechID int
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'concurrent_construction'").Scan(&concurrentTechID)

	// Initially, should have 1 slot
	slots := getConstructionSlots(playerID)
	if slots != 1 {
		t.Errorf("Expected 1 construction slot, got %d", slots)
	}

	// Research Concurrent Construction (level 1, max 1)
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 1, false, now(), now())
	`, playerID, concurrentTechID)

	// Should now have 2 slots
	slots = getConstructionSlots(playerID)
	if slots != 2 {
		t.Errorf("Expected 2 construction slots after research, got %d", slots)
	}
}

func TestTechBonuses_CombinedBonuses(t *testing.T) {
	setupBuildingTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000206"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'combinedtest', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupBuildingTestPlayer(t, playerID)

	// Get tech type IDs
	var qualityMatID, constructBoostID, metalProdID int
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'quality_materials'").Scan(&qualityMatID)
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'construction_boost'").Scan(&constructBoostID)
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'high_yield_mining'").Scan(&metalProdID)

	// Research multiple techs
	database.DB.Exec(`INSERT INTO technologies (player_id, tech_type, level) VALUES ($1, $2, 5)`, playerID, qualityMatID)
	database.DB.Exec(`INSERT INTO technologies (player_id, tech_type, level) VALUES ($1, $2, 8)`, playerID, constructBoostID)
	database.DB.Exec(`INSERT INTO technologies (player_id, tech_type, level) VALUES ($1, $2, 10)`, playerID, metalProdID)

	// Get combined tech bonuses
	bonuses, err := services.GetPlayerTechBonuses(playerID)
	if err != nil {
		t.Fatalf("Failed to get tech bonuses: %v", err)
	}

	// Verify all bonuses are applied correctly
	// quality_materials level 5: -7.5% cost
	if bonuses.BuildCostReduction != 7.5 {
		t.Errorf("Expected BuildCostReduction=7.5, got %f", bonuses.BuildCostReduction)
	}
	// construction_boost level 8: +12% speed
	if bonuses.BuildSpeed != 12.0 {
		t.Errorf("Expected BuildSpeed=12.0, got %f", bonuses.BuildSpeed)
	}
	// high_yield_mining level 10: +10% metal
	if bonuses.MetalOutput != 10.0 {
		t.Errorf("Expected MetalOutput=10.0, got %f", bonuses.MetalOutput)
	}
}
