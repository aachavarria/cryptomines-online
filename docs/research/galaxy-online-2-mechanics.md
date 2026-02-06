# Galaxy Online 2 (GO2) - Mechanics Research Document

> **Purpose**: Exhaustive reference of Galaxy Online 2 game mechanics to inform the Cryptomines Online Game Design Document.
>
> **Sources**: [Galaxy Online II Wiki (Fandom)](https://galaxyonlineii.fandom.com/wiki/Galaxy_Online_II_Wiki)
>
> **Status**: Galaxy Online II was a Facebook-based sci-fi strategy game by IGG. The game is no longer in service but has extensive wiki documentation.

---

## Table of Contents

1. [Game Overview](#1-game-overview)
2. [Resource System](#2-resource-system)
3. [Buildings](#3-buildings)
4. [Technology / Research](#4-technology--research)
5. [Ship & Fleet System](#5-ship--fleet-system)
6. [Combat System](#6-combat-system)
7. [Commander System](#7-commander-system)
8. [Planets & Colonies](#8-planets--colonies)
9. [Corps (Alliances)](#9-corps-alliances)
10. [Instances (PvE)](#10-instances-pve)
11. [Progression & Quest System](#11-progression--quest-system)
12. [Key Formulas & Scaling](#12-key-formulas--scaling)
13. [Monetization](#13-monetization)
14. [Design Takeaways for Cryptomines Online](#14-design-takeaways-for-cryptomines-online)

---

## 1. Game Overview

Galaxy Online II is a space strategy game blending tactical combat, base building, and ship design. Inspired by the Chinese Doujin game "Super Starship Wars OL", it incorporates elements of chess, shogi, and go with sci-fi naval combat.

**Core Gameplay Loop:**
1. Build and upgrade base buildings (ground + space)
2. Research technologies to unlock capabilities
3. Design custom ships from modular components
4. Recruit and enhance commanders
5. Run PvE instances for blueprints and resources
6. Build fleets and engage in PvP combat
7. Join a Corp for resource planets and group content
8. Expand and optimize resource production

**Two Base Locations:**
- **Ground Base (Planet)**: Resource buildings, civic structures, military facilities
- **Space Base (Space Station)**: Defensive structures, orbital weapons

---

## 2. Resource System

### 2.1 Primary Resources

| Resource | Producer Building | Max Buildings | Relative Rate | Primary Use |
|----------|-------------------|---------------|---------------|-------------|
| **Gold** | Residential Area | 8 | 1.0x (highest) | Research, trading, upgrades |
| **He3** (Helium-3) | He3 Extractor | 8 | ~0.5-0.6x Gold | Ship fuel, building upgrades |
| **Metal** | Metal Collector | 8 | ~0.4x Gold | Ship manufacturing, building upgrades |

### 2.2 Currencies

| Currency | Acquisition | Use |
|----------|-------------|-----|
| **Mall Points (MP)** | Premium purchase (real money) | Premium shop, auctions, commander cards |
| **Vouchers** | Daily from friends repairing structures (max 8/day), events | Speedups, some mall items |
| **Badges** | Restricted Instances, tasks | Web Mall purchases |
| **Honor Points** | League Matches (15-255 pts), Pirate Challenges | Honor shop purchases |
| **Champion Points** | Championship matches (max 150/day) | Championship shop |
| **Corsairs' Gold** | Events, trade-ins | Bionic Chips, special items |

### 2.3 Resource Storage & Collection

- Resources accumulate in the **Resource Warehouse** and must be manually collected
- Collection recommended every **12 hours or less** to prevent overflow
- Warehouse has limited capacity; upgrading increases storage
- PvP winner receives **20% of loser's resources** (excluding warehouse contents)
- Resources can be obtained from: production, PvP raids, instances, quests, resource packs

### 2.4 He3 Extractor - Production Per Level (1-24)

| Level | He3 Output/hr | Metal Cost | He3 Cost | Gold Cost | Build Time | Civic Req |
|-------|---------------|-----------|----------|-----------|-----------|-----------|
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

### 2.5 Metal Collector - Production Per Level (1-24)

| Level | Metal Output/hr | Level | Metal Output/hr |
|-------|-----------------|-------|-----------------|
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

**Key Scaling Pattern**: Resources roughly double every ~7 levels. Costs scale exponentially (~1.75x per level).

---

## 3. Buildings

### 3.1 Ground Base Buildings

#### Resource Buildings

| Building | Function | Max Count | Max Level |
|----------|----------|-----------|-----------|
| **Metal Collector** | Produces Metal | 8 | 24 |
| **He3 Extractor** | Produces He3 | 8 | 24 |
| **Residential Area** | Produces Gold (highest output) | 8 | 24 |
| **Resource Warehouse** | Stores all resources | Multiple | 24 |

#### Core / Administrative Buildings

| Building | Function | Max Level |
|----------|----------|-----------|
| **Civic Center** | Main hub; determines max level of all other buildings | 12 |
| **Alliance Center** | Enables Corp membership and features | - |
| **Trading Center** | Player-to-player trading and auctions | - |
| **Galaxy Transporter** | Inter-system resource transport | - |
| **Compound Center** | Commander card merging and enhancement | - |
| **Technology Center** | Research facility (7 science trees) | 12 |
| **Radar** | Detection of incoming attacks | - |

#### Military Buildings

| Building | Function | Max Level |
|----------|----------|-----------|
| **Ship Factory** | Constructs ships; holds 20 designs, 5 production slots | 24 |
| **Spacedock** | Ship berthing and fleet management | - |
| **Command Center** | Commander recruitment (60 max commanders at Lv71+) | 12 |
| **Weapon Research Center** | Develops weapons and modules | - |
| **Recycling Plant** | Recovers resources from scrapped ships | - |

#### Landscaping / Decorative

Casino Resort, Beacon, Monument, Fountain, Library, Theater, Park, College, Hospital, Shopping Center, Statue, Santa Sculpture (morale/aesthetic bonuses).

### 3.2 Space Base Buildings

| Building | Function |
|----------|----------|
| **Space Station** | Main orbital structure; must stay within 1 level of Civic Center |
| **Meteor Star** | Orbital defense structure |
| **Particle Cannon** | Energy weapon defense |
| **Anti-Aircraft Gun** | Anti-air defense |
| **Thor's Cannon** | Advanced heavy defense cannon |
| **Celestial Base** | Advanced orbital facility |

### 3.3 Civic Center Upgrade Table (Hub Building)

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

### 3.4 Technology Center Upgrade Table

| Level | Civic Req | Build Time | Metal | He3 | Gold | Research Time Reduction |
|-------|-----------|-----------|-------|-----|------|------------------------|
| 1 | 1 | 0:01:40 | 450 | 420 | 650 | 3% |
| 5 | 4 | 2:21:24 | 29,068 | 27,130 | 41,987 | 15% |
| 12 | 11 | 6338:12:12 | 48,935,424 | 45,673,062 | 70,684,501 | 36% |

### 3.5 Building System Rules

- **Construction Slots**: Limited concurrent construction (default 2); use Construction Cards for additional slots
- **Mutual Dependencies**: Civic Center <-> Space Station must stay within 1 level of each other
- **Prerequisite Chain**: Civic Center level gates most building upgrades
- **Speed Modifiers**: Logistics Construction Science research reduces build times

---

## 4. Technology / Research

### 4.1 Seven Science Trees

| Science | Focus | Key Unlocks |
|---------|-------|-------------|
| **Logistics Construction** | Economy, construction speed, ship building, resource dev | Construction Boost, Quality Materials, High Yield Mining/Chemistry/Investing, Expanded Capacity, Sync Shipbuilding |
| **Planetary Defense** | Space station defense systems | Station weapon upgrades, defense modules |
| **Ballistics Science** | Ballistic (gun) weapons | Ballistic damage, crits, scatter, penetration |
| **Directional Science** | Beam/directional weapons | Beam damage, range, accuracy |
| **Missile Science** | Missile weapons | Missile damage, range, AoE, scattering |
| **Ship-Based Science** | Fighter/ship-based weapons | Fighter damage, interception |
| **Ship Defense Science** | Shields, structure, armor | Shield/structure upgrades, damage mitigation |

### 4.2 Ballistics Science - Full Tech Tree

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

### 4.3 Ship Defense Science - Full Tech Tree

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

### 4.4 Logistics Construction Science (Priority Techs)

Recommended priority for new players (to Level 6):
1. **Construction Boost** - Reduces building construction time
2. **Quality Materials** - Reduces building resource costs
3. **High Yield Mining** - Increases Metal production
4. **High Yield Chemistry** - Increases He3 production
5. **High Yield Investing** - Increases Gold production
6. **Expanded Capacity** - Increases Resource Warehouse capacity
7. **Sync Shipbuilding** - Unlocks 5th Ship Factory production slot

### 4.5 Research Time Scaling

- Technology Center level reduces research time (3% per level, max 36% at Lv 12)
- Research costs increase exponentially per tech level
- Only one research can be active per tree at a time

---

## 5. Ship & Fleet System

### 5.1 Hull Types (Rock-Paper-Scissors)

| Hull | Effective Stack | Bonus vs | Penalty vs | Armor Types |
|------|----------------|----------|------------|-------------|
| **Frigate (F)** | 1,100 | Battleship (+5%) | Cruiser (-5%) | Nano, Chrome, Regen, Neutralizing |
| **Cruiser (C)** | 1,000 | Frigate (+5%) | Battleship (-5%) | Nano, Chrome, Regen, Neutralizing |
| **Battleship (B)** | 900 | Cruiser (+5%) | Frigate (-5%) | Nano, Chrome, Regen, Neutralizing |
| **Flagship** | Special | Fleet-wide bonuses | - | Special |
| **Humaroid-Flagship** | Special | Advanced bonuses | - | Special |

### 5.2 Ship Modules

Ships are customizable with modular components:

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

### 5.3 Ship Stats

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

### 5.4 Fleet Composition

**Fleet Grid**: 3x3 grid = 9 stacks maximum
**Stack Size**: 3,000 ships per stack (max 27,000 ships per fleet)
**Rule**: One ship design per stack, do NOT mix weapon types within a fleet

**Grid Positions & Attack Power:**

```
+------+----------+------+
| Head | Shoulder | Shoulder |   First Rank:  100% attack
+------+----------+------+
| Flank | Glasshouse | Flank | Second Rank:  90% attack
+------+----------+------+
| Rear  |   Tail   | Rear  |   Third Rank:  75% attack
+------+----------+------+
```

- **Glasshouse** (center-middle): Most protected position
- **Shoulders**: Most vulnerable positions
- **Fleet speed** = speed of the slowest ship

### 5.5 Fleet Formations

| Formation | Description |
|-----------|-------------|
| **Phalanx** | All 9 stacks filled; maximum firepower |
| **Diamond** | 5 stacks (center + 4 adjacent); protects central stack |
| **Battle Line** | 6 stacks (first 2 ranks); frontal defense |
| **Skirmish** | Sparse; minimizes scatter damage |
| **Tee Forward** | T-shape; balanced offense |
| **Enfilade** | Side-focused fire |
| **Tee Reverse** | Reversed T-shape |

### 5.6 Combat Readiness

```
Combat Readiness = (Ships in Stack / Effective Stack) x 100%
```

Example: 900 Battleships in a stack of 900 = 100% readiness. Only ships up to the effective stack number can attack each round.

### 5.7 Targeting Commands

Fleets can be ordered to target:
- Max/Min attack power
- Max/Min durability
- Closest proximity
- By commander rank

### 5.8 Blueprints

- **Required** to build any ship
- Primary source: Instances (PvE) via Treasure Boxes
- Also from: Auction House, Corp bonuses, events
- Treasure Box: **10% chance** of blueprint, split equally among available blueprints
- Avoid upgrading blueprints unless required for next level

### 5.9 Ship Factory

- Max 20 ship designs stored
- 5 production slots (5th requires Sync Shipbuilding research)
- Max production: 2,000,000 ships at a time
- Speed increases 1-60% with levels 1-24

---

## 6. Combat System

### 6.1 Combat Resolution (8 Steps)

#### Step 1: Attacker Fires
```
Attacks = (Attack modules per ship) x (Ships in effective stack) x (Hit chance)

Hit Chance influenced by:
  + Weapon base accuracy
  + Attacker steering stat (~4% per point)
  + Attacker accuracy stat (~1% per 12 points)
  - Defender agility (~4% per point)
  - Defender commander dodge stat
```

#### Step 2: Interceptors Fire
- **PPC (Powered Pulse Cannon)**: 55% chance to shoot down each incoming attack
- Intercept modules can counter missiles and ship-based weapons
- Defender's ship defense stat and commander dodge affect interception

#### Step 3: Calculate Damage
- Weapon damage randomly generated within min-max range
- Double-hits and critical hits applied here
- Modified by armor type (defender) vs damage type (attacker)

#### Step 4: Damage Negation
- Non-EOS shields reduce damage (Heat Diffusion, Daedalus Control, Energy Armor)
- Reduces scatter and piercing effects

#### Step 5: Shield Penetration
- Ballistic and directional weapons can pierce shields
- Penetration chance modified by tech research

#### Step 6: Deal Damage to Shields
- EOS Phase Shift: 30% chance to absorb double damage (with Damage Mitigation tech)
- Only effective stack triggers EOS defensively
- Combined shield total from all 3,000 ships absorbs damage

#### Step 7: Assign Damage to Hull
```
Ships Destroyed = floor(Remaining Damage / (Individual Ship Structure x Stability%))
```

#### Step 8: Calculate Scatter Damage
- Based on weapon type, technology, and Sandora module
- Scatter damage **cannot be absorbed or mitigated by defenses**
- Missile Exaltation tech: "higher total structure than target deal 54% scatter = 432% unpreventable bonus damage"

### 6.2 Battle Duration

- Minimum: **20 rounds** + number of fleets/buildings present
- Maximum: **99 rounds**
- Only **ballistic weapons** fire every round; others require cooldown

### 6.3 Ship Type Advantage Matrix

```
Frigate  ---(+5%)--->  Battleship
Cruiser  ---(+5%)--->  Frigate
Battleship ---(+5%)--> Cruiser
```

### 6.4 Armor Types vs Damage Types

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

### 6.5 Commander Impact on Combat

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

### 6.6 PvP Mechanics

- Scout with single-ship fleet before major attacks
- Winner receives **20% of loser's resources** (not warehouse contents)
- Use "synchronize arrival" to coordinate fleet timing
- Ships must be dismissed or truced when offline for protection

### 6.7 Combat Losses by Mode

| Mode | Ships Lost | He3 Lost |
|------|-----------|----------|
| Normal/Restricted Instances | Yes | Yes |
| Trial/Constellation Instances | No | Yes |
| League/Arena/Championships | No | No |
| PvP (Attack Neighbors) | Yes | Yes |

---

## 7. Commander System

### 7.1 Commander Rarity Tiers

| Tier | Power Level | Acquisition |
|------|-------------|-------------|
| **Common** | Lowest | Free recruitment (with cooldown) |
| **Skill** | Low | Lucky draws, events |
| **Super** | Medium | Lucky Wheel, instances, merging |
| **Legendary** | High | Restricted Instances, Lucky Wheel, events |
| **Divine** | Highest | Premium events, special merging |

### 7.2 Commander Attributes

| Attribute | Effect |
|-----------|--------|
| **Accuracy** | Increases weapon hit chance |
| **Dodge** | Reduces opponent hit rate |
| **Speed** | Determines attack order; affects successive strike chance |
| **Electron** | Increases Critical Hit Rate and Critical Damage |

### 7.3 Star Rank & Effective Stack

Commander cards can be merged in the Compound Center to increase Star Rank. This is the **single biggest upgrade in the game** because Star Rank directly increases the Effective Stack (how many ships can attack per round).

### 7.4 Skill Types

| Skill Color | Type |
|-------------|------|
| Blue | Defense Skill |
| Red | Attack Skill |
| Green | Energy Skill |

### 7.5 Enhancement Systems

- **Gems**: Attach to commander to increase dodge, attack, defense, electron (max 3 diamonds)
- **Bionic Chips**: Advanced enhancement via Cybernetics Center
  - Max 5 chips per commander
  - No duplicate chip types
  - Types: Max Planetary, Max Ballistics, Negator, etc.
  - Purchased with Corsairs' Gold or Mall Points

### 7.6 Recruitment Mechanics

| Method | Cost | Notes |
|--------|------|-------|
| Free Recruitment | Free | 3-hour cooldown between draws |
| Quick Recruitment | 8 Mall Points | No cooldown |
| Commander Cards | 100 Mall Points | Specific commander choice |
| Auction House | Variable | Player-to-player trading |
| Events | Free/Variable | Special events and updates |

### 7.7 Command Center Building

| Level | Civic Req | Cooldown | Build Time | Metal | He3 | Gold |
|-------|-----------|----------|-----------|-------|-----|------|
| 1 | 1 | 3:00:00 | 0:00:40 | 600 | 450 | 500 |
| 2 | 2 | 2:50:00 | 0:02:01 | 1,692 | 1,269 | 1,410 |
| 3 | 3 | 2:40:00 | 0:06:06 | 4,788 | 3,591 | 3,990 |
| 4 | 4 | 2:30:00 | 0:18:33 | 13,599 | 10,199 | 11,332 |
| 5 | 5 | 2:20:00 | 0:56:34 | 38,757 | 29,068 | 32,297 |

Max commanders at Player Level 71+: **60**

---

## 8. Planets & Colonies

### 8.1 Resource Bonus Planets (RBPs)

RBPs are the primary territorial mechanic. They are heavily fortified locations evenly spaced throughout the galaxy.

**Galaxy Layout**: 7x7+ grid of zones, each with 1 RBP in the center, spaced 60 movement squares apart.

### 8.2 RBP Bonuses

| Level Range | Bonus per Level | Cumulative at Top |
|-------------|----------------|-------------------|
| 1-10 | 5% base + 0.5%/level | 10% |
| 11-20 | 1% per level | 20% |
| 21-30 | 1.5% per level | 35% |
| 31-40 | 2% per level | 55% |
| ... | Continues scaling | ... |
| 100 | - | **280%** |

Bonuses apply to: Resource production, Research speed, Shipbuilding speed

### 8.3 RBP Defense

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

### 8.4 Conquest Mechanics

- Only **Corps** (not individuals) can attack RBPs
- **72-hour protection** after capture
- **24-hour battle phase** when vulnerable
- Successful takeover resets timer to 72 hours
- Corps have **99 combat rounds** to complete takeover
- Fleets cannot be recalled mid-battle
- On conquest: defenses reset to Level 1, Space Station retains level

**Multi-Corp Attacks**: Planet goes to the corp with the most kills from a single fleet type. Scoring: 1 point per ship destroyed, per structure destroyed, per Space Station destroyed.

### 8.5 RBP Control Limits

Corps can control one planet per Corp Level (Lv 10 corp = max 10 RBPs).

### 8.6 Upgrading RBPs

- Upgraded using **Corp Wealth** (from member donations)
- Space Station upgradeable to Level 100
- Defenses (Meteor Stars, Cannons, etc.) upgradeable to Level 10
- Total wealth for Level 100 Station: **24,902,439**

---

## 9. Corps (Alliances)

### 9.1 Overview

Corps are groups of players that work together as a military/economic unit. They are essential for mid-to-late game content.

### 9.2 Corp Features

| Feature | Description |
|---------|-------------|
| **Resource Bonus Planets** | Territorial control for production bonuses |
| **Corp Mall** | Shared ship shop (better inventory at higher levels) |
| **Corp Warehouse** | Shared storage facility |
| **Corp Merging Center** | Shared commander merging |
| **Pirate Planets** | Corp-level NPC combat challenges |
| **Galactic Wars** | Corp vs Corp large-scale conflict |

### 9.3 Donation System

| Metric | Value |
|--------|-------|
| **Contribution Points** | 1 point per 10,000 resources donated |
| **Mall Points conversion** | 1:1 ratio (much more efficient) |
| **Daily donation range** | 20-200 contribution points |
| **Max daily donation** | 2,000,000 resources = 200 points |

### 9.4 Corp Bonuses

- Resource production bonus (based on Corp level + RBP count)
- Science research acceleration
- Shipbuilding speed boost
- Access to Corp Mall inventory
- Protection from attacks (deterrent)

### 9.5 Corp Levels

Higher Corp level = more RBPs controllable (1 per level), better Mall inventory, stronger bonuses.

---

## 10. Instances (PvE)

### 10.1 Instance Types

| Type | Description | Losses |
|------|-------------|--------|
| **Normal** | Standard progression (difficulty scaling) | Ships + He3 |
| **Restricted** | Limited access, commander rewards | Ships + He3 |
| **Scenario (Trial)** | Trial-based combat scenarios | He3 only |
| **Constellation** | Advanced endgame content | He3 only |
| **Humaroid-Battles** | Special Collision Chaos content | Varies |

### 10.2 Rewards

- **Treasure Boxes**: 10% blueprint chance, split equally among available blueprints
  - Example: 5 possible blueprints = 2% each
- **Resources**: Gold, Metal, He3
- **Commanders**: Restricted Instances (levels 8-10)
- **Badges**: From Restricted Instances
- **Blueprints**: Instance-specific drops

### 10.3 Progression

- Higher levels yield better rewards
- Players can repeat instances for farming
- Instance Viewer tools available at krtools.info and inst.war2go.ru

---

## 11. Progression & Quest System

### 11.1 New Player Progression (Recommended)

**Day 1:**
1. Build every available building including space defenses
2. Prioritize Gold resource buildings first
3. Reserve one construction slot for quest requirements
4. Borrow ships from allies instead of building immediately

**Week 1:**
1. Focus on ONE weapon type + ONE ship class (recommended: ballistic + frigates OR missiles + battleships)
2. Upgrade resource buildings to level 14
3. Priority order: Gold > Warehouse > Civic Center/Space Station > He3/Metal
4. Advance Logistics Construction Science to level 6
5. Obtain key modules: EOS Phase Shift shield, TCE, AME
6. Do NOT build ships until EOS + TCE + AME are acquired

**Month 1:**
1. Start League participation
2. Run Restricted Instances
3. Begin Space Raids
4. Practice PvP against inactive neighbors
5. Invest in Galaxy Transporter
6. Farm champion points through League level 5+
7. Consider secondary account for resource support

### 11.2 Development Quests (Tutorial Chain)

Sample progression:

| Quest | Requirement | Reward Highlights |
|-------|------------|-------------------|
| Collecting Resources | Harvest Warehouse | 450M, 950H, 500G, Loudspeaker |
| Lv1 Technology Center | Build Tech Center | 2,250M, 2,100H, 3,250G |
| Metal Production | Build Metal Collector Lv1 | Super Transmission Engine Blueprint |
| He3 Production | Build He3 Extractor Lv1 | Estrella Blueprint |
| Lv1 Ship Factory | Build Ship Factory | Typhoon Blueprint |
| Lv1 Command Center | Build Command Center | Energy Shield Booster Blueprint |
| Design a Ship | Complete ship design | Resources |
| Build a Fleet | Create 1 fleet | Resources |
| Lv2 Space Station | Upgrade SS to Lv2 | 10,465M, 9,660H, 13,685G |

### 11.3 Daily Activities

| Activity | Reward |
|----------|--------|
| Daily Quests | 5,000 Gold + accumulated point rewards |
| Friend Visits | Up to 8 vouchers per 24 hours (from repairs) |
| Resource Collection | Manual harvest from warehouses |
| Instance Farming | Blueprints, resources, commanders |
| League Matches | Honor Points (15-255 per match) |

---

## 12. Key Formulas & Scaling

### 12.1 Resource Cost Scaling

Based on Civic Center data analysis:
```
Cost(Level N) ~ Cost(Level 1) x 3.03^(N-1)

Approximate multiplier per level: ~3.03x resources, ~2.87x build time
```

### 12.2 Production Scaling

Based on He3 Extractor data:
```
Production(Level N) ~ Production(Level 1) x 1.134^(N-1)

Growth rate: ~13.4% increase per level
Level 1 to Level 24: ~19.6x output increase
```

### 12.3 Effective Stack Formula

```
Effective Stack = Base Stack + Commander Star Rank Bonus

Base Stack:
  Frigate = 1,100
  Cruiser = 1,000
  Battleship = 900
```

### 12.4 Combat Readiness

```
Combat Readiness = min(Ships in Stack, Effective Stack) / Effective Stack x 100%
```

### 12.5 Defense Power

```
Defense Power = Stack Size x (Total Shields + Total Structure)
Full readiness requires 3,000 ships per stack
```

### 12.6 Ships Destroyed per Attack

```
Ships Destroyed = floor(Net Damage / (Ship Structure x Stability%))
```

### 12.7 Scatter Damage

```
Scatter Damage = Weapon Damage x Scatter% (from tech)
  - Bypasses ALL defenses
  - Applied to adjacent stacks
  - Missile Exaltation: up to 432% bonus unpreventable damage
```

### 12.8 Hit Chance

```
Hit Chance = Base Accuracy
  + (Attacker Steering x 4%)
  + (Attacker Accuracy / 12)
  - (Defender Agility x 4%)
  - Defender Commander Dodge
```

### 12.9 RBP Bonus Scaling

```
Levels 1-10:  5% + (Level x 0.5%)
Levels 11-20: Previous + (Level x 1%)
Levels 21-30: Previous + (Level x 1.5%)
Levels 31-40: Previous + (Level x 2%)
...continues...
Level 100: 280% total bonus
```

### 12.10 Corp Donation

```
Contribution Points = floor(Resources Donated / 10,000)
Mall Points: 1 MP = 1 Contribution Point
Max Daily: 200 points (2,000,000 resources)
```

---

## 13. Monetization

### 13.1 Free-to-Play Model

"Money helps with two things: **rarity** and **time**."

| Purchase Type | Effect |
|---------------|--------|
| **Mall Points** | Premium currency for exclusive items, speedups, commanders |
| **Speedups** | Accelerate building, research, ship production |
| **Exclusive Blueprints** | Ships not available through free progression |
| **Exclusive Commanders** | Higher rarity commanders |
| **Convenience** | Extra construction slots, instant recruitment |

### 13.2 Free Player Progression

Free players can access most content through:
- Vouchers (daily free currency)
- Instance farming
- League/Championship rewards
- Corp bonuses
- Time investment

---

## 14. Design Takeaways for Cryptomines Online

### 14.1 Core Systems to Adapt

1. **Triple Resource Economy**: Metal/He3/Gold maps well to crypto-themed resources (e.g., Hash Power / Energy / Tokens)
2. **Exponential Cost Scaling**: ~3x per building level creates meaningful progression gates
3. **Production Growth**: ~13% per level keeps upgrades rewarding without being explosive
4. **Hub Building Dependency**: Civic Center gates all other buildings - simple, effective progression control
5. **Modular Ship Design**: Deep customization is a major engagement driver
6. **Rock-Paper-Scissors Combat**: Simple hull type advantages (Frigate>BS>Cruiser>Frigate) with +/-5% modifiers
7. **Commander System**: Gacha-like card collection with meaningful combat impact through Star Rank
8. **Corp/Alliance Territorial Control**: RBPs provide meaningful group objectives

### 14.2 Key Balance Lessons

1. **Effective Stack** is the single most impactful stat - controls how many units attack per round
2. **Scatter damage** bypasses all defenses - important balance lever
3. **Shield vs Structure** creates dual-layer defense (absorb then destroy)
4. **Weapon cooldowns** create meaningful weapon type tradeoffs
5. **Formation positioning** adds tactical depth (100%/90%/75% attack rows)
6. **Blueprint gating** controls ship availability independently of resources

### 14.3 Monetization Lessons

1. Time acceleration is primary monetization (construction, research, production)
2. Rarity gating (commanders, blueprints) as secondary monetization
3. Free players can compete but with much longer timelines
4. Premium does NOT grant exclusive combat power, just faster access

### 14.4 Engagement Patterns

1. **Short-term**: Daily quests, resource collection, instance farming
2. **Medium-term**: Technology research, ship design optimization, fleet building
3. **Long-term**: Corp wars, RBP control, commander collection, championship ranking

---

## Sources

- [Galaxy Online II Wiki - Main](https://galaxyonlineii.fandom.com/wiki/Galaxy_Online_II_Wiki)
- [Beginner's Guide](https://galaxyonlineii.fandom.com/wiki/Beginner's_Guide)
- [Beginner FAQ](https://galaxyonlineii.fandom.com/wiki/Beginner_FAQ)
- [Guide To Advancing Quickly](https://galaxyonlineii.fandom.com/wiki/Guide_To_Advancing_Quickly)
- [Combat Mechanics](https://galaxyonlineii.fandom.com/wiki/Combat_Mechanics)
- [Fleet Design](https://galaxyonlineii.fandom.com/wiki/Fleet_Design)
- [Composite Ship Table](https://galaxyonlineii.fandom.com/wiki/Composite_Ship_Table)
- [Technology Center](https://galaxyonlineii.fandom.com/wiki/Technology_Center)
- [Ballistics Science](https://galaxyonlineii.fandom.com/wiki/Ballistics_Science)
- [Ship Defense Science](https://galaxyonlineii.fandom.com/wiki/Ship_Defense_Science)
- [Currency](https://galaxyonlineii.fandom.com/wiki/Currency)
- [Civic Center](https://galaxyonlineii.fandom.com/wiki/Civic_Center)
- [He3 Extractor](https://galaxyonlineii.fandom.com/wiki/He3_Extractor)
- [Metal Collector](https://galaxyonlineii.fandom.com/wiki/Metal_Collector)
- [Corps](https://galaxyonlineii.fandom.com/wiki/Corps)
- [Resource Bonus Planets](https://galaxyonlineii.fandom.com/wiki/Resource_Bonus_Planets)
- [Upgrading Resource Planets](https://galaxyonlineii.fandom.com/wiki/Upgrading_Resource_Planets)
- [Commander Cards](https://galaxyonlineii.fandom.com/wiki/Commander_Cards)
- [Command Center](https://galaxyonlineii.fandom.com/wiki/Command_Center)
- [Development Quests](https://galaxyonlineii.fandom.com/wiki/Development_Quests)
- [Instances](https://galaxyonlineii.fandom.com/wiki/Instances)
- [Galaxy Online II for Beginners](https://galaxyonlineii.fandom.com/wiki/Galaxy_Online_II_for_beginners!)
- [Walkthrough](https://galaxyonlineii.fandom.com/wiki/Walkthrough)
