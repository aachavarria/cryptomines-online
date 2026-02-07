-- Quest System Migration
-- Tables: quest_types, player_quests, daily_quest_progress
-- Seed Data: main quest chain, side quests, daily quests

-- =============================================================================
-- 3.22 Quest Types (Reference Data)
-- =============================================================================
CREATE TABLE quest_types (
    id SERIAL PRIMARY KEY,
    quest_key TEXT UNIQUE NOT NULL,
    category TEXT NOT NULL CHECK (category IN ('main', 'side', 'daily')),
    display_name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    requirement_type TEXT NOT NULL CHECK (requirement_type IN (
        'build_building', 'upgrade_building', 'harvest_resources', 'research_tech',
        'use_blueprint', 'recruit_commander', 'create_ship_design', 'build_ships',
        'create_fleet', 'replenish_ammo', 'build_defense', 'complete_instance',
        'reach_production', 'reach_storage', 'own_ships', 'upgrade_blueprints',
        'upgrade_modules', 'reach_star_rank', 'recycle_ships', 'use_speedup',
        'login', 'use_truce_card', 'donate_resources', 'send_message',
        'add_friend', 'increase_bag_slot', 'use_resource_pack', 'grow_comsats'
    )),
    requirement_target TEXT NOT NULL DEFAULT '',
    -- e.g., 'metal_collector', 'technology_center', or '' for generic
    requirement_value INTEGER NOT NULL DEFAULT 1,
    -- e.g., build level, quantity threshold, production rate target
    chain_order INTEGER,
    -- for main quests: sequential position (1-28+); for side quests: tier level; NULL for daily
    prerequisite_quest_id INTEGER REFERENCES quest_types(id) ON DELETE SET NULL,
    -- main quests: previous quest in chain; side quests: previous tier; daily: NULL
    reward_metal BIGINT NOT NULL DEFAULT 0,
    reward_he3 BIGINT NOT NULL DEFAULT 0,
    reward_gold BIGINT NOT NULL DEFAULT 0,
    reward_item_json JSONB NOT NULL DEFAULT '[]',
    -- e.g., [{"type": "blueprint", "blueprint_id": 5}, {"type": "item", "item_key": "loudspeaker", "quantity": 1}]
    phase INTEGER NOT NULL DEFAULT 1 CHECK (phase IN (1, 2, 3)),
    -- which implementation phase this quest belongs to
    is_active BOOLEAN NOT NULL DEFAULT true,
    -- false for deferred quests (social/chat-dependent)

    CONSTRAINT chk_chain_order_positive CHECK (chain_order IS NULL OR chain_order >= 1),
    CONSTRAINT chk_rewards_non_negative CHECK (reward_metal >= 0 AND reward_he3 >= 0 AND reward_gold >= 0)
);

CREATE INDEX idx_quest_types_category ON quest_types (category);
CREATE INDEX idx_quest_types_phase ON quest_types (phase);
CREATE INDEX idx_quest_types_chain ON quest_types (category, chain_order);

-- =============================================================================
-- 3.23 Player Quests (Progress Tracking)
-- =============================================================================
CREATE TABLE player_quests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    quest_type_id INTEGER NOT NULL REFERENCES quest_types(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'locked'
        CHECK (status IN ('locked', 'available', 'in_progress', 'completed', 'claimed')),
    progress_value INTEGER NOT NULL DEFAULT 0,
    -- current progress toward requirement_value (e.g., ships built, production reached)
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    claimed_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_player_quest UNIQUE (player_id, quest_type_id),
    CONSTRAINT chk_progress_non_negative CHECK (progress_value >= 0),
    CONSTRAINT chk_status_consistency CHECK (
        (status = 'locked' AND completed_at IS NULL AND claimed_at IS NULL) OR
        (status = 'available' AND completed_at IS NULL AND claimed_at IS NULL) OR
        (status = 'in_progress' AND completed_at IS NULL AND claimed_at IS NULL) OR
        (status = 'completed' AND completed_at IS NOT NULL AND claimed_at IS NULL) OR
        (status = 'claimed' AND completed_at IS NOT NULL AND claimed_at IS NOT NULL)
    )
);

CREATE INDEX idx_player_quests_player ON player_quests (player_id);
CREATE INDEX idx_player_quests_status ON player_quests (player_id, status);
CREATE INDEX idx_player_quests_quest_type ON player_quests (quest_type_id);

-- =============================================================================
-- 3.24 Daily Quest Progress
-- =============================================================================
CREATE TABLE daily_quest_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    quest_date DATE NOT NULL DEFAULT CURRENT_DATE,
    daily_points INTEGER NOT NULL DEFAULT 0,
    quests_completed_json JSONB NOT NULL DEFAULT '{}',
    -- e.g., {"daily_login": true, "collect_dues": true, "stockpiling": 2}
    tier_rewards_claimed_json JSONB NOT NULL DEFAULT '[]',
    -- e.g., ["bronze", "silver"]
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_player_daily UNIQUE (player_id, quest_date),
    CONSTRAINT chk_points_non_negative CHECK (daily_points >= 0 AND daily_points <= 100)
);

CREATE INDEX idx_daily_quest_player ON daily_quest_progress (player_id);
CREATE INDEX idx_daily_quest_date ON daily_quest_progress (quest_date DESC);

-- =============================================================================
-- 3.25 Seed Data: Main Quest Chain (Phase 1)
-- =============================================================================
INSERT INTO quest_types (quest_key, category, display_name, requirement_type, requirement_target, requirement_value, chain_order, reward_metal, reward_he3, reward_gold, reward_item_json, phase, is_active) VALUES
-- Phase 1 Main Quests (building/resource focus)
('main_01_collecting_resources', 'main', 'Collecting Resources', 'harvest_resources', 'resource_warehouse', 1, 1, 450, 950, 500, '[{"type":"item","item_key":"loudspeaker","quantity":1}]', 1, true),
('main_03_tech_center', 'main', 'Level 1 Technology Center', 'build_building', 'technology_center', 1, 3, 2250, 2100, 3250, '[{"type":"item","item_key":"loudspeaker","quantity":1}]', 1, true),
('main_04_research', 'main', 'Level 1 Technological Research', 'research_tech', 'concurrent_construction', 1, 4, 500, 480, 600, '[{"type":"item","item_key":"construction_card","quantity":1}]', 1, true),
('main_05_metal_production', 'main', 'Metal Production', 'build_building', 'metal_collector', 1, 5, 425, 530, 425, '[{"type":"blueprint","blueprint_key":"super_transmission_engine"}]', 1, true),
('main_06_blueprints_1', 'main', 'Blueprints 1', 'use_blueprint', 'super_transmission_engine', 1, 6, 515, 1190, 575, '[{"type":"item","item_key":"truce_card","quantity":1}]', 1, true),
('main_07_he3_production', 'main', 'He3 Production', 'build_building', 'he3_extractor', 1, 7, 475, 400, 475, '[{"type":"blueprint","blueprint_key":"estrella"}]', 1, true),
('main_08_blueprints_2', 'main', 'Blueprints 2', 'use_blueprint', 'estrella', 1, 8, 520, 1150, 585, '[{"type":"blueprint","blueprint_key":"ship_reinforcement_facility"}]', 1, true),
('main_09_residential', 'main', 'Creating Residential Area', 'build_building', 'residential_area', 1, 9, 390, 360, 325, '[{"type":"item","item_key":"loudspeaker","quantity":1}]', 1, true),
('main_23_space_station', 'main', 'Level 2 Space Station', 'upgrade_building', 'space_station', 2, 23, 10465, 9660, 13685, '[]', 1, true),
('main_24_space_defense', 'main', 'Space Defense 1', 'build_defense', '', 1, 24, 750, 1600, 880, '[]', 1, true),
('main_25_metal_lv2', 'main', 'Level 2 Metal Collector', 'upgrade_building', 'metal_collector', 2, 25, 815, 690, 815, '[{"type":"item","item_key":"metal_mining_boost","quantity":1}]', 1, true),
('main_26_he3_lv2', 'main', 'Level 2 He3 Extractor', 'upgrade_building', 'he3_extractor', 2, 26, 730, 910, 730, '[{"type":"item","item_key":"he3_mining_boost","quantity":1}]', 1, true),
('main_27_residential_lv2', 'main', 'Level 2 Residential Area', 'upgrade_building', 'residential_area', 2, 27, 670, 620, 560, '[{"type":"item","item_key":"extra_tax","quantity":1}]', 1, true),
-- Phase 2 Main Quests (military focus)
('main_10_ship_factory', 'main', 'Level 1 Ship Factory', 'build_building', 'ship_factory', 1, 10, 600, 475, 550, '[{"type":"blueprint","blueprint_key":"typhoon"}]', 2, true),
('main_11_command_center', 'main', 'Level 1 Command Center', 'build_building', 'command_center', 1, 11, 3000, 2250, 2500, '[{"type":"blueprint","blueprint_key":"energy_shield_booster"}]', 2, true),
('main_12_recruit', 'main', 'Recruit Commanders', 'recruit_commander', '', 1, 12, 605, 1250, 685, '[{"type":"item","item_key":"revival_card","quantity":1}]', 2, true),
('main_13_design_ship', 'main', 'Design a Ship', 'create_ship_design', '', 1, 13, 585, 1235, 650, '[{"type":"item","item_key":"loudspeaker","quantity":1}]', 2, true),
('main_14_build_ships', 'main', 'Ship Building', 'build_ships', '', 1, 14, 600, 1200, 680, '[{"type":"blueprint","blueprint_key":"anti_aircraft_cannon"}]', 2, true),
('main_15_fleet', 'main', 'Build a Fleet', 'create_fleet', '', 1, 15, 610, 1300, 690, '[{"type":"item","item_key":"loudspeaker","quantity":1}]', 2, true),
('main_16_logistics', 'main', 'Wartime Logistics', 'replenish_ammo', '', 1, 16, 2000, 3000, 2000, '[{"type":"item","item_key":"truce_card","quantity":1}]', 2, true),
('main_17_weapon_research', 'main', 'Level 1 Weapon Research Center', 'build_building', 'weapon_research_center', 1, 17, 2500, 1500, 2250, '[{"type":"blueprint","blueprint_key":"starlight_missile_pod"}]', 2, true),
-- Deferred Main Quests (social/chat-dependent)
('main_02_loud_and_clear', 'main', 'Loud and Clear', 'send_message', 'world_channel', 1, 2, 460, 980, 520, '[]', 3, false),
('main_18_bigger_bags', 'main', 'Bigger Bags', 'increase_bag_slot', '', 1, 18, 620, 1310, 700, '[{"type":"item","item_key":"primary_metal_pack","quantity":1}]', 3, false),
('main_19_resource_pack', 'main', 'Resource Pack', 'use_resource_pack', '', 1, 19, 680, 1450, 750, '[{"type":"item","item_key":"galaxy_transfer","quantity":1}]', 3, false),
('main_20_growing_resources', 'main', 'Growing Resources', 'grow_comsats', '', 1, 20, 700, 1480, 780, '[{"type":"item","item_key":"healing_card","quantity":1}]', 3, false),
('main_21_adding_friends', 'main', 'Adding Friends', 'add_friend', '', 1, 21, 710, 1500, 800, '[]', 3, false),
('main_22_mail_system', 'main', 'The Mail System', 'send_message', 'mail', 1, 22, 720, 1520, 820, '[{"type":"item","item_key":"loudspeaker","quantity":1}]', 3, false),
('main_28_ship_factory_lv2', 'main', 'Level 2 Ship Factory', 'upgrade_building', 'ship_factory', 2, 28, 0, 0, 0, '[]', 2, true);

-- Set prerequisite chain (main_01 -> main_03 -> main_04 -> ... sequential)
-- Note: Skip deferred quests in Phase 1/2 chain by linking past them
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_01_collecting_resources') WHERE quest_key = 'main_03_tech_center';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_03_tech_center') WHERE quest_key = 'main_04_research';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_04_research') WHERE quest_key = 'main_05_metal_production';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_05_metal_production') WHERE quest_key = 'main_06_blueprints_1';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_06_blueprints_1') WHERE quest_key = 'main_07_he3_production';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_07_he3_production') WHERE quest_key = 'main_08_blueprints_2';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_08_blueprints_2') WHERE quest_key = 'main_09_residential';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_09_residential') WHERE quest_key = 'main_10_ship_factory';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_10_ship_factory') WHERE quest_key = 'main_11_command_center';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_11_command_center') WHERE quest_key = 'main_12_recruit';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_12_recruit') WHERE quest_key = 'main_13_design_ship';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_13_design_ship') WHERE quest_key = 'main_14_build_ships';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_14_build_ships') WHERE quest_key = 'main_15_fleet';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_15_fleet') WHERE quest_key = 'main_16_logistics';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_16_logistics') WHERE quest_key = 'main_17_weapon_research';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_17_weapon_research') WHERE quest_key = 'main_23_space_station';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_23_space_station') WHERE quest_key = 'main_24_space_defense';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_24_space_defense') WHERE quest_key = 'main_25_metal_lv2';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_25_metal_lv2') WHERE quest_key = 'main_26_he3_lv2';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_26_he3_lv2') WHERE quest_key = 'main_27_residential_lv2';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_27_residential_lv2') WHERE quest_key = 'main_28_ship_factory_lv2';

-- =============================================================================
-- 3.26 Seed Data: Side Quests (Phase 1)
-- =============================================================================
INSERT INTO quest_types (quest_key, category, display_name, requirement_type, requirement_target, requirement_value, chain_order, reward_metal, reward_he3, reward_gold, reward_item_json, phase, is_active) VALUES
('side_harvest_time_1', 'side', 'Harvest Time I', 'reach_production', 'metal', 2180, 1, 500, 400, 500, '[]', 1, true),
('side_harvest_time_2', 'side', 'Harvest Time II', 'reach_production', 'metal', 4360, 2, 1500, 1200, 1500, '[]', 1, true),
('side_harvest_time_3', 'side', 'Harvest Time III', 'reach_production', 'metal', 8720, 3, 4500, 3600, 4500, '[]', 1, true),
('side_gathering_he3_1', 'side', 'Gathering He3 I', 'reach_production', 'he3', 2360, 1, 400, 500, 400, '[]', 1, true),
('side_gathering_he3_2', 'side', 'Gathering He3 II', 'reach_production', 'he3', 4720, 2, 1200, 1500, 1200, '[]', 1, true),
('side_gathering_he3_3', 'side', 'Gathering He3 III', 'reach_production', 'he3', 9440, 3, 3600, 4500, 3600, '[]', 1, true),
('side_raising_morale_1', 'side', 'Raising Morale I', 'reach_production', 'gold', 2800, 1, 400, 400, 500, '[]', 1, true),
('side_raising_morale_2', 'side', 'Raising Morale II', 'reach_production', 'gold', 5600, 2, 1200, 1200, 1500, '[]', 1, true),
('side_raising_morale_3', 'side', 'Raising Morale III', 'reach_production', 'gold', 11200, 3, 3600, 3600, 4500, '[]', 1, true),
('side_plentiful_resources_1', 'side', 'Plentiful Resources I', 'reach_storage', '', 50000, 1, 1000, 1000, 1000, '[]', 1, true),
('side_plentiful_resources_2', 'side', 'Plentiful Resources II', 'reach_storage', '', 200000, 2, 3000, 3000, 3000, '[]', 1, true),
('side_plentiful_resources_3', 'side', 'Plentiful Resources III', 'reach_storage', '', 1000000, 3, 10000, 10000, 10000, '[]', 1, true);

-- Set side quest tier chains
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'side_harvest_time_1') WHERE quest_key = 'side_harvest_time_2';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'side_harvest_time_2') WHERE quest_key = 'side_harvest_time_3';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'side_gathering_he3_1') WHERE quest_key = 'side_gathering_he3_2';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'side_gathering_he3_2') WHERE quest_key = 'side_gathering_he3_3';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'side_raising_morale_1') WHERE quest_key = 'side_raising_morale_2';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'side_raising_morale_2') WHERE quest_key = 'side_raising_morale_3';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'side_plentiful_resources_1') WHERE quest_key = 'side_plentiful_resources_2';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'side_plentiful_resources_2') WHERE quest_key = 'side_plentiful_resources_3';

-- =============================================================================
-- 3.27 Seed Data: Daily Quests
-- =============================================================================
INSERT INTO quest_types (quest_key, category, display_name, requirement_type, requirement_target, requirement_value, chain_order, reward_metal, reward_he3, reward_gold, reward_item_json, phase, is_active) VALUES
('daily_login', 'daily', 'Daily Log In', 'login', '', 1, NULL, 0, 0, 0, '[]', 1, true),
('daily_collect_dues', 'daily', 'Collect Your Dues', 'harvest_resources', '', 1, NULL, 0, 0, 0, '[]', 1, true),
('daily_need_for_speed', 'daily', 'Need for Speed', 'use_speedup', 'construction', 1, NULL, 0, 0, 0, '[]', 1, true),
('daily_stockpiling', 'daily', 'Stockpiling', 'harvest_resources', 'resource_warehouse', 3, NULL, 0, 0, 0, '[]', 1, true),
('daily_donations', 'daily', 'Donations', 'donate_resources', '', 200000, NULL, 0, 0, 0, '[]', 2, true),
('daily_restricted_instances', 'daily', 'Restricted Instances', 'complete_instance', 'restricted', 2, NULL, 0, 0, 5000, '[]', 2, true);
