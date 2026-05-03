import { useState, type CSSProperties } from 'react'
import { createPortal } from 'react-dom'
import { Rocket, X, Plus, Trash2, Pencil } from 'lucide-react'
import { useFleets } from '../../hooks/useFleets.ts'
import { useShipDesigns } from '../../hooks/useShipDesigns.ts'
import type { Fleet, FleetStack, ShipDesign, CreateFleetRequest } from '../../types'
import { formatNumber } from '../../hooks/useCountdown.ts'

const FORMATIONS = ['Phalanx', 'Diamond', 'Arrow', 'Defensive', 'Spread']
const TARGETING = ['Weakest First', 'Strongest First', 'Random', 'Closest']
const GRID_LABELS = [
  ['Front-L', 'Front-C', 'Front-R'],
  ['Mid-L', 'Mid-C', 'Mid-R'],
  ['Back-L', 'Back-C', 'Back-R'],
]

const styles: Record<string, CSSProperties> = {
  panel: {
    maxWidth: 720,
    margin: '0 auto',
    display: 'flex',
    flexDirection: 'column',
    gap: 'var(--sp-4)',
  },
  header: {
    display: 'flex',
    alignItems: 'center',
    gap: 'var(--sp-3)',
  },
  headerIcon: {
    width: 44,
    height: 44,
    borderRadius: 'var(--r-md)',
    background: 'var(--ds-teal-tint)',
    color: 'var(--ds-teal)',
    display: 'inline-grid',
    placeItems: 'center',
  },
  loading: {
    display: 'flex',
    alignItems: 'center',
    gap: 'var(--sp-3)',
    justifyContent: 'center',
  },
  createRow: {
    display: 'flex',
    gap: 'var(--sp-2)',
    alignItems: 'center',
  },
  fleetList: {
    display: 'flex',
    flexDirection: 'column',
    gap: 'var(--sp-2)',
  },
  cardHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 'var(--sp-2)',
  },
  cardName: {
    display: 'inline-flex',
    alignItems: 'center',
    gap: 'var(--sp-2)',
    fontWeight: 600,
    color: 'var(--ds-text)',
  },
  cardInfo: {
    display: 'flex',
    flexWrap: 'wrap',
    gap: 'var(--sp-3)',
    fontSize: 'var(--fs-sm)',
    color: 'var(--ds-text-muted)',
    marginBottom: 'var(--sp-3)',
  },
  cardActions: {
    display: 'flex',
    gap: 'var(--sp-2)',
  },
  config: {
    display: 'flex',
    gap: 'var(--sp-3)',
    marginBottom: 'var(--sp-4)',
  },
  formGroup: {
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
    gap: 'var(--sp-1)',
  },
  grid: {
    display: 'flex',
    flexDirection: 'column',
    gap: 'var(--sp-1)',
    marginBottom: 'var(--sp-4)',
  },
  gridRow: {
    display: 'flex',
    gap: 'var(--sp-1)',
  },
  cell: {
    flex: 1,
    aspectRatio: '1.2',
    minHeight: 88,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 2,
    position: 'relative',
    cursor: 'pointer',
    padding: 'var(--sp-2)',
    textAlign: 'center',
  },
  cellOccupied: {
    background: 'var(--ds-teal-tint)',
    borderColor: 'var(--ds-teal)',
  },
  cellTarget: {
    background: 'var(--ds-orange-tint)',
    borderColor: 'var(--ds-orange)',
  },
  cellRemove: {
    position: 'absolute',
    top: 4,
    right: 4,
  },
  assignForm: {
    borderColor: 'var(--ds-teal)',
    background: 'var(--ds-teal-tint)',
  },
  assignRow: {
    display: 'flex',
    gap: 'var(--sp-2)',
    alignItems: 'center',
    flexWrap: 'wrap',
  },
}

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

  return createPortal(
    <div className="ds-modal-backdrop" onClick={e => { if (e.target === e.currentTarget) onClose() }}>
      <div className="ds-modal ds-modal--lg" role="dialog" aria-modal="true">
        <div className="ds-modal-header">
          <h2 className="ds-modal-title">Fleet: {fleet.name}</h2>
          <button className="ds-btn-icon" onClick={onClose} aria-label="Close">
            <X size={16} strokeWidth={2} aria-hidden="true" />
          </button>
        </div>
        <div className="ds-modal-body">
          {/* Formation & Targeting */}
          <div style={styles.config}>
            <div style={styles.formGroup}>
              <label className="ds-caption">Formation</label>
              <select
                className="ds-select"
                value={fleet.formation}
                onChange={e => onUpdate(fleet.id, { formation: e.target.value })}
              >
                {FORMATIONS.map(f => <option key={f} value={f}>{f}</option>)}
              </select>
            </div>
            <div style={styles.formGroup}>
              <label className="ds-caption">Targeting</label>
              <select
                className="ds-select"
                value={fleet.targeting_command}
                onChange={e => onUpdate(fleet.id, { targeting_command: e.target.value })}
              >
                {TARGETING.map(t => <option key={t} value={t}>{t}</option>)}
              </select>
            </div>
          </div>

          {/* 3x3 Grid */}
          <div style={styles.grid}>
            {[0, 1, 2].map(row => (
              <div key={row} style={styles.gridRow}>
                {[0, 1, 2].map(col => {
                  const stack = getStack(row, col)
                  const isAssignTarget = assignTarget?.row === row && assignTarget?.col === col
                  const cellStyle: CSSProperties = {
                    ...styles.cell,
                    ...(stack ? styles.cellOccupied : {}),
                    ...(isAssignTarget ? styles.cellTarget : {}),
                  }
                  return (
                    <div
                      key={col}
                      className="ds-card"
                      style={cellStyle}
                      onClick={() => {
                        if (!stack) setAssignTarget({ row, col })
                      }}
                    >
                      <div className="ds-caption">{GRID_LABELS[row][col]}</div>
                      {stack ? (
                        <>
                          <div style={{ fontSize: 'var(--fs-sm)', fontWeight: 600 }}>
                            {stack.ship_design_name || 'Ships'}
                          </div>
                          <div className="ds-mono" style={{ color: 'var(--ds-teal-dark)', fontWeight: 600 }}>
                            x{formatNumber(stack.ship_count)}
                          </div>
                          <button
                            className="ds-btn-icon ds-btn-icon--sm"
                            style={styles.cellRemove}
                            onClick={e => { e.stopPropagation(); onRemove(fleet.id, row, col) }}
                            aria-label="Remove stack"
                          >
                            <X size={12} strokeWidth={2} aria-hidden="true" />
                          </button>
                        </>
                      ) : (
                        <div className="ds-text-soft" style={{ fontSize: 'var(--fs-caption)' }}>Empty</div>
                      )}
                    </div>
                  )
                })}
              </div>
            ))}
          </div>

          {/* Assign Stack Form */}
          {assignTarget && (
            <div className="ds-card" style={styles.assignForm}>
              <div className="ds-h3" style={{ marginBottom: 'var(--sp-3)' }}>
                Assign Ships to {GRID_LABELS[assignTarget.row][assignTarget.col]}
              </div>
              <div style={styles.assignRow}>
                <select
                  className="ds-select"
                  style={{ flex: 1, minWidth: 180 }}
                  value={selectedDesign}
                  onChange={e => setSelectedDesign(e.target.value)}
                >
                  <option value="">-- Select Design --</option>
                  {designs.map(d => {
                    const built = (d as ShipDesign & { ships_built?: number }).ships_built ?? 0
                    return (
                      <option key={d.id} value={d.id}>
                        {d.name} ({built} available)
                      </option>
                    )
                  })}
                </select>
                <input
                  type="number"
                  className="ds-input"
                  style={{ width: 100 }}
                  value={shipCount}
                  min={1}
                  max={3000}
                  onChange={e => setShipCount(Math.max(1, Math.min(3000, parseInt(e.target.value) || 1)))}
                />
                <button
                  className="ds-btn ds-btn-primary ds-btn--sm"
                  onClick={handleAssign}
                  disabled={!selectedDesign || shipCount < 1}
                >
                  Assign
                </button>
                <button
                  className="ds-btn ds-btn-ghost ds-btn--sm"
                  onClick={() => setAssignTarget(null)}
                >
                  Cancel
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>,
    document.body,
  )
}

export default function FleetPanel() {
  const { fleets, loading, error, create, update, remove, addStack, removeStackFromFleet } = useFleets()
  const { designs } = useShipDesigns()
  const [editFleet, setEditFleet] = useState<Fleet | null>(null)
  const [creating, setCreating] = useState(false)
  const [newName, setNewName] = useState('')

  if (loading) {
    return (
      <div className="ds-panel" style={styles.loading}>
        <div className="loading-spinner" />
        <span className="ds-text-muted">Loading Fleets...</span>
      </div>
    )
  }

  async function handleCreate() {
    if (!newName) return
    setCreating(true)
    try {
      const req: CreateFleetRequest = {
        name: newName,
        formation: 'Phalanx',
        targeting_command: 'Weakest First',
      }
      try {
        const planetId = typeof localStorage !== 'undefined' && typeof localStorage.getItem === 'function'
          ? localStorage.getItem('current_planet_id')
          : null
        if (planetId) req.planet_id = planetId
      } catch {
        // localStorage may be unavailable in some test environments
      }
      const fleet = await create(req)
      setNewName('')
      setEditFleet(fleet)
    } finally {
      setCreating(false)
    }
  }

  async function handleAssign(fleetId: string, designId: string, row: number, col: number, count: number) {
    await addStack(fleetId, {
      ship_design_id: designId,
      grid_row: row,
      grid_col: col,
      ship_count: count,
    })
    const updated = fleets.find(f => f.id === fleetId)
    if (updated) setEditFleet({ ...updated })
  }

  async function handleRemoveStack(fleetId: string, row: number, col: number) {
    await removeStackFromFleet(fleetId, row, col)
    const updated = fleets.find(f => f.id === fleetId)
    if (updated) setEditFleet({ ...updated })
  }

  return (
    <div className="ds-panel" style={styles.panel}>
      <div style={styles.header}>
        <span style={styles.headerIcon} aria-hidden="true">
          <Rocket size={22} strokeWidth={1.75} />
        </span>
        <div>
          <h2 className="ds-h2" style={{ margin: 0 }}>Fleet Management</h2>
          <div className="ds-text-muted" style={{ fontSize: 'var(--fs-sm)' }}>
            {fleets.length} fleet{fleets.length !== 1 ? 's' : ''}
          </div>
        </div>
      </div>

      {error && (
        <div className="ds-badge ds-badge--danger" role="alert">{error}</div>
      )}

      <div style={styles.createRow}>
        <input
          className="ds-input"
          value={newName}
          placeholder="New fleet name..."
          onChange={e => setNewName(e.target.value)}
          onKeyDown={e => { if (e.key === 'Enter') handleCreate() }}
        />
        <button
          className="ds-btn ds-btn-primary"
          onClick={handleCreate}
          disabled={!newName || creating}
        >
          <Plus size={16} strokeWidth={2} aria-hidden="true" />
          {creating ? '...' : 'Create Fleet'}
        </button>
      </div>

      <div>
        {fleets.length === 0 ? (
          <div className="ds-text-muted" style={{ textAlign: 'center', padding: 'var(--sp-6)' }}>
            No fleets yet. Create a fleet and assign ships.
          </div>
        ) : (
          <div style={styles.fleetList}>
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

  // Fleet speed = slowest ship design's total_movement (min MOV across stacks)
  const fleetMOV = fleet.stacks && fleet.stacks.length > 0
    ? Math.min(...fleet.stacks.filter(s => s.total_movement !== undefined).map(s => s.total_movement ?? 0))
    : 0

  const statusBadge =
    fleet.status === 'stationed' ? 'ds-badge--teal' :
    fleet.status === 'traveling' ? 'ds-badge--info' :
    fleet.status === 'combat' ? 'ds-badge--danger' :
    'ds-badge--neutral'

  return (
    <div className="ds-list-item" onClick={onEdit}>
      <div style={styles.cardHeader}>
        <span style={styles.cardName}>
          <Rocket size={16} strokeWidth={1.75} aria-hidden="true" />
          {fleet.name}
        </span>
        <span className={`ds-badge ${statusBadge}`}>{fleet.status}</span>
      </div>
      <div style={styles.cardInfo}>
        <span>{stackCount}/9 positions</span>
        <span>{formatNumber(totalShips)} ships</span>
        <span>{fleet.formation}</span>
        {fleetMOV > 0 && (
          <span title="Fleet speed (slowest ship)">
            MOV: <span className="ds-mono">{fleetMOV}</span>
          </span>
        )}
      </div>
      <div style={styles.cardActions}>
        <button
          className="ds-btn ds-btn-ghost ds-btn--sm"
          onClick={e => { e.stopPropagation(); onEdit() }}
        >
          <Pencil size={14} strokeWidth={2} aria-hidden="true" />
          Edit
        </button>
        <button
          className="ds-btn ds-btn-danger ds-btn--sm"
          onClick={e => { e.stopPropagation(); onDisband() }}
        >
          <Trash2 size={14} strokeWidth={2} aria-hidden="true" />
          Disband
        </button>
      </div>
    </div>
  )
}
