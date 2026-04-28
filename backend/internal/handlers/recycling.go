package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
)

// applyDevModeRecycling checks if dev-mode is enabled and returns 5 seconds if true, otherwise returns the original seconds
func applyDevModeRecycling(r *http.Request, seconds int) int {
	if r.Header.Get("X-Dev-Mode") == "true" {
		return 5
	}
	return seconds
}

type availableShip struct {
	ID         string `json:"id"`
	DesignName string `json:"design_name"`
	HullClass  string `json:"hull_class"`
	Quantity   int    `json:"quantity"`
}

type startRecycleRequest struct {
	ShipDesignID string `json:"ship_design_id"`
}

type startRecycleResponse struct {
	JobID           string    `json:"job_id"`
	MetalGained     int64     `json:"metal_gained"`
	He3Gained       int64     `json:"he3_gained"`
	GoldGained      int64     `json:"gold_gained"`
	DurationSeconds int       `json:"duration_seconds"`
	CompletedAt     time.Time `json:"completed_at"`
}

// StartRecycle handles POST /api/recycling-plant/recycle
func StartRecycle(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	var req startRecycleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Verify ship design belongs to player and has available quantity
	var availableQuantity int
	err := database.DB.QueryRow(`
		SELECT quantity
		FROM ships
		WHERE ship_design_id = $1 AND player_id = $2
	`, req.ShipDesignID, playerID).Scan(&availableQuantity)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"ship design not found or no ships available"}`, http.StatusNotFound)
		} else {
			log.Printf("Failed to get ship quantity: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	// Check if there are ships available
	if availableQuantity <= 0 {
		http.Error(w, `{"error":"no ships available to recycle"}`, http.StatusConflict)
		return
	}

	// Get ship design to calculate resources
	var hullTypeID int
	var moduleIDs []byte
	err = database.DB.QueryRow(`
		SELECT hull_type_id, module_ids
		FROM ship_designs
		WHERE id = $1
	`, req.ShipDesignID).Scan(&hullTypeID, &moduleIDs)

	if err != nil {
		log.Printf("Failed to get ship design: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Get hull type costs
	var hullMetal, hullHe3, hullGold int64
	err = database.DB.QueryRow(`
		SELECT base_metal_cost, base_he3_cost, base_gold_cost
		FROM hull_types
		WHERE id = $1
	`, hullTypeID).Scan(&hullMetal, &hullHe3, &hullGold)

	if err != nil {
		log.Printf("Failed to get hull type costs: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Parse module IDs and sum costs
	var modules []int
	if len(moduleIDs) > 0 {
		if err := json.Unmarshal(moduleIDs, &modules); err != nil {
			log.Printf("Failed to parse module IDs: %v", err)
		}
	}

	var moduleMetal, moduleHe3, moduleGold int64
	if len(modules) > 0 {
		rows, err := database.DB.Query(`
			SELECT COALESCE(SUM(metal_cost), 0), COALESCE(SUM(he3_cost), 0), COALESCE(SUM(gold_cost), 0)
			FROM module_types
			WHERE id = ANY($1)
		`, modules)
		if err == nil {
			defer rows.Close()
			if rows.Next() {
				rows.Scan(&moduleMetal, &moduleHe3, &moduleGold)
			}
		}
	}

	// Total costs (hull + modules) with 70% recovery rate
	totalMetal := (hullMetal + moduleMetal) * 70 / 100
	totalHe3 := (hullHe3 + moduleHe3) * 70 / 100
	totalGold := (hullGold + moduleGold) * 70 / 100

	// Duration: base 60 seconds + 10 seconds per 1000 total cost
	totalCost := totalMetal + totalHe3 + totalGold
	durationSeconds := 60 + int(totalCost/1000)*10
	if durationSeconds > 3600 {
		durationSeconds = 3600 // Max 1 hour
	}

	effectiveRecycleTime := applyDevModeRecycling(r, durationSeconds)
	completedAt := time.Now().Add(time.Duration(effectiveRecycleTime) * time.Second)

	// Create recycling job and delete ship instance
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var jobID string
	err = tx.QueryRow(`
		INSERT INTO recycling_jobs (player_id, ship_design_id, metal_gained, he3_gained, gold_gained, duration_seconds, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, playerID, req.ShipDesignID, totalMetal, totalHe3, totalGold, durationSeconds, completedAt).Scan(&jobID)

	if err != nil {
		log.Printf("Failed to create recycling job: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Decrement ship quantity by 1
	_, err = tx.Exec(`
		UPDATE ships
		SET quantity = quantity - 1
		WHERE player_id = $1 AND ship_design_id = $2
	`, playerID, req.ShipDesignID)
	if err != nil {
		log.Printf("Failed to decrement ship quantity: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(startRecycleResponse{
		JobID:           jobID,
		MetalGained:     totalMetal,
		He3Gained:       totalHe3,
		GoldGained:      totalGold,
		DurationSeconds: durationSeconds,
		CompletedAt:     completedAt,
	})
}

// ListRecyclingJobs handles GET /api/recycling-plant/jobs
func ListRecyclingJobs(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	rows, err := database.DB.Query(`
		SELECT id, ship_design_id, metal_gained, he3_gained, gold_gained,
		       duration_seconds, started_at, completed_at, collected
		FROM recycling_jobs
		WHERE player_id = $1 AND collected = false
		ORDER BY started_at DESC
	`, playerID)
	if err != nil {
		log.Printf("Failed to list recycling jobs: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	jobs := []models.RecyclingJob{}
	for rows.Next() {
		var job models.RecyclingJob
		err := rows.Scan(
			&job.ID, &job.ShipDesignID, &job.MetalGained, &job.He3Gained, &job.GoldGained,
			&job.DurationSeconds, &job.StartedAt, &job.CompletedAt, &job.Collected,
		)
		if err != nil {
			log.Printf("Failed to scan recycling job: %v", err)
			continue
		}
		jobs = append(jobs, job)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// CollectRecycle handles POST /api/recycling-plant/collect/{id}
func CollectRecycle(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	jobID := r.PathValue("id")

	// Get job and verify ownership
	var metalGained, he3Gained, goldGained int64
	var completedAt sql.NullTime
	var collected bool

	err := database.DB.QueryRow(`
		SELECT metal_gained, he3_gained, gold_gained, completed_at, collected
		FROM recycling_jobs
		WHERE id = $1 AND player_id = $2
	`, jobID, playerID).Scan(&metalGained, &he3Gained, &goldGained, &completedAt, &collected)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"job not found"}`, http.StatusNotFound)
		} else {
			log.Printf("Failed to get recycling job: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	if collected {
		http.Error(w, `{"error":"already collected"}`, http.StatusConflict)
		return
	}

	if !completedAt.Valid || completedAt.Time.After(time.Now()) {
		http.Error(w, `{"error":"not yet completed"}`, http.StatusConflict)
		return
	}

	// Award resources and mark as collected
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Award to homeworld
	tx.Exec(`
		UPDATE resources
		SET metal = metal + $1, he3 = he3 + $2, gold = gold + $3, updated_at = now()
		WHERE planet_id = (
			SELECT p.id FROM planets p WHERE p.player_id = $4 AND p.is_homeworld = true LIMIT 1
		)
	`, metalGained, he3Gained, goldGained, playerID)

	// Mark as collected
	_, err = tx.Exec(`
		UPDATE recycling_jobs SET collected = true WHERE id = $1
	`, jobID)
	if err != nil {
		log.Printf("Failed to mark as collected: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"metal_gained": metalGained,
		"he3_gained":   he3Gained,
		"gold_gained":  goldGained,
	})
}

// ListAvailableShips handles GET /api/ship-instances/available
// Returns ship designs that have ships not currently deployed in fleets.
// available = ships.quantity − (ships of the same design assigned to fleets).
func ListAvailableShips(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	rows, err := database.DB.Query(`
		WITH deployed AS (
			SELECT fs.ship_design_id, COALESCE(SUM(fs.ship_count), 0) AS qty
			FROM fleet_stacks fs
			JOIN fleets f ON f.id = fs.fleet_id
			WHERE f.player_id = $1
			GROUP BY fs.ship_design_id
		)
		SELECT s.ship_design_id, sd.name, ht.hull_class,
		       GREATEST(s.quantity - COALESCE(d.qty, 0), 0) AS available
		FROM ships s
		JOIN ship_designs sd ON s.ship_design_id = sd.id
		JOIN hull_types ht ON sd.hull_type_id = ht.id
		LEFT JOIN deployed d ON d.ship_design_id = s.ship_design_id
		WHERE s.player_id = $1 AND s.quantity > 0
		  AND GREATEST(s.quantity - COALESCE(d.qty, 0), 0) > 0
		ORDER BY ht.hull_class, sd.name
	`, playerID)

	if err != nil {
		log.Printf("Failed to list available ships: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type shipWithQuantity struct {
		ID         string `json:"id"`
		DesignName string `json:"design_name"`
		HullClass  string `json:"hull_class"`
		Quantity   int    `json:"quantity"`
	}

	ships := []shipWithQuantity{}
	for rows.Next() {
		var ship shipWithQuantity
		err := rows.Scan(&ship.ID, &ship.DesignName, &ship.HullClass, &ship.Quantity)
		if err != nil {
			log.Printf("Failed to scan ship: %v", err)
			continue
		}
		ships = append(ships, ship)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ships)
}

// CancelRecycle handles DELETE /api/recycling-plant/jobs/{id}
func CancelRecycle(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	jobID := r.PathValue("id")

	// Verify ownership and that it's not collected
	var shipDesignID string
	var collected bool

	err := database.DB.QueryRow(`
		SELECT ship_design_id, collected
		FROM recycling_jobs
		WHERE id = $1 AND player_id = $2
	`, jobID, playerID).Scan(&shipDesignID, &collected)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"job not found"}`, http.StatusNotFound)
		} else {
			log.Printf("Failed to get recycling job: %v", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	if collected {
		http.Error(w, `{"error":"already collected"}`, http.StatusConflict)
		return
	}

	// Cannot restore ship - just delete the job
	// (In GO2, canceling recycling does NOT restore the ship)
	_, err = database.DB.Exec(`DELETE FROM recycling_jobs WHERE id = $1`, jobID)
	if err != nil {
		log.Printf("Failed to delete recycling job: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
