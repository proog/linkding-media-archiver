package configuration

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func ReadConfiguration() Configuration {
	return Configuration{
		LinkdingBaseUrl:       os.Getenv("LDMA_BASEURL"),
		LinkdingToken:         os.Getenv("LDMA_TOKEN"),
		BundleId:              getLinkdingBundleId(),
		LogLevel:              os.Getenv("LDMA_LOG_LEVEL"),
		ScanInterval:          getScanInterval(),
		SkipExistingBookmarks: getSkipExistingBookmarks(),
		Tags:                  getLinkdingTags(),
		YtdlpFormat:           os.Getenv("LDMA_FORMAT"),
	}
}

func getLinkdingTags() []string {
	tagsEnv := os.Getenv("LDMA_TAGS")
	return strings.Fields(tagsEnv)
}

func getLinkdingBundleId() int {
	bundleId, err := strconv.Atoi(os.Getenv("LDMA_BUNDLE_ID"))

	if bundleId <= 0 || err != nil {
		bundleId = 0
	}

	return bundleId
}

func getScanInterval() time.Duration {
	interval, err := strconv.Atoi(os.Getenv("LDMA_SCAN_INTERVAL"))

	if interval <= 0 || err != nil {
		interval = 3600
	}

	return time.Duration(interval) * time.Second
}

func getSkipExistingBookmarks() bool {
	skip, err := strconv.ParseBool(os.Getenv("LDMA_SKIP_EXISTING_BOOKMARKS"))
	return err == nil && skip
}
