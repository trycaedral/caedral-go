package caedral

import (
	"context"
	"fmt"
)

// ChatService provides chat completion APIs.
type ChatService struct {
	Completions CompletionsService
	client      *Client
}

// CompletionsService creates chat completions.
type CompletionsService struct {
	client *Client
}

// Create performs a non-streaming chat completion.
func (s CompletionsService) Create(ctx context.Context, req ChatCompletionRequest) (*ChatCompletion, error) {
	if req.Stream {
		return nil, fmt.Errorf("caedral: use CreateStream when req.Stream is true")
	}
	req.Stream = false
	var out ChatCompletion
	if err := s.client.doPostJSON(ctx, "/v1/chat/completions", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
