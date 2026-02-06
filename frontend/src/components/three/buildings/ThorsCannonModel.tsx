import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Group, Mesh, MeshStandardMaterial } from 'three'

interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
}

/**
 * Thor's Cannon: orbital strike electromagnetic railgun (2x2).
 * Two parallel vertical rails with projectile cradle, power coils,
 * and lightning-like energy arcs. The most imposing defense weapon.
 */
export default function ThorsCannonModel({
  position = [0, 0, 0],
  scale = 1,
  level = 1,
}: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const cradleRef = useRef<Mesh>(null)
  const arc1Ref = useRef<Mesh>(null)
  const arc2Ref = useRef<Mesh>(null)
  const arc3Ref = useRef<Mesh>(null)
  const coil1Ref = useRef<Mesh>(null)
  const coil2Ref = useRef<Mesh>(null)
  const coil3Ref = useRef<Mesh>(null)
  const leftCapRef = useRef<Mesh>(null)
  const rightCapRef = useRef<Mesh>(null)

  const color = '#44ccff'

  useFrame((state) => {
    const t = state.clock.elapsedTime
    // Energy arcs flicker rapidly like lightning
    const arcRefs = [arc1Ref, arc2Ref, arc3Ref]
    arcRefs.forEach((ref, i) => {
      if (ref.current) {
        const mat = ref.current.material as MeshStandardMaterial
        const flicker = Math.sin(t * 15 + i * 3) * Math.cos(t * 23 + i * 7)
        mat.opacity = 0.3 + Math.abs(flicker) * 0.5
        mat.emissiveIntensity = 1.0 + Math.abs(flicker) * 2.0
      }
    })
    // Power coils pulse upward in sequence
    const coilRefs = [coil1Ref, coil2Ref, coil3Ref]
    coilRefs.forEach((ref, i) => {
      if (ref.current) {
        const mat = ref.current.material as MeshStandardMaterial
        const phase = (t * (1 / 1.5) - i * 0.33) % 1
        mat.emissiveIntensity = 0.3 + Math.max(0, Math.sin(phase * Math.PI)) * 0.5
      }
    })
    // Rail cap electrodes pulse together
    if (leftCapRef.current && rightCapRef.current) {
      const intensity = 0.5 + Math.sin(t * (Math.PI * 2 / 2)) * 0.5
      ;(leftCapRef.current.material as MeshStandardMaterial).emissiveIntensity = intensity
      ;(rightCapRef.current.material as MeshStandardMaterial).emissiveIntensity = intensity
    }
    // Projectile cradle bobs slightly - magnetic levitation
    if (cradleRef.current) {
      cradleRef.current.position.y = 3.0 + Math.sin(t * (Math.PI * 2 / 3)) * 0.1
    }
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Power base - heavy foundation */}
      <mesh position={[0, 0.75, 0]} castShadow>
        <cylinderGeometry args={[2.5, 2.5, 1.5, 8]} />
        <meshStandardMaterial
          color="#1a2a3a"
          emissive={color}
          emissiveIntensity={0.1}
          metalness={0.85}
          roughness={0.25}
        />
      </mesh>

      {/* Base vents */}
      {[0, Math.PI / 2, Math.PI, Math.PI * 3 / 2].map((angle, i) => (
        <mesh
          key={`vent-${i}`}
          position={[Math.cos(angle) * 2.3, 0.75, Math.sin(angle) * 2.3]}
          rotation={[0, -angle, 0]}
        >
          <boxGeometry args={[0.6, 0.2, 0.05]} />
          <meshStandardMaterial
            color={color}
            emissive={color}
            emissiveIntensity={0.3}
          />
        </mesh>
      ))}

      {/* Left rail */}
      <mesh position={[-1.0, 4.5, 0]} castShadow>
        <boxGeometry args={[0.4, 6.0, 0.4]} />
        <meshStandardMaterial
          color="#3a4a5a"
          metalness={0.95}
          roughness={0.1}
        />
      </mesh>

      {/* Right rail */}
      <mesh position={[1.0, 4.5, 0]} castShadow>
        <boxGeometry args={[0.4, 6.0, 0.4]} />
        <meshStandardMaterial
          color="#3a4a5a"
          metalness={0.95}
          roughness={0.1}
        />
      </mesh>

      {/* Left rail top cap - electrode terminal */}
      <mesh ref={leftCapRef} position={[-1.0, 7.65, 0]}>
        <boxGeometry args={[0.6, 0.3, 0.6]} />
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={0.8}
          toneMapped={false}
        />
      </mesh>

      {/* Right rail top cap - electrode terminal */}
      <mesh ref={rightCapRef} position={[1.0, 7.65, 0]}>
        <boxGeometry args={[0.6, 0.3, 0.6]} />
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={0.8}
          toneMapped={false}
        />
      </mesh>

      {/* Projectile cradle between rails */}
      <mesh ref={cradleRef} position={[0, 3.0, 0]} castShadow>
        <boxGeometry args={[0.6, 0.8, 0.6]} />
        <meshStandardMaterial
          color="#2a3a4a"
          emissive="#66ddff"
          emissiveIntensity={0.3}
          metalness={0.85}
          roughness={0.2}
        />
      </mesh>

      {/* Power coils - 3 stacked around the rails */}
      {[0.8, 1.6, 2.4].map((y, i) => {
        const coilRefs = [coil1Ref, coil2Ref, coil3Ref]
        return (
          <mesh
            key={`coil-${i}`}
            ref={coilRefs[i]}
            position={[0, y, 0]}
            rotation={[Math.PI / 2, 0, 0]}
          >
            <torusGeometry args={[1.5, 0.1, 6, 12]} />
            <meshStandardMaterial
              color="#44aacc"
              emissive={color}
              emissiveIntensity={0.5}
              metalness={0.8}
              roughness={0.2}
            />
          </mesh>
        )
      })}

      {/* Energy arc connectors between rails at coil heights */}
      {[0.8, 1.6, 2.4].map((y, i) => {
        const arcRefs = [arc1Ref, arc2Ref, arc3Ref]
        return (
          <mesh
            key={`arc-${i}`}
            ref={arcRefs[i]}
            position={[0, y, 0]}
          >
            <boxGeometry args={[1.6, 0.04, 0.04]} />
            <meshStandardMaterial
              color="#88eeff"
              emissive="#88eeff"
              emissiveIntensity={2.0}
              transparent
              opacity={0.6}
              toneMapped={false}
            />
          </mesh>
        )
      })}
    </group>
  )
}
