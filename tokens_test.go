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

	created, err := c.Tokens.Create(t.Context(), uniqueName("token"), tokenExpiresNever)
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Tokens.Revoke(t.Context(), created.Token.ID) })

	return created
}

func TestTokensService_Create(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	created := createTestToken(t, client)

	assert.NotEmpty(created.SecretKey)
	assert.NotZero(created.Token.ID)
	assert.False(created.Token.Revoked)

	// The new secret works as a bearer token on its own.
	other := New(os.Getenv("LINKWARDEN_URL"), created.SecretKey)
	_, err := other.Collections.List(t.Context())
	assert.NoError(err)
}

func TestTokensService_List(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	created := createTestToken(t, client)

	tokens, err := client.Tokens.List(t.Context())
	assert.NoError(err)

	var names []string
	for _, tok := range tokens {
		names = append(names, tok.Name)
	}
	assert.Contains(names, created.Token.Name)
}

func TestTokensService_Revoke(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	created := createTestToken(t, client)

	assert.NoError(client.Tokens.Revoke(t.Context(), created.Token.ID))

	revoked := New(os.Getenv("LINKWARDEN_URL"), created.SecretKey)
	_, err := revoked.Collections.List(t.Context())
	assert.Error(err)
}

func TestTokensService_Revoke_Unauthorized(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	created := createTestToken(t, client)

	unauthorized := New(os.Getenv("LINKWARDEN_URL"), "invalid-token")
	err := unauthorized.Tokens.Revoke(t.Context(), created.Token.ID)
	assert.Error(err)

	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
	assert.Equal(401, apiErr.StatusCode)
}
