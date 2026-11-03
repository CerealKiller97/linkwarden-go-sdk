package linkwarden

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// PublicUser is the public profile of a user, as returned by PublicService.User.
type PublicUser struct {
	ID                  int     `json:"id"`
	Name                string  `json:"name"`
	Username            string  `json:"username"`
	Email               *string `json:"email,omitempty"`
	Image               string  `json:"image,omitempty"`
	ArchiveAsScreenshot bool    `json:"archiveAsScreenshot"`
	ArchiveAsMonolith   bool    `json:"archiveAsMonolith"`
	ArchiveAsPDF        bool    `json:"archiveAsPDF"`
}

// ListPublicLinksOptions filters/paginates PublicService.CollectionLinks.
type ListPublicLinksOptions struct {
	// Sort is left nil to omit the query parameter and let the API pick its default.
	Sort              *int
	Cursor            int
	PinnedOnly        bool
	SearchQueryString string
	// The SearchBy* fields narrow which fields SearchQueryString matches against.
	SearchByName        bool
	SearchByURL         bool
	SearchByDescription bool
	SearchByTextContent bool
	SearchByTags        bool
}

func (o ListPublicLinksOptions) toQuery(collectionID int) url.Values {
	q := url.Values{}
	q.Set("collectionId", strconv.Itoa(collectionID))
	if o.Sort != nil {
		q.Set("sort", strconv.Itoa(*o.Sort))
	}
	setIfPositive(q, "cursor", o.Cursor)
	setIfNonEmpty(q, "searchQueryString", o.SearchQueryString)
	for key, on := range map[string]bool{
		"pinnedOnly":          o.PinnedOnly,
		"searchByName":        o.SearchByName,
		"searchByUrl":         o.SearchByURL,
		"searchByDescription": o.SearchByDescription,
		"searchByTextContent": o.SearchByTextContent,
		"searchByTags":        o.SearchByTags,
	} {
		if on {
			q.Set(key, "true")
		}
	}

	return q
}

// PublicService accesses the unauthenticated /api/v1/public endpoints, which only
// return data from public collections.
type PublicService struct {
	c *Client
}

// Collection returns a public collection. It doesn't require authentication.
func (s *PublicService) Collection(ctx context.Context, id int) (*Collection, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/public/collections/" + strconv.Itoa(id),
	})
	if err != nil {
		return nil, err
	}

	collection, err := decodeResponseEnvelope[Collection](body)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

// CollectionLinks returns one page of links in a public collection, filtered by
// opts. It doesn't require authentication.
func (s *PublicService) CollectionLinks(ctx context.Context, collectionID int, opts ListPublicLinksOptions) ([]Link, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/public/collections/links",
		query:  opts.toQuery(collectionID),
	})
	if err != nil {
		return nil, err
	}

	data, err := decodeDataEnvelope[struct {
		Links []Link `json:"links"`
	}](body)
	if err != nil {
		return nil, err
	}

	return data.Links, nil
}

// CollectionTags returns the tags used in a public collection. It doesn't require
// authentication.
func (s *PublicService) CollectionTags(ctx context.Context, collectionID int) ([]TagInfo, error) {
	q := url.Values{}
	q.Set("collectionId", strconv.Itoa(collectionID))

	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/public/collections/tags",
		query:  q,
	})
	if err != nil {
		return nil, err
	}

	data, err := decodeDataEnvelope[struct {
		Tags []TagInfo `json:"tags"`
	}](body)
	if err != nil {
		return nil, err
	}

	return data.Tags, nil
}

// Link returns a link from a public collection. It doesn't require authentication.
func (s *PublicService) Link(ctx context.Context, id int) (*Link, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/public/links/" + strconv.Itoa(id),
	})
	if err != nil {
		return nil, err
	}

	link, err := decodeResponseEnvelope[Link](body)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

// User returns a user's public profile. It doesn't require authentication.
func (s *PublicService) User(ctx context.Context, id int) (*PublicUser, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/public/users/" + strconv.Itoa(id),
	})
	if err != nil {
		return nil, err
	}

	user, err := decodeResponseEnvelope[PublicUser](body)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
