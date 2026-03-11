package config

import (
	"fmt"
	"os"
)

// Config holds runtime configuration loaded from environment variables.
type Config struct {
	Port        string
	Env         string
	Database    DatabaseConfig
	Redis       RedisConfig
}

// Load reads environment variables and returns a Config with sane defaults.
func Load() (Config, error) {
	cfg := Config{
		Port:        getEnv("PHOENIX_PORT", "8080"),
		Env:         getEnv("PHOENIX_ENV", "development"),
		Database:    DatabaseConfig{},
		Redis:       RedisConfig{},
	}

	var err error
	cfg.Database, err = LoadDatabaseConfig()
	if err != nil {
		return Config{}, fmt.Errorf("database config: %w", err)
	}

	cfg.Redis, err = LoadRedisConfig()
	if err != nil {
		return Config{}, fmt.Errorf("redis config: %w", err)
	}

	if cfg.Port == "" {
		return Config{}, fmt.Errorf("port cannot be empty")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
