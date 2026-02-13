# Cryptomines Online - Final Scope Definition

**Date:** 2026-02-07
**Goal:** Scope cut para terminar el juego en tiempo razonable

---

## ✅ FEATURES IN SCOPE

### **Already Implemented (Phase 0-2 + Quest System)**

- ✅ Buildings (22 types, multi-tile, 3D models)
- ✅ Resources (warehouse manual collect, production rates)
- ✅ Construction (queue, timers, upgrades)
- ✅ Ships (75 hulls: 30 frigates, 30 cruisers, 15 battleships)
- ✅ Modules (97 modules across 11 categories)
- ✅ Ship Factory (24 levels, 5 production slots)
- ✅ Ship Design (blueprints, module assignment)
- ✅ Fleets (3x3 grid, 27k max ships)
- ✅ Spacedock (12 levels, repair placeholder)
- ✅ Quest System (22 main, 12 side, 6 daily)
- ✅ Normal Instances (30 levels, placeholder combat)
- ✅ Auth (guest login with auto-create player)
- ✅ 3D UI (isometric grid, placement, camera)

### **CRITICAL (Must Have)**

#### 1. Research System ✅
- **Scope:** ALL 7 tech trees, ALL 111 techs
- **Backend:** 4 endpoints (list, start, cancel, speedup)
- **Auto-complete:** Timer completion
- **Effects:** Apply tech bonuses (production, unlock modules, etc.)
- **Trees:**
  - Ballistics Science (16 techs)
  - Ship Defense Science (19 techs)
  - Logistics Construction Science (11 techs)
  - Directional Science (18 techs)
  - Missile Science (18 techs)
  - Ship-Based Science (17 techs)
  - Planetary Defense Science (12 techs)

#### 2. Blueprint Research System ✅
- **Scope:** Upgrade hulls/modules from tier 1 → tier 2 → tier 3
- **Building:** Weapon Research Center (controls research slots)
- **Costs:** baseCost = 10000 × targetLevel (Metal/He3/Gold)
- **Time:** 3600 seconds × targetLevel (1 hour per level)
- **Backend:**
  - Endpoint: `POST /api/blueprints/{id}/research` - Start research (PARTIAL)
  - Worker: Auto-complete when `research_finish_at` reached (MISSING)
  - Service: Helper functions to parse `{base}_i/ii/iii` naming (MISSING)
  - Validation: Check tier unlocked in ship design (MISSING)
- **Frontend:**
  - BlueprintPanel: Research button + level stars (★★☆) (MISSING)
  - Ship Design: Filter hulls/modules by unlocked tiers (MISSING)
  - Research Display: Active research timer (MISSING)
- **Pattern:**
  - Hulls: `weikes_i` → `weikes_ii` → `weikes_iii`
  - Modules: same `name`, tier 1 → 2 → 3

#### 3. Resource Auto-Production ✅
- **Warehouse auto-accumulation:**
  - Add fields: `warehouse_metal`, `warehouse_he3`, `warehouse_gold`
  - Worker/cron updates warehouse every X minutes
  - Cap at `storage_capacity`
- **Collect endpoint:** Transfer warehouse → player resources
- **Review:** Current implementation antes de modificar

---

### **HIGH (Core Gameplay)**

#### 4. Commander System ✅
**IN:**
- Rarity: Common, Skill, Super (3 tiers)
- Stats: Accuracy, Dodge, Speed, Electron (todos)
- Gacha system usando Gold (no Mall Points)
- Star Rank modernizado (auto-merge dupes como gacha moderno)
- Compound Center building
- Assign commanders to fleets
- Command Center recruitment

**OUT:**
- Legendary tier
- Divine tier
- Commander Skills (Blue/Red/Green active skills)
- Gems
- Bionic Chips
- Wounded/Dead states (commanders immortales)
- Healing Card
- Revival Card

**Implementation:**
- Commanders inician como items en inventario
- Al usar → consume item → unlock en tabla commanders
- Max 60 commanders

#### 5. Combat System ✅
**Scope:** COMPLETO - 8-phase combat resolution
- Phase 1-8: Positioning → Targeting → Damage calc → Shields → Armor → HP → Casualties → Loot
- Ship type advantage (Frigate/Cruiser/Battleship rock-paper-scissors)
- Armor types vs Damage types (effectiveness matrix)
- Commander bonuses (accuracy, dodge, speed, electron)
- Combat reports (round-by-round detailed log)
- Instance combat (ya existe placeholder, implementar real)

#### 6. PvP Combat ✅
**Scope:** COMPLETO
- Attack neighbor planets
- Fleet travel time (usa SP - Space Points)
- Combat resolution (reusa combat system)
- Loot 20% resources (excluding warehouse)
- Radar building (incoming attack detection/warning)
- Defense fleets
- PvP combat reports
- Truce cards (protection shields)

#### 7. Space Station Defense Buildings ✅
**Buildings:** 5 defensive structures
- Meteor Star
- Particle Cannon
- Anti-Aircraft Gun
- Thor's Cannon
- Celestial Base

**Requirements:**
- 3D models (5 procedural buildings)
- Combat integration (auto-attack incoming fleets)
- Planetary Defense Science techs (12 techs)

#### 8. Recycling Plant ✅
**Scope:** Simple implementation
- Select ships to scrap
- Recover X% resources (ej: 70% of build cost)
- Delete ships from database

---

### **MEDIUM (Quality of Life)**

#### 9. Inventory System ✅
**Items totales:** ~17 types + blueprints + commander cards

**Resource Packs (8 items):**
- Gold Pack (+30k)
- Advanced Gold Pack (+100k)
- Primary/Junior/Senior Metal Pack (+50k/150k/300k)
- Primary/Junior/Senior He3 Pack (+50k/150k/300k)

**Resource Boosts (6 items):**
- Construction Card (+3 building slots, 72h)
- MVP Tool (+20% all production/build/repair, 7 días)
- Extra Tax / Adv (+30% / +100% Gold, 12h/24h)
- Metal Mining Boost (+30% Metal, 12h)
- He3 Mining Boost (+30% He3, 12h)

**Battle Items (3 items):**
- SP Card (+10 Space Points instant)
- Truce Card (12h protection)
- Adv Truce Card (72h protection)

**Consumables:**
- **Blueprints:** Item → Inventory → Use → Unlock en player_blueprints
- **Commander Cards:** Item → Inventory → Use → Unlock commander

**Implementation:**
- Tabla `player_items` (player_id, item_key, quantity)
- UI: Inventory panel
- Use logic per item type
- Fix blueprint/commander flow (actualmente directo, debe pasar por inventory)

#### 10. World Chat ✅
**Scope:** Basic world chat only
- Single world channel (no corps chat, no private messages)
- Send message
- View messages (recent 100)
- No moderation tools (basic profanity filter)

**OUT:**
- Friends system
- Mail system
- Private messages
- Corps chat (Corps system OUT)
- Loudspeaker item (no longer needed)

#### 11. Production Polish ✅
**IN:**
- Better error messages
- Loading states on all buttons
- Complete tooltips (techs, modules, buildings explain what they do)
- Smooth animations
- Tutorial tooltips (first-time hints)
- Sound effects (build complete, combat, etc.)

**OUT:**
- Mobile responsive polish (desktop only)

---

## ❌ FEATURES OUT OF SCOPE

### **Sistemas Completos Cortados:**

1. ❌ **Trading Center** - Player trading, auctions, marketplace
2. ✅ **Alliance Center / Corps System** - Corps, donations, RBPs, Galaxy Map (Phase 4, Feb 13)
3. ❌ **Decorative Buildings** - 12 aesthetic buildings (Casino, Fountain, etc.)
4. ❌ **Friends System** - Add friends, friend bonuses
5. ❌ **Mail System** - Private messaging
6. ❌ **League/Championship** - PvP tournaments, rankings
7. ❌ **Automated Tests** - Go tests, Vitest tests
8. ❌ **Mobile Responsive** - Desktop only

### **Items/Features Cortados:**

9. ❌ **Planet Transformation Packs** - Visual planet changes (6 items)
10. ❌ **Commander Enhancement Items** - Memory Chip, Merge Chip, Resetting Card, Sealing Card
11. ❌ **Healing/Revival Cards** - Commanders immortales
12. ❌ **Galaxy Transfer** - Move planet location
13. ❌ **Loudspeaker** - World chat item (chat es free)
14. ❌ **Corps Certificate** - Create corps
15. ❌ **Passport** - Extra Restricted Instance
16. ❌ **Constellation Pass / Scroll Chests** - Constellation Instances
17. ❌ **Treasure Boxes** - Instance rewards direct
18. ❌ **Legendary/Divine Commander Tiers** - Solo Common/Skill/Super
19. ❌ **Commander Skills** - Blue/Red/Green active skills
20. ❌ **Gems/Bionic Chips** - Commander enhancements

### **Quests Deferred (6 main quests OUT):**

21. ❌ **main_02:** Loud and Clear (world chat message) - chat es free now
22. ❌ **main_18:** Bigger Bags (inventory slots) - no bag limit
23. ❌ **main_19:** Resource Pack (use resource pack) - simplified
24. ❌ **main_20:** Growing Resources (grow comsats) - feature not implemented
25. ❌ **main_21:** Adding Friends (add friend) - friends OUT
26. ❌ **main_22:** Mail System (send mail) - mail OUT

**Total active main quests:** 22 (28 - 6 deferred)

---

## 📊 Feature Count Summary

### **Systems IN:**
- ✅ Research (7 trees, 111 techs)
- ✅ Blueprint Research (upgrade hulls/modules tier 1→2→3 at Weapon Research Center)
- ✅ Resource Auto-Production (warehouse accumulation)
- ✅ Commander System (gacha, 3 tiers, star rank, compound center)
- ✅ Combat System (8-phase full mechanics)
- ✅ PvP Combat (attack, radar, defense)
- ✅ Space Defense Buildings (5 buildings)
- ✅ Recycling Plant (scrap ships)
- ✅ Inventory System (17 item types + blueprints + commanders)
- ✅ World Chat (basic)
- ✅ Production Polish (minus mobile)

**Total: 11 major systems**

### **Systems OUT:**
- ❌ Trading Center
- ❌ Corps/Alliance
- ❌ Decorative Buildings
- ❌ Friends/Mail
- ❌ League/Championship
- ❌ Tests
- ❌ Mobile
- ❌ Commander enhancements (gems/chips/skills)

**Total: 8 major systems cut**

---

## 🎮 Final Game Scope

### **What Players Can Do:**

**Core Loop:**
1. Build/upgrade buildings on planet
2. Research 111 techs across 7 trees
3. Collect resources from warehouse
4. Design ships (hulls + modules)
5. Build ships in Ship Factory
6. Recruit commanders (gacha with Gold)
7. Merge duplicate commanders (star rank)
8. Assign commanders to fleets
9. Complete 30 Normal Instances (PvE combat)
10. Attack other players (PvP combat)
11. Defend with Space Station defenses
12. Complete 22 main quests + 12 side quests + 6 daily quests
13. Use inventory items (resource packs, boosts, truce cards)
14. Chat in world channel
15. Recycle old ships

**End Game:**
- Max out buildings (Civic Center Lv12)
- Complete all 111 techs
- Collect all blueprints
- Optimize ship designs
- Climb PvP rankings (informal, no league system)
- Complete all quests

---

## 🚀 Implementation Priority

### **Phase A: Foundation (Week 1-2)**
1. Research System backend (4 endpoints + auto-complete)
2. Resource Auto-Production (warehouse fields + worker)
3. Inventory System (tabla + UI + use logic)
4. Fix blueprint/commander flow (inventory → use → unlock)

### **Phase B: Combat (Week 3-4)**
5. Combat System 8-phase (combat engine core)
6. Instance combat integration (replace placeholder)
7. Combat reports UI

### **Phase C: Military (Week 5-6)**
8. Commander System (gacha + star rank + compound center)
9. Space Defense Buildings (5 models + combat integration)
10. PvP Combat (attack + radar + defense fleets)
11. Recycling Plant

### **Phase D: Polish (Week 7-8)**
12. World Chat (basic implementation)
13. Production Polish (errors, loading, tooltips, animations, sounds)
14. Quest system integration (auto-progress missing types)
15. Bug fixes + balancing

**Estimated Total: 8 weeks (2 months)**

---

## 📝 Notes

- **Scope creep:** If something not on this list is requested, default answer is OUT unless critical bug
- **Quest system:** Already implemented but needs integration with new systems (commanders, combat)
- **3D Models needed:** 5 Space Defense buildings + Command Center + Compound Center + Recycling Plant = 8 new models
- **Blueprint/Commander flow:** Major refactor needed (currently bypass inventory)
- **Combat engine:** Most complex piece, needs careful implementation + testing
- **Balancing:** Tech costs, commander stats, combat formulas need playtesting

---

**Last Updated:** 2026-02-07
**Status:** LOCKED - No changes without explicit approval
