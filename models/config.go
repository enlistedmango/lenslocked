package models

import (
	"os"
	"strconv"
)

type Config struct {
	// Database
	Postgres PostgresConfig

	// FiveManage
	FiveManage struct {
		APIKey string
		Debug  bool
	}

	// Security
	CSRF struct {
		Key    string
		Secure bool
	}

	// Session
	Session struct {
		Key string
	}
}

func LoadConfig() Config {
	return Config{
		Postgres: GetPostgresConfig(),
		FiveManage: struct {
			APIKey string
			Debug  bool
		}{
			APIKey: os.Getenv("FIVEMANAGE_API_KEY"),
			Debug:  getEnvBool("FIVEMANAGE_DEBUG", true),
		},
		CSRF: struct {
			Key    string
			Secure bool
		}{
			Key:    os.Getenv("CSRF_KEY"),
			Secure: getEnvBool("CSRF_SECURE", false),
		},
		Session: struct {
			Key string
		}{
			Key: os.Getenv("SESSION_KEY"),
		},
	}
}

func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			return boolValue
		}
	}
	return defaultValue
}
