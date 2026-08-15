package linkwarden

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFaviconService_Get(t *testing.T) {
	assert := require.New(t)

	var gotQuery string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/getFavicon", r.URL.Path)
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "image/x-icon")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("fake-icon-bytes"))
	})

	favicon, err := client.Favicon.Get(t.Context(), "https://example.com")
	assert.NoError(err)
	assert.Equal("url=https%3A%2F%2Fexample.com", gotQuery)
	assert.Equal("image/x-icon", favicon.ContentType)
	assert.Equal("fake-icon-bytes", string(favicon.Data))
}

func TestFaviconService_Get_NoContent(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	favicon, err := client.Favicon.Get(t.Context(), "https://example.com")
	assert.NoError(err)
	assert.Empty(favicon.Data)
}
