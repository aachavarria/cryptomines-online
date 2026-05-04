import { useState, useEffect } from 'react'
import { listItemTypes, type ItemTypeInfo } from '../services/api.ts'

export function useItemTypes() {
  const [itemTypes, setItemTypes] = useState<ItemTypeInfo[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false
    listItemTypes()
      .then(data => {
        if (!cancelled) {
          setItemTypes(data)
          setError(null)
        }
      })
      .catch(() => {
        if (!cancelled) setError('Failed to load item types')
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [])

  return { itemTypes, loading, error }
}
