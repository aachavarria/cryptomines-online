import { useState, useEffect, useMemo } from 'react'
import { createPortal } from 'react-dom'
import { useResearch } from '../../hooks/useResearch.ts'
import { useCountdown, formatDuration, formatNumber } from '../../hooks/useCountdown.ts'
import type { TechTree, TechWithProgress } from '../../types'

const TREES: { key: TechTree; label: string }[] = [
  { key: 'logistics_construction', label: 'Logistics' },
  { key: 'ballistics_science', label: 'Ballistics' },
  { key: 'directional_science', label: 'Directional' },
  { key: 'missile_science', label: 'Missile' },
  { key: 'ship_based_science', label: 'Ship-Based' },
  { key: 'ship_defense_science', label: 'Ship Defense' },
  { key: 'planetary_defense', label: 'Planetary' },
]

interface ResearchPanelProps {
  onClose: () => void
}

export default function ResearchPanel({ onClose }: ResearchPanelProps) {
  const [tab, setTab] = useState<TechTree>('logistics_construction')
  const [confirmTech, setConfirmTech] = useState<TechWithProgress | null>(null)
  const [actionPending, setActionPending] = useState(false)
  const [toast, setToast] = useState<string | null>(null)

  const { trees, active, loading, error, start, cancel, speedup, refreshAll } = useResearch()

  // ESC to close
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === 'Escape') {
        if (confirmTech) {
          setConfirmTech(null)
        } else {
          onClose()
        }
      }
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose, confirmTech])

  // Poll for active research completion
  useEffect(() => {
    if (!active) return
    const interval = setInterval(() => {
      const finishTime = new Date(active.research_finish_at).getTime()
      if (Date.now() >= finishTime) {
        refreshAll()
      }
    }, 2000)
    return () => clearInterval(interval)
  }, [active, refreshAll])

  const currentTreeTechs = trees[tab] || []

  // Group techs by tiers based on prerequisite depth
  const tiers = useMemo(() => buildTiers(currentTreeTechs), [currentTreeTechs])

  // Check which tree has active research
  const activeTreeKey = active?.tree || null

  async function handleStartResearch(tech: TechWithProgress) {
    setConfirmTech(tech)
  }

  async function handleConfirmResearch() {
    if (!confirmTech) return
    setActionPending(true)
    try {
      await start(confirmTech.id)
      setToast(`Research started: ${confirmTech.display_name}`)
      setTimeout(() => setToast(null), 3000)
      setConfirmTech(null)
    } catch {
      setToast('Failed to start research')
      setTimeout(() => setToast(null), 3000)
    } finally {
      setActionPending(false)
    }
  }

  async function handleCancel() {
    if (!active) return
    setActionPending(true)
    try {
      await cancel(active.tech_type_id)
      setToast('Research cancelled')
      setTimeout(() => setToast(null), 3000)
    } catch {
      setToast('Failed to cancel research')
      setTimeout(() => setToast(null), 3000)
    } finally {
      setActionPending(false)
    }
  }

  async function handleSpeedup() {
    if (!active) return
    setActionPending(true)
    try {
      await speedup(active.tech_type_id, 30)
      setToast('Research accelerated by 30 minutes')
      setTimeout(() => setToast(null), 3000)
    } catch {
      setToast('Failed to speed up research')
      setTimeout(() => setToast(null), 3000)
    } finally {
      setActionPending(false)
    }
  }

  return createPortal(
    <div className="research-backdrop" onClick={e => { if (e.target === e.currentTarget) onClose() }}>
      <div className="research-panel">
        {/* Header */}
        <div className="research-header">
          <span className="research-title">
            TECHNOLOGY CENTER
          </span>
          <button className="p2-modal-close" onClick={onClose}>X</button>
        </div>

        {/* Tabs */}
        <div className="research-tabs">
          {TREES.map(t => (
            <button
              key={t.key}
              className={`research-tab ${tab === t.key ? 'active' : ''}`}
              onClick={() => setTab(t.key)}
            >
              {t.label}
              {activeTreeKey === t.key && <span className="research-tab-dot" />}
            </button>
          ))}
        </div>

        {/* Toast */}
        {toast && (
          <div className="quest-reward-flash">{toast}</div>
        )}

        {/* Body */}
        <div className="research-body">
          {loading ? (
            <div className="p2-panel-loading">
              <div className="loading-spinner" /><span>Loading research...</span>
            </div>
          ) : error ? (
            <div className="p2-panel-error">{error}</div>
          ) : currentTreeTechs.length === 0 ? (
            <div className="research-empty">No technologies available in this tree.</div>
          ) : (
            <div className="research-tree-grid">
              {tiers.map((tierTechs, tierIndex) => (
                <div key={tierIndex} className="research-tree-tier">
                  {tierTechs.map(tech => (
                    <TechCard
                      key={tech.id}
                      tech={tech}
                      active={active}
                      onResearch={handleStartResearch}
                    />
                  ))}
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Active Research Bar */}
        {active && active.tree === tab ? (
          <ActiveResearchBar
            active={active}
            onSpeedup={handleSpeedup}
            onCancel={handleCancel}
            pending={actionPending}
          />
        ) : (
          <div className="research-no-active">
            {active ? `Active research in ${getTreeLabel(active.tree)}` : 'No active research'}
          </div>
        )}
      </div>

      {/* Confirm Modal */}
      {confirmTech && (
        <ResearchConfirmModal
          tech={confirmTech}
          onConfirm={handleConfirmResearch}
          onCancel={() => setConfirmTech(null)}
          pending={actionPending}
          hasActiveResearch={active !== null}
        />
      )}
    </div>,
    document.body,
  )
}

// ============ Tech Card ============

function TechCard({
  tech,
  active,
  onResearch,
}: {
  tech: TechWithProgress
  active: { tech_type_id: number } | null
  onResearch: (tech: TechWithProgress) => void
}) {
  const isMaxed = tech.current_level >= tech.max_level
  const isResearching = tech.is_researching
  const isLocked = getCardState(tech) === 'locked'
  const isAvailable = !isMaxed && !isResearching && !isLocked
  const isCompleted = tech.current_level > 0 && !isMaxed

  const stateClass = isResearching
    ? 'tech-card--researching'
    : isMaxed
      ? 'tech-card--maxed'
      : isLocked
        ? 'tech-card--locked'
        : isCompleted
          ? 'tech-card--completed'
          : isAvailable
            ? 'tech-card--available'
            : ''

  const effectText = formatEffect(tech)
  const nextEffectText = !isMaxed ? formatNextEffect(tech) : null

  const hasActiveElsewhere = active !== null && active.tech_type_id !== tech.id

  return (
    <div className={`tech-card ${stateClass}`}>
      <div className="tech-card-header">
        <span className="tech-card-name">{tech.display_name}</span>
        {isMaxed ? (
          <span className="tech-card-badge max-badge">MAX</span>
        ) : isResearching ? (
          <span className="tech-card-badge researching-badge">...</span>
        ) : (
          <span className="tech-card-level">
            Lv {tech.current_level}/{tech.max_level}
          </span>
        )}
      </div>

      {effectText && (
        <div className="tech-card-effect">{effectText}</div>
      )}

      {isResearching && tech.research_finish_at && (
        <TechCardResearchProgress finishAt={tech.research_finish_at} />
      )}

      {!isMaxed && !isResearching && !isLocked && (
        <>
          {nextEffectText && (
            <div className="tech-card-next">Next: {nextEffectText}</div>
          )}
          {tech.cost_next_level && (
            <div className="tech-card-cost">
              <span className="gold-icon" />
              <span className="tech-card-cost-value">
                {formatNumber(tech.cost_next_level.gold)} Gold
              </span>
            </div>
          )}
          {tech.time_next_level_seconds != null && (
            <div className="tech-card-time">
              Time: {formatDuration(tech.time_next_level_seconds)}
            </div>
          )}
          <button
            className="tech-card-btn"
            onClick={() => onResearch(tech)}
            disabled={hasActiveElsewhere}
          >
            {hasActiveElsewhere ? 'BUSY' : 'RESEARCH'}
          </button>
        </>
      )}

      {isLocked && (
        <div className="tech-card-prereq">
          {getPrereqText(tech)}
        </div>
      )}
    </div>
  )
}

function TechCardResearchProgress({ finishAt }: { finishAt: string }) {
  const countdown = useCountdown(finishAt)
  const now = Date.now()
  const end = new Date(finishAt).getTime()
  // Rough progress calculation - we don't know original duration but can estimate
  const remaining = Math.max(0, end - now)
  // We'll show the countdown instead of percentage since we don't have start time
  const pctText = countdown || 'Complete!'

  return (
    <div className="tech-card-progress">
      <div className="tech-card-progress-bar">
        <div
          className="tech-card-progress-fill"
          style={{ width: remaining > 0 ? '60%' : '100%' }}
        />
      </div>
      <span className="tech-card-progress-text">{pctText}</span>
    </div>
  )
}

// ============ Active Research Bar ============

function ActiveResearchBar({
  active,
  onSpeedup,
  onCancel,
  pending,
}: {
  active: { tech_type_id: number; display_name: string; level: number; research_finish_at: string; tree: string }
  onSpeedup: () => void
  onCancel: () => void
  pending: boolean
}) {
  const countdown = useCountdown(active.research_finish_at)
  const now = Date.now()
  const end = new Date(active.research_finish_at).getTime()
  const totalEstimate = Math.max(end - now, 0)
  // We can't know the original total from the API response alone, but we can
  // show remaining time as the progress. For a better UX, show percentage based
  // on a rough estimate (original time might be in the tree data).
  const pct = totalEstimate <= 0 ? 100 : Math.max(5, 100 - Math.round((totalEstimate / (totalEstimate + 60000)) * 100))

  return (
    <div className="research-active-bar">
      <div className="research-active-info">
        <div>
          <span className="research-active-name">{active.display_name}</span>
          <span className="research-active-level"> Lv {active.level - 1} &rarr; {active.level}</span>
        </div>
        <span className="research-active-timer">{countdown || 'Completing...'}</span>
      </div>
      <div className="research-active-progress">
        <div className="research-progress-bar">
          <div className="research-progress-fill" style={{ width: `${pct}%` }} />
        </div>
        <div className="research-active-actions">
          <button
            className="research-speedup-btn"
            onClick={onSpeedup}
            disabled={pending}
          >
            SPEEDUP
          </button>
          <button
            className="research-cancel-btn"
            onClick={onCancel}
            disabled={pending}
          >
            CANCEL
          </button>
        </div>
      </div>
    </div>
  )
}

// ============ Confirmation Modal ============

function ResearchConfirmModal({
  tech,
  onConfirm,
  onCancel,
  pending,
  hasActiveResearch,
}: {
  tech: TechWithProgress
  onConfirm: () => void
  onCancel: () => void
  pending: boolean
  hasActiveResearch: boolean
}) {
  const nextLevel = tech.current_level + 1
  const goldCost = tech.cost_next_level?.gold ?? 0
  const timeSeconds = tech.time_next_level_seconds ?? 0

  const effectText = formatEffect(tech)
  const nextEffectText = formatNextEffect(tech)

  return createPortal(
    <div className="research-confirm-backdrop" onClick={e => { if (e.target === e.currentTarget) onCancel() }}>
      <div className="research-confirm-modal">
        <div className="research-confirm-header">
          Research: {tech.display_name} Lv {nextLevel}
        </div>
        <div className="research-confirm-body">
          {effectText && (
            <div className="research-confirm-row">
              <span className="research-confirm-label">Current</span>
              <span className="research-confirm-value">{effectText}</span>
            </div>
          )}
          {nextEffectText && (
            <div className="research-confirm-row">
              <span className="research-confirm-label">Next</span>
              <span className="research-confirm-value" style={{ color: '#44ccff' }}>{nextEffectText}</span>
            </div>
          )}

          <div className="research-confirm-sep" />

          <div className="research-confirm-row">
            <span className="research-confirm-label">Cost</span>
            <span className="research-confirm-value gold">
              {formatNumber(goldCost)} Gold
            </span>
          </div>
          <div className="research-confirm-row">
            <span className="research-confirm-label">Time</span>
            <span className="research-confirm-value">
              {formatDuration(timeSeconds)}
            </span>
          </div>

          {hasActiveResearch && (
            <>
              <div className="research-confirm-sep" />
              <div className="research-confirm-row">
                <span className="research-confirm-label" style={{ color: 'var(--accent-warning)' }}>
                  Warning: This will replace current research
                </span>
              </div>
            </>
          )}
        </div>
        <div className="research-confirm-actions">
          <button className="research-confirm-cancel" onClick={onCancel}>
            CANCEL
          </button>
          <button
            className="research-confirm-ok"
            onClick={onConfirm}
            disabled={pending}
          >
            {pending ? '...' : 'CONFIRM'}
          </button>
        </div>
      </div>
    </div>,
    document.body,
  )
}

// ============ Helpers ============

function getTreeLabel(tree: string): string {
  const found = TREES.find(t => t.key === tree)
  return found?.label || tree
}

function getCardState(tech: TechWithProgress): 'locked' | 'available' | 'researching' | 'completed' | 'maxed' {
  if (tech.is_researching) return 'researching'
  if (tech.current_level >= tech.max_level) return 'maxed'
  // A tech is locked if it has prerequisites and cost_next_level is null
  // (the backend returns null cost when prerequisites aren't met)
  if (tech.prerequisites.length > 0 && tech.cost_next_level === null) return 'locked'
  if (tech.current_level > 0) return 'completed'
  return 'available'
}

function getPrereqText(tech: TechWithProgress): string {
  if (tech.prerequisites.length === 0) return ''
  return 'Requires: ' + tech.prerequisites
    .map(p => `${p.tech.replace(/_/g, ' ')} Lv ${p.level}`)
    .join(', ')
}

function formatEffect(tech: TechWithProgress): string {
  if (tech.current_level === 0) return ''
  const eff = tech.effects
  if (!eff || !eff.type) return ''
  const perLevel = eff.per_level ?? 0
  const total = perLevel * tech.current_level
  const unit = eff.unit === 'percent' ? '%' : (eff.unit || '')
  const label = eff.type.replace(/_/g, ' ')
  return `+${total}${unit} ${label}`
}

function formatNextEffect(tech: TechWithProgress): string {
  const eff = tech.effects
  if (!eff || !eff.type) return ''
  const perLevel = eff.per_level ?? 0
  const nextTotal = perLevel * (tech.current_level + 1)
  const unit = eff.unit === 'percent' ? '%' : (eff.unit || '')
  const label = eff.type.replace(/_/g, ' ')
  return `+${nextTotal}${unit} ${label}`
}

function buildTiers(techs: TechWithProgress[]): TechWithProgress[][] {
  if (techs.length === 0) return []

  // Build a map of tech name -> tech
  const byName = new Map<string, TechWithProgress>()
  for (const t of techs) {
    byName.set(t.name, t)
  }

  // Calculate depth for each tech (longest prerequisite chain)
  const depthCache = new Map<string, number>()
  function getDepth(techName: string): number {
    if (depthCache.has(techName)) return depthCache.get(techName)!
    const tech = byName.get(techName)
    if (!tech || tech.prerequisites.length === 0) {
      depthCache.set(techName, 0)
      return 0
    }
    let maxDepth = 0
    for (const prereq of tech.prerequisites) {
      if (byName.has(prereq.tech)) {
        maxDepth = Math.max(maxDepth, getDepth(prereq.tech) + 1)
      }
    }
    depthCache.set(techName, maxDepth)
    return maxDepth
  }

  // Calculate all depths
  for (const t of techs) {
    getDepth(t.name)
  }

  // Group by depth
  const tierMap = new Map<number, TechWithProgress[]>()
  for (const t of techs) {
    const d = depthCache.get(t.name) ?? 0
    const existing = tierMap.get(d) || []
    existing.push(t)
    tierMap.set(d, existing)
  }

  // Sort by tier index
  const maxTier = Math.max(...Array.from(tierMap.keys()))
  const result: TechWithProgress[][] = []
  for (let i = 0; i <= maxTier; i++) {
    if (tierMap.has(i)) {
      result.push(tierMap.get(i)!)
    }
  }

  return result
}
