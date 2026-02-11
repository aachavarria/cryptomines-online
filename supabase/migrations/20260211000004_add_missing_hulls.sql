-- ============================================================================
-- MIGRATION: Add Missing Ship Hulls (27 hulls total)
-- Created: 2026-02-11
-- Description: Adds 6 missing cruiser hulls and 21 missing battleship hulls
--              from Galaxy Online 2 validation research
-- ============================================================================
-- Source: SHIP_HULLS_VALIDATION_REPORT.md
-- Reference: docs/validation/SHIP_HULLS_VALIDATION_REPORT.md
-- ============================================================================

-- ============================================================================
-- PART 1: MISSING CRUISER HULLS (6 hulls = 2 lines × 3 tiers)
-- ============================================================================
-- These are advanced cruiser lines with higher stats than standard cruisers
-- Both use Regen armor type
-- ============================================================================

INSERT INTO hull_types (name, display_name, hull_class, tier, armor_type, base_shield, base_structure, base_stability, base_defense, installation_slots, base_agility, base_movement, base_storage, base_metal_cost, base_he3_cost, base_gold_cost, base_build_time_seconds, description) VALUES

-- ========================================
-- Chimera Capra (Advanced Regen Cruiser)
-- ========================================
-- GO2 Wiki: Shields 1,310-1,790 | Structure 7,300-10,000
-- Using 2x multiplier for our scaling: Shields 2,620-3,580 | Structure 14,600-20,000
('chimera_capra_i',   'Chimera Capra I',   'cruiser', 1, 'regen', 2620, 14600, 100.00, 0.00, 200, 0, 0, 0, 8000,  5500,  4000,  60, 'Advanced regen cruiser with high structure'),
('chimera_capra_ii',  'Chimera Capra II',  'cruiser', 2, 'regen', 3100, 17300, 100.00, 0.00, 280, 0, 0, 0, 24000, 16500, 12000, 120, 'Improved Chimera Capra hull'),
('chimera_capra_iii', 'Chimera Capra III', 'cruiser', 3, 'regen', 3580, 20000, 100.00, 0.00, 360, 0, 0, 0, 72000, 49500, 36000, 240, 'Advanced Chimera Capra hull'),

-- ========================================
-- Ultra Gwyar (Advanced Regen Cruiser)
-- ========================================
-- GO2 Wiki: Shields 2,300-2,800 | Structure 10,950-15,000
-- Using 2x multiplier for our scaling: Shields 4,600-5,600 | Structure 21,900-30,000
('ultra_gwyar_i',     'Ultra Gwyar I',     'cruiser', 1, 'regen', 4600, 21900, 100.00, 0.00, 220, 0, 0, 0, 10000, 7000,  5000,  70, 'Elite regen cruiser with highest stats'),
('ultra_gwyar_ii',    'Ultra Gwyar II',    'cruiser', 2, 'regen', 5200, 25950, 100.00, 0.00, 300, 0, 0, 0, 30000, 21000, 15000, 140, 'Improved Ultra Gwyar hull'),
('ultra_gwyar_iii',   'Ultra Gwyar III',   'cruiser', 3, 'regen', 5600, 30000, 100.00, 0.00, 380, 0, 0, 0, 90000, 63000, 45000, 280, 'Advanced Ultra Gwyar hull');

-- ============================================================================
-- PART 2: MISSING BATTLESHIP HULLS (21 hulls = 7 lines × 3 tiers)
-- ============================================================================
-- Adding 7 missing battleship lines to complete the battleship roster
-- Armor types: 3 Chrome, 2 Regen, 1 Neutralizing, 1 Nano
-- ============================================================================

INSERT INTO hull_types (name, display_name, hull_class, tier, armor_type, base_shield, base_structure, base_stability, base_defense, installation_slots, base_agility, base_movement, base_storage, base_metal_cost, base_he3_cost, base_gold_cost, base_build_time_seconds, description) VALUES

-- ========================================
-- Howler (Chrome Battleship)
-- ========================================
-- Heavy chrome armor battleship with balanced stats
('howler_i',          'Howler I',          'battleship', 1, 'chrome', 400,  4400, 100.00, 0.00, 170, 0, 0, 0, 3200,  2200,  1600,  28, 'Chrome armor battleship'),
('howler_ii',         'Howler II',         'battleship', 2, 'chrome', 800,  8800, 100.00, 0.00, 250, 0, 0, 0, 9600,  6600,  4800,  56, 'Improved Howler hull'),
('howler_iii',        'Howler III',        'battleship', 3, 'chrome', 1600, 17600,100.00, 0.00, 370, 0, 0, 0, 28800, 19800, 14400, 112, 'Advanced Howler hull'),

-- ========================================
-- Whirlpool (Regen Battleship)
-- ========================================
-- Regenerative armor battleship with high capacity
('whirlpool_i',       'Whirlpool I',       'battleship', 1, 'regen', 420,  4600, 100.00, 0.00, 175, 0, 0, 0, 3400,  2300,  1700,  30, 'Regen armor battleship'),
('whirlpool_ii',      'Whirlpool II',      'battleship', 2, 'regen', 840,  9200, 100.00, 0.00, 255, 0, 0, 0, 10200, 6900,  5100,  60, 'Improved Whirlpool hull'),
('whirlpool_iii',     'Whirlpool III',     'battleship', 3, 'regen', 1680, 18400,100.00, 0.00, 375, 0, 0, 0, 30600, 20700, 15300, 120, 'Advanced Whirlpool hull'),

-- ========================================
-- Cerberus (Neutralizing Battleship)
-- ========================================
-- Neutralizing armor battleship with specialized defense
('cerberus_i',        'Cerberus I',        'battleship', 1, 'neutralizing', 440,  4800, 100.00, 0.00, 180, 0, 0, 0, 3600,  2400,  1800,  32, 'Neutralizing armor battleship'),
('cerberus_ii',       'Cerberus II',       'battleship', 2, 'neutralizing', 880,  9600, 100.00, 0.00, 260, 0, 0, 0, 10800, 7200,  5400,  64, 'Improved Cerberus hull'),
('cerberus_iii',      'Cerberus III',      'battleship', 3, 'neutralizing', 1760, 19200,100.00, 0.00, 380, 0, 0, 0, 32400, 21600, 16200, 128, 'Advanced Cerberus hull'),

-- ========================================
-- Genesis (Regen Battleship)
-- ========================================
-- Advanced regenerative armor battleship
('genesis_i',         'Genesis I',         'battleship', 1, 'regen', 460,  5000, 100.00, 0.00, 185, 0, 0, 0, 3800,  2600,  1900,  34, 'Advanced regen battleship'),
('genesis_ii',        'Genesis II',        'battleship', 2, 'regen', 920,  10000,100.00, 0.00, 265, 0, 0, 0, 11400, 7800,  5700,  68, 'Improved Genesis hull'),
('genesis_iii',       'Genesis III',       'battleship', 3, 'regen', 1840, 20000,100.00, 0.00, 385, 0, 0, 0, 34200, 23400, 17100, 136, 'Advanced Genesis hull'),

-- ========================================
-- Tiamat (Chrome Battleship)
-- ========================================
-- Elite chrome armor battleship with massive firepower
('tiamat_i',          'Tiamat I',          'battleship', 1, 'chrome', 480,  5200, 100.00, 0.00, 190, 0, 0, 0, 4000,  2800,  2000,  36, 'Elite chrome battleship'),
('tiamat_ii',         'Tiamat II',         'battleship', 2, 'chrome', 960,  10400,100.00, 0.00, 270, 0, 0, 0, 12000, 8400,  6000,  72, 'Improved Tiamat hull'),
('tiamat_iii',        'Tiamat III',        'battleship', 3, 'chrome', 1920, 20800,100.00, 0.00, 390, 0, 0, 0, 36000, 25200, 18000, 144, 'Advanced Tiamat hull'),

-- ========================================
-- Chimera Viper (Nano Battleship)
-- ========================================
-- Advanced nano armor battleship
('chimera_viper_i',   'Chimera Viper I',   'battleship', 1, 'nano', 500,  5400, 100.00, 0.00, 195, 0, 0, 0, 4500,  3000,  2200,  38, 'Advanced nano battleship'),
('chimera_viper_ii',  'Chimera Viper II',  'battleship', 2, 'nano', 1000, 10800,100.00, 0.00, 275, 0, 0, 0, 13500, 9000,  6600,  76, 'Improved Chimera Viper hull'),
('chimera_viper_iii', 'Chimera Viper III', 'battleship', 3, 'nano', 2000, 21600,100.00, 0.00, 395, 0, 0, 0, 40500, 27000, 19800, 152, 'Advanced Chimera Viper hull'),

-- ========================================
-- Ultra Calas (Chrome Battleship)
-- ========================================
-- Ultimate chrome armor battleship with highest stats
('ultra_calas_i',     'Ultra Calas I',     'battleship', 1, 'chrome', 520,  5600, 100.00, 0.00, 200, 0, 0, 0, 5000,  3500,  2500,  40, 'Ultimate chrome battleship'),
('ultra_calas_ii',    'Ultra Calas II',    'battleship', 2, 'chrome', 1040, 11200,100.00, 0.00, 280, 0, 0, 0, 15000, 10500, 7500,  80, 'Improved Ultra Calas hull'),
('ultra_calas_iii',   'Ultra Calas III',   'battleship', 3, 'chrome', 2080, 22400,100.00, 0.00, 400, 0, 0, 0, 45000, 31500, 22500, 160, 'Advanced Ultra Calas hull');

-- ============================================================================
-- PART 3: BLUEPRINTS FOR NEW HULLS
-- ============================================================================
-- Adding blueprints for all new hulls to the blueprints table
-- ============================================================================

INSERT INTO blueprints (name, blueprint_type, hull_type_id, source, description) VALUES
-- Cruiser Blueprints
('Chimera Capra Blueprint', 'hull', (SELECT id FROM hull_types WHERE name = 'chimera_capra_i'), 'instance', 'Unlocks Chimera Capra cruiser line'),
('Ultra Gwyar Blueprint',   'hull', (SELECT id FROM hull_types WHERE name = 'ultra_gwyar_i'),   'instance', 'Unlocks Ultra Gwyar cruiser line'),

-- Battleship Blueprints
('Howler Blueprint',        'hull', (SELECT id FROM hull_types WHERE name = 'howler_i'),        'instance', 'Unlocks Howler battleship line'),
('Whirlpool Blueprint',     'hull', (SELECT id FROM hull_types WHERE name = 'whirlpool_i'),     'instance', 'Unlocks Whirlpool battleship line'),
('Cerberus Blueprint',      'hull', (SELECT id FROM hull_types WHERE name = 'cerberus_i'),      'instance', 'Unlocks Cerberus battleship line'),
('Genesis Blueprint',       'hull', (SELECT id FROM hull_types WHERE name = 'genesis_i'),       'instance', 'Unlocks Genesis battleship line'),
('Tiamat Blueprint',        'hull', (SELECT id FROM hull_types WHERE name = 'tiamat_i'),        'instance', 'Unlocks Tiamat battleship line'),
('Chimera Viper Blueprint', 'hull', (SELECT id FROM hull_types WHERE name = 'chimera_viper_i'), 'instance', 'Unlocks Chimera Viper battleship line'),
('Ultra Calas Blueprint',   'hull', (SELECT id FROM hull_types WHERE name = 'ultra_calas_i'),   'instance', 'Unlocks Ultra Calas battleship line');

-- ============================================================================
-- VERIFICATION QUERIES (for testing)
-- ============================================================================
-- SELECT hull_class, armor_type, COUNT(*) as count
-- FROM hull_types
-- GROUP BY hull_class, armor_type
-- ORDER BY hull_class, armor_type;
--
-- Expected results after migration:
-- Frigates: 30 total (15 nano, 15 neutralizing)
-- Cruisers: 36 total (15 chrome, 21 regen) - was 30, now +6
-- Battleships: 36 total (varied) - was 15, now +21
-- TOTAL: 102 hulls
-- ============================================================================

-- ============================================================================
-- SUMMARY
-- ============================================================================
-- Added 2 advanced cruiser lines:
--   - Chimera Capra (Regen) - 3 tiers
--   - Ultra Gwyar (Regen) - 3 tiers
--
-- Added 7 battleship lines:
--   - Howler (Chrome) - 3 tiers
--   - Whirlpool (Regen) - 3 tiers
--   - Cerberus (Neutralizing) - 3 tiers
--   - Genesis (Regen) - 3 tiers
--   - Tiamat (Chrome) - 3 tiers
--   - Chimera Viper (Nano) - 3 tiers
--   - Ultra Calas (Chrome) - 3 tiers
--
-- Total: 27 new hulls (6 cruisers + 21 battleships)
-- All hulls include proper armor types, balanced stats, and blueprints
-- ============================================================================
