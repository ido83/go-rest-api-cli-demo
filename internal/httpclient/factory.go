package httpclient

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"
	"time"

	"go-rest-api-cli-demo/internal/auth"
)

// Config holds all data needed to build a request/client.
type Config struct {
	Method        string
	URL           string
	Headers       map[string]string
	Body          []byte
	Timeout       time.Duration
	Auth          auth.Strategy
	SkipTLSVerify bool
}

// Factory builds *http.Request + *http.Client from Config.
type Factory struct{}

func (Factory) Build(cfg Config) (*http.Request, *http.Client, error) {
	var bodyReader io.Reader
	if len(cfg.Body) > 0 {
		bodyReader = strings.NewReader(string(cfg.Body))
	}

	req, err := http.NewRequest(cfg.Method, cfg.URL, bodyReader)
	if err != nil {
		return nil, nil, err
	}

	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}

	if cfg.Auth != nil {
		cfg.Auth.Apply(req)
	}

	transport := &http.Transport{}
	if cfg.SkipTLSVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // be careful in prod
	}

	client := &http.Client{
		Timeout:   cfg.Timeout,
		Transport: transport,
	}

	return req, client, nil
}
