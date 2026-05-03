import { useState, useEffect } from 'react'
import { createPortal } from 'react-dom'
import { Package, Box, Zap, Shield, FileText, User, X } from 'lucide-react'
import { getInventory, useItem } from '../../services/api'
import { useResources } from '../../hooks/useResources'
import LoadingButton from '../common/LoadingButton'
import type { InventoryItem } from '../../types/inventory'

interface InventoryPanelProps {
  onClose: () => void
}

const CATEGORY_ICONS: Record<InventoryItem['category'], typeof Package> = {
  resource_pack: Package,
  boost: Zap,
  battle: Shield,
  blueprint: FileText,
  commander: User,
}

function CategoryIcon({ category, size = 20 }: { category: InventoryItem['category']; size?: number }) {
  const Icon = CATEGORY_ICONS[category] || Box
  return <Icon size={size} strokeWidth={1.75} />
}

const CATEGORY_LABELS: Record<string, string> = {
  resource_pack: 'Resource Packs',
  boost: 'Boosts',
  battle: 'Battle Items',
  blueprint: 'Blueprints',
  commander: 'Commanders',
}

const CATEGORY_DESCRIPTIONS: Record<string, string> = {
  resource_pack: 'Instant resource grants',
  boost: 'Timed production and construction buffs',
  battle: 'Space Points and protection cards',
  blueprint: 'Unlock new ship blueprints',
  commander: 'Recruit new commanders',
}

export default function InventoryPanel({ onClose }: InventoryPanelProps) {
  const { refreshResources } = useResources()
  const [items, setItems] = useState<InventoryItem[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [selectedItem, setSelectedItem] = useState<InventoryItem | null>(null)

  const fetchInventory = async () => {
    try {
      const data = await getInventory()
      setItems(data)
      setError(null)
    } catch {
      setError('Failed to load inventory')
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
    setLoading(true)
    try {
      await useItem(item.id)
      await fetchInventory()
      await refreshResources()
      setSelectedItem(null)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unknown error'
      setError(`Failed to use item: ${message}`)
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

  const totalCount = items.reduce((sum, item) => sum + item.quantity, 0)

  return createPortal(
    <div
      className="ds-modal-backdrop"
      onClick={e => {
        if (e.target === e.currentTarget) onClose()
      }}
    >
      <div
        className="ds-modal ds-modal--lg inv-modal"
        role="dialog"
        aria-modal="true"
        aria-label="Inventory"
      >
        <div className="ds-modal-header">
          <div className="inv-title">
            <Package size={22} strokeWidth={1.75} />
            <h2 className="ds-modal-title">Inventory</h2>
            <span className="ds-badge ds-badge--neutral inv-count">
              <span className="ds-mono">{totalCount}</span> items
            </span>
          </div>
          <button
            className="ds-btn-icon"
            onClick={onClose}
            aria-label="Close inventory"
          >
            <X size={18} strokeWidth={1.75} />
          </button>
        </div>

        <div className="ds-modal-body inv-body">
          {error && (
            <div className="ds-badge ds-badge--danger inv-error">{error}</div>
          )}

          {items.length === 0 ? (
            <div className="inv-empty">
              <Box size={40} strokeWidth={1.5} className="ds-text-soft" />
              <p>Your inventory is empty.</p>
              <p className="ds-text-muted">Complete quests to earn items.</p>
            </div>
          ) : (
            <div className="inv-layout">
              <div className="inv-categories">
                {Object.entries(groupedItems).map(([category, categoryItems]) => {
                  const cat = category as InventoryItem['category']
                  const subtotal = categoryItems.reduce((s, i) => s + i.quantity, 0)
                  return (
                    <section key={category} className="ds-card inv-category">
                      <header className="inv-cat-header">
                        <div className="inv-cat-title">
                          <CategoryIcon category={cat} size={18} />
                          <h3 className="ds-h3">
                            {CATEGORY_LABELS[category] || category}
                          </h3>
                        </div>
                        <span className="ds-badge ds-badge--teal ds-mono">{subtotal}</span>
                      </header>
                      <p className="ds-text-muted inv-cat-desc">
                        {CATEGORY_DESCRIPTIONS[category]}
                      </p>

                      <div className="inv-grid">
                        {categoryItems.map(item => {
                          const isSelected = selectedItem?.id === item.id
                          return (
                            <button
                              key={item.id}
                              type="button"
                              className={`ds-list-item inv-item ${isSelected ? 'is-selected' : ''}`}
                              aria-selected={isSelected}
                              onClick={() => setSelectedItem(item)}
                            >
                              <span className="inv-item-icon">
                                <CategoryIcon category={cat} size={22} />
                              </span>
                              <span className="inv-item-info">
                                <span className="inv-item-name">{item.display_name}</span>
                                <span className="inv-item-qty ds-mono ds-text-muted">×{item.quantity}</span>
                              </span>
                            </button>
                          )
                        })}
                      </div>
                    </section>
                  )
                })}
              </div>

              {selectedItem && (
                <aside className="ds-card inv-detail">
                  <div className="inv-detail-header">
                    <h3 className="ds-h3">{selectedItem.display_name}</h3>
                    <button
                      className="ds-btn-icon ds-btn-icon--sm"
                      onClick={() => setSelectedItem(null)}
                      aria-label="Close detail"
                    >
                      <X size={16} strokeWidth={1.75} />
                    </button>
                  </div>
                  <div className="inv-detail-body">
                    <div className="inv-detail-icon" aria-hidden="true">
                      <CategoryIcon category={selectedItem.category} size={48} />
                    </div>
                    <p className="inv-detail-desc">{selectedItem.description}</p>
                    <dl className="inv-detail-meta">
                      <div className="inv-detail-row">
                        <dt className="ds-caption">Category</dt>
                        <dd>{CATEGORY_LABELS[selectedItem.category]}</dd>
                      </div>
                      <div className="inv-detail-row">
                        <dt className="ds-caption">Quantity</dt>
                        <dd className="ds-mono">{selectedItem.quantity}</dd>
                      </div>
                    </dl>
                  </div>
                  <div className="inv-detail-actions">
                    <LoadingButton
                      className="ds-btn ds-btn-secondary ds-btn--block"
                      onClick={() => handleUse(selectedItem)}
                      loading={loading}
                    >
                      Use Item
                    </LoadingButton>
                  </div>
                </aside>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Inline-scoped styles via styled tag are not used; using utility classes from primitives.
          Layout-only helpers below are local. */}
      <style>{`
        .inv-modal { width: min(960px, 95vw); }
        .inv-title { display: flex; align-items: center; gap: var(--sp-2); }
        .inv-title h2 { margin: 0; }
        .inv-count { margin-left: var(--sp-2); }
        .inv-body { padding: var(--sp-5) var(--sp-6); }
        .inv-error { display: inline-flex; margin-bottom: var(--sp-3); }
        .inv-empty {
          display: flex; flex-direction: column; align-items: center;
          gap: var(--sp-2); padding: var(--sp-10) var(--sp-4);
          color: var(--ds-text-muted);
        }
        .inv-empty p { margin: 0; }
        .inv-layout {
          display: grid;
          grid-template-columns: minmax(0, 1fr) 320px;
          gap: var(--sp-5);
        }
        @media (max-width: 1024px) {
          .inv-layout { grid-template-columns: minmax(0, 1fr); }
        }
        .inv-categories { display: flex; flex-direction: column; gap: var(--sp-4); min-width: 0; }
        .inv-category { padding: var(--sp-4); }
        .inv-cat-header {
          display: flex; align-items: center; justify-content: space-between;
          gap: var(--sp-3);
        }
        .inv-cat-title {
          display: flex; align-items: center; gap: var(--sp-2);
          color: var(--ds-text);
        }
        .inv-cat-title h3 { margin: 0; }
        .inv-cat-desc { margin: var(--sp-1) 0 var(--sp-3); font-size: var(--fs-sm); }
        .inv-grid {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
          gap: var(--sp-2);
        }
        .inv-item {
          display: flex; align-items: center; gap: var(--sp-3);
          padding: var(--sp-3); width: 100%;
          text-align: left; cursor: pointer;
          border-color: transparent;
        }
        .inv-item-icon {
          width: 36px; height: 36px;
          display: grid; place-items: center;
          background: var(--ds-surface-2);
          border-radius: var(--r-md);
          color: var(--ds-teal-dark);
          flex-shrink: 0;
        }
        .inv-item-info { display: flex; flex-direction: column; min-width: 0; gap: 2px; }
        .inv-item-name {
          font-size: var(--fs-body); font-weight: var(--fw-medium);
          white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
        }
        .inv-item-qty { font-size: var(--fs-caption); }

        .inv-detail {
          align-self: flex-start; position: sticky; top: 0;
          padding: 0; overflow: hidden;
        }
        .inv-detail-header {
          display: flex; align-items: center; justify-content: space-between;
          padding: var(--sp-3) var(--sp-4);
          border-bottom: 1px solid var(--ds-border);
        }
        .inv-detail-header h3 { margin: 0; }
        .inv-detail-body {
          padding: var(--sp-5);
          display: flex; flex-direction: column; align-items: center; gap: var(--sp-3);
          text-align: center;
        }
        .inv-detail-icon {
          width: 72px; height: 72px;
          background: var(--ds-surface-3);
          border-radius: var(--r-lg);
          display: grid; place-items: center;
          color: var(--ds-teal-dark);
        }
        .inv-detail-desc {
          font-size: var(--fs-sm); line-height: var(--lh-normal);
          color: var(--ds-text); margin: 0;
        }
        .inv-detail-meta {
          margin: 0; align-self: stretch;
          display: flex; flex-direction: column; gap: var(--sp-2);
        }
        .inv-detail-row {
          display: flex; justify-content: space-between; align-items: center;
          font-size: var(--fs-sm);
        }
        .inv-detail-row dt { margin: 0; }
        .inv-detail-row dd { margin: 0; font-weight: var(--fw-medium); }
        .inv-detail-actions {
          padding: var(--sp-4);
          border-top: 1px solid var(--ds-border);
          background: var(--ds-surface-2);
        }
      `}</style>
    </div>,
    document.body,
  )
}
