# Cryptomines Online - Implementation Plan

**Date:** 2026-02-07
**Goal:** Roadmap definitivo para terminar el juego en 8 semanas
**Team:** All agents use Sonnet model

---

## 📅 Timeline Overview

**Total Duration:** 8 weeks (2 months)
**Weekly Commitment:** 5-6 days/week
**Target Completion:** April 2026

---

## 🎯 Phase A: Foundation Systems (Week 1-2)

**Duration:** 10-12 days
**Goal:** Implementar sistemas base que desbloquean el resto

### **A1: Research System (4-5 days)**

**Backend:**
- [ ] Migration: `research_progress` table (player_id, tech_type_id, level, is_researching, research_finish_at)
- [ ] Endpoint: `GET /api/research` - list all techs with player progress
- [ ] Endpoint: `GET /api/research/trees/{tree}` - get tech tree (ballistics, ship_defense, etc.)
- [ ] Endpoint: `POST /api/research/start` - start research (validate prereqs, resources, tech center level)
- [ ] Endpoint: `POST /api/research/cancel` - cancel active research (refund 50% resources)
- [ ] Endpoint: `POST /api/research/speedup` - spend vouchers/gold to reduce time
- [ ] Worker: Auto-complete research when timer expires
- [ ] Apply tech effects: production bonuses, module unlocks, build speed reductions

**Frontend:**
- [ ] ResearchPanel already exists, wire up API calls
- [ ] Show tech tree graph with locked/unlocked/researching states
- [ ] Start/cancel/speedup buttons
- [ ] Countdown timer for active research

**QA:**
- [ ] All 111 techs researchable
- [ ] Prereq validation works
- [ ] Effects apply correctly (test production boost, unlock modules)

---

### **A2: Resource Auto-Production (2-3 days)**

**Database:**
- [ ] Migration: Add `warehouse_metal`, `warehouse_he3`, `warehouse_gold` to `resources` table
- [ ] Migration: Set initial warehouse values = 0

**Backend:**
- [ ] Worker/cron: Every 5 minutes, update warehouse resources:
  ```sql
  UPDATE resources SET
    warehouse_metal = LEAST(warehouse_metal + (metal_per_hour * hours_elapsed / 12), storage_capacity),
    warehouse_he3 = LEAST(warehouse_he3 + (he3_per_hour * hours_elapsed / 12), storage_capacity),
    warehouse_gold = LEAST(warehouse_gold + (gold_per_hour * hours_elapsed / 12), storage_capacity),
    last_collect_at = NOW()
  WHERE ...
  ```
- [ ] Update `POST /api/planets/{id}/resources/collect`:
  - Transfer warehouse → player resources
  - Reset warehouse to 0
  - Return new balances

**Frontend:**
- [ ] Update ResourceHUD to show warehouse amounts (separate from player resources)
- [ ] Collect button tooltip: "Warehouse: X Metal, Y He3, Z Gold"
- [ ] Visual indicator when warehouse near cap

**QA:**
- [ ] Warehouse fills over time (check after 10min, 1hr, 12hr)
- [ ] Caps at storage_capacity
- [ ] Collect transfers correctly

---

### **A3: Inventory System (4-5 days)**

**Database:**
- [ ] Migration: `player_items` table (player_id, item_key TEXT, quantity INT, created_at, updated_at)
- [ ] Migration: `item_types` reference table (item_key, display_name, description, category, effects_json)
- [ ] Seed data: 17 item types (resource packs, boosts, cards, SP card)

**Backend:**
- [ ] Endpoint: `GET /api/inventory` - list player items
- [ ] Endpoint: `POST /api/inventory/use` - use item (consume, apply effect)
- [ ] Item use logic:
  - Resource Packs: instant resources
  - Boosts: apply temporary buff (MVP Tool, Mining Boosts, Extra Tax, Construction Card)
  - Truce Cards: set protection_until timestamp
  - SP Card: +10 SP to fleet

**Frontend:**
- [ ] InventoryPanel component (grid of items with quantities)
- [ ] Item tooltip (description, effects)
- [ ] Use button (with confirmation)
- [ ] Active buffs display (show remaining time for MVP Tool, etc.)

**QA:**
- [ ] Items consumable
- [ ] Effects apply correctly
- [ ] Boosts expire after duration

---

### **A4: Fix Blueprint/Commander Flow (2 days)**

**Problem:** Currently blueprints/commanders bypass inventory

**Backend:**
- [ ] Quest rewards: Give blueprint/commander items to inventory (not direct unlock)
- [ ] Endpoint: `POST /api/inventory/use-blueprint` - consume item → unlock in player_blueprints
- [ ] Endpoint: `POST /api/inventory/use-commander-card` - consume item → unlock commander

**Frontend:**
- [ ] Blueprint item in inventory → "Use" button → unlocks in Blueprints panel
- [ ] Commander card item in inventory → "Use" button → unlocks in (future) Commander panel

**QA:**
- [ ] Quest rewards give items (not direct unlocks)
- [ ] Using items unlocks blueprints/commanders

---

### **A5: Blueprint Research System (3-4 days)**

**Scope:** Upgrade hulls/modules from tier 1 → tier 2 → tier 3 at Weapon Research Center

**Backend:**
- [ ] Create `blueprint_service.go`:
  - `GetAvailableHullTiers()` - Parse `{base}_i/ii/iii` pattern
  - `GetAvailableModuleTiers()` - Filter by name + tier
  - `CanPlayerUseHull()` - Validate tier unlocked
  - `CanPlayerUseModule()` - Validate tier unlocked
- [ ] Create `blueprint_worker.go`:
  - `CompleteFinishedResearch()` - Check `research_finish_at`, update `research_level`
  - Call from main worker every 1 minute
- [ ] Update `ship_designs.go`:
  - Validate hull tier in `CreateShipDesign()`
  - Validate module tiers in `CreateShipDesign()`
  - Same validations in `UpdateShipDesign()`
- [ ] Optional endpoints:
  - `GET /api/blueprints/{id}/available-tiers` - List unlocked tiers
  - `POST /api/blueprints/research/cancel` - Cancel active research (refund 50%)
  - `POST /api/blueprints/research/speedup` - Spend gold/vouchers to accelerate

**Frontend:**
- [ ] Update `BlueprintPanel.tsx`:
  - Visual indicator of research_level (★★☆ stars)
  - "Research" button (only if research_level < 3)
  - Confirmation modal with costs + time
  - Call `POST /api/blueprints/{id}/research`
  - Display countdown timer for active research
- [ ] Update `ShipDesignPanel.tsx`:
  - Filter hull dropdown by unlocked tiers
  - Filter module lists by unlocked tiers
  - Show locked tiers in gray with lock icon + tooltip
  - Error handling for tier validation failures
- [ ] Create `ResearchProgressPanel.tsx` (optional):
  - Show active blueprint research
  - Progress bar with countdown
  - "Speed Up" button

**QA:**
- [ ] Research starts correctly (costs deducted, timer set)
- [ ] Worker completes research automatically when timer expires
- [ ] research_level increments (1→2→3)
- [ ] Ship design validates tier restrictions
- [ ] Cannot use tier 2 hull without research_level >= 2
- [ ] Cannot use tier 3 modules without research_level >= 3
- [ ] Quest progress updates on research complete

---

## 🎯 Phase B: Combat Engine (Week 3-4)

**Duration:** 10-14 days
**Goal:** Implementar combat 8-phase resolution

### **B1: Combat Engine Core (6-7 days)**

**Backend:**
- [ ] `combat_engine.go`: 8-phase combat resolution
  - Phase 1: Calculate effective stacks (base + commander bonus)
  - Phase 2: Ship type advantage (+5%/-5%)
  - Phase 3: Determine attack order (speed, random)
  - Phase 4: Calculate hit chance (accuracy vs dodge, clamp 5%-95%)
  - Phase 5: Calculate damage (weapon dmg * type advantage * armor effectiveness)
  - Phase 6: Apply damage (shields first, then structure)
  - Phase 7: Calculate casualties (destroyed ships)
  - Phase 8: Loot (PvP only, 20% resources)
- [ ] Armor vs Damage effectiveness matrix:
  ```
  Chrome: weak to Explosive, resists Kinetic
  Regen: weak to Heat, resists Explosive
  Nano: weak to Magnetic, resists Heat
  Neutralizing: weak to Kinetic, resists Magnetic
  ```
- [ ] Combat duration: Min 20 rounds, max 99 rounds
- [ ] Victory condition: One side ships = 0 OR 99 rounds elapsed (higher remaining HP wins)

**Data structures:**
- [ ] `CombatState` struct (attacker/defender fleets, round number, logs)
- [ ] `CombatRound` struct (attacks, damage dealt, ships destroyed)
- [ ] `CombatResult` struct (winner, loot, attacker/defender casualties)

**QA:**
- [ ] Unit tests for damage calculation
- [ ] Unit tests for armor effectiveness
- [ ] Unit tests for hit chance clamping
- [ ] Simulate 100 combats, verify no crashes

---

### **B2: Instance Combat Integration (2-3 days)**

**Backend:**
- [ ] Update `POST /api/instances/{id}/attempt`:
  - Load player fleet + instance NPC fleet
  - Call combat engine
  - Save combat report
  - Award rewards (blueprints, resources) if victory
  - Consume He3 from fleet
- [ ] `combat_reports` table: Store full round-by-round log

**Frontend:**
- [ ] InstancePanel: Show combat result (win/loss, casualties, loot)
- [ ] Combat report button → modal with round-by-round details

**QA:**
- [ ] Instance combat uses real combat engine (not placeholder)
- [ ] Rewards granted on victory
- [ ] Combat report saved and viewable

---

### **B3: Combat Reports UI (2 days)**

**Frontend:**
- [ ] CombatReportModal component
- [ ] Round-by-round display:
  - Round X: Attacker deals Y damage, Defender loses Z ships
- [ ] Summary: Total casualties, loot, duration
- [ ] Visual: Ship icons, damage numbers

**QA:**
- [ ] Reports readable and accurate
- [ ] Shows all phases

---

## 🎯 Phase C: Military Systems (Week 5-6)

**Duration:** 10-14 days
**Goal:** Commanders, Space Defenses, PvP, Recycling

### **C1: Commander System (5-6 days)**

**Database:**
- [ ] Migration: `commander_types` reference table (name, rarity, base_accuracy, base_dodge, base_speed, base_electron)
- [ ] Seed data: 20-30 commander types (Common, Skill, Super rarity)
- [ ] Migration: `commanders` table (player_id, commander_type_id, star_rank, accuracy, dodge, speed, electron)
- [ ] Migration: Update `fleets` table: add `commander_id` FK

**Backend:**
- [ ] Endpoint: `GET /api/commanders` - list player commanders
- [ ] Endpoint: `POST /api/commanders/recruit` - gacha draw (costs Gold)
  - Random rarity: 50% Common, 35% Skill, 15% Super
  - Random commander from rarity pool
  - If duplicate → give as item to inventory (for merging)
  - If new → unlock commander
- [ ] Endpoint: `POST /api/commanders/merge` - merge duplicate cards
  - Consume X duplicate cards → increase star_rank
  - Star rank increases stats + effective stack bonus
- [ ] Endpoint: `PUT /api/fleets/{id}/assign-commander` - assign commander to fleet
- [ ] Combat engine: Apply commander bonuses (accuracy, dodge, speed, electron, effective stack)

**Frontend:**
- [ ] CommandCenterPanel (gacha UI, cost display, draw button)
- [ ] CommandersListPanel (grid of owned commanders, stats)
- [ ] CompoundCenterPanel (merge UI, select dupes, merge button)
- [ ] FleetPanel: Commander dropdown (assign/remove)

**QA:**
- [ ] Gacha works (correct rarity rates)
- [ ] Duplicates go to inventory
- [ ] Merge increases star rank
- [ ] Commander bonuses apply in combat

---

### **C2: Space Station Defense Buildings (3-4 days)**

**Database:**
- [ ] Buildings already in seed: Meteor Star, Particle Cannon, Anti-Aircraft Gun, Thor's Cannon, Celestial Base
- [ ] Add building levels data if missing

**Backend:**
- [ ] Endpoint: `POST /api/planets/{id}/buildings` - construct defense buildings
- [ ] Defense stats: HP, attack, range (per building type/level)
- [ ] Combat engine: Defenses auto-attack incoming PvP fleets
  - Defenses attack before fleet reaches planet
  - Player can use Planetary weapons (Lander Module) to destroy defenses

**Frontend:**
- [ ] 3D models: 5 defense buildings (procedural generation)
- [ ] Place on grid (same as other buildings)
- [ ] Stats tooltip (HP, attack, range)

**QA:**
- [ ] Buildings placeable
- [ ] Defenses attack in PvP combat
- [ ] Planetary weapons can destroy defenses

---

### **C3: PvP Combat (4-5 days)**

**Database:**
- [ ] Migration: `planets` table already has `protection_until` field
- [ ] Migration: Add `pvp_attacks` table (attacker_id, defender_id, fleet_id, travel_finish_at, status)

**Backend:**
- [ ] Endpoint: `POST /api/pvp/attack` - send fleet to attack neighbor
  - Calculate travel time based on distance + fleet speed
  - Set fleet status = 'traveling'
  - Create pvp_attacks row with travel_finish_at
- [ ] Endpoint: `POST /api/pvp/recall` - recall traveling fleet
- [ ] Worker: Auto-resolve PvP combat when travel_finish_at reached
  - Load attacker fleet + defender planet defenses + defender fleets
  - Combat: Defenses attack first, then fleet vs fleet
  - Loot 20% resources (excluding warehouse) if attacker wins
  - Save combat report
  - Return fleet to homeworld
- [ ] Radar detection: If defender has Radar, show incoming attack warning

**Frontend:**
- [ ] Universe map (simplified): Show nearby planets
- [ ] Attack button on enemy planet
- [ ] Fleet travel animation (or timer)
- [ ] Incoming attacks panel (Radar view)
- [ ] PvP combat reports

**QA:**
- [ ] Attack sends fleet
- [ ] Combat resolves automatically
- [ ] Loot awarded on win
- [ ] Radar shows incoming attacks

---

### **C4: Recycling Plant (1-2 days)**

**Database:**
- [ ] Building already in seed

**Backend:**
- [ ] Endpoint: `POST /api/recycle` - scrap ships
  - Input: ship_design_id, quantity
  - Calculate resource recovery (70% of build cost)
  - Delete ships from database
  - Award resources to player

**Frontend:**
- [ ] RecyclingPlantPanel (list ship designs, quantity input, recycle button)
- [ ] Confirmation modal (show recovery amount)

**QA:**
- [ ] Ships deleted
- [ ] Resources awarded

---

## 🎯 Phase D: Polish & Integration (Week 7-8)

**Duration:** 10-14 days
**Goal:** World chat, polish, integration, balancing

### **D1: World Chat (2-3 days)**

**Database:**
- [ ] Migration: `world_chat_messages` table (player_id, message TEXT, created_at)

**Backend:**
- [ ] Endpoint: `GET /api/chat/world` - get recent 100 messages
- [ ] Endpoint: `POST /api/chat/world` - send message
- [ ] Profanity filter (basic)

**Frontend:**
- [ ] WorldChatPanel (message list, input field, send button)
- [ ] Auto-refresh every 10 seconds

**QA:**
- [ ] Messages send/receive
- [ ] Profanity filter works

---

### **D2: Production Polish (5-7 days)**

**Error Messages:**
- [ ] Replace generic "failed to X" with specific errors:
  - "Insufficient resources" (show what's needed)
  - "Prerequisites not met" (show what tech/building needed)
  - "Already researching" (show active research)

**Loading States:**
- [ ] All buttons show spinner when pending
- [ ] Disable buttons during API calls

**Tooltips:**
- [ ] All techs: Show prereqs, costs, effects
- [ ] All modules: Show stats, damage type, range
- [ ] All buildings: Show costs, benefits, requirements

**Animations:**
- [ ] Button hover effects
- [ ] Modal slide-in/fade-in
- [ ] Resource number count-up animation
- [ ] Construction progress bar smooth

**Tutorial Tooltips:**
- [ ] First-time hints (dismissable):
  - "Click here to collect resources"
  - "Research techs to unlock modules"
  - "Assign commanders to boost fleets"

**Sound Effects:**
- [ ] Building complete: "ding"
- [ ] Research complete: "chime"
- [ ] Combat victory: "fanfare"
- [ ] Button click: "click"
- [ ] Error: "error beep"

**QA:**
- [ ] All error messages helpful
- [ ] All buttons have loading states
- [ ] All tooltips complete
- [ ] Animations smooth
- [ ] Sounds play correctly

---

### **D3: Quest Integration (2-3 days)**

**Missing auto-progress types:**
- [ ] `recruit_commander` - wire to commander recruitment
- [ ] `replenish_ammo` - (skip, no ammo system)
- [ ] `build_defense` - wire to defense building construction
- [ ] `complete_instance` - wire to instance victory
- [ ] `donate_resources` - (skip, corps OUT)
- [ ] `login` - wire to player login
- [ ] `use_speedup` - wire to research/construction speedup

**QA:**
- [ ] All IN-SCOPE quest types auto-progress
- [ ] Rewards claimable

---

### **D4: Balancing & Bug Fixes (3-4 days)**

**Balancing:**
- [ ] Combat: Test fleet compositions, adjust if needed
- [ ] Economy: Resource production rates vs costs (should take ~1 week to max Civic Center)
- [ ] Commander gacha: Rarity rates feel fair

**Bug Fixes:**
- [ ] Full playthrou

gh (create account → max out buildings → complete quests → PvP)
- [ ] Fix all critical bugs
- [ ] Fix UI glitches

**Final QA:**
- [ ] All systems integrated
- [ ] No game-breaking bugs
- [ ] Game is fun and completable

---

## 🚀 Launch Checklist

**Pre-Launch:**
- [ ] All Phase A-D tasks complete
- [ ] Full playthrough test (0 to endgame)
- [ ] Performance test (1000 ships in fleet, combat simulation)
- [ ] Security review (SQL injection, XSS checks)

**Launch Day:**
- [ ] Deploy to production server
- [ ] Test production environment
- [ ] Monitor for errors
- [ ] Be ready for hotfixes

**Post-Launch:**
- [ ] Gather player feedback
- [ ] Fix critical bugs within 24h
- [ ] Plan Phase 2 features (based on scope cuts)

---

## 📊 Task Breakdown Summary

| Phase | Tasks | Days | Team Focus |
|-------|-------|------|------------|
| A | Foundation (Research, Blueprint Research, Auto-Production, Inventory, Fix Blueprints) | 15-16 | backend-dev, frontend-dev |
| B | Combat Engine (8-phase, Instances, Reports) | 10 | backend-dev, frontend-dev |
| C | Military (Commanders, Defenses, PvP, Recycling) | 13 | backend-dev, frontend-dev, 3d-assets |
| D | Polish (Chat, UX, Integration, Balancing) | 13 | frontend-dev, qa-agent |
| **Total** | **51-52 days** | **~8 weeks** | **All agents** |

---

## 🎯 Success Criteria

**Game is "done" when:**
1. ✅ All 111 techs researchable
2. ✅ All 75 hulls + 97 modules usable
3. ✅ Blueprint research works (upgrade tier 1→2→3 at Weapon Research Center)
4. ✅ Combat system works (PvE instances + PvP)
5. ✅ Commander system works (gacha + merge + assign)
6. ✅ 22 main quests + 12 side quests + 6 daily quests completable
7. ✅ Inventory system works (17 items usable)
8. ✅ World chat works
9. ✅ UI polished (errors, loading, tooltips, sounds)
10. ✅ Full playthrough possible (0 to endgame in ~2 weeks of play)
11. ✅ No critical bugs

---

**Document Status:** FINAL
**Last Updated:** 2026-02-07
**Ready to Start:** YES 🚀
