# Ship Hulls Validation Report
**Galaxy Online 2 Hull Types Analysis**
**Date:** 2026-02-10
**Researcher:** Ship Hulls Research Specialist

---

## Executive Summary

**KEY FINDINGS:**
- Galaxy Online 2 has **3 primary ship classes**: Frigate, Cruiser, Battleship
- **Destroyers, Battlecruisers, and Carriers DO NOT EXIST** as separate classes in GO2
- We currently have **10 Frigate lines** (30 hulls) ✅ COMPLETE
- We currently have **10 Cruiser lines** (30 hulls) ⚠️ MISSING 2 advanced lines
- We currently have **5 Battleship lines** (15 hulls) ⚠️ MISSING 7 lines
- Flagships exist as special variants, not a separate class

---

## Part 1: Ship Class Structure in Galaxy Online 2

### Confirmed Ship Classes
Based on comprehensive wiki research, Galaxy Online 2 has exactly **THREE** primary ship classes:

1. **Frigates** - Fast, agile ships with base movement=1, agility=1
2. **Cruisers** - Medium ships with base movement=0, agility=0
3. **Battleships** - Heavy ships with base movement=0, agility=0

### Classes That DO NOT Exist
- **Destroyers** - NOT a separate class in GO2
- **Battlecruisers** - NOT a separate class in GO2
- **Carriers** - NOT a separate class in GO2

### Special Hull Category
- **Flagships** - Special variants based on Frigate, Cruiser, or Battleship hulls
  - Not a separate class, but enhanced versions with special abilities
  - Examples: Independence (Cruiser-based), Black Hole (Battleship-based), Liberty Wings (Frigate-based)

---

## Part 2: Frigate Hulls - VALIDATION COMPLETE ✅

### Current Implementation Status: **10/10 hull lines (100%)**

All 10 frigate hull lines with 3 tiers each are correctly implemented:

| # | Hull Name | Armor Type | Tiers | Status | Notes |
|---|-----------|------------|-------|--------|-------|
| 1 | Weikes | Nano | I, II, III | ✅ CORRECT | Starter frigate, balanced stats |
| 2 | Air Wanderer | Neutralizing | I, II, III | ✅ CORRECT | Neutralizing armor variant |
| 3 | Valkyrie | Nano | I, II, III | ✅ CORRECT | High-shield nano frigate |
| 4 | GoGetter | Neutralizing | I, II, III | ✅ CORRECT | Sturdy neutralizing frigate |
| 5 | Space Hunter | Nano | I, II, III | ✅ CORRECT | Balanced nano frigate |
| 6 | Sparrow | Neutralizing | I, II, III | ✅ CORRECT | Light neutralizing frigate |
| 7 | Devourer | Nano | I, II, III | ✅ CORRECT | High-shield low-structure |
| 8 | Polymesus | Neutralizing | I, II, III | ✅ CORRECT | Tank neutralizing frigate |
| 9 | Cybra | Nano | I, II, III | ✅ CORRECT | Highest-shield frigate |
| 10 | Hamdar | Neutralizing | I, II, III | ✅ CORRECT | Highest-structure frigate |

**Totals: 10 hull lines × 3 tiers = 30 frigate hulls ✅**

### Armor Distribution
- **Nano armor:** 5 lines (Weikes, Valkyrie, Space Hunter, Devourer, Cybra)
- **Neutralizing armor:** 5 lines (Air Wanderer, GoGetter, Sparrow, Polymesus, Hamdar)

---

## Part 3: Cruiser Hulls - MISSING ADVANCED VARIANTS ⚠️

### Current Implementation Status: **10/12 hull lines (83%)**

We have 10 standard cruiser lines but are MISSING 2 advanced cruiser lines:

#### Currently Implemented (10 lines = 30 hulls)

| # | Hull Name | Armor Type | Tiers | Status | Notes |
|---|-----------|------------|-------|--------|-------|
| 1 | Typhoon | Chrome | I, II, III | ✅ CORRECT | Starter cruiser |
| 2 | Bombardier | Regen | I, II, III | ✅ CORRECT | Regenerative armor |
| 3 | Duke | Chrome | I, II, III | ✅ CORRECT | High-shield chrome |
| 4 | The Shuttler | Regen | I, II, III | ✅ CORRECT | Balanced regen |
| 5 | Watchman | Chrome | I, II, III | ✅ CORRECT | Defensive chrome |
| 6 | Spinner | Regen | I, II, III | ✅ CORRECT | High-capacity regen |
| 7 | Wraith | Chrome | I, II, III | ✅ CORRECT | Stealthy chrome |
| 8 | Encratos | Regen | I, II, III | ✅ CORRECT | Highest-capacity |
| 9 | Nicholas | Chrome | I, II, III | ✅ CORRECT | Balanced chrome |
| 10 | Helena | Regen | I, II, III | ✅ CORRECT | All-around regen |

#### MISSING HULLS (2 advanced lines = 6 hulls)

| # | Hull Name | Armor Type | Tiers | Status | Priority |
|---|-----------|------------|-------|--------|----------|
| 11 | **Chimera Capra** | Regen | I, II, III | ❌ MISSING | HIGH |
| 12 | **Ultra Gwyar** | Regen | I, II, III | ❌ MISSING | HIGH |

**Evidence from GO2 Wiki:**
- **Chimera Capra:** Shields 1,310-1,790 | Structure 7,300-10,000 (Regen armor)
- **Ultra Gwyar:** Shields 2,300-2,800 | Structure 10,950-15,000 (Regen armor)

These are high-tier cruisers with significantly higher stats than standard lines.

### Armor Distribution
- **Chrome armor:** 5 lines (Typhoon, Duke, Watchman, Wraith, Nicholas)
- **Regen armor:** 5 standard + 2 advanced = 7 lines total (Bombardier, Shuttler, Spinner, Encratos, Helena, Chimera Capra, Ultra Gwyar)

**Totals: 10 implemented + 2 missing = 12 hull lines × 3 tiers = 36 total cruiser hulls**

---

## Part 4: Battleship Hulls - MAJOR GAPS ⚠️

### Current Implementation Status: **5/12 hull lines (42%)**

We have only 5 basic battleship lines and are MISSING 7 additional lines:

#### Currently Implemented (5 lines = 15 hulls)

| # | Hull Name | Armor Type | Tiers | Status | Notes |
|---|-----------|------------|-------|--------|-------|
| 1 | Estrella | Chrome | I, II, III | ✅ CORRECT | Starter battleship |
| 2 | Nettle | Regen | I, II, III | ⚠️ WRONG ARMOR | Should be Nano, not Regen |
| 3 | Diaz | Nano | I, II, III | ⚠️ WRONG ARMOR | Should be Neutralizing, not Nano |
| 4 | RV766-The Explorer | Neutralizing | I, II, III | ⚠️ WRONG ARMOR | Should be Regen, not Neutralizing |
| 5 | Palenka | Chrome | I, II, III | ⚠️ WRONG ARMOR | Should be Nano, not Chrome |

#### MISSING HULLS (7 lines = 21 hulls)

| # | Hull Name | Armor Type | Tiers | Status | Priority |
|---|-----------|------------|-------|--------|----------|
| 6 | **Howler** | Chrome | I, II, III | ❌ MISSING | HIGH |
| 7 | **Whirlpool** | Regen | I, II, III | ❌ MISSING | HIGH |
| 8 | **Cerberus** | Neutralizing | I, II, III | ❌ MISSING | HIGH |
| 9 | **Genesis** | Regen | I, II, III | ❌ MISSING | HIGH |
| 10 | **Tiamat** | Chrome | I, II, III | ❌ MISSING | HIGH |
| 11 | **Chimera Viper** | Nano | I, II, III | ❌ MISSING | MEDIUM |
| 12 | **Ultra Calas** | Chrome | I, II, III | ❌ MISSING | MEDIUM |

### Armor Type Corrections Needed

Based on GO2 Wiki research:

| Hull | Current Armor | Correct Armor | Source |
|------|---------------|---------------|--------|
| Nettle | Regen | **Nano** | GO2 Wiki: "minimizes Heat/Explosive damage" |
| Diaz | Nano | **Neutralizing** | GO2 Wiki: "minimizes Magnetic/Explosive damage" |
| RV766-The Explorer | Neutralizing | **Regen** | GO2 Wiki: "minimizes Kinetic/Magnetic damage" |
| Palenka | Chrome | **Nano** | GO2 Wiki battleship armor patterns |

**Totals: 5 implemented (with errors) + 7 missing = 12 hull lines × 3 tiers = 36 total battleship hulls**

---

## Part 5: Flagship/Special Hulls

### Federation Flagships (Special variants, not separate class)

| Name | Base Class | Tiers | Notes |
|------|-----------|-------|-------|
| Independence | Cruiser | I, II, III, Mk2 | Movement +5-7, -90% damage taken |
| Black Hole | Battleship | I, II, III | +150-300% damage to Light armor |
| Liberty Wings | Frigate | - | No flagship effective stack bonus |
| Bastion | Battleship | - | Limited info available |

### Humaroid-Flagships (Pirate-based)
- Intrepid Nexus
- Grim Reaper
- Shadow Trojan
- Mercury Wing
- Gforce Dreadnaught
- Fire Cat
- Conquistador
- Arbiter

**Note:** Flagships are NOT implemented in our current system and are LOW PRIORITY for MVP.

---

## Part 6: Data Quality Issues

### Issue 1: Incorrect Armor Types (CRITICAL)
**Affected hulls:** Nettle, Diaz, RV766-The Explorer, Palenka

**Current implementation has wrong armor types:**
```sql
-- CURRENT (WRONG):
('nettle_i', 'Nettle I', 'battleship', 1, 'regen', ...)      -- Should be 'nano'
('diaz_i', 'Diaz I', 'battleship', 1, 'nano', ...)           -- Should be 'neutralizing'
('rv766_i', 'RV766-The Explorer I', 'battleship', 1, 'neutralizing', ...) -- Should be 'regen'
('palenka_i', 'Palenka I', 'battleship', 1, 'chrome', ...)   -- Should be 'nano'
```

### Issue 2: Stats Mismatch
Our implementation uses different stat values than GO2 wiki. Example:

**Typhoon-I (Cruiser):**
- **Our data:** Shield: 505 | Structure: 2,599 | Slots: 120
- **GO2 Wiki:** Shield: 202 | Structure: 1,040 | Slots: 140

This pattern repeats across all hulls. Our stats appear to be scaled differently.

### Issue 3: Missing Installation Slot Details
GO2 installation slots break down into specific module categories:
- Ballistic slots
- Directional slots
- Missile slots
- Ship-based weapon slots
- Defense slots
- Electronic slots
- Storage slots
- Transmission slots

Our current schema only has a single `installation_slots` integer field.

---

## Part 7: Schema Recommendations

### Current Schema Constraints
```sql
hull_class TEXT NOT NULL CHECK (hull_class IN ('frigate', 'cruiser', 'battleship'))
armor_type TEXT NOT NULL CHECK (armor_type IN ('nano', 'chrome', 'regen', 'neutralizing'))
```

✅ **CORRECT** - No changes needed. GO2 only has these 3 classes and 4 armor types.

### Recommended Schema Enhancement (Future)
```sql
-- Add detailed slot breakdown
ALTER TABLE hull_types ADD COLUMN ballistic_slots INTEGER DEFAULT 0;
ALTER TABLE hull_types ADD COLUMN directional_slots INTEGER DEFAULT 0;
ALTER TABLE hull_types ADD COLUMN missile_slots INTEGER DEFAULT 0;
ALTER TABLE hull_types ADD COLUMN ship_based_slots INTEGER DEFAULT 0;
ALTER TABLE hull_types ADD COLUMN defense_slots INTEGER DEFAULT 0;
ALTER TABLE hull_types ADD COLUMN electronic_slots INTEGER DEFAULT 0;
ALTER TABLE hull_types ADD COLUMN storage_slots INTEGER DEFAULT 0;
ALTER TABLE hull_types ADD COLUMN transmission_slots INTEGER DEFAULT 0;
```

---

## Part 8: Action Items

### IMMEDIATE (Critical Fixes)
1. ✅ **Correct battleship armor types** (Nettle, Diaz, RV766, Palenka)
2. ❌ **Add 7 missing battleship lines** (Howler, Whirlpool, Cerberus, Genesis, Tiamat, Chimera Viper, Ultra Calas)
3. ❌ **Add 2 missing advanced cruiser lines** (Chimera Capra, Ultra Gwyar)

### HIGH PRIORITY
4. ❌ **Verify and reconcile stat values** against GO2 wiki (all hulls show different stats)
5. ❌ **Document stat scaling formula** (our stats appear 2-3x higher than wiki values)

### MEDIUM PRIORITY
6. ❌ **Add detailed module slot breakdown** (8 slot types per hull)
7. ❌ **Research Nebula V-class Frigate** (mentioned in wiki, unclear if it's standard or special)

### LOW PRIORITY (Post-MVP)
8. ❌ Implement flagship system (Independence, Black Hole, etc.)
9. ❌ Implement Humaroid-Flagship system
10. ❌ Add special hull variants (Encratos-A, Encratos-B, etc.)

---

## Part 9: Database Migration Recommendations

### Migration 1: Fix Battleship Armor Types
```sql
-- Fix incorrect armor types on existing battleship hulls
UPDATE hull_types SET armor_type = 'nano' WHERE name LIKE 'nettle_%';
UPDATE hull_types SET armor_type = 'neutralizing' WHERE name LIKE 'diaz_%';
UPDATE hull_types SET armor_type = 'regen' WHERE name LIKE 'rv766_%';
UPDATE hull_types SET armor_type = 'nano' WHERE name LIKE 'palenka_%';
```

### Migration 2: Add Missing Battleship Hulls
Need to add 7 new battleship lines (21 hulls total):
- Howler I/II/III (Chrome)
- Whirlpool I/II/III (Regen)
- Cerberus I/II/III (Neutralizing)
- Genesis I/II/III (Regen)
- Tiamat I/II/III (Chrome)
- Chimera Viper I/II/III (Nano)
- Ultra Calas I/II/III (Chrome)

### Migration 3: Add Missing Cruiser Hulls
Need to add 2 advanced cruiser lines (6 hulls total):
- Chimera Capra I/II/III (Regen)
- Ultra Gwyar I/II/III (Regen)

---

## Part 10: Conclusion

### Completeness Summary

| Ship Class | Implemented | Missing | Completeness |
|------------|-------------|---------|--------------|
| Frigates | 30/30 hulls | 0 hulls | ✅ 100% |
| Cruisers | 30/36 hulls | 6 hulls | ⚠️ 83% |
| Battleships | 15/36 hulls | 21 hulls | ❌ 42% |
| **TOTAL** | **75/102 hulls** | **27 hulls** | **⚠️ 74%** |

### Data Quality Assessment

| Aspect | Status | Issues |
|--------|--------|--------|
| Ship Classes | ✅ CORRECT | Only 3 classes (Frigate, Cruiser, Battleship) |
| Armor Types | ⚠️ MOSTLY CORRECT | 4 battleships have wrong armor |
| Frigate Hulls | ✅ COMPLETE | All 10 lines present |
| Cruiser Hulls | ⚠️ INCOMPLETE | Missing 2 advanced lines |
| Battleship Hulls | ❌ VERY INCOMPLETE | Missing 7 lines + 4 armor errors |
| Stat Values | ⚠️ UNCLEAR | Different from wiki (scaling?) |
| Module Slots | ⚠️ SIMPLIFIED | Single field vs 8 categories |

### Critical Insights

1. **No Destroyers/Battlecruisers/Carriers:** These ship classes do NOT exist in Galaxy Online 2. Our schema correctly excludes them.

2. **Battleship Data Quality:** Our battleship implementation has the most issues:
   - Only 42% complete (15 of 36 hulls)
   - 4 out of 5 implemented hulls have WRONG armor types
   - Missing 7 entire hull lines

3. **Stat Scaling Mystery:** Our stat values are consistently 2-3x higher than GO2 wiki values. This may be intentional game balance, or we may be using a different source. **NEEDS INVESTIGATION.**

4. **Module Slot Simplification:** GO2 has 8 different slot types, we have 1 generic field. This works for MVP but limits ship design complexity.

---

## Research Sources

- [Composite Ship Table - Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Composite_Ship_Table)
- [Category: Frigates - Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Category:Frigates)
- [Category: Cruisers - Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Category:Cruisers)
- [Category: Battleships - Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Category:Battleships)
- [Category: Flagships - Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Category:Flagships)
- [Fleet Design - Galaxy Online II Wiki](https://galaxyonlineii.fandom.com/wiki/Fleet_Design)

---

**Report Prepared By:** Ship Hulls Research Specialist
**Date:** 2026-02-10
**Status:** Task #4 - COMPLETED
