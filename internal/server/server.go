package server

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"

	"github.com/ezphlow/go_template/internal/config"
	"github.com/ezphlow/go_template/internal/health"
	"github.com/ezphlow/go_template/internal/hello"
	"github.com/ezphlow/go_template/internal/models"
	"github.com/ezphlow/go_template/internal/telemetry"
)

// Server is the HTTP process: mux, middleware, and graceful shutdown.
type Server struct {
	cfg     config.Config
	log     *slog.Logger
	health  *health.Checker
	greeter models.Greeter
	http    *http.Server
	stopTel func()
}

// New wires routes and middleware. Call Run to listen.
func New(cfg config.Config, log *slog.Logger, greeter models.Greeter) *Server {
	if log == nil {
		log = slog.Default()
	}
	if greeter == nil {
		greeter = hello.Service{}
	}
	if cfg.Telemetry.Version == "" {
		cfg.Telemetry.Version = "dev"
	}

	h := health.New()
	s := &Server{cfg: cfg, log: log, health: h, greeter: greeter}

	stopTel, err := telemetry.Start(cfg.Telemetry, log)
	if err != nil {
		log.Error("start datadog tracer", "err", err)
		stopTel = func() {}
	}
	s.stopTel = stopTel

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.Live)
	mux.HandleFunc("/readyz", h.Readyz)
	mux.HandleFunc("/v1/hello", s.hello)

	s.http = &http.Server{
		Addr:              cfg.Addr,
		Handler:           telemetry.WrapHTTP(s.withLogging(mux), cfg.Telemetry),
		ReadHeaderTimeout: 5 * time.Second,
	}
	return s
}

// Run listens until ctx is cancelled, then shuts down.
func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		s.log.Info("listening", "addr", s.cfg.Addr)
		s.health.Ready()
		err := s.http.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		s.health.NotReady()
		s.log.Info("shutting down")
		shutCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()
		if err := s.http.Shutdown(shutCtx); err != nil {
			s.stopTelemetry()
			return err
		}
		s.stopTelemetry()
		return <-errCh
	case err := <-errCh:
		s.stopTelemetry()
		return err
	}
}

func (s *Server) stopTelemetry() {
	if s.stopTel != nil {
		s.stopTel()
		s.stopTel = nil
	}
}

func (s *Server) hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	span, _ := tracer.StartSpanFromContext(r.Context(), "hello.greet")
	defer span.Finish()
	span.SetTag("hello.name", name)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s.greeter.Greet(name))
}

func (s *Server) withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		s.log.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", sw.status,
			"dur_ms", time.Since(start).Milliseconds(),
		)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
