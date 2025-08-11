package ytdlp

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadMedia(t *testing.T) {
	ytdlp := Ytdlp{DownloadDir: t.TempDir()}
	paths, err := ytdlp.DownloadMedia("https://www.youtube.com/watch?v=RWGTIIO2QiQ")

	if err != nil {
		t.Fatal(err)
	}

	if len(paths) != 1 {
		t.Fatalf("Expected single file, was %d", len(paths))
	}

	stat, err := os.Stat(paths[0])

	if err != nil {
		t.Fatal(err)
	}

	expectedName := "Tokyo_Station_JY-01_-_Yamanote_Train_Melody-[RWGTIIO2QiQ].webm"
	actualName := filepath.Base(stat.Name())

	if actualName != expectedName {
		t.Fatalf("Unexpected filename: %s", actualName)
	}
}

func TestDownloadMediaWithFormatSelection(t *testing.T) {
	ytdlp := Ytdlp{DownloadDir: t.TempDir(), Format: "bestaudio[ext=m4a]"}
	paths, err := ytdlp.DownloadMedia("https://www.youtube.com/watch?v=RWGTIIO2QiQ")

	if err != nil {
		t.Fatal(err)
	}

	if len(paths) != 1 {
		t.Fatalf("Expected single file, was %d", len(paths))
	}

	stat, err := os.Stat(paths[0])

	if err != nil {
		t.Fatal(err)
	}

	expectedName := "Tokyo_Station_JY-01_-_Yamanote_Train_Melody-[RWGTIIO2QiQ].m4a"
	actualName := filepath.Base(stat.Name())

	if actualName != expectedName {
		t.Fatalf("Unexpected filename: %s", actualName)
	}
}

func TestDownloadMediaWithMultipleFiles(t *testing.T) {
	ytdlp := Ytdlp{DownloadDir: t.TempDir()}
	paths, err := ytdlp.DownloadMedia("https://pontus.granstrom.me/scrappy/")

	if err != nil {
		t.Fatal(err)
	}

	if len(paths) < 2 {
		t.Fatalf("Expected multiple files, was %d", len(paths))
	}

	for _, path := range paths {
		_, err := os.Stat(path)

		if err != nil {
			t.Fatal(err)
		}
	}
}
