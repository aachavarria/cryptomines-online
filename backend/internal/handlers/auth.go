package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
)

type guestResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int         `json:"expires_in"`
	Player       guestPlayer `json:"player"`
}

type guestPlayer struct {
	ID          string    `json:"id"`
	AnonymousID string    `json:"anonymous_id"`
	Level       int       `json:"level"`
	CreatedAt   time.Time `json:"created_at"`
}

// GuestAuth handles POST /api/auth/guest
// Creates an anonymous user via Supabase Auth, then sets up player data.
func GuestAuth(w http.ResponseWriter, r *http.Request) {
	// 1. Call Supabase Auth to create anonymous user
	authResp, err := database.SignUpAnonymous()
	if err != nil {
		log.Printf("Supabase Auth signup failed: %v", err)
		// Fallback: create player with local anonymous ID if Supabase Auth is unavailable
		handleFallbackGuestAuth(w, r)
		return
	}

	supabaseUserID := authResp.User.ID
	accessToken := authResp.AccessToken
	anonID := "supabase_" + supabaseUserID

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// 2. Create player using Supabase auth user ID as the player ID
	var playerID string
	var createdAt time.Time
	err = tx.QueryRow(
		`INSERT INTO players (id, anonymous_id) VALUES ($1, $2) RETURNING id, created_at`,
		supabaseUserID, anonID,
	).Scan(&playerID, &createdAt)
	if err != nil {
		log.Printf("Failed to create player: %v", err)
		http.Error(w, `{"error":"failed to create player"}`, http.StatusInternalServerError)
		return
	}

	// 3. Create homeworld
	if err := createHomeworld(tx, playerID); err != nil {
		log.Printf("Failed to create homeworld: %v", err)
		http.Error(w, `{"error":"failed to create homeworld"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Use the Supabase access_token
	token := accessToken
	if token == "" {
		log.Printf("Supabase Auth did not return access token")
		http.Error(w, `{"error":"authentication failed"}`, http.StatusInternalServerError)
		return
	}

	resp := guestResponse{
		Token:        token,
		RefreshToken: authResp.RefreshToken,
		ExpiresIn:    authResp.ExpiresIn,
		Player: guestPlayer{
			ID:          playerID,
			AnonymousID: anonID,
			Level:       1,
			CreatedAt:   createdAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// handleFallbackGuestAuth creates a guest without Supabase Auth (for local dev without Auth service).
func handleFallbackGuestAuth(w http.ResponseWriter, r *http.Request) {
	anonID := generateAnonymousID()

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var playerID string
	var createdAt time.Time
	err = tx.QueryRow(
		`INSERT INTO players (anonymous_id) VALUES ($1) RETURNING id, created_at`,
		anonID,
	).Scan(&playerID, &createdAt)
	if err != nil {
		log.Printf("Failed to create player: %v", err)
		http.Error(w, `{"error":"failed to create player"}`, http.StatusInternalServerError)
		return
	}

	if err := createHomeworld(tx, playerID); err != nil {
		log.Printf("Failed to create homeworld: %v", err)
		http.Error(w, `{"error":"failed to create homeworld"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Try to create Supabase Auth user for this player
	authResp, err := database.SignUpAnonymous()
	var token string
	if err != nil {
		log.Printf("Warning: Could not create Supabase Auth user in fallback: %v", err)
		http.Error(w, `{"error":"authentication service unavailable"}`, http.StatusServiceUnavailable)
		return
	} else {
		token = authResp.AccessToken
		// Update player ID to match Supabase user ID
		_, err = database.DB.Exec(`UPDATE players SET id = $1 WHERE id = $2`, authResp.User.ID, playerID)
		if err != nil {
			log.Printf("Failed to update player ID: %v", err)
		}
		playerID = authResp.User.ID
	}

	resp := guestResponse{
		Token:        token,
		RefreshToken: authResp.RefreshToken,
		ExpiresIn:    authResp.ExpiresIn,
		Player: guestPlayer{
			ID:          playerID,
			AnonymousID: anonID,
			Level:       1,
			CreatedAt:   createdAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// createHomeworld creates the homeworld planet, initial buildings, and resources
// for a new player within a transaction.
func createHomeworld(tx *sql.Tx, playerID string) error {
	// Find random unoccupied position
	posX, posY, err := findRandomPosition(tx)
	if err != nil {
		return err
	}

	var planetID string
	err = tx.QueryRow(
		`INSERT INTO planets (player_id, name, position_x, position_y, is_homeworld)
		 VALUES ($1, 'Homeworld', $2, $3, true) RETURNING id`,
		playerID, posX, posY,
	).Scan(&planetID)
	if err != nil {
		return err
	}

	// Create initial buildings (all at Lv1, not upgrading)
	// Grid positions chosen so buildings don't overlap on the 16x16 grid.
	// Sizes from buildingConfig.ts: civic_center 3x3, metal_collector 2x2,
	// he3_extractor 2x2, residential_area 2x2, resource_warehouse 3x2, space_station 3x3.
	type initBuilding struct {
		name    string
		gridCol int
		gridRow int
	}
	initialBuildings := []initBuilding{
		{"civic_center", 5, 5},
		{"metal_collector", 2, 2},
		{"he3_extractor", 9, 2},
		{"residential_area", 2, 9},
		{"resource_warehouse", 9, 9},
		{"space_station", 5, 1},
	}
	for _, b := range initialBuildings {
		var btID int
		err = tx.QueryRow(`SELECT id FROM building_types WHERE name = $1`, b.name).Scan(&btID)
		if err != nil {
			return err
		}
		_, err = tx.Exec(
			`INSERT INTO buildings (planet_id, building_type, level, grid_col, grid_row) VALUES ($1, $2, 1, $3, $4)`,
			planetID, btID, b.gridCol, b.gridRow,
		)
		if err != nil {
			return err
		}
	}

	// Initial production rates from wiki Lv1 data:
	// Metal Collector Lv1: 1080/hr, He3 Extractor Lv1: 1180/hr, Residential Area Lv1: 1300/hr
	// Resource Warehouse Lv1 storage: 10000
	_, err = tx.Exec(
		`INSERT INTO resources (planet_id, metal, he3, gold, metal_per_hour, he3_per_hour, gold_per_hour, storage_capacity)
		 VALUES ($1, 5000, 5000, 10000, 1080, 1180, 1300, 10000)`,
		planetID,
	)
	return err
}

func generateAnonymousID() string {
	b := make([]byte, 12)
	rand.Read(b)
	return "guest_" + hex.EncodeToString(b)
}

func findRandomPosition(tx *sql.Tx) (int, int, error) {
	maxAttempts := 100
	for i := 0; i < maxAttempts; i++ {
		x, _ := rand.Int(rand.Reader, big.NewInt(10000))
		y, _ := rand.Int(rand.Reader, big.NewInt(10000))
		posX := int(x.Int64()) + 1
		posY := int(y.Int64()) + 1

		var exists bool
		err := tx.QueryRow(
			`SELECT EXISTS(SELECT 1 FROM planets WHERE position_x = $1 AND position_y = $2)`,
			posX, posY,
		).Scan(&exists)
		if err != nil {
			return 0, 0, err
		}
		if !exists {
			return posX, posY, nil
		}
	}
	return 0, 0, sql.ErrNoRows
}
