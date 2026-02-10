package handlers

import (
	"testing"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/services"
)

// setupShipTestDB initializes the test database connection.
func setupShipTestDB(t *testing.T) {
	if database.DB == nil {
		if err := database.InitSupabase(); err != nil {
			t.Fatalf("Failed to init database: %v", err)
		}
	}
}

// cleanupShipTestPlayer removes test data for a player.
func cleanupShipTestPlayer(t *testing.T, playerID string) {
	database.DB.Exec("DELETE FROM ships WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM ship_designs WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM technologies WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM players WHERE id = $1", playerID)
}

func TestTechBonuses_ShipCostReduction(t *testing.T) {
	setupShipTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000301"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'shipcosttest', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupShipTestPlayer(t, playerID)

	// Get tech type ID for Ship Building Logistics (ship cost reduction)
	var shipLogisticsID int
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'ship_building_logistics'").Scan(&shipLogisticsID)

	// Set Ship Building Logistics to level 10 = -15% cost
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 10, false, now(), now())
	`, playerID, shipLogisticsID)

	// Get tech bonuses
	bonuses, err := services.GetPlayerTechBonuses(playerID)
	if err != nil {
		t.Fatalf("Failed to get tech bonuses: %v", err)
	}

	// Verify ship cost reduction bonus
	if bonuses.ShipBuildCostReduction != 15.0 {
		t.Errorf("Expected ShipBuildCostReduction=15.0, got %f", bonuses.ShipBuildCostReduction)
	}

	// Test cost reduction calculation
	baseCost := int64(1000)
	costReduction := bonuses.ShipBuildCostReduction / 100.0
	effectiveCost := int64(float64(baseCost) * (1.0 - costReduction))
	expectedCost := int64(850) // 1000 * 0.85 = 850

	if effectiveCost != expectedCost {
		t.Errorf("Expected effective cost=%d, got %d", expectedCost, effectiveCost)
	}
}

func TestTechBonuses_ShipBuildSpeed(t *testing.T) {
	setupShipTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000302"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'shipspeedtest', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupShipTestPlayer(t, playerID)

	// Get tech type ID for Ship Building Boost
	var shipBoostID int
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'ship_building_boost'").Scan(&shipBoostID)

	// Set Ship Building Boost to level 10 = +15% speed
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 10, false, now(), now())
	`, playerID, shipBoostID)

	// Get tech bonuses
	bonuses, err := services.GetPlayerTechBonuses(playerID)
	if err != nil {
		t.Fatalf("Failed to get tech bonuses: %v", err)
	}

	// Verify ship build speed bonus
	if bonuses.ShipBuildSpeed != 15.0 {
		t.Errorf("Expected ShipBuildSpeed=15.0, got %f", bonuses.ShipBuildSpeed)
	}

	// Test speed bonus with factory bonus (cumulative)
	factoryBonus := 10 // Factory level 5 gives 10% bonus
	techBonus := int(bonuses.ShipBuildSpeed)
	totalBonus := factoryBonus + techBonus // 25% total

	baseTime := 100
	effectiveTime := services.ShipBuildTime(baseTime, totalBonus, 1)
	expectedTime := 75 // 100 * (1 - 0.25) = 75

	if effectiveTime != expectedTime {
		t.Errorf("Expected effective time=%d, got %d", expectedTime, effectiveTime)
	}
}

func TestTechBonuses_ShipProductionSlots(t *testing.T) {
	setupShipTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000303"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'shipslotstest', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupShipTestPlayer(t, playerID)

	// Factory level 5 gives 2 slots (base)
	factoryLevel := 5
	baseSlots := getShipProductionSlots(playerID, factoryLevel)
	if baseSlots != 2 {
		t.Errorf("Expected 2 base slots at factory level 5, got %d", baseSlots)
	}

	// Get tech type ID for Sync Shipbuilding
	var syncShipID int
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'sync_shipbuilding'").Scan(&syncShipID)

	// Research Sync Shipbuilding (level 1, max 1) - adds 5th slot
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 1, false, now(), now())
	`, playerID, syncShipID)

	// Should now have 3 slots (2 base + 1 tech)
	totalSlots := getShipProductionSlots(playerID, factoryLevel)
	if totalSlots != 3 {
		t.Errorf("Expected 3 total slots after Sync Shipbuilding, got %d", totalSlots)
	}
}

func TestTechBonuses_ShipCombined(t *testing.T) {
	setupShipTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000304"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'shipcombinedtest', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupShipTestPlayer(t, playerID)

	// Get tech type IDs
	var shipBoostID, shipLogisticsID, syncShipID int
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'ship_building_boost'").Scan(&shipBoostID)
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'ship_building_logistics'").Scan(&shipLogisticsID)
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'sync_shipbuilding'").Scan(&syncShipID)

	// Research all ship techs
	database.DB.Exec(`INSERT INTO technologies (player_id, tech_type, level) VALUES ($1, $2, 10)`, playerID, shipBoostID)
	database.DB.Exec(`INSERT INTO technologies (player_id, tech_type, level) VALUES ($1, $2, 8)`, playerID, shipLogisticsID)
	database.DB.Exec(`INSERT INTO technologies (player_id, tech_type, level) VALUES ($1, $2, 1)`, playerID, syncShipID)

	// Get combined tech bonuses
	bonuses, err := services.GetPlayerTechBonuses(playerID)
	if err != nil {
		t.Fatalf("Failed to get tech bonuses: %v", err)
	}

	// Verify all bonuses
	// ship_building_boost level 10: +15% speed
	if bonuses.ShipBuildSpeed != 15.0 {
		t.Errorf("Expected ShipBuildSpeed=15.0, got %f", bonuses.ShipBuildSpeed)
	}
	// ship_building_logistics level 8: -12% cost
	if bonuses.ShipBuildCostReduction != 12.0 {
		t.Errorf("Expected ShipBuildCostReduction=12.0, got %f", bonuses.ShipBuildCostReduction)
	}

	// Check slots bonus
	factoryLevel := 5 // 2 base slots
	totalSlots := getShipProductionSlots(playerID, factoryLevel)
	if totalSlots != 3 {
		t.Errorf("Expected 3 total slots (2 base + 1 tech), got %d", totalSlots)
	}

	// Test combined cost and speed impact
	baseCost := int64(10000)
	baseTime := 1000

	// Apply cost reduction
	costReduction := bonuses.ShipBuildCostReduction / 100.0
	effectiveCost := int64(float64(baseCost) * (1.0 - costReduction))
	expectedCost := int64(8800) // 10000 * 0.88 = 8800

	// Apply speed bonus (with factory bonus)
	factoryBonus := 10 // Factory level 5
	totalSpeedBonus := factoryBonus + int(bonuses.ShipBuildSpeed) // 25% total
	effectiveTime := services.ShipBuildTime(baseTime, totalSpeedBonus, 1)
	expectedTime := 750 // 1000 * (1 - 0.25) = 750

	if effectiveCost != expectedCost {
		t.Errorf("Expected combined cost=%d, got %d", expectedCost, effectiveCost)
	}
	if effectiveTime != expectedTime {
		t.Errorf("Expected combined time=%d, got %d", expectedTime, effectiveTime)
	}
}

func TestTechBonuses_ShipRepairPercent(t *testing.T) {
	setupShipTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000305"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'shiprepairtest', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupShipTestPlayer(t, playerID)

	// Get tech type ID for Repair Technology
	var repairTechID int
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'repair_technology'").Scan(&repairTechID)

	// Set Repair Technology to level 10 = +10% repair
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 10, false, now(), now())
	`, playerID, repairTechID)

	// Get tech bonuses
	bonuses, err := services.GetPlayerTechBonuses(playerID)
	if err != nil {
		t.Fatalf("Failed to get tech bonuses: %v", err)
	}

	// Verify ship repair percent bonus
	if bonuses.ShipRepairPercent != 10.0 {
		t.Errorf("Expected ShipRepairPercent=10.0, got %f", bonuses.ShipRepairPercent)
	}
}

func TestGetShipProductionSlots_NoFactory(t *testing.T) {
	setupShipTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000306"

	// Create test player (no factory)
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'nofactorytest', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupShipTestPlayer(t, playerID)

	// Should default to 1 slot if no factory data
	slots := getShipProductionSlots(playerID, 0)
	if slots != 1 {
		t.Errorf("Expected 1 default slot, got %d", slots)
	}
}

func TestGetShipProductionSlots_MaxCapEnforced(t *testing.T) {
	setupShipTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000307"

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'maxslotstest', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupShipTestPlayer(t, playerID)

	// Factory level 12 gives 4 slots, tech gives +1 = 5 total (hard cap)
	var syncShipID int
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'sync_shipbuilding'").Scan(&syncShipID)
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level) VALUES ($1, $2, 1)
	`, playerID, syncShipID)

	totalSlots := getShipProductionSlots(playerID, 12)
	if totalSlots != 5 {
		t.Errorf("Expected hard cap of 5 slots, got %d", totalSlots)
	}
}
