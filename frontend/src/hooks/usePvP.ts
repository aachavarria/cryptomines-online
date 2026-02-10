import { useState, useCallback } from 'react'
import { searchPlanets, attackPlanet } from '../services/api.ts'
import type { PvPSearchResult, AttackPlanetResponse } from '../types'
import { useGameContext } from '../contexts/GameContext.tsx'

export function usePvP() {
  const { dispatch } = useGameContext()
  const [searchResults, setSearchResults] = useState<PvPSearchResult[]>([])
  const [searching, setSearching] = useState(false)
  const [attacking, setAttacking] = useState(false)
  const [lastBattle, setLastBattle] = useState<AttackPlanetResponse | null>(null)

  const search = useCallback(async (query: string) => {
    if (!query.trim()) {
      setSearchResults([])
      return
    }

    setSearching(true)
    try {
      const results = await searchPlanets(query)
      setSearchResults(results)
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Search failed'
      dispatch({ type: 'SET_ERROR', payload: message })
      setSearchResults([])
    } finally {
      setSearching(false)
    }
  }, [dispatch])

  const attack = useCallback(
    async (defenderPlanetId: string, fleetIds: string[]) => {
      setAttacking(true)
      try {
        const result = await attackPlanet({
          defender_planet_id: defenderPlanetId,
          fleet_ids: fleetIds,
        })
        setLastBattle(result)
        return result
      } catch (error: any) {
        const message = error.response?.data?.error || 'Attack failed'
        dispatch({ type: 'SET_ERROR', payload: message })
        return null
      } finally {
        setAttacking(false)
      }
    },
    [dispatch]
  )

  return {
    searchResults,
    searching,
    attacking,
    lastBattle,
    search,
    attack,
    clearResults: () => setSearchResults([]),
    clearBattle: () => setLastBattle(null),
  }
}
