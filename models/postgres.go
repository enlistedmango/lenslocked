package models

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"

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
	// Check for Railway's DATABASE_URL first
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		u, err := url.Parse(dbURL)
		if err != nil {
			panic(err)
		}
		password, _ := u.User.Password()

		// Check if we're running inside Railway's network
		host := u.Hostname()
		if strings.Contains(host, "autorack.proxy.rlwy.net") {
			// We're connecting from outside Railway, use public hostname
			return PostgresConfig{
				Host:     host,
				Port:     u.Port(),
				User:     u.User.Username(),
				Password: password,
				Database: u.Path[1:], // remove leading "/"
				SSLMode:  "require",  // Railway requires SSL for external connections
			}
		} else {
			// We're inside Railway's network, use internal hostname
			return PostgresConfig{
				Host:     "postgres.railway.internal",
				Port:     "5432",
				User:     u.User.Username(),
				Password: password,
				Database: u.Path[1:],
				SSLMode:  "disable", // Internal connections don't need SSL
			}
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
