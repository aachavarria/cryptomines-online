package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
)

type instanceWithBlueprints struct {
	models.Instance
	BlueprintIDs []int `json:"blueprint_ids"`
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instanceWithBlueprints{Instance: inst, BlueprintIDs: bpIDs})
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
		if len(stacks) == 0 {
			http.Error(w, `{"error":"fleet has no ships: `+fid+`"}`, http.StatusConflict)
			return
		}
	}

	// Simplified combat resolution for Phase 2
	// In Phase 2, player always wins normal instances (combat engine is future work)
	// But we still apply the reward formulas from GDD 8.10.4
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := "attacker_win"
	totalRounds := rng.Intn(30) + 5

	// GDD 8.10.4: Resource rewards
	instanceNum := inst.Difficulty
	metalReward := int64(instanceNum)*500 + int64(rng.Intn(instanceNum*200+1))
	he3Reward := int64(instanceNum)*400 + int64(rng.Intn(instanceNum*150+1))
	goldReward := int64(instanceNum)*600 + int64(rng.Intn(instanceNum*250+1))

	// Blueprint roll: 10% chance
	var blueprintDrop *int
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

	// Create combat report
	lootJSON, _ := json.Marshal(map[string]int64{
		"metal": metalReward, "he3": he3Reward, "gold": goldReward,
	})
	var reportID string
	err = tx.QueryRow(
		`INSERT INTO combat_reports (attacker_id, defender_id, combat_type, result, total_rounds, loot_json, he3_consumed)
		 VALUES ($1, $1, 'instance_normal', $2, $3, $4, 0)
		 RETURNING id`,
		playerID, result, totalRounds, string(lootJSON),
	).Scan(&reportID)
	if err != nil {
		log.Printf("Failed to create combat report: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Update instance progress
	tx.Exec(
		`INSERT INTO instance_progress (player_id, instance_id, completed, attempts, last_attempt_at)
		 VALUES ($1, $2, true, 1, now())
		 ON CONFLICT (player_id, instance_id) DO UPDATE
		 SET completed = true, attempts = instance_progress.attempts + 1, last_attempt_at = now()`,
		playerID, inst.ID,
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

