import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Group, Mesh, MeshStandardMaterial } from 'three'

interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
}

/**
 * Metal Collector: industrial mining rig (resource, 2x2).
 * Central drill shaft with support frame, angled excavation arm
 * with rotating drum, ore containers at base. Rugged, heavy-industry feel.
 */
export default function MetalCollectorModel({ position = [0, 0, 0], scale = 1, level = 1 }: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const drumRef = useRef<Mesh>(null)
  const shaftRef = useRef<Mesh>(null)
  const light0Ref = useRef<Mesh>(null)
  const light1Ref = useRef<Mesh>(null)

  useFrame((state) => {
    const t = state.clock.elapsedTime
    // Excavation drum rotates (local X-axis, period ~2s)
    if (drumRef.current) {
      drumRef.current.rotation.x = t * (Math.PI * 2 / 2)
    }
    // Drill shaft housing rotates slowly (Y-axis, period ~8s)
    if (shaftRef.current) {
      shaftRef.current.rotation.y = t * (Math.PI * 2 / 8)
    }
    // Ore container indicator lights blink alternating (period ~1.5s)
    if (light0Ref.current) {
      const mat = light0Ref.current.material as MeshStandardMaterial
      mat.emissiveIntensity = Math.sin(t * (Math.PI * 2 / 1.5)) > 0 ? 1.5 : 0.2
    }
    if (light1Ref.current) {
      const mat = light1Ref.current.material as MeshStandardMaterial
      mat.emissiveIntensity = Math.cos(t * (Math.PI * 2 / 1.5)) > 0 ? 1.5 : 0.2
    }
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* 4 support frame legs (slight inward tilt) */}
      {[
        [-1.5, 0, -1.5],
        [1.5, 0, -1.5],
        [-1.5, 0, 1.5],
        [1.5, 0, 1.5],
      ].map((pos, i) => (
        <mesh key={`leg-${i}`} position={[pos[0] * 0.9, 1.5, pos[2] * 0.9]} castShadow>
          <boxGeometry args={[0.15, 3.0, 0.15]} />
          <meshStandardMaterial
            color="#2a3a2a"
            metalness={0.8}
            roughness={0.35}
          />
        </mesh>
      ))}

      {/* Cross beams at top connecting legs */}
      {[
        { pos: [0, 3.0, -1.3] as [number, number, number], size: [2.7, 0.1, 0.1] as [number, number, number] },
        { pos: [0, 3.0, 1.3] as [number, number, number], size: [2.7, 0.1, 0.1] as [number, number, number] },
        { pos: [-1.3, 3.0, 0] as [number, number, number], size: [0.1, 0.1, 2.7] as [number, number, number] },
        { pos: [1.3, 3.0, 0] as [number, number, number], size: [0.1, 0.1, 2.7] as [number, number, number] },
      ].map((beam, i) => (
        <mesh key={`beam-${i}`} position={beam.pos}>
          <boxGeometry args={beam.size} />
          <meshStandardMaterial
            color="#2a3a2a"
            metalness={0.8}
            roughness={0.35}
          />
        </mesh>
      ))}

      {/* Drill shaft housing (slowly rotating) */}
      <mesh ref={shaftRef} position={[0, 1.5, 0]} castShadow>
        <cylinderGeometry args={[0.8, 0.8, 3.0, 8]} />
        <meshStandardMaterial
          color="#1a2a1a"
          emissive="#22cc66"
          emissiveIntensity={0.1}
          metalness={0.8}
          roughness={0.3}
        />
      </mesh>

      {/* Drill head (inverted cone, partially embedded in ground) */}
      <mesh position={[0, -0.3, 0]} rotation={[Math.PI, 0, 0]}>
        <coneGeometry args={[0.6, 1.0, 8]} />
        <meshStandardMaterial
          color="#667766"
          metalness={0.9}
          roughness={0.15}
        />
      </mesh>

      {/* Excavation arm (angled 30deg down from top of shaft) */}
      <group position={[0, 2.8, 0]} rotation={[0, 0, -Math.PI * 30 / 180]}>
        <mesh position={[1.75, 0, 0]}>
          <boxGeometry args={[3.5, 0.2, 0.4]} />
          <meshStandardMaterial
            color="#2a3a2a"
            metalness={0.8}
            roughness={0.35}
          />
        </mesh>
        {/* Drum at arm end (rotating) */}
        <mesh ref={drumRef} position={[3.3, 0, 0]}>
          <cylinderGeometry args={[0.5, 0.5, 0.8, 6]} />
          <meshStandardMaterial
            color="#3a4a3a"
            metalness={0.7}
            roughness={0.3}
          />
        </mesh>
      </group>

      {/* Accent light on top of frame */}
      <mesh position={[0, 3.2, 0]}>
        <sphereGeometry args={[0.08, 6, 6]} />
        <meshStandardMaterial
          color="#44ff88"
          emissive="#44ff88"
          emissiveIntensity={1.5}
          toneMapped={false}
        />
      </mesh>

      {/* 2 ore containers at base */}
      <mesh position={[-2.0, 0.4, 1.0]} castShadow>
        <boxGeometry args={[1.0, 0.8, 0.8]} />
        <meshStandardMaterial
          color="#445533"
          emissive="#22cc66"
          emissiveIntensity={0.15}
          metalness={0.7}
          roughness={0.4}
        />
      </mesh>
      <mesh position={[2.0, 0.4, -1.0]} castShadow>
        <boxGeometry args={[1.0, 0.8, 0.8]} />
        <meshStandardMaterial
          color="#445533"
          emissive="#22cc66"
          emissiveIntensity={0.15}
          metalness={0.7}
          roughness={0.4}
        />
      </mesh>

      {/* Ore container indicator lights */}
      <mesh ref={light0Ref} position={[-2.0, 0.9, 1.0]}>
        <sphereGeometry args={[0.06, 6, 6]} />
        <meshStandardMaterial
          color="#44ff88"
          emissive="#44ff88"
          emissiveIntensity={1.5}
          toneMapped={false}
        />
      </mesh>
      <mesh ref={light1Ref} position={[2.0, 0.9, -1.0]}>
        <sphereGeometry args={[0.06, 6, 6]} />
        <meshStandardMaterial
          color="#44ff88"
          emissive="#44ff88"
          emissiveIntensity={1.5}
          toneMapped={false}
        />
      </mesh>
    </group>
  )
}
