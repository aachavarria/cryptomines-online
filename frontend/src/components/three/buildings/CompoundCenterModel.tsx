import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Group, Mesh, MeshStandardMaterial } from 'three'

interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
}

/**
 * Compound Center: administrative complex (core, 2x2).
 * Three interconnected rectangular modules at varying heights,
 * connected by elevated walkways. Central holographic beam display.
 * Window strips glow in alternating patterns.
 */
export default function CompoundCenterModel({ position = [0, 0, 0], scale = 1, level = 1 }: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const holoRef = useRef<Mesh>(null)
  const windowARef = useRef<Mesh>(null)
  const windowBRef = useRef<Mesh>(null)
  const windowCRef = useRef<Mesh>(null)

  useFrame((state) => {
    const t = state.clock.elapsedTime
    // Hologram beam rotates (Y-axis, period ~6s) and pulses (opacity 0.2-0.4)
    if (holoRef.current) {
      holoRef.current.rotation.y = t * (Math.PI * 2 / 6)
      const mat = holoRef.current.material as MeshStandardMaterial
      mat.opacity = 0.3 + Math.sin(t * 1.5) * 0.1
    }
    // Window strips glow in alternating patterns (staggered pulse, period ~4s)
    const wins = [windowARef, windowBRef, windowCRef]
    wins.forEach((ref, i) => {
      if (ref.current) {
        const mat = ref.current.material as MeshStandardMaterial
        mat.emissiveIntensity = 0.5 + Math.sin(t * (Math.PI * 2 / 4) + i * (Math.PI * 2 / 3)) * 0.25
      }
    })
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Module A (tallest) */}
      <mesh position={[-1.5, 2.0, -1.0]} castShadow>
        <boxGeometry args={[2.5, 4.0, 2.5]} />
        <meshStandardMaterial
          color="#1a2244"
          emissive="#4488ff"
          emissiveIntensity={0.1}
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>
      {/* Module A rooftop */}
      <mesh position={[-1.5, 4.1, -1.0]}>
        <boxGeometry args={[2.6, 0.15, 2.6]} />
        <meshStandardMaterial
          color="#2a3366"
          metalness={0.85}
          roughness={0.2}
        />
      </mesh>
      {/* Module A window strip */}
      <mesh ref={windowARef} position={[-0.22, 2.0, -1.0]}>
        <boxGeometry args={[0.05, 2.5, 1.5]} />
        <meshStandardMaterial
          color="#4488ff"
          emissive="#4488ff"
          emissiveIntensity={0.5}
          transparent
          opacity={0.5}
        />
      </mesh>

      {/* Module B (medium) */}
      <mesh position={[1.5, 1.5, -1.0]} castShadow>
        <boxGeometry args={[2.0, 3.0, 2.0]} />
        <meshStandardMaterial
          color="#1a2244"
          emissive="#4488ff"
          emissiveIntensity={0.1}
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>
      {/* Module B rooftop */}
      <mesh position={[1.5, 3.1, -1.0]}>
        <boxGeometry args={[2.1, 0.15, 2.1]} />
        <meshStandardMaterial
          color="#2a3366"
          metalness={0.85}
          roughness={0.2}
        />
      </mesh>
      {/* Module B window strip */}
      <mesh ref={windowBRef} position={[1.5, 1.5, 0.02]}>
        <boxGeometry args={[1.2, 0.05, 0.05]} />
        <meshStandardMaterial
          color="#4488ff"
          emissive="#4488ff"
          emissiveIntensity={0.5}
          transparent
          opacity={0.5}
        />
      </mesh>

      {/* Module C (shortest) */}
      <mesh position={[0, 1.0, 1.5]} castShadow>
        <boxGeometry args={[2.5, 2.0, 2.0]} />
        <meshStandardMaterial
          color="#1a2244"
          emissive="#4488ff"
          emissiveIntensity={0.1}
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>
      {/* Module C rooftop */}
      <mesh position={[0, 2.1, 1.5]}>
        <boxGeometry args={[2.6, 0.15, 2.1]} />
        <meshStandardMaterial
          color="#2a3366"
          metalness={0.85}
          roughness={0.2}
        />
      </mesh>
      {/* Module C window strip */}
      <mesh ref={windowCRef} position={[0, 1.0, 2.52]}>
        <boxGeometry args={[1.5, 0.05, 0.05]} />
        <meshStandardMaterial
          color="#4488ff"
          emissive="#4488ff"
          emissiveIntensity={0.5}
          transparent
          opacity={0.5}
        />
      </mesh>

      {/* Walkway A-B (connecting at Y=2.5) */}
      <mesh position={[0, 2.5, -1.0]}>
        <boxGeometry args={[2.0, 0.15, 0.6]} />
        <meshStandardMaterial
          color="#3a4466"
          metalness={0.9}
          roughness={0.2}
        />
      </mesh>

      {/* Walkway A-C (connecting at Y=1.8) */}
      <mesh position={[-0.75, 1.8, 0.25]}>
        <boxGeometry args={[0.6, 0.15, 2.0]} />
        <meshStandardMaterial
          color="#3a4466"
          metalness={0.9}
          roughness={0.2}
        />
      </mesh>

      {/* Courtyard hologram beam */}
      <mesh ref={holoRef} position={[0, 1.5, 0]}>
        <cylinderGeometry args={[0.5, 0.5, 2.0, 8, 1, true]} />
        <meshStandardMaterial
          color="#66bbff"
          emissive="#66bbff"
          emissiveIntensity={0.8}
          transparent
          opacity={0.3}
          side={2}
          toneMapped={false}
        />
      </mesh>
    </group>
  )
}
