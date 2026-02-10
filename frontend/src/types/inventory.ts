import type { Resource } from './game'

export interface InventoryItem {
  id: string
  item_key: string
  quantity: number
  display_name: string
  category: 'resource_pack' | 'boost' | 'battle' | 'blueprint' | 'commander'
  description: string
  icon_name?: string
  acquired_at: string
}

export interface UseItemResponse {
  effect: string // Human-readable effect description
  resources?: Resource // Updated resources (if applicable)
}
