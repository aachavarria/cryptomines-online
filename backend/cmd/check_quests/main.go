package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("SUPABASE_DB_URL")
	if dbURL == "" {
		dbURL = "postgresql://postgres:postgres@localhost:54322/postgres?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("=== BLUEPRINT SYSTEM STRUCTURE ===\n")

	// Show blueprints table structure
	fmt.Println("1. BLUEPRINTS TABLE (master list of all blueprints)")
	rows, _ := db.Query(`
		SELECT id, name, blueprint_type, hull_type_id, module_type_id
		FROM blueprints
		ORDER BY blueprint_type, id
		LIMIT 10
	`)
	fmt.Println("   ID | Name                    | Type   | Hull ID | Module ID")
	fmt.Println("   ---|-------------------------|--------|---------|----------")
	for rows.Next() {
		var id, hullID, moduleID sql.NullInt64
		var name, bpType string
		rows.Scan(&id, &name, &bpType, &hullID, &moduleID)
		hullStr := "NULL"
		if hullID.Valid {
			hullStr = fmt.Sprintf("%d", hullID.Int64)
		}
		modStr := "NULL"
		if moduleID.Valid {
			modStr = fmt.Sprintf("%d", moduleID.Int64)
		}
		fmt.Printf("   %3d | %-23s | %-6s | %-7s | %s\n", id.Int64, name, bpType, hullStr, modStr)
	}
	rows.Close()

	// Show player_blueprints structure
	fmt.Println("\n2. PLAYER_BLUEPRINTS TABLE (what each player has unlocked)")
	rows, _ = db.Query(`
		SELECT column_name, data_type
		FROM information_schema.columns
		WHERE table_name = 'player_blueprints'
		ORDER BY ordinal_position
	`)
	fmt.Println("   Column Name      | Data Type")
	fmt.Println("   -----------------|----------")
	for rows.Next() {
		var col, dtype string
		rows.Scan(&col, &dtype)
		fmt.Printf("   %-16s | %s\n", col, dtype)
	}
	rows.Close()

	// Show the relationship
	fmt.Println("\n3. HOW IT WORKS:")
	fmt.Println("   blueprints.id (PK)")
	fmt.Println("        ↓")
	fmt.Println("   player_blueprints.blueprint_id (FK)")
	fmt.Println("")
	fmt.Println("   When blueprint_type = 'hull':")
	fmt.Println("     → blueprints.hull_type_id points to hull_types.id")
	fmt.Println("")
	fmt.Println("   When blueprint_type = 'module':")
	fmt.Println("     → blueprints.module_type_id points to module_types.id")
}
