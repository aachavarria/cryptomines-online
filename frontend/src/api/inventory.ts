import { apiRequest } from './api'
import type { InventoryItem, UseItemResponse } from '../types/inventory'

/**
 * Get player's inventory items
 */
export async function getInventory(): Promise<InventoryItem[]> {
  return apiRequest<InventoryItem[]>('GET', '/api/inventory')
}

/**
 * Use an item from inventory
 * @param itemId - The ID of the inventory item to use
 */
export async function useItem(itemId: string): Promise<UseItemResponse> {
  return apiRequest<UseItemResponse>('POST', `/api/inventory/${itemId}/use`)
}
