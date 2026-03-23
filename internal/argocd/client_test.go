package argocd

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListApplicationsStripsAndPaginates(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{
			"items":[
				{"metadata":{"name":"app-a","namespace":"argocd","labels":{"env":"dev"},"creationTimestamp":"2026-01-01"},"spec":{"project":"default","source":{"repoURL":"a"},"destination":{"server":"s"}},"status":{"sync":{"status":"Synced"},"health":{"status":"Healthy"},"summary":{"images":["a"]}}},
				{"metadata":{"name":"app-b","namespace":"argocd","labels":{"env":"prod"},"creationTimestamp":"2026-01-02"},"spec":{"project":"default","source":{"repoURL":"b"},"destination":{"server":"s"}},"status":{"sync":{"status":"OutOfSync"},"health":{"status":"Degraded"},"summary":{"images":["b"]}}}
			],
			"metadata":{"resourceVersion":"123"}
		}`)
	}))
	defer server.Close()

	client := NewClient(server.URL, "token")
	outAny, err := client.ListApplications(context.Background(), map[string]any{
		"limit":  1,
		"offset": 1,
	})
	if err != nil {
		t.Fatalf("ListApplications returned error: %v", err)
	}

	out := outAny.(map[string]any)
	items := out["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	item := items[0].(map[string]any)
	meta := item["metadata"].(map[string]any)
	if meta["name"] != "app-b" {
		t.Fatalf("unexpected app name: %#v", meta["name"])
	}
	resultMeta := out["metadata"].(map[string]any)
	if resultMeta["totalItems"].(int) != 2 {
		t.Fatalf("unexpected totalItems: %#v", resultMeta["totalItems"])
	}
	if resultMeta["hasMore"].(bool) {
		t.Fatal("expected hasMore to be false")
	}
}

func TestDeleteApplicationPassesQueryParameters(t *testing.T) {
	t.Parallel()

	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	defer server.Close()

	client := NewClient(server.URL, "token")
	_, err := client.DeleteApplication(context.Background(), "demo", map[string]any{
		"appNamespace":      "argocd",
		"cascade":           true,
		"propagationPolicy": "Foreground",
	})
	if err != nil {
		t.Fatalf("DeleteApplication returned error: %v", err)
	}
	expected := "appNamespace=argocd&cascade=true&propagationPolicy=Foreground"
	if gotQuery != expected {
		t.Fatalf("unexpected query: %q", gotQuery)
	}
}
