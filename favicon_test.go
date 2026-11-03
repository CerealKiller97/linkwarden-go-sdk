package linkwarden

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFaviconServiceGet(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)

	favicon, err := client.Favicon.Get(t.Context(), "https://go.dev")
	assert.NoError(err)
	assert.NotEmpty(favicon.Data)
	assert.Contains(favicon.ContentType, "image/")
}
