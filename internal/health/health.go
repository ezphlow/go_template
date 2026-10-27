package health

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

// Checker reports liveness and readiness.
type Checker struct {
	started time.Time
	ready   atomic.Bool
}

// New returns a checker that is live immediately and ready after Ready().
func New() *Checker {
	c := &Checker{started: time.Now().UTC()}
	return c
}

// Ready marks the process as able to take traffic.
func (c *Checker) Ready() {
	c.ready.Store(true)
}

// NotReady marks the process as unable to take traffic (e.g. during shutdown).
func (c *Checker) NotReady() {
	c.ready.Store(false)
}

// Live is GET /healthz — process is up.
func (c *Checker) Live(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"started": c.started.Format(time.RFC3339),
	})
}

// Readyz is GET /readyz — process can serve traffic.
func (c *Checker) Readyz(w http.ResponseWriter, _ *http.Request) {
	if !c.ready.Load() {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"status": "not_ready",
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ready",
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
