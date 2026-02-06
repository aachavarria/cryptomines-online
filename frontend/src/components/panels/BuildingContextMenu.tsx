import { useEffect, useRef } from 'react'
import { useGameContext } from '../../contexts/GameContext.tsx'
import { useBuildings } from '../../hooks/useBuildings.ts'
import { useResources } from '../../hooks/useResources.ts'
import type { BuildingWithType } from '../../types/index.ts'

export default function BuildingContextMenu() {
  const {
    state,
    deselectAll,
    openDetailPanel,
    enterMoveMode,
  } = useGameContext()
  const { upgrade } = useBuildings()
  const { collect } = useResources()
  const menuRef = useRef<HTMLDivElement>(null)

  const { selectedBuilding, showContextMenu, contextMenuScreenPos } = state

  // Close on click outside
  useEffect(() => {
    if (!showContextMenu) return
    function handleClick(e: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        deselectAll()
      }
    }
    // Delay to avoid the opening click closing it immediately
    const timer = setTimeout(() => {
      window.addEventListener('mousedown', handleClick)
    }, 50)
    return () => {
      clearTimeout(timer)
      window.removeEventListener('mousedown', handleClick)
    }
  }, [showContextMenu, deselectAll])

  // Close on ESC
  useEffect(() => {
    if (!showContextMenu) return
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') deselectAll()
    }
    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [showContextMenu, deselectAll])

  if (!showContextMenu || !selectedBuilding || !contextMenuScreenPos) return null

  const isWarehouse = selectedBuilding.type_name === 'resource_warehouse'

  // Clamp position so menu stays on screen
  const menuWidth = 120
  const menuHeight = isWarehouse ? 180 : 140
  const x = Math.min(contextMenuScreenPos.x, window.innerWidth - menuWidth - 10)
  const y = Math.min(contextMenuScreenPos.y, window.innerHeight - menuHeight - 10)

  return (
    <div
      ref={menuRef}
      className="building-context-menu"
      style={{ left: x, top: y }}
    >
      <div className="bcm-header">{selectedBuilding.display_name}</div>
      <button
        className="bcm-btn"
        onClick={() => openDetailPanel()}
      >
        View
      </button>
      <button
        className="bcm-btn"
        onClick={() => enterMoveMode(selectedBuilding.id)}
      >
        Move
      </button>
      <UpgradeButton building={selectedBuilding} onUpgrade={upgrade} />
      {isWarehouse && (
        <button
          className="bcm-btn bcm-harvest"
          onClick={async () => {
            await collect()
            deselectAll()
          }}
        >
          Harvest
        </button>
      )}
    </div>
  )
}

function UpgradeButton({
  building,
  onUpgrade,
}: {
  building: BuildingWithType
  onUpgrade: (id: string) => Promise<void>
}) {
  const { deselectAll } = useGameContext()
  const isMaxLevel = building.level >= building.max_level
  const isUpgrading = building.is_upgrading

  async function handleUpgrade() {
    if (isMaxLevel || isUpgrading) return
    await onUpgrade(building.id)
    deselectAll()
  }

  let label = 'Upgrade'
  let disabled = false
  if (isMaxLevel) {
    label = 'Max Level'
    disabled = true
  } else if (isUpgrading) {
    label = 'Upgrading...'
    disabled = true
  }

  return (
    <button
      className="bcm-btn"
      onClick={handleUpgrade}
      disabled={disabled}
    >
      {label}
    </button>
  )
}
