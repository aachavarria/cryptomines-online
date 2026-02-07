import { useEffect, useState } from 'react'
import { supabase } from '../lib/supabase'
import type { User } from '@supabase/supabase-js'

interface PlayerInfo {
  id: string
  anonymous_id: string
  level: number
  created_at: string
}

export function useAuth() {
  const [player, setPlayer] = useState<PlayerInfo | null>(null)
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    // Check for existing session
    supabase.auth.getSession().then(({ data: { session } }) => {
      if (session?.user) {
        setUser(session.user)
        setPlayer({
          id: session.user.id,
          anonymous_id: `supabase_${session.user.id}`,
          level: 1,
          created_at: session.user.created_at,
        })
        setLoading(false)
      } else {
        // Auto guest login
        supabase.auth.signInAnonymously()
          .then(({ data, error }) => {
            if (error) throw error
            if (data.user) {
              setUser(data.user)
              setPlayer({
                id: data.user.id,
                anonymous_id: `supabase_${data.user.id}`,
                level: 1,
                created_at: data.user.created_at,
              })
            }
          })
          .catch((err) => {
            console.error('Auth error:', err)
            setError('Failed to connect to server')
          })
          .finally(() => {
            setLoading(false)
          })
      }
    })

    // Listen for auth changes (including token refresh)
    const { data: { subscription } } = supabase.auth.onAuthStateChange((_event, session) => {
      if (session?.user) {
        setUser(session.user)
        setPlayer({
          id: session.user.id,
          anonymous_id: `supabase_${session.user.id}`,
          level: 1,
          created_at: session.user.created_at,
        })
      } else {
        setUser(null)
        setPlayer(null)
      }
    })

    return () => {
      subscription.unsubscribe()
    }
  }, [])

  return { player, user, loading, error }
}
