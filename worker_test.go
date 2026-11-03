package linkwarden

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// The test user is the instance's first user, which Linkwarden treats as the admin.

func TestWorkerServiceStats(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	createTestLink(t, client, createTestCollection(t, client).ID)

	// Act
	stats, err := client.Worker.Stats(t.Context())

	// Assert
	assert.NoError(err)
	assert.Positive(stats.Link.Pending + stats.Link.Done + stats.Link.Failed)
}

func TestWorkerServiceResetPreservation(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)

	// Act
	msg, err := client.Worker.ResetPreservation(t.Context(), PreservationResetBroken)

	// Assert
	assert.NoError(err)
	assert.NotEmpty(msg)
}

func TestWorkerServiceStatsNotAdmin(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	user, password := createTestUser(t, client)
	asUser := signInAs(t, client, user, password)

	// Act
	_, err := asUser.Worker.Stats(t.Context())

	// Assert
	assert.Error(err)
	var apiErr *APIError
	assert.ErrorAs(err, &apiErr)
	assert.Contains([]int{http.StatusUnauthorized, http.StatusForbidden}, apiErr.StatusCode)
}
