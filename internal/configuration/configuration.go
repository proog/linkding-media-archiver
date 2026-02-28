package configuration

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func ReadConfiguration() Configuration {
	return Configuration{
		LinkdingBaseUrl:       requireEnv("LDMA_BASEURL"),
		LinkdingToken:         requireEnv("LDMA_TOKEN"),
		BundleId:              getLinkdingBundleId(),
		LogLevel:              getLogLevel(),
		ScanInterval:          getScanInterval(),
		SkipExistingBookmarks: getSkipExistingBookmarks(),
		Tags:                  getLinkdingTags(),
		YtdlpFormat:           getYtdlpFormat(),
	}
}

func requireEnv(name string) string {
	value := os.Getenv(name)

	if value == "" {
		log.Fatalf("Environment variable %s is required", name)
	}

	return value
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

func getYtdlpFormat() string {
	return os.Getenv("LDMA_FORMAT")
}

func getLogLevel() string {
	return os.Getenv("LDMA_LOG_LEVEL")
}
