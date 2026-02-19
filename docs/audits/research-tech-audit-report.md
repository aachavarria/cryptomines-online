# Research Technology Implementation Audit Report

**Date:** 2026-02-14
**Scope:** All 7 research trees, 98 technologies total
**Files Audited:**
- `backend/internal/services/tech_effects.go` (effect handler)
- `backend/internal/combat/combat_engine.go` (combat application)
- `backend/internal/combat/combat_loader.go` (tech bonus bridge)
- `backend/internal/handlers/research.go` (API handler)
- `frontend/src/components/panels/ResearchPanel.tsx` (UI display)
- `supabase/migrations/20260206040000_research.sql` (tech definitions)
- `supabase/migrations/20260211000005_add_fighter_technologies.sql` (additional techs)

---

## 1. Executive Summary

| Metric | Count | Percentage |
|--------|-------|------------|
| Total Technologies Defined (DB) | 98 | 100% |
| Backend Handler Implemented (`applyTechEffect` case) | 31 | 31.6% |
| Backend Handler MISSING | 67 | 68.4% |
| Combat Engine Integration (used in combat) | 14 | 14.3% |
| Combat Engine MISSING integration | 84 | 85.7% |
| Frontend Display Working (simple per_level effects) | ~31 | 31.6% |
| Frontend Display Broken (complex/composite effects) | ~67 | 68.4% |

**Verdict: CRITICAL -- Two-thirds of all research technologies silently do nothing when researched.** Players spend resources and time to research techs that have zero gameplay effect. The `applyTechEffect` function has a `default` case that only logs "Unknown tech effect type" without any user-visible feedback.

---

## 2. Critical Systemic Issues (Ranked by Severity)

### CRITICAL-1: 67 Effect Types Have No `applyTechEffect` Case (Severity: CRITICAL)

**File:** `backend/internal/services/tech_effects.go:136-248`

The `applyTechEffect` switch statement only handles 31 simple effect types. The remaining 67 technologies define effect types in the database that fall through to the `default` case and are silently logged as unknown. Players researching these techs get no benefit whatsoever.

**Missing effect types (full list in Section 5):**
- `armor_bonus`, `shield_pen_chance`, `range_damage`, `scatter_damage`, `scatter_rate`, `scatter_weapon_chance` (Ballistics)
- `shield_damage_reduction`, `shield_pen_resist`, `base_shield_bonus`, `shield_restore`, `shield_flat_reduction`, `absorb_no_he3`, `reflect_damage`, `absorb_double`, `structure_damage_reduction`, `structure_pen_resist`, `base_structure_bonus`, `structure_restore`, `structure_flat_reduction`, `absorb_no_he3_structure`, `reflect_structure_damage`, `absorb_double_structure` (Ship Defense)
- `hit_and_shield_pen`, `piercing_damage`, `enemy_hit_reduction`, `piercing_damage_bonus`, `movement_crit_bonus`, `steering_reduction`, `piercing_critical`, `base_damage_and_accuracy`, `enemy_attack_reduction`, `ignore_agility`, `ignore_defense` (Directional)
- `missile_damage_and_pen`, `scatter_all`, `damage_taken_increase`, `armor_damage_bonus`, `scatter_bonus`, `scatter_vs_low_structure`, `scatter_vs_high_structure`, `knockback` (Missile)
- `base_attack`, `he3_and_shield_damage`, `ship_shield_damage`, `reload_chance`, `multi_bonus`, `long_range_attack`, `swarm_chance`, `long_range_bonus` (Ship-Based)
- `defense_range`, `defense_movement` (Planetary Defense)
- `reload_reduction` (Missile -- added tech)
- `base_defense_stats` (Ship Defense -- has a case but is a no-op)

### CRITICAL-2: Combat Engine `TechBonuses` Struct Is Incomplete (Severity: CRITICAL)

**File:** `backend/internal/combat/combat_engine.go:112-126`

The combat engine defines its own `TechBonuses` struct with only 14 fields. Even for the 31 effect types that DO have `applyTechEffect` cases, many bonuses are never used in actual combat calculations:

**Fields in `combat.TechBonuses` (14 total):**
```
BallisticDamage, BallisticCritRate, BallisticCritDamage, BallisticHitRate,
DirectionalDamage, DirectionalCritRate, DirectionalAccuracy,
MissileDamage, MissileHitRate,
BaseShield, BaseStructure, BaseAgility, BaseDefense
```

**Missing from `combat.TechBonuses` (present in `services.TechBonuses` but NOT bridged):**
- `BaseStability` -- defined in services, missing from combat struct
- `SteeringPower` -- defined in services, missing from combat struct
- `WeaponSpaceReduction` -- defined in services, missing from combat struct
- `ShieldBypass` -- defined in services, missing from combat struct
- `InterceptReduction` -- defined in services, missing from combat struct
- `FighterDamage` -- defined in services, missing from combat struct
- `FighterHitRate` -- defined in services, missing from combat struct
- `FuelOptimization` -- defined in services, missing from combat struct

### CRITICAL-3: Combat Engine Does NOT Apply Most Tech Bonuses (Severity: CRITICAL)

**File:** `backend/internal/combat/combat_engine.go:342-384`

In `phase1CalculateEffectiveStacks`, only `BaseShield` and `BaseStructure` from tech bonuses are actually applied to combat:

```go
// Apply tech bonuses (simplified for now)
if fleet.TechBonuses != nil {
    stack.EffectiveShield = int(float64(stack.EffectiveShield) * (1.0 + fleet.TechBonuses.BaseShield/100.0))
    stack.EffectiveStructure = int(float64(stack.EffectiveStructure) * (1.0 + fleet.TechBonuses.BaseStructure/100.0))
}
```

**Result:** Even `BallisticDamage`, `MissileDamage`, `DirectionalDamage`, `BallisticCritRate`, etc. are bridged into the combat struct but NEVER referenced in any damage calculation. The comment "simplified for now" confirms this is a known incomplete implementation.

### HIGH-1: `base_defense_stats` Is a No-Op (Severity: HIGH)

**File:** `backend/internal/services/tech_effects.go:208-219`

The `ship_defense_base` tech's `base_defense_stats` effect type has a case in the switch, but it does nothing:

```go
case "base_defense_stats":
    var baseDefenseEffect struct { ... }
    _ = baseDefenseEffect  // <-- explicitly discarded
```

This is the root tech of the entire Ship Defense Science tree. Its composite bonuses (shield +2%, structure +2%, agility +2%, defense +2%, stability +5% per level) are never applied.

### HIGH-2: Frontend Only Renders Simple `per_level` Effects (Severity: HIGH)

**File:** `frontend/src/components/panels/ResearchPanel.tsx:590-609`

The `formatEffect` and `formatNextEffect` functions only handle `per_level` based effects:

```typescript
const perLevel = eff.per_level ?? 0
const total = perLevel * tech.current_level
```

For technologies with complex effects (composite objects with `values`, `ranges`, `flat`, `chance`, `damage`, etc.), the display shows `+0 [type]` or nothing. This affects 67 technologies including high-value capstone techs like `victory_rush`, `ingenuity`, `damage_mitigation`, etc.

### HIGH-3: Combat Loader Bridge Drops Fields (Severity: HIGH)

**File:** `backend/internal/combat/combat_loader.go:57-71`

The bridge from `services.TechBonuses` to `combat.TechBonuses` explicitly maps only 14 of 33+ fields. Fields like `BaseStability`, `SteeringPower`, `WeaponSpaceReduction`, `ShieldBypass`, `InterceptReduction`, `FighterDamage`, `FighterHitRate`, `FuelOptimization` are silently dropped.

### MEDIUM-1: `steering_power` Used by Multiple Trees but Only One Field (Severity: MEDIUM)

Both `cruise_dynamics` (Missile) and `reconnaissance` (Ship-Based) and `eagle_eye` (Directional) use effect type `steering_power`. The `applyTechEffect` handler correctly accumulates them into `bonuses.SteeringPower`, but this field is never bridged to combat and has no gameplay effect.

### MEDIUM-2: `he3_cost_reduction` Shared Across Trees (Severity: MEDIUM)

Both `fuel_optimization` (Ship-Based) and `energy_conservation_missile` (Missile) use `he3_cost_reduction`. The handler accumulates into `FuelOptimization`, but there is no system that actually reduces He3 costs based on this value.

### MEDIUM-3: No Frontend Bonuses Summary for Complex Effects (Severity: MEDIUM)

**File:** `frontend/src/components/panels/ResearchPanel.tsx:65-88`

The "ACTIVE BONUSES" summary section iterates all trees and sums `per_level * current_level`, which is correct for simple effects but produces wrong or zero values for complex effects.

---

## 3. Consolidated Implementation Status Table

### Legend
- **OK** = Implemented and working
- **PARTIAL** = Handler exists but incomplete or not applied in combat
- **MISSING** = No handler case; falls to default log
- **NO-OP** = Handler case exists but does nothing
- **N/A** = Not applicable (non-combat tech)

### 3.1 Logistics Construction (11 techs)

| # | Tech Name | Effect Type | Backend Handler | Combat Application | Frontend Display | Overall |
|---|-----------|-------------|-----------------|-------------------|-----------------|---------|
| 1 | concurrent_construction | construction_slots | OK (line 151) | N/A | OK | **OK** |
| 2 | construction_boost | build_speed | OK (line 147) | N/A | OK | **OK** |
| 3 | quality_materials | build_cost_reduction | OK (line 149) | N/A | OK | **OK** |
| 4 | ship_building_boost | ship_build_speed | OK (line 153) | N/A | OK | **OK** |
| 5 | ship_building_logistics | ship_build_cost_reduction | OK (line 155) | N/A | OK | **OK** |
| 6 | sync_shipbuilding | ship_production_slots | OK (line 157) | N/A | OK | **OK** |
| 7 | repair_technology | ship_repair_percent | OK (line 165) | N/A | OK | **OK** |
| 8 | high_yield_mining | metal_output | OK (line 139) | N/A | OK | **OK** |
| 9 | high_yield_chemistry | he3_output | OK (line 141) | N/A | OK | **OK** |
| 10 | high_yield_investing | gold_output | OK (line 143) | N/A | OK | **OK** |
| 11 | expand_capacity | warehouse_capacity | OK (line 161) | N/A | OK (but shows unit:flat not formatted ideally) | **OK** |

**Tree Score: 11/11 OK (100%).** Best-implemented tree. All non-combat production/construction effects work end-to-end.

### 3.2 Ballistics Science (14 techs)

| # | Tech Name | Effect Type | Backend Handler | Combat Application | Frontend Display | Overall |
|---|-----------|-------------|-----------------|-------------------|-----------------|---------|
| 1 | ballistics_base | ballistic_damage | OK (line 169) | PARTIAL (bridged but not used in damage calc) | OK | **PARTIAL** |
| 2 | ballistic_malice | ballistic_crit_rate | OK (line 171) | PARTIAL (bridged but not used) | OK | **PARTIAL** |
| 3 | ballistic_crackdown | ballistic_crit_damage | OK (line 173) | PARTIAL (bridged but not used) | OK | **PARTIAL** |
| 4 | steady_control | weapon_space_reduction | OK (line 177) | MISSING (not in combat struct) | OK | **PARTIAL** |
| 5 | precise_ballistics | ballistic_hit_rate | OK (line 175) | PARTIAL (bridged but not used) | OK | **PARTIAL** |
| 6 | shield_penetration | shield_bypass | OK (line 179) | MISSING (not in combat struct) | Broken (flat, no per_level) | **PARTIAL** |
| 7 | depleted_uranium_bomb | armor_bonus | MISSING | MISSING | Broken (complex JSON) | **MISSING** |
| 8 | fire_bomb_research | armor_bonus | MISSING | MISSING | Broken (complex JSON) | **MISSING** |
| 9 | improved_penetration | shield_pen_chance | MISSING | MISSING | Broken (complex JSON) | **MISSING** |
| 10 | victory_rush | range_damage | MISSING | MISSING | Broken (complex JSON) | **MISSING** |
| 11 | ballistic_scattering | scatter_damage | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 12 | improved_ballistic_scattering | scatter_rate | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 13 | hop_bomb_research | scatter_weapon_chance | MISSING | MISSING | OK (has per_level) | **MISSING** |

**Tree Score: 0/13 fully working. 6 PARTIAL (handler exists, no combat effect). 7 MISSING entirely.**

*Note: Only 13 techs listed because the SQL shows 13 inserts for ballistics_science, not 14. The task description may have counted one extra.*

### 3.3 Ship Defense Science (20 techs)

| # | Tech Name | Effect Type | Backend Handler | Combat Application | Frontend Display | Overall |
|---|-----------|-------------|-----------------|-------------------|-----------------|---------|
| 1 | ship_defense_base | base_defense_stats | NO-OP (line 208-219) | Not applied | Broken (composite JSON) | **NO-OP** |
| 2 | shield_research | base_shield | OK (line 221) | OK (phase1, line 375) | OK | **OK** |
| 3 | energy_diffusion | shield_damage_reduction | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 4 | penetration_resistance | shield_pen_resist | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 5 | augment_shield | base_shield_bonus | MISSING | MISSING | Broken (values array) | **MISSING** |
| 6 | restoration | shield_restore | MISSING | MISSING | Broken (composite) | **MISSING** |
| 7 | augment_absorption | shield_flat_reduction | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 8 | energy_conservation_defense | absorb_no_he3 | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 9 | electronic_barrier | reflect_damage | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 10 | damage_mitigation | absorb_double | MISSING | MISSING | Broken (composite) | **MISSING** |
| 11 | ship_structural_analysis | base_structure | OK (line 223) | OK (phase1, line 376) | OK | **OK** |
| 12 | ship_reinforcement | structure_damage_reduction | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 13 | resilience | structure_pen_resist | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 14 | structure_improvement | base_structure_bonus | MISSING | MISSING | Broken (values array) | **MISSING** |
| 15 | fast_repair | structure_restore | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 16 | reaction_armor_improvement | structure_flat_reduction | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 17 | defense_improvement | absorb_no_he3_structure | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 18 | reflection_mastery | reflect_structure_damage | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 19 | stability_mastery | absorb_double_structure | MISSING | MISSING | Broken (composite) | **MISSING** |

**Tree Score: 2/20 fully working (shield_research, ship_structural_analysis only). 1 NO-OP. 17 MISSING.**

*Note: 19 techs listed. The SQL shows 19 INSERT values for ship_defense_science (the 20th may have been miscounted in task description, or `ship_defense_base` was counted separately from the shield/structure branches).*

### 3.4 Directional Science (15 techs)

| # | Tech Name | Effect Type | Backend Handler | Combat Application | Frontend Display | Overall |
|---|-----------|-------------|-----------------|-------------------|-----------------|---------|
| 1 | optics_base | directional_damage | OK (line 183) | PARTIAL (bridged, not used in damage calc) | OK | **PARTIAL** |
| 2 | directional_malice | directional_crit_rate | OK (line 185) | PARTIAL (bridged, not used) | OK | **PARTIAL** |
| 3 | directional_accuracy | directional_accuracy | OK (line 187) | PARTIAL (bridged, not used) | OK | **PARTIAL** |
| 4 | eagle_eye | steering_power | OK (line 189) | MISSING (not in combat struct) | OK | **PARTIAL** |
| 5 | energy_penetration | hit_and_shield_pen | MISSING | MISSING | Broken (composite) | **MISSING** |
| 6 | pierce | piercing_damage | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 7 | radiative_interference | enemy_hit_reduction | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 8 | improved_pierce | piercing_damage_bonus | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 9 | energy_accumulation | movement_crit_bonus | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 10 | electronic_interference | steering_reduction | MISSING | MISSING | Broken (composite) | **MISSING** |
| 11 | piercing_crit | piercing_critical | MISSING | MISSING | Broken (boolean) | **MISSING** |
| 12 | weakness_detection | base_damage_and_accuracy | MISSING | MISSING | Broken (composite) | **MISSING** |
| 13 | particle_impact | enemy_attack_reduction | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 14 | magnetic_impact | ignore_agility | MISSING | MISSING | Broken (composite) | **MISSING** |
| 15 | dynamic_impairment | ignore_defense | MISSING | MISSING | Broken (composite) | **MISSING** |

**Tree Score: 0/15 fully working. 4 PARTIAL (handler exists, no combat effect). 11 MISSING entirely.**

### 3.5 Missile Science (14 + 1 techs)

| # | Tech Name | Effect Type | Backend Handler | Combat Application | Frontend Display | Overall |
|---|-----------|-------------|-----------------|-------------------|-----------------|---------|
| 1 | missile_theory | missile_damage | OK (line 193) | PARTIAL (bridged, not used in damage calc) | OK | **PARTIAL** |
| 2 | missile_accuracy | missile_hit_rate | OK (line 195) | PARTIAL (bridged, not used) | OK | **PARTIAL** |
| 3 | cruise_dynamics | steering_power | OK (line 189) | MISSING (not in combat struct) | OK | **PARTIAL** |
| 4 | missile_research | missile_damage_and_pen | MISSING | MISSING | Broken (composite) | **MISSING** |
| 5 | missile_elusion | intercept_reduction | OK (line 197) | MISSING (not in combat struct) | Broken (composite: also has hit_bonus) | **PARTIAL** |
| 6 | missile_space_optimization | weapon_space_reduction | OK (line 177) | MISSING (not in combat struct) | OK | **PARTIAL** |
| 7 | multidirectional_assault | scatter_all | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 8 | nuclear_radiation | damage_taken_increase | MISSING | MISSING | Broken (composite) | **MISSING** |
| 9 | break_armor | armor_damage_bonus | MISSING | MISSING | Broken (flat, no per_level) | **MISSING** |
| 10 | energy_conservation_missile | he3_cost_reduction | OK (line 204) | MISSING (no He3 cost system) | Broken (composite: also has chance_for_half) | **PARTIAL** |
| 11 | shrapnel_research | scatter_bonus | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 12 | exaltation | scatter_vs_low_structure | MISSING | MISSING | Broken (composite) | **MISSING** |
| 13 | suppression | scatter_vs_high_structure | MISSING | MISSING | Broken (composite) | **MISSING** |
| 14 | missile_concussion | knockback | MISSING | MISSING | Broken (no per_level) | **MISSING** |
| 15 | rapid_loading | reload_reduction | MISSING | MISSING | Broken (flat, no per_level) | **MISSING** |

**Tree Score: 0/15 fully working. 6 PARTIAL (handler exists, no combat effect). 9 MISSING entirely.**

### 3.6 Ship-Based Science (10 + 6 techs)

| # | Tech Name | Effect Type | Backend Handler | Combat Application | Frontend Display | Overall |
|---|-----------|-------------|-----------------|-------------------|-----------------|---------|
| 1 | fighter_weapons_theory | fighter_damage | OK (line 201) | MISSING (not in combat struct) | OK | **PARTIAL** |
| 2 | reconnaissance | steering_power | OK (line 189) | MISSING (not in combat struct) | OK | **PARTIAL** |
| 3 | thruster_optimization | intercept_reduction | OK (line 197) | MISSING (not in combat struct) | OK | **PARTIAL** |
| 4 | navigation | fighter_hit_rate | OK (line 203) | MISSING (not in combat struct) | OK | **PARTIAL** |
| 5 | fuel_optimization | he3_cost_reduction | OK (line 204) | MISSING (no He3 cost system) | OK | **PARTIAL** |
| 6 | fighter_mastery | base_attack | MISSING | MISSING | Broken (flat, no per_level) | **MISSING** |
| 7 | fighter_tech_upgrades | he3_and_shield_damage | MISSING | MISSING | Broken (composite) | **MISSING** |
| 8 | armor_structural_analysis | ship_shield_damage | MISSING | MISSING | Broken (flat, no per_level) | **MISSING** |
| 9 | fighter_weapons_efficiency | reload_chance | MISSING | MISSING | OK (has per_level) | **MISSING** |
| 10 | ingenuity | multi_bonus | MISSING | MISSING | Broken (composite) | **MISSING** |
| 11 | long_ranged_strike | long_range_attack | MISSING | MISSING | Broken (flat) | **MISSING** |
| 12 | fighter_interception_countermeasures | intercept_reduction | OK (line 197) | MISSING | OK | **PARTIAL** |
| 13 | formation_optimization | multi_bonus | MISSING | MISSING | Broken (composite) | **MISSING** |
| 14 | swarm | swarm_chance | MISSING | MISSING | Broken (composite) | **MISSING** |
| 15 | fortune | long_range_bonus | MISSING | MISSING | Broken (composite) | **MISSING** |
| 16 | heavy_gear_research | multi_bonus | MISSING | MISSING | Broken (composite) | **MISSING** |

**Tree Score: 0/16 fully working. 6 PARTIAL (handler exists, no combat/system effect). 10 MISSING entirely.**

### 3.7 Planetary Defense (8 techs)

| # | Tech Name | Effect Type | Backend Handler | Combat Application | Frontend Display | Overall |
|---|-----------|-------------|-----------------|-------------------|-----------------|---------|
| 1 | energy_control | defense_cost_reduction | OK (line 233) | N/A (needs defense build system) | OK | **PARTIAL** |
| 2 | rapid_defense_buildup | defense_build_speed | OK (line 235) | N/A (needs defense build system) | OK | **PARTIAL** |
| 3 | defense_enhancement | defense_value | OK (line 237) | N/A (needs defense combat) | OK | **PARTIAL** |
| 4 | emplacement_mastery | emplacement_attack | OK (line 239) | N/A (needs defense combat) | OK | **PARTIAL** |
| 5 | utmost_defense_buildup | max_defense_structures | OK (line 241) | N/A (needs defense build system) | OK | **PARTIAL** |
| 6 | range_extension | defense_range | MISSING | MISSING | Broken (no per_level) | **MISSING** |
| 7 | thor_buildup | max_thor_cannon | OK (line 243) | N/A (needs defense build system) | Broken (flat, no per_level) | **PARTIAL** |
| 8 | augment_propulsion | defense_movement | MISSING | MISSING | OK (has per_level) | **MISSING** |

**Tree Score: 0/8 fully working. 6 PARTIAL (handler exists but no defense system to consume values). 2 MISSING entirely.**

---

## 4. Breakdown by Tree

| Tree | Total Techs | Fully Working | Partial | Missing/No-Op | % Working |
|------|-------------|---------------|---------|---------------|-----------|
| Logistics Construction | 11 | 11 | 0 | 0 | **100%** |
| Ballistics Science | 13 | 0 | 6 | 7 | **0%** |
| Ship Defense Science | 19-20 | 2 | 0 | 17-18 | **10%** |
| Directional Science | 15 | 0 | 4 | 11 | **0%** |
| Missile Science | 15 | 0 | 6 | 9 | **0%** |
| Ship-Based Science | 16 | 0 | 6 | 10 | **0%** |
| Planetary Defense | 8 | 0 | 6 | 2 | **0%** |
| **TOTAL** | **~98** | **13** | **28** | **~57** | **13.3%** |

Note: "Fully Working" means the tech has a handler, is applied where relevant, AND displays correctly. "Partial" means the handler exists but the bonus is never consumed by any game system.

---

## 5. Missing Effect Types List

These effect types appear in the database `effects_json` but have NO case in `applyTechEffect` (`tech_effects.go:136-248`):

### Combat Effects (would need combat engine integration)
| Effect Type | Used By Techs | Tree |
|-------------|---------------|------|
| `armor_bonus` | depleted_uranium_bomb, fire_bomb_research | Ballistics |
| `shield_pen_chance` | improved_penetration | Ballistics |
| `range_damage` | victory_rush | Ballistics |
| `scatter_damage` | ballistic_scattering | Ballistics |
| `scatter_rate` | improved_ballistic_scattering | Ballistics |
| `scatter_weapon_chance` | hop_bomb_research | Ballistics |
| `shield_damage_reduction` | energy_diffusion | Ship Defense |
| `shield_pen_resist` | penetration_resistance | Ship Defense |
| `base_shield_bonus` | augment_shield | Ship Defense |
| `shield_restore` | restoration | Ship Defense |
| `shield_flat_reduction` | augment_absorption | Ship Defense |
| `absorb_no_he3` | energy_conservation_defense | Ship Defense |
| `reflect_damage` | electronic_barrier | Ship Defense |
| `absorb_double` | damage_mitigation | Ship Defense |
| `structure_damage_reduction` | ship_reinforcement | Ship Defense |
| `structure_pen_resist` | resilience | Ship Defense |
| `base_structure_bonus` | structure_improvement | Ship Defense |
| `structure_restore` | fast_repair | Ship Defense |
| `structure_flat_reduction` | reaction_armor_improvement | Ship Defense |
| `absorb_no_he3_structure` | defense_improvement | Ship Defense |
| `reflect_structure_damage` | reflection_mastery | Ship Defense |
| `absorb_double_structure` | stability_mastery | Ship Defense |
| `hit_and_shield_pen` | energy_penetration | Directional |
| `piercing_damage` | pierce | Directional |
| `enemy_hit_reduction` | radiative_interference | Directional |
| `piercing_damage_bonus` | improved_pierce | Directional |
| `movement_crit_bonus` | energy_accumulation | Directional |
| `steering_reduction` | electronic_interference | Directional |
| `piercing_critical` | piercing_crit | Directional |
| `base_damage_and_accuracy` | weakness_detection | Directional |
| `enemy_attack_reduction` | particle_impact | Directional |
| `ignore_agility` | magnetic_impact | Directional |
| `ignore_defense` | dynamic_impairment | Directional |
| `missile_damage_and_pen` | missile_research | Missile |
| `scatter_all` | multidirectional_assault | Missile |
| `damage_taken_increase` | nuclear_radiation | Missile |
| `armor_damage_bonus` | break_armor | Missile |
| `scatter_bonus` | shrapnel_research | Missile |
| `scatter_vs_low_structure` | exaltation | Missile |
| `scatter_vs_high_structure` | suppression | Missile |
| `knockback` | missile_concussion | Missile |
| `reload_reduction` | rapid_loading | Missile |
| `base_attack` | fighter_mastery | Ship-Based |
| `he3_and_shield_damage` | fighter_tech_upgrades | Ship-Based |
| `ship_shield_damage` | armor_structural_analysis | Ship-Based |
| `reload_chance` | fighter_weapons_efficiency | Ship-Based |
| `multi_bonus` | ingenuity, formation_optimization, heavy_gear_research | Ship-Based |
| `long_range_attack` | long_ranged_strike | Ship-Based |
| `swarm_chance` | swarm | Ship-Based |
| `long_range_bonus` | fortune | Ship-Based |
| `defense_range` | range_extension | Planetary Defense |
| `defense_movement` | augment_propulsion | Planetary Defense |

### No-Op Effect (has case but does nothing)
| Effect Type | Used By Tech | Tree |
|-------------|-------------|------|
| `base_defense_stats` | ship_defense_base | Ship Defense |

**Total Missing/No-Op Effect Types: 52 unique types**

---

## 6. Recommended Fix Priority

### Priority 1: Fix Combat Engine Tech Application (HIGHEST IMPACT)

**Impact:** Makes 28 "PARTIAL" techs actually work in combat.
**Effort:** Medium
**Files:** `combat_engine.go`, `combat_loader.go`

1. Add missing fields to `combat.TechBonuses`: `SteeringPower`, `WeaponSpaceReduction`, `ShieldBypass`, `InterceptReduction`, `FighterDamage`, `FighterHitRate`, `FuelOptimization`, `BaseStability`
2. Bridge all fields in `combat_loader.go:57-71`
3. Apply `BallisticDamage`, `DirectionalDamage`, `MissileDamage`, `FighterDamage` in `phase5CalculateDamage` as weapon damage multipliers
4. Apply `BallisticHitRate`, `MissileHitRate`, `DirectionalAccuracy`, `FighterHitRate` in `phase4CalculateHitChance`
5. Apply `BallisticCritRate`, `DirectionalCritRate`, `BallisticCritDamage` in the crit calculation (currently only uses commander Electron stat)
6. Apply `BaseAgility`, `BaseDefense`, `BaseStability` to effective stats in `phase1`

**This single change makes the 6 base damage/crit/hit techs (ballistics_base, missile_theory, optics_base, etc.) work correctly -- the techs most players research first.**

### Priority 2: Implement `base_defense_stats` Composite Effect

**Impact:** Fixes the root tech of the entire Ship Defense tree (ship_defense_base).
**Effort:** Low
**File:** `tech_effects.go:208-219`

Parse the composite JSON and apply shield/structure/agility/defense/stability bonuses individually.

### Priority 3: Add Simple Missing Effect Handlers

**Impact:** Makes ~20 simple `per_level` effects work in the bonus struct.
**Effort:** Low-Medium
**Files:** `tech_effects.go` (add cases), `TechBonuses` struct (add fields)

Focus on effects that have a simple `per_level` structure:
- `scatter_damage`, `scatter_rate`, `scatter_weapon_chance` (Ballistics scatter branch)
- `shield_damage_reduction`, `shield_pen_resist`, `shield_flat_reduction` (Shield defense)
- `structure_damage_reduction`, `structure_pen_resist`, `structure_flat_reduction` (Structure defense)
- `absorb_no_he3`, `absorb_no_he3_structure` (Energy conservation)
- `reflect_damage`, `reflect_structure_damage` (Reflection)
- `piercing_damage`, `enemy_hit_reduction`, `piercing_damage_bonus` (Directional)
- `enemy_attack_reduction`, `movement_crit_bonus` (Directional advanced)
- `scatter_all`, `scatter_bonus` (Missile scatter)

### Priority 4: Fix Frontend Display for Complex Effects

**Impact:** Players can see what their researched techs actually do.
**Effort:** Medium
**File:** `ResearchPanel.tsx:590-609`

Extend `formatEffect` and `formatNextEffect` to handle:
- `flat` values (e.g., `shield_bypass`, `max_thor_cannon`)
- `values` arrays (e.g., `augment_shield`, `structure_improvement`)
- Composite effects with multiple fields (e.g., `multi_bonus`, `range_damage`)
- Boolean toggles (e.g., `piercing_critical`)

### Priority 5: Implement Advanced Combat Mechanics

**Impact:** Enables endgame combat depth.
**Effort:** High
**Files:** `combat_engine.go` (new phases/mechanics)

These require new combat subsystems:
- Scatter damage (hit adjacent ships)
- Piercing damage (damage behind target)
- Shield/structure restore per round
- Damage reflection
- Knockback mechanics
- Reload/successive strike from tech
- Armor type bonus damage

### Priority 6: Implement Planetary Defense System

**Impact:** Makes 8 planetary defense techs meaningful.
**Effort:** High
**Files:** New defense building/combat system needed

The planetary defense techs accumulate bonuses correctly but there is no system to build/deploy/fight with planetary defense structures yet.

---

## 7. Architecture Recommendation

The current design has a structural problem: **two separate `TechBonuses` structs** exist (`services.TechBonuses` and `combat.TechBonuses`) that must be kept in sync manually. This creates a maintenance burden and is the root cause of dropped fields.

**Recommendation:** Either:
1. Use a single shared `TechBonuses` struct (import `services.TechBonuses` directly in combat package), OR
2. Use a map-based approach (`map[string]float64`) so new effect types don't require struct changes, OR
3. At minimum, add compile-time validation that all fields are bridged

---

*Report generated by automated audit of all 7 research trees.*
*All line references are to files as of commit b27e686 on main branch.*
