package caedral

import "context"

// AudioService generates speech audio.
type AudioService struct {
	client *Client
}

// Generate creates speech from text input.
func (s AudioService) Generate(ctx context.Context, req AudioGenerateRequest) (*AudioGenerateResponse, error) {
	var out AudioGenerateResponse
	if err := s.client.doPostJSON(ctx, "/v1/audio/speech", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
