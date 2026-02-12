# SpotiFLAC CLI Tool Documentation

The SpotiFLAC CLI tool provides headless access to all SpotiFLAC functionality for server environments.

## Installation

The CLI binary is built as part of the deployment process:

```bash
go build -o spotiflac cmd/cli/main.go
```

Or using the `launch.sh` script which builds both server and CLI.

## Configuration

The CLI tool reads all settings from `config.yml`. All download preferences, service settings, and paths configured in the config file will be respected.

### Configuration Location

By default, the CLI looks for `config.yml` in the current directory. You can specify a custom location:

```bash
spotiflac --config /path/to/config.yml download track [url]
```

## Commands

### Download Commands

#### Download Single Track

```bash
spotiflac download track <spotify-url>
```

Example:
```bash
spotiflac download track https://open.spotify.com/track/abc123
```

#### Download Album

```bash
spotiflac download album <spotify-url>
```

Example:
```bash
spotiflac download album https://open.spotify.com/album/xyz789
```

#### Download Playlist

```bash
spotiflac download playlist <spotify-url>
```

Example:
```bash
spotiflac download playlist https://open.spotify.com/playlist/def456
```

### Configuration Commands

#### Get Configuration Value

```bash
spotiflac config get <key>
```

Examples:
```bash
spotiflac config get download.path
spotiflac config get services.default_service
spotiflac config get download.audio_format
```

#### Set Configuration Value

```bash
spotiflac config set <key> <value>
```

Examples:
```bash
spotiflac config set download.path /home/user/Music
spotiflac config set services.default_service tidal
spotiflac config set download.audio_format LOSSLESS
```

## Global Flags

- `--config, -c`: Specify config file path (default: `config.yml`)
- `--json, -j`: Output in JSON format (useful for scripting)

## JSON Output Mode

For scripting and automation, use JSON output mode:

```bash
spotiflac --json download track https://open.spotify.com/track/abc123
```

Output:
```json
{"success": true, "track": "Song Name", "spotify_id": "abc123"}
```

## Examples

### Batch Download from File

```bash
#!/bin/bash
while read -r url; do
    spotiflac download track "$url"
done < spotify_urls.txt
```

### Check Configuration Before Download

```bash
DOWNLOAD_PATH=$(spotiflac config get download.path)
echo "Downloads will be saved to: $DOWNLOAD_PATH"
spotiflac download playlist https://open.spotify.com/playlist/xyz
```

### Change Service and Download

```bash
spotiflac config set services.default_service qobuz
spotiflac download track https://open.spotify.com/track/abc123
```

## Respecting Frontend Settings

The CLI tool fully respects all settings configured via the web frontend:

- **Download Path**: Files saved to configured directory
- **Audio Format**: LOSSLESS, or specific quality levels
- **Filename Format**: Follows configured naming pattern
- **Embed Lyrics**: Automatically embeds if enabled
- **Cover Art**: Embeds maximum quality if configured
- **Default Service**: Uses configured streaming service (Tidal/Qobuz/Amazon)

## Progress Indicators

Downloads show a progress bar when not in JSON mode:

```
Fetching metadata for: https://open.spotify.com/track/abc123
Track: Song Name - Artist Name
Starting download...
Downloading ██████████████████████████ 100%
Download complete: Song Name
```

## Tips

1. **Use absolute paths** when setting download paths from CLI
2. **Test with JSON mode** to ensure commands work before scripting
3. **Check configuration** before batch operations with `spotiflac config get`
4. The CLI shares the same database as the web frontend, so download history is unified
