import { useState, useEffect } from 'react'
import { listCommanders, type Commander } from '../api/commanders'

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
    } catch (err: any) {
      setError(err.message || 'Failed to load commanders')
      console.error('Failed to fetch commanders:', err)
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
