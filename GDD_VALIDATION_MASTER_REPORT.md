# GDD Validation Master Report
## CryptoMines Online - Comprehensive Validation Against Galaxy Online 2

**Date:** 2026-02-10
**Team:** 5 Researchers + Lead Investigator
**Status:** COMPLETE
**Total Findings:** 150+ items identified

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Critical Bugs (Blocking)](#1-critical-bugs-blocking)
3. [Missing Content (Data)](#2-missing-content-data)
4. [Missing Mechanics (Systems)](#3-missing-mechanics-systems)
5. [Incorrect Data](#4-incorrect-data)
6. [Quest System Issues](#5-quest-system-issues)
7. [Out-of-Scope (Documented)](#6-out-of-scope-documented)
8. [Implementation Priorities](#7-implementation-priorities)
9. [Developer Tasks](#8-developer-tasks)
10. [Sources & References](#9-sources--references)

---

## Executive Summary

### Validation Coverage
- ✅ **Ship Hulls:** 74% complete (75/102 hulls)
- ✅ **Technologies:** 93% complete (91/98+ techs)
- ✅ **Ship Modules:** 27% strategic coverage (37/135+ lines)
- ✅ **Combat Mechanics:** 38% implemented (18/47 mechanics)
- ✅ **Buildings:** Validated against GO2 wiki
- ✅ **Quest/Inventory Systems:** Validated and issues identified

### Overall Status
**GOOD FOUNDATION** - Core systems implemented with strategic coverage of GO2 content. Critical bugs identified that need immediate fixing. Missing content is documented and prioritized.

### Critical Issues Count
- 🔴 **P0 (Critical Bugs):** 9 issues - BLOCK FUNCTIONALITY
- 🟠 **P1 (High Priority):** 15 issues - IMPACT GAMEPLAY
- 🟡 **P2 (Medium Priority):** 30+ issues - ENHANCE EXPERIENCE
- 🟢 **P3 (Low Priority):** 100+ items - FUTURE EXPANSION

---

## 1. CRITICAL BUGS (BLOCKING)

These bugs prevent core functionality and MUST be fixed immediately.

### 🔴 BUG #1: Inventory System Inaccessible
**Status:** CRITICAL - System implemented but not accessible
**Impact:** Quest tutorial broken, players cannot use items/blueprints

**Problem:**
- Backend inventory.go EXISTS with GET /api/inventory and POST /api/inventory/{id}/use
- Frontend InventoryPanel.tsx EXISTS with "Use" button functionality
- **BUT: No "Inventory" button in SideNav.tsx**
- Players cannot access inventory → Cannot complete quests that require using items

**Affected Quests:**
- main_06: Requires "use_blueprint" (super_transmission_engine)
- main_08: Requires "use_blueprint" (estrella)
- main_18: Requires "increase_bag_slot"
- main_19: Requires "use_resource_pack"

**Fix:** Add Inventory button to SideNav.tsx (5 min fix)

**Developer Task:** [TASK-001]

---

### 🔴 BUG #2: Quest Rewards Using Out-of-Scope Items
**Status:** CRITICAL - Rewards not being delivered
**Impact:** Players complete quests but receive no items

**Problem:**
Quests configured with items that are DOCUMENTED as out-of-scope:
- ❌ `loudspeaker` - Used in 6 quests (World Chat no longer needs it)
- ❌ `revival_card` - main_12 (Commanders immortal)
- ❌ `healing_card` - main_20 (Commanders immortal)
- ❌ `galaxy_transfer` - main_19 (Galaxy system out-of-scope)

**Code Behavior:**
ClaimQuest validates item_key exists in item_types (line 332-336):
```go
var exists bool
err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM item_types WHERE item_key = $1)`, itemKey).Scan(&exists)
if err != nil || !exists {
    log.Printf("Item key %s does not exist in item_types, skipping", itemKey)
    continue  // ← Silently skips reward
}
```

**Fix:** Replace out-of-scope items with valid in-scope items:
- loudspeaker → construction_card or primary_metal_pack
- revival_card → truce_card
- healing_card → extra_tax or metal_mining_boost
- galaxy_transfer → primary_metal_pack

**Developer Task:** [TASK-002]

---

### 🔴 BUG #3: Quests Missing Descriptions
**Status:** CRITICAL - Players don't know what to do
**Impact:** Zero guidance on quest objectives

**Problem:**
- quest_types table HAS `description` column
- Backend SELECT includes qt.description (line 24)
- **BUT: INSERT statements have NO description values**
- All quests have NULL descriptions

**Example:**
```sql
INSERT INTO quest_types (quest_key, category, display_name, ...) VALUES
('main_01_collecting_resources', 'main', 'Collecting Resources', ...);
-- Missing: description field entirely
```

**Expected (from GO2):**
```
"Your planet is rich in natural resources, and a good Commander must know
how to control their inventory space. Harvest your Resource Warehouse to
collect resources."
```

**Fix:** Add descriptions to all 28 main quests, 12 side quests, 6 daily quests

**Developer Task:** [TASK-003]

---

### 🔴 BUG #4: Resource Warehouse Limit Incorrect
**Status:** FIXED (but document for reference)
**Impact:** Players could build 4 warehouses instead of 1

**Problem:**
- Database had `max_count_per_planet = 4`
- GO2 wiki confirms **maximum 1 Resource Warehouse per planet**

**Fix Applied:**
```sql
-- supabase/migrations/20260206005232_phase1_mvp.sql line 762
-- BEFORE: max_count_per_planet = 4
-- AFTER:  max_count_per_planet = 1
```

**Status:** ✅ FIXED

---

### 🔴 BUG #5: Battleship Armor Types Incorrect
**Status:** NOT FIXED - 4 battleships have wrong armor
**Impact:** Incorrect damage calculations in combat

**Problem:**
4 battleships have wrong armor_type values:

| Hull | Current Armor | Correct Armor | Source |
|------|---------------|---------------|--------|
| Nettle I/II/III | regen | **nano** | GO2 Wiki |
| Diaz I/II/III | nano | **neutralizing** | GO2 Wiki |
| RV766-The Explorer I/II/III | neutralizing | **regen** | GO2 Wiki |
| Palenka I/II/III | chrome | **nano** | GO2 Wiki |

**Impact:** Armor effectiveness matrix will be wrong (nano reduces Heat+Explosive, neutralizing reduces Magnetic+Explosive, etc.)

**Fix:**
```sql
UPDATE hull_types SET armor_type = 'nano' WHERE name LIKE 'nettle_%';
UPDATE hull_types SET armor_type = 'neutralizing' WHERE name LIKE 'diaz_%';
UPDATE hull_types SET armor_type = 'regen' WHERE name LIKE 'rv766_%';
UPDATE hull_types SET armor_type = 'nano' WHERE name LIKE 'palenka_%';
```

**Developer Task:** [TASK-004]

---

### 🔴 BUG #6: No Progress Tracking Display
**Status:** NOT FIXED - Backend has data, frontend doesn't show it
**Impact:** Players don't know quest progress (5/10 enemies defeated, etc.)

**Problem:**
- Backend returns `progress_value` and `requirement_value` (line 26-27 quests.go)
- Frontend QuestPanel receives the data
- **BUT: No UI component displays progress**

**Expected Display:**
```
Quest: Ship Building
Build 10 ships (3/10 completed)
[=======---] 30%
```

**Fix:** Add progress bar/counter to QuestPanel.tsx

**Developer Task:** [TASK-005]

---

### 🔴 BUG #7: Radar Max Level Incorrect
**Status:** NOT FIXED - max_level set to 9 instead of 10
**Impact:** Players cannot unlock Radar Level 10 (final level)

**Problem:**
- Database has `max_level = 9` for Radar building
- GO2 wiki confirms **Radar goes up to Level 10**
- Level 10 requires: Civic Center Lv10, Tech Center Lv9

**Current DB (line 770):**
```sql
('radar', 'Radar', 'core', 'ground', 450, 400, 550, 60, 3.0300, 2.8700, 0, 1.0000, 9, 1, ...)
```

**Fix:**
```sql
UPDATE building_types SET max_level = 10 WHERE building_key = 'radar';
```

**Source:** https://galaxyonlineii.fandom.com/wiki/Radar

**Developer Task:** [TASK-017]

---

### 🔴 BUG #8: Ship Factory Base Costs Wrong
**Status:** NOT FIXED - Level 1 costs completely incorrect
**Impact:** Ship Factory is 3× more expensive than it should be

**Problem:**
Ship Factory Level 1 costs are dramatically wrong:

| Resource | Current (DB) | Correct (GO2) | Error |
|----------|-------------|---------------|--------|
| Metal | 600 | **206** | +291% |
| He3 | 450 | **163** | +276% |
| Gold | 500 | **189** | +265% |
| Time (sec) | 200 | **110** | +82% |

**Current DB (line 772):**
```sql
('ship_factory', 'Ship Factory', 'military', 'ground', 600, 450, 500, 200, ...)
```

**Fix:**
```sql
UPDATE building_types
SET base_cost_metal = 206,
    base_cost_he3 = 163,
    base_cost_gold = 189,
    base_time_seconds = 110
WHERE building_key = 'ship_factory';
```

**Source:** https://galaxyonlineii.fandom.com/wiki/Ship_Factory

**Developer Task:** [TASK-018]

---

### 🔴 BUG #9: Galaxy Transporter Wrong Concept
**Status:** NOT FIXED - Completely wrong building implementation
**Impact:** Galaxy Transporter has wrong function, costs, and progression

**Problem:**
Galaxy Transporter is **NOT** a resource transport building. It's the **Inter-Galactic League** (PvP) access building.

**Incorrect Implementation:**
- Max Level: 12 (WRONG - should be **1**)
- Costs: 350/300/450 (WRONG - should be **20000/20000/20000**)
- Build Time: 100 seconds (WRONG - should be **86400 seconds / 24 hours**)
- Function: Unknown (WRONG - should be "Access Inter-Galactic League, 10 free PvP matches/day")

**Current DB (line 768):**
```sql
('galaxy_transporter', 'Galaxy Transporter', 'core', 'ground', 350, 300, 450, 100, 3.0300, 2.8700, 0, 1.0000, 12, 1, ...)
```

**Fix:**
```sql
UPDATE building_types
SET base_cost_metal = 20000,
    base_cost_he3 = 20000,
    base_cost_gold = 20000,
    base_time_seconds = 86400,
    max_level = 1,
    cost_mult_metal = 1.0000,
    cost_mult_he3 = 1.0000,
    cost_mult_time = 1.0000,
    description = 'Allows access to Inter-Galactic League. 10 free matches per day.'
WHERE building_key = 'galaxy_transporter';
```

**Source:** https://galaxyonlineii.fandom.com/wiki/Galaxy_Transporter

**Developer Task:** [TASK-019]

---

## 2. MISSING CONTENT (DATA)

Content that exists in GO2 but is missing from our implementation.

### 2.1 Ship Hulls - 27 Hulls Missing

**Coverage:** 75/102 hulls (74%)

#### Frigates: ✅ 100% Complete
All 10 frigate lines with 3 tiers each (30 hulls total):
- Weikes, Air Wanderer, Valkyrie, GoGetter, Space Hunter
- Sparrow, Devourer, Polymesus, Cybra, Hamdar

#### Cruisers: ⚠️ 83% Complete (6 missing)
Missing 2 advanced cruiser lines:

| Hull Name | Armor | Tiers | Priority | Stats (GO2) |
|-----------|-------|-------|----------|-------------|
| Chimera Capra | Regen | I, II, III | HIGH | Shields 1,310-1,790 / Structure 7,300-10,000 |
| Ultra Gwyar | Regen | I, II, III | HIGH | Shields 2,300-2,800 / Structure 10,950-15,000 |

**Why Missing:** These are high-tier cruisers with significantly higher stats than standard lines.

#### Battleships: ❌ 42% Complete (21 missing)
Missing 7 battleship lines (21 hulls total):

| Hull Name | Armor | Tiers | Priority | Notes |
|-----------|-------|-------|----------|-------|
| Howler | Chrome | I, II, III | HIGH | Standard battleship line |
| Whirlpool | Regen | I, II, III | HIGH | Standard battleship line |
| Cerberus | Neutralizing | I, II, III | HIGH | Standard battleship line |
| Genesis | Regen | I, II, III | HIGH | Standard battleship line |
| Tiamat | Chrome | I, II, III | HIGH | Standard battleship line |
| Chimera Viper | Nano | I, II, III | MEDIUM | Advanced variant |
| Ultra Calas | Chrome | I, II, III | MEDIUM | Advanced variant |

**Developer Task:** [TASK-006]

---

### 2.2 Technologies - 7 Advanced Fighter Techs Missing

**Coverage:** 91/98+ techs (93%)

#### Missing Fighter Technologies (Ship-Based Science Tree)

| Tech | Levels | Effect | Cost (Lv1) | Priority |
|------|--------|--------|------------|----------|
| Long-ranged Strike | 1 | +3-15% attack at 6-10 slots distance | 98,955 Gold, 4:15:00 | HIGH |
| Fighter Interception Countermeasures | 1-5 | -1-5% intercept rate | 135,381-561,161 Gold | HIGH |
| Formation Optimization | 1-2 | +5-10% crit dmg, -5-10% weapon space | 287,437 Gold | HIGH |
| Swarm | 1-3 | 10-30% chance +8-25% attack at +15-35% He3 | 630,845 Gold | HIGH |
| Fortune | 1-3 | +10-30% double damage/crit in long-range | 885,076 Gold | HIGH |
| Heavy Gear Research | 1-2+ | Enhanced reload/shield dmg/range/intercept | 10M Gold, 111h | MEDIUM |

#### Missing Missile Technology

| Tech | Effect | Status |
|------|--------|--------|
| Rapid Loading | Reduces reload time by 1 round | Partial wiki data |

**Impact:** Fighter weapon endgame tree incomplete. Current implementation covers core fighter mechanics but lacks advanced optimization.

**Developer Task:** [TASK-007]

---

### 2.3 Ship Modules - Strategic Coverage

**Coverage:** 37 module lines / 111 variants (27% by count, STRATEGIC by design)

**Design Decision:** We implemented representative modules from EVERY category rather than all 250+ GO2 modules.

#### Coverage by Category

| Category | GO2 Modules | CMO Modules | Coverage | Status |
|----------|-------------|-------------|----------|--------|
| Ballistic Weapons | 18 lines (54 modules) | 3 lines | 17% | ✅ ACCEPTABLE |
| Directional Weapons | 19 lines (57 modules) | 2 lines | 11% | ✅ ACCEPTABLE |
| Missile Weapons | 20 lines (60 modules) | 2 lines | 10% | ✅ ACCEPTABLE |
| Ship-Based Weapons | 19 lines (57 modules) | 2 lines | 11% | ✅ ACCEPTABLE |
| Structure Defense | 13 lines | 6 lines | 46% | ✅ GOOD |
| Shield Defense | 20 lines | 8 lines | 40% | ✅ GOOD |
| Air Defense | 6 lines | 3 lines | 50% | ✅ GOOD |
| Electronic | 8 lines | 5 lines | 63% | ✅ EXCELLENT |
| Storage | 3 lines | 2 lines | 67% | ✅ EXCELLENT |
| Transmission | 6 lines | 3 lines | 50% | ✅ GOOD |

**Validation Result:** ✅ **APPROVED** - Strategic coverage across all categories with sufficient variety for gameplay.

**Note:** Low weapon coverage is INTENTIONAL simplification. We have heat/kinetic variants and tier progression without overwhelming players with 250+ module choices.

---

### 2.4 Buildings - Validation Complete

**Status:** ✅ All core buildings implemented correctly

#### Ground Base Buildings (Implemented)

**Resources:**
- ✅ Metal Collector (max level 24, max 8 per planet)
- ✅ He3 Extractor (max level 24, max 8 per planet)
- ✅ Residential Area (max level 24, max 8 per planet)
- ✅ Resource Warehouse (max level 24, **max 1 per planet** - FIXED)

**Core:**
- ✅ Civic Center (max level 12, max 1 per planet)
- ✅ Technology Center (max level 12, max 1 per planet)
- ✅ Weapon Research Center (max level 12, max 1 per planet)

**Military:**
- ✅ Command Center (max level 10, max 1 per planet)
- ✅ Ship Factory (max level 10, max 1 per planet)
- ✅ Recycling Plant (max level 10, max 1 per planet)

#### Space Base Buildings

**Main:**
- ✅ Space Station (max level 10)

**Defense:** (Module 7 - Out of current scope)
- 🔲 Meteor Star
- 🔲 Particle Cannon
- 🔲 Anti-aircraft Gun
- 🔲 Thor's Cannon

**City Services:** (Phase 3+)
- 🔲 Alliance Center
- 🔲 Trading Center
- 🔲 Galaxy Transporter

**Landscaping:** (Out of scope - decorative only)
- 🔲 Casino Resort, Beacon, Monument, Fountain, etc.

**Developer Task:** None - Building implementation is correct for Phase 2 scope.

---

## 3. MISSING MECHANICS (SYSTEMS)

Game mechanics and systems that are not fully implemented.

### 3.1 Combat Mechanics - 29 Missing Mechanics

**Coverage:** 18/47 mechanics (38%)

#### ❌ Critical Hit System (Priority: HIGH)
**Status:** NOT IMPLEMENTED
**Impact:** Commander Electron stat is useless

**GO2 Mechanic:**
```
Critical Hit Chance = 5% + (Commander Electron / 200)
Critical Hit Damage = Base Damage × 1.5
```

**Example:**
- Commander with 100 Electron → 5% + (100/200) = 55% crit chance
- Commander with 200 Electron → 5% + (200/200) = 105% crit (always crit)

**Current Implementation:** None - crits never happen

**Developer Task:** [TASK-008]

---

#### ❌ Successive Strikes (Priority: HIGH)
**Status:** NOT IMPLEMENTED
**Impact:** Commander Speed stat underutilized

**GO2 Mechanic:**
```
Successive Strike Chance = Commander Speed / 500
Effect: Attack twice in same round
```

**Example:**
- Commander with 250 Speed → 250/500 = 50% chance to attack twice

**Current Implementation:** None - ships attack once per round

**Developer Task:** [TASK-009]

---

#### ❌ Position-Based Attack Modifiers (Priority: HIGH)
**Status:** NOT IMPLEMENTED
**Impact:** Fleet formation has no tactical meaning

**GO2 Mechanic:**
```
First Rank (Front-L, Shoulders):  100% attack power
Second Rank (Flanks):              90% attack power
Third Rank (Rear/Tail):            75% attack power
```

**Current Implementation:** All positions attack at 100%

**Developer Task:** [TASK-010]

---

#### ❌ Formation Bonuses (Priority: HIGH)
**Status:** NOT IMPLEMENTED
**Impact:** Formations are cosmetic only

**GO2 Formations:**

| Formation | Active Positions | Attack Bonus | Defense Bonus |
|-----------|-----------------|--------------|---------------|
| Phalanx | All 9 | +10% | +5% |
| Diamond | 5 (center + adjacent) | +5% | +15% |
| Battle Line | 6 (first 2 ranks) | 0% | +20% |
| Skirmish | Sparse | +15% | -10% |

**Current Implementation:** Formations exist but provide no bonuses

**Developer Task:** [TASK-011]

---

#### ❌ Shield Penetration (Priority: MEDIUM)
**Status:** NOT IMPLEMENTED
**Impact:** Shield Penetration tech does nothing

**GO2 Mechanic:**
- Ballistic/Directional weapons can bypass shields
- Base 15% penetration (tech increases to 30-45%)
- Penetrated damage goes directly to hull

**Current Implementation:** All damage hits shields first

**Developer Task:** [TASK-012]

---

#### ❌ Scatter Damage (Priority: MEDIUM)
**Status:** NOT IMPLEMENTED
**Impact:** Missile scatter techs useless

**GO2 Mechanic:**
```
Scatter Damage: Cannot be absorbed by shields or mitigated
Missile Tech "Exaltation": If total structure > target, deal 54% scatter = 432% unpreventable bonus damage
```

**Current Implementation:** None

**Developer Task:** [TASK-013]

---

#### ❌ Interceptor Mechanics (Priority: LOW)
**Status:** NOT IMPLEMENTED
**Impact:** Air Defense modules useless

**GO2 Mechanic:**
- PPC (Powered Pulse Cannon): 55% chance to shoot down each incoming attack
- Defense stat affects interception rate

**Current Implementation:** None - Module 6+ feature

---

#### ❌ EOS Phase Shift (Priority: LOW)
**Status:** NOT IMPLEMENTED

**GO2 Mechanic:**
- EOS Phase Shift Shield: 30% chance to absorb double damage
- Only effective stack benefits

**Current Implementation:** None

---

#### ❌ Weapon Expertise Grades (Priority: MEDIUM)
**Status:** NOT IMPLEMENTED
**Impact:** Commander specialization missing

**GO2 Mechanic:**
```
Weapon Expertise: S/A/B/C/D grades
S Grade: +30% damage with that weapon type
D Grade: -30% damage with that weapon type
```

**Current Implementation:** None - all commanders equal with all weapons

**Developer Task:** [TASK-014]

---

#### ✅ Currently Implemented Combat Mechanics

**Phase 1-8 Framework:**
- ✅ Calculate effective stacks (Effective Stack formula: 300 + Star Rank × 50)
- ✅ Ship type advantage (Frigate > Cruiser > Battleship > Frigate, +/-5% dmg)
- ✅ Determine attack order by speed
- ✅ Calculate hit chance (simplified formula)
- ✅ Calculate damage (random within min-max range)
- ✅ Apply damage to shields
- ✅ Calculate casualties (overflow to hull)
- ✅ Calculate loot

**Combat Resolution Status:** Basic framework complete, advanced mechanics missing.

---

### 3.2 Quest System Mechanics

#### ❌ Progress Tracking Display (Priority: HIGH)
**Status:** Data exists, UI missing
**See:** [BUG #6](#🔴-bug-6-no-progress-tracking-display)

#### ❌ Quest Chaining Validation (Priority: MEDIUM)
**Status:** Unclear if prerequisite_quest_id validation works

**To Test:**
- Can player claim main_03 before main_01?
- Do side quest tiers enforce prerequisites?

**Developer Task:** [TASK-015]

---

### 3.3 Commander System

**Status:** Module 4 - Coming Soon (Phase 2)

#### Missing Features (Expected in Module 4)
- 🔲 Commander recruitment (gacha system)
- 🔲 Commander cards (Common/Skill/Super, 0-15 stars)
- 🔲 Commander attributes (Accuracy, Dodge, Speed, Electron)
- 🔲 Fleet assignment (1 commander per fleet)
- 🔲 Effective Stack bonus (300 + Star Rank × 50)
- 🔲 Critical hits (Electron stat)
- 🔲 Successive strikes (Speed stat)

**Note:** Already documented as "Coming Soon" in UI. Out of scope for current validation.

---

### 3.4 Galaxy/Map System

**Status:** Module 6+ - Coming Soon (Phase 3+)

#### Missing Features (Expected in Module 6+)
- 🔲 7x7+ galaxy grid
- 🔲 Resource Bonus Planets (RBPs)
- 🔲 Territory control
- 🔲 Galaxy navigation
- 🔲 Planet colonization

**Note:** Already documented as "Coming Soon" in UI. Out of scope for current validation.

---

## 4. INCORRECT DATA

Data that exists but has wrong values.

### 4.1 Battleship Armor Types
**See:** [BUG #5](#🔴-bug-5-battleship-armor-types-incorrect)

4 battleships (12 hulls) have wrong armor_type values.

---

### 4.2 Resource Warehouse Limit
**See:** [BUG #4](#🔴-bug-4-resource-warehouse-limit-incorrect)

✅ FIXED: Changed from 4 to 1 per planet.

---

### 4.3 Hull Stats Scaling

**Issue:** Our hull stats are 2-3× higher than GO2 wiki values.

**Example - Typhoon I (Cruiser):**
- **Our data:** Shield 505 / Structure 2,599 / Slots 120
- **GO2 Wiki:** Shield 202 / Structure 1,040 / Slots 140

**Possible Explanations:**
1. Intentional game balance adjustment
2. Different data source than wiki
3. Stat inflation for blockchain economy

**Status:** Needs investigation - may be intentional design decision.

**Developer Task:** [TASK-016]

---

## 5. QUEST SYSTEM ISSUES

Comprehensive breakdown of quest-related problems.

### 5.1 Quest Rewards with Out-of-Scope Items
**See:** [BUG #2](#🔴-bug-2-quest-rewards-using-out-of-scope-items)

### 5.2 Missing Quest Descriptions
**See:** [BUG #3](#🔴-bug-3-quests-missing-descriptions)

### 5.3 Quest Requirements Documentation

#### Implemented Requirement Types

| Type | Target | Description | Status |
|------|--------|-------------|--------|
| harvest_resources | resource_warehouse | Collect from warehouse | ✅ Works |
| build_building | building_type | Construct building | ✅ Works |
| upgrade_building | building_type | Upgrade to level X | ✅ Works |
| research_tech | tech_key | Complete tech research | ✅ Works |
| use_blueprint | blueprint_key | Use blueprint item | ❌ Blocked (no inventory access) |
| recruit_commander | - | Recruit a commander | 🔲 Not implemented (Module 4) |
| create_ship_design | - | Design a ship | ✅ Works |
| build_ships | - | Build X ships | ✅ Works |
| create_fleet | - | Create a fleet | ✅ Works |
| replenish_ammo | - | Replenish fleet ammo | ✅ Works |
| build_defense | - | Build defense structure | 🔲 Not implemented (Module 7) |
| send_message | channel | Send chat message | 🔲 Chat not fully implemented |
| increase_bag_slot | - | Expand inventory | 🔲 Not implemented |
| use_resource_pack | - | Use resource pack | ❌ Blocked (no inventory access) |
| grow_comsats | - | Grow comsats | 🔲 Not implemented |
| add_friend | - | Add a friend | 🔲 Social not implemented |
| reach_production | resource_type | Reach production rate | ✅ Works |
| reach_storage | - | Reach storage capacity | ✅ Works |
| login | - | Daily login | ✅ Works |
| use_speedup | target | Use speedup item | 🔲 Not implemented |
| donate_resources | - | Donate to alliance | 🔲 Alliance not implemented |
| complete_instance | instance_type | Complete instance | 🔲 Instances not implemented |

### 5.4 Deferred Quests (Phase 3)

These quests are marked `is_active = false` and require features not in Phase 2:

```sql
-- Deferred Main Quests (social/chat-dependent)
('main_02_loud_and_clear', ..., 'send_message', 'world_channel', ..., false),
('main_18_bigger_bags', ..., 'increase_bag_slot', ..., false),
('main_19_resource_pack', ..., 'use_resource_pack', ..., false),
('main_20_growing_resources', ..., 'grow_comsats', ..., false),
('main_21_adding_friends', ..., 'add_friend', ..., false),
('main_22_mail_system', ..., 'send_message', 'mail', ..., false),
```

**Status:** Correctly deferred. Main quest chain skips around these (sequential prerequisites bypass deferred quests).

---

## 6. OUT-OF-SCOPE (DOCUMENTED)

Content that is documented as intentionally excluded from Phase 2.

### 6.1 Ship Hulls - Special & Flagship

#### OUT: Special Hulls
- Special Hull Frigates (P1-P3, Fleetfoot, Hedgehog, Erotes, Exodus, etc.)
- Special Hull Cruisers (Daybreak, Last Stand, Shadow Guardian, etc.)
- Special Hull Battleships (Aggressive Warlord, Alliance Admiral, etc.)

**Reason:** Obtained via Badge Points from Restricted Instances (endgame content).

#### OUT: Flagships

**Federation Flagships:**
- Liberty Wings (Frigate) - Leo Constellation
- Independence I/II/III (Cruiser) - Instance 9 / Constellation
- Black Hole I/II/III (Battleship) - Capricorn Constellation

**Humaroid Flagships:**
- Intrepid Nexus, Grim Reaper, Shadow Trojan, Firecat, Mercury Wing, GForce's Dreadnaught, Conquistador, Arbiter

**Reason:** Require Constellation Instance rewards / Blueprint Shreds (endgame progression).

**Documentation:** `/docs/research/go2-phase2-scope-filter.md` (lines 31-56)

---

### 6.2 Items - Out-of-Scope

**Documented in:** `/docs/planning/module-9-inventory-system-plan.md` (line 34)

#### OUT: Commander Enhancement Items
- Gems (commander equipment)
- Chips (commander upgrades)
- Merge items (commander fusion)

**Reason:** Complex commander progression system (Phase 3+).

#### OUT: Battle Items (Commander-Related)
- Healing Card (commanders immortal in our version)
- Revival Card (commanders immortal in our version)

**Reason:** Commanders don't die in our simplified system.

#### OUT: Social Items
- Loudspeaker (World Chat no longer needs item to send messages)
- Galaxy Transfer (Galaxy system not implemented)

**Reason:** Simplified social mechanics for Phase 2.

#### OUT: Planet Items
- Planet Transformation Packs (cut from scope entirely)

**Reason:** Planet types not implemented.

---

### 6.3 Buildings - Out-of-Scope

#### OUT: Defense Buildings (Module 7)
- Meteor Star
- Particle Cannon
- Anti-aircraft Gun
- Thor's Cannon

**Status:** Deferred to Module 7 (Planetary Defense).

#### OUT: City Services (Phase 3+)
- Alliance Center (alliance system)
- Trading Center (auction house)
- Galaxy Transporter (galaxy navigation)

#### OUT: Landscaping (Decorative)
- Casino Resort, Beacon, Monument, Fountain, Library, Theater, Park, College, Hospital, Shopping Center, Statue

**Reason:** Purely cosmetic, no gameplay impact.

---

### 6.4 Modules - All In Scope

**Key Finding:** ALL modules in GO2 are hull-agnostic. No modules are restricted by hull type.

**Implication:** Our strategic module coverage (27% by count) is sufficient since we have representatives from every category with proper tier progression.

---

### 6.5 Game Mechanics - Out-of-Scope

#### OUT: Restricted Instances
- Require special mechanics (loss limits, special rewards)
- Badge Points system

#### OUT: Scenario Instances
- Multi-stage epic content
- Blueprint Shreds drops

#### OUT: Constellation System
- Zodiac constellation battles
- Divine Scrolls
- Flagship rewards

#### OUT: Humaroid System
- Humaroid enemies
- Blueprint Shred farming

#### OUT: Alliance/Corp Full Features
- Alliance wars
- Territory control
- Resource donations
- Alliance instances

#### OUT: Auction House / Trading
- Player-to-player trading
- Item auction system
- Blueprint trading

#### OUT: Mall/Shop
- Real-money purchases
- Mall Points currency
- Premium items

---

## 7. IMPLEMENTATION PRIORITIES

Recommended order for fixing issues.

### Priority 0: Critical (Blocking Gameplay)

**Timeline:** Immediate (1-2 days)

1. **[TASK-001]** Add Inventory button to SideNav - 30 minutes
2. **[TASK-002]** Fix quest rewards (replace out-of-scope items) - 2 hours
3. **[TASK-003]** Add quest descriptions (46 quests) - 4 hours
4. **[TASK-004]** Fix battleship armor types (4 hulls) - 30 minutes
5. **[TASK-005]** Add progress tracking UI to quests - 2 hours

**Total P0:** ~9 hours

---

### Priority 1: High (Core Gameplay Impact)

**Timeline:** Week 1 (3-5 days)

6. **[TASK-006]** Add missing hulls (27 hulls: 6 cruisers + 21 battleships) - 8 hours
7. **[TASK-007]** Add missing fighter technologies (7 techs) - 4 hours
8. **[TASK-008]** Implement critical hit system - 6 hours
9. **[TASK-009]** Implement successive strikes - 4 hours
10. **[TASK-010]** Implement position-based attack modifiers - 3 hours
11. **[TASK-011]** Implement formation bonuses - 4 hours
12. **[TASK-014]** Implement weapon expertise grades - 6 hours

**Total P1:** ~35 hours

---

### Priority 2: Medium (Enhanced Experience)

**Timeline:** Week 2-3 (5-10 days)

13. **[TASK-012]** Implement shield penetration - 5 hours
14. **[TASK-013]** Implement scatter damage - 4 hours
15. **[TASK-015]** Test/fix quest chaining validation - 2 hours
16. **[TASK-016]** Investigate hull stats scaling discrepancy - 3 hours
17. Add missing cruiser hulls stats and details - 4 hours
18. Add missing battleship hulls stats and details - 6 hours

**Total P2:** ~24 hours

---

### Priority 3: Low (Future Expansion)

**Timeline:** Post-launch / Phase 3

- Interceptor mechanics (air defense)
- EOS Phase Shift shields
- Additional weapon modules (expand beyond strategic coverage)
- Commander system (Module 4)
- Galaxy/Map system (Module 6+)
- Defense buildings (Module 7)
- Alliance features (Phase 3+)

---

## 8. DEVELOPER TASKS

Detailed task breakdowns for implementation.

### TASK-001: Add Inventory Button to SideNav
**Priority:** P0 - Critical
**Estimated Time:** 30 minutes
**Blocking:** Quest tutorial, item usage

**Files to Modify:**
- `frontend/src/components/layout/SideNav.tsx`

**Changes:**
1. Add to NAV_ITEMS array:
```typescript
{ id: 'inventory', icon: '🎒', label: 'Items', route: null, locked: false },
```

2. Add state:
```typescript
const [inventoryOpen, setInventoryOpen] = useState(false)
```

3. Add handler in handleClick:
```typescript
if (item.id === 'inventory') {
  setInventoryOpen(prev => !prev)
  setQuestOpen(false)
  setResearchOpen(false)
  setChatOpen(false)
  return
}
```

4. Add to isActive:
```typescript
if (item.id === 'inventory') return inventoryOpen
```

5. Import and render InventoryPanel:
```typescript
import InventoryPanel from '../panels/InventoryPanel.tsx'
// ...
{inventoryOpen && <InventoryPanel onClose={() => setInventoryOpen(false)} />}
```

**Testing:**
- Click Inventory button in SideNav
- Verify InventoryPanel opens
- Verify items display correctly
- Test "Use" button functionality

---

### TASK-002: Fix Quest Rewards (Replace Out-of-Scope Items)
**Priority:** P0 - Critical
**Estimated Time:** 2 hours
**Blocking:** Quest rewards not delivered

**Files to Modify:**
- `supabase/migrations/20260206030000_quests.sql`

**Changes:**

Replace out-of-scope items with valid in-scope alternatives:

```sql
-- main_01: loudspeaker → construction_card
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"construction_card","quantity":1}]'
WHERE quest_key = 'main_01_collecting_resources';

-- main_03: loudspeaker → primary_metal_pack
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"primary_metal_pack","quantity":1}]'
WHERE quest_key = 'main_03_tech_center';

-- main_09: loudspeaker → extra_tax
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"extra_tax","quantity":1}]'
WHERE quest_key = 'main_09_residential';

-- main_12: revival_card → truce_card
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"truce_card","quantity":1}]'
WHERE quest_key = 'main_12_recruit';

-- main_13: loudspeaker → metal_mining_boost
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"metal_mining_boost","quantity":1}]'
WHERE quest_key = 'main_13_design_ship';

-- main_15: loudspeaker → he3_mining_boost
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"he3_mining_boost","quantity":1}]'
WHERE quest_key = 'main_15_fleet';

-- main_19: galaxy_transfer → primary_metal_pack
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"primary_metal_pack","quantity":1}]'
WHERE quest_key = 'main_19_resource_pack';

-- main_20: healing_card → extra_tax
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"extra_tax","quantity":1}]'
WHERE quest_key = 'main_20_growing_resources';

-- main_22: loudspeaker → truce_card
UPDATE quest_types SET reward_item_json = '[{"type":"item","item_key":"truce_card","quantity":1}]'
WHERE quest_key = 'main_22_mail_system';
```

**Create Migration:**
```bash
# Create new migration file
supabase/migrations/20260211000000_fix_quest_rewards.sql
```

**Testing:**
- Complete main_01 and claim reward
- Verify construction_card appears in inventory
- Test all 9 affected quests

---

### TASK-003: Add Quest Descriptions
**Priority:** P0 - Critical
**Estimated Time:** 4 hours
**Blocking:** Player guidance

**Files to Modify:**
- `supabase/migrations/20260206030000_quests.sql`
- OR create new migration: `20260211000001_add_quest_descriptions.sql`

**Quest Descriptions to Add:**

Based on GO2 wiki and quest objectives, add descriptions for all quests.

**Example Format:**
```sql
UPDATE quest_types SET description = 'Your planet is rich in natural resources. Harvest your Resource Warehouse to collect stored resources and begin building your empire.'
WHERE quest_key = 'main_01_collecting_resources';

UPDATE quest_types SET description = 'Research is the key to galactic dominance. Build a Technology Center to unlock the research tree and begin developing new technologies.'
WHERE quest_key = 'main_03_tech_center';

UPDATE quest_types SET description = 'Start your first research project. Research Concurrent Construction technology to unlock an additional construction slot, allowing you to build multiple structures simultaneously.'
WHERE quest_key = 'main_04_research';

-- ... Continue for all 46 quests (28 main + 12 side + 6 daily)
```

**Full Quest Description List:**
*(Will be provided in separate document due to length)*

**Testing:**
- Open Quests panel
- Verify descriptions display for all quests
- Check formatting and line breaks

---

### TASK-004: Fix Battleship Armor Types
**Priority:** P0 - Critical
**Estimated Time:** 30 minutes
**Blocking:** Combat damage calculations

**Files to Modify:**
- Create migration: `supabase/migrations/20260211000002_fix_battleship_armor_types.sql`

**Changes:**
```sql
-- Fix Nettle: regen → nano
UPDATE hull_types SET armor_type = 'nano'
WHERE name LIKE 'nettle_%';

-- Fix Diaz: nano → neutralizing
UPDATE hull_types SET armor_type = 'neutralizing'
WHERE name LIKE 'diaz_%';

-- Fix RV766-The Explorer: neutralizing → regen
UPDATE hull_types SET armor_type = 'regen'
WHERE name LIKE 'rv766_%';

-- Fix Palenka: chrome → nano
UPDATE hull_types SET armor_type = 'nano'
WHERE name LIKE 'palenka_%';
```

**Verification Query:**
```sql
SELECT name, armor_type
FROM hull_types
WHERE name LIKE 'nettle_%'
   OR name LIKE 'diaz_%'
   OR name LIKE 'rv766_%'
   OR name LIKE 'palenka_%'
ORDER BY name;
```

**Expected Results:**
- nettle_i, nettle_ii, nettle_iii: nano
- diaz_i, diaz_ii, diaz_iii: neutralizing
- rv766_i, rv766_ii, rv766_iii: regen
- palenka_i, palenka_ii, palenka_iii: nano

**Testing:**
- Run migration
- Verify armor types with query above
- Test combat with these hulls
- Verify damage effectiveness matrix applies correctly

---

### TASK-005: Add Progress Tracking UI
**Priority:** P0 - Critical
**Estimated Time:** 2 hours
**Blocking:** Player feedback on quest progress

**Files to Modify:**
- `frontend/src/components/panels/QuestPanel.tsx`

**Changes:**

Add progress display to quest items:

```typescript
// Add progress bar component
const ProgressBar = ({ current, required }: { current: number, required: number }) => {
  const percentage = Math.min(100, (current / required) * 100)

  return (
    <div className="quest-progress">
      <div className="progress-text">
        {current}/{required} ({Math.round(percentage)}%)
      </div>
      <div className="progress-bar-bg">
        <div
          className="progress-bar-fill"
          style={{ width: `${percentage}%` }}
        />
      </div>
    </div>
  )
}

// Use in quest rendering
{quest.status === 'active' && quest.requirement_value > 0 && (
  <ProgressBar
    current={quest.progress_value}
    required={quest.requirement_value}
  />
)}
```

**CSS to Add:**
```css
.quest-progress {
  margin-top: 8px;
}

.progress-text {
  font-size: 12px;
  color: #888;
  margin-bottom: 4px;
}

.progress-bar-bg {
  height: 6px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
  overflow: hidden;
}

.progress-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, #4CAF50, #8BC34A);
  transition: width 0.3s ease;
}
```

**Testing:**
- Start quest that requires building 5 structures
- Build 1 structure, check progress shows "1/5 (20%)"
- Complete quest, verify progress shows "5/5 (100%)"
- Test with various quest types

---

### TASK-006: Add Missing Hulls (27 Total)
**Priority:** P1 - High
**Estimated Time:** 8 hours
**Impact:** Content completeness

**Files to Modify:**
- `supabase/migrations/20260206020000_phase2_ships.sql`
- OR create: `20260211000003_add_missing_hulls.sql`

**Required Data for Each Hull:**
- name, display_name, hull_class, tier
- armor_type, base_shield, base_structure
- installation_slots, cargo_capacity
- base_speed, base_agility, base_effective_stack

**Missing Cruiser Hulls (6 hulls):**

1. **Chimera Capra I/II/III** (Regen armor)
   - Stats from GO2 wiki (need to scale to our system)

2. **Ultra Gwyar I/II/III** (Regen armor)
   - Stats from GO2 wiki (need to scale to our system)

**Missing Battleship Hulls (21 hulls):**

3. **Howler I/II/III** (Chrome armor)
4. **Whirlpool I/II/III** (Regen armor)
5. **Cerberus I/II/III** (Neutralizing armor)
6. **Genesis I/II/III** (Regen armor)
7. **Tiamat I/II/III** (Chrome armor)
8. **Chimera Viper I/II/III** (Nano armor)
9. **Ultra Calas I/II/III** (Chrome armor)

**INSERT Template:**
```sql
INSERT INTO hull_types (name, display_name, hull_class, tier, armor_type,
    base_shield, base_structure, installation_slots, cargo_capacity,
    base_speed, base_agility, base_effective_stack) VALUES
('howler_i', 'Howler I', 'battleship', 1, 'chrome',
    [shield], [structure], [slots], [cargo], 0, 0, 900),
-- ... repeat for all tiers and hulls
```

**Data Source:**
- GO2 wiki: https://galaxyonlineii.fandom.com/wiki/Category:Battleships
- Scale stats using our existing battleship ratios
- Hulls Validation Report: `/SHIP_HULLS_VALIDATION_REPORT.md`

**Testing:**
- Verify all 27 hulls appear in blueprints list
- Check armor types are correct
- Test creating ship design with new hulls
- Verify stats display correctly

---

### TASK-007: Add Missing Fighter Technologies
**Priority:** P1 - High
**Estimated Time:** 4 hours
**Impact:** Tech tree completeness

**Files to Modify:**
- `supabase/migrations/20260206040000_research.sql`
- OR create: `20260211000004_add_fighter_techs.sql`

**Technologies to Add:**

All stats from Technology Validation Report.

```sql
INSERT INTO technology_types (tech_key, tree, display_name, max_level,
    prerequisite_techs, base_cost_metal, base_cost_he3, base_cost_gold,
    base_time_seconds, cost_multiplier, time_multiplier, effects_json) VALUES

-- 1. Long-ranged Strike
('long_ranged_strike', 'ship_based', 'Long-ranged Strike', 1,
    '["fuel_optimization_3", "thruster_optimization_3", "fighter_mastery_1"]',
    0, 0, 98955, 15300, 1.0, 1.0,
    '[{"type":"attack_bonus","condition":"distance_6_to_10","min_percent":3,"max_percent":15}]'),

-- 2. Fighter Interception Countermeasures (5 levels)
('fighter_intercept_counter', 'ship_based', 'Fighter Interception Countermeasures', 5,
    '["thruster_optimization_3", "navigation_5", "fighter_mastery_1"]',
    0, 0, 135381, 16200, 1.43, 1.34,
    '[{"type":"intercept_reduction","per_level":1,"unit":"percent"}]'),

-- 3. Formation Optimization (2 levels)
('formation_optimization', 'ship_based', 'Formation Optimization', 2,
    '["fighter_intercept_counter_5"]',
    0, 0, 287437, 21600, 1.16, 1.34,
    '[{"type":"crit_damage","per_level":5},{"type":"weapon_space_reduction","per_level":-5}]'),

-- 4. Swarm (3 levels)
('swarm', 'ship_based', 'Swarm', 3,
    '["fighter_weapons_theory_10", "formation_optimization_2"]',
    0, 0, 630845, 43200, 1.33, 1.34,
    '[{"type":"chance_attack_boost","chance_percent_range":[10,30],"attack_boost_range":[8,25],"he3_cost_increase_range":[15,35]}]'),

-- 5. Fortune (3 levels)
('fortune', 'ship_based', 'Fortune', 3,
    '["thruster_optimization_5", "long_ranged_strike_1", "swarm_3"]',
    0, 0, 885076, 57600, 1.33, 1.34,
    '[{"type":"long_range_bonus","double_damage_percent_range":[10,30],"crit_rate_percent_range":[10,30]}]'),

-- 6. Heavy Gear Research (2 levels)
('heavy_gear_research', 'ship_based', 'Heavy Gear Research', 2,
    '["ingenuity_1"]',
    0, 0, 10000000, 399952, 1.0, 1.0,
    '[{"type":"reload_enhancement"},{"type":"shield_damage_boost"},{"type":"attack_range_increase"},{"type":"interception_boost"}]');
```

**Prerequisite Setup:**
```sql
-- Link prerequisites (run after INSERT)
UPDATE technology_types SET prerequisite_tech_ids = ARRAY[
  (SELECT id FROM technology_types WHERE tech_key = 'fuel_optimization'),
  (SELECT id FROM technology_types WHERE tech_key = 'thruster_optimization'),
  (SELECT id FROM technology_types WHERE tech_key = 'fighter_mastery')
] WHERE tech_key = 'long_ranged_strike';

-- ... repeat for all techs
```

**Testing:**
- Verify techs appear in research panel under Ship-Based Science
- Check prerequisites are enforced
- Test research completion
- Verify effects apply correctly (if combat mechanics implemented)

---

### TASK-008: Implement Critical Hit System
**Priority:** P1 - High
**Estimated Time:** 6 hours
**Impact:** Combat depth, Commander Electron stat utility

**Files to Modify:**
- `backend/internal/combat/combat_engine.go`

**Implementation:**

Add critical hit calculation in damage phase:

```go
// In calculateDamage function (around line 410)
func (e *Engine) calculateDamage(attacker *FleetStack, defender *FleetStack) int {
    baseDamage := randomInt(attacker.MinDamage, attacker.MaxDamage)

    // CRITICAL HIT CHECK
    critChance := 0.05 // Base 5%
    if attacker.Commander != nil {
        critChance += float64(attacker.Commander.Electron) / 200.0
    }

    isCrit := rand.Float64() < critChance
    if isCrit {
        baseDamage = int(float64(baseDamage) * 1.5) // +50% damage on crit
        e.log("CRITICAL HIT! %s deals %d damage (base: %d)",
            attacker.Name, baseDamage, baseDamage/1.5)
    }

    return baseDamage
}
```

**Add Commander to FleetStack:**
```go
type FleetStack struct {
    // ... existing fields
    Commander *Commander `json:"commander,omitempty"`
}

type Commander struct {
    Accuracy int
    Dodge    int
    Speed    int
    Electron int
}
```

**Testing:**
- Create commander with 100 Electron (55% crit chance)
- Run 100 attacks, verify ~55% are crits
- Verify crit damage is 1.5× base damage
- Test with 0 Electron (5% base crit)

---

### TASK-009: Implement Successive Strikes
**Priority:** P1 - High
**Estimated Time:** 4 hours
**Impact:** Combat depth, Commander Speed stat utility

**Files to Modify:**
- `backend/internal/combat/combat_engine.go`

**Implementation:**

```go
// In attack phase (around line 385)
func (e *Engine) executeAttacks(attacker *FleetStack, defender *FleetStack) {
    // Normal attack
    damage1 := e.calculateDamage(attacker, defender)
    e.applyDamage(defender, damage1)

    // SUCCESSIVE STRIKE CHECK
    successiveChance := 0.0
    if attacker.Commander != nil {
        successiveChance = float64(attacker.Commander.Speed) / 500.0
    }

    if rand.Float64() < successiveChance {
        // Second attack!
        damage2 := e.calculateDamage(attacker, defender)
        e.applyDamage(defender, damage2)
        e.log("SUCCESSIVE STRIKE! %s attacks again for %d damage",
            attacker.Name, damage2)
    }
}
```

**Testing:**
- Commander with 250 Speed (50% chance)
- Run 100 attacks, verify ~50% trigger second attack
- Verify both attacks deal full damage
- Test with 0 Speed (no successive strikes)

---

### TASK-010: Implement Position-Based Attack Modifiers
**Priority:** P1 - High
**Estimated Time:** 3 hours
**Impact:** Fleet formation tactics

**Files to Modify:**
- `backend/internal/combat/combat_engine.go`

**Implementation:**

```go
// Add position attack modifier
func getPositionAttackModifier(position string) float64 {
    switch position {
    case "front_left", "front_center", "front_right": // First Rank
        return 1.00 // 100% attack
    case "mid_left", "mid_center", "mid_right": // Second Rank
        return 0.90 // 90% attack
    case "back_left", "back_center", "back_right": // Third Rank
        return 0.75 // 75% attack
    default:
        return 1.00
    }
}

// In calculateDamage
func (e *Engine) calculateDamage(attacker *FleetStack) int {
    baseDamage := randomInt(attacker.MinDamage, attacker.MaxDamage)

    // Apply position modifier
    positionMod := getPositionAttackModifier(attacker.Position)
    finalDamage := int(float64(baseDamage) * positionMod)

    return finalDamage
}
```

**Add Position to FleetStack:**
```go
type FleetStack struct {
    // ... existing fields
    Position string `json:"position"` // "front_left", "mid_center", etc.
}
```

**Testing:**
- Place stack in front_left: verify 100% damage
- Place stack in mid_center: verify 90% damage (glasshouse)
- Place stack in back_center: verify 75% damage (tail)

---

### TASK-011: Implement Formation Bonuses
**Priority:** P1 - High
**Estimated Time:** 4 hours
**Impact:** Formation tactical choices

**Files to Modify:**
- `backend/internal/combat/combat_engine.go`
- `backend/internal/models/fleet.go`

**Implementation:**

```go
// Formation bonus definitions
type FormationBonus struct {
    AttackPercent  int
    DefensePercent int
}

func getFormationBonus(formation string) FormationBonus {
    switch formation {
    case "phalanx":
        return FormationBonus{AttackPercent: 10, DefensePercent: 5}
    case "diamond":
        return FormationBonus{AttackPercent: 5, DefensePercent: 15}
    case "battle_line":
        return FormationBonus{AttackPercent: 0, DefensePercent: 20}
    case "skirmish":
        return FormationBonus{AttackPercent: 15, DefensePercent: -10}
    case "tee_forward":
        return FormationBonus{AttackPercent: 5, DefensePercent: 10}
    case "enfilade":
        return FormationBonus{AttackPercent: 12, DefensePercent: 3}
    case "tee_reverse":
        return FormationBonus{AttackPercent: 3, DefensePercent: 12}
    default:
        return FormationBonus{AttackPercent: 0, DefensePercent: 0}
    }
}

// Apply in combat
func (e *Engine) calculateDamage(attacker *FleetStack) int {
    baseDamage := randomInt(attacker.MinDamage, attacker.MaxDamage)

    // Apply formation attack bonus
    formationBonus := getFormationBonus(attacker.Fleet.Formation)
    attackMod := 1.0 + (float64(formationBonus.AttackPercent) / 100.0)

    finalDamage := int(float64(baseDamage) * attackMod)
    return finalDamage
}

func (e *Engine) calculateDefense(defender *FleetStack) int {
    baseDefense := defender.BaseDefense

    // Apply formation defense bonus
    formationBonus := getFormationBonus(defender.Fleet.Formation)
    defenseMod := 1.0 + (float64(formationBonus.DefensePercent) / 100.0)

    finalDefense := int(float64(baseDefense) * defenseMod)
    return finalDefense
}
```

**Testing:**
- Phalanx: verify +10% attack, +5% defense
- Diamond: verify +5% attack, +15% defense
- Skirmish: verify +15% attack, -10% defense (negative defense)

---

### TASK-012 through TASK-016

*(Continued in separate implementation guide due to length)*

---

### TASK-017: Fix Radar Max Level
**Priority:** P0 - Critical
**Estimated Time:** 5 minutes
**Blocking:** Players cannot unlock Radar Level 10

**Files to Modify:**
- `supabase/migrations/20260206005232_phase1_mvp.sql` (line 770)
- OR create new migration: `20260211000002_fix_building_stats.sql`

**Changes:**

```sql
-- Fix Radar max_level: 9 → 10
UPDATE building_types
SET max_level = 10
WHERE building_key = 'radar';
```

**Source:** https://galaxyonlineii.fandom.com/wiki/Radar

**Testing:**
- Build Radar to Level 9
- Verify Level 10 is available to upgrade
- Verify requirements: Civic Center Lv10, Tech Center Lv9
- Verify upgrade works correctly

---

### TASK-018: Fix Ship Factory Base Costs
**Priority:** P0 - Critical
**Estimated Time:** 10 minutes
**Blocking:** Ship Factory is 3× more expensive than intended

**Files to Modify:**
- `supabase/migrations/20260206005232_phase1_mvp.sql` (line 772)
- OR create new migration: `20260211000002_fix_building_stats.sql`

**Changes:**

```sql
-- Fix Ship Factory base costs to match GO2 Level 1
UPDATE building_types
SET base_cost_metal = 206,
    base_cost_he3 = 163,
    base_cost_gold = 189,
    base_time_seconds = 110
WHERE building_key = 'ship_factory';
```

**Verification:**

| Resource | Before | After | GO2 Wiki |
|----------|--------|-------|----------|
| Metal | 600 | **206** | ✅ 206 |
| He3 | 450 | **163** | ✅ 163 |
| Gold | 500 | **189** | ✅ 189 |
| Time (sec) | 200 | **110** | ✅ 110 |

**Source:** https://galaxyonlineii.fandom.com/wiki/Ship_Factory

**Testing:**
- Start building Ship Factory Level 1
- Verify costs: 206 Metal, 163 He3, 189 Gold
- Verify build time: 110 seconds (~2 minutes)
- Verify building completes successfully

---

### TASK-019: Fix Galaxy Transporter Implementation
**Priority:** P0 - Critical
**Estimated Time:** 15 minutes
**Blocking:** Galaxy Transporter has completely wrong function and costs

**Files to Modify:**
- `supabase/migrations/20260206005232_phase1_mvp.sql` (line 768)
- OR create new migration: `20260211000002_fix_building_stats.sql`

**Changes:**

```sql
-- Galaxy Transporter is Inter-Galactic League (PvP) access building
-- NOT a resource transport building
-- Single level only, expensive, 24 hour build

UPDATE building_types
SET base_cost_metal = 20000,
    base_cost_he3 = 20000,
    base_cost_gold = 20000,
    base_time_seconds = 86400,  -- 24 hours
    max_level = 1,
    cost_mult_metal = 1.0000,   -- No progression (single level)
    cost_mult_he3 = 1.0000,
    cost_mult_time = 1.0000,
    description = 'Allows access to Inter-Galactic League. 10 free matches per day.'
WHERE building_key = 'galaxy_transporter';
```

**Verification:**

| Field | Before | After | GO2 Wiki |
|-------|--------|-------|----------|
| Metal | 350 | **20000** | ✅ 20000 |
| He3 | 300 | **20000** | ✅ 20000 |
| Gold | 450 | **20000** | ✅ 20000 |
| Time | 100s | **86400s (24h)** | ✅ 24 hours |
| Max Level | 12 | **1** | ✅ 1 |
| Function | Unknown | **PvP League Access** | ✅ |

**Source:** https://galaxyonlineii.fandom.com/wiki/Galaxy_Transporter

**Testing:**
- Start building Galaxy Transporter
- Verify costs: 20000/20000/20000
- Verify build time: 24 hours
- Verify only 1 level available (no upgrade button after completion)
- Verify description mentions Inter-Galactic League

**Note:** The actual PvP League functionality may be out-of-scope for MVP. This task only fixes the building stats to match GO2. The building can exist without implementing the full PvP system.

---

## 9. SOURCES & REFERENCES

### Research Team
- **hulls-researcher** - Ship Hulls Validation
- **technology-researcher** - Technology Trees Validation
- **modules-researcher** - Ship Modules Validation
- **combat-researcher** - Combat Mechanics Validation
- **buildings-researcher** - Buildings Validation
- **team-lead** - Quest/Inventory Investigation & Report Consolidation

### Source Documents

#### Internal Documentation
- `/SHIP_HULLS_VALIDATION_REPORT.md` - Comprehensive hull analysis
- `/BUILDINGS_VALIDATION_REPORT.md` - Buildings validation and bug findings
- `/docs/research/technology-validation-report.md` - Tech tree analysis
- `/docs/research/ship-modules-validation-report.md` - Module catalog
- `/docs/research/combat-fleet-mechanics-validation.md` - Combat mechanics
- `/docs/research/go2-phase2-scope-filter.md` - Out-of-scope definitions
- `/docs/planning/module-9-inventory-system-plan.md` - Inventory system design
- `/findings.md` - Original issue log
- `/QA_REPORT.md` - QA assessment

#### Galaxy Online 2 Wiki Sources

**Buildings:**
- [Resource Warehouse](https://galaxyonlineii.fandom.com/wiki/Resource_Warehouse)
- [Radar](https://galaxyonlineii.fandom.com/wiki/Radar)
- [Ship Factory](https://galaxyonlineii.fandom.com/wiki/Ship_Factory)
- [Galaxy Transporter](https://galaxyonlineii.fandom.com/wiki/Galaxy_Transporter)
- [Alliance Center](https://galaxyonlineii.fandom.com/wiki/Alliance_Center)
- [Compound Center](https://galaxyonlineii.fandom.com/wiki/Compound_Center)
- [Category: Buildings](https://galaxyonlineii.fandom.com/wiki/Category:Buildings)

**Ships:**
- [Category: Frigates](https://galaxyonlineii.fandom.com/wiki/Category:Frigates)
- [Category: Cruisers](https://galaxyonlineii.fandom.com/wiki/Category:Cruisers)
- [Category: Battleships](https://galaxyonlineii.fandom.com/wiki/Category:Battleships)
- [Composite Ship Table](https://galaxyonlineii.fandom.com/wiki/Composite_Ship_Table)
- [Development Quests](https://galaxyonlineii.fandom.com/wiki/Development_Quests)
- [Category: Items](https://galaxyonlineii.fandom.com/wiki/Category:Items)
- [Logistics Construction Science](https://galaxyonlineii.fandom.com/wiki/Logistics_Construction_Science)
- [Ballistics Science](https://galaxyonlineii.fandom.com/wiki/Ballistics_Science)
- [Directional Science](https://galaxyonlineii.fandom.com/wiki/Directional_Science)
- [Missile Science](https://galaxyonlineii.fandom.com/wiki/Missile_Science)
- [Ship-based Science](https://galaxyonlineii.fandom.com/wiki/Ship-based_Science)
- [Ship Defense Science](https://galaxyonlineii.fandom.com/wiki/Ship_Defense_Science)
- [Planetary Defense Science](https://galaxyonlineii.fandom.com/wiki/Planetary_Defense_Science)

---

## Appendix A: Quick Reference Statistics

### Completeness Summary

| Category | Implemented | Missing | % Complete | Status |
|----------|-------------|---------|------------|--------|
| **Frigates** | 30/30 | 0 | 100% | ✅ Complete |
| **Cruisers** | 30/36 | 6 | 83% | ⚠️ Good |
| **Battleships** | 15/36 | 21 | 42% | ❌ Incomplete |
| **Technologies** | 91/98+ | 7+ | 93% | ✅ Excellent |
| **Modules** | 37/135+ lines | 98+ lines | 27% | ✅ Strategic |
| **Combat Mechanics** | 18/47 | 29 | 38% | ⚠️ Basic |
| **Buildings** | 22 core (3 bugs) | 12 landscaping | 65% | ⚠️ Bugs Found |

### Bug Priority Distribution

| Priority | Count | Total Hours |
|----------|-------|-------------|
| P0 (Critical) | 9 | ~9.5 hours |
| P1 (High) | 12 | ~35 hours |
| P2 (Medium) | 6 | ~24 hours |
| P3 (Low) | Many | TBD |

### Content Gaps

- **Hulls:** 27 missing (6 cruisers, 21 battleships)
- **Technologies:** 7 missing (all fighter tree)
- **Modules:** Strategic coverage sufficient (no gaps)
- **Mechanics:** 29 missing (combat enhancements)

---

**END OF REPORT**

*Last Updated: 2026-02-10*
*Version: 1.0*
*Status: Ready for Developer Implementation*
