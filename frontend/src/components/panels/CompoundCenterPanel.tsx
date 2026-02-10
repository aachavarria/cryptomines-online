import { useState, useEffect } from 'react'
import { useCommanders } from '../../hooks/useCommanders'
import { getInventory, mergeCommander, type Commander } from '../../services/api'
import type { InventoryItem } from '../../types/inventory'
import LoadingButton from '../common/LoadingButton.tsx'
import './CompoundCenterPanel.css'
import '../../styles/common.css'

export default function CompoundCenterPanel() {
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

  return (
    <div className="compound-center-panel">
      <div className="panel-header">
        <h2>Compound Center - Commander Merging</h2>
        <div className="info-text">
          Merge duplicate commanders to increase star rank (max 15 stars)
        </div>
      </div>

      {error && <div className="p2-error-msg">{error}</div>}

      {mergeableCommanders.length === 0 ? (
        <div className="empty-state">
          <p>No commanders available for merging.</p>
          <p className="hint">
            Recruit duplicate commanders at the Command Center to unlock merging.
          </p>
        </div>
      ) : (
        <div className="compound-content">
          {/* Commander Selection List */}
          <div className="commanders-list">
            <h3>Select Commander to Merge</h3>
            <div className="commanders-grid">
              {mergeableCommanders.map(commander => {
                const duplicates = getDuplicateCount(commander)
                return (
                  <div
                    key={commander.id}
                    className={`commander-item ${selectedCommander?.id === commander.id ? 'selected' : ''} ${commander.rarity}`}
                    onClick={() => setSelectedCommander(commander)}
                  >
                    <div className="item-header">
                      <div className="name">{commander.name}</div>
                      <div className="duplicates-badge">
                        {duplicates} duplicate{duplicates > 1 ? 's' : ''}
                      </div>
                    </div>

                    <div className="star-rank">
                      {'★'.repeat(commander.star_rank)}
                      {'☆'.repeat(Math.min(5, 15 - commander.star_rank))}
                      {commander.star_rank > 5 && ` +${commander.star_rank - 5}`}
                    </div>

                    <div className={`rarity-badge ${commander.rarity}`}>
                      {commander.rarity.toUpperCase()}
                    </div>
                  </div>
                )
              })}
            </div>
          </div>

          {/* Merge Panel */}
          {selectedCommander && (
            <div className="merge-panel">
              <h3>Merge {selectedCommander.name}</h3>

              <div className="merge-info">
                <div className="current-rank">
                  <label>Current Star Rank:</label>
                  <div className="rank-display">
                    <div className="stars">
                      {'★'.repeat(selectedCommander.star_rank)}
                      {'☆'.repeat(15 - selectedCommander.star_rank)}
                    </div>
                    <div className="rank-number">{selectedCommander.star_rank} / 15</div>
                  </div>
                  <div className="bonus-text">{getStarBonus(selectedCommander.star_rank)}</div>
                </div>

                <div className="merge-arrow">→</div>

                <div className="new-rank">
                  <label>New Star Rank:</label>
                  <div className="rank-display">
                    <div className="stars">
                      {'★'.repeat(Math.min(selectedCommander.star_rank + mergeQuantity, 15))}
                      {'☆'.repeat(Math.max(0, 15 - selectedCommander.star_rank - mergeQuantity))}
                    </div>
                    <div className="rank-number">
                      {Math.min(selectedCommander.star_rank + mergeQuantity, 15)} / 15
                    </div>
                  </div>
                  <div className="bonus-text">
                    {getStarBonus(Math.min(selectedCommander.star_rank + mergeQuantity, 15))}
                  </div>
                </div>
              </div>

              <div className="merge-controls">
                <div className="quantity-control">
                  <label>Duplicates to Use:</label>
                  <div className="quantity-input">
                    <button
                      onClick={() => setMergeQuantity(Math.max(1, mergeQuantity - 1))}
                      disabled={mergeQuantity <= 1}
                    >
                      -
                    </button>
                    <input
                      type="number"
                      min="1"
                      max={Math.min(getDuplicateCount(selectedCommander), 15 - selectedCommander.star_rank)}
                      value={mergeQuantity}
                      onChange={e => setMergeQuantity(Math.max(1, parseInt(e.target.value) || 1))}
                    />
                    <button
                      onClick={() => setMergeQuantity(Math.min(
                        getDuplicateCount(selectedCommander),
                        15 - selectedCommander.star_rank,
                        mergeQuantity + 1
                      ))}
                      disabled={mergeQuantity >= getDuplicateCount(selectedCommander) || selectedCommander.star_rank + mergeQuantity >= 15}
                    >
                      +
                    </button>
                  </div>
                  <div className="available-text">
                    {getDuplicateCount(selectedCommander)} available
                  </div>
                </div>

                <LoadingButton
                  className="merge-btn"
                  onClick={handleMerge}
                  loading={merging}
                  disabled={selectedCommander.star_rank >= 15 || getDuplicateCount(selectedCommander) < 1}
                >
                  Merge {mergeQuantity} Duplicate{mergeQuantity > 1 ? 's' : ''}
                </LoadingButton>
              </div>

              {selectedCommander.star_rank >= 15 && (
                <div className="max-rank-notice">
                  ★ This commander is already at maximum star rank (15) ★
                </div>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  )
}
