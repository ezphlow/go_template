package telemetry

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ezphlow/go_template/internal/config"
)

func TestStartDisabled(t *testing.T) {
	stop, err := Start(config.Telemetry{Enabled: false}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("Start disabled: %v", err)
	}
	stop()
}

func TestWrapHTTPDisabledPassthrough(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	h := WrapHTTP(inner, config.Telemetry{Enabled: false})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/hello", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
}

func TestNormalizeAgentAddr(t *testing.T) {
	if got := normalizeAgentAddr("datadog-agent"); got != "datadog-agent:8126" {
		t.Errorf("host only = %q", got)
	}
	if got := normalizeAgentAddr("127.0.0.1:8126"); got != "127.0.0.1:8126" {
		t.Errorf("full addr = %q", got)
	}
}
