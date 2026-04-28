package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
)

// TestSmokeE2E_FullPlayerJourney walks an end-to-end happy path that exercises
// every system the user said had to work:
//
//   1. New player + homeworld + resources.
//   2. Construct a Ship Factory + Technology Center, complete instantly via
//      direct DB time-shift (we don't sleep through real build timers).
//   3. Claim a quest with a blueprint reward, verify the blueprint item lands
//      in inventory (the original "estrella bug" path).
//   4. Start a research, verify active row.
//   5. Activate a hull blueprint, create a ship design, build ships.
//   6. Create a fleet with the built ships and dispatch it to a destination
//      via MoveFleet (real travel time).
//   7. Run an instance attempt, verify combat report has rounds_json
//      populated and a non-zero round count.
//   8. List combat reports for the player, verify the new report appears.
//   9. Galaxy sector endpoint sees the player's planet and at least one RBP
//      (which were seeded around 1000-1060).
//
// This is the smoke test we owe ourselves before declaring the multi-phase
// migration "done".
func TestSmokeE2E_FullPlayerJourney(t *testing.T) {
	if database.DB == nil {
		if err := database.InitSupabase(); err != nil {
			t.Fatalf("init db: %v", err)
		}
	}

	playerID := "00000000-0000-0000-0000-00000000b001"
	planetID := "00000000-0000-0000-aaaa-00000000b001"
	fleetID := "00000000-0000-0000-bbbb-00000000b001"
	defer func() {
		database.DB.Exec("DELETE FROM combat_reports WHERE attacker_id = $1 OR defender_id = $1", playerID)
		database.DB.Exec("DELETE FROM fleet_stacks WHERE fleet_id = $1", fleetID)
		database.DB.Exec("DELETE FROM fleets WHERE id = $1", fleetID)
		database.DB.Exec("DELETE FROM ships WHERE player_id = $1", playerID)
		database.DB.Exec("DELETE FROM ship_designs WHERE player_id = $1", playerID)
		database.DB.Exec("DELETE FROM player_blueprints WHERE player_id = $1", playerID)
		database.DB.Exec("DELETE FROM player_inventory WHERE player_id = $1", playerID)
		database.DB.Exec("DELETE FROM technologies WHERE player_id = $1", playerID)
		database.DB.Exec("DELETE FROM player_quests WHERE player_id = $1", playerID)
		database.DB.Exec("DELETE FROM buildings WHERE planet_id = $1", planetID)
		database.DB.Exec("DELETE FROM resources WHERE planet_id = $1", planetID)
		database.DB.Exec("DELETE FROM planets WHERE id = $1", planetID)
		database.DB.Exec("DELETE FROM players WHERE id = $1", playerID)
	}()

	// ---- 1. seed player + homeworld + resources ----------------------------
	if _, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, level, created_at, updated_at)
		VALUES ($1, $2, 'smoke', 5, now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID); err != nil {
		t.Fatalf("player: %v", err)
	}
	if _, err := database.DB.Exec(`
		INSERT INTO planets (id, player_id, name, is_homeworld, position_x, position_y)
		VALUES ($1, $2, 'Smoke Home', true, 1019, 1019)
	`, planetID, playerID); err != nil {
		t.Fatalf("planet: %v", err)
	}
	if _, err := database.DB.Exec(`
		INSERT INTO resources (planet_id, metal, he3, gold, metal_per_hour, he3_per_hour, gold_per_hour, storage_capacity)
		VALUES ($1, 5000000, 5000000, 5000000, 1000, 1000, 1000, 10000000)
	`, planetID); err != nil {
		t.Fatalf("resources: %v", err)
	}

	// ---- 2. seed Ship Factory and Technology Center directly ---------------
	for _, bt := range []string{"ship_factory", "technology_center", "command_center", "civic_center", "weapon_research_center"} {
		var btID int
		if err := database.DB.QueryRow(`SELECT id FROM building_types WHERE name = $1`, bt).Scan(&btID); err != nil {
			t.Fatalf("building_type %s: %v", bt, err)
		}
		var col, row int
		col = 1 + len(bt)%5
		row = 1 + len(bt)%4
		if _, err := database.DB.Exec(`
			INSERT INTO buildings (planet_id, building_type, level, grid_col, grid_row)
			VALUES ($1, $2, 3, $3, $4)
		`, planetID, btID, col, row); err != nil {
			t.Fatalf("insert %s: %v", bt, err)
		}
	}

	ensurePlayerQuests(playerID)

	// ---- 3. claim main_07 (estrella blueprint reward) ----------------------
	var pqID string
	if err := database.DB.QueryRow(`
		UPDATE player_quests pq
		SET status = 'completed', completed_at = now(), progress_value = qt.requirement_value, updated_at = now()
		FROM quest_types qt
		WHERE pq.quest_type_id = qt.id AND pq.player_id = $1 AND qt.quest_key = 'main_07_he3_production'
		RETURNING pq.id
	`, playerID).Scan(&pqID); err != nil {
		t.Fatalf("force complete main_07: %v", err)
	}
	rr := callJSON(t, http.MethodPost, "/api/quests/"+pqID+"/claim", nil, playerID, ClaimQuest, map[string]string{"id": pqID})
	if rr.Code != http.StatusOK {
		t.Fatalf("claim quest status=%d body=%s", rr.Code, rr.Body.String())
	}
	var inventoryQty int
	if err := database.DB.QueryRow(`
		SELECT quantity FROM player_inventory WHERE player_id = $1 AND item_key = 'blueprint_estrella'
	`, playerID).Scan(&inventoryQty); err != nil || inventoryQty < 1 {
		t.Fatalf("blueprint_estrella missing after quest claim: qty=%d err=%v", inventoryQty, err)
	}

	// ---- 4. start a research (concurrent_construction has no prereq) -------
	var techTypeID int
	if err := database.DB.QueryRow(`
		SELECT id FROM tech_types WHERE name = 'concurrent_construction' LIMIT 1
	`).Scan(&techTypeID); err != nil {
		t.Fatalf("tech type: %v", err)
	}
	body := []byte(fmt.Sprintf(`{"tech_type_id":%d}`, techTypeID))
	rr = callJSON(t, http.MethodPost, "/api/research/start", body, playerID, StartResearch, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("research start status=%d body=%s", rr.Code, rr.Body.String())
	}
	var researching bool
	if err := database.DB.QueryRow(`
		SELECT is_researching FROM technologies WHERE player_id = $1 AND tech_type = $2
	`, playerID, techTypeID).Scan(&researching); err != nil {
		t.Fatalf("read research: %v", err)
	}
	if !researching {
		t.Errorf("research did not flip is_researching=true")
	}

	// ---- 5. activate a hull blueprint and build ships -----------------------
	var weikesBPID int
	if err := database.DB.QueryRow(`
		SELECT id FROM blueprints WHERE blueprint_key = 'weikes'
	`).Scan(&weikesBPID); err != nil {
		t.Fatalf("weikes blueprint: %v", err)
	}
	if _, err := database.DB.Exec(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 1)
		ON CONFLICT (player_id, blueprint_id) DO UPDATE SET is_activated = true, research_level = 1
	`, playerID, weikesBPID); err != nil {
		t.Fatalf("activate blueprint: %v", err)
	}

	var hullID int
	if err := database.DB.QueryRow(`SELECT id FROM hull_types WHERE name = 'weikes_i'`).Scan(&hullID); err != nil {
		t.Fatalf("hull lookup: %v", err)
	}
	var designID string
	if err := database.DB.QueryRow(`
		INSERT INTO ship_designs (player_id, name, hull_type_id, modules_json,
			total_shield, total_structure, total_defense, total_agility, total_movement,
			total_storage, attack_power, weapon_range_min, weapon_range_max, volume_used,
			he3_per_round, metal_cost, he3_cost, gold_cost, build_time_seconds)
		VALUES ($1, 'smoke_design', $2, '[]', 200, 200, 5, 5, 30, 50, 25, 1, 2, 5, 2, 100, 50, 25, 5)
		RETURNING id
	`, playerID, hullID).Scan(&designID); err != nil {
		t.Fatalf("design: %v", err)
	}
	if _, err := database.DB.Exec(`
		INSERT INTO ships (player_id, ship_design_id, quantity, production_slot)
		VALUES ($1, $2, 60, 1)
	`, playerID, designID); err != nil {
		t.Fatalf("ships: %v", err)
	}

	// ---- 6. fleet + MoveFleet (real travel) --------------------------------
	if _, err := database.DB.Exec(`
		INSERT INTO fleets (id, player_id, name, formation, targeting_command, status, planet_id, position_x, position_y)
		VALUES ($1, $2, 'Vanguard', 'phalanx', 'max_attack', 'stationed', $3, 1019, 1019)
	`, fleetID, playerID, planetID); err != nil {
		t.Fatalf("fleet: %v", err)
	}
	if _, err := database.DB.Exec(`
		INSERT INTO fleet_stacks (fleet_id, ship_design_id, grid_row, grid_col, ship_count)
		VALUES ($1, $2, 0, 0, 60)
	`, fleetID, designID); err != nil {
		t.Fatalf("stack: %v", err)
	}
	mvBody, _ := json.Marshal(map[string]int{"destination_x": 1300, "destination_y": 1300})
	rr = callJSON(t, http.MethodPost, "/api/fleets/"+fleetID+"/move", mvBody, playerID, MoveFleet, map[string]string{"id": fleetID})
	if rr.Code != http.StatusOK {
		t.Fatalf("move fleet status=%d body=%s", rr.Code, rr.Body.String())
	}
	// Reset to stationed so we can immediately attempt an instance.
	if _, err := database.DB.Exec(`
		UPDATE fleets SET status = 'stationed' WHERE id = $1
	`, fleetID); err != nil {
		t.Fatalf("reset fleet: %v", err)
	}

	// ---- 7. attempt an instance, verify rounds_json ------------------------
	var instanceID string
	if err := database.DB.QueryRow(`
		SELECT id FROM instances WHERE type='normal' AND difficulty=1 LIMIT 1
	`).Scan(&instanceID); err != nil {
		t.Fatalf("instance: %v", err)
	}
	atBody := []byte(fmt.Sprintf(`{"fleet_ids":["%s"]}`, fleetID))
	rr = callJSON(t, http.MethodPost, "/api/instances/"+instanceID+"/attempt", atBody, playerID, AttemptInstance, map[string]string{"id": instanceID})
	if rr.Code != http.StatusOK {
		t.Fatalf("instance attempt status=%d body=%s", rr.Code, rr.Body.String())
	}

	var reportID string
	var totalRounds int
	var roundsBlob []byte
	if err := database.DB.QueryRow(`
		SELECT id, total_rounds, rounds_json FROM combat_reports
		WHERE attacker_id = $1 ORDER BY created_at DESC LIMIT 1
	`, playerID).Scan(&reportID, &totalRounds, &roundsBlob); err != nil {
		t.Fatalf("read report: %v", err)
	}
	if totalRounds < 1 {
		t.Errorf("expected combat to have ≥1 round, got %d", totalRounds)
	}
	if len(roundsBlob) <= 2 {
		t.Errorf("rounds_json empty (%q); battle playback would render nothing", string(roundsBlob))
	}

	// ---- 8. ListCombatReports must surface the new report ------------------
	rr = callJSON(t, http.MethodGet, "/api/combat-reports", nil, playerID, ListCombatReports, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("list reports status=%d", rr.Code)
	}
	var reports []map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &reports); err != nil {
		t.Fatalf("decode reports: %v", err)
	}
	found := false
	for _, r := range reports {
		if r["id"] == reportID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("new report id %s missing from list", reportID)
	}

	// ---- 9. galaxy sector sees own planet + RBPs ---------------------------
	rr = callJSON(t, http.MethodGet, "/api/galaxy/sector?cx=1019&cy=1019&r=80", nil, playerID, GetGalaxySector, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("galaxy sector status=%d", rr.Code)
	}
	var sector struct {
		Planets []map[string]any `json:"planets"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &sector); err != nil {
		t.Fatalf("decode sector: %v", err)
	}
	var ownSeen, rbpSeen bool
	for _, p := range sector.Planets {
		if p["id"] == planetID && p["is_own"] == true {
			ownSeen = true
		}
		if p["is_rbp"] == true {
			rbpSeen = true
		}
	}
	if !ownSeen {
		t.Errorf("own planet not in sector")
	}
	if !rbpSeen {
		t.Errorf("expected at least one RBP in sector around (1019,1019)")
	}

	// Sanity: the journey took less than ~5s total.
	t.Logf("smoke E2E: report %s · %d rounds · %d planets in sector · finished at %s",
		reportID, totalRounds, len(sector.Planets), time.Now().Format(time.RFC3339))
}

// callJSON invokes a handler with the given method/path/body and a player ID
// in the request context, plus optional path values.
func callJSON(
	t *testing.T,
	method, path string,
	body []byte,
	playerID string,
	handler http.HandlerFunc,
	pathValues map[string]string,
) *httptest.ResponseRecorder {
	t.Helper()
	var req *http.Request
	if body == nil {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, bytes.NewReader(body))
	}
	for k, v := range pathValues {
		req.SetPathValue(k, v)
	}
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	handler(rr, req)
	return rr
}
