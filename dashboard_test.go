package linkwarden

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDashboardService_GetV1(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/dashboard", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":[{"id":1,"name":"Go","url":"https://go.dev"}]}`))
	})

	links, err := client.Dashboard.GetV1(t.Context())
	assert.NoError(err)
	assert.Len(links, 1)
	assert.Equal("Go", links[0].Name)
}

func TestDashboardService_GetV2(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v2/dashboard", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true,"message":"Success","data":{"links":[{"id":1,"name":"Go","url":"https://go.dev"}],"numberOfPinnedLinks":2,"numberOfTags":5}}`))
	})

	data, err := client.Dashboard.GetV2(t.Context())
	assert.NoError(err)
	assert.Len(data.Links, 1)
	assert.Equal(2, data.NumberOfPinnedLinks)
	assert.Equal(5, data.NumberOfTags)
}

func TestDashboardService_UpdateV2Layout(t *testing.T) {
	assert := require.New(t)

	var gotBody []DashboardLayoutItem
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPut, r.Method)
		assert.Equal("/api/v2/dashboard", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true,"message":"Success","data":{"links":[],"numberOfPinnedLinks":0,"numberOfTags":0}}`))
	})

	collectionID := 3
	layout := []DashboardLayoutItem{
		{Type: "STATS", Enabled: true, Order: 0},
		{Type: "COLLECTION", CollectionID: &collectionID, Enabled: true, Order: 1},
	}

	data, err := client.Dashboard.UpdateV2Layout(t.Context(), layout)
	assert.NoError(err)
	assert.NotNil(data)
	assert.Len(gotBody, 2)
	assert.Equal("COLLECTION", gotBody[1].Type)
	assert.Equal(3, *gotBody[1].CollectionID)
}
