package linkwarden

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFaviconServiceGet(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)

	// Act
	favicon, err := client.Favicon.Get(t.Context(), "https://go.dev")

	// Assert
	assert.NoError(err)
	assert.NotEmpty(favicon.Data)
	assert.Contains(favicon.ContentType, "image/")
}
