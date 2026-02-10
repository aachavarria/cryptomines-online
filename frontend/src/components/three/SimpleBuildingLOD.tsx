import { TILE_WORLD_SIZE } from '../../contexts/GameContext.tsx'
import { getBuildingSize } from '../../config/buildingConfig.ts'

interface SimpleBuildingLODProps {
  typeName: string
  color: string
  levelScale: number
}

export default function SimpleBuildingLOD({ typeName, color, levelScale }: SimpleBuildingLODProps) {
  const size = getBuildingSize(typeName)
  const width = size.cols * TILE_WORLD_SIZE * 0.6
  const depth = size.rows * TILE_WORLD_SIZE * 0.6
  const height = (typeName === 'civic_center' ? 5 : 3) * levelScale

  return (
    <mesh position={[0, height / 2, 0]} castShadow>
      <boxGeometry args={[width, height, depth]} />
      <meshStandardMaterial
        color={color}
        emissive={color}
        emissiveIntensity={0.15}
        metalness={0.7}
        roughness={0.3}
      />
    </mesh>
  )
}
