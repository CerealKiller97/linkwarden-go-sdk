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

	password := uniqueName("password")
	user, err := c.Users.Create(t.Context(), CreateUserRequest{
		Name:     "Test User",
		Username: uniqueName("user"),
		Password: password,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Users.Delete(t.Context(), user.ID) })

	return user, password
}

func TestUsersService_Me(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)

	me, err := client.Users.Me(t.Context())
	assert.NoError(err)
	assert.NotZero(me.ID)
	assert.NotEmpty(me.Username)
}

func TestUsersService_Get(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	me, err := client.Users.Me(t.Context())
	assert.NoError(err)

	got, err := client.Users.Get(t.Context(), me.ID)
	assert.NoError(err)
	assert.Equal(me.ID, got.ID)
	assert.Equal(me.Username, got.Username)
}

func TestUsersService_CreateListDelete(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	user, _ := createTestUser(t, client)
	assert.NotZero(user.ID)

	users, err := client.Users.List(t.Context())
	assert.NoError(err)

	var ids []int
	for _, u := range users {
		ids = append(ids, u.ID)
	}
	assert.Contains(ids, user.ID)

	assert.NoError(client.Users.Delete(t.Context(), user.ID))

	users, err = client.Users.List(t.Context())
	assert.NoError(err)
	for _, u := range users {
		assert.NotEqual(user.ID, u.ID)
	}
}

func TestUsersService_Update(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	user, password := createTestUser(t, client)

	// Sign in as the new user, so the update can't touch the shared test account.
	token, err := client.Session.Create(t.Context(), CreateSessionRequest{Username: user.Username, Password: password})
	assert.NoError(err)
	asUser := New(os.Getenv("LINKWARDEN_URL"), token)

	updated, err := asUser.Users.Update(t.Context(), user.ID, UpdateUserRequest{
		Name:                  new("Renamed User"),
		PreventDuplicateLinks: new(true),
	})
	assert.NoError(err)
	assert.Equal("Renamed User", *updated.Name)

	got, err := asUser.Users.Me(t.Context())
	assert.NoError(err)
	assert.Equal("Renamed User", *got.Name)
	assert.Equal(user.Username, got.Username)
	assert.True(got.PreventDuplicateLinks)
}

func TestUsersService_UpdatePreferences(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	user, password := createTestUser(t, client)

	token, err := client.Session.Create(t.Context(), CreateSessionRequest{Username: user.Username, Password: password})
	assert.NoError(err)
	asUser := New(os.Getenv("LINKWARDEN_URL"), token)

	prefs, err := asUser.Users.UpdatePreferences(t.Context(), user.ID, UserPreferences{
		Theme:            "light",
		ReadableFontSize: "18px",
	})
	assert.NoError(err)
	assert.Equal("light", prefs.Theme)
	assert.Equal("18px", prefs.ReadableFontSize)
}
