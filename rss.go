package linkwarden

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

// RSSSubscription is a Linkwarden RSS feed subscription.
type RSSSubscription struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	URL          string      `json:"url"`
	OwnerID      int         `json:"ownerId,omitempty"`
	CollectionID int         `json:"collectionId,omitempty"`
	CreatedAt    time.Time   `json:"createdAt,omitzero"`
	UpdatedAt    time.Time   `json:"updatedAt,omitzero"`
	Collection   *Collection `json:"collection,omitempty"`
}

// CreateRSSRequest is the payload for RSSService.Create. Either CollectionID or
// CollectionName may be used to target the destination collection.
type CreateRSSRequest struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	CollectionID   int    `json:"collectionId,omitempty"`
	CollectionName string `json:"collectionName,omitempty"`
}

// RSSService accesses the /api/v1/rss endpoints.
type RSSService struct {
	c *Client
}

func (s *RSSService) List(ctx context.Context) ([]RSSSubscription, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/rss",
	})
	if err != nil {
		return nil, err
	}

	return decodeResponseEnvelope[[]RSSSubscription](body)
}

func (s *RSSService) Create(ctx context.Context, req CreateRSSRequest) (*RSSSubscription, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPost,
		path:     "/api/v1/rss",
		jsonBody: req,
	})
	if err != nil {
		return nil, err
	}

	sub, err := decodeResponseEnvelope[RSSSubscription](body)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (s *RSSService) Delete(ctx context.Context, id int) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/api/v1/rss/" + strconv.Itoa(id),
	})

	return err
}
