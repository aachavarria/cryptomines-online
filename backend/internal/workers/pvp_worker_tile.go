package workers

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cryptomines-online/backend/internal/combat/engine"
	"github.com/cryptomines-online/backend/internal/database"
	"github.com/cryptomines-online/backend/internal/services"
)

// useTileEngine reports whether the new tile-based combat engine should
// resolve PvP attacks. Toggled by COMBAT_ENGINE=tile (default = legacy).
//
// During Sprints 2 and 3 the legacy 8-phase engine remains the default so
// production behavior is unchanged. Sprint 4 flips the default to tile;
// Sprint 5 deletes the legacy code path entirely.
func useTileEngine() bool {
	return strings.EqualFold(os.Getenv("COMBAT_ENGINE"), "tile")
}

// resolveAttackTile is the tile-based replacement for resolveAttack
// (legacy). Loads everything via engine.Setup, runs Battle.Execute, and
// persists casualties + a v2 rounds_json into combat_reports.
//
// This intentionally duplicates the legacy resolveAttack body rather than
// branching mid-function: keeps the diff in pvp_worker.go to one feature
// flag check, which simplifies review and the eventual cleanup PR.
func resolveAttackTile(attackID, attackerID, defenderID, defenderPlanetID string, fleetIDs []string, travelSeconds int) {
	// Galaxy positions for attacker and defender.
	atkX, atkY, err := GetPlanetPosition(getAttackerHomeworld(attackerID))
	if err != nil {
		log.Printf("[PvP Worker tile] attacker home: %v", err)
		atkX, atkY = 0, 0
	}
	defX, defY, err := GetPlanetPosition(defenderPlanetID)
	if err != nil {
		log.Printf("[PvP Worker tile] defender pos: %v", err)
		defX, defY = 0, 0
	}

	atkTB, _ := services.GetPlayerTechBonuses(attackerID)
	defTB, _ := services.GetPlayerTechBonuses(defenderID)

	defendingFleets, _ := loadStationedFleetIDs(defenderPlanetID, defenderID)

	battle, err := engine.Setup(engine.SetupParams{
		AttackerID:             attackerID,
		DefenderID:             defenderID,
		DefenderPlanetID:       defenderPlanetID,
		AttackerFleetIDs:       fleetIDs,
		AttackerHomePlanetID:   getAttackerHomeworld(attackerID),
		AttackerGalaxyX:        atkX,
		AttackerGalaxyY:        atkY,
		DefenderGalaxyX:        defX,
		DefenderGalaxyY:        defY,
		DefenderTechBonuses:    defTB,
		AttackerTechBonuses:    atkTB,
		IncludeDefendingFleets: defendingFleets,
	}, time.Now().UnixNano())
	if err != nil {
		log.Printf("[PvP Worker tile] setup failed, falling back to legacy: %v", err)
		// Safe fallback: legacy resolver still lives in pvp_worker.go.
		resolveAttack(attackID, attackerID, defenderID, defenderPlanetID, fleetIDs, travelSeconds)
		return
	}

	rep := battle.Execute()
	totalRounds := rep.Summary.TotalRounds

	result := "draw"
	switch rep.End.Winner {
	case "attacker":
		result = "attacker_win"
	case "defender":
		result = "defender_win"
	}

	var loot *pvpLootW
	if result == "attacker_win" {
		loot = calculateLootW(defenderPlanetID, nil, fleetIDs)
		awardLootW(attackerID, loot)
		deductLootW(defenderPlanetID, loot)
		resetDefenseBuildingsW(defenderPlanetID)
	}

	roundsBlob, err := json.Marshal(rep)
	if err != nil {
		roundsBlob = []byte("{}")
	}

	attackerHe3 := calculateCombatHe3(fleetIDs, totalRounds)
	reportID := createPvPReportW(attackerID, defenderID, result, totalRounds, loot, attackerHe3, 0, roundsBlob)

	services.UpdateQuestProgress(attackerID, "pvp_attack", "attack", 1)
	if result == "attacker_win" {
		services.UpdateQuestProgress(attackerID, "pvp_win", "win", 1)
	}

	markAttackResolved(attackID, reportID, fleetIDs, travelSeconds)
	log.Printf("[PvP Worker tile] Attack %s resolved: %s (rounds: %d)", attackID, result, totalRounds)
}

// loadStationedFleetIDs returns IDs of all stationed fleets owned by the
// defender on the given planet. Used by Setup to build defending fleets.
func loadStationedFleetIDs(planetID, defenderID string) ([]string, error) {
	rows, err := database.DB.Query(`
		SELECT id FROM fleets
		WHERE player_id = $1 AND planet_id = $2 AND status = 'stationed'
	`, defenderID, planetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			out = append(out, id)
		}
	}
	return out, nil
}

// getAttackerHomeworld returns the attacker's homeworld planet ID for
// galaxy position lookup. Mirrors the legacy helper used by handlers.
func getAttackerHomeworld(playerID string) string {
	var pid string
	_ = database.DB.QueryRow(`
		SELECT id FROM planets WHERE player_id = $1 AND is_homeworld = true LIMIT 1
	`, playerID).Scan(&pid)
	return pid
}
