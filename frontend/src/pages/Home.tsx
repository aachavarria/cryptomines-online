import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth.ts'
import { listPlanets } from '../services/api.ts'

export default function Home() {
  const { player, loading, error: authError } = useAuth()
  const navigate = useNavigate()
  const [status, setStatus] = useState<string>('Initializing...')

  useEffect(() => {
    if (loading) {
      setStatus('Connecting to server...')
      return
    }
    if (authError) {
      setStatus(authError)
      return
    }
    if (!player) return

    setStatus('Loading your colony...')

    listPlanets()
      .then((planets) => {
        if (planets.length > 0) {
          const homeworld = planets.find((p) => p.is_homeworld) ?? planets[0]
          localStorage.setItem('current_planet_id', homeworld.id)
          navigate(`/planet/${homeworld.id}`, { replace: true })
        } else {
          setStatus('No planets found. Something went wrong.')
        }
      })
      .catch(() => {
        setStatus('Failed to load planets. Is the backend running?')
      })
  }, [player, loading, authError, navigate])

  return (
    <div className="loading-screen">
      <h2>Cryptomines Online</h2>
      <div className="loading-spinner" />
      <p>{status}</p>
      {authError && (
        <button className="hud-collect-btn" onClick={() => window.location.reload()}>
          Retry
        </button>
      )}
    </div>
  )
}
