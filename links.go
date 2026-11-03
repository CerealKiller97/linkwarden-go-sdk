package linkwarden

import (
	"cmp"
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

// BulkUpdateLinksNewData holds the changes LinksService.BulkUpdate applies to each
// link. Leave CollectionID at 0 to keep the links in their current collections.
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

// List returns one page of links matching opts. To get the next page, set
// opts.Cursor to the ID of the last link returned.
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

// Get returns the link with the given ID.
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

// Create saves a new link. Set req.Collection.ID to choose its collection; tags
// are matched by name and created if they don't exist yet.
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

// Update changes the fields set in req and keeps the rest. The API replaces the whole
// link (it requires its collection and full tag list), so Update first fetches the link
// and fills every unset field from it. Tags and PinnedBy, when set, replace the link's
// current lists; a Collection with only an ID moves the link to that collection.
func (s *LinksService) Update(ctx context.Context, id int, req UpdateLinkRequest) (*Link, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	payload, err := s.mergeUpdate(ctx, current, req)
	if err != nil {
		return nil, err
	}

	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPut,
		path:     "/api/v1/links/" + strconv.Itoa(id),
		jsonBody: payload,
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

// updateLinkPayload is the full link the update endpoint expects.
type updateLinkPayload struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	URL         string               `json:"url"`
	Description string               `json:"description"`
	Icon        string               `json:"icon,omitempty"`
	IconWeight  string               `json:"iconWeight,omitempty"`
	Color       string               `json:"color,omitempty"`
	Collection  updateLinkCollection `json:"collection"`
	Tags        []Tag                `json:"tags"`
	PinnedBy    []PinnedBy           `json:"pinnedBy,omitempty"`
}

type updateLinkCollection struct {
	ID      int `json:"id"`
	OwnerID int `json:"ownerId"`
}

func (s *LinksService) mergeUpdate(ctx context.Context, current *Link, req UpdateLinkRequest) (updateLinkPayload, error) {
	p := updateLinkPayload{
		ID:          current.ID,
		Name:        cmp.Or(req.Name, current.Name),
		URL:         cmp.Or(req.URL, current.URL),
		Description: cmp.Or(req.Description, current.Description),
		Icon:        cmp.Or(req.Icon, current.Icon),
		IconWeight:  cmp.Or(req.IconWeight, current.IconWeight),
		Color:       cmp.Or(req.Color, current.Color),
		Collection:  updateLinkCollection{ID: current.Collection.ID, OwnerID: current.Collection.OwnerID},
		Tags:        current.Tags,
		PinnedBy:    current.PinnedBy,
	}

	if req.Tags != nil {
		p.Tags = req.Tags
	}
	if p.Tags == nil {
		p.Tags = []Tag{}
	}
	if req.PinnedBy != nil {
		p.PinnedBy = req.PinnedBy
	}

	if id := req.Collection.ID; id != 0 && id != current.Collection.ID {
		ownerID := req.Collection.OwnerID
		if ownerID == 0 {
			target, err := s.c.Collections.Get(ctx, id)
			if err != nil {
				return p, err
			}
			ownerID = target.OwnerID
		}
		p.Collection = updateLinkCollection{ID: id, OwnerID: ownerID}
	}

	return p, nil
}

// Delete deletes the link with the given ID.
func (s *LinksService) Delete(ctx context.Context, id int) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/api/v1/links/" + strconv.Itoa(id),
	})

	return err
}

// BulkUpdate applies req.NewData to every link in req.Links: it can move them to
// another collection and add tags. Set req.RemovePreviousTags to replace the links'
// tags instead of adding to them.
func (s *LinksService) BulkUpdate(ctx context.Context, req BulkUpdateLinksRequest) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPut,
		path:     "/api/v1/links",
		jsonBody: req,
	})

	return err
}

// BulkDelete deletes the links with the given IDs.
func (s *LinksService) BulkDelete(ctx context.Context, ids []int) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method:   http.MethodDelete,
		path:     "/api/v1/links",
		jsonBody: map[string]any{"linkIds": ids},
	})

	return err
}

// Archive asks the server to archive the link with the given ID again. It returns
// once the request is accepted, not when the new archive is ready.
func (s *LinksService) Archive(ctx context.Context, id int) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method: http.MethodPut,
		path:   "/api/v1/links/" + strconv.Itoa(id) + "/archive",
	})

	return err
}

// DeleteArchive deletes the archived files of the links with the given IDs. The links
// themselves are kept.
func (s *LinksService) DeleteArchive(ctx context.Context, ids []int) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method:   http.MethodDelete,
		path:     "/api/v1/links/archive",
		jsonBody: map[string]any{"linkIds": ids},
	})

	return err
}

// Highlights returns the highlights on the link with the given ID.
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
