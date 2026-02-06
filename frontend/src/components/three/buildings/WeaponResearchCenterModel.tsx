import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Group, Mesh, MeshStandardMaterial } from 'three'

interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
}

/**
 * Weapon Research Center (military, 2x2): Bunker-like base with containment sphere
 * held in a cradle of support arms. Hazard markings and experimental aesthetic.
 */
export default function WeaponResearchCenterModel({ position = [0, 0, 0], scale = 1, level = 1 }: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const sphereRef = useRef<Mesh>(null)
  const ringRef = useRef<Mesh>(null)
  const hazardRef = useRef<Mesh>(null)
  const ventRefs = useRef<Mesh[]>([])

  useFrame((state) => {
    const t = state.clock.elapsedTime
    // Containment sphere pulses ominously (period ~2s)
    if (sphereRef.current) {
      const mat = sphereRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 0.3 + Math.sin(t * (Math.PI * 2 / 2)) * 0.25
    }
    // Containment ring rotates (period ~3s)
    if (ringRef.current) {
      ringRef.current.rotation.y = t * (Math.PI * 2 / 3)
    }
    // Hazard light blinks (on/off, period ~0.8s)
    if (hazardRef.current) {
      const mat = hazardRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = Math.sin(t * (Math.PI * 2 / 0.8)) > 0 ? 2.0 : 0.2
    }
    // Exhaust vents glow intermittently (period ~4s)
    ventRefs.current.forEach((mesh, i) => {
      if (mesh) {
        const mat = mesh.material as MeshStandardMaterial
        mat.emissiveIntensity = 0.1 + Math.max(0, Math.sin(t * (Math.PI * 2 / 4) + i * Math.PI)) * 0.4
      }
    })
  })

  const addVentRef = (el: Mesh | null, index: number) => {
    if (el) ventRefs.current[index] = el
  }

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Bunker base - heavy armored */}
      <mesh position={[0, 1.0, 0]} castShadow>
        <boxGeometry args={[6.5, 2.0, 6.5]} />
        <meshStandardMaterial
          color="#2a1a1a"
          emissive="#ff4444"
          emissiveIntensity={0.08}
          metalness={0.85}
          roughness={0.35}
        />
      </mesh>

      {/* Armored ridges on front face */}
      {[0.5, 1.0, 1.5].map((y, i) => (
        <mesh key={`ridge-f-${i}`} position={[0, y, 3.28]}>
          <boxGeometry args={[6.5, 0.15, 0.3]} />
          <meshStandardMaterial color="#3a2222" metalness={0.9} roughness={0.25} />
        </mesh>
      ))}
      {/* Armored ridges on back face */}
      {[0.5, 1.0, 1.5].map((y, i) => (
        <mesh key={`ridge-b-${i}`} position={[0, y, -3.28]}>
          <boxGeometry args={[6.5, 0.15, 0.3]} />
          <meshStandardMaterial color="#3a2222" metalness={0.9} roughness={0.25} />
        </mesh>
      ))}

      {/* Support cradle arms (4, angling inward from corners to meet at containment sphere) */}
      {[
        { x: -2.8, z: -2.8, ry: Math.PI / 4 },
        { x: 2.8, z: -2.8, ry: -Math.PI / 4 },
        { x: -2.8, z: 2.8, ry: (3 * Math.PI) / 4 },
        { x: 2.8, z: 2.8, ry: -(3 * Math.PI) / 4 },
      ].map((arm, i) => (
        <mesh
          key={`arm-${i}`}
          position={[arm.x * 0.55, 3.0, arm.z * 0.55]}
          rotation={[0, arm.ry, -Math.PI / 6]}
          castShadow
        >
          <boxGeometry args={[0.2, 2.5, 0.2]} />
          <meshStandardMaterial color="#444455" metalness={0.85} roughness={0.2} />
        </mesh>
      ))}

      {/* Containment sphere - test chamber */}
      <mesh ref={sphereRef} position={[0, 4.0, 0]}>
        <sphereGeometry args={[1.0, 12, 10]} />
        <meshStandardMaterial
          color="#ff6644"
          emissive="#ff4444"
          emissiveIntensity={0.6}
          transparent
          opacity={0.4}
        />
      </mesh>

      {/* Containment ring around sphere */}
      <mesh ref={ringRef} position={[0, 4.0, 0]} rotation={[Math.PI / 2, 0, 0]}>
        <torusGeometry args={[1.2, 0.08, 6, 16]} />
        <meshStandardMaterial
          color="#ff8844"
          emissive="#ff8844"
          emissiveIntensity={0.8}
          toneMapped={false}
        />
      </mesh>

      {/* Exhaust vents on back face */}
      <mesh ref={(el) => addVentRef(el, 0)} position={[-1.0, 1.2, -3.3]}>
        <boxGeometry args={[0.8, 0.3, 0.1]} />
        <meshStandardMaterial
          color="#ff6633"
          emissive="#ff4444"
          emissiveIntensity={0.2}
          transparent
          opacity={0.5}
        />
      </mesh>
      <mesh ref={(el) => addVentRef(el, 1)} position={[1.0, 1.2, -3.3]}>
        <boxGeometry args={[0.8, 0.3, 0.1]} />
        <meshStandardMaterial
          color="#ff6633"
          emissive="#ff4444"
          emissiveIntensity={0.2}
          transparent
          opacity={0.5}
        />
      </mesh>

      {/* Hazard light on front face */}
      <mesh ref={hazardRef} position={[0, 1.2, 3.3]}>
        <sphereGeometry args={[0.1, 8, 6]} />
        <meshStandardMaterial
          color="#ffaa00"
          emissive="#ffaa00"
          emissiveIntensity={2.0}
          toneMapped={false}
        />
      </mesh>
    </group>
  )
}
