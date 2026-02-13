# Docker Image Publishing & Deployment

## Publishing Options

SpotiFLAC supports two ways to distribute Docker images:

### Option 1: GitHub Container Registry (GHCR) - Recommended

**Advantages:**
- ✅ Free for public repositories
- ✅ Automatic builds via GitHub Actions
- ✅ Multi-platform support (amd64, arm64)
- ✅ Integrated with GitHub repo
- ✅ No separate account needed

**Setup:**

1. **Enable GitHub Actions** (already configured in `.github/workflows/docker-publish.yml`)

2. **Push to main branch** - Image builds automatically:
   ```bash
   git push origin main
   ```

3. **Image is published to:**
   ```
   ghcr.io/raindancer118/spotiflac-server:latest
   ghcr.io/raindancer118/spotiflac-server:main
   ghcr.io/raindancer118/spotiflac-server:sha-<commit>
   ```

4. **For tagged releases:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
   Creates additional tags:
   - `ghcr.io/raindancer118/spotiflac-server:v1.0.0`
   - `ghcr.io/raindancer118/spotiflac-server:1.0.0`
   - `ghcr.io/raindancer118/spotiflac-server:1.0`
   - `ghcr.io/raindancer118/spotiflac-server:1`

### Option 2: Docker Hub

**Advantages:**
- ✅ More widely known
- ✅ Better discoverability

**Setup:**

1. **Create Docker Hub account** at https://hub.docker.com

2. **Login locally:**
   ```bash
   docker login
   ```

3. **Build and tag:**
   ```bash
   docker build -t yourusername/spotiflac:latest .
   ```

4. **Push:**
   ```bash
   docker push yourusername/spotiflac:latest
   ```

## Pulling on Server

### From GitHub Container Registry (Public)

```bash
# Pull latest version
docker pull ghcr.io/raindancer118/spotiflac-server:latest

# Pull specific version
docker pull ghcr.io/raindancer118/spotiflac-server:v1.0.0

# Using docker-compose
docker-compose pull
docker-compose up -d
```

### From GitHub Container Registry (Private Repo)

1. **Create Personal Access Token (PAT):**
   - Go to: GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
   - Create token with `read:packages` scope

2. **Login on server:**
   ```bash
   echo $PAT_TOKEN | docker login ghcr.io -u YOUR_GITHUB_USERNAME --password-stdin
   ```

3. **Pull image:**
   ```bash
   docker pull ghcr.io/raindancer118/spotiflac-server:latest
   ```

### From Docker Hub

```bash
docker pull yourusername/spotiflac:latest
```

## Updated docker-compose.yml for Registry

Instead of building locally, pull from registry:

```yaml
services:
  spotiflac:
    image: ghcr.io/raindancer118/spotiflac-server:latest
    # Remove 'build' section
    container_name: spotiflac-server
    restart: unless-stopped
    # ... rest of config
```

## Server Deployment Workflow

### Initial Setup on Server

```bash
# 1. Create deployment directory
mkdir -p /opt/spotiflac
cd /opt/spotiflac

# 2. Create docker-compose.yml (using registry image)
cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  spotiflac:
    image: ghcr.io/raindancer118/spotiflac-server:latest
    container_name: spotiflac-server
    restart: unless-stopped
    
    ports:
      - "8080:8080"
    
    volumes:
      - ./downloads:/app/downloads
      - ./data:/app/data
    
    environment:
      - SERVER_PORT=8080
      - TZ=Europe/Berlin
EOF

# 3. Start
docker-compose up -d
```

### Updating to Latest Version

```bash
cd /opt/spotiflac
docker-compose pull
docker-compose up -d
```

That's it! Docker automatically restarts with the new image.

## Manual Build & Push (Alternative)

If you prefer to build and push manually:

```bash
# 1. Build for multiple platforms
docker buildx create --use
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --tag ghcr.io/raindancer118/spotiflac-server:latest \
  --tag ghcr.io/raindancer118/spotiflac-server:v1.0.0 \
  --push \
  .

# 2. Or build for single platform
docker build -t ghcr.io/raindancer118/spotiflac-server:latest .
docker push ghcr.io/raindancer118/spotiflac-server:latest
```

## Image Size Verification

After building, verify the optimized size:

```bash
docker images | grep spotiflac
```

Expected size: **~100-150 MB** (thanks to multi-stage build)

## Multi-Platform Support

The GitHub Actions workflow builds for:
- **linux/amd64** - Standard x86_64 servers
- **linux/arm64** - ARM servers (e.g., AWS Graviton, Raspberry Pi 4)

Pull the appropriate architecture automatically:
```bash
docker pull ghcr.io/raindancer118/spotiflac-server:latest
# Docker automatically selects the right platform
```

## Troubleshooting

### Image pull failed (private repo)

**Error:** `denied: permission_denied`

**Solution:** Login to GHCR first:
```bash
echo $GITHUB_TOKEN | docker login ghcr.io -u YOUR_USERNAME --password-stdin
```

### Image pull is slow

**Solution:** Use a closer registry mirror or build locally once and cache.

### Wrong architecture pulled

**Solution:** Explicitly specify platform:
```bash
docker pull --platform linux/amd64 ghcr.io/raindancer118/spotiflac-server:latest
```

## Best Practices

1. **Use version tags** in production instead of `latest`
2. **Set up automatic updates** with Watchtower (optional):
   ```bash
   docker run -d \
     --name watchtower \
     -v /var/run/docker.sock:/var/run/docker.sock \
     containrrr/watchtower \
     spotiflac-server
   ```
3. **Monitor image updates** via GitHub Releases
4. **Test updates** in development before production
