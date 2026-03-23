package mcpserver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGeneratedHTTPServerEndToEnd(t *testing.T) {
	argocdServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/applications":
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{
				"items": [
					{
						"metadata": {
							"name": "demo-app",
							"namespace": "argocd",
							"labels": {"env":"dev"},
							"creationTimestamp": "2026-01-01T00:00:00Z"
						},
						"spec": {
							"project": "default",
							"source": {"repoURL":"https://example.com/repo","path":"apps/demo","targetRevision":"main"},
							"destination": {"server":"https://kubernetes.default.svc","namespace":"demo"}
						},
						"status": {
							"sync": {"status":"Synced"},
							"health": {"status":"Healthy"},
							"summary": {"images":["demo:v1"]}
						}
					}
				],
				"metadata": {"resourceVersion":"1"}
			}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer argocdServer.Close()

	t.Setenv("ARGOCD_BASE_URL", argocdServer.URL)
	t.Setenv("ARGOCD_API_TOKEN", "test-token")

	handler, err := newHTTPHandler()
	if err != nil {
		t.Fatalf("newHTTPHandler returned error: %v", err)
	}

	mcpServer := httptest.NewServer(handler)
	defer mcpServer.Close()

	client := sdkmcp.NewClient(&sdkmcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err := client.Connect(ctx, &sdkmcp.StreamableClientTransport{
		Endpoint:             mcpServer.URL + "/mcp",
		HTTPClient:           mcpServer.Client(),
		DisableStandaloneSSE: true,
	}, nil)
	if err != nil {
		t.Fatalf("connect returned error: %v", err)
	}
	defer session.Close()

	tools, err := session.ListTools(ctx, &sdkmcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("ListTools returned error: %v", err)
	}

	foundListApplications := false
	for _, tool := range tools.Tools {
		if tool.Name == "list_applications" {
			foundListApplications = true
			break
		}
	}
	if !foundListApplications {
		t.Fatal("expected list_applications tool to be present")
	}

	callResult, err := session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name:      "list_applications",
		Arguments: map[string]any{"limit": 1},
	})
	if err != nil {
		t.Fatalf("CallTool returned error: %v", err)
	}
	if len(callResult.Content) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(callResult.Content))
	}

	textContent, ok := callResult.Content[0].(*sdkmcp.TextContent)
	if !ok {
		t.Fatalf("expected text content, got %T", callResult.Content[0])
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(textContent.Text), &payload); err != nil {
		t.Fatalf("failed to decode tool result: %v", err)
	}
	items := payload["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	item := items[0].(map[string]any)
	meta := item["metadata"].(map[string]any)
	if meta["name"] != "demo-app" {
		t.Fatalf("unexpected app name: %#v", meta["name"])
	}
}
