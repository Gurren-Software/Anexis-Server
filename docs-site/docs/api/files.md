---
sidebar_position: 2
slug: /api/files
---

# Files API

## Upload File

Upload a file to storage.

### Request (SaaS Mode)

```bash
curl -X POST http://localhost:8080/api/v1/files/upload \
  -H "Authorization: Bearer <access_token>" \
  -F "file=@/path/to/document.pdf"
```

### Request (Standalone Mode)

```bash
curl -X POST http://localhost:8080/api/v1/files/upload \
  -H "X-API-Key: your-api-key" \
  -F "file=@/path/to/document.pdf"
```

### Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `file` | file | The file to upload (multipart/form-data) |
| `parent_id` | uuid | Parent folder ID (optional) |
| `compress` | bool | Enable compression (optional) |

### Response

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "document.pdf",
    "mime_type": "application/pdf",
    "size": 1048576,
    "parent_id": null,
    "storage_key": "user-id/uuid",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## List Files

List files and folders.

```http
GET /api/v1/files?parent_id=<uuid>&page=1&per_page=20
Authorization: Bearer <access_token>
```

### Query Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `parent_id` | uuid | Filter by parent folder |
| `page` | int | Page number (default: 1) |
| `per_page` | int | Items per page (default: 20) |
| `search` | string | Search by name |

### Response

```json
{
  "success": true,
  "data": {
    "files": [
      {
        "id": "uuid",
        "name": "document.pdf",
        "mime_type": "application/pdf",
        "size": 1048576,
        "is_folder": false,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "per_page": 20,
      "total": 1
    }
  }
}
```

---

## Get File

Get file details.

```http
GET /api/v1/files/:id
Authorization: Bearer <access_token>
```

### Response

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "document.pdf",
    "mime_type": "application/pdf",
    "size": 1048576,
    "parent_id": null,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## Download File

Download a file.

```bash
curl -O http://localhost:8080/api/v1/files/:id/download \
  -H "Authorization: Bearer <access_token>"
```

---

## Create Folder

Create a new folder.

```http
POST /api/v1/files/folder
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "My Documents",
  "parent_id": null
}
```

### Response

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "My Documents",
    "is_folder": true,
    "parent_id": null,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## Rename File/Folder

```http
PUT /api/v1/files/:id/rename
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "New Name.pdf"
}
```

---

## Move File/Folder

```http
PUT /api/v1/files/:id/move
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "parent_id": "new-parent-uuid"
}
```

---

## Delete File/Folder

```http
DELETE /api/v1/files/:id
Authorization: Bearer <access_token>
```

### Response

```json
{
  "success": true,
  "message": "File deleted successfully"
}
```