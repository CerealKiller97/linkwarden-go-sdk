package linkwarden

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCollectionsService_List(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/collections", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":[{"id":1,"name":"Unorganized"}]}`))
	})

	collections, err := client.Collections.List(t.Context())
	assert.NoError(err)
	assert.Len(collections, 1)
	assert.Equal("Unorganized", collections[0].Name)
}

func TestCollectionsService_Get(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/collections/2", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":2,"name":"Reading List"}}`))
	})

	collection, err := client.Collections.Get(t.Context(), 2)
	assert.NoError(err)
	assert.Equal(2, collection.ID)
	assert.Equal("Reading List", collection.Name)
}

func TestCollectionsService_Create(t *testing.T) {
	assert := require.New(t)

	var gotBody CreateCollectionRequest
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		// The SDK deliberately posts to the un-suffixed path, matching the
		// proven-working live-instance behavior rather than the OpenAPI spec's
		// (likely mis-documented) {id}-suffixed create path.
		assert.Equal("/api/v1/collections", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":10,"name":"New Collection"}}`))
	})

	collection, err := client.Collections.Create(t.Context(), CreateCollectionRequest{Name: "New Collection"})
	assert.NoError(err)
	assert.Equal(10, collection.ID)
	assert.Equal("New Collection", gotBody.Name)
}

func TestCollectionsService_Update(t *testing.T) {
	assert := require.New(t)

	var gotBody UpdateCollectionRequest
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPut, r.Method)
		assert.Equal("/api/v1/collections/3", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":3,"name":"Renamed"}}`))
	})

	collection, err := client.Collections.Update(t.Context(), 3, UpdateCollectionRequest{Name: "Renamed", IsPublic: true})
	assert.NoError(err)
	assert.Equal("Renamed", collection.Name)
	assert.True(gotBody.IsPublic)
}

func TestCollectionsService_Delete(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/v1/collections/4", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Collection deleted."}`))
	})

	assert.NoError(client.Collections.Delete(t.Context(), 4))
}

func TestCollectionsService_Delete_Conflict(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"message":"Cannot delete collection with links."}`))
	})

	err := client.Collections.Delete(t.Context(), 4)
	assert.Error(err)
	assert.True(IsConflict(err))
}
