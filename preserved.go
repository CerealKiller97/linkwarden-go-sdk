package linkwarden

import (
	"context"
	"mime"
	"net/http"
	"net/url"
	"strconv"
)

// PreservedFormatURL is a short-lived signed URL for an archived format, as
// returned by PreservedService.Token.
type PreservedFormatURL struct {
	URL string `json:"url"`
}

// PreservedService accesses the /api/v1/preserved endpoints, which serve archived
// formats from the instance's user content domain (USER_CONTENT_DOMAIN).
type PreservedService struct {
	c *Client
}

// Token returns a short-lived signed URL for loading the given archived format of
// a link from the user content domain.
func (s *PreservedService) Token(ctx context.Context, linkID int, format ArchiveFormat) (*PreservedFormatURL, error) {
	q := url.Values{}
	q.Set("linkId", strconv.Itoa(linkID))
	q.Set("format", strconv.Itoa(int(format)))

	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/preserved/token",
		query:  q,
	})
	if err != nil {
		return nil, err
	}

	signed, err := decodeResponseEnvelope[PreservedFormatURL](body)
	if err != nil {
		return nil, err
	}

	return &signed, nil
}

// View streams the archived format a signed token (from Token) points at. The route
// is only served on the user content domain, so the Client must be constructed
// with that domain as its base URL. If download is true, the server sends it as an
// attachment and Filename is populated from Content-Disposition.
func (s *PreservedService) View(ctx context.Context, token string, download bool) (*ArchiveFile, error) {
	q := url.Values{}
	q.Set("token", token)
	if download {
		q.Set("download", "1")
	}

	resp, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/preserved/view",
		query:  q,
	})
	if err != nil {
		return nil, err
	}

	file := &ArchiveFile{
		Data:        body,
		ContentType: resp.Header.Get("Content-Type"),
	}

	if _, params, perr := mime.ParseMediaType(resp.Header.Get("Content-Disposition")); perr == nil {
		file.Filename = params["filename"]
	}

	return file, nil
}
