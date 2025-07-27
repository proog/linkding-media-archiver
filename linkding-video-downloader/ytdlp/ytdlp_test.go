package ytdlp

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadVideo(t *testing.T) {
	ytdlp := Ytdlp{DownloadDir: t.TempDir()}
	path, err := ytdlp.DownloadVideo("https://www.youtube.com/watch?v=RWGTIIO2QiQ")

	if err != nil {
		t.Fatal(err)
	}

	stat, err := os.Stat(path)

	if err != nil {
		t.Fatal(err)
	}

	expectedName := "Tokyo_Station_JY-01_-_Yamanote_Train_Melody-[RWGTIIO2QiQ].webm"
	actualName := filepath.Base(stat.Name())

	if actualName != expectedName {
		t.Fatalf("Unexpected filename: %s", actualName)
	}
}
