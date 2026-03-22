package ytdlp

type Ytdlp struct {
	DownloadDir string
	Format      string
}

type DownloadResult struct {
	Title       string
	Description string
	Tags        []string
	Paths       []string
}

type jsonDump struct {
	Title              string              `json:"title"`
	Description        string              `json:"description"`
	Tags               []string            `json:"tags"`
	Entries            []dumpEntry         `json:"entries"`
	RequestedDownloads []requestedDownload `json:"requested_downloads"`
}

type dumpEntry struct {
	Title              string              `json:"title"`
	Description        string              `json:"description"`
	Tags               []string            `json:"tags"`
	RequestedDownloads []requestedDownload `json:"requested_downloads"`
}

type requestedDownload struct {
	FilePath string `json:"filepath"`
}
