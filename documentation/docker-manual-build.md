# Docker Image Build & Push Guide

## Quick Commands

After logging in to GHCR:

```bash
# Build Docker image (Go 1.24-RC)
docker build -t ghcr.io/raindancer118/spotiflac-server:latest -t ghcr.io/raindancer118/spotiflac-server:main .

# Push to GHCR
docker push ghcr.io/raindancer118/spotiflac-server:latest
docker push ghcr.io/raindancer118/spotiflac-server:main
```

## GHCR Login

1. Create GitHub Personal Access Token:
   - Go to: https://github.com/settings/tokens
   - Click "Generate new token (classic)"
   - Select scopes: `write:packages`, `read:packages`
   - Generate and copy token

2. Login to GHCR:
   ```bash
   echo YOUR_TOKEN | docker login ghcr.io -u Raindancer118 --password-stdin
   ```

## Why Go 1.24-RC?

Dependencies require Go 1.24+:
- `gin-contrib/sse@v1.1.0`
- `golang.org/x/crypto@v0.45.0`

Go 1.24 is RC (Release Candidate) - stable enough for use.
