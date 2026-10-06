package configuration

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func ReadConfiguration() Configuration {
	return Configuration{
		LinkdingBaseUrl:       getEnvOrFile("LDMA_BASEURL"),
		LinkdingToken:         getEnvOrFile("LDMA_TOKEN"),
		BundleId:              getLinkdingBundleId(),
		LogLevel:              getEnvOrFile("LDMA_LOG_LEVEL"),
		ScanInterval:          getScanInterval(),
		SkipExistingBookmarks: getSkipExistingBookmarks(),
		Tags:                  getLinkdingTags(),
		UpdateBookmarkText:    getUpdateBookmarkText(),
		YtdlpFormat:           getEnvOrFile("LDMA_FORMAT"),
	}
}

func getLinkdingTags() []string {
	tagsEnv := getEnvOrFile("LDMA_TAGS")
	return strings.Fields(tagsEnv)
}

func getLinkdingBundleId() int {
	bundleId, err := strconv.Atoi(getEnvOrFile("LDMA_BUNDLE_ID"))

	if bundleId <= 0 || err != nil {
		bundleId = 0
	}

	return bundleId
}

func getScanInterval() time.Duration {
	interval, err := strconv.Atoi(getEnvOrFile("LDMA_SCAN_INTERVAL"))

	if interval <= 0 || err != nil {
		interval = 3600
	}

	return time.Duration(interval) * time.Second
}

func getUpdateBookmarkText() bool {
	update, err := strconv.ParseBool(getEnvOrFile("LDMA_UPDATE_BOOKMARK_TEXT"))
	return err == nil && update
}

func getSkipExistingBookmarks() bool {
	skip, err := strconv.ParseBool(getEnvOrFile("LDMA_SKIP_EXISTING_BOOKMARKS"))
	return err == nil && skip
}

func getEnvOrFile(key string) string {
	if filePath := os.Getenv(key + "_FILE"); filePath != "" {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(content))
	}
	return os.Getenv(key)
}
