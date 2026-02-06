import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Group, Mesh, MeshStandardMaterial } from 'three'

interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
}

/**
 * Civic Center: grand domed capitol building (core, 3x3).
 * Central hexagonal tower topped with hemisphere dome, four wings at 90-degree intervals,
 * apex spire with pulsing beacon. The heart of the colony.
 */
export default function CivicCenterModel({ position = [0, 0, 0], scale = 1, level = 1 }: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const domeRef = useRef<Mesh>(null)
  const beaconRef = useRef<Mesh>(null)
  const ringGroupRef = useRef<Group>(null)

  useFrame((state) => {
    const t = state.clock.elapsedTime
    // Dome pulses gently (emissiveIntensity 0.3-0.6, period ~3s)
    if (domeRef.current) {
      const mat = domeRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 0.45 + Math.sin(t * (Math.PI * 2 / 3)) * 0.15
    }
    // Apex beacon pulses brightly (1.0-3.0, period ~1.5s)
    if (beaconRef.current) {
      const mat = beaconRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 2.0 + Math.sin(t * (Math.PI * 2 / 1.5)) * 1.0
    }
    // Ring lights rotate slowly (period ~10s)
    if (ringGroupRef.current) {
      ringGroupRef.current.rotation.y = t * (Math.PI * 2 / 10)
    }
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Hexagonal base platform */}
      <mesh position={[0, 0.3, 0]} castShadow>
        <cylinderGeometry args={[5, 5, 0.6, 6]} />
        <meshStandardMaterial
          color="#1a2244"
          emissive="#4488ff"
          emissiveIntensity={0.15}
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* Central tower */}
      <mesh position={[0, 3.1, 0]} castShadow>
        <cylinderGeometry args={[2.0, 2.0, 5.0, 8]} />
        <meshStandardMaterial
          color="#1a2244"
          emissive="#4488ff"
          emissiveIntensity={0.15}
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* Dome (translucent half-sphere) */}
      <mesh ref={domeRef} position={[0, 5.6, 0]}>
        <sphereGeometry args={[2.2, 10, 8, 0, Math.PI * 2, 0, Math.PI / 2]} />
        <meshStandardMaterial
          color="#6699ff"
          emissive="#6699ff"
          emissiveIntensity={0.4}
          transparent
          opacity={0.5}
          metalness={0.7}
          roughness={0.2}
        />
      </mesh>

      {/* Four wings at 0/90/180/270 degrees */}
      {[0, Math.PI / 2, Math.PI, Math.PI * 1.5].map((angle, i) => (
        <group key={`wing-${i}`} rotation={[0, angle, 0]}>
          {/* Wing body */}
          <mesh position={[3.5, 1.35, 0]} castShadow>
            <boxGeometry args={[4.0, 1.5, 1.5]} />
            <meshStandardMaterial
              color="#2a3366"
              emissive="#4488ff"
              emissiveIntensity={0.1}
              metalness={0.85}
              roughness={0.2}
            />
          </mesh>
          {/* Wing cap */}
          <mesh position={[5.5, 1.35, 0]}>
            <cylinderGeometry args={[0.6, 0.6, 1.5, 6]} />
            <meshStandardMaterial
              color="#2a3366"
              emissive="#66bbff"
              emissiveIntensity={2.0}
              toneMapped={false}
              metalness={0.85}
              roughness={0.2}
            />
          </mesh>
        </group>
      ))}

      {/* Apex spire */}
      <mesh position={[0, 6.8, 0]}>
        <cylinderGeometry args={[0.08, 0.08, 2.0, 6]} />
        <meshStandardMaterial
          color="#4488ff"
          metalness={0.9}
          roughness={0.15}
        />
      </mesh>

      {/* Apex beacon */}
      <mesh ref={beaconRef} position={[0, 7.9, 0]}>
        <sphereGeometry args={[0.15, 6, 6]} />
        <meshStandardMaterial
          color="#66bbff"
          emissive="#66bbff"
          emissiveIntensity={2.0}
          toneMapped={false}
        />
      </mesh>

      {/* Rotating ring lights on base platform */}
      <group ref={ringGroupRef} position={[0, 0.7, 0]}>
        {[0, Math.PI / 2, Math.PI, Math.PI * 1.5].map((angle, i) => (
          <mesh key={`ring-${i}`} position={[Math.cos(angle) * 4.2, 0, Math.sin(angle) * 4.2]}>
            <sphereGeometry args={[0.1, 6, 6]} />
            <meshStandardMaterial
              color="#66bbff"
              emissive="#66bbff"
              emissiveIntensity={1.5}
              toneMapped={false}
            />
          </mesh>
        ))}
      </group>
    </group>
  )
}
