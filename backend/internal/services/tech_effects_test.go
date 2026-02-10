package services

import (
	"encoding/json"
	"testing"

	"github.com/cryptomines-online/backend/internal/database"
)

// setupTestDB initializes the test database connection.
// Tests should use the actual Supabase local instance for integration testing.
func setupTestDB(t *testing.T) {
	if database.DB == nil {
		if err := database.InitSupabase(); err != nil {
			t.Fatalf("Failed to init database: %v", err)
		}
	}
}

// cleanupTestPlayer removes test data for a player.
func cleanupTestPlayer(t *testing.T, playerID string) {
	database.DB.Exec("DELETE FROM technologies WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM players WHERE id = $1", playerID)
}

func TestGetPlayerTechBonuses_NoResearch(t *testing.T) {
	setupTestDB(t)

	// Create test player
	playerID := "00000000-0000-0000-0000-000000000001"
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'testuser1', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupTestPlayer(t, playerID)

	bonuses, err := GetPlayerTechBonuses(playerID)
	if err != nil {
		t.Fatalf("GetPlayerTechBonuses failed: %v", err)
	}

	// All bonuses should be zero
	if bonuses.MetalOutput != 0 {
		t.Errorf("Expected MetalOutput=0, got %f", bonuses.MetalOutput)
	}
	if bonuses.ConstructionSlots != 0 {
		t.Errorf("Expected ConstructionSlots=0, got %d", bonuses.ConstructionSlots)
	}
	if bonuses.BuildSpeed != 0 {
		t.Errorf("Expected BuildSpeed=0, got %f", bonuses.BuildSpeed)
	}
}

func TestGetPlayerTechBonuses_ProductionBonuses(t *testing.T) {
	setupTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000002"
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'testuser2', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupTestPlayer(t, playerID)

	// Get tech type IDs
	var metalTechID, he3TechID, goldTechID int
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'high_yield_mining'").Scan(&metalTechID)
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'high_yield_chemistry'").Scan(&he3TechID)
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'high_yield_investing'").Scan(&goldTechID)

	// Insert completed research
	// high_yield_mining: +1% metal per level, level 5 = +5%
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 5, false, now(), now())
	`, playerID, metalTechID)

	// high_yield_chemistry: +1% He3 per level, level 3 = +3%
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 3, false, now(), now())
	`, playerID, he3TechID)

	// high_yield_investing: +1% Gold per level, level 7 = +7%
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 7, false, now(), now())
	`, playerID, goldTechID)

	bonuses, err := GetPlayerTechBonuses(playerID)
	if err != nil {
		t.Fatalf("GetPlayerTechBonuses failed: %v", err)
	}

	// Check production bonuses
	if bonuses.MetalOutput != 5.0 {
		t.Errorf("Expected MetalOutput=5.0, got %f", bonuses.MetalOutput)
	}
	if bonuses.He3Output != 3.0 {
		t.Errorf("Expected He3Output=3.0, got %f", bonuses.He3Output)
	}
	if bonuses.GoldOutput != 7.0 {
		t.Errorf("Expected GoldOutput=7.0, got %f", bonuses.GoldOutput)
	}
}

func TestGetPlayerTechBonuses_ConstructionBonuses(t *testing.T) {
	setupTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000003"
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'testuser3', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupTestPlayer(t, playerID)

	// Get tech type IDs
	var concurrentTechID, buildBoostID, shipBoostID, syncShipID int
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'concurrent_construction'").Scan(&concurrentTechID)
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'construction_boost'").Scan(&buildBoostID)
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'ship_building_boost'").Scan(&shipBoostID)
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'sync_shipbuilding'").Scan(&syncShipID)

	// concurrent_construction: +1 slot, max level 1
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 1, false, now(), now())
	`, playerID, concurrentTechID)

	// construction_boost: +1.5% per level, level 10 = +15%
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 10, false, now(), now())
	`, playerID, buildBoostID)

	// ship_building_boost: +1.5% per level, level 8 = +12%
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 8, false, now(), now())
	`, playerID, shipBoostID)

	// sync_shipbuilding: +1 ship slot, max level 1
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 1, false, now(), now())
	`, playerID, syncShipID)

	bonuses, err := GetPlayerTechBonuses(playerID)
	if err != nil {
		t.Fatalf("GetPlayerTechBonuses failed: %v", err)
	}

	// Check construction bonuses
	if bonuses.ConstructionSlots != 1 {
		t.Errorf("Expected ConstructionSlots=1, got %d", bonuses.ConstructionSlots)
	}
	if bonuses.BuildSpeed != 15.0 {
		t.Errorf("Expected BuildSpeed=15.0, got %f", bonuses.BuildSpeed)
	}
	if bonuses.ShipBuildSpeed != 12.0 {
		t.Errorf("Expected ShipBuildSpeed=12.0, got %f", bonuses.ShipBuildSpeed)
	}
	if bonuses.ShipProductionSlots != 1 {
		t.Errorf("Expected ShipProductionSlots=1, got %d", bonuses.ShipProductionSlots)
	}
}

func TestGetPlayerTechBonuses_CombatBonuses(t *testing.T) {
	setupTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000004"
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'testuser4', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupTestPlayer(t, playerID)

	// Get tech type IDs
	var ballisticsID, missilesID, directionalID, fighterID int
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'ballistics_base'").Scan(&ballisticsID)
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'missile_theory'").Scan(&missilesID)
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'optics_base'").Scan(&directionalID)
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'fighter_weapons_theory'").Scan(&fighterID)

	// ballistics_base: +5% per level, level 10 = +50%
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 10, false, now(), now())
	`, playerID, ballisticsID)

	// missile_theory: +4% per level, level 6 = +24%
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 6, false, now(), now())
	`, playerID, missilesID)

	// optics_base: +5% per level, level 8 = +40%
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 8, false, now(), now())
	`, playerID, directionalID)

	// fighter_weapons_theory: +3% per level, level 7 = +21%
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 7, false, now(), now())
	`, playerID, fighterID)

	bonuses, err := GetPlayerTechBonuses(playerID)
	if err != nil {
		t.Fatalf("GetPlayerTechBonuses failed: %v", err)
	}

	// Check combat bonuses
	if bonuses.BallisticDamage != 50.0 {
		t.Errorf("Expected BallisticDamage=50.0, got %f", bonuses.BallisticDamage)
	}
	if bonuses.MissileDamage != 24.0 {
		t.Errorf("Expected MissileDamage=24.0, got %f", bonuses.MissileDamage)
	}
	if bonuses.DirectionalDamage != 40.0 {
		t.Errorf("Expected DirectionalDamage=40.0, got %f", bonuses.DirectionalDamage)
	}
	if bonuses.FighterDamage != 21.0 {
		t.Errorf("Expected FighterDamage=21.0, got %f", bonuses.FighterDamage)
	}
}

func TestGetPlayerTechBonuses_WarehouseCapacity(t *testing.T) {
	setupTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000005"
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'testuser5', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupTestPlayer(t, playerID)

	// Get tech type ID
	var expandCapID int
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'expand_capacity'").Scan(&expandCapID)

	// expand_capacity: +50,000 per level, max 7, level 5 = +250,000
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 5, false, now(), now())
	`, playerID, expandCapID)

	bonuses, err := GetPlayerTechBonuses(playerID)
	if err != nil {
		t.Fatalf("GetPlayerTechBonuses failed: %v", err)
	}

	// Check warehouse capacity
	if bonuses.WarehouseCapacity != 250000 {
		t.Errorf("Expected WarehouseCapacity=250000, got %d", bonuses.WarehouseCapacity)
	}
}

func TestGetTechLevel_Exists(t *testing.T) {
	setupTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000101"
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'testlevel1', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupTestPlayer(t, playerID)

	// Get tech type ID for ballistics
	var ballisticsID int
	database.DB.QueryRow("SELECT id FROM tech_types WHERE name = 'ballistics_base'").Scan(&ballisticsID)

	// Insert level 7
	database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching, created_at, updated_at)
		VALUES ($1, $2, 7, false, now(), now())
	`, playerID, ballisticsID)

	level, err := GetTechLevel(playerID, "ballistics_base")
	if err != nil {
		t.Fatalf("GetTechLevel failed: %v", err)
	}

	if level != 7 {
		t.Errorf("Expected level=7, got %d", level)
	}
}

func TestGetTechLevel_NotResearched(t *testing.T) {
	setupTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000102"
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'testlevel2', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupTestPlayer(t, playerID)

	level, err := GetTechLevel(playerID, "ballistics_base")
	if err != nil {
		t.Fatalf("GetTechLevel failed: %v", err)
	}

	if level != 0 {
		t.Errorf("Expected level=0 for unresearched tech, got %d", level)
	}
}

func TestGetTechLevel_InvalidTech(t *testing.T) {
	setupTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000103"
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'testlevel3', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	defer cleanupTestPlayer(t, playerID)

	level, err := GetTechLevel(playerID, "nonexistent_tech")
	if err != nil {
		t.Fatalf("GetTechLevel should return 0 for invalid tech, got error: %v", err)
	}

	if level != 0 {
		t.Errorf("Expected level=0 for invalid tech, got %d", level)
	}
}

func TestApplyTechEffect_UnknownType(t *testing.T) {
	bonuses := &TechBonuses{}

	effect := &techEffectData{
		Type:     "unknown_effect_type",
		PerLevel: 10.0,
		Unit:     "percent",
	}

	// Should log but not panic
	applyTechEffect(bonuses, effect, 5, "test_tech")

	// All bonuses should remain zero
	if bonuses.MetalOutput != 0 {
		t.Errorf("Unknown effect should not modify bonuses")
	}
}

func TestTechBonuses_JSONSerialization(t *testing.T) {
	bonuses := &TechBonuses{
		MetalOutput:       10.5,
		ConstructionSlots: 2,
		BallisticDamage:   35.0,
		WarehouseCapacity: 150000,
	}

	// Test JSON marshaling
	data, err := json.Marshal(bonuses)
	if err != nil {
		t.Fatalf("Failed to marshal bonuses: %v", err)
	}

	// Test JSON unmarshaling
	var decoded TechBonuses
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal bonuses: %v", err)
	}

	if decoded.MetalOutput != bonuses.MetalOutput {
		t.Errorf("MetalOutput mismatch after JSON round-trip")
	}
	if decoded.ConstructionSlots != bonuses.ConstructionSlots {
		t.Errorf("ConstructionSlots mismatch after JSON round-trip")
	}
	if decoded.BallisticDamage != bonuses.BallisticDamage {
		t.Errorf("BallisticDamage mismatch after JSON round-trip")
	}
	if decoded.WarehouseCapacity != bonuses.WarehouseCapacity {
		t.Errorf("WarehouseCapacity mismatch after JSON round-trip")
	}
}
