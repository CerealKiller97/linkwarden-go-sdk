package linkwarden

import (
	"context"
	"net/http"
)

// DashboardV2Data is the payload returned by the v2 dashboard endpoints.
// CollectionLinks maps a collection ID (as a string key) to that collection's links,
// for each collection section in the dashboard layout.
type DashboardV2Data struct {
	Links               []Link            `json:"links"`
	CollectionLinks     map[string][]Link `json:"collectionLinks,omitempty"`
	NumberOfPinnedLinks int               `json:"numberOfPinnedLinks"`
	NumberOfTags        int               `json:"numberOfTags"`
}

// DashboardLayoutItem configures one widget in the v2 dashboard layout, as sent to
// DashboardService.UpdateV2Layout.
type DashboardLayoutItem struct {
	// Type is one of "STATS", "RECENT_LINKS", "PINNED_LINKS", "COLLECTION".
	Type         string `json:"type"`
	CollectionID *int   `json:"collectionId,omitempty"`
	Enabled      bool   `json:"enabled"`
	// Order is always sent: the API stores it as-is, and rejects a section without one.
	Order int `json:"order"`
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

// GetV2 returns the v2 dashboard: the user's section layout and the data each
// section shows.
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

// UpdateV2Layout replaces the dashboard's section layout and returns the updated
// dashboard.
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
