import { useState, useEffect, useCallback } from 'react'
import {
  getResearch,
  getResearchTree,
  startResearch,
  cancelResearch,
  speedupResearch,
} from '../services/api.ts'
import type {
  TechTree,
  ResearchTreeResponse,
  TechWithProgress,
  ActiveResearch,
  StartResearchResponse,
  CancelResearchResponse,
  SpeedupResearchResponse,
} from '../types'

export function useResearch() {
  const [trees, setTrees] = useState<Record<string, TechWithProgress[]>>({})
  const [active, setActive] = useState<ActiveResearch | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refreshAll = useCallback(async () => {
    try {
      setLoading(true)
      const data = await getResearch()
      const treeMap: Record<string, TechWithProgress[]> = {}
      for (const t of data.trees) {
        treeMap[t.tree] = t.techs
      }
      setTrees(treeMap)
      setActive(data.active ?? null)
      setError(null)
    } catch {
      setError('Failed to load research data')
    } finally {
      setLoading(false)
    }
  }, [])

  const refreshTree = useCallback(async (tree: TechTree) => {
    try {
      const data: ResearchTreeResponse = await getResearchTree(tree)
      setTrees(prev => ({ ...prev, [data.tree]: data.techs }))
      setError(null)
    } catch {
      setError('Failed to load tree data')
    }
  }, [])

  useEffect(() => {
    refreshAll()
  }, [refreshAll])

  const start = useCallback(async (techTypeId: number): Promise<StartResearchResponse> => {
    const result = await startResearch(techTypeId)
    await refreshAll()
    return result
  }, [refreshAll])

  const cancel = useCallback(async (techTypeId: number): Promise<CancelResearchResponse> => {
    const result = await cancelResearch(techTypeId)
    await refreshAll()
    return result
  }, [refreshAll])

  const speedup = useCallback(async (techTypeId: number, minutes: number): Promise<SpeedupResearchResponse> => {
    const result = await speedupResearch(techTypeId, minutes)
    await refreshAll()
    return result
  }, [refreshAll])

  return {
    trees,
    active,
    loading,
    error,
    refreshAll,
    refreshTree,
    start,
    cancel,
    speedup,
  }
}
