---
sidebar_position: 4
slug: /deployment/self-hosted
---

# Self-Hosted Deployment

Deploy Anexis Server for personal or organizational self-hosting.

## Why Self-Host?

- **Full Control**: Your data, your infrastructure
- **No Cloud Costs**: Use your own storage
- **Privacy**: No third-party data handling
- **Customization**: Configure for your needs

---

## Quick Start

### 1. Prerequisites

- Docker & Docker Compose
- 2GB RAM minimum
- 10GB storage (for files)

### 2. Run

```bash
# Clone repository
git clone https://github.com/Treefle-labs/anexis-server.git
cd anexis-server

# Start self-hosted stack
docker-compose -f docker-compose.selfhosted.yml up -d
```

### 3. Access

- API: `http://localhost:8080`
- Health: `http://localhost:8080/health`

---

## Configuration

### Storage Options

#### Local Storage (Default)

```yaml
STORAGE_PROVIDER=local
STORAGE_LOCAL_PATH=/app/data/storage
```

Data stored in Docker volume.

#### S3-Compatible Storage

```yaml
STORAGE_PROVIDER=s3
S3_ENDPOINT=https://s3.amazonaws.com
S3_REGION=us-east-1
S3_BUCKET=my-backup-bucket
S3_ACCESS_KEY=your_access_key
S3_SECRET_KEY=your_secret_key
```

Works with:
- AWS S3
- MinIO
- DigitalOcean Spaces
- Wasabi
- Backblaze B2 (via S3 compatibility)

### Authentication

In standalone mode, use API key:

```bash
# Set your API key
ANEXIS_API_KEY=your-secret-key
```

All requests use `X-API-Key` header:

```bash
curl -X POST http://localhost:8080/api/v1/files/upload \
  -H "X-API-Key: your-secret-key" \
  -F "file=@document.pdf"
```

---

## Use Cases

### Personal Cloud Storage

Replace Google Drive, Dropbox, etc.:

```bash
# Your own cloud
docker-compose -f docker-compose.selfhosted.yml up -d
```

### Backup Server

Use S3 backend for backups:

```yaml
STORAGE_PROVIDER=s3
S3_BUCKET=my-backups
```

### Media Server

Store photos, videos locally:

```yaml
STORAGE_PROVIDER=local
STORAGE_LOCAL_PATH=/mnt/media/anexis
```

Mount external drive:

```yaml
volumes:
  - /media/my硬盘:/app/data/storage
```

---

## Data Management

### Backup Database

```bash
docker-compose -f docker-compose.selfhosted.yml exec postgres pg_dump -U anexis anexis > backup.sql
```

### Restore Database

```bash
docker-compose -f docker-compose.selfhosted.yml exec -T postgres psql -U anexis anexis < backup.sql
```

### Backup Files

```bash
docker volume ls | grep anexis
docker run --rm -v anexis_anexis_data:/data -v $(pwd):/backup alpine tar czf /backup/anexis-files.tar.gz /data
```

---

## Update

```bash
# Pull latest
git pull

# Rebuild
docker-compose -f docker-compose.selfhosted.yml build

# Restart
docker-compose -f docker-compose.selfhosted.yml up -d
```

---

## Troubleshooting

### Cannot Connect

Check logs:
```bash
docker-compose -f docker-compose.selfhosted.yml logs
```

### Storage Full

- Check disk space: `df -h`
- Clean temp files: `docker system prune`

### Database Connection Failed

Ensure PostgreSQL is healthy:
```bash
docker-compose -f docker-compose.selfhosted.yml ps
```