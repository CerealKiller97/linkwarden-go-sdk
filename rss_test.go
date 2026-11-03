package linkwarden

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const testFeedURL = "https://go.dev/blog/feed.atom"

func createTestSubscription(t *testing.T, c *Client, collectionID int) *RSSSubscription {
	t.Helper()
	assert := require.New(t)

	sub, err := c.RSS.Create(t.Context(), CreateRSSRequest{
		Name:         uniqueName("feed"),
		URL:          testFeedURL,
		CollectionID: collectionID,
	})
	assert.NoError(err)
	t.Cleanup(func() { _ = c.RSS.Delete(t.Context(), sub.ID) })

	return sub
}

func TestRSSServiceCreate(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	collection := createTestCollection(t, client)

	// Act
	sub, err := client.RSS.Create(t.Context(), CreateRSSRequest{
		Name:         uniqueName("feed"),
		URL:          testFeedURL,
		CollectionID: collection.ID,
	})

	// Assert
	assert.NoError(err)
	t.Cleanup(func() { _ = client.RSS.Delete(t.Context(), sub.ID) })
	assert.NotZero(sub.ID)
	assert.Equal(testFeedURL, sub.URL)
	assert.Equal(collection.ID, sub.CollectionID)
}

func TestRSSServiceList(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	collection := createTestCollection(t, client)
	sub := createTestSubscription(t, client, collection.ID)

	// Act
	subs, err := client.RSS.List(t.Context())

	// Assert
	assert.NoError(err)
	var ids []int
	for _, s := range subs {
		ids = append(ids, s.ID)
	}
	assert.Contains(ids, sub.ID)
}

func TestRSSServiceDelete(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	collection := createTestCollection(t, client)
	sub := createTestSubscription(t, client, collection.ID)

	// Act
	err := client.RSS.Delete(t.Context(), sub.ID)

	// Assert
	assert.NoError(err)
	subs, err := client.RSS.List(t.Context())
	assert.NoError(err)
	for _, s := range subs {
		assert.NotEqual(sub.ID, s.ID, "deleted subscription should not be listed")
	}
}

func TestRSSServiceDeleteNotFound(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)

	// Act
	err := client.RSS.Delete(t.Context(), 999999999)

	// Assert
	assert.True(IsNotFound(err))
}
