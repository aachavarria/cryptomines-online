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
 * Anti-Aircraft Gun: rapid-fire flak turret (1x1).
 * Compact rotating platform with quad barrels pointing skyward,
 * ammo feed drum, and visible rotation bearing ring.
 */
export default function AntiAircraftGunModel({
  position = [0, 0, 0],
  scale = 1,
  level = 1,
  animate = true,
}: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const turretRef = useRef<Group>(null)
  const drumRef = useRef<Mesh>(null)
  const tip1Ref = useRef<Mesh>(null)
  const tip2Ref = useRef<Mesh>(null)
  const tip3Ref = useRef<Mesh>(null)
  const tip4Ref = useRef<Mesh>(null)

  const color = '#44ccff'

  useFrame((state) => {
    if (!animate) return
    const t = state.clock.elapsedTime
    // Turret rotates rapidly - fast scanning
    if (turretRef.current) {
      turretRef.current.rotation.y = t * (Math.PI * 2 / 3)
    }
    // Barrel tips flicker simulating muzzle flash
    const tipRefs = [tip1Ref, tip2Ref, tip3Ref, tip4Ref]
    tipRefs.forEach((ref, i) => {
      if (ref.current) {
        const mat = ref.current.material as MeshStandardMaterial
        mat.emissiveIntensity = 0.5 + Math.abs(Math.sin(t * 12 + i * 1.7)) * 1.5
      }
    })
    // Ammo drum rotating element
    if (drumRef.current) {
      drumRef.current.rotation.y = t * (Math.PI * 2 / 1)
    }
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Base platform */}
      <mesh position={[0, 0.25, 0]} castShadow>
        <cylinderGeometry args={[1.2, 1.2, 0.5, 8]} />
        <meshStandardMaterial
          color="#1a2a3a"
          emissive={color}
          emissiveIntensity={0.08}
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* Rotation ring - visible bearing */}
      <mesh position={[0, 0.5, 0]} rotation={[Math.PI / 2, 0, 0]}>
        <torusGeometry args={[1.0, 0.08, 6, 12]} />
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={0.3}
          metalness={0.85}
          roughness={0.2}
        />
      </mesh>

      {/* Rotating turret group */}
      <group ref={turretRef} position={[0, 1.0, 0]}>
        {/* Turret body */}
        <mesh castShadow>
          <boxGeometry args={[1.0, 0.8, 1.0]} />
          <meshStandardMaterial
            color="#2a3a4a"
            metalness={0.85}
            roughness={0.2}
          />
        </mesh>

        {/* Barrel cluster mount */}
        <mesh position={[0, 0.5, 0.2]} rotation={[-Math.PI * 60 / 180, 0, 0]}>
          <cylinderGeometry args={[0.25, 0.25, 0.3, 6]} />
          <meshStandardMaterial
            color="#3a4a5a"
            metalness={0.85}
            roughness={0.2}
          />
        </mesh>

        {/* Quad barrels - angled 60deg upward, 2x2 grid */}
        <group position={[0, 0.6, 0.3]} rotation={[-Math.PI * 60 / 180, 0, 0]}>
          {[[-0.15, -0.15], [0.15, -0.15], [-0.15, 0.15], [0.15, 0.15]].map(
            ([xOff, zOff], i) => {
              const tipRefs = [tip1Ref, tip2Ref, tip3Ref, tip4Ref]
              return (
                <group key={`barrel-${i}`}>
                  <mesh position={[xOff, 1.0, zOff]} rotation={[Math.PI / 2, 0, 0]}>
                    <cylinderGeometry args={[0.06, 0.06, 2.0, 6]} />
                    <meshStandardMaterial
                      color="#555566"
                      metalness={0.95}
                      roughness={0.1}
                    />
                  </mesh>
                  {/* Barrel tip - muzzle heat */}
                  <mesh ref={tipRefs[i]} position={[xOff, 2.0, zOff]}>
                    <sphereGeometry args={[0.04, 6, 6]} />
                    <meshStandardMaterial
                      color="#ffaa44"
                      emissive="#ffaa44"
                      emissiveIntensity={1.0}
                      toneMapped={false}
                    />
                  </mesh>
                </group>
              )
            }
          )}
        </group>

        {/* Ammo drum */}
        <mesh position={[0.6, 0.1, 0]}>
          <cylinderGeometry args={[0.3, 0.3, 0.5, 8]} />
          <meshStandardMaterial
            color="#3a4a5a"
            metalness={0.8}
            roughness={0.25}
          />
        </mesh>

        {/* Ammo drum rotating indicator */}
        <mesh ref={drumRef} position={[0.6, 0.4, 0]}>
          <cylinderGeometry args={[0.08, 0.08, 0.1, 4]} />
          <meshStandardMaterial
            color={color}
            emissive={color}
            emissiveIntensity={0.6}
          />
        </mesh>

        {/* Ammo feed belt */}
        <mesh position={[0.35, 0.3, 0]}>
          <boxGeometry args={[0.5, 0.08, 0.08]} />
          <meshStandardMaterial
            color="#3a4a5a"
            metalness={0.8}
            roughness={0.3}
          />
        </mesh>
      </group>
    </group>
  )
}
