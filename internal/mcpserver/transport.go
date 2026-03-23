package mcpserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/argoproj-labs/mcp-for-argocd/internal/logging"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("expected one of: stdio, sse, http")
	}

	switch args[0] {
	case "stdio":
		return connectStdio()
	case "sse":
		port, err := parsePort(args[1:])
		if err != nil {
			return err
		}
		return connectSSE(port)
	case "http":
		port, err := parsePort(args[1:])
		if err != nil {
			return err
		}
		return connectHTTP(port)
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func connectStdio() error {
	cfg, err := envConfig()
	if err != nil {
		return err
	}
	server, err := NewSDKServer(cfg)
	if err != nil {
		return err
	}
	logging.Logger.Info("connecting to stdio transport")
	return server.Server.Run(context.Background(), &sdkmcp.StdioTransport{})
}

func connectSSE(port int) error {
	handler, err := newSSEHandler()
	if err != nil {
		return err
	}
	addr := ":" + strconv.Itoa(port)
	logging.Logger.Info("connecting to sse transport", "port", port)
	return http.ListenAndServe(addr, instrumentHTTPHandler("argocd-mcp-sse", handler))
}

func newSSEHandler() (http.Handler, error) {
	cfg, err := envConfig()
	if err != nil {
		return nil, err
	}
	server, err := NewSDKServer(cfg)
	if err != nil {
		return nil, err
	}
	sseHandler := sdkmcp.NewSSEHandler(func(*http.Request) *sdkmcp.Server {
		return server.Server
	}, nil)
	mux := http.NewServeMux()
	mux.Handle("/sse", sseHandler)
	return mux, nil
}

func connectHTTP(port int) error {
	handler, err := newHTTPHandler()
	if err != nil {
		return err
	}
	addr := ":" + strconv.Itoa(port)
	logging.Logger.Info("connecting to http stream transport", "port", port)
	return http.ListenAndServe(addr, instrumentHTTPHandler("argocd-mcp-http", handler))
}

func newHTTPHandler() (http.Handler, error) {
	cfg, err := envConfig()
	if err != nil {
		return nil, err
	}
	server, err := NewSDKServer(cfg)
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.Handle("/mcp", server.Handler)
	return mux, nil
}

func parsePort(args []string) (int, error) {
	port := 3000
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--port", "-port":
			if i+1 >= len(args) {
				return 0, fmt.Errorf("missing value for --port")
			}
			value, err := strconv.Atoi(args[i+1])
			if err != nil {
				return 0, fmt.Errorf("invalid port: %w", err)
			}
			port = value
			i++
		default:
			return 0, fmt.Errorf("unknown argument: %s", args[i])
		}
	}
	return port, nil
}

func envConfig() (Config, error) {
	cfg := Config{
		ArgoCDBaseURL:  os.Getenv("ARGOCD_BASE_URL"),
		ArgoCDAPIToken: os.Getenv("ARGOCD_API_TOKEN"),
	}
	if cfg.ArgoCDBaseURL == "" || cfg.ArgoCDAPIToken == "" {
		return Config{}, fmt.Errorf("ARGOCD_BASE_URL and ARGOCD_API_TOKEN must be set")
	}
	return cfg, nil
}

func instrumentHTTPHandler(name string, next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, name)
}

func newHTTPHandlerForTest(h http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/mcp", h)
	return mux
}
