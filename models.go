package caedral

import "context"

// ModelsService lists available models.
type ModelsService struct {
	client *Client
}

// List returns the public model catalog.
func (s ModelsService) List(ctx context.Context) (*ModelListResponse, error) {
	var out ModelListResponse
	if err := s.client.doGet(ctx, "/v1/models", &out); err != nil {
		return nil, err
	}
	return &out, nil
}
