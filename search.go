package linkwarden

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// SearchOptions filters/paginates SearchService.Links.
type SearchOptions struct {
	SearchQueryString string
	// Sort is one of the Sort constants. It's left nil to omit the query parameter and
	// let the API pick its default.
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

// Links returns one page of links matching opts.SearchQueryString. It returns nil
// when nothing matches.
func (s *SearchService) Links(ctx context.Context, opts SearchOptions) ([]Link, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/search",
		query:  opts.toQuery(),
	})
	if err != nil {
		return nil, err
	}

	return decodeSearchLinks(body)
}

// decodeSearchLinks reads the links out of a /api/v1/search response. With no
// matches the API sends "data": [] instead of {"links": [...]}, so that returns nil.
func decodeSearchLinks(body []byte) ([]Link, error) {
	raw, err := decodeDataEnvelope[json.RawMessage](body)
	if err != nil {
		return nil, err
	}

	if bytes.HasPrefix(bytes.TrimSpace(raw), []byte("[")) {
		return nil, nil
	}

	var data searchData
	if err := decodeJSON(raw, &data); err != nil {
		return nil, err
	}

	return data.Links, nil
}
