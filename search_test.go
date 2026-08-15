package linkwarden

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchService_Links(t *testing.T) {
	assert := require.New(t)

	var gotQuery string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/search", r.URL.Path)
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"Success","data":{"nextCursor":7,"links":[{"id":1,"name":"Go","url":"https://go.dev"}]}}`))
	})

	sort := 2
	links, err := client.Search.Links(t.Context(), SearchOptions{SearchQueryString: "go", Sort: &sort, Cursor: 3})
	assert.NoError(err)
	assert.Equal("cursor=3&searchQueryString=go&sort=2", gotQuery)
	assert.Len(links, 1)
	assert.Equal("Go", links[0].Name)
}

func TestSearchService_Links_OmitsUnsetQueryParams(t *testing.T) {
	assert := require.New(t)

	var gotQuery string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"Success","data":{"links":[]}}`))
	})

	_, err := client.Search.Links(t.Context(), SearchOptions{})
	assert.NoError(err)
	assert.Equal("", gotQuery)
}
