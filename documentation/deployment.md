# SpotiFLAC Server Deployment Guide

This guide covers deploying SpotiFLAC Server on Ubuntu Server 22.04 or later.

## Prerequisites

- Ubuntu Server 22.04+ (headless)
- Go 1.21+ installed
- Node.js 18+ and npm installed
- Nginx installed
- Git installed
- Minimum 2GB RAM recommended
- Sufficient disk space for downloads

## Installation Steps

### 1. Create Dedicated User

Following the principle of least privilege (Rule #10), create a dedicated user:

```bash
sudo useradd -r -s /bin/bash -d /opt/spotiflac spotiflac
sudo mkdir -p /opt/spotiflac
sudo chown spotiflac:spotiflac /opt/spotiflac
```

### 2. Clone Repository

```bash
sudo -u spotiflac git clone https://github.com/afkarxyz/SpotiFLAC.git /opt/spotiflac
```

### 3. Configure Application

```bash
cd /opt/spotiflac
sudo -u spotiflac cp config.yml.example config.yml
sudo -u spotiflac nano config.yml
```

Edit `config.yml` to set:
- Download path
- Server port (default: 8080)
- CORS origins
- Service preferences

### 4. Initial Build

```bash
cd /opt/spotiflac
sudo -u spotiflac ./launch.sh
```

This will:
- Install Go dependencies
- Build backend server
- Build CLI tool
- Install npm dependencies
- Build frontend for production

### 5. Configure Nginx

Copy the Nginx configuration:

```bash
sudo cp nginx.conf.example /etc/nginx/sites-available/spotiflac
```

Edit the configuration:

```bash
sudo nano /etc/nginx/sites-available/spotiflac
```

Change `server_name` to your domain or IP address.

Enable the site:

```bash
sudo ln -s /etc/nginx/sites-available/spotiflac /etc/nginx/sites-enabled/
sudo nginx -t  # Test configuration
sudo systemctl reload nginx
```

### 6. Configure Systemd Service

Install the systemd service:

```bash
sudo cp spotiflac-server.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable spotiflac-server
sudo systemctl start spotiflac-server
```

Check status:

```bash
sudo systemctl status spotiflac-server
```

### 7. Configure Firewall

Allow HTTP and HTTPS traffic:

```bash
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

### 8. SSL/TLS Configuration (Recommended)

Install certbot:

```bash
sudo apt install certbot python3-certbot-nginx
```

Obtain SSL certificate:

```bash
sudo certbot --nginx -d your-domain.com
```

Certbot will automatically configure Nginx for HTTPS.

## Accessing the Application

### Web Frontend

Open your browser and navigate to:
- HTTP: `http://your-server-ip/` or `http://your-domain.com/`
- HTTPS: `https://your-domain.com/` (after SSL setup)

### CLI Tool

SSH into the server and use the CLI:

```bash
cd /opt/spotiflac
./spotiflac download track https://open.spotify.com/track/abc123
```

## Updating the Application

To update to the latest version:

```bash
cd /opt/spotiflac
sudo -u spotiflac ./launch.sh
```

The `launch.sh` script will:
1. Pull latest code from Git
2. Rebuild backend and frontend
3. Restart the service

## Monitoring

### View Logs

```bash
# Service logs
sudo journalctl -u spotiflac-server -f

# Nginx logs
sudo tail -f /var/log/nginx/spotiflac_access.log
sudo tail -f /var/log/nginx/spotiflac_error.log
```

### Check Service Status

```bash
sudo systemctl status spotiflac-server
```

### Check Disk Space

Monitor download directory for available space:

```bash
df -h /path/to/downloads
```

## Troubleshooting

### Service Won't Start

1. Check logs:
   ```bash
   sudo journalctl -u spotiflac-server -n 50
   ```

2. Verify configuration:
   ```bash
   cd /opt/spotiflac
   cat config.yml
   ```

3. Check port availability:
   ```bash
   sudo netstat -tulpn | grep 8080
   ```

### Frontend Won't Load

1. Check Nginx configuration:
   ```bash
   sudo nginx -t
   ```

2. Verify frontend was built:
   ```bash
   ls -la /opt/spotiflac/frontend/dist/
   ```

3. Check Nginx logs:
   ```bash
   sudo tail -f /var/log/nginx/spotiflac_error.log
   ```

### Downloads Fail

1. Check FFmpeg installation:
   ```bash
   ffmpeg -version
   ```

2. Verify download path permissions:
   ```bash
   ls -ld /path/to/downloads
   ```

3. Check available disk space:
   ```bash
   df -h /path/to/downloads
   ```

### Database Issues

The application uses SQLite for history storage. Database file location is configured in `config.yml`.

If database is corrupted:

```bash
cd /opt/spotiflac
rm SpotiFLAC.db  # Backup first if needed
sudo systemctl restart spotiflac-server
```

## Security Recommendations

Following the security rules:

1. **Run as non-root user** (Rule #10: Least Privilege)
   - Service runs as `spotiflac` user
   - Systemd configuration enforces this

2. **Use HTTPS** (Rule #12: Encrypt data in transit)
   - Configure SSL/TLS with certbot
   - Redirect HTTP to HTTPS

3. **Keep software updated** (Rule #17: Patching is mandatory)
   -Run updates regularly:
   ```bash
   sudo apt update && sudo apt upgrade
   cd /opt/spotiflac && sudo -u spotiflac ./launch.sh
   ```

4. **Monitor logs** (Rule #16: Audit logs)
   - Regularly review system and application logs
   - Set up log rotation

5. **Firewall configuration**
   - Only expose necessary ports (80, 443)
   - Consider IP whitelisting if applicable

6. **Secure configuration**
   - Set proper file permissions:
     ```bash
     sudo chown -R spotiflac:spotiflac /opt/spotiflac
     sudo chmod 640 /opt/spotiflac/config.yml
     ```

## Performance Tuning

### Adjust Resource Limits

Edit systemd service to adjust memory limits:

```bash
sudo systemctl edit spotiflac-server
```

Add:
```ini
[Service]
MemoryLimit=4G
```

### Nginx Caching

The provided Nginx configuration includes caching for static assets. Adjust cache duration in `nginx.conf` if needed.

## Backup

### Configuration Backup

```bash
sudo cp /opt/spotiflac/config.yml /backup/location/
```

### Database Backup

```bash
sudo cp /opt/spotiflac/SpotiFLAC.db /backup/location/
```

### Automated Backups

Add to crontab:

```bash
0 2 * * * cp /opt/spotiflac/config.yml /backup/location/config.yml.$(date +\%Y\%m\%d)
0 2 * * * cp /opt/spotiflac/SpotiFLAC.db /backup/location/SpotiFLAC.db.$(date +\%Y\%m\%d)
```

## Support

For issues and questions, refer to:
- GitHub Issues: https://github.com/afkarxyz/SpotiFLAC/issues
- Documentation: `/opt/spotiflac/documentation/`
