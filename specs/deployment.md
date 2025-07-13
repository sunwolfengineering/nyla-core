# Nyla Analytics - Deployment Specification (Core)

## Overview

Nyla Analytics Core is designed for simple self-hosting with minimal operational overhead. The entire application is packaged as a single container with SQLite for data storage, optimized for single-site analytics.

**Deployment:**
Self-hosted single binary with SQLite database, packaged as a Docker container for simple deployment.

## Container Design

### Base Image

```dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o nyla

FROM alpine:3.19
COPY --from=builder /app/nyla /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/nyla"]
```

### Multi-stage Build

1. Build stage
   - Compile Go binary
   - Bundle frontend assets
   - Generate static files

2. Runtime stage
   - Minimal Alpine base
   - SQLite and CA certificates
   - Configuration files
   - Data volume mounts

## Configuration

### Environment Variables

```bash
# Server (Core)
NYLA_HOST=0.0.0.0
NYLA_PORT=3000
NYLA_ENV=production
NYLA_EDITION=core  # Set to 'core' for self-hosted

# Database (SQLite only in core)
NYLA_DB_PATH=/data/nyla.db
NYLA_DB_BACKUP_PATH=/backup

# Security
NYLA_API_KEY=nyla_key_xxx
NYLA_ENCRYPTION_KEY=xxx
NYLA_ALLOWED_ORIGINS=https://yourdomain.com

# Privacy (Core defaults)
NYLA_IP_ANONYMIZATION=true
NYLA_RETENTION_DAYS=90
NYLA_RESPECT_DNT=true
NYLA_SITE_NAME="My Site"

# Feature Flags (Core)
NYLA_ENABLE_MULTI_SITE=false  # Always false in core
NYLA_ENABLE_TEAMS=false       # Always false in core

# Logging
NYLA_LOG_LEVEL=info
NYLA_LOG_FORMAT=json
```

### Configuration File

```yaml
# config.yaml
server:
  host: 0.0.0.0
  port: 3000
  read_timeout: 5s
  write_timeout: 10s
  shutdown_timeout: 30s

database:
  path: /data/nyla.db
  backup:
    path: /backup
    interval: 24h
    retain: 7
  vacuum_interval: 168h

security:
  api_key: nyla_key_xxx
  encryption_key: xxx
  allowed_origins:
    - https://app.getnyla.app
    - https://dashboard.getnyla.app
  rate_limits:
    collect: 100/minute
    query: 60/minute

privacy:
  ip_anonymization: true
  retention_days: 90
  respect_dnt: true
  pii_patterns:
    - email
    - phone
    - credit_card

logging:
  level: info
  format: json
  output: stdout
```

## Directory Structure

```
/
├── usr/
│   └── local/
│       └── bin/
│           └── nyla
├── etc/
│   └── nyla/
│       ├── config.yaml
│       └── ssl/
├── data/
│   └── nyla.db
└── backup/
    └── nyla-YYYY-MM-DD.db
```

## Volume Management

### Data Volume

```yaml
volumes:
  - /path/to/data:/data
  - /path/to/backup:/backup
  - /path/to/config:/etc/nyla
```

### Backup Strategy

1. SQLite Online Backup API
2. Daily snapshots
3. Retention policy
4. Integrity verification
5. Optional encryption

## Resource Requirements

### Minimum Requirements

- CPU: 1 core
- RAM: 512MB
- Disk: 1GB

### Recommended

- CPU: 2 cores
- RAM: 1GB
- Disk: 10GB

## Health Monitoring

### Health Check Endpoint

```
GET /health
```

```json
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime": "24h",
  "database": "connected",
  "metrics": {
    "events_today": 1234,
    "active_sites": 2
  }
}
```

### Metrics

- Prometheus endpoint: `/metrics`
- Basic system metrics
- Application metrics
- Custom metrics

## Logging

### Log Format

```json
{
  "timestamp": "2024-03-14T15:09:26Z",
  "level": "info",
  "message": "Request processed",
  "method": "POST",
  "path": "/v1/collect",
  "duration_ms": 42,
  "status": 200
}
```

### Log Levels

- error: Errors requiring attention
- warn: Warning conditions
- info: Normal operations
- debug: Detailed debugging

## Security

### SSL/TLS

- Auto-HTTPS with Let's Encrypt for *.getnyla.app domains
- HTTP/2 support
- Modern cipher suites
- Perfect forward secrecy

### Headers

```nginx
add_header Strict-Transport-Security "max-age=63072000";
add_header X-Frame-Options "DENY";
add_header X-Content-Type-Options "nosniff";
add_header Content-Security-Policy "default-src 'self'";
```

## Docker Compose

```yaml
version: '3.8'

services:
  nyla:
    image: nyla/analytics:latest
    container_name: nyla
    restart: unless-stopped
    environment:
      - NYLA_ENV=production
      - NYLA_PORT=3000
    volumes:
      - ./data:/data
      - ./backup:/backup
      - ./config:/etc/nyla
    ports:
      - "3000:3000"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

## Deployment Steps

1. Initial Setup
```bash
# Create directories
mkdir -p data backup config
chmod 700 data backup

# Generate configuration
nyla init > config/config.yaml

# Generate secrets
nyla generate-keys > config/secrets
```

2. Database Setup
```bash
# Initialize database
nyla migrate up

# Verify schema
nyla db-check
```

3. Container Launch
```bash
# Pull image
docker pull nyla/analytics:latest

# Start service
docker-compose up -d

# Verify health
curl http://localhost:3000/health
```

## Backup Procedures

### Automated Backups

```bash
#!/bin/bash
# backup.sh
DATE=$(date +%Y-%m-%d)
sqlite3 /data/nyla.db ".backup '/backup/nyla-$DATE.db'"
find /backup -name "nyla-*.db" -mtime +7 -delete
```

### Manual Backup

```bash
# Create backup
nyla backup create

# List backups
nyla backup list

# Restore from backup
nyla backup restore <backup-file>
```

## Monitoring

### Key Metrics

1. System
   - CPU usage
   - Memory usage
   - Disk I/O
   - Network I/O

2. Application
   - Request rate
   - Error rate
   - Response times
   - Active sessions

3. Database
   - Size
   - Write rate
   - Read rate
   - Query times

### Alerting

1. Critical
   - Service down
   - Database errors
   - High error rate
   - Disk space low

2. Warning
   - High latency
   - Backup failures
   - Certificate expiry
   - Resource usage

## Maintenance

### Regular Tasks

1. Daily
   - Health check
   - Backup verification
   - Log rotation

2. Weekly
   - Database vacuum
   - Index optimization
   - Metric aggregation

3. Monthly
   - Security updates
   - SSL renewal check
   - Resource review

### Upgrades

1. Preparation
   - Backup verification
   - Version compatibility
   - Downtime window

2. Process
   - Stop service
   - Backup data
   - Update image
   - Run migrations
   - Start service
   - Verify health

3. Rollback
   - Stop service
   - Restore backup
   - Revert image
   - Start service
   - Verify health 