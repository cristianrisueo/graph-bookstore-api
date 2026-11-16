package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Store encapsula el pool de conexiones a la base de datos.
type Store struct {
	db *sql.DB
}

// Códigos de error de Postgres para violaciones de clave foránea. Desde
// Postgres 18, ON DELETE RESTRICT lanza 23001 (restrict_violation); antes de
// esa versión, y también con ON DELETE NO ACTION o en inserciones, lanza
// 23503 (foreign_key_violation).
const (
	pgRestrictViolation   = "23001"
	pgForeignKeyViolation = "23503"
)

// New abre un pool de conexiones a PostgreSQL y lo verifica con un ping.
func New(dsn string) (*Store, error) {
	// Abre el pool (todavía no se establece ninguna conexión).
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Configura el pool de conexiones.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Comprueba que la base de datos es accesible.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Store{db: db}, nil
}

// Ping comprueba que la base de datos sigue siendo accesible.
func (s *Store) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// Close cierra el pool de conexiones.
func (s *Store) Close() error {
	return s.db.Close()
}
