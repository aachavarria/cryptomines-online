import { useState, useEffect } from 'react'
import { createPortal } from 'react-dom'
import { useGameContext } from '../../contexts/GameContext'
import { getInventory, useItem } from '../../services/api'
import type { InventoryItem } from '../../types/inventory'
import './InventoryPanel.css'

interface InventoryPanelProps {
  onClose: () => void
}

export default function InventoryPanel({ onClose }: InventoryPanelProps) {
  const { refreshResources } = useGameContext()
  const [items, setItems] = useState<InventoryItem[]>([])
  const [loading, setLoading] = useState(false)
  const [selectedItem, setSelectedItem] = useState<InventoryItem | null>(null)
  const [hoveredItem, setHoveredItem] = useState<InventoryItem | null>(null)

  const fetchInventory = async () => {
    try {
      const data = await getInventory()
      setItems(data)
    } catch (err) {
      console.error('Failed to fetch inventory:', err)
    }
  }

  useEffect(() => {
    fetchInventory()
  }, [])

  // ESC to close
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose])

  const handleUse = async (item: InventoryItem) => {
    if (!confirm(`Use ${item.display_name}?\n\n${item.description}`)) return

    setLoading(true)
    try {
      const result = await useItem(item.id)

      // Show success notification
      alert(`✓ ${result.effect}`)

      // Refresh inventory + resources
      await fetchInventory()
      await refreshResources()
    } catch (err: any) {
      alert(`✗ Failed to use item: ${err.message || 'Unknown error'}`)
    } finally {
      setLoading(false)
    }
  }

  // Group items by category
  const groupedItems = items.reduce((acc, item) => {
    if (!acc[item.category]) acc[item.category] = []
    acc[item.category].push(item)
    return acc
  }, {} as Record<string, InventoryItem[]>)

  const categoryLabels: Record<string, string> = {
    resource_pack: 'Resource Packs',
    boost: 'Boosts',
    battle: 'Battle Items',
    blueprint: 'Blueprints',
    commander: 'Commanders'
  }

  const categoryDescriptions: Record<string, string> = {
    resource_pack: 'Instant resource grants',
    boost: 'Timed production and construction buffs',
    battle: 'Space Points and protection cards',
    blueprint: 'Unlock new ship blueprints',
    commander: 'Recruit new commanders'
  }

  return createPortal(
    <div className="inventory-backdrop" onClick={e => { if (e.target === e.currentTarget) onClose() }}>
      <div className="inventory-panel">
        <div className="inventory-header">
          <h2>Inventory</h2>
          <div className="inventory-stats">
            <span>{items.reduce((sum, item) => sum + item.quantity, 0)} items</span>
          </div>
          <button className="close-btn" onClick={onClose}>×</button>
        </div>

      {items.length === 0 ? (
        <div className="inventory-empty">
          <p>Your inventory is empty.</p>
          <p className="hint">Complete quests to earn items!</p>
        </div>
      ) : (
        <div className="inventory-categories">
          {Object.entries(groupedItems).map(([category, categoryItems]) => (
            <div key={category} className="inventory-category">
              <div className="category-header">
                <h3>{categoryLabels[category] || category}</h3>
                <span className="category-count">{categoryItems.reduce((sum, item) => sum + item.quantity, 0)}</span>
              </div>
              <p className="category-description">{categoryDescriptions[category]}</p>

              <div className="item-grid">
                {categoryItems.map(item => (
                  <div
                    key={item.id}
                    className={`inventory-item ${selectedItem?.id === item.id ? 'selected' : ''}`}
                    onClick={() => setSelectedItem(item)}
                    onMouseEnter={() => setHoveredItem(item)}
                    onMouseLeave={() => setHoveredItem(null)}
                  >
                    <div className="item-icon">
                      {/* Placeholder icon based on category */}
                      {category === 'resource_pack' && '📦'}
                      {category === 'boost' && '⚡'}
                      {category === 'battle' && '🛡️'}
                      {category === 'blueprint' && '📜'}
                      {category === 'commander' && '👤'}
                    </div>
                    <div className="item-info">
                      <div className="item-name">{item.display_name}</div>
                      <div className="item-quantity">×{item.quantity}</div>
                    </div>

                    {/* Tooltip on hover */}
                    {hoveredItem?.id === item.id && (
                      <div className="item-tooltip">
                        <div className="tooltip-title">{item.display_name}</div>
                        <div className="tooltip-description">{item.description}</div>
                        <div className="tooltip-category">{categoryLabels[item.category]}</div>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Selected Item Detail Panel */}
      {selectedItem && (
        <div className="item-detail">
          <div className="detail-header">
            <h3>{selectedItem.display_name}</h3>
            <button
              className="close-btn"
              onClick={() => setSelectedItem(null)}
            >
              ×
            </button>
          </div>
          <div className="detail-body">
            <div className="detail-icon">
              {selectedItem.category === 'resource_pack' && '📦'}
              {selectedItem.category === 'boost' && '⚡'}
              {selectedItem.category === 'battle' && '🛡️'}
              {selectedItem.category === 'blueprint' && '📜'}
              {selectedItem.category === 'commander' && '👤'}
            </div>
            <p className="detail-description">{selectedItem.description}</p>
            <div className="detail-meta">
              <span>Category: {categoryLabels[selectedItem.category]}</span>
              <span>Quantity: {selectedItem.quantity}</span>
            </div>
          </div>
          <div className="detail-actions">
            <button
              className="use-btn"
              onClick={() => handleUse(selectedItem)}
              disabled={loading}
            >
              {loading ? 'Using...' : 'Use Item'}
            </button>
          </div>
        </div>
      )}
      </div>
    </div>,
    document.body
  )
}
