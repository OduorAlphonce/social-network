package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	DatabasePath     string
	AppEnv           string
	AllowedOrigin 	 string
	MigrationsDir	 string
}

var App Config

// Load reads environment variables (optionally from .env) and populates the
// package-level App configuration value.
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	App = Config{
		Port:             getEnv("PORT", "8080"),
		DatabasePath:      mustGetEnv("DATABASE_PATH"),
		AppEnv:           getEnv("APP_ENV", "development"),
		AllowedOrigin:           getEnv("ALLOWED_ORIGIN", "*"),
		MigrationsDir:           mustGetEnv("MIGRATIONS_DIR"),
	}
}

// getEnv returns the environment variable value or the provided fallback.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// mustGetEnv returns the value for a required env var or logs/fails when
// it's missing.
func mustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		log.Fatalf("Required enviroment variable %s is not set", key)
	}
	return value
}
