package auth

// auth_test.go - unit tests for auth strategies (Basic, Bearer).
//
// How to run from repo root:
//   go test ./...
// or just this package:
//   go test ./internal/auth

import (
	"encoding/base64"
	"net/http"
	"strings"
	"testing"
)

func TestBasicAuthStrategy(t *testing.T) {
	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	a := Basic{User: "user", Pass: "pass"}
	a.Apply(req)

	got := req.Header.Get("Authorization")
	if got == "" {
		t.Fatalf("expected Authorization header to be set, got empty")
	}

	if !strings.HasPrefix(got, "Basic ") {
		t.Fatalf("expected Authorization header to start with 'Basic ', got %q", got)
	}

	// Check that the credentials are correct (user:pass, base64)
	encoded := strings.TrimPrefix(got, "Basic ")
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("failed to base64-decode Basic auth payload: %v", err)
	}
	if string(decodedBytes) != "user:pass" {
		t.Errorf("expected decoded credentials 'user:pass', got %q", string(decodedBytes))
	}
}

func TestBearerAuthStrategy(t *testing.T) {
	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	a := Bearer{Token: "TOKEN123"}
	a.Apply(req)

	got := req.Header.Get("Authorization")
	if got != "Bearer TOKEN123" {
		t.Errorf("expected 'Bearer TOKEN123', got %q", got)
	}
}

func TestNoAuthStrategy(t *testing.T) {
	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	var a NoAuth
	a.Apply(req)

	if got := req.Header.Get("Authorization"); got != "" {
		t.Errorf("expected no Authorization header for NoAuth, got %q", got)
	}
}
