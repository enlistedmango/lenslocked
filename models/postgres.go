package models

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"time"

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
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s "+
			"connect_timeout=30 "+
			"application_name=lenslocked "+
			"fallback_application_name=lenslocked-app",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
	)
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
	// Check for Railway's DATABASE_URL first
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		// For Railway, prefer individual connection parameters
		if os.Getenv("PGHOST") != "" {
			return PostgresConfig{
				Host:     os.Getenv("PGHOST"),
				Port:     envOr("PGPORT", "5432"),
				User:     os.Getenv("PGUSER"),
				Password: os.Getenv("PGPASSWORD"),
				Database: os.Getenv("PGDATABASE"),
				SSLMode:  "disable", // Internal Railway connections don't need SSL
			}
		}

		// Fallback to DATABASE_URL parsing
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
			Database: u.Path[1:],
			SSLMode:  "disable",
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

func OpenWithRetry(config PostgresConfig, maxAttempts int, delay time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("Database connection attempt %d of %d\n", attempt, maxAttempts)
		db, err = sql.Open("pgx", config.String())
		if err != nil {
			fmt.Printf("Error opening database: %v\n", err)
			time.Sleep(delay)
			continue
		}

		// Try to ping
		err = db.Ping()
		if err == nil {
			fmt.Printf("Successfully connected to database on attempt %d\n", attempt)
			return db, nil
		}

		fmt.Printf("Failed to ping database: %v\n", err)
		db.Close()
		time.Sleep(delay)
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxAttempts, err)
}
