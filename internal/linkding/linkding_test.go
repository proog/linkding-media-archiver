package linkding

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"testing/iotest"
	"time"

	"github.com/joho/godotenv"
)

const validTag = "video"
const validTag2 = "music"
const validBundleId = 2

func TestGetBookmarksWithTags(t *testing.T) {
	tests := [][]string{
		{},
		{validTag},
		{validTag, validTag2},
		{"3922EBE2-A681-4EAD-A5EA-89FF4B2CBCBE", validTag2},
	}

	for _, test := range tests {
		t.Run(strings.Join(test, ","), func(t *testing.T) {
			client := getClient(t)
			bookmarks, err := client.GetBookmarks(BookmarksQuery{Tags: test})
			check(t, err)

			if len(bookmarks) == 0 {
				t.Fatal("No bookmarks found")
			}

			if len(test) > 0 {
				for _, bookmark := range bookmarks {
					hasTag := slices.ContainsFunc(bookmark.TagNames, func(tag string) bool { return slices.Contains(test, tag) })

					if !hasTag {
						t.Errorf("Bookmark %d did not have the expected tags", bookmark.Id)
					}
				}
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

	bookmarks, err := client.GetBookmarks(BookmarksQuery{Tags: []string{validTag}})
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

	bookmarks, err := client.GetBookmarks(BookmarksQuery{Tags: []string{validTag}})
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

func TestGetUserProfile(t *testing.T) {
	client := getClient(t)
	profile, err := client.GetUserProfile()
	check(t, err)

	if len(strings.SplitN(profile.Version, ".", 3)) != 3 {
		t.Fatalf("Expected version to be a semver, got %s", profile.Version)
	}
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
