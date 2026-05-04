package engine

import (
	"testing"

	"github.com/cryptomines-online/backend/internal/combat/tilemap"
)

// fleet builds a single-stack fleet unit for testing.
func fleet(id string, side tilemap.Side, tile tilemap.Position, ships, dmgMin, dmgMax, minR, maxR, mov, speed int, weapon WeaponClass) *Unit {
	return &Unit{
		ID: id, Side: side, Kind: UnitFleet, FleetID: id,
		Tile: tile, Footprint: []tilemap.Position{tile},
		Speed: speed, Movement: mov,
		Phalanx: Phalanx{Stacks: []*Stack{{
			ID: id + "_s0", GridRow: 0, GridCol: 0,
			ShipDesignID: "test", HullClass: HullFrigate,
			ShipCount: ships,
			WeaponClass: weapon, WeaponName: string(weapon),
			DamageType: DamageKinetic,
			MinRange: minR, MaxRange: maxR,
			BaseDamageMin: dmgMin, BaseDamageMax: dmgMax,
			HitChance: 0.95,
			ShieldHP:  ships * 10, MaxShield: ships * 10,
			StructHP:  ships * 50, MaxStruct: ships * 50,
		}}},
	}
}

func building(id string, tile tilemap.Position, hp, dmgMin, dmgMax, minR, maxR, speed int, weapon WeaponClass) *Unit {
	return &Unit{
		ID: id, Side: tilemap.SideDefender, Kind: UnitBuilding, BuildingID: id,
		TypeName: "particle_cannon",
		Tile: tile, Footprint: []tilemap.Position{tile},
		Speed: speed, Movement: 0,
		Phalanx: Phalanx{Stacks: []*Stack{{
			ID: id + "_s0", GridRow: 0, GridCol: 0,
			HullClass: HullBuilding,
			ShipCount: hp,
			WeaponClass: weapon, WeaponName: string(weapon),
			DamageType: DamageKinetic,
			MinRange: minR, MaxRange: maxR,
			BaseDamageMin: dmgMin, BaseDamageMax: dmgMax,
			HitChance: 0.95,
			ShieldHP:  0, MaxShield: 0,
			StructHP:  hp, MaxStruct: hp,
		}}},
	}
}

func TestEngine_AttackerWipesLoneDefender(t *testing.T) {
	world := tilemap.NewBattleMap()
	atk := fleet("atk", tilemap.SideAttacker, tilemap.Position{0, 0}, 100, 50, 50, 1, 5, 3, 10, WeaponBallistic)
	def := building("def", tilemap.Position{5, 5}, 200, 10, 10, 1, 3, 5, WeaponPlanetary)

	if err := world.SetUnitTiles(atk.ID, atk.Footprint); err != nil {
		t.Fatal(err)
	}
	if err := world.SetUnitTiles(def.ID, def.Footprint); err != nil {
		t.Fatal(err)
	}

	b := NewBattle(world, []*Unit{atk, def}, 42, 50)
	rep := b.Execute()

	if rep.End.Winner != tilemap.SideAttacker {
		t.Errorf("expected attacker win, got %s (%s)", rep.End.Winner, rep.End.Reason)
	}
	if !atk.IsAlive() {
		t.Error("attacker should still be alive")
	}
	if def.IsAlive() {
		t.Error("defender should be dead")
	}
}

func TestEngine_DefenderHoldsAtRoundCap(t *testing.T) {
	world := tilemap.NewBattleMap()
	// Tiny attacker, indestructible defender (huge HP, fast cooldown).
	atk := fleet("atk", tilemap.SideAttacker, tilemap.Position{0, 0}, 5, 1, 1, 1, 2, 1, 5, WeaponBallistic)
	def := building("def", tilemap.Position{18, 18}, 100000, 100, 100, 1, 30, 5, WeaponPlanetaryAOE)
	_ = world.SetUnitTiles(atk.ID, atk.Footprint)
	_ = world.SetUnitTiles(def.ID, def.Footprint)

	b := NewBattle(world, []*Unit{atk, def}, 1, 25)
	rep := b.Execute()

	if rep.End.Winner != tilemap.SideDefender {
		t.Errorf("expected defender win on cap or wipe, got %s (%s)", rep.End.Winner, rep.End.Reason)
	}
}

func TestEngine_RecordsMovement(t *testing.T) {
	world := tilemap.NewBattleMap()
	atk := fleet("atk", tilemap.SideAttacker, tilemap.Position{0, 0}, 100, 50, 50, 1, 3, 3, 10, WeaponBallistic)
	def := building("def", tilemap.Position{19, 19}, 1000, 10, 10, 1, 1, 5, WeaponPlanetary)
	_ = world.SetUnitTiles(atk.ID, atk.Footprint)
	_ = world.SetUnitTiles(def.ID, def.Footprint)

	b := NewBattle(world, []*Unit{atk, def}, 99, 5)
	rep := b.Execute()

	movesSeen := 0
	for _, r := range rep.Rounds {
		for _, a := range r.Actions {
			if a.Type == "move" && a.UnitID == "atk" {
				movesSeen++
			}
		}
	}
	if movesSeen == 0 {
		t.Error("expected attacker to move at least once")
	}
}
