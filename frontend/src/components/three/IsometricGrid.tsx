import { useMemo, useCallback } from 'react'
import * as THREE from 'three'
import type { ThreeEvent } from '@react-three/fiber'
import {
  GRID_SIZE,
  TILE_WORLD_SIZE,
  gridToWorld,
  worldToGrid,
  getBuildingTiles,
  canPlaceBuilding,
  type GridPosition,
} from '../../contexts/GameContext.tsx'

interface IsometricGridProps {
  visible: boolean
  occupiedTiles: Set<string>
  buildingTypeName: string | null
  cursorGridPos: GridPosition | null
  onTileClick: (pos: GridPosition) => void
  onCursorMove: (pos: GridPosition | null) => void
}

// Half-size of each square tile
const HALF = TILE_WORLD_SIZE / 2

export default function IsometricGrid({
  visible,
  occupiedTiles,
  buildingTypeName,
  cursorGridPos,
  onTileClick,
  onCursorMove,
}: IsometricGridProps) {
  // 1. Generate wireframe line segments geometry (square tiles)
  const gridLinesGeometry = useMemo(() => {
    const points: number[] = []
    for (let col = 0; col < GRID_SIZE; col++) {
      for (let row = 0; row < GRID_SIZE; row++) {
        const [cx, , cz] = gridToWorld(col, row)
        // Square tile corners
        const tl = [cx - HALF, 0, cz - HALF]
        const tr = [cx + HALF, 0, cz - HALF]
        const br = [cx + HALF, 0, cz + HALF]
        const bl = [cx - HALF, 0, cz + HALF]
        // 4 edges per square as line segment pairs
        points.push(...tl, ...tr)
        points.push(...tr, ...br)
        points.push(...br, ...bl)
        points.push(...bl, ...tl)
      }
    }
    const geo = new THREE.BufferGeometry()
    geo.setAttribute('position', new THREE.Float32BufferAttribute(points, 3))
    return geo
  }, [])

  // 2. Compute highlight tiles under cursor
  const highlightTiles = useMemo(() => {
    if (!cursorGridPos || !buildingTypeName) return []
    const tiles = getBuildingTiles(buildingTypeName, cursorGridPos.col, cursorGridPos.row)
    const valid = canPlaceBuilding(buildingTypeName, cursorGridPos.col, cursorGridPos.row, occupiedTiles)
    return tiles.map(t => ({
      ...t,
      worldPos: gridToWorld(t.col, t.row),
      valid,
      inBounds: t.col >= 0 && t.col < GRID_SIZE && t.row >= 0 && t.row < GRID_SIZE,
    }))
  }, [cursorGridPos, buildingTypeName, occupiedTiles])

  // 3. Ground plane sizing (covers entire grid area)
  const groundPlaneSize = GRID_SIZE * TILE_WORLD_SIZE + TILE_WORLD_SIZE

  // Handle pointer move on ground plane -> convert to grid coords
  const handlePointerMove = useCallback((e: ThreeEvent<PointerEvent>) => {
    e.stopPropagation()
    const point = e.point
    const gridPos = worldToGrid(point.x, point.z)
    if (gridPos.col >= 0 && gridPos.col < GRID_SIZE && gridPos.row >= 0 && gridPos.row < GRID_SIZE) {
      onCursorMove(gridPos)
    } else {
      onCursorMove(null)
    }
  }, [onCursorMove])

  return (
    <group position={[0, 0.05, 0]}>
      {/* Wireframe grid lines — always visible */}
      <lineSegments geometry={gridLinesGeometry}>
        <lineBasicMaterial color="#224466" transparent opacity={0.12} depthWrite={false} />
      </lineSegments>

      {/* Highlight tiles under cursor — only during placement */}
      {visible && highlightTiles.filter(t => t.inBounds).map(tile => (
        <SquareTile
          key={`hl-${tile.col},${tile.row}`}
          position={tile.worldPos}
          occupied={!tile.valid}
          onClick={() => {
            if (tile.valid && cursorGridPos) onTileClick(cursorGridPos)
          }}
        />
      ))}

      {/* Invisible ground plane for raycasting — only during placement */}
      {visible && (
        <mesh
          rotation={[-Math.PI / 2, 0, 0]}
          position={[0, -0.01, 0]}
          onPointerMove={handlePointerMove}
          onPointerLeave={() => onCursorMove(null)}
          onClick={(e) => {
            e.stopPropagation()
            if (cursorGridPos && buildingTypeName) {
              const valid = canPlaceBuilding(buildingTypeName, cursorGridPos.col, cursorGridPos.row, occupiedTiles)
              if (valid) onTileClick(cursorGridPos)
            }
          }}
        >
          <planeGeometry args={[groundPlaneSize, groundPlaneSize]} />
          <meshBasicMaterial visible={false} />
        </mesh>
      )}
    </group>
  )
}

function SquareTile({
  position,
  occupied,
  onClick,
}: {
  position: [number, number, number]
  occupied: boolean
  onClick: () => void
}) {
  const color = occupied ? '#ff4444' : '#22cc66'

  return (
    <mesh
      position={position}
      rotation={[-Math.PI / 2, 0, 0]}
      onClick={(e) => {
        e.stopPropagation()
        onClick()
      }}
    >
      <planeGeometry args={[TILE_WORLD_SIZE, TILE_WORLD_SIZE]} />
      <meshStandardMaterial
        color={color}
        emissive={color}
        emissiveIntensity={0.2}
        transparent
        opacity={0.25}
        side={THREE.DoubleSide}
        depthWrite={false}
      />
    </mesh>
  )
}
