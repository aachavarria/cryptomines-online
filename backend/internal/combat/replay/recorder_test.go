package replay

import (
	"encoding/json"
	"testing"

	"github.com/cryptomines-online/backend/internal/combat/tilemap"
)

func TestRecorder_RoundTripsToJSON(t *testing.T) {
	r := NewRecorder()
	r.SetMap(20, 20, []TerrainView{
		{Kind: "meteor_star", Position: tilemap.Position{5, 5}, HP: 100, MaxHP: 100},
	})
	r.AddInitialUnit(UnitState{
		UnitID: "atk_0", Side: tilemap.SideAttacker, Kind: UnitFleet,
		FleetID: "f1", Tile: tilemap.Position{0, 0},
	})
	r.BeginRound(1)
	r.LogMove("atk_0", []tilemap.Position{{0, 0}, {1, 1}, {2, 2}})
	r.LogFire(Action{
		Attacker: "atk_0", Target: "def_b1", Weapon: "missile",
		Damage: 1500, Casualties: 3, Crit: false,
	})
	r.LogStackUpdate("def_b1", 0, 1, 0)
	r.LogDestroy("def_b1")
	r.SetEnd(tilemap.SideAttacker, EndAllDefendersDestroy)
	r.SetSummary(Summary{TotalRounds: 1, AttackerCasualties: 0, DefenderCasualties: 3, BuildingsDestroyed: 1})

	rep := r.Build()
	if rep.SchemaVersion != CurrentSchemaVersion {
		t.Errorf("schema_version = %d, want %d", rep.SchemaVersion, CurrentSchemaVersion)
	}
	if len(rep.Rounds) != 1 || len(rep.Rounds[0].Actions) != 4 {
		t.Errorf("unexpected actions: %+v", rep.Rounds)
	}

	// Round-trip through JSON to confirm tags are correct.
	raw, err := json.Marshal(rep)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got Report
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.End.Winner != tilemap.SideAttacker || got.End.Reason != EndAllDefendersDestroy {
		t.Errorf("end state lost in roundtrip: %+v", got.End)
	}
	if got.Rounds[0].Actions[0].Type != ActionMove {
		t.Errorf("first action type = %s, want move", got.Rounds[0].Actions[0].Type)
	}
}

func TestRecorder_NoMovementZeroLengthPathSkipped(t *testing.T) {
	r := NewRecorder()
	r.BeginRound(1)
	r.LogMove("u1", []tilemap.Position{{0, 0}}) // single tile = no movement
	rep := r.Build()
	if len(rep.Rounds[0].Actions) != 0 {
		t.Errorf("zero-length move should be skipped, got %+v", rep.Rounds[0].Actions)
	}
}
