package linkwarden

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinksServiceCreate(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	url := uniqueURL()
	tagName := uniqueName("tag")

	link, err := client.Links.Create(t.Context(), CreateLinkRequest{
		Name:        "New",
		URL:         url,
		Description: "created by the SDK tests",
		Tags:        []Tag{{Name: tagName}},
		Collection:  Collection{ID: collection.ID},
	})
	assert.NoError(err)
	t.Cleanup(func() { _ = client.Links.Delete(t.Context(), link.ID) })
	t.Cleanup(func() { deleteTagByName(t, client, tagName) })

	assert.NotZero(link.ID)
	assert.Equal("New", link.Name)
	assert.Equal(url, link.URL)
	assert.Equal(collection.ID, link.CollectionID)
	assert.Len(link.Tags, 1)
	assert.Equal(tagName, link.Tags[0].Name)
}

func TestLinksServiceGet(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	link := createTestLink(t, client, collection.ID)

	got, err := client.Links.Get(t.Context(), link.ID)
	assert.NoError(err)
	assert.Equal(link.ID, got.ID)
	assert.Equal(link.URL, got.URL)
	assert.Equal(collection.ID, got.Collection.ID)
}

func TestLinksServiceGetInaccessible(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)

	_, err := client.Links.Get(t.Context(), 999999999)
	assert.Error(err)

	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
	assert.NotEmpty(apiErr.Message)
}

func TestLinksServiceListFiltersByCollection(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	other := createTestCollection(t, client)
	first := createTestLink(t, client, collection.ID)
	second := createTestLink(t, client, collection.ID)
	createTestLink(t, client, other.ID)

	links, err := client.Links.List(t.Context(), ListLinksOptions{CollectionID: collection.ID})
	assert.NoError(err)
	assert.ElementsMatch([]int{first.ID, second.ID}, linkIDs(links))
}

func TestLinksServiceListFiltersByTag(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	tag := createTestTag(t, client)

	tagged, err := client.Links.Create(t.Context(), CreateLinkRequest{
		Name:       uniqueName("tagged"),
		URL:        uniqueURL(),
		Tags:       []Tag{{Name: tag.Name}},
		Collection: Collection{ID: collection.ID},
	})
	assert.NoError(err)
	t.Cleanup(func() { _ = client.Links.Delete(t.Context(), tagged.ID) })
	createTestLink(t, client, collection.ID)

	links, err := client.Links.List(t.Context(), ListLinksOptions{TagID: tag.ID})
	assert.NoError(err)
	assert.Equal([]int{tagged.ID}, linkIDs(links))
}

func TestLinksServiceUpdate(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	link := createTestLink(t, client, collection.ID)

	updated, err := client.Links.Update(t.Context(), link.ID, UpdateLinkRequest{
		Name:        "Renamed",
		Description: "updated",
	})
	assert.NoError(err)
	assert.Equal(link.ID, updated.ID)

	got, err := client.Links.Get(t.Context(), link.ID)
	assert.NoError(err)
	assert.Equal("Renamed", got.Name)
	assert.Equal("updated", got.Description)
}

func TestLinksServiceUpdateKeepsTagsAndCollection(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	tag := createTestTag(t, client)
	link, err := client.Links.Create(t.Context(), CreateLinkRequest{
		Name:       uniqueName("link"),
		URL:        uniqueURL(),
		Tags:       []Tag{{Name: tag.Name}},
		Collection: Collection{ID: collection.ID},
	})
	assert.NoError(err)
	t.Cleanup(func() { _ = client.Links.Delete(t.Context(), link.ID) })

	_, err = client.Links.Update(t.Context(), link.ID, UpdateLinkRequest{Name: "Renamed"})
	assert.NoError(err)

	got, err := client.Links.Get(t.Context(), link.ID)
	assert.NoError(err)
	assert.Equal("Renamed", got.Name)
	assert.Equal(link.URL, got.URL)
	assert.Equal(collection.ID, got.Collection.ID)
	assert.Len(got.Tags, 1)
	assert.Equal(tag.Name, got.Tags[0].Name)
}

func TestLinksServiceUpdateMovesToCollection(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	from := createTestCollection(t, client)
	to := createTestCollection(t, client)
	link := createTestLink(t, client, from.ID)

	_, err := client.Links.Update(t.Context(), link.ID, UpdateLinkRequest{Collection: Collection{ID: to.ID}})
	assert.NoError(err)

	got, err := client.Links.Get(t.Context(), link.ID)
	assert.NoError(err)
	assert.Equal(to.ID, got.Collection.ID)
}

func TestLinksServiceDelete(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	link := createTestLink(t, client, collection.ID)

	assert.NoError(client.Links.Delete(t.Context(), link.ID))

	links, err := client.Links.List(t.Context(), ListLinksOptions{CollectionID: collection.ID})
	assert.NoError(err)
	assert.Empty(links)
}

func TestLinksServiceBulkUpdate(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	from := createTestCollection(t, client)
	to := createTestCollection(t, client)
	first := createTestLink(t, client, from.ID)
	second := createTestLink(t, client, from.ID)
	tag := createTestTag(t, client)

	err := client.Links.BulkUpdate(t.Context(), BulkUpdateLinksRequest{
		Links:              []Link{*first, *second},
		RemovePreviousTags: true,
		NewData: BulkUpdateLinksNewData{
			CollectionID: to.ID,
			Tags:         []Tag{{Name: tag.Name}},
		},
	})
	assert.NoError(err)

	moved, err := client.Links.List(t.Context(), ListLinksOptions{CollectionID: to.ID})
	assert.NoError(err)
	assert.ElementsMatch([]int{first.ID, second.ID}, linkIDs(moved))
	for _, l := range moved {
		assert.Len(l.Tags, 1)
		assert.Equal(tag.Name, l.Tags[0].Name)
	}
}

func TestLinksServiceBulkDelete(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	first := createTestLink(t, client, collection.ID)
	second := createTestLink(t, client, collection.ID)
	kept := createTestLink(t, client, collection.ID)

	assert.NoError(client.Links.BulkDelete(t.Context(), []int{first.ID, second.ID}))

	links, err := client.Links.List(t.Context(), ListLinksOptions{CollectionID: collection.ID})
	assert.NoError(err)
	assert.Equal([]int{kept.ID}, linkIDs(links))
}

func TestLinksServiceArchive(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	link := createTestLink(t, client, collection.ID)

	assert.NoError(client.Links.Archive(t.Context(), link.ID))
}

func TestLinksServiceDeleteArchive(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	link := createTestLink(t, client, collection.ID)

	assert.NoError(client.Links.DeleteArchive(t.Context(), []int{link.ID}))
}

func TestLinksServiceHighlights(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	link := createTestLink(t, client, collection.ID)

	highlights, err := client.Links.Highlights(t.Context(), link.ID)
	assert.NoError(err)
	assert.Empty(highlights)

	created := createTestHighlight(t, client, link.ID)

	highlights, err = client.Links.Highlights(t.Context(), link.ID)
	assert.NoError(err)
	assert.Len(highlights, 1)
	assert.Equal(created.ID, highlights[0].ID)
}

func deleteTagByName(t *testing.T, c *Client, name string) {
	t.Helper()

	tags, err := c.Tags.List(t.Context(), ListTagsOptions{Search: name})
	if err != nil {
		return
	}
	for _, tag := range tags {
		if tag.Name == name {
			_ = c.Tags.Delete(t.Context(), tag.ID)
		}
	}
}
