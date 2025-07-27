package linkding

import (
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"strings"
)

var extensionMap = map[string]string{
	".flv":  "video/x-flv",
	".mkv":  "video/x-matroska",
	".mp4":  "video/mp4",
	".webm": "video/webm",
}

var mimeTypes = slices.Collect(maps.Values(extensionMap))

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
