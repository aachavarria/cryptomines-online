import { useEffect, useRef } from 'react'
import { Eye, Move, Hammer, Sparkles, Users } from 'lucide-react'
import { useGameContext } from '../../contexts/GameContext.tsx'
import { useBuildings } from '../../hooks/useBuildings.ts'
import { useResources } from '../../hooks/useResources.ts'
import type { BuildingWithType } from '../../types/index.ts'

const ICON_SIZE = 16

export default function BuildingContextMenu() {
  const {
    state,
    deselectAll,
    openDetailPanel,
    openCommandCenterPanel,
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
  const isCommandCenter = selectedBuilding.type_name === 'command_center'

  // Clamp position so menu stays on screen (~140px wide per §7.3).
  const menuWidth = 140
  const baseEntries = 3 // View + Move + Upgrade
  const extraEntries = (isWarehouse ? 1 : 0) + (isCommandCenter ? 1 : 0)
  // Header (~28px) + items (~32px each) + padding.
  const menuHeight = 36 + (baseEntries + extraEntries) * 34 + 16
  const x = Math.min(contextMenuScreenPos.x, window.innerWidth - menuWidth - 10)
  const y = Math.min(contextMenuScreenPos.y, window.innerHeight - menuHeight - 10)

  return (
    <div
      ref={menuRef}
      className="building-context-menu"
      style={{ left: x, top: y, width: menuWidth }}
    >
      <div className="bcm-header">{selectedBuilding.display_name}</div>
      <button
        className="ds-btn ds-btn-ghost ds-btn--block bcm-item"
        onClick={() => openDetailPanel()}
      >
        <Eye size={ICON_SIZE} strokeWidth={1.75} aria-hidden="true" />
        <span>View</span>
      </button>
      <button
        className="ds-btn ds-btn-ghost ds-btn--block bcm-item"
        onClick={() => enterMoveMode(selectedBuilding.id)}
      >
        <Move size={ICON_SIZE} strokeWidth={1.75} aria-hidden="true" />
        <span>Move</span>
      </button>
      <UpgradeButton building={selectedBuilding} onUpgrade={upgrade} />
      {isWarehouse && (
        <button
          className="ds-btn ds-btn-ghost ds-btn--block bcm-item"
          onClick={async () => {
            await collect()
            deselectAll()
          }}
        >
          <Sparkles size={ICON_SIZE} strokeWidth={1.75} aria-hidden="true" />
          <span>Harvest</span>
        </button>
      )}
      {isCommandCenter && (
        <button
          className="ds-btn ds-btn-ghost ds-btn--block bcm-item"
          onClick={() => openCommandCenterPanel()}
        >
          <Users size={ICON_SIZE} strokeWidth={1.75} aria-hidden="true" />
          <span>Recruit</span>
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
    try {
      await onUpgrade(building.id)
      deselectAll()
    } catch {
      // Error is handled by useBuildings hook via dispatch SET_ERROR
    }
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
      className="ds-btn ds-btn-ghost ds-btn--block bcm-item"
      onClick={handleUpgrade}
      disabled={disabled}
    >
      <Hammer size={ICON_SIZE} strokeWidth={1.75} aria-hidden="true" />
      <span>{label}</span>
    </button>
  )
}
