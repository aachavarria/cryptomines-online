package services

import (
	"testing"
)

// GO2 wiki data: Metal Collector
// Base costs: Metal=85, He3=106, Gold=85
// Cost multiplier: 1.75
// Base time: 40s, Time multiplier: 1.75
// Base production: 1080/hr, Production multiplier: 1.134
var (
	mcBaseMetal int64   = 85
	mcBaseHe3   int64   = 106
	mcBaseGold  int64   = 85
	mcCostMult  float64 = 1.75
	mcBaseTime  int     = 40
	mcTimeMult  float64 = 1.75
	mcBaseProd  int     = 1080
	mcProdMult  float64 = 1.134
)

func TestUpgradeCost_Level1(t *testing.T) {
	// At level 1, multiplier = 1.75^0 = 1.0, so costs equal base costs
	metal, he3, gold := UpgradeCost(mcBaseMetal, mcBaseHe3, mcBaseGold, mcCostMult, 1)
	if metal != 85 {
		t.Errorf("Metal Collector Lv1 metal cost: got %d, want 85", metal)
	}
	if he3 != 106 {
		t.Errorf("Metal Collector Lv1 he3 cost: got %d, want 106", he3)
	}
	if gold != 85 {
		t.Errorf("Metal Collector Lv1 gold cost: got %d, want 85", gold)
	}
}

func TestUpgradeCost_Level2(t *testing.T) {
	// At level 2, multiplier = 1.75^1 = 1.75
	// Expected: 85*1.75=148.75->149, 106*1.75=185.5->186, 85*1.75=148.75->149
	metal, he3, gold := UpgradeCost(mcBaseMetal, mcBaseHe3, mcBaseGold, mcCostMult, 2)
	if metal != 149 {
		t.Errorf("Metal Collector Lv2 metal cost: got %d, want 149", metal)
	}
	if he3 != 186 {
		t.Errorf("Metal Collector Lv2 he3 cost: got %d, want 186", he3)
	}
	if gold != 149 {
		t.Errorf("Metal Collector Lv2 gold cost: got %d, want 149", gold)
	}
}

func TestUpgradeCost_Level5(t *testing.T) {
	// At level 5, multiplier = 1.75^4 = 9.37890625
	// Expected: 85*9.3789..=797.2..->797, 106*9.3789..=994.2..->994, 85*9.3789..=797.2..->797
	metal, he3, gold := UpgradeCost(mcBaseMetal, mcBaseHe3, mcBaseGold, mcCostMult, 5)
	if metal != 797 {
		t.Errorf("Metal Collector Lv5 metal cost: got %d, want 797", metal)
	}
	if he3 != 994 {
		t.Errorf("Metal Collector Lv5 he3 cost: got %d, want 994", he3)
	}
	if gold != 797 {
		t.Errorf("Metal Collector Lv5 gold cost: got %d, want 797", gold)
	}
}

// GO2 wiki data: Civic Center
// Base costs: Metal=550, He3=480, Gold=600
// Cost multiplier: 3.03, Time multiplier: 2.87
// Base time: 300s
func TestUpgradeCost_CivicCenter(t *testing.T) {
	var ccBaseMetal int64 = 550
	var ccBaseHe3 int64 = 480
	var ccBaseGold int64 = 600
	var ccCostMult float64 = 3.03

	tests := []struct {
		level     int
		wantMetal int64
		wantHe3   int64
		wantGold  int64
	}{
		{1, 550, 480, 600},
		{2, 1667, 1454, 1818}, // 550*3.03=1666.5->1667
		{3, 5049, 4407, 5509}, // 550*3.03^2=550*9.1809=5049.495->5049
	}

	for _, tt := range tests {
		metal, he3, gold := UpgradeCost(ccBaseMetal, ccBaseHe3, ccBaseGold, ccCostMult, tt.level)
		if metal != tt.wantMetal {
			t.Errorf("CC Lv%d metal: got %d, want %d", tt.level, metal, tt.wantMetal)
		}
		if he3 != tt.wantHe3 {
			t.Errorf("CC Lv%d he3: got %d, want %d", tt.level, he3, tt.wantHe3)
		}
		if gold != tt.wantGold {
			t.Errorf("CC Lv%d gold: got %d, want %d", tt.level, gold, tt.wantGold)
		}
	}
}

func TestUpgradeTime_Level1(t *testing.T) {
	// Level 1: 40 * 1.75^0 = 40
	got := UpgradeTime(mcBaseTime, mcTimeMult, 1)
	if got != 40 {
		t.Errorf("Metal Collector Lv1 build time: got %d, want 40", got)
	}
}

func TestUpgradeTime_Level2(t *testing.T) {
	// Level 2: 40 * 1.75^1 = 70
	got := UpgradeTime(mcBaseTime, mcTimeMult, 2)
	if got != 70 {
		t.Errorf("Metal Collector Lv2 build time: got %d, want 70", got)
	}
}

func TestUpgradeTime_Level5(t *testing.T) {
	// Level 5: 40 * 1.75^4 = 40 * 9.37890625 = 375.15625 -> 375
	got := UpgradeTime(mcBaseTime, mcTimeMult, 5)
	if got != 375 {
		t.Errorf("Metal Collector Lv5 build time: got %d, want 375", got)
	}
}

func TestUpgradeTime_CivicCenter(t *testing.T) {
	// CC: base=300, mult=2.87
	// Lv1: 300
	// Lv2: 300*2.87 = 861
	// Lv3: 300*2.87^2 = 300*8.2369 = 2471.07 -> 2471
	tests := []struct {
		level int
		want  int
	}{
		{1, 300},
		{2, 861},
		{3, 2471},
	}
	for _, tt := range tests {
		got := UpgradeTime(300, 2.87, tt.level)
		if got != tt.want {
			t.Errorf("CC Lv%d build time: got %d, want %d", tt.level, got, tt.want)
		}
	}
}

func TestProductionRate_Level1(t *testing.T) {
	// Level 1: 1080 * 1.134^0 = 1080
	got := ProductionRate(mcBaseProd, mcProdMult, 1)
	if got != 1080 {
		t.Errorf("Metal Collector Lv1 production: got %d, want 1080", got)
	}
}

func TestProductionRate_Level2(t *testing.T) {
	// Level 2: 1080 * 1.134^1 = 1224.72 -> 1225
	got := ProductionRate(mcBaseProd, mcProdMult, 2)
	if got != 1225 {
		t.Errorf("Metal Collector Lv2 production: got %d, want 1225", got)
	}
}

func TestProductionRate_Level5(t *testing.T) {
	// Level 5: 1080 * 1.134^4 = 1080 * 1.65367... = 1785.97 -> 1786
	got := ProductionRate(mcBaseProd, mcProdMult, 5)
	if got != 1786 {
		t.Errorf("Metal Collector Lv5 production: got %d, want 1786", got)
	}
}

func TestProductionRate_He3Extractor(t *testing.T) {
	// He3 Extractor: base=1180, mult=1.134
	// Lv1: 1180
	// Lv5: 1180 * 1.134^4 = 1180 * 1.65367... = 1951.33 -> 1951
	tests := []struct {
		level int
		want  int64
	}{
		{1, 1180},
		{2, 1338}, // 1180*1.134=1338.12->1338
		{5, 1951},
	}
	for _, tt := range tests {
		got := ProductionRate(1180, 1.134, tt.level)
		if got != tt.want {
			t.Errorf("He3 Extractor Lv%d production: got %d, want %d", tt.level, got, tt.want)
		}
	}
}

func TestProductionRate_ResidentialArea(t *testing.T) {
	// Residential Area: base=1300, mult=1.134
	// Lv1: 1300
	tests := []struct {
		level int
		want  int64
	}{
		{1, 1300},
		{2, 1474}, // 1300*1.134=1474.2->1474
	}
	for _, tt := range tests {
		got := ProductionRate(1300, 1.134, tt.level)
		if got != tt.want {
			t.Errorf("Residential Lv%d production: got %d, want %d", tt.level, got, tt.want)
		}
	}
}

func TestUpgradeCost_ZeroCosts(t *testing.T) {
	// Buildings with zero base costs should remain zero at any level
	metal, he3, gold := UpgradeCost(0, 0, 0, 3.03, 10)
	if metal != 0 || he3 != 0 || gold != 0 {
		t.Errorf("Zero base costs should produce zero: got %d/%d/%d", metal, he3, gold)
	}
}

func TestUpgradeTime_Level1AlwaysBase(t *testing.T) {
	// At level 1 with any multiplier, time should equal base
	tests := []struct {
		base int
		mult float64
	}{
		{40, 1.75},
		{300, 2.87},
		{200, 3.03},
		{100, 1.50},
	}
	for _, tt := range tests {
		got := UpgradeTime(tt.base, tt.mult, 1)
		if got != tt.base {
			t.Errorf("Lv1 time with base=%d mult=%.2f: got %d, want %d", tt.base, tt.mult, got, tt.base)
		}
	}
}
