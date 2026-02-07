import { useState, useEffect } from 'react'
import { createPortal } from 'react-dom'
import { useQuests } from '../../hooks/useQuests.ts'
import { formatNumber } from '../../hooks/useCountdown.ts'
import type {
  PlayerQuestWithType,
  DailyQuestEntry,
  DailyTierReward,
} from '../../types'

type Tab = 'main' | 'side' | 'daily'

interface QuestPanelProps {
  onClose: () => void
}

export default function QuestPanel({ onClose }: QuestPanelProps) {
  const [tab, setTab] = useState<Tab>('main')
  const { quests, daily, loading, error, refreshDaily, claim, claimTier } = useQuests()
  const [claimingId, setClaimingId] = useState<string | null>(null)
  const [claimingTier, setClaimingTier] = useState<string | null>(null)
  const [rewardFlash, setRewardFlash] = useState<string | null>(null)

  // Load daily quests when switching to daily tab
  useEffect(() => {
    if (tab === 'daily' && !daily) {
      refreshDaily()
    }
  }, [tab, daily, refreshDaily])

  // ESC to close
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose])

  async function handleClaim(questId: string) {
    setClaimingId(questId)
    try {
      const result = await claim(questId)
      setRewardFlash(
        `+${formatNumber(result.rewards.metal)}M  +${formatNumber(result.rewards.he3)}H  +${formatNumber(result.rewards.gold)}G`
      )
      setTimeout(() => setRewardFlash(null), 3000)
    } catch {
      // error handled by hook
    } finally {
      setClaimingId(null)
    }
  }

  async function handleClaimTier(tier: string) {
    setClaimingTier(tier)
    try {
      const result = await claimTier(tier)
      setRewardFlash(`Tier reward: ${result.reward.quantity}x ${result.reward.type}`)
      setTimeout(() => setRewardFlash(null), 3000)
    } catch {
      // error handled by hook
    } finally {
      setClaimingTier(null)
    }
  }

  // Count claimable per tab for tab badges
  const mainClaimable = quests?.main_quests.filter(q => q.status === 'completed').length || 0
  const sideClaimable = quests?.side_quests.filter(q => q.status === 'completed').length || 0
  const dailyClaimable = daily?.tier_rewards.filter(
    t => !t.claimed && daily.daily_points >= t.points_required,
  ).length || 0

  return createPortal(
    <div className="quest-backdrop" onClick={e => { if (e.target === e.currentTarget) onClose() }}>
      <div className="quest-panel">
        {/* Header */}
        <div className="quest-header">
          <span className="quest-title">QUEST LOG</span>
          <button className="p2-modal-close" onClick={onClose}>X</button>
        </div>

        {/* Tabs */}
        <div className="quest-tabs">
          <button
            className={`quest-tab ${tab === 'main' ? 'active' : ''}`}
            onClick={() => setTab('main')}
          >
            Main Quests
            {mainClaimable > 0 && <span className="quest-tab-badge">{mainClaimable}</span>}
          </button>
          <button
            className={`quest-tab ${tab === 'side' ? 'active' : ''}`}
            onClick={() => setTab('side')}
          >
            Side Quests
            {sideClaimable > 0 && <span className="quest-tab-badge">{sideClaimable}</span>}
          </button>
          <button
            className={`quest-tab ${tab === 'daily' ? 'active' : ''}`}
            onClick={() => setTab('daily')}
          >
            Daily Quests
            {dailyClaimable > 0 && <span className="quest-tab-badge">{dailyClaimable}</span>}
          </button>
        </div>

        {/* Reward flash */}
        {rewardFlash && (
          <div className="quest-reward-flash">{rewardFlash}</div>
        )}

        {/* Content */}
        <div className="quest-body">
          {loading ? (
            <div className="p2-panel-loading">
              <div className="loading-spinner" /><span>Loading quests...</span>
            </div>
          ) : error ? (
            <div className="p2-panel-error">{error}</div>
          ) : (
            <>
              {tab === 'main' && quests && (
                <MainQuestTab
                  quests={quests.main_quests}
                  currentQuest={quests.current_main_quest}
                  onClaim={handleClaim}
                  claimingId={claimingId}
                />
              )}
              {tab === 'side' && quests && (
                <SideQuestTab
                  quests={quests.side_quests}
                  onClaim={handleClaim}
                  claimingId={claimingId}
                />
              )}
              {tab === 'daily' && daily && (
                <DailyQuestTab
                  daily={daily}
                  onClaimTier={handleClaimTier}
                  claimingTier={claimingTier}
                />
              )}
              {tab === 'daily' && !daily && (
                <div className="p2-panel-loading">
                  <div className="loading-spinner" /><span>Loading daily quests...</span>
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </div>,
    document.body,
  )
}

// ============ Main Quest Tab ============

function MainQuestTab({
  quests,
  currentQuest,
  onClaim,
  claimingId,
}: {
  quests: PlayerQuestWithType[]
  currentQuest: PlayerQuestWithType | null
  onClaim: (id: string) => void
  claimingId: string | null
}) {
  const claimed = quests.filter(q => q.status === 'claimed')
  const locked = quests.filter(q => q.status === 'locked')

  return (
    <div className="quest-content">
      {/* Current Quest */}
      {currentQuest && (
        <div className="quest-section">
          <div className="quest-section-title">CURRENT MAIN QUEST</div>
          <QuestCard quest={currentQuest} onClaim={onClaim} claimingId={claimingId} highlight />
        </div>
      )}

      {/* Completed but unclaimed quests */}
      {quests.filter(q => q.status === 'completed' && q.id !== currentQuest?.id).map(q => (
        <div key={q.id} className="quest-section">
          <QuestCard quest={q} onClaim={onClaim} claimingId={claimingId} highlight />
        </div>
      ))}

      {/* Claimed quests */}
      {claimed.length > 0 && (
        <div className="quest-section">
          <div className="quest-section-title">COMPLETED QUESTS</div>
          {claimed.map(q => (
            <QuestCard key={q.id} quest={q} onClaim={onClaim} claimingId={claimingId} />
          ))}
        </div>
      )}

      {/* Locked quests */}
      {locked.length > 0 && (
        <div className="quest-section">
          <div className="quest-section-title">UPCOMING QUESTS</div>
          {locked.map(q => (
            <QuestCard key={q.id} quest={q} onClaim={onClaim} claimingId={claimingId} />
          ))}
        </div>
      )}

      {quests.length === 0 && (
        <div className="p2-empty-state">No main quests available.</div>
      )}
    </div>
  )
}

// ============ Side Quest Tab ============

function SideQuestTab({
  quests,
  onClaim,
  claimingId,
}: {
  quests: PlayerQuestWithType[]
  onClaim: (id: string) => void
  claimingId: string | null
}) {
  // Group side quests by category-like prefix (derive from quest_key)
  const groups = groupSideQuests(quests)

  return (
    <div className="quest-content">
      {groups.length === 0 && (
        <div className="p2-empty-state">No side quests available.</div>
      )}
      {groups.map(group => (
        <div key={group.label} className="quest-section">
          <div className="quest-section-title">{group.label}</div>
          {group.quests.map(q => (
            <QuestCard key={q.id} quest={q} onClaim={onClaim} claimingId={claimingId} />
          ))}
        </div>
      ))}
    </div>
  )
}

function groupSideQuests(quests: PlayerQuestWithType[]): { label: string; quests: PlayerQuestWithType[] }[] {
  const groupMap = new Map<string, PlayerQuestWithType[]>()
  for (const q of quests) {
    // Derive a group label from the quest display_name (strip trailing tier numbers/roman numerals)
    const label = deriveGroupLabel(q.display_name)
    const existing = groupMap.get(label) || []
    existing.push(q)
    groupMap.set(label, existing)
  }
  return Array.from(groupMap.entries()).map(([label, quests]) => ({ label, quests }))
}

function deriveGroupLabel(name: string): string {
  // Strip trailing roman numerals or numbers (e.g. "Harvest Time II" -> "Harvest Time")
  const stripped = name.replace(/\s+(I{1,3}V?|IV|VI{0,3}|IX|X{0,3}|\d+)\s*$/, '')
  return stripped || name
}

// ============ Daily Quest Tab ============

function DailyQuestTab({
  daily,
  onClaimTier,
  claimingTier,
}: {
  daily: DailyQuestsResponse
  onClaimTier: (tier: string) => void
  claimingTier: string | null
}) {
  const maxPoints = daily.tier_rewards.length > 0
    ? daily.tier_rewards[daily.tier_rewards.length - 1].points_required
    : 70

  return (
    <div className="quest-content">
      {/* Today's quests */}
      <div className="quest-section">
        <div className="quest-section-title-row">
          <span className="quest-section-title">TODAY'S QUESTS</span>
          <span className="quest-daily-points">
            Points: <strong>{daily.daily_points}</strong>/{maxPoints}
          </span>
        </div>
        <div className="quest-daily-list">
          {daily.quests.map(q => (
            <DailyQuestRow key={q.quest_key} quest={q} />
          ))}
        </div>
      </div>

      {/* Tier rewards */}
      <div className="quest-section">
        <div className="quest-section-title">REWARD TIERS</div>
        <div className="quest-tier-list">
          {daily.tier_rewards.map(tier => (
            <TierRewardRow
              key={tier.tier}
              tier={tier}
              currentPoints={daily.daily_points}
              onClaim={onClaimTier}
              claiming={claimingTier === tier.tier}
            />
          ))}
        </div>
      </div>

      <div className="quest-daily-reset">
        Resets at 1:00 PM server time
      </div>
    </div>
  )
}

function DailyQuestRow({ quest }: { quest: DailyQuestEntry }) {
  const isMultiStep = quest.progress !== undefined && quest.required !== undefined

  return (
    <div className={`quest-daily-row ${quest.completed ? 'completed' : ''}`}>
      <span className="quest-daily-check">
        {quest.completed ? '\u2713' : '\u25CB'}
      </span>
      <span className="quest-daily-name">
        {quest.display_name}
        {isMultiStep && ` (${quest.progress}/${quest.required})`}
      </span>
      <span className="quest-daily-pts">
        +{quest.points} pt{quest.points !== 1 ? 's' : ''}
        {isMultiStep && quest.points_per && ' each'}
      </span>
    </div>
  )
}

function TierRewardRow({
  tier,
  currentPoints,
  onClaim,
  claiming,
}: {
  tier: DailyTierReward
  currentPoints: number
  onClaim: (tier: string) => void
  claiming: boolean
}) {
  const reachable = currentPoints >= tier.points_required
  const claimable = reachable && !tier.claimed

  const TIER_REWARDS: Record<string, string> = {
    bronze: 'Loudspeaker',
    silver: 'Resource Box',
    gold: 'SP Card',
    diamond: 'Raw Gemstone',
  }

  return (
    <div className={`quest-tier-row ${tier.claimed ? 'claimed' : ''} ${claimable ? 'claimable' : ''} ${!reachable ? 'locked' : ''}`}>
      <span className="quest-tier-pts">{tier.points_required} pts</span>
      <span className={`quest-tier-name tier-${tier.tier}`}>
        {tier.tier.charAt(0).toUpperCase() + tier.tier.slice(1)}
      </span>
      <span className="quest-tier-reward">{TIER_REWARDS[tier.tier] || ''}</span>
      <span className="quest-tier-status">
        {tier.claimed ? (
          <span className="quest-badge claimed">CLAIMED</span>
        ) : claimable ? (
          <button
            className="quest-claim-btn"
            onClick={() => onClaim(tier.tier)}
            disabled={claiming}
          >
            {claiming ? '...' : 'CLAIM'}
          </button>
        ) : (
          <QuestProgressMini current={currentPoints} required={tier.points_required} />
        )}
      </span>
    </div>
  )
}

// ============ Shared Components ============

function QuestCard({
  quest,
  onClaim,
  claimingId,
  highlight,
}: {
  quest: PlayerQuestWithType
  onClaim: (id: string) => void
  claimingId: string | null
  highlight?: boolean
}) {
  const isClaimed = quest.status === 'claimed'
  const isCompleted = quest.status === 'completed'
  const isLocked = quest.status === 'locked'
  const isAvailable = quest.status === 'available'
  const progressPct = quest.requirement_value > 0
    ? Math.min(100, (quest.progress_value / quest.requirement_value) * 100)
    : 0

  const rewards = parseRewards(quest)

  return (
    <div className={`quest-card ${isClaimed ? 'dimmed' : ''} ${isCompleted ? 'complete' : ''} ${isLocked ? 'locked' : ''} ${highlight ? 'highlight' : ''}`}>
      <div className="quest-card-header">
        <span className="quest-card-icon">
          {isClaimed && '\u2713'}
          {isCompleted && '!'}
          {isAvailable && '\u25B6'}
          {isLocked && '\u{1F512}'}
        </span>
        <span className="quest-card-name">{quest.display_name}</span>
        <span className="quest-card-badge-area">
          {isClaimed && <span className="quest-badge claimed">CLAIMED</span>}
          {isLocked && <span className="quest-badge locked">LOCKED</span>}
          {isCompleted && !isClaimed && <span className="quest-badge complete">COMPLETE</span>}
        </span>
      </div>

      {!isClaimed && !isLocked && (
        <>
          <div className="quest-card-desc">{quest.description}</div>
          {(isAvailable || isCompleted) && (
            <div className="quest-card-progress">
              <div className="quest-progress-bar">
                <div
                  className="quest-progress-fill"
                  style={{ width: `${progressPct}%` }}
                />
              </div>
              <span className="quest-progress-text">
                {formatNumber(quest.progress_value)}/{formatNumber(quest.requirement_value)}
              </span>
            </div>
          )}
          <div className="quest-card-rewards">
            <span className="quest-reward-label">Rewards:</span>
            {rewards}
          </div>
        </>
      )}

      {isCompleted && (
        <button
          className="quest-claim-btn"
          onClick={() => onClaim(quest.id)}
          disabled={claimingId === quest.id}
        >
          {claimingId === quest.id ? 'Claiming...' : 'CLAIM REWARD'}
        </button>
      )}
    </div>
  )
}

function parseRewards(quest: PlayerQuestWithType) {
  const parts: JSX.Element[] = []
  if (quest.reward_metal > 0) {
    parts.push(
      <span key="metal" className="quest-res">
        <span className="quest-res-dot metal" />
        {formatNumber(quest.reward_metal)}
      </span>
    )
  }
  if (quest.reward_he3 > 0) {
    parts.push(
      <span key="he3" className="quest-res">
        <span className="quest-res-dot he3" />
        {formatNumber(quest.reward_he3)}
      </span>
    )
  }
  if (quest.reward_gold > 0) {
    parts.push(
      <span key="gold" className="quest-res">
        <span className="quest-res-dot gold" />
        {formatNumber(quest.reward_gold)}
      </span>
    )
  }
  if (quest.reward_item_json) {
    try {
      const items = JSON.parse(quest.reward_item_json)
      if (Array.isArray(items)) {
        items.forEach((item: { type: string; quantity: number }, i: number) => {
          parts.push(
            <span key={`item-${i}`} className="quest-res item">
              {item.quantity}x {item.type}
            </span>
          )
        })
      }
    } catch {
      // ignore parse errors
    }
  }
  return <span className="quest-rewards-inline">{parts}</span>
}

function QuestProgressMini({ current, required }: { current: number; required: number }) {
  const pct = Math.min(100, (current / required) * 100)
  return (
    <div className="quest-progress-mini">
      <div className="quest-progress-mini-bar">
        <div className="quest-progress-mini-fill" style={{ width: `${pct}%` }} />
      </div>
      <span className="quest-progress-mini-text">{current}/{required}</span>
    </div>
  )
}
