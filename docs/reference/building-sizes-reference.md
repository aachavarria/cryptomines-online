# Building Sizes Reference
**CryptoMines Online - Building Grid Layout**

Last Updated: 2026-02-11

---

## 📐 Grid System Overview

- **Total Grid:** 20×20 tiles
- **Tile Size:** 4 world units (3D space)
- **System:** Isometric 3D (Three.js)
- **Pixels:** Dynamic - adjusts to viewport (no fixed pixels)

---

## 📊 Complete Building List with Sizes

| Building | Display Name | Size | Category | Base | Max/Planet |
|----------|--------------|------|----------|------|------------|
| **MAJOR STRUCTURES (3×3)** | | | | | |
| `civic_center` | Civic Center | **3×3** | Core | Ground | 1 |
| `space_station` | Space Station | **3×3** | Space | Space | 1 |
| **LARGE FACILITIES (3×2)** | | | | | |
| `ship_factory` | Ship Factory | **3×2** | Military | Ground | 1 |
| `resource_warehouse` | Resource Warehouse | **3×2** | Resource | Ground | 1 |
| `spacedock` | Spacedock | **3×2** | Military | Ground | 1 |
| `technology_center` | Technology Center | **3×2** | Core | Ground | 1 |
| **STANDARD BUILDINGS (2×2)** | | | | | |
| `metal_collector` | Metal Collector | **2×2** | Resource | Ground | Unlimited |
| `he3_extractor` | He3 Extractor | **2×2** | Resource | Ground | Unlimited |
| `residential_area` | Residential Area | **2×2** | Resource | Ground | Unlimited |
| `alliance_center` | Alliance Center | **2×2** | Core | Ground | 1 |
| `trading_center` | Trading Center | **2×2** | Core | Ground | 1 |
| `galaxy_transporter` | Galaxy Transporter | **2×2** | Core | Ground | 1 |
| `compound_center` | Compound Center | **2×2** | Core | Ground | 1 |
| `command_center` | Command Center | **2×2** | Military | Ground | 1 |
| `weapon_research_center` | Weapon Research Center | **2×2** | Military | Ground | 1 |
| `recycling_plant` | Recycling Plant | **2×2** | Military | Ground | 1 |
| `thors_cannon` | Thor's Cannon | **2×2** | Defense | Space | 4 |
| `celestial_base` | Celestial Base | **2×2** | Space | Space | 1 |
| **NARROW STRUCTURES (1×2)** | | | | | |
| `particle_cannon` | Particle Cannon | **1×2** | Defense | Space | 4 |
| **SMALL STRUCTURES (1×1)** | | | | | |
| `radar` | Radar | **1×1** | Core | Ground | 1 |
| `meteor_star` | Meteor Star | **1×1** | Space | Space | 4 |
| `anti_aircraft_gun` | Anti-Aircraft Gun | **1×1** | Defense | Space | 4 |

---

## 📏 Summary by Size

### 3×3 (9 tiles) - 2 buildings
- Civic Center
- Space Station

### 3×2 (6 tiles) - 4 buildings
- Ship Factory
- Resource Warehouse
- Spacedock
- Technology Center

### 2×2 (4 tiles) - 14 buildings
- Metal Collector
- He3 Extractor
- Residential Area
- Alliance Center
- Trading Center
- Galaxy Transporter
- Compound Center
- Command Center
- Weapon Research Center
- Recycling Plant
- Thor's Cannon
- Celestial Base

### 1×2 (2 tiles) - 1 building
- Particle Cannon

### 1×1 (1 tile) - 3 buildings
- Radar
- Meteor Star
- Anti-Aircraft Gun

---

## 🎮 Pixel Conversion (Approximate)

The system uses **world units** in 3D space, not fixed pixels.

**Conversion factors:**
- 1 tile = 4 world units
- Pixels vary based on:
  - Viewport size (player's screen)
  - Camera zoom level
  - Camera distance to planet

**Approximate screen space (1920×1080 @ normal zoom):**
- 1×1 tile ≈ 60-80px
- 2×2 tiles ≈ 120-160px
- 3×3 tiles ≈ 180-240px

*Note: These values are approximate and change with zoom/camera position.*

---

## 📐 Grid Capacity Analysis

**Total grid:** 20×20 = 400 tiles

**Maximum theoretical buildings (all 1×1):** 400 buildings

**Realistic capacity:**
- 2× 3×3 buildings = 18 tiles
- 4× 3×2 buildings = 24 tiles
- ~50× 2×2 buildings = 200 tiles
- ~20× 1×1 buildings = 20 tiles
- **Total:** ~262 tiles used, 138 tiles free

**Practical limit:** ~70-80 buildings per planet (including paths/spacing)

---

## 🔗 Related Files

- **Frontend config:** `frontend/src/config/buildingConfig.ts`
- **Grid system:** `frontend/src/contexts/GameContext.tsx`
- **Database schema:** `supabase/migrations/20260206005232_phase1_mvp.sql`

---

## 📝 Notes

- Buildings cannot overlap
- Some buildings require specific prerequisites
- Space buildings (space_station, meteor_star, etc.) are placed separately from ground buildings
- Resource buildings (metal_collector, he3_extractor, residential_area) can be built multiple times
- Most core/military buildings are limited to 1 per planet
