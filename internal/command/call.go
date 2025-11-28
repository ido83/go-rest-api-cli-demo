package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"go-rest-api-cli-demo/internal/auth"
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
		urlStr       = fs.String("url", "", "Request URL")
		inlineJSON   = fs.String("data", "", "Inline JSON body")
		jsonFilePath = fs.String("json-file", "", "Path to JSON file with extra payload")
		timeoutSec   = fs.Int("timeout", 30, "Timeout in seconds")
		insecure     = fs.Bool("insecure", false, "Skip TLS verification (NOT recommended for prod)")

		authType = fs.String("auth", "none", "Auth: none|basic|bearer")
		user     = fs.String("user", "", "Username for basic auth")
		pass     = fs.String("pass", "", "Password for basic auth")
		token    = fs.String("token", "", "Bearer token")
	)

	headers := HeaderFlag{} // initialized non-nil

	fs.Var(&headers, "header", "HTTP header 'Key: Value' (can be repeated)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *urlStr == "" {
		return fmt.Errorf("--url is required")
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

	var body []byte
	if len(fileMap) > 0 || len(inlineMap) > 0 {
		merged := payload.Merge(fileMap, inlineMap)
		body, err = json.Marshal(merged)
		if err != nil {
			return fmt.Errorf("marshalling merged JSON: %w", err)
		}

		// Ensure Content-Type if not set
		if _, ok := headers["Content-Type"]; !ok {
			headers["Content-Type"] = "application/json"
		}
	}

	// Choose auth strategy
	var authStrategy auth.Strategy = auth.NoAuth{}
	switch strings.ToLower(*authType) {
	case "basic":
		authStrategy = auth.Basic{User: *user, Pass: *pass}
	case "bearer":
		authStrategy = auth.Bearer{Token: *token}
	case "none":
		// default
	default:
		return fmt.Errorf("unknown auth type: %s", *authType)
	}

	cfg := httpclient.Config{
		Method:        strings.ToUpper(*method),
		URL:           *urlStr,
		Headers:       headers,
		Body:          body,
		Timeout:       time.Duration(*timeoutSec) * time.Second,
		Auth:          authStrategy,
		SkipTLSVerify: *insecure,
	}

	req, client, err := c.Factory.Build(cfg)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	// Print request
	fmt.Println("=== Request ===")
	fmt.Printf("%s %s\n", req.Method, req.URL.String())
	for k, v := range req.Header {
		fmt.Printf("%s: %s\n", k, strings.Join(v, ", "))
	}
	if len(body) > 0 {
		fmt.Println()
		fmt.Println("Body:")
		fmt.Println(string(body))
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	fmt.Println("\n=== Response ===")
	fmt.Printf("Status: %s\n", resp.Status)
	for k, v := range resp.Header {
		fmt.Printf("%s: %s\n", k, strings.Join(v, ", "))
	}
	fmt.Println()
	fmt.Println(string(respBody))

	return nil
}
