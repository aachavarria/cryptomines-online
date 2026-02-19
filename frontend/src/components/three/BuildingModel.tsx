import { useRef, useMemo, Suspense, useState, useEffect } from 'react'
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
  const hoverLightRef = useRef<THREE.SpotLight>(null)
  const countdown = useCountdown(building.upgrade_finish_at)

  // Check dev mode from localStorage
  const [showDebug, setShowDebug] = useState(() =>
    localStorage.getItem('dev_mode') === 'true'
  )

  useEffect(() => {
    const checkDevMode = () => {
      setShowDebug(localStorage.getItem('dev_mode') === 'true')
    }
    window.addEventListener('storage', checkDevMode)
    // Also check periodically for same-tab changes
    const interval = setInterval(checkDevMode, 1000)
    return () => {
      window.removeEventListener('storage', checkDevMode)
      clearInterval(interval)
    }
  }, [])

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

  const shouldAnimate = isSelected || building.is_upgrading

  // Building is under construction if upgrading to level 1 (initial construction)
  const isUnderConstruction = building.is_upgrading && building.level === 0

  useFrame((state) => {
    // Subtle idle breathing bob (only when selected or upgrading, NOT on hover)
    if (groupRef.current && shouldAnimate) {
      const t = state.clock.elapsedTime
      groupRef.current.position.y = position[1] + Math.sin(t * 0.8 + seed * 6.28) * 0.05
    }

    // Pulsing hover light
    if (hoverLightRef.current && isHovered) {
      const t = state.clock.elapsedTime
      // Oscillate between 2 and 8 intensity
      hoverLightRef.current.intensity = 5 + Math.sin(t * 3) * 3
    }
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
      {/* DEBUG helpers (controlled by DEV button in SideNav) */}
      {showDebug && (
        <>
          {/* Ground level reference (red circle at Y=0) */}
          <mesh position={[0, 0, 0]} rotation={[-Math.PI / 2, 0, 0]}>
            <circleGeometry args={[0.5, 16]} />
            <meshBasicMaterial color="#ff0000" transparent opacity={0.5} />
          </mesh>

          {/* Building footprint (shows expected size based on buildingConfig) */}
          <mesh position={[0, 0.01, 0]} rotation={[-Math.PI / 2, 0, 0]}>
            <planeGeometry args={[size.cols * TILE_WORLD_SIZE, size.rows * TILE_WORLD_SIZE]} />
            <meshBasicMaterial color="#ff0000" transparent opacity={0.3} side={THREE.DoubleSide} />
          </mesh>
        </>
      )}

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
          <TypeModel
            scale={levelScale}
            level={building.level}
            animate={shouldAnimate}
            isUnderConstruction={isUnderConstruction}
          />
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

      {/* Pulsing hover spotlight */}
      {isHovered && (
        <>
          <spotLight
            ref={hoverLightRef}
            position={[8, baseHeight * levelScale + 5, 8]}
            target={groupRef.current || undefined}
            color="#ffffff"
            intensity={50}
            angle={Math.PI / 3}
            penumbra={0.5}
            distance={20}
            decay={2}
          />
        </>
      )}

      {/* Hover tooltip */}
      {isHovered && (
        <Html
          position={[0, baseHeight * levelScale + 1.5, 0]}
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
