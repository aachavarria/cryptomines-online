package services

import (
	"math"
)

// UpgradeCost calculates the cost to upgrade a building from currentLevel to currentLevel+1.
// Formula: Cost(N) = BaseCost * CostMultiplier^(N-1) where N = target level
// For resource buildings, we use the level reference tables from the DB.
// This function uses the generic formula as fallback.
func UpgradeCost(baseCostMetal, baseCostHe3, baseCostGold int64, costMultiplier float64, targetLevel int) (metal, he3, gold int64) {
	multiplier := math.Pow(costMultiplier, float64(targetLevel-1))
	metal = int64(math.Round(float64(baseCostMetal) * multiplier))
	he3 = int64(math.Round(float64(baseCostHe3) * multiplier))
	gold = int64(math.Round(float64(baseCostGold) * multiplier))
	return
}

// UpgradeTime calculates the build time in seconds for upgrading to targetLevel.
// Formula: Time(N) = BaseTime * TimeMultiplier^(N-1)
func UpgradeTime(baseTimeSeconds int, timeMultiplier float64, targetLevel int) int {
	return int(math.Round(float64(baseTimeSeconds) * math.Pow(timeMultiplier, float64(targetLevel-1))))
}

// ProductionRate calculates the production per hour at a given level.
// Formula: Production(N) = BaseProd * ProductionMultiplier^(N-1)
func ProductionRate(baseProduction int, productionMultiplier float64, level int) int64 {
	return int64(math.Round(float64(baseProduction) * math.Pow(productionMultiplier, float64(level-1))))
}
