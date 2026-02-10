package combat

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/services"
)

// LoadPlayerFleet loads a player's fleet from the database and converts it to combat format
func LoadPlayerFleet(fleetID, playerID string) (*Fleet, error) {
	// Get fleet info
	var fleet Fleet
	var commanderID sql.NullString
	var formation, targeting string

	err := database.DB.QueryRow(`
		SELECT id, player_id, formation, targeting_command, commander_id
		FROM fleets
		WHERE id = $1 AND player_id = $2
	`, fleetID, playerID).Scan(
		&fleet.FleetID,
		&fleet.PlayerID,
		&formation,
		&targeting,
		&commanderID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load fleet: %w", err)
	}

	fleet.Formation = formation
	fleet.Targeting = targeting
	fleet.Side = "attacker"

	// Load commander bonuses if assigned
	if commanderID.Valid {
		commander, err := loadCommanderBonuses(commanderID.String)
		if err != nil {
			log.Printf("Failed to load commander bonuses: %v", err)
		} else {
			fleet.CommanderBonus = commander
		}
	}

	// Load tech bonuses
	techBonuses, err := services.GetPlayerTechBonuses(playerID)
	if err != nil {
		log.Printf("Failed to load tech bonuses: %v", err)
		techBonuses = &services.TechBonuses{}
	}

	// Convert services.TechBonuses to combat.TechBonuses
	fleet.TechBonuses = &TechBonuses{
		BallisticDamage:     techBonuses.BallisticDamage,
		BallisticCritRate:   techBonuses.BallisticCritRate,
		BallisticCritDamage: techBonuses.BallisticCritDamage,
		BallisticHitRate:    techBonuses.BallisticHitRate,
		DirectionalDamage:   techBonuses.DirectionalDamage,
		DirectionalCritRate: techBonuses.DirectionalCritRate,
		DirectionalAccuracy: techBonuses.DirectionalAccuracy,
		MissileDamage:       techBonuses.MissileDamage,
		MissileHitRate:      techBonuses.MissileHitRate,
		BaseShield:          techBonuses.BaseShield,
		BaseStructure:       techBonuses.BaseStructure,
		BaseAgility:         techBonuses.BaseAgility,
		BaseDefense:         techBonuses.BaseDefense,
	}

	// Load fleet stacks
	stacks, err := loadFleetStacks(fleetID)
	if err != nil {
		return nil, fmt.Errorf("failed to load fleet stacks: %w", err)
	}
	fleet.Stacks = stacks

	return &fleet, nil
}

// LoadInstanceFleet creates a combat fleet from instance enemy configuration
func LoadInstanceFleet(instanceID string) (*Fleet, error) {
	// Get instance enemy fleets JSON
	var enemyFleetsJSON []byte
	err := database.DB.QueryRow(`
		SELECT enemy_fleets_json FROM instances WHERE id = $1
	`, instanceID).Scan(&enemyFleetsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to load instance: %w", err)
	}

	// Parse enemy fleets (assuming JSON structure: [{"hull_type": "frigate_i", "quantity": 100, ...}])
	var enemyConfig []struct {
		HullType  string `json:"hull_type"`
		Quantity  int    `json:"quantity"`
		GridRow   int    `json:"grid_row"`
		GridCol   int    `json:"grid_col"`
	}

	if err := json.Unmarshal(enemyFleetsJSON, &enemyConfig); err != nil {
		return nil, fmt.Errorf("failed to parse enemy fleets: %w", err)
	}

	fleet := &Fleet{
		PlayerID:       "npc",
		FleetID:        "instance_" + instanceID,
		CommanderBonus: nil, // NPCs don't have commanders
		TechBonuses:    &TechBonuses{},
		Stacks:         []*FleetStack{},
		Formation:      "phalanx",
		Targeting:      "max_attack",
		Side:           "defender",
	}

	// Convert enemy config to fleet stacks
	for i, enemy := range enemyConfig {
		stack, err := createStackFromHullType(enemy.HullType, enemy.Quantity, enemy.GridRow, enemy.GridCol)
		if err != nil {
			log.Printf("Failed to create stack from hull %s: %v", enemy.HullType, err)
			continue
		}
		stack.ID = fmt.Sprintf("enemy_%d", i)
		fleet.Stacks = append(fleet.Stacks, stack)
	}

	return fleet, nil
}

// loadFleetStacks loads all stacks for a fleet and converts them to combat format
func loadFleetStacks(fleetID string) ([]*FleetStack, error) {
	rows, err := database.DB.Query(`
		SELECT fs.id, fs.ship_design_id, fs.grid_row, fs.grid_col, fs.ship_count,
		       sd.hull_type_id, sd.total_weapon_damage, sd.total_defense, sd.total_speed,
		       sd.total_accuracy, sd.total_dodge, sd.total_shield, sd.total_structure
		FROM fleet_stacks fs
		JOIN ship_designs sd ON fs.ship_design_id = sd.id
		WHERE fs.fleet_id = $1 AND fs.ship_count > 0
	`, fleetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stacks []*FleetStack
	for rows.Next() {
		var stack FleetStack
		var hullTypeID int

		err := rows.Scan(
			&stack.ID, &stack.ShipDesignID, &stack.GridRow, &stack.GridCol, &stack.ShipCount,
			&hullTypeID,
			&stack.BaseAttack, &stack.BaseDefense, &stack.BaseSpeed,
			&stack.BaseAccuracy, &stack.BaseDodge, &stack.BaseShield, &stack.BaseStructure,
		)
		if err != nil {
			log.Printf("Failed to scan fleet stack: %v", err)
			continue
		}

		// Get hull type info for ship type classification
		stack.ShipType, stack.DamageType, stack.ArmorType = getHullTypeInfo(hullTypeID)

		stacks = append(stacks, &stack)
	}

	return stacks, nil
}

// loadCommanderBonuses loads commander stats and converts to combat bonuses
func loadCommanderBonuses(commanderID string) (*CommanderBonus, error) {
	var accuracy, dodge, speed, electron int

	err := database.DB.QueryRow(`
		SELECT accuracy, dodge, speed, electron
		FROM commanders
		WHERE id = $1
	`, commanderID).Scan(&accuracy, &dodge, &speed, &electron)
	if err != nil {
		return nil, err
	}

	// TODO: Calculate effective stack bonus from commander star rank
	// For now, use electron as effective stack bonus (simplified)
	effectiveStackBonus := float64(electron) * 0.1 // 10% per 10 electron

	return &CommanderBonus{
		Accuracy:      accuracy,
		Dodge:         dodge,
		Speed:         speed,
		Electron:      electron,
		EffectiveStack: effectiveStackBonus,
	}, nil
}

// getHullTypeInfo retrieves ship type, damage type, and armor type from hull_types table
func getHullTypeInfo(hullTypeID int) (ShipType, DamageType, ArmorType) {
	var classification, primaryWeapon, armorClass string

	err := database.DB.QueryRow(`
		SELECT classification, primary_weapon_type, armor_class
		FROM hull_types
		WHERE id = $1
	`, hullTypeID).Scan(&classification, &primaryWeapon, &armorClass)

	if err != nil {
		log.Printf("Failed to get hull type info: %v", err)
		return ShipTypeFrigate, DamageKinetic, ArmorChrome
	}

	// Map classification to ShipType
	shipType := mapClassificationToShipType(classification)

	// Map primary weapon to DamageType
	damageType := mapWeaponToDamageType(primaryWeapon)

	// Map armor class to ArmorType
	armorType := mapArmorClassToArmorType(armorClass)

	return shipType, damageType, armorType
}

// createStackFromHullType creates a fleet stack from hull type name (for NPC fleets)
func createStackFromHullType(hullTypeName string, quantity, gridRow, gridCol int) (*FleetStack, error) {
	var hullTypeID int
	var classification, primaryWeapon, armorClass string
	var baseShield, baseStructure, baseSpeed, baseAccuracy, baseDodge int

	// Get hull type base stats
	err := database.DB.QueryRow(`
		SELECT id, classification, primary_weapon_type, armor_class,
		       base_shield, base_structure, base_speed, base_accuracy, base_dodge
		FROM hull_types
		WHERE name = $1
	`, hullTypeName).Scan(
		&hullTypeID, &classification, &primaryWeapon, &armorClass,
		&baseShield, &baseStructure, &baseSpeed, &baseAccuracy, &baseDodge,
	)
	if err != nil {
		return nil, fmt.Errorf("hull type not found: %s", hullTypeName)
	}

	shipType := mapClassificationToShipType(classification)
	damageType := mapWeaponToDamageType(primaryWeapon)
	armorType := mapArmorClassToArmorType(armorClass)

	// For NPC ships, use simplified base attack (can be enhanced later)
	baseAttack := baseShield / 10
	baseDefense := baseStructure / 10

	return &FleetStack{
		ShipCount:     quantity,
		ShipType:      shipType,
		DamageType:    damageType,
		ArmorType:     armorType,
		BaseAttack:    baseAttack,
		BaseDefense:   baseDefense,
		BaseSpeed:     baseSpeed,
		BaseAccuracy:  baseAccuracy,
		BaseDodge:     baseDodge,
		BaseShield:    baseShield,
		BaseStructure: baseStructure,
		GridRow:       gridRow,
		GridCol:       gridCol,
	}, nil
}

// Mapping helpers

func mapClassificationToShipType(classification string) ShipType {
	switch classification {
	case "frigate":
		return ShipTypeFrigate
	case "cruiser":
		return ShipTypeCruiser
	case "battleship":
		return ShipTypeBattleship
	default:
		return ShipTypeFrigate
	}
}

func mapWeaponToDamageType(weapon string) DamageType {
	switch weapon {
	case "kinetic":
		return DamageKinetic
	case "explosive":
		return DamageExplosive
	case "heat":
		return DamageHeat
	case "magnetic":
		return DamageMagnetic
	default:
		return DamageKinetic
	}
}

func mapArmorClassToArmorType(armorClass string) ArmorType {
	switch armorClass {
	case "chrome":
		return ArmorChrome
	case "regen":
		return ArmorRegen
	case "nano":
		return ArmorNano
	case "neutralizing":
		return ArmorNeutralizing
	default:
		return ArmorChrome
	}
}
