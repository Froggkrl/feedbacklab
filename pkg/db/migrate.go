package db

import (
	"log"

	// Import pgx driver for PostgreSQL.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Migrate runs database migrations from the specified directory.
func Migrate(dbURL, migrationsDir string) error {
	db, err := goose.OpenDBWithDriver("pgx", dbURL)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	log.Println("migrations applied successfully")
	return nil
}
