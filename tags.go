package linkwarden

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// TagInfo is a Linkwarden tag resource (as returned by list/get/update).
type TagInfo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	OwnerID   int       `json:"ownerId,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitzero"`
	UpdatedAt time.Time `json:"updatedAt,omitzero"`
	Count     *TagCount `json:"_count,omitempty"`
}

type TagCount struct {
	Links int `json:"links"`
}

// ListTagsOptions filters/paginates TagsService.List.
type ListTagsOptions struct {
	Sort   *int
	Cursor int
	Search string
}

func (o ListTagsOptions) toQuery() url.Values {
	q := url.Values{}
	if o.Sort != nil {
		q.Set("sort", strconv.Itoa(*o.Sort))
	}
	setIfPositive(q, "cursor", o.Cursor)
	setIfNonEmpty(q, "search", o.Search)

	return q
}

type tagListData struct {
	Tags       []TagInfo `json:"tags"`
	NextCursor *int      `json:"nextCursor"`
}

// CreateOrUpdateTagInput is one entry in a TagsService.CreateOrUpdateBulk request.
// The archive settings are left nil to keep the tag's (or the user's) defaults.
type CreateOrUpdateTagInput struct {
	Label                   string `json:"label"`
	ArchiveAsScreenshot     *bool  `json:"archiveAsScreenshot,omitempty"`
	ArchiveAsMonolith       *bool  `json:"archiveAsMonolith,omitempty"`
	ArchiveAsPDF            *bool  `json:"archiveAsPDF,omitempty"`
	ArchiveAsReadable       *bool  `json:"archiveAsReadable,omitempty"`
	ArchiveAsWaybackMachine *bool  `json:"archiveAsWaybackMachine,omitempty"`
	AITag                   *bool  `json:"aiTag,omitempty"`
}

// TagsService accesses the /api/v1/tags endpoints.
type TagsService struct {
	c *Client
}

func (s *TagsService) List(ctx context.Context, opts ListTagsOptions) ([]TagInfo, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/tags",
		query:  opts.toQuery(),
	})
	if err != nil {
		return nil, err
	}

	data, err := decodeDataEnvelope[tagListData](body)
	if err != nil {
		return nil, err
	}

	return data.Tags, nil
}

func (s *TagsService) Get(ctx context.Context, id int) (*TagInfo, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/tags/" + strconv.Itoa(id),
	})
	if err != nil {
		return nil, err
	}

	tag, err := decodeResponseEnvelope[TagInfo](body)
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (s *TagsService) CreateOrUpdateBulk(ctx context.Context, tags []CreateOrUpdateTagInput) ([]TagInfo, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPost,
		path:     "/api/v1/tags",
		jsonBody: map[string]any{"tags": tags},
	})
	if err != nil {
		return nil, err
	}

	return decodeResponseEnvelope[[]TagInfo](body)
}

func (s *TagsService) Update(ctx context.Context, id int, name string) (*TagInfo, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPut,
		path:     "/api/v1/tags/" + strconv.Itoa(id),
		jsonBody: map[string]any{"name": name},
	})
	if err != nil {
		return nil, err
	}

	tag, err := decodeResponseEnvelope[TagInfo](body)
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (s *TagsService) Delete(ctx context.Context, id int) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/api/v1/tags/" + strconv.Itoa(id),
	})

	return err
}

func (s *TagsService) BulkDelete(ctx context.Context, ids []int) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method:   http.MethodDelete,
		path:     "/api/v1/tags",
		jsonBody: map[string]any{"tagIds": ids},
	})

	return err
}
