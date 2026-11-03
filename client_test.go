package linkwarden

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew_TrimsTrailingSlash(t *testing.T) {
	assert := require.New(t)

	c := New("https://example.com/", "tok")
	assert.Equal("https://example.com", c.baseURL)
}

func TestRequest_AuthenticatesWithToken(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)

	_, err := client.Collections.List(t.Context())
	assert.NoError(err)
}

func TestRequest_InvalidTokenReturnsAPIError(t *testing.T) {
	assert := require.New(t)

	newTestClient(t)
	client := New(os.Getenv("LINKWARDEN_URL"), "invalid-token")

	_, err := client.Collections.List(t.Context())
	assert.Error(err)

	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
	assert.Equal(http.StatusUnauthorized, apiErr.StatusCode)
	assert.NotEmpty(apiErr.Message)
	assert.NotEmpty(apiErr.RawBody)
}

func TestNewAPIError_ReadsMessageField(t *testing.T) {
	assert := require.New(t)

	err := newAPIError(http.StatusBadRequest, []byte(`{"message":"Bad input."}`))
	assert.Equal("Bad input.", err.Message)
	assert.Equal("linkwarden: 400 Bad input.", err.Error())
}

func TestNewAPIError_ReadsResponseField(t *testing.T) {
	assert := require.New(t)

	err := newAPIError(http.StatusNotFound, []byte(`{"response":"Tag not found."}`))
	assert.Equal("Tag not found.", err.Message)
	assert.True(IsNotFound(err))
	assert.False(IsConflict(err))
}

func TestNewAPIError_FallsBackToStatusText(t *testing.T) {
	assert := require.New(t)

	err := newAPIError(http.StatusInternalServerError, []byte(`Internal Server Error`))
	assert.Equal(http.StatusText(http.StatusInternalServerError), err.Message)
}

func TestDecodeResponseEnvelope(t *testing.T) {
	assert := require.New(t)

	got, err := decodeResponseEnvelope[string]([]byte(`{"response":"hello"}`))
	assert.NoError(err)
	assert.Equal("hello", got)
}

func TestDecodeDataEnvelope(t *testing.T) {
	assert := require.New(t)

	got, err := decodeDataEnvelope[string]([]byte(`{"data":"hello","message":"Success","success":true}`))
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
			require.Equal(t, tt.want, fileContentType(tt.filename, tt.head))
		})
	}
}
