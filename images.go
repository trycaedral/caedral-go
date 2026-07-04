package caedral

import "context"

// ImagesService generates images.
type ImagesService struct {
	client *Client
}

// Generate creates an image from a prompt.
func (s ImagesService) Generate(ctx context.Context, req ImageGenerateRequest) (*ImageGenerateResponse, error) {
	var out ImageGenerateResponse
	if err := s.client.doPostJSON(ctx, "/v1/images/generations", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
