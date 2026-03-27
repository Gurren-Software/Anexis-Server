---
sidebar_position: 3
---

# Installation

Choose your installation method based on your needs.

## Installation Methods

| Method | Best For | Complexity |
|--------|----------|------------|
| Docker Compose | Quick start, self-hosted | Easy |
| Docker | Custom deployments | Medium |
| Binary (Release) | Production servers | Easy |
| Source | Development | Hard |

---

## Docker Compose (Recommended)

### Quick Start

```bash
# Clone repository
git clone https://github.com/Treefle-labs/anexis-server.git
cd anexis-server

# Start self-hosted
docker-compose -f docker-compose.selfhosted.yml up -d
```

### Development Mode

```bash
docker-compose up -d
```

---

## Docker

### Pull Image

```bash
docker pull ghcr.io/treefle-labs/anexis-server:latest
```

### Run Container

```bash
docker run -d \
  --name anexis-server \
  -p 8080:8080 \
  -e SERVER_MODE=standalone \
  -e STORAGE_PROVIDER=local \
  -e DB_HOST=postgres \
  -e ANEXIS_API_KEY=your-key \
  ghcr.io/treefle-labs/anexis-server:latest
```

---

## Binary (Pre-built)

Download the pre-built binary from GitHub Releases:

```bash
# Download latest release (amd64)
curl -L https://github.com/Treefle-labs/anexis-server/releases/latest/download/anexis-server-linux-amd64 -o anexis-server
chmod +x anexis-server

# Run
./anexis-server
```

Available architectures:
- `anexis-server-linux-amd64` - x86_64
- `anexis-server-linux-arm64` - ARM64 (Raspberry Pi, etc.)

---

## Build from Source

### Prerequisites

- Go 1.24+
- PostgreSQL 16

### Build

```bash
# Clone
git clone https://github.com/Treefle-labs/anexis-server.git
cd anexis-server

# Build
cd apps/api
go build -o anexis-server ./cmd/server

# Run
./anexis-server
```

---

## Next Steps

- [Configuration](./configuration) - Customize your setup
- [Deployment](./deployment/docker) - Production deployment
- [API Reference](./api/auth) - Start using the API