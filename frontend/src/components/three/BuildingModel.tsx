import { useRef, useMemo, Suspense } from 'react'
import { useFrame, useThree } from '@react-three/fiber'
import { Html } from '@react-three/drei'
import * as THREE from 'three'
import type { Group } from 'three'
import SelectionRing from './SelectionRing.tsx'
import UpgradeProgressBar3D from './UpgradeProgressBar3D.tsx'
import SimpleBuildingLOD from './SimpleBuildingLOD.tsx'
import { CATEGORY_COLORS, BUILDING_ABBREVIATIONS, TILE_WORLD_SIZE } from '../../contexts/GameContext.tsx'
import { useCountdown } from '../../hooks/useCountdown.ts'
import { getBuildingSize } from '../../config/buildingConfig.ts'
import { BUILDING_MODELS } from './buildings/index.ts'
import type { BuildingWithType } from '../../types/index.ts'

interface BuildingModelProps {
  building: BuildingWithType
  position: [number, number, number]
  isSelected: boolean
  isHovered: boolean
  isMoving: boolean
  onClick: (screenPos: { x: number; y: number }) => void
  onPointerOver: () => void
  onPointerOut: () => void
}

export default function BuildingModel({
  building,
  position,
  isSelected,
  isHovered,
  isMoving,
  onClick,
  onPointerOver,
  onPointerOut,
}: BuildingModelProps) {
  const groupRef = useRef<Group>(null)
  const countdown = useCountdown(building.upgrade_finish_at)

  const color = CATEGORY_COLORS[building.category] || '#888888'
  const abbr = BUILDING_ABBREVIATIONS[building.type_name] || '??'

  // Scale based on level (1.0 to 1.2)
  const levelScale = 1 + (building.level / building.max_level) * 0.2
  // Height varies by type
  const baseHeight = building.type_name === 'civic_center' ? 5 : 3

  const TypeModel = BUILDING_MODELS[building.type_name]
  const size = getBuildingSize(building.type_name)
  const baseRadius = Math.max(size.cols, size.rows) * TILE_WORLD_SIZE * 0.45
  const { camera } = useThree()
  const zoom = (camera as THREE.OrthographicCamera).zoom
  const useLOD = zoom < 5

  const emissiveIntensity = isSelected ? 0.4 : isHovered ? 0.25 : 0.1

  // Unique seed per building so idle animations aren't all synced
  const seed = useMemo(() => {
    let h = 0
    for (let i = 0; i < building.id.length; i++) {
      h = (h * 31 + building.id.charCodeAt(i)) | 0
    }
    return (h & 0xff) / 255
  }, [building.id])

  const shouldAnimate = isSelected || isHovered || building.is_upgrading

  useFrame((state) => {
    if (!groupRef.current || !shouldAnimate) return
    const t = state.clock.elapsedTime
    // Subtle idle breathing bob (only for active buildings)
    groupRef.current.position.y = position[1] + Math.sin(t * 0.8 + seed * 6.28) * 0.05
  })

  return (
    <group
      ref={groupRef}
      position={position}
      onClick={(e) => {
        e.stopPropagation()
        // Get screen position for context menu
        const rect = (e.nativeEvent.target as HTMLElement)?.getBoundingClientRect?.()
        const screenX = e.nativeEvent instanceof PointerEvent ? e.nativeEvent.clientX : (rect?.left ?? 0)
        const screenY = e.nativeEvent instanceof PointerEvent ? e.nativeEvent.clientY : (rect?.top ?? 0)
        onClick({ x: screenX, y: screenY })
      }}
      onPointerOver={(e) => {
        e.stopPropagation()
        onPointerOver()
        document.body.style.cursor = 'pointer'
      }}
      onPointerOut={() => {
        onPointerOut()
        document.body.style.cursor = 'default'
      }}
    >
      {/* Base platform (circle with enough segments to look smooth from isometric view) */}
      <mesh position={[0, 0.05, 0]} rotation={[-Math.PI / 2, 0, 0]} receiveShadow>
        <circleGeometry args={[baseRadius, 32]} />
        <meshStandardMaterial
          color="#1a1a2e"
          emissive={color}
          emissiveIntensity={0.15}
          metalness={0.5}
          roughness={0.6}
        />
      </mesh>

      {/* Building model — LOD: simplified box at far zoom, full model at close zoom */}
      {useLOD ? (
        <SimpleBuildingLOD typeName={building.type_name} color={color} levelScale={levelScale} />
      ) : TypeModel ? (
        <Suspense fallback={
          <mesh position={[0, baseHeight * levelScale / 2, 0]} castShadow>
            <boxGeometry args={[2.2, baseHeight, 2.2]} />
            <meshStandardMaterial color={color} emissive={color} emissiveIntensity={0.15} metalness={0.7} roughness={0.3} />
          </mesh>
        }>
          <TypeModel scale={levelScale} level={building.level} animate={shouldAnimate} />
        </Suspense>
      ) : (
        <mesh
          position={[0, baseHeight * levelScale / 2, 0]}
          castShadow
          scale={[levelScale, levelScale, levelScale]}
        >
          <boxGeometry args={[2.2, baseHeight, 2.2]} />
          <meshStandardMaterial
            color={color}
            emissive={color}
            emissiveIntensity={emissiveIntensity}
            metalness={0.7}
            roughness={0.3}
            transparent={isMoving}
            opacity={isMoving ? 0.4 : 1}
          />
        </mesh>
      )}

      {/* Abbreviation text on top face */}
      <Html
        position={[0, baseHeight * levelScale + 0.2, 0]}
        center
        occlude={false}
        style={{ pointerEvents: 'none' }}
      >
        <div style={{
          color: '#000',
          fontSize: '11px',
          fontWeight: 700,
          fontFamily: 'Rajdhani, sans-serif',
          textAlign: 'center',
          opacity: 0.7,
        }}>
          {abbr}
        </div>
      </Html>

      {/* Upgrading scaffolding */}
      {building.is_upgrading && (
        <mesh position={[0, baseHeight * levelScale / 2, 0]}>
          <boxGeometry args={[2.8, baseHeight * levelScale + 0.5, 2.8]} />
          <meshStandardMaterial
            color="#ffaa22"
            wireframe
            transparent
            opacity={0.4}
          />
        </mesh>
      )}

      {/* Upgrade progress bar (yellow, above building) */}
      {building.is_upgrading && (
        <UpgradeProgressBar3D
          building={building}
          yOffset={baseHeight * levelScale + 1}
        />
      )}

      {/* Selection ring */}
      <SelectionRing visible={isSelected} radius={baseRadius + 0.5} />

      {/* Hover tooltip - only visible on hover, not permanently */}
      {isHovered && !isSelected && (
        <Html
          position={[0, baseHeight * levelScale + 2, 0]}
          center
          occlude={false}
          style={{ pointerEvents: 'none' }}
        >
          <div className="hover-tooltip-3d">
            <span>Lv: {building.is_upgrading ? building.level + 1 : building.level} {building.display_name}</span>
            {building.is_upgrading && countdown && (
              <span className="hover-tooltip-timer">{countdown}</span>
            )}
          </div>
        </Html>
      )}
    </group>
  )
}
