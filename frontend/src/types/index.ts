export interface Player {
  id: string
  anonymous_id: string
  username: string | null
  level: number
  experience: number
  created_at: string
}

export interface Planet {
  id: string
  player_id: string
  name: string
  position_x: number
  position_y: number
  is_homeworld: boolean
  is_rbp: boolean
  rbp_level: number
  created_at: string
  updated_at: string
}

export interface BuildingWithType {
  id: string
  planet_id: string
  building_type: number
  grid_col: number
  grid_row: number
  level: number
  is_upgrading: boolean
  upgrade_finish_at: string | null
  created_at: string
  updated_at: string
  type_name: string
  display_name: string
  category: string
  base_type: string
  max_level: number
}

export interface Resource {
  id: string
  planet_id: string
  metal: number
  he3: number
  gold: number
  metal_per_hour: number
  he3_per_hour: number
  gold_per_hour: number
  storage_capacity: number
  warehouse_metal: number
  warehouse_he3: number
  warehouse_gold: number
  warehouse_capacity: number
  last_collected_at: string
  updated_at: string
}

export interface ResourcesResponse extends Resource {
  pending_metal: number
  pending_he3: number
  pending_gold: number
}

export interface CollectResponse {
  collected: {
    metal: number
    he3: number
    gold: number
  }
  resources: Resource
}

export interface BuildingResponse {
  building: {
    id: string
    planet_id: string
    building_type: number
    grid_col: number
    grid_row: number
    level: number
    is_upgrading: boolean
    upgrade_finish_at: string | null
    created_at: string
    updated_at: string
  }
  resources: {
    metal: number
    he3: number
    gold: number
  }
}

export interface GuestAuthResponse {
  token: string
  player: {
    id: string
    anonymous_id: string
    level: number
    created_at: string
  }
}

export interface BuildingTypeData {
  name: string
  display_name: string
  category: string
  base_cost_metal: number
  base_cost_he3: number
  base_cost_gold: number
  base_time_seconds: number
  cost_multiplier: number
  time_multiplier: number
  max_level: number
  max_count_per_planet: number
  civic_center_req_per_level: boolean
  base_production_per_hour: number
  production_multiplier: number
}

// ============ Phase 2 Types ============

// Hull types reference data
export interface HullType {
  id: number
  name: string
  display_name: string
  hull_class: 'frigate' | 'cruiser' | 'battleship'
  tier: number
  armor_type: 'nano' | 'chrome' | 'regen' | 'neutralizing'
  base_shield: number
  base_structure: number
  base_stability: number
  base_defense: number
  installation_slots: number
  base_agility: number
  base_movement: number
  base_storage: number
  base_metal_cost: number
  base_he3_cost: number
  base_gold_cost: number
  base_build_time_seconds: number
  description: string
}

// Module types reference data
export interface ModuleType {
  id: number
  name: string
  display_name: string
  category: 'ballistic' | 'directional' | 'missile' | 'ship_based' | 'planetary' | 'structure' | 'shield' | 'air_defense' | 'electronic' | 'storage' | 'transmission'
  tier: number
  damage_type: 'kinetic' | 'heat' | 'explosive' | 'magnetic' | 'siege' | null
  min_damage: number
  max_damage: number
  weapon_range_min: number
  weapon_range_max: number
  cooldown: number
  he3_per_round: number
  volume: number
  max_per_ship: number
  effects_json: string
  metal_cost: number
  he3_cost: number
  gold_cost: number
  build_time_seconds: number
  description: string
}

// Ship design
export interface ShipDesign {
  id: string
  player_id: string
  name: string
  hull_type_id: number
  hull_name?: string
  hull_class?: string
  modules_json: string
  total_shield: number
  total_structure: number
  total_defense: number
  total_agility: number
  total_movement: number
  total_storage: number
  attack_power: number
  weapon_range_min: number
  weapon_range_max: number
  volume_used: number
  he3_per_round: number
  metal_cost: number
  he3_cost: number
  gold_cost: number
  build_time_seconds: number
  created_at: string
  updated_at: string
}

export interface ShipDesignModule {
  module_type_id: number
  quantity: number
  placement_order: number
}

export interface CreateShipDesignRequest {
  name: string
  hull_type_id: number
  modules: ShipDesignModule[]
}

// Ship Factory
export interface ShipFactoryStatus {
  level: number
  production_slots: number
  speed_bonus_pct: number
}

export interface ProductionSlot {
  slot: number
  in_use: boolean
  ship_design_id?: string
  design_name?: string
  quantity?: number
  build_finish_at?: string
}

export interface BuildShipRequest {
  ship_design_id: string
  quantity: number
  production_slot: number
}

export interface BuildShipResponse {
  slot: number
  ship_design_id: string
  quantity: number
  build_finish_at: string
  resources_spent: { metal: number; he3: number; gold: number }
}

// Blueprints
export interface Blueprint {
  id: number
  name: string
  display_name: string
  blueprint_type: 'hull' | 'module'
  hull_type_id: number | null
  module_type_id: number | null
  hull_class?: string
  module_category?: string
  source: 'instance' | 'quest' | 'auction' | 'trafficker' | 'mall'
  research_level: number
  description: string
}

export interface PlayerBlueprint {
  id: string
  player_id: string
  blueprint_id: number
  is_activated: boolean
  research_level: number
  acquired_at: string
  blueprint_name: string
  blueprint_type: 'hull' | 'module'
  hull_type_id: number | null
  module_type_id: number | null
  source: string
  is_researching?: boolean
  research_finish_at?: string
}

export interface ActiveBlueprintResearch {
  id: string
  player_blueprint_id: string
  target_level: number
  is_researching: boolean
  research_finish_at: string
  blueprint_id: number
  blueprint_name: string
  blueprint_type: 'hull' | 'module'
}

// Fleets
export interface Fleet {
  id: string
  player_id: string
  name: string
  formation: string
  targeting_command: string
  commander_id: string | null
  status: 'stationed' | 'traveling' | 'combat'
  stacks: FleetStack[]
  created_at: string
}

export interface FleetStack {
  fleet_id: string
  ship_design_id: string
  ship_design_name?: string
  hull_class?: string
  grid_row: number
  grid_col: number
  ship_count: number
}

export interface CreateFleetRequest {
  name: string
  formation?: string
  targeting_command?: string
}

export interface AssignStackRequest {
  ship_design_id: string
  grid_row: number
  grid_col: number
  ship_count: number
}

export interface AssignStackResponse {
  stack: FleetStack
  available_ships: number
}

// Instances
export interface Instance {
  id: number
  name: string
  type: 'normal'
  difficulty: number
  required_level: number
  max_fleets: number
  ships_lost_on_defeat: boolean
  he3_lost_on_defeat: boolean
  exp_reward: number
  enemy_fleets_json: string
  rewards_json: string
  description: string
}

export interface InstanceDetail extends Instance {
  enemy_fleets: InstanceEnemyFleet[]
  blueprint_pool: number[]
}

export interface InstanceEnemyFleet {
  hull_class: string
  ship_count: number
  power_estimate: number
}

export interface InstanceProgress {
  instance_id: number
  completed: boolean
  best_result?: string
}

export interface InstanceAttemptResponse {
  report_id: string
  result: 'attacker_win' | 'defender_win'
  total_rounds: number
  exp_gained: number
  treasure_box: {
    resources: { metal: number; he3: number; gold: number }
    blueprint_id: number | null
  }
  losses: {
    ships_destroyed: Record<string, number>
    he3_consumed: number
  }
}

// Spacedock
export interface SpacedockStatus {
  level: number
  repair_pct: number
  active_repairs: number
}

export interface SpacedockRepair {
  id: string
  player_id: string
  ship_design_id: string
  destroyed_count: number
  repaired_count: number
  repair_finish_at: string | null
  combat_report_id: string | null
  created_at: string
}

// ============ Quest Types ============

export interface PlayerQuestWithType {
  id: string
  quest_key: string
  display_name: string
  description: string
  category: 'main' | 'side'
  status: 'locked' | 'available' | 'completed' | 'claimed'
  progress_value: number
  requirement_value: number
  chain_order: number
  reward_metal: number
  reward_he3: number
  reward_gold: number
  reward_item_json: string | null
}

export interface QuestsResponse {
  main_quests: PlayerQuestWithType[]
  side_quests: PlayerQuestWithType[]
  current_main_quest: PlayerQuestWithType | null
}

export interface DailyQuestEntry {
  quest_key: string
  display_name: string
  completed: boolean
  points: number
  progress?: number
  required?: number
  points_per?: number
  max_points?: number
}

export interface DailyTierReward {
  tier: string
  points_required: number
  claimed: boolean
}

export interface DailyQuestsResponse {
  date: string
  daily_points: number
  quests: DailyQuestEntry[]
  tier_rewards: DailyTierReward[]
}

export interface ClaimQuestResponse {
  quest_key: string
  rewards: {
    metal: number
    he3: number
    gold: number
    items: string | null
  }
  next_quest_unlocked: string | null
}

export interface ClaimDailyTierResponse {
  tier: string
  reward: {
    type: string
    quantity: number
  }
}

// ============ Research/Tech Tree Types ============

export type TechTree =
  | 'logistics_construction'
  | 'ballistics_science'
  | 'ship_defense_science'
  | 'directional_science'
  | 'missile_science'
  | 'ship_based_science'
  | 'planetary_defense'

export interface TechType {
  id: number
  name: string
  display_name: string
  tree: TechTree
  max_level: number
  prerequisites: { tech: string; level: number }[]
  base_cost_gold: number
  cost_multiplier: number
  base_time_seconds: number
  time_multiplier: number
  effects: {
    type: string
    per_level?: number
    unit?: string
    [key: string]: any
  }
  description: string
}

export interface TechWithProgress extends TechType {
  current_level: number
  is_researching: boolean
  research_finish_at: string | null
  cost_next_level: { gold: number } | null
  time_next_level_seconds: number | null
}

export interface ResearchTreeResponse {
  tree: string
  techs: TechWithProgress[]
}

export interface ResearchAllResponse {
  trees: ResearchTreeResponse[]
  active: ActiveResearch | null
}

export interface ActiveResearch {
  tech_type_id: number
  tech_name: string
  display_name: string
  tree: string
  level: number
  is_researching: boolean
  research_finish_at: string
}

export interface StartResearchResponse {
  technology: {
    id: string
    tech_type: number
    level: number
    is_researching: boolean
    research_finish_at: string
  }
  resources: {
    metal: number
    he3: number
    gold: number
  }
}

export interface CancelResearchResponse {
  tech_type_id: number
  refunded_gold: number
  resources: {
    metal: number
    he3: number
    gold: number
  }
}

export interface SpeedupResearchResponse {
  technology: {
    id: string
    tech_type: number
    is_researching: boolean
    research_finish_at: string
  }
  vouchers_spent: number
}

// ============ Recycling Plant Types ============

export interface RecyclingJob {
  id: string
  player_id: string
  ship_design_id: string
  metal_gained: number
  he3_gained: number
  gold_gained: number
  duration_seconds: number
  started_at: string
  completed_at: string
  collected: boolean
}

export interface StartRecycleResponse {
  job_id: string
  metal_gained: number
  he3_gained: number
  gold_gained: number
  duration_seconds: number
  completed_at: string
}

export interface CollectRecycleResponse {
  success: boolean
  metal_gained: number
  he3_gained: number
  gold_gained: number
}

export interface AvailableShip {
  id: string
  design_name: string
  hull_class: string
  quantity: number
}

// ============ Chat System Types ============

export interface ChatMessage {
  id: string
  player_id: string
  player_name: string
  message: string
  channel: 'world' | 'alliance'
  created_at: string
}

export interface SendMessageRequest {
  message: string
  channel: 'world' | 'alliance'
}

export interface SendMessageResponse {
  success: boolean
  message_id: string
}

// ============ PvP Combat Types ============

export interface PvPSearchResult {
  planet_id: string
  planet_name: string
  player_id: string
  player_name: string
}

export interface PvPLoot {
  metal: number
  he3: number
  gold: number
}

export interface PvPBattleLosses {
  ships_destroyed: Record<string, number>
  he3_consumed: number
}

export interface AttackPlanetRequest {
  defender_planet_id: string
  fleet_ids: string[]
}

export interface AttackPlanetResponse {
  report_id: string
  result: 'attacker_win' | 'defender_win' | 'draw'
  total_rounds: number
  loot_gained: PvPLoot | null
  attacker_losses: PvPBattleLosses
  defender_losses: PvPBattleLosses
}

// Building types that can be constructed (from seed data)
export const BUILDING_TYPES = [
  { name: 'metal_collector', display_name: 'Metal Collector', category: 'resource' },
  { name: 'he3_extractor', display_name: 'He3 Extractor', category: 'resource' },
  { name: 'residential_area', display_name: 'Residential Area', category: 'resource' },
  { name: 'resource_warehouse', display_name: 'Resource Warehouse', category: 'resource' },
  { name: 'civic_center', display_name: 'Civic Center', category: 'core' },
  { name: 'technology_center', display_name: 'Technology Center', category: 'core' },
  { name: 'alliance_center', display_name: 'Alliance Center', category: 'core' },
  { name: 'trading_center', display_name: 'Trading Center', category: 'core' },
  { name: 'ship_factory', display_name: 'Ship Factory', category: 'military' },
  { name: 'spacedock', display_name: 'Spacedock', category: 'military' },
  { name: 'command_center', display_name: 'Command Center', category: 'military' },
  { name: 'weapon_research_center', display_name: 'Weapon Research Center', category: 'military' },
  { name: 'space_station', display_name: 'Space Station', category: 'space' },
  { name: 'meteor_star', display_name: 'Meteor Star', category: 'defense' },
  { name: 'particle_cannon', display_name: 'Particle Cannon', category: 'defense' },
  { name: 'anti_aircraft_gun', display_name: 'Anti-Aircraft Gun', category: 'defense' },
  { name: 'thors_cannon', display_name: "Thor's Cannon", category: 'defense' },
] as const
