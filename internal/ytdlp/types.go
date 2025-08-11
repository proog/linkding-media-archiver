package ytdlp

type Ytdlp struct {
	DownloadDir string
	// MaxHeight limits the maximum video height (in pixels).
	// When 0 or less, no explicit height cap is applied.
	MaxHeight int
}
