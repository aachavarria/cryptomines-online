import { useRef, useState, useMemo } from 'react'
import { useFrame } from '@react-three/fiber'
import { Html } from '@react-three/drei'
import type { Group } from 'three'
import { FrigateModel, CruiserModel, BattleshipModel } from './ships/index.ts'
import type { Fleet } from '../../types'

interface FleetMarkerProps {
  fleet: Fleet
  position: [number, number, number]
  onClick: () => void
}

// Pick a representative ship model based on the fleet's largest hull class.
function pickModel(fleet: Fleet): 'frigate' | 'cruiser' | 'battleship' {
  const classes = (fleet.stacks ?? []).map(s => s.hull_class).filter(Boolean) as string[]
  if (classes.includes('battleship')) return 'battleship'
  if (classes.includes('cruiser')) return 'cruiser'
  return 'frigate'
}

export default function FleetMarker({ fleet, position, onClick }: FleetMarkerProps) {
  const groupRef = useRef<Group>(null)
  const [hovered, setHovered] = useState(false)
  const modelKind = useMemo(() => pickModel(fleet), [fleet])
  const totalShips = useMemo(
    () => (fleet.stacks ?? []).reduce((sum, s) => sum + (s.ship_count ?? 0), 0),
    [fleet],
  )

  useFrame((state) => {
    if (!groupRef.current) return
    const t = state.clock.elapsedTime
    // Idle hover: gentle bob and slow yaw rotation
    groupRef.current.position.y = position[1] + Math.sin(t * 0.8 + position[0] * 0.3) * 0.15
    groupRef.current.rotation.y = t * 0.2
  })

  const color = fleet.status === 'traveling' ? '#ffaa44'
    : fleet.status === 'combat' ? '#ff4444'
    : '#44ccff'

  return (
    <group
      ref={groupRef}
      position={position}
      onClick={(e) => {
        e.stopPropagation()
        onClick()
      }}
      onPointerOver={(e) => {
        e.stopPropagation()
        setHovered(true)
        document.body.style.cursor = 'pointer'
      }}
      onPointerOut={() => {
        setHovered(false)
        document.body.style.cursor = 'default'
      }}
    >
      {modelKind === 'battleship' && <BattleshipModel scale={0.5} color={color} animated={hovered} />}
      {modelKind === 'cruiser' && <CruiserModel scale={0.6} color={color} animated={hovered} />}
      {modelKind === 'frigate' && <FrigateModel scale={0.7} color={color} animated={hovered} />}

      {/* Always-visible fleet label */}
      <Html
        position={[0, 1.6, 0]}
        center
        occlude={false}
        style={{ pointerEvents: 'none' }}
      >
        <div className={`fleet-marker-label ${hovered ? 'hovered' : ''}`}>
          <div className="fleet-marker-name">{fleet.name}</div>
          <div className="fleet-marker-meta">
            {totalShips > 0 ? `${totalShips} ships` : 'empty'}
            {fleet.status !== 'stationed' && ` · ${fleet.status}`}
          </div>
        </div>
      </Html>
    </group>
  )
}
