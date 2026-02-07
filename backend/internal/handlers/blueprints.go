package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
	"github.com/cryptomines-online/backend/internal/services"
)

// ListBlueprints handles GET /api/blueprints
// Lists all blueprints in the game (hull + module).
func ListBlueprints(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(
		`SELECT id, name, blueprint_type, hull_type_id, module_type_id,
		        source, research_level, description
		 FROM blueprints ORDER BY blueprint_type, name`,
	)
	if err != nil {
		log.Printf("Failed to list blueprints: %v", err)
		http.Error(w, `{"error":"failed to list blueprints"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	blueprints := []models.Blueprint{}
	for rows.Next() {
		var b models.Blueprint
		err := rows.Scan(
			&b.ID, &b.Name, &b.BlueprintType, &b.HullTypeID, &b.ModuleTypeID,
			&b.Source, &b.ResearchLevel, &b.Description,
		)
		if err != nil {
			log.Printf("Failed to scan blueprint: %v", err)
			continue
		}
		blueprints = append(blueprints, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blueprints)
}

// ListMyBlueprints handles GET /api/blueprints/mine
// Lists player's unlocked blueprints with blueprint info.
func ListMyBlueprints(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	rows, err := database.DB.Query(
		`SELECT pb.id, pb.player_id, pb.blueprint_id, pb.is_activated,
		        pb.research_level, pb.acquired_at,
		        b.name, b.blueprint_type, b.hull_type_id, b.module_type_id, b.source
		 FROM player_blueprints pb
		 JOIN blueprints b ON pb.blueprint_id = b.id
		 WHERE pb.player_id = $1
		 ORDER BY b.blueprint_type, b.name`, playerID,
	)
	if err != nil {
		log.Printf("Failed to list player blueprints: %v", err)
		http.Error(w, `{"error":"failed to list player blueprints"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	blueprints := []models.PlayerBlueprintWithInfo{}
	for rows.Next() {
		var pb models.PlayerBlueprintWithInfo
		err := rows.Scan(
			&pb.ID, &pb.PlayerID, &pb.BlueprintID, &pb.IsActivated,
			&pb.ResearchLevel, &pb.AcquiredAt,
			&pb.BlueprintName, &pb.BlueprintType, &pb.HullTypeID, &pb.ModuleTypeID, &pb.Source,
		)
		if err != nil {
			log.Printf("Failed to scan player blueprint: %v", err)
			continue
		}
		blueprints = append(blueprints, pb)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blueprints)
}

// ActivateBlueprint handles POST /api/blueprints/{id}/activate
// Activates an unactivated blueprint.
func ActivateBlueprint(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	blueprintID := r.PathValue("id")

	var pb models.PlayerBlueprint
	err := database.DB.QueryRow(
		`SELECT id, player_id, blueprint_id, is_activated, research_level, acquired_at
		 FROM player_blueprints
		 WHERE blueprint_id = $1 AND player_id = $2`, blueprintID, playerID,
	).Scan(&pb.ID, &pb.PlayerID, &pb.BlueprintID, &pb.IsActivated, &pb.ResearchLevel, &pb.AcquiredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"blueprint not found"}`, http.StatusNotFound)
			return
		}
		log.Printf("Failed to get player blueprint: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if pb.IsActivated {
		http.Error(w, `{"error":"blueprint already activated"}`, http.StatusConflict)
		return
	}

	err = database.DB.QueryRow(
		`UPDATE player_blueprints
		 SET is_activated = true
		 WHERE id = $1
		 RETURNING id, player_id, blueprint_id, is_activated, research_level, acquired_at`,
		pb.ID,
	).Scan(&pb.ID, &pb.PlayerID, &pb.BlueprintID, &pb.IsActivated, &pb.ResearchLevel, &pb.AcquiredAt)
	if err != nil {
		log.Printf("Failed to activate blueprint: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Get blueprint name for quest tracking
	var blueprintName string
	err = database.DB.QueryRow(`SELECT name FROM blueprints WHERE id = $1`, blueprintID).Scan(&blueprintName)
	if err == nil {
		// Update quest progress for using/activating a blueprint
		services.UpdateQuestProgress(playerID, "use_blueprint", blueprintName, 1)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pb)
}

// ResearchBlueprint handles POST /api/blueprints/{id}/research
// Starts researching a blueprint at the Weapon Research Center.
func ResearchBlueprint(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	blueprintID := r.PathValue("id")

	// Get the player blueprint
	var pb models.PlayerBlueprint
	err := database.DB.QueryRow(
		`SELECT id, player_id, blueprint_id, is_activated, research_level, acquired_at
		 FROM player_blueprints
		 WHERE blueprint_id = $1 AND player_id = $2`, blueprintID, playerID,
	).Scan(&pb.ID, &pb.PlayerID, &pb.BlueprintID, &pb.IsActivated, &pb.ResearchLevel, &pb.AcquiredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"blueprint not found"}`, http.StatusNotFound)
			return
		}
		log.Printf("Failed to get player blueprint: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if !pb.IsActivated {
		http.Error(w, `{"error":"blueprint must be activated first"}`, http.StatusConflict)
		return
	}

	if pb.ResearchLevel >= 3 {
		http.Error(w, `{"error":"blueprint already at max research level"}`, http.StatusConflict)
		return
	}

	// Check if already researching this blueprint
	var activeResearch int
	err = database.DB.QueryRow(
		`SELECT COUNT(*) FROM blueprint_research
		 WHERE player_id = $1 AND is_researching = true`, playerID,
	).Scan(&activeResearch)
	if err != nil {
		log.Printf("Failed to check active research: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Check Weapon Research Center level for research slot
	var wrcLevel int
	err = database.DB.QueryRow(
		`SELECT COALESCE(MAX(b.level), 0)
		 FROM buildings b
		 JOIN building_types bt ON b.building_type = bt.id
		 JOIN planets p ON b.planet_id = p.id
		 WHERE p.player_id = $1 AND bt.name = 'weapon_research_center'`, playerID,
	).Scan(&wrcLevel)
	if err != nil {
		log.Printf("Failed to check WRC level: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if wrcLevel < 1 {
		http.Error(w, `{"error":"weapon research center required"}`, http.StatusConflict)
		return
	}

	// 1 research slot per WRC (only 1 active at a time per WRC)
	if activeResearch >= 1 {
		http.Error(w, `{"error":"research slot in use"}`, http.StatusConflict)
		return
	}

	targetLevel := pb.ResearchLevel + 1

	// Research costs scale with level
	baseCost := int64(10000) * int64(targetLevel)
	metalCost := baseCost
	he3Cost := baseCost * 3 / 4
	goldCost := baseCost / 2

	// Deduct resources from homeworld
	var remaining struct{ Metal, He3, Gold int64 }
	err = database.DB.QueryRow(
		`UPDATE resources
		 SET metal = metal - $1, he3 = he3 - $2, gold = gold - $3, updated_at = now()
		 WHERE planet_id = (
		     SELECT p.id FROM planets p WHERE p.player_id = $4 AND p.is_homeworld = true LIMIT 1
		 ) AND metal >= $1 AND he3 >= $2 AND gold >= $3
		 RETURNING metal, he3, gold`,
		metalCost, he3Cost, goldCost, playerID,
	).Scan(&remaining.Metal, &remaining.He3, &remaining.Gold)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"insufficient resources"}`, http.StatusConflict)
			return
		}
		log.Printf("Failed to deduct resources: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Research time: 1 hour per level
	researchTimeSec := 3600 * targetLevel

	type researchResponse struct {
		ID               string `json:"id"`
		PlayerBlueprintID string `json:"player_blueprint_id"`
		TargetLevel      int    `json:"target_level"`
		IsResearching    bool   `json:"is_researching"`
		ResearchFinishAt string `json:"research_finish_at"`
		MetalCost        int64  `json:"metal_cost"`
		He3Cost          int64  `json:"he3_cost"`
		GoldCost         int64  `json:"gold_cost"`
	}

	var resp researchResponse
	err = database.DB.QueryRow(
		`INSERT INTO blueprint_research (player_id, player_blueprint_id, target_level, is_researching, research_finish_at, metal_cost, he3_cost, gold_cost)
		 VALUES ($1, $2, $3, true, now() + ($4 || ' seconds')::interval, $5, $6, $7)
		 RETURNING id, player_blueprint_id, target_level, is_researching, research_finish_at::text, metal_cost, he3_cost, gold_cost`,
		playerID, pb.ID, targetLevel, researchTimeSec, metalCost, he3Cost, goldCost,
	).Scan(&resp.ID, &resp.PlayerBlueprintID, &resp.TargetLevel, &resp.IsResearching, &resp.ResearchFinishAt, &resp.MetalCost, &resp.He3Cost, &resp.GoldCost)
	if err != nil {
		log.Printf("Failed to start research: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
