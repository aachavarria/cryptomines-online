-- Migration: Grant missing blueprint rewards from completed quests
-- Bug: Quest reward system didn't support module blueprints, only hull blueprints
-- Impact: Players who completed main_05 (and possibly other quests) didn't receive blueprint rewards
-- Date: 2026-02-11

-- Step 1: Create the blueprint item type for super_transmission_engine if it doesn't exist
INSERT INTO item_types (item_key, display_name, category, description, blueprint_id)
SELECT
    'blueprint_super_transmission_engine',
    'Blueprint: super_transmission_engine',
    'blueprint',
    'Unlock super_transmission_engine blueprint',
    b.id
FROM blueprints b
WHERE b.blueprint_type = 'module'
  AND b.module_type_id = (SELECT id FROM module_types WHERE name = 'super_transmission_engine' AND tier = 0)
ON CONFLICT (item_key) DO NOTHING;

-- Step 2: Grant the blueprint to players who completed main_05 but didn't receive it
-- Only grant if the quest is completed AND the player doesn't already have the item
INSERT INTO player_inventory (player_id, item_key, quantity)
SELECT DISTINCT
    pq.player_id,
    'blueprint_super_transmission_engine',
    1
FROM player_quests pq
JOIN quest_types qt ON pq.quest_type_id = qt.id
WHERE qt.quest_key = 'main_05_metal_production'
  AND pq.is_completed = true
  AND NOT EXISTS (
      SELECT 1 FROM player_inventory pi
      WHERE pi.player_id = pq.player_id
        AND pi.item_key = 'blueprint_super_transmission_engine'
  )
ON CONFLICT (player_id, item_key) DO NOTHING;

-- Log the fix
COMMENT ON TABLE item_types IS 'Fixed missing module blueprint rewards for quest main_05 (2026-02-11)';
