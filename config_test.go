package linkwarden

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigServiceGet(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)

	// Act
	cfg, err := client.Config.Get(t.Context())

	// Assert
	assert.NoError(err)
	assert.NotNil(cfg.InstanceVersion)
	assert.NotEmpty(*cfg.InstanceVersion)
}
