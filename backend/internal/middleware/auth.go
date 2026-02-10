package middleware

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const PlayerIDKey contextKey = "player_id"

// Auth is a middleware that validates JWT tokens from the Authorization header.
// Uses Supabase JWT secret for token validation.
// If the player doesn't exist in the database, it creates them automatically.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
			return
		}

		// Get Supabase JWT secret
		supabaseSecret := os.Getenv("SUPABASE_JWT_SECRET")
		if supabaseSecret == "" {
			// Default Supabase local dev JWT secret
			supabaseSecret = "super-secret-jwt-token-with-at-least-32-characters-long"
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(supabaseSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error":"invalid token claims"}`, http.StatusUnauthorized)
			return
		}

		playerID, ok := claims["sub"].(string)
		if !ok || playerID == "" {
			http.Error(w, `{"error":"invalid token subject"}`, http.StatusUnauthorized)
			return
		}

		// Ensure player exists in database
		if err := ensurePlayerExists(playerID); err != nil {
			log.Printf("Failed to ensure player exists: %v", err)
			http.Error(w, `{"error":"failed to initialize player"}`, http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), PlayerIDKey, playerID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetPlayerID extracts the player ID from the request context.
func GetPlayerID(r *http.Request) string {
	if id, ok := r.Context().Value(PlayerIDKey).(string); ok {
		return id
	}
	return ""
}

// ensurePlayerExists checks if a player exists in the database.
// If not, creates the player with their homeworld and initial resources.
// This is idempotent - safe to call multiple times.
func ensurePlayerExists(playerID string) error {
	// Check if player already exists
	var exists bool
	err := database.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM players WHERE id = $1)`, playerID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return nil // Player already exists
	}

	// Create player with homeworld in a transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create player
	anonID := "supabase_" + playerID
	_, err = tx.Exec(`INSERT INTO players (id, anonymous_id) VALUES ($1, $2)`, playerID, anonID)
	if err != nil {
		return err
	}

	// Create homeworld
	if err := createHomeworld(tx, playerID); err != nil {
		return err
	}

	return tx.Commit()
}

// createHomeworld creates the initial homeworld planet for a new player
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

	// Create initial buildings
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

	// Initial resources
	_, err = tx.Exec(
		`INSERT INTO resources (planet_id, metal, he3, gold, metal_per_hour, he3_per_hour, gold_per_hour, storage_capacity)
		 VALUES ($1, 5000, 5000, 10000, 1080, 1180, 1300, 10000)`,
		planetID,
	)
	return err
}

// findRandomPosition finds an unoccupied position on the map
func findRandomPosition(tx *sql.Tx) (int, int, error) {
	// Simple strategy: use player_id hash to generate deterministic but distributed positions
	// In production, you'd use a proper spatial algorithm
	var count int
	err := tx.QueryRow(`SELECT COUNT(*) FROM planets`).Scan(&count)
	if err != nil {
		return 0, 0, err
	}

	// Spread players in a grid pattern
	gridSize := 100
	x := (count % gridSize) * 100 + 50
	y := (count / gridSize) * 100 + 50

	return x, y, nil
}
