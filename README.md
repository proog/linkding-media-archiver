# Linkding Media Archiver

[![GitHub Release](https://img.shields.io/github/v/release/proog/linkding-media-archiver?sort=semver&style=for-the-badge&logo=github&label=Latest%20Release)](https://github.com/proog/linkding-media-archiver/releases)
[![Docker Image Version](https://img.shields.io/docker/v/proog/linkding-media-archiver?sort=semver&style=for-the-badge&logo=docker&label=Docker%20Hub)](https://hub.docker.com/r/proog/linkding-media-archiver)

Automatically download media for your [Linkding](https://linkding.link/) bookmarks!

## What it is

Linkding can automatically create HTML snapshots of your bookmarks to guard against link rot. Linkding Media Archiver supplements this feature by automatically downloading audio and video files for your bookmarks, such as YouTube videos or SoundCloud songs, and adding them to Linkding as additional assets.

## How it works

Linkding Media Archiver retrieves bookmarks that do not already have a media file attached and attempts to download one using [yt-dlp](https://github.com/yt-dlp/yt-dlp). If successful, the file is uploaded to Linkding as a bookmark asset. This process repeats on a configurable schedule with any bookmarks that have been added or changed since the previous run. If Linkding Media Archiver is restarted, it will retrieve all bookmarks again.

As yt-dlp is used to download media, [any site supported by yt-dlp](https://github.com/yt-dlp/yt-dlp/blob/master/supportedsites.md) should work. Please report a bug if Linkding Media Archiver fails to use a file that yt-dlp provides. yt-dlp's default format selection is used, which generally means the highest quality available in any file type, unless otherwise specified via the `LDMA_FORMAT` environment variable. Multiple files (such as YouTube playlists) are supported and will be added as multiple assets.

> [!WARNING]
> yt-dlp supports many arbitrary websites with its "generic extractor", which might cause Linkding Media Archiver to add media to unexpected bookmarks — for instance, a promotional video on a product landing page. For this reason, it is highly recommended to limit the bookmark selection to one or more tags using the `LDMA_TAGS` environment variable. For more advanced filtering, it is also possible to filter by [bundle](https://github.com/sissbruecker/linkding/pull/1097) with `LDMA_BUNDLE_ID`.

## Usage

Linkding Media Archiver requires [yt-dlp](https://github.com/yt-dlp/yt-dlp) and a [Linkding](https://linkding.link/) instance to work. The easiest way to run it is by using [the Docker image](https://hub.docker.com/r/proog/linkding-media-archiver), which includes yt-dlp. See `docker-compose.example.yml` for an example Docker Compose setup that combines Linkding and Linkding Media Archiver. Alternatively, it can be run as a binary by cloning the repository and compiling from source.

```sh
# Docker Compose (preferred, see docker-compose.example.yml)
docker compose up

# Docker
docker run --rm -e LDMA_BASEURL="http://localhost:9090" -e LDMA_TOKEN="abcd1234" proog/linkding-media-archiver [-n] [-s]

# Binary
go build -o ./linkding-media-archiver ./cmd
LDMA_BASEURL="http://localhost:9090" LDMA_TOKEN="abcd1234" ./linkding-media-archiver [-n] [-s]
```

### Flags

- `-n` Dry run: download media but do not actually upload it to Linkding
- `-s` Single run: exit after processing bookmarks once

### Environment variables

| Name                 | Example                            | Default                | Description                                                                                                                                         |
| -------------------- | ---------------------------------- | ---------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| `LDMA_BASEURL`       | `http://linkding.example.com:9090` | None **(required)**    | Base URL of your Linkding instance                                                                                                                  |
| `LDMA_TOKEN`         | `{random 40 char token}`           | None **(required)**    | Auth token from the Linkding integration page                                                                                                       |
| `LDMA_TAGS`          | `video music youtube`              | None (all bookmarks)   | Only process bookmarks with any of these tags (space separated, omit the #)                                                                         |
| `LDMA_BUNDLE_ID`     | `42`                               | None (all bookmarks)   | Only process bookmarks matching this [bundle](https://github.com/sissbruecker/linkding/pull/1097) (get the id from the url when editing the bundle) |
| `LDMA_FORMAT`        | `best[ext=mp4]`                    | None (yt-dlp defaults) | Format selection expression ([see yt-dlp docs](https://github.com/yt-dlp/yt-dlp?tab=readme-ov-file#format-selection))                               |
| `LDMA_SCAN_INTERVAL` | `600` (10 mins)                    | `3600` (1 hour)        | Schedule to check for new bookmarks                                                                                                                 |
| `LDMA_LOG_LEVEL`     | `DEBUG`                            | `INFO`                 | Log level, useful for troubleshooting                                                                                                               |
