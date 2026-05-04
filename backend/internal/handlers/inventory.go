package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
	"github.com/cryptomines-online/backend/internal/services"
)

// InventoryItem represents an item in player's inventory with metadata from item_types
type InventoryItem struct {
	ID          string `json:"id"`
	ItemKey     string `json:"item_key"`
	Quantity    int    `json:"quantity"`
	DisplayName string `json:"display_name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	IconName    *string `json:"icon_name"`
	AcquiredAt  time.Time `json:"acquired_at"`
}

// UseItemResponse represents the response from using an item
type UseItemResponse struct {
	Effect    string          `json:"effect"`    // Human-readable effect description
	Resources *models.Resource `json:"resources,omitempty"` // Updated resources (if applicable)
}

// GetInventory handles GET /api/inventory
// Returns all items in player's inventory with metadata from item_types
func GetInventory(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	// Query inventory with item_types metadata
	rows, err := database.DB.Query(`
		SELECT pi.id, pi.item_key, pi.quantity, pi.acquired_at,
		       it.display_name, it.category, it.description, it.icon_name
		FROM player_inventory pi
		JOIN item_types it ON pi.item_key = it.item_key
		WHERE pi.player_id = $1 AND pi.quantity > 0
		ORDER BY it.category, it.display_name
	`, playerID)
	if err != nil {
		log.Printf("Failed to query inventory: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []InventoryItem{}
	for rows.Next() {
		var item InventoryItem
		err := rows.Scan(&item.ID, &item.ItemKey, &item.Quantity, &item.AcquiredAt,
			&item.DisplayName, &item.Category, &item.Description, &item.IconName)
		if err != nil {
			log.Printf("Failed to scan inventory row: %v", err)
			continue
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating inventory rows: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// ItemTypeInfo is the public metadata for an item type, used by tooltips
// and reward previews. Excludes player-scoped fields.
type ItemTypeInfo struct {
	ItemKey     string  `json:"item_key"`
	DisplayName string  `json:"display_name"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	IconName    *string `json:"icon_name"`
}

// ListItemTypes handles GET /api/item-types
// Returns the static catalog of item types (display_name, description, etc).
// Public — no auth required, used by quest reward tooltips and similar UI.
func ListItemTypes(w http.ResponseWriter, _ *http.Request) {
	rows, err := database.DB.Query(`
		SELECT item_key, display_name, category, description, icon_name
		FROM item_types
		ORDER BY category, display_name
	`)
	if err != nil {
		log.Printf("Failed to query item_types: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []ItemTypeInfo{}
	for rows.Next() {
		var it ItemTypeInfo
		if err := rows.Scan(&it.ItemKey, &it.DisplayName, &it.Category, &it.Description, &it.IconName); err != nil {
			log.Printf("Failed to scan item_type row: %v", err)
			continue
		}
		items = append(items, it)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// UseItem handles POST /api/inventory/{id}/use
// Consumes item and applies effect based on category
func UseItem(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	itemID := r.PathValue("id")

	// Begin transaction (atomicity: deduct item + apply effect)
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Get inventory item with lock (prevent double-use)
	var invItem struct {
		ItemKey  string
		Quantity int
		Category string
	}
	err = tx.QueryRow(`
		SELECT pi.item_key, pi.quantity, it.category
		FROM player_inventory pi
		JOIN item_types it ON pi.item_key = it.item_key
		WHERE pi.id = $1 AND pi.player_id = $2
		FOR UPDATE
	`, itemID, playerID).Scan(&invItem.ItemKey, &invItem.Quantity, &invItem.Category)

	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"item not found"}`, http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Failed to get inventory item: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Check quantity
	if invItem.Quantity <= 0 {
		http.Error(w, `{"error":"item out of stock"}`, http.StatusBadRequest)
		return
	}

	// Dispatch to category-specific use logic
	var effectMsg string
	var resources *models.Resource

	switch invItem.Category {
	case "resource_pack":
		effectMsg, resources, err = useResourcePack(tx, playerID, invItem.ItemKey)
	case "boost":
		effectMsg, err = useBoost(tx, playerID, invItem.ItemKey)
	case "battle":
		effectMsg, err = useBattleItem(tx, playerID, invItem.ItemKey)
	case "blueprint":
		effectMsg, err = useBlueprint(tx, playerID, invItem.ItemKey)
	case "commander":
		effectMsg, err = useCommanderCard(tx, playerID, invItem.ItemKey)
	default:
		err = fmt.Errorf("unknown item category: %s", invItem.Category)
	}

	if err != nil {
		log.Printf("Failed to apply item effect: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// Deduct item quantity (or delete if last one)
	if invItem.Quantity == 1 {
		_, err = tx.Exec(`DELETE FROM player_inventory WHERE id = $1`, itemID)
	} else {
		_, err = tx.Exec(`
			UPDATE player_inventory
			SET quantity = quantity - 1, updated_at = now()
			WHERE id = $1
		`, itemID)
	}
	if err != nil {
		log.Printf("Failed to deduct item: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UseItemResponse{
		Effect:    effectMsg,
		Resources: resources,
	})
}

// useResourcePack grants resources instantly (Metal, He3, Gold)
func useResourcePack(tx *sql.Tx, playerID, itemKey string) (string, *models.Resource, error) {
	// Get resource pack details
	var resourceType string
	var resourceAmount int64
	err := tx.QueryRow(`
		SELECT resource_type, resource_amount
		FROM item_types
		WHERE item_key = $1
	`, itemKey).Scan(&resourceType, &resourceAmount)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get resource pack details: %w", err)
	}

	// Get player's planet
	var planetID string
	err = tx.QueryRow(`SELECT id FROM planets WHERE player_id = $1 LIMIT 1`, playerID).Scan(&planetID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get planet: %w", err)
	}

	// Update resources
	var updateQuery string
	switch resourceType {
	case "metal":
		updateQuery = "UPDATE resources SET metal = metal + $1, updated_at = now() WHERE planet_id = $2 RETURNING *"
	case "he3":
		updateQuery = "UPDATE resources SET he3 = he3 + $1, updated_at = now() WHERE planet_id = $2 RETURNING *"
	case "gold":
		updateQuery = "UPDATE resources SET gold = gold + $1, updated_at = now() WHERE planet_id = $2 RETURNING *"
	default:
		return "", nil, fmt.Errorf("unknown resource type: %s", resourceType)
	}

	var res models.Resource
	err = tx.QueryRow(updateQuery, resourceAmount, planetID).Scan(
		&res.ID, &res.PlanetID, &res.Metal, &res.He3, &res.Gold,
		&res.MetalPerHour, &res.He3PerHour, &res.GoldPerHour,
		&res.StorageCapacity, &res.LastCollectedAt, &res.UpdatedAt,
		&res.WarehouseMetal, &res.WarehouseHe3, &res.WarehouseGold, &res.LastWarehouseUpdate,
	)
	if err != nil {
		return "", nil, fmt.Errorf("failed to grant resources: %w", err)
	}

	// Update quest progress
	services.UpdateQuestProgress(playerID, "use_item", "resource_pack", 1)

	effectMsg := fmt.Sprintf("Granted %d %s", resourceAmount, resourceType)
	return effectMsg, &res, nil
}

// useBoost creates an active buff (Construction Card, MVP Tool, etc.)
func useBoost(tx *sql.Tx, playerID, itemKey string) (string, error) {
	// Get boost details
	var boostType string
	var boostValue, durationHours int
	var displayName string
	err := tx.QueryRow(`
		SELECT boost_type, boost_value, duration_hours, display_name
		FROM item_types
		WHERE item_key = $1
	`, itemKey).Scan(&boostType, &boostValue, &durationHours, &displayName)
	if err != nil {
		return "", fmt.Errorf("failed to get boost details: %w", err)
	}

	// Calculate expiry
	expiresAt := time.Now().Add(time.Duration(durationHours) * time.Hour)

	// Insert or replace buff (ON CONFLICT replaces existing buff of same type)
	_, err = tx.Exec(`
		INSERT INTO active_buffs (player_id, buff_type, buff_value, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (player_id, buff_type)
		DO UPDATE SET buff_value = EXCLUDED.buff_value, expires_at = EXCLUDED.expires_at, created_at = now()
	`, playerID, boostType, boostValue, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to create buff: %w", err)
	}

	// Update quest progress
	services.UpdateQuestProgress(playerID, "use_item", "boost", 1)

	effectMsg := fmt.Sprintf("%s activated for %d hours", displayName, durationHours)
	return effectMsg, nil
}

// useBattleItem applies instant battle effects (SP Card, Truce Card)
func useBattleItem(tx *sql.Tx, playerID, itemKey string) (string, error) {
	// Get battle item details
	var battleEffect string
	var battleValue int
	var displayName string
	err := tx.QueryRow(`
		SELECT battle_effect, battle_value, display_name
		FROM item_types
		WHERE item_key = $1
	`, itemKey).Scan(&battleEffect, &battleValue, &displayName)
	if err != nil {
		return "", fmt.Errorf("failed to get battle item details: %w", err)
	}

	var effectMsg string

	switch battleEffect {
	case "sp_grant":
		// SP Card: Grant Space Points (cap at max SP)
		_, err = tx.Exec(`
			UPDATE players SET space_points = LEAST(space_points + $1, max_space_points)
			WHERE id = $2
		`, battleValue, playerID)
		if err != nil {
			return "", fmt.Errorf("failed to grant SP: %w", err)
		}
		effectMsg = fmt.Sprintf("Granted %d Space Points", battleValue)

	case "protection":
		// Truce Card: Grant protection from attacks for N hours
		// Cannot use if player has outgoing pending attacks
		var hasPending bool
		err = tx.QueryRow(`
			SELECT EXISTS(SELECT 1 FROM pending_attacks WHERE attacker_id = $1 AND status = 'traveling')
		`, playerID).Scan(&hasPending)
		if err == nil && hasPending {
			return "", fmt.Errorf("cannot activate truce while you have fleets attacking")
		}

		// Cannot use if incoming attacks are already in transit
		var hasIncoming bool
		err = tx.QueryRow(`
			SELECT EXISTS(SELECT 1 FROM pending_attacks WHERE defender_id = $1 AND status = 'traveling')
		`, playerID).Scan(&hasIncoming)
		if err == nil && hasIncoming {
			return "", fmt.Errorf("cannot activate truce during an imminent attack")
		}

		// Set protection_until on player's homeworld planet
		duration := time.Duration(battleValue) * time.Hour
		protUntil := time.Now().Add(duration)
		_, err = tx.Exec(`
			UPDATE planets SET protection_until = $1, updated_at = now()
			WHERE player_id = $2 AND is_homeworld = true
		`, protUntil, playerID)
		if err != nil {
			return "", fmt.Errorf("failed to activate protection: %w", err)
		}
		effectMsg = fmt.Sprintf("%d-hour truce protection activated until %s", battleValue, protUntil.Format("15:04 Jan 2"))

	default:
		return "", fmt.Errorf("unknown battle effect: %s", battleEffect)
	}

	// Update quest progress
	services.UpdateQuestProgress(playerID, "use_item", "battle", 1)

	return effectMsg, nil
}

// useBlueprint unlocks a blueprint (insert into player_blueprints)
func useBlueprint(tx *sql.Tx, playerID, itemKey string) (string, error) {
	// Get blueprint ID from item_types
	var blueprintID int
	var displayName string
	err := tx.QueryRow(`
		SELECT blueprint_id, display_name
		FROM item_types
		WHERE item_key = $1 AND category = 'blueprint'
	`, itemKey).Scan(&blueprintID, &displayName)
	if err != nil {
		return "", fmt.Errorf("failed to get blueprint details: %w", err)
	}

	// Check if player already has this blueprint
	var exists bool
	err = tx.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM player_blueprints WHERE player_id = $1 AND blueprint_id = $2)
	`, playerID, blueprintID).Scan(&exists)
	if err != nil {
		return "", fmt.Errorf("failed to check blueprint ownership: %w", err)
	}

	if exists {
		return "", fmt.Errorf("you already own this blueprint")
	}

	// Insert into player_blueprints with is_activated = true, research_level = 1
	_, err = tx.Exec(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 1)
	`, playerID, blueprintID)
	if err != nil {
		return "", fmt.Errorf("failed to unlock blueprint: %w", err)
	}

	// Get module/hull name for quest tracking (quest targets use snake_case names from module_types/hull_types)
	var questTarget string
	err = tx.QueryRow(`
		SELECT COALESCE(mt.name, ht.name, b.name)
		FROM blueprints b
		LEFT JOIN module_types mt ON b.module_type_id = mt.id
		LEFT JOIN hull_types ht ON b.hull_type_id = ht.id
		WHERE b.id = $1`, blueprintID).Scan(&questTarget)
	if err == nil {
		services.UpdateQuestProgress(playerID, "use_blueprint", questTarget, 1)
	}

	effectMsg := fmt.Sprintf("Blueprint unlocked: %s", displayName)
	return effectMsg, nil
}

// useCommanderCard unlocks a commander (insert into commanders or add star rank)
func useCommanderCard(tx *sql.Tx, playerID, itemKey string) (string, error) {
	// Get commander card details
	var commanderTypeName string
	var displayName string
	err := tx.QueryRow(`
		SELECT commander_type, display_name
		FROM item_types
		WHERE item_key = $1 AND category = 'commander'
	`, itemKey).Scan(&commanderTypeName, &displayName)
	if err != nil {
		return "", fmt.Errorf("failed to get commander card details: %w", err)
	}

	// Look up the commander type from reference table
	var ctID int
	var ctName, ctRarity string
	var ctAccuracy, ctDodge, ctSpeed, ctElectron int
	err = tx.QueryRow(`
		SELECT id, name, rarity, base_accuracy, base_dodge, base_speed, base_electron
		FROM commander_types
		WHERE name = $1
	`, commanderTypeName).Scan(&ctID, &ctName, &ctRarity, &ctAccuracy, &ctDodge, &ctSpeed, &ctElectron)
	if err != nil {
		return "", fmt.Errorf("commander type not found: %s", commanderTypeName)
	}

	// Check commander count limit
	var commanderCount int
	err = tx.QueryRow(`SELECT COUNT(*) FROM commanders WHERE player_id = $1`, playerID).Scan(&commanderCount)
	if err != nil {
		return "", fmt.Errorf("failed to count commanders: %w", err)
	}
	if commanderCount >= 60 {
		return "", fmt.Errorf("max commanders reached (60/60)")
	}

	// Check if player already owns this commander
	var existingID string
	err = tx.QueryRow(`SELECT id FROM commanders WHERE player_id = $1 AND name = $2`, playerID, ctName).Scan(&existingID)

	var effectMsg string
	if err == sql.ErrNoRows {
		// New commander — insert
		_, err = tx.Exec(`
			INSERT INTO commanders (player_id, name, rarity, star_rank, accuracy, dodge, speed, electron, is_deployed)
			VALUES ($1, $2, $3, 0, $4, $5, $6, $7, false)
		`, playerID, ctName, ctRarity, ctAccuracy, ctDodge, ctSpeed, ctElectron)
		if err != nil {
			return "", fmt.Errorf("failed to create commander: %w", err)
		}
		effectMsg = fmt.Sprintf("Commander unlocked: %s (%s)", displayName, ctRarity)
	} else if err == nil {
		// Duplicate — increase star rank by 1 (max 15)
		_, err = tx.Exec(`
			UPDATE commanders SET star_rank = LEAST(star_rank + 1, 15), updated_at = now()
			WHERE id = $1
		`, existingID)
		if err != nil {
			return "", fmt.Errorf("failed to upgrade commander: %w", err)
		}
		effectMsg = fmt.Sprintf("Duplicate %s! Star rank increased by 1.", displayName)
	} else {
		return "", fmt.Errorf("failed to check commander: %w", err)
	}

	// Update quest progress
	services.UpdateQuestProgress(playerID, "unlock_commander", "commander", 1)

	return effectMsg, nil
}

// GetActiveBuffs handles GET /api/player/buffs
// Returns active buffs with remaining time, cleaning up expired ones
func GetActiveBuffs(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	// Clean up expired buffs
	database.DB.Exec(`DELETE FROM active_buffs WHERE player_id = $1 AND expires_at <= now()`, playerID)

	type ActiveBuff struct {
		BuffType  string `json:"buff_type"`
		BuffValue int    `json:"buff_value"`
		ExpiresAt string `json:"expires_at"`
	}

	rows, err := database.DB.Query(`
		SELECT buff_type, buff_value, expires_at
		FROM active_buffs
		WHERE player_id = $1 AND expires_at > now()
		ORDER BY expires_at ASC
	`, playerID)
	if err != nil {
		log.Printf("GetActiveBuffs: query failed: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	buffs := []ActiveBuff{}
	for rows.Next() {
		var b ActiveBuff
		var expiresAt time.Time
		if err := rows.Scan(&b.BuffType, &b.BuffValue, &expiresAt); err != nil {
			continue
		}
		b.ExpiresAt = expiresAt.Format(time.RFC3339)
		buffs = append(buffs, b)
	}

	// Also get truce protection status from planet
	var protectionUntil sql.NullTime
	database.DB.QueryRow(`
		SELECT protection_until FROM planets
		WHERE player_id = $1 AND is_homeworld = true
	`, playerID).Scan(&protectionUntil)

	type BuffsResponse struct {
		Buffs           []ActiveBuff `json:"buffs"`
		ProtectionUntil *string      `json:"protection_until"`
	}

	resp := BuffsResponse{Buffs: buffs}
	if protectionUntil.Valid && protectionUntil.Time.After(time.Now()) {
		t := protectionUntil.Time.Format(time.RFC3339)
		resp.ProtectionUntil = &t
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
