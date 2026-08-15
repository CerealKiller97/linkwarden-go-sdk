package linkwarden

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRSSService_List(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/rss", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":[{"id":1,"name":"Blog","url":"https://blog.example.com/rss","collectionId":2}]}`))
	})

	subs, err := client.RSS.List(t.Context())
	assert.NoError(err)
	assert.Len(subs, 1)
	assert.Equal("Blog", subs[0].Name)
	assert.Equal(2, subs[0].CollectionID)
}

func TestRSSService_Create(t *testing.T) {
	assert := require.New(t)

	var gotBody CreateRSSRequest
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/api/v1/rss", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":3,"name":"News","url":"https://news.example.com/rss"}}`))
	})

	sub, err := client.RSS.Create(t.Context(), CreateRSSRequest{Name: "News", URL: "https://news.example.com/rss", CollectionID: 5})
	assert.NoError(err)
	assert.Equal(3, sub.ID)
	assert.Equal("News", gotBody.Name)
	assert.Equal(5, gotBody.CollectionID)
}

func TestRSSService_Delete(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/v1/rss/9", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"RSS subscription deleted."}`))
	})

	assert.NoError(client.RSS.Delete(t.Context(), 9))
}
