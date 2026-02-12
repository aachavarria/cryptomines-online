-- ============================================================================
-- MIGRATION: Fix Building Stats (Bugs #7, #8, #9)
-- Created: 2026-02-11
-- Description: Fixes incorrect building stats found during GDD validation
--              - Radar: max_level should be 10, not 9
--              - Ship Factory: costs are 3x too expensive
--              - Galaxy Transporter: wrong concept entirely (IGL access, not transport)
-- ============================================================================

-- ----------------------------------------------------------------------------
-- BUG #7: Fix Radar max_level from 9 to 10
-- Source: https://galaxyonlineii.fandom.com/wiki/Radar
-- ----------------------------------------------------------------------------
UPDATE building_types
SET max_level = 10
WHERE name = 'radar';

-- ----------------------------------------------------------------------------
-- BUG #8: Fix Ship Factory base costs (3× too expensive)
-- Current: 600/450/500 Metal/He3/Gold
-- Correct: 206/163/189 (GO2 Level 1 verified costs)
-- Source: https://galaxyonlineii.fandom.com/wiki/Ship_Factory
-- ----------------------------------------------------------------------------
UPDATE building_types
SET base_cost_metal = 206,
    base_cost_he3 = 163,
    base_cost_gold = 189,
    base_time_seconds = 110
WHERE name = 'ship_factory';

-- ----------------------------------------------------------------------------
-- BUG #9: Fix Galaxy Transporter (wrong concept entirely)
-- This is NOT a resource transport building - it's the Inter-Galactic League
-- (PvP) access building. Should be expensive, single-level, 24-hour build.
--
-- Wrong: max_level=12, costs 350/300/450, build time 100s
-- Correct: max_level=1, costs 20000/20000/20000, build time 24 hours
-- Source: https://galaxyonlineii.fandom.com/wiki/Galaxy_Transporter
-- ----------------------------------------------------------------------------
UPDATE building_types
SET base_cost_metal = 20000,
    base_cost_he3 = 20000,
    base_cost_gold = 20000,
    base_time_seconds = 86400,  -- 24 hours
    max_level = 1,
    cost_multiplier = 1.0000,   -- No cost scaling (only 1 level)
    time_multiplier = 1.0000,   -- No time scaling (only 1 level)
    description = 'Allows access to Inter-Galactic League. 10 free matches per day.'
WHERE name = 'galaxy_transporter';

-- ============================================================================
-- VERIFICATION QUERIES (for testing)
-- ============================================================================
-- SELECT name, max_level FROM building_types WHERE name = 'radar';
-- Expected: max_level = 10
--
-- SELECT name, base_cost_metal, base_cost_he3, base_cost_gold, base_time_seconds
-- FROM building_types WHERE name = 'ship_factory';
-- Expected: 206, 163, 189, 110
--
-- SELECT name, base_cost_metal, base_cost_he3, base_cost_gold, base_time_seconds, max_level
-- FROM building_types WHERE name = 'galaxy_transporter';
-- Expected: 20000, 20000, 20000, 86400, 1
-- ============================================================================
