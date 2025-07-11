package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
)

// Open returns a *sql.DB ready for use.
// dsn example: "postgres://user:pass@localhost:5432/marketplace?sslmode=disable"
func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	// Reasonable production-ish settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(1 * time.Hour)

	// Ping with a 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return db, nil
}
