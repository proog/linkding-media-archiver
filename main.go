package main

import (
	"flag"
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

	isDryRun := flag.Bool("n", false, "Dry run: download videos but do not actually upload them to Linkding")
	isSingleRun := flag.Bool("s", false, "Single run: exit after processing bookmarks instead")
	flag.Parse()

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

	cleanupAndExit := func(code int) {
		os.RemoveAll(tempdir)
		os.Exit(code)
	}

	onInterrupt(cleanupAndExit)

	ytdlp := ytdlp.NewYtdlp(tempdir)
	tag := getLinkdingTag()
	interval := getScanInterval()
	sleep := time.NewTicker(time.Duration(interval) * time.Second)

	var lastScan time.Time

	// Run immediately and then on every tick
	for ; true; <-sleep.C {
		query := linkding.BookmarksQuery{Tag: tag, ModifiedSince: lastScan}
		err := processBookmarks(client, ytdlp, &query, *isDryRun)

		if err != nil {
			slog.Error("Error processing bookmarks", "error", err)
		}

		if *isSingleRun {
			cleanupAndExit(0)
		}

		lastScan = time.Now()
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

func onInterrupt(cleanup func(int)) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cleanup(1)
	}()
}
