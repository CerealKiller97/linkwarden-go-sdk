package linkwarden

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMigrationsService_Export(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	collection := createTestCollection(t, client)
	link := createTestLink(t, client, collection.ID)

	data, err := client.Migrations.Export(t.Context())
	assert.NoError(err)
	assert.Contains(string(data), collection.Name)
	assert.Contains(string(data), link.URL)
}

func TestMigrationsService_Import(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t)
	folder := uniqueName("imported")
	url := uniqueURL()

	bookmarks := fmt.Sprintf(`<!DOCTYPE NETSCAPE-Bookmark-file-1>
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><H3>%s</H3>
    <DL><p>
        <DT><A HREF="%s">Imported link</A>
    </DL><p>
</DL><p>
`, folder, url)

	_, err := client.Migrations.Import(t.Context(), ImportFormatBookmarksHTML, bookmarks)
	assert.NoError(err)

	// The import runs in the background; wait for its folder to show up as a collection.
	var imported *Collection
	assert.Eventually(func() bool {
		collections, err := client.Collections.List(t.Context())
		if err != nil {
			return false
		}
		for i := range collections {
			if collections[i].Name == folder {
				imported = &collections[i]
				return true
			}
		}

		return false
	}, time.Minute, time.Second)
	t.Cleanup(func() { _ = client.Collections.Delete(t.Context(), imported.ID) })

	links, err := client.Links.List(t.Context(), ListLinksOptions{CollectionID: imported.ID})
	assert.NoError(err)
	assert.Len(links, 1)
	assert.Equal(url, links[0].URL)
}
