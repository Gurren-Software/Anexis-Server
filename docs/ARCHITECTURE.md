# Anexis Server Architecture

This document describes the vertical slice architecture used in Anexis Server.

## Overview

Anexis Server uses a **vertical slice architecture** organized as a **monorepo**. This approach groups code by feature rather than by technical layer, making it easier to understand, test, and scale individual features.

## Monorepo Structure

```
anexis-server/
├── packages/                    # Shared packages
│   └── database/                # Database package (shared across services)
│
├── apps/                        # Applications
│   └── api/                     # API server
│
├── nginx/                       # Load balancer configuration
├── scripts/                     # Helper scripts
└── docs/                        # Documentation
```

## Packages

### database

The `packages/database` package is a shared module containing:

- **Connection management** with connection pooling
- **Domain models** (User, File, Link, MigrationJob, BackupJob)
- **Migration tooling** via Atlas

This package is imported by all services that need database access, ensuring:
- Single source of truth for models
- Consistent database configuration
- Shared migrations

## API Server Architecture

The API server in `apps/api` follows vertical slice architecture:

```
apps/api/
├── cmd/
│   └── server/
│       └── main.go              # Entry point with DI
│
├── internal/
│   ├── config/                  # Environment configuration
│   │
│   ├── infrastructure/          # Cross-cutting concerns
│   │   ├── http/
│   │   │   ├── server.go        # HTTP server with graceful shutdown
│   │   │   ├── middleware/      # JWT, CORS, logging
│   │   │   └── response/        # Standardized responses
│   │   │
│   │   └── storage/             # Storage abstraction
│   │       ├── provider.go      # Interface
│   │       └── backblaze/       # B2 implementation
│   │
│   └── features/                # VERTICAL SLICES
│       ├── auth/
│       ├── files/
│       ├── links/
│       ├── migration/
│       └── backup/
│
└── docs/                        # Swagger documentation
```

## Feature Slices

Each feature slice is self-contained:

```
features/<feature>/
├── dto.go           # Request/Response types
├── repository.go    # Database operations
├── service.go       # Business logic
├── handler.go       # HTTP handlers
└── routes.go        # Route registration
```

### Benefits

1. **High Cohesion** - All code for a feature is in one place
2. **Low Coupling** - Features don't depend on each other
3. **Easy Testing** - Each slice can be tested in isolation
4. **Independent Scaling** - Features can be split into microservices if needed
5. **Clear Ownership** - Teams can own specific slices

## Data Flow

```
Request → Router → Handler → Service → Repository → Database
                      ↓
                   Response
```

1. **Router** routes to appropriate handler
2. **Handler** validates input, calls service, formats response
3. **Service** contains business logic, orchestrates operations
4. **Repository** handles database queries

## Infrastructure Layer

Cross-cutting concerns are in `internal/infrastructure/`:

### HTTP
- Server setup with graceful shutdown
- JWT/CORS/Logging middleware
- Standardized response helpers

### Storage
- Provider interface for abstraction
- Backblaze B2 implementation
- Can add S3, GCS, etc. implementations

## Dependency Injection

Dependencies are injected in `main.go`:

```go
// Repositories
authRepo := auth.NewRepository(db.DB)
filesRepo := files.NewRepository(db.DB)

// Services
authService := auth.NewService(authRepo, jwtSecret, expiration)
filesService := files.NewService(filesRepo, storage, authRepo)

// Handlers
authHandler := auth.NewHandler(authService)
filesHandler := files.NewHandler(filesService)
```

## Scaling

### Database Scaling
- `packages/database` is shared, so all API instances use the same schema
- Connection pooling is configured for multiple instances

### API Scaling
- Stateless design allows horizontal scaling
- Docker Compose production config runs 3+ replicas
- Nginx load balancer distributes traffic

### Storage Scaling
- Backblaze B2 handles file storage at scale
- Storage provider interface allows switching providers

## Adding a New Feature

1. Create folder: `apps/api/internal/features/<feature>/`
2. Add files: `dto.go`, `repository.go`, `service.go`, `handler.go`, `routes.go`
3. Add models to `packages/database/models/` if needed
4. Register routes in `main.go`
5. Generate migrations: `make migrate-new NAME=add_<feature>`

## Testing Strategy

```
Unit Tests     → Service logic, Repository queries
Integration    → API endpoints with test database
E2E Tests      → Full flow with Docker environment
```

Run tests:
```bash
make test           # All tests
make test-coverage  # With coverage report
```
