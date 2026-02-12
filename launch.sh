#!/bin/bash
# SpotiFLAC Server Deployment Script
# Following rule #5: Automated deployment with git pull, build, and service restart
# This script should be run on the Ubuntu server to deploy updates

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
APP_DIR="/opt/spotiflac"
USER="spotiflac"
SERVICE_NAME="spotiflac-server"

echo "=== SpotiFLAC Server Deployment ==="

# Check if running as root
if [ "$EUID" -eq 0 ]; then
    echo -e "${YELLOW}Warning: Running as root. Consider running as the spotiflac user.${NC}"
fi

# 1. Git Operations
echo -e "${GREEN}[1/7] Pulling latest changes...${NC}"
git fetch origin
git pull origin main

# Check if pull was successful
if [ $? -ne 0 ]; then
    echo -e "${RED}Git pull failed. Aborting deployment.${NC}"
    exit 1
fi

# 2. Generate .env from config.yml (if needed)
echo -e "${GREEN}[2/7] Checking configuration...${NC}"
if [ ! -f "config.yml" ]; then
    echo -e "${YELLOW}Warning: config.yml not found. Creating from defaults...${NC}"
    if [ -f "config.yml.example" ]; then
        cp config.yml.example config.yml
        echo -e "${GREEN}Created config.yml from example. Please edit before continuing.${NC}"
        exit 0
    else
        echo -e "${RED}No config.yml or example found!${NC}"
        exit 1
    fi
fi

# 3. Install Go dependencies
echo -e "${GREEN}[3/7] Installing Go dependencies...${NC}"
go mod download
go mod tidy

# 4. Build backend server
echo -e "${GREEN}[4/7] Building backend server...${NC}"
go build -o spotiflac-server cmd/server/main.go
if [ $? -ne 0 ]; then
    echo -e "${RED}Server build failed!${NC}"
    exit 1
fi

# Build CLI tool
go build -o spotiflac cmd/cli/main.go
if [ $? -ne 0 ]; then
    echo -e "${RED}CLI build failed!${NC}"
    exit 1
fi

echo -e "${GREEN}Build successful!${NC}"

# 5. Build frontend
echo -e "${GREEN}[5/7] Building frontend...${NC}"
cd frontend

# Install/update npm dependencies
if [ ! -d "node_modules" ] || [ "package.json" -nt "node_modules" ]; then
    echo "Installing npm dependencies..."
    npm install
fi

# Build for production
npm run build

if [ $? -ne 0 ]; then
    echo -e "${RED}Frontend build failed!${NC}"
    cd ..
    exit 1
fi

echo -e "${GREEN}Frontend build successful!${NC}"
cd ..

# 6. Data migration (if needed)
echo -e "${GREEN}[6/7] Checking for data migrations...${NC}"
# Add any database migration logic here if needed
# For now, SQLite database is backward compatible

# 7. Restart service
echo -e "${GREEN}[7/7] Restarting service...${NC}"

# Check if systemd service exists
if systemctl list-units --full -all | grep -Fq "$SERVICE_NAME.service"; then
    echo "Restarting systemd service..."
    systemctl restart $SERVICE_NAME

    # Check service status
    if systemctl is-active --quiet $SERVICE_NAME; then
        echo -e "${GREEN}Service restarted successfully!${NC}"
    else
        echo -e "${RED}Service failed to start. Check logs with: journalctl -u $SERVICE_NAME -n 50${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}Systemd service not found. Service management skipped.${NC}"
    echo -e "${YELLOW}To install the service:${NC}"
    echo "  1. Copy spotiflac-server.service to /etc/systemd/system/"
    echo "  2. Run: systemctl daemon-reload"
    echo "  3. Run: systemctl enable spotiflac-server"
    echo "  4. Run: systemctl start spotiflac-server"
fi

echo ""
echo -e "${GREEN}=== Deployment Complete ===${NC}"
echo ""
echo "Server binary: ./spotiflac-server"
echo "CLI binary: ./spotiflac"
echo "Frontend dist: ./frontend/dist"
echo ""
echo "To start manually: ./spotiflac-server"
echo "To check status: systemctl status $SERVICE_NAME"
echo "To view logs: journalctl -u $SERVICE_NAME -f"
