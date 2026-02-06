package models

import "time"

type Planet struct {
	ID               string     `json:"id"`
	PlayerID         string     `json:"player_id"`
	Name             string     `json:"name"`
	PositionX        int        `json:"position_x"`
	PositionY        int        `json:"position_y"`
	IsHomeworld      bool       `json:"is_homeworld"`
	IsRBP            bool       `json:"is_rbp"`
	RBPLevel         int        `json:"rbp_level"`
	ControllingCorpID *string   `json:"controlling_corp_id"`
	ProtectionUntil  *time.Time `json:"protection_until"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
