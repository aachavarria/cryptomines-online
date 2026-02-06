import { useRef } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Mesh, MeshStandardMaterial } from 'three'

interface SelectionRingProps {
  visible: boolean
  radius?: number
}

export default function SelectionRing({ visible, radius = 2.5 }: SelectionRingProps) {
  const ref = useRef<Mesh>(null)

  useFrame((state) => {
    if (ref.current && visible) {
      ref.current.rotation.y += 0.005
      const t = state.clock.elapsedTime
      const mat = ref.current.material as MeshStandardMaterial
      mat.opacity = 0.5 + Math.sin(t * 2) * 0.15
    }
  })

  if (!visible) return null

  return (
    <mesh ref={ref} rotation={[-Math.PI / 2, 0, 0]} position={[0, 0.1, 0]}>
      <torusGeometry args={[radius, 0.12, 8, 32]} />
      <meshStandardMaterial
        color="#4488ff"
        emissive="#4488ff"
        emissiveIntensity={0.8}
        transparent
        opacity={0.7}
        depthWrite={false}
      />
    </mesh>
  )
}
