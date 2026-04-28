package combat

import (
	"testing"

	"github.com/cryptomines-online/backend/internal/services"
)

// makeTargetStack returns a minimally-populated stack for selectTarget tests.
func makeTargetStack(id string, attack, durability, row int) *FleetStack {
	return &FleetStack{
		ID:               id,
		ShipCount:        10,
		CurrentShips:     10,
		EffectiveAttack:  attack,
		CurrentShield:    durability / 2,
		CurrentStructure: durability / 2,
		GridRow:          row,
		GridCol:          0,
	}
}

// runSelect runs selectTarget for a given attacker targeting strategy.
// The attacker stack is at row 0; defender stacks are passed in.
func runSelect(targeting string, attacker *FleetStack, defenders []*FleetStack) *FleetStack {
	ce := NewCombatEngine(1)
	state := &CombatState{
		Attacker: &Fleet{Stacks: []*FleetStack{attacker}, Targeting: targeting, Side: "attacker"},
		Defender: &Fleet{Stacks: defenders, Targeting: "max_attack", Side: "defender"},
	}
	return ce.selectTarget(state, attacker)
}

func TestSelectTarget_MaxAttack(t *testing.T) {
	a := makeTargetStack("a", 50, 100, 0)
	weak := makeTargetStack("d_weak", 20, 200, 0)
	strong := makeTargetStack("d_strong", 80, 100, 1)
	got := runSelect("max_attack", a, []*FleetStack{weak, strong})
	if got != strong {
		t.Fatalf("max_attack should pick highest EffectiveAttack, got %v", got.ID)
	}
}

func TestSelectTarget_MinAttack(t *testing.T) {
	a := makeTargetStack("a", 50, 100, 0)
	weak := makeTargetStack("d_weak", 20, 200, 0)
	strong := makeTargetStack("d_strong", 80, 100, 1)
	got := runSelect("min_attack", a, []*FleetStack{weak, strong})
	if got != weak {
		t.Fatalf("min_attack should pick lowest EffectiveAttack, got %v", got.ID)
	}
}

func TestSelectTarget_MaxDurability(t *testing.T) {
	a := makeTargetStack("a", 50, 100, 0)
	soft := makeTargetStack("d_soft", 50, 50, 0)
	tough := makeTargetStack("d_tough", 50, 800, 1)
	got := runSelect("max_durability", a, []*FleetStack{soft, tough})
	if got != tough {
		t.Fatalf("max_durability should pick most HP, got %v", got.ID)
	}
}

func TestSelectTarget_MinDurability(t *testing.T) {
	a := makeTargetStack("a", 50, 100, 0)
	soft := makeTargetStack("d_soft", 50, 50, 0)
	tough := makeTargetStack("d_tough", 50, 800, 1)
	got := runSelect("min_durability", a, []*FleetStack{soft, tough})
	if got != soft {
		t.Fatalf("min_durability should pick least HP, got %v", got.ID)
	}
}

func TestSelectTarget_Closest(t *testing.T) {
	a := makeTargetStack("a", 50, 100, 0)
	near := makeTargetStack("d_near", 50, 100, 1)
	far := makeTargetStack("d_far", 50, 100, 4)
	got := runSelect("closest", a, []*FleetStack{far, near})
	if got != near {
		t.Fatalf("closest should pick nearest GridRow, got %v", got.ID)
	}
}

func TestSelectTarget_DeadStacksSkipped(t *testing.T) {
	a := makeTargetStack("a", 50, 100, 0)
	dead := makeTargetStack("d_dead", 999, 999, 0)
	dead.CurrentShips = 0
	alive := makeTargetStack("d_alive", 10, 10, 1)
	got := runSelect("max_attack", a, []*FleetStack{dead, alive})
	if got != alive {
		t.Fatalf("dead targets must be filtered, got %v", got.ID)
	}
}

func TestSelectTarget_NoAliveTargetsReturnsNil(t *testing.T) {
	a := makeTargetStack("a", 50, 100, 0)
	d := makeTargetStack("d", 50, 100, 0)
	d.CurrentShips = 0
	got := runSelect("max_attack", a, []*FleetStack{d})
	if got != nil {
		t.Fatalf("expected nil when no alive targets, got %v", got.ID)
	}
}

// Sanity: with the unused services.TechBonuses import wired in, ensure the
// fleet shape used by the production handler still compiles in this package.
func TestSelectTarget_FleetShapeCompat(t *testing.T) {
	_ = &Fleet{TechBonuses: &services.TechBonuses{}}
}
