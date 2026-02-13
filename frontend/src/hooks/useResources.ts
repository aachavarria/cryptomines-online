import { useEffect, useCallback } from 'react'
import { getResources, collectResources, listBuildings } from '../services/api.ts'
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

  // Unified collect: flushes production → warehouse → main
  const collect = useCallback(async () => {
    if (!planetId) return
    try {
      const result = await collectResources(planetId)
      dispatch({
        type: 'SET_RESOURCES',
        payload: {
          ...result.resources,
          warehouse_capacity: state.resources?.warehouse_capacity || 0,
          pending_metal: 0,
          pending_he3: 0,
          pending_gold: 0,
        },
      })
      // Refresh buildings to detect completed upgrades
      const buildings = await listBuildings(planetId)
      dispatch({ type: 'SET_BUILDINGS', payload: buildings })
      return result.collected
    } catch {
      dispatch({ type: 'SET_ERROR', payload: 'Failed to collect resources' })
      return null
    }
  }, [planetId, dispatch, state.resources])

  return {
    resources: state.resources,
    collect,
    refreshResources: fetchResources,
  }
}
