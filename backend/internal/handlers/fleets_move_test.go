package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
)

// TestMoveFleet_RealTravelTime verifies that the previously-stub MoveFleet
// endpoint now schedules an arrival time proportional to distance, not
// instant ("now() + interval '1 second'"). It uses a far-away destination
// so the travel time is unambiguously larger than the old hardcoded second.
func TestMoveFleet_RealTravelTime(t *testing.T) {
	if database.DB == nil {
		if err := database.InitSupabase(); err != nil {
			t.Fatalf("init db: %v", err)
		}
	}

	playerID := "00000000-0000-0000-0000-000000009101"
	planetID := "00000000-0000-0000-aaaa-000000009101"
	fleetID := "00000000-0000-0000-bbbb-000000009101"
	defer func() {
		database.DB.Exec("DELETE FROM fleet_stacks WHERE fleet_id = $1", fleetID)
		database.DB.Exec("DELETE FROM fleets WHERE id = $1", fleetID)
		database.DB.Exec("DELETE FROM resources WHERE planet_id = $1", planetID)
		database.DB.Exec("DELETE FROM planets WHERE id = $1", planetID)
		database.DB.Exec("DELETE FROM players WHERE id = $1", playerID)
	}()

	if _, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'movetest', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID); err != nil {
		t.Fatalf("create player: %v", err)
	}
	if _, err := database.DB.Exec(`
		INSERT INTO planets (id, player_id, name, is_homeworld, position_x, position_y, created_at, updated_at)
		VALUES ($1, $2, 'Origin', true, 1500, 1500, now(), now())
		ON CONFLICT (id) DO NOTHING
	`, planetID, playerID); err != nil {
		t.Fatalf("create planet: %v", err)
	}
	if _, err := database.DB.Exec(`
		INSERT INTO resources (planet_id, metal, he3, gold)
		VALUES ($1, 0, 100000, 0)
	`, planetID); err != nil {
		t.Fatalf("create resources: %v", err)
	}

	// Stationed fleet with a single non-empty stack so MoveFleet doesn't reject.
	if _, err := database.DB.Exec(`
		INSERT INTO fleets (id, player_id, name, formation, targeting_command, status, planet_id, position_x, position_y)
		VALUES ($1, $2, 'Mover', 'phalanx', 'max_attack', 'stationed', $3, 1500, 1500)
	`, fleetID, playerID, planetID); err != nil {
		t.Fatalf("create fleet: %v", err)
	}

	// Need a ship_design + stack with ship_count > 0 so the handler accepts the move.
	var hullID int
	if err := database.DB.QueryRow(`SELECT id FROM hull_types WHERE name='weikes_i'`).Scan(&hullID); err != nil {
		t.Fatalf("hull lookup: %v", err)
	}
	var designID string
	if err := database.DB.QueryRow(`
		INSERT INTO ship_designs (player_id, name, hull_type_id, modules_json,
			total_shield, total_structure, total_defense, total_agility, total_movement,
			total_storage, attack_power, weapon_range_min, weapon_range_max, volume_used,
			he3_per_round, metal_cost, he3_cost, gold_cost, build_time_seconds)
		VALUES ($1, 'mover_design', $2, '[]', 100, 100, 5, 5, 30, 50, 20, 1, 2, 5, 2, 100, 50, 25, 10)
		RETURNING id
	`, playerID, hullID).Scan(&designID); err != nil {
		t.Fatalf("design: %v", err)
	}
	if _, err := database.DB.Exec(`
		INSERT INTO fleet_stacks (fleet_id, ship_design_id, grid_row, grid_col, ship_count)
		VALUES ($1, $2, 0, 0, 5)
	`, fleetID, designID); err != nil {
		t.Fatalf("stack: %v", err)
	}

	// Far away destination — should produce a multi-second ETA.
	body, _ := json.Marshal(map[string]int{
		"destination_x": 1900,
		"destination_y": 1900,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/fleets/"+fleetID+"/move", bytes.NewReader(body))
	req.SetPathValue("id", fleetID)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	before := time.Now()
	MoveFleet(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}

	var resp struct {
		Status     string    `json:"status"`
		ArrivalAt  time.Time `json:"arrival_at"`
		DestinationX int     `json:"destination_x"`
		DestinationY int     `json:"destination_y"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v body=%s", err, rr.Body.String())
	}
	if resp.Status != "traveling" {
		t.Errorf("status=%s, want traveling", resp.Status)
	}
	if resp.DestinationX != 1900 || resp.DestinationY != 1900 {
		t.Errorf("destination wrong: %+v", resp)
	}

	travelSeconds := resp.ArrivalAt.Sub(before).Seconds()
	if travelSeconds < 2 {
		t.Errorf("expected travel time > 2s for 400-unit hop, got %.2fs (still using stub?)", travelSeconds)
	}
}
