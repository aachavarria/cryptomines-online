import { useMemo } from 'react'
import * as THREE from 'three'
import { useTexture } from '@react-three/drei'
import { GRID_SIZE, TILE_WORLD_SIZE, type BaseView } from '../../contexts/GameContext.tsx'

interface PlanetSurfaceProps {
  base: BaseView
}

export default function PlanetSurface({ base }: PlanetSurfaceProps) {
  if (base === 'space') {
    return <SpaceBaseSurface />
  }
  return <GroundBaseSurface />
}

function GroundBaseSurface() {
  const terrainTexture = useTexture('/assets/terrains/terrain.jpg')

  useMemo(() => {
    terrainTexture.wrapS = THREE.RepeatWrapping
    terrainTexture.wrapT = THREE.RepeatWrapping
    terrainTexture.repeat.set(16, 16)
  }, [terrainTexture])

  const platformRadius = GRID_SIZE * TILE_WORLD_SIZE * 1.0
  const gridSize = GRID_SIZE * TILE_WORLD_SIZE

  return (
    <group>
      <mesh rotation={[-Math.PI / 2, 0, 0]} receiveShadow position={[0, -0.1, 0]}>
        <circleGeometry args={[platformRadius, 64]} />
        <meshStandardMaterial
          map={terrainTexture}
          color="#ffffff"
          metalness={0.1}
          roughness={0.9}
        />
      </mesh>

      <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, -0.09, 0]}>
        <planeGeometry args={[gridSize, gridSize]} />
        <meshStandardMaterial
          map={terrainTexture}
          color="#ffffff"
          emissiveIntensity={0.15}
          emissive="#ffffff"
          metalness={0.1}
          roughness={0.9}
          transparent
          opacity={0.15}
        />
      </mesh>

      <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, -0.05, 0]}>
        <ringGeometry args={[platformRadius - 2, platformRadius, 64]} />
        <meshStandardMaterial
          color="#aa6633"
          emissive="#aa6633"
          emissiveIntensity={0.2}
          transparent
          opacity={0.3}
        />
      </mesh>
    </group>
  )
}

function SpaceBaseSurface() {
  const platformRadius = GRID_SIZE * TILE_WORLD_SIZE * 1.0
  const gridSize = GRID_SIZE * TILE_WORLD_SIZE

  return (
    <group>
      {/* Dark orbital platform — translucent steel-blue */}
      <mesh rotation={[-Math.PI / 2, 0, 0]} receiveShadow position={[0, -0.1, 0]}>
        <circleGeometry args={[platformRadius, 64]} />
        <meshStandardMaterial
          color="#0a1430"
          emissive="#1a2c5c"
          emissiveIntensity={0.25}
          metalness={0.7}
          roughness={0.3}
          transparent
          opacity={0.85}
        />
      </mesh>

      {/* Energy grid overlay */}
      <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, -0.09, 0]}>
        <planeGeometry args={[gridSize, gridSize]} />
        <meshStandardMaterial
          color="#3388ff"
          emissive="#3388ff"
          emissiveIntensity={0.3}
          metalness={0.5}
          roughness={0.4}
          transparent
          opacity={0.18}
        />
      </mesh>

      {/* Bright cyan edge ring (orbital perimeter marker) */}
      <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, -0.05, 0]}>
        <ringGeometry args={[platformRadius - 2, platformRadius, 64]} />
        <meshStandardMaterial
          color="#44ccff"
          emissive="#44ccff"
          emissiveIntensity={0.6}
          transparent
          opacity={0.55}
        />
      </mesh>

      {/* Distant planet "below" the platform — gives spatial context */}
      <mesh position={[0, -120, 0]}>
        <sphereGeometry args={[60, 48, 48]} />
        <meshStandardMaterial
          color="#2a4a6c"
          emissive="#1a3a5c"
          emissiveIntensity={0.2}
          roughness={1.0}
        />
      </mesh>
    </group>
  )
}
