import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Group, Mesh, MeshStandardMaterial } from 'three'

interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
  animate?: boolean
}

/**
 * Alliance Center: diplomatic meeting hall (core, 2x2).
 * Circular base with hexagonal pillar ring supporting a floating platform.
 * Holographic globe hovers in center. Platform edge has orbiting lights.
 */
export default function AllianceCenterModel({ position = [0, 0, 0], scale = 1, level = 1, animate = true }: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const globeRef = useRef<Mesh>(null)
  const globeWireRef = useRef<Mesh>(null)
  const orbitRef = useRef<Group>(null)

  useFrame((state) => {
    if (!animate) return
    const t = state.clock.elapsedTime
    // Globe rotates (Y-axis, period ~6s)
    if (globeRef.current) {
      globeRef.current.rotation.y = t * (Math.PI * 2 / 6)
    }
    if (globeWireRef.current) {
      globeWireRef.current.rotation.y = t * (Math.PI * 2 / 6)
      // Globe pulses (opacity 0.3-0.5, period ~2s)
      const mat = globeWireRef.current.material as MeshStandardMaterial
      mat.opacity = 0.4 + Math.sin(t * Math.PI) * 0.1
    }
    // Platform edge lights orbit (period ~8s)
    if (orbitRef.current) {
      orbitRef.current.rotation.y = t * (Math.PI * 2 / 8)
    }
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Circular base */}
      <mesh position={[0, 0.4, 0]} castShadow>
        <cylinderGeometry args={[3.2, 3.2, 0.8, 12]} />
        <meshStandardMaterial
          color="#1a2244"
          emissive="#4488ff"
          emissiveIntensity={0.1}
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* 6 pillars in hexagonal arrangement */}
      {Array.from({ length: 6 }).map((_, i) => {
        const angle = (i / 6) * Math.PI * 2
        const x = Math.cos(angle) * 2.5
        const z = Math.sin(angle) * 2.5
        return (
          <group key={`pillar-${i}`}>
            <mesh position={[x, 2.55, z]} castShadow>
              <cylinderGeometry args={[0.15, 0.15, 3.5, 6]} />
              <meshStandardMaterial
                color="#3a4466"
                metalness={0.85}
                roughness={0.2}
              />
            </mesh>
            {/* Connecting arches between adjacent pillars */}
            {i < 6 && (
              <mesh
                position={[
                  (x + Math.cos(((i + 1) / 6) * Math.PI * 2) * 2.5) / 2,
                  4.2,
                  (z + Math.sin(((i + 1) / 6) * Math.PI * 2) * 2.5) / 2,
                ]}
                rotation={[0, -angle - Math.PI / 6, 0]}
              >
                <boxGeometry args={[1.3, 0.08, 0.08]} />
                <meshStandardMaterial
                  color="#3a4466"
                  metalness={0.85}
                  roughness={0.2}
                />
              </mesh>
            )}
          </group>
        )
      })}

      {/* Floating platform */}
      <mesh position={[0, 3.8, 0]} castShadow>
        <cylinderGeometry args={[2.0, 2.0, 0.3, 8]} />
        <meshStandardMaterial
          color="#2a3355"
          emissive="#4488ff"
          emissiveIntensity={0.2}
          metalness={0.8}
          roughness={0.2}
        />
      </mesh>

      {/* Holographic globe (solid) */}
      <mesh ref={globeRef} position={[0, 5.0, 0]}>
        <sphereGeometry args={[0.8, 12, 8]} />
        <meshStandardMaterial
          color="#4488ff"
          emissive="#4488ff"
          emissiveIntensity={0.8}
          transparent
          opacity={0.4}
        />
      </mesh>

      {/* Globe wireframe overlay */}
      <mesh ref={globeWireRef} position={[0, 5.0, 0]}>
        <sphereGeometry args={[0.82, 12, 8]} />
        <meshStandardMaterial
          color="#88bbff"
          emissive="#88bbff"
          emissiveIntensity={1.0}
          wireframe
          transparent
          opacity={0.4}
          toneMapped={false}
        />
      </mesh>

      {/* Orbiting platform edge lights */}
      <group ref={orbitRef} position={[0, 4.0, 0]}>
        {Array.from({ length: 6 }).map((_, i) => {
          const angle = (i / 6) * Math.PI * 2
          return (
            <mesh key={`orbit-${i}`} position={[Math.cos(angle) * 1.9, 0, Math.sin(angle) * 1.9]}>
              <sphereGeometry args={[0.06, 6, 6]} />
              <meshStandardMaterial
                color="#66bbff"
                emissive="#66bbff"
                emissiveIntensity={1.5}
                toneMapped={false}
              />
            </mesh>
          )
        })}
      </group>
    </group>
  )
}
