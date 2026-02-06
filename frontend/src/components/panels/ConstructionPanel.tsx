import { useState, useEffect } from 'react'
import { useGameContext, CATEGORY_COLORS, BUILDING_ABBREVIATIONS } from '../../contexts/GameContext.tsx'
import { formatNumber } from '../../hooks/useCountdown.ts'

const CATEGORIES = [
  { id: 'resource', label: 'Resource' },
  { id: 'core', label: 'Core' },
  { id: 'military', label: 'Military' },
  { id: 'space', label: 'Space' },
]

export default function ConstructionPanel() {
  const { state, closeConstructPanel, enterPlacementMode } = useGameContext()
  const [activeCategory, setActiveCategory] = useState('resource')

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

  const filteredTypes = state.buildingTypes.filter(bt => {
    if (activeCategory === 'space') {
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
          <span className="cm-title">Build New Structure</span>
          <button className="cm-close" onClick={closeConstructPanel}>X</button>
        </div>

        <div className="cm-tabs">
          {CATEGORIES.map(cat => (
            <button
              key={cat.id}
              className={`cm-tab ${cat.id} ${activeCategory === cat.id ? 'active' : ''}`}
              onClick={() => setActiveCategory(cat.id)}
            >
              {cat.label}
            </button>
          ))}
        </div>

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
