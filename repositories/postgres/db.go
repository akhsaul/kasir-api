package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"kasir-api/config"
	"kasir-api/helpers/logger"
)

// DB wraps *sql.DB for PostgreSQL storage (pgx driver).
type DB struct {
	*sql.DB
}

// NewDB creates a new PostgreSQL connection and runs migrations.
func NewDB(cfg *config.DatabaseConfig) (*DB, error) {
	dsn := cfg.DSN()
	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}
	// Disable prepared statement cache so it works with Supabase/Supavisor (transaction mode pooler).
	connConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	connStr := stdlib.RegisterConnConfig(connConfig)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	pg := &DB{DB: db}
	if err := pg.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	logger.Info("Database connected successfully")

	return pg, nil
}

// migrate creates tables if they don't exist.
func (db *DB) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			price INTEGER NOT NULL CHECK (price > 0),
			stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
			category_id INTEGER REFERENCES categories(id)
		)`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS category_id INTEGER REFERENCES categories(id)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}
