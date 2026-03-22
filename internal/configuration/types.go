package configuration

import "time"

type Configuration struct {
	LinkdingBaseUrl       string
	LinkdingToken         string
	BundleId              int
	LogLevel              string
	ScanInterval          time.Duration
	SkipExistingBookmarks bool
	Tags                  []string
	UpdateBookmarkText    bool
	YtdlpFormat           string
}
