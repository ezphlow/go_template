---
name: migrate-go-layout
description: >-
  Migrates an existing Golang repository onto the ezphlow go_template layout
  (cmd/internal, Makefile, golangci-lint, git lint hooks). Use when the user
  asks to migrate, convert, restyle, or bootstrap a Go repo to go_template;
  add lint git hooks; or align another Golang service with this skeleton.
---

# Migrate a Go repo to go_template

Template root (source of truth):

`/home/ezphlow/go/github.com/ezphlow/go_template`

Copy files from that tree. Do not invent a parallel layout.

Target layout:

```
cmd/<binary>/main.go
internal/<pkg>/          # config, server, domain — never a root main.go
pkg/                     # only if the module is a public library
Makefile                 # lint test build hooks
.golangci.yml
.githooks/pre-commit     # gofmt + golangci-lint
.githooks/pre-push       # go test -race
.github/workflows/ci.yml
.cursor/rules/go-service.mdc
.cursor/skills/migrate-go-layout/
```

## Workflow

Copy this checklist and tick as you go:

```
- [ ] 1. Inventory target
- [ ] 2. Plan mapping (show the user, then apply)
- [ ] 3. Move code into cmd/ + internal/
- [ ] 4. Drop in scaffolding (Makefile, lint, hooks, CI)
- [ ] 5. Fix module path + imports
- [ ] 6. Install hooks
- [ ] 7. Audit + lint + test
```

### 1. Inventory

In the target repo, collect:

- `go.mod` module path and Go version (raise to **1.21+** if lower)
- Where `package main` / `func main` live (root, `cmd/`, `main.go` + sibling packages)
- Existing Makefile / justfile / `.golangci.yml` / CI / git hooks
- Env prefix and health routes (`/healthz` is the convention)
- Tests: `go test ./...` baseline **before** moving files

Run the auditor (read-only):

```sh
TEMPLATE=/home/ezphlow/go/github.com/ezphlow/go_template
"$TEMPLATE/.cursor/skills/migrate-go-layout/scripts/audit-layout.sh" "$PWD"
```

### 2. Plan mapping

Write a short plan before moving files:

- Binary name → `cmd/<binary>/`
- Each existing top-level Go package → `internal/<same>` unless it is intentionally public (`pkg/`)
- Config → `internal/config` with `APP_*` (keep the old prefix only if the user forbids a rename)
- HTTP wiring → `internal/server` (or keep their router package and move it under `internal/`)
- Keep domain package names; do not flatten everything into one folder

If the repo is already close (has `cmd/` + `internal/`), skip the move and only add missing scaffolding.

### 3. Move code

- `git mv` when possible so history follows.
- Root `main.go` → `cmd/<binary>/main.go`. Leave `main` thin: load config, construct server, signals, `os.Exit`.
- Root-level library packages (`config`, `service`, `logger`, …) → `internal/<name>/`.
- Update import paths in the same change. Do not leave a transitional broken tree.
- Do not vendor. Do not add dependencies the service does not already use.

### 4. Scaffolding

Copy these from the template and substitute the **target module path** and **binary name**:

| File | Replace |
| --- | --- |
| `Makefile` | `MODULE`, `BINARY`, `CMD` |
| `.golangci.yml` | `local-prefixes` |
| `.githooks/pre-commit`, `.githooks/pre-push` | none (generic) |
| `.gitignore` | merge; keep their extras |
| `.editorconfig` | copy if missing |
| `.github/workflows/ci.yml` | copy if they have no Go CI |
| `Dockerfile` | `cmd/<binary>` path |
| `.cursor/rules/go-service.mdc` | copy |
| `.cursor/skills/migrate-go-layout/` | copy so the migrated repo stays the template |

If they already have a Makefile, **add** `lint`, `test`, `hooks`, `fmt` targets instead of overwriting custom deploy/docker targets.

If they already have `.golangci.yml`, enable the template's linter set (drop deprecated `ifshort`) and set `local-prefixes` to their module. This layout ships golangci-lint **v1** config (v1.64.x).

### 5. Module and imports

- Keep the existing module path unless the user asks to rename.
- `gofmt` / `goimports` with `-local <module>`.
- `go mod tidy`.

### 6. Install git hooks

```sh
make hooks
# git config core.hooksPath .githooks
```

Confirm `.githooks/pre-commit` and `pre-push` are executable. Do not use husky.

### 7. Verify

```sh
"$TEMPLATE/.cursor/skills/migrate-go-layout/scripts/audit-layout.sh" "$PWD"
make lint
make test
```

If `golangci-lint` is missing, install v1.64.8 (see template README) and re-run. Do not claim the migration is done while lint or tests fail.

If a check cannot run (no Go toolchain), say so and leave the files in place.

## Do not

- Rewrite business logic while migrating.
- Force `pkg/` for private code.
- Commit unless the user asks.
- Replace their Docker/K8s/deploy files except the build path (`cmd/...`).
- Add pre-commit (Python) or lefthook unless they already use them — this layout uses `core.hooksPath=.githooks`.

## New repo instead of migrate

If the user wants a new service, copy the template tree, then do step 4 substitutions only. Keep the sample `/v1/hello` until they replace it.

## Reference

- File-by-file notes: [reference.md](reference.md)
- Auditor: [scripts/audit-layout.sh](scripts/audit-layout.sh)
