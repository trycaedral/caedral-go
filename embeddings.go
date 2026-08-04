package caedral

import (
	"context"
	"strings"
)

// EmbeddingsService creates text embeddings.
type EmbeddingsService struct {
	client *Client
}

// Create generates embeddings for the given input.
func (s EmbeddingsService) Create(ctx context.Context, req EmbeddingCreateRequest) (*EmbeddingCreateResponse, error) {
	model := strings.TrimSpace(req.Model)
	if model == "" {
		model = DefaultEmbeddingModel
	}
	if err := validateEmbeddingModel(model); err != nil {
		return nil, err
	}

	dimensions := DefaultEmbeddingDimensions
	if req.Dimensions != nil {
		dimensions = *req.Dimensions
	}
	if err := validateEmbeddingDimensions(dimensions); err != nil {
		return nil, err
	}

	body := EmbeddingCreateRequest{
		Model:      model,
		Input:      req.Input,
		Dimensions: &dimensions,
	}

	var out EmbeddingCreateResponse
	if err := s.client.doPostJSON(ctx, "/v1/embeddings", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
