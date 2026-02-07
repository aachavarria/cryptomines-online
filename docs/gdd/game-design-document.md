# Cryptomines Online - Game Design Document

> **Version**: 3.3
> **Last Updated**: 2026-02-06
> **Status**: Draft
> **Directive**: 1:1 faithful reproduction of Galaxy Online 2 mechanics. Only the game name and visual theme are original. All mechanics, naming, formulas, and systems must match GO2.

---

## Table of Contents

1. [Game Vision](#1-game-vision)
2. [Core Mechanics](#2-core-mechanics)
3. [Data Models (SQL)](#3-data-models-sql)
4. [Game Formulas](#4-game-formulas)
5. [API Endpoints](#5-api-endpoints)
6. [MVP Prioritization (Phase 1)](#6-mvp-prioritization-phase-1)
7. [User Flows](#7-user-flows)
8. [Phase 2: Ships, Fleets & Combat](#8-phase-2-ships-fleets--combat)

---

## 1. Game Vision

### 1.1 Elevator Pitch

**Cryptomines Online** is a 1:1 browser-based recreation of Galaxy Online 2, the classic sci-fi strategy game by IGG. Players build a ground base and space station, research technologies across 7 science trees, design custom modular ships, recruit and merge commanders, run PvE instances for blueprints, build fleets on a 3x3 tactical grid, and fight for Resource Bonus Planets in Corps. The game name and visual theme are original; all mechanics are faithful to GO2.

### 1.2 Core Concept

A persistent online space strategy game. Players manage two base locations (Ground Base and Space Station), gather three resources (Metal, He3, Gold), research seven science trees, design ships from modular components, recruit commanders through a gacha-card system, and join Corps to compete for territorial control of Resource Bonus Planets.

### 1.3 What Is Original vs GO2

| Aspect | Original (Cryptomines Online) | Faithful to GO2 |
|--------|-------------------------------|-----------------|
| Game name | Cryptomines Online | - |
| Visual theme / art style | Crypto-space aesthetic | - |
| UI framework | Modern React + WebSockets | - |
| Backend architecture | Go + Supabase | - |
| Authentication | Anonymous/guest play | - |
| All game mechanics | - | Faithful 1:1 copy |
| Resource names | - | Metal, He3, Gold |
| Building names | - | Civic Center, He3 Extractor, etc. |
| Formulas & scaling | - | Exact GO2 formulas |
| Combat system | - | 8-step GO2 resolution |
| Commander system | - | GO2 rarity tiers, Star Rank |
| Tech trees | - | GO2 seven science trees |
| Corp system | - | GO2 Corps with RBPs |

### 1.4 Target Audience

- **Primary**: Strategy game fans (25-40) who played GO2, OGame, Travian, or similar browser strategy games
- **Secondary**: Idle/incremental game enthusiasts looking for strategic depth
- **Tertiary**: Nostalgic GO2 players wanting to experience the game again

### 1.5 Platform & Tech Stack

- **Platform**: Web browser (desktop + mobile responsive)
- **Backend**: Go (with Air hot reload)
- **Frontend**: React / TypeScript (Vite)
- **Database**: Supabase (PostgreSQL + Auth + Realtime)
- **Real-time**: Supabase Realtime + WebSockets

---

## 2. Core Mechanics

All mechanics in this section are faithful to Galaxy Online 2 as documented in the research document.

### 2.1 Resource System

#### 2.1.1 Primary Resources

| Resource | Producer Building | Max Buildings | Relative Rate | Primary Use |
|----------|-------------------|---------------|---------------|-------------|
| **Metal** | Metal Collector | 8 | ~0.4x Gold | Ship manufacturing, building upgrades |
| **He3** (Helium-3) | He3 Extractor | 8 | ~0.5-0.6x Gold | Ship fuel, building upgrades |
| **Gold** | Residential Area | 8 | 1.0x (highest) | Research, trading, upgrades |

#### 2.1.2 Currencies

| Currency | Acquisition | Use |
|----------|-------------|-----|
| **Mall Points (MP)** | Premium purchase (real money) | Premium shop, auctions, commander cards |
| **Vouchers** | Daily from friends repairing structures (max 8/day), events | Speedups, some mall items |
| **Badges** | Restricted Instances, tasks | Web Mall purchases |
| **Honor Points** | League Matches (15-255 pts), Pirate Challenges | Honor shop purchases |
| **Champion Points** | Championship matches (max 150/day) | Championship shop |
| **Corsairs' Gold** | Events, trade-ins | Bionic Chips, special items |

#### 2.1.3 Resource Storage & Collection

- Resources accumulate in the **Resource Warehouse** and must be manually collected
- Collection recommended every **12 hours or less** to prevent overflow
- Warehouse has limited capacity; upgrading increases storage
- PvP winner receives **20% of loser's resources** (excluding warehouse contents)
- Resources can be obtained from: production, PvP raids, instances, quests, resource packs

### 2.2 Building System

GO2 has **two base locations**: a Ground Base (planet) and a Space Base (space station).

#### 2.2.1 Ground Base - Resource Buildings

| Building | Function | Max Count | Max Level |
|----------|----------|-----------|-----------|
| **Metal Collector** | Produces Metal | 8 | 24 |
| **He3 Extractor** | Produces He3 | 8 | 24 |
| **Residential Area** | Produces Gold (highest output) | 8 | 24 |
| **Resource Warehouse** | Stores all resources; increases capacity | Multiple | 24 |

#### 2.2.2 Ground Base - Core / Administrative Buildings

| Building | Function | Max Level |
|----------|----------|-----------|
| **Civic Center** | Main hub; determines max level of all other buildings | 12 |
| **Alliance Center** | Enables Corp membership and features | [NEEDS RESEARCH: max level] |
| **Trading Center** | Player-to-player trading and auctions | [NEEDS RESEARCH: max level] |
| **Galaxy Transporter** | Inter-system resource transport | [NEEDS RESEARCH: max level] |
| **Compound Center** | Commander card merging and enhancement | [NEEDS RESEARCH: max level] |
| **Technology Center** | Research facility (7 science trees) | 12 |
| **Radar** | Detection of incoming attacks | [NEEDS RESEARCH: max level] |

#### 2.2.3 Ground Base - Military Buildings

| Building | Function | Max Level |
|----------|----------|-----------|
| **Ship Factory** | Constructs ships; holds 20 designs, 5 production slots | 24 |
| **Spacedock** | Ship berthing and fleet management | [NEEDS RESEARCH: max level] |
| **Command Center** | Commander recruitment (60 max commanders at Lv71+) | 12 |
| **Weapon Research Center** | Develops weapons and modules | [NEEDS RESEARCH: max level] |
| **Recycling Plant** | Recovers resources from scrapped ships | [NEEDS RESEARCH: max level] |

#### 2.2.4 Ground Base - Landscaping / Decorative

Casino Resort, Beacon, Monument, Fountain, Library, Theater, Park, College, Hospital, Shopping Center, Statue, Santa Sculpture (morale/aesthetic bonuses).

[NEEDS RESEARCH: exact list of decorative buildings, their effects, costs]

#### 2.2.5 Space Base Buildings

| Building | Function |
|----------|----------|
| **Space Station** | Main orbital structure; must stay within 1 level of Civic Center |
| **Meteor Star** | Orbital defense structure |
| **Particle Cannon** | Energy weapon defense |
| **Anti-Aircraft Gun** | Anti-air defense |
| **Thor's Cannon** | Advanced heavy defense cannon |
| **Celestial Base** | Advanced orbital facility |

[NEEDS RESEARCH: Space Base building max levels, costs per level, prerequisites]

#### 2.2.6 Civic Center Upgrade Table (Hub Building)

| Level | Space Station Req | Build Time | Metal | He3 | Gold |
|-------|-------------------|-----------|-------|-----|------|
| 1 | 1 | 00:05:00 | 550 | 480 | 600 |
| 2 | 1 | 00:14:06 | 1,661 | 1,450 | 1,812 |
| 3 | 2 | 00:39:54 | 5,033 | 4,392 | 5,490 |
| 4 | 3 | 01:53:19 | 15,300 | 13,353 | 16,691 |
| 5 | 4 | 05:22:58 | 46,664 | 40,725 | 50,907 |
| 6 | 5 | 15:23:42 | 142,793 | 124,619 | 155,774 |
| 7 | 6 | 44:11:03 | 438,375 | 382,582 | 478,227 |
| 8 | 7 | 127:15:00 | 1,350,194 | 1,178,351 | 1,472,939 |
| 9 | 8 | 367:45:09 | 4,172,100 | 3,641,105 | 4,551,382 |
| 10 | 9 | 1066:28:57 | 12,933,509 | 11,287,426 | 14,109,283 |
| 11 | 10 | 3103:27:51 | 40,223,214 | 35,103,895 | 43,879,869 |
| 12 | 11 | 9062:06:56 | 125,496,426 | 109,524,154 | 136,905,192 |

**Scaling**: ~3x resource cost per level; ~2.9x build time per level.

#### 2.2.7 Technology Center Upgrade Table

| Level | Civic Req | Build Time | Metal | He3 | Gold | Research Time Reduction |
|-------|-----------|-----------|-------|-----|------|------------------------|
| 1 | 1 | 00:01:40 | 450 | 420 | 650 | 3% |
| 2 | 2 | 00:05:02 | 1,269 | 1,184 | 1,833 | 6% |
| 3 | 3 | 00:15:15 | 3,591 | 3,352 | 5,187 | 9% |
| 4 | 4 | 00:46:22 | 10,199 | 9,519 | 14,732 | 12% |
| 5 | 5 | 02:21:24 | 29,068 | 27,130 | 41,987 | 15% |
| 6 | 6 | 07:12:42 | 83,134 | 77,591 | 120,082 | 18% |
| 7 | 7 | 22:08:24 | 238,594 | 222,688 | 344,636 | 21% |
| 8 | 8 | 68:11:30 | 687,150 | 641,340 | 992,550 | 24% |
| 9 | 9 | 210:42:44 | 1,985,864 | 1,853,473 | 2,868,470 | 27% |
| 10 | 10 | 653:12:27 | 5,759,006 | 5,375,072 | 8,318,564 | 30% |
| 11 | 11 | 2031:28:32 | 16,758,707 | 15,641,460 | 24,207,021 | 33% |
| 12 | 12 | 6338:12:12 | 48,935,424 | 45,673,062 | 70,684,501 | 36% |

**Formula:** `EffectiveResearchTime = BaseResearchTime * (1 - TechCenterLevel * 0.03)`

Build times/costs are reduced by Construction Boost and Quality Materials techs from the Logistics Construction tree.

#### 2.2.8 Command Center Building

| Level | Civic Req | Cooldown | Build Time | Metal | He3 | Gold |
|-------|-----------|----------|-----------|-------|-----|------|
| 1 | 1 | 3:00:00 | 0:00:40 | 600 | 450 | 500 |
| 2 | 2 | 2:50:00 | 0:02:01 | 1,692 | 1,269 | 1,410 |
| 3 | 3 | 2:40:00 | 0:06:06 | 4,788 | 3,591 | 3,990 |
| 4 | 4 | 2:30:00 | 0:18:33 | 13,599 | 10,199 | 11,332 |
| 5 | 5 | 2:20:00 | 0:56:34 | 38,757 | 29,068 | 32,297 |

[NEEDS RESEARCH: Command Center levels 6-12 exact costs and cooldowns]

Max commanders at Player Level 71+: **60**

#### 2.2.9 Building Grid & Footprints

Buildings are placed on a **20x20 isometric grid** (400 tiles). Each building type occupies a multi-tile footprint (cols x rows):

| Footprint | Buildings |
|-----------|-----------|
| **3x3** | Civic Center, Space Station |
| **3x2** | Ship Factory, Resource Warehouse, Spacedock, Technology Center |
| **2x2** | Metal Collector, He3 Extractor, Residential Area, Alliance Center, Trading Center, Galaxy Transporter, Compound Center, Command Center, Weapon Research Center, Recycling Plant, Thor's Cannon, Celestial Base |
| **1x2** | Particle Cannon |
| **1x1** | Radar, Meteor Star, Anti-Aircraft Gun |

A building at grid position (col, row) occupies all tiles from (col, row) to (col + width - 1, row + height - 1). All tiles in the footprint must be empty and within bounds for valid placement. See UI/UX Specification Section 5.3 for detailed placement mechanics.

#### 2.2.10 Building System Rules

- **Construction Slots**: Limited concurrent construction (default 2); use Construction Cards for additional slots
- **Mutual Dependencies**: Civic Center <-> Space Station must stay within 1 level of each other
- **Prerequisite Chain**: Civic Center level gates most building upgrades
- **Speed Modifiers**: Logistics Construction Science research reduces build times

### 2.3 Technology / Research

#### 2.3.1 Seven Science Trees

| Science | Focus | Key Unlocks |
|---------|-------|-------------|
| **Logistics Construction** | Economy, construction speed, ship building, resource dev | Construction Boost, Quality Materials, High Yield Mining/Chemistry/Investing, Expanded Capacity, Sync Shipbuilding |
| **Planetary Defense** | Space station defense systems | Station weapon upgrades, defense modules |
| **Ballistics Science** | Ballistic (gun) weapons | Ballistic damage, crits, scatter, penetration |
| **Directional Science** | Beam/directional weapons | Beam damage, range, accuracy |
| **Missile Science** | Missile weapons | Missile damage, range, AoE, scattering |
| **Ship-Based Science** | Fighter/ship-based weapons | Fighter damage, interception |
| **Ship Defense Science** | Shields, structure, armor | Shield/structure upgrades, damage mitigation |

#### 2.3.2 Ballistics Science - Full Tech Tree

```
Ballistics (Base) [Lv 1-10]
  No prereq | +5% ballistic damage per level
  Lv1: 541 Gold, 0:01:57 | Lv10: 18,808 Gold, 39:18:21

  +-- Ballistic Malice [Lv 1-5] (req: Ballistics Lv 3)
  |     +1% critical hit rate per level
  |     Lv1: 7,558 Gold, 0:23:48 | Lv5: 31,327 Gold, 6:48:55
  |
  |     +-- Ballistic Crackdown [Lv 1-2] (req: Ballistics Lv 3)
  |     |     +10% critical damage per level
  |     |     Lv1: 29,628 Gold, 1:12:15 | Lv2: 34,303 Gold, 1:36:51
  |     |
  |     +-- Steady Control [Lv 1-5] (req: Ballistics Lv 6, Malice Lv 3)
  |           Reduces weapon space by 2-10%
  |
  |           +-- Precise Ballistics [Lv 1-5] (req: Ballistics Lv 8, Steady Lv 3)
  |                 +1-5% hit rate
  |
  +-- Shield Penetration [Lv 1] (req: Malice Lv 5, Crackdown Lv 1, Precise Lv 1)
  |     154,616 Gold, 4:15:00 | 15% shield bypass damage
  |
  +-- Depleted Uranium Bomb [Lv 1-3] (req: Crackdown Lv 2, Penetration Lv 1)
  |     +10-30% vs Neutral armor, +1-3% vs Light armor
  |
  +-- Fire Bomb Research [Lv 1-3] (req: Crackdown Lv 2, Penetration Lv 1)
  |     +10-30% vs Regen armor, +1-3% vs Light armor
  |
  +-- Improved Penetration [Lv 1-3] (req: DU Bomb Lv 1, Fire Bomb Lv 1)
  |     +1-3% vs Light armor, 3-10% shield pen chance
  |
  +-- Victory Rush [Lv 1] (req: DU Lv 3, Fire Lv 3, Imp.Pen. Lv 3)
  |     1,918,521 Gold, 62:20:00
  |     Range-dependent damage: 220%/180%/150%/120% + 5% crit rate/damage
  |
  +-- Ballistic Scattering [Lv 1-5] (req: Ballistics Lv 10, Malice Lv 5, Steady Lv 5, Precise Lv 3)
  |     5-25% scattering damage to adjacent ships
  |
  +-- Improved Ballistic Scattering [Lv 1-3] (req: Precise Lv 5, Scattering Lv 3)
  |     +8-25% scattering rate
  |
  +-- Hop Bomb Research [Lv 1-5] (req: Scattering Lv 5, Imp. Scatter Lv 2)
        3-15% chance to deal 100% weapon damage as scatter
```

#### 2.3.3 Ship Defense Science - Full Tech Tree

**Shield Branch:**
```
Ship Defense Tech (Base) [Lv 1-2]
  Lv1: 2,846 Gold, 0:54:00 | +2% base shield/structure/agility/defense, +5% stability
  Lv2: 3,295 Gold, 1:12:23 | +5% base stats, +10% stability

  +-- Shield Research [Lv 1-5] (req: Base Lv 1)
  |     +1-5% base shield
  |     Lv1: 10,392G | Lv5: 43,076G, 12:53:09
  |
  |     +-- Energy Diffusion [Lv 1-3] (req: Base Lv 2, Shield Lv 3)
  |     |     Each shield module reduces damage by 1-3
  |     |     Lv1: 46,765G, 3:00:00
  |     |
  |     +-- Penetration Resistance [Lv 1-2] (req: Shield Lv 5, Diffusion Lv 1)
  |     |     -3-7% enemy shield penetration chance
  |     |
  |     +-- Augment Shield [Lv 1-3] (req: Diffusion Lv 2, Pen.Res. Lv 1)
  |     |     +6-20% base shield
  |     |
  |     +-- Restoration [Lv 1-2] (req: Diffusion Lv 3, Aug.Shield Lv 2)
  |     |     +30-60% shield restore per round, +1-2% interception chance
  |     |
  |     +-- Augment Absorption [Lv 1-2] (req: Shield Lv 5, Diffusion Lv 1)
  |     |     Shield modules reduce damage by 2-5
  |     |
  |     +-- Energy Conservation [Lv 1-3] (req: Diffusion Lv 2, Aug.Abs. Lv 1)
  |     |     +3-10% chance absorb damage without He3
  |     |
  |     +-- Electronic Barrier [Lv 1-2] (req: Diffusion Lv 3, E.Cons. Lv 2)
  |           Reflect 5-10% damage before shields drop to 0
  |
  +-- Damage Mitigation [Lv 1-3] (req: Aug.Shield Lv 3, Restoration Lv 2, E.Cons. Lv 3, E.Barrier Lv 2)
        10-30% absorb double damage, 15-45% lower collateral
        Lv3: 1,426,024G, 94:06:23
```

**Structure Branch:**
```
  +-- Ship Structural Analysis [Lv 1-2] (req: Base Lv 1)
  |     +1-2% base structure
  |
  |     +-- Ship Reinforcement [Lv 1-3] (req: Base Lv 2, Analysis Lv 3)
  |     |     Each structure module reduces damage by 1-3
  |     |
  |     +-- Resilience [Lv 1-2] (req: Analysis Lv 5, Reinforce Lv 1)
  |     |     -3-7% enemy structure penetration chance
  |     |
  |     +-- Structure Improvement [Lv 1-3] (req: Reinforce Lv 2, Resilience Lv 1)
  |     |     +6-20% base structure
  |     |
  |     +-- Fast Repair [Lv 1-2] (req: Reinforce Lv 3, Struct.Imp. Lv 2)
  |     |     +30-60% structure restore per round
  |     |
  |     +-- Reaction Armor Improvement [Lv 1-2] (req: Analysis Lv 5, Reinforce Lv 1)
  |     |     Structure modules reduce damage by 2-5
  |     |
  |     +-- Defense Improvement [Lv 1-3] (req: Reinforce Lv 2, Reaction Lv 1)
  |     |     +3-10% absorb damage without He3
  |     |
  |     +-- Reflection Mastery [Lv 1-2] (req: Reinforce Lv 3, Def.Imp. Lv 2)
  |           Reflect 5-10% damage before structure drops to 0
  |
  +-- Stability Mastery [Lv 1-3] (req: Struct.Imp Lv 3, Fast Repair Lv 2, Def.Imp Lv 3, Reflection Lv 2)
        10-30% absorb double damage, 15-45% lower collateral
```

#### 2.3.4 Logistics Construction Science - Full Tech Tree

The most important tree for new players. Covers economy, construction, shipbuilding, and resource production. **Total techs: 11 | All costs in Gold only.**

```
Concurrent Construction [Lv 1] (no prereq)
  Adds 1 construction slot | 1,000G, 0:00:20

  +-- Construction Boost [Lv 1-10] (req: Concurrent Construction Lv 1)
  |     +1-15% building construction speed
  |     Lv1: 2,400G, 0:08:00 | Lv10: 1,463,112G, 56:55:01
  |
  |     +-- Quality Materials [Lv 1-10] (req: Construction Boost Lv 3)
  |           -1-15% building resource costs
  |           Lv1: 1,200G, 0:04:00 | Lv10: 512,250G, 28:27:30

Ship Building Boost [Lv 1-10] (no prereq)
  +1-15% shipbuilding speed
  Lv1: 990G, 0:03:18 | Lv10: 1,463,112G, 81:17:02

  +-- Ship Building Logistics [Lv 1-10] (req: Ship Building Boost Lv 2)
  |     -1-15% ship construction resource costs
  |     Lv1: 1,386G, 0:04:37 | Lv10: 2,048,358G, 113:47:52
  |
  |     +-- Sync Shipbuilding [Lv 1] (req: Ship Building Logistics Lv 4)
  |           Adds 1 shipbuilding slot (5th production slot)
  |           174,000G, 9:40:00
  |
  |           +-- Repair Technology [Lv 1-10] (req: Sync Shipbuilding Lv 1)
  |                 +1-10% ship repair percentage
  |                 Lv1: 5,310G, 0:17:42 | Lv10: 2,266,722G, 125:55:45

High Yield Mining [Lv 1-10] (no prereq)
  +1-10% Metal output
  Lv1: 1,740G, 0:05:48 | Lv10: 742,764G, 41:15:53

  +-- High Yield Chemistry [Lv 1-10] (req: High Yield Mining Lv 2)
  |     +1-10% He3 output
  |     Lv1: 2,400G, 0:08:00 | Lv10: 1,024,500G, 56:55:00
  |
  |     +-- High Yield Investing [Lv 1-10] (req: High Yield Chemistry Lv 2)
  |           +1-10% Gold output
  |           Lv1: 3,570G, 0:11:54 | Lv10: 1,523,952G, 84:39:50
  |
  |           +-- Expand Capacity [Lv 1-7+] (req: High Yield Investing Lv 4)
  |                 +50,000-350,000 warehouse storage per level
  |                 Lv1: 3,540G, 0:11:48
```

**Summary Table:**

| Tech | Max Lv | Prerequisites | Effect per Level | Lv1 Cost | Max Lv Cost |
|------|--------|---------------|-----------------|----------|-------------|
| Concurrent Construction | 1 | None | +1 construction slot | 1,000G | - |
| Construction Boost | 10 | Concurrent Lv1 | +1-15% build speed | 2,400G | 1,463,112G |
| Quality Materials | 10 | Const.Boost Lv3 | -1-15% build costs | 1,200G | 512,250G |
| Ship Building Boost | 10 | None | +1-15% ship build speed | 990G | 1,463,112G |
| Ship Building Logistics | 10 | SBB Lv2 | -1-15% ship build costs | 1,386G | 2,048,358G |
| Sync Shipbuilding | 1 | SBL Lv4 | +1 ship production slot | 174,000G | - |
| Repair Technology | 10 | Sync Lv1 | +1-10% repair % | 5,310G | 2,266,722G |
| High Yield Mining | 10 | None | +1-10% Metal output | 1,740G | 742,764G |
| High Yield Chemistry | 10 | HYM Lv2 | +1-10% He3 output | 2,400G | 1,024,500G |
| High Yield Investing | 10 | HYC Lv2 | +1-10% Gold output | 3,570G | 1,523,952G |
| Expand Capacity | 7+ | HYI Lv4 | +50k-350k storage | 3,540G | TBD |

#### 2.3.5 Directional Science - Full Tech Tree

Enhances beam/directional weapons. Range 2-5 (extendable to 2-6). Lowest He3 usage. Piercing damage hits all ships vertically. **Total techs: 15 | All costs in Gold only.**

```
Optics (Base) [Lv 1-10] (no prereq)
  +5% directional weapon damage per level (max 50%)
  Lv1: 541G, 0:01:57 | Lv10: 18,808G, 39:18:21

  +-- Directional Malice [Lv 1-5] (req: Optics Lv 3)
  |     +1% critical strike rate per level (max 5%)
  |     Lv1: 7,558G, 0:23:48 | Lv5: 31,327G, 6:48:55
  |
  +-- Directional Accuracy [Lv 1-5] (req: Optics Lv 3)
  |     +1% accuracy per level (max 5%)
  |     Lv1: 7,558G, 0:23:48 | Lv5: 31,327G, 6:48:55
  |
  +-- Eagle Eye [Lv 1-2] (req: Optics Lv 6, Dir.Accuracy Lv 1)
  |     +10% steering power per level (max 20%)
  |     Lv1: 22,724G, 0:42:30 | Lv2: 26,309G, 0:56:58
  |
  +-- Energy Penetration [Lv 1] (req: Dir.Accuracy Lv 2, Eagle Eye Lv 1)
  |     +8% hit rate AND 8% shield penetration chance
  |     62,584G, 1:42:00
  |
  +-- Pierce [Lv 1-5] (req: Optics Lv 10, Malice Lv 5, Accuracy Lv 5, E.Pen Lv 1)
  |     Piercing damage through target rows, +3% per level
  |     Lv1: 87,919G, 1:22:27 | Lv5: 364,426G, 23:36:36
  |
  +-- Radiative Interference [Lv 1-5] (req: Optics Lv 6, Accuracy Lv 2)
  |     Reduces enemy hit rate by 2% per level (stackable, max 25%)
  |     Lv1: 43,116G, 2:33:00 | Lv5: 178,716G, 43:48:45
  |
  +-- Improved Pierce [Lv 1-2] (req: Pierce Lv 3, Eagle Eye Lv 2)
  |     +5% additional piercing damage per level
  |     Lv1: 240,487G, 4:57:30 | Lv2: 278,425G, 6:38:46
  |
  +-- Energy Accumulation [Lv 1-3] (req: Eagle Eye Lv 2, E.Pen Lv 1)
  |     +2-6% critical bonus when fleet moves
  |     Lv1: 136,845G, 3:19:45 | Lv3: 242,368G, 10:26:35
  |
  +-- Electronic Interference [Lv 1-3] (req: Accuracy Lv 2, Rad.Interf. Lv 3)
  |     5% chance to reduce enemy steering by 10%; -3-9% weapon space
  |     Lv1: 87,581G, 3:19:45 | Lv3: 155,115G, 10:26:35
  |
  +-- Piercing Crit [Lv 1] (req: Pierce Lv 5, Imp.Pierce Lv 1)
  |     Enables critical piercing vs horizontally-aligned ships
  |     391,234G, 7:05:00
  |
  +-- Weakness Detection [Lv 1-3] (req: E.Accum Lv 3, Elec.Interf. Lv 3)
  |     +5% base damage per level; 3-10% bonus accuracy
  |     Lv1: 349,054G, 5:38:18 | Lv3: 618,216G, 17:41:12
  |
  +-- Particle Impact Tech [Lv 1-3] (req: Imp.Pierce Lv 2, Piercing Crit Lv 1)
  |     Reduces enemy attack power by 10-30% based on piercing damage, lasts 2 rounds
  |     Lv1: 646,368G, 11:20:00 | Lv3: 1,144,795G, 35:33:04
  |
  +-- Magnetic Impact [Lv 1-3] (req: Accuracy Lv 5, Rad.Interf. Lv 5, Weakness Lv 2)
  |     Ignores 8-25% enemy agility with 10-30% chance
  |     Lv1: 996,120G, 26:55:00 | Lv3: 1,764,254G, 84:26:04
  |
  +-- Dynamic Impairment [Lv 1-3] (req: Particle Lv 3, Weakness Lv 3, Magnetic Lv 3)
        Ignores 8-25% defense; reduces movement 1-3 with 6-20% chance
        Lv1: 1,431,605G, 34:42:30 | Lv3: 2,535,549G, 108:52:33
```

#### 2.3.6 Missile Science - Full Tech Tree

Enhances missile/guided weapons. Scatters damage across ALL enemy ships. Range 3-6. **Total techs: 16 | All costs in Gold only.**

```
Missile Theory (Base) [Lv 1-10] (no prereq)
  +4% missile damage per level (max 40%)
  Lv1: 541G, 0:01:57 | Lv10: 18,808G, 39:18:21

  +-- Missile Accuracy [Lv 1-5] (req: Missile Theory Lv 3)
  |     +2% hit rate per level (max 10%)
  |     Lv1: 7,558G, 0:23:48 | Lv5: 31,327G, 6:48:55
  |
  +-- Cruise Dynamics [Lv 1-2] (req: Missile Theory Lv 3)
  |     +10% steering power per level (max 20%)
  |     Lv1: 11,603G, 0:56:06 | Lv2: 13,434G, 1:15:12
  |
  +-- Missile Research [Lv 1-3] (req: Theory Lv 6, Cruise Lv 1)
  |     +3-10% base damage, +1-5% shield pen chance
  |     Lv1: 38,831G, 2:04:06 | Lv3: 68,773G, 6:29:17
  |
  +-- Missile Elusion [Lv 1-3] (req: Cruise Lv 2, Missile Research Lv 1)
  |     -3-10% interception rate, +3-9% hit rate
  |     Lv1: 78,750G, 2:41:30 | Lv3: 139,478G, 8:26:37
  |
  +-- Missile Space Optimization [Lv 1-4] (req: Research Lv 2, Elusion Lv 2)
  |     -5-20% weapon space per level
  |     Lv1: 233,515G, 4:40:30 | Lv4: 632,710G, 34:19:16
  |
  +-- Multidirectional Assault [Lv 1-5] (req: Theory Lv 10, Accuracy Lv 3, Cruise Lv 2, Research Lv 3)
  |     Scatters 6-30% damage across ALL enemy ships
  |     Lv1: 80,796G, 2:50:00 | Lv5: 334,903G, 48:40:50
  |
  +-- Nuclear Radiation Research [Lv 1-5] (req: Accuracy Lv 4, Multidirect. Lv 1)
  |     2-10% increase to target fleet's damage taken, 4-20% chance
  |     Lv1: 112,916G, 2:16:00 | Lv5: 468,041G, 38:56:40
  |
  +-- Break Armor [Lv 1] (req: Multidirect. Lv 2, Nuc.Rad. Lv 3)
  |     +5% damage to all armor types
  |     287,437G, 7:05:00
  |
  +-- Energy Conservation [Lv 1-4] (req: Multidirect. Lv 3, Nuc.Rad. Lv 5, Break Armor Lv 1)
  |     4-18% chance to reduce He3 cost by 50%
  |     Lv1: 428,575G, 8:30:00 | Lv4: 1,161,222G, 62:24:06
  |
  +-- Shrapnel Research [Lv 1-2] (req: Accuracy Lv 5, Multidirect. Lv 3)
  |     +2-4% more scattering damage
  |     Lv1: 157,172G, 4:23:30 | Lv2: 181,964G, 5:53:11
  |
  +-- Exaltation [Lv 1-4] (req: Multidirect. Lv 5, Shrapnel Lv 1)
  |     +20-80% scatter vs lower-structure fleets, 2-8% He3 reduction chance
  |     Lv1: 224,495G, 4:19:15 | Lv4: 608,270G, 31:43:15
  |
  +-- Suppression [Lv 1-4] (req: Shrapnel Lv 2, Exaltation Lv 4)
  |     +3-12% scatter rate vs higher-structure fleets, +3-12% crit chance
  |     Lv1: 484,776G, 6:22:30 | Lv4: 1,313,502G, 46:48:05
  |
  +-- Missile Concussion [Lv 1] (req: E.Conservation Lv 2, Research Lv 3, Elusion Lv 3, Space Opt. Lv 4)
  |     25% chance to push target fleet back 2 spaces (once/round)
  |     1,351,976G, 49:35:00
  |
  +-- Rapid Loading [Lv ?] (wiki details truncated)
  |     Reduces reload time by one round
  |
  +-- Perfect Storm [Lv ?] (wiki details truncated)
        Capstone tech
```

#### 2.3.7 Ship-Based Science - Full Tech Tree

Enhances fighter-based weapons. Fighters have unique interception mechanics. **Total techs: 16 | All costs in Gold only.**

```
Fighter Weapons Theory (Base) [Lv 1-10] (no prereq)
  +3% fighter weapon damage per level (max 30%)
  Lv1: 541G, 0:01:57 | Lv10: 18,808G, 39:18:21

  +-- Reconnaissance [Lv 1-2] (req: FWT Lv 3)
  |     +10% steering power per level (max 20%)
  |
  +-- Thruster Optimization [Lv 1-5] (req: FWT Lv 3)
  |     -1% intercept rate per level (max -5%)
  |
  +-- Navigation [Lv 1-5] (req: FWT Lv 3)
  |     +1% hit rate per level (max 5%)
  |
  +-- Fuel Optimization [Lv 1-5] (req: FWT Lv 6, Reconnaissance Lv 1)
  |     -1-5% He3 costs
  |
  +-- Fighter Tech Upgrades [Lv 1-3] (req: Reconnaissance Lv 2, Fuel Opt. Lv 5)
  |     -2-6% He3 costs; +3-10% damage vs shielded enemies
  |
  +-- Fighter Mastery [Lv 1] (req: FWT Lv 6, Navigation Lv 3)
  |     +5% base attack power
  |
  +-- Armor Structural Analysis [Lv 1] (req: Fighter Tech Upgrades Lv 3)
  |     +10% damage vs ships and shields
  |
  +-- Fighter-based Weapons Efficiency [Lv 1-3] (req: FWT Lv 10, Armor Analysis Lv 1)
  |     +10% chance per level to finish reloading after attacks (max 30%)
  |
  +-- Long-ranged Strike [Lv 1] (req: Fuel Opt. Lv 3, Thruster Lv 3, Mastery Lv 1)
  |     +3-15% attack at distances 6-10 slots
  |
  +-- Fighter Interception Countermeasures [Lv 1-5] (req: Thruster Lv 3, Nav Lv 5, Mastery Lv 1)
  |     -1-5% intercept rate
  |
  +-- Formation Optimization [Lv 1-2] (req: FIC Lv 5)
  |     +5-10% critical damage; -5-10% occupied space
  |
  +-- Swarm [Lv 1-3] (req: FWT Lv 10, Formation Opt. Lv 2)
  |     10-30% chance: +8-25% attack with 15-35% He3 surcharge
  |
  +-- Fortune [Lv 1-3] (req: Thruster Lv 5, Long-ranged Lv 1, Swarm Lv 3)
  |     10-30% double-damage in long-range; 10-30% crit boost
  |
  +-- Ingenuity [Lv 1] (req: Weapons Efficiency Lv 3, Fortune Lv 3)
  |     +5% swarm trigger, +5% attack, -5% He3, +20% intercept rate, +10% unshielded damage
  |     2,045,150G, 70:50:00
  |
  +-- Heavy Gear Research [Lv 1-2] (req: Ingenuity Lv 1)
        Enhanced reloading, shield damage, attack range, interception
        Lv1: 10,000,000G
```

#### 2.3.8 Planetary Defense Science - Full Tech Tree

Enhances space station defense structures. **Total techs: 8 | All costs in Gold only.**

```
Energy Control (Base) [Lv 1-10] (no prereq)
  -1-10% resource costs for defensive structures
  Lv1: 2,500G, 0:08:20 | Lv10: 1,280,000G, 71:06:40

  +-- Rapid Defense Buildup [Lv 1-10] (req: Energy Control Lv 1)
  |     +1-10% defense construction speed
  |     Lv1: 3,000G, 0:10:00 | Lv10: 1,536,000G, 85:20:00
  |
  +-- Defense Enhancement [Lv 1-10] (req: Energy Control Lv 3, Rapid Lv 3)
  |     +1-10% defensive value of all structures
  |     Lv1: 1,500G, 0:05:00 | Lv10: 768,000G, 42:40:00
  |
  +-- Emplacement Mastery [Lv 1-10] (req: Rapid Lv 5)
  |     +1-10% emplacement attack power
  |     Lv1: 2,250G, 0:07:30 | Lv10: 1,152,000G, 64:00:00
  |
  +-- Utmost Defense Buildup [Lv 1-10] (req: Energy Control Lv 5, Defense Enh. Lv 5, Emplacement Lv 3)
  |     +1-10% max number of defensive structures
  |     Lv1: 1,750G, 0:05:50 | Lv10: 896,000G, 49:46:40
  |
  +-- Range Extension [Lv 1-2] (req: Rapid Lv 8, Emplacement Lv 5)
  |     Increases attack range of Particle Cannons and Anti-Aircraft Guns
  |     Lv1: 300,000G, 16:40:00 | Lv2: 600,000G, 33:20:00
  |
  +-- Thor Buildup [Lv 1] (req: Energy Control Lv 8, Utmost Lv 5, Range Ext. Lv 1)
  |     +1 max Thor's Cannon allowed
  |     800,000G, 44:26:40
  |
  +-- Augment Propulsion [Lv 1-2] (req: Emplacement Lv 8, Thor Buildup Lv 1)
        +1-2 ship movement speed when defending own planet only
        Lv1: 400,000G, 22:13:20 | Lv2: 900,000G, 50:00:00
```

**Planetary Defense Cost Reference (Energy Control & Rapid Defense per level):**

| Tech | Lv | Gold Cost | Research Time |
|------|----|-----------|---------------|
| Energy Control | 1 | 2,500 | 00:08:20 |
| Energy Control | 2 | 5,000 | 00:16:40 |
| Energy Control | 3 | 10,000 | 00:33:20 |
| Energy Control | 4 | 20,000 | 01:06:40 |
| Energy Control | 5 | 40,000 | 02:13:20 |
| Energy Control | 6 | 80,000 | 04:26:40 |
| Energy Control | 7 | 160,000 | 08:53:20 |
| Energy Control | 8 | 320,000 | 17:46:40 |
| Energy Control | 9 | 640,000 | 35:33:20 |
| Energy Control | 10 | 1,280,000 | 71:06:40 |
| Rapid Defense | 1 | 3,000 | 00:10:00 |
| Rapid Defense | 2 | 6,000 | 00:20:00 |
| Rapid Defense | 3 | 12,000 | 00:40:00 |
| Rapid Defense | 4 | 24,000 | 01:20:00 |
| Rapid Defense | 5 | 48,000 | 02:40:00 |
| Rapid Defense | 6 | 96,000 | 05:20:00 |
| Rapid Defense | 7 | 192,000 | 10:40:00 |
| Rapid Defense | 8 | 384,000 | 21:20:00 |
| Rapid Defense | 9 | 768,000 | 42:40:00 |
| Rapid Defense | 10 | 1,536,000 | 85:20:00 |
| Range Extension | 1 | 300,000 | 16:40:00 |
| Range Extension | 2 | 600,000 | 33:20:00 |
| Thor Buildup | 1 | 800,000 | 44:26:40 |
| Augment Propulsion | 1 | 400,000 | 22:13:20 |
| Augment Propulsion | 2 | 900,000 | 50:00:00 |

#### 2.3.9 Research Rules

- Only one research can be active per tree at a time (7 trees = up to 7 concurrent researches)
- **All tech research costs Gold only** -- no Metal or He3
- Technology Center level reduces research time by 3% per level (max 36% at Lv 12)
- Research costs increase exponentially per tech level (~1.53x cost per level, ~2.34x time per level)
- Technologies have prerequisite chains within their tree

#### 2.3.10 Weapon Research Center

Separate building from the Technology Center. Handles **blueprint research** for ship hulls and modules.

**Grid Size:** 2x2

| Level | Civic Req | Build Time | Metal | He3 | Gold | Research Time Reduction |
|-------|-----------|-----------|-------|------|------|------------------------|
| 1 | 1 | 0:00:40 | 500 | 300 | 450 | 3% |
| 6 | 6 | 2:53:05 | 92,371 | 55,422 | 83,134 | 18% |
| 12 | 12 | 2535:16:53 | 54,372,693 | 32,623,616 | 48,935,424 | 36% |

Time reduction formula: identical to Technology Center (3% per level).

**Blueprint Research Levels:**

| Research Level | Requirement | Effect | Cost Multiplier |
|----------------|-------------|--------|-----------------|
| 1 (Base) | Blueprint acquired | Base stats | 1.0x |
| 2 | WRC Level 6 | +10% stats | 2.0x |
| 3 | WRC Level 10 | +25% stats | 5.0x |

#### 2.3.11 Research Acceleration (Speedups)

| Research Type | Cost per 30 min Speedup | Notes |
|---------------|------------------------|-------|
| Technology Research | 3 vouchers/MP | Standard rate |
| Blueprint Research | 8 vouchers/MP | Most expensive; highest priority for speedups |
| Building Construction | 3 vouchers/MP | Friends provide 2% free acceleration per building |
| Instant Completion | 1 voucher per 10 min remaining | Alternative to incremental speedup |

### 2.4 Ship & Fleet System

#### 2.4.1 Hull Types (Rock-Paper-Scissors)

| Hull | Effective Stack | Bonus vs | Penalty vs | Armor Types |
|------|----------------|----------|------------|-------------|
| **Frigate (F)** | 1,100 | Battleship (+5%) | Cruiser (-5%) | Nano, Chrome, Regen, Neutralizing |
| **Cruiser (C)** | 1,000 | Frigate (+5%) | Battleship (-5%) | Nano, Chrome, Regen, Neutralizing |
| **Battleship (B)** | 900 | Cruiser (+5%) | Frigate (-5%) | Nano, Chrome, Regen, Neutralizing |

#### 2.4.2 Ship Modules

**Attack Modules:**

| Category | Weapon Type | Range | Cooldown | Notes |
|----------|-------------|-------|----------|-------|
| Ballistic | Guns/Cannons | Short | 0 (fires every round) | Best for consistent DPS |
| Directional | Beams/Lasers | Medium | 1-2 rounds | Pierce shields |
| Missile | Missiles | Long (min 5) | 3-4 rounds | High burst, interceptable |
| Ship-Based | Fighters | Max (min 6) | High | Interceptable, high range |
| Planetary | Siege weapons | - | - | Only damages structures |

**Defense Modules:**
- **Shield Modules**: EOS Phase Shift, Heat Diffusion, Daedalus Control System, Energy Armor
- **Structure Modules**: Hull plating, reinforcement
- **Air Defense**: Powered Pulse Cannon (PPC) - 55% intercept chance per incoming attack

**Auxiliary Modules:**
- **Electronic**: Targeting, sensors
- **Storage**: Resource/ammo capacity
- **Transmission**: Movement speed (TCE, AME, Super Transmission Engine)

[NEEDS RESEARCH: complete module list with stats, space requirements, and costs]

#### 2.4.3 Ship Stats

| Stat | Function |
|------|----------|
| **Shield** | Absorbs damage before hull (e.g., Weikes-I Frigate: 270) |
| **Structure** | Hull HP (e.g., Weikes-I: 770, Typhoon-I Cruiser: 1,040) |
| **Stability** | Reduces effective damage to hull (harder to destroy) |
| **Defense** | Damage reduction percentage |
| **Installation Slots** | Module capacity (e.g., Frigates ~100, Cruisers ~140) |
| **Agility** | Dodge chance, reduces enemy hit rate |
| **Movement (MOV)** | Fleet travel speed |
| **Storage** | Cargo capacity |

[NEEDS RESEARCH: complete ship hull table with base stats for all hull types]

#### 2.4.4 Fleet Composition

- **Fleet Grid**: 3x3 grid = 9 stacks maximum
- **Stack Size**: 3,000 ships per stack (max 27,000 ships per fleet)
- **Rule**: One ship design per stack, do NOT mix weapon types within a fleet

**Grid Positions & Attack Power:**

```
+------+----------+----------+
| Head | Shoulder | Shoulder |   First Rank:  100% attack
+------+----------+----------+
| Flank| Glasshouse| Flank   |   Second Rank:  90% attack
+------+----------+----------+
| Rear |   Tail   |  Rear    |   Third Rank:  75% attack
+------+----------+----------+
```

- **Glasshouse** (center-middle): Most protected position
- **Shoulders**: Most vulnerable positions
- **Fleet speed** = speed of the slowest ship

#### 2.4.5 Fleet Formations

| Formation | Description |
|-----------|-------------|
| **Phalanx** | All 9 stacks filled; maximum firepower |
| **Diamond** | 5 stacks (center + 4 adjacent); protects central stack |
| **Battle Line** | 6 stacks (first 2 ranks); frontal defense |
| **Skirmish** | Sparse; minimizes scatter damage |
| **Tee Forward** | T-shape; balanced offense |
| **Enfilade** | Side-focused fire |
| **Tee Reverse** | Reversed T-shape |

[NEEDS RESEARCH: exact grid positions activated for each formation]

#### 2.4.6 Targeting Commands

Fleets can be ordered to target:
- Max/Min attack power
- Max/Min durability
- Closest proximity
- By commander rank

#### 2.4.7 Blueprints

- **Required** to build any ship
- Primary source: Instances (PvE) via Treasure Boxes
- Also from: Auction House, Corp bonuses, events
- Treasure Box: **10% chance** of blueprint, split equally among available blueprints
- Avoid upgrading blueprints unless required for next level

#### 2.4.8 Ship Factory

- Max 20 ship designs stored
- 5 production slots (5th requires Sync Shipbuilding research)
- Max production: 2,000,000 ships at a time
- Speed increases 1-60% with levels 1-24

[NEEDS RESEARCH: Ship Factory speed bonus per level table]

### 2.5 Combat System

#### 2.5.1 Combat Resolution (8 Steps per Round)

**Step 1: Attacker Fires**
```
Attacks = (Attack modules per ship) x (Ships in effective stack) x (Hit chance)

Hit Chance influenced by:
  + Weapon base accuracy
  + Attacker steering stat (~4% per point)
  + Attacker accuracy stat (~1% per 12 points)
  - Defender agility (~4% per point)
  - Defender commander dodge stat
```

**Step 2: Interceptors Fire**
- **PPC (Powered Pulse Cannon)**: 55% chance to shoot down each incoming attack
- Intercept modules can counter missiles and ship-based weapons
- Defender's ship defense stat and commander dodge affect interception

**Step 3: Calculate Damage**
- Weapon damage randomly generated within min-max range
- Double-hits and critical hits applied here
- Modified by armor type (defender) vs damage type (attacker)

**Step 4: Damage Negation**
- Non-EOS shields reduce damage (Heat Diffusion, Daedalus Control, Energy Armor)
- Reduces scatter and piercing effects

**Step 5: Shield Penetration**
- Ballistic and directional weapons can pierce shields
- Penetration chance modified by tech research

**Step 6: Deal Damage to Shields**
- EOS Phase Shift: 30% chance to absorb double damage (with Damage Mitigation tech)
- Only effective stack triggers EOS defensively
- Combined shield total from all 3,000 ships absorbs damage

**Step 7: Assign Damage to Hull**
```
Ships Destroyed = floor(Remaining Damage / (Individual Ship Structure x Stability%))
```

**Step 8: Calculate Scatter Damage**
- Based on weapon type, technology, and Sandora module
- Scatter damage **cannot be absorbed or mitigated by defenses**
- Missile Exaltation tech: "higher total structure than target deal 54% scatter = 432% unpreventable bonus damage"

#### 2.5.2 Battle Duration

- Minimum: **20 rounds** + number of fleets/buildings present
- Maximum: **99 rounds**
- Only **ballistic weapons** fire every round; others require cooldown

#### 2.5.3 Ship Type Advantage Matrix

```
Frigate  ---(+5%)--->  Battleship
Cruiser  ---(+5%)--->  Frigate
Battleship ---(+5%)--> Cruiser
```

#### 2.5.4 Armor Types vs Damage Types

| Armor Type | Strong Against | Weak Against |
|------------|---------------|--------------|
| Nano | - | - |
| Chrome | Kinetic | - |
| Regen | Light Energy | - |
| Neutralizing | - | - |

Tech-specific bonuses:
- Depleted Uranium Bomb: +10-30% vs Neutral armor
- Fire Bomb: +10-30% vs Regen armor
- Both give +1-3% vs Light armor

[NEEDS RESEARCH: full armor type damage modifier matrix with exact percentages]

#### 2.5.5 Commander Impact on Combat

**Weapon Expertise:**

| Grade | Damage Modifier |
|-------|----------------|
| S | +30% |
| A | +10% |
| B | 0% (neutral) |
| C | -10% |
| D | -30% |

**Ship Expertise:**

| Grade | Damage Dealt | Damage Received |
|-------|-------------|-----------------|
| S | +10% | -10% |
| A | +5% | -10% |
| B | 0% | 0% |
| C | -5% | +10% |
| D | -10% | +10% |

#### 2.5.6 Combat Losses by Mode

| Mode | Ships Lost | He3 Lost |
|------|-----------|----------|
| Normal/Restricted Instances | Yes | Yes |
| Trial/Constellation Instances | No | Yes |
| League/Arena/Championships | No | No |
| PvP (Attack Neighbors) | Yes | Yes |

#### 2.5.7 PvP Mechanics

- Scout with single-ship fleet before major attacks
- Winner receives **20% of loser's resources** (not warehouse contents)
- Use "synchronize arrival" to coordinate fleet timing
- Ships must be dismissed or truced when offline for protection

### 2.6 Commander System

#### 2.6.1 Commander Rarity Tiers (GO2 Original)

| Tier | Power Level | Acquisition |
|------|-------------|-------------|
| **Common** | Lowest | Free recruitment (with cooldown) |
| **Skill** | Low | Lucky draws, events |
| **Super** | Medium | Lucky Wheel, instances, merging |
| **Legendary** | High | Restricted Instances, Lucky Wheel, events |
| **Divine** | Highest | Premium events, special merging |

#### 2.6.2 Commander Attributes

| Attribute | Effect |
|-----------|--------|
| **Accuracy** | Increases weapon hit chance |
| **Dodge** | Reduces opponent hit rate |
| **Speed** | Determines attack order; affects successive strike chance |
| **Electron** | Increases Critical Hit Rate and Critical Damage |

#### 2.6.3 Star Rank & Effective Stack

Commander cards can be merged in the **Compound Center** to increase Star Rank. This is the **single biggest upgrade in the game** because Star Rank directly increases the Effective Stack (how many ships can attack per round).

[NEEDS RESEARCH: exact Star Rank to Effective Stack bonus table]

#### 2.6.4 Skill Types

| Skill Color | Type |
|-------------|------|
| Blue | Defense Skill |
| Red | Attack Skill |
| Green | Energy Skill |

[NEEDS RESEARCH: full skill list with effects]

#### 2.6.5 Enhancement Systems

- **Gems**: Attach to commander to increase dodge, attack, defense, electron (max 3 diamonds)
- **Bionic Chips**: Advanced enhancement via Cybernetics Center
  - Max 5 chips per commander
  - No duplicate chip types
  - Types: Max Planetary, Max Ballistics, Negator, etc.
  - Purchased with Corsairs' Gold or Mall Points

[NEEDS RESEARCH: complete gem types, bionic chip types, and their stat values]

#### 2.6.6 Recruitment Mechanics

| Method | Cost | Notes |
|--------|------|-------|
| Free Recruitment | Free | 3-hour cooldown between draws |
| Quick Recruitment | 8 Mall Points | No cooldown |
| Commander Cards | 100 Mall Points | Specific commander choice |
| Auction House | Variable | Player-to-player trading |
| Events | Free/Variable | Special events and updates |

### 2.7 Planets & Colonies

#### 2.7.1 Resource Bonus Planets (RBPs)

RBPs are the primary territorial mechanic. They are heavily fortified locations evenly spaced throughout the galaxy.

**Galaxy Layout**: 7x7+ grid of zones, each with 1 RBP in the center, spaced 60 movement squares apart.

#### 2.7.2 RBP Bonuses

| Level Range | Bonus per Level | Cumulative at Top |
|-------------|----------------|-------------------|
| 1-10 | 5% base + 0.5%/level | 10% |
| 11-20 | 1% per level | 20% |
| 21-30 | 1.5% per level | 35% |
| 31-40 | 2% per level | 55% |
| ... | Continues scaling | ... |
| 100 | - | **280%** |

Bonuses apply to: Resource production, Research speed, Shipbuilding speed

#### 2.7.3 RBP Defense

**Initial NPC Defenses:**
- 5 fleets of 4,500-7,200 Level 6 ships each
- Need fleets dealing ~1 million damage per swing to overcome

**Defensive Structures (per planet):**

| Structure | Count | HP Range (Lv1-10) | Attack Range |
|-----------|-------|-------------------|--------------|
| Meteor Stars | 63 | 40K - 20.48M | - |
| Particle Cannons | 8 | - | 8-17 tiles, 10K-450K dmg |
| Anti-Aircraft Guns | 12 | 16K - 8.19M | 30 tiles, 50K-500K dmg |
| Thor's Cannons | 5 | - | 50K-1.35M dmg |

**Fleet Capacity:**

| RBP Level | Max Defense Fleets |
|-----------|-------------------|
| 1 | 10 |
| 5 | 10 |
| 10 | 15 |
| 15 | 20 |
| ... | +5 per 5 levels |
| 100 | 105 |

#### 2.7.4 Conquest Mechanics

- Only **Corps** (not individuals) can attack RBPs
- **72-hour protection** after capture
- **24-hour battle phase** when vulnerable
- Successful takeover resets timer to 72 hours
- Corps have **99 combat rounds** to complete takeover
- Fleets cannot be recalled mid-battle
- On conquest: defenses reset to Level 1, Space Station retains level

**Multi-Corp Attacks**: Planet goes to the corp with the most kills from a single fleet type. Scoring: 1 point per ship destroyed, per structure destroyed, per Space Station destroyed.

#### 2.7.5 RBP Control Limits

Corps can control one planet per Corp Level (Lv 10 corp = max 10 RBPs).

#### 2.7.6 Upgrading RBPs

- Upgraded using **Corp Wealth** (from member donations)
- Space Station upgradeable to Level 100
- Defenses (Meteor Stars, Cannons, etc.) upgradeable to Level 10
- Total wealth for Level 100 Station: **24,902,439**

### 2.8 Corps (Alliances)

#### 2.8.1 Overview

Corps are groups of players that work together as a military/economic unit. They are essential for mid-to-late game content.

#### 2.8.2 Corp Features

| Feature | Description |
|---------|-------------|
| **Resource Bonus Planets** | Territorial control for production bonuses |
| **Corp Mall** | Shared ship shop (better inventory at higher levels) |
| **Corp Warehouse** | Shared storage facility |
| **Corp Merging Center** | Shared commander merging |
| **Pirate Planets** | Corp-level NPC combat challenges |
| **Galactic Wars** | Corp vs Corp large-scale conflict |

#### 2.8.3 Donation System

| Metric | Value |
|--------|-------|
| **Contribution Points** | 1 point per 10,000 resources donated |
| **Mall Points conversion** | 1:1 ratio (much more efficient) |
| **Daily donation range** | 20-200 contribution points |
| **Max daily donation** | 2,000,000 resources = 200 points |

#### 2.8.4 Corp Bonuses

- Resource production bonus (based on Corp level + RBP count)
- Science research acceleration
- Shipbuilding speed boost
- Access to Corp Mall inventory
- Protection from attacks (deterrent)

#### 2.8.5 Corp Levels

Higher Corp level = more RBPs controllable (1 per level), better Mall inventory, stronger bonuses.

[NEEDS RESEARCH: Corp level requirements, max level, member limits per level]

### 2.9 Instances (PvE)

#### 2.9.1 Instance Types

| Type | Description | Losses |
|------|-------------|--------|
| **Normal** | Standard progression (difficulty scaling) | Ships + He3 |
| **Restricted** | Limited access, commander rewards | Ships + He3 |
| **Scenario (Trial)** | Trial-based combat scenarios | He3 only |
| **Constellation** | Advanced endgame content | He3 only |

#### 2.9.2 Rewards

- **Treasure Boxes**: 10% blueprint chance, split equally among available blueprints
  - Example: 5 possible blueprints = 2% each
- **Resources**: Gold, Metal, He3
- **Commanders**: Restricted Instances (levels 8-10)
- **Badges**: From Restricted Instances
- **Blueprints**: Instance-specific drops

[NEEDS RESEARCH: instance list with levels, enemy compositions, and reward tables]

### 2.10 Progression & Quest System

GO2 has two main quest categories: **Development Quests** (one-time completion) and **Daily Quests** (repeatable). The quest system serves as the primary tutorial and progression guide. Quest rewards sustain players for approximately one week before natural resource production takes over.

**Key mechanic**: Players claim main quest rewards immediately but can SAVE side quest rewards until needed (they act as resource banks).

#### 2.10.1 Development Quests - Main Quest Chain

Main quests must be completed in sequential order. Each quest unlocks the next.

| # | Quest Name | Requirement | Metal | He3 | Gold | Special Reward |
|---|-----------|------------|-------|-----|------|---------------|
| 1 | Collecting Resources | Harvest Resource Warehouse | 450 | 950 | 500 | Loudspeaker |
| 2 | Loud and Clear | Say something in World Channel | 460 | 980 | 520 | - |
| 3 | Level 1 Technology Center | Build Technology Center Lv1 | 2,250 | 2,100 | 3,250 | Loudspeaker |
| 4 | Level 1 Technological Research | Research Concurrent Construction Lv1 | 500 | 480 | 600 | Construction Card |
| 5 | Metal Production | Build Metal Collector Lv1 | 425 | 530 | 425 | **Super Transmission Engine BP** |
| 6 | Blueprints 1 | Use Super Transmission Engine blueprint | 515 | 1,190 | 575 | Truce Card |
| 7 | He3 Production | Build He3 Extractor Lv1 | 475 | 400 | 475 | **Estrella BP** |
| 8 | Blueprints 2 | Use Estrella blueprint | 520 | 1,150 | 585 | Ship Reinforcement Facility BP |
| 9 | Creating Residential Area | Build Residential Area Lv1 | 390 | 360 | 325 | Loudspeaker |
| 10 | Level 1 Ship Factory | Build Ship Factory Lv1 | 600 | 475 | 550 | **Typhoon BP** |
| 11 | Level 1 Command Center | Build Command Center Lv1 | 3,000 | 2,250 | 2,500 | **Energy Shield Booster BP** |
| 12 | Recruit Commanders | Recruit 1 Commander | 605 | 1,250 | 685 | Revival Card |
| 13 | Design a Ship | Complete a ship design | 585 | 1,235 | 650 | Loudspeaker |
| 14 | Ship Building | Build Lv1 ships | 600 | 1,200 | 680 | **Anti-Aircraft Cannon BP** |
| 15 | Build a Fleet | Create 1 Fleet | 610 | 1,300 | 690 | Loudspeaker |
| 16 | Wartime Logistics | Replenish Ammunition | 2,000 | 3,000 | 2,000 | Truce Card |
| 17 | Level 1 Weapon Research Center | Build Weapon Research Center Lv1 | 2,500 | 1,500 | 2,250 | **Starlight Missile Pod BP** |
| 18 | Bigger Bags | Increase bag slot by 1 | 620 | 1,310 | 700 | Primary Metal Pack |
| 19 | Resource Pack | Use a Resource Pack | 680 | 1,450 | 750 | Galaxy Transfer |
| 20 | Growing Resources | Grow resources on Comsats | 700 | 1,480 | 780 | Healing Card |
| 21 | Adding Friends | Add 1 Friend | 710 | 1,500 | 800 | - |
| 22 | The Mail System | Send an email | 720 | 1,520 | 820 | Loudspeaker |
| 23 | Level 2 Space Station | Upgrade Space Station to Lv2 | 10,465 | 9,660 | 13,685 | - |
| 24 | Space Defense 1 | Build 1 defensive structure | 750 | 1,600 | 880 | - |
| 25 | Level 2 Metal Collector | Upgrade Metal Collector to Lv2 | 815 | 690 | 815 | Metal Mining Boost |
| 26 | Level 2 He3 Extractor | Upgrade He3 Extractor to Lv2 | 730 | 910 | 730 | He3 Mining Boost |
| 27 | Level 2 Residential Area | Upgrade Residential Area to Lv2 | 670 | 620 | 560 | Extra Tax |
| 28 | Level 2 Ship Factory | Build additional Ship Factories | - | - | - | - |

**Quest phase mapping:**
- Quests 1-9: Phase 1 building fundamentals (resource production, tech, residential)
- Quests 10-17: Phase 2 military systems (ships, fleets, combat, weapons)
- Quests 18-22: Social/utility features (inventory, friends, mail)
- Quests 23-28: Upgrade/expansion cycle

**Blueprint rewards in main quest chain** (critical starter blueprints):
1. **Super Transmission Engine** (Quest 5) - Transmission module
2. **Estrella** (Quest 7) - Ship hull
3. **Ship Reinforcement Facility** (Quest 8) - Structure module
4. **Typhoon** (Quest 10) - Ship hull
5. **Energy Shield Booster** (Quest 11) - Shield module
6. **Anti-Aircraft Cannon** (Quest 14) - Air Defense module
7. **Starlight Missile Pod** (Quest 17) - Missile weapon

#### 2.10.2 Development Quests - Side Quests

Side quests run in parallel with main quests. Completing one tier unlocks the next harder tier. There are approximately 24 side quest categories organized into groups:

**Resource Production Side Quests:**

| Category | Requirement Pattern | Example Tiers |
|----------|-------------------|---------------|
| Harvest Time | Increase Metal productivity to X | Lv1: 2,180/hr ... Lv9+: 12,970+/hr |
| Gathering He3 | Increase He3 productivity to X | Up to Lv20: 131,120 He3/hr |
| Raising Morale | Increase Gold productivity to X | Lv1: 500 reward ... Lv9: 40,000 reward |
| Plentiful Resources | Store X amount of Metal, He3, Gold | Escalating thresholds |

**Military Side Quests:**

| Category | Requirement Pattern |
|----------|-------------------|
| Building Ships | Own a set number of ships |
| Ship Research | Upgrade ship blueprints N times |
| Parts Research | Upgrade module parts N times |
| Recruitment | Recruit N commanders |
| Beefing Up Defenses | Build N defensive structures |
| Military Ranks | Own a commander card with star level X |

**Social/Economic Side Quests:**

| Category | Requirement Pattern |
|----------|-------------------|
| Friendly Faces | Add N friends |
| Free Trade | Reach trading volume X in auction house |
| Contribution | Reach N contribution points with corps |
| Staying Green | Recycle N ships |
| Joining Forces | Join/participate in corps activities |

**Speedup Side Quests:**

| Category | Requirement | Reward Example |
|----------|------------|---------------|
| Speed up Tech Advancement | Use speedups on tech research | 2,000M, 3,000H, 2,200G + 5 Vouchers |
| Speed up Module Blueprint Research | Use speedups on module research | 2,600M, 4,200H, 3,000G + 5 Vouchers |
| Speed up Shipbuilding | Use speedups on ship construction | 2,200M, 3,500H, 2,500G + 5 Vouchers |
| Speed up Construction | Use speedups on buildings | 2,500M, 4,000H, 2,800G + 5 Vouchers |

**Instance/Combat Side Quests:**

| Category | Requirement Pattern |
|----------|-------------------|
| Instances | Complete N instances |
| Peace Agreement | Use N truce cards |

**Side quest reward scaling**: Rewards increase significantly with tiers. Example: Ship Blueprint Research Tier 1 (upgrade 1 BP) = ~5,000 resources; Tier max (upgrade 90 BPs) = ~3,000,000 resources. Special items (Galaxy Transfer, Construction Card, Sealing Card, Advanced Galaxy Transfer) appear at higher tiers.

#### 2.10.3 Daily Quests

Daily quests reset at 1:00 PM server time. Completing them earns points toward accumulated tier rewards.

| Quest | Requirement | Points |
|-------|------------|--------|
| Daily Log In | Log in to the game | 10 |
| Collect Your Dues | Harvest your own asteroids once | 4 |
| Finders Keepers | Harvest matured asteroids of 5 friends | 2 per friend (max 10) |
| Helping Hand | Speed up construction/upgrade for 5 friends | 2 per building (max 10) |
| Need for Speed | Speed up your own construction/upgrade | 3 |
| Stockpiling | Harvest from Resource Warehouse 3 times | 1 per harvest (max 3) |
| Base Construction | Give wrenches to 5 friends' bases | 2 per building (max 10) |
| Voucher Rush | Harvest matured vouchers from Celestial Industrial Base | 4 |
| Donations | Donate 200,000 resources to Corps | 6 |
| Restricted Instances | Attempt 2 Restricted Instances (win not required) | 5 per attempt (max 10) |

**Maximum daily points**: ~70 (if all quests completed fully)

**Accumulated Point Reward Tiers:**

| Points | Tier | Reward |
|--------|------|--------|
| 10 | Bronze | Random: Loudspeaker (50%), Resource Box (20%), SP Card (30%) |
| 30 | Silver | Random: Loudspeaker (40%), Resource Box (30%), SP Card (10%), Resource Packs (10%), He3/Metal Packs (10%) |
| 50 | Gold | Random: Resource Box (30%), SP Card (30%), Resource Packs (10%), Mining Boosts (5% ea), Cards/Items (2-5%) |
| 70 | Diamond | **Raw Gemstone (100%)** - guaranteed |

**Key mechanics:**
- Points reset daily but earned rewards persist
- Reaching a threshold doesn't consume points (reaching 70 gets ALL four tiers)
- The 70-point Diamond tier always gives a Raw Gemstone (valuable crafting material)

#### 2.10.4 Quest State Machine

```
MAIN QUEST STATES:
  locked -> available -> in_progress -> completed -> claimed
  - locked: prerequisite quest not yet completed
  - available: prerequisite met, player can start tracking
  - in_progress: player is working on requirement
  - completed: requirement met, reward not yet claimed
  - claimed: reward collected (terminal state)

SIDE QUEST STATES:
  locked -> available -> completed -> claimed -> (next tier unlocked)
  - locked: category not yet accessible
  - available: current tier requirement visible
  - completed: tier requirement met
  - claimed: tier reward collected, next tier becomes available

DAILY QUEST STATES (per day):
  available -> completed -> claimed
  - Resets to available at 1:00 PM server time daily
  - Points accumulate toward tier thresholds
  - Tier rewards are claimable once threshold is reached
```

#### 2.10.5 Quest Phasing for Cryptomines Online

**Phase 1 Quests** (building/resource focus):
- Main quest chain: Quests 1, 3, 5, 7, 9, 23, 25-27 (building progression only)
- Resource side quests: Harvest Time, Gathering He3, Raising Morale, Plentiful Resources
- Daily quests: Daily Login, Collect Your Dues, Stockpiling (simplified set)

**Phase 2 Quests** (military focus):
- Main quest chain extension: Quests 10-17, 24 (military systems)
- Military side quests: Building Ships, Ship Research, Parts Research, Recruitment, Beefing Up Defenses
- Instance side quests: Complete N instances
- Daily quests extension: Restricted Instances, Donations

**Deferred quests** (require social/chat systems):
- Quests 2, 21, 22 (Loud and Clear, Adding Friends, The Mail System)
- Side quests: Friendly Faces, Free Trade, Contribution, Joining Forces
- Daily quests: Finders Keepers, Helping Hand, Base Construction

#### 2.10.6 Daily Activities Summary

| Activity | Reward |
|----------|--------|
| Daily Quests | Points toward tier rewards (Bronze/Silver/Gold/Diamond) |
| Friend Visits | Up to 8 vouchers per 24 hours (from repairs) |
| Resource Collection | Manual harvest from warehouses |
| Instance Farming | Blueprints, resources, commanders |
| League Matches | Honor Points (15-255 per match) |

### 2.11 Monetization

#### 2.11.1 Free-to-Play Model

"Money helps with two things: **rarity** and **time**."

| Purchase Type | Effect |
|---------------|--------|
| **Mall Points** | Premium currency for exclusive items, speedups, commanders |
| **Speedups** | Accelerate building, research, ship production |
| **Exclusive Blueprints** | Ships not available through free progression |
| **Exclusive Commanders** | Higher rarity commanders |
| **Convenience** | Extra construction slots, instant recruitment |

#### 2.11.2 Free Player Progression

Free players can access most content through:
- Vouchers (daily free currency)
- Instance farming
- League/Championship rewards
- Corp bonuses
- Time investment

---

## 3. Data Models (SQL)

All schemas target **Supabase (PostgreSQL)**. Ready for migration. Resource columns use GO2 names: `metal`, `he3`, `gold`.

### 3.1 Players

```sql
CREATE TABLE players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    anonymous_id TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE,
    level INTEGER NOT NULL DEFAULT 1,
    experience BIGINT NOT NULL DEFAULT 0,
    mall_points INTEGER NOT NULL DEFAULT 0,
    vouchers INTEGER NOT NULL DEFAULT 0,
    honor_points INTEGER NOT NULL DEFAULT 0,
    champion_points INTEGER NOT NULL DEFAULT 0,
    badges INTEGER NOT NULL DEFAULT 0,
    corsairs_gold INTEGER NOT NULL DEFAULT 0,
    tutorial_step INTEGER NOT NULL DEFAULT 0,
    is_online BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_login TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_players_anonymous_id ON players (anonymous_id);
CREATE INDEX idx_players_username ON players (username) WHERE username IS NOT NULL;
CREATE INDEX idx_players_last_login ON players (last_login);
```

### 3.2 Planets

```sql
CREATE TABLE planets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    name TEXT NOT NULL DEFAULT 'New Colony',
    position_x INTEGER NOT NULL,
    position_y INTEGER NOT NULL,
    is_homeworld BOOLEAN NOT NULL DEFAULT false,
    -- RBP fields (null for player homeworlds)
    is_rbp BOOLEAN NOT NULL DEFAULT false,
    rbp_level INTEGER NOT NULL DEFAULT 0,
    controlling_corp_id UUID,
    protection_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_planet_position UNIQUE (position_x, position_y)
);

CREATE INDEX idx_planets_player_id ON planets (player_id);
CREATE INDEX idx_planets_position ON planets (position_x, position_y);
CREATE INDEX idx_planets_rbp ON planets (is_rbp) WHERE is_rbp = true;
```

### 3.3 Building Types (Reference Data)

```sql
CREATE TABLE building_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    display_name TEXT NOT NULL,
    category TEXT NOT NULL CHECK (category IN ('resource', 'core', 'military', 'defense', 'space', 'decorative')),
    base TEXT NOT NULL CHECK (base IN ('ground', 'space')) DEFAULT 'ground',
    base_cost_metal BIGINT NOT NULL DEFAULT 0,
    base_cost_he3 BIGINT NOT NULL DEFAULT 0,
    base_cost_gold BIGINT NOT NULL DEFAULT 0,
    base_time_seconds INTEGER NOT NULL DEFAULT 60,
    cost_multiplier NUMERIC(6,4) NOT NULL DEFAULT 3.0300,
    time_multiplier NUMERIC(6,4) NOT NULL DEFAULT 2.8700,
    base_production_per_hour INTEGER NOT NULL DEFAULT 0,
    production_multiplier NUMERIC(6,4) NOT NULL DEFAULT 1.1340,
    max_level INTEGER NOT NULL DEFAULT 24,
    max_count_per_planet INTEGER NOT NULL DEFAULT 1,
    prerequisite_building TEXT,
    prerequisite_level INTEGER NOT NULL DEFAULT 0,
    civic_center_req_per_level BOOLEAN NOT NULL DEFAULT true,
    description TEXT NOT NULL DEFAULT '',

    CONSTRAINT chk_positive_costs CHECK (
        base_cost_metal >= 0 AND base_cost_he3 >= 0 AND base_cost_gold >= 0
    ),
    CONSTRAINT chk_positive_multipliers CHECK (
        cost_multiplier > 0 AND time_multiplier > 0
    )
);
```

### 3.4 Buildings

```sql
CREATE TABLE buildings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    planet_id UUID NOT NULL REFERENCES planets(id) ON DELETE CASCADE,
    building_type INTEGER NOT NULL REFERENCES building_types(id),
    level INTEGER NOT NULL DEFAULT 1,
    is_upgrading BOOLEAN NOT NULL DEFAULT false,
    upgrade_finish_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_level_positive CHECK (level >= 1),
    CONSTRAINT chk_upgrade_consistency CHECK (
        (is_upgrading = true AND upgrade_finish_at IS NOT NULL) OR
        (is_upgrading = false AND upgrade_finish_at IS NULL)
    )
);

CREATE INDEX idx_buildings_planet_id ON buildings (planet_id);
CREATE INDEX idx_buildings_type ON buildings (building_type);
CREATE INDEX idx_buildings_upgrading ON buildings (is_upgrading) WHERE is_upgrading = true;
```

### 3.5 Resources

```sql
CREATE TABLE resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    planet_id UUID UNIQUE NOT NULL REFERENCES planets(id) ON DELETE CASCADE,
    metal BIGINT NOT NULL DEFAULT 0,
    he3 BIGINT NOT NULL DEFAULT 0,
    gold BIGINT NOT NULL DEFAULT 0,
    metal_per_hour BIGINT NOT NULL DEFAULT 0,
    he3_per_hour BIGINT NOT NULL DEFAULT 0,
    gold_per_hour BIGINT NOT NULL DEFAULT 0,
    storage_capacity BIGINT NOT NULL DEFAULT 100000,
    last_collected_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_resources_non_negative CHECK (
        metal >= 0 AND he3 >= 0 AND gold >= 0
    ),
    CONSTRAINT chk_rates_non_negative CHECK (
        metal_per_hour >= 0 AND he3_per_hour >= 0 AND gold_per_hour >= 0
    )
);
```

### 3.6 Tech Types (Reference Data)

```sql
CREATE TABLE tech_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    display_name TEXT NOT NULL,
    tree TEXT NOT NULL CHECK (tree IN (
        'logistics_construction', 'planetary_defense', 'ballistics_science',
        'directional_science', 'missile_science', 'ship_based_science', 'ship_defense_science'
    )),
    prerequisites_json JSONB NOT NULL DEFAULT '[]',
    -- Note: In GO2 all tech research costs Gold only (metal=0, he3=0).
    -- Schema retains metal/he3 columns for flexibility.
    base_cost_metal BIGINT NOT NULL DEFAULT 0,
    base_cost_he3 BIGINT NOT NULL DEFAULT 0,
    base_cost_gold BIGINT NOT NULL DEFAULT 0,
    cost_multiplier NUMERIC(6,4) NOT NULL DEFAULT 1.5300,
    base_time_seconds INTEGER NOT NULL DEFAULT 117,
    time_multiplier NUMERIC(6,4) NOT NULL DEFAULT 2.3400,
    max_level INTEGER NOT NULL DEFAULT 10,
    effects_json JSONB NOT NULL DEFAULT '{}',
    description TEXT NOT NULL DEFAULT ''
);

CREATE INDEX idx_tech_types_tree ON tech_types (tree);
```

**prerequisites_json format:** `[{"tech": "ballistics_base", "level": 3}, ...]`

**effects_json format:** `{"type": "ballistic_damage", "per_level": 5, "unit": "percent"}`

### 3.7 Technologies

```sql
CREATE TABLE technologies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    tech_type INTEGER NOT NULL REFERENCES tech_types(id),
    level INTEGER NOT NULL DEFAULT 0,
    is_researching BOOLEAN NOT NULL DEFAULT false,
    research_finish_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_player_tech UNIQUE (player_id, tech_type),
    CONSTRAINT chk_level_non_negative CHECK (level >= 0),
    CONSTRAINT chk_research_consistency CHECK (
        (is_researching = true AND research_finish_at IS NOT NULL) OR
        (is_researching = false AND research_finish_at IS NULL)
    )
);

CREATE INDEX idx_technologies_player_id ON technologies (player_id);
CREATE INDEX idx_technologies_researching ON technologies (is_researching) WHERE is_researching = true;
```

### 3.8 Ship Designs

```sql
CREATE TABLE ship_designs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    name TEXT NOT NULL DEFAULT 'New Design',
    hull_type TEXT NOT NULL CHECK (hull_type IN ('frigate', 'cruiser', 'battleship')),
    armor_type TEXT NOT NULL DEFAULT 'nano' CHECK (armor_type IN ('nano', 'chrome', 'regen', 'neutralizing')),
    modules_json JSONB NOT NULL DEFAULT '[]',
    stats_json JSONB NOT NULL DEFAULT '{}',
    blueprint_id TEXT,
    metal_cost BIGINT NOT NULL DEFAULT 0,
    he3_cost BIGINT NOT NULL DEFAULT 0,
    gold_cost BIGINT NOT NULL DEFAULT 0,
    build_time_seconds INTEGER NOT NULL DEFAULT 60,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_ship_designs_player_id ON ship_designs (player_id);
CREATE INDEX idx_ship_designs_hull ON ship_designs (hull_type);
```

### 3.9 Ships

```sql
CREATE TABLE ships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    ship_design_id UUID NOT NULL REFERENCES ship_designs(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 0,
    is_building BOOLEAN NOT NULL DEFAULT false,
    build_quantity INTEGER NOT NULL DEFAULT 0,
    build_finish_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_player_ship_design UNIQUE (player_id, ship_design_id),
    CONSTRAINT chk_quantity_non_negative CHECK (quantity >= 0),
    CONSTRAINT chk_build_consistency CHECK (
        (is_building = true AND build_finish_at IS NOT NULL AND build_quantity > 0) OR
        (is_building = false AND build_finish_at IS NULL AND build_quantity = 0)
    )
);

CREATE INDEX idx_ships_player_id ON ships (player_id);
```

### 3.10 Commanders

```sql
CREATE TABLE commanders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    rarity TEXT NOT NULL CHECK (rarity IN ('common', 'skill', 'super', 'legendary', 'divine')),
    star_rank INTEGER NOT NULL DEFAULT 0,
    accuracy INTEGER NOT NULL DEFAULT 0,
    dodge INTEGER NOT NULL DEFAULT 0,
    speed INTEGER NOT NULL DEFAULT 0,
    electron INTEGER NOT NULL DEFAULT 0,
    weapon_expertise JSONB NOT NULL DEFAULT '{"ballistic":"B","directional":"B","missile":"B","ship_based":"B"}',
    ship_expertise JSONB NOT NULL DEFAULT '{"frigate":"B","cruiser":"B","battleship":"B"}',
    skills_json JSONB NOT NULL DEFAULT '[]',
    gems_json JSONB NOT NULL DEFAULT '[]',
    bionic_chips_json JSONB NOT NULL DEFAULT '[]',
    is_deployed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_star_rank CHECK (star_rank >= 0 AND star_rank <= 15),
    CONSTRAINT chk_stats_non_negative CHECK (
        accuracy >= 0 AND dodge >= 0 AND speed >= 0 AND electron >= 0
    )
);

CREATE INDEX idx_commanders_player_id ON commanders (player_id);
CREATE INDEX idx_commanders_rarity ON commanders (rarity);
```

### 3.11 Fleets

```sql
CREATE TABLE fleets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    name TEXT NOT NULL DEFAULT 'Fleet',
    formation TEXT NOT NULL DEFAULT 'phalanx'
        CHECK (formation IN ('phalanx', 'diamond', 'battle_line', 'skirmish', 'tee_forward', 'enfilade', 'tee_reverse')),
    commander_id UUID REFERENCES commanders(id) ON DELETE SET NULL,
    targeting_command TEXT NOT NULL DEFAULT 'max_attack'
        CHECK (targeting_command IN ('max_attack', 'min_attack', 'max_durability', 'min_durability', 'closest', 'by_commander_rank')),
    grid_json JSONB NOT NULL DEFAULT '[[null,null,null],[null,null,null],[null,null,null]]',
    status TEXT NOT NULL DEFAULT 'stationed'
        CHECK (status IN ('stationed', 'traveling', 'combat', 'returning', 'dismissed')),
    planet_id UUID REFERENCES planets(id) ON DELETE SET NULL,
    position_x INTEGER,
    position_y INTEGER,
    destination_x INTEGER,
    destination_y INTEGER,
    arrival_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_travel_consistency CHECK (
        (status IN ('traveling', 'returning') AND destination_x IS NOT NULL
            AND destination_y IS NOT NULL AND arrival_at IS NOT NULL) OR
        (status IN ('stationed', 'combat', 'dismissed'))
    )
);

CREATE INDEX idx_fleets_player_id ON fleets (player_id);
CREATE INDEX idx_fleets_status ON fleets (status);
CREATE INDEX idx_fleets_commander ON fleets (commander_id) WHERE commander_id IS NOT NULL;
```

### 3.12 Combat Reports

```sql
CREATE TABLE combat_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attacker_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    defender_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    attacker_fleet_id UUID REFERENCES fleets(id) ON DELETE SET NULL,
    defender_fleet_id UUID REFERENCES fleets(id) ON DELETE SET NULL,
    combat_type TEXT NOT NULL DEFAULT 'pvp'
        CHECK (combat_type IN ('pvp', 'instance_normal', 'instance_restricted', 'instance_trial', 'instance_constellation', 'league', 'championship', 'rbp_attack')),
    result TEXT NOT NULL CHECK (result IN ('attacker_win', 'defender_win', 'draw')),
    total_rounds INTEGER NOT NULL DEFAULT 0,
    rounds_json JSONB NOT NULL DEFAULT '[]',
    loot_json JSONB NOT NULL DEFAULT '{}',
    attacker_losses_json JSONB NOT NULL DEFAULT '{}',
    defender_losses_json JSONB NOT NULL DEFAULT '{}',
    he3_consumed BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_combat_reports_attacker ON combat_reports (attacker_id);
CREATE INDEX idx_combat_reports_defender ON combat_reports (defender_id);
CREATE INDEX idx_combat_reports_created ON combat_reports (created_at DESC);
CREATE INDEX idx_combat_reports_type ON combat_reports (combat_type);
```

### 3.13 Instances (PvE Content)

```sql
CREATE TABLE instances (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('normal', 'restricted', 'trial', 'constellation')),
    difficulty INTEGER NOT NULL DEFAULT 1,
    required_level INTEGER NOT NULL DEFAULT 1,
    ships_lost_on_defeat BOOLEAN NOT NULL DEFAULT true,
    he3_lost_on_defeat BOOLEAN NOT NULL DEFAULT true,
    enemy_fleets_json JSONB NOT NULL DEFAULT '[]',
    rewards_json JSONB NOT NULL DEFAULT '{}',
    treasure_box_blueprints_json JSONB NOT NULL DEFAULT '[]',
    description TEXT NOT NULL DEFAULT '',

    CONSTRAINT chk_difficulty_positive CHECK (difficulty >= 1),
    CONSTRAINT chk_level_positive CHECK (required_level >= 1)
);

CREATE INDEX idx_instances_type ON instances (type);
CREATE INDEX idx_instances_difficulty ON instances (difficulty);
```

### 3.14 Instance Progress

```sql
CREATE TABLE instance_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    instance_id INTEGER NOT NULL REFERENCES instances(id) ON DELETE CASCADE,
    completed BOOLEAN NOT NULL DEFAULT false,
    attempts INTEGER NOT NULL DEFAULT 0,
    best_score INTEGER NOT NULL DEFAULT 0,
    last_attempt_at TIMESTAMPTZ,

    CONSTRAINT uq_player_instance UNIQUE (player_id, instance_id),
    CONSTRAINT chk_attempts_non_negative CHECK (attempts >= 0),
    CONSTRAINT chk_score_non_negative CHECK (best_score >= 0)
);

CREATE INDEX idx_instance_progress_player ON instance_progress (player_id);
```

### 3.15 Corps (Alliances)

```sql
CREATE TABLE corps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    tag TEXT UNIQUE NOT NULL,
    leader_id UUID NOT NULL REFERENCES players(id),
    level INTEGER NOT NULL DEFAULT 1,
    wealth BIGINT NOT NULL DEFAULT 0,
    max_members INTEGER NOT NULL DEFAULT 20,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_tag_length CHECK (char_length(tag) BETWEEN 2 AND 5),
    CONSTRAINT chk_name_length CHECK (char_length(name) BETWEEN 3 AND 30),
    CONSTRAINT chk_level_positive CHECK (level >= 1),
    CONSTRAINT chk_wealth_non_negative CHECK (wealth >= 0)
);

CREATE INDEX idx_corps_leader ON corps (leader_id);
CREATE INDEX idx_corps_name ON corps (name);
```

### 3.16 Corp Members

```sql
CREATE TABLE corp_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    corp_id UUID NOT NULL REFERENCES corps(id) ON DELETE CASCADE,
    player_id UUID UNIQUE NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    role TEXT NOT NULL DEFAULT 'member'
        CHECK (role IN ('leader', 'officer', 'member')),
    contribution_points BIGINT NOT NULL DEFAULT 0,
    daily_contribution_today INTEGER NOT NULL DEFAULT 0,
    last_contribution_date DATE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_corp_player UNIQUE (corp_id, player_id),
    CONSTRAINT chk_contribution_non_negative CHECK (contribution_points >= 0),
    CONSTRAINT chk_daily_contribution CHECK (daily_contribution_today >= 0 AND daily_contribution_today <= 200)
);

CREATE INDEX idx_corp_members_corp ON corp_members (corp_id);
CREATE INDEX idx_corp_members_player ON corp_members (player_id);
```

### 3.17 Blueprints

```sql
CREATE TABLE blueprints (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    hull_type TEXT NOT NULL CHECK (hull_type IN ('frigate', 'cruiser', 'battleship')),
    base_shield INTEGER NOT NULL DEFAULT 0,
    base_structure INTEGER NOT NULL DEFAULT 0,
    base_stability NUMERIC(5,2) NOT NULL DEFAULT 100.00,
    base_defense NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    installation_slots INTEGER NOT NULL DEFAULT 100,
    base_agility INTEGER NOT NULL DEFAULT 0,
    base_movement INTEGER NOT NULL DEFAULT 0,
    base_storage INTEGER NOT NULL DEFAULT 0,
    source TEXT NOT NULL DEFAULT 'instance',
    description TEXT NOT NULL DEFAULT ''
);

CREATE INDEX idx_blueprints_hull ON blueprints (hull_type);
```

### 3.18 Player Blueprints (Unlocked)

```sql
CREATE TABLE player_blueprints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    blueprint_id INTEGER NOT NULL REFERENCES blueprints(id) ON DELETE CASCADE,
    acquired_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_player_blueprint UNIQUE (player_id, blueprint_id)
);

CREATE INDEX idx_player_blueprints_player ON player_blueprints (player_id);
```

### 3.19 Seed Data: Building Types

```sql
INSERT INTO building_types (name, display_name, category, base, base_cost_metal, base_cost_he3, base_cost_gold, base_time_seconds, cost_multiplier, time_multiplier, base_production_per_hour, production_multiplier, max_level, max_count_per_planet, description) VALUES
-- Ground Base: Resource Buildings
('metal_collector',      'Metal Collector',      'resource', 'ground', 95,  80,  95,  40,  3.0300, 2.8700, 1080, 1.1340, 24, 8, 'Produces Metal'),
('he3_extractor',        'He3 Extractor',        'resource', 'ground', 95,  80,  95,  40,  3.0300, 2.8700, 1180, 1.1340, 24, 8, 'Produces He3 (Helium-3)'),
('residential_area',     'Residential Area',     'resource', 'ground', 95,  80,  95,  40,  3.0300, 2.8700, 1400, 1.1340, 24, 8, 'Produces Gold (highest output)'),
('resource_warehouse',   'Resource Warehouse',   'resource', 'ground', 120, 100, 120, 60,  3.0300, 2.8700, 0,    1.0000, 24, 4, 'Stores all resources; increases capacity'),
-- Ground Base: Core / Administrative Buildings
('civic_center',         'Civic Center',         'core', 'ground', 550, 480, 600, 300, 3.0300, 2.8700, 0, 1.0000, 12, 1, 'Main hub; determines max level of all other buildings'),
('technology_center',    'Technology Center',    'core', 'ground', 450, 420, 650, 100, 3.0300, 2.8700, 0, 1.0000, 12, 1, 'Research facility (7 science trees); 3% research time reduction per level'),
('alliance_center',      'Alliance Center',      'core', 'ground', 400, 350, 500, 120, 3.0300, 2.8700, 0, 1.0000, 12, 1, 'Enables Corp membership and features'),
('trading_center',       'Trading Center',       'core', 'ground', 300, 250, 400, 90,  3.0300, 2.8700, 0, 1.0000, 12, 1, 'Player-to-player trading and auctions'),
('galaxy_transporter',   'Galaxy Transporter',   'core', 'ground', 350, 300, 450, 100, 3.0300, 2.8700, 0, 1.0000, 12, 1, 'Inter-system resource transport'),
('compound_center',      'Compound Center',      'core', 'ground', 400, 350, 500, 100, 3.0300, 2.8700, 0, 1.0000, 12, 1, 'Commander card merging and enhancement'),
('radar',                'Radar',                'core', 'ground', 200, 300, 200, 80,  3.0300, 2.8700, 0, 1.0000, 12, 1, 'Detection of incoming attacks'),
-- Ground Base: Military Buildings
('ship_factory',         'Ship Factory',         'military', 'ground', 600, 450, 500, 200, 3.0300, 2.8700, 0, 1.0000, 24, 1, 'Constructs ships; holds 20 designs, 5 production slots'),
('spacedock',            'Spacedock',            'military', 'ground', 500, 400, 450, 150, 3.0300, 2.8700, 0, 1.0000, 12, 1, 'Ship berthing and fleet management'),
('command_center',       'Command Center',       'military', 'ground', 600, 450, 500, 40,  3.0300, 2.8700, 0, 1.0000, 12, 1, 'Commander recruitment (60 max at Lv71+)'),
('weapon_research_center','Weapon Research Center','military','ground', 500, 400, 600, 120, 3.0300, 2.8700, 0, 1.0000, 12, 1, 'Develops weapons and modules'),
('recycling_plant',      'Recycling Plant',      'military', 'ground', 400, 350, 400, 100, 3.0300, 2.8700, 0, 1.0000, 12, 1, 'Recovers resources from scrapped ships'),
-- Space Base Buildings
('space_station',        'Space Station',        'space', 'space', 550, 480, 600, 300, 3.0300, 2.8700, 0, 1.0000, 12, 1, 'Main orbital structure; must stay within 1 level of Civic Center'),
('meteor_star',          'Meteor Star',          'defense', 'space', 300, 250, 300, 120, 3.0300, 2.8700, 0, 1.0000, 10, 63, 'Orbital defense structure'),
('particle_cannon',      'Particle Cannon',      'defense', 'space', 400, 350, 400, 150, 3.0300, 2.8700, 0, 1.0000, 10, 8, 'Energy weapon defense'),
('anti_aircraft_gun',    'Anti-Aircraft Gun',    'defense', 'space', 350, 300, 350, 130, 3.0300, 2.8700, 0, 1.0000, 10, 12, 'Anti-air defense'),
('thors_cannon',         'Thor''s Cannon',       'defense', 'space', 500, 450, 500, 180, 3.0300, 2.8700, 0, 1.0000, 10, 5, 'Advanced heavy defense cannon'),
('celestial_base',       'Celestial Base',       'space', 'space', 600, 500, 600, 200, 3.0300, 2.8700, 0, 1.0000, 12, 1, 'Advanced orbital facility');
```

### 3.20 He3 Extractor - Complete Level Data

This is the reference data for He3 Extractor. Other resource buildings follow the same pattern.

```sql
-- Reference: He3 Extractor per-level data (from GO2 wiki)
-- This data can be used to validate formulas or as a lookup table
CREATE TABLE he3_extractor_levels (
    level INTEGER PRIMARY KEY,
    he3_output_per_hour INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL,
    civic_center_req INTEGER NOT NULL
);

INSERT INTO he3_extractor_levels VALUES
(1,  1180,   95,       80,       95,       40,      1),
(2,  1215,   163,      138,      163,      69,      1),
(3,  1264,   283,      238,      283,      119,     2),
(4,  1327,   492,      414,      492,      207,     2),
(5,  1407,   861,      725,      861,      362,     3),
(6,  1505,   1515,     1276,     1515,     638,     3),
(7,  1626,   2681,     2258,     2681,     1129,    4),
(8,  1772,   4773,     4019,     4773,     2010,    4),
(9,  1949,   8544,     7195,     8544,     3597,    5),
(10, 2164,   15379,    12950,    15379,    6475,    5),
(11, 2423,   27835,    23440,    27835,    11720,   6),
(12, 2738,   50660,    42661,    50660,    21331,   6),
(13, 3122,   92708,    78070,    92708,    39035,   7),
(14, 3590,   170583,   143649,   170583,   71824,   7),
(15, 4164,   315579,   265750,   315579,   132875,  8),
(16, 4872,   586976,   494296,   586976,   247148,  8),
(17, 5749,   1097645,  924333,   1097645,  462167,  9),
(18, 6842,   2063573,  1737746,  2063573,  868873,  9),
(19, 8210,   3900154,  3284340,  3900154,  1642170, 10),
(20, 9934,   7410292,  6240246,  7410292,  3120123, 10),
(21, 12120,  14153658, 11918870, 14153658, 5959435, 11),
(22, 14907,  27175024, 22884230, 27175024, 11442115, 11),
(23, 18485,  52447796, 44166565, 52447796, 22083283, 12),
(24, 23106,  101748723,85683135, 101748723,42841568, 12);
```

### 3.21 Metal Collector - Production Reference

```sql
CREATE TABLE metal_collector_levels (
    level INTEGER PRIMARY KEY,
    metal_output_per_hour INTEGER NOT NULL
);

INSERT INTO metal_collector_levels VALUES
(1,  1080),  (2,  1112),  (3,  1157),  (4,  1215),
(5,  1288),  (6,  1378),  (7,  1488),  (8,  1622),
(9,  1784),  (10, 1980),  (11, 2218),  (12, 2506),
(13, 2857),  (14, 3286),  (15, 3812),  (16, 4459),
(17, 5262),  (18, 6262),  (19, 7514),  (20, 9092),
(21, 11093), (22, 13644), (23, 16919), (24, 21148);
```

[NEEDS RESEARCH: Residential Area (Gold) per-level production table]
[NEEDS RESEARCH: Resource Warehouse per-level storage capacity table]
[NEEDS RESEARCH: Metal Collector per-level cost and build time table]

### 3.22 Quest Types (Reference Data)

```sql
CREATE TABLE quest_types (
    id SERIAL PRIMARY KEY,
    quest_key TEXT UNIQUE NOT NULL,
    category TEXT NOT NULL CHECK (category IN ('main', 'side', 'daily')),
    display_name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    requirement_type TEXT NOT NULL CHECK (requirement_type IN (
        'build_building', 'upgrade_building', 'harvest_resources', 'research_tech',
        'use_blueprint', 'recruit_commander', 'create_ship_design', 'build_ships',
        'create_fleet', 'replenish_ammo', 'build_defense', 'complete_instance',
        'reach_production', 'reach_storage', 'own_ships', 'upgrade_blueprints',
        'upgrade_modules', 'reach_star_rank', 'recycle_ships', 'use_speedup',
        'login', 'use_truce_card', 'donate_resources', 'send_message',
        'add_friend', 'increase_bag_slot', 'use_resource_pack', 'grow_comsats'
    )),
    requirement_target TEXT NOT NULL DEFAULT '',
    -- e.g., 'metal_collector', 'technology_center', or '' for generic
    requirement_value INTEGER NOT NULL DEFAULT 1,
    -- e.g., build level, quantity threshold, production rate target
    chain_order INTEGER,
    -- for main quests: sequential position (1-28+); for side quests: tier level; NULL for daily
    prerequisite_quest_id INTEGER REFERENCES quest_types(id) ON DELETE SET NULL,
    -- main quests: previous quest in chain; side quests: previous tier; daily: NULL
    reward_metal BIGINT NOT NULL DEFAULT 0,
    reward_he3 BIGINT NOT NULL DEFAULT 0,
    reward_gold BIGINT NOT NULL DEFAULT 0,
    reward_item_json JSONB NOT NULL DEFAULT '[]',
    -- e.g., [{"type": "blueprint", "blueprint_id": 5}, {"type": "item", "item_key": "loudspeaker", "quantity": 1}]
    phase INTEGER NOT NULL DEFAULT 1 CHECK (phase IN (1, 2, 3)),
    -- which implementation phase this quest belongs to
    is_active BOOLEAN NOT NULL DEFAULT true,
    -- false for deferred quests (social/chat-dependent)

    CONSTRAINT chk_chain_order_positive CHECK (chain_order IS NULL OR chain_order >= 1),
    CONSTRAINT chk_rewards_non_negative CHECK (reward_metal >= 0 AND reward_he3 >= 0 AND reward_gold >= 0)
);

CREATE INDEX idx_quest_types_category ON quest_types (category);
CREATE INDEX idx_quest_types_phase ON quest_types (phase);
CREATE INDEX idx_quest_types_chain ON quest_types (category, chain_order);
```

### 3.23 Player Quests (Progress Tracking)

```sql
CREATE TABLE player_quests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    quest_type_id INTEGER NOT NULL REFERENCES quest_types(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'locked'
        CHECK (status IN ('locked', 'available', 'in_progress', 'completed', 'claimed')),
    progress_value INTEGER NOT NULL DEFAULT 0,
    -- current progress toward requirement_value (e.g., ships built, production reached)
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    claimed_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_player_quest UNIQUE (player_id, quest_type_id),
    CONSTRAINT chk_progress_non_negative CHECK (progress_value >= 0),
    CONSTRAINT chk_status_consistency CHECK (
        (status = 'locked' AND completed_at IS NULL AND claimed_at IS NULL) OR
        (status = 'available' AND completed_at IS NULL AND claimed_at IS NULL) OR
        (status = 'in_progress' AND completed_at IS NULL AND claimed_at IS NULL) OR
        (status = 'completed' AND completed_at IS NOT NULL AND claimed_at IS NULL) OR
        (status = 'claimed' AND completed_at IS NOT NULL AND claimed_at IS NOT NULL)
    )
);

CREATE INDEX idx_player_quests_player ON player_quests (player_id);
CREATE INDEX idx_player_quests_status ON player_quests (player_id, status);
CREATE INDEX idx_player_quests_quest_type ON player_quests (quest_type_id);
```

### 3.24 Daily Quest Progress

```sql
CREATE TABLE daily_quest_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    quest_date DATE NOT NULL DEFAULT CURRENT_DATE,
    daily_points INTEGER NOT NULL DEFAULT 0,
    quests_completed_json JSONB NOT NULL DEFAULT '{}',
    -- e.g., {"daily_login": true, "collect_dues": true, "stockpiling": 2}
    tier_rewards_claimed_json JSONB NOT NULL DEFAULT '[]',
    -- e.g., ["bronze", "silver"]
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_player_daily UNIQUE (player_id, quest_date),
    CONSTRAINT chk_points_non_negative CHECK (daily_points >= 0 AND daily_points <= 100)
);

CREATE INDEX idx_daily_quest_player ON daily_quest_progress (player_id);
CREATE INDEX idx_daily_quest_date ON daily_quest_progress (quest_date DESC);
```

### 3.25 Seed Data: Main Quest Chain (Phase 1)

```sql
INSERT INTO quest_types (quest_key, category, display_name, requirement_type, requirement_target, requirement_value, chain_order, reward_metal, reward_he3, reward_gold, reward_item_json, phase, is_active) VALUES
-- Phase 1 Main Quests (building/resource focus)
('main_01_collecting_resources', 'main', 'Collecting Resources', 'harvest_resources', 'resource_warehouse', 1, 1, 450, 950, 500, '[{"type":"item","item_key":"loudspeaker","quantity":1}]', 1, true),
('main_03_tech_center', 'main', 'Level 1 Technology Center', 'build_building', 'technology_center', 1, 3, 2250, 2100, 3250, '[{"type":"item","item_key":"loudspeaker","quantity":1}]', 1, true),
('main_04_research', 'main', 'Level 1 Technological Research', 'research_tech', 'concurrent_construction', 1, 4, 500, 480, 600, '[{"type":"item","item_key":"construction_card","quantity":1}]', 1, true),
('main_05_metal_production', 'main', 'Metal Production', 'build_building', 'metal_collector', 1, 5, 425, 530, 425, '[{"type":"blueprint","blueprint_key":"super_transmission_engine"}]', 1, true),
('main_06_blueprints_1', 'main', 'Blueprints 1', 'use_blueprint', 'super_transmission_engine', 1, 6, 515, 1190, 575, '[{"type":"item","item_key":"truce_card","quantity":1}]', 1, true),
('main_07_he3_production', 'main', 'He3 Production', 'build_building', 'he3_extractor', 1, 7, 475, 400, 475, '[{"type":"blueprint","blueprint_key":"estrella"}]', 1, true),
('main_08_blueprints_2', 'main', 'Blueprints 2', 'use_blueprint', 'estrella', 1, 8, 520, 1150, 585, '[{"type":"blueprint","blueprint_key":"ship_reinforcement_facility"}]', 1, true),
('main_09_residential', 'main', 'Creating Residential Area', 'build_building', 'residential_area', 1, 9, 390, 360, 325, '[{"type":"item","item_key":"loudspeaker","quantity":1}]', 1, true),
('main_23_space_station', 'main', 'Level 2 Space Station', 'upgrade_building', 'space_station', 2, 23, 10465, 9660, 13685, '[]', 1, true),
('main_24_space_defense', 'main', 'Space Defense 1', 'build_defense', '', 1, 24, 750, 1600, 880, '[]', 1, true),
('main_25_metal_lv2', 'main', 'Level 2 Metal Collector', 'upgrade_building', 'metal_collector', 2, 25, 815, 690, 815, '[{"type":"item","item_key":"metal_mining_boost","quantity":1}]', 1, true),
('main_26_he3_lv2', 'main', 'Level 2 He3 Extractor', 'upgrade_building', 'he3_extractor', 2, 26, 730, 910, 730, '[{"type":"item","item_key":"he3_mining_boost","quantity":1}]', 1, true),
('main_27_residential_lv2', 'main', 'Level 2 Residential Area', 'upgrade_building', 'residential_area', 2, 27, 670, 620, 560, '[{"type":"item","item_key":"extra_tax","quantity":1}]', 1, true),
-- Phase 2 Main Quests (military focus)
('main_10_ship_factory', 'main', 'Level 1 Ship Factory', 'build_building', 'ship_factory', 1, 10, 600, 475, 550, '[{"type":"blueprint","blueprint_key":"typhoon"}]', 2, true),
('main_11_command_center', 'main', 'Level 1 Command Center', 'build_building', 'command_center', 1, 11, 3000, 2250, 2500, '[{"type":"blueprint","blueprint_key":"energy_shield_booster"}]', 2, true),
('main_12_recruit', 'main', 'Recruit Commanders', 'recruit_commander', '', 1, 12, 605, 1250, 685, '[{"type":"item","item_key":"revival_card","quantity":1}]', 2, true),
('main_13_design_ship', 'main', 'Design a Ship', 'create_ship_design', '', 1, 13, 585, 1235, 650, '[{"type":"item","item_key":"loudspeaker","quantity":1}]', 2, true),
('main_14_build_ships', 'main', 'Ship Building', 'build_ships', '', 1, 14, 600, 1200, 680, '[{"type":"blueprint","blueprint_key":"anti_aircraft_cannon"}]', 2, true),
('main_15_fleet', 'main', 'Build a Fleet', 'create_fleet', '', 1, 15, 610, 1300, 690, '[{"type":"item","item_key":"loudspeaker","quantity":1}]', 2, true),
('main_16_logistics', 'main', 'Wartime Logistics', 'replenish_ammo', '', 1, 16, 2000, 3000, 2000, '[{"type":"item","item_key":"truce_card","quantity":1}]', 2, true),
('main_17_weapon_research', 'main', 'Level 1 Weapon Research Center', 'build_building', 'weapon_research_center', 1, 17, 2500, 1500, 2250, '[{"type":"blueprint","blueprint_key":"starlight_missile_pod"}]', 2, true),
-- Deferred Main Quests (social/chat-dependent)
('main_02_loud_and_clear', 'main', 'Loud and Clear', 'send_message', 'world_channel', 1, 2, 460, 980, 520, '[]', 3, false),
('main_18_bigger_bags', 'main', 'Bigger Bags', 'increase_bag_slot', '', 1, 18, 620, 1310, 700, '[{"type":"item","item_key":"primary_metal_pack","quantity":1}]', 3, false),
('main_19_resource_pack', 'main', 'Resource Pack', 'use_resource_pack', '', 1, 19, 680, 1450, 750, '[{"type":"item","item_key":"galaxy_transfer","quantity":1}]', 3, false),
('main_20_growing_resources', 'main', 'Growing Resources', 'grow_comsats', '', 1, 20, 700, 1480, 780, '[{"type":"item","item_key":"healing_card","quantity":1}]', 3, false),
('main_21_adding_friends', 'main', 'Adding Friends', 'add_friend', '', 1, 21, 710, 1500, 800, '[]', 3, false),
('main_22_mail_system', 'main', 'The Mail System', 'send_message', 'mail', 1, 22, 720, 1520, 820, '[{"type":"item","item_key":"loudspeaker","quantity":1}]', 3, false),
('main_28_ship_factory_lv2', 'main', 'Level 2 Ship Factory', 'upgrade_building', 'ship_factory', 2, 28, 0, 0, 0, '[]', 2, true);

-- Set prerequisite chain (main_01 -> main_03 -> main_04 -> ... sequential)
-- Note: Skip deferred quests in Phase 1/2 chain by linking past them
-- This is done via UPDATE after INSERT to resolve forward references
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_01_collecting_resources') WHERE quest_key = 'main_03_tech_center';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_03_tech_center') WHERE quest_key = 'main_04_research';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_04_research') WHERE quest_key = 'main_05_metal_production';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_05_metal_production') WHERE quest_key = 'main_06_blueprints_1';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_06_blueprints_1') WHERE quest_key = 'main_07_he3_production';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_07_he3_production') WHERE quest_key = 'main_08_blueprints_2';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_08_blueprints_2') WHERE quest_key = 'main_09_residential';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_09_residential') WHERE quest_key = 'main_10_ship_factory';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_10_ship_factory') WHERE quest_key = 'main_11_command_center';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_11_command_center') WHERE quest_key = 'main_12_recruit';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_12_recruit') WHERE quest_key = 'main_13_design_ship';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_13_design_ship') WHERE quest_key = 'main_14_build_ships';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_14_build_ships') WHERE quest_key = 'main_15_fleet';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_15_fleet') WHERE quest_key = 'main_16_logistics';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_16_logistics') WHERE quest_key = 'main_17_weapon_research';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_17_weapon_research') WHERE quest_key = 'main_23_space_station';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_23_space_station') WHERE quest_key = 'main_24_space_defense';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_24_space_defense') WHERE quest_key = 'main_25_metal_lv2';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_25_metal_lv2') WHERE quest_key = 'main_26_he3_lv2';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_26_he3_lv2') WHERE quest_key = 'main_27_residential_lv2';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'main_27_residential_lv2') WHERE quest_key = 'main_28_ship_factory_lv2';
```

### 3.26 Seed Data: Side Quests (Phase 1)

```sql
-- Resource production side quests (Phase 1, tier 1 only - additional tiers seeded incrementally)
INSERT INTO quest_types (quest_key, category, display_name, requirement_type, requirement_target, requirement_value, chain_order, reward_metal, reward_he3, reward_gold, reward_item_json, phase, is_active) VALUES
('side_harvest_time_1', 'side', 'Harvest Time I', 'reach_production', 'metal', 2180, 1, 500, 400, 500, '[]', 1, true),
('side_harvest_time_2', 'side', 'Harvest Time II', 'reach_production', 'metal', 4360, 2, 1500, 1200, 1500, '[]', 1, true),
('side_harvest_time_3', 'side', 'Harvest Time III', 'reach_production', 'metal', 8720, 3, 4500, 3600, 4500, '[]', 1, true),
('side_gathering_he3_1', 'side', 'Gathering He3 I', 'reach_production', 'he3', 2360, 1, 400, 500, 400, '[]', 1, true),
('side_gathering_he3_2', 'side', 'Gathering He3 II', 'reach_production', 'he3', 4720, 2, 1200, 1500, 1200, '[]', 1, true),
('side_gathering_he3_3', 'side', 'Gathering He3 III', 'reach_production', 'he3', 9440, 3, 3600, 4500, 3600, '[]', 1, true),
('side_raising_morale_1', 'side', 'Raising Morale I', 'reach_production', 'gold', 2800, 1, 400, 400, 500, '[]', 1, true),
('side_raising_morale_2', 'side', 'Raising Morale II', 'reach_production', 'gold', 5600, 2, 1200, 1200, 1500, '[]', 1, true),
('side_raising_morale_3', 'side', 'Raising Morale III', 'reach_production', 'gold', 11200, 3, 3600, 3600, 4500, '[]', 1, true),
('side_plentiful_resources_1', 'side', 'Plentiful Resources I', 'reach_storage', '', 50000, 1, 1000, 1000, 1000, '[]', 1, true),
('side_plentiful_resources_2', 'side', 'Plentiful Resources II', 'reach_storage', '', 200000, 2, 3000, 3000, 3000, '[]', 1, true),
('side_plentiful_resources_3', 'side', 'Plentiful Resources III', 'reach_storage', '', 1000000, 3, 10000, 10000, 10000, '[]', 1, true);

-- Set side quest tier chains
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'side_harvest_time_1') WHERE quest_key = 'side_harvest_time_2';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'side_harvest_time_2') WHERE quest_key = 'side_harvest_time_3';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'side_gathering_he3_1') WHERE quest_key = 'side_gathering_he3_2';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'side_gathering_he3_2') WHERE quest_key = 'side_gathering_he3_3';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'side_raising_morale_1') WHERE quest_key = 'side_raising_morale_2';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'side_raising_morale_2') WHERE quest_key = 'side_raising_morale_3';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'side_plentiful_resources_1') WHERE quest_key = 'side_plentiful_resources_2';
UPDATE quest_types SET prerequisite_quest_id = (SELECT id FROM quest_types WHERE quest_key = 'side_plentiful_resources_2') WHERE quest_key = 'side_plentiful_resources_3';
```

### 3.27 Seed Data: Daily Quests

```sql
INSERT INTO quest_types (quest_key, category, display_name, requirement_type, requirement_target, requirement_value, chain_order, reward_metal, reward_he3, reward_gold, reward_item_json, phase, is_active) VALUES
('daily_login', 'daily', 'Daily Log In', 'login', '', 1, NULL, 0, 0, 0, '[]', 1, true),
('daily_collect_dues', 'daily', 'Collect Your Dues', 'harvest_resources', '', 1, NULL, 0, 0, 0, '[]', 1, true),
('daily_need_for_speed', 'daily', 'Need for Speed', 'use_speedup', 'construction', 1, NULL, 0, 0, 0, '[]', 1, true),
('daily_stockpiling', 'daily', 'Stockpiling', 'harvest_resources', 'resource_warehouse', 3, NULL, 0, 0, 0, '[]', 1, true),
('daily_donations', 'daily', 'Donations', 'donate_resources', '', 200000, NULL, 0, 0, 0, '[]', 2, true),
('daily_restricted_instances', 'daily', 'Restricted Instances', 'complete_instance', 'restricted', 2, NULL, 0, 0, 5000, '[]', 2, true);
```

### 3.28 Seed Data: Tech Types (Logistics Construction)

```sql
INSERT INTO tech_types (name, display_name, tree, max_level, prerequisites_json, base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier, base_time_seconds, time_multiplier, effects_json, description) VALUES
('concurrent_construction', 'Concurrent Construction', 'logistics_construction', 1, '[]', 0, 0, 1000, 1.0000, 20, 1.0000, '{"type":"construction_slots","per_level":1}', 'Adds 1 construction slot'),
('construction_boost', 'Construction Boost', 'logistics_construction', 10, '[{"tech":"concurrent_construction","level":1}]', 0, 0, 2400, 1.5300, 480, 2.3400, '{"type":"build_speed","per_level":1.5,"unit":"percent"}', '+1-15% building construction speed'),
('quality_materials', 'Quality Materials', 'logistics_construction', 10, '[{"tech":"construction_boost","level":3}]', 0, 0, 1200, 1.5300, 240, 2.3400, '{"type":"build_cost_reduction","per_level":1.5,"unit":"percent"}', '-1-15% building resource costs'),
('ship_building_boost', 'Ship Building Boost', 'logistics_construction', 10, '[]', 0, 0, 990, 1.5300, 198, 2.3400, '{"type":"ship_build_speed","per_level":1.5,"unit":"percent"}', '+1-15% shipbuilding speed'),
('ship_building_logistics', 'Ship Building Logistics', 'logistics_construction', 10, '[{"tech":"ship_building_boost","level":2}]', 0, 0, 1386, 1.5300, 277, 2.3400, '{"type":"ship_build_cost_reduction","per_level":1.5,"unit":"percent"}', '-1-15% ship construction resource costs'),
('sync_shipbuilding', 'Sync Shipbuilding', 'logistics_construction', 1, '[{"tech":"ship_building_logistics","level":4}]', 0, 0, 174000, 1.0000, 34800, 1.0000, '{"type":"ship_production_slots","per_level":1}', 'Adds 1 shipbuilding slot (5th production slot)'),
('repair_technology', 'Repair Technology', 'logistics_construction', 10, '[{"tech":"sync_shipbuilding","level":1}]', 0, 0, 5310, 1.5300, 1062, 2.3400, '{"type":"ship_repair_percent","per_level":1,"unit":"percent"}', '+1-10% ship repair percentage'),
('high_yield_mining', 'High Yield Mining', 'logistics_construction', 10, '[]', 0, 0, 1740, 1.5300, 348, 2.3400, '{"type":"metal_output","per_level":1,"unit":"percent"}', '+1-10% Metal output'),
('high_yield_chemistry', 'High Yield Chemistry', 'logistics_construction', 10, '[{"tech":"high_yield_mining","level":2}]', 0, 0, 2400, 1.5300, 480, 2.3400, '{"type":"he3_output","per_level":1,"unit":"percent"}', '+1-10% He3 output'),
('high_yield_investing', 'High Yield Investing', 'logistics_construction', 10, '[{"tech":"high_yield_chemistry","level":2}]', 0, 0, 3570, 1.5300, 714, 2.3400, '{"type":"gold_output","per_level":1,"unit":"percent"}', '+1-10% Gold output'),
('expand_capacity', 'Expand Capacity', 'logistics_construction', 7, '[{"tech":"high_yield_investing","level":4}]', 0, 0, 3540, 1.5300, 708, 2.3400, '{"type":"warehouse_capacity","per_level":50000,"unit":"flat"}', '+50,000-350,000 warehouse storage per level');
```

### 3.29 Seed Data: Tech Types (Ballistics Science)

```sql
INSERT INTO tech_types (name, display_name, tree, max_level, prerequisites_json, base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier, base_time_seconds, time_multiplier, effects_json, description) VALUES
('ballistics_base', 'Ballistics', 'ballistics_science', 10, '[]', 0, 0, 541, 1.5300, 117, 2.3400, '{"type":"ballistic_damage","per_level":5,"unit":"percent"}', '+5% ballistic damage per level'),
('ballistic_malice', 'Ballistic Malice', 'ballistics_science', 5, '[{"tech":"ballistics_base","level":3}]', 0, 0, 7558, 1.4300, 1428, 1.3400, '{"type":"ballistic_crit_rate","per_level":1,"unit":"percent"}', '+1% critical hit rate per level'),
('ballistic_crackdown', 'Ballistic Crackdown', 'ballistics_science', 2, '[{"tech":"ballistics_base","level":3}]', 0, 0, 29628, 1.1600, 4335, 1.3400, '{"type":"ballistic_crit_damage","per_level":10,"unit":"percent"}', '+10% critical damage per level'),
('steady_control', 'Steady Control Tech', 'ballistics_science', 5, '[{"tech":"ballistics_base","level":6},{"tech":"ballistic_malice","level":3}]', 0, 0, 26108, 1.4300, 3366, 1.3400, '{"type":"weapon_space_reduction","per_level":2,"unit":"percent"}', '-2% weapon space per level'),
('precise_ballistics', 'Precise Ballistics', 'ballistics_science', 5, '[{"tech":"ballistics_base","level":8},{"tech":"steady_control","level":3}]', 0, 0, 51100, 1.4300, 4080, 1.3400, '{"type":"ballistic_hit_rate","per_level":1,"unit":"percent"}', '+1% hit rate per level'),
('shield_penetration', 'Shield Penetration', 'ballistics_science', 1, '[{"tech":"ballistic_malice","level":5},{"tech":"ballistic_crackdown","level":1},{"tech":"precise_ballistics","level":1}]', 0, 0, 154616, 1.0000, 15300, 1.0000, '{"type":"shield_bypass","flat":15,"unit":"percent"}', '15% shield bypass damage'),
('depleted_uranium_bomb', 'Depleted Uranium Bomb', 'ballistics_science', 3, '[{"tech":"ballistic_crackdown","level":2},{"tech":"shield_penetration","level":1}]', 0, 0, 203249, 1.3300, 12750, 1.3400, '{"type":"armor_bonus","neutral":10,"light":1,"unit":"percent_per_level"}', '+10-30% vs Neutral armor, +1-3% vs Light armor'),
('fire_bomb_research', 'Fire Bomb Research', 'ballistics_science', 3, '[{"tech":"ballistic_crackdown","level":2},{"tech":"shield_penetration","level":1}]', 0, 0, 203249, 1.3300, 12750, 1.3400, '{"type":"armor_bonus","regen":10,"light":1,"unit":"percent_per_level"}', '+10-30% vs Regen armor, +1-3% vs Light armor'),
('improved_penetration', 'Improved Penetration', 'ballistics_science', 3, '[{"tech":"depleted_uranium_bomb","level":1},{"tech":"fire_bomb_research","level":1}]', 0, 0, 484776, 1.3300, 22950, 1.3400, '{"type":"shield_pen_chance","per_level":3,"light_bonus":1,"unit":"percent"}', '+1-3% vs Light armor, 3-10% shield pen chance'),
('victory_rush', 'Victory Rush', 'ballistics_science', 1, '[{"tech":"depleted_uranium_bomb","level":3},{"tech":"fire_bomb_research","level":3},{"tech":"improved_penetration","level":3}]', 0, 0, 1918521, 1.0000, 224400, 1.0000, '{"type":"range_damage","ranges":[220,180,150,120],"crit_rate":5,"crit_damage":5}', 'Range damage: 220%/180%/150%/120%, +5% crit rate/damage'),
('ballistic_scattering', 'Ballistic Scattering', 'ballistics_science', 5, '[{"tech":"ballistics_base","level":10},{"tech":"ballistic_malice","level":5},{"tech":"steady_control","level":5},{"tech":"precise_ballistics","level":3}]', 0, 0, 119765, 1.4300, 9180, 1.3400, '{"type":"scatter_damage","per_level":5,"unit":"percent"}', '5-25% scatter damage to adjacent ships'),
('improved_ballistic_scattering', 'Improved Ballistic Scattering', 'ballistics_science', 3, '[{"tech":"precise_ballistics","level":5},{"tech":"ballistic_scattering","level":3}]', 0, 0, 240487, 1.3300, 17850, 1.3400, '{"type":"scatter_rate","per_level":8,"unit":"percent"}', '+8-25% scattering rate'),
('hop_bomb_research', 'Hop Bomb Research', 'ballistics_science', 5, '[{"tech":"ballistic_scattering","level":5},{"tech":"improved_ballistic_scattering","level":2}]', 0, 0, 349930, 1.4300, 20400, 1.3400, '{"type":"scatter_weapon_chance","per_level":3,"unit":"percent"}', '3-15% chance to deal 100% weapon damage as scatter');
```

### 3.30 Seed Data: Tech Types (Ship Defense Science)

```sql
INSERT INTO tech_types (name, display_name, tree, max_level, prerequisites_json, base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier, base_time_seconds, time_multiplier, effects_json, description) VALUES
('ship_defense_base', 'Ship Defense Tech', 'ship_defense_science', 2, '[]', 0, 0, 2846, 1.1600, 3240, 1.3400, '{"type":"base_defense_stats","shield":2,"structure":2,"agility":2,"defense":2,"stability":5,"unit":"percent_per_level"}', '+2-5% base shield/structure/agility/defense, +5-10% stability'),
('shield_research', 'Shield Research', 'ship_defense_science', 5, '[{"tech":"ship_defense_base","level":1}]', 0, 0, 10392, 1.4300, 2700, 1.3400, '{"type":"base_shield","per_level":1,"unit":"percent"}', '+1-5% base shield'),
('energy_diffusion', 'Energy Diffusion Tech', 'ship_defense_science', 3, '[{"tech":"ship_defense_base","level":2},{"tech":"shield_research","level":3}]', 0, 0, 46765, 1.3300, 10800, 1.3400, '{"type":"shield_damage_reduction","per_level":1,"unit":"flat"}', 'Each shield module reduces damage by 1-3 points'),
('penetration_resistance', 'Penetration Resistance', 'ship_defense_science', 2, '[{"tech":"shield_research","level":5},{"tech":"energy_diffusion","level":1}]', 0, 0, 117576, 1.1600, 21600, 1.3400, '{"type":"shield_pen_resist","per_level":3,"unit":"percent"}', '-3-7% enemy shield penetration chance'),
('augment_shield', 'Augment Shield', 'ship_defense_science', 3, '[{"tech":"energy_diffusion","level":2},{"tech":"penetration_resistance","level":1}]', 0, 0, 248747, 1.3300, 39600, 1.3400, '{"type":"base_shield_bonus","values":[6,12,20],"unit":"percent"}', '+6-20% base shield'),
('restoration', 'Restoration', 'ship_defense_science', 2, '[{"tech":"energy_diffusion","level":3},{"tech":"augment_shield","level":2}]', 0, 0, 540000, 1.1600, 90000, 1.3400, '{"type":"shield_restore","per_level":30,"intercept":1,"unit":"percent"}', '+30-60% shield restore per round, +1-2% interception'),
('augment_absorption', 'Augment Absorption', 'ship_defense_science', 2, '[{"tech":"shield_research","level":5},{"tech":"energy_diffusion","level":1}]', 0, 0, 117576, 1.1600, 21600, 1.3400, '{"type":"shield_flat_reduction","per_level":2,"unit":"flat"}', 'Shield modules reduce damage by 2-5 points'),
('energy_conservation_defense', 'Energy Conservation', 'ship_defense_science', 3, '[{"tech":"energy_diffusion","level":2},{"tech":"augment_absorption","level":1}]', 0, 0, 248747, 1.3300, 39600, 1.3400, '{"type":"absorb_no_he3","per_level":3,"unit":"percent"}', '+3-10% chance absorb damage without He3'),
('electronic_barrier', 'Electronic Barrier', 'ship_defense_science', 2, '[{"tech":"energy_diffusion","level":3},{"tech":"energy_conservation_defense","level":2}]', 0, 0, 540000, 1.1600, 90000, 1.3400, '{"type":"reflect_damage","per_level":5,"unit":"percent"}', 'Reflect 5-10% damage before shields drop to 0'),
('damage_mitigation', 'Damage Mitigation', 'ship_defense_science', 3, '[{"tech":"augment_shield","level":3},{"tech":"restoration","level":2},{"tech":"energy_conservation_defense","level":3},{"tech":"electronic_barrier","level":2}]', 0, 0, 805152, 1.3300, 108000, 1.3400, '{"type":"absorb_double","per_level":10,"collateral_reduction":15,"unit":"percent"}', '10-30% absorb double damage, 15-45% lower collateral'),
-- Structure branch
('ship_structural_analysis', 'Ship Structural Analysis', 'ship_defense_science', 5, '[{"tech":"ship_defense_base","level":1}]', 0, 0, 10392, 1.4300, 2700, 1.3400, '{"type":"base_structure","per_level":1,"unit":"percent"}', '+1-5% base structure'),
('ship_reinforcement', 'Ship Reinforcement Tech', 'ship_defense_science', 3, '[{"tech":"ship_defense_base","level":2},{"tech":"ship_structural_analysis","level":3}]', 0, 0, 46765, 1.3300, 10800, 1.3400, '{"type":"structure_damage_reduction","per_level":1,"unit":"flat"}', 'Each structure module reduces damage by 1-3'),
('resilience', 'Resilience', 'ship_defense_science', 2, '[{"tech":"ship_structural_analysis","level":5},{"tech":"ship_reinforcement","level":1}]', 0, 0, 117576, 1.1600, 21600, 1.3400, '{"type":"structure_pen_resist","per_level":3,"unit":"percent"}', '-3-7% enemy structure penetration chance'),
('structure_improvement', 'Structure Improvement', 'ship_defense_science', 3, '[{"tech":"ship_reinforcement","level":2},{"tech":"resilience","level":1}]', 0, 0, 248747, 1.3300, 39600, 1.3400, '{"type":"base_structure_bonus","values":[6,12,20],"unit":"percent"}', '+6-20% base structure'),
('fast_repair', 'Fast Repair', 'ship_defense_science', 2, '[{"tech":"ship_reinforcement","level":3},{"tech":"structure_improvement","level":2}]', 0, 0, 540000, 1.1600, 90000, 1.3400, '{"type":"structure_restore","per_level":30,"unit":"percent"}', '+30-60% structure restore per round'),
('reaction_armor_improvement', 'Reaction Armor Improvement', 'ship_defense_science', 2, '[{"tech":"ship_structural_analysis","level":5},{"tech":"ship_reinforcement","level":1}]', 0, 0, 117576, 1.1600, 21600, 1.3400, '{"type":"structure_flat_reduction","per_level":2,"unit":"flat"}', 'Structure modules reduce damage by 2-5'),
('defense_improvement', 'Defense Improvement', 'ship_defense_science', 3, '[{"tech":"ship_reinforcement","level":2},{"tech":"reaction_armor_improvement","level":1}]', 0, 0, 248747, 1.3300, 39600, 1.3400, '{"type":"absorb_no_he3_structure","per_level":3,"unit":"percent"}', '+3-10% absorb damage without He3'),
('reflection_mastery', 'Reflection Mastery', 'ship_defense_science', 2, '[{"tech":"ship_reinforcement","level":3},{"tech":"defense_improvement","level":2}]', 0, 0, 540000, 1.1600, 90000, 1.3400, '{"type":"reflect_structure_damage","per_level":5,"unit":"percent"}', 'Reflect 5-10% damage before structure drops to 0'),
('stability_mastery', 'Stability Mastery', 'ship_defense_science', 3, '[{"tech":"structure_improvement","level":3},{"tech":"fast_repair","level":2},{"tech":"defense_improvement","level":3},{"tech":"reflection_mastery","level":2}]', 0, 0, 805152, 1.3300, 108000, 1.3400, '{"type":"absorb_double_structure","per_level":10,"collateral_reduction":15,"unit":"percent"}', '10-30% absorb double damage, 15-45% lower collateral');
```

### 3.31 Seed Data: Tech Types (Directional, Missile, Ship-Based, Planetary Defense)

```sql
-- Directional Science (15 techs)
INSERT INTO tech_types (name, display_name, tree, max_level, prerequisites_json, base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier, base_time_seconds, time_multiplier, effects_json, description) VALUES
('optics_base', 'Optics', 'directional_science', 10, '[]', 0, 0, 541, 1.5300, 117, 2.3400, '{"type":"directional_damage","per_level":5,"unit":"percent"}', '+5% directional weapon damage per level'),
('directional_malice', 'Directional Malice', 'directional_science', 5, '[{"tech":"optics_base","level":3}]', 0, 0, 7558, 1.4300, 1428, 1.3400, '{"type":"directional_crit_rate","per_level":1,"unit":"percent"}', '+1% critical strike rate per level'),
('directional_accuracy', 'Directional Accuracy', 'directional_science', 5, '[{"tech":"optics_base","level":3}]', 0, 0, 7558, 1.4300, 1428, 1.3400, '{"type":"directional_accuracy","per_level":1,"unit":"percent"}', '+1% accuracy per level'),
('eagle_eye', 'Eagle Eye', 'directional_science', 2, '[{"tech":"optics_base","level":6},{"tech":"directional_accuracy","level":1}]', 0, 0, 22724, 1.1600, 2550, 1.3400, '{"type":"steering_power","per_level":10,"unit":"percent"}', '+10% steering power per level'),
('energy_penetration', 'Energy Penetration', 'directional_science', 1, '[{"tech":"directional_accuracy","level":2},{"tech":"eagle_eye","level":1}]', 0, 0, 62584, 1.0000, 6120, 1.0000, '{"type":"hit_and_shield_pen","hit":8,"shield_pen":8,"unit":"percent"}', '+8% hit rate AND 8% shield penetration chance'),
('pierce', 'Pierce', 'directional_science', 5, '[{"tech":"optics_base","level":10},{"tech":"directional_malice","level":5},{"tech":"directional_accuracy","level":5},{"tech":"energy_penetration","level":1}]', 0, 0, 87919, 1.4300, 4947, 1.3400, '{"type":"piercing_damage","per_level":3,"unit":"percent"}', 'Piercing damage through target rows, +3% per level'),
('radiative_interference', 'Radiative Interference', 'directional_science', 5, '[{"tech":"optics_base","level":6},{"tech":"directional_accuracy","level":2}]', 0, 0, 43116, 1.4300, 9180, 1.3400, '{"type":"enemy_hit_reduction","per_level":2,"unit":"percent"}', 'Reduces enemy hit rate by 2% per level'),
('improved_pierce', 'Improved Pierce', 'directional_science', 2, '[{"tech":"pierce","level":3},{"tech":"eagle_eye","level":2}]', 0, 0, 240487, 1.1600, 17850, 1.3400, '{"type":"piercing_damage_bonus","per_level":5,"unit":"percent"}', '+5% additional piercing damage per level'),
('energy_accumulation', 'Energy Accumulation', 'directional_science', 3, '[{"tech":"eagle_eye","level":2},{"tech":"energy_penetration","level":1}]', 0, 0, 136845, 1.3300, 11985, 1.3400, '{"type":"movement_crit_bonus","per_level":2,"unit":"percent"}', '+2-6% critical bonus when fleet moves'),
('electronic_interference', 'Electronic Interference', 'directional_science', 3, '[{"tech":"directional_accuracy","level":2},{"tech":"radiative_interference","level":3}]', 0, 0, 87581, 1.3300, 11985, 1.3400, '{"type":"steering_reduction","chance":5,"per_level":3,"unit":"percent"}', '5% chance reduce enemy steering by 10%; -3-9% weapon space'),
('piercing_crit', 'Piercing Crit', 'directional_science', 1, '[{"tech":"pierce","level":5},{"tech":"improved_pierce","level":1}]', 0, 0, 391234, 1.0000, 25500, 1.0000, '{"type":"piercing_critical","enabled":true}', 'Enables critical piercing vs horizontally-aligned ships'),
('weakness_detection', 'Weakness Detection', 'directional_science', 3, '[{"tech":"energy_accumulation","level":3},{"tech":"electronic_interference","level":3}]', 0, 0, 349054, 1.3300, 20298, 1.3400, '{"type":"base_damage_and_accuracy","damage":5,"accuracy":3,"unit":"percent_per_level"}', '+5% base damage per level; 3-10% bonus accuracy'),
('particle_impact', 'Particle Impact Tech', 'directional_science', 3, '[{"tech":"improved_pierce","level":2},{"tech":"piercing_crit","level":1}]', 0, 0, 646368, 1.3300, 40800, 1.3400, '{"type":"enemy_attack_reduction","per_level":10,"duration":2,"unit":"percent"}', 'Reduces enemy attack power by 10-30%, lasts 2 rounds'),
('magnetic_impact', 'Magnetic Impact', 'directional_science', 3, '[{"tech":"directional_accuracy","level":5},{"tech":"radiative_interference","level":5},{"tech":"weakness_detection","level":2}]', 0, 0, 996120, 1.3300, 96900, 1.3400, '{"type":"ignore_agility","per_level":8,"chance":10,"unit":"percent"}', 'Ignores 8-25% enemy agility with 10-30% chance'),
('dynamic_impairment', 'Dynamic Impairment', 'directional_science', 3, '[{"tech":"particle_impact","level":3},{"tech":"weakness_detection","level":3},{"tech":"magnetic_impact","level":3}]', 0, 0, 1431605, 1.3300, 124950, 1.3400, '{"type":"ignore_defense","per_level":8,"mov_reduction":1,"chance":6,"unit":"percent"}', 'Ignores 8-25% defense; reduces movement 1-3');

-- Missile Science (14 confirmed techs)
INSERT INTO tech_types (name, display_name, tree, max_level, prerequisites_json, base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier, base_time_seconds, time_multiplier, effects_json, description) VALUES
('missile_theory', 'Missile Theory', 'missile_science', 10, '[]', 0, 0, 541, 1.5300, 117, 2.3400, '{"type":"missile_damage","per_level":4,"unit":"percent"}', '+4% missile damage per level'),
('missile_accuracy', 'Missile Accuracy', 'missile_science', 5, '[{"tech":"missile_theory","level":3}]', 0, 0, 7558, 1.4300, 1428, 1.3400, '{"type":"missile_hit_rate","per_level":2,"unit":"percent"}', '+2% hit rate per level'),
('cruise_dynamics', 'Cruise Dynamics', 'missile_science', 2, '[{"tech":"missile_theory","level":3}]', 0, 0, 11603, 1.1600, 3366, 1.3400, '{"type":"steering_power","per_level":10,"unit":"percent"}', '+10% steering power per level'),
('missile_research', 'Missile Research', 'missile_science', 3, '[{"tech":"missile_theory","level":6},{"tech":"cruise_dynamics","level":1}]', 0, 0, 38831, 1.3300, 7446, 1.3400, '{"type":"missile_damage_and_pen","damage":3,"pen":1,"unit":"percent_per_level"}', '+3-10% base damage, +1-5% shield pen chance'),
('missile_elusion', 'Missile Elusion', 'missile_science', 3, '[{"tech":"cruise_dynamics","level":2},{"tech":"missile_research","level":1}]', 0, 0, 78750, 1.3300, 9690, 1.3400, '{"type":"intercept_reduction","per_level":3,"hit_bonus":3,"unit":"percent"}', '-3-10% interception rate, +3-9% hit rate'),
('missile_space_optimization', 'Missile Space Optimization', 'missile_science', 4, '[{"tech":"missile_research","level":2},{"tech":"missile_elusion","level":2}]', 0, 0, 233515, 1.3300, 16830, 1.3400, '{"type":"weapon_space_reduction","per_level":5,"unit":"percent"}', '-5-20% weapon space per level'),
('multidirectional_assault', 'Multidirectional Assault', 'missile_science', 5, '[{"tech":"missile_theory","level":10},{"tech":"missile_accuracy","level":3},{"tech":"cruise_dynamics","level":2},{"tech":"missile_research","level":3}]', 0, 0, 80796, 1.4300, 10200, 1.3400, '{"type":"scatter_all","per_level":6,"unit":"percent"}', 'Scatters 6-30% damage across ALL enemy ships'),
('nuclear_radiation', 'Nuclear Radiation Research', 'missile_science', 5, '[{"tech":"missile_accuracy","level":4},{"tech":"multidirectional_assault","level":1}]', 0, 0, 112916, 1.4300, 8160, 1.3400, '{"type":"damage_taken_increase","per_level":2,"chance":4,"unit":"percent"}', '2-10% increase to target damage taken, 4-20% chance'),
('break_armor', 'Break Armor', 'missile_science', 1, '[{"tech":"multidirectional_assault","level":2},{"tech":"nuclear_radiation","level":3}]', 0, 0, 287437, 1.0000, 25500, 1.0000, '{"type":"armor_damage_bonus","flat":5,"unit":"percent"}', '+5% damage to all armor types'),
('energy_conservation_missile', 'Energy Conservation', 'missile_science', 4, '[{"tech":"multidirectional_assault","level":3},{"tech":"nuclear_radiation","level":5},{"tech":"break_armor","level":1}]', 0, 0, 428575, 1.3300, 30600, 1.3400, '{"type":"he3_cost_reduction","per_level":4,"chance_for_half":true,"unit":"percent"}', '4-18% chance to reduce He3 cost by 50%'),
('shrapnel_research', 'Shrapnel Research', 'missile_science', 2, '[{"tech":"missile_accuracy","level":5},{"tech":"multidirectional_assault","level":3}]', 0, 0, 157172, 1.1600, 15810, 1.3400, '{"type":"scatter_bonus","per_level":2,"unit":"percent"}', '+2-4% more scattering damage'),
('exaltation', 'Exaltation', 'missile_science', 4, '[{"tech":"multidirectional_assault","level":5},{"tech":"shrapnel_research","level":1}]', 0, 0, 224495, 1.3300, 15555, 1.3400, '{"type":"scatter_vs_low_structure","per_level":20,"he3_reduction":2,"unit":"percent"}', '+20-80% scatter vs lower-structure fleets'),
('suppression', 'Suppression', 'missile_science', 4, '[{"tech":"shrapnel_research","level":2},{"tech":"exaltation","level":4}]', 0, 0, 484776, 1.3300, 22950, 1.3400, '{"type":"scatter_vs_high_structure","per_level":3,"crit":3,"unit":"percent"}', '+3-12% scatter vs higher-structure fleets, +3-12% crit'),
('missile_concussion', 'Missile Concussion', 'missile_science', 1, '[{"tech":"energy_conservation_missile","level":2},{"tech":"missile_research","level":3},{"tech":"missile_elusion","level":3},{"tech":"missile_space_optimization","level":4}]', 0, 0, 1351976, 1.0000, 178500, 1.0000, '{"type":"knockback","chance":25,"distance":2,"per_round":true}', '25% chance to push target back 2 spaces');

-- Ship-Based Science (core 10 techs with confirmed costs)
INSERT INTO tech_types (name, display_name, tree, max_level, prerequisites_json, base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier, base_time_seconds, time_multiplier, effects_json, description) VALUES
('fighter_weapons_theory', 'Fighter Weapons Theory', 'ship_based_science', 10, '[]', 0, 0, 541, 1.5300, 117, 2.3400, '{"type":"fighter_damage","per_level":3,"unit":"percent"}', '+3% fighter weapon damage per level'),
('reconnaissance', 'Reconnaissance', 'ship_based_science', 2, '[{"tech":"fighter_weapons_theory","level":3}]', 0, 0, 7558, 1.1600, 1428, 1.3400, '{"type":"steering_power","per_level":10,"unit":"percent"}', '+10% steering power per level'),
('thruster_optimization', 'Thruster Optimization', 'ship_based_science', 5, '[{"tech":"fighter_weapons_theory","level":3}]', 0, 0, 7558, 1.4300, 1428, 1.3400, '{"type":"intercept_reduction","per_level":1,"unit":"percent"}', '-1% intercept rate per level'),
('navigation', 'Navigation', 'ship_based_science', 5, '[{"tech":"fighter_weapons_theory","level":3}]', 0, 0, 7558, 1.4300, 1428, 1.3400, '{"type":"fighter_hit_rate","per_level":1,"unit":"percent"}', '+1% hit rate per level'),
('fuel_optimization', 'Fuel Optimization', 'ship_based_science', 5, '[{"tech":"fighter_weapons_theory","level":6},{"tech":"reconnaissance","level":1}]', 0, 0, 26108, 1.4300, 3366, 1.3400, '{"type":"he3_cost_reduction","per_level":1,"unit":"percent"}', '-1-5% He3 costs'),
('fighter_mastery', 'Fighter Mastery', 'ship_based_science', 1, '[{"tech":"fighter_weapons_theory","level":6},{"tech":"navigation","level":3}]', 0, 0, 51100, 1.0000, 4080, 1.0000, '{"type":"base_attack","flat":5,"unit":"percent"}', '+5% base attack power'),
('fighter_tech_upgrades', 'Fighter Tech Upgrades', 'ship_based_science', 3, '[{"tech":"reconnaissance","level":2},{"tech":"fuel_optimization","level":5}]', 0, 0, 136845, 1.3300, 11985, 1.3400, '{"type":"he3_and_shield_damage","he3":2,"shield_damage":3,"unit":"percent_per_level"}', '-2-6% He3; +3-10% damage vs shielded enemies'),
('armor_structural_analysis', 'Armor Structural Analysis', 'ship_based_science', 1, '[{"tech":"fighter_tech_upgrades","level":3}]', 0, 0, 391234, 1.0000, 25500, 1.0000, '{"type":"ship_shield_damage","flat":10,"unit":"percent"}', '+10% damage vs ships and shields'),
('fighter_weapons_efficiency', 'Fighter-based Weapons Efficiency', 'ship_based_science', 3, '[{"tech":"fighter_weapons_theory","level":10},{"tech":"armor_structural_analysis","level":1}]', 0, 0, 349054, 1.3300, 20298, 1.3400, '{"type":"reload_chance","per_level":10,"unit":"percent"}', '+10% chance per level to finish reloading after attacks'),
('ingenuity', 'Ingenuity', 'ship_based_science', 1, '[{"tech":"fighter_weapons_efficiency","level":3}]', 0, 0, 2045150, 1.0000, 255000, 1.0000, '{"type":"multi_bonus","swarm":5,"attack":5,"he3":-5,"intercept":20,"unshielded":10}', '+5% swarm/attack, -5% He3, +20% intercept, +10% unshielded dmg');

-- Planetary Defense Science (8 techs)
INSERT INTO tech_types (name, display_name, tree, max_level, prerequisites_json, base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier, base_time_seconds, time_multiplier, effects_json, description) VALUES
('energy_control', 'Energy Control', 'planetary_defense', 10, '[]', 0, 0, 2500, 2.0000, 500, 2.0000, '{"type":"defense_cost_reduction","per_level":1,"unit":"percent"}', '-1-10% resource costs for defensive structures'),
('rapid_defense_buildup', 'Rapid Defense Buildup', 'planetary_defense', 10, '[{"tech":"energy_control","level":1}]', 0, 0, 3000, 2.0000, 600, 2.0000, '{"type":"defense_build_speed","per_level":1,"unit":"percent"}', '+1-10% defense construction speed'),
('defense_enhancement', 'Defense Enhancement', 'planetary_defense', 10, '[{"tech":"energy_control","level":3},{"tech":"rapid_defense_buildup","level":3}]', 0, 0, 1500, 2.0000, 300, 2.0000, '{"type":"defense_value","per_level":1,"unit":"percent"}', '+1-10% defensive value of all structures'),
('emplacement_mastery', 'Emplacement Mastery', 'planetary_defense', 10, '[{"tech":"rapid_defense_buildup","level":5}]', 0, 0, 2250, 2.0000, 450, 2.0000, '{"type":"emplacement_attack","per_level":1,"unit":"percent"}', '+1-10% emplacement attack power'),
('utmost_defense_buildup', 'Utmost Defense Buildup', 'planetary_defense', 10, '[{"tech":"energy_control","level":5},{"tech":"defense_enhancement","level":5},{"tech":"emplacement_mastery","level":3}]', 0, 0, 1750, 2.0000, 350, 2.0000, '{"type":"max_defense_structures","per_level":1,"unit":"percent"}', '+1-10% max defensive structures'),
('range_extension', 'Range Extension', 'planetary_defense', 2, '[{"tech":"rapid_defense_buildup","level":8},{"tech":"emplacement_mastery","level":5}]', 0, 0, 300000, 2.0000, 60000, 2.0000, '{"type":"defense_range","applies_to":["particle_cannon","anti_aircraft_gun"]}', 'Increases attack range of Particle Cannons and Anti-Aircraft Guns'),
('thor_buildup', 'Thor Buildup', 'planetary_defense', 1, '[{"tech":"energy_control","level":8},{"tech":"utmost_defense_buildup","level":5},{"tech":"range_extension","level":1}]', 0, 0, 800000, 1.0000, 160000, 1.0000, '{"type":"max_thor_cannon","flat":1}', '+1 max Thor''s Cannon allowed'),
('augment_propulsion', 'Augment Propulsion', 'planetary_defense', 2, '[{"tech":"emplacement_mastery","level":8},{"tech":"thor_buildup","level":1}]', 0, 0, 400000, 2.2500, 80000, 2.2500, '{"type":"defense_movement","per_level":1,"unit":"flat"}', '+1-2 ship movement speed when defending own planet');
```

---

## 4. Game Formulas

All formulas are from Galaxy Online 2. These are approximations derived from the research data. Where exact per-level data exists (Section 3.20, 3.21), the lookup tables should be used instead.

### 4.1 Building Cost Scaling (Approximation)

```
Cost(N) ~ BaseCost * 3.03^(N-1)

Where:
  N = target level (1-based)
  BaseCost = base_cost_metal, base_cost_he3, or base_cost_gold from building_types

Note: This is an approximation. The actual GO2 values don't follow a perfect
exponential. For Civic Center, the multiplier is ~3.03x resources per level.
Use lookup tables (Section 3.20) for exact values where available.
```

### 4.2 Building Time Scaling (Approximation)

```
Time(N) ~ BaseTime * 2.87^(N-1)

Where:
  N = target level
  BaseTime = base_time_seconds from building_types
  Result is in seconds

Note: Approximation. ~2.9x build time per level based on Civic Center data.
Use lookup tables for exact values where available.
```

### 4.3 Resource Production Scaling (Approximation)

```
Production(N) ~ BaseProduction * 1.134^(N-1)

Where:
  N = current building level
  BaseProduction = base_production_per_hour from building_types

Growth: ~13.4% increase per level.
Level 1 to Level 24: ~19.6x output increase.

Note: Use He3 Extractor and Metal Collector lookup tables for exact values.
```

### 4.4 Total Planet Production Rate

```
PlanetRate(resource) = SUM(Production(level_i)) for each building_i of that resource type

Example: 8x He3 Extractors at level 10 = 8 * 2,164 = 17,312 He3/hr
```

### 4.5 Resource Collection (Idle Accumulation)

```
Collected = min(
    Rate * HoursElapsed,
    StorageCapacity - CurrentStored
)

Where:
  Rate = metal_per_hour, he3_per_hour, or gold_per_hour
  HoursElapsed = (now - last_collected_at) in fractional hours
  StorageCapacity = storage_capacity from resources table
  CurrentStored = current metal, he3, or gold
```

### 4.6 Research Cost Scaling

**All tech research costs Gold only** -- no Metal or He3 for any of the 7 science trees.

**Base techs (Lv 1-10) across all weapon trees share identical cost scaling:**

| Level | Gold Cost | Base Time |
|-------|-----------|-----------|
| 1 | 541 | 00:01:57 |
| 2 | 628 | 00:02:38 |
| 3 | 959 | 00:06:08 |
| 4 | 1,467 | 00:14:21 |
| 5 | 2,244 | 00:33:35 |
| 6 | 3,434 | 01:18:37 |
| 7 | 5,253 | 03:03:59 |
| 8 | 8,036 | 07:10:34 |
| 9 | 12,294 | 16:47:41 |
| 10 | 18,808 | 39:18:21 |

```
Approximate formulas (base techs across all weapon trees):

GoldCost(level) ~= 541 * 1.53^(level-1)
BaseTime(level) ~= 117 * 2.34^(level-1)   // in seconds

Average multipliers per level:
  Cost:  ~1.53x
  Time:  ~2.34x

Higher-tier techs have different base costs but similar exponential
scaling within their own level ranges.
```

**Estimated total Gold investment per tree (all techs maxed):**

| Tree | Approx Total Gold |
|------|-------------------|
| Logistics Construction | ~12,000,000 |
| Ballistics | ~6,500,000 |
| Directional | ~8,500,000 |
| Missile | ~8,000,000 |
| Ship-Based | ~14,000,000 |
| Ship Defense | ~8,000,000 |
| Planetary Defense | ~12,000,000 |

### 4.7 Research Time Reduction (Technology Center)

```
ResearchTimeReduction = TechnologyCenterLevel * 0.03

Max reduction: 36% at Level 12

EffectiveResearchTime = BaseResearchTime * (1 - ResearchTimeReduction)
```

### 4.8 Combat: Hit Chance

```
HitChance = BaseAccuracy
    + (AttackerSteering * 0.04)
    + (AttackerAccuracy / 12)
    - (DefenderAgility * 0.04)
    - DefenderCommanderDodge

Clamped to [5%, 95%] (assumed; [NEEDS RESEARCH: exact min/max hit chance])
```

### 4.9 Combat: Ships Destroyed

```
ShipsDestroyed = floor(RemainingDamage / (ShipStructure * StabilityPercent))

Where:
  RemainingDamage = damage after shields are depleted
  ShipStructure = individual ship's structure HP
  StabilityPercent = 1.0 + stability_bonus (from tech/modules)
```

### 4.10 Combat: Scatter Damage

```
ScatterDamage = WeaponDamage * ScatterPercent

Where:
  ScatterPercent = sum of scatter bonuses from tech research
    - Ballistic Scattering: 5-25%
    - Improved Ballistic Scattering: +8-25%
    - Hop Bomb: 3-15% chance for 100% weapon damage as scatter
  Scatter damage bypasses ALL defenses (shields, structure bonuses, armor)
  Applied to adjacent stacks in the 3x3 grid
  Missile Exaltation: up to 432% unpreventable bonus damage
```

### 4.11 Combat: Effective Stack

```
EffectiveStack = BaseStack + CommanderStarRankBonus

BaseStack:
  Frigate    = 1,100
  Cruiser    = 1,000
  Battleship =   900
```

[NEEDS RESEARCH: exact Star Rank to Effective Stack bonus formula/table]

### 4.12 Combat: Combat Readiness

```
CombatReadiness = min(ShipsInStack, EffectiveStack) / EffectiveStack * 100%

Only ships up to the Effective Stack count can attack each round.
Example: 900 Battleships with EffectiveStack=900 = 100% readiness.
Example: 450 Battleships with EffectiveStack=900 = 50% readiness (only 450 attack).
```

### 4.13 Combat: Type Advantage

```
DamageModifier:
  Frigate vs Battleship:  +5%
  Cruiser vs Frigate:     +5%
  Battleship vs Cruiser:  +5%
  Reverse matchups:       -5%
  Same type:               0%
```

### 4.14 Combat: Defense Power

```
DefensePower = StackSize * (TotalShields + TotalStructure)
Full readiness requires 3,000 ships per stack
```

### 4.15 PvP Loot

```
Loot = DefenderUncollectedResources * 0.20

Resources stored in the warehouse (already collected) are NOT lootable.
```

### 4.16 Corp Donation

```
ContributionPoints = floor(ResourcesDonated / 10000)
MallPoints: 1 MP = 1 Contribution Point
MaxDaily: 200 points (2,000,000 resources)
```

### 4.17 RBP Bonus Scaling

```
Levels 1-10:  5% + (Level * 0.5%)
Levels 11-20: Previous + (Level * 1%)
Levels 21-30: Previous + (Level * 1.5%)
Levels 31-40: Previous + (Level * 2%)
...continues...
Level 100: 280% total bonus

Bonuses apply to: Resource production, Research speed, Shipbuilding speed
```

### 4.18 PPC Interception

```
InterceptChance = 0.55 per PPC module per incoming attack
```

### 4.19 EOS Phase Shift

```
With Damage Mitigation tech:
  30% chance to absorb double damage (only on effective stack)
```

---

## 5. API Endpoints

All endpoints are prefixed with `/api`. Authentication is via JWT Bearer token in the `Authorization` header.

### 5.1 Auth

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/auth/guest` | Create anonymous guest account; returns JWT | No |
| `POST` | `/api/auth/refresh` | Refresh an expiring JWT | Yes |

**POST /api/auth/guest**
```json
// Request: (empty body)
// Response 201:
{
    "token": "eyJhbG...",
    "player": {
        "id": "uuid",
        "anonymous_id": "guest_abc123",
        "level": 1,
        "created_at": "2026-02-05T00:00:00Z"
    }
}
```

### 5.2 Player

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/player/me` | Get current player profile | Yes |
| `PUT` | `/api/player/me` | Update player profile (username) | Yes |

### 5.3 Planets

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/planets` | List all planets owned by player | Yes |
| `GET` | `/api/planets/:id` | Get planet details (ground + space base) | Yes |

### 5.4 Buildings

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/planets/:id/buildings` | List all buildings on planet (ground + space) | Yes |
| `POST` | `/api/planets/:id/buildings` | Construct a new building | Yes |
| `POST` | `/api/planets/:id/buildings/:buildingId/upgrade` | Start upgrading a building | Yes |
| `POST` | `/api/planets/:id/buildings/:buildingId/cancel` | Cancel an in-progress upgrade | Yes |

**POST /api/planets/:id/buildings/:buildingId/upgrade**
```json
// Request: (empty body; building ID in URL)
// Response 200:
{
    "building": {
        "id": "uuid",
        "building_type": 1,
        "level": 3,
        "is_upgrading": true,
        "upgrade_finish_at": "2026-02-05T01:30:00Z"
    },
    "resources": {
        "metal": 45000,
        "he3": 32000,
        "gold": 51000
    }
}
```

### 5.5 Resources

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/planets/:id/resources` | Get current resources and production rates | Yes |
| `POST` | `/api/planets/:id/resources/collect` | Collect accumulated resources from warehouse | Yes |

**POST /api/planets/:id/resources/collect**
```json
// Request: (empty body)
// Response 200:
{
    "collected": {
        "metal": 12500,
        "he3": 8300,
        "gold": 15200
    },
    "resources": {
        "metal": 62500,
        "he3": 48300,
        "gold": 75200,
        "metal_per_hour": 8640,
        "he3_per_hour": 9440,
        "gold_per_hour": 11200,
        "storage_capacity": 500000,
        "last_collected_at": "2026-02-05T12:00:00Z"
    }
}
```

### 5.6 Research

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/research` | List all research progress for player | Yes |
| `GET` | `/api/research/trees` | Get all 7 tech trees with prerequisites and costs | Yes |
| `GET` | `/api/research/trees/:tree` | Get single tree (e.g., `ballistics_science`) | Yes |
| `POST` | `/api/research/start` | Start researching a technology | Yes |
| `POST` | `/api/research/cancel` | Cancel in-progress research | Yes |
| `POST` | `/api/research/speedup` | Spend vouchers/MP to accelerate research | Yes |

**POST /api/research/start**
```json
// Request:
{ "tech_type_id": 5 }

// Response 200:
{
    "technology": {
        "id": "uuid",
        "tech_type": 5,
        "level": 2,
        "is_researching": true,
        "research_finish_at": "2026-02-05T03:45:00Z"
    },
    "resources": {
        "metal": 50000,
        "he3": 42000,
        "gold": 60000
    }
}
```

**POST /api/research/speedup**
```json
// Request:
{ "tech_type_id": 5, "speedup_minutes": 30 }

// Response 200:
{
    "technology": {
        "id": "uuid",
        "tech_type": 5,
        "is_researching": true,
        "research_finish_at": "2026-02-05T03:15:00Z"
    },
    "vouchers_spent": 3
}
// Cost: 3 vouchers/MP per 30 minutes of reduction
```

**GET /api/research/trees/:tree**
```json
// Response 200:
{
    "tree": "ballistics_science",
    "techs": [
        {
            "id": 12,
            "name": "ballistics_base",
            "display_name": "Ballistics",
            "max_level": 10,
            "prerequisites": [],
            "current_level": 6,
            "is_researching": false,
            "cost_next_level": { "gold": 3434 },
            "time_next_level_seconds": 4717,
            "effects": { "type": "ballistic_damage", "per_level": 5, "unit": "percent" }
        }
    ]
}
```

### 5.7 Ship Designs

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/ship-designs` | List all player ship designs (max 20) | Yes |
| `POST` | `/api/ship-designs` | Create a new ship design (requires blueprint) | Yes |
| `PUT` | `/api/ship-designs/:id` | Update a ship design | Yes |
| `DELETE` | `/api/ship-designs/:id` | Delete a ship design | Yes |

### 5.8 Ships

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/ships` | List all player ships by design | Yes |
| `POST` | `/api/ships/build` | Start building ships in Ship Factory | Yes |

**POST /api/ships/build**
```json
// Request:
{
    "ship_design_id": "uuid",
    "quantity": 500,
    "production_slot": 1
}

// Response 200:
{
    "ships": {
        "id": "uuid",
        "ship_design_id": "uuid",
        "quantity": 1200,
        "is_building": true,
        "build_quantity": 500,
        "build_finish_at": "2026-02-05T06:00:00Z"
    },
    "resources": { "metal": 30000, "he3": 25000, "gold": 40000 }
}
```

### 5.9 Blueprints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/blueprints` | List all available blueprints in the game | Yes |
| `GET` | `/api/blueprints/mine` | List blueprints unlocked by this player | Yes |

### 5.10 Fleets

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/fleets` | List all player fleets | Yes |
| `POST` | `/api/fleets` | Create a new fleet | Yes |
| `PUT` | `/api/fleets/:id` | Update fleet composition/formation/targeting | Yes |
| `DELETE` | `/api/fleets/:id` | Disband a fleet | Yes |
| `POST` | `/api/fleets/:id/move` | Send fleet to coordinates | Yes |
| `POST` | `/api/fleets/:id/recall` | Recall a traveling fleet | Yes |
| `POST` | `/api/fleets/:id/dismiss` | Dismiss fleet (protection when offline) | Yes |

### 5.11 Combat

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/combat/attack` | Attack a target with a fleet (PvP) | Yes |
| `POST` | `/api/combat/scout` | Scout a target (single-ship fleet) | Yes |
| `GET` | `/api/combat/reports` | List combat reports (paginated) | Yes |
| `GET` | `/api/combat/reports/:id` | Get detailed combat report | Yes |

**POST /api/combat/attack**
```json
// Request:
{
    "fleet_id": "uuid",
    "target_planet_id": "uuid"
}

// Response 200:
{
    "report_id": "uuid",
    "result": "attacker_win",
    "total_rounds": 25,
    "loot": { "metal": 5000, "he3": 3200, "gold": 7800 },
    "attacker_losses": { "ships_destroyed": 120, "he3_consumed": 45000 },
    "defender_losses": { "ships_destroyed": 85, "he3_consumed": 32000 }
}
```

### 5.12 Instances (PvE)

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/instances` | List available instances by type | Yes |
| `GET` | `/api/instances/:id` | Get instance details and enemy fleets | Yes |
| `POST` | `/api/instances/:id/attempt` | Attempt an instance with a fleet | Yes |
| `GET` | `/api/instances/progress` | Get player's instance progress | Yes |

### 5.13 Commanders

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/commanders` | List all player commanders | Yes |
| `POST` | `/api/commanders/recruit` | Free recruitment (3hr cooldown) | Yes |
| `POST` | `/api/commanders/recruit/quick` | Quick recruitment (8 Mall Points) | Yes |
| `POST` | `/api/commanders/:id/merge` | Merge commander cards for star rank (Compound Center) | Yes |
| `PUT` | `/api/commanders/:id` | Update commander (equip gems, bionic chips, assign skills) | Yes |

### 5.14 Corp

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/corp` | Get player's corp info | Yes |
| `POST` | `/api/corp` | Create a new corp | Yes |
| `POST` | `/api/corp/join` | Join a corp by ID | Yes |
| `POST` | `/api/corp/leave` | Leave current corp | Yes |
| `GET` | `/api/corp/members` | List corp members | Yes |
| `POST` | `/api/corp/donate` | Donate resources to corp (contribution points) | Yes |
| `PUT` | `/api/corp/members/:id/role` | Change a member's role (leader/officer only) | Yes |
| `GET` | `/api/corp/search` | Search corps by name/tag | Yes |
| `GET` | `/api/corp/mall` | Browse corp mall inventory | Yes |
| `POST` | `/api/corp/rbp/:id/attack` | Initiate RBP attack (corp leaders only) | Yes |

### 5.15 Trading

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/trading/listings` | Browse auction house listings | Yes |
| `POST` | `/api/trading/listings` | Create a listing (sell) | Yes |
| `POST` | `/api/trading/listings/:id/buy` | Purchase a listing | Yes |
| `DELETE` | `/api/trading/listings/:id` | Cancel own listing | Yes |

### 5.16 Recycling

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/ships/recycle` | Scrap ships to recover resources (Recycling Plant) | Yes |

### 5.17 Quests

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/quests` | Get all player quests (main + side) with current status | Yes |
| `GET` | `/api/quests/main` | Get main quest chain with progress | Yes |
| `GET` | `/api/quests/side` | Get side quests with tier progress | Yes |
| `POST` | `/api/quests/:id/claim` | Claim reward for a completed quest | Yes |
| `GET` | `/api/quests/daily` | Get today's daily quest progress and tier status | Yes |
| `POST` | `/api/quests/daily/:quest_key/complete` | Mark a daily quest as completed | Yes |
| `POST` | `/api/quests/daily/claim-tier` | Claim a daily tier reward (bronze/silver/gold/diamond) | Yes |

**GET /api/quests**
```json
// Response 200:
{
    "main_quests": [
        {
            "id": "uuid",
            "quest_key": "main_01_collecting_resources",
            "display_name": "Collecting Resources",
            "status": "claimed",
            "progress_value": 1,
            "requirement_value": 1,
            "chain_order": 1
        },
        {
            "id": "uuid",
            "quest_key": "main_03_tech_center",
            "display_name": "Level 1 Technology Center",
            "status": "available",
            "progress_value": 0,
            "requirement_value": 1,
            "chain_order": 3
        }
    ],
    "side_quests": [
        {
            "id": "uuid",
            "quest_key": "side_harvest_time_1",
            "display_name": "Harvest Time I",
            "status": "in_progress",
            "progress_value": 1500,
            "requirement_value": 2180,
            "chain_order": 1,
            "category_group": "resource_production"
        }
    ],
    "current_main_quest": {
        "quest_key": "main_03_tech_center",
        "display_name": "Level 1 Technology Center",
        "description": "Build Technology Center Lv1",
        "progress_value": 0,
        "requirement_value": 1
    }
}
```

**POST /api/quests/:id/claim**
```json
// Response 200:
{
    "quest_key": "main_01_collecting_resources",
    "rewards": {
        "metal": 450,
        "he3": 950,
        "gold": 500,
        "items": [{"type": "item", "item_key": "loudspeaker", "quantity": 1}]
    },
    "next_quest_unlocked": "main_03_tech_center"
}

// Validations:
// - Quest belongs to player
// - Quest status is 'completed' (requirement met, not yet claimed)
// - For main quests: prerequisite quest must be claimed
// - Reward resources are added to player's planet resources
// - Blueprint rewards are added to player_blueprints
// - Item rewards are added to player inventory
// - Quest status transitions to 'claimed'
// - Next quest in chain transitions from 'locked' to 'available'
```

**GET /api/quests/daily**
```json
// Response 200:
{
    "date": "2026-02-06",
    "daily_points": 17,
    "quests": [
        { "quest_key": "daily_login", "display_name": "Daily Log In", "completed": true, "points": 10 },
        { "quest_key": "daily_collect_dues", "display_name": "Collect Your Dues", "completed": false, "points": 4 },
        { "quest_key": "daily_need_for_speed", "display_name": "Need for Speed", "completed": true, "points": 3 },
        { "quest_key": "daily_stockpiling", "display_name": "Stockpiling", "completed": false, "progress": 1, "required": 3, "points_per": 1, "max_points": 3 }
    ],
    "tier_rewards": [
        { "tier": "bronze", "points_required": 10, "claimed": true },
        { "tier": "silver", "points_required": 30, "claimed": false },
        { "tier": "gold", "points_required": 50, "claimed": false },
        { "tier": "diamond", "points_required": 70, "claimed": false }
    ]
}
```

**POST /api/quests/daily/claim-tier**
```json
// Request:
{
    "tier": "bronze"
}

// Response 200:
{
    "tier": "bronze",
    "reward": {
        "type": "loudspeaker",
        "quantity": 1
    }
}

// Validations:
// - Player has enough daily_points for the tier threshold
// - Tier has not already been claimed today
// - Reward is randomly rolled from the tier's loot table
// - Tier is added to tier_rewards_claimed_json
```

---

## 6. MVP Prioritization (Phase 1)

### 6.1 MVP Goal

A **minimum playable game** where a player can log in anonymously, build up a ground base and space station, produce resources, construct/upgrade buildings, and research technology. No combat, fleets, or corps in Phase 1.

### 6.2 MVP Features (Must Have)

| Priority | Feature | Tables Used | Endpoints |
|----------|---------|-------------|-----------|
| P0 | Guest authentication | `players` | `POST /api/auth/guest` |
| P0 | Auto-create homeworld (ground base + space station) | `planets`, `resources`, `buildings` | (server-side on auth) |
| P0 | View planet and buildings (ground + space) | `buildings`, `building_types` | `GET /api/planets`, `GET /api/planets/:id/buildings` |
| P0 | Construct new buildings | `buildings` | `POST /api/planets/:id/buildings` |
| P0 | Upgrade buildings (Civic Center <-> Space Station dependency) | `buildings` | `POST /api/planets/:id/buildings/:id/upgrade` |
| P0 | Resource production (idle) with Metal, He3, Gold | `resources` | `GET /api/planets/:id/resources` |
| P0 | Collect resources from Resource Warehouse | `resources` | `POST /api/planets/:id/resources/collect` |
| P1 | View all 7 research trees | `tech_types`, `technologies` | `GET /api/research/trees`, `GET /api/research` |
| P1 | Start research (Logistics Construction tree first) | `technologies` | `POST /api/research/start` |
| P1 | Player profile (set username) | `players` | `GET /api/player/me`, `PUT /api/player/me` |

### 6.3 MVP Tech Trees (Phase 1)

Only 2 research trees active in MVP:

1. **Logistics Construction** - Construction Boost, Quality Materials, High Yield Mining/Chemistry/Investing, Expanded Capacity
2. **Ship Defense Science** - Basic defensive tech (prepares for Phase 2 combat)

### 6.4 MVP Seed Data

On guest account creation, the server automatically:

1. Creates a `player` record with `anonymous_id`
2. Creates a `planet` (homeworld) at a random unoccupied position, `is_homeworld = true`
3. Creates a `resources` row with starter resources: 5,000 Metal, 5,000 He3, 10,000 Gold
4. Creates initial **ground base** buildings:
   - 1x Civic Center (Lv 1)
   - 1x Metal Collector (Lv 1)
   - 1x He3 Extractor (Lv 1)
   - 1x Residential Area (Lv 1)
   - 1x Resource Warehouse (Lv 1)
5. Creates initial **space base** building:
   - 1x Space Station (Lv 1)
6. Recalculates `resources.metal_per_hour`, `he3_per_hour`, `gold_per_hour` based on initial buildings

### 6.5 Phase Roadmap

| Phase | Features |
|-------|----------|
| **Phase 1 (MVP)** | Auth, buildings (ground + space), resources (Metal/He3/Gold), 2 tech trees |
| **Phase 2** | Ship Factory (24 lvl), Ship Design (Frigate/Cruiser/Battleship + modules), Spacedock (12 lvl), Fleet System (3x3), Blueprints (ship + module), 30 Normal Instances, Combat (8-phase) |
| **Phase 3** | Full 7 tech trees + Research UI, Weapon Research Center, Command Center, commander recruitment/merging, PvP combat, combat reports, Radar |
| **Phase 4** | Corps, RBPs, corp donations, Galactic Wars, Alliance Center |
| **Phase 5** | Trading Center, Recycling Plant, Restricted/Constellation Instances, Weapon Research Center |
| **Phase 6** | League/Championship, Honor/Champion Points shops, decorative buildings, Galaxy Transporter |

### 6.6 MVP Database Migration

Only these tables are required for MVP Phase 1:

```
players, planets, building_types, buildings, resources, tech_types, technologies,
he3_extractor_levels, metal_collector_levels
```

Remaining tables are added in later phase migrations.

---

## 7. User Flows

### 7.1 Onboarding Flow

```
1. User opens the app in browser
2. Frontend calls POST /api/auth/guest
3. Server creates:
   - Player record with anonymous_id
   - Homeworld planet at random position
   - Ground base: Civic Center, Metal Collector, He3 Extractor, Residential Area, Resource Warehouse (all Lv 1)
   - Space base: Space Station (Lv 1)
   - Resources: 5,000 Metal, 5,000 He3, 10,000 Gold
4. Server returns JWT token
5. Frontend stores token, redirects to planet view
6. Tutorial overlay (matches GO2 Development Quests):
   a. "This is your Civic Center - it controls all other buildings"
   b. "Build more Metal Collectors, He3 Extractors, and Residential Areas"
   c. "Collect resources from your Resource Warehouse regularly"
   d. "Upgrade your Civic Center to unlock higher building levels"
   e. "Keep your Space Station within 1 level of your Civic Center"
7. Tutorial tracks quest chain: Collecting Resources -> Lv1 Technology Center -> Metal Production -> He3 Production -> etc.
```

### 7.2 Daily Gameplay Loop

```
1. Login (auto via stored JWT or new guest session)
2. COLLECT: Harvest accumulated resources from Resource Warehouse
3. BUILD: Start building upgrades (remember: Civic Center <-> Space Station within 1 level)
   - Check Civic Center level prerequisites
   - Verify resource availability (Metal, He3, Gold)
   - Max 2 concurrent construction slots (or more with Construction Cards)
4. RESEARCH: Start or check research progress in Technology Center
   - Browse 7 science trees
   - Only 1 research active per tree
5. WAIT/IDLE: Resources accumulate while offline (cap at warehouse storage)
6. RETURN: Repeat from step 2
```

### 7.3 Building Upgrade Flow

```
1. Player views planet buildings list (ground base tab + space base tab)
2. Player taps a building to see details + upgrade cost
3. Frontend shows: current level, next level stats, Metal/He3/Gold cost, time
4. Player confirms upgrade
5. Frontend calls POST /api/planets/:id/buildings/:buildingId/upgrade
6. Server validates:
   a. Player owns this planet
   b. Building is not already upgrading
   c. Building level < max_level
   d. Civic Center level >= required prerequisite (civic_center_req)
   e. For Civic Center: Space Station >= current Civic Center level
   f. For Space Station: Civic Center >= current Space Station level
   g. Player has sufficient Metal, He3, and Gold
   h. Construction slot available (max 2 concurrent)
7. Server deducts resources, sets is_upgrading=true, calculates upgrade_finish_at
8. Frontend shows countdown timer
9. When timer completes:
   a. Server finalizes: level++, is_upgrading=false, upgrade_finish_at=null
   b. If resource building: recalculate planet production rates
   c. If Technology Center: apply research time reduction (3% per level)
```

### 7.4 Resource Collection Flow

```
1. Player views planet resources panel
2. Frontend shows:
   - Current stored: Metal / He3 / Gold
   - Production rate: X/hr for each
   - Pending collection: calculated from (rate * time_since_last_collect)
   - Storage capacity and fill percentage
3. Player taps "Collect" on Resource Warehouse
4. Frontend calls POST /api/planets/:id/resources/collect
5. Server calculates:
   a. elapsed = now - last_collected_at (in hours)
   b. For each resource: pending = rate_per_hour * elapsed
   c. For each resource: collected = min(pending, storage_capacity - current)
   d. Updates stored amounts, sets last_collected_at = now
6. Frontend shows collection animation + new totals
```

### 7.5 Research Flow

```
1. Player opens Technology Center panel
2. Frontend calls GET /api/research/trees to show all 7 science trees
3. Player selects a technology to research
4. Frontend shows: current level, next level effects, Gold cost (most techs cost Gold only), time
5. Player confirms
6. Frontend calls POST /api/research/start with tech_type_id
7. Server validates:
   a. Tech prerequisites met (other techs at required levels)
   b. No other research active in this tree
   c. Technology Center level sufficient
   d. Player has sufficient resources
8. Server deducts resources, sets is_researching=true, calculates research_finish_at
9. Frontend shows countdown
10. On completion: level++, apply effects (e.g., production bonuses, build speed)
```

---

## Appendix A: [NEEDS RESEARCH] Tracker

Items marked [NEEDS RESEARCH] throughout this document that need investigation:

| Section | Topic | Priority |
|---------|-------|----------|
| 2.2.2 | Alliance Center, Trading Center, Galaxy Transporter, Compound Center, Radar max levels | Medium |
| 2.2.3 | Spacedock, Weapon Research Center, Recycling Plant max levels | Medium |
| 2.2.4 | Decorative buildings: exact list, effects, costs | Low |
| 2.2.5 | Space Base building max levels, costs per level, prerequisites | High |
| ~~2.2.7~~ | ~~Technology Center levels 2-4, 6-11 exact costs~~ | ~~RESOLVED (Section 2.2.7 complete Lv1-12 table)~~ |
| 2.2.8 | Command Center levels 6-12 exact costs and cooldowns | Medium |
| ~~2.3.4~~ | ~~Logistics Construction full tech tree with prerequisites, costs, max levels~~ | ~~RESOLVED (Section 2.3.4 full tree + summary table + seed data 3.28)~~ |
| ~~2.3.5~~ | ~~Directional Science, Missile Science, Ship-Based Science, Planetary Defense full tech trees~~ | ~~RESOLVED (Sections 2.3.5-2.3.8 + seed data 3.31)~~ |
| 2.4.2 | Complete module list with stats, space requirements, and costs | High |
| 2.4.3 | Complete ship hull table with base stats for all hull types | High |
| 2.4.5 | Exact grid positions activated for each formation | Medium |
| 2.4.8 | Ship Factory speed bonus per level table | Medium |
| 2.5.4 | Full armor type damage modifier matrix with exact percentages | Medium |
| 2.6.3 | Exact Star Rank to Effective Stack bonus table | High |
| 2.6.4 | Full skill list with effects | Medium |
| 2.6.5 | Complete gem types, bionic chip types, and their stat values | Low |
| 2.8.5 | Corp level requirements, max level, member limits per level | Medium |
| 2.9.2 | Instance list with levels, enemy compositions, and reward tables | Low |
| ~~2.10.1~~ | ~~Complete quest chain list beyond initial quests~~ | ~~RESOLVED (Section 2.10.1-2.10.6, 3.22-3.27, 5.17)~~ |
| 3.21 | Residential Area (Gold) per-level production table | High |
| 3.21 | Resource Warehouse per-level storage capacity table | High |
| 3.21 | Metal Collector per-level cost and build time table | Medium |
| ~~4.6~~ | ~~Exact research cost scaling multiplier per level~~ | ~~RESOLVED (Section 4.6 formula: ~1.53x cost, ~2.34x time + per-level table)~~ |
| 4.8 | Exact min/max hit chance clamp values | Medium |
| 4.11 | Exact Star Rank to Effective Stack bonus formula/table | High |
| 2.3.5 | Ballistics: Demolition Warhead, Range Extension, Artillery Specialization full details | Low |
| 2.3.6 | Missile: Rapid Loading and Perfect Storm full details | Low |
| 2.3.7 | Ship-Based: remaining tech costs (Reconnaissance through Fortune) | Low |
| 2.3.8 | Planetary Defense: per-level costs for Defense Enhancement, Emplacement Mastery, Utmost Defense | Low |
| 2.3.10 | Weapon Research Center Lv 2-5, 7-11 build costs | Medium |
| 2.3.4 | Logistics: Expand Capacity exact max level (7+ confirmed, exact cap unknown) | Low |

## Appendix B: He3 Extractor Per-Level Reference

| Level | He3/hr | Metal Cost | He3 Cost | Gold Cost | Build Time | Civic Req |
|-------|--------|-----------|----------|-----------|-----------|-----------|
| 1 | 1,180 | 95 | 80 | 95 | 0:00:40 | 1 |
| 2 | 1,215 | 163 | 138 | 163 | 0:01:09 | 1 |
| 3 | 1,264 | 283 | 238 | 283 | 0:01:59 | 2 |
| 4 | 1,327 | 492 | 414 | 492 | 0:03:27 | 2 |
| 5 | 1,407 | 861 | 725 | 861 | 0:06:02 | 3 |
| 6 | 1,505 | 1,515 | 1,276 | 1,515 | 0:10:38 | 3 |
| 7 | 1,626 | 2,681 | 2,258 | 2,681 | 0:18:49 | 4 |
| 8 | 1,772 | 4,773 | 4,019 | 4,773 | 0:33:30 | 4 |
| 9 | 1,949 | 8,544 | 7,195 | 8,544 | 0:59:57 | 5 |
| 10 | 2,164 | 15,379 | 12,950 | 15,379 | 1:47:55 | 5 |
| 11 | 2,423 | 27,835 | 23,440 | 27,835 | 3:15:20 | 6 |
| 12 | 2,738 | 50,660 | 42,661 | 50,660 | 5:55:31 | 6 |
| 13 | 3,122 | 92,708 | 78,070 | 92,708 | 10:50:35 | 7 |
| 14 | 3,590 | 170,583 | 143,649 | 170,583 | 19:57:04 | 7 |
| 15 | 4,164 | 315,579 | 265,750 | 315,579 | 36:54:35 | 8 |
| 16 | 4,872 | 586,976 | 494,296 | 586,976 | 68:39:08 | 8 |
| 17 | 5,749 | 1,097,645 | 924,333 | 1,097,645 | 128:22:47 | 9 |
| 18 | 6,842 | 2,063,573 | 1,737,746 | 2,063,573 | 241:21:13 | 9 |
| 19 | 8,210 | 3,900,154 | 3,284,340 | 3,900,154 | 456:09:30 | 10 |
| 20 | 9,934 | 7,410,292 | 6,240,246 | 7,410,292 | 866:42:03 | 10 |
| 21 | 12,120 | 14,153,658 | 11,918,870 | 14,153,658 | 1655:23:55 | 11 |
| 22 | 14,907 | 27,175,024 | 22,884,230 | 27,175,024 | 3178:21:55 | 11 |
| 23 | 18,485 | 52,447,796 | 44,166,565 | 52,447,796 | 6134:14:42 | 12 |
| 24 | 23,106 | 101,748,723 | 85,683,135 | 101,748,723 | 11900:26:08 | 12 |

## Appendix C: Metal Collector Per-Level Production

| Level | Metal/hr | Level | Metal/hr |
|-------|---------|-------|---------|
| 1 | 1,080 | 13 | 2,857 |
| 2 | 1,112 | 14 | 3,286 |
| 3 | 1,157 | 15 | 3,812 |
| 4 | 1,215 | 16 | 4,459 |
| 5 | 1,288 | 17 | 5,262 |
| 6 | 1,378 | 18 | 6,262 |
| 7 | 1,488 | 19 | 7,514 |
| 8 | 1,622 | 20 | 9,092 |
| 9 | 1,784 | 21 | 11,093 |
| 10 | 1,980 | 22 | 13,644 |
| 11 | 2,218 | 23 | 16,919 |
| 12 | 2,506 | 24 | 21,148 |

## Appendix D: Combat Resolution Pseudocode

```
FOR EACH ROUND (min 20, max 99 rounds):
  FOR EACH ATTACKING STACK (ordered by commander speed):
    1. Count attacks = weapon_modules * min(ships, effective_stack) * hit_chance
    2. Defender PPC interceptors roll (55% per PPC module per attack)
    3. Surviving attacks roll damage [weapon_min, weapon_max]
    4. Apply critical hits (electron stat + tech bonuses)
    5. Apply type advantage (+/-5% Frigate/Cruiser/Battleship triangle)
    6. Apply weapon expertise modifier (S:+30%, A:+10%, B:0%, C:-10%, D:-30%)
    7. Apply ship expertise modifier (S:+10%/-10%, A:+5%/-10%, etc.)
    8. Non-EOS shields reduce damage (Heat Diffusion, Daedalus, Energy Armor)
    9. Shield penetration check (tech-based %)
    10. EOS Phase Shift: 30% chance to absorb double damage (Damage Mitigation tech)
    11. Remaining damage applied to combined shield HP pool (3,000 ships * shield_per_ship)
    12. Overflow -> ships_destroyed = floor(overflow / (structure * stability))
    13. Scatter damage = weapon_dmg * scatter% -> hits adjacent stacks (NO defense)

  CHECK: All stacks on one side destroyed -> battle ends early

AFTER BATTLE:
  PvP: Winner gets 20% of loser's uncollected resources
  Instances: Treasure Box roll (10% blueprint chance)
  Generate combat_report with full round-by-round data
  Apply He3 consumption for ships used
  Normal/Restricted Instances: destroyed ships are permanently lost
  Trial/Constellation: no ship loss, only He3 consumed
  League/Championship: no ship loss, no He3 consumed
```

---

## 8. Phase 2: Ships, Fleets & Combat

### 8.0 Phase 2 Scope

**IN SCOPE:**
- Ship Factory building (24 levels, 5 production slots)
- Ship Design: Frigate, Cruiser, Battleship (standard hulls only)
- All modules (~396 total, hull-agnostic; initial subset for seed data)
- Module Blueprints (all categories)
- Ship Blueprints for the 3 standard hull classes
- Blueprint Research (Weapon Research Center, Levels 1-3)
- Spacedock building (12 levels, PvP ship repair)
- Fleet System (3x3 grid, 27,000 max ships)
- Commander assignment to fleets (core mechanics only)
- 30 Normal Instances (PvE)
- Combat system (8-phase resolution)

**OUT OF SCOPE (deferred to later phases or removed):**
- Special Hulls (require Badge Points from Restricted Instances)
- Restricted Instances, Scenario Instances, Constellation Instances
- Commander acquisition loop (Lucky Wheel, Mall Point draws, Compound Center merging)
- Commander Cards gacha system

**PERMANENTLY OUT OF SCOPE:**
- Flagships / Humaroid Flagships (not implemented; only Frigate, Cruiser, Battleship)
- Blueprint Shreds (Flagship-only crafting material)
- Light Armor (Flagship-exclusive armor type)

---

### 8.1 Ship Factory

#### 8.1.1 Overview

The Ship Factory is the production facility for constructing ships. It is a military building in the Ground Base.

| Property | Value |
|----------|-------|
| **Max Level** | 24 |
| **Design Capacity** | 20 ship blueprints stored |
| **Production Slots** | 5 total |
| **Max Production Run** | 2,000,000 ships per slot |
| **Prerequisite** | Civic Center (level-gated) |

#### 8.1.2 Production Slot Unlocking

| Slot | Requirement |
|------|-------------|
| Slot 1 | Ship Factory Level 1 |
| Slot 2 | Ship Factory Level 4 |
| Slot 3 | Ship Factory Level 8 |
| Slot 4 | Ship Factory Level 12 |
| Slot 5 | "Sync Shipbuilding" tech (Logistics Construction tree) |

#### 8.1.3 Speed Bonus Progression

Ship Factory level increases ship production speed.

| Level | Speed Bonus | Level | Speed Bonus |
|-------|-------------|-------|-------------|
| 1 | 1% | 13 | 30% |
| 2 | 3% | 14 | 33% |
| 3 | 5% | 15 | 36% |
| 4 | 7% | 16 | 39% |
| 5 | 9% | 17 | 42% |
| 6 | 11% | 18 | 45% |
| 7 | 14% | 19 | 48% |
| 8 | 17% | 20 | 51% |
| 9 | 20% | 21 | 54% |
| 10 | 23% | 22 | 56% |
| 11 | 25% | 23 | 58% |
| 12 | 27% | 24 | 60% |

**Formula (approximation):**
```
SpeedBonus(level) = level * 2.5 (capped at 60%)
```

#### 8.1.4 Ship Build Time

```
ShipBuildTime = BaseShipBuildTime * (1 - SpeedBonus/100) * (1 - ConstructionBoostTech/100)

BaseShipBuildTime = sum of all module build times on the design
  + 1 second per Atomic Framework
  + 1 second per Orbital Shield

Effective production time for a batch:
  BatchTime = ShipBuildTime * Quantity
```

#### 8.1.5 Ship Factory Level Data (Reference Table)

```sql
CREATE TABLE ship_factory_levels (
    level INTEGER PRIMARY KEY,
    civic_center_req INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds BIGINT NOT NULL,
    speed_bonus_pct INTEGER NOT NULL,
    production_slots INTEGER NOT NULL
);

INSERT INTO ship_factory_levels VALUES
(1,  1,  600,       450,       500,       200,      1,  1),
(2,  1,  1818,      1364,      1515,      574,      3,  1),
(3,  2,  5509,      4132,      4590,      1647,     5,  1),
(4,  2,  16692,     12519,     13910,     4727,     7,  2),
(5,  3,  50576,     37932,     42147,     13566,    9,  2),
(6,  3,  153244,    114933,    127704,    38934,    11, 2),
(7,  4,  464330,    348248,    386942,    111760,   14, 2),
(8,  4,  1406920,   1055190,   1172434,   320751,   17, 3),
(9,  5,  4264968,   3198726,   3554140,   920555,   20, 3),
(10, 5,  12924853,  9693640,   10770711,  2641993,  23, 3),
(11, 6,  39170305,  29377729,  32641855,  7582522,  25, 3),
(12, 6,  118706025, 89029519,  98921822,  21761838, 27, 4),
(13, 7,  359699416, 269774562, 299749513, 62454475, 30, 4),
(14, 7,  1090000000,817500000, 908333333, 179244223,33, 4),
(15, 8,  3303000000,2477250000,2752500000,514370880,36, 4),
(16, 8,  10008000000,7506000000,8340000000,1476244425,39,4),
(17, 9,  30324000000,22743000000,25270000000,4236821340,42,4),
(18, 9,  91882000000,68911500000,76568333333,12159717246,45,4),
(19, 10, 278400000000,208800000000,232000000000,34898378376,48,4),
(20, 10, 843600000000,632700000000,702870000000,100158345738,51,4),
(21, 11, 2556100000000,1917075000000,2130083333333,287454452118,54,4),
(22, 11, 7744983000000,5808737250000,6454152500000,824914057499,56,4),
(23, 12, 23467297000000,17600473000000,19556081111111,2367503784943,58,4),
(24, 12, 71105509000000,53329132000000,59254591111111,6792715862065,60,4);
```

**Note:** Levels 13+ costs are extrapolated from the base formula `Cost(N) = BaseCost * 3.03^(N-1)`. Exact GO2 wiki values should replace these when available.

---

### 8.2 Ship Design System

#### 8.2.1 Hull Types (Standard - Phase 2)

Three hull classes with rock-paper-scissors balance:

| Hull Class | Lines | Tiers | Total Hulls | Armor Types | Base Eff. Stack | Movement | Counter | Countered By |
|------------|-------|-------|-------------|-------------|-----------------|----------|---------|--------------|
| **Frigate** | 10 | I-III | 30 | Nano, Neutralizing | 1,100 | 1 (default) | Battleship (+5%) | Cruiser (-5%) |
| **Cruiser** | 10 | I-III | 30 | Chrome, Regen | 1,000 | 0 | Frigate (+5%) | Battleship (-5%) |
| **Battleship** | 5+ | I-III | 15+ | All four | 900 | 0 | Cruiser (+5%) | Frigate (-5%) |

#### 8.2.2 Frigate Hull Lines

| Hull Line | Tier I Shield/Structure/Slots | Tier II Shield/Structure/Slots | Tier III Shield/Structure/Slots | Armor |
|-----------|-------------------------------|--------------------------------|----------------------------------|-------|
| Weikes | 270 / 770 / 80 | 540 / 1,540 / 120 | 1,078 / 3,080 / 180 | Nano |
| Air Wanderer | 290 / 810 / 85 | 580 / 1,620 / 125 | 1,150 / 3,240 / 190 | Neutralizing |
| Valkyrie | 310 / 750 / 90 | 620 / 1,500 / 130 | 1,230 / 3,000 / 195 | Nano |
| GoGetter | 260 / 800 / 82 | 520 / 1,600 / 122 | 1,040 / 3,200 / 182 | Neutralizing |
| Space Hunter | 300 / 790 / 88 | 600 / 1,580 / 128 | 1,200 / 3,160 / 188 | Nano |
| Sparrow | 280 / 760 / 78 | 560 / 1,520 / 118 | 1,120 / 3,040 / 178 | Neutralizing |
| Devourer | 320 / 740 / 92 | 640 / 1,480 / 132 | 1,280 / 2,960 / 198 | Nano |
| Polymesus | 250 / 820 / 84 | 500 / 1,640 / 124 | 1,000 / 3,280 / 184 | Neutralizing |
| Cybra | 340 / 780 / 95 | 680 / 1,560 / 135 | 1,360 / 3,120 / 200 | Nano |
| Hamdar | 275 / 830 / 86 | 550 / 1,660 / 126 | 1,100 / 3,320 / 186 | Neutralizing |

**Frigate Characteristics:** Highest shields relative to structure, fastest (Movement 1 by default), lowest module capacity, +1% base critical hit bonus.

#### 8.2.3 Cruiser Hull Lines

| Hull Line | Tier I Shield/Structure/Slots | Tier II Shield/Structure/Slots | Tier III Shield/Structure/Slots | Armor |
|-----------|-------------------------------|--------------------------------|----------------------------------|-------|
| Typhoon | 505 / 2,599 / 120 | 1,010 / 5,198 / 180 | 2,020 / 10,396 / 270 | Chrome |
| Bombardier | 480 / 2,700 / 125 | 960 / 5,400 / 185 | 1,920 / 10,800 / 275 | Regen |
| Duke | 520 / 2,500 / 115 | 1,040 / 5,000 / 175 | 2,080 / 10,000 / 265 | Chrome |
| The Shuttler | 490 / 2,650 / 122 | 980 / 5,300 / 182 | 1,960 / 10,600 / 272 | Regen |
| Watchman | 530 / 2,550 / 118 | 1,060 / 5,100 / 178 | 2,120 / 10,200 / 268 | Chrome |
| Spinner | 470 / 2,750 / 128 | 940 / 5,500 / 188 | 1,880 / 11,000 / 278 | Regen |
| Wraith | 540 / 2,450 / 112 | 1,080 / 4,900 / 172 | 2,160 / 9,800 / 262 | Chrome |
| Encratos | 500 / 2,800 / 130 | 1,000 / 5,600 / 190 | 2,000 / 11,200 / 280 | Regen |
| Nicholas | 510 / 2,620 / 119 | 1,020 / 5,240 / 179 | 2,040 / 10,480 / 269 | Chrome |
| Helena | 495 / 2,680 / 124 | 990 / 5,360 / 184 | 1,980 / 10,720 / 274 | Regen |

**Cruiser Characteristics:** Balanced shields and structure, moderate module capacity, better fuel storage than Frigates.

#### 8.2.4 Battleship Hull Lines

| Hull Line | Tier I Shield/Structure/Slots | Tier II Shield/Structure/Slots | Tier III Shield/Structure/Slots | Armor |
|-----------|-------------------------------|--------------------------------|----------------------------------|-------|
| Estrella | 380 / 4,200 / 160 | 760 / 8,400 / 240 | 1,520 / 16,800 / 360 | Chrome |
| Nettle | 350 / 4,500 / 165 | 700 / 9,000 / 245 | 1,400 / 18,000 / 370 | Regen |
| Diaz | 400 / 4,100 / 155 | 800 / 8,200 / 235 | 1,600 / 16,400 / 350 | Nano |
| RV766-The Explorer | 370 / 4,400 / 170 | 740 / 8,800 / 250 | 1,480 / 17,600 / 375 | Neutralizing |
| Palenka | 360 / 4,350 / 162 | 720 / 8,700 / 242 | 1,440 / 17,400 / 365 | Chrome |

**Battleship Characteristics:** Highest structure and module capacity, greatest fuel storage, lowest shields, all four armor types available.

**Note:** Hull stats for tiers II and III are extrapolated (2x and 4x tier I stats). Exact GO2 values should replace these when confirmed. Installation slots scale ~1.5x per tier.

#### 8.2.5 Module Categories & Initial Subset

With ~396 total modules, Phase 2 seeds a representative subset (2-3 per weapon class, 1-2 per defense/auxiliary type). Full module catalog to be expanded incrementally.

**CRITICAL DESIGN RULES:**
1. **No mixing weapon classes** on the same ship -- widely regarded as ineffective
2. **Module placement order matters** -- see Section 8.2.7
3. **Modules are hull-agnostic** -- constraint is volume/installation slots, not hull type
4. **Some modules limited to 1 per ship** (e.g., Extreme Counterattack, Nano Station Warehouse)
5. **Approximately half storage to weapons, half to defense** is recommended

#### 8.2.6 Initial Module Subset (Seed Data)

**Attack Modules - Ballistic (Range 1-2, Cooldown 0):**

| Module | Tier | Damage Type | Min Dmg | Max Dmg | Volume | He3/Round | Build Time (s) |
|--------|------|-------------|---------|---------|--------|-----------|----------------|
| Rapid Fire | I | Kinetic | 12 | 18 | 8 | 2 | 3 |
| Rapid Fire | II | Kinetic | 24 | 36 | 12 | 4 | 5 |
| Rapid Fire | III | Kinetic | 48 | 72 | 18 | 8 | 8 |
| Taskmaster | I | Heat | 14 | 20 | 9 | 2 | 3 |
| Taskmaster | II | Heat | 28 | 40 | 14 | 4 | 5 |
| Taskmaster | III | Heat | 56 | 80 | 20 | 8 | 8 |
| Gatling Cannon | I | Kinetic | 16 | 22 | 10 | 3 | 4 |
| Gatling Cannon | II | Kinetic | 32 | 44 | 15 | 6 | 6 |
| Gatling Cannon | III | Kinetic | 64 | 88 | 22 | 12 | 10 |

**Attack Modules - Directional (Range 2-5, Cooldown 1):**

| Module | Tier | Damage Type | Min Dmg | Max Dmg | Volume | He3/Round | Build Time (s) |
|--------|------|-------------|---------|---------|--------|-----------|----------------|
| Cluster Laser Transmitter | I | Heat | 30 | 45 | 12 | 4 | 5 |
| Cluster Laser Transmitter | II | Heat | 60 | 90 | 18 | 8 | 8 |
| Cluster Laser Transmitter | III | Heat | 120 | 180 | 26 | 16 | 12 |
| Magneto Pulsar | I | Magnetic | 35 | 50 | 14 | 5 | 5 |
| Magneto Pulsar | II | Magnetic | 70 | 100 | 20 | 10 | 8 |
| Magneto Pulsar | III | Magnetic | 140 | 200 | 28 | 20 | 12 |

**Attack Modules - Missile (Range 5-8, Cooldown 3):**

| Module | Tier | Damage Type | Min Dmg | Max Dmg | Volume | He3/Round | Build Time (s) |
|--------|------|-------------|---------|---------|--------|-----------|----------------|
| Rocket Frame | I | Explosive | 80 | 120 | 16 | 8 | 6 |
| Rocket Frame | II | Explosive | 160 | 240 | 24 | 16 | 10 |
| Rocket Frame | III | Explosive | 320 | 480 | 34 | 32 | 15 |
| Starlight Missile Pod | I | Explosive | 90 | 135 | 18 | 10 | 7 |
| Starlight Missile Pod | II | Explosive | 180 | 270 | 26 | 20 | 11 |
| Starlight Missile Pod | III | Explosive | 360 | 540 | 36 | 40 | 16 |

**Attack Modules - Ship-Based Weapons (Range 6-10, Cooldown 4):**

| Module | Tier | Damage Type | Min Dmg | Max Dmg | Volume | He3/Round | Build Time (s) |
|--------|------|-------------|---------|---------|--------|-----------|----------------|
| Streamliner | I | Kinetic | 120 | 180 | 20 | 16 | 8 |
| Streamliner | II | Kinetic | 240 | 360 | 30 | 32 | 12 |
| Streamliner | III | Kinetic | 480 | 720 | 42 | 64 | 18 |
| Golem | I | Magnetic | 130 | 200 | 22 | 18 | 9 |
| Golem | II | Magnetic | 260 | 400 | 32 | 36 | 13 |
| Golem | III | Magnetic | 520 | 800 | 44 | 72 | 19 |

**Attack Modules - Planetary Weapons (Range 1-2, Cooldown 1):**

| Module | Tier | Damage Type | Min Dmg | Max Dmg | Volume | He3/Round | Build Time (s) |
|--------|------|-------------|---------|---------|--------|-----------|----------------|
| Lander Module | I | Siege | 50 | 75 | 10 | 4 | 4 |
| Lander Module | II | Siege | 100 | 150 | 16 | 8 | 7 |
| Lander Module | III | Siege | 200 | 300 | 24 | 16 | 11 |

**Note:** Planetary Weapons ONLY damage defensive structures, NOT ships.

**Defense Modules - Structure:**

| Module | Tier | Effect | Volume | Limit/Ship | Build Time (s) |
|--------|------|--------|--------|------------|----------------|
| Atomic Framework | - | +500 structure | 6 | Unlimited | 1 |
| Ship Reinforcement Facility | I | Each reduces damage by 1 | 8 | Unlimited | 4 |
| Ship Reinforcement Facility | II | Each reduces damage by 2 | 12 | Unlimited | 6 |
| Ship Reinforcement Facility | III | Each reduces damage by 3 | 18 | Unlimited | 9 |
| Quick Reaction Armor | I | 5% reflect damage | 10 | 1 | 5 |
| Quick Reaction Armor | II | 10% reflect damage | 15 | 1 | 8 |
| Quick Reaction Armor | III | 15% reflect damage | 22 | 1 | 12 |
| Reflective Plating | I | +3% defense | 8 | 1 | 4 |
| Reflective Plating | II | +6% defense | 12 | 1 | 7 |
| Reflective Plating | III | +10% defense | 18 | 1 | 10 |
| Energy Armor | I | +200 structure, 5 He3/round | 10 | 1 | 5 |
| Energy Armor | II | +400 structure, 10 He3/round | 16 | 1 | 8 |
| Energy Armor | III | +600 structure, 15 He3/round | 24 | 1 | 12 |
| Daedalus Control System | I | +3% structure, damage reduction | 12 | 1 | 6 |
| Daedalus Control System | II | +6% structure, damage reduction | 18 | 1 | 9 |
| Daedalus Control System | III | +10% structure, damage reduction | 26 | 1 | 13 |

**Defense Modules - Shields:**

| Module | Tier | Effect | Volume | Limit/Ship | Build Time (s) |
|--------|------|--------|--------|------------|----------------|
| Orbital Shield | - | +300 shield | 6 | Unlimited | 1 |
| Energy Shield Booster | I | +5% shield effectiveness | 8 | 1 | 4 |
| Energy Shield Booster | II | +10% shield effectiveness | 12 | 1 | 7 |
| Energy Shield Booster | III | +15% shield effectiveness | 18 | 1 | 10 |
| Particle Stun Shield | I | Kinetic shield; reduces ballistic/SBW dmg by 5 | 10 | Unlimited | 5 |
| Particle Stun Shield | II | Kinetic shield; reduces ballistic/SBW dmg by 10 | 15 | Unlimited | 8 |
| Particle Stun Shield | III | Kinetic shield; reduces ballistic/SBW dmg by 15 | 22 | Unlimited | 12 |
| Heat Diffusion Shield | I | Heat shield; reduces ballistic/directional/SBW dmg by 5 | 10 | Unlimited | 5 |
| Heat Diffusion Shield | II | Heat shield; reduces ballistic/directional/SBW dmg by 10 | 15 | Unlimited | 8 |
| Heat Diffusion Shield | III | Heat shield; reduces ballistic/directional/SBW dmg by 15 | 22 | Unlimited | 12 |
| Space-Time Magnetic Shield | I | Magnetic shield; reduces directional/SBW dmg by 5 | 10 | Unlimited | 5 |
| Space-Time Magnetic Shield | II | Magnetic shield; reduces directional/SBW dmg by 10 | 15 | Unlimited | 8 |
| Space-Time Magnetic Shield | III | Magnetic shield; reduces directional/SBW dmg by 15 | 22 | Unlimited | 12 |
| Detonator Shield | I | Explosive shield; reduces missile/SBW dmg by 5 | 10 | Unlimited | 5 |
| Detonator Shield | II | Explosive shield; reduces missile/SBW dmg by 10 | 15 | Unlimited | 8 |
| Detonator Shield | III | Explosive shield; reduces missile/SBW dmg by 15 | 22 | Unlimited | 12 |
| Shield Regenerator | I | +10% shield restore/round | 8 | 1 | 4 |
| Shield Regenerator | II | +20% shield restore/round | 12 | 1 | 7 |
| Shield Regenerator | III | +30% shield restore/round | 18 | 1 | 10 |
| EOS Phase Shift Shield | I | 10% absorb double dmg | 14 | 1 | 7 |
| EOS Phase Shift Shield | II | 20% absorb double dmg | 20 | 1 | 10 |
| EOS Phase Shift Shield | III | 30% absorb double dmg | 28 | 1 | 14 |

**Defense Modules - Air Defense:**

| Module | Tier | Effect | Volume | Limit/Ship | Build Time (s) |
|--------|------|--------|--------|------------|----------------|
| Anti-Aircraft Cannon | I | 15% intercept missiles | 8 | Unlimited | 4 |
| Anti-Aircraft Cannon | II | 25% intercept missiles | 12 | Unlimited | 7 |
| Anti-Aircraft Cannon | III | 35% intercept missiles | 18 | Unlimited | 10 |
| Powered Pulse Cannon | I | 35% intercept any attack | 10 | Unlimited | 5 |
| Powered Pulse Cannon | II | 45% intercept any attack | 15 | Unlimited | 8 |
| Powered Pulse Cannon | III | 55% intercept any attack | 22 | Unlimited | 12 |
| Extreme Counterattack | I | Reflect 20% intercepted damage | 12 | 1 | 6 |
| Extreme Counterattack | II | Reflect 35% intercepted damage | 18 | 1 | 9 |
| Extreme Counterattack | III | Reflect 50% intercepted damage | 26 | 1 | 13 |

**Auxiliary Modules - Electronic (1 of each type per ship):**

| Module | Tier | Effect | Volume | Build Time (s) |
|--------|------|--------|--------|----------------|
| Agility Booster | I | +1 agility | 6 | 3 |
| Agility Booster | II | +2 agility | 9 | 5 |
| Agility Booster | III | +3 agility | 13 | 8 |
| Infrared Scanner | I | +1 steering | 6 | 3 |
| Infrared Scanner | II | +2 steering | 9 | 5 |
| Infrared Scanner | III | +3 steering | 13 | 8 |
| ECM Booster | I | +5% dodge chance | 7 | 4 |
| ECM Booster | II | +10% dodge chance | 10 | 6 |
| ECM Booster | III | +15% dodge chance | 14 | 9 |
| Auto Target System | I | +5% hit chance | 7 | 4 |
| Auto Target System | II | +10% hit chance | 10 | 6 |
| Auto Target System | III | +15% hit chance | 14 | 9 |
| Time Dilation Module | I | +3% critical hit rate | 8 | 5 |
| Time Dilation Module | II | +6% critical hit rate | 12 | 7 |
| Time Dilation Module | III | +10% critical hit rate | 16 | 10 |

**Auxiliary Modules - Storage:**

| Module | Tier | Effect | Volume | Limit/Ship | Build Time (s) |
|--------|------|--------|--------|------------|----------------|
| Station Warehouse | - | +200 He3 storage | 4 | Unlimited | 2 |
| Nano Station Warehouse | - | +500 He3 storage, +50 module capacity | 8 | 1 | 4 |

**Auxiliary Modules - Transmission:**

| Module | Tier | Effect | Volume | Build Time (s) |
|--------|------|--------|--------|----------------|
| Super Transmission Engine | - | +1 movement | 8 | 3 |
| Team Combat Engine | I | +1 movement, +1 agility | 10 | 5 |
| Team Combat Engine | II | +1 movement, +2 agility | 14 | 7 |
| Team Combat Engine | III | +1 movement, +3 agility | 18 | 10 |
| Anti-Matter Engine | I | +2 movement, +1 agility | 12 | 6 |
| Anti-Matter Engine | II | +2 movement, +2 agility | 16 | 8 |
| Anti-Matter Engine | III | +2 movement, +3 agility | 20 | 11 |

#### 8.2.7 Module Placement Order

The order modules are placed affects effectiveness. Optimal order (first to last):

```
1.  Reflective Plating
2.  Engines (Team Combat Engine, Anti-Matter Engine, EOS Phase Shift Engine)
3.  Electronic Modules (Agility Booster, Infrared Scanner, ECM, Auto Target, Time Dilation)
4.  Maintenance Facilities (Ship Reinforcement, Shield Regenerator)
5.  Air Defense (Anti-Aircraft Cannon, Powered Pulse Cannon, Extreme Counterattack)
6.  Ship-Based Weapons (SBW) -- MUST be before Daedalus/Quick Reaction Armor
7.  Extreme Counterattack
8.  Quick Reaction Armor
9.  Daedalus Control System
10. Shield Modules (EOS Phase Shift, Detonator, Space-Time Magnetic, Heat Diffusion, Particle Stun)
11. Energy Shield Booster
12. Energy Armor
13. Non-SBW Weapons (Ballistic, Directional, Missile)
```

**Placement-Independent Modules** (function identically regardless of position):
- Nano Station Warehouse
- Station Warehouse
- Atomic Framework
- Orbital Shield
- Super Transmission Engine

#### 8.2.8 Ship Design Rules

1. **Max 20 designs** stored in Ship Factory
2. **Name restrictions:** No spaces; only periods (.), dashes (-), and underscores (_) allowed
3. **Blueprint required:** Player must own the hull blueprint before creating a design
4. **Module blueprints required:** Each module requires its blueprint to be owned
5. **Volume constraint:** Sum of all module volumes must not exceed hull installation_slots
6. **Per-ship module limits** enforced (e.g., max 1 Extreme Counterattack)
7. **Design is immutable once ships are built** -- must create a new design to change

#### 8.2.9 Ship Stats Calculation

```
Stats computed from hull base stats + equipped modules:

TotalShield    = hull.base_shield + SUM(shield_bonuses)
TotalStructure = hull.base_structure + SUM(structure_bonuses) + (atomic_framework_count * 500)
TotalDefense   = hull.base_defense + SUM(defense_bonuses)
TotalAgility   = hull.base_agility + SUM(agility_bonuses)
TotalMovement  = hull.base_movement + SUM(movement_bonuses)
TotalStorage   = hull.base_storage + SUM(storage_bonuses)
AttackPower    = SUM(module_avg_damage * module_count)
Range          = MAX(weapon_range) across equipped weapons (should be single weapon class)
BuildTime      = SUM(module_build_times) + (atomic_framework_count * 1) + (orbital_shield_count * 1)

He3PerRound = SUM(module_he3_per_round)

ResourceCost:
  metal_cost = hull_base_metal + SUM(module_metal_costs)
  he3_cost   = hull_base_he3 + SUM(module_he3_costs)
  gold_cost  = hull_base_gold + SUM(module_gold_costs)
```

---

### 8.3 Blueprints

#### 8.3.1 Ship Blueprints

Players must own a hull blueprint before designing ships with that hull. Blueprints are obtained from:

1. **Initial Quests:** Starter blueprints (Weikes, Typhoon, Estrella) given during tutorial
2. **Normal Instances:** 10% chance from Treasure Box; instance-specific drops
3. **Auction House:** Purchased from other players with Gold
4. **Galactic Trafficker:** Purchased with Corsairs' Gold

**Frigate Blueprints (10):** Weikes, Air Wanderer, Valkyrie, GoGetter, Space Hunter, Sparrow, Devourer, Polymesus, Cybra, Hamdar

**Cruiser Blueprints (10):** Typhoon, Bombardier, Duke, The Shuttler, Watchman, Spinner, Wraith, Encratos, Nicholas, Helena

**Battleship Blueprints (5):** Estrella, Nettle, Diaz, RV766-The Explorer, Palenka

#### 8.3.2 Module Blueprints

Every module requires its blueprint. Module blueprints are obtained from the same sources as ship blueprints. Each module line has a single blueprint that covers all 3 tiers.

#### 8.3.3 Blueprint Activation

Blueprints from instances arrive "unactivated" in the player's inventory. Player must activate the blueprint before using it. Unactivated blueprints can be sold on the Auction House.

#### 8.3.4 Blueprint Research (Weapon Research Center)

Module blueprints can be researched/upgraded in the Weapon Research Center building to enhance module stats.

| Research Level | Requirement | Effect | Cost Multiplier |
|----------------|-------------|--------|-----------------|
| Level 1 | Own blueprint | Base stats | 1x |
| Level 2 | Level 1 complete | +10% module stats | 3x base cost |
| Level 3 | Level 2 complete | +25% module stats | 9x base cost |

**Research Time:** Scales with module tier and research level. Weapon Research Center level reduces time.

#### 8.3.5 Starter Blueprints (Given During Tutorial)

On completing the Phase 2 tutorial chain, players receive:
- **Weikes** (Frigate) hull blueprint
- **Typhoon** (Cruiser) hull blueprint
- **Estrella** (Battleship) hull blueprint
- **Rapid Fire** module blueprint (Ballistic)
- **Energy Shield Booster** module blueprint (Shield)
- **Super Transmission Engine** module blueprint (Transmission)

---

### 8.4 Spacedock

#### 8.4.1 Overview

The Spacedock is a military building for repairing ships destroyed in **PvP combat only**. Ships lost in Normal Instances are permanently destroyed and CANNOT be repaired.

| Property | Value |
|----------|-------|
| **Max Level** | 12 |
| **Function** | Repair PvP-destroyed ships |
| **Storage** | 2 pages (10 oldest designs) |
| **Prerequisite** | Civic Center + Ship Factory (level-gated) |

#### 8.4.2 Repair Mechanics

- **Repair Percentage** = individual probability per destroyed ship
- Each destroyed ship has an independent roll based on the Spacedock's repair percentage
- **CRITICAL:** If more than 10 different ship designs are destroyed simultaneously, only the 10 oldest designs are recoverable. Newer designs beyond 10 suffer 100% losses.

#### 8.4.3 Spacedock Level Data

```sql
CREATE TABLE spacedock_levels (
    level INTEGER PRIMARY KEY,
    civic_center_req INTEGER NOT NULL,
    ship_factory_req INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL,
    repair_pct NUMERIC(5,2) NOT NULL
);

INSERT INTO spacedock_levels VALUES
(1,  1,  1,  500,        400,        450,        44,        1.00),
(2,  2,  3,  1515,       1212,       1364,       126,       2.70),
(3,  3,  5,  4590,       3672,       4131,       362,       4.50),
(4,  4,  7,  13910,      11128,      12519,      1039,      6.30),
(5,  5,  9,  42147,      33718,      37932,      2982,      8.10),
(6,  6,  11, 127704,     102163,     114933,     8559,      10.00),
(7,  7,  13, 386942,     309554,     348248,     24565,     11.80),
(8,  8,  15, 1172434,    937947,     1055190,    70501,     13.60),
(9,  9,  17, 3554140,    2843312,    3198726,    202339,    15.50),
(10, 10, 19, 10770711,   8616569,    9693640,    580714,    17.30),
(11, 11, 21, 32641855,   26113484,   29377729,   1666651,   18.60),
(12, 12, 23, 98921822,   79137458,   89029519,   10036800,  20.00);
```

#### 8.4.4 Repair Acceleration

Players can spend 10 Mall Points to reduce repair time by 10%.

---

### 8.5 Fleet System

#### 8.5.1 Fleet Grid (3x3)

```
Position Layout:
+----------------+----------------+----------------+
| Left Shoulder  |     Head       | Right Shoulder |   First Rank:  100% attack
+----------------+----------------+----------------+
| Left Flank     |  Glasshouse    | Right Flank    |   Second Rank:  90% attack
+----------------+----------------+----------------+
| Left Rear      |     Tail       | Right Rear     |   Third Rank:   75% attack
+----------------+----------------+----------------+

Grid indices (row, col):
  [0,0] [0,1] [0,2]    <- First Rank
  [1,0] [1,1] [1,2]    <- Second Rank
  [2,0] [2,1] [2,2]    <- Third Rank
```

| Property | Value |
|----------|-------|
| **Positions** | 9 stacks (3x3 grid) |
| **Ships per Stack** | Max 3,000 of a single design |
| **Max Fleet Size** | 27,000 ships (9 x 3,000) |
| **Design Rule** | One ship design per stack; cannot mix designs |

#### 8.5.2 Position Attack Power Modifiers

| Rank | Positions | Attack Power |
|------|-----------|-------------|
| First (row 0) | Head, Left Shoulder, Right Shoulder | 100% |
| Second (row 1) | Left Flank, Glasshouse, Right Flank | 90% |
| Third (row 2) | Left Rear, Tail, Right Rear | 75% |

**Glasshouse** (center [1,1]) is the single most protected position. Typically houses glass cannon stacks.

#### 8.5.3 Fleet Formations

| Formation | Active Positions | Description |
|-----------|-----------------|-------------|
| Phalanx | All 9 | Full grid; maximum firepower |
| Diamond | [0,1], [1,0], [1,1], [1,2], [2,1] | 5 stacks; protects central |
| Battle Line | [0,0]-[0,2], [1,0]-[1,2] | 6 stacks; frontal defense |
| Skirmish | [0,0], [0,2], [1,1], [2,0], [2,2] | 5 stacks; minimizes scatter |
| Tee Forward | [0,0]-[0,2], [1,1], [2,1] | 5 stacks; T-shape offense |
| Enfilade | [0,0], [1,0], [1,1], [1,2], [2,0] | 5 stacks; side-focused |
| Tee Reverse | [0,1], [1,1], [2,0]-[2,2] | 5 stacks; reversed T |

#### 8.5.4 Fleet Properties

```
FleetSpeed = MIN(ship_movement) across all stacks
  (Fleet moves at the speed of the slowest ship)

FleetAttackPower = SUM(stack_attack * position_modifier) across all stacks

FleetDefensePower = SUM(stack_ships * (stack_shield + stack_structure)) across all stacks
```

#### 8.5.5 Targeting Commands

Fleets can be ordered to prioritize targets:

| Command | Target Selection |
|---------|-----------------|
| Max Attack | Enemy fleet with highest total attack |
| Min Attack | Enemy fleet with lowest total attack |
| Max Durability | Enemy fleet with highest total durability |
| Min Durability | Enemy fleet with lowest total durability |
| Closest | Nearest enemy fleet by position |
| By Commander Rank | Enemy fleet with highest commander star rank |

#### 8.5.6 Fleet Limits

| Player Level | Max Fleets |
|-------------|------------|
| 1-10 | 2 |
| 11-20 | 3 |
| 21-30 | 4 |
| 31-40 | 5 |
| 41-50 | 6 |
| 51+ | 8 |

#### 8.5.7 Commander Assignment

Each fleet has one commander. Commander impacts:
- **Effective Stack** via Star Rank bonus
- **Weapon Expertise** (S/A/B/C/D/F grades, +30% to -30% damage modifier)
- **Ship Expertise** (S/A/B/C/D/F grades, +10%/-10% to -10%/+10% dealt/received)
- **Attributes:** Accuracy (hit chance), Dodge (evasion), Speed (attack order), Electron (crit rate/dmg)

---

### 8.6 Combat System (8-Phase Resolution)

#### 8.6.1 Combat Overview

Combat occurs when fleets engage in PvP or PvE (instances). Each combat consists of multiple rounds, each resolved through 8 sequential phases.

| Property | Value |
|----------|-------|
| **Min Rounds** | 20 + number of fleets/buildings present |
| **Max Rounds** | 99 |
| **Resolution** | Per-stack, ordered by commander Speed |
| **End Condition** | All stacks on one side destroyed OR max rounds reached |

#### 8.6.2 Effective Stack

Effective Stack determines how many ships in a stack can attack per round.

```
EffectiveStack = BaseStack + CommanderStarRankBonus

BaseStack:
  Frigate:    1,100
  Cruiser:    1,000
  Battleship:   900

AttackingShips = min(ShipsInStack, EffectiveStack)

CombatReadiness = AttackingShips / EffectiveStack * 100%
```

**Commander Star Rank Bonus (estimated):**

| Star Rank | Effective Stack Bonus |
|-----------|----------------------|
| 0 | +0 |
| 1 | +100 |
| 2 | +200 |
| 3 | +350 |
| 4 | +550 |
| 5 | +800 |
| 6 | +1,100 |
| 7 | +1,500 |
| 8 | +2,000 |
| 9 | +2,600 |
| 10 | +3,300 |
| 11 | +4,200 |
| 12 | +5,300 |
| 13 | +6,600 |
| 14 | +8,100 |
| 15 | +10,000 |

#### 8.6.3 The 8 Phases (Per Round, Per Attacking Stack)

**Phase 1: Attacker Fires**
```
Attacks = WeaponModulesPerShip * min(ShipsInStack, EffectiveStack) * HitChance

HitChance = BaseWeaponAccuracy
    + (AttackerSteering * 0.04)
    + (AttackerCommanderAccuracy / 12)
    - (DefenderAgility * 0.04)
    - (DefenderCommanderDodge * 0.02)
    + AutoTargetBonus
    - ECMBonus

HitChance clamped to [5%, 95%]

Weapons only fire when target is within range.
Weapons on cooldown skip this phase.
```

**Phase 2: Interceptors Fire**
```
FOR EACH incoming attack:
  FOR EACH defender air defense module:
    IF rand() < InterceptChance:
      Attack is intercepted (destroyed)
      IF Extreme Counterattack equipped:
        ReflectedDamage = InterceptedDamage * ReflectPercent
        Apply to attacker

PPC Intercept Chance: 55% per PPC-III (35% PPC-I, 45% PPC-II)
Anti-Aircraft: Missiles only (15%/25%/35%)
```

**Phase 3: Calculate Damage**
```
FOR EACH surviving attack:
  RawDamage = rand(WeaponMinDmg, WeaponMaxDmg)

  // Critical hit check
  CritChance = BaseCritRate + ElectronBonus + TimeDilationBonus + FrigateCritBonus
  IF rand() < CritChance:
    RawDamage *= (1.5 + CritDamageBonus)

  // Type advantage
  RawDamage *= (1.0 + TypeAdvantageModifier)  // +/-5%

  // Commander weapon expertise
  RawDamage *= WeaponExpertiseMultiplier  // S:1.3, A:1.1, B:1.0, C:0.9, D:0.7

  // Commander ship expertise (damage dealt modifier)
  RawDamage *= ShipExpertiseDealModifier  // S:1.1, A:1.05, B:1.0, C:0.95, D:0.9

  // Blueprint research bonus
  RawDamage *= (1.0 + BlueprintResearchBonus)  // +0/10/25%

  TotalRoundDamage += RawDamage
```

**Phase 4: Damage Negation**
```
FOR EACH shield module on defender:
  TotalRoundDamage -= ShieldDamageReduction
  // Particle Stun: reduces Kinetic damage
  // Heat Diffusion: reduces Heat damage
  // Space-Time Magnetic: reduces Magnetic damage
  // Detonator: reduces Explosive damage
  // Daedalus/Energy Armor: flat damage reduction

FOR EACH structure module on defender:
  TotalRoundDamage -= StructureDamageReduction

// Reflective Plating defense bonus
TotalRoundDamage *= (1.0 - DefensePercent)

// Commander ship expertise (damage received modifier)
TotalRoundDamage *= ShipExpertiseReceiveModifier  // S:0.9, A:0.9, B:1.0, C:1.1, D:1.1
```

**Phase 5: Shield Penetration**
```
// Ballistic and directional weapons can bypass shields
PenetrationChance = BasePenetrationChance + TechBonus - DefenderPenResistance

IF rand() < PenetrationChance:
  PenetratingDamage = TotalRoundDamage * PenetrationPercent  // 15% base
  // PenetratingDamage goes directly to hull (Phase 7)
  TotalRoundDamage -= PenetratingDamage
```

**Phase 6: Deal Damage to Shields**
```
TotalShieldHP = ShipsInStack * ShieldPerShip * (1 + ShieldBonuses)

// EOS Phase Shift check
IF EOSEquipped AND rand() < EOSAbsorbChance:  // 30% with tech
  DamageAbsorbed = TotalRoundDamage * 2  // absorbs double
  TotalShieldHP -= DamageAbsorbed
ELSE:
  TotalShieldHP -= TotalRoundDamage

// Shield regeneration
TotalShieldHP += ShieldRegenRate * ShipsInStack

IF TotalShieldHP < 0:
  OverflowDamage = abs(TotalShieldHP)
  // Overflow goes to Phase 7
ELSE:
  OverflowDamage = 0
```

**Phase 7: Assign Damage to Hull**
```
TotalHullDamage = OverflowDamage + PenetratingDamage

// Quick Reaction Armor reflect
IF QuickReactionEquipped:
  ReflectedDamage = TotalHullDamage * ReflectPercent
  Apply ReflectedDamage to attacker
  TotalHullDamage -= ReflectedDamage

ShipsDestroyed = floor(TotalHullDamage / (ShipStructure * StabilityMultiplier))

StabilityMultiplier = 1.0 + (StabilityTechBonus / 100)

RemainingShips = max(0, ShipsInStack - ShipsDestroyed)
```

**Phase 8: Calculate Scatter Damage**
```
ScatterDamage = WeaponDamage * ScatterPercent  // From tech research

// Scatter bypasses ALL defenses (shields, structure bonuses, armor)
FOR EACH adjacent stack on defender grid:
  AdjacentShipsDestroyed = floor(ScatterDamage / AdjacentShipStructure)
  AdjacentStack.ships -= AdjacentShipsDestroyed
```

#### 8.6.4 Armor Type vs Damage Type Matrix

| Damage Type | vs Chrome | vs Regen | vs Nano | vs Neutralizing | vs No Armor |
|-------------|-----------|----------|---------|-----------------|-------------|
| **Kinetic** | -20% | 0% | 0% | +20% | 0% |
| **Heat (Solar)** | 0% | +20% | 0% | 0% | 0% |
| **Explosive** | +20% | 0% | -20% | 0% | 0% |
| **Magnetic** | 0% | 0% | +20% | -20% | 0% |

Chrome: Strong vs Kinetic, Weak vs Explosive
Regen: Weak vs Heat
Nano: Strong vs Explosive, Weak vs Magnetic
Neutralizing: Strong vs Magnetic, Weak vs Kinetic

#### 8.6.5 Weapon Properties Summary

| Weapon Class | Range | Cooldown | Power | He3 Cost | Shield Pierce | Interceptable |
|-------------|-------|----------|-------|----------|---------------|---------------|
| Ballistic | 1-2 | 0 | Lowest | Lowest | Yes (with tech) | No |
| Directional | 2-5 | 1 | Mid-Low | Low | Yes (with tech) | No |
| Missile | 5-8 | 3 | Mid-High | High (2x) | No | Yes |
| Ship-Based | 6-10 | 4 | Highest | Highest (4x) | No | Yes |
| Planetary | 1-2 | 1 | N/A | Low | N/A | No |

#### 8.6.6 Combat Losses by Mode

| Mode | Ships Lost | He3 Lost | Repairable |
|------|-----------|----------|------------|
| Normal Instance | Yes | Yes | NO (permanent) |
| PvP Attack | Yes | Yes | Yes (Spacedock) |

#### 8.6.7 He3 Consumption

```
He3Consumed = TotalRounds * SUM(
    StackShips * SUM(module_he3_per_round for each module)
) for each stack in fleet

He3 is consumed regardless of win/loss.
Ships must have enough He3 storage for the battle duration.
```

---

### 8.7 Normal Instances (PvE)

#### 8.7.1 Overview

30 Normal Instances providing progressive PvE challenges. All completable with Frigate/Cruiser/Battleship fleets (no Special Hulls or Flagships required).

| Property | Value |
|----------|-------|
| **Total Instances** | 30 |
| **Repeatable** | Yes (data does NOT reset) |
| **Ship Loss** | Permanent (no Spacedock repair) |
| **He3 Loss** | Yes |
| **Rewards** | Treasure Box (resources + 10% blueprint chance) |

#### 8.7.2 Instance List

| # | Name | Max Fleets | EXP | Min Level | Blueprint Drops |
|---|------|------------|-----|-----------|-----------------|
| 1 | Ancestral Recall | 3 | 180 | 1 | Weikes, Rapid Fire |
| 2 | Deadzone | 4 | 500 | 3 | Typhoon, Taskmaster |
| 3 | Bravery | 4 | 1,000 | 5 | Estrella, Cluster Laser Transmitter |
| 4 | Dark Frontier | 5 | 1,800 | 8 | Air Wanderer, Gatling Cannon |
| 5 | Crimson Nebula | 5 | 2,800 | 10 | Bombardier, Magneto Pulsar |
| 6 | Astral Rift | 6 | 4,000 | 13 | Valkyrie, Rocket Frame |
| 7 | Void Passage | 6 | 5,500 | 15 | Duke, Starlight Missile Pod |
| 8 | Stellar Tempest | 7 | 7,200 | 18 | GoGetter, Streamliner |
| 9 | Iron Bastion | 7 | 9,200 | 20 | The Shuttler, Golem |
| 10 | Phantom Corridor | 8 | 11,500 | 23 | Space Hunter, Particle Stun Shield |
| 11 | Warp Anomaly | 8 | 14,000 | 25 | Watchman, Heat Diffusion Shield |
| 12 | Shattered Core | 9 | 16,800 | 28 | Sparrow, Space-Time Magnetic Shield |
| 13 | Nova Storm | 9 | 20,000 | 30 | Spinner, Detonator Shield |
| 14 | Eclipse Point | 10 | 23,500 | 33 | Devourer, Anti-Aircraft Cannon |
| 15 | Binary Star | 10 | 27,200 | 35 | Wraith, Powered Pulse Cannon |
| 16 | Gravity Well | 11 | 31,200 | 38 | Polymesus, Ship Reinforcement Facility |
| 17 | Plasma Fields | 11 | 35,500 | 40 | Encratos, Quick Reaction Armor |
| 18 | Ion Storm | 12 | 40,000 | 43 | Cybra, Reflective Plating |
| 19 | Cosmic Tides | 12 | 44,800 | 45 | Nicholas, Daedalus Control System |
| 20 | Asteroid Belt | 12 | 49,800 | 48 | Hamdar, Energy Armor |
| 21 | Dark Matter | 13 | 52,000 | 50 | Nettle, EOS Phase Shift Shield |
| 22 | Quantum Flux | 13 | 54,500 | 52 | Helena, Shield Regenerator |
| 23 | Supernova Rim | 13 | 57,000 | 54 | Diaz, Energy Shield Booster |
| 24 | Pulsar Gate | 14 | 59,500 | 56 | RV766-The Explorer, Agility Booster |
| 25 | Neutron Field | 14 | 62,000 | 58 | Palenka, Infrared Scanner |
| 26 | Singularity | 14 | 64,500 | 60 | ECM Booster, Auto Target System |
| 27 | Event Horizon | 15 | 67,000 | 62 | Time Dilation Module, Extreme Counterattack |
| 28 | Omega Rift | 15 | 69,500 | 64 | Team Combat Engine, Anti-Matter Engine |
| 29 | Final Frontier | 15 | 71,500 | 66 | Lander Module, Nano Station Warehouse |
| 30 | Triumphant Glory | 15 | 73,500 | 68 | ALL blueprints (equal chance) |

**Note:** Instance names for 4-29 are original (GO2 wiki did not have complete names). Blueprint drops are designed to provide progression gating.

#### 8.7.3 Treasure Box Mechanics

```
On instance completion:
  1. Award EXP to player
  2. Award resource reward (Metal + He3 + Gold, scaling with instance level)
  3. Roll Treasure Box:
     - 90% chance: Small resource bonus (10-50% of base reward)
     - 10% chance: Blueprint drop
       - Blueprint selected equally from instance's blueprint pool
       - Example: Instance 1 has 2 blueprints -> 5% chance each
```

#### 8.7.4 Enemy Fleet Compositions

Enemy fleets scale with instance number. Each instance has pre-defined enemy fleet compositions:

```
Instance Enemy Scaling:
  Ships per stack = 200 + (instance_number * 150)
  Enemy hull tier = ceil(instance_number / 10)  // I for 1-10, II for 11-20, III for 21-30
  Number of enemy fleets = instance max_fleets
  Modules per enemy ship = 3 + floor(instance_number / 5)

Enemy fleet targeting follows the same rules as player targeting:
  Missile fleets -> target player's highest attack fleet
  Ship-based fleets -> target player's highest durability fleet
```

---

### 8.8 Phase 2 Data Models (SQL)

New tables and modifications for Phase 2. These are additive (Phase 1 tables remain unchanged).

#### 8.8.1 Hull Types (Reference Data)

```sql
CREATE TABLE hull_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    display_name TEXT NOT NULL,
    hull_class TEXT NOT NULL CHECK (hull_class IN ('frigate', 'cruiser', 'battleship')),
    tier INTEGER NOT NULL CHECK (tier BETWEEN 1 AND 3),
    armor_type TEXT NOT NULL CHECK (armor_type IN ('nano', 'chrome', 'regen', 'neutralizing')),
    base_shield INTEGER NOT NULL DEFAULT 0,
    base_structure INTEGER NOT NULL DEFAULT 0,
    base_stability NUMERIC(5,2) NOT NULL DEFAULT 100.00,
    base_defense NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    installation_slots INTEGER NOT NULL DEFAULT 100,
    base_agility INTEGER NOT NULL DEFAULT 0,
    base_movement INTEGER NOT NULL DEFAULT 0,
    base_storage INTEGER NOT NULL DEFAULT 0,
    base_metal_cost BIGINT NOT NULL DEFAULT 0,
    base_he3_cost BIGINT NOT NULL DEFAULT 0,
    base_gold_cost BIGINT NOT NULL DEFAULT 0,
    base_build_time_seconds INTEGER NOT NULL DEFAULT 10,
    description TEXT NOT NULL DEFAULT ''
);

CREATE INDEX idx_hull_types_class ON hull_types (hull_class);
CREATE INDEX idx_hull_types_tier ON hull_types (tier);
```

#### 8.8.2 Module Types (Reference Data)

```sql
CREATE TABLE module_types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    display_name TEXT NOT NULL,
    category TEXT NOT NULL CHECK (category IN (
        'ballistic', 'directional', 'missile', 'ship_based', 'planetary',
        'structure', 'shield', 'air_defense',
        'electronic', 'storage', 'transmission'
    )),
    tier INTEGER NOT NULL DEFAULT 1 CHECK (tier BETWEEN 0 AND 3),
    -- 0 = no tiers (e.g., Atomic Framework)
    damage_type TEXT CHECK (damage_type IN ('kinetic', 'heat', 'explosive', 'magnetic', 'siege', NULL)),
    min_damage INTEGER NOT NULL DEFAULT 0,
    max_damage INTEGER NOT NULL DEFAULT 0,
    weapon_range_min INTEGER NOT NULL DEFAULT 0,
    weapon_range_max INTEGER NOT NULL DEFAULT 0,
    cooldown INTEGER NOT NULL DEFAULT 0,
    he3_per_round INTEGER NOT NULL DEFAULT 0,
    volume INTEGER NOT NULL DEFAULT 1,
    max_per_ship INTEGER NOT NULL DEFAULT 0,
    -- 0 = unlimited
    effects_json JSONB NOT NULL DEFAULT '{}',
    -- e.g., {"shield_bonus": 300, "agility_bonus": 1, "movement_bonus": 1}
    metal_cost BIGINT NOT NULL DEFAULT 0,
    he3_cost BIGINT NOT NULL DEFAULT 0,
    gold_cost BIGINT NOT NULL DEFAULT 0,
    build_time_seconds INTEGER NOT NULL DEFAULT 1,
    description TEXT NOT NULL DEFAULT '',

    CONSTRAINT uq_module_name_tier UNIQUE (name, tier),
    CONSTRAINT chk_damage_range CHECK (min_damage <= max_damage),
    CONSTRAINT chk_weapon_range CHECK (weapon_range_min <= weapon_range_max)
);

CREATE INDEX idx_module_types_category ON module_types (category);
CREATE INDEX idx_module_types_tier ON module_types (tier);
```

#### 8.8.3 Ship Blueprints (Reference Data - replaces existing `blueprints` table)

The existing `blueprints` table from Phase 1 is replaced with a more complete structure:

```sql
-- Drop old blueprints table and recreate
DROP TABLE IF EXISTS player_blueprints;
DROP TABLE IF EXISTS blueprints;

CREATE TABLE blueprints (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    blueprint_type TEXT NOT NULL CHECK (blueprint_type IN ('hull', 'module')),
    hull_type_id INTEGER REFERENCES hull_types(id) ON DELETE SET NULL,
    module_type_id INTEGER REFERENCES module_types(id) ON DELETE SET NULL,
    source TEXT NOT NULL DEFAULT 'instance'
        CHECK (source IN ('instance', 'quest', 'auction', 'trafficker', 'mall')),
    research_level INTEGER NOT NULL DEFAULT 0 CHECK (research_level BETWEEN 0 AND 3),
    description TEXT NOT NULL DEFAULT '',

    CONSTRAINT chk_blueprint_ref CHECK (
        (blueprint_type = 'hull' AND hull_type_id IS NOT NULL AND module_type_id IS NULL) OR
        (blueprint_type = 'module' AND module_type_id IS NOT NULL AND hull_type_id IS NULL)
    )
);

CREATE INDEX idx_blueprints_type ON blueprints (blueprint_type);
CREATE INDEX idx_blueprints_hull ON blueprints (hull_type_id) WHERE hull_type_id IS NOT NULL;
CREATE INDEX idx_blueprints_module ON blueprints (module_type_id) WHERE module_type_id IS NOT NULL;
```

#### 8.8.4 Player Blueprints (Unlocked)

```sql
CREATE TABLE player_blueprints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    blueprint_id INTEGER NOT NULL REFERENCES blueprints(id) ON DELETE CASCADE,
    is_activated BOOLEAN NOT NULL DEFAULT false,
    research_level INTEGER NOT NULL DEFAULT 1 CHECK (research_level BETWEEN 1 AND 3),
    acquired_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_player_blueprint UNIQUE (player_id, blueprint_id)
);

CREATE INDEX idx_player_blueprints_player ON player_blueprints (player_id);
CREATE INDEX idx_player_blueprints_activated ON player_blueprints (is_activated) WHERE is_activated = true;
```

#### 8.8.5 Ship Designs (Updated)

The existing `ship_designs` table is updated to reference the new hull_types and store module placement:

```sql
-- Drop and recreate with proper references
DROP TABLE IF EXISTS ships;
DROP TABLE IF EXISTS ship_designs;

CREATE TABLE ship_designs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    hull_type_id INTEGER NOT NULL REFERENCES hull_types(id),
    modules_json JSONB NOT NULL DEFAULT '[]',
    -- Array of {module_type_id, quantity, placement_order}
    -- Stats computed from hull + modules
    total_shield INTEGER NOT NULL DEFAULT 0,
    total_structure INTEGER NOT NULL DEFAULT 0,
    total_defense NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    total_agility INTEGER NOT NULL DEFAULT 0,
    total_movement INTEGER NOT NULL DEFAULT 0,
    total_storage INTEGER NOT NULL DEFAULT 0,
    attack_power INTEGER NOT NULL DEFAULT 0,
    weapon_range_min INTEGER NOT NULL DEFAULT 0,
    weapon_range_max INTEGER NOT NULL DEFAULT 0,
    volume_used INTEGER NOT NULL DEFAULT 0,
    he3_per_round INTEGER NOT NULL DEFAULT 0,
    metal_cost BIGINT NOT NULL DEFAULT 0,
    he3_cost BIGINT NOT NULL DEFAULT 0,
    gold_cost BIGINT NOT NULL DEFAULT 0,
    build_time_seconds INTEGER NOT NULL DEFAULT 10,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_design_name CHECK (name ~ '^[a-zA-Z0-9._-]+$'),
    CONSTRAINT chk_volume CHECK (volume_used >= 0)
);

CREATE INDEX idx_ship_designs_player ON ship_designs (player_id);
CREATE INDEX idx_ship_designs_hull ON ship_designs (hull_type_id);
```

#### 8.8.6 Ships (Updated)

```sql
CREATE TABLE ships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    ship_design_id UUID NOT NULL REFERENCES ship_designs(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 0,
    is_building BOOLEAN NOT NULL DEFAULT false,
    build_quantity INTEGER NOT NULL DEFAULT 0,
    build_finish_at TIMESTAMPTZ,
    production_slot INTEGER NOT NULL DEFAULT 1 CHECK (production_slot BETWEEN 1 AND 5),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_player_ship_design UNIQUE (player_id, ship_design_id),
    CONSTRAINT chk_quantity_non_negative CHECK (quantity >= 0),
    CONSTRAINT chk_build_consistency CHECK (
        (is_building = true AND build_finish_at IS NOT NULL AND build_quantity > 0) OR
        (is_building = false AND build_finish_at IS NULL AND build_quantity = 0)
    )
);

CREATE INDEX idx_ships_player ON ships (player_id);
CREATE INDEX idx_ships_building ON ships (is_building) WHERE is_building = true;
```

#### 8.8.7 Fleets (Updated)

```sql
DROP TABLE IF EXISTS fleets;

CREATE TABLE fleets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    name TEXT NOT NULL DEFAULT 'Fleet',
    formation TEXT NOT NULL DEFAULT 'phalanx'
        CHECK (formation IN ('phalanx', 'diamond', 'battle_line', 'skirmish', 'tee_forward', 'enfilade', 'tee_reverse')),
    commander_id UUID REFERENCES commanders(id) ON DELETE SET NULL,
    targeting_command TEXT NOT NULL DEFAULT 'max_attack'
        CHECK (targeting_command IN ('max_attack', 'min_attack', 'max_durability', 'min_durability', 'closest', 'by_commander_rank')),
    status TEXT NOT NULL DEFAULT 'stationed'
        CHECK (status IN ('stationed', 'traveling', 'combat', 'returning', 'dismissed')),
    planet_id UUID REFERENCES planets(id) ON DELETE SET NULL,
    position_x INTEGER,
    position_y INTEGER,
    destination_x INTEGER,
    destination_y INTEGER,
    arrival_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_travel_consistency CHECK (
        (status IN ('traveling', 'returning') AND destination_x IS NOT NULL
            AND destination_y IS NOT NULL AND arrival_at IS NOT NULL) OR
        (status IN ('stationed', 'combat', 'dismissed'))
    )
);

CREATE INDEX idx_fleets_player ON fleets (player_id);
CREATE INDEX idx_fleets_status ON fleets (status);
CREATE INDEX idx_fleets_commander ON fleets (commander_id) WHERE commander_id IS NOT NULL;
```

#### 8.8.8 Fleet Stacks (New - replaces grid_json)

```sql
CREATE TABLE fleet_stacks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id UUID NOT NULL REFERENCES fleets(id) ON DELETE CASCADE,
    ship_design_id UUID NOT NULL REFERENCES ship_designs(id),
    grid_row INTEGER NOT NULL CHECK (grid_row BETWEEN 0 AND 2),
    grid_col INTEGER NOT NULL CHECK (grid_col BETWEEN 0 AND 2),
    ship_count INTEGER NOT NULL DEFAULT 0 CHECK (ship_count BETWEEN 0 AND 3000),

    CONSTRAINT uq_fleet_position UNIQUE (fleet_id, grid_row, grid_col),
    CONSTRAINT chk_ship_count CHECK (ship_count >= 0 AND ship_count <= 3000)
);

CREATE INDEX idx_fleet_stacks_fleet ON fleet_stacks (fleet_id);
CREATE INDEX idx_fleet_stacks_design ON fleet_stacks (ship_design_id);
```

#### 8.8.9 Instance Blueprints (Junction Table)

```sql
CREATE TABLE instance_blueprints (
    instance_id INTEGER NOT NULL REFERENCES instances(id) ON DELETE CASCADE,
    blueprint_id INTEGER NOT NULL REFERENCES blueprints(id) ON DELETE CASCADE,

    PRIMARY KEY (instance_id, blueprint_id)
);
```

#### 8.8.10 Spacedock Repairs (New)

```sql
CREATE TABLE spacedock_repairs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    ship_design_id UUID NOT NULL REFERENCES ship_designs(id),
    destroyed_count INTEGER NOT NULL DEFAULT 0,
    repaired_count INTEGER NOT NULL DEFAULT 0,
    repair_finish_at TIMESTAMPTZ,
    combat_report_id UUID REFERENCES combat_reports(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_counts CHECK (destroyed_count >= 0 AND repaired_count >= 0 AND repaired_count <= destroyed_count)
);

CREATE INDEX idx_spacedock_repairs_player ON spacedock_repairs (player_id);
CREATE INDEX idx_spacedock_repairs_repairing ON spacedock_repairs (repair_finish_at) WHERE repair_finish_at IS NOT NULL;
```

#### 8.8.11 Blueprint Research Progress (New)

```sql
CREATE TABLE blueprint_research (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    player_blueprint_id UUID NOT NULL REFERENCES player_blueprints(id) ON DELETE CASCADE,
    target_level INTEGER NOT NULL CHECK (target_level BETWEEN 2 AND 3),
    is_researching BOOLEAN NOT NULL DEFAULT false,
    research_finish_at TIMESTAMPTZ,
    metal_cost BIGINT NOT NULL DEFAULT 0,
    he3_cost BIGINT NOT NULL DEFAULT 0,
    gold_cost BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_research_consistency CHECK (
        (is_researching = true AND research_finish_at IS NOT NULL) OR
        (is_researching = false AND research_finish_at IS NULL)
    )
);

CREATE INDEX idx_blueprint_research_player ON blueprint_research (player_id);
CREATE INDEX idx_blueprint_research_active ON blueprint_research (is_researching) WHERE is_researching = true;
```

---

### 8.9 Phase 2 API Endpoints

All endpoints prefixed with `/api`. Auth required (JWT Bearer token).

#### 8.9.1 Ship Factory

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/ship-factory` | Get Ship Factory status (level, slots, active builds) |
| `GET` | `/api/ship-factory/slots` | List all 5 production slots and their status |
| `POST` | `/api/ship-factory/build` | Start building ships in a production slot |
| `POST` | `/api/ship-factory/cancel/:slot` | Cancel build in a specific slot |

**POST /api/ship-factory/build**
```json
// Request:
{
    "ship_design_id": "uuid",
    "quantity": 500,
    "production_slot": 1
}

// Response 200:
{
    "slot": 1,
    "ship_design_id": "uuid",
    "quantity": 500,
    "build_finish_at": "2026-02-06T06:00:00Z",
    "resources_spent": { "metal": 250000, "he3": 200000, "gold": 300000 }
}

// Validations:
// - Ship Factory exists and is not upgrading
// - Production slot is unlocked (level/tech check)
// - Slot is not already in use
// - quantity <= 2,000,000
// - Player has required resources
// - Ship design exists and belongs to player
```

#### 8.9.2 Ship Designs

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/ship-designs` | List all player ship designs (max 20) |
| `POST` | `/api/ship-designs` | Create a new ship design |
| `PUT` | `/api/ship-designs/:id` | Update a ship design (only if no ships built) |
| `DELETE` | `/api/ship-designs/:id` | Delete a ship design (only if no ships exist) |
| `GET` | `/api/ship-designs/:id/stats` | Calculate stats for a design |

**POST /api/ship-designs**
```json
// Request:
{
    "name": "Frigate-Ballistic-V1",
    "hull_type_id": 1,
    "modules": [
        { "module_type_id": 1, "quantity": 5, "placement_order": 1 },
        { "module_type_id": 20, "quantity": 2, "placement_order": 2 },
        { "module_type_id": 35, "quantity": 1, "placement_order": 3 }
    ]
}

// Response 201:
{
    "design": {
        "id": "uuid",
        "name": "Frigate-Ballistic-V1",
        "hull_type_id": 1,
        "total_shield": 570,
        "total_structure": 1770,
        "attack_power": 90,
        "weapon_range_min": 1,
        "weapon_range_max": 2,
        "volume_used": 72,
        "he3_per_round": 10,
        "build_time_seconds": 30
    }
}

// Validations:
// - Player has < 20 designs
// - Name matches regex ^[a-zA-Z0-9._-]+$
// - Player owns hull blueprint (activated)
// - Player owns all module blueprints (activated)
// - Total volume <= hull installation_slots
// - Per-ship module limits respected
// - Max 20 characters in name
```

#### 8.9.3 Blueprints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/blueprints` | List all blueprints in the game (hull + module) |
| `GET` | `/api/blueprints/mine` | List player's unlocked blueprints |
| `POST` | `/api/blueprints/:id/activate` | Activate an unactivated blueprint |
| `POST` | `/api/blueprints/:id/research` | Start researching a blueprint (Weapon Research Center) |

#### 8.9.4 Fleets

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/fleets` | List all player fleets |
| `POST` | `/api/fleets` | Create a new fleet |
| `PUT` | `/api/fleets/:id` | Update fleet (formation, targeting, commander) |
| `DELETE` | `/api/fleets/:id` | Disband a fleet (ships return to pool) |
| `POST` | `/api/fleets/:id/assign-stack` | Assign ships to a grid position |
| `POST` | `/api/fleets/:id/remove-stack` | Remove ships from a grid position |
| `POST` | `/api/fleets/:id/move` | Send fleet to coordinates |
| `POST` | `/api/fleets/:id/recall` | Recall a traveling fleet |
| `POST` | `/api/fleets/:id/dismiss` | Dismiss fleet (offline protection) |

**POST /api/fleets/:id/assign-stack**
```json
// Request:
{
    "ship_design_id": "uuid",
    "grid_row": 0,
    "grid_col": 1,
    "ship_count": 3000
}

// Response 200:
{
    "stack": {
        "fleet_id": "uuid",
        "ship_design_id": "uuid",
        "grid_row": 0,
        "grid_col": 1,
        "ship_count": 3000
    },
    "available_ships": 12000
}

// Validations:
// - Fleet belongs to player
// - Fleet is stationed (not traveling/combat)
// - Grid position is empty
// - Position is valid for current formation
// - Player has enough unassigned ships of this design
// - ship_count <= 3000
// - One design per stack
```

#### 8.9.5 Instances

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/instances` | List all normal instances |
| `GET` | `/api/instances/:id` | Get instance details (enemy fleets, rewards) |
| `POST` | `/api/instances/:id/attempt` | Attempt an instance with selected fleets |
| `GET` | `/api/instances/progress` | Get player's instance completion progress |

**POST /api/instances/:id/attempt**
```json
// Request:
{
    "fleet_ids": ["uuid1", "uuid2", "uuid3"]
}

// Response 200:
{
    "report_id": "uuid",
    "result": "attacker_win",
    "total_rounds": 35,
    "exp_gained": 1000,
    "treasure_box": {
        "resources": { "metal": 5000, "he3": 3000, "gold": 8000 },
        "blueprint_id": 15  // or null if no blueprint drop
    },
    "losses": {
        "ships_destroyed": { "design_uuid_1": 45, "design_uuid_2": 120 },
        "he3_consumed": 85000
    }
}

// Validations:
// - Instance exists and is type 'normal'
// - Player level >= instance min_level
// - fleet_ids.length <= instance max_fleets
// - All fleets belong to player and are stationed
// - All fleets have at least 1 stack with ships
```

#### 8.9.6 Spacedock

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/spacedock` | Get Spacedock status (level, active repairs) |
| `GET` | `/api/spacedock/repairs` | List all pending/completed repairs |
| `POST` | `/api/spacedock/repair` | Start repairing ships from a PvP loss |
| `POST` | `/api/spacedock/accelerate` | Spend Mall Points to accelerate repair |

#### 8.9.7 Hull & Module Reference Data

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/hull-types` | List all hull types with stats |
| `GET` | `/api/module-types` | List all module types with stats |
| `GET` | `/api/module-types?category=ballistic` | Filter modules by category |

---

### 8.10 Phase 2 Formulas

#### 8.10.1 Ship Build Cost (per ship)

```
MetalCost  = hull.base_metal_cost + SUM(module.metal_cost * module.quantity)
He3Cost    = hull.base_he3_cost + SUM(module.he3_cost * module.quantity)
GoldCost   = hull.base_gold_cost + SUM(module.gold_cost * module.quantity)

BatchMetalCost = MetalCost * batch_quantity
BatchHe3Cost   = He3Cost * batch_quantity
BatchGoldCost  = GoldCost * batch_quantity
```

#### 8.10.2 Ship Build Time (per batch)

```
BaseShipTime = hull.base_build_time_seconds + SUM(module.build_time_seconds * module.quantity)
                + (atomic_framework_count * 1) + (orbital_shield_count * 1)

SpeedReduction = ShipFactorySpeedBonus + ConstructionBoostTech
EffectiveTime = BaseShipTime * (1 - SpeedReduction / 100)

BatchTime = EffectiveTime * batch_quantity
```

#### 8.10.3 Spacedock Repair

```
FOR EACH destroyed ship of a given design:
  IF rand() < (spacedock_repair_pct / 100):
    ship is recovered
  ELSE:
    ship is permanently lost

RepairTime = destroyed_count * 10 seconds  // base time per ship
EffectiveRepairTime = RepairTime * (1 - ConstructionBoostTech / 100)
```

#### 8.10.4 Instance Rewards

```
ResourceReward:
  metal = instance_number * 500 + rand(0, instance_number * 200)
  he3   = instance_number * 400 + rand(0, instance_number * 150)
  gold  = instance_number * 600 + rand(0, instance_number * 250)

EXP = instance.exp_reward

BlueprintRoll:
  IF rand() < 0.10:  // 10% chance
    blueprint = random_choice(instance.blueprint_pool)
    Award unactivated blueprint to player
```

#### 8.10.5 Fleet Speed

```
FleetSpeed = MIN(design.total_movement) across all stacks in fleet
TravelTime = distance / FleetSpeed * BASE_TRAVEL_TIME_PER_UNIT
```

---

### 8.11 Phase 2 User Flows

#### 8.11.1 Ship Design Flow

```
1. Player opens Ship Factory -> "Design" tab
2. Player clicks "New Design"
3. Select hull class (Frigate / Cruiser / Battleship)
4. Select specific hull (e.g., Weikes-I) -- must own activated blueprint
5. Design screen shows:
   - Hull stats (shield, structure, slots, armor, movement)
   - Module palette (filtered to owned/activated module blueprints)
   - Module placement area (with placement order shown)
   - Running total: volume used / available, calculated stats
6. Player drags modules into design:
   - Weapons fill from bottom
   - Shields fill above weapons
   - Electronic/engine at top
   - Placement order enforced
7. Player names the design (no spaces, . - _ only)
8. Player clicks "Save"
9. Frontend calls POST /api/ship-designs
10. Server validates volume, blueprints, limits, name
11. Design saved with computed stats
```

#### 8.11.2 Ship Building Flow

```
1. Player opens Ship Factory -> "Build" tab
2. Shows 5 production slots (locked/unlocked/active)
3. Player selects an empty unlocked slot
4. Player selects a ship design from their list
5. Player enters quantity (1 to 2,000,000)
6. Frontend shows:
   - Cost per ship and total cost
   - Build time per ship and total time
   - Speed bonus from Ship Factory level
7. Player confirms
8. Frontend calls POST /api/ship-factory/build
9. Server validates:
   - Slot is unlocked and empty
   - Player has resources
   - Quantity within limit
10. Production starts, countdown shown
11. On completion: ships added to player's ship pool
```

#### 8.11.3 Fleet Creation Flow

```
1. Player opens Fleet panel (airplane icon in bottom-right)
2. Player clicks "New Fleet"
3. Fleet editor shows:
   - 3x3 grid with position names
   - Formation selector (Phalanx, Diamond, etc.)
   - Targeting command dropdown
   - Commander assignment slot
4. Player selects a grid position
5. Player chooses ship design from available ships
6. Player sets quantity (1-3000)
7. Ships assigned to position
8. Repeat for other positions
9. Player assigns commander (optional)
10. Player selects formation and targeting command
11. Player clicks "Save Fleet"
12. Frontend calls POST /api/fleets + POST /api/fleets/:id/assign-stack for each position
```

#### 8.11.4 Instance Battle Flow

```
1. Player opens Instance panel
2. Shows list of 30 Normal Instances with:
   - Name, level requirement, max fleets, status (completed/uncompleted)
3. Player selects an instance
4. Instance detail shows:
   - Enemy fleet preview (hull types, approximate strength)
   - Rewards (EXP, resource range, blueprint pool)
   - Max fleets allowed
5. Player selects fleets to send (up to max_fleets)
6. Player clicks "Attack"
7. Frontend calls POST /api/instances/:id/attempt
8. Server resolves combat (8-phase system):
   a. Player fleets vs enemy fleets
   b. Full combat simulation
   c. Generate combat report
9. Result screen shows:
   - Win/Loss
   - Round-by-round summary
   - Ships lost (permanently destroyed)
   - He3 consumed
   - Rewards earned (EXP, resources, blueprint if lucky)
10. Player collects rewards
11. Instance marked as completed in progress tracker
```

#### 8.11.5 Spacedock Repair Flow (PvP only)

```
1. After PvP loss, destroyed ships appear in Spacedock
2. Player opens Spacedock panel
3. Shows list of destroyed ship batches:
   - Design name, quantity destroyed, repair chance %
4. Player clicks "Repair All"
5. Frontend calls POST /api/spacedock/repair
6. Server rolls repair for each destroyed ship individually
7. Repair timer starts
8. On completion:
   - Successfully repaired ships return to player's ship pool
   - Failed repairs are permanently lost
9. Player can accelerate with Mall Points (10 MP = -10% time)
```

---

### 8.12 Phase 2 Seed Data

#### 8.12.1 Tutorial Quest Chain Extension

On first Ship Factory construction, Phase 2 tutorial begins:

| Step | Requirement | Reward |
|------|------------|--------|
| Build Ship Factory | Construct Ship Factory Lv 1 | Weikes Blueprint, Rapid Fire Blueprint |
| Design a Ship | Create first ship design | 5,000 Metal, 5,000 He3, 5,000 Gold |
| Build Ships | Build 100 ships | Typhoon Blueprint |
| Build Spacedock | Construct Spacedock Lv 1 | Estrella Blueprint |
| Create a Fleet | Create first fleet with ships | Energy Shield Booster Blueprint, Super Transmission Engine Blueprint |
| First Instance | Complete Instance 1 (Ancestral Recall) | 10,000 Metal, 10,000 He3, 15,000 Gold |

#### 8.12.2 Initial Module Seed Count

Phase 2 seeds a representative module subset:

| Category | Module Lines | Total with Tiers |
|----------|-------------|-----------------|
| Ballistic | 3 (Rapid Fire, Taskmaster, Gatling Cannon) | 9 |
| Directional | 2 (Cluster Laser Transmitter, Magneto Pulsar) | 6 |
| Missile | 2 (Rocket Frame, Starlight Missile Pod) | 6 |
| Ship-Based | 2 (Streamliner, Golem) | 6 |
| Planetary | 1 (Lander Module) | 3 |
| Structure | 5 (Atomic Framework + 4 lines) | 13 |
| Shield | 8 (Orbital Shield + 7 lines) | 21 |
| Air Defense | 3 lines | 9 |
| Electronic | 5 lines | 15 |
| Storage | 2 (Station Warehouse, Nano Station Warehouse) | 2 |
| Transmission | 3 (Super Transmission Engine, Team Combat Engine, Anti-Matter Engine) | 7 |
| **TOTAL** | **36 module lines** | **97 individual modules** |

The remaining ~300 modules are added incrementally via content updates, keeping the instance blueprint drop system as the primary progression gate.

---

### 8.13 Phase 2 Migration Checklist

New tables to create in Phase 2 migration:

```
-- New reference tables
hull_types
module_types
ship_factory_levels
spacedock_levels

-- Modified tables (drop and recreate)
blueprints (new structure with blueprint_type)
player_blueprints (add is_activated, research_level)
ship_designs (new structure with hull_type_id, computed stats)
ships (add production_slot)
fleets (remove grid_json)

-- New tables
fleet_stacks
instance_blueprints
spacedock_repairs
blueprint_research
```

**Seed data required:**
- 75+ hull_types rows (25 hull lines x 3 tiers)
- 97 module_types rows (initial subset)
- 24 ship_factory_levels rows
- 12 spacedock_levels rows
- 100+ blueprints rows (25 hull + 36 module line blueprints)
- 30 instance rows (Normal Instances) with instance_blueprints mappings
- 6 starter blueprints for tutorial quest rewards
