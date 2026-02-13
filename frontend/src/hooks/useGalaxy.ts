import { useState, useEffect, useCallback } from 'react'
import { getGalaxyMap, attackRBP } from '../services/api'
import type { GalaxyZone, AttackRBPRequest, AttackRBPResponse } from '../types/corps'
import { useGameContext } from '../contexts/GameContext'

export function useGalaxy() {
  const { dispatch } = useGameContext()
  const [zones, setZones] = useState<GalaxyZone[]>([])
  const [selectedZone, setSelectedZone] = useState<GalaxyZone | null>(null)
  const [loading, setLoading] = useState(false)
  const [attacking, setAttacking] = useState(false)

  const fetchMap = useCallback(async () => {
    setLoading(true)
    try {
      const data = await getGalaxyMap()
      setZones(data.zones)
    } catch (error: unknown) {
      const axiosErr = error as { response?: { data?: { error?: string } } }
      const message = axiosErr.response?.data?.error || 'Failed to fetch galaxy map'
      dispatch({ type: 'SET_ERROR', payload: message })
    } finally {
      setLoading(false)
    }
  }, [dispatch])

  const attack = useCallback(
    async (rbpPlanetId: string, req: AttackRBPRequest): Promise<AttackRBPResponse | null> => {
      setAttacking(true)
      try {
        const response = await attackRBP(rbpPlanetId, req)
        await fetchMap() // Refresh map after attack
        return response
      } catch (error: unknown) {
        const axiosErr = error as { response?: { data?: { error?: string } } }
        const message = axiosErr.response?.data?.error || 'Failed to attack RBP'
        dispatch({ type: 'SET_ERROR', payload: message })
        return null
      } finally {
        setAttacking(false)
      }
    },
    [dispatch, fetchMap],
  )

  const selectZone = useCallback((zone: GalaxyZone | null) => {
    setSelectedZone(zone)
  }, [])

  // Initial fetch
  useEffect(() => {
    fetchMap()
  }, [fetchMap])

  return {
    zones,
    selectedZone,
    loading,
    attacking,
    fetchMap,
    attack,
    selectZone,
    refresh: fetchMap,
  }
}
