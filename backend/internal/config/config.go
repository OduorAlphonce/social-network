package config

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Port is the HTTP server's TCP port. It uses PORT or --port and defaults
	// to 8080.
	Port string
	// DatabasePath is the required SQLite database file path. It uses
	// DATABASE_PATH or --database-path.
	DatabasePath string
	// BaseAddress is the HTTP server's host or IP address. It uses BASE_ADDRESS
	// or --base-address and defaults to localhost.
	BaseAddress string
	// AppEnv identifies the runtime environment. It uses APP_ENV and defaults
	// to development.
	AppEnv string
	// AllowedOrigin is the origin permitted to make cross-origin requests. It
	// uses ALLOWED_ORIGIN and defaults to "*".
	AllowedOrigin string
	// MigrationsDir is the required directory containing database migrations.
	// It uses MIGRATIONS_DIR.
	MigrationsDir string
}

var App Config

// Load reads configuration from .env, environment variables, and command-line
// options, in increasing order of precedence.
func Load() error {
	return load(os.Args[1:])
}

func load(args []string) error {
	if err := godotenv.Load(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("load .env: %w", err)
		}
	}

	cfg := Config{
		Port:          getEnv("PORT", "8080"),
		DatabasePath:  getEnv("DATABASE_PATH", ""),
		BaseAddress:   getEnv("BASE_ADDRESS", "localhost"),
		AppEnv:        getEnv("APP_ENV", "development"),
		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "*"),
		MigrationsDir: getEnv("MIGRATIONS_DIR", ""),
	}

	flags := flag.NewFlagSet("server", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&cfg.Port, "port", cfg.Port, "application port")
	flags.StringVar(&cfg.DatabasePath, "database-path", cfg.DatabasePath, "database file path")
	flags.StringVar(&cfg.BaseAddress, "base-address", cfg.BaseAddress, "application base address")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("parse command-line options: %w", err)
	}

	if err := validate(cfg); err != nil {
		return err
	}

	App = cfg
	return nil
}

// getEnv returns the environment variable value or the provided fallback.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func validate(cfg Config) error {
	var missing []string
	if cfg.DatabasePath == "" {
		missing = append(missing, "DATABASE_PATH")
	}
	if cfg.MigrationsDir == "" {
		missing = append(missing, "MIGRATIONS_DIR")
	}
	if len(missing) > 0 {
		return fmt.Errorf("required configuration is not set: %v", missing)
	}
	return nil
}
