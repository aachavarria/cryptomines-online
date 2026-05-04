package spawn

import (
	"testing"

	"github.com/cryptomines-online/backend/internal/combat/tilemap"
)

func TestPickCorner_Quadrants(t *testing.T) {
	cases := []struct {
		name                                          string
		ax, ay, dx, dy                                int
		want                                          Corner
	}{
		{"target east-south of attacker → top-left", 100, 100, 200, 200, CornerTopLeft},
		{"target west-south of attacker → top-right", 200, 100, 100, 200, CornerTopRight},
		{"target east-north of attacker → bottom-left", 100, 200, 200, 100, CornerBottomLeft},
		{"target west-north of attacker → bottom-right", 200, 200, 100, 100, CornerBottomRight},
		{"identical positions default to top-left", 100, 100, 100, 100, CornerTopLeft},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := PickCorner(c.ax, c.ay, c.dx, c.dy)
			if got != c.want {
				t.Errorf("PickCorner(%d,%d,%d,%d) = %s, want %s",
					c.ax, c.ay, c.dx, c.dy, got, c.want)
			}
		})
	}
}

func TestCorner_Anchor(t *testing.T) {
	cases := []struct {
		c    Corner
		want tilemap.Position
	}{
		{CornerTopLeft, tilemap.Position{0, 0}},
		{CornerTopRight, tilemap.Position{19, 0}},
		{CornerBottomLeft, tilemap.Position{0, 19}},
		{CornerBottomRight, tilemap.Position{19, 19}},
	}
	for _, c := range cases {
		if got := c.c.Anchor(); !got.Equals(c.want) {
			t.Errorf("%s.Anchor() = %s, want %s", c.c, got, c.want)
		}
	}
}

func TestEdgeTiles_FromTopLeft(t *testing.T) {
	tiles := EdgeTiles(CornerTopLeft, 5)
	if len(tiles) != 5 {
		t.Fatalf("expected 5 tiles, got %d", len(tiles))
	}
	for i, p := range tiles {
		if !p.Equals(tilemap.Position{Col: i, Row: 0}) {
			t.Errorf("tile %d = %s, want (%d,0)", i, p, i)
		}
	}
}

func TestEdgeTiles_WrapsInward(t *testing.T) {
	// Asking for 25 tiles (more than width 20) forces an inward step.
	tiles := EdgeTiles(CornerTopLeft, 25)
	if len(tiles) != 25 {
		t.Fatalf("expected 25 tiles, got %d", len(tiles))
	}
	if !tiles[20].Equals(tilemap.Position{Col: 0, Row: 1}) {
		t.Errorf("tile 20 (first inward) = %s, want (0,1)", tiles[20])
	}
}

func TestDefenderSpawnTiles_AdjacentToStation(t *testing.T) {
	m := tilemap.NewBattleMap()
	// Pretend Space Station occupies (10,10)..(12,12).
	stationAnchor := tilemap.Position{10, 10}
	if err := m.SetUnitTiles("station", tilemap.Footprint(stationAnchor, 3, 3)); err != nil {
		t.Fatalf("station setup: %v", err)
	}
	tiles := DefenderSpawnTiles(m, stationAnchor, 3, 3, 5)
	if len(tiles) != 5 {
		t.Fatalf("expected 5 spawn tiles, got %d", len(tiles))
	}
	for _, p := range tiles {
		if m.UnitAt(p) != "" {
			t.Errorf("spawn tile %s overlaps station", p)
		}
	}
}
