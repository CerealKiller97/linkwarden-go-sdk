package linkwarden

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigServiceGet(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)

	cfg, err := client.Config.Get(t.Context())
	assert.NoError(err)
	assert.NotNil(cfg.InstanceVersion)
	assert.NotEmpty(*cfg.InstanceVersion)
}
