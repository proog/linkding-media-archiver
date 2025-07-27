package main

import (
	"linkding-video-downloader/linkding"
	"linkding-video-downloader/ytdlp"
	"log/slog"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	client, err := linkding.NewClient(os.Getenv("LD_BASEURL"), os.Getenv("LD_TOKEN"))

	if err != nil {
		panic(err)
	}

	bookmarks, err := client.GetBookmarks(os.Getenv("LD_TAG"))

	if err != nil {
		panic(err)
	}

	tempdir, err := os.MkdirTemp(os.TempDir(), "videos")

	if err != nil {
		panic(err)
	}

	defer os.Remove(tempdir)
	ytdlp := ytdlp.NewYtdlp(tempdir)

	concurrency := 3
	sem := make(chan int, concurrency)
	var wg sync.WaitGroup

	for _, bookmark := range bookmarks {
		wg.Add(1)
		sem <- bookmark.Id

		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			addVideoAsset(client, ytdlp, &bookmark)
		}()
	}

	wg.Wait()
}

func addVideoAsset(client *linkding.Client, ytdlp *ytdlp.Ytdlp, bookmark *linkding.Bookmark) error {
	logger := slog.With("bookmarkId", bookmark.Id)
	logger.Info("Processing bookmark")

	// TODO: check if video asset already exists

	// assets, err := ld.GetBookmarkAssets(bookmark.Id)

	path, err := ytdlp.DownloadVideo(bookmark.Url)

	if err != nil {
		logger.Error("Failed to download video", "error", err)
		return err
	}

	file, err := os.Open(path)

	if err != nil {
		logger.Error("Failed to open video file", "path", path, "error", err)
		return err
	}

	defer file.Close()
	defer os.Remove(path)

	asset, err := client.AddBookmarkAsset(bookmark.Id, file)

	if err != nil {
		logger.Error("Failed to add asset", "error", err)
		return err
	}

	logger.Info("Bookmark processed", "assetId", asset.Id)
	return nil
}
