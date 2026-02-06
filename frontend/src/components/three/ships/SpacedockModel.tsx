import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Group, Mesh, MeshStandardMaterial } from 'three'

interface SpacedockModelProps {
  position?: [number, number, number]
  scale?: number
  emissiveIntensity?: number
}

/**
 * Spacedock: orbital repair facility building on the planet surface.
 * Circular docking ring with repair arms, energy conduits, and shield generators.
 * Visually different from Ship Factory - more open/ring-shaped vs enclosed hangar.
 */
export default function SpacedockModel({
  position = [0, 0, 0],
  scale = 1,
  emissiveIntensity = 0.15,
}: SpacedockModelProps) {
  const groupRef = useRef<Group>(null)
  const ringRef = useRef<Mesh>(null)
  const pulseRef = useRef<Mesh>(null)

  const color = '#44ccff' // space category

  useFrame((state) => {
    const t = state.clock.elapsedTime
    if (groupRef.current) {
      groupRef.current.position.y = position[1] + Math.sin(t * 0.7) * 0.03
    }
    // Slowly rotating docking ring
    if (ringRef.current) {
      ringRef.current.rotation.y = t * 0.15
    }
    // Pulsing energy core
    if (pulseRef.current) {
      const mat = pulseRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 0.8 + Math.sin(t * 1.5) * 0.4
    }
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Central pillar / support column */}
      <mesh position={[0, 2.0, 0]} castShadow>
        <cylinderGeometry args={[0.4, 0.6, 4.0, 8]} />
        <meshStandardMaterial
          color="#1a2a3a"
          emissive={color}
          emissiveIntensity={emissiveIntensity}
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* Rotating docking ring (upper) */}
      <group ref={ringRef} position={[0, 3.6, 0]}>
        <mesh rotation={[Math.PI / 2, 0, 0]}>
          <torusGeometry args={[1.6, 0.12, 6, 16]} />
          <meshStandardMaterial
            color="#2a3a4a"
            emissive={color}
            emissiveIntensity={emissiveIntensity * 1.5}
            metalness={0.85}
            roughness={0.2}
          />
        </mesh>

        {/* Docking clamps - 4 evenly spaced */}
        {[0, Math.PI / 2, Math.PI, (3 * Math.PI) / 2].map((angle, i) => (
          <group
            key={`clamp-${i}`}
            position={[
              Math.cos(angle) * 1.6,
              0,
              Math.sin(angle) * 1.6,
            ]}
            rotation={[0, -angle, 0]}
          >
            <mesh castShadow>
              <boxGeometry args={[0.3, 0.15, 0.15]} />
              <meshStandardMaterial
                color="#3a4a5a"
                metalness={0.9}
                roughness={0.15}
              />
            </mesh>
            {/* Clamp light */}
            <mesh position={[0.18, 0, 0]}>
              <sphereGeometry args={[0.04, 6, 6]} />
              <meshStandardMaterial
                color="#44ffaa"
                emissive="#44ffaa"
                emissiveIntensity={1.5}
                toneMapped={false}
              />
            </mesh>
          </group>
        ))}
      </group>

      {/* Lower support ring (static) */}
      <mesh position={[0, 1.0, 0]} rotation={[Math.PI / 2, 0, 0]}>
        <torusGeometry args={[1.2, 0.1, 6, 12]} />
        <meshStandardMaterial
          color="#2a3a4a"
          emissive={color}
          emissiveIntensity={emissiveIntensity}
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* Energy core (center of upper ring) */}
      <mesh ref={pulseRef} position={[0, 3.6, 0]}>
        <sphereGeometry args={[0.25, 12, 12]} />
        <meshStandardMaterial
          color="#44ccff"
          emissive="#44ccff"
          emissiveIntensity={0.8}
          transparent
          opacity={0.6}
          toneMapped={false}
        />
      </mesh>

      {/* Support struts connecting rings to pillar */}
      {[0, (2 * Math.PI) / 3, (4 * Math.PI) / 3].map((angle, i) => (
        <mesh
          key={`strut-${i}`}
          position={[
            Math.cos(angle) * 0.6,
            2.3,
            Math.sin(angle) * 0.6,
          ]}
          rotation={[0, 0, Math.PI / 6]}
          castShadow
        >
          <boxGeometry args={[0.06, 1.8, 0.06]} />
          <meshStandardMaterial
            color="#3a4a5a"
            metalness={0.85}
            roughness={0.2}
          />
        </mesh>
      ))}

      {/* Shield generator dishes (2, opposite sides) */}
      <group position={[-0.5, 4.2, 0]} rotation={[0, 0, -Math.PI / 6]}>
        <mesh castShadow>
          <cylinderGeometry args={[0.0, 0.2, 0.12, 8]} />
          <meshStandardMaterial
            color="#3a4a5a"
            metalness={0.9}
            roughness={0.15}
          />
        </mesh>
      </group>
      <group position={[0.5, 4.2, 0]} rotation={[0, 0, Math.PI / 6]}>
        <mesh castShadow>
          <cylinderGeometry args={[0.0, 0.2, 0.12, 8]} />
          <meshStandardMaterial
            color="#3a4a5a"
            metalness={0.9}
            roughness={0.15}
          />
        </mesh>
      </group>

      {/* Base platform energy conduit lines */}
      {[0, Math.PI / 2, Math.PI, (3 * Math.PI) / 2].map((angle, i) => (
        <mesh
          key={`conduit-${i}`}
          position={[
            Math.cos(angle) * 0.8,
            0.05,
            Math.sin(angle) * 0.8,
          ]}
          rotation={[-Math.PI / 2, 0, angle]}
        >
          <boxGeometry args={[0.04, 0.6, 0.02]} />
          <meshStandardMaterial
            color={color}
            emissive={color}
            emissiveIntensity={0.6}
            transparent
            opacity={0.6}
          />
        </mesh>
      ))}

      {/* Status lights */}
      <mesh position={[0, 4.5, 0]}>
        <sphereGeometry args={[0.05, 6, 6]} />
        <meshStandardMaterial
          color="#44ffaa"
          emissive="#44ffaa"
          emissiveIntensity={2}
          toneMapped={false}
        />
      </mesh>
    </group>
  )
}
