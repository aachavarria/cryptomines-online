package workers

import (
	"testing"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/google/uuid"
)

func setupBlueprintWorkerTest(t *testing.T) (playerID string, cleanup func()) {
	if err := database.InitSupabase(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	playerID = uuid.New().String()
	anonID := "anon-" + uuid.New().String()
	username := "test-user-" + uuid.New().String()

	// Create test player
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username)
		VALUES ($1, $2, $3)
	`, playerID, anonID, username)
	if err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}

	// Create test planet (minimal fields required)
	planetID := uuid.New().String()
	_, err = database.DB.Exec(`
		INSERT INTO planets (id, player_id, name, is_homeworld, position_x, position_y)
		VALUES ($1, $2, 'Test Planet', true, $3, $4)
	`, planetID, playerID, 1000+time.Now().Unix()%1000, 1000+time.Now().Unix()%1000)
	if err != nil {
		t.Fatalf("Failed to create test planet: %v", err)
	}

	cleanup = func() {
		database.DB.Exec("DELETE FROM blueprint_research WHERE player_id = $1", playerID)
		database.DB.Exec("DELETE FROM player_blueprints WHERE player_id = $1", playerID)
		database.DB.Exec("DELETE FROM planets WHERE player_id = $1", playerID)
		database.DB.Exec("DELETE FROM players WHERE id = $1", playerID)
	}

	return playerID, cleanup
}

func TestApplyCompletedBlueprintResearch(t *testing.T) {
	playerID, cleanup := setupBlueprintWorkerTest(t)
	defer cleanup()

	// Get a blueprint to test with
	var blueprintID int
	err := database.DB.QueryRow(`
		SELECT id FROM blueprints WHERE blueprint_type = 'hull' LIMIT 1
	`).Scan(&blueprintID)
	if err != nil {
		t.Fatalf("Failed to find test blueprint: %v", err)
	}

	// Create player blueprint (activated, level 1)
	var playerBlueprintID string
	err = database.DB.QueryRow(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 1)
		RETURNING id
	`, playerID, blueprintID).Scan(&playerBlueprintID)
	if err != nil {
		t.Fatalf("Failed to create player blueprint: %v", err)
	}

	// Create blueprint research that finished 1 minute ago
	pastTime := time.Now().Add(-1 * time.Minute)
	var researchID string
	err = database.DB.QueryRow(`
		INSERT INTO blueprint_research (player_id, player_blueprint_id, target_level, is_researching, research_finish_at, created_at)
		VALUES ($1, $2, 2, true, $3, now())
		RETURNING id
	`, playerID, playerBlueprintID, pastTime).Scan(&researchID)
	if err != nil {
		t.Fatalf("Failed to create blueprint research: %v", err)
	}

	// Run the worker function
	ApplyCompletedBlueprintResearch()

	// Verify player_blueprints was updated to level 2
	var researchLevel int
	err = database.DB.QueryRow(`
		SELECT research_level FROM player_blueprints WHERE id = $1
	`, playerBlueprintID).Scan(&researchLevel)
	if err != nil {
		t.Fatalf("Failed to query player blueprint: %v", err)
	}
	if researchLevel != 2 {
		t.Errorf("Expected research_level = 2, got %d", researchLevel)
	}

	// Verify blueprint_research was marked as complete
	var isResearching bool
	var finishAt *time.Time
	err = database.DB.QueryRow(`
		SELECT is_researching, research_finish_at FROM blueprint_research WHERE id = $1
	`, researchID).Scan(&isResearching, &finishAt)
	if err != nil {
		t.Fatalf("Failed to query blueprint research: %v", err)
	}
	if isResearching {
		t.Errorf("Expected is_researching = false, got true")
	}
	if finishAt != nil {
		t.Errorf("Expected research_finish_at = NULL, got %v", finishAt)
	}
}

func TestApplyCompletedBlueprintResearch_NoCompleted(t *testing.T) {
	playerID, cleanup := setupBlueprintWorkerTest(t)
	defer cleanup()

	// Get a blueprint to test with
	var blueprintID int
	err := database.DB.QueryRow(`
		SELECT id FROM blueprints WHERE blueprint_type = 'hull' LIMIT 1
	`).Scan(&blueprintID)
	if err != nil {
		t.Fatalf("Failed to find test blueprint: %v", err)
	}

	// Create player blueprint (activated, level 1)
	var playerBlueprintID string
	err = database.DB.QueryRow(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 1)
		RETURNING id
	`, playerID, blueprintID).Scan(&playerBlueprintID)
	if err != nil {
		t.Fatalf("Failed to create player blueprint: %v", err)
	}

	// Create blueprint research that finishes in the future
	futureTime := time.Now().Add(1 * time.Hour)
	_, err = database.DB.Exec(`
		INSERT INTO blueprint_research (player_id, player_blueprint_id, target_level, is_researching, research_finish_at, created_at)
		VALUES ($1, $2, 2, true, $3, now())
	`, playerID, playerBlueprintID, futureTime)
	if err != nil {
		t.Fatalf("Failed to create blueprint research: %v", err)
	}

	// Run the worker function (should do nothing)
	ApplyCompletedBlueprintResearch()

	// Verify player_blueprints is still level 1
	var researchLevel int
	err = database.DB.QueryRow(`
		SELECT research_level FROM player_blueprints WHERE id = $1
	`, playerBlueprintID).Scan(&researchLevel)
	if err != nil {
		t.Fatalf("Failed to query player blueprint: %v", err)
	}
	if researchLevel != 1 {
		t.Errorf("Expected research_level = 1 (unchanged), got %d", researchLevel)
	}
}

func TestApplyCompletedBlueprintResearch_MultipleCompletions(t *testing.T) {
	playerID, cleanup := setupBlueprintWorkerTest(t)
	defer cleanup()

	// Get two blueprints to test with
	var blueprintID1, blueprintID2 int
	rows, err := database.DB.Query(`
		SELECT id FROM blueprints WHERE blueprint_type = 'hull' LIMIT 2
	`)
	if err != nil {
		t.Fatalf("Failed to find test blueprints: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&blueprintID1)
	}
	if rows.Next() {
		rows.Scan(&blueprintID2)
	}

	// Create player blueprints (activated, level 1)
	var playerBlueprintID1, playerBlueprintID2 string
	err = database.DB.QueryRow(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 1)
		RETURNING id
	`, playerID, blueprintID1).Scan(&playerBlueprintID1)
	if err != nil {
		t.Fatalf("Failed to create player blueprint 1: %v", err)
	}

	err = database.DB.QueryRow(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 1)
		RETURNING id
	`, playerID, blueprintID2).Scan(&playerBlueprintID2)
	if err != nil {
		t.Fatalf("Failed to create player blueprint 2: %v", err)
	}

	// Create two blueprint researches that both finished
	pastTime := time.Now().Add(-1 * time.Minute)
	_, err = database.DB.Exec(`
		INSERT INTO blueprint_research (player_id, player_blueprint_id, target_level, is_researching, research_finish_at, created_at)
		VALUES ($1, $2, 2, true, $3, now())
	`, playerID, playerBlueprintID1, pastTime)
	if err != nil {
		t.Fatalf("Failed to create blueprint research 1: %v", err)
	}

	_, err = database.DB.Exec(`
		INSERT INTO blueprint_research (player_id, player_blueprint_id, target_level, is_researching, research_finish_at, created_at)
		VALUES ($1, $2, 2, true, $3, now())
	`, playerID, playerBlueprintID2, pastTime)
	if err != nil {
		t.Fatalf("Failed to create blueprint research 2: %v", err)
	}

	// Run the worker function
	ApplyCompletedBlueprintResearch()

	// Verify both player_blueprints were updated to level 2
	var researchLevel1, researchLevel2 int
	err = database.DB.QueryRow(`
		SELECT research_level FROM player_blueprints WHERE id = $1
	`, playerBlueprintID1).Scan(&researchLevel1)
	if err != nil {
		t.Fatalf("Failed to query player blueprint 1: %v", err)
	}
	if researchLevel1 != 2 {
		t.Errorf("Blueprint 1: Expected research_level = 2, got %d", researchLevel1)
	}

	err = database.DB.QueryRow(`
		SELECT research_level FROM player_blueprints WHERE id = $1
	`, playerBlueprintID2).Scan(&researchLevel2)
	if err != nil {
		t.Fatalf("Failed to query player blueprint 2: %v", err)
	}
	if researchLevel2 != 2 {
		t.Errorf("Blueprint 2: Expected research_level = 2, got %d", researchLevel2)
	}

	// Verify both blueprint_research records were marked as complete
	var count int
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM blueprint_research
		WHERE player_id = $1 AND is_researching = false AND research_finish_at IS NULL
	`, playerID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count completed research: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected 2 completed research records, got %d", count)
	}
}

func TestApplyCompletedBlueprintResearch_Level2to3(t *testing.T) {
	playerID, cleanup := setupBlueprintWorkerTest(t)
	defer cleanup()

	// Get a blueprint to test with
	var blueprintID int
	err := database.DB.QueryRow(`
		SELECT id FROM blueprints WHERE blueprint_type = 'hull' LIMIT 1
	`).Scan(&blueprintID)
	if err != nil {
		t.Fatalf("Failed to find test blueprint: %v", err)
	}

	// Create player blueprint (activated, level 2)
	var playerBlueprintID string
	err = database.DB.QueryRow(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 2)
		RETURNING id
	`, playerID, blueprintID).Scan(&playerBlueprintID)
	if err != nil {
		t.Fatalf("Failed to create player blueprint: %v", err)
	}

	// Create blueprint research to upgrade from level 2 to level 3
	pastTime := time.Now().Add(-1 * time.Minute)
	var researchID string
	err = database.DB.QueryRow(`
		INSERT INTO blueprint_research (player_id, player_blueprint_id, target_level, is_researching, research_finish_at, created_at)
		VALUES ($1, $2, 3, true, $3, now())
		RETURNING id
	`, playerID, playerBlueprintID, pastTime).Scan(&researchID)
	if err != nil {
		t.Fatalf("Failed to create blueprint research: %v", err)
	}

	// Run the worker function
	ApplyCompletedBlueprintResearch()

	// Verify player_blueprints was updated to level 3
	var researchLevel int
	err = database.DB.QueryRow(`
		SELECT research_level FROM player_blueprints WHERE id = $1
	`, playerBlueprintID).Scan(&researchLevel)
	if err != nil {
		t.Fatalf("Failed to query player blueprint: %v", err)
	}
	if researchLevel != 3 {
		t.Errorf("Expected research_level = 3, got %d", researchLevel)
	}

	// Verify blueprint_research was marked as complete
	var isResearching bool
	err = database.DB.QueryRow(`
		SELECT is_researching FROM blueprint_research WHERE id = $1
	`, researchID).Scan(&isResearching)
	if err != nil {
		t.Fatalf("Failed to query blueprint research: %v", err)
	}
	if isResearching {
		t.Errorf("Expected is_researching = false, got true")
	}
}
