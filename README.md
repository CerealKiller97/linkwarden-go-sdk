# linkwarden-go-sdk

[![Go Reference](https://pkg.go.dev/badge/github.com/CerealKiller97/linkwarden-go-sdk.svg)](https://pkg.go.dev/github.com/CerealKiller97/linkwarden-go-sdk)
[![test](https://github.com/CerealKiller97/linkwarden-go-sdk/actions/workflows/test.yml/badge.svg)](https://github.com/CerealKiller97/linkwarden-go-sdk/actions/workflows/test.yml)

A typed Go client for the [Linkwarden](https://linkwarden.app) API. It covers links,
collections, tags, search, archives, migrations, RSS, highlights, tokens and more, and depends
only on the standard library.

It powers [linkwarden-cli](https://github.com/CerealKiller97/linkwarden-cli).

## Install

Requires **Go 1.27+**.

```bash
go get github.com/CerealKiller97/linkwarden-go-sdk
```

The package is named `linkwarden`.

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func main() {
	ctx := context.Background()

	client := linkwarden.New(
		"https://cloud.linkwarden.app",
		os.Getenv("LINKWARDEN_TOKEN"),
		linkwarden.WithHTTPClient(&http.Client{Timeout: 15 * time.Second}),
	)

	// Create a collection…
	col, err := client.Collections.Create(ctx, linkwarden.CreateCollectionRequest{
		Name:  "Reading List",
		Color: "#f59e0b",
	})
	if err != nil {
		log.Fatal(err)
	}

	// …and save a tagged link into it.
	link, err := client.Links.Create(ctx, linkwarden.CreateLinkRequest{
		Name:       "Effective Go",
		URL:        "https://go.dev/doc/effective_go",
		Collection: linkwarden.Collection{ID: col.ID},
		Tags:       []linkwarden.Tag{{Name: "golang"}, {Name: "docs"}},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("saved #%d → %s\n", link.ID, link.URL)
}
```

<details>
<summary><b>More SDK examples</b></summary>

**Search and list**

```go
hits, err := client.Search.Links(ctx, linkwarden.SearchOptions{SearchQueryString: "kubernetes"})

pinned, err := client.Links.List(ctx, linkwarden.ListLinksOptions{
	CollectionID: 7,
	PinnedOnly:   true,
})
```

**Bulk move and retag**

```go
err := client.Links.BulkUpdate(ctx, linkwarden.BulkUpdateLinksRequest{
	Links:              []linkwarden.Link{{ID: 412}, {ID: 418}},
	RemovePreviousTags: true,
	NewData: linkwarden.BulkUpdateLinksNewData{
		CollectionID: 7,
		Tags:         []linkwarden.Tag{{Name: "weekend"}},
	},
})
```

**Download an archived PDF**

```go
pdf, err := client.Archives.Get(ctx, link.ID, linkwarden.ArchiveFormatPDF, false)
if err != nil {
	log.Fatal(err)
}
os.WriteFile("effective-go.pdf", pdf.Data, 0o644)
```

**Back up everything**

```go
data, err := client.Migrations.Export(ctx)
if err != nil {
	log.Fatal(err)
}
os.WriteFile("backup.json", data, 0o600)
```

**Handle API errors**

```go
_, err := client.Links.Get(ctx, 999999)
switch {
case linkwarden.IsNotFound(err):
	fmt.Println("no such link")
case linkwarden.IsConflict(err):
	fmt.Println("already exists")
case err != nil:
	var apiErr *linkwarden.APIError
	if errors.As(err, &apiErr) {
		fmt.Println(apiErr.StatusCode, apiErr.Message)
	}
}
```

</details>

| Service              | What it covers                                        |
| -------------------- | ----------------------------------------------------- |
| `client.Links`       | CRUD, bulk update/delete, archive, highlights         |
| `client.Collections` | CRUD, nesting, sharing                                |
| `client.Tags`        | List, bulk create, rename, delete                     |
| `client.Search`      | Full-text link search                                 |
| `client.Archives`    | Get and upload PNG / JPEG / PDF / Readability / HTML  |
| `client.Migrations`  | Export and import (Linkwarden, HTML, Pocket, …)       |
| `client.RSS`         | Feed subscriptions                                    |
| `client.Highlights`  | Create and delete highlights                          |
| `client.Tokens`      | Create, list, revoke API tokens                       |
| `client.Dashboard`   | v1/v2 dashboard data and layout                       |
| `client.Config`      | Public instance configuration                         |
| `client.Favicon`     | Favicon lookup for any URL                            |

Options: `WithHTTPClient(*http.Client)` for timeouts, proxies or custom transports, and `WithLogger(*slog.Logger)` for internal diagnostics (defaults to a text handler on stdout; pass `slog.New(slog.DiscardHandler)` to silence it). The SDK depends only on the standard library.


## Development

The tests run against a real, throwaway Linkwarden instance instead of mocked responses. Start
one in Docker (it also creates a test user and API token), then run the tests against it:

```bash
scripts/dev-linkwarden.sh up
eval "$(scripts/dev-linkwarden.sh env)"
go test ./... -race
golangci-lint run        # lint (config in .golangci.yml)
```

Without `LINKWARDEN_URL` and `LINKWARDEN_TOKEN` the tests are skipped locally, and fail in CI.
`scripts/dev-linkwarden.sh reset` deletes the instance and its data.

CI runs tests, `golangci-lint`, `gosec` and `govulncheck` on pushes to `main` and on pull requests.

Releases are plain SemVer git tags (`v0.1.0`, `v0.2.0`, …). The module stays on `v0.x` until the
API is stable.

## License

[GPL-3.0](LICENSE)
