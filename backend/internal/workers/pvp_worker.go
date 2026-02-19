package workers

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"time"

	"github.com/cryptomines-online/backend/internal/combat"
	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/services"
	"github.com/lib/pq"
)

// StartPvPWorker runs background jobs for PvP travel & combat resolution.
// Job 1: Resolve attacks where fleets have arrived (every 5s)
// Job 2: Return fleets that finished returning (every 5s)
// Job 3: Reset SP daily
func StartPvPWorker() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("[PvP Worker] Started (5s interval)")

	for range ticker.C {
		resolveArrivedAttacks()
		returnArrivedFleets()
		resetDailySP()
	}
}

// resolveArrivedAttacks processes pending_attacks where arrival_at <= now()
func resolveArrivedAttacks() {
	rows, err := database.DB.Query(`
		SELECT id, attacker_id, defender_id, defender_planet_id, fleet_ids, travel_seconds
		FROM pending_attacks
		WHERE status = 'traveling' AND arrival_at <= now()
		ORDER BY arrival_at ASC
		LIMIT 10
	`)
	if err != nil {
		log.Printf("[PvP Worker] Failed to query pending attacks: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			attackID        string
			attackerID      string
			defenderID      string
			defenderPlanet  string
			fleetIDs        []string
			travelSeconds   int
		)
		if err := rows.Scan(&attackID, &attackerID, &defenderID, &defenderPlanet, pq.Array(&fleetIDs), &travelSeconds); err != nil {
			log.Printf("[PvP Worker] Failed to scan pending attack: %v", err)
			continue
		}

		log.Printf("[PvP Worker] Resolving attack %s: %s → %s", attackID, attackerID, defenderID)
		resolveAttack(attackID, attackerID, defenderID, defenderPlanet, fleetIDs, travelSeconds)
	}
}

// resolveAttack runs combat for a single pending attack
func resolveAttack(attackID, attackerID, defenderID, defenderPlanetID string, fleetIDs []string, travelSeconds int) {
	// Load attacker fleets
	var attackerStacks []*combat.FleetStack
	for _, fid := range fleetIDs {
		fleet, err := combat.LoadPlayerFleet(fid, attackerID)
		if err != nil {
			log.Printf("[PvP Worker] Failed to load attacker fleet %s: %v", fid, err)
			continue
		}
		attackerStacks = append(attackerStacks, fleet.Stacks...)
	}

	if len(attackerStacks) == 0 {
		log.Printf("[PvP Worker] No attacker stacks for attack %s, marking resolved", attackID)
		markAttackResolved(attackID, "", fleetIDs, travelSeconds)
		return
	}

	// Create combined attacker fleet
	attackerFleet := &combat.Fleet{
		PlayerID:       attackerID,
		FleetID:        "combined_attacker",
		CommanderBonus: nil,
		TechBonuses:    &services.TechBonuses{},
		Stacks:         attackerStacks,
		Formation:      "phalanx",
		Targeting:      "max_attack",
		Side:           "attacker",
	}

	// Load attacker tech bonuses
	if tb, err := services.GetPlayerTechBonuses(attackerID); err == nil {
		attackerFleet.TechBonuses = tb
	}

	// Load defender fleets (all stationed fleets on the planet)
	defenderFleet, err := loadDefenderFleetsForWorker(defenderPlanetID, defenderID)
	if err != nil {
		log.Printf("[PvP Worker] Failed to load defender fleets: %v", err)
		markAttackResolved(attackID, "", fleetIDs, travelSeconds)
		return
	}

	var result string
	var totalRounds int
	var loot *pvpLootW
	var attackerHe3, defenderHe3 int64

	if len(defenderFleet.Stacks) == 0 {
		// Auto-win
		result = "attacker_win"
		totalRounds = 0
		loot = calculateLootW(defenderPlanetID, attackerStacks, fleetIDs)
		awardLootW(attackerID, loot)
		deductLootW(defenderPlanetID, loot)
		resetDefenseBuildingsW(defenderPlanetID)
	} else {
		// Run combat engine
		engine := combat.NewCombatEngine(time.Now().UnixNano())
		combatResult, err := engine.ExecuteCombat(attackerFleet, defenderFleet)
		if err != nil {
			log.Printf("[PvP Worker] Combat execution failed: %v", err)
			markAttackResolved(attackID, "", fleetIDs, travelSeconds)
			return
		}

		if combatResult.Winner == "attacker" {
			result = "attacker_win"
		} else if combatResult.Winner == "defender" {
			result = "defender_win"
		} else {
			result = "draw"
		}
		totalRounds = combatResult.TotalRounds

		// Apply casualties to attacker fleets
		for _, stack := range attackerFleet.Stacks {
			destroyed := stack.ShipCount - stack.CurrentShips
			if destroyed > 0 {
				database.DB.Exec(`UPDATE fleet_stacks SET ship_count = $1 WHERE id = $2`, stack.CurrentShips, stack.ID)
			}
		}

		// Apply casualties to defender fleets
		for _, stack := range defenderFleet.Stacks {
			destroyed := stack.ShipCount - stack.CurrentShips
			if destroyed > 0 {
				database.DB.Exec(`UPDATE fleet_stacks SET ship_count = $1 WHERE id = $2`, stack.CurrentShips, stack.ID)
			}
		}

		if result == "attacker_win" {
			loot = calculateLootW(defenderPlanetID, attackerStacks, fleetIDs)
			awardLootW(attackerID, loot)
			deductLootW(defenderPlanetID, loot)
			resetDefenseBuildingsW(defenderPlanetID)
		}

		// He3 consumption = sum(design.he3_per_round * avg_ships) * totalRounds
		attackerHe3 = calculateCombatHe3(fleetIDs, totalRounds)
		defenderHe3 = int64(0) // defender He3 tracked separately if needed
	}

	// Create combat report
	reportID := createPvPReportW(attackerID, defenderID, result, totalRounds, loot, attackerHe3, defenderHe3)

	// Update quest progress
	services.UpdateQuestProgress(attackerID, "pvp_attack", "attack", 1)
	if result == "attacker_win" {
		services.UpdateQuestProgress(attackerID, "pvp_win", "win", 1)
	}

	markAttackResolved(attackID, reportID, fleetIDs, travelSeconds)
	log.Printf("[PvP Worker] Attack %s resolved: %s (rounds: %d)", attackID, result, totalRounds)
}

// markAttackResolved updates attack status and sets fleets to returning
func markAttackResolved(attackID, reportID string, fleetIDs []string, travelSeconds int) {
	returnSeconds := travelSeconds / 2
	if returnSeconds < 15 {
		returnSeconds = 15
	}

	var reportIDPtr *string
	if reportID != "" {
		reportIDPtr = &reportID
	}

	database.DB.Exec(`
		UPDATE pending_attacks
		SET status = 'resolved', combat_report_id = $1,
		    return_at = now() + ($2 || ' seconds')::interval
		WHERE id = $3
	`, reportIDPtr, returnSeconds, attackID)

	// Set attacker fleets to returning
	for _, fid := range fleetIDs {
		database.DB.Exec(`
			UPDATE fleets
			SET status = 'returning', arrival_at = now() + ($1 || ' seconds')::interval, updated_at = now()
			WHERE id = $2 AND status = 'traveling'
		`, returnSeconds, fid)
	}
}

// returnArrivedFleets transitions returning fleets back to stationed
func returnArrivedFleets() {
	result, err := database.DB.Exec(`
		UPDATE fleets
		SET status = 'stationed', destination_x = NULL, destination_y = NULL, arrival_at = NULL, updated_at = now()
		WHERE status = 'returning' AND arrival_at <= now()
	`)
	if err != nil {
		log.Printf("[PvP Worker] Failed to return fleets: %v", err)
		return
	}
	if n, _ := result.RowsAffected(); n > 0 {
		log.Printf("[PvP Worker] Returned %d fleets to stationed", n)
	}
}

// resetDailySP resets space points for players whose last reset was before today
func resetDailySP() {
	today := time.Now().Truncate(24 * time.Hour)
	result, err := database.DB.Exec(`
		UPDATE players
		SET space_points = max_space_points, sp_last_reset = now()
		WHERE sp_last_reset < $1
	`, today)
	if err != nil {
		log.Printf("[PvP Worker] Failed to reset SP: %v", err)
		return
	}
	if n, _ := result.RowsAffected(); n > 0 {
		log.Printf("[PvP Worker] Reset SP for %d players", n)
	}
}

// --- Helper types and functions (worker-local to avoid circular deps) ---

type pvpLootW struct {
	Metal int64 `json:"metal"`
	He3   int64 `json:"he3"`
	Gold  int64 `json:"gold"`
}

func calculateLootW(defenderPlanetID string, _ []*combat.FleetStack, fleetIDs []string) *pvpLootW {
	var metal, he3, gold int64
	err := database.DB.QueryRow(`
		SELECT metal, he3, gold FROM resources WHERE planet_id = $1
	`, defenderPlanetID).Scan(&metal, &he3, &gold)
	if err != nil {
		return &pvpLootW{}
	}

	lootMetal := metal * 20 / 100
	lootHe3 := he3 * 20 / 100
	lootGold := gold * 20 / 100

	// Cap loot by fleet cargo capacity (sum of total_storage * ship_count across all stacks)
	cargoCap := getFleetCargoCapacity(fleetIDs)
	if cargoCap > 0 {
		if lootMetal > cargoCap {
			lootMetal = cargoCap
		}
		if lootHe3 > cargoCap {
			lootHe3 = cargoCap
		}
		if lootGold > cargoCap {
			lootGold = cargoCap
		}
	}

	return &pvpLootW{Metal: lootMetal, He3: lootHe3, Gold: lootGold}
}

// getFleetCargoCapacity returns total cargo from all ship designs in the given fleets
func getFleetCargoCapacity(fleetIDs []string) int64 {
	var totalCargo int64
	for _, fid := range fleetIDs {
		rows, err := database.DB.Query(`
			SELECT sd.total_storage, fs.ship_count
			FROM fleet_stacks fs
			JOIN ship_designs sd ON sd.id = fs.ship_design_id
			WHERE fs.fleet_id = $1 AND fs.ship_count > 0
		`, fid)
		if err != nil {
			continue
		}
		for rows.Next() {
			var storage, count int
			if err := rows.Scan(&storage, &count); err != nil {
				continue
			}
			totalCargo += int64(storage) * int64(count)
		}
		rows.Close()
	}
	return totalCargo
}

func awardLootW(attackerID string, loot *pvpLootW) {
	if loot == nil {
		return
	}
	database.DB.Exec(`
		UPDATE resources SET metal = metal + $1, he3 = he3 + $2, gold = gold + $3, updated_at = now()
		WHERE planet_id = (SELECT id FROM planets WHERE player_id = $4 AND is_homeworld = true LIMIT 1)
	`, loot.Metal, loot.He3, loot.Gold, attackerID)
}

func deductLootW(defenderPlanetID string, loot *pvpLootW) {
	if loot == nil {
		return
	}
	database.DB.Exec(`
		UPDATE resources SET metal = GREATEST(0, metal - $1), he3 = GREATEST(0, he3 - $2),
		    gold = GREATEST(0, gold - $3), updated_at = now()
		WHERE planet_id = $4
	`, loot.Metal, loot.He3, loot.Gold, defenderPlanetID)
}

func resetDefenseBuildingsW(planetID string) {
	defenseTypes := []string{"meteor_star", "particle_cannon", "anti_aircraft_gun", "thors_cannon", "celestial_base"}
	for _, typeName := range defenseTypes {
		database.DB.Exec(`
			UPDATE buildings SET level = 0, is_upgrading = false, upgrade_finish_at = NULL, updated_at = now()
			WHERE planet_id = $1 AND building_type_id = (SELECT id FROM building_types WHERE name = $2) AND level > 0
		`, planetID, typeName)
	}
}

func createPvPReportW(attackerID, defenderID, result string, totalRounds int, loot *pvpLootW, attackerHe3, defenderHe3 int64) string {
	lootJSON := []byte("{}")
	if loot != nil {
		lootJSON, _ = json.Marshal(map[string]int64{"metal": loot.Metal, "he3": loot.He3, "gold": loot.Gold})
	}

	var reportID string
	database.DB.QueryRow(`
		INSERT INTO combat_reports (attacker_id, defender_id, combat_type, result, total_rounds, loot_json, he3_consumed)
		VALUES ($1, $2, 'pvp', $3, $4, $5, $6) RETURNING id
	`, attackerID, defenderID, result, totalRounds, string(lootJSON), attackerHe3+defenderHe3).Scan(&reportID)

	return reportID
}

func loadDefenderFleetsForWorker(planetID, defenderID string) (*combat.Fleet, error) {
	// Get tech bonuses first — needed for defense buildings and DefenseMovement
	tb, err := services.GetPlayerTechBonuses(defenderID)
	if err != nil {
		log.Printf("[PvP Worker] Failed to get defender tech bonuses: %v", err)
		tb = &services.TechBonuses{}
	}

	rows, err := database.DB.Query(`
		SELECT id FROM fleets WHERE player_id = $1 AND planet_id = $2 AND status = 'stationed'
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
		fleet, err := combat.LoadPlayerFleet(fleetID, defenderID)
		if err != nil {
			continue
		}
		allStacks = append(allStacks, fleet.Stacks...)
	}

	// Apply DefenseMovement: flat speed bonus to all defender ship stacks
	if tb.DefenseMovement > 0 {
		for _, stack := range allStacks {
			stack.BaseSpeed += tb.DefenseMovement
		}
	}

	// Load defense buildings with tech bonuses applied to their stats
	defenseStacks, err := loadDefenseBuildingsForWorker(planetID, tb)
	if err == nil {
		allStacks = append(allStacks, defenseStacks...)
	}

	fleet := &combat.Fleet{
		PlayerID:       defenderID,
		FleetID:        "combined_defender",
		CommanderBonus: nil,
		TechBonuses:    tb,
		Stacks:         allStacks,
		Formation:      "phalanx",
		Targeting:      "max_attack",
		Side:           "defender",
	}

	return fleet, nil
}

// loadDefenseBuildingsForWorker loads defense buildings with tech bonuses applied.
func loadDefenseBuildingsForWorker(planetID string, tb *services.TechBonuses) ([]*combat.FleetStack, error) {
	rows, err := database.DB.Query(`
		SELECT b.id, bt.name, b.level
		FROM buildings b
		JOIN building_types bt ON bt.id = b.building_type
		WHERE b.planet_id = $1
		  AND bt.category = 'defense'
		  AND b.is_upgrading = false
		  AND b.level > 0
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
		stack := buildingToStackW(buildingID, buildingName, level, stackID, tb)
		if stack != nil {
			stacks = append(stacks, stack)
			stackID++
		}
	}
	return stacks, nil
}

// buildingToStackW converts a defense building to a combat stack with tech bonuses.
// Mirrors the logic in handlers/defense_helpers.go buildingToStack.
func buildingToStackW(buildingID, buildingName string, level, stackID int, tb *services.TechBonuses) *combat.FleetStack {
	var baseAttack, baseDefense, baseShield, baseStructure, baseSpeed, baseAccuracy, baseDodge int
	l2 := level * level

	switch buildingName {
	case "space_station":
		baseAttack = 5000 * l2
		baseDefense = 10000 * l2
		baseShield = 200000 * l2
		baseStructure = 300000 * l2
		baseSpeed = 20
		baseAccuracy = 60
		baseDodge = 10
	case "particle_cannon":
		baseAttack = 10000 * l2
		baseDefense = 5000 * l2
		baseShield = 50000 * l2
		baseStructure = 50000 * l2
		baseSpeed = 40
		baseAccuracy = 90
		baseDodge = 20
	case "anti_aircraft_gun":
		baseAttack = 7500 * l2
		baseDefense = 6000 * l2
		baseShield = 60000 * l2
		baseStructure = 100000 * l2
		baseSpeed = 50
		baseAccuracy = 100
		baseDodge = 30
	case "meteor_star":
		baseAttack = 5000 * l2
		baseDefense = 8000 * l2
		baseShield = 80000 * l2
		baseStructure = 120000 * l2
		baseSpeed = 35
		baseAccuracy = 75
		baseDodge = 25
	case "thors_cannon":
		baseAttack = 20000 * l2
		baseDefense = 10000 * l2
		baseShield = 100000 * l2
		baseStructure = 150000 * l2
		baseSpeed = 10
		baseAccuracy = 80
		baseDodge = 5
	default:
		return nil
	}

	// Apply tech bonuses
	if tb != nil {
		if tb.DefenseValue > 0 {
			mult := 1.0 + tb.DefenseValue/100.0
			baseDefense = int(float64(baseDefense) * mult)
			baseShield = int(float64(baseShield) * mult)
			baseStructure = int(float64(baseStructure) * mult)
			baseDodge = int(float64(baseDodge) * mult)
		}
		if tb.EmplacementAttack > 0 {
			switch buildingName {
			case "particle_cannon", "anti_aircraft_gun", "thors_cannon", "meteor_star":
				baseAttack = int(float64(baseAttack) * (1.0 + tb.EmplacementAttack/100.0))
			}
		}
		if tb.DefenseRange > 0 {
			switch buildingName {
			case "particle_cannon", "anti_aircraft_gun":
				baseAccuracy += tb.DefenseRange
			}
		}
	}

	return &combat.FleetStack{
		ID:               buildingID,
		ShipDesignID:     "",
		ShipType:         combat.ShipTypeCruiser,
		DamageType:       combat.DamageExplosive,
		ArmorType:        combat.ArmorChrome,
		ShipCount:        1,
		CurrentShips:     1,
		BaseAttack:       baseAttack,
		BaseDefense:      baseDefense,
		BaseSpeed:        baseSpeed,
		BaseAccuracy:     baseAccuracy,
		BaseDodge:        baseDodge,
		BaseShield:       baseShield,
		BaseStructure:    baseStructure,
		CurrentShield:    baseShield,
		CurrentStructure: baseStructure,
		GridRow:          stackID / 3,
		GridCol:          stackID % 3,
	}
}

// CalculateTravelTime returns travel seconds based on distance and fleet speed.
// Used by both the handler (to set arrival_at) and externally.
func CalculateTravelTime(originX, originY, destX, destY int, fleetSpeed int) int {
	dx := float64(destX - originX)
	dy := float64(destY - originY)
	distance := math.Sqrt(dx*dx + dy*dy)
	if distance < 1 {
		distance = 1
	}

	speed := fleetSpeed
	if speed < 1 {
		speed = 1
	}

	// BASE_TRAVEL_FACTOR = 300 (5 minutes per distance unit at speed 1)
	travelSeconds := int(distance / float64(speed) * 300)

	// Clamp: min 30s, max 3600s (1 hour)
	if travelSeconds < 30 {
		travelSeconds = 30
	}
	if travelSeconds > 3600 {
		travelSeconds = 3600
	}

	return travelSeconds
}

// GetFleetSpeed returns the minimum TotalMovement across all ship designs in the given fleets.
func GetFleetSpeed(fleetIDs []string) int {
	minSpeed := 999
	for _, fid := range fleetIDs {
		rows, err := database.DB.Query(`
			SELECT sd.total_movement
			FROM fleet_stacks fs
			JOIN ship_designs sd ON sd.id = fs.ship_design_id
			WHERE fs.fleet_id = $1 AND fs.ship_count > 0
		`, fid)
		if err != nil {
			continue
		}
		for rows.Next() {
			var mov int
			if err := rows.Scan(&mov); err != nil {
				continue
			}
			if mov < minSpeed {
				minSpeed = mov
			}
		}
		rows.Close()
	}
	if minSpeed == 999 {
		return 1
	}
	return minSpeed
}

// GetPlanetPosition returns planet coordinates
func GetPlanetPosition(planetID string) (int, int, error) {
	var x, y int
	err := database.DB.QueryRow(`SELECT position_x, position_y FROM planets WHERE id = $1`, planetID).Scan(&x, &y)
	return x, y, err
}

// calculateCombatHe3 calculates He3 consumed during combat based on weapon he3_per_round
func calculateCombatHe3(fleetIDs []string, totalRounds int) int64 {
	var totalHe3PerRound int64
	for _, fid := range fleetIDs {
		rows, err := database.DB.Query(`
			SELECT sd.he3_per_round, fs.ship_count
			FROM fleet_stacks fs
			JOIN ship_designs sd ON sd.id = fs.ship_design_id
			WHERE fs.fleet_id = $1 AND fs.ship_count > 0
		`, fid)
		if err != nil {
			continue
		}
		for rows.Next() {
			var he3pr, count int
			if err := rows.Scan(&he3pr, &count); err != nil {
				continue
			}
			totalHe3PerRound += int64(he3pr) * int64(count)
		}
		rows.Close()
	}
	return totalHe3PerRound * int64(totalRounds)
}

// DeductTravelHe3 deducts He3 fuel from attacker resources when fleet departs.
// Formula: sum(design.he3_per_round * ship_count) * ceil(travel_minutes)
func DeductTravelHe3(playerID string, fleetIDs []string, travelSeconds int) int64 {
	var totalHe3PerRound int64
	for _, fid := range fleetIDs {
		rows, err := database.DB.Query(`
			SELECT sd.he3_per_round, fs.ship_count
			FROM fleet_stacks fs
			JOIN ship_designs sd ON sd.id = fs.ship_design_id
			WHERE fs.fleet_id = $1 AND fs.ship_count > 0
		`, fid)
		if err != nil {
			continue
		}
		for rows.Next() {
			var he3pr, count int
			if err := rows.Scan(&he3pr, &count); err != nil {
				continue
			}
			totalHe3PerRound += int64(he3pr) * int64(count)
		}
		rows.Close()
	}

	travelMinutes := int64(travelSeconds+59) / 60 // ceil
	he3Cost := totalHe3PerRound * travelMinutes
	if he3Cost <= 0 {
		return 0
	}

	database.DB.Exec(`
		UPDATE resources SET he3 = GREATEST(0, he3 - $1), updated_at = now()
		WHERE planet_id = (SELECT id FROM planets WHERE player_id = $2 AND is_homeworld = true LIMIT 1)
	`, he3Cost, playerID)

	return he3Cost
}

// Unused import guard
var _ = sql.ErrNoRows
