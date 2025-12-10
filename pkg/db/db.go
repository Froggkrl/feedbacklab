// Package db provides database connection and migration utilities.
package db

import (
	"log"

	// Import pgx driver for PostgreSQL.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// Connect establishes a connection to the database using the provided URL.
func Connect(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", databaseURL)
	if err != nil {
		return nil, err
	}
	log.Println("database connected successfully")
	return db, nil
}
