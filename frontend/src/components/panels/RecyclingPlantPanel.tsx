import { useState, useMemo } from 'react'
import { useRecycling } from '../../hooks/useRecycling.ts'
import { formatNumber, useCountdown } from '../../hooks/useCountdown.ts'
import LoadingButton from '../common/LoadingButton.tsx'
import '../../styles/phase2.css'
import '../../styles/common.css'

export default function RecyclingPlantPanel() {
  const { jobs, availableShips, loading, recycle, collect, cancel } = useRecycling()
  const [selectedShipId, setSelectedShipId] = useState<string>('')
  const [confirmRecycle, setConfirmRecycle] = useState(false)

  // Separate active and completed jobs
  const { activeJobs, completedJobs } = useMemo(() => {
    const now = new Date()
    const active = jobs.filter((j) => new Date(j.completed_at) > now && !j.collected)
    const completed = jobs.filter((j) => new Date(j.completed_at) <= now && !j.collected)
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
    <div className="panel recycling-panel">
      <div className="panel-header">
        <h2>Recycling Plant</h2>
        <p className="panel-subtitle">Recycle unwanted ships for resources</p>
      </div>

      <div className="panel-content">
        {/* Available Ships Section */}
        <div className="rp-section">
          <h3>Available Ships for Recycling</h3>
          {availableShips.length === 0 ? (
            <div className="rp-empty">
              <p>No ships available for recycling</p>
              <p className="rp-hint">
                Ships must be unassigned from fleets to be recycled.
                Visit the Fleet panel to dismiss ships first.
              </p>
            </div>
          ) : (
            <div className="rp-ships-grid">
              {availableShips.map((ship) => (
                <div
                  key={ship.id}
                  className={`rp-ship-card ${selectedShipId === ship.id ? 'selected' : ''}`}
                  onClick={() => setSelectedShipId(ship.id)}
                >
                  <div className="rp-ship-name">{ship.design_name}</div>
                  <div className="rp-ship-class">{ship.hull_class}</div>
                  <div className="rp-ship-id">ID: {ship.id.slice(0, 8)}</div>
                </div>
              ))}
            </div>
          )}

          {selectedShipId && (
            <div className="rp-recycle-actions">
              <LoadingButton
                className="btn btn-danger"
                onClick={() => setConfirmRecycle(true)}
                loading={loading}
              >
                Recycle Selected Ship
              </LoadingButton>
              <button className="btn btn-secondary" onClick={() => setSelectedShipId('')}>
                Cancel
              </button>
            </div>
          )}
        </div>

        {/* Active Jobs Section */}
        {activeJobs.length > 0 && (
          <div className="rp-section">
            <h3>Active Recycling Jobs ({activeJobs.length})</h3>
            <div className="rp-jobs-list">
              {activeJobs.map((job) => (
                <RecyclingJobCard key={job.id} job={job} onCancel={handleCancel} />
              ))}
            </div>
          </div>
        )}

        {/* Completed Jobs Section */}
        {completedJobs.length > 0 && (
          <div className="rp-section">
            <h3>Completed Jobs ({completedJobs.length})</h3>
            <div className="rp-jobs-list">
              {completedJobs.map((job) => (
                <CompletedJobCard key={job.id} job={job} onCollect={handleCollect} />
              ))}
            </div>
          </div>
        )}

        {activeJobs.length === 0 && completedJobs.length === 0 && (
          <div className="rp-empty-jobs">
            <p>No active recycling jobs</p>
            <p className="rp-hint">Select a ship above to start recycling</p>
          </div>
        )}
      </div>

      {/* Confirmation Modal */}
      {confirmRecycle && (
        <div className="modal-overlay" onClick={() => setConfirmRecycle(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h3>Confirm Recycling</h3>
            <p>
              Are you sure you want to recycle this ship? The ship will be destroyed and cannot be
              recovered.
            </p>
            <p className="rp-modal-warning">⚠️ This action is permanent!</p>
            <div className="modal-actions">
              <LoadingButton className="btn btn-danger" onClick={handleRecycle} loading={loading}>
                Confirm Recycle
              </LoadingButton>
              <button className="btn btn-secondary" onClick={() => setConfirmRecycle(false)}>
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

function RecyclingJobCard({ job, onCancel }: { job: any; onCancel: (id: string) => void }) {
  const countdown = useCountdown(job.completed_at)
  const progress = useMemo(() => {
    const start = new Date(job.started_at).getTime()
    const end = new Date(job.completed_at).getTime()
    const now = Date.now()
    return Math.min(100, ((now - start) / (end - start)) * 100)
  }, [job.started_at, job.completed_at])

  return (
    <div className="rp-job-card active">
      <div className="rp-job-header">
        <span className="rp-job-status">🔄 Recycling</span>
        <span className="rp-job-time">{countdown}</span>
      </div>
      <div className="rp-job-progress">
        <div className="rp-progress-bar">
          <div className="rp-progress-fill" style={{ width: `${progress}%` }} />
        </div>
      </div>
      <div className="rp-job-rewards">
        <span>🪙 {formatNumber(job.metal_gained)} M</span>
        <span>⚡ {formatNumber(job.he3_gained)} H3</span>
        <span>💰 {formatNumber(job.gold_gained)} G</span>
      </div>
      <LoadingButton className="btn btn-small btn-secondary" onClick={() => onCancel(job.id)}>
        Cancel
      </LoadingButton>
    </div>
  )
}

function CompletedJobCard({ job, onCollect }: { job: any; onCollect: (id: string) => void }) {
  return (
    <div className="rp-job-card completed">
      <div className="rp-job-header">
        <span className="rp-job-status completed">✅ Complete</span>
      </div>
      <div className="rp-job-rewards">
        <span>🪙 {formatNumber(job.metal_gained)} M</span>
        <span>⚡ {formatNumber(job.he3_gained)} H3</span>
        <span>💰 {formatNumber(job.gold_gained)} G</span>
      </div>
      <LoadingButton className="btn btn-small btn-primary" onClick={() => onCollect(job.id)}>
        Collect Resources
      </LoadingButton>
    </div>
  )
}
