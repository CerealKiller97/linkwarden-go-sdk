package linkwarden

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// newTestClient spins up an httptest.Server driven by handler and returns a Client
// pointed at it. The server is closed automatically when the test ends.
func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return New(server.URL, "test-token")
}

func TestNew_TrimsTrailingSlash(t *testing.T) {
	assert := require.New(t)

	c := New("https://example.com/", "tok")
	assert.Equal("https://example.com", c.baseURL)
}

func TestRequest_SetsAuthAndAcceptHeaders(t *testing.T) {
	assert := require.New(t)

	var gotAuth, gotAccept string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotAccept = r.Header.Get("Accept")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":[]}`))
	})

	_, err := client.Collections.List(t.Context())
	assert.NoError(err)
	assert.Equal("Bearer test-token", gotAuth)
	assert.Equal("application/json", gotAccept)
}

func TestRequest_JSONBodySetsContentType(t *testing.T) {
	assert := require.New(t)

	var gotContentType string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{}}`))
	})

	_, err := client.Collections.Create(t.Context(), CreateCollectionRequest{Name: "x"})
	assert.NoError(err)
	assert.Equal("application/json", gotContentType)
}

func TestRequest_NonOKStatusReturnsAPIError(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"response":"Collection not found."}`))
	})

	_, err := client.Collections.Get(t.Context(), 42)
	assert.Error(err)

	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
	assert.Equal(http.StatusNotFound, apiErr.StatusCode)
	assert.Equal("Collection not found.", apiErr.Message)
	assert.True(IsNotFound(err))
	assert.False(IsConflict(err))
}

func TestRequest_ErrorMessageFallsBackToStatusText(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`not json`))
	})

	_, err := client.Collections.Get(t.Context(), 1)
	assert.Error(err)

	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
	assert.Equal(http.StatusText(http.StatusInternalServerError), apiErr.Message)
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
