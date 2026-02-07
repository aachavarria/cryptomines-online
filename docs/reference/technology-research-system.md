# Galaxy Online 2 - Technology & Research System Reference

> Research compiled from GO2 wiki (galaxyonlineii.fandom.com) and player guides.
> Date: 2026-02-06

---

## Table of Contents

1. [Overview](#1-overview)
2. [Technology Center Building](#2-technology-center-building)
3. [Seven Science Trees](#3-seven-science-trees)
   - 3.1 Logistics Construction Science
   - 3.2 Ballistics Science
   - 3.3 Directional Science
   - 3.4 Missile Science
   - 3.5 Ship-Based Science
   - 3.6 Ship Defense Science
   - 3.7 Planetary Defense Science
4. [Weapon Research Center](#4-weapon-research-center)
5. [Research Acceleration](#5-research-acceleration)
6. [Research Cost Analysis](#6-research-cost-analysis)
7. [Recommended Research Order](#7-recommended-research-order)
8. [Implementation Recommendations](#8-implementation-recommendations)

---

## 1. Overview

GO2 has two distinct research systems:

1. **Technology Center** - Seven science trees covering economy, combat, and defense
2. **Weapon Research Center** - Blueprint research for ship hulls and modules

### Core Rules
- Only one research can be active per tree at a time (7 trees = up to 7 concurrent researches)
- Technologies have prerequisite chains within their tree
- Technology Center level reduces research time by 3% per level (max 36% at Lv 12)
- All tech research costs Gold only (no Metal or He3)
- Blueprint research in Weapon Research Center is separate and costs different resources

---

## 2. Technology Center Building

**Purpose:** Research facility for all 7 science trees. Each level reduces research time by 3%.

**Grid Size:** 3x2 (same as Ship Factory, Resource Warehouse, Spacedock)

| Level | Civic Req | Build Time | Metal | He3 | Gold | Research Time Reduction |
|-------|-----------|-----------|-------|------|------|------------------------|
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

**Note:** Build times/costs are reduced by Construction Boost and Quality Materials techs.

---

## 3. Seven Science Trees

### 3.1 Logistics Construction Science

The most important research for new players. Covers economy, construction, shipbuilding, and resource production.

**Total techs: 11 | All costs in Gold only**

#### Tech Tree Structure

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

#### Summary Table

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

---

### 3.2 Ballistics Science

Enhances ballistic (kinetic/gun) weapons. Weapon range 1-2, lowest He3 cost.

**Total techs: ~16 | All costs in Gold only**

#### Tech Tree Structure

```
Ballistics (Base) [Lv 1-10] (no prereq)
  +5% ballistic damage per level (max 50%)
  Lv1: 541G, 0:01:57 | Lv10: 18,808G, 39:18:21

  +-- Ballistic Malice [Lv 1-5] (req: Ballistics Lv 3)
  |     +1% critical hit rate per level (max 5%)
  |     Lv1: 7,558G, 0:23:48 | Lv5: 31,327G, 6:48:55
  |
  +-- Ballistic Crackdown [Lv 1-2] (req: Ballistics Lv 3)
  |     +10% critical damage per level (max 20%)
  |     Lv1: 29,628G, 1:12:15 | Lv2: 34,303G, 1:36:51
  |
  +-- Steady Control Tech [Lv 1-5] (req: Ballistics Lv 6, Malice Lv 3)
  |     -2% weapon space per level (max -10%)
  |     Lv1: 26,108G, 0:56:06 | Lv5: 108,217G, 16:03:52
  |
  |     +-- Precise Ballistics [Lv 1-5] (req: Ballistics Lv 8, Steady Lv 3)
  |           +1% hit rate per level (max 5%)
  |           Lv1: 51,100G, 1:08:00 | Lv5: 211,811G, 19:28:20
  |
  +-- Shield Penetration [Lv 1] (req: Malice Lv 5, Crackdown Lv 1, Precise Lv 1)
  |     15% shield bypass damage
  |     154,616G, 4:15:00
  |
  +-- Depleted Uranium Bomb [Lv 1-3] (req: Crackdown Lv 2, Penetration Lv 1)
  |     +10-30% vs Neutral armor, +1-3% vs Light armor
  |     Lv1: 203,249G, 3:32:30 | Lv3: 359,977G, 11:06:35
  |
  +-- Fire Bomb Research [Lv 1-3] (req: Crackdown Lv 2, Penetration Lv 1)
  |     +10-30% vs Regen armor, +1-3% vs Light armor
  |     Lv1: 203,249G, 3:32:30 | Lv3: 359,977G, 11:06:35
  |
  +-- Improved Penetration [Lv 1-3] (req: DU Bomb Lv 1, Fire Bomb Lv 1)
  |     +1-3% vs Light armor, 3-10% shield pen chance
  |     Lv1: 484,776G, 6:22:30 | Lv3: 858,602G, 19:59:52
  |
  +-- Victory Rush [Lv 1] (req: DU Lv 3, Fire Lv 3, Imp.Pen. Lv 3)
  |     Range damage: 220%/180%/150%/120% at ranges 1-4
  |     +5% crit rate, +5% crit damage
  |     1,918,521G, 62:20:00
  |
  +-- Ballistic Scattering [Lv 1-5] (req: Ballistics Lv 10, Malice Lv 5, Steady Lv 5, Precise Lv 3)
  |     5-25% scatter damage to ships in same vertical row
  |     Lv1: 119,765G, 2:33:00 | Lv5: 496,433G, 43:48:45
  |
  +-- Improved Ballistic Scattering [Lv 1-3] (req: Precise Lv 5, Scattering Lv 3)
  |     +8-25% scattering rate
  |     Lv1: 240,487G, 4:57:30 | Lv3: 425,931G, 15:33:13
  |
  +-- Hop Bomb Research [Lv 1-5] (req: Scattering Lv 5, Imp. Scatter Lv 2)
  |     3-15% chance to deal 100% weapon damage as scatter
  |     Lv1: 349,930G, 5:40:00 | Lv5: 1,450,474G, 97:21:40
  |
  +-- Demolition Warhead Research [Lv ?] (req: Imp. Scatter Lv 3, Hop Bomb Lv 3)
  |     [Details not fully confirmed from wiki]
  |
  +-- Range Extension [Lv ?]
  |     [Referenced in Victory Rush description]
  |
  +-- Artillery Specialization [Lv ?]
       [Listed in wiki contents, details truncated]
```

---

### 3.3 Directional Science

Enhances beam/directional weapons. Range 2-5 (extendable to 2-6). Lowest He3 usage. Piercing damage hits all ships vertically.

**Total techs: 15 | All costs in Gold only**

#### Tech Tree Structure

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

---

### 3.4 Missile Science

Enhances missile/guided weapons. Scatters damage across ALL enemy ships. Range 3-6.

**Total techs: 16 | All costs in Gold only**

#### Tech Tree Structure

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
  +-- Rapid Loading [Lv ?] (details truncated from wiki)
  |     Reduces reload time by one round
  |
  +-- Perfect Storm [Lv ?] (details truncated from wiki)
        [Capstone tech, details unavailable]
```

---

### 3.5 Ship-Based Science

Enhances fighter-based weapons. Fighters have unique interception mechanics.

**Total techs: 16 | All costs in Gold only**

#### Tech Tree Structure

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

---

### 3.6 Ship Defense Science

Two parallel branches: Shields and Structure. Must be researched before offensive trees.

**Total techs: ~20 | All costs in Gold only**

#### Shield Branch

```
Ship Defense Tech (Base) [Lv 1-2] (no prereq)
  Lv1: +2% base shield/structure/agility/defense, +5% stability | 2,846G, 0:54:00
  Lv2: +5% base stats, +10% stability | 3,295G, 1:12:23

  +-- Shield Research [Lv 1-5] (req: Base Lv 1)
  |     +1-5% base shield
  |     Lv1: 10,392G, 0:45:00 | Lv5: 43,076G, 12:53:09
  |
  |     +-- Energy Diffusion Tech [Lv 1-3] (req: Base Lv 2, Shield Lv 3)
  |     |     Each shield module reduces damage by 1-3 points
  |     |     Lv1: 46,765G, 3:00:00 | Lv3: 82,828G, 9:24:39
  |     |
  |     +-- Penetration Resistance [Lv 1-2] (req: Shield Lv 5, Diffusion Lv 1)
  |     |     -3-7% enemy shield penetration chance
  |     |     Lv1: 117,576G, 6:00:00 | Lv2: 136,122G, 8:02:32
  |     |
  |     +-- Augment Shield [Lv 1-3] (req: Diffusion Lv 2, Pen.Res. Lv 1)
  |     |     +6-20% base shield
  |     |     Lv1: 248,747G, 11:00:00 | Lv3: 440,562G, 34:30:21
  |     |
  |     +-- Restoration [Lv 1-2] (req: Diffusion Lv 3, Aug.Shield Lv 2)
  |     |     +30-60% shield restore per round, +1-2% interception
  |     |     Lv1: 540,000G, 25:00:00 | Lv2: 625,178G, 33:30:32
  |     |
  |     +-- Augment Absorption [Lv 1-2] (req: Shield Lv 5, Diffusion Lv 1)
  |     |     Shield modules reduce damage by 2-5 points
  |     |     Lv1: 117,576G, 6:00:00 | Lv2: 136,122G, 8:02:32
  |     |
  |     +-- Energy Conservation [Lv 1-3] (req: Diffusion Lv 2, Aug.Abs. Lv 1)
  |     |     +3-10% chance absorb damage without He3
  |     |     Lv1: 248,747G, 11:00:00 | Lv3: 440,562G, 34:30:21
  |     |
  |     +-- Electronic Barrier [Lv 1-2] (req: Diffusion Lv 3, E.Cons. Lv 2)
  |           Reflect 5-10% damage before shields drop
  |           Lv1: 540,000G, 25:00:00 | Lv2: 625,178G, 33:30:32

  +-- Damage Mitigation [Lv 1-3] (req: Aug.Shield Lv 3, Restoration Lv 2, E.Cons. Lv 3, E.Barrier Lv 2)
        10-30% absorb double damage, 15-45% lower collateral
        Lv1: 805,152G, 30:00:00 | Lv3: 1,426,024G, 94:06:23
```

#### Structure Branch

```
  +-- Ship Structural Analysis [Lv 1-5] (req: Base Lv 1)
  |     +1-5% base structure
  |     Lv1: 10,392G, 0:45:00 | Lv5: 43,076G, 12:53:09
  |
  |     +-- Ship Reinforcement Tech [Lv 1-3] (req: Base Lv 2, Analysis Lv 3)
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

  +-- Stability Mastery [Lv 1-3] (req: Struct.Imp Lv 3, Fast Repair Lv 2, Def.Imp Lv 3, Reflection Lv 2)
        10-30% absorb double damage, 15-45% lower collateral
```

---

### 3.7 Planetary Defense Science

Enhances space station defense structures. Important for defending against attacks.

**Total techs: 8 | All costs in Gold only**

#### Tech Tree Structure

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

#### Planetary Defense Cost Table (Complete)

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

---

## 4. Weapon Research Center

**Purpose:** Research/upgrade ship hull blueprints and module blueprints. Separate from Technology Center.

**Grid Size:** 2x2

### Building Upgrade Table

| Level | Civic Req | Build Time | Metal | He3 | Gold | Research Time Reduction |
|-------|-----------|-----------|-------|------|------|------------------------|
| 1 | 1 | 0:00:40 | 500 | 300 | 450 | 3% |
| 6 | 6 | 2:53:05 | 92,371 | 55,422 | 83,134 | 18% |
| 12 | 12 | 2535:16:53 | 54,372,693 | 32,623,616 | 48,935,424 | 36% |

**Time Reduction:** 3% per level, identical formula to Technology Center.

### Two Research Tracks

1. **Ship Blueprint Research** - Upgrade existing hull blueprints for improved stats
2. **Module Blueprint Research** - Upgrade modules for improved stats (+0/10/25% at levels 1/2/3)

### Blueprint Research Levels (from GDD Section 8.3.4)

| Research Level | Requirement | Effect | Cost Multiplier |
|----------------|-------------|--------|-----------------|
| 1 (Base) | Blueprint acquired | Base stats | 1.0x |
| 2 | WRC Level 6 | +10% stats | 2.0x |
| 3 | WRC Level 10 | +25% stats | 5.0x |

---

## 5. Research Acceleration

### Technology Research Speedup
- Cost: **3 vouchers/MP per 30 minutes** of time reduction
- Alternative: **1 voucher per 10 minutes** remaining for instant completion

### Blueprint Research Speedup
- Cost: **8 vouchers/MP per 30 minutes** of time reduction (more expensive than tech)
- Wiki recommends spending vouchers on this rather than buildings/tech due to long times

### Building Construction Speedup
- Cost: **3 vouchers/MP per 30 minutes** of time reduction
- Friends can provide 2% free acceleration per building (one friend visit per building)

### Strategic Note
Blueprint research takes the longest and is the most expensive to accelerate, making it the highest-priority target for speedup items.

---

## 6. Research Cost Analysis

### Observed Patterns

**All tech research costs Gold only** - no Metal or He3.

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

**Scaling ratio analysis** (Level N+1 / Level N):
- Lv1->2: ~1.16x cost, ~1.35x time
- Lv2->3: ~1.53x cost, ~2.33x time
- Lv3->4: ~1.53x cost, ~2.34x time
- Lv4->5: ~1.53x cost, ~2.34x time
- Lv5->6: ~1.53x cost, ~2.34x time
- Lv6->7: ~1.53x cost, ~2.34x time
- Average multiplier: ~1.53x cost per level, ~2.34x time per level

**Approximate formula (base techs):**
```
GoldCost(level) ~= 541 * 1.53^(level-1)
BaseTime(level) ~= 117 * 2.34^(level-1)  // seconds
```

**Higher-tier techs** have different base costs but similar exponential scaling within their own level ranges.

### Total Gold Investment by Tree (estimated)

| Tree | Approx Total Gold (all techs maxed) |
|------|-------------------------------------|
| Logistics Construction | ~12,000,000 |
| Ballistics | ~6,500,000 |
| Directional | ~8,500,000 |
| Missile | ~8,000,000 |
| Ship-Based | ~14,000,000 |
| Ship Defense | ~8,000,000 |
| Planetary Defense | ~12,000,000 |

---

## 7. Recommended Research Order (GO2 Meta)

Based on player guides and wiki recommendations:

### Phase 1: Economy Foundation (Logistics)
1. Concurrent Construction Lv 1
2. Construction Boost to Lv 6
3. Quality Materials to Lv 6
4. High Yield Mining to Lv 6
5. High Yield Chemistry to Lv 6
6. High Yield Investing to Lv 6
7. Expand Capacity to Lv 3-4

### Phase 2: Ship Production
8. Ship Building Boost to Lv 6
9. Ship Building Logistics to Lv 4
10. Sync Shipbuilding Lv 1

### Phase 3: Choose ONE Weapon Tree
- **Ballistics** + Frigates = best for instance farming
- **Missiles** + Battleships = good for PvP scatter damage
- **Directional** = lowest He3, good for long fights
- **Ship-Based** = unique mechanics but harder to optimize

### Phase 4: Ship Defense
11. Ship Defense Tech Lv 2
12. Shield Research to Lv 5
13. Energy Diffusion to Lv 3
14. Continue down shield or structure branch based on fleet composition

### Phase 5: Advanced Weapon + Planetary Defense
- Push chosen weapon tree to capstone techs
- Start Planetary Defense for station security

---

## 8. Implementation Recommendations

### What Already Exists in GDD
- Technology Center building with full upgrade table (Section 2.2.7)
- Ballistics Science full tech tree (Section 2.3.2)
- Ship Defense Science full tech tree (Section 2.3.3)
- Logistics Construction priority techs (Section 2.3.4, partial)
- Research rules (Section 2.3.6)
- SQL schema for tech_types and technologies (Sections 3.6-3.7)
- API endpoints for research (Section 5.6)
- Research flow documentation (Section 7.5)
- Weapon Research Center building and blueprint research (Section 8.3.4)
- Blueprint research SQL schema (Section 8.8.11)

### What Needs to Be Added to GDD
1. **Logistics Construction** - Full tech tree with prerequisites, costs, max levels (currently partial)
2. **Directional Science** - Complete tech tree (currently marked NEEDS RESEARCH)
3. **Missile Science** - Complete tech tree (currently marked NEEDS RESEARCH)
4. **Ship-Based Science** - Complete tech tree (currently marked NEEDS RESEARCH)
5. **Planetary Defense** - Complete tech tree (currently marked NEEDS RESEARCH)
6. **Research cost formula** - The ~1.53x cost multiplier pattern (Section 4.6 NEEDS RESEARCH)
7. **tech_types seed data** - SQL INSERT statements for all ~90+ technologies across 7 trees

### Phase 3 Implementation Priority

**High priority (implement first):**
- Logistics Construction full tree (already partially in game, critical for economy)
- Ship Defense Science (already in GDD, needed before combat)
- Research cost formula in Section 4.6

**Medium priority:**
- Ballistics Science (already in GDD)
- One additional weapon tree (Directional or Missile)
- Weapon Research Center building + blueprint research

**Lower priority (can defer to Phase 4/5):**
- Remaining weapon trees
- Planetary Defense (only matters after space station defense is implemented)
- Capstone techs for all trees (Victory Rush, Ingenuity, Dynamic Impairment, etc.)

### SQL Seed Data Notes

The `tech_types` table needs entries for all technologies. Key fields:
- `tree` enum: logistics, ballistics, directional, missile, ship_based, ship_defense, planetary_defense
- `key` unique identifier: e.g., 'concurrent_construction', 'construction_boost', 'ballistics_base'
- `max_level`: varies per tech (1, 2, 3, 4, 5, 10)
- `prerequisites`: JSON array of {tech_key, required_level} pairs
- `base_gold_cost`: Lv 1 cost
- `gold_cost_multiplier`: ~1.16-1.53x depending on tier
- `base_research_time_seconds`: Lv 1 time
- `time_multiplier`: ~2.34x per level for base techs
- `effect_type`: enum or JSON describing what the tech modifies
- `effect_per_level`: numeric bonus per level

### Data Gaps Still Remaining
- Logistics Construction: Expand Capacity max level (7+, exact cap unknown)
- Ballistics: Demolition Warhead, Range Extension, Artillery Specialization full details
- Missile: Rapid Loading and Perfect Storm full details
- Ship Defense: Structure branch detailed costs (not on wiki)
- Planetary Defense: Defense Enhancement, Emplacement Mastery, Utmost Defense per-level costs (only Lv1 and Lv10 ranges shown)
- Technology Center Lv 2-4, 6-11 build costs (listed in GDD NEEDS RESEARCH)
- Weapon Research Center Lv 2-5, 7-11 build costs

---

## Sources

- Galaxy Online II Wiki - Technology Center: https://galaxyonlineii.fandom.com/wiki/Technology_Center
- Galaxy Online II Wiki - Ballistics Science: https://galaxyonlineii.fandom.com/wiki/Ballistics_Science
- Galaxy Online II Wiki - Directional Science: https://galaxyonlineii.fandom.com/wiki/Directional_Science
- Galaxy Online II Wiki - Missile Science: https://galaxyonlineii.fandom.com/wiki/Missile_Science
- Galaxy Online II Wiki - Ship-based Science: https://galaxyonlineii.fandom.com/wiki/Ship-based_Science
- Galaxy Online II Wiki - Ship Defense Science: https://galaxyonlineii.fandom.com/wiki/Ship_Defense_Science
- Galaxy Online II Wiki - Planetary Defense Science: https://galaxyonlineii.fandom.com/wiki/Planetary_Defense_Science
- Galaxy Online II Wiki - Weapon Research Center: https://galaxyonlineii.fandom.com/wiki/Weapon_Research_Center
- Galaxy Online II Wiki - Accelerate: https://galaxyonlineii.fandom.com/wiki/Accelerate_(Buildings_and_Research)
- Galaxy Online II Wiki - Guide To Advancing Quickly: https://galaxyonlineii.fandom.com/wiki/Guide_To_Advancing_Quickly
- Galaxy Online II Wiki - Beginner's Guide: https://galaxyonlineii.fandom.com/wiki/Beginner's_Guide
