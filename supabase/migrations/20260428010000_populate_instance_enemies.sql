-- =============================================================================
-- Populate enemy_fleets_json for the 30 seeded normal instances.
--
-- Why: every instance shipped with enemy_fleets_json='[]', so attempting any
-- instance auto-wins (the combat engine resolves with zero defender stacks).
-- The QA E2E pass surfaced this; the schema has been valid all along, just
-- empty. We scale ship counts with difficulty and rotate hull lines so
-- different instances feel distinct.
--
-- Tiers used:
--   difficulty  1-10  → tier-I frigates with light cruiser support
--   difficulty 11-20  → tier-II cruiser-heavy fleets
--   difficulty 21-30  → tier-III battleships flanked by cruisers/frigates
-- grid_row 0 = front, 1 = mid, 2 = back (matches the 3x3 combat grid)
-- =============================================================================

UPDATE instances SET enemy_fleets_json = jsonb_build_array(
    jsonb_build_object('hull_type', 'weikes_i',     'quantity', 40 + difficulty * 15, 'grid_row', 0, 'grid_col', 0),
    jsonb_build_object('hull_type', 'air_wanderer_i','quantity', 20 + difficulty * 10, 'grid_row', 1, 'grid_col', 0)
) WHERE type = 'normal' AND difficulty BETWEEN 1 AND 5;

UPDATE instances SET enemy_fleets_json = jsonb_build_array(
    jsonb_build_object('hull_type', 'typhoon_i',    'quantity', 20 + difficulty * 8,  'grid_row', 0, 'grid_col', 0),
    jsonb_build_object('hull_type', 'weikes_i',     'quantity', 60 + difficulty * 12, 'grid_row', 1, 'grid_col', 0),
    jsonb_build_object('hull_type', 'sparrow_i',    'quantity', 30 + difficulty * 8,  'grid_row', 2, 'grid_col', 0)
) WHERE type = 'normal' AND difficulty BETWEEN 6 AND 10;

UPDATE instances SET enemy_fleets_json = jsonb_build_array(
    jsonb_build_object('hull_type', 'duke_i',       'quantity', 30 + difficulty * 6,  'grid_row', 0, 'grid_col', 0),
    jsonb_build_object('hull_type', 'weikes_ii',    'quantity', 40 + difficulty * 10, 'grid_row', 1, 'grid_col', 0),
    jsonb_build_object('hull_type', 'wraith_i',     'quantity', 20 + difficulty * 5,  'grid_row', 2, 'grid_col', 0)
) WHERE type = 'normal' AND difficulty BETWEEN 11 AND 15;

UPDATE instances SET enemy_fleets_json = jsonb_build_array(
    jsonb_build_object('hull_type', 'typhoon_ii',   'quantity', 25 + difficulty * 5,  'grid_row', 0, 'grid_col', 0),
    jsonb_build_object('hull_type', 'spinner_i',    'quantity', 30 + difficulty * 6,  'grid_row', 1, 'grid_col', 0),
    jsonb_build_object('hull_type', 'sparrow_ii',   'quantity', 20 + difficulty * 5,  'grid_row', 2, 'grid_col', 0)
) WHERE type = 'normal' AND difficulty BETWEEN 16 AND 20;

UPDATE instances SET enemy_fleets_json = jsonb_build_array(
    jsonb_build_object('hull_type', 'estrella_i',   'quantity', 5  + difficulty * 2,  'grid_row', 0, 'grid_col', 0),
    jsonb_build_object('hull_type', 'duke_ii',      'quantity', 20 + difficulty * 4,  'grid_row', 1, 'grid_col', 0),
    jsonb_build_object('hull_type', 'weikes_iii',   'quantity', 40 + difficulty * 6,  'grid_row', 2, 'grid_col', 0)
) WHERE type = 'normal' AND difficulty BETWEEN 21 AND 25;

UPDATE instances SET enemy_fleets_json = jsonb_build_array(
    jsonb_build_object('hull_type', 'estrella_ii',  'quantity', 8  + difficulty * 2,  'grid_row', 0, 'grid_col', 0),
    jsonb_build_object('hull_type', 'nettle_i',     'quantity', 5  + difficulty * 2,  'grid_row', 0, 'grid_col', 1),
    jsonb_build_object('hull_type', 'wraith_ii',    'quantity', 25 + difficulty * 4,  'grid_row', 1, 'grid_col', 0),
    jsonb_build_object('hull_type', 'sparrow_iii',  'quantity', 30 + difficulty * 5,  'grid_row', 2, 'grid_col', 0)
) WHERE type = 'normal' AND difficulty BETWEEN 26 AND 30;
