package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/models"
	"github.com/cryptomines-online/backend/internal/services"
)

const maxBuildQuantity = 2000000

type shipFactoryStatus struct {
	Level          int  `json:"level"`
	IsUpgrading    bool `json:"is_upgrading"`
	SpeedBonusPct  int  `json:"speed_bonus_pct"`
	ProductionSlots int `json:"production_slots"`
	ActiveBuilds   int  `json:"active_builds"`
}

type productionSlotStatus struct {
	Slot          int        `json:"slot"`
	InUse         bool       `json:"in_use"`
	ShipDesignID  *string    `json:"ship_design_id,omitempty"`
	DesignName    *string    `json:"design_name,omitempty"`
	Quantity      *int       `json:"quantity,omitempty"`
	BuildFinishAt *time.Time `json:"build_finish_at,omitempty"`
}

type buildRequest struct {
	ShipDesignID   string `json:"ship_design_id"`
	Quantity       int    `json:"quantity"`
	ProductionSlot int    `json:"production_slot"`
}

// GetShipFactory handles GET /api/ship-factory
func GetShipFactory(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	// Complete any finished builds first
	completeFinishedBuilds(playerID)

	// Get Ship Factory building status
	var level int
	var isUpgrading bool
	err := database.DB.QueryRow(
		`SELECT COALESCE(b.level, 0), COALESCE(b.is_upgrading, false)
		 FROM buildings b
		 JOIN building_types bt ON b.building_type = bt.id
		 JOIN planets p ON b.planet_id = p.id
		 WHERE p.player_id = $1 AND bt.name = 'ship_factory'
		 ORDER BY b.level DESC LIMIT 1`, playerID,
	).Scan(&level, &isUpgrading)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"ship factory not built"}`, http.StatusNotFound)
			return
		}
		log.Printf("Failed to get ship factory: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Get factory level data
	var sfl models.ShipFactoryLevel
	err = database.DB.QueryRow(
		`SELECT level, civic_center_req, metal_cost, he3_cost, gold_cost,
		        build_time_seconds, speed_bonus_pct, production_slots
		 FROM ship_factory_levels WHERE level = $1`, level,
	).Scan(&sfl.Level, &sfl.CivicCenterReq, &sfl.MetalCost, &sfl.He3Cost, &sfl.GoldCost,
		&sfl.BuildTime, &sfl.SpeedBonusPct, &sfl.ProductionSlots)
	if err != nil {
		log.Printf("Failed to get factory level data: %v", err)
		// Use defaults
		sfl.SpeedBonusPct = 0
		sfl.ProductionSlots = 1
	}

	// Count active builds
	var activeBuilds int
	database.DB.QueryRow(
		`SELECT COUNT(*) FROM ships
		 WHERE player_id = $1 AND is_building = true`, playerID,
	).Scan(&activeBuilds)

	status := shipFactoryStatus{
		Level:           level,
		IsUpgrading:     isUpgrading,
		SpeedBonusPct:   sfl.SpeedBonusPct,
		ProductionSlots: sfl.ProductionSlots,
		ActiveBuilds:    activeBuilds,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// GetShipFactorySlots handles GET /api/ship-factory/slots
func GetShipFactorySlots(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	completeFinishedBuilds(playerID)

	// Get factory level to determine available slots
	var level int
	err := database.DB.QueryRow(
		`SELECT COALESCE(b.level, 0)
		 FROM buildings b
		 JOIN building_types bt ON b.building_type = bt.id
		 JOIN planets p ON b.planet_id = p.id
		 WHERE p.player_id = $1 AND bt.name = 'ship_factory'
		 ORDER BY b.level DESC LIMIT 1`, playerID,
	).Scan(&level)
	if err != nil {
		http.Error(w, `{"error":"ship factory not built"}`, http.StatusNotFound)
		return
	}

	var productionSlots int
	err = database.DB.QueryRow(
		`SELECT production_slots FROM ship_factory_levels WHERE level = $1`, level,
	).Scan(&productionSlots)
	if err != nil {
		productionSlots = 1
	}

	// Get active builds
	rows, err := database.DB.Query(
		`SELECT s.production_slot, s.ship_design_id, sd.name, s.build_quantity, s.build_finish_at
		 FROM ships s
		 JOIN ship_designs sd ON s.ship_design_id = sd.id
		 WHERE s.player_id = $1 AND s.is_building = true`, playerID,
	)
	if err != nil {
		log.Printf("Failed to get active builds: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	activeSlots := map[int]productionSlotStatus{}
	for rows.Next() {
		var slot int
		var designID string
		var designName string
		var qty int
		var finishAt time.Time
		if err := rows.Scan(&slot, &designID, &designName, &qty, &finishAt); err != nil {
			continue
		}
		activeSlots[slot] = productionSlotStatus{
			Slot:          slot,
			InUse:         true,
			ShipDesignID:  &designID,
			DesignName:    &designName,
			Quantity:      &qty,
			BuildFinishAt: &finishAt,
		}
	}

	// Build response for all slots
	slots := []productionSlotStatus{}
	for i := 1; i <= productionSlots; i++ {
		if s, ok := activeSlots[i]; ok {
			slots = append(slots, s)
		} else {
			slots = append(slots, productionSlotStatus{Slot: i, InUse: false})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slots)
}

// BuildShips handles POST /api/ship-factory/build
func BuildShips(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)

	var req buildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 || req.Quantity > maxBuildQuantity {
		http.Error(w, `{"error":"quantity must be between 1 and 2000000"}`, http.StatusBadRequest)
		return
	}
	if req.ProductionSlot < 1 || req.ProductionSlot > 5 {
		http.Error(w, `{"error":"production slot must be between 1 and 5"}`, http.StatusBadRequest)
		return
	}

	completeFinishedBuilds(playerID)

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Get Ship Factory level and status
	var factoryLevel int
	var factoryUpgrading bool
	err = tx.QueryRow(
		`SELECT COALESCE(b.level, 0), COALESCE(b.is_upgrading, false)
		 FROM buildings b
		 JOIN building_types bt ON b.building_type = bt.id
		 JOIN planets p ON b.planet_id = p.id
		 WHERE p.player_id = $1 AND bt.name = 'ship_factory'
		 ORDER BY b.level DESC LIMIT 1`, playerID,
	).Scan(&factoryLevel, &factoryUpgrading)
	if err != nil {
		http.Error(w, `{"error":"ship factory not built"}`, http.StatusNotFound)
		return
	}
	if factoryUpgrading {
		http.Error(w, `{"error":"ship factory is upgrading"}`, http.StatusConflict)
		return
	}

	// Check production slot is unlocked
	var sfl models.ShipFactoryLevel
	err = tx.QueryRow(
		`SELECT production_slots, speed_bonus_pct FROM ship_factory_levels WHERE level = $1`,
		factoryLevel,
	).Scan(&sfl.ProductionSlots, &sfl.SpeedBonusPct)
	if err != nil {
		log.Printf("Failed to get factory level data: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if req.ProductionSlot > sfl.ProductionSlots {
		http.Error(w, `{"error":"production slot not unlocked"}`, http.StatusConflict)
		return
	}

	// Check slot is not in use
	var slotInUse bool
	err = tx.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM ships WHERE player_id = $1 AND is_building = true AND production_slot = $2)`,
		playerID, req.ProductionSlot,
	).Scan(&slotInUse)
	if err != nil {
		log.Printf("Failed to check slot: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if slotInUse {
		http.Error(w, `{"error":"production slot already in use"}`, http.StatusConflict)
		return
	}

	// Get ship design (must belong to player)
	var design models.ShipDesign
	err = tx.QueryRow(
		`SELECT id, player_id, metal_cost, he3_cost, gold_cost, build_time_seconds
		 FROM ship_designs WHERE id = $1 AND player_id = $2`,
		req.ShipDesignID, playerID,
	).Scan(&design.ID, &design.PlayerID, &design.MetalCost, &design.He3Cost, &design.GoldCost, &design.BuildTimeSeconds)
	if err != nil {
		http.Error(w, `{"error":"ship design not found"}`, http.StatusNotFound)
		return
	}

	// Calculate batch cost (GDD 8.10.1)
	batchMetal, batchHe3, batchGold := services.ShipBuildCost(
		design.MetalCost, design.He3Cost, design.GoldCost, req.Quantity,
	)

	// Deduct resources from homeworld
	var remaining struct{ Metal, He3, Gold int64 }
	err = tx.QueryRow(
		`UPDATE resources
		 SET metal = metal - $1, he3 = he3 - $2, gold = gold - $3, updated_at = now()
		 WHERE planet_id = (
		     SELECT p.id FROM planets p WHERE p.player_id = $4 AND p.is_homeworld = true LIMIT 1
		 ) AND metal >= $1 AND he3 >= $2 AND gold >= $3
		 RETURNING metal, he3, gold`,
		batchMetal, batchHe3, batchGold, playerID,
	).Scan(&remaining.Metal, &remaining.He3, &remaining.Gold)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"insufficient resources"}`, http.StatusConflict)
			return
		}
		log.Printf("Failed to deduct resources: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Calculate batch build time (GDD 8.10.2)
	batchTimeSec := services.ShipBuildTime(design.BuildTimeSeconds, sfl.SpeedBonusPct, req.Quantity)
	finishAt := time.Now().Add(time.Duration(batchTimeSec) * time.Second)

	// Upsert ships row: if player already has ships of this design, update; otherwise insert
	var ship models.Ship
	err = tx.QueryRow(
		`INSERT INTO ships (player_id, ship_design_id, quantity, is_building, build_quantity, build_finish_at, production_slot)
		 VALUES ($1, $2, 0, true, $3, $4, $5)
		 ON CONFLICT (player_id, ship_design_id) DO UPDATE
		 SET is_building = true, build_quantity = $3, build_finish_at = $4, production_slot = $5, updated_at = now()
		 RETURNING id, player_id, ship_design_id, quantity, is_building, build_quantity, build_finish_at, production_slot, updated_at`,
		playerID, req.ShipDesignID, req.Quantity, finishAt, req.ProductionSlot,
	).Scan(&ship.ID, &ship.PlayerID, &ship.ShipDesignID, &ship.Quantity, &ship.IsBuilding,
		&ship.BuildQuantity, &ship.BuildFinishAt, &ship.ProductionSlot, &ship.UpdatedAt)
	if err != nil {
		log.Printf("Failed to start ship build: %v", err)
		http.Error(w, `{"error":"failed to start build"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	type buildResponse struct {
		Slot           int       `json:"slot"`
		ShipDesignID   string    `json:"ship_design_id"`
		Quantity       int       `json:"quantity"`
		BuildFinishAt  time.Time `json:"build_finish_at"`
		ResourcesSpent struct {
			Metal int64 `json:"metal"`
			He3   int64 `json:"he3"`
			Gold  int64 `json:"gold"`
		} `json:"resources_spent"`
	}

	resp := buildResponse{
		Slot:          req.ProductionSlot,
		ShipDesignID:  req.ShipDesignID,
		Quantity:      req.Quantity,
		BuildFinishAt: finishAt,
	}
	resp.ResourcesSpent.Metal = batchMetal
	resp.ResourcesSpent.He3 = batchHe3
	resp.ResourcesSpent.Gold = batchGold

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// CancelShipBuild handles POST /api/ship-factory/cancel/{slot}
func CancelShipBuild(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	slotStr := r.PathValue("slot")
	slot, err := strconv.Atoi(slotStr)
	if err != nil || slot < 1 || slot > 5 {
		http.Error(w, `{"error":"invalid slot number"}`, http.StatusBadRequest)
		return
	}

	// Find the building ship in this slot
	var shipID string
	err = database.DB.QueryRow(
		`SELECT id FROM ships
		 WHERE player_id = $1 AND is_building = true AND production_slot = $2`,
		playerID, slot,
	).Scan(&shipID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"no build in this slot"}`, http.StatusNotFound)
			return
		}
		log.Printf("Failed to find build: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Cancel (no resource refund per GO2 mechanics)
	_, err = database.DB.Exec(
		`UPDATE ships SET is_building = false, build_quantity = 0, build_finish_at = NULL, updated_at = now()
		 WHERE id = $1`, shipID,
	)
	if err != nil {
		log.Printf("Failed to cancel build: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Cancelled bool `json:"cancelled"`
		Slot      int  `json:"slot"`
	}{Cancelled: true, Slot: slot})
}

// completeFinishedBuilds auto-completes any ship builds that have finished.
func completeFinishedBuilds(playerID string) {
	now := time.Now()
	_, err := database.DB.Exec(
		`UPDATE ships
		 SET quantity = quantity + build_quantity,
		     is_building = false,
		     build_quantity = 0,
		     build_finish_at = NULL,
		     updated_at = now()
		 WHERE player_id = $1 AND is_building = true AND build_finish_at <= $2`,
		playerID, now,
	)
	if err != nil {
		log.Printf("Failed to complete finished builds: %v", err)
	}
}
