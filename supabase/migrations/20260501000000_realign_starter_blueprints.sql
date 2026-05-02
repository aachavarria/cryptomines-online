-- =============================================================================
-- Realign starter blueprint rewards in main quest chain to GDD Section 8.12.1
--
-- Why: Quest 10 (Ship Factory) gave Typhoon hull but no offensive module, so
-- new players could design a ship hull but had nothing to attack with until
-- Quest 17 (Weapon Research Center) which unlocked the first weapon (Starlight
-- Missile Pod). This left a 7-quest gap with no playable combat capability.
--
-- The original Galaxy Online 2 tutorial paired the first hull with the first
-- weapon at Ship Factory construction (Weikes + Rapid Fire). This migration
-- restores that combo and redistributes the remaining starter blueprints so
-- that each major military milestone unlocks a coherent set:
--   - Q10 Ship Factory       -> Weikes hull + Rapid Fire weapon (combat MVP)
--   - Q11 Command Center     -> Typhoon cruiser
--   - Q14 Build 100 ships    -> Estrella battleship
--   - Q17 Weapon Research    -> Energy Shield + Anti-Aircraft + Starlight Missile
--
-- Phase 1 also gets a small reshuffle so the early "use blueprint" loop has
-- a structure module (Ship Reinforcement Facility) instead of the Estrella
-- battleship hull which is too heavy for a brand-new player.
-- =============================================================================

-- Quest 7 (He3 Production): Estrella -> Ship Reinforcement Facility
UPDATE quest_types
SET reward_item_json = '[{"type":"blueprint","blueprint_key":"ship_reinforcement_facility"}]'
WHERE quest_key = 'main_07_he3_production';

-- Quest 8 (Blueprints 2): retarget the "use_blueprint" requirement to
-- ship_reinforcement_facility (was estrella) and award Atomic Framework BP
UPDATE quest_types
SET requirement_target = 'ship_reinforcement_facility',
    reward_item_json   = '[{"type":"blueprint","blueprint_key":"atomic_framework"}]'
WHERE quest_key = 'main_08_blueprints_2';

-- Quest 10 (Ship Factory): Typhoon -> Weikes hull + Rapid Fire weapon
UPDATE quest_types
SET reward_item_json = '[{"type":"blueprint","blueprint_key":"weikes"},{"type":"blueprint","blueprint_key":"rapid_fire"}]'
WHERE quest_key = 'main_10_ship_factory';

-- Quest 11 (Command Center): Energy Shield Booster -> Typhoon cruiser
UPDATE quest_types
SET reward_item_json = '[{"type":"blueprint","blueprint_key":"typhoon"}]'
WHERE quest_key = 'main_11_command_center';

-- Quest 14 (Build Ships): Anti-Aircraft Cannon -> Estrella battleship
UPDATE quest_types
SET reward_item_json = '[{"type":"blueprint","blueprint_key":"estrella"}]'
WHERE quest_key = 'main_14_build_ships';

-- Quest 17 (Weapon Research Center): Starlight Missile only -> military bundle
-- (Energy Shield Booster + Anti-Aircraft Cannon + Starlight Missile Pod)
UPDATE quest_types
SET reward_item_json = '[{"type":"blueprint","blueprint_key":"energy_shield_booster"},{"type":"blueprint","blueprint_key":"anti_aircraft_cannon"},{"type":"blueprint","blueprint_key":"starlight_missile_pod"}]'
WHERE quest_key = 'main_17_weapon_research';

-- Sanity check: every blueprint_key referenced above must exist. If any are
-- missing this migration will raise so the deploy fails loudly instead of
-- silently shipping broken quest rewards.
DO $$
DECLARE
    missing_key TEXT;
BEGIN
    FOR missing_key IN
        SELECT k FROM (VALUES
            ('weikes'),
            ('rapid_fire'),
            ('typhoon'),
            ('estrella'),
            ('ship_reinforcement_facility'),
            ('atomic_framework'),
            ('energy_shield_booster'),
            ('anti_aircraft_cannon'),
            ('starlight_missile_pod'),
            ('super_transmission_engine')
        ) AS required(k)
        WHERE NOT EXISTS (
            SELECT 1 FROM blueprints WHERE blueprint_key = required.k
        )
    LOOP
        RAISE EXCEPTION 'Missing blueprint_key in blueprints table: %', missing_key;
    END LOOP;
END $$;
