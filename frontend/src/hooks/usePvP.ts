import { useState, useCallback, useEffect, useRef } from 'react'
import { searchPlanets, attackPlanet, getPendingAttacks, cancelAttack, getIncomingAttacks, getPlayerSP } from '../services/api.ts'
import type { PvPSearchResult, AttackPlanetResponse, PendingAttack, RadarIncomingResponse, PlayerSP } from '../types'
import { useGameContext } from '../contexts/GameContext.tsx'

export function usePvP() {
  const { dispatch } = useGameContext()
  const [searchResults, setSearchResults] = useState<PvPSearchResult[]>([])
  const [searching, setSearching] = useState(false)
  const [attacking, setAttacking] = useState(false)
  const [lastDispatch, setLastDispatch] = useState<AttackPlanetResponse | null>(null)
  const [pendingAttacks, setPendingAttacks] = useState<PendingAttack[]>([])
  const [incoming, setIncoming] = useState<RadarIncomingResponse | null>(null)
  const [playerSP, setPlayerSP] = useState<PlayerSP | null>(null)
  const [cooldownSeconds, setCooldownSeconds] = useState(0)
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null)
  const pollRef = useRef<ReturnType<typeof setInterval> | null>(null)

  // Countdown timer for cooldown
  useEffect(() => {
    if (cooldownSeconds <= 0) {
      if (timerRef.current) clearInterval(timerRef.current)
      return
    }
    timerRef.current = setInterval(() => {
      setCooldownSeconds(prev => {
        if (prev <= 1) {
          if (timerRef.current) clearInterval(timerRef.current)
          return 0
        }
        return prev - 1
      })
    }, 1000)
    return () => { if (timerRef.current) clearInterval(timerRef.current) }
  }, [cooldownSeconds > 0]) // eslint-disable-line react-hooks/exhaustive-deps

  // Poll pending attacks + incoming + SP every 10s
  useEffect(() => {
    const fetchAll = async () => {
      try {
        const [pending, inc, sp] = await Promise.all([
          getPendingAttacks(),
          getIncomingAttacks(),
          getPlayerSP(),
        ])
        setPendingAttacks(pending)
        setIncoming(inc)
        setPlayerSP(sp)
      } catch {
        // silent fail on polling
      }
    }

    fetchAll()
    pollRef.current = setInterval(fetchAll, 10000)
    return () => { if (pollRef.current) clearInterval(pollRef.current) }
  }, [])

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
        setLastDispatch(result)
        // Refresh pending attacks
        const pending = await getPendingAttacks()
        setPendingAttacks(pending)
        // Refresh SP
        const sp = await getPlayerSP()
        setPlayerSP(sp)
        return result
      } catch (error: unknown) {
        const axiosErr = error as { response?: { data?: { error?: string; remaining_seconds?: number } } }
        const data = axiosErr.response?.data
        if (data?.remaining_seconds) {
          setCooldownSeconds(data.remaining_seconds)
        }
        const message = data?.error || 'Attack failed'
        dispatch({ type: 'SET_ERROR', payload: message })
        return null
      } finally {
        setAttacking(false)
      }
    },
    [dispatch]
  )

  const cancel = useCallback(async (attackId: string) => {
    try {
      await cancelAttack(attackId)
      const pending = await getPendingAttacks()
      setPendingAttacks(pending)
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Cancel failed'
      dispatch({ type: 'SET_ERROR', payload: message })
    }
  }, [dispatch])

  const refreshPending = useCallback(async () => {
    try {
      const pending = await getPendingAttacks()
      setPendingAttacks(pending)
    } catch {
      // silent
    }
  }, [])

  return {
    searchResults,
    searching,
    attacking,
    lastDispatch,
    pendingAttacks,
    incoming,
    playerSP,
    cooldownSeconds,
    search,
    attack,
    cancel,
    refreshPending,
    clearResults: () => setSearchResults([]),
    clearDispatch: () => setLastDispatch(null),
  }
}
