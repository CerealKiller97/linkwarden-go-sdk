package linkwarden

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokensService_List(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/tokens", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":[{"id":1,"name":"CI token","isSession":false}]}`))
	})

	tokens, err := client.Tokens.List(t.Context())
	assert.NoError(err)
	assert.Len(tokens, 1)
	assert.Equal("CI token", tokens[0].Name)
}

func TestTokensService_Create(t *testing.T) {
	assert := require.New(t)

	var gotBody struct {
		Name    string `json:"name"`
		Expires int    `json:"expires"`
	}
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/api/v1/tokens", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"secretKey":"secret-abc","token":{"id":1,"name":"CI token"}}}`))
	})

	created, err := client.Tokens.Create(t.Context(), "CI token", 2)
	assert.NoError(err)
	assert.Equal("secret-abc", created.SecretKey)
	assert.Equal("CI token", gotBody.Name)
	assert.Equal(2, gotBody.Expires)
}

func TestTokensService_Revoke(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/v1/tokens/5", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Token revoked."}`))
	})

	assert.NoError(client.Tokens.Revoke(t.Context(), 5))
}

func TestTokensService_Revoke_Unauthorized(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"Unauthorized"}`))
	})

	err := client.Tokens.Revoke(t.Context(), 5)
	assert.Error(err)

	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
	assert.Equal(http.StatusUnauthorized, apiErr.StatusCode)
}
