package server

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/ezphlow/go_template/internal/config"
	"github.com/ezphlow/go_template/internal/health"
)

// Server is the HTTP process: mux, middleware, and graceful shutdown.
type Server struct {
	cfg    config.Config
	log    *slog.Logger
	health *health.Checker
	http   *http.Server
}

// New wires routes and middleware. Call Run to listen.
func New(cfg config.Config, log *slog.Logger) *Server {
	if log == nil {
		log = slog.Default()
	}
	h := health.New()
	s := &Server{cfg: cfg, log: log, health: h}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.Live)
	mux.HandleFunc("/readyz", h.Readyz)
	mux.HandleFunc("/v1/hello", s.hello)

	s.http = &http.Server{
		Addr:              cfg.Addr,
		Handler:           s.withLogging(mux),
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
			return err
		}
		return <-errCh
	case err := <-errCh:
		return err
	}
}

func (s *Server) hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "hello, " + name,
	})
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
