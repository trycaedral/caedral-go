package caedral

import "context"

// EmbeddingsService creates text embeddings.
type EmbeddingsService struct {
	client *Client
}

// Create generates embeddings for the given input.
func (s EmbeddingsService) Create(ctx context.Context, req EmbeddingCreateRequest) (*EmbeddingCreateResponse, error) {
	var out EmbeddingCreateResponse
	if err := s.client.doPostJSON(ctx, "/v1/embeddings", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
