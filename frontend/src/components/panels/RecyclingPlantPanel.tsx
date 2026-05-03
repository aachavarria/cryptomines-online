import { useState, useMemo } from 'react'
import { Recycle, AlertTriangle, X, Coins } from 'lucide-react'
import { useRecycling } from '../../hooks/useRecycling.ts'
import { formatNumber, useCountdown } from '../../hooks/useCountdown.ts'
import LoadingButton from '../common/LoadingButton.tsx'
import type { RecyclingJob } from '../../types'

export default function RecyclingPlantPanel() {
  const { jobs, availableShips, loading, error, recycle, collect, cancel } = useRecycling()
  const [selectedShipId, setSelectedShipId] = useState<string>('')
  const [confirmRecycle, setConfirmRecycle] = useState(false)

  const { activeJobs, completedJobs } = useMemo(() => {
    const now = new Date()
    const active = jobs.filter(j => new Date(j.completed_at) > now && !j.collected)
    const completed = jobs.filter(j => new Date(j.completed_at) <= now && !j.collected)
    return { activeJobs: active, completedJobs: completed }
  }, [jobs])

  async function handleRecycle() {
    if (!selectedShipId) return
    const success = await recycle(selectedShipId)
    if (success) {
      setSelectedShipId('')
      setConfirmRecycle(false)
    }
  }

  async function handleCollect(jobId: string) {
    await collect(jobId)
  }

  async function handleCancel(jobId: string) {
    if (confirm('Cancel recycling? The ship cannot be recovered.')) {
      await cancel(jobId)
    }
  }

  return (
    <div className="ds-panel rp-panel">
      <header className="rp-header">
        <div className="rp-header-icon">
          <Recycle size={22} strokeWidth={1.75} />
        </div>
        <div className="rp-header-text">
          <h2 className="ds-h2">Recycling Plant</h2>
          <p className="ds-text-muted rp-subtitle">Recycle unwanted ships for resources</p>
        </div>
      </header>

      {error && (
        <div className="ds-badge ds-badge--danger rp-error" role="alert">
          {error}
        </div>
      )}

      {/* Available Ships */}
      <section className="rp-section">
        <h3 className="ds-h3">Available Ships</h3>
        {availableShips.length === 0 ? (
          <div className="ds-card rp-empty">
            <p>No ships available for recycling.</p>
            <p className="ds-text-muted">
              Ships must be unassigned from fleets to be recycled. Visit the Fleet panel to dismiss
              ships first.
            </p>
          </div>
        ) : (
          <div className="rp-ships-grid">
            {availableShips.map(ship => {
              const isSelected = selectedShipId === ship.id
              return (
                <button
                  key={ship.id}
                  type="button"
                  className={`ds-list-item rp-ship-card ${isSelected ? 'is-selected' : ''}`}
                  aria-selected={isSelected}
                  onClick={() => setSelectedShipId(ship.id)}
                >
                  <span className="rp-ship-name">{ship.design_name}</span>
                  <span className="ds-badge ds-badge--neutral">{ship.hull_class}</span>
                  <span className="ds-text-soft rp-ship-id ds-mono">
                    {ship.id.slice(0, 8)}
                  </span>
                </button>
              )
            })}
          </div>
        )}

        {selectedShipId && (
          <div className="rp-actions">
            <LoadingButton
              className="ds-btn ds-btn-secondary"
              onClick={() => setConfirmRecycle(true)}
              loading={loading}
            >
              Recycle Selected Ship
            </LoadingButton>
            <button className="ds-btn ds-btn-ghost" onClick={() => setSelectedShipId('')}>
              Clear
            </button>
          </div>
        )}
      </section>

      {/* Active Jobs */}
      {activeJobs.length > 0 && (
        <section className="rp-section">
          <h3 className="ds-h3">Active Jobs ({activeJobs.length})</h3>
          <div className="rp-jobs-list">
            {activeJobs.map(job => (
              <RecyclingJobCard key={job.id} job={job} onCancel={handleCancel} />
            ))}
          </div>
        </section>
      )}

      {/* Completed Jobs */}
      {completedJobs.length > 0 && (
        <section className="rp-section">
          <h3 className="ds-h3">Completed ({completedJobs.length})</h3>
          <div className="rp-jobs-list">
            {completedJobs.map(job => (
              <CompletedJobCard key={job.id} job={job} onCollect={handleCollect} />
            ))}
          </div>
        </section>
      )}

      {activeJobs.length === 0 && completedJobs.length === 0 && availableShips.length > 0 && (
        <div className="ds-card rp-empty">
          <p>No active recycling jobs.</p>
          <p className="ds-text-muted">Select a ship above to start recycling.</p>
        </div>
      )}

      {confirmRecycle && (
        <div
          className="ds-modal-backdrop"
          onClick={e => {
            if (e.target === e.currentTarget) setConfirmRecycle(false)
          }}
        >
          <div className="ds-modal ds-modal--sm" role="dialog" aria-modal="true">
            <div className="ds-modal-header">
              <div className="rp-confirm-title">
                <AlertTriangle size={20} strokeWidth={1.75} />
                <h3 className="ds-modal-title">Confirm Recycling</h3>
              </div>
              <button
                className="ds-btn-icon"
                onClick={() => setConfirmRecycle(false)}
                aria-label="Close"
              >
                <X size={18} strokeWidth={1.75} />
              </button>
            </div>
            <div className="ds-modal-body">
              <p>
                Are you sure you want to recycle this ship? The ship will be destroyed and cannot
                be recovered.
              </p>
              <p className="rp-modal-warning">
                <AlertTriangle size={14} strokeWidth={2} /> This action is permanent.
              </p>
            </div>
            <div className="ds-modal-footer">
              <button className="ds-btn ds-btn-ghost" onClick={() => setConfirmRecycle(false)}>
                Cancel
              </button>
              <LoadingButton
                className="ds-btn ds-btn-danger"
                onClick={handleRecycle}
                loading={loading}
              >
                Confirm Recycle
              </LoadingButton>
            </div>
          </div>
        </div>
      )}

      <style>{`
        .rp-panel { display: flex; flex-direction: column; gap: var(--sp-4); max-width: 960px; }
        .rp-header { display: flex; align-items: center; gap: var(--sp-3); }
        .rp-header-icon {
          width: 44px; height: 44px;
          display: grid; place-items: center;
          background: var(--ds-teal-tint);
          color: var(--ds-teal-dark);
          border-radius: var(--r-lg);
        }
        .rp-header-text h2 { margin: 0; }
        .rp-subtitle { margin: 2px 0 0; font-size: var(--fs-sm); }

        .rp-error { display: inline-flex; }

        .rp-section { display: flex; flex-direction: column; gap: var(--sp-3); }
        .rp-section h3 { margin: 0; }

        .rp-empty {
          color: var(--ds-text);
        }
        .rp-empty p { margin: 0 0 var(--sp-1); }
        .rp-empty p:last-child { font-size: var(--fs-sm); }

        .rp-ships-grid {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
          gap: var(--sp-2);
        }
        .rp-ship-card {
          display: flex; flex-direction: column; gap: var(--sp-1);
          padding: var(--sp-3);
          width: 100%;
          align-items: flex-start; text-align: left;
          border-color: transparent;
        }
        .rp-ship-name { font-weight: var(--fw-semibold); color: var(--ds-text); }
        .rp-ship-id { font-size: var(--fs-caption); }

        .rp-actions { display: flex; gap: var(--sp-2); margin-top: var(--sp-2); }

        .rp-jobs-list { display: flex; flex-direction: column; gap: var(--sp-2); }
        .rp-modal-warning {
          display: inline-flex; align-items: center; gap: var(--sp-1);
          color: var(--ds-danger); font-size: var(--fs-sm); margin-top: var(--sp-2);
        }
        .rp-confirm-title { display: flex; align-items: center; gap: var(--sp-2); color: var(--ds-text); }
        .rp-confirm-title h3 { margin: 0; }
      `}</style>
    </div>
  )
}

function RewardRow({ job }: { job: RecyclingJob }) {
  return (
    <div className="rp-rewards">
      <span className="rp-reward">
        <span className="ds-resource-dot ds-resource-dot--metal" />
        <span className="ds-mono">{formatNumber(job.metal_gained)}</span>
        <span className="ds-text-muted">M</span>
      </span>
      <span className="rp-reward">
        <span className="ds-resource-dot ds-resource-dot--he3" />
        <span className="ds-mono">{formatNumber(job.he3_gained)}</span>
        <span className="ds-text-muted">H3</span>
      </span>
      <span className="rp-reward">
        <span className="ds-resource-dot ds-resource-dot--gold" />
        <span className="ds-mono">{formatNumber(job.gold_gained)}</span>
        <span className="ds-text-muted">G</span>
      </span>
      <style>{`
        .rp-rewards { display: flex; flex-wrap: wrap; gap: var(--sp-3); font-size: var(--fs-sm); }
        .rp-reward { display: inline-flex; align-items: center; gap: var(--sp-1); }
      `}</style>
    </div>
  )
}

function RecyclingJobCard({
  job,
  onCancel,
}: {
  job: RecyclingJob
  onCancel: (id: string) => void
}) {
  const countdown = useCountdown(job.completed_at)
  const progress = useMemo(() => {
    const start = new Date(job.started_at).getTime()
    const end = new Date(job.completed_at).getTime()
    const now = Date.now()
    return Math.min(100, ((now - start) / (end - start)) * 100)
  }, [job.started_at, job.completed_at])

  return (
    <div className="ds-card rp-job-card">
      <div className="rp-job-header">
        <span className="rp-job-status">
          <Recycle size={14} strokeWidth={1.75} /> Recycling
        </span>
        <span className="ds-mono ds-text-muted">{countdown}</span>
      </div>
      <div className="ds-bar">
        <div className="ds-bar-fill" style={{ width: `${progress}%` }} />
      </div>
      <RewardRow job={job} />
      <button className="ds-btn ds-btn-ghost ds-btn--sm" onClick={() => onCancel(job.id)}>
        Cancel
      </button>

      <style>{`
        .rp-job-card { display: flex; flex-direction: column; gap: var(--sp-2); }
        .rp-job-header {
          display: flex; align-items: center; justify-content: space-between;
          font-size: var(--fs-sm);
        }
        .rp-job-status {
          display: inline-flex; align-items: center; gap: var(--sp-1);
          color: var(--ds-teal-dark); font-weight: var(--fw-semibold);
        }
      `}</style>
    </div>
  )
}

function CompletedJobCard({
  job,
  onCollect,
}: {
  job: RecyclingJob
  onCollect: (id: string) => void
}) {
  return (
    <div className="ds-card rp-job-card">
      <div className="rp-job-header">
        <span className="ds-badge ds-badge--success">
          <Coins size={12} strokeWidth={2} /> Complete
        </span>
      </div>
      <RewardRow job={job} />
      <LoadingButton
        className="ds-btn ds-btn-secondary ds-btn--sm"
        onClick={() => onCollect(job.id)}
      >
        Collect Resources
      </LoadingButton>

      <style>{`
        .rp-job-header {
          display: flex; align-items: center; justify-content: space-between;
          font-size: var(--fs-sm);
        }
      `}</style>
    </div>
  )
}
