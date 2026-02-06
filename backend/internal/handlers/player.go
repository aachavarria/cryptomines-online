package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
)

// PlayerMe handles GET /api/player/me
func PlayerMe(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	var p models.Player
	err := database.DB.QueryRow(
		`SELECT id, anonymous_id, username, level, experience,
		        mall_points, vouchers, honor_points, champion_points,
		        badges, corsairs_gold, tutorial_step, is_online,
		        created_at, last_login, updated_at
		 FROM players WHERE id = $1`, playerID,
	).Scan(
		&p.ID, &p.AnonymousID, &p.Username, &p.Level, &p.Experience,
		&p.MallPoints, &p.Vouchers, &p.HonorPoints, &p.ChampionPoints,
		&p.Badges, &p.CorsairsGold, &p.TutorialStep, &p.IsOnline,
		&p.CreatedAt, &p.LastLogin, &p.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to get player: %v", err)
		http.Error(w, `{"error":"player not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
