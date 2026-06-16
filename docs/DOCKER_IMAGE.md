# Anexis Server Docker Image

Self-hostable cloud file storage API with local, S3-compatible, and Backblaze B2 storage support.

Anexis Server can run as:

- **Standalone mode**: single-instance/self-hosted API protected by `X-API-Key`.
- **SaaS mode**: multi-user API with JWT authentication.

## Image

```text
ghcr.io/gurren-software/anexis-server
```

Use a released tag when deploying:

```bash
docker pull ghcr.io/gurren-software/anexis-server:<tag>
```

Example:

```bash
docker pull ghcr.io/gurren-software/anexis-server:v0.1.0
```

## Quick Start

This example runs Anexis Server in standalone mode with local file storage.

You need a PostgreSQL database. The example assumes PostgreSQL is reachable from the container as `postgres`.

```bash
docker run -d \
  --name anexis-server \
  -p 8080:8080 \
  -v anexis_data:/app/data \
  -e SERVER_MODE=standalone \
  -e SERVER_HOST=0.0.0.0 \
  -e SERVER_PORT=8080 \
  -e STORAGE_PROVIDER=local \
  -e STORAGE_LOCAL_PATH=/app/data/storage \
  -e ANEXIS_API_KEY=change-this-api-key \
  -e JWT_SECRET=change-this-jwt-secret \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=anexis \
  -e DB_PASSWORD=anexis_password \
  -e DB_NAME=anexis \
  -e DB_SSLMODE=disable \
  ghcr.io/gurren-software/anexis-server:<tag>
```

The API will be available at:

```text
http://localhost:8080
```

Health check:

```bash
curl http://localhost:8080/health
```

## Docker Compose

Recommended standalone deployment with PostgreSQL and persistent local storage:

```yaml
services:
  anexis:
    image: ghcr.io/gurren-software/anexis-server:<tag>
    container_name: anexis-server
    ports:
      - "8080:8080"
    volumes:
      - anexis_data:/app/data
    environment:
      SERVER_MODE: standalone
      SERVER_HOST: 0.0.0.0
      SERVER_PORT: 8080
      STORAGE_PROVIDER: local
      STORAGE_LOCAL_PATH: /app/data/storage
      ANEXIS_API_KEY: change-this-api-key
      JWT_SECRET: change-this-jwt-secret
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: anexis
      DB_PASSWORD: anexis_password
      DB_NAME: anexis
      DB_SSLMODE: disable
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    container_name: anexis-postgres
    environment:
      POSTGRES_USER: anexis
      POSTGRES_PASSWORD: anexis_password
      POSTGRES_DB: anexis
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U anexis"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  anexis_data:
  postgres_data:
```

Start it:

```bash
docker compose up -d
```

View logs:

```bash
docker logs -f anexis-server
```

## Authentication

### Standalone Mode

Standalone mode is the simplest self-hosted setup.

```text
SERVER_MODE=standalone
ANEXIS_API_KEY=<your-secret-api-key>
```

Protected requests must include:

```text
X-API-Key: <your-secret-api-key>
```

### SaaS Mode

SaaS mode enables user registration/login and JWT-authenticated routes.

```text
SERVER_MODE=saas
JWT_SECRET=<long-random-secret>
```

Protected requests must include:

```text
Authorization: Bearer <access_token>
```

## Storage Providers

### Local Storage

Best for simple self-hosted deployments.

```text
STORAGE_PROVIDER=local
STORAGE_LOCAL_PATH=/app/data/storage
```

Mount `/app/data` as a persistent volume:

```text
-v anexis_data:/app/data
```

Without a volume, uploaded files are lost when the container is removed.

### S3-Compatible Storage

Use this for AWS S3, MinIO, Cloudflare R2, or another S3-compatible provider.

```text
STORAGE_PROVIDER=s3
S3_ENDPOINT=https://s3.example.com
S3_REGION=us-east-1
S3_BUCKET=anexis
S3_ACCESS_KEY=<access-key>
S3_SECRET_KEY=<secret-key>
S3_FORCE_PATH_STYLE=false
S3_BASE_PATH=
```

For MinIO, you usually want:

```text
S3_FORCE_PATH_STYLE=true
```

### Backblaze B2

```text
STORAGE_PROVIDER=b2
B2_APPLICATION_KEY_ID=<key-id>
B2_APPLICATION_KEY=<application-key>
B2_BUCKET_NAME=<bucket-name>
```

## Environment Variables

### Server

| Variable | Default | Description |
| --- | --- | --- |
| `SERVER_HOST` | `0.0.0.0` | Address the API binds to inside the container. |
| `SERVER_PORT` | `8080` | API port inside the container. |
| `ENVIRONMENT` | `development` | Use `production` for deployments. |
| `DEBUG` | `true` | Use `false` for production. |
| `SERVER_MODE` | `saas` | `standalone` or `saas`. |

### Database

| Variable | Default | Description |
| --- | --- | --- |
| `DB_HOST` | `localhost` | PostgreSQL hostname. |
| `DB_PORT` | `5432` | PostgreSQL port. |
| `DB_USER` | `postgres` | PostgreSQL user. |
| `DB_PASSWORD` | empty | PostgreSQL password. |
| `DB_NAME` | `anexis` | PostgreSQL database. |
| `DB_SSLMODE` | `disable` | PostgreSQL SSL mode. |

### Security

| Variable | Required | Description |
| --- | --- | --- |
| `JWT_SECRET` | Yes in SaaS | Secret used to sign JWT tokens. |
| `ANEXIS_API_KEY` | Yes in standalone | API key for standalone deployments. |

Use long random values for secrets:

```bash
openssl rand -hex 32
```

### Limits

| Variable | Default | Description |
| --- | --- | --- |
| `RATE_LIMIT_REQUESTS` | `100` | Requests allowed per window. |
| `RATE_LIMIT_WINDOW_SECONDS` | `60` | Rate-limit window in seconds. |
| `MAX_UPLOAD_SIZE_MB` | `100` | Maximum upload size in MiB. |

## Volumes

For local storage deployments, persist:

```text
/app/data
```

The image creates:

```text
/app/data/storage
/app/data/temp
/app/data/uploads
```

## API Documentation

Once running, Swagger is available at:

```text
http://localhost:8080/swagger/index.html
```

Health endpoint:

```text
http://localhost:8080/health
```

## Common Operations

Stop the container:

```bash
docker stop anexis-server
```

Remove the container:

```bash
docker rm anexis-server
```

Follow logs:

```bash
docker logs -f anexis-server
```

Open a shell:

```bash
docker exec -it anexis-server sh
```

## Troubleshooting

### The API cannot connect to PostgreSQL

Check:

- `DB_HOST` matches the PostgreSQL service/container name.
- PostgreSQL is healthy.
- `DB_USER`, `DB_PASSWORD`, and `DB_NAME` match the database container.
- Both containers are on the same Docker network.

### Uploads disappear after recreating the container

You are probably using local storage without a persistent volume.

Mount `/app/data`:

```text
-v anexis_data:/app/data
```

### Protected routes return `401`

For standalone mode, send:

```text
X-API-Key: <ANEXIS_API_KEY>
```

For SaaS mode, send:

```text
Authorization: Bearer <access_token>
```

### The container is unhealthy

Check the health endpoint from inside the host:

```bash
curl http://localhost:8080/health
```

Then inspect logs:

```bash
docker logs anexis-server
```
