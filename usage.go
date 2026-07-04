package caedral

import "context"

// UsageService returns account usage summaries.
type UsageService struct {
	client *Client
}

// Get returns the current account usage summary.
func (s UsageService) Get(ctx context.Context) (*UsageSummary, error) {
	var out UsageSummary
	if err := s.client.doGet(ctx, "/v1/usage", &out); err != nil {
		return nil, err
	}
	return &out, nil
}
