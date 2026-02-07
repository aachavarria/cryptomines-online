package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
	"github.com/cryptomines-online/backend/internal/services"
)

// GetResources handles GET /api/planets/{id}/resources
func GetResources(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	planetID := r.PathValue("id")

	if !verifyPlanetOwnership(planetID, playerID) {
		http.Error(w, `{"error":"planet not found"}`, http.StatusNotFound)
		return
	}

	// First, check for any completed building upgrades and apply them
	applyCompletedUpgrades(planetID)

	var res models.Resource
	err := database.DB.QueryRow(
		`SELECT id, planet_id, metal, he3, gold, metal_per_hour, he3_per_hour,
		        gold_per_hour, storage_capacity, last_collected_at, updated_at
		 FROM resources WHERE planet_id = $1`, planetID,
	).Scan(
		&res.ID, &res.PlanetID, &res.Metal, &res.He3, &res.Gold,
		&res.MetalPerHour, &res.He3PerHour, &res.GoldPerHour,
		&res.StorageCapacity, &res.LastCollectedAt, &res.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to get resources: %v", err)
		http.Error(w, `{"error":"resources not found"}`, http.StatusNotFound)
		return
	}

	// Calculate pending accumulated resources (not yet collected)
	elapsed := time.Since(res.LastCollectedAt).Hours()
	pendingMetal := int64(float64(res.MetalPerHour) * elapsed)
	pendingHe3 := int64(float64(res.He3PerHour) * elapsed)
	pendingGold := int64(float64(res.GoldPerHour) * elapsed)

	type resourcesResponse struct {
		models.Resource
		PendingMetal int64 `json:"pending_metal"`
		PendingHe3   int64 `json:"pending_he3"`
		PendingGold  int64 `json:"pending_gold"`
	}

	resp := resourcesResponse{
		Resource:     res,
		PendingMetal: pendingMetal,
		PendingHe3:   pendingHe3,
		PendingGold:  pendingGold,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type collectResponse struct {
	Collected collectAmounts  `json:"collected"`
	Resources models.Resource `json:"resources"`
}

type collectAmounts struct {
	Metal int64 `json:"metal"`
	He3   int64 `json:"he3"`
	Gold  int64 `json:"gold"`
}

// CollectResources handles POST /api/planets/{id}/resources/collect
func CollectResources(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	planetID := r.PathValue("id")

	if !verifyPlanetOwnership(planetID, playerID) {
		http.Error(w, `{"error":"planet not found"}`, http.StatusNotFound)
		return
	}

	// Apply completed upgrades first
	applyCompletedUpgrades(planetID)

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Get current resources with lock
	var res models.Resource
	err = tx.QueryRow(
		`SELECT id, planet_id, metal, he3, gold, metal_per_hour, he3_per_hour,
		        gold_per_hour, storage_capacity, last_collected_at, updated_at
		 FROM resources WHERE planet_id = $1 FOR UPDATE`, planetID,
	).Scan(
		&res.ID, &res.PlanetID, &res.Metal, &res.He3, &res.Gold,
		&res.MetalPerHour, &res.He3PerHour, &res.GoldPerHour,
		&res.StorageCapacity, &res.LastCollectedAt, &res.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to get resources: %v", err)
		http.Error(w, `{"error":"resources not found"}`, http.StatusNotFound)
		return
	}

	// Calculate accumulated resources
	elapsed := time.Since(res.LastCollectedAt).Hours()
	collectedMetal := int64(float64(res.MetalPerHour) * elapsed)
	collectedHe3 := int64(float64(res.He3PerHour) * elapsed)
	collectedGold := int64(float64(res.GoldPerHour) * elapsed)

	// Cap resources at storage capacity
	newMetal := min64(res.Metal+collectedMetal, res.StorageCapacity)
	newHe3 := min64(res.He3+collectedHe3, res.StorageCapacity)
	newGold := min64(res.Gold+collectedGold, res.StorageCapacity)

	actualMetal := newMetal - res.Metal
	actualHe3 := newHe3 - res.He3
	actualGold := newGold - res.Gold

	// Update resources
	now := time.Now()
	err = tx.QueryRow(
		`UPDATE resources
		 SET metal = $1, he3 = $2, gold = $3, last_collected_at = $4, updated_at = $4
		 WHERE planet_id = $5
		 RETURNING id, planet_id, metal, he3, gold, metal_per_hour, he3_per_hour,
		           gold_per_hour, storage_capacity, last_collected_at, updated_at`,
		newMetal, newHe3, newGold, now, planetID,
	).Scan(
		&res.ID, &res.PlanetID, &res.Metal, &res.He3, &res.Gold,
		&res.MetalPerHour, &res.He3PerHour, &res.GoldPerHour,
		&res.StorageCapacity, &res.LastCollectedAt, &res.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to update resources: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Update quest progress for harvesting resources
	services.UpdateQuestProgress(playerID, "harvest_resources", "resource_warehouse", 1)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collectResponse{
		Collected: collectAmounts{Metal: actualMetal, He3: actualHe3, Gold: actualGold},
		Resources: res,
	})
}

// applyCompletedUpgrades checks for buildings that have finished upgrading,
// completes their upgrade, and recalculates production rates.
func applyCompletedUpgrades(planetID string) {
	now := time.Now()

	// Find and complete all finished upgrades
	rows, err := database.DB.Query(`
		SELECT b.id, bt.name, b.level
		FROM buildings b
		JOIN building_types bt ON b.building_type = bt.id
		WHERE b.planet_id = $1 AND b.is_upgrading = true AND b.upgrade_finish_at <= $2
	`, planetID, now)
	if err != nil {
		log.Printf("Failed to query completed upgrades: %v", err)
		return
	}

	type completedBuilding struct {
		id          string
		buildingType string
		newLevel     int
	}
	var completed []completedBuilding

	for rows.Next() {
		var id, buildingType string
		var currentLevel int
		rows.Scan(&id, &buildingType, &currentLevel)
		completed = append(completed, completedBuilding{
			id:          id,
			buildingType: buildingType,
			newLevel:     currentLevel + 1,
		})
	}
	rows.Close()

	if len(completed) == 0 {
		return
	}

	// Get player ID for quest tracking
	var playerID string
	err = database.DB.QueryRow(`
		SELECT player_id FROM planets WHERE id = $1
	`, planetID).Scan(&playerID)
	if err != nil {
		log.Printf("Failed to get player ID for quest tracking: %v", err)
		return
	}

	// Apply the upgrades
	_, err = database.DB.Exec(`
		UPDATE buildings
		SET level = level + 1, is_upgrading = false, upgrade_finish_at = NULL, updated_at = now()
		WHERE planet_id = $1 AND is_upgrading = true AND upgrade_finish_at <= $2
	`, planetID, now)
	if err != nil {
		log.Printf("Failed to apply completed upgrades: %v", err)
		return
	}

	// Update quest progress for each completed building
	for _, building := range completed {
		// Track quest progress for building construction (level 1)
		if building.newLevel == 1 {
			services.UpdateQuestProgress(playerID, "build_building", building.buildingType, 1)
		}
		// Track quest progress for building upgrades (level > 1)
		if building.newLevel > 1 {
			services.UpdateQuestProgress(playerID, "upgrade_building", building.buildingType, 1)
		}
	}

	recalculateProductionRates(planetID)
}

// recalculateProductionRates recalculates metal_per_hour, he3_per_hour, gold_per_hour,
// and storage_capacity based on current buildings on the planet.
func recalculateProductionRates(planetID string) {
	var metalPerHour, he3PerHour, goldPerHour int64
	var storageCapacity int64

	// Sum up production from all resource-producing buildings
	rows, err := database.DB.Query(
		`SELECT bt.name, b.level, bt.base_production_per_hour, bt.production_multiplier
		 FROM buildings b
		 JOIN building_types bt ON b.building_type = bt.id
		 WHERE b.planet_id = $1 AND b.is_upgrading = false`, planetID,
	)
	if err != nil {
		log.Printf("Failed to query buildings for production: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var level int
		var baseProd int
		var prodMult float64
		if err := rows.Scan(&name, &level, &baseProd, &prodMult); err != nil {
			continue
		}

		switch name {
		case "metal_collector":
			metalPerHour += lookupProductionFromDB("metal_collector_levels", "metal_output_per_hour", level)
		case "he3_extractor":
			he3PerHour += lookupProductionFromDB("he3_extractor_levels", "he3_output_per_hour", level)
		case "residential_area":
			goldPerHour += lookupProductionFromDB("residential_area_levels", "gold_output_per_hour", level)
		case "resource_warehouse":
			cap := lookupStorageFromDB(level)
			if cap > storageCapacity {
				storageCapacity = cap
			}
		}
	}

	// Use a reasonable default if no warehouse exists
	if storageCapacity == 0 {
		storageCapacity = 100000
	}

	_, err = database.DB.Exec(
		`UPDATE resources
		 SET metal_per_hour = $1, he3_per_hour = $2, gold_per_hour = $3,
		     storage_capacity = $4, updated_at = now()
		 WHERE planet_id = $5`,
		metalPerHour, he3PerHour, goldPerHour, storageCapacity, planetID,
	)
	if err != nil {
		log.Printf("Failed to update production rates: %v", err)
	}
}

func lookupProductionFromDB(table, column string, level int) int64 {
	var output int64
	// Using parameterized level but table/column names are hardcoded from switch statement
	query := "SELECT " + column + " FROM " + table + " WHERE level = $1"
	err := database.DB.QueryRow(query, level).Scan(&output)
	if err != nil {
		return 0
	}
	return output
}

func lookupStorageFromDB(level int) int64 {
	var capacity int64
	err := database.DB.QueryRow(
		`SELECT storage_capacity FROM resource_warehouse_levels WHERE level = $1`, level,
	).Scan(&capacity)
	if err != nil {
		return 0
	}
	return capacity
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
