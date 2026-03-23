# Argo CD MCP Server

This repository is now a Go implementation of an Argo CD MCP server built on top of `loom` and `loom-mcp`.

It exposes Argo CD operations as MCP tools and uses the framework-generated MCP surface for:
- tool schemas
- strict argument validation
- MCP SDK server wiring
- adapter telemetry hooks

The hand-written code is intentionally narrow: Argo CD API calls, response shaping, and small transport bootstrapping.

## Features

- Go implementation with generated MCP surface from `loom` and `loom-mcp`
- `stdio`, streamable HTTP, and SSE transports
- OpenTelemetry-compatible HTTP instrumentation
- Read-only mode via `MCP_READ_ONLY=true`
- End-to-end tested generated HTTP MCP server path

## Tools

Application management:
- `list_applications`
- `get_application`
- `create_application`
- `update_application`
- `delete_application`
- `sync_application`

Resource management:
- `get_application_resource_tree`
- `get_application_managed_resources`
- `get_application_workload_logs`
- `get_application_events`
- `get_resource_events`
- `get_resources`
- `get_resource_actions`
- `run_resource_action`

## Configuration

Required environment variables:

```bash
export ARGOCD_BASE_URL="https://argocd.example.com"
export ARGOCD_API_TOKEN="your-token"
```

Optional environment variables:

```bash
export MCP_READ_ONLY="true"
export LOG_LEVEL="debug"
```

`MCP_READ_ONLY=true` removes these tools from the MCP server:
- `create_application`
- `update_application`
- `delete_application`
- `sync_application`
- `run_resource_action`

## Running

Build:

```bash
go build ./cmd/argocd-mcp
```

Run over stdio:

```bash
go run ./cmd/argocd-mcp stdio
```

Run over streamable HTTP on port `3000`:

```bash
go run ./cmd/argocd-mcp http --port 3000
```

Run over SSE on port `3000`:

```bash
go run ./cmd/argocd-mcp sse --port 3000
```

HTTP endpoints:
- streamable HTTP: `POST/GET/DELETE /mcp`
- SSE: `GET /sse`

## MCP Client Example

VS Code `mcp.json` using stdio:

```json
{
  "servers": {
    "argocd-mcp": {
      "type": "stdio",
      "command": "/absolute/path/to/argocd-mcp",
      "args": ["stdio"],
      "env": {
        "ARGOCD_BASE_URL": "https://argocd.example.com",
        "ARGOCD_API_TOKEN": "your-token"
      }
    }
  }
}
```

## Development

Generate framework code from the design:

```bash
loom gen github.com/argoproj-labs/mcp-for-argocd/design
```

Test:

```bash
make test
```

Build:

```bash
make build
```

Lint and verify:

```bash
make lint
make complexity
make verify
```

Container image:

```bash
docker build -t argocd-mcp .
```

## Architecture

Key pieces:
- [design/design.go](/Users/luca/code/mcp-for-argocd/design/design.go): design-first source of truth
- [gen/](/Users/luca/code/mcp-for-argocd/gen): generated `loom` and `loom-mcp` output
- [internal/mcpserver/service.go](/Users/luca/code/mcp-for-argocd/internal/mcpserver/service.go): generated service implementation adapter
- [internal/argocd/](/Users/luca/code/mcp-for-argocd/internal/argocd): Argo CD HTTP client and response shaping

## Notes

- This repo no longer depends on Node.js, TypeScript, or `pnpm`.
- The MCP contract is framework-generated; avoid hand-editing `gen/`.
- If you change the design, regenerate before committing.
