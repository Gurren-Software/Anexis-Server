---
sidebar_position: 2
---

# Getting Started

This guide will help you get Anexis Server up and running quickly.

## Prerequisites

- Go 1.24+ (for development)
- Docker & Docker Compose (for containerized deployment)
- PostgreSQL 16 (external or via Docker)

## Quick Start with Docker

The fastest way to get started is using Docker Compose:

```bash
# Clone the repository
git clone https://github.com/Treefle-labs/anexis-server.git
cd anexis-server

# Start with docker-compose (development)
docker-compose up -d
```

## Local Development

### 1. Clone and Setup

```bash
git clone https://github.com/Treefle-labs/anexis-server.git
cd anexis-server
```

### 2. Configure Environment

Copy the example environment file and customize:

```bash
cp .env.example .env
```

### 3. Start Database

```bash
# Using Docker
docker run -d \
  --name anexis-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=your_password \
  -e POSTGRES_DB=anexis \
  -p 5432:5432 \
  postgres:16-alpine
```

### 4. Run the Server

```bash
cd apps/api
go run cmd/server/main.go
```

The server will start at `http://localhost:8080`

## Self-Hosted Quick Start

For self-hosted deployment with local storage:

```bash
docker-compose -f docker-compose.selfhosted.yml up -d
```

This starts:
- Anexis API server (standalone mode)
- PostgreSQL database

## Next Steps

- [Configuration Guide](./configuration) - Customize your setup
- [Deployment](./deployment/docker) - Production deployment options
- [API Reference](./api/auth) - Start using the API