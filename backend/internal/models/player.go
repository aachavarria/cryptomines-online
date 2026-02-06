package models

import "time"

type Player struct {
	ID              string    `json:"id"`
	AnonymousID     string    `json:"anonymous_id"`
	Username        *string   `json:"username"`
	Level           int       `json:"level"`
	Experience      int64     `json:"experience"`
	MallPoints      int       `json:"mall_points"`
	Vouchers        int       `json:"vouchers"`
	HonorPoints     int       `json:"honor_points"`
	ChampionPoints  int       `json:"champion_points"`
	Badges          int       `json:"badges"`
	CorsairsGold    int       `json:"corsairs_gold"`
	TutorialStep    int       `json:"tutorial_step"`
	IsOnline        bool      `json:"is_online"`
	CreatedAt       time.Time `json:"created_at"`
	LastLogin       time.Time `json:"last_login"`
	UpdatedAt       time.Time `json:"updated_at"`
}
