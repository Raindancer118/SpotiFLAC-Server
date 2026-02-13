# Docker Deployment Guide

## How SpotiFLAC Works

**Important:** SpotiFLAC does **NOT require Spotify API credentials**! 

The app works by:
1. **Input:** You provide a Spotify URL (track, album, or playlist)
2. **Metadata:** App fetches metadata from Spotify (public API, no credentials needed)
3. **Download:** Music is downloaded from **Tidal, Qobuz, or Amazon Music** (via public APIs)
4. **Tagging:** Downloaded files are tagged with Spotify metadata

**Supported streaming services for downloads:**
- 🎵 **Tidal** (Default, FLAC quality)
- 🎵 **Qobuz** (Hi-Res FLAC)
- 🎵 **Amazon Music** (Fallback)

## Quick Start

1. **Start with Docker Compose:**
   ```bash
   docker-compose up -d
   ```

2. **Access the web interface:**
   ```
   http://localhost:8080
   ```

That's it! No configuration needed for basic usage.

## Advanced Configuration (Optional)

If you want to customize settings, create a `.env` file:

```bash
cp .env.example .env
nano .env
```

## Docker Build

### Using Docker Compose (Recommended)

```bash
# Build and start
docker-compose up -d

# View logs
docker-compose logs -f spotiflac

# Stop
docker-compose down

# Rebuild after code changes
docker-compose up -d --build
```

### Using Docker CLI

```bash
# Build image
docker build -t spotiflac:latest .

# Run container
docker run -d \
  --name spotiflac \
  -p 8080:8080 \
  -v $(pwd)/downloads:/app/downloads \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/config.yml:/app/config.yml:ro \
  -e SPOTIFY_CLIENT_ID=your_client_id \
  -e SPOTIFY_CLIENT_SECRET=your_client_secret \
  spotiflac:latest

# View logs
docker logs -f spotiflac

# Stop container
docker stop spotiflac
docker rm spotiflac
```

## Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SERVER_PORT` | HTTP server port | `8080` | No |
| `WEB_DIR` | Frontend files directory | `/app/web` | No |
| `DATABASE_PATH` | SQLite database path | `/app/data/spotiflac.db` | No |
| `LOG_LEVEL` | Logging level (debug/info/warn/error) | `info` | No |
| `TZ` | Timezone | `UTC` | No |
| `TIDAL_API_URL` | Custom Tidal API endpoint | Public API | No |
| `QOBUZ_API_URL` | Custom Qobuz API endpoint | Public API | No |
| `AMAZON_API_URL` | Custom Amazon Music API endpoint | Public API | No |

### Volumes

| Host Path | Container Path | Purpose |
|-----------|---------------|---------|
| `./downloads` | `/app/downloads` | Downloaded music files |
| `./data` | `/app/data` | Database and persistent data |
| `./config.yml` | `/app/config.yml` | Optional configuration file |

## Using with Nginx

1. **Uncomment Nginx service** in `docker-compose.yml`:
   ```yaml
   # Remove the profiles line under nginx service
   ```

2. **Configure SSL certificates** (optional):
   ```bash
   mkdir -p ssl
   # Copy your certificates to ./ssl/
   ```

3. **Edit** `nginx.conf.example` and save as `nginx.conf`

4. **Start with Nginx:**
   ```bash
   docker-compose up -d
   ```

Access via:
- HTTP: `http://localhost`
- HTTPS: `https://localhost` (if SSL configured)

## CLI Usage in Docker

Run CLI commands inside the container:

```bash
# Download a track
docker exec spotiflac-server /app/spotiflac-cli download \
  "https://open.spotify.com/track/xxxxx"

# Download a playlist
docker exec spotiflac-server /app/spotiflac-cli download \
  "https://open.spotify.com/playlist/xxxxx" \
  --output /app/downloads/my-playlist

# Check version
docker exec spotiflac-server /app/spotiflac-cli version
```

## Health Check

The container includes a health check that runs every 30 seconds:

```bash
# Check health status
docker inspect --format='{{.State.Health.Status}}' spotiflac-server

# View health check logs
docker inspect --format='{{range .State.Health.Log}}{{.Output}}{{end}}' spotiflac-server
```

## Resource Limits

Default resource limits (adjust in `docker-compose.yml`):
- **CPU Limit:** 2 cores
- **Memory Limit:** 2 GB
- **CPU Reservation:** 0.5 cores
- **Memory Reservation:** 512 MB

## Troubleshooting

### Container won't start

```bash
# Check logs
docker-compose logs spotiflac

# Check container status
docker-compose ps
```

### Permission issues with volumes

```bash
# Fix permissions
sudo chown -R 1000:1000 downloads/ data/
```

### FFmpeg not found

FFmpeg is included in the Docker image. If you see errors:

```bash
# Verify FFmpeg is available
docker exec spotiflac-server ffmpeg -version
```

### Database locked errors

Stop all containers and restart:

```bash
docker-compose down
docker-compose up -d
```

## Updating

```bash
# Pull latest code
git pull

# Rebuild and restart
docker-compose down
docker-compose up -d --build
```

## Production Deployment

For production use:

1. **Use a reverse proxy** (Nginx/Traefik) with SSL
2. **Set proper resource limits** in `docker-compose.yml`
3. **Configure log rotation**
4. **Set up automated backups** for `/app/data`
5. **Use Docker secrets** for sensitive credentials
6. **Enable firewall** and restrict access as needed

### Example with Traefik

```yaml
services:
  spotiflac:
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.spotiflac.rule=Host(`music.example.com`)"
      - "traefik.http.routers.spotiflac.entrypoints=websecure"
      - "traefik.http.routers.spotiflac.tls.certresolver=letsencrypt"
```

## Security Notes

⚠️ **Important:**
- Never commit `.env` file to version control
- Use strong, unique Spotify API credentials
- Run container as non-root user (default: `spotiflac`)
- Keep Docker and images up to date
- Restrict network access with firewall rules
- Use HTTPS in production

## Multi-Architecture Support

To build for different architectures:

```bash
# Build for ARM64 (e.g., Raspberry Pi)
docker buildx build --platform linux/arm64 -t spotiflac:arm64 .

# Build for multiple platforms
docker buildx build --platform linux/amd64,linux/arm64 -t spotiflac:latest .
```
