package linkwarden

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCollectionsServiceCreate(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	name := uniqueName("collection")

	collection, err := client.Collections.Create(t.Context(), CreateCollectionRequest{
		Name:        name,
		Description: "created by the SDK tests",
		Color:       "#f59e0b",
	})
	assert.NoError(err)
	t.Cleanup(func() { _ = client.Collections.Delete(t.Context(), collection.ID) })

	assert.NotZero(collection.ID)
	assert.Equal(name, collection.Name)
	assert.Equal("created by the SDK tests", collection.Description)
	assert.Equal("#f59e0b", collection.Color)
}

func TestCollectionsServiceCreateWithParent(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	parent := createTestCollection(t, client)

	child, err := client.Collections.Create(t.Context(), CreateCollectionRequest{
		Name:     uniqueName("child"),
		ParentID: parent.ID,
	})
	assert.NoError(err)
	t.Cleanup(func() { _ = client.Collections.Delete(t.Context(), child.ID) })

	assert.NotNil(child.ParentID)
	assert.Equal(parent.ID, *child.ParentID)
}

func TestCollectionsServiceList(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)

	collections, err := client.Collections.List(t.Context())
	assert.NoError(err)

	var found *Collection
	for i := range collections {
		if collections[i].ID == collection.ID {
			found = &collections[i]
		}
	}
	assert.NotNil(found, "new collection should be listed")
	assert.Equal(collection.Name, found.Name)
}

func TestCollectionsServiceGet(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)

	got, err := client.Collections.Get(t.Context(), collection.ID)
	assert.NoError(err)
	assert.Equal(collection.ID, got.ID)
	assert.Equal(collection.Name, got.Name)
}

func TestCollectionsServiceGetNotFound(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)

	_, err := client.Collections.Get(t.Context(), 999999999)
	assert.True(IsNotFound(err))
}

func TestCollectionsServiceUpdate(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	renamed := uniqueName("renamed")

	updated, err := client.Collections.Update(t.Context(), collection.ID, UpdateCollectionRequest{
		Name:        renamed,
		Description: "updated",
		Color:       "#10b981",
	})
	assert.NoError(err)
	assert.Equal(renamed, updated.Name)

	got, err := client.Collections.Get(t.Context(), collection.ID)
	assert.NoError(err)
	assert.Equal(renamed, got.Name)
	assert.Equal("updated", got.Description)
	assert.Equal("#10b981", got.Color)
}

func TestCollectionsServiceUpdateKeepsUnsetFields(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection, err := client.Collections.Create(t.Context(), CreateCollectionRequest{
		Name:        uniqueName("collection"),
		Description: "kept",
		Color:       "#f59e0b",
	})
	assert.NoError(err)
	t.Cleanup(func() { _ = client.Collections.Delete(t.Context(), collection.ID) })

	_, err = client.Collections.Update(t.Context(), collection.ID, UpdateCollectionRequest{Color: "#10b981"})
	assert.NoError(err)

	got, err := client.Collections.Get(t.Context(), collection.ID)
	assert.NoError(err)
	assert.Equal(collection.Name, got.Name)
	assert.Equal("kept", got.Description)
	assert.Equal("#10b981", got.Color)
}

func TestCollectionsServiceUpdateIsPublic(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)

	_, err := client.Collections.Update(t.Context(), collection.ID, UpdateCollectionRequest{IsPublic: new(true)})
	assert.NoError(err)
	got, err := client.Collections.Get(t.Context(), collection.ID)
	assert.NoError(err)
	assert.True(got.IsPublic)

	// Leaving IsPublic nil keeps it public…
	_, err = client.Collections.Update(t.Context(), collection.ID, UpdateCollectionRequest{Color: "#10b981"})
	assert.NoError(err)
	got, err = client.Collections.Get(t.Context(), collection.ID)
	assert.NoError(err)
	assert.True(got.IsPublic)

	// …and an explicit false makes it private again.
	_, err = client.Collections.Update(t.Context(), collection.ID, UpdateCollectionRequest{IsPublic: new(false)})
	assert.NoError(err)
	got, err = client.Collections.Get(t.Context(), collection.ID)
	assert.NoError(err)
	assert.False(got.IsPublic)
}

func TestCollectionsServiceDelete(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection, err := client.Collections.Create(t.Context(), CreateCollectionRequest{Name: uniqueName("collection")})
	assert.NoError(err)

	assert.NoError(client.Collections.Delete(t.Context(), collection.ID))

	collections, err := client.Collections.List(t.Context())
	assert.NoError(err)
	for _, c := range collections {
		assert.NotEqual(collection.ID, c.ID, "deleted collection should not be listed")
	}
}

func TestCollectionsServiceDeleteInaccessible(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)

	err := client.Collections.Delete(t.Context(), 999999999)
	assert.Error(err)

	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
}
