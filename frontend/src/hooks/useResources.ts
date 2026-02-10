import { useEffect, useCallback } from 'react'
import { getResources, collectResources, collectWarehouse, listBuildings } from '../services/api.ts'
import { useGameContext } from '../contexts/GameContext.tsx'

export function useResources() {
  const { state, dispatch } = useGameContext()
  const planetId = state.currentPlanet?.id

  const fetchResources = useCallback(async () => {
    if (!planetId) return
    try {
      const r = await getResources(planetId)
      dispatch({ type: 'SET_RESOURCES', payload: r })
    } catch {
      // Silent fail on refresh
    }
  }, [planetId, dispatch])

  // Auto-refresh every 30s
  useEffect(() => {
    if (!planetId) return
    fetchResources()
    const interval = setInterval(fetchResources, 30000)
    return () => clearInterval(interval)
  }, [planetId, fetchResources])

  const collect = useCallback(async () => {
    if (!planetId) return
    try {
      const result = await collectResources(planetId)
      dispatch({
        type: 'SET_RESOURCES',
        payload: {
          ...result.resources,
          pending_metal: 0,
          pending_he3: 0,
          pending_gold: 0,
        },
      })
      // Refresh buildings to detect completed upgrades
      const buildings = await listBuildings(planetId)
      dispatch({ type: 'SET_BUILDINGS', payload: buildings })
    } catch {
      dispatch({ type: 'SET_ERROR', payload: 'Failed to collect resources' })
    }
  }, [planetId, dispatch])

  const collectWarehouseResources = useCallback(async () => {
    if (!planetId) return null
    try {
      const result = await collectWarehouse()
      dispatch({
        type: 'SET_RESOURCES',
        payload: {
          ...result.resources,
          warehouse_metal: 0,
          warehouse_he3: 0,
          warehouse_gold: 0,
        },
      })
      return result.collected
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to collect warehouse'
      dispatch({ type: 'SET_ERROR', payload: message })
      return null
    }
  }, [planetId, dispatch])

  return {
    resources: state.resources,
    collect,
    collectWarehouse: collectWarehouseResources,
    refreshResources: fetchResources,
  }
}
