package engine

import (
	"math/rand"

	"github.com/cryptomines-online/backend/internal/combat/replay"
)

// damageOutcome is the result of one stack's volley against another stack.
type damageOutcome struct {
	HitChance         float64
	Hit               bool
	RawDamage         int
	TypeAdvantage     float64
	ArmorEffective    float64
	PositionModifier  float64
	ShieldDamage      int
	StructureDamage   int
	PenetratingDamage int
	Casualties        int
	Crit              bool
}

// resolveAttack runs the GDD §8.6.9 sub-phases for one shooter stack vs.
// one target stack. Returns the outcome and mutates target stack state
// (ShieldHP, StructHP, ShipCount).
//
// The eight sub-phases are conceptually:
//
//	P1 fires — ships * weapon damage range, gated by cooldown (caller-checked)
//	P2 interceptors — TODO (PPC roll); skipped in this baseline
//	P3 hit chance — accuracy vs dodge; clamp 5–95%
//	P4 type advantage — frigate/cruiser/battleship triangle
//	P5 armor effectiveness — damage type vs armor
//	P6 position modifier — phalanx row attack power
//	P7 apply to shields then structure (with optional penetration)
//	P8 convert damage to casualties (ships killed)
//
// Scatter is handled separately by ApplyScatter.
func resolveAttack(rng *rand.Rand, attacker, target *Stack) damageOutcome {
	out := damageOutcome{}
	if attacker == nil || target == nil || attacker.ShipCount <= 0 || target.ShipCount <= 0 {
		return out
	}

	// Phase 3: Hit chance.
	hit := attacker.HitChance
	if hit <= 0 {
		hit = 0.7
	}
	hit += float64(attacker.Accuracy)/1000.0 - float64(target.Dodge)/1000.0
	if hit < 0.05 {
		hit = 0.05
	}
	if hit > 0.95 {
		hit = 0.95
	}
	out.HitChance = hit
	if rng.Float64() > hit {
		return out
	}
	out.Hit = true

	// Phase 1: Raw damage roll × shooting ships.
	span := attacker.BaseDamageMax - attacker.BaseDamageMin
	if span < 0 {
		span = 0
	}
	per := attacker.BaseDamageMin
	if span > 0 {
		per += rng.Intn(span + 1)
	}
	raw := float64(per * attacker.ShipCount)

	// Phase 4: Type advantage.
	out.TypeAdvantage = typeAdvantage(attacker.HullClass, target.HullClass)
	raw *= out.TypeAdvantage

	// Phase 5: Armor effectiveness.
	out.ArmorEffective = armorEffectiveness(attacker.DamageType, target.ArmorType)
	raw *= out.ArmorEffective

	// Phase 6: Position attack modifier.
	out.PositionModifier = positionAttackModifier(attacker.GridRow)
	raw *= out.PositionModifier

	out.RawDamage = int(raw)

	// Phase 7: Shield penetration roll.
	totalDmg := out.RawDamage
	if attacker.ShieldPenChance > 0 && rng.Float64() < attacker.ShieldPenChance {
		out.PenetratingDamage = int(float64(totalDmg) * attacker.ShieldPenPct)
		totalDmg -= out.PenetratingDamage
	}

	// Phase 7: Shields first, then structure.
	if target.ShieldHP > 0 {
		out.ShieldDamage = totalDmg
		if out.ShieldDamage > target.ShieldHP {
			out.ShieldDamage = target.ShieldHP
		}
		target.ShieldHP -= out.ShieldDamage
		totalDmg -= out.ShieldDamage
	}
	leftover := totalDmg + out.PenetratingDamage
	if leftover > 0 && target.StructHP > 0 {
		out.StructureDamage = leftover
		if out.StructureDamage > target.StructHP {
			out.StructureDamage = target.StructHP
		}
		target.StructHP -= out.StructureDamage
	}

	// Phase 8: Casualties from structure damage.
	out.Casualties = casualtiesFrom(out.StructureDamage, target)
	if out.Casualties > target.ShipCount {
		out.Casualties = target.ShipCount
	}
	target.ShipCount -= out.Casualties

	return out
}

// typeAdvantage implements the rock-paper-scissors triangle:
//
//	frigate    >  ship_based — n/a here; we model hull-vs-hull only
//	cruiser    >  frigate      → 1.05
//	battleship >  cruiser      → 1.05
//	frigate    >  battleship   → 1.05
//	(reverse pairings → 0.95)
//	mirror or unknown          → 1.0
func typeAdvantage(attacker, target HullClass) float64 {
	if attacker == HullBuilding || target == HullBuilding {
		return 1.0
	}
	if attacker == target {
		return 1.0
	}
	switch {
	case attacker == HullCruiser && target == HullFrigate:
		return 1.05
	case attacker == HullBattleship && target == HullCruiser:
		return 1.05
	case attacker == HullFrigate && target == HullBattleship:
		return 1.05
	case attacker == HullFrigate && target == HullCruiser:
		return 0.95
	case attacker == HullCruiser && target == HullBattleship:
		return 0.95
	case attacker == HullBattleship && target == HullFrigate:
		return 0.95
	}
	return 1.0
}

// armorEffectiveness mirrors the legacy 4×5 matrix (GDD §8.6.14).
func armorEffectiveness(d DamageType, a ArmorType) float64 {
	switch a {
	case ArmorChrome:
		switch d {
		case DamageExplosive:
			return 1.25
		case DamageKinetic:
			return 0.75
		}
	case ArmorRegen:
		switch d {
		case DamageHeat:
			return 1.25
		case DamageExplosive:
			return 0.75
		}
	case ArmorNano:
		switch d {
		case DamageMagnetic:
			return 1.25
		case DamageHeat:
			return 0.75
		}
	case ArmorNeutralizing:
		switch d {
		case DamageKinetic:
			return 1.25
		case DamageMagnetic:
			return 0.75
		}
	}
	return 1.0
}

// positionAttackModifier from the phalanx row (GDD §8.6.5).
func positionAttackModifier(row int) float64 {
	switch row {
	case 0:
		return 1.00
	case 1:
		return 0.90
	case 2:
		return 0.75
	}
	return 1.00
}

// casualtiesFrom converts structure damage into ships destroyed. If the
// stack has zero structure HP per ship configured (e.g., a building),
// 1 damage = 1 HP off the bag and ShipCount tracks HP directly.
func casualtiesFrom(structureDamage int, target *Stack) int {
	if structureDamage <= 0 || target.ShipCount <= 0 {
		return 0
	}
	hpPerShip := 1
	if target.MaxStruct > 0 && target.ShipCount > 0 {
		hpPerShip = max(1, target.MaxStruct/target.ShipCount)
	}
	if hpPerShip <= 1 {
		return structureDamage
	}
	return structureDamage / hpPerShip
}

// applyScatterToPhalanx distributes scatter damage to stacks adjacent to
// the primary target inside the target's 3x3 phalanx (GDD §8.6.7
// phase 8). Scatter bypasses shields and structure modifiers.
func applyScatterToPhalanx(rec *replay.Recorder, attackerID string, targetUnit *Unit, primaryStack *Stack, scatterDamage int) {
	if scatterDamage <= 0 || primaryStack == nil {
		return
	}
	for _, s := range targetUnit.Phalanx.Stacks {
		if s == nil || s.ShipCount <= 0 || s == primaryStack {
			continue
		}
		dr := s.GridRow - primaryStack.GridRow
		dc := s.GridCol - primaryStack.GridCol
		if abs(dr) > 1 || abs(dc) > 1 {
			continue
		}
		killed := casualtiesFrom(scatterDamage, s)
		if killed > s.ShipCount {
			killed = s.ShipCount
		}
		s.ShipCount -= killed
		if rec != nil {
			rec.LogScatter(attackerID, targetUnit.ID, scatterDamage)
			rec.LogStackUpdate(targetUnit.ID, s.GridRow, s.GridCol, s.ShipCount)
		}
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
