package services

import (
	"database/sql"
	"log"
)

// LevelCost holds the cost and build time for a specific building level.
type LevelCost struct {
	MetalCost        int64
	He3Cost          int64
	GoldCost         int64
	BuildTimeSeconds int
}

// lookupTableMap maps building type names to their level reference table names.
var lookupTableMap = map[string]string{
	"metal_collector":         "metal_collector_levels",
	"he3_extractor":           "he3_extractor_levels",
	"residential_area":        "residential_area_levels",
	"resource_warehouse":      "resource_warehouse_levels",
	"civic_center":            "civic_center_levels",
	"technology_center":       "technology_center_levels",
	"command_center":          "command_center_levels",
	"space_station":           "space_station_levels",
	"weapon_research_center":  "weapon_research_center_levels",
	"alliance_center":         "alliance_center_levels",
	"trading_center":          "trading_center_levels",
	"radar":                   "radar_levels",
	"spacedock":               "spacedock_levels",
	"recycling_plant":         "recycling_plant_levels",
	"meteor_star":             "meteor_star_levels",
	"particle_cannon":         "particle_cannon_levels",
	"anti_aircraft_gun":       "anti_aircraft_gun_levels",
	"thors_cannon":            "thors_cannon_levels",
}

// GetBuildingLevelCost looks up the exact cost for a building at a given level.
// It first tries the level reference table (exact wiki data).
// If no lookup table exists or the level is not found, it falls back to the formula.
func GetBuildingLevelCost(db *sql.DB, buildingTypeName string, level int, baseMetal, baseHe3, baseGold int64, costMult float64, baseTime int, timeMult float64) LevelCost {
	tableName, hasLookup := lookupTableMap[buildingTypeName]
	if hasLookup {
		cost, err := lookupLevelCost(db, tableName, level)
		if err == nil {
			return cost
		}
		log.Printf("Lookup table miss for %s level %d, falling back to formula: %v", buildingTypeName, level, err)
	}

	// Fallback to formula
	metal, he3, gold := UpgradeCost(baseMetal, baseHe3, baseGold, costMult, level)
	buildTime := UpgradeTime(baseTime, timeMult, level)
	return LevelCost{
		MetalCost:        metal,
		He3Cost:          he3,
		GoldCost:         gold,
		BuildTimeSeconds: buildTime,
	}
}

func lookupLevelCost(db *sql.DB, tableName string, level int) (LevelCost, error) {
	// All level tables have metal_cost, he3_cost, gold_cost, build_time_seconds columns
	query := "SELECT metal_cost, he3_cost, gold_cost, build_time_seconds FROM " + tableName + " WHERE level = $1"
	var cost LevelCost
	err := db.QueryRow(query, level).Scan(&cost.MetalCost, &cost.He3Cost, &cost.GoldCost, &cost.BuildTimeSeconds)
	if err != nil {
		return LevelCost{}, err
	}
	return cost, nil
}
