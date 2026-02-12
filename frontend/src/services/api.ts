import axios from 'axios'
import type {
  GuestAuthResponse,
  Planet,
  BuildingWithType,
  ResourcesResponse,
  CollectResponse,
  BuildingResponse,
  BuildingTypeData,
  HullType,
  ModuleType,
  ShipDesign,
  CreateShipDesignRequest,
  ShipFactoryStatus,
  ProductionSlot,
  BuildShipRequest,
  BuildShipResponse,
  Blueprint,
  PlayerBlueprint,
  ActiveBlueprintResearch,
  Fleet,
  CreateFleetRequest,
  AssignStackRequest,
  AssignStackResponse,
  Instance,
  InstanceDetail,
  InstanceProgress,
  InstanceAttemptResponse,
  SpacedockStatus,
  SpacedockRepair,
  QuestsResponse,
  DailyQuestsResponse,
  ClaimQuestResponse,
  ClaimDailyTierResponse,
  ResearchAllResponse,
  ResearchTreeResponse,
  StartResearchResponse,
  CancelResearchResponse,
  SpeedupResearchResponse,
  ActiveResearch,
  RecyclingJob,
  StartRecycleResponse,
  CollectRecycleResponse,
  AvailableShip,
  ChatMessage,
  SendMessageResponse,
  PvPSearchResult,
  AttackPlanetRequest,
  AttackPlanetResponse,
} from '../types'
import type { InventoryItem, UseItemResponse } from '../types/inventory'

const api = axios.create({
  baseURL: '/api',
})

// Request interceptor to add auth token and dev-mode header
api.interceptors.request.use(async (config) => {
  const { data: { session } } = await import('../lib/supabase').then(m => m.supabase.auth.getSession())
  if (session?.access_token) {
    config.headers.Authorization = `Bearer ${session.access_token}`
  }

  // Add dev-mode header if enabled
  const devMode = localStorage.getItem('dev_mode') === 'true'
  if (devMode) {
    config.headers['X-Dev-Mode'] = 'true'
  }

  return config
})

// Auth - now just a wrapper, actual auth is handled by Supabase client
export async function guestLogin(): Promise<GuestAuthResponse> {
  // This function is now deprecated - auth is handled by useAuth hook with Supabase client
  // Keeping for backwards compatibility
  throw new Error('Use Supabase client directly via useAuth hook')
}

// Health
export async function healthCheck(): Promise<{ status: string }> {
  const { data } = await api.get<{ status: string }>('/health')
  return data
}

// Planets
export async function listPlanets(): Promise<Planet[]> {
  const { data } = await api.get<Planet[]>('/planets')
  return data
}

export async function getPlanet(id: string): Promise<Planet> {
  const { data } = await api.get<Planet>(`/planets/${id}`)
  return data
}

// Buildings
export async function listBuildings(planetId: string): Promise<BuildingWithType[]> {
  const { data } = await api.get<BuildingWithType[]>(`/planets/${planetId}/buildings`)
  return data
}

export async function constructBuilding(
  planetId: string,
  buildingType: string,
  gridCol: number,
  gridRow: number,
): Promise<BuildingResponse> {
  const { data } = await api.post<BuildingResponse>(`/planets/${planetId}/buildings`, {
    building_type: buildingType,
    grid_col: gridCol,
    grid_row: gridRow,
  })
  return data
}

export async function moveBuilding(
  planetId: string,
  buildingId: string,
  gridCol: number,
  gridRow: number,
): Promise<void> {
  await api.put(`/planets/${planetId}/buildings/${buildingId}/move`, {
    grid_col: gridCol,
    grid_row: gridRow,
  })
}

export async function upgradeBuilding(
  planetId: string,
  buildingId: string,
): Promise<BuildingResponse> {
  const { data } = await api.post<BuildingResponse>(
    `/planets/${planetId}/buildings/${buildingId}/upgrade`,
  )
  return data
}

// Resources
export async function getResources(planetId: string): Promise<ResourcesResponse> {
  const { data } = await api.get<ResourcesResponse>(`/planets/${planetId}/resources`)
  return data
}

export async function collectResources(planetId: string): Promise<CollectResponse> {
  const { data } = await api.post<CollectResponse>(`/planets/${planetId}/resources/collect`)
  return data
}

export async function collectWarehouse(): Promise<CollectResponse> {
  const { data } = await api.post<CollectResponse>('/resources/collect-warehouse')
  return data
}

// Building Types
export async function getBuildingTypes(): Promise<BuildingTypeData[]> {
  const { data } = await api.get<BuildingTypeData[]>('/building-types')
  return data
}

// Cancel upgrade
export async function cancelUpgrade(planetId: string, buildingId: string): Promise<void> {
  await api.post(`/planets/${planetId}/buildings/${buildingId}/cancel`)
}

// ============ Phase 2: Hull & Module Reference ============

export async function getHullTypes(): Promise<HullType[]> {
  const { data } = await api.get<HullType[]>('/hull-types')
  return data
}

export async function getModuleTypes(category?: string): Promise<ModuleType[]> {
  const params = category ? { category } : {}
  const { data } = await api.get<ModuleType[]>('/module-types', { params })
  return data
}

// ============ Phase 2: Ship Designs ============

export async function listShipDesigns(): Promise<ShipDesign[]> {
  const { data } = await api.get<ShipDesign[]>('/ship-designs')
  return data
}

export async function createShipDesign(req: CreateShipDesignRequest): Promise<{ design: ShipDesign }> {
  const { data } = await api.post<{ design: ShipDesign }>('/ship-designs', req)
  return data
}

export async function updateShipDesign(id: string, req: CreateShipDesignRequest): Promise<{ design: ShipDesign }> {
  const { data } = await api.put<{ design: ShipDesign }>(`/ship-designs/${id}`, req)
  return data
}

export async function deleteShipDesign(id: string): Promise<void> {
  await api.delete(`/ship-designs/${id}`)
}

export async function getShipDesignStats(id: string): Promise<ShipDesign> {
  const { data } = await api.get<ShipDesign>(`/ship-designs/${id}/stats`)
  return data
}

// ============ Phase 2: Ship Factory ============

export async function getShipFactory(): Promise<ShipFactoryStatus> {
  const { data } = await api.get<ShipFactoryStatus>('/ship-factory')
  return data
}

export async function getShipFactorySlots(): Promise<ProductionSlot[]> {
  const { data } = await api.get<ProductionSlot[]>('/ship-factory/slots')
  return data
}

export async function buildShips(req: BuildShipRequest): Promise<BuildShipResponse> {
  const { data } = await api.post<BuildShipResponse>('/ship-factory/build', req)
  return data
}

export async function cancelShipBuild(slot: number): Promise<void> {
  await api.post(`/ship-factory/cancel/${slot}`)
}

// ============ Phase 2: Blueprints ============

export async function listBlueprints(): Promise<Blueprint[]> {
  const { data } = await api.get<Blueprint[]>('/blueprints')
  return data
}

export async function listMyBlueprints(): Promise<PlayerBlueprint[]> {
  const { data } = await api.get<PlayerBlueprint[]>('/blueprints/mine')
  return data
}

export async function activateBlueprint(id: number): Promise<void> {
  await api.post(`/blueprints/${id}/activate`)
}

export async function researchBlueprint(id: number): Promise<void> {
  await api.post(`/blueprints/${id}/research`)
}

export async function getActiveBlueprintResearch(): Promise<ActiveBlueprintResearch | null> {
  const { data } = await api.get<ActiveBlueprintResearch | null>('/blueprint-research/active')
  return data
}

// ============ Phase 2: Fleets ============

export async function listFleets(): Promise<Fleet[]> {
  const { data } = await api.get<Fleet[]>('/fleets')
  return data
}

export async function createFleet(req: CreateFleetRequest): Promise<Fleet> {
  const { data } = await api.post<Fleet>('/fleets', req)
  return data
}

export async function updateFleet(id: string, updates: Partial<CreateFleetRequest>): Promise<Fleet> {
  const { data } = await api.put<Fleet>(`/fleets/${id}`, updates)
  return data
}

export async function deleteFleet(id: string): Promise<void> {
  await api.delete(`/fleets/${id}`)
}

export async function assignStack(fleetId: string, req: AssignStackRequest): Promise<AssignStackResponse> {
  const { data } = await api.post<AssignStackResponse>(`/fleets/${fleetId}/assign-stack`, req)
  return data
}

export async function removeStack(fleetId: string, row: number, col: number): Promise<void> {
  await api.post(`/fleets/${fleetId}/remove-stack`, { grid_row: row, grid_col: col })
}

export async function dismissFleet(id: string): Promise<void> {
  await api.post(`/fleets/${id}/dismiss`)
}

// ============ Phase 2: Instances ============

export async function listInstances(): Promise<Instance[]> {
  const { data } = await api.get<Instance[]>('/instances')
  return data
}

export async function getInstanceDetail(id: number): Promise<InstanceDetail> {
  const { data } = await api.get<InstanceDetail>(`/instances/${id}`)
  return data
}

export async function attemptInstance(id: number, fleetIds: string[]): Promise<InstanceAttemptResponse> {
  const { data } = await api.post<InstanceAttemptResponse>(`/instances/${id}/attempt`, { fleet_ids: fleetIds })
  return data
}

export async function getInstanceProgress(): Promise<InstanceProgress[]> {
  const { data } = await api.get<InstanceProgress[]>('/instances/progress')
  return data
}

// ============ Phase 2: Spacedock ============

export async function getSpacedock(): Promise<SpacedockStatus> {
  const { data } = await api.get<SpacedockStatus>('/spacedock')
  return data
}

export async function getSpacedockRepairs(): Promise<SpacedockRepair[]> {
  const { data } = await api.get<SpacedockRepair[]>('/spacedock/repairs')
  return data
}

export async function startRepair(): Promise<void> {
  await api.post('/spacedock/repair')
}

export async function accelerateRepair(): Promise<void> {
  await api.post('/spacedock/accelerate')
}

// ============ Quests ============

export async function getQuests(): Promise<QuestsResponse> {
  const { data } = await api.get<QuestsResponse>('/quests')
  return data
}

export async function getDailyQuests(): Promise<DailyQuestsResponse> {
  const { data } = await api.get<DailyQuestsResponse>('/quests/daily')
  return data
}

export async function claimQuest(questId: string): Promise<ClaimQuestResponse> {
  const { data } = await api.post<ClaimQuestResponse>(`/quests/${questId}/claim`)
  return data
}

export async function claimDailyTier(tier: string): Promise<ClaimDailyTierResponse> {
  const { data } = await api.post<ClaimDailyTierResponse>('/quests/daily/claim-tier', { tier })
  return data
}

// ============ Research ============

export async function getResearch(): Promise<ResearchAllResponse> {
  const { data } = await api.get<ResearchAllResponse>('/research')
  return data
}

export async function getResearchTree(tree: string): Promise<ResearchTreeResponse> {
  const { data } = await api.get<ResearchTreeResponse>(`/research/trees/${tree}`)
  return data
}

export async function getActiveResearch(): Promise<ActiveResearch | null> {
  const { data } = await api.get<ActiveResearch | null>('/research/active')
  return data
}

export async function startResearch(techTypeId: number): Promise<StartResearchResponse> {
  const { data } = await api.post<StartResearchResponse>('/research/start', { tech_type_id: techTypeId })
  return data
}

export async function cancelResearch(techTypeId: number): Promise<CancelResearchResponse> {
  const { data } = await api.post<CancelResearchResponse>('/research/cancel', { tech_type_id: techTypeId })
  return data
}

export async function speedupResearch(techTypeId: number, speedupMinutes: number): Promise<SpeedupResearchResponse> {
  const { data } = await api.post<SpeedupResearchResponse>('/research/speedup', {
    tech_type_id: techTypeId,
    speedup_minutes: speedupMinutes,
  })
  return data
}

// ============ Recycling Plant ============

export async function listRecyclingJobs(): Promise<RecyclingJob[]> {
  const { data } = await api.get<RecyclingJob[]>('/recycling-plant/jobs')
  return data
}

export async function startRecycle(shipInstanceId: string): Promise<StartRecycleResponse> {
  const { data } = await api.post<StartRecycleResponse>('/recycling-plant/recycle', {
    ship_instance_id: shipInstanceId,
  })
  return data
}

export async function collectRecycle(jobId: string): Promise<CollectRecycleResponse> {
  const { data } = await api.post<CollectRecycleResponse>(`/recycling-plant/collect/${jobId}`)
  return data
}

export async function cancelRecycle(jobId: string): Promise<void> {
  await api.delete(`/recycling-plant/jobs/${jobId}`)
}

export async function listAvailableShipsForRecycling(): Promise<AvailableShip[]> {
  // Get all ship instances not in fleets
  const { data } = await api.get<AvailableShip[]>('/ship-instances/available')
  return data
}

// Combat Reports
export interface CombatReport {
  id: string
  attacker_id: string
  defender_id: string
  combat_type: string
  result: string
  total_rounds: number
  loot_json?: string
  rounds_json?: string
  he3_consumed: number
  created_at: string
}

export async function listCombatReports(): Promise<CombatReport[]> {
  const { data } = await api.get<CombatReport[]>('/combat-reports')
  return data
}

export async function getCombatReport(id: string): Promise<CombatReport> {
  const { data } = await api.get<CombatReport>(`/combat-reports/${id}`)
  return data
}

// ============ Chat System ============

export async function getChatMessages(channel: 'world' | 'alliance' = 'world', before?: string): Promise<ChatMessage[]> {
  const params = new URLSearchParams({ channel })
  if (before) params.append('before', before)
  const { data } = await api.get<ChatMessage[]>(`/chat/messages?${params.toString()}`)
  return data
}

export async function sendChatMessage(message: string, channel: 'world' | 'alliance' = 'world'): Promise<SendMessageResponse> {
  const { data } = await api.post<SendMessageResponse>('/chat/send', { message, channel })
  return data
}

// ============ PvP Combat ============

export async function searchPlanets(query: string): Promise<PvPSearchResult[]> {
  const { data} = await api.get<PvPSearchResult[]>(`/pvp/search?query=${encodeURIComponent(query)}`)
  return data
}

export async function attackPlanet(req: AttackPlanetRequest): Promise<AttackPlanetResponse> {
  const { data } = await api.post<AttackPlanetResponse>('/pvp/attack', req)
  return data
}

// ============ Commanders ============

export interface Commander {
  id: string
  player_id: string
  name: string
  rarity: 'common' | 'skill' | 'super'
  star_rank: number
  accuracy: number
  dodge: number
  speed: number
  electron: number
  weapon_expertise?: string
  ship_expertise?: string
  is_deployed: boolean
  created_at: string
  updated_at: string
}

export interface RecruitResponse {
  commander: Commander
  is_duplicate: boolean
  message: string
}

export interface MergeResponse {
  new_star_rank: number
  duplicates_used: number
  remaining_duplicates: number
  message: string
}

export async function recruitCommander(): Promise<RecruitResponse> {
  const { data } = await api.post<RecruitResponse>('/commanders/recruit')
  return data
}

export async function listCommanders(): Promise<Commander[]> {
  const { data } = await api.get<Commander[]>('/commanders')
  return data
}

export async function mergeCommander(commanderId: string, quantity: number): Promise<MergeResponse> {
  const { data } = await api.post<MergeResponse>('/commanders/merge', {
    commander_id: commanderId,
    quantity,
  })
  return data
}

export async function assignCommander(fleetId: string, commanderId: string): Promise<{ message: string }> {
  const { data } = await api.post<{ message: string }>(`/fleets/${fleetId}/assign-commander`, {
    commander_id: commanderId,
  })
  return data
}

export async function unassignCommander(fleetId: string): Promise<{ message: string }> {
  const { data } = await api.post<{ message: string }>(`/fleets/${fleetId}/unassign-commander`)
  return data
}

export async function dismissCommander(commanderId: string): Promise<{ message: string }> {
  const { data } = await api.delete<{ message: string }>(`/commanders/${commanderId}`)
  return data
}

// ============ Inventory ============

export async function getInventory(): Promise<InventoryItem[]> {
  const { data } = await api.get<InventoryItem[]>('/inventory')
  return data
}

export async function useItem(itemId: string): Promise<UseItemResponse> {
  const { data } = await api.post<UseItemResponse>(`/inventory/${itemId}/use`)
  return data
}

// ============ Combat Reports Helpers ============

export interface LootData {
  metal: number
  he3: number
  gold: number
}

export interface AttackData {
  AttackerStackID: string
  DefenderStackID: string
  AttackerSide: string
  DefenderSide: string
  Hit: boolean
  Damage: number
  ShieldDamage: number
  StructureDamage: number
  ShipsDestroyed: number
}

export interface RoundData {
  RoundNumber: number
  Attacks: AttackData[]
  Casualties: Record<string, number>
}

export const combatReportsAPI = {
  list: listCombatReports,
  get: getCombatReport,

  parseLoot: (report: CombatReport): LootData | null => {
    if (!report.loot_json) return null
    try {
      return JSON.parse(report.loot_json) as LootData
    } catch {
      return null
    }
  },

  parseRounds: (report: CombatReport): RoundData[] | null => {
    if (!report.rounds_json) return null
    try {
      return JSON.parse(report.rounds_json) as RoundData[]
    } catch {
      return null
    }
  },
}

export default api
