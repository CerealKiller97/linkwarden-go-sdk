package linkwarden

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// The test instance has no USER_CONTENT_DOMAIN, so the preserved endpoints can only be
// checked for the error they return; the happy path needs a separately hosted domain.

func TestPreservedService_Token_NoUserContentDomain(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	link := createTestLink(t, client, createTestCollection(t, client).ID)

	_, err := client.Preserved.Token(t.Context(), link.ID, ArchiveFormatPNG)
	assert.Error(err)

	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
	assert.Equal(http.StatusBadRequest, apiErr.StatusCode)
}

func TestAvatarService_Get_NoImage(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	user, _ := createTestUser(t, client)

	_, err := client.Avatar.Get(t.Context(), user.ID)
	assert.True(IsNotFound(err), "got %v", err)
}
