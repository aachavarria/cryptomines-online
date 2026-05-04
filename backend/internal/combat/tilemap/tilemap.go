// Package tilemap provides core spatial primitives for tile-based combat
// (GDD §8.6). The battle map is a 20x20 grid that mirrors the defender's
// Space Base layout. All combat positioning, movement, and weapon range
// calculations operate on this grid.
package tilemap

import "fmt"

// Width and Height of the battle map. Matches the Space Base build grid
// (GDD §2.2.9). Buildings keep their build-time tile positions when combat
// starts; attacking fleets enter from a corner and traverse tiles per round.
const (
	Width  = 20
	Height = 20
)

// Blockable is the read-only interface the pathfinder and spawn helpers
// use to ask "is this tile blocked for this unit?". BattleMap implements
// this. Callers may pass "" as the unit ID for occupancy-only checks.
type Blockable interface {
	IsBlocked(p Position, movingUnitID string) bool
}

// Side identifies which faction a unit belongs to.
type Side string

const (
	SideAttacker Side = "attacker"
	SideDefender Side = "defender"
)

// Position is a tile coordinate on the battle map.
type Position struct {
	Col int `json:"col"`
	Row int `json:"row"`
}

// InBounds reports whether the position lies inside the battle map.
func (p Position) InBounds() bool {
	return p.Col >= 0 && p.Col < Width && p.Row >= 0 && p.Row < Height
}

// Equals returns true when two positions point at the same tile.
func (p Position) Equals(o Position) bool {
	return p.Col == o.Col && p.Row == o.Row
}

// String renders a position as "(col,row)".
func (p Position) String() string {
	return fmt.Sprintf("(%d,%d)", p.Col, p.Row)
}

// Chebyshev distance between two tiles. GO2 uses 8-directional movement
// (king moves), so distance = max(|dx|, |dy|). This is the metric used for
// weapon range and movement budgets.
func Chebyshev(a, b Position) int {
	dx := a.Col - b.Col
	if dx < 0 {
		dx = -dx
	}
	dy := a.Row - b.Row
	if dy < 0 {
		dy = -dy
	}
	if dx > dy {
		return dx
	}
	return dy
}

// Manhattan distance is provided for callers that need it (rare — the engine
// itself uses Chebyshev). |dx| + |dy|.
func Manhattan(a, b Position) int {
	dx := a.Col - b.Col
	if dx < 0 {
		dx = -dx
	}
	dy := a.Row - b.Row
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

// Neighbors8 returns the up-to-8 in-bounds tiles surrounding p (Chebyshev
// neighborhood). Used for AoE patterns (Anti-Aircraft Gun) and pathing.
func Neighbors8(p Position) []Position {
	out := make([]Position, 0, 8)
	for dc := -1; dc <= 1; dc++ {
		for dr := -1; dr <= 1; dr++ {
			if dc == 0 && dr == 0 {
				continue
			}
			n := Position{Col: p.Col + dc, Row: p.Row + dr}
			if n.InBounds() {
				out = append(out, n)
			}
		}
	}
	return out
}

// Footprint enumerates every tile a multi-tile entity covers given its
// anchor (top-left) position and size. Mirrors GDD §2.2.9 building
// footprints. Tiles outside the map are silently dropped.
func Footprint(anchor Position, cols, rows int) []Position {
	out := make([]Position, 0, cols*rows)
	for dc := 0; dc < cols; dc++ {
		for dr := 0; dr < rows; dr++ {
			t := Position{Col: anchor.Col + dc, Row: anchor.Row + dr}
			if t.InBounds() {
				out = append(out, t)
			}
		}
	}
	return out
}

// TerrainKind classifies the obstacle at a tile.
type TerrainKind string

const (
	TerrainEmpty      TerrainKind = ""
	TerrainMeteorStar TerrainKind = "meteor_star"
)

// Terrain is a stationary obstacle on the battle map. Meteor Stars block
// movement entirely (GDD §8.6.2) and have HP — attackers must destroy them
// or path around. Other terrain kinds may be added later.
type Terrain struct {
	Kind     TerrainKind `json:"kind"`
	Position Position    `json:"position"`
	HP       int         `json:"hp"`
	MaxHP    int         `json:"max_hp"`
}

// Blocks reports whether this terrain blocks movement.
func (t Terrain) Blocks() bool {
	if t.HP <= 0 {
		return false
	}
	switch t.Kind {
	case TerrainMeteorStar:
		return true
	default:
		return false
	}
}

// BattleMap is the spatial state of an in-progress combat. Unit and terrain
// occupancy are tracked by tile so the engine can answer:
//   - "Is this tile blocked for movement?"
//   - "What unit (if any) sits on this tile?"
//   - "Is there line-of-sight from A to B?"
//
// The map is owned by a single combat run and mutated as units move and die.
type BattleMap struct {
	// terrain indexed by tile key "col,row"
	terrain map[string]*Terrain
	// unit occupancy indexed by tile key. Multi-tile defense buildings
	// register every tile of their footprint so any tile lookup resolves.
	occupancy map[string]string // tile key -> unit ID
}

// NewBattleMap creates an empty battle map.
func NewBattleMap() *BattleMap {
	return &BattleMap{
		terrain:   make(map[string]*Terrain),
		occupancy: make(map[string]string),
	}
}

func tileKey(p Position) string {
	return fmt.Sprintf("%d,%d", p.Col, p.Row)
}

// AddTerrain places terrain on the map. Returns an error if the position is
// out of bounds or already occupied by other terrain.
func (m *BattleMap) AddTerrain(t *Terrain) error {
	if !t.Position.InBounds() {
		return fmt.Errorf("terrain position %s out of bounds", t.Position)
	}
	k := tileKey(t.Position)
	if _, exists := m.terrain[k]; exists {
		return fmt.Errorf("terrain already at %s", t.Position)
	}
	m.terrain[k] = t
	return nil
}

// TerrainAt returns the terrain on a tile, or nil if empty.
func (m *BattleMap) TerrainAt(p Position) *Terrain {
	return m.terrain[tileKey(p)]
}

// AllTerrain returns a snapshot of every terrain entity on the map.
func (m *BattleMap) AllTerrain() []*Terrain {
	out := make([]*Terrain, 0, len(m.terrain))
	for _, t := range m.terrain {
		out = append(out, t)
	}
	return out
}

// DamageTerrain applies HP damage to terrain at the given tile. If HP drops
// to 0, the terrain is removed and the tile becomes passable. Returns true
// when the terrain was destroyed, false otherwise (including when no
// terrain existed there).
func (m *BattleMap) DamageTerrain(p Position, dmg int) bool {
	t := m.TerrainAt(p)
	if t == nil {
		return false
	}
	t.HP -= dmg
	if t.HP <= 0 {
		t.HP = 0
		delete(m.terrain, tileKey(p))
		return true
	}
	return false
}

// SetUnitTiles registers a unit as occupying every tile in its footprint.
// Use a single-element slice for fleets (1 tile) and the full footprint for
// multi-tile defense buildings. Returns an error if any tile is already
// occupied by another unit or by blocking terrain.
func (m *BattleMap) SetUnitTiles(unitID string, tiles []Position) error {
	for _, p := range tiles {
		if !p.InBounds() {
			return fmt.Errorf("unit %s tile %s out of bounds", unitID, p)
		}
		if t := m.TerrainAt(p); t != nil && t.Blocks() {
			return fmt.Errorf("unit %s tile %s blocked by terrain", unitID, p)
		}
		k := tileKey(p)
		if existing, taken := m.occupancy[k]; taken && existing != unitID {
			return fmt.Errorf("unit %s tile %s already occupied by %s", unitID, p, existing)
		}
	}
	for _, p := range tiles {
		m.occupancy[tileKey(p)] = unitID
	}
	return nil
}

// ClearUnit removes a unit's tile occupancy entirely. Used when a unit is
// destroyed.
func (m *BattleMap) ClearUnit(unitID string) {
	for k, owner := range m.occupancy {
		if owner == unitID {
			delete(m.occupancy, k)
		}
	}
}

// MoveUnit relocates a single-tile unit from one position to another. The
// destination must be in bounds, free of blocking terrain, and either empty
// or already owned by this unit. Returns an error otherwise.
func (m *BattleMap) MoveUnit(unitID string, from, to Position) error {
	if !to.InBounds() {
		return fmt.Errorf("move %s: destination %s out of bounds", unitID, to)
	}
	if t := m.TerrainAt(to); t != nil && t.Blocks() {
		return fmt.Errorf("move %s: destination %s blocked by terrain", unitID, to)
	}
	dk := tileKey(to)
	if owner, taken := m.occupancy[dk]; taken && owner != unitID {
		return fmt.Errorf("move %s: destination %s occupied by %s", unitID, to, owner)
	}
	delete(m.occupancy, tileKey(from))
	m.occupancy[dk] = unitID
	return nil
}

// UnitAt returns the unit ID on a tile, or "" if empty.
func (m *BattleMap) UnitAt(p Position) string {
	return m.occupancy[tileKey(p)]
}

// IsBlocked reports whether a tile cannot be entered by the given unit. A
// tile is blocked if it is out of bounds, has blocking terrain, or is
// occupied by a different unit. The unit's own tiles are not blocked for
// itself (so move-in-place is a no-op).
func (m *BattleMap) IsBlocked(p Position, movingUnitID string) bool {
	if !p.InBounds() {
		return true
	}
	if t := m.TerrainAt(p); t != nil && t.Blocks() {
		return true
	}
	if owner, taken := m.occupancy[tileKey(p)]; taken && owner != movingUnitID {
		return true
	}
	return false
}

// HasLineOfSight reports whether a straight Bresenham line from a to b is
// free of blocking terrain. Unit occupancy does NOT block LOS — only
// terrain (Meteor Stars). Endpoints themselves are not tested.
func (m *BattleMap) HasLineOfSight(a, b Position) bool {
	for _, p := range bresenham(a, b) {
		if p.Equals(a) || p.Equals(b) {
			continue
		}
		if t := m.TerrainAt(p); t != nil && t.Blocks() {
			return false
		}
	}
	return true
}

// bresenham returns the tiles on a straight line from a to b inclusive.
func bresenham(a, b Position) []Position {
	dx := b.Col - a.Col
	if dx < 0 {
		dx = -dx
	}
	dy := b.Row - a.Row
	if dy < 0 {
		dy = -dy
	}
	sx, sy := 1, 1
	if a.Col > b.Col {
		sx = -1
	}
	if a.Row > b.Row {
		sy = -1
	}
	err := dx - dy
	cur := a
	out := []Position{cur}
	for !cur.Equals(b) {
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			cur.Col += sx
		}
		if e2 < dx {
			err += dx
			cur.Row += sy
		}
		out = append(out, cur)
	}
	return out
}
