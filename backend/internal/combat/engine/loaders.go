package engine

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/cryptomines-online/backend/internal/combat/tilemap"
	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/services"
)

// defenseBuildingProfile is the static per-type configuration the loader
// uses to translate a row in `buildings` into an engine.Unit. Range,
// AoE, and footprint follow GDD §2.2.9 + §8.6.4.
type defenseBuildingProfile struct {
	weapon     WeaponClass
	weaponName string
	dmgType    DamageType
	armor      ArmorType
	minRange   int
	maxRange   int
	cooldown   int
	speed      int // initiative
	cols, rows int // footprint
}

var defenseProfiles = map[string]defenseBuildingProfile{
	"space_station": {
		weapon: WeaponPlanetary, weaponName: "station_turret",
		dmgType: DamageExplosive, armor: ArmorChrome,
		minRange: 1, maxRange: 5, cooldown: 1, speed: 20, cols: 3, rows: 3,
	},
	"particle_cannon": {
		weapon: WeaponPlanetary, weaponName: "particle_cannon",
		dmgType: DamageMagnetic, armor: ArmorChrome,
		minRange: 8, maxRange: 17, cooldown: 1, speed: 40, cols: 1, rows: 2,
	},
	"anti_aircraft_gun": {
		weapon: WeaponPlanetaryAOE, weaponName: "anti_aircraft_gun",
		dmgType: DamageExplosive, armor: ArmorChrome,
		minRange: 1, maxRange: 30, cooldown: 1, speed: 50, cols: 1, rows: 1,
	},
	"thors_cannon": {
		weapon: WeaponPlanetary, weaponName: "thors_cannon",
		dmgType: DamageHeat, armor: ArmorChrome,
		minRange: 1, maxRange: 25, cooldown: 2, speed: 10, cols: 2, rows: 2,
	},
	"celestial_base": {
		weapon: WeaponPlanetary, weaponName: "celestial_base",
		dmgType: DamageExplosive, armor: ArmorChrome,
		minRange: 1, maxRange: 8, cooldown: 1, speed: 25, cols: 2, rows: 2,
	},
}

// damageTable returns the (min,max) per-shot damage and the (shield,structure)
// HP for a defense building of the given type at the given level. Numbers
// follow GDD §2.7.3 RBP defense scaling (level^2). Tech bonuses are
// applied by the caller.
func damageTable(name string, level int) (dmgMin, dmgMax, shield, structHP int) {
	l2 := level * level
	switch name {
	case "space_station":
		return 4000 * l2, 6000 * l2, 200000 * l2, 300000 * l2
	case "particle_cannon":
		return 10000 * l2, 45000 * l2, 50000 * l2, 50000 * l2
	case "anti_aircraft_gun":
		return 50000 * l2, 500000 * l2, 60000 * l2, 100000 * l2
	case "thors_cannon":
		return 50000 * l2, 1350000 * l2, 100000 * l2, 150000 * l2
	case "celestial_base":
		return 5000 * l2, 8000 * l2, 80000 * l2, 120000 * l2
	}
	return 0, 0, 0, 0
}

// LoadDefenseUnits loads all defense buildings on a planet (excluding
// Meteor Stars, which load as terrain) and converts them to engine
// units placed on their build-time tile positions. Tech bonuses are
// applied to HP and damage.
//
// Multi-tile buildings register every tile in their footprint with the
// battle map (multiple Footprint entries → single Unit). The tile
// returned in Unit.Tile is the anchor (top-left).
func LoadDefenseUnits(planetID string, tb *services.TechBonuses) ([]*Unit, error) {
	rows, err := database.DB.Query(`
		SELECT b.id, bt.name, b.level, b.grid_col, b.grid_row
		FROM buildings b
		JOIN building_types bt ON bt.id = b.building_type
		WHERE b.planet_id = $1
		  AND bt.category = 'defense'
		  AND b.is_upgrading = false
		  AND b.level > 0
		  AND bt.name <> 'meteor_star'
	`, planetID)
	if err != nil {
		return nil, fmt.Errorf("LoadDefenseUnits query: %w", err)
	}
	defer rows.Close()

	units := []*Unit{}
	for rows.Next() {
		var bid, name string
		var level, gc, gr int
		if err := rows.Scan(&bid, &name, &level, &gc, &gr); err != nil {
			log.Printf("LoadDefenseUnits scan: %v", err)
			continue
		}
		profile, ok := defenseProfiles[name]
		if !ok {
			continue
		}
		dmgMin, dmgMax, shield, structHP := damageTable(name, level)
		applyTechToDefense(name, tb, &dmgMin, &dmgMax, &shield, &structHP)

		anchor := tilemap.Position{Col: gc, Row: gr}
		footprint := tilemap.Footprint(anchor, profile.cols, profile.rows)
		unit := &Unit{
			ID: "def_" + bid, Side: tilemap.SideDefender, Kind: UnitBuilding,
			BuildingID: bid, TypeName: name,
			Tile: anchor, Footprint: footprint,
			Speed: profile.speed, Movement: 0,
			Phalanx: Phalanx{Stacks: []*Stack{{
				ID: bid + "_s0", GridRow: 0, GridCol: 0,
				HullClass:   HullBuilding,
				ShipCount:   structHP, // for buildings, ShipCount tracks HP
				WeaponClass: profile.weapon, WeaponName: profile.weaponName,
				DamageType:  profile.dmgType, ArmorType: profile.armor,
				MinRange:    profile.minRange, MaxRange: profile.maxRange,
				Cooldown:    profile.cooldown,
				BaseDamageMin: dmgMin, BaseDamageMax: dmgMax,
				HitChance:   0.85,
				ShieldHP:    shield, MaxShield: shield,
				StructHP:    structHP, MaxStruct: structHP,
				Accuracy:    rangeAccuracyBonus(name, tb),
			}}},
		}
		units = append(units, unit)
	}
	return units, nil
}

func applyTechToDefense(name string, tb *services.TechBonuses, dmgMin, dmgMax, shield, structHP *int) {
	if tb == nil {
		return
	}
	if tb.DefenseValue > 0 {
		mult := 1.0 + tb.DefenseValue/100.0
		*shield = int(float64(*shield) * mult)
		*structHP = int(float64(*structHP) * mult)
	}
	if tb.EmplacementAttack > 0 {
		switch name {
		case "particle_cannon", "anti_aircraft_gun", "thors_cannon":
			mult := 1.0 + tb.EmplacementAttack/100.0
			*dmgMin = int(float64(*dmgMin) * mult)
			*dmgMax = int(float64(*dmgMax) * mult)
		}
	}
}

func rangeAccuracyBonus(name string, tb *services.TechBonuses) int {
	if tb == nil {
		return 0
	}
	switch name {
	case "particle_cannon", "anti_aircraft_gun":
		return tb.DefenseRange
	}
	return 0
}

// LoadTerrain returns Meteor Stars on the planet as terrain entities
// (GDD §8.6.2). They block movement and have HP scaled by level.
func LoadTerrain(planetID string) ([]*tilemap.Terrain, error) {
	rows, err := database.DB.Query(`
		SELECT b.level, b.grid_col, b.grid_row
		FROM buildings b
		JOIN building_types bt ON bt.id = b.building_type
		WHERE b.planet_id = $1
		  AND bt.name = 'meteor_star'
		  AND b.is_upgrading = false
		  AND b.level > 0
	`, planetID)
	if err != nil {
		return nil, fmt.Errorf("LoadTerrain query: %w", err)
	}
	defer rows.Close()

	out := []*tilemap.Terrain{}
	for rows.Next() {
		var level, gc, gr int
		if err := rows.Scan(&level, &gc, &gr); err != nil {
			continue
		}
		hp := 40000 * level * level // GDD §2.7.3: 40K..20.48M across levels 1..10
		out = append(out, &tilemap.Terrain{
			Kind:     tilemap.TerrainMeteorStar,
			Position: tilemap.Position{Col: gc, Row: gr},
			HP:       hp, MaxHP: hp,
		})
	}
	return out, nil
}

// LoadFleetUnit reads one fleet from the database and constructs a
// Unit with a populated phalanx. Tile is left zero — the spawn helper
// assigns it. Side is provided by the caller (attacker / defender).
//
// Stack stats are translated:
//
//	ship_designs.attack_power     → BaseDamageMin/Max (split as ±20%)
//	ship_designs.total_shield     → ShieldHP (multiplied by ShipCount)
//	ship_designs.total_structure  → StructHP  (multiplied by ShipCount)
//	hull_types.base_movement      → Movement (used at unit level)
//	hull_types.base_agility       → Dodge
//	hull_types.hull_class         → HullClass + ArmorType
func LoadFleetUnit(fleetID, playerID string, side tilemap.Side) (*Unit, error) {
	var name, formation, targeting string
	var commanderID sql.NullString
	err := database.DB.QueryRow(`
		SELECT name, formation, targeting_command, commander_id
		FROM fleets WHERE id = $1 AND player_id = $2
	`, fleetID, playerID).Scan(&name, &formation, &targeting, &commanderID)
	if err != nil {
		return nil, fmt.Errorf("LoadFleetUnit fleet: %w", err)
	}
	commanderSpeed := 0
	if commanderID.Valid {
		_ = database.DB.QueryRow(`SELECT speed FROM commanders WHERE id = $1`, commanderID.String).Scan(&commanderSpeed)
	}

	rows, err := database.DB.Query(`
		SELECT fs.id, fs.ship_design_id, fs.grid_row, fs.grid_col, fs.ship_count,
		       sd.attack_power, sd.total_shield, sd.total_structure,
		       ht.hull_class, ht.armor_type, ht.base_movement, ht.base_agility
		FROM fleet_stacks fs
		JOIN ship_designs sd ON sd.id = fs.ship_design_id
		JOIN hull_types ht ON ht.id = sd.hull_type_id
		WHERE fs.fleet_id = $1 AND fs.ship_count > 0
	`, fleetID)
	if err != nil {
		return nil, fmt.Errorf("LoadFleetUnit stacks: %w", err)
	}
	defer rows.Close()

	unit := &Unit{
		ID: "atk_" + fleetID, Side: side, Kind: UnitFleet, FleetID: fleetID,
		TargetCmd: targeting, Movement: 999, // narrowed below to MIN(stack.Movement)
		Speed:     commanderSpeed,
		Phalanx:   Phalanx{Stacks: []*Stack{}},
	}
	if side == tilemap.SideDefender {
		unit.ID = "def_" + fleetID
	}

	for rows.Next() {
		var stackID, designID, hullClass, armor string
		var gr, gc, count, attackPower, shield, structHP, mov, agility int
		if err := rows.Scan(&stackID, &designID, &gr, &gc, &count,
			&attackPower, &shield, &structHP,
			&hullClass, &armor, &mov, &agility); err != nil {
			log.Printf("LoadFleetUnit scan: %v", err)
			continue
		}
		dmgMin := int(float64(attackPower) * 0.80)
		dmgMax := int(float64(attackPower) * 1.20)
		hull := mapHullClass(hullClass)
		s := &Stack{
			ID: stackID, GridRow: gr, GridCol: gc,
			ShipDesignID: designID, HullClass: hull,
			ShipCount:    count,
			WeaponClass:  defaultWeaponClassForHull(hull),
			WeaponName:   "ship_weapons",
			DamageType:   DamageKinetic,
			ArmorType:    mapArmor(armor),
			MinRange:     1, MaxRange: 4,
			BaseDamageMin: dmgMin, BaseDamageMax: dmgMax,
			HitChance:    0.80,
			ShieldHP:     shield * count, MaxShield: shield * count,
			StructHP:     structHP * count, MaxStruct: structHP * count,
			Accuracy:     0,
			Dodge:        agility,
			Movement:     mov,
		}
		unit.Phalanx.Stacks = append(unit.Phalanx.Stacks, s)
		if mov > 0 && mov < unit.Movement {
			unit.Movement = mov
		}
	}
	if len(unit.Phalanx.Stacks) == 0 {
		unit.Movement = 0
	}
	if unit.Movement == 999 {
		unit.Movement = 1
	}
	_ = name
	_ = formation
	return unit, nil
}

func mapHullClass(s string) HullClass {
	switch s {
	case "frigate":
		return HullFrigate
	case "cruiser":
		return HullCruiser
	case "battleship":
		return HullBattleship
	}
	return HullFrigate
}

func mapArmor(s string) ArmorType {
	switch s {
	case "chrome":
		return ArmorChrome
	case "regen":
		return ArmorRegen
	case "nano":
		return ArmorNano
	case "neutralizing":
		return ArmorNeutralizing
	}
	return ArmorNone
}

func defaultWeaponClassForHull(h HullClass) WeaponClass {
	switch h {
	case HullFrigate:
		return WeaponBallistic
	case HullCruiser:
		return WeaponDirectional
	case HullBattleship:
		return WeaponMissile
	}
	return WeaponBallistic
}

// SpaceStationAnchor returns the (col,row) and footprint dims of the
// defender's Space Station building, used as the focal point for
// defending-fleet spawns. Returns ok=false when no Space Station exists.
func SpaceStationAnchor(planetID string) (anchor tilemap.Position, cols, rows int, ok bool) {
	var gc, gr int
	err := database.DB.QueryRow(`
		SELECT b.grid_col, b.grid_row
		FROM buildings b
		JOIN building_types bt ON bt.id = b.building_type
		WHERE b.planet_id = $1 AND bt.name = 'space_station'
		  AND b.is_upgrading = false AND b.level > 0
		LIMIT 1
	`, planetID).Scan(&gc, &gr)
	if err != nil {
		return tilemap.Position{}, 0, 0, false
	}
	return tilemap.Position{Col: gc, Row: gr}, 3, 3, true
}
