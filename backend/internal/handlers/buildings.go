package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/errs"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
	"github.com/cryptomines-online/backend/internal/services"
)

// applyDevMode checks if dev-mode is enabled and returns 5 seconds if true, otherwise returns the original seconds
func applyDevMode(r *http.Request, seconds int) int {
	if r.Header.Get("X-Dev-Mode") == "true" {
		return 5
	}
	return seconds
}

// getConstructionSlots calculates the number of construction slots available to a player.
// Base: 1 slot, +1 per level of Concurrent Construction tech, + active buffs from items.
func getConstructionSlots(playerID string) int {
	baseSlots := 1

	// Get Concurrent Construction tech level
	var techLevel int
	err := database.DB.QueryRow(`
		SELECT COALESCE(t.level, 0)
		FROM tech_types tt
		LEFT JOIN technologies t ON t.tech_type = tt.id AND t.player_id = $1
		WHERE tt.name = 'concurrent_construction'
	`, playerID).Scan(&techLevel)

	if err != nil {
		log.Printf("Failed to get construction slots: %v", err)
		return baseSlots
	}

	// Check for active construction_slots buff (from Construction Card items)
	var buffSlots int
	err = database.DB.QueryRow(`
		SELECT COALESCE(buff_value, 0)
		FROM active_buffs
		WHERE player_id = $1 AND buff_type = 'construction_slots' AND expires_at > now()
	`, playerID).Scan(&buffSlots)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Failed to get construction slots buff: %v", err)
	}

	return baseSlots + techLevel + buffSlots
}

// ListBuildings handles GET /api/planets/{id}/buildings
func ListBuildings(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	planetID := r.PathValue("id")

	if !verifyPlanetOwnership(planetID, playerID) {
		http.Error(w, `{"error":"planet not found"}`, http.StatusNotFound)
		return
	}

	// Auto-complete finished upgrades
	applyCompletedUpgrades(planetID)

	rows, err := database.DB.Query(
		`SELECT b.id, b.planet_id, b.building_type, b.grid_col, b.grid_row, b.level, b.is_upgrading,
		        b.upgrade_finish_at, b.created_at, b.updated_at,
		        bt.name, bt.display_name, bt.category, bt.base, bt.max_level
		 FROM buildings b
		 JOIN building_types bt ON b.building_type = bt.id
		 WHERE b.planet_id = $1
		 ORDER BY bt.category, bt.name`, planetID,
	)
	if err != nil {
		log.Printf("Failed to list buildings: %v", err)
		http.Error(w, `{"error":"failed to list buildings"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	buildings := []models.BuildingWithType{}
	for rows.Next() {
		var b models.BuildingWithType
		err := rows.Scan(
			&b.ID, &b.PlanetID, &b.BuildingType, &b.GridCol, &b.GridRow, &b.Level, &b.IsUpgrading,
			&b.UpgradeFinishAt, &b.CreatedAt, &b.UpdatedAt,
			&b.TypeName, &b.DisplayName, &b.Category, &b.BaseName, &b.MaxLevel,
		)
		if err != nil {
			log.Printf("Failed to scan building: %v", err)
			continue
		}
		buildings = append(buildings, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildings)
}

type constructRequest struct {
	BuildingTypeName string `json:"building_type"`
	GridCol          int    `json:"grid_col"`
	GridRow          int    `json:"grid_row"`
}

type buildingResponse struct {
	Building  models.Building `json:"building"`
	Resources resourceState   `json:"resources"`
}

type resourceState struct {
	Metal int64 `json:"metal"`
	He3   int64 `json:"he3"`
	Gold  int64 `json:"gold"`
}

// ConstructBuilding handles POST /api/planets/{id}/buildings
func ConstructBuilding(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	planetID := r.PathValue("id")

	if !verifyPlanetOwnership(planetID, playerID) {
		http.Error(w, `{"error":"planet not found"}`, http.StatusNotFound)
		return
	}

	var req constructRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Get building type
	var bt models.BuildingType
	err := database.DB.QueryRow(
		`SELECT id, name, display_name, category, base, base_cost_metal, base_cost_he3,
		        base_cost_gold, base_time_seconds, cost_multiplier, time_multiplier,
		        base_production_per_hour, production_multiplier, max_level,
		        max_count_per_planet, prerequisite_building, prerequisite_level,
		        civic_center_req_per_level, description
		 FROM building_types WHERE name = $1`, req.BuildingTypeName,
	).Scan(
		&bt.ID, &bt.Name, &bt.DisplayName, &bt.Category, &bt.Base,
		&bt.BaseCostMetal, &bt.BaseCostHe3, &bt.BaseCostGold,
		&bt.BaseTimeSeconds, &bt.CostMultiplier, &bt.TimeMultiplier,
		&bt.BaseProductionPerHour, &bt.ProductionMultiplier, &bt.MaxLevel,
		&bt.MaxCountPerPlanet, &bt.PrerequisiteBuilding, &bt.PrerequisiteLevel,
		&bt.CivicCenterReqPerLevel, &bt.Description,
	)
	if err != nil {
		http.Error(w, `{"error":"unknown building type"}`, http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Check construction slot limit (max 2 simultaneous upgrades/constructions)
	var activeBuilds int
	err = tx.QueryRow(
		`SELECT COUNT(*) FROM buildings WHERE planet_id = $1 AND is_upgrading = true`,
		planetID,
	).Scan(&activeBuilds)
	if err != nil {
		log.Printf("Failed to count active builds: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	maxSlots := getConstructionSlots(playerID)
	if activeBuilds >= maxSlots {
		errs.Conflict(
			fmt.Sprintf("All %d construction slots are in use", maxSlots),
			map[string]interface{}{
				"slots_used":      activeBuilds,
				"slots_available": maxSlots,
				"hint":            "Upgrade Concurrent Construction tech or wait for current construction to finish",
			},
		).WriteJSON(w, http.StatusConflict)
		return
	}

	// Check max count per planet
	var currentCount int
	err = tx.QueryRow(
		`SELECT COUNT(*) FROM buildings WHERE planet_id = $1 AND building_type = $2`,
		planetID, bt.ID,
	).Scan(&currentCount)
	if err != nil {
		log.Printf("Failed to count buildings: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if currentCount >= bt.MaxCountPerPlanet {
		errs.MaxCountReached(
			fmt.Sprintf("Maximum %s count reached (%d/%d)", bt.DisplayName, currentCount, bt.MaxCountPerPlanet),
			bt.MaxCountPerPlanet,
			currentCount,
		).WriteJSON(w, http.StatusConflict)
		return
	}

	// Check prerequisite building
	if bt.PrerequisiteBuilding != nil {
		var prereqLevel int
		err = tx.QueryRow(
			`SELECT COALESCE(MAX(b.level), 0)
			 FROM buildings b
			 JOIN building_types bt2 ON b.building_type = bt2.id
			 WHERE b.planet_id = $1 AND bt2.name = $2`,
			planetID, *bt.PrerequisiteBuilding,
		).Scan(&prereqLevel)
		if err != nil {
			log.Printf("Failed to check prerequisite: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
		if prereqLevel < bt.PrerequisiteLevel {
			var prereqDisplayName string
			database.DB.QueryRow(`SELECT display_name FROM building_types WHERE name = $1`, *bt.PrerequisiteBuilding).Scan(&prereqDisplayName)
			if prereqDisplayName == "" {
				prereqDisplayName = *bt.PrerequisiteBuilding
			}
			errs.PrerequisiteNotMet(
				fmt.Sprintf("%s requires %s level %d (current: %d)", bt.DisplayName, prereqDisplayName, bt.PrerequisiteLevel, prereqLevel),
				prereqDisplayName,
				prereqLevel,
				bt.PrerequisiteLevel,
			).WriteJSON(w, http.StatusConflict)
			return
		}
	}

	// Get cost from lookup table (exact wiki data) or fallback to formula
	levelCost := services.GetBuildingLevelCost(
		database.DB, bt.Name, 1,
		bt.BaseCostMetal, bt.BaseCostHe3, bt.BaseCostGold,
		bt.CostMultiplier, bt.BaseTimeSeconds, bt.TimeMultiplier,
	)

	// Apply tech bonuses (build cost reduction and build speed)
	techBonuses, err := services.GetPlayerTechBonuses(playerID)
	if err != nil {
		log.Printf("Failed to get tech bonuses: %v", err)
		// Continue without bonuses rather than failing
		techBonuses = &services.TechBonuses{}
	}

	// Apply cost reduction bonus
	costReduction := techBonuses.BuildCostReduction / 100.0
	levelCost.MetalCost = int64(float64(levelCost.MetalCost) * (1.0 - costReduction))
	levelCost.He3Cost = int64(float64(levelCost.He3Cost) * (1.0 - costReduction))
	levelCost.GoldCost = int64(float64(levelCost.GoldCost) * (1.0 - costReduction))

	// Apply build speed bonus (reduces time)
	speedBonus := techBonuses.BuildSpeed / 100.0
	levelCost.BuildTimeSeconds = int(float64(levelCost.BuildTimeSeconds) * (1.0 - speedBonus))
	if levelCost.BuildTimeSeconds < 1 {
		levelCost.BuildTimeSeconds = 1
	}

	// Deduct resources with detailed error handling
	var res resourceState
	res.Metal, res.He3, res.Gold, err = checkAndDeductResources(
		tx, planetID,
		levelCost.MetalCost, levelCost.He3Cost, levelCost.GoldCost,
		fmt.Sprintf("%s level 1", bt.DisplayName),
	)
	if err != nil {
		if err.Error() == fmt.Sprintf("insufficient resources for %s level 1", bt.DisplayName) {
			errs.InsufficientResources(
				fmt.Sprintf("Not enough resources to construct %s", bt.DisplayName),
				map[string]int64{
					"metal": levelCost.MetalCost,
					"he3":   levelCost.He3Cost,
					"gold":  levelCost.GoldCost,
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

	// Validate grid coordinates
	if req.GridCol < 0 || req.GridRow < 0 {
		http.Error(w, `{"error":"invalid grid coordinates"}`, http.StatusBadRequest)
		return
	}

	// Create building (starts at level 0, upgrading to level 1)
	// applyCompletedUpgrades does level+1 when done, so 0->1 for initial construction
	effectiveTime := applyDevMode(r, levelCost.BuildTimeSeconds)
	finishAt := time.Now().Add(time.Duration(effectiveTime) * time.Second)
	var building models.Building
	err = tx.QueryRow(
		`INSERT INTO buildings (planet_id, building_type, grid_col, grid_row, level, is_upgrading, upgrade_finish_at)
		 VALUES ($1, $2, $3, $4, 0, true, $5)
		 RETURNING id, planet_id, building_type, grid_col, grid_row, level, is_upgrading, upgrade_finish_at, created_at, updated_at`,
		planetID, bt.ID, req.GridCol, req.GridRow, finishAt,
	).Scan(
		&building.ID, &building.PlanetID, &building.BuildingType, &building.GridCol, &building.GridRow,
		&building.Level, &building.IsUpgrading, &building.UpgradeFinishAt, &building.CreatedAt, &building.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to create building: %v", err)
		http.Error(w, `{"error":"failed to create building"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(buildingResponse{Building: building, Resources: res})
}

// UpgradeBuilding handles POST /api/planets/{id}/buildings/{buildingId}/upgrade
func UpgradeBuilding(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	planetID := r.PathValue("id")
	buildingID := r.PathValue("buildingId")

	if !verifyPlanetOwnership(planetID, playerID) {
		http.Error(w, `{"error":"planet not found"}`, http.StatusNotFound)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Check construction slot limit
	var activeBuilds int
	err = tx.QueryRow(
		`SELECT COUNT(*) FROM buildings WHERE planet_id = $1 AND is_upgrading = true`,
		planetID,
	).Scan(&activeBuilds)
	if err != nil {
		log.Printf("Failed to count active builds: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	maxSlots := getConstructionSlots(playerID)
	if activeBuilds >= maxSlots {
		errs.Conflict(
			fmt.Sprintf("All %d construction slots are in use", maxSlots),
			map[string]interface{}{
				"slots_used":      activeBuilds,
				"slots_available": maxSlots,
				"hint":            "Upgrade Concurrent Construction tech or wait for current construction to finish",
			},
		).WriteJSON(w, http.StatusConflict)
		return
	}

	// Get the building with its type info
	var b models.Building
	var bt models.BuildingType
	err = tx.QueryRow(
		`SELECT b.id, b.planet_id, b.building_type, b.grid_col, b.grid_row, b.level, b.is_upgrading,
		        b.upgrade_finish_at, b.created_at, b.updated_at,
		        bt.id, bt.name, bt.base_cost_metal, bt.base_cost_he3, bt.base_cost_gold,
		        bt.base_time_seconds, bt.cost_multiplier, bt.time_multiplier,
		        bt.max_level, bt.civic_center_req_per_level
		 FROM buildings b
		 JOIN building_types bt ON b.building_type = bt.id
		 WHERE b.id = $1 AND b.planet_id = $2`,
		buildingID, planetID,
	).Scan(
		&b.ID, &b.PlanetID, &b.BuildingType, &b.GridCol, &b.GridRow, &b.Level, &b.IsUpgrading,
		&b.UpgradeFinishAt, &b.CreatedAt, &b.UpdatedAt,
		&bt.ID, &bt.Name, &bt.BaseCostMetal, &bt.BaseCostHe3, &bt.BaseCostGold,
		&bt.BaseTimeSeconds, &bt.CostMultiplier, &bt.TimeMultiplier,
		&bt.MaxLevel, &bt.CivicCenterReqPerLevel,
	)
	if err != nil {
		http.Error(w, `{"error":"building not found"}`, http.StatusNotFound)
		return
	}

	// Get building display name
	var displayName string
	database.DB.QueryRow(`SELECT display_name FROM building_types WHERE id = $1`, bt.ID).Scan(&displayName)
	if displayName == "" {
		displayName = bt.Name
	}

	// Check if already upgrading
	if b.IsUpgrading {
		var finishAt string
		if b.UpgradeFinishAt != nil {
			finishAt = b.UpgradeFinishAt.Format("2006-01-02 15:04:05")
		}
		errs.Conflict(
			fmt.Sprintf("%s is already upgrading", displayName),
			map[string]interface{}{
				"building_id":      b.ID,
				"current_level":    b.Level,
				"upgrade_finish_at": finishAt,
			},
		).WriteJSON(w, http.StatusConflict)
		return
	}

	// Check max level
	targetLevel := b.Level + 1
	if targetLevel > bt.MaxLevel {
		errs.MaxLevelReached(
			fmt.Sprintf("%s is already at max level %d", displayName, bt.MaxLevel),
			bt.MaxLevel,
		).WriteJSON(w, http.StatusConflict)
		return
	}

	// Check civic center requirement
	if bt.CivicCenterReqPerLevel && bt.Name != "civic_center" {
		var ccLevel int
		err = tx.QueryRow(
			`SELECT COALESCE(MAX(b.level), 0)
			 FROM buildings b
			 JOIN building_types bt2 ON b.building_type = bt2.id
			 WHERE b.planet_id = $1 AND bt2.name = 'civic_center'`,
			planetID,
		).Scan(&ccLevel)
		if err != nil {
			log.Printf("Failed to check civic center: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
		if ccLevel < targetLevel {
			errs.PrerequisiteNotMet(
				fmt.Sprintf("Civic Center must be level %d to upgrade %s to level %d (current: %d)", targetLevel, displayName, targetLevel, ccLevel),
				"Civic Center",
				ccLevel,
				targetLevel,
			).WriteJSON(w, http.StatusConflict)
			return
		}
	}

	// Get cost from lookup table (exact wiki data) or fallback to formula
	levelCost := services.GetBuildingLevelCost(
		database.DB, bt.Name, targetLevel,
		bt.BaseCostMetal, bt.BaseCostHe3, bt.BaseCostGold,
		bt.CostMultiplier, bt.BaseTimeSeconds, bt.TimeMultiplier,
	)

	// Apply tech bonuses (build cost reduction and build speed)
	techBonuses, err := services.GetPlayerTechBonuses(playerID)
	if err != nil {
		log.Printf("Failed to get tech bonuses: %v", err)
		// Continue without bonuses rather than failing
		techBonuses = &services.TechBonuses{}
	}

	// Apply cost reduction bonus
	costReduction := techBonuses.BuildCostReduction / 100.0
	levelCost.MetalCost = int64(float64(levelCost.MetalCost) * (1.0 - costReduction))
	levelCost.He3Cost = int64(float64(levelCost.He3Cost) * (1.0 - costReduction))
	levelCost.GoldCost = int64(float64(levelCost.GoldCost) * (1.0 - costReduction))

	// Apply build speed bonus (reduces time)
	speedBonus := techBonuses.BuildSpeed / 100.0
	levelCost.BuildTimeSeconds = int(float64(levelCost.BuildTimeSeconds) * (1.0 - speedBonus))
	if levelCost.BuildTimeSeconds < 1 {
		levelCost.BuildTimeSeconds = 1
	}

	// Deduct resources with detailed error handling
	var res resourceState
	res.Metal, res.He3, res.Gold, err = checkAndDeductResources(
		tx, planetID,
		levelCost.MetalCost, levelCost.He3Cost, levelCost.GoldCost,
		fmt.Sprintf("%s level %d", displayName, targetLevel),
	)
	if err != nil {
		if err.Error() == fmt.Sprintf("insufficient resources for %s level %d", displayName, targetLevel) {
			errs.InsufficientResources(
				fmt.Sprintf("Not enough resources to upgrade %s to level %d", displayName, targetLevel),
				map[string]int64{
					"metal": levelCost.MetalCost,
					"he3":   levelCost.He3Cost,
					"gold":  levelCost.GoldCost,
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

	// Start upgrade
	effectiveUpgradeTime := applyDevMode(r, levelCost.BuildTimeSeconds)
	finishAt := time.Now().Add(time.Duration(effectiveUpgradeTime) * time.Second)
	err = tx.QueryRow(
		`UPDATE buildings
		 SET is_upgrading = true, upgrade_finish_at = $1, updated_at = now()
		 WHERE id = $2
		 RETURNING id, planet_id, building_type, grid_col, grid_row, level, is_upgrading, upgrade_finish_at, created_at, updated_at`,
		finishAt, buildingID,
	).Scan(
		&b.ID, &b.PlanetID, &b.BuildingType, &b.GridCol, &b.GridRow, &b.Level, &b.IsUpgrading,
		&b.UpgradeFinishAt, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to start upgrade: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildingResponse{Building: b, Resources: res})
}

// ListBuildingTypes handles GET /api/building-types (public)
func ListBuildingTypes(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(
		`SELECT name, display_name, category, base_cost_metal, base_cost_he3,
		        base_cost_gold, base_time_seconds, cost_multiplier, time_multiplier,
		        max_level, max_count_per_planet, civic_center_req_per_level,
		        base_production_per_hour, production_multiplier
		 FROM building_types ORDER BY category, name`,
	)
	if err != nil {
		log.Printf("Failed to list building types: %v", err)
		http.Error(w, `{"error":"failed to list building types"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type buildingTypeRow struct {
		Name                   string  `json:"name"`
		DisplayName            string  `json:"display_name"`
		Category               string  `json:"category"`
		BaseCostMetal          int64   `json:"base_cost_metal"`
		BaseCostHe3            int64   `json:"base_cost_he3"`
		BaseCostGold           int64   `json:"base_cost_gold"`
		BaseTimeSeconds        int     `json:"base_time_seconds"`
		CostMultiplier         float64 `json:"cost_multiplier"`
		TimeMultiplier         float64 `json:"time_multiplier"`
		MaxLevel               int     `json:"max_level"`
		MaxCountPerPlanet      int     `json:"max_count_per_planet"`
		CivicCenterReqPerLevel bool    `json:"civic_center_req_per_level"`
		BaseProductionPerHour  int     `json:"base_production_per_hour"`
		ProductionMultiplier   float64 `json:"production_multiplier"`
	}

	types := []buildingTypeRow{}
	for rows.Next() {
		var bt buildingTypeRow
		err := rows.Scan(
			&bt.Name, &bt.DisplayName, &bt.Category,
			&bt.BaseCostMetal, &bt.BaseCostHe3, &bt.BaseCostGold,
			&bt.BaseTimeSeconds, &bt.CostMultiplier, &bt.TimeMultiplier,
			&bt.MaxLevel, &bt.MaxCountPerPlanet, &bt.CivicCenterReqPerLevel,
			&bt.BaseProductionPerHour, &bt.ProductionMultiplier,
		)
		if err != nil {
			log.Printf("Failed to scan building type: %v", err)
			continue
		}
		types = append(types, bt)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types)
}

// CancelUpgrade handles POST /api/planets/{id}/buildings/{buildingId}/cancel
func CancelUpgrade(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	planetID := r.PathValue("id")
	buildingID := r.PathValue("buildingId")

	if !verifyPlanetOwnership(planetID, playerID) {
		http.Error(w, `{"error":"planet not found"}`, http.StatusNotFound)
		return
	}

	// Get the building and verify it's upgrading
	var b models.Building
	err := database.DB.QueryRow(
		`SELECT id, planet_id, building_type, grid_col, grid_row, level, is_upgrading,
		        upgrade_finish_at, created_at, updated_at
		 FROM buildings
		 WHERE id = $1 AND planet_id = $2`,
		buildingID, planetID,
	).Scan(
		&b.ID, &b.PlanetID, &b.BuildingType, &b.GridCol, &b.GridRow, &b.Level, &b.IsUpgrading,
		&b.UpgradeFinishAt, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		http.Error(w, `{"error":"building not found"}`, http.StatusNotFound)
		return
	}

	if !b.IsUpgrading {
		http.Error(w, `{"error":"building is not upgrading"}`, http.StatusConflict)
		return
	}

	// Cancel the upgrade (no resource refund)
	err = database.DB.QueryRow(
		`UPDATE buildings
		 SET is_upgrading = false, upgrade_finish_at = NULL, updated_at = now()
		 WHERE id = $1
		 RETURNING id, planet_id, building_type, grid_col, grid_row, level, is_upgrading, upgrade_finish_at, created_at, updated_at`,
		buildingID,
	).Scan(
		&b.ID, &b.PlanetID, &b.BuildingType, &b.GridCol, &b.GridRow, &b.Level, &b.IsUpgrading,
		&b.UpgradeFinishAt, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to cancel upgrade: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Building models.Building `json:"building"`
	}{Building: b})
}

type moveRequest struct {
	GridCol int `json:"grid_col"`
	GridRow int `json:"grid_row"`
}

// MoveBuilding handles PUT /api/planets/{id}/buildings/{buildingId}/move
func MoveBuilding(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	planetID := r.PathValue("id")
	buildingID := r.PathValue("buildingId")

	if !verifyPlanetOwnership(planetID, playerID) {
		http.Error(w, `{"error":"planet not found"}`, http.StatusNotFound)
		return
	}

	var req moveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.GridCol < 0 || req.GridRow < 0 {
		http.Error(w, `{"error":"invalid grid coordinates"}`, http.StatusBadRequest)
		return
	}

	var b models.Building
	err := database.DB.QueryRow(
		`UPDATE buildings
		 SET grid_col = $1, grid_row = $2, updated_at = now()
		 WHERE id = $3 AND planet_id = $4
		 RETURNING id, planet_id, building_type, grid_col, grid_row, level, is_upgrading, upgrade_finish_at, created_at, updated_at`,
		req.GridCol, req.GridRow, buildingID, planetID,
	).Scan(
		&b.ID, &b.PlanetID, &b.BuildingType, &b.GridCol, &b.GridRow, &b.Level, &b.IsUpgrading,
		&b.UpgradeFinishAt, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"building not found"}`, http.StatusNotFound)
			return
		}
		log.Printf("Failed to move building: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Building models.Building `json:"building"`
	}{Building: b})
}

func verifyPlanetOwnership(planetID, playerID string) bool {
	var exists bool
	err := database.DB.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM planets WHERE id = $1 AND player_id = $2)`,
		planetID, playerID,
	).Scan(&exists)
	return err == nil && exists
}

// checkAndDeductResources checks if resources are sufficient and deducts them, returning detailed error if insufficient
func checkAndDeductResources(tx *sql.Tx, planetID string, metalCost, he3Cost, goldCost int64, actionName string) (metal, he3, gold int64, err error) {
	// Get current resources
	var currentMetal, currentHe3, currentGold int64
	err = tx.QueryRow(`SELECT metal, he3, gold FROM resources WHERE planet_id = $1`, planetID).
		Scan(&currentMetal, &currentHe3, &currentGold)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get resources: %w", err)
	}

	// Check sufficiency
	if currentMetal < metalCost || currentHe3 < he3Cost || currentGold < goldCost {
		return currentMetal, currentHe3, currentGold, fmt.Errorf("insufficient resources for %s", actionName)
	}

	// Deduct
	err = tx.QueryRow(
		`UPDATE resources
		 SET metal = metal - $1, he3 = he3 - $2, gold = gold - $3, updated_at = now()
		 WHERE planet_id = $4
		 RETURNING metal, he3, gold`,
		metalCost, he3Cost, goldCost, planetID,
	).Scan(&metal, &he3, &gold)

	return metal, he3, gold, err
}
