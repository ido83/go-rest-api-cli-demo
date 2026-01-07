package httpclient

// factory_test.go - tests HTTP client factory (headers, auth, TLS).
//
// How to run from repo root:
//   go test ./...
// or just this package:
//   go test ./internal/httpclient

import (
	"go-rest-api-cli-demo/internal/auth"
	"net/http"
	"testing"
	"time"
)

func TestFactoryBuildWithHeadersAndAuth(t *testing.T) {
	f := Factory{}

	cfg := Config{
		Method: "POST",
		URL:    "https://example.com/api",
		Headers: map[string]string{
			"X-Test": "abc",
		},
		Timeout:       5 * time.Second,
		Auth:          auth.Bearer{Token: "TOKEN456"},
		SkipTLSVerify: true,
	}

	req, client, err := f.Build(cfg)
	if err != nil {
		t.Fatalf("Factory.Build failed: %v", err)
	}

	if req.Method != "POST" {
		t.Errorf("expected method POST, got %s", req.Method)
	}
	if req.URL.String() != "https://example.com/api" {
		t.Errorf("expected URL https://example.com/api, got %s", req.URL.String())
	}

	if got := req.Header.Get("X-Test"); got != "abc" {
		t.Errorf("expected header X-Test == abc, got %q", got)
	}

	if got := req.Header.Get("Authorization"); got != "Bearer TOKEN456" {
		t.Errorf("expected Authorization == 'Bearer TOKEN456', got %q", got)
	}

	// Check timeout
	if client.Timeout != 5*time.Second {
		t.Errorf("expected client timeout 5s, got %v", client.Timeout)
	}

	// Check SkipTLSVerify is true
	tr, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected *http.Transport, got %T", client.Transport)
	}
	if tr.TLSClientConfig == nil || !tr.TLSClientConfig.InsecureSkipVerify {
		t.Errorf("expected InsecureSkipVerify == true, got %#v", tr.TLSClientConfig)
	}
}
