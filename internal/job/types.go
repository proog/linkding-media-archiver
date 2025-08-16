package job

import "time"

type JobConfiguration struct {
	Tags     []string
	BundleId int
	IsDryRun bool
	LastScan time.Time
}
