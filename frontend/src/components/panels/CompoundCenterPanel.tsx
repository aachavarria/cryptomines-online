import { useState, useEffect } from 'react'
import { ArrowRight, Plus, Minus, Sparkles, X } from 'lucide-react'
import { useCommanders } from '../../hooks/useCommanders'
import { getInventory, mergeCommander, type Commander } from '../../services/api'
import type { InventoryItem } from '../../types/inventory'
import LoadingButton from '../common/LoadingButton.tsx'

// Map backend rarity → design-system rarity badge variant.
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

interface CompoundCenterPanelProps {
  onClose?: () => void
}

export default function CompoundCenterPanel({ onClose }: CompoundCenterPanelProps = {}) {
  const { commanders, refresh: refreshCommanders } = useCommanders()
  const [inventory, setInventory] = useState<InventoryItem[]>([])
  const [selectedCommander, setSelectedCommander] = useState<Commander | null>(null)
  const [mergeQuantity, setMergeQuantity] = useState(1)
  const [merging, setMerging] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Fetch inventory to get commander cards
  const fetchInventory = async () => {
    try {
      const items = await getInventory()
      setInventory(items.filter(item => item.category === 'commander'))
      setError(null)
    } catch {
      setError('Failed to load commander cards from inventory')
    }
  }

  useEffect(() => {
    fetchInventory()
  }, [])

  // Find duplicate cards for selected commander
  const getDuplicateCount = (commander: Commander): number => {
    const commanderItemKey = `commander_${commander.name}`
    const item = inventory.find(i => i.item_key === commanderItemKey)
    return item?.quantity || 0
  }

  // Commanders with duplicates available for merging
  const mergeableCommanders = commanders.filter(c => getDuplicateCount(c) > 0 && c.star_rank < 15)

  const handleMerge = async () => {
    if (!selectedCommander) return

    const duplicates = getDuplicateCount(selectedCommander)
    if (mergeQuantity > duplicates) {
      alert(`Not enough duplicates (have ${duplicates}, need ${mergeQuantity})`)
      return
    }

    const newRank = Math.min(selectedCommander.star_rank + mergeQuantity, 15)
    if (!confirm(`Merge ${mergeQuantity} duplicate(s) to increase star rank from ${selectedCommander.star_rank} to ${newRank}?`)) return

    setMerging(true)
    try {
      const result = await mergeCommander(selectedCommander.id, mergeQuantity)
      alert(`✓ ${result.message}`)

      // Refresh data
      await refreshCommanders()
      await fetchInventory()

      // Update selected commander
      const updated = commanders.find(c => c.id === selectedCommander.id)
      if (updated) {
        setSelectedCommander(updated)
      }

      setMergeQuantity(1)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unknown error'
      setError(`Merge failed: ${message}`)
    } finally {
      setMerging(false)
    }
  }

  const getStarBonus = (stars: number): string => {
    // Each star grants +2% to all stats (simplified)
    const bonus = stars * 2
    return `+${bonus}% to all stats`
  }

  const handleBackdropClick = (e: React.MouseEvent<HTMLDivElement>) => {
    if (onClose && e.target === e.currentTarget) onClose()
  }

  return (
    <div
      className="ds-modal-backdrop"
      role="dialog"
      aria-modal="true"
      onClick={handleBackdropClick}
    >
      <div className="ds-modal ds-modal--lg">
        <div className="ds-modal-header">
          <div>
            <h2 className="ds-modal-title">Compound Center</h2>
            <div
              className="ds-text-muted"
              style={{ fontSize: 'var(--fs-sm)', marginTop: 2 }}
            >
              Merge duplicate commanders to increase star rank (max <span className="ds-mono">15</span>).
            </div>
          </div>
          {onClose && (
            <button
              type="button"
              className="ds-btn-icon"
              aria-label="Close"
              onClick={onClose}
            >
              <X size={18} strokeWidth={2} />
            </button>
          )}
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

          {mergeableCommanders.length === 0 ? (
            <div
              style={{
                textAlign: 'center',
                padding: 'var(--sp-8) var(--sp-4)',
                color: 'var(--ds-text-muted)',
              }}
            >
              <p style={{ margin: 0 }}>No commanders available for merging.</p>
              <p
                className="ds-text-soft"
                style={{ marginTop: 'var(--sp-2)', fontSize: 'var(--fs-sm)' }}
              >
                Recruit duplicate commanders at the Command Center to unlock merging.
              </p>
            </div>
          ) : (
            <div
              style={{
                display: 'grid',
                gridTemplateColumns: 'minmax(0, 1fr) minmax(0, 1.2fr)',
                gap: 'var(--sp-5)',
              }}
            >
              {/* Commander Selection List */}
              <div>
                <div
                  className="ds-caption"
                  style={{ marginBottom: 'var(--sp-3)' }}
                >
                  Select Commander
                </div>
                <div
                  style={{
                    display: 'flex',
                    flexDirection: 'column',
                    gap: 'var(--sp-2)',
                    maxHeight: 480,
                    overflowY: 'auto',
                    paddingRight: 'var(--sp-1)',
                  }}
                >
                  {mergeableCommanders.map(commander => {
                    const duplicates = getDuplicateCount(commander)
                    const isSelected = selectedCommander?.id === commander.id
                    return (
                      <div
                        key={commander.id}
                        role="button"
                        tabIndex={0}
                        aria-selected={isSelected}
                        className="ds-list-item"
                        onClick={() => setSelectedCommander(commander)}
                        onKeyDown={(e) => {
                          if (e.key === 'Enter' || e.key === ' ') {
                            e.preventDefault()
                            setSelectedCommander(commander)
                          }
                        }}
                      >
                        <div className="ds-row--between" style={{ marginBottom: 6 }}>
                          <span style={{ fontWeight: 600, color: 'var(--ds-text)' }}>
                            {commander.name}
                          </span>
                          <span className="ds-badge ds-badge--teal">
                            <span className="ds-mono">{duplicates}</span>
                            {' '}dupe{duplicates > 1 ? 's' : ''}
                          </span>
                        </div>

                        <div
                          className="ds-mono"
                          style={{
                            color: 'var(--ds-orange)',
                            fontSize: 'var(--fs-sm)',
                            letterSpacing: '0.1em',
                            marginBottom: 6,
                          }}
                        >
                          {'★'.repeat(commander.star_rank)}
                          <span style={{ color: 'var(--ds-text-soft)' }}>
                            {'☆'.repeat(Math.min(5, 15 - commander.star_rank))}
                          </span>
                          {commander.star_rank > 5 && (
                            <span
                              className="ds-text-muted"
                              style={{ marginLeft: 6, letterSpacing: 0 }}
                            >
                              +{commander.star_rank - 5}
                            </span>
                          )}
                        </div>

                        <span className={`ds-badge ds-badge--rarity-${rarityVariant(commander.rarity)}`}>
                          {commander.rarity.toUpperCase()}
                        </span>
                      </div>
                    )
                  })}
                </div>
              </div>

              {/* Merge Panel */}
              {selectedCommander ? (
                <div className="ds-card ds-stack">
                  <h3 className="ds-h3" style={{ margin: 0, textAlign: 'center' }}>
                    Merge {selectedCommander.name}
                  </h3>

                  <div
                    style={{
                      display: 'grid',
                      gridTemplateColumns: '1fr auto 1fr',
                      gap: 'var(--sp-3)',
                      alignItems: 'center',
                    }}
                  >
                    <div className="ds-card" style={{ background: 'var(--ds-surface)' }}>
                      <div className="ds-caption" style={{ marginBottom: 'var(--sp-2)' }}>
                        Current
                      </div>
                      <div
                        className="ds-mono"
                        style={{
                          color: 'var(--ds-orange)',
                          fontSize: 'var(--fs-sm)',
                          letterSpacing: '0.12em',
                          textAlign: 'center',
                          marginBottom: 4,
                        }}
                      >
                        {'★'.repeat(selectedCommander.star_rank)}
                        <span style={{ color: 'var(--ds-text-soft)' }}>
                          {'☆'.repeat(15 - selectedCommander.star_rank)}
                        </span>
                      </div>
                      <div
                        className="ds-mono"
                        style={{ textAlign: 'center', fontWeight: 600 }}
                      >
                        {selectedCommander.star_rank} / 15
                      </div>
                      <div
                        className="ds-text-muted"
                        style={{
                          textAlign: 'center',
                          fontSize: 'var(--fs-caption)',
                          marginTop: 4,
                          fontStyle: 'italic',
                        }}
                      >
                        {getStarBonus(selectedCommander.star_rank)}
                      </div>
                    </div>

                    <ArrowRight
                      size={22}
                      strokeWidth={2}
                      color="var(--ds-teal)"
                      aria-hidden="true"
                    />

                    <div
                      className="ds-card"
                      style={{
                        background: 'var(--ds-teal-tint)',
                        borderColor: 'var(--ds-teal)',
                      }}
                    >
                      <div className="ds-caption" style={{ marginBottom: 'var(--sp-2)' }}>
                        New
                      </div>
                      <div
                        className="ds-mono"
                        style={{
                          color: 'var(--ds-orange)',
                          fontSize: 'var(--fs-sm)',
                          letterSpacing: '0.12em',
                          textAlign: 'center',
                          marginBottom: 4,
                        }}
                      >
                        {'★'.repeat(Math.min(selectedCommander.star_rank + mergeQuantity, 15))}
                        <span style={{ color: 'var(--ds-text-soft)' }}>
                          {'☆'.repeat(Math.max(0, 15 - selectedCommander.star_rank - mergeQuantity))}
                        </span>
                      </div>
                      <div
                        className="ds-mono"
                        style={{ textAlign: 'center', fontWeight: 600 }}
                      >
                        {Math.min(selectedCommander.star_rank + mergeQuantity, 15)} / 15
                      </div>
                      <div
                        className="ds-text-muted"
                        style={{
                          textAlign: 'center',
                          fontSize: 'var(--fs-caption)',
                          marginTop: 4,
                          fontStyle: 'italic',
                        }}
                      >
                        {getStarBonus(Math.min(selectedCommander.star_rank + mergeQuantity, 15))}
                      </div>
                    </div>
                  </div>

                  <div className="ds-stack">
                    <label
                      className="ds-caption"
                      htmlFor="merge-quantity-input"
                    >
                      Duplicates to use
                    </label>
                    <div className="ds-row" style={{ gap: 'var(--sp-2)' }}>
                      <button
                        type="button"
                        className="ds-btn-icon"
                        aria-label="Decrease"
                        onClick={() => setMergeQuantity(Math.max(1, mergeQuantity - 1))}
                        disabled={mergeQuantity <= 1}
                      >
                        <Minus size={14} strokeWidth={2} />
                      </button>
                      <input
                        id="merge-quantity-input"
                        className="ds-input ds-mono"
                        type="number"
                        min={1}
                        max={Math.min(getDuplicateCount(selectedCommander), 15 - selectedCommander.star_rank)}
                        value={mergeQuantity}
                        onChange={e => setMergeQuantity(Math.max(1, parseInt(e.target.value) || 1))}
                        style={{ textAlign: 'center', flex: 1 }}
                      />
                      <button
                        type="button"
                        className="ds-btn-icon"
                        aria-label="Increase"
                        onClick={() => setMergeQuantity(Math.min(
                          getDuplicateCount(selectedCommander),
                          15 - selectedCommander.star_rank,
                          mergeQuantity + 1
                        ))}
                        disabled={mergeQuantity >= getDuplicateCount(selectedCommander) || selectedCommander.star_rank + mergeQuantity >= 15}
                      >
                        <Plus size={14} strokeWidth={2} />
                      </button>
                    </div>
                    <div
                      className="ds-text-muted"
                      style={{
                        textAlign: 'center',
                        fontSize: 'var(--fs-caption)',
                      }}
                    >
                      <span className="ds-mono">{getDuplicateCount(selectedCommander)}</span> available
                    </div>

                    <LoadingButton
                      className="ds-btn-secondary ds-btn--block"
                      onClick={handleMerge}
                      loading={merging}
                      disabled={selectedCommander.star_rank >= 15 || getDuplicateCount(selectedCommander) < 1}
                    >
                      Compound <span className="ds-mono">{mergeQuantity}</span>
                      {' '}Duplicate{mergeQuantity > 1 ? 's' : ''}
                    </LoadingButton>
                  </div>

                  {selectedCommander.star_rank >= 15 && (
                    <div
                      className="ds-row"
                      style={{
                        gap: 'var(--sp-2)',
                        padding: 'var(--sp-3)',
                        background: 'var(--ds-orange-tint)',
                        color: 'var(--ds-orange-strong)',
                        borderRadius: 'var(--r-md)',
                        justifyContent: 'center',
                        fontWeight: 600,
                        fontSize: 'var(--fs-sm)',
                      }}
                    >
                      <Sparkles size={16} strokeWidth={2} />
                      Maximum star rank (15) reached
                    </div>
                  )}
                </div>
              ) : (
                <div
                  className="ds-card"
                  style={{
                    display: 'grid',
                    placeItems: 'center',
                    color: 'var(--ds-text-muted)',
                    minHeight: 200,
                    fontSize: 'var(--fs-sm)',
                  }}
                >
                  Select a commander to begin compounding.
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
