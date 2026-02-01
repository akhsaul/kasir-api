package postgres

import (
	"database/sql"
	"fmt"
	"kasir-api/helpers/logger"
	_ "github.com/lib/pq"
	"kasir-api/config"
)

// DB wraps *sql.DB for PostgreSQL storage.
type DB struct {
	*sql.DB
}

// NewDB creates a new PostgreSQL connection and runs migrations.
func NewDB(cfg *config.DatabaseConfig) (*DB, error) {
	dsn := cfg.DSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	pg := &DB{DB: db}
	if err := pg.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

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
