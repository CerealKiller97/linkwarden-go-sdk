package linkwarden

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Tag is a Linkwarden tag, as attached to a link.
type Tag struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
}

// Link is a Linkwarden bookmark.
type Link struct {
	ID            int        `json:"id,omitempty"`
	Name          string     `json:"name"`
	Type          string     `json:"type,omitempty"`
	Description   string     `json:"description,omitempty"`
	CreatedByID   int        `json:"createdById,omitempty"`
	CollectionID  int        `json:"collectionId,omitempty"`
	Icon          string     `json:"icon,omitempty"`
	IconWeight    string     `json:"iconWeight,omitempty"`
	Color         string     `json:"color,omitempty"`
	URL           string     `json:"url"`
	TextContent   string     `json:"textContent,omitempty"`
	Preview       string     `json:"preview,omitempty"`
	Image         string     `json:"image,omitempty"`
	PDF           string     `json:"pdf,omitempty"`
	Readable      string     `json:"readable,omitempty"`
	Monolith      string     `json:"monolith,omitempty"`
	LastPreserved time.Time  `json:"lastPreserved,omitzero"`
	ImportDate    *time.Time `json:"importDate,omitempty"`
	CreatedAt     time.Time  `json:"createdAt,omitzero"`
	UpdatedAt     time.Time  `json:"updatedAt,omitzero"`
	Tags          []Tag      `json:"tags,omitempty"`
	Collection    Collection `json:"collection,omitzero"`
	PinnedBy      []PinnedBy `json:"pinnedBy,omitempty"`
}

// CreateLinkRequest is the payload for LinksService.Create.
type CreateLinkRequest struct {
	Name        string     `json:"name"`
	URL         string     `json:"url"`
	Type        string     `json:"type,omitempty"`
	Description string     `json:"description,omitempty"`
	Tags        []Tag      `json:"tags,omitempty"`
	Collection  Collection `json:"collection"`
}

// UpdateLinkRequest is the payload for LinksService.Update.
type UpdateLinkRequest struct {
	ID          int        `json:"id"`
	Name        string     `json:"name,omitempty"`
	URL         string     `json:"url,omitempty"`
	Description string     `json:"description,omitempty"`
	Icon        string     `json:"icon,omitempty"`
	IconWeight  string     `json:"iconWeight,omitempty"`
	Color       string     `json:"color,omitempty"`
	Tags        []Tag      `json:"tags,omitempty"`
	Collection  Collection `json:"collection,omitzero"`
	PinnedBy    []PinnedBy `json:"pinnedBy,omitempty"`
}

// ListLinksOptions filters/paginates LinksService.List.
type ListLinksOptions struct {
	// Sort is left nil to omit the query parameter and let the API pick its default.
	Sort              *int
	Cursor            int
	CollectionID      int
	TagID             int
	PinnedOnly        bool
	SearchQueryString string
}

func (o ListLinksOptions) toQuery() url.Values {
	q := url.Values{}
	if o.Sort != nil {
		q.Set("sort", strconv.Itoa(*o.Sort))
	}
	setIfPositive(q, "cursor", o.Cursor)
	setIfPositive(q, "collectionId", o.CollectionID)
	setIfPositive(q, "tagId", o.TagID)
	if o.PinnedOnly {
		q.Set("pinnedOnly", "true")
	}
	setIfNonEmpty(q, "searchQueryString", o.SearchQueryString)

	return q
}

// BulkUpdateLinksRequest is the payload for LinksService.BulkUpdate.
type BulkUpdateLinksRequest struct {
	Links              []Link                 `json:"links"`
	RemovePreviousTags bool                   `json:"removePreviousTags,omitempty"`
	NewData            BulkUpdateLinksNewData `json:"newData,omitzero"`
}

type BulkUpdateLinksNewData struct {
	CollectionID int   `json:"collectionId,omitempty"`
	Tags         []Tag `json:"tags,omitempty"`
}

// Highlight is a highlighted passage on a link.
type Highlight struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId"`
	LinkID      int       `json:"linkId"`
	Color       string    `json:"color"`
	Comment     string    `json:"comment,omitempty"`
	StartOffset int       `json:"startOffset"`
	EndOffset   int       `json:"endOffset"`
	Text        string    `json:"text"`
	CreatedAt   time.Time `json:"createdAt,omitzero"`
	UpdatedAt   time.Time `json:"updatedAt,omitzero"`
}

// LinksService accesses the /api/v1/links endpoints.
type LinksService struct {
	c *Client
}

func (s *LinksService) List(ctx context.Context, opts ListLinksOptions) ([]Link, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/links",
		query:  opts.toQuery(),
	})
	if err != nil {
		return nil, err
	}

	return decodeResponseEnvelope[[]Link](body)
}

func (s *LinksService) Get(ctx context.Context, id int) (*Link, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/links/" + strconv.Itoa(id),
	})
	if err != nil {
		return nil, err
	}

	link, err := decodeResponseEnvelope[Link](body)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (s *LinksService) Create(ctx context.Context, req CreateLinkRequest) (*Link, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPost,
		path:     "/api/v1/links",
		jsonBody: req,
	})
	if err != nil {
		return nil, err
	}

	link, err := decodeResponseEnvelope[Link](body)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (s *LinksService) Update(ctx context.Context, id int, req UpdateLinkRequest) (*Link, error) {
	req.ID = id

	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPut,
		path:     "/api/v1/links/" + strconv.Itoa(id),
		jsonBody: req,
	})
	if err != nil {
		return nil, err
	}

	var link Link
	if err := decodeJSON(body, &link); err != nil {
		return nil, err
	}

	return &link, nil
}

func (s *LinksService) Delete(ctx context.Context, id int) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/api/v1/links/" + strconv.Itoa(id),
	})

	return err
}

func (s *LinksService) BulkUpdate(ctx context.Context, req BulkUpdateLinksRequest) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPut,
		path:     "/api/v1/links",
		jsonBody: req,
	})

	return err
}

func (s *LinksService) BulkDelete(ctx context.Context, ids []int) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method:   http.MethodDelete,
		path:     "/api/v1/links",
		jsonBody: map[string]any{"linkIds": ids},
	})

	return err
}

func (s *LinksService) Archive(ctx context.Context, id int) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method: http.MethodPut,
		path:   "/api/v1/links/" + strconv.Itoa(id) + "/archive",
	})

	return err
}

func (s *LinksService) DeleteArchive(ctx context.Context, ids []int) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method:   http.MethodDelete,
		path:     "/api/v1/links/archive",
		jsonBody: map[string]any{"linkIds": ids},
	})

	return err
}

func (s *LinksService) Highlights(ctx context.Context, linkID int) ([]Highlight, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/links/" + strconv.Itoa(linkID) + "/highlights",
	})
	if err != nil {
		return nil, err
	}

	return decodeResponseEnvelope[[]Highlight](body)
}
