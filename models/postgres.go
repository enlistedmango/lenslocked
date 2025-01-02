package models

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func Open(config PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.String())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}
	return db, nil
}

// Renamed from DefaultPostgresConfig to GetPostgresConfig
func GetPostgresConfig() PostgresConfig {
	// Check for DATABASE_URL first (Heroku)
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		u, err := url.Parse(dbURL)
		if err != nil {
			panic(err)
		}
		password, _ := u.User.Password()
		return PostgresConfig{
			Host:     u.Hostname(),
			Port:     u.Port(),
			User:     u.User.Username(),
			Password: password,
			Database: u.Path[1:], // remove leading "/"
			SSLMode:  "require",  // Heroku requires SSL
		}
	}

	// Fallback to individual env vars for local development
	return PostgresConfig{
		Host:     envOr("PSQL_HOST", "localhost"),
		Port:     envOr("PSQL_PORT", "5432"),
		User:     envOr("PSQL_USER", "postgres"),
		Password: envOr("PSQL_PASSWORD", "postgres"),
		Database: envOr("PSQL_DATABASE", "lenslocked"),
		SSLMode:  envOr("PSQL_SSLMODE", "disable"),
	}
}

// Helper function to get env variable with default
func envOr(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
