package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cryptomines-online/backend/internal/database"
)

// TechBonuses represents all active technology bonuses for a player.
// These bonuses are calculated from completed research levels.
type TechBonuses struct {
	// Production bonuses (percent)
	MetalOutput  float64 `json:"metal_output"`
	He3Output    float64 `json:"he3_output"`
	GoldOutput   float64 `json:"gold_output"`

	// Construction bonuses
	BuildSpeed           float64 `json:"build_speed"`             // percent bonus to building construction speed
	BuildCostReduction   float64 `json:"build_cost_reduction"`    // percent reduction in building costs
	ConstructionSlots    int     `json:"construction_slots"`      // additional construction slots
	ShipBuildSpeed       float64 `json:"ship_build_speed"`        // percent bonus to ship construction speed
	ShipBuildCostReduction float64 `json:"ship_build_cost_reduction"` // percent reduction in ship costs
	ShipProductionSlots  int     `json:"ship_production_slots"`   // additional shipbuilding slots

	// Capacity bonuses
	WarehouseCapacity int64 `json:"warehouse_capacity"` // flat bonus to warehouse storage

	// Ship repair
	ShipRepairPercent float64 `json:"ship_repair_percent"` // percent of max HP repaired

	// Combat bonuses - Ballistics
	BallisticDamage      float64 `json:"ballistic_damage"`
	BallisticCritRate    float64 `json:"ballistic_crit_rate"`
	BallisticCritDamage  float64 `json:"ballistic_crit_damage"`
	BallisticHitRate     float64 `json:"ballistic_hit_rate"`
	WeaponSpaceReduction float64 `json:"weapon_space_reduction"`
	ShieldBypass         float64 `json:"shield_bypass"`

	// Combat bonuses - Directional
	DirectionalDamage   float64 `json:"directional_damage"`
	DirectionalCritRate float64 `json:"directional_crit_rate"`
	DirectionalAccuracy float64 `json:"directional_accuracy"`
	SteeringPower       float64 `json:"steering_power"`

	// Combat bonuses - Missile
	MissileDamage     float64 `json:"missile_damage"`
	MissileHitRate    float64 `json:"missile_hit_rate"`
	InterceptReduction float64 `json:"intercept_reduction"`

	// Combat bonuses - Fighter
	FighterDamage   float64 `json:"fighter_damage"`
	FighterHitRate  float64 `json:"fighter_hit_rate"`
	FuelOptimization float64 `json:"fuel_optimization"` // He3 cost reduction

	// Defense bonuses
	BaseShield     float64 `json:"base_shield"`     // percent bonus
	BaseStructure  float64 `json:"base_structure"`  // percent bonus
	BaseAgility    float64 `json:"base_agility"`    // percent bonus
	BaseDefense    float64 `json:"base_defense"`    // percent bonus
	BaseStability  float64 `json:"base_stability"`  // percent bonus

	// Planetary defense bonuses
	DefenseCostReduction   float64 `json:"defense_cost_reduction"`
	DefenseBuildSpeed      float64 `json:"defense_build_speed"`
	DefenseValue           float64 `json:"defense_value"`
	EmplacementAttack      float64 `json:"emplacement_attack"`
	MaxDefenseStructures   float64 `json:"max_defense_structures"`
	MaxThorCannon          int     `json:"max_thor_cannon"`
}

// techEffectData represents the parsed effects_json structure from tech_types.
type techEffectData struct {
	Type     string  `json:"type"`
	PerLevel float64 `json:"per_level"`
	Unit     string  `json:"unit"`
	Flat     float64 `json:"flat"`
}

// GetPlayerTechBonuses calculates and returns all active technology bonuses for a player.
// It queries all completed research and aggregates the effects.
func GetPlayerTechBonuses(playerID string) (*TechBonuses, error) {
	bonuses := &TechBonuses{
		ConstructionSlots: 0,
		ShipProductionSlots: 0,
		MaxThorCannon: 0,
	}

	// Query all completed techs with their effects
	rows, err := database.DB.Query(`
		SELECT t.level, tt.effects_json, tt.name
		FROM technologies t
		JOIN tech_types tt ON t.tech_type = tt.id
		WHERE t.player_id = $1 AND t.level > 0
		ORDER BY tt.name
	`, playerID)
	if err != nil {
		log.Printf("Failed to query player research: %v", err)
		return nil, fmt.Errorf("failed to query player research: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var level int
		var effectsJSON []byte
		var techName string

		if err := rows.Scan(&level, &effectsJSON, &techName); err != nil {
			log.Printf("Failed to scan tech row: %v", err)
			continue
		}

		// Parse effects JSON
		var effects techEffectData
		if err := json.Unmarshal(effectsJSON, &effects); err != nil {
			log.Printf("Failed to parse effects for tech %s: %v", techName, err)
			continue
		}

		// Apply effect based on type
		applyTechEffect(bonuses, &effects, level, techName)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating tech rows: %v", err)
		return nil, fmt.Errorf("error iterating tech rows: %w", err)
	}

	return bonuses, nil
}

// applyTechEffect applies a single tech effect to the bonuses struct.
func applyTechEffect(bonuses *TechBonuses, effect *techEffectData, level int, techName string) {
	switch effect.Type {
	// Production
	case "metal_output":
		bonuses.MetalOutput += effect.PerLevel * float64(level)
	case "he3_output":
		bonuses.He3Output += effect.PerLevel * float64(level)
	case "gold_output":
		bonuses.GoldOutput += effect.PerLevel * float64(level)

	// Construction
	case "build_speed":
		bonuses.BuildSpeed += effect.PerLevel * float64(level)
	case "build_cost_reduction":
		bonuses.BuildCostReduction += effect.PerLevel * float64(level)
	case "construction_slots":
		bonuses.ConstructionSlots += int(effect.PerLevel * float64(level))
	case "ship_build_speed":
		bonuses.ShipBuildSpeed += effect.PerLevel * float64(level)
	case "ship_build_cost_reduction":
		bonuses.ShipBuildCostReduction += effect.PerLevel * float64(level)
	case "ship_production_slots":
		bonuses.ShipProductionSlots += int(effect.PerLevel * float64(level))

	// Capacity
	case "warehouse_capacity":
		bonuses.WarehouseCapacity += int64(effect.PerLevel * float64(level))

	// Ship repair
	case "ship_repair_percent":
		bonuses.ShipRepairPercent += effect.PerLevel * float64(level)

	// Ballistics
	case "ballistic_damage":
		bonuses.BallisticDamage += effect.PerLevel * float64(level)
	case "ballistic_crit_rate":
		bonuses.BallisticCritRate += effect.PerLevel * float64(level)
	case "ballistic_crit_damage":
		bonuses.BallisticCritDamage += effect.PerLevel * float64(level)
	case "ballistic_hit_rate":
		bonuses.BallisticHitRate += effect.PerLevel * float64(level)
	case "weapon_space_reduction":
		bonuses.WeaponSpaceReduction += effect.PerLevel * float64(level)
	case "shield_bypass":
		bonuses.ShieldBypass += effect.Flat

	// Directional
	case "directional_damage":
		bonuses.DirectionalDamage += effect.PerLevel * float64(level)
	case "directional_crit_rate":
		bonuses.DirectionalCritRate += effect.PerLevel * float64(level)
	case "directional_accuracy":
		bonuses.DirectionalAccuracy += effect.PerLevel * float64(level)
	case "steering_power":
		bonuses.SteeringPower += effect.PerLevel * float64(level)

	// Missile
	case "missile_damage":
		bonuses.MissileDamage += effect.PerLevel * float64(level)
	case "missile_hit_rate":
		bonuses.MissileHitRate += effect.PerLevel * float64(level)
	case "intercept_reduction":
		bonuses.InterceptReduction += effect.PerLevel * float64(level)

	// Fighter
	case "fighter_damage":
		bonuses.FighterDamage += effect.PerLevel * float64(level)
	case "fighter_hit_rate":
		bonuses.FighterHitRate += effect.PerLevel * float64(level)
	case "he3_cost_reduction":
		bonuses.FuelOptimization += effect.PerLevel * float64(level)

	// Defense stats
	case "base_defense_stats":
		// Ship Defense Tech provides multiple bonuses
		var baseDefenseEffect struct {
			Type      string  `json:"type"`
			Shield    float64 `json:"shield"`
			Structure float64 `json:"structure"`
			Agility   float64 `json:"agility"`
			Defense   float64 `json:"defense"`
			Stability float64 `json:"stability"`
		}
		// This is a composite effect, already handled separately
		_ = baseDefenseEffect
	case "base_shield":
		bonuses.BaseShield += effect.PerLevel * float64(level)
	case "base_structure":
		bonuses.BaseStructure += effect.PerLevel * float64(level)
	case "base_agility":
		bonuses.BaseAgility += effect.PerLevel * float64(level)
	case "base_defense":
		bonuses.BaseDefense += effect.PerLevel * float64(level)
	case "base_stability":
		bonuses.BaseStability += effect.PerLevel * float64(level)

	// Planetary defense
	case "defense_cost_reduction":
		bonuses.DefenseCostReduction += effect.PerLevel * float64(level)
	case "defense_build_speed":
		bonuses.DefenseBuildSpeed += effect.PerLevel * float64(level)
	case "defense_value":
		bonuses.DefenseValue += effect.PerLevel * float64(level)
	case "emplacement_attack":
		bonuses.EmplacementAttack += effect.PerLevel * float64(level)
	case "max_defense_structures":
		bonuses.MaxDefenseStructures += effect.PerLevel * float64(level)
	case "max_thor_cannon":
		bonuses.MaxThorCannon += int(effect.Flat)

	default:
		// Unknown effect type - log for future implementation
		log.Printf("Unknown tech effect type: %s (tech: %s)", effect.Type, techName)
	}
}

// GetTechLevel returns the current level of a specific technology for a player.
// Returns 0 if the technology hasn't been researched yet.
// This is used for prerequisite checks and condition evaluation.
func GetTechLevel(playerID string, techName string) (int, error) {
	var level int
	err := database.DB.QueryRow(`
		SELECT COALESCE(t.level, 0)
		FROM tech_types tt
		LEFT JOIN technologies t ON t.tech_type = tt.id AND t.player_id = $1
		WHERE tt.name = $2
	`, playerID, techName).Scan(&level)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		log.Printf("Failed to get tech level for %s: %v", techName, err)
		return 0, fmt.Errorf("failed to get tech level: %w", err)
	}

	return level, nil
}
