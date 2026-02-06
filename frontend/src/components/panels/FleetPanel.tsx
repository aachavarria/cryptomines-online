import { useState } from 'react'
import { useFleets } from '../../hooks/useFleets.ts'
import { useShipDesigns } from '../../hooks/useShipDesigns.ts'
import type { Fleet, FleetStack, ShipDesign } from '../../types'
import { formatNumber } from '../../hooks/useCountdown.ts'

const FORMATIONS = ['Phalanx', 'Diamond', 'Arrow', 'Defensive', 'Spread']
const TARGETING = ['Weakest First', 'Strongest First', 'Random', 'Closest']
const GRID_LABELS = [
  ['Front-L', 'Front-C', 'Front-R'],
  ['Mid-L', 'Mid-C', 'Mid-R'],
  ['Back-L', 'Back-C', 'Back-R'],
]

interface FleetEditorProps {
  fleet: Fleet
  designs: ShipDesign[]
  onAssign: (fleetId: string, designId: string, row: number, col: number, count: number) => Promise<void>
  onRemove: (fleetId: string, row: number, col: number) => Promise<void>
  onUpdate: (id: string, updates: { formation?: string; targeting_command?: string }) => Promise<void>
  onClose: () => void
}

function FleetEditor({ fleet, designs, onAssign, onRemove, onUpdate, onClose }: FleetEditorProps) {
  const [assignTarget, setAssignTarget] = useState<{ row: number; col: number } | null>(null)
  const [selectedDesign, setSelectedDesign] = useState('')
  const [shipCount, setShipCount] = useState(100)

  function getStack(row: number, col: number): FleetStack | undefined {
    return fleet.stacks?.find(s => s.grid_row === row && s.grid_col === col)
  }

  async function handleAssign() {
    if (!assignTarget || !selectedDesign || shipCount < 1) return
    await onAssign(fleet.id, selectedDesign, assignTarget.row, assignTarget.col, shipCount)
    setAssignTarget(null)
    setSelectedDesign('')
    setShipCount(100)
  }

  return (
    <div className="p2-modal-backdrop" onClick={e => { if (e.target === e.currentTarget) onClose() }}>
      <div className="p2-modal p2-modal-lg">
        <div className="p2-modal-header">
          <span className="p2-modal-title">Fleet: {fleet.name}</span>
          <button className="p2-modal-close" onClick={onClose}>X</button>
        </div>
        <div className="p2-modal-body">
          {/* Formation & Targeting */}
          <div className="fleet-config">
            <div className="p2-form-group">
              <label className="p2-label">Formation</label>
              <select
                className="p2-select"
                value={fleet.formation}
                onChange={e => onUpdate(fleet.id, { formation: e.target.value })}
              >
                {FORMATIONS.map(f => <option key={f} value={f}>{f}</option>)}
              </select>
            </div>
            <div className="p2-form-group">
              <label className="p2-label">Targeting</label>
              <select
                className="p2-select"
                value={fleet.targeting_command}
                onChange={e => onUpdate(fleet.id, { targeting_command: e.target.value })}
              >
                {TARGETING.map(t => <option key={t} value={t}>{t}</option>)}
              </select>
            </div>
          </div>

          {/* 3x3 Grid */}
          <div className="fleet-grid">
            {[0, 1, 2].map(row => (
              <div key={row} className="fleet-grid-row">
                {[0, 1, 2].map(col => {
                  const stack = getStack(row, col)
                  const isAssignTarget = assignTarget?.row === row && assignTarget?.col === col
                  return (
                    <div
                      key={col}
                      className={`fleet-cell ${stack ? 'occupied' : 'empty'} ${isAssignTarget ? 'target' : ''}`}
                      onClick={() => {
                        if (!stack) setAssignTarget({ row, col })
                      }}
                    >
                      <div className="fleet-cell-label">{GRID_LABELS[row][col]}</div>
                      {stack ? (
                        <>
                          <div className="fleet-cell-design">{stack.ship_design_name || 'Ships'}</div>
                          <div className="fleet-cell-count">x{formatNumber(stack.ship_count)}</div>
                          <button
                            className="fleet-cell-remove"
                            onClick={e => { e.stopPropagation(); onRemove(fleet.id, row, col) }}
                          >
                            X
                          </button>
                        </>
                      ) : (
                        <div className="fleet-cell-empty">Empty</div>
                      )}
                    </div>
                  )
                })}
              </div>
            ))}
          </div>

          {/* Assign Stack Form */}
          {assignTarget && (
            <div className="fleet-assign-form">
              <div className="p2-section-title">
                Assign Ships to {GRID_LABELS[assignTarget.row][assignTarget.col]}
              </div>
              <div className="fleet-assign-row">
                <select
                  className="p2-select"
                  value={selectedDesign}
                  onChange={e => setSelectedDesign(e.target.value)}
                >
                  <option value="">-- Select Design --</option>
                  {designs.map(d => (
                    <option key={d.id} value={d.id}>
                      {d.name} ({d.ships_built} available)
                    </option>
                  ))}
                </select>
                <input
                  type="number"
                  className="p2-input p2-input-sm"
                  value={shipCount}
                  min={1}
                  max={3000}
                  onChange={e => setShipCount(Math.max(1, Math.min(3000, parseInt(e.target.value) || 1)))}
                />
                <button
                  className="p2-btn p2-btn-primary p2-btn-sm"
                  onClick={handleAssign}
                  disabled={!selectedDesign || shipCount < 1}
                >
                  Assign
                </button>
                <button
                  className="p2-btn p2-btn-secondary p2-btn-sm"
                  onClick={() => setAssignTarget(null)}
                >
                  Cancel
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default function FleetPanel() {
  const { fleets, loading, error, create, update, remove, addStack, removeStackFromFleet } = useFleets()
  const { designs } = useShipDesigns()
  const [editFleet, setEditFleet] = useState<Fleet | null>(null)
  const [creating, setCreating] = useState(false)
  const [newName, setNewName] = useState('')

  if (loading) {
    return <div className="p2-panel-loading"><div className="loading-spinner" /><span>Loading Fleets...</span></div>
  }

  async function handleCreate() {
    if (!newName) return
    setCreating(true)
    try {
      const fleet = await create({ name: newName, formation: 'Phalanx', targeting_command: 'Weakest First' })
      setNewName('')
      setEditFleet(fleet)
    } finally {
      setCreating(false)
    }
  }

  async function handleAssign(fleetId: string, designId: string, row: number, col: number, count: number) {
    const result = await addStack(fleetId, {
      ship_design_id: designId,
      grid_row: row,
      grid_col: col,
      ship_count: count,
    })
    // Refresh the edit fleet data
    const updated = fleets.find(f => f.id === fleetId)
    if (updated) setEditFleet({ ...updated })
  }

  async function handleRemoveStack(fleetId: string, row: number, col: number) {
    await removeStackFromFleet(fleetId, row, col)
    const updated = fleets.find(f => f.id === fleetId)
    if (updated) setEditFleet({ ...updated })
  }

  return (
    <div className="p2-panel">
      <div className="p2-panel-header">
        <div className="p2-panel-icon fleet-icon">FL</div>
        <div>
          <div className="p2-panel-title">Fleet Management</div>
          <div className="p2-panel-subtitle">{fleets.length} fleet{fleets.length !== 1 ? 's' : ''}</div>
        </div>
      </div>

      {error && <div className="p2-error-msg">{error}</div>}

      {/* Create Fleet */}
      <div className="p2-section">
        <div className="fleet-create-row">
          <input
            className="p2-input"
            value={newName}
            placeholder="New fleet name..."
            onChange={e => setNewName(e.target.value)}
            onKeyDown={e => { if (e.key === 'Enter') handleCreate() }}
          />
          <button
            className="p2-btn p2-btn-primary"
            onClick={handleCreate}
            disabled={!newName || creating}
          >
            {creating ? '...' : 'Create Fleet'}
          </button>
        </div>
      </div>

      {/* Fleet List */}
      <div className="p2-section">
        {fleets.length === 0 ? (
          <div className="p2-empty-state">No fleets yet. Create a fleet and assign ships.</div>
        ) : (
          <div className="fleet-list">
            {fleets.map(fleet => (
              <FleetCard
                key={fleet.id}
                fleet={fleet}
                onEdit={() => setEditFleet(fleet)}
                onDisband={() => remove(fleet.id)}
              />
            ))}
          </div>
        )}
      </div>

      {editFleet && (
        <FleetEditor
          fleet={editFleet}
          designs={designs}
          onAssign={handleAssign}
          onRemove={handleRemoveStack}
          onUpdate={update}
          onClose={() => setEditFleet(null)}
        />
      )}
    </div>
  )
}

function FleetCard({ fleet, onEdit, onDisband }: {
  fleet: Fleet
  onEdit: () => void
  onDisband: () => void
}) {
  const stackCount = fleet.stacks?.length || 0
  const totalShips = fleet.stacks?.reduce((sum, s) => sum + s.ship_count, 0) || 0

  return (
    <div className="fleet-card" onClick={onEdit}>
      <div className="fleet-card-header">
        <span className="fleet-card-name">{fleet.name}</span>
        <span className={`fleet-card-status ${fleet.status}`}>{fleet.status}</span>
      </div>
      <div className="fleet-card-info">
        <span>{stackCount}/9 positions</span>
        <span>{formatNumber(totalShips)} ships</span>
        <span>{fleet.formation}</span>
      </div>
      <div className="fleet-card-actions">
        <button className="p2-btn p2-btn-primary p2-btn-xs" onClick={e => { e.stopPropagation(); onEdit() }}>Edit</button>
        <button className="p2-btn p2-btn-danger p2-btn-xs" onClick={e => { e.stopPropagation(); onDisband() }}>Disband</button>
      </div>
    </div>
  )
}
