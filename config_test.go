package linkwarden

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigService_Get(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/config", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"DISABLE_REGISTRATION":true,"MAX_FILE_BUFFER":100,"USER_CONTENT_DOMAIN":"https://usercontent.example.com","INSTANCE_VERSION":"2.10.0"}}`))
	})

	cfg, err := client.Config.Get(t.Context())
	assert.NoError(err)
	assert.NotNil(cfg.DisableRegistration)
	assert.True(*cfg.DisableRegistration)
	assert.Equal(100, *cfg.MaxFileBuffer)
	assert.Equal("https://usercontent.example.com", *cfg.UserContentDomain)
	assert.Equal("2.10.0", *cfg.InstanceVersion)
	assert.Nil(cfg.Admin)
}
