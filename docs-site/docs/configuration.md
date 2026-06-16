---
sidebar_position: 4
---

# Configuration

Anexis Server is configured via environment variables. This guide covers all available options.

## Server Mode

| Variable | Values | Default | Description |
|----------|--------|---------|-------------|
| `SERVER_MODE` | `saas`, `standalone` | `saas` | Deployment mode |
| `SERVER_HOST` | string | `0.0.0.0` | Server bind address |
| `SERVER_PORT` | string | `8080` | Server port |
| `ENVIRONMENT` | `development`, `production` | `development` | Runtime environment |
| `DEBUG` | `true`, `false` | `true` | Enable debug mode |

### SaaS Mode
Multi-user mode with JWT authentication and storage quotas.

### Standalone Mode
Self-hosted mode with API key authentication, no quotas.

## Authentication

### JWT (SaaS Mode)

| Variable | Default | Description |
|----------|---------|-------------|
| `JWT_SECRET` | - | Secret key for signing JWT tokens |
| `JWT_EXPIRATION_HOURS` | `24` | Token expiration time |

### API Key (Standalone Mode)

| Variable | Description |
|----------|-------------|
| `ANEXIS_API_KEY` | API key for standalone mode authentication |

## Storage Provider

| Variable | Values | Default | Description |
|----------|--------|---------|-------------|
| `STORAGE_PROVIDER` | `local`, `b2`, `s3` | `local` | Storage backend |

### Local Storage

| Variable | Default | Description |
|----------|---------|-------------|
| `STORAGE_LOCAL_PATH` | `./data/storage` | Path for local file storage |

### Backblaze B2

| Variable | Description |
|----------|-------------|
| `B2_APPLICATION_KEY_ID` | B2 key ID |
| `B2_APPLICATION_KEY` | B2 application key |
| `B2_BUCKET_NAME` | B2 bucket name |
| `B2_BUCKET_ID` | B2 bucket ID (optional) |

### S3-Compatible

| Variable | Default | Description |
|----------|---------|-------------|
| `S3_ENDPOINT` | - | S3 endpoint URL |
| `S3_REGION` | `us-east-1` | AWS region |
| `S3_BUCKET` | - | S3 bucket name |
| `S3_ACCESS_KEY` | - | AWS access key |
| `S3_SECRET_KEY` | - | AWS secret key |
| `S3_FORCE_PATH_STYLE` | `false` | Use path-style addressing |
| `S3_BASE_PATH` | - | Base path prefix |

## Database

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | - | Database password |
| `DB_NAME` | `anexis` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |

## File Upload

| Variable | Default | Description |
|----------|---------|-------------|
| `MAX_UPLOAD_SIZE_MB` | `100` | Maximum upload size in MB |

## Rate Limiting

| Variable | Default | Description |
|----------|---------|-------------|
| `RATE_LIMIT_REQUESTS` | `100` | Maximum requests per window |
| `RATE_LIMIT_WINDOW_SECONDS` | `60` | Rate limit window in seconds |

## Example Configuration

### SaaS Mode with B2

```bash
SERVER_MODE=saas
STORAGE_PROVIDER=b2
B2_APPLICATION_KEY_ID=your_key_id
B2_APPLICATION_KEY=your_key
B2_BUCKET_NAME=your_bucket
JWT_SECRET=your_jwt_secret
```

### Standalone Mode with Local Storage

```bash
SERVER_MODE=standalone
STORAGE_PROVIDER=local
STORAGE_LOCAL_PATH=/data/anexis
ANEXIS_API_KEY=your_api_key
```

### Standalone Mode with S3

```bash
SERVER_MODE=standalone
STORAGE_PROVIDER=s3
S3_ENDPOINT=https://s3.amazonaws.com
S3_REGION=us-east-1
S3_BUCKET=your-bucket
S3_ACCESS_KEY=your_access_key
S3_SECRET_KEY=your_secret_key
ANEXIS_API_KEY=your_api_key
```