-- Migration: Commander Types (30 Commanders)
-- Module 4: Commander System
-- Date: 2026-02-07
-- Priority: HIGH (Phase C: Military)

-- ============================================================================
-- TABLE: commander_types (Reference Table for Gacha System)
-- ============================================================================
-- Stores all commander templates for recruitment gacha
-- Players recruit commanders at Command Center using Gold (10,000 per recruitment)
-- Drop rates: Common 50%, Skill 35%, Super 15%

CREATE TABLE commander_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL, -- Internal key (e.g., 'rookie_pilot')
    display_name TEXT NOT NULL, -- Display name (e.g., 'Rookie Pilot')
    rarity TEXT NOT NULL CHECK (rarity IN ('common', 'skill', 'super')),

    -- Base stats (applied when commander is recruited)
    base_accuracy INTEGER NOT NULL DEFAULT 0 CHECK (base_accuracy >= 0),
    base_dodge INTEGER NOT NULL DEFAULT 0 CHECK (base_dodge >= 0),
    base_speed INTEGER NOT NULL DEFAULT 0 CHECK (base_speed >= 0),
    base_electron INTEGER NOT NULL DEFAULT 0 CHECK (base_electron >= 0),

    description TEXT NOT NULL DEFAULT '',
    avatar_url TEXT, -- Future: commander portrait images (Phase 3+)
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================================
-- INDEXES
-- ============================================================================
CREATE INDEX idx_commander_types_rarity ON commander_types (rarity);

-- ============================================================================
-- SEED DATA: 30 Commanders
-- ============================================================================

-- COMMON COMMANDERS (50% drop rate) - 15 commanders
-- Stat range: 10-17 per stat (total ~50-60)
-- Balanced, low-tier commanders
INSERT INTO commander_types (name, display_name, rarity, base_accuracy, base_dodge, base_speed, base_electron, description) VALUES
('rookie_pilot', 'Rookie Pilot', 'common', 10, 10, 10, 10, 'Inexperienced but eager commander with balanced stats'),
('veteran_soldier', 'Veteran Soldier', 'common', 15, 12, 11, 12, 'Battle-tested infantry officer with focus on accuracy'),
('junior_navigator', 'Junior Navigator', 'common', 12, 15, 13, 10, 'Skilled at evasive maneuvers'),
('tactical_officer', 'Tactical Officer', 'common', 14, 11, 12, 13, 'Strategic planner with electron focus'),
('cargo_captain', 'Cargo Captain', 'common', 11, 13, 14, 12, 'Transport specialist with speed advantage'),
('scout_commander', 'Scout Commander', 'common', 13, 16, 15, 11, 'Fast reconnaissance expert'),
('defensive_captain', 'Defensive Captain', 'common', 12, 17, 10, 11, 'Defensive specialist with high dodge'),
('assault_leader', 'Assault Leader', 'common', 16, 11, 13, 14, 'Offensive-focused commander'),
('supply_officer', 'Supply Officer', 'common', 11, 12, 11, 15, 'Resource management specialist'),
('patrol_chief', 'Patrol Chief', 'common', 14, 14, 12, 12, 'Balanced patrol leader'),
('training_instructor', 'Training Instructor', 'common', 13, 13, 13, 13, 'Perfectly balanced stats for training'),
('repair_specialist', 'Repair Specialist', 'common', 10, 14, 11, 16, 'Engineering focus with high electron'),
('communications_officer', 'Communications Officer', 'common', 12, 12, 16, 14, 'Speed specialist for quick response'),
('mining_foreman', 'Mining Foreman', 'common', 11, 11, 12, 12, 'Resource operations commander'),
('logistics_manager', 'Logistics Manager', 'common', 13, 12, 14, 13, 'Supply chain expert');

-- SKILL COMMANDERS (35% drop rate) - 10 commanders
-- Stat range: 18-28 per stat (total ~80-100)
-- Medium-tier, specialized commanders
INSERT INTO commander_types (name, display_name, rarity, base_accuracy, base_dodge, base_speed, base_electron, description) VALUES
('ace_pilot', 'Ace Pilot', 'skill', 22, 20, 21, 19, 'Elite fighter pilot with excellent all-around stats'),
('fleet_captain', 'Fleet Captain', 'skill', 21, 19, 20, 22, 'Experienced fleet commander with electron focus'),
('combat_veteran', 'Combat Veteran', 'skill', 25, 18, 19, 21, 'Accuracy specialist with extensive combat experience'),
('evasion_master', 'Evasion Master', 'skill', 19, 26, 22, 20, 'Dodge specialist - master of evasive tactics'),
('speed_demon', 'Speed Demon', 'skill', 20, 21, 27, 19, 'Speed specialist - unmatched agility'),
('tech_genius', 'Tech Genius', 'skill', 19, 20, 20, 28, 'Electron specialist - advanced technology expert'),
('strike_commander', 'Strike Commander', 'skill', 24, 19, 22, 23, 'Offensive expert with high accuracy'),
('shield_specialist', 'Shield Specialist', 'skill', 18, 25, 19, 24, 'Defensive expert with dodge and electron focus'),
('tactical_genius', 'Tactical Genius', 'skill', 23, 22, 21, 25, 'Balanced elite commander'),
('战术大师', '战术大师', 'skill', 22, 23, 23, 24, 'Tactical master (Chinese name variant for diversity)');

-- SUPER COMMANDERS (15% drop rate) - 5 commanders
-- Stat range: 31-41 per stat (total ~130-150)
-- High-tier, elite commanders
INSERT INTO commander_types (name, display_name, rarity, base_accuracy, base_dodge, base_speed, base_electron, description) VALUES
('grand_admiral', 'Grand Admiral', 'super', 35, 32, 33, 34, 'Supreme fleet commander with exceptional balanced stats'),
('legendary_ace', 'Legendary Ace', 'super', 38, 31, 35, 33, 'Unmatched accuracy - legendary pilot'),
('shadow_ghost', 'Shadow Ghost', 'super', 31, 39, 36, 32, 'Master of evasion - nearly impossible to hit'),
('lightning_strike', 'Lightning Strike', 'super', 33, 33, 40, 34, 'Unparalleled speed - strikes before the enemy reacts'),
('quantum_commander', 'Quantum Commander', 'super', 32, 34, 34, 41, 'Advanced technology expert with quantum computing mastery');

-- Total: 30 commanders (15 Common + 10 Skill + 5 Super)

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================
ALTER TABLE commander_types ENABLE ROW LEVEL SECURITY;

-- commander_types: Public read access (all players can see all commander types for gacha)
CREATE POLICY commander_types_select_all ON commander_types FOR SELECT USING (true);

-- ============================================================================
-- COMMENTS
-- ============================================================================
COMMENT ON TABLE commander_types IS 'Reference table for commander gacha system. 30 commander templates with base stats.';
COMMENT ON COLUMN commander_types.rarity IS 'Commander rarity: common (50%), skill (35%), super (15%). Legendary/Divine out of scope.';
COMMENT ON COLUMN commander_types.base_accuracy IS 'Base accuracy stat (hit rate bonus). Common: 10-17, Skill: 18-28, Super: 31-41.';
COMMENT ON COLUMN commander_types.base_dodge IS 'Base dodge stat (evasion bonus). Affects enemy hit rate.';
COMMENT ON COLUMN commander_types.base_speed IS 'Base speed stat (turn order). Higher speed acts first in combat.';
COMMENT ON COLUMN commander_types.base_electron IS 'Base electron stat (tech bonus). Affects special abilities and shields.';
