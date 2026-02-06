package models

import "time"

type Resource struct {
	ID              string    `json:"id"`
	PlanetID        string    `json:"planet_id"`
	Metal           int64     `json:"metal"`
	He3             int64     `json:"he3"`
	Gold            int64     `json:"gold"`
	MetalPerHour    int64     `json:"metal_per_hour"`
	He3PerHour      int64     `json:"he3_per_hour"`
	GoldPerHour     int64     `json:"gold_per_hour"`
	StorageCapacity int64     `json:"storage_capacity"`
	LastCollectedAt time.Time `json:"last_collected_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
