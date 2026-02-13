-- ============================================================================
-- Phase 4 Corps & Galaxy Migration - Cryptomines Online
-- Tables: corps, corp_members, corp_donations, galaxy_zones, rbp_attacks
-- Seed data: System player, 49 RBP planets (7x7 grid), galaxy zones
-- ============================================================================

-- =========================
-- 1. CORPS
-- =========================
CREATE TABLE corps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    tag TEXT UNIQUE NOT NULL,
    leader_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    level INTEGER NOT NULL DEFAULT 1,
    wealth BIGINT NOT NULL DEFAULT 0,
    max_members INTEGER NOT NULL DEFAULT 20,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_tag_length CHECK (char_length(tag) >= 3 AND char_length(tag) <= 5),
    CONSTRAINT chk_level_positive CHECK (level >= 1),
    CONSTRAINT chk_wealth_non_negative CHECK (wealth >= 0),
    CONSTRAINT chk_max_members_positive CHECK (max_members > 0)
);

CREATE INDEX idx_corps_name ON corps (name);
CREATE INDEX idx_corps_tag ON corps (tag);
CREATE INDEX idx_corps_leader ON corps (leader_id);

-- =========================
-- 2. CORP MEMBERS
-- =========================
CREATE TABLE corp_members (
    corp_id UUID NOT NULL REFERENCES corps(id) ON DELETE CASCADE,
    player_id UUID NOT NULL UNIQUE REFERENCES players(id) ON DELETE CASCADE,
    role TEXT NOT NULL CHECK (role IN ('leader', 'officer', 'member')),
    contribution_points BIGINT NOT NULL DEFAULT 0,
    daily_contribution INTEGER NOT NULL DEFAULT 0,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (corp_id, player_id),
    CONSTRAINT chk_contribution_non_negative CHECK (contribution_points >= 0),
    CONSTRAINT chk_daily_contribution_non_negative CHECK (daily_contribution >= 0)
);

CREATE INDEX idx_corp_members_player ON corp_members (player_id);
CREATE INDEX idx_corp_members_role ON corp_members (corp_id, role);

-- =========================
-- 3. CORP DONATIONS
-- =========================
CREATE TABLE corp_donations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    corp_id UUID NOT NULL REFERENCES corps(id) ON DELETE CASCADE,
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    metal BIGINT NOT NULL DEFAULT 0,
    he3 BIGINT NOT NULL DEFAULT 0,
    gold BIGINT NOT NULL DEFAULT 0,
    contribution_points INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_resources_non_negative CHECK (metal >= 0 AND he3 >= 0 AND gold >= 0),
    CONSTRAINT chk_contribution_positive CHECK (contribution_points > 0)
);

CREATE INDEX idx_corp_donations_corp ON corp_donations (corp_id);
CREATE INDEX idx_corp_donations_player ON corp_donations (player_id);
CREATE INDEX idx_corp_donations_created ON corp_donations (created_at DESC);

-- =========================
-- 4. GALAXY ZONES
-- =========================
CREATE TABLE galaxy_zones (
    zone_x INTEGER NOT NULL,
    zone_y INTEGER NOT NULL,
    rbp_planet_id UUID NOT NULL REFERENCES planets(id) ON DELETE CASCADE,

    PRIMARY KEY (zone_x, zone_y),
    CONSTRAINT chk_zone_coords CHECK (zone_x >= 0 AND zone_y >= 0)
);

CREATE INDEX idx_galaxy_zones_rbp ON galaxy_zones (rbp_planet_id);

-- =========================
-- 5. RBP ATTACKS
-- =========================
CREATE TABLE rbp_attacks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rbp_planet_id UUID NOT NULL REFERENCES planets(id) ON DELETE CASCADE,
    attacking_corp_id UUID NOT NULL REFERENCES corps(id) ON DELETE CASCADE,
    combat_report_id UUID REFERENCES combat_reports(id) ON DELETE SET NULL,
    total_kills INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_total_kills_non_negative CHECK (total_kills >= 0)
);

CREATE INDEX idx_rbp_attacks_planet ON rbp_attacks (rbp_planet_id);
CREATE INDEX idx_rbp_attacks_corp ON rbp_attacks (attacking_corp_id);
CREATE INDEX idx_rbp_attacks_created ON rbp_attacks (created_at DESC);

-- =========================
-- 6. ALTER PLANETS - Add Foreign Key
-- =========================
-- The controlling_corp_id column already exists from phase1_mvp migration
-- We just need to add the foreign key constraint
ALTER TABLE planets
ADD CONSTRAINT fk_planets_controlling_corp
FOREIGN KEY (controlling_corp_id) REFERENCES corps(id) ON DELETE SET NULL;

-- =========================
-- 7. SEED DATA - System Player
-- =========================
-- Create a system player to own RBP planets
INSERT INTO players (
    id,
    anonymous_id,
    username,
    level,
    experience,
    tutorial_step,
    is_online,
    created_at,
    last_login
) VALUES (
    '00000000-0000-0000-0000-000000000000',
    'SYSTEM',
    'SYSTEM',
    999,
    0,
    999,
    false,
    now(),
    now()
) ON CONFLICT (id) DO NOTHING;

-- =========================
-- 8. SEED DATA - RBP Planets (7x7 Grid)
-- =========================
-- Create 49 RBP planets in a 7x7 grid (zone 0,0 to 6,6)
-- Each RBP planet is located at position (1000 + zone_x*10, 1000 + zone_y*10)
-- RBP level varies from 1-7 based on (zone_x + zone_y) % 7 + 1
-- All planets get 72-hour protection initially

DO $$
DECLARE
    system_player_id UUID := '00000000-0000-0000-0000-000000000000';
    zone_x INT;
    zone_y INT;
    pos_x INT;
    pos_y INT;
    rbp_lvl INT;
    planet_name TEXT;
    new_planet_id UUID;
BEGIN
    FOR zone_x IN 0..6 LOOP
        FOR zone_y IN 0..6 LOOP
            -- Calculate position (offset to avoid player planets)
            pos_x := 1000 + zone_x * 10;
            pos_y := 1000 + zone_y * 10;

            -- Calculate RBP level (1-7)
            rbp_lvl := ((zone_x + zone_y) % 7) + 1;

            -- Generate planet name
            planet_name := 'RBP-' || zone_x || zone_y;

            -- Insert RBP planet
            INSERT INTO planets (
                player_id,
                name,
                position_x,
                position_y,
                is_homeworld,
                is_rbp,
                rbp_level,
                protection_until
            ) VALUES (
                system_player_id,
                planet_name,
                pos_x,
                pos_y,
                false,
                true,
                rbp_lvl,
                now() + interval '72 hours'
            ) RETURNING id INTO new_planet_id;

            -- Create resources entry for the planet
            INSERT INTO resources (planet_id, metal, he3, gold)
            VALUES (new_planet_id, 0, 0, 0);

            -- Link zone to planet
            INSERT INTO galaxy_zones (zone_x, zone_y, rbp_planet_id)
            VALUES (zone_x, zone_y, new_planet_id);
        END LOOP;
    END LOOP;
END $$;
