import { useRef, Suspense, useEffect, useMemo } from "react";
import { useGLTF } from "@react-three/drei";
import type { Group } from "three";
import * as THREE from "three";

interface BuildingComponentProps {
  position?: [number, number, number];
  scale?: number;
  level?: number;
  animate?: boolean;
  isUnderConstruction?: boolean;
  hologramColor?: number;
}

/**
 * Metal Collector: Loads the GLB model for metal extraction building (resource, 2x2).
 * Uses pre-made 3D model from assets/gbl/buildings/metal-extractor.glb
 * Applies hologram effect when under construction.
 */
function MetalCollectorGLB({
  position = [0, 0, 0],
  scale = 1,
  isUnderConstruction = false,
  hologramColor = 0x00ffff,
}: BuildingComponentProps) {
  const groupRef = useRef<Group>(null);
  const { scene } = useGLTF("/assets/gbl/buildings/metal-extractor.glb");

  // Clone the scene to avoid modifying the cached original
  const clonedScene = useMemo(() => scene.clone(), [scene]);

  // Store original materials to restore them later
  const originalMaterials = useRef<
    Map<THREE.Mesh, THREE.Material | THREE.Material[]>
  >(new Map());

  useEffect(() => {
    if (!clonedScene) return;

    clonedScene.traverse((child) => {
      if ((child as THREE.Mesh).isMesh) {
        const mesh = child as THREE.Mesh;

        if (isUnderConstruction) {
          // Save original material if not already saved
          if (!originalMaterials.current.has(mesh)) {
            originalMaterials.current.set(mesh, mesh.material);
          }

          // Apply hologram material
          mesh.material = new THREE.MeshBasicMaterial({
            color: hologramColor,
            transparent: true,
            opacity: hologramColor === 0xffffff ? 0.25 : 0.5,
            wireframe: false,
          });
        } else {
          // Restore original material
          const original = originalMaterials.current.get(mesh);
          if (original) {
            mesh.material = original;
          }
        }
      }
    });
  }, [clonedScene, isUnderConstruction, hologramColor]);

  return (
    <group ref={groupRef} position={position} scale={scale}>
      {/* GLB models from Tripo3D need scaling adjustment */}
      {/* Offset Y to compensate for center pivot (move up so base is at ground) */}
      <primitive object={clonedScene} scale={10} position={[0, 0, 0]} />
    </group>
  );
}

export default function MetalCollectorModel(props: BuildingComponentProps) {
  return (
    <Suspense fallback={null}>
      <MetalCollectorGLB {...props} />
    </Suspense>
  );
}

// Preload the model for better performance
useGLTF.preload("/assets/gbl/buildings/metal-extractor.glb");
