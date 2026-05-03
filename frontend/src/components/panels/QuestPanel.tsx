import { useState, useEffect, useMemo, type ReactElement, type ReactNode } from 'react'
import { createPortal } from 'react-dom'
import { useQuests } from '../../hooks/useQuests.ts'
import { useBlueprints } from '../../hooks/useBlueprints.ts'
import { useItemTypes } from '../../hooks/useItemTypes.ts'
import { formatNumber } from '../../hooks/useCountdown.ts'
import {
  ScrollText,
  Trophy,
  CircleCheck,
  Lock,
  ChevronRight,
  X,
  Sparkles,
} from 'lucide-react'
import type {
  PlayerQuestWithType,
  DailyQuestEntry,
  DailyTierReward,
  DailyQuestsResponse,
  Blueprint,
} from '../../types'
import type { ItemTypeInfo } from '../../services/api.ts'

type BlueprintLookup = Record<string, Blueprint>
type ItemLookup = Record<string, ItemTypeInfo>

function snakeToTitle(s: string): string {
  return s
    .split('_')
    .filter(Boolean)
    .map(w => w.charAt(0).toUpperCase() + w.slice(1))
    .join(' ')
}

type Tab = 'main' | 'side' | 'daily'

interface QuestPanelProps {
  onClose: () => void
}

export default function QuestPanel({ onClose }: QuestPanelProps) {
  const [tab, setTab] = useState<Tab>('main')
  const [selectedId, setSelectedId] = useState<string | null>(null)
  const [selectedTier, setSelectedTier] = useState<string | null>(null)
  const { quests, daily, loading, error, refreshDaily, claim, claimTier } = useQuests()
  const { allBlueprints } = useBlueprints()
  const { itemTypes } = useItemTypes()
  const [claimingId, setClaimingId] = useState<string | null>(null)
  const [claimingTier, setClaimingTier] = useState<string | null>(null)
  const [rewardFlash, setRewardFlash] = useState<string | null>(null)

  const blueprintLookup = useMemo<BlueprintLookup>(() => {
    const map: BlueprintLookup = {}
    for (const bp of allBlueprints) {
      if (bp.blueprint_key) map[bp.blueprint_key] = bp
    }
    return map
  }, [allBlueprints])

  const itemLookup = useMemo<ItemLookup>(() => {
    const map: ItemLookup = {}
    for (const it of itemTypes) map[it.item_key] = it
    return map
  }, [itemTypes])

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

  const visibleQuests: PlayerQuestWithType[] = useMemo(() => {
    if (!quests) return []
    if (tab === 'main') return quests.main_quests
    if (tab === 'side') return quests.side_quests
    return []
  }, [quests, tab])

  // Auto-select first quest when switching tab or quest list arrives
  useEffect(() => {
    if (tab === 'main' || tab === 'side') {
      if (!selectedId && visibleQuests.length > 0) {
        const preferred = visibleQuests.find(q => q.status === 'completed')
          ?? visibleQuests.find(q => q.status === 'available')
          ?? visibleQuests[0]
        setSelectedId(preferred.id)
      } else if (selectedId && !visibleQuests.find(q => q.id === selectedId)) {
        setSelectedId(visibleQuests[0]?.id ?? null)
      }
    }
  }, [tab, visibleQuests, selectedId])

  const selectedQuest = visibleQuests.find(q => q.id === selectedId) ?? null

  return createPortal(
    <div className="quest-backdrop" onClick={e => { if (e.target === e.currentTarget) onClose() }}>
      <div className="quest-panel ds-modal ds-modal--lg">
        {/* Header */}
        <div className="ds-modal-header">
          <div className="ds-row">
            <ScrollText size={22} strokeWidth={1.75} style={{ color: 'var(--ds-teal)' }} />
            <h2 className="ds-modal-title">Quest Log</h2>
          </div>
          <button className="ds-btn-icon" onClick={onClose} aria-label="Close">
            <X size={16} />
          </button>
        </div>

        {/* Tabs */}
        <div className="ds-tabs quest-tabs">
          <TabButton active={tab === 'main'} count={mainClaimable} onClick={() => { setTab('main'); setSelectedId(null) }}>
            Main
          </TabButton>
          <TabButton active={tab === 'side'} count={sideClaimable} onClick={() => { setTab('side'); setSelectedId(null) }}>
            Side
          </TabButton>
          <TabButton active={tab === 'daily'} count={dailyClaimable} onClick={() => { setTab('daily'); setSelectedTier(null) }}>
            Daily
          </TabButton>
        </div>

        {/* Reward flash */}
        {rewardFlash && (
          <div className="quest-reward-flash">
            <Sparkles size={14} /> <span className="ds-mono">{rewardFlash}</span>
          </div>
        )}

        {/* Content */}
        <div className="quest-body">
          {loading ? (
            <div className="quest-loading">
              <div className="loading-spinner" /><span>Loading quests...</span>
            </div>
          ) : error ? (
            <div className="quest-error">{error}</div>
          ) : tab === 'daily' ? (
            daily ? (
              <DailyQuestTab
                daily={daily}
                onClaimTier={handleClaimTier}
                claimingTier={claimingTier}
                selectedTier={selectedTier}
                onSelectTier={setSelectedTier}
              />
            ) : (
              <div className="quest-loading">
                <div className="loading-spinner" /><span>Loading daily quests...</span>
              </div>
            )
          ) : (
            <QuestTwoColumn
              quests={visibleQuests}
              selectedId={selectedId}
              onSelect={setSelectedId}
              selectedQuest={selectedQuest}
              onClaim={handleClaim}
              claimingId={claimingId}
              blueprintLookup={blueprintLookup}
              itemLookup={itemLookup}
            />
          )}
        </div>
      </div>
    </div>,
    document.body,
  )
}

function TabButton({
  active,
  count,
  onClick,
  children,
}: {
  active: boolean
  count: number
  onClick: () => void
  children: ReactNode
}) {
  return (
    <button className="ds-tab" aria-selected={active} onClick={onClick}>
      {children}
      {count > 0 && (
        <span className="ds-badge ds-badge--orange quest-tab-badge">{count}</span>
      )}
    </button>
  )
}

// ============ Two-column layout (Main/Side) ============

function QuestTwoColumn({
  quests,
  selectedId,
  onSelect,
  selectedQuest,
  onClaim,
  claimingId,
  blueprintLookup,
  itemLookup,
}: {
  quests: PlayerQuestWithType[]
  selectedId: string | null
  onSelect: (id: string) => void
  selectedQuest: PlayerQuestWithType | null
  onClaim: (id: string) => void
  claimingId: string | null
  blueprintLookup: BlueprintLookup
  itemLookup: ItemLookup
}) {
  if (quests.length === 0) {
    return <div className="quest-empty">No quests available.</div>
  }
  return (
    <div className="quest-two-col">
      <div className="quest-list">
        {quests.map(q => (
          <QuestListItem
            key={q.id}
            quest={q}
            selected={q.id === selectedId}
            onClick={() => onSelect(q.id)}
          />
        ))}
      </div>
      <div className="quest-detail">
        {selectedQuest ? (
          <QuestDetail
            quest={selectedQuest}
            onClaim={onClaim}
            claimingId={claimingId}
            blueprintLookup={blueprintLookup}
            itemLookup={itemLookup}
          />
        ) : (
          <div className="quest-empty">Select a quest to see details.</div>
        )}
      </div>
    </div>
  )
}

function QuestListItem({
  quest,
  selected,
  onClick,
}: {
  quest: PlayerQuestWithType
  selected: boolean
  onClick: () => void
}) {
  const status = quest.status
  const isClaimed = status === 'claimed'
  const isCompleted = status === 'completed'
  const isLocked = status === 'locked'

  const Icon = isClaimed ? CircleCheck : isLocked ? Lock : isCompleted ? Trophy : ChevronRight
  const iconColor = isClaimed
    ? 'var(--ds-text-soft)'
    : isCompleted
      ? 'var(--ds-orange)'
      : isLocked
        ? 'var(--ds-text-soft)'
        : 'var(--ds-teal)'

  const badge = isClaimed ? (
    <span className="ds-badge ds-badge--neutral">CLAIMED</span>
  ) : isCompleted ? (
    <span className="ds-badge ds-badge--success">COMPLETE</span>
  ) : isLocked ? (
    <span className="ds-badge ds-badge--neutral">LOCKED</span>
  ) : (
    <span className="ds-badge ds-badge--teal">ACTIVE</span>
  )

  return (
    <button
      className={`ds-list-item quest-list-item ${selected ? 'is-selected' : ''} ${isLocked || isClaimed ? 'is-dim' : ''}`}
      aria-selected={selected}
      onClick={onClick}
    >
      <div className="quest-list-item-icon" style={{ color: iconColor }}>
        <Icon size={18} strokeWidth={1.75} />
      </div>
      <div className="quest-list-item-name">{quest.display_name}</div>
      {badge}
    </button>
  )
}

function QuestDetail({
  quest,
  onClaim,
  claimingId,
  blueprintLookup,
  itemLookup,
}: {
  quest: PlayerQuestWithType
  onClaim: (id: string) => void
  claimingId: string | null
  blueprintLookup: BlueprintLookup
  itemLookup: ItemLookup
}) {
  const isClaimed = quest.status === 'claimed'
  const isCompleted = quest.status === 'completed'
  const isLocked = quest.status === 'locked'
  const total = Math.max(1, quest.requirement_value)
  const current = Math.min(total, quest.progress_value)

  const rewards = parseRewards(quest, blueprintLookup, itemLookup)

  return (
    <div className="ds-panel quest-detail-panel">
      <div className="quest-detail-header">
        <h3 className="ds-h3">{quest.display_name}</h3>
        {isClaimed && <span className="ds-badge ds-badge--neutral">CLAIMED</span>}
        {isCompleted && <span className="ds-badge ds-badge--success">COMPLETE</span>}
        {isLocked && <span className="ds-badge ds-badge--neutral">LOCKED</span>}
        {!isClaimed && !isCompleted && !isLocked && <span className="ds-badge ds-badge--teal">ACTIVE</span>}
      </div>

      {quest.description && (
        <p className="quest-detail-desc">{quest.description}</p>
      )}

      {!isLocked && (
        <div className="quest-detail-progress">
          <div className="ds-row--between" style={{ marginBottom: 'var(--sp-1)' }}>
            <span className="ds-caption">Objective</span>
            <span className="ds-mono">{formatNumber(current)} / {formatNumber(total)}</span>
          </div>
          <SegmentedBar current={current} total={total} />
        </div>
      )}

      <div className="quest-detail-rewards">
        <div className="ds-caption">Rewards</div>
        <div className="quest-rewards-row">{rewards}</div>
      </div>

      {isCompleted && (
        <button
          className="ds-btn-primary ds-btn--block"
          onClick={() => onClaim(quest.id)}
          disabled={claimingId === quest.id}
        >
          {claimingId === quest.id ? 'CLAIMING...' : 'CLAIM REWARD'}
        </button>
      )}
    </div>
  )
}

function SegmentedBar({ current, total }: { current: number; total: number }) {
  // For small totals, show segmented; otherwise simple progress.
  const SEGMENT_LIMIT = 10
  if (total <= SEGMENT_LIMIT) {
    const segments: ReactElement[] = []
    for (let i = 0; i < total; i++) {
      const filled = i < current
      segments.push(
        <span key={i} className={`quest-bar-seg ${filled ? 'is-filled' : ''}`} />,
      )
    }
    return <div className="quest-bar-segmented">{segments}</div>
  }
  const pct = Math.min(100, (current / total) * 100)
  return (
    <div className="ds-bar">
      <div className="ds-bar-fill" style={{ width: `${pct}%` }} />
    </div>
  )
}

// ============ Daily Quest Tab ============

function DailyQuestTab({
  daily,
  onClaimTier,
  claimingTier,
  selectedTier,
  onSelectTier,
}: {
  daily: DailyQuestsResponse
  onClaimTier: (tier: string) => void
  claimingTier: string | null
  selectedTier: string | null
  onSelectTier: (tier: string | null) => void
}) {
  const maxPoints = daily.tier_rewards.length > 0
    ? daily.tier_rewards[daily.tier_rewards.length - 1].points_required
    : 70
  const pct = Math.min(100, (daily.daily_points / maxPoints) * 100)

  return (
    <div className="quest-daily">
      <div className="ds-panel quest-daily-summary">
        <div className="ds-row--between">
          <span className="ds-caption">Daily Points</span>
          <span><span className="ds-mono">{daily.daily_points}</span><span className="ds-text-muted"> / {maxPoints}</span></span>
        </div>
        <div className="ds-bar" style={{ marginTop: 'var(--sp-2)' }}>
          <div className="ds-bar-fill" style={{ width: `${pct}%` }} />
        </div>
        <div className="quest-daily-reset ds-text-soft">
          Resets at 1:00 PM server time
        </div>
      </div>

      <div className="quest-daily-cols">
        <div className="quest-daily-tasks">
          <div className="ds-caption" style={{ marginBottom: 'var(--sp-2)' }}>Today's Quests</div>
          <div className="quest-daily-list">
            {daily.quests.map(q => (
              <DailyQuestRow key={q.quest_key} quest={q} />
            ))}
          </div>
        </div>

        <div className="quest-tier-list">
          <div className="ds-caption" style={{ marginBottom: 'var(--sp-2)' }}>Reward Tiers</div>
          {daily.tier_rewards.map(tier => (
            <TierRewardRow
              key={tier.tier}
              tier={tier}
              currentPoints={daily.daily_points}
              onClaim={onClaimTier}
              claiming={claimingTier === tier.tier}
              selected={selectedTier === tier.tier}
              onSelect={() => onSelectTier(tier.tier === selectedTier ? null : tier.tier)}
            />
          ))}
        </div>
      </div>
    </div>
  )
}

function DailyQuestRow({ quest }: { quest: DailyQuestEntry }) {
  const isMultiStep = quest.progress !== undefined && quest.required !== undefined

  return (
    <div className={`quest-daily-row ${quest.completed ? 'is-completed' : ''}`}>
      <span className="quest-daily-check">
        {quest.completed ? <CircleCheck size={16} style={{ color: 'var(--ds-success)' }} /> : <span className="quest-daily-circle" />}
      </span>
      <span className="quest-daily-name">
        {quest.display_name}
        {isMultiStep && <span className="ds-text-muted ds-mono"> ({quest.progress}/{quest.required})</span>}
      </span>
      <span className="quest-daily-pts ds-mono">
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
  selected,
  onSelect,
}: {
  tier: DailyTierReward
  currentPoints: number
  onClaim: (tier: string) => void
  claiming: boolean
  selected: boolean
  onSelect: () => void
}) {
  const reachable = currentPoints >= tier.points_required
  const claimable = reachable && !tier.claimed

  const TIER_REWARDS: Record<string, string> = {
    bronze: 'Loudspeaker',
    silver: 'Resource Box',
    gold: 'SP Card',
    diamond: 'Raw Gemstone',
  }

  const stateBadge = tier.claimed ? (
    <span className="ds-badge ds-badge--neutral">CLAIMED</span>
  ) : claimable ? (
    <span className="ds-badge ds-badge--success">READY</span>
  ) : (
    <span className="ds-badge ds-badge--neutral">{tier.points_required - currentPoints} pts</span>
  )

  return (
    <div
      className={`ds-list-item quest-tier-row ${selected ? 'is-selected' : ''} ${!reachable && !tier.claimed ? 'is-dim' : ''}`}
      onClick={onSelect}
      aria-selected={selected}
      role="button"
      tabIndex={0}
      onKeyDown={e => { if (e.key === 'Enter' || e.key === ' ') onSelect() }}
    >
      <span className={`quest-tier-name tier-${tier.tier}`}>
        {tier.tier.charAt(0).toUpperCase() + tier.tier.slice(1)}
      </span>
      <span className="quest-tier-pts ds-mono ds-text-muted">{tier.points_required} pts</span>
      <span className="quest-tier-reward">{TIER_REWARDS[tier.tier] || ''}</span>
      <span className="quest-tier-status">
        {claimable ? (
          <button
            className="ds-btn-primary ds-btn--sm"
            onClick={e => { e.stopPropagation(); onClaim(tier.tier) }}
            disabled={claiming}
          >
            {claiming ? '...' : 'CLAIM'}
          </button>
        ) : stateBadge}
      </span>
    </div>
  )
}

// ============ Reward parsing ============

type RewardItem = {
  type?: string
  blueprint_key?: string
  item_key?: string
  quantity?: number
}

function parseRewards(
  quest: PlayerQuestWithType,
  blueprintLookup: BlueprintLookup,
  itemLookup: ItemLookup,
) {
  const parts: ReactElement[] = []
  if (quest.reward_metal > 0) {
    parts.push(
      <span key="metal" className="quest-reward-chip">
        <span className="ds-resource-dot ds-resource-dot--metal" />
        <span className="ds-mono">+{formatNumber(quest.reward_metal)}</span>
        <span className="ds-text-muted">Metal</span>
      </span>,
    )
  }
  if (quest.reward_he3 > 0) {
    parts.push(
      <span key="he3" className="quest-reward-chip">
        <span className="ds-resource-dot ds-resource-dot--he3" />
        <span className="ds-mono">+{formatNumber(quest.reward_he3)}</span>
        <span className="ds-text-muted">He3</span>
      </span>,
    )
  }
  if (quest.reward_gold > 0) {
    parts.push(
      <span key="gold" className="quest-reward-chip">
        <span className="ds-resource-dot ds-resource-dot--gold" />
        <span className="ds-mono">+{formatNumber(quest.reward_gold)}</span>
        <span className="ds-text-muted">Gold</span>
      </span>,
    )
  }
  if (quest.reward_item_json) {
    try {
      const raw = quest.reward_item_json
      const items: RewardItem[] = Array.isArray(raw)
        ? (raw as RewardItem[])
        : typeof raw === 'string'
          ? (JSON.parse(raw) as RewardItem[])
          : []
      if (Array.isArray(items)) {
        items.forEach((item, i) => {
          if (item.type === 'blueprint' && item.blueprint_key) {
            const bp = blueprintLookup[item.blueprint_key]
            const label = bp?.display_name || snakeToTitle(item.blueprint_key)
            parts.push(
              <span key={`item-${i}`} className="quest-reward-chip">
                <span className="ds-badge ds-badge--info">BP</span>
                <span>{label}</span>
              </span>,
            )
          } else if (item.type === 'item' && item.item_key) {
            const qty = item.quantity ?? 1
            const meta = itemLookup[item.item_key]
            const label = meta?.display_name || snakeToTitle(item.item_key)
            parts.push(
              <span key={`item-${i}`} className="quest-reward-chip">
                {qty > 1 && <span className="ds-mono">{qty}x</span>}
                <span>{label}</span>
              </span>,
            )
          } else {
            const qty = item.quantity ?? 1
            const fallbackKey = item.blueprint_key || item.item_key || item.type || 'item'
            parts.push(
              <span key={`item-${i}`} className="quest-reward-chip">
                {qty > 1 && <span className="ds-mono">{qty}x</span>}
                <span>{snakeToTitle(fallbackKey)}</span>
              </span>,
            )
          }
        })
      }
    } catch {
      // ignore parse errors
    }
  }
  if (parts.length === 0) {
    return <span className="ds-text-soft">No rewards</span>
  }
  return <>{parts}</>
}
