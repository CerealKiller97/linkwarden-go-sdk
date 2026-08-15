package linkwarden

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHighlightsService_Create(t *testing.T) {
	assert := require.New(t)

	var gotBody CreateHighlightRequest
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/api/v1/highlights", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":4,"userId":1,"linkId":9,"color":"yellow","text":"quoted text","startOffset":0,"endOffset":11}}`))
	})

	highlight, err := client.Highlights.Create(t.Context(), CreateHighlightRequest{
		Color:       "yellow",
		Text:        "quoted text",
		StartOffset: 0,
		EndOffset:   11,
		LinkID:      9,
	})
	assert.NoError(err)
	assert.Equal(4, highlight.ID)
	assert.Equal(9, gotBody.LinkID)
	assert.Equal("yellow", gotBody.Color)
}

func TestHighlightsService_Delete(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/v1/highlights/4", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":4}`))
	})

	assert.NoError(client.Highlights.Delete(t.Context(), 4))
}
