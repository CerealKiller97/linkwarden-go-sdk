package linkwarden

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const testFeedURL = "https://go.dev/blog/feed.atom"

func createTestSubscription(t *testing.T, c *Client, collectionID int) *RSSSubscription {
	t.Helper()

	sub, err := c.RSS.Create(t.Context(), CreateRSSRequest{
		Name:         uniqueName("feed"),
		URL:          testFeedURL,
		CollectionID: collectionID,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.RSS.Delete(t.Context(), sub.ID) })

	return sub
}

func TestRSSService_Create(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)

	sub := createTestSubscription(t, client, collection.ID)
	assert.NotZero(sub.ID)
	assert.Equal(testFeedURL, sub.URL)
	assert.Equal(collection.ID, sub.CollectionID)
}

func TestRSSService_List(t *testing.T) {
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

func TestRSSService_Delete(t *testing.T) {
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

func TestRSSService_Delete_NotFound(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)

	err := client.RSS.Delete(t.Context(), 999999999)
	assert.True(IsNotFound(err))
}
