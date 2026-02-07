import { useState, useEffect, useCallback } from 'react'
import {
  getQuests,
  getDailyQuests,
  claimQuest,
  claimDailyTier,
} from '../services/api.ts'
import type {
  PlayerQuestWithType,
  QuestsResponse,
  DailyQuestsResponse,
  ClaimQuestResponse,
  ClaimDailyTierResponse,
} from '../types'

export function useQuests() {
  const [quests, setQuests] = useState<QuestsResponse | null>(null)
  const [daily, setDaily] = useState<DailyQuestsResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const refreshQuests = useCallback(async () => {
    try {
      const data = await getQuests()
      setQuests(data)
      setError(null)
    } catch {
      setError('Failed to load quests')
    } finally {
      setLoading(false)
    }
  }, [])

  const refreshDaily = useCallback(async () => {
    try {
      const data = await getDailyQuests()
      setDaily(data)
      setError(null)
    } catch {
      setError('Failed to load daily quests')
    }
  }, [])

  useEffect(() => {
    refreshQuests()
  }, [refreshQuests])

  // Auto-refresh quests every 5 seconds to show progress updates
  useEffect(() => {
    const interval = setInterval(() => {
      refreshQuests()
    }, 5000)
    return () => clearInterval(interval)
  }, [refreshQuests])

  const claim = useCallback(async (questId: string): Promise<ClaimQuestResponse> => {
    const result = await claimQuest(questId)
    await refreshQuests()
    return result
  }, [refreshQuests])

  const claimTier = useCallback(async (tier: string): Promise<ClaimDailyTierResponse> => {
    const result = await claimDailyTier(tier)
    await refreshDaily()
    return result
  }, [refreshDaily])

  // Count claimable quests for notification badge
  const claimableCount = getClaimableCount(quests, daily)

  return {
    quests,
    daily,
    loading,
    error,
    refreshQuests,
    refreshDaily,
    claim,
    claimTier,
    claimableCount,
  }
}

function getClaimableCount(
  quests: QuestsResponse | null,
  daily: DailyQuestsResponse | null,
): number {
  let count = 0
  if (quests) {
    count += quests.main_quests.filter(q => q.status === 'completed').length
    count += quests.side_quests.filter(q => q.status === 'completed').length
  }
  if (daily) {
    count += daily.tier_rewards.filter(
      t => !t.claimed && daily.daily_points >= t.points_required,
    ).length
  }
  return count
}
