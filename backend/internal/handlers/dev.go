package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/services"
)

// devCheck validates that X-Dev-Mode header is set. Returns false if blocked.
func devCheck(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("X-Dev-Mode") != "true" {
		http.Error(w, `{"error":"dev mode not enabled"}`, http.StatusForbidden)
		return false
	}
	return true
}

// DevResetPlayer handles POST /api/dev/reset
// Resets player to initial state: clears buildings, techs, quests, commanders,
// ships, fleets, blueprints, recycling, combat reports, inventory, buffs, chat.
// Resets resources to 5000/5000/10000.
func DevResetPlayer(w http.ResponseWriter, r *http.Request) {
	if !devCheck(w, r) {
		return
	}
	playerID := middleware.GetPlayerID(r)

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("DevResetPlayer: tx begin failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Get the player's planet IDs
	rows, err := tx.Query(`SELECT id FROM planets WHERE player_id = $1`, playerID)
	if err != nil {
		log.Printf("DevResetPlayer: query planets failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	var planetIDs []string
	for rows.Next() {
		var pid string
		if err := rows.Scan(&pid); err != nil {
			rows.Close()
			log.Printf("DevResetPlayer: scan planet failed: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
		planetIDs = append(planetIDs, pid)
	}
	rows.Close()

	// Delete buildings for all planets, then re-insert defaults
	type initBuilding struct {
		name    string
		gridCol int
		gridRow int
	}
	defaultBuildings := []initBuilding{
		{"civic_center", 5, 5},
		{"metal_collector", 2, 2},
		{"he3_extractor", 9, 2},
		{"residential_area", 2, 9},
		{"resource_warehouse", 9, 9},
		{"space_station", 5, 1},
	}
	for _, pid := range planetIDs {
		if _, err := tx.Exec(`DELETE FROM buildings WHERE planet_id = $1`, pid); err != nil {
			log.Printf("DevResetPlayer: delete buildings failed: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
		// Re-insert default buildings at level 1
		for _, b := range defaultBuildings {
			var btID int
			if err := tx.QueryRow(`SELECT id FROM building_types WHERE name = $1`, b.name).Scan(&btID); err != nil {
				log.Printf("DevResetPlayer: lookup building type %s failed: %v", b.name, err)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
				return
			}
			if _, err := tx.Exec(
				`INSERT INTO buildings (planet_id, building_type, level, grid_col, grid_row) VALUES ($1, $2, 1, $3, $4)`,
				pid, btID, b.gridCol, b.gridRow,
			); err != nil {
				log.Printf("DevResetPlayer: insert default building %s failed: %v", b.name, err)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
				return
			}
		}
	}

	// Reset resources for all planets (with correct Lv1 production rates)
	for _, pid := range planetIDs {
		if _, err := tx.Exec(`
			UPDATE resources SET
				metal = 5000, he3 = 5000, gold = 10000,
				metal_per_hour = 1080, he3_per_hour = 1180, gold_per_hour = 1300,
				storage_capacity = 10000,
				warehouse_metal = 0, warehouse_he3 = 0, warehouse_gold = 0,
				last_collected_at = now(), last_warehouse_update = now(), updated_at = now()
			WHERE planet_id = $1`, pid); err != nil {
			log.Printf("DevResetPlayer: reset resources failed: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
	}

	// Delete technologies
	if _, err := tx.Exec(`DELETE FROM technologies WHERE player_id = $1`, playerID); err != nil {
		log.Printf("DevResetPlayer: delete technologies failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Delete blueprint research (references player_blueprints, must go first)
	if _, err := tx.Exec(`DELETE FROM blueprint_research WHERE player_id = $1`, playerID); err != nil {
		log.Printf("DevResetPlayer: delete blueprint_research failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Delete player blueprints
	if _, err := tx.Exec(`DELETE FROM player_blueprints WHERE player_id = $1`, playerID); err != nil {
		log.Printf("DevResetPlayer: delete player_blueprints failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Delete ship instances (references ship_designs and fleet_stacks)
	if _, err := tx.Exec(`DELETE FROM ship_instances WHERE player_id = $1`, playerID); err != nil {
		log.Printf("DevResetPlayer: delete ship_instances failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Delete fleet stacks (must delete before fleets due to FK)
	if _, err := tx.Exec(`
		DELETE FROM fleet_stacks WHERE fleet_id IN (
			SELECT id FROM fleets WHERE player_id = $1
		)`, playerID); err != nil {
		log.Printf("DevResetPlayer: delete fleet_stacks failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Delete fleets
	if _, err := tx.Exec(`DELETE FROM fleets WHERE player_id = $1`, playerID); err != nil {
		log.Printf("DevResetPlayer: delete fleets failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Delete spacedock repairs (references ship_designs)
	if _, err := tx.Exec(`DELETE FROM spacedock_repairs WHERE player_id = $1`, playerID); err != nil {
		log.Printf("DevResetPlayer: delete spacedock_repairs failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Delete ships (the old quantity-based table, references ship_designs)
	if _, err := tx.Exec(`DELETE FROM ships WHERE player_id = $1`, playerID); err != nil {
		log.Printf("DevResetPlayer: delete ships failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Delete ship designs
	if _, err := tx.Exec(`DELETE FROM ship_designs WHERE player_id = $1`, playerID); err != nil {
		log.Printf("DevResetPlayer: delete ship_designs failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Delete commanders
	if _, err := tx.Exec(`DELETE FROM commanders WHERE player_id = $1`, playerID); err != nil {
		log.Printf("DevResetPlayer: delete commanders failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Delete recycling jobs
	if _, err := tx.Exec(`DELETE FROM recycling_jobs WHERE player_id = $1`, playerID); err != nil {
		log.Printf("DevResetPlayer: delete recycling_jobs failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Delete combat reports (player can be attacker or defender)
	if _, err := tx.Exec(`DELETE FROM combat_reports WHERE attacker_id = $1 OR defender_id = $1`, playerID); err != nil {
		log.Printf("DevResetPlayer: delete combat_reports failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Reset quest progress
	if _, err := tx.Exec(`DELETE FROM player_quests WHERE player_id = $1`, playerID); err != nil {
		log.Printf("DevResetPlayer: delete player_quests failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if _, err := tx.Exec(`DELETE FROM daily_quest_progress WHERE player_id = $1`, playerID); err != nil {
		log.Printf("DevResetPlayer: delete daily_quest_progress failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Delete inventory
	if _, err := tx.Exec(`DELETE FROM player_inventory WHERE player_id = $1`, playerID); err != nil {
		log.Printf("DevResetPlayer: delete player_inventory failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Delete active buffs
	if _, err := tx.Exec(`DELETE FROM active_buffs WHERE player_id = $1`, playerID); err != nil {
		log.Printf("DevResetPlayer: delete active_buffs failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Delete chat messages
	if _, err := tx.Exec(`DELETE FROM chat_messages WHERE player_id = $1`, playerID); err != nil {
		log.Printf("DevResetPlayer: delete chat_messages failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("DevResetPlayer: commit failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// DevGiveResources handles POST /api/dev/give-resources
// Body: {"metal": int, "he3": int, "gold": int}
// Adds the specified amounts to the player's homeworld resources.
func DevGiveResources(w http.ResponseWriter, r *http.Request) {
	if !devCheck(w, r) {
		return
	}
	playerID := middleware.GetPlayerID(r)

	var body struct {
		Metal int64 `json:"metal"`
		He3   int64 `json:"he3"`
		Gold  int64 `json:"gold"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(`
		UPDATE resources SET
			metal = metal + $1,
			he3 = he3 + $2,
			gold = gold + $3,
			updated_at = now()
		WHERE planet_id = (
			SELECT id FROM planets WHERE player_id = $4 AND is_homeworld = true LIMIT 1
		)`, body.Metal, body.He3, body.Gold, playerID)
	if err != nil {
		log.Printf("DevGiveResources: update failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error":"no homeworld found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// DevGiveItem handles POST /api/dev/give-item
// Body: {"item_key": string, "quantity": int}
// Upserts the item into the player's inventory.
func DevGiveItem(w http.ResponseWriter, r *http.Request) {
	if !devCheck(w, r) {
		return
	}
	playerID := middleware.GetPlayerID(r)

	var body struct {
		ItemKey  string `json:"item_key"`
		Quantity int    `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if body.ItemKey == "" || body.Quantity <= 0 {
		http.Error(w, `{"error":"item_key and positive quantity required"}`, http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec(`
		INSERT INTO player_inventory (player_id, item_key, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (player_id, item_key) DO UPDATE SET
			quantity = player_inventory.quantity + $3,
			updated_at = now()
	`, playerID, body.ItemKey, body.Quantity)
	if err != nil {
		log.Printf("DevGiveItem: upsert failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// DevCompleteConstructions handles POST /api/dev/complete-constructions
// Instantly completes all upgrading buildings for the player.
func DevCompleteConstructions(w http.ResponseWriter, r *http.Request) {
	if !devCheck(w, r) {
		return
	}
	playerID := middleware.GetPlayerID(r)

	result, err := database.DB.Exec(`
		UPDATE buildings SET
			level = level + 1,
			is_upgrading = false,
			upgrade_finish_at = NULL,
			updated_at = now()
		WHERE planet_id IN (
			SELECT id FROM planets WHERE player_id = $1
		) AND is_upgrading = true
	`, playerID)
	if err != nil {
		log.Printf("DevCompleteConstructions: update failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"completed": rowsAffected,
	})
}

// DevCompleteResearch handles POST /api/dev/complete-research
// Instantly completes all active research for the player.
func DevCompleteResearch(w http.ResponseWriter, r *http.Request) {
	if !devCheck(w, r) {
		return
	}
	playerID := middleware.GetPlayerID(r)

	// Get names of researching techs before completing (for quest progress)
	rows, err := database.DB.Query(`
		SELECT tt.name FROM technologies t
		JOIN tech_types tt ON t.tech_type_id = tt.id
		WHERE t.player_id = $1 AND t.is_researching = true
	`, playerID)
	if err == nil {
		var techNames []string
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err == nil {
				techNames = append(techNames, name)
			}
		}
		rows.Close()

		// Complete the research
		result, err2 := database.DB.Exec(`
			UPDATE technologies SET
				level = level + 1,
				is_researching = false,
				research_finish_at = NULL,
				updated_at = now()
			WHERE player_id = $1 AND is_researching = true
		`, playerID)
		if err2 != nil {
			log.Printf("DevCompleteResearch: update failed: %v", err2)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Update quest progress for each completed research
		for _, name := range techNames {
			services.UpdateQuestProgress(playerID, "research_tech", name, 1)
		}

		rowsAffected, _ := result.RowsAffected()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":    "ok",
			"completed": rowsAffected,
		})
		return
	}

	// Fallback if query fails
	result, err := database.DB.Exec(`
		UPDATE technologies SET
			level = level + 1,
			is_researching = false,
			research_finish_at = NULL,
			updated_at = now()
		WHERE player_id = $1 AND is_researching = true
	`, playerID)
	if err != nil {
		log.Printf("DevCompleteResearch: update failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":    "ok",
		"completed": rowsAffected,
	})
}

// DevCompleteShipBuilds handles POST /api/dev/complete-ship-builds
// Instantly completes all active ship production slots.
// For each active build, creates ship_instances for the remaining quantity.
func DevCompleteShipBuilds(w http.ResponseWriter, r *http.Request) {
	if !devCheck(w, r) {
		return
	}
	playerID := middleware.GetPlayerID(r)

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("DevCompleteShipBuilds: tx begin failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Find active builds in ships table
	rows, err := tx.Query(`
		SELECT id, ship_design_id, build_quantity
		FROM ships
		WHERE player_id = $1 AND is_building = true
	`, playerID)
	if err != nil {
		log.Printf("DevCompleteShipBuilds: query failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	type buildInfo struct {
		ID           string
		ShipDesignID string
		BuildQty     int
	}
	var builds []buildInfo
	for rows.Next() {
		var b buildInfo
		if err := rows.Scan(&b.ID, &b.ShipDesignID, &b.BuildQty); err != nil {
			rows.Close()
			log.Printf("DevCompleteShipBuilds: scan failed: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
		builds = append(builds, b)
	}
	rows.Close()

	totalCreated := 0
	for _, b := range builds {
		// Get hull_type_id from ship_design
		var hullTypeID int
		err := tx.QueryRow(`SELECT hull_type_id FROM ship_designs WHERE id = $1`, b.ShipDesignID).Scan(&hullTypeID)
		if err != nil {
			log.Printf("DevCompleteShipBuilds: get hull_type failed: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Create ship instances for the build quantity
		for i := 0; i < b.BuildQty; i++ {
			_, err := tx.Exec(`
				INSERT INTO ship_instances (player_id, ship_design_id, hull_type_id)
				VALUES ($1, $2, $3)
			`, playerID, b.ShipDesignID, hullTypeID)
			if err != nil {
				log.Printf("DevCompleteShipBuilds: insert ship_instance failed: %v", err)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
				return
			}
		}
		totalCreated += b.BuildQty

		// Mark build as complete
		_, err = tx.Exec(`
			UPDATE ships SET
				quantity = quantity + build_quantity,
				is_building = false,
				build_quantity = 0,
				build_finish_at = NULL,
				updated_at = now()
			WHERE id = $1
		`, b.ID)
		if err != nil {
			log.Printf("DevCompleteShipBuilds: update ships failed: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("DevCompleteShipBuilds: commit failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "ok",
		"ships_created": totalCreated,
	})
}

// DevGiveBlueprint handles POST /api/dev/give-blueprint
// Body: {"blueprint_id": int, "level": int}
// Upserts a player blueprint with specified research level.
func DevGiveBlueprint(w http.ResponseWriter, r *http.Request) {
	if !devCheck(w, r) {
		return
	}
	playerID := middleware.GetPlayerID(r)

	var body struct {
		BlueprintID int `json:"blueprint_id"`
		Level       int `json:"level"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if body.BlueprintID <= 0 {
		http.Error(w, `{"error":"blueprint_id required"}`, http.StatusBadRequest)
		return
	}
	if body.Level < 1 || body.Level > 3 {
		body.Level = 1
	}

	_, err := database.DB.Exec(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, $3)
		ON CONFLICT (player_id, blueprint_id) DO UPDATE SET
			is_activated = true,
			research_level = $3
	`, playerID, body.BlueprintID, body.Level)
	if err != nil {
		log.Printf("DevGiveBlueprint: upsert failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Track quest progress for use_blueprint (quest targets use snake_case from module_types/hull_types)
	var questTarget string
	if err := database.DB.QueryRow(`
		SELECT COALESCE(mt.name, ht.name, b.name)
		FROM blueprints b
		LEFT JOIN module_types mt ON b.module_type_id = mt.id
		LEFT JOIN hull_types ht ON b.hull_type_id = ht.id
		WHERE b.id = $1`, body.BlueprintID).Scan(&questTarget); err == nil {
		services.UpdateQuestProgress(playerID, "use_blueprint", questTarget, 1)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// DevGiveCommander handles POST /api/dev/give-commander
// Body: {"commander_type_id": int, "star_rank": int}
// Creates a new commander from the commander_types reference table.
func DevGiveCommander(w http.ResponseWriter, r *http.Request) {
	if !devCheck(w, r) {
		return
	}
	playerID := middleware.GetPlayerID(r)

	var body struct {
		CommanderTypeID int `json:"commander_type_id"`
		StarRank        int `json:"star_rank"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if body.CommanderTypeID <= 0 {
		http.Error(w, `{"error":"commander_type_id required"}`, http.StatusBadRequest)
		return
	}
	if body.StarRank < 0 || body.StarRank > 15 {
		body.StarRank = 0
	}

	// Look up the commander type
	var name, rarity string
	var accuracy, dodge, speed, electron int
	err := database.DB.QueryRow(`
		SELECT display_name, rarity, base_accuracy, base_dodge, base_speed, base_electron
		FROM commander_types WHERE id = $1
	`, body.CommanderTypeID).Scan(&name, &rarity, &accuracy, &dodge, &speed, &electron)
	if err != nil {
		log.Printf("DevGiveCommander: lookup commander_type failed: %v", err)
		http.Error(w, `{"error":"commander type not found"}`, http.StatusNotFound)
		return
	}

	var commanderID string
	err = database.DB.QueryRow(`
		INSERT INTO commanders (player_id, name, rarity, star_rank, accuracy, dodge, speed, electron, is_deployed)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, false)
		RETURNING id
	`, playerID, name, rarity, body.StarRank, accuracy, dodge, speed, electron).Scan(&commanderID)
	if err != nil {
		log.Printf("DevGiveCommander: insert failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "ok",
		"commander_id": commanderID,
	})
}

// DevGiveShips handles POST /api/dev/give-ships
// Body: {"ship_design_id": string, "quantity": int}
// Creates individual ship instances for the specified design.
func DevGiveShips(w http.ResponseWriter, r *http.Request) {
	if !devCheck(w, r) {
		return
	}
	playerID := middleware.GetPlayerID(r)

	var body struct {
		ShipDesignID string `json:"ship_design_id"`
		Quantity     int    `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if body.ShipDesignID == "" || body.Quantity <= 0 {
		http.Error(w, `{"error":"ship_design_id and positive quantity required"}`, http.StatusBadRequest)
		return
	}

	// Get hull_type_id from ship_design and verify ownership
	var hullTypeID int
	err := database.DB.QueryRow(`
		SELECT hull_type_id FROM ship_designs WHERE id = $1 AND player_id = $2
	`, body.ShipDesignID, playerID).Scan(&hullTypeID)
	if err != nil {
		log.Printf("DevGiveShips: design lookup failed: %v", err)
		http.Error(w, `{"error":"ship design not found or not owned by player"}`, http.StatusNotFound)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("DevGiveShips: tx begin failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	for i := 0; i < body.Quantity; i++ {
		_, err := tx.Exec(`
			INSERT INTO ship_instances (player_id, ship_design_id, hull_type_id)
			VALUES ($1, $2, $3)
		`, playerID, body.ShipDesignID, hullTypeID)
		if err != nil {
			log.Printf("DevGiveShips: insert ship_instance failed: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("DevGiveShips: commit failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "ok",
		"ships_created": body.Quantity,
	})
}

// DevSetBuildingLevel handles POST /api/dev/set-building-level
// Body: {"building_id": string, "level": int}
// Sets a building's level directly.
func DevSetBuildingLevel(w http.ResponseWriter, r *http.Request) {
	if !devCheck(w, r) {
		return
	}
	playerID := middleware.GetPlayerID(r)

	var body struct {
		BuildingID string `json:"building_id"`
		Level      int    `json:"level"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if body.BuildingID == "" || body.Level < 1 {
		http.Error(w, `{"error":"building_id and level >= 1 required"}`, http.StatusBadRequest)
		return
	}

	// Check max level for the building type
	var maxLevel int
	var btName string
	err := database.DB.QueryRow(`
		SELECT bt.max_level, bt.name FROM buildings b
		JOIN building_types bt ON b.building_type = bt.id
		WHERE b.id = $1 AND b.planet_id IN (
			SELECT id FROM planets WHERE player_id = $2
		)
	`, body.BuildingID, playerID).Scan(&maxLevel, &btName)
	if err != nil {
		http.Error(w, `{"error":"building not found"}`, http.StatusNotFound)
		return
	}
	if body.Level > maxLevel {
		http.Error(w, `{"error":"level exceeds max level `+fmt.Sprintf("%d", maxLevel)+` for `+btName+`"}`, http.StatusConflict)
		return
	}

	result, err := database.DB.Exec(`
		UPDATE buildings SET
			level = $1,
			is_upgrading = false,
			upgrade_finish_at = NULL,
			updated_at = now()
		WHERE id = $2 AND planet_id IN (
			SELECT id FROM planets WHERE player_id = $3
		)
	`, body.Level, body.BuildingID, playerID)
	if err != nil {
		log.Printf("DevSetBuildingLevel: update failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error":"building not found or not owned by player"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// DevCompleteRecycling handles POST /api/dev/complete-recycling
// Instantly completes all active recycling jobs for the player.
func DevCompleteRecycling(w http.ResponseWriter, r *http.Request) {
	if !devCheck(w, r) {
		return
	}
	playerID := middleware.GetPlayerID(r)

	result, err := database.DB.Exec(`
		UPDATE recycling_jobs SET
			completed_at = now()
		WHERE player_id = $1 AND completed_at IS NULL
	`, playerID)
	if err != nil {
		log.Printf("DevCompleteRecycling: update failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"completed": rowsAffected,
	})
}

// DevCreateCorp handles POST /api/dev/create-corp
// Body: {"name": string, "tag": string}
// Creates a corp for the player instantly (bypasses Alliance Center check)
func DevCreateCorp(w http.ResponseWriter, r *http.Request) {
	if !devCheck(w, r) {
		return
	}
	playerID := middleware.GetPlayerID(r)

	var body struct {
		Name string `json:"name"`
		Tag  string `json:"tag"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if body.Name == "" {
		body.Name = "Dev Corp"
	}
	if body.Tag == "" {
		body.Tag = "DEV"
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("DevCreateCorp: tx begin failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Delete existing corp membership
	_, err = tx.Exec(`DELETE FROM corp_members WHERE player_id = $1`, playerID)
	if err != nil {
		log.Printf("DevCreateCorp: delete existing membership failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Insert into corps
	var corpID string
	err = tx.QueryRow(`
		INSERT INTO corps (name, tag, description, leader_id, level, wealth, max_members)
		VALUES ($1, $2, 'Dev corp created instantly', $3, 1, 0, 50)
		RETURNING id
	`, body.Name, body.Tag, playerID).Scan(&corpID)

	if err != nil {
		log.Printf("DevCreateCorp: insert corps failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Insert into corp_members
	_, err = tx.Exec(`
		INSERT INTO corp_members (corp_id, player_id, role, contribution_points, daily_contribution)
		VALUES ($1, $2, 'leader', 0, 0)
	`, corpID, playerID)

	if err != nil {
		log.Printf("DevCreateCorp: insert corp_members failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("DevCreateCorp: commit failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"corp_id": corpID,
	})
}

// DevJoinCorp handles POST /api/dev/join-corp
// Body: {"corp_id": string}
// Joins a corp instantly (bypasses Alliance Center check)
func DevJoinCorp(w http.ResponseWriter, r *http.Request) {
	if !devCheck(w, r) {
		return
	}
	playerID := middleware.GetPlayerID(r)

	var body struct {
		CorpID string `json:"corp_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if body.CorpID == "" {
		http.Error(w, `{"error":"corp_id required"}`, http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("DevJoinCorp: tx begin failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Delete existing corp membership
	_, err = tx.Exec(`DELETE FROM corp_members WHERE player_id = $1`, playerID)
	if err != nil {
		log.Printf("DevJoinCorp: delete existing membership failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Insert into corp_members
	_, err = tx.Exec(`
		INSERT INTO corp_members (corp_id, player_id, role, contribution_points, daily_contribution)
		VALUES ($1, $2, 'member', 0, 0)
	`, body.CorpID, playerID)

	if err != nil {
		log.Printf("DevJoinCorp: insert failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("DevJoinCorp: commit failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// DevGiveCorpWealth handles POST /api/dev/give-corp-wealth
// Body: {"corp_id": string, "wealth": int64}
// Adds wealth to a corp directly
func DevGiveCorpWealth(w http.ResponseWriter, r *http.Request) {
	if !devCheck(w, r) {
		return
	}

	var body struct {
		CorpID string `json:"corp_id"`
		Wealth int64  `json:"wealth"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if body.CorpID == "" || body.Wealth <= 0 {
		http.Error(w, `{"error":"corp_id and positive wealth required"}`, http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("DevGiveCorpWealth: tx begin failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Add wealth
	_, err = tx.Exec(`
		UPDATE corps
		SET wealth = wealth + $1, updated_at = now()
		WHERE id = $2
	`, body.Wealth, body.CorpID)

	if err != nil {
		log.Printf("DevGiveCorpWealth: update failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Recalculate corp level
	var totalWealth int64
	err = tx.QueryRow(`SELECT wealth FROM corps WHERE id = $1`, body.CorpID).Scan(&totalWealth)
	if err == nil {
		newLevel := services.CalculateCorpLevel(totalWealth)
		_, err = tx.Exec(`UPDATE corps SET level = $1 WHERE id = $2`, newLevel, body.CorpID)
		if err != nil {
			log.Printf("DevGiveCorpWealth: update level failed: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("DevGiveCorpWealth: commit failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
