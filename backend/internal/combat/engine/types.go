// Package engine implements the tile-based tactical combat loop described
// in GDD §8.6. It composes:
//
//   - tilemap (spatial primitives)
//   - pathing (A* movement)
//   - spawn   (corner picker, defender placement)
//   - ai      (target selection, destination picking — see ai.go)
//   - replay  (action recorder for the rounds_json schema)
//
// The 8 damage sub-phases (GDD §8.6.9) are applied per-shot inside the
// fire step. They are the only piece reused from the legacy engine: the
// math is correct, only the spatial framing changed.
package engine

import (
	"github.com/cryptomines-online/backend/internal/combat/tilemap"
)

// WeaponClass classifies a stack's primary weapon for range, AI bias, and
// AoE behavior. Mirrors the GDD §8.6.4 weapon table.
type WeaponClass string

const (
	WeaponBallistic   WeaponClass = "ballistic"
	WeaponDirectional WeaponClass = "directional"
	WeaponMissile     WeaponClass = "missile"
	WeaponShipBased   WeaponClass = "ship_based"
	WeaponPlanetary   WeaponClass = "planetary" // defense buildings (single-target)
	WeaponPlanetaryAOE WeaponClass = "planetary_aoe" // Anti-Aircraft Gun: 3x3 AoE
)

// DamageType is the elemental class of a weapon's damage.
type DamageType string

const (
	DamageKinetic   DamageType = "kinetic"
	DamageHeat      DamageType = "heat"
	DamageExplosive DamageType = "explosive"
	DamageMagnetic  DamageType = "magnetic"
)

// ArmorType is the protective class of a hull.
type ArmorType string

const (
	ArmorNone         ArmorType = ""
	ArmorChrome       ArmorType = "chrome"
	ArmorRegen        ArmorType = "regen"
	ArmorNano         ArmorType = "nano"
	ArmorNeutralizing ArmorType = "neutralizing"
)

// HullClass categorizes ships for the rock-paper-scissors triangle.
type HullClass string

const (
	HullFrigate    HullClass = "frigate"
	HullCruiser    HullClass = "cruiser"
	HullBattleship HullClass = "battleship"
	HullBuilding   HullClass = "building"
)

// Stack is a single phalanx slot inside a fleet — or, for defense
// buildings, a synthetic single-stack representation. Combat math
// operates per stack.
//
// The fields are deliberately a smaller, flat subset of the legacy
// combat.FleetStack — only the values the tile-based engine actually
// reads are kept. Loaders translate from the database row + tech bonuses
// into these fields once at battle start.
type Stack struct {
	ID            string
	GridRow       int    // 0..2 within the parent fleet's 3x3 (0 for buildings)
	GridCol       int    // 0..2 within the parent fleet's 3x3 (0 for buildings)
	ShipDesignID  string // empty for buildings
	HullClass     HullClass
	ShipCount     int    // current ships alive (or building HP for buildings)

	// Weapon profile
	WeaponClass   WeaponClass
	WeaponName    string
	DamageType    DamageType
	MinRange      int // tiles (Chebyshev)
	MaxRange      int // tiles (Chebyshev)
	Cooldown      int // rounds between volleys (0 = every round)
	CooldownLeft  int // rounds remaining; ready when 0
	BaseDamageMin int // per-hit damage range
	BaseDamageMax int

	// Per-ship combat stats (multiplied by ShipCount during phases)
	Accuracy   int       // per-ship accuracy contribution
	Dodge      int       // per-ship dodge contribution
	ShieldHP   int       // total shield HP for the stack (current)
	MaxShield  int       // max shield HP (for restoration)
	StructHP   int       // total structure HP for the stack (current)
	MaxStruct  int       // max structure HP (for restoration)
	ArmorType  ArmorType
	HitChance  float64 // base weapon hit chance (0..1); 0 → engine default 0.7

	// Movement (only used at unit level via MIN(stack.Movement))
	Movement int

	// Advanced combat properties (a subset of the legacy FleetStack;
	// expanded as needed). Zero values mean "no effect".
	ScatterPct      float64 // % of damage that splashes within target's phalanx
	ShieldPenChance float64 // chance to bypass shields entirely with PenetrationPct
	ShieldPenPct    float64 // % of damage that penetrates shields when triggered
}

// Phalanx is the per-fleet 3x3 layout of stacks (GDD §8.6.5). For defense
// buildings, the phalanx contains exactly one stack at (0,0).
type Phalanx struct {
	Stacks []*Stack // 1..9 elements; sparse layout (formation rules apply)
}

// UnitKind separates mobile fleets from stationary buildings.
type UnitKind string

const (
	UnitFleet    UnitKind = "fleet"
	UnitBuilding UnitKind = "building"
)

// Unit is a single combatant on the battle map. A fleet occupies one tile
// and contains a phalanx; a building keeps its multi-tile footprint and
// has a one-stack phalanx.
type Unit struct {
	ID         string
	Side       tilemap.Side
	Kind       UnitKind
	FleetID    string // for Kind=fleet
	BuildingID string // for Kind=building
	TypeName   string // e.g. "particle_cannon"

	Tile      tilemap.Position   // anchor; for buildings, top-left of footprint
	Footprint []tilemap.Position // tiles covered (1 for fleets, N for buildings)

	Phalanx     Phalanx
	Speed       int  // commander speed → initiative; buildings have a fixed value
	Movement    int  // tiles per round (fleet level = MIN(stack.Movement)); 0 for buildings
	TargetCmd   string // "max_attack" | "min_attack" | "max_durability" | "min_durability" | "closest" | "by_commander_rank"

	// Aggregate HP for buildings; for fleets, sum of stack ship counts.
	dead bool
}

// IsAlive reports whether any stack still has ships (or HP for buildings).
func (u *Unit) IsAlive() bool {
	if u.dead {
		return false
	}
	for _, s := range u.Phalanx.Stacks {
		if s != nil && s.ShipCount > 0 {
			return true
		}
	}
	return false
}

// MarkDead unconditionally flags the unit as eliminated.
func (u *Unit) MarkDead() {
	u.dead = true
}

// TotalShips sums alive ships across the phalanx. For buildings this is
// the building's HP.
func (u *Unit) TotalShips() int {
	total := 0
	for _, s := range u.Phalanx.Stacks {
		if s != nil {
			total += s.ShipCount
		}
	}
	return total
}

// TotalAttack approximates the unit's offensive output for AI biasing.
// Sum of (mid-damage × ships) across stacks.
func (u *Unit) TotalAttack() int {
	total := 0
	for _, s := range u.Phalanx.Stacks {
		if s == nil || s.ShipCount <= 0 {
			continue
		}
		mid := (s.BaseDamageMin + s.BaseDamageMax) / 2
		total += mid * s.ShipCount
	}
	return total
}

// TotalDurability approximates how much damage the unit can absorb (sum
// of ships × hp-per-ship). Currently used as a stand-in: detailed shield
// + structure math is delegated to the damage sub-phases.
func (u *Unit) TotalDurability() int {
	// Without per-stack shield/structure here, fall back on raw count.
	// AI biasing is approximate; the precise math runs in the fire step.
	return u.TotalShips()
}
