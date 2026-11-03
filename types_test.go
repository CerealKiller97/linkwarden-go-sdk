package linkwarden

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPinnedBy_UnmarshalJSON(t *testing.T) {
	assert := require.New(t)

	var link Link
	assert.NoError(json.Unmarshal([]byte(`{"pinnedBy":[1,{"id":2},null]}`), &link))
	assert.Len(link.PinnedBy, 3)
	assert.Equal(1, *link.PinnedBy[0].ID)
	assert.Equal(2, *link.PinnedBy[1].ID)
	assert.Nil(link.PinnedBy[2].ID)
}
