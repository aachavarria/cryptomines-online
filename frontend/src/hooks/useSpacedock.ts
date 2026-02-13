import { useState, useEffect, useCallback, useRef } from 'react'
import { getSpacedock, getSpacedockRepairs, startRepair } from '../services/api.ts'
import type { SpacedockStatus, SpacedockRepair } from '../types'
import axios from 'axios'

export function useSpacedock() {
  const [status, setStatus] = useState<SpacedockStatus | null>(null)
  const [repairs, setRepairs] = useState<SpacedockRepair[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const notBuiltRef = useRef(false)

  const refresh = useCallback(async () => {
    if (notBuiltRef.current) return
    try {
      const [s, r] = await Promise.all([getSpacedock(), getSpacedockRepairs()])
      setStatus(s)
      setRepairs(r)
      setError(null)
    } catch (e) {
      if (axios.isAxiosError(e) && e.response?.status === 404) {
        notBuiltRef.current = true
      } else {
        setError('Failed to load Spacedock')
      }
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
