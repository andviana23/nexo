package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	// Get migration files to apply
	migrations := []string{
		"042_commission_rules.up.sql",
		"043_commission_periods.up.sql",
		"044_advances.up.sql",
		"045_commission_items.up.sql",
	}

	migrationsDir := "migrations"

	for _, migration := range migrations {
		filePath := filepath.Join(migrationsDir, migration)

		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading %s: %v", migration, err)
			continue
		}

		fmt.Printf("Applying %s...\n", migration)

		_, err = conn.Exec(ctx, string(content))
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				fmt.Printf("  ⚠️  %s - already applied (skipped)\n", migration)
			} else {
				log.Printf("  ❌ Error applying %s: %v\n", migration, err)
			}
		} else {
			fmt.Printf("  ✅ %s applied successfully\n", migration)
		}
	}

	fmt.Println("\n✅ Migration process completed!")
}
