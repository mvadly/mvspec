package util

import (
	"encoding/json"
	"strings"
)

var RefPrefixes = []string{"response.", "models.", "util."}

func CleanRef(ref string) string {
	for _, prefix := range RefPrefixes {
		ref = strings.ReplaceAll(ref, prefix, "")
	}
	return ref
}

func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func ExtractPathFromRouter(router string) string {
	parts := strings.Fields(router)
	if len(parts) >= 1 {
		path := parts[0]
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		return path
	}
	return ""
}

func ExtractMethodFromHandler(handler string) string {
	if handler == "" {
		return ""
	}
	parts := strings.Split(handler, ".")
	if len(parts) > 0 {
		method := parts[len(parts)-1]
		method = strings.TrimPrefix(method, "*")
		return method
	}
	return ""
}

func ParseJSONExample(jsonStr string) interface{} {
	if jsonStr == "" {
		return nil
	}
	var result interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil
	}
	return result
}