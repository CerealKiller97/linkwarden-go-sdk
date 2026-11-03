package linkwarden

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// User is a Linkwarden user account. The endpoints return different subsets of
// these fields, so anything an endpoint doesn't send is left at its zero value.
// Subscription, ParentSubscription and DashboardSections are kept as raw JSON since
// their shapes aren't part of the published API schema.
type User struct {
	ID                      int             `json:"id"`
	Name                    *string         `json:"name,omitempty"`
	Username                string          `json:"username,omitempty"`
	Email                   *string         `json:"email,omitempty"`
	EmailVerified           *string         `json:"emailVerified,omitempty"`
	UnverifiedNewEmail      *string         `json:"unverifiedNewEmail,omitempty"`
	Image                   *string         `json:"image,omitempty"`
	Locale                  string          `json:"locale,omitempty"`
	CollectionOrder         []int           `json:"collectionOrder,omitempty"`
	LinksRouteTo            string          `json:"linksRouteTo,omitempty"`
	PreventDuplicateLinks   bool            `json:"preventDuplicateLinks,omitempty"`
	ArchiveAsScreenshot     bool            `json:"archiveAsScreenshot,omitempty"`
	ArchiveAsMonolith       bool            `json:"archiveAsMonolith,omitempty"`
	ArchiveAsPDF            bool            `json:"archiveAsPDF,omitempty"`
	ArchiveAsWaybackMachine bool            `json:"archiveAsWaybackMachine,omitempty"`
	IsPrivate               bool            `json:"isPrivate,omitempty"`
	ReferredBy              *string         `json:"referredBy,omitempty"`
	WhitelistedUsers        []string        `json:"whitelistedUsers,omitempty"`
	HasUnIndexedLinks       bool            `json:"hasUnIndexedLinks,omitempty"`
	CreatedAt               time.Time       `json:"createdAt,omitzero"`
	UpdatedAt               time.Time       `json:"updatedAt,omitzero"`
	DashboardSections       json.RawMessage `json:"dashboardSections,omitempty"`
	Subscription            json.RawMessage `json:"subscription,omitempty"`
	ParentSubscription      json.RawMessage `json:"parentSubscription,omitempty"`
}

// CreateUserRequest is the payload for UsersService.Create. Whether Email or
// Username is required depends on the instance's configuration.
type CreateUserRequest struct {
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Invite   bool   `json:"invite,omitempty"`
}

// UpdateUserRequest is the payload for UsersService.Update. Nil fields are omitted.
type UpdateUserRequest struct {
	Name                    *string `json:"name,omitempty"`
	Email                   *string `json:"email,omitempty"`
	Username                *string `json:"username,omitempty"`
	Image                   *string `json:"image,omitempty"`
	Password                *string `json:"password,omitempty"`
	NewPassword             *string `json:"newPassword,omitempty"`
	OldPassword             *string `json:"oldPassword,omitempty"`
	ArchiveAsScreenshot     *bool   `json:"archiveAsScreenshot,omitempty"`
	ArchiveAsPDF            *bool   `json:"archiveAsPDF,omitempty"`
	ArchiveAsMonolith       *bool   `json:"archiveAsMonolith,omitempty"`
	ArchiveAsWaybackMachine *bool   `json:"archiveAsWaybackMachine,omitempty"`
	Locale                  *string `json:"locale,omitempty"`
	IsPrivate               *bool   `json:"isPrivate,omitempty"`
	PreventDuplicateLinks   *bool   `json:"preventDuplicateLinks,omitempty"`
	CollectionOrder         []int   `json:"collectionOrder,omitempty"`
	// LinksRouteTo is "ORIGINAL" or "CUSTOM".
	LinksRouteTo     *string  `json:"linksRouteTo,omitempty"`
	WhitelistedUsers []string `json:"whitelistedUsers,omitempty"`
	ReferredBy       *string  `json:"referredBy,omitempty"`
}

// UserPreferences holds a user's reading and theme preferences. As a request to
// UsersService.UpdatePreferences, empty fields are omitted.
type UserPreferences struct {
	ID int `json:"id,omitempty"`
	// Theme is "dark", "light" or "auto".
	Theme              string `json:"theme,omitempty"`
	ReadableFontFamily string `json:"readableFontFamily,omitempty"`
	ReadableFontSize   string `json:"readableFontSize,omitempty"`
	ReadableLineHeight string `json:"readableLineHeight,omitempty"`
	ReadableLineWidth  string `json:"readableLineWidth,omitempty"`
}

// UsersService accesses the /api/v1/users endpoints.
type UsersService struct {
	c *Client
}

// List returns the users on the instance that the authenticated user can see.
func (s *UsersService) List(ctx context.Context) ([]User, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/users",
	})
	if err != nil {
		return nil, err
	}

	return decodeResponseEnvelope[[]User](body)
}

// Create registers a new user. It doesn't require authentication unless the
// instance restricts registration.
func (s *UsersService) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPost,
		path:     "/api/v1/users",
		jsonBody: req,
	})
	if err != nil {
		return nil, err
	}

	user, err := decodeResponseEnvelope[User](body)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Me returns the authenticated user.
func (s *UsersService) Me(ctx context.Context) (*User, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/users/me",
	})
	if err != nil {
		return nil, err
	}

	user, err := decodeResponseEnvelope[User](body)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Get returns a user's account data. Linkwarden v2.16 answers with the authenticated
// user whatever the ID, so this is only useful with the caller's own ID (see Me).
func (s *UsersService) Get(ctx context.Context, id int) (*User, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/users/" + strconv.Itoa(id),
	})
	if err != nil {
		return nil, err
	}

	user, err := decodeResponseEnvelope[User](body)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update changes the fields set in req and keeps the rest. The API applies the update
// to the authenticated user and requires a username, so when req.Username is nil
// Update sends the user's current one.
func (s *UsersService) Update(ctx context.Context, id int, req UpdateUserRequest) (*User, error) {
	if req.Username == nil {
		me, err := s.Me(ctx)
		if err != nil {
			return nil, err
		}
		req.Username = &me.Username
	}

	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPut,
		path:     "/api/v1/users/" + strconv.Itoa(id),
		jsonBody: req,
	})
	if err != nil {
		return nil, err
	}

	user, err := decodeResponseEnvelope[User](body)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Delete deletes a user account and all its data.
func (s *UsersService) Delete(ctx context.Context, id int) error {
	_, _, err := s.c.request(ctx, requestOptions{
		method: http.MethodDelete,
		path:   "/api/v1/users/" + strconv.Itoa(id),
	})

	return err
}

// UpdatePreferences changes the reading and theme preferences of the user with the
// given ID and returns the result. Empty fields in prefs are left unchanged, and
// prefs.ID is ignored.
func (s *UsersService) UpdatePreferences(ctx context.Context, id int, prefs UserPreferences) (*UserPreferences, error) {
	prefs.ID = 0

	_, body, err := s.c.request(ctx, requestOptions{
		method:   http.MethodPut,
		path:     "/api/v1/users/" + strconv.Itoa(id) + "/preference",
		jsonBody: prefs,
	})
	if err != nil {
		return nil, err
	}

	updated, err := decodeResponseEnvelope[UserPreferences](body)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}
