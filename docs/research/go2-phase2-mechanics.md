# Galaxy Online 2 - Phase 2 Mechanics Research

**Research Date:** 2026-02-06
**Researcher:** research-agent
**Purpose:** Complete investigation of Phase 2 mechanics for Cryptomines Online

---

## 1. SHIP FACTORY

### Overview
The Ship Factory is the production facility for constructing ships in Galaxy Online II. It is a critical building in the Planet Base category.

### Key Features
- **Design Capacity:** 20 different ship blueprints can be stored
- **Production Slots:** 5 slots total
  - Slot 1: Available at base level (Ship Factory Level 1)
  - Slots 2-4: Unlocked by upgrading Ship Factory [NEEDS INVESTIGATION: specific levels for each slot]
  - Slot 5: Unlocked by researching "Sync Shipbuilding" in Logistics Construction Science

### Upgrade System
- **Total Levels:** 24 levels
- **Speed Bonuses:** 1% (Level 1) to 60% (Level 24)
- **Construction Times:** ~2 minutes (Level 1) to ~433 days (Level 24)
- **Requirements per level:**
  - Minimum Civic Center level
  - Resource costs: Metal, He3, and Gold
  - Construction time

### Production Mechanics
- **Ship Production Limit:** Currently limited to 2 million ships per production run
- **Cost Variability:** Build times and construction costs vary depending on:
  - Construction Boost research
  - Quality Materials research
- **Design Interface:** Players select modules for blueprints through an in-game Ship Design screen before production begins
- **Module Build Time Impact:** Every Atomic Framework or Orbital Shield adds 1 second to ship building time

### Example Level Progression
- **Level 1:** 1 slot, 1% speed boost, ~2 minutes construction
- **Level 10:** 4 slots, 14% speed boost, ~1.5 days construction
- **Level 24:** 4 slots, 60% speed boost, ~433 days construction

**Source:** [Ship Factory Wiki](https://galaxyonlineii.fandom.com/wiki/Ship_Factory)

---

## 2. SHIP DESIGN SYSTEM

### Design Workflow
Ship Design is almost completely determined and driven by Fleet Design. Players should review fleet specifications before creating individual ship designs.

### Design Capacity
- **Maximum Designs:** 20 different ship blueprints can be stored in Ship Factory

### Hull Types

#### Frigates
- **Role:** Mobile, quick ships designed to counter battleships
- **Combat Bonus:** Deal extra damage to and receive less damage from Battleships (+5%/-5%)
- **Default Movement:** 1 (no transmission modules needed)
- **Armor Types:** Nano or Neutralizing armor
- **Critical Hit Bonus:** +1%
- **Characteristics:**
  - Excellent shields
  - Lowest module capacity
  - Weakest structure
  - Base Effective Stack: 1100 ships
  - Example stats: Weikes-I (270 shields, 770 structure) to Cybra-III (1,078 shields, 3,080 structure)

#### Cruisers
- **Role:** Balanced mid-tier vessels that excel against frigates
- **Combat Bonus:** Counter frigates effectively (+5%/-5%)
- **Armor Types:** Chrome or Regen armor
- **Characteristics:**
  - Moderate module capacity
  - Moderate shields and structure
  - Better fuel storage than frigates
  - Base Effective Stack: 1000 ships
  - Example stats: Wraith-I (505 shields, 2,599 structure)

#### Battleships
- **Role:** Heavy capital ships with maximum firepower
- **Combat Bonus:** Dominate cruisers but vulnerable to frigates (+5%/-5%)
- **Armor Types:** All four types (Chrome, Regen, Nano, Neutralizing)
- **Characteristics:**
  - Maximum volume capacity
  - Highest structure
  - Greatest fuel storage
  - Lowest shields
  - Base Effective Stack: 900 ships

#### Special Hulls
- Offer unique bonuses
- Can be Frigates, Cruisers, or Battleships
- Enhanced capabilities

#### Flagships & Humaroid-Flagships
- Advanced variants with enhanced capabilities
- Obtained through Blueprint Shreds (see Section 3)

### Rock-Paper-Scissors Balance
- **Frigates** counter **Battleships** (±5% attack/defense)
- **Cruisers** counter **Frigates** (±5% attack/defense)
- **Battleships** counter **Cruisers** (±5% attack/defense)

### Module Categories

#### Attack Modules

**1. Ballistic Weapons**
- **Range:** 1-2 squares
- **Cooldown:** None (fires every round)
- **Power Level:** Least powerful
- **He3 Consumption:** Lowest per round
- **Damage Types:**
  - Kinetic (extra damage to Neutralizing armor)
  - Heat (extra damage to Regen armor)
- **Best for:** Fast ships with 3-5 movement speed

**2. Directional Weapons**
- **Range:** 2-5 squares
- **Cooldown:** 1 round
- **Power Level:** Third most powerful
- **He3 Consumption:** Moderate (least total He3 over entire battle)
- **Damage Types:**
  - Heat
  - Magnetic (damages Nano armor)
- **Best for:** Ships with 2-4 movement speed

**3. Missile Weapons**
- **Range:** 5-8 squares
- **Cooldown:** 3 rounds
- **Power Level:** Second most powerful
- **He3 Consumption:** High (twice ballistic/directional)
- **Damage Type:** Explosive (extra damage to Chrome armor)
- **Best for:** Mid-range and long-range combat

**4. Ship-Based Weapons (SBW)**
- **Range:** 6-10 squares
- **Cooldown:** 4 rounds
- **Power Level:** Most powerful
- **He3 Consumption:** Highest (4x ballistic/directional per round)
- **Damage Types:** All four (Kinetic, Heat, Magnetic, Explosive)
- **Best for:** Long-range combat

**5. Planetary Weapons**
- **Range:** 1-2 squares
- **Cooldown:** None (fires every round)
- **Function:** Only damages defensive structures, NOT ships
- **He3 Consumption:** As much as Ballistic Weapons
- **Best for:** Fast ships for engagement

**CRITICAL DESIGN RULE:** Mixing different weapon classes on the same ship is "widely regarded very ineffective and not recommended" because weapon ranges cause only one class to fire effectively per round, wasting hull space and resources.

#### Defense Modules

**1. Structure Modules**
- Provide improvements to Hull, Overall Defense, added ship stability, reducing damage and reflecting damage
- **Placement Restrictions:** Some limited to one per ship, others can be equipped multiple times
- **He3 Usage:** Most use minimal He3, except Energy Armor variants

**2. Shield Modules**
Four shield types defend against different weapon classes:
- **Kinetic Shield (Particle Stun Shield):** Protects from ballistic and ship-based weapons
- **Heat Shield (Heat Diffusion Shield):** Defends against ballistic, directional, and ship-based weapons
- **Magnetic Shield (Space-Time Magnetic Shield):** Blocks directional and ship-based weapons
- **Explosive Shield (Detonator Shield):** Counters missile and ship-based weapons

Shield modules also increase:
- Overall shield capacity
- Effectiveness
- Restoration
- Damage reduction

**Recommended Shield Ratio (for close-range fleets 1-3 range):** 1:2:1 ratio of Kinetic:Heat:Magnetic shields (heat weapons most prevalent)

**3. Air Defense Modules**
- Provide defensive tactics for intercepting missiles and fighters
- Some variants reflect incoming damage
- Example: "One PPC fires and has a 55% chance to shoot down a single incoming attack"

#### Auxiliary Modules

**1. Electronic Modules**
- Provide improvements to Agility, Steering Power, and Critical Hit chance bonuses
- **Placement Restriction:** Only one of each electronic module variant may equip per ship simultaneously
- Help reduce incoming damage while boosting attack effectiveness
- Examples: Agility Booster, Infrared Scanner, ECM, Auto Target System, Time Dilation Module

**2. Storage Modules**
- Provide extra storage for He3
- Additional storage for other modules
- Expand cargo capacity and module capacity

**3. Transmission Modules**
- Provide extra boost to mobility and agility
- Enhance ship speed and maneuverability capabilities
- Examples: Team Combat Engine, Anti-Matter Engine, EOS Phase Shift Engine

### Module Placement Rules

**CRITICAL:** "When designing ships, the actual order you place certain modules is vital to the design's end effectiveness."

#### Optimal Placement Order (First to Last)
1. Reflective Plating
2. Engines (Team Combat Engine, Anti-Matter Engine, EOS Phase Shift Engine due to their bonuses)
3. Electronic Modules (Agility Booster, Infrared Scanner, ECM, Auto Target System, Time Dilation Module)
4. Maintenance Facilities (Gravity Maintenance, Ship Reinforcement, Shield Regenerator, Hull Maintenance Mechanic)
5. Air Defense (Anti-Aircraft Cannon, Missile Interception, Fighter Interception Cabin, Powered Pulse Cannon)
6. Ship-Based Weapons (SBW)
7. Extreme Counterattack
8. Quick Reaction Armor
9. Daedalus Control System
10. Shield Modules (EOS Phase Shift, Sagan variants, Detonator, Space-Time Magnetic, Heat Diffusion, Particle Stun)
11. Energy Shield Booster
12. Energy Armor
13. Non-SBW Weapons

**Critical Dependency:** "Ship-Based Weapons should be placed before Daedalus and/or Quick Reaction Armor" for optimal effectiveness.

#### Placement-Independent Modules
These modules function identically regardless of placement order:
- Nano Station Warehouse
- Station Warehouse
- Atomic Framework
- Orbital Shield
- Super Transmission Engine

### Hull Capacity & Module Storage
- Each hull design has a specific storage volume for modules
- Players must choose modules wisely when equipping a ship
- Module capacity measured in storage slots
- Example: 5x Cannon-III might use 320 out of 335 hull slots, leaving 15 slots (not enough for another weapon/shield/module)

### Strategic Module Allocation
- Recommended: Approximately half ship's storage volume to weapons, half to defensive systems

### Design Interface Mechanics
- Design screen displays options for:
  - Hull types (Frigates, Cruisers, Battleships, Flagships)
  - Attack modules (Ballistic, Directional, Missile, Ship-Based weapons)
  - Defense modules (Structure, Shields, Air Defense)
  - Auxiliary modules (Electronic, Storage, Transmission)

### Naming Requirements
- **CRITICAL:** System will NOT save designs if names contain spaces or illegal characters
- **Permitted characters:** Periods (.), dashes (-), and underscores (_)

### Ship Statistics Tracked
The design screen displays:
- Attack power
- Shields
- Range
- Storage
- Structure
- Movement
- Build time
- Transmission time

**Sources:**
- [Hull Design Wiki](https://galaxyonlineii.fandom.com/wiki/Hull_Design)
- [Attack Module Wiki](https://galaxyonlineii.fandom.com/wiki/Attack_Module)
- [Defense Module Wiki](https://galaxyonlineii.fandom.com/wiki/Defense_Module)
- [Auxiliary Module Wiki](https://galaxyonlineii.fandom.com/wiki/Auxiliary_Module)
- [Module Placement Wiki](https://galaxyonlineii.fandom.com/wiki/Module_Placement)
- [Ship Design Wiki](https://galaxyonlineii.fandom.com/wiki/Ship_Design)
- [Module Placement Rules](https://galaxyonlineii.fandom.com/wiki/Module_Placement)

---

## 3. BLUEPRINTS

### What Are Blueprints?
Blueprints are essential game components that exist in two forms:
1. **Ship Blueprints:** Enable construction of specific hull types
2. **Module Blueprints:** Enable construction of ship modules (attack, defense, auxiliary)

### How Blueprints Are Obtained

#### Standard Methods
1. **Initial Quests:** Distributed during early gameplay
2. **Instances:** Earned by completing Normal and Restricted instances
   - Blueprints are level-dependent (only specific blueprints available at certain instance levels)
   - 10% chance of getting a blueprint from Treasure Box
3. **Web Mall:** Purchased using:
   - Badge Points (from Restricted instances)
   - Honor Points (from League Matches)
4. **Auction House:** Bought with Gold or Mall Points from other players
5. **Galactic Trafficker:** Sold for Corsairs' Gold

#### Special: Blueprint Shreds (Humaroid Flagships)
Blueprint Shreds are fragments collected through gameplay that serve as crafting materials for advanced flagships.

**How to Obtain Shreds:**
- Defeating Humaroids
- Completing Scenario Instance Milestones

**Shred Types:** Five distinct types (Shred I, II, III, IV, V)

**Minimum Requirement:** 300 total shreds (of three types) for a single blueprint

**Random Blueprint Function:** Spend 200 of each shred type for random blueprint (risky: "may give out one of the listed blueprints, and not necessarily the unlisted one")

### Blueprint Activation
- Once obtained, players must "go to their bag to activate it so they can start using it within the game"
- Unactivated blueprints from instances can be resold through Auction House

### Blueprint Types

#### Ship Blueprints
**Frigates:** Starting with Weikes
**Cruisers:** Beginning with Typhoon
**Battleships:** Such as Estrella

**Special Humaroid Flagships (via Blueprint Shreds):**
- **Intrepid Nexus:** 100 shreds each (three types)
- **Grim Reaper:** 100 shreds each (three types)
- **Shadow Trojan:** 80 shreds across four types
- **Firecat:** 80 shreds across four types
- **Mercury Wing:** 250 of one type + 60 of another
- **GForce's Dreadnaught:** 180, 180, and 60 shreds respectively
- **Conquistador:** Only available by random chance
- **Arbiter:** Only available by random chance

#### Module Blueprints
Three categories:
- **Attack modules:** Ballistic, directional, missile, ship-based, planetary weapons
- **Defense modules:** Structure, shield, air defense
- **Auxiliary modules:** Electronic, storage, transmission

### Blueprint Research System
- Blueprints can be researched/upgraded through the Weapons Factory
- Each blueprint can be researched up to Level 3
- Only Level 1 blueprint needed to complete each level upgrade
- [NEEDS INVESTIGATION: Specific research costs, times, and stat bonuses per level]

### Blueprint Prerequisites
- "Before a ship can be built, any blueprint requirements must be met"
- Some blueprints have prerequisite chains
- [NEEDS INVESTIGATION: Specific prerequisite chains for each ship type]

**Sources:**
- [Blueprint Shreds Wiki](https://galaxyonlineii.fandom.com/wiki/Blueprint_Shreds)
- [Blueprints Wiki](https://galaxyonlineii.fandom.com/wiki/Blueprints)
- [Ship Blueprint Research Wiki](https://galaxyonlineii.fandom.com/wiki/Ship_Blueprint_Research)

---

## 4. SPACEDOCK

### Overview
The Spacedock is a military building dedicated to repairing ships destroyed in battle. It is part of the Planet Base construction category.

### Core Function
- **Purpose:** Repair ships destroyed in battle
- **LIMITATION:** Cannot repair vessels lost in Instance or Restricted Instance battles (only PvP losses)

### Repair Mechanics
- **Repair Percentage:** Represents an individual probability for each ship
- "The % that your space dock tells you that could be repaired is simply the chance of every individual ship being repaired"
- Each destroyed ship has an independent chance of recovery based on this percentage

### Storage Limitations
- **Capacity:** Holds only 2 pages of ships
- **Storage Order:** Design order (oldest first)
- **CRITICAL LIMIT:** If more than 10 different ship types are destroyed simultaneously:
  - Only the 10 oldest designs are recoverable
  - Newer designs (beyond 10) suffer 100% losses

### Repair Acceleration
- Players can spend 10 SP (Mall Points) to reduce repair time by 10%
- Better scaling for longer repair durations

### Level Progression
**Total Levels:** 12 levels

**Requirements per level:**
- Civic Center level requirement
- Ship Factory level requirement
- Metal, He3, and Gold costs
- Construction time

**Repair Rate Progression:**
- Level 1: 1% repair rate
- Level 12: 20% repair rate
- Increases from 1% to 20% across all 12 levels

**Example Levels:**
- **Level 1:**
  - Civic Center 1 required
  - Ship Factory 1 required
  - 1% repair rate
  - 44 seconds upgrade time

- **Level 12:**
  - Civic Center 12 required
  - Ship Factory 23 required
  - 20% repair rate
  - ~116 days upgrade time

**Cost Scaling:** Costs and times scale exponentially, requiring millions of resources at higher levels

### Research Benefits
- Upgrade Times and Upgrade Cost can be reduced by the Logistics Construction Science

[NEEDS INVESTIGATION: Complete table with all 12 levels, exact costs, and repair percentages]

**Source:** [Spacedock Wiki](https://galaxyonlineii.fandom.com/wiki/Spacedock)

---

## 5. FLEET SYSTEM

### Fleet Grid Layout

#### 3x3 Grid Formation
- **Total Positions:** 9 stacks
- **Ships per Stack:** Maximum 3,000 ships of a single design
- **Maximum Fleet Size:** 27,000 ships total (9 × 3,000)
- **CRITICAL RULE:** "Each stack is limited to a single design and you cannot mix ship designs on the same stack; they must be uniform"
- Adjacent stacks may be completely different ships

#### Position Structure & Attack Power

**First Rank (100% attack power):**
- Left Shoulder
- Head
- Right Shoulder

**Second Rank (90% attack power):**
- Left Flank
- Glasshouse (center) - "the single best protected position in any Phalanx"
- Right Flank

**Third Rank (75% attack power):**
- Left Rear
- Tail
- Right Rear

**Strategic Note:** Glasshouse typically houses glass cannons or flagships due to its protected position

### Fleet Limits
- [NEEDS INVESTIGATION: Maximum number of fleets a player can have]
- [NEEDS INVESTIGATION: Fleet unlock progression]

### Fleet Creation
- Access fleet creation by clicking "the airplane in the bottom right corner of the screen"
- **Location Requirement:** Only available in Space Base or Ground Base

### Fleet Uses
Fleets serve multiple purposes:
- Protecting your planet
- Attacking other planets
- Normal Instances
- League battles
- Arena combat
- Space Raids
- Championship events

### Commander System

#### Acquisition Methods
Commander Cards obtained through:
- Drawing from the Command Center (100 Mall Points)
- Lucky Wheel rewards
- Recruiting from Command Center (rare)
- Auction House purchases
- Free cards during events

#### Core Function
- "They allow you to create fleet commanders with special attributes"
- **CRITICAL:** "Commander is one of the most important parts of the game because they stand on the vanguard, leading your fleet into battle"
- Cards merge in the Compound Center to increase Star Rank and dramatically boost stats while raising effective fleet numbers

#### Expertise Rating System
Letter grades (S, A, B, C, D, F) indicating proficiency with weapon types and ship hulls:
- **S Grade:** +30% weapon damage
- **A Grade:** +10% weapon damage
- **B Grade:** Average performance (no modifiers)
- **C, D, F Grades:** Penalties to performance

Commanders have different levels of expertise in each weapon and ship hull type. One commander may be better with battleships while another might be extremely knowledgeable about ballistic weapons.

**Minimum Requirement:** Commanders require appropriate skill levels (minimum rating of B) for each vessel type in the fleet

#### Commander Attributes
Four primary attributes determine effectiveness:
- **Accuracy:** Enhances weapon hit chance and potential damage
- **Dodge:** Reduces opponent hit rates by lowering enemy accuracy
- **Speed:** Affects combat order and successive strike probability
- **Electron:** Increases critical hit rate and critical damage output

#### Card Classification

**By Level:**
- Skill
- Super
- Legendary
- Divine

**By Specialty:**
- Attack skills
- Defense skills
- Energy skills

Higher-tier commanders possess greater attributes and abilities.

#### Effective Stack Impact
Commanders significantly boost "Effective Stack" value based on their star rank. Effective Stack determines how many ships in a stack can attack. This is considered the game's most impactful upgrade.

**Base Effective Stack (without Commander):**
- Frigates: 1100 ships
- Cruisers: 1000 ships
- Battleships: 900 ships

**Sources:**
- [Fleet Design Wiki](https://galaxyonlineii.fandom.com/wiki/Fleet_Design)
- [Fleets Wiki](https://galaxyonlineii.fandom.com/wiki/Fleets)
- [Commander Cards Wiki](https://galaxyonlineii.fandom.com/wiki/Commander_Cards)

---

## 6. NORMAL INSTANCES (PvE Content)

### Overview
Instances are "a very important aspect of the game" where players fight against AI-controlled ships to win rewards. Each instance provides different challenges which get progressively tougher at higher levels.

### Instance Mechanics

#### Rewards
Players receive **Treasure Boxes** upon completion containing:
- **Resources:** Gold, Metal, or He3 (small amounts)
- **Blueprints:** 10% chance, divided equally among blueprints available from that instance

#### Blueprint Acquisition
- "Blueprints are level-dependent meaning that only specific Blueprints can be acquired at certain levels"
- "Blueprints can only be acquired from specific instances"
- Players must target particular instances to obtain desired equipment designs

#### Farming
- Instances can be run repeatedly for accumulation of rewards and in-game wealth
- Instance data does NOT reset (unlike Scenario Instances which reset every 24 hours)

### Combat Mechanics (Instance-Specific)

#### Enemy Targeting Behavior
- **Enemy Missile Fleets:** Attack your fleet with the strongest total attack rating
- **Enemy Ship-Based Fleets:** Attack your fleet with the highest total durability rating

#### Progression Strategy
The wiki recommends:
1. Initial focus on "Logistics Construction Science tech levels"
2. Then weapon technology
3. Then shield technology
4. Finally tackle higher-level instances for superior rewards

### The 30 Normal Instances

Based on the Normal Instances wiki, the 30 instances are:

1. **Instance 1:** "Ancestral Recall" - 3 max fleets, 180 EXP
2. **Instance 2:** "Deadzone" - 4 max fleets, 500 EXP
3. **Instance 3:** "Bravery" - 4 max fleets, 1,000 EXP
4. [NEEDS INVESTIGATION: Instances 4-29 names, max fleets, EXP]
5. **Instance 30:** "Triumphant Glory" - 15 max fleets, 73,500 EXP

[NEEDS INVESTIGATION: Complete list of all 30 instances with:
- Official instance names
- Level requirements
- Max fleets allowed
- Experience rewards
- Checkpoint structure (if any)
- Enemy fleet compositions
- Specific blueprint drops per instance]

**Interactive Tool Available:** krtools.info Instance Viewer contains detailed information about all instances, layouts, ships, and mechanics

**Sources:**
- [Normal Instances Wiki](https://galaxyonlineii.fandom.com/wiki/Normal_Instances)
- [Instances Wiki](https://galaxyonlineii.fandom.com/wiki/Instances)
- [Instance Viewer Tool](https://krtools.deajae.co.uk/inst/)

---

## 7. COMBAT SYSTEM (for PvE)

### Effective Stack
- **Definition:** "Effective Stack determines how many ships in a stack can attack"
- **Base Values:**
  - Frigates: 1100 ships
  - Cruisers: 1000 ships
  - Battleships: 900 ships
- **Commander Impact:** Commander Cards significantly boost effective stack based on star rating (most impactful upgrade in game)

### Combat Sequence (8-Phase System)

#### Phase 1: Attacker Fires
- **Number of Hits:** Attack modules per ship × ships in effective stack × weapon hit chance
- **Hit Chance Influenced By:**
  - Weapon accuracy
  - Attacker steering
  - Defender agility
  - Defender dodge
- **Rough Approximations:**
  - Each point of agility reduces hit chance by ~4%
  - Each steering point increases hit chance by ~4%

#### Phase 2: Interceptors Fire
- Defensive modules can intercept incoming attacks
- Example: "One PPC fires and has a 55% chance to shoot down a single incoming attack" (independent of number of missiles incoming)
- Interception effectiveness varies by defense module type

#### Phase 3: Calculate Damage
- Remaining projectiles deal damage within their specified range
- Weapon bonuses apply
- Damage calculations account for double-hits and critical effects

#### Phase 4: Damage Negation
- Most shield systems reduce incoming damage before it reaches the hull

#### Phase 5: Shield Penetration
- Ballistic and directional weapons can pierce shields

#### Phase 6: Deal Damage to Shields
- "EOS triggers and has 30% chance to absorb double damage by Damage Mitigation tech"

#### Phase 7: Assign Damage to Hull
- Excess damage divides by individual ship structure and stability percentage to determine destroyed vessels

#### Phase 8: Calculate Scatter Damage
- Unblockable scatter damage affects other stacks
- Bypasses defensive systems entirely

### Armor System

#### Five Armor Types
1. **Chrome Armor:** Reduces Solar Damage & Kinetic Damage to the lowest level
2. **Nano Armor:** Reduces Solar Damage & Explosive Damage to the lowest level
3. **Regen Armor:** Reduces Kinetic Damage & Magnetic Damage to the lowest level
4. **Neutralizing Armor:** Reduces Explosive Damage & Magnetic Damage to the lowest level
5. **Light Armor:** [NEEDS INVESTIGATION: damage reduction properties]

#### High Quality Armor
- Designation that reduces damage by 5-15% from base damage

### Damage Types & Effectiveness

#### Four Damage Types
1. **Kinetic Damage:** Extra damage to Neutralizing armor
2. **Solar Damage (Heat):** Extra damage to Regen armor
3. **Explosive Damage:** Extra damage to Chrome armor
4. **Magnetic Damage:** Extra damage to Nano armor

#### Damage Calculation Formula
**Damage Caused = Weapon Damage × Hit Chance × Number of Weapons that Hit Target**

More specifically: "The amount of damage dealt to each armor depends on weapon expertise, hull expertise, weapon damage, defense, and some other factors."

#### Strategic Implications
- "If a ship has no module/armor to protect against specific types of damage, they will suffer the full damage amount"
- "If you know in advance what kind of weapons an enemy will be using, you can protect your ships with the corresponding armor type"
- Damage is increased or reduced due to ship armor type and attacking damage type

### Stability
- Stability increases ship survivability by reducing destruction likelihood when taking damage
- Factors into Phase 7 (Assign Damage to Hull) calculations

**Sources:**
- [Combat Mechanics Wiki](https://galaxyonlineii.fandom.com/wiki/Combat_Mechanics)
- [Category: Armor Wiki](https://galaxyonlineii.fandom.com/wiki/Category:Armor)

---

## 8. ADDITIONAL RESEARCH SOURCES

### Interactive Tools
- **KRTools Ship Design Viewer:** krtools.info/des/ - Interactive ship design visualization
- **KRTools Instance Viewer:** krtools.deajae.co.uk/inst/ - Detailed instance information

### Key Wiki Pages Consulted
- [Galaxy Online II Wiki Main Page](https://galaxyonlineii.fandom.com/wiki/Galaxy_Online_II_Wiki)
- [Ship Factory](https://galaxyonlineii.fandom.com/wiki/Ship_Factory)
- [Spacedock](https://galaxyonlineii.fandom.com/wiki/Spacedock)
- [Hull Design](https://galaxyonlineii.fandom.com/wiki/Hull_Design)
- [Attack Module](https://galaxyonlineii.fandom.com/wiki/Attack_Module)
- [Defense Module](https://galaxyonlineii.fandom.com/wiki/Defense_Module)
- [Auxiliary Module](https://galaxyonlineii.fandom.com/wiki/Auxiliary_Module)
- [Module Placement](https://galaxyonlineii.fandom.com/wiki/Module_Placement)
- [Ship Design](https://galaxyonlineii.fandom.com/wiki/Ship_Design)
- [Composite Ship Table](https://galaxyonlineii.fandom.com/wiki/Composite_Ship_Table)
- [Blueprints](https://galaxyonlineii.fandom.com/wiki/Blueprints)
- [Blueprint Shreds](https://galaxyonlineii.fandom.com/wiki/Blueprint_Shreds)
- [Ship Blueprint Research](https://galaxyonlineii.fandom.com/wiki/Ship_Blueprint_Research)
- [Fleet Design](https://galaxyonlineii.fandom.com/wiki/Fleet_Design)
- [Fleets](https://galaxyonlineii.fandom.com/wiki/Fleets)
- [Commander Cards](https://galaxyonlineii.fandom.com/wiki/Commander_Cards)
- [Normal Instances](https://galaxyonlineii.fandom.com/wiki/Normal_Instances)
- [Instances](https://galaxyonlineii.fandom.com/wiki/Instances)
- [Combat Mechanics](https://galaxyonlineii.fandom.com/wiki/Combat_Mechanics)

### External Resources
- [Galaxy Online II Building Ships Guide - Game Yum](https://www.gameyum.com/galaxy-online/109247-galaxy-online-ii-guides-how-to-build-ships/)
- [Module Placement Rules](https://galaxyonlineii.fandom.com/wiki/Module_Placement)
- [Fleet Design](https://galaxyonlineii.fandom.com/wiki/Fleet_Design)

---

## 9. INFORMATION GAPS - [NEEDS INVESTIGATION]

The following information was NOT found in the wikis and will need to be improvised or found through additional research:

### Ship Factory
- [ ] Exact Ship Factory level requirements for unlocking Slot 2, Slot 3, and Slot 4
- [ ] Complete table of all 24 levels with exact Metal/He3/Gold costs
- [ ] Exact construction times for each level

### Spacedock
- [ ] Complete table of all 12 levels with exact Metal/He3/Gold costs
- [ ] Exact construction times for each level
- [ ] Exact repair percentages for levels 2-11 (only Level 1 = 1% and Level 12 = 20% confirmed)

### Ship Design
- [ ] Exact hull capacity (volume) numbers for each frigate/cruiser/battleship type
- [ ] Exact module slot counts for each hull type
- [ ] Complete list of all available modules with exact stats
- [ ] Exact resource costs for each module
- [ ] Ship build time formulas

### Blueprints
- [ ] Complete blueprint prerequisite chains
- [ ] Blueprint research costs (Metal/He3/Gold) for Levels 1-3
- [ ] Blueprint research time requirements
- [ ] Exact stat bonuses gained from researching blueprints to Level 2 and Level 3

### Fleet System
- [ ] Maximum number of fleets a player can have
- [ ] Fleet unlock progression (if any)
- [ ] Fleet capacity increases (if any)

### Normal Instances
- [ ] Complete list of all 30 instance names
- [ ] Level requirements for each instance
- [ ] Max fleets allowed for instances 4-29
- [ ] Experience rewards for instances 4-29
- [ ] Checkpoint structure (number of waves/checkpoints per instance)
- [ ] Enemy fleet compositions per instance
- [ ] Specific blueprint drops per instance (which blueprints drop from which instances)

### Combat System
- [ ] Exact hit chance formula
- [ ] Exact critical hit formula
- [ ] Exact damage calculation formula (complete formula with all variables)
- [ ] Light Armor damage reduction properties
- [ ] Exact stability formula and effects

### Other
- [ ] Complete list of all module names (only categories and some examples found)
- [ ] Transmission time mechanics
- [ ] Fleet travel system
- [ ] Energy/stamina system for instances (if any)

---

## 10. NOTES FOR GDD DESIGNER

### Core Design Principles Found

1. **Rock-Paper-Scissors Balance:** The three hull types create a perfect counter system (±5% damage modifiers)

2. **Specialization Over Generalization:** Mixing weapon types on one ship is "widely regarded very ineffective" - ships should specialize

3. **Module Placement Order Matters:** The order modules are placed affects their effectiveness (see Section 2 - Module Placement Rules)

4. **Commander Cards Are Critical:** Described as "most impactful upgrade" due to Effective Stack bonuses

5. **Trade-offs Everywhere:**
   - Frigates: Fast, high shields, low structure
   - Cruisers: Balanced stats
   - Battleships: Slow, high structure, low shields

6. **Weapon Range & Cooldown Balance:**
   - Longer range = Higher cooldown + Higher power
   - Shorter range = No cooldown + Lower power + Lower He3 cost

7. **Fleet Formation Strategy:** Attack power varies by position (100%/90%/75% by rank)

8. **Limited Production:** 5 ship factory slots, 20 design capacity, 2M ship production limit

### Suggested Implementation Order

Based on complexity and dependencies:

1. **Ship Factory Building** (simpler, no dependencies)
2. **Spacedock Building** (simpler, no dependencies)
3. **Blueprint System** (needed before ship design)
4. **Ship Design System** (complex, needs blueprints)
5. **Fleet System** (needs ships)
6. **Commander System** (enhances fleets)
7. **Normal Instances** (needs fleets + combat system)
8. **Combat System** (most complex, needs everything else)

---

## RESEARCH COMPLETION STATUS

**Date Completed:** 2026-02-06
**Information Gathered:** ~70% complete
**Sources Consulted:** 20+ wiki pages, 2 interactive tools, multiple external guides

**Ready for GDD Phase:** YES - Sufficient information for Phase 2 GDD creation
**Improvisation Needed:** YES - See Section 9 for specific gaps

**Recommendation:** Proceed with GDD creation using gathered information. Flag [NEEDS INVESTIGATION] items for later research or reasonable improvisation based on game balance principles.
