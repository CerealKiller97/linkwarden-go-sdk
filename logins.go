package linkwarden

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// LoginConfig describes the sign-in methods a Linkwarden instance offers. The
// flag fields are documented as nullable strings, so they're passed through as-is.
type LoginConfig struct {
	CredentialsEnabled   *string      `json:"credentialsEnabled"`
	EmailEnabled         *string      `json:"emailEnabled"`
	RegistrationDisabled *string      `json:"registrationDisabled"`
	ButtonAuths          []ButtonAuth `json:"buttonAuths"`
}

// ButtonAuth is one SSO provider offered on the login page.
type ButtonAuth struct {
	Method string `json:"method"`
	Name   string `json:"name"`
}

// LoginsService accesses the /api/v1/logins endpoint.
type LoginsService struct {
	c *Client
}

// Get returns the instance's login configuration. Unlike most endpoints, the
// response is not wrapped in an envelope.
func (s *LoginsService) Get(ctx context.Context) (*LoginConfig, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/logins",
	})
	if err != nil {
		return nil, err
	}

	var cfg LoginConfig
	if err := json.Unmarshal(body, &cfg); err != nil {
		return nil, fmt.Errorf("decode login config: %w", err)
	}

	return &cfg, nil
}
