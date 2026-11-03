# linkwarden-go-sdk

[![Go Reference](https://pkg.go.dev/badge/github.com/CerealKiller97/linkwarden-go-sdk.svg)](https://pkg.go.dev/github.com/CerealKiller97/linkwarden-go-sdk)
[![test](https://github.com/CerealKiller97/linkwarden-go-sdk/actions/workflows/test.yml/badge.svg)](https://github.com/CerealKiller97/linkwarden-go-sdk/actions/workflows/test.yml)
[![lint](https://github.com/CerealKiller97/linkwarden-go-sdk/actions/workflows/lint.yml/badge.svg)](https://github.com/CerealKiller97/linkwarden-go-sdk/actions/workflows/lint.yml)
[![gosec](https://github.com/CerealKiller97/linkwarden-go-sdk/actions/workflows/gosec.yml/badge.svg)](https://github.com/CerealKiller97/linkwarden-go-sdk/actions/workflows/gosec.yml)
[![govulncheck](https://github.com/CerealKiller97/linkwarden-go-sdk/actions/workflows/govulncheck.yml/badge.svg)](https://github.com/CerealKiller97/linkwarden-go-sdk/actions/workflows/govulncheck.yml)

A typed Go client for the [Linkwarden](https://linkwarden.app) API. It covers links,
collections, tags, search, archives, migrations, RSS, highlights, tokens and more, and depends
only on the standard library.

## Install

Requires **Go 1.27+**. Each release is tested against a specific Linkwarden version, listed in
[CHANGELOG.md](CHANGELOG.md).

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

### Services

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
| `client.Users`       | Current user, registration, profile and preferences   |
| `client.Session`     | Sign in with username and password for a token        |
| `client.Auth`        | Password reset and email verification                 |
| `client.Logins`      | Sign-in methods the instance offers                   |
| `client.Public`      | Read-only access to public collections, links, users  |
| `client.Avatar`      | User avatar images                                    |
| `client.Preserved`   | Signed URLs for archives on the user content domain   |
| `client.Worker`      | Preservation queue stats and resets (admin only)      |

Options: `WithHTTPClient(*http.Client)` for timeouts, proxies or custom transports, and `WithLogger(*slog.Logger)` for internal diagnostics. The client is silent by default: it discards all log output unless you pass a logger. The SDK depends only on the standard library.

## Examples

The snippets below assume the `client` and `ctx` from the example above.

### Links

<details>
<summary><b>Search, filter and list</b></summary>

```go
hits, err := client.Search.Links(ctx, linkwarden.SearchOptions{SearchQueryString: "kubernetes"})

pinned, err := client.Links.List(ctx, linkwarden.ListLinksOptions{
	CollectionID: 7,
	PinnedOnly:   true,
})
```

</details>

<details>
<summary><b>Page through every link</b></summary>

Lists are cursor-paginated: pass the ID of the last link you received to get the next page.

```go
var all []linkwarden.Link
opts := linkwarden.ListLinksOptions{}
for {
	page, err := client.Links.List(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}
	if len(page) == 0 {
		break
	}
	all = append(all, page...)
	opts.Cursor = page[len(page)-1].ID
}
```

</details>

<details>
<summary><b>Update, archive and delete</b></summary>

```go
// Only the fields you set are changed; the rest are kept from the current link.
link, err := client.Links.Update(ctx, 412, linkwarden.UpdateLinkRequest{
	Name: "Effective Go (updated)",
	Tags: []linkwarden.Tag{{Name: "golang"}},
})

// Re-run the archiver (screenshot, PDF, readability…) for a link.
err = client.Links.Archive(ctx, 412)

err = client.Links.Delete(ctx, 412)
err = client.Links.BulkDelete(ctx, []int{418, 419, 420})
```

</details>

<details>
<summary><b>Bulk move and retag</b></summary>

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

</details>

### Collections and tags

<details>
<summary><b>Nested collections and sharing</b></summary>

```go
parent, err := client.Collections.Create(ctx, linkwarden.CreateCollectionRequest{Name: "Programming"})

child, err := client.Collections.Create(ctx, linkwarden.CreateCollectionRequest{
	Name:     "Go",
	ParentID: parent.ID,
})

// Make it public and give user 3 read/write access.
_, err = client.Collections.Update(ctx, child.ID, linkwarden.UpdateCollectionRequest{
	IsPublic: new(true),
	Members: []linkwarden.Member{
		{UserID: 3, CanCreate: true, CanUpdate: true},
	},
})
```

</details>

<details>
<summary><b>Manage tags</b></summary>

```go
tags, err := client.Tags.List(ctx, linkwarden.ListTagsOptions{Search: "go"})
for _, t := range tags {
	if t.Count != nil {
		fmt.Printf("%-20s %d links\n", t.Name, t.Count.Links)
	}
}

_, err = client.Tags.CreateOrUpdateBulk(ctx, []linkwarden.CreateOrUpdateTagInput{
	{Label: "to-read"},
	{Label: "reference"},
})

_, err = client.Tags.Update(ctx, 12, "golang") // rename
err = client.Tags.BulkDelete(ctx, []int{13, 14})
```

</details>

### Archives, highlights and RSS

<details>
<summary><b>Download an archived PDF</b></summary>

```go
pdf, err := client.Archives.Get(ctx, link.ID, linkwarden.ArchiveFormatPDF, false)
if err != nil {
	log.Fatal(err)
}
os.WriteFile("effective-go.pdf", pdf.Data, 0o644)
```

</details>

<details>
<summary><b>Upload your own archive file</b></summary>

```go
f, err := os.Open("snapshot.png")
if err != nil {
	log.Fatal(err)
}
defer f.Close()

// Attach it to an existing link…
_, err = client.Archives.Upload(ctx, link.ID, linkwarden.ArchiveFormatPNG, false, "snapshot.png", f)

// …or create a new link from a file in one call.
// _, err = client.Archives.UploadStandalone(ctx, linkwarden.ArchiveFormatPDF, false,
// 	"https://example.com/paper", "paper.pdf", pdfReader)
```

</details>

<details>
<summary><b>Highlights</b></summary>

```go
h, err := client.Highlights.Create(ctx, linkwarden.CreateHighlightRequest{
	LinkID:      link.ID,
	Color:       "yellow",
	Text:        "Don't communicate by sharing memory",
	Comment:     "the classic",
	StartOffset: 120,
	EndOffset:   155,
})

highlights, err := client.Links.Highlights(ctx, link.ID)

err = client.Highlights.Delete(ctx, h.ID)
```

</details>

<details>
<summary><b>Subscribe to an RSS feed</b></summary>

New feed items are saved into the collection automatically.

```go
feed, err := client.RSS.Create(ctx, linkwarden.CreateRSSRequest{
	Name:           "Go Blog",
	URL:            "https://go.dev/blog/feed.atom",
	CollectionName: "Go",
})

feeds, err := client.RSS.List(ctx)
err = client.RSS.Delete(ctx, feed.ID)
```

</details>

### Import, export and account

<details>
<summary><b>Back up everything</b></summary>

```go
data, err := client.Migrations.Export(ctx)
if err != nil {
	log.Fatal(err)
}
os.WriteFile("backup.json", data, 0o600)
```

</details>

<details>
<summary><b>Import browser bookmarks</b></summary>

```go
html, err := os.ReadFile("bookmarks.html") // exported from Chrome, Firefox, Safari…
if err != nil {
	log.Fatal(err)
}

msg, err := client.Migrations.Import(ctx, linkwarden.ImportFormatBookmarksHTML, string(html))
fmt.Println(msg) // the import runs in the background on the server
```

Other formats: `ImportFormatLinkwardenJSON`, `ImportFormatWallabagJSON`,
`ImportFormatOmnivoreJSON`, `ImportFormatPocketJSON`.

</details>

<details>
<summary><b>API tokens</b></summary>

```go
created, err := client.Tokens.Create(ctx, "ci-bot", 0)
if err != nil {
	log.Fatal(err)
}
fmt.Println(created.SecretKey) // only returned once — store it now

tokens, err := client.Tokens.List(ctx)
err = client.Tokens.Revoke(ctx, created.Token.ID)
```

</details>

<details>
<summary><b>Sign in with a username and password</b></summary>

`Session`, `Auth`, `Logins`, `Public` and `Config` work without a token.

```go
anon := linkwarden.New(baseURL, "")

token, err := anon.Session.Create(ctx, linkwarden.CreateSessionRequest{
	Username: "dev",
	Password: os.Getenv("LINKWARDEN_PASSWORD"),
})
if err != nil {
	log.Fatal(err)
}

client := linkwarden.New(baseURL, token)
me, err := client.Users.Me(ctx)
```

</details>

<details>
<summary><b>Read a public collection</b></summary>

```go
anon := linkwarden.New(baseURL, "")

col, err := anon.Public.Collection(ctx, 7)
links, err := anon.Public.CollectionLinks(ctx, col.ID, linkwarden.ListPublicLinksOptions{
	SearchQueryString: "go",
	SearchByName:      true,
	SearchByTags:      true,
})
```

</details>

<details>
<summary><b>Instance config and favicons</b></summary>

```go
cfg, err := client.Config.Get(ctx)
if err == nil && cfg.InstanceVersion != nil {
	fmt.Println("Linkwarden", *cfg.InstanceVersion)
}

icon, err := client.Favicon.Get(ctx, "https://go.dev")
os.WriteFile("go.ico", icon.Data, 0o644)
```

</details>

### Client behaviour

<details>
<summary><b>Timeouts, cancellation and logging</b></summary>

```go
client := linkwarden.New(baseURL, token,
	linkwarden.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
	linkwarden.WithLogger(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))),
)

// Logging is off by default. Pass your app's logger to opt in:
//   linkwarden.WithLogger(slog.Default())

// Every call takes a context, so per-request deadlines work too.
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
links, err := client.Links.List(ctx, linkwarden.ListLinksOptions{})
```

</details>

<details>
<summary><b>Handle API errors</b></summary>

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

### Releases

Releases are SemVer git tags (`v1.0.0`, `v1.1.0`, …). Changes are listed in
[CHANGELOG.md](CHANGELOG.md), and each release there names the Linkwarden version its tests
passed against.

That version is pinned in [`.linkwarden-version`](.linkwarden-version). CI and
`scripts/dev-linkwarden.sh` both run it. To bump it, change the file and fix anything the
tests catch.

To cut a release:

1. In `CHANGELOG.md`, which follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
   rename `[Unreleased]` to `## [x.y.z] - YYYY-MM-DD`, keep its
   `🧪 Tested against Linkwarden: <version>` line, start a new empty `[Unreleased]`, and update
   the compare links at the bottom. The comment at the top of the file lists the section
   headings to use.
2. Check the notes: `scripts/release-notes.sh vX.Y.Z`
3. Merge that to `main`, then tag and push: `git tag -a vX.Y.Z -m vX.Y.Z && git push origin vX.Y.Z`

The release workflow then checks the changelog entry, runs the tests on the tag against the
pinned Linkwarden version, and publishes a GitHub release with that entry as its notes. It
fails if the entry is missing or names a different Linkwarden version.

## License

[GPL-3.0](LICENSE)
