import { useState, useEffect, useCallback } from 'react'
import {
  listRecyclingJobs,
  startRecycle,
  collectRecycle,
  cancelRecycle,
  listAvailableShipsForRecycling,
} from '../services/api.ts'
import type { RecyclingJob, AvailableShip } from '../types'
import { useGameContext } from '../contexts/GameContext.tsx'

export function useRecycling() {
  const { dispatch } = useGameContext()
  const [jobs, setJobs] = useState<RecyclingJob[]>([])
  const [availableShips, setAvailableShips] = useState<AvailableShip[]>([])
  const [loading, setLoading] = useState(false)

  const fetchJobs = useCallback(async () => {
    try {
      const data = await listRecyclingJobs()
      setJobs(data)
    } catch (error) {
      console.error('Failed to fetch recycling jobs:', error)
    }
  }, [])

  const fetchAvailableShips = useCallback(async () => {
    try {
      const data = await listAvailableShipsForRecycling()
      setAvailableShips(data)
    } catch (error) {
      // Endpoint not implemented yet (Task #70)
      // Silently fail and show empty list
      setAvailableShips([])
    }
  }, [])

  const recycle = useCallback(
    async (shipInstanceId: string) => {
      setLoading(true)
      try {
        await startRecycle(shipInstanceId)
        await fetchJobs()
        await fetchAvailableShips()
        return true
      } catch (error) {
        const message = error instanceof Error ? error.message : 'Failed to start recycling'
        dispatch({ type: 'SET_ERROR', payload: message })
        return false
      } finally {
        setLoading(false)
      }
    },
    [dispatch, fetchJobs, fetchAvailableShips]
  )

  const collect = useCallback(
    async (jobId: string) => {
      setLoading(true)
      try {
        const result = await collectRecycle(jobId)
        await fetchJobs()
        dispatch({
          type: 'SET_ERROR',
          payload: `Collected: ${result.metal_gained} Metal, ${result.he3_gained} He3, ${result.gold_gained} Gold`,
        })
        return true
      } catch (error) {
        const message = error instanceof Error ? error.message : 'Failed to collect'
        dispatch({ type: 'SET_ERROR', payload: message })
        return false
      } finally {
        setLoading(false)
      }
    },
    [dispatch, fetchJobs]
  )

  const cancel = useCallback(
    async (jobId: string) => {
      setLoading(true)
      try {
        await cancelRecycle(jobId)
        await fetchJobs()
        return true
      } catch (error) {
        const message = error instanceof Error ? error.message : 'Failed to cancel'
        dispatch({ type: 'SET_ERROR', payload: message })
        return false
      } finally {
        setLoading(false)
      }
    },
    [dispatch, fetchJobs]
  )

  // Auto-refresh every 10 seconds
  useEffect(() => {
    fetchJobs()
    fetchAvailableShips()
    const interval = setInterval(() => {
      fetchJobs()
    }, 10000)
    return () => clearInterval(interval)
  }, [fetchJobs, fetchAvailableShips])

  return {
    jobs,
    availableShips,
    loading,
    recycle,
    collect,
    cancel,
    refresh: fetchJobs,
  }
}
