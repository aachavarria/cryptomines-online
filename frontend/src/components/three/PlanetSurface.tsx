import { useMemo } from 'react'
import * as THREE from 'three'
import { GRID_SIZE, TILE_WORLD_SIZE } from '../../contexts/GameContext.tsx'

export default function PlanetSurface() {
  // Create a subtle isometric grid texture
  const gridTexture = useMemo(() => {
    const size = 1024
    const canvas = document.createElement('canvas')
    canvas.width = size
    canvas.height = size
    const ctx = canvas.getContext('2d')!
    ctx.fillStyle = '#1a1a2e'
    ctx.fillRect(0, 0, size, size)

    // Draw subtle grid lines
    ctx.strokeStyle = 'rgba(42, 42, 78, 0.4)'
    ctx.lineWidth = 0.5
    const step = size / 32
    for (let i = 0; i <= 32; i++) {
      const pos = i * step
      ctx.beginPath()
      ctx.moveTo(pos, 0)
      ctx.lineTo(pos, size)
      ctx.stroke()
      ctx.beginPath()
      ctx.moveTo(0, pos)
      ctx.lineTo(size, pos)
      ctx.stroke()
    }

    const texture = new THREE.CanvasTexture(canvas)
    texture.wrapS = THREE.RepeatWrapping
    texture.wrapT = THREE.RepeatWrapping
    texture.repeat.set(8, 8)
    return texture
  }, [])

  // Calculate platform size to cover the square grid (centered at origin)
  // Grid spans from -40 to +40 on X/Z, diagonal ~57, add margin
  const platformRadius = GRID_SIZE * TILE_WORLD_SIZE * 1.0

  return (
    <group>
      {/* Main surface */}
      <mesh rotation={[-Math.PI / 2, 0, 0]} receiveShadow position={[0, -0.1, 0]}>
        <circleGeometry args={[platformRadius, 64]} />
        <meshStandardMaterial
          map={gridTexture}
          color="#1a1a2e"
          metalness={0.3}
          roughness={0.7}
        />
      </mesh>
      {/* Edge glow ring */}
      <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, -0.05, 0]}>
        <ringGeometry args={[platformRadius - 2, platformRadius, 64]} />
        <meshStandardMaterial
          color="#2244aa"
          emissive="#2244aa"
          emissiveIntensity={0.3}
          transparent
          opacity={0.4}
        />
      </mesh>
    </group>
  )
}
