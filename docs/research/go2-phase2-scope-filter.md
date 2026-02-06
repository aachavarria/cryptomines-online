# Phase 2 Scope Filter: Frigate/Cruiser/Battleship Compatibility

**Date:** 2026-02-06
**Purpose:** Separate IN SCOPE (compatible with Frigate/Cruiser/Battleship) from OUT OF SCOPE (requires Special Hulls/Flagships) for Phase 2 implementation.

---

## 1. HULL TYPES

### IN SCOPE - Standard Hulls (3 classes)

All standard hulls follow the rock-paper-scissors balance (+/-5% damage modifiers).

**Frigates** (10 hull lines, each with tiers I-III):
- Weikes, Air Wanderer, Valkyrie, GoGetter, Space Hunter
- Sparrow, Devourer, Polymesus, Cybra, Hamdar
- Armor types: Nano, Neutralizing
- Base Effective Stack: 1100

**Cruisers** (10 hull lines, each with tiers I-III):
- Typhoon, Bombardier, Duke, The Shuttler, Watchman
- Spinner, Wraith, Encratos, Nicholas, Helena
- Armor types: Chrome, Regen
- Base Effective Stack: 1000

**Battleships** (5+ hull lines, each with tiers I-III):
- Estrella, Nettle, Diaz, RV766-The Explorer, Palenka
- Armor types: All four (Chrome, Regen, Nano, Neutralizing)
- Base Effective Stack: 900

### OUT OF SCOPE - Special Hulls

Special Hulls are obtained via Badge Points (from Restricted Instances) or special blueprints. They have unique bonuses, great volume capacity, best shields, and most He3 storage.

**Special Hull Frigates:** Industrial Ships (P1-P3), Fleetfoot, Hedgehog, Erotes series, Exodus series, etc.
**Special Hull Cruisers:** Daybreak, Last Stand, Shadow Guardian, Quick Assault, Arctic Airboard, etc.
**Special Hull Battleships:** Aggressive Warlord, Alliance Admiral, Presidio of Glory, Novas Ark, etc.

### OUT OF SCOPE - Flagships

**Federation Flagships:**
- Liberty Wings (Frigate) - from Leo Constellation
- Independence I/II/III (Cruiser) - from Instance 9 / Constellation
- Black Hole I/II/III (Battleship) - from Capricorn Constellation

**Humaroid Flagships** (obtained via Blueprint Shreds):
- Intrepid Nexus (100 each Shreds I, II, III)
- Grim Reaper (100 each Shreds I, II, III)
- Shadow Trojan (80 each Shreds I, II, III, IV)
- Firecat (80 each Shreds I, II, III, IV + 50 Shred V)
- Mercury Wing (250 Shred I + 60 Shred II)
- GForce's Dreadnaught (180, 180, 60 shreds)
- Conquistador (random only)
- Arbiter (random only)

**Why OUT OF SCOPE:** Humaroid Flagships have special rules (max 1 stack per fleet), require Blueprint Shreds from endgame content (Humaroids, Scenario Instance milestones), and provide fleet-wide abilities. Federation Flagships require Constellation Instance rewards. Both are advanced progression systems beyond Phase 2.

---

## 2. MODULES

### Key Finding: ALL modules are hull-agnostic

After reviewing every module category on the GO2 wiki, **no modules have hull-type restrictions**. All modules can be installed on any hull type (Frigate, Cruiser, Battleship, Special, or Flagship). The only constraints are:

1. **Volume/slot capacity** - Each hull has limited installation slots; higher-tier hulls have more
2. **Per-ship limits** - Some modules limited to 1 per ship (e.g., Extreme Counterattack, Nano Station Warehouse)
3. **Blueprint requirements** - Modules need blueprints before they can be used in designs

### IN SCOPE - All Module Categories

**Attack Modules - Ballistic Weapons** (18 weapon lines x 3 tiers = 54 modules):
- Rapid Fire, Taskmaster, Tempest Rapid Fire, Wiseman Scatterbomb
- Gatling Cannon, Blaze Auto-Cannon, Typhoon Speed Cannon, Tenho Dimensional Bomb
- Flash Sonicbomb, Calamity Howitzer, Deathhowl Shockwave, Ares Widowmaker
- Azrael, BloodyMary, Tempest, Judgment, Bloodspur, Shooting Star
- Range: 1-2, Cooldown: 0, Damage: Kinetic or Heat

**Attack Modules - Directional Weapons** (19 weapon lines x 3 tiers = 57 modules):
- Cluster Laser Transmitter, Magneto Pulsar, Pulse Laser Transmitter, Guided Laser Bomb
- The Sumo, Magneto Bomb, Plasma Energy Cannon, Positron Bomb
- Meltdown, Heartstopper, "Princes" Lightspear, "Tomahawk" Godfist
- Flamer, Nemesis, Cyclone, Avalanche, Phantom Ray, Scorpion Stinger, Tracking Bolt
- Range: 2-5, Cooldown: 1, Damage: Heat or Magnetic

**Attack Modules - Missile Weapons** (20 weapon lines x 3 tiers = 60 modules):
- Rocket Frame, Starlight Missile Pod, Rocket Booster, Razor Missile Pod
- Refined Rocket Frame, Poacher Firebomb, Boomerang Rocket Rack, Hellfire Missile Pod
- Destroyer Rocket Rack, Thunder Missile Pod, Destroyer Nuclear Rack, Doomsday Nuclear Launcher
- Dragon Slayer, Wipeout, Terminator, Zeus, Shadowflare, Red Comet, Apate, Spitfire
- Range: 5-8, Cooldown: 3, Damage: Explosive

**Attack Modules - Ship-Based Weapons (SBW)** (19 weapon lines x 3 tiers = 57 modules):
- Kinetic: Streamliner, Guardian, Pandora, Titan, Dusk Interceptor
- Magnetic: Golem, Leopard Streamliner, Warhammer Streamliner, Nebula, Dark Specter
- Heat: Hunter Streamliner, Dragonturtle, Hammerhead, Space Citadel, Chronus Wing
- Explosive: Ladybug, Nomad Streamliner, Hornet Swarm, Widowmaker
- Range: 6-10, Cooldown: 4, Damage: All four types

**Attack Modules - Planetary Weapons** (6 weapon lines x 3 tiers = 18 modules):
- Lander Module, Molotov Cocktail, Heavy Duty, Star Shooter, Smasher, Doomsday
- Range: 1-2, Cooldown: 1, Only damages defensive structures

**Defense Modules - Structure** (13 module lines, most with tiers I-III):
- Atomic Framework (no tiers)
- Ship Reinforcement Facility, Gravity Maint. Facility, Hull Maintenance Mechanic
- Quick Reaction Armor, Reflective Plating, Energy Armor
- Daedalus Control System, Armstrong Radion Armor, Copernicus Gyroscope
- "Icarus" Control System, Drexler Ratchet, Tensor Armor

**Defense Modules - Shields** (20 module lines, most with tiers I-III):
- Orbital Shield (no tiers)
- Energy Shield Booster, Particle Stun Shield, Heat Diffusion Shield
- Space-Time Magnetic Shield, Detonator Shield, Shield Regenerator
- EOS Phase Shift Shield
- Sagan series: Kinetics, Heat, Magnetic, Anti-Explosive
- Hertz Regenerator, Jack-O-Lantern
- Tyson series: Kinetic, Heat, Magnetic, Anti-Explosion
- Twilight Armor, Zeroth Accelerator

**Defense Modules - Air Defense** (6 module lines x 3 tiers = 18 modules):
- Anti-Aircraft Cannon, Missile Interception, Fighter Interception Cabin
- Powered Pulse Cannon, Extreme Counterattack (max 1 per ship)
- Aldrin Particle Cannon

**Auxiliary Modules - Electronic** (8 module lines x 3 tiers = 24 modules):
- Agility Booster, Infrared Scanner, ECM Booster, Auto Target System
- Time Dilation Module, Armstrong Core, Artemis, Infrared Seeker
- Limit: 1 of each type per ship

**Auxiliary Modules - Storage** (3 module lines):
- Station Warehouse (no limit per ship)
- Nano Station Warehouse (max 1 per ship)
- Drexler Nano Cells (max 1 per ship)

**Auxiliary Modules - Transmission** (6 module lines, progression chain):
- Super Transmission Engine -> Team Combat Engine -> Anti-matter Engine
- -> Eos Phase Shift Engine -> Wheeler Engine -> Vapor Engine

### Module Totals IN SCOPE
- Attack: ~246 modules (across 5 weapon classes)
- Defense: ~117 modules (structure + shields + air defense)
- Auxiliary: ~33 modules (electronic + storage + transmission)
- **Total: ~396 individual modules**

### OUT OF SCOPE - Module-Related Items
- None. All modules work with standard hulls. The limiting factor is hull volume, not hull type.

---

## 3. BLUEPRINTS

### IN SCOPE - Standard Ship Blueprints

**Frigate Blueprints:** Weikes, Air Wanderer, Valkyrie, GoGetter, Space Hunter, Sparrow, Devourer, Polymesus, Cybra, Hamdar (10 blueprints)

**Cruiser Blueprints:** Typhoon, Bombardier, Duke, The Shuttler, Watchman, Spinner, Wraith, Encratos, Nicholas, Helena (10 blueprints)

**Battleship Blueprints:** Estrella, Nettle, Diaz, RV766-The Explorer, Palenka (5+ blueprints)

**Obtained via:** Initial quests, Normal Instance drops (10% from Treasure Box), Auction House, Galactic Trafficker

### IN SCOPE - Module Blueprints

All module blueprints are in scope since all modules work with standard hulls. Obtained via Normal Instances, quests, Auction House.

Examples from wiki:
- Attack: Rapid Fire, Taskmaster, Cluster Laser Transmitter, Magneto Pulsar, Rocket Frame, Starlight Missile Pod, Streamliner, Golem, Lander Module
- Defense: Ship Reinforcement Facility, Energy Shield Booster, Anti-Aircraft Cannon
- Auxiliary: Super Transmission Engine

### IN SCOPE - Blueprint Research (Weapons Factory)

- Blueprints can be researched/upgraded to Level 3
- Only Level 1 blueprint needed to start each upgrade
- Applies to all module blueprints (in scope)

### OUT OF SCOPE - Blueprint Shreds

Blueprint Shreds are fragments for crafting Humaroid Flagships:
- 5 shred types (Shred I through V)
- Obtained from: Defeating Humaroids, Scenario Instance milestones
- Required: 300+ total shreds for a single flagship
- Random Blueprint function: 200 of each type (1000 total)

**Why OUT OF SCOPE:** Blueprint Shreds exist solely to craft Flagships, which are out of scope.

### OUT OF SCOPE - Special Hull Blueprints

- Obtained via Badge Points (from Restricted Instances)
- Not available from Normal Instances

---

## 4. NORMAL INSTANCES (PvE)

### IN SCOPE - All 30 Normal Instances

Normal Instances do NOT require Special Hulls or Flagships. They are the standard PvE progression content. All 30 can be completed with Frigate/Cruiser/Battleship fleets.

**Known instances:**
1. Ancestral Recall - 3 max fleets, 180 EXP
2. Deadzone - 4 max fleets, 500 EXP
3. Bravery - 4 max fleets, 1,000 EXP
4-29. [Names need investigation - see research gaps]
30. Triumphant Glory - 15 max fleets, 73,500 EXP

**Rewards:** Treasure Boxes containing resources (Gold, Metal, He3) and Blueprints (10% chance)
**Farming:** Instances can be run repeatedly, data does NOT reset
**Blueprint drops are level-dependent** - specific blueprints only from specific instances

### OUT OF SCOPE - Instance Types

**Restricted Instances:**
- Daily challenges (3 free/day + 1 with Passport)
- Reward: Badge Points (used for Special Hull blueprints), Constellation Passes
- 10 difficulty levels, 10-25 max fleets
- **Why OUT:** Badge Points feed Special Hull acquisition; Constellation Passes feed Constellation Instances

**Scenario Instances (Trial Instances):**
- 10 sequential story-based scenarios, reset daily
- Reward: Blueprint Shreds, cards, resources
- **Why OUT:** Blueprint Shreds are for Flagships (out of scope)

**Constellation Instances:**
- Accessed via Constellation Passes (from Restricted Instances)
- Reward: Federation Flagship blueprints
- **Why OUT:** Federation Flagships are out of scope

---

## 5. BUILDINGS

### IN SCOPE

**Ship Factory** (24 levels):
- Core building for ship production
- 5 production slots (slot 5 via "Sync Shipbuilding" research)
- 20 design capacity
- No dependency on Special Hulls or Flagships

**Spacedock** (12 levels):
- Repairs ships lost in PvP
- 1%-20% repair rate (Level 1-12)
- Cannot repair Instance losses
- No dependency on Special Hulls or Flagships

### OUT OF SCOPE

**Compound Center** (mentioned in Commander system):
- Used for merging Commander Cards AND crafting Humaroid Flagships from Blueprint Shreds
- The flagship-crafting portion is out of scope
- Commander Card merging is partially in scope (see Section 6)

---

## 6. FLEET & COMMANDER SYSTEM

### IN SCOPE

**Fleet System:**
- 3x3 grid formation (9 stacks, 3000 ships each, 27000 max)
- Position-based attack power (100%/90%/75% by rank)
- Single design per stack rule
- Fleet creation, fleet uses (planet defense, attacks, Normal Instances)

**Commander System (Partial):**
- Commander Cards with 4 attributes (Accuracy, Dodge, Speed, Electron)
- Expertise ratings (S/A/B/C/D/F) for weapon types AND hull types
- Effective Stack bonuses based on star rank
- Card tiers: Skill, Super, Legendary, Divine
- Specialties: Attack, Defense, Energy

### OUT OF SCOPE (Commander-Related)

- Commander Card acquisition via Lucky Wheel (gacha system)
- Commander Card acquisition via 100 Mall Points draws
- Compound Center merging for star rank increases (endgame optimization)
- Mall Points economy in general

**Note:** The Commander system is deeply intertwined with all content. For Phase 2, we include the core mechanics (assign commander to fleet, expertise affects performance, effective stack) but defer the acquisition/upgrade loop which depends on Mall Points and endgame systems.

---

## 7. COMBAT SYSTEM

### IN SCOPE

The entire combat system works identically for standard and special hulls:

- 8-phase combat sequence
- Effective Stack mechanics
- Weapon range/cooldown/damage type system
- Armor system (Chrome, Regen, Nano, Neutralizing)
- 4 damage types (Kinetic, Heat/Solar, Explosive, Magnetic)
- Shield mechanics
- Module placement order importance
- Hit chance, critical hits, stability

### OUT OF SCOPE

- Light Armor (exclusive to Flagships)
- Flagship fleet-wide abilities (e.g., Conquistador's scattering rate bonus)
- Humaroid Flagship 1-stack-per-fleet restriction
- Federation Flagship effective stack bonus of +100

---

## 8. SCOPE SUMMARY

### Phase 2 IN SCOPE Checklist

| Category | In Scope | Details |
|----------|----------|---------|
| Hull Types | 3 classes | Frigate (10 lines), Cruiser (10 lines), Battleship (5+ lines), each with I/II/III tiers |
| Modules | ALL | ~396 modules across 11 categories, no hull restrictions |
| Ship Blueprints | ~25+ | Standard hull blueprints for all 3 classes |
| Module Blueprints | ALL | All module blueprints (no hull restrictions) |
| Blueprint Research | YES | Weapons Factory upgrade system (Levels 1-3) |
| Normal Instances | All 30 | Standard PvE progression, no special hull requirements |
| Ship Factory | YES | 24 levels, 5 slots, 20 designs |
| Spacedock | YES | 12 levels, PvP ship repair |
| Fleet System | YES | 3x3 grid, 27000 max ships, position-based power |
| Commander System | PARTIAL | Core mechanics only (assignment, expertise, effective stack) |
| Combat System | YES | Full 8-phase system with standard armor types |

### Phase 2 OUT OF SCOPE Checklist

| Category | Out of Scope | Reason |
|----------|-------------|--------|
| Special Hulls | ALL | Require Badge Points from Restricted Instances |
| Federation Flagships | ALL | Require Constellation Instance rewards |
| Humaroid Flagships | ALL | Require Blueprint Shreds from endgame content |
| Blueprint Shreds | ALL | Only for Flagship crafting |
| Restricted Instances | ALL | Feed Special Hull acquisition loop |
| Scenario Instances | ALL | Feed Blueprint Shred acquisition |
| Constellation Instances | ALL | Feed Federation Flagship acquisition |
| Light Armor type | YES | Exclusive to Flagships |
| Commander acquisition loop | PARTIAL | Mall Points, Lucky Wheel, Compound Center merging |
| Flagship abilities | ALL | Fleet-wide bonuses, special restrictions |

---

## 9. IMPLEMENTATION NOTES

### Volume Matters, Not Hull Type
The key insight is that modules are NOT restricted by hull type. The constraint is **volume/installation slots**. Frigates have fewer slots than Battleships, so they can equip fewer modules. This is the core design constraint for ship design.

### Module Count is Large
With ~396 individual modules, Phase 2 needs a strategy for module progression:
- Early game: Basic modules (tier I of each category)
- Mid game: Advanced modules (tier II-III)
- Blueprint drops from instances gate access to higher-tier modules

### Instance-Blueprint Mapping is Critical
Each Normal Instance drops specific blueprints. This mapping drives player progression. The exact mapping is a [NEEDS INVESTIGATION] item from the research.

### Recommended Phase 2 Implementation Subset
Given the large module count, the GDD designer may want to start with a representative subset:
- 2-3 weapons per weapon class (1 per damage type)
- 1-2 per defense subcategory
- 1 per auxiliary subcategory
- Scale up in later phases
