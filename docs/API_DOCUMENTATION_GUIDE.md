# API Documentation Guide

## Overview

FITS Backend provides comprehensive API documentation using **Swagger UI** - automatically generated from code annotations. No manual HTML files are required; documentation is kept in sync with code.

## Accessing Documentation

**URL:** `http://localhost:8080/docs`

The documentation is automatically generated from Go code annotations using the `swag` tool.

## Features

- **Code-first documentation** - annotations in Go code generate OpenAPI spec
- **Interactive API testing** - test endpoints directly in browser
- **Automatic synchronization** - docs always match code
- **No manual HTML** - everything generated from source
- **Complete API coverage** - all endpoints, models, and responses documented
- **Authentication support** - test protected endpoints with JWT tokens

## Quick Start

### Step 1: Start the Server

```bash
make run
# or
go run cmd/server/main.go
```

### Step 2: Bootstrap Admin (First Run)

```bash
curl -X POST http://localhost:8080/api/v1/bootstrap/init \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "SecurePassword123!",
    "email": "admin@fits.example.com"
  }'
```

### Step 3: Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "SecurePassword123!"
  }'
```

Save the `access_token` from the response.

### Step 4: Explore the API

Open your browser and navigate to `http://localhost:8080/docs` to view the complete API documentation.

## Authentication

All protected endpoints require a JWT Bearer token:

```
Authorization: Bearer <your-access-token>
```

### Getting a Token

1. **Bootstrap** (first time only):
   - Endpoint: `POST /api/v1/bootstrap/init`
   - Creates admin account
   - Can only be called once

2. **Login**:
   - Endpoint: `POST /api/v1/auth/login`
   - Returns `access_token` (valid 1h) and `refresh_token` (valid 30 days)

3. **Refresh**:
   - Endpoint: `POST /api/v1/auth/refresh`
   - Use `refresh_token` to get new `access_token`

## API Coverage

The documentation covers:

- Authentication & Authorization
- Student Management (CRUD + Pagination)
- Teacher Management (CRUD + Pagination)
- Invitation System
- Health Checks

## Updating Documentation

Documentation is auto-generated from Go code annotations. The process is:

### 1. Add Annotations to Code

All endpoints use special comments that `swag` reads:

```go
// @Summary Create a student
// @Description Creates a new student in the system
// @Tags students
// @Accept json
// @Produce json
// @Param student body CreateStudentRequest true "Student data"
// @Success 201 {object} StudentResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/student [post]
// @Security BearerAuth
func (h *Handler) Create(c *fiber.Ctx) error {
    // implementation...
}
```

### 2. Regenerate Documentation

After adding/modifying annotations:

```bash
# Regenerate OpenAPI specification from code
~/go/bin/swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

# Or use Makefile
make swagger
```

This generates `docs/swagger.json` and `docs/swagger.yaml` from your code annotations.

### 3. View Changes

Restart the server and visit `http://localhost:8080/docs` to see updated documentation.

```bash
make run
```

The Swagger UI automatically reads the generated specification - no manual HTML editing needed.

## API Endpoints

### Authentication
- `POST /api/v1/bootstrap/init` - Initialize admin account
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/logout` - Logout
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/admin/invite` - Create invitation (Admin only)

### Students
- `POST /api/v1/student` - Create student (Admin only)
- `GET /api/v1/student/:uuid` - Get student by UUID
- `PUT /api/v1/student/:uuid` - Update student (Admin only)
- `DELETE /api/v1/student/:uuid` - Delete student (Admin only)
- `GET /api/v1/student` - List all students (with pagination)

### Teachers
- `POST /api/v1/teacher` - Create teacher (Admin only)
- `GET /api/v1/teacher/:uuid` - Get teacher by UUID
- `PUT /api/v1/teacher/:uuid` - Update teacher (Admin only)
- `DELETE /api/v1/teacher/:uuid` - Delete teacher (Admin only)
- `GET /api/v1/teacher` - List all teachers (with pagination)

### System
- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics (requires secret)

## Troubleshooting

### "401 Unauthorized" errors
- Check if your token is expired (1h validity)
- Verify Bearer token format: `Bearer <token>`
- Use refresh token to get new access token

### "CORS errors" in browser
- Check `allowed_origins` in `config.toml`
- Default: `*` (all origins for development)

### Documentation not updating
- Regenerate Swagger docs with swag
- Clear browser cache
- Restart server

## Additional Resources

- **Main README:** `../README.md`
- **API Endpoints:** `docs/api/API.md`
- **Contributing:** `../CONTRIBUTING.md`
- **Architecture:** `docs/architecture/ARCHITECTURE.md`

## How It Works

FITS Backend uses a **code-first documentation approach**:

1. **Developers annotate code** with special comments (e.g., `@Summary`, `@Router`)
2. **`swag` tool scans code** and generates OpenAPI 3.0 specification
3. **Swagger UI renders** the specification as interactive documentation
4. **Documentation stays in sync** with code automatically

This eliminates manual HTML writing and ensures docs never go out of sync with implementation.

## Summary

FITS Backend uses **Swagger UI with code-generated documentation**. All endpoints are documented through Go code annotations, ensuring documentation accuracy and eliminating manual HTML maintenance.

**Start here:** `http://localhost:8080/docs`
