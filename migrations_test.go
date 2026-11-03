package linkwarden

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrationsService_Export(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/migration", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"username":"bob","collections":[]}`))
	})

	data, err := client.Migrations.Export(t.Context())
	assert.NoError(err)
	assert.JSONEq(`{"username":"bob","collections":[]}`, string(data))
}

func TestMigrationsService_Import(t *testing.T) {
	assert := require.New(t)

	var gotBody struct {
		Format int    `json:"format"`
		Data   string `json:"data"`
	}
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/api/v1/migration", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Import started."}`))
	})

	msg, err := client.Migrations.Import(t.Context(), ImportFormatBookmarksHTML, "<html>bookmarks</html>")
	assert.NoError(err)
	assert.Equal("Import started.", msg)
	assert.Equal(1, gotBody.Format)
	assert.Equal("<html>bookmarks</html>", gotBody.Data)
}
