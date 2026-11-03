package linkwarden

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSearchServiceLinks(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	collection := createTestCollection(t, client)
	// One word with no separators, so the search can't match other links on a shared prefix.
	word := strings.ReplaceAll(uniqueName("searchable"), "-", "")
	link, err := client.Links.Create(t.Context(), CreateLinkRequest{
		Name:       word,
		URL:        uniqueURL(),
		Collection: Collection{ID: collection.ID},
	})
	assert.NoError(err)
	t.Cleanup(func() { _ = client.Links.Delete(t.Context(), link.ID) })

	// Act
	// New links reach the search index asynchronously, so retry until it finds one.
	var found []Link
	assert.Eventually(func() bool {
		links, err := client.Search.Links(t.Context(), SearchOptions{SearchQueryString: word})
		if err != nil {
			return false
		}
		found = links

		return len(links) > 0
	}, time.Minute, time.Second)

	// Assert
	assert.Equal([]int{link.ID}, linkIDs(found))
}

func TestSearchServiceLinksNoMatches(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)

	// Act
	links, err := client.Search.Links(t.Context(), SearchOptions{SearchQueryString: uniqueName("nothing")})

	// Assert
	assert.NoError(err)
	assert.Empty(links)
}
