package linkwarden

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagsService_List(t *testing.T) {
	assert := require.New(t)

	var gotQuery string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/tags", r.URL.Path)
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"Success","success":true,"data":{"nextCursor":42,"tags":[{"id":1,"name":"Technology"}]}}`))
	})

	tags, err := client.Tags.List(t.Context(), ListTagsOptions{Search: "tech", Cursor: 5})
	assert.NoError(err)
	assert.Equal("cursor=5&search=tech", gotQuery)
	assert.Len(tags, 1)
	assert.Equal("Technology", tags[0].Name)
}

func TestTagsService_Get(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/tags/1", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":1,"name":"Technology"}}`))
	})

	tag, err := client.Tags.Get(t.Context(), 1)
	assert.NoError(err)
	assert.Equal("Technology", tag.Name)
}

func TestTagsService_CreateOrUpdateBulk(t *testing.T) {
	assert := require.New(t)

	var gotBody struct {
		Tags []CreateOrUpdateTagInput `json:"tags"`
	}
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/api/v1/tags", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":[{"id":2,"name":"Reading"}]}`))
	})

	tags, err := client.Tags.CreateOrUpdateBulk(t.Context(), []CreateOrUpdateTagInput{{Label: "Reading"}})
	assert.NoError(err)
	assert.Len(tags, 1)
	assert.Equal("Reading", gotBody.Tags[0].Label)
}

func TestTagsService_Update(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPut, r.Method)
		assert.Equal("/api/v1/tags/3", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":3,"name":"Renamed"}}`))
	})

	tag, err := client.Tags.Update(t.Context(), 3, "Renamed")
	assert.NoError(err)
	assert.Equal("Renamed", tag.Name)
}

func TestTagsService_Delete(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/v1/tags/3", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Tag deleted."}`))
	})

	assert.NoError(client.Tags.Delete(t.Context(), 3))
}

func TestTagsService_BulkDelete(t *testing.T) {
	assert := require.New(t)

	var gotBody struct {
		TagIDs []int `json:"tagIds"`
	}
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/v1/tags", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Tags deleted."}`))
	})

	assert.NoError(client.Tags.BulkDelete(t.Context(), []int{1, 2}))
	assert.Equal([]int{1, 2}, gotBody.TagIDs)
}
