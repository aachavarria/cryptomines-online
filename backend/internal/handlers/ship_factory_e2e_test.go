package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
)

// =============================================================================
// Ship Factory E2E Tests
//
// These tests cover the full HTTP-level flow for ship factory build/cancel,
// ship_designs CRUD, and tier gating. They write real rows to the test DB
// and clean up after themselves.
//
// Player UUID range: 5000-6999 (per test plan, no collisions with other suites).
// =============================================================================

// setupShipFactoryE2E initializes the test DB.
func setupShipFactoryE2E(t *testing.T) {
	t.Helper()
	if database.DB == nil {
		if err := database.InitSupabase(); err != nil {
			t.Fatalf("Failed to init database: %v", err)
		}
	}
}

// cleanupShipFactoryE2EPlayer removes all dependent rows for a player in
// the correct order (children before parents).
func cleanupShipFactoryE2EPlayer(t *testing.T, playerID string) {
	t.Helper()
	// Children of fleets/ships first, then their parents, then planet & player
	database.DB.Exec("DELETE FROM ship_instances WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM fleet_stacks WHERE fleet_id IN (SELECT id FROM fleets WHERE player_id = $1)", playerID)
	database.DB.Exec("DELETE FROM fleets WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM ships WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM ship_designs WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM player_blueprints WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM technologies WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM buildings WHERE planet_id IN (SELECT id FROM planets WHERE player_id = $1)", playerID)
	database.DB.Exec("DELETE FROM resources WHERE planet_id IN (SELECT id FROM planets WHERE player_id = $1)", playerID)
	database.DB.Exec("DELETE FROM planets WHERE player_id = $1", playerID)
	database.DB.Exec("DELETE FROM players WHERE id = $1", playerID)
}

// shipFactoryE2EFixture seeds a player + homeworld + resources, plus optional
// Ship Factory + Civic Center buildings, plus the Weikes hull blueprint.
// Returns the player_id, planet_id, hull_type_id, and tier 1 module_type_id
// (with its blueprint already activated for the player).
type shipFactoryE2EFixture struct {
	PlayerID         string
	PlanetID         string
	HullTypeID       int
	ModuleTypeID     int
	HullBlueprintID  int
	ModuleBlueprintID int
}

// seedShipFactoryE2EFixture inserts the baseline player+planet+resources
// (NOTE: resources table has no created_at — only updated_at and last_collected_at).
func seedShipFactoryE2EFixture(t *testing.T, playerID, planetID string, opts struct {
	WithShipFactory      bool
	ShipFactoryLevel     int
	WithResources        bool
	Metal, He3, Gold     int64
	WeikesBlueprintLevel int // 0 = don't give blueprint
	ModuleBlueprintLevel int // 0 = don't give module blueprint
}) shipFactoryE2EFixture {
	t.Helper()

	_, err := database.DB.Exec(`
		INSERT INTO players (id, anonymous_id, username, level)
		VALUES ($1, $2, $3, 50)
		ON CONFLICT (id) DO NOTHING
	`, playerID, "anon-"+playerID, "sf-e2e-"+playerID[len(playerID)-12:])
	if err != nil {
		t.Fatalf("create player: %v", err)
	}

	// Use playerID-derived position to avoid uq_planet_position collision.
	// Take last 4 hex chars and convert to int -> position offset.
	posSeed, _ := strconv.ParseInt(playerID[len(playerID)-4:], 16, 64)
	posX := 100 + int(posSeed%500)
	posY := 100 + int((posSeed/500)%500)
	_, err = database.DB.Exec(`
		INSERT INTO planets (id, player_id, name, is_homeworld, position_x, position_y)
		VALUES ($1, $2, 'SF E2E Planet', true, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`, planetID, playerID, posX, posY)
	if err != nil {
		t.Fatalf("create planet: %v", err)
	}

	if opts.WithResources {
		// resources has no created_at, only updated_at + last_collected_at
		_, err = database.DB.Exec(`
			INSERT INTO resources (planet_id, metal, he3, gold)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (planet_id) DO UPDATE
			SET metal = EXCLUDED.metal, he3 = EXCLUDED.he3, gold = EXCLUDED.gold,
			    updated_at = now()
		`, planetID, opts.Metal, opts.He3, opts.Gold)
		if err != nil {
			t.Fatalf("create resources: %v", err)
		}
	}

	if opts.WithShipFactory {
		// Need a Civic Center first because ship_factory has prerequisite_building='civic_center'
		// (we just need the row to exist for the JOIN — building costs are not checked at HTTP layer here)
		var civicCenterID, shipFactoryID int
		database.DB.QueryRow(`SELECT id FROM building_types WHERE name = 'civic_center'`).Scan(&civicCenterID)
		database.DB.QueryRow(`SELECT id FROM building_types WHERE name = 'ship_factory'`).Scan(&shipFactoryID)
		// buildings.building_type is the FK column (NOT building_type_id)
		database.DB.Exec(`
			INSERT INTO buildings (planet_id, building_type, grid_col, grid_row, level, is_upgrading)
			VALUES ($1, $2, 0, 0, 5, false)
			ON CONFLICT DO NOTHING
		`, planetID, civicCenterID)
		database.DB.Exec(`
			INSERT INTO buildings (planet_id, building_type, grid_col, grid_row, level, is_upgrading)
			VALUES ($1, $2, 1, 0, $3, false)
			ON CONFLICT DO NOTHING
		`, planetID, shipFactoryID, opts.ShipFactoryLevel)
	}

	fix := shipFactoryE2EFixture{
		PlayerID: playerID,
		PlanetID: planetID,
	}

	// Look up Weikes hull and module IDs for tests that want them
	database.DB.QueryRow(`SELECT id FROM hull_types WHERE name = 'weikes_i'`).Scan(&fix.HullTypeID)
	database.DB.QueryRow(`
		SELECT id FROM blueprints
		WHERE blueprint_type = 'hull' AND hull_type_id = (SELECT id FROM hull_types WHERE name = 'weikes_i')
	`).Scan(&fix.HullBlueprintID)
	database.DB.QueryRow(`SELECT id FROM module_types WHERE tier = 1 LIMIT 1`).Scan(&fix.ModuleTypeID)
	database.DB.QueryRow(`
		SELECT b.id FROM blueprints b
		JOIN module_types mt ON b.module_type_id = mt.id
		WHERE b.blueprint_type = 'module' AND mt.id = $1
	`, fix.ModuleTypeID).Scan(&fix.ModuleBlueprintID)

	if opts.WeikesBlueprintLevel > 0 && fix.HullBlueprintID > 0 {
		database.DB.Exec(`
			INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
			VALUES ($1, $2, true, $3)
			ON CONFLICT (player_id, blueprint_id) DO UPDATE
			SET is_activated = true, research_level = EXCLUDED.research_level
		`, playerID, fix.HullBlueprintID, opts.WeikesBlueprintLevel)
	}

	if opts.ModuleBlueprintLevel > 0 && fix.ModuleBlueprintID > 0 {
		database.DB.Exec(`
			INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
			VALUES ($1, $2, true, $3)
			ON CONFLICT (player_id, blueprint_id) DO UPDATE
			SET is_activated = true, research_level = EXCLUDED.research_level
		`, playerID, fix.ModuleBlueprintID, opts.ModuleBlueprintLevel)
	}

	return fix
}

// makeShipDesignDirectly inserts a ship_design row directly (bypasses CreateShipDesign
// validation) so build tests can use a known design ID without exercising tier checks.
func makeShipDesignDirectly(t *testing.T, playerID string, hullTypeID int, name string, metalCost, he3Cost, goldCost int64) string {
	t.Helper()
	var designID string
	err := database.DB.QueryRow(`
		INSERT INTO ship_designs (
			player_id, name, hull_type_id, modules_json,
			total_shield, total_structure, total_defense, total_agility,
			total_movement, total_storage, attack_power,
			weapon_range_min, weapon_range_max, volume_used,
			he3_per_round, metal_cost, he3_cost, gold_cost, build_time_seconds
		) VALUES ($1, $2, $3, '[]', 100, 100, 10, 10, 5, 100, 50, 1, 5, 10, 5, $4, $5, $6, 60)
		RETURNING id
	`, playerID, name, hullTypeID, metalCost, he3Cost, goldCost).Scan(&designID)
	if err != nil {
		t.Fatalf("insert ship_design: %v", err)
	}
	return designID
}

// callBuildShips invokes BuildShips with proper auth context; returns code and parsed body.
func callBuildShips(t *testing.T, playerID, designID string, qty, slot int) (int, map[string]any) {
	t.Helper()
	body, _ := json.Marshal(buildRequest{
		ShipDesignID:   designID,
		Quantity:       qty,
		ProductionSlot: slot,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/ship-factory/build", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	BuildShips(rr, req)
	var parsed map[string]any
	if rr.Body.Len() > 0 {
		_ = json.Unmarshal(rr.Body.Bytes(), &parsed)
	}
	return rr.Code, parsed
}

// =============================================================================
// 1a. Full build flow: factory + resources + design -> BuildShips deducts and queues
// =============================================================================
func TestShipFactoryE2E_BuildFlow_Success(t *testing.T) {
	setupShipFactoryE2E(t)
	playerID := "00000000-0000-0000-0000-000000005001"
	planetID := "00000000-0000-0000-0001-000000005001"
	defer cleanupShipFactoryE2EPlayer(t, playerID)

	fix := seedShipFactoryE2EFixture(t, playerID, planetID, struct {
		WithShipFactory      bool
		ShipFactoryLevel     int
		WithResources        bool
		Metal, He3, Gold     int64
		WeikesBlueprintLevel int
		ModuleBlueprintLevel int
	}{
		WithShipFactory: true, ShipFactoryLevel: 1,
		WithResources: true, Metal: 100000, He3: 100000, Gold: 100000,
		WeikesBlueprintLevel: 1, ModuleBlueprintLevel: 1,
	})

	// Ship_design: cost 1000/500/200 per ship.
	designID := makeShipDesignDirectly(t, fix.PlayerID, fix.HullTypeID, "weikes-build", 1000, 500, 200)

	// Build 2 ships
	code, body := callBuildShips(t, fix.PlayerID, designID, 2, 1)
	if code != http.StatusOK {
		t.Fatalf("expected 200 from BuildShips, got %d: %v", code, body)
	}

	// Verify a row in ships table marks the build as in-progress in slot 1
	var buildQty, slot int
	var isBuilding bool
	err := database.DB.QueryRow(`
		SELECT build_quantity, is_building, production_slot FROM ships
		WHERE player_id = $1 AND ship_design_id = $2
	`, fix.PlayerID, designID).Scan(&buildQty, &isBuilding, &slot)
	if err != nil {
		t.Fatalf("expected ships row to exist: %v", err)
	}
	if !isBuilding {
		t.Errorf("expected ships.is_building=true, got false")
	}
	if buildQty != 2 {
		t.Errorf("expected build_quantity=2, got %d", buildQty)
	}
	if slot != 1 {
		t.Errorf("expected production_slot=1, got %d", slot)
	}

	// Verify resources were deducted: started with 100000 metal,
	// build cost = ShipBuildCost(1000, 500, 200, 2). Don't compute exact amount —
	// just assert a deduction happened.
	var metal, he3, gold int64
	database.DB.QueryRow(`SELECT metal, he3, gold FROM resources WHERE planet_id=$1`, fix.PlanetID).
		Scan(&metal, &he3, &gold)
	if metal >= 100000 {
		t.Errorf("expected metal deducted from 100000, got %d", metal)
	}
	if he3 >= 100000 {
		t.Errorf("expected he3 deducted from 100000, got %d", he3)
	}
	if gold >= 100000 {
		t.Errorf("expected gold deducted from 100000, got %d", gold)
	}
}

// =============================================================================
// 1b. Build rejected when no Ship Factory exists
// =============================================================================
func TestShipFactoryE2E_BuildRejected_NoFactory(t *testing.T) {
	setupShipFactoryE2E(t)
	playerID := "00000000-0000-0000-0000-000000005002"
	planetID := "00000000-0000-0000-0001-000000005002"
	defer cleanupShipFactoryE2EPlayer(t, playerID)

	fix := seedShipFactoryE2EFixture(t, playerID, planetID, struct {
		WithShipFactory      bool
		ShipFactoryLevel     int
		WithResources        bool
		Metal, He3, Gold     int64
		WeikesBlueprintLevel int
		ModuleBlueprintLevel int
	}{
		WithShipFactory: false,
		WithResources:   true, Metal: 100000, He3: 100000, Gold: 100000,
		WeikesBlueprintLevel: 1, ModuleBlueprintLevel: 1,
	})

	designID := makeShipDesignDirectly(t, fix.PlayerID, fix.HullTypeID, "noFactory", 1000, 500, 200)

	code, _ := callBuildShips(t, fix.PlayerID, designID, 2, 1)
	if code != http.StatusNotFound {
		t.Errorf("expected 404 (ship factory not built), got %d", code)
	}
}

// =============================================================================
// 1c. Build rejected when ship_design references hull not unlocked by player
// =============================================================================
func TestShipFactoryE2E_BuildRejected_DesignNotOwned(t *testing.T) {
	setupShipFactoryE2E(t)
	playerID := "00000000-0000-0000-0000-000000005003"
	otherPlayerID := "00000000-0000-0000-0000-000000005103"
	planetID := "00000000-0000-0000-0001-000000005003"
	otherPlanetID := "00000000-0000-0000-0001-000000005103"
	defer cleanupShipFactoryE2EPlayer(t, playerID)
	defer cleanupShipFactoryE2EPlayer(t, otherPlayerID)

	fix := seedShipFactoryE2EFixture(t, playerID, planetID, struct {
		WithShipFactory      bool
		ShipFactoryLevel     int
		WithResources        bool
		Metal, He3, Gold     int64
		WeikesBlueprintLevel int
		ModuleBlueprintLevel int
	}{
		WithShipFactory: true, ShipFactoryLevel: 1,
		WithResources: true, Metal: 100000, He3: 100000, Gold: 100000,
	})

	// The OTHER player owns the design (different player_id).
	otherFix := seedShipFactoryE2EFixture(t, otherPlayerID, otherPlanetID, struct {
		WithShipFactory      bool
		ShipFactoryLevel     int
		WithResources        bool
		Metal, He3, Gold     int64
		WeikesBlueprintLevel int
		ModuleBlueprintLevel int
	}{
		WeikesBlueprintLevel: 1,
	})

	otherDesignID := makeShipDesignDirectly(t, otherFix.PlayerID, otherFix.HullTypeID, "other-design", 1000, 500, 200)

	// Player A tries to build with player B's design ID — must 404.
	code, _ := callBuildShips(t, fix.PlayerID, otherDesignID, 1, 1)
	if code != http.StatusNotFound {
		t.Errorf("expected 404 (ship design not found), got %d", code)
	}
}

// =============================================================================
// 1d. Build rejected when player has insufficient resources
// =============================================================================
func TestShipFactoryE2E_BuildRejected_InsufficientResources(t *testing.T) {
	setupShipFactoryE2E(t)
	playerID := "00000000-0000-0000-0000-000000005004"
	planetID := "00000000-0000-0000-0001-000000005004"
	defer cleanupShipFactoryE2EPlayer(t, playerID)

	fix := seedShipFactoryE2EFixture(t, playerID, planetID, struct {
		WithShipFactory      bool
		ShipFactoryLevel     int
		WithResources        bool
		Metal, He3, Gold     int64
		WeikesBlueprintLevel int
		ModuleBlueprintLevel int
	}{
		WithShipFactory: true, ShipFactoryLevel: 1,
		WithResources: true, Metal: 10, He3: 10, Gold: 10, // way too few
		WeikesBlueprintLevel: 1, ModuleBlueprintLevel: 1,
	})

	designID := makeShipDesignDirectly(t, fix.PlayerID, fix.HullTypeID, "broke", 100000, 100000, 100000)

	code, body := callBuildShips(t, fix.PlayerID, designID, 5, 1)
	if code != http.StatusConflict {
		t.Errorf("expected 409 (insufficient resources), got %d: %v", code, body)
	}
}

// =============================================================================
// 1e. CancelShipBuild clears the slot
//
// Note: per current handler, cancel does NOT refund (matches GO2 mechanics).
// Test verifies the slot is cleared so a new build can take it.
// =============================================================================
func TestShipFactoryE2E_CancelShipBuild_ClearsSlot(t *testing.T) {
	setupShipFactoryE2E(t)
	playerID := "00000000-0000-0000-0000-000000005005"
	planetID := "00000000-0000-0000-0001-000000005005"
	defer cleanupShipFactoryE2EPlayer(t, playerID)

	fix := seedShipFactoryE2EFixture(t, playerID, planetID, struct {
		WithShipFactory      bool
		ShipFactoryLevel     int
		WithResources        bool
		Metal, He3, Gold     int64
		WeikesBlueprintLevel int
		ModuleBlueprintLevel int
	}{
		WithShipFactory: true, ShipFactoryLevel: 1,
		WithResources: true, Metal: 100000, He3: 100000, Gold: 100000,
		WeikesBlueprintLevel: 1, ModuleBlueprintLevel: 1,
	})

	designID := makeShipDesignDirectly(t, fix.PlayerID, fix.HullTypeID, "cancelme", 1000, 500, 200)

	if code, _ := callBuildShips(t, fix.PlayerID, designID, 2, 1); code != http.StatusOK {
		t.Fatalf("BuildShips setup failed: code=%d", code)
	}

	// Capture pre-cancel resources.
	var preMetal int64
	database.DB.QueryRow(`SELECT metal FROM resources WHERE planet_id=$1`, fix.PlanetID).Scan(&preMetal)

	// Now cancel slot 1
	req := httptest.NewRequest(http.MethodPost, "/api/ship-factory/cancel/1", nil)
	req.SetPathValue("slot", "1")
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, fix.PlayerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	CancelShipBuild(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("CancelShipBuild expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	// Verify ships row no longer marks slot 1 as in-progress.
	var isBuilding bool
	var buildQty int
	database.DB.QueryRow(`
		SELECT is_building, build_quantity FROM ships
		WHERE player_id=$1 AND ship_design_id=$2
	`, fix.PlayerID, designID).Scan(&isBuilding, &buildQty)
	if isBuilding {
		t.Errorf("expected is_building=false after cancel, got true")
	}
	if buildQty != 0 {
		t.Errorf("expected build_quantity=0 after cancel, got %d", buildQty)
	}

	// Resources are NOT refunded by current handler (GO2 mechanics).
	var postMetal int64
	database.DB.QueryRow(`SELECT metal FROM resources WHERE planet_id=$1`, fix.PlanetID).Scan(&postMetal)
	if postMetal != preMetal {
		t.Errorf("expected no resource refund: pre=%d, post=%d", preMetal, postMetal)
	}
}

// =============================================================================
// 1f. After completion: build with finish_at in the past auto-completes
// when GetShipFactorySlots is polled (handler calls completeFinishedBuilds).
// Verifies that the worker logic moves build_quantity into quantity.
//
// Note: ship_instances rows are NOT created by the ship factory pipeline —
// the ships table just stores aggregate quantity. The migrations show
// ship_instances are created by separate flows (e.g. instance combat / spacedock).
// =============================================================================
func TestShipFactoryE2E_BuildCompletes_OnPoll(t *testing.T) {
	setupShipFactoryE2E(t)
	playerID := "00000000-0000-0000-0000-000000005006"
	planetID := "00000000-0000-0000-0001-000000005006"
	defer cleanupShipFactoryE2EPlayer(t, playerID)

	fix := seedShipFactoryE2EFixture(t, playerID, planetID, struct {
		WithShipFactory      bool
		ShipFactoryLevel     int
		WithResources        bool
		Metal, He3, Gold     int64
		WeikesBlueprintLevel int
		ModuleBlueprintLevel int
	}{
		WithShipFactory: true, ShipFactoryLevel: 1,
		WithResources: true, Metal: 100000, He3: 100000, Gold: 100000,
		WeikesBlueprintLevel: 1, ModuleBlueprintLevel: 1,
	})

	designID := makeShipDesignDirectly(t, fix.PlayerID, fix.HullTypeID, "completes", 1000, 500, 200)

	if code, _ := callBuildShips(t, fix.PlayerID, designID, 3, 1); code != http.StatusOK {
		t.Fatalf("BuildShips setup failed: code=%d", code)
	}

	// Force build_finish_at into the past so completeFinishedBuilds() flushes it.
	_, err := database.DB.Exec(`
		UPDATE ships SET build_finish_at = now() - interval '1 hour'
		WHERE player_id = $1 AND ship_design_id = $2
	`, fix.PlayerID, designID)
	if err != nil {
		t.Fatalf("backdate build_finish_at: %v", err)
	}

	// Poll GetShipFactorySlots — its handler internally calls completeFinishedBuilds.
	req := httptest.NewRequest(http.MethodGet, "/api/ship-factory/slots", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, fix.PlayerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	GetShipFactorySlots(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("GetShipFactorySlots: code=%d body=%s", rr.Code, rr.Body.String())
	}

	// Verify the build flushed: quantity=3, is_building=false.
	var qty int
	var isBuilding bool
	err = database.DB.QueryRow(`
		SELECT quantity, is_building FROM ships
		WHERE player_id=$1 AND ship_design_id=$2
	`, fix.PlayerID, designID).Scan(&qty, &isBuilding)
	if err != nil {
		t.Fatalf("read ships: %v", err)
	}
	if qty != 3 {
		t.Errorf("expected quantity=3 after completion, got %d", qty)
	}
	if isBuilding {
		t.Errorf("expected is_building=false after completion, got true")
	}
}

// =============================================================================
// 1g-i. ListShipDesigns happy path (empty + with rows)
// =============================================================================
func TestShipFactoryE2E_ListShipDesigns_Empty(t *testing.T) {
	setupShipFactoryE2E(t)
	playerID := "00000000-0000-0000-0000-000000005007"
	planetID := "00000000-0000-0000-0001-000000005007"
	defer cleanupShipFactoryE2EPlayer(t, playerID)

	seedShipFactoryE2EFixture(t, playerID, planetID, struct {
		WithShipFactory      bool
		ShipFactoryLevel     int
		WithResources        bool
		Metal, He3, Gold     int64
		WeikesBlueprintLevel int
		ModuleBlueprintLevel int
	}{})

	req := httptest.NewRequest(http.MethodGet, "/api/ship-designs", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	ListShipDesigns(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("ListShipDesigns: code=%d", rr.Code)
	}
	var out []models.ShipDesign
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("parse list: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty list, got %d", len(out))
	}
}

func TestShipFactoryE2E_ListShipDesigns_WithRows(t *testing.T) {
	setupShipFactoryE2E(t)
	playerID := "00000000-0000-0000-0000-000000005008"
	planetID := "00000000-0000-0000-0001-000000005008"
	defer cleanupShipFactoryE2EPlayer(t, playerID)

	fix := seedShipFactoryE2EFixture(t, playerID, planetID, struct {
		WithShipFactory      bool
		ShipFactoryLevel     int
		WithResources        bool
		Metal, He3, Gold     int64
		WeikesBlueprintLevel int
		ModuleBlueprintLevel int
	}{
		WeikesBlueprintLevel: 1,
	})

	makeShipDesignDirectly(t, fix.PlayerID, fix.HullTypeID, "alpha", 100, 50, 25)
	makeShipDesignDirectly(t, fix.PlayerID, fix.HullTypeID, "beta", 200, 100, 50)

	req := httptest.NewRequest(http.MethodGet, "/api/ship-designs", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	ListShipDesigns(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("ListShipDesigns: code=%d", rr.Code)
	}
	var out []models.ShipDesign
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("parse list: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 designs, got %d", len(out))
	}
}

// =============================================================================
// 1g (continued). CreateShipDesign happy path: blueprint owned + activated
// =============================================================================
func TestShipFactoryE2E_CreateShipDesign_HappyPath(t *testing.T) {
	setupShipFactoryE2E(t)
	playerID := "00000000-0000-0000-0000-000000005009"
	planetID := "00000000-0000-0000-0001-000000005009"
	defer cleanupShipFactoryE2EPlayer(t, playerID)

	fix := seedShipFactoryE2EFixture(t, playerID, planetID, struct {
		WithShipFactory      bool
		ShipFactoryLevel     int
		WithResources        bool
		Metal, He3, Gold     int64
		WeikesBlueprintLevel int
		ModuleBlueprintLevel int
	}{
		WeikesBlueprintLevel: 1, ModuleBlueprintLevel: 1,
	})

	if fix.ModuleTypeID == 0 {
		t.Skip("no tier-1 module seeded; cannot run CreateShipDesign happy path")
	}

	body, _ := json.Marshal(createDesignRequest{
		Name:       "happy-design",
		HullTypeID: fix.HullTypeID,
		Modules: []models.DesignModule{
			{ModuleTypeID: fix.ModuleTypeID, Quantity: 1, PlacementOrder: 0},
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/ship-designs", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	CreateShipDesign(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("CreateShipDesign expected 201, got %d: %s", rr.Code, rr.Body.String())
	}
}

// =============================================================================
// 1g (rejection). CreateShipDesign rejected when hull blueprint not owned
// =============================================================================
func TestShipFactoryE2E_CreateShipDesign_HullBlueprintNotOwned(t *testing.T) {
	setupShipFactoryE2E(t)
	playerID := "00000000-0000-0000-0000-000000005010"
	planetID := "00000000-0000-0000-0001-000000005010"
	defer cleanupShipFactoryE2EPlayer(t, playerID)

	// No blueprint for hull
	fix := seedShipFactoryE2EFixture(t, playerID, planetID, struct {
		WithShipFactory      bool
		ShipFactoryLevel     int
		WithResources        bool
		Metal, He3, Gold     int64
		WeikesBlueprintLevel int
		ModuleBlueprintLevel int
	}{
		ModuleBlueprintLevel: 1,
	})

	if fix.ModuleTypeID == 0 {
		t.Skip("no tier-1 module seeded")
	}

	body, _ := json.Marshal(createDesignRequest{
		Name:       "no-bp",
		HullTypeID: fix.HullTypeID,
		Modules: []models.DesignModule{
			{ModuleTypeID: fix.ModuleTypeID, Quantity: 1, PlacementOrder: 0},
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/ship-designs", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	CreateShipDesign(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("expected 409 (hull blueprint not owned), got %d: %s", rr.Code, rr.Body.String())
	}
}

// =============================================================================
// 1h. Tier filter via HTTP: tier-2 module on tier-1 hull rejected.
//
// This is the E2E (HTTP) complement to ship_designs_tier_test.go which already
// covers it at the function level. We exercise CreateShipDesign through the
// full handler with a tier-2 module and expect 403.
// =============================================================================
func TestShipFactoryE2E_TierFiltering_Tier2ModuleOnTier1Hull(t *testing.T) {
	setupShipFactoryE2E(t)
	playerID := "00000000-0000-0000-0000-000000005011"
	planetID := "00000000-0000-0000-0001-000000005011"
	defer cleanupShipFactoryE2EPlayer(t, playerID)

	fix := seedShipFactoryE2EFixture(t, playerID, planetID, struct {
		WithShipFactory      bool
		ShipFactoryLevel     int
		WithResources        bool
		Metal, He3, Gold     int64
		WeikesBlueprintLevel int
		ModuleBlueprintLevel int
	}{
		WeikesBlueprintLevel: 1, // tier-1 hull blueprint
	})

	// Find a tier-2 module and its blueprint chain.
	var tier2ModuleID int
	var tier2ModuleName string
	err := database.DB.QueryRow(`
		SELECT id, name FROM module_types WHERE tier = 2 LIMIT 1
	`).Scan(&tier2ModuleID, &tier2ModuleName)
	if err != nil {
		t.Skip("no tier-2 module seeded; cannot validate tier gating")
	}

	// Find the tier-1 form of the same module line and its blueprint
	tier1Name := tier2ModuleName
	if len(tier2ModuleName) > 3 && tier2ModuleName[len(tier2ModuleName)-3:] == "_ii" {
		tier1Name = tier2ModuleName[:len(tier2ModuleName)-3] + "_i"
	}
	var tier1ModuleBlueprintID int
	if err := database.DB.QueryRow(`
		SELECT b.id FROM blueprints b
		JOIN module_types mt ON b.module_type_id = mt.id
		WHERE b.blueprint_type = 'module' AND mt.name = $1
	`, tier1Name).Scan(&tier1ModuleBlueprintID); err != nil {
		t.Skipf("no tier-1 blueprint companion for module %s", tier2ModuleName)
	}
	// Activate at research_level=1 only — does NOT unlock tier 2.
	database.DB.Exec(`
		INSERT INTO player_blueprints (player_id, blueprint_id, is_activated, research_level)
		VALUES ($1, $2, true, 1)
		ON CONFLICT (player_id, blueprint_id) DO UPDATE
		SET is_activated = true, research_level = 1
	`, playerID, tier1ModuleBlueprintID)

	body, _ := json.Marshal(createDesignRequest{
		Name:       "tier2bad",
		HullTypeID: fix.HullTypeID,
		Modules: []models.DesignModule{
			{ModuleTypeID: tier2ModuleID, Quantity: 1, PlacementOrder: 0},
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/ship-designs", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	CreateShipDesign(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Errorf("expected 403 (module tier not unlocked), got %d: %s", rr.Code, rr.Body.String())
	}
}

// =============================================================================
// Bonus: BuildShips request validation — invalid quantity / slot
// =============================================================================
func TestShipFactoryE2E_BuildRejected_InvalidQuantity(t *testing.T) {
	setupShipFactoryE2E(t)
	playerID := "00000000-0000-0000-0000-000000005012"
	planetID := "00000000-0000-0000-0001-000000005012"
	defer cleanupShipFactoryE2EPlayer(t, playerID)

	fix := seedShipFactoryE2EFixture(t, playerID, planetID, struct {
		WithShipFactory      bool
		ShipFactoryLevel     int
		WithResources        bool
		Metal, He3, Gold     int64
		WeikesBlueprintLevel int
		ModuleBlueprintLevel int
	}{
		WithShipFactory: true, ShipFactoryLevel: 1,
		WithResources: true, Metal: 100000, He3: 100000, Gold: 100000,
		WeikesBlueprintLevel: 1,
	})
	designID := makeShipDesignDirectly(t, fix.PlayerID, fix.HullTypeID, "qty0", 1000, 500, 200)

	// quantity=0 must 400
	if code, _ := callBuildShips(t, fix.PlayerID, designID, 0, 1); code != http.StatusBadRequest {
		t.Errorf("quantity=0 expected 400, got %d", code)
	}
	// quantity above max must 400
	if code, _ := callBuildShips(t, fix.PlayerID, designID, 3000000, 1); code != http.StatusBadRequest {
		t.Errorf("quantity over max expected 400, got %d", code)
	}
	// slot=0 must 400
	if code, _ := callBuildShips(t, fix.PlayerID, designID, 1, 0); code != http.StatusBadRequest {
		t.Errorf("slot=0 expected 400, got %d", code)
	}
}

// =============================================================================
// Bonus: BuildShips rejected when slot is already in use
// =============================================================================
func TestShipFactoryE2E_BuildRejected_SlotInUse(t *testing.T) {
	setupShipFactoryE2E(t)
	playerID := "00000000-0000-0000-0000-000000005013"
	planetID := "00000000-0000-0000-0001-000000005013"
	defer cleanupShipFactoryE2EPlayer(t, playerID)

	fix := seedShipFactoryE2EFixture(t, playerID, planetID, struct {
		WithShipFactory      bool
		ShipFactoryLevel     int
		WithResources        bool
		Metal, He3, Gold     int64
		WeikesBlueprintLevel int
		ModuleBlueprintLevel int
	}{
		WithShipFactory: true, ShipFactoryLevel: 1,
		WithResources: true, Metal: 100000, He3: 100000, Gold: 100000,
		WeikesBlueprintLevel: 1,
	})
	d1 := makeShipDesignDirectly(t, fix.PlayerID, fix.HullTypeID, "first", 100, 50, 25)
	d2 := makeShipDesignDirectly(t, fix.PlayerID, fix.HullTypeID, "second", 100, 50, 25)

	if code, _ := callBuildShips(t, fix.PlayerID, d1, 1, 1); code != http.StatusOK {
		t.Fatalf("first build: code=%d", code)
	}
	if code, _ := callBuildShips(t, fix.PlayerID, d2, 1, 1); code != http.StatusConflict {
		t.Errorf("expected 409 (slot already in use), got %d", code)
	}
}

// =============================================================================
// Bonus: GetShipFactory returns 404 when no factory
// =============================================================================
func TestShipFactoryE2E_GetShipFactory_NoFactory(t *testing.T) {
	setupShipFactoryE2E(t)
	playerID := "00000000-0000-0000-0000-000000005014"
	planetID := "00000000-0000-0000-0001-000000005014"
	defer cleanupShipFactoryE2EPlayer(t, playerID)

	seedShipFactoryE2EFixture(t, playerID, planetID, struct {
		WithShipFactory      bool
		ShipFactoryLevel     int
		WithResources        bool
		Metal, He3, Gold     int64
		WeikesBlueprintLevel int
		ModuleBlueprintLevel int
	}{})

	req := httptest.NewRequest(http.MethodGet, "/api/ship-factory", nil)
	ctx := context.WithValue(req.Context(), middleware.PlayerIDKey, playerID)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	GetShipFactory(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404 when no factory, got %d", rr.Code)
	}
}

