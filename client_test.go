package linkwarden

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTrimsTrailingSlash(t *testing.T) {
	// Arrange
	assert := require.New(t)

	// Act
	c := New("https://example.com/", "tok")

	// Assert
	assert.Equal("https://example.com", c.baseURL)
}

func TestNewDiscardsLogsByDefault(t *testing.T) {
	// Arrange
	assert := require.New(t)

	// Act
	c := New("https://example.com", "tok")

	// Assert
	assert.False(c.logger.Enabled(context.Background(), slog.LevelError))
}

func TestWithLoggerOverridesDefault(t *testing.T) {
	// Arrange
	assert := require.New(t)
	buf := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(buf, nil))

	// Act
	c := New("https://example.com", "tok", WithLogger(logger))
	c.logger.Error("boom")

	// Assert
	assert.Same(logger, c.logger)
	assert.Contains(buf.String(), "boom")
}

func TestWithLoggerIgnoresNil(t *testing.T) {
	// Arrange
	assert := require.New(t)

	// Act
	c := New("https://example.com", "tok", WithLogger(nil))

	// Assert
	assert.NotNil(c.logger)
}

func TestRequestAuthenticatesWithToken(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)

	// Act
	_, err := client.Collections.List(t.Context())

	// Assert
	assert.NoError(err)
}

func TestRequestInvalidTokenReturnsAPIError(t *testing.T) {
	// Arrange
	assert := require.New(t)
	newTestClient(t)
	client := New(os.Getenv("LINKWARDEN_URL"), "invalid-token")

	// Act
	_, err := client.Collections.List(t.Context())

	// Assert
	assert.Error(err)
	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
	assert.Equal(http.StatusUnauthorized, apiErr.StatusCode)
	assert.NotEmpty(apiErr.Message)
	assert.NotEmpty(apiErr.RawBody)
}

func TestNewAPIErrorReadsMessageField(t *testing.T) {
	// Arrange
	assert := require.New(t)

	// Act
	err := newAPIError(http.StatusBadRequest, []byte(`{"message":"Bad input."}`))

	// Assert
	assert.Equal("Bad input.", err.Message)
	assert.Equal("linkwarden: 400 Bad input.", err.Error())
}

func TestNewAPIErrorReadsResponseField(t *testing.T) {
	// Arrange
	assert := require.New(t)

	// Act
	err := newAPIError(http.StatusNotFound, []byte(`{"response":"Tag not found."}`))

	// Assert
	assert.Equal("Tag not found.", err.Message)
	assert.True(IsNotFound(err))
	assert.False(IsConflict(err))
}

func TestNewAPIErrorFallsBackToStatusText(t *testing.T) {
	// Arrange
	assert := require.New(t)

	// Act
	err := newAPIError(http.StatusInternalServerError, []byte(`Internal Server Error`))

	// Assert
	assert.Equal(http.StatusText(http.StatusInternalServerError), err.Message)
}

func TestDecodeResponseEnvelope(t *testing.T) {
	// Arrange
	assert := require.New(t)

	// Act
	got, err := decodeResponseEnvelope[string]([]byte(`{"response":"hello"}`))

	// Assert
	assert.NoError(err)
	assert.Equal("hello", got)
}

func TestDecodeDataEnvelope(t *testing.T) {
	// Arrange
	assert := require.New(t)

	// Act
	got, err := decodeDataEnvelope[string]([]byte(`{"data":"hello","message":"Success","success":true}`))

	// Assert
	assert.NoError(err)
	assert.Equal("hello", got)
}

func TestFileContentType(t *testing.T) {
	png := []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR")

	tests := []struct {
		name     string
		filename string
		head     []byte
		want     string
	}{
		{"from extension", "page.pdf", nil, "application/pdf"},
		{"drops parameters", "page.html", nil, "text/html"},
		{"sniffed without extension", "screenshot", png, "image/png"},
		{"sniffed with unknown extension", "screenshot.unknown", png, "image/png"},
		{"unrecognised content", "blob", []byte{0x00, 0x01, 0x02}, "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			assert := require.New(t)

			// Act
			got := fileContentType(tt.filename, tt.head)

			// Assert
			assert.Equal(tt.want, got)
		})
	}
}
