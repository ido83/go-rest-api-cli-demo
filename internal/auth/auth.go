package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Strategy is the auth interface.
type Strategy interface {
	Apply(req *http.Request)
}

type NoAuth struct{}

func (NoAuth) Apply(req *http.Request) {}

type Basic struct {
	User string
	Pass string
}

func (b Basic) Apply(req *http.Request) {
	req.SetBasicAuth(b.User, b.Pass)
}

type Bearer struct {
	Token string
}

func (b Bearer) Apply(req *http.Request) {
	if b.Token != "" {
		req.Header.Set("Authorization", "Bearer "+b.Token)
	}
}

// OAuth2ClientCredentials implements the OAuth 2.0 Client Credentials grant.
// Call FetchToken once before using Apply.
type OAuth2ClientCredentials struct {
	TokenURL     string
	ClientID     string
	ClientSecret string
	Scopes       []string
	cachedToken  string
}

// TokenResponse holds the fields returned by a token endpoint.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	Scope       string `json:"scope,omitempty"`
	Error       string `json:"error,omitempty"`
	ErrorDesc   string `json:"error_description,omitempty"`
}

// FetchTokenFull performs the client_credentials grant and returns the full token response.
func (o *OAuth2ClientCredentials) FetchTokenFull() (*TokenResponse, error) {
	if o.TokenURL == "" {
		return nil, fmt.Errorf("oauth2: token URL is required")
	}
	if o.ClientID == "" {
		return nil, fmt.Errorf("oauth2: client ID is required")
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", o.ClientID)
	form.Set("client_secret", o.ClientSecret)
	if len(o.Scopes) > 0 {
		form.Set("scope", strings.Join(o.Scopes, " "))
	}

	resp, err := http.PostForm(o.TokenURL, form)
	if err != nil {
		return nil, fmt.Errorf("oauth2: token request failed: %w", err)
	}
	defer resp.Body.Close()

	var tr TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, fmt.Errorf("oauth2: failed to decode token response: %w", err)
	}

	if tr.Error != "" {
		return nil, fmt.Errorf("oauth2: token error %q: %s", tr.Error, tr.ErrorDesc)
	}
	if tr.AccessToken == "" {
		return nil, fmt.Errorf("oauth2: token response contained no access_token (HTTP %d)", resp.StatusCode)
	}

	return &tr, nil
}

// FetchToken performs the client_credentials grant and caches the access token.
func (o *OAuth2ClientCredentials) FetchToken() error {
	tr, err := o.FetchTokenFull()
	if err != nil {
		return err
	}
	o.cachedToken = tr.AccessToken
	return nil
}

// Apply sets the cached Bearer token on the request.
func (o *OAuth2ClientCredentials) Apply(req *http.Request) {
	if o.cachedToken != "" {
		req.Header.Set("Authorization", "Bearer "+o.cachedToken)
	}
}
