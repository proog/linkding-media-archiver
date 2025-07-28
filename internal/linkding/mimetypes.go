package linkding

import (
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"strings"
)

var extensionMap = map[string]string{
	".3g2":  "video/3gpp2",
	".3gp":  "video/3gpp",
	".aac":  "audio/aac",
	".avi":  "video/x-msvideo",
	".flac": "audio/flac",
	".flv":  "video/x-flv",
	".m4a":  "audio/mp4",
	".m4v":  "video/x-m4v",
	".mkv":  "video/x-matroska",
	".mov":  "video/quicktime",
	".mp3":  "audio/mpeg",
	".mp4":  "video/mp4",
	".ogg":  "audio/ogg",
	".ogv":  "video/ogg",
	".opus": "audio/opus",
	".wav":  "audio/wav",
	".weba": "audio/webm",
	".webm": "video/webm",
	".wmv":  "video/x-ms-wmv",
}

var mimeTypes = slices.Compact(slices.Collect(maps.Values(extensionMap)))

func GetMimeType(fileName string) (mimeType string, err error) {
	ext := strings.ToLower(filepath.Ext(fileName))
	mimeType, ok := extensionMap[ext]

	if !ok {
		err = fmt.Errorf("unknown MIME type for %s", fileName)
	}

	return
}

func IsKnownMimeType(mimeType string) bool {
	return slices.Contains(mimeTypes, strings.ToLower(mimeType))
}
