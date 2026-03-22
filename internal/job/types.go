package job

import "time"

type JobConfiguration struct {
	Tags               []string
	BundleId           int
	UpdateBookmarkText bool
	IsDryRun           bool
	LastScan           time.Time
}
