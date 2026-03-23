SHELL := /bin/zsh

GO ?= go
GOLANGCI_LINT ?= golangci-lint
GO_PACKAGES := ./...
GO_DIRS := cmd design internal
GO_FILES := $(shell find $(GO_DIRS) -type f -name '*.go')

.PHONY: help fmt check-fmt vet lint complexity test build generate verify ci

help:
	@echo "Available targets:"
	@echo "  make fmt         - Format Go source files"
	@echo "  make check-fmt   - Fail if Go files are not formatted"
	@echo "  make vet         - Run go vet"
	@echo "  make lint        - Run golangci-lint"
	@echo "  make complexity  - Run complexity-focused linters"
	@echo "  make test        - Run Go tests"
	@echo "  make build       - Build all Go packages"
	@echo "  make generate    - Regenerate loom/loom-mcp code"
	@echo "  make verify      - Run the full verification suite"
	@echo "  make ci          - Alias for verify"

fmt:
	@gofmt -w $(GO_FILES)

check-fmt:
	@unformatted="$$(gofmt -l $(GO_FILES))"; \
	if [[ -n "$$unformatted" ]]; then \
		echo "These files need gofmt:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

vet:
	@$(GO) vet $(GO_PACKAGES)

lint:
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { echo "golangci-lint is required for 'make lint'"; exit 1; }
	@$(GOLANGCI_LINT) run $(GO_PACKAGES)

complexity:
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { echo "golangci-lint is required for 'make complexity'"; exit 1; }
	@$(GOLANGCI_LINT) run --enable-only gocyclo --enable-only gocognit $(GO_PACKAGES)

test:
	@$(GO) test $(GO_PACKAGES)

build:
	@$(GO) build $(GO_PACKAGES)

generate:
	@command -v loom >/dev/null 2>&1 || { echo "loom is required for 'make generate'"; exit 1; }
	@loom gen github.com/argoproj-labs/mcp-for-argocd/design

verify: check-fmt vet lint complexity test build

ci: verify
