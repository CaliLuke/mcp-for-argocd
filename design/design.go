package design

import (
	. "github.com/CaliLuke/loom-mcp/dsl"
	. "github.com/CaliLuke/loom/dsl"
)

var _ = API("argocdmcp", func() {
	Title("Argo CD MCP Server")
	Description("A Go port of the Argo CD MCP server.")
})

var ApplicationNamespace = Type("ApplicationNamespace", func() {
	Attribute("value", String, `The namespace where the ArgoCD application resource will be created.
This is the namespace of the Application resource itself, not the destination namespace for the application's resources.
You can specify any valid Kubernetes namespace (e.g., 'argocd', 'argocd-apps', 'my-namespace', etc.).
The default ArgoCD namespace is typically 'argocd', but you can use any namespace you prefer.`)
})

var ResourceRef = Type("ResourceRef", func() {
	Attribute("uid", String)
	Attribute("kind", String)
	Attribute("namespace", String)
	Attribute("name", String)
	Attribute("version", String)
	Attribute("group", String)
	Required("uid", "kind", "namespace", "name", "version", "group")
})

var ApplicationMetadata = Type("ApplicationMetadata", func() {
	Attribute("name", String)
	Attribute("namespace", String)
	Required("name", "namespace")
})

var ApplicationSource = Type("ApplicationSource", func() {
	Attribute("repoURL", String)
	Attribute("path", String)
	Attribute("targetRevision", String)
	Required("repoURL", "path", "targetRevision")
})

var ApplicationAutomated = Type("ApplicationAutomated", func() {
	Attribute("prune", Boolean)
	Attribute("selfHeal", Boolean)
	Required("prune", "selfHeal")
})

var ApplicationBackoff = Type("ApplicationBackoff", func() {
	Attribute("duration", String)
	Attribute("maxDuration", String)
	Attribute("factor", Int)
	Required("duration", "maxDuration", "factor")
})

var ApplicationRetry = Type("ApplicationRetry", func() {
	Attribute("limit", Int)
	Attribute("backoff", ApplicationBackoff)
	Required("limit", "backoff")
})

var ApplicationSyncPolicy = Type("ApplicationSyncPolicy", func() {
	Attribute("syncOptions", ArrayOf(String))
	Attribute("automated", ApplicationAutomated)
	Attribute("retry", ApplicationRetry)
	Required("syncOptions", "retry")
})

var ApplicationDestination = Type("ApplicationDestination", func() {
	Attribute("server", String)
	Attribute("namespace", String)
	Attribute("name", String)
	Description(`The destination of the application.
Only one of server or name must be specified.`)
})

var ApplicationSpec = Type("ApplicationSpec", func() {
	Attribute("project", String)
	Attribute("source", ApplicationSource)
	Attribute("syncPolicy", ApplicationSyncPolicy)
	Attribute("destination", ApplicationDestination)
	Required("project", "source", "syncPolicy", "destination")
})

var ListApplicationsPayload = Type("ListApplicationsPayload", func() {
	Attribute("search", String, `Search applications by name. This is a partial match on the application name and does not support glob patterns (e.g. "*"). Optional.`)
	Attribute("limit", Int, "Maximum number of applications to return. Use this to reduce token usage when there are many applications. Optional.")
	Attribute("offset", Int, "Number of applications to skip before returning results. Use with limit for pagination. Optional.")
})

var GetApplicationPayload = Type("GetApplicationPayload", func() {
	Attribute("applicationName", String)
	Attribute("applicationNamespace", String)
	Required("applicationName")
})

var GetApplicationResourceTreePayload = Type("GetApplicationResourceTreePayload", func() {
	Attribute("applicationName", String)
	Required("applicationName")
})

var GetApplicationManagedResourcesPayload = Type("GetApplicationManagedResourcesPayload", func() {
	Attribute("applicationName", String)
	Attribute("kind", String, `Filter by Kubernetes resource kind (e.g., "ConfigMap", "Secret", "Deployment")`)
	Attribute("namespace", String, "Filter by Kubernetes namespace")
	Attribute("name", String, "Filter by resource name")
	Attribute("version", String, "Filter by resource API version")
	Attribute("group", String, "Filter by API group")
	Attribute("appNamespace", String, "Filter by Argo CD application namespace")
	Attribute("project", String, "Filter by Argo CD project")
	Required("applicationName")
})

var GetApplicationWorkloadLogsPayload = Type("GetApplicationWorkloadLogsPayload", func() {
	Attribute("applicationName", String)
	Attribute("applicationNamespace", String)
	Attribute("resourceRef", ResourceRef)
	Attribute("container", String)
	Required("applicationName", "applicationNamespace", "resourceRef", "container")
})

var GetApplicationEventsPayload = Type("GetApplicationEventsPayload", func() {
	Attribute("applicationName", String)
	Required("applicationName")
})

var GetResourceEventsPayload = Type("GetResourceEventsPayload", func() {
	Attribute("applicationName", String)
	Attribute("applicationNamespace", String)
	Attribute("resourceUID", String)
	Attribute("resourceNamespace", String)
	Attribute("resourceName", String)
	Required("applicationName", "applicationNamespace", "resourceUID", "resourceNamespace", "resourceName")
})

var GetResourcesPayload = Type("GetResourcesPayload", func() {
	Attribute("applicationName", String)
	Attribute("applicationNamespace", String)
	Attribute("resourceRefs", ArrayOf(ResourceRef))
	Required("applicationName", "applicationNamespace")
})

var GetResourceActionsPayload = Type("GetResourceActionsPayload", func() {
	Attribute("applicationName", String)
	Attribute("applicationNamespace", String)
	Attribute("resourceRef", ResourceRef)
	Required("applicationName", "applicationNamespace", "resourceRef")
})

var Application = Type("Application", func() {
	Attribute("metadata", ApplicationMetadata)
	Attribute("spec", ApplicationSpec)
	Required("metadata", "spec")
})

var CreateApplicationPayload = Type("CreateApplicationPayload", func() {
	Attribute("application", Application)
	Required("application")
})

var UpdateApplicationPayload = Type("UpdateApplicationPayload", func() {
	Attribute("applicationName", String)
	Attribute("application", Application)
	Required("applicationName", "application")
})

var DeleteApplicationPayload = Type("DeleteApplicationPayload", func() {
	Attribute("applicationName", String)
	Attribute("applicationNamespace", String, "The namespace where the application is located. Required if application is not in the default namespace.")
	Attribute("cascade", Boolean, "Whether to cascade the deletion to child resources")
	Attribute("propagationPolicy", String, `Deletion propagation policy (e.g., "Foreground", "Background", "Orphan")`)
	Required("applicationName")
})

var SyncApplicationPayload = Type("SyncApplicationPayload", func() {
	Attribute("applicationName", String)
	Attribute("applicationNamespace", String, "The namespace where the application is located. Required if application is not in the default namespace.")
	Attribute("dryRun", Boolean, "Perform a dry run sync without applying changes")
	Attribute("prune", Boolean, "Remove resources that are no longer defined in the source")
	Attribute("revision", String, "Sync to a specific revision instead of the latest")
	Attribute("syncOptions", ArrayOf(String), `Additional sync options (e.g., ["CreateNamespace=true", "PrunePropagationPolicy=foreground"])`)
	Required("applicationName")
})

var RunResourceActionPayload = Type("RunResourceActionPayload", func() {
	Attribute("applicationName", String)
	Attribute("applicationNamespace", String)
	Attribute("resourceRef", ResourceRef)
	Attribute("action", String)
	Required("applicationName", "applicationNamespace", "resourceRef", "action")
})

var _ = Service("argocd", func() {
	Description("Argo CD MCP tool service.")
	MCP("argocd-mcp", "1.0.0", ProtocolVersion("2025-06-18"))
	JSONRPC(func() {
		POST("/mcp")
	})

	Method("list_applications", func() {
		Description("list_applications returns list of applications")
		Payload(ListApplicationsPayload)
		Result(String)
		Tool("list_applications", "list_applications returns list of applications")
		JSONRPC(func() {})
	})

	Method("get_application", func() {
		Description("get_application returns application by application name. Optionally specify the application namespace to get applications from non-default namespaces.")
		Payload(GetApplicationPayload)
		Result(String)
		Tool("get_application", "get_application returns application by application name. Optionally specify the application namespace to get applications from non-default namespaces.")
		JSONRPC(func() {})
	})

	Method("get_application_resource_tree", func() {
		Description("get_application_resource_tree returns resource tree for application by application name")
		Payload(GetApplicationResourceTreePayload)
		Result(String)
		Tool("get_application_resource_tree", "get_application_resource_tree returns resource tree for application by application name")
		JSONRPC(func() {})
	})

	Method("get_application_managed_resources", func() {
		Description(`get_application_managed_resources returns managed resources for application by application name with optional filtering. Use filters to avoid token limits with large applications. Examples: kind="ConfigMap" for config maps only, namespace="production" for specific namespace, or combine multiple filters.`)
		Payload(GetApplicationManagedResourcesPayload)
		Result(String)
		Tool("get_application_managed_resources", `get_application_managed_resources returns managed resources for application by application name with optional filtering. Use filters to avoid token limits with large applications. Examples: kind="ConfigMap" for config maps only, namespace="production" for specific namespace, or combine multiple filters.`)
		JSONRPC(func() {})
	})

	Method("get_application_workload_logs", func() {
		Description("get_application_workload_logs returns logs for application workload (Deployment, StatefulSet, Pod, etc.) by application name and resource ref and optionally container name")
		Payload(GetApplicationWorkloadLogsPayload)
		Result(String)
		Tool("get_application_workload_logs", "get_application_workload_logs returns logs for application workload (Deployment, StatefulSet, Pod, etc.) by application name and resource ref and optionally container name")
		JSONRPC(func() {})
	})

	Method("get_application_events", func() {
		Description("get_application_events returns events for application by application name")
		Payload(GetApplicationEventsPayload)
		Result(String)
		Tool("get_application_events", "get_application_events returns events for application by application name")
		JSONRPC(func() {})
	})

	Method("get_resource_events", func() {
		Description("get_resource_events returns events for a resource that is managed by an application")
		Payload(GetResourceEventsPayload)
		Result(String)
		Tool("get_resource_events", "get_resource_events returns events for a resource that is managed by an application")
		JSONRPC(func() {})
	})

	Method("get_resources", func() {
		Description("get_resources return manifests for resources specified by resourceRefs. If resourceRefs is empty or not provided, fetches all resources managed by the application.")
		Payload(GetResourcesPayload)
		Result(String)
		Tool("get_resources", "get_resources return manifests for resources specified by resourceRefs. If resourceRefs is empty or not provided, fetches all resources managed by the application.")
		JSONRPC(func() {})
	})

	Method("get_resource_actions", func() {
		Description("get_resource_actions returns actions for a resource that is managed by an application")
		Payload(GetResourceActionsPayload)
		Result(String)
		Tool("get_resource_actions", "get_resource_actions returns actions for a resource that is managed by an application")
		JSONRPC(func() {})
	})

	Method("create_application", func() {
		Description("create_application creates a new ArgoCD application in the specified namespace. The application.metadata.namespace field determines where the Application resource will be created (e.g., \"argocd\", \"argocd-apps\", or any custom namespace).")
		Payload(CreateApplicationPayload)
		Result(String)
		Tool("create_application", "create_application creates a new ArgoCD application in the specified namespace. The application.metadata.namespace field determines where the Application resource will be created (e.g., \"argocd\", \"argocd-apps\", or any custom namespace).")
		JSONRPC(func() {})
	})

	Method("update_application", func() {
		Description("update_application updates application")
		Payload(UpdateApplicationPayload)
		Result(String)
		Tool("update_application", "update_application updates application")
		JSONRPC(func() {})
	})

	Method("delete_application", func() {
		Description("delete_application deletes application. Specify applicationNamespace if the application is in a non-default namespace to avoid permission errors.")
		Payload(DeleteApplicationPayload)
		Result(String)
		Tool("delete_application", "delete_application deletes application. Specify applicationNamespace if the application is in a non-default namespace to avoid permission errors.")
		JSONRPC(func() {})
	})

	Method("sync_application", func() {
		Description("sync_application syncs application. Specify applicationNamespace if the application is in a non-default namespace to avoid permission errors.")
		Payload(SyncApplicationPayload)
		Result(String)
		Tool("sync_application", "sync_application syncs application. Specify applicationNamespace if the application is in a non-default namespace to avoid permission errors.")
		JSONRPC(func() {})
	})

	Method("run_resource_action", func() {
		Description("run_resource_action runs an action on a resource")
		Payload(RunResourceActionPayload)
		Result(String)
		Tool("run_resource_action", "run_resource_action runs an action on a resource")
		JSONRPC(func() {})
	})
})
