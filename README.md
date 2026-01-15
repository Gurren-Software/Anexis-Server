# Anexis Server

A cloud file storage server built with Go, featuring vertical slice architecture, Backblaze B2 integration, and scalable deployment with load balancing.

## Features

- 🔐 **JWT Authentication** - Secure user registration, login, and token refresh
- 📁 **File Management** - Upload, download, folders, rename, move, delete with compression
- 🔗 **Access Links** - Permanent, temporal, streaming, and download links with password protection
- 🔄 **Provider Migration** - Import files from Google Drive, Amazon S3, OneDrive, Dropbox
- 💾 **Backup/Restore** - Export all user data as downloadable ZIP archives
- 🐳 **Docker Ready** - Development and production configurations with Nginx load balancer
- 📖 **Swagger Docs** - Auto-generated API documentation

## Architecture

```
anexis-server/
├── packages/
│   └── database/           # Shared database package
│       ├── config.go       # Connection pooling
│       └── models/         # User, File, Link, MigrationJob, BackupJob
│
├── apps/
│   └── api/                # API server (scalable)
│       ├── cmd/server/     # Entry point
│       └── internal/
│           ├── config/     # Environment config
│           ├── infrastructure/
│           │   ├── http/   # Server, middleware, responses
│           │   └── storage/# Provider interface + Backblaze
│           └── features/   # Vertical slices
│               ├── auth/
│               ├── files/
│               ├── links/
│               ├── migration/
│               └── backup/
│
├── nginx/                  # Load balancer config
├── docker-compose.yml      # Development
└── docker-compose.prod.yml # Production (3 replicas + LB)
```

## Quick Start

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- PostgreSQL 16+
- Backblaze B2 account

### Setup

```bash
# Clone the repository
git clone https://github.com/Treefle-labs/anexis-server.git
cd anexis-server

# Copy environment template
cp .env.example .env
# Edit .env with your credentials

# Start development environment
make dev

# Or start production environment
make prod
```

### Development

```bash
# Build the API server
make build

# Run tests
make test

# Generate Swagger documentation
make swagger

# Run database migrations
make migrate

# View all available commands
make help
```

## API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login and get tokens |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| GET | `/api/v1/auth/me` | Get current user profile |
| PUT | `/api/v1/auth/password` | Change password |

### Files
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/files` | List files |
| POST | `/api/v1/files/upload` | Upload a file |
| POST | `/api/v1/files/folder` | Create folder |
| GET | `/api/v1/files/:id` | Get file details |
| GET | `/api/v1/files/:id/download` | Download file |
| PUT | `/api/v1/files/:id/rename` | Rename file/folder |
| PUT | `/api/v1/files/:id/move` | Move file/folder |
| DELETE | `/api/v1/files/:id` | Delete file/folder |

### Links
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/links` | List access links |
| POST | `/api/v1/links` | Create access link |
| PUT | `/api/v1/links/:id` | Update link |
| DELETE | `/api/v1/links/:id` | Delete link |
| GET | `/api/v1/links/:token/access` | Access file via link |
| GET | `/api/v1/links/:token/stream` | Get streaming URL |

### Migration
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/migration` | List migration jobs |
| POST | `/api/v1/migration` | Start migration |
| GET | `/api/v1/migration/:id` | Get job status |
| POST | `/api/v1/migration/:id/cancel` | Cancel job |

### Backup
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/backup` | List backup jobs |
| POST | `/api/v1/backup/export` | Start export |
| GET | `/api/v1/backup/:id` | Get job status |
| GET | `/api/v1/backup/:id/download` | Get download URL |

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | API server port | `8080` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | - |
| `DB_NAME` | Database name | `anexis` |
| `JWT_SECRET` | JWT signing key | - |
| `JWT_EXPIRATION_HOURS` | Token expiration | `24` |
| `B2_APPLICATION_KEY_ID` | Backblaze key ID | - |
| `B2_APPLICATION_KEY` | Backblaze application key | - |
| `B2_BUCKET_NAME` | Backblaze bucket name | - |

## Deployment

### Production with Load Balancing

```bash
# Start with 3 API replicas behind Nginx
docker-compose -f docker-compose.prod.yml up -d

# Scale to more replicas
docker-compose -f docker-compose.prod.yml up -d --scale api=5
```

### Health Check

```bash
curl http://localhost:8080/health
# {"status":"ok","service":"anexis-api"}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📧 Email: support@treefle-labs.com
- 🐛 Issues: [GitHub Issues](https://github.com/Treefle-labs/anexis-server/issues)
