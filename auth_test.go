package linkwarden

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// newAnonymousClient returns a client with no token, for the unauthenticated endpoints.
func newAnonymousClient(t *testing.T) *Client {
	t.Helper()

	newTestClient(t) // skips or fails the test when no instance is configured

	return New(os.Getenv("LINKWARDEN_URL"), "")
}

func TestSessionServiceCreate(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	user, password := createTestUser(t, client)

	// Act
	token, err := newAnonymousClient(t).Session.Create(t.Context(), CreateSessionRequest{
		Username:    user.Username,
		Password:    password,
		SessionName: "sdk-test",
	})

	// Assert
	assert.NoError(err)
	assert.NotEmpty(token)
	// The session token works as a bearer token.
	me, err := New(os.Getenv("LINKWARDEN_URL"), token).Users.Me(t.Context())
	assert.NoError(err)
	assert.Equal(user.ID, me.ID)
}

func TestSessionServiceCreateWrongPassword(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	user, _ := createTestUser(t, client)

	// Act
	_, err := newAnonymousClient(t).Session.Create(t.Context(), CreateSessionRequest{
		Username: user.Username,
		Password: "wrong-password",
	})

	// Assert
	assert.Error(err)
	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
}

func TestAuthServiceInvalidInput(t *testing.T) {
	client := newAnonymousClient(t)

	// The test instance has no email provider and no reset tokens, so each call is
	// rejected; this checks the requests reach the right endpoints and decode errors.
	tests := map[string]func() (string, error){
		"ForgotPassword": func() (string, error) {
			return client.Auth.ForgotPassword(t.Context(), uniqueName("nobody")+"@example.com")
		},
		"ResetPassword": func() (string, error) {
			return client.Auth.ResetPassword(t.Context(), "invalid-token", "a-long-enough-password")
		},
		"VerifyEmail": func() (string, error) {
			return client.Auth.VerifyEmail(t.Context(), "invalid-token")
		},
	}

	for name, call := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			assert := require.New(t)

			// Act
			_, err := call()

			// Assert
			assert.Error(err)
			var apiErr *APIError
			assert.ErrorAs(err, &apiErr)
			assert.Equal(http.StatusBadRequest, apiErr.StatusCode)
			assert.NotEmpty(apiErr.Message)
		})
	}
}

func TestLoginsServiceGet(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newAnonymousClient(t)

	// Act
	cfg, err := client.Logins.Get(t.Context())

	// Assert
	assert.NoError(err)
	assert.NotNil(cfg.CredentialsEnabled)
	assert.Equal("true", *cfg.CredentialsEnabled)
}
