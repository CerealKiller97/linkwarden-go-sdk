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
// Since Linkwarden v2.17 it succeeds whether or not an account has that email, so
// the answer doesn't reveal which emails are registered. Resetting the password signs
// the user out everywhere.
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

// RequestVerificationEmail emails a 6-digit verification code to the unverified
// account with the given email address or username. Sign in with that code in
// CreateSessionRequest.Code to verify the email. It fails with a 400 APIError when
// the instance has no email server configured or the email is already verified.
func (s *AuthService) RequestVerificationEmail(ctx context.Context, emailOrUsername string) (string, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPost,
		path:     "/api/v1/auth/request-verification-email",
		jsonBody: map[string]any{"email": emailOrUsername},
	})
	if err != nil {
		return "", err
	}

	return decodeResponseEnvelope[string](body)
}

// VerifyEmail confirms a new email address after an email change, using the token
// from the verification email. Since Linkwarden v2.17, a new account's email is
// verified with a code instead: see RequestVerificationEmail.
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
