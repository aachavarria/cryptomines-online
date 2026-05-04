package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/errs"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
	"github.com/cryptomines-online/backend/internal/services"
)

const (
	RecruitmentCostGold = 10000 // Cost per recruitment in Gold
	MaxCommanders       = 60    // Maximum commanders per player
)

// Commander represents a player's owned commander
type Commander struct {
	ID               string    `json:"id"`
	PlayerID         string    `json:"player_id"`
	Name             string    `json:"name"`
	Rarity           string    `json:"rarity"`
	StarRank         int       `json:"star_rank"`
	Accuracy         int       `json:"accuracy"`
	Dodge            int       `json:"dodge"`
	Speed            int       `json:"speed"`
	Electron         int       `json:"electron"`
	WeaponExpertise  *string   `json:"weapon_expertise"`
	ShipExpertise    *string   `json:"ship_expertise"`
	IsDeployed       bool      `json:"is_deployed"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// RecruitCommander handles POST /api/commanders/recruit
// Gacha system: Roll for rarity (50% common, 35% skill, 15% super), pick random commander
// Cost: 10,000 Gold per recruitment
// Cooldown: Based on Command Center level (3h at Lv1, 1h10m at Lv12)
// Duplicates: Create commander card item in inventory (for merging later)
func RecruitCommander(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Get player's planet and resources
	var planetID string
	var resources models.Resource
	err = tx.QueryRow(`
		SELECT p.id, r.id, r.planet_id, r.metal, r.he3, r.gold,
		       r.metal_per_hour, r.he3_per_hour, r.gold_per_hour, r.storage_capacity,
		       r.warehouse_metal, r.warehouse_he3, r.warehouse_gold,
		       r.last_warehouse_update, r.last_collected_at, r.updated_at
		FROM planets p
		JOIN resources r ON r.planet_id = p.id
		WHERE p.player_id = $1 AND p.is_homeworld = true
		FOR UPDATE OF r
	`, playerID).Scan(
		&planetID,
		&resources.ID, &resources.PlanetID, &resources.Metal, &resources.He3, &resources.Gold,
		&resources.MetalPerHour, &resources.He3PerHour, &resources.GoldPerHour, &resources.StorageCapacity,
		&resources.WarehouseMetal, &resources.WarehouseHe3, &resources.WarehouseGold,
		&resources.LastWarehouseUpdate, &resources.LastCollectedAt, &resources.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to get resources: %v", err)
		http.Error(w, `{"error":"planet not found"}`, http.StatusNotFound)
		return
	}

	// Promote any buildings whose construction has finished so a freshly-built
	// Command Center counts toward the level check below.
	applyCompletedUpgrades(planetID)

	// Check Gold cost
	if resources.Gold < RecruitmentCostGold {
		errs.InsufficientResources(
			fmt.Sprintf("Not enough Gold to recruit commander (need %d, have %d)", RecruitmentCostGold, resources.Gold),
			map[string]int64{"gold": RecruitmentCostGold},
			map[string]int64{"gold": resources.Gold},
		).WriteJSON(w, http.StatusPaymentRequired)
		return
	}

	// Check commander count limit
	var commanderCount int
	err = tx.QueryRow(`SELECT COUNT(*) FROM commanders WHERE player_id = $1`, playerID).Scan(&commanderCount)
	if err != nil {
		log.Printf("Failed to count commanders: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if commanderCount >= MaxCommanders {
		http.Error(w, fmt.Sprintf(`{"error":"max commanders reached (%d/%d)"}`, commanderCount, MaxCommanders), http.StatusConflict)
		return
	}

	// Check Command Center cooldown based on level. Building must be fully
	// constructed (level >= 1, not still upgrading from 0).
	var ccLevel int
	err = tx.QueryRow(`
		SELECT b.level FROM buildings b
		JOIN building_types bt ON b.building_type = bt.id
		WHERE b.planet_id = $1 AND bt.name = 'command_center'
		ORDER BY b.level DESC
		LIMIT 1
	`, planetID).Scan(&ccLevel)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"command center not built"}`, http.StatusConflict)
		return
	}
	if err != nil {
		log.Printf("Failed to query command center level: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if ccLevel < 1 {
		http.Error(w, `{"error":"command center still under construction"}`, http.StatusConflict)
		return
	}

	// Cooldown per GO2: Lv1=3h, Lv2=2h50m, ... Lv5=2h20m, then continues at 10min/level, min 1h10m at Lv12
	cooldownMinutes := 180 - (ccLevel-1)*10
	if cooldownMinutes < 70 {
		cooldownMinutes = 70 // minimum 1h10m at max Command Center level (Lv12)
	}
	cooldownDuration := time.Duration(cooldownMinutes) * time.Minute

	var lastRecruitAt sql.NullTime
	err = tx.QueryRow(`
		SELECT MAX(created_at) FROM commanders WHERE player_id = $1
	`, playerID).Scan(&lastRecruitAt)
	if err == nil && lastRecruitAt.Valid {
		if time.Since(lastRecruitAt.Time) < cooldownDuration {
			remaining := cooldownDuration - time.Since(lastRecruitAt.Time)
			http.Error(w, fmt.Sprintf(`{"error":"recruitment on cooldown","remaining_seconds":%d}`, int(remaining.Seconds())), http.StatusConflict)
			return
		}
	}

	// Gacha roll: Determine rarity
	rarity := rollCommanderRarity()

	// Pick random commander of chosen rarity
	var commanderType struct {
		ID          int
		Name        string
		DisplayName string
		Rarity      string
		Accuracy    int
		Dodge       int
		Speed       int
		Electron    int
	}
	err = tx.QueryRow(`
		SELECT id, name, display_name, rarity, base_accuracy, base_dodge, base_speed, base_electron
		FROM commander_types
		WHERE rarity = $1
		ORDER BY RANDOM()
		LIMIT 1
	`, rarity).Scan(
		&commanderType.ID, &commanderType.Name, &commanderType.DisplayName, &commanderType.Rarity,
		&commanderType.Accuracy, &commanderType.Dodge, &commanderType.Speed, &commanderType.Electron,
	)
	if err != nil {
		log.Printf("Failed to get commander type: %v", err)
		http.Error(w, `{"error":"failed to recruit commander"}`, http.StatusInternalServerError)
		return
	}

	// Check if player already owns this commander
	var existingID string
	err = tx.QueryRow(`
		SELECT id FROM commanders
		WHERE player_id = $1 AND name = $2
	`, playerID, commanderType.Name).Scan(&existingID)

	var recruitedCommander Commander
	var isDuplicate bool

	if err == sql.ErrNoRows {
		// New commander - insert into commanders table. New recruits start at
		// star rank 1; duplicates merged later promote rank up to 15.
		err = tx.QueryRow(`
			INSERT INTO commanders (player_id, name, rarity, star_rank, accuracy, dodge, speed, electron, is_deployed)
			VALUES ($1, $2, $3, 1, $4, $5, $6, $7, false)
			RETURNING id, player_id, name, rarity, star_rank, accuracy, dodge, speed, electron, is_deployed, created_at, updated_at
		`, playerID, commanderType.Name, commanderType.Rarity,
			commanderType.Accuracy, commanderType.Dodge, commanderType.Speed, commanderType.Electron).Scan(
			&recruitedCommander.ID, &recruitedCommander.PlayerID, &recruitedCommander.Name,
			&recruitedCommander.Rarity, &recruitedCommander.StarRank,
			&recruitedCommander.Accuracy, &recruitedCommander.Dodge, &recruitedCommander.Speed, &recruitedCommander.Electron,
			&recruitedCommander.IsDeployed, &recruitedCommander.CreatedAt, &recruitedCommander.UpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to create commander: %v", err)
			http.Error(w, `{"error":"failed to create commander"}`, http.StatusInternalServerError)
			return
		}
		isDuplicate = false

	} else if err != nil {
		log.Printf("Failed to check existing commander: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return

	} else {
		// Duplicate - give commander card item to inventory
		commanderItemKey := "commander_" + commanderType.Name
		isDuplicate = true

		// Create dynamic commander item in item_types (if not exists)
		_, err = tx.Exec(`
			INSERT INTO item_types (item_key, display_name, category, description, commander_type)
			VALUES ($1, $2, 'commander', $3, $4)
			ON CONFLICT (item_key) DO NOTHING
		`, commanderItemKey, "Commander: "+commanderType.DisplayName,
			"Use to unlock "+commanderType.DisplayName+" for merging", commanderType.Name)
		if err != nil {
			log.Printf("Failed to create commander item type: %v", err)
		}

		// Add to inventory
		_, err = tx.Exec(`
			INSERT INTO player_inventory (player_id, item_key, quantity)
			VALUES ($1, $2, 1)
			ON CONFLICT (player_id, item_key)
			DO UPDATE SET quantity = player_inventory.quantity + 1, updated_at = now()
		`, playerID, commanderItemKey)
		if err != nil {
			log.Printf("Failed to add commander card to inventory: %v", err)
			http.Error(w, `{"error":"failed to add commander card"}`, http.StatusInternalServerError)
			return
		}

		// Return existing commander info (for display purposes)
		err = tx.QueryRow(`
			SELECT id, player_id, name, rarity, star_rank, accuracy, dodge, speed, electron, is_deployed, created_at, updated_at
			FROM commanders
			WHERE id = $1
		`, existingID).Scan(
			&recruitedCommander.ID, &recruitedCommander.PlayerID, &recruitedCommander.Name,
			&recruitedCommander.Rarity, &recruitedCommander.StarRank,
			&recruitedCommander.Accuracy, &recruitedCommander.Dodge, &recruitedCommander.Speed, &recruitedCommander.Electron,
			&recruitedCommander.IsDeployed, &recruitedCommander.CreatedAt, &recruitedCommander.UpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to get existing commander: %v", err)
		}
	}

	// Deduct Gold
	_, err = tx.Exec(`
		UPDATE resources
		SET gold = gold - $1, updated_at = now()
		WHERE planet_id = $2
	`, RecruitmentCostGold, planetID)
	if err != nil {
		log.Printf("Failed to deduct gold: %v", err)
		http.Error(w, `{"error":"failed to deduct gold"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Update quest progress
	services.UpdateQuestProgress(playerID, "recruit_commander", "commander", 1)

	// Response
	type recruitResponse struct {
		Commander   Commander `json:"commander"`
		IsDuplicate bool      `json:"is_duplicate"`
		Message     string    `json:"message"`
	}

	message := fmt.Sprintf("Recruited %s (%s rarity)", commanderType.DisplayName, commanderType.Rarity)
	if isDuplicate {
		message = fmt.Sprintf("Duplicate %s! Commander card added to inventory for merging.", commanderType.DisplayName)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(recruitResponse{
		Commander:   recruitedCommander,
		IsDuplicate: isDuplicate,
		Message:     message,
	})
}

// rollCommanderRarity performs gacha roll for rarity
// Drop rates: Common 50%, Skill 35%, Super 15%
func rollCommanderRarity() string {
	roll := rand.Float64()
	if roll < 0.50 {
		return "common"
	} else if roll < 0.85 { // 0.50 + 0.35
		return "skill"
	}
	return "super" // 0.85 + 0.15
}

// ListCommanders handles GET /api/commanders
// Returns all commanders owned by player
func ListCommanders(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	rows, err := database.DB.Query(`
		SELECT id, player_id, name, rarity, star_rank, accuracy, dodge, speed, electron,
		       weapon_expertise, ship_expertise, is_deployed, created_at, updated_at
		FROM commanders
		WHERE player_id = $1
		ORDER BY star_rank DESC, rarity DESC, name ASC
	`, playerID)
	if err != nil {
		log.Printf("Failed to query commanders: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	commanders := []Commander{}
	for rows.Next() {
		var cmd Commander
		err := rows.Scan(
			&cmd.ID, &cmd.PlayerID, &cmd.Name, &cmd.Rarity, &cmd.StarRank,
			&cmd.Accuracy, &cmd.Dodge, &cmd.Speed, &cmd.Electron,
			&cmd.WeaponExpertise, &cmd.ShipExpertise, &cmd.IsDeployed,
			&cmd.CreatedAt, &cmd.UpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan commander: %v", err)
			continue
		}
		commanders = append(commanders, cmd)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commanders)
}

// MergeCommanders handles POST /api/commanders/merge
// Consumes duplicate commander cards from inventory to increase star_rank
// Max star_rank: 15
func MergeCommanders(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	var req struct {
		CommanderID string `json:"commander_id"`
		Quantity    int    `json:"quantity"` // Number of duplicates to consume
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		http.Error(w, `{"error":"quantity must be positive"}`, http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Get commander with lock
	var cmd Commander
	err = tx.QueryRow(`
		SELECT id, player_id, name, rarity, star_rank, accuracy, dodge, speed, electron, is_deployed
		FROM commanders
		WHERE id = $1 AND player_id = $2
		FOR UPDATE
	`, req.CommanderID, playerID).Scan(
		&cmd.ID, &cmd.PlayerID, &cmd.Name, &cmd.Rarity, &cmd.StarRank,
		&cmd.Accuracy, &cmd.Dodge, &cmd.Speed, &cmd.Electron, &cmd.IsDeployed,
	)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"commander not found"}`, http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Failed to get commander: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Check star rank limit
	if cmd.StarRank >= 15 {
		http.Error(w, `{"error":"commander already at max star rank (15)"}`, http.StatusConflict)
		return
	}

	// Check duplicate cards in inventory
	commanderItemKey := "commander_" + cmd.Name
	var inventoryQty int
	var inventoryID string
	err = tx.QueryRow(`
		SELECT id, quantity
		FROM player_inventory
		WHERE player_id = $1 AND item_key = $2
		FOR UPDATE
	`, playerID, commanderItemKey).Scan(&inventoryID, &inventoryQty)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"no duplicate commander cards available"}`, http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Failed to get inventory item: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Check sufficient quantity
	if inventoryQty < req.Quantity {
		http.Error(w, fmt.Sprintf(`{"error":"insufficient duplicates (have %d, need %d)"}`, inventoryQty, req.Quantity), http.StatusConflict)
		return
	}

	// Calculate new star rank (cap at 15)
	newStarRank := cmd.StarRank + req.Quantity
	if newStarRank > 15 {
		newStarRank = 15
		req.Quantity = 15 - cmd.StarRank // Only consume what's needed
	}

	// Update commander star rank
	_, err = tx.Exec(`
		UPDATE commanders
		SET star_rank = $1, updated_at = now()
		WHERE id = $2
	`, newStarRank, req.CommanderID)
	if err != nil {
		log.Printf("Failed to update commander: %v", err)
		http.Error(w, `{"error":"failed to merge"}`, http.StatusInternalServerError)
		return
	}

	// Consume duplicate cards from inventory
	if inventoryQty == req.Quantity {
		_, err = tx.Exec(`DELETE FROM player_inventory WHERE id = $1`, inventoryID)
	} else {
		_, err = tx.Exec(`
			UPDATE player_inventory
			SET quantity = quantity - $1, updated_at = now()
			WHERE id = $2
		`, req.Quantity, inventoryID)
	}
	if err != nil {
		log.Printf("Failed to consume duplicates: %v", err)
		http.Error(w, `{"error":"failed to consume duplicates"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Update quest progress
	services.UpdateQuestProgress(playerID, "merge_commander", "commander", 1)

	// Response
	type mergeResponse struct {
		NewStarRank      int    `json:"new_star_rank"`
		DuplicatesUsed   int    `json:"duplicates_used"`
		RemainingDuplicates int `json:"remaining_duplicates"`
		Message          string `json:"message"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mergeResponse{
		NewStarRank:         newStarRank,
		DuplicatesUsed:      req.Quantity,
		RemainingDuplicates: inventoryQty - req.Quantity,
		Message:             fmt.Sprintf("Commander merged! Star rank: %d → %d", cmd.StarRank, newStarRank),
	})
}

// AssignCommander handles POST /api/fleets/{id}/assign-commander
// Assigns commander to fleet (sets commander_id, marks commander as deployed)
func AssignCommander(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	fleetID := r.PathValue("id")

	var req struct {
		CommanderID string `json:"commander_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Verify fleet ownership
	var fleetExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM fleets WHERE id = $1 AND player_id = $2)`, fleetID, playerID).Scan(&fleetExists)
	if err != nil || !fleetExists {
		http.Error(w, `{"error":"fleet not found"}`, http.StatusNotFound)
		return
	}

	// Verify commander ownership and check if already deployed
	var isDeployed bool
	err = tx.QueryRow(`
		SELECT is_deployed FROM commanders
		WHERE id = $1 AND player_id = $2
		FOR UPDATE
	`, req.CommanderID, playerID).Scan(&isDeployed)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"commander not found"}`, http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Failed to get commander: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if isDeployed {
		http.Error(w, `{"error":"commander already deployed to another fleet"}`, http.StatusConflict)
		return
	}

	// Assign commander to fleet
	_, err = tx.Exec(`
		UPDATE fleets
		SET commander_id = $1, updated_at = now()
		WHERE id = $2
	`, req.CommanderID, fleetID)
	if err != nil {
		log.Printf("Failed to assign commander: %v", err)
		http.Error(w, `{"error":"failed to assign commander"}`, http.StatusInternalServerError)
		return
	}

	// Mark commander as deployed
	_, err = tx.Exec(`
		UPDATE commanders
		SET is_deployed = true, updated_at = now()
		WHERE id = $1
	`, req.CommanderID)
	if err != nil {
		log.Printf("Failed to mark commander deployed: %v", err)
		http.Error(w, `{"error":"failed to deploy commander"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Commander assigned to fleet"})
}

// UnassignCommander handles POST /api/fleets/{id}/unassign-commander
// Removes commander from fleet (clears commander_id, marks commander as not deployed)
func UnassignCommander(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	fleetID := r.PathValue("id")

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Get fleet with commander_id
	var commanderID sql.NullString
	err = tx.QueryRow(`
		SELECT commander_id FROM fleets
		WHERE id = $1 AND player_id = $2
		FOR UPDATE
	`, fleetID, playerID).Scan(&commanderID)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"fleet not found"}`, http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Failed to get fleet: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if !commanderID.Valid {
		http.Error(w, `{"error":"no commander assigned to this fleet"}`, http.StatusConflict)
		return
	}

	// Unassign commander from fleet
	_, err = tx.Exec(`
		UPDATE fleets
		SET commander_id = NULL, updated_at = now()
		WHERE id = $1
	`, fleetID)
	if err != nil {
		log.Printf("Failed to unassign commander: %v", err)
		http.Error(w, `{"error":"failed to unassign commander"}`, http.StatusInternalServerError)
		return
	}

	// Mark commander as not deployed
	_, err = tx.Exec(`
		UPDATE commanders
		SET is_deployed = false, updated_at = now()
		WHERE id = $1
	`, commanderID.String)
	if err != nil {
		log.Printf("Failed to mark commander undeployed: %v", err)
		http.Error(w, `{"error":"failed to undeploy commander"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Commander unassigned from fleet"})
}

// DismissCommander handles DELETE /api/commanders/{id}
// Deletes commander (cannot dismiss if deployed)
func DismissCommander(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	commanderID := r.PathValue("id")

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Check if commander is deployed
	var isDeployed bool
	err = tx.QueryRow(`
		SELECT is_deployed FROM commanders
		WHERE id = $1 AND player_id = $2
		FOR UPDATE
	`, commanderID, playerID).Scan(&isDeployed)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"commander not found"}`, http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Failed to get commander: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if isDeployed {
		http.Error(w, `{"error":"cannot dismiss deployed commander - unassign from fleet first"}`, http.StatusConflict)
		return
	}

	// Delete commander
	_, err = tx.Exec(`DELETE FROM commanders WHERE id = $1`, commanderID)
	if err != nil {
		log.Printf("Failed to delete commander: %v", err)
		http.Error(w, `{"error":"failed to dismiss commander"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Commander dismissed"})
}
