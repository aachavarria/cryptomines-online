package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
)

// TestCombatReports_RoundsJSONPersisted runs an instance attempt and
// then verifies the resulting combat_reports row carries a non-empty
// rounds_json so the UI battle playback has data to replay.
func TestCombatReports_RoundsJSONPersisted(t *testing.T) {
	if database.DB == nil {
		if err := database.InitSupabase(); err != nil {
			t.Fatalf("init db: %v", err)
		}
	}

	playerID := "00000000-0000-0000-0000-00000000a201"
	planetID := "00000000-0000-0000-aaaa-00000000a201"
	fleetID := "00000000-0000-0000-bbbb-00000000a201"
	defer func() {
		database.DB.Exec("DELETE FROM combat_reports WHERE attacker_id = $1", playerID)
		database.DB.Exec("DELETE FROM fleet_stacks WHERE fleet_id = $1", fleetID)
		database.DB.Exec("DELETE FROM fleets WHERE id = $1", fleetID)
		database.DB.Exec("DELETE FROM ships WHERE player_id = $1", playerID)
		database.DB.Exec("DELETE FROM ship_designs WHERE player_id = $1", playerID)
		database.DB.Exec("DELETE FROM resources WHERE planet_id = $1", planetID)
		database.DB.Exec("DELETE FROM planets WHERE id = $1", planetID)
		database.DB.Exec("DELETE FROM players WHERE id = $1", playerID)
	}()

	if _, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, level, created_at, updated_at)
		VALUES ($1, $2, 'roundsjson', 5, now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID); err != nil {
		t.Fatalf("player: %v", err)
	}
	if _, err := database.DB.Exec(`
		INSERT INTO planets (id, player_id, name, is_homeworld, position_x, position_y, created_at, updated_at)
		VALUES ($1, $2, 'rj-home', true, 1700, 1700, now(), now())
	`, planetID, playerID); err != nil {
		t.Fatalf("planet: %v", err)
	}

	var hullID int
	if err := database.DB.QueryRow(`SELECT id FROM hull_types WHERE name='weikes_i'`).Scan(&hullID); err != nil {
		t.Fatalf("hull: %v", err)
	}
	var designID string
	if err := database.DB.QueryRow(`
		INSERT INTO ship_designs (player_id, name, hull_type_id, modules_json,
			total_shield, total_structure, total_defense, total_agility, total_movement,
			total_storage, attack_power, weapon_range_min, weapon_range_max, volume_used,
			he3_per_round, metal_cost, he3_cost, gold_cost, build_time_seconds)
		VALUES ($1, 'rj_design', $2, '[]', 200, 200, 5, 5, 30, 50, 30, 1, 2, 5, 2, 100, 50, 25, 10)
		RETURNING id
	`, playerID, hullID).Scan(&designID); err != nil {
		t.Fatalf("design: %v", err)
	}
	if _, err := database.DB.Exec(`
		INSERT INTO fleets (id, player_id, name, formation, targeting_command, status, planet_id, position_x, position_y)
		VALUES ($1, $2, 'rj_fleet', 'phalanx', 'max_attack', 'stationed', $3, 1700, 1700)
	`, fleetID, playerID, planetID); err != nil {
		t.Fatalf("fleet: %v", err)
	}
	if _, err := database.DB.Exec(`
		INSERT INTO fleet_stacks (fleet_id, ship_design_id, grid_row, grid_col, ship_count)
		VALUES ($1, $2, 0, 0, 80)
	`, fleetID, designID); err != nil {
		t.Fatalf("stack: %v", err)
	}

	var instanceID string
	if err := database.DB.QueryRow(`
		SELECT id FROM instances WHERE type='normal' AND difficulty=1 LIMIT 1
	`).Scan(&instanceID); err != nil {
		t.Fatalf("instance: %v", err)
	}

	body := []byte(fmt.Sprintf(`{"fleet_ids":["%s"]}`, fleetID))
	req := httptest.NewRequest(http.MethodPost, "/api/instances/"+instanceID+"/attempt", newReader(body))
	req.SetPathValue("id", instanceID)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	AttemptInstance(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("attempt status=%d body=%s", rr.Code, rr.Body.String())
	}

	var reportID string
	var roundsJSON []byte
	var totalRounds int
	if err := database.DB.QueryRow(`
		SELECT id, rounds_json, total_rounds FROM combat_reports
		WHERE attacker_id = $1 ORDER BY created_at DESC LIMIT 1
	`, playerID).Scan(&reportID, &roundsJSON, &totalRounds); err != nil {
		t.Fatalf("report: %v", err)
	}
	if len(roundsJSON) == 0 || string(roundsJSON) == "[]" {
		t.Fatalf("rounds_json was empty (%q); battle playback will have nothing to render", string(roundsJSON))
	}

	var rounds []map[string]any
	if err := json.Unmarshal(roundsJSON, &rounds); err != nil {
		t.Fatalf("rounds_json not valid JSON: %v", err)
	}
	if totalRounds > 0 && len(rounds) == 0 {
		t.Errorf("totalRounds=%d but rounds_json has %d entries", totalRounds, len(rounds))
	}
	if len(rounds) > 0 {
		first := rounds[0]
		if _, ok := first["RoundNumber"]; !ok {
			t.Errorf("first round missing RoundNumber field, got keys: %v", keys(first))
		}
		if _, ok := first["Attacks"]; !ok {
			t.Errorf("first round missing Attacks field, got keys: %v", keys(first))
		}
	}
}

func keys(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
