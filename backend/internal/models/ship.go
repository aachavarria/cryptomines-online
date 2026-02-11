package models

import "time"

type HullType struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	DisplayName       string  `json:"display_name"`
	HullClass         string  `json:"hull_class"`
	Tier              int     `json:"tier"`
	ArmorType         string  `json:"armor_type"`
	BaseShield        int     `json:"base_shield"`
	BaseStructure     int     `json:"base_structure"`
	BaseStability     float64 `json:"base_stability"`
	BaseDefense       float64 `json:"base_defense"`
	InstallationSlots int     `json:"installation_slots"`
	BaseAgility       int     `json:"base_agility"`
	BaseMovement      int     `json:"base_movement"`
	BaseStorage       int     `json:"base_storage"`
	BaseMetalCost     int64   `json:"base_metal_cost"`
	BaseHe3Cost       int64   `json:"base_he3_cost"`
	BaseGoldCost      int64   `json:"base_gold_cost"`
	BaseBuildTime     int     `json:"base_build_time_seconds"`
	Description       string  `json:"description"`
}

type ModuleType struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	DisplayName    string `json:"display_name"`
	Category       string `json:"category"`
	Tier           int    `json:"tier"`
	DamageType     *string `json:"damage_type"`
	MinDamage      int    `json:"min_damage"`
	MaxDamage      int    `json:"max_damage"`
	WeaponRangeMin int    `json:"weapon_range_min"`
	WeaponRangeMax int    `json:"weapon_range_max"`
	Cooldown       int    `json:"cooldown"`
	He3PerRound    int    `json:"he3_per_round"`
	Volume         int    `json:"volume"`
	MaxPerShip     int    `json:"max_per_ship"`
	EffectsJSON    string `json:"effects_json"`
	MetalCost      int64  `json:"metal_cost"`
	He3Cost        int64  `json:"he3_cost"`
	GoldCost       int64  `json:"gold_cost"`
	BuildTime      int    `json:"build_time_seconds"`
	Description    string `json:"description"`
}

type Blueprint struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	DisplayName    string `json:"display_name"`
	BlueprintType  string `json:"blueprint_type"`
	HullTypeID     *int   `json:"hull_type_id"`
	ModuleTypeID   *int   `json:"module_type_id"`
	HullClass      string `json:"hull_class,omitempty"`
	ModuleCategory string `json:"module_category,omitempty"`
	Source         string `json:"source"`
	ResearchLevel  int    `json:"research_level"`
	Description    string `json:"description"`
}

type PlayerBlueprint struct {
	ID            string    `json:"id"`
	PlayerID      string    `json:"player_id"`
	BlueprintID   int       `json:"blueprint_id"`
	IsActivated   bool      `json:"is_activated"`
	ResearchLevel int       `json:"research_level"`
	AcquiredAt    time.Time `json:"acquired_at"`
}

type PlayerBlueprintWithInfo struct {
	PlayerBlueprint
	BlueprintName string  `json:"blueprint_name"`
	BlueprintType string  `json:"blueprint_type"`
	HullTypeID    *int    `json:"hull_type_id"`
	ModuleTypeID  *int    `json:"module_type_id"`
	Source        string  `json:"source"`
}

type ShipDesign struct {
	ID              string    `json:"id"`
	PlayerID        string    `json:"player_id"`
	Name            string    `json:"name"`
	HullTypeID      int       `json:"hull_type_id"`
	ModulesJSON     string    `json:"modules_json"`
	TotalShield     int       `json:"total_shield"`
	TotalStructure  int       `json:"total_structure"`
	TotalDefense    float64   `json:"total_defense"`
	TotalAgility    int       `json:"total_agility"`
	TotalMovement   int       `json:"total_movement"`
	TotalStorage    int       `json:"total_storage"`
	AttackPower     int       `json:"attack_power"`
	WeaponRangeMin  int       `json:"weapon_range_min"`
	WeaponRangeMax  int       `json:"weapon_range_max"`
	VolumeUsed      int       `json:"volume_used"`
	He3PerRound     int       `json:"he3_per_round"`
	MetalCost       int64     `json:"metal_cost"`
	He3Cost         int64     `json:"he3_cost"`
	GoldCost        int64     `json:"gold_cost"`
	BuildTimeSeconds int      `json:"build_time_seconds"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type DesignModule struct {
	ModuleTypeID   int `json:"module_type_id"`
	Quantity       int `json:"quantity"`
	PlacementOrder int `json:"placement_order"`
}

type Ship struct {
	ID            string     `json:"id"`
	PlayerID      string     `json:"player_id"`
	ShipDesignID  string     `json:"ship_design_id"`
	Quantity      int        `json:"quantity"`
	IsBuilding    bool       `json:"is_building"`
	BuildQuantity int        `json:"build_quantity"`
	BuildFinishAt *time.Time `json:"build_finish_at"`
	ProductionSlot int       `json:"production_slot"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type Fleet struct {
	ID               string     `json:"id"`
	PlayerID         string     `json:"player_id"`
	Name             string     `json:"name"`
	Formation        string     `json:"formation"`
	CommanderID      *string    `json:"commander_id"`
	TargetingCommand string     `json:"targeting_command"`
	Status           string     `json:"status"`
	PlanetID         *string    `json:"planet_id"`
	PositionX        *int       `json:"position_x"`
	PositionY        *int       `json:"position_y"`
	DestinationX     *int       `json:"destination_x"`
	DestinationY     *int       `json:"destination_y"`
	ArrivalAt        *time.Time `json:"arrival_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type FleetStack struct {
	ID           string `json:"id"`
	FleetID      string `json:"fleet_id"`
	ShipDesignID string `json:"ship_design_id"`
	GridRow      int    `json:"grid_row"`
	GridCol      int    `json:"grid_col"`
	ShipCount    int    `json:"ship_count"`
}

type Instance struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Type                string `json:"type"`
	Difficulty          int    `json:"difficulty"`
	RequiredLevel       int    `json:"required_level"`
	MaxFleets           int    `json:"max_fleets"`
	ShipsLostOnDefeat   bool   `json:"ships_lost_on_defeat"`
	He3LostOnDefeat     bool   `json:"he3_lost_on_defeat"`
	ExpReward           int    `json:"exp_reward"`
	EnemyFleetsJSON     string `json:"enemy_fleets_json"`
	RewardsJSON         string `json:"rewards_json"`
	Description         string `json:"description"`
}

type InstanceProgress struct {
	ID            string     `json:"id"`
	PlayerID      string     `json:"player_id"`
	InstanceID    int        `json:"instance_id"`
	Completed     bool       `json:"completed"`
	Attempts      int        `json:"attempts"`
	BestScore     int        `json:"best_score"`
	LastAttemptAt *time.Time `json:"last_attempt_at"`
}

type SpacedockRepair struct {
	ID              string     `json:"id"`
	PlayerID        string     `json:"player_id"`
	ShipDesignID    string     `json:"ship_design_id"`
	DestroyedCount  int        `json:"destroyed_count"`
	RepairedCount   int        `json:"repaired_count"`
	RepairFinishAt  *time.Time `json:"repair_finish_at"`
	CombatReportID  *string    `json:"combat_report_id"`
	CreatedAt       time.Time  `json:"created_at"`
}

type ShipFactoryLevel struct {
	Level          int   `json:"level"`
	CivicCenterReq int   `json:"civic_center_req"`
	MetalCost      int64 `json:"metal_cost"`
	He3Cost        int64 `json:"he3_cost"`
	GoldCost       int64 `json:"gold_cost"`
	BuildTime      int   `json:"build_time_seconds"`
	SpeedBonusPct  int   `json:"speed_bonus_pct"`
	ProductionSlots int  `json:"production_slots"`
}
