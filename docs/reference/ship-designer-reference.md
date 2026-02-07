# Ship Designer Reference: KRTools Layout Analysis & Cryptomines Adaptation

> **Source**: KRTools Ship Designer (https://krtools.deajae.co.uk/des/)
> **Purpose**: Reference for redesigning the Ship Design Editor modal
> **Audience**: GDD agent (spec updates), frontend-dev (implementation)
> **Date**: 2026-02-06

---

## 1. KRTools Reference Layout

### 1.1 Overall Structure (725x715px)

The KRTools ship designer is a single-screen tool with absolute-positioned panels.
All interaction happens in one view -- no page navigation, no tabs switching views.

```
+--+------+-------------------------------------------+-----------+
|  | Hull |  Module Type Tabs                         |           |
|  | List |  [Attack] [Defense] [Auxiliary]            |           |
|C |      |  Module Sub-Category Icons                |  Ship     |
|L |(4-5  |  [bal][dir][mis][sbw][pla]                |  Design   |
|A | hull |                                           |  Preview  |
|S | thumb|  +------------------------------------+   |           |
|S | nails|  | Module List (scrollable)           |   | Hull Name |
|  | with |  | [icon] Module Name       vol:XX    |   | Des. Name |
|T | page |  | [icon] Module Name       vol:XX    |   |           |
|A | nav) |  | [icon] Module Name       vol:XX    |   | +-------+ |
|B |      |  | (198x41px per item, 40x40 icon)    |   | |Module | |
|S |      |  +------------------------------------+   | |Grid   | |
|  |      |                                           | |20x20  | |
|F |      |  [Save Design] [Link Design] [Show Des.]  | |icons  | |
|R |      |                                           | +-------+ |
|I |      |                                           |           |
|G |71px  |              214px                         | [====] V  |
|  |      |                                           | Volume Bar|
+--+------+-------------------------------------------+-----------+
|          Design Stats (2-column grid)               | Cost      |
|  Attack: ___ | Shield: ___   Atk/Rnd: ___           | M: ___    |
|  Structure: ___ | Steering: ___  Agility: ___       | He3: ___  |
|  Storage: ___ | Stability: ___  Mobility: ___       | Au: ___   |
|  Defense: ___                                       |           |
+-----------------------------------------------------+-----------+
  100px offset                318px                      ~280px
```

### 1.2 Panel Dimensions & Positions

| Panel | Position | Size | Description |
|-------|----------|------|-------------|
| `#hullCategorySelector` | left:0, top:70px | 87px wide | Vertical tabs: Frigate/Cruiser/Battleship icons |
| `#hullSelector` | left:92px, top:52px | 71x267px | Hull thumbnail list with pager (4-5 visible) |
| `#hullPager` | below hull list | 71px wide | Previous/Next arrows for hull pagination |
| `#moduleTypeSelector` | left:241px, top:0 | 159px wide | 3 tabs (43px each): Attack/Defense/Auxiliary |
| `#moduleCategorySelector` | left:240px, top:46px | variable | Sub-category icons (24x24px inline) |
| `#moduleSelector` | left:211px, top:77px | 214x231px | Scrollable module list |
| `#shipDesign` | left:481px, top:10px | 235px wide | Design preview (right panel) |
| `#selectedModules` | inside shipDesign | 226x145px | Grid of 20x20 module icons on ship silhouette |
| `#designStats` | left:100px, top:403px | 318px wide | 2-column stat display (155px per column) |
| `#designCost` | left:445px, top:435px | ~280px | Metal / He3 / Gold cost |

### 1.3 Interaction Flow

1. **Select hull class** (left tabs) -- filters hull list
2. **Select hull** (hull thumbnails) -- loads hull stats, shows silhouette in preview
3. **Select module type** (top tabs: Attack/Defense/Auxiliary)
4. **Select module sub-category** (icon row below type tabs)
5. **Click module from list** -- adds to ship design (appears as 20x20 icon in grid)
6. **Click installed module** -- removes from design
7. **Stats update live** as modules are added/removed
8. **Volume bar** shows used/max capacity with progress bar
9. **Save/Link/Show** buttons for persistence

---

## 2. Module Category Taxonomy

### 2.1 Attack Modules (5 sub-categories)

| Sub-Category | CSS Class | Weapon Types | Range | Cooldown |
|--------------|-----------|-------------|-------|----------|
| **Ballistic** | `.ballistic` | Rapid Fire, Taskmaster, Gatling Cannon | 1-2 | 0 |
| **Directional** | `.directional` | Cluster Laser Transmitter, Magneto Pulsar | 2-5 | 1 |
| **Missile** | `.missle` (sic) | Rocket Frame, Starlight Missile Pod | 5-8 | 3 |
| **Ship-Based** | `.shipbased` | Streamliner, Golem | 6-10 | 4 |
| **Planetary** | `.planetary` | Lander Module | 1-2 | 1 |

### 2.2 Defense Modules (3 sub-categories)

| Sub-Category | CSS Class | Module Types |
|--------------|-----------|-------------|
| **Structure** | `.structure` | Atomic Framework, Ship Reinforcement Facility, Quick Reaction Armor, Reflective Plating, Energy Armor, Daedalus Control System |
| **Shield** | `.shield` | Orbital Shield, Energy Shield Booster, Particle Stun Shield, Heat Diffusion Shield, Space-Time Magnetic Shield, Detonator Shield, Shield Regenerator, EOS Phase Shift Shield |
| **Air Defense** | `.airdefense` | Anti-Aircraft Cannon, Powered Pulse Cannon, Extreme Counterattack |

### 2.3 Auxiliary Modules (3 sub-categories)

| Sub-Category | CSS Class | Module Types |
|--------------|-----------|-------------|
| **Electronic** | `.electronic` | Agility Booster, Infrared Scanner, ECM Booster, Auto Target System, Time Dilation Module |
| **Storage** | `.storage` | Station Warehouse, Nano Station Warehouse |
| **Transmission** | `.transmission` | Super Transmission Engine, Team Combat Engine, Anti-Matter Engine |

### 2.4 Sub-Category Icon Mapping

Each sub-category has a 24x24 icon. For our implementation, we use text abbreviations or custom SVG icons:

| Sub-Category | Abbreviation | Suggested Icon |
|--------------|-------------|----------------|
| Ballistic | BAL | Bullet/shell |
| Directional | DIR | Laser beam |
| Missile | MIS | Rocket |
| Ship-Based | SBW | Crosshair |
| Planetary | PLN | Globe with arrow |
| Structure | STR | Shield with rivets |
| Shield | SHD | Energy dome |
| Air Defense | AAD | Interceptor |
| Electronic | ELC | Circuit board |
| Storage | STO | Container |
| Transmission | TRN | Engine/thruster |

---

## 3. Adaptation Plan: KRTools to Cryptomines Online

### 3.1 What We KEEP (Faithful Recreation)

| Feature | KRTools | Our Version |
|---------|---------|-------------|
| Hull class tabs | 5 classes (frigate/cruiser/battleship/flagship/special) | 3 classes only (frigate/cruiser/battleship) |
| Module type tabs | Attack/Defense/Auxiliary | Same 3 types |
| Module sub-category filters | 11 sub-categories with icons | Same 11 sub-categories |
| Module list with icons | 198x41px items, 40x40 icon | Adapted sizing for our modal |
| Volume bar | HTML5 `<progress>` with counter | Same concept, styled to theme |
| Stats display | 2-column, 10 stat fields | Same stats: Attack, Shield, Atk/Rnd, Structure, Steering, Agility, Storage, Stability, Mobility, Defense |
| Cost display | Metal / He3 / Gold | Same 3 resources |
| Live stat calculation | Updates on every module add/remove | Same behavior |
| Single-screen layout | All panels visible simultaneously | Same -- no sub-pages or wizards |

### 3.2 What We ADAPT

| Feature | KRTools Original | Our Adaptation | Reason |
|---------|-----------------|----------------|--------|
| Hull preview | 2D ship silhouette background | **3D model** rendered with React Three Fiber | We have 5 GLB models; 3D is our visual identity |
| Module grid overlay | 20x20 icons over 2D silhouette | Module list beside 3D model | 3D model replaces the silhouette area |
| Visual theme | Original GO2 game UI (grey/blue panels) | Dark sci-fi (dark blues `#0a0e1a`, glowing cyan/amber borders) | Matches our existing Cryptomines theme |
| Rendering context | Standalone page | React modal via `createPortal` | Fits our existing panel/modal architecture |
| Hull thumbnails | Small image thumbnails with pager | Hull cards with name + stat summary | No hull sprites -- use text cards |
| Design save | Save/Link/Show (3 buttons) | Save only (backend persistence via API) | Link sharing not in Phase 2 scope |
| Module items | Background image sprite + 40x40 icon | CSS-styled cards with tier badge + icon placeholder | No original sprite assets |

### 3.3 What We IGNORE

| Feature | Reason |
|---------|--------|
| Flagship hull class | Not in Phase 2 scope (Section 8.2.1 GDD) |
| Special hull class | Not in Phase 2 scope |
| `#options` panel | KRTools-specific settings (sim options) |
| `#quick` panel | Quick-build presets -- not needed |
| `#shieldnegation` / `#hullnegation` | Advanced sim stats -- not in scope |
| `#chips` panel | GO2 commander chip system -- separate feature |
| Link Design button | No URL sharing in Phase 2 |
| Show Designs button | Our design list is a separate panel (`ShipDesignPanel`) |

### 3.4 What We ADD (Not in KRTools)

| Feature | Description |
|---------|-------------|
| **Blueprint gating** | Hulls/modules greyed out if player lacks activated blueprint |
| **3D ship preview** | Rotating 3D model replaces 2D silhouette |
| **Design name input** | Inline name field with validation (KRTools has separate dialog) |
| **Module placement order hints** | Visual ordering guide per GDD Section 8.2.7 |
| **Per-ship limit indicators** | Show "1/ship" badge on limited modules |

---

## 4. 3D Model Mapping

### 4.1 Available Models

| File | Path | Size | Visual Style |
|------|------|------|-------------|
| `nave1.glb` | `/assets/gbl/ships/nave1.glb` | 970 KB | Small fighter silhouette |
| `nave2.glb` | `/assets/gbl/ships/nave2.glb` | 1.0 MB | Medium fighter |
| `nave3.glb` | `/assets/gbl/ships/nave3.glb` | 1.1 MB | Medium-large vessel |
| `nave4.glb` | `/assets/gbl/ships/nave4.glb` | 2.5 MB | Large capital ship |
| `nave5.glb` | `/assets/gbl/ships/nave5.glb` | 2.1 MB | Large warship |

### 4.2 Hull Class to Model Assignment

Models are assigned by hull class. Within each class, hull lines can map to different models for variety:

| Hull Class | Primary Model | Alt. Model | Rationale |
|------------|--------------|------------|-----------|
| **Frigate** (10 lines) | `nave1.glb` | `nave2.glb` | Smallest models for smallest class |
| **Cruiser** (10 lines) | `nave2.glb` | `nave3.glb` | Medium models for mid-class |
| **Battleship** (5 lines) | `nave4.glb` | `nave5.glb` | Largest models for capital class |

**Suggested per-line assignment:**

```
Frigates:
  nave1.glb: Weikes, Valkyrie, Space Hunter, Devourer, Cybra      (5 lines)
  nave2.glb: Air Wanderer, GoGetter, Sparrow, Polymesus, Hamdar   (5 lines)

Cruisers:
  nave2.glb: Typhoon, Duke, Watchman, Wraith, Nicholas            (5 lines)
  nave3.glb: Bombardier, The Shuttler, Spinner, Encratos, Helena   (5 lines)

Battleships:
  nave4.glb: Estrella, Diaz, Palenka                              (3 lines)
  nave5.glb: Nettle, RV766-The Explorer                           (2 lines)
```

### 4.3 3D Preview Panel Specs

| Property | Value |
|----------|-------|
| Canvas size | 226x200px (replaces 226x145px silhouette area) |
| Camera | Perspective, FOV 45, distance auto-fit to model bounds |
| Rotation | Auto-rotate Y-axis (0.5 rad/s), drag to orbit |
| Lighting | Ambient (0.4 intensity) + directional (0.8, from top-right) |
| Background | Transparent (shows modal background) |
| Model color | Tinted by hull class color (frigate: `#44aaff`, cruiser: `#cc8844`, battleship: `#aa4444`) |

---

## 5. Blueprint Gating UX

### 5.1 Hull Blueprint Gating

Players must own an **activated** hull blueprint to select that hull for design.

| State | Visual Treatment | Interaction |
|-------|-----------------|-------------|
| **Owned & Activated** | Normal appearance, full opacity | Clickable, selects hull |
| **Owned but NOT Activated** | 50% opacity, amber border, "Activate" badge | Click shows "Activate this blueprint first" tooltip |
| **Not Owned** | 30% opacity, greyed out, lock icon overlay | Click shows "Blueprint required" tooltip with acquisition hints |

### 5.2 Module Blueprint Gating

Same pattern as hull gating, applied to the module list:

| State | Visual Treatment | Interaction |
|-------|-----------------|-------------|
| **Owned & Activated** | Normal module card | Click adds to design |
| **Not Owned** | Greyed out, lock icon, name visible but stats hidden | Click shows "Blueprint required" tooltip |

### 5.3 Starter Blueprints (Always Available)

Per GDD Section 8.3.5, tutorial completion grants these blueprints (always shown as owned):

**Hull Blueprints:**
- Weikes (Frigate)
- Typhoon (Cruiser)
- Estrella (Battleship)

**Module Blueprints:**
- Rapid Fire (Ballistic Attack)
- Energy Shield Booster (Shield Defense)
- Super Transmission Engine (Transmission Auxiliary)

### 5.4 Blueprint Check API

The existing `useBlueprints()` hook provides `hasActivated(blueprintId: number) => boolean`. The designer modal should call this for each hull and module when rendering the selection lists.

---

## 6. Proposed Layout for Cryptomines Implementation

### 6.1 Modal Layout (Adapted from KRTools)

Our modal targets ~800x600px (responsive down to 700x500px).

```
+------------------------------------------------------------------+
| [X]  NEW SHIP DESIGN                                              |
+--------+--------------------------+------------------------------+
| HULL   | MODULE SELECTION         | DESIGN PREVIEW               |
| CLASS  |                          |                              |
| [FRG]  | [Attack] [Def] [Aux]    | Hull: Weikes Tier I          |
| [CRS]  | [BAL][DIR][MIS][SBW][PL]| Design: MyFrigate-V1         |
| [BTL]  |                          |                              |
|        | +----------------------+ | +----------+                 |
| HULL   | | Rapid Fire III       | | |          |                 |
| LIST   | |   ATK:48-72  v:18   | | |   3D     |                 |
|        | | Taskmaster II        | | |  Model   |                 |
| Weikes | |   ATK:28-40  v:14   | | | Preview  |                 |
| Air W. | | Gatling Cannon I     | | |  (R3F)   |                 |
| Valky. | |   ATK:16-22  v:10   | | |          |                 |
| GoGett | +----------------------+ | +----------+                 |
| Sp.Hnt |                          |                              |
| Sparrw |                          | Volume: [========  ] 68/80   |
| Devour |                          |                              |
| ...    |                          | INSTALLED MODULES            |
|        |                          | Rapid Fire III   x4   72v [-]|
|        |                          | Orbital Shield   x2   12v [-]|
|        |                          |                              |
+--------+--------------------------+------------------------------+
| STATS                                              | COST        |
| Attack: 288    | Shield: 870    | Atk/Rnd: 288    | M: 12,400   |
| Structure: 770 | Steering: 0    | Agility: 0      | He3: 8,200  |
| Storage: 0     | Stability: 0   | Mobility: 1     | Au: 9,800   |
| Defense: 0                                         |             |
+----------------------------------------------------+-------------+
|                    [ SAVE DESIGN ]                                |
+------------------------------------------------------------------+
```

### 6.2 Three-Column Layout

| Column | Width | Content |
|--------|-------|---------|
| **Left** | ~100px | Hull class tabs (vertical) + hull list (scrollable) |
| **Center** | ~300px | Module type tabs + sub-category icons + module list (scrollable) |
| **Right** | ~300px | 3D preview canvas + volume bar + installed modules list |
| **Bottom** | full width | Stats grid (2 columns) + cost display + save button |

### 6.3 Key CSS Classes (New)

```
.sd-editor          -- modal body (flexbox row)
.sd-col-hull        -- left column
.sd-col-modules     -- center column
.sd-col-preview     -- right column
.sd-class-tab       -- hull class tab button (vertical)
.sd-hull-card       -- hull selection card
.sd-hull-locked     -- greyed out hull (no blueprint)
.sd-mod-type-tabs   -- attack/defense/auxiliary tab bar
.sd-mod-subcat      -- sub-category icon row
.sd-mod-item        -- module list item (clickable)
.sd-mod-locked      -- greyed out module (no blueprint)
.sd-preview-canvas  -- React Three Fiber canvas container
.sd-volume-bar      -- volume progress bar
.sd-installed-list  -- installed modules with +/- controls
.sd-stats-grid      -- bottom stats 2-column grid
.sd-cost-display    -- resource cost display
```

---

## 7. Implementation Notes

### 7.1 Current State of ShipDesignPanel.tsx

The existing implementation at `frontend/src/components/panels/ShipDesignPanel.tsx` has:
- Basic 2-column layout (hull list left, modules right)
- Hull class filter tabs (frigate/cruiser/battleship)
- Module add/remove with quantity tracking
- Volume calculation and validation
- Name validation (alphanumeric + `.` `-` `_`, max 20 chars)
- `createPortal` modal rendering
- Blueprint check via `useBlueprints()` hook (imported but **not yet used for gating**)

**Gaps to address:**
1. No module type/sub-category filtering (shows flat list)
2. No 3D preview (no Three.js canvas)
3. No blueprint gating visuals (lock icons, opacity)
4. Stats calculation is incomplete (only shield/structure/attack/he3)
5. No cost display
6. No module placement order guidance
7. Missing stats: Steering, Agility, Storage, Stability, Mobility, Defense

### 7.2 Data Dependencies

The redesign requires these data fields from the API:

**From `hull_types` table:**
- `hull_class` (frigate/cruiser/battleship)
- `base_shield`, `base_structure`, `installation_slots`
- `base_agility`, `base_movement`, `base_storage` (needed for new stats)
- `model_file` (new field -- or derived from hull class mapping above)

**From `module_types` table:**
- `category` (attack/defense/auxiliary)
- `sub_category` (ballistic/directional/missile/shipbased/planetary/structure/shield/airdefense/electronic/storage/transmission)
- `tier` (I/II/III)
- `volume`
- All stat fields from `effects_json`
- `limit_per_ship` (null = unlimited, 1 = unique)
- `metal_cost`, `he3_cost`, `gold_cost`

**From `player_blueprints` table (via `useBlueprints` hook):**
- `blueprint_type` (hull/module)
- `reference_id` (hull_type_id or module_type_id)
- `is_activated` (boolean)

### 7.3 Stat Calculation Reference

Per GDD Section 8.2.9, the full stat set for the design stats panel:

```
Attack     = SUM(module avg_damage * quantity)
Shield     = hull.base_shield + SUM(shield_bonuses)
Atk/Rnd    = Attack (same as total attack for display)
Structure  = hull.base_structure + SUM(structure_bonuses) + (atomic_framework_count * 500)
Steering   = SUM(steering_bonuses)
Agility    = hull.base_agility + SUM(agility_bonuses)
Storage    = hull.base_storage + SUM(storage_bonuses)
Stability  = 0 (base -- modified by modules if any)
Mobility   = hull.base_movement + SUM(movement_bonuses)
Defense    = hull.base_defense + SUM(defense_bonuses)

Cost:
  Metal = hull.base_metal + SUM(module.metal_cost * quantity)
  He3   = hull.base_he3 + SUM(module.he3_cost * quantity)
  Gold  = hull.base_gold + SUM(module.gold_cost * quantity)
```

### 7.4 Module Placement Order (Visual Guide)

The designer should display a numbered placement order guide. Modules should be auto-sorted in the installed list according to optimal placement order (GDD Section 8.2.7):

```
 1. Reflective Plating
 2. Engines (Team Combat, Anti-Matter, EOS Phase Shift)
 3. Electronics (Agility Booster, Infrared Scanner, ECM, Auto Target, Time Dilation)
 4. Maintenance (Ship Reinforcement, Shield Regenerator)
 5. Air Defense (Anti-Aircraft Cannon, Powered Pulse Cannon)
 6. Ship-Based Weapons
 7. Extreme Counterattack
 8. Quick Reaction Armor
 9. Daedalus Control System
10. Shields (EOS Phase Shift, Detonator, Space-Time Magnetic, Heat Diffusion, Particle Stun)
11. Energy Shield Booster
12. Energy Armor
13. Non-SBW Weapons (Ballistic, Directional, Missile)
```

Placement-independent modules (any position): Station Warehouse, Nano Station Warehouse, Atomic Framework, Orbital Shield, Super Transmission Engine.

---

## 8. Summary Checklist

For the **GDD agent**:
- [ ] Add `sub_category` enum to module_types SQL schema if not present
- [ ] Add `model_file` field to hull_types or create hull-to-model mapping table
- [ ] Document the 11 sub-categories in the module taxonomy section
- [ ] Add UI/UX spec section for the ship designer modal layout

For the **frontend-dev**:
- [ ] Restructure `DesignEditor` from 2-column to 3-column layout
- [ ] Add module type tabs (Attack/Defense/Auxiliary)
- [ ] Add module sub-category filter icons
- [ ] Integrate React Three Fiber for 3D hull preview
- [ ] Implement blueprint gating visuals (lock/grey/opacity)
- [ ] Add full stat calculation (all 10 stats + 3 costs)
- [ ] Add volume progress bar
- [ ] Auto-sort installed modules by placement order
- [ ] Add per-ship limit badges on unique modules
