package models

import "time"

type BuildingType struct {
	ID                     int     `json:"id"`
	Name                   string  `json:"name"`
	DisplayName            string  `json:"display_name"`
	Category               string  `json:"category"`
	Base                   string  `json:"base"`
	BaseCostMetal          int64   `json:"base_cost_metal"`
	BaseCostHe3            int64   `json:"base_cost_he3"`
	BaseCostGold           int64   `json:"base_cost_gold"`
	BaseTimeSeconds        int     `json:"base_time_seconds"`
	CostMultiplier         float64 `json:"cost_multiplier"`
	TimeMultiplier         float64 `json:"time_multiplier"`
	BaseProductionPerHour  int     `json:"base_production_per_hour"`
	ProductionMultiplier   float64 `json:"production_multiplier"`
	MaxLevel               int     `json:"max_level"`
	MaxCountPerPlanet      int     `json:"max_count_per_planet"`
	PrerequisiteBuilding   *string `json:"prerequisite_building"`
	PrerequisiteLevel      int     `json:"prerequisite_level"`
	CivicCenterReqPerLevel bool    `json:"civic_center_req_per_level"`
	Description            string  `json:"description"`
}

type Building struct {
	ID              string     `json:"id"`
	PlanetID        string     `json:"planet_id"`
	BuildingType    int        `json:"building_type"`
	Level           int        `json:"level"`
	IsUpgrading     bool       `json:"is_upgrading"`
	UpgradeFinishAt *time.Time `json:"upgrade_finish_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type BuildingWithType struct {
	Building
	TypeName    string `json:"type_name"`
	DisplayName string `json:"display_name"`
	Category    string `json:"category"`
	BaseName    string `json:"base_type"`
	MaxLevel    int    `json:"max_level"`
}
