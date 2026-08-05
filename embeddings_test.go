package caedral_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/trycaedral/caedral-go"
)

func newEmbeddingsTestClient(t *testing.T, handler http.HandlerFunc) (*caedral.Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	client, err := caedral.NewClient("cd_live_test", caedral.WithBaseURL(server.URL))
	if err != nil {
		server.Close()
		t.Fatalf("new client: %v", err)
	}
	return client, server
}

func TestEmbeddingsCreateDefaults(t *testing.T) {
	var gotBody map[string]any
	client, server := newEmbeddingsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/embeddings" {
			t.Errorf("path = %q, want /v1/embeddings", r.URL.Path)
		}
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if err := json.Unmarshal(raw, &gotBody); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"object": "list",
			"model": "caedral-embed-e1-small-v1",
			"data": [{"object": "embedding", "index": 0, "embedding": [0.1]}],
			"usage": {"prompt_tokens": 1, "total_tokens": 1, "completion_tokens": 0}
		}`))
	})
	defer server.Close()

	ctx := context.Background()
	_, err := client.Embeddings.Create(ctx, caedral.EmbeddingCreateRequest{
		Input: "query: hello",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if gotBody["model"] != caedral.DefaultEmbeddingModel {
		t.Fatalf("model = %v, want %q", gotBody["model"], caedral.DefaultEmbeddingModel)
	}
	dims, ok := gotBody["dimensions"].(float64)
	if !ok {
		t.Fatalf("dimensions type = %T, want float64", gotBody["dimensions"])
	}
	if int(dims) != caedral.DefaultEmbeddingDimensions {
		t.Fatalf("dimensions = %d, want %d", int(dims), caedral.DefaultEmbeddingDimensions)
	}
	if gotBody["input"] != "query: hello" {
		t.Fatalf("input = %v, want %q", gotBody["input"], "query: hello")
	}
}

func TestEmbeddingsCreateRejectsUnsupportedModel(t *testing.T) {
	client, server := newEmbeddingsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("request should not reach server")
	})
	defer server.Close()

	ctx := context.Background()
	_, err := client.Embeddings.Create(ctx, caedral.EmbeddingCreateRequest{
		Model: "BAAI/bge-m3",
		Input: "test",
	})
	if err == nil {
		t.Fatal("expected error for unsupported model")
	}
	if err.Error() != "unsupported embedding model: BAAI/bge-m3" {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestEmbeddingsCreateRejectsUnsupportedDimensions(t *testing.T) {
	client, server := newEmbeddingsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("request should not reach server")
	})
	defer server.Close()

	dims := 512
	ctx := context.Background()
	_, err := client.Embeddings.Create(ctx, caedral.EmbeddingCreateRequest{
		Input:      "test",
		Dimensions: &dims,
	})
	if err == nil {
		t.Fatal("expected error for unsupported dimensions")
	}
	if err.Error() != "unsupported embedding dimensions: 512" {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestEmbeddingsCreateLegacyAlias(t *testing.T) {
	var gotBody map[string]any
	client, server := newEmbeddingsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if err := json.Unmarshal(raw, &gotBody); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"object": "list",
			"model": "caedral-embed-e1-small-v1",
			"data": [{"object": "embedding", "index": 0, "embedding": [0.1]}],
			"usage": {"prompt_tokens": 1, "total_tokens": 1, "completion_tokens": 0}
		}`))
	})
	defer server.Close()

	ctx := context.Background()
	_, err := client.Embeddings.Create(ctx, caedral.EmbeddingCreateRequest{
		Model: caedral.LegacyEmbeddingModelAlias,
		Input: "legacy alias text",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if gotBody["model"] != caedral.LegacyEmbeddingModelAlias {
		t.Fatalf("model = %v, want %q", gotBody["model"], caedral.LegacyEmbeddingModelAlias)
	}
}

func TestEmbeddingsCreateInputTypeQuery(t *testing.T) {
	var gotBody map[string]any
	client, server := newEmbeddingsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if err := json.Unmarshal(raw, &gotBody); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"object": "list",
			"model": "caedral-embed-e1-small-v1",
			"data": [{"object": "embedding", "index": 0, "embedding": [0.1]}],
			"usage": {"prompt_tokens": 1, "total_tokens": 1, "completion_tokens": 0}
		}`))
	})
	defer server.Close()

	ctx := context.Background()
	_, err := client.Embeddings.Create(ctx, caedral.EmbeddingCreateRequest{
		Input:     "what is semantic search?",
		InputType: "search_query",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if gotBody["input_type"] != "search_query" {
		t.Fatalf("input_type = %v, want %q", gotBody["input_type"], "search_query")
	}
}

func TestEmbeddingsCreateInputTypeDocument(t *testing.T) {
	var gotBody map[string]any
	client, server := newEmbeddingsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if err := json.Unmarshal(raw, &gotBody); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"object": "list",
			"model": "caedral-embed-e1-small-v1",
			"data": [{"object": "embedding", "index": 0, "embedding": [0.1]}],
			"usage": {"prompt_tokens": 1, "total_tokens": 1, "completion_tokens": 0}
		}`))
	})
	defer server.Close()

	ctx := context.Background()
	_, err := client.Embeddings.Create(ctx, caedral.EmbeddingCreateRequest{
		Input:     "document body for indexing",
		InputType: "search_document",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if gotBody["input_type"] != "search_document" {
		t.Fatalf("input_type = %v, want %q", gotBody["input_type"], "search_document")
	}
}

func TestEmbeddingsCreateEncodingFormatFloat(t *testing.T) {
	var gotBody map[string]any
	client, server := newEmbeddingsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if err := json.Unmarshal(raw, &gotBody); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"object": "list",
			"model": "caedral-embed-e1-small-v1",
			"data": [{"object": "embedding", "index": 0, "embedding": [0.1]}],
			"usage": {"prompt_tokens": 1, "total_tokens": 1, "completion_tokens": 0}
		}`))
	})
	defer server.Close()

	ctx := context.Background()
	_, err := client.Embeddings.Create(ctx, caedral.EmbeddingCreateRequest{
		Input:          "float vectors",
		EncodingFormat: "float",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if gotBody["encoding_format"] != "float" {
		t.Fatalf("encoding_format = %v, want %q", gotBody["encoding_format"], "float")
	}
}

func TestEmbeddingsCreateEncodingFormatBase64(t *testing.T) {
	var gotBody map[string]any
	client, server := newEmbeddingsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if err := json.Unmarshal(raw, &gotBody); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"object": "list",
			"model": "caedral-embed-e1-small-v1",
			"data": [{"object": "embedding", "index": 0, "embedding": [0.1]}],
			"usage": {"prompt_tokens": 1, "total_tokens": 1, "completion_tokens": 0}
		}`))
	})
	defer server.Close()

	ctx := context.Background()
	_, err := client.Embeddings.Create(ctx, caedral.EmbeddingCreateRequest{
		Input:          "base64 vectors",
		EncodingFormat: "base64",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if gotBody["encoding_format"] != "base64" {
		t.Fatalf("encoding_format = %v, want %q", gotBody["encoding_format"], "base64")
	}
}
