package linkwarden

import (
	"context"
	"net/http"
	"net/url"
)

// Favicon is a binary favicon image retrieved from FaviconService.Get.
type Favicon struct {
	Data        []byte
	ContentType string
}

// FaviconService accesses the /api/v1/getFavicon endpoint.
type FaviconService struct {
	c *Client
}

func (s *FaviconService) Get(ctx context.Context, targetURL string) (*Favicon, error) {
	q := url.Values{}
	q.Set("url", targetURL)

	resp, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/getFavicon",
		query:  q,
	})
	if err != nil {
		return nil, err
	}

	return &Favicon{Data: body, ContentType: resp.Header.Get("Content-Type")}, nil
}
