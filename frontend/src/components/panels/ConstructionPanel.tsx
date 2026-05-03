import { useState, useEffect } from 'react'
import { X, Hammer } from 'lucide-react'
import { useGameContext, BUILDING_ABBREVIATIONS } from '../../contexts/GameContext.tsx'
import { formatNumber } from '../../hooks/useCountdown.ts'

const GROUND_CATEGORIES = [
  { id: 'resource', label: 'Resource' },
  { id: 'core', label: 'Core' },
  { id: 'military', label: 'Military' },
]

// Space base has no Station tab — Space Station is auto-created at planet
// creation, and the rest are defenses (+ Celestial Base, an end-game structure
// shown together with defenses).
const SPACE_CATEGORIES = [
  { id: 'defense', label: 'Defense' },
]

// Buildings auto-created at planet creation — never shown in the Build panel.
const AUTO_CREATED = new Set(['civic_center', 'space_station'])

// Returns 'space' or 'ground' for a building type, falling back to category
// when the backend hasn't been restarted to expose `base`.
function resolveBase(bt: { base?: string; category: string }): 'space' | 'ground' {
  if (bt.base === 'space') return 'space'
  if (bt.base === 'ground') return 'ground'
  if (bt.category === 'space' || bt.category === 'defense') return 'space'
  return 'ground'
}

export default function ConstructionPanel() {
  const { state, closeConstructPanel, enterPlacementMode } = useGameContext()
  const categories = state.currentBase === 'space' ? SPACE_CATEGORIES : GROUND_CATEGORIES
  const [activeCategory, setActiveCategory] = useState(categories[0].id)

  // Reset active category when base changes
  useEffect(() => {
    setActiveCategory(categories[0].id)
  }, [state.currentBase, categories])

  // Close on Escape
  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') closeConstructPanel()
    }
    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [closeConstructPanel])

  // Count existing buildings per type
  const buildingCounts: Record<string, number> = {}
  for (const b of state.buildings) {
    buildingCounts[b.type_name] = (buildingCounts[b.type_name] || 0) + 1
  }

  // Buildings whose mechanics are out of scope per docs/planning/final-scope.md.
  // Hidden from construction so players don't waste resources on no-op buildings.
  const HIDDEN = new Set(['trading_center', 'galaxy_transporter', 'compound_center'])

  const filteredTypes = state.buildingTypes.filter(bt => {
    if (HIDDEN.has(bt.name)) return false
    if (AUTO_CREATED.has(bt.name)) return false
    if (resolveBase(bt) !== state.currentBase) return false
    if (state.currentBase === 'space') {
      // In space view show all buildable space-base structures (Celestial
      // Base + 4 defenses) regardless of the active tab — there is only
      // one tab.
      return bt.category === 'space' || bt.category === 'defense'
    }
    return bt.category === activeCategory
  })

  // Per-category total counts (for the (built/typesShown) tab badge)
  function categoryCount(catId: string): { built: number; types: number } {
    let built = 0
    let types = 0
    for (const bt of state.buildingTypes) {
      if (HIDDEN.has(bt.name)) continue
      if (AUTO_CREATED.has(bt.name)) continue
      if (resolveBase(bt) !== state.currentBase) continue
      if (bt.category !== catId) continue
      types += 1
      built += buildingCounts[bt.name] || 0
    }
    return { built, types }
  }

  function handleBuild(typeName: string) {
    // Enter placement mode - user will click a tile on the isometric grid
    // PlanetScene.handleTileClick will call construct() when a tile is clicked
    enterPlacementMode(typeName)
  }

  return (
    <div
      className="ds-modal-backdrop"
      role="dialog"
      aria-modal="true"
      onClick={(e) => {
        if (e.target === e.currentTarget) closeConstructPanel()
      }}
    >
      <div className="ds-modal ds-modal--lg" style={{ maxWidth: 'min(760px, 95vw)' }}>
        <div className="ds-modal-header">
          <h2 className="ds-modal-title">
            Build on {state.currentBase === 'space' ? 'Space Base' : 'Ground Base'}
          </h2>
          <button
            type="button"
            className="ds-btn-icon"
            aria-label="Close"
            onClick={closeConstructPanel}
          >
            <X size={18} strokeWidth={2} />
          </button>
        </div>

        {state.currentBase === 'ground' && (
          <div className="ds-tabs" style={{ padding: '0 var(--sp-6)' }}>
            {categories.map(cat => {
              const { built, types } = categoryCount(cat.id)
              const selected = activeCategory === cat.id
              return (
                <button
                  key={cat.id}
                  type="button"
                  role="tab"
                  aria-selected={selected}
                  className="ds-tab"
                  onClick={() => setActiveCategory(cat.id)}
                >
                  {cat.label}
                  <span
                    className="ds-mono ds-text-muted"
                    style={{ marginLeft: 6, fontSize: 'var(--fs-caption)' }}
                  >
                    ({built}/{types})
                  </span>
                </button>
              )
            })}
          </div>
        )}

        <div className="ds-modal-body">
          <div
            style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(2, minmax(0, 1fr))',
              gap: 'var(--sp-4)',
            }}
          >
            {filteredTypes.map(bt => {
              const count = buildingCounts[bt.name] || 0
              const maxCount = bt.max_count_per_planet
              const isMaxed = count >= maxCount
              const abbr = BUILDING_ABBREVIATIONS[bt.name] || '??'
              return (
                <div
                  key={bt.name}
                  className="ds-card"
                  style={{
                    display: 'flex',
                    flexDirection: 'column',
                    gap: 'var(--sp-3)',
                    opacity: isMaxed ? 0.6 : 1,
                  }}
                >
                  <div className="ds-row" style={{ gap: 'var(--sp-3)' }}>
                    <div
                      style={{
                        width: 44,
                        height: 44,
                        borderRadius: 'var(--r-md)',
                        background: 'var(--ds-surface-3)',
                        border: '1px solid var(--ds-border)',
                        color: 'var(--ds-text-muted)',
                        display: 'grid',
                        placeItems: 'center',
                        fontWeight: 600,
                        fontSize: 13,
                        letterSpacing: '0.04em',
                        flexShrink: 0,
                      }}
                      aria-hidden="true"
                    >
                      {abbr}
                    </div>
                    <div style={{ flex: 1, minWidth: 0 }}>
                      <div
                        style={{
                          fontSize: 'var(--fs-h3)',
                          fontWeight: 600,
                          color: 'var(--ds-text)',
                          lineHeight: 'var(--lh-tight)',
                        }}
                      >
                        {bt.display_name}
                      </div>
                      <div
                        className="ds-text-muted"
                        style={{ fontSize: 'var(--fs-sm)', marginTop: 2 }}
                      >
                        Built <span className="ds-mono">{count}</span>/<span className="ds-mono">{maxCount}</span>
                      </div>
                    </div>
                  </div>

                  <div className="ds-row" style={{ gap: 'var(--sp-3)', flexWrap: 'wrap' }}>
                    <span className="ds-row" style={{ gap: 6 }}>
                      <span className="ds-resource-dot ds-resource-dot--metal" />
                      <span className="ds-mono">{formatNumber(bt.base_cost_metal)}</span>
                    </span>
                    <span className="ds-row" style={{ gap: 6 }}>
                      <span className="ds-resource-dot ds-resource-dot--he3" />
                      <span className="ds-mono">{formatNumber(bt.base_cost_he3)}</span>
                    </span>
                    <span className="ds-row" style={{ gap: 6 }}>
                      <span className="ds-resource-dot ds-resource-dot--gold" />
                      <span className="ds-mono">{formatNumber(bt.base_cost_gold)}</span>
                    </span>
                  </div>

                  <button
                    type="button"
                    className="ds-btn-secondary ds-btn--block"
                    onClick={() => handleBuild(bt.name)}
                    disabled={isMaxed}
                  >
                    <Hammer size={14} strokeWidth={2} />
                    {isMaxed ? 'Max Built' : 'Build'}
                  </button>
                </div>
              )
            })}
          </div>
        </div>
      </div>
    </div>
  )
}
