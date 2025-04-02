package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *sql.DB) error {
	migrations, err := filepath.Glob("./internal/database/migrations/*.up.sql")
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		content, err := os.ReadFile(migration)
		if err != nil {
			return err
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute %s: %w", migration, err)
		}
	}

	return nil
}
