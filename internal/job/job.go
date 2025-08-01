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

	bookmarks, err := getBookmarks(client, config)
	if err != nil {
		return
	}

	if len(bookmarks) == 0 {
		logger.Info("No bookmarks to process")
		return
	}

	logger.Info("Processing bookmarks", "count", len(bookmarks))

	var wg sync.WaitGroup
	succeeded := make(chan linkding.Bookmark, len(bookmarks))
	failed := make(chan linkding.Bookmark, len(bookmarks))

	for _, bookmark := range bookmarks {
		paths, err := downloadMedia(client, ytdlp, bookmark)

		if err != nil {
			failed <- bookmark
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := uploadMedia(client, bookmark, paths, config.IsDryRun); err != nil {
				failed <- bookmark
				return
			}

			succeeded <- bookmark
		}()
	}

	wg.Wait()

	logger.Info("Done processing bookmarks", "succeeded", len(succeeded), "failed", len(failed))

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

func downloadMedia(client *linkding.Client, ytdlp *ytdlp.Ytdlp, bookmark linkding.Bookmark) ([]string, error) {
	logger := slog.With("bookmarkId", bookmark.Id)
	assets, err := client.GetBookmarkAssets(bookmark.Id)

	if err != nil {
		logger.Error("Failed to fetch bookmark assets")
		return nil, err
	}

	mediaAssetIndex := slices.IndexFunc(assets, func(asset linkding.Asset) bool {
		return asset.AssetType == "upload" && linkding.IsKnownMimeType(asset.ContentType)
	})

	if mediaAssetIndex > -1 {
		logger.Info("Media asset already exists", "assetId", assets[mediaAssetIndex].Id)
		return nil, nil
	}

	logger.Info("Downloading media")
	paths, err := ytdlp.DownloadMedia(bookmark.Url)

	if err != nil {
		logger.Error("Failed to download media", "error", err)
		return nil, err
	}

	logger.Info("Media downloaded successfully", "paths", paths)
	return paths, nil
}

func uploadMedia(client *linkding.Client, bookmark linkding.Bookmark, paths []string, isDryRun bool) error {
	logger := slog.With("bookmarkId", bookmark.Id, "isDryRun", isDryRun)

	for _, path := range paths {
		logger.Info("Adding asset", "path", path)

		file, err := os.Open(path)

		if err != nil {
			logger.Error("Failed to open media file", "path", path, "error", err)
			return err
		}

		defer file.Close()
		defer os.Remove(path)

		asset, err := uploadAsset(client, bookmark, file, isDryRun)

		if err != nil {
			logger.Error("Failed to add asset", "path", path, "error", err)
			return err
		}

		logger.Info("Asset added successfully", "path", path, "assetId", asset.Id)
	}

	return nil
}

func uploadAsset(client *linkding.Client, bookmark linkding.Bookmark, file *os.File, isDryRun bool) (*linkding.Asset, error) {
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
