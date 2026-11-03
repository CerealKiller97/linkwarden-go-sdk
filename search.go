package linkwarden

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// SearchOptions filters/paginates SearchService.Links.
type SearchOptions struct {
	SearchQueryString string
	// Sort is left nil to omit the query parameter and let the API pick its default.
	Sort   *int
	Cursor int
}

func (o SearchOptions) toQuery() url.Values {
	q := url.Values{}
	setIfNonEmpty(q, "searchQueryString", o.SearchQueryString)
	if o.Sort != nil {
		q.Set("sort", strconv.Itoa(*o.Sort))
	}
	setIfPositive(q, "cursor", o.Cursor)

	return q
}

type searchData struct {
	NextCursor *int   `json:"nextCursor"`
	Links      []Link `json:"links"`
}

// SearchService accesses the /api/v1/search endpoint.
type SearchService struct {
	c *Client
}

func (s *SearchService) Links(ctx context.Context, opts SearchOptions) ([]Link, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/search",
		query:  opts.toQuery(),
	})
	if err != nil {
		return nil, err
	}

	data, err := decodeDataEnvelope[searchData](body)
	if err != nil {
		return nil, err
	}

	return data.Links, nil
}
