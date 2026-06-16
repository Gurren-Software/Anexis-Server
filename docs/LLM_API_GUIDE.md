# LLM API Implementation Guide

This file is written for coding agents and LLMs that need to understand, extend, or consume the Anexis Server API.

## Project Shape

Anexis Server is a Go API server using Gin, Gorm, PostgreSQL, and pluggable storage providers.

Main modules:

- `apps/api`: HTTP API, business features, storage providers, server entrypoint.
- `packages/database`: shared database config and Gorm models.
- `docs-site`: Docusaurus documentation site.

Go workspace:

- `go.work` includes `./apps/api` and `./packages/database`.
- Do not run `go test ./...` or `golangci-lint run ./...` from the repo root. Use explicit module paths:

```bash
go test ./apps/api/... ./packages/database/...
golangci-lint run ./apps/api/... ./packages/database/...
```

## API Base

Default local base URL:

```text
http://localhost:8080/api/v1
```

Health endpoint:

```text
GET /health
```

Swagger endpoint:

```text
GET /swagger/index.html
```

## Auth Modes

The server supports two modes:

### SaaS mode

Set:

```text
SERVER_MODE=saas
```

Auth flow:

1. Register with `POST /api/v1/auth/register`.
2. Login with `POST /api/v1/auth/login`.
3. Send `Authorization: Bearer <access_token>` to protected routes.

### Standalone mode

Set:

```text
SERVER_MODE=standalone
ANEXIS_API_KEY=<secret>
```

Protected routes should use:

```text
X-API-Key: <secret>
```

Implementation detail: route registration chooses different auth middleware in `apps/api/cmd/server/main.go`.

## Response Envelope

Successful responses use:

```json
{
  "success": true,
  "data": {}
}
```

Paginated/list responses can include:

```json
{
  "success": true,
  "data": [],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

Errors use:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable message",
    "details": "Optional validation details"
  }
}
```

Response helpers live in:

```text
apps/api/internal/infrastructure/http/response
```

Use those helpers instead of writing ad-hoc JSON.

## Route Map

All feature routes are registered under `/api/v1`.

### Auth

Defined in:

```text
apps/api/internal/features/auth
```

Routes:

| Method | Path | Auth | Notes |
| --- | --- | --- | --- |
| `POST` | `/auth/register` | Public in SaaS only | Create user. |
| `POST` | `/auth/login` | Public in SaaS only | Get access and refresh tokens. |
| `POST` | `/auth/refresh` | Public in SaaS only | Refresh access token. |
| `GET` | `/auth/me` | Protected | Get current user. |
| `PUT` | `/auth/password` | Protected | Change password. |

### Files

Defined in:

```text
apps/api/internal/features/files
```

Routes:

| Method | Path | Notes |
| --- | --- | --- |
| `GET` | `/files` | List files, supports pagination/search. |
| `POST` | `/files/upload` | Multipart upload field: `file`. |
| `POST` | `/files/folder` | Create folder. |
| `GET` | `/files/:id` | Get file metadata. |
| `GET` | `/files/:id/download` | Download file. |
| `PUT` | `/files/:id/rename` | Rename file/folder. |
| `PUT` | `/files/:id/move` | Move file/folder. |
| `DELETE` | `/files/:id` | Delete file/folder. |

### Links

Defined in:

```text
apps/api/internal/features/links
```

Routes:

| Method | Path | Auth | Notes |
| --- | --- | --- | --- |
| `GET` | `/links/:token/access` | Public | Access shared file. |
| `GET` | `/links/:token/stream` | Public | Get stream URL. |
| `GET` | `/links` | Protected | List links. |
| `POST` | `/links` | Protected | Create link. |
| `PUT` | `/links/:id` | Protected | Update link. |
| `DELETE` | `/links/:id` | Protected | Delete link. |

### Migration

Defined in:

```text
apps/api/internal/features/migration
```

Routes:

| Method | Path | Notes |
| --- | --- | --- |
| `GET` | `/migration` | List migration jobs. |
| `POST` | `/migration` | Start migration. |
| `GET` | `/migration/:id` | Get migration job. |
| `POST` | `/migration/:id/cancel` | Cancel migration job. |

Provider stubs live in:

```text
apps/api/internal/features/migration/providers
```

Current provider names: `google`, `amazon`, `microsoft`, `dropbox`.

### Backup

Defined in:

```text
apps/api/internal/features/backup
```

Routes:

| Method | Path | Notes |
| --- | --- | --- |
| `GET` | `/backup` | List backup jobs. |
| `POST` | `/backup/export` | Start export job. |
| `GET` | `/backup/:id` | Get backup job. |
| `GET` | `/backup/:id/download` | Get backup download URL. |

## Storage Providers

Interface:

```text
apps/api/internal/infrastructure/storage/provider.go
```

Implementations:

- `local`: filesystem storage.
- `s3`: S3-compatible storage.
- `backblaze`: Backblaze B2.

When adding a provider:

1. Implement `storage.Provider`.
2. Add config fields in `apps/api/internal/config/config.go`.
3. Wire creation in `apps/api/cmd/server/storage_factory.go`.
4. Add docs for required environment variables.
5. Add tests for any pure logic and adapter behavior that can be tested without cloud credentials.

## Database Models

Models live in:

```text
packages/database/models
```

Primary models:

- `User`
- `File`
- `Link`
- `MigrationJob`
- `BackupJob`

Repository code lives inside each feature package in `repository.go`.

## Implementation Conventions

When adding or changing an endpoint:

1. Add request/response DTOs in `dto.go`.
2. Add business behavior in `service.go`.
3. Add persistence in `repository.go`.
4. Add HTTP binding/error mapping in `handler.go`.
5. Register the route in `routes.go`.
6. Use response helpers from `internal/infrastructure/http/response`.
7. Return typed errors from services when handlers need specific HTTP mappings.
8. Add tests for DTO conversion, validation-independent logic, storage behavior, or repository behavior.
9. Update Swagger comments on handlers.
10. Regenerate Swagger docs if the project workflow requires it.

## Error Handling Rules

- Do not ignore returned errors unless the operation is explicitly best-effort.
- If a background job update fails, log the error.
- Use `errors.Is` in handlers to map known service errors to API responses.
- Never leak secrets in API responses or logs.

## Test Commands

Run all Go tests:

```bash
GOCACHE=/tmp/go-build-cache go test -v -race ./apps/api/... ./packages/database/...
```

Run lint:

```bash
golangci-lint run ./apps/api/... ./packages/database/...
```

Run coverage for packages with current unit coverage:

```bash
make test-coverage
```

Coverage artifacts are ignored by Git via:

```text
*.coverage*
```

## Docker Commands

Build API image:

```bash
docker build -t ghcr.io/gurren-software/anexis-server:latest -f apps/api/Dockerfile .
```

Run self-hosted stack:

```bash
docker compose -f docker-compose.selfhosted.yml up --build
```
