package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/errs"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
	"github.com/cryptomines-online/backend/internal/services"
)

// applyDevMode checks if dev-mode is enabled and returns 5 seconds if true, otherwise returns the original seconds
func applyDevModeResearch(r *http.Request, seconds int) int {
	if r.Header.Get("X-Dev-Mode") == "true" {
		return 5
	}
	return seconds
}

// techResearchResponse is a single tech with player progress for API responses.
type techResearchResponse struct {
	ID                int             `json:"id"`
	Name              string          `json:"name"`
	DisplayName       string          `json:"display_name"`
	Tree              string          `json:"tree"`
	MaxLevel          int             `json:"max_level"`
	Prerequisites     json.RawMessage `json:"prerequisites"`
	CurrentLevel      int             `json:"current_level"`
	IsResearching     bool            `json:"is_researching"`
	ResearchFinishAt  *time.Time      `json:"research_finish_at,omitempty"`
	CostNextLevel     *costInfo       `json:"cost_next_level,omitempty"`
	TimeNextLevel     *int            `json:"time_next_level_seconds,omitempty"`
	Effects           json.RawMessage `json:"effects"`
	Description       string          `json:"description"`
}

type costInfo struct {
	Metal int64 `json:"metal"`
	He3   int64 `json:"he3"`
	Gold  int64 `json:"gold"`
}

// ListResearch handles GET /api/research
// Returns all 7 trees with player progress.
func ListResearch(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	applyCompletedResearch(playerID)

	techs := getTreesWithProgress(playerID, "")
	if techs == nil {
		http.Error(w, `{"error":"failed to load research"}`, http.StatusInternalServerError)
		return
	}

	// Group by tree
	treeMap := map[string][]techResearchResponse{}
	for _, t := range techs {
		treeMap[t.Tree] = append(treeMap[t.Tree], t)
	}

	// Convert map to array of tree responses
	type treeResponse struct {
		Tree  string                 `json:"tree"`
		Techs []techResearchResponse `json:"techs"`
	}

	var trees []treeResponse
	for treeName, techs := range treeMap {
		trees = append(trees, treeResponse{
			Tree:  treeName,
			Techs: techs,
		})
	}

	// Get active research
	var active *activeResearchInfo
	rows, err := database.DB.Query(`
		SELECT t.id, t.player_id, t.tech_type, t.level, t.is_researching,
		       t.research_finish_at, t.created_at, t.updated_at,
		       tt.name, tt.display_name, tt.tree
		FROM technologies t
		JOIN tech_types tt ON t.tech_type = tt.id
		WHERE t.player_id = $1 AND t.is_researching = true
		LIMIT 1`, playerID)
	
	if err == nil {
		defer rows.Close()
		if rows.Next() {
			var a activeResearchInfo
			rows.Scan(
				&a.ID, &a.PlayerID, &a.TechType, &a.Level, &a.IsResearching,
				&a.ResearchFinishAt, &a.CreatedAt, &a.UpdatedAt,
				&a.TechName, &a.DisplayName, &a.Tree,
			)
			active = &a
		}
	}

	type response struct {
		Trees  []treeResponse      `json:"trees"`
		Active *activeResearchInfo `json:"active"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{Trees: trees, Active: active})
}

type activeResearchInfo struct {
	models.Technology
	TechName    string `json:"tech_name"`
	DisplayName string `json:"display_name"`
	Tree        string `json:"tree"`
}

// GetResearchTree handles GET /api/research/trees/{tree}
// Returns a single tree with player progress.
func GetResearchTree(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	treeName := r.PathValue("tree")

	validTrees := map[string]bool{
		"logistics_construction": true,
		"planetary_defense":      true,
		"ballistics_science":     true,
		"directional_science":    true,
		"missile_science":        true,
		"ship_based_science":     true,
		"ship_defense_science":   true,
	}
	if !validTrees[treeName] {
		http.Error(w, `{"error":"invalid tree name"}`, http.StatusBadRequest)
		return
	}

	applyCompletedResearch(playerID)

	techs := getTreesWithProgress(playerID, treeName)
	if techs == nil {
		http.Error(w, `{"error":"failed to load research tree"}`, http.StatusInternalServerError)
		return
	}

	type response struct {
		Tree  string                 `json:"tree"`
		Techs []techResearchResponse `json:"techs"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{Tree: treeName, Techs: techs})
}

// StartResearch handles POST /api/research/start
func StartResearch(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	var req struct {
		TechTypeID int `json:"tech_type_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	applyCompletedResearch(playerID)

	// Get tech type
	var tt models.TechType
	err := database.DB.QueryRow(
		`SELECT id, name, display_name, tree, prerequisites_json,
		        base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier,
		        base_time_seconds, time_multiplier, max_level, effects_json, description
		 FROM tech_types WHERE id = $1`, req.TechTypeID,
	).Scan(
		&tt.ID, &tt.Name, &tt.DisplayName, &tt.Tree, &tt.PrerequisitesJSON,
		&tt.BaseCostMetal, &tt.BaseCostHe3, &tt.BaseCostGold, &tt.CostMultiplier,
		&tt.BaseTimeSeconds, &tt.TimeMultiplier, &tt.MaxLevel, &tt.EffectsJSON, &tt.Description,
	)
	if err != nil {
		http.Error(w, `{"error":"unknown tech type"}`, http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Check if player already has active research in this tree
	var activeCount int
	err = tx.QueryRow(
		`SELECT COUNT(*) FROM technologies t
		 JOIN tech_types tt ON t.tech_type = tt.id
		 WHERE t.player_id = $1 AND tt.tree = $2 AND t.is_researching = true`,
		playerID, tt.Tree,
	).Scan(&activeCount)
	if err != nil {
		log.Printf("Failed to check active research: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if activeCount > 0 {
		errs.Conflict(
			fmt.Sprintf("Already researching a technology in the %s tree", tt.Tree),
			map[string]interface{}{
				"tree": tt.Tree,
				"hint": "Wait for current research to complete or cancel it first",
			},
		).WriteJSON(w, http.StatusConflict)
		return
	}

	// Check if player has a Technology Center built
	techCenterLevel := getTechCenterLevel(tx, playerID)
	if techCenterLevel < 1 {
		errs.Conflict(
			"You must build a Technology Center before you can conduct research",
			map[string]interface{}{
				"required_building": "Technology Center",
				"hint":              "Build a Technology Center on your homeworld to unlock research",
			},
		).WriteJSON(w, http.StatusConflict)
		return
	}

	// Get current level for this tech (or 0 if not started)
	var currentLevel int
	var techID *string
	err = tx.QueryRow(
		`SELECT id, level FROM technologies
		 WHERE player_id = $1 AND tech_type = $2`,
		playerID, tt.ID,
	).Scan(&techID, &currentLevel)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Failed to get current tech level: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Check max level
	targetLevel := currentLevel + 1
	if targetLevel > tt.MaxLevel {
		errs.MaxLevelReached(
			fmt.Sprintf("%s is already at max level %d", tt.DisplayName, tt.MaxLevel),
			tt.MaxLevel,
		).WriteJSON(w, http.StatusConflict)
		return
	}

	// Check prerequisites
	var prereqs []models.TechPrerequisite
	if err := json.Unmarshal(tt.PrerequisitesJSON, &prereqs); err != nil {
		log.Printf("Failed to parse prerequisites: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	for _, prereq := range prereqs {
		var prereqLevel int
		err = tx.QueryRow(
			`SELECT COALESCE(t.level, 0) FROM tech_types tt
			 LEFT JOIN technologies t ON t.tech_type = tt.id AND t.player_id = $1
			 WHERE tt.name = $2`,
			playerID, prereq.Tech,
		).Scan(&prereqLevel)
		if err != nil {
			log.Printf("Failed to check prereq %s: %v", prereq.Tech, err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
		if prereqLevel < prereq.Level {
			// Get prerequisite display name
			var prereqDisplayName string
			database.DB.QueryRow(`SELECT display_name FROM tech_types WHERE name = $1`, prereq.Tech).Scan(&prereqDisplayName)
			if prereqDisplayName == "" {
				prereqDisplayName = prereq.Tech
			}
			errs.PrerequisiteNotMet(
				fmt.Sprintf("%s requires %s level %d (current: %d)", tt.DisplayName, prereqDisplayName, prereq.Level, prereqLevel),
				prereqDisplayName,
				prereqLevel,
				prereq.Level,
			).WriteJSON(w, http.StatusConflict)
			return
		}
	}

	// Calculate cost for target level
	goldCost := calcLevelCost(tt.BaseCostGold, tt.CostMultiplier, targetLevel)
	metalCost := calcLevelCost(tt.BaseCostMetal, tt.CostMultiplier, targetLevel)
	he3Cost := calcLevelCost(tt.BaseCostHe3, tt.CostMultiplier, targetLevel)

	// Calculate research time (with Tech Center bonus)
	baseTime := calcLevelTime(tt.BaseTimeSeconds, tt.TimeMultiplier, targetLevel)
	// techCenterLevel already retrieved earlier for validation
	reduction := float64(techCenterLevel) * 0.03
	effectiveTime := int(math.Round(float64(baseTime) * (1.0 - reduction)))
	if effectiveTime < 1 {
		effectiveTime = 1
	}
	// Apply dev-mode if enabled
	effectiveTime = applyDevModeResearch(r, effectiveTime)

	// Get homeworld ID
	var homeworldID string
	err = tx.QueryRow(`SELECT id FROM planets WHERE player_id = $1 AND is_homeworld = true LIMIT 1`, playerID).Scan(&homeworldID)
	if err != nil {
		log.Printf("Failed to get homeworld: %v", err)
		errs.InternalError("Failed to find homeworld").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	// Deduct resources from homeworld with detailed error
	var res resourceState
	res.Metal, res.He3, res.Gold, err = checkAndDeductResources(
		tx, homeworldID,
		metalCost, he3Cost, goldCost,
		fmt.Sprintf("%s level %d", tt.DisplayName, targetLevel),
	)
	if err != nil {
		if err.Error() == fmt.Sprintf("insufficient resources for %s level %d", tt.DisplayName, targetLevel) {
			errs.InsufficientResources(
				fmt.Sprintf("Not enough resources to research %s level %d", tt.DisplayName, targetLevel),
				map[string]int64{
					"metal": metalCost,
					"he3":   he3Cost,
					"gold":  goldCost,
				},
				map[string]int64{
					"metal": res.Metal,
					"he3":   res.He3,
					"gold":  res.Gold,
				},
			).WriteJSON(w, http.StatusConflict)
			return
		}
		log.Printf("Failed to deduct resources: %v", err)
		errs.InternalError("Failed to deduct resources").WriteJSON(w, http.StatusInternalServerError)
		return
	}

	// Start research (upsert technologies row)
	finishAt := time.Now().Add(time.Duration(effectiveTime) * time.Second)
	var tech models.Technology
	err = tx.QueryRow(
		`INSERT INTO technologies (player_id, tech_type, level, is_researching, research_finish_at)
		 VALUES ($1, $2, $3, true, $4)
		 ON CONFLICT (player_id, tech_type)
		 DO UPDATE SET is_researching = true, research_finish_at = $4, updated_at = now()
		 RETURNING id, player_id, tech_type, level, is_researching, research_finish_at, created_at, updated_at`,
		playerID, tt.ID, currentLevel, finishAt,
	).Scan(
		&tech.ID, &tech.PlayerID, &tech.TechType, &tech.Level,
		&tech.IsResearching, &tech.ResearchFinishAt, &tech.CreatedAt, &tech.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to start research: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	type response struct {
		Technology models.Technology `json:"technology"`
		Resources  resourceState    `json:"resources"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{Technology: tech, Resources: res})
}

// CancelResearch handles POST /api/research/cancel
func CancelResearch(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	var req struct {
		TechTypeID int `json:"tech_type_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	var tech models.Technology
	err := database.DB.QueryRow(
		`UPDATE technologies
		 SET is_researching = false, research_finish_at = NULL, updated_at = now()
		 WHERE player_id = $1 AND tech_type = $2 AND is_researching = true
		 RETURNING id, player_id, tech_type, level, is_researching, research_finish_at, created_at, updated_at`,
		playerID, req.TechTypeID,
	).Scan(
		&tech.ID, &tech.PlayerID, &tech.TechType, &tech.Level,
		&tech.IsResearching, &tech.ResearchFinishAt, &tech.CreatedAt, &tech.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"no active research for this tech"}`, http.StatusConflict)
			return
		}
		log.Printf("Failed to cancel research: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Technology models.Technology `json:"technology"`
	}{Technology: tech})
}

// SpeedupResearch handles POST /api/research/speedup
// Cost: 3 vouchers/MP per 30 minutes of time reduction.
func SpeedupResearch(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	var req struct {
		TechTypeID     int `json:"tech_type_id"`
		SpeedupMinutes int `json:"speedup_minutes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.SpeedupMinutes <= 0 {
		http.Error(w, `{"error":"speedup_minutes must be positive"}`, http.StatusBadRequest)
		return
	}

	// Calculate voucher cost: 3 per 30 minutes (round up)
	voucherCost := int(math.Ceil(float64(req.SpeedupMinutes) / 30.0 * 3.0))

	var tech models.Technology
	err := database.DB.QueryRow(
		`SELECT id, player_id, tech_type, level, is_researching, research_finish_at, created_at, updated_at
		 FROM technologies
		 WHERE player_id = $1 AND tech_type = $2 AND is_researching = true`,
		playerID, req.TechTypeID,
	).Scan(
		&tech.ID, &tech.PlayerID, &tech.TechType, &tech.Level,
		&tech.IsResearching, &tech.ResearchFinishAt, &tech.CreatedAt, &tech.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"no active research for this tech"}`, http.StatusConflict)
			return
		}
		log.Printf("Failed to get research: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Reduce finish time
	newFinishAt := tech.ResearchFinishAt.Add(-time.Duration(req.SpeedupMinutes) * time.Minute)

	// If new finish time is in the past, complete immediately
	now := time.Now()
	if newFinishAt.Before(now) {
		newFinishAt = now
	}

	err = database.DB.QueryRow(
		`UPDATE technologies
		 SET research_finish_at = $1, updated_at = now()
		 WHERE id = $2
		 RETURNING id, player_id, tech_type, level, is_researching, research_finish_at, created_at, updated_at`,
		newFinishAt, tech.ID,
	).Scan(
		&tech.ID, &tech.PlayerID, &tech.TechType, &tech.Level,
		&tech.IsResearching, &tech.ResearchFinishAt, &tech.CreatedAt, &tech.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to speedup research: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Auto-complete if time is now
	applyCompletedResearch(playerID)

	// Re-read after possible completion
	database.DB.QueryRow(
		`SELECT id, player_id, tech_type, level, is_researching, research_finish_at, created_at, updated_at
		 FROM technologies WHERE id = $1`, tech.ID,
	).Scan(
		&tech.ID, &tech.PlayerID, &tech.TechType, &tech.Level,
		&tech.IsResearching, &tech.ResearchFinishAt, &tech.CreatedAt, &tech.UpdatedAt,
	)

	type response struct {
		Technology    models.Technology `json:"technology"`
		VouchersSpent int              `json:"vouchers_spent"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{Technology: tech, VouchersSpent: voucherCost})
}

// GetActiveResearch handles GET /api/research/active
// Returns all currently active researches across all trees.
func GetActiveResearch(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	applyCompletedResearch(playerID)

	rows, err := database.DB.Query(
		`SELECT t.id, t.player_id, t.tech_type, t.level, t.is_researching,
		        t.research_finish_at, t.created_at, t.updated_at,
		        tt.name, tt.display_name, tt.tree
		 FROM technologies t
		 JOIN tech_types tt ON t.tech_type = tt.id
		 WHERE t.player_id = $1 AND t.is_researching = true`, playerID,
	)
	if err != nil {
		log.Printf("Failed to get active research: %v", err)
		http.Error(w, `{"error":"failed to get active research"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type activeResearch struct {
		models.Technology
		TechName    string `json:"tech_name"`
		DisplayName string `json:"display_name"`
		Tree        string `json:"tree"`
	}

	var active []activeResearch
	for rows.Next() {
		var a activeResearch
		err := rows.Scan(
			&a.ID, &a.PlayerID, &a.TechType, &a.Level, &a.IsResearching,
			&a.ResearchFinishAt, &a.CreatedAt, &a.UpdatedAt,
			&a.TechName, &a.DisplayName, &a.Tree,
		)
		if err != nil {
			log.Printf("Failed to scan active research: %v", err)
			continue
		}
		active = append(active, a)
	}

	if active == nil {
		active = []activeResearch{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(active)
}

// applyCompletedResearch auto-completes any finished research timers.
func applyCompletedResearch(playerID string) {
	now := time.Now()

	// Find completed research before updating
	rows, err := database.DB.Query(`
		SELECT t.id, tt.name, t.level
		FROM technologies t
		JOIN tech_types tt ON t.tech_type = tt.id
		WHERE t.player_id = $1 AND t.is_researching = true AND t.research_finish_at <= $2
	`, playerID, now)
	if err != nil {
		log.Printf("Failed to query completed research: %v", err)
		return
	}

	type completedResearch struct {
		id       string
		techName string
		newLevel int
	}
	var completed []completedResearch

	for rows.Next() {
		var id, techName string
		var currentLevel int
		rows.Scan(&id, &techName, &currentLevel)
		completed = append(completed, completedResearch{
			id:       id,
			techName: techName,
			newLevel: currentLevel + 1,
		})
	}
	rows.Close()

	if len(completed) == 0 {
		return
	}

	// Apply the research completion
	_, err = database.DB.Exec(`
		UPDATE technologies
		SET level = level + 1, is_researching = false, research_finish_at = NULL, updated_at = now()
		WHERE player_id = $1 AND is_researching = true AND research_finish_at <= $2
	`, playerID, now)
	if err != nil {
		log.Printf("Failed to apply completed research: %v", err)
		return
	}

	// Update quest progress for each completed research
	for _, research := range completed {
		services.UpdateQuestProgress(playerID, "research_tech", research.techName, 1)
	}
}

// getTreesWithProgress loads tech types with player progress.
// If treeName is empty, returns all trees.
func getTreesWithProgress(playerID, treeName string) []techResearchResponse {
	var query string
	var args []interface{}

	if treeName == "" {
		query = `SELECT tt.id, tt.name, tt.display_name, tt.tree, tt.max_level,
		                tt.prerequisites_json, tt.base_cost_gold, tt.cost_multiplier,
		                tt.base_time_seconds, tt.time_multiplier, tt.effects_json, tt.description,
		                COALESCE(t.level, 0), COALESCE(t.is_researching, false),
		                t.research_finish_at
		         FROM tech_types tt
		         LEFT JOIN technologies t ON t.tech_type = tt.id AND t.player_id = $1
		         ORDER BY tt.tree, tt.id`
		args = []interface{}{playerID}
	} else {
		query = `SELECT tt.id, tt.name, tt.display_name, tt.tree, tt.max_level,
		                tt.prerequisites_json, tt.base_cost_gold, tt.cost_multiplier,
		                tt.base_time_seconds, tt.time_multiplier, tt.effects_json, tt.description,
		                COALESCE(t.level, 0), COALESCE(t.is_researching, false),
		                t.research_finish_at
		         FROM tech_types tt
		         LEFT JOIN technologies t ON t.tech_type = tt.id AND t.player_id = $1
		         WHERE tt.tree = $2
		         ORDER BY tt.id`
		args = []interface{}{playerID, treeName}
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Failed to query tech trees: %v", err)
		return nil
	}
	defer rows.Close()

	techCenterLevel := getTechCenterLevelFromDB(playerID)

	var techs []techResearchResponse
	for rows.Next() {
		var t techResearchResponse
		var baseCostGold int64
		var costMult, timeMult float64
		var baseTimeSec int

		err := rows.Scan(
			&t.ID, &t.Name, &t.DisplayName, &t.Tree, &t.MaxLevel,
			&t.Prerequisites, &baseCostGold, &costMult,
			&baseTimeSec, &timeMult, &t.Effects, &t.Description,
			&t.CurrentLevel, &t.IsResearching, &t.ResearchFinishAt,
		)
		if err != nil {
			log.Printf("Failed to scan tech: %v", err)
			continue
		}

		// Calculate next level cost/time if not at max
		if t.CurrentLevel < t.MaxLevel && !t.IsResearching {
			nextLevel := t.CurrentLevel + 1
			goldCost := calcLevelCost(baseCostGold, costMult, nextLevel)
			t.CostNextLevel = &costInfo{Metal: 0, He3: 0, Gold: goldCost}

			baseTime := calcLevelTime(baseTimeSec, timeMult, nextLevel)
			reduction := float64(techCenterLevel) * 0.03
			effectiveTime := int(math.Round(float64(baseTime) * (1.0 - reduction)))
			if effectiveTime < 1 {
				effectiveTime = 1
			}
			t.TimeNextLevel = &effectiveTime
		}

		techs = append(techs, t)
	}

	if techs == nil {
		techs = []techResearchResponse{}
	}

	return techs
}

// calcLevelCost computes gold cost for a given level.
// Cost(level) = baseCost * multiplier^(level-1)
func calcLevelCost(baseCost int64, multiplier float64, level int) int64 {
	if level <= 1 {
		return baseCost
	}
	return int64(math.Round(float64(baseCost) * math.Pow(multiplier, float64(level-1))))
}

// calcLevelTime computes base research time in seconds for a given level.
// Time(level) = baseTime * multiplier^(level-1)
func calcLevelTime(baseTimeSeconds int, multiplier float64, level int) int {
	if level <= 1 {
		return baseTimeSeconds
	}
	return int(math.Round(float64(baseTimeSeconds) * math.Pow(multiplier, float64(level-1))))
}

// getTechCenterLevel gets the player's Technology Center building level within a transaction.
func getTechCenterLevel(tx *sql.Tx, playerID string) int {
	var level int
	err := tx.QueryRow(
		`SELECT COALESCE(MAX(b.level), 0)
		 FROM buildings b
		 JOIN building_types bt ON b.building_type = bt.id
		 JOIN planets p ON b.planet_id = p.id
		 WHERE p.player_id = $1 AND bt.name = 'technology_center'`,
		playerID,
	).Scan(&level)
	if err != nil {
		return 0
	}
	return level
}

// getTechCenterLevelFromDB gets the player's Technology Center level directly from DB.
func getTechCenterLevelFromDB(playerID string) int {
	var level int
	err := database.DB.QueryRow(
		`SELECT COALESCE(MAX(b.level), 0)
		 FROM buildings b
		 JOIN building_types bt ON b.building_type = bt.id
		 JOIN planets p ON b.planet_id = p.id
		 WHERE p.player_id = $1 AND bt.name = 'technology_center'`,
		playerID,
	).Scan(&level)
	if err != nil {
		return 0
	}
	return level
}
