package linkding

import (
	"os"
	"path/filepath"
	"testing"
	"testing/iotest"

	"github.com/joho/godotenv"
)

func TestGetBookmarks(t *testing.T) {
	godotenv.Load("../.env")
	client, err := NewClient(os.Getenv("LD_BASEURL"), os.Getenv("LD_TOKEN"))
	check(t, err)

	bookmarks, err := client.GetBookmarks("video")
	check(t, err)

	if len(bookmarks) == 0 {
		t.Fatal("No bookmarks found")
	}
}

func TestAddBookmarkAsset(t *testing.T) {
	const expectedDisplayName = "test-asset.txt"
	const expectedContentType = "text/plain"
	expectedContent := []byte("Test content")

	godotenv.Load("../.env")
	client, err := NewClient(os.Getenv("LD_BASEURL"), os.Getenv("LD_TOKEN"))
	check(t, err)

	file, err := os.Create(filepath.Join(t.TempDir(), expectedDisplayName))
	check(t, err)
	defer file.Close()

	file.Write(expectedContent)
	file.Sync()
	file.Seek(0, 0)

	bookmarks, err := client.GetBookmarks("video")
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

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
