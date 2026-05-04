package tilemap

import "testing"

func TestPosition_InBounds(t *testing.T) {
	cases := []struct {
		p    Position
		want bool
	}{
		{Position{0, 0}, true},
		{Position{19, 19}, true},
		{Position{20, 0}, false},
		{Position{-1, 5}, false},
		{Position{5, 20}, false},
	}
	for _, c := range cases {
		if got := c.p.InBounds(); got != c.want {
			t.Errorf("%s.InBounds() = %v, want %v", c.p, got, c.want)
		}
	}
}

func TestChebyshev(t *testing.T) {
	cases := []struct {
		a, b Position
		want int
	}{
		{Position{0, 0}, Position{0, 0}, 0},
		{Position{0, 0}, Position{3, 0}, 3},
		{Position{0, 0}, Position{0, 5}, 5},
		{Position{0, 0}, Position{3, 5}, 5},
		{Position{2, 2}, Position{5, 7}, 5},
		{Position{10, 10}, Position{4, 6}, 6},
	}
	for _, c := range cases {
		if got := Chebyshev(c.a, c.b); got != c.want {
			t.Errorf("Chebyshev(%s,%s) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestNeighbors8_Center(t *testing.T) {
	got := Neighbors8(Position{5, 5})
	if len(got) != 8 {
		t.Fatalf("expected 8 neighbors at center, got %d", len(got))
	}
}

func TestNeighbors8_Corner(t *testing.T) {
	got := Neighbors8(Position{0, 0})
	if len(got) != 3 {
		t.Fatalf("expected 3 neighbors at (0,0), got %d", len(got))
	}
}

func TestFootprint_3x3(t *testing.T) {
	tiles := Footprint(Position{2, 3}, 3, 3)
	if len(tiles) != 9 {
		t.Fatalf("expected 9 tiles for 3x3 footprint, got %d", len(tiles))
	}
}

func TestFootprint_ClipsOutOfBounds(t *testing.T) {
	// Anchor near the edge, footprint extends past the map.
	tiles := Footprint(Position{18, 18}, 3, 3)
	for _, p := range tiles {
		if !p.InBounds() {
			t.Errorf("Footprint returned out-of-bounds tile %s", p)
		}
	}
	// 2x2 fits at (18,18); 3x3 should clip to 4 tiles.
	if len(tiles) != 4 {
		t.Errorf("expected 4 tiles after clipping at (18,18) 3x3, got %d", len(tiles))
	}
}

func TestBattleMap_TerrainAndOccupancy(t *testing.T) {
	m := NewBattleMap()
	ms := &Terrain{Kind: TerrainMeteorStar, Position: Position{5, 5}, HP: 1000, MaxHP: 1000}
	if err := m.AddTerrain(ms); err != nil {
		t.Fatalf("AddTerrain: %v", err)
	}
	// Cannot place a unit on top of a Meteor Star.
	if err := m.SetUnitTiles("u1", []Position{{5, 5}}); err == nil {
		t.Fatal("expected SetUnitTiles to fail on blocked tile")
	}
	// Place a unit on (4,4).
	if err := m.SetUnitTiles("u1", []Position{{4, 4}}); err != nil {
		t.Fatalf("SetUnitTiles: %v", err)
	}
	if got := m.UnitAt(Position{4, 4}); got != "u1" {
		t.Errorf("UnitAt(4,4) = %q, want u1", got)
	}
	// Move u1 to (4,5).
	if err := m.MoveUnit("u1", Position{4, 4}, Position{4, 5}); err != nil {
		t.Fatalf("MoveUnit: %v", err)
	}
	if m.UnitAt(Position{4, 4}) != "" || m.UnitAt(Position{4, 5}) != "u1" {
		t.Errorf("after MoveUnit, occupancy not updated: from=%q to=%q",
			m.UnitAt(Position{4, 4}), m.UnitAt(Position{4, 5}))
	}
}

func TestBattleMap_DamageTerrain(t *testing.T) {
	m := NewBattleMap()
	ms := &Terrain{Kind: TerrainMeteorStar, Position: Position{1, 1}, HP: 100, MaxHP: 100}
	_ = m.AddTerrain(ms)

	if destroyed := m.DamageTerrain(Position{1, 1}, 50); destroyed {
		t.Error("terrain should not be destroyed at half HP")
	}
	if destroyed := m.DamageTerrain(Position{1, 1}, 60); !destroyed {
		t.Error("terrain should be destroyed when HP <= 0")
	}
	if got := m.TerrainAt(Position{1, 1}); got != nil {
		t.Error("destroyed terrain should be removed from map")
	}
}

func TestBattleMap_LineOfSightBlockedByMeteor(t *testing.T) {
	m := NewBattleMap()
	_ = m.AddTerrain(&Terrain{Kind: TerrainMeteorStar, Position: Position{5, 5}, HP: 100, MaxHP: 100})

	if m.HasLineOfSight(Position{2, 2}, Position{8, 8}) {
		t.Error("LOS through Meteor Star should be blocked")
	}
	if !m.HasLineOfSight(Position{0, 0}, Position{0, 10}) {
		t.Error("clear LOS along column should pass")
	}
}

func TestBattleMap_IsBlocked(t *testing.T) {
	m := NewBattleMap()
	_ = m.AddTerrain(&Terrain{Kind: TerrainMeteorStar, Position: Position{3, 3}, HP: 100, MaxHP: 100})
	_ = m.SetUnitTiles("u1", []Position{{4, 4}})

	if !m.IsBlocked(Position{3, 3}, "u2") {
		t.Error("Meteor Star tile should be blocked")
	}
	if !m.IsBlocked(Position{4, 4}, "u2") {
		t.Error("tile occupied by other unit should be blocked")
	}
	if m.IsBlocked(Position{4, 4}, "u1") {
		t.Error("a unit should not be blocked from its own tile")
	}
	if !m.IsBlocked(Position{-1, 0}, "u2") {
		t.Error("out-of-bounds tile should be blocked")
	}
}
