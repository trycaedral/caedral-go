package caedral

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// CreateStream performs a streaming chat completion.
// Chunks are sent on the returned channel; the channel closes when the stream ends.
// Non-nil errors are sent on errCh (buffered, capacity 1) when streaming fails.
func (s CompletionsService) CreateStream(ctx context.Context, req ChatCompletionRequest) (<-chan ChatCompletionChunk, <-chan error, error) {
	req.Stream = true
	resp, err := s.client.doPostStream(ctx, "/v1/chat/completions", req)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		raw, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, nil, fmt.Errorf("caedral: read error response: %w", readErr)
		}
		return nil, nil, parseAPIError(resp.StatusCode, raw)
	}

	chunks := make(chan ChatCompletionChunk)
	errCh := make(chan error, 1)

	go func() {
		defer close(chunks)
		defer close(errCh)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
			}

			line := strings.TrimSpace(scanner.Text())
			if line == "" || !strings.HasPrefix(line, "data:") {
				continue
			}
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if data == "" || data == "[DONE]" {
				continue
			}

			var chunk ChatCompletionChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				errCh <- fmt.Errorf("caedral: decode stream chunk: %w", err)
				return
			}
			chunks <- chunk
		}

		if err := scanner.Err(); err != nil {
			errCh <- fmt.Errorf("caedral: stream read: %w", err)
		}
	}()

	return chunks, errCh, nil
}

// ChatCompletionStream wraps a streaming response with Recv/Close helpers.
type ChatCompletionStream struct {
	chunks <-chan ChatCompletionChunk
	errCh  <-chan error
	ctx    context.Context
}

// OpenStream opens a streaming completion using an iterator-style API.
func (s CompletionsService) OpenStream(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionStream, error) {
	chunks, errCh, err := s.CreateStream(ctx, req)
	if err != nil {
		return nil, err
	}
	return &ChatCompletionStream{chunks: chunks, errCh: errCh, ctx: ctx}, nil
}

// Recv returns the next chunk or io.EOF when the stream completes.
func (s *ChatCompletionStream) Recv() (ChatCompletionChunk, error) {
	select {
	case <-s.ctx.Done():
		return ChatCompletionChunk{}, s.ctx.Err()
	case chunk, ok := <-s.chunks:
		if !ok {
			select {
			case err := <-s.errCh:
				if err != nil {
					return ChatCompletionChunk{}, err
				}
			default:
			}
			return ChatCompletionChunk{}, io.EOF
		}
		return chunk, nil
	}
}

// Close drains any remaining chunks.
func (s *ChatCompletionStream) Close() error {
	for range s.chunks {
	}
	return nil
}
