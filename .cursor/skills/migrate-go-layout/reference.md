# go_template file map

Canonical copies live in `/home/ezphlow/go/github.com/ezphlow/go_template`.

## Runtime

| Path | Role |
| --- | --- |
| `cmd/service/main.go` | Signals, config, `server.Run` |
| `internal/config` | `APP_ADDR`, `APP_LOG_LEVEL`, `APP_SHUTDOWN_TIMEOUT` |
| `internal/server` | `net/http` mux, request log, graceful shutdown |
| `internal/health` | `GET /healthz`, `GET /readyz` |

Routes use Go 1.22+ method patterns (`GET /healthz`). Health is live immediately; ready flips on listen / shutdown.

## Tooling

| Path | Role |
| --- | --- |
| `.golangci.yml` | v1 config; `gofmt` + `goimports` + staticcheck / revive / gocritic |
| `.githooks/pre-commit` | staged `.go` gofmt check + `golangci-lint run ./...` |
| `.githooks/pre-push` | `go test ./... -race` |
| `Makefile` | `hooks`, `lint`, `test`, `build`, `serve`, `ci` |
| `.github/workflows/ci.yml` | setup-go + golangci-lint-action v6 (v1.64.8) + test/build |

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
| Flat files next to `main.go` | packages under `internal/`, wired from `cmd/` |
| Already `cmd/` + `internal/` | add hooks + lint + Makefile targets only |
| Library-only module (no main) | keep packages; add lint/hooks/Makefile; skip `cmd/` unless they want a sample |

## Env

Prefer `APP_*`. If the existing service has a well-known prefix (`CAPE_`, `WORKFLOWER_`), keep it and document it in the repo rule instead of silently renaming production env vars.
