import { useState, useEffect } from 'react'
import { listCommanders, type Commander } from '../services/api'

export function useCommanders() {
  const [commanders, setCommanders] = useState<Commander[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchCommanders = async () => {
    try {
      setLoading(true)
      const data = await listCommanders()
      setCommanders(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load commanders')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchCommanders()
  }, [])

  return {
    commanders,
    loading,
    error,
    refresh: fetchCommanders,
  }
}
