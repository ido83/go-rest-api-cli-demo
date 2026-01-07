package payload

// payload_test.go - tests JSON helpers & hashing in internal/payload.
//
// How to run tests (from repo root):
//   go test ./...
// or just this folder:
//   go test ./test

import (
	"os"
	"path/filepath"
	"testing"
)

// helper to create a temp JSON file
func writeTempJSONFile(t *testing.T, dir string, content string) string {
	t.Helper()
	f, err := os.CreateTemp(dir, "payload-*.json")
	if err != nil {
		t.Fatalf("CreateTemp failed: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed writing temp json: %v", err)
	}
	return f.Name()
}

func TestLoadJSONFileAndParseJSONInline(t *testing.T) {
	dir := t.TempDir()
	path := writeTempJSONFile(t, dir, `{"name":"test","value":42}`)

	m, err := LoadJSONFile(path)
	if err != nil {
		t.Fatalf("LoadJSONFile failed: %v", err)
	}
	if got := m["name"]; got != "test" {
		t.Errorf("expected name == %q, got %v", "test", got)
	}
	if got := m["value"]; got != float64(42) { // JSON numbers become float64
		t.Errorf("expected value == 42, got %v", got)
	}

	inline := `{"inline":true}`
	im, err := ParseJSONInline(inline)
	if err != nil {
		t.Fatalf("ParseJSONInline failed: %v", err)
	}
	if got := im["inline"]; got != true {
		t.Errorf("expected inline == true, got %v", got)
	}
}

func TestMerge(t *testing.T) {
	base := map[string]interface{}{
		"a": 1,
		"b": "base",
	}
	override := map[string]interface{}{
		"b": "override",
		"c": true,
	}

	merged := Merge(base, override)

	if merged["a"] != 1.0 && merged["a"] != 1 { // depending on source, may be int/float64
		t.Errorf("expected merged[a] == 1, got %v", merged["a"])
	}
	if merged["b"] != "override" {
		t.Errorf("expected merged[b] == override, got %v", merged["b"])
	}
	if merged["c"] != true {
		t.Errorf("expected merged[c] == true, got %v", merged["c"])
	}
}

func TestNormalizeHashAlgo(t *testing.T) {
	cases := map[string]string{
		"md5":      "md5",
		"MD5":      "md5",
		" md5 ":    "md5",
		"sha1":     "sha1",
		"SHA-1":    "sha1",
		" sha-1 ":  "sha1",
		"sha256":   "sha256",
		"SHA-256":  "sha256",
		" sha-256": "sha256",
	}

	for in, want := range cases {
		got := NormalizeHashAlgo(in)
		if got != want {
			t.Errorf("NormalizeHashAlgo(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestComputeFileHash_KnownContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.bin")

	// well-known test payload
	content := []byte("hello world")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	// Expected hashes for "hello world"
	const (
		wantMD5    = "5eb63bbbe01eeed093cb22bb8f5acdc3"
		wantSHA1   = "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"
		wantSHA256 = "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	)

	tests := []struct {
		algo string
		want string
	}{
		{"md5", wantMD5},
		{"MD5", wantMD5},
		{"sha-1", wantSHA1},
		{"SHA-1", wantSHA1},
		{"sha-256", wantSHA256},
		{"SHA-256", wantSHA256},
	}

	for _, tt := range tests {
		got, err := ComputeFileHash(path, tt.algo)
		if err != nil {
			t.Fatalf("ComputeFileHash(%q) unexpected error: %v", tt.algo, err)
		}
		if got != tt.want {
			t.Errorf("ComputeFileHash(%q) = %q, want %q", tt.algo, got, tt.want)
		}
	}
}
