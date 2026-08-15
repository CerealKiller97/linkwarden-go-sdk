package linkwarden

import (
	"context"
	"net/http"
	"strconv"
)

// CreateHighlightRequest is the payload for HighlightsService.Create. The API treats
// this as a create-or-update: a highlight matching the same link/offsets is updated
// in place rather than duplicated.
type CreateHighlightRequest struct {
	Color       string `json:"color"`
	Comment     string `json:"comment,omitempty"`
	StartOffset int    `json:"startOffset"`
	EndOffset   int    `json:"endOffset"`
	Text        string `json:"text"`
	LinkID      int    `json:"linkId"`
}

// HighlightsService accesses the /api/v1/highlights endpoints. Listing highlights for
// a link is exposed as LinksService.Highlights, since there is no list-all endpoint.
type HighlightsService struct {
	c *Client
}

func (s *HighlightsService) Create(ctx context.Context, req CreateHighlightRequest) (*Highlight, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPost,
		path:     "/api/v1/highlights",
		jsonBody: req,
	})
	if err != nil {
		return nil, err
	}

	highlight, err := decodeResponseEnvelope[Highlight](body)
	if err != nil {
		return nil, err
	}

	return &highlight, nil
}

func (s *HighlightsService) Delete(ctx context.Context, id int) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/api/v1/highlights/" + strconv.Itoa(id),
	})

	return err
}
