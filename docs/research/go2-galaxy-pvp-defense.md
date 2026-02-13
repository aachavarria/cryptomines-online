# Galaxy Online 2 - Galaxy Map, PvP, Defense & Fleet Movement Research

**Research Date:** February 13, 2026
**Purpose:** Document individual player mechanics for galaxy navigation, PvP attacks, defense, and fleet movement
**Scope:** Individual player mechanics only (excludes Corps/Alliance features, RBPs, League/Championship)

---

## 1. Galaxy Map (Individual Player View)

### Overview
The galaxy map is accessed by clicking the "Galaxy" icon in the game interface. Players can browse planets and identify targets from this view.

### Player Location & Coordinates
- Players have **coordinates** that identify their planet location in the galaxy
- Coordinates are visible when viewing a player's planet in galaxy view
- **[UNCONFIRMED]** Exact coordinate system format (zones, sectors, systems) - not found in wiki

### Finding Other Players
Players can find targets to attack through the galaxy view:
1. Click the Galaxy icon to open galaxy view
2. Browse/navigate the map to locate other players
3. Click on a planet to view player information:
   - Player ID
   - Player coordinates
   - Corps affiliation
   - League ranking

### Radar for Monitoring Movement
- Click Galaxy view icon → green radar symbol → 5 tabs at top left
- Shows ships in transit (incoming and outgoing fleets)
- Radar building level determines how much information is revealed about incoming attacks (see Defense section)

### Search/Scan Mechanics
**[UNCONFIRMED]** - No specific "search" or "scan" feature documented in the wiki beyond browsing the galaxy map manually.

**Sources:**
- [Beginner FAQ | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Beginner_FAQ)
- [Attacking Neighbors (PvP) | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Attacking_Neighbors_(PvP))

---

## 2. Space Station (Player's Own)

### Role & Purpose
- The Space Station is "the only entrance your enemies have to attack your planet"
- Provides access to:
  - Celestial Bases
  - Subsidiary Territories
  - Instance Maps
- Functions as the planet's primary defensive hub

### Space Station vs Ground Base
**Two separate building locations exist:**

**Ground Base Buildings:**
- Civic Center
- Metal Collector
- He3 Extractor
- Residential Area (Gold production)
- Resource Warehouse
- Ship Factory
- Spacedock
- Technology Center
- Trading Center
- Alliance Center
- Compound Center
- Galaxy Transporter
- Landscaping structures (Casino Resort, Beacon, Monument, Fountain, Library, Theater, Park, College, Hospital, Shopping Center, Statue, Santa Sculpture)

**Space Base (Orbital) Buildings:**
- Space Station (main structure)
- Meteor Star (defense)
- Particle Cannon (defense)
- Anti-Aircraft Gun (defense)
- Thor's Cannon (defense)
- Celestial Base

### Space Station Stats
| Level | HP | Notes |
|-------|-----|-------|
| 1 | 100 | Minimum |
| 12 | 1,200 | Maximum |

**[UNCONFIRMED]** Civic Center relationship - mentioned that they should be upgraded together, but exact "must stay within 1 level" rule not confirmed in sources.

### Celestial Base
- Serves as the planet's "Headstone"
- Allows resource collection from subsidiary territories
- Can be repaired

**Sources:**
- [Orbital Bases | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Orbital_Bases)
- [Walkthrough | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Walkthrough)
- [Guide To Advancing Quickly | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Guide_To_Advancing_Quickly)

---

## 3. Defense Buildings

### Overview
"All enemy ships MUST defeat all Orbital Defenses before they can start attacking your Space Station."

Even having a few high-level orbital defenses provides precious time for allies to retaliate and defend.

### 1. Meteor Star

**Purpose:** Defensive structure that blocks fleet movement paths

**Mechanics:**
- Prevents fleets from traveling in a particular direction
- Enemy fleets must destroy Meteor Stars before proceeding to their target
- Does NOT deal direct damage - purely a blocking/delaying structure

**Stats:**

| Level | HP | Defense | Build Time | Cost (Metal/He3/Gold) |
|-------|-----|---------|-----------|----------------------|
| 1 | 200 | 50 | 20 sec | 50/45/48 |
| 2 | 600 | 50 | - | - |
| 3 | 2,000 | 50 | - | - |
| 9 | 400,000 | 50 | 24+ hours | 220,652/198,586/211,826 |
| 10 | 30,000,000 | 120 | 71+ hours | 5,437,269/4,893,542/5,219,779 |
| 11 | 60,000,000 | 120 | - | - |
| 12 | 90,000,000 | 120 | 604 hours | Very high |

### 2. Particle Cannon

**Purpose:** Long-range continuous-fire weapon with highest range but weakest attack

**Mechanics:**
- "The highest range available" among defense buildings
- Continuous fire (no cooldown mentioned)
- Effective against shields when lowered

**Stats:**

| Level | Attack | HP | Range (squares) | Build Time | Cost (Metal/He3/Gold) |
|-------|--------|-----|----------------|-----------|----------------------|
| 1 | 5,000 | 10,000 | 5 | 45 sec | Low |
| 2-8 | Scaling | Scaling | 6-12 | Scaling | Scaling |
| 9 | 70,000 | 1,400,000 | 13 | 94+ hours | High |
| 10 | 1,000,000 | 5,000,000 | ~14-15 | 71+ hours | 10,581,962/9,523,665/10,158,111 |
| 11 | 1,500,000 | 10,000,000 | ~16-17 | Very long | Very high |
| 12 | 2,000,000 | 20,000,000 | ~18 | Very long | 102,678,894+ Metal |

### 3. Anti-Aircraft Gun

**Purpose:** High-damage area-effect weapon

**Mechanics:**
- "The most powerful attack second to the Thor Cannon"
- **Cooldown: 1 round**
- Attacks "one square in all directions" at point of impact (area-of-effect)
- Strikes multiple ships simultaneously within its attack range

**Stats:**

| Level | Attack | HP | Range (squares) | Build Time |
|-------|--------|-----|----------------|-----------|
| 1 | 25,000 | 20,000 | 3 | Short |
| 2-8 | Scaling | Scaling | 4-10 | Scaling |
| 9 | 225,000 | 3,200,000 | 11 | Long |
| 10 | 2,000,000 | 10,000,000 | ~12-13 | Very long |
| 11 | 2,500,000 | 30,000,000 | ~14-15 | Very long |
| 12 | 3,000,000 | 50,000,000 | ~16-17 | Very long |

### 4. Thor's Cannon

**Purpose:** Ultimate attack weapon with most powerful attack and radius

**Mechanics:**
- "The ultimate attack weapon"
- Most powerful attack among all defenses
- Largest attack radius
- **Cooldown: 2 rounds**
- Attacks **ALL enemy ships within its attack range simultaneously**

**Stats:**

| Level | Attack | HP | Range (squares) | Build Time | Cost (Metal) |
|-------|--------|-----|----------------|-----------|-------------|
| 1 | 10,000 | 200,000 | 10 | Short | Low |
| 2-8 | Scaling | Scaling | 11-17 | Scaling | Scaling |
| 9 | 800,000 | 7,000,000 | 18 | Very long | High |
| 10 | 3,000,000 | 30,000,000 | ~19-20 | Very long | Very high |
| 11 | 4,500,000 | 60,000,000 | ~21-22 | Very long | Very high |
| 12 | 6,000,000 | 90,000,000 | ~23-24 | Very long | 238,000,000+ |

### Research Bonus
**Emplacement Mastery** research provides up to **10% bonus to attack values** for all defense cannons.

### Defense Building Notes
- When a player is successfully attacked and defeated, **all defenses reset automatically at no cost**
- Reconstruction time scales with building level
- Level 10+ Thor cannons take **several days** to rebuild
- Build times and costs affected by Construction Boost and Quality Materials research

### Defense Strategy: Corner Wall
- Build Meteor Stars and Particle Cannons around each corner
- Prevents enemy fleet movement
- Player's own fleets can position behind the wall
- Very effective: enemy fleets must destroy the wall first to reach your ships

**Sources:**
- [Combat Mechanics | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Combat_Mechanics)
- [Composite Orbital Defenses Table | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Composite_Orbital_Defenses_Table)
- [Defense Strategies | Galaxy Online II Wiki](https://pylonnexus.com/galaxyonlineiifandomcom/wiki/Defense_Strategies.html)

---

## 4. PvP Attack Mechanics (Individual Player)

### Attack Flow

**Step 1: Target Selection**
1. Click Galaxy icon
2. Locate target planet
3. Click planet to view info (ID, coordinates, corps, league rank)
4. Select "Send to attack with your fleets"

**Step 2: Scout Attack (Recommended)**
- **Cannot see enemy space station defenses unless they are your friend**
- "The only way is to send a scout attack"
- Most players send "a small expedition fleet of low value ships before they send in their main fleet"
- Scout attack reveals defensive structures and fleet composition

**Step 3: Send Main Fleet**
- Select fleet(s) to send
- Confirm attack

**Step 4: Travel Time**
- Fleet travels from your planet to target
- **[UNCONFIRMED]** Exact travel time formula not documented
- Travel time likely based on distance and fleet speed
- **Return travel time = 50% of outbound journey time**

**Step 5: Combat**
- Battle initiates upon arrival
- See Combat Mechanics section below

**Step 6: Return Home**
- If victorious, fleets "will return home on their own"
- Fleets appear at "one of the 4 wormholes located at the corners of your space base"
- Return time = half of outbound travel time
- Loot delivered via in-game mail

### Space Points (SP) Cost

**SP Usage:**
- Required for: Sending ships to attack planets, accelerating ship repairs, harvesting resources from subsidiary territories, assisting friends (costs 0 SP despite requiring at least 1)
- SP bar **resets to full every day at midnight** (server time)
- Can use SP cards to get more SP if you run out

**[UNCONFIRMED]** Exact SP cost per attack - not documented in wiki. Scenario instances charge **1 SP** + variable Gold fee, but PvP attack cost not specified.

### Combat Sequence

**Minimum Duration:** 20 rounds
**Maximum Duration:** 20 + additional rounds based on defender fleets, attacker fleets, and defensive structures

**Critical Combat Constraint:**
- "If you cannot complete the battle before those rounds complete, you will lose the battle"
- Helium-3 (He3) fuels ship attacks
- If He3 depletes, ships sit idle through remaining rounds
- Players "always run the risk of running out of power before you run out of targets"

**Combat Order:**
1. Attacker fleet vs Defender orbital defenses (Meteor Stars, Particle Cannons, Anti-Aircraft Guns, Thor's Cannons)
2. Attacker fleet vs Defender fleets (if any stationed)
3. Attacker fleet vs Space Station (if defenses destroyed)

**8-Phase Combat Round (per round):**
1. Attacker fires
2. Interceptors activate (defense modules shoot down incoming attacks)
3. Damage calculation
4. Damage negation (shields, defensive systems)
5. Shield penetration (certain weapons bypass shields)
6. Hull damage (destroys ships based on structure and stability)
7. Scatter damage (unmitigated bonus damage affects secondary stacks)
8. Defender fires (same phases 1-7)

**Hit Chance Mechanics:**
- Every 1 agility reduces hit chance by 4%
- Every 1 steering increases hit chance by 4%
- Weapon type, accuracy stats, agility, and dodge stats all affect hit chance

**Effective Stack:**
- Determines how many ships in a stack can attack per round
- Base effective stack: Frigates (1,100), Cruisers (1,000), Battleships (900)
- Commander cards significantly increase effective stack based on star rating
- A 3,000-ship stack does NOT attack all at once - only effective stack size attacks

**Ship Destruction:**
- Damage dealt to shields that exceeds shield capacity is assigned to ship's structure
- Divided by (individual ship structure × stability percentage), rounded down = number of ships destroyed
- **Stability** makes it harder for ships to be destroyed (see Stability section)

### Victory Conditions

**Attacker Wins If:**
- All defender fleets destroyed AND
- All orbital defenses destroyed AND
- Space Station destroyed
- Battle completed within round limit

**Attacker Loses If:**
- Cannot complete battle before round limit expires
- All attacking ships destroyed
- Runs out of He3 fuel before completing objectives

### Loot Mechanics

**Successful Attack:**
- "If you win the battle, you will receive a mail that contains 20% of the losing planet's harvested resources, along with a battle report"
- Loot = **20% of loser's resources NOT including what is in their warehouse**
- Resources in the Resource Warehouse are **protected from looting**

**[UNCONFIRMED]** Specific caps on loot amounts - wiki mentions 20% rule but doesn't specify per-resource caps or total caps.

**Battle Report Contents:**
- How many ships were sent
- How many ships were destroyed
- Name of ship design that shot down the most ships
- Name of player who owns that design
- Number of ships shot down by that design

### Ships After Combat

**Destroyed Ships:**
- Ships lost in PvP can be repaired via Spacedock
- Spacedock repair percentage depends on:
  - Spacedock level (1% at Lv1 → 20% at Lv12)
  - Repair Technology research level
- Repair is **probabilistic**: each ship has independent chance of recovery equal to displayed percentage
- "Think of it as a dice with 100 sides, and there is a green mark on the amount of sides equal to the % that it tells you could be repaired"
- Spacedock stores **only 2 pages of ships** (~20 designs)
- Ships stored by design age (oldest first)
- If more than 10 different ship types destroyed, **only the 10 oldest remain recoverable**
- Newer designs beyond 10 suffer **100% losses**
- Ships lost in Instances or Restricted Instances **cannot be repaired**

**Repair Speed-Up:**
- Spend **10 SP to reduce repair time by 10%**
- Better scaling for longer repair durations

**Damaged Ships:**
- Ships that survive combat return with their surviving numbers
- No "damage" state between operational and destroyed

### Escalation Risk

Attacking corporation members may trigger corp intervention:
- Corp members can retaliate
- May escalate into large-scale warfare involving allied corporations
- Significant strategic risk factor

**Sources:**
- [Attacking Neighbors (PvP) | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Attacking_Neighbors_(PvP))
- [Combat Mechanics | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Combat_Mechanics)
- [Spacedock | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Spacedock)
- [Stability | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Stability)
- [Beginner FAQ | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Beginner_FAQ)

---

## 5. Being Attacked / Defense

### Radar Building - Early Warning System

**Primary Function:**
- Provides information on attacking fleets
- Information found in "Check Transition List" window
- Higher radar level = more information revealed about incoming enemies

**Detection Capabilities by Level:**

| Radar Level | Detection Time | Information Revealed |
|-------------|---------------|---------------------|
| 1 | 30 minutes | Arrival Time |
| 2 | 1 hour | Arrival Time |
| 3 | 1.5 hours | Arrival Time + Origin Coordinates |
| 4 | 2 hours | Arrival Time + Origin Coordinates |
| 5 | 2.5 hours | Arrival Time + Origin Coordinates + Fleet Number Strength |
| 6 | 3 hours | Arrival Time + Origin Coordinates + Fleet Number Strength |
| 7 | 3.5 hours | Arrival Time + Origin Coordinates + Fleet Number Strength + Fleet Types |
| 8 | 4 hours | Arrival Time + Origin Coordinates + Fleet Number Strength + Fleet Types |
| 9 | 4.5 hours | Arrival Time + Origin Coordinates + Fleet Number Strength + Fleet Types + Commanders |

**Construction Requirements:**
- Level 1 requires: Civic Center Lv2, Technology Center Lv1
- Level 9 requires: Civic Center Lv10, Technology Center Lv9

**How to Check Incoming Fleets:**
1. Click Galaxy view icon
2. Click green radar symbol
3. Click on 5 tabs at top left
4. View incoming fleet information based on radar level

### Recalling/Dismissing Fleets

**[UNCONFIRMED]** - Wiki does not document fleet recall mechanics during incoming attack detection. Unknown if fleets can be recalled mid-travel when incoming attack is detected.

### Truce Card - Protection Mechanic

**Truce Card:**
- Prevents **all attacks for 12 hours** starting when activated
- **Cannot be used before or during an imminent attack**
- No attacks may be made BY a player protected by Truce Card (protection goes both ways)
- Obtained from: New Year Chest, Pirate Chest 2 (drops from Humaroid instances 13-14)

**Other Protection Items (functions unknown):**
- Healing Card (New Year Chest, Pirate Chest 1)
- Revival Card (New Year Chest, Pirate Chest 1)

**[UNCONFIRMED]** "Peace Shield" mechanics beyond Truce Card - no other protection items documented.

### Defense Fleets

**Fleet Uses:**
- "Fleets are used for protecting your planet, attacking other planets, instances, League, Arena, Space Raids, and for the Championship"

**Fleet Capacity:**
- Each fleet: 9 stacks of 3,000 ships = **27,000 ships maximum per fleet**
- **[UNCONFIRMED]** Maximum number of fleets per player - one example mentions 25 fleets, but no hard limit documented

**Fleet Positioning:**
- Fleets use a 3×3 grid with 9 positions called "stacks"
- Three ranks with varying attack power:
  - **First Rank (100% attack)**: Left/Right Shoulders and Head
  - **Second Rank (90% attack)**: Flanks and Glasshouse (central position)
  - **Third Rank (75% attack)**: Rear positions
- **Glasshouse** = "the single best protected position in any Phalanx"
- Typically houses glass cannons or flagships needing protection

**Defense Scaling:**
- Defense = Stack size × shields and struct values
- Fleet reaches full defensive readiness only at max stack size (3,000 ships each)
- Adjacent stacks can receive up to **30% of damage** through scattering effects

**Enemy Fleet Targeting:**
Six targeting criteria exist:
1. Maximum attack power
2. Minimum attack power
3. Maximum durability
4. Minimum durability
5. Closest distance
6. Highest commander rank

Players can manipulate fleet composition/positioning so fleets fulfill none of these criteria, making them unlikely targets when other fleets are present.

**Fleet Roles:**
- **Tanks:** Ships biased towards defense, usually with single weapons module consistent with fleet type
- **Glass Cannons:** High offense, low defense
- **Balanced:** Mix of offense and defense

**[UNCONFIRMED]** How to explicitly "set" defense fleets or garrison fleets - wiki mentions fleets protect your planet but doesn't detail stationing/garrison mechanics.

### What Happens When Offline and Attacked

**Defense Participation:**
- Orbital defenses (Meteor Stars, Particle Cannons, Anti-Aircraft Guns, Thor's Cannons) auto-fire when attacked
- Defense buildings follow their mechanics (cooldowns, ranges, targeting)
- Fleets stationed at planet participate in defense
- **[UNCONFIRMED]** Explicit "offline defense mode" or special mechanics when player is offline

**After Defeat:**
- All orbital defenses **reset automatically at no cost**
- Higher-level defenses take longer to rebuild (days for Lv10+ Thor cannons)
- Loot taken: 20% of resources not in warehouse
- Ships lost can be repaired via Spacedock (percentage recovery based on Spacedock level)

**Sources:**
- [Radar | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Radar)
- [Loot Chests | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Loot_Chests)
- [Fleet Design | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Fleet_Design)
- [Fleets | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Fleets)

---

## 6. Fleet Movement

### Wormholes & Return Travel

**Return Mechanics:**
- Victorious fleets "will return home on their own"
- Appear at "one of the 4 wormholes located at the corners of your space base"
- Return travel requires **half the outbound journey time**

**[UNCONFIRMED]** Whether wormholes are used for outbound travel or only returns.

### Speed Formula

**[UNCONFIRMED]** - Exact speed formula not documented in wiki.

**Known Factors:**
- Base ship speed (varies by design)
- Commander bonuses (likely affect speed)
- Technology bonuses (Ship Defense Science research tree exists)
- Distance between planets affects travel time

**Fleet Speed Types (from Fleet Strategies):**
- **Ballistic Fleet:** Movement should be **5+** to catch missile/fighter opponents
- **Missile Fleet:** Speed of **3-4+** necessary
- Ship design directly affects fleet speed

### Fleet Recall Mid-Travel

**[UNCONFIRMED]** - Wiki does not document ability to recall fleets mid-travel.

### Fleet Coordination

**[UNCONFIRMED]** - Specific mechanics for multiple fleets attacking same target simultaneously not documented.

**Known:**
- Players can have multiple fleets (at least 25+ mentioned)
- Fleets can be sent to attack other planets
- Multiple attacks may trigger corp/alliance responses

### He3 Fuel Consumption

**Critical Mechanic for Fleet Operations:**

**Fuel for Combat:**
- Firing weapons consumes He3
- He3 cost calculated by weapon blueprint cost
- Example: Rocket Frame-III has He3 cost of 0.07
  - 5 weapons × 0.07 = 0.35 He3 per fire
  - 100 ships = 35 He3 per volley
- Ballistic and missile weapons have higher He3 costs
- Risk: "You always run the risk of running out of power before you run out of targets"

**Fuel Depletion Consequence:**
- Ships with depleted He3 sit idle through remaining combat rounds
- Can cause battle loss even with surviving ships

**Mitigation:**
- Choose weapons with lower He3 consumption
- Research Ship Defense Science (top row) to reduce He3 expenditure
- Ensure adequate He3 supply before attacks

**Sources:**
- [Attacking Neighbors (PvP) | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Attacking_Neighbors_(PvP))
- [Fleet Strategies | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Fleet_Strategies)
- [He3 | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/He3)

---

## 7. Ship Combat Mechanics Deep Dive

### Stability System

**What is Stability?**
- Ship attribute that "makes it harder for ships to be destroyed"
- Determines threshold (as multiplier of hull strength) needed to destroy one ship
- Example: 400% stability + 1,000 structure per ship
  - Each ship requires 4,000 damage to destroy (1,000 × 400%)
  - 10,000 damage destroys only 2 ships (not 3)

**Stability Sources:**
- **Hull selection:** +40-60%
- **Modules:** Gravity Maintenance Facility, Ship Reinforcement Facility, Copernicus Gyroscope, Daedalus Control System, Icarus Control System
- **Research:** Ship Defense Science technology tree (+60%)
- **Commander enhancements:** Gems and bionic chips
- **Daedalus Control System special:** 30% chance per attack to grant additional +100% stability boost

### Ship Class Advantages

**Rock-Paper-Scissors Balance:**
- **Frigates > Battleships**
- **Battleships > Cruisers**
- **Cruisers > Frigates**

### Fleet Design Philosophies

**1. Ballistic Fleet (Close Combat) - "Victory"**
- Weapons: Short-range ballistics (1-2 range) + laser guns (2-5 range)
- Settings: Commander on 'next target' + 'close combat'
- Movement: **5+** (to catch missile/fighter opponents)
- Defense: Strong shields and armor essential
- Counters: Rocket and fighter fleets

**2. Missile Fleet (Long Range) - "Concordia"**
- Weapons: Rockets (5-8 range) + fighter carriers (6-10 range) - kept in separate fleets
- Settings: Commander on 'next target' + 'far/long range combat'
- Speed: **3-4+**
- Defense: Anti-rocket flak defenses, glass cannon builds
- Tactics: Hit-and-run

**3. Mixed Fleet (Balanced) - "Exelsior"**
- Weapons: Rockets (5-8 range) + laser guns (2-5 range)
- Range: Mid-range engagement
- Defense: Medium armor and shields, emphasis on shields (bonus attributes)
- Pairs well with ballistic "Victory" fleets

**Critical Design Rule:**
"NEVER MAKE A FLEET THAT HAS MORE THAN ONE TYPE OF WEAPON" - mixed weapon types reduce efficiency dramatically.

### Key Combat Elements

**Important Stats:**
- **Ship stability:** Prevents excessive casualties when shields fail
- **Agility:** Reduces incoming damage (-4% hit chance per 1 agility)
- **Steering:** Increases hit chance (+4% hit chance per 1 steering)
- **Balance:** Balance between shields and armor matters more than maximizing either

**Interception:**
- Defense modules shoot down incoming attacks
- Example: "One PPC fires and has a 55% chance to shoot down a single incoming attack"

**Sources:**
- [Fleet Strategies | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Fleet_Strategies)
- [Stability | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Stability)
- [Combat Mechanics | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Combat_Mechanics)

---

## 8. Resource Warehouse Protection

### Warehouse Function
- Resource Warehouse "will collect your resources until you harvest from it"
- Can upgrade to **level 24**
- Storage capacity: 10,000 of each resource (Lv1) → 20,000,000 (Lv24)

### PvP Loot Protection

**Critical Rule:**
- When you win a battle, you get "20% of the loser's resources, **not including what is in his warehouse**"
- Resources stored in the warehouse are **protected from looting**
- Only harvested resources (outside warehouse) can be stolen

**[UNCONFIRMED]** Exact mechanics of what counts as "in warehouse" vs "harvested resources" - warehouse may auto-collect or require manual harvesting.

### Capacity Scaling

**Research Bonus:**
- "Expand Capacity" research under Logistics Construction Science tree
- Affects storage values

**Sources:**
- [Resource Warehouse | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Resource_Warehouse)
- [Beginner FAQ | Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Beginner_FAQ)

---

## Summary of Key Findings

### ✅ Well-Documented Mechanics
1. **Radar system** - 9 levels with clear detection progression (30min to 4.5hr early warning)
2. **Defense buildings** - 4 types with detailed stats, ranges, cooldowns, and HP values
3. **PvP attack flow** - Scout → Attack → Combat → Loot → Return
4. **Loot mechanics** - 20% of resources not in warehouse
5. **Combat phases** - 8-phase combat system with hit chance, interception, damage, shields, hull damage
6. **Stability system** - Detailed mechanics for ship destruction probability
7. **Spacedock repair** - Probabilistic repair (1%-20% based on level)
8. **Truce Card** - 12-hour protection (cannot use during imminent attack)
9. **Wormhole returns** - 4 wormholes, 50% return travel time
10. **He3 fuel** - Weapons consume fuel, risk of running out during combat

### ⚠️ Partially Documented
1. **Galaxy coordinates** - Exist and are visible, but format/system not detailed
2. **Space Station** - Documented as main entrance, but Civic Center relationship unclear
3. **Defense fleets** - Fleets protect planet, but garrison/stationing mechanics not explicit
4. **Fleet speed** - Design affects speed (3-5+ mentioned), but no formula
5. **SP cost for attacks** - Used but exact cost not specified (instances cost 1 SP)

### ❌ Not Found / Unconfirmed
1. **Galaxy coordinate system format** (zones, sectors, systems)
2. **Search/scan mechanics** for finding players
3. **Travel time formula** (distance × speed calculation)
4. **Fleet recall mid-travel**
5. **Multiple fleet coordination** mechanics
6. **Exact SP cost per PvP attack**
7. **Space Station level requirements table**
8. **Space Station ↔ Civic Center "within 1 level" rule**
9. **Maximum number of fleets per player**
10. **Explicit "defense fleet" stationing/garrison commands**
11. **Offline defense mode** specifics
12. **Peace Shield** mechanics beyond Truce Card
13. **Loot caps** per resource type
14. **Warehouse auto-collect** vs manual harvest

---

## Research Sources

All information sourced from Galaxy Online II Wiki (Fandom):
- [Combat Mechanics](https://galaxyonlineii.fandom.com/wiki/Combat_Mechanics)
- [Composite Orbital Defenses Table](https://galaxyonlineii.fandom.com/wiki/Composite_Orbital_Defenses_Table)
- [Radar](https://galaxyonlineii.fandom.com/wiki/Radar)
- [Attacking Neighbors (PvP)](https://galaxyonlineii.fandom.com/wiki/Attacking_Neighbors_(PvP))
- [Beginner FAQ](https://galaxyonlineii.fandom.com/wiki/Beginner_FAQ)
- [Walkthrough](https://galaxyonlineii.fandom.com/wiki/Walkthrough)
- [Fleet Design](https://galaxyonlineii.fandom.com/wiki/Fleet_Design)
- [Fleet Strategies](https://galaxyonlineii.fandom.com/wiki/Fleet_Strategies)
- [Fleets](https://galaxyonlineii.fandom.com/wiki/Fleets)
- [Stability](https://galaxyonlineii.fandom.com/wiki/Stability)
- [Spacedock](https://galaxyonlineii.fandom.com/wiki/Spacedock)
- [Loot Chests](https://galaxyonlineii.fandom.com/wiki/Loot_Chests)
- [Resource Warehouse](https://galaxyonlineii.fandom.com/wiki/Resource_Warehouse)
- [Orbital Bases](https://galaxyonlineii.fandom.com/wiki/Orbital_Bases)
- [Guide To Advancing Quickly](https://galaxyonlineii.fandom.com/wiki/Guide_To_Advancing_Quickly)
- [He3](https://galaxyonlineii.fandom.com/wiki/He3)
- [Defense Strategies](https://pylonnexus.com/galaxyonlineiifandomcom/wiki/Defense_Strategies.html)

---

**Note:** Galaxy Online II is no longer in service. A remake called SplitWars: Online exists with new graphics and active developers.
