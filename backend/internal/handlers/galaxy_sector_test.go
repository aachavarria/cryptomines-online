package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
)

func setupGalaxyTestDB(t *testing.T) {
	if database.DB == nil {
		if err := database.InitSupabase(); err != nil {
			t.Fatalf("Failed to init database: %v", err)
		}
	}
}

// TestGetGalaxySector_ReturnsRBPsAndOwnPlanet verifies the sector endpoint
// returns RBPs (already seeded in 49-cell grid around 1000,1000) and flags
// the caller's own planet via is_own=true.
func TestGetGalaxySector_ReturnsRBPsAndOwnPlanet(t *testing.T) {
	setupGalaxyTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000801"
	planetID := "00000000-0000-0000-aaaa-000000000801"
	defer func() {
		database.DB.Exec("DELETE FROM resources WHERE planet_id = $1", planetID)
		database.DB.Exec("DELETE FROM planets WHERE id = $1", planetID)
		database.DB.Exec("DELETE FROM players WHERE id = $1", playerID)
	}()

	if _, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'galaxytest', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID); err != nil {
		t.Fatalf("create player: %v", err)
	}
	// Place inside the seeded 7x7 RBP region (1000..1060).
	if _, err := database.DB.Exec(`
		INSERT INTO planets (id, player_id, name, is_homeworld, position_x, position_y, created_at, updated_at)
		VALUES ($1, $2, 'GalaxyHome', true, 1015, 1015, now(), now())
		ON CONFLICT (id) DO NOTHING
	`, planetID, playerID); err != nil {
		t.Fatalf("create planet: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/galaxy/sector?cx=1015&cy=1015&r=50", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	GetGalaxySector(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}

	var body struct {
		CenterX int            `json:"center_x"`
		CenterY int            `json:"center_y"`
		Radius  int            `json:"radius"`
		Planets []SectorPlanet `json:"planets"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v body=%s", err, rr.Body.String())
	}
	if body.CenterX != 1015 || body.CenterY != 1015 || body.Radius != 50 {
		t.Errorf("center/radius wrong: %+v", body)
	}

	var ownFound, rbpFound bool
	for _, p := range body.Planets {
		if p.ID == planetID {
			if !p.IsOwn {
				t.Errorf("own planet not flagged: %+v", p)
			}
			if !p.IsHomeworld {
				t.Errorf("own homeworld flag missing: %+v", p)
			}
			ownFound = true
		}
		if p.IsRBP {
			rbpFound = true
		}
	}
	if !ownFound {
		t.Errorf("own planet missing from sector results")
	}
	if !rbpFound {
		t.Errorf("expected at least one RBP in this sector (RBPs seeded at 1000..1060)")
	}
}

func TestGetGalaxySector_RadiusClamped(t *testing.T) {
	setupGalaxyTestDB(t)

	playerID := "00000000-0000-0000-0000-000000000802"
	defer func() { database.DB.Exec("DELETE FROM players WHERE id = $1", playerID) }()
	database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, created_at, updated_at)
		VALUES ($1, $2, 'rclamp', now(), now())
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID)

	req := httptest.NewRequest(http.MethodGet, "/api/galaxy/sector?cx=1000&cy=1000&r=99999", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	GetGalaxySector(rr, req)

	var body struct {
		Radius int `json:"radius"`
	}
	json.Unmarshal(rr.Body.Bytes(), &body)
	if body.Radius != 1000 {
		t.Errorf("radius should be clamped to 1000, got %d", body.Radius)
	}
}
