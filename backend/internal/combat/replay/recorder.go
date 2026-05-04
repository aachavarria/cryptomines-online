package replay

import "github.com/cryptomines-online/backend/internal/combat/tilemap"

// Recorder accumulates Actions during combat and finalizes them into a
// Report. The engine calls one method per observable event; the recorder
// does not interpret semantics. Cheap to instantiate per battle.
type Recorder struct {
	report      Report
	currentRound *Round
}

// NewRecorder starts a fresh recording for a single battle.
func NewRecorder() *Recorder {
	return &Recorder{
		report: Report{
			SchemaVersion: CurrentSchemaVersion,
			Rounds:        []Round{},
		},
	}
}

// SetMap captures the static battlefield layout.
func (r *Recorder) SetMap(width, height int, terrain []TerrainView) {
	r.report.BattleMap = BattleMap{Width: width, Height: height, Terrain: terrain}
}

// AddInitialUnit appends one unit to the spawn snapshot.
func (r *Recorder) AddInitialUnit(u UnitState) {
	r.report.InitialUnits = append(r.report.InitialUnits, u)
}

// BeginRound starts a new round bucket. Actions logged after this go
// into the new round.
func (r *Recorder) BeginRound(n int) {
	r.report.Rounds = append(r.report.Rounds, Round{Round: n, Actions: []Action{}})
	r.currentRound = &r.report.Rounds[len(r.report.Rounds)-1]
}

// LogMove records a fleet's tile-by-tile path for this round.
func (r *Recorder) LogMove(unitID string, path []tilemap.Position) {
	if len(path) < 2 {
		return // no movement
	}
	from := path[0]
	to := path[len(path)-1]
	r.append(Action{
		Type:   ActionMove,
		UnitID: unitID,
		From:   &from,
		To:     &to,
		Path:   path,
	})
}

// LogFire records a single attacker→target damage event.
func (r *Recorder) LogFire(a Action) {
	a.Type = ActionFire
	r.append(a)
}

// LogAOE records that a weapon hit a primary target plus surrounding tiles.
func (r *Recorder) LogAOE(source string, center tilemap.Position, tiles []tilemap.Position, damagePct float64) {
	r.append(Action{
		Type:      ActionAOE,
		Source:    source,
		Center:    &center,
		Tiles:     tiles,
		DamagePct: damagePct,
	})
}

// LogScatter records phalanx-internal splash (GDD §8.6.7 phase 8).
func (r *Recorder) LogScatter(source, target string, damage int) {
	r.append(Action{
		Type:   ActionScatter,
		Source: source,
		Target: target,
		Damage: damage,
	})
}

// LogStackUpdate records a per-phalanx-cell ship count change.
func (r *Recorder) LogStackUpdate(unitID string, gridRow, gridCol, newShipCount int) {
	r.append(Action{
		Type:         ActionStackUpdate,
		UnitID:       unitID,
		StackGridRow: gridRow,
		StackGridCol: gridCol,
		NewShipCount: newShipCount,
	})
}

// LogTerrainHit records damage to a terrain tile (e.g., shooting a
// Meteor Star to clear a path).
func (r *Recorder) LogTerrainHit(source string, p tilemap.Position, damage int) {
	r.append(Action{
		Type:       ActionTerrainHit,
		Source:     source,
		TerrainPos: &p,
		Damage:     damage,
	})
}

// LogDestroy records that a unit was eliminated this round.
func (r *Recorder) LogDestroy(unitID string) {
	r.append(Action{Type: ActionDestroy, UnitID: unitID})
}

// SetEnd records the final outcome.
func (r *Recorder) SetEnd(winner tilemap.Side, reason EndReason) {
	r.report.End = EndState{Winner: winner, Reason: reason}
}

// SetSummary records the headline post-battle stats.
func (r *Recorder) SetSummary(s Summary) {
	r.report.Summary = s
}

// Build returns the finalized Report by value (caller may mutate freely).
func (r *Recorder) Build() Report {
	return r.report
}

func (r *Recorder) append(a Action) {
	if r.currentRound == nil {
		// Defensive: append a synthetic round 0 for early actions.
		r.report.Rounds = append(r.report.Rounds, Round{Round: 0, Actions: []Action{}})
		r.currentRound = &r.report.Rounds[len(r.report.Rounds)-1]
	}
	r.currentRound.Actions = append(r.currentRound.Actions, a)
}
