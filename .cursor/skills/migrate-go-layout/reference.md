# go_template file map

Canonical copies live in `/home/ezphlow/go/github.com/ezphlow/go_template`.

## Runtime

| Path | Role |
| --- | --- |
| `cmd/service/main.go` | Signals, config, `server.Run` |
| `internal/config` | `APP_ADDR`, `APP_LOG_LEVEL`, `APP_SHUTDOWN_TIMEOUT` |
| `internal/models` | Shared types + interfaces (`Hello`, `Greeter`) |
| `internal/hello` | Production implementations of those interfaces |
| `internal/mocks` | Hand-written `XxxFn` test doubles (`mocks.Greeter`) |
| `internal/server` | `net/http` mux, request log, graceful shutdown |
| `internal/health` | `GET /healthz`, `GET /readyz` |
| `internal/telemetry` | Datadog `dd-trace-go` v2: `Start`, HTTP wrap, runtime metrics |

Routes use Go 1.22+ method patterns (`GET /healthz`). Health is live immediately; ready flips on listen / shutdown.

## Tooling

| Path | Role |
| --- | --- |
| `.golangci.yml` | v2 config; formatters `gofmt` + `goimports`; staticcheck / revive / gocritic |
| `.githooks/post-commit` | gofmt + `golangci-lint run ./...` after the commit |
| `.githooks/pre-push` | `go test ./... -race` |
| `Makefile` | `hooks`, `lint`, `test`, `build`, `serve`, `ci` |
| `.github/workflows/ci.yml` | setup-go + golangci-lint-action v8 (v2.15.0) + test/build |

`make hooks` sets `core.hooksPath=.githooks` (repo-local, not global). Bypass: `SKIP_GIT_HOOKS=1`.

## Cursor

| Path | Role |
| --- | --- |
| `.cursor/rules/go-service.mdc` | always-on layout rule |
| `.cursor/skills/migrate-go-layout/` | this skill |

Copy the skill to `~/.cursor/skills/migrate-go-layout/` for other workspaces.

## Typical source shapes → target

| Before | After |
| --- | --- |
| `main.go` at repo root | `cmd/<name>/main.go` |
| `config/`, `service/`, `logger/` at root | `internal/config`, `internal/service`, `internal/logger` |
| `entities/` or `models/` at root | `internal/models` (types + interfaces only) |
| `mock/` or `mocks/` at root | `internal/mocks` (`XxxFn` stubs; not generated into `models`) |
| Flat files next to `main.go` | packages under `internal/`, wired from `cmd/` |
| Already `cmd/` + `internal/` | add hooks + lint + Makefile targets only |
| Library-only module (no main) | keep packages; add lint/hooks/Makefile; skip `cmd/` unless they want a sample |

## Env

Prefer `APP_*`. If the existing service has a well-known prefix (`CAPE_`, `WORKFLOWER_`), keep it and document it in the repo rule instead of silently renaming production env vars.
