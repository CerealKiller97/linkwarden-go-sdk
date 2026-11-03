package linkwarden

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// tokenExpiresNever is the API's "expires" preset for a token that never expires.
const tokenExpiresNever = 4

func createTestToken(t *testing.T, c *Client) *CreatedToken {
	t.Helper()
	assert := require.New(t)

	created, err := c.Tokens.Create(t.Context(), uniqueName("token"), tokenExpiresNever)
	assert.NoError(err)
	t.Cleanup(func() { _ = c.Tokens.Revoke(t.Context(), created.Token.ID) })

	return created
}

func TestTokensServiceCreate(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)

	// Act
	created, err := client.Tokens.Create(t.Context(), uniqueName("token"), tokenExpiresNever)

	// Assert
	assert.NoError(err)
	t.Cleanup(func() { _ = client.Tokens.Revoke(t.Context(), created.Token.ID) })
	assert.NotEmpty(created.SecretKey)
	assert.NotZero(created.Token.ID)
	assert.False(created.Token.Revoked)
	// The new secret works as a bearer token on its own.
	_, err = New(os.Getenv("LINKWARDEN_URL"), created.SecretKey).Collections.List(t.Context())
	assert.NoError(err)
}

func TestTokensServiceList(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	created := createTestToken(t, client)

	// Act
	tokens, err := client.Tokens.List(t.Context())

	// Assert
	assert.NoError(err)
	var names []string
	for _, tok := range tokens {
		names = append(names, tok.Name)
	}
	assert.Contains(names, created.Token.Name)
}

func TestTokensServiceRevoke(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	created := createTestToken(t, client)

	// Act
	err := client.Tokens.Revoke(t.Context(), created.Token.ID)

	// Assert
	assert.NoError(err)
	_, err = New(os.Getenv("LINKWARDEN_URL"), created.SecretKey).Collections.List(t.Context())
	assert.Error(err)
}

func TestTokensServiceRevokeUnauthorized(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	created := createTestToken(t, client)
	unauthorized := New(os.Getenv("LINKWARDEN_URL"), "invalid-token")

	// Act
	err := unauthorized.Tokens.Revoke(t.Context(), created.Token.ID)

	// Assert
	assert.Error(err)
	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
	assert.Equal(401, apiErr.StatusCode)
}
