---
sidebar_position: 4
slug: /api/migration
---

# Migration API

Import files from external cloud providers.

## Supported Providers

- Google Drive
- Amazon S3
- Microsoft OneDrive
- Dropbox

## Start Migration

```http
POST /api/v1/migration
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "provider": "google_drive",
  "credentials": {
    "access_token": "ya29.a0..."
  }
}
```

### Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `provider` | string | `google_drive`, `amazon_s3`, `microsoft`, `dropbox` |
| `credentials` | object | Provider-specific credentials |

### Response

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "provider": "google_drive",
    "status": "pending",
    "progress": 0,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## List Migrations

```http
GET /api/v1/migration
Authorization: Bearer <access_token>
```

### Response

```json
{
  "success": true,
  "data": {
    "jobs": [
      {
        "id": "uuid",
        "provider": "google_drive",
        "status": "completed",
        "progress": 100,
        "files_imported": 50,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

## Get Migration Status

```http
GET /api/v1/migration/:id
Authorization: Bearer <access_token>
```

---

## Cancel Migration

```http
POST /api/v1/migration/:id/cancel
Authorization: Bearer <access_token>
```