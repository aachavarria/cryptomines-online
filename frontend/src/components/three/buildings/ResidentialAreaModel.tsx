import { useRef, useMemo } from 'react'
import { useFrame } from '@react-three/fiber'
import * as THREE from 'three'
import type { Group, Mesh, MeshStandardMaterial } from 'three'

interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
  animate?: boolean
}

/**
 * Residential Area (resource, 2x2): Cluster of dome-shaped habitat pods
 * connected by walkways. Warm interior lighting visible through windows.
 */
export default function ResidentialAreaModel({ position = [0, 0, 0], scale = 1, level = 1, animate = true }: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const windowRefs = useRef<Mesh[]>([])
  const orbitLightRef = useRef<Mesh>(null)

  // Dome geometries
  const largeDomeGeo = useMemo(() => new THREE.SphereGeometry(1.5, 10, 8, 0, Math.PI * 2, 0, Math.PI / 2), [])
  const medDomeGeo = useMemo(() => new THREE.SphereGeometry(1.1, 10, 8, 0, Math.PI * 2, 0, Math.PI / 2), [])
  const smallDomeGeoA = useMemo(() => new THREE.SphereGeometry(0.8, 8, 6, 0, Math.PI * 2, 0, Math.PI / 2), [])
  const smallDomeGeoB = useMemo(() => new THREE.SphereGeometry(0.7, 8, 6, 0, Math.PI * 2, 0, Math.PI / 2), [])

  useFrame((state) => {
    if (!animate) return
    const t = state.clock.elapsedTime
    // Window flicker
    windowRefs.current.forEach((mesh, i) => {
      if (mesh) {
        const mat = mesh.material as MeshStandardMaterial
        mat.emissiveIntensity = 0.6 + Math.sin(t * 2.3 + i * 1.7) * 0.2
      }
    })
    // Orbiting patrol light around large dome
    if (orbitLightRef.current) {
      const angle = t * (Math.PI * 2 / 12)
      orbitLightRef.current.position.x = -0.5 + Math.cos(angle) * 1.6
      orbitLightRef.current.position.z = -0.3 + Math.sin(angle) * 1.6
    }
  })

  const addWindowRef = (el: Mesh | null, index: number) => {
    if (el) windowRefs.current[index] = el
  }

  const windowMat = (
    <meshStandardMaterial
      color="#ffdd88"
      emissive="#ffdd88"
      emissiveIntensity={0.8}
      transparent
      opacity={0.6}
    />
  )

  const domeMat = (
    <meshStandardMaterial
      color="#2a3a2a"
      emissive="#22cc66"
      emissiveIntensity={0.08}
      metalness={0.6}
      roughness={0.35}
    />
  )

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Large dome */}
      <mesh geometry={largeDomeGeo} position={[-0.5, 0, -0.3]} castShadow>
        {domeMat}
      </mesh>
      {/* Base ring - large dome */}
      <mesh position={[-0.5, 0.1, -0.3]} rotation={[Math.PI / 2, 0, 0]}>
        <torusGeometry args={[1.5, 0.06, 6, 12]} />
        <meshStandardMaterial color="#22cc66" emissive="#22cc66" emissiveIntensity={0.4} />
      </mesh>

      {/* Medium dome */}
      <mesh geometry={medDomeGeo} position={[1.3, 0, 0.5]} castShadow>
        {domeMat}
      </mesh>
      <mesh position={[1.3, 0.1, 0.5]} rotation={[Math.PI / 2, 0, 0]}>
        <torusGeometry args={[1.1, 0.06, 6, 10]} />
        <meshStandardMaterial color="#22cc66" emissive="#22cc66" emissiveIntensity={0.4} />
      </mesh>

      {/* Small dome A */}
      <mesh geometry={smallDomeGeoA} position={[0.3, 0, 1.8]} castShadow>
        {domeMat}
      </mesh>
      <mesh position={[0.3, 0.1, 1.8]} rotation={[Math.PI / 2, 0, 0]}>
        <torusGeometry args={[0.8, 0.06, 6, 8]} />
        <meshStandardMaterial color="#22cc66" emissive="#22cc66" emissiveIntensity={0.4} />
      </mesh>

      {/* Small dome B */}
      <mesh geometry={smallDomeGeoB} position={[-1.5, 0, 1.2]} castShadow>
        {domeMat}
      </mesh>
      <mesh position={[-1.5, 0.1, 1.2]} rotation={[Math.PI / 2, 0, 0]}>
        <torusGeometry args={[0.7, 0.06, 6, 8]} />
        <meshStandardMaterial color="#22cc66" emissive="#22cc66" emissiveIntensity={0.4} />
      </mesh>

      {/* Walkway: large dome -> medium dome */}
      <mesh position={[0.4, 0.3, 0.1]} rotation={[0, -0.4, 0]} castShadow>
        <boxGeometry args={[1.8, 0.4, 0.4]} />
        <meshStandardMaterial color="#3a4a3a" metalness={0.75} roughness={0.3} />
      </mesh>

      {/* Walkway: medium dome -> small dome A */}
      <mesh position={[0.8, 0.3, 1.15]} rotation={[0, 0.7, 0]} castShadow>
        <boxGeometry args={[1.4, 0.4, 0.4]} />
        <meshStandardMaterial color="#3a4a3a" metalness={0.75} roughness={0.3} />
      </mesh>

      {/* Walkway: large dome -> small dome B */}
      <mesh position={[-1.0, 0.3, 0.45]} rotation={[0, 0.8, 0]} castShadow>
        <boxGeometry args={[1.2, 0.4, 0.4]} />
        <meshStandardMaterial color="#3a4a3a" metalness={0.75} roughness={0.3} />
      </mesh>

      {/* Windows - large dome */}
      <mesh ref={(el) => addWindowRef(el, 0)} position={[-0.5, 0.7, 1.1]}>
        <boxGeometry args={[0.3, 0.3, 0.02]} />
        {windowMat}
      </mesh>
      <mesh ref={(el) => addWindowRef(el, 1)} position={[-1.8, 0.6, -0.3]}>
        <boxGeometry args={[0.02, 0.3, 0.3]} />
        {windowMat}
      </mesh>

      {/* Windows - medium dome */}
      <mesh ref={(el) => addWindowRef(el, 2)} position={[1.3, 0.5, 1.5]}>
        <boxGeometry args={[0.3, 0.3, 0.02]} />
        {windowMat}
      </mesh>
      <mesh ref={(el) => addWindowRef(el, 3)} position={[2.2, 0.5, 0.5]}>
        <boxGeometry args={[0.02, 0.3, 0.3]} />
        {windowMat}
      </mesh>

      {/* Windows - small domes */}
      <mesh ref={(el) => addWindowRef(el, 4)} position={[0.3, 0.4, 2.5]}>
        <boxGeometry args={[0.25, 0.25, 0.02]} />
        {windowMat}
      </mesh>
      <mesh ref={(el) => addWindowRef(el, 5)} position={[-1.5, 0.35, 1.85]}>
        <boxGeometry args={[0.25, 0.25, 0.02]} />
        {windowMat}
      </mesh>

      {/* Orbiting patrol light */}
      <mesh ref={orbitLightRef} position={[-0.5, 0.15, 1.3]}>
        <sphereGeometry args={[0.06, 6, 6]} />
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
