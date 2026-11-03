package linkwarden

import (
	"context"
	"net/http"
)

// PublicConfig is the publicly-readable runtime configuration exposed by a
// Linkwarden instance. Fields are pointers since the API marks all of them
// nullable/optional.
type PublicConfig struct {
	DisableRegistration       *bool   `json:"DISABLE_REGISTRATION,omitempty"`
	Admin                     *int    `json:"ADMIN,omitempty"`
	RSSPollingIntervalMinutes *int    `json:"RSS_POLLING_INTERVAL_MINUTES,omitempty"`
	EmailProvider             *bool   `json:"EMAIL_PROVIDER,omitempty"`
	MaxFileBuffer             *int    `json:"MAX_FILE_BUFFER,omitempty"`
	UserContentDomain         *string `json:"USER_CONTENT_DOMAIN,omitempty"`
	AIEnabled                 *bool   `json:"AI_ENABLED,omitempty"`
	InstanceVersion           *string `json:"INSTANCE_VERSION,omitempty"`
}

// ConfigService accesses the /api/v1/config endpoint.
type ConfigService struct {
	c *Client
}

func (s *ConfigService) Get(ctx context.Context) (*PublicConfig, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/config",
	})
	if err != nil {
		return nil, err
	}

	cfg, err := decodeResponseEnvelope[PublicConfig](body)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
