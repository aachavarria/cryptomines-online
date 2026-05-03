import { useState, useEffect } from 'react'
import { createPortal } from 'react-dom'
import { FileText, Sparkles, Lock, Star, X, FlaskConical } from 'lucide-react'
import LoadingButton from '../common/LoadingButton'
import { useBlueprints } from '../../hooks/useBlueprints.ts'
import { useCountdown, formatDuration, formatNumber } from '../../hooks/useCountdown.ts'
import { devGiveBlueprint } from '../../services/api.ts'
import type { Blueprint, PlayerBlueprint } from '../../types'

const BP_TYPE_LABELS: Record<string, string> = {
  hull: 'Hull',
  module: 'Module',
}

// Map blueprint metadata to a rarity bucket so visuals match the design system.
// Rarity ordering: hull battleship/T3 = epic/legendary, cruiser = rare, frigate = uncommon,
// modules default to common, planetary = legendary.
function getBlueprintRarity(bp: Blueprint): 'common' | 'uncommon' | 'rare' | 'epic' | 'legendary' {
  if (bp.blueprint_type === 'hull') {
    if (bp.hull_class === 'battleship') return 'epic'
    if (bp.hull_class === 'cruiser') return 'rare'
    if (bp.hull_class === 'frigate') return 'uncommon'
    return 'common'
  }
  if (bp.module_category === 'planetary') return 'legendary'
  if (bp.module_category === 'missile' || bp.module_category === 'ship_based') return 'rare'
  if (bp.module_category === 'directional' || bp.module_category === 'shield') return 'uncommon'
  return 'common'
}

export default function BlueprintPanel() {
  const {
    allBlueprints,
    myBlueprints,
    activeResearch,
    loading,
    error,
    activate,
    research,
    hasActivated,
    hasOwned,
    refresh,
  } = useBlueprints()
  const [filter, setFilter] = useState<'all' | 'hull' | 'module'>('all')
  const [activating, setActivating] = useState<number | null>(null)
  const [researchingBp, setResearchingBp] = useState<number | null>(null)
  const [confirmResearch, setConfirmResearch] =
    useState<{ bp: Blueprint; playerBp: PlayerBlueprint } | null>(null)
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
    return (
      <div className="ds-panel bp-loading">
        <div className="loading-spinner" />
        <span>Loading Blueprints...</span>
      </div>
    )
  }

  const filtered =
    filter === 'all' ? allBlueprints : allBlueprints.filter(bp => bp.blueprint_type === filter)

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
      // ignore
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
      // hook handles
    } finally {
      setResearchingBp(null)
    }
  }

  const activatedCount = myBlueprints.filter(bp => bp.is_activated).length
  const ownedCount = myBlueprints.length

  return (
    <div className="ds-panel bp-panel">
      <header className="bp-header">
        <div className="bp-header-icon">
          <FileText size={22} strokeWidth={1.75} />
        </div>
        <div className="bp-header-text">
          <h2 className="ds-h2">Blueprints</h2>
          <p className="ds-text-muted bp-header-sub">
            <span className="ds-mono">{activatedCount}</span> activated /{' '}
            <span className="ds-mono">{ownedCount}</span> owned
          </p>
        </div>
      </header>

      {error && (
        <div className="ds-badge ds-badge--danger bp-error" role="alert">
          {error}
        </div>
      )}

      {activeResearch && <ActiveResearchBar active={activeResearch} />}

      <div className="ds-tabs bp-tabs" role="tablist">
        {(['all', 'hull', 'module'] as const).map(f => (
          <button
            key={f}
            className="ds-tab"
            role="tab"
            aria-selected={filter === f}
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
          const rarity = getBlueprintRarity(bp)

          return (
            <div
              key={bp.id}
              className={`ds-card bp-card${activated ? ' bp-card--activated' : ''}${!owned ? ' bp-card--locked' : ''}`}
            >
              <div className="bp-card-tags">
                <span className={`ds-badge ds-badge--rarity-${rarity}`}>
                  {BP_TYPE_LABELS[bp.blueprint_type]}
                </span>
                {bp.hull_class && (
                  <span className="ds-badge ds-badge--neutral">{bp.hull_class}</span>
                )}
                {bp.module_category && (
                  <span className="ds-badge ds-badge--neutral">{bp.module_category}</span>
                )}
              </div>

              <div className="bp-card-name">{bp.display_name}</div>

              {owned && (
                <div className="bp-stars" aria-label={`Research level ${researchLevel} of 3`}>
                  {[1, 2, 3].map(level => (
                    <Star
                      key={level}
                      size={14}
                      strokeWidth={1.75}
                      className={level <= researchLevel ? 'bp-star bp-star--filled' : 'bp-star'}
                      fill={level <= researchLevel ? 'currentColor' : 'none'}
                    />
                  ))}
                </div>
              )}

              <div className="bp-card-status">
                {!owned ? (
                  <>
                    <span className="ds-badge ds-badge--neutral bp-locked-badge">
                      <Lock size={12} strokeWidth={2} /> Not Owned
                    </span>
                    {isDevMode && (
                      <LoadingButton
                        className="ds-btn ds-btn-ghost ds-btn--sm"
                        onClick={() => handleDevGive(bp.id)}
                        loading={devGiving === bp.id}
                      >
                        DEV Get
                      </LoadingButton>
                    )}
                  </>
                ) : !activated ? (
                  <LoadingButton
                    className="ds-btn ds-btn-secondary ds-btn--sm"
                    onClick={() => handleActivate(bp.id)}
                    loading={activating === bp.id}
                  >
                    Activate
                  </LoadingButton>
                ) : researchLevel >= 3 ? (
                  <span className="ds-badge ds-badge--orange">Max Research</span>
                ) : isResearching ? (
                  <span className="ds-badge ds-badge--info">Researching…</span>
                ) : canResearch && playerBp ? (
                  <LoadingButton
                    className="ds-btn ds-btn-secondary ds-btn--sm"
                    onClick={() => handleResearch(bp, playerBp)}
                    loading={researchingBp === bp.id}
                  >
                    {`Research Lv${researchLevel + 1}`}
                  </LoadingButton>
                ) : activeResearch ? (
                  <span className="ds-badge ds-badge--warning">WRC Busy</span>
                ) : (
                  <span className="ds-badge ds-badge--success">Activated</span>
                )}
              </div>
            </div>
          )
        })}
      </div>

      {confirmResearch && (
        <ResearchConfirmModal
          bp={confirmResearch.bp}
          playerBp={confirmResearch.playerBp}
          onConfirm={handleConfirmResearch}
          onCancel={() => setConfirmResearch(null)}
          pending={researchingBp !== null}
        />
      )}

      <style>{`
        .bp-panel { display: flex; flex-direction: column; gap: var(--sp-4); max-width: 1100px; }
        .bp-loading { display: flex; align-items: center; gap: var(--sp-3); }
        .bp-header {
          display: flex; align-items: center; gap: var(--sp-3);
        }
        .bp-header-icon {
          width: 44px; height: 44px;
          display: grid; place-items: center;
          background: var(--ds-teal-tint);
          color: var(--ds-teal-dark);
          border-radius: var(--r-lg);
        }
        .bp-header-text h2 { margin: 0; }
        .bp-header-sub { font-size: var(--fs-sm); margin: 2px 0 0; }

        .bp-error { display: inline-flex; }

        .bp-tabs { margin-top: var(--sp-1); }

        .bp-grid {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
          gap: var(--sp-3);
        }
        .bp-card {
          display: flex; flex-direction: column; gap: var(--sp-2);
          padding: var(--sp-4);
          transition: border-color var(--motion-fast), box-shadow var(--motion-fast);
        }
        .bp-card--activated { border-color: var(--ds-teal); }
        .bp-card--locked { opacity: 0.7; }
        .bp-card-tags { display: flex; flex-wrap: wrap; gap: var(--sp-1); }
        .bp-card-name {
          font-size: var(--fs-h3); font-weight: var(--fw-semibold);
          color: var(--ds-text);
        }
        .bp-stars {
          display: inline-flex; gap: 2px;
          color: var(--ds-text-soft);
        }
        .bp-star { color: var(--ds-text-soft); }
        .bp-star--filled { color: var(--ds-orange); }
        .bp-card-status {
          margin-top: auto; padding-top: var(--sp-2);
          display: flex; align-items: center; gap: var(--sp-2); flex-wrap: wrap;
        }
        .bp-locked-badge { display: inline-flex; align-items: center; gap: 4px; }
      `}</style>
    </div>
  )
}

// ============ Active Research Bar ============

function ActiveResearchBar({
  active,
}: {
  active: { blueprint_name: string; target_level: number; research_finish_at: string }
}) {
  const countdown = useCountdown(active.research_finish_at)
  const now = Date.now()
  const end = new Date(active.research_finish_at).getTime()
  const remaining = Math.max(0, end - now)
  const totalTime = 3600 * active.target_level * 1000
  const elapsed = totalTime - remaining
  const pct = Math.min(100, Math.max(5, Math.round((elapsed / totalTime) * 100)))

  return (
    <div className="ds-card bp-active">
      <div className="bp-active-row">
        <div className="bp-active-info">
          <FlaskConical size={18} strokeWidth={1.75} className="bp-active-icon" />
          <span className="bp-active-name">{active.blueprint_name}</span>
          <span className="ds-badge ds-badge--teal">Lv {active.target_level}</span>
        </div>
        <span className="ds-mono bp-active-timer">{countdown || 'Completing…'}</span>
      </div>
      <div className="ds-bar bp-active-bar">
        <div className="ds-bar-fill" style={{ width: `${pct}%` }} />
      </div>

      <style>{`
        .bp-active { display: flex; flex-direction: column; gap: var(--sp-2); }
        .bp-active-row { display: flex; align-items: center; justify-content: space-between; gap: var(--sp-3); }
        .bp-active-info { display: flex; align-items: center; gap: var(--sp-2); min-width: 0; }
        .bp-active-icon { color: var(--ds-teal); flex-shrink: 0; }
        .bp-active-name { font-weight: var(--fw-semibold); color: var(--ds-text); }
        .bp-active-timer { color: var(--ds-text-muted); font-size: var(--fs-sm); }
      `}</style>
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
  const he3Cost = Math.floor((baseCost * 3) / 4)
  const goldCost = Math.floor(baseCost / 2)
  const researchTime = 3600 * targetLevel

  return createPortal(
    <div
      className="ds-modal-backdrop"
      onClick={e => {
        if (e.target === e.currentTarget) onCancel()
      }}
    >
      <div className="ds-modal ds-modal--sm bp-confirm" role="dialog" aria-modal="true">
        <div className="ds-modal-header">
          <div className="bp-confirm-title">
            <Sparkles size={20} strokeWidth={1.75} />
            <h3 className="ds-modal-title">
              Research: {bp.display_name} Lv{targetLevel}
            </h3>
          </div>
          <button className="ds-btn-icon" onClick={onCancel} aria-label="Cancel">
            <X size={18} strokeWidth={1.75} />
          </button>
        </div>
        <div className="ds-modal-body bp-confirm-body">
          <dl className="bp-confirm-list">
            <Row label="Current Level" value={`Lv ${playerBp.research_level}`} />
            <Row label="Target Level" value={`Lv ${targetLevel}`} accent="orange" />
          </dl>
          <hr className="ds-divider" />
          <dl className="bp-confirm-list">
            <Row
              label={
                <span>
                  <span className="ds-resource-dot ds-resource-dot--metal" /> Metal
                </span>
              }
              value={formatNumber(metalCost)}
            />
            <Row
              label={
                <span>
                  <span className="ds-resource-dot ds-resource-dot--he3" /> He3
                </span>
              }
              value={formatNumber(he3Cost)}
            />
            <Row
              label={
                <span>
                  <span className="ds-resource-dot ds-resource-dot--gold" /> Gold
                </span>
              }
              value={formatNumber(goldCost)}
            />
            <Row label="Time" value={formatDuration(researchTime)} />
          </dl>
          <p className="ds-text-muted bp-confirm-note">
            Research will be conducted at the Weapon Research Center.
          </p>
        </div>
        <div className="ds-modal-footer">
          <button className="ds-btn ds-btn-ghost" onClick={onCancel}>
            Cancel
          </button>
          <LoadingButton
            className="ds-btn ds-btn-secondary"
            onClick={onConfirm}
            loading={pending}
          >
            Start Research
          </LoadingButton>
        </div>

        <style>{`
          .bp-confirm-title { display: flex; align-items: center; gap: var(--sp-2); color: var(--ds-text); }
          .bp-confirm-title h3 { margin: 0; }
          .bp-confirm-body { display: flex; flex-direction: column; gap: var(--sp-3); }
          .bp-confirm-list { display: flex; flex-direction: column; gap: var(--sp-2); margin: 0; }
          .bp-confirm-row {
            display: flex; align-items: center; justify-content: space-between;
            font-size: var(--fs-sm);
          }
          .bp-confirm-row dt, .bp-confirm-row dd { margin: 0; }
          .bp-confirm-row dt { color: var(--ds-text-muted); display: inline-flex; align-items: center; gap: var(--sp-2); }
          .bp-confirm-row dd { font-weight: var(--fw-semibold); }
          .bp-confirm-row dd.accent-orange { color: var(--ds-orange-strong); }
          .bp-confirm-note { font-size: var(--fs-sm); margin: 0; }
        `}</style>
      </div>
    </div>,
    document.body,
  )
}

function Row({
  label,
  value,
  accent,
}: {
  label: React.ReactNode
  value: React.ReactNode
  accent?: 'orange'
}) {
  return (
    <div className="bp-confirm-row">
      <dt>{label}</dt>
      <dd className={accent === 'orange' ? 'accent-orange' : ''}>{value}</dd>
    </div>
  )
}
