package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/errs"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
)

const maxShipsPerStack = 3000

var validFormations = map[string]bool{
	"phalanx": true, "diamond": true, "battle_line": true,
	"skirmish": true, "tee_forward": true, "enfilade": true, "tee_reverse": true,
}

var validTargeting = map[string]bool{
	"max_attack": true, "min_attack": true, "max_durability": true,
	"min_durability": true, "closest": true, "by_commander_rank": true,
}

type createFleetRequest struct {
	Name      string `json:"name"`
	PlanetID  string `json:"planet_id"`
	Formation string `json:"formation,omitempty"`
}

type updateFleetRequest struct {
	Name             *string `json:"name,omitempty"`
	Formation        *string `json:"formation,omitempty"`
	TargetingCommand *string `json:"targeting_command,omitempty"`
	CommanderID      *string `json:"commander_id,omitempty"`
}

type assignStackRequest struct {
	ShipDesignID string `json:"ship_design_id"`
	GridRow      int    `json:"grid_row"`
	GridCol      int    `json:"grid_col"`
	ShipCount    int    `json:"ship_count"`
}

type removeStackRequest struct {
	GridRow int `json:"grid_row"`
	GridCol int `json:"grid_col"`
}

type moveFleetRequest struct {
	DestinationX int `json:"destination_x"`
	DestinationY int `json:"destination_y"`
}

type fleetWithStacks struct {
	models.Fleet
	Stacks []models.FleetStack `json:"stacks"`
}

// ListFleets handles GET /api/fleets
func ListFleets(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	rows, err := database.DB.Query(
		`SELECT id, player_id, name, formation, commander_id, targeting_command,
		        status, planet_id, position_x, position_y,
		        destination_x, destination_y, arrival_at, created_at, updated_at
		 FROM fleets WHERE player_id = $1
		 ORDER BY created_at`, playerID,
	)
	if err != nil {
		log.Printf("Failed to list fleets: %v", err)
		http.Error(w, `{"error":"failed to list fleets"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	fleets := []fleetWithStacks{}
	for rows.Next() {
		var f models.Fleet
		err := rows.Scan(
			&f.ID, &f.PlayerID, &f.Name, &f.Formation, &f.CommanderID, &f.TargetingCommand,
			&f.Status, &f.PlanetID, &f.PositionX, &f.PositionY,
			&f.DestinationX, &f.DestinationY, &f.ArrivalAt, &f.CreatedAt, &f.UpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan fleet: %v", err)
			continue
		}

		stacks := getFleetStacks(f.ID)
		fleets = append(fleets, fleetWithStacks{Fleet: f, Stacks: stacks})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fleets)
}

// CreateFleet handles POST /api/fleets
func CreateFleet(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	var req createFleetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		req.Name = "Fleet"
	}
	if req.Formation == "" {
		req.Formation = "phalanx"
	}
	if !validFormations[req.Formation] {
		http.Error(w, `{"error":"invalid formation"}`, http.StatusBadRequest)
		return
	}

	// Verify planet ownership
	if req.PlanetID != "" && !verifyPlanetOwnership(req.PlanetID, playerID) {
		errs.NotFound(fmt.Sprintf("Planet %s not found or does not belong to you", req.PlanetID)).WriteJSON(w, http.StatusNotFound)
		return
	}

	var f models.Fleet
	var planetID *string
	if req.PlanetID != "" {
		planetID = &req.PlanetID
	}

	err := database.DB.QueryRow(
		`INSERT INTO fleets (player_id, name, formation, planet_id)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, player_id, name, formation, commander_id, targeting_command,
		           status, planet_id, position_x, position_y,
		           destination_x, destination_y, arrival_at, created_at, updated_at`,
		playerID, req.Name, req.Formation, planetID,
	).Scan(
		&f.ID, &f.PlayerID, &f.Name, &f.Formation, &f.CommanderID, &f.TargetingCommand,
		&f.Status, &f.PlanetID, &f.PositionX, &f.PositionY,
		&f.DestinationX, &f.DestinationY, &f.ArrivalAt, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to create fleet: %v", err)
		http.Error(w, `{"error":"failed to create fleet"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fleetWithStacks{Fleet: f, Stacks: []models.FleetStack{}})
}

// UpdateFleet handles PUT /api/fleets/{id}
func UpdateFleet(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	fleetID := r.PathValue("id")

	fleet, err := getOwnedFleet(fleetID, playerID)
	if err != nil {
		http.Error(w, `{"error":"fleet not found"}`, http.StatusNotFound)
		return
	}

	if fleet.Status != "stationed" {
		http.Error(w, `{"error":"fleet must be stationed to modify"}`, http.StatusConflict)
		return
	}

	var req updateFleetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Formation != nil && !validFormations[*req.Formation] {
		http.Error(w, `{"error":"invalid formation"}`, http.StatusBadRequest)
		return
	}
	if req.TargetingCommand != nil && !validTargeting[*req.TargetingCommand] {
		http.Error(w, `{"error":"invalid targeting command"}`, http.StatusBadRequest)
		return
	}

	// Build dynamic update
	name := fleet.Name
	formation := fleet.Formation
	targeting := fleet.TargetingCommand
	commanderID := fleet.CommanderID

	if req.Name != nil {
		name = *req.Name
	}
	if req.Formation != nil {
		formation = *req.Formation
	}
	if req.TargetingCommand != nil {
		targeting = *req.TargetingCommand
	}
	if req.CommanderID != nil {
		if *req.CommanderID == "" {
			commanderID = nil
		} else {
			commanderID = req.CommanderID
		}
	}

	var f models.Fleet
	err = database.DB.QueryRow(
		`UPDATE fleets SET name = $1, formation = $2, targeting_command = $3, commander_id = $4, updated_at = now()
		 WHERE id = $5
		 RETURNING id, player_id, name, formation, commander_id, targeting_command,
		           status, planet_id, position_x, position_y,
		           destination_x, destination_y, arrival_at, created_at, updated_at`,
		name, formation, targeting, commanderID, fleetID,
	).Scan(
		&f.ID, &f.PlayerID, &f.Name, &f.Formation, &f.CommanderID, &f.TargetingCommand,
		&f.Status, &f.PlanetID, &f.PositionX, &f.PositionY,
		&f.DestinationX, &f.DestinationY, &f.ArrivalAt, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to update fleet: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	stacks := getFleetStacks(f.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fleetWithStacks{Fleet: f, Stacks: stacks})
}

// DeleteFleet handles DELETE /api/fleets/{id}
// Disbands a fleet; ships return to pool.
func DeleteFleet(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	fleetID := r.PathValue("id")

	fleet, err := getOwnedFleet(fleetID, playerID)
	if err != nil {
		http.Error(w, `{"error":"fleet not found"}`, http.StatusNotFound)
		return
	}

	if fleet.Status != "stationed" && fleet.Status != "dismissed" {
		http.Error(w, `{"error":"fleet must be stationed or dismissed to disband"}`, http.StatusConflict)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Return ships from stacks to pool
	_, err = tx.Exec(
		`UPDATE ships SET quantity = ships.quantity + fs.ship_count
		 FROM fleet_stacks fs
		 WHERE ships.ship_design_id = fs.ship_design_id
		   AND ships.player_id = $1
		   AND fs.fleet_id = $2`, playerID, fleetID,
	)
	if err != nil {
		log.Printf("Failed to return ships: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Delete stacks and fleet (cascade handles stacks)
	_, err = tx.Exec(`DELETE FROM fleets WHERE id = $1`, fleetID)
	if err != nil {
		log.Printf("Failed to delete fleet: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Disbanded bool `json:"disbanded"`
	}{Disbanded: true})
}

// AssignStack handles POST /api/fleets/{id}/assign-stack
func AssignStack(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	fleetID := r.PathValue("id")

	fleet, err := getOwnedFleet(fleetID, playerID)
	if err != nil {
		http.Error(w, `{"error":"fleet not found"}`, http.StatusNotFound)
		return
	}

	if fleet.Status != "stationed" {
		http.Error(w, `{"error":"fleet must be stationed"}`, http.StatusConflict)
		return
	}

	var req assignStackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.GridRow < 0 || req.GridRow > 2 || req.GridCol < 0 || req.GridCol > 2 {
		http.Error(w, `{"error":"grid position must be 0-2 for row and col"}`, http.StatusBadRequest)
		return
	}
	if req.ShipCount <= 0 || req.ShipCount > maxShipsPerStack {
		http.Error(w, `{"error":"ship count must be between 1 and 3000"}`, http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Check position is empty
	var positionOccupied bool
	err = tx.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM fleet_stacks WHERE fleet_id = $1 AND grid_row = $2 AND grid_col = $3)`,
		fleetID, req.GridRow, req.GridCol,
	).Scan(&positionOccupied)
	if err != nil {
		log.Printf("Failed to check position: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if positionOccupied {
		http.Error(w, `{"error":"grid position already occupied"}`, http.StatusConflict)
		return
	}

	// Check player has enough unassigned ships of this design
	var availableShips int
	err = tx.QueryRow(
		`SELECT COALESCE(quantity, 0) FROM ships
		 WHERE player_id = $1 AND ship_design_id = $2 AND is_building = false`,
		playerID, req.ShipDesignID,
	).Scan(&availableShips)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"no ships of this design"}`, http.StatusConflict)
			return
		}
		log.Printf("Failed to check ships: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if availableShips < req.ShipCount {
		http.Error(w, `{"error":"not enough unassigned ships"}`, http.StatusConflict)
		return
	}

	// Deduct ships from pool
	_, err = tx.Exec(
		`UPDATE ships SET quantity = quantity - $1, updated_at = now()
		 WHERE player_id = $2 AND ship_design_id = $3`,
		req.ShipCount, playerID, req.ShipDesignID,
	)
	if err != nil {
		log.Printf("Failed to deduct ships: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Create fleet stack
	var stack models.FleetStack
	err = tx.QueryRow(
		`INSERT INTO fleet_stacks (fleet_id, ship_design_id, grid_row, grid_col, ship_count)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, fleet_id, ship_design_id, grid_row, grid_col, ship_count`,
		fleetID, req.ShipDesignID, req.GridRow, req.GridCol, req.ShipCount,
	).Scan(&stack.ID, &stack.FleetID, &stack.ShipDesignID, &stack.GridRow, &stack.GridCol, &stack.ShipCount)
	if err != nil {
		log.Printf("Failed to create fleet stack: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	remaining := availableShips - req.ShipCount

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Stack          models.FleetStack `json:"stack"`
		AvailableShips int               `json:"available_ships"`
	}{Stack: stack, AvailableShips: remaining})
}

// RemoveStack handles POST /api/fleets/{id}/remove-stack
func RemoveStack(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	fleetID := r.PathValue("id")

	fleet, err := getOwnedFleet(fleetID, playerID)
	if err != nil {
		http.Error(w, `{"error":"fleet not found"}`, http.StatusNotFound)
		return
	}

	if fleet.Status != "stationed" {
		http.Error(w, `{"error":"fleet must be stationed"}`, http.StatusConflict)
		return
	}

	var req removeStackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Get stack info
	var stack models.FleetStack
	err = tx.QueryRow(
		`SELECT id, fleet_id, ship_design_id, grid_row, grid_col, ship_count
		 FROM fleet_stacks WHERE fleet_id = $1 AND grid_row = $2 AND grid_col = $3`,
		fleetID, req.GridRow, req.GridCol,
	).Scan(&stack.ID, &stack.FleetID, &stack.ShipDesignID, &stack.GridRow, &stack.GridCol, &stack.ShipCount)
	if err != nil {
		http.Error(w, `{"error":"no stack at this position"}`, http.StatusNotFound)
		return
	}

	// Return ships to pool
	_, err = tx.Exec(
		`UPDATE ships SET quantity = quantity + $1, updated_at = now()
		 WHERE player_id = $2 AND ship_design_id = $3`,
		stack.ShipCount, playerID, stack.ShipDesignID,
	)
	if err != nil {
		log.Printf("Failed to return ships: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Delete stack
	_, err = tx.Exec(`DELETE FROM fleet_stacks WHERE id = $1`, stack.ID)
	if err != nil {
		log.Printf("Failed to delete stack: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Removed bool `json:"removed"`
	}{Removed: true})
}

// MoveFleet handles POST /api/fleets/{id}/move
func MoveFleet(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	fleetID := r.PathValue("id")

	fleet, err := getOwnedFleet(fleetID, playerID)
	if err != nil {
		http.Error(w, `{"error":"fleet not found"}`, http.StatusNotFound)
		return
	}

	if fleet.Status != "stationed" {
		http.Error(w, `{"error":"fleet must be stationed to move"}`, http.StatusConflict)
		return
	}

	// Check fleet has at least 1 stack
	stacks := getFleetStacks(fleetID)
	if len(stacks) == 0 {
		http.Error(w, `{"error":"fleet has no ships"}`, http.StatusConflict)
		return
	}

	var req moveFleetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// For Phase 2, travel is instant (stub for future implementation)
	var f models.Fleet
	err = database.DB.QueryRow(
		`UPDATE fleets SET status = 'traveling', destination_x = $1, destination_y = $2,
		        arrival_at = now() + interval '1 second', updated_at = now()
		 WHERE id = $3
		 RETURNING id, player_id, name, formation, commander_id, targeting_command,
		           status, planet_id, position_x, position_y,
		           destination_x, destination_y, arrival_at, created_at, updated_at`,
		req.DestinationX, req.DestinationY, fleetID,
	).Scan(
		&f.ID, &f.PlayerID, &f.Name, &f.Formation, &f.CommanderID, &f.TargetingCommand,
		&f.Status, &f.PlanetID, &f.PositionX, &f.PositionY,
		&f.DestinationX, &f.DestinationY, &f.ArrivalAt, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to move fleet: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fleetWithStacks{Fleet: f, Stacks: stacks})
}

// RecallFleet handles POST /api/fleets/{id}/recall
func RecallFleet(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	fleetID := r.PathValue("id")

	fleet, err := getOwnedFleet(fleetID, playerID)
	if err != nil {
		http.Error(w, `{"error":"fleet not found"}`, http.StatusNotFound)
		return
	}

	if fleet.Status != "traveling" {
		http.Error(w, `{"error":"fleet must be traveling to recall"}`, http.StatusConflict)
		return
	}

	var f models.Fleet
	err = database.DB.QueryRow(
		`UPDATE fleets SET status = 'returning', arrival_at = now() + interval '1 second', updated_at = now()
		 WHERE id = $1
		 RETURNING id, player_id, name, formation, commander_id, targeting_command,
		           status, planet_id, position_x, position_y,
		           destination_x, destination_y, arrival_at, created_at, updated_at`,
		fleetID,
	).Scan(
		&f.ID, &f.PlayerID, &f.Name, &f.Formation, &f.CommanderID, &f.TargetingCommand,
		&f.Status, &f.PlanetID, &f.PositionX, &f.PositionY,
		&f.DestinationX, &f.DestinationY, &f.ArrivalAt, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to recall fleet: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	stacks := getFleetStacks(f.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fleetWithStacks{Fleet: f, Stacks: stacks})
}

// DismissFleet handles POST /api/fleets/{id}/dismiss
func DismissFleet(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	fleetID := r.PathValue("id")

	fleet, err := getOwnedFleet(fleetID, playerID)
	if err != nil {
		http.Error(w, `{"error":"fleet not found"}`, http.StatusNotFound)
		return
	}

	if fleet.Status != "stationed" {
		http.Error(w, `{"error":"fleet must be stationed to dismiss"}`, http.StatusConflict)
		return
	}

	var f models.Fleet
	err = database.DB.QueryRow(
		`UPDATE fleets SET status = 'dismissed', updated_at = now()
		 WHERE id = $1
		 RETURNING id, player_id, name, formation, commander_id, targeting_command,
		           status, planet_id, position_x, position_y,
		           destination_x, destination_y, arrival_at, created_at, updated_at`,
		fleetID,
	).Scan(
		&f.ID, &f.PlayerID, &f.Name, &f.Formation, &f.CommanderID, &f.TargetingCommand,
		&f.Status, &f.PlanetID, &f.PositionX, &f.PositionY,
		&f.DestinationX, &f.DestinationY, &f.ArrivalAt, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to dismiss fleet: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	stacks := getFleetStacks(f.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fleetWithStacks{Fleet: f, Stacks: stacks})
}

// Helper functions

func getOwnedFleet(fleetID, playerID string) (models.Fleet, error) {
	var f models.Fleet
	err := database.DB.QueryRow(
		`SELECT id, player_id, name, formation, commander_id, targeting_command,
		        status, planet_id, position_x, position_y,
		        destination_x, destination_y, arrival_at, created_at, updated_at
		 FROM fleets WHERE id = $1 AND player_id = $2`, fleetID, playerID,
	).Scan(
		&f.ID, &f.PlayerID, &f.Name, &f.Formation, &f.CommanderID, &f.TargetingCommand,
		&f.Status, &f.PlanetID, &f.PositionX, &f.PositionY,
		&f.DestinationX, &f.DestinationY, &f.ArrivalAt, &f.CreatedAt, &f.UpdatedAt,
	)
	return f, err
}

func getFleetStacks(fleetID string) []models.FleetStack {
	rows, err := database.DB.Query(
		`SELECT id, fleet_id, ship_design_id, grid_row, grid_col, ship_count
		 FROM fleet_stacks WHERE fleet_id = $1
		 ORDER BY grid_row, grid_col`, fleetID,
	)
	if err != nil {
		log.Printf("Failed to get fleet stacks: %v", err)
		return []models.FleetStack{}
	}
	defer rows.Close()

	stacks := []models.FleetStack{}
	for rows.Next() {
		var s models.FleetStack
		if err := rows.Scan(&s.ID, &s.FleetID, &s.ShipDesignID, &s.GridRow, &s.GridCol, &s.ShipCount); err != nil {
			continue
		}
		stacks = append(stacks, s)
	}
	return stacks
}
