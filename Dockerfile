# SpotiFLAC Dockerfile
# Multi-stage build for optimal image size

# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy all frontend files first (postinstall script needs scripts/generate-icon.js)
COPY frontend/ .

# Install pnpm and dependencies
RUN npm install -g pnpm && pnpm install --frozen-lockfile

# Build frontend for production
RUN pnpm run build

# Stage 2: Build Go Backend
FROM golang:1.24-rc-alpine AS backend-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Set GOTOOLCHAIN to auto to allow using newer Go versions
ENV GOTOOLCHAIN=auto

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o spotiflac-server ./cmd/server

# Build the CLI binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o spotiflac-cli ./cmd/cli

# Stage 3: Runtime Image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    ffmpeg \
    tzdata

# Create non-root user
RUN addgroup -g 1000 spotiflac && \
    adduser -D -u 1000 -G spotiflac spotiflac

WORKDIR /app

# Copy binaries from builder stages
COPY --from=backend-builder /app/spotiflac-server /app/spotiflac-server
COPY --from=backend-builder /app/spotiflac-cli /app/spotiflac-cli
COPY --from=frontend-builder /app/frontend/dist /app/web

# Create necessary directories with correct permissions
RUN mkdir -p /app/downloads /app/data && \
    chown -R spotiflac:spotiflac /app

# Switch to non-root user
USER spotiflac

# Expose HTTP port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Set environment variables
ENV SERVER_MODE=server \
    SERVER_PORT=8080 \
    WEB_DIR=/app/web

# Run the server
CMD ["/app/spotiflac-server"]
