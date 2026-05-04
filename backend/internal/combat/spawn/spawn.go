// Package spawn determines where attacker and defender units enter the
// 20x20 battle map at the start of combat (GDD §8.6.2). Spawn rules:
//
//   - Attacker fleets enter from one of the 4 corners of the battle map,
//     selected by the galaxy attack vector (sign of Δx, Δy from attacker's
//     home to the defender's planet).
//   - Multiple attacker fleets fan out along the entry edge starting at
//     the corner.
//   - Defending fleets spawn on free tiles adjacent to the Space Station
//     building (or, if no Space Station exists, near the map center).
package spawn

import (
	"github.com/cryptomines-online/backend/internal/combat/tilemap"
)

// Corner identifies one of the 4 entry corners of the battle map.
type Corner string

const (
	CornerTopLeft     Corner = "top_left"     // (0,0)
	CornerTopRight    Corner = "top_right"    // (Width-1, 0)
	CornerBottomLeft  Corner = "bottom_left"  // (0, Height-1)
	CornerBottomRight Corner = "bottom_right" // (Width-1, Height-1)
)

// Anchor returns the tile position of a corner.
func (c Corner) Anchor() tilemap.Position {
	switch c {
	case CornerTopLeft:
		return tilemap.Position{Col: 0, Row: 0}
	case CornerTopRight:
		return tilemap.Position{Col: tilemap.Width - 1, Row: 0}
	case CornerBottomLeft:
		return tilemap.Position{Col: 0, Row: tilemap.Height - 1}
	case CornerBottomRight:
		return tilemap.Position{Col: tilemap.Width - 1, Row: tilemap.Height - 1}
	}
	return tilemap.Position{Col: 0, Row: 0}
}

// PickCorner selects the entry corner from the galaxy-space attack vector
// (defender minus attacker). Implements the table in GDD §8.6.2:
//
//	Δx ≥ 0, Δy ≥ 0 → top-left      (attacker is to the left and above target)
//	Δx < 0, Δy ≥ 0 → top-right
//	Δx ≥ 0, Δy < 0 → bottom-left
//	Δx < 0, Δy < 0 → bottom-right
//
// Rationale: the attacker arrives "from" their direction, which mirrors
// the sign of the vector. Identical positions default to top-left.
func PickCorner(attackerHomeX, attackerHomeY, defenderPlanetX, defenderPlanetY int) Corner {
	dx := defenderPlanetX - attackerHomeX
	dy := defenderPlanetY - attackerHomeY
	switch {
	case dx >= 0 && dy >= 0:
		return CornerTopLeft
	case dx < 0 && dy >= 0:
		return CornerTopRight
	case dx >= 0 && dy < 0:
		return CornerBottomLeft
	default:
		return CornerBottomRight
	}
}

// EdgeTiles returns up to N spawn tiles along the entry edge starting from
// the corner anchor. The walk direction follows the corner's primary edge
// (top corners walk right/left along row 0; bottom corners walk along the
// last row). When N exceeds the edge length the walk wraps inward by 1
// row, then 2 rows, etc.
//
// Used to fan multiple attacker fleets across the entry edge.
func EdgeTiles(c Corner, n int) []tilemap.Position {
	if n <= 0 {
		return nil
	}
	anchor := c.Anchor()
	out := make([]tilemap.Position, 0, n)
	dCol, dRow := edgeWalkDelta(c)

	// Walk along the primary edge first.
	for i := 0; len(out) < n && i < tilemap.Width; i++ {
		p := tilemap.Position{Col: anchor.Col + i*dCol, Row: anchor.Row}
		if p.InBounds() {
			out = append(out, p)
		}
	}
	// Then walk inward (further rows from the edge), repeating along
	// the row direction. This is rarely reached — most battles have <20
	// attacker fleets — but defends against edge cases.
	for inward := 1; len(out) < n && inward < tilemap.Height; inward++ {
		row := anchor.Row + inward*dRow
		for i := 0; len(out) < n && i < tilemap.Width; i++ {
			p := tilemap.Position{Col: anchor.Col + i*dCol, Row: row}
			if p.InBounds() {
				out = append(out, p)
			}
		}
	}
	return out
}

func edgeWalkDelta(c Corner) (int, int) {
	switch c {
	case CornerTopLeft:
		return 1, 1 // walk right, then down for inward rows
	case CornerTopRight:
		return -1, 1
	case CornerBottomLeft:
		return 1, -1
	case CornerBottomRight:
		return -1, -1
	}
	return 1, 1
}

// DefenderSpawnTiles returns up to N tiles adjacent to the defender's
// Space Station footprint that are free of terrain and other units. If
// stationAnchor is the zero value (no Space Station present), spawns
// expand from the map center.
//
// world is consulted to skip blocked tiles. If fewer than N free tiles
// exist, returns whatever was found.
func DefenderSpawnTiles(world tilemap.Blockable, stationAnchor tilemap.Position, stationCols, stationRows, n int) []tilemap.Position {
	if n <= 0 {
		return nil
	}
	out := make([]tilemap.Position, 0, n)

	center := tilemap.Position{Col: tilemap.Width / 2, Row: tilemap.Height / 2}
	if stationCols > 0 && stationRows > 0 {
		center = tilemap.Position{
			Col: stationAnchor.Col + stationCols/2,
			Row: stationAnchor.Row + stationRows/2,
		}
	}

	// Concentric Chebyshev rings from radius 1 outward.
	visited := map[tilemap.Position]bool{}
	for radius := 1; len(out) < n && radius < tilemap.Width; radius++ {
		for dc := -radius; dc <= radius; dc++ {
			for dr := -radius; dr <= radius; dr++ {
				// Only the outer ring of this radius.
				if absInt(dc) != radius && absInt(dr) != radius {
					continue
				}
				p := tilemap.Position{Col: center.Col + dc, Row: center.Row + dr}
				if !p.InBounds() || visited[p] {
					continue
				}
				visited[p] = true
				if world.IsBlocked(p, "") {
					continue
				}
				out = append(out, p)
				if len(out) >= n {
					return out
				}
			}
		}
	}
	return out
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
