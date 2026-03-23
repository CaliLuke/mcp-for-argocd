# Repository Guidelines

## Project Structure & Module Organization

This repository is a Go implementation of an Argo CD MCP server built with `loom` and `loom-mcp`.

- `cmd/argocd-mcp/`: binary entrypoint.
- `design/`: design-first source of truth for generated APIs and MCP surface.
- `gen/`: generated code from `loom gen`. Do not hand-edit.
- `internal/argocd/`: Argo CD HTTP client and response shaping.
- `internal/mcpserver/`: service implementation and transport wiring.
- `internal/logging/`: shared logging setup.

Tests live next to the code they cover, for example `internal/mcpserver/http_integration_test.go`.

## Build, Test, and Development Commands

Use the `Makefile` as the default entrypoint:

- `make fmt`: format hand-written Go files with `gofmt`.
- `make lint`: run `golangci-lint`.
- `make complexity`: enforce `gocyclo` and `gocognit` thresholds.
- `make test`: run `go test ./...`.
- `make build`: build all packages.
- `make verify`: run formatting, vet, lint, complexity, tests, and build.
- `make generate`: regenerate `gen/` from `design/design.go`.

Run locally with:

- `go run ./cmd/argocd-mcp stdio`
- `go run ./cmd/argocd-mcp http --port 3000`

## Coding Style & Naming Conventions

Target Go 1.26. Keep code `gofmt`-clean and lint-clean. Use standard Go naming:

- exported identifiers: `CamelCase`
- internal helpers: `camelCase`
- tests: `TestXxx`

Prefer small functions, explicit error handling, and narrow packages. Keep framework-generated concerns in `gen/` and business logic in `internal/`.

## Testing Guidelines

Write table-driven Go tests where practical. Keep unit tests beside the package they exercise and add integration tests for MCP flows when transport behavior changes. Run `make verify` before opening a PR.

## Commit & Pull Request Guidelines

Recent history uses short imperative subjects and occasional conventional prefixes such as `fix:` and `chore:`. Follow that style:

- `fix: handle empty Argo CD responses`
- `chore: regenerate loom output`

PRs should include a short summary, the reason for the change, any config or behavior changes, and the verification performed. Link the relevant issue when one exists.

## Security & Configuration Tips

Set `ARGOCD_BASE_URL` and `ARGOCD_API_TOKEN` locally. Do not commit secrets. Use `MCP_READ_ONLY=true` when testing against shared Argo CD environments.
