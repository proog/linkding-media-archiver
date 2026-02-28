package main

import (
	"flag"
	"linkding-media-archiver/internal/configuration"
	"linkding-media-archiver/internal/job"
	"linkding-media-archiver/internal/linkding"
	"linkding-media-archiver/internal/logging"
	"linkding-media-archiver/internal/semver"
	"linkding-media-archiver/internal/ytdlp"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	isDryRun := flag.Bool("n", false, "Dry run: download media but do not actually upload it to Linkding")
	isSingleRun := flag.Bool("s", false, "Single run: exit after processing bookmarks once")
	flag.Parse()

	config := configuration.ReadConfiguration()

	logger := logging.NewLogger(config.LogLevel)
	slog.SetDefault(logger)

	client := createLinkdingClient(config)
	minVersion := semver.Semver{Major: 1, Minor: 44}
	checkLinkdingVersion(client, minVersion)

	tempdir := createTempDir()
	cleanupAndExit := func(code int) {
		os.RemoveAll(tempdir)
		os.Exit(code)
	}
	onInterrupt(cleanupAndExit)

	ytdlp := ytdlp.NewYtdlp(tempdir, config.YtdlpFormat)
	sleep := time.NewTicker(config.ScanInterval)

	var lastScan time.Time

	// Run immediately and then on every tick
	for ; true; <-sleep.C {
		timeBeforeRun := time.Now()

		if config.SkipExistingBookmarks && lastScan.IsZero() {
			logger.Info("Waiting for initial scan", "scanInterval", config.ScanInterval)
			lastScan = timeBeforeRun
			continue
		}

		jobConfig := job.JobConfiguration{
			Tags:     config.Tags,
			BundleId: config.BundleId,
			IsDryRun: *isDryRun,
			LastScan: lastScan,
		}
		err := job.ProcessBookmarks(client, ytdlp, jobConfig)

		if err == nil {
			lastScan = timeBeforeRun // Only update last scan time when bookmarks were actually processed
		} else {
			logger.Error("Error processing bookmarks", "error", err)
		}

		if *isSingleRun {
			cleanupAndExit(0)
		}

		logger.Info("Waiting for next scan", "scanInterval", config.ScanInterval)
	}
}

func createLinkdingClient(config configuration.Configuration) *linkding.Client {
	client, err := linkding.NewClient(config.LinkdingBaseUrl, config.LinkdingToken)

	if err != nil {
		log.Fatal(err)
	}

	return client
}

func checkLinkdingVersion(client *linkding.Client, minVersion semver.Semver) {
	profile, err := client.GetUserProfile()

	if err != nil {
		log.Fatal(err)
	}

	version, err := semver.Parse(profile.Version)

	if err != nil {
		log.Fatal(err)
	}

	if semver.Compare(version, minVersion) == -1 {
		log.Fatalf("Please upgrade Linkding: found version %s, but this program requires at least version %s", version, minVersion)
	}
}

func createTempDir() string {
	tempdir, err := os.MkdirTemp(os.TempDir(), "media")

	if err != nil {
		log.Fatal(err)
	}

	return tempdir
}

func onInterrupt(cleanup func(int)) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cleanup(1)
	}()
}
