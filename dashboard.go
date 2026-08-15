package linkwarden

import (
	"context"
	"encoding/json"
	"net/http"
)

// DashboardV2Data is the payload returned by the v2 dashboard endpoints.
// CollectionLinks is kept as raw JSON since its per-collection shape isn't part of
// the published API schema.
type DashboardV2Data struct {
	Links               []Link          `json:"links"`
	CollectionLinks     json.RawMessage `json:"collectionLinks,omitempty"`
	NumberOfPinnedLinks int             `json:"numberOfPinnedLinks"`
	NumberOfTags        int             `json:"numberOfTags"`
}

// DashboardLayoutItem configures one widget in the v2 dashboard layout, as sent to
// DashboardService.UpdateV2Layout.
type DashboardLayoutItem struct {
	// Type is one of "STATS", "RECENT_LINKS", "PINNED_LINKS", "COLLECTION".
	Type         string `json:"type"`
	CollectionID *int   `json:"collectionId,omitempty"`
	Enabled      bool   `json:"enabled"`
	Order        int    `json:"order,omitempty"`
}

// DashboardService accesses the /api/v1/dashboard and /api/v2/dashboard endpoints.
type DashboardService struct {
	c *Client
}

// GetV1 returns the v1 dashboard's recent links.
func (s *DashboardService) GetV1(ctx context.Context) ([]Link, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/dashboard",
	})
	if err != nil {
		return nil, err
	}

	return decodeResponseEnvelope[[]Link](body)
}

func (s *DashboardService) GetV2(ctx context.Context) (*DashboardV2Data, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v2/dashboard",
	})
	if err != nil {
		return nil, err
	}

	data, err := decodeDataEnvelope[DashboardV2Data](body)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *DashboardService) UpdateV2Layout(ctx context.Context, layout []DashboardLayoutItem) (*DashboardV2Data, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPut,
		path:     "/api/v2/dashboard",
		jsonBody: layout,
	})
	if err != nil {
		return nil, err
	}

	data, err := decodeDataEnvelope[DashboardV2Data](body)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
