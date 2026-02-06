import { useMemo } from 'react'
import * as THREE from 'three'
import { getBuildingSize } from '../../config/buildingConfig.ts'
import { TILE_WORLD_SIZE, gridToWorld } from '../../contexts/GameContext.tsx'
import type { GridPosition } from '../../contexts/GameContext.tsx'

interface GhostPreviewProps {
  typeName: string       // building type being placed
  gridPosition: GridPosition | null  // current cursor grid position (anchor)
  isValid: boolean       // whether placement is valid (green) or not (red)
  visible: boolean
}

export default function GhostPreview({ typeName, gridPosition, isValid, visible }: GhostPreviewProps) {
  const size = getBuildingSize(typeName)

  // Calculate world position (center of the multi-tile footprint)
  const worldPos = useMemo(() => {
    if (!gridPosition) return null
    const centerCol = gridPosition.col + (size.cols - 1) / 2
    const centerRow = gridPosition.row + (size.rows - 1) / 2
    return gridToWorld(centerCol, centerRow)
  }, [gridPosition, size])

  // Ghost box dimensions in world units (square grid — axis-aligned)
  const boxWidth = size.cols * TILE_WORLD_SIZE * 0.95
  const boxDepth = size.rows * TILE_WORLD_SIZE * 0.95
  const boxHeight = Math.max(size.cols, size.rows) * 2 + 2 // taller for bigger buildings

  const color = isValid ? '#22cc66' : '#ff4444'

  if (!visible || !worldPos) return null

  return (
    <group position={worldPos}>
      {/* Ghost building shape */}
      <mesh position={[0, boxHeight / 2, 0]}>
        <boxGeometry args={[boxWidth, boxHeight, boxDepth]} />
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={0.3}
          transparent
          opacity={0.25}
          depthWrite={false}
          side={THREE.DoubleSide}
        />
      </mesh>

      {/* Wireframe outline */}
      <mesh position={[0, boxHeight / 2, 0]}>
        <boxGeometry args={[boxWidth, boxHeight, boxDepth]} />
        <meshStandardMaterial
          color={color}
          wireframe
          transparent
          opacity={0.5}
          depthWrite={false}
        />
      </mesh>

      {/* Base footprint indicator */}
      <mesh position={[0, 0.05, 0]} rotation={[-Math.PI / 2, 0, 0]}>
        <planeGeometry args={[boxWidth, boxDepth]} />
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={0.4}
          transparent
          opacity={0.3}
          depthWrite={false}
          side={THREE.DoubleSide}
        />
      </mesh>
    </group>
  )
}
