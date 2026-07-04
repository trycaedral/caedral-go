package caedral

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultMaxRetries = 3
	defaultTimeout    = 120 * time.Second
)

// Client is the Caedral API client.
type Client struct {
	apiKey     string
	baseURL    string
	http       *http.Client
	maxRetries int
	timeout    time.Duration

	Chat       ChatService
	Models     ModelsService
	Usage      UsageService
	Embeddings EmbeddingsService
	Images     ImagesService
	Audio      AudioService
	Rerank     RerankService
}

// NewClient creates a Caedral API client.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("caedral: apiKey is required")
	}

	c := &Client{
		apiKey:     strings.TrimSpace(apiKey),
		baseURL:    defaultBaseURL,
		http:       &http.Client{Timeout: defaultTimeout},
		maxRetries: defaultMaxRetries,
		timeout:    defaultTimeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.http.Timeout == 0 {
		c.http.Timeout = c.timeout
	}

	c.Chat = ChatService{
		client:      c,
		Completions: CompletionsService{client: c},
	}
	c.Models = ModelsService{client: c}
	c.Usage = UsageService{client: c}
	c.Embeddings = EmbeddingsService{client: c}
	c.Images = ImagesService{client: c}
	c.Audio = AudioService{client: c}
	c.Rerank = RerankService{client: c}

	return c, nil
}

func trimTrailingSlash(s string) string {
	return strings.TrimRight(strings.TrimSpace(s), "/")
}

func (c *Client) doGet(ctx context.Context, path string, out any) error {
	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		lastErr = c.doJSON(ctx, http.MethodGet, path, nil, out)
		if lastErr == nil || !shouldRetry(lastErr, attempt, c.maxRetries) {
			return lastErr
		}
		time.Sleep(backoff(attempt))
	}
	return lastErr
}

func (c *Client) doPostJSON(ctx context.Context, path string, body any, out any) error {
	return c.doJSON(ctx, http.MethodPost, path, body, out)
}

func (c *Client) doPostStream(ctx context.Context, path string, body any) (*http.Response, error) {
	return c.doRaw(ctx, http.MethodPost, path, body)
}

func (c *Client) doJSON(ctx context.Context, method, path string, body any, out any) error {
	resp, err := c.doRaw(ctx, method, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("caedral: read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return parseAPIError(resp.StatusCode, raw)
	}

	if out == nil || len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("caedral: decode response: %w", err)
	}
	return nil
}

func (c *Client) doRaw(ctx context.Context, method, path string, body any) (*http.Response, error) {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("caedral: encode request: %w", err)
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return nil, fmt.Errorf("caedral: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("caedral: request failed: %w", err)
	}
	return resp, nil
}

func parseAPIError(statusCode int, raw []byte) error {
	var parsed any
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &parsed); err != nil {
			parsed = string(raw)
		}
	}
	return NewAPIError(statusCode, parsed)
}

func shouldRetry(err error, attempt, maxRetries int) bool {
	if attempt >= maxRetries {
		return false
	}
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == 502 || apiErr.StatusCode == 503
	}
	return true
}

func backoff(attempt int) time.Duration {
	return time.Duration(100*(1<<attempt)) * time.Millisecond
}
