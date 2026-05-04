package pathing

import (
	"testing"

	"github.com/cryptomines-online/backend/internal/combat/tilemap"
)

func TestFindPath_StraightLine(t *testing.T) {
	m := tilemap.NewBattleMap()
	path := FindPath(m, "u1", tilemap.Position{0, 0}, tilemap.Position{4, 4}, -1)
	if path == nil {
		t.Fatal("expected path on empty map")
	}
	// Chebyshev distance = 4 → 5 nodes total (start + 4 steps).
	if len(path) != 5 {
		t.Errorf("expected 5 tiles in path, got %d: %v", len(path), path)
	}
	if !path[0].Equals(tilemap.Position{0, 0}) || !path[len(path)-1].Equals(tilemap.Position{4, 4}) {
		t.Errorf("path endpoints wrong: %v", path)
	}
}

func TestFindPath_AroundMeteorStar(t *testing.T) {
	m := tilemap.NewBattleMap()
	// Block the diagonal between (0,0) and (4,4) with a meteor at (2,2).
	_ = m.AddTerrain(&tilemap.Terrain{
		Kind: tilemap.TerrainMeteorStar, Position: tilemap.Position{2, 2}, HP: 100, MaxHP: 100,
	})
	path := FindPath(m, "u1", tilemap.Position{0, 0}, tilemap.Position{4, 4}, -1)
	if path == nil {
		t.Fatal("expected path around meteor")
	}
	for _, p := range path {
		if p.Equals(tilemap.Position{2, 2}) {
			t.Errorf("path passes through blocked tile: %v", path)
		}
	}
}

func TestFindPath_TruncatesToMaxSteps(t *testing.T) {
	m := tilemap.NewBattleMap()
	path := FindPath(m, "u1", tilemap.Position{0, 0}, tilemap.Position{10, 10}, 3)
	if path == nil {
		t.Fatal("expected truncated path")
	}
	if len(path)-1 != 3 {
		t.Errorf("expected 3 steps, got %d steps (path %v)", len(path)-1, path)
	}
}

func TestFindPath_GoalOccupied_ApproachesAdjacent(t *testing.T) {
	m := tilemap.NewBattleMap()
	_ = m.SetUnitTiles("enemy", []tilemap.Position{{5, 5}})
	path := FindPath(m, "u1", tilemap.Position{0, 0}, tilemap.Position{5, 5}, -1)
	if path == nil {
		t.Fatal("expected approach path")
	}
	last := path[len(path)-1]
	if last.Equals(tilemap.Position{5, 5}) {
		t.Errorf("path ended on occupied goal tile: %v", path)
	}
	if tilemap.Chebyshev(last, tilemap.Position{5, 5}) != 1 {
		t.Errorf("expected path to end adjacent to goal, ended at %v", last)
	}
}

func TestFindPath_NoPathReturnsNil(t *testing.T) {
	m := tilemap.NewBattleMap()
	// Box (0,0) in completely with meteors.
	for _, p := range []tilemap.Position{{0, 1}, {1, 0}, {1, 1}} {
		_ = m.AddTerrain(&tilemap.Terrain{
			Kind: tilemap.TerrainMeteorStar, Position: p, HP: 100, MaxHP: 100,
		})
	}
	path := FindPath(m, "u1", tilemap.Position{0, 0}, tilemap.Position{10, 10}, -1)
	if path != nil {
		t.Errorf("expected nil path when boxed in, got %v", path)
	}
}

func TestFindPath_StartEqualsGoal(t *testing.T) {
	m := tilemap.NewBattleMap()
	path := FindPath(m, "u1", tilemap.Position{5, 5}, tilemap.Position{5, 5}, -1)
	if len(path) != 1 || !path[0].Equals(tilemap.Position{5, 5}) {
		t.Errorf("start==goal path should be [start], got %v", path)
	}
}
