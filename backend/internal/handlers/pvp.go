package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cryptomines-online/backend/internal/combat"
	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/errs"
	"github.com/cryptomines-online/backend/internal/middleware"
	"github.com/cryptomines-online/backend/internal/services"
)

type attackPlanetRequest struct {
	DefenderPlanetID string   `json:"defender_planet_id"`
	FleetIDs         []string `json:"fleet_ids"`
}

type attackPlanetResponse struct {
	ReportID        string        `json:"report_id"`
	Result          string        `json:"result"`
	TotalRounds     int           `json:"total_rounds"`
	LootGained      *pvpLoot      `json:"loot_gained"`
	AttackerLosses  *battleLosses `json:"attacker_losses"`
	DefenderLosses  *battleLosses `json:"defender_losses"`
}

type pvpLoot struct {
	Metal int64 `json:"metal"`
	He3   int64 `json:"he3"`
	Gold  int64 `json:"gold"`
}

// AttackPlanet handles POST /api/pvp/attack
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
	err := database.DB.QueryRow(`
		SELECT player_id, name FROM planets WHERE id = $1
	`, req.DefenderPlanetID).Scan(&defenderID, &defenderPlanetName)

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

	// Check attack cooldown (5 minutes per target)
	var lastAttack sql.NullTime
	err = database.DB.QueryRow(`
		SELECT MAX(created_at) FROM combat_reports
		WHERE attacker_id = $1 AND defender_id = $2 AND combat_type = 'pvp'
	`, attackerID, defenderID).Scan(&lastAttack)

	if err == nil && lastAttack.Valid {
		cooldown := 5 * time.Minute
		if time.Since(lastAttack.Time) < cooldown {
			remaining := cooldown - time.Since(lastAttack.Time)
			http.Error(w, `{"error":"attack cooldown active","remaining_seconds":`+string(rune(int(remaining.Seconds())))+`}`, http.StatusConflict)
			return
		}
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
		if len(stacks) == 0 {
			http.Error(w, `{"error":"fleet has no ships: `+fid+`"}`, http.StatusConflict)
			return
		}
	}

	// Load attacker fleet
	var attackerStacks []*combat.FleetStack
	for _, fid := range req.FleetIDs {
		attackerFleet, err := combat.LoadPlayerFleet(fid, attackerID)
		if err != nil {
			log.Printf("Failed to load attacker fleet %s: %v", fid, err)
			http.Error(w, `{"error":"failed to load attacker fleet"}`, http.StatusInternalServerError)
			return
		}
		attackerStacks = append(attackerStacks, attackerFleet.Stacks...)
	}

	// Create combined attacker fleet
	attackerFleet := &combat.Fleet{
		PlayerID:       attackerID,
		FleetID:        "combined_attacker",
		CommanderBonus: nil,
		TechBonuses:    &combat.TechBonuses{},
		Stacks:         attackerStacks,
		Formation:      "phalanx",
		Targeting:      "max_attack",
		Side:           "attacker",
	}

	// Load attacker tech bonuses
	attackerTechBonuses, err := services.GetPlayerTechBonuses(attackerID)
	if err == nil {
		attackerFleet.TechBonuses = &combat.TechBonuses{
			BallisticDamage:     attackerTechBonuses.BallisticDamage,
			BallisticCritRate:   attackerTechBonuses.BallisticCritRate,
			BallisticCritDamage: attackerTechBonuses.BallisticCritDamage,
			BallisticHitRate:    attackerTechBonuses.BallisticHitRate,
			DirectionalDamage:   attackerTechBonuses.DirectionalDamage,
			DirectionalCritRate: attackerTechBonuses.DirectionalCritRate,
			DirectionalAccuracy: attackerTechBonuses.DirectionalAccuracy,
			MissileDamage:       attackerTechBonuses.MissileDamage,
			MissileHitRate:      attackerTechBonuses.MissileHitRate,
			BaseShield:          attackerTechBonuses.BaseShield,
			BaseStructure:       attackerTechBonuses.BaseStructure,
			BaseAgility:         attackerTechBonuses.BaseAgility,
			BaseDefense:         attackerTechBonuses.BaseDefense,
		}
	}

	// Load defender fleets (all stationed fleets on the planet)
	defenderFleets, err := loadDefenderFleets(req.DefenderPlanetID, defenderID)
	if err != nil {
		log.Printf("Failed to load defender fleets: %v", err)
		http.Error(w, `{"error":"failed to load defender fleets"}`, http.StatusInternalServerError)
		return
	}

	// If no defender fleets, auto-win for attacker
	if len(defenderFleets.Stacks) == 0 {
		// Auto-win: calculate loot and create report
		loot := calculatePvPLoot(req.DefenderPlanetID, attackerStacks)
		awardLoot(attackerID, loot)

		reportID := createPvPReport(attackerID, defenderID, "attacker_win", 0, loot, 0, 0)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(attackPlanetResponse{
			ReportID:       reportID,
			Result:         "attacker_win",
			TotalRounds:    0,
			LootGained:     loot,
			AttackerLosses: &battleLosses{ShipsDestroyed: map[string]int{}, He3Consumed: 0},
			DefenderLosses: &battleLosses{ShipsDestroyed: map[string]int{}, He3Consumed: 0},
		})
		return
	}

	// Execute PvP combat
	engine := combat.NewCombatEngine(time.Now().UnixNano())
	combatResult, err := engine.ExecuteCombat(attackerFleet, defenderFleets)
	if err != nil {
		log.Printf("PvP combat execution failed: %v", err)
		http.Error(w, `{"error":"combat execution failed"}`, http.StatusInternalServerError)
		return
	}

	// Determine result
	var result string
	if combatResult.Winner == "attacker" {
		result = "attacker_win"
	} else if combatResult.Winner == "defender" {
		result = "defender_win"
	} else {
		result = "draw"
	}

	// Apply casualties to attacker fleets
	attackerShipsDestroyed := make(map[string]int)
	for _, stack := range attackerFleet.Stacks {
		destroyed := stack.ShipCount - stack.CurrentShips
		if destroyed > 0 {
			attackerShipsDestroyed[stack.ID] = destroyed
			database.DB.Exec(`UPDATE fleet_stacks SET ship_count = $1 WHERE id = $2`, stack.CurrentShips, stack.ID)
		}
	}

	// Apply casualties to defender fleets
	defenderShipsDestroyed := make(map[string]int)
	for _, stack := range defenderFleets.Stacks {
		destroyed := stack.ShipCount - stack.CurrentShips
		if destroyed > 0 {
			defenderShipsDestroyed[stack.ID] = destroyed
			database.DB.Exec(`UPDATE fleet_stacks SET ship_count = $1 WHERE id = $2`, stack.CurrentShips, stack.ID)
		}
	}

	// Calculate loot (only on attacker victory)
	var loot *pvpLoot
	if result == "attacker_win" {
		loot = calculatePvPLoot(req.DefenderPlanetID, attackerStacks)
		awardLoot(attackerID, loot)
		// Deduct from defender
		deductLoot(req.DefenderPlanetID, loot)
	}

	// He3 consumption (simplified: 1 He3 per ship lost)
	attackerHe3 := int64(combatResult.AttackerCasualties)
	defenderHe3 := int64(combatResult.DefenderCasualties)

	// Create combat report
	reportID := createPvPReport(attackerID, defenderID, result, combatResult.TotalRounds, loot, attackerHe3, defenderHe3)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attackPlanetResponse{
		ReportID:       reportID,
		Result:         result,
		TotalRounds:    combatResult.TotalRounds,
		LootGained:     loot,
		AttackerLosses: &battleLosses{ShipsDestroyed: attackerShipsDestroyed, He3Consumed: attackerHe3},
		DefenderLosses: &battleLosses{ShipsDestroyed: defenderShipsDestroyed, He3Consumed: defenderHe3},
	})
}

// loadDefenderFleets loads all stationed fleets on a planet + defense buildings
func loadDefenderFleets(planetID, defenderID string) (*combat.Fleet, error) {
	// Get all stationed fleets on this planet
	rows, err := database.DB.Query(`
		SELECT id FROM fleets
		WHERE player_id = $1 AND planet_id = $2 AND status = 'stationed'
	`, defenderID, planetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allStacks []*combat.FleetStack
	for rows.Next() {
		var fleetID string
		if err := rows.Scan(&fleetID); err != nil {
			continue
		}

		defenderFleet, err := combat.LoadPlayerFleet(fleetID, defenderID)
		if err != nil {
			log.Printf("Failed to load defender fleet %s: %v", fleetID, err)
			continue
		}
		allStacks = append(allStacks, defenderFleet.Stacks...)
	}

	// Load defense buildings as defensive stacks
	defenseStacks, err := loadDefenseBuildings(planetID)
	if err != nil {
		log.Printf("Failed to load defense buildings: %v", err)
	} else {
		allStacks = append(allStacks, defenseStacks...)
	}

	// Create combined defender fleet
	fleet := &combat.Fleet{
		PlayerID:       defenderID,
		FleetID:        "combined_defender",
		CommanderBonus: nil,
		TechBonuses:    &combat.TechBonuses{},
		Stacks:         allStacks,
		Formation:      "phalanx",
		Targeting:      "max_attack",
		Side:           "defender",
	}

	// Load defender tech bonuses
	defenderTechBonuses, err := services.GetPlayerTechBonuses(defenderID)
	if err == nil {
		fleet.TechBonuses = &combat.TechBonuses{
			BallisticDamage:     defenderTechBonuses.BallisticDamage,
			BallisticCritRate:   defenderTechBonuses.BallisticCritRate,
			BallisticCritDamage: defenderTechBonuses.BallisticCritDamage,
			BallisticHitRate:    defenderTechBonuses.BallisticHitRate,
			DirectionalDamage:   defenderTechBonuses.DirectionalDamage,
			DirectionalCritRate: defenderTechBonuses.DirectionalCritRate,
			DirectionalAccuracy: defenderTechBonuses.DirectionalAccuracy,
			MissileDamage:       defenderTechBonuses.MissileDamage,
			MissileHitRate:      defenderTechBonuses.MissileHitRate,
			BaseShield:          defenderTechBonuses.BaseShield,
			BaseStructure:       defenderTechBonuses.BaseStructure,
			BaseAgility:         defenderTechBonuses.BaseAgility,
			BaseDefense:         defenderTechBonuses.BaseDefense,
		}
	}

	return fleet, nil
}

// loadDefenseBuildings converts defense buildings into combat stacks
func loadDefenseBuildings(planetID string) ([]*combat.FleetStack, error) {
	// Get defense buildings on this planet
	rows, err := database.DB.Query(`
		SELECT b.id, bt.name, b.level
		FROM buildings b
		JOIN building_types bt ON bt.name = b.building_type
		WHERE b.planet_id = $1
		  AND bt.type = 'defense'
		  AND b.construction_end_time IS NULL
	`, planetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stacks []*combat.FleetStack
	stackID := 0
	for rows.Next() {
		var buildingID, buildingName string
		var level int

		if err := rows.Scan(&buildingID, &buildingName, &level); err != nil {
			continue
		}

		// Convert building to combat stack (simplified stats based on building type and level)
		stack := buildingToStack(buildingID, buildingName, level, stackID)
		if stack != nil {
			stacks = append(stacks, stack)
			stackID++
		}
	}

	return stacks, nil
}

// buildingToStack converts a defense building to a combat stack
func buildingToStack(buildingID, buildingName string, level, stackID int) *combat.FleetStack {
	// Base stats by building type (simplified)
	var baseAttack, baseDefense, baseShield, baseStructure, baseSpeed, baseAccuracy, baseDodge int

	switch buildingName {
	case "space_station":
		// Space Station: High HP, low damage
		baseAttack = 50 * level
		baseDefense = 100 * level
		baseShield = 2000 * level
		baseStructure = 3000 * level
		baseSpeed = 20
		baseAccuracy = 60
		baseDodge = 10

	case "particle_cannon":
		// Particle Cannon: High damage, medium HP
		baseAttack = 200 * level
		baseDefense = 50 * level
		baseShield = 500 * level
		baseStructure = 800 * level
		baseSpeed = 40
		baseAccuracy = 90
		baseDodge = 20

	case "anti_aircraft_gun":
		// Anti-Aircraft Gun: High accuracy, AoE
		baseAttack = 150 * level
		baseDefense = 60 * level
		baseShield = 600 * level
		baseStructure = 1000 * level
		baseSpeed = 50
		baseAccuracy = 100
		baseDodge = 30

	case "meteor_star":
		// Meteor Star: Balanced
		baseAttack = 100 * level
		baseDefense = 80 * level
		baseShield = 800 * level
		baseStructure = 1200 * level
		baseSpeed = 35
		baseAccuracy = 75
		baseDodge = 25

	case "thors_cannon":
		// Thor's Cannon: Massive damage, low speed
		baseAttack = 500 * level
		baseDefense = 100 * level
		baseShield = 1000 * level
		baseStructure = 1500 * level
		baseSpeed = 10
		baseAccuracy = 80
		baseDodge = 5

	default:
		return nil
	}

	return &combat.FleetStack{
		ID:                buildingID,
		ShipDesignID:      "", // Not applicable for buildings
		ShipType:          combat.ShipTypeCruiser, // Treat as cruiser class for type advantage
		DamageType:        combat.DamageExplosive, // Defense buildings use explosive damage
		ArmorType:         combat.ArmorChrome,     // Chrome armor
		ShipCount:         1,                      // 1 building = 1 "ship"
		CurrentShips:      1,
		BaseAttack:        baseAttack,
		BaseDefense:       baseDefense,
		BaseSpeed:         baseSpeed,
		BaseAccuracy:      baseAccuracy,
		BaseDodge:         baseDodge,
		BaseShield:        baseShield,
		BaseStructure:     baseStructure,
		CurrentShield:     baseShield,
		CurrentStructure:  baseStructure,
		GridRow:           stackID / 3, // Simple grid placement
		GridCol:           stackID % 3,
	}
}

// calculatePvPLoot calculates 20% of defender's resources, capped by attacker cargo
func calculatePvPLoot(defenderPlanetID string, attackerStacks []*combat.FleetStack) *pvpLoot {
	// Get defender planet resources
	var metal, he3, gold int64
	err := database.DB.QueryRow(`
		SELECT metal, he3, gold FROM resources WHERE planet_id = $1
	`, defenderPlanetID).Scan(&metal, &he3, &gold)

	if err != nil {
		return &pvpLoot{Metal: 0, He3: 0, Gold: 0}
	}

	// 20% loot rate
	lootMetal := metal * 20 / 100
	lootHe3 := he3 * 20 / 100
	lootGold := gold * 20 / 100

	// TODO: Calculate attacker cargo capacity and cap loot
	// For now, cap at 1 million each
	if lootMetal > 1000000 {
		lootMetal = 1000000
	}
	if lootHe3 > 1000000 {
		lootHe3 = 1000000
	}
	if lootGold > 1000000 {
		lootGold = 1000000
	}

	return &pvpLoot{Metal: lootMetal, He3: lootHe3, Gold: lootGold}
}

// awardLoot awards loot to attacker's homeworld
func awardLoot(attackerID string, loot *pvpLoot) {
	if loot == nil {
		return
	}

	database.DB.Exec(`
		UPDATE resources
		SET metal = metal + $1, he3 = he3 + $2, gold = gold + $3, updated_at = now()
		WHERE planet_id = (
			SELECT p.id FROM planets p WHERE p.player_id = $4 AND p.is_homeworld = true LIMIT 1
		)
	`, loot.Metal, loot.He3, loot.Gold, attackerID)
}

// deductLoot deducts loot from defender's planet
func deductLoot(defenderPlanetID string, loot *pvpLoot) {
	if loot == nil {
		return
	}

	database.DB.Exec(`
		UPDATE resources
		SET metal = GREATEST(0, metal - $1),
		    he3 = GREATEST(0, he3 - $2),
		    gold = GREATEST(0, gold - $3),
		    updated_at = now()
		WHERE planet_id = $4
	`, loot.Metal, loot.He3, loot.Gold, defenderPlanetID)
}

// createPvPReport creates a combat report for PvP battle
func createPvPReport(attackerID, defenderID, result string, totalRounds int, loot *pvpLoot, attackerHe3, defenderHe3 int64) string {
	lootJSON := []byte("{}")
	if loot != nil {
		lootJSON, _ = json.Marshal(map[string]int64{
			"metal": loot.Metal,
			"he3":   loot.He3,
			"gold":  loot.Gold,
		})
	}

	var reportID string
	database.DB.QueryRow(`
		INSERT INTO combat_reports (attacker_id, defender_id, combat_type, result, total_rounds, loot_json, he3_consumed)
		VALUES ($1, $2, 'pvp', $3, $4, $5, $6)
		RETURNING id
	`, attackerID, defenderID, result, totalRounds, string(lootJSON), attackerHe3+defenderHe3).Scan(&reportID)

	return reportID
}

// SearchPlanets handles GET /api/pvp/search?query=planetname
func SearchPlanets(w http.ResponseWriter, r *http.Request) {
	playerID := middleware.GetPlayerID(r)
	query := r.URL.Query().Get("query")

	if query == "" {
		http.Error(w, `{"error":"query parameter required"}`, http.StatusBadRequest)
		return
	}

	// Search for planets (exclude own planets)
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
