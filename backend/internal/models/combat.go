package models

import "time"

// CombatReport represents a combat report from instance or PvP battle
type CombatReport struct {
	ID          string    `json:"id"`
	AttackerID  string    `json:"attacker_id"`
	DefenderID  string    `json:"defender_id"`
	CombatType  string    `json:"combat_type"` // instance_normal, pvp, etc.
	Result      string    `json:"result"`      // attacker_win, defender_win, draw
	TotalRounds int       `json:"total_rounds"`
	LootJSON    []byte    `json:"loot_json,omitempty"`
	RoundsJSON  []byte    `json:"rounds_json,omitempty"`
	He3Consumed int64     `json:"he3_consumed"`
	CreatedAt   time.Time `json:"created_at"`
}

// RecyclingJob represents a ship recycling job
type RecyclingJob struct {
	ID              string    `json:"id"`
	PlayerID        string    `json:"player_id"`
	ShipDesignID    string    `json:"ship_design_id"`
	MetalGained     int64     `json:"metal_gained"`
	He3Gained       int64     `json:"he3_gained"`
	GoldGained      int64     `json:"gold_gained"`
	DurationSeconds int       `json:"duration_seconds"`
	StartedAt       time.Time `json:"started_at"`
	CompletedAt     time.Time `json:"completed_at"`
	Collected       bool      `json:"collected"`
}
