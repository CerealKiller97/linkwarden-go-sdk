package linkwarden

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinksService_List(t *testing.T) {
	assert := require.New(t)

	var gotMethod, gotPath, gotQuery string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":[{"id":1,"name":"Example","url":"https://example.com"}]}`))
	})

	sort := 2
	links, err := client.Links.List(t.Context(), ListLinksOptions{
		Sort:         &sort,
		CollectionID: 3,
		PinnedOnly:   true,
	})
	assert.NoError(err)
	assert.Equal(http.MethodGet, gotMethod)
	assert.Equal("/api/v1/links", gotPath)
	assert.Equal("collectionId=3&pinnedOnly=true&sort=2", gotQuery)
	assert.Len(links, 1)
	assert.Equal("Example", links[0].Name)
}

func TestLinksService_Get(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/links/5", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":5,"name":"Example","url":"https://example.com"}}`))
	})

	link, err := client.Links.Get(t.Context(), 5)
	assert.NoError(err)
	assert.Equal(5, link.ID)
}

func TestLinksService_Get_NotFound(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"response":"Link not found."}`))
	})

	_, err := client.Links.Get(t.Context(), 999)
	assert.Error(err)
	assert.True(IsNotFound(err))
}

func TestLinksService_Create(t *testing.T) {
	assert := require.New(t)

	var gotBody CreateLinkRequest
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/api/v1/links", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":9,"name":"New","url":"https://example.com"}}`))
	})

	link, err := client.Links.Create(t.Context(), CreateLinkRequest{
		Name:       "New",
		URL:        "https://example.com",
		Collection: Collection{ID: 2},
	})
	assert.NoError(err)
	assert.Equal(9, link.ID)
	assert.Equal("New", gotBody.Name)
	assert.Equal(2, gotBody.Collection.ID)
}

func TestLinksService_Update(t *testing.T) {
	assert := require.New(t)

	var gotBody UpdateLinkRequest
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPut, r.Method)
		assert.Equal("/api/v1/links/7", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":7,"name":"Renamed","url":"https://example.com"}`))
	})

	link, err := client.Links.Update(t.Context(), 7, UpdateLinkRequest{Name: "Renamed"})
	assert.NoError(err)
	assert.Equal(7, link.ID)
	assert.Equal("Renamed", link.Name)
	assert.Equal(7, gotBody.ID)
}

func TestLinksService_Delete(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/v1/links/3", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Link deleted."}`))
	})

	assert.NoError(client.Links.Delete(t.Context(), 3))
}

func TestLinksService_BulkUpdate(t *testing.T) {
	assert := require.New(t)

	var gotBody BulkUpdateLinksRequest
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPut, r.Method)
		assert.Equal("/api/v1/links", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Links updated."}`))
	})

	err := client.Links.BulkUpdate(t.Context(), BulkUpdateLinksRequest{
		Links:              []Link{{ID: 1}, {ID: 2}},
		RemovePreviousTags: true,
	})
	assert.NoError(err)
	assert.Len(gotBody.Links, 2)
	assert.True(gotBody.RemovePreviousTags)
}

func TestLinksService_BulkDelete(t *testing.T) {
	assert := require.New(t)

	var gotBody struct {
		LinkIDs []int `json:"linkIds"`
	}
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/v1/links", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Links deleted."}`))
	})

	assert.NoError(client.Links.BulkDelete(t.Context(), []int{1, 2, 3}))
	assert.Equal([]int{1, 2, 3}, gotBody.LinkIDs)
}

func TestLinksService_Archive(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPut, r.Method)
		assert.Equal("/api/v1/links/4/archive", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Link archived."}`))
	})

	assert.NoError(client.Links.Archive(t.Context(), 4))
}

func TestLinksService_DeleteArchive(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/v1/links/archive", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Archived files deleted."}`))
	})

	assert.NoError(client.Links.DeleteArchive(t.Context(), []int{4}))
}

func TestLinksService_Highlights(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/links/6/highlights", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":[{"id":1,"linkId":6,"text":"quote","color":"yellow","startOffset":0,"endOffset":5}]}`))
	})

	highlights, err := client.Links.Highlights(t.Context(), 6)
	assert.NoError(err)
	assert.Len(highlights, 1)
	assert.Equal("quote", highlights[0].Text)
}
