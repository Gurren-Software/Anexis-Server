---
sidebar_position: 1
slug: /api/auth
---

# Authentication API

## Register (SaaS Mode Only)

Register a new user account.

### Request

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepass123"
}
```

### Response

```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "name": "John Doe",
      "email": "john@example.com",
      "storage_used": 0,
      "storage_quota": 5368709120
    },
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

---

## Login

Authenticate and get access token.

### SaaS Mode (JWT)

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepass123"
}
```

### Standalone Mode (API Key)

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "X-API-Key: your-secret-api-key"
```

### Response

```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

---

## Refresh Token

Refresh an expired access token.

```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGc..."
}
```

### Response

```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

---

## Get Current User

```http
GET /api/v1/auth/me
Authorization: Bearer <access_token>
```

### Response

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john@example.com",
    "storage_used": 1048576,
    "storage_quota": 5368709120,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## Change Password

```http
PUT /api/v1/auth/password
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "current_password": "oldpass123",
  "new_password": "newpass456"
}
```

### Response

```json
{
  "success": true,
  "message": "Password changed successfully"
}
```