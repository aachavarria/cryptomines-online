import { useRef, useMemo } from 'react'
import { useFrame } from '@react-three/fiber'
import * as THREE from 'three'
import type { Group } from 'three'

interface FrigateModelProps {
  position?: [number, number, number]
  color?: string
  scale?: number
  animated?: boolean
}

/**
 * Frigate: small, fast, angular fighter-like silhouette.
 * Narrow fuselage with swept-back wings and prominent engines.
 */
export default function FrigateModel({
  position = [0, 0, 0],
  color = '#44aaff',
  scale = 1,
  animated = true,
}: FrigateModelProps) {
  const groupRef = useRef<Group>(null)

  // Procedural hull geometry - elongated wedge shape
  const hullGeometry = useMemo(() => {
    const shape = new THREE.Shape()
    // Top-down profile: pointed nose, widening to wings, tapering to engines
    shape.moveTo(0, 1.8)       // nose tip
    shape.lineTo(-0.3, 1.0)   // left nose edge
    shape.lineTo(-0.4, 0.2)   // left body
    shape.lineTo(-1.2, -0.4)  // left wing tip
    shape.lineTo(-1.0, -0.8)  // left wing trailing edge
    shape.lineTo(-0.3, -0.6)  // left engine mount
    shape.lineTo(-0.25, -1.4) // left engine nozzle
    shape.lineTo(0.25, -1.4)  // right engine nozzle
    shape.lineTo(0.3, -0.6)   // right engine mount
    shape.lineTo(1.0, -0.8)   // right wing trailing edge
    shape.lineTo(1.2, -0.4)   // right wing tip
    shape.lineTo(0.4, 0.2)    // right body
    shape.lineTo(0.3, 1.0)    // right nose edge
    shape.closePath()

    const extrudeSettings = {
      depth: 0.3,
      bevelEnabled: true,
      bevelThickness: 0.08,
      bevelSize: 0.05,
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
      // Gentle banking idle animation
      groupRef.current.rotation.z = Math.sin(t * 0.6) * 0.05
      groupRef.current.position.y = position[1] + Math.sin(t * 1.2) * 0.03
    }
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Main hull */}
      <mesh geometry={hullGeometry} castShadow>
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={0.15}
          metalness={0.8}
          roughness={0.2}
        />
      </mesh>

      {/* Cockpit canopy */}
      <mesh position={[0, 0.2, 0.5]} castShadow>
        <sphereGeometry args={[0.18, 8, 6, 0, Math.PI * 2, 0, Math.PI / 2]} />
        <meshStandardMaterial
          color="#88ddff"
          emissive="#88ddff"
          emissiveIntensity={0.4}
          metalness={0.9}
          roughness={0.1}
          transparent
          opacity={0.7}
        />
      </mesh>

      {/* Left engine glow */}
      <mesh position={[-0.25, 0, -1.4]}>
        <cylinderGeometry args={[0.08, 0.12, 0.15, 8]} />
        <meshStandardMaterial
          color="#ff6622"
          emissive="#ff6622"
          emissiveIntensity={1.5}
          toneMapped={false}
        />
      </mesh>

      {/* Right engine glow */}
      <mesh position={[0.25, 0, -1.4]}>
        <cylinderGeometry args={[0.08, 0.12, 0.15, 8]} />
        <meshStandardMaterial
          color="#ff6622"
          emissive="#ff6622"
          emissiveIntensity={1.5}
          toneMapped={false}
        />
      </mesh>

      {/* Wing tip lights */}
      <mesh position={[-1.1, 0, -0.4]}>
        <sphereGeometry args={[0.04, 6, 6]} />
        <meshStandardMaterial
          color="#ff2222"
          emissive="#ff2222"
          emissiveIntensity={2}
          toneMapped={false}
        />
      </mesh>
      <mesh position={[1.1, 0, -0.4]}>
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
