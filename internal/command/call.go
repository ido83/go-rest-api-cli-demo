package command

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"go-rest-api-cli-demo/internal/auth"
	cfgstore "go-rest-api-cli-demo/internal/config"
	"go-rest-api-cli-demo/internal/httpclient"
	"go-rest-api-cli-demo/internal/payload"
)

// CallCommand = "call" subcommand.
type CallCommand struct {
	Factory httpclient.Factory
}

func NewCallCommand(factory httpclient.Factory) *CallCommand {
	return &CallCommand{Factory: factory}
}

func (c *CallCommand) Name() string        { return "call" }
func (c *CallCommand) Description() string { return "Execute a REST API call" }

func (c *CallCommand) Run(args []string) error {
	fs := flag.NewFlagSet("call", flag.ContinueOnError)

	var (
		method       = fs.String("method", "GET", "HTTP method (GET, POST, PUT, DELETE, PATCH...)")
		urlStr       = fs.String("url", "", "Request URL (absolute or relative, when using --profile)")
		profileName  = fs.String("profile", "", "Profile name to use from config")
		inlineJSON   = fs.String("data", "", "Inline JSON body")
		jsonFilePath = fs.String("json-file", "", "Path to JSON file with extra payload")
		timeoutSec   = fs.Int("timeout", 30, "Timeout in seconds")
		insecure     = fs.Bool("insecure", false, "Skip TLS verification (NOT recommended for prod)")

		authType = fs.String("auth", "none", "Auth: none|basic|bearer")
		user     = fs.String("user", "", "Username for basic auth")
		pass     = fs.String("pass", "", "Password for basic auth")
		token    = fs.String("token", "", "Bearer token")

		pretty    = fs.Bool("pretty", false, "Pretty-print JSON responses")
		raw       = fs.Bool("raw", false, "Print only response body (no status/headers)")
		jsonOnly  = fs.Bool("json-only", false, "If response is JSON, print only JSON body")
		outPath   = fs.String("out", "", "Write response body to file")
		retries   = fs.Int("retries", 0, "Number of retries on failure (network/5xx)")
		retryWait = fs.Int("retry-delay", 1, "Delay between retries in seconds")
	)

	headers := HeaderFlag{} // initialized non-nil
	fs.Var(&headers, "header", "HTTP header 'Key: Value' (can be repeated)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *urlStr == "" {
		return fmt.Errorf("--url is required")
	}

	// Load profiles if requested
	var (
		baseURLFromProfile string
		profileHeaders     map[string]string
		profileAuthType    string
		profileUser        string
		profilePass        string
		profileToken       string
	)
	if *profileName != "" {
		cfg, err := cfgstore.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		p, ok := cfg.Profiles[*profileName]
		if !ok {
			return fmt.Errorf("profile %q not found", *profileName)
		}
		baseURLFromProfile = p.BaseURL
		profileHeaders = p.Headers
		profileAuthType = strings.ToLower(p.AuthType)
		profileUser = p.User
		profilePass = p.Pass
		profileToken = p.Token
	}

	// Effective URL (profile base URL + relative path)
	finalURL := *urlStr
	if baseURLFromProfile != "" &&
		!strings.HasPrefix(strings.ToLower(*urlStr), "http://") &&
		!strings.HasPrefix(strings.ToLower(*urlStr), "https://") {
		finalURL = strings.TrimRight(baseURLFromProfile, "/") + "/" + strings.TrimLeft(*urlStr, "/")
	}

	// JSON: load file + inline, merge
	fileMap := map[string]interface{}{}
	inlineMap := map[string]interface{}{}
	var err error

	if *jsonFilePath != "" {
		fileMap, err = payload.LoadJSONFile(*jsonFilePath)
		if err != nil {
			return fmt.Errorf("loading json-file: %w", err)
		}
	}

	if *inlineJSON != "" {
		inlineMap, err = payload.ParseJSONInline(*inlineJSON)
		if err != nil {
			return fmt.Errorf("parsing inline JSON: %w", err)
		}
	}

	// Merge headers: profile headers first, then CLI overrides
	effectiveHeaders := HeaderFlag{}
	if profileHeaders != nil {
		for k, v := range profileHeaders {
			effectiveHeaders[k] = v
		}
	}
	for k, v := range headers {
		effectiveHeaders[k] = v
	}

	var body []byte
	if len(fileMap) > 0 || len(inlineMap) > 0 {
		merged := payload.Merge(fileMap, inlineMap)
		body, err = json.Marshal(merged)
		if err != nil {
			return fmt.Errorf("marshalling merged JSON: %w", err)
		}

		// Ensure Content-Type if not set
		if _, ok := effectiveHeaders["Content-Type"]; !ok {
			effectiveHeaders["Content-Type"] = "application/json"
		}
	}

	// Choose auth strategy (profile defaults + CLI overrides)
	finalAuthType := strings.ToLower(*authType)
	finalUser := *user
	finalPass := *pass
	finalToken := *token

	// Use profile defaults if CLI didn't override
	if *profileName != "" {
		if (finalAuthType == "" || finalAuthType == "none") && profileAuthType != "" && profileAuthType != "none" {
			finalAuthType = profileAuthType
		}
		if finalUser == "" && profileUser != "" {
			finalUser = profileUser
		}
		if finalPass == "" && profilePass != "" {
			finalPass = profilePass
		}
		if finalToken == "" && profileToken != "" {
			finalToken = profileToken
		}
	}

	var authStrategy auth.Strategy = auth.NoAuth{}
	switch finalAuthType {
	case "basic":
		authStrategy = auth.Basic{User: finalUser, Pass: finalPass}
	case "bearer":
		authStrategy = auth.Bearer{Token: finalToken}
	case "", "none":
		// default no auth
	default:
		return fmt.Errorf("unknown auth type: %s", finalAuthType)
	}

	cfg := httpclient.Config{
		Method:        strings.ToUpper(*method),
		URL:           finalURL,
		Headers:       effectiveHeaders,
		Body:          body,
		Timeout:       time.Duration(*timeoutSec) * time.Second,
		Auth:          authStrategy,
		SkipTLSVerify: *insecure,
	}

	// Print request preview (once)
	reqPreview, _, err := c.Factory.Build(cfg)
	if err != nil {
		return fmt.Errorf("build request preview: %w", err)
	}

	fmt.Println("=== Request ===")
	fmt.Printf("%s %s\n", reqPreview.Method, reqPreview.URL.String())
	for k, v := range reqPreview.Header {
		fmt.Printf("%s: %s\n", k, strings.Join(v, ", "))
	}
	if len(body) > 0 {
		fmt.Println()
		fmt.Println("Body:")
		fmt.Println(string(body))
	}

	// Retry logic
	attempts := *retries + 1
	if attempts < 1 {
		attempts = 1
	}

	var (
		resp    *http.Response
		lastErr error
	)

	for i := 0; i < attempts; i++ {
		req, client, err := c.Factory.Build(cfg)
		if err != nil {
			return fmt.Errorf("build request: %w", err)
		}

		resp, err = client.Do(req)
		if err != nil {
			lastErr = err
		} else if resp.StatusCode >= 500 && resp.StatusCode <= 599 {
			lastErr = fmt.Errorf("received HTTP %d", resp.StatusCode)
		} else {
			lastErr = nil
			break
		}

		if resp != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}

		if i < attempts-1 {
			time.Sleep(time.Duration(*retryWait) * time.Second)
		}
	}

	if lastErr != nil {
		return fmt.Errorf("request failed after %d attempt(s): %w", attempts, lastErr)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	// Decide how to print based on flags
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	isJSON := strings.HasPrefix(contentType, "application/json")

	bodyToPrint := respBody
	if isJSON && *pretty {
		var buf bytes.Buffer
		if err := json.Indent(&buf, respBody, "", "  "); err == nil {
			bodyToPrint = buf.Bytes()
		}
	}

	// json-only overrides raw if both set
	if *jsonOnly {
		// Only print body (pretty if requested)
		fmt.Println(string(bodyToPrint))
	} else if *raw {
		// Raw body only
		fmt.Println(string(bodyToPrint))
	} else {
		// Default: status + headers + body
		fmt.Println("\n=== Response ===")
		fmt.Printf("Status: %s\n", resp.Status)
		for k, v := range resp.Header {
			fmt.Printf("%s: %s\n", k, strings.Join(v, ", "))
		}
		fmt.Println()
		fmt.Println(string(bodyToPrint))
	}

	// Save to file if requested
	if *outPath != "" {
		if err := os.WriteFile(*outPath, bodyToPrint, 0o644); err != nil {
			return fmt.Errorf("failed to write response to file: %w", err)
		}
	}

	return nil
}
