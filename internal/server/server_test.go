package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ezphlow/go_template/internal/config"
	"github.com/ezphlow/go_template/internal/hello"
	"github.com/ezphlow/go_template/internal/mocks"
	"github.com/ezphlow/go_template/internal/models"
)

func newTestServer(t *testing.T, g models.Greeter) *Server {
	t.Helper()
	if g == nil {
		g = hello.Service{}
	}
	return New(config.Config{Addr: ":0"}, slog.New(slog.NewTextHandler(io.Discard, nil)), g)
}

func TestHello(t *testing.T) {
	s := newTestServer(t, hello.Service{})

	rec := httptest.NewRecorder()
	s.http.Handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/hello", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var body models.Hello
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json: %v", err)
	}
	if body.Message != "hello, world" {
		t.Errorf("message = %q, want hello, world", body.Message)
	}

	rec = httptest.NewRecorder()
	s.http.Handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/hello?name=ez", nil))
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json: %v", err)
	}
	if body.Message != "hello, ez" {
		t.Errorf("message = %q, want hello, ez", body.Message)
	}
}

func TestHelloUsesGreeter(t *testing.T) {
	s := newTestServer(t, mocks.Greeter{
		GreetFn: func(name string) models.Hello {
			return models.Hello{Message: "mocked:" + name}
		},
	})
	rec := httptest.NewRecorder()
	s.http.Handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/hello?name=ez", nil))
	var body models.Hello
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json: %v", err)
	}
	if body.Message != "mocked:ez" {
		t.Errorf("message = %q, want mocked:ez", body.Message)
	}
}

func TestHealthzMounted(t *testing.T) {
	s := newTestServer(t, nil)
	rec := httptest.NewRecorder()
	s.http.Handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}
