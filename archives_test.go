package linkwarden

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArchivesService_Get(t *testing.T) {
	assert := require.New(t)

	var gotQuery string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/archives/7", r.URL.Path)
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("%PDF-1.4 fake content"))
	})

	file, err := client.Archives.Get(t.Context(), 7, ArchiveFormatPDF, true)
	assert.NoError(err)
	assert.Equal("format=2&preview=true", gotQuery)
	assert.Equal("application/pdf", file.ContentType)
	assert.Equal("%PDF-1.4 fake content", string(file.Data))
}

func TestArchivesService_Get_NotFound(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"response":"Link not found."}`))
	})

	_, err := client.Archives.Get(t.Context(), 999, ArchiveFormatPNG, false)
	assert.Error(err)
	assert.True(IsNotFound(err))
}

func TestArchivesService_Upload(t *testing.T) {
	assert := require.New(t)

	var gotQuery, gotContentType string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/api/v1/archives/3", r.URL.Path)
		gotQuery = r.URL.RawQuery
		gotContentType = r.Header.Get("Content-Type")

		assert.NoError(r.ParseMultipartForm(1 << 20))
		file, _, err := r.FormFile("file")
		assert.NoError(err)
		defer file.Close()

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":3,"name":"Example","url":"https://example.com"}}`))
	})

	link, err := client.Archives.Upload(t.Context(), 3, ArchiveFormatPNG, false, "screenshot.png", strings.NewReader("fake-png-bytes"))
	assert.NoError(err)
	assert.Equal(3, link.ID)
	assert.Equal("format=0", gotQuery)
	assert.Contains(gotContentType, "multipart/form-data")
}

func TestArchivesService_UploadStandalone(t *testing.T) {
	assert := require.New(t)

	var gotURLField string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/api/v1/archives", r.URL.Path)
		assert.Equal("format=4&preview=true", r.URL.RawQuery)

		assert.NoError(r.ParseMultipartForm(1 << 20))
		gotURLField = r.FormValue("url")
		file, _, err := r.FormFile("file")
		assert.NoError(err)
		defer file.Close()

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":9,"name":"Standalone","url":"https://example.com/page"}}`))
	})

	link, err := client.Archives.UploadStandalone(t.Context(), ArchiveFormatMonolithHTML, true, "https://example.com/page", "page.html", strings.NewReader("<html></html>"))
	assert.NoError(err)
	assert.Equal(9, link.ID)
	assert.Equal("https://example.com/page", gotURLField)
}
