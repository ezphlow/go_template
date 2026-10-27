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
	Telemetry       Telemetry
}

// Telemetry configures the Datadog tracer started with the HTTP server.
type Telemetry struct {
	Enabled   bool
	Service   string
	Env       string
	Version   string
	AgentAddr string
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

	c.Telemetry = Telemetry{
		Enabled:   parseBool(firstEnv("APP_DD_ENABLED", "DD_TRACE_ENABLED"), true),
		Service:   firstEnv("APP_DD_SERVICE", "DD_SERVICE", "go_template"),
		Env:       firstEnv("APP_DD_ENV", "DD_ENV", "dev"),
		Version:   firstEnv("APP_DD_VERSION", "DD_VERSION", ""),
		AgentAddr: firstEnv("APP_DD_AGENT_ADDR", "DD_AGENT_HOST", ""),
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

func firstEnv(keysAndFallback ...string) string {
	if len(keysAndFallback) == 0 {
		return ""
	}
	fallback := keysAndFallback[len(keysAndFallback)-1]
	for _, key := range keysAndFallback[:len(keysAndFallback)-1] {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}
	return fallback
}

func parseBool(v string, fallback bool) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}
