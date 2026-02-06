import { useState, useMemo, useEffect } from 'react'
import SceneLighting from './SceneLighting.tsx'
import Starfield from './Starfield.tsx'
import PlanetSurface from './PlanetSurface.tsx'
import BuildingModel from './BuildingModel.tsx'
import IsometricGrid from './IsometricGrid.tsx'
import GhostPreview from './GhostPreview.tsx'
import CameraController from './CameraController.tsx'
import {
  useGameContext,
  autoAssignPositions,
  getBuildingTiles,
  canPlaceBuilding,
  getBuildingWorldCenter,
  type GridPosition,
} from '../../contexts/GameContext.tsx'
import { useBuildings } from '../../hooks/useBuildings.ts'

export default function PlanetScene() {
  const { state, selectBuilding, dispatch, placeBuilding, exitPlacementMode } = useGameContext()
  const { construct } = useBuildings()
  const [cursorGridPos, setCursorGridPos] = useState<GridPosition | null>(null)
  const {
    buildings,
    selectedBuilding,
    hoveredBuilding,
    showConstructPanel,
    showDetailPanel,
    cameraTarget,
    placementMode,
    buildingPositions,
  } = state

  // Auto-assign positions for buildings that don't have one
  useEffect(() => {
    if (buildings.length === 0) return
    const newPositions = autoAssignPositions(buildings, buildingPositions)
    // Only update if there are new assignments
    const hasNew = buildings.some(b => !buildingPositions[b.id] && newPositions[b.id])
    if (hasNew) {
      dispatch({ type: 'SET_BUILDING_POSITIONS', payload: newPositions })
    }
  }, [buildings, buildingPositions, dispatch])

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
      construct(typeName)
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
      placeBuilding(placementMode.movingBuildingId, pos)
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

  return (
    <>
      <SceneLighting />
      <Starfield />
      <PlanetSurface />

      {/* Isometric grid - visible during placement mode */}
      <IsometricGrid
        visible={placementMode.active}
        occupiedTiles={occupiedTiles}
        buildingTypeName={placementMode.buildingTypeName || null}
        cursorGridPos={cursorGridPos}
        onTileClick={handleTileClick}
        onCursorMove={setCursorGridPos}
      />

      {/* Ghost preview during placement */}
      <GhostPreview
        typeName={placementMode.buildingTypeName || ''}
        gridPosition={cursorGridPos}
        isValid={
          cursorGridPos && placementMode.buildingTypeName
            ? canPlaceBuilding(placementMode.buildingTypeName, cursorGridPos.col, cursorGridPos.row, occupiedTiles)
            : false
        }
        visible={placementMode.active && !!placementMode.buildingTypeName && !!cursorGridPos}
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

      <CameraController
        target={cameraTarget}
        enabled={!modalOpen}
      />
    </>
  )
}
