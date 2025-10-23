---
title: "System Design"
weight: 2
---

# System Design

## Design Philosophy

FITS Backend is designed with the following principles:

### 1. Simplicity First

- Clear, readable code over clever optimizations
- Explicit over implicit behavior
- Standard library where possible

### 2. Security by Default

- All endpoints require authentication unless explicitly public
- Defense in depth with multiple security layers
- Secure defaults in configuration

### 3. Developer Experience

- Comprehensive API documentation
- Clear error messages
- Consistent response formats
- Easy local development setup

### 4. Production Ready

- Graceful shutdown
- Health checks
- Prometheus metrics
- Structured logging
- Rate limiting

## Domain Model

### User Roles

```
┌─────────────────────────────────────┐
│              User (Abstract)        │
│  - id: uint                         │
│  - email: string                    │
│  - password_hash: string            │
│  - role: string                     │
│  - created_at: timestamp            │
│  - updated_at: timestamp            │
└─────────────────┬───────────────────┘
                  │
     ┌────────────┼────────────┐
     │            │            │
┌────┴────┐  ┌────┴───┐  ┌─────┴──┐
│  Admin  │  │Teacher │  │Student │
└─────────┘  └────────┘  └────────┘
```

### Core Entities

#### User

Base entity for authentication.

**Fields:**

- `id`: Primary key
- `email`: Unique email address
- `password_hash`: Bcrypt hash
- `role`: User role (admin, teacher, student)
- `is_active`: Account status
- `created_at`: Creation timestamp
- `updated_at`: Last update timestamp

#### Session

Tracks active user sessions.

**Fields:**

- `id`: Primary key
- `user_id`: Foreign key to User
- `token_hash`: Hashed JWT token
- `refresh_token_hash`: Hashed refresh token
- `expires_at`: Session expiration
- `ip_address`: Client IP
- `user_agent`: Client user agent
- `created_at`: Session start

#### Invitation

Manages user registration invitations.

**Fields:**

- `id`: Primary key
- `email`: Invitee email
- `role`: Role to assign
- `token`: Unique invitation token
- `expires_at`: Invitation expiration
- `created_by`: Admin who created invitation
- `used_at`: When invitation was used

#### Student

Extends User for student-specific data.

**Fields:**

- `id`: Primary key
- `user_id`: Foreign key to User
- `first_name`: First name
- `last_name`: Last name
- `student_id`: Unique student identifier
- Additional student-specific fields

#### Teacher

Extends User for teacher-specific data.

**Fields:**

- `id`: Primary key
- `user_id`: Foreign key to User
- `first_name`: First name
- `last_name`: Last name
- `department`: Teacher department
- Additional teacher-specific fields

## API Design

### RESTful Principles

FITS follows REST conventions:

- **Resources**: Nouns, not verbs (`/students`, not `/getStudents`)
- **HTTP Methods**: Standard semantics
  - `GET`: Retrieve resource(s)
  - `POST`: Create resource
  - `PUT`: Update entire resource
  - `PATCH`: Update partial resource
  - `DELETE`: Remove resource
- **Status Codes**: Appropriate HTTP status codes
- **HATEOAS**: Links to related resources (future enhancement)

### URL Structure

```
https://api.example.com/api/v1/{resource}/{id}/{sub-resource}
```

**Examples:**

- `GET /api/v1/student` - List all students
- `GET /api/v1/student/{id}` - Get specific student
- `POST /api/v1/student` - Create student
- `PUT /api/v1/student/{id}` - Update student
- `DELETE /api/v1/student/{id}` - Delete student

### Response Format

#### Success Response

```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "student@example.com",
    "role": "student"
  },
  "metadata": {
    "timestamp": "2025-10-23T18:00:00Z"
  }
}
```

#### Error Response

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
  },
  "metadata": {
    "timestamp": "2025-10-23T18:00:00Z"
  }
}
```

#### Paginated Response

```json
{
  "success": true,
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### Error Handling

Error codes are consistent across the API:

| Code                  | HTTP Status | Description                       |
| --------------------- | ----------- | --------------------------------- |
| `VALIDATION_ERROR`    | 400         | Invalid input data                |
| `UNAUTHORIZED`        | 401         | Missing or invalid authentication |
| `FORBIDDEN`           | 403         | Insufficient permissions          |
| `NOT_FOUND`           | 404         | Resource not found                |
| `CONFLICT`            | 409         | Resource conflict (duplicate)     |
| `RATE_LIMIT_EXCEEDED` | 429         | Too many requests                 |
| `INTERNAL_ERROR`      | 500         | Server error                      |

See `internal/common/errors/errors.go:1` for implementation.

## Security Design

### Authentication Flow

1. **Bootstrap**: Initial admin created via `/api/v1/bootstrap/init`
2. **Invitation**: Admin creates invitation via `/api/v1/admin/invite`
3. **Registration**: User registers via invitation token
4. **Login**: User authenticates with email/password
5. **Token**: Server issues JWT access token + refresh token
6. **Authorization**: Token included in `Authorization: Bearer {token}` header

### Password Security

- **Hashing**: Bcrypt with configurable cost (default: 12)
- **Validation**: Minimum 8 characters, complexity requirements
- **Storage**: Never stored in plain text
- **Transmission**: Only accepted over HTTPS in production

Implementation: `pkg/crypto/password.go:1`

### JWT Token Design

**Access Token:**

- Short-lived (default: 15 minutes)
- Contains user ID and role
- Signed with HS256
- Cannot be revoked individually

**Refresh Token:**

- Long-lived (default: 7 days)
- Stored in database with hash
- Can be revoked
- Used to obtain new access token

**Claims Structure:**

```json
{
  "sub": "123", // User ID
  "email": "user@example.com",
  "role": "student",
  "exp": 1698012345, // Expiration
  "iat": 1698012045 // Issued at
}
```

Implementation: `pkg/crypto/jwt.go:1`

### Session Management

Sessions are tracked in the database with:

- Token hash (not the actual token)
- IP address
- User agent
- Expiration timestamp

**Logout**: Deletes session from database
**Logout All**: Deletes all user sessions
**Session Cleanup**: Automatic cleanup of expired sessions

Implementation: `internal/domain/auth/auth_service.go:1`

## Rate Limiting Design

### Two-Tier Approach

#### Tier 1: Global IP-Based Rate Limiting

- Applied to ALL requests
- Prevents basic DoS attacks
- Default: 100 requests/minute per IP
- Implemented with Fiber's built-in limiter

#### Tier 2: Per-User Role-Based Rate Limiting

- Applied after authentication
- Different limits per role
- More granular control
- Automatic cleanup of expired entries

**Rate Limits by Role:**

```go
Admin:    1000 requests/minute
Teacher:   300 requests/minute
Student:   100 requests/minute
```

Implementation: `internal/middleware/ratelimit.go:1`

### Rate Limit Response

```json
{
  "success": false,
  "error": "Rate limit exceeded",
  "details": "Maximum 100 requests per minute allowed",
  "code": 429
}
```

## Database Design

### Schema

```sql
-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sessions table
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    refresh_token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Invitations table
CREATE TABLE invitations (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_by INTEGER REFERENCES users(id),
    used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Students table
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    student_id VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Teachers table
CREATE TABLE teachers (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    department VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Migrations

Migrations are handled automatically by GORM on startup.

Location: `pkg/database/migrations.go:1`

### Indexes

Critical indexes for performance:

- `users.email` - Unique index for login lookups
- `sessions.user_id` - Foreign key index
- `sessions.expires_at` - For cleanup queries
- `students.student_id` - Unique index
- `invitations.token` - Unique index for validation

## Configuration Management

Configuration is loaded from `configs/config.toml` with fallback to environment variables.

### Configuration Structure

```toml
[server]
host = "0.0.0.0"
port = 8080
read_timeout = "30s"
write_timeout = "30s"
rate_limit = 100
allowed_origins = "*"

[database]
host = "localhost"
port = 5432
user = "fits"
password = "password"
database = "fits_db"

[jwt]
secret = "your-secret-key"
access_token_expiry = "15m"
refresh_token_expiry = "168h"

[logging]
level = "info"
format = "json"

[tls]
enabled = false
cert_file = ""
key_file = ""
```

See `internal/config/config.go:1` for full structure.

## Logging Design

### Structured Logging

Uses Uber Zap for structured, high-performance logging.

**Log Levels:**

- `DEBUG`: Detailed debugging information
- `INFO`: General informational messages
- `WARN`: Warning messages
- `ERROR`: Error messages
- `FATAL`: Critical errors (exits application)

**Log Formats:**

- `json`: Machine-readable JSON (production)
- `console`: Human-readable text (development)

**Example Log Entry:**

```json
{
  "level": "info",
  "ts": "2025-10-23T18:00:00.000Z",
  "caller": "server/main.go:69",
  "msg": "FITS Backend starting up",
  "version": "1.0.1",
  "log_level": "info"
}
```

Implementation: `pkg/logger/logger.go:1`

## Testing Design

### Test Organization

```
domain/auth/
├── auth_service.go
├── auth_service_test.go      # Unit tests
├── repository.go
└── repository_gorm_test.go   # Integration tests
```

### Testing Strategy

1. **Unit Tests**: Test business logic in isolation
2. **Integration Tests**: Test database operations with SQLite
3. **Test Fixtures**: Shared test data and utilities
4. **Mocking**: Interface-based design enables easy mocking

### Example Test

```go
func TestAuthService_Login(t *testing.T) {
    // Arrange
    repo := NewMockRepository()
    service := NewAuthService(repo, jwtService, config)

    // Act
    token, err := service.Login(ctx, "test@example.com", "password")

    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
}
```

## Extensibility Points

### Adding New Domain

1. Create directory: `internal/domain/{entity}/`
2. Define models: `model.go`
3. Define repository interface: `repository.go`
4. Implement GORM repository: `repository_gorm.go`
5. Implement business logic: `service.go`
6. Implement HTTP handlers: `handler.go`
7. Register routes in `cmd/server/main.go`

### Adding Middleware

1. Create middleware: `internal/middleware/{name}.go`
2. Implement `fiber.Handler` interface
3. Apply in `cmd/server/main.go`

### Adding Configuration

1. Update struct: `internal/config/config.go`
2. Add to TOML: `configs/config.toml`
3. Document in `.env.example`

## Performance Characteristics

### Request Processing Time

- Average: < 10ms for simple queries
- P95: < 50ms
- P99: < 100ms

### Database Connections

- Pool size: 25 connections (default)
- Max idle: 5 connections
- Max lifetime: 5 minutes

### Memory Usage

- Baseline: ~50MB
- Under load: ~200MB
- Rate limiter: ~1KB per active user

### Concurrency

- Fiber handles multiple concurrent requests
- Database connection pooling
- Goroutine-safe rate limiter
