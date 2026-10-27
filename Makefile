GO              ?= go
MODULE          := github.com/ezphlow/go_template
BINARY          := bin/service
CMD             := ./cmd/service
GOPATH_BIN      := $(shell $(GO) env GOPATH)/bin
export PATH := $(GOPATH_BIN):$(PATH)
GOLANGCI_LINT   ?= golangci-lint
GOLANGCI_VERSION ?= v2.15.0

.PHONY: all hooks hooks-uninstall lint lint-fast fmt test build serve tidy docker ci clean

all: lint test build

## Git hooks (lint after commit, tests on push)
hooks:
	git config core.hooksPath .githooks
	chmod +x .githooks/post-commit .githooks/pre-push
	@echo "Installed git hooks via core.hooksPath=.githooks"

hooks-uninstall:
	git config --unset core.hooksPath || true
	@echo "Uninstalled repo git hooks"

lint:
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { \
		echo "golangci-lint not found. Install: https://golangci-lint.run/welcome/install/"; \
		echo "  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_VERSION)"; \
		exit 1; \
	}
	$(GOLANGCI_LINT) run ./...

lint-fast:
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { echo "golangci-lint required"; exit 1; }
	$(GOLANGCI_LINT) run ./...

fmt:
	$(GO) fmt ./...
	@command -v goimports >/dev/null 2>&1 && goimports -w -local $(MODULE) . || true

test:
	$(GO) test ./... -race -count=1 -timeout 60s

build:
	mkdir -p bin
	$(GO) build -ldflags "-X main.Version=$${VERSION:-dev}" -o $(BINARY) $(CMD)

serve: build
	./$(BINARY)

tidy:
	$(GO) mod tidy

docker:
	docker build -t go-template:local .

ci: lint test build

clean:
	rm -rf bin
	$(GO) clean ./...
