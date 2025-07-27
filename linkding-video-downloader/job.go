package main

import (
	"linkding-video-downloader/linkding"
	"linkding-video-downloader/ytdlp"
	"log/slog"
	"os"
	"slices"
	"sync"
)

func processBookmarks(client *linkding.Client, ytdlp *ytdlp.Ytdlp, tag string) (err error) {
	const concurrency = 4

	bookmarks, err := client.GetBookmarks(tag)
	if err != nil {
		return
	}

	slog.Info("Processing bookmarks", "count", len(bookmarks), "concurrency", concurrency)

	var wg sync.WaitGroup
	jobs := make(chan *linkding.Bookmark, concurrency)
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

	return
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
		return asset.AssetType == "upload" && linkding.IsKnownMimeType(asset.ContentType)
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
