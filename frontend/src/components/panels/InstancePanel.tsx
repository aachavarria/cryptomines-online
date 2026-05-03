import { useState, type CSSProperties } from 'react'
import { createPortal } from 'react-dom'
import { MapPin, Layers, Sword, Check, X } from 'lucide-react'
import { useInstances } from '../../hooks/useInstances.ts'
import { useFleets } from '../../hooks/useFleets.ts'
import { formatNumber } from '../../hooks/useCountdown.ts'
import type { Instance, InstanceDetail, InstanceAttemptResponse, Fleet } from '../../types'

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
    background: 'var(--ds-orange-tint)',
    color: 'var(--ds-orange)',
    display: 'inline-grid',
    placeItems: 'center',
  },
  list: {
    display: 'flex',
    flexDirection: 'column',
    gap: 'var(--sp-2)',
    maxHeight: 'calc(100vh - 240px)',
    overflowY: 'auto',
  },
  listItem: {
    display: 'flex',
    alignItems: 'center',
    gap: 'var(--sp-3)',
  },
  difficulty: {
    width: 36,
    height: 36,
    borderRadius: 'var(--r-md)',
    background: 'var(--ds-surface-3)',
    color: 'var(--ds-text-muted)',
    display: 'inline-grid',
    placeItems: 'center',
    fontWeight: 700,
    flexShrink: 0,
  },
  rowBetween: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    gap: 'var(--sp-3)',
  },
  modalSection: {
    display: 'flex',
    flexDirection: 'column',
    gap: 'var(--sp-2)',
    paddingBottom: 'var(--sp-4)',
    borderBottom: '1px solid var(--ds-border)',
    marginBottom: 'var(--sp-4)',
  },
  resultBanner: {
    padding: 'var(--sp-5)',
    borderRadius: 'var(--r-lg)',
    textAlign: 'center',
    fontSize: 'var(--fs-h1)',
    fontWeight: 700,
    letterSpacing: '0.06em',
    textTransform: 'uppercase',
    marginBottom: 'var(--sp-4)',
  },
  fleetList: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, minmax(180px, 1fr))',
    gap: 'var(--sp-2)',
  },
  resourceRow: {
    display: 'inline-flex',
    alignItems: 'center',
    gap: 'var(--sp-2)',
  },
  resDot: {
    width: 10,
    height: 10,
    borderRadius: '50%',
    display: 'inline-block',
  },
}

function classBadge(hullClass: string): string {
  // hullClass: scout, frigate, destroyer, cruiser, battleship, dreadnought, etc.
  const lower = hullClass.toLowerCase()
  if (lower.includes('dreadnought') || lower.includes('battleship')) return 'ds-badge--danger'
  if (lower.includes('cruiser') || lower.includes('destroyer')) return 'ds-badge--orange'
  if (lower.includes('frigate')) return 'ds-badge--info'
  return 'ds-badge--neutral'
}

interface InstanceDetailModalProps {
  instance: Instance
  detail: InstanceDetail | null
  loadingDetail: boolean
  fleets: Fleet[]
  onAttempt: (instanceId: number, fleetIds: string[]) => Promise<InstanceAttemptResponse>
  onClose: () => void
}

function InstanceDetailModal({
  instance, detail, loadingDetail, fleets, onAttempt, onClose,
}: InstanceDetailModalProps) {
  const [selectedFleets, setSelectedFleets] = useState<string[]>([])
  const [attacking, setAttacking] = useState(false)
  const [result, setResult] = useState<InstanceAttemptResponse | null>(null)

  function toggleFleet(id: string) {
    if (selectedFleets.includes(id)) {
      setSelectedFleets(selectedFleets.filter(f => f !== id))
    } else if (selectedFleets.length < instance.max_fleets) {
      setSelectedFleets([...selectedFleets, id])
    }
  }

  async function handleAttack() {
    if (selectedFleets.length === 0) return
    setAttacking(true)
    try {
      const res = await onAttempt(instance.id, selectedFleets)
      setResult(res)
    } catch {
      // Error handled by hook
    } finally {
      setAttacking(false)
    }
  }

  const stationedFleets = fleets.filter(f => f.status === 'stationed' && (f.stacks?.length || 0) > 0)
  const isVictory = result?.result === 'attacker_win'

  return createPortal(
    <div className="ds-modal-backdrop" onClick={e => { if (e.target === e.currentTarget) onClose() }}>
      <div className="ds-modal" role="dialog" aria-modal="true">
        <div className="ds-modal-header">
          <h2 className="ds-modal-title" style={{ display: 'inline-flex', alignItems: 'center', gap: 'var(--sp-2)' }}>
            <MapPin size={20} strokeWidth={1.75} aria-hidden="true" />
            {instance.name}
          </h2>
          <button className="ds-btn-icon" onClick={onClose} aria-label="Close">
            <X size={16} strokeWidth={2} aria-hidden="true" />
          </button>
        </div>
        <div className="ds-modal-body">
          {result ? (
            /* Battle Result */
            <div>
              <div
                style={{
                  ...styles.resultBanner,
                  background: isVictory ? 'var(--ds-success-tint)' : 'var(--ds-danger-tint)',
                  color: isVictory ? '#15803D' : '#B91C1C',
                  border: `2px solid ${isVictory ? 'var(--ds-success)' : 'var(--ds-danger)'}`,
                }}
              >
                {isVictory ? 'Victory' : 'Defeat'}
              </div>
              <div className="ds-card" style={{ display: 'flex', flexDirection: 'column', gap: 'var(--sp-2)' }}>
                <div style={styles.rowBetween}>
                  <span className="ds-text-muted">Rounds</span>
                  <span className="ds-mono">{result.total_rounds}</span>
                </div>
                <div style={styles.rowBetween}>
                  <span className="ds-text-muted">EXP Gained</span>
                  <span className="ds-mono" style={{ color: 'var(--ds-teal-dark)', fontWeight: 600 }}>
                    +{formatNumber(result.exp_gained)}
                  </span>
                </div>
                {result.treasure_box && (
                  <>
                    <div className="ds-caption" style={{ marginTop: 'var(--sp-2)' }}>Rewards</div>
                    <div style={styles.rowBetween}>
                      <span style={styles.resourceRow}>
                        <span style={{ ...styles.resDot, background: 'var(--ds-metal)' }} />
                        Metal
                      </span>
                      <span className="ds-mono">+{formatNumber(result.treasure_box.resources.metal)}</span>
                    </div>
                    <div style={styles.rowBetween}>
                      <span style={styles.resourceRow}>
                        <span style={{ ...styles.resDot, background: 'var(--ds-he3)' }} />
                        He3
                      </span>
                      <span className="ds-mono">+{formatNumber(result.treasure_box.resources.he3)}</span>
                    </div>
                    <div style={styles.rowBetween}>
                      <span style={styles.resourceRow}>
                        <span style={{ ...styles.resDot, background: 'var(--ds-gold)' }} />
                        Gold
                      </span>
                      <span className="ds-mono">+{formatNumber(result.treasure_box.resources.gold)}</span>
                    </div>
                    {result.treasure_box.blueprint_id && (
                      <div style={{ ...styles.rowBetween, color: 'var(--ds-orange-strong)' }}>
                        <span>Blueprint Drop!</span>
                        <span className="ds-mono">ID: {result.treasure_box.blueprint_id}</span>
                      </div>
                    )}
                  </>
                )}
                {Object.keys(result.losses.ships_destroyed).length > 0 && (
                  <>
                    <div className="ds-caption" style={{ marginTop: 'var(--sp-2)' }}>Losses</div>
                    {Object.entries(result.losses.ships_destroyed).map(([designId, count]) => (
                      <div key={designId} style={styles.rowBetween}>
                        <span className="ds-text-muted">Design {designId.slice(0, 8)}...</span>
                        <span className="ds-mono" style={{ color: 'var(--ds-danger)' }}>-{count} ships</span>
                      </div>
                    ))}
                    <div style={styles.rowBetween}>
                      <span className="ds-text-muted">He3 Consumed</span>
                      <span className="ds-mono" style={{ color: 'var(--ds-danger)' }}>
                        {formatNumber(result.losses.he3_consumed)}
                      </span>
                    </div>
                  </>
                )}
              </div>
              <button
                className="ds-btn ds-btn-secondary ds-btn--block"
                style={{ marginTop: 'var(--sp-4)' }}
                onClick={onClose}
              >
                Close
              </button>
            </div>
          ) : (
            /* Pre-battle screen */
            <>
              <div style={styles.modalSection}>
                <div style={styles.rowBetween}>
                  <span className="ds-text-muted">Level Requirement</span>
                  <span className="ds-mono">{instance.required_level}</span>
                </div>
                <div style={styles.rowBetween}>
                  <span className="ds-text-muted">Max Fleets</span>
                  <span className="ds-mono">{instance.max_fleets}</span>
                </div>
                <div style={styles.rowBetween}>
                  <span className="ds-text-muted">EXP Reward</span>
                  <span className="ds-mono" style={{ color: 'var(--ds-teal-dark)', fontWeight: 600 }}>
                    {formatNumber(instance.exp_reward)}
                  </span>
                </div>
              </div>

              {loadingDetail ? (
                <div style={{ display: 'flex', justifyContent: 'center', padding: 'var(--sp-4)' }}>
                  <div className="loading-spinner" />
                </div>
              ) : detail && (
                <div style={styles.modalSection}>
                  <div className="ds-caption">Enemy Forces</div>
                  {detail.enemy_fleets.map((ef, i) => (
                    <div key={i} style={styles.rowBetween}>
                      <span className={`ds-badge ${classBadge(ef.hull_class)}`}>{ef.hull_class}</span>
                      <span className="ds-mono">x{ef.ship_count}</span>
                      <span className="ds-text-muted ds-mono" style={{ fontSize: 'var(--fs-sm)' }}>
                        ~{formatNumber(ef.power_estimate)} power
                      </span>
                    </div>
                  ))}
                </div>
              )}

              <div style={{ marginBottom: 'var(--sp-4)' }}>
                <div className="ds-caption" style={{ marginBottom: 'var(--sp-2)' }}>
                  Select Fleets ({selectedFleets.length}/{instance.max_fleets})
                </div>
                {stationedFleets.length === 0 ? (
                  <div className="ds-text-muted" style={{ textAlign: 'center', padding: 'var(--sp-4)' }}>
                    No stationed fleets available. Create a fleet first.
                  </div>
                ) : (
                  <div style={styles.fleetList}>
                    {stationedFleets.map(f => {
                      const isSelected = selectedFleets.includes(f.id)
                      const totalShips = f.stacks?.reduce((sum, s) => sum + s.ship_count, 0) || 0
                      return (
                        <button
                          key={f.id}
                          className="ds-list-item"
                          aria-selected={isSelected}
                          style={{
                            textAlign: 'left',
                            border: 'none',
                            background: 'transparent',
                            font: 'inherit',
                            padding: 'var(--sp-3)',
                            cursor: 'pointer',
                            borderRadius: 'var(--r-lg)',
                            ...(isSelected ? { borderColor: 'var(--ds-orange)', borderWidth: 2, borderStyle: 'solid' } : {}),
                          }}
                          onClick={() => toggleFleet(f.id)}
                        >
                          <div style={{ fontWeight: 600 }}>{f.name}</div>
                          <div className="ds-text-muted ds-mono" style={{ fontSize: 'var(--fs-sm)', marginTop: 4 }}>
                            {formatNumber(totalShips)} ships
                          </div>
                        </button>
                      )
                    })}
                  </div>
                )}
              </div>

              <button
                className="ds-btn ds-btn-primary ds-btn--block"
                onClick={handleAttack}
                disabled={selectedFleets.length === 0 || attacking}
              >
                <Sword size={16} strokeWidth={2} aria-hidden="true" />
                {attacking ? 'Attacking...' : 'Attack'}
              </button>
            </>
          )}
        </div>
      </div>
    </div>,
    document.body,
  )
}

export default function InstancePanel() {
  const { instances, loading, error, getDetail, attempt, isCompleted } = useInstances()
  const { fleets } = useFleets()
  const [selectedInstance, setSelectedInstance] = useState<Instance | null>(null)
  const [instanceDetail, setInstanceDetail] = useState<InstanceDetail | null>(null)
  const [loadingDetail, setLoadingDetail] = useState(false)

  if (loading) {
    return (
      <div className="ds-panel" style={{ display: 'flex', alignItems: 'center', gap: 'var(--sp-3)', justifyContent: 'center' }}>
        <div className="loading-spinner" />
        <span className="ds-text-muted">Loading Instances...</span>
      </div>
    )
  }

  async function handleSelect(inst: Instance) {
    setSelectedInstance(inst)
    setLoadingDetail(true)
    setInstanceDetail(null)
    try {
      const detail = await getDetail(inst.id)
      setInstanceDetail(detail)
    } catch {
      // Error handled
    } finally {
      setLoadingDetail(false)
    }
  }

  const completedCount = instances.filter(i => isCompleted(i.id)).length

  return (
    <div className="ds-panel" style={styles.panel}>
      <div style={styles.header}>
        <span style={styles.headerIcon} aria-hidden="true">
          <Layers size={22} strokeWidth={1.75} />
        </span>
        <div>
          <h2 className="ds-h2" style={{ margin: 0 }}>Instances</h2>
          <div className="ds-text-muted" style={{ fontSize: 'var(--fs-sm)' }}>
            Normal Instances - <span className="ds-mono">{completedCount}/{instances.length}</span> completed
          </div>
        </div>
      </div>

      {error && (
        <div className="ds-badge ds-badge--danger" role="alert">{error}</div>
      )}

      <div style={styles.list}>
        {instances.length === 0 && (
          <div className="ds-text-muted" style={{ textAlign: 'center', padding: 'var(--sp-6)' }}>
            No instances available.
          </div>
        )}
        {instances.map(inst => {
          const completed = isCompleted(inst.id)
          return (
            <div
              key={inst.id}
              className="ds-list-item"
              style={{
                ...styles.listItem,
                ...(completed ? { opacity: 0.7 } : {}),
              }}
              onClick={() => handleSelect(inst)}
            >
              <div style={styles.difficulty}>{inst.difficulty}</div>
              <div style={{ flex: 1, minWidth: 0 }}>
                <div style={{ fontWeight: 600 }}>{inst.name}</div>
                <div className="ds-text-muted" style={{ fontSize: 'var(--fs-sm)', marginTop: 2 }}>
                  Lv {inst.required_level} | {inst.max_fleets} fleet{inst.max_fleets > 1 ? 's' : ''}
                </div>
              </div>
              <div style={{ display: 'flex', alignItems: 'center', gap: 'var(--sp-2)' }}>
                <span className="ds-badge ds-badge--teal">
                  +{formatNumber(inst.exp_reward)} EXP
                </span>
                {completed && (
                  <span
                    className="ds-badge ds-badge--success"
                    style={{ display: 'inline-flex', alignItems: 'center', gap: 4 }}
                  >
                    <Check size={12} strokeWidth={2} aria-hidden="true" />
                    Done
                  </span>
                )}
              </div>
            </div>
          )
        })}
      </div>

      {selectedInstance && (
        <InstanceDetailModal
          instance={selectedInstance}
          detail={instanceDetail}
          loadingDetail={loadingDetail}
          fleets={fleets}
          onAttempt={attempt}
          onClose={() => { setSelectedInstance(null); setInstanceDetail(null) }}
        />
      )}
    </div>
  )
}
