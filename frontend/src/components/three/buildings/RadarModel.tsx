import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import { DoubleSide } from 'three'
import type { Group, Mesh, MeshStandardMaterial } from 'three'

interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
  animate?: boolean
}

/**
 * Radar (military, 1x1): Slender scanning tower with a large rotating
 * parabolic dish at the top. Classic radar silhouette -- tall and narrow.
 */
export default function RadarModel({ position = [0, 0, 0], scale = 1, level = 1, animate = true }: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const dishRef = useRef<Group>(null)
  const feedTipRef = useRef<Mesh>(null)
  const baseLightRef = useRef<Mesh>(null)

  useFrame((state) => {
    if (!animate) return
    const t = state.clock.elapsedTime
    // Dish rotates (Y-axis, period ~4s) - the signature radar spin
    if (dishRef.current) {
      dishRef.current.rotation.y = t * (Math.PI * 2 / 4)
    }
    // Feed tip glows with each sweep (spikes when facing forward)
    if (feedTipRef.current) {
      const mat = feedTipRef.current.material as MeshStandardMaterial
      const angle = (t * (Math.PI * 2 / 4)) % (Math.PI * 2)
      // Spike when facing Z+ (forward), fade otherwise
      const forwardness = Math.max(0, Math.cos(angle))
      mat.emissiveIntensity = 0.5 + forwardness * 2.5
    }
    // Base indicator light blinks (period ~2s)
    if (baseLightRef.current) {
      const mat = baseLightRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = Math.sin(t * (Math.PI * 2 / 2)) > 0 ? 1.5 : 0.3
    }
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Equipment base box */}
      <mesh position={[0, 0.5, 0]} castShadow>
        <boxGeometry args={[2.0, 1.0, 2.0]} />
        <meshStandardMaterial
          color="#2a1a1a"
          emissive="#ff4444"
          emissiveIntensity={0.1}
          metalness={0.8}
          roughness={0.3}
        />
      </mesh>

      {/* Mast */}
      <mesh position={[0, 2.75, 0]} castShadow>
        <cylinderGeometry args={[0.15, 0.15, 3.5, 6]} />
        <meshStandardMaterial
          color="#444455"
          metalness={0.9}
          roughness={0.15}
        />
      </mesh>

      {/* Support guy-wires (3 thin struts from base edges to mid-mast) */}
      {[0, (2 * Math.PI) / 3, (4 * Math.PI) / 3].map((angle, i) => {
        const x = Math.cos(angle) * 0.9
        const z = Math.sin(angle) * 0.9
        return (
          <mesh
            key={`wire-${i}`}
            position={[x * 0.5, 1.5, z * 0.5]}
            rotation={[
              Math.atan2(z, 2.0) * 0.4,
              -angle,
              Math.PI / 7,
            ]}
          >
            <boxGeometry args={[0.02, 2.5, 0.02]} />
            <meshStandardMaterial color="#555566" metalness={0.85} roughness={0.2} />
          </mesh>
        )
      })}

      {/* Rotating dish assembly */}
      <group ref={dishRef} position={[0, 4.5, 0]}>
        {/* Dish mount pivot */}
        <mesh>
          <boxGeometry args={[0.3, 0.3, 0.3]} />
          <meshStandardMaterial color="#3a2222" metalness={0.85} roughness={0.2} />
        </mesh>

        {/* Parabolic dish (inverted shallow cone) */}
        <mesh position={[0, 0.1, 0.4]} rotation={[-Math.PI / 2, 0, 0]}>
          <coneGeometry args={[1.5, 0.5, 12, 1, true]} />
          <meshStandardMaterial
            color="#3a2222"
            emissive="#ff4444"
            emissiveIntensity={0.15}
            metalness={0.8}
            roughness={0.25}
            side={DoubleSide}
          />
        </mesh>

        {/* Feed horn extending from dish center forward */}
        <mesh position={[0, 0.1, 0.8]} rotation={[Math.PI / 2, 0, 0]}>
          <cylinderGeometry args={[0.06, 0.06, 0.8, 6]} />
          <meshStandardMaterial color="#444455" metalness={0.9} roughness={0.15} />
        </mesh>

        {/* Feed tip sphere */}
        <mesh ref={feedTipRef} position={[0, 0.1, 1.2]}>
          <sphereGeometry args={[0.08, 8, 6]} />
          <meshStandardMaterial
            color="#ff4444"
            emissive="#ff4444"
            emissiveIntensity={1.5}
            toneMapped={false}
          />
        </mesh>
      </group>

      {/* Base indicator light */}
      <mesh ref={baseLightRef} position={[0.8, 1.1, 0.8]}>
        <sphereGeometry args={[0.06, 6, 6]} />
        <meshStandardMaterial
          color="#ff2222"
          emissive="#ff2222"
          emissiveIntensity={1.5}
          toneMapped={false}
        />
      </mesh>
    </group>
  )
}
