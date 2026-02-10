import { useState } from 'react'
import { useCommanders } from '../../hooks/useCommanders'
import { useResources } from '../../hooks/useResources'
import { recruitCommander, type RecruitResponse } from '../../services/api'
import LoadingButton from '../common/LoadingButton.tsx'
import './CommandCenterPanel.css'
import '../../styles/common.css'

export default function CommandCenterPanel() {
  const { refreshResources } = useResources()
  const { commanders, refresh: refreshCommanders } = useCommanders()
  const [recruiting, setRecruiting] = useState(false)
  const [result, setResult] = useState<RecruitResponse | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [cooldown, setCooldown] = useState(0) // TODO: Implement cooldown tracking

  const handleRecruit = async () => {
    if (!confirm('Recruit commander for 10,000 Gold?')) return

    setRecruiting(true)
    setResult(null)
    setError(null)

    try {
      const response = await recruitCommander()
      setResult(response)

      // Refresh commanders list and resources
      await refreshCommanders()
      await refreshResources()
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unknown error'
      setError(`Recruitment failed: ${message}`)
    } finally {
      setRecruiting(false)
    }
  }

  return (
    <div className="command-center-panel">
      <div className="panel-header">
        <h2>Command Center - Recruitment</h2>
        <div className="commander-count">
          {commanders.length} / 60 Commanders
        </div>
      </div>

      {error && <div className="p2-error-msg">{error}</div>}

      <div className="recruitment-section">
        <div className="recruitment-info">
          <div className="cost-display">
            <span className="label">Recruitment Cost:</span>
            <span className="value gold">10,000 Gold</span>
          </div>

          {cooldown > 0 && (
            <div className="cooldown-display">
              <span className="label">Cooldown:</span>
              <span className="value">{Math.floor(cooldown / 60)}m {cooldown % 60}s</span>
            </div>
          )}

          <div className="drop-rates">
            <h4>Drop Rates:</h4>
            <div className="rates-list">
              <div className="rate common">Common: 50%</div>
              <div className="rate skill">Skill: 35%</div>
              <div className="rate super">Super: 15%</div>
            </div>
          </div>
        </div>

        <LoadingButton
          className="recruit-btn"
          onClick={handleRecruit}
          loading={recruiting}
          disabled={cooldown > 0 || commanders.length >= 60}
        >
          {cooldown > 0 ? 'On Cooldown' : 'Recruit Commander'}
        </LoadingButton>
      </div>

      {result && (
        <div className={`recruitment-result ${result.is_duplicate ? 'duplicate' : 'new'}`}>
          {result.is_duplicate ? (
            <div className="duplicate-result">
              <div className="result-icon">📦</div>
              <h3>Duplicate Commander!</h3>
              <p className="message">{result.message}</p>
              <div className="hint">
                Go to <strong>Compound Center</strong> to merge duplicates and increase star rank!
              </div>
            </div>
          ) : (
            <div className="new-commander-result">
              <div className="result-icon">✨</div>
              <h3>New Commander Recruited!</h3>

              <div className="commander-card">
                <div className={`rarity-badge ${result.commander.rarity}`}>
                  {result.commander.rarity.toUpperCase()}
                </div>

                <h4 className="commander-name">{result.commander.name}</h4>

                <div className="star-rank">
                  {'★'.repeat(result.commander.star_rank)}
                  {'☆'.repeat(15 - result.commander.star_rank)}
                </div>

                <div className="stats-grid">
                  <div className="stat">
                    <span className="stat-label">ACC</span>
                    <span className="stat-value">{result.commander.accuracy}</span>
                  </div>
                  <div className="stat">
                    <span className="stat-label">DOD</span>
                    <span className="stat-value">{result.commander.dodge}</span>
                  </div>
                  <div className="stat">
                    <span className="stat-label">SPD</span>
                    <span className="stat-value">{result.commander.speed}</span>
                  </div>
                  <div className="stat">
                    <span className="stat-label">ELEC</span>
                    <span className="stat-value">{result.commander.electron}</span>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      )}

      {commanders.length === 0 && !result && (
        <div className="empty-state">
          <p>You have no commanders yet.</p>
          <p className="hint">Recruit your first commander to begin building your fleet!</p>
        </div>
      )}
    </div>
  )
}
