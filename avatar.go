package linkwarden

import (
	"context"
	"net/http"
	"strconv"
)

// Avatar is a binary avatar image retrieved from AvatarService.Get.
type Avatar struct {
	Data        []byte
	ContentType string
}

// AvatarService accesses the /api/v1/avatar endpoint.
type AvatarService struct {
	c *Client
}

func (s *AvatarService) Get(ctx context.Context, userID int) (*Avatar, error) {
	resp, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/avatar/" + strconv.Itoa(userID),
	})
	if err != nil {
		return nil, err
	}

	return &Avatar{Data: body, ContentType: resp.Header.Get("Content-Type")}, nil
}
