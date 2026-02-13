package services

import (
	"database/sql"
	"log"
	"math"

	"github.com/cryptomines-online/backend/internal/database"
)

// CorpBonuses holds aggregated corp bonuses from RBPs
type CorpBonuses struct {
	TotalRBPBonus float64 `json:"total_rbp_bonus"`
	RBPCount      int     `json:"rbp_count"`
	CorpLevel     int     `json:"corp_level"`
}

// CalculateRBPBonus returns the resource bonus percentage for a given RBP level.
// GO2 scaling:
//   - Lv1-10:  5% + 0.5% per level  (5.5% to 10%)
//   - Lv11-20: 10% + 1% per level above 10  (11% to 20%)
//   - Lv21-30: 20% + 2% per level above 20  (22% to 40%)
//   - Lv31-50: 40% + 3% per level above 30  (43% to 100%)
//   - Lv51-100: 100% + 4% per level above 50  (104% to 300%)
//
// Maximum bonus: 300% at Lv100
func CalculateRBPBonus(rbpLevel int) float64 {
	if rbpLevel < 1 {
		return 0
	}

	if rbpLevel <= 10 {
		// 5% + 0.5% per level
		return 5.0 + float64(rbpLevel)*0.5
	} else if rbpLevel <= 20 {
		// 10% + 1% per level above 10
		return 10.0 + float64(rbpLevel-10)*1.0
	} else if rbpLevel <= 30 {
		// 20% + 2% per level above 20
		return 20.0 + float64(rbpLevel-20)*2.0
	} else if rbpLevel <= 50 {
		// 40% + 3% per level above 30
		return 40.0 + float64(rbpLevel-30)*3.0
	} else if rbpLevel <= 100 {
		// 100% + 4% per level above 50
		return 100.0 + float64(rbpLevel-50)*4.0
	}

	// Cap at level 100 bonus (300%)
	return 300.0
}

// GetCorpBonuses aggregates all RBP bonuses for a corp
func GetCorpBonuses(corpID string) (*CorpBonuses, error) {
	bonuses := &CorpBonuses{
		TotalRBPBonus: 0,
		RBPCount:      0,
		CorpLevel:     1,
	}

	// Get corp level and wealth
	var wealth int64
	err := database.DB.QueryRow(`
		SELECT level, wealth
		FROM corps
		WHERE id = $1
	`, corpID).Scan(&bonuses.CorpLevel, &wealth)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Corp not found: %s", corpID)
			return nil, err
		}
		log.Printf("Failed to query corp: %v", err)
		return nil, err
	}

	// Recalculate corp level from wealth (in case it's out of sync)
	calculatedLevel := CalculateCorpLevel(wealth)
	if calculatedLevel != bonuses.CorpLevel {
		bonuses.CorpLevel = calculatedLevel
	}

	// Get all RBPs controlled by this corp and sum their bonuses
	rows, err := database.DB.Query(`
		SELECT rbp_level
		FROM planets
		WHERE controlling_corp_id = $1 AND is_rbp = true
	`, corpID)
	if err != nil {
		log.Printf("Failed to query corp RBPs: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rbpLevel int
		if err := rows.Scan(&rbpLevel); err != nil {
			log.Printf("Failed to scan RBP level: %v", err)
			continue
		}

		bonuses.RBPCount++
		bonuses.TotalRBPBonus += CalculateRBPBonus(rbpLevel)
	}

	return bonuses, nil
}

// GetCorpRBPCount returns how many RBPs a corp controls
func GetCorpRBPCount(corpID string) (int, error) {
	var count int
	err := database.DB.QueryRow(`
		SELECT COUNT(*)
		FROM planets
		WHERE controlling_corp_id = $1 AND is_rbp = true
	`, corpID).Scan(&count)
	if err != nil {
		log.Printf("Failed to get RBP count for corp %s: %v", corpID, err)
		return 0, err
	}

	return count, nil
}

// CanCorpControlMoreRBPs checks if corp can control more RBPs (max = corp level)
func CanCorpControlMoreRBPs(corpID string) (bool, error) {
	var level int
	var rbpCount int

	// Get corp level
	err := database.DB.QueryRow(`
		SELECT level FROM corps WHERE id = $1
	`, corpID).Scan(&level)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Corp not found: %s", corpID)
			return false, err
		}
		log.Printf("Failed to query corp level: %v", err)
		return false, err
	}

	// Get current RBP count
	rbpCount, err = GetCorpRBPCount(corpID)
	if err != nil {
		return false, err
	}

	return rbpCount < level, nil
}

// CalculateCorpLevel returns the corp level based on wealth
// GO2: level = floor(wealth / 100000) + 1, capped at 50
func CalculateCorpLevel(wealth int64) int {
	if wealth < 0 {
		wealth = 0
	}

	level := int(math.Floor(float64(wealth)/100000.0)) + 1

	return min(level, 50)
}
