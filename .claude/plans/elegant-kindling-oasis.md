# Three.js Frontend Performance Optimization Plan

## Context

El frontend tiene **~50 useFrame callbacks por frame** (22 building models + wrapper + utilities), Date calculations cada frame en progress bars, shadow maps sobredimensionados, y 22 building models cargados estáticamente. Con 20+ edificios en pantalla, esto impacta FPS especialmente en zoom out y dispositivos modestos.

**Objetivo**: Reducir callbacks por frame un ~70%, optimizar quick wins, y preparar code splitting + LOD para escalabilidad.

---

## Phase 1: Quick Wins (bajo riesgo, impacto inmediato)

### 1.1 Fix UpgradeProgressBar3D — cache Date calculations
**Archivo**: `frontend/src/components/three/UpgradeProgressBar3D.tsx`

- Añadir `useMemo` para cachear `finishTime`, `startTime`, `totalDuration` (depender de `building.upgrade_finish_at`, `building.updated_at`)
- El `useFrame` solo hace `Date.now()` y aritmética simple
- Elimina 3x `new Date().getTime()` por frame por edificio upgradeando

### 1.2 Shadow map 2048→1024
**Archivo**: `frontend/src/components/three/SceneLighting.tsx`

- Cambiar `shadow-mapSize-width={2048}` y `shadow-mapSize-height={2048}` → `1024`
- Reduce 4x la memoria de shadow map (16MB → 4MB GPU)
- No visible en orthographic a los zooms disponibles (2-30)

### 1.3 SelectionRing — guard useFrame cuando invisible
**Archivo**: `frontend/src/components/three/SelectionRing.tsx`

- Línea 13: ya tiene `if (ref.current && visible)` pero el useFrame sigue registrado
- Mover el `if (!visible) return null` ANTES del useFrame hook → NO, hooks no se pueden condicionar
- En su lugar: ya está bien optimizado (early return en useFrame + null render). No tocar.

### 1.4 Tooltip — ya es condicional
**Archivo**: `frontend/src/components/three/BuildingModel.tsx` línea 169

- `{isHovered && !isSelected && (<Html>...` — ya solo monta el Html cuando hover
- Solo 1 building puede estar hovered a la vez → solo 1 Html tooltip max
- **No necesita cambio**. El abbreviation `<Html>` (línea 126-142) sí está siempre montado pero es necesario.

---

## Phase 2: Animation Throttling (impacto alto, riesgo medio)

### Enfoque: `animate` prop + early return en building models

Cada building model recibe un nuevo prop `animate: boolean`. Cuando `false`, su `useFrame` hace early return. BuildingModel.tsx decide quién anima basándose en prioridad.

### 2.1 Añadir `animate` prop a BuildingComponentProps
**Archivo**: `frontend/src/components/three/buildings/index.ts`

```ts
export interface BuildingComponentProps {
  position?: [number, number, number]
  scale?: number
  level?: number
  animate?: boolean  // NEW
}
```

### 2.2 Actualizar los 22 building models
**Archivos**: Todos en `frontend/src/components/three/buildings/*.tsx`

Patrón para cada archivo (ejemplo CivicCenterModel.tsx):
```ts
// ANTES
export default function CivicCenterModel({ position = [0, 0, 0], scale = 1, level = 1 }: BuildingComponentProps) {
  useFrame((state) => {
    const t = state.clock.elapsedTime
    // ...animations
  })

// DESPUÉS
export default function CivicCenterModel({ position = [0, 0, 0], scale = 1, level = 1, animate = true }: BuildingComponentProps) {
  useFrame((state) => {
    if (!animate) return  // ← early return, callback sigue registrado pero es O(1)
    const t = state.clock.elapsedTime
    // ...animations sin cambios
  })
```

Los 22 archivos a modificar:
- CivicCenterModel, TechnologyCenterModel, AllianceCenterModel, TradingCenterModel
- GalaxyTransporterModel, CompoundCenterModel, MetalCollectorModel, He3ExtractorModel
- ResidentialAreaModel, ResourceWarehouseModel, ShipFactoryBuildingModel, SpacedockBuildingModel
- CommandCenterModel, WeaponResearchCenterModel, RadarModel, SpaceStationModel
- MeteorStarModel, ParticleCannonModel, AntiAircraftGunModel, ThorsCannonModel
- CelestialBaseModel, RecyclingPlantModel

### 2.3 BuildingModel.tsx — compute priority y pasar `animate`
**Archivo**: `frontend/src/components/three/BuildingModel.tsx`

```ts
// Determinar si este building merece animación completa
const shouldAnimate = isSelected || isHovered || building.is_upgrading

// Pasar a TypeModel
<TypeModel scale={levelScale} level={building.level} animate={shouldAnimate} />
```

**Resultado**: Solo el edificio selected/hovered/upgrading corre animaciones completas. Los demás quedan "congelados" en su último frame — visualmente imperceptible porque las animaciones son sutiles (pulsos, rotaciones lentas).

### 2.4 BuildingModel breathing — throttle no-priority
**Archivo**: `frontend/src/components/three/BuildingModel.tsx` líneas 61-67

```ts
useFrame((state) => {
  if (groupRef.current) {
    if (!isSelected && !isHovered) return  // solo breathing para selected/hovered
    const t = state.clock.elapsedTime
    groupRef.current.position.y = position[1] + Math.sin(t * 0.8 + seed * 6.28) * 0.05
  }
})
```

**Nota**: El breathing es tan sutil (0.05 units) que nadie nota si un edificio no-seleccionado no lo tiene.

---

## Phase 3: Code Splitting de Building Models

### 3.1 Lazy imports en buildings/index.ts
**Archivo**: `frontend/src/components/three/buildings/index.ts`

```ts
import { lazy, type ComponentType } from 'react'

const CivicCenterModel = lazy(() => import('./CivicCenterModel'))
const TechnologyCenterModel = lazy(() => import('./TechnologyCenterModel'))
// ... 22 lazy imports

export const BUILDING_MODELS: Record<string, ComponentType<BuildingComponentProps>> = {
  civic_center: CivicCenterModel,
  // ...
}
```

### 3.2 Suspense fallback en BuildingModel.tsx
**Archivo**: `frontend/src/components/three/BuildingModel.tsx`

Wrappear `<TypeModel>` en `<Suspense>` con fallback de box coloreado:

```tsx
import { Suspense } from 'react'

// Fallback: simple colored box matching building category
const BuildingFallback = ({ color, height }: { color: string; height: number }) => (
  <mesh position={[0, height / 2, 0]}>
    <boxGeometry args={[2, height, 2]} />
    <meshStandardMaterial color={color} emissive={color} emissiveIntensity={0.15} />
  </mesh>
)

// En el render:
{TypeModel ? (
  <Suspense fallback={<BuildingFallback color={color} height={baseHeight} />}>
    <TypeModel scale={levelScale} level={building.level} animate={shouldAnimate} />
  </Suspense>
) : (
  <BuildingFallback color={color} height={baseHeight} />
)}
```

**Resultado**: Bundle inicial reduce ~80KB gzipped. Cada building type se carga la primera vez que aparece.

---

## Phase 4: LOD System (zoom-based detail)

### 4.1 Crear SimpleBuildingLOD component
**Nuevo archivo**: `frontend/src/components/three/SimpleBuildingLOD.tsx`

Un componente genérico que renderiza una versión simplificada: 1 mesh con el color de categoría, scaled al tamaño del building. Para usar cuando zoom < threshold.

### 4.2 LOD en BuildingModel.tsx
**Archivo**: `frontend/src/components/three/BuildingModel.tsx`

```tsx
import { useThree } from '@react-three/fiber'

// Dentro del componente:
const { camera } = useThree()
const zoom = (camera as THREE.OrthographicCamera).zoom
const useLOD = zoom < 5  // threshold: zoom 5 ≈ edificios de ~8px

// En render:
{useLOD ? (
  <SimpleBuildingLOD color={color} size={size} />
) : (
  <Suspense fallback={<BuildingFallback ... />}>
    <TypeModel scale={levelScale} level={building.level} animate={shouldAnimate && !useLOD} />
  </Suspense>
)}
```

**Resultado**: A zoom máximo out (zoom=2), todos los buildings renderizan 1 mesh en vez de 15-30.

---

## Resumen de Impacto

| Fase | Esfuerzo | Impacto Performance | Archivos |
|------|----------|-------------------|----------|
| Phase 1: Quick Wins | 30 min | 10-15% (shadow + Date) | 2 archivos |
| Phase 2: Animation Throttle | 2-3 hrs | 40-50% (reduce work/frame) | 23 archivos (cambio mecánico) |
| Phase 3: Code Splitting | 1 hr | 15-20% initial load | 2 archivos |
| Phase 4: LOD System | 2 hrs | 20-30% at zoom out | 2 archivos nuevos |

## Orden de Implementación

```
Phase 1 → Phase 2 → Phase 3 → Phase 4
(cada fase es independiente y verificable)
```

## Verificación

Después de cada fase:
1. `cd frontend && npm run build` — sin errores de compilación
2. `npm run dev` — abrir en browser, verificar visualmente:
   - Edificios se ven igual (colores, animaciones de selected/hovered)
   - Placement mode funciona (grid, ghost, construct)
   - Context menu funciona (click building → menu aparece)
   - Zoom in/out funcional
3. Chrome DevTools → Performance tab → Record 5s:
   - Frame time < 16.67ms (60fps target)
   - Comparar useFrame callback count antes/después
4. Phase 4 específico: zoom out máximo → verificar LOD swap sin pop visual

## Archivos Clave

- `frontend/src/components/three/UpgradeProgressBar3D.tsx` — Phase 1
- `frontend/src/components/three/SceneLighting.tsx` — Phase 1
- `frontend/src/components/three/BuildingModel.tsx` — Phase 2, 3, 4
- `frontend/src/components/three/buildings/index.ts` — Phase 2, 3
- `frontend/src/components/three/buildings/*.tsx` (22 files) — Phase 2
- `frontend/src/components/three/SimpleBuildingLOD.tsx` — Phase 4 (nuevo)
- `frontend/src/components/three/PlanetScene.tsx` — referencia, no se modifica
