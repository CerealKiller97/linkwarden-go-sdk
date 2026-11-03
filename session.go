package linkwarden

import (
	"context"
	"net/http"
)

// CreateSessionRequest is the payload for SessionService.Create.
type CreateSessionRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	SessionName string `json:"sessionName,omitempty"`
}

// SessionService accesses the /api/v1/session endpoint.
type SessionService struct {
	c *Client
}

// Create signs in with a username and password, returning a session token usable
// as the bearer token for New. A 401 means the account's email isn't verified yet.
func (s *SessionService) Create(ctx context.Context, req CreateSessionRequest) (string, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPost,
		path:     "/api/v1/session",
		jsonBody: req,
	})
	if err != nil {
		return "", err
	}

	session, err := decodeResponseEnvelope[struct {
		Token string `json:"token"`
	}](body)
	if err != nil {
		return "", err
	}

	return session.Token, nil
}
