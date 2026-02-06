import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Group, Mesh, MeshStandardMaterial } from 'three'

interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
}

/**
 * Technology Center: research laboratory with satellite dish (core, 3x2).
 * Rectangular building with angular roof, parabolic dish, holographic ring,
 * and pulsing data cores.
 */
export default function TechnologyCenterModel({ position = [0, 0, 0], scale = 1, level = 1 }: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const holoRingRef = useRef<Mesh>(null)
  const dishGroupRef = useRef<Group>(null)
  const dataCore0Ref = useRef<Mesh>(null)
  const dataCore1Ref = useRef<Mesh>(null)
  const dataCore2Ref = useRef<Mesh>(null)

  useFrame((state) => {
    const t = state.clock.elapsedTime
    // Holographic ring rotates (Y-axis, period ~4s) and bobs
    if (holoRingRef.current) {
      holoRingRef.current.rotation.y = t * (Math.PI * 2 / 4)
      holoRingRef.current.position.y = 6.0 + Math.sin(t * 1.5) * 0.15
    }
    // Dish rotates slowly (period ~20s)
    if (dishGroupRef.current) {
      dishGroupRef.current.rotation.y = t * (Math.PI * 2 / 20)
    }
    // Data cores pulse sequentially (stagger 0.3-0.8)
    const cores = [dataCore0Ref, dataCore1Ref, dataCore2Ref]
    cores.forEach((ref, i) => {
      if (ref.current) {
        const mat = ref.current.material as MeshStandardMaterial
        mat.emissiveIntensity = 0.55 + Math.sin(t * 2 + i * (Math.PI * 2 / 3)) * 0.25
      }
    })
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Main body */}
      <mesh position={[0, 1.75, 0]} castShadow>
        <boxGeometry args={[8.0, 3.5, 5.5]} />
        <meshStandardMaterial
          color="#1a2a44"
          emissive="#4488ff"
          emissiveIntensity={0.12}
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* Angular roof - left panel */}
      <mesh position={[-1.2, 3.7, 0]} rotation={[0, 0, Math.PI * 15 / 180]} castShadow>
        <boxGeometry args={[4.5, 0.3, 5.5]} />
        <meshStandardMaterial
          color="#2a3a55"
          metalness={0.9}
          roughness={0.2}
        />
      </mesh>
      {/* Angular roof - right panel */}
      <mesh position={[1.2, 3.7, 0]} rotation={[0, 0, -Math.PI * 15 / 180]} castShadow>
        <boxGeometry args={[4.5, 0.3, 5.5]} />
        <meshStandardMaterial
          color="#2a3a55"
          metalness={0.9}
          roughness={0.2}
        />
      </mesh>

      {/* Satellite dish assembly (slowly rotating) */}
      <group ref={dishGroupRef} position={[0, 4.3, 0]}>
        {/* Dish support */}
        <mesh>
          <cylinderGeometry args={[0.15, 0.15, 1.5, 6]} />
          <meshStandardMaterial
            color="#3a4a66"
            metalness={0.9}
            roughness={0.2}
          />
        </mesh>
        {/* Parabolic dish (inverted cone) */}
        <mesh position={[0, 0.8, 0]} rotation={[Math.PI, 0, 0]}>
          <coneGeometry args={[1.8, 0.5, 10, 1, true]} />
          <meshStandardMaterial
            color="#3a4a66"
            metalness={0.95}
            roughness={0.15}
            side={2}
          />
        </mesh>
      </group>

      {/* Holographic ring floating above dish */}
      <mesh ref={holoRingRef} position={[0, 6.0, 0]}>
        <torusGeometry args={[1.0, 0.06, 8, 24]} />
        <meshStandardMaterial
          color="#66ccff"
          emissive="#66ccff"
          emissiveIntensity={1.2}
          transparent
          opacity={0.5}
          toneMapped={false}
        />
      </mesh>

      {/* 3 data cores along one side */}
      {[-1.8, 0, 1.8].map((z, i) => {
        const refs = [dataCore0Ref, dataCore1Ref, dataCore2Ref]
        return (
          <mesh key={`core-${i}`} ref={refs[i]} position={[-4.2, 0.8, z]}>
            <cylinderGeometry args={[0.3, 0.3, 1.2, 6]} />
            <meshStandardMaterial
              color="#224488"
              emissive="#4488ff"
              emissiveIntensity={0.6}
              metalness={0.8}
              roughness={0.25}
            />
          </mesh>
        )
      })}
    </group>
  )
}
