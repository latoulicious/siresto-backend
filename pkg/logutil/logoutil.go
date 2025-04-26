package logutil

// Route returns a consistent structured format for route logging
func Route(method, path string) map[string]interface{} {
	return map[string]interface{}{
		"method": method,
		"path":   path,
		"source": "router",
		"action": "route.registered",
		"entity": "route",
	}
}

// MainCall standardizes logs from main package
func MainCall(action, entity string, extra map[string]interface{}) map[string]interface{} {
	fields := map[string]interface{}{
		"source": "main",
		"action": action,
		"entity": entity,
	}
	for k, v := range extra {
		fields[k] = v
	}
	return fields
}

// HandlerCall standardizes log entries from HTTP handlers (controllers)
func HandlerCall(action, entity string, extra map[string]interface{}) map[string]interface{} {
	fields := map[string]interface{}{
		"source": "handler",
		"action": action,
		"entity": entity,
	}
	for k, v := range extra {
		fields[k] = v
	}
	return fields
}

// ServiceCall standardizes logs from services
func ServiceCall(action, entity string, extra map[string]interface{}) map[string]interface{} {
	fields := map[string]interface{}{
		"source": "service",
		"action": action,
		"entity": entity,
	}
	for k, v := range extra {
		fields[k] = v
	}
	return fields
}

// RepoCall standardizes logs from repositories
func RepoCall(action, entity string, extra map[string]interface{}) map[string]interface{} {
	fields := map[string]interface{}{
		"source": "repository",
		"action": action,
		"entity": entity,
	}
	for k, v := range extra {
		fields[k] = v
	}
	return fields
}
