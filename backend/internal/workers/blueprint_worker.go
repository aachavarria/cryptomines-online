package workers

import (
	"log"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/services"
)

// ApplyCompletedBlueprintResearch auto-completes any finished blueprint research timers.
// This function checks all active blueprint research and completes those that have finished.
func ApplyCompletedBlueprintResearch() {
	now := time.Now()

	// Find completed blueprint research before updating
	rows, err := database.DB.Query(`
		SELECT br.id, br.player_id, br.player_blueprint_id, br.target_level, b.name
		FROM blueprint_research br
		JOIN player_blueprints pb ON br.player_blueprint_id = pb.id
		JOIN blueprints b ON pb.blueprint_id = b.id
		WHERE br.is_researching = true AND br.research_finish_at <= $1
	`, now)
	if err != nil {
		log.Printf("Failed to query completed blueprint research: %v", err)
		return
	}

	type completedBlueprintResearch struct {
		researchID        string
		playerID          string
		playerBlueprintID string
		targetLevel       int
		blueprintName     string
	}
	var completed []completedBlueprintResearch

	for rows.Next() {
		var c completedBlueprintResearch
		if err := rows.Scan(&c.researchID, &c.playerID, &c.playerBlueprintID, &c.targetLevel, &c.blueprintName); err != nil {
			log.Printf("Failed to scan blueprint research row: %v", err)
			continue
		}
		completed = append(completed, c)
	}
	rows.Close()

	if len(completed) == 0 {
		return
	}

	log.Printf("Processing %d completed blueprint research operations", len(completed))

	// Process each completed research
	for _, c := range completed {
		// Begin transaction for this research
		tx, err := database.DB.Begin()
		if err != nil {
			log.Printf("Failed to begin transaction for blueprint research %s: %v", c.researchID, err)
			continue
		}

		// Update player_blueprints to new research level
		_, err = tx.Exec(`
			UPDATE player_blueprints
			SET research_level = $1
			WHERE id = $2
		`, c.targetLevel, c.playerBlueprintID)
		if err != nil {
			tx.Rollback()
			log.Printf("Failed to update player blueprint %s to level %d: %v", c.playerBlueprintID, c.targetLevel, err)
			continue
		}

		// Update blueprint_research to mark as complete
		_, err = tx.Exec(`
			UPDATE blueprint_research
			SET is_researching = false, research_finish_at = NULL
			WHERE id = $1
		`, c.researchID)
		if err != nil {
			tx.Rollback()
			log.Printf("Failed to complete blueprint research %s: %v", c.researchID, err)
			continue
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit blueprint research completion %s: %v", c.researchID, err)
			continue
		}

		log.Printf("Completed blueprint research: player=%s, blueprint=%s, level=%d",
			c.playerID, c.blueprintName, c.targetLevel)

		// Update quest progress (blueprint research quest)
		err = services.UpdateQuestProgress(c.playerID, "research_blueprint", c.blueprintName, 1)
		if err != nil {
			log.Printf("Failed to update quest progress for blueprint research: %v", err)
			// Don't fail the completion if quest update fails
		}
	}
}

// StartBlueprintWorker starts a background worker that periodically checks for completed blueprint research.
// It runs every 30 seconds to auto-complete finished research timers.
func StartBlueprintWorker() {
	ticker := time.NewTicker(30 * time.Second)
	log.Println("Blueprint research worker started (checking every 30 seconds)")

	// Run once immediately on startup
	go ApplyCompletedBlueprintResearch()

	// Then run on ticker
	go func() {
		for range ticker.C {
			ApplyCompletedBlueprintResearch()
		}
	}()
}
