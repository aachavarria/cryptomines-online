package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/models"
)

// ListHullTypes handles GET /api/hull-types
func ListHullTypes(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(
		`SELECT id, name, display_name, hull_class, tier, armor_type,
		        base_shield, base_structure, base_stability, base_defense,
		        installation_slots, base_agility, base_movement, base_storage,
		        base_metal_cost, base_he3_cost, base_gold_cost,
		        base_build_time_seconds, description
		 FROM hull_types ORDER BY hull_class, tier, name`,
	)
	if err != nil {
		log.Printf("Failed to list hull types: %v", err)
		http.Error(w, `{"error":"failed to list hull types"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	hulls := []models.HullType{}
	for rows.Next() {
		var h models.HullType
		err := rows.Scan(
			&h.ID, &h.Name, &h.DisplayName, &h.HullClass, &h.Tier, &h.ArmorType,
			&h.BaseShield, &h.BaseStructure, &h.BaseStability, &h.BaseDefense,
			&h.InstallationSlots, &h.BaseAgility, &h.BaseMovement, &h.BaseStorage,
			&h.BaseMetalCost, &h.BaseHe3Cost, &h.BaseGoldCost,
			&h.BaseBuildTime, &h.Description,
		)
		if err != nil {
			log.Printf("Failed to scan hull type: %v", err)
			continue
		}
		hulls = append(hulls, h)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hulls)
}

// ListModuleTypes handles GET /api/module-types
// Supports ?category=ballistic filter
func ListModuleTypes(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	baseQuery := `SELECT id, name, display_name, category, tier, damage_type,
	                     min_damage, max_damage, weapon_range_min, weapon_range_max,
	                     cooldown, he3_per_round, volume, max_per_ship, effects_json,
	                     metal_cost, he3_cost, gold_cost, build_time_seconds, description
	              FROM module_types`

	var rows *sql.Rows
	var err error
	if category != "" {
		rows, err = database.DB.Query(baseQuery+` WHERE category = $1 ORDER BY name, tier`, category)
	} else {
		rows, err = database.DB.Query(baseQuery + ` ORDER BY category, name, tier`)
	}
	if err != nil {
		log.Printf("Failed to list module types: %v", err)
		http.Error(w, `{"error":"failed to list module types"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	modules := []models.ModuleType{}
	for rows.Next() {
		var m models.ModuleType
		err := rows.Scan(
			&m.ID, &m.Name, &m.DisplayName, &m.Category, &m.Tier, &m.DamageType,
			&m.MinDamage, &m.MaxDamage, &m.WeaponRangeMin, &m.WeaponRangeMax,
			&m.Cooldown, &m.He3PerRound, &m.Volume, &m.MaxPerShip, &m.EffectsJSON,
			&m.MetalCost, &m.He3Cost, &m.GoldCost, &m.BuildTime, &m.Description,
		)
		if err != nil {
			log.Printf("Failed to scan module type: %v", err)
			continue
		}
		modules = append(modules, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modules)
}
