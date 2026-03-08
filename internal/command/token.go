package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"go-rest-api-cli/internal/auth"
	cfgstore "go-rest-api-cli/internal/config"
)

// TokenCommand fetches an OAuth2 access token and prints or saves it.
type TokenCommand struct{}

func NewTokenCommand() *TokenCommand {
	return &TokenCommand{}
}

func (t *TokenCommand) Name() string        { return "token" }
func (t *TokenCommand) Description() string { return "Fetch an OAuth2 access token" }

func (t *TokenCommand) Run(args []string) error {
	fs := flag.NewFlagSet("token", flag.ContinueOnError)

	profile := fs.String("profile", "", "Profile name to load OAuth2 settings from")

	oauth2TokenURL     := fs.String("oauth2-token-url", "", "OAuth2 token endpoint URL")
	oauth2ClientID     := fs.String("oauth2-client-id", "", "OAuth2 client ID")
	oauth2ClientSecret := fs.String("oauth2-client-secret", "", "OAuth2 client secret")
	oauth2Scopes       := fs.String("oauth2-scopes", "", "OAuth2 space-separated scopes")

	printJSON := fs.Bool("json", false, "Print full token response as JSON (access_token, token_type, expires_in, scope)")
	quiet     := fs.Bool("quiet", false, "Suppress console output (useful with --out)")
	outPath   := fs.String("out", "", "Write token (or JSON) to file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Load profile defaults
	finalTokenURL     := *oauth2TokenURL
	finalClientID     := *oauth2ClientID
	finalClientSecret := *oauth2ClientSecret
	finalScopes       := *oauth2Scopes

	if *profile != "" {
		cfg, err := cfgstore.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		p, ok := cfg.Profiles[*profile]
		if !ok {
			return fmt.Errorf("profile %q not found", *profile)
		}
		if finalTokenURL == "" {
			finalTokenURL = p.OAuth2TokenURL
		}
		if finalClientID == "" {
			finalClientID = p.OAuth2ClientID
		}
		if finalClientSecret == "" {
			finalClientSecret = p.OAuth2ClientSecret
		}
		if finalScopes == "" {
			finalScopes = p.OAuth2Scopes
		}
	}

	var scopes []string
	if finalScopes != "" {
		scopes = strings.Fields(finalScopes)
	}

	o := &auth.OAuth2ClientCredentials{
		TokenURL:     finalTokenURL,
		ClientID:     finalClientID,
		ClientSecret: finalClientSecret,
		Scopes:       scopes,
	}

	tr, err := o.FetchTokenFull()
	if err != nil {
		return err
	}

	// Build the output string
	var output string
	if *printJSON {
		raw, err := json.MarshalIndent(tr, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal token response: %w", err)
		}
		output = string(raw)
	} else {
		output = tr.AccessToken
	}

	if !*quiet {
		fmt.Println(output)
	}

	if *outPath != "" {
		if err := os.WriteFile(*outPath, []byte(output), 0o600); err != nil {
			return fmt.Errorf("write token to file: %w", err)
		}
		if !*quiet {
			fmt.Fprintf(os.Stderr, "token written to %s\n", *outPath)
		}
	}

	return nil
}
