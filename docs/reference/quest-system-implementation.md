# Quest System Implementation Summary

**Date:** 2026-02-06
**Status:** ✅ COMPLETED (implemented by user)
**Migration:** `20260206030000_quests.sql`

---

## Overview

Complete quest system with main quest chain, side quests, and daily quests. Integrated with all game actions (building, research, ships, blueprints). Auto-progress tracking via `quest_service.UpdateQuestProgress()`.

---

## Database Schema

### Tables (3)

1. **quest_types** (reference data)
   - Fields: id, quest_key, category (main/side/daily), display_name, description, requirement_type, requirement_target, requirement_value, chain_order, prerequisite_quest_id, reward_metal, reward_he3, reward_gold, reward_item_json, phase, is_active
   - Seed: 46 rows (28 main + 12 side + 6 daily)

2. **player_quests** (progress tracking)
   - Fields: id, player_id, quest_type_id, status (locked/available/in_progress/completed/claimed), progress_value, started_at, completed_at, claimed_at, updated_at
   - Auto-initialized on first ListQuests call

3. **daily_quest_progress** (daily points)
   - Fields: id, player_id, quest_date, daily_points, quests_completed_json, tier_rewards_claimed_json, created_at, updated_at
   - Auto-created daily per player

---

## Quest Content

### Main Quest Chain (28 quests)

**Phase 1 (9 active):** 01, 03, 04, 05, 06, 07, 08, 09, 23, 24, 25, 26, 27
- Focus: building progression, resource production, tech research
- Blueprint rewards: Super Transmission Engine, Estrella, Ship Reinforcement Facility

**Phase 2 (9 active):** 10, 11, 12, 13, 14, 15, 16, 17, 28
- Focus: military systems (Ship Factory, Command Center, commanders, ships, fleets)
- Blueprint rewards: Typhoon, Energy Shield Booster, Anti-Aircraft Cannon, Starlight Missile Pod

**Deferred (6 inactive):** 02, 18, 19, 20, 21, 22
- Require social/chat systems (world channel, friends, mail, inventory)
- is_active = false

### Side Quests (12 quests, 4 categories x 3 tiers)

1. **Harvest Time I/II/III** - reach_production metal (2180 → 4360 → 8720)
2. **Gathering He3 I/II/III** - reach_production he3 (2360 → 4720 → 9440)
3. **Raising Morale I/II/III** - reach_production gold (2800 → 5600 → 11200)
4. **Plentiful Resources I/II/III** - reach_storage (50k → 200k → 1M)

Tier 1 unlocked at start, Tier 2/3 unlock after claiming previous tier.

### Daily Quests (6 quests)

| Quest Key | Display Name | Requirement | Points | Phase |
|-----------|--------------|-------------|--------|-------|
| daily_login | Daily Log In | login | 10 | 1 |
| daily_collect_dues | Collect Your Dues | harvest_resources (1x) | 4 | 1 |
| daily_need_for_speed | Need for Speed | use_speedup | 3 | 1 |
| daily_stockpiling | Stockpiling | harvest_resources (3x) | 1 per harvest (max 3) | 1 |
| daily_donations | Donations | donate_resources (200k) | 6 | 2 |
| daily_restricted_instances | Restricted Instances | complete_instance (2x) | 5 per attempt (max 10) | 2 |

**Daily Tier Rewards:**
- Bronze: 10 pts → loudspeaker
- Silver: 30 pts → resource_box
- Gold: 50 pts → sp_card
- Diamond: 70 pts → raw_gemstone

---

## Backend Endpoints (4)

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | /api/quests | ListQuests | Returns main_quests, side_quests, current_main_quest |
| GET | /api/quests/daily | GetDailyQuests | Returns today's daily progress, quests, tier_rewards |
| POST | /api/quests/{id}/claim | ClaimQuest | Claims reward, awards resources, unlocks next quest |
| POST | /api/quests/daily/claim-tier | ClaimDailyTier | Claims tier reward (bronze/silver/gold/diamond) |

### Quest Service

**`quest_service.UpdateQuestProgress(playerID, requirementType, requirementTarget, incrementBy)`**
- Called from handlers after actions complete
- Finds all active quests matching (type, target)
- Increments progress_value
- Auto-marks as 'completed' if requirement met
- Integrated in: buildings, research, blueprints, ships, fleets

**Integration points:**
- `buildings.go`: build_building, upgrade_building, harvest_resources
- `research.go`: research_tech (on complete)
- `blueprints.go`: use_blueprint (on activate)
- `ship_designs.go`: create_ship_design
- `ship_factory.go`: build_ships (on start production)
- `fleets.go`: create_fleet
- `instances.go`: complete_instance (when implemented)

---

## Frontend Components

### QuestPanel.tsx
- 3 tabs: Main Quests, Side Quests, Daily Quests
- Visual badges showing claimable count per tab
- Progress bars for multi-step quests (e.g., stockpiling 2/3)
- Claim buttons with loading state
- Reward flash notifications ("+500M +400H +500G")
- ESC to close

### useQuests.ts Hook
```typescript
const {
  quests,          // { main_quests, side_quests, current_main_quest }
  daily,           // { date, daily_points, quests, tier_rewards }
  loading,
  error,
  refreshDaily,    // () => Promise<void>
  claim,           // (questId) => Promise<{quest_key, rewards, next_quest_unlocked}>
  claimTier        // (tier) => Promise<{tier, reward}>
} = useQuests()
```

### Types (types/index.ts)
- PlayerQuestWithType
- DailyQuestEntry
- DailyTierReward

---

## Reward Flow

### Claiming Main/Side Quest:
1. POST /api/quests/{id}/claim
2. Backend validates status = 'completed'
3. Updates player_quests.status → 'claimed', claimed_at = now
4. Awards resources to homeworld planet: `UPDATE resources SET metal += reward_metal, ...`
5. Unlocks next quest in chain: finds quest with prerequisite_quest_id = current quest, sets status → 'available'
6. Returns rewards + next_quest_unlocked

### Claiming Daily Tier:
1. POST /api/quests/daily/claim-tier with {"tier": "bronze"}
2. Backend validates daily_points >= required_points
3. Validates tier not already claimed
4. Appends tier to tier_rewards_claimed_json
5. Returns tier reward item (loudspeaker, resource_box, etc.)

---

## Quest Progression Examples

### Main Quest Flow:
1. Player starts: main_01 status = 'available', rest = 'locked'
2. Player harvests resources → UpdateQuestProgress("harvest_resources", "resource_warehouse", 1)
3. main_01 progress 1/1 → status = 'completed'
4. Player clicks Claim → awards 450M/950H/500G, unlocks main_03 (status → 'available')
5. Player builds Technology Center Lv1 → main_03 completes
6. Repeat until end of chain

### Side Quest Flow:
1. Player starts: All Tier 1 side quests = 'available' (no prerequisite)
2. Player reaches 2180 Metal/hr production → Harvest Time I completes
3. Player claims reward → Harvest Time II unlocks (Tier 2)

### Daily Quest Flow:
1. Player logs in → daily_login completed → 10 points
2. Player harvests once → daily_collect_dues completed → +4 points (14 total)
3. Player harvests 3x from warehouse → daily_stockpiling 3/3 → +3 points (17 total)
4. Player claims Bronze tier (10 pts required, has 17) → receives loudspeaker

---

## Known Limitations / Future Work

1. **Auto-progress not yet wired for all requirements:**
   - recruit_commander (Phase 2, needs Command Center implementation)
   - replenish_ammo (Phase 2, needs fleet ammo system)
   - build_defense (Phase 2, needs Space Base implementation)
   - complete_instance (Phase 2, needs Instance attempt flow integration)
   - donate_resources (Phase 2+, needs Corp system)
   - login (daily, needs session tracking)
   - use_speedup (daily, needs speedup item system)

2. **Blueprint rewards implemented but not validated:**
   - Main quest chain grants 7 blueprints as rewards
   - Need to verify blueprint IDs match seed data in blueprints table

3. **Item rewards not yet consumable:**
   - loudspeaker, construction_card, truce_card, etc. stored in reward_item_json
   - Need inventory/item system to store and use these

4. **No QA report yet:**
   - Quest system implemented by user outside of team workflow
   - Recommend creating `/docs/qa/quest-system-report.md` with integration tests

---

## GDD Documentation

**Section 2.10** - Progression & Quest System (already documented)
- Main quest chain table (28 quests)
- Side quest categories
- Daily quest point system
- Quest state machine
- Quest phasing (Phase 1/2/deferred)

---

## Next Steps

1. ✅ **Document** - This file created
2. ⏳ **QA** - Create QA report, verify integration with all actions
3. ⏳ **Complete auto-progress** - Wire up remaining requirement types (recruit, ammo, defense, instance, donate, login, speedup)
4. ⏳ **Validate blueprints** - Verify quest rewards match blueprint seed data
5. ⏳ **UI polish** - Quest icon in SideNav, notification dot when quests claimable
6. ⏳ **Item system** - Implement inventory to store/use quest reward items

---

**Implementation Quality:** ⭐⭐⭐⭐⭐ (5/5)
- Clean separation: quest_types (reference) vs player_quests (state)
- Proper transactions for claim flow
- Auto-unlock next quest in chain
- Daily point accumulation with tier rewards
- Well-structured frontend with loading states and error handling
