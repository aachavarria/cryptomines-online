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
 * Recycling Plant: debris processing facility (2x2).
 * Funnel-shaped intake hopper on top feeding into boxy processing core,
 * conveyor belts extending from sides, output chute, and status indicator.
 */
export default function RecyclingPlantModel({
  position = [0, 0, 0],
  scale = 1,
  level = 1,
  animate = true,
}: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const hopperRimRef = useRef<Mesh>(null)
  const statusRef = useRef<Mesh>(null)
  const rollerRefs = useRef<(Mesh | null)[]>([])
  const hopperGlowRef = useRef<Mesh>(null)

  const color = '#44ccff'

  useFrame((state) => {
    if (!animate) return
    const t = state.clock.elapsedTime
    // Hopper rim pulses - intake active
    if (hopperRimRef.current) {
      const mat = hopperRimRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 0.3 + Math.sin(t * (Math.PI * 2 / 2)) * 0.15
    }
    // Conveyor rollers rotate
    rollerRefs.current.forEach((ref) => {
      if (ref) {
        ref.rotation.z += 0.1
      }
    })
    // Status light alternates green/cyan
    if (statusRef.current) {
      const mat = statusRef.current.material as MeshStandardMaterial
      const cycle = Math.sin(t * (Math.PI * 2 / 3))
      if (cycle > 0) {
        mat.color.setHex(0x44ff88)
        mat.emissive.setHex(0x44ff88)
      } else {
        mat.color.setHex(0x44ccff)
        mat.emissive.setHex(0x44ccff)
      }
    }
    // Faint glow inside hopper pulses
    if (hopperGlowRef.current) {
      const mat = hopperGlowRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 0.4 + Math.sin(t * 3) * 0.3
    }
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Processing core - main body */}
      <mesh position={[0, 1.25, 0]} castShadow>
        <boxGeometry args={[4.5, 2.5, 4.5]} />
        <meshStandardMaterial
          color="#1a2a3a"
          emissive={color}
          emissiveIntensity={0.1}
          metalness={0.75}
          roughness={0.3}
        />
      </mesh>

      {/* Intake hopper - inverted cone/funnel on top */}
      <mesh position={[0, 3.25, 0]} castShadow>
        <cylinderGeometry args={[1.5, 0.5, 1.5, 8]} />
        <meshStandardMaterial
          color="#2a3a4a"
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* Hopper glow inside */}
      <mesh ref={hopperGlowRef} position={[0, 3.0, 0]}>
        <cylinderGeometry args={[0.4, 0.4, 0.5, 8]} />
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={0.5}
          transparent
          opacity={0.4}
        />
      </mesh>

      {/* Hopper rim */}
      <mesh
        ref={hopperRimRef}
        position={[0, 4.0, 0]}
        rotation={[Math.PI / 2, 0, 0]}
      >
        <torusGeometry args={[1.5, 0.1, 6, 12]} />
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={0.4}
          metalness={0.8}
          roughness={0.2}
        />
      </mesh>

      {/* Left conveyor belt */}
      <mesh position={[-3.5, 1.0, 0]} rotation={[0, 0, Math.PI * 5 / 180]}>
        <boxGeometry args={[3.0, 0.1, 0.6]} />
        <meshStandardMaterial
          color="#3a4a4a"
          metalness={0.7}
          roughness={0.3}
        />
      </mesh>

      {/* Left conveyor rollers */}
      {[-2.5, -3.5, -4.5].map((x, i) => (
        <mesh
          key={`lr-${i}`}
          ref={(el) => { rollerRefs.current[i] = el }}
          position={[x, 1.1, 0]}
          rotation={[Math.PI / 2, 0, 0]}
        >
          <cylinderGeometry args={[0.08, 0.08, 0.6, 6]} />
          <meshStandardMaterial
            color="#555566"
            metalness={0.9}
            roughness={0.15}
          />
        </mesh>
      ))}

      {/* Right conveyor belt */}
      <mesh position={[3.5, 1.0, 0]} rotation={[0, 0, -Math.PI * 5 / 180]}>
        <boxGeometry args={[3.0, 0.1, 0.6]} />
        <meshStandardMaterial
          color="#3a4a4a"
          metalness={0.7}
          roughness={0.3}
        />
      </mesh>

      {/* Right conveyor rollers */}
      {[2.5, 3.5, 4.5].map((x, i) => (
        <mesh
          key={`rr-${i}`}
          ref={(el) => { rollerRefs.current[i + 3] = el }}
          position={[x, 1.1, 0]}
          rotation={[Math.PI / 2, 0, 0]}
        >
          <cylinderGeometry args={[0.08, 0.08, 0.6, 6]} />
          <meshStandardMaterial
            color="#555566"
            metalness={0.9}
            roughness={0.15}
          />
        </mesh>
      ))}

      {/* Output chute - back, angled down */}
      <mesh position={[0, 0.8, -2.8]} rotation={[Math.PI * 15 / 180, 0, 0]} castShadow>
        <boxGeometry args={[1.0, 0.5, 0.5]} />
        <meshStandardMaterial
          color="#2a3a2a"
          metalness={0.75}
          roughness={0.3}
        />
      </mesh>

      {/* Processing status indicator light */}
      <mesh ref={statusRef} position={[0, 1.8, 2.28]}>
        <cylinderGeometry args={[0.15, 0.15, 0.3, 8]} />
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
