package payload

import (
	"encoding/json"
	"os"
	"strings"
)

// LoadJSONFile loads JSON into a map[string]interface{}.
func LoadJSONFile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(string(data)) == "" {
		return map[string]interface{}{}, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// ParseJSONInline parses inline JSON (string) into a map.
func ParseJSONInline(s string) (map[string]interface{}, error) {
	if strings.TrimSpace(s) == "" {
		return map[string]interface{}{}, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return nil, err
	}
	return m, nil
}

// Merge merges two maps; inlineMap overrides fileMap keys.
func Merge(fileMap, inlineMap map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	for k, v := range fileMap {
		out[k] = v
	}
	for k, v := range inlineMap {
		out[k] = v
	}
	return out
}
