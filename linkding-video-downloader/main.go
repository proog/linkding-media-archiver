package main

import (
	"linkding-video-downloader/linkding"
	"linkding-video-downloader/logging"
	"linkding-video-downloader/ytdlp"
	"log/slog"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

const maxConcurrency = 4

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	logger := logging.NewLogger()
	slog.SetDefault(logger)

	client, err := linkding.NewClient(os.Getenv("LD_BASEURL"), os.Getenv("LD_TOKEN"))
	if err != nil {
		panic(err)
	}

	tempdir, err := os.MkdirTemp(os.TempDir(), "videos")
	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(tempdir)
	ytdlp := ytdlp.NewYtdlp(tempdir)

	tag := os.Getenv("LD_TAG")
	if tag == "" {
		tag = "video"
	}

	bookmarks, err := client.GetBookmarks(tag)
	if err != nil {
		panic(err)
	}

	processBookmarks(client, ytdlp, bookmarks)
}

func processBookmarks(client *linkding.Client, ytdlp *ytdlp.Ytdlp, bookmarks []linkding.Bookmark) {
	jobs := make(chan *linkding.Bookmark, maxConcurrency)
	var wg sync.WaitGroup

	slog.Info("Processing bookmarks", "count", len(bookmarks))
	failedCount := 0

	for _, bookmark := range bookmarks {
		wg.Add(1)
		jobs <- &bookmark

		go func() {
			defer wg.Done()
			defer func() { <-jobs }()

			if err := processBookmark(client, ytdlp, &bookmark); err != nil {
				failedCount++
			}
		}()
	}

	wg.Wait()

	slog.Info("Done processing bookmarks", "succeeded", len(bookmarks)-failedCount, "failed", failedCount)
}

func processBookmark(client *linkding.Client, ytdlp *ytdlp.Ytdlp, bookmark *linkding.Bookmark) (err error) {
	logger := slog.With("bookmarkId", bookmark.Id)
	logger.Info("Processing bookmark")

	assets, err := client.GetBookmarkAssets(bookmark.Id)

	if err != nil {
		logger.Error("Failed to fetch bookmark assets")
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
