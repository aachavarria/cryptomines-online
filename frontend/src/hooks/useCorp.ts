import { useState, useEffect, useCallback } from 'react'
import {
  getCorp,
  createCorp,
  joinCorp,
  leaveCorp,
  listCorpMembers,
  donateResources,
  updateMemberRole,
  searchCorps,
} from '../services/api'
import type {
  CorpWithBonuses,
  CorpMember,
  CreateCorpRequest,
  DonateRequest,
  CorpSearchResult,
} from '../types/corps'
import { useGameContext } from '../contexts/GameContext'

export function useCorp() {
  const { dispatch } = useGameContext()
  const [corpData, setCorpData] = useState<CorpWithBonuses | null>(null)
  const [members, setMembers] = useState<CorpMember[]>([])
  const [searchResults, setSearchResults] = useState<CorpSearchResult[]>([])
  const [loading, setLoading] = useState(false)
  const [membersLoading, setMembersLoading] = useState(false)
  const [searchLoading, setSearchLoading] = useState(false)

  const fetchCorp = useCallback(async () => {
    setLoading(true)
    try {
      const data = await getCorp()
      setCorpData(data)
    } catch (error: unknown) {
      const axiosErr = error as { response?: { data?: { error?: string } } }
      const message = axiosErr.response?.data?.error || 'Failed to fetch corp data'
      dispatch({ type: 'SET_ERROR', payload: message })
    } finally {
      setLoading(false)
    }
  }, [dispatch])

  const fetchMembers = useCallback(async () => {
    setMembersLoading(true)
    try {
      const data = await listCorpMembers()
      setMembers(data)
    } catch (error: unknown) {
      const axiosErr = error as { response?: { data?: { error?: string } } }
      const message = axiosErr.response?.data?.error || 'Failed to fetch corp members'
      dispatch({ type: 'SET_ERROR', payload: message })
    } finally {
      setMembersLoading(false)
    }
  }, [dispatch])

  const create = useCallback(
    async (req: CreateCorpRequest) => {
      setLoading(true)
      try {
        await createCorp(req)
        await fetchCorp()
        return true
      } catch (error: unknown) {
        const axiosErr = error as { response?: { data?: { error?: string } } }
        const message = axiosErr.response?.data?.error || 'Failed to create corp'
        dispatch({ type: 'SET_ERROR', payload: message })
        return false
      } finally {
        setLoading(false)
      }
    },
    [dispatch, fetchCorp],
  )

  const join = useCallback(
    async (corpId: string) => {
      setLoading(true)
      try {
        await joinCorp(corpId)
        await fetchCorp()
        return true
      } catch (error: unknown) {
        const axiosErr = error as { response?: { data?: { error?: string } } }
        const message = axiosErr.response?.data?.error || 'Failed to join corp'
        dispatch({ type: 'SET_ERROR', payload: message })
        return false
      } finally {
        setLoading(false)
      }
    },
    [dispatch, fetchCorp],
  )

  const leave = useCallback(async () => {
    setLoading(true)
    try {
      await leaveCorp()
      setCorpData(null)
      setMembers([])
      return true
    } catch (error: unknown) {
      const axiosErr = error as { response?: { data?: { error?: string } } }
      const message = axiosErr.response?.data?.error || 'Failed to leave corp'
      dispatch({ type: 'SET_ERROR', payload: message })
      return false
    } finally {
      setLoading(false)
    }
  }, [dispatch])

  const donate = useCallback(
    async (req: DonateRequest) => {
      setLoading(true)
      try {
        await donateResources(req)
        await fetchCorp()
        await fetchMembers()
        return true
      } catch (error: unknown) {
        const axiosErr = error as { response?: { data?: { error?: string } } }
        const message = axiosErr.response?.data?.error || 'Failed to donate resources'
        dispatch({ type: 'SET_ERROR', payload: message })
        return false
      } finally {
        setLoading(false)
      }
    },
    [dispatch, fetchCorp, fetchMembers],
  )

  const updateRole = useCallback(
    async (memberId: string, role: string) => {
      setLoading(true)
      try {
        await updateMemberRole(memberId, role)
        await fetchMembers()
        return true
      } catch (error: unknown) {
        const axiosErr = error as { response?: { data?: { error?: string } } }
        const message = axiosErr.response?.data?.error || 'Failed to update member role'
        dispatch({ type: 'SET_ERROR', payload: message })
        return false
      } finally {
        setLoading(false)
      }
    },
    [dispatch, fetchMembers],
  )

  const search = useCallback(
    async (query: string) => {
      if (!query.trim()) {
        setSearchResults([])
        return
      }
      setSearchLoading(true)
      try {
        const data = await searchCorps(query)
        setSearchResults(data)
      } catch (error: unknown) {
        const axiosErr = error as { response?: { data?: { error?: string } } }
        const message = axiosErr.response?.data?.error || 'Failed to search corps'
        dispatch({ type: 'SET_ERROR', payload: message })
        setSearchResults([])
      } finally {
        setSearchLoading(false)
      }
    },
    [dispatch],
  )

  const refresh = useCallback(async () => {
    await fetchCorp()
    if (corpData?.corp) {
      await fetchMembers()
    }
  }, [fetchCorp, fetchMembers, corpData])

  // Auto-refresh corp data every 30s
  useEffect(() => {
    fetchCorp()
    const interval = setInterval(() => {
      fetchCorp()
    }, 30000)
    return () => clearInterval(interval)
  }, [fetchCorp])

  return {
    corpData,
    members,
    searchResults,
    loading,
    membersLoading,
    searchLoading,
    fetchCorp,
    fetchMembers,
    create,
    join,
    leave,
    donate,
    updateRole,
    search,
    refresh,
  }
}
