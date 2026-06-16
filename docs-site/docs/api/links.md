---
sidebar_position: 3
slug: /api/links
---

# Links API

Create and manage access links for sharing files.

## Create Link

Create a new access link.

```http
POST /api/v1/links
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "file_id": "uuid",
  "type": "download",
  "access_type": "public",
  "expires_at": "2024-12-31T23:59:59Z"
}
```

### Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `file_id` | uuid | The file to share |
| `type` | string | `download`, `stream`, `permanent`, `temporal` |
| `access_type` | string | `public`, `password` |
| `password` | string | Password for password-protected links |
| `expires_at` | datetime | Expiration time (optional) |

### Response

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "token": "abc123...",
    "url": "http://localhost:8080/api/v1/links/abc123.../access",
    "type": "download",
    "access_type": "public",
    "expires_at": "2024-12-31T23:59:59Z",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## List Links

```http
GET /api/v1/links
Authorization: Bearer <access_token>
```

### Response

```json
{
  "success": true,
  "data": {
    "links": [
      {
        "id": "uuid",
        "token": "abc123...",
        "file_id": "uuid",
        "type": "download",
        "access_type": "public",
        "expires_at": "2024-12-31T23:59:59Z",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

## Update Link

```http
PUT /api/v1/links/:id
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "access_type": "password",
  "password": "newpassword"
}
```

---

## Delete Link

```http
DELETE /api/v1/links/:id
Authorization: Bearer <access_token>
```

---

## Access Link (Public)

Access a file via a link token (no auth required).

```http
GET /api/v1/links/:token/access
```

### Response

Returns the file content directly.

---

## Stream Link (Public)

Get a streaming URL for a file.

```http
GET /api/v1/links/:token/stream
```

### Response

```json
{
  "success": true,
  "data": {
    "stream_url": "http://localhost:8080/api/v1/storage/..."
  }
}
```