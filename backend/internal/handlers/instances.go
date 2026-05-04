package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/cryptomines-online/backend/internal/combat"
	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
	"github.com/cryptomines-online/backend/internal/services"
)

// enemyFleetSummary is the parsed/aggregated view of enemy_fleets_json that
// the InstanceDetail modal expects. The raw JSON groups entries per hull_type
// + grid position; here we collapse to one row per hull_class and estimate
// rough power so the player sees a digestible preview.
type enemyFleetSummary struct {
	HullClass     string `json:"hull_class"`
	ShipCount     int    `json:"ship_count"`
	PowerEstimate int    `json:"power_estimate"`
}

type instanceDetailResponse struct {
	models.Instance
	EnemyFleets   []enemyFleetSummary `json:"enemy_fleets"`
	BlueprintPool []int               `json:"blueprint_pool"`
	BlueprintIDs  []int               `json:"blueprint_ids"` // legacy alias
}

type attemptRequest struct {
	FleetIDs []string `json:"fleet_ids"`
}

type attemptResponse struct {
	ReportID    string        `json:"report_id"`
	Result      string        `json:"result"`
	TotalRounds int           `json:"total_rounds"`
	ExpGained   int           `json:"exp_gained"`
	TreasureBox *treasureBox  `json:"treasure_box"`
	Losses      *battleLosses `json:"losses"`
}

type treasureBox struct {
	Resources struct {
		Metal int64 `json:"metal"`
		He3   int64 `json:"he3"`
		Gold  int64 `json:"gold"`
	} `json:"resources"`
	BlueprintID *int `json:"blueprint_id"`
}

type battleLosses struct {
	ShipsDestroyed map[string]int `json:"ships_destroyed"`
	He3Consumed    int64          `json:"he3_consumed"`
}

// ListInstances handles GET /api/instances
func ListInstances(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(
		`SELECT id, name, type, difficulty, required_level, max_fleets,
		        ships_lost_on_defeat, he3_lost_on_defeat, exp_reward,
		        enemy_fleets_json, rewards_json, description
		 FROM instances WHERE type = 'normal'
		 ORDER BY difficulty`,
	)
	if err != nil {
		log.Printf("Failed to list instances: %v", err)
		http.Error(w, `{"error":"failed to list instances"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	instances := []models.Instance{}
	for rows.Next() {
		var inst models.Instance
		err := rows.Scan(
			&inst.ID, &inst.Name, &inst.Type, &inst.Difficulty, &inst.RequiredLevel,
			&inst.MaxFleets, &inst.ShipsLostOnDefeat, &inst.He3LostOnDefeat,
			&inst.ExpReward, &inst.EnemyFleetsJSON, &inst.RewardsJSON, &inst.Description,
		)
		if err != nil {
			log.Printf("Failed to scan instance: %v", err)
			continue
		}
		instances = append(instances, inst)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instances)
}

// GetInstance handles GET /api/instances/{id}
func GetInstance(w http.ResponseWriter, r *http.Request) {
	instanceID := r.PathValue("id")

	var inst models.Instance
	err := database.DB.QueryRow(
		`SELECT id, name, type, difficulty, required_level, max_fleets,
		        ships_lost_on_defeat, he3_lost_on_defeat, exp_reward,
		        enemy_fleets_json, rewards_json, description
		 FROM instances WHERE id = $1`, instanceID,
	).Scan(
		&inst.ID, &inst.Name, &inst.Type, &inst.Difficulty, &inst.RequiredLevel,
		&inst.MaxFleets, &inst.ShipsLostOnDefeat, &inst.He3LostOnDefeat,
		&inst.ExpReward, &inst.EnemyFleetsJSON, &inst.RewardsJSON, &inst.Description,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"instance not found"}`, http.StatusNotFound)
			return
		}
		log.Printf("Failed to get instance: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Get blueprint drop pool
	bpRows, err := database.DB.Query(
		`SELECT blueprint_id FROM instance_blueprints WHERE instance_id = $1`, inst.ID,
	)
	if err != nil {
		log.Printf("Failed to get instance blueprints: %v", err)
	}
	defer func() {
		if bpRows != nil {
			bpRows.Close()
		}
	}()

	bpIDs := []int{}
	if bpRows != nil {
		for bpRows.Next() {
			var bpID int
			if err := bpRows.Scan(&bpID); err == nil {
				bpIDs = append(bpIDs, bpID)
			}
		}
	}

	enemyFleets := summarizeEnemyFleets(inst.EnemyFleetsJSON)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instanceDetailResponse{
		Instance:      inst,
		EnemyFleets:   enemyFleets,
		BlueprintPool: bpIDs,
		BlueprintIDs:  bpIDs,
	})
}

// summarizeEnemyFleets parses enemy_fleets_json into a per-hull_class summary.
// The raw shape is [{"hull_type": "weikes_i", "quantity": 55, "grid_row": 0,
// "grid_col": 0}, ...]. We look up hull_class + tier from hull_types and
// estimate power as quantity × tier × 100 (rough ballpark, good enough for a
// pre-battle preview).
func summarizeEnemyFleets(enemyFleetsJSON string) []enemyFleetSummary {
	if enemyFleetsJSON == "" {
		return []enemyFleetSummary{}
	}
	var raw []struct {
		HullType string `json:"hull_type"`
		Quantity int    `json:"quantity"`
	}
	if err := json.Unmarshal([]byte(enemyFleetsJSON), &raw); err != nil {
		log.Printf("Failed to parse enemy_fleets_json: %v", err)
		return []enemyFleetSummary{}
	}
	if len(raw) == 0 {
		return []enemyFleetSummary{}
	}

	type agg struct {
		count int
		power int
	}
	byClass := map[string]*agg{}
	for _, e := range raw {
		var hullClass string
		var tier int
		err := database.DB.QueryRow(
			`SELECT hull_class, tier FROM hull_types WHERE name = $1`, e.HullType,
		).Scan(&hullClass, &tier)
		if err != nil {
			// Unknown hull — skip silently rather than crash the whole preview.
			continue
		}
		a := byClass[hullClass]
		if a == nil {
			a = &agg{}
			byClass[hullClass] = a
		}
		a.count += e.Quantity
		a.power += e.Quantity * tier * 100
	}

	out := make([]enemyFleetSummary, 0, len(byClass))
	// Preserve a stable order: frigate, cruiser, battleship.
	for _, cls := range []string{"frigate", "cruiser", "battleship"} {
		if a, ok := byClass[cls]; ok {
			out = append(out, enemyFleetSummary{HullClass: cls, ShipCount: a.count, PowerEstimate: a.power})
			delete(byClass, cls)
		}
	}
	for cls, a := range byClass {
		out = append(out, enemyFleetSummary{HullClass: cls, ShipCount: a.count, PowerEstimate: a.power})
	}
	return out
}

// AttemptInstance handles POST /api/instances/{id}/attempt
func AttemptInstance(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	instanceID := r.PathValue("id")

	var req attemptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if len(req.FleetIDs) == 0 {
		http.Error(w, `{"error":"at least one fleet required"}`, http.StatusBadRequest)
		return
	}

	// Get instance
	var inst models.Instance
	err := database.DB.QueryRow(
		`SELECT id, name, type, difficulty, required_level, max_fleets,
		        ships_lost_on_defeat, he3_lost_on_defeat, exp_reward,
		        enemy_fleets_json, rewards_json, description
		 FROM instances WHERE id = $1 AND type = 'normal'`, instanceID,
	).Scan(
		&inst.ID, &inst.Name, &inst.Type, &inst.Difficulty, &inst.RequiredLevel,
		&inst.MaxFleets, &inst.ShipsLostOnDefeat, &inst.He3LostOnDefeat,
		&inst.ExpReward, &inst.EnemyFleetsJSON, &inst.RewardsJSON, &inst.Description,
	)
	if err != nil {
		http.Error(w, `{"error":"instance not found or not normal type"}`, http.StatusNotFound)
		return
	}

	// Check player level
	var playerLevel int
	err = database.DB.QueryRow(`SELECT level FROM players WHERE id = $1`, playerID).Scan(&playerLevel)
	if err != nil {
		log.Printf("Failed to get player level: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if playerLevel < inst.RequiredLevel {
		http.Error(w, `{"error":"player level too low"}`, http.StatusConflict)
		return
	}

	// Check fleet count
	if len(req.FleetIDs) > inst.MaxFleets {
		http.Error(w, `{"error":"too many fleets for this instance"}`, http.StatusConflict)
		return
	}

	// Validate all fleets belong to player, are stationed, and have ships
	for _, fid := range req.FleetIDs {
		fleet, err := getOwnedFleet(fid, playerID)
		if err != nil {
			http.Error(w, `{"error":"fleet not found: `+fid+`"}`, http.StatusNotFound)
			return
		}
		if fleet.Status != "stationed" {
			http.Error(w, `{"error":"fleet not stationed: `+fid+`"}`, http.StatusConflict)
			return
		}
		stacks := getFleetStacks(fid)
		// Check fleet has at least one stack with ships (ship_count > 0)
		hasShips := false
		for _, s := range stacks {
			if s.ShipCount > 0 {
				hasShips = true
				break
			}
		}
		if !hasShips {
			http.Error(w, `{"error":"fleet has no ships: `+fid+`"}`, http.StatusConflict)
			return
		}
	}

	// Combat resolution using real combat engine.
	// Load player fleets, merge their stacks into a single combat fleet and
	// promote the strongest commander among them to lead the combined force —
	// otherwise multi-fleet attempts would lose every commander bonus.
	var allStacks []*combat.FleetStack
	var leadCommander *combat.CommanderBonus
	var leadScore int
	var totalHe3Consumed int64
	for _, fid := range req.FleetIDs {
		playerFleet, err := combat.LoadPlayerFleet(fid, playerID)
		if err != nil {
			log.Printf("Failed to load player fleet %s: %v", fid, err)
			http.Error(w, `{"error":"failed to load player fleet"}`, http.StatusInternalServerError)
			return
		}
		allStacks = append(allStacks, playerFleet.Stacks...)
		if c := playerFleet.CommanderBonus; c != nil {
			score := c.Accuracy + c.Dodge + c.Speed + c.Electron
			if leadCommander == nil || score > leadScore {
				leadCommander = c
				leadScore = score
			}
		}
	}

	attackerFleet := &combat.Fleet{
		PlayerID:       playerID,
		FleetID:        "combined",
		CommanderBonus: leadCommander,
		TechBonuses:    &services.TechBonuses{},
		Stacks:         allStacks,
		Formation:      "phalanx",
		Targeting:      "max_attack",
		Side:           "attacker",
	}

	// Load tech bonuses for combined fleet
	if techBonuses, err := services.GetPlayerTechBonuses(playerID); err == nil {
		attackerFleet.TechBonuses = techBonuses
	}

	// Load instance enemy fleet
	defenderFleet, err := combat.LoadInstanceFleet(instanceID)
	if err != nil {
		log.Printf("Failed to load instance fleet: %v", err)
		http.Error(w, `{"error":"failed to load instance fleet"}`, http.StatusInternalServerError)
		return
	}

	// Execute combat
	engine := combat.NewCombatEngine(time.Now().UnixNano())
	combatResult, err := engine.ExecuteCombat(attackerFleet, defenderFleet)
	if err != nil {
		log.Printf("Combat execution failed: %v", err)
		http.Error(w, `{"error":"combat execution failed"}`, http.StatusInternalServerError)
		return
	}

	// Determine result string for database
	var result string
	if combatResult.Winner == "attacker" {
		result = "attacker_win"
	} else if combatResult.Winner == "defender" {
		result = "defender_win"
	} else {
		result = "draw"
	}
	totalRounds := combatResult.TotalRounds

	// Apply casualties to player fleet stacks (combat engine modifies stacks in-place)
	shipsDestroyed := make(map[string]int)
	for _, stack := range attackerFleet.Stacks {
		destroyed := stack.ShipCount - stack.CurrentShips
		if destroyed > 0 {
			shipsDestroyed[stack.ID] = destroyed
			// Update database fleet_stacks
			database.DB.Exec(
				`UPDATE fleet_stacks SET ship_count = $1 WHERE id = $2`,
				stack.CurrentShips, stack.ID,
			)
		}
	}

	// Calculate He3 consumption (simplified: 1 He3 per ship lost)
	totalHe3Consumed = int64(combatResult.AttackerCasualties)

	// Calculate rewards (only on victory)
	var metalReward, he3Reward, goldReward int64
	var blueprintDrop *int

	if result == "attacker_win" {
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))

		// GDD 8.10.4: Resource rewards
		instanceNum := inst.Difficulty
		metalReward = int64(instanceNum)*500 + int64(rng.Intn(instanceNum*200+1))
		he3Reward = int64(instanceNum)*400 + int64(rng.Intn(instanceNum*150+1))
		goldReward = int64(instanceNum)*600 + int64(rng.Intn(instanceNum*250+1))

		// Blueprint roll: 10% chance
		if rng.Float64() < 0.10 {
			bpRows, err := database.DB.Query(
				`SELECT blueprint_id FROM instance_blueprints WHERE instance_id = $1`, inst.ID,
			)
			if err == nil {
				defer bpRows.Close()
				bpIDs := []int{}
				for bpRows.Next() {
					var bpID int
					if bpRows.Scan(&bpID) == nil {
						bpIDs = append(bpIDs, bpID)
					}
				}
				if len(bpIDs) > 0 {
					chosenBP := bpIDs[rng.Intn(len(bpIDs))]
					blueprintDrop = &chosenBP

					// Award unactivated blueprint to player (ignore if already has it)
					database.DB.Exec(
						`INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
						 VALUES ($1, $2, false, 1)
						 ON CONFLICT (player_id, blueprint_id) DO NOTHING`,
						playerID, chosenBP,
					)
				}
			}
		}
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Award resources to homeworld
	tx.Exec(
		`UPDATE resources SET metal = metal + $1, he3 = he3 + $2, gold = gold + $3, updated_at = now()
		 WHERE planet_id = (
		     SELECT p.id FROM planets p WHERE p.player_id = $4 AND p.is_homeworld = true LIMIT 1
		 )`, metalReward, he3Reward, goldReward, playerID,
	)

	// Award EXP to player
	tx.Exec(
		`UPDATE players SET experience = experience + $1, updated_at = now() WHERE id = $2`,
		inst.ExpReward, playerID,
	)

	// Create combat report. Persist the round-by-round log so the UI battle
	// playback can replay every attack instead of only showing summary stats.
	lootJSON, _ := json.Marshal(map[string]int64{
		"metal": metalReward, "he3": he3Reward, "gold": goldReward,
	})
	roundsJSON, _ := json.Marshal(combatResult.DetailedRounds)
	if len(roundsJSON) == 0 {
		roundsJSON = []byte("[]")
	}
	var reportID string
	err = tx.QueryRow(
		`INSERT INTO combat_reports (attacker_id, defender_id, combat_type, result, total_rounds, loot_json, rounds_json, he3_consumed)
		 VALUES ($1, $1, 'instance_normal', $2, $3, $4, $5, $6)
		 RETURNING id`,
		playerID, result, totalRounds, string(lootJSON), string(roundsJSON), totalHe3Consumed,
	).Scan(&reportID)
	if err != nil {
		log.Printf("Failed to create combat report: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Update instance progress (completed only on victory)
	isCompleted := result == "attacker_win"
	tx.Exec(
		`INSERT INTO instance_progress (player_id, instance_id, completed, attempts, last_attempt_at)
		 VALUES ($1, $2, $3, 1, now())
		 ON CONFLICT (player_id, instance_id) DO UPDATE
		 SET completed = CASE WHEN $3 = true THEN true ELSE instance_progress.completed END,
		     attempts = instance_progress.attempts + 1,
		     last_attempt_at = now()`,
		playerID, inst.ID, isCompleted,
	)

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	tb := &treasureBox{}
	tb.Resources.Metal = metalReward
	tb.Resources.He3 = he3Reward
	tb.Resources.Gold = goldReward
	tb.BlueprintID = blueprintDrop

	resp := attemptResponse{
		ReportID:    reportID,
		Result:      result,
		TotalRounds: totalRounds,
		ExpGained:   inst.ExpReward,
		TreasureBox: tb,
		Losses:      &battleLosses{ShipsDestroyed: map[string]int{}, He3Consumed: 0},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetInstanceProgress handles GET /api/instances/progress
func GetInstanceProgress(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	rows, err := database.DB.Query(
		`SELECT ip.id, ip.player_id, ip.instance_id, ip.completed, ip.attempts,
		        ip.best_score, ip.last_attempt_at
		 FROM instance_progress ip
		 WHERE ip.player_id = $1
		 ORDER BY ip.instance_id`, playerID,
	)
	if err != nil {
		log.Printf("Failed to get instance progress: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	progress := []models.InstanceProgress{}
	for rows.Next() {
		var p models.InstanceProgress
		err := rows.Scan(&p.ID, &p.PlayerID, &p.InstanceID, &p.Completed, &p.Attempts, &p.BestScore, &p.LastAttemptAt)
		if err != nil {
			log.Printf("Failed to scan instance progress: %v", err)
			continue
		}
		progress = append(progress, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}

