package main

import (
	"flag"
	"linkding-media-archiver/internal/job"
	"linkding-media-archiver/internal/linkding"
	"linkding-media-archiver/internal/logging"
	"linkding-media-archiver/internal/ytdlp"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	isDryRun := flag.Bool("n", false, "Dry run: download media but do not actually upload it to Linkding")
	isSingleRun := flag.Bool("s", false, "Single run: exit after processing bookmarks once")
	flag.Parse()

	logger := logging.NewLogger()
	slog.SetDefault(logger)

	client, err := linkding.NewClient(os.Getenv("LD_BASEURL"), os.Getenv("LD_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	tempdir, err := os.MkdirTemp(os.TempDir(), "media")
	if err != nil {
		log.Fatal(err)
	}

	cleanupAndExit := func(code int) {
		os.RemoveAll(tempdir)
		os.Exit(code)
	}

	onInterrupt(cleanupAndExit)

	ytdlp := ytdlp.NewYtdlp(tempdir)
	tags := getLinkdingTags()
	interval := getScanInterval()
	sleep := time.NewTicker(time.Duration(interval) * time.Second)

	var lastScan time.Time

	// Run immediately and then on every tick
	for ; true; <-sleep.C {
		config := job.JobConfiguration{Tags: tags, IsDryRun: *isDryRun, LastScan: lastScan}
		err := job.ProcessBookmarks(client, ytdlp, config)

		if err != nil {
			logger.Error("Error processing bookmarks", "error", err)
		}

		if *isSingleRun {
			cleanupAndExit(0)
		}

		lastScan = time.Now()
		logger.Info("Waiting for next scan", "intervalSeconds", interval)
	}
}

func getLinkdingTags() []string {
	tagsEnv := os.Getenv("LD_TAGS")
	return strings.Fields(tagsEnv)
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
