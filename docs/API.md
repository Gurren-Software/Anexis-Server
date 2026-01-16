# Anexis Server API Documentation

This document provides detailed information about the Anexis Server API.

## Base URL

- **Development**: `http://localhost:8080`
- **Production**: Your configured domain

## Authentication

All protected endpoints require a JWT token in the `Authorization` header:

```
Authorization: Bearer <your_access_token>
```

### Getting a Token

```bash
# Register a new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepass123"
  }'

# Login to get tokens
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123"
  }'

# Response:
# {
#   "access_token": "eyJhbG...",
#   "token_type": "Bearer",
#   "expires_in": 86400,
#   "refresh_token": "eyJhbG..."
# }
```

---

## Files API

### Upload a File

```bash
curl -X POST http://localhost:8080/api/v1/files/upload \
  -H "Authorization: Bearer <token>" \
  -F "file=@/path/to/document.pdf" \
  -F "compress=true" \
  -F "description=My important document"
```

### List Files

```bash
curl http://localhost:8080/api/v1/files \
  -H "Authorization: Bearer <token>"

# With pagination and search
curl "http://localhost:8080/api/v1/files?page=1&per_page=20&search=document" \
  -H "Authorization: Bearer <token>"
```

### Download a File

```bash
curl -O http://localhost:8080/api/v1/files/123/download \
  -H "Authorization: Bearer <token>"
```

### Create a Folder

```bash
curl -X POST http://localhost:8080/api/v1/files/folder \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Documents",
    "parent_id": null
  }'
```

---

## Links API

### Create an Access Link

```bash
# Permanent public link
curl -X POST http://localhost:8080/api/v1/links \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "file_id": 123,
    "type": "permanent",
    "access_type": "public"
  }'

# Temporal password-protected link
curl -X POST http://localhost:8080/api/v1/links \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "file_id": 123,
    "type": "temporal",
    "access_type": "restricted",
    "password": "secretpass",
    "expires_in": 86400,
    "max_downloads": 5
  }'
```

### Access a File via Link

```bash
# Public link
curl -O "http://localhost:8080/api/v1/links/abc123token/access"

# Password-protected link
curl -O "http://localhost:8080/api/v1/links/abc123token/access?password=secretpass"
```

### Get Streaming URL

```bash
curl "http://localhost:8080/api/v1/links/abc123token/stream"

# Response:
# { "stream_url": "https://backblaze-url..." }
```

---

## Migration API

### Start a Migration from Google Drive

```bash
curl -X POST http://localhost:8080/api/v1/migration \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "google",
    "access_token": "<google_oauth_token>",
    "refresh_token": "<google_refresh_token>"
  }'
```

### Check Migration Status

```bash
curl http://localhost:8080/api/v1/migration/1 \
  -H "Authorization: Bearer <token>"

# Response:
# {
#   "id": 1,
#   "provider": "google",
#   "status": "running",
#   "total_files": 150,
#   "processed_files": 75,
#   "progress": 50.0
# }
```

### Cancel a Migration

```bash
curl -X POST http://localhost:8080/api/v1/migration/1/cancel \
  -H "Authorization: Bearer <token>"
```

---

## Backup API

### Export All Data

```bash
curl -X POST http://localhost:8080/api/v1/backup/export \
  -H "Authorization: Bearer <token>"

# Response:
# {
#   "id": 1,
#   "type": "export",
#   "status": "pending"
# }
```

### Download Backup Archive

```bash
# Get download URL
curl http://localhost:8080/api/v1/backup/1/download \
  -H "Authorization: Bearer <token>"

# Response:
# { "download_url": "https://..." }

# Download the archive
curl -O "<download_url>"
```

---

## Error Responses

All errors follow this format:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable message",
    "details": "Additional details (optional)"
  }
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Missing or invalid token |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `VALIDATION_ERROR` | 422 | Invalid request data |
| `INTERNAL_ERROR` | 500 | Server error |
| `QUOTA_EXCEEDED` | 400 | Storage quota exceeded |
| `LINK_EXPIRED` | 400 | Access link has expired |
| `DOWNLOAD_LIMIT` | 400 | Download limit reached |

---

## Rate Limiting

- **Default**: 100 requests per minute per IP
- **File Downloads**: No rate limit
- Headers included in response:
  - `X-RateLimit-Limit`
  - `X-RateLimit-Remaining`
  - `X-RateLimit-Reset`

---

## Swagger UI

Interactive API documentation is available at:

```
http://localhost:8080/swagger/index.html
```

Generate/update Swagger docs:

```bash
make swagger
```
