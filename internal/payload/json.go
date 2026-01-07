package payload

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
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

// NormalizeHashAlgo normalizes algo names so we can accept things like
// "sha-256", "SHA256", "Sha-1", etc.
func NormalizeHashAlgo(algo string) string {
	s := strings.ToLower(strings.TrimSpace(algo))
	s = strings.ReplaceAll(s, "-", "")
	return s
}

// supportedHashAlgos returns the list of human-facing supported algorithm names.
func SupportedHashAlgos() []string {
	return []string{"md5", "sha-1", "sha-256"}
}

// newHasher returns a hash.Hash for the given normalized algorithm name.
// normalized must be the result of NormalizeHashAlgo.
func newHasher(normalized string) (hash.Hash, error) {
	switch normalized {
	case "md5":
		return md5.New(), nil
	case "sha1":
		// #nosec G505 – sha1 allowed here for compatibility / legacy use
		return sha1.New(), nil
	case "sha256":
		return sha256.New(), nil
	default:
		return nil, fmt.Errorf("unsupported hash algorithm %q (supported: md5, sha-1, sha-256)", normalized)
	}
}

// ComputeFileHash calculates the hash of a file using the requested algorithm
// (md5, sha-1, sha-256) and returns it as a hex string.
func ComputeFileHash(path string, algo string) (string, error) {
	normalized := NormalizeHashAlgo(algo)

	h, err := newHasher(normalized)
	if err != nil {
		return "", err
	}

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
