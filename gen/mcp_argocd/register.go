package mcpargocd

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/CaliLuke/loom-mcp/runtime/agent/planner"
	agentsruntime "github.com/CaliLuke/loom-mcp/runtime/agent/runtime"
	"github.com/CaliLuke/loom-mcp/runtime/agent/telemetry"
	"github.com/CaliLuke/loom-mcp/runtime/agent/tools"
	mcpruntime "github.com/CaliLuke/loom-mcp/runtime/mcp"
	"github.com/CaliLuke/loom-mcp/runtime/mcp/retry"
	"github.com/modelcontextprotocol/go-sdk/jsonrpc"
)

// ArgocdArgocdMcpToolsetToolSpecs contains the tool specifications for the argocd-mcp toolset.
var ArgocdArgocdMcpToolsetToolSpecs = []tools.ToolSpec{
	{
		Name:        "list_applications",
		Service:     "argocd",
		Toolset:     "argocd.argocd-mcp",
		Description: "list_applications returns list of applications",
		Payload: tools.TypeSpec{
			Name:   "*argocd.ListApplicationsPayload",
			Schema: []byte("{\"type\":\"object\",\"properties\":{\"limit\":{\"type\":\"integer\",\"description\":\"Maximum number of applications to return. Use this to reduce token usage when there are many applications. Optional.\"},\"offset\":{\"type\":\"integer\",\"description\":\"Number of applications to skip before returning results. Use with limit for pagination. Optional.\"},\"search\":{\"type\":\"string\",\"description\":\"Search applications by name. This is a partial match on the application name and does not support glob patterns (e.g. \\\"*\\\"). Optional.\"}},\"additionalProperties\":false}"),
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
		Result: tools.TypeSpec{
			Name:   "string",
			Schema: nil,
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
	},
	{
		Name:        "get_application",
		Service:     "argocd",
		Toolset:     "argocd.argocd-mcp",
		Description: "get_application returns application by application name. Optionally specify the application namespace to get applications from non-default namespaces.",
		Payload: tools.TypeSpec{
			Name:   "*argocd.GetApplicationPayload",
			Schema: []byte("{\"type\":\"object\",\"required\":[\"applicationName\"],\"properties\":{\"applicationName\":{\"type\":\"string\"},\"applicationNamespace\":{\"type\":\"string\"}},\"additionalProperties\":false}"),
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
		Result: tools.TypeSpec{
			Name:   "string",
			Schema: nil,
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
	},
	{
		Name:        "get_application_resource_tree",
		Service:     "argocd",
		Toolset:     "argocd.argocd-mcp",
		Description: "get_application_resource_tree returns resource tree for application by application name",
		Payload: tools.TypeSpec{
			Name:   "*argocd.GetApplicationResourceTreePayload",
			Schema: []byte("{\"type\":\"object\",\"required\":[\"applicationName\"],\"properties\":{\"applicationName\":{\"type\":\"string\"}},\"additionalProperties\":false}"),
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
		Result: tools.TypeSpec{
			Name:   "string",
			Schema: nil,
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
	},
	{
		Name:        "get_application_managed_resources",
		Service:     "argocd",
		Toolset:     "argocd.argocd-mcp",
		Description: "get_application_managed_resources returns managed resources for application by application name with optional filtering. Use filters to avoid token limits with large applications. Examples: kind=\"ConfigMap\" for config maps only, namespace=\"production\" for specific namespace, or combine multiple filters.",
		Payload: tools.TypeSpec{
			Name:   "*argocd.GetApplicationManagedResourcesPayload",
			Schema: []byte("{\"type\":\"object\",\"required\":[\"applicationName\"],\"properties\":{\"appNamespace\":{\"type\":\"string\",\"description\":\"Filter by Argo CD application namespace\"},\"applicationName\":{\"type\":\"string\"},\"group\":{\"type\":\"string\",\"description\":\"Filter by API group\"},\"kind\":{\"type\":\"string\",\"description\":\"Filter by Kubernetes resource kind (e.g., \\\"ConfigMap\\\", \\\"Secret\\\", \\\"Deployment\\\")\"},\"name\":{\"type\":\"string\",\"description\":\"Filter by resource name\"},\"namespace\":{\"type\":\"string\",\"description\":\"Filter by Kubernetes namespace\"},\"project\":{\"type\":\"string\",\"description\":\"Filter by Argo CD project\"},\"version\":{\"type\":\"string\",\"description\":\"Filter by resource API version\"}},\"additionalProperties\":false}"),
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
		Result: tools.TypeSpec{
			Name:   "string",
			Schema: nil,
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
	},
	{
		Name:        "get_application_workload_logs",
		Service:     "argocd",
		Toolset:     "argocd.argocd-mcp",
		Description: "get_application_workload_logs returns logs for application workload (Deployment, StatefulSet, Pod, etc.) by application name and resource ref and optionally container name",
		Payload: tools.TypeSpec{
			Name:   "*argocd.GetApplicationWorkloadLogsPayload",
			Schema: []byte("{\"type\":\"object\",\"required\":[\"applicationName\",\"applicationNamespace\",\"resourceRef\",\"container\"],\"properties\":{\"applicationName\":{\"type\":\"string\"},\"applicationNamespace\":{\"type\":\"string\"},\"container\":{\"type\":\"string\"},\"resourceRef\":{\"type\":\"object\",\"required\":[\"uid\",\"kind\",\"namespace\",\"name\",\"version\",\"group\"],\"properties\":{\"group\":{\"type\":\"string\"},\"kind\":{\"type\":\"string\"},\"name\":{\"type\":\"string\"},\"namespace\":{\"type\":\"string\"},\"uid\":{\"type\":\"string\"},\"version\":{\"type\":\"string\"}},\"additionalProperties\":false}},\"additionalProperties\":false}"),
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
		Result: tools.TypeSpec{
			Name:   "string",
			Schema: nil,
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
	},
	{
		Name:        "get_application_events",
		Service:     "argocd",
		Toolset:     "argocd.argocd-mcp",
		Description: "get_application_events returns events for application by application name",
		Payload: tools.TypeSpec{
			Name:   "*argocd.GetApplicationEventsPayload",
			Schema: []byte("{\"type\":\"object\",\"required\":[\"applicationName\"],\"properties\":{\"applicationName\":{\"type\":\"string\"}},\"additionalProperties\":false}"),
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
		Result: tools.TypeSpec{
			Name:   "string",
			Schema: nil,
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
	},
	{
		Name:        "get_resource_events",
		Service:     "argocd",
		Toolset:     "argocd.argocd-mcp",
		Description: "get_resource_events returns events for a resource that is managed by an application",
		Payload: tools.TypeSpec{
			Name:   "*argocd.GetResourceEventsPayload",
			Schema: []byte("{\"type\":\"object\",\"required\":[\"applicationName\",\"applicationNamespace\",\"resourceUID\",\"resourceNamespace\",\"resourceName\"],\"properties\":{\"applicationName\":{\"type\":\"string\"},\"applicationNamespace\":{\"type\":\"string\"},\"resourceName\":{\"type\":\"string\"},\"resourceNamespace\":{\"type\":\"string\"},\"resourceUID\":{\"type\":\"string\"}},\"additionalProperties\":false}"),
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
		Result: tools.TypeSpec{
			Name:   "string",
			Schema: nil,
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
	},
	{
		Name:        "get_resources",
		Service:     "argocd",
		Toolset:     "argocd.argocd-mcp",
		Description: "get_resources return manifests for resources specified by resourceRefs. If resourceRefs is empty or not provided, fetches all resources managed by the application.",
		Payload: tools.TypeSpec{
			Name:   "*argocd.GetResourcesPayload",
			Schema: []byte("{\"type\":\"object\",\"required\":[\"applicationName\",\"applicationNamespace\"],\"properties\":{\"applicationName\":{\"type\":\"string\"},\"applicationNamespace\":{\"type\":\"string\"},\"resourceRefs\":{\"type\":\"array\",\"items\":{\"type\":\"object\",\"required\":[\"uid\",\"kind\",\"namespace\",\"name\",\"version\",\"group\"],\"properties\":{\"group\":{\"type\":\"string\"},\"kind\":{\"type\":\"string\"},\"name\":{\"type\":\"string\"},\"namespace\":{\"type\":\"string\"},\"uid\":{\"type\":\"string\"},\"version\":{\"type\":\"string\"}},\"additionalProperties\":false}}},\"additionalProperties\":false}"),
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
		Result: tools.TypeSpec{
			Name:   "string",
			Schema: nil,
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
	},
	{
		Name:        "get_resource_actions",
		Service:     "argocd",
		Toolset:     "argocd.argocd-mcp",
		Description: "get_resource_actions returns actions for a resource that is managed by an application",
		Payload: tools.TypeSpec{
			Name:   "*argocd.GetResourceActionsPayload",
			Schema: []byte("{\"type\":\"object\",\"required\":[\"applicationName\",\"applicationNamespace\",\"resourceRef\"],\"properties\":{\"applicationName\":{\"type\":\"string\"},\"applicationNamespace\":{\"type\":\"string\"},\"resourceRef\":{\"type\":\"object\",\"required\":[\"uid\",\"kind\",\"namespace\",\"name\",\"version\",\"group\"],\"properties\":{\"group\":{\"type\":\"string\"},\"kind\":{\"type\":\"string\"},\"name\":{\"type\":\"string\"},\"namespace\":{\"type\":\"string\"},\"uid\":{\"type\":\"string\"},\"version\":{\"type\":\"string\"}},\"additionalProperties\":false}},\"additionalProperties\":false}"),
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
		Result: tools.TypeSpec{
			Name:   "string",
			Schema: nil,
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
	},
	{
		Name:        "create_application",
		Service:     "argocd",
		Toolset:     "argocd.argocd-mcp",
		Description: "create_application creates a new ArgoCD application in the specified namespace. The application.metadata.namespace field determines where the Application resource will be created (e.g., \"argocd\", \"argocd-apps\", or any custom namespace).",
		Payload: tools.TypeSpec{
			Name:   "*argocd.CreateApplicationPayload",
			Schema: []byte("{\"type\":\"object\",\"required\":[\"application\"],\"properties\":{\"application\":{\"type\":\"object\",\"required\":[\"metadata\",\"spec\"],\"properties\":{\"metadata\":{\"type\":\"object\",\"required\":[\"name\",\"namespace\"],\"properties\":{\"name\":{\"type\":\"string\"},\"namespace\":{\"type\":\"string\"}},\"additionalProperties\":false},\"spec\":{\"type\":\"object\",\"required\":[\"project\",\"source\",\"syncPolicy\",\"destination\"],\"properties\":{\"destination\":{\"type\":\"object\",\"description\":\"The destination of the application.\\nOnly one of server or name must be specified.\",\"properties\":{\"name\":{\"type\":\"string\"},\"namespace\":{\"type\":\"string\"},\"server\":{\"type\":\"string\"}},\"additionalProperties\":false},\"project\":{\"type\":\"string\"},\"source\":{\"type\":\"object\",\"required\":[\"repoURL\",\"path\",\"targetRevision\"],\"properties\":{\"path\":{\"type\":\"string\"},\"repoURL\":{\"type\":\"string\"},\"targetRevision\":{\"type\":\"string\"}},\"additionalProperties\":false},\"syncPolicy\":{\"type\":\"object\",\"required\":[\"syncOptions\",\"retry\"],\"properties\":{\"automated\":{\"type\":\"object\",\"required\":[\"prune\",\"selfHeal\"],\"properties\":{\"prune\":{\"type\":\"boolean\"},\"selfHeal\":{\"type\":\"boolean\"}},\"additionalProperties\":false},\"retry\":{\"type\":\"object\",\"required\":[\"limit\",\"backoff\"],\"properties\":{\"backoff\":{\"type\":\"object\",\"required\":[\"duration\",\"maxDuration\",\"factor\"],\"properties\":{\"duration\":{\"type\":\"string\"},\"factor\":{\"type\":\"integer\"},\"maxDuration\":{\"type\":\"string\"}},\"additionalProperties\":false},\"limit\":{\"type\":\"integer\"}},\"additionalProperties\":false},\"syncOptions\":{\"type\":\"array\",\"items\":{\"type\":\"string\"}}},\"additionalProperties\":false}},\"additionalProperties\":false}},\"additionalProperties\":false}},\"additionalProperties\":false}"),
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
		Result: tools.TypeSpec{
			Name:   "string",
			Schema: nil,
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
	},
	{
		Name:        "update_application",
		Service:     "argocd",
		Toolset:     "argocd.argocd-mcp",
		Description: "update_application updates application",
		Payload: tools.TypeSpec{
			Name:   "*argocd.UpdateApplicationPayload",
			Schema: []byte("{\"type\":\"object\",\"required\":[\"applicationName\",\"application\"],\"properties\":{\"application\":{\"type\":\"object\",\"required\":[\"metadata\",\"spec\"],\"properties\":{\"metadata\":{\"type\":\"object\",\"required\":[\"name\",\"namespace\"],\"properties\":{\"name\":{\"type\":\"string\"},\"namespace\":{\"type\":\"string\"}},\"additionalProperties\":false},\"spec\":{\"type\":\"object\",\"required\":[\"project\",\"source\",\"syncPolicy\",\"destination\"],\"properties\":{\"destination\":{\"type\":\"object\",\"description\":\"The destination of the application.\\nOnly one of server or name must be specified.\",\"properties\":{\"name\":{\"type\":\"string\"},\"namespace\":{\"type\":\"string\"},\"server\":{\"type\":\"string\"}},\"additionalProperties\":false},\"project\":{\"type\":\"string\"},\"source\":{\"type\":\"object\",\"required\":[\"repoURL\",\"path\",\"targetRevision\"],\"properties\":{\"path\":{\"type\":\"string\"},\"repoURL\":{\"type\":\"string\"},\"targetRevision\":{\"type\":\"string\"}},\"additionalProperties\":false},\"syncPolicy\":{\"type\":\"object\",\"required\":[\"syncOptions\",\"retry\"],\"properties\":{\"automated\":{\"type\":\"object\",\"required\":[\"prune\",\"selfHeal\"],\"properties\":{\"prune\":{\"type\":\"boolean\"},\"selfHeal\":{\"type\":\"boolean\"}},\"additionalProperties\":false},\"retry\":{\"type\":\"object\",\"required\":[\"limit\",\"backoff\"],\"properties\":{\"backoff\":{\"type\":\"object\",\"required\":[\"duration\",\"maxDuration\",\"factor\"],\"properties\":{\"duration\":{\"type\":\"string\"},\"factor\":{\"type\":\"integer\"},\"maxDuration\":{\"type\":\"string\"}},\"additionalProperties\":false},\"limit\":{\"type\":\"integer\"}},\"additionalProperties\":false},\"syncOptions\":{\"type\":\"array\",\"items\":{\"type\":\"string\"}}},\"additionalProperties\":false}},\"additionalProperties\":false}},\"additionalProperties\":false},\"applicationName\":{\"type\":\"string\"}},\"additionalProperties\":false}"),
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
		Result: tools.TypeSpec{
			Name:   "string",
			Schema: nil,
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
	},
	{
		Name:        "delete_application",
		Service:     "argocd",
		Toolset:     "argocd.argocd-mcp",
		Description: "delete_application deletes application. Specify applicationNamespace if the application is in a non-default namespace to avoid permission errors.",
		Payload: tools.TypeSpec{
			Name:   "*argocd.DeleteApplicationPayload",
			Schema: []byte("{\"type\":\"object\",\"required\":[\"applicationName\"],\"properties\":{\"applicationName\":{\"type\":\"string\"},\"applicationNamespace\":{\"type\":\"string\",\"description\":\"The namespace where the application is located. Required if application is not in the default namespace.\"},\"cascade\":{\"type\":\"boolean\",\"description\":\"Whether to cascade the deletion to child resources\"},\"propagationPolicy\":{\"type\":\"string\",\"description\":\"Deletion propagation policy (e.g., \\\"Foreground\\\", \\\"Background\\\", \\\"Orphan\\\")\"}},\"additionalProperties\":false}"),
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
		Result: tools.TypeSpec{
			Name:   "string",
			Schema: nil,
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
	},
	{
		Name:        "sync_application",
		Service:     "argocd",
		Toolset:     "argocd.argocd-mcp",
		Description: "sync_application syncs application. Specify applicationNamespace if the application is in a non-default namespace to avoid permission errors.",
		Payload: tools.TypeSpec{
			Name:   "*argocd.SyncApplicationPayload",
			Schema: []byte("{\"type\":\"object\",\"required\":[\"applicationName\"],\"properties\":{\"applicationName\":{\"type\":\"string\"},\"applicationNamespace\":{\"type\":\"string\",\"description\":\"The namespace where the application is located. Required if application is not in the default namespace.\"},\"dryRun\":{\"type\":\"boolean\",\"description\":\"Perform a dry run sync without applying changes\"},\"prune\":{\"type\":\"boolean\",\"description\":\"Remove resources that are no longer defined in the source\"},\"revision\":{\"type\":\"string\",\"description\":\"Sync to a specific revision instead of the latest\"},\"syncOptions\":{\"type\":\"array\",\"description\":\"Additional sync options (e.g., [\\\"CreateNamespace=true\\\", \\\"PrunePropagationPolicy=foreground\\\"])\",\"items\":{\"type\":\"string\"}}},\"additionalProperties\":false}"),
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
		Result: tools.TypeSpec{
			Name:   "string",
			Schema: nil,
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
	},
	{
		Name:        "run_resource_action",
		Service:     "argocd",
		Toolset:     "argocd.argocd-mcp",
		Description: "run_resource_action runs an action on a resource",
		Payload: tools.TypeSpec{
			Name:   "*argocd.RunResourceActionPayload",
			Schema: []byte("{\"type\":\"object\",\"required\":[\"applicationName\",\"applicationNamespace\",\"resourceRef\",\"action\"],\"properties\":{\"action\":{\"type\":\"string\"},\"applicationName\":{\"type\":\"string\"},\"applicationNamespace\":{\"type\":\"string\"},\"resourceRef\":{\"type\":\"object\",\"required\":[\"uid\",\"kind\",\"namespace\",\"name\",\"version\",\"group\"],\"properties\":{\"group\":{\"type\":\"string\"},\"kind\":{\"type\":\"string\"},\"name\":{\"type\":\"string\"},\"namespace\":{\"type\":\"string\"},\"uid\":{\"type\":\"string\"},\"version\":{\"type\":\"string\"}},\"additionalProperties\":false}},\"additionalProperties\":false}"),
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
		Result: tools.TypeSpec{
			Name:   "string",
			Schema: nil,
			Codec: tools.JSONCodec[any]{
				ToJSON: func(v any) ([]byte, error) {
					return json.Marshal(v)
				},
				FromJSON: func(data []byte) (any, error) {
					if len(data) == 0 {
						return nil, nil
					}
					var out any
					if err := json.Unmarshal(data, &out); err != nil {
						return nil, err
					}
					return out, nil
				},
			},
		},
	},
}

// RegisterArgocdArgocdMcpToolset registers the argocd-mcp toolset with the runtime.
// The caller parameter provides the MCP client for making remote calls.
func RegisterArgocdArgocdMcpToolset(ctx context.Context, rt *agentsruntime.Runtime, caller mcpruntime.Caller) error {
	if rt == nil {
		return errors.New("runtime is required")
	}
	if caller == nil {
		return errors.New("mcp caller is required")
	}

	exec := func(ctx context.Context, call planner.ToolRequest) (planner.ToolResult, error) {
		fullName := call.Name
		toolName := string(fullName)
		const suitePrefix = "argocd.argocd-mcp" + "."
		if strings.HasPrefix(toolName, suitePrefix) {
			toolName = toolName[len(suitePrefix):]
		}

		payload, err := json.Marshal(call.Payload)
		if err != nil {
			return planner.ToolResult{Name: fullName}, err
		}

		resp, err := caller.CallTool(ctx, mcpruntime.CallRequest{
			Suite:   "argocd.argocd-mcp",
			Tool:    toolName,
			Payload: payload,
		})
		if err != nil {
			return ArgocdArgocdMcpToolsetHandleError(fullName, err), nil
		}

		var value any
		if len(resp.Result) > 0 {
			if err := json.Unmarshal(resp.Result, &value); err != nil {
				return planner.ToolResult{Name: fullName}, err
			}
		}

		var toolTelemetry *telemetry.ToolTelemetry
		if len(resp.Structured) > 0 {
			var structured any
			if err := json.Unmarshal(resp.Structured, &structured); err != nil {
				return planner.ToolResult{Name: fullName}, err
			}
			toolTelemetry = &telemetry.ToolTelemetry{
				Extra: map[string]any{"structured": structured},
			}
		}

		return planner.ToolResult{
			Name:      fullName,
			Result:    value,
			Telemetry: toolTelemetry,
		}, nil
	}

	return rt.RegisterToolset(agentsruntime.ToolsetRegistration{
		Name:        "argocd.argocd-mcp",
		Description: "Argo CD MCP tool service.",
		Execute: func(ctx context.Context, call *planner.ToolRequest) (*planner.ToolResult, error) {
			if call == nil {
				return nil, errors.New("tool request is nil")
			}
			out, err := exec(ctx, *call)
			if err != nil {
				return nil, err
			}
			return &out, nil
		},
		Specs:            ArgocdArgocdMcpToolsetToolSpecs,
		DecodeInExecutor: true,
	})
}

// ArgocdArgocdMcpToolsetHandleError converts an error into a tool result with appropriate retry hints.
func ArgocdArgocdMcpToolsetHandleError(toolName tools.Ident, err error) planner.ToolResult {
	result := planner.ToolResult{
		Name:  toolName,
		Error: planner.ToolErrorFromError(err),
	}
	if hint := ArgocdArgocdMcpToolsetRetryHint(toolName, err); hint != nil {
		result.RetryHint = hint
	}
	return result
}

// ArgocdArgocdMcpToolsetRetryHint determines if an error should trigger a retry and returns appropriate hints.
func ArgocdArgocdMcpToolsetRetryHint(toolName tools.Ident, err error) *planner.RetryHint {
	key := string(toolName)
	var retryErr *retry.RetryableError
	if errors.As(err, &retryErr) {
		return &planner.RetryHint{
			Reason:         planner.RetryReasonInvalidArguments,
			Tool:           toolName,
			Message:        retryErr.Prompt,
			RestrictToTool: true,
		}
	}
	var rpcErr *jsonrpc.Error
	if errors.As(err, &rpcErr) {
		switch rpcErr.Code {
		case jsonrpc.CodeInvalidParams:
			// Schema and example are known at generation time - use switch for direct lookup
			var schemaJSON, example string
			switch key {
			case "list_applications":
				schemaJSON = "{\"type\":\"object\",\"properties\":{\"limit\":{\"type\":\"integer\",\"description\":\"Maximum number of applications to return. Use this to reduce token usage when there are many applications. Optional.\"},\"offset\":{\"type\":\"integer\",\"description\":\"Number of applications to skip before returning results. Use with limit for pagination. Optional.\"},\"search\":{\"type\":\"string\",\"description\":\"Search applications by name. This is a partial match on the application name and does not support glob patterns (e.g. \\\"*\\\"). Optional.\"}},\"additionalProperties\":false}"
				example = "{\"limit\":1,\"offset\":1,\"search\":\"abc123\"}"
			case "get_application":
				schemaJSON = "{\"type\":\"object\",\"required\":[\"applicationName\"],\"properties\":{\"applicationName\":{\"type\":\"string\"},\"applicationNamespace\":{\"type\":\"string\"}},\"additionalProperties\":false}"
				example = "{\"applicationName\":\"abc123\",\"applicationNamespace\":\"abc123\"}"
			case "get_application_resource_tree":
				schemaJSON = "{\"type\":\"object\",\"required\":[\"applicationName\"],\"properties\":{\"applicationName\":{\"type\":\"string\"}},\"additionalProperties\":false}"
				example = "{\"applicationName\":\"abc123\"}"
			case "get_application_managed_resources":
				schemaJSON = "{\"type\":\"object\",\"required\":[\"applicationName\"],\"properties\":{\"appNamespace\":{\"type\":\"string\",\"description\":\"Filter by Argo CD application namespace\"},\"applicationName\":{\"type\":\"string\"},\"group\":{\"type\":\"string\",\"description\":\"Filter by API group\"},\"kind\":{\"type\":\"string\",\"description\":\"Filter by Kubernetes resource kind (e.g., \\\"ConfigMap\\\", \\\"Secret\\\", \\\"Deployment\\\")\"},\"name\":{\"type\":\"string\",\"description\":\"Filter by resource name\"},\"namespace\":{\"type\":\"string\",\"description\":\"Filter by Kubernetes namespace\"},\"project\":{\"type\":\"string\",\"description\":\"Filter by Argo CD project\"},\"version\":{\"type\":\"string\",\"description\":\"Filter by resource API version\"}},\"additionalProperties\":false}"
				example = "{\"appNamespace\":\"abc123\",\"applicationName\":\"abc123\",\"group\":\"abc123\",\"kind\":\"abc123\",\"name\":\"abc123\",\"namespace\":\"abc123\",\"project\":\"abc123\",\"version\":\"abc123\"}"
			case "get_application_workload_logs":
				schemaJSON = "{\"type\":\"object\",\"required\":[\"applicationName\",\"applicationNamespace\",\"resourceRef\",\"container\"],\"properties\":{\"applicationName\":{\"type\":\"string\"},\"applicationNamespace\":{\"type\":\"string\"},\"container\":{\"type\":\"string\"},\"resourceRef\":{\"type\":\"object\",\"required\":[\"uid\",\"kind\",\"namespace\",\"name\",\"version\",\"group\"],\"properties\":{\"group\":{\"type\":\"string\"},\"kind\":{\"type\":\"string\"},\"name\":{\"type\":\"string\"},\"namespace\":{\"type\":\"string\"},\"uid\":{\"type\":\"string\"},\"version\":{\"type\":\"string\"}},\"additionalProperties\":false}},\"additionalProperties\":false}"
				example = "{\"applicationName\":\"abc123\",\"applicationNamespace\":\"abc123\",\"container\":\"abc123\",\"resourceRef\":{\"group\":\"abc123\",\"kind\":\"abc123\",\"name\":\"abc123\",\"namespace\":\"abc123\",\"uid\":\"abc123\",\"version\":\"abc123\"}}"
			case "get_application_events":
				schemaJSON = "{\"type\":\"object\",\"required\":[\"applicationName\"],\"properties\":{\"applicationName\":{\"type\":\"string\"}},\"additionalProperties\":false}"
				example = "{\"applicationName\":\"abc123\"}"
			case "get_resource_events":
				schemaJSON = "{\"type\":\"object\",\"required\":[\"applicationName\",\"applicationNamespace\",\"resourceUID\",\"resourceNamespace\",\"resourceName\"],\"properties\":{\"applicationName\":{\"type\":\"string\"},\"applicationNamespace\":{\"type\":\"string\"},\"resourceName\":{\"type\":\"string\"},\"resourceNamespace\":{\"type\":\"string\"},\"resourceUID\":{\"type\":\"string\"}},\"additionalProperties\":false}"
				example = "{\"applicationName\":\"abc123\",\"applicationNamespace\":\"abc123\",\"resourceName\":\"abc123\",\"resourceNamespace\":\"abc123\",\"resourceUID\":\"abc123\"}"
			case "get_resources":
				schemaJSON = "{\"type\":\"object\",\"required\":[\"applicationName\",\"applicationNamespace\"],\"properties\":{\"applicationName\":{\"type\":\"string\"},\"applicationNamespace\":{\"type\":\"string\"},\"resourceRefs\":{\"type\":\"array\",\"items\":{\"type\":\"object\",\"required\":[\"uid\",\"kind\",\"namespace\",\"name\",\"version\",\"group\"],\"properties\":{\"group\":{\"type\":\"string\"},\"kind\":{\"type\":\"string\"},\"name\":{\"type\":\"string\"},\"namespace\":{\"type\":\"string\"},\"uid\":{\"type\":\"string\"},\"version\":{\"type\":\"string\"}},\"additionalProperties\":false}}},\"additionalProperties\":false}"
				example = "{\"applicationName\":\"abc123\",\"applicationNamespace\":\"abc123\",\"resourceRefs\":[{\"group\":\"abc123\",\"kind\":\"abc123\",\"name\":\"abc123\",\"namespace\":\"abc123\",\"uid\":\"abc123\",\"version\":\"abc123\"}]}"
			case "get_resource_actions":
				schemaJSON = "{\"type\":\"object\",\"required\":[\"applicationName\",\"applicationNamespace\",\"resourceRef\"],\"properties\":{\"applicationName\":{\"type\":\"string\"},\"applicationNamespace\":{\"type\":\"string\"},\"resourceRef\":{\"type\":\"object\",\"required\":[\"uid\",\"kind\",\"namespace\",\"name\",\"version\",\"group\"],\"properties\":{\"group\":{\"type\":\"string\"},\"kind\":{\"type\":\"string\"},\"name\":{\"type\":\"string\"},\"namespace\":{\"type\":\"string\"},\"uid\":{\"type\":\"string\"},\"version\":{\"type\":\"string\"}},\"additionalProperties\":false}},\"additionalProperties\":false}"
				example = "{\"applicationName\":\"abc123\",\"applicationNamespace\":\"abc123\",\"resourceRef\":{\"group\":\"abc123\",\"kind\":\"abc123\",\"name\":\"abc123\",\"namespace\":\"abc123\",\"uid\":\"abc123\",\"version\":\"abc123\"}}"
			case "create_application":
				schemaJSON = "{\"type\":\"object\",\"required\":[\"application\"],\"properties\":{\"application\":{\"type\":\"object\",\"required\":[\"metadata\",\"spec\"],\"properties\":{\"metadata\":{\"type\":\"object\",\"required\":[\"name\",\"namespace\"],\"properties\":{\"name\":{\"type\":\"string\"},\"namespace\":{\"type\":\"string\"}},\"additionalProperties\":false},\"spec\":{\"type\":\"object\",\"required\":[\"project\",\"source\",\"syncPolicy\",\"destination\"],\"properties\":{\"destination\":{\"type\":\"object\",\"description\":\"The destination of the application.\\nOnly one of server or name must be specified.\",\"properties\":{\"name\":{\"type\":\"string\"},\"namespace\":{\"type\":\"string\"},\"server\":{\"type\":\"string\"}},\"additionalProperties\":false},\"project\":{\"type\":\"string\"},\"source\":{\"type\":\"object\",\"required\":[\"repoURL\",\"path\",\"targetRevision\"],\"properties\":{\"path\":{\"type\":\"string\"},\"repoURL\":{\"type\":\"string\"},\"targetRevision\":{\"type\":\"string\"}},\"additionalProperties\":false},\"syncPolicy\":{\"type\":\"object\",\"required\":[\"syncOptions\",\"retry\"],\"properties\":{\"automated\":{\"type\":\"object\",\"required\":[\"prune\",\"selfHeal\"],\"properties\":{\"prune\":{\"type\":\"boolean\"},\"selfHeal\":{\"type\":\"boolean\"}},\"additionalProperties\":false},\"retry\":{\"type\":\"object\",\"required\":[\"limit\",\"backoff\"],\"properties\":{\"backoff\":{\"type\":\"object\",\"required\":[\"duration\",\"maxDuration\",\"factor\"],\"properties\":{\"duration\":{\"type\":\"string\"},\"factor\":{\"type\":\"integer\"},\"maxDuration\":{\"type\":\"string\"}},\"additionalProperties\":false},\"limit\":{\"type\":\"integer\"}},\"additionalProperties\":false},\"syncOptions\":{\"type\":\"array\",\"items\":{\"type\":\"string\"}}},\"additionalProperties\":false}},\"additionalProperties\":false}},\"additionalProperties\":false}},\"additionalProperties\":false}"
				example = "{\"application\":{\"metadata\":{\"name\":\"abc123\",\"namespace\":\"abc123\"},\"spec\":{\"destination\":{\"name\":\"abc123\",\"namespace\":\"abc123\",\"server\":\"abc123\"},\"project\":\"abc123\",\"source\":{\"path\":\"abc123\",\"repoURL\":\"abc123\",\"targetRevision\":\"abc123\"},\"syncPolicy\":{\"automated\":{\"prune\":false,\"selfHeal\":false},\"retry\":{\"backoff\":{\"duration\":\"abc123\",\"factor\":1,\"maxDuration\":\"abc123\"},\"limit\":1},\"syncOptions\":[\"abc123\"]}}}}"
			case "update_application":
				schemaJSON = "{\"type\":\"object\",\"required\":[\"applicationName\",\"application\"],\"properties\":{\"application\":{\"type\":\"object\",\"required\":[\"metadata\",\"spec\"],\"properties\":{\"metadata\":{\"type\":\"object\",\"required\":[\"name\",\"namespace\"],\"properties\":{\"name\":{\"type\":\"string\"},\"namespace\":{\"type\":\"string\"}},\"additionalProperties\":false},\"spec\":{\"type\":\"object\",\"required\":[\"project\",\"source\",\"syncPolicy\",\"destination\"],\"properties\":{\"destination\":{\"type\":\"object\",\"description\":\"The destination of the application.\\nOnly one of server or name must be specified.\",\"properties\":{\"name\":{\"type\":\"string\"},\"namespace\":{\"type\":\"string\"},\"server\":{\"type\":\"string\"}},\"additionalProperties\":false},\"project\":{\"type\":\"string\"},\"source\":{\"type\":\"object\",\"required\":[\"repoURL\",\"path\",\"targetRevision\"],\"properties\":{\"path\":{\"type\":\"string\"},\"repoURL\":{\"type\":\"string\"},\"targetRevision\":{\"type\":\"string\"}},\"additionalProperties\":false},\"syncPolicy\":{\"type\":\"object\",\"required\":[\"syncOptions\",\"retry\"],\"properties\":{\"automated\":{\"type\":\"object\",\"required\":[\"prune\",\"selfHeal\"],\"properties\":{\"prune\":{\"type\":\"boolean\"},\"selfHeal\":{\"type\":\"boolean\"}},\"additionalProperties\":false},\"retry\":{\"type\":\"object\",\"required\":[\"limit\",\"backoff\"],\"properties\":{\"backoff\":{\"type\":\"object\",\"required\":[\"duration\",\"maxDuration\",\"factor\"],\"properties\":{\"duration\":{\"type\":\"string\"},\"factor\":{\"type\":\"integer\"},\"maxDuration\":{\"type\":\"string\"}},\"additionalProperties\":false},\"limit\":{\"type\":\"integer\"}},\"additionalProperties\":false},\"syncOptions\":{\"type\":\"array\",\"items\":{\"type\":\"string\"}}},\"additionalProperties\":false}},\"additionalProperties\":false}},\"additionalProperties\":false},\"applicationName\":{\"type\":\"string\"}},\"additionalProperties\":false}"
				example = "{\"application\":{\"metadata\":{\"name\":\"abc123\",\"namespace\":\"abc123\"},\"spec\":{\"destination\":{\"name\":\"abc123\",\"namespace\":\"abc123\",\"server\":\"abc123\"},\"project\":\"abc123\",\"source\":{\"path\":\"abc123\",\"repoURL\":\"abc123\",\"targetRevision\":\"abc123\"},\"syncPolicy\":{\"automated\":{\"prune\":false,\"selfHeal\":false},\"retry\":{\"backoff\":{\"duration\":\"abc123\",\"factor\":1,\"maxDuration\":\"abc123\"},\"limit\":1},\"syncOptions\":[\"abc123\"]}}},\"applicationName\":\"abc123\"}"
			case "delete_application":
				schemaJSON = "{\"type\":\"object\",\"required\":[\"applicationName\"],\"properties\":{\"applicationName\":{\"type\":\"string\"},\"applicationNamespace\":{\"type\":\"string\",\"description\":\"The namespace where the application is located. Required if application is not in the default namespace.\"},\"cascade\":{\"type\":\"boolean\",\"description\":\"Whether to cascade the deletion to child resources\"},\"propagationPolicy\":{\"type\":\"string\",\"description\":\"Deletion propagation policy (e.g., \\\"Foreground\\\", \\\"Background\\\", \\\"Orphan\\\")\"}},\"additionalProperties\":false}"
				example = "{\"applicationName\":\"abc123\",\"applicationNamespace\":\"abc123\",\"cascade\":false,\"propagationPolicy\":\"abc123\"}"
			case "sync_application":
				schemaJSON = "{\"type\":\"object\",\"required\":[\"applicationName\"],\"properties\":{\"applicationName\":{\"type\":\"string\"},\"applicationNamespace\":{\"type\":\"string\",\"description\":\"The namespace where the application is located. Required if application is not in the default namespace.\"},\"dryRun\":{\"type\":\"boolean\",\"description\":\"Perform a dry run sync without applying changes\"},\"prune\":{\"type\":\"boolean\",\"description\":\"Remove resources that are no longer defined in the source\"},\"revision\":{\"type\":\"string\",\"description\":\"Sync to a specific revision instead of the latest\"},\"syncOptions\":{\"type\":\"array\",\"description\":\"Additional sync options (e.g., [\\\"CreateNamespace=true\\\", \\\"PrunePropagationPolicy=foreground\\\"])\",\"items\":{\"type\":\"string\"}}},\"additionalProperties\":false}"
				example = "{\"applicationName\":\"abc123\",\"applicationNamespace\":\"abc123\",\"dryRun\":false,\"prune\":false,\"revision\":\"abc123\",\"syncOptions\":[\"abc123\"]}"
			case "run_resource_action":
				schemaJSON = "{\"type\":\"object\",\"required\":[\"applicationName\",\"applicationNamespace\",\"resourceRef\",\"action\"],\"properties\":{\"action\":{\"type\":\"string\"},\"applicationName\":{\"type\":\"string\"},\"applicationNamespace\":{\"type\":\"string\"},\"resourceRef\":{\"type\":\"object\",\"required\":[\"uid\",\"kind\",\"namespace\",\"name\",\"version\",\"group\"],\"properties\":{\"group\":{\"type\":\"string\"},\"kind\":{\"type\":\"string\"},\"name\":{\"type\":\"string\"},\"namespace\":{\"type\":\"string\"},\"uid\":{\"type\":\"string\"},\"version\":{\"type\":\"string\"}},\"additionalProperties\":false}},\"additionalProperties\":false}"
				example = "{\"action\":\"abc123\",\"applicationName\":\"abc123\",\"applicationNamespace\":\"abc123\",\"resourceRef\":{\"group\":\"abc123\",\"kind\":\"abc123\",\"name\":\"abc123\",\"namespace\":\"abc123\",\"uid\":\"abc123\",\"version\":\"abc123\"}}"
			}
			prompt := retry.BuildRepairPrompt("tools/call:"+key, rpcErr.Message, example, schemaJSON)
			return &planner.RetryHint{
				Reason:         planner.RetryReasonInvalidArguments,
				Tool:           toolName,
				Message:        prompt,
				RestrictToTool: true,
			}
		case jsonrpc.CodeMethodNotFound:
			return &planner.RetryHint{
				Reason:  planner.RetryReasonToolUnavailable,
				Tool:    toolName,
				Message: rpcErr.Message,
			}
		}
	}
	return nil
}
