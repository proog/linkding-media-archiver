package main

import (
	"linkding-video-downloader/linkding"
	"linkding-video-downloader/ytdlp"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

func processBookmarks(client *linkding.Client, ytdlp *ytdlp.Ytdlp, query *linkding.BookmarksQuery, isDryRun bool) (err error) {
	logger := slog.With("tag", query.Tag, "isDryRun", isDryRun)

	const concurrency = 4

	bookmarks, err := client.GetBookmarks(query)
	if err != nil {
		return
	}

	logger.Info("Processing bookmarks", "count", len(bookmarks))

	var wg sync.WaitGroup
	bookmarkJobs := make(chan linkding.Bookmark, len(bookmarks))
	failedCount := 0

	for w := 1; w <= concurrency; w++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for bookmark := range bookmarkJobs {
				if err := processBookmark(client, ytdlp, &bookmark, isDryRun); err != nil {
					failedCount++
				}
			}
		}()
	}

	for _, bookmark := range bookmarks {
		bookmarkJobs <- bookmark
	}

	close(bookmarkJobs)
	wg.Wait()

	logger.Info("Done processing bookmarks", "succeeded", len(bookmarks)-failedCount, "failed", failedCount)

	return
}

func processBookmark(client *linkding.Client, ytdlp *ytdlp.Ytdlp, bookmark *linkding.Bookmark, isDryRun bool) (err error) {
	logger := slog.With("bookmarkId", bookmark.Id, "isDryRun", isDryRun)
	logger.Info("Processing bookmark")

	assets, err := client.GetBookmarkAssets(bookmark.Id)

	if err != nil {
		logger.Error("Failed to fetch bookmark assets")
		return
	}

	mediaAssetIndex := slices.IndexFunc(assets, func(asset linkding.Asset) bool {
		return asset.AssetType == "upload" && linkding.IsKnownMimeType(asset.ContentType)
	})

	if mediaAssetIndex > -1 {
		logger.Info("Media asset already exists", "assetId", assets[mediaAssetIndex].Id)
		return
	}

	path, err := ytdlp.DownloadMedia(bookmark.Url)

	if err != nil {
		logger.Error("Failed to download media", "error", err)
		return
	}

	file, err := os.Open(path)

	if err != nil {
		logger.Error("Failed to open media file", "path", path, "error", err)
		return
	}

	defer file.Close()
	defer os.Remove(path)

	var asset *linkding.Asset

	if isDryRun {
		mimeType, err := linkding.GetMimeType(file.Name())
		if err != nil {
			return err
		}

		asset = &linkding.Asset{Id: -1, AssetType: "upload", ContentType: mimeType, DisplayName: "Simulated Asset" + filepath.Ext(file.Name())}
	} else {
		asset, err = client.AddBookmarkAsset(bookmark.Id, file)

		if err != nil {
			logger.Error("Failed to add asset", "error", err)
			return
		}
	}

	logger.Info("Bookmark processed successfully", "assetId", asset.Id)
	return
}
