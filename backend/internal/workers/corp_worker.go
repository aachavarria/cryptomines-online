package workers

import (
	"log"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/services"
)

// UpdateCorpDailyReset resets daily contribution for all corp members
// This should run once per day, but we check if the last reset was > 24h ago
func UpdateCorpDailyReset() {
	now := time.Now()

	// Reset daily_contribution to 0 for all members
	// We don't track last reset time, so we just reset everyone
	result, err := database.DB.Exec(`
		UPDATE corp_members
		SET daily_contribution = 0
		WHERE daily_contribution > 0
	`)
	if err != nil {
		log.Printf("Failed to reset daily contributions: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("Reset daily contributions for %d corp members at %s", rowsAffected, now.Format(time.RFC3339))
	}
}

// UpdateCorpLevels recalculates corp levels based on their current wealth
func UpdateCorpLevels() {
	rows, err := database.DB.Query(`
		SELECT id, wealth, level
		FROM corps
		ORDER BY id
	`)
	if err != nil {
		log.Printf("Failed to query corps for level update: %v", err)
		return
	}
	defer rows.Close()

	type corpUpdate struct {
		id       string
		wealth   int64
		oldLevel int
		newLevel int
	}

	var updates []corpUpdate
	for rows.Next() {
		var u corpUpdate
		if err := rows.Scan(&u.id, &u.wealth, &u.oldLevel); err != nil {
			log.Printf("Failed to scan corp row: %v", err)
			continue
		}

		// Calculate new level
		u.newLevel = services.CalculateCorpLevel(u.wealth)

		// Only update if level changed
		if u.newLevel != u.oldLevel {
			updates = append(updates, u)
		}
	}

	if len(updates) == 0 {
		return
	}

	log.Printf("Updating %d corp levels", len(updates))

	// Update each corp
	successCount := 0
	for _, u := range updates {
		_, err := database.DB.Exec(`
			UPDATE corps
			SET level = $1, updated_at = $2
			WHERE id = $3
		`, u.newLevel, time.Now(), u.id)

		if err != nil {
			log.Printf("Failed to update corp %s level from %d to %d: %v", u.id, u.oldLevel, u.newLevel, err)
			continue
		}

		log.Printf("Updated corp %s level from %d to %d (wealth: %d)", u.id, u.oldLevel, u.newLevel, u.wealth)
		successCount++
	}

	log.Printf("Successfully updated %d/%d corp levels", successCount, len(updates))
}

// RemoveExpiredRBPProtection removes protection from RBP planets where protection_until has passed
func RemoveExpiredRBPProtection() {
	now := time.Now()

	result, err := database.DB.Exec(`
		UPDATE planets
		SET protection_until = NULL
		WHERE is_rbp = true
		  AND protection_until IS NOT NULL
		  AND protection_until < $1
	`, now)
	if err != nil {
		log.Printf("Failed to remove expired RBP protection: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("Removed protection from %d RBP planets at %s", rowsAffected, now.Format(time.RFC3339))
	}
}

// StartCorpWorker starts a background worker that runs corp-related maintenance tasks.
// It runs every 1 hour to:
// 1. Reset daily_contribution (if needed)
// 2. Update corp levels based on wealth
// 3. Remove expired RBP protection
func StartCorpWorker() {
	ticker := time.NewTicker(1 * time.Hour)
	log.Println("Corp worker started (running every 1 hour)")

	// Run once immediately on startup
	go func() {
		UpdateCorpDailyReset()
		UpdateCorpLevels()
		RemoveExpiredRBPProtection()
	}()

	// Then run on ticker
	go func() {
		for range ticker.C {
			UpdateCorpDailyReset()
			UpdateCorpLevels()
			RemoveExpiredRBPProtection()
		}
	}()
}
