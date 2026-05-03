import { useState, useMemo, useEffect } from 'react'
import SceneLighting from './SceneLighting.tsx'
import Starfield from './Starfield.tsx'
import PlanetSurface from './PlanetSurface.tsx'
import BuildingModel from './BuildingModel.tsx'
import IsometricGrid from './IsometricGrid.tsx'
import GhostPreview from './GhostPreview.tsx'
import CameraController from './CameraController.tsx'
import FleetMarker from './FleetMarker.tsx'
import {
  useGameContext,
  getBuildingTiles,
  canPlaceBuilding,
  getBuildingWorldCenter,
  gridToWorld,
  GRID_SIZE,
  type GridPosition,
} from '../../contexts/GameContext.tsx'
import { useBuildings } from '../../hooks/useBuildings.ts'
import { useFleets } from '../../hooks/useFleets.ts'

export default function PlanetScene() {
  const { state, selectBuilding, dispatch, exitPlacementMode } = useGameContext()
  const { construct, move } = useBuildings()
  const { fleets } = useFleets()
  const [cursorGridPos, setCursorGridPos] = useState<GridPosition | null>(null)
  const {
    buildings: allBuildings,
    selectedBuilding,
    hoveredBuilding,
    showConstructPanel,
    showDetailPanel,
    cameraTarget,
    placementMode,
    currentBase,
  } = state

  // Only render buildings that belong to the currently-viewed base
  const buildings = useMemo(
    () => allBuildings.filter(b => (b.base_type === 'space' ? 'space' : 'ground') === currentBase),
    [allBuildings, currentBase],
  )

  // Derive buildingPositions from server-side grid_col/grid_row
  const buildingPositions = useMemo(() => {
    const positions: Record<string, GridPosition> = {}
    for (const b of buildings) {
      positions[b.id] = { col: b.grid_col, row: b.grid_row }
    }
    return positions
  }, [buildings])

  // Compute occupied tile keys (multi-tile aware)
  const occupiedTiles = useMemo(() => {
    const set = new Set<string>()
    for (const [buildingId, pos] of Object.entries(buildingPositions)) {
      const building = buildings.find(b => b.id === buildingId)
      if (building) {
        const tiles = getBuildingTiles(building.type_name, pos.col, pos.row)
        for (const t of tiles) {
          set.add(`${t.col},${t.row}`)
        }
      } else {
        set.add(`${pos.col},${pos.row}`)
      }
    }
    return set
  }, [buildingPositions, buildings])

  // Resolve building type name for both construct and move modes
  const activeBuildingTypeName = useMemo(() => {
    if (placementMode.buildingTypeName) return placementMode.buildingTypeName
    if (placementMode.movingBuildingId) {
      const b = buildings.find(b => b.id === placementMode.movingBuildingId)
      return b?.type_name || null
    }
    return null
  }, [placementMode, buildings])

  // Check if cursor position is a valid placement (accounts for move mode exclusion)
  const isCursorValid = useMemo(() => {
    if (!cursorGridPos || !activeBuildingTypeName) return false
    return canPlaceBuilding(
      activeBuildingTypeName, cursorGridPos.col, cursorGridPos.row, occupiedTiles,
      placementMode.movingBuildingId, buildingPositions, buildings
    )
  }, [cursorGridPos, activeBuildingTypeName, occupiedTiles, placementMode.movingBuildingId, buildingPositions, buildings])

  // Handle tile click during placement mode
  function handleTileClick(pos: GridPosition) {
    if (!placementMode.active) return

    if (placementMode.buildingTypeName) {
      // Validate multi-tile placement
      if (!canPlaceBuilding(placementMode.buildingTypeName, pos.col, pos.row, occupiedTiles)) {
        return  // invalid placement, ignore
      }
      const typeName = placementMode.buildingTypeName
      exitPlacementMode()
      construct(typeName, pos.col, pos.row)
    } else if (placementMode.movingBuildingId) {
      // For move: find the building's type_name
      const movingBuilding = buildings.find(b => b.id === placementMode.movingBuildingId)
      if (movingBuilding) {
        if (!canPlaceBuilding(
          movingBuilding.type_name, pos.col, pos.row, occupiedTiles,
          placementMode.movingBuildingId, buildingPositions, buildings
        )) {
          return
        }
      }
      const buildingId = placementMode.movingBuildingId
      exitPlacementMode()
      move(buildingId, pos.col, pos.row)
    }
  }

  // Handle ESC during placement mode
  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape' && placementMode.active) {
        exitPlacementMode()
      }
    }
    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [placementMode.active, exitPlacementMode])

  const modalOpen = showConstructPanel || showDetailPanel

  // Fleets visible on the space base: filter to current planet (if known) and
  // those not currently traveling away from us. If a fleet has no planet_id
  // (legacy data), still show it so the player isn't left with a blank base.
  const currentPlanetId = typeof window !== 'undefined'
    ? localStorage.getItem('current_planet_id')
    : null
  const stationedHere = useMemo(() => {
    if (currentBase !== 'space') return []
    return fleets.filter(f => {
      if (f.status !== 'stationed') return false
      if (!currentPlanetId) return true
      if (!f.planet_id) return true
      return f.planet_id === currentPlanetId
    })
  }, [fleets, currentBase, currentPlanetId])

  // Lay fleets out along the top edge of the grid, hovering above the platform.
  // Spacing of 2 tiles keeps them legible without colliding with buildings.
  const fleetSlots = useMemo(() => {
    const slots: Array<{ fleet: typeof stationedHere[number]; pos: [number, number, number] }> = []
    const startCol = 1
    const row = 0
    const spacing = 2
    stationedHere.forEach((f, i) => {
      const col = startCol + i * spacing
      if (col >= GRID_SIZE - 1) return
      const [wx, , wz] = gridToWorld(col, row)
      slots.push({ fleet: f, pos: [wx, 2.5, wz] })
    })
    return slots
  }, [stationedHere])

  function openFleetManagement() {
    window.dispatchEvent(
      new CustomEvent('game:navigate', { detail: { path: '/military?tab=fleets' } }),
    )
  }

  return (
    <>
      <SceneLighting />
      <Starfield />
      <PlanetSurface base={currentBase} />

      {/* Isometric grid - visible during placement mode */}
      <IsometricGrid
        visible={placementMode.active}
        occupiedTiles={occupiedTiles}
        buildingTypeName={activeBuildingTypeName}
        cursorGridPos={cursorGridPos}
        onTileClick={handleTileClick}
        onCursorMove={setCursorGridPos}
      />

      {/* Ghost preview during placement */}
      <GhostPreview
        typeName={activeBuildingTypeName || ''}
        gridPosition={cursorGridPos}
        isValid={isCursorValid}
        visible={placementMode.active && !!activeBuildingTypeName && !!cursorGridPos}
      />

      {/* Buildings */}
      {buildings.map((building) => {
        const gridPos = buildingPositions[building.id]
        if (!gridPos) return null
        const worldPos = getBuildingWorldCenter(building.type_name, gridPos.col, gridPos.row)
        const isMoving = placementMode.active && placementMode.movingBuildingId === building.id

        return (
          <BuildingModel
            key={building.id}
            building={building}
            position={worldPos}
            isSelected={selectedBuilding?.id === building.id}
            isHovered={hoveredBuilding === building.id}
            isMoving={isMoving}
            onClick={(screenPos) => {
              if (!placementMode.active) {
                selectBuilding(building, screenPos)
              }
            }}
            onPointerOver={() => dispatch({ type: 'HOVER_BUILDING', payload: building.id })}
            onPointerOut={() => dispatch({ type: 'HOVER_BUILDING', payload: null })}
          />
        )
      })}

      {/* Fleet markers - only on the space base */}
      {currentBase === 'space' && fleetSlots.map(({ fleet, pos }) => (
        <FleetMarker
          key={fleet.id}
          fleet={fleet}
          position={pos}
          onClick={openFleetManagement}
        />
      ))}

      <CameraController
        target={cameraTarget}
        enabled={!modalOpen}
      />
    </>
  )
}
