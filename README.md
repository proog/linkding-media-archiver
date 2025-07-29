# Linkding Media Archiver

Automatically download media for your [Linkding](https://linkding.link/) bookmarks!

## What it is

Linkding can automatically create HTML snapshots of your bookmarks to guard against link rot. Linkding Media Archiver supplements this feature by automatically downloading audio and video files for your bookmarks, such as YouTube videos or SoundCloud songs, and adding them to Linkding as additional assets.

## How it works

Linkding Media Archiver retrieves all bookmarks with a specific tag that do not already have a media file attached and attempts to download one using [yt-dlp](https://github.com/yt-dlp/yt-dlp). If successful, the file is uploaded to Linkding as a bookmark asset. This process repeats on a configurable schedule with any bookmarks that have been added or changed since the previous run. If Linkding Media Archiver is restarted, it will retrieve all bookmarks again.

As yt-dlp is used to download media, [any site supported by yt-dlp](https://github.com/yt-dlp/yt-dlp/blob/master/supportedsites.md) should work. Please report a bug if Linkding Media Archiver fails to add a file that yt-dlp provides.

## Limitations

- Downloading multiple files per bookmark (such as YouTube playlists) is not supported
- yt-dlp uses default format and quality settings (i.e. the highest quality available and no fixed media formats)

## Usage

```sh
./linkding-media-archiver [-n] [-s]
```

### Flags

- `-n` Dry run: download media but do not actually upload it to Linkding
- `-s` Single run: exit after processing bookmarks once

### Environment variables

| Name                    | Example                            | Default         | Description                                   |
| ----------------------- | ---------------------------------- | --------------- | --------------------------------------------- |
| `LD_BASEURL`            | `http://linkding.example.com:9090` | None            | Base URL of your Linkding instance            |
| `LD_TOKEN`              | `{random 40 char token}`           | None            | Auth token from the Linkding integration page |
| `LD_TAG`                | `media`                            | `video`         | Process bookmarks with this tag (omit the #)  |
| `SCAN_INTERVAL_SECONDS` | `600` (10 mins)                    | `3600` (1 hour) | Schedule to checking for new bookmarks        |
| `LOG_LEVEL`             | `DEBUG`                            | `INFO`          | Log level, useful for troubleshooting         |
