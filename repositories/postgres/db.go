package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	"kasir-api/config"
	"kasir-api/db"
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
	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	pg := &DB{DB: sqlDB}
	if err := pg.runMigrations(cfg.DBName); err != nil {
		sqlDB.Close() //nolint:errcheck // best-effort close on migration failure
		return nil, fmt.Errorf("migrate: %w", err)
	}

	logger.Info("Database connected successfully")

	return pg, nil
}

// runMigrations applies all pending database migrations using golang-migrate.
func (d *DB) runMigrations(dbName string) error {
	sourceDriver, err := iofs.New(db.Migrations, "migrations")
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}

	dbDriver, err := migratepg.WithInstance(d.DB, &migratepg.Config{
		DatabaseName: dbName,
	})
	if err != nil {
		return fmt.Errorf("create migration db driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, dbName, dbDriver)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations: %w", err)
	}

	version, dirty, _ := m.Version()
	logger.Info("Migrations applied — version: %d, dirty: %v", version, dirty)

	return nil
}
