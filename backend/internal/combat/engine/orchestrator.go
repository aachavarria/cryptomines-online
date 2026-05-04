package engine

import (
	"fmt"

	"github.com/cryptomines-online/backend/internal/combat/spawn"
	"github.com/cryptomines-online/backend/internal/combat/tilemap"
	"github.com/cryptomines-online/backend/internal/services"
)

// SetupParams describes one battle setup. The orchestrator (RunPvP)
// loads everything from the database and produces a ready-to-execute
// Battle.
type SetupParams struct {
	AttackerID            string
	DefenderID            string
	DefenderPlanetID      string
	AttackerFleetIDs      []string
	AttackerHomePlanetID  string
	AttackerGalaxyX       int
	AttackerGalaxyY       int
	DefenderGalaxyX       int
	DefenderGalaxyY       int
	DefenderTechBonuses   *services.TechBonuses
	AttackerTechBonuses   *services.TechBonuses
	IncludeDefendingFleets []string // optional fleets owned by defender that defend the planet
}

// Setup loads everything needed for a tile-based PvP battle and returns
// a primed *Battle. Caller invokes Battle.Execute() to resolve.
//
// Steps:
//   1. New 20×20 BattleMap.
//   2. Load Meteor Star terrain → BattleMap.AddTerrain.
//   3. Load defender's defense buildings → Units, register on map.
//   4. Load defender's defending fleets → Units, spawn near Space Station.
//   5. Load attacker's fleets → Units, spawn from corner picked by
//      galaxy attack vector.
//   6. Construct *Battle with the populated unit list.
func Setup(p SetupParams, seed int64) (*Battle, error) {
	world := tilemap.NewBattleMap()

	// Terrain (Meteor Stars).
	terrain, err := LoadTerrain(p.DefenderPlanetID)
	if err != nil {
		return nil, fmt.Errorf("setup terrain: %w", err)
	}
	for _, t := range terrain {
		_ = world.AddTerrain(t)
	}

	units := []*Unit{}

	// Defense buildings as stationary units.
	defenseUnits, err := LoadDefenseUnits(p.DefenderPlanetID, p.DefenderTechBonuses)
	if err != nil {
		return nil, fmt.Errorf("setup defense: %w", err)
	}
	for _, u := range defenseUnits {
		if err := world.SetUnitTiles(u.ID, u.Footprint); err != nil {
			// If a building's footprint conflicts (e.g. terrain on a tile),
			// skip it — the worker will surface the issue in logs.
			continue
		}
		units = append(units, u)
	}

	// Defending fleets (optional). Spawn near Space Station.
	stationAnchor, sCols, sRows, hasStation := SpaceStationAnchor(p.DefenderPlanetID)
	defSpawnTiles := spawn.DefenderSpawnTiles(world, stationAnchor, sCols, sRows, len(p.IncludeDefendingFleets))
	for i, fid := range p.IncludeDefendingFleets {
		u, err := LoadFleetUnit(fid, p.DefenderID, tilemap.SideDefender)
		if err != nil {
			continue
		}
		if i >= len(defSpawnTiles) {
			break
		}
		u.Tile = defSpawnTiles[i]
		u.Footprint = []tilemap.Position{u.Tile}
		if err := world.SetUnitTiles(u.ID, u.Footprint); err != nil {
			continue
		}
		units = append(units, u)
	}
	_ = hasStation

	// Attacker fleets entering from the corner.
	corner := spawn.PickCorner(p.AttackerGalaxyX, p.AttackerGalaxyY, p.DefenderGalaxyX, p.DefenderGalaxyY)
	atkTiles := spawn.EdgeTiles(corner, len(p.AttackerFleetIDs)*2)
	atkIdx := 0
	for _, fid := range p.AttackerFleetIDs {
		u, err := LoadFleetUnit(fid, p.AttackerID, tilemap.SideAttacker)
		if err != nil {
			continue
		}
		// Find the next free tile in atkTiles.
		var placed bool
		for ; atkIdx < len(atkTiles); atkIdx++ {
			t := atkTiles[atkIdx]
			if !world.IsBlocked(t, "") {
				u.Tile = t
				u.Footprint = []tilemap.Position{t}
				if err := world.SetUnitTiles(u.ID, u.Footprint); err != nil {
					continue
				}
				atkIdx++
				placed = true
				break
			}
		}
		if !placed {
			continue
		}
		units = append(units, u)
	}

	if len(units) == 0 {
		return nil, fmt.Errorf("setup: no units placed on battle map")
	}

	return NewBattle(world, units, seed, 0), nil
}
