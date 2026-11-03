package linkwarden

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDashboardServiceGetV1(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	collection := createTestCollection(t, client)
	link := createTestLink(t, client, collection.ID)

	// Act
	links, err := client.Dashboard.GetV1(t.Context())

	// Assert
	assert.NoError(err)
	assert.Contains(linkIDs(links), link.ID)
}

func TestDashboardServiceGetV2(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	collection := createTestCollection(t, client)
	link := createTestLink(t, client, collection.ID)

	// Act
	data, err := client.Dashboard.GetV2(t.Context())

	// Assert
	assert.NoError(err)
	assert.Contains(linkIDs(data.Links), link.ID)
}

func TestDashboardServiceGetV2CollectionLinks(t *testing.T) {
	// Arrange
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

	// Act
	data, err := client.Dashboard.GetV2(t.Context())

	// Assert
	assert.NoError(err)
	assert.Equal([]int{link.ID}, linkIDs(data.CollectionLinks[strconv.Itoa(collection.ID)]))
}

func TestDashboardServiceUpdateV2Layout(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)

	// Act
	_, err := client.Dashboard.UpdateV2Layout(t.Context(), []DashboardLayoutItem{
		{Type: "STATS", Enabled: true, Order: 0},
		{Type: "RECENT_LINKS", Enabled: true, Order: 1},
		{Type: "PINNED_LINKS", Enabled: false, Order: 2},
	})

	// Assert
	assert.NoError(err)
}
