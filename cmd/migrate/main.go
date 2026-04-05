package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"subscription-service/internal/repository"
)

func main() {
	pool, err := repository.InitDB()
	if err != nil {
		log.Fatalf("❌ DB connection error: %v", err)
	}
	defer pool.Close()

	migrationPath := filepath.Join("migrations", "000001_init_subscriptions.up.sql")
	sqlBytes, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Fatalf("❌ Read migration file: %v", err)
	}

	_, err = pool.Exec(context.Background(), string(sqlBytes))
	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	log.Println("✅ Migrations applied successfully")
}
