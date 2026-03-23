package argocd

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPRequestBuildsQueryAndHeaders(t *testing.T) {
	t.Parallel()

	var gotAuth string
	var gotPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.String()
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "secret-token")
	if _, err := client.Get(context.Background(), "/api/v1/applications", map[string]any{"search": "demo"}); err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	if gotAuth != "Bearer secret-token" {
		t.Fatalf("unexpected auth header: %q", gotAuth)
	}
	if gotPath != "/api/v1/applications?search=demo" {
		t.Fatalf("unexpected path: %q", gotPath)
	}
}

func TestHTTPGetStreamParsesLineDelimitedJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, "{\"result\":{\"line\":1}}\n{\"result\":{\"line\":2}}\n")
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "secret-token")
	var lines []any
	err := client.GetStream(context.Background(), "/stream", nil, func(chunk any) {
		lines = append(lines, chunk)
	})
	if err != nil {
		t.Fatalf("GetStream returned error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestHTTPClientReturnsDecodeErrorForInvalidJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, "not-json")
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "secret-token")
	_, err := client.Get(context.Background(), "/broken", nil)
	if err == nil || !strings.Contains(err.Error(), "invalid character") {
		t.Fatalf("unexpected error: %v", err)
	}
}
