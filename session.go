package linkwarden

import (
	"context"
	"net/http"
)

// CreateSessionRequest is the payload for SessionService.Create.
type CreateSessionRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	// Code is the 6-digit code emailed to a user whose email isn't verified yet (see
	// AuthService.RequestVerificationEmail). Sent with the password, it verifies the
	// email and signs in. A wrong or expired code fails with a 400 APIError whose Code
	// is ErrorCodeInvalidVerificationCode. Requires Linkwarden v2.17 or later.
	Code        string `json:"code,omitempty"`
	SessionName string `json:"sessionName,omitempty"`
}

// SessionService accesses the /api/v1/session endpoint.
type SessionService struct {
	c *Client
}

// Create signs in with a username and password, returning a session token usable
// as the bearer token for New.
//
// When the account's email isn't verified yet, it fails with a 401 APIError whose
// Code is ErrorCodeEmailNotVerified and Linkwarden emails a verification code; call
// Create again with that code in req.Code. After too many wrong passwords it fails
// with a 429 APIError whose Code is ErrorCodeTooManyAttempts and whose RetryAfter
// says how long to wait.
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

// Delete signs out: it revokes the session token the client was created with, so it
// can't be used again. It only works for session tokens from Create (or a mobile
// sign-in). An access token from TokensService fails with a 404 APIError and stays
// valid; revoke it with TokensService.Revoke instead. Requires Linkwarden v2.17 or
// later.
func (s *SessionService) Delete(ctx context.Context) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/api/v1/session",
	})

	return err
}
