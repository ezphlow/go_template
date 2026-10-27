package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ezphlow/go_template/internal/config"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	return New(config.Config{Addr: ":0"}, slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func TestHello(t *testing.T) {
	s := newTestServer(t)

	rec := httptest.NewRecorder()
	s.http.Handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/hello", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json: %v", err)
	}
	if body["message"] != "hello, world" {
		t.Errorf("message = %q, want hello, world", body["message"])
	}

	rec = httptest.NewRecorder()
	s.http.Handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/hello?name=ez", nil))
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json: %v", err)
	}
	if body["message"] != "hello, ez" {
		t.Errorf("message = %q, want hello, ez", body["message"])
	}
}

func TestHealthzMounted(t *testing.T) {
	s := newTestServer(t)
	rec := httptest.NewRecorder()
	s.http.Handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}
