import { useState, useEffect, useCallback } from 'react'
import { listBlueprints, listMyBlueprints, activateBlueprint, researchBlueprint, getActiveBlueprintResearch } from '../services/api.ts'
import type { Blueprint, PlayerBlueprint, ActiveBlueprintResearch } from '../types'

export function useBlueprints() {
  const [allBlueprints, setAllBlueprints] = useState<Blueprint[]>([])
  const [myBlueprints, setMyBlueprints] = useState<PlayerBlueprint[]>([])
  const [activeResearch, setActiveResearch] = useState<ActiveBlueprintResearch | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    try {
      const [all, mine, active] = await Promise.all([
        listBlueprints(),
        listMyBlueprints(),
        getActiveBlueprintResearch()
      ])
      setAllBlueprints(all)
      setMyBlueprints(mine)
      setActiveResearch(active)
      setError(null)
    } catch {
      setError('Failed to load blueprints')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    refresh()
  }, [refresh])

  const activate = useCallback(async (id: number) => {
    try {
      await activateBlueprint(id)
      await refresh()
    } catch {
      setError('Failed to activate blueprint')
    }
  }, [refresh])

  const hasActivated = useCallback((blueprintId: number): boolean => {
    return myBlueprints.some(bp => bp.blueprint_id === blueprintId && bp.is_activated)
  }, [myBlueprints])

  const hasOwned = useCallback((blueprintId: number): boolean => {
    return myBlueprints.some(bp => bp.blueprint_id === blueprintId)
  }, [myBlueprints])

  const hasHullBlueprint = useCallback((hullTypeId: number): boolean => {
    return myBlueprints.some(bp => bp.blueprint_type === 'hull' && bp.hull_type_id === hullTypeId && bp.is_activated)
  }, [myBlueprints])

  const hasModuleBlueprint = useCallback((moduleTypeId: number): boolean => {
    return myBlueprints.some(bp => bp.blueprint_type === 'module' && bp.module_type_id === moduleTypeId && bp.is_activated)
  }, [myBlueprints])

  const research = useCallback(async (id: number) => {
    try {
      await researchBlueprint(id)
      await refresh()
    } catch {
      setError('Failed to start research')
      throw new Error('Failed to start research')
    }
  }, [refresh])

  return { allBlueprints, myBlueprints, activeResearch, loading, error, activate, research, hasActivated, hasOwned, hasHullBlueprint, hasModuleBlueprint, refresh }
}
