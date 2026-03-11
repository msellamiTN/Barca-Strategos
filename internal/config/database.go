package config

import (
	"fmt"
	"os"
)

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	DSN string
}

// LoadDatabaseConfig loads database configuration from environment.
func LoadDatabaseConfig() (DatabaseConfig, error) {
	dsn := os.Getenv("PHOENIX_DATABASE_URL")
	if dsn == "" {
		return DatabaseConfig{}, fmt.Errorf("PHOENIX_DATABASE_URL is required")
	}
	return DatabaseConfig{DSN: dsn}, nil
}

// RedisConfig holds Redis connection settings.
type RedisConfig struct {
	Addr string
}

// LoadRedisConfig loads Redis configuration from environment.
func LoadRedisConfig() (RedisConfig, error) {
	addr := os.Getenv("PHOENIX_REDIS_URL")
	if addr == "" {
		addr = "redis:6379"
	}
	return RedisConfig{Addr: addr}, nil
}
