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
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)

	sub := createTestSubscription(t, client, collection.ID)
	assert.NotZero(sub.ID)
	assert.Equal(testFeedURL, sub.URL)
	assert.Equal(collection.ID, sub.CollectionID)
}

func TestRSSServiceList(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	sub := createTestSubscription(t, client, collection.ID)

	subs, err := client.RSS.List(t.Context())
	assert.NoError(err)

	var ids []int
	for _, s := range subs {
		ids = append(ids, s.ID)
	}
	assert.Contains(ids, sub.ID)
}

func TestRSSServiceDelete(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	sub := createTestSubscription(t, client, collection.ID)

	assert.NoError(client.RSS.Delete(t.Context(), sub.ID))

	subs, err := client.RSS.List(t.Context())
	assert.NoError(err)
	for _, s := range subs {
		assert.NotEqual(sub.ID, s.ID, "deleted subscription should not be listed")
	}
}

func TestRSSServiceDeleteNotFound(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)

	err := client.RSS.Delete(t.Context(), 999999999)
	assert.True(IsNotFound(err))
}
