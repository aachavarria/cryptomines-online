import { useState } from 'react'
import { UserCircle2, X } from 'lucide-react'
import { useCommanders } from '../../hooks/useCommanders'
import { assignCommander, unassignCommander, type Commander } from '../../services/api'

interface FleetCommanderAssignmentProps {
  fleetId: string
  currentCommanderId?: string | null
  onAssigned?: () => void
  onUnassigned?: () => void
}

// Map backend rarity → DS rarity badge variant.
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

export default function FleetCommanderAssignment({
  fleetId,
  currentCommanderId,
  onAssigned,
  onUnassigned,
}: FleetCommanderAssignmentProps) {
  const { commanders, refresh: refreshCommanders } = useCommanders()
  const [showModal, setShowModal] = useState(false)
  const [assigning, setAssigning] = useState(false)

  const currentCommander = commanders.find(c => c.id === currentCommanderId)
  const availableCommanders = commanders.filter(
    c => !c.is_deployed || c.id === currentCommanderId,
  )

  const handleAssign = async (commander: Commander) => {
    if (!confirm(`Assign ${commander.name} to this fleet?`)) return

    setAssigning(true)
    try {
      await assignCommander(fleetId, commander.id)
      await refreshCommanders()
      setShowModal(false)
      onAssigned?.()
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Unknown error'
      alert(`Failed to assign commander: ${message}`)
    } finally {
      setAssigning(false)
    }
  }

  const handleUnassign = async () => {
    if (!currentCommander) return
    if (!confirm(`Unassign ${currentCommander.name} from this fleet?`)) return

    setAssigning(true)
    try {
      await unassignCommander(fleetId)
      await refreshCommanders()
      onUnassigned?.()
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Unknown error'
      alert(`Failed to unassign commander: ${message}`)
    } finally {
      setAssigning(false)
    }
  }

  return (
    <div style={{ marginTop: 'var(--sp-5)' }}>
      <div className="ds-card ds-stack">
        <div className="ds-caption">Commander</div>

        {currentCommander ? (
          <div className="ds-stack">
            <div className="ds-row" style={{ gap: 'var(--sp-3)', flexWrap: 'wrap' }}>
              <span style={{ fontWeight: 700, fontSize: 'var(--fs-h3)', color: 'var(--ds-text)' }}>
                {currentCommander.name}
              </span>
              <span className={`ds-badge ds-badge--rarity-${rarityVariant(currentCommander.rarity)}`}>
                {currentCommander.rarity.toUpperCase()}
              </span>
              <span
                className="ds-mono"
                style={{ color: 'var(--ds-orange)', letterSpacing: '0.1em' }}
                aria-label={`${currentCommander.star_rank} stars`}
              >
                {'★'.repeat(currentCommander.star_rank)}
              </span>
            </div>

            <div
              style={{
                display: 'grid',
                gridTemplateColumns: 'repeat(4, 1fr)',
                gap: 'var(--sp-3)',
                padding: 'var(--sp-3)',
                background: 'var(--ds-surface)',
                border: '1px solid var(--ds-border)',
                borderRadius: 'var(--r-md)',
              }}
            >
              {[
                ['ACC', currentCommander.accuracy],
                ['DOD', currentCommander.dodge],
                ['SPD', currentCommander.speed],
                ['ELEC', currentCommander.electron],
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
                      fontWeight: 700,
                      color: 'var(--ds-teal-dark)',
                    }}
                  >
                    {value}
                  </span>
                </div>
              ))}
            </div>

            <div className="ds-row" style={{ gap: 'var(--sp-2)' }}>
              <button
                type="button"
                className="ds-btn-ghost"
                style={{ flex: 1 }}
                onClick={() => setShowModal(true)}
                disabled={assigning}
              >
                Change
              </button>
              <button
                type="button"
                className="ds-btn-danger"
                style={{ flex: 1 }}
                onClick={handleUnassign}
                disabled={assigning}
              >
                Unassign
              </button>
            </div>
          </div>
        ) : (
          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              gap: 'var(--sp-3)',
              padding: 'var(--sp-6) var(--sp-4)',
            }}
          >
            <UserCircle2
              size={48}
              strokeWidth={1.5}
              color="var(--ds-text-soft)"
              aria-hidden="true"
            />
            <div className="ds-text-muted">No commander assigned</div>
            <button
              type="button"
              className="ds-btn-secondary"
              onClick={() => setShowModal(true)}
              disabled={assigning}
            >
              Assign Commander
            </button>
          </div>
        )}
      </div>

      {showModal && (
        <div
          className="ds-modal-backdrop"
          role="dialog"
          aria-modal="true"
          onClick={e => {
            if (e.target === e.currentTarget) setShowModal(false)
          }}
        >
          <div className="ds-modal">
            <div className="ds-modal-header">
              <h3 className="ds-modal-title">Select Commander</h3>
              <button
                type="button"
                className="ds-btn-icon"
                aria-label="Close"
                onClick={() => setShowModal(false)}
              >
                <X size={18} strokeWidth={2} />
              </button>
            </div>

            <div className="ds-modal-body">
              {availableCommanders.length === 0 ? (
                <div
                  style={{
                    textAlign: 'center',
                    padding: 'var(--sp-8) var(--sp-4)',
                    color: 'var(--ds-text-muted)',
                  }}
                >
                  <p style={{ margin: 0 }}>No commanders available.</p>
                  <p
                    className="ds-text-soft"
                    style={{ marginTop: 'var(--sp-2)', fontSize: 'var(--fs-sm)' }}
                  >
                    Visit the Command Center to recruit commanders.
                  </p>
                </div>
              ) : (
                <div
                  style={{
                    display: 'flex',
                    flexDirection: 'column',
                    gap: 'var(--sp-2)',
                  }}
                >
                  {availableCommanders.map(commander => {
                    const isCurrent = commander.id === currentCommanderId
                    return (
                      <div
                        key={commander.id}
                        role="button"
                        tabIndex={0}
                        aria-selected={isCurrent}
                        className="ds-list-item"
                        onClick={() => handleAssign(commander)}
                        onKeyDown={(e) => {
                          if (e.key === 'Enter' || e.key === ' ') {
                            e.preventDefault()
                            handleAssign(commander)
                          }
                        }}
                      >
                        <div className="ds-row--between" style={{ marginBottom: 6 }}>
                          <span style={{ fontWeight: 700, color: 'var(--ds-text)' }}>
                            {commander.name}
                          </span>
                          {isCurrent && (
                            <span className="ds-badge ds-badge--orange">Current</span>
                          )}
                        </div>

                        <div
                          className="ds-row"
                          style={{ gap: 'var(--sp-2)', marginBottom: 6 }}
                        >
                          <span className={`ds-badge ds-badge--rarity-${rarityVariant(commander.rarity)}`}>
                            {commander.rarity.toUpperCase()}
                          </span>
                          <span
                            className="ds-mono"
                            style={{
                              color: 'var(--ds-orange)',
                              fontSize: 'var(--fs-sm)',
                              letterSpacing: '0.1em',
                            }}
                          >
                            {'★'.repeat(commander.star_rank)}
                          </span>
                        </div>

                        <div
                          className="ds-row"
                          style={{
                            gap: 'var(--sp-3)',
                            flexWrap: 'wrap',
                            fontSize: 'var(--fs-sm)',
                            color: 'var(--ds-text-muted)',
                          }}
                        >
                          <span>ACC <span className="ds-mono">{commander.accuracy}</span></span>
                          <span>DOD <span className="ds-mono">{commander.dodge}</span></span>
                          <span>SPD <span className="ds-mono">{commander.speed}</span></span>
                          <span>ELEC <span className="ds-mono">{commander.electron}</span></span>
                        </div>
                      </div>
                    )
                  })}
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
