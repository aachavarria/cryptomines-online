import { useState } from 'react'
import { useBlueprints } from '../../hooks/useBlueprints.ts'
import type { Blueprint, PlayerBlueprint } from '../../types'

const BP_TYPE_LABELS: Record<string, string> = {
  hull: 'Hull',
  module: 'Module',
}

export default function BlueprintPanel() {
  const { allBlueprints, myBlueprints, loading, error, activate, hasActivated, hasOwned } = useBlueprints()
  const [filter, setFilter] = useState<'all' | 'hull' | 'module'>('all')
  const [activating, setActivating] = useState<number | null>(null)

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

  return (
    <div className="p2-panel">
      <div className="p2-panel-header">
        <div className="p2-panel-icon bp-icon">BP</div>
        <div>
          <div className="p2-panel-title">Blueprints</div>
          <div className="p2-panel-subtitle">
            {myBlueprints.filter(bp => bp.activated).length} activated / {myBlueprints.length} owned
          </div>
        </div>
      </div>

      {error && <div className="p2-error-msg">{error}</div>}

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

          return (
            <div
              key={bp.id}
              className={`bp-card ${activated ? 'activated' : owned ? 'owned' : 'locked'}`}
            >
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
              <div className="bp-card-status">
                {activated ? (
                  <span className="bp-status-active">Activated</span>
                ) : owned ? (
                  <button
                    className="p2-btn p2-btn-success p2-btn-sm"
                    onClick={() => handleActivate(bp.id)}
                    disabled={activating === bp.id}
                  >
                    {activating === bp.id ? '...' : 'Activate'}
                  </button>
                ) : (
                  <span className="bp-status-locked">Not Owned</span>
                )}
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}
