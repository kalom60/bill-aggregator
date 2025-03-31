package database

import (
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string) {
	absPath, err := filepath.Abs("internal/database/migrations")
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	migrationPath := "file://" + absPath

	m, err := migrate.New(migrationPath, dbURL)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Database migrated successfully!")
}
