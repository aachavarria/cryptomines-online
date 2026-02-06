import { useRef, useEffect } from 'react'
import { useThree, useFrame } from '@react-three/fiber'
import { MapControls } from '@react-three/drei'
import * as THREE from 'three'
import type { MapControls as MapControlsImpl } from 'three-stdlib'

interface CameraControllerProps {
  target: [number, number, number] | null
  enabled: boolean
}

// Isometric direction: equal X/Y/Z gives 45° rotation with orthographic
const CAMERA_DIR = new THREE.Vector3(1, 1, 1).normalize()
const CAMERA_DISTANCE = 100
const CAMERA_OFFSET = CAMERA_DIR.clone().multiplyScalar(CAMERA_DISTANCE)

export default function CameraController({ target, enabled }: CameraControllerProps) {
  const controlsRef = useRef<MapControlsImpl>(null)
  const { camera } = useThree()
  const isAnimating = useRef(false)
  const animStart = useRef(0)
  const startTarget = useRef(new THREE.Vector3())
  const endTarget = useRef(new THREE.Vector3())
  const startZoom = useRef(10)
  const endZoom = useRef(10)

  useEffect(() => {
    // Set initial camera position (isometric angle via orthographic)
    camera.position.copy(CAMERA_OFFSET)
    camera.lookAt(0, 0, 0)
  }, [camera])

  useEffect(() => {
    if (target && controlsRef.current) {
      isAnimating.current = true
      animStart.current = performance.now()
      startTarget.current.copy(controlsRef.current.target)
      endTarget.current.set(target[0], target[1], target[2])
    }
  }, [target, camera])

  useFrame(() => {
    if (isAnimating.current && controlsRef.current) {
      const elapsed = performance.now() - animStart.current
      const duration = 500
      let t = Math.min(elapsed / duration, 1)
      // Ease in-out
      t = t < 0.5 ? 2 * t * t : 1 - Math.pow(-2 * t + 2, 2) / 2

      // Lerp the controls target
      controlsRef.current.target.lerpVectors(startTarget.current, endTarget.current, t)
      // Keep camera at fixed offset from target
      camera.position.copy(controlsRef.current.target).add(CAMERA_OFFSET)
      controlsRef.current.update()

      if (t >= 1) {
        isAnimating.current = false
      }
    }
  })

  return (
    <MapControls
      ref={controlsRef}
      enabled={enabled && !isAnimating.current}
      enableRotate={false}
      enableDamping
      dampingFactor={0.08}
      panSpeed={1.0}
      zoomSpeed={0.5}
      minZoom={2}
      maxZoom={30}
      screenSpacePanning={true}
    />
  )
}
