package models

import "time"

type Resource struct {
	ID                   string    `json:"id"`
	PlanetID             string    `json:"planet_id"`
	Metal                int64     `json:"metal"`
	He3                  int64     `json:"he3"`
	Gold                 int64     `json:"gold"`
	MetalPerHour         int64     `json:"metal_per_hour"`
	He3PerHour           int64     `json:"he3_per_hour"`
	GoldPerHour          int64     `json:"gold_per_hour"`
	StorageCapacity      int64     `json:"storage_capacity"`
	WarehouseMetal       int64     `json:"warehouse_metal"`
	WarehouseHe3         int64     `json:"warehouse_he3"`
	WarehouseGold        int64     `json:"warehouse_gold"`
	WarehouseCapacity    int64     `json:"warehouse_capacity"`
	LastWarehouseUpdate  time.Time `json:"last_warehouse_update"`
	LastCollectedAt      time.Time `json:"last_collected_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
