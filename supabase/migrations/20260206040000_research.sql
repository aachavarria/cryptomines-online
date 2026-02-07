-- Research System Migration
-- Tables: tech_types, technologies (player_research)
-- Seed Data: ~90 tech_types across 7 science trees

-- =============================================================================
-- 3.6 Tech Types (Reference Data)
-- =============================================================================
DROP TABLE IF EXISTS technologies CASCADE;
DROP TABLE IF EXISTS tech_types CASCADE;

CREATE TABLE tech_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    display_name TEXT NOT NULL,
    tree TEXT NOT NULL CHECK (tree IN (
        'logistics_construction', 'planetary_defense', 'ballistics_science',
        'directional_science', 'missile_science', 'ship_based_science', 'ship_defense_science'
    )),
    prerequisites_json JSONB NOT NULL DEFAULT '[]',
    -- Note: In GO2 all tech research costs Gold only (metal=0, he3=0).
    -- Schema retains metal/he3 columns for flexibility.
    base_cost_metal BIGINT NOT NULL DEFAULT 0,
    base_cost_he3 BIGINT NOT NULL DEFAULT 0,
    base_cost_gold BIGINT NOT NULL DEFAULT 0,
    cost_multiplier NUMERIC(6,4) NOT NULL DEFAULT 1.5300,
    base_time_seconds INTEGER NOT NULL DEFAULT 117,
    time_multiplier NUMERIC(6,4) NOT NULL DEFAULT 2.3400,
    max_level INTEGER NOT NULL DEFAULT 10,
    effects_json JSONB NOT NULL DEFAULT '{}',
    description TEXT NOT NULL DEFAULT ''
);

CREATE INDEX idx_tech_types_tree ON tech_types (tree);

-- =============================================================================
-- 3.7 Technologies (Player Research Progress)
-- =============================================================================
CREATE TABLE technologies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    tech_type INTEGER NOT NULL REFERENCES tech_types(id),
    level INTEGER NOT NULL DEFAULT 0,
    is_researching BOOLEAN NOT NULL DEFAULT false,
    research_finish_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_player_tech UNIQUE (player_id, tech_type),
    CONSTRAINT chk_level_non_negative CHECK (level >= 0),
    CONSTRAINT chk_research_consistency CHECK (
        (is_researching = true AND research_finish_at IS NOT NULL) OR
        (is_researching = false AND research_finish_at IS NULL)
    )
);

CREATE INDEX idx_technologies_player_id ON technologies (player_id);
CREATE INDEX idx_technologies_researching ON technologies (is_researching) WHERE is_researching = true;

-- =============================================================================
-- 3.28 Seed Data: Tech Types (Logistics Construction Science) - 11 techs
-- =============================================================================
INSERT INTO tech_types (name, display_name, tree, max_level, prerequisites_json, base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier, base_time_seconds, time_multiplier, effects_json, description) VALUES
('concurrent_construction', 'Concurrent Construction', 'logistics_construction', 1, '[]', 0, 0, 1000, 1.0000, 20, 1.0000, '{"type":"construction_slots","per_level":1}', 'Adds 1 construction slot'),
('construction_boost', 'Construction Boost', 'logistics_construction', 10, '[{"tech":"concurrent_construction","level":1}]', 0, 0, 2400, 1.5300, 480, 2.3400, '{"type":"build_speed","per_level":1.5,"unit":"percent"}', '+1-15% building construction speed'),
('quality_materials', 'Quality Materials', 'logistics_construction', 10, '[{"tech":"construction_boost","level":3}]', 0, 0, 1200, 1.5300, 240, 2.3400, '{"type":"build_cost_reduction","per_level":1.5,"unit":"percent"}', '-1-15% building resource costs'),
('ship_building_boost', 'Ship Building Boost', 'logistics_construction', 10, '[]', 0, 0, 990, 1.5300, 198, 2.3400, '{"type":"ship_build_speed","per_level":1.5,"unit":"percent"}', '+1-15% shipbuilding speed'),
('ship_building_logistics', 'Ship Building Logistics', 'logistics_construction', 10, '[{"tech":"ship_building_boost","level":2}]', 0, 0, 1386, 1.5300, 277, 2.3400, '{"type":"ship_build_cost_reduction","per_level":1.5,"unit":"percent"}', '-1-15% ship construction resource costs'),
('sync_shipbuilding', 'Sync Shipbuilding', 'logistics_construction', 1, '[{"tech":"ship_building_logistics","level":4}]', 0, 0, 174000, 1.0000, 34800, 1.0000, '{"type":"ship_production_slots","per_level":1}', 'Adds 1 shipbuilding slot (5th production slot)'),
('repair_technology', 'Repair Technology', 'logistics_construction', 10, '[{"tech":"sync_shipbuilding","level":1}]', 0, 0, 5310, 1.5300, 1062, 2.3400, '{"type":"ship_repair_percent","per_level":1,"unit":"percent"}', '+1-10% ship repair percentage'),
('high_yield_mining', 'High Yield Mining', 'logistics_construction', 10, '[]', 0, 0, 1740, 1.5300, 348, 2.3400, '{"type":"metal_output","per_level":1,"unit":"percent"}', '+1-10% Metal output'),
('high_yield_chemistry', 'High Yield Chemistry', 'logistics_construction', 10, '[{"tech":"high_yield_mining","level":2}]', 0, 0, 2400, 1.5300, 480, 2.3400, '{"type":"he3_output","per_level":1,"unit":"percent"}', '+1-10% He3 output'),
('high_yield_investing', 'High Yield Investing', 'logistics_construction', 10, '[{"tech":"high_yield_chemistry","level":2}]', 0, 0, 3570, 1.5300, 714, 2.3400, '{"type":"gold_output","per_level":1,"unit":"percent"}', '+1-10% Gold output'),
('expand_capacity', 'Expand Capacity', 'logistics_construction', 7, '[{"tech":"high_yield_investing","level":4}]', 0, 0, 3540, 1.5300, 708, 2.3400, '{"type":"warehouse_capacity","per_level":50000,"unit":"flat"}', '+50,000-350,000 warehouse storage per level');

-- =============================================================================
-- 3.29 Seed Data: Tech Types (Ballistics Science) - 14 techs
-- =============================================================================
INSERT INTO tech_types (name, display_name, tree, max_level, prerequisites_json, base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier, base_time_seconds, time_multiplier, effects_json, description) VALUES
('ballistics_base', 'Ballistics', 'ballistics_science', 10, '[]', 0, 0, 541, 1.5300, 117, 2.3400, '{"type":"ballistic_damage","per_level":5,"unit":"percent"}', '+5% ballistic damage per level'),
('ballistic_malice', 'Ballistic Malice', 'ballistics_science', 5, '[{"tech":"ballistics_base","level":3}]', 0, 0, 7558, 1.4300, 1428, 1.3400, '{"type":"ballistic_crit_rate","per_level":1,"unit":"percent"}', '+1% critical hit rate per level'),
('ballistic_crackdown', 'Ballistic Crackdown', 'ballistics_science', 2, '[{"tech":"ballistics_base","level":3}]', 0, 0, 29628, 1.1600, 4335, 1.3400, '{"type":"ballistic_crit_damage","per_level":10,"unit":"percent"}', '+10% critical damage per level'),
('steady_control', 'Steady Control Tech', 'ballistics_science', 5, '[{"tech":"ballistics_base","level":6},{"tech":"ballistic_malice","level":3}]', 0, 0, 26108, 1.4300, 3366, 1.3400, '{"type":"weapon_space_reduction","per_level":2,"unit":"percent"}', '-2% weapon space per level'),
('precise_ballistics', 'Precise Ballistics', 'ballistics_science', 5, '[{"tech":"ballistics_base","level":8},{"tech":"steady_control","level":3}]', 0, 0, 51100, 1.4300, 4080, 1.3400, '{"type":"ballistic_hit_rate","per_level":1,"unit":"percent"}', '+1% hit rate per level'),
('shield_penetration', 'Shield Penetration', 'ballistics_science', 1, '[{"tech":"ballistic_malice","level":5},{"tech":"ballistic_crackdown","level":1},{"tech":"precise_ballistics","level":1}]', 0, 0, 154616, 1.0000, 15300, 1.0000, '{"type":"shield_bypass","flat":15,"unit":"percent"}', '15% shield bypass damage'),
('depleted_uranium_bomb', 'Depleted Uranium Bomb', 'ballistics_science', 3, '[{"tech":"ballistic_crackdown","level":2},{"tech":"shield_penetration","level":1}]', 0, 0, 203249, 1.3300, 12750, 1.3400, '{"type":"armor_bonus","neutral":10,"light":1,"unit":"percent_per_level"}', '+10-30% vs Neutral armor, +1-3% vs Light armor'),
('fire_bomb_research', 'Fire Bomb Research', 'ballistics_science', 3, '[{"tech":"ballistic_crackdown","level":2},{"tech":"shield_penetration","level":1}]', 0, 0, 203249, 1.3300, 12750, 1.3400, '{"type":"armor_bonus","regen":10,"light":1,"unit":"percent_per_level"}', '+10-30% vs Regen armor, +1-3% vs Light armor'),
('improved_penetration', 'Improved Penetration', 'ballistics_science', 3, '[{"tech":"depleted_uranium_bomb","level":1},{"tech":"fire_bomb_research","level":1}]', 0, 0, 484776, 1.3300, 22950, 1.3400, '{"type":"shield_pen_chance","per_level":3,"light_bonus":1,"unit":"percent"}', '+1-3% vs Light armor, 3-10% shield pen chance'),
('victory_rush', 'Victory Rush', 'ballistics_science', 1, '[{"tech":"depleted_uranium_bomb","level":3},{"tech":"fire_bomb_research","level":3},{"tech":"improved_penetration","level":3}]', 0, 0, 1918521, 1.0000, 224400, 1.0000, '{"type":"range_damage","ranges":[220,180,150,120],"crit_rate":5,"crit_damage":5}', 'Range damage: 220%/180%/150%/120%, +5% crit rate/damage'),
('ballistic_scattering', 'Ballistic Scattering', 'ballistics_science', 5, '[{"tech":"ballistics_base","level":10},{"tech":"ballistic_malice","level":5},{"tech":"steady_control","level":5},{"tech":"precise_ballistics","level":3}]', 0, 0, 119765, 1.4300, 9180, 1.3400, '{"type":"scatter_damage","per_level":5,"unit":"percent"}', '5-25% scatter damage to adjacent ships'),
('improved_ballistic_scattering', 'Improved Ballistic Scattering', 'ballistics_science', 3, '[{"tech":"precise_ballistics","level":5},{"tech":"ballistic_scattering","level":3}]', 0, 0, 240487, 1.3300, 17850, 1.3400, '{"type":"scatter_rate","per_level":8,"unit":"percent"}', '+8-25% scattering rate'),
('hop_bomb_research', 'Hop Bomb Research', 'ballistics_science', 5, '[{"tech":"ballistic_scattering","level":5},{"tech":"improved_ballistic_scattering","level":2}]', 0, 0, 349930, 1.4300, 20400, 1.3400, '{"type":"scatter_weapon_chance","per_level":3,"unit":"percent"}', '3-15% chance to deal 100% weapon damage as scatter');

-- =============================================================================
-- 3.30 Seed Data: Tech Types (Ship Defense Science) - 20 techs
-- =============================================================================
INSERT INTO tech_types (name, display_name, tree, max_level, prerequisites_json, base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier, base_time_seconds, time_multiplier, effects_json, description) VALUES
('ship_defense_base', 'Ship Defense Tech', 'ship_defense_science', 2, '[]', 0, 0, 2846, 1.1600, 3240, 1.3400, '{"type":"base_defense_stats","shield":2,"structure":2,"agility":2,"defense":2,"stability":5,"unit":"percent_per_level"}', '+2-5% base shield/structure/agility/defense, +5-10% stability'),
-- Shield branch
('shield_research', 'Shield Research', 'ship_defense_science', 5, '[{"tech":"ship_defense_base","level":1}]', 0, 0, 10392, 1.4300, 2700, 1.3400, '{"type":"base_shield","per_level":1,"unit":"percent"}', '+1-5% base shield'),
('energy_diffusion', 'Energy Diffusion Tech', 'ship_defense_science', 3, '[{"tech":"ship_defense_base","level":2},{"tech":"shield_research","level":3}]', 0, 0, 46765, 1.3300, 10800, 1.3400, '{"type":"shield_damage_reduction","per_level":1,"unit":"flat"}', 'Each shield module reduces damage by 1-3 points'),
('penetration_resistance', 'Penetration Resistance', 'ship_defense_science', 2, '[{"tech":"shield_research","level":5},{"tech":"energy_diffusion","level":1}]', 0, 0, 117576, 1.1600, 21600, 1.3400, '{"type":"shield_pen_resist","per_level":3,"unit":"percent"}', '-3-7% enemy shield penetration chance'),
('augment_shield', 'Augment Shield', 'ship_defense_science', 3, '[{"tech":"energy_diffusion","level":2},{"tech":"penetration_resistance","level":1}]', 0, 0, 248747, 1.3300, 39600, 1.3400, '{"type":"base_shield_bonus","values":[6,12,20],"unit":"percent"}', '+6-20% base shield'),
('restoration', 'Restoration', 'ship_defense_science', 2, '[{"tech":"energy_diffusion","level":3},{"tech":"augment_shield","level":2}]', 0, 0, 540000, 1.1600, 90000, 1.3400, '{"type":"shield_restore","per_level":30,"intercept":1,"unit":"percent"}', '+30-60% shield restore per round, +1-2% interception'),
('augment_absorption', 'Augment Absorption', 'ship_defense_science', 2, '[{"tech":"shield_research","level":5},{"tech":"energy_diffusion","level":1}]', 0, 0, 117576, 1.1600, 21600, 1.3400, '{"type":"shield_flat_reduction","per_level":2,"unit":"flat"}', 'Shield modules reduce damage by 2-5 points'),
('energy_conservation_defense', 'Energy Conservation', 'ship_defense_science', 3, '[{"tech":"energy_diffusion","level":2},{"tech":"augment_absorption","level":1}]', 0, 0, 248747, 1.3300, 39600, 1.3400, '{"type":"absorb_no_he3","per_level":3,"unit":"percent"}', '+3-10% chance absorb damage without He3'),
('electronic_barrier', 'Electronic Barrier', 'ship_defense_science', 2, '[{"tech":"energy_diffusion","level":3},{"tech":"energy_conservation_defense","level":2}]', 0, 0, 540000, 1.1600, 90000, 1.3400, '{"type":"reflect_damage","per_level":5,"unit":"percent"}', 'Reflect 5-10% damage before shields drop to 0'),
('damage_mitigation', 'Damage Mitigation', 'ship_defense_science', 3, '[{"tech":"augment_shield","level":3},{"tech":"restoration","level":2},{"tech":"energy_conservation_defense","level":3},{"tech":"electronic_barrier","level":2}]', 0, 0, 805152, 1.3300, 108000, 1.3400, '{"type":"absorb_double","per_level":10,"collateral_reduction":15,"unit":"percent"}', '10-30% absorb double damage, 15-45% lower collateral'),
-- Structure branch
('ship_structural_analysis', 'Ship Structural Analysis', 'ship_defense_science', 5, '[{"tech":"ship_defense_base","level":1}]', 0, 0, 10392, 1.4300, 2700, 1.3400, '{"type":"base_structure","per_level":1,"unit":"percent"}', '+1-5% base structure'),
('ship_reinforcement', 'Ship Reinforcement Tech', 'ship_defense_science', 3, '[{"tech":"ship_defense_base","level":2},{"tech":"ship_structural_analysis","level":3}]', 0, 0, 46765, 1.3300, 10800, 1.3400, '{"type":"structure_damage_reduction","per_level":1,"unit":"flat"}', 'Each structure module reduces damage by 1-3'),
('resilience', 'Resilience', 'ship_defense_science', 2, '[{"tech":"ship_structural_analysis","level":5},{"tech":"ship_reinforcement","level":1}]', 0, 0, 117576, 1.1600, 21600, 1.3400, '{"type":"structure_pen_resist","per_level":3,"unit":"percent"}', '-3-7% enemy structure penetration chance'),
('structure_improvement', 'Structure Improvement', 'ship_defense_science', 3, '[{"tech":"ship_reinforcement","level":2},{"tech":"resilience","level":1}]', 0, 0, 248747, 1.3300, 39600, 1.3400, '{"type":"base_structure_bonus","values":[6,12,20],"unit":"percent"}', '+6-20% base structure'),
('fast_repair', 'Fast Repair', 'ship_defense_science', 2, '[{"tech":"ship_reinforcement","level":3},{"tech":"structure_improvement","level":2}]', 0, 0, 540000, 1.1600, 90000, 1.3400, '{"type":"structure_restore","per_level":30,"unit":"percent"}', '+30-60% structure restore per round'),
('reaction_armor_improvement', 'Reaction Armor Improvement', 'ship_defense_science', 2, '[{"tech":"ship_structural_analysis","level":5},{"tech":"ship_reinforcement","level":1}]', 0, 0, 117576, 1.1600, 21600, 1.3400, '{"type":"structure_flat_reduction","per_level":2,"unit":"flat"}', 'Structure modules reduce damage by 2-5'),
('defense_improvement', 'Defense Improvement', 'ship_defense_science', 3, '[{"tech":"ship_reinforcement","level":2},{"tech":"reaction_armor_improvement","level":1}]', 0, 0, 248747, 1.3300, 39600, 1.3400, '{"type":"absorb_no_he3_structure","per_level":3,"unit":"percent"}', '+3-10% absorb damage without He3'),
('reflection_mastery', 'Reflection Mastery', 'ship_defense_science', 2, '[{"tech":"ship_reinforcement","level":3},{"tech":"defense_improvement","level":2}]', 0, 0, 540000, 1.1600, 90000, 1.3400, '{"type":"reflect_structure_damage","per_level":5,"unit":"percent"}', 'Reflect 5-10% damage before structure drops to 0'),
('stability_mastery', 'Stability Mastery', 'ship_defense_science', 3, '[{"tech":"structure_improvement","level":3},{"tech":"fast_repair","level":2},{"tech":"defense_improvement","level":3},{"tech":"reflection_mastery","level":2}]', 0, 0, 805152, 1.3300, 108000, 1.3400, '{"type":"absorb_double_structure","per_level":10,"collateral_reduction":15,"unit":"percent"}', '10-30% absorb double damage, 15-45% lower collateral');

-- =============================================================================
-- 3.31 Seed Data: Tech Types (Directional Science) - 15 techs
-- =============================================================================
INSERT INTO tech_types (name, display_name, tree, max_level, prerequisites_json, base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier, base_time_seconds, time_multiplier, effects_json, description) VALUES
('optics_base', 'Optics', 'directional_science', 10, '[]', 0, 0, 541, 1.5300, 117, 2.3400, '{"type":"directional_damage","per_level":5,"unit":"percent"}', '+5% directional weapon damage per level'),
('directional_malice', 'Directional Malice', 'directional_science', 5, '[{"tech":"optics_base","level":3}]', 0, 0, 7558, 1.4300, 1428, 1.3400, '{"type":"directional_crit_rate","per_level":1,"unit":"percent"}', '+1% critical strike rate per level'),
('directional_accuracy', 'Directional Accuracy', 'directional_science', 5, '[{"tech":"optics_base","level":3}]', 0, 0, 7558, 1.4300, 1428, 1.3400, '{"type":"directional_accuracy","per_level":1,"unit":"percent"}', '+1% accuracy per level'),
('eagle_eye', 'Eagle Eye', 'directional_science', 2, '[{"tech":"optics_base","level":6},{"tech":"directional_accuracy","level":1}]', 0, 0, 22724, 1.1600, 2550, 1.3400, '{"type":"steering_power","per_level":10,"unit":"percent"}', '+10% steering power per level'),
('energy_penetration', 'Energy Penetration', 'directional_science', 1, '[{"tech":"directional_accuracy","level":2},{"tech":"eagle_eye","level":1}]', 0, 0, 62584, 1.0000, 6120, 1.0000, '{"type":"hit_and_shield_pen","hit":8,"shield_pen":8,"unit":"percent"}', '+8% hit rate AND 8% shield penetration chance'),
('pierce', 'Pierce', 'directional_science', 5, '[{"tech":"optics_base","level":10},{"tech":"directional_malice","level":5},{"tech":"directional_accuracy","level":5},{"tech":"energy_penetration","level":1}]', 0, 0, 87919, 1.4300, 4947, 1.3400, '{"type":"piercing_damage","per_level":3,"unit":"percent"}', 'Piercing damage through target rows, +3% per level'),
('radiative_interference', 'Radiative Interference', 'directional_science', 5, '[{"tech":"optics_base","level":6},{"tech":"directional_accuracy","level":2}]', 0, 0, 43116, 1.4300, 9180, 1.3400, '{"type":"enemy_hit_reduction","per_level":2,"unit":"percent"}', 'Reduces enemy hit rate by 2% per level'),
('improved_pierce', 'Improved Pierce', 'directional_science', 2, '[{"tech":"pierce","level":3},{"tech":"eagle_eye","level":2}]', 0, 0, 240487, 1.1600, 17850, 1.3400, '{"type":"piercing_damage_bonus","per_level":5,"unit":"percent"}', '+5% additional piercing damage per level'),
('energy_accumulation', 'Energy Accumulation', 'directional_science', 3, '[{"tech":"eagle_eye","level":2},{"tech":"energy_penetration","level":1}]', 0, 0, 136845, 1.3300, 11985, 1.3400, '{"type":"movement_crit_bonus","per_level":2,"unit":"percent"}', '+2-6% critical bonus when fleet moves'),
('electronic_interference', 'Electronic Interference', 'directional_science', 3, '[{"tech":"directional_accuracy","level":2},{"tech":"radiative_interference","level":3}]', 0, 0, 87581, 1.3300, 11985, 1.3400, '{"type":"steering_reduction","chance":5,"per_level":3,"unit":"percent"}', '5% chance reduce enemy steering by 10%; -3-9% weapon space'),
('piercing_crit', 'Piercing Crit', 'directional_science', 1, '[{"tech":"pierce","level":5},{"tech":"improved_pierce","level":1}]', 0, 0, 391234, 1.0000, 25500, 1.0000, '{"type":"piercing_critical","enabled":true}', 'Enables critical piercing vs horizontally-aligned ships'),
('weakness_detection', 'Weakness Detection', 'directional_science', 3, '[{"tech":"energy_accumulation","level":3},{"tech":"electronic_interference","level":3}]', 0, 0, 349054, 1.3300, 20298, 1.3400, '{"type":"base_damage_and_accuracy","damage":5,"accuracy":3,"unit":"percent_per_level"}', '+5% base damage per level; 3-10% bonus accuracy'),
('particle_impact', 'Particle Impact Tech', 'directional_science', 3, '[{"tech":"improved_pierce","level":2},{"tech":"piercing_crit","level":1}]', 0, 0, 646368, 1.3300, 40800, 1.3400, '{"type":"enemy_attack_reduction","per_level":10,"duration":2,"unit":"percent"}', 'Reduces enemy attack power by 10-30%, lasts 2 rounds'),
('magnetic_impact', 'Magnetic Impact', 'directional_science', 3, '[{"tech":"directional_accuracy","level":5},{"tech":"radiative_interference","level":5},{"tech":"weakness_detection","level":2}]', 0, 0, 996120, 1.3300, 96900, 1.3400, '{"type":"ignore_agility","per_level":8,"chance":10,"unit":"percent"}', 'Ignores 8-25% enemy agility with 10-30% chance'),
('dynamic_impairment', 'Dynamic Impairment', 'directional_science', 3, '[{"tech":"particle_impact","level":3},{"tech":"weakness_detection","level":3},{"tech":"magnetic_impact","level":3}]', 0, 0, 1431605, 1.3300, 124950, 1.3400, '{"type":"ignore_defense","per_level":8,"mov_reduction":1,"chance":6,"unit":"percent"}', 'Ignores 8-25% defense; reduces movement 1-3');

-- =============================================================================
-- 3.31 Seed Data: Tech Types (Missile Science) - 14 techs
-- =============================================================================
INSERT INTO tech_types (name, display_name, tree, max_level, prerequisites_json, base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier, base_time_seconds, time_multiplier, effects_json, description) VALUES
('missile_theory', 'Missile Theory', 'missile_science', 10, '[]', 0, 0, 541, 1.5300, 117, 2.3400, '{"type":"missile_damage","per_level":4,"unit":"percent"}', '+4% missile damage per level'),
('missile_accuracy', 'Missile Accuracy', 'missile_science', 5, '[{"tech":"missile_theory","level":3}]', 0, 0, 7558, 1.4300, 1428, 1.3400, '{"type":"missile_hit_rate","per_level":2,"unit":"percent"}', '+2% hit rate per level'),
('cruise_dynamics', 'Cruise Dynamics', 'missile_science', 2, '[{"tech":"missile_theory","level":3}]', 0, 0, 11603, 1.1600, 3366, 1.3400, '{"type":"steering_power","per_level":10,"unit":"percent"}', '+10% steering power per level'),
('missile_research', 'Missile Research', 'missile_science', 3, '[{"tech":"missile_theory","level":6},{"tech":"cruise_dynamics","level":1}]', 0, 0, 38831, 1.3300, 7446, 1.3400, '{"type":"missile_damage_and_pen","damage":3,"pen":1,"unit":"percent_per_level"}', '+3-10% base damage, +1-5% shield pen chance'),
('missile_elusion', 'Missile Elusion', 'missile_science', 3, '[{"tech":"cruise_dynamics","level":2},{"tech":"missile_research","level":1}]', 0, 0, 78750, 1.3300, 9690, 1.3400, '{"type":"intercept_reduction","per_level":3,"hit_bonus":3,"unit":"percent"}', '-3-10% interception rate, +3-9% hit rate'),
('missile_space_optimization', 'Missile Space Optimization', 'missile_science', 4, '[{"tech":"missile_research","level":2},{"tech":"missile_elusion","level":2}]', 0, 0, 233515, 1.3300, 16830, 1.3400, '{"type":"weapon_space_reduction","per_level":5,"unit":"percent"}', '-5-20% weapon space per level'),
('multidirectional_assault', 'Multidirectional Assault', 'missile_science', 5, '[{"tech":"missile_theory","level":10},{"tech":"missile_accuracy","level":3},{"tech":"cruise_dynamics","level":2},{"tech":"missile_research","level":3}]', 0, 0, 80796, 1.4300, 10200, 1.3400, '{"type":"scatter_all","per_level":6,"unit":"percent"}', 'Scatters 6-30% damage across ALL enemy ships'),
('nuclear_radiation', 'Nuclear Radiation Research', 'missile_science', 5, '[{"tech":"missile_accuracy","level":4},{"tech":"multidirectional_assault","level":1}]', 0, 0, 112916, 1.4300, 8160, 1.3400, '{"type":"damage_taken_increase","per_level":2,"chance":4,"unit":"percent"}', '2-10% increase to target damage taken, 4-20% chance'),
('break_armor', 'Break Armor', 'missile_science', 1, '[{"tech":"multidirectional_assault","level":2},{"tech":"nuclear_radiation","level":3}]', 0, 0, 287437, 1.0000, 25500, 1.0000, '{"type":"armor_damage_bonus","flat":5,"unit":"percent"}', '+5% damage to all armor types'),
('energy_conservation_missile', 'Energy Conservation', 'missile_science', 4, '[{"tech":"multidirectional_assault","level":3},{"tech":"nuclear_radiation","level":5},{"tech":"break_armor","level":1}]', 0, 0, 428575, 1.3300, 30600, 1.3400, '{"type":"he3_cost_reduction","per_level":4,"chance_for_half":true,"unit":"percent"}', '4-18% chance to reduce He3 cost by 50%'),
('shrapnel_research', 'Shrapnel Research', 'missile_science', 2, '[{"tech":"missile_accuracy","level":5},{"tech":"multidirectional_assault","level":3}]', 0, 0, 157172, 1.1600, 15810, 1.3400, '{"type":"scatter_bonus","per_level":2,"unit":"percent"}', '+2-4% more scattering damage'),
('exaltation', 'Exaltation', 'missile_science', 4, '[{"tech":"multidirectional_assault","level":5},{"tech":"shrapnel_research","level":1}]', 0, 0, 224495, 1.3300, 15555, 1.3400, '{"type":"scatter_vs_low_structure","per_level":20,"he3_reduction":2,"unit":"percent"}', '+20-80% scatter vs lower-structure fleets'),
('suppression', 'Suppression', 'missile_science', 4, '[{"tech":"shrapnel_research","level":2},{"tech":"exaltation","level":4}]', 0, 0, 484776, 1.3300, 22950, 1.3400, '{"type":"scatter_vs_high_structure","per_level":3,"crit":3,"unit":"percent"}', '+3-12% scatter vs higher-structure fleets, +3-12% crit'),
('missile_concussion', 'Missile Concussion', 'missile_science', 1, '[{"tech":"energy_conservation_missile","level":2},{"tech":"missile_research","level":3},{"tech":"missile_elusion","level":3},{"tech":"missile_space_optimization","level":4}]', 0, 0, 1351976, 1.0000, 178500, 1.0000, '{"type":"knockback","chance":25,"distance":2,"per_round":true}', '25% chance to push target back 2 spaces');

-- =============================================================================
-- 3.31 Seed Data: Tech Types (Ship-Based Science) - 10 core techs
-- =============================================================================
INSERT INTO tech_types (name, display_name, tree, max_level, prerequisites_json, base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier, base_time_seconds, time_multiplier, effects_json, description) VALUES
('fighter_weapons_theory', 'Fighter Weapons Theory', 'ship_based_science', 10, '[]', 0, 0, 541, 1.5300, 117, 2.3400, '{"type":"fighter_damage","per_level":3,"unit":"percent"}', '+3% fighter weapon damage per level'),
('reconnaissance', 'Reconnaissance', 'ship_based_science', 2, '[{"tech":"fighter_weapons_theory","level":3}]', 0, 0, 7558, 1.1600, 1428, 1.3400, '{"type":"steering_power","per_level":10,"unit":"percent"}', '+10% steering power per level'),
('thruster_optimization', 'Thruster Optimization', 'ship_based_science', 5, '[{"tech":"fighter_weapons_theory","level":3}]', 0, 0, 7558, 1.4300, 1428, 1.3400, '{"type":"intercept_reduction","per_level":1,"unit":"percent"}', '-1% intercept rate per level'),
('navigation', 'Navigation', 'ship_based_science', 5, '[{"tech":"fighter_weapons_theory","level":3}]', 0, 0, 7558, 1.4300, 1428, 1.3400, '{"type":"fighter_hit_rate","per_level":1,"unit":"percent"}', '+1% hit rate per level'),
('fuel_optimization', 'Fuel Optimization', 'ship_based_science', 5, '[{"tech":"fighter_weapons_theory","level":6},{"tech":"reconnaissance","level":1}]', 0, 0, 26108, 1.4300, 3366, 1.3400, '{"type":"he3_cost_reduction","per_level":1,"unit":"percent"}', '-1-5% He3 costs'),
('fighter_mastery', 'Fighter Mastery', 'ship_based_science', 1, '[{"tech":"fighter_weapons_theory","level":6},{"tech":"navigation","level":3}]', 0, 0, 51100, 1.0000, 4080, 1.0000, '{"type":"base_attack","flat":5,"unit":"percent"}', '+5% base attack power'),
('fighter_tech_upgrades', 'Fighter Tech Upgrades', 'ship_based_science', 3, '[{"tech":"reconnaissance","level":2},{"tech":"fuel_optimization","level":5}]', 0, 0, 136845, 1.3300, 11985, 1.3400, '{"type":"he3_and_shield_damage","he3":2,"shield_damage":3,"unit":"percent_per_level"}', '-2-6% He3; +3-10% damage vs shielded enemies'),
('armor_structural_analysis', 'Armor Structural Analysis', 'ship_based_science', 1, '[{"tech":"fighter_tech_upgrades","level":3}]', 0, 0, 391234, 1.0000, 25500, 1.0000, '{"type":"ship_shield_damage","flat":10,"unit":"percent"}', '+10% damage vs ships and shields'),
('fighter_weapons_efficiency', 'Fighter-based Weapons Efficiency', 'ship_based_science', 3, '[{"tech":"fighter_weapons_theory","level":10},{"tech":"armor_structural_analysis","level":1}]', 0, 0, 349054, 1.3300, 20298, 1.3400, '{"type":"reload_chance","per_level":10,"unit":"percent"}', '+10% chance per level to finish reloading after attacks'),
('ingenuity', 'Ingenuity', 'ship_based_science', 1, '[{"tech":"fighter_weapons_efficiency","level":3}]', 0, 0, 2045150, 1.0000, 255000, 1.0000, '{"type":"multi_bonus","swarm":5,"attack":5,"he3":-5,"intercept":20,"unshielded":10}', '+5% swarm/attack, -5% He3, +20% intercept, +10% unshielded dmg');

-- =============================================================================
-- 3.31 Seed Data: Tech Types (Planetary Defense Science) - 8 techs
-- =============================================================================
INSERT INTO tech_types (name, display_name, tree, max_level, prerequisites_json, base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier, base_time_seconds, time_multiplier, effects_json, description) VALUES
('energy_control', 'Energy Control', 'planetary_defense', 10, '[]', 0, 0, 2500, 2.0000, 500, 2.0000, '{"type":"defense_cost_reduction","per_level":1,"unit":"percent"}', '-1-10% resource costs for defensive structures'),
('rapid_defense_buildup', 'Rapid Defense Buildup', 'planetary_defense', 10, '[{"tech":"energy_control","level":1}]', 0, 0, 3000, 2.0000, 600, 2.0000, '{"type":"defense_build_speed","per_level":1,"unit":"percent"}', '+1-10% defense construction speed'),
('defense_enhancement', 'Defense Enhancement', 'planetary_defense', 10, '[{"tech":"energy_control","level":3},{"tech":"rapid_defense_buildup","level":3}]', 0, 0, 1500, 2.0000, 300, 2.0000, '{"type":"defense_value","per_level":1,"unit":"percent"}', '+1-10% defensive value of all structures'),
('emplacement_mastery', 'Emplacement Mastery', 'planetary_defense', 10, '[{"tech":"rapid_defense_buildup","level":5}]', 0, 0, 2250, 2.0000, 450, 2.0000, '{"type":"emplacement_attack","per_level":1,"unit":"percent"}', '+1-10% emplacement attack power'),
('utmost_defense_buildup', 'Utmost Defense Buildup', 'planetary_defense', 10, '[{"tech":"energy_control","level":5},{"tech":"defense_enhancement","level":5},{"tech":"emplacement_mastery","level":3}]', 0, 0, 1750, 2.0000, 350, 2.0000, '{"type":"max_defense_structures","per_level":1,"unit":"percent"}', '+1-10% max defensive structures'),
('range_extension', 'Range Extension', 'planetary_defense', 2, '[{"tech":"rapid_defense_buildup","level":8},{"tech":"emplacement_mastery","level":5}]', 0, 0, 300000, 2.0000, 60000, 2.0000, '{"type":"defense_range","applies_to":["particle_cannon","anti_aircraft_gun"]}', 'Increases attack range of Particle Cannons and Anti-Aircraft Guns'),
('thor_buildup', 'Thor Buildup', 'planetary_defense', 1, '[{"tech":"energy_control","level":8},{"tech":"utmost_defense_buildup","level":5},{"tech":"range_extension","level":1}]', 0, 0, 800000, 1.0000, 160000, 1.0000, '{"type":"max_thor_cannon","flat":1}', '+1 max Thor''s Cannon allowed'),
('augment_propulsion', 'Augment Propulsion', 'planetary_defense', 2, '[{"tech":"emplacement_mastery","level":8},{"tech":"thor_buildup","level":1}]', 0, 0, 400000, 2.2500, 80000, 2.2500, '{"type":"defense_movement","per_level":1,"unit":"flat"}', '+1-2 ship movement speed when defending own planet');
