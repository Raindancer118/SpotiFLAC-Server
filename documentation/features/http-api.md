# SpotiFLAC HTTP API Documentation

The SpotiFLAC HTTP API provides access to all application functionality via RESTful endpoints.

## Base URL

```
http://localhost:8080/api
```

## Authentication

Currently, the API does not require authentication. For production deployments, consider adding authentication via reverse proxy or middleware.

## Endpoints

### Health Check

#### GET /health

Check server health status.

**Response:**
```json
{
  "status": "healthy",
  "time": 1708000000
}
```

---

### Spotify Metadata

#### POST /api/spotify/metadata

Fetch metadata from a Spotify URL.

**Request:**
```json
{
  "url": "https://open.spotify.com/track/abc123",
  "batch": false,
  "delay": 1.0,
  "timeout": 300.0
}
```

**Response:**
```json
{
  "track": {
    "name": "Song Title",
    "artists": "Artist Name",
    "album_name": "Album Name",
    "spotify_id": "abc123",
    "duration": 240,
    "images": "https://...",
    ...
  }
}
```

#### POST /api/spotify/search

Search Spotify for tracks, albums, artists, or playlists.

**Request:**
```json
{
  "query": "search query",
  "limit": 10
}
```

#### POST /api/spotify/search-by-type

Search Spotify by specific type.

**Request:**
```json
{
  "query": "search query",
  "search_type": "track",
  "limit": 50,
  "offset": 0
}
```

Types: `track`, `album`, `artist`, `playlist`

#### POST /api/spotify/streaming-urls

Get streaming URLs for a Spotify track.

**Request:**
```json
{
  "spotify_track_id": "abc123",
  "region": "US"
}
```

---

### Downloads

#### POST /api/download/track

Download a track.

**Request:**
```json
{
  "service": "tidal",
  "track_name": "Song Title",
  "artist_name": "Artist Name",
  "album_name": "Album Name",
  "spotify_id": "abc123",
  "output_dir": "/path/to/music",
  "audio_format": "LOSSLESS",
  "filename_format": "title-artist",
  "embed_lyrics": true,
  ...
}
```

**Response:**
```json
{
  "success": true,
  "message": "Download completed successfully",
  "file": "/path/to/music/Song Title - Artist Name.flac",
  "already_exists": false,
  "item_id": "abc123-1708000000"
}
```

#### GET /api/download/queue

Get current download queue status.

**Response:**
```json
{
  "queue": [
    {
      "item_id": "abc123-1708000000",
      "track_name": "Song Title",
      "artist_name": "Artist Name",
      "status": "downloading",
      "progress": 45.5,
      ...
    }
  ],
  "is_downloading": true
}
```

#### GET /api/download/progress

Get download progress for current item.

**Response:**
```json
{
  "is_downloading": true,
  "mb_downloaded": 15.5,
  "mb_total": 35.2,
  "percentage": 44.0,
  "speed_mbps": 2.5,
  ...
}
```

#### POST /api/download/queue/clear

Clear completed downloads from queue.

#### POST /api/download/queue/clear-all

Clear all items from download queue.

#### POST /api/download/queue/cancel-all

Cancel all queued items.

--- 

### History

#### GET /api/history/downloads

Get download history.

**Response:**
```json
[
  {
    "id": "uuid",
    "spotify_id": "abc123",
    "title": "Song Title",
    "artists": "Artist Name",
    "album": "Album Name",
    "duration_str": "4:00",
    "cover_url": "https://...",
    "quality": "24-bit/96.0kHz",
    "format": "FLAC",
    "path": "/path/to/file.flac",
    "timestamp": 1708000000
  }
]
```

#### POST /api/history/downloads/clear

Clear download history.

#### DELETE /api/history/downloads/:id

Delete specific history item.

---

### Settings

#### GET /api/settings

Get current settings.

**Response:**
```json
{
  "downloadPath": "/home/user/Music",
  "filenameFormat": "title-artist",
  "audioFormat": "LOSSLESS",
  "embedLyrics": true,
  "embedMaxQualityCover": true,
  "defaultService": "tidal",
  "theme": "default",
  "themeMode": "dark",
  ...
}
```

#### POST /api/settings

Update settings.

**Request:**
```json
{
  "downloadPath": "/new/path",
  "audioFormat": "LOSSLESS",
  "embedLyrics": true,
  ...
}
```

**Response:**
```json
{
  "success": true
}
```

#### GET /api/defaults

Get default system values.

**Response:**
```json
{
  "downloadPath": "/home/user/Music"
}
```

---

### System

#### GET /api/system/ffmpeg/status

Check if FFmpeg is installed.

**Response:**
```json
{
  "installed": true
}
```

---

### Analysis

#### POST /api/analysis/track

Analyze an audio file.

**Request:**
```json
{
  "file_path": "/path/to/file.flac"
}
```

**Response:**
```json
{
  "sample_rate": 96000,
  "bits_per_sample": 24,
  "duration": 240.5,
  "bitrate": 2500,
  ...
}
```

---

## WebSocket

### Endpoint: /ws

Connect to WebSocket for real-time updates.

**Message Types:**

```json
{
  "type": "download_progress",
  "data": {
    "is_downloading": true,
    "percentage": 45.5,
    ...
  }
}
```

```json
{
  "type": "queue_update",
  "data": {
    "queue": [...],
    "is_downloading": true
  }
}
```

**Client Messages:**

```json
{
  "type": "ping"
}
```

```json
{
  "type": "request_status"
}
```

---

## Error Responses

Errors return appropriate HTTP status codes with error messages:

```json
{
  "error": "Error description"
}
```

Common status codes:
- `400 Bad Request` - Invalid request format or parameters
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error occurred

---

## CORS

The API supports CORS. Allowed origins are configured in `config.yml`:

```yaml
server:
  cors_origins:
    - "http://localhost:5173"
    - "http://localhost:8080"
```

---

## Rate Limiting

Currently, no rate limiting is implemented. For production, consider adding rate limiting at the reverse proxy level.
