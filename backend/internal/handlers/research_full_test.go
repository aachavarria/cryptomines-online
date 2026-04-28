package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
)

// =============================================================================
// Research test helpers
// =============================================================================

// setupResearchTestDB initializes the shared DB connection.
func setupResearchTestDB(t *testing.T) {
	if database.DB == nil {
		if err := database.InitSupabase(); err != nil {
			t.Fatalf("Failed to init database: %v", err)
		}
	}
}

// cleanupResearchPlayer removes all data for a research-test player. Order:
// child tables first, then parents.
func cleanupResearchPlayer(t *testing.T, playerID string) {
	t.Helper()
	database.DB.Exec("DELETE FROM player_inventory WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM technologies WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM player_quests WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM daily_quest_progress WHERE player_id = $1", playerID)
	database.DB.Exec(`DELETE FROM buildings WHERE planet_id IN (SELECT id FROM planets WHERE player_id = $1)`, playerID)
	database.DB.Exec(`DELETE FROM resources WHERE planet_id IN (SELECT id FROM planets WHERE player_id = $1)`, playerID)
	database.DB.Exec("DELETE FROM planets WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM players WHERE id = $1", playerID)
}

// seedResearchPlayer creates a player + planet + resources row + a Technology
// Center (so research is unlocked). Returns (planetID, techCenterID).
//
// gold: how much gold to seed on the homeworld (other resources start at 0).
// techCenterLevel: pass 1+ to satisfy StartResearch's "must have a Technology
// Center" gate; pass 0 to specifically test that gate.
func seedResearchPlayer(t *testing.T, playerID string, gold int64, techCenterLevel int) (string, string) {
	t.Helper()

	if _, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID, "rfull-"+playerID[len(playerID)-8:]); err != nil {
		t.Fatalf("create player: %v", err)
	}

	tail := playerID[len(playerID)-4:]
	parsed, err := strconv.ParseInt(tail, 16, 32)
	if err != nil {
		t.Fatalf("parse position tail: %v", err)
	}
	posX := 9000 + int(parsed%4000)
	posY := 9000 + int((parsed/4000)%4000)

	planetID := "00000000-0000-0000-bbbb-" + playerID[len(playerID)-12:]
	if _, err := database.DB.Exec(`
		INSERT INTO planets (id, player_id, name, is_homeworld, position_x, position_y, created_at, updated_at)
		VALUES ($1, $2, 'RFullHome', true, $3, $4, now(), now())
		ON CONFLICT (id) DO NOTHING
	`, planetID, playerID, posX, posY); err != nil {
		t.Fatalf("create planet: %v", err)
	}

	// resources: no created_at column. Use updated_at and last_collected_at.
	if _, err := database.DB.Exec(`
		INSERT INTO resources (planet_id, metal, he3, gold)
		VALUES ($1, 0, 0, $2)
		ON CONFLICT (planet_id) DO NOTHING
	`, planetID, gold); err != nil {
		t.Fatalf("create resources: %v", err)
	}
	// In case the row pre-existed (it shouldn't), ensure gold is what we want.
	if _, err := database.DB.Exec(
		`UPDATE resources SET gold = $1, updated_at = now() WHERE planet_id = $2`,
		gold, planetID,
	); err != nil {
		t.Fatalf("set gold: %v", err)
	}

	techCenterID := ""
	if techCenterLevel > 0 {
		var btID int
		if err := database.DB.QueryRow(
			`SELECT id FROM building_types WHERE name = 'technology_center'`,
		).Scan(&btID); err != nil {
			t.Fatalf("technology_center building_type missing: %v", err)
		}
		if err := database.DB.QueryRow(`
			INSERT INTO buildings (planet_id, building_type, grid_col, grid_row, level, is_upgrading)
			VALUES ($1, $2, 4, 4, $3, false)
			RETURNING id
		`, planetID, btID, techCenterLevel).Scan(&techCenterID); err != nil {
			t.Fatalf("create technology_center: %v", err)
		}
	}

	return planetID, techCenterID
}

// invokeStartResearch dispatches StartResearch with a given tech_type_id.
func invokeStartResearch(t *testing.T, playerID string, techTypeID int) (int, []byte) {
	t.Helper()
	body, _ := json.Marshal(map[string]any{"tech_type_id": techTypeID})
	req := httptest.NewRequest(http.MethodPost, "/api/research/start", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	StartResearch(rr, req)
	return rr.Code, rr.Body.Bytes()
}

// invokeCancelResearch dispatches CancelResearch with a given tech_type_id.
func invokeCancelResearch(t *testing.T, playerID string, techTypeID int) (int, []byte) {
	t.Helper()
	body, _ := json.Marshal(map[string]any{"tech_type_id": techTypeID})
	req := httptest.NewRequest(http.MethodPost, "/api/research/cancel", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	CancelResearch(rr, req)
	return rr.Code, rr.Body.Bytes()
}

// invokeSpeedupResearch dispatches SpeedupResearch.
func invokeSpeedupResearch(t *testing.T, playerID string, techTypeID, minutes int) (int, []byte) {
	t.Helper()
	body, _ := json.Marshal(map[string]any{
		"tech_type_id":    techTypeID,
		"speedup_minutes": minutes,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/research/speedup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	SpeedupResearch(rr, req)
	return rr.Code, rr.Body.Bytes()
}

// fetchTechTypeID returns the id of a tech_type by name.
func fetchTechTypeID(t *testing.T, name string) int {
	t.Helper()
	var id int
	if err := database.DB.QueryRow(`SELECT id FROM tech_types WHERE name = $1`, name).Scan(&id); err != nil {
		t.Fatalf("lookup tech_type %s: %v", name, err)
	}
	return id
}

// =============================================================================
// (a) ListResearch returns trees with seeded techs
// =============================================================================

// TestListResearch_ContainsExpectedTrees calls ListResearch and verifies all
// 7 seeded trees appear with at least one tech each.
func TestListResearch_ContainsExpectedTrees(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003000"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 0, 0)

	req := httptest.NewRequest(http.MethodGet, "/api/research", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	ListResearch(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("ListResearch: status=%d body=%s", rr.Code, rr.Body.String())
	}

	var body struct {
		Trees []struct {
			Tree  string `json:"tree"`
			Techs []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"techs"`
		} `json:"trees"`
		Active *struct {
			Tree string `json:"tree"`
		} `json:"active"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v body=%s", err, rr.Body.String())
	}
	want := map[string]bool{
		"logistics_construction": false,
		"planetary_defense":      false,
		"ballistics_science":     false,
		"directional_science":    false,
		"missile_science":        false,
		"ship_based_science":     false,
		"ship_defense_science":   false,
	}
	for _, tr := range body.Trees {
		if _, ok := want[tr.Tree]; ok {
			want[tr.Tree] = true
		}
		if len(tr.Techs) == 0 {
			t.Errorf("tree %q has zero techs", tr.Tree)
		}
	}
	for treeName, found := range want {
		if !found {
			t.Errorf("expected tree %q in response, not found", treeName)
		}
	}
	if body.Active != nil {
		t.Errorf("expected no active research for fresh player, got tree=%q", body.Active.Tree)
	}
}

// =============================================================================
// (b) GetResearchTree for each real tree
// =============================================================================

// TestGetResearchTree_AllTrees calls GetResearchTree for each of the 7 valid
// tree names and verifies a non-empty tech list.
func TestGetResearchTree_AllTrees(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003001"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 0, 0)

	trees := []string{
		"logistics_construction",
		"planetary_defense",
		"ballistics_science",
		"directional_science",
		"missile_science",
		"ship_based_science",
		"ship_defense_science",
	}
	for _, treeName := range trees {
		t.Run(treeName, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/research/trees/"+treeName, nil)
			req.SetPathValue("tree", treeName)
			ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()
			GetResearchTree(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
			}

			var body struct {
				Tree  string `json:"tree"`
				Techs []struct {
					Name string `json:"name"`
					Tree string `json:"tree"`
				} `json:"techs"`
			}
			if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if body.Tree != treeName {
				t.Errorf("response tree=%q want %q", body.Tree, treeName)
			}
			if len(body.Techs) == 0 {
				t.Errorf("tree %q returned zero techs", treeName)
			}
			for _, tch := range body.Techs {
				if tch.Tree != treeName {
					t.Errorf("tech %q has tree=%q in tree-scoped response", tch.Name, tch.Tree)
				}
			}
		})
	}
}

// TestGetResearchTree_InvalidTree returns 400 for an unknown tree.
func TestGetResearchTree_InvalidTree(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003002"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 0, 0)

	req := httptest.NewRequest(http.MethodGet, "/api/research/trees/not_a_tree", nil)
	req.SetPathValue("tree", "not_a_tree")
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	GetResearchTree(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("got status=%d want 400 for invalid tree", rr.Code)
	}
}

// =============================================================================
// (c) StartResearch happy path
// =============================================================================

// TestStartResearch_ConcurrentConstruction_HappyPath researches the
// 'concurrent_construction' tech (max_level=1, base_cost_gold=1000, no prereqs)
// and verifies a row was upserted into 'technologies' with is_researching=true.
func TestStartResearch_ConcurrentConstruction_HappyPath(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003100"
	defer cleanupResearchPlayer(t, playerID)
	planetID, _ := seedResearchPlayer(t, playerID, 5000, 1)

	techID := fetchTechTypeID(t, "concurrent_construction")

	code, body := invokeStartResearch(t, playerID, techID)
	if code != http.StatusOK {
		t.Fatalf("StartResearch: status=%d body=%s", code, string(body))
	}

	// Verify row in 'technologies' (the actual active-research storage).
	var level int
	var isResearching bool
	var finishAt *time.Time
	if err := database.DB.QueryRow(`
		SELECT level, is_researching, research_finish_at
		FROM technologies
		WHERE player_id = $1 AND tech_type = $2
	`, playerID, techID).Scan(&level, &isResearching, &finishAt); err != nil {
		t.Fatalf("read technologies: %v", err)
	}
	if !isResearching {
		t.Errorf("is_researching=false after StartResearch, want true")
	}
	if finishAt == nil {
		t.Errorf("research_finish_at is null after StartResearch")
	}
	// Level on the row stays at currentLevel until applyCompletedResearch runs.
	// concurrent_construction starts at 0, so row.level == 0 during research.
	if level != 0 {
		t.Errorf("level after start: got %d, want 0 (advances on completion)", level)
	}

	// Verify gold was deducted.
	var gold int64
	database.DB.QueryRow(`SELECT gold FROM resources WHERE planet_id = $1`, planetID).Scan(&gold)
	if gold != 4000 { // 5000 - 1000
		t.Errorf("gold after deduct: got %d, want 4000", gold)
	}
}

// TestStartResearch_NoTechCenter_Returns409 verifies the Technology Center gate.
// A player without a Tech Center cannot start research.
func TestStartResearch_NoTechCenter_Returns409(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003101"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 5000, 0) // techCenterLevel=0

	techID := fetchTechTypeID(t, "concurrent_construction")
	code, body := invokeStartResearch(t, playerID, techID)
	if code != http.StatusConflict {
		t.Errorf("got status=%d want 409 for no tech center; body=%s", code, string(body))
	}
}

// TestStartResearch_AlreadyResearchingInTree_Returns409 starts one research,
// then attempts to start a second tech in the same tree → expects 409.
func TestStartResearch_AlreadyResearchingInTree_Returns409(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003102"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 100000, 1)

	first := fetchTechTypeID(t, "concurrent_construction")
	if code, body := invokeStartResearch(t, playerID, first); code != http.StatusOK {
		t.Fatalf("first start: status=%d body=%s", code, string(body))
	}

	// ship_building_boost is in the same tree (logistics_construction) and has no prereq.
	second := fetchTechTypeID(t, "ship_building_boost")
	code, _ := invokeStartResearch(t, playerID, second)
	if code != http.StatusConflict {
		t.Errorf("second start (same tree): got %d want 409", code)
	}
}

// =============================================================================
// (d) StartResearch rejects when prerequisites unmet
// =============================================================================

// TestStartResearch_PrerequisiteUnmet_Returns409 tries to research
// 'construction_boost' (requires concurrent_construction lvl 1) without having
// researched the prerequisite. Expects 409.
func TestStartResearch_PrerequisiteUnmet_Returns409(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003200"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 100000, 1)

	techID := fetchTechTypeID(t, "construction_boost")
	code, body := invokeStartResearch(t, playerID, techID)
	if code != http.StatusConflict {
		t.Errorf("got status=%d want 409 for unmet prereq; body=%s", code, string(body))
	}

	// Sanity: no row in 'technologies' should exist for this tech.
	var n int
	database.DB.QueryRow(`
		SELECT COUNT(*) FROM technologies WHERE player_id = $1 AND tech_type = $2
	`, playerID, techID).Scan(&n)
	if n != 0 {
		t.Errorf("expected 0 rows for failed start, got %d", n)
	}
}

// TestStartResearch_PrerequisiteMet_Succeeds verifies the inverse: after
// researching concurrent_construction to lvl 1, construction_boost can start.
func TestStartResearch_PrerequisiteMet_Succeeds(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003201"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 100000, 1)

	prereqID := fetchTechTypeID(t, "concurrent_construction")
	// Pre-populate the prereq at level 1 (skip the actual research timer).
	if _, err := database.DB.Exec(`
		INSERT INTO technologies (player_id, tech_type, level, is_researching)
		VALUES ($1, $2, 1, false)
		ON CONFLICT (player_id, tech_type) DO UPDATE SET level=1, is_researching=false, research_finish_at=NULL
	`, playerID, prereqID); err != nil {
		t.Fatalf("seed prereq: %v", err)
	}

	techID := fetchTechTypeID(t, "construction_boost")
	code, body := invokeStartResearch(t, playerID, techID)
	if code != http.StatusOK {
		t.Errorf("start with prereq met: status=%d body=%s", code, string(body))
	}
}

// =============================================================================
// (e) StartResearch rejects when insufficient resources
// =============================================================================

// TestStartResearch_InsufficientGold_Returns409 attempts research with 0 gold
// (concurrent_construction costs 1000). Expects 409 conflict and no
// technologies row.
func TestStartResearch_InsufficientGold_Returns409(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003300"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 0, 1) // 0 gold, has tech center

	techID := fetchTechTypeID(t, "concurrent_construction")
	code, body := invokeStartResearch(t, playerID, techID)
	if code != http.StatusConflict {
		t.Errorf("got status=%d want 409 for insufficient gold; body=%s", code, string(body))
	}

	var n int
	database.DB.QueryRow(`
		SELECT COUNT(*) FROM technologies WHERE player_id = $1 AND tech_type = $2 AND is_researching = true
	`, playerID, techID).Scan(&n)
	if n != 0 {
		t.Errorf("expected no active research after failed start, got %d rows", n)
	}
}

// TestStartResearch_BarelyInsufficientGold_Returns409 verifies cost-comparison
// boundary: 999 gold (vs cost 1000) still rejects.
func TestStartResearch_BarelyInsufficientGold_Returns409(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003301"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 999, 1)

	techID := fetchTechTypeID(t, "concurrent_construction")
	code, _ := invokeStartResearch(t, playerID, techID)
	if code != http.StatusConflict {
		t.Errorf("got status=%d want 409 for 999 gold (cost=1000)", code)
	}
}

// =============================================================================
// (f) CancelResearch deactivates the active row
// =============================================================================

// TestCancelResearch_RemovesActiveFlag starts research, cancels it, and verifies
// the row remains but is_researching=false and research_finish_at IS NULL.
//
// NOTE: The current handler does NOT refund the spent gold. This test asserts
// the actual behaviour (no refund). If refund support is added later, update
// the assertion accordingly.
func TestCancelResearch_RemovesActiveFlag(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003400"
	defer cleanupResearchPlayer(t, playerID)
	planetID, _ := seedResearchPlayer(t, playerID, 5000, 1)

	techID := fetchTechTypeID(t, "concurrent_construction")
	if code, body := invokeStartResearch(t, playerID, techID); code != http.StatusOK {
		t.Fatalf("start: status=%d body=%s", code, string(body))
	}

	// Sanity: row is researching now.
	var goldAfterStart int64
	database.DB.QueryRow(`SELECT gold FROM resources WHERE planet_id = $1`, planetID).Scan(&goldAfterStart)
	if goldAfterStart != 4000 {
		t.Fatalf("gold after start: got %d want 4000", goldAfterStart)
	}

	if code, body := invokeCancelResearch(t, playerID, techID); code != http.StatusOK {
		t.Fatalf("cancel: status=%d body=%s", code, string(body))
	}

	// Row should still exist (handler updates, doesn't delete) but
	// is_researching=false and research_finish_at IS NULL (chk_research_consistency).
	var isResearching bool
	var finishAt *time.Time
	if err := database.DB.QueryRow(`
		SELECT is_researching, research_finish_at
		FROM technologies
		WHERE player_id = $1 AND tech_type = $2
	`, playerID, techID).Scan(&isResearching, &finishAt); err != nil {
		t.Fatalf("read row after cancel: %v", err)
	}
	if isResearching {
		t.Errorf("is_researching=true after cancel")
	}
	if finishAt != nil {
		t.Errorf("research_finish_at not null after cancel: %v", finishAt)
	}

	// Verify no refund (current behaviour).
	var goldAfterCancel int64
	database.DB.QueryRow(`SELECT gold FROM resources WHERE planet_id = $1`, planetID).Scan(&goldAfterCancel)
	if goldAfterCancel != goldAfterStart {
		t.Errorf("gold after cancel changed: %d -> %d (no refund expected)", goldAfterStart, goldAfterCancel)
	}
}

// TestCancelResearch_NoActive_Returns409 cancels with no active research.
func TestCancelResearch_NoActive_Returns409(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003401"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 0, 1)

	techID := fetchTechTypeID(t, "concurrent_construction")
	code, _ := invokeCancelResearch(t, playerID, techID)
	if code != http.StatusConflict {
		t.Errorf("got status=%d want 409 for no active research", code)
	}
}

// =============================================================================
// (g) SpeedupResearch reduces remaining time
// =============================================================================

// TestSpeedupResearch_ReducesFinishTime starts a long-running research
// (ballistics_base, 117s base time, max_level 10), pushes research_finish_at
// far into the future, then calls speedup with 30 minutes. Verifies
// research_finish_at moved earlier by ~30 minutes and vouchers_spent=3.
//
// NOTE: The handler currently doesn't actually consume an inventory item even
// though the GDD references "Speedup" cards — it just reports voucher_cost.
// This test seeds a marker row in player_inventory anyway so any future code
// that DOES consume an item has a stack to use, and we assert the time
// reduction (the actual observable effect).
func TestSpeedupResearch_ReducesFinishTime(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003500"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 100000, 1)

	// Use ballistics_base (10000 gold base, but cheap enough for 100k).
	techID := fetchTechTypeID(t, "ballistics_base")

	// Start research.
	if code, body := invokeStartResearch(t, playerID, techID); code != http.StatusOK {
		t.Fatalf("start: status=%d body=%s", code, string(body))
	}

	// Force research_finish_at to exactly 1 hour from now so the speedup math
	// is deterministic regardless of the seeded base_time_seconds.
	target := time.Now().Add(1 * time.Hour)
	if _, err := database.DB.Exec(`
		UPDATE technologies
		SET research_finish_at = $1, updated_at = now()
		WHERE player_id = $2 AND tech_type = $3 AND is_researching = true
	`, target, playerID, techID); err != nil {
		t.Fatalf("set finish_at: %v", err)
	}

	// Seed a notional speedup item in inventory so the *invariant* holds for
	// future implementations that may decrement quantity. We use construction_card
	// (which is a real seeded item_key) as a proxy stack.
	if _, err := database.DB.Exec(`
		INSERT INTO player_inventory (player_id, item_key, quantity)
		VALUES ($1, 'construction_card', 5)
		ON CONFLICT (player_id, item_key) DO UPDATE SET quantity = 5
	`, playerID); err != nil {
		t.Fatalf("seed inventory: %v", err)
	}

	// Speedup by 30 min → expect 3 vouchers spent and finish_at moved earlier.
	code, body := invokeSpeedupResearch(t, playerID, techID, 30)
	if code != http.StatusOK {
		t.Fatalf("speedup: status=%d body=%s", code, string(body))
	}

	var resp struct {
		Technology struct {
			ResearchFinishAt *time.Time `json:"research_finish_at"`
			IsResearching    bool       `json:"is_researching"`
		} `json:"technology"`
		VouchersSpent int `json:"vouchers_spent"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("decode: %v body=%s", err, string(body))
	}
	if resp.VouchersSpent != 3 {
		t.Errorf("vouchers_spent: got %d want 3 (3 per 30 min)", resp.VouchersSpent)
	}

	// The new finish_at should be ~30 min after now (1 hour - 30 min).
	var actualFinish *time.Time
	if err := database.DB.QueryRow(`
		SELECT research_finish_at FROM technologies
		WHERE player_id = $1 AND tech_type = $2
	`, playerID, techID).Scan(&actualFinish); err != nil {
		t.Fatalf("read finish_at: %v", err)
	}
	if actualFinish == nil {
		t.Fatalf("finish_at is null")
	}
	expected := target.Add(-30 * time.Minute)
	delta := actualFinish.Sub(expected)
	if delta < -10*time.Second || delta > 10*time.Second {
		t.Errorf("finish_at after speedup: %v, expected ~%v (delta=%v)",
			actualFinish, expected, delta)
	}
}

// TestSpeedupResearch_NegativeMinutes_Returns400 verifies validation.
func TestSpeedupResearch_NegativeMinutes_Returns400(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003501"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 100000, 1)

	techID := fetchTechTypeID(t, "concurrent_construction")
	if code, _ := invokeStartResearch(t, playerID, techID); code != http.StatusOK {
		t.Fatalf("start failed")
	}
	code, _ := invokeSpeedupResearch(t, playerID, techID, -5)
	if code != http.StatusBadRequest {
		t.Errorf("speedup -5min: got status=%d want 400", code)
	}
}

// TestSpeedupResearch_NoActiveResearch_Returns409 verifies error path.
func TestSpeedupResearch_NoActiveResearch_Returns409(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003502"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 100000, 1)

	techID := fetchTechTypeID(t, "concurrent_construction")
	code, _ := invokeSpeedupResearch(t, playerID, techID, 10)
	if code != http.StatusConflict {
		t.Errorf("speedup with no active: got status=%d want 409", code)
	}
}

// TestSpeedupResearch_PastFinishTime_AutoCompletes pushes speedup minutes large
// enough that finish_at would land in the past — handler clamps to now() and
// applyCompletedResearch flips level up.
func TestSpeedupResearch_PastFinishTime_AutoCompletes(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003503"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 100000, 1)

	techID := fetchTechTypeID(t, "concurrent_construction")
	if code, _ := invokeStartResearch(t, playerID, techID); code != http.StatusOK {
		t.Fatalf("start failed")
	}

	// Force finish_at to 1 minute from now, then speedup by 60 minutes.
	if _, err := database.DB.Exec(`
		UPDATE technologies
		SET research_finish_at = $1
		WHERE player_id = $2 AND tech_type = $3
	`, time.Now().Add(1*time.Minute), playerID, techID); err != nil {
		t.Fatalf("set finish_at: %v", err)
	}

	code, body := invokeSpeedupResearch(t, playerID, techID, 60)
	if code != http.StatusOK {
		t.Fatalf("speedup: status=%d body=%s", code, string(body))
	}
	_ = body

	// After speedup with overflow, the research should be auto-completed:
	// level should be 1 (was 0), is_researching=false.
	var level int
	var isResearching bool
	if err := database.DB.QueryRow(`
		SELECT level, is_researching FROM technologies
		WHERE player_id = $1 AND tech_type = $2
	`, playerID, techID).Scan(&level, &isResearching); err != nil {
		t.Fatalf("read tech: %v", err)
	}
	if isResearching {
		t.Errorf("is_researching=true after large speedup; expected auto-complete")
	}
	if level != 1 {
		t.Errorf("level after auto-complete: got %d want 1", level)
	}
}

// =============================================================================
// (h) GetActiveResearch returns active research row
// =============================================================================

// TestGetActiveResearch_ReturnsStartedRow starts research and verifies
// GetActiveResearch returns a slice with one entry containing the tech name.
func TestGetActiveResearch_ReturnsStartedRow(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003600"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 100000, 1)

	techID := fetchTechTypeID(t, "concurrent_construction")
	if code, _ := invokeStartResearch(t, playerID, techID); code != http.StatusOK {
		t.Fatalf("start failed")
	}

	// Push finish_at far into future so applyCompletedResearch doesn't auto-complete.
	if _, err := database.DB.Exec(`
		UPDATE technologies SET research_finish_at = $1
		WHERE player_id = $2 AND tech_type = $3
	`, time.Now().Add(2*time.Hour), playerID, techID); err != nil {
		t.Fatalf("push finish_at: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/research/active", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	GetActiveResearch(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var rows []struct {
		TechName      string `json:"tech_name"`
		Tree          string `json:"tree"`
		IsResearching bool   `json:"is_researching"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &rows); err != nil {
		t.Fatalf("decode: %v body=%s", err, rr.Body.String())
	}
	if len(rows) != 1 {
		t.Fatalf("got %d active rows, want 1", len(rows))
	}
	if rows[0].TechName != "concurrent_construction" {
		t.Errorf("tech_name: got %q want concurrent_construction", rows[0].TechName)
	}
	if rows[0].Tree != "logistics_construction" {
		t.Errorf("tree: got %q want logistics_construction", rows[0].Tree)
	}
	if !rows[0].IsResearching {
		t.Errorf("is_researching=false in active list (expected true)")
	}
}

// TestGetActiveResearch_NoActive_ReturnsEmpty verifies the empty-list path.
func TestGetActiveResearch_NoActive_ReturnsEmpty(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003601"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 0, 1)

	req := httptest.NewRequest(http.MethodGet, "/api/research/active", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	GetActiveResearch(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var rows []map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &rows); err != nil {
		t.Fatalf("decode: %v body=%s", err, rr.Body.String())
	}
	if len(rows) != 0 {
		t.Errorf("got %d active rows for fresh player, want 0", len(rows))
	}
}

// =============================================================================
// extra: ListResearch.Active reflects active research
// =============================================================================

// TestListResearch_ActiveFieldPopulated starts research and verifies the
// top-level "active" field of /api/research is non-null and points at the
// right tech.
func TestListResearch_ActiveFieldPopulated(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003700"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 100000, 1)

	techID := fetchTechTypeID(t, "concurrent_construction")
	if code, _ := invokeStartResearch(t, playerID, techID); code != http.StatusOK {
		t.Fatalf("start failed")
	}
	// Keep it active.
	database.DB.Exec(`
		UPDATE technologies SET research_finish_at = $1
		WHERE player_id = $2 AND tech_type = $3
	`, time.Now().Add(2*time.Hour), playerID, techID)

	req := httptest.NewRequest(http.MethodGet, "/api/research", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	ListResearch(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var body struct {
		Active *struct {
			TechName string `json:"tech_name"`
			Tree     string `json:"tree"`
		} `json:"active"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Active == nil {
		t.Fatalf("expected active to be non-null after start")
	}
	if body.Active.TechName != "concurrent_construction" {
		t.Errorf("active.tech_name: got %q want concurrent_construction", body.Active.TechName)
	}
}

// =============================================================================
// guard against accidentally undocumented behaviour
// =============================================================================

// TestStartResearch_UnknownTechTypeID_Returns400 verifies the malformed-input
// path: a tech_type_id that doesn't exist returns 400.
func TestStartResearch_UnknownTechTypeID_Returns400(t *testing.T) {
	setupResearchTestDB(t)
	playerID := "00000000-0000-0000-0000-000000003800"
	defer cleanupResearchPlayer(t, playerID)
	seedResearchPlayer(t, playerID, 100000, 1)

	code, body := invokeStartResearch(t, playerID, 999999)
	if code != http.StatusBadRequest {
		t.Errorf("got status=%d want 400 for unknown tech_type_id; body=%s", code, string(body))
	}
}

