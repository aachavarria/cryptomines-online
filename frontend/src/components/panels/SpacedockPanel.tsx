import { useSpacedock } from '../../hooks/useSpacedock.ts'
import { useCountdown, formatNumber } from '../../hooks/useCountdown.ts'
import type { SpacedockRepair } from '../../types'

function RepairTimer({ finishAt }: { finishAt: string | null }) {
  const timeLeft = useCountdown(finishAt)
  if (!timeLeft) return <span className="dock-repair-done">Complete</span>
  return <span className="dock-repair-timer">{timeLeft}</span>
}

export default function SpacedockPanel() {
  const { status, repairs, loading, error, repair } = useSpacedock()

  if (loading) {
    return <div className="p2-panel-loading"><div className="loading-spinner" /><span>Loading Spacedock...</span></div>
  }

  if (error || !status) {
    return <div className="p2-panel-error">{error || 'Spacedock not found. Build one first.'}</div>
  }

  const pendingRepairs = repairs.filter(r => r.repair_finish_at && new Date(r.repair_finish_at) > new Date())
  const completedRepairs = repairs.filter(r => !r.repair_finish_at || new Date(r.repair_finish_at) <= new Date())

  return (
    <div className="p2-panel">
      <div className="p2-panel-header">
        <div className="p2-panel-icon dock-icon">SD</div>
        <div>
          <div className="p2-panel-title">Spacedock</div>
          <div className="p2-panel-subtitle">
            Level {status.level} | Repair Rate: {status.repair_pct}%
          </div>
        </div>
      </div>

      {/* Active Repairs */}
      <div className="p2-section">
        <div className="p2-section-title">
          Active Repairs ({pendingRepairs.length})
        </div>
        {pendingRepairs.length === 0 ? (
          <div className="p2-empty-state">No ships being repaired.</div>
        ) : (
          <div className="dock-repair-list">
            {pendingRepairs.map(r => (
              <RepairCard key={r.id} repair={r} />
            ))}
          </div>
        )}
      </div>

      {/* Repair Button */}
      <div className="p2-section">
        <button
          className="p2-btn p2-btn-primary p2-btn-full"
          onClick={repair}
          disabled={pendingRepairs.length > 0}
        >
          Repair All Damaged Ships
        </button>
        <div className="dock-info-text">
          Ships destroyed in PvP battles can be recovered here.
          Repair rate: {status.repair_pct}% chance per ship.
        </div>
      </div>

      {/* Completed Repairs */}
      {completedRepairs.length > 0 && (
        <div className="p2-section">
          <div className="p2-section-title">Completed</div>
          <div className="dock-repair-list">
            {completedRepairs.map(r => (
              <div key={r.id} className="dock-repair-card dock-repair-done-card">
                <div className="dock-repair-name">Design {r.ship_design_id.slice(0, 8)}...</div>
                <div className="dock-repair-stats">
                  <span>Destroyed: {r.destroyed_count}</span>
                  <span className="dock-recovered">Recovered: {r.repaired_count}</span>
                  <span className="dock-lost">Lost: {r.destroyed_count - r.repaired_count}</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

function RepairCard({ repair }: { repair: SpacedockRepair }) {
  return (
    <div className="dock-repair-card">
      <div className="dock-repair-header">
        <span className="dock-repair-name">Design {repair.ship_design_id.slice(0, 8)}...</span>
        <RepairTimer finishAt={repair.repair_finish_at} />
      </div>
      <div className="dock-repair-stats">
        <span>Destroyed: {repair.destroyed_count}</span>
        <span>Expected recovery: ~{repair.repaired_count}</span>
      </div>
      <div className="dock-repair-bar">
        <div className="dock-repair-bar-fill" />
      </div>
    </div>
  )
}
