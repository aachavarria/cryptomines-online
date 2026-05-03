import { useMemo, useEffect, Suspense } from 'react'
import * as THREE from 'three'
import { getBuildingSize } from '../../config/buildingConfig.ts'
import { TILE_WORLD_SIZE, gridToWorld } from '../../contexts/GameContext.tsx'
import type { GridPosition } from '../../contexts/GameContext.tsx'
import { BUILDING_MODELS } from './buildings/index.ts'
import { colors3d, colors3dOpacity } from '../../config/colors3d.ts'

interface GhostPreviewProps {
  typeName: string       // building type being placed
  gridPosition: GridPosition | null  // current cursor grid position (anchor)
  isValid: boolean       // whether placement is valid (green) or not (red)
  visible: boolean
}

export default function GhostPreview({ typeName, gridPosition, isValid, visible }: GhostPreviewProps) {
  const size = getBuildingSize(typeName)
  const TypeModel = BUILDING_MODELS[typeName]

  // Calculate world position (center of the multi-tile footprint)
  const worldPos = useMemo(() => {
    if (!gridPosition) return null
    const centerCol = gridPosition.col + (size.cols - 1) / 2
    const centerRow = gridPosition.row + (size.rows - 1) / 2
    return gridToWorld(centerCol, centerRow)
  }, [gridPosition, size])

  // Ghost box dimensions in world units (for fallback)
  const boxWidth = size.cols * TILE_WORLD_SIZE * 0.95
  const boxDepth = size.rows * TILE_WORLD_SIZE * 0.95
  const boxHeight = Math.max(size.cols, size.rows) * 2 + 2

  // Per design-system §7.6: valid = teal-success, invalid = danger.
  const color = isValid ? colors3d.tileValid : colors3d.tileInvalid
  const footprintOpacity = isValid ? colors3dOpacity.tileValid : colors3dOpacity.tileInvalid

  if (!visible || !worldPos) return null

  return (
    <group position={worldPos}>
      {/* 3D Model in hologram mode */}
      {TypeModel ? (
        <Suspense fallback={
          <mesh position={[0, boxHeight / 2, 0]}>
            <boxGeometry args={[boxWidth, boxHeight, boxDepth]} />
            <meshStandardMaterial
              color={color}
              transparent
              opacity={0.3}
              wireframe
            />
          </mesh>
        }>
          <HologramBuilding
            TypeModel={TypeModel}
            color={color}
            isValid={isValid}
          />
        </Suspense>
      ) : (
        // Fallback to box if model not found
        <>
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
        </>
      )}

      {/* Base footprint indicator */}
      <mesh position={[0, 0.05, 0]} rotation={[-Math.PI / 2, 0, 0]}>
        <planeGeometry args={[boxWidth, boxDepth]} />
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={0.4}
          transparent
          opacity={footprintOpacity}
          depthWrite={false}
          side={THREE.DoubleSide}
        />
      </mesh>
    </group>
  )
}

// Component to apply hologram effect to building models
function HologramBuilding({ TypeModel, color, isValid: _isValid }: {
  TypeModel: any,
  color: number,
  isValid: boolean
}) {
  useEffect(() => {
    // Apply hologram material to all meshes in the model
    return () => {
      // Cleanup if needed
    }
  }, [])

  return (
    <TypeModel
      position={[0, 0, 0]}
      scale={1}
      level={1}
      animate={false}
      isUnderConstruction={true}
      hologramColor={color}
    />
  )
}
