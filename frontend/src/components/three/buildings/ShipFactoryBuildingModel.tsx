import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Group, Mesh, MeshStandardMaterial } from 'three'

interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
}

/**
 * Ship Factory building (military, 3x2): Tall hangar with assembly gantry,
 * rotating crane arm, and glowing production bays.
 * Building-context version following the building component pattern.
 */
export default function ShipFactoryBuildingModel({ position = [0, 0, 0], scale = 1, level = 1 }: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const craneRef = useRef<Group>(null)
  const bayLightRef = useRef<Mesh>(null)

  const color = '#ff4444'

  useFrame((state) => {
    const t = state.clock.elapsedTime
    // Slowly rotating crane arm (period ~3.3s)
    if (craneRef.current) {
      craneRef.current.rotation.y = t * (Math.PI * 2 / 3.3)
    }
    // Pulsing production bay light
    if (bayLightRef.current) {
      const mat = bayLightRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 0.6 + Math.sin(t * 2) * 0.3
    }
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Main hangar body */}
      <mesh position={[0, 1.8, 0]} castShadow>
        <boxGeometry args={[2.6, 3.6, 2.2]} />
        <meshStandardMaterial
          color="#2a1a1a"
          emissive={color}
          emissiveIntensity={0.15}
          metalness={0.75}
          roughness={0.3}
        />
      </mesh>

      {/* Roof - angled industrial */}
      <mesh position={[0, 3.8, 0]} castShadow>
        <boxGeometry args={[2.8, 0.3, 2.4]} />
        <meshStandardMaterial
          color="#3a2222"
          emissive={color}
          emissiveIntensity={0.07}
          metalness={0.8}
          roughness={0.2}
        />
      </mesh>

      {/* Assembly gantry tower */}
      <mesh position={[0, 5.0, 0]} castShadow>
        <boxGeometry args={[0.3, 2.2, 0.3]} />
        <meshStandardMaterial color="#444455" metalness={0.9} roughness={0.15} />
      </mesh>

      {/* Rotating crane arm */}
      <group ref={craneRef} position={[0, 5.8, 0]}>
        <mesh castShadow>
          <boxGeometry args={[2.4, 0.12, 0.12]} />
          <meshStandardMaterial color="#555566" metalness={0.9} roughness={0.15} />
        </mesh>
        {/* Crane hook/light */}
        <mesh position={[1.0, -0.15, 0]}>
          <sphereGeometry args={[0.06, 6, 6]} />
          <meshStandardMaterial
            color="#ffaa22"
            emissive="#ffaa22"
            emissiveIntensity={1.5}
            toneMapped={false}
          />
        </mesh>
      </group>

      {/* Production bay opening (front) */}
      <mesh position={[0, 1.0, 1.12]} ref={bayLightRef}>
        <boxGeometry args={[1.6, 1.8, 0.05]} />
        <meshStandardMaterial
          color="#ff8844"
          emissive="#ff8844"
          emissiveIntensity={0.6}
          transparent
          opacity={0.6}
        />
      </mesh>

      {/* Side vents - left */}
      {[0.8, 1.6, 2.4].map((y, i) => (
        <mesh key={`lv-${i}`} position={[-1.32, y, 0]}>
          <boxGeometry args={[0.05, 0.1, 0.8]} />
          <meshStandardMaterial
            color="#ff6633"
            emissive="#ff6633"
            emissiveIntensity={0.5}
            transparent
            opacity={0.5}
          />
        </mesh>
      ))}

      {/* Side vents - right */}
      {[0.8, 1.6, 2.4].map((y, i) => (
        <mesh key={`rv-${i}`} position={[1.32, y, 0]}>
          <boxGeometry args={[0.05, 0.1, 0.8]} />
          <meshStandardMaterial
            color="#ff6633"
            emissive="#ff6633"
            emissiveIntensity={0.5}
            transparent
            opacity={0.5}
          />
        </mesh>
      ))}

      {/* Exhaust stacks */}
      <mesh position={[-0.8, 4.2, -0.6]} castShadow>
        <cylinderGeometry args={[0.12, 0.15, 0.8, 6]} />
        <meshStandardMaterial color="#444455" metalness={0.85} roughness={0.2} />
      </mesh>
      <mesh position={[0.8, 4.2, -0.6]} castShadow>
        <cylinderGeometry args={[0.12, 0.15, 0.8, 6]} />
        <meshStandardMaterial color="#444455" metalness={0.85} roughness={0.2} />
      </mesh>

      {/* Warning lights */}
      <mesh position={[-1.2, 4.0, 1.0]}>
        <sphereGeometry args={[0.06, 6, 6]} />
        <meshStandardMaterial
          color="#ff2222"
          emissive="#ff2222"
          emissiveIntensity={2}
          toneMapped={false}
        />
      </mesh>
      <mesh position={[1.2, 4.0, 1.0]}>
        <sphereGeometry args={[0.06, 6, 6]} />
        <meshStandardMaterial
          color="#ff2222"
          emissive="#ff2222"
          emissiveIntensity={2}
          toneMapped={false}
        />
      </mesh>
    </group>
  )
}
