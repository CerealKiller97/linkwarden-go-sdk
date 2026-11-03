package linkwarden

import (
	"context"
	"net/http"
	"net/url"
)

// AuthService accesses the unauthenticated /api/v1/auth endpoints. Each call
// returns the API's status message.
type AuthService struct {
	c *Client
}

// ForgotPassword sends a password reset email to the account with the given email.
func (s *AuthService) ForgotPassword(ctx context.Context, email string) (string, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPost,
		path:     "/api/v1/auth/forgot-password",
		jsonBody: map[string]any{"email": email},
	})
	if err != nil {
		return "", err
	}

	return decodeResponseEnvelope[string](body)
}

// ResetPassword sets a new password (at least 8 characters) using a reset token
// from the ForgotPassword email.
func (s *AuthService) ResetPassword(ctx context.Context, token, password string) (string, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPost,
		path:     "/api/v1/auth/reset-password",
		jsonBody: map[string]any{"token": token, "password": password},
	})
	if err != nil {
		return "", err
	}

	return decodeResponseEnvelope[string](body)
}

// VerifyEmail confirms an email address using the token from the verification email.
func (s *AuthService) VerifyEmail(ctx context.Context, token string) (string, error) {
	q := url.Values{}
	q.Set("token", token)

	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodPost,
		path:   "/api/v1/auth/verify-email",
		query:  q,
	})
	if err != nil {
		return "", err
	}

	return decodeResponseEnvelope[string](body)
}
