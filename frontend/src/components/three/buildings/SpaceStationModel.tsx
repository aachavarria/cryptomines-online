import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Group, Mesh, MeshStandardMaterial } from 'three'

interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
}

/**
 * Space Station: massive orbital defense platform (3x3).
 * Circular multi-ring structure with central command core,
 * two concentric defense rings, turret hardpoints, and shield dome.
 */
export default function SpaceStationModel({
  position = [0, 0, 0],
  scale = 1,
  level = 1,
}: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const innerRingRef = useRef<Mesh>(null)
  const outerRingRef = useRef<Mesh>(null)
  const shieldRef = useRef<Mesh>(null)
  const coreDomeRef = useRef<Mesh>(null)

  const color = '#44ccff'

  useFrame((state) => {
    const t = state.clock.elapsedTime
    // Inner ring rotates clockwise
    if (innerRingRef.current) {
      innerRingRef.current.rotation.y = t * (Math.PI * 2 / 12)
    }
    // Outer ring rotates counter-clockwise
    if (outerRingRef.current) {
      outerRingRef.current.rotation.y = -t * (Math.PI * 2 / 20)
    }
    // Shield dome barely visible pulse
    if (shieldRef.current) {
      const mat = shieldRef.current.material as MeshStandardMaterial
      mat.opacity = 0.05 + Math.sin(t * (Math.PI * 2 / 5)) * 0.05
    }
    // Core dome energy pulse
    if (coreDomeRef.current) {
      const mat = coreDomeRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 0.45 + Math.sin(t * (Math.PI * 2 / 3)) * 0.15
    }
  })

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Base platform */}
      <mesh position={[0, 0.15, 0]} castShadow>
        <cylinderGeometry args={[5.0, 5.0, 0.3, 16]} />
        <meshStandardMaterial
          color="#1a2a3a"
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* Central core */}
      <mesh position={[0, 2.8, 0]} castShadow>
        <cylinderGeometry args={[1.5, 1.5, 5.0, 8]} />
        <meshStandardMaterial
          color="#1a2a3a"
          emissive={color}
          emissiveIntensity={0.15}
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* Core dome */}
      <mesh ref={coreDomeRef} position={[0, 5.3, 0]}>
        <sphereGeometry args={[1.6, 10, 8, 0, Math.PI * 2, 0, Math.PI / 2]} />
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={0.4}
          transparent
          opacity={0.3}
        />
      </mesh>

      {/* Inner ring */}
      <group ref={innerRingRef} position={[0, 2.0, 0]}>
        <mesh rotation={[Math.PI / 2, 0, 0]}>
          <torusGeometry args={[2.8, 0.4, 6, 16]} />
          <meshStandardMaterial
            color="#2a3a4a"
            emissive={color}
            emissiveIntensity={0.1}
            metalness={0.85}
            roughness={0.2}
          />
        </mesh>
        {/* 6 support struts from core to inner ring */}
        {Array.from({ length: 6 }).map((_, i) => {
          const angle = (i * Math.PI * 2) / 6
          return (
            <mesh
              key={`inner-strut-${i}`}
              position={[Math.cos(angle) * 1.4, 0, Math.sin(angle) * 1.4]}
              rotation={[0, -angle + Math.PI / 2, 0]}
              castShadow
            >
              <boxGeometry args={[2.5, 0.15, 0.15]} />
              <meshStandardMaterial
                color="#2a3a4a"
                metalness={0.85}
                roughness={0.2}
              />
            </mesh>
          )
        })}
      </group>

      {/* Outer ring */}
      <group ref={outerRingRef} position={[0, 1.0, 0]}>
        <mesh rotation={[Math.PI / 2, 0, 0]}>
          <torusGeometry args={[4.5, 0.5, 6, 20]} />
          <meshStandardMaterial
            color="#2a3a4a"
            emissive={color}
            emissiveIntensity={0.1}
            metalness={0.85}
            roughness={0.2}
          />
        </mesh>
        {/* 6 outer struts from inner to outer ring */}
        {Array.from({ length: 6 }).map((_, i) => {
          const angle = (i * Math.PI * 2) / 6
          return (
            <mesh
              key={`outer-strut-${i}`}
              position={[Math.cos(angle) * 3.65, 0, Math.sin(angle) * 3.65]}
              rotation={[0, -angle + Math.PI / 2, 0]}
              castShadow
            >
              <boxGeometry args={[2.0, 0.15, 0.15]} />
              <meshStandardMaterial
                color="#2a3a4a"
                metalness={0.85}
                roughness={0.2}
              />
            </mesh>
          )
        })}
        {/* 8 turret hardpoints on outer ring */}
        {Array.from({ length: 8 }).map((_, i) => {
          const angle = (i * Math.PI * 2) / 8
          return (
            <group
              key={`turret-${i}`}
              position={[Math.cos(angle) * 4.5, 0.3, Math.sin(angle) * 4.5]}
              rotation={[0, -angle, 0]}
            >
              {/* Turret base */}
              <mesh castShadow>
                <cylinderGeometry args={[0.2, 0.2, 0.3, 6]} />
                <meshStandardMaterial
                  color="#3a4a5a"
                  metalness={0.9}
                  roughness={0.15}
                />
              </mesh>
              {/* Turret barrel */}
              <mesh position={[0.3, 0.1, 0]} rotation={[0, 0, Math.PI / 2]}>
                <cylinderGeometry args={[0.06, 0.06, 0.5, 6]} />
                <meshStandardMaterial
                  color="#555566"
                  metalness={0.9}
                  roughness={0.15}
                />
              </mesh>
              {/* Turret ready light */}
              <mesh position={[0, 0.2, 0]}>
                <sphereGeometry args={[0.04, 6, 6]} />
                <meshStandardMaterial
                  color={color}
                  emissive={color}
                  emissiveIntensity={1.5}
                  toneMapped={false}
                />
              </mesh>
            </group>
          )
        })}
      </group>

      {/* Shield dome - faint energy shield */}
      <mesh ref={shieldRef} position={[0, 0, 0]}>
        <sphereGeometry args={[5.0, 12, 8, 0, Math.PI * 2, 0, Math.PI / 2]} />
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={0.2}
          transparent
          opacity={0.1}
          side={2}
        />
      </mesh>
    </group>
  )
}
