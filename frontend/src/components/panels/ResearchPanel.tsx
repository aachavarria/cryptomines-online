import { useState, useEffect, useMemo } from 'react'
import { createPortal } from 'react-dom'
import { useResearch } from '../../hooks/useResearch.ts'
import { useCountdown, formatDuration, formatNumber } from '../../hooks/useCountdown.ts'
import {
  FlaskConical,
  Microscope,
  Lightbulb,
  Lock,
  Check,
  Sparkles,
  X,
  Zap,
} from 'lucide-react'
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
  const [selectedId, setSelectedId] = useState<number | null>(null)

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

  // Auto-select first available tech when tab changes
  useEffect(() => {
    if (currentTreeTechs.length === 0) {
      setSelectedId(null)
      return
    }
    if (!selectedId || !currentTreeTechs.find(t => t.id === selectedId)) {
      const first = currentTreeTechs.find(t => getCardState(t) === 'available')
        ?? currentTreeTechs[0]
      setSelectedId(first.id)
    }
  }, [currentTreeTechs, selectedId])

  // Check which tree has active research
  const activeTreeKey = active?.tree || null

  // Calculate active bonuses across all trees
  const activeBonuses = useMemo(() => {
    const results: { text: string; key: string }[] = []
    Object.values(trees).forEach(treeTechs => {
      treeTechs.forEach(tech => {
        if (tech.current_level > 0 && tech.effects?.type) {
          const text = formatEffect(tech)
          if (text) {
            results.push({ text, key: `${tech.id}-${tech.effects.type}` })
          }
        }
      })
    })
    return results
  }, [trees])

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
      <div className="research-panel ds-modal ds-modal--lg">
        {/* Header */}
        <div className="ds-modal-header">
          <div className="ds-row">
            <Microscope size={22} strokeWidth={1.75} style={{ color: 'var(--ds-teal)' }} />
            <h2 className="ds-modal-title">Technology Center</h2>
          </div>
          <button className="ds-btn-icon" onClick={onClose} aria-label="Close">
            <X size={16} />
          </button>
        </div>

        {/* Tabs */}
        <div className="ds-tabs research-tabs">
          {TREES.map(t => (
            <button
              key={t.key}
              className="ds-tab"
              aria-selected={tab === t.key}
              onClick={() => { setTab(t.key); setSelectedId(null) }}
            >
              {t.label}
              {activeTreeKey === t.key && <span className="research-tab-dot" aria-label="active research" />}
            </button>
          ))}
        </div>

        {/* Toast */}
        {toast && (
          <div className="quest-reward-flash">
            <Sparkles size={14} /> {toast}
          </div>
        )}

        {/* Active Bonuses Summary */}
        {activeBonuses.length > 0 && (
          <div className="research-bonuses-summary">
            <div className="ds-caption">Active Bonuses</div>
            <div className="research-bonuses-grid">
              {activeBonuses.map(bonus => (
                <span key={bonus.key} className="ds-badge ds-badge--teal">
                  {bonus.text}
                </span>
              ))}
            </div>
          </div>
        )}

        {/* Body */}
        <div className="research-body">
          {loading ? (
            <div className="quest-loading">
              <div className="loading-spinner" /><span>Loading research...</span>
            </div>
          ) : error ? (
            <div className="quest-error">{error}</div>
          ) : currentTreeTechs.length === 0 ? (
            <div className="quest-empty">No technologies available in this tree.</div>
          ) : (
            <div className="research-tree-grid">
              {tiers.map((tierTechs, tierIndex) => (
                <div key={tierIndex} className="research-tree-tier">
                  {tierIndex > 0 && <div className="research-tier-connector" aria-hidden />}
                  <div className="research-tier-cards">
                    {tierTechs.map(tech => (
                      <TechCard
                        key={tech.id}
                        tech={tech}
                        active={active}
                        onResearch={handleStartResearch}
                        selected={tech.id === selectedId}
                        onSelect={() => setSelectedId(tech.id)}
                      />
                    ))}
                  </div>
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
  selected,
  onSelect,
}: {
  tech: TechWithProgress
  active: { tech_type_id: number } | null
  onResearch: (tech: TechWithProgress) => void
  selected: boolean
  onSelect: () => void
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

  // Pick an icon by tree-agnostic state.
  const Icon = isMaxed ? Sparkles : isLocked ? Lock : isCompleted ? Lightbulb : FlaskConical

  return (
    <div
      className={`ds-card tech-card ${stateClass} ${selected ? 'is-selected' : ''}`}
      onClick={onSelect}
      role="button"
      tabIndex={0}
      onKeyDown={e => { if (e.key === 'Enter' || e.key === ' ') onSelect() }}
    >
      <div className="tech-card-header">
        <div className="tech-card-icon">
          <Icon size={16} strokeWidth={1.75} />
        </div>
        <span className="tech-card-name">{tech.display_name}</span>
        {isMaxed ? (
          <span className="ds-badge ds-badge--orange">MAX</span>
        ) : isResearching ? (
          <span className="ds-badge ds-badge--teal">…</span>
        ) : (
          <span className="ds-mono ds-text-muted tech-card-level">
            {tech.current_level}/{tech.max_level}
          </span>
        )}
      </div>

      {effectText && (
        <div className="tech-card-effect ds-text-muted">{effectText}</div>
      )}

      {isResearching && tech.research_finish_at && (
        <TechCardResearchProgress finishAt={tech.research_finish_at} />
      )}

      {!isMaxed && !isResearching && !isLocked && (
        <>
          {nextEffectText && (
            <div className="tech-card-next">
              <span className="ds-caption">Next</span> <span>{nextEffectText}</span>
            </div>
          )}
          {tech.cost_next_level && (
            <div className="tech-card-cost">
              <span className="ds-resource-dot ds-resource-dot--gold" />
              <span className="ds-mono">{formatNumber(tech.cost_next_level.gold)}</span>
              <span className="ds-text-muted">Gold</span>
            </div>
          )}
          {tech.time_next_level_seconds != null && (
            <div className="tech-card-time ds-text-muted">
              Time <span className="ds-mono">{formatDuration(tech.time_next_level_seconds)}</span>
            </div>
          )}
          <button
            className="ds-btn-secondary ds-btn--sm ds-btn--block"
            onClick={e => { e.stopPropagation(); onResearch(tech) }}
            disabled={hasActiveElsewhere}
          >
            {hasActiveElsewhere ? 'BUSY' : 'RESEARCH'}
          </button>
        </>
      )}

      {isLocked && (
        <div className="tech-card-prereq ds-text-soft">
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
  const remaining = Math.max(0, end - now)
  const pctText = countdown || 'Complete!'

  return (
    <div className="tech-card-progress">
      <div className="ds-bar">
        <div
          className="ds-bar-fill"
          style={{ width: remaining > 0 ? '60%' : '100%' }}
        />
      </div>
      <span className="tech-card-progress-text ds-mono">{pctText}</span>
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
  const pct = totalEstimate <= 0 ? 100 : Math.max(5, 100 - Math.round((totalEstimate / (totalEstimate + 60000)) * 100))

  return (
    <div className="research-active-bar">
      <div className="research-active-info">
        <div className="ds-row">
          <FlaskConical size={16} strokeWidth={1.75} style={{ color: 'var(--ds-teal)' }} />
          <span className="research-active-name">{active.display_name}</span>
          <span className="ds-text-muted ds-mono"> Lv {active.level - 1} → {active.level}</span>
        </div>
        <span className="ds-mono research-active-timer">{countdown || 'Completing…'}</span>
      </div>
      <div className="research-active-progress">
        <div className="ds-bar">
          <div className="ds-bar-fill" style={{ width: `${pct}%` }} />
        </div>
        <div className="research-active-actions">
          <button
            className="ds-btn-ghost ds-btn--sm"
            onClick={onSpeedup}
            disabled={pending}
          >
            <Zap size={14} /> SPEEDUP
          </button>
          <button
            className="ds-btn-danger ds-btn--sm"
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
    <div className="ds-modal-backdrop research-confirm-backdrop" onClick={e => { if (e.target === e.currentTarget) onCancel() }}>
      <div className="ds-modal ds-modal--sm research-confirm-modal">
        <div className="ds-modal-header">
          <h3 className="ds-modal-title">Research: {tech.display_name} Lv {nextLevel}</h3>
        </div>
        <div className="ds-modal-body">
          <div className="research-confirm-rows">
            {effectText && (
              <div className="research-confirm-row">
                <span className="ds-text-muted">Current</span>
                <span>{effectText}</span>
              </div>
            )}
            {nextEffectText && (
              <div className="research-confirm-row">
                <span className="ds-text-muted">Next</span>
                <span style={{ color: 'var(--ds-teal-dark)' }}>{nextEffectText}</span>
              </div>
            )}
            <hr className="ds-divider" />
            <div className="research-confirm-row">
              <span className="ds-text-muted">Cost</span>
              <span><span className="ds-mono">{formatNumber(goldCost)}</span> Gold</span>
            </div>
            <div className="research-confirm-row">
              <span className="ds-text-muted">Time</span>
              <span className="ds-mono">{formatDuration(timeSeconds)}</span>
            </div>
            {hasActiveResearch && (
              <>
                <hr className="ds-divider" />
                <div className="research-confirm-row">
                  <span style={{ color: 'var(--ds-orange-strong)' }}>
                    Warning: This will replace current research
                  </span>
                </div>
              </>
            )}
          </div>
        </div>
        <div className="ds-modal-footer">
          <button className="ds-btn-ghost" onClick={onCancel}>CANCEL</button>
          <button
            className="ds-btn-primary"
            onClick={onConfirm}
            disabled={pending}
          >
            {pending ? '…' : <><Check size={14} /> CONFIRM</>}
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

function getUnitSuffix(eff: { unit?: string }): string {
  if (!eff.unit) return ''
  if (eff.unit === 'percent' || eff.unit === 'percent_per_level') return '%'
  if (eff.unit === 'flat') return ''
  return ` ${eff.unit}`
}

function formatLabel(type: string): string {
  return type.replace(/_/g, ' ')
}

/** Format the effect value at a specific level index (0-based) for the given effect object. */
function formatEffectAtLevel(eff: Record<string, any>, levelIdx: number): string {
  if (!eff || !eff.type) return ''
  const unit = getUnitSuffix(eff)
  const label = formatLabel(eff.type)

  if (typeof eff.enabled === 'boolean') {
    return eff.enabled ? `${label} enabled` : ''
  }

  if (Array.isArray(eff.values)) {
    if (levelIdx < 0) return ''
    const val = levelIdx < eff.values.length ? eff.values[levelIdx] : eff.values[eff.values.length - 1]
    return `+${val}${unit} ${label}`
  }

  if (eff.type === 'range_damage' && Array.isArray(eff.ranges)) {
    const parts = [`Range dmg: ${eff.ranges.join('/')}%`]
    if (eff.crit_rate) parts.push(`+${eff.crit_rate}% crit rate`)
    if (eff.crit_damage) parts.push(`+${eff.crit_damage}% crit dmg`)
    return parts.join(', ')
  }

  if (eff.type === 'multi_bonus') {
    const skip = new Set(['type', 'unit'])
    const parts: string[] = []
    for (const [k, v] of Object.entries(eff)) {
      if (skip.has(k)) continue
      if (typeof v !== 'number') continue
      const perLevel = k.endsWith('_per_level')
      const cleanKey = formatLabel(perLevel ? k.replace(/_per_level$/, '') : k)
      const val = perLevel ? v * (levelIdx + 1) : v
      const sign = val >= 0 ? '+' : ''
      parts.push(`${sign}${val}${unit} ${cleanKey}`)
    }
    return parts.join(', ')
  }

  if (Array.isArray(eff.applies_to)) {
    return `${label}: ${eff.applies_to.map((s: string) => formatLabel(s)).join(', ')}`
  }

  if (eff.flat !== undefined) {
    const sign = eff.flat >= 0 ? '+' : ''
    return `${sign}${eff.flat}${unit} ${label}`
  }

  if (eff.per_level !== undefined) {
    const total = eff.per_level * (levelIdx + 1)
    const sign = total >= 0 ? '+' : ''
    const parts = [`${sign}${total}${unit} ${label}`]

    const skip = new Set(['type', 'per_level', 'unit'])
    for (const [k, v] of Object.entries(eff)) {
      if (skip.has(k)) continue
      if (typeof v !== 'number') continue
      const perLevel = k.endsWith('_per_level')
      const cleanKey = formatLabel(perLevel ? k.replace(/_per_level$/, '') : k)
      const val = perLevel ? v * (levelIdx + 1) : v
      const s = val >= 0 ? '+' : ''
      parts.push(`${s}${val}${unit} ${cleanKey}`)
    }
    if (typeof eff.chance_for_half === 'boolean' && eff.chance_for_half) {
      parts.push('50% cost reduction')
    }
    return parts.join(', ')
  }

  const skip = new Set(['type', 'unit'])
  const parts: string[] = []
  for (const [k, v] of Object.entries(eff)) {
    if (skip.has(k)) continue
    if (typeof v !== 'number') continue
    const val = v * (levelIdx + 1)
    const sign = val >= 0 ? '+' : ''
    parts.push(`${sign}${val}${unit} vs ${formatLabel(k)}`)
  }
  if (parts.length > 0) return parts.join(', ')

  return label
}

function formatEffect(tech: TechWithProgress): string {
  if (tech.current_level === 0) return ''
  return formatEffectAtLevel(tech.effects, tech.current_level - 1)
}

function formatNextEffect(tech: TechWithProgress): string {
  return formatEffectAtLevel(tech.effects, tech.current_level)
}

function buildTiers(techs: TechWithProgress[]): TechWithProgress[][] {
  if (techs.length === 0) return []

  const byName = new Map<string, TechWithProgress>()
  for (const t of techs) {
    byName.set(t.name, t)
  }

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

  for (const t of techs) {
    getDepth(t.name)
  }

  const tierMap = new Map<number, TechWithProgress[]>()
  for (const t of techs) {
    const d = depthCache.get(t.name) ?? 0
    const existing = tierMap.get(d) || []
    existing.push(t)
    tierMap.set(d, existing)
  }

  const maxTier = Math.max(...Array.from(tierMap.keys()))
  const result: TechWithProgress[][] = []
  for (let i = 0; i <= maxTier; i++) {
    if (tierMap.has(i)) {
      result.push(tierMap.get(i)!)
    }
  }

  return result
}
