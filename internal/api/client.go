package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// contextKey is an unexported type for context keys in this package.
type contextKey struct{}

// APIError represents an HTTP error response from the API.
type APIError struct {
	StatusCode int
	Body       string
	Method     string
	Path       string
}

// IsAuth reports whether this error was returned for an auth-related status code.
// main.go uses it to map APIError → exit code 2 (auth) vs 1 (generic API).
func (e *APIError) IsAuth() bool {
	return e.StatusCode == 401 || e.StatusCode == 403
}

func (e *APIError) Error() string {
	statusText := http.StatusText(e.StatusCode)
	msg := fmt.Sprintf("API error %d %s: %s %s", e.StatusCode, statusText, e.Method, e.Path)
	if e.Body != "" {
		msg += fmt.Sprintf(": %s", e.Body)
	}
	switch e.StatusCode {
	case 401:
		msg += " (hint: check INTERVALS_API_KEY - may be expired or wrong scope)"
	case 403:
		msg += " (hint: check the permissions on this API key)"
	case 404:
		msg += " (hint: resource not found - check the athlete ID or resource ID; run the relevant list command to see valid IDs)"
	case 409:
		msg += " (hint: concurrent modification - refetch and retry)"
	case 422:
		msg += " (hint: validation failed server-side; run with --dry-run to inspect the body)"
	case 429:
		msg += " (hint: rate limited - slow down requests)"
	case 500:
		msg += " (hint: usually a malformed request body; run with --dry-run to inspect)"
	case 503:
		msg += " (hint: transient - retry with backoff)"
	}
	return msg
}

// Client is an HTTP client for the intervals.icu API.
type Client struct {
	baseURL    string
	apiKey     string
	athleteID  string
	httpClient *http.Client
}

// NewClient creates a new Client with the given base URL, API key, and athlete ID.
func NewClient(baseURL, apiKey, athleteID string) *Client {
	return &Client{
		baseURL:   baseURL,
		apiKey:    apiKey,
		athleteID: athleteID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// WithContext stores a Client in the context.
func WithContext(ctx context.Context, client *Client) context.Context {
	return context.WithValue(ctx, contextKey{}, client)
}

// FromContext retrieves a Client from the context, or nil if not present.
func FromContext(ctx context.Context) *Client {
	c, _ := ctx.Value(contextKey{}).(*Client)
	return c
}

// Do executes an HTTP request using context.Background().
func (c *Client) Do(method, pathTemplate string, params map[string]string, body []byte) (*http.Response, error) {
	return c.DoWithContext(context.Background(), method, pathTemplate, params, body)
}

// DoWithContext executes an HTTP request with the provided context.
// It performs path substitution, sets query params, adds Basic Auth, and handles errors.
func (c *Client) DoWithContext(ctx context.Context, method, pathTemplate string, params map[string]string, body []byte) (*http.Response, error) {
	// Copy params so we don't mutate the caller's map
	remaining := make(map[string]string)
	for k, v := range params {
		remaining[k] = v
	}

	// Perform path substitution
	path := pathTemplate

	// {id} and {athleteId} → athleteID (built-in substitutions)
	path = strings.ReplaceAll(path, "{id}", c.athleteID)
	path = strings.ReplaceAll(path, "{athleteId}", c.athleteID)

	// Other {key} → params[key], consumed from remaining
	for k, v := range remaining {
		placeholder := "{" + k + "}"
		if strings.Contains(path, placeholder) {
			path = strings.ReplaceAll(path, placeholder, v)
			delete(remaining, k)
		}
	}

	// Build URL with remaining params as query string
	rawURL := c.baseURL + path
	if len(remaining) > 0 {
		q := url.Values{}
		for k, v := range remaining {
			q.Set(k, v)
		}
		rawURL += "?" + q.Encode()
	}

	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Basic Auth: user=API_KEY, pass=apiKey
	req.SetBasicAuth("API_KEY", c.apiKey)

	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	// Handle error responses
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Body:       strings.TrimSpace(string(respBody)),
			Method:     method,
			Path:       path,
		}
	}

	return resp, nil
}

// DryRun returns a string describing the request without executing it.
func (c *Client) DryRun(method, pathTemplate string, params map[string]string, body []byte) string {
	// Copy params
	remaining := make(map[string]string)
	for k, v := range params {
		remaining[k] = v
	}

	// Path substitution
	path := pathTemplate
	path = strings.ReplaceAll(path, "{id}", c.athleteID)
	path = strings.ReplaceAll(path, "{athleteId}", c.athleteID)

	for k, v := range remaining {
		placeholder := "{" + k + "}"
		if strings.Contains(path, placeholder) {
			path = strings.ReplaceAll(path, placeholder, v)
			delete(remaining, k)
		}
	}

	rawURL := c.baseURL + path
	if len(remaining) > 0 {
		q := url.Values{}
		for k, v := range remaining {
			q.Set(k, v)
		}
		rawURL += "?" + q.Encode()
	}

	result := fmt.Sprintf("%s %s", method, rawURL)
	if len(body) > 0 {
		preview := string(body)
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		result += fmt.Sprintf(" body: %s", preview)
	}
	return result
}
