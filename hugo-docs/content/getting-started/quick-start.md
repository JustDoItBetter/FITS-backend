---
title: "Quick Start"
weight: 2
---

# Quick Start Guide

Get your first API requests running in 5 minutes.

## Prerequisites

- FITS Backend running (see [Installation Guide](/getting-started/installation/))
- `curl` or Postman installed
- Server accessible at `http://localhost:8080`

## Step 1: Bootstrap Admin User

The first step is to create an admin user. This can only be done once when no users exist.

```bash
curl -X POST http://localhost:8080/api/v1/bootstrap/init \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "SecurePassword123!",
    "first_name": "Admin",
    "last_name": "User"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Bootstrap admin created successfully",
  "data": {
    "user": {
      "id": 1,
      "email": "admin@example.com",
      "role": "admin",
      "is_active": true
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Important:**
- Save the `token` - you'll need it for authenticated requests
- This endpoint can only be called once
- Subsequent calls will return an error

## Step 2: Login

Authenticate with your admin credentials:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "SecurePassword123!"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "admin@example.com",
      "role": "admin"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2025-10-23T18:15:00Z"
  }
}
```

## Step 3: Set Authentication Token

Export the token for subsequent requests:

```bash
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

Or in PowerShell:
```powershell
$TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## Step 4: Create an Invitation

As admin, create an invitation for a new teacher:

```bash
curl -X POST http://localhost:8080/api/v1/admin/invite \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "email": "teacher@example.com",
    "role": "teacher"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "invitation": {
      "id": 1,
      "email": "teacher@example.com",
      "role": "teacher",
      "token": "inv_abc123def456",
      "expires_at": "2025-10-30T18:00:00Z"
    }
  }
}
```

**Note:** Save the invitation `token` to share with the new user.

## Step 5: Register with Invitation

The invited user can now register:

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "invitation_token": "inv_abc123def456",
    "password": "TeacherPass123!",
    "first_name": "John",
    "last_name": "Teacher"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 2,
      "email": "teacher@example.com",
      "role": "teacher"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

## Step 6: Create a Student

Now let's create a student record:

```bash
curl -X POST http://localhost:8080/api/v1/student \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "first_name": "Jane",
    "last_name": "Student",
    "email": "jane@example.com",
    "student_id": "S001"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "first_name": "Jane",
    "last_name": "Student",
    "email": "jane@example.com",
    "student_id": "S001",
    "created_at": "2025-10-23T18:00:00Z",
    "updated_at": "2025-10-23T18:00:00Z"
  }
}
```

## Step 7: List Students

Retrieve all students:

```bash
curl -X GET "http://localhost:8080/api/v1/student?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "first_name": "Jane",
      "last_name": "Student",
      "email": "jane@example.com",
      "student_id": "S001",
      "created_at": "2025-10-23T18:00:00Z",
      "updated_at": "2025-10-23T18:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

## Step 8: Get Single Student

Retrieve a specific student:

```bash
curl -X GET http://localhost:8080/api/v1/student/1 \
  -H "Authorization: Bearer $TOKEN"
```

## Step 9: Update Student

Update student information:

```bash
curl -X PUT http://localhost:8080/api/v1/student/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "first_name": "Jane",
    "last_name": "Smith",
    "email": "jane.smith@example.com",
    "student_id": "S001"
  }'
```

## Step 10: Check Health

Verify the system is healthy:

```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok",
  "database": "connected",
  "time": "2025-10-23T18:00:00Z"
}
```

## Complete Example Script

Here's a complete bash script for the quick start:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "1. Bootstrap Admin User"
BOOTSTRAP_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/bootstrap/init \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "SecurePassword123!",
    "first_name": "Admin",
    "last_name": "User"
  }')

echo $BOOTSTRAP_RESPONSE | jq '.'
TOKEN=$(echo $BOOTSTRAP_RESPONSE | jq -r '.data.token')

echo -e "\n2. Create Invitation"
INVITE_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/admin/invite \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "email": "teacher@example.com",
    "role": "teacher"
  }')

echo $INVITE_RESPONSE | jq '.'
INVITATION_TOKEN=$(echo $INVITE_RESPONSE | jq -r '.data.invitation.token')

echo -e "\n3. Register Teacher"
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d "{
    \"invitation_token\": \"$INVITATION_TOKEN\",
    \"password\": \"TeacherPass123!\",
    \"first_name\": \"John\",
    \"last_name\": \"Teacher\"
  }")

echo $REGISTER_RESPONSE | jq '.'

echo -e "\n4. Create Student"
STUDENT_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/student \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "first_name": "Jane",
    "last_name": "Student",
    "email": "jane@example.com",
    "student_id": "S001"
  }')

echo $STUDENT_RESPONSE | jq '.'

echo -e "\n5. List Students"
curl -s -X GET "$BASE_URL/api/v1/student?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo -e "\n6. Check Health"
curl -s $BASE_URL/health | jq '.'

echo -e "\nQuick start completed!"
```

Save as `quick-start.sh`, make executable, and run:

```bash
chmod +x quick-start.sh
./quick-start.sh
```

## Using Postman

### Import Collection

The repository includes a Postman collection at `docs/FITS-Backend-API.postman_collection.json`.

1. Open Postman
2. Import → Upload Files
3. Select `docs/FITS-Backend-API.postman_collection.json`

### Set Environment

Create a Postman environment with:

```json
{
  "name": "FITS Local",
  "values": [
    {
      "key": "base_url",
      "value": "http://localhost:8080",
      "enabled": true
    },
    {
      "key": "token",
      "value": "",
      "enabled": true
    }
  ]
}
```

The collection automatically sets the `token` variable after login.

## Interactive API Documentation

The easiest way to explore the API is through Swagger UI:

1. Open browser: http://localhost:8080/docs
2. Click "Authorize" button
3. Enter: `Bearer YOUR_TOKEN_HERE`
4. Try out any endpoint directly in the browser

## Common Operations

### Refresh Token

When your access token expires:

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

### Logout

Invalidate your current session:

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer $TOKEN"
```

### Logout All Sessions

Invalidate all sessions for your user:

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout-all \
  -H "Authorization: Bearer $TOKEN"
```

### Delete Student

```bash
curl -X DELETE http://localhost:8080/api/v1/student/1 \
  -H "Authorization: Bearer $TOKEN"
```

## Error Handling

All errors follow a consistent format:

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid email format",
    "details": {
      "field": "email",
      "value": "invalid-email"
    }
  }
}
```

Common HTTP status codes:
- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Validation error
- `401 Unauthorized` - Missing or invalid token
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

## Next Steps

- [API Reference](/api/endpoints/) - Detailed endpoint documentation
- [Authentication Guide](/api/authentication/) - Deep dive into auth flow
- [Development Setup](/development/local-setup/) - Set up your dev environment
- [Swagger UI](http://localhost:8080/docs) - Interactive API explorer

## Troubleshooting

### Token Expired

Access tokens expire after 15 minutes. Use the refresh token to get a new one.

### 401 Unauthorized

Make sure you're including the `Authorization: Bearer {token}` header.

### 403 Forbidden

Your role doesn't have permission for this endpoint. Check RBAC rules.

### 429 Rate Limit Exceeded

You've exceeded the rate limit. Wait a minute and try again.

### Bootstrap Already Initialized

The bootstrap endpoint can only be called once. Use the login endpoint instead.
