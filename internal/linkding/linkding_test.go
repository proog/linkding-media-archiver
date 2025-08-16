package linkding

import (
	"os"
	"path/filepath"
	"testing"
	"testing/iotest"
	"time"

	"github.com/joho/godotenv"
)

const validTag = "video"
const validBundleId = 2

func TestGetBookmarksWithTag(t *testing.T) {
	tests := []string{
		"",
		validTag,
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			client := getClient(t)
			bookmarks, err := client.GetBookmarks(BookmarksQuery{Tag: test})
			check(t, err)

			if len(bookmarks) == 0 {
				t.Fatal("No bookmarks found")
			}
		})
	}
}

func TestGetBookmarksWithBundle(t *testing.T) {
	client := getClient(t)
	allBookmarks, err := client.GetBookmarks(BookmarksQuery{})
	check(t, err)

	bookmarks, err := client.GetBookmarks(BookmarksQuery{BundleId: validBundleId})
	check(t, err)

	if !(len(bookmarks) < len(allBookmarks)) {
		t.Fatal("Bundle query did not filter the result")
	}
}

func TestGetBookmarksModifiedSince(t *testing.T) {
	tests := []time.Time{
		{},
		time.Now().Add(-1 * time.Hour),
	}

	for _, test := range tests {
		t.Run(test.Format(time.DateTime), func(t *testing.T) {
			client := getClient(t)
			bookmarks, err := client.GetBookmarks(BookmarksQuery{ModifiedSince: test})
			check(t, err)

			if len(bookmarks) == 0 {
				t.Fatal("No bookmarks found")
			}
		})
	}
}

func TestGetBookmarkAssets(t *testing.T) {
	client := getClient(t)

	bookmarks, err := client.GetBookmarks(BookmarksQuery{Tag: validTag})
	check(t, err)

	bookmark := bookmarks[0]
	assets, err := client.GetBookmarkAssets(bookmark.Id)
	check(t, err)

	if len(assets) == 0 {
		t.Fatal("No assets found")
	}
}

func TestAddBookmarkAsset(t *testing.T) {
	const expectedDisplayName = "test-asset.mp4"
	const expectedContentType = "video/mp4"
	expectedContent := []byte("Test content")

	client := getClient(t)

	file, err := os.Create(filepath.Join(t.TempDir(), expectedDisplayName))
	check(t, err)
	defer file.Close()

	file.Write(expectedContent)
	file.Sync()
	file.Seek(0, 0)

	bookmarks, err := client.GetBookmarks(BookmarksQuery{Tag: validTag})
	check(t, err)

	bookmark := bookmarks[0]
	asset, err := client.AddBookmarkAsset(bookmark.Id, file)
	check(t, err)

	if asset.DisplayName != expectedDisplayName {
		t.Fatalf("Expected display name %s, got %s", expectedDisplayName, asset.DisplayName)
	}

	if asset.ContentType != expectedContentType {
		t.Fatalf("Expected content type %s, got %s", expectedContentType, asset.ContentType)
	}

	download, err := client.DownloadBookmarkAsset(bookmark.Id, asset.Id)
	check(t, err)
	defer download.Close()

	err = iotest.TestReader(download, expectedContent)
	check(t, err)
}

func getClient(t *testing.T) *Client {
	godotenv.Load("../../.env")

	client, err := NewClient(os.Getenv("LDMA_BASEURL"), os.Getenv("LDMA_TOKEN"))
	check(t, err)

	return client
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
