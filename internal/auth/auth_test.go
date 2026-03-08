package auth

// auth_test.go - unit tests for auth strategies (Basic, Bearer, OAuth2).
//
// How to run from repo root:
//   go test ./...
// or just this package:
//   go test ./internal/auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestOAuth2ClientCredentials_FetchAndApply(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad form", http.StatusBadRequest)
			return
		}
		if r.FormValue("grant_type") != "client_credentials" {
			http.Error(w, "wrong grant_type", http.StatusBadRequest)
			return
		}
		if r.FormValue("client_id") != "myid" || r.FormValue("client_secret") != "mysecret" {
			http.Error(w, "wrong credentials", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TokenResponse{AccessToken: "test-token-abc", TokenType: "Bearer"})
	}))
	defer server.Close()

	o := &OAuth2ClientCredentials{
		TokenURL:     server.URL + "/token",
		ClientID:     "myid",
		ClientSecret: "mysecret",
	}

	if err := o.FetchToken(); err != nil {
		t.Fatalf("FetchToken failed: %v", err)
	}

	req, _ := http.NewRequest("GET", "https://api.example.com/data", nil)
	o.Apply(req)

	got := req.Header.Get("Authorization")
	if got != "Bearer test-token-abc" {
		t.Errorf("expected 'Bearer test-token-abc', got %q", got)
	}
}

func TestOAuth2ClientCredentials_FetchWithScopes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.FormValue("scope") != "read write" {
			http.Error(w, "wrong scope", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TokenResponse{AccessToken: "scoped-token", TokenType: "Bearer"})
	}))
	defer server.Close()

	o := &OAuth2ClientCredentials{
		TokenURL:     server.URL + "/token",
		ClientID:     "id",
		ClientSecret: "secret",
		Scopes:       []string{"read", "write"},
	}

	if err := o.FetchToken(); err != nil {
		t.Fatalf("FetchToken failed: %v", err)
	}

	req, _ := http.NewRequest("GET", "https://api.example.com/data", nil)
	o.Apply(req)

	if got := req.Header.Get("Authorization"); got != "Bearer scoped-token" {
		t.Errorf("expected 'Bearer scoped-token', got %q", got)
	}
}

func TestOAuth2ClientCredentials_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TokenResponse{Error: "invalid_client", ErrorDesc: "Unknown client"})
	}))
	defer server.Close()

	o := &OAuth2ClientCredentials{
		TokenURL:     server.URL + "/token",
		ClientID:     "bad",
		ClientSecret: "bad",
	}

	err := o.FetchToken()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid_client") {
		t.Errorf("expected error to contain 'invalid_client', got %q", err.Error())
	}
}

func TestOAuth2ClientCredentials_MissingTokenURL(t *testing.T) {
	o := &OAuth2ClientCredentials{ClientID: "id", ClientSecret: "secret"}
	if err := o.FetchToken(); err == nil {
		t.Fatal("expected error for missing TokenURL")
	}
}

func TestOAuth2ClientCredentials_MissingClientID(t *testing.T) {
	o := &OAuth2ClientCredentials{TokenURL: "http://example.com/token", ClientSecret: "secret"}
	if err := o.FetchToken(); err == nil {
		t.Fatal("expected error for missing ClientID")
	}
}

func TestOAuth2ClientCredentials_FetchTokenFull(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TokenResponse{
			AccessToken: "full-token-xyz",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
			Scope:       "read write",
		})
	}))
	defer server.Close()

	o := &OAuth2ClientCredentials{
		TokenURL:     server.URL + "/token",
		ClientID:     "id",
		ClientSecret: "secret",
	}

	tr, err := o.FetchTokenFull()
	if err != nil {
		t.Fatalf("FetchTokenFull failed: %v", err)
	}
	if tr.AccessToken != "full-token-xyz" {
		t.Errorf("expected access_token 'full-token-xyz', got %q", tr.AccessToken)
	}
	if tr.TokenType != "Bearer" {
		t.Errorf("expected token_type 'Bearer', got %q", tr.TokenType)
	}
	if tr.ExpiresIn != 3600 {
		t.Errorf("expected expires_in 3600, got %d", tr.ExpiresIn)
	}
	if tr.Scope != "read write" {
		t.Errorf("expected scope 'read write', got %q", tr.Scope)
	}
}
