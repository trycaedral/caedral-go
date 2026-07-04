package caedral

import "context"

// RerankService reranks documents by relevance to a query.
type RerankService struct {
	client *Client
}

// Create reranks documents for the given query.
func (s RerankService) Create(ctx context.Context, req RerankCreateRequest) (*RerankCreateResponse, error) {
	var out RerankCreateResponse
	if err := s.client.doPostJSON(ctx, "/v1/rerank", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
