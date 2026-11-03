package linkwarden

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestHighlight(t *testing.T, c *Client, linkID int) *Highlight {
	t.Helper()
	assert := require.New(t)

	highlight, err := c.Highlights.Create(t.Context(), CreateHighlightRequest{
		Color:       "yellow",
		Comment:     "worth rereading",
		StartOffset: 0,
		EndOffset:   5,
		Text:        "Hello",
		LinkID:      linkID,
	})
	assert.NoError(err)
	t.Cleanup(func() { _ = c.Highlights.Delete(t.Context(), highlight.ID) })

	return highlight
}

func TestHighlightsServiceCreate(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	link := createTestLink(t, client, collection.ID)

	highlight := createTestHighlight(t, client, link.ID)
	assert.NotZero(highlight.ID)
	assert.Equal(link.ID, highlight.LinkID)
	assert.Equal("yellow", highlight.Color)
	assert.Equal("worth rereading", highlight.Comment)
	assert.Equal("Hello", highlight.Text)
	assert.Equal(5, highlight.EndOffset)
}

func TestHighlightsServiceDelete(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	link := createTestLink(t, client, collection.ID)
	highlight := createTestHighlight(t, client, link.ID)

	assert.NoError(client.Highlights.Delete(t.Context(), highlight.ID))

	highlights, err := client.Links.Highlights(t.Context(), link.ID)
	assert.NoError(err)
	assert.Empty(highlights)
}
