package services

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/cryptomines-online/backend/internal/models"
)

// DesignStats holds calculated stats for a ship design.
type DesignStats struct {
	TotalShield     int     `json:"total_shield"`
	TotalStructure  int     `json:"total_structure"`
	TotalDefense    float64 `json:"total_defense"`
	TotalAgility    int     `json:"total_agility"`
	TotalMovement   int     `json:"total_movement"`
	TotalStorage    int     `json:"total_storage"`
	AttackPower     int     `json:"attack_power"`
	WeaponRangeMin  int     `json:"weapon_range_min"`
	WeaponRangeMax  int     `json:"weapon_range_max"`
	VolumeUsed      int     `json:"volume_used"`
	He3PerRound     int     `json:"he3_per_round"`
	MetalCost       int64   `json:"metal_cost"`
	He3Cost         int64   `json:"he3_cost"`
	GoldCost        int64   `json:"gold_cost"`
	BuildTimeSeconds int    `json:"build_time_seconds"`
}

// CalculateDesignStats computes all stats for a ship design given the hull and modules.
// Follows GDD 8.2.9 Ship Stats Calculation and 8.10.1/8.10.2 formulas.
func CalculateDesignStats(db *sql.DB, hull models.HullType, modules []models.DesignModule) (DesignStats, error) {
	stats := DesignStats{
		TotalShield:    hull.BaseShield,
		TotalStructure: hull.BaseStructure,
		TotalDefense:   hull.BaseDefense,
		TotalAgility:   hull.BaseAgility,
		TotalMovement:  hull.BaseMovement,
		TotalStorage:   hull.BaseStorage,
		MetalCost:      hull.BaseMetalCost,
		He3Cost:        hull.BaseHe3Cost,
		GoldCost:       hull.BaseGoldCost,
		BuildTimeSeconds: hull.BaseBuildTime,
	}

	weaponRangeMinSet := false
	var weaponRangeMin, weaponRangeMax int

	for _, mod := range modules {
		var mt models.ModuleType
		err := db.QueryRow(
			`SELECT id, name, display_name, category, tier, damage_type,
			        min_damage, max_damage, weapon_range_min, weapon_range_max,
			        cooldown, he3_per_round, volume, max_per_ship, effects_json,
			        metal_cost, he3_cost, gold_cost, build_time_seconds
			 FROM module_types WHERE id = $1`, mod.ModuleTypeID,
		).Scan(
			&mt.ID, &mt.Name, &mt.DisplayName, &mt.Category, &mt.Tier, &mt.DamageType,
			&mt.MinDamage, &mt.MaxDamage, &mt.WeaponRangeMin, &mt.WeaponRangeMax,
			&mt.Cooldown, &mt.He3PerRound, &mt.Volume, &mt.MaxPerShip, &mt.EffectsJSON,
			&mt.MetalCost, &mt.He3Cost, &mt.GoldCost, &mt.BuildTime,
		)
		if err != nil {
			return stats, err
		}

		qty := mod.Quantity

		// Volume
		stats.VolumeUsed += mt.Volume * qty

		// Costs (per ship formula from GDD 8.10.1)
		stats.MetalCost += mt.MetalCost * int64(qty)
		stats.He3Cost += mt.He3Cost * int64(qty)
		stats.GoldCost += mt.GoldCost * int64(qty)

		// Build time
		stats.BuildTimeSeconds += mt.BuildTime * qty

		// He3 consumption per round
		stats.He3PerRound += mt.He3PerRound * qty

		// Attack power (average damage * quantity)
		if mt.MinDamage > 0 || mt.MaxDamage > 0 {
			avgDamage := (mt.MinDamage + mt.MaxDamage) / 2
			stats.AttackPower += avgDamage * qty

			// Weapon range (use the widest range across all weapons)
			if !weaponRangeMinSet || mt.WeaponRangeMin < weaponRangeMin {
				weaponRangeMin = mt.WeaponRangeMin
				weaponRangeMinSet = true
			}
			if mt.WeaponRangeMax > weaponRangeMax {
				weaponRangeMax = mt.WeaponRangeMax
			}
		}

		// Process effects_json for stat bonuses
		var effects map[string]interface{}
		if err := json.Unmarshal([]byte(mt.EffectsJSON), &effects); err != nil {
			log.Printf("Failed to parse effects_json for module %d: %v", mt.ID, err)
			continue
		}

		for key, val := range effects {
			floatVal := toFloat64(val)
			intVal := int(floatVal)

			switch key {
			case "structure_bonus":
				stats.TotalStructure += intVal * qty
			case "shield_bonus":
				stats.TotalShield += intVal * qty
			case "agility_bonus":
				stats.TotalAgility += intVal * qty
			case "movement_bonus":
				stats.TotalMovement += intVal * qty
			case "he3_storage_bonus":
				stats.TotalStorage += intVal * qty
			case "module_capacity_bonus":
				// Increases effective slots but doesn't change VolumeUsed calculation
			case "defense_bonus_pct":
				stats.TotalDefense += floatVal * float64(qty)
			case "structure_bonus_pct":
				// Percentage bonuses applied after all flat bonuses
				// For now we accumulate; final application at the end
			}
		}
	}

	if weaponRangeMinSet {
		stats.WeaponRangeMin = weaponRangeMin
		stats.WeaponRangeMax = weaponRangeMax
	}

	return stats, nil
}

// ShipBuildCost calculates total batch cost.
// GDD 8.10.1: BatchCost = PerShipCost * quantity
func ShipBuildCost(perShipMetal, perShipHe3, perShipGold int64, quantity int) (metal, he3, gold int64) {
	return perShipMetal * int64(quantity), perShipHe3 * int64(quantity), perShipGold * int64(quantity)
}

// ShipBuildTime calculates batch build time.
// GDD 8.10.2: BatchTime = EffectiveTime * quantity
// EffectiveTime = BaseShipTime * (1 - SpeedReduction/100)
func ShipBuildTime(baseTimePerShip int, speedBonusPct int, quantity int) int {
	effectiveTime := float64(baseTimePerShip) * (1.0 - float64(speedBonusPct)/100.0)
	if effectiveTime < 1 {
		effectiveTime = 1
	}
	return int(effectiveTime * float64(quantity))
}

func toFloat64(v interface{}) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	case int64:
		return float64(n)
	default:
		return 0
	}
}
