package argocd

import (
	"context"
)

type Client struct {
	client *HTTPClient
}

func NewClient(baseURL, apiToken string) *Client {
	return &Client{client: NewHTTPClient(baseURL, apiToken)}
}

func (c *Client) ListApplications(ctx context.Context, params map[string]any) (any, error) {
	query := map[string]any{}
	if search, ok := params["search"].(string); ok && search != "" {
		query["search"] = search
	}
	resp, err := c.client.Get(ctx, "/api/v1/applications", query)
	if err != nil {
		return nil, err
	}
	body, _ := resp.Body.(map[string]any)
	rawItems, _ := body["items"].([]any)
	stripped := make([]any, 0, len(rawItems))
	for _, rawItem := range rawItems {
		app, _ := rawItem.(map[string]any)
		metadata, _ := app["metadata"].(map[string]any)
		spec, _ := app["spec"].(map[string]any)
		status, _ := app["status"].(map[string]any)
		stripped = append(stripped, map[string]any{
			"metadata": map[string]any{
				"name":              metadata["name"],
				"namespace":         metadata["namespace"],
				"labels":            metadata["labels"],
				"creationTimestamp": metadata["creationTimestamp"],
			},
			"spec": map[string]any{
				"project":     spec["project"],
				"source":      spec["source"],
				"destination": spec["destination"],
			},
			"status": map[string]any{
				"sync":    status["sync"],
				"health":  status["health"],
				"summary": status["summary"],
			},
		})
	}

	start := asInt(params["offset"])
	end := len(stripped)
	if limit := asInt(params["limit"]); limit > 0 {
		end = start + limit
	}
	if start > len(stripped) {
		start = len(stripped)
	}
	if end > len(stripped) {
		end = len(stripped)
	}
	items := stripped[start:end]

	meta, _ := body["metadata"].(map[string]any)
	return map[string]any{
		"items": items,
		"metadata": map[string]any{
			"resourceVersion": meta["resourceVersion"],
			"totalItems":      len(stripped),
			"returnedItems":   len(items),
			"hasMore":         end < len(stripped),
		},
	}, nil
}

func (c *Client) GetApplication(ctx context.Context, applicationName, appNamespace string) (any, error) {
	var query map[string]any
	if appNamespace != "" {
		query = map[string]any{"appNamespace": appNamespace}
	}
	resp, err := c.client.Get(ctx, "/api/v1/applications/"+applicationName, query)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *Client) CreateApplication(ctx context.Context, application any) (any, error) {
	resp, err := c.client.Post(ctx, "/api/v1/applications", nil, application)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *Client) UpdateApplication(ctx context.Context, applicationName string, application any) (any, error) {
	resp, err := c.client.Put(ctx, "/api/v1/applications/"+applicationName, nil, application)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *Client) DeleteApplication(ctx context.Context, applicationName string, options map[string]any) (any, error) {
	resp, err := c.client.Delete(ctx, "/api/v1/applications/"+applicationName, options)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *Client) SyncApplication(ctx context.Context, applicationName string, options map[string]any) (any, error) {
	resp, err := c.client.Post(ctx, "/api/v1/applications/"+applicationName+"/sync", nil, options)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *Client) GetApplicationResourceTree(ctx context.Context, applicationName string) (any, error) {
	resp, err := c.client.Get(ctx, "/api/v1/applications/"+applicationName+"/resource-tree", nil)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *Client) GetApplicationManagedResources(ctx context.Context, applicationName string, filters map[string]any) (any, error) {
	resp, err := c.client.Get(ctx, "/api/v1/applications/"+applicationName+"/managed-resources", filters)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *Client) GetWorkloadLogs(ctx context.Context, applicationName, applicationNamespace string, resourceRef map[string]any, container string) (any, error) {
	logs := make([]any, 0)
	err := c.client.GetStream(ctx, "/api/v1/applications/"+applicationName+"/logs", map[string]any{
		"appNamespace": applicationNamespace,
		"namespace":    resourceRef["namespace"],
		"resourceName": resourceRef["name"],
		"group":        resourceRef["group"],
		"kind":         resourceRef["kind"],
		"version":      resourceRef["version"],
		"follow":       false,
		"tailLines":    100,
		"container":    container,
	}, func(chunk any) {
		logs = append(logs, chunk)
	})
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (c *Client) GetApplicationEvents(ctx context.Context, applicationName string) (any, error) {
	resp, err := c.client.Get(ctx, "/api/v1/applications/"+applicationName+"/events", nil)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *Client) GetResource(ctx context.Context, applicationName, applicationNamespace string, resourceRef map[string]any) (any, error) {
	resp, err := c.client.Get(ctx, "/api/v1/applications/"+applicationName+"/resource", map[string]any{
		"appNamespace": applicationNamespace,
		"namespace":    resourceRef["namespace"],
		"resourceName": resourceRef["name"],
		"group":        resourceRef["group"],
		"kind":         resourceRef["kind"],
		"version":      resourceRef["version"],
	})
	if err != nil {
		return nil, err
	}
	body, _ := resp.Body.(map[string]any)
	return body["manifest"], nil
}

func (c *Client) GetResourceEvents(ctx context.Context, applicationName, applicationNamespace, resourceUID, resourceNamespace, resourceName string) (any, error) {
	resp, err := c.client.Get(ctx, "/api/v1/applications/"+applicationName+"/events", map[string]any{
		"appNamespace":      applicationNamespace,
		"resourceNamespace": resourceNamespace,
		"resourceUID":       resourceUID,
		"resourceName":      resourceName,
	})
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *Client) GetResourceActions(ctx context.Context, applicationName, applicationNamespace string, resourceRef map[string]any) (any, error) {
	resp, err := c.client.Get(ctx, "/api/v1/applications/"+applicationName+"/resource/actions", map[string]any{
		"appNamespace": applicationNamespace,
		"namespace":    resourceRef["namespace"],
		"resourceName": resourceRef["name"],
		"group":        resourceRef["group"],
		"kind":         resourceRef["kind"],
		"version":      resourceRef["version"],
	})
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *Client) RunResourceAction(ctx context.Context, applicationName, applicationNamespace string, resourceRef map[string]any, action string) (any, error) {
	resp, err := c.client.Post(ctx, "/api/v1/applications/"+applicationName+"/resource/actions", map[string]any{
		"appNamespace": applicationNamespace,
		"namespace":    resourceRef["namespace"],
		"resourceName": resourceRef["name"],
		"group":        resourceRef["group"],
		"kind":         resourceRef["kind"],
		"version":      resourceRef["version"],
	}, action)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func asInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case int32:
		return int(n)
	case int64:
		return int(n)
	case float64:
		return int(n)
	default:
		return 0
	}
}
