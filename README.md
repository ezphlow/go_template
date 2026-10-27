# go_template

Sample Go HTTP service and the layout to copy (or migrate) other Golang repos onto.

```
cmd/<name>/main.go     thin entrypoint
internal/config        APP_* env
internal/models        shared types + interfaces
internal/hello         production implementations
internal/mocks         test doubles (XxxFn stubs)
internal/server        HTTP + graceful shutdown
internal/health        /healthz, /readyz
internal/telemetry     Datadog tracer (dd-trace-go v2)
.githooks              post-commit lint, pre-push test
.golangci.yml          golangci-lint
.cursor/skills         migrate-go-layout skill
```

## Run

```sh
make hooks    # once per clone: lint after commit, tests on push
make test
make serve    # :8080  GET /healthz  GET /v1/hello?name=ez
```

Env: `APP_ADDR` (default `:8080`), `APP_LOG_LEVEL` (`debug`/`info`/`warn`/`error`), `APP_SHUTDOWN_TIMEOUT` (default `10s`).

Datadog (started in `server.New`): `APP_DD_ENABLED` / `DD_TRACE_ENABLED` (default `true`), `APP_DD_SERVICE` / `DD_SERVICE` (default `go_template`), `APP_DD_ENV` / `DD_ENV` (default `dev`), `APP_DD_AGENT_ADDR` / `DD_AGENT_HOST` (default `localhost:8126`). Without an agent, traces are dropped. Tests construct `config.Config` with `Telemetry.Enabled` false.

Skip hooks: `SKIP_GIT_HOOKS=1`. `git commit --no-verify` does **not** skip `post-commit`.

## New service from this repo

1. Clone or copy this tree.
2. Rename the module in `go.mod`, `.golangci.yml` `local-prefixes`, and `Makefile` `MODULE`.
3. Rename `cmd/service` if you want a different binary name.
4. `make hooks tidy test`

## Migrate an existing Go repo

Use the Cursor skill [`.cursor/skills/migrate-go-layout/SKILL.md`](.cursor/skills/migrate-go-layout/SKILL.md). Copy it to `~/.cursor/skills/migrate-go-layout/` so it is available in other workspaces:

```sh
mkdir -p ~/.cursor/skills
cp -R .cursor/skills/migrate-go-layout ~/.cursor/skills/
```

Then in the target repo: *migrate this Go repo to the go_template layout*.

Audit a repo without changing it:

```sh
.cursor/skills/migrate-go-layout/scripts/audit-layout.sh /path/to/repo
```

## Commands

| Target | What |
| --- | --- |
| `make lint` | golangci-lint |
| `make test` | `go test ./... -race` |
| `make build` | `bin/service` |
| `make fmt` | gofmt + goimports |
| `make hooks` | `git config core.hooksPath .githooks` |
| `make ci` | lint + test + build |
| `make docker` | local image `go-template:local` |

Requires Go 1.26+ (for `dd-trace-go` v2) and [golangci-lint](https://golangci-lint.run/welcome/install/) v2.15.x on `$PATH`.
