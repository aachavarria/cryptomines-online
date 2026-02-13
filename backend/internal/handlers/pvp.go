package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/errs"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/workers"
	"github.com/lib/pq"
)

type attackPlanetRequest struct {
	DefenderPlanetID string   `json:"defender_planet_id"`
	FleetIDs         []string `json:"fleet_ids"`
}

// New async response — fleet dispatched, combat resolved later by worker
type attackDispatchResponse struct {
	PendingAttackID string `json:"pending_attack_id"`
	TravelSeconds   int    `json:"travel_seconds"`
	ArrivalAt       string `json:"arrival_at"`
	FleetsDispatched int   `json:"fleets_dispatched"`
}

type pvpLoot struct {
	Metal int64 `json:"metal"`
	He3   int64 `json:"he3"`
	Gold  int64 `json:"gold"`
}

// AttackPlanet handles POST /api/pvp/attack
// Now dispatches fleets with travel time instead of resolving combat instantly.
// Combat is resolved by the PvP worker when fleets arrive.
func AttackPlanet(w http.ResponseWriter, r *http.Request) {
	attackerID := middleware.GetPlayerID(r)

	var req attackPlanetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if len(req.FleetIDs) == 0 {
		http.Error(w, `{"error":"at least one fleet required"}`, http.StatusBadRequest)
		return
	}

	// Get defender planet and owner
	var defenderID string
	var defenderPlanetName string
	var defenderPosX, defenderPosY int
	err := database.DB.QueryRow(`
		SELECT p.player_id, p.name, p.position_x, p.position_y FROM planets p WHERE p.id = $1
	`, req.DefenderPlanetID).Scan(&defenderID, &defenderPlanetName, &defenderPosX, &defenderPosY)

	if err != nil {
		if err == sql.ErrNoRows {
			errs.NotFound(fmt.Sprintf("Target planet %s not found", req.DefenderPlanetID)).WriteJSON(w, http.StatusNotFound)
		} else {
			log.Printf("Failed to get defender planet: %v", err)
			errs.InternalError("Failed to find target planet").WriteJSON(w, http.StatusInternalServerError)
		}
		return
	}

	// Cannot attack own planet
	if defenderID == attackerID {
		http.Error(w, `{"error":"cannot attack own planet"}`, http.StatusConflict)
		return
	}

	// Check attacker SP
	var sp int
	err = database.DB.QueryRow(`SELECT space_points FROM players WHERE id = $1`, attackerID).Scan(&sp)
	if err != nil {
		log.Printf("Failed to get SP: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if sp < 1 {
		http.Error(w, `{"error":"not enough Space Points (need 1 SP)"}`, http.StatusConflict)
		return
	}

	// Check defender truce protection
	var protUntil sql.NullTime
	err = database.DB.QueryRow(`SELECT protection_until FROM planets WHERE id = $1`, req.DefenderPlanetID).Scan(&protUntil)
	if err == nil && protUntil.Valid && protUntil.Time.After(time.Now()) {
		remaining := int(time.Until(protUntil.Time).Seconds())
		http.Error(w, fmt.Sprintf(`{"error":"target under truce protection","remaining_seconds":%d}`, remaining), http.StatusConflict)
		return
	}

	// Check attacker truce protection (bidirectional — can't attack while protected)
	var attackerProtUntil sql.NullTime
	err = database.DB.QueryRow(`
		SELECT protection_until FROM planets WHERE player_id = $1 AND is_homeworld = true
	`, attackerID).Scan(&attackerProtUntil)
	if err == nil && attackerProtUntil.Valid && attackerProtUntil.Time.After(time.Now()) {
		http.Error(w, `{"error":"cannot attack while under truce protection"}`, http.StatusConflict)
		return
	}

	// Check attacker cooldown (5 minutes per target)
	var lastAttack sql.NullTime
	err = database.DB.QueryRow(`
		SELECT MAX(created_at) FROM combat_reports
		WHERE attacker_id = $1 AND defender_id = $2 AND combat_type = 'pvp'
	`, attackerID, defenderID).Scan(&lastAttack)

	if err == nil && lastAttack.Valid {
		cooldown := 5 * time.Minute
		if time.Since(lastAttack.Time) < cooldown {
			remaining := cooldown - time.Since(lastAttack.Time)
			http.Error(w, fmt.Sprintf(`{"error":"attack cooldown active","remaining_seconds":%d}`, int(remaining.Seconds())), http.StatusConflict)
			return
		}
	}

	// Check no existing pending attack against same target
	var existingAttack bool
	database.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM pending_attacks
		WHERE attacker_id = $1 AND defender_planet_id = $2 AND status = 'traveling')
	`, attackerID, req.DefenderPlanetID).Scan(&existingAttack)
	if existingAttack {
		http.Error(w, `{"error":"already have an attack in transit to this target"}`, http.StatusConflict)
		return
	}

	// Validate attacker fleets
	for _, fid := range req.FleetIDs {
		fleet, err := getOwnedFleet(fid, attackerID)
		if err != nil {
			http.Error(w, `{"error":"fleet not found: `+fid+`"}`, http.StatusNotFound)
			return
		}
		if fleet.Status != "stationed" {
			http.Error(w, `{"error":"fleet not stationed: `+fid+`"}`, http.StatusConflict)
			return
		}
		stacks := getFleetStacks(fid)
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

	// Get attacker planet position
	attackerPosX, attackerPosY, err := workers.GetPlanetPosition(getAttackerHomeworld(attackerID))
	if err != nil {
		log.Printf("Failed to get attacker position: %v", err)
		http.Error(w, `{"error":"failed to get attacker position"}`, http.StatusInternalServerError)
		return
	}

	// Calculate fleet speed (slowest ship)
	fleetSpeed := workers.GetFleetSpeed(req.FleetIDs)

	// Calculate travel time
	travelSeconds := workers.CalculateTravelTime(attackerPosX, attackerPosY, defenderPosX, defenderPosY, fleetSpeed)
	arrivalAt := time.Now().Add(time.Duration(travelSeconds) * time.Second)

	// Deduct 1 SP
	database.DB.Exec(`UPDATE players SET space_points = space_points - 1 WHERE id = $1`, attackerID)

	// Deduct He3 fuel for travel
	workers.DeductTravelHe3(attackerID, req.FleetIDs, travelSeconds)

	// Create pending attack
	var pendingID string
	err = database.DB.QueryRow(`
		INSERT INTO pending_attacks (attacker_id, defender_id, defender_planet_id, fleet_ids, travel_seconds, arrival_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, attackerID, defenderID, req.DefenderPlanetID, pq.Array(req.FleetIDs), travelSeconds, arrivalAt).Scan(&pendingID)
	if err != nil {
		log.Printf("Failed to create pending attack: %v", err)
		http.Error(w, `{"error":"failed to dispatch attack"}`, http.StatusInternalServerError)
		return
	}

	// Set fleets to traveling
	for _, fid := range req.FleetIDs {
		database.DB.Exec(`
			UPDATE fleets SET status = 'traveling', destination_x = $1, destination_y = $2,
			       arrival_at = $3, updated_at = now()
			WHERE id = $4
		`, defenderPosX, defenderPosY, arrivalAt, fid)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attackDispatchResponse{
		PendingAttackID:  pendingID,
		TravelSeconds:    travelSeconds,
		ArrivalAt:        arrivalAt.Format(time.RFC3339),
		FleetsDispatched: len(req.FleetIDs),
	})
}

// CancelAttack handles POST /api/pvp/cancel/{id}
// Cancels a pending attack and recalls fleets
func CancelAttack(w http.ResponseWriter, r *http.Request) {
	attackerID := middleware.GetPlayerID(r)
	attackID := r.PathValue("id")

	// Get pending attack
	var status string
	var fleetIDs []string
	var departAt time.Time
	err := database.DB.QueryRow(`
		SELECT status, fleet_ids, depart_at FROM pending_attacks
		WHERE id = $1 AND attacker_id = $2
	`, attackID, attackerID).Scan(&status, pq.Array(&fleetIDs), &departAt)

	if err != nil {
		http.Error(w, `{"error":"attack not found"}`, http.StatusNotFound)
		return
	}
	if status != "traveling" {
		http.Error(w, `{"error":"attack cannot be cancelled (already resolved)"}`, http.StatusConflict)
		return
	}

	// Calculate return time (time elapsed so far * 0.5, min 15s)
	elapsed := time.Since(departAt).Seconds()
	returnSeconds := int(elapsed * 0.5)
	if returnSeconds < 15 {
		returnSeconds = 15
	}

	// Cancel the attack
	database.DB.Exec(`UPDATE pending_attacks SET status = 'cancelled' WHERE id = $1`, attackID)

	// Set fleets to returning
	for _, fid := range fleetIDs {
		database.DB.Exec(`
			UPDATE fleets SET status = 'returning',
			       arrival_at = now() + ($1 || ' seconds')::interval, updated_at = now()
			WHERE id = $2 AND status = 'traveling'
		`, returnSeconds, fid)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"cancelled":      true,
		"return_seconds": returnSeconds,
	})
}

// GetPendingAttacks handles GET /api/pvp/pending
// Returns all outgoing attacks for the player
func GetPendingAttacks(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	rows, err := database.DB.Query(`
		SELECT pa.id, pa.defender_id, pa.defender_planet_id, pa.fleet_ids, pa.status,
		       pa.travel_seconds, pa.depart_at, pa.arrival_at, pa.return_at, pa.combat_report_id,
		       p.name as planet_name, pl.username as defender_name
		FROM pending_attacks pa
		JOIN planets p ON p.id = pa.defender_planet_id
		JOIN players pl ON pl.id = pa.defender_id
		WHERE pa.attacker_id = $1
		ORDER BY pa.created_at DESC
		LIMIT 20
	`, playerID)
	if err != nil {
		log.Printf("Failed to get pending attacks: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type pendingAttack struct {
		ID              string    `json:"id"`
		DefenderID      string    `json:"defender_id"`
		DefenderPlanet  string    `json:"defender_planet_id"`
		FleetIDs        []string  `json:"fleet_ids"`
		Status          string    `json:"status"`
		TravelSeconds   int       `json:"travel_seconds"`
		DepartAt        time.Time `json:"depart_at"`
		ArrivalAt       time.Time `json:"arrival_at"`
		ReturnAt        *time.Time `json:"return_at"`
		CombatReportID  *string   `json:"combat_report_id"`
		PlanetName      string    `json:"planet_name"`
		DefenderName    string    `json:"defender_name"`
	}

	attacks := []pendingAttack{}
	for rows.Next() {
		var a pendingAttack
		if err := rows.Scan(&a.ID, &a.DefenderID, &a.DefenderPlanet, pq.Array(&a.FleetIDs),
			&a.Status, &a.TravelSeconds, &a.DepartAt, &a.ArrivalAt, &a.ReturnAt,
			&a.CombatReportID, &a.PlanetName, &a.DefenderName); err != nil {
			log.Printf("Failed to scan pending attack: %v", err)
			continue
		}
		attacks = append(attacks, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attacks)
}

// GetIncomingAttacks handles GET /api/radar/incoming
// Returns incoming attacks based on player's Radar building level
func GetIncomingAttacks(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	// Get player's homeworld
	var planetID string
	err := database.DB.QueryRow(`SELECT id FROM planets WHERE player_id = $1 AND is_homeworld = true`, playerID).Scan(&planetID)
	if err != nil {
		http.Error(w, `{"error":"planet not found"}`, http.StatusInternalServerError)
		return
	}

	// Get radar building level
	var radarLevel int
	err = database.DB.QueryRow(`
		SELECT COALESCE(b.level, 0)
		FROM buildings b
		WHERE b.planet_id = $1 AND b.building_type = 'radar' AND b.level > 0
		LIMIT 1
	`, planetID).Scan(&radarLevel)
	if err != nil {
		radarLevel = 0
	}

	// No radar = no detection
	if radarLevel == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"radar_level":      0,
			"incoming_attacks": []interface{}{},
		})
		return
	}

	// Get detection advance time (seconds before arrival)
	var detectAdvance int
	err = database.DB.QueryRow(`SELECT detect_advance_seconds FROM radar_levels WHERE level = $1`, radarLevel).Scan(&detectAdvance)
	if err != nil {
		detectAdvance = 1800 // default 30 min
	}

	// Query incoming attacks within detection window
	rows, err := database.DB.Query(`
		SELECT pa.id, pa.attacker_id, pa.fleet_ids, pa.arrival_at, pa.depart_at,
		       ap.position_x as origin_x, ap.position_y as origin_y,
		       apl.username as attacker_name
		FROM pending_attacks pa
		JOIN planets ap ON ap.player_id = pa.attacker_id AND ap.is_homeworld = true
		JOIN players apl ON apl.id = pa.attacker_id
		WHERE pa.defender_id = $1 AND pa.status = 'traveling'
		  AND pa.arrival_at <= now() + ($2 || ' seconds')::interval
		ORDER BY pa.arrival_at ASC
	`, playerID, detectAdvance)
	if err != nil {
		log.Printf("Failed to get incoming attacks: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type incomingAttack struct {
		ID           string    `json:"id"`
		ArrivalAt    time.Time `json:"arrival_at"`
		// Revealed based on radar level:
		OriginX      *int      `json:"origin_x,omitempty"`      // Lv3+
		OriginY      *int      `json:"origin_y,omitempty"`      // Lv3+
		FleetCount   *int      `json:"fleet_count,omitempty"`   // Lv5+
		AttackerName *string   `json:"attacker_name,omitempty"` // Lv7+
	}

	attacks := []incomingAttack{}
	for rows.Next() {
		var attackID, attackerID string
		var fleetIDs []string
		var arrivalAt, departAt time.Time
		var originX, originY int
		var attackerName string

		if err := rows.Scan(&attackID, &attackerID, pq.Array(&fleetIDs), &arrivalAt, &departAt,
			&originX, &originY, &attackerName); err != nil {
			continue
		}

		incoming := incomingAttack{
			ID:        attackID,
			ArrivalAt: arrivalAt,
		}

		// Radar Lv3+: reveal origin coordinates
		if radarLevel >= 3 {
			incoming.OriginX = &originX
			incoming.OriginY = &originY
		}

		// Radar Lv5+: reveal fleet count
		if radarLevel >= 5 {
			fleetCount := len(fleetIDs)
			incoming.FleetCount = &fleetCount
		}

		// Radar Lv7+: reveal attacker name
		if radarLevel >= 7 {
			incoming.AttackerName = &attackerName
		}

		attacks = append(attacks, incoming)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"radar_level":      radarLevel,
		"detect_advance":   detectAdvance,
		"incoming_attacks": attacks,
	})
}

// GetPlayerSP handles GET /api/player/sp
func GetPlayerSP(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	var sp, maxSP int
	err := database.DB.QueryRow(`SELECT space_points, max_space_points FROM players WHERE id = $1`, playerID).Scan(&sp, &maxSP)
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"space_points":     sp,
		"max_space_points": maxSP,
	})
}

// SearchPlanets handles GET /api/pvp/search?query=planetname
func SearchPlanets(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	query := r.URL.Query().Get("query")

	if query == "" {
		http.Error(w, `{"error":"query parameter required"}`, http.StatusBadRequest)
		return
	}

	rows, err := database.DB.Query(`
		SELECT p.id, p.name, p.player_id, pl.username
		FROM planets p
		JOIN players pl ON pl.id = p.player_id
		WHERE p.player_id != $1 AND (p.name ILIKE $2 OR pl.username ILIKE $2)
		LIMIT 50
	`, playerID, "%"+query+"%")

	if err != nil {
		log.Printf("Failed to search planets: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type searchResult struct {
		PlanetID   string `json:"planet_id"`
		PlanetName string `json:"planet_name"`
		PlayerID   string `json:"player_id"`
		PlayerName string `json:"player_name"`
	}

	results := []searchResult{}
	for rows.Next() {
		var r searchResult
		if err := rows.Scan(&r.PlanetID, &r.PlanetName, &r.PlayerID, &r.PlayerName); err != nil {
			continue
		}
		results = append(results, r)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// getAttackerHomeworld returns the homeworld planet ID for a player
func getAttackerHomeworld(playerID string) string {
	var planetID string
	database.DB.QueryRow(`SELECT id FROM planets WHERE player_id = $1 AND is_homeworld = true`, playerID).Scan(&planetID)
	return planetID
}
