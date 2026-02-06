import { useEffect, useCallback, useRef } from 'react'
import {
  listBuildings,
  upgradeBuilding,
  constructBuilding,
} from '../services/api.ts'
import { useGameContext } from '../contexts/GameContext.tsx'

export function useBuildings() {
  const { state, dispatch } = useGameContext()
  const planetId = state.currentPlanet?.id
  const refreshTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const fetchBuildings = useCallback(async () => {
    if (!planetId) return
    try {
      const b = await listBuildings(planetId)
      dispatch({ type: 'SET_BUILDINGS', payload: b })
    } catch {
      dispatch({ type: 'SET_ERROR', payload: 'Failed to load buildings' })
    }
  }, [planetId, dispatch])

  useEffect(() => {
    fetchBuildings()
  }, [fetchBuildings])

  // Auto-refresh when the nearest upgrade/construction finishes
  useEffect(() => {
    if (refreshTimerRef.current) {
      clearTimeout(refreshTimerRef.current)
      refreshTimerRef.current = null
    }

    const upgrading = state.buildings.filter(b => b.is_upgrading && b.upgrade_finish_at)
    if (upgrading.length === 0) return

    // Find the soonest finish time
    const now = Date.now()
    let soonest = Infinity
    for (const b of upgrading) {
      const finishMs = new Date(b.upgrade_finish_at!).getTime()
      const remaining = finishMs - now
      if (remaining < soonest) soonest = remaining
    }

    // Schedule refresh slightly after completion (1s buffer for server clock drift)
    const delay = Math.max(1000, soonest + 1000)
    refreshTimerRef.current = setTimeout(() => {
      fetchBuildings()
    }, delay)

    return () => {
      if (refreshTimerRef.current) {
        clearTimeout(refreshTimerRef.current)
      }
    }
  }, [state.buildings, fetchBuildings])

  const upgrade = useCallback(async (buildingId: string) => {
    if (!planetId) return
    try {
      await upgradeBuilding(planetId, buildingId)
      await fetchBuildings()
    } catch {
      dispatch({ type: 'SET_ERROR', payload: 'Failed to upgrade building' })
    }
  }, [planetId, fetchBuildings, dispatch])

  const construct = useCallback(async (buildingType: string) => {
    if (!planetId) return
    try {
      await constructBuilding(planetId, buildingType)
      await fetchBuildings()
    } catch {
      dispatch({ type: 'SET_ERROR', payload: 'Failed to construct building' })
    }
  }, [planetId, fetchBuildings, dispatch])

  return {
    buildings: state.buildings,
    upgrade,
    construct,
    refreshBuildings: fetchBuildings,
  }
}
