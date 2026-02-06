import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Group, Mesh, MeshStandardMaterial } from 'three'

interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
}

/**
 * Command Center (military, 2x2): Stepped pyramid with flat-topped command deck,
 * holographic tactical display, antenna arrays, and armored viewport strips.
 */
export default function CommandCenterModel({ position = [0, 0, 0], scale = 1, level = 1 }: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const tacticalRef = useRef<Mesh>(null)
  const antennaRefs = useRef<Mesh[]>([])
  const viewportRefs = useRef<Mesh[]>([])

  useFrame((state) => {
    const t = state.clock.elapsedTime
    // Tactical display rotates slowly (period ~8s)
    if (tacticalRef.current) {
      tacticalRef.current.rotation.y = t * (Math.PI * 2 / 8)
      // Tactical display pulses (emissiveIntensity 0.5-1.0, period ~3s)
      const mat = tacticalRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 0.5 + Math.sin(t * (Math.PI * 2 / 3)) * 0.25
    }
    // Antenna tip lights blink in alternating pairs (period ~1s)
    antennaRefs.current.forEach((mesh, i) => {
      if (mesh) {
        const mat = mesh.material as MeshStandardMaterial
        const phase = i % 2 === 0 ? 0 : Math.PI
        mat.emissiveIntensity = Math.sin(t * (Math.PI * 2) + phase) > 0 ? 2.0 : 0.3
      }
    })
    // Viewport strips cycle brightness (staggered wave, period ~4s)
    viewportRefs.current.forEach((mesh, i) => {
      if (mesh) {
        const mat = mesh.material as MeshStandardMaterial
        mat.emissiveIntensity = 0.3 + Math.sin(t * (Math.PI * 2 / 4) + i * 1.2) * 0.2
      }
    })
  })

  const addAntennaRef = (el: Mesh | null, index: number) => {
    if (el) antennaRefs.current[index] = el
  }

  const addViewportRef = (el: Mesh | null, index: number) => {
    if (el) viewportRefs.current[index] = el
  }

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Bottom tier - wide armored base */}
      <mesh position={[0, 0.75, 0]} castShadow>
        <boxGeometry args={[7.0, 1.5, 7.0]} />
        <meshStandardMaterial
          color="#2a1a1a"
          emissive="#ff4444"
          emissiveIntensity={0.1}
          metalness={0.8}
          roughness={0.3}
        />
      </mesh>

      {/* Middle tier - stepped inward */}
      <mesh position={[0, 2.25, 0]} castShadow>
        <boxGeometry args={[5.0, 1.5, 5.0]} />
        <meshStandardMaterial
          color="#2a1a1a"
          emissive="#ff4444"
          emissiveIntensity={0.1}
          metalness={0.8}
          roughness={0.3}
        />
      </mesh>

      {/* Command deck - top tier */}
      <mesh position={[0, 3.6, 0]} castShadow>
        <boxGeometry args={[3.5, 1.2, 3.5]} />
        <meshStandardMaterial
          color="#3a2222"
          emissive="#ff4444"
          emissiveIntensity={0.15}
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* Viewport strips on command deck (4 faces) */}
      {[
        { pos: [0, 3.6, 1.78] as [number, number, number], size: [3.2, 0.4, 0.05] as [number, number, number] },
        { pos: [0, 3.6, -1.78] as [number, number, number], size: [3.2, 0.4, 0.05] as [number, number, number] },
        { pos: [1.78, 3.6, 0] as [number, number, number], size: [0.05, 0.4, 3.2] as [number, number, number] },
        { pos: [-1.78, 3.6, 0] as [number, number, number], size: [0.05, 0.4, 3.2] as [number, number, number] },
      ].map((strip, i) => (
        <mesh key={`vp-${i}`} ref={(el) => addViewportRef(el, i)} position={strip.pos}>
          <boxGeometry args={strip.size} />
          <meshStandardMaterial
            color="#ff6644"
            emissive="#ff6644"
            emissiveIntensity={0.5}
            transparent
            opacity={0.5}
          />
        </mesh>
      ))}

      {/* Display projector column */}
      <mesh position={[0, 4.6, 0]}>
        <cylinderGeometry args={[0.15, 0.15, 0.8, 6]} />
        <meshStandardMaterial
          color="#3a2222"
          emissive="#ff4444"
          emissiveIntensity={0.2}
          metalness={0.85}
          roughness={0.2}
        />
      </mesh>

      {/* Holographic tactical display */}
      <mesh ref={tacticalRef} position={[0, 5.5, 0]} rotation={[Math.PI / 2, 0, 0]}>
        <cylinderGeometry args={[1.0, 1.0, 0.02, 12]} />
        <meshStandardMaterial
          color="#ff6644"
          emissive="#ff6644"
          emissiveIntensity={0.8}
          transparent
          opacity={0.4}
          toneMapped={false}
        />
      </mesh>

      {/* 4 antenna masts at middle tier corners */}
      {[
        [-2.3, 0, -2.3],
        [2.3, 0, -2.3],
        [-2.3, 0, 2.3],
        [2.3, 0, 2.3],
      ].map((offset, i) => (
        <group key={`ant-${i}`}>
          {/* Mast */}
          <mesh position={[offset[0], 3.75, offset[2]]}>
            <cylinderGeometry args={[0.06, 0.06, 1.5, 6]} />
            <meshStandardMaterial color="#555566" metalness={0.9} roughness={0.15} />
          </mesh>
          {/* Antenna tip light */}
          <mesh ref={(el) => addAntennaRef(el, i)} position={[offset[0], 4.55, offset[2]]}>
            <sphereGeometry args={[0.06, 6, 6]} />
            <meshStandardMaterial
              color="#ff2222"
              emissive="#ff2222"
              emissiveIntensity={2.0}
              toneMapped={false}
            />
          </mesh>
        </group>
      ))}
    </group>
  )
}
