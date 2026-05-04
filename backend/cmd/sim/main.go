// sim is a headless battle simulator for the tile-based combat engine.
// It does not touch the database; it constructs synthetic units from a
// YAML/JSON scenario, runs Battle.Execute, and prints the result.
//
// Usage:
//
//	go run ./cmd/sim -scenario testdata/missile_vs_thors.json -seed 42 -verbose
//	go run ./cmd/sim -bench 1000 -scenario testdata/missile_vs_thors.json
//
// Used by Stream E for balancing (range/HP tuning) and benchmarking.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cryptomines-online/backend/internal/combat/engine"
	"github.com/cryptomines-online/backend/internal/combat/replay"
	"github.com/cryptomines-online/backend/internal/combat/tilemap"
)

// scenarioFile mirrors the on-disk shape used by the CLI. Only fields
// the CLI needs are documented here; engine types are reused for the
// rest.
type scenarioFile struct {
	Name     string         `json:"name"`
	Terrain  []terrainEntry `json:"terrain"`
	Units    []unitEntry    `json:"units"`
	RoundCap int            `json:"round_cap"`
}

type terrainEntry struct {
	Kind string `json:"kind"`
	Col  int    `json:"col"`
	Row  int    `json:"row"`
	HP   int    `json:"hp"`
}

type unitEntry struct {
	ID         string       `json:"id"`
	Side       string       `json:"side"` // "attacker" | "defender"
	Kind       string       `json:"kind"` // "fleet" | "building"
	Tile       [2]int       `json:"tile"`
	Speed      int          `json:"speed"`
	Movement   int          `json:"movement"`
	TargetCmd  string       `json:"target_cmd"`
	Stacks     []stackEntry `json:"stacks"`
	Footprint  [][2]int     `json:"footprint,omitempty"`
}

type stackEntry struct {
	GridRow       int     `json:"grid_row"`
	GridCol       int     `json:"grid_col"`
	HullClass     string  `json:"hull_class"`
	ShipCount     int     `json:"ship_count"`
	WeaponClass   string  `json:"weapon_class"`
	WeaponName    string  `json:"weapon_name"`
	DamageType    string  `json:"damage_type"`
	ArmorType     string  `json:"armor_type"`
	MinRange      int     `json:"min_range"`
	MaxRange      int     `json:"max_range"`
	Cooldown      int     `json:"cooldown"`
	BaseDamageMin int     `json:"base_damage_min"`
	BaseDamageMax int     `json:"base_damage_max"`
	HitChance     float64 `json:"hit_chance"`
	ShieldHP      int     `json:"shield_hp"`
	StructHP      int     `json:"struct_hp"`
	Accuracy      int     `json:"accuracy"`
	Dodge         int     `json:"dodge"`
}

func main() {
	scenarioPath := flag.String("scenario", "", "path to JSON scenario file")
	seed := flag.Int64("seed", 0, "RNG seed (0 = use time.Now)")
	bench := flag.Int("bench", 0, "run N battles and print timing stats")
	verbose := flag.Bool("verbose", false, "print full Report JSON")
	flag.Parse()

	if *scenarioPath == "" {
		log.Fatal("missing -scenario")
	}
	raw, err := os.ReadFile(*scenarioPath)
	if err != nil {
		log.Fatalf("read scenario: %v", err)
	}
	var sc scenarioFile
	if err := json.Unmarshal(raw, &sc); err != nil {
		log.Fatalf("parse scenario: %v", err)
	}

	if *bench > 0 {
		runBench(&sc, *bench, *seed)
		return
	}

	useSeed := *seed
	if useSeed == 0 {
		useSeed = time.Now().UnixNano()
	}
	rep := runOne(&sc, useSeed)
	if *verbose {
		out, _ := json.MarshalIndent(rep, "", "  ")
		fmt.Println(string(out))
		return
	}
	fmt.Printf("scenario : %s\n", sc.Name)
	fmt.Printf("winner   : %s\n", rep.End.Winner)
	fmt.Printf("reason   : %s\n", rep.End.Reason)
	fmt.Printf("rounds   : %d\n", rep.Summary.TotalRounds)
	fmt.Printf("buildings_destroyed: %d\n", rep.Summary.BuildingsDestroyed)
}

func runOne(sc *scenarioFile, seed int64) replay.Report {
	world := tilemap.NewBattleMap()
	for _, t := range sc.Terrain {
		_ = world.AddTerrain(&tilemap.Terrain{
			Kind:     tilemap.TerrainKind(t.Kind),
			Position: tilemap.Position{Col: t.Col, Row: t.Row},
			HP:       t.HP, MaxHP: t.HP,
		})
	}
	units := []*engine.Unit{}
	for _, ue := range sc.Units {
		u := unitFromEntry(ue)
		if err := world.SetUnitTiles(u.ID, u.Footprint); err != nil {
			log.Printf("place %s: %v", u.ID, err)
			continue
		}
		units = append(units, u)
	}
	b := engine.NewBattle(world, units, seed, sc.RoundCap)
	return b.Execute()
}

func runBench(sc *scenarioFile, n int, seed int64) {
	if seed == 0 {
		seed = 1
	}
	atkWins, defWins := 0, 0
	totalRounds := 0
	start := time.Now()
	for i := 0; i < n; i++ {
		rep := runOne(sc, seed+int64(i))
		switch rep.End.Winner {
		case "attacker":
			atkWins++
		case "defender":
			defWins++
		}
		totalRounds += rep.Summary.TotalRounds
	}
	elapsed := time.Since(start)
	fmt.Printf("scenario          : %s\n", sc.Name)
	fmt.Printf("battles           : %d\n", n)
	fmt.Printf("attacker_winrate  : %.1f%%\n", 100.0*float64(atkWins)/float64(n))
	fmt.Printf("defender_winrate  : %.1f%%\n", 100.0*float64(defWins)/float64(n))
	fmt.Printf("avg_rounds        : %.1f\n", float64(totalRounds)/float64(n))
	fmt.Printf("total_time        : %s\n", elapsed)
	fmt.Printf("avg_time_per_battle: %s\n", elapsed/time.Duration(n))
}

func unitFromEntry(ue unitEntry) *engine.Unit {
	side := tilemap.SideAttacker
	if ue.Side == "defender" {
		side = tilemap.SideDefender
	}
	kind := engine.UnitFleet
	if ue.Kind == "building" {
		kind = engine.UnitBuilding
	}
	tile := tilemap.Position{Col: ue.Tile[0], Row: ue.Tile[1]}
	footprint := []tilemap.Position{tile}
	if len(ue.Footprint) > 0 {
		footprint = make([]tilemap.Position, 0, len(ue.Footprint))
		for _, f := range ue.Footprint {
			footprint = append(footprint, tilemap.Position{Col: f[0], Row: f[1]})
		}
	}
	stacks := make([]*engine.Stack, 0, len(ue.Stacks))
	for i, se := range ue.Stacks {
		stacks = append(stacks, &engine.Stack{
			ID: fmt.Sprintf("%s_s%d", ue.ID, i),
			GridRow: se.GridRow, GridCol: se.GridCol,
			HullClass: engine.HullClass(se.HullClass),
			ShipCount: se.ShipCount,
			WeaponClass: engine.WeaponClass(se.WeaponClass),
			WeaponName: se.WeaponName,
			DamageType: engine.DamageType(se.DamageType),
			ArmorType: engine.ArmorType(se.ArmorType),
			MinRange: se.MinRange, MaxRange: se.MaxRange,
			Cooldown: se.Cooldown,
			BaseDamageMin: se.BaseDamageMin, BaseDamageMax: se.BaseDamageMax,
			HitChance: se.HitChance,
			ShieldHP: se.ShieldHP, MaxShield: se.ShieldHP,
			StructHP: se.StructHP, MaxStruct: se.StructHP,
			Accuracy: se.Accuracy, Dodge: se.Dodge,
		})
	}
	return &engine.Unit{
		ID: ue.ID, Side: side, Kind: kind,
		Tile: tile, Footprint: footprint,
		Speed: ue.Speed, Movement: ue.Movement,
		TargetCmd: ue.TargetCmd,
		Phalanx: engine.Phalanx{Stacks: stacks},
	}
}
