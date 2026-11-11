package linkwarden

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// createTestUser registers a new user with a random password, returning it and the
// password. The calling client must be the instance admin to delete it afterwards.
func createTestUser(t *testing.T, c *Client) (*User, string) {
	t.Helper()
	assert := require.New(t)

	password := uniqueName("password")
	user, err := c.Users.Create(t.Context(), CreateUserRequest{
		Name:     "Test User",
		Username: uniqueName("user"),
		Password: password,
	})
	assert.NoError(err)
	t.Cleanup(func() { _ = c.Users.Delete(t.Context(), user.ID) })

	return user, password
}

// signInAs returns a client authenticated as user, so a test can change that user's
// account without touching the shared test account.
func signInAs(t *testing.T, c *Client, user *User, password string) *Client {
	t.Helper()
	assert := require.New(t)

	token, err := c.Session.Create(t.Context(), CreateSessionRequest{Username: user.Username, Password: password})
	assert.NoError(err)

	return New(os.Getenv("LINKWARDEN_URL"), token)
}

func TestUsersServiceMe(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)

	// Act
	me, err := client.Users.Me(t.Context())

	// Assert
	assert.NoError(err)
	assert.NotZero(me.ID)
	assert.NotEmpty(me.Username)
}

func TestUsersServiceGet(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	me, err := client.Users.Me(t.Context())
	assert.NoError(err)

	// Act
	got, err := client.Users.Get(t.Context(), me.ID)

	// Assert
	assert.NoError(err)
	assert.Equal(me.ID, got.ID)
	assert.Equal(me.Username, got.Username)
}

func TestUsersServiceCreate(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	username := uniqueName("user")

	// Act
	user, err := client.Users.Create(t.Context(), CreateUserRequest{
		Name:     "Test User",
		Username: username,
		Password: uniqueName("password"),
	})

	// Assert
	assert.NoError(err)
	t.Cleanup(func() { _ = client.Users.Delete(t.Context(), user.ID) })
	assert.NotZero(user.ID)
	assert.Equal(username, user.Username)
}

func TestUsersServiceList(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	user, _ := createTestUser(t, client)

	// Act
	users, err := client.Users.List(t.Context())

	// Assert
	assert.NoError(err)
	var ids []int
	for _, u := range users {
		ids = append(ids, u.ID)
	}
	assert.Contains(ids, user.ID)
}

func TestUsersServiceDelete(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	user, _ := createTestUser(t, client)

	// Act
	err := client.Users.Delete(t.Context(), user.ID)

	// Assert
	assert.NoError(err)
	users, err := client.Users.List(t.Context())
	assert.NoError(err)
	for _, u := range users {
		assert.NotEqual(user.ID, u.ID)
	}
}

func TestUsersServiceUpdate(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	user, password := createTestUser(t, client)
	asUser := signInAs(t, client, user, password)

	// Act
	updated, err := asUser.Users.Update(t.Context(), user.ID, UpdateUserRequest{
		Name:                  new("Renamed User"),
		PreventDuplicateLinks: new(true),
	})

	// Assert
	assert.NoError(err)
	assert.Equal("Renamed User", *updated.Name)
	got, err := asUser.Users.Me(t.Context())
	assert.NoError(err)
	assert.Equal("Renamed User", *got.Name)
	assert.Equal(user.Username, got.Username)
	assert.True(got.PreventDuplicateLinks)
}

func TestUsersServiceUpdatePreferences(t *testing.T) {
	// Arrange
	assert := require.New(t)
	client := newTestClient(t)
	user, password := createTestUser(t, client)
	asUser := signInAs(t, client, user, password)

	shareUsageAnalytics := false

	// Act
	prefs, err := asUser.Users.UpdatePreferences(t.Context(), user.ID, UserPreferences{
		Theme:               "light",
		ReadableFontSize:    "18px",
		ShareUsageAnalytics: &shareUsageAnalytics,
	})

	// Assert
	assert.NoError(err)
	assert.Equal("light", prefs.Theme)
	assert.Equal("18px", prefs.ReadableFontSize)
	assert.NotNil(prefs.ShareUsageAnalytics)
	assert.False(*prefs.ShareUsageAnalytics)
}
