package linkwarden

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// The test instance has no USER_CONTENT_DOMAIN, so the preserved endpoints can only be
// checked for the error they return; the happy path needs a separately hosted domain.

func TestPreservedServiceTokenNoUserContentDomain(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	link := createTestLink(t, client, createTestCollection(t, client).ID)

	// Act
	_, err := client.Preserved.Token(t.Context(), link.ID, ArchiveFormatPNG)

	// Assert
	assert.Error(err)
	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
	assert.Equal(http.StatusBadRequest, apiErr.StatusCode)
}

func TestAvatarServiceGetNoImage(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	user, _ := createTestUser(t, client)

	// Act
	_, err := client.Avatar.Get(t.Context(), user.ID)

	// Assert
	assert.True(IsNotFound(err), "got %v", err)
}
