package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// Registers the "pgx" driver for database/sql.
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Store wraps the database connection pool.
type Store struct {
	db *sql.DB
}

// New opens a connection pool to PostgreSQL and verifies it with a ping.
func New(dsn string) (*Store, error) {
	// Open the pool (no connection is made yet).
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Configure the connection pool.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Make sure the database is reachable.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Store{db: db}, nil
}

// Ping checks that the database is still reachable.
func (s *Store) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// Close closes the connection pool.
func (s *Store) Close() error {
	return s.db.Close()
}
