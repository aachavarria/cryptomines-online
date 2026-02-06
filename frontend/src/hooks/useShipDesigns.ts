import { useState, useEffect, useCallback } from 'react'
import {
  listShipDesigns,
  createShipDesign,
  deleteShipDesign,
  getHullTypes,
  getModuleTypes,
} from '../services/api.ts'
import type { ShipDesign, HullType, ModuleType, CreateShipDesignRequest } from '../types'

export function useShipDesigns() {
  const [designs, setDesigns] = useState<ShipDesign[]>([])
  const [hullTypes, setHullTypes] = useState<HullType[]>([])
  const [moduleTypes, setModuleTypes] = useState<ModuleType[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    try {
      const [d, h, m] = await Promise.all([
        listShipDesigns(),
        getHullTypes(),
        getModuleTypes(),
      ])
      setDesigns(d)
      setHullTypes(h)
      setModuleTypes(m)
      setError(null)
    } catch {
      setError('Failed to load ship designs')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    refresh()
  }, [refresh])

  const create = useCallback(async (req: CreateShipDesignRequest) => {
    try {
      await createShipDesign(req)
      await refresh()
    } catch {
      setError('Failed to create design')
      throw new Error('create failed')
    }
  }, [refresh])

  const remove = useCallback(async (id: string) => {
    try {
      await deleteShipDesign(id)
      await refresh()
    } catch {
      setError('Cannot delete design (ships may exist)')
    }
  }, [refresh])

  return { designs, hullTypes, moduleTypes, loading, error, create, remove, refresh }
}
