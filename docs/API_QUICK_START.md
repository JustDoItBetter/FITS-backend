# FITS Backend API - Quick Start Guide

**Version**: 1.0.0
**Documentation**: [OpenAPI Specification](./openapi.yaml)
**Swagger UI**: http://localhost:8080/swagger/index.html

---

## Table of Contents

1. [Getting Started](#getting-started)
2. [Authentication](#authentication)
3. [Common Workflows](#common-workflows)
4. [API Examples](#api-examples)
5. [Error Handling](#error-handling)
6. [Rate Limiting](#rate-limiting)
7. [Testing Tools](#testing-tools)

---

## Getting Started

### Base URLs

| Environment | URL |
|-------------|-----|
| **Local** | http://localhost:8080 |
| **Staging** | https://staging-api.fits.example.com |
| **Production** | https://api.fits.example.com |

### Quick Health Check

```bash
curl http://localhost:8080/health
```

**Response**:
```json
{
  "status": "ok",
  "database": "connected",
  "time": "2025-10-22T14:30:00Z"
}
```

---

## Authentication

### 1. Bootstrap Admin (One-Time Setup)

```bash
curl -X POST http://localhost:8080/api/v1/bootstrap/init
```

**Response**:
```json
{
  "success": true,
  "message": "Admin certificate generated successfully",
  "data": {
    "admin_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "message": "Admin certificate generated successfully",
    "public_key_path": "./configs/keys/admin.pub"
  }
}
```

**⚠️ Save the `admin_token` securely! This can only be done once.**

---

### 2. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "max.mustermann",
    "password": "SecurePassword123!"
  }'
```

**Response**:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "role": "student",
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

---

### 3. Using Access Tokens

Include the access token in the `Authorization` header for protected endpoints:

```bash
curl http://localhost:8080/api/v1/student \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

### 4. Refresh Tokens

When your access token expires (after 1 hour), use the refresh token:

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

---

### 5. Logout

Invalidate all refresh tokens:

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## Common Workflows

### Workflow 1: Create a New Teacher (Admin)

```bash
# Step 1: Login as admin
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"AdminPassword123!"}' \
  | jq -r '.data.access_token')

# Step 2: Create invitation for teacher
INVITE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/admin/invite \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "anna.schmidt@example.com",
    "first_name": "Anna",
    "last_name": "Schmidt",
    "role": "teacher",
    "department": "Computer Science"
  }')

# Extract invitation link
INVITE_LINK=$(echo $INVITE_RESPONSE | jq -r '.data.invitation_link')
echo "Send this link to the teacher: $INVITE_LINK"
```

---

### Workflow 2: Complete Teacher Registration

```bash
# Teacher receives invitation link and extracts token
INVITE_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Step 1: View invitation details
curl -s http://localhost:8080/api/v1/invite/$INVITE_TOKEN | jq

# Step 2: Complete registration
curl -X POST http://localhost:8080/api/v1/invite/$INVITE_TOKEN/complete \
  -H "Content-Type: application/json" \
  -d '{
    "username": "anna.schmidt",
    "password": "TeacherPass456!"
  }'

# Step 3: Login with new credentials
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "anna.schmidt",
    "password": "TeacherPass456!"
  }'
```

---

### Workflow 3: Create a Student

```bash
# Login as admin
TOKEN="YOUR_ADMIN_TOKEN"

# Create student
curl -X POST http://localhost:8080/api/v1/student \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "max.mustermann@example.com",
    "first_name": "Max",
    "last_name": "Mustermann",
    "teacher_id": "550e8400-e29b-41d4-a716-446655440010"
  }'
```

---

### Workflow 4: List Students with Pagination

```bash
# Get first page (20 items)
curl "http://localhost:8080/api/v1/student?page=1&limit=20" | jq

# Get second page
curl "http://localhost:8080/api/v1/student?page=2&limit=20" | jq

# Get with custom limit
curl "http://localhost:8080/api/v1/student?page=1&limit=50" | jq
```

---

## API Examples

### Create Teacher

**Request**:
```bash
curl -X POST http://localhost:8080/api/v1/teacher \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "peter.mueller@example.com",
    "first_name": "Peter",
    "last_name": "Mueller",
    "department": "Mathematics"
  }'
```

**Response**:
```json
{
  "success": true,
  "message": "Teacher created successfully",
  "data": {
    "uuid": "550e8400-e29b-41d4-a716-446655440010",
    "email": "peter.mueller@example.com",
    "first_name": "Peter",
    "last_name": "Mueller",
    "department": "Mathematics",
    "created_at": "2025-10-22T14:30:00Z",
    "updated_at": "2025-10-22T14:30:00Z"
  }
}
```

---

### Update Student

**Request**:
```bash
curl -X PUT http://localhost:8080/api/v1/student/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "max.new@example.com",
    "first_name": "Maximilian"
  }'
```

**Response**:
```json
{
  "success": true,
  "message": "Student updated successfully",
  "data": {
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "email": "max.new@example.com",
    "first_name": "Maximilian",
    "last_name": "Mustermann",
    "teacher_id": "550e8400-e29b-41d4-a716-446655440010",
    "created_at": "2025-09-30T12:00:00Z",
    "updated_at": "2025-10-22T14:35:00Z"
  }
}
```

---

### Delete Teacher

**Request**:
```bash
curl -X DELETE http://localhost:8080/api/v1/teacher/550e8400-e29b-41d4-a716-446655440010 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response**: `204 No Content` (successful deletion)

---

### Get Student by UUID

**Request**:
```bash
curl http://localhost:8080/api/v1/student/550e8400-e29b-41d4-a716-446655440000
```

**Response**:
```json
{
  "success": true,
  "data": {
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "email": "max.mustermann@example.com",
    "first_name": "Max",
    "last_name": "Mustermann",
    "teacher_id": "550e8400-e29b-41d4-a716-446655440010",
    "created_at": "2025-09-30T12:00:00Z",
    "updated_at": "2025-09-30T12:00:00Z"
  }
}
```

---

## Error Handling

### Common Error Responses

#### 400 Bad Request
```json
{
  "success": false,
  "code": 400,
  "error": "Bad Request",
  "details": "Invalid UUID format"
}
```

#### 401 Unauthorized
```json
{
  "success": false,
  "code": 401,
  "error": "Unauthorized",
  "details": "Invalid or expired token"
}
```

#### 403 Forbidden
```json
{
  "success": false,
  "code": 403,
  "error": "Forbidden",
  "details": "Admin role required"
}
```

#### 404 Not Found
```json
{
  "success": false,
  "code": 404,
  "error": "Not Found",
  "details": "Student not found"
}
```

#### 409 Conflict
```json
{
  "success": false,
  "code": 409,
  "error": "Conflict",
  "details": "Email already exists"
}
```

#### 422 Validation Error
```json
{
  "success": false,
  "code": 422,
  "error": "Validation Error",
  "details": "Password must be at least 8 characters long"
}
```

#### 500 Internal Server Error
```json
{
  "success": false,
  "code": 500,
  "error": "Internal Server Error",
  "details": "An unexpected error occurred"
}
```

---

## Rate Limiting

### Global Rate Limits (Per IP)
- **100 requests per minute** for all endpoints
- Returns `429 Too Many Requests` when exceeded

### Per-User Rate Limits (Authenticated)
| Role | Limit |
|------|-------|
| **Admin** | 200 requests/minute |
| **Teacher** | 100 requests/minute |
| **Student** | 50 requests/minute |

### Rate Limit Response
```json
{
  "success": false,
  "error": "Rate limit exceeded",
  "details": "Maximum 100 requests per minute allowed",
  "code": 429
}
```

---

## Testing Tools

### Using Postman

1. **Import OpenAPI Spec**:
   - File → Import
   - Select `docs/openapi.yaml`
   - All endpoints imported with examples

2. **Set Environment Variables**:
   ```
   BASE_URL: http://localhost:8080
   ACCESS_TOKEN: (set after login)
   REFRESH_TOKEN: (set after login)
   ```

3. **Collection Pre-request Script**:
   ```javascript
   // Auto-refresh expired tokens
   if (pm.environment.get("ACCESS_TOKEN") && isTokenExpired()) {
       refreshToken();
   }
   ```

---

### Using cURL Scripts

**Save tokens to file**:
```bash
# Login and save tokens
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"max.mustermann","password":"SecurePassword123!"}' \
  | jq -r '.data | "ACCESS_TOKEN=\(.access_token)\nREFRESH_TOKEN=\(.refresh_token)"' \
  > ~/.fits_tokens

# Source tokens
source ~/.fits_tokens

# Use token
curl http://localhost:8080/api/v1/student \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

---

### Using HTTPie

```bash
# Install HTTPie
pip install httpie

# Login
http POST http://localhost:8080/api/v1/auth/login \
  username=max.mustermann password=SecurePassword123!

# Use token
http GET http://localhost:8080/api/v1/student \
  Authorization:"Bearer YOUR_ACCESS_TOKEN"
```

---

### Using Swagger UI (Built-in)

1. Start the API server
2. Open http://localhost:8080/swagger/index.html
3. Click "Authorize" button
4. Enter: `Bearer YOUR_ACCESS_TOKEN`
5. Try out any endpoint with pre-filled examples

---

## Complete Example: Student Management Flow

```bash
#!/bin/bash
# Complete student management workflow

BASE_URL="http://localhost:8080"

# 1. Bootstrap admin (only needed once)
echo "=== Bootstrap Admin ==="
ADMIN_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/bootstrap/init)
ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | jq -r '.data.admin_token')
echo "Admin token: $ADMIN_TOKEN"

# 2. Create teacher invitation
echo -e "\n=== Create Teacher Invitation ==="
TEACHER_INVITE=$(curl -s -X POST $BASE_URL/api/v1/admin/invite \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "anna.schmidt@example.com",
    "first_name": "Anna",
    "last_name": "Schmidt",
    "role": "teacher",
    "department": "Computer Science"
  }')
INVITE_TOKEN=$(echo $TEACHER_INVITE | jq -r '.data.invitation_token')
echo "Invitation created: $(echo $TEACHER_INVITE | jq -r '.data.invitation_link')"

# 3. Teacher completes registration
echo -e "\n=== Complete Teacher Registration ==="
curl -s -X POST $BASE_URL/api/v1/invite/$INVITE_TOKEN/complete \
  -H "Content-Type: application/json" \
  -d '{
    "username": "anna.schmidt",
    "password": "TeacherPass456!"
  }' | jq

# 4. Teacher logs in
echo -e "\n=== Teacher Login ==="
TEACHER_LOGIN=$(curl -s -X POST $BASE_URL/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "anna.schmidt",
    "password": "TeacherPass456!"
  }')
TEACHER_TOKEN=$(echo $TEACHER_LOGIN | jq -r '.data.access_token')
TEACHER_UUID=$(echo $TEACHER_LOGIN | jq -r '.data.user_id')
echo "Teacher UUID: $TEACHER_UUID"

# 5. Admin creates student
echo -e "\n=== Create Student ==="
STUDENT=$(curl -s -X POST $BASE_URL/api/v1/student \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"max.mustermann@example.com\",
    \"first_name\": \"Max\",
    \"last_name\": \"Mustermann\",
    \"teacher_id\": \"$TEACHER_UUID\"
  }")
STUDENT_UUID=$(echo $STUDENT | jq -r '.data.uuid')
echo "Student created: $STUDENT_UUID"

# 6. List all students
echo -e "\n=== List Students ==="
curl -s "$BASE_URL/api/v1/student?page=1&limit=10" | jq

# 7. Get specific student
echo -e "\n=== Get Student Details ==="
curl -s "$BASE_URL/api/v1/student/$STUDENT_UUID" | jq

# 8. Update student
echo -e "\n=== Update Student ==="
curl -s -X PUT "$BASE_URL/api/v1/student/$STUDENT_UUID" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "max.new@example.com"
  }' | jq

echo -e "\n=== Workflow Complete ==="
```

---

## Troubleshooting

### Issue: 401 Unauthorized

**Solution**: Check token expiry and refresh if needed
```bash
# Check token expiry
echo $ACCESS_TOKEN | jwt decode -

# Refresh token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}"
```

---

### Issue: 403 Forbidden

**Solution**: Verify you have the required role (admin, teacher, student)
```bash
# Check your role in the token
echo $ACCESS_TOKEN | jwt decode - | jq '.role'
```

---

### Issue: 422 Validation Error

**Solution**: Check request body matches schema requirements
- Email must be valid format
- Password must be at least 8 characters
- UUIDs must be valid format
- Required fields must not be empty

---

### Issue: 429 Rate Limit Exceeded

**Solution**: Wait 60 seconds or use authenticated requests with higher limits

---

## Additional Resources

- **OpenAPI Spec**: [docs/openapi.yaml](./openapi.yaml)
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health
- **Improvement Report**: [docs/OPENAPI_IMPROVEMENTS_REPORT.md](./OPENAPI_IMPROVEMENTS_REPORT.md)

---

**Last Updated**: 2025-10-22
**API Version**: 1.0.0
**OpenAPI Version**: 3.0.3
