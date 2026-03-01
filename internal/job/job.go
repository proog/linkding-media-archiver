package job

import (
	"linkding-media-archiver/internal/linkding"
	"linkding-media-archiver/internal/ytdlp"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

func ProcessBookmarks(client *linkding.Client, ytdlp *ytdlp.Ytdlp, config JobConfiguration) (err error) {
	logger := slog.With("tags", config.Tags, "bundleId", config.BundleId, "isDryRun", config.IsDryRun)

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
		hasAsset, err := hasMediaAsset(client, bookmark)
		if err != nil {
			failed <- bookmark
			continue
		}

		if hasAsset {
			succeeded <- bookmark
			continue
		}

		result, err := downloadMedia(ytdlp, bookmark)
		if err != nil {
			failed <- bookmark
			continue
		}

		wg.Go(func() {
			if err := uploadMedia(client, bookmark, result.Paths, config.IsDryRun); err != nil {
				failed <- bookmark
				return
			}

			if config.UpdateBookmarkText {
				if err := updateBookmark(client, bookmark, *result, config.IsDryRun); err != nil {
					failed <- bookmark
					return
				}
			}

			succeeded <- bookmark
		})
	}

	wg.Wait()

	logger.Info("Done processing bookmarks", "succeeded", len(succeeded), "failed", len(failed))

	return
}

func getBookmarks(client *linkding.Client, config JobConfiguration) ([]linkding.Bookmark, error) {
	query := linkding.BookmarksQuery{Tags: config.Tags, BundleId: config.BundleId, ModifiedSince: config.LastScan}
	return client.GetBookmarks(query)
}

func hasMediaAsset(client *linkding.Client, bookmark linkding.Bookmark) (bool, error) {
	logger := slog.With("bookmarkId", bookmark.Id)
	assets, err := client.GetBookmarkAssets(bookmark.Id)

	if err != nil {
		logger.Error("Failed to fetch bookmark assets")
		return false, err
	}

	mediaAssetIndex := slices.IndexFunc(assets, func(asset linkding.Asset) bool {
		return asset.AssetType == "upload" && linkding.IsKnownMimeType(asset.ContentType)
	})

	if mediaAssetIndex > -1 {
		logger.Info("Media asset already exists", "assetId", assets[mediaAssetIndex].Id)
		return true, nil
	}

	return false, nil
}

func downloadMedia(ytdlp *ytdlp.Ytdlp, bookmark linkding.Bookmark) (*ytdlp.DownloadResult, error) {
	logger := slog.With("bookmarkId", bookmark.Id)
	logger.Info("Downloading media")
	result, err := ytdlp.DownloadMedia(bookmark.Url)

	if err != nil {
		logger.Error("Failed to download media", "error", err)
		return nil, err
	}

	logger.Info("Media downloaded successfully", "result", result)
	return result, nil
}

func uploadMedia(client *linkding.Client, bookmark linkding.Bookmark, paths []string, isDryRun bool) error {
	logger := slog.With("bookmarkId", bookmark.Id, "isDryRun", isDryRun)

	for _, path := range paths {
		logger := logger.With("path", path)
		logger.Info("Adding asset")

		file, err := os.Open(path)

		if err != nil {
			logger.Error("Failed to open media file", "error", err)
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

func updateBookmark(client *linkding.Client, bookmark linkding.Bookmark, result ytdlp.DownloadResult, isDryRun bool) error {
	logger := slog.With("bookmarkId", bookmark.Id, "isDryRun", isDryRun)
	update := linkding.BookmarkUpdate{Title: bookmark.Title, Description: bookmark.Description}

	if strings.TrimSpace(result.Title) != "" {
		update.Title = result.Title
	}
	if strings.TrimSpace(result.Description) != "" {
		update.Description = result.Description
	}

	hasChanges := update.Title != bookmark.Title || update.Description != bookmark.Description
	if !hasChanges {
		logger.Info("Skipping bookmark update as there are no changes")
		return nil
	}

	logger.Info("Updating bookmark", "title", update.Title, "description", update.Description, "oldTitle", bookmark.Title, "oldDescription", bookmark.Description)

	if !isDryRun {
		if _, err := client.UpdateBookmark(bookmark.Id, update); err != nil {
			logger.Error("Failed to update bookmark", "error", err)
			return err
		}
	}

	logger.Info("Updated bookmark")
	return nil
}
