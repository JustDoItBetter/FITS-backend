# Swagger UI Guide

> Interactive API documentation and testing interface for FITS Backend.

## Overview

Swagger UI provides a web-based interface to explore and test the FITS Backend API. It automatically generates documentation from OpenAPI specifications and allows you to execute API calls directly from your browser.

## Accessing Swagger UI

### Primary URL

```
http://localhost:8080/docs
```

### Alternative URLs

- `http://localhost:8080/swagger/index.html` - Direct Swagger UI
- `http://localhost:8080/api` - Redirects to documentation

## Quick Start

### 1. Open Swagger UI

Navigate to `http://localhost:8080/docs` in your web browser.

### 2. Test Unauthenticated Endpoint

Try the health check endpoint first:

1. Find the `health` section
2. Locate `GET /health`
3. Click **"Try it out"**
4. Click **"Execute"**

Expected response (200 OK):

```json
{
  "status": "ok",
  "database": "connected",
  "time": "2025-10-22T12:00:00Z"
}
```

### 3. List Resources

Test the student list endpoint:

1. Find the `Students` section
2. Locate `GET /api/v1/student`
3. Click **"Try it out"**
4. Adjust pagination parameters (optional)
5. Click **"Execute"**

Response includes paginated student list.

## Authentication

Most endpoints require authentication. Here's how to authorize in Swagger UI.

### Step 1: Obtain Access Token

#### Option A: Bootstrap (First Time Only)

If the system is fresh:

1. Find `POST /api/v1/bootstrap/init`
2. Click **"Try it out"**
3. Click **"Execute"**
4. Copy the `admin_token` from response

**Note:** This only works once. After initial bootstrap, use login instead.

#### Option B: Login

For subsequent access:

1. Find `POST /api/v1/auth/login`
2. Click **"Try it out"**
3. Enter credentials:

```json
{
  "username": "your_username",
  "password": "your_password"
}
```

4. Click **"Execute"**
5. Copy the `access_token` from response

### Step 2: Authorize Swagger UI

1. Click the **"Authorize"** button at the top (lock icon)
2. Enter your token with `Bearer` prefix:
   ```
   Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
   ```
3. Click **"Authorize"**
4. Click **"Close"**

### Step 3: Test Protected Endpoints

Endpoints requiring authentication now have the lock icon unlocked. Try:

- `POST /api/v1/student` - Create student (admin only)
- `GET /api/v1/student/{uuid}` - Get student details
- `PUT /api/v1/student/{uuid}` - Update student (admin only)

## Common Workflows

### Creating a Teacher

1. Authorize with admin token
2. Find `POST /api/v1/admin/invite` under Invitations
3. Enter teacher details:

```json
{
  "email": "teacher@example.com",
  "first_name": "Jane",
  "last_name": "Smith",
  "role": "teacher",
  "department": "Computer Science"
}
```

4. Execute and copy `invitation_token`
5. Use `POST /api/v1/invite/{token}/complete` to finalize registration

### Creating a Student

1. Authorize with admin token
2. Find `POST /api/v1/admin/invite` under Invitations
3. Enter student details:

```json
{
  "email": "student@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "role": "student",
  "teacher_uuid": "teacher-uuid-here"
}
```

4. Execute and copy `invitation_token`
5. Use `POST /api/v1/invite/{token}/complete` to finalize registration

### Refreshing Tokens

When access token expires:

1. Find `POST /api/v1/auth/refresh`
2. Enter your refresh token:

```json
{
  "refresh_token": "your_refresh_token"
}
```

3. Execute to get new access token
4. Re-authorize with new token

## API Sections

### Authentication

- `POST /api/v1/bootstrap/init` - Initialize admin (once)
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - Logout (requires auth)
- `POST /api/v1/auth/refresh` - Refresh access token

### Invitations

- `POST /api/v1/admin/invite` - Create invitation (admin)
- `GET /api/v1/invite/{token}` - Get invitation details
- `POST /api/v1/invite/{token}/complete` - Complete registration

### Students

- `POST /api/v1/student` - Create student (admin)
- `GET /api/v1/student` - List students (paginated)
- `GET /api/v1/student/{uuid}` - Get student details
- `PUT /api/v1/student/{uuid}` - Update student (admin)
- `DELETE /api/v1/student/{uuid}` - Delete student (admin)

### Teachers

- `POST /api/v1/teacher` - Create teacher (admin)
- `GET /api/v1/teacher` - List teachers (paginated)
- `GET /api/v1/teacher/{uuid}` - Get teacher details
- `PUT /api/v1/teacher/{uuid}` - Update teacher (admin)
- `DELETE /api/v1/teacher/{uuid}` - Delete teacher (admin)

### System

- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics (requires secret)

## Pagination Parameters

List endpoints support pagination:

- `page` - Page number (default: 1)
- `page_size` - Items per page (default: 10, max: 100)
- `sort` - Sort field (optional)
- `order` - Sort order: `asc` or `desc` (optional)

Example:
```
GET /api/v1/student?page=2&page_size=20&sort=last_name&order=asc
```

## Response Formats

### Success Response

```json
{
  "success": true,
  "data": { ... }
}
```

### Error Response

```json
{
  "success": false,
  "error": "Error message",
  "details": "Detailed error description",
  "code": 400
}
```

### Paginated Response

```json
{
  "success": true,
  "data": {
    "items": [ ... ],
    "pagination": {
      "current_page": 1,
      "page_size": 10,
      "total_items": 50,
      "total_pages": 5
    }
  }
}
```

## Troubleshooting

### "Try it out" Button Not Working

**Problem:** Clicking "Try it out" does nothing

**Solution:** Ensure you're accessing via `http://localhost:8080` (NOT `http://0.0.0.0:8080`). Browsers treat `0.0.0.0` as untrustworthy.

### "Failed to fetch" Error

**Problem:** Request fails with "Failed to fetch"

**Possible causes:**

1. **Server not running**
   - Check: `curl http://localhost:8080/health`
   - Fix: `make run`

2. **CORS issue**
   - Verify you're using `localhost:8080` not `0.0.0.0:8080`
   - Clear browser cache (Ctrl+Shift+R)

3. **Network issue**
   - Check firewall settings
   - Verify port 8080 is not blocked

### "Unauthorized" Error (401)

**Problem:** Getting 401 errors on protected endpoints

**Solution:**

1. Verify you clicked **"Authorize"** button
2. Check token format includes `Bearer` prefix
3. Ensure token hasn't expired (access tokens expire after 1 hour)
4. Refresh token if needed: `POST /api/v1/auth/refresh`

### "Forbidden" Error (403)

**Problem:** Getting 403 errors

**Cause:** Insufficient permissions for your role

**Solution:**

- Ensure you're using admin token for admin-only endpoints
- Check role-based access control (RBAC) requirements
- Verify resource ownership (students/teachers can only access their own data)

### Bootstrap Already Initialized

**Problem:** Bootstrap returns "admin already initialized"

**This is normal** - Bootstrap can only run once. Use login instead:

```bash
POST /api/v1/auth/login
```

## Regenerating Documentation

If you modify API endpoints or annotations:

```bash
# Regenerate Swagger docs
make swagger

# Or manually
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

# Restart server
make run
```

Documentation updates automatically on server restart.

## Best Practices

### Testing

1. **Start with health check** - Verify server is running
2. **Test unauthenticated endpoints** - Ensure basic connectivity
3. **Obtain token** - Bootstrap or login
4. **Authorize Swagger UI** - Click Authorize button
5. **Test protected endpoints** - Try CRUD operations
6. **Check responses** - Verify expected data structure

### Security

1. **Never commit tokens** - Access tokens are sensitive
2. **Use localhost** - Not `0.0.0.0` for browser testing
3. **Refresh regularly** - Access tokens expire after 1 hour
4. **Logout when done** - Invalidate refresh tokens

### Development

1. **Keep docs in sync** - Regenerate after API changes
2. **Add examples** - Include request/response examples in annotations
3. **Document errors** - Specify possible error codes
4. **Version API** - Use `/api/v1/` prefix for versioning

## Related Documentation

- [Authentication Guide](authentication.md) - JWT and RBAC details
- [API Quick Start](../API_QUICK_START.md) - Using the API
- [Security Guide](security.md) - Security best practices
- [Testing Guide](testing.md) - Automated testing

## Technical Details

### Swagger Version

- **Swagger UI**: Latest (auto-updated)
- **OpenAPI**: 3.0 specification
- **Generator**: swag (github.com/swaggo/swag)

### Endpoints

- **Swagger UI**: `/docs`, `/swagger/index.html`
- **OpenAPI JSON**: `/swagger/doc.json`
- **OpenAPI YAML**: `/swagger/swagger.yaml`

### Configuration

Swagger is configured in `cmd/server/main.go`:

```go
// @title FITS Backend API
// @version 1.0
// @description REST API for FITS training management system
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

## License

This project is licensed under the GNU General Public License v3.0 - see [LICENSE](../../LICENSE) for details.

---

**Status:** Production Ready | **Last Updated:** 2025-10-22
