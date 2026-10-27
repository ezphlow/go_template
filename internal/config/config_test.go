package config

import (
	"log/slog"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("APP_ADDR", "")
	t.Setenv("APP_SHUTDOWN_TIMEOUT", "")
	t.Setenv("APP_LOG_LEVEL", "")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if c.Addr != ":8080" {
		t.Errorf("Addr = %q, want :8080", c.Addr)
	}
	if c.ShutdownTimeout != 10*time.Second {
		t.Errorf("ShutdownTimeout = %v, want 10s", c.ShutdownTimeout)
	}
	if c.LogLevel != slog.LevelInfo {
		t.Errorf("LogLevel = %v, want info", c.LogLevel)
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("APP_ADDR", ":9090")
	t.Setenv("APP_SHUTDOWN_TIMEOUT", "2s")
	t.Setenv("APP_LOG_LEVEL", "debug")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if c.Addr != ":9090" {
		t.Errorf("Addr = %q, want :9090", c.Addr)
	}
	if c.ShutdownTimeout != 2*time.Second {
		t.Errorf("ShutdownTimeout = %v, want 2s", c.ShutdownTimeout)
	}
	if c.LogLevel != slog.LevelDebug {
		t.Errorf("LogLevel = %v, want debug", c.LogLevel)
	}
}

func TestLoadInvalidTimeout(t *testing.T) {
	t.Setenv("APP_SHUTDOWN_TIMEOUT", "nope")
	if _, err := Load(); err == nil {
		t.Fatal("expected error for invalid APP_SHUTDOWN_TIMEOUT")
	}
}

func TestValidateEmptyAddr(t *testing.T) {
	c := Config{Addr: " ", ShutdownTimeout: time.Second}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for empty addr")
	}
}
