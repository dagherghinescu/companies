package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// NewDBClient connectrs to a postgres DB and returns a new sql.DB object.
func NewDBClient(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}

	return db, nil
}
