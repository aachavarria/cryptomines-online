# Game Mechanics Research Findings

**Researcher:** GDD Researcher
**Date:** 2026-02-10
**Task:** Investigate and document game mechanics from GDD and Galaxy Online
**Status:** Complete

---

## 1. Warehouse Limit

### Finding: Maximum 4 Resource Warehouses Per Player

**Evidence:**
- Database seed data (`/supabase/migrations/20260206005232_phase1_mvp.sql`, line 762):
  ```sql
  ('resource_warehouse', 'Resource Warehouse', 'resource', 'ground',
   380, 370, 480, 35, 1.7500, 1.7500, 0, 1.0000, 24, 4,
   'civic_center', 1, 'Stores all resources; increases capacity')
  ```
  - `max_count_per_planet = 4`

- GDD documentation (`/home/yurei/cryptomines-online/docs/gdd/game-design-document.md`, line 111):
  - Shows "Multiple" in the table but doesn't specify exact count
  - The actual implementation in database is 4

- Galaxy Online 2 source (`/home/yurei/cryptomines-online/docs/research/galaxy-online-2-mechanics.md`, line 140):
  - States "Multiple" without specific number

**Related Code:**
- `/home/yurei/cryptomines-online/docs/planning/module-3-resource-auto-production-plan.md`, line 625:
  - "Multiple warehouses → uses highest storage_capacity (existing logic)"
  - This indicates the system takes the maximum storage capacity when multiple warehouses exist

**Conclusion:**
- **Maximum warehouses: 4 per planet**
- Storage capacity stacks (or uses highest - needs clarification)
- Each warehouse can be upgraded to level 24

---

## 2. Technology Center Requirements

### Finding: Technology Center Required for Research

**Building Details:**
- **Max Level:** 12
- **Purpose:** Research facility for 7 science trees
- **Bonus:** 3% research time reduction per level (max 36% at Level 12)

**Evidence:**
- GDD (`/home/yurei/cryptomines-online/docs/gdd/game-design-document.md`, line 122):
  ```
  Technology Center | Research facility (7 science trees) | 12
  ```

- Research time reduction (`line 687`):
  ```
  Technology Center level reduces research time by 3% per level (max 36% at Lv 12)
  ```

- Database seed data shows Technology Center levels 1-12 with research time reduction

**Research Rules (from GDD line 684-688):**
- Only one research can be active per tree at a time (7 trees = up to 7 concurrent researches)
- All tech research costs Gold only (no Metal or He3)
- Technology Center level reduces research time by 3% per level (max 36% at Lv 12)
- Research costs increase exponentially per tech level (~1.53x cost per level, ~2.34x time per level)
- Technologies have prerequisite chains within their tree

**Implementation Status:**
- Task #6 pending: "Add Technology Center requirement validation for research"
- System needs to validate Technology Center level before allowing research

**Formula:**
```
ResearchTimeReduction = TechnologyCenterLevel * 0.03
Max reduction: 36% at Level 12
```

---

## 3. Weapon Center (Weapon Research Center) Requirements

### Finding: Weapon Research Center for Blueprint Research

**Building Details:**
- **Official Name:** Weapon Research Center
- **Max Level:** 12
- **Grid Size:** 2x2
- **Purpose:** Blueprint research for ship hulls and modules
- **Bonus:** 3% blueprint research time reduction per level (max 36% at Level 12)

**Evidence:**
- GDD (`/home/yurei/cryptomines-online/docs/gdd/game-design-document.md`, line 693):
  ```
  Weapon Research Center
  Separate building from the Technology Center.
  Handles blueprint research for ship hulls and modules.
  ```

- Time reduction formula (line 703):
  ```
  Time reduction formula: identical to Technology Center (3% per level)
  ```

- Levels table (lines 698-702):
  ```
  | Level | Civic Req | Build Time | Metal | He3 | Gold | Research Time Reduction |
  | 1     | 1         | 0:00:40    | 500   | 300 | 450  | 3%                     |
  | 6     | 6         | 2:53:05    | 92,371| 55,422| 83,134| 18%                  |
  | 12    | 12        | 2535:16:53 | 54,372,693| 32,623,616| 48,935,424| 36%     |
  ```

**Blueprint Research Levels (from GDD):**
- Research levels 1-10 for blueprints
- Each level improves blueprint quality/stats
- Separate from Technology Center research

**Implementation Status:**
- Task #7 pending: "Add Weapon Center requirement validation for military research"
- System needs to validate Weapon Research Center level before allowing blueprint research

**Distinction:**
- **Technology Center:** Science/tech research (7 trees)
- **Weapon Research Center:** Blueprint research (ship hulls and modules)

---

## 4. Fleet Formation Mechanics and Slot System

### Finding: 3x3 Grid Formation with 9 Stacks Maximum

**Fleet Composition (from GO2 mechanics, line 414-416):**
- **Fleet Grid:** 3x3 grid = 9 stacks maximum
- **Stack Size:** 3,000 ships per stack (max 27,000 ships per fleet)
- **Rule:** One ship design per stack, do NOT mix weapon types within a fleet

**Grid Positions & Attack Power (lines 420-432):**
```
+------+----------+----------+
| Head | Shoulder | Shoulder |   First Rank:  100% attack
+------+----------+----------+
| Flank| Glasshouse| Flank   |   Second Rank:  90% attack
+------+----------+----------+
| Rear |   Tail   | Rear     |   Third Rank:  75% attack
+------+----------+----------+
```

**Position Details:**
- **Glasshouse** (center-middle): Most protected position
- **Shoulders:** Most vulnerable positions
- **Fleet speed:** Speed of the slowest ship

**Fleet Formations (lines 436-444):**
| Formation | Description |
|-----------|-------------|
| Phalanx | All 9 stacks filled; maximum firepower |
| Diamond | 5 stacks (center + 4 adjacent); protects central stack |
| Battle Line | 6 stacks (first 2 ranks); frontal defense |
| Skirmish | Sparse; minimizes scatter damage |
| Tee Forward | T-shape; balanced offense |
| Enfilade | Side-focused fire |
| Tee Reverse | Reversed T-shape |

**Implementation Status:**
- Task #14 pending: "Create fleet formation UI with 6 ship slots"
- **DISCREPANCY:** Task mentions "6 ship slots" but GO2 mechanics specify 9 slots (3x3 grid)
- **Recommendation:** Clarify if we're using simplified 6-slot or full 9-slot system

**Combat Readiness Formula (line 449-452):**
```
Combat Readiness = (Ships in Stack / Effective Stack) x 100%
```

Example: 900 Battleships in a stack of 900 = 100% readiness. Only ships up to the effective stack number can attack each round.

---

## 5. Commander System Mechanics

### Finding: Gacha-Based Commander Recruitment with Star Ranking

**Overview:**
- Commanders provide combat bonuses and increase Effective Stack
- Star ranking system (0-15 stars) through merging duplicates
- Gacha recruitment system with 3 rarity tiers

**Rarity Tiers:**
- **Common:** 50% drop rate (15 commanders)
- **Skill:** 35% drop rate (10 commanders)
- **Super:** 15% drop rate (5 commanders)
- **Legendary/Divine:** OUT of scope for MVP

**Commander Stats:**
| Attribute | Effect |
|-----------|--------|
| Accuracy | Increases weapon hit chance |
| Dodge | Reduces opponent hit rate |
| Speed | Determines attack order; affects successive strike chance |
| Electron | Increases Critical Hit Rate and Critical Damage |

**Star Rank System (from Module 4 plan):**
- **Max Star Rank:** 15
- **Stat Bonus Formula:** `Star N: Base × (1 + N×0.10)`
  - Star 0: Base stats (no bonus)
  - Star 1: Base × 1.10 (+10%)
  - Star 5: Base × 1.50 (+50%)
  - Star 15: Base × 2.50 (+150%)

**Effective Stack Bonus (MOST IMPORTANT):**
```
Base Effective Stack (no commander): 300 ships
With Commander:
  Effective Stack = 300 + (Star Rank × 50)

Star 0: 300 (no bonus)
Star 1: 350 (+50)
Star 5: 550 (+250)
Star 10: 800 (+500)
Star 15: 1050 (+750)
```

**This means:** A Star 15 commander allows 3.5× more ships to attack per round!

**Recruitment Mechanics:**
- **Building:** Command Center (12 levels)
- **Cost:** 10,000 Gold per recruitment
- **Cooldown:** 3 hours (Level 1) → 1h10m (Level 12)
- **Max Commanders:** 60 per player

**Merging System:**
- **Building:** Compound Center
- **Merge Costs (estimated):**
  - Star 0 → Star 1: 2 duplicates + 5,000 Gold
  - Star 1 → Star 2: 3 duplicates + 10,000 Gold
  - Star N → Star N+1: (N+2) duplicates + (5000 × 2^N) Gold

**Implementation Status:**
- Task #15 pending: "Complete commanders feature [blocked by #1]"
- Full implementation plan exists at `/home/yurei/cryptomines-online/docs/planning/module-4-commander-system-plan.md`

**Combat Impact:**
- Weapon Expertise: S grade = +30% damage, D grade = -30% damage
- Ship Expertise: S grade = +10% dealt / -10% received
- See GO2 mechanics lines 558-575 for full details

---

## 6. Galaxy/Map System Mechanics

### Finding: 7x7+ Grid Galaxy with Resource Bonus Planets (RBPs)

**Galaxy Layout (from GO2 mechanics, line 666):**
- **Structure:** 7x7+ grid of zones
- **RBPs:** Each zone has 1 RBP in the center, spaced 60 movement squares apart
- **Purpose:** Territorial control for production bonuses

**RBP Bonuses (lines 670-678):**
| Level Range | Bonus per Level | Cumulative at Top |
|-------------|----------------|-------------------|
| 1-10 | 5% base + 0.5%/level | 10% |
| 11-20 | 1% per level | 20% |
| 21-30 | 1.5% per level | 35% |
| 31-40 | 2% per level | 55% |
| ... | Continues scaling | ... |
| 100 | - | **280%** |

**Bonuses Apply To:**
- Resource production
- Research speed
- Shipbuilding speed

**RBP Defense (lines 683-703):**
- **Initial NPC Defenses:** 5 fleets of 4,500-7,200 Level 6 ships each
- **Defensive Structures:**
  - Meteor Stars: 63 units
  - Particle Cannons: 8 units
  - Anti-Aircraft Guns: 12 units
  - Thor's Cannons: 5 units
- **Fleet Capacity:** Scales from 10 (Level 1) to 105 (Level 100)

**Conquest Mechanics (lines 706-716):**
- Only **Corps** (not individuals) can attack RBPs
- **72-hour protection** after capture
- **24-hour battle phase** when vulnerable
- **99 combat rounds** to complete takeover
- Fleets cannot be recalled mid-battle

**Corp Control:**
- Corps can control one planet per Corp Level (Lv 10 corp = max 10 RBPs)
- Upgraded using Corp Wealth (from member donations)

**Implementation Status:**
- Task #16 pending: "Complete galaxy/map feature [blocked by #1]"
- Large, complex system - likely Phase 2 or later

**Galactic Features:**
- **Galaxy Transporter:** Building for inter-system resource transport
- **Movement:** Measured in movement squares
- **Distance:** 60 movement squares between RBPs

---

## Summary of Key Findings

### 1. Warehouse Limit
**Answer:** Maximum **4 Resource Warehouses** per planet
- Database implementation: `max_count_per_planet = 4`
- Storage capacity stacks or uses highest value

### 2. Technology Center Requirements
**Answer:** Required for research, provides 3% time reduction per level (max 36% at Level 12)
- Separate from Weapon Research Center
- Handles 7 science trees
- Only Gold cost for research

### 3. Weapon Center Requirements
**Answer:** Weapon Research Center required for blueprint research
- 3% blueprint research time reduction per level (max 36% at Level 12)
- Separate building from Technology Center
- Handles ship hull and module blueprints

### 4. Fleet Formation Mechanics
**Answer:** 3x3 grid with 9 stack slots (not 6)
- Each stack holds up to 3,000 ships
- Attack power varies by rank: 100%/90%/75%
- 7 formation types available

### 5. Commander System
**Answer:** Gacha recruitment with star ranking
- 3 rarity tiers (Common/Skill/Super)
- Star rank 0-15 increases stats and Effective Stack
- Most important upgrade: Effective Stack bonus
- Max 60 commanders per player

### 6. Galaxy/Map System
**Answer:** 7x7+ grid with Resource Bonus Planets
- RBPs provide up to 280% bonus at Level 100
- Corp-based territorial control
- 60 movement squares between RBPs
- Complex defense and conquest mechanics

---

## Discrepancies Found

### 1. Fleet Formation Slots
- **Task #14:** Mentions "6 ship slots"
- **GO2 Mechanics:** Specifies 9 slots (3x3 grid)
- **Recommendation:** Verify if simplified 6-slot or full 9-slot implementation

### 2. GDD vs Database
- **GDD:** Shows "Multiple" warehouses without specific count
- **Database:** Implements exactly 4 warehouses
- **Recommendation:** Update GDD to specify "4" instead of "Multiple"

---

## Files Referenced

1. `/home/yurei/cryptomines-online/docs/gdd/game-design-document.md`
2. `/home/yurei/cryptomines-online/docs/research/galaxy-online-2-mechanics.md`
3. `/home/yurei/cryptomines-online/docs/planning/module-4-commander-system-plan.md`
4. `/home/yurei/cryptomines-online/supabase/migrations/20260206005232_phase1_mvp.sql`
5. `/home/yurei/cryptomines-online/docs/planning/module-3-resource-auto-production-plan.md`

---

## Recommendations for Implementation

### Immediate Actions

1. **Update Task #4:** "Fix multiple warehouses validation"
   - Enforce max 4 warehouses per planet
   - Clarify storage capacity stacking behavior

2. **Update Task #6:** "Add Technology Center requirement validation for research"
   - Validate Technology Center level before allowing tech research
   - Apply 3% time reduction per level

3. **Update Task #7:** "Add Weapon Center requirement validation for military research"
   - Validate Weapon Research Center level before blueprint research
   - Apply 3% time reduction per level

4. **Update Task #14:** "Create fleet formation UI with 6 ship slots"
   - Verify correct slot count (6 vs 9)
   - Implement position-based attack power modifiers

5. **Update Task #15:** "Complete commanders feature"
   - Implement gacha recruitment system
   - Add star ranking and merging
   - Integrate Effective Stack bonuses

6. **Update Task #16:** "Complete galaxy/map feature"
   - Large scope - consider breaking into sub-tasks
   - Prioritize core galaxy grid and RBP placement
   - Defer complex conquest mechanics to later phase

---

**Research Complete: 2026-02-10**
