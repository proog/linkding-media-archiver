package ytdlp

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestDownloadMedia(t *testing.T) {
	ytdlp := Ytdlp{DownloadDir: t.TempDir()}
	results, err := ytdlp.DownloadMedia("https://www.youtube.com/watch?v=RWGTIIO2QiQ")

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected single result, was %d", len(results))
	}

	result := results[0]
	stat, err := os.Stat(result.Path)

	if err != nil {
		t.Fatal(err)
	}

	expectedName := "Tokyo_Station_JY-01_-_Yamanote_Train_Melody-[RWGTIIO2QiQ].webm"
	actualName := filepath.Base(stat.Name())

	if actualName != expectedName {
		t.Errorf("Unexpected filename: %s", actualName)
	}

	if result.Metadata.Title != "Tokyo Station (JY-01) - Yamanote Train Melody" {
		t.Errorf("Unexpected title: %s", result.Metadata.Title)
	}

	if !strings.HasPrefix(result.Metadata.Description, "The Tokyo station train melody for the JR Yamanote Line.") {
		t.Errorf("Unexpected description: %s", result.Metadata.Description)
	}

	expectedTags := []string{"tokyo station melody", "yamanote line jingle", "駅のメロディー"}
	for _, tag := range expectedTags {
		if !slices.Contains(result.Metadata.Tags, tag) {
			t.Errorf("Expected tag %s was not found in %s", tag, strings.Join(result.Metadata.Tags, ","))
		}
	}
}

func TestDownloadMediaWithFormatSelection(t *testing.T) {
	ytdlp := Ytdlp{DownloadDir: t.TempDir(), Format: "bestaudio[ext=m4a]"}
	results, err := ytdlp.DownloadMedia("https://www.youtube.com/watch?v=RWGTIIO2QiQ")

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected single result, was %d", len(results))
	}

	result := results[0]
	stat, err := os.Stat(result.Path)

	if err != nil {
		t.Fatal(err)
	}

	expectedName := "Tokyo_Station_JY-01_-_Yamanote_Train_Melody-[RWGTIIO2QiQ].m4a"
	actualName := filepath.Base(stat.Name())

	if actualName != expectedName {
		t.Errorf("Unexpected filename: %s", actualName)
	}
}

func TestDownloadMediaWithMultipleFiles(t *testing.T) {
	ytdlp := Ytdlp{DownloadDir: t.TempDir()}
	results, err := ytdlp.DownloadMedia("https://pontus.granstrom.me/scrappy/")

	if err != nil {
		t.Fatal(err)
	}

	if len(results) < 2 {
		t.Errorf("Expected multiple results, was %d", len(results))
	}

	for _, result := range results {
		if _, err := os.Stat(result.Path); err != nil {
			t.Error(err)
		}
	}
}
