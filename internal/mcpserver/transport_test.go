package mcpserver

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestParsePort(t *testing.T) {
	t.Parallel()

	port, err := parsePort([]string{"--port", "9000"})
	if err != nil {
		t.Fatalf("parsePort returned error: %v", err)
	}
	if port != 9000 {
		t.Fatalf("expected port 9000, got %d", port)
	}
}

func TestEnvConfig(t *testing.T) {
	t.Setenv("ARGOCD_BASE_URL", "https://argocd.example.com")
	t.Setenv("ARGOCD_API_TOKEN", "secret")

	cfg, err := envConfig()
	if err != nil {
		t.Fatalf("envConfig returned error: %v", err)
	}
	if cfg.ArgoCDBaseURL != "https://argocd.example.com" {
		t.Fatalf("unexpected base url: %q", cfg.ArgoCDBaseURL)
	}
	if cfg.ArgoCDAPIToken != "secret" {
		t.Fatalf("unexpected token: %q", cfg.ArgoCDAPIToken)
	}
}

func TestEnvConfigRequiresVariables(t *testing.T) {
	oldBaseURL := os.Getenv("ARGOCD_BASE_URL")
	oldToken := os.Getenv("ARGOCD_API_TOKEN")
	defer func() {
		_ = os.Setenv("ARGOCD_BASE_URL", oldBaseURL)
		_ = os.Setenv("ARGOCD_API_TOKEN", oldToken)
	}()
	_ = os.Unsetenv("ARGOCD_BASE_URL")
	_ = os.Unsetenv("ARGOCD_API_TOKEN")

	if _, err := envConfig(); err == nil {
		t.Fatal("expected envConfig to fail when env vars are missing")
	}
}

func TestHTTPHandlerMountsAtMCP(t *testing.T) {
	t.Parallel()

	handler := newHTTPHandlerForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
}

func TestSSEHandlerRequiresEnvConfig(t *testing.T) {
	oldBaseURL := os.Getenv("ARGOCD_BASE_URL")
	oldToken := os.Getenv("ARGOCD_API_TOKEN")
	defer func() {
		_ = os.Setenv("ARGOCD_BASE_URL", oldBaseURL)
		_ = os.Setenv("ARGOCD_API_TOKEN", oldToken)
	}()
	_ = os.Unsetenv("ARGOCD_BASE_URL")
	_ = os.Unsetenv("ARGOCD_API_TOKEN")

	if _, err := newSSEHandler(); err == nil {
		t.Fatal("expected newSSEHandler to fail without env config")
	}
}
