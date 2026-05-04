// Package replay defines the on-the-wire schema persisted in
// combat_reports.rounds_json (GDD §3.12). The frontend's BattlePlayback
// component consumes this JSON to render an animated tile-based replay.
//
// schema_version is bumped whenever the action vocabulary changes. Old
// reports remain readable forever; the frontend dispatches on the version.
package replay

import "github.com/cryptomines-online/backend/internal/combat/tilemap"

// CurrentSchemaVersion identifies the tile-based replay schema introduced
// alongside the §8.6 rewrite. Legacy abstract-engine reports use
// version 1; tile-based reports use version 2.
const CurrentSchemaVersion = 2

// EndReason explains why the battle ended.
type EndReason string

const (
	EndAttackerWiped       EndReason = "attacker_wiped"
	EndDefenderWiped       EndReason = "defender_wiped"
	EndRoundCapReached     EndReason = "round_cap_reached"
	EndAllDefendersDestroy EndReason = "all_defenders_destroyed"
	EndAllAttackersDestroy EndReason = "all_attackers_destroyed"
)

// UnitKind distinguishes mobile fleets from stationary defense buildings.
type UnitKind string

const (
	UnitFleet    UnitKind = "fleet"
	UnitBuilding UnitKind = "building"
)

// Report is the full payload stored in combat_reports.rounds_json.
type Report struct {
	SchemaVersion int          `json:"schema_version"`
	BattleMap     BattleMap    `json:"battle_map"`
	InitialUnits  []UnitState  `json:"initial_units"`
	Rounds        []Round      `json:"rounds"`
	End           EndState     `json:"end"`
	Summary       Summary      `json:"summary"`
}

// BattleMap captures the static layout: terrain and dimensions.
type BattleMap struct {
	Width   int           `json:"width"`
	Height  int           `json:"height"`
	Terrain []TerrainView `json:"terrain"`
}

// TerrainView is the persisted shape of a single terrain entity.
type TerrainView struct {
	Kind     string           `json:"kind"`
	Position tilemap.Position `json:"position"`
	HP       int              `json:"hp"`
	MaxHP    int              `json:"max_hp"`
}

// UnitState snapshots a unit's starting state for the replay viewer to
// hydrate its scene.
type UnitState struct {
	UnitID     string           `json:"unit_id"`
	Side       tilemap.Side     `json:"side"`
	Kind       UnitKind         `json:"kind"`
	FleetID    string           `json:"fleet_id,omitempty"`
	BuildingID string           `json:"building_id,omitempty"`
	TypeName   string           `json:"type_name,omitempty"`
	Tile       tilemap.Position `json:"tile"`
	Footprint  []tilemap.Position `json:"footprint,omitempty"`
	HP         int              `json:"hp,omitempty"`
	Phalanx    []PhalanxStack   `json:"phalanx,omitempty"`
}

// PhalanxStack is a per-fleet exposure-grid entry (GDD §8.6.5). Stored on
// the InitialUnits entries so the replay can render each ship class in its
// position. Damage events that change ship counts emit per-stack updates
// inside Round.Actions.
type PhalanxStack struct {
	GridRow      int    `json:"grid_row"`
	GridCol      int    `json:"grid_col"`
	ShipDesignID string `json:"ship_design_id"`
	HullClass    string `json:"hull_class"`
	ShipCount    int    `json:"ship_count"`
}

// Round is one complete cycle of every unit acting once.
type Round struct {
	Round   int      `json:"round"`
	Actions []Action `json:"actions"`
}

// ActionType is the discriminator for the union-shape Action.
type ActionType string

const (
	ActionMove        ActionType = "move"
	ActionFire        ActionType = "fire"
	ActionAOE         ActionType = "aoe"
	ActionScatter     ActionType = "scatter"
	ActionStackUpdate ActionType = "stack_update"
	ActionTerrainHit  ActionType = "terrain_hit"
	ActionDestroy     ActionType = "destroy"
)

// Action is a discriminated record describing one observable event in a
// round. Fields are sparsely populated based on Type — JSON consumers
// should branch on Type before reading optional fields.
type Action struct {
	Type ActionType `json:"type"`

	// Movement.
	UnitID string             `json:"unit_id,omitempty"`
	From   *tilemap.Position  `json:"from,omitempty"`
	To     *tilemap.Position  `json:"to,omitempty"`
	Path   []tilemap.Position `json:"path,omitempty"`

	// Fire.
	Attacker     string `json:"attacker,omitempty"`
	Target       string `json:"target,omitempty"`
	Weapon       string `json:"weapon,omitempty"`
	Damage       int    `json:"damage,omitempty"`
	Casualties   int    `json:"casualties,omitempty"`
	Crit         bool   `json:"crit,omitempty"`
	StackGridRow int    `json:"stack_grid_row,omitempty"`
	StackGridCol int    `json:"stack_grid_col,omitempty"`

	// AoE / Scatter.
	Source     string             `json:"source,omitempty"`
	Center     *tilemap.Position  `json:"center,omitempty"`
	Tiles      []tilemap.Position `json:"tiles,omitempty"`
	DamagePct  float64            `json:"damage_pct,omitempty"`

	// Stack update (per-phalanx-cell ship count change after damage).
	NewShipCount int `json:"new_ship_count,omitempty"`

	// Terrain damage.
	TerrainPos *tilemap.Position `json:"terrain_pos,omitempty"`
}

// EndState records the final outcome.
type EndState struct {
	Winner tilemap.Side `json:"winner"`
	Reason EndReason    `json:"reason"`
}

// Summary keeps the headline stats redundantly so a UI can render a
// post-battle card without scanning every action.
type Summary struct {
	TotalRounds         int   `json:"total_rounds"`
	AttackerCasualties  int   `json:"attacker_casualties"`
	DefenderCasualties  int   `json:"defender_casualties"`
	BuildingsDestroyed  int   `json:"buildings_destroyed"`
	He3Consumed         int64 `json:"he3_consumed"`
}
