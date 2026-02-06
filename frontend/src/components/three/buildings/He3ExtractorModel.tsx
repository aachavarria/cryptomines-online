import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Group, Mesh, MeshStandardMaterial } from 'three'

interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
}

/**
 * He3 Extractor: gas extraction tower (resource, 2x2).
 * Tall vertical column with bulbous collection tank at top,
 * radiating pipes to ground-level processing ring. Sleek and vertical.
 * Gas venting effect from the top, orbiting light on processing ring.
 */
export default function He3ExtractorModel({ position = [0, 0, 0], scale = 1, level = 1 }: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const innerGlowRef = useRef<Mesh>(null)
  const vent0Ref = useRef<Mesh>(null)
  const vent1Ref = useRef<Mesh>(null)
  const vent2Ref = useRef<Mesh>(null)
  const orbitLightRef = useRef<Group>(null)

  useFrame((state) => {
    const t = state.clock.elapsedTime
    // Tank bulb inner glow pulses (emissiveIntensity 0.6-1.2, period ~2.5s)
    if (innerGlowRef.current) {
      const mat = innerGlowRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 0.9 + Math.sin(t * (Math.PI * 2 / 2.5)) * 0.3
    }
    // Vent nozzles glow intermittently (staggered, period ~1s each)
    const vents = [vent0Ref, vent1Ref, vent2Ref]
    vents.forEach((ref, i) => {
      if (ref.current) {
        const mat = ref.current.material as MeshStandardMaterial
        const phase = (t + i * 0.33) % 1.0
        mat.emissiveIntensity = phase < 0.5 ? 1.2 : 0.3
      }
    })
    // Processing ring orbiting light point (period ~5s)
    if (orbitLightRef.current) {
      orbitLightRef.current.rotation.y = t * (Math.PI * 2 / 5)
    }
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Ground processing ring */}
      <mesh position={[0, 0.5, 0]} rotation={[Math.PI / 2, 0, 0]}>
        <torusGeometry args={[2.5, 0.2, 8, 12]} />
        <meshStandardMaterial
          color="#2a3a2a"
          metalness={0.85}
          roughness={0.25}
        />
      </mesh>

      {/* Central column */}
      <mesh position={[0, 2.45, 0]} castShadow>
        <cylinderGeometry args={[0.5, 0.5, 4.5, 8]} />
        <meshStandardMaterial
          color="#1a2a1a"
          emissive="#22cc66"
          emissiveIntensity={0.12}
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* Collection tank bulb (translucent) */}
      <mesh position={[0, 5.2, 0]}>
        <sphereGeometry args={[1.2, 10, 8]} />
        <meshStandardMaterial
          color="#2a4a2a"
          emissive="#22cc66"
          emissiveIntensity={0.3}
          transparent
          opacity={0.6}
          metalness={0.6}
          roughness={0.3}
        />
      </mesh>

      {/* Tank inner glow (smaller nested sphere) */}
      <mesh ref={innerGlowRef} position={[0, 5.2, 0]}>
        <sphereGeometry args={[0.7, 8, 6]} />
        <meshStandardMaterial
          color="#44ff88"
          emissive="#44ff88"
          emissiveIntensity={1.0}
          transparent
          opacity={0.4}
          toneMapped={false}
        />
      </mesh>

      {/* Tank cap (cone on top) */}
      <mesh position={[0, 6.5, 0]}>
        <coneGeometry args={[0.3, 0.6, 6]} />
        <meshStandardMaterial
          color="#2a4a2a"
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* 4 connecting pipes from ring up to tank at 45deg */}
      {[0, Math.PI / 2, Math.PI, Math.PI * 1.5].map((angle, i) => {
        const x = Math.cos(angle) * 1.25
        const z = Math.sin(angle) * 1.25
        return (
          <mesh
            key={`pipe-${i}`}
            position={[x, 2.75, z]}
            rotation={[
              angle === 0 || angle === Math.PI ? 0 : (angle > Math.PI ? 0.78 : -0.78),
              0,
              angle === 0 ? -0.78 : angle === Math.PI ? 0.78 : 0,
            ]}
          >
            <cylinderGeometry args={[0.1, 0.1, 3.8, 6]} />
            <meshStandardMaterial
              color="#3a4a3a"
              metalness={0.8}
              roughness={0.25}
            />
          </mesh>
        )
      })}

      {/* 3 vent exhaust nozzles on top */}
      {[-0.25, 0, 0.25].map((xOff, i) => {
        const refs = [vent0Ref, vent1Ref, vent2Ref]
        return (
          <mesh key={`vent-${i}`} ref={refs[i]} position={[xOff, 7.0, 0]}>
            <coneGeometry args={[0.15, 0.4, 6]} />
            <meshStandardMaterial
              color="#44ff88"
              emissive="#44ff88"
              emissiveIntensity={1.2}
              toneMapped={false}
            />
          </mesh>
        )
      })}

      {/* Pressure gauges on column sides */}
      <mesh position={[0.55, 3.0, 0]} rotation={[0, 0, Math.PI / 2]}>
        <cylinderGeometry args={[0.2, 0.2, 0.05, 8]} />
        <meshStandardMaterial
          color="#3a4a3a"
          metalness={0.8}
          roughness={0.3}
        />
      </mesh>
      <mesh position={[-0.55, 3.5, 0]} rotation={[0, 0, Math.PI / 2]}>
        <cylinderGeometry args={[0.2, 0.2, 0.05, 8]} />
        <meshStandardMaterial
          color="#3a4a3a"
          metalness={0.8}
          roughness={0.3}
        />
      </mesh>

      {/* Orbiting light on processing ring */}
      <group ref={orbitLightRef} position={[0, 0.5, 0]}>
        <mesh position={[2.5, 0.3, 0]}>
          <sphereGeometry args={[0.08, 6, 6]} />
          <meshStandardMaterial
            color="#44ff88"
            emissive="#44ff88"
            emissiveIntensity={1.5}
            toneMapped={false}
          />
        </mesh>
      </group>
    </group>
  )
}
