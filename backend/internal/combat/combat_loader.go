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

	fleet.TechBonuses = techBonuses

	// Load fleet stacks
	stacks, err := loadFleetStacks(fleetID)
	if err != nil {
		return nil, fmt.Errorf("failed to load fleet stacks: %w", err)
	}
	fleet.Stacks = stacks

	// Populate advanced combat properties from tech bonuses
	applyAdvancedCombatProps(&fleet)

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
		TechBonuses:    &services.TechBonuses{},
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
	// Speed/accuracy/dodge are NOT denormalised on ship_designs — they live on
	// hull_types as base values and are buffed at combat time by commander +
	// tech bonuses. We pull base_movement → speed, base_agility → dodge, and
	// default accuracy to 100 (modifiers stack on top).
	rows, err := database.DB.Query(`
		SELECT fs.id, fs.ship_design_id, fs.grid_row, fs.grid_col, fs.ship_count,
		       sd.hull_type_id, sd.attack_power, sd.total_defense,
		       ht.base_movement, 100 AS base_accuracy, ht.base_agility,
		       sd.total_shield, sd.total_structure
		FROM fleet_stacks fs
		JOIN ship_designs sd ON fs.ship_design_id = sd.id
		JOIN hull_types ht ON ht.id = sd.hull_type_id
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
		// total_defense is numeric(5,2) — scan into float64 then truncate.
		var totalDefense float64

		err := rows.Scan(
			&stack.ID, &stack.ShipDesignID, &stack.GridRow, &stack.GridCol, &stack.ShipCount,
			&hullTypeID,
			&stack.BaseAttack, &totalDefense, &stack.BaseSpeed,
			&stack.BaseAccuracy, &stack.BaseDodge, &stack.BaseShield, &stack.BaseStructure,
		)
		if err != nil {
			log.Printf("Failed to scan fleet stack: %v", err)
			continue
		}
		stack.BaseDefense = int(totalDefense)

		// Get hull type info for ship type classification
		stack.ShipType, stack.DamageType, stack.ArmorType, stack.WeaponCategory = getHullTypeInfo(hullTypeID)

		stacks = append(stacks, &stack)
	}

	return stacks, nil
}

// loadCommanderBonuses loads commander stats and converts to combat bonuses
func loadCommanderBonuses(commanderID string) (*CommanderBonus, error) {
	var accuracy, dodge, speed, electron, starRank int
	var weaponExpertise, shipExpertise sql.NullString

	err := database.DB.QueryRow(`
		SELECT accuracy, dodge, speed, electron, star_rank,
		       weapon_expertise, ship_expertise
		FROM commanders
		WHERE id = $1
	`, commanderID).Scan(&accuracy, &dodge, &speed, &electron, &starRank,
		&weaponExpertise, &shipExpertise)
	if err != nil {
		return nil, err
	}

	// EffectiveStack = 0.1 * star_rank (GDD formula)
	effectiveStackBonus := 0.1 * float64(starRank)

	bonus := &CommanderBonus{
		Accuracy:       accuracy,
		Dodge:          dodge,
		Speed:          speed,
		Electron:       electron,
		StarRank:       starRank,
		EffectiveStack: effectiveStackBonus,
	}
	if weaponExpertise.Valid {
		bonus.WeaponExpertise = weaponExpertise.String
	}
	if shipExpertise.Valid {
		bonus.ShipExpertise = shipExpertise.String
	}

	return bonus, nil
}

// getHullTypeInfo retrieves ship type and armor type from hull_types table.
// Damage type and weapon category are determined by equipped modules (via ship_designs),
// so this function returns defaults for those values.
func getHullTypeInfo(hullTypeID int) (ShipType, DamageType, ArmorType, WeaponCategory) {
	var hullClass, armorType string

	err := database.DB.QueryRow(`
		SELECT hull_class, armor_type
		FROM hull_types
		WHERE id = $1
	`, hullTypeID).Scan(&hullClass, &armorType)

	if err != nil {
		log.Printf("Failed to get hull type info: %v", err)
		return ShipTypeFrigate, DamageKinetic, ArmorChrome, WeaponBallistic
	}

	// Map hull_class to ShipType
	shipType := mapClassificationToShipType(hullClass)

	// Map armor_type to ArmorType
	armor := mapArmorClassToArmorType(armorType)

	// Damage type and weapon category come from modules, not the hull itself.
	// Default to kinetic/ballistic; player ships get these from ship_designs.
	return shipType, DamageKinetic, armor, WeaponBallistic
}

// createStackFromHullType creates a fleet stack from hull type name (for NPC fleets)
func createStackFromHullType(hullTypeName string, quantity, gridRow, gridCol int) (*FleetStack, error) {
	var hullTypeID int
	var hullClass, armorTypeStr string
	var baseShield, baseStructure, baseAgility, baseMovement int

	// Get hull type base stats
	err := database.DB.QueryRow(`
		SELECT id, hull_class, armor_type,
		       base_shield, base_structure, base_agility, base_movement
		FROM hull_types
		WHERE name = $1
	`, hullTypeName).Scan(
		&hullTypeID, &hullClass, &armorTypeStr,
		&baseShield, &baseStructure, &baseAgility, &baseMovement,
	)
	if err != nil {
		return nil, fmt.Errorf("hull type not found: %s", hullTypeName)
	}

	shipType := mapClassificationToShipType(hullClass)
	armorType := mapArmorClassToArmorType(armorTypeStr)

	// For NPC ships, use simplified base attack (can be enhanced later)
	baseAttack := baseShield / 10
	baseDefense := baseStructure / 10

	// Damage type and weapon category come from modules; NPC hulls default to kinetic/ballistic
	return &FleetStack{
		ShipCount:      quantity,
		ShipType:       shipType,
		DamageType:     DamageKinetic,
		ArmorType:      armorType,
		WeaponCategory: WeaponBallistic,
		BaseAttack:     baseAttack,
		BaseDefense:    baseDefense,
		BaseSpeed:      baseMovement,
		BaseAgility:    baseAgility,
		BaseShield:     baseShield,
		BaseStructure:  baseStructure,
		GridRow:        gridRow,
		GridCol:        gridCol,
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

func mapWeaponToCategory(weapon string) WeaponCategory {
	switch weapon {
	case "kinetic":
		return WeaponBallistic
	case "explosive":
		return WeaponMissile
	case "heat":
		return WeaponDirectional
	case "magnetic":
		return WeaponFighter
	default:
		return WeaponBallistic
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

// applyAdvancedCombatProps copies tech-derived advanced combat properties onto each fleet stack.
// These values are fleet-wide (from the tech tree) and apply equally to all stacks.
func applyAdvancedCombatProps(fleet *Fleet) {
	tb := fleet.TechBonuses
	if tb == nil {
		return
	}

	for _, stack := range fleet.Stacks {
		// Scatter / AoE
		stack.ScatterDamage = tb.ScatterDamage
		stack.ScatterAll = tb.ScatterAll
		stack.ScatterRate = tb.ScatterRate / 100.0 // tech stores as percent, engine expects 0-1
		stack.ScatterBonus = tb.ScatterBonus
		stack.ScatterVsHighStructure = tb.ScatterVsHighStructure

		// Piercing
		stack.PiercingDamage = tb.PiercingDamage
		stack.PiercingDamageBonus = tb.PiercingDamageBonus
		stack.PiercingCritical = tb.PiercingCriticalEnabled

		// Restoration
		stack.ShieldRestore = tb.ShieldRestore
		stack.StructureRestore = tb.StructureRestore
		stack.AbsorbDouble = tb.AbsorbDouble / 100.0 // tech stores as percent, engine expects 0-1

		// Reflection
		stack.ReflectDamage = tb.ReflectDamage
		stack.ReflectStructureDamage = tb.ReflectStructureDamage

		// Positional: knockback
		stack.Knockback = tb.KnockbackDistance

		// Positional: range damage
		if len(tb.RangeDamageRanges) > 0 {
			stack.RangeDamage = make(map[int]float64)
			for dist, mult := range tb.RangeDamageRanges {
				stack.RangeDamage[dist] = 1.0 + mult/100.0 // tech stores bonus as percent
			}
		}

		// Debuffs
		stack.EnemyAttackReduction = tb.EnemyAttackReduction
		if tb.EnemyAttackReduction > 0 {
			stack.EnemyAttackReductionRounds = 3 // default duration
		}
	}
}
