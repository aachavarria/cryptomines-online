# GDD Requirements Checklist

**Generated:** 2026-02-13
**Source:** game-design-document.md v3.3 + final-scope.md
**Scope:** Only features marked IN SCOPE in final-scope.md

---

## 1. Buildings System

### 1.1 Ground Base - Resource Buildings
- [ ] REQ-B001: Metal Collector produces Metal (GDD 2.2.1)
  - Max count: 8 per planet
  - Max level: 24
  - Footprint: 2x2
- [ ] REQ-B002: He3 Extractor produces He3 (GDD 2.2.1)
  - Max count: 8 per planet
  - Max level: 24
  - Footprint: 2x2
- [ ] REQ-B003: Residential Area produces Gold (highest output) (GDD 2.2.1)
  - Max count: 8 per planet
  - Max level: 24
  - Footprint: 2x2
- [ ] REQ-B004: Resource Warehouse stores resources (GDD 2.2.1)
  - Multiple allowed
  - Max level: 24
  - Footprint: 3x2

### 1.2 Ground Base - Core Buildings
- [ ] REQ-B005: Civic Center is main hub (GDD 2.2.2, 2.2.6)
  - Max level: 12
  - Footprint: 3x3
  - Formula: Lv1 = 550 Metal, 480 He3, 600 Gold, 5 min
  - Scaling: ~3x resource cost per level, ~2.9x build time
- [ ] REQ-B006: Technology Center enables research (GDD 2.2.2, 2.2.7)
  - Max level: 12
  - Footprint: 3x2
  - Effect: 3% research time reduction per level (max 36% at Lv12)
  - Formula: `EffectiveResearchTime = BaseResearchTime * (1 - TechCenterLevel * 0.03)`
- [ ] REQ-B007: Radar detects incoming attacks (GDD 2.2.2)
  - Max level: [NEEDS RESEARCH]
  - Footprint: 1x1
  - Scope: IN

### 1.3 Ground Base - Military Buildings
- [ ] REQ-B008: Ship Factory constructs ships (GDD 2.2.3, 2.4.8)
  - Max level: 24
  - Footprint: 3x2
  - Design slots: 20 max
  - Production slots: 5 (5th requires Sync Shipbuilding tech)
  - Max production: 2,000,000 ships at a time
  - Speed bonus: 1-60% with levels 1-24
- [ ] REQ-B009: Spacedock manages fleets (GDD 2.2.3)
  - Max level: 12
  - Footprint: 3x2
  - Scope: IN (repair placeholder)
- [ ] REQ-B010: Command Center recruits commanders (GDD 2.2.3, 2.2.8)
  - Max level: 12
  - Footprint: 2x2
  - Cooldown reduction: 3:00:00 at Lv1 → 2:20:00 at Lv5
  - Max commanders: 60 (at player level 71+)
- [ ] REQ-B011: Weapon Research Center develops blueprints (GDD 2.2.3, 2.3.10)
  - Max level: 12
  - Footprint: 2x2
  - Research slots: tied to building level
  - Research time reduction: 3% per level (max 36% at Lv12)
- [ ] REQ-B012: Recycling Plant scraps ships (GDD 2.2.3, 2.6.3)
  - Max level: [NEEDS RESEARCH]
  - Footprint: 2x2
  - Recovery: 70% of total build cost (hull + modules)
  - Scope: IN
- [ ] REQ-B013: Compound Center merges commanders (GDD 2.2.2, 2.6.3)
  - Max level: [NEEDS RESEARCH]
  - Footprint: 2x2
  - Function: Star Rank merging
  - Scope: IN

### 1.4 Space Base - Defense Buildings
- [ ] REQ-B014: Space Station is main orbital structure (GDD 2.2.5)
  - Max level: 100 (for RBPs)
  - Footprint: 3x3
  - Rule: Must stay within 1 level of Civic Center
- [ ] REQ-B015: Meteor Star orbital defense (GDD 2.2.5)
  - Max level: 10
  - Footprint: 1x1
  - HP Range: 40K (Lv1) - 20.48M (Lv10)
  - Combat integration: YES
- [ ] REQ-B016: Particle Cannon energy defense (GDD 2.2.5)
  - Max level: 10
  - Footprint: 1x2
  - Range: 8-17 tiles
  - Damage: 10K-450K
  - Combat integration: YES
- [ ] REQ-B017: Anti-Aircraft Gun anti-air defense (GDD 2.2.5)
  - Max level: 10
  - Footprint: 1x1
  - Range: 30 tiles
  - Damage: 50K-500K
  - HP Range: 16K (Lv1) - 8.19M (Lv10)
  - Combat integration: YES
- [ ] REQ-B018: Thor's Cannon heavy defense (GDD 2.2.5)
  - Max level: 10
  - Footprint: 2x2
  - Damage: 50K-1.35M
  - Combat integration: YES
- [ ] REQ-B019: Celestial Base advanced facility (GDD 2.2.5)
  - Max level: 10
  - Footprint: 2x2
  - Combat integration: YES

### 1.5 Building Grid & Placement
- [ ] REQ-B020: 20x20 isometric grid (400 tiles) (GDD 2.2.9)
  - Detail: Buildings placed on (col, row) coordinates
  - Footprint validation: All tiles in footprint must be empty and in bounds
- [ ] REQ-B021: Multi-tile building footprints (GDD 2.2.9)
  - 3x3: Civic Center, Space Station
  - 3x2: Ship Factory, Resource Warehouse, Spacedock, Technology Center
  - 2x2: Most military/admin buildings
  - 1x2: Particle Cannon
  - 1x1: Radar, Meteor Star, Anti-Aircraft Gun

### 1.6 Building System Rules
- [ ] REQ-B022: Construction slots (GDD 2.2.10)
  - Default: 2 concurrent construction slots
  - Construction Card: +3 slots for 72 hours
- [ ] REQ-B023: Civic Center <-> Space Station mutual dependency (GDD 2.2.10)
  - Must stay within 1 level of each other
- [ ] REQ-B024: Civic Center level gates upgrades (GDD 2.2.10)
  - Detail: Most buildings require Civic Center level N to upgrade to N

---

## 2. Resources System

### 2.1 Primary Resources
- [ ] REQ-R001: Metal resource (GDD 2.1.1)
  - Producer: Metal Collector (max 8)
  - Relative rate: ~0.4x Gold
  - Use: Ship manufacturing, building upgrades
- [ ] REQ-R002: He3 resource (GDD 2.1.1)
  - Producer: He3 Extractor (max 8)
  - Relative rate: ~0.5-0.6x Gold
  - Use: Ship fuel, building upgrades
- [ ] REQ-R003: Gold resource (GDD 2.1.1)
  - Producer: Residential Area (max 8)
  - Relative rate: 1.0x (highest)
  - Use: Research, trading, upgrades

### 2.2 Resource Storage & Collection
- [ ] REQ-R004: Warehouse accumulation (GDD 2.1.3, final-scope Phase A)
  - Resources accumulate in warehouse (auto-production)
  - Fields: warehouse_metal, warehouse_he3, warehouse_gold
  - Cap: storage_capacity
  - Worker: Updates warehouse every 5 minutes
- [ ] REQ-R005: Manual collection (GDD 2.1.3)
  - Endpoint: Transfer warehouse → player resources
  - Recommended: Collect every 12 hours or less
- [ ] REQ-R006: PvP loot (GDD 2.1.3, 2.5.7)
  - Winner receives 20% of loser's resources
  - Excludes warehouse contents
  - Max cap: 1,000,000 (hardcoded, cargo capacity NOT implemented)

---

## 3. Research System

### 3.1 Seven Science Trees
- [ ] REQ-T001: Total of 111 techs across 7 trees (GDD 2.3.1, final-scope)
  - Ballistics Science: 16 techs
  - Ship Defense Science: 19 techs
  - Logistics Construction Science: 11 techs
  - Directional Science: 18 techs
  - Missile Science: 18 techs
  - Ship-Based Science: 17 techs
  - Planetary Defense Science: 12 techs

### 3.2 Ballistics Science (16 techs)
- [ ] REQ-T002: Ballistics (Base) Lv 1-10 (GDD 2.3.2)
  - Effect: +5% ballistic damage per level
  - Lv1: 541 Gold, 1:57 | Lv10: 18,808 Gold, 39:18:21
- [ ] REQ-T003: Ballistic Malice Lv 1-5 (requires Ballistics Lv 3)
  - Effect: +1% critical hit rate per level
- [ ] REQ-T004: Ballistic Crackdown Lv 1-2 (requires Ballistics Lv 3)
  - Effect: +10% critical damage per level
- [ ] REQ-T005: Steady Control Lv 1-5 (requires Ballistics Lv 6, Malice Lv 3)
  - Effect: Reduces weapon space by 2-10%
- [ ] REQ-T006: Precise Ballistics Lv 1-5 (requires Ballistics Lv 8, Steady Lv 3)
  - Effect: +1-5% hit rate
- [ ] REQ-T007: Shield Penetration Lv 1 (requires Malice Lv 5, Crackdown Lv 1, Precise Lv 1)
  - Effect: 15% shield bypass damage
  - Cost: 154,616 Gold, 4:15:00
- [ ] REQ-T008: Depleted Uranium Bomb Lv 1-3 (requires Crackdown Lv 2, Penetration Lv 1)
  - Effect: +10-30% vs Neutral armor, +1-3% vs Light armor
- [ ] REQ-T009: Fire Bomb Research Lv 1-3 (requires Crackdown Lv 2, Penetration Lv 1)
  - Effect: +10-30% vs Regen armor, +1-3% vs Light armor
- [ ] REQ-T010: Improved Penetration Lv 1-3 (requires DU Bomb Lv 1, Fire Bomb Lv 1)
  - Effect: +1-3% vs Light armor, 3-10% shield pen chance
- [ ] REQ-T011: Victory Rush Lv 1 (requires DU Lv 3, Fire Lv 3, Imp.Pen. Lv 3)
  - Effect: Range-dependent damage: 220%/180%/150%/120% + 5% crit rate/damage
  - Cost: 1,918,521 Gold, 62:20:00
- [ ] REQ-T012: Ballistic Scattering Lv 1-5 (requires Ballistics Lv 10, Malice Lv 5, Steady Lv 5, Precise Lv 3)
  - Effect: 5-25% scattering damage to adjacent ships
- [ ] REQ-T013: Improved Ballistic Scattering Lv 1-3 (requires Precise Lv 5, Scattering Lv 3)
  - Effect: +8-25% scattering rate
- [ ] REQ-T014: Hop Bomb Research Lv 1-5 (requires Scattering Lv 5, Imp. Scatter Lv 2)
  - Effect: 3-15% chance to deal 100% weapon damage as scatter

### 3.3 Ship Defense Science (19 techs)
- [ ] REQ-T015: Ship Defense Tech (Base) Lv 1-2 (GDD 2.3.3)
  - Lv1: +2% base shield/structure/agility/defense, +5% stability
  - Lv2: +5% base stats, +10% stability
- [ ] REQ-T016: Shield Research Lv 1-5 (requires Base Lv 1)
  - Effect: +1-5% base shield
- [ ] REQ-T017: Energy Diffusion Lv 1-3 (requires Base Lv 2, Shield Lv 3)
  - Effect: Each shield module reduces damage by 1-3
- [ ] REQ-T018: Penetration Resistance Lv 1-2 (requires Shield Lv 5, Diffusion Lv 1)
  - Effect: -3-7% enemy shield penetration chance
- [ ] REQ-T019: Augment Shield Lv 1-3 (requires Diffusion Lv 2, Pen.Res. Lv 1)
  - Effect: +6-20% base shield
- [ ] REQ-T020: Restoration Lv 1-2 (requires Diffusion Lv 3, Aug.Shield Lv 2)
  - Effect: +30-60% shield restore per round, +1-2% interception chance
- [ ] REQ-T021: Augment Absorption Lv 1-2 (requires Shield Lv 5, Diffusion Lv 1)
  - Effect: Shield modules reduce damage by 2-5
- [ ] REQ-T022: Energy Conservation Lv 1-3 (requires Diffusion Lv 2, Aug.Abs. Lv 1)
  - Effect: +3-10% chance absorb damage without He3
- [ ] REQ-T023: Electronic Barrier Lv 1-2 (requires Diffusion Lv 3, E.Cons. Lv 2)
  - Effect: Reflect 5-10% damage before shields drop to 0
- [ ] REQ-T024: Damage Mitigation Lv 1-3 (requires Aug.Shield Lv 3, Restoration Lv 2, E.Cons. Lv 3, E.Barrier Lv 2)
  - Effect: 10-30% absorb double damage, 15-45% lower collateral
- [ ] REQ-T025: Ship Structural Analysis Lv 1-2 (requires Base Lv 1)
  - Effect: +1-2% base structure
- [ ] REQ-T026: Ship Reinforcement Lv 1-3 (requires Base Lv 2, Analysis Lv 3)
  - Effect: Each structure module reduces damage by 1-3
- [ ] REQ-T027: Resilience Lv 1-2 (requires Analysis Lv 5, Reinforce Lv 1)
  - Effect: -3-7% enemy structure penetration chance
- [ ] REQ-T028: Structure Improvement Lv 1-3 (requires Reinforce Lv 2, Resilience Lv 1)
  - Effect: +6-20% base structure
- [ ] REQ-T029: Fast Repair Lv 1-2 (requires Reinforce Lv 3, Struct.Imp. Lv 2)
  - Effect: +30-60% structure restore per round
- [ ] REQ-T030: Reaction Armor Improvement Lv 1-2 (requires Analysis Lv 5, Reinforce Lv 1)
  - Effect: Structure modules reduce damage by 2-5
- [ ] REQ-T031: Defense Improvement Lv 1-3 (requires Reinforce Lv 2, Reaction Lv 1)
  - Effect: +3-10% absorb damage without He3
- [ ] REQ-T032: Reflection Mastery Lv 1-2 (requires Reinforce Lv 3, Def.Imp. Lv 2)
  - Effect: Reflect 5-10% damage before structure drops to 0
- [ ] REQ-T033: Stability Mastery Lv 1-3 (requires Struct.Imp Lv 3, Fast Repair Lv 2, Def.Imp Lv 3, Reflection Lv 2)
  - Effect: 10-30% absorb double damage, 15-45% lower collateral

### 3.4 Logistics Construction Science (11 techs)
- [ ] REQ-T034: Concurrent Construction Lv 1 (GDD 2.3.4)
  - Effect: +1 construction slot
  - Cost: 1,000 Gold, 0:00:20
- [ ] REQ-T035: Construction Boost Lv 1-10 (requires Concurrent Lv 1)
  - Effect: +1-15% building construction speed
  - Lv1: 2,400G, 0:08:00 | Lv10: 1,463,112G, 56:55:01
- [ ] REQ-T036: Quality Materials Lv 1-10 (requires Construction Boost Lv 3)
  - Effect: -1-15% building resource costs
  - Lv1: 1,200G, 0:04:00 | Lv10: 512,250G, 28:27:30
- [ ] REQ-T037: Ship Building Boost Lv 1-10 (no prereq)
  - Effect: +1-15% shipbuilding speed
  - Lv1: 990G, 0:03:18 | Lv10: 1,463,112G, 81:17:02
- [ ] REQ-T038: Ship Building Logistics Lv 1-10 (requires Ship Building Boost Lv 2)
  - Effect: -1-15% ship construction resource costs
  - Lv1: 1,386G, 0:04:37 | Lv10: 2,048,358G, 113:47:52
- [ ] REQ-T039: Sync Shipbuilding Lv 1 (requires Ship Building Logistics Lv 4)
  - Effect: Adds 5th shipbuilding slot
  - Cost: 174,000 Gold, 9:40:00
- [ ] REQ-T040: Repair Technology Lv 1-10 (requires Sync Shipbuilding Lv 1)
  - Effect: +1-10% ship repair percentage
  - Lv1: 5,310G, 0:17:42 | Lv10: 2,266,722G, 125:55:45
- [ ] REQ-T041: High Yield Mining Lv 1-10 (no prereq)
  - Effect: +1-10% Metal output
  - Lv1: 1,740G, 0:05:48 | Lv10: 742,764G, 41:15:53
- [ ] REQ-T042: High Yield Chemistry Lv 1-10 (requires High Yield Mining Lv 2)
  - Effect: +1-10% He3 output
  - Lv1: 2,400G, 0:08:00 | Lv10: 1,024,500G, 56:55:00
- [ ] REQ-T043: High Yield Investing Lv 1-10 (requires High Yield Chemistry Lv 2)
  - Effect: +1-10% Gold output
  - Lv1: 3,570G, 0:11:54 | Lv10: 1,523,952G, 84:39:50
- [ ] REQ-T044: Expand Capacity Lv 1-7+ (requires High Yield Investing Lv 4)
  - Effect: +50,000-350,000 warehouse storage per level
  - Lv1: 3,540G, 0:11:48

### 3.5 Directional Science (18 techs)
- [ ] REQ-T045: Optics (Base) Lv 1-10 (GDD 2.3.5)
  - Effect: +5% directional weapon damage per level (max 50%)
- [ ] REQ-T046: Directional Malice Lv 1-5
- [ ] REQ-T047: Directional Accuracy Lv 1-5
- [ ] REQ-T048: Eagle Eye Lv 1-2
- [ ] REQ-T049: Energy Penetration Lv 1
- [ ] REQ-T050: Pierce Lv 1-5
- [ ] REQ-T051: Radiative Interference Lv 1-5
- [ ] REQ-T052: Improved Pierce Lv 1-2
- [ ] REQ-T053: Energy Accumulation Lv 1-3
- [ ] REQ-T054: Electronic Interference Lv 1-3
- [ ] REQ-T055: Piercing Crit Lv 1
- [ ] REQ-T056: Weakness Detection Lv 1-3
- [ ] REQ-T057: Particle Impact Tech Lv 1-3
- [ ] REQ-T058: Magnetic Impact Lv 1-3
- [ ] REQ-T059: Dynamic Impairment Lv 1-3

### 3.6 Missile Science (18 techs)
- [ ] REQ-T060: Missile Theory (Base) Lv 1-10 (GDD 2.3.6)
  - Effect: +4% missile damage per level (max 40%)
- [ ] REQ-T061: Missile Accuracy Lv 1-5
  - Effect: +2% hit rate per level (max 10%)
- [ ] REQ-T062-T077: [Remaining 16 Missile Science techs from GDD 2.3.6]
  - Scope: IN (all 18 techs)

### 3.7 Ship-Based Science (17 techs)
- [ ] REQ-T078-T094: [All 17 Ship-Based Science techs from GDD 2.3.7]
  - Scope: IN (all 17 techs)

### 3.8 Planetary Defense Science (12 techs)
- [ ] REQ-T095-T106: [All 12 Planetary Defense Science techs from GDD 2.3.8]
  - Scope: IN (all 12 techs)

### 3.9 Research System Rules
- [ ] REQ-T107: Research prerequisites validation (GDD 2.3.9)
  - Detail: Cannot start tech without prerequisite techs completed
- [ ] REQ-T108: Research slots (GDD 2.3.9)
  - Detail: Technology Center level determines concurrent research slots
- [ ] REQ-T109: Auto-completion worker (final-scope Phase A)
  - Detail: Worker checks `research_finish_at` and auto-completes techs
- [ ] REQ-T110: Tech effects application (final-scope Phase A)
  - Detail: Apply tech bonuses (production, unlock modules, etc.)
- [ ] REQ-T111: Research time reduction by Tech Center (GDD 2.2.7)
  - Formula: `EffectiveResearchTime = BaseResearchTime * (1 - TechCenterLevel * 0.03)`

---

## 4. Blueprint Research System

### 4.1 Blueprint Tier Upgrades
- [ ] REQ-BP001: Upgrade hulls/modules from tier 1 → 2 → 3 (GDD 2.3.10, final-scope)
  - Pattern: `weikes_i` → `weikes_ii` → `weikes_iii`
  - Modules: same `name`, tier 1 → 2 → 3
- [ ] REQ-BP002: Weapon Research Center controls research slots (GDD 2.3.10)
  - Max level: 12
  - Research slots: tied to building level
- [ ] REQ-BP003: Blueprint research costs (GDD 2.3.10, final-scope)
  - Formula: `baseCost = 10000 × targetLevel` (Metal/He3/Gold)
  - Time: 3600 seconds × targetLevel (1 hour per level)
- [ ] REQ-BP004: Blueprint research levels (GDD 2.3.10)
  - Level 1 (Base): Blueprint acquired, base stats, 1.0x cost
  - Level 2: WRC Lv 6 required, +10% stats, 2.0x cost
  - Level 3: WRC Lv 10 required, +25% stats, 5.0x cost
- [ ] REQ-BP005: Auto-complete worker (final-scope Phase A)
  - Detail: Auto-complete when `research_finish_at` reached
- [ ] REQ-BP006: Tier validation in ship design (final-scope Phase A)
  - Detail: Only show unlocked tiers in ship design panel

---

## 5. Ships System

### 5.1 Hull Types
- [ ] REQ-S001: Frigate hull type (GDD 2.4.1)
  - Effective stack: 1,100
  - Bonus: +5% vs Battleship
  - Penalty: -5% vs Cruiser
  - Armor types: Nano, Chrome, Regen, Neutralizing
- [ ] REQ-S002: Cruiser hull type (GDD 2.4.1)
  - Effective stack: 1,000
  - Bonus: +5% vs Frigate
  - Penalty: -5% vs Battleship
- [ ] REQ-S003: Battleship hull type (GDD 2.4.1)
  - Effective stack: 900
  - Bonus: +5% vs Cruiser
  - Penalty: -5% vs Frigate
- [ ] REQ-S004: Total of 75 hulls (final-scope)
  - 30 Frigates
  - 30 Cruisers
  - 15 Battleships

### 5.2 Ship Modules
- [ ] REQ-S005: Total of 97 modules across 11 categories (final-scope)
  - Ballistic weapons
  - Directional weapons
  - Missile weapons
  - Ship-Based weapons
  - Planetary weapons
  - Shield modules
  - Structure modules
  - Air Defense modules (PPC: 55% intercept chance)
  - Electronic modules
  - Storage modules
  - Transmission modules (TCE, AME, Super Transmission Engine)

### 5.3 Ship Stats
- [ ] REQ-S006: Ship stats calculation (GDD 2.4.3)
  - Shield: Absorbs damage before hull
  - Structure: Hull HP
  - Stability: Reduces effective damage to hull
  - Defense: Damage reduction percentage
  - Installation Slots: Module capacity
  - Agility: Dodge chance, reduces enemy hit rate
  - Movement (MOV): Fleet travel speed
  - Storage: Cargo capacity

### 5.4 Ship Design
- [ ] REQ-S007: Ship Factory design slots (GDD 2.4.8, final-scope)
  - Max 20 ship designs stored
  - Hull + module assignment
  - Blueprint required to build
- [ ] REQ-S008: Ship Factory production (GDD 2.4.8)
  - 5 production slots (5th requires Sync Shipbuilding tech)
  - Max production: 2,000,000 ships at a time
  - Speed bonus: 1-60% with levels 1-24

### 5.5 Blueprints
- [ ] REQ-S009: Blueprint acquisition (GDD 2.4.7, final-scope)
  - Total: 62 blueprints (25 hulls + 37 modules)
  - Primary source: Instances (PvE)
  - Treasure Box: 10% chance of blueprint
  - Also: Auction House (OUT), Corp bonuses, events
- [ ] REQ-S010: Blueprint unlock flow (final-scope Module 9)
  - Item → Inventory → Use → Unlock in player_blueprints
  - Pattern: Currently direct, must pass through inventory

---

## 6. Fleet System

### 6.1 Fleet Composition
- [ ] REQ-F001: Fleet grid 3x3 (GDD 2.4.4)
  - Max 9 stacks per fleet
  - Max 3,000 ships per stack
  - Max 27,000 ships per fleet
- [ ] REQ-F002: Grid positions & attack power (GDD 2.4.4)
  - First Rank: 100% attack
  - Second Rank: 90% attack
  - Third Rank: 75% attack
  - Glasshouse (center-middle): Most protected
  - Shoulders: Most vulnerable
- [ ] REQ-F003: Fleet speed (GDD 2.4.4)
  - Formula: Speed of slowest ship in fleet

### 6.2 Fleet Formations
- [ ] REQ-F004: Formation types (GDD 2.4.5)
  - Phalanx: All 9 stacks filled
  - Diamond: 5 stacks (center + 4 adjacent)
  - Battle Line: 6 stacks (first 2 ranks)
  - Skirmish: Sparse
  - Tee Forward: T-shape
  - Enfilade: Side-focused
  - Tee Reverse: Reversed T-shape

### 6.3 Fleet Commands
- [ ] REQ-F005: Targeting commands (GDD 2.4.6)
  - Max/Min attack power
  - Max/Min durability
  - Closest proximity
  - By commander rank

---

## 7. Combat System

### 7.1 Combat Resolution (8-Phase)
- [ ] REQ-C001: Phase 1 - Attacker Fires (GDD 2.5.1)
  - Formula: `Attacks = (Attack modules per ship) × (Ships in effective stack) × (Hit chance)`
  - Hit chance influenced by: Weapon accuracy, Steering (~4%/pt), Accuracy (~1%/12pt), Agility (~4%/pt), Commander dodge
- [ ] REQ-C002: Phase 2 - Interceptors Fire (GDD 2.5.1)
  - PPC: 55% chance to shoot down each incoming attack
  - Intercept modules counter missiles and ship-based weapons
- [ ] REQ-C003: Phase 3 - Calculate Damage (GDD 2.5.1)
  - Weapon damage: random within min-max range
  - Double-hits and critical hits applied
  - Modified by armor type vs damage type
- [ ] REQ-C004: Phase 4 - Damage Negation (GDD 2.5.1)
  - Non-EOS shields reduce damage
  - Reduces scatter and piercing effects
- [ ] REQ-C005: Phase 5 - Shield Penetration (GDD 2.5.1)
  - Ballistic and directional weapons can pierce shields
  - Penetration chance modified by tech research
- [ ] REQ-C006: Phase 6 - Deal Damage to Shields (GDD 2.5.1)
  - EOS Phase Shift: 30% chance absorb double damage
  - Only effective stack triggers EOS defensively
  - Combined shield total from all 3,000 ships absorbs damage
- [ ] REQ-C007: Phase 7 - Assign Damage to Hull (GDD 2.5.1)
  - Formula: `Ships Destroyed = floor(Remaining Damage / (Individual Ship Structure × Stability%))`
- [ ] REQ-C008: Phase 8 - Calculate Scatter Damage (GDD 2.5.1)
  - Based on weapon type, technology, Sandora module
  - Scatter damage CANNOT be absorbed or mitigated by defenses

### 7.2 Battle Duration
- [ ] REQ-C009: Combat rounds (GDD 2.5.2)
  - Minimum: 20 rounds + number of fleets/buildings present
  - Maximum: 99 rounds
  - Only ballistic weapons fire every round; others require cooldown

### 7.3 Ship Type Advantage
- [ ] REQ-C010: Rock-paper-scissors mechanics (GDD 2.5.3)
  - Frigate → +5% vs Battleship
  - Cruiser → +5% vs Frigate
  - Battleship → +5% vs Cruiser

### 7.4 Armor Types vs Damage Types
- [ ] REQ-C011: Armor effectiveness matrix (GDD 2.5.4)
  - Nano: neutral
  - Chrome: Strong vs Kinetic
  - Regen: Strong vs Light Energy
  - Neutralizing: neutral
  - Tech bonuses: DU Bomb (+10-30% vs Neutral), Fire Bomb (+10-30% vs Regen)

### 7.5 Commander Impact
- [ ] REQ-C012: Weapon Expertise grades (GDD 2.5.5)
  - S: +30% damage
  - A: +10% damage
  - B: 0% (neutral)
  - C: -10% damage
  - D: -30% damage
- [ ] REQ-C013: Ship Expertise grades (GDD 2.5.5)
  - S: +10% dealt, -10% received
  - A: +5% dealt, -10% received
  - B: 0% / 0%
  - C: -5% dealt, +10% received
  - D: -10% dealt, +10% received

### 7.6 Combat Losses
- [ ] REQ-C014: Loss modes (GDD 2.5.6)
  - Normal/Restricted Instances: Ships lost YES, He3 lost YES
  - Trial/Constellation Instances (OUT): Ships lost NO, He3 lost YES
  - League/Arena/Championships (OUT): Ships lost NO, He3 lost NO
  - PvP (Attack Neighbors): Ships lost YES, He3 lost YES

---

## 8. PvP Combat

### 8.1 PvP Mechanics
- [ ] REQ-PVP001: Attack neighbor planets (GDD 2.5.7, final-scope)
  - Scope: IN
- [ ] REQ-PVP002: Fleet travel time (GDD 2.5.7)
  - Uses SP (Space Points)
- [ ] REQ-PVP003: Combat resolution (final-scope Phase C)
  - Reuses 8-phase combat system
- [ ] REQ-PVP004: Loot 20% resources (GDD 2.5.7, final-scope)
  - Winner receives 20% of loser's resources
  - Excludes warehouse contents
  - Max cap: 1,000,000 (cargo capacity NOT implemented)
- [ ] REQ-PVP005: Attack cooldown (final-scope Module 6)
  - 5-minute cooldown between attacks
- [ ] REQ-PVP006: Radar building (GDD 2.2.2, final-scope)
  - Detects incoming attacks
  - Provides warning
- [ ] REQ-PVP007: Defense fleets (final-scope Phase C)
  - Assign fleets to defend planet
- [ ] REQ-PVP008: PvP combat reports (final-scope Phase C)
  - Round-by-round detailed log
- [ ] REQ-PVP009: Truce cards (final-scope Module 9)
  - Truce Card: 12h protection
  - Adv Truce Card: 72h protection

---

## 9. Commander System

### 9.1 Commander Rarity (Simplified)
- [ ] REQ-CMD001: 3 rarity tiers (final-scope Module 4)
  - Common: Lowest power
  - Skill: Low power
  - Super: Medium power
  - Legendary/Divine: OUT OF SCOPE
- [ ] REQ-CMD002: Gacha recruitment rates (GDD 2.6.3, final-scope)
  - 50% Common
  - 35% Skill
  - 15% Super

### 9.2 Commander Attributes
- [ ] REQ-CMD003: 4 stat attributes (GDD 2.6.2)
  - Accuracy: Increases weapon hit chance
  - Dodge: Reduces opponent hit rate
  - Speed: Determines attack order; affects successive strike chance
  - Electron: Increases Critical Hit Rate and Critical Damage

### 9.3 Star Rank System
- [ ] REQ-CMD004: Star Rank merging (GDD 2.6.3, final-scope)
  - Auto-merge duplicate commanders (modern gacha system)
  - Star Rank increases Effective Stack
  - Building: Compound Center
  - Max 60 commanders (at player level 71+)

### 9.4 Recruitment Mechanics
- [ ] REQ-CMD005: Recruitment methods (GDD 2.6.6, final-scope)
  - Free Recruitment: 3-hour cooldown, uses Gold (NO Mall Points)
  - Quick Recruitment: uses Gold (NO Mall Points)
  - Commander Cards: Item → Inventory → Use → Unlock commander
  - Auction House: OUT OF SCOPE
  - Events: IN SCOPE
- [ ] REQ-CMD006: Max commanders (GDD 2.2.8)
  - 60 commanders max (at player level 71+)

### 9.5 Commander Enhancements (OUT OF SCOPE)
- [ ] REQ-CMD007: Commander Skills (Blue/Red/Green) - OUT (final-scope)
- [ ] REQ-CMD008: Gems - OUT (final-scope)
- [ ] REQ-CMD009: Bionic Chips - OUT (final-scope)
- [ ] REQ-CMD010: Wounded/Dead states - OUT (commanders immortal) (final-scope)

---

## 10. Recycling Plant

### 10.1 Ship Scrapping
- [ ] REQ-REC001: Select ships to scrap (final-scope Module 8)
  - Scope: Simple implementation
- [ ] REQ-REC002: Recover 70% resources (GDD 2.6.3, final-scope)
  - 70% of total build cost (hull + modules)
- [ ] REQ-REC003: Delete ships from database (final-scope Module 8)
  - Permanent deletion

---

## 11. Inventory System

### 11.1 Resource Packs (8 items)
- [ ] REQ-INV001: Gold Pack (+30k) (final-scope Module 9)
- [ ] REQ-INV002: Advanced Gold Pack (+100k)
- [ ] REQ-INV003: Primary Metal Pack (+50k)
- [ ] REQ-INV004: Junior Metal Pack (+150k)
- [ ] REQ-INV005: Senior Metal Pack (+300k)
- [ ] REQ-INV006: Primary He3 Pack (+50k)
- [ ] REQ-INV007: Junior He3 Pack (+150k)
- [ ] REQ-INV008: Senior He3 Pack (+300k)

### 11.2 Resource Boosts (6 items)
- [ ] REQ-INV009: Construction Card (+3 building slots, 72h)
- [ ] REQ-INV010: MVP Tool (+20% all production/build/repair, 7 days)
- [ ] REQ-INV011: Extra Tax (+30% Gold, 12h)
- [ ] REQ-INV012: Adv Extra Tax (+100% Gold, 24h)
- [ ] REQ-INV013: Metal Mining Boost (+30% Metal, 12h)
- [ ] REQ-INV014: He3 Mining Boost (+30% He3, 12h)

### 11.3 Battle Items (3 items)
- [ ] REQ-INV015: SP Card (+10 Space Points instant)
- [ ] REQ-INV016: Truce Card (12h protection)
- [ ] REQ-INV017: Adv Truce Card (72h protection)

### 11.4 Consumables
- [ ] REQ-INV018: Blueprint items (final-scope Module 9)
  - Item → Inventory → Use → Unlock in player_blueprints
- [ ] REQ-INV019: Commander Cards (final-scope Module 9)
  - Item → Inventory → Use → Unlock commander

### 11.5 Inventory Implementation
- [ ] REQ-INV020: Table `player_items` (player_id, item_key, quantity) (final-scope)
- [ ] REQ-INV021: Inventory panel UI (final-scope)
- [ ] REQ-INV022: Use logic per item type (final-scope)
  - 5 use functions implemented

---

## 12. World Chat

### 12.1 Basic World Chat
- [ ] REQ-CHAT001: Single world channel (final-scope Module 10)
  - NO corps chat, NO private messages
- [ ] REQ-CHAT002: Send message (final-scope)
  - Rate limiting: 3s between messages
- [ ] REQ-CHAT003: View recent messages (final-scope)
  - Display recent 100 messages
  - Pagination
- [ ] REQ-CHAT004: Profanity filter (final-scope)
  - 100 words filter
  - Unit tests: 19 subtests pass

### 12.2 Chat Features (OUT OF SCOPE)
- [ ] REQ-CHAT005: Friends system - OUT
- [ ] REQ-CHAT006: Mail system - OUT
- [ ] REQ-CHAT007: Private messages - OUT
- [ ] REQ-CHAT008: Corps chat - OUT (Corps moved to Phase 4, but corps chat still OUT)
- [ ] REQ-CHAT009: Loudspeaker item - OUT (chat is free)

---

## 13. Corps System (Phase 4 - Feb 13, 2026)

### 13.1 Corp Basics
- [ ] REQ-CORP001: Create corp (final-scope Phase 4)
  - Alliance Center building required
  - Corp creation endpoint
- [ ] REQ-CORP002: Join/leave corp (final-scope Phase 4)
  - Join endpoint
  - Leave endpoint
- [ ] REQ-CORP003: Corp roles (final-scope Phase 4)
  - Leader: Full control
  - Officer: Moderate control
  - Member: Basic access

### 13.2 Corp Donations & Wealth
- [ ] REQ-CORP004: Donation system (GDD 2.8.3, final-scope Phase 4)
  - Max 200 pts/day (2M resources)
  - Donate Metal/He3/Gold
- [ ] REQ-CORP005: Corp wealth (final-scope Phase 4)
  - Accumulates from donations
  - Used to upgrade RBPs
- [ ] REQ-CORP006: Corp levels (final-scope Phase 4)
  - Corp level increases with wealth
  - Corp level = max RBPs controlled

### 13.3 Corp Features (OUT OF SCOPE)
- [ ] REQ-CORP007: Corp Mall - OUT
- [ ] REQ-CORP008: Corp Warehouse - OUT
- [ ] REQ-CORP009: Corp Merging Center - OUT
- [ ] REQ-CORP010: Pirate Planets - OUT
- [ ] REQ-CORP011: Galactic Wars - OUT

---

## 14. Galaxy Map & RBPs (Phase 4 - Feb 13, 2026)

### 14.1 Galaxy Map
- [ ] REQ-GAL001: 7x7 zone grid (GDD 2.7.1, final-scope Phase 4)
  - 49 zones total
  - 1 RBP per zone (center of each zone)
  - Spacing: 60 movement squares apart
- [ ] REQ-GAL002: Galaxy Map panel (final-scope Phase 4)
  - HTML/CSS grid panel
  - Display all 49 RBPs
  - Show ownership, level, bonuses

### 14.2 RBP Bonuses
- [ ] REQ-GAL003: RBP bonus scaling (GDD 2.7.2, final-scope Phase 4)
  - Formula: GO2-faithful
  - Lv 1-10: 5% base + 0.5%/level = 10% at Lv10
  - Lv 11-20: 1% per level = 20% at Lv20
  - Lv 21-30: 1.5% per level = 35% at Lv30
  - Lv 31-40: 2% per level = 55% at Lv40
  - Lv 100: 280% total
- [ ] REQ-GAL004: Bonus application (GDD 2.7.2)
  - Resource production
  - Research speed
  - Shipbuilding speed

### 14.3 RBP Defense
- [ ] REQ-GAL005: Initial NPC defenses (GDD 2.7.3)
  - 5 fleets of 4,500-7,200 Level 6 ships each
- [ ] REQ-GAL006: Defensive structures (GDD 2.7.3)
  - Meteor Stars: 63, HP 40K (Lv1) - 20.48M (Lv10)
  - Particle Cannons: 8, Range 8-17 tiles, Damage 10K-450K
  - Anti-Aircraft Guns: 12, HP 16K (Lv1) - 8.19M (Lv10), Range 30 tiles, Damage 50K-500K
  - Thor's Cannons: 5, Damage 50K-1.35M
- [ ] REQ-GAL007: Fleet capacity by level (GDD 2.7.3)
  - Lv 1: 10 defense fleets
  - Lv 5: 10 defense fleets
  - Lv 10: 15 defense fleets
  - Lv 15: 20 defense fleets
  - +5 per 5 levels
  - Lv 100: 105 defense fleets

### 14.4 RBP Conquest
- [ ] REQ-GAL008: Only Corps can attack RBPs (GDD 2.7.4, final-scope Phase 4)
  - Not individuals
- [ ] REQ-GAL009: Conquest mechanics (GDD 2.7.4, final-scope Phase 4)
  - Reuses 8-phase combat engine
  - 99 combat rounds max
  - Fleets cannot be recalled mid-battle
- [ ] REQ-GAL010: Protection timers (GDD 2.7.4)
  - 72-hour protection after capture
  - 24-hour battle phase when vulnerable
  - Successful takeover resets to 72 hours
- [ ] REQ-GAL011: Conquest reset (GDD 2.7.4)
  - On conquest: defenses reset to Level 1
  - Space Station retains level
- [ ] REQ-GAL012: Multi-Corp attacks (GDD 2.7.4)
  - Planet goes to corp with most kills from single fleet type
  - Scoring: 1 pt per ship destroyed, per structure destroyed, per Space Station destroyed

### 14.5 RBP Control & Upgrades
- [ ] REQ-GAL013: Control limits (GDD 2.7.5)
  - Corps can control one planet per Corp Level
  - Lv 10 corp = max 10 RBPs
- [ ] REQ-GAL014: Upgrading RBPs (GDD 2.7.6)
  - Uses Corp Wealth (from member donations)
  - Space Station upgradeable to Level 100
  - Defenses upgradeable to Level 10
  - Total wealth for Lv 100 Station: 24,902,439

---

## 15. Quest System

### 15.1 Main Quests
- [ ] REQ-Q001: 22 active main quests (final-scope)
  - Total designed: 28
  - 6 deferred (OUT OF SCOPE): main_02, main_18, main_19, main_20, main_21, main_22
- [ ] REQ-Q002: Auto-progress integration (final-scope)
  - Buildings
  - Research
  - Blueprints
  - Ships
  - Commanders
  - Combat

### 15.2 Side Quests
- [ ] REQ-Q003: 12 side quests (final-scope)
  - Auto-progress integration

### 15.3 Daily Quests
- [ ] REQ-Q004: 6 daily quests (final-scope)
  - Points: 10-70 per quest
  - Tier rewards: bronze/silver/gold/diamond
- [ ] REQ-Q005: Daily system (final-scope)
  - Daily reset
  - Points accumulation
  - Tier reward thresholds

---

## 16. Production Polish

### 16.1 UI/UX Polish
- [ ] REQ-POL001: Better error messages (final-scope Module 11)
  - All hooks use GameContext SET_ERROR
- [ ] REQ-POL002: Loading states on all buttons (final-scope Module 11)
  - LoadingButton standardized across all panels
- [ ] REQ-POL003: Complete tooltips (final-scope Module 11)
  - Techs explain what they do
  - Modules explain stats
  - Buildings explain bonuses
- [ ] REQ-POL004: Smooth animations (final-scope Module 11)
  - Scope: IN
- [ ] REQ-POL005: Tutorial tooltips (final-scope Module 11)
  - First-time hints
- [ ] REQ-POL006: Sound effects (final-scope Module 11)
  - Build complete
  - Combat
  - etc.

### 16.2 Polish (OUT OF SCOPE)
- [ ] REQ-POL007: Mobile responsive polish - OUT (desktop only)

---

## 17. OUT OF SCOPE SYSTEMS

### 17.1 Major Systems Cut
- [ ] REQ-OUT001: Trading Center - OUT (final-scope)
- [ ] REQ-OUT002: Decorative Buildings (12 types) - OUT (final-scope)
- [ ] REQ-OUT003: Friends System - OUT (final-scope)
- [ ] REQ-OUT004: Mail System - OUT (final-scope)
- [ ] REQ-OUT005: League/Championship - OUT (final-scope)
- [ ] REQ-OUT006: Automated Tests - OUT (final-scope)
- [ ] REQ-OUT007: Mobile Responsive - OUT (desktop only) (final-scope)

### 17.2 Items/Features Cut
- [ ] REQ-OUT008: Planet Transformation Packs (6 items) - OUT
- [ ] REQ-OUT009: Commander Enhancement Items (Memory Chip, Merge Chip, etc.) - OUT
- [ ] REQ-OUT010: Healing/Revival Cards - OUT (commanders immortal)
- [ ] REQ-OUT011: Galaxy Transfer - OUT
- [ ] REQ-OUT012: Loudspeaker - OUT (chat is free)
- [ ] REQ-OUT013: Corps Certificate - OUT
- [ ] REQ-OUT014: Passport (Extra Restricted Instance) - OUT
- [ ] REQ-OUT015: Constellation Pass/Scroll Chests - OUT
- [ ] REQ-OUT016: Treasure Boxes - OUT (instance rewards direct)
- [ ] REQ-OUT017: Legendary/Divine Commander Tiers - OUT
- [ ] REQ-OUT018: Commander Skills (Blue/Red/Green) - OUT
- [ ] REQ-OUT019: Gems/Bionic Chips - OUT

---

## Summary Statistics

**Total Requirements Extracted:** 400+ (exact count TBD after full GDD parse)

**Systems IN SCOPE (11 major):**
1. Buildings (22 types, 24 requirements)
2. Resources (6 requirements)
3. Research (111 techs, 111+ requirements)
4. Blueprint Research (6 requirements)
5. Ships (10 requirements)
6. Fleets (5 requirements)
7. Combat (14 requirements)
8. PvP (9 requirements)
9. Commanders (6 requirements in scope, 4 out)
10. Recycling Plant (3 requirements)
11. Inventory (22 requirements)
12. World Chat (9 requirements, 5 out)
13. Corps (6 requirements, 5 out)
14. Galaxy Map & RBPs (14 requirements)
15. Quests (5 requirements)
16. Production Polish (7 requirements, 1 out)

**Systems OUT OF SCOPE (8 major):**
1. Trading Center
2. Decorative Buildings
3. Friends/Mail
4. League/Championship
5. Tests
6. Mobile
7. Commander enhancements (gems/chips/skills)
8. Advanced items (19 item types cut)

**Total Requirements OUT:** 20+ explicitly marked

---

## Notes

- **[NEEDS RESEARCH] Tags:** Some requirements have incomplete data from GDD (exact costs, max levels, complete lists)
- **Phase 4 (Feb 13, 2026):** Corps & Galaxy Map implemented but NOT yet QA tested
- **Formulas Included:** All key formulas from GDD preserved (build costs, scaling, combat resolution, RBP bonuses)
- **GO2 Fidelity:** All requirement names, values, mechanics faithful to Galaxy Online 2 as documented in GDD

---

**Next Steps:**
1. QA team: Use this checklist to validate backend implementation
2. QA team: Use this checklist to validate frontend implementation
3. Fill in [NEEDS RESEARCH] gaps with wiki research
4. Create test cases for each REQ-XXX item
5. Run WebMCP test flows against each system
