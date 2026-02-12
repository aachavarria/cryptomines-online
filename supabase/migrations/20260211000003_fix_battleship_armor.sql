-- ============================================================================
-- MIGRATION: Fix Battleship Armor Types (Bug #5)
-- Created: 2026-02-11
-- Description: Fixes incorrect armor_type values for 4 battleship hull lines.
--              Wrong armor types cause incorrect damage calculations in combat.
-- ============================================================================
-- Source: SHIP_HULLS_VALIDATION_REPORT.md
-- Reference: GDD_VALIDATION_MASTER_REPORT.md Bug #5
-- ============================================================================

-- ----------------------------------------------------------------------------
-- ARMOR TYPE CORRECTIONS
-- ----------------------------------------------------------------------------
-- Nettle I/II/III: regen → nano
-- Diaz I/II/III: nano → neutralizing
-- RV766-The Explorer I/II/III: neutralizing → regen
-- Palenka I/II/III: chrome → nano
-- ----------------------------------------------------------------------------

-- Fix Nettle I/II/III (regen → nano)
-- Current: 'regen' | Correct: 'nano'
-- High-structure battleship with nano armor protection
UPDATE hull_types
SET armor_type = 'nano'
WHERE name IN ('nettle_i', 'nettle_ii', 'nettle_iii');

-- Fix Diaz I/II/III (nano → neutralizing)
-- Current: 'nano' | Correct: 'neutralizing'
-- Neutralizing armor provides better defense against certain weapons
UPDATE hull_types
SET armor_type = 'neutralizing'
WHERE name IN ('diaz_i', 'diaz_ii', 'diaz_iii');

-- Fix RV766-The Explorer I/II/III (neutralizing → regen)
-- Current: 'neutralizing' | Correct: 'regen'
-- Explorer variant with regenerative armor capabilities
UPDATE hull_types
SET armor_type = 'regen'
WHERE name IN ('rv766_i', 'rv766_ii', 'rv766_iii');

-- Fix Palenka I/II/III (chrome → nano)
-- Current: 'chrome' | Correct: 'nano'
-- Heavy battleship with nano armor protection
UPDATE hull_types
SET armor_type = 'nano'
WHERE name IN ('palenka_i', 'palenka_ii', 'palenka_iii');

-- ============================================================================
-- VERIFICATION QUERIES (for testing)
-- ============================================================================
-- SELECT name, hull_class, tier, armor_type
-- FROM hull_types
-- WHERE name LIKE 'nettle_%' OR name LIKE 'diaz_%'
--    OR name LIKE 'rv766_%' OR name LIKE 'palenka_%'
-- ORDER BY name;
--
-- Expected results:
-- diaz_i, diaz_ii, diaz_iii: armor_type = 'neutralizing'
-- nettle_i, nettle_ii, nettle_iii: armor_type = 'nano'
-- palenka_i, palenka_ii, palenka_iii: armor_type = 'nano'
-- rv766_i, rv766_ii, rv766_iii: armor_type = 'regen'
-- ============================================================================

-- ============================================================================
-- IMPACT ANALYSIS
-- ============================================================================
-- Armor effectiveness matrix determines damage reduction in combat:
-- - Nano armor: Better vs certain weapon types
-- - Chrome armor: Different damage reduction profile
-- - Regen armor: Regenerative capabilities
-- - Neutralizing armor: Specialized damage mitigation
--
-- These corrections ensure proper combat balance and accurate damage calculations.
-- ============================================================================
