package engine

import (
	"github.com/cryptomines-online/backend/internal/combat/tilemap"
)

// chooseTarget picks an enemy unit for the attacker to fire on this round,
// among those within the attacker's effective weapon range. Returns nil
// when no enemy is in range or all enemies are dead.
//
// Targeting precedence (GDD §8.6.11):
//  1. Manual targeting_command set on the attacking fleet (player's order).
//  2. Otherwise, weapon-class default AI bias.
//
// Defense buildings always use their weapon-class default (they have no
// commander-issued targeting_command).
func chooseTarget(attacker *Unit, candidates []*Unit, world *tilemap.BattleMap) *Unit {
	primary := primaryStack(attacker)
	if primary == nil {
		return nil
	}
	inRange := []*Unit{}
	for _, c := range candidates {
		if c == nil || !c.IsAlive() || c.Side == attacker.Side {
			continue
		}
		if !weaponInRange(attacker, c, primary, world) {
			continue
		}
		inRange = append(inRange, c)
	}
	if len(inRange) == 0 {
		return nil
	}
	cmd := attacker.TargetCmd
	if attacker.Kind == UnitBuilding || cmd == "" {
		cmd = defaultCommandForWeapon(primary.WeaponClass)
	}
	return applyCommand(attacker, inRange, cmd)
}

func primaryStack(u *Unit) *Stack {
	for _, s := range u.Phalanx.Stacks {
		if s != nil && s.ShipCount > 0 {
			return s
		}
	}
	return nil
}

// weaponInRange reports whether attacker can fire on target this round.
// Tests Chebyshev distance, weapon range bounds, and line-of-sight (LOS
// is blocked only by Meteor Stars per GDD §8.6.4).
func weaponInRange(attacker, target *Unit, stack *Stack, world *tilemap.BattleMap) bool {
	if stack.CooldownLeft > 0 {
		return false
	}
	d := tilemap.Chebyshev(attacker.Tile, target.Tile)
	if d < stack.MinRange || d > stack.MaxRange {
		return false
	}
	return world.HasLineOfSight(attacker.Tile, target.Tile)
}

func defaultCommandForWeapon(w WeaponClass) string {
	switch w {
	case WeaponMissile:
		return "max_attack"
	case WeaponShipBased:
		return "max_durability"
	default: // ballistic, directional, planetary, planetary_aoe
		return "closest"
	}
}

func applyCommand(attacker *Unit, candidates []*Unit, cmd string) *Unit {
	if len(candidates) == 0 {
		return nil
	}
	best := candidates[0]
	score := scoreFor(attacker, best, cmd)
	for _, c := range candidates[1:] {
		s := scoreFor(attacker, c, cmd)
		if s > score {
			best = c
			score = s
		}
	}
	return best
}

func scoreFor(attacker, target *Unit, cmd string) int {
	switch cmd {
	case "max_attack":
		return target.TotalAttack()
	case "min_attack":
		return -target.TotalAttack()
	case "max_durability":
		return target.TotalDurability()
	case "min_durability":
		return -target.TotalDurability()
	case "closest":
		return -tilemap.Chebyshev(attacker.Tile, target.Tile)
	case "by_commander_rank":
		return target.Speed
	default:
		return -tilemap.Chebyshev(attacker.Tile, target.Tile)
	}
}

// pickDestination decides the tile a fleet wants to be on at the end of
// its movement this round (GDD §8.6.10). Heuristic, in priority order:
//  1. If no enemy is within max_range of any candidate stack → close on
//     the nearest enemy.
//  2. If our weapon has min_range > 1 (Missile, Ship-Based) and an enemy
//     is too close → kite to a tile at >= min_range.
//  3. Otherwise, hold position.
//
// Returns the unit's current tile when no movement is preferred.
func pickDestination(unit *Unit, enemies []*Unit, _ *tilemap.BattleMap) tilemap.Position {
	primary := primaryStack(unit)
	if primary == nil || unit.Movement <= 0 || len(enemies) == 0 {
		return unit.Tile
	}

	nearest, nearestDist := nearestEnemy(unit, enemies)
	if nearest == nil {
		return unit.Tile
	}

	// 1. Out of range entirely → close in.
	if nearestDist > primary.MaxRange {
		return moveToward(unit.Tile, nearest.Tile, unit.Movement)
	}

	// 2. Inside min_range → step back. Choose a tile at exactly min_range
	//    in the direction away from the threat.
	if primary.MinRange > 1 && nearestDist < primary.MinRange {
		return moveAway(unit.Tile, nearest.Tile, primary.MinRange-nearestDist)
	}

	// 3. Already in a good firing slot.
	return unit.Tile
}

func nearestEnemy(unit *Unit, enemies []*Unit) (*Unit, int) {
	var best *Unit
	bestDist := -1
	for _, e := range enemies {
		if e == nil || !e.IsAlive() || e.Side == unit.Side {
			continue
		}
		d := tilemap.Chebyshev(unit.Tile, e.Tile)
		if bestDist < 0 || d < bestDist {
			best = e
			bestDist = d
		}
	}
	return best, bestDist
}

// moveToward returns the position `steps` tiles closer to target, clamped
// to map bounds. The actual A* path is computed later — this is the
// destination hint.
func moveToward(from, target tilemap.Position, steps int) tilemap.Position {
	dx := signOf(target.Col - from.Col)
	dy := signOf(target.Row - from.Row)
	dist := tilemap.Chebyshev(from, target)
	if dist <= steps {
		// Try to land adjacent to target (1-tile gap).
		return tilemap.Position{Col: target.Col - dx, Row: target.Row - dy}
	}
	return tilemap.Position{
		Col: clamp(from.Col+dx*steps, 0, tilemap.Width-1),
		Row: clamp(from.Row+dy*steps, 0, tilemap.Height-1),
	}
}

func moveAway(from, threat tilemap.Position, steps int) tilemap.Position {
	dx := signOf(from.Col - threat.Col)
	dy := signOf(from.Row - threat.Row)
	if dx == 0 && dy == 0 {
		dx = 1 // arbitrary direction when stacked on the threat
	}
	return tilemap.Position{
		Col: clamp(from.Col+dx*steps, 0, tilemap.Width-1),
		Row: clamp(from.Row+dy*steps, 0, tilemap.Height-1),
	}
}

func signOf(v int) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	default:
		return 0
	}
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
