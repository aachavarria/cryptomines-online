package services

import (
	"log"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
)

// UpdateQuestProgress checks and updates quest progress for a player
// based on the action they performed.
// It finds all active quests matching the requirement type and target,
// increments their progress, and marks them as completed if the requirement is met.
func UpdateQuestProgress(playerID, requirementType, requirementTarget string, incrementBy int) error {
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction for quest progress: %v", err)
		return err
	}
	defer tx.Rollback()

	// Find all active quests that match this action
	// If requirement_target is empty in the quest, match any target for that type
	rows, err := tx.Query(`
		SELECT pq.id, pq.progress_value, qt.requirement_value
		FROM player_quests pq
		JOIN quest_types qt ON pq.quest_type_id = qt.id
		WHERE pq.player_id = $1
		  AND pq.status = 'available'
		  AND qt.requirement_type = $2
		  AND (qt.requirement_target = $3 OR qt.requirement_target = '')
		  AND qt.is_active = true
		FOR UPDATE OF pq
	`, playerID, requirementType, requirementTarget)
	if err != nil {
		log.Printf("Failed to query quests for progress update: %v", err)
		return err
	}
	defer rows.Close()

	type questUpdate struct {
		id              string
		newProgress     int
		requirementValue int
	}

	var updates []questUpdate
	for rows.Next() {
		var id string
		var currentProgress, requirementValue int
		if err := rows.Scan(&id, &currentProgress, &requirementValue); err != nil {
			log.Printf("Failed to scan quest row: %v", err)
			continue
		}

		newProgress := currentProgress + incrementBy
		updates = append(updates, questUpdate{
			id:              id,
			newProgress:     newProgress,
			requirementValue: requirementValue,
		})
	}
	rows.Close()

	// Update each quest
	now := time.Now()
	for _, update := range updates {
		if update.newProgress >= update.requirementValue {
			// Mark as completed
			_, err = tx.Exec(`
				UPDATE player_quests
				SET progress_value = $1, status = 'completed', completed_at = $2, updated_at = $2
				WHERE id = $3
			`, update.newProgress, now, update.id)
			if err != nil {
				log.Printf("Failed to complete quest %s: %v", update.id, err)
				return err
			}
			log.Printf("Quest %s completed for player %s", update.id, playerID)
		} else {
			// Just update progress
			_, err = tx.Exec(`
				UPDATE player_quests
				SET progress_value = $1, updated_at = $2
				WHERE id = $3
			`, update.newProgress, now, update.id)
			if err != nil {
				log.Printf("Failed to update quest progress %s: %v", update.id, err)
				return err
			}
			log.Printf("Quest %s progress updated to %d/%d for player %s", 
				update.id, update.newProgress, update.requirementValue, playerID)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit quest progress update: %v", err)
		return err
	}

	return nil
}

// CheckAndCompleteQuests marks specific quests as completed if their progress meets the requirement.
// This is useful for batch checking after complex operations.
func CheckAndCompleteQuests(playerID string, questTypeIDs []int) error {
	if len(questTypeIDs) == 0 {
		return nil
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction for quest completion check: %v", err)
		return err
	}
	defer tx.Rollback()

	// Build the query with IN clause
	query := `
		UPDATE player_quests pq
		SET status = 'completed', completed_at = $1, updated_at = $1
		FROM quest_types qt
		WHERE pq.quest_type_id = qt.id
		  AND pq.player_id = $2
		  AND pq.status = 'available'
		  AND pq.progress_value >= qt.requirement_value
		  AND pq.quest_type_id = ANY($3)
	`

	now := time.Now()
	result, err := tx.Exec(query, now, playerID, questTypeIDs)
	if err != nil {
		log.Printf("Failed to check and complete quests: %v", err)
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected > 0 {
		log.Printf("Completed %d quests for player %s", rowsAffected, playerID)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit quest completion: %v", err)
		return err
	}

	return nil
}

// SyncBuildingQuests checks all existing buildings and updates quest progress
// for any "build_building" quests that should already be completed
func SyncBuildingQuests(playerID string) error {
	// Single query: find all quests that should be completed based on existing buildings
	// This avoids nested queries inside a transaction (Go sql.Tx uses single connection)
	rows, err := database.DB.Query(`
		SELECT pq.id, qt.requirement_value
		FROM player_quests pq
		JOIN quest_types qt ON pq.quest_type_id = qt.id
		WHERE pq.player_id = $1
		  AND pq.status = 'available'
		  AND qt.requirement_type = 'build_building'
		  AND qt.is_active = true
		  AND EXISTS (
			SELECT 1
			FROM buildings b
			JOIN building_types bt ON b.building_type = bt.id
			JOIN planets p ON b.planet_id = p.id
			WHERE p.player_id = $1
			  AND bt.name = qt.requirement_target
			  AND b.level >= qt.requirement_value
		  )
	`, playerID)
	if err != nil {
		log.Printf("Failed to query syncable quests: %v", err)
		return err
	}

	type questUpdate struct {
		id    string
		value int
	}
	var updates []questUpdate
	for rows.Next() {
		var u questUpdate
		if err := rows.Scan(&u.id, &u.value); err != nil {
			continue
		}
		updates = append(updates, u)
	}
	rows.Close()

	// Apply updates
	for _, u := range updates {
		_, err := database.DB.Exec(`
			UPDATE player_quests
			SET progress_value = $1, status = 'completed', completed_at = now(), updated_at = now()
			WHERE id = $2
		`, u.value, u.id)
		if err != nil {
			log.Printf("Failed to auto-complete quest %s: %v", u.id, err)
		} else {
			log.Printf("Auto-completed quest %s for player %s", u.id, playerID)
		}
	}

	return nil
}
