import { useState, useEffect, useCallback } from 'react'
import {
  getShipFactory,
  getShipFactorySlots,
  buildShips,
  cancelShipBuild,
} from '../services/api.ts'
import type { ShipFactoryStatus, ProductionSlot, BuildShipRequest } from '../types'

export function useShipFactory() {
  const [factory, setFactory] = useState<ShipFactoryStatus | null>(null)
  const [slots, setSlots] = useState<ProductionSlot[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    try {
      const [f, s] = await Promise.all([getShipFactory(), getShipFactorySlots()])
      setFactory(f)
      setSlots(s)
      setError(null)
    } catch {
      setError('Failed to load Ship Factory')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    refresh()
    const interval = setInterval(refresh, 10000)
    return () => clearInterval(interval)
  }, [refresh])

  const build = useCallback(async (req: BuildShipRequest) => {
    try {
      await buildShips(req)
      await refresh()
    } catch {
      setError('Failed to start ship production')
      throw new Error('build failed')
    }
  }, [refresh])

  const cancel = useCallback(async (slot: number) => {
    try {
      await cancelShipBuild(slot)
      await refresh()
    } catch {
      setError('Failed to cancel production')
    }
  }, [refresh])

  return { factory, slots, loading, error, build, cancel, refresh }
}
