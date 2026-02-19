import { useMemo } from 'react'
import * as THREE from 'three'
import { useTexture } from '@react-three/drei'
import { GRID_SIZE, TILE_WORLD_SIZE } from '../../contexts/GameContext.tsx'

export default function PlanetSurface() {
  // Load the terrain texture
  const terrainTexture = useTexture('/assets/terrains/terrain.jpg')

  // Configure texture for seamless repeat
  useMemo(() => {
    terrainTexture.wrapS = THREE.RepeatWrapping
    terrainTexture.wrapT = THREE.RepeatWrapping
    // Adjust repeat to match the surface size (more repeats = smaller tiles)
    terrainTexture.repeat.set(16, 16)
  }, [terrainTexture])

  // Calculate platform size to cover the square grid (centered at origin)
  // Grid spans from -40 to +40 on X/Z, diagonal ~57, add margin
  const platformRadius = GRID_SIZE * TILE_WORLD_SIZE * 1.0

  // Grid area size (square)
  const gridSize = GRID_SIZE * TILE_WORLD_SIZE

  return (
    <group>
      {/* Main surface */}
      <mesh rotation={[-Math.PI / 2, 0, 0]} receiveShadow position={[0, -0.1, 0]}>
        <circleGeometry args={[platformRadius, 64]} />
        <meshStandardMaterial
          map={terrainTexture}
          color="#ffffff"
          metalness={0.1}
          roughness={0.9}
        />
      </mesh>

      {/* Lighter grid area overlay (SQUARE) */}
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

      {/* Edge glow ring */}
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
