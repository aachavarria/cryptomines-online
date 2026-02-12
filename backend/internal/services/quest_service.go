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
	rows, err := tx.Query(`
		SELECT pq.id, pq.progress_value, qt.requirement_value
		FROM player_quests pq
		JOIN quest_types qt ON pq.quest_type_id = qt.id
		WHERE pq.player_id = $1
		  AND pq.status = 'available'
		  AND qt.requirement_type = $2
		  AND qt.requirement_target = $3
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
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction for quest sync: %v", err)
		return err
	}
	defer tx.Rollback()

	// Get all buildings the player has at level >= 1
	rows, err := tx.Query(`
		SELECT building_type, level
		FROM player_buildings
		WHERE player_id = $1 AND level >= 1
	`, playerID)
	if err != nil {
		log.Printf("Failed to query player buildings: %v", err)
		return err
	}
	defer rows.Close()

	type buildingInfo struct {
		buildingType string
		level        int
	}
	var buildings []buildingInfo
	for rows.Next() {
		var b buildingInfo
		if err := rows.Scan(&b.buildingType, &b.level); err != nil {
			continue
		}
		buildings = append(buildings, b)
	}
	rows.Close()

	// For each building, update any "build_building" quests
	for _, building := range buildings {
		// Find active quests requiring this building at level 1
		questRows, err := tx.Query(`
			SELECT pq.id, pq.progress_value, qt.requirement_value
			FROM player_quests pq
			JOIN quest_types qt ON pq.quest_type_id = qt.id
			WHERE pq.player_id = $1
			  AND pq.status = 'available'
			  AND qt.requirement_type = 'build_building'
			  AND qt.requirement_target = $2
			  AND qt.requirement_value <= $3
			  AND qt.is_active = true
		`, playerID, building.buildingType, building.level)
		if err != nil {
			continue
		}

		for questRows.Next() {
			var questID string
			var currentProgress, requirementValue int
			if err := questRows.Scan(&questID, &currentProgress, &requirementValue); err != nil {
				continue
			}

			// If progress is less than requirement, complete it
			if currentProgress < requirementValue {
				_, err = tx.Exec(`
					UPDATE player_quests
					SET progress_value = $1, status = 'completed', completed_at = now(), updated_at = now()
					WHERE id = $2
				`, requirementValue, questID)
				if err == nil {
					log.Printf("Auto-completed quest %s for player %s (building %s already exists)",
						questID, playerID, building.buildingType)
				}
			}
		}
		questRows.Close()
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit quest sync: %v", err)
		return err
	}

	return nil
}
