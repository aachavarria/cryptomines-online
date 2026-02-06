import { useRef, useMemo } from 'react'
import { useFrame } from '@react-three/fiber'
import * as THREE from 'three'
import type { Group } from 'three'

interface CruiserModelProps {
  position?: [number, number, number]
  color?: string
  scale?: number
  animated?: boolean
}

/**
 * Cruiser: medium-sized, balanced warship.
 * Wider hull than frigate, side turret mounts, twin engine blocks.
 */
export default function CruiserModel({
  position = [0, 0, 0],
  color = '#cc8844',
  scale = 1,
  animated = true,
}: CruiserModelProps) {
  const groupRef = useRef<Group>(null)

  // Main hull - broader, more angular than frigate
  const hullGeometry = useMemo(() => {
    const shape = new THREE.Shape()
    // Top-down profile: chisel-nosed, broad mid-section, tapered stern
    shape.moveTo(0, 2.4)        // nose point
    shape.lineTo(-0.5, 1.6)    // left nose bevel
    shape.lineTo(-0.8, 0.6)    // left forward hull
    shape.lineTo(-1.0, -0.2)   // left mid hull (widest)
    shape.lineTo(-0.9, -1.0)   // left rear hull
    shape.lineTo(-0.7, -1.6)   // left engine block outer
    shape.lineTo(-0.5, -2.0)   // left engine nozzle
    shape.lineTo(-0.2, -1.8)   // left engine inner
    shape.lineTo(0.2, -1.8)    // right engine inner
    shape.lineTo(0.5, -2.0)    // right engine nozzle
    shape.lineTo(0.7, -1.6)    // right engine block outer
    shape.lineTo(0.9, -1.0)    // right rear hull
    shape.lineTo(1.0, -0.2)    // right mid hull
    shape.lineTo(0.8, 0.6)     // right forward hull
    shape.lineTo(0.5, 1.6)     // right nose bevel
    shape.closePath()

    const extrudeSettings = {
      depth: 0.5,
      bevelEnabled: true,
      bevelThickness: 0.1,
      bevelSize: 0.08,
      bevelSegments: 2,
    }
    const geo = new THREE.ExtrudeGeometry(shape, extrudeSettings)
    geo.center()
    geo.rotateX(Math.PI / 2)
    return geo
  }, [])

  useFrame((state) => {
    if (groupRef.current && animated) {
      const t = state.clock.elapsedTime
      // Slow, steady drift - cruisers are stable
      groupRef.current.rotation.z = Math.sin(t * 0.4) * 0.02
      groupRef.current.position.y = position[1] + Math.sin(t * 0.7) * 0.04
    }
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Main hull */}
      <mesh geometry={hullGeometry} castShadow>
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={0.12}
          metalness={0.75}
          roughness={0.25}
        />
      </mesh>

      {/* Bridge superstructure */}
      <mesh position={[0, 0.35, 0.8]} castShadow>
        <boxGeometry args={[0.5, 0.2, 0.6]} />
        <meshStandardMaterial
          color="#2a2a3e"
          emissive={color}
          emissiveIntensity={0.2}
          metalness={0.8}
          roughness={0.2}
        />
      </mesh>

      {/* Bridge viewport */}
      <mesh position={[0, 0.42, 1.0]}>
        <boxGeometry args={[0.4, 0.06, 0.15]} />
        <meshStandardMaterial
          color="#88ddff"
          emissive="#88ddff"
          emissiveIntensity={0.6}
          metalness={0.9}
          roughness={0.1}
          transparent
          opacity={0.8}
        />
      </mesh>

      {/* Left turret mount */}
      <group position={[-0.8, 0.2, -0.2]}>
        <mesh castShadow>
          <cylinderGeometry args={[0.15, 0.18, 0.15, 8]} />
          <meshStandardMaterial
            color="#333344"
            metalness={0.9}
            roughness={0.2}
          />
        </mesh>
        {/* Barrel */}
        <mesh position={[0, 0.08, 0.3]} rotation={[Math.PI / 2, 0, 0]}>
          <cylinderGeometry args={[0.04, 0.04, 0.4, 6]} />
          <meshStandardMaterial
            color="#555566"
            metalness={0.9}
            roughness={0.1}
          />
        </mesh>
      </group>

      {/* Right turret mount */}
      <group position={[0.8, 0.2, -0.2]}>
        <mesh castShadow>
          <cylinderGeometry args={[0.15, 0.18, 0.15, 8]} />
          <meshStandardMaterial
            color="#333344"
            metalness={0.9}
            roughness={0.2}
          />
        </mesh>
        <mesh position={[0, 0.08, 0.3]} rotation={[Math.PI / 2, 0, 0]}>
          <cylinderGeometry args={[0.04, 0.04, 0.4, 6]} />
          <meshStandardMaterial
            color="#555566"
            metalness={0.9}
            roughness={0.1}
          />
        </mesh>
      </group>

      {/* Left engine glow */}
      <mesh position={[-0.35, 0, -2.0]}>
        <cylinderGeometry args={[0.1, 0.16, 0.2, 8]} />
        <meshStandardMaterial
          color="#ff6622"
          emissive="#ff6622"
          emissiveIntensity={1.2}
          toneMapped={false}
        />
      </mesh>

      {/* Right engine glow */}
      <mesh position={[0.35, 0, -2.0]}>
        <cylinderGeometry args={[0.1, 0.16, 0.2, 8]} />
        <meshStandardMaterial
          color="#ff6622"
          emissive="#ff6622"
          emissiveIntensity={1.2}
          toneMapped={false}
        />
      </mesh>

      {/* Side running lights */}
      <mesh position={[-1.0, 0, -0.2]}>
        <sphereGeometry args={[0.04, 6, 6]} />
        <meshStandardMaterial
          color="#ff2222"
          emissive="#ff2222"
          emissiveIntensity={2}
          toneMapped={false}
        />
      </mesh>
      <mesh position={[1.0, 0, -0.2]}>
        <sphereGeometry args={[0.04, 6, 6]} />
        <meshStandardMaterial
          color="#22ff22"
          emissive="#22ff22"
          emissiveIntensity={2}
          toneMapped={false}
        />
      </mesh>
    </group>
  )
}
