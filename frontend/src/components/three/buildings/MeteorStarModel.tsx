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
 * Meteor Star: point-defense laser turret (1x1).
 * Low-profile hexagonal armored base with rotating turret housing
 * twin-barrel laser emitters and targeting sensor.
 */
export default function MeteorStarModel({
  position = [0, 0, 0],
  scale = 1,
  level = 1,
  animate = true,
}: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const turretRef = useRef<Group>(null)
  const sensorRef = useRef<Mesh>(null)
  const leftTipRef = useRef<Mesh>(null)
  const rightTipRef = useRef<Mesh>(null)

  const color = '#44ccff'

  useFrame((state) => {
    if (!animate) return
    const t = state.clock.elapsedTime
    // Turret rotates scanning for targets
    if (turretRef.current) {
      turretRef.current.rotation.y = t * (Math.PI * 2 / 5)
    }
    // Barrel tips glow with slight flicker
    if (leftTipRef.current) {
      const mat = leftTipRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 1.8 + Math.sin(t * 8) * 0.4
    }
    if (rightTipRef.current) {
      const mat = rightTipRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 1.8 + Math.cos(t * 8) * 0.4
    }
    // Sensor blinks
    if (sensorRef.current) {
      const mat = sensorRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = Math.sin(t * (Math.PI * 2 / 1.5)) > 0 ? 1.5 : 0.2
    }
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Hexagonal armored base */}
      <mesh position={[0, 0.5, 0]} castShadow>
        <cylinderGeometry args={[1.5, 1.5, 1.0, 6]} />
        <meshStandardMaterial
          color="#1a2a3a"
          emissive={color}
          emissiveIntensity={0.1}
          metalness={0.85}
          roughness={0.25}
        />
      </mesh>

      {/* Armor plates around base */}
      {[0, Math.PI * 2 / 3, Math.PI * 4 / 3].map((angle, i) => (
        <mesh
          key={`armor-${i}`}
          position={[
            Math.cos(angle) * 1.3,
            0.5,
            Math.sin(angle) * 1.3,
          ]}
          rotation={[0, -angle, 0]}
          castShadow
        >
          <boxGeometry args={[0.8, 0.6, 0.05]} />
          <meshStandardMaterial
            color="#2a3a4a"
            metalness={0.8}
            roughness={0.2}
          />
        </mesh>
      ))}

      {/* Rotating turret group */}
      <group ref={turretRef} position={[0, 1.2, 0]}>
        {/* Turret housing - squashed half sphere */}
        <mesh scale={[1, 0.7, 1]}>
          <sphereGeometry args={[0.8, 8, 6, 0, Math.PI * 2, 0, Math.PI / 2]} />
          <meshStandardMaterial
            color="#2a3a4a"
            emissive={color}
            emissiveIntensity={0.15}
            metalness={0.85}
            roughness={0.2}
          />
        </mesh>

        {/* Twin barrels - angled 20deg upward */}
        <group rotation={[-Math.PI * 20 / 180, 0, 0]} position={[0, 0.2, 0.4]}>
          {/* Left barrel */}
          <mesh position={[-0.2, 0, 0]} rotation={[Math.PI / 2, 0, 0]}>
            <cylinderGeometry args={[0.08, 0.08, 1.5, 6]} />
            <meshStandardMaterial
              color="#555566"
              metalness={0.9}
              roughness={0.15}
            />
          </mesh>
          {/* Left barrel tip */}
          <mesh ref={leftTipRef} position={[-0.2, 0, 0.75]}>
            <sphereGeometry args={[0.05, 6, 6]} />
            <meshStandardMaterial
              color="#44ffff"
              emissive="#44ffff"
              emissiveIntensity={2.0}
              toneMapped={false}
            />
          </mesh>
          {/* Right barrel */}
          <mesh position={[0.2, 0, 0]} rotation={[Math.PI / 2, 0, 0]}>
            <cylinderGeometry args={[0.08, 0.08, 1.5, 6]} />
            <meshStandardMaterial
              color="#555566"
              metalness={0.9}
              roughness={0.15}
            />
          </mesh>
          {/* Right barrel tip */}
          <mesh ref={rightTipRef} position={[0.2, 0, 0.75]}>
            <sphereGeometry args={[0.05, 6, 6]} />
            <meshStandardMaterial
              color="#44ffff"
              emissive="#44ffff"
              emissiveIntensity={2.0}
              toneMapped={false}
            />
          </mesh>
        </group>

        {/* Targeting sensor */}
        <mesh ref={sensorRef} position={[0, 0.5, 0.1]}>
          <boxGeometry args={[0.2, 0.1, 0.1]} />
          <meshStandardMaterial
            color="#ff4444"
            emissive="#ff4444"
            emissiveIntensity={1.0}
            toneMapped={false}
          />
        </mesh>
      </group>
    </group>
  )
}
