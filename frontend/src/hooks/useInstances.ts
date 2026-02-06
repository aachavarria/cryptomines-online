import { useState, useEffect, useCallback } from 'react'
import {
  listInstances,
  getInstanceDetail,
  attemptInstance,
  getInstanceProgress,
} from '../services/api.ts'
import type { Instance, InstanceDetail, InstanceProgress, InstanceAttemptResponse } from '../types'

export function useInstances() {
  const [instances, setInstances] = useState<Instance[]>([])
  const [progress, setProgress] = useState<InstanceProgress[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    try {
      const [inst, prog] = await Promise.all([listInstances(), getInstanceProgress()])
      setInstances(inst)
      setProgress(prog)
      setError(null)
    } catch {
      setError('Failed to load instances')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    refresh()
  }, [refresh])

  const getDetail = useCallback(async (id: number): Promise<InstanceDetail> => {
    return getInstanceDetail(id)
  }, [])

  const attempt = useCallback(async (id: number, fleetIds: string[]): Promise<InstanceAttemptResponse> => {
    try {
      const result = await attemptInstance(id, fleetIds)
      await refresh()
      return result
    } catch {
      setError('Failed to attempt instance')
      throw new Error('attempt failed')
    }
  }, [refresh])

  const isCompleted = useCallback((instanceId: number): boolean => {
    return progress.some(p => p.instance_id === instanceId && p.completed)
  }, [progress])

  return { instances, progress, loading, error, getDetail, attempt, isCompleted, refresh }
}
