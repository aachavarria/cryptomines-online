import { useState } from 'react'
import { useCommanders } from '../../hooks/useCommanders'
import { assignCommander, unassignCommander, type Commander } from '../../api/commanders'
import './FleetCommanderAssignment.css'

interface FleetCommanderAssignmentProps {
  fleetId: string
  currentCommanderId?: string | null
  onAssigned?: () => void
  onUnassigned?: () => void
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

  // Get current commander
  const currentCommander = commanders.find(c => c.id === currentCommanderId)

  // Available commanders (not deployed)
  const availableCommanders = commanders.filter(c => !c.is_deployed || c.id === currentCommanderId)

  const handleAssign = async (commander: Commander) => {
    if (!confirm(`Assign ${commander.name} to this fleet?`)) return

    setAssigning(true)
    try {
      await assignCommander(fleetId, commander.id)
      alert(`✓ ${commander.name} assigned to fleet`)

      await refreshCommanders()
      setShowModal(false)
      onAssigned?.()
    } catch (err: any) {
      alert(`✗ Failed to assign commander: ${err.message}`)
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
      alert(`✓ ${currentCommander.name} unassigned`)

      await refreshCommanders()
      onUnassigned?.()
    } catch (err: any) {
      alert(`✗ Failed to unassign commander: ${err.message}`)
    } finally {
      setAssigning(false)
    }
  }

  return (
    <div className="fleet-commander-assignment">
      <div className="commander-slot">
        <div className="slot-label">Commander:</div>

        {currentCommander ? (
          <div className="assigned-commander">
            <div className="commander-info">
              <div className="commander-name">{currentCommander.name}</div>
              <div className={`rarity-badge ${currentCommander.rarity}`}>
                {currentCommander.rarity.toUpperCase()}
              </div>
              <div className="star-rank">
                {'★'.repeat(currentCommander.star_rank)}
              </div>
            </div>

            <div className="commander-stats">
              <div className="stat">
                <span className="label">ACC</span>
                <span className="value">{currentCommander.accuracy}</span>
              </div>
              <div className="stat">
                <span className="label">DOD</span>
                <span className="value">{currentCommander.dodge}</span>
              </div>
              <div className="stat">
                <span className="label">SPD</span>
                <span className="value">{currentCommander.speed}</span>
              </div>
              <div className="stat">
                <span className="label">ELEC</span>
                <span className="value">{currentCommander.electron}</span>
              </div>
            </div>

            <div className="commander-actions">
              <button
                className="change-btn"
                onClick={() => setShowModal(true)}
                disabled={assigning}
              >
                Change
              </button>
              <button
                className="unassign-btn"
                onClick={handleUnassign}
                disabled={assigning}
              >
                Unassign
              </button>
            </div>
          </div>
        ) : (
          <div className="no-commander">
            <div className="empty-icon">👤</div>
            <div className="empty-text">No commander assigned</div>
            <button
              className="assign-btn"
              onClick={() => setShowModal(true)}
              disabled={assigning}
            >
              Assign Commander
            </button>
          </div>
        )}
      </div>

      {/* Commander Selection Modal */}
      {showModal && (
        <div className="commander-modal-backdrop" onClick={() => setShowModal(false)}>
          <div className="commander-modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Select Commander</h3>
              <button className="close-btn" onClick={() => setShowModal(false)}>×</button>
            </div>

            <div className="modal-body">
              {availableCommanders.length === 0 ? (
                <div className="no-commanders">
                  <p>No commanders available.</p>
                  <p className="hint">Visit the Command Center to recruit commanders.</p>
                </div>
              ) : (
                <div className="commanders-list">
                  {availableCommanders.map(commander => (
                    <div
                      key={commander.id}
                      className={`commander-option ${commander.rarity} ${commander.id === currentCommanderId ? 'current' : ''}`}
                      onClick={() => handleAssign(commander)}
                    >
                      <div className="option-header">
                        <div className="name">{commander.name}</div>
                        {commander.id === currentCommanderId && (
                          <div className="current-badge">CURRENT</div>
                        )}
                      </div>

                      <div className="option-details">
                        <div className={`rarity-badge ${commander.rarity}`}>
                          {commander.rarity.toUpperCase()}
                        </div>
                        <div className="star-rank">
                          {'★'.repeat(commander.star_rank)}
                        </div>
                      </div>

                      <div className="option-stats">
                        <span>ACC: {commander.accuracy}</span>
                        <span>DOD: {commander.dodge}</span>
                        <span>SPD: {commander.speed}</span>
                        <span>ELEC: {commander.electron}</span>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
