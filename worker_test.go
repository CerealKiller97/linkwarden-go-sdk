package linkwarden

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// The test user is the instance's first user, which Linkwarden treats as the admin.

func TestWorkerServiceStats(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	createTestLink(t, client, createTestCollection(t, client).ID)

	stats, err := client.Worker.Stats(t.Context())
	assert.NoError(err)
	assert.Positive(stats.Link.Pending + stats.Link.Done + stats.Link.Failed)
}

func TestWorkerServiceResetPreservation(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)

	msg, err := client.Worker.ResetPreservation(t.Context(), PreservationResetBroken)
	assert.NoError(err)
	assert.NotEmpty(msg)
}

func TestWorkerServiceStatsNotAdmin(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	user, password := createTestUser(t, client)

	token, err := client.Session.Create(t.Context(), CreateSessionRequest{Username: user.Username, Password: password})
	assert.NoError(err)

	_, err = New(client.baseURL, token).Worker.Stats(t.Context())
	assert.Error(err)

	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
	assert.Contains([]int{http.StatusUnauthorized, http.StatusForbidden}, apiErr.StatusCode)
}
