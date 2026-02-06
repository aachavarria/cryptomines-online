import { BUILDING_TYPES } from '../types'
import type { BuildingWithType } from '../types'

interface ConstructPanelProps {
  buildings: BuildingWithType[]
  onConstruct: (buildingType: string) => void
  constructing: string | null
}

export default function ConstructPanel({ buildings, onConstruct, constructing }: ConstructPanelProps) {
  // Find civic center level
  const civicCenter = buildings.find((b) => b.type_name === 'civic_center')
  const ccLevel = civicCenter?.level ?? 0

  // Group existing buildings by type_name to show count
  const buildingCounts = buildings.reduce<Record<string, number>>((acc, b) => {
    acc[b.type_name] = (acc[b.type_name] || 0) + 1
    return acc
  }, {})

  return (
    <div className="construct-panel">
      <h3 className="panel-title">Construct New Building</h3>
      <p className="panel-subtitle">Civic Center Lv {ccLevel}</p>
      <div className="construct-grid">
        {BUILDING_TYPES.map((bt) => {
          const count = buildingCounts[bt.name] || 0
          return (
            <div key={bt.name} className="construct-card">
              <div className="construct-info">
                <span className="construct-name">{bt.display_name}</span>
                <span className="construct-category">{bt.category}</span>
                {count > 0 && <span className="construct-count">Built: {count}</span>}
              </div>
              <button
                className="btn-construct"
                onClick={() => onConstruct(bt.name)}
                disabled={constructing === bt.name}
              >
                {constructing === bt.name ? 'Building...' : 'Build'}
              </button>
            </div>
          )
        })}
      </div>
    </div>
  )
}
