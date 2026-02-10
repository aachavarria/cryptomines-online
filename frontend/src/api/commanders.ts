import { apiRequest } from './api'

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

/**
 * Recruit a new commander (gacha)
 * Cost: 10,000 Gold
 */
export async function recruitCommander(): Promise<RecruitResponse> {
  return apiRequest<RecruitResponse>('POST', '/api/commanders/recruit')
}

/**
 * Get all owned commanders
 */
export async function listCommanders(): Promise<Commander[]> {
  return apiRequest<Commander[]>('GET', '/api/commanders')
}

/**
 * Merge commander (consume duplicates to increase star rank)
 */
export async function mergeCommander(commanderId: string, quantity: number): Promise<MergeResponse> {
  return apiRequest<MergeResponse>('POST', '/api/commanders/merge', {
    commander_id: commanderId,
    quantity,
  })
}

/**
 * Assign commander to fleet
 */
export async function assignCommander(fleetId: string, commanderId: string): Promise<{ message: string }> {
  return apiRequest<{ message: string }>('POST', `/api/fleets/${fleetId}/assign-commander`, {
    commander_id: commanderId,
  })
}

/**
 * Unassign commander from fleet
 */
export async function unassignCommander(fleetId: string): Promise<{ message: string }> {
  return apiRequest<{ message: string }>('POST', `/api/fleets/${fleetId}/unassign-commander`)
}

/**
 * Dismiss commander (delete)
 */
export async function dismissCommander(commanderId: string): Promise<{ message: string }> {
  return apiRequest<{ message: string }>('DELETE', `/api/commanders/${commanderId}`)
}
