package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLive(t *testing.T) {
	c := New()
	rec := httptest.NewRecorder()
	c.Live(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status = %v, want ok", body["status"])
	}
}

func TestReadyz(t *testing.T) {
	c := New()
	rec := httptest.NewRecorder()
	c.Readyz(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("before Ready() status = %d, want 503", rec.Code)
	}

	c.Ready()
	rec = httptest.NewRecorder()
	c.Readyz(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("after Ready() status = %d, want 200", rec.Code)
	}

	c.NotReady()
	rec = httptest.NewRecorder()
	c.Readyz(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("after NotReady() status = %d, want 503", rec.Code)
	}
}
