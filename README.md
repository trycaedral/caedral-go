# Caedral Go SDK

Official Go client for the [Caedral API](https://caedral.com). OpenAI-compatible request shapes with idiomatic Go patterns (`context.Context`, functional options, channel-based streaming).

> **Module path:** This repository uses `github.com/trycaedral/caedral-go` as a placeholder import path. Replace it with your published module path (for example `github.com/your-org/caedral-go`) before releasing.

## Installation

```bash
go get github.com/trycaedral/caedral-go
```

Local development:

```bash
cd sdk-go
go mod tidy
```

## Quickstart

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/trycaedral/caedral-go"
)

func main() {
    client, err := caedral.NewClient(
        "cd_live_...",
        caedral.WithBaseURL("http://localhost:5001"),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    prompt := "Hello!"
    completion, err := client.Chat.Completions.Create(ctx, caedral.ChatCompletionRequest{
        Model: "caedral-titan",
        Messages: []caedral.ChatMessage{
            {Role: "user", Content: &prompt},
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(completion.Choices[0].Message["content"])
}
```

Run the included example:

```bash
cd sdk-go/examples/quickstart
CAEDRAL_API_KEY=cd_live_... CAEDRAL_BASE_URL=http://localhost:5001 go run .
```

## Configuration

Functional options for `caedral.NewClient`:

| Option | Default | Description |
|--------|---------|-------------|
| `WithBaseURL(url)` | `https://api.caedral.com` | API gateway base URL |
| `WithHTTPClient(c)` | `http.Client{Timeout: 120s}` | Custom HTTP client |
| `WithMaxRetries(n)` | `3` | Retries for idempotent GET requests |
| `WithTimeout(d)` | `120s` | Request timeout |

## Methods

### Chat completions

```go
completion, err := client.Chat.Completions.Create(ctx, caedral.ChatCompletionRequest{
    Model: "caedral-olympus",
    Messages: []caedral.ChatMessage{
        {Role: "user", Content: &prompt},
    },
})
```

**Streaming (channels):**

```go
chunks, errCh, err := client.Chat.Completions.CreateStream(ctx, req)
if err != nil { log.Fatal(err) }

for chunk := range chunks {
    if len(chunk.Choices) > 0 {
        if text, ok := chunk.Choices[0].Delta["content"].(string); ok {
            fmt.Print(text)
        }
    }
}
if streamErr := <-errCh; streamErr != nil {
    log.Fatal(streamErr)
}
```

**Streaming (iterator-style):**

```go
stream, err := client.Chat.Completions.OpenStream(ctx, req)
if err != nil { log.Fatal(err) }
defer stream.Close()

for {
    chunk, err := stream.Recv()
    if err == io.EOF { break }
    if err != nil { log.Fatal(err) }
    _ = chunk
}
```

### Models

```go
models, err := client.Models.List(ctx)
```

### Usage

```go
usage, err := client.Usage.Get(ctx)
fmt.Println(usage.WeeklyPool.Remaining)
```

### Embeddings

```go
result, err := client.Embeddings.Create(ctx, caedral.EmbeddingCreateRequest{
    Model: "caedral-embed",
    Input: "Caedral unifies frontier models behind one API.",
})
```

### Images

```go
image, err := client.Images.Generate(ctx, caedral.ImageGenerateRequest{
    Model:  "caedral-vision",
    Prompt: "A minimal geometric logo on a dark background",
})
```

### Audio

```go
audio, err := client.Audio.Generate(ctx, caedral.AudioGenerateRequest{
    Model: "caedral-voice",
    Input: "Welcome to Caedral.",
    Voice: "alloy",
})
```

### Rerank

```go
ranked, err := client.Rerank.Create(ctx, caedral.RerankCreateRequest{
    Model:     "caedral-rerank",
    Query:     "billing and subscriptions",
    Documents: []string{"Caedral pricing tiers", "Local gateway port 5001"},
    TopN:      intPtr(2),
})
```

## Error handling

```go
completion, err := client.Chat.Completions.Create(ctx, req)
if err != nil {
    var apiErr *caedral.APIError
    if errors.As(err, &apiErr) {
        log.Printf("status=%d type=%s msg=%s", apiErr.StatusCode, apiErr.Type, apiErr.Message)
    }
    return err
}
```

## Integration tests

Requires a running local gateway on port **5001** and `DATABASE_URL` in the repo root `.env` (tests create a temporary API key automatically).

```bash
cd sdk-go
go test ./... -v
```

Optional environment variables:

| Variable | Description |
|----------|-------------|
| `CAEDRAL_BASE_URL` | Gateway URL (default `http://localhost:5001`) |
| `CAEDRAL_TEST_API_KEY` | Use an existing key instead of creating one |

## Async client

An `AsyncCaedral` client is planned as a fast-follow. The synchronous client covers all endpoints today.

## License

MIT
