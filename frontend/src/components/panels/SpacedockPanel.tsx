import { Anchor, ParkingCircle, Wrench } from 'lucide-react'
import { useSpacedock } from '../../hooks/useSpacedock.ts'
import { useCountdown } from '../../hooks/useCountdown.ts'
import type { SpacedockRepair } from '../../types'

function RepairTimer({ finishAt }: { finishAt: string | null }) {
  const timeLeft = useCountdown(finishAt)
  if (!timeLeft) return <span className="ds-badge ds-badge--success">Complete</span>
  return <span className="ds-mono ds-text-muted">{timeLeft}</span>
}

export default function SpacedockPanel() {
  const { status, repairs, loading, error, repair } = useSpacedock()

  if (loading) {
    return (
      <div className="ds-panel sd-loading">
        <div className="loading-spinner" />
        <span>Loading Spacedock...</span>
      </div>
    )
  }

  if (error || !status) {
    return (
      <div className="ds-panel">
        <p className="ds-text-muted">{error || 'Spacedock not found. Build one first.'}</p>
      </div>
    )
  }

  const pendingRepairs = repairs.filter(
    r => r.repair_finish_at && new Date(r.repair_finish_at) > new Date(),
  )
  const completedRepairs = repairs.filter(
    r => !r.repair_finish_at || new Date(r.repair_finish_at) <= new Date(),
  )

  return (
    <div className="ds-panel sd-panel">
      <header className="sd-header">
        <div className="sd-header-icon">
          <Anchor size={22} strokeWidth={1.75} />
        </div>
        <div className="sd-header-text">
          <h2 className="ds-h2">Spacedock</h2>
          <p className="ds-text-muted sd-subtitle">
            Level <span className="ds-mono">{status.level}</span> · Repair rate{' '}
            <span className="ds-mono">{status.repair_pct}%</span>
          </p>
        </div>
      </header>

      <section className="sd-section">
        <h3 className="ds-h3">Active Repairs ({pendingRepairs.length})</h3>
        {pendingRepairs.length === 0 ? (
          <div className="ds-card sd-empty">
            <ParkingCircle size={28} strokeWidth={1.5} />
            <p>No ships being repaired.</p>
          </div>
        ) : (
          <div className="sd-list">
            {pendingRepairs.map(r => (
              <RepairCard key={r.id} repair={r} />
            ))}
          </div>
        )}
      </section>

      <section className="sd-section">
        <button
          className="ds-btn ds-btn-secondary ds-btn--block"
          onClick={repair}
          disabled={pendingRepairs.length > 0}
        >
          <Wrench size={16} strokeWidth={2} />
          Repair All Damaged Ships
        </button>
        <p className="ds-text-muted sd-info">
          Ships destroyed in PvP battles can be recovered here. Repair rate:{' '}
          <span className="ds-mono">{status.repair_pct}%</span> chance per ship.
        </p>
      </section>

      {completedRepairs.length > 0 && (
        <section className="sd-section">
          <h3 className="ds-h3">Completed</h3>
          <div className="sd-list">
            {completedRepairs.map(r => (
              <div key={r.id} className="ds-card sd-repair-card sd-repair--done">
                <div className="sd-repair-name ds-mono">
                  Design {r.ship_design_id.slice(0, 8)}…
                </div>
                <div className="sd-repair-stats">
                  <span>
                    Destroyed: <span className="ds-mono">{r.destroyed_count}</span>
                  </span>
                  <span className="sd-recovered">
                    Recovered: <span className="ds-mono">{r.repaired_count}</span>
                  </span>
                  <span className="sd-lost">
                    Lost:{' '}
                    <span className="ds-mono">{r.destroyed_count - r.repaired_count}</span>
                  </span>
                </div>
              </div>
            ))}
          </div>
        </section>
      )}

      <style>{`
        .sd-panel { display: flex; flex-direction: column; gap: var(--sp-4); max-width: 720px; }
        .sd-loading { display: flex; align-items: center; gap: var(--sp-3); }
        .sd-header { display: flex; align-items: center; gap: var(--sp-3); }
        .sd-header-icon {
          width: 44px; height: 44px;
          display: grid; place-items: center;
          background: var(--ds-teal-tint);
          color: var(--ds-teal-dark);
          border-radius: var(--r-lg);
        }
        .sd-header-text h2 { margin: 0; }
        .sd-subtitle { margin: 2px 0 0; font-size: var(--fs-sm); }
        .sd-section { display: flex; flex-direction: column; gap: var(--sp-3); }
        .sd-section h3 { margin: 0; }
        .sd-empty {
          display: flex; flex-direction: column; align-items: center; gap: var(--sp-2);
          color: var(--ds-text-muted); padding: var(--sp-5);
        }
        .sd-empty p { margin: 0; }
        .sd-list { display: flex; flex-direction: column; gap: var(--sp-2); }
        .sd-info { margin-top: var(--sp-2); font-size: var(--fs-sm); }
        .sd-repair-name { color: var(--ds-text); font-weight: var(--fw-semibold); }
        .sd-repair-stats {
          display: flex; flex-wrap: wrap; gap: var(--sp-3);
          font-size: var(--fs-sm); color: var(--ds-text);
        }
        .sd-recovered { color: var(--ds-success); }
        .sd-lost { color: var(--ds-danger); }
      `}</style>
    </div>
  )
}

function RepairCard({ repair }: { repair: SpacedockRepair }) {
  return (
    <div className="ds-card sd-repair-card">
      <div className="sd-repair-head">
        <span className="ds-mono sd-repair-name">
          Design {repair.ship_design_id.slice(0, 8)}…
        </span>
        <RepairTimer finishAt={repair.repair_finish_at} />
      </div>
      <div className="sd-repair-stats">
        <span>
          Destroyed: <span className="ds-mono">{repair.destroyed_count}</span>
        </span>
        <span>
          Expected recovery:{' '}
          <span className="ds-mono">~{repair.repaired_count}</span>
        </span>
      </div>
      <div className="ds-bar">
        <div className="ds-bar-fill" style={{ width: '50%' }} />
      </div>

      <style>{`
        .sd-repair-card { display: flex; flex-direction: column; gap: var(--sp-2); }
        .sd-repair-head {
          display: flex; align-items: center; justify-content: space-between;
        }
        .sd-repair-name { color: var(--ds-text); font-weight: var(--fw-semibold); }
      `}</style>
    </div>
  )
}
