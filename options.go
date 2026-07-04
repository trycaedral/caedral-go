package caedral

import (
	"net/http"
	"time"
)

const defaultBaseURL = "https://api.caedral.com"

// Option configures a Client.
type Option func(*Client)

// WithBaseURL sets the API gateway base URL (e.g. http://localhost:5001).
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = trimTrailingSlash(baseURL)
	}
}

// WithHTTPClient supplies a custom http.Client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.http = httpClient
	}
}

// WithMaxRetries sets retry attempts for idempotent GET requests.
func WithMaxRetries(maxRetries int) Option {
	return func(c *Client) {
		c.maxRetries = maxRetries
	}
}

// WithTimeout sets the per-request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}
