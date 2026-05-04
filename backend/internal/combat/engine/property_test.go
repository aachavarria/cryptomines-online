package engine

import (
	"testing"

	"github.com/cryptomines-online/backend/internal/combat/tilemap"
)

// TestProp_BattleAlwaysTerminates runs many random small battles and
// verifies that every one returns a non-empty Report with a winner and
// a finite round count ≤ roundCap.
func TestProp_BattleAlwaysTerminates(t *testing.T) {
	const trials = 50
	for i := 0; i < trials; i++ {
		world := tilemap.NewBattleMap()
		atk := fleet("atk", tilemap.SideAttacker, tilemap.Position{0, 0}, 50+i, 100, 200, 1, 4, 2, 10, WeaponBallistic)
		def := building("def", tilemap.Position{10 + i%5, 10 + i%5}, 5000+i*100, 50, 200, 1, 5, 5, WeaponPlanetary)
		_ = world.SetUnitTiles(atk.ID, atk.Footprint)
		_ = world.SetUnitTiles(def.ID, def.Footprint)

		b := NewBattle(world, []*Unit{atk, def}, int64(i+1), 30)
		rep := b.Execute()

		if rep.End.Winner == "" {
			t.Fatalf("trial %d: no winner recorded", i)
		}
		if rep.Summary.TotalRounds <= 0 || rep.Summary.TotalRounds > 30 {
			t.Errorf("trial %d: total_rounds out of range: %d", i, rep.Summary.TotalRounds)
		}
		if len(rep.Rounds) == 0 {
			t.Errorf("trial %d: rounds slice empty", i)
		}
	}
}

// TestProp_AttackerWipedLosesEarly: when attacker has no chance, the
// battle ends with defender_win at or before roundCap.
func TestProp_AttackerWipedLosesEarly(t *testing.T) {
	world := tilemap.NewBattleMap()
	// 1-ship attacker (1 HP each) vs a heavy cannon next door.
	atk := fleet("atk", tilemap.SideAttacker, tilemap.Position{0, 0}, 1, 1, 1, 1, 2, 1, 5, WeaponBallistic)
	atk.Phalanx.Stacks[0].ShieldHP = 1
	atk.Phalanx.Stacks[0].MaxShield = 1
	atk.Phalanx.Stacks[0].StructHP = 1
	atk.Phalanx.Stacks[0].MaxStruct = 1
	def := building("def", tilemap.Position{2, 2}, 1000, 5000, 5000, 1, 5, 50, WeaponPlanetary)

	_ = world.SetUnitTiles(atk.ID, atk.Footprint)
	_ = world.SetUnitTiles(def.ID, def.Footprint)

	b := NewBattle(world, []*Unit{atk, def}, 1, 30)
	rep := b.Execute()
	if rep.End.Winner != tilemap.SideDefender {
		t.Errorf("expected defender win, got %s", rep.End.Winner)
	}
	if rep.Summary.TotalRounds > 5 {
		t.Errorf("expected attacker to die fast, took %d rounds", rep.Summary.TotalRounds)
	}
}

// TestProp_RoundCapEnforced: configure mutually-non-lethal units (both
// hit but both have effectively infinite HP); battle must reach the
// roundCap and award defender_win.
func TestProp_RoundCapEnforced(t *testing.T) {
	world := tilemap.NewBattleMap()
	atk := fleet("atk", tilemap.SideAttacker, tilemap.Position{0, 0}, 100, 1, 1, 1, 1, 0, 5, WeaponBallistic)
	def := building("def", tilemap.Position{19, 19}, 1000, 1, 1, 1, 1, 5, WeaponPlanetary)
	// Both far apart with range 1 → never come into firing range; round
	// cap forces defender win.
	atk.Phalanx.Stacks[0].MinRange = 100
	atk.Phalanx.Stacks[0].MaxRange = 100
	def.Phalanx.Stacks[0].MinRange = 100
	def.Phalanx.Stacks[0].MaxRange = 100
	atk.Movement = 0

	_ = world.SetUnitTiles(atk.ID, atk.Footprint)
	_ = world.SetUnitTiles(def.ID, def.Footprint)

	b := NewBattle(world, []*Unit{atk, def}, 7, 25)
	rep := b.Execute()
	if rep.End.Reason != "round_cap_reached" {
		t.Errorf("expected round_cap_reached, got %s", rep.End.Reason)
	}
	if rep.End.Winner != tilemap.SideDefender {
		t.Errorf("expected defender win at cap, got %s", rep.End.Winner)
	}
	if rep.Summary.TotalRounds != 25 {
		t.Errorf("expected total_rounds=25, got %d", rep.Summary.TotalRounds)
	}
}

// BenchmarkBattle_SmallScenario measures wall time of one battle in the
// engine. Target: < 50ms per battle (PvP worker SLA).
func BenchmarkBattle_SmallScenario(b *testing.B) {
	for i := 0; i < b.N; i++ {
		world := tilemap.NewBattleMap()
		atk := fleet("atk", tilemap.SideAttacker, tilemap.Position{0, 0}, 200, 100, 200, 1, 5, 3, 10, WeaponBallistic)
		def := building("def", tilemap.Position{10, 10}, 500000, 100, 1000, 1, 8, 5, WeaponPlanetary)
		_ = world.SetUnitTiles(atk.ID, atk.Footprint)
		_ = world.SetUnitTiles(def.ID, def.Footprint)
		batt := NewBattle(world, []*Unit{atk, def}, int64(i+1), 50)
		_ = batt.Execute()
	}
}
