---
sidebar_position: 5
slug: /api/backup
---

# Backup API

Export and restore your data.

## Start Export

Create a backup export of all your files.

```http
POST /api/v1/backup/export
Authorization: Bearer <access_token>
```

### Response

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "type": "export",
    "status": "pending",
    "progress": 0,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## List Backups

```http
GET /api/v1/backup
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
        "type": "export",
        "status": "completed",
        "progress": 100,
        "file_size": 104857600,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

## Get Backup Status

```http
GET /api/v1/backup/:id
Authorization: Bearer <access_token>
```

---

## Download Backup

Get the download URL for a completed backup.

```http
GET /api/v1/backup/:id/download
Authorization: Bearer <access_token>
```

### Response

```json
{
  "success": true,
  "data": {
    "download_url": "http://localhost:8080/api/v1/storage/...",
    "expires_at": "2024-01-02T00:00:00Z"
  }
}
```