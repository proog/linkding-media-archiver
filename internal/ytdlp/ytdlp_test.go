package ytdlp

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestDownloadMedia(t *testing.T) {
	ytdlp := NewYtdlp(t.TempDir(), "")
	result, err := ytdlp.DownloadMedia("https://www.youtube.com/watch?v=RWGTIIO2QiQ")

	if err != nil {
		t.Fatal(err)
	}

	if len(result.Paths) != 1 {
		t.Fatalf("Expected single file, was %d", len(result.Paths))
	}

	stat, err := os.Stat(result.Paths[0])

	if err != nil {
		t.Fatal(err)
	}

	expectedName := "Tokyo_Station_JY-01_-_Yamanote_Train_Melody-[RWGTIIO2QiQ].webm"
	actualName := filepath.Base(stat.Name())

	if actualName != expectedName {
		t.Errorf("Unexpected filename: %s", actualName)
	}

	if result.Title != "Tokyo Station (JY-01) - Yamanote Train Melody" {
		t.Errorf("Unexpected title: %s", result.Title)
	}

	if !strings.HasPrefix(result.Description, "The Tokyo station train melody for the JR Yamanote Line.") {
		t.Errorf("Unexpected description: %s", result.Description)
	}

	expectedTags := []string{"tokyo station melody", "yamanote line jingle", "駅のメロディー"}
	for _, tag := range expectedTags {
		if !slices.Contains(result.Tags, tag) {
			t.Errorf("Expected tag %s was not found in %s", tag, strings.Join(result.Tags, ","))
		}
	}
}

func TestDownloadMediaWithFormatSelection(t *testing.T) {
	ytdlp := NewYtdlp(t.TempDir(), "bestaudio[ext=m4a]")
	result, err := ytdlp.DownloadMedia("https://www.youtube.com/watch?v=RWGTIIO2QiQ")

	if err != nil {
		t.Fatal(err)
	}

	if len(result.Paths) != 1 {
		t.Fatalf("Expected single file, was %d", len(result.Paths))
	}

	stat, err := os.Stat(result.Paths[0])

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
	ytdlp := NewYtdlp(t.TempDir(), "")
	result, err := ytdlp.DownloadMedia("https://www.youtube.com/playlist?list=PLSBoMdEkRnhQCyNGzVR66TgY93bJcTfsc")

	if err != nil {
		t.Fatal(err)
	}

	if len(result.Paths) < 2 {
		t.Errorf("Expected multiple files, was %d", len(result.Paths))
	}

	if result.Title != "Weslace and Zromitman" {
		t.Errorf("Unexpected title: %s", result.Title)
	}

	if !strings.HasPrefix(result.Description, "The most pleasant of pleasant days begin with a nice hearty breakfast to put that spring in your step.") {
		t.Errorf("Unexpected description: %s", result.Description)
	}

	for _, path := range result.Paths {
		if _, err := os.Stat(path); err != nil {
			t.Error(err)
		}
	}
}
