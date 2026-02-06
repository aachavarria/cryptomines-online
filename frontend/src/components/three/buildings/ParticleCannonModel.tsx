import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Group, Mesh, MeshStandardMaterial } from 'three'

interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
}

/**
 * Particle Cannon: heavy beam weapon emplacement (1x2).
 * Long tracked base with massive barrel, energy capacitor bank,
 * and particle focusing ring at the barrel tip.
 */
export default function ParticleCannonModel({
  position = [0, 0, 0],
  scale = 1,
  level = 1,
}: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const turretRef = useRef<Group>(null)
  const emitterRef = useRef<Mesh>(null)
  const cap1Ref = useRef<Mesh>(null)
  const cap2Ref = useRef<Mesh>(null)
  const cap3Ref = useRef<Mesh>(null)
  const shroudRef = useRef<Mesh>(null)

  const color = '#44ccff'

  useFrame((state) => {
    const t = state.clock.elapsedTime
    // Turret + barrel slow aiming sweep
    if (turretRef.current) {
      turretRef.current.rotation.y = Math.sin(t * (Math.PI * 2 / 10)) * 0.6
    }
    // Emitter ring rapid pulse - charging effect
    if (emitterRef.current) {
      const mat = emitterRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 1.0 + Math.sin(t * (Math.PI * 2 / 0.5)) * 0.75
    }
    // Capacitors pulse in sequence
    const capRefs = [cap1Ref, cap2Ref, cap3Ref]
    capRefs.forEach((ref, i) => {
      if (ref.current) {
        const mat = ref.current.material as MeshStandardMaterial
        const phase = (t - i * 0.67) % 2
        mat.emissiveIntensity = 0.1 + Math.max(0, Math.sin(phase * Math.PI)) * 0.5
      }
    })
    // Barrel shroud glows briefly every ~6s
    if (shroudRef.current) {
      const mat = shroudRef.current.material as MeshStandardMaterial
      const cycle = t % 6
      mat.emissiveIntensity = cycle < 0.5 ? 0.5 + (0.5 - cycle) * 1.5 : 0.2
    }
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Tracked base - long, low chassis */}
      <mesh position={[0, 0.4, 0]} castShadow>
        <boxGeometry args={[3.0, 0.8, 6.5]} />
        <meshStandardMaterial
          color="#1a2a3a"
          emissive={color}
          emissiveIntensity={0.08}
          metalness={0.8}
          roughness={0.3}
        />
      </mesh>

      {/* Track wheels - 2 per side */}
      {[[-1.3, 0.3, -2.0], [-1.3, 0.3, 2.0], [1.3, 0.3, -2.0], [1.3, 0.3, 2.0]].map(
        ([x, y, z], i) => (
          <mesh key={`wheel-${i}`} position={[x, y, z]} rotation={[0, 0, Math.PI / 2]}>
            <cylinderGeometry args={[0.4, 0.4, 0.3, 8]} />
            <meshStandardMaterial
              color="#2a3a4a"
              metalness={0.85}
              roughness={0.3}
            />
          </mesh>
        )
      )}

      {/* Turret pivot + barrel group */}
      <group ref={turretRef} position={[0, 1.3, 0]}>
        {/* Turret pivot base */}
        <mesh castShadow>
          <cylinderGeometry args={[1.0, 1.0, 1.0, 8]} />
          <meshStandardMaterial
            color="#2a3a4a"
            metalness={0.85}
            roughness={0.2}
          />
        </mesh>

        {/* Main barrel - angled 15deg upward */}
        <group rotation={[-Math.PI * 15 / 180, 0, 0]} position={[0, 0.3, 0.5]}>
          {/* Barrel shroud at base */}
          <mesh ref={shroudRef} position={[0, 0, 0]} rotation={[Math.PI / 2, 0, 0]}>
            <cylinderGeometry args={[0.6, 0.6, 1.0, 8]} />
            <meshStandardMaterial
              color="#2a3a4a"
              emissive={color}
              emissiveIntensity={0.2}
              metalness={0.85}
              roughness={0.2}
            />
          </mesh>

          {/* Main barrel */}
          <mesh position={[0, 0, 2.0]} rotation={[Math.PI / 2, 0, 0]} castShadow>
            <cylinderGeometry args={[0.25, 0.4, 4.0, 8]} />
            <meshStandardMaterial
              color="#3a4a5a"
              metalness={0.9}
              roughness={0.15}
            />
          </mesh>

          {/* Emitter ring at barrel tip */}
          <mesh ref={emitterRef} position={[0, 0, 4.0]} rotation={[Math.PI / 2, 0, 0]}>
            <torusGeometry args={[0.3, 0.05, 6, 12]} />
            <meshStandardMaterial
              color="#44ffff"
              emissive="#44ffff"
              emissiveIntensity={1.5}
              toneMapped={false}
            />
          </mesh>
        </group>
      </group>

      {/* Capacitor bank - 3 cylinders along one side */}
      {[-1.5, 0, 1.5].map((z, i) => {
        const refs = [cap1Ref, cap2Ref, cap3Ref]
        return (
          <mesh
            key={`cap-${i}`}
            ref={refs[i]}
            position={[-1.2, 1.0, z]}
            castShadow
          >
            <cylinderGeometry args={[0.25, 0.25, 1.2, 6]} />
            <meshStandardMaterial
              color="#2a4a5a"
              emissive={color}
              emissiveIntensity={0.3}
              metalness={0.8}
              roughness={0.2}
            />
          </mesh>
        )
      })}
    </group>
  )
}
