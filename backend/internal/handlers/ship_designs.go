package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
	"github.com/cryptomines-online/backend/internal/services"
)

const maxDesigns = 20
const maxDesignNameLen = 20

var designNameRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

type createDesignRequest struct {
	Name       string                `json:"name"`
	HullTypeID int                   `json:"hull_type_id"`
	Modules    []models.DesignModule `json:"modules"`
}

// ListShipDesigns handles GET /api/ship-designs
func ListShipDesigns(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	rows, err := database.DB.Query(
		`SELECT id, player_id, name, hull_type_id, modules_json,
		        total_shield, total_structure, total_defense, total_agility,
		        total_movement, total_storage, attack_power,
		        weapon_range_min, weapon_range_max, volume_used,
		        he3_per_round, metal_cost, he3_cost, gold_cost,
		        build_time_seconds, created_at, updated_at
		 FROM ship_designs WHERE player_id = $1
		 ORDER BY created_at DESC`, playerID,
	)
	if err != nil {
		log.Printf("Failed to list ship designs: %v", err)
		http.Error(w, `{"error":"failed to list ship designs"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	designs := []models.ShipDesign{}
	for rows.Next() {
		var d models.ShipDesign
		err := rows.Scan(
			&d.ID, &d.PlayerID, &d.Name, &d.HullTypeID, &d.ModulesJSON,
			&d.TotalShield, &d.TotalStructure, &d.TotalDefense, &d.TotalAgility,
			&d.TotalMovement, &d.TotalStorage, &d.AttackPower,
			&d.WeaponRangeMin, &d.WeaponRangeMax, &d.VolumeUsed,
			&d.He3PerRound, &d.MetalCost, &d.He3Cost, &d.GoldCost,
			&d.BuildTimeSeconds, &d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan ship design: %v", err)
			continue
		}
		designs = append(designs, d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(designs)
}

// CreateShipDesign handles POST /api/ship-designs
func CreateShipDesign(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	var req createDesignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate name
	if len(req.Name) == 0 || len(req.Name) > maxDesignNameLen {
		http.Error(w, `{"error":"name must be 1-20 characters"}`, http.StatusBadRequest)
		return
	}
	if !designNameRegex.MatchString(req.Name) {
		http.Error(w, `{"error":"name must match ^[a-zA-Z0-9._-]+$"}`, http.StatusBadRequest)
		return
	}

	// Check design count
	var designCount int
	err := database.DB.QueryRow(
		`SELECT COUNT(*) FROM ship_designs WHERE player_id = $1`, playerID,
	).Scan(&designCount)
	if err != nil {
		log.Printf("Failed to count designs: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if designCount >= maxDesigns {
		http.Error(w, `{"error":"maximum 20 designs reached"}`, http.StatusConflict)
		return
	}

	// Get hull type
	var hull models.HullType
	err = database.DB.QueryRow(
		`SELECT id, name, display_name, hull_class, tier, armor_type,
		        base_shield, base_structure, base_stability, base_defense,
		        installation_slots, base_agility, base_movement, base_storage,
		        base_metal_cost, base_he3_cost, base_gold_cost,
		        base_build_time_seconds, description
		 FROM hull_types WHERE id = $1`, req.HullTypeID,
	).Scan(
		&hull.ID, &hull.Name, &hull.DisplayName, &hull.HullClass, &hull.Tier, &hull.ArmorType,
		&hull.BaseShield, &hull.BaseStructure, &hull.BaseStability, &hull.BaseDefense,
		&hull.InstallationSlots, &hull.BaseAgility, &hull.BaseMovement, &hull.BaseStorage,
		&hull.BaseMetalCost, &hull.BaseHe3Cost, &hull.BaseGoldCost,
		&hull.BaseBuildTime, &hull.Description,
	)
	if err != nil {
		http.Error(w, `{"error":"hull type not found"}`, http.StatusBadRequest)
		return
	}

	// Check player owns activated hull blueprint
	// Hull blueprints reference the Tier I hull; any hull in the same line requires only the line blueprint
	// Find the base hull name (strip _i, _ii, _iii suffix to find the line name)
	var hullBlueprintExists bool
	err = database.DB.QueryRow(
		`SELECT EXISTS(
		     SELECT 1 FROM player_blueprints pb
		     JOIN blueprints bp ON pb.blueprint_id = bp.id
		     JOIN hull_types ht ON bp.hull_type_id = ht.id
		     WHERE pb.player_id = $1 AND pb.is_activated = true
		       AND bp.blueprint_type = 'hull'
		       AND ht.name = (
		           SELECT ht2.name FROM hull_types ht2
		           WHERE ht2.hull_class = $2
		             AND ht2.tier = 1
		             AND REPLACE(REPLACE(REPLACE($3, '_iii', ''), '_ii', ''), '_i', '') =
		                 REPLACE(ht2.name, '_i', '')
		       )
		 )`, playerID, hull.HullClass, hull.Name,
	).Scan(&hullBlueprintExists)
	if err != nil {
		log.Printf("Failed to check hull blueprint: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if !hullBlueprintExists {
		http.Error(w, `{"error":"hull blueprint not owned or not activated"}`, http.StatusConflict)
		return
	}

	// Check hull tier is unlocked via blueprint research
	canUseHull, err := services.CanUseHullTier(playerID, req.HullTypeID)
	if err != nil {
		log.Printf("Failed to check hull tier: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if !canUseHull {
		http.Error(w, `{"error":"hull tier not unlocked - research blueprint to unlock"}`, http.StatusForbidden)
		return
	}

	// Validate modules
	if len(req.Modules) == 0 {
		http.Error(w, `{"error":"design must have at least one module"}`, http.StatusBadRequest)
		return
	}

	// Check each module blueprint is owned and activated, check per-ship limits
	moduleQuantities := map[string]int{} // module_name -> total quantity
	for _, mod := range req.Modules {
		var mt models.ModuleType
		err = database.DB.QueryRow(
			`SELECT id, name, max_per_ship FROM module_types WHERE id = $1`, mod.ModuleTypeID,
		).Scan(&mt.ID, &mt.Name, &mt.MaxPerShip)
		if err != nil {
			http.Error(w, `{"error":"module type not found"}`, http.StatusBadRequest)
			return
		}

		// Check module blueprint owned and activated
		var moduleBlueprintExists bool
		err = database.DB.QueryRow(
			`SELECT EXISTS(
			     SELECT 1 FROM player_blueprints pb
			     JOIN blueprints bp ON pb.blueprint_id = bp.id
			     JOIN module_types mtt ON bp.module_type_id = mtt.id
			     WHERE pb.player_id = $1 AND pb.is_activated = true
			       AND bp.blueprint_type = 'module'
			       AND mtt.name = $2
			 )`, playerID, mt.Name,
		).Scan(&moduleBlueprintExists)
		if err != nil {
			log.Printf("Failed to check module blueprint: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
		if !moduleBlueprintExists {
			http.Error(w, `{"error":"module blueprint not owned or not activated: `+mt.Name+`"}`, http.StatusConflict)
			return
		}

		// Check module tier is unlocked via blueprint research
		canUseModule, err := services.CanUseModuleTier(playerID, mod.ModuleTypeID)
		if err != nil {
			log.Printf("Failed to check module tier: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
		if !canUseModule {
			http.Error(w, `{"error":"module tier not unlocked - research blueprint to unlock"}`, http.StatusForbidden)
			return
		}

		// Check per-ship limit (0 = unlimited)
		moduleQuantities[mt.Name] += mod.Quantity
		if mt.MaxPerShip > 0 && moduleQuantities[mt.Name] > mt.MaxPerShip {
			http.Error(w, `{"error":"exceeded max per ship limit for module: `+mt.Name+`"}`, http.StatusConflict)
			return
		}
	}

	// Calculate stats
	stats, err := services.CalculateDesignStats(database.DB, hull, req.Modules)
	if err != nil {
		log.Printf("Failed to calculate design stats: %v", err)
		http.Error(w, `{"error":"failed to calculate design stats"}`, http.StatusInternalServerError)
		return
	}

	// Check volume fits
	if stats.VolumeUsed > hull.InstallationSlots {
		http.Error(w, `{"error":"total volume exceeds hull installation slots"}`, http.StatusConflict)
		return
	}

	// Serialize modules to JSON
	modulesJSON, err := json.Marshal(req.Modules)
	if err != nil {
		log.Printf("Failed to marshal modules: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Insert design
	var design models.ShipDesign
	err = database.DB.QueryRow(
		`INSERT INTO ship_designs (
		     player_id, name, hull_type_id, modules_json,
		     total_shield, total_structure, total_defense, total_agility,
		     total_movement, total_storage, attack_power,
		     weapon_range_min, weapon_range_max, volume_used,
		     he3_per_round, metal_cost, he3_cost, gold_cost, build_time_seconds
		 ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		 RETURNING id, player_id, name, hull_type_id, modules_json,
		           total_shield, total_structure, total_defense, total_agility,
		           total_movement, total_storage, attack_power,
		           weapon_range_min, weapon_range_max, volume_used,
		           he3_per_round, metal_cost, he3_cost, gold_cost,
		           build_time_seconds, created_at, updated_at`,
		playerID, req.Name, req.HullTypeID, string(modulesJSON),
		stats.TotalShield, stats.TotalStructure, stats.TotalDefense, stats.TotalAgility,
		stats.TotalMovement, stats.TotalStorage, stats.AttackPower,
		stats.WeaponRangeMin, stats.WeaponRangeMax, stats.VolumeUsed,
		stats.He3PerRound, stats.MetalCost, stats.He3Cost, stats.GoldCost, stats.BuildTimeSeconds,
	).Scan(
		&design.ID, &design.PlayerID, &design.Name, &design.HullTypeID, &design.ModulesJSON,
		&design.TotalShield, &design.TotalStructure, &design.TotalDefense, &design.TotalAgility,
		&design.TotalMovement, &design.TotalStorage, &design.AttackPower,
		&design.WeaponRangeMin, &design.WeaponRangeMax, &design.VolumeUsed,
		&design.He3PerRound, &design.MetalCost, &design.He3Cost, &design.GoldCost,
		&design.BuildTimeSeconds, &design.CreatedAt, &design.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to insert ship design: %v", err)
		http.Error(w, `{"error":"failed to create ship design"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		Design models.ShipDesign `json:"design"`
	}{Design: design})
}

// UpdateShipDesign handles PUT /api/ship-designs/{id}
// Only allowed if no ships have been built with this design.
func UpdateShipDesign(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	designID := r.PathValue("id")

	// Verify ownership
	var existingDesign models.ShipDesign
	err := database.DB.QueryRow(
		`SELECT id, player_id FROM ship_designs WHERE id = $1 AND player_id = $2`,
		designID, playerID,
	).Scan(&existingDesign.ID, &existingDesign.PlayerID)
	if err != nil {
		http.Error(w, `{"error":"design not found"}`, http.StatusNotFound)
		return
	}

	// Check no ships built
	var shipCount int
	err = database.DB.QueryRow(
		`SELECT COALESCE(SUM(quantity + build_quantity), 0) FROM ships WHERE ship_design_id = $1`,
		designID,
	).Scan(&shipCount)
	if err != nil {
		log.Printf("Failed to check ships: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if shipCount > 0 {
		http.Error(w, `{"error":"cannot modify design with existing ships"}`, http.StatusConflict)
		return
	}

	var req createDesignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate name
	if len(req.Name) == 0 || len(req.Name) > maxDesignNameLen {
		http.Error(w, `{"error":"name must be 1-20 characters"}`, http.StatusBadRequest)
		return
	}
	if !designNameRegex.MatchString(req.Name) {
		http.Error(w, `{"error":"name must match ^[a-zA-Z0-9._-]+$"}`, http.StatusBadRequest)
		return
	}

	// Get hull type
	var hull models.HullType
	err = database.DB.QueryRow(
		`SELECT id, name, display_name, hull_class, tier, armor_type,
		        base_shield, base_structure, base_stability, base_defense,
		        installation_slots, base_agility, base_movement, base_storage,
		        base_metal_cost, base_he3_cost, base_gold_cost,
		        base_build_time_seconds, description
		 FROM hull_types WHERE id = $1`, req.HullTypeID,
	).Scan(
		&hull.ID, &hull.Name, &hull.DisplayName, &hull.HullClass, &hull.Tier, &hull.ArmorType,
		&hull.BaseShield, &hull.BaseStructure, &hull.BaseStability, &hull.BaseDefense,
		&hull.InstallationSlots, &hull.BaseAgility, &hull.BaseMovement, &hull.BaseStorage,
		&hull.BaseMetalCost, &hull.BaseHe3Cost, &hull.BaseGoldCost,
		&hull.BaseBuildTime, &hull.Description,
	)
	if err != nil {
		http.Error(w, `{"error":"hull type not found"}`, http.StatusBadRequest)
		return
	}

	// Check hull tier is unlocked via blueprint research
	canUseHull, err := services.CanUseHullTier(playerID, req.HullTypeID)
	if err != nil {
		log.Printf("Failed to check hull tier: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if !canUseHull {
		http.Error(w, `{"error":"hull tier not unlocked - research blueprint to unlock"}`, http.StatusForbidden)
		return
	}

	// Validate modules
	if len(req.Modules) == 0 {
		http.Error(w, `{"error":"design must have at least one module"}`, http.StatusBadRequest)
		return
	}

	// Check each module tier is unlocked
	for _, mod := range req.Modules {
		canUseModule, err := services.CanUseModuleTier(playerID, mod.ModuleTypeID)
		if err != nil {
			log.Printf("Failed to check module tier: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
		if !canUseModule {
			http.Error(w, `{"error":"module tier not unlocked - research blueprint to unlock"}`, http.StatusForbidden)
			return
		}
	}

	// Calculate stats
	stats, err := services.CalculateDesignStats(database.DB, hull, req.Modules)
	if err != nil {
		log.Printf("Failed to calculate design stats: %v", err)
		http.Error(w, `{"error":"failed to calculate design stats"}`, http.StatusInternalServerError)
		return
	}

	if stats.VolumeUsed > hull.InstallationSlots {
		http.Error(w, `{"error":"total volume exceeds hull installation slots"}`, http.StatusConflict)
		return
	}

	modulesJSON, err := json.Marshal(req.Modules)
	if err != nil {
		log.Printf("Failed to marshal modules: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	var design models.ShipDesign
	err = database.DB.QueryRow(
		`UPDATE ship_designs SET
		     name = $1, hull_type_id = $2, modules_json = $3,
		     total_shield = $4, total_structure = $5, total_defense = $6, total_agility = $7,
		     total_movement = $8, total_storage = $9, attack_power = $10,
		     weapon_range_min = $11, weapon_range_max = $12, volume_used = $13,
		     he3_per_round = $14, metal_cost = $15, he3_cost = $16, gold_cost = $17,
		     build_time_seconds = $18, updated_at = now()
		 WHERE id = $19
		 RETURNING id, player_id, name, hull_type_id, modules_json,
		           total_shield, total_structure, total_defense, total_agility,
		           total_movement, total_storage, attack_power,
		           weapon_range_min, weapon_range_max, volume_used,
		           he3_per_round, metal_cost, he3_cost, gold_cost,
		           build_time_seconds, created_at, updated_at`,
		req.Name, req.HullTypeID, string(modulesJSON),
		stats.TotalShield, stats.TotalStructure, stats.TotalDefense, stats.TotalAgility,
		stats.TotalMovement, stats.TotalStorage, stats.AttackPower,
		stats.WeaponRangeMin, stats.WeaponRangeMax, stats.VolumeUsed,
		stats.He3PerRound, stats.MetalCost, stats.He3Cost, stats.GoldCost,
		stats.BuildTimeSeconds, designID,
	).Scan(
		&design.ID, &design.PlayerID, &design.Name, &design.HullTypeID, &design.ModulesJSON,
		&design.TotalShield, &design.TotalStructure, &design.TotalDefense, &design.TotalAgility,
		&design.TotalMovement, &design.TotalStorage, &design.AttackPower,
		&design.WeaponRangeMin, &design.WeaponRangeMax, &design.VolumeUsed,
		&design.He3PerRound, &design.MetalCost, &design.He3Cost, &design.GoldCost,
		&design.BuildTimeSeconds, &design.CreatedAt, &design.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to update ship design: %v", err)
		http.Error(w, `{"error":"failed to update ship design"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Design models.ShipDesign `json:"design"`
	}{Design: design})
}

// DeleteShipDesign handles DELETE /api/ship-designs/{id}
// Only allowed if no ships exist for this design.
func DeleteShipDesign(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	designID := r.PathValue("id")

	// Verify ownership
	var exists bool
	err := database.DB.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM ship_designs WHERE id = $1 AND player_id = $2)`,
		designID, playerID,
	).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, `{"error":"design not found"}`, http.StatusNotFound)
		return
	}

	// Check no ships exist
	var shipCount int
	err = database.DB.QueryRow(
		`SELECT COALESCE(SUM(quantity + build_quantity), 0) FROM ships WHERE ship_design_id = $1`,
		designID,
	).Scan(&shipCount)
	if err != nil {
		log.Printf("Failed to check ships: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if shipCount > 0 {
		http.Error(w, `{"error":"cannot delete design with existing ships"}`, http.StatusConflict)
		return
	}

	// Check no fleet stacks reference this design
	var stackCount int
	err = database.DB.QueryRow(
		`SELECT COUNT(*) FROM fleet_stacks WHERE ship_design_id = $1`, designID,
	).Scan(&stackCount)
	if err != nil {
		log.Printf("Failed to check fleet stacks: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if stackCount > 0 {
		http.Error(w, `{"error":"cannot delete design used in fleets"}`, http.StatusConflict)
		return
	}

	_, err = database.DB.Exec(`DELETE FROM ship_designs WHERE id = $1`, designID)
	if err != nil {
		log.Printf("Failed to delete ship design: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Deleted bool `json:"deleted"`
	}{Deleted: true})
}

// GetDesignStats handles GET /api/ship-designs/{id}/stats
// Recalculates and returns stats for a design.
func GetDesignStats(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	designID := r.PathValue("id")

	var design models.ShipDesign
	err := database.DB.QueryRow(
		`SELECT id, player_id, name, hull_type_id, modules_json,
		        total_shield, total_structure, total_defense, total_agility,
		        total_movement, total_storage, attack_power,
		        weapon_range_min, weapon_range_max, volume_used,
		        he3_per_round, metal_cost, he3_cost, gold_cost,
		        build_time_seconds, created_at, updated_at
		 FROM ship_designs WHERE id = $1 AND player_id = $2`,
		designID, playerID,
	).Scan(
		&design.ID, &design.PlayerID, &design.Name, &design.HullTypeID, &design.ModulesJSON,
		&design.TotalShield, &design.TotalStructure, &design.TotalDefense, &design.TotalAgility,
		&design.TotalMovement, &design.TotalStorage, &design.AttackPower,
		&design.WeaponRangeMin, &design.WeaponRangeMax, &design.VolumeUsed,
		&design.He3PerRound, &design.MetalCost, &design.He3Cost, &design.GoldCost,
		&design.BuildTimeSeconds, &design.CreatedAt, &design.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"design not found"}`, http.StatusNotFound)
			return
		}
		log.Printf("Failed to get design: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(design)
}
