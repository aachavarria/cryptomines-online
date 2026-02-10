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
 * Celestial Base: fortified defense outpost (2x2).
 * Star-shaped (5-pointed) base platform with central watchtower,
 * shield projector ring, armored walls with embrasures.
 */
export default function CelestialBaseModel({
  position = [0, 0, 0],
  scale = 1,
  level = 1,
  animate = true,
}: BuildingComponentProps) {
  const groupRef = useRef<Group>(null)
  const shieldRingRef = useRef<Mesh>(null)
  const embrasureRefs = useRef<(Mesh | null)[]>([])
  const windowRef = useRef<Mesh>(null)

  const color = '#44ccff'

  // 5-pointed star shape for extrusion
  const starGeometry = useMemo(() => {
    const shape = new THREE.Shape()
    const outerRadius = 3.5
    const innerRadius = 2.0
    const points = 5

    for (let i = 0; i < points * 2; i++) {
      const radius = i % 2 === 0 ? outerRadius : innerRadius
      const angle = (i * Math.PI) / points - Math.PI / 2
      const x = Math.cos(angle) * radius
      const y = Math.sin(angle) * radius
      if (i === 0) {
        shape.moveTo(x, y)
      } else {
        shape.lineTo(x, y)
      }
    }
    shape.closePath()

    const extrudeSettings = {
      depth: 1.0,
      bevelEnabled: false,
    }
    const geo = new THREE.ExtrudeGeometry(shape, extrudeSettings)
    geo.rotateX(-Math.PI / 2)
    geo.translate(0, 0.5, 0)
    return geo
  }, [])

  useFrame((state) => {
    if (!animate) return
    const t = state.clock.elapsedTime
    // Shield ring rotates
    if (shieldRingRef.current) {
      shieldRingRef.current.rotation.y = t * (Math.PI * 2 / 6)
      // Shield ring pulses
      const mat = shieldRingRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 0.5 + Math.sin(t * (Math.PI * 2 / 3)) * 0.25
    }
    // Embrasure lights blink in sequence around the star
    embrasureRefs.current.forEach((ref, i) => {
      if (ref) {
        const mat = ref.material as MeshStandardMaterial
        const phase = (t * 2 - i * 0.5) % 2.5
        mat.emissiveIntensity = phase < 0.3 ? 2.0 : 0.5
      }
    })
    // Tower observation windows glow with faint flicker
    if (windowRef.current) {
      const mat = windowRef.current.material as MeshStandardMaterial
      mat.emissiveIntensity = 0.3 + Math.sin(t * 7) * 0.1
    }
  })

  // Calculate star point positions for walls and lights
  const starPoints = useMemo(() => {
    const points: [number, number][] = []
    for (let i = 0; i < 5; i++) {
      const angle = (i * Math.PI * 2) / 5 - Math.PI / 2
      points.push([Math.cos(angle) * 3.5, Math.sin(angle) * 3.5])
    }
    return points
  }, [])

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* Star-shaped base platform */}
      <mesh geometry={starGeometry} castShadow>
        <meshStandardMaterial
          color="#1a2a3a"
          emissive={color}
          emissiveIntensity={0.1}
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* Central watchtower */}
      <mesh position={[0, 2.25, 0]} castShadow>
        <cylinderGeometry args={[0.8, 0.8, 3.5, 6]} />
        <meshStandardMaterial
          color="#2a3a4a"
          emissive={color}
          emissiveIntensity={0.12}
          metalness={0.8}
          roughness={0.25}
        />
      </mesh>

      {/* Tower roof - watchtower cone */}
      <mesh position={[0, 4.4, 0]} castShadow>
        <coneGeometry args={[0.9, 0.8, 6]} />
        <meshStandardMaterial
          color="#3a4a5a"
          metalness={0.85}
          roughness={0.2}
        />
      </mesh>

      {/* Observation windows on tower - 4 sides */}
      {[0, Math.PI / 2, Math.PI, Math.PI * 3 / 2].map((angle, i) => (
        <mesh
          key={`window-${i}`}
          ref={i === 0 ? windowRef : undefined}
          position={[Math.cos(angle) * 0.81, 3.2, Math.sin(angle) * 0.81]}
          rotation={[0, -angle, 0]}
        >
          <boxGeometry args={[0.4, 0.2, 0.02]} />
          <meshStandardMaterial
            color={color}
            emissive={color}
            emissiveIntensity={0.4}
            transparent
            opacity={0.5}
          />
        </mesh>
      ))}

      {/* Shield projector ring */}
      <mesh
        ref={shieldRingRef}
        position={[0, 2.5, 0]}
        rotation={[Math.PI / 2, 0, 0]}
      >
        <torusGeometry args={[1.5, 0.1, 6, 16]} />
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={0.8}
          transparent
          opacity={0.4}
        />
      </mesh>

      {/* 5 wall segments at each star point */}
      {starPoints.map(([px, pz], i) => {
        const angle = Math.atan2(pz, px)
        return (
          <group key={`wall-${i}`}>
            {/* Wall */}
            <mesh
              position={[px * 0.85, 1.25, pz * 0.85]}
              rotation={[0, -angle, 0]}
              castShadow
            >
              <boxGeometry args={[0.8, 1.5, 0.2]} />
              <meshStandardMaterial
                color="#2a3a4a"
                metalness={0.85}
                roughness={0.2}
              />
            </mesh>
            {/* Embrasure light at wall top */}
            <mesh
              ref={(el) => { embrasureRefs.current[i] = el }}
              position={[px * 0.85, 2.1, pz * 0.85]}
            >
              <sphereGeometry args={[0.06, 6, 6]} />
              <meshStandardMaterial
                color="#44ffaa"
                emissive="#44ffaa"
                emissiveIntensity={1.5}
                toneMapped={false}
              />
            </mesh>
          </group>
        )
      })}
    </group>
  )
}
