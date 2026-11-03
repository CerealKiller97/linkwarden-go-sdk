package linkwarden

import (
	"context"
	"net/http"
)

// WorkerStats reports the background worker's preservation and search indexing
// queues, as returned by WorkerService.Stats.
type WorkerStats struct {
	Link struct {
		Pending int `json:"pending"`
		Done    int `json:"done"`
		Failed  int `json:"failed"`
	} `json:"link"`
	Search struct {
		Pending int `json:"pending"`
		Done    int `json:"done"`
	} `json:"search"`
}

// PreservationResetAction selects which links WorkerService.ResetPreservation clears.
type PreservationResetAction string

const (
	// PreservationResetAll clears every link's preserved files and re-preserves them.
	PreservationResetAll PreservationResetAction = "allAndRePreserve"
	// PreservationResetBroken clears only links whose archives are broken.
	PreservationResetBroken PreservationResetAction = "allBroken"
)

// WorkerService accesses the /api/v1/worker endpoints. Both require the server
// admin user.
type WorkerService struct {
	c *Client
}

func (s *WorkerService) Stats(ctx context.Context) (*WorkerStats, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/worker",
	})
	if err != nil {
		return nil, err
	}

	stats, err := decodeDataEnvelope[WorkerStats](body)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// ResetPreservation starts clearing preserved files in bulk, returning the API's
// status message.
func (s *WorkerService) ResetPreservation(ctx context.Context, action PreservationResetAction) (string, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodDelete,
		path:     "/api/v1/worker/preservation",
		jsonBody: map[string]any{"action": action},
	})
	if err != nil {
		return "", err
	}

	return decodeResponseEnvelope[string](body)
}
