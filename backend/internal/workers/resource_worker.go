package workers

import (
	"log"
	"time"

	"github.com/cryptomines-online/backend/internal/database"
)

// UpdateWarehouseResources calculates and updates warehouse resources for all planets.
// This function processes resource production since the last warehouse update.
func UpdateWarehouseResources() {
	now := time.Now()

	// Query resources that need updating
	// We update all resources, but the calculation naturally handles the elapsed time
	rows, err := database.DB.Query(`
		SELECT r.id, r.planet_id, r.metal_per_hour, r.he3_per_hour, r.gold_per_hour,
		       r.warehouse_metal, r.warehouse_he3, r.warehouse_gold,
		       r.metal, r.he3, r.gold, r.storage_capacity,
		       r.last_warehouse_update,
		       EXTRACT(EPOCH FROM ($1 - r.last_warehouse_update)) / 3600.0 as elapsed_hours
		FROM resources r
		WHERE r.last_warehouse_update < $1
		ORDER BY r.last_warehouse_update ASC
		LIMIT 1000
	`, now)
	if err != nil {
		log.Printf("Failed to query resources for warehouse update: %v", err)
		return
	}
	defer rows.Close()

	type resourceUpdate struct {
		id                 string
		planetID           string
		metalPerHour       int64
		he3PerHour         int64
		goldPerHour        int64
		warehouseMetal     int64
		warehouseHe3       int64
		warehouseGold      int64
		currentMetal       int64
		currentHe3         int64
		currentGold        int64
		storageCapacity    int64
		lastWarehouseUpdate time.Time
		elapsedHours       float64
	}

	var updates []resourceUpdate
	for rows.Next() {
		var u resourceUpdate
		err := rows.Scan(
			&u.id, &u.planetID, &u.metalPerHour, &u.he3PerHour, &u.goldPerHour,
			&u.warehouseMetal, &u.warehouseHe3, &u.warehouseGold,
			&u.currentMetal, &u.currentHe3, &u.currentGold, &u.storageCapacity,
			&u.lastWarehouseUpdate, &u.elapsedHours,
		)
		if err != nil {
			log.Printf("Failed to scan resource row: %v", err)
			continue
		}
		updates = append(updates, u)
	}

	if len(updates) == 0 {
		return
	}

	log.Printf("Updating warehouse for %d resources", len(updates))

	// Process each resource update
	successCount := 0
	for _, u := range updates {
		// Calculate production since last update
		metalProduced := int64(float64(u.metalPerHour) * u.elapsedHours)
		he3Produced := int64(float64(u.he3PerHour) * u.elapsedHours)
		goldProduced := int64(float64(u.goldPerHour) * u.elapsedHours)

		// Calculate new warehouse amounts
		newWarehouseMetal := u.warehouseMetal + metalProduced
		newWarehouseHe3 := u.warehouseHe3 + he3Produced
		newWarehouseGold := u.warehouseGold + goldProduced

		// Cap at storage capacity (warehouse + current resources cannot exceed capacity)
		maxWarehouseMetal := u.storageCapacity - u.currentMetal
		maxWarehouseHe3 := u.storageCapacity - u.currentHe3
		maxWarehouseGold := u.storageCapacity - u.currentGold

		if newWarehouseMetal > maxWarehouseMetal {
			newWarehouseMetal = maxWarehouseMetal
		}
		if newWarehouseHe3 > maxWarehouseHe3 {
			newWarehouseHe3 = maxWarehouseHe3
		}
		if newWarehouseGold > maxWarehouseGold {
			newWarehouseGold = maxWarehouseGold
		}

		// Ensure non-negative
		if newWarehouseMetal < 0 {
			newWarehouseMetal = 0
		}
		if newWarehouseHe3 < 0 {
			newWarehouseHe3 = 0
		}
		if newWarehouseGold < 0 {
			newWarehouseGold = 0
		}

		// Update database
		_, err := database.DB.Exec(`
			UPDATE resources
			SET warehouse_metal = $1,
			    warehouse_he3 = $2,
			    warehouse_gold = $3,
			    last_warehouse_update = $4
			WHERE id = $5
		`, newWarehouseMetal, newWarehouseHe3, newWarehouseGold, now, u.id)

		if err != nil {
			log.Printf("Failed to update warehouse for resource %s (planet %s): %v", u.id, u.planetID, err)
			continue
		}

		successCount++
	}

	log.Printf("Successfully updated warehouse for %d/%d resources", successCount, len(updates))
}

// StartResourceWorker starts a background worker that periodically updates warehouse resources.
// It runs every 5 minutes to auto-accumulate production into warehouses.
func StartResourceWorker() {
	ticker := time.NewTicker(5 * time.Minute)
	log.Println("Resource warehouse worker started (updating every 5 minutes)")

	// Run once immediately on startup
	go UpdateWarehouseResources()

	// Then run on ticker
	go func() {
		for range ticker.C {
			UpdateWarehouseResources()
		}
	}()
}
