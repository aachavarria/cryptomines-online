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
} from '../types'

const api = axios.create({
  baseURL: '/api',
})

// Attach JWT to every request if available
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Auth
export async function guestLogin(): Promise<GuestAuthResponse> {
  const { data } = await api.post<GuestAuthResponse>('/auth/guest')
  localStorage.setItem('token', data.token)
  localStorage.setItem('player_id', data.player.id)
  return data
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
): Promise<BuildingResponse> {
  const { data } = await api.post<BuildingResponse>(`/planets/${planetId}/buildings`, {
    building_type: buildingType,
  })
  return data
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

// Building Types
export async function getBuildingTypes(): Promise<BuildingTypeData[]> {
  const { data } = await api.get<BuildingTypeData[]>('/building-types')
  return data
}

// Cancel upgrade
export async function cancelUpgrade(planetId: string, buildingId: string): Promise<any> {
  const { data } = await api.post(`/planets/${planetId}/buildings/${buildingId}/cancel`)
  return data
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

export default api
