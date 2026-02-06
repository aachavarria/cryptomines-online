import { useState, useEffect, useCallback } from 'react'
import { getSpacedock, getSpacedockRepairs, startRepair } from '../services/api.ts'
import type { SpacedockStatus, SpacedockRepair } from '../types'

export function useSpacedock() {
  const [status, setStatus] = useState<SpacedockStatus | null>(null)
  const [repairs, setRepairs] = useState<SpacedockRepair[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refresh = useCallback(async () => {
    try {
      const [s, r] = await Promise.all([getSpacedock(), getSpacedockRepairs()])
      setStatus(s)
      setRepairs(r)
      setError(null)
    } catch {
      setError('Failed to load Spacedock')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    refresh()
    const interval = setInterval(refresh, 10000)
    return () => clearInterval(interval)
  }, [refresh])

  const repair = useCallback(async () => {
    try {
      await startRepair()
      await refresh()
    } catch {
      setError('Failed to start repair')
    }
  }, [refresh])

  return { status, repairs, loading, error, repair, refresh }
}
