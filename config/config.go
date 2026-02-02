package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds application configuration.
type Config struct {
	DB DatabaseConfig
	Server
}

// DatabaseConfig holds PostgreSQL connection settings.
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	Enabled  bool
}

// Server holds server configuration.
type Server struct {
	Port string
}

// DSN returns the PostgreSQL connection string.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
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

	cfg := &Config{
		DB: DatabaseConfig{
			Host:     host,
			Port:     v.GetInt("DB_PORT"),
			User:     v.GetString("DB_USER"),
			Password: v.GetString("DB_PASSWORD"),
			DBName:   v.GetString("DB_NAME"),
			SSLMode:  sslmode,
			Enabled:  v.GetBool("USE_POSTGRES"),
		},
		Server: Server{
			Port: v.GetString("PORT"),
		},
	}

	return cfg, nil
}