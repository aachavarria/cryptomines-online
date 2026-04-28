export interface Corp {
  id: string
  name: string
  tag: string
  leader_id: string
  level: number
  wealth: number
  max_members: number
  description: string
  created_at: string
  updated_at: string
}

export interface CorpBonuses {
  total_rbp_bonus: number
  rbp_count: number
  corp_level: number
}

export interface CorpWithBonuses {
  corp: Corp | null
  bonuses: CorpBonuses | null
  role: string | null
}

export interface CorpMember {
  player_id: string
  player_name: string
  role: 'leader' | 'officer' | 'member'
  contribution_points: number
  daily_contribution: number
  joined_at: string
}

export interface CreateCorpRequest {
  name: string
  tag: string
  description: string
}

export interface DonateRequest {
  metal: number
  he3: number
  gold: number
}

export interface CorpSearchResult {
  id: string
  name: string
  tag: string
  level: number
  member_count: number
  max_members: number
  description: string
}

export interface GalaxyZone {
  zone_x: number
  zone_y: number
  rbp_planet_id: string
  rbp_name: string
  rbp_level: number
  controlling_corp_id: string | null
  controlling_corp_name: string | null
  controlling_corp_tag: string | null
  protection_until: string | null
  is_protected: boolean
}

export interface GalaxyMapResponse {
  zones: GalaxyZone[]
}

export interface SectorPlanet {
  id: string
  name: string
  position_x: number
  position_y: number
  is_homeworld: boolean
  is_rbp: boolean
  rbp_level: number
  is_own: boolean
  owner_id?: string
  owner_name?: string
  controlling_corp?: {
    id: string
    name: string
    tag: string
  }
  protection_until?: string
  defense_strength: number
}

export interface GalaxySectorResponse {
  center_x: number
  center_y: number
  radius: number
  planets: SectorPlanet[]
}

export interface AttackRBPRequest {
  fleet_ids: string[]
}

export interface AttackRBPResponse {
  report_id: string
  result: 'attacker_win' | 'defender_win' | 'draw'
  total_rounds: number
  conquered: boolean
}
