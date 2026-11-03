package linkwarden

import (
	"net/url"
	"strconv"
	"time"
)

// Member represents a collection member's access grant, shared between
// collection read/update payloads.
type Member struct {
	UserID       int        `json:"userId"`
	CollectionID int        `json:"collectionId,omitempty"`
	CanCreate    bool       `json:"canCreate"`
	CanUpdate    bool       `json:"canUpdate"`
	CanDelete    bool       `json:"canDelete"`
	CreatedAt    time.Time  `json:"createdAt,omitzero"`
	UpdatedAt    time.Time  `json:"updatedAt,omitzero"`
	User         MemberUser `json:"user,omitzero"`
}

type MemberUser struct {
	Username string  `json:"username,omitempty"`
	Name     string  `json:"name,omitempty"`
	Image    *string `json:"image,omitempty"`
}

// PinnedBy represents an entry in a link's pinnedBy list.
type PinnedBy struct {
	ID *int `json:"id,omitempty"`
}

func setIfNonEmpty(q url.Values, key, val string) {
	if val != "" {
		q.Set(key, val)
	}
}

func setIfPositive(q url.Values, key string, v int) {
	if v > 0 {
		q.Set(key, strconv.Itoa(v))
	}
}
