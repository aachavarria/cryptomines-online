import { useRef, useMemo } from 'react'
import { useFrame } from '@react-three/fiber'
import type { Mesh } from 'three'
import type { BuildingWithType } from '../../types/index.ts'

interface UpgradeProgressBar3DProps {
  building: BuildingWithType
  yOffset: number
}

export default function UpgradeProgressBar3D({ building, yOffset }: UpgradeProgressBar3DProps) {
  const fillRef = useRef<Mesh>(null)
  const barWidth = 2.5
  const barHeight = 0.2

  const { startTime, totalDuration } = useMemo(() => {
    const finish = building.upgrade_finish_at ? new Date(building.upgrade_finish_at).getTime() : 0
    const start = new Date(building.updated_at).getTime()
    return { startTime: start, totalDuration: finish - start }
  }, [building.upgrade_finish_at, building.updated_at])

  useFrame(() => {
    if (!fillRef.current || !building.upgrade_finish_at || totalDuration <= 0) return

    const elapsed = Date.now() - startTime
    const progress = Math.min(1, Math.max(0, elapsed / totalDuration))

    fillRef.current.scale.x = progress
    fillRef.current.position.x = -(barWidth * (1 - progress)) / 2
  })

  return (
    <group position={[0, yOffset, 0]} rotation={[0, -Math.PI / 4, 0]}>
      {/* Background bar */}
      <mesh>
        <planeGeometry args={[barWidth, barHeight]} />
        <meshBasicMaterial color="#333333" transparent opacity={0.7} depthWrite={false} />
      </mesh>
      {/* Fill bar */}
      <mesh ref={fillRef} position={[0, 0, 0.01]}>
        <planeGeometry args={[barWidth, barHeight]} />
        <meshBasicMaterial color="#ffaa22" transparent opacity={0.9} depthWrite={false} />
      </mesh>
    </group>
  )
}
