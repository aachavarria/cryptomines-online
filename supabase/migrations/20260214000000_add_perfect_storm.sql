-- ============================================================================
-- MIGRATION: Add Perfect Storm Capstone Technology (Missile Science)
-- Created: 2026-02-14
-- Description: Adds the Perfect Storm capstone tech to complete the Missile
--              Science tree. Wiki details were truncated; values estimated
--              from other capstone techs (Victory Rush, Ingenuity).
-- ============================================================================
-- Reference: docs/gdd/game-design-document.md line 554
-- GDD note: "Perfect Storm [Lv ?] (wiki details truncated) — Capstone tech"
-- Partial wiki data for Rapid Loading (sibling): "Enable Guided weapons to
--   finish reloading 1 round faster. Also increase 10% chance of triggering..."
-- ============================================================================

INSERT INTO tech_types (
    name, display_name, tree, max_level,
    prerequisites_json,
    base_cost_metal, base_cost_he3, base_cost_gold,
    cost_multiplier, base_time_seconds, time_multiplier,
    effects_json, description
) VALUES (
    'perfect_storm',
    'Perfect Storm',
    'missile_science',
    1,
    '[{"tech":"suppression","level":4},{"tech":"missile_concussion","level":1},{"tech":"rapid_loading","level":1}]',
    0, 0, 2000000,
    1.0000, 250000, 1.0000,
    '{"type":"multi_bonus","missile_damage":10,"scatter_all":10,"missile_hit_rate":5,"intercept_reduction":5,"unit":"percent"}',
    '+10% missile damage, +10% scatter to all, +5% hit rate, -5% interception (capstone)'
);
