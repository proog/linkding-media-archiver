package main

import (
	"linkding-video-downloader/linkding"
	"linkding-video-downloader/ytdlp"
	"log/slog"
	"os"
	"slices"
	"strings"
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

	defer os.RemoveAll(tempdir)
	ytdlp := ytdlp.NewYtdlp(tempdir)

	concurrency := 3
	sem := make(chan int, concurrency)
	var wg sync.WaitGroup

	slog.Info("Processing bookmarks", "count", len(bookmarks))

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

	slog.Info("Done processing bookmarks", "count", len(bookmarks))
}

func addVideoAsset(client *linkding.Client, ytdlp *ytdlp.Ytdlp, bookmark *linkding.Bookmark) (err error) {
	logger := slog.With("bookmarkId", bookmark.Id)
	logger.Info("Processing bookmark")

	assets, err := client.GetBookmarkAssets(bookmark.Id)

	if err != nil {
		return
	}

	videoAssetIndex := slices.IndexFunc(assets, func(asset linkding.Asset) bool {
		return asset.AssetType == "upload" && strings.HasPrefix(asset.ContentType, "video/")
	})

	if videoAssetIndex > -1 {
		logger.Info("Video asset already exists", "assetId", assets[videoAssetIndex].Id)
		return
	}

	path, err := ytdlp.DownloadVideo(bookmark.Url)

	if err != nil {
		logger.Error("Failed to download video", "error", err)
		return
	}

	file, err := os.Open(path)

	if err != nil {
		logger.Error("Failed to open video file", "path", path, "error", err)
		return
	}

	defer file.Close()
	defer os.Remove(path)

	asset, err := client.AddBookmarkAsset(bookmark.Id, file)

	if err != nil {
		logger.Error("Failed to add asset", "error", err)
		return
	}

	logger.Info("Bookmark processed successfully", "assetId", asset.Id)
	return
}
