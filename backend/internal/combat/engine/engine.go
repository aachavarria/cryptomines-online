package engine

import (
	"math/rand"
	"sort"

	"github.com/cryptomines-online/backend/internal/combat/pathing"
	"github.com/cryptomines-online/backend/internal/combat/replay"
	"github.com/cryptomines-online/backend/internal/combat/tilemap"
)

// Battle is one fully-prepared combat instance ready to run. Construct it
// with NewBattle, then call Execute to resolve.
type Battle struct {
	Map      *tilemap.BattleMap
	Units    []*Unit
	rng      *rand.Rand
	rec      *replay.Recorder
	roundCap int
}

// NewBattle wires together a battle from a fully-spawned set of units and
// a populated battle map. Caller is responsible for placing units on the
// map (via tilemap.SetUnitTiles) and registering terrain before calling.
//
// roundCap follows GDD §8.6.1:
//
//	min(99, 20 + #defenderFleets + #attackerFleets + #defenseStructures)
//
// Pass 0 to compute it automatically from the unit list.
func NewBattle(world *tilemap.BattleMap, units []*Unit, seed int64, roundCap int) *Battle {
	if roundCap <= 0 {
		roundCap = computeRoundCap(units)
	}
	return &Battle{
		Map:      world,
		Units:    units,
		rng:      rand.New(rand.NewSource(seed)),
		rec:      replay.NewRecorder(),
		roundCap: roundCap,
	}
}

func computeRoundCap(units []*Unit) int {
	cap := 20
	for _, u := range units {
		cap++
		_ = u
	}
	if cap > 99 {
		cap = 99
	}
	return cap
}

// Execute runs the round loop until one side wins or the round cap is
// reached. Returns the recorded report (always non-empty) and the
// outcome.
func (b *Battle) Execute() replay.Report {
	b.captureInitialState()

	endReason := replay.EndRoundCapReached
	winner := tilemap.SideAttacker

	for round := 1; round <= b.roundCap; round++ {
		b.rec.BeginRound(round)
		order := b.initiativeOrder()

		for _, u := range order {
			if !u.IsAlive() {
				continue
			}
			b.takeTurn(u)
			if w, done := b.checkVictory(); done {
				winner = w
				endReason = victoryReason(w)
				b.finalize(winner, endReason, round)
				return b.rec.Build()
			}
		}
		b.tickCooldowns()
	}

	// Round cap reached: defender wins by default (attacker failed to clear).
	winner = tilemap.SideDefender
	endReason = replay.EndRoundCapReached
	b.finalize(winner, endReason, b.roundCap)
	return b.rec.Build()
}

func victoryReason(w tilemap.Side) replay.EndReason {
	if w == tilemap.SideAttacker {
		return replay.EndAllDefendersDestroy
	}
	return replay.EndAllAttackersDestroy
}

// captureInitialState writes the spawn snapshot and terrain to the
// replay recorder before the first round.
func (b *Battle) captureInitialState() {
	terrain := []replay.TerrainView{}
	for _, t := range b.Map.AllTerrain() {
		terrain = append(terrain, replay.TerrainView{
			Kind: string(t.Kind), Position: t.Position, HP: t.HP, MaxHP: t.MaxHP,
		})
	}
	b.rec.SetMap(tilemap.Width, tilemap.Height, terrain)

	for _, u := range b.Units {
		var phalanx []replay.PhalanxStack
		for _, s := range u.Phalanx.Stacks {
			if s == nil {
				continue
			}
			phalanx = append(phalanx, replay.PhalanxStack{
				GridRow: s.GridRow, GridCol: s.GridCol,
				ShipDesignID: s.ShipDesignID, HullClass: string(s.HullClass),
				ShipCount: s.ShipCount,
			})
		}
		kind := replay.UnitFleet
		if u.Kind == UnitBuilding {
			kind = replay.UnitBuilding
		}
		b.rec.AddInitialUnit(replay.UnitState{
			UnitID: u.ID, Side: u.Side, Kind: kind,
			FleetID: u.FleetID, BuildingID: u.BuildingID, TypeName: u.TypeName,
			Tile: u.Tile, Footprint: u.Footprint,
			HP: u.TotalShips(), Phalanx: phalanx,
		})
	}
}

// initiativeOrder sorts living units by Speed desc; defender wins ties
// (GDD §8.6.6). Buildings have a fixed Speed value provided by the caller.
func (b *Battle) initiativeOrder() []*Unit {
	alive := make([]*Unit, 0, len(b.Units))
	for _, u := range b.Units {
		if u.IsAlive() {
			alive = append(alive, u)
		}
	}
	sort.SliceStable(alive, func(i, j int) bool {
		if alive[i].Speed != alive[j].Speed {
			return alive[i].Speed > alive[j].Speed
		}
		// Defender first on ties.
		if alive[i].Side != alive[j].Side {
			return alive[i].Side == tilemap.SideDefender
		}
		return alive[i].ID < alive[j].ID
	})
	return alive
}

// takeTurn runs one unit's full turn: move (fleets only), acquire target,
// fire, AoE/scatter, then tick.
func (b *Battle) takeTurn(u *Unit) {
	enemies := b.enemiesOf(u.Side)

	// 1. MOVE — fleets only.
	if u.Kind == UnitFleet && u.Movement > 0 {
		dest := pickDestination(u, enemies, b.Map)
		if !dest.Equals(u.Tile) {
			path := pathing.FindPath(b.Map, u.ID, u.Tile, dest, u.Movement)
			if len(path) > 1 {
				newTile := path[len(path)-1]
				if err := b.Map.MoveUnit(u.ID, u.Tile, newTile); err == nil {
					u.Tile = newTile
					u.Footprint = []tilemap.Position{newTile}
					b.rec.LogMove(u.ID, path)
				}
			}
		}
	}

	// 2. ACQUIRE TARGET.
	target := chooseTarget(u, enemies, b.Map)
	if target == nil {
		return
	}

	// 3. FIRE — every active stack on this unit fires once at the chosen
	//    target. The 8 sub-phases run inside resolveAttack (damage.go).
	totalCasualties := 0
	totalDamage := 0
	for _, attacker := range u.Phalanx.Stacks {
		if attacker == nil || attacker.ShipCount <= 0 || attacker.CooldownLeft > 0 {
			continue
		}
		victim := primaryStack(target)
		if victim == nil {
			break
		}
		out := resolveAttack(b.rng, attacker, victim)
		attacker.CooldownLeft = attacker.Cooldown
		if !out.Hit {
			continue
		}
		totalDamage += out.RawDamage
		totalCasualties += out.Casualties

		b.rec.LogFire(replay.Action{
			Attacker: u.ID, Target: target.ID,
			Weapon: attacker.WeaponName,
			Damage: out.RawDamage, Casualties: out.Casualties, Crit: out.Crit,
			StackGridRow: attacker.GridRow, StackGridCol: attacker.GridCol,
		})
		b.rec.LogStackUpdate(target.ID, victim.GridRow, victim.GridCol, victim.ShipCount)

		// Scatter inside the target's phalanx (GDD §8.6.7 phase 8).
		if attacker.ScatterPct > 0 {
			scatter := int(float64(out.RawDamage) * attacker.ScatterPct)
			applyScatterToPhalanx(b.rec, u.ID, target, victim, scatter)
		}

		// Battle-map AoE — Anti-Aircraft Gun: 50% damage to 8 surrounding
		// tiles (GDD §8.6.4). Applies once per attacker volley, against
		// the original target's tile.
		if attacker.WeaponClass == WeaponPlanetaryAOE {
			neighbors := tilemap.Neighbors8(target.Tile)
			b.rec.LogAOE(u.ID, target.Tile, neighbors, 0.5)
			for _, n := range neighbors {
				id := b.Map.UnitAt(n)
				if id == "" || id == target.ID {
					continue
				}
				neighbor := b.findUnit(id)
				if neighbor == nil || neighbor.Side == u.Side || !neighbor.IsAlive() {
					continue
				}
				splashTarget := primaryStack(neighbor)
				if splashTarget == nil {
					continue
				}
				splashOut := resolveAttack(b.rng, halfDamageProxy(attacker), splashTarget)
				_ = splashOut
				if !neighbor.IsAlive() {
					neighbor.MarkDead()
					b.Map.ClearUnit(neighbor.ID)
					b.rec.LogDestroy(neighbor.ID)
				}
			}
		}
	}
	_ = totalCasualties
	_ = totalDamage

	// 5. Mark destroyed targets.
	if !target.IsAlive() {
		target.MarkDead()
		b.Map.ClearUnit(target.ID)
		b.rec.LogDestroy(target.ID)
	}
}

// halfDamageProxy returns a copy of stack `s` with damage halved, used
// to apply AoE splash at 50% (Anti-Aircraft Gun). The copy is
// cooldown-cleared so resolveAttack does not gate it.
func halfDamageProxy(s *Stack) *Stack {
	c := *s
	c.BaseDamageMin /= 2
	c.BaseDamageMax /= 2
	c.CooldownLeft = 0
	return &c
}

func (b *Battle) enemiesOf(side tilemap.Side) []*Unit {
	out := make([]*Unit, 0, len(b.Units))
	for _, u := range b.Units {
		if u.Side != side && u.IsAlive() {
			out = append(out, u)
		}
	}
	return out
}

func (b *Battle) findUnit(id string) *Unit {
	for _, u := range b.Units {
		if u.ID == id {
			return u
		}
	}
	return nil
}

// checkVictory inspects living units. Returns (winningSide, true) when
// one side has no living units left. Returns (_, false) otherwise.
func (b *Battle) checkVictory() (tilemap.Side, bool) {
	atkAlive, defAlive := 0, 0
	for _, u := range b.Units {
		if !u.IsAlive() {
			continue
		}
		if u.Side == tilemap.SideAttacker {
			atkAlive++
		} else {
			defAlive++
		}
	}
	switch {
	case atkAlive == 0 && defAlive == 0:
		return tilemap.SideDefender, true
	case atkAlive == 0:
		return tilemap.SideDefender, true
	case defAlive == 0:
		return tilemap.SideAttacker, true
	default:
		return "", false
	}
}

// tickCooldowns decrements each stack's CooldownLeft at the end of the
// round.
func (b *Battle) tickCooldowns() {
	for _, u := range b.Units {
		for _, s := range u.Phalanx.Stacks {
			if s != nil && s.CooldownLeft > 0 {
				s.CooldownLeft--
			}
		}
	}
}

func (b *Battle) finalize(winner tilemap.Side, reason replay.EndReason, totalRounds int) {
	b.rec.SetEnd(winner, reason)
	atkCas, defCas, blds := 0, 0, 0
	for _, u := range b.Units {
		// Casualties: original ships - current. We don't track originals
		// here, so this summary is a conservative placeholder; integration
		// layer can override Summary post-hoc with exact numbers.
		_ = u
	}
	for _, u := range b.Units {
		if !u.IsAlive() && u.Kind == UnitBuilding && u.Side == tilemap.SideDefender {
			blds++
		}
	}
	b.rec.SetSummary(replay.Summary{
		TotalRounds:        totalRounds,
		AttackerCasualties: atkCas,
		DefenderCasualties: defCas,
		BuildingsDestroyed: blds,
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
