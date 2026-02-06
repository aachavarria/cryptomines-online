package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
)

// ListPlanets handles GET /api/planets
func ListPlanets(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	rows, err := database.DB.Query(
		`SELECT id, player_id, name, position_x, position_y,
		        is_homeworld, is_rbp, rbp_level, controlling_corp_id,
		        protection_until, created_at, updated_at
		 FROM planets WHERE player_id = $1 ORDER BY created_at`, playerID,
	)
	if err != nil {
		log.Printf("Failed to list planets: %v", err)
		http.Error(w, `{"error":"failed to list planets"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	planets := []models.Planet{}
	for rows.Next() {
		var p models.Planet
		err := rows.Scan(
			&p.ID, &p.PlayerID, &p.Name, &p.PositionX, &p.PositionY,
			&p.IsHomeworld, &p.IsRBP, &p.RBPLevel, &p.ControllingCorpID,
			&p.ProtectionUntil, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan planet: %v", err)
			continue
		}
		planets = append(planets, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(planets)
}

// GetPlanet handles GET /api/planets/{id}
func GetPlanet(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	planetID := r.PathValue("id")

	var p models.Planet
	err := database.DB.QueryRow(
		`SELECT id, player_id, name, position_x, position_y,
		        is_homeworld, is_rbp, rbp_level, controlling_corp_id,
		        protection_until, created_at, updated_at
		 FROM planets WHERE id = $1 AND player_id = $2`, planetID, playerID,
	).Scan(
		&p.ID, &p.PlayerID, &p.Name, &p.PositionX, &p.PositionY,
		&p.IsHomeworld, &p.IsRBP, &p.RBPLevel, &p.ControllingCorpID,
		&p.ProtectionUntil, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to get planet: %v", err)
		http.Error(w, `{"error":"planet not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
