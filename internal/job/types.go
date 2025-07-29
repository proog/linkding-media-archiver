package job

import "time"

type JobConfiguration struct {
	Tags     []string
	IsDryRun bool
	LastScan time.Time
}
