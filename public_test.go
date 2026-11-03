package linkwarden

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// createPublicTestCollection creates a public collection holding one tagged link.
func createPublicTestCollection(t *testing.T, c *Client) (*Collection, *Link) {
	t.Helper()
	assert := require.New(t)

	collection := createTestCollection(t, c)
	_, err := c.Collections.Update(t.Context(), collection.ID, UpdateCollectionRequest{IsPublic: new(true)})
	assert.NoError(err)

	link, err := c.Links.Create(t.Context(), CreateLinkRequest{
		Name:       uniqueName("link"),
		URL:        uniqueURL(),
		Collection: Collection{ID: collection.ID},
		Tags:       []Tag{{Name: uniqueName("tag")}},
	})
	assert.NoError(err)
	t.Cleanup(func() { _ = c.Links.Delete(t.Context(), link.ID) })

	return collection, link
}

func TestPublicServiceCollection(t *testing.T) {
	// Arrange
	assert := require.New(t)
	collection, _ := createPublicTestCollection(t, newTestClient(t))
	anon := newAnonymousClient(t)

	// Act
	got, err := anon.Public.Collection(t.Context(), collection.ID)

	// Assert
	assert.NoError(err)
	assert.Equal(collection.ID, got.ID)
	assert.Equal(collection.Name, got.Name)
}

func TestPublicServiceCollectionPrivate(t *testing.T) {
	// Arrange
	assert := require.New(t)
	collection := createTestCollection(t, newTestClient(t))
	anon := newAnonymousClient(t)

	// Act
	_, err := anon.Public.Collection(t.Context(), collection.ID)

	// Assert
	assert.Error(err)
}

func TestPublicServiceCollectionLinks(t *testing.T) {
	// Arrange
	assert := require.New(t)
	collection, link := createPublicTestCollection(t, newTestClient(t))
	anon := newAnonymousClient(t)

	// Act
	links, err := anon.Public.CollectionLinks(t.Context(), collection.ID, ListPublicLinksOptions{})

	// Assert
	assert.NoError(err)
	assert.Equal([]int{link.ID}, linkIDs(links))
}

func TestPublicServiceCollectionTags(t *testing.T) {
	// Arrange
	assert := require.New(t)
	collection, link := createPublicTestCollection(t, newTestClient(t))
	anon := newAnonymousClient(t)

	// Act
	tags, err := anon.Public.CollectionTags(t.Context(), collection.ID)

	// Assert
	assert.NoError(err)
	var names []string
	for _, tag := range tags {
		names = append(names, tag.Name)
	}
	assert.Contains(names, link.Tags[0].Name)
}

func TestPublicServiceLink(t *testing.T) {
	// Arrange
	assert := require.New(t)
	_, link := createPublicTestCollection(t, newTestClient(t))
	anon := newAnonymousClient(t)

	// Act
	got, err := anon.Public.Link(t.Context(), link.ID)

	// Assert
	assert.NoError(err)
	assert.Equal(link.ID, got.ID)
	assert.Equal(link.URL, got.URL)
}

func TestPublicServiceUser(t *testing.T) {
	// Arrange
	assert := require.New(t)
	me, err := newTestClient(t).Users.Me(t.Context())
	assert.NoError(err)
	anon := newAnonymousClient(t)

	// Act
	got, err := anon.Public.User(t.Context(), me.ID)

	// Assert
	assert.NoError(err)
	assert.Equal(me.ID, got.ID)
	assert.Equal(me.Username, got.Username)
}
