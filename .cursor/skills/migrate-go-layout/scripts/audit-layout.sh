#!/usr/bin/env bash
# Report how close a Go repo is to go_template layout. Exit 1 if required pieces are missing.
set -euo pipefail

root="${1:-.}"
root="$(cd "$root" && pwd)"
fail=0

ok() { printf '  ok   %s\n' "$1"; }
bad() { printf '  miss %s\n' "$1"; fail=1; }
info() { printf '  info %s\n' "$1"; }

echo "audit: $root"

if [[ -f "$root/go.mod" ]]; then
  ok "go.mod ($(awk '/^module /{print $2}' "$root/go.mod"))"
else
  bad "go.mod"
fi

mains="$(find "$root" -name vendor -prune -o -name '*.go' -print 2>/dev/null | xargs grep -l '^func main(' 2>/dev/null || true)"
if echo "$mains" | grep -q '/cmd/[^/]\+/'; then
  ok "package main under cmd/<name>/"
elif [[ -n "$mains" ]]; then
  bad "package main should live under cmd/<name>/ (found: $(echo "$mains" | tr '\n' ' '))"
else
  info "no func main (library module?)"
fi

if [[ -f "$root/main.go" ]]; then
  bad "root main.go should move to cmd/<name>/main.go"
fi

if [[ -d "$root/internal" ]]; then
  ok "internal/"
else
  bad "internal/ (private service code)"
fi

if [[ -f "$root/Makefile" ]]; then
  if grep -qE '^lint:' "$root/Makefile" && grep -qE '^hooks:' "$root/Makefile"; then
    ok "Makefile has lint + hooks"
  else
    bad "Makefile missing lint: and/or hooks: targets"
  fi
else
  bad "Makefile"
fi

if [[ -f "$root/.golangci.yml" ]] || [[ -f "$root/.golangci.yaml" ]]; then
  ok ".golangci.yml"
else
  bad ".golangci.yml"
fi

if [[ -x "$root/.githooks/post-commit" ]]; then
  ok ".githooks/post-commit (executable)"
else
  bad ".githooks/post-commit executable"
fi

if [[ -x "$root/.githooks/pre-push" ]]; then
  ok ".githooks/pre-push (executable)"
else
  bad ".githooks/pre-push executable"
fi

if [[ -d "$root/.git" ]]; then
  hooks_path="$(git -C "$root" config --get core.hooksPath || true)"
  if [[ "$hooks_path" == ".githooks" ]]; then
    ok "git core.hooksPath=.githooks"
  else
    bad "git core.hooksPath is '${hooks_path:-unset}' (run: make hooks)"
  fi
else
  info "not a git repo — skip hooksPath check"
fi

if [[ -f "$root/.github/workflows/ci.yml" ]] || [[ -f "$root/.github/workflows/ci.yaml" ]]; then
  ok ".github/workflows/ci.yml"
else
  info "no .github/workflows/ci.yml"
fi

if [[ -d "$root/.cursor/skills/migrate-go-layout" ]]; then
  ok ".cursor/skills/migrate-go-layout"
else
  info "skill not copied into this repo"
fi

echo
if [[ "$fail" -ne 0 ]]; then
  echo "audit: FAIL (missing required layout pieces)"
  exit 1
fi
echo "audit: PASS"
exit 0
