package mcpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	genargocd "github.com/argoproj-labs/mcp-for-argocd/gen/argocd"
	mcpargocd "github.com/argoproj-labs/mcp-for-argocd/gen/mcp_argocd"
	argocdclient "github.com/argoproj-labs/mcp-for-argocd/internal/argocd"
	"github.com/argoproj-labs/mcp-for-argocd/internal/logging"
)

type Config struct {
	ArgoCDBaseURL  string
	ArgoCDAPIToken string
}

type Service struct {
	client   *argocdclient.Client
	readOnly bool
}

func NewService(cfg Config) *Service {
	return &Service{
		client:   argocdclient.NewClient(cfg.ArgoCDBaseURL, cfg.ArgoCDAPIToken),
		readOnly: strings.EqualFold(strings.TrimSpace(os.Getenv("MCP_READ_ONLY")), "true"),
	}
}

func NewSDKServer(cfg Config) (*mcpargocd.SDKServer, error) {
	service := NewService(cfg)
	server, err := mcpargocd.NewSDKServer(service, sdkServerOptions())
	if err != nil {
		return nil, err
	}
	if service.readOnly {
		server.Server.RemoveTools(
			"create_application",
			"update_application",
			"delete_application",
			"sync_application",
			"run_resource_action",
		)
	}
	return server, nil
}

type contextKey string

const (
	contextKeyRequestMethod contextKey = "request_method"
	contextKeyRequestPath   contextKey = "request_path"
	contextKeySessionID     contextKey = "mcp_session_id"
)

func sdkServerOptions() *mcpargocd.SDKServerOptions {
	return &mcpargocd.SDKServerOptions{
		Adapter: &mcpargocd.MCPAdapterOptions{
			Logger: func(ctx context.Context, event string, details any) {
				logging.Logger.DebugContext(ctx, "mcp adapter event", "event", event, "details", details)
			},
			TelemetryName: "github.com/argoproj-labs/mcp-for-argocd/mcp",
			ToolCallInterceptors: []mcpargocd.ToolCallInterceptor{
				func(ctx context.Context, info mcpargocd.ToolCallInterceptorInfo, payload *mcpargocd.ToolsCallPayload, stream mcpargocd.ToolsCallServerStream, next mcpargocd.ToolCallHandler) (bool, error) {
					start := time.Now()
					toolName := ""
					if info != nil {
						toolName = info.Tool()
					}
					logging.Logger.InfoContext(ctx, "mcp tool call started",
						"tool", toolName,
						"request_method", ctx.Value(contextKeyRequestMethod),
						"request_path", ctx.Value(contextKeyRequestPath),
						"session_id", ctx.Value(contextKeySessionID),
					)
					toolErr, err := next(ctx, payload, stream)
					logging.Logger.InfoContext(ctx, "mcp tool call completed",
						"tool", toolName,
						"duration_ms", time.Since(start).Milliseconds(),
						"tool_error", toolErr,
						"error", errString(err),
					)
					return toolErr, err
				},
			},
		},
		RequestContext: func(ctx context.Context, r *http.Request) context.Context {
			if r == nil {
				return ctx
			}
			ctx = context.WithValue(ctx, contextKeyRequestMethod, r.Method)
			if r.URL != nil {
				ctx = context.WithValue(ctx, contextKeyRequestPath, r.URL.Path)
			}
			if sessionID := r.Header.Get("Mcp-Session-Id"); sessionID != "" {
				ctx = context.WithValue(ctx, contextKeySessionID, sessionID)
			}
			return ctx
		},
	}
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func (s *Service) ListApplications(ctx context.Context, p *genargocd.ListApplicationsPayload) (string, error) {
	params := map[string]any{}
	if p != nil {
		if p.Search != nil {
			params["search"] = *p.Search
		}
		if p.Limit != nil {
			params["limit"] = *p.Limit
		}
		if p.Offset != nil {
			params["offset"] = *p.Offset
		}
	}
	result, err := s.client.ListApplications(ctx, params)
	return stringify(result, err)
}

func (s *Service) GetApplication(ctx context.Context, p *genargocd.GetApplicationPayload) (string, error) {
	var namespace string
	if p != nil && p.ApplicationNamespace != nil {
		namespace = *p.ApplicationNamespace
	}
	result, err := s.client.GetApplication(ctx, p.ApplicationName, namespace)
	return stringify(result, err)
}

func (s *Service) GetApplicationResourceTree(ctx context.Context, p *genargocd.GetApplicationResourceTreePayload) (string, error) {
	result, err := s.client.GetApplicationResourceTree(ctx, p.ApplicationName)
	return stringify(result, err)
}

func (s *Service) GetApplicationManagedResources(ctx context.Context, p *genargocd.GetApplicationManagedResourcesPayload) (string, error) {
	filters := map[string]any{}
	if p != nil {
		if p.Kind != nil {
			filters["kind"] = *p.Kind
		}
		if p.Namespace != nil {
			filters["namespace"] = *p.Namespace
		}
		if p.Name != nil {
			filters["name"] = *p.Name
		}
		if p.Version != nil {
			filters["version"] = *p.Version
		}
		if p.Group != nil {
			filters["group"] = *p.Group
		}
		if p.AppNamespace != nil {
			filters["appNamespace"] = *p.AppNamespace
		}
		if p.Project != nil {
			filters["project"] = *p.Project
		}
	}
	result, err := s.client.GetApplicationManagedResources(ctx, p.ApplicationName, filters)
	return stringify(result, err)
}

func (s *Service) GetApplicationWorkloadLogs(ctx context.Context, p *genargocd.GetApplicationWorkloadLogsPayload) (string, error) {
	resourceRef, err := resourceRefToMap(p.ResourceRef)
	if err != nil {
		return "", err
	}
	result, err := s.client.GetWorkloadLogs(ctx, p.ApplicationName, p.ApplicationNamespace, resourceRef, p.Container)
	return stringify(result, err)
}

func (s *Service) GetApplicationEvents(ctx context.Context, p *genargocd.GetApplicationEventsPayload) (string, error) {
	result, err := s.client.GetApplicationEvents(ctx, p.ApplicationName)
	return stringify(result, err)
}

func (s *Service) GetResourceEvents(ctx context.Context, p *genargocd.GetResourceEventsPayload) (string, error) {
	result, err := s.client.GetResourceEvents(ctx, p.ApplicationName, p.ApplicationNamespace, p.ResourceUID, p.ResourceNamespace, p.ResourceName)
	return stringify(result, err)
}

func (s *Service) GetResources(ctx context.Context, p *genargocd.GetResourcesPayload) (string, error) {
	resourceRefs := make([]map[string]any, 0, len(p.ResourceRefs))
	for _, ref := range p.ResourceRefs {
		converted, err := resourceRefToMap(ref)
		if err != nil {
			return "", err
		}
		resourceRefs = append(resourceRefs, converted)
	}
	payload := map[string]any{
		"applicationName":      p.ApplicationName,
		"applicationNamespace": p.ApplicationNamespace,
	}
	if len(resourceRefs) > 0 {
		payload["resourceRefs"] = resourceRefs
	}

	applicationName := p.ApplicationName
	applicationNamespace := p.ApplicationNamespace
	if len(resourceRefs) == 0 {
		treeAny, err := s.client.GetApplicationResourceTree(ctx, applicationName)
		if err != nil {
			return "", err
		}
		tree, _ := treeAny.(map[string]any)
		nodes, _ := tree["nodes"].([]any)
		resourceRefs = make([]map[string]any, 0, len(nodes))
		for _, nodeAny := range nodes {
			node, _ := nodeAny.(map[string]any)
			resourceRefs = append(resourceRefs, map[string]any{
				"uid":       node["uid"],
				"version":   node["version"],
				"group":     node["group"],
				"kind":      node["kind"],
				"name":      node["name"],
				"namespace": node["namespace"],
			})
		}
	}
	out := make([]any, 0, len(resourceRefs))
	for _, ref := range resourceRefs {
		resource, err := s.client.GetResource(ctx, applicationName, applicationNamespace, ref)
		if err != nil {
			return "", err
		}
		out = append(out, resource)
	}
	return stringify(out, nil)
}

func (s *Service) GetResourceActions(ctx context.Context, p *genargocd.GetResourceActionsPayload) (string, error) {
	resourceRef, err := resourceRefToMap(p.ResourceRef)
	if err != nil {
		return "", err
	}
	result, err := s.client.GetResourceActions(ctx, p.ApplicationName, p.ApplicationNamespace, resourceRef)
	return stringify(result, err)
}

func (s *Service) CreateApplication(ctx context.Context, p *genargocd.CreateApplicationPayload) (string, error) {
	if s.readOnly {
		return "", errors.New("create_application is disabled in read-only mode")
	}
	result, err := s.client.CreateApplication(ctx, applicationToMap(p.Application))
	return stringify(result, err)
}

func (s *Service) UpdateApplication(ctx context.Context, p *genargocd.UpdateApplicationPayload) (string, error) {
	if s.readOnly {
		return "", errors.New("update_application is disabled in read-only mode")
	}
	result, err := s.client.UpdateApplication(ctx, p.ApplicationName, applicationToMap(p.Application))
	return stringify(result, err)
}

func (s *Service) DeleteApplication(ctx context.Context, p *genargocd.DeleteApplicationPayload) (string, error) {
	if s.readOnly {
		return "", errors.New("delete_application is disabled in read-only mode")
	}
	options := map[string]any{}
	if p.ApplicationNamespace != nil {
		options["appNamespace"] = *p.ApplicationNamespace
	}
	if p.Cascade != nil {
		options["cascade"] = *p.Cascade
	}
	if p.PropagationPolicy != nil {
		options["propagationPolicy"] = *p.PropagationPolicy
	}
	if len(options) == 0 {
		options = nil
	}
	result, err := s.client.DeleteApplication(ctx, p.ApplicationName, options)
	return stringify(result, err)
}

func (s *Service) SyncApplication(ctx context.Context, p *genargocd.SyncApplicationPayload) (string, error) {
	if s.readOnly {
		return "", errors.New("sync_application is disabled in read-only mode")
	}
	options := map[string]any{}
	if p.ApplicationNamespace != nil {
		options["appNamespace"] = *p.ApplicationNamespace
	}
	if p.DryRun != nil {
		options["dryRun"] = *p.DryRun
	}
	if p.Prune != nil {
		options["prune"] = *p.Prune
	}
	if p.Revision != nil {
		options["revision"] = *p.Revision
	}
	if len(p.SyncOptions) > 0 {
		options["syncOptions"] = p.SyncOptions
	}
	if len(options) == 0 {
		options = nil
	}
	result, err := s.client.SyncApplication(ctx, p.ApplicationName, options)
	return stringify(result, err)
}

func (s *Service) RunResourceAction(ctx context.Context, p *genargocd.RunResourceActionPayload) (string, error) {
	if s.readOnly {
		return "", errors.New("run_resource_action is disabled in read-only mode")
	}
	resourceRef, err := resourceRefToMap(p.ResourceRef)
	if err != nil {
		return "", err
	}
	result, err := s.client.RunResourceAction(ctx, p.ApplicationName, p.ApplicationNamespace, resourceRef, p.Action)
	return stringify(result, err)
}

func stringify(value any, err error) (string, error) {
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func resourceRefToMap(ref *genargocd.ResourceRef) (map[string]any, error) {
	if ref == nil {
		return nil, fmt.Errorf("resourceRef is required")
	}
	return map[string]any{
		"uid":       ref.UID,
		"kind":      ref.Kind,
		"namespace": ref.Namespace,
		"name":      ref.Name,
		"version":   ref.Version,
		"group":     ref.Group,
	}, nil
}

func applicationToMap(app *genargocd.Application) map[string]any {
	if app == nil {
		return nil
	}
	result := map[string]any{}
	if app.Metadata != nil {
		result["metadata"] = map[string]any{
			"name":      app.Metadata.Name,
			"namespace": app.Metadata.Namespace,
		}
	}
	if app.Spec != nil {
		spec := map[string]any{
			"project": app.Spec.Project,
		}
		if app.Spec.Source != nil {
			spec["source"] = map[string]any{
				"repoURL":        app.Spec.Source.RepoURL,
				"path":           app.Spec.Source.Path,
				"targetRevision": app.Spec.Source.TargetRevision,
			}
		}
		if app.Spec.SyncPolicy != nil {
			syncPolicy := map[string]any{
				"syncOptions": app.Spec.SyncPolicy.SyncOptions,
			}
			if app.Spec.SyncPolicy.Automated != nil {
				syncPolicy["automated"] = map[string]any{
					"prune":    app.Spec.SyncPolicy.Automated.Prune,
					"selfHeal": app.Spec.SyncPolicy.Automated.SelfHeal,
				}
			}
			if app.Spec.SyncPolicy.Retry != nil {
				retry := map[string]any{
					"limit": app.Spec.SyncPolicy.Retry.Limit,
				}
				if app.Spec.SyncPolicy.Retry.Backoff != nil {
					retry["backoff"] = map[string]any{
						"duration":    app.Spec.SyncPolicy.Retry.Backoff.Duration,
						"maxDuration": app.Spec.SyncPolicy.Retry.Backoff.MaxDuration,
						"factor":      app.Spec.SyncPolicy.Retry.Backoff.Factor,
					}
				}
				syncPolicy["retry"] = retry
			}
			spec["syncPolicy"] = syncPolicy
		}
		if app.Spec.Destination != nil {
			destination := map[string]any{}
			if app.Spec.Destination.Server != nil {
				destination["server"] = *app.Spec.Destination.Server
			}
			if app.Spec.Destination.Namespace != nil {
				destination["namespace"] = *app.Spec.Destination.Namespace
			}
			if app.Spec.Destination.Name != nil {
				destination["name"] = *app.Spec.Destination.Name
			}
			spec["destination"] = destination
		}
		result["spec"] = spec
	}
	return result
}
