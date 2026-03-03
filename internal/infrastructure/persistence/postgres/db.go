package postgres

import (
	"fmt"

	"github.com/contractiq/contractiq/internal/infrastructure/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewConnection creates a new PostgreSQL database connection pool.
func NewConnection(cfg config.DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}
