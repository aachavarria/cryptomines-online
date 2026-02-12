-- Add quest descriptions to all 46 quests
-- Bug: All quests have NULL descriptions, leaving players with zero guidance

-- =============================================================================
-- Main Quest Descriptions (28 quests)
-- =============================================================================

-- Phase 1: Resource & Building Focus (13 quests)
UPDATE quest_types SET description = 'Your planet is rich in natural resources. Harvest your Resource Warehouse to collect stored resources and begin building your empire.'
WHERE quest_key = 'main_01_collecting_resources';

UPDATE quest_types SET description = 'Research is the key to galactic dominance. Build a Technology Center to unlock the research tree and gain access to powerful technologies.'
WHERE quest_key = 'main_03_tech_center';

UPDATE quest_types SET description = 'Start your first research project. Research Concurrent Construction to unlock an additional construction slot, allowing you to build multiple structures at once.'
WHERE quest_key = 'main_04_research';

UPDATE quest_types SET description = 'Metal is the foundation of any space empire. Build a Metal Collector to begin harvesting this essential resource for construction and ship production.'
WHERE quest_key = 'main_05_metal_production';

UPDATE quest_types SET description = 'Blueprints unlock advanced ship modules. Use the Super Transmission Engine blueprint to upgrade your ship designs with faster engines.'
WHERE quest_key = 'main_06_blueprints_1';

UPDATE quest_types SET description = 'Helium-3 is the fuel of the future. Build a He3 Extractor to begin producing this critical resource for advanced technologies and ship operations.'
WHERE quest_key = 'main_07_he3_production';

UPDATE quest_types SET description = 'Expand your blueprint collection. Use the Estrella blueprint to unlock powerful new ship designs for your fleet.'
WHERE quest_key = 'main_08_blueprints_2';

UPDATE quest_types SET description = 'A thriving population generates valuable gold income. Build a Residential Area to house your citizens and boost your economy.'
WHERE quest_key = 'main_09_residential';

UPDATE quest_types SET description = 'Expand your station capabilities. Upgrade your Space Station to level 2 to unlock higher-tier buildings and research.'
WHERE quest_key = 'main_23_space_station';

UPDATE quest_types SET description = 'Defend your territory from attackers. Build defensive structures to protect your planet from enemy raids.'
WHERE quest_key = 'main_24_space_defense';

UPDATE quest_types SET description = 'Increase your metal production. Upgrade your Metal Collector to level 2 to double your metal income and support larger construction projects.'
WHERE quest_key = 'main_25_metal_lv2';

UPDATE quest_types SET description = 'Boost your He3 output. Upgrade your He3 Extractor to level 2 to increase fuel production for your growing fleet.'
WHERE quest_key = 'main_26_he3_lv2';

UPDATE quest_types SET description = 'Grow your population and income. Upgrade your Residential Area to level 2 to generate more gold for your empire.'
WHERE quest_key = 'main_27_residential_lv2';

-- Phase 2: Military Focus (8 quests)
UPDATE quest_types SET description = 'Begin your military expansion. Build a Ship Factory to start producing warships for conquest and defense.'
WHERE quest_key = 'main_10_ship_factory';

UPDATE quest_types SET description = 'Leadership is essential for victory. Build a Command Center to recruit commanders who will lead your fleets into battle.'
WHERE quest_key = 'main_11_command_center';

UPDATE quest_types SET description = 'Every fleet needs a commander. Recruit your first commander to lead your ships and gain combat bonuses.'
WHERE quest_key = 'main_12_recruit';

UPDATE quest_types SET description = 'Design your first custom warship. Create a ship design using your blueprints and modules to prepare for production.'
WHERE quest_key = 'main_13_design_ship';

UPDATE quest_types SET description = 'Produce your first warship. Build ships from your designs to assemble a combat-ready fleet.'
WHERE quest_key = 'main_14_build_ships';

UPDATE quest_types SET description = 'Organize your military forces. Create a fleet and assign a commander to prepare for space combat operations.'
WHERE quest_key = 'main_15_fleet';

UPDATE quest_types SET description = 'Ammunition is the lifeblood of war. Replenish your fleet''s ammo to ensure they''re ready for extended combat missions.'
WHERE quest_key = 'main_16_logistics';

UPDATE quest_types SET description = 'Unlock advanced weaponry. Build a Weapon Research Center to research cutting-edge weapon technologies for your ships.'
WHERE quest_key = 'main_17_weapon_research';

UPDATE quest_types SET description = 'Scale up ship production. Upgrade your Ship Factory to level 2 to build more powerful vessels and produce ships faster.'
WHERE quest_key = 'main_28_ship_factory_lv2';

-- Phase 3: Social/Deferred Quests (6 quests)
UPDATE quest_types SET description = 'Make your voice heard across the galaxy. Send a message in the world channel to connect with other players.'
WHERE quest_key = 'main_02_loud_and_clear';

UPDATE quest_types SET description = 'Expand your carrying capacity. Increase your bag slots to store more items and resources for future use.'
WHERE quest_key = 'main_18_bigger_bags';

UPDATE quest_types SET description = 'Put your resources to work. Use a resource pack to instantly gain valuable materials for construction and research.'
WHERE quest_key = 'main_19_resource_pack';

UPDATE quest_types SET description = 'Invest in long-term growth. Grow your comsats to establish permanent resource income streams across the galaxy.'
WHERE quest_key = 'main_20_growing_resources';

UPDATE quest_types SET description = 'Build your network. Add a friend to your contacts list to coordinate strategies and support each other in the galaxy.'
WHERE quest_key = 'main_21_adding_friends';

UPDATE quest_types SET description = 'Master interstellar communication. Send a private message using the mail system to coordinate with allies.'
WHERE quest_key = 'main_22_mail_system';

-- =============================================================================
-- Side Quest Descriptions (12 quests)
-- =============================================================================

-- Harvest Time Series (Metal Production)
UPDATE quest_types SET description = 'Boost your metal production to 2,180 per hour. Upgrade your Metal Collectors to reach this production milestone.'
WHERE quest_key = 'side_harvest_time_1';

UPDATE quest_types SET description = 'Double your metal output to 4,360 per hour. Continue upgrading Metal Collectors to maintain industrial growth.'
WHERE quest_key = 'side_harvest_time_2';

UPDATE quest_types SET description = 'Achieve elite metal production at 8,720 per hour. Maximize your Metal Collector levels for massive industrial capacity.'
WHERE quest_key = 'side_harvest_time_3';

-- Gathering He3 Series (He3 Production)
UPDATE quest_types SET description = 'Increase He3 production to 2,360 per hour. Upgrade your He3 Extractors to fuel your growing military operations.'
WHERE quest_key = 'side_gathering_he3_1';

UPDATE quest_types SET description = 'Expand He3 output to 4,720 per hour. Advanced ship systems require substantial fuel reserves.'
WHERE quest_key = 'side_gathering_he3_2';

UPDATE quest_types SET description = 'Reach peak He3 production at 9,440 per hour. Maximum fuel output supports the largest fleet operations.'
WHERE quest_key = 'side_gathering_he3_3';

-- Raising Morale Series (Gold Production)
UPDATE quest_types SET description = 'Grow your economy to produce 2,800 gold per hour. Upgrade Residential Areas to increase population and tax income.'
WHERE quest_key = 'side_raising_morale_1';

UPDATE quest_types SET description = 'Boost gold production to 5,600 per hour. A prosperous population funds your military and research efforts.'
WHERE quest_key = 'side_raising_morale_2';

UPDATE quest_types SET description = 'Achieve maximum gold output at 11,200 per hour. Peak economic growth enables rapid expansion and dominance.'
WHERE quest_key = 'side_raising_morale_3';

-- Plentiful Resources Series (Storage Capacity)
UPDATE quest_types SET description = 'Reach 50,000 total storage capacity. Upgrade warehouses to stockpile resources for major construction projects.'
WHERE quest_key = 'side_plentiful_resources_1';

UPDATE quest_types SET description = 'Expand storage to 200,000 capacity. Large reserves protect against raids and enable long-term planning.'
WHERE quest_key = 'side_plentiful_resources_2';

UPDATE quest_types SET description = 'Achieve 1,000,000 storage capacity. Massive stockpiles support empire-scale operations and sustained warfare.'
WHERE quest_key = 'side_plentiful_resources_3';

-- =============================================================================
-- Daily Quest Descriptions (6 quests)
-- =============================================================================

UPDATE quest_types SET description = 'Log in today to earn daily quest points. Consistency is rewarded in the galaxy.'
WHERE quest_key = 'daily_login';

UPDATE quest_types SET description = 'Harvest resources from any production building. Collect your daily production to keep resources flowing.'
WHERE quest_key = 'daily_collect_dues';

UPDATE quest_types SET description = 'Use a construction speedup item to instantly complete building projects. Time is a valuable resource.'
WHERE quest_key = 'daily_need_for_speed';

UPDATE quest_types SET description = 'Harvest your Resource Warehouse 3 times today. Regular collection maximizes resource accumulation.'
WHERE quest_key = 'daily_stockpiling';

UPDATE quest_types SET description = 'Donate 200,000 resources to your alliance. Support your allies and earn valuable alliance rewards.'
WHERE quest_key = 'daily_donations';

UPDATE quest_types SET description = 'Complete 2 restricted instance battles. Test your fleet strength in challenging PvE combat scenarios.'
WHERE quest_key = 'daily_restricted_instances';
