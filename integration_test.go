package caedral_test

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/caedral/caedral-go"
)

var (
	sharedClient   *caedral.Client
	sharedCleanup  func()
	sharedInitOnce sync.Once
	sharedInitErr  error
)

func TestMain(m *testing.M) {
	loadRootEnv()
	code := m.Run()
	if sharedCleanup != nil {
		sharedCleanup()
	}
	os.Exit(code)
}

func baseURL() string {
	if v := os.Getenv("CAEDRAL_BASE_URL"); v != "" {
		return v
	}
	return "http://localhost:5001"
}

func loadRootEnv() {
	root := filepath.Join("..", ".env")
	data, err := os.ReadFile(root)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

func sharedTestClient(t *testing.T) *caedral.Client {
	t.Helper()

	sharedInitOnce.Do(func() {
		apiKey := os.Getenv("CAEDRAL_TEST_API_KEY")
		if apiKey == "" {
			fixture, err := createTestAPIKey(t)
			if err != nil {
				sharedInitErr = err
				return
			}
			apiKey = fixture.RawKey
			sharedCleanup = fixture.Cleanup
		} else {
			sharedCleanup = func() {}
		}

		client, err := caedral.NewClient(apiKey, caedral.WithBaseURL(baseURL()))
		if err != nil {
			sharedInitErr = err
			return
		}
		sharedClient = client
	})

	if sharedInitErr != nil {
		t.Fatalf("shared test client: %v", sharedInitErr)
	}
	return sharedClient
}

func isRetryableUpstream(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "Upstream error") ||
		strings.Contains(msg, "502") ||
		strings.Contains(msg, "503")
}

func retryUpstream(t *testing.T, fn func() error) {
	t.Helper()
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
		}
		lastErr = fn()
		if lastErr == nil {
			return
		}
		if !isRetryableUpstream(lastErr) {
			t.Fatal(lastErr)
		}
	}
	t.Fatal(lastErr)
}

func TestListModels(t *testing.T) {
	client := sharedTestClient(t)

	ctx := context.Background()
	models, err := client.Models.List(ctx)
	if err != nil {
		t.Fatalf("models.list: %v", err)
	}
	if models.Object != "list" {
		t.Fatalf("expected list object, got %q", models.Object)
	}
	if len(models.Data) < 4 {
		t.Fatalf("expected at least 4 models, got %d", len(models.Data))
	}

	required := map[string]bool{
		"caedral-base":       false,
		"caedral-titan":      false,
		"caedral-olympus":    false,
		"caedral-primordial": false,
	}
	for _, model := range models.Data {
		if _, ok := required[model.ID]; ok {
			required[model.ID] = true
		}
		if model.OwnedBy != "caedral" {
			t.Fatalf("expected owned_by caedral, got %q", model.OwnedBy)
		}
	}
	for id, found := range required {
		if !found {
			t.Fatalf("missing model %s", id)
		}
	}
}

func TestChatCompletion(t *testing.T) {
	client := sharedTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	content := "Reply with exactly: GO SDK OK"
	var completion *caedral.ChatCompletion
	retryUpstream(t, func() error {
		var err error
		completion, err = client.Chat.Completions.Create(ctx, caedral.ChatCompletionRequest{
			Model: "caedral-base",
			Messages: []caedral.ChatMessage{
				{Role: "user", Content: &content},
			},
		})
		return err
	})
	if completion.Object != "chat.completion" {
		t.Fatalf("unexpected object %q", completion.Object)
	}
	if completion.Model != "caedral-base" {
		t.Fatalf("unexpected model %q", completion.Model)
	}
	if len(completion.Choices) == 0 || completion.Choices[0].Message["content"] == nil {
		t.Fatalf("expected assistant content")
	}
}

func TestChatStreaming(t *testing.T) {
	client := sharedTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	content := "Count to 3."
	var chunks <-chan caedral.ChatCompletionChunk
	var errCh <-chan error
	retryUpstream(t, func() error {
		var err error
		chunks, errCh, err = client.Chat.Completions.CreateStream(ctx, caedral.ChatCompletionRequest{
			Model: "caedral-base",
			Messages: []caedral.ChatMessage{
				{Role: "user", Content: &content},
			},
		})
		return err
	})

	var text strings.Builder
	for chunk := range chunks {
		if chunk.Object != "chat.completion.chunk" {
			t.Fatalf("unexpected chunk object %q", chunk.Object)
		}
		if len(chunk.Choices) > 0 {
			if delta, ok := chunk.Choices[0].Delta["content"].(string); ok {
				text.WriteString(delta)
			}
		}
	}
	if err, ok := <-errCh; ok && err != nil {
		t.Fatalf("stream error: %v", err)
	}
	if text.Len() == 0 {
		t.Fatalf("expected streamed text")
	}
}

func TestUsageGet(t *testing.T) {
	client := sharedTestClient(t)

	ctx := context.Background()
	usage, err := client.Usage.Get(ctx)
	if err != nil {
		t.Fatalf("usage.get: %v", err)
	}
	if usage.AccountStatus == "" || usage.Plan == "" {
		t.Fatalf("unexpected usage payload: %+v", usage)
	}
}

func TestInvalidAPIKey(t *testing.T) {
	client, err := caedral.NewClient(
		"cd_live_invalid_integration_test_key",
		caedral.WithBaseURL(baseURL()),
	)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	ctx := context.Background()
	content := "Hello"
	_, err = client.Chat.Completions.Create(ctx, caedral.ChatCompletionRequest{
		Model: "caedral-base",
		Messages: []caedral.ChatMessage{
			{Role: "user", Content: &content},
		},
	})
	var apiErr *caedral.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %v", err)
	}
	if apiErr.Type != "invalid_api_key" || apiErr.StatusCode != 401 {
		t.Fatalf("unexpected api error: %+v", apiErr)
	}

	_, err = client.Usage.Get(ctx)
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError on usage, got %v", err)
	}
	if apiErr.Type != "invalid_api_key" {
		t.Fatalf("unexpected usage error: %+v", apiErr)
	}
}

func TestOpenStreamRecv(t *testing.T) {
	client := sharedTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	content := "Say hi."
	var stream *caedral.ChatCompletionStream
	retryUpstream(t, func() error {
		var err error
		stream, err = client.Chat.Completions.OpenStream(ctx, caedral.ChatCompletionRequest{
			Model: "caedral-base",
			Messages: []caedral.ChatMessage{
				{Role: "user", Content: &content},
			},
		})
		return err
	})
	defer stream.Close()

	var parts int
	for {
		_, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("recv: %v", err)
		}
		parts++
	}
	if parts == 0 {
		t.Fatalf("expected stream chunks")
	}
}
