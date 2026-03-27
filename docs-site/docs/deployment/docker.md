---
sidebar_position: 1
slug: /deployment/docker
---

# Docker Deployment

Deploy Anexis Server using Docker.

## Prerequisites

- Docker 20.10+
- Docker Compose 2.0+ (for compose files)

## Using Pre-built Image

Pull the latest image from GitHub Container Registry:

```bash
docker pull ghcr.io/treefle-labs/anexis-server:latest
```

Or use a specific version:

```bash
docker pull ghcr.io/treefle-labs/anexis-server:v1.0.0
```

## Build Image (Optional)

If you want to build yourself:

```bash
# Clone and build
git clone https://github.com/Treefle-labs/anexis-server.git
cd anexis-server

# Build the image
docker build -t anexis-server:latest -f apps/api/Dockerfile .
```

## Run Container

### Basic Run (Local Storage)

```bash
docker run -d \
  --name anexis-server \
  -p 8080:8080 \
  -e SERVER_MODE=standalone \
  -e STORAGE_PROVIDER=local \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=password \
  -e DB_NAME=anexis \
  -e ANEXIS_API_KEY=your-api-key \
  -v anexis-data:/app/data \
  anexis-server:latest
```

### With B2 Storage

```bash
docker run -d \
  --name anexis-server \
  -p 8080:8080 \
  -e SERVER_MODE=saas \
  -e STORAGE_PROVIDER=b2 \
  -e B2_APPLICATION_KEY_ID=your_key_id \
  -e B2_APPLICATION_KEY=your_key \
  -e B2_BUCKET_NAME=your_bucket \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=password \
  -e DB_NAME=anexis \
  anexis-server:latest
```

### With S3 Storage

```bash
docker run -d \
  --name anexis-server \
  -p 8080:8080 \
  -e SERVER_MODE=standalone \
  -e STORAGE_PROVIDER=s3 \
  -e S3_ENDPOINT=https://s3.amazonaws.com \
  -e S3_REGION=us-east-1 \
  -e S3_BUCKET=your-bucket \
  -e S3_ACCESS_KEY=your_access_key \
  -e S3_SECRET_KEY=your_secret_key \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=password \
  -e DB_NAME=anexis \
  anexis-server:latest
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `SERVER_MODE` | `saas` or `standalone` | Yes |
| `STORAGE_PROVIDER` | `local`, `b2`, `s3` | Yes |
| `DB_HOST` | Database host | Yes |
| `DB_PORT` | Database port | Yes |
| `DB_USER` | Database user | Yes |
| `DB_PASSWORD` | Database password | Yes |
| `DB_NAME` | Database name | Yes |
| `ANEXIS_API_KEY` | API key (standalone mode) | No |
| `JWT_SECRET` | JWT secret (saas mode) | No |

## Volumes

| Volume | Description |
|--------|-------------|
| `/app/data/storage` | Local storage (when STORAGE_PROVIDER=local) |
| `/app/data/temp` | Temporary files |
| `/app/data/uploads` | Upload staging |

## Health Check

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "ok",
  "service": "anexis-api",
  "mode": "standalone",
  "storage": "local"
}
```