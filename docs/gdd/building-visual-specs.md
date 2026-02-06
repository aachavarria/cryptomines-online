# Building Visual Specifications - Procedural Three.js Models

> All 22 building types for Cryptomines Online (Galaxy Online 2).
> Each building is built from Three.js primitives (BoxGeometry, CylinderGeometry, ConeGeometry, SphereGeometry, TorusGeometry, ExtrudeGeometry).
> Coordinate system: Y is up. Buildings sit at Y=0 on the isometric grid.

## Grid & Sizing Reference

- **TILE_WORLD_SIZE = 4** (4 world units per tile)
- A 2x2 building occupies 8x8 world units on the XZ plane
- A 3x3 building occupies 12x12 world units on the XZ plane
- Building height should generally be 4-8 world units (taller for important buildings)
- All dimensions below are in world units

## Material Convention

Established in the existing ship/building models:
- **metalness**: 0.7-0.9 for metallic surfaces
- **roughness**: 0.15-0.3 for polished metal
- **emissive**: category color at low intensity (0.1-0.2) for base glow
- **accent lights**: emissiveIntensity 1.5-2.0, toneMapped=false for bloom/glow
- **transparent panels**: opacity 0.5-0.7 for windows, energy fields, etc.

## Category Colors

| Category | Hex       | Buildings |
|----------|-----------|-----------|
| Resource | `#22cc66` | metal_collector, he3_extractor, residential_area, resource_warehouse |
| Core     | `#4488ff` | civic_center, technology_center, alliance_center, trading_center, galaxy_transporter, compound_center |
| Military | `#ff4444` | ship_factory, spacedock, command_center, weapon_research_center, radar |
| Defense  | `#44ccff` | space_station, meteor_star, particle_cannon, anti_aircraft_gun, thors_cannon, celestial_base, recycling_plant |

---

## 1. Civic Center (core, 3x3)

**Footprint**: ~10x10 world units, height ~8-9

**Shape**: Grand domed capitol building. Central hexagonal tower topped with a large hemisphere dome. Four symmetrical wing extensions radiate from the base at 90-degree intervals. A spire antenna extends from the dome apex.

**Geometry**:
- Base platform: CylinderGeometry (radius 5, height 0.6, segments 6) — hexagonal
- Central tower: CylinderGeometry (radius 2.0, height 5.0, segments 8)
- Dome: SphereGeometry (radius 2.2, half-sphere) atop the tower
- Four wings: BoxGeometry (4.0 x 1.5 x 1.5) extending outward from tower base at 0/90/180/270 degrees
- Wing caps: CylinderGeometry (radius 0.6, height 1.5, segments 6) at each wing tip
- Apex spire: CylinderGeometry (radius 0.08, height 2.0) + SphereGeometry (radius 0.15) beacon on top

**Colors**:
- Base/tower: dark navy `#1a2244`, emissive `#4488ff` at 0.15
- Dome: translucent blue `#6699ff`, opacity 0.5, emissiveIntensity 0.4
- Wings: medium blue `#2a3366`, metalness 0.85
- Accent lights: bright blue `#66bbff`, emissiveIntensity 2.0 at wing tips and dome rim

**Animation**:
- Dome pulses gently (emissiveIntensity oscillates 0.3-0.6, period ~3s)
- Apex beacon pulses brightly (emissiveIntensity 1.0-3.0, period ~1.5s)
- Four ring lights on the base platform rotate slowly (period ~10s)

---

## 2. Technology Center (core, 3x2)

**Footprint**: ~10x7 world units, height ~6-7

**Shape**: Research laboratory with satellite dish. Rectangular main building with angular roof panels. A large parabolic dish on the roof points skyward. Holographic projection ring floats above the dish.

**Geometry**:
- Main body: BoxGeometry (8.0 x 3.5 x 5.5) — long and wide
- Angular roof: two BoxGeometry (4.5 x 0.3 x 5.5) tilted +-15deg forming a V-ridge
- Satellite dish: ConeGeometry (radius 1.8, height 0.5, open-ended) inverted, on roof center
- Dish support: CylinderGeometry (radius 0.15, height 1.5)
- Holographic ring: TorusGeometry (radius 1.0, tube 0.06) floating 1.5 units above dish
- Data cores: 3x CylinderGeometry (radius 0.3, height 1.2) along one side, evenly spaced

**Colors**:
- Body: dark steel blue `#1a2a44`, emissive `#4488ff` at 0.12
- Roof panels: `#2a3a55`, metalness 0.9
- Dish: `#3a4a66`, metalness 0.95
- Holographic ring: `#66ccff`, emissive `#66ccff`, emissiveIntensity 1.2, transparent, opacity 0.5
- Data cores: `#224488`, emissive `#4488ff`, emissiveIntensity 0.6

**Animation**:
- Holographic ring rotates (Y-axis, period ~4s) and bobs up/down slightly
- Data cores pulse sequentially (emissiveIntensity stagger 0.3-0.8)
- Dish rotates very slowly on Y-axis (period ~20s)

---

## 3. Alliance Center (core, 2x2)

**Footprint**: ~7x7 world units, height ~5-6

**Shape**: Diplomatic meeting hall. Circular base with a ring of tall, thin pillars supporting a floating platform. A holographic globe hovers in the center above the platform.

**Geometry**:
- Circular base: CylinderGeometry (radius 3.2, height 0.8, segments 12)
- 6 pillars: CylinderGeometry (radius 0.15, height 3.5) arranged in a hexagonal pattern at radius 2.5
- Floating platform: CylinderGeometry (radius 2.0, height 0.3, segments 8) at Y=3.8
- Holographic globe: SphereGeometry (radius 0.8, widthSegments 12, heightSegments 8) at Y=5.0
- Globe wireframe overlay: same sphere but wireframe material
- Connecting arches between adjacent pillars: thin BoxGeometry arcs

**Colors**:
- Base: `#1a2244`, emissive `#4488ff` at 0.1
- Pillars: `#3a4466`, metalness 0.85
- Platform: `#2a3355`, emissive `#4488ff` at 0.2
- Globe: `#4488ff`, emissive `#4488ff`, emissiveIntensity 0.8, transparent, opacity 0.4
- Globe wireframe: `#88bbff`, emissiveIntensity 1.0

**Animation**:
- Globe rotates on Y-axis (period ~6s)
- Globe pulses (opacity 0.3-0.5, period ~2s)
- Platform edge lights orbit (6 small spheres at platform rim, rotating period ~8s)

---

## 4. Trading Center (core, 2x2)

**Footprint**: ~7x7 world units, height ~5

**Shape**: Marketplace/exchange building. Octagonal bazaar structure with an open-air feel. Two levels, the upper level has a retractable dome (shown partially open). Cargo containers stacked around the base.

**Geometry**:
- Octagonal base: CylinderGeometry (radius 3.0, height 2.0, segments 8)
- Upper level: CylinderGeometry (radius 2.5, height 1.5, segments 8) atop the base
- Partial dome: SphereGeometry (radius 2.5, half-sphere, phiLength ~4.5) offset to show it "opening"
- 4 cargo containers: BoxGeometry (1.0 x 0.6 x 0.6) scattered at base perimeter
- Landing pad markings: flat ring (TorusGeometry, radius 3.2, tube 0.04) on ground
- Trade beacon: CylinderGeometry (radius 0.1, height 1.5) + sphere on top of dome

**Colors**:
- Base/upper: `#1a2244`, emissive `#4488ff` at 0.12
- Dome: `#3a4a77`, transparent, opacity 0.4, emissive `#4488ff` at 0.3
- Cargo containers: mixed — `#335522`, `#553322`, `#223355` (resource-colored goods)
- Trade beacon: `#ffcc44`, emissive `#ffcc44`, emissiveIntensity 1.5
- Landing ring: `#4488ff`, emissive `#4488ff`, emissiveIntensity 0.5

**Animation**:
- Trade beacon blinks (emissiveIntensity 0.5-2.0, period ~1s)
- Landing ring pulses (emissiveIntensity 0.3-0.7, period ~3s)

---

## 5. Galaxy Transporter (core, 2x2)

**Footprint**: ~7x7 world units, height ~7

**Shape**: Teleportation gateway. Two tall vertical pylons flanking a circular energy portal. The portal ring floats between the pylons. Energy arcs connect pylons to the ring. Otherworldly, gateway aesthetic.

**Geometry**:
- Left pylon: BoxGeometry (0.8 x 6.0 x 0.8) at X=-2.5
- Right pylon: BoxGeometry (0.8 x 6.0 x 0.8) at X=+2.5
- Pylon caps: ConeGeometry (radius 0.5, height 1.0) on top of each pylon
- Portal ring: TorusGeometry (radius 2.0, tube 0.2, segments 24) — vertical, between pylons
- Portal fill: CircleGeometry (radius 1.8) inside the ring — transparent energy field
- Base platform: BoxGeometry (7.0 x 0.4 x 4.0)
- Energy conduits: 4x thin BoxGeometry arcs from pylons to ring at 45-degree angles

**Colors**:
- Pylons: `#1a2244`, emissive `#4488ff` at 0.15, metalness 0.85
- Pylon caps: `#2a3366`, emissive `#6699ff` at 0.3
- Portal ring: `#4488ff`, emissive `#4488ff`, emissiveIntensity 1.0
- Portal fill: `#88bbff`, emissive `#88bbff`, emissiveIntensity 0.6, transparent, opacity 0.3
- Conduits: `#66aaff`, emissive `#66aaff`, emissiveIntensity 0.8
- Base: `#1a1a2a`

**Animation**:
- Portal ring rotates on Z-axis (facing camera, period ~5s)
- Portal fill ripples (opacity 0.2-0.4, emissiveIntensity 0.4-0.8, period ~2s)
- Pylon cap lights pulse alternately (period ~1.5s)

---

## 6. Compound Center (core, 2x2)

**Footprint**: ~7x7 world units, height ~4.5

**Shape**: Administrative complex. A cluster of three interconnected rectangular modules of varying heights, connected by elevated walkways. Central courtyard with a holographic display.

**Geometry**:
- Module A (tallest): BoxGeometry (2.5 x 4.0 x 2.5) at offset (-1.5, 0, -1.0)
- Module B (medium): BoxGeometry (2.0 x 3.0 x 2.0) at offset (1.5, 0, -1.0)
- Module C (shortest): BoxGeometry (2.5 x 2.0 x 2.0) at offset (0, 0, 1.5)
- Walkway A-B: BoxGeometry (2.0 x 0.15 x 0.6) connecting at Y=2.5
- Walkway A-C: BoxGeometry (0.6 x 0.15 x 2.0) connecting at Y=1.8
- Courtyard hologram: flat CylinderGeometry (radius 0.5, height 2.0, open) at center — transparent beam

**Colors**:
- Modules: `#1a2244`, emissive `#4488ff` at 0.1, metalness 0.8
- Module rooftops: `#2a3366`, slightly brighter
- Walkways: `#3a4466`, metalness 0.9
- Window strips: BoxGeometry strips on module sides — `#4488ff`, emissive, emissiveIntensity 0.5, transparent, opacity 0.5
- Hologram beam: `#66bbff`, emissive, emissiveIntensity 0.8, transparent, opacity 0.3

**Animation**:
- Hologram beam rotates (Y-axis, period ~6s) and pulses (opacity 0.2-0.4)
- Window strips glow in alternating patterns (staggered pulse, period ~4s)

---

## 7. Metal Collector (resource, 2x2)

**Footprint**: ~7x7 world units, height ~4-5

**Shape**: Industrial mining rig. A central drill shaft descends into the ground, surrounded by processing conveyors. An angled arm supports a rotating excavation drum. Ore containers sit at the base. Rugged, heavy-industry feel.

**Geometry**:
- Drill shaft housing: CylinderGeometry (radius 0.8, height 3.0, segments 8) at center
- Drill head: ConeGeometry (radius 0.6, height 1.0) protruding downward (inverted, partially embedded)
- Support frame: 4x BoxGeometry (0.15 x 3.0 x 0.15) legs at corners, slight inward tilt
- Cross beams: 4x BoxGeometry (3.0 x 0.1 x 0.1) connecting legs at top
- Excavation arm: BoxGeometry (3.5 x 0.2 x 0.4) angled from top of shaft, tilted 30deg down
- Drum: CylinderGeometry (radius 0.5, height 0.8, segments 6) at arm end
- 2 ore containers: BoxGeometry (1.0 x 0.8 x 0.8) at base perimeter

**Colors**:
- Frame/legs: `#2a3a2a`, metalness 0.8, roughness 0.35 (weathered industrial)
- Drill shaft: `#1a2a1a`, emissive `#22cc66` at 0.1
- Drum: `#3a4a3a`, metalness 0.7
- Ore containers: `#445533`, emissive `#22cc66` at 0.15
- Accent lights: `#44ff88`, emissiveIntensity 1.5 — on top of frame and arm tip
- Drill head: `#667766`, metalness 0.9

**Animation**:
- Excavation drum rotates (local X-axis, period ~2s)
- Drill shaft housing rotates slowly (Y-axis, period ~8s)
- Ore container indicator lights blink (alternating, period ~1.5s)

---

## 8. He3 Extractor (resource, 2x2)

**Footprint**: ~7x7 world units, height ~5-6

**Shape**: Gas extraction tower. Tall vertical column with a bulbous collection tank at the top. Radiating pipes connect the tank to a circular ground-level processing ring. Gas venting effect from the top. More vertical and sleek than the Metal Collector.

**Geometry**:
- Central column: CylinderGeometry (radius 0.5, height 4.5, segments 8)
- Collection tank (bulb): SphereGeometry (radius 1.2) at top of column
- Tank cap: ConeGeometry (radius 0.3, height 0.6) on top of bulb
- Ground processing ring: TorusGeometry (radius 2.5, tube 0.2, segments 12) at Y=0.5
- 4 connecting pipes: CylinderGeometry (radius 0.1, height 3.0) angled from ring up to tank at 45deg
- Vent exhaust: 3x small ConeGeometry (radius 0.15, height 0.4) on top — gas nozzles
- Pressure gauges: 2x small CylinderGeometry (radius 0.2, height 0.05) on column sides

**Colors**:
- Column: `#1a2a1a`, emissive `#22cc66` at 0.12
- Tank bulb: `#2a4a2a`, transparent, opacity 0.6, emissive `#22cc66` at 0.3
- Tank inner glow: nested smaller SphereGeometry — `#44ff88`, emissive, emissiveIntensity 1.0, opacity 0.4
- Processing ring: `#2a3a2a`, metalness 0.85
- Pipes: `#3a4a3a`, metalness 0.8
- Vents: `#44ff88`, emissive, emissiveIntensity 1.2

**Animation**:
- Tank bulb inner glow pulses (emissiveIntensity 0.6-1.2, period ~2.5s) — simulates gas fill level
- Vent nozzles glow intermittently (staggered, period ~1s each)
- Processing ring has orbiting light point (period ~5s)

---

## 9. Residential Area (resource, 2x2)

**Footprint**: ~7x7 world units, height ~3.5

**Shape**: Habitat modules. A cluster of 3-4 dome-shaped habitat pods of varying sizes arranged organically, connected by short enclosed walkways at ground level. Warm interior lighting visible through windows. Most civilian-looking building.

**Geometry**:
- Large dome: SphereGeometry (radius 1.5, half-sphere) at offset (-0.5, 0, -0.3)
- Medium dome: SphereGeometry (radius 1.1, half-sphere) at offset (1.3, 0, 0.5)
- Small dome A: SphereGeometry (radius 0.8, half-sphere) at offset (0.3, 0, 1.8)
- Small dome B: SphereGeometry (radius 0.7, half-sphere) at offset (-1.5, 0, 1.2)
- Walkways: BoxGeometry tubes (0.4 x 0.4 x length) connecting adjacent domes at Y=0.3
- Windows: small BoxGeometry (0.3 x 0.3 x 0.02) inset on dome surfaces, 2-3 per dome
- Base ring per dome: TorusGeometry (matching radius, tube 0.06) at Y=0.1

**Colors**:
- Dome exteriors: `#2a3a2a`, emissive `#22cc66` at 0.08, metalness 0.6 (less metallic — habitation)
- Windows: `#ffdd88`, emissive `#ffdd88`, emissiveIntensity 0.8, transparent, opacity 0.6 (warm interior light)
- Walkways: `#3a4a3a`, metalness 0.75
- Base rings: `#22cc66`, emissive, emissiveIntensity 0.4

**Animation**:
- Windows glow warmly with subtle flicker (emissiveIntensity 0.6-1.0, slightly random per window)
- One small light orbits the large dome perimeter at ground level (period ~12s) — patrol/ambient

---

## 10. Resource Warehouse (resource, 3x2)

**Footprint**: ~10x7 world units, height ~4

**Shape**: Large storage depot. Long, low-profile warehouse with a barrel-vault roof. Loading dock doors on one long side. Stacked cargo containers visible outside. Conveyor track runs along the roof ridge. Utilitarian and blocky.

**Geometry**:
- Main structure: BoxGeometry (9.0 x 2.5 x 5.5)
- Barrel vault roof: half-CylinderGeometry (radius 3.0, height 9.0) rotated and placed on top — or approximate with a wide BoxGeometry + slight curve via ExtrudeGeometry
- Loading dock doors: 3x BoxGeometry (1.5 x 1.8 x 0.1) evenly spaced on the front face
- Cargo stacks (outside): 4-5x BoxGeometry (0.8 x 0.6 x 0.6) and (0.6 x 0.4 x 0.4) scattered near dock
- Roof conveyor: BoxGeometry (8.0 x 0.15 x 0.3) along the ridge
- Side vents: 4x BoxGeometry (0.05 x 0.3 x 1.0) on each long side

**Colors**:
- Walls: `#2a3a2a`, emissive `#22cc66` at 0.08, metalness 0.7
- Roof: `#3a4a3a`, metalness 0.8
- Dock doors: `#44aa66`, emissive `#22cc66`, emissiveIntensity 0.4, transparent, opacity 0.5
- Cargo containers: mixed greens and browns — `#446633`, `#555544`, `#336644`
- Conveyor: `#4a5a4a`, metalness 0.85
- Vents: `#22cc66`, emissive, emissiveIntensity 0.3

**Animation**:
- Loading dock doors pulse gently (emissiveIntensity 0.3-0.6, staggered) — simulating activity
- Conveyor has a small marker light that slides back and forth (period ~4s)

---

## 11. Ship Factory (military, 3x2)

**Already implemented** — see `frontend/src/components/three/ships/ShipFactoryModel.tsx`.

**Summary for reference**:
- Tall hangar body (BoxGeometry 2.6 x 3.6 x 2.2), industrial roof
- Assembly gantry tower + rotating crane arm (period ~3.3s)
- Production bay opening (front face, glowing orange, pulsing)
- Side vents (3 per side), exhaust stacks, warning lights
- Color: military red `#ff4444`, dark body `#2a1a1a`

---

## 12. Spacedock (defense, 3x2)

**Already implemented** — see `frontend/src/components/three/ships/SpacedockModel.tsx`.

**Summary for reference**:
- Central support pillar (CylinderGeometry, tapered)
- Rotating upper docking ring (TorusGeometry, 4 docking clamps with green lights)
- Static lower support ring
- Pulsing energy core (sphere in ring center)
- Support struts, shield generator dishes, ground conduits
- Color: defense cyan `#44ccff`, dark body `#1a2a3a`

---

## 13. Command Center (military, 2x2)

**Footprint**: ~7x7 world units, height ~5-6

**Shape**: Military tactical operations hub. Angular, aggressive architecture. Stepped pyramid base with a flat-topped command deck. Large holographic tactical display projected above the roof. Antenna arrays on corners. Armored aesthetic with sharp edges.

**Geometry**:
- Bottom tier: BoxGeometry (7.0 x 1.5 x 7.0) — wide armored base
- Middle tier: BoxGeometry (5.0 x 1.5 x 5.0) — stepped inward
- Command deck: BoxGeometry (3.5 x 1.2 x 3.5) — top tier
- Tactical display: flat CylinderGeometry (radius 1.0, height 0.02) hovering at Y=5.5, horizontal — holographic table
- Display projector: CylinderGeometry (radius 0.15, height 0.8) connecting deck roof to display
- 4 antenna masts: CylinderGeometry (radius 0.06, height 1.5) at the four corners of middle tier
- Viewport strip: BoxGeometry (3.2 x 0.4 x 0.05) on each face of command deck (4 total)

**Colors**:
- Base tiers: dark red-black `#2a1a1a`, emissive `#ff4444` at 0.1, metalness 0.8
- Command deck: `#3a2222`, emissive `#ff4444` at 0.15
- Tactical display: `#ff6644`, emissive, emissiveIntensity 0.8, transparent, opacity 0.4
- Viewport strips: `#ff6644`, emissive, emissiveIntensity 0.5, transparent, opacity 0.5
- Antenna masts: `#555566`, metalness 0.9
- Antenna tips: `#ff2222`, emissive, emissiveIntensity 2.0

**Animation**:
- Tactical display rotates slowly (Y-axis, period ~8s)
- Tactical display pulses (emissiveIntensity 0.5-1.0, period ~3s)
- Antenna tip lights blink in alternating pairs (period ~1s)
- Viewport strips cycle brightness (staggered wave, period ~4s)

---

## 14. Weapon Research Center (military, 2x2)

**Footprint**: ~7x7 world units, height ~5

**Shape**: Secretive weapons lab. Bunker-like lower section with heavy armor plating. Upper section features a containment sphere (weapons test chamber) held in a cradle of support arms. Hazard markings. Feels dangerous and experimental.

**Geometry**:
- Bunker base: BoxGeometry (6.5 x 2.0 x 6.5) with slight bevel — heavy, armored
- Armored ridges: 3x BoxGeometry (6.5 x 0.15 x 0.3) horizontal bands on each face
- Support cradle: 4x BoxGeometry (0.2 x 2.0 x 0.2) arms angling inward from base corners, meeting at Y=4.5
- Containment sphere: SphereGeometry (radius 1.0, segments 12) at Y=4.0 — the test chamber
- Containment ring: TorusGeometry (radius 1.2, tube 0.08) around the sphere, horizontal
- Exhaust vents: 2x BoxGeometry (0.8 x 0.3 x 0.1) on the back face
- Hazard light: SphereGeometry (radius 0.1) on the front face at mid-height

**Colors**:
- Bunker: `#2a1a1a`, emissive `#ff4444` at 0.08, metalness 0.85, roughness 0.35 (armored)
- Armor ridges: `#3a2222`, metalness 0.9
- Support arms: `#444455`, metalness 0.85
- Containment sphere: `#ff6644`, emissive `#ff4444`, emissiveIntensity 0.6, transparent, opacity 0.4
- Containment ring: `#ff8844`, emissive, emissiveIntensity 0.8
- Hazard light: `#ffaa00`, emissive, emissiveIntensity 2.0

**Animation**:
- Containment sphere pulses ominously (emissiveIntensity 0.3-0.8, period ~2s)
- Containment ring rotates (Y-axis, period ~3s)
- Hazard light blinks (on/off, period ~0.8s)
- Exhaust vents emit faint glow intermittently (period ~4s)

---

## 15. Radar (military, 1x1)

**Footprint**: ~3.5x3.5 world units, height ~5

**Shape**: Scanning tower. Slender vertical mast with a large rotating parabolic dish at the top. Small equipment box at the base. Classic radar silhouette — tall and narrow, instantly recognizable.

**Geometry**:
- Equipment base: BoxGeometry (2.0 x 1.0 x 2.0)
- Mast: CylinderGeometry (radius 0.15, height 3.5, segments 6)
- Dish mount: BoxGeometry (0.3 x 0.3 x 0.3) at mast top — pivot joint
- Parabolic dish: ConeGeometry (radius 1.5, height 0.5, openEnded) — inverted, shallow cone shape
- Dish feed horn: CylinderGeometry (radius 0.06, height 0.8) extending from dish center forward
- Feed tip: SphereGeometry (radius 0.08) at feed horn end
- Support guy-wires: 3x thin BoxGeometry (0.02 x 2.5 x 0.02) angled from base edges to mid-mast

**Colors**:
- Base box: `#2a1a1a`, emissive `#ff4444` at 0.1
- Mast: `#444455`, metalness 0.9
- Dish: `#3a2222`, emissive `#ff4444` at 0.15, metalness 0.8
- Feed tip: `#ff4444`, emissive, emissiveIntensity 1.5
- Guy-wires: `#555566`, metalness 0.85

**Animation**:
- **Dish rotates** (Y-axis, period ~4s) — the signature radar spin
- Feed tip glows with each sweep (emissiveIntensity spikes to 3.0 when facing forward, fades to 0.5)
- Base indicator light blinks (period ~2s)

---

## 16. Space Station (defense, 3x3)

**Footprint**: ~10x10 world units, height ~8-9

**Shape**: Massive orbital defense platform (ground representation). Circular multi-ring structure. Central command core surrounded by two concentric defense rings at different heights. Turret hardpoints on the outer ring. Shield generator dome on top. The largest and most imposing defense building.

**Geometry**:
- Central core: CylinderGeometry (radius 1.5, height 5.0, segments 8)
- Core dome: SphereGeometry (radius 1.6, half-sphere) on top of core
- Inner ring: TorusGeometry (radius 2.8, tube 0.4, segments 16) at Y=2.0
- Outer ring: TorusGeometry (radius 4.5, tube 0.5, segments 20) at Y=1.0
- 6 support struts: BoxGeometry (2.5 x 0.15 x 0.15) connecting inner ring to core, evenly spaced
- 6 outer struts: BoxGeometry (2.0 x 0.15 x 0.15) connecting outer to inner ring
- 8 turret hardpoints: CylinderGeometry (radius 0.2, height 0.3) + CylinderGeometry (radius 0.06, height 0.5, barrel) on outer ring
- Shield dome: SphereGeometry (radius 5.0, half-sphere), transparent — faint energy shield
- Base platform: CylinderGeometry (radius 5.0, height 0.3, segments 16)

**Colors**:
- Core: `#1a2a3a`, emissive `#44ccff` at 0.15, metalness 0.8
- Core dome: `#44ccff`, emissive, emissiveIntensity 0.4, transparent, opacity 0.3
- Rings: `#2a3a4a`, emissive `#44ccff` at 0.1, metalness 0.85
- Turrets: `#3a4a5a`, metalness 0.9
- Shield dome: `#44ccff`, emissive, emissiveIntensity 0.2, transparent, opacity 0.1
- Turret barrels: `#555566`
- Platform: `#1a2a3a`

**Animation**:
- Inner ring rotates clockwise (Y-axis, period ~12s)
- Outer ring rotates counter-clockwise (Y-axis, period ~20s)
- Shield dome pulses barely visible (opacity 0.05-0.15, period ~5s)
- Core dome energy pulses (emissiveIntensity 0.3-0.6, period ~3s)
- Turrets have blinking ready-lights (staggered, period ~2s)

---

## 17. Meteor Star (defense, 1x1)

**Footprint**: ~3.5x3.5 world units, height ~4

**Shape**: Point-defense laser turret. Low-profile armored base with a rotating turret housing a twin-barrel laser emitter. Short and compact. Looks like it could track and shoot down incoming projectiles.

**Geometry**:
- Armored base: CylinderGeometry (radius 1.5, height 1.0, segments 6) — hexagonal, squat
- Turret housing: SphereGeometry (radius 0.8, half-sphere, slightly squashed via scale) at Y=1.2
- Twin barrels: 2x CylinderGeometry (radius 0.08, height 1.5) parallel, offset +-0.2, angled 20deg upward
- Barrel tips: 2x SphereGeometry (radius 0.05) — emitter nodes
- Targeting sensor: small BoxGeometry (0.2 x 0.1 x 0.1) on top of housing
- Armor plates: 3x BoxGeometry (0.8 x 0.6 x 0.05) around base hexagon sides

**Colors**:
- Base: `#1a2a3a`, emissive `#44ccff` at 0.1, metalness 0.85
- Turret housing: `#2a3a4a`, emissive `#44ccff` at 0.15
- Barrels: `#555566`, metalness 0.9
- Barrel tips: `#44ffff`, emissive, emissiveIntensity 2.0, toneMapped false
- Sensor: `#ff4444`, emissive, emissiveIntensity 1.0
- Armor plates: `#2a3a4a`, metalness 0.8

**Animation**:
- Turret (housing + barrels) rotates on Y-axis (period ~5s) — scanning for targets
- Barrel tips glow steadily with slight flicker
- Sensor light blinks (period ~1.5s)

---

## 18. Particle Cannon (defense, 1x2)

**Footprint**: ~3.5x7 world units, height ~6

**Shape**: Heavy beam weapon emplacement. Long and narrow. Massive single barrel mounted on a tracked base. Energy capacitor bank lines one side. The barrel is thick and imposing, pointed at a slight upward angle. Most aggressive-looking defense structure.

**Geometry**:
- Tracked base: BoxGeometry (3.0 x 0.8 x 6.5) — long, low chassis
- Track wheels: 4x CylinderGeometry (radius 0.4, height 0.3) on sides, 2 per side
- Turret pivot: CylinderGeometry (radius 1.0, height 1.0, segments 8) at center of base
- Main barrel: CylinderGeometry (radius 0.4, height 4.0, tapered to 0.25 at tip) — angled 15deg upward
- Barrel shroud: CylinderGeometry (radius 0.6, height 1.0) at barrel base — heat dissipation
- Barrel tip emitter: TorusGeometry (radius 0.3, tube 0.05) at barrel end — particle focusing ring
- Capacitor bank: 3x CylinderGeometry (radius 0.25, height 1.2) in a row along one side of base

**Colors**:
- Base/chassis: `#1a2a3a`, emissive `#44ccff` at 0.08, metalness 0.8
- Turret pivot: `#2a3a4a`, metalness 0.85
- Barrel: `#3a4a5a`, metalness 0.9, roughness 0.15
- Barrel shroud: `#2a3a4a`, emissive `#44ccff` at 0.2
- Emitter ring: `#44ffff`, emissive, emissiveIntensity 1.5, toneMapped false
- Capacitors: `#2a4a5a`, emissive `#44ccff` at 0.3

**Animation**:
- Turret + barrel rotates slowly on Y-axis (period ~10s) — aiming sweep
- Emitter ring pulses rapidly (emissiveIntensity 1.0-2.5, period ~0.5s) — charging effect
- Capacitors pulse in sequence (bottom to top, period ~2s each) — energy charging
- Barrel shroud glows when capacitors are "full" (brief intensity spike every ~6s)

---

## 19. Anti-Aircraft Gun (defense, 1x1)

**Footprint**: ~3.5x3.5 world units, height ~3.5

**Shape**: Rapid-fire flak turret. Compact rotating platform with quad barrels pointing skyward. Ammo feed drum on the side. Lower and more compact than Meteor Star. Visually distinct with its quad-barrel cluster.

**Geometry**:
- Base platform: CylinderGeometry (radius 1.2, height 0.5, segments 8)
- Rotation ring: TorusGeometry (radius 1.0, tube 0.08) at Y=0.5 — visible bearing
- Turret body: BoxGeometry (1.0 x 0.8 x 1.0) at Y=1.0
- Quad barrels: 4x CylinderGeometry (radius 0.06, height 2.0) in a 2x2 grid pattern (offset +-0.15), angled 60deg upward
- Barrel cluster mount: CylinderGeometry (radius 0.25, height 0.3) connecting barrels to turret
- Ammo drum: CylinderGeometry (radius 0.3, height 0.5) on turret side
- Ammo feed belt: BoxGeometry (0.5 x 0.08 x 0.08) connecting drum to barrel mount

**Colors**:
- Base: `#1a2a3a`, emissive `#44ccff` at 0.08
- Turret body: `#2a3a4a`, metalness 0.85
- Barrels: `#555566`, metalness 0.95
- Barrel tips: `#ffaa44`, emissive, emissiveIntensity 1.0 (muzzle heat)
- Ammo drum: `#3a4a5a`, metalness 0.8
- Rotation ring: `#44ccff`, emissive, emissiveIntensity 0.3

**Animation**:
- Turret (body + barrels) rotates rapidly on Y-axis (period ~3s) — fast scanning
- Barrel tips flicker (simulating muzzle flash, random emissiveIntensity 0.5-2.0)
- Ammo drum has a small rotating element (period ~1s)

---

## 20. Thor's Cannon (defense, 2x2)

**Footprint**: ~7x7 world units, height ~7

**Shape**: Orbital strike weapon. Towering electromagnetic railgun. Two parallel rails extending vertically with a projectile cradle between them. Massive power coils wrap around the base. Lightning-like energy arcs between the rails. The most imposing defense weapon.

**Geometry**:
- Power base: CylinderGeometry (radius 2.5, height 1.5, segments 8) — heavy foundation
- Left rail: BoxGeometry (0.4 x 6.0 x 0.4) at X=-1.0
- Right rail: BoxGeometry (0.4 x 6.0 x 0.4) at X=+1.0
- Rail top caps: BoxGeometry (0.6 x 0.3 x 0.6) on each rail top — electrode terminals
- Projectile cradle: BoxGeometry (0.6 x 0.8 x 0.6) between rails at Y=3.0
- Power coils: 3x TorusGeometry (radius 1.5, tube 0.1) stacked at Y=0.8, 1.6, 2.4 around the rails
- Energy arc connectors: 3x thin BoxGeometry (1.6 x 0.04 x 0.04) between rails at coil heights
- Base vents: 4x BoxGeometry (0.6 x 0.2 x 0.05) on base perimeter

**Colors**:
- Base: `#1a2a3a`, emissive `#44ccff` at 0.1, metalness 0.85
- Rails: `#3a4a5a`, metalness 0.95, roughness 0.1 — highly polished
- Rail caps: `#44ccff`, emissive, emissiveIntensity 0.8
- Projectile cradle: `#2a3a4a`, emissive `#66ddff` at 0.3
- Power coils: `#44aacc`, emissive `#44ccff`, emissiveIntensity 0.5
- Energy arcs: `#88eeff`, emissive, emissiveIntensity 2.0, transparent, opacity 0.6
- Vents: `#44ccff`, emissive, emissiveIntensity 0.3

**Animation**:
- Energy arcs flicker (opacity 0.3-0.8, emissiveIntensity 1.0-3.0, rapid random — simulating lightning)
- Power coils pulse upward in sequence (bottom to top, period ~1.5s) — charging cycle
- Rail cap electrodes pulse together (emissiveIntensity 0.5-1.5, period ~2s)
- Projectile cradle bobs slightly (Y +-0.1, period ~3s) — magnetic levitation effect

---

## 21. Celestial Base (defense, 2x2)

**Footprint**: ~7x7 world units, height ~5

**Shape**: Fortified defense outpost. Star-shaped (5-pointed) base platform with watchtower at center. Shield projector ring encircles the tower at mid-height. Armored walls with embrasures. A mini-fortress on the planet surface.

**Geometry**:
- Star base: ExtrudeGeometry from a 5-pointed star Shape (outer radius 3.5, inner radius 2.0, height 1.0)
- Central tower: CylinderGeometry (radius 0.8, height 3.5, segments 6) at center
- Tower top: ConeGeometry (radius 0.9, height 0.8) — watchtower roof
- Shield projector ring: TorusGeometry (radius 1.5, tube 0.1) at Y=2.5
- 5 wall segments: BoxGeometry (0.8 x 1.5 x 0.2) at each star point, standing upright
- 5 embrasure lights: SphereGeometry (radius 0.06) at each wall segment top
- Observation window: BoxGeometry (0.4 x 0.2 x 0.02) on tower, 4 sides

**Colors**:
- Star base: `#1a2a3a`, emissive `#44ccff` at 0.1, metalness 0.8
- Tower: `#2a3a4a`, emissive `#44ccff` at 0.12
- Tower roof: `#3a4a5a`, metalness 0.85
- Shield ring: `#44ccff`, emissive, emissiveIntensity 0.8, transparent, opacity 0.4
- Walls: `#2a3a4a`, metalness 0.85
- Embrasure lights: `#44ffaa`, emissive, emissiveIntensity 1.5
- Observation windows: `#44ccff`, emissive, emissiveIntensity 0.4, transparent, opacity 0.5

**Animation**:
- Shield ring rotates (Y-axis, period ~6s)
- Shield ring pulses (emissiveIntensity 0.5-1.0, period ~3s)
- Embrasure lights blink in sequence around the star (period ~0.5s per light)
- Tower observation windows glow steadily with faint flicker

---

## 22. Recycling Plant (defense, 2x2)

**Footprint**: ~7x7 world units, height ~4

**Shape**: Debris processing facility. Funnel-shaped intake hopper on top feeding into a boxy processing core. Conveyor belts extend outward from the sides. Crushed material output chute. Industrial but cleaner than Metal Collector. Recycles combat debris.

**Geometry**:
- Processing core: BoxGeometry (4.5 x 2.5 x 4.5) — main body
- Intake hopper: ConeGeometry (radius 1.5 top, radius 0.5 bottom, height 1.5) inverted on top — funnel
- Hopper rim: TorusGeometry (radius 1.5, tube 0.1) at hopper top edge
- Left conveyor: BoxGeometry (3.0 x 0.1 x 0.6) extending from left side, slightly angled down
- Right conveyor: BoxGeometry (3.0 x 0.1 x 0.6) extending from right side, slightly angled down
- Conveyor rollers: 6x CylinderGeometry (radius 0.08, height 0.6) spaced along each conveyor
- Output chute: BoxGeometry (1.0 x 0.5 x 0.5) protruding from the back, angled down
- Processing indicator: CylinderGeometry (radius 0.15, height 0.3) on core front — status light

**Colors**:
- Core body: `#1a2a3a`, emissive `#44ccff` at 0.1, metalness 0.75
- Hopper: `#2a3a4a`, metalness 0.8
- Hopper rim: `#44ccff`, emissive, emissiveIntensity 0.4
- Conveyors: `#3a4a4a`, metalness 0.7
- Rollers: `#555566`, metalness 0.9
- Output chute: `#2a3a2a` (greenish tint — processed materials)
- Status light: `#44ff88`, emissive, emissiveIntensity 1.5 (green = active)

**Animation**:
- Hopper rim pulses (emissiveIntensity 0.3-0.6, period ~2s) — intake active
- Conveyor rollers rotate (local Z-axis, period ~1s) — material moving
- Status light alternates green/cyan (period ~3s)
- Faint glow inside hopper pulses (simulating material being processed)

---

## Implementation Notes for 3D Asset Developer

### Component Pattern
Follow the existing pattern from `ShipFactoryModel.tsx` and `SpacedockModel.tsx`:
- React functional component with `useRef` for animated parts
- `useFrame` hook for per-frame animations
- Props: `position`, `scale`, `emissiveIntensity`
- Use `castShadow` on substantial geometry, skip on tiny accent lights
- Use `toneMapped={false}` on bright emissive materials for bloom effect

### File Naming Convention
Place each model in `frontend/src/components/three/buildings/`:
- `CivicCenterModel.tsx`
- `TechnologyCenterModel.tsx`
- `MetalCollectorModel.tsx`
- etc.

### Performance Budget
- Target 50-150 triangles per building (before Three.js segment tessellation)
- Keep segment counts low: 6-8 for cylinders, 8-12 for spheres
- Use `useMemo` for static geometry that does not change
- Limit animated refs to 2-3 per building
- Avoid creating new materials in render loop — define in JSX or useMemo

### Silhouette Distinctiveness Checklist
Each building should be identifiable by silhouette alone:
- **Tall & narrow**: Radar, He3 Extractor, Thor's Cannon
- **Domed**: Civic Center, Residential Area, Space Station
- **Boxy/industrial**: Ship Factory, Resource Warehouse, Metal Collector
- **Ring/portal**: Spacedock, Galaxy Transporter
- **Turret/weapon**: Meteor Star, Anti-Aircraft Gun, Particle Cannon
- **Tower**: Command Center (stepped pyramid), Celestial Base (star + tower)
- **Spherical element**: Weapon Research Center (containment), Technology Center (dish)
- **Multi-module cluster**: Compound Center, Alliance Center (pillars + globe)
- **Funnel/hopper**: Recycling Plant
- **Open/bazaar**: Trading Center (partial dome)

### Level-of-Detail Scaling
Buildings at higher levels could scale slightly larger (1.0 + level * 0.02) and increase emissive intensity. This is handled by the parent rendering system, not individual models.
