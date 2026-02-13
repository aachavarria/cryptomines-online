package handlers

import (
	"log"

	"github.com/cryptomines-online/backend/internal/combat"
	"github.com/cryptomines-online/backend/internal/database"
)

// loadDefenseBuildings converts defense buildings on a planet into combat stacks.
// Shared between PvP, RBP, and worker code.
func loadDefenseBuildings(planetID string) ([]*combat.FleetStack, error) {
	rows, err := database.DB.Query(`
		SELECT b.id, bt.name, b.level
		FROM buildings b
		JOIN building_types bt ON bt.name = b.building_type
		WHERE b.planet_id = $1
		  AND bt.type = 'defense'
		  AND b.construction_end_time IS NULL
	`, planetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stacks []*combat.FleetStack
	stackID := 0
	for rows.Next() {
		var buildingID, buildingName string
		var level int

		if err := rows.Scan(&buildingID, &buildingName, &level); err != nil {
			continue
		}

		stack := buildingToStack(buildingID, buildingName, level, stackID)
		if stack != nil {
			stacks = append(stacks, stack)
			stackID++
		}
	}

	return stacks, nil
}

// buildingToStack converts a defense building to a combat stack
// Stats use level^2 scaling to match GO2 power scale (e.g. Particle Cannon Lv10 = 1M attack, 5M HP)
func buildingToStack(buildingID, buildingName string, level, stackID int) *combat.FleetStack {
	var baseAttack, baseDefense, baseShield, baseStructure, baseSpeed, baseAccuracy, baseDodge int
	l2 := level * level // level^2 scaling

	switch buildingName {
	case "space_station":
		// Tankiest: huge HP, moderate attack
		baseAttack = 5000 * l2
		baseDefense = 10000 * l2
		baseShield = 200000 * l2
		baseStructure = 300000 * l2
		baseSpeed = 20
		baseAccuracy = 60
		baseDodge = 10
	case "particle_cannon":
		// High attack, moderate HP (GDD: Lv10 = 1M attack, 5M HP)
		baseAttack = 10000 * l2
		baseDefense = 5000 * l2
		baseShield = 50000 * l2
		baseStructure = 50000 * l2
		baseSpeed = 40
		baseAccuracy = 90
		baseDodge = 20
	case "anti_aircraft_gun":
		// Anti-air specialist: high accuracy, balanced stats
		baseAttack = 7500 * l2
		baseDefense = 6000 * l2
		baseShield = 60000 * l2
		baseStructure = 100000 * l2
		baseSpeed = 50
		baseAccuracy = 100
		baseDodge = 30
	case "meteor_star":
		// All-rounder: balanced attack and defense
		baseAttack = 5000 * l2
		baseDefense = 8000 * l2
		baseShield = 80000 * l2
		baseStructure = 120000 * l2
		baseSpeed = 35
		baseAccuracy = 75
		baseDodge = 25
	case "thors_cannon":
		// Strongest attack: massive damage, slow
		baseAttack = 20000 * l2
		baseDefense = 10000 * l2
		baseShield = 100000 * l2
		baseStructure = 150000 * l2
		baseSpeed = 10
		baseAccuracy = 80
		baseDodge = 5
	default:
		return nil
	}

	return &combat.FleetStack{
		ID:               buildingID,
		ShipDesignID:     "",
		ShipType:         combat.ShipTypeCruiser,
		DamageType:       combat.DamageExplosive,
		ArmorType:        combat.ArmorChrome,
		ShipCount:        1,
		CurrentShips:     1,
		BaseAttack:       baseAttack,
		BaseDefense:      baseDefense,
		BaseSpeed:        baseSpeed,
		BaseAccuracy:     baseAccuracy,
		BaseDodge:        baseDodge,
		BaseShield:       baseShield,
		BaseStructure:    baseStructure,
		CurrentShield:    baseShield,
		CurrentStructure: baseStructure,
		GridRow:          stackID / 3,
		GridCol:          stackID % 3,
	}
}

// resetDefenseBuildings resets all defense buildings on a planet to level 0 after PvP defeat.
// Per GO2 wiki: "all defenses reset automatically at no cost" after a successful attack.
func resetDefenseBuildings(planetID string) {
	defenseTypes := []string{"meteor_star", "particle_cannon", "anti_aircraft_gun", "thors_cannon", "celestial_base"}
	for _, typeName := range defenseTypes {
		_, err := database.DB.Exec(`
			UPDATE buildings SET level = 0, is_upgrading = false, upgrade_finish_at = NULL, updated_at = now()
			WHERE planet_id = $1 AND building_type_id = (
				SELECT id FROM building_types WHERE name = $2
			) AND level > 0
		`, planetID, typeName)
		if err != nil {
			log.Printf("Failed to reset defense building %s: %v", typeName, err)
		}
	}
}
