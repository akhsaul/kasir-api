package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds application configuration.
type Config struct {
	DB        DatabaseConfig
	Server    ServerConfig
	RateLimit RateLimitConfig
}

// RateLimitConfig holds rate limiting settings.
type RateLimitConfig struct {
	Rate  float64 // requests per second per IP
	Burst int     // maximum burst size
}

// DatabaseConfig holds PostgreSQL connection settings.
type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	SSLMode      string
	Enabled      bool
	MaxOpenConns int
	MaxIdleConns int
}

// ServerConfig holds server configuration.
type ServerConfig struct {
	Port string
}

// DSN returns the PostgreSQL connection string.
// Uses default_query_exec_mode=simple_protocol to disable prepared statement cache,
// required when using Supabase/Supavisor (transaction mode pooler).
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s default_query_exec_mode=simple_protocol",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// Load reads configuration: .env file first, then env vars as fallback.
func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		v.SetConfigFile(".env")
		v.SetConfigType("env")
		_ = v.ReadInConfig()
	}

	host := v.GetString("DB_HOST")
	sslmode := v.GetString("DB_SSLMODE")
	if sslmode == "" {
		if host == "localhost" || host == "127.0.0.1" {
			sslmode = "disable"
		} else {
			sslmode = "require"
		}
	}

	maxOpenConns := v.GetInt("DB_MAX_OPEN_CONNS")
	if maxOpenConns <= 0 {
		maxOpenConns = 25
	}
	maxIdleConns := v.GetInt("DB_MAX_IDLE_CONNS")
	if maxIdleConns <= 0 {
		maxIdleConns = 5
	}

	port := v.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	rateLimit := v.GetFloat64("RATE_LIMIT_RPS")
	if rateLimit <= 0 {
		rateLimit = 10 // 10 requests per second default
	}
	rateBurst := v.GetInt("RATE_LIMIT_BURST")
	if rateBurst <= 0 {
		rateBurst = 20
	}

	cfg := &Config{
		RateLimit: RateLimitConfig{
			Rate:  rateLimit,
			Burst: rateBurst,
		},
		DB: DatabaseConfig{
			Host:         host,
			Port:         v.GetInt("DB_PORT"),
			User:         v.GetString("DB_USER"),
			Password:     v.GetString("DB_PASSWORD"),
			DBName:       v.GetString("DB_NAME"),
			SSLMode:      sslmode,
			Enabled:      v.GetBool("USE_POSTGRES"),
			MaxOpenConns: maxOpenConns,
			MaxIdleConns: maxIdleConns,
		},
		Server: ServerConfig{
			Port: port,
		},
	}

	return cfg, nil
}
