package job

import "time"

type JobConfiguration struct {
	Tag      string
	IsDryRun bool
	LastScan time.Time
}
