package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Port string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
