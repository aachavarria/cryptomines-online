import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Mesh, MeshStandardMaterial } from 'three'

interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
  animate?: boolean
}

/**
 * Galaxy Transporter: teleportation gateway (core, 2x2).
 * Two tall vertical pylons flanking a circular energy portal ring.
 * Portal fill ripples with energy. Energy conduits connect pylons to ring.
 */
export default function GalaxyTransporterModel({ position = [0, 0, 0], scale = 1, level = 1, animate = true }: BuildingComponentProps) {
  const portalRingRef = useRef<Mesh>(null)
  const portalFillRef = useRef<Mesh>(null)
  const leftCapRef = useRef<Mesh>(null)
  const rightCapRef = useRef<Mesh>(null)

  useFrame((state) => {
    if (!animate) return
    const t = state.clock.elapsedTime
    // Portal ring rotates on Z-axis (period ~5s)
    if (portalRingRef.current) {
      portalRingRef.current.rotation.z = t * (Math.PI * 2 / 5)
    }
    // Portal fill ripples (opacity 0.2-0.4, emissiveIntensity 0.4-0.8, period ~2s)
    if (portalFillRef.current) {
      const mat = portalFillRef.current.material as MeshStandardMaterial
      mat.opacity = 0.3 + Math.sin(t * Math.PI) * 0.1
      mat.emissiveIntensity = 0.6 + Math.sin(t * Math.PI) * 0.2
    }
    // Pylon caps pulse alternately (period ~1.5s)
    if (leftCapRef.current) {
      const mat = leftCapRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 0.5 + Math.sin(t * (Math.PI * 2 / 1.5)) * 0.3
    }
    if (rightCapRef.current) {
      const mat = rightCapRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 0.5 + Math.cos(t * (Math.PI * 2 / 1.5)) * 0.3
    }
  })

  return (
    <group position={position} scale={scale}>
      {/* Base platform */}
      <mesh position={[0, 0.2, 0]} castShadow>
        <boxGeometry args={[7.0, 0.4, 4.0]} />
        <meshStandardMaterial
          color="#1a1a2a"
          metalness={0.8}
          roughness={0.3}
        />
      </mesh>

      {/* Left pylon */}
      <mesh position={[-2.5, 3.2, 0]} castShadow>
        <boxGeometry args={[0.8, 6.0, 0.8]} />
        <meshStandardMaterial
          color="#1a2244"
          emissive="#4488ff"
          emissiveIntensity={0.15}
          metalness={0.85}
          roughness={0.2}
        />
      </mesh>
      {/* Left pylon cap */}
      <mesh ref={leftCapRef} position={[-2.5, 6.7, 0]}>
        <coneGeometry args={[0.5, 1.0, 6]} />
        <meshStandardMaterial
          color="#2a3366"
          emissive="#6699ff"
          emissiveIntensity={0.3}
          metalness={0.85}
          roughness={0.2}
        />
      </mesh>

      {/* Right pylon */}
      <mesh position={[2.5, 3.2, 0]} castShadow>
        <boxGeometry args={[0.8, 6.0, 0.8]} />
        <meshStandardMaterial
          color="#1a2244"
          emissive="#4488ff"
          emissiveIntensity={0.15}
          metalness={0.85}
          roughness={0.2}
        />
      </mesh>
      {/* Right pylon cap */}
      <mesh ref={rightCapRef} position={[2.5, 6.7, 0]}>
        <coneGeometry args={[0.5, 1.0, 6]} />
        <meshStandardMaterial
          color="#2a3366"
          emissive="#6699ff"
          emissiveIntensity={0.3}
          metalness={0.85}
          roughness={0.2}
        />
      </mesh>

      {/* Portal ring (vertical torus) */}
      <mesh ref={portalRingRef} position={[0, 3.8, 0]} rotation={[0, Math.PI / 2, 0]}>
        <torusGeometry args={[2.0, 0.2, 8, 24]} />
        <meshStandardMaterial
          color="#4488ff"
          emissive="#4488ff"
          emissiveIntensity={1.0}
          toneMapped={false}
          metalness={0.7}
          roughness={0.2}
        />
      </mesh>

      {/* Portal fill (energy field) */}
      <mesh ref={portalFillRef} position={[0, 3.8, 0]} rotation={[0, Math.PI / 2, 0]}>
        <circleGeometry args={[1.8, 16]} />
        <meshStandardMaterial
          color="#88bbff"
          emissive="#88bbff"
          emissiveIntensity={0.6}
          transparent
          opacity={0.3}
          side={2}
        />
      </mesh>

      {/* Energy conduits - 4 arcs from pylons to ring */}
      {[
        { pos: [-1.3, 5.2, 0] as [number, number, number], rot: [0, 0, 0.6] as [number, number, number] },
        { pos: [-1.3, 2.4, 0] as [number, number, number], rot: [0, 0, -0.6] as [number, number, number] },
        { pos: [1.3, 5.2, 0] as [number, number, number], rot: [0, 0, -0.6] as [number, number, number] },
        { pos: [1.3, 2.4, 0] as [number, number, number], rot: [0, 0, 0.6] as [number, number, number] },
      ].map((conduit, i) => (
        <mesh key={`conduit-${i}`} position={conduit.pos} rotation={conduit.rot}>
          <boxGeometry args={[1.8, 0.06, 0.06]} />
          <meshStandardMaterial
            color="#66aaff"
            emissive="#66aaff"
            emissiveIntensity={0.8}
            toneMapped={false}
          />
        </mesh>
      ))}
    </group>
  )
}
