package linkwarden

import (
	"context"
	"net/http"
)

// ImportFormat identifies the source format for MigrationsService.Import.
type ImportFormat int

const (
	ImportFormatLinkwardenJSON ImportFormat = 0
	ImportFormatBookmarksHTML  ImportFormat = 1
	ImportFormatWallabagJSON   ImportFormat = 2
	ImportFormatOmnivoreJSON   ImportFormat = 3
	ImportFormatPocketJSON     ImportFormat = 4
)

// MigrationsService accesses the /api/v1/migration endpoints.
type MigrationsService struct {
	c *Client
}

// Export returns the raw JSON export payload (user info, collections, links), for
// writing straight to a backup file.
func (s *MigrationsService) Export(ctx context.Context) ([]byte, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/migration",
	})
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Import uploads a raw import payload (in the given source format) for background
// processing, returning the API's status message.
func (s *MigrationsService) Import(ctx context.Context, format ImportFormat, data string) (string, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/api/v1/migration",
		jsonBody: map[string]any{
			"format": int(format),
			"data":   data,
		},
	})
	if err != nil {
		return "", err
	}

	return decodeResponseEnvelope[string](body)
}
