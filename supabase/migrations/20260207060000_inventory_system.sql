-- Migration: Inventory System (Items, Resource Packs, Boosts, Battle Items)
-- Module 9: Inventory System
-- Date: 2026-02-07
-- Priority: MEDIUM (Quality of Life - Phase A: Foundation)

-- ============================================================================
-- TABLE: item_types (Static Reference Table)
-- ============================================================================
-- Stores all consumable item definitions (17 base types + dynamic blueprint/commander items)
-- Categories: resource_pack, boost, battle, blueprint, commander

CREATE TABLE item_types (
    item_key TEXT PRIMARY KEY,
    display_name TEXT NOT NULL,
    category TEXT NOT NULL CHECK (category IN ('resource_pack', 'boost', 'battle', 'blueprint', 'commander')),
    description TEXT NOT NULL,
    icon_name TEXT, -- Placeholder for future icon assets (Phase 3+)

    -- Resource Pack fields (category = 'resource_pack')
    resource_type TEXT CHECK (resource_type IN ('metal', 'he3', 'gold')),
    resource_amount INTEGER DEFAULT 0 CHECK (resource_amount >= 0),

    -- Boost fields (category = 'boost')
    boost_type TEXT CHECK (boost_type IN ('construction_slots', 'production_pct', 'gold_production_pct', 'metal_production_pct', 'he3_production_pct', 'ship_build_speed_pct', 'repair_speed_pct')),
    boost_value INTEGER DEFAULT 0 CHECK (boost_value >= 0),
    duration_hours INTEGER DEFAULT 0 CHECK (duration_hours >= 0),

    -- Battle Item fields (category = 'battle')
    battle_effect TEXT CHECK (battle_effect IN ('sp_grant', 'protection')),
    battle_value INTEGER DEFAULT 0 CHECK (battle_value >= 0),

    -- Blueprint/Commander fields (dynamically created)
    blueprint_id INTEGER REFERENCES blueprints(id),
    commander_type TEXT, -- e.g., 'random_common', 'random_skill', 'specific_zeus'

    -- Meta fields
    is_consumable BOOLEAN NOT NULL DEFAULT true, -- All items consumed on use (no reusable items yet)
    max_stack INTEGER NOT NULL DEFAULT 999, -- Maximum quantity per stack (unused, inventory unlimited)
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================================
-- TABLE: player_inventory (Player Item Stacks)
-- ============================================================================
-- Stores players' item quantities. Items stack by item_key.
-- Unique constraint prevents duplicate entries (player_id, item_key).

CREATE TABLE player_inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    item_key TEXT NOT NULL REFERENCES item_types(item_key),
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity >= 0),
    acquired_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_player_item UNIQUE (player_id, item_key)
);

-- ============================================================================
-- TABLE: active_buffs (Timed Production/Construction Boosts)
-- ============================================================================
-- Tracks active buffs from boost items (Construction Card, MVP Tool, etc.)
-- Only one buff of each type can be active per player (UNIQUE constraint).
-- Buffs expire when expires_at <= now().

CREATE TABLE active_buffs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    buff_type TEXT NOT NULL CHECK (buff_type IN ('construction_slots', 'production_pct', 'gold_production_pct', 'metal_production_pct', 'he3_production_pct', 'ship_build_speed_pct', 'repair_speed_pct')),
    buff_value INTEGER NOT NULL CHECK (buff_value >= 0),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_player_buff UNIQUE (player_id, buff_type)
);

-- ============================================================================
-- INDEXES
-- ============================================================================
CREATE INDEX idx_item_types_category ON item_types (category);
CREATE INDEX idx_player_inventory_player ON player_inventory (player_id);
CREATE INDEX idx_player_inventory_item ON player_inventory (item_key);
CREATE INDEX idx_active_buffs_player ON active_buffs (player_id);
CREATE INDEX idx_active_buffs_expiry ON active_buffs (expires_at);

-- ============================================================================
-- SEED DATA: 17 Base Item Types
-- ============================================================================

-- Resource Packs (8 items) - Instant resource grants
INSERT INTO item_types (item_key, display_name, category, description, resource_type, resource_amount) VALUES
('gold_pack', 'Gold Pack', 'resource_pack', 'Grants 30,000 Gold instantly', 'gold', 30000),
('advanced_gold_pack', 'Advanced Gold Pack', 'resource_pack', 'Grants 100,000 Gold instantly', 'gold', 100000),
('primary_metal_pack', 'Primary Metal Pack', 'resource_pack', 'Grants 50,000 Metal instantly', 'metal', 50000),
('junior_metal_pack', 'Junior Metal Pack', 'resource_pack', 'Grants 150,000 Metal instantly', 'metal', 150000),
('senior_metal_pack', 'Senior Metal Pack', 'resource_pack', 'Grants 300,000 Metal instantly', 'metal', 300000),
('primary_he3_pack', 'Primary He3 Pack', 'resource_pack', 'Grants 50,000 He3 instantly', 'he3', 50000),
('junior_he3_pack', 'Junior He3 Pack', 'resource_pack', 'Grants 150,000 He3 instantly', 'he3', 150000),
('senior_he3_pack', 'Senior He3 Pack', 'resource_pack', 'Grants 300,000 He3 instantly', 'he3', 300000);

-- Resource Boosts (6 items) - Timed production/construction buffs
INSERT INTO item_types (item_key, display_name, category, description, boost_type, boost_value, duration_hours) VALUES
('construction_card', 'Construction Card', 'boost', 'Grants +3 construction slots for 72 hours', 'construction_slots', 3, 72),
('mvp_tool', 'MVP Tool', 'boost', 'Grants +20% to all production, building, and repair for 7 days', 'production_pct', 20, 168),
('extra_tax', 'Extra Tax', 'boost', 'Grants +30% Gold production for 12 hours', 'gold_production_pct', 30, 12),
('advanced_extra_tax', 'Advanced Extra Tax', 'boost', 'Grants +100% Gold production for 24 hours', 'gold_production_pct', 100, 24),
('metal_mining_boost', 'Metal Mining Boost', 'boost', 'Grants +30% Metal production for 12 hours', 'metal_production_pct', 30, 12),
('he3_mining_boost', 'He3 Mining Boost', 'boost', 'Grants +30% He3 production for 12 hours', 'he3_production_pct', 30, 12);

-- Battle Items (3 items) - Instant effects (SP, protection)
INSERT INTO item_types (item_key, display_name, category, description, battle_effect, battle_value) VALUES
('sp_card', 'SP Card', 'battle', 'Grants 10 Space Points instantly', 'sp_grant', 10),
('truce_card', 'Truce Card', 'battle', 'Grants 12-hour protection from attacks', 'protection', 12),
('advanced_truce_card', 'Advanced Truce Card', 'battle', 'Grants 72-hour protection from attacks', 'protection', 72);

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================
ALTER TABLE item_types ENABLE ROW LEVEL SECURITY;
ALTER TABLE player_inventory ENABLE ROW LEVEL SECURITY;
ALTER TABLE active_buffs ENABLE ROW LEVEL SECURITY;

-- item_types: Public read access (all players can see all item definitions)
CREATE POLICY item_types_select_all ON item_types FOR SELECT USING (true);

-- player_inventory: Players can only see their own items
CREATE POLICY player_inventory_select_own ON player_inventory FOR SELECT
    USING (player_id IN (SELECT id FROM players WHERE anonymous_id = auth.uid()::text));

-- active_buffs: Players can only see their own buffs
CREATE POLICY active_buffs_select_own ON active_buffs FOR SELECT
    USING (player_id IN (SELECT id FROM players WHERE anonymous_id = auth.uid()::text));

-- ============================================================================
-- COMMENTS
-- ============================================================================
COMMENT ON TABLE item_types IS 'Static reference table for all consumable items (17 base types + dynamic blueprint/commander items)';
COMMENT ON TABLE player_inventory IS 'Player item stacks. Items consumed on use (quantity decremented). UNIQUE per player_id + item_key.';
COMMENT ON TABLE active_buffs IS 'Active timed buffs from boost items. Only one buff of each type per player. Expires when expires_at <= now().';

COMMENT ON COLUMN item_types.category IS 'Item category: resource_pack (instant), boost (timed buff), battle (SP/protection), blueprint (unlock), commander (unlock)';
COMMENT ON COLUMN item_types.boost_type IS 'Type of buff for boost items: construction_slots, production_pct, gold_production_pct, etc.';
COMMENT ON COLUMN item_types.battle_effect IS 'Effect type for battle items: sp_grant (add SP), protection (truce card)';
COMMENT ON COLUMN player_inventory.quantity IS 'Stack quantity. Decremented on use. Row deleted when reaches 0.';
COMMENT ON COLUMN active_buffs.buff_type IS 'Must match item_types.boost_type. UNIQUE constraint prevents stacking same buff type.';
COMMENT ON COLUMN active_buffs.expires_at IS 'Buff expires when expires_at <= now(). Check before applying buff bonuses.';
