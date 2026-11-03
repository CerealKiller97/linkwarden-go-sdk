package linkwarden

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/stretchr/testify/require"
)

// testPNG returns a small, valid PNG image to upload as an archive.
func testPNG(t *testing.T) []byte {
	t.Helper()
	assert := require.New(t)

	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for x := range 4 {
		for y := range 4 {
			img.Set(x, y, color.RGBA{R: 245, G: 158, B: 11, A: 255})
		}
	}

	var buf bytes.Buffer
	assert.NoError(png.Encode(&buf, img))

	return buf.Bytes()
}

func TestArchivesServiceUpload(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	collection := createTestCollection(t, client)
	link := createTestLink(t, client, collection.ID)

	// Act
	updated, err := client.Archives.Upload(t.Context(), link.ID, ArchiveFormatPNG, false, "screenshot.png", bytes.NewReader(testPNG(t)))

	// Assert
	assert.NoError(err)
	assert.Equal(link.ID, updated.ID)
	assert.NotEmpty(updated.Image)
}

func TestArchivesServiceGet(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	collection := createTestCollection(t, client)
	link := createTestLink(t, client, collection.ID)
	data := testPNG(t)
	_, err := client.Archives.Upload(t.Context(), link.ID, ArchiveFormatPNG, false, "screenshot.png", bytes.NewReader(data))
	assert.NoError(err)

	// Act
	file, err := client.Archives.Get(t.Context(), link.ID, ArchiveFormatPNG, false)

	// Assert
	assert.NoError(err)
	assert.Equal("image/png", file.ContentType)
	assert.Equal(data, file.Data)
}

func TestArchivesServiceGetInaccessible(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)

	// Act
	_, err := client.Archives.Get(t.Context(), 999999999, ArchiveFormatPNG, false)

	// Assert
	assert.Error(err)
	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
}

func TestArchivesServiceUploadStandalone(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	url := uniqueURL()

	// Act
	link, err := client.Archives.UploadStandalone(t.Context(), ArchiveFormatPNG, false, url, "screenshot.png", bytes.NewReader(testPNG(t)))

	// Assert
	assert.NoError(err)
	t.Cleanup(func() { _ = client.Links.Delete(t.Context(), link.ID) })
	assert.NotZero(link.ID)
	assert.Equal(url, link.URL)
}

// A file without a recognisable extension is uploaded with its sniffed type, since the
// API rejects application/octet-stream parts.
func TestArchivesServiceUploadSniffsContentType(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	link := createTestLink(t, client, createTestCollection(t, client).ID)
	data := testPNG(t)

	// Act
	_, err := client.Archives.Upload(t.Context(), link.ID, ArchiveFormatPNG, false, "screenshot", bytes.NewReader(data))

	// Assert
	assert.NoError(err)
	file, err := client.Archives.Get(t.Context(), link.ID, ArchiveFormatPNG, false)
	assert.NoError(err)
	assert.Equal(data, file.Data)
}
