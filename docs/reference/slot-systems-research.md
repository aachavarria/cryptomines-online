# GO2 Slot & Capacity Systems Research

**Date:** 2026-02-06
**Status:** Research Complete
**Sources:** Galaxy Online II Fandom Wiki, GameYum guides, KRTools, Google Sites FAQ, community guides

---

## 1. Planet Building Slots (Grid / Placement)

### Research Findings

- **Grid Size:** The planet base uses an isometric tile grid. The exact grid dimensions are not explicitly documented in any public wiki source for GO2. Community references mention large base layouts with dozens of buildings placed.
- **No "building slot" limit:** GO2 does NOT use a numbered slot system for buildings. Instead, buildings occupy multi-tile footprints on a 2D grid. Players can place as many buildings as physically fit on the grid, constrained by footprint size and available empty tiles.
- **Civic Center does NOT unlock building slots.** The Civic Center gates the *level* of other buildings (e.g., Civic Center Lv3 required to upgrade Metal Collector to Lv5), but does not unlock additional grid slots. All grid tiles are available from the start.
- **Building count limits:** Some building types have fixed max counts (e.g., 8 Metal Collectors, 8 He3 Extractors, 8 Residential Areas), but this is a per-type cap, not a grid slot limit.

### GDD Status

**Current GDD (Section 2.2.9):** Documents a 20x20 isometric grid (400 tiles) with multi-tile footprints. This is **well-documented and accurate** to GO2's approach -- placement is spatial, not slot-based.

**Gap:** None significant. The 20x20 size is our design decision (exact GO2 grid dimensions are undocumented publicly). The footprint system matches GO2.

---

## 2. Construction Queue Slots

### Research Findings

- **Default construction slots:** 2 concurrent building upgrades by default.
- **Concurrent Construction tech:** Researching "Concurrent Construction" (Logistics Construction tree, Lv1) adds +1 slot, bringing the total to **3 permanent construction slots**.
- **Construction Cards:** A consumable mall item that "activates 3 extra slots in the Construction Schedule for 72 hours." With a Construction Card active, total becomes **5 or 6 construction slots** (2 base + 1 from tech + 3 from card = 6 max, or 2 base + 3 from card = 5 without tech).
- **"5 construction slots":** Multiple guides reference "use all 5 construction slots" which matches: 2 default + 3 from Construction Card = 5 (pre-research), or 3 permanent + 3 card = 6 (post-research). The common advice of "5 slots" likely refers to the Construction Card being active before Concurrent Construction research is done.

### GDD Status

**Current GDD (Section 2.2.10):** States "Limited concurrent construction (default 2); use Construction Cards for additional slots." Also documents Concurrent Construction tech at max_level=1 adding +1 slot.

**Current GDD (Section 7.2):** States "Max 2 concurrent construction slots (or more with Construction Cards)."

**Current GDD (Section 3.28 seed data):** `concurrent_construction` has max_level=1, effects `{"type":"construction_slots","per_level":1}`.

**Assessment:** The GDD is **accurate**. Default 2, +1 from Concurrent Construction (permanent), +3 from Construction Card (temporary 72h). The progression is:
- Base: 2 slots
- After Concurrent Construction Lv1: 3 slots
- With Construction Card active: +3 temporary (up to 5 or 6 total)

**Gap:** The GDD mentions Construction Cards but doesn't specify they add exactly 3 extra slots for 72 hours. This detail should be added to the GDD mall items section when documented.

---

## 3. Ship Factory Production Slots

### Research Findings

- **Total production slots:** 5 maximum
- **Slot unlocking by Ship Factory level:**
  - Slot 1: Ship Factory Level 1
  - Slot 2: Ship Factory Level 4
  - Slot 3: Ship Factory Level 8
  - Slot 4: Ship Factory Level 12
  - Slot 5: "Sync Shipbuilding" tech (Logistics Construction tree, requires Ship Building Logistics Lv4)
- **Max production per slot:** 2,000,000 ships
- **Design capacity:** 20 ship blueprints stored in the factory
- **Speed bonus:** Ship Factory level increases production speed from 1% (Lv1) to 60% (Lv24)

### GDD Status

**Current GDD (Section 8.1):** Fully documented with correct slot unlocking table, speed bonus table (Lv1-24), SQL schema for `ship_factory_levels`, and production slot column in the level data.

**Assessment:** **Fully accurate and complete.** No gaps.

---

## 4. Fleet Slots (Max Fleets per Player)

### Research Findings

- **Fleet count is tied to player level**, NOT to Command Center level.
- **Command Center** handles commander recruitment (draw/recruit commanders), NOT fleet count.
- **Max commanders at Player Level 71+:** 60 commanders = 60 fleets.
- The exact per-level breakdown below 71 is not comprehensively documented in public sources. One screenshot showed a Lv84 player with 51 active fleets.
- Fleet count scales with player level; the progression is gradual.

### GDD Status

**Current GDD (Section 8.5.6):**

| Player Level | Max Fleets |
|-------------|------------|
| 1-10 | 2 |
| 11-20 | 3 |
| 21-30 | 4 |
| 31-40 | 5 |
| 41-50 | 6 |
| 51+ | 8 |

**Assessment:** This table is a **reasonable approximation** for our game design but likely does NOT match GO2 exactly. GO2 scales to 60 commanders/fleets at level 71+, implying a much finer-grained and higher-count progression. The GDD's max of 8 fleets at level 51+ is a significant simplification compared to GO2's 60-fleet endgame.

**Gap / Design Decision Needed:** The GDD fleet limit table should be reviewed. Options:
1. **Keep simplified** (current): Max 8 fleets. Simpler for our implementation.
2. **Match GO2 more closely:** Scale to 60 fleets by level 71+, with per-level unlocks (roughly +1 fleet per 1-2 player levels). This creates much more endgame fleet management complexity.

**Recommendation:** Keep the simplified table for Phase 2 launch. Revisit when implementing PvP / high-level content. Note in the GDD that GO2's endgame goes to 60 fleets.

---

## 5. Fleet Formation Grid Slots

### Research Findings

- **Grid:** 3x3 = **9 position slots** per fleet
- **Max ships per stack:** 3,000
- **Max fleet size:** 27,000 ships (9 x 3,000)
- **Rule:** One ship design per stack (no mixing within a stack)
- **Formations** determine which of the 9 positions are active:
  - Phalanx: All 9
  - Diamond: 5 (center cross)
  - Battle Line: 6 (first 2 ranks)
  - Skirmish: 5 (corners + center)
  - Tee Forward: 5
  - Enfilade: 5
  - Tee Reverse: 5
- **Attack power modifiers by rank:**
  - First Rank (row 0): 100%
  - Second Rank (row 1): 90%
  - Third Rank (row 2): 75%

### GDD Status

**Current GDD (Sections 2.4.4, 8.5.1-8.5.3):** Fully documented with grid layout, position names, attack modifiers, formation active positions, SQL schema with `fleet_stacks` table.

**Assessment:** **Fully accurate and complete.** No gaps. The formation active positions are documented with exact grid coordinates in section 8.5.3.

---

## 6. Spacedock Repair Slots

### Research Findings

- **No concurrent repair slot limit.** The Spacedock does NOT use a "repair slot" system in the traditional sense.
- **Repair is probabilistic:** After PvP combat, destroyed ships go to the Spacedock. Each individual ship has an independent probability of being repaired based on the Spacedock's repair percentage (1% at Lv1, 20% at Lv12).
- **Storage limit:** 2 pages maximum, storing the 10 oldest ship designs. If >10 different designs are destroyed simultaneously, designs beyond the 10 oldest are permanently lost.
- **Repair time:** Ships take time to repair. Players can spend 10 SP (Mall Points) to reduce repair time by 10%.
- **PvP only:** Ships lost in Normal Instances are permanently destroyed and CANNOT be repaired.
- **Repair Technology** (Logistics Construction tech): Increases repair percentage by +1-10% across 10 levels.

### GDD Status

**Current GDD (Sections 8.4, 8.8.10):** Fully documented with level data table (12 levels), repair percentages, storage limits, SQL schema for `spacedock_repairs`.

**Assessment:** **Fully accurate and complete.** The GDD correctly models the probabilistic repair system, the 10-design storage limit, PvP-only restriction, and Repair Technology enhancement.

**Minor note:** The `spacedock_repairs` table in the GDD (section 8.8.10) includes a `repair_finish_at` timestamp, implying a time-based repair queue. This is consistent with GO2 where repairs take time after the probability roll succeeds.

---

## 7. Research Slots

### Research Findings

- **One research active per tree at a time.**
- **7 trees = up to 7 concurrent researches.**
- Researching in one tree does NOT block research in other trees.
- Technology Center level reduces research time by 3% per level (36% at Lv12).
- Research costs are Gold only.

### GDD Status

**Current GDD (Section 2.3.9):** "Only one research can be active per tree at a time (7 trees = up to 7 concurrent researches)"

**Current GDD (Section 7.2):** "Only 1 research active per tree"

**Assessment:** **Fully accurate and complete.** No gaps.

---

## Summary: GDD Gap Analysis

| System | GDD Status | Action Needed |
|--------|-----------|---------------|
| Planet Building Grid | Complete (20x20, footprint-based) | None |
| Construction Queue Slots | Complete (default 2, +1 from tech) | Minor: Document Construction Card = +3 for 72h |
| Ship Factory Production Slots | Complete (5 slots, level-gated) | None |
| Fleet Slots (Max Fleets) | Simplified (max 8) vs GO2's 60 | Design decision: keep simplified or expand |
| Fleet Formation Grid | Complete (3x3, 9 positions) | None |
| Spacedock Repair | Complete (probabilistic, 12 levels) | None |
| Research Slots | Complete (1 per tree, 7 concurrent) | None |

### Key Findings

1. **Construction Cards** add exactly 3 extra temporary slots for 72 hours. This is a premium/mall mechanic.
2. **GO2 endgame has 60 fleets** (at player level 71+), far more than our GDD's simplified 8-fleet max. This is a deliberate simplification that should be noted.
3. **No building slot system** -- GO2 uses spatial grid placement, matching our GDD's 20x20 grid approach.
4. **Research is 1-per-tree** (7 concurrent), confirmed across multiple sources.
5. All slot systems in our GDD are **accurate to GO2** except fleet count, which is intentionally simplified.

---

## Sources

- [Galaxy Online II Wiki - Guide To Advancing Quickly](https://galaxyonlineii.fandom.com/wiki/Guide_To_Advancing_Quickly)
- [Galaxy Online II Wiki - Command Center](https://galaxyonlineii.fandom.com/wiki/Command_Center)
- [Galaxy Online II Wiki - Fleet Design](https://galaxyonlineii.fandom.com/wiki/Fleet_Design)
- [Galaxy Online II Wiki - Spacedock](https://galaxyonlineii.fandom.com/wiki/Spacedock)
- [Galaxy Online II Wiki - Civic Center](https://galaxyonlineii.fandom.com/wiki/Civic_Center)
- [Galaxy Online II Wiki - Technology Center](https://galaxyonlineii.fandom.com/wiki/Technology_Center)
- [Galaxy Online II Wiki - Development Items](https://galaxyonlineii.fandom.com/wiki/Development_Items)
- [Galaxy Online II Wiki - Construction](https://galaxyonlineii.fandom.com/wiki/Construction)
- [Galaxy Online II Wiki - Beginner's Guide](https://galaxyonlineii.fandom.com/wiki/Beginner's_Guide)
- [Galaxy Online II Wiki - Normal Instances](https://galaxyonlineii.fandom.com/wiki/Normal_Instances)
- [GameYum - Galaxy Online II Structures Guide](https://www.gameyum.com/galaxy-online/115041-galaxy-online-ii-guides-building-structures/)
- [Google Sites - Galaxy Online II FAQ](https://sites.google.com/site/galaxyonlineii/faq)
- [KRTools - Commander Planner](https://krtools.deajae.co.uk/cc/)
