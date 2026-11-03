package linkwarden

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

// Collection is a Linkwarden collection (a folder of links).
type Collection struct {
	ID          int              `json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	Description string           `json:"description,omitempty"`
	Icon        string           `json:"icon,omitempty"`
	IconWeight  string           `json:"iconWeight,omitempty"`
	Color       string           `json:"color,omitempty"`
	ParentID    *int             `json:"parentId,omitempty"`
	IsPublic    bool             `json:"isPublic,omitempty"`
	OwnerID     int              `json:"ownerId,omitempty"`
	CreatedByID int              `json:"createdById,omitempty"`
	CreatedAt   time.Time        `json:"createdAt,omitzero"`
	UpdatedAt   time.Time        `json:"updatedAt,omitzero"`
	Count       *CollectionCount `json:"_count,omitempty"`
	Parent      *CollectionRef   `json:"parent,omitempty"`
	Members     []Member         `json:"members,omitempty"`
}

type CollectionCount struct {
	Links int `json:"links"`
}

type CollectionRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CreateCollectionRequest is the payload for CollectionsService.Create.
type CreateCollectionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
	Icon        string `json:"icon,omitempty"`
	IconWeight  string `json:"iconWeight,omitempty"`
	ParentID    int    `json:"parentId,omitempty"`
}

// UpdateCollectionRequest is the payload for CollectionsService.Update.
type UpdateCollectionRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Color       string   `json:"color,omitempty"`
	Icon        string   `json:"icon,omitempty"`
	IconWeight  string   `json:"iconWeight,omitempty"`
	ParentID    int      `json:"parentId,omitempty"`
	IsPublic    bool     `json:"isPublic,omitempty"`
	Members     []Member `json:"members,omitempty"`
}

// CollectionsService accesses the /api/v1/collections endpoints.
type CollectionsService struct {
	c *Client
}

func (s *CollectionsService) List(ctx context.Context) ([]Collection, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/collections",
	})
	if err != nil {
		return nil, err
	}

	return decodeResponseEnvelope[[]Collection](body)
}

func (s *CollectionsService) Get(ctx context.Context, id int) (*Collection, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/collections/" + strconv.Itoa(id),
	})
	if err != nil {
		return nil, err
	}

	collection, err := decodeResponseEnvelope[Collection](body)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

// Create creates a new collection. It posts to the un-suffixed /api/v1/collections
// path (rather than the {id}-suffixed path the OpenAPI spec documents for create)
// because that's what the live API actually accepts.
func (s *CollectionsService) Create(ctx context.Context, req CreateCollectionRequest) (*Collection, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPost,
		path:     "/api/v1/collections",
		jsonBody: req,
	})
	if err != nil {
		return nil, err
	}

	collection, err := decodeResponseEnvelope[Collection](body)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

func (s *CollectionsService) Update(ctx context.Context, id int, req UpdateCollectionRequest) (*Collection, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPut,
		path:     "/api/v1/collections/" + strconv.Itoa(id),
		jsonBody: req,
	})
	if err != nil {
		return nil, err
	}

	collection, err := decodeResponseEnvelope[Collection](body)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

func (s *CollectionsService) Delete(ctx context.Context, id int) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/api/v1/collections/" + strconv.Itoa(id),
	})

	return err
}
