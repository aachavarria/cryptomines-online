import { useState, useEffect } from 'react'
import { createPortal } from 'react-dom'
import LoadingButton from '../common/LoadingButton'
import { useBlueprints } from '../../hooks/useBlueprints.ts'
import { useCountdown, formatDuration, formatNumber } from '../../hooks/useCountdown.ts'
import { devGiveBlueprint } from '../../services/api.ts'
import type { Blueprint, PlayerBlueprint } from '../../types'

const BP_TYPE_LABELS: Record<string, string> = {
  hull: 'Hull',
  module: 'Module',
}

export default function BlueprintPanel() {
  const { allBlueprints, myBlueprints, activeResearch, loading, error, activate, research, hasActivated, hasOwned, refresh } = useBlueprints()
  const [filter, setFilter] = useState<'all' | 'hull' | 'module'>('all')
  const [activating, setActivating] = useState<number | null>(null)
  const [researchingBp, setResearchingBp] = useState<number | null>(null)
  const [confirmResearch, setConfirmResearch] = useState<{ bp: Blueprint; playerBp: PlayerBlueprint } | null>(null)
  const [devGiving, setDevGiving] = useState<number | null>(null)
  const isDevMode = localStorage.getItem('dev_mode') === 'true'

  // Poll for active research completion
  useEffect(() => {
    if (!activeResearch) return
    const interval = setInterval(() => {
      const finishTime = new Date(activeResearch.research_finish_at).getTime()
      if (Date.now() >= finishTime) {
        refresh()
      }
    }, 2000)
    return () => clearInterval(interval)
  }, [activeResearch, refresh])

  if (loading) {
    return <div className="p2-panel-loading"><div className="loading-spinner" /><span>Loading Blueprints...</span></div>
  }

  const filtered = filter === 'all'
    ? allBlueprints
    : allBlueprints.filter(bp => bp.blueprint_type === filter)

  async function handleActivate(bpId: number) {
    setActivating(bpId)
    try {
      await activate(bpId)
    } finally {
      setActivating(null)
    }
  }

  async function handleDevGive(bpId: number) {
    setDevGiving(bpId)
    try {
      await devGiveBlueprint(bpId, 1)
      await refresh()
    } catch {
      // dev endpoint failed
    } finally {
      setDevGiving(null)
    }
  }

  async function handleResearch(bp: Blueprint, playerBp: PlayerBlueprint) {
    setConfirmResearch({ bp, playerBp })
  }

  async function handleConfirmResearch() {
    if (!confirmResearch) return
    setResearchingBp(confirmResearch.bp.id)
    try {
      await research(confirmResearch.bp.id)
      setConfirmResearch(null)
    } catch {
      // Error handled by hook
    } finally {
      setResearchingBp(null)
    }
  }

  return (
    <div className="p2-panel">
      <div className="p2-panel-header">
        <div className="p2-panel-icon bp-icon">BP</div>
        <div>
          <div className="p2-panel-title">Blueprints</div>
          <div className="p2-panel-subtitle">
            {myBlueprints.filter(bp => bp.is_activated).length} activated / {myBlueprints.length} owned
          </div>
        </div>
      </div>

      {error && <div className="p2-error-msg">{error}</div>}

      {/* Active Research Bar */}
      {activeResearch && (
        <ActiveResearchBar active={activeResearch} />
      )}

      <div className="p2-tabs">
        {(['all', 'hull', 'module'] as const).map(f => (
          <button
            key={f}
            className={`p2-tab ${filter === f ? 'active' : ''}`}
            onClick={() => setFilter(f)}
          >
            {f === 'all' ? 'All' : f.charAt(0).toUpperCase() + f.slice(1) + 's'}
          </button>
        ))}
      </div>

      <div className="bp-grid">
        {filtered.map(bp => {
          const owned = hasOwned(bp.id)
          const activated = hasActivated(bp.id)
          const playerBp = myBlueprints.find(pbp => pbp.blueprint_id === bp.id)
          const researchLevel = playerBp?.research_level ?? 0
          const isResearching = activeResearch?.blueprint_id === bp.id
          const canResearch = activated && researchLevel < 3 && !isResearching && !activeResearch

          const cardClass = isResearching
            ? 'bp-card bp-card--researching'
            : researchLevel >= 3
              ? 'bp-card bp-card--maxed'
              : activated
                ? 'bp-card activated'
                : owned
                  ? 'bp-card owned'
                  : 'bp-card locked'

          return (
            <div key={bp.id} className={cardClass}>
              <div className="bp-card-type">
                <span className={`bp-type-badge ${bp.blueprint_type}`}>
                  {BP_TYPE_LABELS[bp.blueprint_type]}
                </span>
                {bp.hull_class && (
                  <span className={`bp-class-badge ${bp.hull_class}`}>
                    {bp.hull_class}
                  </span>
                )}
                {bp.module_category && (
                  <span className="bp-cat-badge">{bp.module_category}</span>
                )}
              </div>

              <div className="bp-card-name">{bp.display_name}</div>

              {/* Research Level Stars */}
              {owned && (
                <div className="bp-research-stars">
                  {[1, 2, 3].map(level => (
                    <span
                      key={level}
                      className={`bp-star ${level <= researchLevel ? 'filled' : 'empty'}`}
                    >
                      ★
                    </span>
                  ))}
                </div>
              )}

              <div className="bp-card-status">
                {!owned ? (
                  <>
                    <span className="bp-status-locked">🔒 Not Owned</span>
                    {isDevMode && (
                      <LoadingButton
                        className="p2-btn p2-btn-dev p2-btn-sm"
                        onClick={() => handleDevGive(bp.id)}
                        loading={devGiving === bp.id}
                      >
                        DEV Get
                      </LoadingButton>
                    )}
                  </>
                ) : !activated ? (
                  <LoadingButton
                    className="p2-btn p2-btn-success p2-btn-sm"
                    onClick={() => handleActivate(bp.id)}
                    loading={activating === bp.id}
                  >
                    Activate
                  </LoadingButton>
                ) : researchLevel >= 3 ? (
                  <span className="bp-status-maxed">Max Research</span>
                ) : isResearching ? (
                  <span className="bp-status-researching">Researching...</span>
                ) : canResearch && playerBp ? (
                  <LoadingButton
                    className="p2-btn p2-btn-research p2-btn-sm"
                    onClick={() => handleResearch(bp, playerBp)}
                    loading={researchingBp === bp.id}
                  >
                    {`Research Lv${researchLevel + 1}`}
                  </LoadingButton>
                ) : activeResearch ? (
                  <span className="bp-status-busy">WRC Busy</span>
                ) : (
                  <span className="bp-status-active">Activated</span>
                )}
              </div>
            </div>
          )
        })}
      </div>

      {/* Research Confirmation Modal */}
      {confirmResearch && (
        <ResearchConfirmModal
          bp={confirmResearch.bp}
          playerBp={confirmResearch.playerBp}
          onConfirm={handleConfirmResearch}
          onCancel={() => setConfirmResearch(null)}
          pending={researchingBp !== null}
        />
      )}
    </div>
  )
}

// ============ Active Research Bar ============

function ActiveResearchBar({ active }: { active: { blueprint_name: string; target_level: number; research_finish_at: string } }) {
  const countdown = useCountdown(active.research_finish_at)
  const now = Date.now()
  const end = new Date(active.research_finish_at).getTime()
  const remaining = Math.max(0, end - now)
  // Rough progress estimate based on 1 hour per level
  const totalTime = 3600 * active.target_level * 1000
  const elapsed = totalTime - remaining
  const pct = Math.min(100, Math.max(5, Math.round((elapsed / totalTime) * 100)))

  return (
    <div className="bp-active-research">
      <div className="bp-active-info">
        <div>
          <span className="bp-active-name">{active.blueprint_name}</span>
          <span className="bp-active-level"> → Lv {active.target_level}</span>
        </div>
        <span className="bp-active-timer">{countdown || 'Completing...'}</span>
      </div>
      <div className="bp-progress-bar">
        <div className="bp-progress-fill" style={{ width: `${pct}%` }} />
      </div>
    </div>
  )
}

// ============ Research Confirmation Modal ============

function ResearchConfirmModal({
  bp,
  playerBp,
  onConfirm,
  onCancel,
  pending,
}: {
  bp: Blueprint
  playerBp: PlayerBlueprint
  onConfirm: () => void
  onCancel: () => void
  pending: boolean
}) {
  const targetLevel = playerBp.research_level + 1
  const baseCost = 10000 * targetLevel
  const metalCost = baseCost
  const he3Cost = Math.floor(baseCost * 3 / 4)
  const goldCost = Math.floor(baseCost / 2)
  const researchTime = 3600 * targetLevel // 1 hour per level

  return createPortal(
    <div className="bp-confirm-backdrop" onClick={e => { if (e.target === e.currentTarget) onCancel() }}>
      <div className="bp-confirm-modal">
        <div className="bp-confirm-header">
          Research: {bp.display_name} Lv{targetLevel}
        </div>
        <div className="bp-confirm-body">
          <div className="bp-confirm-row">
            <span className="bp-confirm-label">Current Level</span>
            <span className="bp-confirm-value">Lv {playerBp.research_level}</span>
          </div>
          <div className="bp-confirm-row">
            <span className="bp-confirm-label">Target Level</span>
            <span className="bp-confirm-value bp-next">Lv {targetLevel}</span>
          </div>

          <div className="bp-confirm-sep" />

          <div className="bp-confirm-row">
            <span className="bp-confirm-label">Metal</span>
            <span className="bp-confirm-value">{formatNumber(metalCost)}</span>
          </div>
          <div className="bp-confirm-row">
            <span className="bp-confirm-label">He3</span>
            <span className="bp-confirm-value">{formatNumber(he3Cost)}</span>
          </div>
          <div className="bp-confirm-row">
            <span className="bp-confirm-label">Gold</span>
            <span className="bp-confirm-value gold">{formatNumber(goldCost)}</span>
          </div>
          <div className="bp-confirm-row">
            <span className="bp-confirm-label">Time</span>
            <span className="bp-confirm-value">{formatDuration(researchTime)}</span>
          </div>

          <div className="bp-confirm-note">
            Research will be conducted at the Weapon Research Center.
          </div>
        </div>
        <div className="bp-confirm-actions">
          <button className="bp-confirm-cancel" onClick={onCancel}>
            CANCEL
          </button>
          <LoadingButton
            className="bp-confirm-ok"
            onClick={onConfirm}
            loading={pending}
          >
            START RESEARCH
          </LoadingButton>
        </div>
      </div>
    </div>,
    document.body,
  )
}
