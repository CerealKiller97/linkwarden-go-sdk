package linkwarden

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDashboardService_GetV1(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	link := createTestLink(t, client, collection.ID)

	links, err := client.Dashboard.GetV1(t.Context())
	assert.NoError(err)
	assert.Contains(linkIDs(links), link.ID)
}

func TestDashboardService_GetV2(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	link := createTestLink(t, client, collection.ID)

	data, err := client.Dashboard.GetV2(t.Context())
	assert.NoError(err)
	assert.Contains(linkIDs(data.Links), link.ID)
}

func TestDashboardService_GetV2_CollectionLinks(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	link := createTestLink(t, client, collection.ID)

	// Collection links are only returned for collections shown on the dashboard.
	_, err := client.Dashboard.UpdateV2Layout(t.Context(), []DashboardLayoutItem{
		{Type: "RECENT_LINKS", Enabled: true, Order: 0},
		{Type: "COLLECTION", CollectionID: &collection.ID, Enabled: true, Order: 1},
	})
	assert.NoError(err)
	t.Cleanup(func() {
		_, _ = client.Dashboard.UpdateV2Layout(t.Context(), []DashboardLayoutItem{
			{Type: "RECENT_LINKS", Enabled: true, Order: 0},
		})
	})

	data, err := client.Dashboard.GetV2(t.Context())
	assert.NoError(err)
	assert.Equal([]int{link.ID}, linkIDs(data.CollectionLinks[strconv.Itoa(collection.ID)]))
}

func TestDashboardService_UpdateV2Layout(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)

	_, err := client.Dashboard.UpdateV2Layout(t.Context(), []DashboardLayoutItem{
		{Type: "STATS", Enabled: true, Order: 0},
		{Type: "RECENT_LINKS", Enabled: true, Order: 1},
		{Type: "PINNED_LINKS", Enabled: false, Order: 2},
	})
	assert.NoError(err)
}
