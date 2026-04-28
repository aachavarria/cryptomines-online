package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
)

// setupQuestTestDB initializes the test DB connection.
func setupQuestTestDB(t *testing.T) {
	if database.DB == nil {
		if err := database.InitSupabase(); err != nil {
			t.Fatalf("Failed to init database: %v", err)
		}
	}
}

func cleanupQuestTestPlayer(t *testing.T, playerID string) {
	t.Helper()
	database.DB.Exec("DELETE FROM player_inventory WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM player_quests WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM resources WHERE planet_id IN (SELECT id FROM planets WHERE player_id = $1)", playerID)
	database.DB.Exec("DELETE FROM planets WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM players WHERE id = $1", playerID)
}

func seedQuestPlayer(t *testing.T, playerID string) {
	t.Helper()
	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'questtest', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)
	if err != nil {
		t.Fatalf("create player: %v", err)
	}

	planetID := "00000000-0000-0000-aaaa-" + playerID[len(playerID)-12:]
	_, err = database.DB.Exec(`
		INSERT INTO planets (id, player_id, name, is_homeworld, position_x, position_y, created_at, updated_at)
		VALUES ($1, $2, 'Test', true, 100, 100, now(), now())
		ON CONFLICT (id) DO NOTHING
	`, planetID, playerID)
	if err != nil {
		t.Fatalf("create planet: %v", err)
	}

	_, err = database.DB.Exec(`
		INSERT INTO resources (planet_id, metal, he3, gold)
		VALUES ($1, 0, 0, 0)
	`, planetID)
	if err != nil {
		t.Fatalf("create resources: %v", err)
	}

	ensurePlayerQuests(playerID)
}

// markQuestCompleted forces a player_quest into 'completed' so it can be claimed.
func markQuestCompleted(t *testing.T, playerID, questKey string) string {
	t.Helper()
	var pqID string
	err := database.DB.QueryRow(`
		UPDATE player_quests pq
		SET status = 'completed',
		    progress_value = qt.requirement_value,
		    completed_at = now(),
		    updated_at = now()
		FROM quest_types qt
		WHERE pq.quest_type_id = qt.id AND pq.player_id = $1 AND qt.quest_key = $2
		RETURNING pq.id
	`, playerID, questKey).Scan(&pqID)
	if err != nil {
		t.Fatalf("mark quest %s completed: %v", questKey, err)
	}
	return pqID
}

func claimQuest(t *testing.T, playerID, pqID string) (int, map[string]any) {
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

// TestClaimQuest_BlueprintReward_Estrella reproduces the original bug
// "Failed to find blueprint estrella: sql: no rows in result set" and verifies
// the blueprint item is granted via player_inventory after the migration.
func TestClaimQuest_BlueprintReward_Estrella(t *testing.T) {
	setupQuestTestDB(t)
	playerID := "00000000-0000-0000-0000-000000000901"
	defer cleanupQuestTestPlayer(t, playerID)
	seedQuestPlayer(t, playerID)

	pqID := markQuestCompleted(t, playerID, "main_07_he3_production")
	code, _ := claimQuest(t, playerID, pqID)
	if code != http.StatusOK {
		t.Fatalf("claim main_07: status=%d", code)
	}

	var qty int
	err := database.DB.QueryRow(`
		SELECT quantity FROM player_inventory
		WHERE player_id = $1 AND item_key = 'blueprint_estrella'
	`, playerID).Scan(&qty)
	if err != nil {
		t.Fatalf("expected blueprint_estrella in inventory, got: %v", err)
	}
	if qty != 1 {
		t.Errorf("blueprint_estrella quantity: got %d, want 1", qty)
	}

	var status string
	if err := database.DB.QueryRow(`
		SELECT status FROM player_quests WHERE id = $1
	`, pqID).Scan(&status); err != nil {
		t.Fatalf("read claimed status: %v", err)
	}
	if status != "claimed" {
		t.Errorf("quest status: got %q, want claimed", status)
	}
}

// TestClaimQuest_AllBlueprintRewards verifies every quest with a blueprint
// reward references a blueprint_key that resolves and gets inserted into
// player_inventory. Catches the entire class of bug, not just estrella.
func TestClaimQuest_AllBlueprintRewards(t *testing.T) {
	setupQuestTestDB(t)
	playerID := "00000000-0000-0000-0000-000000000902"
	defer cleanupQuestTestPlayer(t, playerID)
	seedQuestPlayer(t, playerID)

	rows, err := database.DB.Query(`
		SELECT qt.quest_key,
		       (jsonb_array_elements(qt.reward_item_json::jsonb) ->> 'blueprint_key') AS bp_key
		FROM quest_types qt
		WHERE qt.is_active = true
		  AND qt.reward_item_json::jsonb @> '[{"type":"blueprint"}]'
	`)
	if err != nil {
		t.Fatalf("query blueprint rewards: %v", err)
	}
	defer rows.Close()

	type pair struct{ Quest, Key string }
	var entries []pair
	for rows.Next() {
		var p pair
		if err := rows.Scan(&p.Quest, &p.Key); err != nil {
			t.Fatalf("scan: %v", err)
		}
		if p.Key != "" {
			entries = append(entries, p)
		}
	}
	if len(entries) == 0 {
		t.Fatal("no quests with blueprint rewards found")
	}

	for _, e := range entries {
		t.Run(e.Quest+"->"+e.Key, func(t *testing.T) {
			pqID := markQuestCompleted(t, playerID, e.Quest)
			code, _ := claimQuest(t, playerID, pqID)
			if code != http.StatusOK {
				t.Fatalf("claim %s: status=%d", e.Quest, code)
			}
			var qty int
			err := database.DB.QueryRow(`
				SELECT quantity FROM player_inventory
				WHERE player_id = $1 AND item_key = $2
			`, playerID, "blueprint_"+e.Key).Scan(&qty)
			if err != nil {
				t.Fatalf("blueprint_%s missing from inventory: %v", e.Key, err)
			}
			if qty < 1 {
				t.Errorf("blueprint_%s qty: got %d, want >=1", e.Key, qty)
			}
		})
	}
}

// TestClaimQuest_IdempotentBlueprintLookup verifies that every blueprint
// referenced by an active quest can be resolved by blueprint_key.
// This is a fast pure-DB invariant — runs in ms, no HTTP.
func TestClaimQuest_BlueprintKeysAreSeeded(t *testing.T) {
	setupQuestTestDB(t)

	rows, err := database.DB.Query(`
		SELECT DISTINCT (jsonb_array_elements(qt.reward_item_json::jsonb) ->> 'blueprint_key') AS bp_key
		FROM quest_types qt
		WHERE qt.is_active = true
		  AND qt.reward_item_json::jsonb @> '[{"type":"blueprint"}]'
	`)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	defer rows.Close()

	var missing []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			continue
		}
		if key == "" {
			continue
		}
		var id int
		err := database.DB.QueryRow(
			`SELECT id FROM blueprints WHERE blueprint_key = $1`, key,
		).Scan(&id)
		if err != nil {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		t.Fatalf("blueprint_keys referenced by quests but missing in blueprints table: %v", missing)
	}
}
