package engine

import (
	"math/rand"
	"testing"
)

func mkStack(ships, perShip, shield, structHP int, dmg DamageType, hull HullClass, armor ArmorType) *Stack {
	return &Stack{
		ID: "s", HullClass: hull, ShipCount: ships,
		DamageType: dmg, ArmorType: armor,
		WeaponClass: WeaponBallistic, WeaponName: "ballistic",
		BaseDamageMin: perShip, BaseDamageMax: perShip,
		HitChance: 1.0,
		ShieldHP:  shield, MaxShield: shield,
		StructHP:  structHP, MaxStruct: structHP,
	}
}

func TestResolveAttack_KineticVsChrome_ResistsDamage(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	atk := mkStack(10, 100, 0, 0, DamageKinetic, HullFrigate, ArmorNone)
	def := mkStack(10, 0, 1000, 1000, DamageKinetic, HullFrigate, ArmorChrome)

	out := resolveAttack(rng, atk, def)
	if !out.Hit {
		t.Fatal("expected guaranteed hit")
	}
	if out.ArmorEffective != 0.75 {
		t.Errorf("Chrome should resist Kinetic at 0.75x, got %v", out.ArmorEffective)
	}
}

func TestResolveAttack_ExplosiveVsChrome_AmplifiesDamage(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	atk := mkStack(10, 100, 0, 0, DamageExplosive, HullFrigate, ArmorNone)
	def := mkStack(10, 0, 1000, 1000, DamageExplosive, HullFrigate, ArmorChrome)

	out := resolveAttack(rng, atk, def)
	if out.ArmorEffective != 1.25 {
		t.Errorf("Chrome should amplify Explosive at 1.25x, got %v", out.ArmorEffective)
	}
}

func TestResolveAttack_TypeAdvantageTriangle(t *testing.T) {
	cases := []struct {
		atk, def HullClass
		want     float64
	}{
		{HullCruiser, HullFrigate, 1.05},
		{HullBattleship, HullCruiser, 1.05},
		{HullFrigate, HullBattleship, 1.05},
		{HullFrigate, HullCruiser, 0.95},
		{HullCruiser, HullBattleship, 0.95},
		{HullBattleship, HullFrigate, 0.95},
		{HullFrigate, HullFrigate, 1.0},
	}
	for _, c := range cases {
		got := typeAdvantage(c.atk, c.def)
		if got != c.want {
			t.Errorf("typeAdvantage(%s,%s) = %v, want %v", c.atk, c.def, got, c.want)
		}
	}
}

func TestResolveAttack_PositionAttackModifier(t *testing.T) {
	if positionAttackModifier(0) != 1.00 {
		t.Errorf("front rank should be 1.00")
	}
	if positionAttackModifier(1) != 0.90 {
		t.Errorf("middle rank should be 0.90")
	}
	if positionAttackModifier(2) != 0.75 {
		t.Errorf("back rank should be 0.75")
	}
}

func TestResolveAttack_ShieldsAbsorbBeforeStructure(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	atk := mkStack(1, 100, 0, 0, DamageKinetic, HullFrigate, ArmorNone)
	def := mkStack(10, 0, 60, 1000, DamageKinetic, HullFrigate, ArmorNone)

	out := resolveAttack(rng, atk, def)
	if !out.Hit {
		t.Fatal("expected hit")
	}
	if out.ShieldDamage == 0 {
		t.Error("expected shield damage to be applied first")
	}
	if def.ShieldHP > 0 && out.StructureDamage > 0 {
		t.Error("structure damage should only apply when shields are depleted")
	}
}

func TestResolveAttack_ShieldPenetrationBypassesShields(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	atk := mkStack(1, 1000, 0, 0, DamageKinetic, HullFrigate, ArmorNone)
	atk.ShieldPenChance = 1.0
	atk.ShieldPenPct = 0.5
	def := mkStack(10, 0, 99999, 99999, DamageKinetic, HullFrigate, ArmorNone)

	out := resolveAttack(rng, atk, def)
	if out.PenetratingDamage == 0 {
		t.Error("expected penetration damage > 0")
	}
}

func TestResolveAttack_LowHitChanceMisses(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	atk := mkStack(1, 100, 0, 0, DamageKinetic, HullFrigate, ArmorNone)
	atk.HitChance = 0.001 // clamped up to 5% by the engine
	def := mkStack(10, 0, 1000, 1000, DamageKinetic, HullFrigate, ArmorNone)
	// Boost target dodge so the effective hit rate stays at the 5% floor.
	def.Dodge = 1000

	misses := 0
	for i := 0; i < 200; i++ {
		def.ShieldHP = 1000
		def.StructHP = 1000
		out := resolveAttack(rng, atk, def)
		if !out.Hit {
			misses++
		}
	}
	// 5% expected hit rate → ~95% misses with sampling noise.
	if misses < 170 {
		t.Errorf("expected mostly misses at 5%% hit chance, got %d/200 misses", misses)
	}
}
