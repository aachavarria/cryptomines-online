import { useState, useEffect, useCallback } from 'react'
import { listAvailableShipsForRecycling } from '../services/api.ts'
import type { AvailableShip } from '../types'

export function useShipInventory() {
  const [ships, setShips] = useState<AvailableShip[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    try {
      const data = await listAvailableShipsForRecycling()
      setShips(data)
      setError(null)
    } catch {
      setError('Failed to load ship inventory')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    refresh()
  }, [refresh])

  return { ships, loading, error, refresh }
}
