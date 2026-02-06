import { useRef, useMemo } from 'react'
import { useFrame } from '@react-three/fiber'
import * as THREE from 'three'
import type { Group } from 'three'

interface BattleshipModelProps {
  position?: [number, number, number]
  color?: string
  scale?: number
  animated?: boolean
}

/**
 * Battleship: large, heavily armored capital ship.
 * Massive hull with layered armor plating, multiple turrets, quad engines.
 */
export default function BattleshipModel({
  position = [0, 0, 0],
  color = '#aa4444',
  scale = 1,
  animated = true,
}: BattleshipModelProps) {
  const groupRef = useRef<Group>(null)

  // Main hull - broad, heavy, imposing
  const hullGeometry = useMemo(() => {
    const shape = new THREE.Shape()
    // Top-down profile: blunt bow, very wide mid, massive stern
    shape.moveTo(0, 3.2)        // bow center
    shape.lineTo(-0.6, 2.6)    // left bow bevel
    shape.lineTo(-1.0, 1.6)    // left forward
    shape.lineTo(-1.4, 0.4)    // left mid forward (widening)
    shape.lineTo(-1.5, -0.6)   // left mid (widest)
    shape.lineTo(-1.4, -1.6)   // left rear
    shape.lineTo(-1.2, -2.4)   // left engine section
    shape.lineTo(-0.8, -2.8)   // left outer engine
    shape.lineTo(-0.6, -3.0)   // left outer engine nozzle
    shape.lineTo(-0.3, -2.6)   // left inner engine gap
    shape.lineTo(0.3, -2.6)    // right inner engine gap
    shape.lineTo(0.6, -3.0)    // right outer engine nozzle
    shape.lineTo(0.8, -2.8)    // right outer engine
    shape.lineTo(1.2, -2.4)    // right engine section
    shape.lineTo(1.4, -1.6)    // right rear
    shape.lineTo(1.5, -0.6)    // right mid (widest)
    shape.lineTo(1.4, 0.4)     // right mid forward
    shape.lineTo(1.0, 1.6)     // right forward
    shape.lineTo(0.6, 2.6)     // right bow bevel
    shape.closePath()

    const extrudeSettings = {
      depth: 0.7,
      bevelEnabled: true,
      bevelThickness: 0.12,
      bevelSize: 0.1,
      bevelSegments: 3,
    }
    const geo = new THREE.ExtrudeGeometry(shape, extrudeSettings)
    geo.center()
    geo.rotateX(Math.PI / 2)
    return geo
  }, [])

  useFrame((state) => {
    if (groupRef.current && animated) {
      const t = state.clock.elapsedTime
      // Very slow, ponderous movement - battleships are massive
      groupRef.current.rotation.z = Math.sin(t * 0.25) * 0.01
      groupRef.current.position.y = position[1] + Math.sin(t * 0.5) * 0.05
    }
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Main hull */}
      <mesh geometry={hullGeometry} castShadow>
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={0.1}
          metalness={0.7}
          roughness={0.3}
        />
      </mesh>

      {/* Upper armor plating layer */}
      <mesh position={[0, 0.4, -0.2]} castShadow>
        <boxGeometry args={[2.0, 0.15, 3.5]} />
        <meshStandardMaterial
          color="#3a2a2a"
          emissive={color}
          emissiveIntensity={0.05}
          metalness={0.85}
          roughness={0.2}
        />
      </mesh>

      {/* Command bridge tower */}
      <group position={[0, 0.5, 1.2]}>
        <mesh castShadow>
          <boxGeometry args={[0.6, 0.4, 0.8]} />
          <meshStandardMaterial
            color="#2a2a3e"
            emissive={color}
            emissiveIntensity={0.15}
            metalness={0.8}
            roughness={0.2}
          />
        </mesh>
        {/* Bridge viewport strip */}
        <mesh position={[0, 0.15, 0.35]}>
          <boxGeometry args={[0.5, 0.08, 0.08]} />
          <meshStandardMaterial
            color="#88ddff"
            emissive="#88ddff"
            emissiveIntensity={0.6}
            transparent
            opacity={0.8}
          />
        </mesh>
        {/* Antenna/sensor mast */}
        <mesh position={[0, 0.35, 0]} rotation={[0, 0, 0]}>
          <cylinderGeometry args={[0.02, 0.02, 0.3, 4]} />
          <meshStandardMaterial color="#666677" metalness={0.9} roughness={0.1} />
        </mesh>
      </group>

      {/* Forward main turret */}
      <group position={[0, 0.35, 2.0]}>
        <mesh castShadow>
          <cylinderGeometry args={[0.22, 0.28, 0.18, 8]} />
          <meshStandardMaterial color="#333344" metalness={0.9} roughness={0.15} />
        </mesh>
        {/* Twin barrels */}
        <mesh position={[-0.08, 0.1, 0.4]} rotation={[Math.PI / 2, 0, 0]}>
          <cylinderGeometry args={[0.04, 0.04, 0.5, 6]} />
          <meshStandardMaterial color="#555566" metalness={0.9} roughness={0.1} />
        </mesh>
        <mesh position={[0.08, 0.1, 0.4]} rotation={[Math.PI / 2, 0, 0]}>
          <cylinderGeometry args={[0.04, 0.04, 0.5, 6]} />
          <meshStandardMaterial color="#555566" metalness={0.9} roughness={0.1} />
        </mesh>
      </group>

      {/* Left broadside turret */}
      <group position={[-1.1, 0.3, -0.4]}>
        <mesh castShadow>
          <cylinderGeometry args={[0.18, 0.22, 0.15, 8]} />
          <meshStandardMaterial color="#333344" metalness={0.9} roughness={0.15} />
        </mesh>
        <mesh position={[0, 0.08, 0.35]} rotation={[Math.PI / 2, 0, 0]}>
          <cylinderGeometry args={[0.035, 0.035, 0.45, 6]} />
          <meshStandardMaterial color="#555566" metalness={0.9} roughness={0.1} />
        </mesh>
      </group>

      {/* Right broadside turret */}
      <group position={[1.1, 0.3, -0.4]}>
        <mesh castShadow>
          <cylinderGeometry args={[0.18, 0.22, 0.15, 8]} />
          <meshStandardMaterial color="#333344" metalness={0.9} roughness={0.15} />
        </mesh>
        <mesh position={[0, 0.08, 0.35]} rotation={[Math.PI / 2, 0, 0]}>
          <cylinderGeometry args={[0.035, 0.035, 0.45, 6]} />
          <meshStandardMaterial color="#555566" metalness={0.9} roughness={0.1} />
        </mesh>
      </group>

      {/* Rear dorsal turret */}
      <group position={[0, 0.35, -1.4]}>
        <mesh castShadow>
          <cylinderGeometry args={[0.2, 0.25, 0.16, 8]} />
          <meshStandardMaterial color="#333344" metalness={0.9} roughness={0.15} />
        </mesh>
        <mesh position={[-0.06, 0.08, -0.35]} rotation={[Math.PI / 2, 0, 0]}>
          <cylinderGeometry args={[0.035, 0.035, 0.4, 6]} />
          <meshStandardMaterial color="#555566" metalness={0.9} roughness={0.1} />
        </mesh>
        <mesh position={[0.06, 0.08, -0.35]} rotation={[Math.PI / 2, 0, 0]}>
          <cylinderGeometry args={[0.035, 0.035, 0.4, 6]} />
          <meshStandardMaterial color="#555566" metalness={0.9} roughness={0.1} />
        </mesh>
      </group>

      {/* Quad engine glows */}
      {[-0.7, -0.3, 0.3, 0.7].map((x, i) => (
        <mesh key={i} position={[x, 0, -2.8]}>
          <cylinderGeometry args={[0.1, 0.15, 0.2, 8]} />
          <meshStandardMaterial
            color="#ff4422"
            emissive="#ff4422"
            emissiveIntensity={1.0}
            toneMapped={false}
          />
        </mesh>
      ))}

      {/* Side armor ridge lines */}
      <mesh position={[-1.45, 0.15, -0.4]} castShadow>
        <boxGeometry args={[0.08, 0.2, 2.8]} />
        <meshStandardMaterial
          color="#4a3333"
          emissive={color}
          emissiveIntensity={0.08}
          metalness={0.85}
          roughness={0.2}
        />
      </mesh>
      <mesh position={[1.45, 0.15, -0.4]} castShadow>
        <boxGeometry args={[0.08, 0.2, 2.8]} />
        <meshStandardMaterial
          color="#4a3333"
          emissive={color}
          emissiveIntensity={0.08}
          metalness={0.85}
          roughness={0.2}
        />
      </mesh>

      {/* Navigation lights */}
      <mesh position={[-1.5, 0, -0.6]}>
        <sphereGeometry args={[0.05, 6, 6]} />
        <meshStandardMaterial
          color="#ff2222"
          emissive="#ff2222"
          emissiveIntensity={2}
          toneMapped={false}
        />
      </mesh>
      <mesh position={[1.5, 0, -0.6]}>
        <sphereGeometry args={[0.05, 6, 6]} />
        <meshStandardMaterial
          color="#22ff22"
          emissive="#22ff22"
          emissiveIntensity={2}
          toneMapped={false}
        />
      </mesh>
      <mesh position={[0, 0.6, 1.5]}>
        <sphereGeometry args={[0.04, 6, 6]} />
        <meshStandardMaterial
          color="#ffffff"
          emissive="#ffffff"
          emissiveIntensity={1.5}
          toneMapped={false}
        />
      </mesh>
    </group>
  )
}
