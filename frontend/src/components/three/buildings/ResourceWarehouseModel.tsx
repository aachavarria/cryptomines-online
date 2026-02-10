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
 * Resource Warehouse (resource, 3x2): Long storage depot with barrel-vault roof,
 * loading dock doors, stacked cargo, and a conveyor track on the roof ridge.
 */
export default function ResourceWarehouseModel({ position = [0, 0, 0], scale = 1, level = 1, animate = true }: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const dockRefs = useRef<Mesh[]>([])
  const conveyorLightRef = useRef<Mesh>(null)

  // Barrel vault roof approximation using half-cylinder
  const vaultGeo = useMemo(() => {
    const geo = new THREE.CylinderGeometry(3.0, 3.0, 9.0, 8, 1, false, 0, Math.PI)
    geo.rotateZ(Math.PI / 2)
    geo.rotateY(Math.PI / 2)
    return geo
  }, [])

  useFrame((state) => {
    if (!animate) return
    const t = state.clock.elapsedTime
    // Loading dock doors pulse staggered
    dockRefs.current.forEach((mesh, i) => {
      if (mesh) {
        const mat = mesh.material as MeshStandardMaterial
        mat.emissiveIntensity = 0.3 + Math.sin(t * 1.5 + i * 0.8) * 0.15
      }
    })
    // Conveyor marker light slides back and forth
    if (conveyorLightRef.current) {
      conveyorLightRef.current.position.x = Math.sin(t * (Math.PI * 2 / 4)) * 3.5
    }
  })

  const addDockRef = (el: Mesh | null, index: number) => {
    if (el) dockRefs.current[index] = el
  }

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Main structure */}
      <mesh position={[0, 1.25, 0]} castShadow>
        <boxGeometry args={[9.0, 2.5, 5.5]} />
        <meshStandardMaterial
          color="#2a3a2a"
          emissive="#22cc66"
          emissiveIntensity={0.08}
          metalness={0.7}
          roughness={0.3}
        />
      </mesh>

      {/* Barrel vault roof */}
      <mesh geometry={vaultGeo} position={[0, 2.5, 0]} castShadow>
        <meshStandardMaterial
          color="#3a4a3a"
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* Loading dock doors (front face) */}
      {[-2.5, 0, 2.5].map((x, i) => (
        <mesh
          key={`dock-${i}`}
          ref={(el) => addDockRef(el, i)}
          position={[x, 0.9, 2.76]}
        >
          <boxGeometry args={[1.5, 1.8, 0.1]} />
          <meshStandardMaterial
            color="#44aa66"
            emissive="#22cc66"
            emissiveIntensity={0.4}
            transparent
            opacity={0.5}
          />
        </mesh>
      ))}

      {/* Cargo stacks near dock */}
      <mesh position={[-3.0, 0.3, 3.5]} castShadow>
        <boxGeometry args={[0.8, 0.6, 0.6]} />
        <meshStandardMaterial color="#446633" metalness={0.5} roughness={0.4} />
      </mesh>
      <mesh position={[-2.2, 0.3, 3.8]} castShadow>
        <boxGeometry args={[0.6, 0.6, 0.6]} />
        <meshStandardMaterial color="#555544" metalness={0.5} roughness={0.4} />
      </mesh>
      <mesh position={[1.5, 0.3, 3.6]} castShadow>
        <boxGeometry args={[0.8, 0.6, 0.6]} />
        <meshStandardMaterial color="#336644" metalness={0.5} roughness={0.4} />
      </mesh>
      <mesh position={[3.2, 0.3, 3.5]} castShadow>
        <boxGeometry args={[0.6, 0.4, 0.4]} />
        <meshStandardMaterial color="#446633" metalness={0.5} roughness={0.4} />
      </mesh>
      <mesh position={[3.2, 0.6, 3.5]} castShadow>
        <boxGeometry args={[0.6, 0.4, 0.4]} />
        <meshStandardMaterial color="#555544" metalness={0.5} roughness={0.4} />
      </mesh>

      {/* Roof conveyor track */}
      <mesh position={[0, 4.6, 0]}>
        <boxGeometry args={[8.0, 0.15, 0.3]} />
        <meshStandardMaterial
          color="#4a5a4a"
          metalness={0.85}
          roughness={0.2}
        />
      </mesh>

      {/* Conveyor marker light */}
      <mesh ref={conveyorLightRef} position={[0, 4.8, 0]}>
        <sphereGeometry args={[0.08, 6, 6]} />
        <meshStandardMaterial
          color="#44ff88"
          emissive="#44ff88"
          emissiveIntensity={1.5}
          toneMapped={false}
        />
      </mesh>

      {/* Side vents - left */}
      {[-3.0, -1.0, 1.0, 3.0].map((x, i) => (
        <mesh key={`lv-${i}`} position={[x, 1.5, -2.78]}>
          <boxGeometry args={[0.05, 0.3, 1.0]} />
          <meshStandardMaterial
            color="#22cc66"
            emissive="#22cc66"
            emissiveIntensity={0.3}
          />
        </mesh>
      ))}

      {/* Side vents - right */}
      {[-3.0, -1.0, 1.0, 3.0].map((x, i) => (
        <mesh key={`rv-${i}`} position={[x, 1.5, 2.78]}>
          <boxGeometry args={[0.05, 0.3, 1.0]} />
          <meshStandardMaterial
            color="#22cc66"
            emissive="#22cc66"
            emissiveIntensity={0.3}
          />
        </mesh>
      ))}
    </group>
  )
}
