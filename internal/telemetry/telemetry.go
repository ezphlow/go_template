package telemetry

import (
	"log/slog"
	"net/http"
	"strings"

	httptrace "github.com/DataDog/dd-trace-go/contrib/net/http/v2"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"

	"github.com/ezphlow/go_template/internal/config"
)

// Start initializes the Datadog tracer (APM + runtime metrics).
// The returned func stops the tracer; it is safe to call when disabled.
func Start(cfg config.Telemetry, log *slog.Logger) (func(), error) {
	if log == nil {
		log = slog.Default()
	}
	if !cfg.Enabled {
		log.Info("datadog tracer disabled")
		return func() {}, nil
	}

	opts := []tracer.StartOption{
		tracer.WithService(cfg.Service),
		tracer.WithEnv(cfg.Env),
		tracer.WithRuntimeMetrics(),
	}
	if cfg.Version != "" {
		opts = append(opts, tracer.WithServiceVersion(cfg.Version))
	}
	if addr := strings.TrimSpace(cfg.AgentAddr); addr != "" {
		opts = append(opts, tracer.WithAgentAddr(normalizeAgentAddr(addr)))
	}

	if err := tracer.Start(opts...); err != nil {
		return func() {}, err
	}
	log.Info("datadog tracer started",
		"service", cfg.Service,
		"env", cfg.Env,
		"version", cfg.Version,
	)
	return tracer.Stop, nil
}

// WrapHTTP attaches Datadog APM to the server mux. Health probes are ignored.
func WrapHTTP(h http.Handler, cfg config.Telemetry) http.Handler {
	if !cfg.Enabled || h == nil {
		return h
	}
	opts := []httptrace.Option{
		httptrace.WithIgnoreRequest(func(r *http.Request) bool {
			switch r.URL.Path {
			case "/healthz", "/readyz":
				return true
			default:
				return false
			}
		}),
	}
	return httptrace.WrapHandler(h, cfg.Service, "", opts...)
}

func normalizeAgentAddr(addr string) string {
	if strings.Contains(addr, ":") {
		return addr
	}
	return addr + ":8126"
}
