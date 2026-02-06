import { useGameContext, gridToWorld } from '../../contexts/GameContext.tsx'
import { useCountdown } from '../../hooks/useCountdown.ts'
import type { BuildingWithType } from '../../types/index.ts'

export default function ConstructionInfoPanel() {
  const { state, focusCamera } = useGameContext()

  const upgradingBuildings = state.buildings.filter(b => b.is_upgrading)
  const maxSlots = 2 // Backend maxConstructionSlots = 2

  if (upgradingBuildings.length === 0) return null

  return (
    <div className="construction-info-panel">
      <div className="cip-header">
        Construction ({upgradingBuildings.length}/{maxSlots})
      </div>
      {upgradingBuildings.map(b => (
        <ConstructionEntry
          key={b.id}
          building={b}
          onClick={() => {
            const pos = state.buildingPositions[b.id]
            if (pos) {
              focusCamera(gridToWorld(pos.col, pos.row))
            }
          }}
        />
      ))}
    </div>
  )
}

function ConstructionEntry({ building, onClick }: { building: BuildingWithType; onClick: () => void }) {
  const countdown = useCountdown(building.upgrade_finish_at)

  // Calculate progress
  let progress = 0
  if (building.upgrade_finish_at) {
    const finishTime = new Date(building.upgrade_finish_at).getTime()
    const now = Date.now()
    const totalDuration = finishTime - new Date(building.updated_at).getTime()
    const elapsed = now - new Date(building.updated_at).getTime()
    if (totalDuration > 0) {
      progress = Math.min(100, (elapsed / totalDuration) * 100)
    }
  }

  return (
    <div className="cip-entry" onClick={onClick}>
      <div className="cip-entry-info">
        <span className="cip-entry-name">
          {building.display_name} Lv: {building.level + 1}
        </span>
        <div className="cip-progress">
          <div className="cip-progress-fill" style={{ width: `${progress}%` }} />
        </div>
      </div>
      <span className="cip-entry-timer">{countdown || 'Done!'}</span>
    </div>
  )
}
