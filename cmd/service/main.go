package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ezphlow/go_template/internal/config"
	"github.com/ezphlow/go_template/internal/hello"
	"github.com/ezphlow/go_template/internal/server"
)

// Version is set at link time: -ldflags "-X main.Version=..."
var Version = "dev"

func main() {
	os.Exit(run())
}

func run() int {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "err", err)
		return 1
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))
	log = log.With("service", "go_template", "version", Version)
	if cfg.Telemetry.Version == "" {
		cfg.Telemetry.Version = Version
	}

	srv := server.New(cfg, log, hello.Service{})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := srv.Run(ctx); err != nil {
		log.Error("server exited", "err", err)
		return 1
	}
	return 0
}
