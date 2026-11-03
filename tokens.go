package linkwarden

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

// Token is a Linkwarden API token, as returned by TokensService.List.
type Token struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	IsSession bool      `json:"isSession"`
	Expires   time.Time `json:"expires"`
	CreatedAt time.Time `json:"createdAt"`
}

// CreatedToken is returned by TokensService.Create; SecretKey is the bearer token
// value and is only ever returned at creation time.
type CreatedToken struct {
	SecretKey string      `json:"secretKey"`
	Token     TokenDetail `json:"token"`
}

type TokenDetail struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	UserID     int        `json:"userId"`
	Token      string     `json:"token"`
	Revoked    bool       `json:"revoked"`
	IsSession  bool       `json:"isSession"`
	Expires    time.Time  `json:"expires"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

// TokensService accesses the /api/v1/tokens endpoints.
type TokensService struct {
	c *Client
}

func (s *TokensService) List(ctx context.Context) ([]Token, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/tokens",
	})
	if err != nil {
		return nil, err
	}

	return decodeResponseEnvelope[[]Token](body)
}

// Create creates a new token named name. expires is passed through verbatim to the
// API as the "expires" field (this instance's existing token-generation flow uses a
// 0-4 expiry-preset index rather than a duration, so the caller controls the value;
// note the OpenAPI spec instead documents it as a duration in seconds, 0 = never).
func (s *TokensService) Create(ctx context.Context, name string, expires int) (*CreatedToken, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPost,
		path:     "/api/v1/tokens",
		jsonBody: map[string]any{"name": name, "expires": expires},
	})
	if err != nil {
		return nil, err
	}

	created, err := decodeResponseEnvelope[CreatedToken](body)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (s *TokensService) Revoke(ctx context.Context, id int) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/api/v1/tokens/" + strconv.Itoa(id),
	})

	return err
}
