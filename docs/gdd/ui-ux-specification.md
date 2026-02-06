# Cryptomines Online - UI/UX Specification

> **Version**: 2.2
> **Last Updated**: 2026-02-06
> **Status**: Phase 1 UI Redesign (Grid v2 — multi-tile buildings)
> **Scope**: Planet/Base View (3D), Building Detail, Construction, Resource HUD
> **Renderer**: Three.js (React Three Fiber) for 3D scenes + HTML/CSS overlays for HUD/panels
> **Directive**: Visual style inspired by Galaxy Online 2 layout and structure. 3D planet view with HTML overlay panels. All mechanics names faithful to GO2.

---

## Table of Contents

1. [Design Philosophy](#1-design-philosophy)
2. [Rendering Architecture: 3D vs 2D](#2-rendering-architecture-3d-vs-2d)
3. [Color Palette & Typography](#3-color-palette--typography)
4. [Layout Architecture](#4-layout-architecture)
5. [Screen: Planet/Base View (3D)](#5-screen-planetbase-view-3d)
6. [Screen: Building Detail Panel (HTML Overlay)](#6-screen-building-detail-panel-html-overlay)
7. [Screen: Construction/Upgrade Panel (HTML Overlay)](#7-screen-constructionupgrade-panel-html-overlay)
8. [Component: Resource HUD (HTML Overlay)](#8-component-resource-hud-html-overlay)
9. [Reusable Components](#9-reusable-components)
10. [3D Assets Required](#10-3d-assets-required)
11. [2D Assets Required](#11-2d-assets-required)
12. [Camera System](#12-camera-system)
13. [3D Visual Effects & Lighting](#13-3d-visual-effects--lighting)
14. [Responsive Design](#14-responsive-design)
15. [API Integration Map](#15-api-integration-map)
16. [Performance Considerations](#16-performance-considerations)

---

## 1. Design Philosophy

### 1.1 Core Principles

- **It must look like a GAME, not a web app.** No admin panels, no CRUD lists.
- The planet base is a **3D scene** rendered with Three.js. Buildings are 3D models placed on a planet surface.
- HUD elements (resources, menus, panels) are **HTML/CSS overlays** rendered on top of the Three.js canvas, exactly like modern web-based 3D games.
- GO2 uses a dark sci-fi theme with metallic panels, glowing borders, and iconography. We replicate that feeling.
- Information density is high but organized -- GO2 shows resources, timers, and status at a glance.
- Every interactive element has visual feedback (hover glow, click effect, state change).

### 1.2 GO2 Visual Reference

GO2's planet view shows:
- **Top bar**: Player info (level, name) + resources (Metal/He3/Gold with icons and rates) always visible
- **Left sidebar**: Navigation icons (Base, Fleet, Tech, Commander, Corp, etc.)
- **Center area**: Visual planet surface with building sprites on an **isometric grid** -- buildings are freely placed by the player on a 20x20 tile grid with multi-tile footprints (not fixed positions/slots)
- **Context menu on click**: Left-clicking a building opens a context menu with **View / Move / Upgrade** buttons (NOT a direct panel opening). View opens the detail/management panel.
- **Bottom-right area**: Build button, Military button, construction info panel with timers
- **Hover tooltips**: Hovering over a building shows "Lv: X [Building Name]"

For Phase 1, we focus on: Top bar + Left sidebar + Center 3D planet view + Context menu + Detail panel (via View) + Construction info.

### 1.3 3D vs 2D Split

The fundamental rule: **3D renders the world, 2D renders the interface.**

| Layer | Renderer | Content |
|-------|----------|---------|
| Background | Three.js | Starfield skybox, ambient space |
| Planet Surface | Three.js | Planet terrain with isometric tile grid, building 3D models, effects |
| Hover Tooltips | Three.js (HTML in 3D via CSS2DRenderer) | "Lv: X [Name]" tooltip on building hover |
| Placement Grid | Three.js | Green/red tile indicators during construction/move mode |
| HUD Overlay | HTML/CSS | Resource bar, side nav, construction info panel |
| Context Menu | HTML/CSS | View/Move/Upgrade popup on building click |
| Panels | HTML/CSS | Building detail panel (opened via View), construction modal |
| Modals | HTML/CSS | Construction selection, confirmations |

---

## 2. Rendering Architecture: 3D vs 2D

### 2.1 Tech Stack

```
Three.js Layer:
  - three (core)
  - @react-three/fiber (React bindings for Three.js)
  - @react-three/drei (helpers: MapControls/custom pan, Html, Environment, etc.)
  - @react-three/postprocessing (optional: bloom, glow effects)

HTML Overlay Layer:
  - Standard React components
  - CSS with custom properties
  - Positioned absolutely over the Three.js canvas
  - pointer-events: none on overlay container, pointer-events: auto on interactive elements
```

### 2.2 Layer Stack (z-index)

```
z-index structure (bottom to top):

z: 0     Three.js Canvas (position: fixed, covers viewport)
z: 10    CSS2DRenderer hover tooltips (Three.js managed, "Lv: X [Name]")
z: 50    SideNav overlay (HTML, left side)
z: 60    Construction Info Panel (HTML, bottom-right, countdown timers)
z: 100   Resource HUD (HTML, fixed top bar)
z: 150   Building Context Menu (HTML, floating near clicked building)
z: 200   Building Detail Panel (HTML, right slide-in, opened via "View")
z: 300   Construction Modal (HTML, centered overlay)
z: 400   Tooltips (HTML, floating)
z: 500   Toast notifications (HTML, top-right)
```

### 2.3 Component Architecture

```
<App>
  <GameProvider>                          // React Context for game state
    <GameShell>                           // Full viewport container

      {/* 3D LAYER - Full viewport canvas */}
      <Canvas>                            // @react-three/fiber
        <PlanetScene>                     // The 3D planet scene
          <SceneLighting />               // Lights setup
          <Skybox />                      // Starfield background
          <PlanetSurface />              // Planet terrain mesh (scaled for 20x20 grid)
          <IsometricGrid />              // 20x20 wireframe diamond grid (visible in placement mode)
          <GhostPreview />               // Translucent building preview following cursor
          <BuildingModels />             // 3D building instances on multi-tile footprints
          <HoverTooltip />               // "Lv: X [Name]" tooltip on hover
          <SelectionIndicator />         // Glow ring on selected building
          <ConstructionEffects />        // Yellow progress bar + particles
          <PlacementIndicator />         // Green/red tile feedback for build/move
          <CameraController />           // Pan-based camera (click-drag to move)
        </PlanetScene>
      </Canvas>

      {/* 2D HTML OVERLAY LAYER - over the canvas */}
      <div className="hud-overlay">
        <ResourceHUD />                   // Fixed top bar
        <SideNav />                       // Left icon navigation
        <BuildingContextMenu />          // View/Move/Upgrade popup on building click
        <ConstructionInfoPanel />        // Bottom-right panel with countdown timers
        <BuildingDetailPanel />          // Right slide-in (opened via View in context menu)
        <ConstructionPanel />            // Modal (when constructing new via Build button)
        <ToastContainer />               // Notification toasts
      </div>

    </GameShell>
  </GameProvider>
</App>
```

### 2.4 Interaction Flow Between 3D and 2D

**Building Click -> Context Menu Flow** (faithful to GO2):
```
User left-clicks 3D building model
  |
  +--> Three.js Raycaster detects intersection
  |
  +--> Identify building from mesh userData
  |
  +--> 3D scene reacts: SelectionIndicator glows around building
  |
  +--> HTML overlay: BuildingContextMenu appears near the building
  |      Shows 3 buttons stacked vertically: [View] [Move] [Upgrade]
  |      (Warehouse shows 4 buttons: [View] [Move] [Upgrade] [Harvest])
  |
  +--> User clicks "View":
  |      +--> BuildingDetailPanel slides in from right
  |      +--> Shows building stats, production, upgrade costs
  |
  +--> User clicks "Move":
  |      +--> Enter move mode: grid tiles become visible
  |      +--> Green tiles = valid, Red tiles = invalid
  |      +--> Click a green tile to reposition building
  |
  +--> User clicks "Upgrade":
  |      +--> API call: POST /api/planets/:id/buildings/:bid/upgrade
  |      +--> Yellow progress bar appears over building in 3D
  |      +--> ConstructionInfoPanel (bottom-right) shows countdown
  |
  +--> User clicks "Harvest" (Warehouse only):
         +--> API call: POST /api/planets/:id/resources/collect
         +--> Resources update in HUD
```

**Hover Tooltip Flow**:
```
User hovers mouse over 3D building model
  |
  +--> Three.js Raycaster detects hover
  |
  +--> HoverTooltip appears: dark tooltip showing "Lv: X [Building Name]"
  |    Example: "Lv: 5 Metal Collector"
  |
  +--> Mouse leaves building: tooltip disappears
```

**Construction Flow**:
```
User clicks "Build" button (lower-right of screen)
  |
  +--> ConstructionPanel modal opens (building selection)
  |
  +--> User selects a building type and clicks "Build"
  |
  +--> Modal closes; placement mode activates:
  |      Wireframe grid (20x20) fades in across surface
  |      Ghost preview (translucent building model) appears at cursor
  |
  +--> Ghost follows cursor, snapping to grid positions
  |      Green ghost + green tiles = valid multi-tile footprint
  |      Red ghost + red tiles = invalid (overlap or out of bounds)
  |
  +--> User clicks a valid (green) position
  |
  +--> API call: POST /api/planets/:id/buildings (with grid col, row)
  |
  +--> New building model spawns at grid position with build-up animation
  |
  +--> Yellow progress bar over building + ConstructionInfoPanel shows countdown
  |
  +--> Wireframe grid fades out, ghost disappears
```

---

## 3. Color Palette & Typography

### 3.1 Color Palette

```
PRIMARY BACKGROUND
  --bg-deep:        #050510     Dark space void (app background, not visible behind 3D)
  --bg-panel:       #0c0c20cc   Panel backgrounds (semi-transparent to show 3D behind)
  --bg-panel-solid: #0c0c20     Panel backgrounds (opaque variant for detail panels)
  --bg-panel-light: #141430cc   Lighter panel variant (hover states)
  --bg-header:      #0a0a1eee   Top bar background (slightly transparent)

BORDERS & FRAMES
  --border-frame:   #1a1a40     Default panel borders
  --border-active:  #2244aa     Active/selected borders
  --border-glow:    #3366cc     Glowing border for important elements

ACCENTS
  --accent-primary: #4488ff     Primary interactive elements (buttons, links)
  --accent-hover:   #66aaff     Hover state
  --accent-success: #22cc66     Success states, "Collect" button
  --accent-warning: #ffaa22     Timers, in-progress states
  --accent-danger:  #ff4444     Error states, insufficient resources
  --accent-gold:    #ffd700     Gold resource, premium, rare items

RESOURCE COLORS (consistent across 2D UI and 3D labels)
  --color-metal:    #8899aa     Metal - steel/silver tone
  --color-he3:      #44ccff     He3 - cyan/energy blue
  --color-gold:     #ffd700     Gold - classic gold

TEXT
  --text-bright:    #e8e8f8     Primary text
  --text-normal:    #b0b0cc     Secondary text
  --text-dim:       #666688     Tertiary/disabled text
  --text-label:     #8888aa     Labels and categories

3D SCENE COLORS (Three.js materials)
  --3d-ground:      #1a1a2e     Planet surface base color
  --3d-ground-grid: #2a2a4e     Grid lines on surface
  --3d-glow-idle:   #334488     Idle building base glow
  --3d-glow-select: #4488ff     Selected building glow
  --3d-glow-upgrade:#ffaa22     Upgrading building glow
  --3d-ambient:     #1a1a3a     Ambient light color
  --3d-key-light:   #aabbff     Key directional light color
  --3d-rim-light:   #4466aa     Rim/back light color

GLOW EFFECTS (CSS)
  --glow-blue:      0 0 12px rgba(68, 136, 255, 0.4)
  --glow-cyan:      0 0 12px rgba(68, 204, 255, 0.4)
  --glow-gold:      0 0 12px rgba(255, 215, 0, 0.3)
  --glow-green:     0 0 12px rgba(34, 204, 102, 0.4)
```

### 3.2 Typography

```
FONT STACK
  Primary:     'Rajdhani', 'Orbitron', 'Segoe UI', system-ui, sans-serif
  Monospace:   'JetBrains Mono', 'Fira Code', monospace  (for numbers/timers)

SIZES
  --font-xs:     0.65rem    (labels, category tags)
  --font-sm:     0.75rem    (secondary info, rates)
  --font-base:   0.875rem   (body text, building names)
  --font-md:     1rem       (section headers, resource values)
  --font-lg:     1.25rem    (page titles)
  --font-xl:     1.5rem     (planet name)
  --font-timer:  0.85rem    (countdown timers - monospace)

WEIGHTS
  Normal: 400    (body text)
  Medium: 500    (labels)
  Semi:   600    (buttons, building names)
  Bold:   700    (resource numbers, timers, headers)
```

### 3.3 Spacing System

```
--space-xs:   0.25rem   (4px)
--space-sm:   0.5rem    (8px)
--space-md:   0.75rem   (12px)
--space-lg:   1rem      (16px)
--space-xl:   1.5rem    (24px)
--space-xxl:  2rem      (32px)

PANEL PADDING:   var(--space-md) var(--space-lg)
CARD PADDING:    var(--space-sm) var(--space-md)
```

---

## 4. Layout Architecture

### 4.1 Overall App Shell

The entire viewport is the Three.js canvas. HTML elements float on top.

```
+------------------------------------------------------------------+
|  [RESOURCE HUD - HTML overlay, fixed top]          z:100          |
|  Logo | Metal: XXX | He3: XXX | Gold: XXX | Storage | [Collect]  |
+------+-----------------------------------------------------------+
|      |                                                            |
| NAV  |     THREE.JS CANVAS (full viewport)                       |
| BAR  |                                                            |
| HTML |     3D Planet surface with isometric grid                  |
| z:50 |     Pan camera (drag to move viewport)                    |
|      |                                                            |
|  []  |     [3D buildings placed on isometric tiles]              |
|  []  |     [Hover tooltip: "Lv: X [Name]" on mouse over]        |
|  []  |     [Context menu on click: View/Move/Upgrade]            |
|  []  |     [Selection glow ring, particles]                      |
|  []  |                                          +----------------+
|      |                                          | Construction   |
+------+------------------------------------------+ Info Panel     |
|  [BUILD]  [MILITARY]                            | Timer 1: 12:45 |
|                     (bottom-right buttons)      | Timer 2: 00:32 |
+------------------------------------------------------------------+
```

When "View" is clicked in the context menu, the detail panel slides in from right:

```
+------+-------------------------------+-----------------------+
|      |                               |  BUILDING DETAIL      |
| NAV  |  3D CANVAS continues          |  HTML overlay z:200   |
|      |  behind the panel              |  360px width          |
|      |  (camera can still pan)       |  Semi-transparent bg  |
|      |                               |                       |
+------+-------------------------------+-----------------------+
```

### 4.2 Navigation Structure

**SideNav** - HTML overlay, vertical icon bar on left side (60px wide):

| Icon | Label | Route/Action | Phase | 3D/2D |
|------|-------|-------------|-------|-------|
| Planet icon | Base | Show planet 3D scene | Phase 1 | 3D scene |
| Flask icon | Research | `/research` (2D page) | Phase 2 | 2D page |
| Ship icon | Fleet | `/fleet` (2D page) | Phase 2 | 2D page |
| Person icon | Commander | `/commanders` (2D page) | Phase 2 | 2D page |
| Map icon | Galaxy | Galaxy map (3D scene) | Phase 3 | 3D scene |
| Shield icon | Corp | `/corp` (2D page) | Phase 3 | 2D page |

For Phase 1, only "Base" is active. Others show as dimmed/locked with tooltip "Coming Soon".

---

## 5. Screen: Planet/Base View (3D)

This is the PRIMARY screen for Phase 1. The planet base is a 3D scene with building models placed on a terrain surface.

### 5.1 Scene Composition

```
THREE.JS SCENE GRAPH:

Scene
  |
  +-- Skybox (starfield cube texture or procedural stars)
  |
  +-- AmbientLight (low intensity, --3d-ambient color)
  |
  +-- DirectionalLight (key light, --3d-key-light, casts shadows)
  |
  +-- DirectionalLight (rim light, --3d-rim-light, opposite side)
  |
  +-- PlanetSurface (ground plane scaled for 20x20 grid)
  |     |
  |     +-- IsometricGrid (20x20 wireframe diamond tiles, visible during placement mode)
  |     +-- PlacementOverlay (green/red multi-tile footprint overlays)
  |     +-- GhostPreview (translucent building model following cursor during placement)
  |
  +-- BuildingGroup
  |     |
  |     +-- Building_0 (3D model placed on a grid tile)
  |     |     +-- SelectionRing (visible when selected)
  |     |     +-- UpgradeProgressBar (yellow bar, visible when upgrading)
  |     |     +-- UpgradeParticles (visible when upgrading)
  |     |
  |     +-- Building_1 ...
  |     +-- Building_N ...
  |
  +-- HoverTooltipAnchor (CSS2DObject, follows hovered building)
  |
  +-- PanCamera (drag-based camera panning, NOT orbit)
```

### 5.2 Planet Surface

**Flat Base with Isometric Grid** (faithful to GO2):

A large flat platform representing the planet surface, viewed from a fixed isometric angle (top-down tilted). Buildings are placed on a **20x20 isometric grid** where each building occupies a multi-tile footprint depending on its type.

- Geometry: `PlaneGeometry`, radius scaled to accommodate the full 20x20 grid
- Material: `MeshStandardMaterial` with dark rocky texture
- Grid: **20x20 isometric diamond tiles** -- hidden by default, wireframe lines appear during placement mode
- Edge: Fading to transparent at the edges, merging into space
- The ground has a slight emissive glow along the grid lines when visible

**Grid Visibility States**:
- **Default**: Grid lines are fully hidden. Buildings sit on the surface naturally with no visible grid.
- **Placement mode** (build or move): A **subtle cyan wireframe** grid becomes visible across the entire 20x20 surface. Tiles under the cursor are highlighted. Color-coded feedback:
  - **Green tiles** = valid placement (all tiles in the building footprint are unoccupied)
  - **Red tiles** = invalid placement (one or more tiles in the footprint are occupied or out of bounds)
- Grid appearance: diamond-shaped (isometric) tiles, thin cyan wireframe lines (`#00cccc`, low opacity ~0.15), brightening under the cursor (~0.4)

### 5.3 Building Placement (Multi-Tile Isometric Grid)

Buildings are **freely placed by the player** on a **20x20 isometric grid** of diamond tiles, NOT on predefined fixed slots. This is faithful to GO2's grid system where players choose where to place each building. Each building type occupies a **multi-tile footprint** (cols x rows).

```
ISOMETRIC GRID VIEW (20x20 diamond tiles):

        /\    /\    /\    /\    /\    ...
       /  \  /  \  /  \  /  \  /  \
      / 0,0\/ 1,0\/ 2,0\/ 3,0\/ 4,0\  ...
      \    /\    /\    /\    /\    /
       \  /  \  /  \  /  \  /  \  /
        \/ 0,1\/ 1,1\/ 2,1\/ 3,1\/   ...
        /\    /\    /\    /\    /\
       /  \  /  \  /  \  /  \  /  \
      / 0,2\/ 1,2\/ 2,2\/ 3,2\/ 4,2\  ...
      \    /\    /\    /\    /\    /
       ...

Grid: 20 columns x 20 rows = 400 tiles.
Tile spacing: ~2 world units between tile centers.
Coordinate: (col, row), 0-indexed, col=0..19, row=0..19.
```

#### 5.3.1 Building Footprints (cols x rows)

Each building type occupies a specific multi-tile footprint on the grid:

| Footprint | Buildings |
|-----------|-----------|
| **3x3** (9 tiles) | Civic Center, Space Station |
| **3x2** (6 tiles) | Ship Factory, Resource Warehouse, Spacedock, Technology Center |
| **2x2** (4 tiles) | Metal Collector, He3 Extractor, Residential Area, Alliance Center, Trading Center, Galaxy Transporter, Compound Center, Command Center, Weapon Research Center, Recycling Plant, Thor's Cannon, Celestial Base |
| **1x2** (2 tiles) | Particle Cannon |
| **1x1** (1 tile) | Radar, Meteor Star, Anti-Aircraft Gun |

A building placed at grid position (col, row) occupies all tiles from (col, row) to (col + width - 1, row + height - 1). For example, a 3x3 Civic Center at (8, 8) occupies the 9 tiles: (8,8), (9,8), (10,8), (8,9), (9,9), (10,9), (8,10), (9,10), (10,10).

#### 5.3.2 Placement Rules

- Buildings occupy their full multi-tile footprint -- all tiles must be empty and within bounds
- **Civic Center** starts pre-placed near the center of the grid (always present)
- All other buildings are placed by the player when constructed
- Buildings can be **moved** to any valid position via the Move option in the context menu
- The 20x20 grid provides 400 tiles total, accommodating all buildings with room to arrange

#### 5.3.3 Ghost Preview Placement System

During placement mode, a **translucent ghost preview** of the building follows the player's cursor across the grid, providing immediate visual feedback:

- **Ghost model**: A semi-transparent copy of the building's 3D model (opacity ~0.5)
- **Snaps to grid**: The ghost locks to the nearest valid grid position as the cursor moves
- **Color feedback**:
  - **Green tint** + green highlighted tiles = valid placement (all footprint tiles empty and in bounds)
  - **Red tint** + red highlighted tiles = invalid placement (overlap with existing building or out of bounds)
- **Cursor highlight**: The tiles covered by the building footprint are highlighted with colored overlays matching the valid/invalid state
- The ghost preview is visible only during placement mode (build or move)
- Press ESC to cancel placement and hide the ghost

#### 5.3.4 Wireframe Grid

When placement mode activates, a **subtle wireframe overlay** appears across the entire 20x20 grid surface:

- **Lines**: Thin cyan wireframe (`#00cccc`) drawn along the isometric diamond tile edges
- **Base opacity**: Very low (~0.15) for ambient visibility
- **Cursor highlight**: Tiles directly under the cursor (matching the building footprint) brighten to ~0.4 opacity
- **Visibility**: Only shown during placement mode; fully hidden in normal gameplay
- **Fade animation**: Grid fades in over 200ms when entering placement mode, fades out over 200ms when exiting

#### 5.3.5 New Building Placement Flow

1. Player clicks "Build" button (bottom-right of screen)
2. Selects building type from ConstructionPanel modal and clicks "Build"
3. Modal closes; **placement mode activates**:
   - Wireframe grid fades in across the 20x20 surface
   - Ghost preview of the selected building type appears at the cursor
4. As the cursor moves, the ghost snaps to grid positions with green/red feedback
5. Player clicks on a valid (green) position to confirm placement
6. API call fires: `POST /api/planets/:id/buildings` with grid coordinates
7. Building appears at that position with a build-up animation (scale 0 -> 1)
8. Wireframe grid fades out, ghost disappears
9. If ESC is pressed during placement, construction is cancelled

#### 5.3.6 Building Move Flow

1. Player clicks building -> context menu -> "Move"
2. **Placement mode activates** with the building's footprint
3. The original building becomes semi-transparent at its current position
4. Ghost preview follows the cursor with green/red feedback (excluding the building's own tiles from collision)
5. Player clicks a valid (green) position to move the building
6. Building animates to the new position
7. Grid and ghost disappear

### 5.4 3D Building Models

Each building type has a **unique procedural 3D model** built entirely in Three.js code (no external GLTF/GLB assets). Models are constructed from Three.js primitives (BoxGeometry, CylinderGeometry, SphereGeometry, etc.) with MeshStandardMaterial and emissive accents.

**Model Specifications**:

| Property | Value |
|----------|-------|
| Format | **Procedural Three.js** (code-generated geometry, no external files) |
| Polygon count | 500-2000 tris per model (low-poly style) |
| Scale | Proportional to building footprint (see 5.3.1); height 2-6 world units |
| Origin | Bottom-center of the model, placed at the center of the building's multi-tile footprint |
| Materials | `MeshStandardMaterial` with emissive accents per category color |
| Style | Low-poly sci-fi with emissive accents (glowing windows, energy lines) |

**Model List** (22 buildings):

| Building | Model Description | Key Visual Features |
|----------|------------------|---------------------|
| **Metal Collector** | Mining rig with drill arm | Rotating drill animation, gray/metallic |
| **He3 Extractor** | Gas collection tower with pipes | Glowing blue tubes, vapor particle effect |
| **Residential Area** | Domed city cluster | Warm yellow window glow, multiple small domes |
| **Resource Warehouse** | Large storage container/silo | Gray metallic crates, industrial look |
| **Civic Center** | Tall central command tower | Blue energy beacon at top, tallest building |
| **Technology Center** | Lab dome with antenna array | Purple/violet glow, holographic display ring |
| **Alliance Center** | Building with two interconnected towers | Two towers with connecting bridge, shield emblem |
| **Trading Center** | Market hall with open sides | Gold accents, rotating trade symbol |
| **Galaxy Transporter** | Portal arch structure | Swirling energy in the arch, teal glow |
| **Compound Center** | Facility with card slot aesthetic | Merging energy beams, green accents |
| **Radar** | Rotating dish on tower | Animated rotating dish, green sweep line |
| **Ship Factory** | Large hangar with construction frame | Red/orange interior glow, ship skeleton visible |
| **Spacedock** | Docking bay with launch pad | Blue landing lights, open bay |
| **Command Center** | Military command building | Gold trim, command screens visible |
| **Weapon Research Center** | Armored lab with weapon mounts | Red accent lights, turret shapes on roof |
| **Recycling Plant** | Processing facility with conveyors | Green glow, rotating recycler mechanism |
| **Space Station** | Orbital ring structure (elevated) | Floats above base, cyan ring, rotating section |
| **Meteor Star** | Rocky asteroid defense node | Rough rock texture, embedded metal plates |
| **Particle Cannon** | Energy turret on swivel | Purple energy buildup effect, tracking animation |
| **Anti-Aircraft Gun** | Rapid-fire flak turret | Red muzzle flash spots, multi-barrel |
| **Thor's Cannon** | Massive rail gun emplacement | Orange lightning arcs, heavy armored base |
| **Celestial Base** | Advanced dome on elevated platform | Cyan crystalline dome, hovering effect |

### 5.5 Building States in 3D

Each building model reflects its state visually:

| State | 3D Visual Effect |
|-------|-----------------|
| **Idle** | Normal materials, subtle emissive pulse (breathing glow) on accent lights |
| **Hover** | Outline highlight + **dark tooltip** appears: "Lv: X [Building Name]" (e.g., "Lv: 5 Metal Collector"). Tooltip disappears when mouse leaves. |
| **Selected** | Bright selection ring at base (emissive torus), **context menu** appears near building (View/Move/Upgrade) |
| **Upgrading** | **Yellow progress bar** displayed over the building sprite + particle sparks rising + scaffolding wireframe |
| **Max Level** | Gold emissive accents instead of default color, star particle effect |
| **Just Built** | Build-up animation: model scales from 0 to 1 over 1 second with particle burst |
| **Moving** | Building becomes semi-transparent, grid tiles visible with green/red indicators |

### 5.6 Hover Tooltip (CSS2DRenderer)

In GO2, buildings do NOT have permanent floating labels. Instead, **hovering the mouse** over a building displays a tooltip. This keeps the planet surface clean and uncluttered.

```
  (mouse hovers over building)
          |
   +---------------------+
   | Lv: 5 Metal Collector |    <-- tooltip appears on hover only
   +---------------------+
          |
   [3D Building Model]
```

**Tooltip Component** (rendered via `@react-three/drei` `<Html>` component):

```tsx
interface HoverTooltipProps {
  name: string           // e.g., "Metal Collector"
  level: number          // e.g., 5
  visible: boolean       // true only when mouse hovers over this building
}
```

Visual:
- Background: semi-transparent dark panel (`rgba(12, 12, 32, 0.9)`)
- Border: 1px solid `--border-frame`
- Text format: `"Lv: {level} {name}"` -- e.g., "Lv: 5 Metal Collector"
- Text: `--text-bright`, small font
- Position: 2 units above building model top
- Appears instantly on mouse hover, disappears on mouse leave
- NOT visible by default -- only on hover
- When building is upgrading: tooltip additionally shows a small yellow timer below the name

### 5.7 Building Context Menu (HTML Overlay)

When a building is left-clicked, a **context menu** appears near the building (not a direct panel opening). This is faithful to GO2's interaction pattern.

```
  [3D Building Model]
          |
   +-------------+
   |    View      |    <-- Opens BuildingDetailPanel
   +-------------+
   |    Move      |    <-- Enters move mode with grid
   +-------------+
   |   Upgrade    |    <-- Starts upgrade (if possible)
   +-------------+

   (Warehouse only: additional button)
   +-------------+
   |   Harvest    |    <-- Collects stored resources
   +-------------+
```

**Component**:
```tsx
interface BuildingContextMenuProps {
  building: BuildingWithType
  position: { x: number; y: number }   // screen position near the building
  onView: () => void                    // opens detail panel
  onMove: () => void                    // enters move mode
  onUpgrade: () => void                 // starts upgrade
  onHarvest?: () => void                // only for Warehouse
  onClose: () => void                   // close menu (click outside)
}
```

Visual:
- 3 vertically stacked blue buttons (View / Move / Upgrade)
- Warehouse has 4 buttons (View / Move / Upgrade / Harvest)
- Background: `--bg-panel` semi-transparent
- Border: 1px solid `--border-active`
- Position: HTML overlay, screen-space position calculated from the 3D building's world position projected to screen
- Closes when: clicking outside, pressing ESC, or clicking any option
- z-index: 150

### 5.8 Selection Indicator

When a building is left-clicked:

- A glowing ring appears at the building's base
  - Geometry: `TorusGeometry` at ground level, radius matching building footprint
  - Material: Emissive `--3d-glow-select` color, pulsing opacity animation
  - Animation: Slow rotation + breathing intensity
- The selected building's emissive intensity increases
- The **BuildingContextMenu** (HTML overlay) appears near the building
- Ring and menu disappear when the context menu is closed

### 5.9 Construction/Upgrade Effects

When a building is upgrading (faithful to GO2):

- **Yellow Progress Bar**: Displayed **over the building** sprite in 3D space
  - Rendered via CSS2DRenderer or a 3D plane above the building
  - Color: `--accent-warning` (yellow/orange)
  - Shows fill percentage based on time elapsed vs total time
  - Visible from any camera angle
- **Scaffolding**: Wireframe box around the building (slightly larger)
  - Material: Dashed wireframe, `--accent-warning` color
  - Fades in when upgrade starts, fades out when complete
- **Sparks**: Particle system emitting small bright particles upward
  - `PointsMaterial` with additive blending
  - Color: warm orange/yellow
  - Count: 20-40 particles, rising and fading
  - Emitter: positioned at building top
- **Base glow**: Ground ring changes to `--3d-glow-upgrade` color

### 5.10 Planet Header (HTML Overlay)

Positioned at top of main content area, below the Resource HUD:

```
+------------------------------------------------------------+
|  [Planet Icon]  Colony Alpha     (245, 128)     CC Lv 5    |
+------------------------------------------------------------+
```

- HTML overlay, position: absolute, top: 64px (below HUD), left: 68px (right of nav)
- Semi-transparent background
- Planet name: `--font-xl`, `--text-bright`, bold
- Coordinates: `--font-sm`, `--text-dim`
- Civic Center level: pill badge, gold border, `--accent-gold`

### 5.11 Warehouse Special Behavior

The **Resource Warehouse** has unique behavior faithful to GO2:

- **Limit**: Only **1 Warehouse per planet** is allowed
- **Storage**: Stores each resource type (Metal, He3, Gold) **separately** with individual capacities
- **Context menu**: Has **4 options** instead of the standard 3:
  - **View** - Opens detail panel showing storage levels per resource
  - **Move** - Repositions on the grid (same as other buildings)
  - **Upgrade** - Increases storage capacity
  - **Harvest** - Collects accumulated resources from the warehouse into the player's main pool

The Warehouse is the mechanism for collecting produced resources. Resource buildings (Metal Collector, He3 Extractor, Residential Area) produce resources that accumulate in the Warehouse. Players must **manually Harvest** via the context menu to add them to their usable resource pool displayed in the HUD.

### 5.12 Construction Info Panel (HTML Overlay, Bottom-Right)

In GO2, the construction info appears in the **lower-right corner** of the screen (not a full-width bottom bar). It shows a list of buildings currently under construction with countdown timers.

```
                                    +---------------------------+
                                    | Metal Collector Lv: 3     |
                                    |           00:12:45        |
                                    +---------------------------+
                                    | He3 Extractor Lv: 1       |
                                    |           00:00:32        |
                                    +---------------------------+
                                    | [Speed Up] icon available |
                                    +---------------------------+
                                           (bottom-right corner)
```

**Component**: `ConstructionInfoPanel`

```tsx
interface ConstructionInfoPanelProps {
  buildings: BuildingWithType[]    // filtered to is_upgrading === true
  maxSlots: number                 // GO2 has 5 construction slots
  onCancel: (buildingId: string) => void
  onFocusBuilding: (buildingId: string) => void   // pan camera to this building
}
```

- Position: fixed bottom-right, z-index 60
- Background: `--bg-panel` (semi-transparent)
- Shows only buildings with `is_upgrading === true`
- Each entry: building name, level, countdown timer
- Timer uses monospace font, updates every second
- Clicking an entry pans the camera to that building
- Compact vertical list, not a full-width bar
- GO2 allows **5 construction slots** simultaneously
- Optional: speed-up button (blue/yellow arrow icon) beside each timer

### 5.13 Data Requirements

This screen requires these API calls:
- `GET /api/planets/:id` - Planet name, coordinates
- `GET /api/planets/:id/buildings` - All buildings with types (drives 3D model placement)
- `GET /api/planets/:id/resources` - Resources for the HUD (passed up)

Auto-refresh: Every 30 seconds for resources, every 5 seconds for building timers (client-side countdown).

---

## 6. Screen: Building Detail Panel (HTML Overlay)

When a user clicks a building in the 3D scene and selects **"View"** from the context menu, a detail panel slides in from the right side as an HTML overlay. The panel does NOT open directly on click -- the context menu always appears first (faithful to GO2).

### 6.1 Layout

```
+------+-------------------------------+-----------------------+
|      |                               |  BUILDING DETAIL      |
|      |   3D CANVAS continues         |  HTML overlay z:200   |
| NAV  |   (camera can still pan)      |                       |
|      |                               |  [Close X]            |
|      |                               |                       |
|      |                               |  [Building Icon 96px] |
|      |                               |  Metal Collector      |
|      |                               |  Level 5 / 24         |
|      |                               |  Category: Resource   |
|      |                               |                       |
|      |                               |  --- PRODUCTION ---   |
|      |                               |  Output: 1,288/hr     |
|      |                               |  Next Lv: 1,378/hr    |
|      |                               |  (+7.0%)              |
|      |                               |                       |
|      |                               |  --- UPGRADE COST --- |
|      |                               |  Metal:  1,355   [OK] |
|      |                               |  He3:    1,690   [OK] |
|      |                               |  Gold:   1,355   [OK] |
|      |                               |  Time:   10m 38s      |
|      |                               |                       |
|      |                               |  [UPGRADE TO LV 6]    |
|      |                               |                       |
|      |                               |  --- INFO ---         |
|      |                               |  "Produces Metal..."  |
|      |                               |                       |
+------+-------------------------------+-----------------------+
```

### 6.2 Panel Width & Positioning

- Desktop: 360px fixed width, right-anchored, full viewport height minus HUD
- Tablet: 320px
- Mobile: Full screen overlay
- Background: `--bg-panel-solid` (opaque, to keep readability against 3D scene)
- Left border: 1px solid `--border-glow` with box-shadow glow
- The 3D scene behind the panel is still visible and interactive (camera panning)

### 6.3 Panel Sections

**Header Section**:
```
+---------------------------------------+
|  [X Close]              [Category]    |
|                                       |
|        [Building Icon 96x96]          |
|                                       |
|        Metal Collector                |
|        Level 5 / 24                   |
|        [=========>    ] 5/24          |
+---------------------------------------+
```

- Close button: top-right, `X` icon
- Category badge: top-right below close, colored pill
- Building 2D icon: centered, 96x96px (same icon as in construction panel, not a 3D render)
- Name: `--font-lg`, bold
- Level: `--font-md`, level/maxLevel with a thin progress bar (5/24 = ~21%)

**Production Section** (only for resource buildings):
```
+---------------------------------------+
|  PRODUCTION                           |
|  Current: 1,288 Metal/hr              |
|  Next Lv: 1,378 Metal/hr (+7.0%)     |
+---------------------------------------+
```

**Special Effect Section** (for non-resource buildings):
```
+---------------------------------------+
|  EFFECT                               |
|  Research Time Reduction: 15%         |
|  Next Lv: 18%  (+3%)                 |
+---------------------------------------+
```

**Upgrade Cost Section**:
```
+---------------------------------------+
|  UPGRADE TO LEVEL 6                   |
|                                       |
|  [Metal Icon]  Metal    1,355   [OK]  |
|  [He3 Icon]    He3      1,690   [OK]  |
|  [Gold Icon]   Gold     1,355   [OK]  |
|  [Clock Icon]  Time     10m 38s       |
|                                       |
|  Requires: Civic Center Lv 3   [OK]   |
|                                       |
|  [ ===== UPGRADE ===== ]              |
+---------------------------------------+
```

- Each resource line: icon, name, cost, status indicator (green check / red X)
- Time: formatted as Xh Xm Xs
- Prerequisites: Civic Center level requirement
- Upgrade button states: same as v1 spec (see 6.4)

### 6.4 Upgrade States

| State | Button | 3D Effect |
|-------|--------|-----------|
| Can upgrade | Blue gradient, "UPGRADE TO LV X" | - |
| Insufficient resources | Gray, disabled, missing in red | - |
| Already upgrading | Yellow, "UPGRADING... HH:MM:SS" with progress | Scaffolding + sparks on 3D model |
| Max level | Gold border, "MAX LEVEL" | Gold emissive on model |
| Missing prerequisite | Gray, disabled, requirement in red | - |
| No construction slot | Gray, disabled, "All slots busy" | - |

### 6.5 Component

```tsx
interface BuildingDetailPanelProps {
  building: BuildingWithType | null        // null = panel hidden
  resources: ResourcesResponse
  upgradeCost: { metal: number; he3: number; gold: number; time_seconds: number }
  isUpgrading: boolean
  canUpgrade: boolean
  insufficientResources: { metal: boolean; he3: boolean; gold: boolean }
  onUpgrade: () => void
  onCancel: () => void
  onClose: () => void
}
```

### 6.6 Data Requirements

- Building data from `GET /api/planets/:id/buildings` (already loaded)
- Resources from `GET /api/planets/:id/resources` (already loaded)
- Upgrade cost: Calculated client-side or fetched
- Uses: `POST /api/planets/:id/buildings/:buildingId/upgrade`
- Uses: `POST /api/planets/:id/buildings/:buildingId/cancel`

---

## 7. Screen: Construction/Upgrade Panel (HTML Overlay)

When clicking the **"Build" button** in the lower-right corner of the screen, a modal overlay appears for selecting what to build. In GO2, this button is always visible in the bottom-right area. There are no fixed "empty slot" markers on the grid -- the player selects a building type first, then chooses where to place it on the isometric grid.

### 7.1 Layout

Same as v1 spec -- this is a pure HTML modal overlay over the 3D canvas:

```
+------------------------------------------------------------+
|  BUILD NEW STRUCTURE                              [X Close] |
+------------------------------------------------------------+
|                                                             |
|  [Tab: Resource] [Tab: Core] [Tab: Military] [Tab: Space]  |
|                                                             |
|  +------------------+  +------------------+                 |
|  | [Icon 48x48]     |  | [Icon 48x48]     |                |
|  | Metal Collector  |  | He3 Extractor    |                 |
|  | Built: 3/8       |  | Built: 2/8       |                 |
|  | M: 95  H: 80     |  | M: 95  H: 80    |                 |
|  | G: 95  T: 40s    |  | G: 95  T: 40s   |                 |
|  | [BUILD]          |  | [BUILD]          |                 |
|  +------------------+  +------------------+                 |
|                                                             |
+------------------------------------------------------------+
```

### 7.2 Category Tabs

| Tab | Category Filter | Color Accent |
|-----|----------------|--------------|
| Resource | `resource` | `--accent-success` (green) |
| Core | `core` | `--accent-primary` (blue) |
| Military | `military` | `--accent-danger` (red) |
| Space | `space`, `defense` | `--color-he3` (cyan) |

### 7.3 Construction Card States

| State | Visual |
|-------|--------|
| Available | Normal card, blue "BUILD" button |
| Max count reached | Dimmed card, "MAX BUILT" badge |
| Missing prerequisite | Dimmed, red prerequisite text, disabled button |
| Insufficient resources | Red cost numbers, disabled button |
| Building in progress | Button shows "BUILDING..." |

### 7.4 After Build Action

When the user clicks "BUILD" on a building card:
1. Modal closes
2. **Placement mode activates**: isometric grid becomes visible with green/red tile indicators
3. Player clicks a **green tile** to place the building
4. API call fires: `POST /api/planets/:id/buildings` with grid coordinates
5. New building model spawns at the chosen tile with a build-up animation (scale 0 -> 1 with particle burst)
6. **Yellow progress bar** appears over the building
7. **Construction Info Panel** (bottom-right) updates with a countdown timer
8. Camera pans to the new building

If the player presses ESC during placement mode, construction is cancelled and the grid hides.

### 7.5 Component

```tsx
interface ConstructionPanelProps {
  buildings: BuildingWithType[]
  buildingTypes: BuildingType[]
  resources: ResourcesResponse
  activeCategory: string
  onCategoryChange: (category: string) => void
  onConstruct: (buildingTypeName: string) => void
  constructingType: string | null
  onClose: () => void
}
```

### 7.6 Presentation

- Centered modal overlay with dark backdrop (`rgba(0,0,0,0.7)`) over 3D canvas
- z-index: 300
- Modal width: 680px desktop, 100% mobile
- Building cards in a 2-column grid
- Close on backdrop click, X button, or Escape

---

## 8. Component: Resource HUD (HTML Overlay)

The Resource HUD is a fixed HTML overlay at the top of the viewport, rendered over the Three.js canvas.

### 8.1 Layout

```
+--------------------------------------------------------------------------+
|                                                                           |
|  [Logo]  CRYPTOMINES    ||  [Metal Icon]  45,230    +8,640/hr           |
|          ONLINE         ||  [He3 Icon]    32,100    +9,440/hr           |
|                         ||  [Gold Icon]   71,500    +11,200/hr          |
|  Civic Center Lv 5     ||                                               |
|                         ||  Storage: 500K / 500K    [COLLECT +12.5K]    |
|                                                                           |
+--------------------------------------------------------------------------+
```

### 8.2 Sections

**Left Section** (Logo + Player Info):
- Game logo/title: "CRYPTOMINES ONLINE" in `--accent-primary`, uppercase, letter-spaced
- Below: Civic Center level badge (gold pill)
- Width: ~200px

**Center Section** (Resources):
- Three resource rows, each with:
  - Resource icon (16x16, color-coded)
  - Current amount (bold, formatted with K/M suffixes)
  - Production rate (/hr, in green, smaller text)
- If pending resources > 0, show total pending below in yellow

**Right Section** (Storage + Collect):
- Storage capacity bar: thin horizontal bar
- Storage text: "current / max"
- Collect button: green gradient, shows pending amount

### 8.3 Visual Styling

- Height: 56px on desktop, 48px on mobile
- Background: `--bg-header` (semi-transparent so starfield slightly shows through)
- Position: fixed top, z-index 100
- Bottom border: 1px solid `--border-active` with box-shadow glow
- `backdrop-filter: blur(8px)` for frosted glass effect over 3D scene
- Resources use `font-variant-numeric: tabular-nums` for stable width

### 8.4 Component

```tsx
interface ResourceHUDProps {
  resources: ResourcesResponse | null
  civicCenterLevel: number
  onCollect: () => void
  collecting: boolean
}
```

### 8.5 Resource Number Formatting

```
< 1,000:        Show exact (e.g., "950")
1,000-999,999:  Show with comma separators (e.g., "45,230")
1M-999M:        Show as X.XXM (e.g., "12.53M")
1B+:            Show as X.XXB (e.g., "1.25B")
```

### 8.6 Data Requirements

- Resources from `GET /api/planets/:id/resources`
- Auto-refresh: every 30 seconds
- Client-side interpolation: Increment displayed values between refreshes based on production rates

---

## 9. Reusable Components

### 9.1 Component Catalog

**3D Components** (inside `<Canvas>`):

| Component | Description | Used In |
|-----------|-------------|---------|
| `PlanetScene` | Root 3D scene with lighting and skybox | Planet View |
| `PlanetSurface` | Ground plane scaled for 20x20 grid | Planet View |
| `IsometricGrid` | 20x20 wireframe diamond grid (visible in placement mode) | Planet View |
| `GhostPreview` | Translucent building preview following cursor during placement | Planet View |
| `PlacementIndicator` | Green/red multi-tile footprint overlays during placement | Planet View |
| `BuildingModel` | Individual 3D building instance on multi-tile footprint | Planet View |
| `HoverTooltip` | "Lv: X [Name]" tooltip on mouse hover | Planet View |
| `SelectionRing` | Glowing ring at selected building base | Planet View |
| `UpgradeProgressBar3D` | Yellow progress bar above upgrading building | Planet View |
| `UpgradeParticles` | Spark particle system for upgrading | Planet View |
| `Skybox` | Starfield background | Planet View |
| `CameraController` | Pan-based camera with zoom (fixed isometric angle) | Planet View |

**2D Components** (HTML overlay):

| Component | Description | Used In |
|-----------|-------------|---------|
| `GameShell` | App wrapper with canvas + overlay layout | App root |
| `ResourceHUD` | Fixed top resource bar | All screens |
| `SideNav` | Left navigation icon bar | All screens |
| `BuildingContextMenu` | View/Move/Upgrade popup on building click | Planet View |
| `ConstructionInfoPanel` | Bottom-right panel with countdown timers | Planet View |
| `BuildingDetailPanel` | Right slide-in detail panel (via "View") | Planet View |
| `ConstructionPanel` | Modal for selecting new building (via "Build" button) | Planet View |
| `ProgressBar` | Thin horizontal progress bar | Multiple |
| `CountdownTimer` | Live countdown (monospace) | Multiple |
| `ResourceCost` | Resource icon + amount + status | Detail, Construction |
| `CategoryTabs` | Segmented tab buttons | Construction Panel |
| `Tooltip` | Hover tooltip | Multiple |
| `Modal` | Centered overlay modal | Construction Panel |
| `Badge` | Small pill badge (level, status) | Multiple |
| `IconButton` | Icon-only button with tooltip | SideNav, panels |
| `GlowButton` | Primary action button with glow | Upgrade, Build, Collect |
| `Toast` | Notification toast | Error/success events |

### 9.2 GameShell (Updated for Three.js)

```tsx
interface GameShellProps {
  children: ReactNode
}

// Key change: manages the Canvas + HTML overlay stack
// Canvas is full viewport, HTML floats on top
// GameShell holds the GameContext provider
```

### 9.3 BuildingModel (3D Component)

```tsx
interface BuildingModelProps {
  buildingType: string              // e.g., 'civic_center'
  gridPosition: { col: number; row: number }  // isometric grid coordinates (top-left of footprint)
  footprint: { cols: number; rows: number }   // multi-tile size from buildingConfig
  position: [number, number, number]  // 3D world position (center of footprint)
  level: number
  maxLevel: number
  isUpgrading: boolean
  isSelected: boolean
  isHovered: boolean
  onClick: () => void
  onPointerOver: () => void
  onPointerOut: () => void
}
```

Behavior:
- Renders the **procedural Three.js model** for `buildingType` (dispatches to type-specific component)
- Placed at the 3D world position corresponding to the **center of its multi-tile footprint**
- Model scale is proportional to footprint size (3x3 buildings are visually larger than 1x1)
- Applies level-dependent scale (higher level = slightly larger, up to 1.2x)
- Idle animation: subtle floating or breathing glow
- Hover: outline effect + HoverTooltip appears ("Lv: X [Name]")
- Selected: SelectionRing child visible, BuildingContextMenu (HTML) appears
- Upgrading: Yellow progress bar + UpgradeParticles + scaffolding wireframe
- Raycaster interaction for click/hover detection

### 9.4 CameraController (3D Component)

```tsx
interface CameraControllerProps {
  target: [number, number, number] | null   // null = default overview, set to pan to building
  enabled: boolean                           // false when modal is open
  enableRotate: false                        // always false -- fixed isometric angle
}
```

See Section 12 for full camera details. Key: pan-based (drag to scroll), NO orbit/rotation.

### 9.5 Other 2D components

ProgressBar, CountdownTimer, ResourceCost, GlowButton, Tooltip, Modal, Badge -- same specs as v1 (see original definitions in previous version, carried forward here unchanged).

---

## 10. 3D Assets Required

### 10.1 Building Models (22 procedural models)

**Format**: Procedural Three.js code -- each building type has a dedicated `.tsx` component that constructs geometry from Three.js primitives (Box, Cylinder, Sphere, Cone, Torus, etc.). No external GLTF/GLB files.

**Style**: Low-poly sci-fi. Think "Stellaris meets Monument Valley" -- clean geometric shapes with emissive accent details. Not photorealistic, not voxel. Smooth low-poly with sharp edges and glowing energy lines.

**Per-Model Spec**:

| Property | Value |
|----------|-------|
| Triangle count | 500 - 2,000 per model |
| Bounding box | Proportional to footprint: width/depth match tile coverage, height 2-6 units |
| Materials | `MeshStandardMaterial` with PBR properties (metalness, roughness, emissive) |
| Origin point | Bottom-center of the multi-tile footprint |
| Coordinate system | Y-up |
| Animations | Optional idle loop via useFrame (rotating parts, pulsing emissive) |

**Complete Model List**:

| # | Building | File Name | Priority | Description |
|---|----------|-----------|----------|-------------|
| 1 | Civic Center | `civic_center.glb` | HIGH | Tall command tower with blue energy beacon at top |
| 2 | Metal Collector | `metal_collector.glb` | HIGH | Mining rig with mechanical drill arm |
| 3 | He3 Extractor | `he3_extractor.glb` | HIGH | Gas collection tower with glowing blue tubes |
| 4 | Residential Area | `residential_area.glb` | HIGH | Domed city cluster with warm yellow window glow |
| 5 | Resource Warehouse | `resource_warehouse.glb` | HIGH | Large storage containers, industrial |
| 6 | Technology Center | `technology_center.glb` | HIGH | Lab dome with antenna array, purple glow |
| 7 | Ship Factory | `ship_factory.glb` | MEDIUM | Large hangar with ship frame inside |
| 8 | Command Center | `command_center.glb` | MEDIUM | Military command, gold trim |
| 9 | Space Station | `space_station.glb` | MEDIUM | Orbital ring, elevated/floating above base |
| 10 | Alliance Center | `alliance_center.glb` | MEDIUM | Two connected towers |
| 11 | Trading Center | `trading_center.glb` | MEDIUM | Market hall, gold exchange symbol |
| 12 | Spacedock | `spacedock.glb` | MEDIUM | Docking bay with landing pad |
| 13 | Radar | `radar.glb` | MEDIUM | Rotating dish tower |
| 14 | Weapon Research Center | `weapon_research_center.glb` | MEDIUM | Armored lab with turret shapes |
| 15 | Recycling Plant | `recycling_plant.glb` | LOW | Processing facility with conveyor |
| 16 | Galaxy Transporter | `galaxy_transporter.glb` | LOW | Portal arch with swirling energy |
| 17 | Compound Center | `compound_center.glb` | LOW | Card merging facility |
| 18 | Meteor Star | `meteor_star.glb` | LOW | Rocky asteroid defense node |
| 19 | Particle Cannon | `particle_cannon.glb` | LOW | Energy turret |
| 20 | Anti-Aircraft Gun | `anti_aircraft_gun.glb` | LOW | Multi-barrel flak turret |
| 21 | Thor's Cannon | `thors_cannon.glb` | LOW | Massive rail gun |
| 22 | Celestial Base | `celestial_base.glb` | LOW | Crystalline dome, hovering |

### 10.2 Environment Assets

| Asset | Format | Description |
|-------|--------|-------------|
| **Skybox** | 6x cube map PNG (2048x2048 each) or HDR equirectangular | Deep space starfield with nebula hints |
| **Planet Surface Texture** | Diffuse + Normal + Roughness (2048x2048) | Dark rocky/metallic terrain |
| **Grid Overlay Texture** | PNG with alpha (1024x1024) | Subtle hexagonal or square grid pattern |
| **Building Platform** | GLTF or procedural | Hexagonal pad where buildings sit |

### 10.3 Effect Assets

| Asset | Type | Description |
|-------|------|-------------|
| **Selection Ring** | Procedural (code) | Emissive torus geometry, animated |
| **Upgrade Sparks** | Particle texture (64x64 PNG) | Small bright dot/star for spark particles |
| **Construction Scaffolding** | Procedural wireframe (code) | Box wireframe overlay |
| **Build-up Effect** | Particle texture (64x64 PNG) | Dust/energy particles for new construction |
| **Ambient Dust** | Particle texture (32x32 PNG) | Floating dust motes in atmosphere |

### 10.4 Placeholder Strategy (Fallback)

All 22 building types have unique procedural models. If a model is missing or fails to load, fall back to a colored box placeholder:

```tsx
// Placeholder 3D building: colored box with abbreviation as 3D text
function PlaceholderBuilding({ type, category }: { type: string; category: string }) {
  const colors: Record<string, string> = {
    resource: '#22cc66',
    core: '#4488ff',
    military: '#ff4444',
    space: '#44ccff',
    defense: '#44ccff',
  }

  return (
    <group>
      {/* Base: simple box geometry */}
      <mesh>
        <boxGeometry args={[3, 3 + Math.random(), 3]} />
        <meshStandardMaterial
          color={colors[category] || '#888888'}
          emissive={colors[category] || '#888888'}
          emissiveIntensity={0.15}
          metalness={0.7}
          roughness={0.3}
        />
      </mesh>
      {/* Hover tooltip handled by CSS2DRenderer */}
    </group>
  )
}
```

Placeholder abbreviations (same as 2D spec):
MC, HE, RA, RW, CC, TC, AC, TR, GT, CP, RD, SF, SD, CM, WR, RP, SS, MS, PC, AA, TH, CB

---

## 11. 2D Assets Required

### 11.1 Building Icons (for panels/HUD, not for 3D)

Each building still needs a 2D icon for use in HTML panels (detail panel header, construction panel cards, construction queue).

Format: PNG with transparency, 64x64 base (scaled to 48px and 96px).

Same list as v1 spec -- 22 building icons, matching the visual style of the 3D models but as flat icons.

### 11.2 Resource Icons (3 icons)

| Resource | Size | Description |
|----------|------|-------------|
| Metal | 16x16, 24x24 | Steel/silver ingot or ore chunk |
| He3 | 16x16, 24x24 | Glowing blue atom/molecule |
| Gold | 16x16, 24x24 | Gold coin or bar |

### 11.3 UI Icons (8+ icons)

| Icon | Description | Used In |
|------|-------------|---------|
| Clock | Timer/countdown | Upgrade times |
| Checkmark | Success/available | Resource cost status |
| X / Cross | Error/insufficient | Resource cost status |
| Plus | Add/build new | Build button, construction |
| Star | Max level | Max level badge |
| Lock | Locked/unavailable | Disabled nav items |
| Arrow Right | Slide panel | Detail panel |
| Chevron | Expand/collapse | Sections |

### 11.4 2D Placeholder Strategy

Same as v1: colored squares with 2-letter abbreviations for building icons until final art is ready.

---

## 12. Camera System

### 12.1 Camera Type

**Pan-based camera** (faithful to GO2). The camera looks down at the planet surface from a fixed isometric angle. The player drags to pan the viewport, NOT to orbit.

In GO2, the planet view uses a **fixed isometric perspective** with camera panning (click-and-drag to scroll around the surface). There is no free camera rotation or orbit.

```tsx
// Camera configuration
const CAMERA_CONFIG = {
  // Orthographic or perspective with fixed angle (isometric-like)
  type: 'perspective',       // can use orthographic for true isometric
  fov: 45,
  near: 0.1,
  far: 1000,
  // Fixed position: looking down at the base from isometric angle
  fixedAngle: { x: -45, y: 45 },  // degrees, isometric-like tilt
  initialPosition: [30, 40, 30],    // x, y, z (fixed angle, adjustable height)
  initialTarget: [0, 0, 0],         // looking at center of base
}
```

### 12.2 Camera Controls

| Control | Action |
|---------|--------|
| Left-click drag (on empty space) | **Pan** the camera across the surface (scroll viewport) |
| Left-click (on building) | Open context menu (NOT camera drag) |
| Scroll wheel | Zoom in/out (move camera closer/further along the fixed angle) |
| Click on ConstructionInfoPanel entry | Pan camera to that building |

**Important**: There is NO orbit/rotation. The camera angle is fixed (isometric). Only panning and zooming are available, like in GO2.

### 12.3 Camera Constraints

```tsx
const PAN_CONSTRAINTS = {
  // Zoom limits
  minZoom: 15,               // Cannot zoom closer than 15 units
  maxZoom: 80,               // Cannot zoom farther than 80 units
  // Pan limits (keep buildings in view)
  panBounds: {
    minX: -40, maxX: 40,
    minZ: -40, maxZ: 40,
  },
  // Smooth panning
  enableDamping: true,
  dampingFactor: 0.08,
  panSpeed: 1.0,
  zoomSpeed: 0.5,
  // Fixed angle (no rotation)
  enableRotate: false,
}
```

### 12.4 Camera Transitions

When the user clicks a building in the ConstructionInfoPanel or after placing a new building, the camera smoothly pans to center on it:

```
Current pan position -----> Target pan position
                      500ms ease-in-out

Target calculation:
  - Center the camera view on the building's world position
  - Keep the same fixed isometric angle and zoom level
  - Smooth interpolation using lerp over ~500ms
```

```tsx
// Pseudo-code for camera pan-to-building
function panToBuilding(buildingPosition: Vector3) {
  const targetLookAt = buildingPosition.clone()
  const targetCameraPos = targetLookAt.clone().add(cameraOffset)  // fixed angle offset

  // Animate over 500ms with ease-in-out
  animateCameraTo(targetCameraPos, targetLookAt, 500)
}
```

### 12.5 Camera States

| State | Behavior |
|-------|----------|
| Overview (default) | Camera at initial position, centered on the base, fixed isometric angle |
| Panning | Camera moves laterally as player drags, maintaining fixed angle |
| Building focused | Camera pans to center on a building (after placement or info panel click) |
| Modal open | Camera controls disabled (prevents clicks going to 3D scene) |
| Placement mode | Camera panning still allowed so player can scroll to find a valid tile |

---

## 13. 3D Visual Effects & Lighting

### 13.1 Lighting Setup

```
Scene Lights:

1. AmbientLight
   Color: #1a1a3a
   Intensity: 0.4
   Purpose: Base illumination so nothing is fully black

2. DirectionalLight (Key Light - simulates a nearby star)
   Color: #aabbff (cool white-blue)
   Intensity: 1.2
   Position: (20, 30, 10)
   Casts shadows: yes
   Shadow map: 2048x2048

3. DirectionalLight (Rim/Fill Light)
   Color: #4466aa (blue tint)
   Intensity: 0.3
   Position: (-15, 10, -20)
   Casts shadows: no

4. PointLights (per building, optional)
   Color: building category color
   Intensity: 0.5
   Distance: 6
   Purpose: Local glow around each building's emissive elements
```

### 13.2 Post-Processing (Optional Enhancement)

Using `@react-three/postprocessing`:

| Effect | Purpose | Priority |
|--------|---------|----------|
| **Bloom** | Glow on emissive materials (building lights, selection ring) | HIGH |
| **Vignette** | Subtle edge darkening for cinematic feel | LOW |
| **ChromaticAberration** | Very subtle, adds sci-fi feel | LOW |

Bloom configuration:
```tsx
<Bloom
  intensity={0.5}
  luminanceThreshold={0.8}
  luminanceSmoothing={0.1}
  mipmapBlur
/>
```

### 13.3 Material Standards

All 3D materials should use PBR (Physically Based Rendering):

| Property | Resource Buildings | Core Buildings | Military Buildings | Space Buildings |
|----------|-------------------|----------------|-------------------|----------------|
| Metalness | 0.6 | 0.5 | 0.8 | 0.7 |
| Roughness | 0.4 | 0.3 | 0.3 | 0.2 |
| Emissive color | Green tint | Blue tint | Red/orange tint | Cyan tint |
| Emissive intensity | 0.1 (idle), 0.3 (selected) | Same | Same | Same |

### 13.4 Animations

| Animation | Type | Description |
|-----------|------|-------------|
| Building idle | Emissive pulse | Slow sine wave on emissiveIntensity (period: 3s) |
| Radar dish rotation | Mesh rotation | Continuous Y-axis rotation (1 revolution per 4s) |
| Selection ring | Rotation + opacity | Slow Y-rotation + breathing opacity (0.5-1.0) |
| Upgrade sparks | Particle system | 30 particles/s, rising, randomized positions |
| Build-up spawn | Scale animation | Scale from 0.01 to 1.0 over 1s with bounce easing |
| Camera transition | Position lerp | 500ms ease-in-out position/target interpolation |
| Skybox stars | Subtle drift | Very slow rotation of skybox (barely noticeable) |
| He3 Extractor glow | Emissive pulse | Blue tube glow pulses faster than standard (period: 2s) |

---

## 14. Responsive Design

### 14.1 Breakpoints

```
DESKTOP:    >= 1200px   Full layout (sidebar + 3D canvas + detail panel)
TABLET:     768-1199px  Narrower sidebar, 3D canvas adjusts
MOBILE:     < 768px     No sidebar, hamburger, simplified 3D
```

### 14.2 Desktop (>= 1200px)

```
+------+-------------------------------------------+-------------------+
| 60px |         THREE.JS CANVAS (fills space)      |    360px          |
| Nav  |         3D Planet + isometric grid          | Detail Panel      |
|      |         Pan camera (drag to scroll)         | (via View, HTML)  |
+------+-------------------------------------------+---+---------------+
|                                              | Construction Info     |
|  [BUILD]  [MILITARY]  (bottom-right)         | Panel (timers, HTML)  |
+----------------------------------------------+-----------------------+
```

### 14.3 Tablet (768-1199px)

```
+----+--------------------------------------------------+
|48px|         THREE.JS CANVAS (fills space)             |
| Nav|         3D Planet + buildings (isometric grid)     |
|    |         Touch pan controls (drag to scroll)        |
+----+--------------------------------------------------+
| [Construction Queue - compact]                         |
+--------------------------------------------------------+
```

- SideNav: 48px, icons only
- Detail Panel: 320px overlay (does not push canvas, floats on top)
- Camera initial distance increased (zoom out more)

### 14.4 Mobile (< 768px)

```
+----------------------------------------------+
| [Ham.] CRYPTOMINES [Resources compact]       |  HUD
+----------------------------------------------+
|                                              |
|     THREE.JS CANVAS (fills viewport)          |
|     3D Planet + buildings (isometric grid)     |
|     Touch: pinch zoom, 1-finger pan            |
|                                              |
+----------------------------------------------+
| [Queue: compact 1-line]                      |
+----------------------------------------------+
```

- No sidebar, hamburger menu opens nav as drawer overlay
- Resource HUD: compact mode (icons + numbers only)
- Detail Panel: full-screen overlay (covers 3D)
- Construction Panel: full-screen modal
- Camera: touch-optimized controls (pinch zoom, single-finger pan to scroll viewport)
- Hover tooltips: tap-and-hold to show building tooltip (no hover on mobile)

### 14.5 3D Performance Tiers

On lower-end devices, reduce 3D complexity:

| Tier | Detection | Adjustments |
|------|-----------|-------------|
| HIGH | Desktop GPU, > 4GB VRAM | Full effects, shadows, bloom, particles |
| MEDIUM | Integrated GPU or tablet | No bloom, reduced shadow map (1024), fewer particles |
| LOW | Mobile or old device | No shadows, no post-processing, simplified materials, no particles |

Detection can use `renderer.capabilities` or a simple FPS check on first 60 frames.

---

## 15. API Integration Map

### 15.1 Screen-to-Endpoint Mapping

| Screen / Component | Endpoints Used | Trigger | Layer |
|---------------------|---------------|---------|-------|
| **ResourceHUD** | `GET /api/planets/:id/resources` | Auto-refresh 30s | HTML |
| **ResourceHUD** (collect) | `POST /api/planets/:id/resources/collect` | Click "Collect" | HTML |
| **Warehouse** (harvest) | `POST /api/planets/:id/resources/collect` | Context menu "Harvest" | HTML |
| **Planet 3D Scene** (load) | `GET /api/planets/:id` | Page mount | 3D |
| **Planet 3D Scene** (buildings) | `GET /api/planets/:id/buildings` | Page mount | 3D + HTML |
| **Context Menu** (upgrade) | `POST /api/planets/:id/buildings/:bid/upgrade` | Context menu "Upgrade" | HTML -> 3D |
| **Building Detail** (cancel) | `POST /api/planets/:id/buildings/:bid/cancel` | Click "Cancel" in info panel | HTML -> 3D |
| **Construction Panel** (build) | `POST /api/planets/:id/buildings` | Grid tile click after selection | HTML -> 3D |
| **Initial Auth** | `POST /api/auth/guest` | App mount | HTML |
| **Player Info** | `GET /api/player/me` | App mount | HTML |
| **Planet List** | `GET /api/planets` | App mount | HTML |

### 15.2 Data Flow

```
App Mount
  |
  +--> POST /api/auth/guest (if no stored token)
  |      +--> Store JWT in localStorage
  |
  +--> GET /api/player/me
  |      +--> Store player data in GameContext
  |
  +--> GET /api/planets
         |
         +--> Set first planet as current
                |
                +--> GET /api/planets/:id
                +--> GET /api/planets/:id/buildings
                +--> GET /api/planets/:id/resources
                      |
                      +--> Initialize 3D scene with building data
                      +--> Render HTML overlay with resource data
                      +--> Start 30s resource refresh interval
```

### 15.3 State Management

```tsx
// GameContext provides:
interface GameState {
  player: Player | null
  planets: Planet[]
  currentPlanet: Planet | null
  buildings: BuildingWithType[]
  resources: ResourcesResponse | null
  selectedBuilding: BuildingWithType | null    // drives 3D selection ring + context menu
  hoveredBuilding: BuildingWithType | null     // drives 3D hover tooltip
  showContextMenu: boolean                     // context menu visible (View/Move/Upgrade)
  showDetailPanel: boolean                     // detail panel visible (opened via "View")
  showConstructPanel: boolean                  // construction modal visible (via "Build" button)
  placementMode: {                             // active during build or move
    active: boolean
    buildingTypeId?: string                    // for new construction
    buildingTypeName?: string                  // building_type_name for footprint lookup
    movingBuildingId?: string                  // for moving existing
    ghostPosition?: { col: number; row: number } // current ghost snap position
    isValid?: boolean                          // green (valid) or red (invalid) ghost state
  }
  cameraTarget: [number, number, number] | null  // drives 3D camera pan
}
```

State changes flow bidirectionally:
- 3D click -> updates `selectedBuilding` -> context menu appears
- Context menu "View" -> `showDetailPanel = true` -> HTML panel slides in
- Context menu "Move" -> `placementMode.active = true` -> grid tiles become visible with green/red
- Context menu "Upgrade" -> API call -> yellow progress bar + construction info panel
- Context menu "Harvest" (Warehouse) -> API call -> resources update in HUD
- "Build" button click -> `showConstructPanel = true` -> construction modal opens
- Build modal select -> `placementMode.active = true` -> wireframe grid fades in, ghost preview follows cursor
- Cursor move in placement mode -> ghost snaps to grid, green/red feedback on footprint tiles
- Click valid position in placement mode -> API call -> new building spawns -> grid + ghost hide

### 15.4 Error Handling UI

| Error Type | UI Response |
|------------|-------------|
| Network error | Red toast notification (HTML), auto-dismiss 5s |
| Insufficient resources | Red highlight on costs in detail panel |
| Max buildings reached | "MAX BUILT" badge in construction panel |
| Missing prerequisite | Red text in detail/construction panels |
| Server error | Red toast |
| Auth expired | Auto-redirect to auth |
| WebGL not supported | Full-screen HTML fallback message |
| 3D model load failure | Use placeholder box, log warning |

### 15.5 Loading States

| Component | Loading State |
|-----------|--------------|
| 3D Scene | Black canvas with centered loading spinner + "Loading base..." text (HTML overlay) |
| Building models | Progressive loading -- placeholder boxes replaced by real models as they load |
| Isometric Grid | Loads instantly (procedural geometry) |
| Building Detail | Skeleton lines |
| Resource HUD | Three placeholder boxes with shimmer |
| Construction Panel | Skeleton cards |
| Context Menu | Instant (HTML, no loading) |
| Camera pan transition | Smooth interpolation (no loading state needed) |

---

## 16. Performance Considerations

### 16.1 3D Scene Optimization

| Technique | Description |
|-----------|-------------|
| **Instanced Meshes** | Use `InstancedMesh` for buildings of the same type (e.g., 8x Metal Collectors share one geometry) |
| **LOD** | Not needed for Phase 1 (low-poly models already optimized) |
| **Texture Atlasing** | Combine building textures into shared atlas where possible |
| **Frustum Culling** | Enabled by default in Three.js -- buildings behind camera are not rendered |
| **Shadow Optimization** | Only the key light casts shadows; shadow map 2048px; buildings close to camera only |
| **Object Pooling** | Reuse particle system instances rather than creating/destroying |
| **Dispose on Unmount** | All geometries, materials, and textures must be disposed when component unmounts |

### 16.2 Target Frame Rate

| Platform | Target FPS |
|----------|-----------|
| Desktop | 60 FPS |
| Tablet | 30-60 FPS |
| Mobile | 30 FPS |

If FPS drops below target, auto-reduce quality tier (see Section 14.5).

### 16.3 Asset Loading Strategy

```
1. Load critical HTML (HUD, nav) immediately
2. Initialize Three.js canvas with skybox + ground plane
3. Load placeholder boxes for all buildings (instant)
4. Async load GLTF models in priority order (HIGH first)
5. Replace placeholder with real model as each loads (fade transition)
6. Load post-processing effects last
```

Use `@react-three/drei`'s `useGLTF.preload()` for high-priority models and `Suspense` for lazy loading.

### 16.4 Memory Budget

| Asset Type | Budget |
|------------|--------|
| Building models (22) | ~5-10 MB total GLTF |
| Textures | ~10-15 MB (compressed) |
| Skybox | ~2-4 MB (cube map) |
| Particle textures | < 0.5 MB |
| **Total** | **~20-30 MB** |

---

## Appendix A: File Structure (Proposed)

```
frontend/src/
  components/
    layout/
      GameShell.tsx           // Canvas + overlay wrapper
      ResourceHUD.tsx          // Fixed top bar (HTML)
      SideNav.tsx              // Left nav bar (HTML)
    three/                     // All Three.js / R3F components
      PlanetScene.tsx          // Root scene component
      PlanetSurface.tsx        // Ground plane mesh (scaled for 20x20 grid)
      IsometricGrid.tsx        // 20x20 wireframe diamond grid overlay
      PlacementIndicator.tsx   // Green/red multi-tile footprint feedback
      GhostPreview.tsx         // Translucent building preview during placement
      BuildingModel.tsx        // Dispatches to type-specific procedural model
      buildings/               // 22 procedural building model components
        CivicCenterModel.tsx
        MetalCollectorModel.tsx
        ... (one per building type)
      HoverTooltip.tsx         // "Lv: X [Name]" tooltip on hover
      SelectionRing.tsx        // Glow ring on selected building
      UpgradeProgressBar3D.tsx // Yellow progress bar above upgrading building
      UpgradeParticles.tsx     // Spark particles
      Skybox.tsx               // Starfield background
      CameraController.tsx     // Pan-based camera (fixed isometric angle)
      SceneLighting.tsx        // Lights setup
      PlaceholderBuilding.tsx  // Colored box placeholder
    panels/
      BuildingContextMenu.tsx  // View/Move/Upgrade popup on click (HTML)
      BuildingDetailPanel.tsx  // Right slide-in, opened via "View" (HTML)
      ConstructionPanel.tsx    // Modal overlay for building selection (HTML)
      ConstructionInfoPanel.tsx // Bottom-right countdown panel (HTML)
    ui/
      ProgressBar.tsx
      CountdownTimer.tsx
      ResourceCost.tsx
      GlowButton.tsx
      Badge.tsx
      Tooltip.tsx
      Modal.tsx
      CategoryTabs.tsx
      IconButton.tsx
      Toast.tsx
    icons/
      BuildingIcon.tsx         // 2D icons for panels
      ResourceIcon.tsx
  contexts/
    GameContext.tsx             // Global game state
  hooks/
    useResources.ts
    useBuildings.ts
    useCountdown.ts
    useCameraPan.ts            // Pan-based camera transition logic
    usePlacementMode.ts        // Grid placement state for build/move
    useIsometricGrid.ts        // Grid coordinate <-> world position conversion
    usePerformanceTier.ts      // Auto-detect quality tier
  pages/
    Planet.tsx                 // Main planet page
    Home.tsx
  services/
    api.ts                     // (existing)
  types/
    index.ts                   // (existing)
  assets/
    models/                    // .glb files for buildings
    textures/                  // Planet surface, skybox, particles
    icons/                     // 2D building icons, resource icons
  styles/
    variables.css
    reset.css
    layout.css
    panels.css
    components.css
    animations.css
  App.tsx
  main.tsx
```

## Appendix B: Animation Specifications

**3D Animations** (Three.js):

| Animation | Duration | Easing | Trigger |
|-----------|----------|--------|---------|
| Building idle glow | 3000ms | sine | Continuous |
| Selection ring rotation | 8000ms | linear | Building selected, continuous |
| Selection ring opacity | 2000ms | sine | Building selected, continuous |
| Camera pan transition | 500ms | ease-in-out | Pan to building |
| Building spawn (scale) | 1000ms | bounce | New construction placed on grid |
| Upgrade scaffolding fade-in | 300ms | ease-out | Upgrade starts |
| Upgrade scaffolding fade-out | 300ms | ease-in | Upgrade completes |
| Upgrade yellow progress bar | continuous | linear | Fill based on time elapsed |
| Upgrade spark lifecycle | 1500ms | linear | Continuous while upgrading |
| Placement grid fade-in | 200ms | ease-out | Enter build/move mode |
| Placement grid fade-out | 200ms | ease-in | Exit build/move mode |
| Building move (reposition) | 400ms | ease-in-out | Building moved to new tile |
| Hover tooltip fade-in | 150ms | ease-out | Mouse enters building |
| Hover tooltip fade-out | 100ms | ease-in | Mouse leaves building |

**2D Animations** (CSS/HTML):

| Animation | Duration | Easing | Trigger |
|-----------|----------|--------|---------|
| Context menu appear | 150ms | ease-out | Building clicked |
| Context menu disappear | 100ms | ease-in | Option selected or click outside |
| Detail panel slide-in | 250ms | ease-out | "View" selected in context menu |
| Detail panel slide-out | 200ms | ease-in | Panel closed |
| Modal fade-in | 200ms | ease-out | Open construction panel |
| Modal fade-out | 150ms | ease-in | Close construction panel |
| Glow button hover | 200ms | ease | Mouse enter |
| Resource number change | 300ms | ease-out | Value update |
| Progress bar fill | 500ms | ease-out | Value change |
| Timer pulse (< 60s) | 1000ms | ease-in-out | Timer < 60s, infinite |
| Skeleton shimmer | 1500ms | linear | Loading state, infinite |
| Toast notification | 300ms in, 200ms out | ease | Error/success event |

## Appendix C: Interaction Summary

| User Action | 3D Response | 2D Response |
|-------------|------------|-------------|
| Hover building (3D) | Outline glow, **tooltip appears**: "Lv: X [Name]" | Cursor changes to pointer |
| Click building (3D) | Selection ring | **Context menu** appears (View/Move/Upgrade) |
| Context menu: "View" | - | Detail panel slides in from right |
| Context menu: "Move" | Grid tiles appear (green/red), building becomes semi-transparent | Context menu closes |
| Context menu: "Upgrade" | Yellow progress bar + scaffolding + sparks | Construction info panel updates with countdown |
| Context menu: "Harvest" (Warehouse) | Brief golden flash on warehouse | Resources update in HUD |
| Click "Build" button (bottom-right) | - | Construction modal opens |
| Click "Build" in modal | Wireframe grid appears + ghost preview follows cursor | Modal closes |
| Move cursor (placement mode) | Ghost snaps to grid, green/red tint + footprint highlight | - |
| Click valid position (placement mode) | New building spawns with animation + yellow progress bar | Construction info panel updates |
| Click "Collect" (2D HUD) | Brief golden flash on all resource buildings | Numbers update, pending resets |
| Click "Cancel" (2D panel) | Scaffolding + sparks + progress bar disappear | Timer removed from info panel |
| Scroll wheel (3D) | Camera zoom in/out (fixed angle) | - |
| Left-drag on empty space (3D) | **Camera pans** (scrolls viewport) | - |
| Click nav item (2D) | Scene changes (future) | Route navigates |
| Timer reaches zero | Scaffolding + progress bar disappear, model "upgrades" (glow) | Timer removed, level updates |
| Click info panel entry (2D) | Camera pans to that building | Highlight in info panel |
| ESC key | Deselect building (ring off), cancel placement mode | Close context menu/panel/modal |

## Appendix D: WebGL Fallback

If WebGL is not available (rare, but possible on old browsers):

```
+-----------------------------------------------+
|                                               |
|   Your browser does not support WebGL.        |
|                                               |
|   Cryptomines Online requires a modern        |
|   browser with WebGL support.                 |
|                                               |
|   Recommended: Chrome, Firefox, Edge          |
|                                               |
+-----------------------------------------------+
```

Display a full-screen HTML message with browser recommendations. Do not attempt a 2D fallback for Phase 1.
