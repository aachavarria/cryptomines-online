import { useEffect, useState } from 'react'
import { guestLogin } from '../services/api'
import type { GuestAuthResponse } from '../types'

export function useAuth() {
  const [player, setPlayer] = useState<GuestAuthResponse['player'] | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const existingToken = localStorage.getItem('token')
    if (existingToken) {
      // Already logged in - retrieve stored player info
      const storedId = localStorage.getItem('player_id')
      if (storedId) {
        setPlayer({ id: storedId, anonymous_id: '', level: 1, created_at: '' })
        setLoading(false)
        return
      }
    }

    // Auto guest login
    guestLogin()
      .then((resp) => {
        setPlayer(resp.player)
      })
      .catch(() => {
        setError('Failed to connect to server')
      })
      .finally(() => {
        setLoading(false)
      })
  }, [])

  return { player, loading, error }
}
