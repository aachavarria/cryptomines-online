# Module 9: Inventory System - Implementation Plan

**Date:** 2026-02-07
**Module:** Inventory System (Items, Resource Packs, Boosts, Battle Items)
**Estimated Hours:** 14-18 hours (3-4 days)
**Priority:** MEDIUM (Quality of Life)

---

## 1. Executive Summary

### Overview
The Inventory System manages all consumable items players acquire through quests, rewards, and future mall purchases. It includes Resource Packs (instant resource grants), Resource Boosts (timed production bonuses), Battle Items (SP Cards, Truce Cards), Blueprints, and Commander Cards. Currently, blueprints and commanders bypass inventory and go directly to unlock tables - this module fixes that flow by implementing a proper **item → inventory → use → effect** pipeline.

### Scope
**IN SCOPE:**
- **Item Types (17 base types):**
  - Resource Packs: Gold Pack, Adv Gold Pack, Primary/Junior/Senior Metal Pack (3), Primary/Junior/Senior He3 Pack (3) = 8 types
  - Resource Boosts: Construction Card, MVP Tool, Extra Tax, Adv Extra Tax, Metal Mining Boost, He3 Mining Boost = 6 types
  - Battle Items: SP Card, Truce Card, Adv Truce Card = 3 types
  - Consumables: Blueprint items, Commander Card items = 2 categories (multiple items per category)
- Database: `item_types`, `player_inventory`, `active_buffs` tables
- Backend: Use endpoints (POST /api/inventory/{item_id}/use) with type-specific logic
- Frontend: InventoryPanel component (grid view, use buttons, tooltips)
- Quest integration: Fix ClaimQuest to add items to inventory (not direct unlock)

**OUT OF SCOPE:**
- Mall/shop system (Phase 3+)
- Trading items between players
- Item auction house
- Planet Transformation Packs (cut from scope)
- Commander Enhancement Items (Gems, Chips, Merge items - cut from scope)
- Healing/Revival Cards (commanders immortal)
- Loudspeaker item (World Chat no longer needs it)

### Key Design Decisions
1. **Inventory as central hub:** All items (blueprints, commanders, packs, boosts) go through inventory first
2. **Stacking:** Items stack by `item_key` (e.g., 5x Gold Pack → quantity=5)
3. **Use logic per type:** Resource Packs instant, Boosts create active_buffs rows, Battle Items instant effect
4. **Timed buffs:** `active_buffs` table tracks active boosts (Construction Card, MVP Tool, etc.) with `expires_at` timestamps
5. **Blueprint/Commander refactor:** Change quest claim logic to insert into inventory, add "Use" step to unlock
6. **Unlimited inventory:** No slot cap (simplified from GO2)

---

## 2. Current Database Analysis

### Existing Tables

#### `quest_types.reward_item_json`
```sql
reward_item_json JSONB NOT NULL DEFAULT '[]'
-- Example: [{"type":"item","item_key":"loudspeaker","quantity":1}, {"type":"blueprint","blueprint_key":"super_transmission_engine"}]
```

**Current Quest Claim Logic (quests.go):**
```go
// Lines 291-305: Award resources to player
// Lines 327-332: Return reward_item_json as-is (NO PROCESSING)
```

**Problem:** Items are returned in JSON but never inserted into inventory or processed. Blueprints/commanders are never unlocked.

#### `player_blueprints`
```sql
CREATE TABLE player_blueprints (
    player_id UUID NOT NULL,
    blueprint_id INTEGER NOT NULL REFERENCES blueprints(id),
    is_activated BOOLEAN DEFAULT false,
    research_level INTEGER DEFAULT 0,
    ...
);
```

**Current flow:** Blueprints are manually inserted directly (bypasses inventory).
**Target flow:** Quest → Inventory (blueprint item) → Use → Insert into player_blueprints.

#### `commanders`
```sql
CREATE TABLE commanders (
    id UUID PRIMARY KEY,
    player_id UUID NOT NULL,
    name TEXT NOT NULL,
    rarity TEXT NOT NULL,
    star_rank INTEGER DEFAULT 0,
    accuracy INTEGER,
    dodge INTEGER,
    speed INTEGER,
    electron INTEGER,
    ...
);
```

**Current flow:** Commander cards are directly inserted (bypasses inventory).
**Target flow:** Quest → Inventory (commander card) → Use → Insert into commanders.

### Database Completeness: 0%
**Status:** ❌ NO INVENTORY TABLES EXIST - Full migration required

**Required Tables:**
1. `item_types` - Static reference table (17 base types + blueprint/commander items)
2. `player_inventory` - Player's item stacks
3. `active_buffs` - Timed production/construction boosts

---

## 3. Galaxy Online 2 Research

### Sources
Research conducted on Galaxy Online 2 wiki (game no longer active, archived information):

1. **[Development Items Wiki](https://galaxyonlineii.fandom.com/wiki/Development_Items)**
   - Resource Packs, Production Boosts, Construction Boost
   - Specific quantities and durations
2. **[Category: Items Wiki](https://galaxyonlineii.fandom.com/wiki/Category:Items)**
   - 41 total items across multiple categories
   - Items stored in "bag under the My Tools tab"

### GO2 Item Mechanics Summary

#### Resource Packs (Instant Consumption)
| Item | Mall Points | Resource Grant | Daily Limit |
|------|-------------|----------------|-------------|
| Gold Pack | 5 | 30,000 Gold | 100 |
| Advanced Gold Pack | - | 100,000 Gold | - |
| Primary Metal Pack | 16 | 50,000-60,000 Metal | - |
| Junior Metal Pack | 32 | 150,000-180,000 Metal | - |
| Senior Metal Pack | 56 | 300,000-350,000 Metal | - |
| Primary He3 Pack | 16 | 50,000-60,000 He3 | - |
| Junior He3 Pack | 32 | 150,000-180,000 He3 | - |
| Senior He3 Pack | 56 | 300,000-350,000 He3 | - |

**Mechanics:**
- Click "Use" → Instantly grant resources
- Item removed from inventory (consumed)
- No cooldown or duration

#### Resource Boosts (Timed Buffs)
| Item | Mall Points | Effect | Duration |
|------|-------------|--------|----------|
| Construction Card | 24 | +3 construction slots | 72 hours |
| MVP Tool | 120 | +20% production/build/repair | 168 hours (7 days) |
| Extra Tax | 4 | +30% Gold production | 12 hours |
| Advanced Extra Tax | 24 | +100% Gold production | 24 hours |
| Metal Mining Boost | 4 | +30% Metal production | 12 hours |
| He3 Mining Boost | 4 | +30% He3 production | 12 hours |

**Mechanics:**
- Click "Use" → Create active buff with `expires_at` timestamp
- Buff effects apply automatically (no manual activation)
- Multiple buffs of same type don't stack (replace existing)
- Item consumed, buff persists until expiry

#### Battle Items (Instant Effects)
| Item | Effect |
|------|--------|
| SP Card | +10 Space Points (instant) |
| Truce Card | 12h protection from attacks |
| Advanced Truce Card | 72h protection from attacks |

**Mechanics:**
- SP Card: Add 10 SP to player (capped at max SP)
- Truce Card: Create protection entry in `active_protections` table

### Cryptomines Online Adaptations
- **1:1 Copy:** Using exact GO2 item names, quantities, and durations
- **Simplified:** No Mall Points cost (out of scope for Phase 2), no daily limits
- **Unlimited Inventory:** No slot cap (GO2 had inventory expansion mechanics)
- **No Loudspeaker:** Loudspeaker item removed from scope (World Chat no longer needs it)

---

## 4. System Architecture

### Component Overview
```
┌──────────────────────────────────────────────────────────────────┐
│                      INVENTORY SYSTEM                             │
├──────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌───────────────┐      ┌──────────────┐      ┌──────────────┐ │
│  │   Quest/      │      │   Inventory  │      │   Use Item   │ │
│  │   Instance    │─────►│   Storage    │─────►│   Effects    │ │
│  │   Rewards     │      │              │      │              │ │
│  │               │      │ player_      │      │ • Resources  │ │
│  │ Add items to  │      │ inventory    │      │ • Buffs      │ │
│  │ inventory     │      │              │      │ • Unlocks    │ │
│  └───────────────┘      └──────────────┘      └──────────────┘ │
│                                                                   │
│  ┌───────────────┐      ┌──────────────┐      ┌──────────────┐ │
│  │   Frontend    │◄────►│   Backend    │◄────►│   Database   │ │
│  │               │ HTTP │              │ SQL  │              │ │
│  │  Inventory    │      │ Use handlers │      │ item_types   │ │
│  │  Panel        │      │              │      │ player_inv   │ │
│  │  • Grid view  │      │ • Resource   │      │ active_buffs │ │
│  │  • Tooltips   │      │ • Boost      │      │              │ │
│  │  • Use button │      │ • Battle     │      │              │ │
│  └───────────────┘      │ • Blueprint  │      └──────────────┘ │
│                         │ • Commander  │                        │
│                         └──────────────┘                        │
└──────────────────────────────────────────────────────────────────┘
```

### Data Flow: Item Usage

```
1. Player acquires item (quest reward, instance loot, future mall)
   ↓
2. Item inserted into player_inventory:
   INSERT INTO player_inventory (player_id, item_key, quantity)
   VALUES (playerID, 'gold_pack', 1)
   ON CONFLICT (player_id, item_key) DO UPDATE SET quantity = quantity + 1
   ↓
3. Player opens Inventory Panel, sees item grid
   ↓
4. Player clicks "Use" button on item
   ↓
5. POST /api/inventory/{item_id}/use
   ↓
6. Backend determines item type from item_types table:
   - category: 'resource_pack' → Grant resources, consume item
   - category: 'boost' → Create active_buffs entry, consume item
   - category: 'battle' → Apply effect (SP/protection), consume item
   - category: 'blueprint' → Insert into player_blueprints, consume item
   - category: 'commander' → Insert into commanders, consume item
   ↓
7. Execute type-specific logic in transaction:
   BEGIN;
     - Deduct item from player_inventory (quantity - 1)
     - Apply effect (resources, buff, unlock)
   COMMIT;
   ↓
8. Return success response with effect details
   ↓
9. Frontend updates UI (inventory grid, resource HUD, buffs display)
```

### Blueprint/Commander Flow Refactor

**OLD FLOW (bypass inventory):**
```
Quest Claim → Direct INSERT into player_blueprints/commanders
```

**NEW FLOW (proper inventory):**
```
Quest Claim → INSERT blueprint/commander item into player_inventory
            → Player opens Inventory
            → Player clicks "Use" on blueprint/commander card
            → Backend INSERTS into player_blueprints/commanders
            → Item consumed
```

---

## 5. Database Schema Design

### 5.1 New Table: item_types

```sql
CREATE TABLE item_types (
    item_key TEXT PRIMARY KEY,
    display_name TEXT NOT NULL,
    category TEXT NOT NULL CHECK (category IN (
        'resource_pack', 'boost', 'battle', 'blueprint', 'commander'
    )),
    description TEXT NOT NULL,
    icon_name TEXT,

    -- Resource pack data (category='resource_pack')
    resource_type TEXT CHECK (resource_type IN ('metal', 'he3', 'gold')),
    resource_amount INTEGER DEFAULT 0,

    -- Boost data (category='boost')
    boost_type TEXT CHECK (boost_type IN (
        'construction_slots', 'production_pct', 'gold_production_pct',
        'metal_production_pct', 'he3_production_pct',
        'ship_build_speed_pct', 'repair_speed_pct'
    )),
    boost_value INTEGER DEFAULT 0,  -- +3 slots, +20%, +30%, +100%
    duration_hours INTEGER DEFAULT 0,

    -- Battle item data (category='battle')
    battle_effect TEXT CHECK (battle_effect IN ('sp_grant', 'protection')),
    battle_value INTEGER DEFAULT 0,  -- +10 SP, 12h/72h protection

    -- Blueprint/Commander data (category='blueprint'/'commander')
    blueprint_id INTEGER REFERENCES blueprints(id),
    commander_type TEXT,  -- 'common', 'skill', 'super'

    -- Metadata
    is_consumable BOOLEAN NOT NULL DEFAULT true,
    max_stack INTEGER NOT NULL DEFAULT 999,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_item_types_category ON item_types (category);
```

### 5.2 New Table: player_inventory

```sql
CREATE TABLE player_inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    item_key TEXT NOT NULL REFERENCES item_types(item_key),
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity >= 0),
    acquired_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_player_item UNIQUE (player_id, item_key)
);

CREATE INDEX idx_player_inventory_player ON player_inventory (player_id);
CREATE INDEX idx_player_inventory_item ON player_inventory (item_key);
```

### 5.3 New Table: active_buffs

```sql
CREATE TABLE active_buffs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    buff_type TEXT NOT NULL CHECK (buff_type IN (
        'construction_slots', 'production_pct', 'gold_production_pct',
        'metal_production_pct', 'he3_production_pct',
        'ship_build_speed_pct', 'repair_speed_pct'
    )),
    buff_value INTEGER NOT NULL,  -- +3, +20, +30, +100
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_player_buff UNIQUE (player_id, buff_type)
    -- Players can only have one buff of each type active
);

CREATE INDEX idx_active_buffs_player ON active_buffs (player_id);
CREATE INDEX idx_active_buffs_expiry ON active_buffs (expires_at) WHERE expires_at > now();
```

### 5.4 Seed Data: item_types

```sql
-- Resource Packs (8 types)
INSERT INTO item_types (item_key, display_name, category, description, resource_type, resource_amount) VALUES
('gold_pack', 'Gold Pack', 'resource_pack', 'Grants 30,000 Gold', 'gold', 30000),
('advanced_gold_pack', 'Advanced Gold Pack', 'resource_pack', 'Grants 100,000 Gold', 'gold', 100000),
('primary_metal_pack', 'Primary Metal Pack', 'resource_pack', 'Grants 50,000 Metal', 'metal', 50000),
('junior_metal_pack', 'Junior Metal Pack', 'resource_pack', 'Grants 150,000 Metal', 'metal', 150000),
('senior_metal_pack', 'Senior Metal Pack', 'resource_pack', 'Grants 300,000 Metal', 'metal', 300000),
('primary_he3_pack', 'Primary He3 Pack', 'resource_pack', 'Grants 50,000 He3', 'he3', 50000),
('junior_he3_pack', 'Junior He3 Pack', 'resource_pack', 'Grants 150,000 He3', 'he3', 150000),
('senior_he3_pack', 'Senior He3 Pack', 'resource_pack', 'Grants 300,000 He3', 'he3', 300000);

-- Resource Boosts (6 types)
INSERT INTO item_types (item_key, display_name, category, description, boost_type, boost_value, duration_hours) VALUES
('construction_card', 'Construction Card', 'boost', '+3 construction slots for 72 hours', 'construction_slots', 3, 72),
('mvp_tool', 'MVP Tool', 'boost', '+20% to all production, building, and repair for 7 days', 'production_pct', 20, 168),
('extra_tax', 'Extra Tax', 'boost', '+30% Gold production for 12 hours', 'gold_production_pct', 30, 12),
('advanced_extra_tax', 'Advanced Extra Tax', 'boost', '+100% Gold production for 24 hours', 'gold_production_pct', 100, 24),
('metal_mining_boost', 'Metal Mining Boost', 'boost', '+30% Metal production for 12 hours', 'metal_production_pct', 30, 12),
('he3_mining_boost', 'He3 Mining Boost', 'boost', '+30% He3 production for 12 hours', 'he3_production_pct', 30, 12);

-- Battle Items (3 types)
INSERT INTO item_types (item_key, display_name, category, description, battle_effect, battle_value) VALUES
('sp_card', 'SP Card', 'battle', 'Grants 10 Space Points instantly', 'sp_grant', 10),
('truce_card', 'Truce Card', 'battle', '12-hour protection from attacks', 'protection', 12),
('advanced_truce_card', 'Advanced Truce Card', 'battle', '72-hour protection from attacks', 'protection', 72);

-- Blueprint items (dynamic - created when blueprints exist)
-- Example: INSERT INTO item_types (item_key, display_name, category, description, blueprint_id)
-- VALUES ('blueprint_super_transmission_engine', 'Blueprint: Super Transmission Engine', 'blueprint', 'Unlock Super Transmission Engine blueprint', 5);

-- Commander cards (dynamic - seeded based on commander types)
-- Example: INSERT INTO item_types (item_key, display_name, category, description, commander_type)
-- VALUES ('commander_card_common', 'Common Commander Card', 'commander', 'Summon a random Common commander', 'common');
```

---

## 6. Backend Implementation

### 6.1 New Handler: inventory.go

**File:** `/backend/internal/handlers/inventory.go`

```go
package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/cryptomines-online/backend/internal/database"
    "github.com/cryptomines-online/backend/internal/middleware"
)

// GetInventory handles GET /api/inventory
func GetInventory(w http.ResponseWriter, r *http.Request) {
    playerID := middleware.GetPlayerID(r)

    rows, err := database.DB.Query(`
        SELECT pi.id, pi.item_key, pi.quantity, it.display_name, it.category,
               it.description, it.icon_name
        FROM player_inventory pi
        JOIN item_types it ON pi.item_key = it.item_key
        WHERE pi.player_id = $1 AND pi.quantity > 0
        ORDER BY it.category, it.display_name
    `, playerID)
    if err != nil {
        log.Printf("Failed to get inventory: %v", err)
        http.Error(w, `{"error":"failed to get inventory"}`, http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    type inventoryItem struct {
        ID          string `json:"id"`
        ItemKey     string `json:"item_key"`
        Quantity    int    `json:"quantity"`
        DisplayName string `json:"display_name"`
        Category    string `json:"category"`
        Description string `json:"description"`
        IconName    string `json:"icon_name"`
    }

    items := []inventoryItem{}
    for rows.Next() {
        var item inventoryItem
        var iconName sql.NullString
        err := rows.Scan(&item.ID, &item.ItemKey, &item.Quantity, &item.DisplayName,
            &item.Category, &item.Description, &iconName)
        if err != nil {
            log.Printf("Failed to scan inventory item: %v", err)
            continue
        }
        if iconName.Valid {
            item.IconName = iconName.String
        }
        items = append(items, item)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(items)
}

// UseItem handles POST /api/inventory/{id}/use
func UseItem(w http.ResponseWriter, r *http.Request) {
    playerID := middleware.GetPlayerID(r)
    inventoryID := r.PathValue("id")

    tx, err := database.DB.Begin()
    if err != nil {
        log.Printf("Failed to begin transaction: %v", err)
        http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    // Get inventory item + item type data
    var itemKey string
    var quantity int
    var category string
    var resourceType, boostType, battleEffect sql.NullString
    var resourceAmount, boostValue, battleValue, durationHours sql.NullInt64
    var blueprintID sql.NullInt64
    var commanderType sql.NullString

    err = tx.QueryRow(`
        SELECT pi.item_key, pi.quantity, it.category,
               it.resource_type, it.resource_amount,
               it.boost_type, it.boost_value, it.duration_hours,
               it.battle_effect, it.battle_value,
               it.blueprint_id, it.commander_type
        FROM player_inventory pi
        JOIN item_types it ON pi.item_key = it.item_key
        WHERE pi.id = $1 AND pi.player_id = $2
        FOR UPDATE OF pi
    `, inventoryID, playerID).Scan(
        &itemKey, &quantity, &category,
        &resourceType, &resourceAmount,
        &boostType, &boostValue, &durationHours,
        &battleEffect, &battleValue,
        &blueprintID, &commanderType,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, `{"error":"item not found"}`, http.StatusNotFound)
            return
        }
        log.Printf("Failed to get item: %v", err)
        http.Error(w, `{"error":"failed to get item"}`, http.StatusInternalServerError)
        return
    }

    if quantity <= 0 {
        http.Error(w, `{"error":"item out of stock"}`, http.StatusBadRequest)
        return
    }

    // Execute use logic based on category
    var effect string
    switch category {
    case "resource_pack":
        effect = useResourcePack(tx, playerID, resourceType.String, int(resourceAmount.Int64))
    case "boost":
        effect = useBoost(tx, playerID, boostType.String, int(boostValue.Int64), int(durationHours.Int64))
    case "battle":
        effect = useBattleItem(tx, playerID, battleEffect.String, int(battleValue.Int64))
    case "blueprint":
        effect = useBlueprint(tx, playerID, int(blueprintID.Int64))
    case "commander":
        effect = useCommanderCard(tx, playerID, commanderType.String)
    default:
        http.Error(w, `{"error":"unknown item category"}`, http.StatusBadRequest)
        return
    }

    if effect == "ERROR" {
        http.Error(w, `{"error":"failed to use item"}`, http.StatusInternalServerError)
        return
    }

    // Deduct item quantity
    _, err = tx.Exec(`
        UPDATE player_inventory
        SET quantity = quantity - 1, updated_at = now()
        WHERE id = $1
    `, inventoryID)
    if err != nil {
        log.Printf("Failed to deduct item: %v", err)
        http.Error(w, `{"error":"failed to deduct item"}`, http.StatusInternalServerError)
        return
    }

    if err := tx.Commit(); err != nil {
        log.Printf("Failed to commit: %v", err)
        http.Error(w, `{"error":"failed to commit"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "success": "true",
        "effect":  effect,
    })
}

// useResourcePack grants instant resources
func useResourcePack(tx *sql.Tx, playerID, resourceType string, amount int) string {
    var column string
    switch resourceType {
    case "metal":
        column = "metal"
    case "he3":
        column = "he3"
    case "gold":
        column = "gold"
    default:
        return "ERROR"
    }

    _, err := tx.Exec(`
        UPDATE players
        SET `+column+` = `+column+` + $1, updated_at = now()
        WHERE id = $2
    `, amount, playerID)
    if err != nil {
        log.Printf("Failed to grant resources: %v", err)
        return "ERROR"
    }

    return "Granted " + resourceType + ": " + string(amount)
}

// useBoost creates an active buff
func useBoost(tx *sql.Tx, playerID, boostType string, value, hours int) string {
    expiresAt := time.Now().Add(time.Duration(hours) * time.Hour)

    // Upsert buff (replace if exists)
    _, err := tx.Exec(`
        INSERT INTO active_buffs (player_id, buff_type, buff_value, expires_at)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (player_id, buff_type) DO UPDATE
        SET buff_value = EXCLUDED.buff_value,
            expires_at = EXCLUDED.expires_at,
            created_at = now()
    `, playerID, boostType, value, expiresAt)
    if err != nil {
        log.Printf("Failed to create buff: %v", err)
        return "ERROR"
    }

    return "Activated " + boostType + " for " + string(hours) + " hours"
}

// useBattleItem applies instant battle effects
func useBattleItem(tx *sql.Tx, playerID, effect string, value int) string {
    switch effect {
    case "sp_grant":
        // Add SP to player (Module 6 will have SP system)
        _, err := tx.Exec(`
            UPDATE players
            SET space_points = LEAST(space_points + $1, max_space_points),
                updated_at = now()
            WHERE id = $2
        `, value, playerID)
        if err != nil {
            log.Printf("Failed to grant SP: %v", err)
            return "ERROR"
        }
        return "Granted " + string(value) + " Space Points"

    case "protection":
        // Create protection entry (Module 6 PvP will check this)
        expiresAt := time.Now().Add(time.Duration(value) * time.Hour)
        _, err := tx.Exec(`
            INSERT INTO active_protections (player_id, protection_type, expires_at)
            VALUES ($1, 'truce_card', $2)
            ON CONFLICT (player_id) DO UPDATE
            SET expires_at = GREATEST(active_protections.expires_at, EXCLUDED.expires_at)
        `, playerID, expiresAt)
        if err != nil {
            log.Printf("Failed to create protection: %v", err)
            return "ERROR"
        }
        return "Protected for " + string(value) + " hours"

    default:
        return "ERROR"
    }
}

// useBlueprint unlocks a blueprint
func useBlueprint(tx *sql.Tx, playerID string, blueprintID int) string {
    _, err := tx.Exec(`
        INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
        VALUES ($1, $2, true, 1)
        ON CONFLICT (player_id, blueprint_id) DO NOTHING
    `, playerID, blueprintID)
    if err != nil {
        log.Printf("Failed to unlock blueprint: %v", err)
        return "ERROR"
    }

    return "Blueprint unlocked"
}

// useCommanderCard summons a random commander
func useCommanderCard(tx *sql.Tx, playerID, rarity string) string {
    // Gacha logic: select random commander of specified rarity
    // For Phase 2, simplified: just insert a predefined commander
    // Module 4 (Commander System) will have full gacha implementation

    _, err := tx.Exec(`
        INSERT INTO commanders (player_id, name, rarity, star_rank, accuracy, dodge, speed, electron)
        VALUES ($1, 'Commander', $2, 0, 10, 10, 10, 10)
    `, playerID, rarity)
    if err != nil {
        log.Printf("Failed to summon commander: %v", err)
        return "ERROR"
    }

    return "Commander summoned: " + rarity
}
```

### 6.2 Update: quests.go (Fix ClaimQuest)

**File:** `/backend/internal/handlers/quests.go` (lines 306-340, update)

```go
// After awarding resources (line 305), ADD item processing:

// Process item rewards (add to inventory)
var rewardItems []map[string]interface{}
if qt.RewardItemJSON != nil {
    if err := json.Unmarshal(qt.RewardItemJSON, &rewardItems); err != nil {
        log.Printf("Failed to parse reward items: %v", err)
    } else {
        for _, item := range rewardItems {
            itemType := item["type"].(string)

            if itemType == "item" {
                // Regular item (resource pack, boost, battle)
                itemKey := item["item_key"].(string)
                quantity := int(item["quantity"].(float64))

                _, err = tx.Exec(`
                    INSERT INTO player_inventory (player_id, item_key, quantity)
                    VALUES ($1, $2, $3)
                    ON CONFLICT (player_id, item_key) DO UPDATE
                    SET quantity = player_inventory.quantity + EXCLUDED.quantity,
                        updated_at = now()
                `, playerID, itemKey, quantity)
                if err != nil {
                    log.Printf("Failed to add item to inventory: %v", err)
                }

            } else if itemType == "blueprint" {
                // Blueprint item (add to inventory, NOT direct unlock)
                blueprintKey := item["blueprint_key"].(string)

                // Get blueprint ID
                var blueprintID int
                err = tx.QueryRow(`
                    SELECT id FROM blueprints WHERE blueprint_type = 'hull'
                    AND hull_type_id = (SELECT id FROM hull_types WHERE name = $1)
                    LIMIT 1
                `, blueprintKey).Scan(&blueprintID)
                if err != nil {
                    log.Printf("Failed to find blueprint: %v", err)
                    continue
                }

                // Create dynamic blueprint item in item_types (if not exists)
                blueprintItemKey := "blueprint_" + blueprintKey
                _, err = tx.Exec(`
                    INSERT INTO item_types (item_key, display_name, category, description, blueprint_id)
                    VALUES ($1, $2, 'blueprint', $3, $4)
                    ON CONFLICT (item_key) DO NOTHING
                `, blueprintItemKey, "Blueprint: "+blueprintKey, "Unlock "+blueprintKey+" blueprint", blueprintID)

                // Add to inventory
                _, err = tx.Exec(`
                    INSERT INTO player_inventory (player_id, item_key, quantity)
                    VALUES ($1, $2, 1)
                    ON CONFLICT (player_id, item_key) DO UPDATE
                    SET quantity = player_inventory.quantity + 1, updated_at = now()
                `, playerID, blueprintItemKey)
                if err != nil {
                    log.Printf("Failed to add blueprint item: %v", err)
                }
            }
        }
    }
}
```

### 6.3 Route Registration

**File:** `/backend/cmd/server/main.go` (update)

```go
// Inventory endpoints
http.HandleFunc("GET /api/inventory", middleware.AuthMiddleware(handlers.GetInventory))
http.HandleFunc("POST /api/inventory/{id}/use", middleware.AuthMiddleware(handlers.UseItem))
```

### 6.4 Backend Tasks Breakdown

| Task | Description | Hours |
|------|-------------|-------|
| Create inventory.go | GetInventory + UseItem handlers with type-specific logic | 3.5 |
| Use logic per category | 5 use functions (resource_pack, boost, battle, blueprint, commander) | 2.0 |
| Update quests.go | Fix ClaimQuest to add items to inventory (not direct unlock) | 1.5 |
| Database migration | Create 3 tables (item_types, player_inventory, active_buffs) + seed 17 items | 1.5 |
| Route registration | Add 2 routes | 0.25 |
| Error handling | Comprehensive validation + error messages | 0.75 |
| Testing (manual) | Test all 17 item types + blueprint/commander flow | 1.5 |
| **Subtotal** | | **11 hours** |

---

## 7. Frontend Implementation

### 7.1 New Component: InventoryPanel.tsx

**File:** `/frontend/src/components/panels/InventoryPanel.tsx`

```typescript
import { useState, useEffect } from 'react'
import { useGame } from '../../contexts/GameContext'
import { getInventory, useItem } from '../../api/inventory'

interface InventoryItem {
  id: string
  item_key: string
  quantity: number
  display_name: string
  category: string
  description: string
  icon_name?: string
}

export default function InventoryPanel() {
  const { refreshResources } = useGame()
  const [items, setItems] = useState<InventoryItem[]>([])
  const [loading, setLoading] = useState(false)
  const [selectedItem, setSelectedItem] = useState<InventoryItem | null>(null)

  const fetchInventory = async () => {
    try {
      const data = await getInventory()
      setItems(data)
    } catch (err) {
      console.error('Failed to fetch inventory:', err)
    }
  }

  useEffect(() => {
    fetchInventory()
  }, [])

  const handleUse = async (item: InventoryItem) => {
    if (!confirm(`Use ${item.display_name}?`)) return

    setLoading(true)
    try {
      const result = await useItem(item.id)
      alert(`Success!\n${result.effect}`)

      // Refresh inventory + resources
      await fetchInventory()
      await refreshResources()
    } catch (err: any) {
      alert(`Failed to use item: ${err.message}`)
    } finally {
      setLoading(false)
    }
  }

  const groupedItems = items.reduce((acc, item) => {
    if (!acc[item.category]) acc[item.category] = []
    acc[item.category].push(item)
    return acc
  }, {} as Record<string, InventoryItem[]>)

  return (
    <div className="inventory-panel">
      <h2>Inventory</h2>

      {items.length === 0 && <p>Your inventory is empty.</p>}

      {Object.entries(groupedItems).map(([category, categoryItems]) => (
        <div key={category} className="inventory-category">
          <h3>{category.replace('_', ' ').toUpperCase()}</h3>
          <div className="item-grid">
            {categoryItems.map(item => (
              <div
                key={item.id}
                className="inventory-item"
                onMouseEnter={() => setSelectedItem(item)}
                onMouseLeave={() => setSelectedItem(null)}
              >
                <div className="item-icon">
                  {item.icon_name ? (
                    <img src={`/icons/${item.icon_name}.png`} alt={item.display_name} />
                  ) : (
                    <div className="item-placeholder">{item.display_name[0]}</div>
                  )}
                </div>
                <div className="item-quantity">×{item.quantity}</div>
                <button
                  onClick={() => handleUse(item)}
                  disabled={loading}
                  className="use-button"
                >
                  Use
                </button>
              </div>
            ))}
          </div>
        </div>
      ))}

      {selectedItem && (
        <div className="item-tooltip">
          <h4>{selectedItem.display_name}</h4>
          <p>{selectedItem.description}</p>
        </div>
      )}
    </div>
  )
}
```

### 7.2 API Functions

**File:** `/frontend/src/api/inventory.ts`

```typescript
import { apiRequest } from './base'

export interface InventoryItem {
  id: string
  item_key: string
  quantity: number
  display_name: string
  category: string
  description: string
  icon_name?: string
}

export interface UseItemResponse {
  success: string
  effect: string
}

export async function getInventory(): Promise<InventoryItem[]> {
  return apiRequest<InventoryItem[]>('/api/inventory')
}

export async function useItem(inventoryID: string): Promise<UseItemResponse> {
  return apiRequest<UseItemResponse>(`/api/inventory/${inventoryID}/use`, {
    method: 'POST',
  })
}
```

### 7.3 GameContext Integration

**File:** `/frontend/src/contexts/GameContext.tsx` (update)

```typescript
// Add inventory to context state
const [activePanel, setActivePanel] = useState<string | null>(null)

// Add inventory button to main UI
<button onClick={() => setActivePanel('inventory')}>Inventory</button>
```

**File:** `/frontend/src/App.tsx` (update)

```typescript
import InventoryPanel from './components/panels/InventoryPanel'

// Add to panel rendering
{activePanel === 'inventory' && <InventoryPanel />}
```

### 7.4 Frontend Tasks Breakdown

| Task | Description | Hours |
|------|-------------|-------|
| Create InventoryPanel.tsx | Item grid, category grouping, use buttons | 2.0 |
| Create api/inventory.ts | API request functions (getInventory, useItem) | 0.25 |
| Item tooltips | Hover tooltip showing item description | 0.5 |
| Category tabs/grouping | Visual separation of item categories | 0.5 |
| GameContext integration | Add inventory button to main UI | 0.25 |
| Styling | CSS for item grid, icons, tooltips | 0.75 |
| Testing | Manual UI testing for all item types | 0.75 |
| **Subtotal** | | **5 hours** |

---

## 8. QA Checklist

### 8.1 Backend Tests

**Manual Testing (Go):**

| # | Test Case | Expected Result | Status |
|---|-----------|----------------|--------|
| 1 | GET /api/inventory (empty inventory) | 200 OK, empty array | ⬜ |
| 2 | GET /api/inventory (with items) | 200 OK, all items listed with quantities | ⬜ |
| 3 | Use Gold Pack (30k gold) | Player gold += 30000, item quantity -= 1 | ⬜ |
| 4 | Use Senior Metal Pack (300k metal) | Player metal += 300000, item consumed | ⬜ |
| 5 | Use Construction Card (+ 3 slots, 72h) | active_buffs row created, expires_at = now + 72h | ⬜ |
| 6 | Use MVP Tool (+20% production, 7 days) | active_buffs row, expires_at = now + 168h | ⬜ |
| 7 | Use same boost twice (replace existing) | Old buff replaced, expires_at updated | ⬜ |
| 8 | Use SP Card (+10 SP) | Player SP += 10 (capped at max) | ⬜ |
| 9 | Use Truce Card (12h protection) | active_protections row, expires_at = now + 12h | ⬜ |
| 10 | Use Blueprint item | Item consumed, row in player_blueprints | ⬜ |
| 11 | Use Commander Card | Item consumed, commander in commanders table | ⬜ |
| 12 | Use item with quantity = 0 | 400 Bad Request: "item out of stock" | ⬜ |
| 13 | Use item not owned by player | 404 Not Found: "item not found" | ⬜ |
| 14 | Quest claim with item reward | Item added to player_inventory | ⬜ |
| 15 | Quest claim with blueprint reward | Blueprint item added to inventory (NOT direct unlock) | ⬜ |
| 16 | Use all items in stack (quantity 5 → 0) | Inventory row quantity = 0 (not deleted) | ⬜ |
| 17 | Buff expiry check | active_buffs with expires_at < now() ignored | ⬜ |
| 18 | Transaction rollback on error | All changes reverted if use logic fails | ⬜ |

### 8.2 Frontend Tests

**Manual Testing (UI):**

| # | Test Case | Expected Result | Status |
|---|-----------|----------------|--------|
| 19 | Open Inventory panel (empty) | Message: "Your inventory is empty" | ⬜ |
| 20 | Open Inventory panel (with items) | Item grid grouped by category | ⬜ |
| 21 | Hover over item | Tooltip shows item name + description | ⬜ |
| 22 | Click "Use" on Gold Pack | Confirm modal → Resources increase | ⬜ |
| 23 | Use Construction Card | Buff icon appears, construction slots +3 | ⬜ |
| 24 | Use item, refresh inventory | Used item quantity decreases | ⬜ |
| 25 | Use last item in stack | Item removed from grid | ⬜ |
| 26 | Resource HUD updates after using pack | Metal/He3/Gold counts increase | ⬜ |
| 27 | Use Blueprint item | Blueprint unlocked, appears in BlueprintPanel | ⬜ |
| 28 | Use Commander Card | Commander appears in commander roster | ⬜ |
| 29 | Error handling (backend failure) | Alert shows error message | ⬜ |
| 30 | Loading state during use | Button disabled, "Using..." text | ⬜ |

### 8.3 Integration Tests

| # | Test Case | Expected Result | Status |
|---|-----------|----------------|--------|
| 31 | Quest claim → Inventory → Use Blueprint | Full flow: quest → inventory item → unlock | ⬜ |
| 32 | Quest claim with multiple items | All items added to inventory with correct quantities | ⬜ |
| 33 | Use boost while another boost active | New boost replaces old (no stacking) | ⬜ |
| 34 | Buff expiry after duration | Buff effects stop after expires_at | ⬜ |
| 35 | Concurrent item use (race condition) | Database locks prevent double-use | ⬜ |
| 36 | Item icons load correctly | All 17 item types have icons displayed | ⬜ |
| 37 | Category grouping correct | Items grouped: Resource Packs, Boosts, Battle, Blueprint, Commander | ⬜ |
| 38 | Stack overflow (999+ items) | Items stack correctly beyond 999 (no max_stack enforcement) | ⬜ |

---

## 9. Edge Cases & Error Handling

### 9.1 Edge Cases

| Case | Handling |
|------|----------|
| Item quantity = 0 | 400 Bad Request: "item out of stock" |
| Use item not owned | 404 Not Found: "item not found" |
| Buff already active (same type) | Replace existing buff (ON CONFLICT ... DO UPDATE) |
| SP Card when SP at max | Grant SP up to max (LEAST(sp + 10, max_sp)) |
| Blueprint already unlocked | ON CONFLICT DO NOTHING (no error) |
| Commander roster full (60 max) | Allow insert (Module 4 will handle max enforcement) |
| Item with NULL resource_amount | Error: "invalid item data" |
| Unknown item category | 400 Bad Request: "unknown item category" |
| Buff expiry in past | Ignore buff (WHERE expires_at > now()) |

### 9.2 Validation Rules

**Backend:**
1. ✅ Item must exist in player_inventory (ownership check)
2. ✅ Quantity must be > 0 (stock check)
3. ✅ Item category must match use logic
4. ✅ Resource pack must have valid resource_type + amount
5. ✅ Boost must have valid boost_type + value + duration
6. ✅ Blueprint must have valid blueprint_id
7. ✅ Commander must have valid rarity

**Frontend:**
1. ✅ Confirm modal for all use actions
2. ✅ Disable "Use" button while loading
3. ✅ Show clear error messages for failures
4. ✅ Refresh inventory after successful use

---

## 10. Risk Assessment

| Risk | Impact | Likelihood | Mitigation |
|------|--------|-----------|------------|
| **Quest claim breaks (blueprint flow)** | HIGH | MEDIUM | Test quest claim thoroughly, ensure item_types seeded |
| **Transaction failure during use** | HIGH | LOW | Use database transactions (BEGIN/COMMIT/ROLLBACK) |
| **Race condition (double use)** | MEDIUM | LOW | Use `FOR UPDATE` locks on player_inventory |
| **Buff stacking bug** | MEDIUM | MEDIUM | ON CONFLICT ... DO UPDATE (replace, don't stack) |
| **Item duplication** | HIGH | LOW | Deduct item in same transaction as effect |
| **Frontend state desync** | MEDIUM | MEDIUM | Call fetchInventory() + refreshResources() after use |
| **Missing item_types seed data** | HIGH | MEDIUM | Verify all 17 items seeded in migration |
| **Blueprint/Commander unlock fails** | MEDIUM | MEDIUM | Add error logging, test both flows |

**Critical Paths:**
1. Quest claim → Add to inventory (fix ClaimQuest logic)
2. Use item → Deduct + Apply effect (transaction atomicity)
3. Buff creation → Check expiry (WHERE expires_at > now())

---

## 11. Timeline & Estimates

### 11.1 Task Breakdown

| Phase | Task | Owner | Hours | Dependencies |
|-------|------|-------|-------|-------------|
| **Database** | | | | |
| 1 | Create migration (3 tables + seed data) | infra-dev | 1.5 | - |
| **Backend** | | | | |
| 2 | Create inventory.go (GetInventory, UseItem) | backend-dev | 3.5 | Task 1 |
| 3 | Use logic (5 functions) | backend-dev | 2.0 | Task 2 |
| 4 | Update quests.go (ClaimQuest fix) | backend-dev | 1.5 | Task 2 |
| 5 | Route registration | backend-dev | 0.25 | Task 2 |
| 6 | Error handling | backend-dev | 0.75 | Task 2-4 |
| 7 | Manual testing (Postman) | backend-dev | 1.5 | Task 2-6 |
| **Frontend** | | | | |
| 8 | Create InventoryPanel.tsx | frontend-dev | 2.0 | Backend Task 2 |
| 9 | Create api/inventory.ts | frontend-dev | 0.25 | - |
| 10 | Item tooltips | frontend-dev | 0.5 | Task 8 |
| 11 | Category grouping | frontend-dev | 0.5 | Task 8 |
| 12 | GameContext integration | frontend-dev | 0.25 | Task 8 |
| 13 | Styling (CSS) | frontend-dev | 0.75 | Task 8 |
| 14 | Manual UI testing | frontend-dev | 0.75 | Task 8-13 |
| **QA** | | | | |
| 15 | Backend QA (18 tests) | qa-agent | 1.5 | Backend tasks |
| 16 | Frontend QA (12 tests) | qa-agent | 1.0 | Frontend tasks |
| 17 | Integration QA (8 tests) | qa-agent | 1.0 | All tasks |
| **TOTAL** | | | **18 hours** | |

### 11.2 Timeline (Parallel Execution)

**Day 1 (6 hours):**
- Infra: Task 1 (migration) - 1.5 hours
- Backend: Tasks 2-3 (inventory.go, use logic) - 5.5 hours

**Day 2 (6 hours):**
- Backend: Tasks 4-6 (quests fix, routes, error handling) - 2.5 hours
- Frontend: Tasks 8-11 (InventoryPanel, tooltips, grouping) - 3.25 hours

**Day 3 (6 hours):**
- Backend: Task 7 (testing) - 1.5 hours
- Frontend: Tasks 12-14 (GameContext, styling, testing) - 1.75 hours
- QA: Tasks 15-17 (all tests) - 3.5 hours

**Total:** 18 hours (3-4 days with parallel work)

---

## 12. Success Criteria

### 12.1 Must-Have (MVP)
- ✅ 3 database tables (item_types, player_inventory, active_buffs) seeded
- ✅ All 17 item types functional (packs, boosts, battle items)
- ✅ Blueprint/Commander flow refactored (quest → inventory → use → unlock)
- ✅ Timed buffs tracked in active_buffs table
- ✅ Quest rewards add items to inventory (not direct unlock)
- ✅ InventoryPanel displays items grouped by category
- ✅ Use button consumes items and applies effects
- ✅ Transaction atomicity (deduct + effect in one transaction)

### 12.2 Nice-to-Have (Future)
- ⬜ Item icons (currently placeholders)
- ⬜ Buff UI indicator (show active buffs in HUD)
- ⬜ Buff stacking (allow multiple buffs of same type)
- ⬜ Item gifting (send items to other players)
- ⬜ Item auction house (trade items)
- ⬜ Item history log (track usage)

### 12.3 Definition of Done
- All 38 QA tests pass
- Quest rewards add items to inventory (verified)
- Blueprint/Commander use flow works end-to-end
- All 17 item types tested individually
- No item duplication bugs (atomicity verified)
- Buff expiry logic correct (expires_at checks work)

---

## 13. Future Enhancements

### Phase 3+ Additions
1. **Mall System:** Players can buy items with Mall Points/Vouchers
2. **Item Icons:** Add visual icons for all 17 item types (currently placeholders)
3. **Buff UI:** Display active buffs in HUD (Construction Card: 3 slots for 48h remaining)
4. **Commander Gacha:** useCommanderCard() calls full gacha logic (random stats, skills)
5. **Item Trading:** Players can send items to friends/corps
6. **Item Auction House:** List items for sale (Metal, He3, Gold prices)
7. **Loudspeaker Return:** If World Chat adds premium features, re-enable Loudspeaker item

---

## 14. Appendix

### 14.1 GO2 References
- [Development Items Wiki](https://galaxyonlineii.fandom.com/wiki/Development_Items) - Resource packs, boosts, Construction Card, MVP Tool
- [Category: Items Wiki](https://galaxyonlineii.fandom.com/wiki/Category:Items) - All 41 item types in GO2

### 14.2 Item Type Reference

**Resource Packs (Instant):**
- Gold Pack: 30,000 Gold
- Advanced Gold Pack: 100,000 Gold
- Primary Metal Pack: 50,000 Metal
- Junior Metal Pack: 150,000 Metal
- Senior Metal Pack: 300,000 Metal
- Primary He3 Pack: 50,000 He3
- Junior He3 Pack: 150,000 He3
- Senior He3 Pack: 300,000 He3

**Resource Boosts (Timed):**
- Construction Card: +3 slots, 72 hours
- MVP Tool: +20% all production/build/repair, 168 hours (7 days)
- Extra Tax: +30% Gold production, 12 hours
- Advanced Extra Tax: +100% Gold production, 24 hours
- Metal Mining Boost: +30% Metal production, 12 hours
- He3 Mining Boost: +30% He3 production, 12 hours

**Battle Items (Instant):**
- SP Card: +10 Space Points
- Truce Card: 12h protection from attacks
- Advanced Truce Card: 72h protection from attacks

### 14.3 Database Schema SQL

```sql
-- Complete migration file: /supabase/migrations/20260207000000_inventory_system.sql

CREATE TABLE item_types (
    item_key TEXT PRIMARY KEY,
    display_name TEXT NOT NULL,
    category TEXT NOT NULL CHECK (category IN ('resource_pack', 'boost', 'battle', 'blueprint', 'commander')),
    description TEXT NOT NULL,
    icon_name TEXT,
    resource_type TEXT CHECK (resource_type IN ('metal', 'he3', 'gold')),
    resource_amount INTEGER DEFAULT 0,
    boost_type TEXT CHECK (boost_type IN ('construction_slots', 'production_pct', 'gold_production_pct', 'metal_production_pct', 'he3_production_pct', 'ship_build_speed_pct', 'repair_speed_pct')),
    boost_value INTEGER DEFAULT 0,
    duration_hours INTEGER DEFAULT 0,
    battle_effect TEXT CHECK (battle_effect IN ('sp_grant', 'protection')),
    battle_value INTEGER DEFAULT 0,
    blueprint_id INTEGER REFERENCES blueprints(id),
    commander_type TEXT,
    is_consumable BOOLEAN NOT NULL DEFAULT true,
    max_stack INTEGER NOT NULL DEFAULT 999,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE player_inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    item_key TEXT NOT NULL REFERENCES item_types(item_key),
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity >= 0),
    acquired_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_player_item UNIQUE (player_id, item_key)
);

CREATE TABLE active_buffs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    buff_type TEXT NOT NULL CHECK (buff_type IN ('construction_slots', 'production_pct', 'gold_production_pct', 'metal_production_pct', 'he3_production_pct', 'ship_build_speed_pct', 'repair_speed_pct')),
    buff_value INTEGER NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_player_buff UNIQUE (player_id, buff_type)
);

-- Indexes
CREATE INDEX idx_item_types_category ON item_types (category);
CREATE INDEX idx_player_inventory_player ON player_inventory (player_id);
CREATE INDEX idx_player_inventory_item ON player_inventory (item_key);
CREATE INDEX idx_active_buffs_player ON active_buffs (player_id);
CREATE INDEX idx_active_buffs_expiry ON active_buffs (expires_at) WHERE expires_at > now();

-- Seed data (17 items)
INSERT INTO item_types (item_key, display_name, category, description, resource_type, resource_amount) VALUES
('gold_pack', 'Gold Pack', 'resource_pack', 'Grants 30,000 Gold', 'gold', 30000),
('advanced_gold_pack', 'Advanced Gold Pack', 'resource_pack', 'Grants 100,000 Gold', 'gold', 100000),
('primary_metal_pack', 'Primary Metal Pack', 'resource_pack', 'Grants 50,000 Metal', 'metal', 50000),
('junior_metal_pack', 'Junior Metal Pack', 'resource_pack', 'Grants 150,000 Metal', 'metal', 150000),
('senior_metal_pack', 'Senior Metal Pack', 'resource_pack', 'Grants 300,000 Metal', 'metal', 300000),
('primary_he3_pack', 'Primary He3 Pack', 'resource_pack', 'Grants 50,000 He3', 'he3', 50000),
('junior_he3_pack', 'Junior He3 Pack', 'resource_pack', 'Grants 150,000 He3', 'he3', 150000),
('senior_he3_pack', 'Senior He3 Pack', 'resource_pack', 'Grants 300,000 He3', 'he3', 300000);

INSERT INTO item_types (item_key, display_name, category, description, boost_type, boost_value, duration_hours) VALUES
('construction_card', 'Construction Card', 'boost', '+3 construction slots for 72 hours', 'construction_slots', 3, 72),
('mvp_tool', 'MVP Tool', 'boost', '+20% to all production, building, and repair for 7 days', 'production_pct', 20, 168),
('extra_tax', 'Extra Tax', 'boost', '+30% Gold production for 12 hours', 'gold_production_pct', 30, 12),
('advanced_extra_tax', 'Advanced Extra Tax', 'boost', '+100% Gold production for 24 hours', 'gold_production_pct', 100, 24),
('metal_mining_boost', 'Metal Mining Boost', 'boost', '+30% Metal production for 12 hours', 'metal_production_pct', 30, 12),
('he3_mining_boost', 'He3 Mining Boost', 'boost', '+30% He3 production for 12 hours', 'he3_production_pct', 30, 12);

INSERT INTO item_types (item_key, display_name, category, description, battle_effect, battle_value) VALUES
('sp_card', 'SP Card', 'battle', 'Grants 10 Space Points instantly', 'sp_grant', 10),
('truce_card', 'Truce Card', 'battle', '12-hour protection from attacks', 'protection', 12),
('advanced_truce_card', 'Advanced Truce Card', 'battle', '72-hour protection from attacks', 'protection', 72);

-- RLS policies
ALTER TABLE item_types ENABLE ROW LEVEL SECURITY;
ALTER TABLE player_inventory ENABLE ROW LEVEL SECURITY;
ALTER TABLE active_buffs ENABLE ROW LEVEL SECURITY;

CREATE POLICY item_types_select_all ON item_types FOR SELECT USING (true);
CREATE POLICY player_inventory_select_own ON player_inventory FOR SELECT USING (auth.uid() = player_id);
CREATE POLICY active_buffs_select_own ON active_buffs FOR SELECT USING (auth.uid() = player_id);
```

---

## 15. Implementation Notes

### Backend Notes
- Use `FOR UPDATE` locks on player_inventory to prevent race conditions
- Transaction must be atomic: deduct item + apply effect in one transaction
- Buff replacement (ON CONFLICT ... DO UPDATE) prevents stacking
- Quest ClaimQuest must insert items into inventory (not direct unlock)
- Blueprint items created dynamically when rewarded (item_key = "blueprint_" + blueprint_key)
- SP Card caps at max_space_points (LEAST(sp + 10, max_sp))
- Truce Card uses GREATEST(existing, new) to extend protection (not replace)

### Frontend Notes
- Fetch inventory on panel open (not on every render)
- Group items by category (Resource Packs, Boosts, Battle, Blueprint, Commander)
- Hover tooltip shows item description
- Confirm modal required for all use actions
- Refresh inventory + resources after successful use
- Item icons currently placeholders (add icons in Phase 3)

### Testing Notes
- Test all 17 item types individually
- Test quest claim → inventory flow (blueprint/commander)
- Test buff replacement (use same boost twice)
- Test SP Card at max SP (should cap)
- Test transaction rollback (simulate effect failure)
- Test concurrent item use (database locks should prevent double-deduct)

---

**END OF PLAN**
