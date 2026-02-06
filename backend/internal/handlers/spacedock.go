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

type spacedockStatus struct {
	Level         int     `json:"level"`
	IsUpgrading   bool    `json:"is_upgrading"`
	RepairPct     float64 `json:"repair_pct"`
	ActiveRepairs int     `json:"active_repairs"`
}

type repairRequest struct {
	CombatReportID string `json:"combat_report_id"`
}

// GetSpacedock handles GET /api/spacedock
func GetSpacedock(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	completeFinishedRepairs(playerID)

	var level int
	var isUpgrading bool
	err := database.DB.QueryRow(
		`SELECT COALESCE(b.level, 0), COALESCE(b.is_upgrading, false)
		 FROM buildings b
		 JOIN building_types bt ON b.building_type = bt.id
		 JOIN planets p ON b.planet_id = p.id
		 WHERE p.player_id = $1 AND bt.name = 'spacedock'
		 ORDER BY b.level DESC LIMIT 1`, playerID,
	).Scan(&level, &isUpgrading)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"spacedock not built"}`, http.StatusNotFound)
			return
		}
		log.Printf("Failed to get spacedock: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	var repairPct float64
	database.DB.QueryRow(
		`SELECT repair_pct FROM spacedock_levels WHERE level = $1`, level,
	).Scan(&repairPct)

	var activeRepairs int
	database.DB.QueryRow(
		`SELECT COUNT(*) FROM spacedock_repairs
		 WHERE player_id = $1 AND repair_finish_at IS NOT NULL AND repair_finish_at > now()`,
		playerID,
	).Scan(&activeRepairs)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spacedockStatus{
		Level:         level,
		IsUpgrading:   isUpgrading,
		RepairPct:     repairPct,
		ActiveRepairs: activeRepairs,
	})
}

// ListRepairs handles GET /api/spacedock/repairs
func ListRepairs(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	completeFinishedRepairs(playerID)

	rows, err := database.DB.Query(
		`SELECT id, player_id, ship_design_id, destroyed_count, repaired_count,
		        repair_finish_at, combat_report_id, created_at
		 FROM spacedock_repairs WHERE player_id = $1
		 ORDER BY created_at DESC`, playerID,
	)
	if err != nil {
		log.Printf("Failed to list repairs: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	repairs := []models.SpacedockRepair{}
	for rows.Next() {
		var rep models.SpacedockRepair
		err := rows.Scan(
			&rep.ID, &rep.PlayerID, &rep.ShipDesignID, &rep.DestroyedCount,
			&rep.RepairedCount, &rep.RepairFinishAt, &rep.CombatReportID, &rep.CreatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan repair: %v", err)
			continue
		}
		repairs = append(repairs, rep)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repairs)
}

// StartRepair handles POST /api/spacedock/repair
// Starts repairing ships from a PvP loss.
func StartRepair(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	var req repairRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Get spacedock level for repair percentage
	var sdLevel int
	err := database.DB.QueryRow(
		`SELECT COALESCE(b.level, 0)
		 FROM buildings b
		 JOIN building_types bt ON b.building_type = bt.id
		 JOIN planets p ON b.planet_id = p.id
		 WHERE p.player_id = $1 AND bt.name = 'spacedock'
		 ORDER BY b.level DESC LIMIT 1`, playerID,
	).Scan(&sdLevel)
	if err != nil || sdLevel < 1 {
		http.Error(w, `{"error":"spacedock not built"}`, http.StatusConflict)
		return
	}

	var repairPct float64
	err = database.DB.QueryRow(
		`SELECT repair_pct FROM spacedock_levels WHERE level = $1`, sdLevel,
	).Scan(&repairPct)
	if err != nil {
		log.Printf("Failed to get repair pct: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Verify combat report exists and belongs to player
	var reportExists bool
	err = database.DB.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM combat_reports WHERE id = $1 AND (attacker_id = $2 OR defender_id = $2))`,
		req.CombatReportID, playerID,
	).Scan(&reportExists)
	if err != nil || !reportExists {
		http.Error(w, `{"error":"combat report not found"}`, http.StatusNotFound)
		return
	}

	// Check if already repaired from this report
	var alreadyRepaired bool
	err = database.DB.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM spacedock_repairs WHERE player_id = $1 AND combat_report_id = $2)`,
		playerID, req.CombatReportID,
	).Scan(&alreadyRepaired)
	if err == nil && alreadyRepaired {
		http.Error(w, `{"error":"already repaired from this battle"}`, http.StatusConflict)
		return
	}

	// For Phase 2: placeholder since we don't have real PvP losses yet
	// The repair will be triggered after PvP combat is implemented
	// For now, respond with the spacedock status
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Message   string  `json:"message"`
		RepairPct float64 `json:"repair_pct"`
	}{
		Message:   "repair started (placeholder for PvP losses)",
		RepairPct: repairPct,
	})
}

// AccelerateRepair handles POST /api/spacedock/accelerate
// Spend Mall Points to accelerate repair.
func AccelerateRepair(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	type accelRequest struct {
		RepairID string `json:"repair_id"`
	}

	var req accelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Get repair
	var rep models.SpacedockRepair
	err := database.DB.QueryRow(
		`SELECT id, player_id, ship_design_id, destroyed_count, repaired_count,
		        repair_finish_at, combat_report_id, created_at
		 FROM spacedock_repairs WHERE id = $1 AND player_id = $2`,
		req.RepairID, playerID,
	).Scan(
		&rep.ID, &rep.PlayerID, &rep.ShipDesignID, &rep.DestroyedCount,
		&rep.RepairedCount, &rep.RepairFinishAt, &rep.CombatReportID, &rep.CreatedAt,
	)
	if err != nil {
		http.Error(w, `{"error":"repair not found"}`, http.StatusNotFound)
		return
	}

	if rep.RepairFinishAt == nil {
		http.Error(w, `{"error":"repair already completed"}`, http.StatusConflict)
		return
	}

	// Cost: 1 Mall Point per 10 minutes remaining
	remaining := time.Until(*rep.RepairFinishAt)
	if remaining <= 0 {
		http.Error(w, `{"error":"repair already finished"}`, http.StatusConflict)
		return
	}

	mallPointCost := int(remaining.Minutes()/10) + 1

	// Deduct mall points
	var remainingMP int
	err = database.DB.QueryRow(
		`UPDATE players SET mall_points = mall_points - $1, updated_at = now()
		 WHERE id = $2 AND mall_points >= $1
		 RETURNING mall_points`,
		mallPointCost, playerID,
	).Scan(&remainingMP)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"insufficient mall points"}`, http.StatusConflict)
			return
		}
		log.Printf("Failed to deduct mall points: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Complete repair immediately
	_, err = database.DB.Exec(
		`UPDATE spacedock_repairs SET repair_finish_at = now() WHERE id = $1`, req.RepairID,
	)
	if err != nil {
		log.Printf("Failed to accelerate repair: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	completeFinishedRepairs(playerID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Accelerated    bool `json:"accelerated"`
		MallPointsUsed int  `json:"mall_points_used"`
		MallPointsLeft int  `json:"mall_points_left"`
	}{
		Accelerated:    true,
		MallPointsUsed: mallPointCost,
		MallPointsLeft: remainingMP,
	})
}

// completeFinishedRepairs auto-completes repairs that are done.
func completeFinishedRepairs(playerID string) {
	now := time.Now()
	rows, err := database.DB.Query(
		`SELECT id, ship_design_id, repaired_count FROM spacedock_repairs
		 WHERE player_id = $1 AND repair_finish_at IS NOT NULL AND repair_finish_at <= $2`,
		playerID, now,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var repairID, designID string
		var repairedCount int
		if rows.Scan(&repairID, &designID, &repairedCount) != nil {
			continue
		}

		// Add repaired ships back to player's pool
		database.DB.Exec(
			`UPDATE ships SET quantity = quantity + $1, updated_at = now()
			 WHERE player_id = $2 AND ship_design_id = $3`,
			repairedCount, playerID, designID,
		)

		// Mark repair as completed
		database.DB.Exec(
			`UPDATE spacedock_repairs SET repair_finish_at = NULL WHERE id = $1`, repairID,
		)
	}
}
