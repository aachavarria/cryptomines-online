import { useState, useEffect } from 'react'
import { useGameContext, CATEGORY_COLORS, BUILDING_ABBREVIATIONS } from '../../contexts/GameContext.tsx'
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

  function handleBuild(typeName: string) {
    // Enter placement mode - user will click a tile on the isometric grid
    // PlanetScene.handleTileClick will call construct() when a tile is clicked
    enterPlacementMode(typeName)
  }

  return (
    <div
      className="construction-modal-backdrop"
      onClick={(e) => {
        if (e.target === e.currentTarget) closeConstructPanel()
      }}
    >
      <div className="construction-modal">
        <div className="cm-header">
          <span className="cm-title">
            Build on {state.currentBase === 'space' ? 'Space Base' : 'Ground Base'}
          </span>
          <button className="cm-close" onClick={closeConstructPanel}>X</button>
        </div>

        {state.currentBase === 'ground' && (
          <div className="cm-tabs">
            {categories.map(cat => (
              <button
                key={cat.id}
                className={`cm-tab ${cat.id} ${activeCategory === cat.id ? 'active' : ''}`}
                onClick={() => setActiveCategory(cat.id)}
              >
                {cat.label}
              </button>
            ))}
          </div>
        )}

        <div className="cm-grid">
          {filteredTypes.map(bt => {
            const count = buildingCounts[bt.name] || 0
            const maxCount = bt.max_count_per_planet
            const isMaxed = count >= maxCount
            const color = CATEGORY_COLORS[bt.category] || '#888'
            const abbr = BUILDING_ABBREVIATIONS[bt.name] || '??'
            return (
              <div key={bt.name} className={`cm-card ${isMaxed ? 'dimmed' : ''}`}>
                <div className="cm-card-header">
                  <div
                    className="cm-card-icon"
                    style={{ background: color, color: '#000' }}
                  >
                    {abbr}
                  </div>
                  <div className="cm-card-info">
                    <div className="cm-card-name">{bt.display_name}</div>
                    <div className="cm-card-count">
                      Built: {count}/{maxCount}
                    </div>
                  </div>
                </div>

                <div className="cm-card-costs">
                  <span className="cm-card-cost">
                    <span className="cm-card-cost-dot metal" />
                    {formatNumber(bt.base_cost_metal)}
                  </span>
                  <span className="cm-card-cost">
                    <span className="cm-card-cost-dot he3" />
                    {formatNumber(bt.base_cost_he3)}
                  </span>
                  <span className="cm-card-cost">
                    <span className="cm-card-cost-dot gold" />
                    {formatNumber(bt.base_cost_gold)}
                  </span>
                </div>

                <button
                  className="cm-card-build-btn"
                  onClick={() => handleBuild(bt.name)}
                  disabled={isMaxed}
                >
                  {isMaxed ? 'MAX BUILT' : 'BUILD'}
                </button>
              </div>
            )
          })}
        </div>
      </div>
    </div>
  )
}
