# go_template

Sample Go HTTP service and the layout to copy (or migrate) other Golang repos onto.

```
cmd/<name>/main.go     thin entrypoint
internal/config        APP_* env
internal/server        HTTP + graceful shutdown
internal/health        /healthz, /readyz
.githooks              pre-commit lint, pre-push test
.golangci.yml          golangci-lint v2
.cursor/skills         migrate-go-layout skill
```

## Run

```sh
make hooks    # once per clone: lint on commit, tests on push
make test
make serve    # :8080  GET /healthz  GET /v1/hello?name=ez
```

Env: `APP_ADDR` (default `:8080`), `APP_LOG_LEVEL` (`debug`/`info`/`warn`/`error`), `APP_SHUTDOWN_TIMEOUT` (default `10s`).

Skip hooks: `SKIP_GIT_HOOKS=1` or `git commit --no-verify`.

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

Requires [golangci-lint](https://golangci-lint.run/welcome/install/) v1.64.x on `$PATH` (Go 1.21+).
