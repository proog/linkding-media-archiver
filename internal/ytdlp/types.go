package ytdlp

type Ytdlp struct {
	DownloadDir string
	Format      string
}

type DownloadResult struct {
	Path     string
	Metadata Metadata
}

type Metadata struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}
