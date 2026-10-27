package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

// EnvPrefix is prepended to every setting name (APP_ADDR, APP_LOG_LEVEL, ...).
const EnvPrefix = "APP"

// Config is process configuration loaded from the environment.
type Config struct {
	Addr            string
	ShutdownTimeout time.Duration
	LogLevel        slog.Level
}

// Load reads APP_* environment variables and applies defaults.
func Load() (Config, error) {
	c := Config{
		Addr:            getenv("APP_ADDR", ":8080"),
		ShutdownTimeout: 10 * time.Second,
		LogLevel:        slog.LevelInfo,
	}

	if v := os.Getenv("APP_SHUTDOWN_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return Config{}, fmt.Errorf("APP_SHUTDOWN_TIMEOUT: %w", err)
		}
		c.ShutdownTimeout = d
	}

	if v := os.Getenv("APP_LOG_LEVEL"); v != "" {
		var level slog.Level
		if err := level.UnmarshalText([]byte(v)); err != nil {
			return Config{}, fmt.Errorf("APP_LOG_LEVEL: %w", err)
		}
		c.LogLevel = level
	}

	return c, c.Validate()
}

// Validate returns an error if the config cannot be used to start the process.
func (c Config) Validate() error {
	if strings.TrimSpace(c.Addr) == "" {
		return fmt.Errorf("APP_ADDR is empty")
	}
	if c.ShutdownTimeout <= 0 {
		return fmt.Errorf("APP_SHUTDOWN_TIMEOUT must be > 0")
	}
	return nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
