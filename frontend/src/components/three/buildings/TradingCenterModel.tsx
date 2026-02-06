import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Mesh, MeshStandardMaterial } from 'three'

interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
}

/**
 * Trading Center: octagonal marketplace/exchange (core, 2x2).
 * Two-level octagonal bazaar with partial retractable dome,
 * cargo containers at base, landing pad ring, trade beacon.
 */
export default function TradingCenterModel({ position = [0, 0, 0], scale = 1, level = 1 }: BuildingComponentProps) {
  const beaconRef = useRef<Mesh>(null)
  const landingRingRef = useRef<Mesh>(null)

  useFrame((state) => {
    const t = state.clock.elapsedTime
    // Trade beacon blinks (0.5-2.0, period ~1s)
    if (beaconRef.current) {
      const mat = beaconRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 1.25 + Math.sin(t * Math.PI * 2) * 0.75
    }
    // Landing ring pulses (0.3-0.7, period ~3s)
    if (landingRingRef.current) {
      const mat = landingRingRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 0.5 + Math.sin(t * (Math.PI * 2 / 3)) * 0.2
    }
  })

  return (
    <group position={position} scale={scale}>
      {/* Landing pad ring on ground */}
      <mesh ref={landingRingRef} position={[0, 0.05, 0]} rotation={[Math.PI / 2, 0, 0]}>
        <torusGeometry args={[3.2, 0.04, 6, 24]} />
        <meshStandardMaterial
          color="#4488ff"
          emissive="#4488ff"
          emissiveIntensity={0.5}
          toneMapped={false}
        />
      </mesh>

      {/* Octagonal base (lower level) */}
      <mesh position={[0, 1.0, 0]} castShadow>
        <cylinderGeometry args={[3.0, 3.0, 2.0, 8]} />
        <meshStandardMaterial
          color="#1a2244"
          emissive="#4488ff"
          emissiveIntensity={0.12}
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* Upper level */}
      <mesh position={[0, 2.75, 0]} castShadow>
        <cylinderGeometry args={[2.5, 2.5, 1.5, 8]} />
        <meshStandardMaterial
          color="#1a2244"
          emissive="#4488ff"
          emissiveIntensity={0.12}
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* Partial dome (retractable, shown partially open) */}
      <mesh position={[0.4, 3.5, 0]} rotation={[0, 0.3, 0]}>
        <sphereGeometry args={[2.5, 10, 8, 0, 4.5, 0, Math.PI / 2]} />
        <meshStandardMaterial
          color="#3a4a77"
          emissive="#4488ff"
          emissiveIntensity={0.3}
          transparent
          opacity={0.4}
          metalness={0.7}
          roughness={0.2}
        />
      </mesh>

      {/* 4 cargo containers scattered at base */}
      {[
        { pos: [2.5, 0.3, 1.5] as [number, number, number], color: '#335522' },
        { pos: [-2.2, 0.3, 1.8] as [number, number, number], color: '#553322' },
        { pos: [2.0, 0.3, -2.0] as [number, number, number], color: '#223355' },
        { pos: [-1.8, 0.3, -2.2] as [number, number, number], color: '#335522' },
      ].map((cargo, i) => (
        <mesh key={`cargo-${i}`} position={cargo.pos} castShadow>
          <boxGeometry args={[1.0, 0.6, 0.6]} />
          <meshStandardMaterial
            color={cargo.color}
            metalness={0.6}
            roughness={0.4}
          />
        </mesh>
      ))}

      {/* Trade beacon on top of dome */}
      <mesh position={[0, 4.5, 0]}>
        <cylinderGeometry args={[0.1, 0.1, 1.5, 6]} />
        <meshStandardMaterial
          color="#ffcc44"
          metalness={0.8}
          roughness={0.2}
        />
      </mesh>
      <mesh ref={beaconRef} position={[0, 5.3, 0]}>
        <sphereGeometry args={[0.15, 6, 6]} />
        <meshStandardMaterial
          color="#ffcc44"
          emissive="#ffcc44"
          emissiveIntensity={1.5}
          toneMapped={false}
        />
      </mesh>
    </group>
  )
}
