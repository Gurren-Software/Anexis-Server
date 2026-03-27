---
sidebar_position: 2
slug: /deployment/docker-compose
---

# Docker Compose Deployment

Use Docker Compose for easier deployment.

## Self-Hosted Quick Start

The fastest way to start self-hosted:

```bash
docker-compose -f docker-compose.selfhosted.yml up -d
```

This starts:
- Anexis API server (standalone mode, local storage)
- PostgreSQL database

Access at `http://localhost:8080`

---

## Development Mode

```bash
docker-compose up -d
```

Access at `http://localhost:8080`

### Configuration

Create a `.env` file:

```bash
# Server
SERVER_MODE=standalone
SERVER_PORT=8080

# Storage
STORAGE_PROVIDER=local
STORAGE_LOCAL_PATH=/app/data/storage

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=anexis
DB_PASSWORD=anexis_password
DB_NAME=anexis

# Auth
ANEXIS_API_KEY=your-secret-api-key
JWT_SECRET=change-this-in-production
```

---

## Production with B2

Create `docker-compose.prod.yml`:

```bash
docker-compose -f docker-compose.prod.yml up -d
```

This starts:
- 3 API replicas (load balanced)
- PostgreSQL
- Nginx load balancer

---

## Custom Compose File

Create your own `docker-compose.yml`:

```yaml
version: '3.8'

services:
  anexis:
    image: anexis-server:latest
    ports:
      - "8080:8080"
    environment:
      - SERVER_MODE=standalone
      - STORAGE_PROVIDER=local
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=your_password
      - DB_NAME=anexis
      - ANEXIS_API_KEY=your-api-key
    volumes:
      - anexis-data:/app/data
    depends_on:
      postgres:
        condition: service_healthy

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: your_password
      POSTGRES_DB: anexis
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  anexis-data:
  postgres-data:
```

---

## Useful Commands

```bash
# Start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down

# Rebuild
docker-compose build
docker-compose up -d

# Scale (production only)
docker-compose up -d --scale anexis=3
```