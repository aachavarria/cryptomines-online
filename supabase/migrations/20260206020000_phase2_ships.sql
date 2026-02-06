-- ============================================================================
-- Phase 2 Migration - Ships, Fleets & Combat
-- New tables: hull_types, module_types, ship_factory_levels, blueprints,
--   player_blueprints, ship_designs, ships, commanders, combat_reports,
--   fleets, fleet_stacks, instances, instance_progress, instance_blueprints,
--   spacedock_repairs, blueprint_research
-- Altered: spacedock_levels (repair_pct INTEGER -> NUMERIC)
-- ============================================================================

-- ============================================================================
-- PART 1: REFERENCE TABLES
-- ============================================================================

-- =========================
-- 1. HULL TYPES
-- =========================
CREATE TABLE hull_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    display_name TEXT NOT NULL,
    hull_class TEXT NOT NULL CHECK (hull_class IN ('frigate', 'cruiser', 'battleship')),
    tier INTEGER NOT NULL CHECK (tier BETWEEN 1 AND 3),
    armor_type TEXT NOT NULL CHECK (armor_type IN ('nano', 'chrome', 'regen', 'neutralizing')),
    base_shield INTEGER NOT NULL DEFAULT 0,
    base_structure INTEGER NOT NULL DEFAULT 0,
    base_stability NUMERIC(5,2) NOT NULL DEFAULT 100.00,
    base_defense NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    installation_slots INTEGER NOT NULL DEFAULT 100,
    base_agility INTEGER NOT NULL DEFAULT 0,
    base_movement INTEGER NOT NULL DEFAULT 0,
    base_storage INTEGER NOT NULL DEFAULT 0,
    base_metal_cost BIGINT NOT NULL DEFAULT 0,
    base_he3_cost BIGINT NOT NULL DEFAULT 0,
    base_gold_cost BIGINT NOT NULL DEFAULT 0,
    base_build_time_seconds INTEGER NOT NULL DEFAULT 10,
    description TEXT NOT NULL DEFAULT ''
);

CREATE INDEX idx_hull_types_class ON hull_types (hull_class);
CREATE INDEX idx_hull_types_tier ON hull_types (tier);

-- =========================
-- 2. MODULE TYPES
-- =========================
CREATE TABLE module_types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    display_name TEXT NOT NULL,
    category TEXT NOT NULL CHECK (category IN (
        'ballistic', 'directional', 'missile', 'ship_based', 'planetary',
        'structure', 'shield', 'air_defense',
        'electronic', 'storage', 'transmission'
    )),
    tier INTEGER NOT NULL DEFAULT 1 CHECK (tier BETWEEN 0 AND 3),
    damage_type TEXT CHECK (damage_type IN ('kinetic', 'heat', 'explosive', 'magnetic', 'siege', NULL)),
    min_damage INTEGER NOT NULL DEFAULT 0,
    max_damage INTEGER NOT NULL DEFAULT 0,
    weapon_range_min INTEGER NOT NULL DEFAULT 0,
    weapon_range_max INTEGER NOT NULL DEFAULT 0,
    cooldown INTEGER NOT NULL DEFAULT 0,
    he3_per_round INTEGER NOT NULL DEFAULT 0,
    volume INTEGER NOT NULL DEFAULT 1,
    max_per_ship INTEGER NOT NULL DEFAULT 0,
    effects_json JSONB NOT NULL DEFAULT '{}',
    metal_cost BIGINT NOT NULL DEFAULT 0,
    he3_cost BIGINT NOT NULL DEFAULT 0,
    gold_cost BIGINT NOT NULL DEFAULT 0,
    build_time_seconds INTEGER NOT NULL DEFAULT 1,
    description TEXT NOT NULL DEFAULT '',

    CONSTRAINT uq_module_name_tier UNIQUE (name, tier),
    CONSTRAINT chk_damage_range CHECK (min_damage <= max_damage),
    CONSTRAINT chk_weapon_range CHECK (weapon_range_min <= weapon_range_max)
);

CREATE INDEX idx_module_types_category ON module_types (category);
CREATE INDEX idx_module_types_tier ON module_types (tier);

-- =========================
-- 3. SHIP FACTORY LEVELS
-- =========================
CREATE TABLE ship_factory_levels (
    level INTEGER PRIMARY KEY,
    civic_center_req INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL,
    speed_bonus_pct INTEGER NOT NULL,
    production_slots INTEGER NOT NULL
);

-- =========================
-- 4. ALTER SPACEDOCK LEVELS (repair_pct INTEGER -> NUMERIC)
-- =========================
ALTER TABLE spacedock_levels ALTER COLUMN repair_pct TYPE NUMERIC(5,2);

-- Update spacedock_levels with more precise repair percentages from GDD 8.4.3
UPDATE spacedock_levels SET repair_pct = 1.00 WHERE level = 1;
UPDATE spacedock_levels SET repair_pct = 2.70 WHERE level = 2;
UPDATE spacedock_levels SET repair_pct = 4.50 WHERE level = 3;
UPDATE spacedock_levels SET repair_pct = 6.30 WHERE level = 4;
UPDATE spacedock_levels SET repair_pct = 8.10 WHERE level = 5;
UPDATE spacedock_levels SET repair_pct = 10.00 WHERE level = 6;
UPDATE spacedock_levels SET repair_pct = 11.80 WHERE level = 7;
UPDATE spacedock_levels SET repair_pct = 13.60 WHERE level = 8;
UPDATE spacedock_levels SET repair_pct = 15.50 WHERE level = 9;
UPDATE spacedock_levels SET repair_pct = 17.30 WHERE level = 10;
UPDATE spacedock_levels SET repair_pct = 18.60 WHERE level = 11;
UPDATE spacedock_levels SET repair_pct = 20.00 WHERE level = 12;

-- Also update the costs to match GDD 8.4.3 (differs from Phase 1 values)
DELETE FROM spacedock_levels;
INSERT INTO spacedock_levels VALUES
(1,  1,  1,  500,        400,        450,        44,        1.00),
(2,  2,  3,  1515,       1212,       1364,       126,       2.70),
(3,  3,  5,  4590,       3672,       4131,       362,       4.50),
(4,  4,  7,  13910,      11128,      12519,      1039,      6.30),
(5,  5,  9,  42147,      33718,      37932,      2982,      8.10),
(6,  6,  11, 127704,     102163,     114933,     8559,      10.00),
(7,  7,  13, 386942,     309554,     348248,     24565,     11.80),
(8,  8,  15, 1172434,    937947,     1055190,    70501,     13.60),
(9,  9,  17, 3554140,    2843312,    3198726,    202339,    15.50),
(10, 10, 19, 10770711,   8616569,    9693640,    580714,    17.30),
(11, 11, 21, 32641855,   26113484,   29377729,   1666651,   18.60),
(12, 12, 23, 98921822,   79137458,   89029519,   10036800,  20.00);

-- ============================================================================
-- PART 2: GAME ENTITY TABLES
-- ============================================================================

-- =========================
-- 5. BLUEPRINTS
-- =========================
CREATE TABLE blueprints (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    blueprint_type TEXT NOT NULL CHECK (blueprint_type IN ('hull', 'module')),
    hull_type_id INTEGER REFERENCES hull_types(id) ON DELETE SET NULL,
    module_type_id INTEGER REFERENCES module_types(id) ON DELETE SET NULL,
    source TEXT NOT NULL DEFAULT 'instance'
        CHECK (source IN ('instance', 'quest', 'auction', 'trafficker', 'mall')),
    research_level INTEGER NOT NULL DEFAULT 0 CHECK (research_level BETWEEN 0 AND 3),
    description TEXT NOT NULL DEFAULT '',

    CONSTRAINT chk_blueprint_ref CHECK (
        (blueprint_type = 'hull' AND hull_type_id IS NOT NULL AND module_type_id IS NULL) OR
        (blueprint_type = 'module' AND module_type_id IS NOT NULL AND hull_type_id IS NULL)
    )
);

CREATE INDEX idx_blueprints_type ON blueprints (blueprint_type);
CREATE INDEX idx_blueprints_hull ON blueprints (hull_type_id) WHERE hull_type_id IS NOT NULL;
CREATE INDEX idx_blueprints_module ON blueprints (module_type_id) WHERE module_type_id IS NOT NULL;

-- =========================
-- 6. PLAYER BLUEPRINTS
-- =========================
CREATE TABLE player_blueprints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    blueprint_id INTEGER NOT NULL REFERENCES blueprints(id) ON DELETE CASCADE,
    is_activated BOOLEAN NOT NULL DEFAULT false,
    research_level INTEGER NOT NULL DEFAULT 1 CHECK (research_level BETWEEN 1 AND 3),
    acquired_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_player_blueprint UNIQUE (player_id, blueprint_id)
);

CREATE INDEX idx_player_blueprints_player ON player_blueprints (player_id);
CREATE INDEX idx_player_blueprints_activated ON player_blueprints (is_activated) WHERE is_activated = true;

-- =========================
-- 7. SHIP DESIGNS
-- =========================
CREATE TABLE ship_designs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    hull_type_id INTEGER NOT NULL REFERENCES hull_types(id),
    modules_json JSONB NOT NULL DEFAULT '[]',
    total_shield INTEGER NOT NULL DEFAULT 0,
    total_structure INTEGER NOT NULL DEFAULT 0,
    total_defense NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    total_agility INTEGER NOT NULL DEFAULT 0,
    total_movement INTEGER NOT NULL DEFAULT 0,
    total_storage INTEGER NOT NULL DEFAULT 0,
    attack_power INTEGER NOT NULL DEFAULT 0,
    weapon_range_min INTEGER NOT NULL DEFAULT 0,
    weapon_range_max INTEGER NOT NULL DEFAULT 0,
    volume_used INTEGER NOT NULL DEFAULT 0,
    he3_per_round INTEGER NOT NULL DEFAULT 0,
    metal_cost BIGINT NOT NULL DEFAULT 0,
    he3_cost BIGINT NOT NULL DEFAULT 0,
    gold_cost BIGINT NOT NULL DEFAULT 0,
    build_time_seconds INTEGER NOT NULL DEFAULT 10,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_design_name CHECK (name ~ '^[a-zA-Z0-9._-]+$'),
    CONSTRAINT chk_volume CHECK (volume_used >= 0)
);

CREATE INDEX idx_ship_designs_player ON ship_designs (player_id);
CREATE INDEX idx_ship_designs_hull ON ship_designs (hull_type_id);

-- =========================
-- 8. SHIPS
-- =========================
CREATE TABLE ships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    ship_design_id UUID NOT NULL REFERENCES ship_designs(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 0,
    is_building BOOLEAN NOT NULL DEFAULT false,
    build_quantity INTEGER NOT NULL DEFAULT 0,
    build_finish_at TIMESTAMPTZ,
    production_slot INTEGER NOT NULL DEFAULT 1 CHECK (production_slot BETWEEN 1 AND 5),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_player_ship_design UNIQUE (player_id, ship_design_id),
    CONSTRAINT chk_quantity_non_negative CHECK (quantity >= 0),
    CONSTRAINT chk_build_consistency CHECK (
        (is_building = true AND build_finish_at IS NOT NULL AND build_quantity > 0) OR
        (is_building = false AND build_finish_at IS NULL AND build_quantity = 0)
    )
);

CREATE INDEX idx_ships_player ON ships (player_id);
CREATE INDEX idx_ships_building ON ships (is_building) WHERE is_building = true;

-- =========================
-- 9. COMMANDERS
-- =========================
CREATE TABLE commanders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    rarity TEXT NOT NULL CHECK (rarity IN ('common', 'skill', 'super', 'legendary', 'divine')),
    star_rank INTEGER NOT NULL DEFAULT 0,
    accuracy INTEGER NOT NULL DEFAULT 0,
    dodge INTEGER NOT NULL DEFAULT 0,
    speed INTEGER NOT NULL DEFAULT 0,
    electron INTEGER NOT NULL DEFAULT 0,
    weapon_expertise JSONB NOT NULL DEFAULT '{"ballistic":"B","directional":"B","missile":"B","ship_based":"B"}',
    ship_expertise JSONB NOT NULL DEFAULT '{"frigate":"B","cruiser":"B","battleship":"B"}',
    skills_json JSONB NOT NULL DEFAULT '[]',
    gems_json JSONB NOT NULL DEFAULT '[]',
    bionic_chips_json JSONB NOT NULL DEFAULT '[]',
    is_deployed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_star_rank CHECK (star_rank >= 0 AND star_rank <= 15),
    CONSTRAINT chk_stats_non_negative CHECK (
        accuracy >= 0 AND dodge >= 0 AND speed >= 0 AND electron >= 0
    )
);

CREATE INDEX idx_commanders_player_id ON commanders (player_id);
CREATE INDEX idx_commanders_rarity ON commanders (rarity);

-- =========================
-- 10. COMBAT REPORTS
-- =========================
CREATE TABLE combat_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attacker_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    defender_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    attacker_fleet_id UUID,
    defender_fleet_id UUID,
    combat_type TEXT NOT NULL DEFAULT 'pvp'
        CHECK (combat_type IN ('pvp', 'instance_normal', 'instance_restricted', 'instance_trial', 'instance_constellation', 'instance_humaroid', 'league', 'championship', 'rbp_attack')),
    result TEXT NOT NULL CHECK (result IN ('attacker_win', 'defender_win', 'draw')),
    total_rounds INTEGER NOT NULL DEFAULT 0,
    rounds_json JSONB NOT NULL DEFAULT '[]',
    loot_json JSONB NOT NULL DEFAULT '{}',
    attacker_losses_json JSONB NOT NULL DEFAULT '{}',
    defender_losses_json JSONB NOT NULL DEFAULT '{}',
    he3_consumed BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_combat_reports_attacker ON combat_reports (attacker_id);
CREATE INDEX idx_combat_reports_defender ON combat_reports (defender_id);
CREATE INDEX idx_combat_reports_created ON combat_reports (created_at DESC);
CREATE INDEX idx_combat_reports_type ON combat_reports (combat_type);

-- =========================
-- 11. FLEETS
-- =========================
CREATE TABLE fleets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    name TEXT NOT NULL DEFAULT 'Fleet',
    formation TEXT NOT NULL DEFAULT 'phalanx'
        CHECK (formation IN ('phalanx', 'diamond', 'battle_line', 'skirmish', 'tee_forward', 'enfilade', 'tee_reverse')),
    commander_id UUID REFERENCES commanders(id) ON DELETE SET NULL,
    targeting_command TEXT NOT NULL DEFAULT 'max_attack'
        CHECK (targeting_command IN ('max_attack', 'min_attack', 'max_durability', 'min_durability', 'closest', 'by_commander_rank')),
    status TEXT NOT NULL DEFAULT 'stationed'
        CHECK (status IN ('stationed', 'traveling', 'combat', 'returning', 'dismissed')),
    planet_id UUID REFERENCES planets(id) ON DELETE SET NULL,
    position_x INTEGER,
    position_y INTEGER,
    destination_x INTEGER,
    destination_y INTEGER,
    arrival_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_travel_consistency CHECK (
        (status IN ('traveling', 'returning') AND destination_x IS NOT NULL
            AND destination_y IS NOT NULL AND arrival_at IS NOT NULL) OR
        (status IN ('stationed', 'combat', 'dismissed'))
    )
);

CREATE INDEX idx_fleets_player ON fleets (player_id);
CREATE INDEX idx_fleets_status ON fleets (status);
CREATE INDEX idx_fleets_commander ON fleets (commander_id) WHERE commander_id IS NOT NULL;

-- Add fleet_id foreign keys to combat_reports now that fleets table exists
ALTER TABLE combat_reports
    ADD CONSTRAINT fk_combat_reports_attacker_fleet
    FOREIGN KEY (attacker_fleet_id) REFERENCES fleets(id) ON DELETE SET NULL;
ALTER TABLE combat_reports
    ADD CONSTRAINT fk_combat_reports_defender_fleet
    FOREIGN KEY (defender_fleet_id) REFERENCES fleets(id) ON DELETE SET NULL;

-- =========================
-- 12. FLEET STACKS
-- =========================
CREATE TABLE fleet_stacks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id UUID NOT NULL REFERENCES fleets(id) ON DELETE CASCADE,
    ship_design_id UUID NOT NULL REFERENCES ship_designs(id),
    grid_row INTEGER NOT NULL CHECK (grid_row BETWEEN 0 AND 2),
    grid_col INTEGER NOT NULL CHECK (grid_col BETWEEN 0 AND 2),
    ship_count INTEGER NOT NULL DEFAULT 0 CHECK (ship_count BETWEEN 0 AND 3000),

    CONSTRAINT uq_fleet_position UNIQUE (fleet_id, grid_row, grid_col),
    CONSTRAINT chk_ship_count CHECK (ship_count >= 0 AND ship_count <= 3000)
);

CREATE INDEX idx_fleet_stacks_fleet ON fleet_stacks (fleet_id);
CREATE INDEX idx_fleet_stacks_design ON fleet_stacks (ship_design_id);

-- =========================
-- 13. INSTANCES
-- =========================
CREATE TABLE instances (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('normal', 'restricted', 'trial', 'constellation', 'humaroid')),
    difficulty INTEGER NOT NULL DEFAULT 1,
    required_level INTEGER NOT NULL DEFAULT 1,
    max_fleets INTEGER NOT NULL DEFAULT 3,
    ships_lost_on_defeat BOOLEAN NOT NULL DEFAULT true,
    he3_lost_on_defeat BOOLEAN NOT NULL DEFAULT true,
    exp_reward INTEGER NOT NULL DEFAULT 0,
    enemy_fleets_json JSONB NOT NULL DEFAULT '[]',
    rewards_json JSONB NOT NULL DEFAULT '{}',
    description TEXT NOT NULL DEFAULT '',

    CONSTRAINT chk_difficulty_positive CHECK (difficulty >= 1),
    CONSTRAINT chk_level_positive CHECK (required_level >= 1)
);

CREATE INDEX idx_instances_type ON instances (type);
CREATE INDEX idx_instances_difficulty ON instances (difficulty);

-- =========================
-- 14. INSTANCE PROGRESS
-- =========================
CREATE TABLE instance_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    instance_id INTEGER NOT NULL REFERENCES instances(id) ON DELETE CASCADE,
    completed BOOLEAN NOT NULL DEFAULT false,
    attempts INTEGER NOT NULL DEFAULT 0,
    best_score INTEGER NOT NULL DEFAULT 0,
    last_attempt_at TIMESTAMPTZ,

    CONSTRAINT uq_player_instance UNIQUE (player_id, instance_id),
    CONSTRAINT chk_attempts_non_negative CHECK (attempts >= 0),
    CONSTRAINT chk_score_non_negative CHECK (best_score >= 0)
);

CREATE INDEX idx_instance_progress_player ON instance_progress (player_id);

-- =========================
-- 15. INSTANCE BLUEPRINTS (Junction Table)
-- =========================
CREATE TABLE instance_blueprints (
    instance_id INTEGER NOT NULL REFERENCES instances(id) ON DELETE CASCADE,
    blueprint_id INTEGER NOT NULL REFERENCES blueprints(id) ON DELETE CASCADE,

    PRIMARY KEY (instance_id, blueprint_id)
);

-- =========================
-- 16. SPACEDOCK REPAIRS
-- =========================
CREATE TABLE spacedock_repairs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    ship_design_id UUID NOT NULL REFERENCES ship_designs(id),
    destroyed_count INTEGER NOT NULL DEFAULT 0,
    repaired_count INTEGER NOT NULL DEFAULT 0,
    repair_finish_at TIMESTAMPTZ,
    combat_report_id UUID REFERENCES combat_reports(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_counts CHECK (destroyed_count >= 0 AND repaired_count >= 0 AND repaired_count <= destroyed_count)
);

CREATE INDEX idx_spacedock_repairs_player ON spacedock_repairs (player_id);
CREATE INDEX idx_spacedock_repairs_repairing ON spacedock_repairs (repair_finish_at) WHERE repair_finish_at IS NOT NULL;

-- =========================
-- 17. BLUEPRINT RESEARCH
-- =========================
CREATE TABLE blueprint_research (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    player_blueprint_id UUID NOT NULL REFERENCES player_blueprints(id) ON DELETE CASCADE,
    target_level INTEGER NOT NULL CHECK (target_level BETWEEN 2 AND 3),
    is_researching BOOLEAN NOT NULL DEFAULT false,
    research_finish_at TIMESTAMPTZ,
    metal_cost BIGINT NOT NULL DEFAULT 0,
    he3_cost BIGINT NOT NULL DEFAULT 0,
    gold_cost BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_research_consistency CHECK (
        (is_researching = true AND research_finish_at IS NOT NULL) OR
        (is_researching = false AND research_finish_at IS NULL)
    )
);

CREATE INDEX idx_blueprint_research_player ON blueprint_research (player_id);
CREATE INDEX idx_blueprint_research_active ON blueprint_research (is_researching) WHERE is_researching = true;

-- ============================================================================
-- PART 3: TRIGGERS
-- ============================================================================

CREATE TRIGGER tr_ship_designs_updated_at
    BEFORE UPDATE ON ship_designs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER tr_ships_updated_at
    BEFORE UPDATE ON ships
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER tr_commanders_updated_at
    BEFORE UPDATE ON commanders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER tr_fleets_updated_at
    BEFORE UPDATE ON fleets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- PART 4: ROW LEVEL SECURITY
-- ============================================================================

-- Reference tables (public read)
ALTER TABLE hull_types ENABLE ROW LEVEL SECURITY;
ALTER TABLE module_types ENABLE ROW LEVEL SECURITY;
ALTER TABLE ship_factory_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE blueprints ENABLE ROW LEVEL SECURITY;
ALTER TABLE instances ENABLE ROW LEVEL SECURITY;

CREATE POLICY hull_types_select_all ON hull_types FOR SELECT USING (true);
CREATE POLICY module_types_select_all ON module_types FOR SELECT USING (true);
CREATE POLICY ship_factory_levels_select_all ON ship_factory_levels FOR SELECT USING (true);
CREATE POLICY blueprints_select_all ON blueprints FOR SELECT USING (true);
CREATE POLICY instances_select_all ON instances FOR SELECT USING (true);

-- Player-owned tables
ALTER TABLE player_blueprints ENABLE ROW LEVEL SECURITY;
ALTER TABLE ship_designs ENABLE ROW LEVEL SECURITY;
ALTER TABLE ships ENABLE ROW LEVEL SECURITY;
ALTER TABLE commanders ENABLE ROW LEVEL SECURITY;
ALTER TABLE combat_reports ENABLE ROW LEVEL SECURITY;
ALTER TABLE fleets ENABLE ROW LEVEL SECURITY;
ALTER TABLE fleet_stacks ENABLE ROW LEVEL SECURITY;
ALTER TABLE instance_progress ENABLE ROW LEVEL SECURITY;
ALTER TABLE instance_blueprints ENABLE ROW LEVEL SECURITY;
ALTER TABLE spacedock_repairs ENABLE ROW LEVEL SECURITY;
ALTER TABLE blueprint_research ENABLE ROW LEVEL SECURITY;

CREATE POLICY player_blueprints_select_own ON player_blueprints
    FOR SELECT USING (player_id = auth.uid());

CREATE POLICY ship_designs_select_own ON ship_designs
    FOR SELECT USING (player_id = auth.uid());

CREATE POLICY ships_select_own ON ships
    FOR SELECT USING (player_id = auth.uid());

CREATE POLICY commanders_select_own ON commanders
    FOR SELECT USING (player_id = auth.uid());

CREATE POLICY combat_reports_select_own ON combat_reports
    FOR SELECT USING (attacker_id = auth.uid() OR defender_id = auth.uid());

CREATE POLICY fleets_select_own ON fleets
    FOR SELECT USING (player_id = auth.uid());

CREATE POLICY fleet_stacks_select_own ON fleet_stacks
    FOR SELECT USING (
        fleet_id IN (SELECT id FROM fleets WHERE player_id = auth.uid())
    );

CREATE POLICY instance_progress_select_own ON instance_progress
    FOR SELECT USING (player_id = auth.uid());

CREATE POLICY instance_blueprints_select_all ON instance_blueprints
    FOR SELECT USING (true);

CREATE POLICY spacedock_repairs_select_own ON spacedock_repairs
    FOR SELECT USING (player_id = auth.uid());

CREATE POLICY blueprint_research_select_own ON blueprint_research
    FOR SELECT USING (player_id = auth.uid());

-- ============================================================================
-- PART 5: SEED DATA - HULL TYPES (25 hull lines x 3 tiers = 75 rows)
-- ============================================================================

-- Frigate Hull Lines (10 lines x 3 tiers = 30)
-- Stats from GDD 8.2.2: Tier I base, Tier II ~2x, Tier III ~4x
-- All frigates: base_movement=1, base_agility=1, base_stability=100
INSERT INTO hull_types (name, display_name, hull_class, tier, armor_type, base_shield, base_structure, base_stability, base_defense, installation_slots, base_agility, base_movement, base_storage, base_metal_cost, base_he3_cost, base_gold_cost, base_build_time_seconds, description) VALUES
-- Weikes
('weikes_i',       'Weikes I',       'frigate', 1, 'nano',          270,  770,  100.00, 0.00, 80,  1, 1, 0, 500,  300,  200,  10, 'Starter frigate with balanced stats'),
('weikes_ii',      'Weikes II',      'frigate', 2, 'nano',          540,  1540, 100.00, 0.00, 120, 1, 1, 0, 1500, 900,  600,  20, 'Improved Weikes hull'),
('weikes_iii',     'Weikes III',     'frigate', 3, 'nano',          1078, 3080, 100.00, 0.00, 180, 1, 1, 0, 4500, 2700, 1800, 40, 'Advanced Weikes hull'),
-- Air Wanderer
('air_wanderer_i', 'Air Wanderer I', 'frigate', 1, 'neutralizing',  290,  810,  100.00, 0.00, 85,  1, 1, 0, 520,  310,  210,  10, 'Neutralizing armor frigate'),
('air_wanderer_ii','Air Wanderer II','frigate', 2, 'neutralizing',  580,  1620, 100.00, 0.00, 125, 1, 1, 0, 1560, 930,  630,  20, 'Improved Air Wanderer hull'),
('air_wanderer_iii','Air Wanderer III','frigate',3, 'neutralizing', 1150, 3240, 100.00, 0.00, 190, 1, 1, 0, 4680, 2790, 1890, 40, 'Advanced Air Wanderer hull'),
-- Valkyrie
('valkyrie_i',     'Valkyrie I',     'frigate', 1, 'nano',          310,  750,  100.00, 0.00, 90,  1, 1, 0, 540,  320,  220,  10, 'High-shield nano frigate'),
('valkyrie_ii',    'Valkyrie II',    'frigate', 2, 'nano',          620,  1500, 100.00, 0.00, 130, 1, 1, 0, 1620, 960,  660,  20, 'Improved Valkyrie hull'),
('valkyrie_iii',   'Valkyrie III',   'frigate', 3, 'nano',          1230, 3000, 100.00, 0.00, 195, 1, 1, 0, 4860, 2880, 1980, 40, 'Advanced Valkyrie hull'),
-- GoGetter
('gogetter_i',     'GoGetter I',     'frigate', 1, 'neutralizing',  260,  800,  100.00, 0.00, 82,  1, 1, 0, 510,  305,  205,  10, 'Sturdy neutralizing frigate'),
('gogetter_ii',    'GoGetter II',    'frigate', 2, 'neutralizing',  520,  1600, 100.00, 0.00, 122, 1, 1, 0, 1530, 915,  615,  20, 'Improved GoGetter hull'),
('gogetter_iii',   'GoGetter III',   'frigate', 3, 'neutralizing',  1040, 3200, 100.00, 0.00, 182, 1, 1, 0, 4590, 2745, 1845, 40, 'Advanced GoGetter hull'),
-- Space Hunter
('space_hunter_i', 'Space Hunter I', 'frigate', 1, 'nano',          300,  790,  100.00, 0.00, 88,  1, 1, 0, 530,  315,  215,  10, 'Balanced nano frigate'),
('space_hunter_ii','Space Hunter II','frigate', 2, 'nano',          600,  1580, 100.00, 0.00, 128, 1, 1, 0, 1590, 945,  645,  20, 'Improved Space Hunter hull'),
('space_hunter_iii','Space Hunter III','frigate',3, 'nano',         1200, 3160, 100.00, 0.00, 188, 1, 1, 0, 4770, 2835, 1935, 40, 'Advanced Space Hunter hull'),
-- Sparrow
('sparrow_i',      'Sparrow I',      'frigate', 1, 'neutralizing',  280,  760,  100.00, 0.00, 78,  1, 1, 0, 490,  295,  195,  10, 'Light neutralizing frigate'),
('sparrow_ii',     'Sparrow II',     'frigate', 2, 'neutralizing',  560,  1520, 100.00, 0.00, 118, 1, 1, 0, 1470, 885,  585,  20, 'Improved Sparrow hull'),
('sparrow_iii',    'Sparrow III',    'frigate', 3, 'neutralizing',  1120, 3040, 100.00, 0.00, 178, 1, 1, 0, 4410, 2655, 1755, 40, 'Advanced Sparrow hull'),
-- Devourer
('devourer_i',     'Devourer I',     'frigate', 1, 'nano',          320,  740,  100.00, 0.00, 92,  1, 1, 0, 550,  330,  225,  10, 'High-shield low-structure frigate'),
('devourer_ii',    'Devourer II',    'frigate', 2, 'nano',          640,  1480, 100.00, 0.00, 132, 1, 1, 0, 1650, 990,  675,  20, 'Improved Devourer hull'),
('devourer_iii',   'Devourer III',   'frigate', 3, 'nano',          1280, 2960, 100.00, 0.00, 198, 1, 1, 0, 4950, 2970, 2025, 40, 'Advanced Devourer hull'),
-- Polymesus
('polymesus_i',    'Polymesus I',    'frigate', 1, 'neutralizing',  250,  820,  100.00, 0.00, 84,  1, 1, 0, 505,  300,  200,  10, 'Tank neutralizing frigate'),
('polymesus_ii',   'Polymesus II',   'frigate', 2, 'neutralizing',  500,  1640, 100.00, 0.00, 124, 1, 1, 0, 1515, 900,  600,  20, 'Improved Polymesus hull'),
('polymesus_iii',  'Polymesus III',  'frigate', 3, 'neutralizing',  1000, 3280, 100.00, 0.00, 184, 1, 1, 0, 4545, 2700, 1800, 40, 'Advanced Polymesus hull'),
-- Cybra
('cybra_i',        'Cybra I',        'frigate', 1, 'nano',          340,  780,  100.00, 0.00, 95,  1, 1, 0, 560,  335,  230,  10, 'Highest-shield frigate'),
('cybra_ii',       'Cybra II',       'frigate', 2, 'nano',          680,  1560, 100.00, 0.00, 135, 1, 1, 0, 1680, 1005, 690,  20, 'Improved Cybra hull'),
('cybra_iii',      'Cybra III',      'frigate', 3, 'nano',          1360, 3120, 100.00, 0.00, 200, 1, 1, 0, 5040, 3015, 2070, 40, 'Advanced Cybra hull'),
-- Hamdar
('hamdar_i',       'Hamdar I',       'frigate', 1, 'neutralizing',  275,  830,  100.00, 0.00, 86,  1, 1, 0, 515,  308,  208,  10, 'Highest-structure frigate'),
('hamdar_ii',      'Hamdar II',      'frigate', 2, 'neutralizing',  550,  1660, 100.00, 0.00, 126, 1, 1, 0, 1545, 924,  624,  20, 'Improved Hamdar hull'),
('hamdar_iii',     'Hamdar III',     'frigate', 3, 'neutralizing',  1100, 3320, 100.00, 0.00, 186, 1, 1, 0, 4635, 2772, 1872, 40, 'Advanced Hamdar hull'),

-- Cruiser Hull Lines (10 lines x 3 tiers = 30)
-- All cruisers: base_movement=0, base_agility=0, base_stability=100
-- Typhoon
('typhoon_i',      'Typhoon I',      'cruiser', 1, 'chrome',        505,  2599, 100.00, 0.00, 120, 0, 0, 0, 1200, 800,  600,  15, 'Starter cruiser'),
('typhoon_ii',     'Typhoon II',     'cruiser', 2, 'chrome',        1010, 5198, 100.00, 0.00, 180, 0, 0, 0, 3600, 2400, 1800, 30, 'Improved Typhoon hull'),
('typhoon_iii',    'Typhoon III',    'cruiser', 3, 'chrome',        2020, 10396,100.00, 0.00, 270, 0, 0, 0, 10800,7200, 5400, 60, 'Advanced Typhoon hull'),
-- Bombardier
('bombardier_i',   'Bombardier I',   'cruiser', 1, 'regen',         480,  2700, 100.00, 0.00, 125, 0, 0, 0, 1250, 830,  625,  15, 'Regenerative armor cruiser'),
('bombardier_ii',  'Bombardier II',  'cruiser', 2, 'regen',         960,  5400, 100.00, 0.00, 185, 0, 0, 0, 3750, 2490, 1875, 30, 'Improved Bombardier hull'),
('bombardier_iii', 'Bombardier III', 'cruiser', 3, 'regen',         1920, 10800,100.00, 0.00, 275, 0, 0, 0, 11250,7470, 5625, 60, 'Advanced Bombardier hull'),
-- Duke
('duke_i',         'Duke I',         'cruiser', 1, 'chrome',        520,  2500, 100.00, 0.00, 115, 0, 0, 0, 1180, 790,  590,  15, 'High-shield chrome cruiser'),
('duke_ii',        'Duke II',        'cruiser', 2, 'chrome',        1040, 5000, 100.00, 0.00, 175, 0, 0, 0, 3540, 2370, 1770, 30, 'Improved Duke hull'),
('duke_iii',       'Duke III',       'cruiser', 3, 'chrome',        2080, 10000,100.00, 0.00, 265, 0, 0, 0, 10620,7110, 5310, 60, 'Advanced Duke hull'),
-- The Shuttler
('the_shuttler_i', 'The Shuttler I', 'cruiser', 1, 'regen',         490,  2650, 100.00, 0.00, 122, 0, 0, 0, 1220, 810,  610,  15, 'Balanced regen cruiser'),
('the_shuttler_ii','The Shuttler II','cruiser', 2, 'regen',         980,  5300, 100.00, 0.00, 182, 0, 0, 0, 3660, 2430, 1830, 30, 'Improved Shuttler hull'),
('the_shuttler_iii','The Shuttler III','cruiser',3, 'regen',        1960, 10600,100.00, 0.00, 272, 0, 0, 0, 10980,7290, 5490, 60, 'Advanced Shuttler hull'),
-- Watchman
('watchman_i',     'Watchman I',     'cruiser', 1, 'chrome',        530,  2550, 100.00, 0.00, 118, 0, 0, 0, 1190, 795,  595,  15, 'Defensive chrome cruiser'),
('watchman_ii',    'Watchman II',    'cruiser', 2, 'chrome',        1060, 5100, 100.00, 0.00, 178, 0, 0, 0, 3570, 2385, 1785, 30, 'Improved Watchman hull'),
('watchman_iii',   'Watchman III',   'cruiser', 3, 'chrome',        2120, 10200,100.00, 0.00, 268, 0, 0, 0, 10710,7155, 5355, 60, 'Advanced Watchman hull'),
-- Spinner
('spinner_i',      'Spinner I',      'cruiser', 1, 'regen',         470,  2750, 100.00, 0.00, 128, 0, 0, 0, 1260, 840,  630,  15, 'High-capacity regen cruiser'),
('spinner_ii',     'Spinner II',     'cruiser', 2, 'regen',         940,  5500, 100.00, 0.00, 188, 0, 0, 0, 3780, 2520, 1890, 30, 'Improved Spinner hull'),
('spinner_iii',    'Spinner III',    'cruiser', 3, 'regen',         1880, 11000,100.00, 0.00, 278, 0, 0, 0, 11340,7560, 5670, 60, 'Advanced Spinner hull'),
-- Wraith
('wraith_i',       'Wraith I',       'cruiser', 1, 'chrome',        540,  2450, 100.00, 0.00, 112, 0, 0, 0, 1170, 780,  585,  15, 'Stealthy chrome cruiser'),
('wraith_ii',      'Wraith II',      'cruiser', 2, 'chrome',        1080, 4900, 100.00, 0.00, 172, 0, 0, 0, 3510, 2340, 1755, 30, 'Improved Wraith hull'),
('wraith_iii',     'Wraith III',     'cruiser', 3, 'chrome',        2160, 9800, 100.00, 0.00, 262, 0, 0, 0, 10530,7020, 5265, 60, 'Advanced Wraith hull'),
-- Encratos
('encratos_i',     'Encratos I',     'cruiser', 1, 'regen',         500,  2800, 100.00, 0.00, 130, 0, 0, 0, 1280, 850,  640,  15, 'Highest-capacity cruiser'),
('encratos_ii',    'Encratos II',    'cruiser', 2, 'regen',         1000, 5600, 100.00, 0.00, 190, 0, 0, 0, 3840, 2550, 1920, 30, 'Improved Encratos hull'),
('encratos_iii',   'Encratos III',   'cruiser', 3, 'regen',         2000, 11200,100.00, 0.00, 280, 0, 0, 0, 11520,7650, 5760, 60, 'Advanced Encratos hull'),
-- Nicholas
('nicholas_i',     'Nicholas I',     'cruiser', 1, 'chrome',        510,  2620, 100.00, 0.00, 119, 0, 0, 0, 1210, 805,  605,  15, 'Balanced chrome cruiser'),
('nicholas_ii',    'Nicholas II',    'cruiser', 2, 'chrome',        1020, 5240, 100.00, 0.00, 179, 0, 0, 0, 3630, 2415, 1815, 30, 'Improved Nicholas hull'),
('nicholas_iii',   'Nicholas III',   'cruiser', 3, 'chrome',        2040, 10480,100.00, 0.00, 269, 0, 0, 0, 10890,7245, 5445, 60, 'Advanced Nicholas hull'),
-- Helena
('helena_i',       'Helena I',       'cruiser', 1, 'regen',         495,  2680, 100.00, 0.00, 124, 0, 0, 0, 1240, 825,  618,  15, 'All-around regen cruiser'),
('helena_ii',      'Helena II',      'cruiser', 2, 'regen',         990,  5360, 100.00, 0.00, 184, 0, 0, 0, 3720, 2475, 1854, 30, 'Improved Helena hull'),
('helena_iii',     'Helena III',     'cruiser', 3, 'regen',         1980, 10720,100.00, 0.00, 274, 0, 0, 0, 11160,7425, 5562, 60, 'Advanced Helena hull'),

-- Battleship Hull Lines (5 lines x 3 tiers = 15)
-- All battleships: base_movement=0, base_agility=0, base_stability=100
-- Estrella
('estrella_i',     'Estrella I',     'battleship', 1, 'chrome',       380,  4200, 100.00, 0.00, 160, 0, 0, 0, 3000, 2000, 1500, 25, 'Starter battleship'),
('estrella_ii',    'Estrella II',    'battleship', 2, 'chrome',       760,  8400, 100.00, 0.00, 240, 0, 0, 0, 9000, 6000, 4500, 50, 'Improved Estrella hull'),
('estrella_iii',   'Estrella III',   'battleship', 3, 'chrome',       1520, 16800,100.00, 0.00, 360, 0, 0, 0, 27000,18000,13500,100, 'Advanced Estrella hull'),
-- Nettle
('nettle_i',       'Nettle I',       'battleship', 1, 'regen',        350,  4500, 100.00, 0.00, 165, 0, 0, 0, 3100, 2060, 1550, 25, 'High-structure regen battleship'),
('nettle_ii',      'Nettle II',      'battleship', 2, 'regen',        700,  9000, 100.00, 0.00, 245, 0, 0, 0, 9300, 6180, 4650, 50, 'Improved Nettle hull'),
('nettle_iii',     'Nettle III',     'battleship', 3, 'regen',        1400, 18000,100.00, 0.00, 370, 0, 0, 0, 27900,18540,13950,100, 'Advanced Nettle hull'),
-- Diaz
('diaz_i',         'Diaz I',         'battleship', 1, 'nano',         400,  4100, 100.00, 0.00, 155, 0, 0, 0, 2950, 1970, 1475, 25, 'Nano armor battleship'),
('diaz_ii',        'Diaz II',        'battleship', 2, 'nano',         800,  8200, 100.00, 0.00, 235, 0, 0, 0, 8850, 5910, 4425, 50, 'Improved Diaz hull'),
('diaz_iii',       'Diaz III',       'battleship', 3, 'nano',         1600, 16400,100.00, 0.00, 350, 0, 0, 0, 26550,17730,13275,100, 'Advanced Diaz hull'),
-- RV766-The Explorer
('rv766_i',        'RV766-The Explorer I',  'battleship', 1, 'neutralizing', 370, 4400, 100.00, 0.00, 170, 0, 0, 0, 3150, 2100, 1575, 25, 'Neutralizing explorer battleship'),
('rv766_ii',       'RV766-The Explorer II', 'battleship', 2, 'neutralizing', 740, 8800, 100.00, 0.00, 250, 0, 0, 0, 9450, 6300, 4725, 50, 'Improved RV766 hull'),
('rv766_iii',      'RV766-The Explorer III','battleship', 3, 'neutralizing', 1480,17600,100.00, 0.00, 375, 0, 0, 0, 28350,18900,14175,100, 'Advanced RV766 hull'),
-- Palenka
('palenka_i',      'Palenka I',      'battleship', 1, 'chrome',       360,  4350, 100.00, 0.00, 162, 0, 0, 0, 3050, 2030, 1525, 25, 'Chrome heavy battleship'),
('palenka_ii',     'Palenka II',     'battleship', 2, 'chrome',       720,  8700, 100.00, 0.00, 242, 0, 0, 0, 9150, 6090, 4575, 50, 'Improved Palenka hull'),
('palenka_iii',    'Palenka III',    'battleship', 3, 'chrome',       1440, 17400,100.00, 0.00, 365, 0, 0, 0, 27450,18270,13725,100, 'Advanced Palenka hull');

-- ============================================================================
-- PART 6: SEED DATA - MODULE TYPES (97 modules)
-- ============================================================================

-- Attack Modules - Ballistic (Range 1-2, Cooldown 0)
INSERT INTO module_types (name, display_name, category, tier, damage_type, min_damage, max_damage, weapon_range_min, weapon_range_max, cooldown, he3_per_round, volume, max_per_ship, effects_json, metal_cost, he3_cost, gold_cost, build_time_seconds, description) VALUES
('rapid_fire',                'Rapid Fire',                'ballistic', 1, 'kinetic',   12,  18,  1, 2, 0, 2,  8,  0, '{}', 100, 60, 40, 3, 'Basic kinetic ballistic weapon'),
('rapid_fire',                'Rapid Fire',                'ballistic', 2, 'kinetic',   24,  36,  1, 2, 0, 4,  12, 0, '{}', 300, 180, 120, 5, 'Improved kinetic ballistic weapon'),
('rapid_fire',                'Rapid Fire',                'ballistic', 3, 'kinetic',   48,  72,  1, 2, 0, 8,  18, 0, '{}', 900, 540, 360, 8, 'Advanced kinetic ballistic weapon'),
('taskmaster',                'Taskmaster',                'ballistic', 1, 'heat',      14,  20,  1, 2, 0, 2,  9,  0, '{}', 110, 66, 44, 3, 'Heat ballistic weapon'),
('taskmaster',                'Taskmaster',                'ballistic', 2, 'heat',      28,  40,  1, 2, 0, 4,  14, 0, '{}', 330, 198, 132, 5, 'Improved heat ballistic weapon'),
('taskmaster',                'Taskmaster',                'ballistic', 3, 'heat',      56,  80,  1, 2, 0, 8,  20, 0, '{}', 990, 594, 396, 8, 'Advanced heat ballistic weapon'),
('gatling_cannon',            'Gatling Cannon',            'ballistic', 1, 'kinetic',   16,  22,  1, 2, 0, 3,  10, 0, '{}', 120, 72, 48, 4, 'Heavy kinetic ballistic weapon'),
('gatling_cannon',            'Gatling Cannon',            'ballistic', 2, 'kinetic',   32,  44,  1, 2, 0, 6,  15, 0, '{}', 360, 216, 144, 6, 'Improved heavy ballistic weapon'),
('gatling_cannon',            'Gatling Cannon',            'ballistic', 3, 'kinetic',   64,  88,  1, 2, 0, 12, 22, 0, '{}', 1080, 648, 432, 10, 'Advanced heavy ballistic weapon'),

-- Attack Modules - Directional (Range 2-5, Cooldown 1)
('cluster_laser_transmitter', 'Cluster Laser Transmitter', 'directional', 1, 'heat',    30,  45,  2, 5, 1, 4,  12, 0, '{}', 150, 90, 60, 5, 'Heat directional weapon'),
('cluster_laser_transmitter', 'Cluster Laser Transmitter', 'directional', 2, 'heat',    60,  90,  2, 5, 1, 8,  18, 0, '{}', 450, 270, 180, 8, 'Improved heat directional weapon'),
('cluster_laser_transmitter', 'Cluster Laser Transmitter', 'directional', 3, 'heat',    120, 180, 2, 5, 1, 16, 26, 0, '{}', 1350, 810, 540, 12, 'Advanced heat directional weapon'),
('magneto_pulsar',            'Magneto Pulsar',            'directional', 1, 'magnetic', 35,  50,  2, 5, 1, 5,  14, 0, '{}', 170, 102, 68, 5, 'Magnetic directional weapon'),
('magneto_pulsar',            'Magneto Pulsar',            'directional', 2, 'magnetic', 70,  100, 2, 5, 1, 10, 20, 0, '{}', 510, 306, 204, 8, 'Improved magnetic directional weapon'),
('magneto_pulsar',            'Magneto Pulsar',            'directional', 3, 'magnetic', 140, 200, 2, 5, 1, 20, 28, 0, '{}', 1530, 918, 612, 12, 'Advanced magnetic directional weapon'),

-- Attack Modules - Missile (Range 5-8, Cooldown 3)
('rocket_frame',              'Rocket Frame',              'missile', 1, 'explosive',    80,  120, 5, 8, 3, 8,  16, 0, '{}', 200, 120, 80, 6, 'Standard explosive missile'),
('rocket_frame',              'Rocket Frame',              'missile', 2, 'explosive',    160, 240, 5, 8, 3, 16, 24, 0, '{}', 600, 360, 240, 10, 'Improved explosive missile'),
('rocket_frame',              'Rocket Frame',              'missile', 3, 'explosive',    320, 480, 5, 8, 3, 32, 34, 0, '{}', 1800, 1080, 720, 15, 'Advanced explosive missile'),
('starlight_missile_pod',     'Starlight Missile Pod',     'missile', 1, 'explosive',    90,  135, 5, 8, 3, 10, 18, 0, '{}', 220, 132, 88, 7, 'Heavy explosive missile pod'),
('starlight_missile_pod',     'Starlight Missile Pod',     'missile', 2, 'explosive',    180, 270, 5, 8, 3, 20, 26, 0, '{}', 660, 396, 264, 11, 'Improved heavy missile pod'),
('starlight_missile_pod',     'Starlight Missile Pod',     'missile', 3, 'explosive',    360, 540, 5, 8, 3, 40, 36, 0, '{}', 1980, 1188, 792, 16, 'Advanced heavy missile pod'),

-- Attack Modules - Ship-Based Weapons (Range 6-10, Cooldown 4)
('streamliner',               'Streamliner',               'ship_based', 1, 'kinetic',  120, 180, 6, 10, 4, 16, 20, 0, '{}', 300, 180, 120, 8, 'Kinetic ship-based weapon'),
('streamliner',               'Streamliner',               'ship_based', 2, 'kinetic',  240, 360, 6, 10, 4, 32, 30, 0, '{}', 900, 540, 360, 12, 'Improved kinetic SBW'),
('streamliner',               'Streamliner',               'ship_based', 3, 'kinetic',  480, 720, 6, 10, 4, 64, 42, 0, '{}', 2700, 1620, 1080, 18, 'Advanced kinetic SBW'),
('golem',                     'Golem',                     'ship_based', 1, 'magnetic',  130, 200, 6, 10, 4, 18, 22, 0, '{}', 330, 198, 132, 9, 'Magnetic ship-based weapon'),
('golem',                     'Golem',                     'ship_based', 2, 'magnetic',  260, 400, 6, 10, 4, 36, 32, 0, '{}', 990, 594, 396, 13, 'Improved magnetic SBW'),
('golem',                     'Golem',                     'ship_based', 3, 'magnetic',  520, 800, 6, 10, 4, 72, 44, 0, '{}', 2970, 1782, 1188, 19, 'Advanced magnetic SBW'),

-- Attack Modules - Planetary (Range 1-2, Cooldown 1)
('lander_module',             'Lander Module',             'planetary', 1, 'siege',      50,  75,  1, 2, 1, 4,  10, 0, '{}', 130, 78, 52, 4, 'Siege weapon for planetary structures'),
('lander_module',             'Lander Module',             'planetary', 2, 'siege',      100, 150, 1, 2, 1, 8,  16, 0, '{}', 390, 234, 156, 7, 'Improved siege weapon'),
('lander_module',             'Lander Module',             'planetary', 3, 'siege',      200, 300, 1, 2, 1, 16, 24, 0, '{}', 1170, 702, 468, 11, 'Advanced siege weapon'),

-- Defense Modules - Structure
('atomic_framework',          'Atomic Framework',          'structure', 0, NULL, 0, 0, 0, 0, 0, 0, 6, 0, '{"structure_bonus": 500}', 50, 30, 20, 1, '+500 structure per module'),
('ship_reinforcement_facility','Ship Reinforcement Facility','structure', 1, NULL, 0, 0, 0, 0, 0, 0, 8, 0, '{"damage_reduction": 1}', 100, 60, 40, 4, 'Reduces damage by 1 per module'),
('ship_reinforcement_facility','Ship Reinforcement Facility','structure', 2, NULL, 0, 0, 0, 0, 0, 0, 12, 0, '{"damage_reduction": 2}', 300, 180, 120, 6, 'Reduces damage by 2 per module'),
('ship_reinforcement_facility','Ship Reinforcement Facility','structure', 3, NULL, 0, 0, 0, 0, 0, 0, 18, 0, '{"damage_reduction": 3}', 900, 540, 360, 9, 'Reduces damage by 3 per module'),
('quick_reaction_armor',      'Quick Reaction Armor',      'structure', 1, NULL, 0, 0, 0, 0, 0, 0, 10, 1, '{"reflect_damage_pct": 5}', 120, 72, 48, 5, '5% reflect damage'),
('quick_reaction_armor',      'Quick Reaction Armor',      'structure', 2, NULL, 0, 0, 0, 0, 0, 0, 15, 1, '{"reflect_damage_pct": 10}', 360, 216, 144, 8, '10% reflect damage'),
('quick_reaction_armor',      'Quick Reaction Armor',      'structure', 3, NULL, 0, 0, 0, 0, 0, 0, 22, 1, '{"reflect_damage_pct": 15}', 1080, 648, 432, 12, '15% reflect damage'),
('reflective_plating',        'Reflective Plating',        'structure', 1, NULL, 0, 0, 0, 0, 0, 0, 8, 1, '{"defense_bonus_pct": 3}', 100, 60, 40, 4, '+3% defense'),
('reflective_plating',        'Reflective Plating',        'structure', 2, NULL, 0, 0, 0, 0, 0, 0, 12, 1, '{"defense_bonus_pct": 6}', 300, 180, 120, 7, '+6% defense'),
('reflective_plating',        'Reflective Plating',        'structure', 3, NULL, 0, 0, 0, 0, 0, 0, 18, 1, '{"defense_bonus_pct": 10}', 900, 540, 360, 10, '+10% defense'),
('energy_armor',              'Energy Armor',              'structure', 1, NULL, 0, 0, 0, 0, 0, 5, 10, 1, '{"structure_bonus": 200}', 130, 78, 52, 5, '+200 structure, 5 He3/round'),
('energy_armor',              'Energy Armor',              'structure', 2, NULL, 0, 0, 0, 0, 0, 10, 16, 1, '{"structure_bonus": 400}', 390, 234, 156, 8, '+400 structure, 10 He3/round'),
('energy_armor',              'Energy Armor',              'structure', 3, NULL, 0, 0, 0, 0, 0, 15, 24, 1, '{"structure_bonus": 600}', 1170, 702, 468, 12, '+600 structure, 15 He3/round'),
('daedalus_control_system',   'Daedalus Control System',   'structure', 1, NULL, 0, 0, 0, 0, 0, 0, 12, 1, '{"structure_bonus_pct": 3, "damage_reduction": 1}', 150, 90, 60, 6, '+3% structure and damage reduction'),
('daedalus_control_system',   'Daedalus Control System',   'structure', 2, NULL, 0, 0, 0, 0, 0, 0, 18, 1, '{"structure_bonus_pct": 6, "damage_reduction": 2}', 450, 270, 180, 9, '+6% structure and damage reduction'),
('daedalus_control_system',   'Daedalus Control System',   'structure', 3, NULL, 0, 0, 0, 0, 0, 0, 26, 1, '{"structure_bonus_pct": 10, "damage_reduction": 3}', 1350, 810, 540, 13, '+10% structure and damage reduction'),

-- Defense Modules - Shields
('orbital_shield',            'Orbital Shield',            'shield', 0, NULL, 0, 0, 0, 0, 0, 0, 6, 0, '{"shield_bonus": 300}', 50, 30, 20, 1, '+300 shield per module'),
('energy_shield_booster',     'Energy Shield Booster',     'shield', 1, NULL, 0, 0, 0, 0, 0, 0, 8, 1, '{"shield_effectiveness_pct": 5}', 100, 60, 40, 4, '+5% shield effectiveness'),
('energy_shield_booster',     'Energy Shield Booster',     'shield', 2, NULL, 0, 0, 0, 0, 0, 0, 12, 1, '{"shield_effectiveness_pct": 10}', 300, 180, 120, 7, '+10% shield effectiveness'),
('energy_shield_booster',     'Energy Shield Booster',     'shield', 3, NULL, 0, 0, 0, 0, 0, 0, 18, 1, '{"shield_effectiveness_pct": 15}', 900, 540, 360, 10, '+15% shield effectiveness'),
('particle_stun_shield',      'Particle Stun Shield',      'shield', 1, NULL, 0, 0, 0, 0, 0, 0, 10, 0, '{"kinetic_damage_reduction": 5}', 120, 72, 48, 5, 'Reduces ballistic/SBW dmg by 5'),
('particle_stun_shield',      'Particle Stun Shield',      'shield', 2, NULL, 0, 0, 0, 0, 0, 0, 15, 0, '{"kinetic_damage_reduction": 10}', 360, 216, 144, 8, 'Reduces ballistic/SBW dmg by 10'),
('particle_stun_shield',      'Particle Stun Shield',      'shield', 3, NULL, 0, 0, 0, 0, 0, 0, 22, 0, '{"kinetic_damage_reduction": 15}', 1080, 648, 432, 12, 'Reduces ballistic/SBW dmg by 15'),
('heat_diffusion_shield',     'Heat Diffusion Shield',     'shield', 1, NULL, 0, 0, 0, 0, 0, 0, 10, 0, '{"heat_damage_reduction": 5}', 120, 72, 48, 5, 'Reduces heat dmg by 5'),
('heat_diffusion_shield',     'Heat Diffusion Shield',     'shield', 2, NULL, 0, 0, 0, 0, 0, 0, 15, 0, '{"heat_damage_reduction": 10}', 360, 216, 144, 8, 'Reduces heat dmg by 10'),
('heat_diffusion_shield',     'Heat Diffusion Shield',     'shield', 3, NULL, 0, 0, 0, 0, 0, 0, 22, 0, '{"heat_damage_reduction": 15}', 1080, 648, 432, 12, 'Reduces heat dmg by 15'),
('space_time_magnetic_shield','Space-Time Magnetic Shield','shield', 1, NULL, 0, 0, 0, 0, 0, 0, 10, 0, '{"magnetic_damage_reduction": 5}', 120, 72, 48, 5, 'Reduces magnetic dmg by 5'),
('space_time_magnetic_shield','Space-Time Magnetic Shield','shield', 2, NULL, 0, 0, 0, 0, 0, 0, 15, 0, '{"magnetic_damage_reduction": 10}', 360, 216, 144, 8, 'Reduces magnetic dmg by 10'),
('space_time_magnetic_shield','Space-Time Magnetic Shield','shield', 3, NULL, 0, 0, 0, 0, 0, 0, 22, 0, '{"magnetic_damage_reduction": 15}', 1080, 648, 432, 12, 'Reduces magnetic dmg by 15'),
('detonator_shield',          'Detonator Shield',          'shield', 1, NULL, 0, 0, 0, 0, 0, 0, 10, 0, '{"explosive_damage_reduction": 5}', 120, 72, 48, 5, 'Reduces explosive dmg by 5'),
('detonator_shield',          'Detonator Shield',          'shield', 2, NULL, 0, 0, 0, 0, 0, 0, 15, 0, '{"explosive_damage_reduction": 10}', 360, 216, 144, 8, 'Reduces explosive dmg by 10'),
('detonator_shield',          'Detonator Shield',          'shield', 3, NULL, 0, 0, 0, 0, 0, 0, 22, 0, '{"explosive_damage_reduction": 15}', 1080, 648, 432, 12, 'Reduces explosive dmg by 15'),
('shield_regenerator',        'Shield Regenerator',        'shield', 1, NULL, 0, 0, 0, 0, 0, 0, 8, 1, '{"shield_restore_pct": 10}', 100, 60, 40, 4, '+10% shield restore/round'),
('shield_regenerator',        'Shield Regenerator',        'shield', 2, NULL, 0, 0, 0, 0, 0, 0, 12, 1, '{"shield_restore_pct": 20}', 300, 180, 120, 7, '+20% shield restore/round'),
('shield_regenerator',        'Shield Regenerator',        'shield', 3, NULL, 0, 0, 0, 0, 0, 0, 18, 1, '{"shield_restore_pct": 30}', 900, 540, 360, 10, '+30% shield restore/round'),
('eos_phase_shift_shield',    'EOS Phase Shift Shield',    'shield', 1, NULL, 0, 0, 0, 0, 0, 0, 14, 1, '{"absorb_double_pct": 10}', 180, 108, 72, 7, '10% absorb double dmg'),
('eos_phase_shift_shield',    'EOS Phase Shift Shield',    'shield', 2, NULL, 0, 0, 0, 0, 0, 0, 20, 1, '{"absorb_double_pct": 20}', 540, 324, 216, 10, '20% absorb double dmg'),
('eos_phase_shift_shield',    'EOS Phase Shift Shield',    'shield', 3, NULL, 0, 0, 0, 0, 0, 0, 28, 1, '{"absorb_double_pct": 30}', 1620, 972, 648, 14, '30% absorb double dmg'),

-- Defense Modules - Air Defense
('anti_aircraft_cannon',      'Anti-Aircraft Cannon',      'air_defense', 1, NULL, 0, 0, 0, 0, 0, 0, 8, 0, '{"intercept_missile_pct": 15}', 100, 60, 40, 4, '15% intercept missiles'),
('anti_aircraft_cannon',      'Anti-Aircraft Cannon',      'air_defense', 2, NULL, 0, 0, 0, 0, 0, 0, 12, 0, '{"intercept_missile_pct": 25}', 300, 180, 120, 7, '25% intercept missiles'),
('anti_aircraft_cannon',      'Anti-Aircraft Cannon',      'air_defense', 3, NULL, 0, 0, 0, 0, 0, 0, 18, 0, '{"intercept_missile_pct": 35}', 900, 540, 360, 10, '35% intercept missiles'),
('powered_pulse_cannon',      'Powered Pulse Cannon',      'air_defense', 1, NULL, 0, 0, 0, 0, 0, 0, 10, 0, '{"intercept_any_pct": 35}', 130, 78, 52, 5, '35% intercept any attack'),
('powered_pulse_cannon',      'Powered Pulse Cannon',      'air_defense', 2, NULL, 0, 0, 0, 0, 0, 0, 15, 0, '{"intercept_any_pct": 45}', 390, 234, 156, 8, '45% intercept any attack'),
('powered_pulse_cannon',      'Powered Pulse Cannon',      'air_defense', 3, NULL, 0, 0, 0, 0, 0, 0, 22, 0, '{"intercept_any_pct": 55}', 1170, 702, 468, 12, '55% intercept any attack'),
('extreme_counterattack',     'Extreme Counterattack',     'air_defense', 1, NULL, 0, 0, 0, 0, 0, 0, 12, 1, '{"reflect_intercepted_pct": 20}', 150, 90, 60, 6, 'Reflect 20% intercepted damage'),
('extreme_counterattack',     'Extreme Counterattack',     'air_defense', 2, NULL, 0, 0, 0, 0, 0, 0, 18, 1, '{"reflect_intercepted_pct": 35}', 450, 270, 180, 9, 'Reflect 35% intercepted damage'),
('extreme_counterattack',     'Extreme Counterattack',     'air_defense', 3, NULL, 0, 0, 0, 0, 0, 0, 26, 1, '{"reflect_intercepted_pct": 50}', 1350, 810, 540, 13, 'Reflect 50% intercepted damage'),

-- Auxiliary Modules - Electronic (1 per ship each)
('agility_booster',           'Agility Booster',           'electronic', 1, NULL, 0, 0, 0, 0, 0, 0, 6, 1, '{"agility_bonus": 1}', 80, 48, 32, 3, '+1 agility'),
('agility_booster',           'Agility Booster',           'electronic', 2, NULL, 0, 0, 0, 0, 0, 0, 9, 1, '{"agility_bonus": 2}', 240, 144, 96, 5, '+2 agility'),
('agility_booster',           'Agility Booster',           'electronic', 3, NULL, 0, 0, 0, 0, 0, 0, 13, 1, '{"agility_bonus": 3}', 720, 432, 288, 8, '+3 agility'),
('infrared_scanner',          'Infrared Scanner',          'electronic', 1, NULL, 0, 0, 0, 0, 0, 0, 6, 1, '{"steering_bonus": 1}', 80, 48, 32, 3, '+1 steering'),
('infrared_scanner',          'Infrared Scanner',          'electronic', 2, NULL, 0, 0, 0, 0, 0, 0, 9, 1, '{"steering_bonus": 2}', 240, 144, 96, 5, '+2 steering'),
('infrared_scanner',          'Infrared Scanner',          'electronic', 3, NULL, 0, 0, 0, 0, 0, 0, 13, 1, '{"steering_bonus": 3}', 720, 432, 288, 8, '+3 steering'),
('ecm_booster',               'ECM Booster',               'electronic', 1, NULL, 0, 0, 0, 0, 0, 0, 7, 1, '{"dodge_chance_pct": 5}', 90, 54, 36, 4, '+5% dodge chance'),
('ecm_booster',               'ECM Booster',               'electronic', 2, NULL, 0, 0, 0, 0, 0, 0, 10, 1, '{"dodge_chance_pct": 10}', 270, 162, 108, 6, '+10% dodge chance'),
('ecm_booster',               'ECM Booster',               'electronic', 3, NULL, 0, 0, 0, 0, 0, 0, 14, 1, '{"dodge_chance_pct": 15}', 810, 486, 324, 9, '+15% dodge chance'),
('auto_target_system',        'Auto Target System',        'electronic', 1, NULL, 0, 0, 0, 0, 0, 0, 7, 1, '{"hit_chance_pct": 5}', 90, 54, 36, 4, '+5% hit chance'),
('auto_target_system',        'Auto Target System',        'electronic', 2, NULL, 0, 0, 0, 0, 0, 0, 10, 1, '{"hit_chance_pct": 10}', 270, 162, 108, 6, '+10% hit chance'),
('auto_target_system',        'Auto Target System',        'electronic', 3, NULL, 0, 0, 0, 0, 0, 0, 14, 1, '{"hit_chance_pct": 15}', 810, 486, 324, 9, '+15% hit chance'),
('time_dilation_module',      'Time Dilation Module',      'electronic', 1, NULL, 0, 0, 0, 0, 0, 0, 8, 1, '{"critical_hit_pct": 3}', 100, 60, 40, 5, '+3% critical hit rate'),
('time_dilation_module',      'Time Dilation Module',      'electronic', 2, NULL, 0, 0, 0, 0, 0, 0, 12, 1, '{"critical_hit_pct": 6}', 300, 180, 120, 7, '+6% critical hit rate'),
('time_dilation_module',      'Time Dilation Module',      'electronic', 3, NULL, 0, 0, 0, 0, 0, 0, 16, 1, '{"critical_hit_pct": 10}', 900, 540, 360, 10, '+10% critical hit rate'),

-- Auxiliary Modules - Storage
('station_warehouse',         'Station Warehouse',         'storage', 0, NULL, 0, 0, 0, 0, 0, 0, 4, 0, '{"he3_storage_bonus": 200}', 40, 24, 16, 2, '+200 He3 storage'),
('nano_station_warehouse',    'Nano Station Warehouse',    'storage', 0, NULL, 0, 0, 0, 0, 0, 0, 8, 1, '{"he3_storage_bonus": 500, "module_capacity_bonus": 50}', 100, 60, 40, 4, '+500 He3 storage, +50 module capacity'),

-- Auxiliary Modules - Transmission
('super_transmission_engine', 'Super Transmission Engine', 'transmission', 0, NULL, 0, 0, 0, 0, 0, 0, 8, 0, '{"movement_bonus": 1}', 100, 60, 40, 3, '+1 movement'),
('team_combat_engine',        'Team Combat Engine',        'transmission', 1, NULL, 0, 0, 0, 0, 0, 0, 10, 1, '{"movement_bonus": 1, "agility_bonus": 1}', 130, 78, 52, 5, '+1 movement, +1 agility'),
('team_combat_engine',        'Team Combat Engine',        'transmission', 2, NULL, 0, 0, 0, 0, 0, 0, 14, 1, '{"movement_bonus": 1, "agility_bonus": 2}', 390, 234, 156, 7, '+1 movement, +2 agility'),
('team_combat_engine',        'Team Combat Engine',        'transmission', 3, NULL, 0, 0, 0, 0, 0, 0, 18, 1, '{"movement_bonus": 1, "agility_bonus": 3}', 1170, 702, 468, 10, '+1 movement, +3 agility'),
('anti_matter_engine',        'Anti-Matter Engine',        'transmission', 1, NULL, 0, 0, 0, 0, 0, 0, 12, 1, '{"movement_bonus": 2, "agility_bonus": 1}', 150, 90, 60, 6, '+2 movement, +1 agility'),
('anti_matter_engine',        'Anti-Matter Engine',        'transmission', 2, NULL, 0, 0, 0, 0, 0, 0, 16, 1, '{"movement_bonus": 2, "agility_bonus": 2}', 450, 270, 180, 8, '+2 movement, +2 agility'),
('anti_matter_engine',        'Anti-Matter Engine',        'transmission', 3, NULL, 0, 0, 0, 0, 0, 0, 20, 1, '{"movement_bonus": 2, "agility_bonus": 3}', 1350, 810, 540, 11, '+2 movement, +3 agility');

-- ============================================================================
-- PART 7: SEED DATA - SHIP FACTORY LEVELS (24 rows)
-- ============================================================================

INSERT INTO ship_factory_levels VALUES
(1,  1,  600,       450,       500,       200,      1,  1),
(2,  1,  1818,      1364,      1515,      574,      3,  1),
(3,  2,  5509,      4132,      4590,      1647,     5,  1),
(4,  2,  16692,     12519,     13910,     4727,     7,  2),
(5,  3,  50576,     37932,     42147,     13566,    9,  2),
(6,  3,  153244,    114933,    127704,    38934,    11, 2),
(7,  4,  464330,    348248,    386942,    111760,   14, 2),
(8,  4,  1406920,   1055190,   1172434,   320751,   17, 3),
(9,  5,  4264968,   3198726,   3554140,   920555,   20, 3),
(10, 5,  12924853,  9693640,   10770711,  2641993,  23, 3),
(11, 6,  39170305,  29377729,  32641855,  7582522,  25, 3),
(12, 6,  118706025, 89029519,  98921822,  21761838, 27, 4),
(13, 7,  359699416, 269774562, 299749513, 62454475, 30, 4),
(14, 7,  1090000000,817500000, 908333333, 179244223,33, 4),
(15, 8,  3303000000,2477250000,2752500000,514370880,36, 4),
(16, 8,  10008000000,7506000000,8340000000,1476244425,39,4),
(17, 9,  30324000000,22743000000,25270000000,4236821340,42,4),
(18, 9,  91882000000,68911500000,76568333333,12159717246,45,4),
(19, 10, 278400000000,208800000000,232000000000,34898378376,48,4),
(20, 10, 843600000000,632700000000,702870000000,100158345738,51,4),
(21, 11, 2556100000000,1917075000000,2130083333333,287454452118,54,4),
(22, 11, 7744983000000,5808737250000,6454152500000,824914057499,56,4),
(23, 12, 23467297000000,17600473000000,19556081111111,2367503784943,58,4),
(24, 12, 71105509000000,53329132000000,59254591111111,6792715862065,60,4);

-- ============================================================================
-- PART 8: SEED DATA - BLUEPRINTS (25 hull + 36 module lines = 61 blueprints)
-- ============================================================================

-- Hull Blueprints (25 hull lines, one blueprint per line referencing Tier I)
-- Blueprint references the Tier I hull; the blueprint unlocks all tiers of that line
INSERT INTO blueprints (name, blueprint_type, hull_type_id, source, description) VALUES
('Weikes Blueprint',       'hull', (SELECT id FROM hull_types WHERE name = 'weikes_i'),       'quest', 'Unlocks Weikes frigate line'),
('Air Wanderer Blueprint', 'hull', (SELECT id FROM hull_types WHERE name = 'air_wanderer_i'), 'instance', 'Unlocks Air Wanderer frigate line'),
('Valkyrie Blueprint',     'hull', (SELECT id FROM hull_types WHERE name = 'valkyrie_i'),     'instance', 'Unlocks Valkyrie frigate line'),
('GoGetter Blueprint',     'hull', (SELECT id FROM hull_types WHERE name = 'gogetter_i'),     'instance', 'Unlocks GoGetter frigate line'),
('Space Hunter Blueprint', 'hull', (SELECT id FROM hull_types WHERE name = 'space_hunter_i'), 'instance', 'Unlocks Space Hunter frigate line'),
('Sparrow Blueprint',      'hull', (SELECT id FROM hull_types WHERE name = 'sparrow_i'),      'instance', 'Unlocks Sparrow frigate line'),
('Devourer Blueprint',     'hull', (SELECT id FROM hull_types WHERE name = 'devourer_i'),     'instance', 'Unlocks Devourer frigate line'),
('Polymesus Blueprint',    'hull', (SELECT id FROM hull_types WHERE name = 'polymesus_i'),    'instance', 'Unlocks Polymesus frigate line'),
('Cybra Blueprint',        'hull', (SELECT id FROM hull_types WHERE name = 'cybra_i'),        'instance', 'Unlocks Cybra frigate line'),
('Hamdar Blueprint',       'hull', (SELECT id FROM hull_types WHERE name = 'hamdar_i'),       'instance', 'Unlocks Hamdar frigate line'),
('Typhoon Blueprint',      'hull', (SELECT id FROM hull_types WHERE name = 'typhoon_i'),      'quest', 'Unlocks Typhoon cruiser line'),
('Bombardier Blueprint',   'hull', (SELECT id FROM hull_types WHERE name = 'bombardier_i'),   'instance', 'Unlocks Bombardier cruiser line'),
('Duke Blueprint',         'hull', (SELECT id FROM hull_types WHERE name = 'duke_i'),         'instance', 'Unlocks Duke cruiser line'),
('The Shuttler Blueprint', 'hull', (SELECT id FROM hull_types WHERE name = 'the_shuttler_i'), 'instance', 'Unlocks The Shuttler cruiser line'),
('Watchman Blueprint',     'hull', (SELECT id FROM hull_types WHERE name = 'watchman_i'),     'instance', 'Unlocks Watchman cruiser line'),
('Spinner Blueprint',      'hull', (SELECT id FROM hull_types WHERE name = 'spinner_i'),      'instance', 'Unlocks Spinner cruiser line'),
('Wraith Blueprint',       'hull', (SELECT id FROM hull_types WHERE name = 'wraith_i'),       'instance', 'Unlocks Wraith cruiser line'),
('Encratos Blueprint',     'hull', (SELECT id FROM hull_types WHERE name = 'encratos_i'),     'instance', 'Unlocks Encratos cruiser line'),
('Nicholas Blueprint',     'hull', (SELECT id FROM hull_types WHERE name = 'nicholas_i'),     'instance', 'Unlocks Nicholas cruiser line'),
('Helena Blueprint',       'hull', (SELECT id FROM hull_types WHERE name = 'helena_i'),       'instance', 'Unlocks Helena cruiser line'),
('Estrella Blueprint',     'hull', (SELECT id FROM hull_types WHERE name = 'estrella_i'),     'quest', 'Unlocks Estrella battleship line'),
('Nettle Blueprint',       'hull', (SELECT id FROM hull_types WHERE name = 'nettle_i'),       'instance', 'Unlocks Nettle battleship line'),
('Diaz Blueprint',         'hull', (SELECT id FROM hull_types WHERE name = 'diaz_i'),         'instance', 'Unlocks Diaz battleship line'),
('RV766-The Explorer Blueprint', 'hull', (SELECT id FROM hull_types WHERE name = 'rv766_i'),  'instance', 'Unlocks RV766-The Explorer battleship line'),
('Palenka Blueprint',      'hull', (SELECT id FROM hull_types WHERE name = 'palenka_i'),      'instance', 'Unlocks Palenka battleship line');

-- Module Blueprints (36 module lines, one blueprint per line referencing Tier I or Tier 0)
-- Each module blueprint unlocks all tiers of that module line
INSERT INTO blueprints (name, blueprint_type, module_type_id, source, description) VALUES
('Rapid Fire Blueprint',                'module', (SELECT id FROM module_types WHERE name = 'rapid_fire' AND tier = 1),                'quest', 'Unlocks Rapid Fire weapon line'),
('Taskmaster Blueprint',                'module', (SELECT id FROM module_types WHERE name = 'taskmaster' AND tier = 1),                'instance', 'Unlocks Taskmaster weapon line'),
('Gatling Cannon Blueprint',            'module', (SELECT id FROM module_types WHERE name = 'gatling_cannon' AND tier = 1),            'instance', 'Unlocks Gatling Cannon weapon line'),
('Cluster Laser Transmitter Blueprint', 'module', (SELECT id FROM module_types WHERE name = 'cluster_laser_transmitter' AND tier = 1), 'instance', 'Unlocks Cluster Laser Transmitter weapon line'),
('Magneto Pulsar Blueprint',            'module', (SELECT id FROM module_types WHERE name = 'magneto_pulsar' AND tier = 1),            'instance', 'Unlocks Magneto Pulsar weapon line'),
('Rocket Frame Blueprint',              'module', (SELECT id FROM module_types WHERE name = 'rocket_frame' AND tier = 1),              'instance', 'Unlocks Rocket Frame weapon line'),
('Starlight Missile Pod Blueprint',     'module', (SELECT id FROM module_types WHERE name = 'starlight_missile_pod' AND tier = 1),     'instance', 'Unlocks Starlight Missile Pod weapon line'),
('Streamliner Blueprint',              'module', (SELECT id FROM module_types WHERE name = 'streamliner' AND tier = 1),               'instance', 'Unlocks Streamliner weapon line'),
('Golem Blueprint',                     'module', (SELECT id FROM module_types WHERE name = 'golem' AND tier = 1),                     'instance', 'Unlocks Golem weapon line'),
('Lander Module Blueprint',             'module', (SELECT id FROM module_types WHERE name = 'lander_module' AND tier = 1),             'instance', 'Unlocks Lander Module weapon line'),
('Atomic Framework Blueprint',          'module', (SELECT id FROM module_types WHERE name = 'atomic_framework' AND tier = 0),          'instance', 'Unlocks Atomic Framework'),
('Ship Reinforcement Facility Blueprint','module', (SELECT id FROM module_types WHERE name = 'ship_reinforcement_facility' AND tier = 1),'instance', 'Unlocks Ship Reinforcement Facility line'),
('Quick Reaction Armor Blueprint',      'module', (SELECT id FROM module_types WHERE name = 'quick_reaction_armor' AND tier = 1),      'instance', 'Unlocks Quick Reaction Armor line'),
('Reflective Plating Blueprint',        'module', (SELECT id FROM module_types WHERE name = 'reflective_plating' AND tier = 1),        'instance', 'Unlocks Reflective Plating line'),
('Energy Armor Blueprint',              'module', (SELECT id FROM module_types WHERE name = 'energy_armor' AND tier = 1),              'instance', 'Unlocks Energy Armor line'),
('Daedalus Control System Blueprint',   'module', (SELECT id FROM module_types WHERE name = 'daedalus_control_system' AND tier = 1),   'instance', 'Unlocks Daedalus Control System line'),
('Orbital Shield Blueprint',            'module', (SELECT id FROM module_types WHERE name = 'orbital_shield' AND tier = 0),            'instance', 'Unlocks Orbital Shield'),
('Energy Shield Booster Blueprint',     'module', (SELECT id FROM module_types WHERE name = 'energy_shield_booster' AND tier = 1),     'quest', 'Unlocks Energy Shield Booster line'),
('Particle Stun Shield Blueprint',      'module', (SELECT id FROM module_types WHERE name = 'particle_stun_shield' AND tier = 1),      'instance', 'Unlocks Particle Stun Shield line'),
('Heat Diffusion Shield Blueprint',     'module', (SELECT id FROM module_types WHERE name = 'heat_diffusion_shield' AND tier = 1),     'instance', 'Unlocks Heat Diffusion Shield line'),
('Space-Time Magnetic Shield Blueprint','module', (SELECT id FROM module_types WHERE name = 'space_time_magnetic_shield' AND tier = 1),'instance', 'Unlocks Space-Time Magnetic Shield line'),
('Detonator Shield Blueprint',          'module', (SELECT id FROM module_types WHERE name = 'detonator_shield' AND tier = 1),          'instance', 'Unlocks Detonator Shield line'),
('Shield Regenerator Blueprint',        'module', (SELECT id FROM module_types WHERE name = 'shield_regenerator' AND tier = 1),        'instance', 'Unlocks Shield Regenerator line'),
('EOS Phase Shift Shield Blueprint',    'module', (SELECT id FROM module_types WHERE name = 'eos_phase_shift_shield' AND tier = 1),    'instance', 'Unlocks EOS Phase Shift Shield line'),
('Anti-Aircraft Cannon Blueprint',      'module', (SELECT id FROM module_types WHERE name = 'anti_aircraft_cannon' AND tier = 1),      'instance', 'Unlocks Anti-Aircraft Cannon line'),
('Powered Pulse Cannon Blueprint',      'module', (SELECT id FROM module_types WHERE name = 'powered_pulse_cannon' AND tier = 1),      'instance', 'Unlocks Powered Pulse Cannon line'),
('Extreme Counterattack Blueprint',     'module', (SELECT id FROM module_types WHERE name = 'extreme_counterattack' AND tier = 1),     'instance', 'Unlocks Extreme Counterattack line'),
('Agility Booster Blueprint',           'module', (SELECT id FROM module_types WHERE name = 'agility_booster' AND tier = 1),           'instance', 'Unlocks Agility Booster line'),
('Infrared Scanner Blueprint',          'module', (SELECT id FROM module_types WHERE name = 'infrared_scanner' AND tier = 1),          'instance', 'Unlocks Infrared Scanner line'),
('ECM Booster Blueprint',               'module', (SELECT id FROM module_types WHERE name = 'ecm_booster' AND tier = 1),               'instance', 'Unlocks ECM Booster line'),
('Auto Target System Blueprint',        'module', (SELECT id FROM module_types WHERE name = 'auto_target_system' AND tier = 1),        'instance', 'Unlocks Auto Target System line'),
('Time Dilation Module Blueprint',      'module', (SELECT id FROM module_types WHERE name = 'time_dilation_module' AND tier = 1),      'instance', 'Unlocks Time Dilation Module line'),
('Station Warehouse Blueprint',         'module', (SELECT id FROM module_types WHERE name = 'station_warehouse' AND tier = 0),         'instance', 'Unlocks Station Warehouse'),
('Nano Station Warehouse Blueprint',    'module', (SELECT id FROM module_types WHERE name = 'nano_station_warehouse' AND tier = 0),    'instance', 'Unlocks Nano Station Warehouse'),
('Super Transmission Engine Blueprint', 'module', (SELECT id FROM module_types WHERE name = 'super_transmission_engine' AND tier = 0), 'quest', 'Unlocks Super Transmission Engine'),
('Team Combat Engine Blueprint',        'module', (SELECT id FROM module_types WHERE name = 'team_combat_engine' AND tier = 1),        'instance', 'Unlocks Team Combat Engine line'),
('Anti-Matter Engine Blueprint',        'module', (SELECT id FROM module_types WHERE name = 'anti_matter_engine' AND tier = 1),        'instance', 'Unlocks Anti-Matter Engine line');

-- ============================================================================
-- PART 9: SEED DATA - INSTANCES (30 Normal Instances)
-- ============================================================================

INSERT INTO instances (name, type, difficulty, required_level, max_fleets, ships_lost_on_defeat, he3_lost_on_defeat, exp_reward, enemy_fleets_json, rewards_json, description) VALUES
('Ancestral Recall',   'normal', 1,  1,  3,  true, true, 180,   '[]', '{"metal": 500, "he3": 500, "gold": 500}', 'First instance; tutorial-level challenge'),
('Deadzone',           'normal', 2,  3,  4,  true, true, 500,   '[]', '{"metal": 1000, "he3": 1000, "gold": 1000}', 'Early game instance'),
('Bravery',            'normal', 3,  5,  4,  true, true, 1000,  '[]', '{"metal": 1500, "he3": 1500, "gold": 1500}', 'First real challenge'),
('Dark Frontier',      'normal', 4,  8,  5,  true, true, 1800,  '[]', '{"metal": 2500, "he3": 2500, "gold": 2500}', 'Into the unknown'),
('Crimson Nebula',     'normal', 5,  10, 5,  true, true, 2800,  '[]', '{"metal": 3500, "he3": 3500, "gold": 3500}', 'Danger in the red mist'),
('Astral Rift',        'normal', 6,  13, 6,  true, true, 4000,  '[]', '{"metal": 5000, "he3": 5000, "gold": 5000}', 'Space-time distortion zone'),
('Void Passage',       'normal', 7,  15, 6,  true, true, 5500,  '[]', '{"metal": 6500, "he3": 6500, "gold": 6500}', 'Journey through the void'),
('Stellar Tempest',    'normal', 8,  18, 7,  true, true, 7200,  '[]', '{"metal": 8500, "he3": 8500, "gold": 8500}', 'Solar storm battlefield'),
('Iron Bastion',       'normal', 9,  20, 7,  true, true, 9200,  '[]', '{"metal": 11000, "he3": 11000, "gold": 11000}', 'Fortified enemy base'),
('Phantom Corridor',   'normal', 10, 23, 8,  true, true, 11500, '[]', '{"metal": 14000, "he3": 14000, "gold": 14000}', 'Ghost fleet territory'),
('Warp Anomaly',       'normal', 11, 25, 8,  true, true, 14000, '[]', '{"metal": 17000, "he3": 17000, "gold": 17000}', 'Unstable warp zone'),
('Shattered Core',     'normal', 12, 28, 9,  true, true, 16800, '[]', '{"metal": 20000, "he3": 20000, "gold": 20000}', 'Destroyed planet core'),
('Nova Storm',         'normal', 13, 30, 9,  true, true, 20000, '[]', '{"metal": 24000, "he3": 24000, "gold": 24000}', 'Supernova aftermath'),
('Eclipse Point',      'normal', 14, 33, 10, true, true, 23500, '[]', '{"metal": 28000, "he3": 28000, "gold": 28000}', 'Shadow of the dark star'),
('Binary Star',        'normal', 15, 35, 10, true, true, 27200, '[]', '{"metal": 32000, "he3": 32000, "gold": 32000}', 'Twin star system'),
('Gravity Well',       'normal', 16, 38, 11, true, true, 31200, '[]', '{"metal": 37000, "he3": 37000, "gold": 37000}', 'Gravitational trap'),
('Plasma Fields',      'normal', 17, 40, 11, true, true, 35500, '[]', '{"metal": 42000, "he3": 42000, "gold": 42000}', 'Supercharged plasma zone'),
('Ion Storm',          'normal', 18, 43, 12, true, true, 40000, '[]', '{"metal": 47000, "he3": 47000, "gold": 47000}', 'Electromagnetic chaos'),
('Cosmic Tides',       'normal', 19, 45, 12, true, true, 44800, '[]', '{"metal": 52000, "he3": 52000, "gold": 52000}', 'Cosmic currents'),
('Asteroid Belt',      'normal', 20, 48, 12, true, true, 49800, '[]', '{"metal": 58000, "he3": 58000, "gold": 58000}', 'Dense asteroid field'),
('Dark Matter',        'normal', 21, 50, 13, true, true, 52000, '[]', '{"metal": 64000, "he3": 64000, "gold": 64000}', 'Dark matter concentration'),
('Quantum Flux',       'normal', 22, 52, 13, true, true, 54500, '[]', '{"metal": 70000, "he3": 70000, "gold": 70000}', 'Quantum instability zone'),
('Supernova Rim',      'normal', 23, 54, 13, true, true, 57000, '[]', '{"metal": 76000, "he3": 76000, "gold": 76000}', 'Edge of the supernova'),
('Pulsar Gate',        'normal', 24, 56, 14, true, true, 59500, '[]', '{"metal": 83000, "he3": 83000, "gold": 83000}', 'Pulsar-powered gateway'),
('Neutron Field',      'normal', 25, 58, 14, true, true, 62000, '[]', '{"metal": 90000, "he3": 90000, "gold": 90000}', 'Neutron star radiation field'),
('Singularity',        'normal', 26, 60, 14, true, true, 64500, '[]', '{"metal": 97000, "he3": 97000, "gold": 97000}', 'Near the singularity'),
('Event Horizon',      'normal', 27, 62, 15, true, true, 67000, '[]', '{"metal": 105000, "he3": 105000, "gold": 105000}', 'Point of no return'),
('Omega Rift',         'normal', 28, 64, 15, true, true, 69500, '[]', '{"metal": 113000, "he3": 113000, "gold": 113000}', 'The final rift'),
('Final Frontier',     'normal', 29, 66, 15, true, true, 71500, '[]', '{"metal": 121000, "he3": 121000, "gold": 121000}', 'Edge of known space'),
('Triumphant Glory',   'normal', 30, 68, 15, true, true, 73500, '[]', '{"metal": 130000, "he3": 130000, "gold": 130000}', 'The ultimate challenge');

-- ============================================================================
-- PART 10: SEED DATA - INSTANCE BLUEPRINTS (blueprint drops per instance)
-- ============================================================================

-- Instance 1: Weikes, Rapid Fire
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Ancestral Recall' AND b.name = 'Weikes Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Ancestral Recall' AND b.name = 'Rapid Fire Blueprint';

-- Instance 2: Typhoon, Taskmaster
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Deadzone' AND b.name = 'Typhoon Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Deadzone' AND b.name = 'Taskmaster Blueprint';

-- Instance 3: Estrella, Cluster Laser Transmitter
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Bravery' AND b.name = 'Estrella Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Bravery' AND b.name = 'Cluster Laser Transmitter Blueprint';

-- Instance 4: Air Wanderer, Gatling Cannon
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Dark Frontier' AND b.name = 'Air Wanderer Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Dark Frontier' AND b.name = 'Gatling Cannon Blueprint';

-- Instance 5: Bombardier, Magneto Pulsar
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Crimson Nebula' AND b.name = 'Bombardier Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Crimson Nebula' AND b.name = 'Magneto Pulsar Blueprint';

-- Instance 6: Valkyrie, Rocket Frame
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Astral Rift' AND b.name = 'Valkyrie Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Astral Rift' AND b.name = 'Rocket Frame Blueprint';

-- Instance 7: Duke, Starlight Missile Pod
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Void Passage' AND b.name = 'Duke Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Void Passage' AND b.name = 'Starlight Missile Pod Blueprint';

-- Instance 8: GoGetter, Streamliner
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Stellar Tempest' AND b.name = 'GoGetter Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Stellar Tempest' AND b.name = 'Streamliner Blueprint';

-- Instance 9: The Shuttler, Golem
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Iron Bastion' AND b.name = 'The Shuttler Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Iron Bastion' AND b.name = 'Golem Blueprint';

-- Instance 10: Space Hunter, Particle Stun Shield
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Phantom Corridor' AND b.name = 'Space Hunter Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Phantom Corridor' AND b.name = 'Particle Stun Shield Blueprint';

-- Instance 11: Watchman, Heat Diffusion Shield
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Warp Anomaly' AND b.name = 'Watchman Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Warp Anomaly' AND b.name = 'Heat Diffusion Shield Blueprint';

-- Instance 12: Sparrow, Space-Time Magnetic Shield
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Shattered Core' AND b.name = 'Sparrow Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Shattered Core' AND b.name = 'Space-Time Magnetic Shield Blueprint';

-- Instance 13: Spinner, Detonator Shield
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Nova Storm' AND b.name = 'Spinner Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Nova Storm' AND b.name = 'Detonator Shield Blueprint';

-- Instance 14: Devourer, Anti-Aircraft Cannon
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Eclipse Point' AND b.name = 'Devourer Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Eclipse Point' AND b.name = 'Anti-Aircraft Cannon Blueprint';

-- Instance 15: Wraith, Powered Pulse Cannon
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Binary Star' AND b.name = 'Wraith Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Binary Star' AND b.name = 'Powered Pulse Cannon Blueprint';

-- Instance 16: Polymesus, Ship Reinforcement Facility
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Gravity Well' AND b.name = 'Polymesus Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Gravity Well' AND b.name = 'Ship Reinforcement Facility Blueprint';

-- Instance 17: Encratos, Quick Reaction Armor
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Plasma Fields' AND b.name = 'Encratos Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Plasma Fields' AND b.name = 'Quick Reaction Armor Blueprint';

-- Instance 18: Cybra, Reflective Plating
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Ion Storm' AND b.name = 'Cybra Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Ion Storm' AND b.name = 'Reflective Plating Blueprint';

-- Instance 19: Nicholas, Daedalus Control System
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Cosmic Tides' AND b.name = 'Nicholas Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Cosmic Tides' AND b.name = 'Daedalus Control System Blueprint';

-- Instance 20: Hamdar, Energy Armor
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Asteroid Belt' AND b.name = 'Hamdar Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Asteroid Belt' AND b.name = 'Energy Armor Blueprint';

-- Instance 21: Nettle, EOS Phase Shift Shield
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Dark Matter' AND b.name = 'Nettle Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Dark Matter' AND b.name = 'EOS Phase Shift Shield Blueprint';

-- Instance 22: Helena, Shield Regenerator
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Quantum Flux' AND b.name = 'Helena Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Quantum Flux' AND b.name = 'Shield Regenerator Blueprint';

-- Instance 23: Diaz, Energy Shield Booster
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Supernova Rim' AND b.name = 'Diaz Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Supernova Rim' AND b.name = 'Energy Shield Booster Blueprint';

-- Instance 24: RV766-The Explorer, Agility Booster
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Pulsar Gate' AND b.name = 'RV766-The Explorer Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Pulsar Gate' AND b.name = 'Agility Booster Blueprint';

-- Instance 25: Palenka, Infrared Scanner
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Neutron Field' AND b.name = 'Palenka Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Neutron Field' AND b.name = 'Infrared Scanner Blueprint';

-- Instance 26: ECM Booster, Auto Target System
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Singularity' AND b.name = 'ECM Booster Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Singularity' AND b.name = 'Auto Target System Blueprint';

-- Instance 27: Time Dilation Module, Extreme Counterattack
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Event Horizon' AND b.name = 'Time Dilation Module Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Event Horizon' AND b.name = 'Extreme Counterattack Blueprint';

-- Instance 28: Team Combat Engine, Anti-Matter Engine
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Omega Rift' AND b.name = 'Team Combat Engine Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Omega Rift' AND b.name = 'Anti-Matter Engine Blueprint';

-- Instance 29: Lander Module, Nano Station Warehouse
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Final Frontier' AND b.name = 'Lander Module Blueprint';
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Final Frontier' AND b.name = 'Nano Station Warehouse Blueprint';

-- Instance 30: ALL blueprints (equal chance)
-- Insert all blueprints as drops for Triumphant Glory
INSERT INTO instance_blueprints (instance_id, blueprint_id)
SELECT i.id, b.id FROM instances i, blueprints b
WHERE i.name = 'Triumphant Glory';
