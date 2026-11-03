package linkwarden

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagsServiceCreateOrUpdateBulk(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	first, second := uniqueName("tag"), uniqueName("tag")

	// Act
	tags, err := client.Tags.CreateOrUpdateBulk(t.Context(), []CreateOrUpdateTagInput{{Label: first}, {Label: second}})

	// Assert
	assert.NoError(err)
	for _, tag := range tags {
		t.Cleanup(func() { _ = client.Tags.Delete(t.Context(), tag.ID) })
	}
	assert.Len(tags, 2)
	assert.ElementsMatch([]string{first, second}, []string{tags[0].Name, tags[1].Name})
}

func TestTagsServiceCreateOrUpdateBulkArchiveSettings(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	label := uniqueName("tag")

	// Act
	tags, err := client.Tags.CreateOrUpdateBulk(t.Context(), []CreateOrUpdateTagInput{{
		Label:               label,
		ArchiveAsScreenshot: new(false),
		ArchiveAsPDF:        new(true),
	}})

	// Assert
	assert.NoError(err)
	assert.Len(tags, 1)
	t.Cleanup(func() { _ = client.Tags.Delete(t.Context(), tags[0].ID) })
	assert.Equal(label, tags[0].Name)
}

func TestTagsServiceList(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	tag := createTestTag(t, client)

	// Act
	tags, err := client.Tags.List(t.Context(), ListTagsOptions{Search: tag.Name})

	// Assert
	assert.NoError(err)
	assert.Len(tags, 1)
	assert.Equal(tag.ID, tags[0].ID)
	assert.Equal(tag.Name, tags[0].Name)
}

func TestTagsServiceGet(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	tag := createTestTag(t, client)

	// Act
	got, err := client.Tags.Get(t.Context(), tag.ID)

	// Assert
	assert.NoError(err)
	assert.Equal(tag.ID, got.ID)
	assert.Equal(tag.Name, got.Name)
}

func TestTagsServiceGetNotFound(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)

	// Act
	_, err := client.Tags.Get(t.Context(), 999999999)

	// Assert
	assert.Error(err)
	assert.True(IsNotFound(err))
}

func TestTagsServiceUpdate(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	tag := createTestTag(t, client)
	renamed := uniqueName("renamed")

	// Act
	updated, err := client.Tags.Update(t.Context(), tag.ID, renamed)

	// Assert
	assert.NoError(err)
	assert.Equal(renamed, updated.Name)
	got, err := client.Tags.Get(t.Context(), tag.ID)
	assert.NoError(err)
	assert.Equal(renamed, got.Name)
}

func TestTagsServiceDelete(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	tag := createTestTag(t, client)

	// Act
	err := client.Tags.Delete(t.Context(), tag.ID)

	// Assert
	assert.NoError(err)
	_, err = client.Tags.Get(t.Context(), tag.ID)
	assert.True(IsNotFound(err))
}

func TestTagsServiceBulkDelete(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	first, second, kept := createTestTag(t, client), createTestTag(t, client), createTestTag(t, client)

	// Act
	err := client.Tags.BulkDelete(t.Context(), []int{first.ID, second.ID})

	// Assert
	assert.NoError(err)
	for _, id := range []int{first.ID, second.ID} {
		_, err := client.Tags.Get(t.Context(), id)
		assert.True(IsNotFound(err))
	}
	_, err = client.Tags.Get(t.Context(), kept.ID)
	assert.NoError(err)
}
