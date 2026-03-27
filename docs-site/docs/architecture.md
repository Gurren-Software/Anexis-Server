---
sidebar_position: 5
---

# Architecture

Anexis Server follows a clean architecture pattern with vertical slices, making it easy to understand, test, and extend.

## Project Structure

```
anexis-server/
├── apps/
│   └── api/                        # Main API server
│       ├── cmd/server/              # Entry point
│       └── internal/
│           ├── config/              # Configuration loading
│           ├── infrastructure/      # Cross-cutting concerns
│           │   ├── http/            # HTTP server, middleware
│           │   └── storage/        # Storage providers
│           └── features/            # Business logic modules
│               ├── auth/            # Authentication
│               ├── files/           # File management
│               ├── links/           # Access links
│               ├── migration/       # Cloud migration
│               └── backup/          # Backup/restore
│
├── packages/
│   └── database/                    # Shared database package
│       ├── config.go                # Connection pooling
│       ├── models/                  # GORM models
│       └── migrations/              # SQL migrations
│
├── docs-site/                      # Documentation (Docusaurus)
├── docker-compose.yml              # Development config
├── docker-compose.prod.yml         # Production config
└── docker-compose.selfhosted.yml  # Self-hosted config
```

## Core Components

### Storage Provider Abstraction

The storage layer is abstracted through the `Provider` interface, allowing easy addition of new storage backends:

```
storage.Provider
├── Upload(key, reader, size, contentType)
├── Download(key) → io.ReadCloser
├── Delete(key)
├── GetURL(key, expiresIn) → string
├── GetStreamURL(key, expiresIn) → string
├── Exists(key) → bool
├── GetMetadata(key) → FileMetadata
├── List(prefix, maxKeys) → []FileMetadata
└── Copy(srcKey, dstKey)
```

**Available Providers:**
- `local` - Filesystem storage
- `b2` - Backblaze B2
- `s3` - S3-compatible (AWS, MinIO, etc.)

### Server Modes

| Mode | Auth | Quotas | Use Case |
|------|------|--------|-----------|
| **SaaS** | JWT | Yes | Multi-tenant cloud service |
| **Standalone** | API Key | No | Self-hosted, single organization |

### Request Flow

```
Client Request
    ↓
Middleware (CORS, Rate Limit, Auth)
    ↓
Feature Handler (auth, files, links, etc.)
    ↓
Service (Business Logic)
    ↓
Repository (Database)
    ↓
Storage Provider (Cloud/Local)
```

## Technology Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.24 |
| HTTP Router | Gin |
| Database | PostgreSQL 16 |
| ORM | GORM |
| Storage | Backblaze B2, S3-compatible |
| Auth | JWT / API Key |
| Docs | Docusaurus |