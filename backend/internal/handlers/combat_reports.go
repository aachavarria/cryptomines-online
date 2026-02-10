package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
)

// ListCombatReports handles GET /api/combat-reports
func ListCombatReports(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	rows, err := database.DB.Query(
		`SELECT id, attacker_id, defender_id, combat_type, result, total_rounds,
		        loot_json, he3_consumed, created_at
		 FROM combat_reports
		 WHERE attacker_id = $1 OR defender_id = $1
		 ORDER BY created_at DESC
		 LIMIT 100`, playerID,
	)
	if err != nil {
		log.Printf("Failed to list combat reports: %v", err)
		http.Error(w, `{"error":"failed to list combat reports"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	reports := []models.CombatReport{}
	for rows.Next() {
		var report models.CombatReport
		err := rows.Scan(
			&report.ID, &report.AttackerID, &report.DefenderID, &report.CombatType,
			&report.Result, &report.TotalRounds, &report.LootJSON, &report.He3Consumed,
			&report.CreatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan combat report: %v", err)
			continue
		}
		reports = append(reports, report)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

// GetCombatReport handles GET /api/combat-reports/{id}
func GetCombatReport(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	reportID := r.PathValue("id")

	var report models.CombatReport
	err := database.DB.QueryRow(
		`SELECT id, attacker_id, defender_id, combat_type, result, total_rounds,
		        loot_json, rounds_json, he3_consumed, created_at
		 FROM combat_reports
		 WHERE id = $1 AND (attacker_id = $2 OR defender_id = $2)`, reportID, playerID,
	).Scan(
		&report.ID, &report.AttackerID, &report.DefenderID, &report.CombatType,
		&report.Result, &report.TotalRounds, &report.LootJSON, &report.RoundsJSON,
		&report.He3Consumed, &report.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"combat report not found"}`, http.StatusNotFound)
			return
		}
		log.Printf("Failed to get combat report: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
