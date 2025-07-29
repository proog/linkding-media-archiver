package job

import (
	"linkding-media-archiver/internal/linkding"
	"linkding-media-archiver/internal/ytdlp"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

func ProcessBookmarks(client *linkding.Client, ytdlp *ytdlp.Ytdlp, config JobConfiguration) (err error) {
	logger := slog.With("tags", config.Tags, "isDryRun", config.IsDryRun)

	const concurrency = 4

	bookmarks, err := getBookmarks(client, config)
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
				if err := processBookmark(client, ytdlp, &bookmark, config.IsDryRun); err != nil {
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

func getBookmarks(client *linkding.Client, config JobConfiguration) ([]linkding.Bookmark, error) {
	bookmarks := make([]linkding.Bookmark, 0, 100)

	tags := config.Tags
	if len(tags) == 0 {
		tags = []string{""}
	}

	for _, tag := range tags {
		query := linkding.BookmarksQuery{Tag: tag, ModifiedSince: config.LastScan}
		bookmarksForTag, err := client.GetBookmarks(query)

		if err != nil {
			return nil, err
		}

		for _, bookmarkForTag := range bookmarksForTag {
			exists := slices.ContainsFunc(bookmarks, func(b linkding.Bookmark) bool { return b.Id == bookmarkForTag.Id })

			if !exists {
				bookmarks = append(bookmarks, bookmarkForTag)
			}
		}
	}

	return bookmarks, nil
}

func processBookmark(client *linkding.Client, ytdlp *ytdlp.Ytdlp, bookmark *linkding.Bookmark, isDryRun bool) error {
	logger := slog.With("bookmarkId", bookmark.Id, "isDryRun", isDryRun)
	logger.Info("Processing bookmark")

	assets, err := client.GetBookmarkAssets(bookmark.Id)

	if err != nil {
		logger.Error("Failed to fetch bookmark assets")
		return err
	}

	mediaAssetIndex := slices.IndexFunc(assets, func(asset linkding.Asset) bool {
		return asset.AssetType == "upload" && linkding.IsKnownMimeType(asset.ContentType)
	})

	if mediaAssetIndex > -1 {
		logger.Info("Media asset already exists", "assetId", assets[mediaAssetIndex].Id)
		return err
	}

	paths, err := ytdlp.DownloadMedia(bookmark.Url)

	if err != nil {
		logger.Error("Failed to download media", "error", err)
		return err
	}

	for _, path := range paths {
		file, err := os.Open(path)

		if err != nil {
			logger.Error("Failed to open media file", "path", path, "error", err)
			return err
		}

		defer file.Close()
		defer os.Remove(path)

		asset, err := uploadAsset(client, bookmark, file, isDryRun)

		if err != nil {
			logger.Error("Failed to add asset", "error", err)
			return err
		}

		logger.Info("Asset added successfully", "assetId", asset.Id)
	}

	logger.Info("Bookmark processed successfully")
	return nil
}

func uploadAsset(client *linkding.Client, bookmark *linkding.Bookmark, file *os.File, isDryRun bool) (*linkding.Asset, error) {
	if isDryRun {
		mimeType, err := linkding.GetMimeType(file.Name())
		if err != nil {
			return nil, err
		}

		asset := &linkding.Asset{Id: -1, AssetType: "upload", ContentType: mimeType, DisplayName: "Simulated Asset" + filepath.Ext(file.Name())}
		return asset, nil
	}

	return client.AddBookmarkAsset(bookmark.Id, file)
}
