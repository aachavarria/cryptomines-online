package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
)

// =============================================================================
// Quest test helpers
// =============================================================================

// cleanupQuestFullPlayer removes ALL data for a quest-full test player.
// Order: child tables first (FK respect), then parents.
func cleanupQuestFullPlayer(t *testing.T, playerID string) {
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

// seedQuestFullPlayer creates a player + homeworld + resources row + initialises
// the player_quests table. Uses unique position_x/position_y derived from the
// last 4 hex chars of the player UUID so two tests never collide.
func seedQuestFullPlayer(t *testing.T, playerID string) string {
	t.Helper()

	if _, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID, "qfull-"+playerID[len(playerID)-8:]); err != nil {
		t.Fatalf("create player: %v", err)
	}

	// Derive a deterministic but unique position from the last 4 hex digits.
	tail := playerID[len(playerID)-4:]
	parsed, err := strconv.ParseInt(tail, 16, 32)
	if err != nil {
		t.Fatalf("parse position tail: %v", err)
	}
	// Offset so we land in a region used by no other test files.
	posX := 5000 + int(parsed%4000)
	posY := 5000 + int((parsed/4000)%4000)

	planetID := "00000000-0000-0000-aaaa-" + playerID[len(playerID)-12:]
	if _, err := database.DB.Exec(`
		INSERT INTO planets (id, player_id, name, is_homeworld, position_x, position_y, created_at, updated_at)
		VALUES ($1, $2, 'QFullHome', true, $3, $4, now(), now())
		ON CONFLICT (id) DO NOTHING
	`, planetID, playerID, posX, posY); err != nil {
		t.Fatalf("create planet: %v", err)
	}

	// resources table has no created_at column — only updated_at, last_collected_at, last_warehouse_update.
	if _, err := database.DB.Exec(`
		INSERT INTO resources (planet_id, metal, he3, gold)
		VALUES ($1, 0, 0, 0)
		ON CONFLICT (planet_id) DO NOTHING
	`, planetID); err != nil {
		t.Fatalf("create resources: %v", err)
	}

	// Populate player_quests rows so we can mark/claim them.
	ensurePlayerQuests(playerID)
	return planetID
}

// markQuestCompletedFull forces a player_quest into 'completed'. It honours the
// chk_status_consistency constraint by setting completed_at NOT NULL.
// Returns the player_quest UUID.
func markQuestCompletedFull(t *testing.T, playerID, questKey string) string {
	t.Helper()
	var pqID string
	err := database.DB.QueryRow(`
		UPDATE player_quests pq
		SET status = 'completed',
		    progress_value = qt.requirement_value,
		    completed_at = now(),
		    claimed_at = NULL,
		    updated_at = now()
		FROM quest_types qt
		WHERE pq.quest_type_id = qt.id
		  AND pq.player_id = $1
		  AND qt.quest_key = $2
		RETURNING pq.id
	`, playerID, questKey).Scan(&pqID)
	if err != nil {
		t.Fatalf("mark quest %s completed: %v", questKey, err)
	}
	return pqID
}

// invokeClaimQuest dispatches ClaimQuest with the given quest UUID and decoded body.
func invokeClaimQuest(t *testing.T, playerID, pqID string) (int, map[string]any) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/quests/"+pqID+"/claim", nil)
	req.SetPathValue("id", pqID)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	ClaimQuest(rr, req)
	var body map[string]any
	if rr.Body.Len() > 0 {
		_ = json.Unmarshal(rr.Body.Bytes(), &body)
	}
	return rr.Code, body
}

// =============================================================================
// (a) ClaimQuest happy path for every active main quest
// =============================================================================

// TestClaimQuest_AllMainQuests_HappyPath iterates every active main quest in
// chain order, marks it completed, claims it, and verifies:
//   - response status is 200
//   - DB status flips to 'claimed' with claimed_at populated
//   - homeworld resources receive metal/he3/gold rewards (cumulative)
//   - inventory items materialise for both 'item' and 'blueprint' rewards
//   - the next quest (linked via prerequisite_quest_id) flips locked->available
func TestClaimQuest_AllMainQuests_HappyPath(t *testing.T) {
	setupQuestTestDB(t)
	playerID := "00000000-0000-0000-0000-000000001000"
	defer cleanupQuestFullPlayer(t, playerID)
	planetID := seedQuestFullPlayer(t, playerID)

	rows, err := database.DB.Query(`
		SELECT qt.id, qt.quest_key, qt.reward_metal, qt.reward_he3, qt.reward_gold, qt.reward_item_json
		FROM quest_types qt
		WHERE qt.category = 'main' AND qt.is_active = true
		ORDER BY qt.chain_order
	`)
	if err != nil {
		t.Fatalf("list main quests: %v", err)
	}
	type questRow struct {
		ID          int
		Key         string
		RewardMetal int64
		RewardHe3   int64
		RewardGold  int64
		Items       []byte
	}
	var quests []questRow
	for rows.Next() {
		var q questRow
		if err := rows.Scan(&q.ID, &q.Key, &q.RewardMetal, &q.RewardHe3, &q.RewardGold, &q.Items); err != nil {
			t.Fatalf("scan: %v", err)
		}
		quests = append(quests, q)
	}
	rows.Close()
	if len(quests) == 0 {
		t.Skip("no active main quests seeded; can't run loop")
	}

	// Track cumulative reward expectations for the homeworld.
	var wantMetal, wantHe3, wantGold int64

	for _, q := range quests {
		t.Run(q.Key, func(t *testing.T) {
			pqID := markQuestCompletedFull(t, playerID, q.Key)
			code, _ := invokeClaimQuest(t, playerID, pqID)
			if code != http.StatusOK {
				t.Fatalf("claim %s: got status=%d", q.Key, code)
			}

			// Verify status='claimed' with both timestamps populated (constraint-honouring).
			var status string
			var completedAt, claimedAt *string
			if err := database.DB.QueryRow(`
				SELECT status, completed_at::text, claimed_at::text
				FROM player_quests WHERE id = $1
			`, pqID).Scan(&status, &completedAt, &claimedAt); err != nil {
				t.Fatalf("read claimed quest: %v", err)
			}
			if status != "claimed" {
				t.Errorf("status: got %q want claimed", status)
			}
			if completedAt == nil || claimedAt == nil {
				t.Errorf("expected completed_at AND claimed_at not null after claim (chk_status_consistency)")
			}

			// Verify cumulative resource rewards on homeworld.
			wantMetal += q.RewardMetal
			wantHe3 += q.RewardHe3
			wantGold += q.RewardGold
			var gotMetal, gotHe3, gotGold int64
			if err := database.DB.QueryRow(
				`SELECT metal, he3, gold FROM resources WHERE planet_id = $1`, planetID,
			).Scan(&gotMetal, &gotHe3, &gotGold); err != nil {
				t.Fatalf("read resources: %v", err)
			}
			if gotMetal != wantMetal || gotHe3 != wantHe3 || gotGold != wantGold {
				t.Errorf("homeworld resources after %s: got (m=%d he3=%d g=%d) want (m=%d he3=%d g=%d)",
					q.Key, gotMetal, gotHe3, gotGold, wantMetal, wantHe3, wantGold)
			}

			// Verify item rewards landed in inventory.
			if len(q.Items) > 0 && string(q.Items) != "[]" {
				var items []map[string]any
				if err := json.Unmarshal(q.Items, &items); err != nil {
					t.Fatalf("parse reward items: %v", err)
				}
				for _, it := range items {
					itemType, _ := it["type"].(string)
					var inventoryKey string
					switch itemType {
					case "item":
						inventoryKey, _ = it["item_key"].(string)
					case "blueprint":
						bk, _ := it["blueprint_key"].(string)
						inventoryKey = "blueprint_" + bk
					case "commander":
						ck, _ := it["commander_key"].(string)
						inventoryKey = "commander_" + ck
					default:
						continue
					}
					if inventoryKey == "" {
						continue
					}
					var qty int
					if err := database.DB.QueryRow(`
						SELECT quantity FROM player_inventory
						WHERE player_id = $1 AND item_key = $2
					`, playerID, inventoryKey).Scan(&qty); err != nil {
						t.Errorf("inventory missing %s for %s: %v", inventoryKey, q.Key, err)
						continue
					}
					if qty < 1 {
						t.Errorf("inventory key %s qty %d (want >=1)", inventoryKey, qty)
					}
				}
			}

			// Verify next quest in chain (if any) was unlocked: locked -> available.
			var nextStatus string
			err := database.DB.QueryRow(`
				SELECT pq.status FROM player_quests pq
				JOIN quest_types nxt ON nxt.id = pq.quest_type_id
				WHERE pq.player_id = $1 AND nxt.prerequisite_quest_id = $2 AND nxt.is_active = true
				LIMIT 1
			`, playerID, q.ID).Scan(&nextStatus)
			if err == nil && nextStatus == "locked" {
				t.Errorf("next quest after %s still locked; should be available", q.Key)
			}
		})
	}
}

// =============================================================================
// (b) ClaimQuest for daily-tier and side quest categories (sanity)
// =============================================================================

// TestClaimQuest_OneSideQuest_RewardsApplied claims one tier-1 side quest and
// verifies resources are awarded.
func TestClaimQuest_OneSideQuest_RewardsApplied(t *testing.T) {
	setupQuestTestDB(t)
	playerID := "00000000-0000-0000-0000-000000001100"
	defer cleanupQuestFullPlayer(t, playerID)
	planetID := seedQuestFullPlayer(t, playerID)

	var sideKey string
	err := database.DB.QueryRow(`
		SELECT quest_key FROM quest_types
		WHERE category = 'side' AND is_active = true AND prerequisite_quest_id IS NULL
		ORDER BY quest_key LIMIT 1
	`).Scan(&sideKey)
	if err != nil {
		t.Skipf("no side quest seeded: %v", err)
	}

	pqID := markQuestCompletedFull(t, playerID, sideKey)
	code, _ := invokeClaimQuest(t, playerID, pqID)
	if code != http.StatusOK {
		t.Fatalf("claim side %s: status=%d", sideKey, code)
	}

	var metal, he3, gold int64
	if err := database.DB.QueryRow(
		`SELECT metal, he3, gold FROM resources WHERE planet_id = $1`, planetID,
	).Scan(&metal, &he3, &gold); err != nil {
		t.Fatalf("read resources: %v", err)
	}
	if metal == 0 && he3 == 0 && gold == 0 {
		t.Errorf("side quest %s gave no resource reward (m/he3/g all 0)", sideKey)
	}
}

// TestClaimQuest_ClaimDailyTier_Bronze claims the bronze daily tier after
// forcibly setting daily_points >= 10. Daily quest *types* don't go through
// ClaimQuest itself; tier rewards go through ClaimDailyTier.
func TestClaimQuest_ClaimDailyTier_Bronze(t *testing.T) {
	setupQuestTestDB(t)
	playerID := "00000000-0000-0000-0000-000000001101"
	defer cleanupQuestFullPlayer(t, playerID)
	seedQuestFullPlayer(t, playerID)

	// Ensure daily progress row exists, then force enough points for bronze (10).
	ensureDailyProgress(playerID)
	if _, err := database.DB.Exec(`
		UPDATE daily_quest_progress
		SET daily_points = 25
		WHERE player_id = $1 AND quest_date = CURRENT_DATE
	`, playerID); err != nil {
		t.Fatalf("set points: %v", err)
	}

	body := []byte(`{"tier":"bronze"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/quests/daily/claim-tier",
		newReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	ClaimDailyTier(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("claim bronze tier: status=%d body=%s", rr.Code, rr.Body.String())
	}

	// Re-claim should return 409 (already claimed).
	rr2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/api/quests/daily/claim-tier",
		newReader(body))
	req2.Header.Set("Content-Type", "application/json")
	req2 = req2.WithContext(ctx)
	ClaimDailyTier(rr2, req2)
	if rr2.Code != http.StatusConflict {
		t.Errorf("re-claim bronze: got status=%d, want 409", rr2.Code)
	}
}

// TestClaimQuest_DailyTier_InsufficientPoints verifies 409 when
// daily_points < required.
func TestClaimQuest_DailyTier_InsufficientPoints(t *testing.T) {
	setupQuestTestDB(t)
	playerID := "00000000-0000-0000-0000-000000001102"
	defer cleanupQuestFullPlayer(t, playerID)
	seedQuestFullPlayer(t, playerID)

	ensureDailyProgress(playerID)
	// Default points are 0, gold tier needs 50.
	body := []byte(`{"tier":"gold"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/quests/daily/claim-tier",
		newReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	ClaimDailyTier(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("got status=%d, want 409 for insufficient points", rr.Code)
	}
}

// =============================================================================
// (c) ClaimQuest rejects when status != 'completed' (returns 409)
// =============================================================================

// TestClaimQuest_NotCompleted_Returns409 attempts to claim an 'available' quest
// and expects 409 Conflict (the handler rejects non-completed claims).
func TestClaimQuest_NotCompleted_Returns409(t *testing.T) {
	setupQuestTestDB(t)
	playerID := "00000000-0000-0000-0000-000000001200"
	defer cleanupQuestFullPlayer(t, playerID)
	seedQuestFullPlayer(t, playerID)

	// Get an 'available' quest (the first main quest is auto-set to available).
	var pqID string
	err := database.DB.QueryRow(`
		SELECT pq.id FROM player_quests pq
		JOIN quest_types qt ON pq.quest_type_id = qt.id
		WHERE pq.player_id = $1 AND pq.status = 'available' AND qt.category = 'main'
		LIMIT 1
	`, playerID).Scan(&pqID)
	if err != nil {
		t.Fatalf("find available quest: %v", err)
	}

	code, _ := invokeClaimQuest(t, playerID, pqID)
	if code != http.StatusConflict {
		t.Errorf("claim non-completed quest: got status=%d, want 409", code)
	}
}

// TestClaimQuest_Locked_Returns409 attempts to claim a 'locked' quest and
// expects 409 (status != 'completed').
func TestClaimQuest_Locked_Returns409(t *testing.T) {
	setupQuestTestDB(t)
	playerID := "00000000-0000-0000-0000-000000001201"
	defer cleanupQuestFullPlayer(t, playerID)
	seedQuestFullPlayer(t, playerID)

	var pqID string
	err := database.DB.QueryRow(`
		SELECT id FROM player_quests
		WHERE player_id = $1 AND status = 'locked'
		LIMIT 1
	`, playerID).Scan(&pqID)
	if err != nil {
		t.Skipf("no locked quest available: %v", err)
	}
	code, _ := invokeClaimQuest(t, playerID, pqID)
	if code != http.StatusConflict {
		t.Errorf("claim locked quest: got status=%d, want 409", code)
	}
}

// TestClaimQuest_AlreadyClaimed_Returns409 marks a quest 'claimed' then tries
// to claim again — should reject with 409.
func TestClaimQuest_AlreadyClaimed_Returns409(t *testing.T) {
	setupQuestTestDB(t)
	playerID := "00000000-0000-0000-0000-000000001202"
	defer cleanupQuestFullPlayer(t, playerID)
	seedQuestFullPlayer(t, playerID)

	pqID := markQuestCompletedFull(t, playerID, "main_01_collecting_resources")
	code, _ := invokeClaimQuest(t, playerID, pqID)
	if code != http.StatusOK {
		t.Fatalf("first claim: status=%d", code)
	}
	code2, _ := invokeClaimQuest(t, playerID, pqID)
	if code2 != http.StatusConflict {
		t.Errorf("second claim: got status=%d, want 409", code2)
	}
}

// =============================================================================
// (d) ClaimQuest rejects when quest belongs to another player (404)
// =============================================================================

// TestClaimQuest_OtherPlayer_Returns404 creates two players, marks player A's
// quest completed, then tries to claim it as player B. Expects 404.
func TestClaimQuest_OtherPlayer_Returns404(t *testing.T) {
	setupQuestTestDB(t)
	playerA := "00000000-0000-0000-0000-000000001300"
	playerB := "00000000-0000-0000-0000-000000001301"
	defer cleanupQuestFullPlayer(t, playerA)
	defer cleanupQuestFullPlayer(t, playerB)
	seedQuestFullPlayer(t, playerA)
	seedQuestFullPlayer(t, playerB)

	pqA := markQuestCompletedFull(t, playerA, "main_01_collecting_resources")

	code, _ := invokeClaimQuest(t, playerB, pqA)
	if code != http.StatusNotFound {
		t.Errorf("claim other player's quest: got status=%d, want 404", code)
	}

	// Sanity: A can still claim it.
	codeA, _ := invokeClaimQuest(t, playerA, pqA)
	if codeA != http.StatusOK {
		t.Errorf("owner re-claim after foreign attempt: got %d, want 200", codeA)
	}
}

// TestClaimQuest_NonExistentID_Returns404 uses a syntactically valid UUID that
// doesn't exist; expects 404.
func TestClaimQuest_NonExistentID_Returns404(t *testing.T) {
	setupQuestTestDB(t)
	playerID := "00000000-0000-0000-0000-000000001302"
	defer cleanupQuestFullPlayer(t, playerID)
	seedQuestFullPlayer(t, playerID)

	bogus := "00000000-0000-0000-0000-aaaaaaaaaaaa"
	code, _ := invokeClaimQuest(t, playerID, bogus)
	if code != http.StatusNotFound {
		t.Errorf("claim non-existent quest: got %d, want 404", code)
	}
}

// =============================================================================
// (e) SyncQuests: build he3_extractor lvl 1 -> main_07 should auto-complete
// =============================================================================

// TestSyncQuests_He3ExtractorCompletesMain07 builds an he3_extractor at level 1
// and verifies that POST /api/quests/sync flips main_07_he3_production from
// 'available' to 'completed' (it's a build_building quest with target=he3_extractor,
// requirement_value=1).
func TestSyncQuests_He3ExtractorCompletesMain07(t *testing.T) {
	setupQuestTestDB(t)
	playerID := "00000000-0000-0000-0000-000000001400"
	defer cleanupQuestFullPlayer(t, playerID)
	planetID := seedQuestFullPlayer(t, playerID)

	// Look up he3_extractor building_type id.
	var btID int
	if err := database.DB.QueryRow(
		`SELECT id FROM building_types WHERE name = 'he3_extractor'`,
	).Scan(&btID); err != nil {
		t.Skipf("he3_extractor building_type missing: %v", err)
	}

	// Insert an he3_extractor at level 1 on the homeworld (not upgrading).
	if _, err := database.DB.Exec(`
		INSERT INTO buildings (planet_id, building_type, grid_col, grid_row, level, is_upgrading, created_at, updated_at)
		VALUES ($1, $2, 5, 5, 1, false, now(), now())
	`, planetID, btID); err != nil {
		t.Fatalf("create he3_extractor: %v", err)
	}

	// SyncQuests only auto-completes quests in 'available'. main_07 starts
	// 'locked' (waits on prerequisite chain). Force it 'available' first to
	// match the runtime path.
	if _, err := database.DB.Exec(`
		UPDATE player_quests pq
		SET status = 'available', completed_at = NULL, claimed_at = NULL, updated_at = now()
		FROM quest_types qt
		WHERE pq.quest_type_id = qt.id
		  AND pq.player_id = $1
		  AND qt.quest_key = 'main_07_he3_production'
	`, playerID); err != nil {
		t.Fatalf("set main_07 available: %v", err)
	}

	// Call SyncQuests.
	req := httptest.NewRequest(http.MethodPost, "/api/quests/sync", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	SyncQuests(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("SyncQuests: status=%d body=%s", rr.Code, rr.Body.String())
	}

	// Verify status flipped to completed (chk_status_consistency: completed_at NOT NULL).
	var status string
	var completedAt *string
	if err := database.DB.QueryRow(`
		SELECT pq.status, pq.completed_at::text
		FROM player_quests pq
		JOIN quest_types qt ON pq.quest_type_id = qt.id
		WHERE pq.player_id = $1 AND qt.quest_key = 'main_07_he3_production'
	`, playerID).Scan(&status, &completedAt); err != nil {
		t.Fatalf("read main_07: %v", err)
	}
	if status != "completed" {
		t.Errorf("main_07 status: got %q, want completed", status)
	}
	if completedAt == nil {
		t.Errorf("main_07 completed_at is null after sync")
	}
}

// TestSyncQuests_NoQualifyingBuildings_NoOp verifies that SyncQuests is a no-op
// when no buildings match: status stays as-is.
func TestSyncQuests_NoQualifyingBuildings_NoOp(t *testing.T) {
	setupQuestTestDB(t)
	playerID := "00000000-0000-0000-0000-000000001401"
	defer cleanupQuestFullPlayer(t, playerID)
	seedQuestFullPlayer(t, playerID)

	req := httptest.NewRequest(http.MethodPost, "/api/quests/sync", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	SyncQuests(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("SyncQuests: status=%d body=%s", rr.Code, rr.Body.String())
	}

	// Make sure no quest accidentally became 'completed' without a matching building.
	var n int
	if err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM player_quests
		WHERE player_id = $1 AND status = 'completed'
	`, playerID).Scan(&n); err != nil {
		t.Fatalf("count completed: %v", err)
	}
	if n != 0 {
		t.Errorf("after sync without buildings, %d quests are 'completed'", n)
	}
}

// TestSyncQuests_BuildingBelowRequirement_NoComplete builds a Civic Center at
// level 1 and verifies a hypothetical quest that needs level 2 doesn't complete.
// Picks any 'upgrade_building' quest with requirement_value >= 2 and ensures it
// stays 'available' when only level-1 building exists.
func TestSyncQuests_BuildingBelowRequirement_NoComplete(t *testing.T) {
	setupQuestTestDB(t)
	playerID := "00000000-0000-0000-0000-000000001402"
	defer cleanupQuestFullPlayer(t, playerID)
	planetID := seedQuestFullPlayer(t, playerID)

	// Find an active 'build_building' quest with requirement_value >= 2.
	var questKey, target string
	var requirement int
	err := database.DB.QueryRow(`
		SELECT qt.quest_key, qt.requirement_target, qt.requirement_value
		FROM quest_types qt
		WHERE qt.is_active = true
		  AND qt.requirement_type = 'build_building'
		  AND qt.requirement_value >= 2
		LIMIT 1
	`).Scan(&questKey, &target, &requirement)
	if err != nil {
		// Try upgrade_building instead.
		err = database.DB.QueryRow(`
			SELECT qt.quest_key, qt.requirement_target, qt.requirement_value
			FROM quest_types qt
			WHERE qt.is_active = true
			  AND qt.requirement_type = 'upgrade_building'
			  AND qt.requirement_value >= 2
			LIMIT 1
		`).Scan(&questKey, &target, &requirement)
		if err != nil {
			t.Skipf("no quest with requirement_value>=2 found")
		}
	}

	// Force the quest to 'available'.
	if _, err := database.DB.Exec(`
		UPDATE player_quests pq
		SET status = 'available', completed_at = NULL, claimed_at = NULL
		FROM quest_types qt
		WHERE pq.quest_type_id = qt.id AND pq.player_id = $1 AND qt.quest_key = $2
	`, playerID, questKey); err != nil {
		t.Fatalf("set available: %v", err)
	}

	// Build target at level 1 (below requirement).
	var btID int
	if err := database.DB.QueryRow(`SELECT id FROM building_types WHERE name = $1`, target).Scan(&btID); err != nil {
		t.Skipf("building_type %s missing: %v", target, err)
	}
	if _, err := database.DB.Exec(`
		INSERT INTO buildings (planet_id, building_type, grid_col, grid_row, level, is_upgrading)
		VALUES ($1, $2, 6, 6, 1, false)
	`, planetID, btID); err != nil {
		t.Fatalf("create building: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/quests/sync", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	SyncQuests(rr, req)

	// Quest should still be 'available' since building level (1) < requirement (>=2).
	var status string
	if err := database.DB.QueryRow(`
		SELECT pq.status FROM player_quests pq
		JOIN quest_types qt ON pq.quest_type_id = qt.id
		WHERE pq.player_id = $1 AND qt.quest_key = $2
	`, playerID, questKey).Scan(&status); err != nil {
		t.Fatalf("read status: %v", err)
	}
	if status == "completed" {
		t.Errorf("quest %s (req level=%d) completed despite only level-1 building", questKey, requirement)
	}
}

// =============================================================================
// (f) ListQuests + GetDailyQuests basic shape
// =============================================================================

// TestListQuests_Shape verifies the response has main_quests, side_quests,
// and current_main_quest fields, with at least one main quest in 'available'.
func TestListQuests_Shape(t *testing.T) {
	setupQuestTestDB(t)
	playerID := "00000000-0000-0000-0000-000000001500"
	defer cleanupQuestFullPlayer(t, playerID)
	seedQuestFullPlayer(t, playerID)

	req := httptest.NewRequest(http.MethodGet, "/api/quests", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	ListQuests(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("ListQuests: status=%d body=%s", rr.Code, rr.Body.String())
	}

	var body struct {
		MainQuests []struct {
			QuestKey string `json:"quest_key"`
			Status   string `json:"status"`
			Category string `json:"category"`
		} `json:"main_quests"`
		SideQuests []struct {
			QuestKey string `json:"quest_key"`
			Status   string `json:"status"`
		} `json:"side_quests"`
		CurrentMainQuest *struct {
			QuestKey string `json:"quest_key"`
		} `json:"current_main_quest"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v body=%s", err, rr.Body.String())
	}
	if len(body.MainQuests) == 0 {
		t.Errorf("expected non-empty main_quests")
	}
	if body.CurrentMainQuest == nil {
		t.Errorf("expected current_main_quest to point at first non-claimed")
	}
	// Side quests should be present (tier 1 with no prereq -> available).
	if len(body.SideQuests) == 0 {
		t.Errorf("expected non-empty side_quests")
	}
}

// TestListQuests_AfterClaim_CurrentAdvances claims main_01 and verifies
// current_main_quest is no longer main_01 (it should advance to next).
func TestListQuests_AfterClaim_CurrentAdvances(t *testing.T) {
	setupQuestTestDB(t)
	playerID := "00000000-0000-0000-0000-000000001501"
	defer cleanupQuestFullPlayer(t, playerID)
	seedQuestFullPlayer(t, playerID)

	pqID := markQuestCompletedFull(t, playerID, "main_01_collecting_resources")
	if code, _ := invokeClaimQuest(t, playerID, pqID); code != http.StatusOK {
		t.Fatalf("first claim status=%d", code)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/quests", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	ListQuests(rr, req)

	var body struct {
		CurrentMainQuest *struct {
			QuestKey string `json:"quest_key"`
		} `json:"current_main_quest"`
	}
	json.Unmarshal(rr.Body.Bytes(), &body)
	if body.CurrentMainQuest == nil {
		t.Fatalf("current_main_quest is nil after claim")
	}
	if body.CurrentMainQuest.QuestKey == "main_01_collecting_resources" {
		t.Errorf("current_main_quest still %q after claim — should have advanced",
			body.CurrentMainQuest.QuestKey)
	}
}

// TestGetDailyQuests_Shape verifies the daily endpoint returns the expected
// structure: date, daily_points, quests, tier_rewards.
func TestGetDailyQuests_Shape(t *testing.T) {
	setupQuestTestDB(t)
	playerID := "00000000-0000-0000-0000-000000001502"
	defer cleanupQuestFullPlayer(t, playerID)
	seedQuestFullPlayer(t, playerID)

	req := httptest.NewRequest(http.MethodGet, "/api/quests/daily", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	GetDailyQuests(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("GetDailyQuests: status=%d body=%s", rr.Code, rr.Body.String())
	}

	var body struct {
		Date        string `json:"date"`
		DailyPoints int    `json:"daily_points"`
		Quests      []struct {
			QuestKey    string `json:"quest_key"`
			DisplayName string `json:"display_name"`
			Points      int    `json:"points"`
		} `json:"quests"`
		TierRewards []struct {
			Tier           string `json:"tier"`
			PointsRequired int    `json:"points_required"`
			Claimed        bool   `json:"claimed"`
		} `json:"tier_rewards"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v body=%s", err, rr.Body.String())
	}
	if body.Date == "" {
		t.Errorf("expected date in response")
	}
	if body.DailyPoints != 0 {
		t.Errorf("expected initial daily_points=0, got %d", body.DailyPoints)
	}
	if len(body.TierRewards) != 4 {
		t.Errorf("expected 4 tier rewards (bronze/silver/gold/diamond), got %d", len(body.TierRewards))
	}
	// Verify 4 tier rewards in the expected order.
	wantTiers := []string{"bronze", "silver", "gold", "diamond"}
	for i, tr := range body.TierRewards {
		if i >= len(wantTiers) {
			break
		}
		if tr.Tier != wantTiers[i] {
			t.Errorf("tier_rewards[%d]=%s, want %s", i, tr.Tier, wantTiers[i])
		}
		if tr.Claimed {
			t.Errorf("tier %s should start unclaimed", tr.Tier)
		}
	}
	// All daily quest entries should have non-empty quest_key + display_name (if any seeded).
	if len(body.Quests) > 0 {
		for _, q := range body.Quests {
			if q.QuestKey == "" || q.DisplayName == "" {
				t.Errorf("daily quest entry missing key/name: %+v", q)
			}
		}
	} else {
		// If no daily quests are seeded as 'active', that's still OK — sanity skip.
		t.Logf("no active daily quests seeded; skipping per-quest assertions")
	}
}

// =============================================================================
// helpers
// =============================================================================

// newReader wraps a byte slice as an io.Reader for httptest body construction.
func newReader(b []byte) io.Reader {
	return bytes.NewReader(b)
}
