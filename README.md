# Linkding Media Archiver

Automatically archive media bookmarked in Linkding using [yt-dlp](https://github.com/yt-dlp/yt-dlp)

## Usage

```sh
./linkding-media-archiver [-n] [-s]
```

### Flags

- `-n` Dry run: download media but do not actually upload it to Linkding
- `-s` Single run: exit after processing bookmarks once

## Environment variables

| Name                    | Example                            | Default         | Description                                   |
| ----------------------- | ---------------------------------- | --------------- | --------------------------------------------- |
| `LD_BASEURL`            | `http://linkding.example.com:9090` | None            | Base URL of your Linkding instance            |
| `LD_TOKEN`              | `{random 40 char token}`           | None            | Auth token from the Linkding integration page |
| `LD_TAG`                | `media`                            | `video`         | Process bookmarks with this tag (omit the #)  |
| `SCAN_INTERVAL_SECONDS` | `600` (10 mins)                    | `3600` (1 hour) | Schedule to checking for new bookmarks        |
| `LOG_LEVEL`             | `DEBUG`                            | `INFO`          | Log level, useful for troubleshooting         |
