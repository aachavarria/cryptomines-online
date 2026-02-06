import { useState, useEffect, useCallback } from 'react'
import {
  listFleets,
  createFleet,
  updateFleet,
  deleteFleet,
  assignStack,
  removeStack,
} from '../services/api.ts'
import type { Fleet, CreateFleetRequest, AssignStackRequest } from '../types'

export function useFleets() {
  const [fleets, setFleets] = useState<Fleet[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    try {
      const f = await listFleets()
      setFleets(f)
      setError(null)
    } catch {
      setError('Failed to load fleets')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    refresh()
  }, [refresh])

  const create = useCallback(async (req: CreateFleetRequest): Promise<Fleet> => {
    try {
      const fleet = await createFleet(req)
      await refresh()
      return fleet
    } catch {
      setError('Failed to create fleet')
      throw new Error('create failed')
    }
  }, [refresh])

  const update = useCallback(async (id: string, updates: Partial<CreateFleetRequest>) => {
    try {
      await updateFleet(id, updates)
      await refresh()
    } catch {
      setError('Failed to update fleet')
    }
  }, [refresh])

  const remove = useCallback(async (id: string) => {
    try {
      await deleteFleet(id)
      await refresh()
    } catch {
      setError('Failed to disband fleet')
    }
  }, [refresh])

  const addStack = useCallback(async (fleetId: string, req: AssignStackRequest) => {
    try {
      const result = await assignStack(fleetId, req)
      await refresh()
      return result
    } catch {
      setError('Failed to assign ships')
      throw new Error('assign failed')
    }
  }, [refresh])

  const removeStackFromFleet = useCallback(async (fleetId: string, row: number, col: number) => {
    try {
      await removeStack(fleetId, row, col)
      await refresh()
    } catch {
      setError('Failed to remove stack')
    }
  }, [refresh])

  return { fleets, loading, error, create, update, remove, addStack, removeStackFromFleet, refresh }
}
