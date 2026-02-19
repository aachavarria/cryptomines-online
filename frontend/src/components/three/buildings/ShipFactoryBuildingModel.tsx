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

function ShipFactoryGLB({
  position = [0, 0, 0],
  scale = 1,
  isUnderConstruction = false,
  hologramColor = 0x00ffff,
}: BuildingComponentProps) {
  const groupRef = useRef<Group>(null);
  const { scene } = useGLTF("/assets/gbl/buildings/ship-factory.glb");

  const clonedScene = useMemo(() => scene.clone(), [scene]);
  const originalMaterials = useRef<
    Map<THREE.Mesh, THREE.Material | THREE.Material[]>
  >(new Map());

  useEffect(() => {
    if (!clonedScene) return;

    clonedScene.traverse((child) => {
      if ((child as THREE.Mesh).isMesh) {
        const mesh = child as THREE.Mesh;

        if (isUnderConstruction) {
          if (!originalMaterials.current.has(mesh)) {
            originalMaterials.current.set(mesh, mesh.material);
          }

          mesh.material = new THREE.MeshBasicMaterial({
            color: hologramColor,
            transparent: true,
            opacity: hologramColor === 0xffffff ? 0.25 : 0.5,
            wireframe: false,
          });
        } else {
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
      <primitive object={clonedScene} scale={10} position={[0, 0, 0]} />
    </group>
  );
}

export default function ShipFactoryBuildingModel(props: BuildingComponentProps) {
  return (
    <Suspense fallback={null}>
      <ShipFactoryGLB {...props} />
    </Suspense>
  );
}

useGLTF.preload("/assets/gbl/buildings/ship-factory.glb");
