-- Fix quest rewards - replace out-of-scope items
-- Bug: 9 quests configured with out-of-scope items (loudspeaker, revival_card, healing_card, galaxy_transfer)
-- Backend silently skips these rewards

-- main_01: loudspeaker → construction_card
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"construction_card","quantity":1}]'
WHERE quest_key = 'main_01_collecting_resources';

-- main_03: loudspeaker → primary_metal_pack
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"primary_metal_pack","quantity":1}]'
WHERE quest_key = 'main_03_tech_center';

-- main_09: loudspeaker → extra_tax
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"extra_tax","quantity":1}]'
WHERE quest_key = 'main_09_residential';

-- main_12: revival_card → truce_card
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"truce_card","quantity":1}]'
WHERE quest_key = 'main_12_recruit';

-- main_13: loudspeaker → metal_mining_boost
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"metal_mining_boost","quantity":1}]'
WHERE quest_key = 'main_13_design_ship';

-- main_15: loudspeaker → he3_mining_boost
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"he3_mining_boost","quantity":1}]'
WHERE quest_key = 'main_15_fleet';

-- main_19: galaxy_transfer → primary_metal_pack
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"primary_metal_pack","quantity":1}]'
WHERE quest_key = 'main_19_resource_pack';

-- main_20: healing_card → extra_tax
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"extra_tax","quantity":1}]'
WHERE quest_key = 'main_20_growing_resources';

-- main_22: loudspeaker → truce_card
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"truce_card","quantity":1}]'
WHERE quest_key = 'main_22_mail_system';
