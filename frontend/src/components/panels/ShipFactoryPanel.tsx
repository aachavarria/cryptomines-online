import { useState } from 'react'
import { createPortal } from 'react-dom'
import { Factory, X, Hammer, Box } from 'lucide-react'
import LoadingButton from '../common/LoadingButton'
import { useShipFactory } from '../../hooks/useShipFactory.ts'
import { useShipDesigns } from '../../hooks/useShipDesigns.ts'
import { useShipInventory } from '../../hooks/useShipInventory.ts'
import { useCountdown, formatDuration, formatNumber } from '../../hooks/useCountdown.ts'
import type { ProductionSlot, ShipDesign } from '../../types'

const HULL_CLASS_TINT: Record<string, string> = {
  frigate: 'var(--ds-info)',
  cruiser: 'var(--ds-orange-strong)',
  battleship: 'var(--ds-danger)',
}

function SlotTimer({ finishAt }: { finishAt: string }) {
  const timeLeft = useCountdown(finishAt)
  if (!timeLeft) return <span className="ds-badge ds-badge--success">Complete</span>
  return <span className="ds-mono ds-text-muted">{timeLeft}</span>
}

interface BuildModalProps {
  slot: number
  designs: ShipDesign[]
  speedBonus: number
  onBuild: (designId: string, quantity: number, slot: number) => Promise<void>
  onClose: () => void
}

function BuildModal({ slot, designs, speedBonus, onBuild, onClose }: BuildModalProps) {
  const [selectedDesign, setSelectedDesign] = useState<string>('')
  const [quantity, setQuantity] = useState(100)
  const [submitting, setSubmitting] = useState(false)

  const design = designs.find(d => d.id === selectedDesign)

  async function handleBuild() {
    if (!selectedDesign || quantity < 1) return
    setSubmitting(true)
    try {
      await onBuild(selectedDesign, quantity, slot)
      onClose()
    } catch {
      setSubmitting(false)
    }
  }

  return createPortal(
    <div
      className="ds-modal-backdrop"
      onClick={e => {
        if (e.target === e.currentTarget) onClose()
      }}
    >
      <div className="ds-modal ds-modal--sm" role="dialog" aria-modal="true">
        <div className="ds-modal-header">
          <div className="sf-modal-title">
            <Hammer size={20} strokeWidth={1.75} />
            <h3 className="ds-modal-title">Build Ships — Slot {slot}</h3>
          </div>
          <button className="ds-btn-icon" onClick={onClose} aria-label="Close">
            <X size={18} strokeWidth={1.75} />
          </button>
        </div>
        <div className="ds-modal-body sf-modal-body">
          <div className="sf-form-group">
            <label className="ds-caption">Ship Design</label>
            <select
              className="ds-select"
              value={selectedDesign}
              onChange={e => setSelectedDesign(e.target.value)}
            >
              <option value="">— Select Design —</option>
              {designs.map(d => (
                <option key={d.id} value={d.id}>
                  {d.name} ({d.hull_class})
                </option>
              ))}
            </select>
          </div>

          <div className="sf-form-group">
            <label className="ds-caption">Quantity (1 – 2,000,000)</label>
            <input
              type="number"
              className="ds-input"
              value={quantity}
              min={1}
              max={2000000}
              onChange={e =>
                setQuantity(Math.max(1, Math.min(2000000, parseInt(e.target.value) || 1)))
              }
            />
          </div>

          {design && (
            <div className="ds-card sf-summary">
              <div className="sf-summary-title">Build Summary</div>
              <SummaryRow label="Time per ship" value={formatDuration(design.build_time_seconds)} />
              <SummaryRow
                label="Total time"
                value={formatDuration(
                  Math.ceil(design.build_time_seconds * quantity * (1 - speedBonus / 100)),
                )}
              />
              {speedBonus > 0 && (
                <SummaryRow label="Speed bonus" value={`-${speedBonus}%`} accent="teal" />
              )}
              <hr className="ds-divider" />
              <SummaryRow
                label={
                  <span className="sf-cost-label">
                    <span className="ds-resource-dot ds-resource-dot--metal" /> Total Metal
                  </span>
                }
                value={formatNumber(design.metal_cost * quantity)}
              />
              <SummaryRow
                label={
                  <span className="sf-cost-label">
                    <span className="ds-resource-dot ds-resource-dot--he3" /> Total He3
                  </span>
                }
                value={formatNumber(design.he3_cost * quantity)}
              />
              <SummaryRow
                label={
                  <span className="sf-cost-label">
                    <span className="ds-resource-dot ds-resource-dot--gold" /> Total Gold
                  </span>
                }
                value={formatNumber(design.gold_cost * quantity)}
              />
            </div>
          )}
        </div>
        <div className="ds-modal-footer">
          <button className="ds-btn ds-btn-ghost" onClick={onClose}>
            Cancel
          </button>
          <LoadingButton
            className="ds-btn ds-btn-secondary"
            onClick={handleBuild}
            disabled={!selectedDesign || quantity < 1}
            loading={submitting}
          >
            Start Production
          </LoadingButton>
        </div>

        <style>{`
          .sf-modal-title { display: flex; align-items: center; gap: var(--sp-2); color: var(--ds-text); }
          .sf-modal-title h3 { margin: 0; }
          .sf-modal-body { display: flex; flex-direction: column; gap: var(--sp-3); }
          .sf-form-group { display: flex; flex-direction: column; gap: 6px; }
          .sf-summary { display: flex; flex-direction: column; gap: var(--sp-2); }
          .sf-summary-title { font-weight: var(--fw-semibold); color: var(--ds-text); }
          .sf-summary-row {
            display: flex; align-items: center; justify-content: space-between;
            font-size: var(--fs-sm);
          }
          .sf-summary-row dt, .sf-summary-row dd { margin: 0; }
          .sf-summary-row dt { color: var(--ds-text-muted); }
          .sf-summary-row dd { font-weight: var(--fw-semibold); color: var(--ds-text); }
          .sf-summary-row dd.accent-teal { color: var(--ds-teal-dark); }
          .sf-cost-label { display: inline-flex; align-items: center; gap: var(--sp-2); }
        `}</style>
      </div>
    </div>,
    document.body,
  )
}

function SummaryRow({
  label,
  value,
  accent,
}: {
  label: React.ReactNode
  value: React.ReactNode
  accent?: 'teal'
}) {
  return (
    <div className="sf-summary-row">
      <dt>{label}</dt>
      <dd className={accent === 'teal' ? 'accent-teal' : ''}>{value}</dd>
    </div>
  )
}

export default function ShipFactoryPanel() {
  const { factory, slots, loading, error, build, cancel } = useShipFactory()
  const { designs } = useShipDesigns()
  const { ships, loading: shipsLoading } = useShipInventory()
  const [buildSlot, setBuildSlot] = useState<number | null>(null)

  if (loading) {
    return (
      <div className="ds-panel sf-loading">
        <div className="loading-spinner" />
        <span>Loading Ship Factory...</span>
      </div>
    )
  }

  if (error || !factory) {
    return (
      <div className="ds-panel">
        <p className="ds-text-muted">{error || 'Ship Factory not found. Build one first.'}</p>
      </div>
    )
  }

  async function handleBuild(designId: string, quantity: number, slot: number) {
    await build({ ship_design_id: designId, quantity, production_slot: slot })
  }

  return (
    <div className="ds-panel sf-panel">
      <header className="sf-header">
        <div className="sf-header-icon">
          <Factory size={22} strokeWidth={1.75} />
        </div>
        <div className="sf-header-text">
          <h2 className="ds-h2">Ship Factory</h2>
          <p className="ds-text-muted sf-subtitle">
            Level <span className="ds-mono">{factory.level}</span> · Speed bonus{' '}
            <span className="ds-mono">{factory.speed_bonus_pct}%</span>
          </p>
        </div>
      </header>

      <section className="sf-section">
        <h3 className="ds-h3">
          Production Slots ({factory.production_slots}/5)
        </h3>
        <div className="sf-slots">
          {slots.map(slot => (
            <SlotCard
              key={slot.slot}
              slot={slot}
              onBuild={() => setBuildSlot(slot.slot)}
              onCancel={() => cancel(slot.slot)}
            />
          ))}
        </div>
      </section>

      <section className="sf-section">
        <h3 className="ds-h3">Ship Inventory</h3>
        {shipsLoading ? (
          <div className="ds-card sf-loading">
            <div className="loading-spinner" />
            <span>Loading ships...</span>
          </div>
        ) : ships.length === 0 ? (
          <div className="ds-card sf-empty">
            <Box size={32} strokeWidth={1.5} />
            <p>No ships constructed yet. Start building ships above.</p>
          </div>
        ) : (
          <div className="sf-inventory">
            {ships.map(ship => (
              <ShipInventoryCard key={ship.id} ship={ship} />
            ))}
          </div>
        )}
      </section>

      {buildSlot !== null && (
        <BuildModal
          slot={buildSlot}
          designs={designs}
          speedBonus={factory.speed_bonus_pct}
          onBuild={handleBuild}
          onClose={() => setBuildSlot(null)}
        />
      )}

      <style>{`
        .sf-panel { display: flex; flex-direction: column; gap: var(--sp-4); max-width: 1100px; }
        .sf-loading { display: flex; align-items: center; gap: var(--sp-3); }
        .sf-header { display: flex; align-items: center; gap: var(--sp-3); }
        .sf-header-icon {
          width: 44px; height: 44px;
          display: grid; place-items: center;
          background: var(--ds-teal-tint);
          color: var(--ds-teal-dark);
          border-radius: var(--r-lg);
        }
        .sf-header-text h2 { margin: 0; }
        .sf-subtitle { margin: 2px 0 0; font-size: var(--fs-sm); }
        .sf-section { display: flex; flex-direction: column; gap: var(--sp-3); }
        .sf-section h3 { margin: 0; }

        .sf-slots {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
          gap: var(--sp-3);
        }
        .sf-empty {
          display: flex; flex-direction: column; align-items: center; gap: var(--sp-2);
          padding: var(--sp-6); color: var(--ds-text-muted);
        }
        .sf-empty p { margin: 0; }

        .sf-inventory {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
          gap: var(--sp-3);
        }
      `}</style>
    </div>
  )
}

function SlotCard({
  slot,
  onBuild,
  onCancel,
}: {
  slot: ProductionSlot
  onBuild: () => void
  onCancel: () => void
}) {
  if (slot.in_use) {
    return (
      <div className="ds-card sf-slot sf-slot--building">
        <div className="sf-slot-head">
          <span className="ds-caption">Slot {slot.slot}</span>
          <span className="ds-badge ds-badge--teal">Building</span>
        </div>
        <div className="sf-slot-design">{slot.design_name || 'Building…'}</div>
        <div className="sf-slot-meta">
          <span className="ds-mono">×{formatNumber(slot.quantity || 0)}</span>
          {slot.build_finish_at && <SlotTimer finishAt={slot.build_finish_at} />}
        </div>
        <button className="ds-btn ds-btn-danger ds-btn--sm" onClick={onCancel}>
          Cancel
        </button>
        <style>{`
          .sf-slot { display: flex; flex-direction: column; gap: var(--sp-2); }
          .sf-slot-head { display: flex; align-items: center; justify-content: space-between; }
          .sf-slot-design { font-weight: var(--fw-semibold); color: var(--ds-text); }
          .sf-slot-meta {
            display: flex; align-items: center; justify-content: space-between;
            font-size: var(--fs-sm);
          }
        `}</style>
      </div>
    )
  }

  return (
    <div className="ds-card sf-slot sf-slot--empty">
      <div className="sf-slot-head">
        <span className="ds-caption">Slot {slot.slot}</span>
        <span className="ds-badge ds-badge--neutral">Idle</span>
      </div>
      <button className="ds-btn ds-btn-secondary ds-btn--sm" onClick={onBuild}>
        Build
      </button>
      <style>{`
        .sf-slot { display: flex; flex-direction: column; gap: var(--sp-2); }
        .sf-slot-head { display: flex; align-items: center; justify-content: space-between; }
      `}</style>
    </div>
  )
}

function ShipInventoryCard({ ship }: { ship: import('../../types').AvailableShip }) {
  const tint = HULL_CLASS_TINT[ship.hull_class] || 'var(--ds-text-muted)'

  return (
    <div className="ds-card sf-inv-card">
      <div className="sf-inv-class" style={{ color: tint, borderColor: tint }}>
        {ship.hull_class.charAt(0).toUpperCase()}
      </div>
      <div className="sf-inv-info">
        <div className="sf-inv-name">{ship.design_name}</div>
        <div className="ds-text-muted sf-inv-hull">
          {ship.hull_class.charAt(0).toUpperCase() + ship.hull_class.slice(1)}
        </div>
      </div>
      <div className="sf-inv-qty">
        <div className="ds-caption">Qty</div>
        <div className="ds-mono sf-inv-qty-val">{formatNumber(ship.quantity)}</div>
      </div>
      <style>{`
        .sf-inv-card {
          display: flex; align-items: center; gap: var(--sp-3);
        }
        .sf-inv-class {
          width: 36px; height: 36px;
          display: grid; place-items: center;
          border-radius: var(--r-md);
          border: 2px solid currentColor;
          font-weight: var(--fw-bold);
          flex-shrink: 0;
        }
        .sf-inv-info { flex: 1; min-width: 0; }
        .sf-inv-name {
          font-weight: var(--fw-semibold); color: var(--ds-text);
          white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
        }
        .sf-inv-hull { font-size: var(--fs-sm); }
        .sf-inv-qty { text-align: right; }
        .sf-inv-qty-val { font-size: var(--fs-h3); font-weight: var(--fw-semibold); color: var(--ds-text); }
      `}</style>
    </div>
  )
}
