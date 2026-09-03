package imgur

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.imgur.com"

// AuthMode defines which authorization strategy to use.
type AuthMode int

const (
	// AuthModeAnonymous sends a Client-ID header and is used for public endpoints.
	AuthModeAnonymous AuthMode = iota
	// AuthModeOAuth sends a Bearer token and is used for user-specific endpoints.
	AuthModeOAuth
)

// Config configures the Imgur client.
type Config struct {
	BaseURL     string
	ClientID    string
	AccessToken string
	Mode        AuthMode
	HTTPClient  *http.Client
}

// Client provides shared request/response behavior for all endpoint services.
type Client struct {
	baseURL     string
	clientID    string
	accessToken string
	mode        AuthMode
	httpClient  *http.Client
}

// NewClient constructs a new shared client.
func NewClient(cfg Config) (*Client, error) {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")

	if cfg.Mode == AuthModeOAuth && cfg.AccessToken == "" {
		return nil, fmt.Errorf("access token is required for oauth mode")
	}
	if cfg.Mode == AuthModeAnonymous && cfg.ClientID == "" {
		return nil, fmt.Errorf("client id is required for anonymous mode")
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	return &Client{
		baseURL:     baseURL,
		clientID:    cfg.ClientID,
		accessToken: cfg.AccessToken,
		mode:        cfg.Mode,
		httpClient:  httpClient,
	}, nil
}

// NewRequest builds an authorized request to an Imgur API path.
func (c *Client) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	switch c.mode {
	case AuthModeOAuth:
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	default:
		req.Header.Set("Authorization", "Client-ID "+c.clientID)
	}

	return req, nil
}

// Do executes an HTTP request and wraps non-2xx responses as APIError.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return resp, nil
	}

	defer resp.Body.Close()
	b, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: "unable to read error body"}
	}

	return nil, &APIError{StatusCode: resp.StatusCode, Body: string(b)}
}
