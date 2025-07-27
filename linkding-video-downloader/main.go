package main

import (
	"linkding-video-downloader/linkding"
	"linkding-video-downloader/logging"
	"linkding-video-downloader/ytdlp"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	logger := logging.NewLogger()
	slog.SetDefault(logger)

	client, err := linkding.NewClient(os.Getenv("LD_BASEURL"), os.Getenv("LD_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	tempdir, err := os.MkdirTemp(os.TempDir(), "videos")
	if err != nil {
		log.Fatal(err)
	}

	onExit(func() { os.RemoveAll(tempdir) })

	ytdlp := ytdlp.NewYtdlp(tempdir)
	tag := getLinkdingTag()
	interval := getScanInterval()
	sleep := time.Tick(time.Duration(interval) * time.Second)

	// Run immediately and then on every tick
	for ; true; <-sleep {
		err := processBookmarks(client, ytdlp, tag)

		if err != nil {
			slog.Error("Error processing bookmarks", "error", err)
		}

		slog.Info("Waiting for next scan", "intervalSeconds", interval)
	}
}

func getLinkdingTag() string {
	tag := os.Getenv("LD_TAG")

	if tag == "" {
		tag = "video"
	}

	return tag
}

func getScanInterval() int {
	intervalSeconds, err := strconv.Atoi(os.Getenv("SCAN_INTERVAL_SECONDS"))

	if intervalSeconds <= 0 || err != nil {
		intervalSeconds = 3600
	}

	return intervalSeconds
}

func onExit(cleanup func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		slog.Info("Exiting...")
		cleanup()
		os.Exit(0)
	}()
}
