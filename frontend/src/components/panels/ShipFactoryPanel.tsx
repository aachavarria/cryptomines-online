import { useState } from 'react'
import { createPortal } from 'react-dom'
import LoadingButton from '../common/LoadingButton'
import { useShipFactory } from '../../hooks/useShipFactory.ts'
import { useShipDesigns } from '../../hooks/useShipDesigns.ts'
import { useShipInventory } from '../../hooks/useShipInventory.ts'
import { useCountdown, formatDuration, formatNumber } from '../../hooks/useCountdown.ts'
import type { ProductionSlot, ShipDesign } from '../../types'

function SlotTimer({ finishAt }: { finishAt: string }) {
  const timeLeft = useCountdown(finishAt)
  if (!timeLeft) return <span className="slot-status-text done">Complete</span>
  return <span className="slot-timer">{timeLeft}</span>
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
    <div className="p2-modal-backdrop" onClick={e => { if (e.target === e.currentTarget) onClose() }}>
      <div className="p2-modal p2-modal-sm">
        <div className="p2-modal-header">
          <span className="p2-modal-title">Build Ships - Slot {slot}</span>
          <button className="p2-modal-close" onClick={onClose}>X</button>
        </div>
        <div className="p2-modal-body">
          <div className="p2-form-group">
            <label className="p2-label">Ship Design</label>
            <select
              className="p2-select"
              value={selectedDesign}
              onChange={e => setSelectedDesign(e.target.value)}
            >
              <option value="">-- Select Design --</option>
              {designs.map(d => (
                <option key={d.id} value={d.id}>
                  {d.name} ({d.hull_class})
                </option>
              ))}
            </select>
          </div>

          <div className="p2-form-group">
            <label className="p2-label">Quantity (1 - 2,000,000)</label>
            <input
              type="number"
              className="p2-input"
              value={quantity}
              min={1}
              max={2000000}
              onChange={e => setQuantity(Math.max(1, Math.min(2000000, parseInt(e.target.value) || 1)))}
            />
          </div>

          {design && (
            <div className="p2-build-summary">
              <div className="p2-summary-title">Build Summary</div>
              <div className="p2-summary-row">
                <span>Time per ship:</span>
                <span className="p2-summary-value">{formatDuration(design.build_time_seconds)}</span>
              </div>
              <div className="p2-summary-row">
                <span>Total time:</span>
                <span className="p2-summary-value">
                  {formatDuration(Math.ceil(design.build_time_seconds * quantity * (1 - speedBonus / 100)))}
                </span>
              </div>
              {speedBonus > 0 && (
                <div className="p2-summary-row p2-bonus">
                  <span>Speed bonus:</span>
                  <span>-{speedBonus}%</span>
                </div>
              )}
              <div className="p2-summary-divider" />
              <div className="p2-summary-row">
                <span>Total Metal:</span>
                <span className="p2-summary-value">{formatNumber(design.metal_cost * quantity)}</span>
              </div>
              <div className="p2-summary-row">
                <span>Total He3:</span>
                <span className="p2-summary-value">{formatNumber(design.he3_cost * quantity)}</span>
              </div>
              <div className="p2-summary-row">
                <span>Total Gold:</span>
                <span className="p2-summary-value">{formatNumber(design.gold_cost * quantity)}</span>
              </div>
            </div>
          )}

          <LoadingButton
            className="p2-btn p2-btn-primary p2-btn-full"
            onClick={handleBuild}
            disabled={!selectedDesign || quantity < 1}
            loading={submitting}
          >
            Start Production
          </LoadingButton>
        </div>
      </div>
    </div>,
    document.body,
  )
}

export default function ShipFactoryPanel() {
  const { factory, slots, loading, error, build, cancel } = useShipFactory()
  const { designs } = useShipDesigns()
  const { ships, loading: shipsLoading } = useShipInventory()
  const [buildSlot, setBuildSlot] = useState<number | null>(null)

  if (loading) {
    return <div className="p2-panel-loading"><div className="loading-spinner" /><span>Loading Ship Factory...</span></div>
  }

  if (error || !factory) {
    return <div className="p2-panel-error">{error || 'Ship Factory not found. Build one first.'}</div>
  }

  async function handleBuild(designId: string, quantity: number, slot: number) {
    await build({ ship_design_id: designId, quantity, production_slot: slot })
  }

  return (
    <div className="p2-panel">
      <div className="p2-panel-header">
        <div className="p2-panel-icon sf-icon">SF</div>
        <div>
          <div className="p2-panel-title">Ship Factory</div>
          <div className="p2-panel-subtitle">
            Level {factory.level} | Speed Bonus: {factory.speed_bonus_pct}%
          </div>
        </div>
      </div>

      <div className="p2-section">
        <div className="p2-section-title">Production Slots ({factory.production_slots}/5)</div>
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
      </div>

      <div className="p2-section">
        <div className="p2-section-title">Ship Inventory</div>
        {shipsLoading ? (
          <div className="p2-panel-loading"><div className="loading-spinner" /><span>Loading ships...</span></div>
        ) : ships.length === 0 ? (
          <div className="p2-empty-state">No ships constructed yet. Start building ships above.</div>
        ) : (
          <div className="sf-inventory">
            {ships.map(ship => (
              <ShipInventoryCard key={ship.id} ship={ship} />
            ))}
          </div>
        )}
      </div>

      {buildSlot !== null && (
        <BuildModal
          slot={buildSlot}
          designs={designs}
          speedBonus={factory.speed_bonus_pct}
          onBuild={handleBuild}
          onClose={() => setBuildSlot(null)}
        />
      )}
    </div>
  )
}

function SlotCard({ slot, onBuild, onCancel }: {
  slot: ProductionSlot
  onBuild: () => void
  onCancel: () => void
}) {
  // Slot is locked if slot number is beyond unlocked count
  // Note: This requires factory.production_slots to be passed or inferred from context
  // For now, we'll check if in_use is false and there's no data

  if (slot.in_use) {
    return (
      <div className="sf-slot sf-slot-building">
        <div className="sf-slot-number">Slot {slot.slot}</div>
        <div className="sf-slot-design">{slot.design_name || 'Building...'}</div>
        <div className="sf-slot-qty">x{formatNumber(slot.quantity || 0)}</div>
        {slot.build_finish_at && <SlotTimer finishAt={slot.build_finish_at} />}
        <button className="p2-btn p2-btn-danger p2-btn-sm" onClick={onCancel}>Cancel</button>
      </div>
    )
  }

  return (
    <div className="sf-slot sf-slot-empty">
      <div className="sf-slot-number">Slot {slot.slot}</div>
      <button className="p2-btn p2-btn-primary p2-btn-sm" onClick={onBuild}>Build</button>
    </div>
  )
}

function ShipInventoryCard({ ship }: { ship: import('../../types').AvailableShip }) {
  const hullClassColors: Record<string, string> = {
    frigate: '#44aaff',
    cruiser: '#cc8844',
    battleship: '#aa4444',
  }

  const hullColor = hullClassColors[ship.hull_class] || '#888'

  return (
    <div className="sf-inv-card">
      <div className="sf-inv-class" style={{ color: hullColor, borderColor: hullColor }}>
        {ship.hull_class.charAt(0).toUpperCase()}
      </div>
      <div className="sf-inv-info">
        <div className="sf-inv-name">{ship.design_name}</div>
        <div className="sf-inv-hull">{ship.hull_class.charAt(0).toUpperCase() + ship.hull_class.slice(1)}</div>
      </div>
      <div className="sf-inv-qty">
        <div className="sf-inv-qty-label">Qty</div>
        <div className="sf-inv-qty-value">{formatNumber(ship.quantity)}</div>
      </div>
    </div>
  )
}
