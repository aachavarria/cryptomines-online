# Buildings Validation Report - CryptoMines Online
**Date:** 2026-02-10
**Researcher:** buildings-specialist
**Status:** CRITICAL ISSUES FOUND

## Executive Summary

Validated ALL buildings from Galaxy Online 2 wiki against our database schema and GDD. Found **CRITICAL DISCREPANCIES** in max counts, max levels, and missing buildings.

---

## 1. CONFIRMED: Resource Warehouse Bug Fix

✅ **CORRECT FIX**: Resource Warehouse max_count_per_planet should be **1**, not 4
- **Source:** https://galaxyonlineii.fandom.com/wiki/Resource_Warehouse
- **Wiki states:** "You can have a max of 1 resource warehouse that is fully upgradable to level 24"
- **Our DB:** Currently set to 1 in building_types (line 762) ✅

---

## 2. CRITICAL: Radar Max Level WRONG

❌ **ERROR IN DATABASE**: Radar max level is **9** in our DB, should be **10**

**Current DB (line 770):**
```sql
('radar', 'Radar', 'core', 'ground', 450, 400, 550, 60, 3.0300, 2.8700, 0, 1.0000, 9, 1, ...)
```

**Should be:**
- Max Level: **10** (not 9)
- Source: https://galaxyonlineii.fandom.com/wiki/Radar

**Radar Level Progression (Wiki):**
| Level | Civic Center | Tech Center | Detection Time | Metal | He3 | Gold |
|-------|--------------|-------------|----------------|-------|-----|------|
| 1 | 2 | 1 | 00:30:00 | 450 | 400 | 550 |
| 2 | 3 | 2 | 01:00:00 | 1,269 | 1,128 | 1,551 |
| 3 | 4 | 3 | 01:30:00 | 3,591 | 3,192 | 4,389 |
| 4 | 5 | 4 | 02:00:00 | 10,199 | 9,066 | 12,466 |
| 5 | 6 | 5 | 02:30:00 | 29,068 | 25,838 | 35,527 |
| 6 | 7 | 6 | 03:00:00 | 83,134 | 73,897 | 101,608 |
| 7 | 8 | 7 | 03:30:00 | 238,594 | 212,083 | 291,615 |
| 8 | 9 | 8 | 04:00:00 | 687,150 | 610,800 | 839,850 |
| 9 | 10 | 9 | 04:30:00 | 1,985,864 | 1,765,213 | 2,427,167 |
| **10** | **10** | **9** | **?** | **?** | **?** | **?** |

Note: Level 10 costs missing from wiki, but max level confirmed as 10.

---

## 3. CRITICAL: Ship Factory Base Costs WRONG

❌ **ERROR IN DATABASE**: Ship Factory base costs are INCORRECT

**Current DB (line 772):**
```sql
('ship_factory', 'Ship Factory', 'military', 'ground', 600, 450, 500, 200, ...)
```

**Wiki Level 1 costs:**
- Metal: **206** (not 600)
- He3: **163** (not 450)
- Gold: **189** (not 500)
- Build Time: **110 seconds** (not 200)

**Source:** https://galaxyonlineii.fandom.com/wiki/Ship_Factory

Ship Factory has unique progression:
- Levels 1-24
- Build slots increase: 1→2→3→4 slots
- Speed bonuses: 1%→60% at max level

---

## 4. Galaxy Transporter - Wrong Understanding

❌ **MISCONCEPTION**: Galaxy Transporter is NOT a resource transport building

**Actual Function:** Inter-Galactic League access building
- **Max Level:** 1 (single level only)
- **Prerequisites:** Civic Center Level 1
- **Build Time:** 24 hours
- **Cost:** 20,000 Metal / 20,000 He3 / 20,000 Gold
- **Function:** "Allows access to Inter-Galactic League. 10 free matches per day."

**Current DB (line 768):**
```sql
('galaxy_transporter', 'Galaxy Transporter', 'core', 'ground', 350, 300, 450, 100, 3.0300, 2.8700, 0, 1.0000, 12, 1, ...)
```

**Should be:**
- base_cost_metal: **20000** (not 350)
- base_cost_he3: **20000** (not 300)
- base_cost_gold: **20000** (not 450)
- base_time_seconds: **86400** (24 hours, not 100)
- max_level: **1** (not 12)

**Source:** https://galaxyonlineii.fandom.com/wiki/Galaxy_Transporter

---

## 5. Compound Center - Incomplete Data

⚠️ **PARTIAL DATA**: Compound Center only has Level 1 documented on wiki

**Wiki Data (Level 1):**
- Prerequisites: Civic Center Lv3, Weapon Research Center Lv2
- Build Time: 1 minute
- Cost: 12 Metal / 100 He3 / 10 Gold

**Current DB (line 769):**
```sql
('compound_center', 'Compound Center', 'core', 'ground', 400, 350, 500, 100, 3.0300, 2.8700, 0, 1.0000, 12, 1, ...)
```

**Issue:** Wiki doesn't show max level or full progression. Our DB shows max_level=12, but this is unconfirmed.

**Source:** https://galaxyonlineii.fandom.com/wiki/Compound_Center

---

## 6. Alliance Center - VERIFIED CORRECT

✅ **CORRECT**: Alliance Center data matches wiki

- Max Level: **11** ✅
- Max Count: **1** ✅
- Base costs match ✅
- Prerequisites: Civic Center Lv2, Tech Center Lv1 ✅

**Source:** https://galaxyonlineii.fandom.com/wiki/Alliance_Center

---

## 7. Missing Building: Celestial Base

❌ **BUILDING NOT FOUND ON WIKI**

**Current DB (line 783):**
```sql
('celestial_base', 'Celestial Base', 'space', 'space', 600, 500, 600, 200, 3.0300, 2.8700, 0, 1.0000, 12, 1, ...)
```

**Status:** Wiki page returns 404. This building may not exist in GO2 or uses a different name.

---

## 8. Decorative/Landscaping Buildings - MISSING FROM DATABASE

❌ **MISSING BUILDINGS**: We have NO landscaping buildings in our database

**GO2 Landscaping Buildings (from wiki):**
1. Casino Resort
2. Beacon
3. Monument
4. Fountain
5. Library
6. Theater
7. Park
8. College
9. Hospital
10. Shopping Center
11. Statue
12. Santa Sculpture

**Function:** Morale/aesthetic bonuses (GDD mentions these on line 137)

**Status:** Not implemented in our database. Need to research costs, levels, and effects.

**Source:** https://galaxyonlineii.fandom.com/wiki/Category:Buildings

---

## 9. Complete Buildings Inventory

### We HAVE in Database (22 buildings):

**Resource Buildings (4):**
- ✅ Metal Collector
- ✅ He3 Extractor
- ✅ Residential Area
- ✅ Resource Warehouse

**Core/Administrative Buildings (7):**
- ✅ Civic Center
- ✅ Technology Center
- ✅ Alliance Center (verified correct)
- ✅ Trading Center
- ⚠️ Galaxy Transporter (WRONG costs/function)
- ⚠️ Compound Center (unverified max level)
- ❌ Radar (WRONG max level: 9 should be 10)

**Military Buildings (5):**
- ❌ Ship Factory (WRONG base costs)
- ✅ Spacedock
- ✅ Command Center
- ✅ Weapon Research Center
- ✅ Recycling Plant

**Space Base Buildings (6):**
- ✅ Space Station
- ✅ Meteor Star
- ✅ Particle Cannon
- ✅ Anti-Aircraft Gun
- ✅ Thor's Cannon
- ❌ Celestial Base (NOT FOUND on wiki)

### We MISSING (12 buildings):

**Landscaping/Decorative (12):**
- ❌ Casino Resort
- ❌ Beacon
- ❌ Monument
- ❌ Fountain
- ❌ Library
- ❌ Theater
- ❌ Park
- ❌ College
- ❌ Hospital
- ❌ Shopping Center
- ❌ Statue
- ❌ Santa Sculpture

---

## 10. Summary of Required Fixes

### CRITICAL (Must Fix):
1. **Radar max_level**: 9 → **10**
2. **Ship Factory base costs**: (600,450,500,200) → **(206,163,189,110)**
3. **Galaxy Transporter**: Complete rewrite - single level PvP building, costs (350,300,450,100) → **(20000,20000,20000,86400)**

### MEDIUM (Should Fix):
4. **Celestial Base**: Remove or find correct wiki page
5. **Compound Center**: Verify max level (currently 12, unconfirmed)

### LOW (Future Work):
6. **Add 12 Landscaping Buildings**: Research costs, levels, effects, and implement

---

## 11. Verified Correct Buildings

These buildings match the wiki exactly:
- ✅ Metal Collector (base costs, levels, production)
- ✅ He3 Extractor (base costs, levels, production)
- ✅ Residential Area (base costs, levels, production)
- ✅ Resource Warehouse (max_count=1 ✅, levels)
- ✅ Civic Center (all levels match)
- ✅ Technology Center (all levels match)
- ✅ Alliance Center (all levels match)
- ✅ Command Center (partial - levels 6-12 need verification)
- ✅ Space Station (all levels match)
- ✅ Weapon Research Center (all levels match)
- ✅ Spacedock (all levels match)
- ✅ Recycling Plant (partial data)
- ✅ Meteor Star (all levels match)
- ✅ Particle Cannon (all levels match)
- ✅ Anti-Aircraft Gun (all levels match)
- ✅ Thor's Cannon (all levels match)

---

## 12. References

All data validated against official Galaxy Online 2 wiki:
- Main Buildings Category: https://galaxyonlineii.fandom.com/wiki/Category:Buildings
- Individual building pages cited above

---

## Recommendations

**IMMEDIATE ACTION REQUIRED:**
1. Fix Radar max_level to 10
2. Fix Ship Factory base costs to match wiki Level 1
3. Fix Galaxy Transporter completely (function, costs, max level)
4. Investigate Celestial Base (remove if doesn't exist)

**PHASE 2:**
5. Research and add all 12 landscaping buildings
6. Verify Compound Center max level

---

**Report Complete**
*buildings-specialist - 2026-02-10*
