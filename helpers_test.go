package linkwarden

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// The tests run against a real Linkwarden instance, set through LINKWARDEN_URL and
// LINKWARDEN_TOKEN. `scripts/dev-linkwarden.sh up` starts a throwaway one locally, and CI
// starts one per run. Without them, tests are skipped locally but fail in CI, so a broken
// workflow can't pass by skipping everything.
//
// Every test creates its own uniquely named data and deletes it afterwards, so tests can
// share one instance and run in any order.

func newTestClient(t *testing.T) *Client {
	t.Helper()

	baseURL, token := os.Getenv("LINKWARDEN_URL"), os.Getenv("LINKWARDEN_TOKEN")
	if baseURL == "" || token == "" {
		if os.Getenv("CI") != "" {
			t.Fatal("LINKWARDEN_URL and LINKWARDEN_TOKEN must be set in CI")
		}
		t.Skip("LINKWARDEN_URL / LINKWARDEN_TOKEN not set; start an instance with scripts/dev-linkwarden.sh up")
	}

	return New(baseURL, token, WithLogger(slog.New(slog.DiscardHandler)))
}

// uniqueName returns prefix followed by a random suffix, for names that must not clash
// with data from other tests or earlier runs.
func uniqueName(prefix string) string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)

	return prefix + "-" + hex.EncodeToString(b)
}

// uniqueURL returns a URL that no other link on the instance uses.
func uniqueURL() string {
	return "https://example.com/" + uniqueName("e2e")
}

func createTestCollection(t *testing.T, c *Client) *Collection {
	t.Helper()
	assert := require.New(t)

	collection, err := c.Collections.Create(t.Context(), CreateCollectionRequest{Name: uniqueName("collection")})
	assert.NoError(err)
	t.Cleanup(func() { _ = c.Collections.Delete(t.Context(), collection.ID) })

	return collection
}

func createTestLink(t *testing.T, c *Client, collectionID int) *Link {
	t.Helper()
	assert := require.New(t)

	link, err := c.Links.Create(t.Context(), CreateLinkRequest{
		Name:       uniqueName("link"),
		URL:        uniqueURL(),
		Collection: Collection{ID: collectionID},
	})
	assert.NoError(err)
	t.Cleanup(func() { _ = c.Links.Delete(t.Context(), link.ID) })

	return link
}

func createTestTag(t *testing.T, c *Client) *TagInfo {
	t.Helper()
	assert := require.New(t)

	tags, err := c.Tags.CreateOrUpdateBulk(t.Context(), []CreateOrUpdateTagInput{{Label: uniqueName("tag")}})
	assert.NoError(err)
	assert.Len(tags, 1)
	t.Cleanup(func() { _ = c.Tags.Delete(t.Context(), tags[0].ID) })

	return &tags[0]
}

func linkIDs(links []Link) []int {
	ids := make([]int, len(links))
	for i, l := range links {
		ids[i] = l.ID
	}

	return ids
}
