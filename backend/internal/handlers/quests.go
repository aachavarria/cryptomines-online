package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
	"github.com/cryptomines-online/backend/internal/services"
)

// ListQuests handles GET /api/quests
// Returns all player quests (main + side) with current status and progress.
func ListQuests(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	// Ensure player has quest rows initialized
	ensurePlayerQuests(playerID)

	rows, err := database.DB.Query(
		`SELECT pq.id, qt.quest_key, qt.display_name, qt.description, qt.category,
		        pq.status, pq.progress_value, qt.requirement_value, qt.chain_order,
		        qt.reward_metal, qt.reward_he3, qt.reward_gold, qt.reward_item_json
		 FROM player_quests pq
		 JOIN quest_types qt ON pq.quest_type_id = qt.id
		 WHERE pq.player_id = $1 AND qt.category IN ('main', 'side') AND qt.is_active = true
		 ORDER BY qt.category, qt.chain_order`, playerID,
	)
	if err != nil {
		log.Printf("Failed to list quests: %v", err)
		http.Error(w, `{"error":"failed to list quests"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var mainQuests []models.PlayerQuestWithType
	var sideQuests []models.PlayerQuestWithType
	for rows.Next() {
		var q models.PlayerQuestWithType
		err := rows.Scan(
			&q.ID, &q.QuestKey, &q.DisplayName, &q.Description, &q.Category,
			&q.Status, &q.ProgressValue, &q.RequirementValue, &q.ChainOrder,
			&q.RewardMetal, &q.RewardHe3, &q.RewardGold, &q.RewardItemJSON,
		)
		if err != nil {
			log.Printf("Failed to scan quest: %v", err)
			continue
		}
		if q.Category == "main" {
			mainQuests = append(mainQuests, q)
		} else {
			sideQuests = append(sideQuests, q)
		}
	}

	// Find current main quest (first non-claimed active quest)
	var currentMainQuest *models.PlayerQuestWithType
	for i := range mainQuests {
		if mainQuests[i].Status != "claimed" {
			currentMainQuest = &mainQuests[i]
			break
		}
	}

	type response struct {
		MainQuests       []models.PlayerQuestWithType `json:"main_quests"`
		SideQuests       []models.PlayerQuestWithType `json:"side_quests"`
		CurrentMainQuest *models.PlayerQuestWithType  `json:"current_main_quest"`
	}

	if mainQuests == nil {
		mainQuests = []models.PlayerQuestWithType{}
	}
	if sideQuests == nil {
		sideQuests = []models.PlayerQuestWithType{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{
		MainQuests:       mainQuests,
		SideQuests:       sideQuests,
		CurrentMainQuest: currentMainQuest,
	})
}

// GetDailyQuests handles GET /api/quests/daily
// Returns today's daily quest progress and tier reward status.
func GetDailyQuests(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	// Ensure today's daily progress row exists
	ensureDailyProgress(playerID)

	var dp models.DailyQuestProgress
	err := database.DB.QueryRow(
		`SELECT id, player_id, quest_date, daily_points, quests_completed_json,
		        tier_rewards_claimed_json, created_at, updated_at
		 FROM daily_quest_progress
		 WHERE player_id = $1 AND quest_date = CURRENT_DATE`, playerID,
	).Scan(
		&dp.ID, &dp.PlayerID, &dp.QuestDate, &dp.DailyPoints,
		&dp.QuestsCompletedJSON, &dp.TierRewardsClaimedJSON,
		&dp.CreatedAt, &dp.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to get daily progress: %v", err)
		http.Error(w, `{"error":"failed to get daily progress"}`, http.StatusInternalServerError)
		return
	}

	// Fetch daily quest type definitions
	rows, err := database.DB.Query(
		`SELECT quest_key, display_name, requirement_type, requirement_value
		 FROM quest_types
		 WHERE category = 'daily' AND is_active = true
		 ORDER BY quest_key`,
	)
	if err != nil {
		log.Printf("Failed to list daily quest types: %v", err)
		http.Error(w, `{"error":"failed to list daily quests"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Parse completed quests map
	var completedMap map[string]interface{}
	json.Unmarshal(dp.QuestsCompletedJSON, &completedMap)
	if completedMap == nil {
		completedMap = map[string]interface{}{}
	}

	// Parse claimed tiers
	var claimedTiers []string
	json.Unmarshal(dp.TierRewardsClaimedJSON, &claimedTiers)

	claimedSet := map[string]bool{}
	for _, t := range claimedTiers {
		claimedSet[t] = true
	}

	// Daily quest point values (from GO2 research)
	pointMap := map[string]int{
		"daily_login":                10,
		"daily_collect_dues":         4,
		"daily_need_for_speed":       3,
		"daily_stockpiling":          1, // per harvest, max 3
		"daily_donations":            6,
		"daily_restricted_instances": 5, // per attempt, max 10
	}

	type dailyQuestEntry struct {
		QuestKey    string `json:"quest_key"`
		DisplayName string `json:"display_name"`
		Completed   bool   `json:"completed"`
		Points      int    `json:"points"`
		Progress    *int   `json:"progress,omitempty"`
		Required    *int   `json:"required,omitempty"`
		PointsPer   *int   `json:"points_per,omitempty"`
		MaxPoints   *int   `json:"max_points,omitempty"`
	}

	var quests []dailyQuestEntry
	for rows.Next() {
		var questKey, displayName, reqType string
		var reqValue int
		if err := rows.Scan(&questKey, &displayName, &reqType, &reqValue); err != nil {
			continue
		}

		pts := pointMap[questKey]
		entry := dailyQuestEntry{
			QuestKey:    questKey,
			DisplayName: displayName,
			Points:      pts,
		}

		// Check completion status
		if val, ok := completedMap[questKey]; ok {
			switch v := val.(type) {
			case bool:
				entry.Completed = v
			case float64:
				// Multi-step quest (e.g., stockpiling)
				progress := int(v)
				entry.Progress = &progress
				entry.Required = &reqValue
				entry.PointsPer = &pts
				maxPts := pts * reqValue
				entry.MaxPoints = &maxPts
				entry.Completed = progress >= reqValue
			}
		}

		quests = append(quests, entry)
	}

	type tierReward struct {
		Tier           string `json:"tier"`
		PointsRequired int    `json:"points_required"`
		Claimed        bool   `json:"claimed"`
	}

	tierRewards := []tierReward{
		{Tier: "bronze", PointsRequired: 10, Claimed: claimedSet["bronze"]},
		{Tier: "silver", PointsRequired: 30, Claimed: claimedSet["silver"]},
		{Tier: "gold", PointsRequired: 50, Claimed: claimedSet["gold"]},
		{Tier: "diamond", PointsRequired: 70, Claimed: claimedSet["diamond"]},
	}

	type response struct {
		Date        string            `json:"date"`
		DailyPoints int               `json:"daily_points"`
		Quests      []dailyQuestEntry `json:"quests"`
		TierRewards []tierReward      `json:"tier_rewards"`
	}

	if quests == nil {
		quests = []dailyQuestEntry{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{
		Date:        dp.QuestDate,
		DailyPoints: dp.DailyPoints,
		Quests:      quests,
		TierRewards: tierRewards,
	})
}

// ClaimQuest handles POST /api/quests/{id}/claim
// Claims reward for a completed quest and unlocks the next quest in chain.
func ClaimQuest(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	questID := r.PathValue("id")

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Get the player quest and its type info
	var pq models.PlayerQuest
	var qt models.QuestType
	err = tx.QueryRow(
		`SELECT pq.id, pq.player_id, pq.quest_type_id, pq.status, pq.progress_value,
		        qt.id, qt.quest_key, qt.category, qt.display_name, qt.requirement_value,
		        qt.reward_metal, qt.reward_he3, qt.reward_gold, qt.reward_item_json
		 FROM player_quests pq
		 JOIN quest_types qt ON pq.quest_type_id = qt.id
		 WHERE pq.id = $1 AND pq.player_id = $2
		 FOR UPDATE OF pq`, questID, playerID,
	).Scan(
		&pq.ID, &pq.PlayerID, &pq.QuestTypeID, &pq.Status, &pq.ProgressValue,
		&qt.ID, &qt.QuestKey, &qt.Category, &qt.DisplayName, &qt.RequirementValue,
		&qt.RewardMetal, &qt.RewardHe3, &qt.RewardGold, &qt.RewardItemJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"quest not found"}`, http.StatusNotFound)
			return
		}
		log.Printf("Failed to get quest: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Verify quest is in 'completed' status
	if pq.Status != "completed" {
		http.Error(w, `{"error":"quest is not completed"}`, http.StatusConflict)
		return
	}

	// Mark quest as claimed
	now := time.Now()
	_, err = tx.Exec(
		`UPDATE player_quests SET status = 'claimed', claimed_at = $1, updated_at = $1
		 WHERE id = $2`, now, questID,
	)
	if err != nil {
		log.Printf("Failed to claim quest: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Award resource rewards to player's homeworld planet
	if qt.RewardMetal > 0 || qt.RewardHe3 > 0 || qt.RewardGold > 0 {
		_, err = tx.Exec(
			`UPDATE resources SET
			    metal = metal + $1, he3 = he3 + $2, gold = gold + $3, updated_at = now()
			 WHERE planet_id = (
			     SELECT id FROM planets WHERE player_id = $4 AND is_homeworld = true LIMIT 1
			 )`, qt.RewardMetal, qt.RewardHe3, qt.RewardGold, playerID,
		)
		if err != nil {
			log.Printf("Failed to award resources: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
	}

	// Process item rewards (add to inventory instead of direct unlock)
	if len(qt.RewardItemJSON) > 0 && string(qt.RewardItemJSON) != "[]" {
		var items []map[string]interface{}
		if err := json.Unmarshal(qt.RewardItemJSON, &items); err != nil {
			log.Printf("Failed to parse reward items: %v", err)
		} else {
			for _, item := range items {
				itemType, hasType := item["type"].(string)
				if !hasType {
					continue
				}

				if itemType == "item" {
					// Regular item (resource pack, boost, battle item)
					itemKey, _ := item["item_key"].(string)
					if itemKey == "" {
						continue
					}
					quantity := 1
					if qty, ok := item["quantity"].(float64); ok {
						quantity = int(qty)
					}

					// Verify item_key exists in item_types before inserting
					var exists bool
					err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM item_types WHERE item_key = $1)`, itemKey).Scan(&exists)
					if err != nil || !exists {
						log.Printf("Item key %s does not exist in item_types, skipping", itemKey)
						continue
					}

					_, err = tx.Exec(`
						INSERT INTO player_inventory (player_id, item_key, quantity)
						VALUES ($1, $2, $3)
						ON CONFLICT (player_id, item_key)
						DO UPDATE SET quantity = player_inventory.quantity + EXCLUDED.quantity, updated_at = now()
					`, playerID, itemKey, quantity)
					if err != nil {
						log.Printf("Failed to add item to inventory: %v", err)
					}

				} else if itemType == "blueprint" {
					// Blueprint item (add to inventory, NOT direct unlock)
					blueprintKey, _ := item["blueprint_key"].(string)
					if blueprintKey == "" {
						continue
					}

					// Get blueprint ID (try hull first, then module)
					var blueprintID int
					err = tx.QueryRow(`
						SELECT id FROM blueprints
						WHERE blueprint_type = 'hull' AND hull_type_id = (
							SELECT id FROM hull_types WHERE name = $1
						)
						LIMIT 1
					`, blueprintKey).Scan(&blueprintID)

					// If not found as hull, try module
					if err == sql.ErrNoRows {
						err = tx.QueryRow(`
							SELECT id FROM blueprints
							WHERE blueprint_type = 'module' AND module_type_id = (
								SELECT id FROM module_types WHERE name = $1 AND tier = 0
							)
							LIMIT 1
						`, blueprintKey).Scan(&blueprintID)
					}

					if err != nil {
						log.Printf("Failed to find blueprint %s: %v", blueprintKey, err)
						continue
					}

					// Create dynamic blueprint item in item_types (if not exists)
					blueprintItemKey := "blueprint_" + blueprintKey
					_, err = tx.Exec(`
						INSERT INTO item_types (item_key, display_name, category, description, blueprint_id)
						VALUES ($1, $2, 'blueprint', $3, $4)
						ON CONFLICT (item_key) DO NOTHING
					`, blueprintItemKey, "Blueprint: "+blueprintKey, "Unlock "+blueprintKey+" blueprint", blueprintID)
					if err != nil {
						log.Printf("Failed to create blueprint item type: %v", err)
						continue
					}

					// Add to inventory (item_type now guaranteed to exist)
					_, err = tx.Exec(`
						INSERT INTO player_inventory (player_id, item_key, quantity)
						VALUES ($1, $2, 1)
						ON CONFLICT (player_id, item_key)
						DO UPDATE SET quantity = player_inventory.quantity + 1, updated_at = now()
					`, playerID, blueprintItemKey)
					if err != nil {
						log.Printf("Failed to add blueprint item to inventory: %v", err)
					}

				} else if itemType == "commander" {
					// Commander card item (add to inventory, NOT direct unlock)
					commanderKey, _ := item["commander_key"].(string)
					if commanderKey == "" {
						continue
					}

					// Create dynamic commander item in item_types (if not exists)
					commanderItemKey := "commander_" + commanderKey
					_, err = tx.Exec(`
						INSERT INTO item_types (item_key, display_name, category, description, commander_type)
						VALUES ($1, $2, 'commander', $3, $4)
						ON CONFLICT (item_key) DO NOTHING
					`, commanderItemKey, "Commander: "+commanderKey, "Unlock "+commanderKey+" commander", commanderKey)
					if err != nil {
						log.Printf("Failed to create commander item type: %v", err)
						continue
					}

					// Add to inventory (item_type now guaranteed to exist)
					_, err = tx.Exec(`
						INSERT INTO player_inventory (player_id, item_key, quantity)
						VALUES ($1, $2, 1)
						ON CONFLICT (player_id, item_key)
						DO UPDATE SET quantity = player_inventory.quantity + 1, updated_at = now()
					`, playerID, commanderItemKey)
					if err != nil {
						log.Printf("Failed to add commander item to inventory: %v", err)
					}
				}
			}
		}
	}

	// Unlock next quest in chain (for main and side quests)
	var nextQuestKey *string
	err = tx.QueryRow(
		`UPDATE player_quests SET status = 'available', updated_at = now()
		 WHERE player_id = $1 AND quest_type_id = (
		     SELECT id FROM quest_types WHERE prerequisite_quest_id = $2 AND is_active = true LIMIT 1
		 ) AND status = 'locked'
		 RETURNING (SELECT quest_key FROM quest_types WHERE id = quest_type_id)`,
		playerID, qt.ID,
	).Scan(&nextQuestKey)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Failed to unlock next quest: %v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	type rewardResponse struct {
		Metal int64           `json:"metal"`
		He3   int64           `json:"he3"`
		Gold  int64           `json:"gold"`
		Items json.RawMessage `json:"items"`
	}

	type response struct {
		QuestKey          string         `json:"quest_key"`
		Rewards           rewardResponse `json:"rewards"`
		NextQuestUnlocked *string        `json:"next_quest_unlocked"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{
		QuestKey: qt.QuestKey,
		Rewards: rewardResponse{
			Metal: qt.RewardMetal,
			He3:   qt.RewardHe3,
			Gold:  qt.RewardGold,
			Items: qt.RewardItemJSON,
		},
		NextQuestUnlocked: nextQuestKey,
	})
}

// ClaimDailyTier handles POST /api/quests/daily/claim-tier
// Claims a daily tier reward (bronze/silver/gold/diamond).
func ClaimDailyTier(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	var req struct {
		Tier string `json:"tier"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate tier
	tierPoints := map[string]int{
		"bronze":  10,
		"silver":  30,
		"gold":    50,
		"diamond": 70,
	}
	requiredPoints, validTier := tierPoints[req.Tier]
	if !validTier {
		http.Error(w, `{"error":"invalid tier"}`, http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Get today's daily progress with lock
	var dp models.DailyQuestProgress
	err = tx.QueryRow(
		`SELECT id, daily_points, tier_rewards_claimed_json
		 FROM daily_quest_progress
		 WHERE player_id = $1 AND quest_date = CURRENT_DATE
		 FOR UPDATE`, playerID,
	).Scan(&dp.ID, &dp.DailyPoints, &dp.TierRewardsClaimedJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"no daily progress for today"}`, http.StatusNotFound)
			return
		}
		log.Printf("Failed to get daily progress: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Check sufficient points
	if dp.DailyPoints < requiredPoints {
		http.Error(w, `{"error":"insufficient daily points"}`, http.StatusConflict)
		return
	}

	// Check not already claimed
	var claimedTiers []string
	json.Unmarshal(dp.TierRewardsClaimedJSON, &claimedTiers)
	for _, t := range claimedTiers {
		if t == req.Tier {
			http.Error(w, `{"error":"tier already claimed"}`, http.StatusConflict)
			return
		}
	}

	// Add tier to claimed list
	claimedTiers = append(claimedTiers, req.Tier)
	claimedJSON, _ := json.Marshal(claimedTiers)

	_, err = tx.Exec(
		`UPDATE daily_quest_progress
		 SET tier_rewards_claimed_json = $1, updated_at = now()
		 WHERE id = $2`, claimedJSON, dp.ID,
	)
	if err != nil {
		log.Printf("Failed to update tier claims: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Tier reward mapping (simplified; from GO2 research)
	type tierRewardItem struct {
		Type     string `json:"type"`
		Quantity int    `json:"quantity"`
	}
	tierRewards := map[string]tierRewardItem{
		"bronze":  {Type: "loudspeaker", Quantity: 1},
		"silver":  {Type: "resource_box", Quantity: 1},
		"gold":    {Type: "sp_card", Quantity: 1},
		"diamond": {Type: "raw_gemstone", Quantity: 1},
	}

	type response struct {
		Tier   string         `json:"tier"`
		Reward tierRewardItem `json:"reward"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{
		Tier:   req.Tier,
		Reward: tierRewards[req.Tier],
	})
}

// ensurePlayerQuests initializes quest rows for a player if they don't exist.
// The first main quest starts as 'available', all others as 'locked'.
// Side quests tier 1 starts as 'available', higher tiers as 'locked'.
func ensurePlayerQuests(playerID string) {
	// Check if player already has quest rows
	var count int
	err := database.DB.QueryRow(
		`SELECT COUNT(*) FROM player_quests WHERE player_id = $1`, playerID,
	).Scan(&count)
	if err != nil || count > 0 {
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction for quest init: %v", err)
		return
	}
	defer tx.Rollback()

	// Re-check inside transaction to avoid race conditions
	err = tx.QueryRow(
		`SELECT COUNT(*) FROM player_quests WHERE player_id = $1`, playerID,
	).Scan(&count)
	if err != nil || count > 0 {
		return
	}

	// Insert all active quest types as player quests
	// Main quests: first in chain = 'available', rest = 'locked'
	// Side quests: tier 1 (chain_order=1, no prerequisite) = 'available', rest = 'locked'
	_, err = tx.Exec(
		`INSERT INTO player_quests (player_id, quest_type_id, status)
		 SELECT $1, qt.id,
		     CASE
		         WHEN qt.category = 'main' AND qt.chain_order = 1 THEN 'available'
		         WHEN qt.category = 'side' AND qt.prerequisite_quest_id IS NULL THEN 'available'
		         ELSE 'locked'
		     END
		 FROM quest_types qt
		 WHERE qt.category IN ('main', 'side') AND qt.is_active = true`, playerID,
	)
	if err != nil {
		log.Printf("Failed to initialize player quests: %v", err)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit quest init: %v", err)
	}
}

// ensureDailyProgress creates today's daily progress row if it doesn't exist.
func ensureDailyProgress(playerID string) {
	_, err := database.DB.Exec(
		`INSERT INTO daily_quest_progress (player_id, quest_date)
		 VALUES ($1, CURRENT_DATE)
		 ON CONFLICT (player_id, quest_date) DO NOTHING`, playerID,
	)
	if err != nil {
		log.Printf("Failed to ensure daily progress: %v", err)
	}
}

// SyncQuests handles POST /api/quests/sync
// Syncs quest progress with existing buildings/research
func SyncQuests(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	err := services.SyncBuildingQuests(playerID)
	if err != nil {
		log.Printf("Failed to sync quests for player %s: %v", playerID, err)
		http.Error(w, `{"error":"failed to sync quests"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Quests synced with existing buildings",
	})
}
