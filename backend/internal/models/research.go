package models

import (
	"encoding/json"
	"time"
)

type TechType struct {
	ID                int             `json:"id"`
	Name              string          `json:"name"`
	DisplayName       string          `json:"display_name"`
	Tree              string          `json:"tree"`
	PrerequisitesJSON json.RawMessage `json:"prerequisites_json"`
	BaseCostMetal     int64           `json:"base_cost_metal"`
	BaseCostHe3       int64           `json:"base_cost_he3"`
	BaseCostGold      int64           `json:"base_cost_gold"`
	CostMultiplier    float64         `json:"cost_multiplier"`
	BaseTimeSeconds   int             `json:"base_time_seconds"`
	TimeMultiplier    float64         `json:"time_multiplier"`
	MaxLevel          int             `json:"max_level"`
	EffectsJSON       json.RawMessage `json:"effects_json"`
	Description       string          `json:"description"`
}

type Technology struct {
	ID               string     `json:"id"`
	PlayerID         string     `json:"player_id"`
	TechType         int        `json:"tech_type"`
	Level            int        `json:"level"`
	IsResearching    bool       `json:"is_researching"`
	ResearchFinishAt *time.Time `json:"research_finish_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type TechPrerequisite struct {
	Tech  string `json:"tech"`
	Level int    `json:"level"`
}
