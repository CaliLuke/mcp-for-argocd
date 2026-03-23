package argocd

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/argoproj-labs/mcp-for-argocd/internal/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

type HTTPResponse struct {
	Status int
	Header http.Header
	Body   any
}

type HTTPClient struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

var tracer = otel.Tracer("github.com/argoproj-labs/mcp-for-argocd/internal/argocd")

func NewHTTPClient(baseURL, apiToken string) *HTTPClient {
	return &HTTPClient{
		baseURL:    baseURL,
		apiToken:   apiToken,
		httpClient: http.DefaultClient,
	}
}

func (c *HTTPClient) Get(ctx context.Context, path string, params map[string]any) (*HTTPResponse, error) {
	return c.request(ctx, http.MethodGet, path, params, nil)
}

func (c *HTTPClient) Post(ctx context.Context, path string, params map[string]any, body any) (*HTTPResponse, error) {
	return c.request(ctx, http.MethodPost, path, params, body)
}

func (c *HTTPClient) Put(ctx context.Context, path string, params map[string]any, body any) (*HTTPResponse, error) {
	return c.request(ctx, http.MethodPut, path, params, body)
}

func (c *HTTPClient) Delete(ctx context.Context, path string, params map[string]any) (*HTTPResponse, error) {
	return c.request(ctx, http.MethodDelete, path, params, nil)
}

func (c *HTTPClient) GetStream(ctx context.Context, path string, params map[string]any, cb func(any)) error {
	ctx, span := tracer.Start(ctx, "argocd.http.stream")
	defer span.End()

	endpoint, err := c.absURL(path, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "url_build_failed")
		logging.Logger.Error("argocd stream url build failed", "path", path, "error", err)
		return err
	}
	span.SetAttributes(
		attribute.String("http.method", http.MethodGet),
		attribute.String("url.full", endpoint.String()),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "request_build_failed")
		logging.Logger.Error("argocd stream request build failed", "method", http.MethodGet, "url", endpoint.String(), "error", err)
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "request_failed")
		logging.Logger.Error("argocd stream request failed", "method", http.MethodGet, "url", endpoint.String(), "error", err)
		return err
	}
	defer resp.Body.Close()
	span.SetAttributes(attribute.Int("http.response.status_code", resp.StatusCode))

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var item map[string]any
		if err := json.Unmarshal([]byte(line), &item); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "decode_failed")
			logging.Logger.Error("argocd stream decode failed", "url", endpoint.String(), "error", err)
			return err
		}
		cb(item["result"])
	}
	if err := scanner.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "stream_scan_failed")
		return err
	}
	span.SetStatus(codes.Ok, "")
	return nil
}

func (c *HTTPClient) request(ctx context.Context, method, path string, params map[string]any, body any) (*HTTPResponse, error) {
	ctx, span := tracer.Start(ctx, "argocd.http.request")
	defer span.End()

	endpoint, err := c.absURL(path, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "url_build_failed")
		logging.Logger.Error("argocd request url build failed", "path", path, "error", err)
		return nil, err
	}
	span.SetAttributes(
		attribute.String("http.method", method),
		attribute.String("url.full", endpoint.String()),
	)
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "body_encode_failed")
			logging.Logger.Error("argocd request body encode failed", "method", method, "url", endpoint.String(), "error", err)
			return nil, err
		}
		bodyReader = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), bodyReader)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "request_build_failed")
		logging.Logger.Error("argocd request build failed", "method", method, "url", endpoint.String(), "error", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "request_failed")
		logging.Logger.Error("argocd request failed", "method", method, "url", endpoint.String(), "error", err)
		return nil, err
	}
	defer resp.Body.Close()
	span.SetAttributes(attribute.Int("http.response.status_code", resp.StatusCode))

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "response_read_failed")
		logging.Logger.Error("argocd response read failed", "method", method, "url", endpoint.String(), "status", resp.StatusCode, "error", err)
		return nil, err
	}

	var parsed any
	if len(bytes.TrimSpace(raw)) > 0 {
		if err := json.Unmarshal(raw, &parsed); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "response_decode_failed")
			logging.Logger.Error("argocd response decode failed", "method", method, "url", endpoint.String(), "status", resp.StatusCode, "error", err)
			return nil, err
		}
	}

	if resp.StatusCode >= http.StatusBadRequest {
		span.SetStatus(codes.Error, http.StatusText(resp.StatusCode))
		logging.Logger.Warn("argocd request returned error status", "method", method, "url", endpoint.String(), "status", resp.StatusCode)
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return &HTTPResponse{
		Status: resp.StatusCode,
		Header: resp.Header.Clone(),
		Body:   parsed,
	}, nil
}

func (c *HTTPClient) absURL(path string, params map[string]any) (*url.URL, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		u, err := url.Parse(path)
		if err != nil {
			return nil, err
		}
		return u, nil
	}
	base, err := url.Parse(c.baseURL)
	if err != nil {
		logging.Logger.Error("invalid argocd base url", "base_url", c.baseURL, "error", err)
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	u, err := base.Parse(path)
	if err != nil {
		return nil, err
	}
	if params != nil {
		query := u.Query()
		for k, v := range params {
			if v == nil {
				continue
			}
			query.Set(k, fmt.Sprint(v))
		}
		u.RawQuery = query.Encode()
	}
	return u, nil
}
