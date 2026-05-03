import { useState, useEffect, useRef } from 'react'
import { createPortal } from 'react-dom'
import { X, Sparkles, Package, Users } from 'lucide-react'
import { useCommanders } from '../../hooks/useCommanders'
import { useResources } from '../../hooks/useResources'
import { useGameContext } from '../../contexts/GameContext.tsx'
import { recruitCommander, type RecruitResponse } from '../../services/api'
import LoadingButton from '../common/LoadingButton.tsx'
import axios from 'axios'

// Map backend rarity → design-system rarity badge variant.
// Backend uses 'common' / 'skill' / 'super'. DS expects
// 'common' | 'rare' | 'epic' | 'legendary' (per spec).
function rarityVariant(rarity: string): 'common' | 'rare' | 'epic' | 'legendary' {
  switch (rarity) {
    case 'super':
      return 'epic'
    case 'skill':
      return 'rare'
    case 'legendary':
      return 'legendary'
    default:
      return 'common'
  }
}

export default function CommandCenterPanel() {
  const { state, closeCommandCenterPanel } = useGameContext()
  const { refreshResources } = useResources()
  const { commanders, refresh: refreshCommanders } = useCommanders()
  const [recruiting, setRecruiting] = useState(false)
  const [result, setResult] = useState<RecruitResponse | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [cooldown, setCooldown] = useState(0)
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null)

  const isOpen = state.showCommandCenterPanel

  useEffect(() => {
    if (!isOpen) return
    function onKey(e: KeyboardEvent) {
      if (e.key === 'Escape') closeCommandCenterPanel()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [isOpen, closeCommandCenterPanel])

  // Countdown timer
  useEffect(() => {
    if (cooldown <= 0) {
      if (timerRef.current) clearInterval(timerRef.current)
      return
    }
    timerRef.current = setInterval(() => {
      setCooldown(prev => {
        if (prev <= 1) {
          if (timerRef.current) clearInterval(timerRef.current)
          return 0
        }
        return prev - 1
      })
    }, 1000)
    return () => { if (timerRef.current) clearInterval(timerRef.current) }
  }, [cooldown > 0]) // eslint-disable-line react-hooks/exhaustive-deps

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
      if (axios.isAxiosError(err) && err.response?.data?.remaining_seconds) {
        setCooldown(err.response.data.remaining_seconds)
      }
      const message = err instanceof Error ? err.message : 'Unknown error'
      setError(`Recruitment failed: ${message}`)
    } finally {
      setRecruiting(false)
    }
  }

  if (!isOpen) return null

  return createPortal(
    <div
      className="ds-modal-backdrop"
      role="dialog"
      aria-modal="true"
      onClick={e => { if (e.target === e.currentTarget) closeCommandCenterPanel() }}
    >
      <div className="ds-modal ds-modal--lg" style={{ maxWidth: 'min(680px, 95vw)' }}>
        <div className="ds-modal-header">
          <div className="ds-row" style={{ gap: 'var(--sp-3)' }}>
            <h2 className="ds-modal-title">Command Center</h2>
            <span className="ds-badge ds-badge--neutral">
              <Users size={12} strokeWidth={2} />
              <span className="ds-mono">{commanders.length}</span>
              {' / '}
              <span className="ds-mono">60</span>
            </span>
          </div>
          <button
            type="button"
            className="ds-btn-icon"
            aria-label="Close"
            onClick={closeCommandCenterPanel}
          >
            <X size={18} strokeWidth={2} />
          </button>
        </div>

        <div className="ds-modal-body">
          {error && (
            <div
              className="ds-badge ds-badge--danger"
              style={{
                display: 'flex',
                width: '100%',
                padding: 'var(--sp-3)',
                marginBottom: 'var(--sp-4)',
              }}
            >
              {error}
            </div>
          )}

          <div className="ds-card ds-stack" style={{ marginBottom: 'var(--sp-4)' }}>
            <div className="ds-row--between">
              <span className="ds-text-muted" style={{ fontSize: 'var(--fs-sm)' }}>
                Recruitment cost
              </span>
              <span className="ds-row" style={{ gap: 6 }}>
                <span className="ds-resource-dot ds-resource-dot--gold" />
                <span className="ds-mono" style={{ fontWeight: 600 }}>10,000</span>
                <span className="ds-text-muted" style={{ fontSize: 'var(--fs-sm)' }}>Gold</span>
              </span>
            </div>

            {cooldown > 0 && (
              <div className="ds-row--between">
                <span className="ds-text-muted" style={{ fontSize: 'var(--fs-sm)' }}>
                  Cooldown
                </span>
                <span className="ds-mono" style={{ color: 'var(--ds-danger)', fontWeight: 600 }}>
                  {cooldown >= 3600 && `${Math.floor(cooldown / 3600)}h `}
                  {Math.floor((cooldown % 3600) / 60)}m {cooldown % 60}s
                </span>
              </div>
            )}

            <div>
              <div className="ds-caption" style={{ marginBottom: 'var(--sp-2)' }}>
                Drop rates
              </div>
              <div className="ds-row" style={{ gap: 'var(--sp-2)' }}>
                <span className="ds-badge ds-badge--rarity-common" style={{ flex: 1, justifyContent: 'center' }}>
                  Common <span className="ds-mono" style={{ marginLeft: 4 }}>50%</span>
                </span>
                <span className="ds-badge ds-badge--rarity-rare" style={{ flex: 1, justifyContent: 'center' }}>
                  Skill <span className="ds-mono" style={{ marginLeft: 4 }}>35%</span>
                </span>
                <span className="ds-badge ds-badge--rarity-epic" style={{ flex: 1, justifyContent: 'center' }}>
                  Super <span className="ds-mono" style={{ marginLeft: 4 }}>15%</span>
                </span>
              </div>
            </div>
          </div>

          {result && (
            <div className="ds-card ds-stack" style={{ marginBottom: 'var(--sp-4)' }}>
              {result.is_duplicate ? (
                <>
                  <div className="ds-row" style={{ gap: 'var(--sp-2)' }}>
                    <Package size={20} strokeWidth={1.75} color="var(--ds-orange-strong)" />
                    <h3 className="ds-h3" style={{ margin: 0 }}>Duplicate Commander</h3>
                  </div>
                  <p style={{ margin: 0, color: 'var(--ds-text)' }}>{result.message}</p>
                  <div
                    className="ds-text-muted"
                    style={{
                      fontSize: 'var(--fs-sm)',
                      padding: 'var(--sp-3)',
                      background: 'var(--ds-orange-tint)',
                      borderRadius: 'var(--r-md)',
                    }}
                  >
                    Visit the <strong>Compound Center</strong> to merge duplicates and increase star rank.
                  </div>
                </>
              ) : (
                <>
                  <div className="ds-row" style={{ gap: 'var(--sp-2)' }}>
                    <Sparkles size={20} strokeWidth={1.75} color="var(--ds-teal)" />
                    <h3 className="ds-h3" style={{ margin: 0 }}>New Commander Recruited</h3>
                  </div>

                  <div
                    className="ds-card"
                    style={{
                      background: 'var(--ds-surface)',
                      textAlign: 'center',
                    }}
                  >
                    <div style={{ marginBottom: 'var(--sp-2)' }}>
                      <span className={`ds-badge ds-badge--rarity-${rarityVariant(result.commander.rarity)}`}>
                        {result.commander.rarity.toUpperCase()}
                      </span>
                    </div>

                    <div
                      className="ds-h3"
                      style={{ margin: '0 0 var(--sp-2)' }}
                    >
                      {result.commander.name}
                    </div>

                    <div
                      className="ds-mono"
                      style={{
                        color: 'var(--ds-orange)',
                        fontSize: 'var(--fs-h3)',
                        letterSpacing: '0.15em',
                        marginBottom: 'var(--sp-3)',
                      }}
                    >
                      {'★'.repeat(result.commander.star_rank)}
                      <span style={{ color: 'var(--ds-text-soft)' }}>
                        {'☆'.repeat(15 - result.commander.star_rank)}
                      </span>
                    </div>

                    <div
                      style={{
                        display: 'grid',
                        gridTemplateColumns: 'repeat(4, 1fr)',
                        gap: 'var(--sp-3)',
                      }}
                    >
                      {[
                        ['ACC', result.commander.accuracy],
                        ['DOD', result.commander.dodge],
                        ['SPD', result.commander.speed],
                        ['ELEC', result.commander.electron],
                      ].map(([label, value]) => (
                        <div
                          key={label as string}
                          style={{
                            display: 'flex',
                            flexDirection: 'column',
                            alignItems: 'center',
                            gap: 4,
                          }}
                        >
                          <span className="ds-caption">{label}</span>
                          <span
                            className="ds-mono"
                            style={{
                              fontSize: 'var(--fs-h3)',
                              fontWeight: 600,
                              color: 'var(--ds-teal-dark)',
                            }}
                          >
                            {value}
                          </span>
                        </div>
                      ))}
                    </div>
                  </div>
                </>
              )}
            </div>
          )}

          {commanders.length === 0 && !result && (
            <div
              style={{
                textAlign: 'center',
                padding: 'var(--sp-6) var(--sp-4)',
                color: 'var(--ds-text-muted)',
              }}
            >
              <p style={{ margin: 0 }}>You have no commanders yet.</p>
              <p
                className="ds-text-soft"
                style={{ marginTop: 'var(--sp-2)', fontSize: 'var(--fs-sm)' }}
              >
                Recruit your first commander to begin building your fleet.
              </p>
            </div>
          )}
        </div>

        <div className="ds-modal-footer">
          <button
            type="button"
            className="ds-btn-ghost"
            onClick={closeCommandCenterPanel}
          >
            Close
          </button>
          <LoadingButton
            className="ds-btn-primary"
            onClick={handleRecruit}
            loading={recruiting}
            disabled={cooldown > 0 || commanders.length >= 60}
          >
            <Sparkles size={14} strokeWidth={2} />
            {cooldown > 0 ? 'On Cooldown' : 'Recruit Commander'}
          </LoadingButton>
        </div>
      </div>
    </div>,
    document.body
  )
}
