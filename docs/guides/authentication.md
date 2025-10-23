# Authentication System Guide

> Complete guide to JWT-based authentication, role-based access control, and invitation-based registration.

## Overview

FITS Backend implements a secure JWT-based authentication system with role-based access control (RBAC). The system supports three user roles with distinct permissions and includes an invitation-based registration flow for secure user onboarding.

## System Components

### 1. Admin Bootstrap System

The bootstrap system initializes the first admin account:

- **RSA Keypair Generation** - Automatic 4096-bit key generation on first run
- **Admin Token** - Permanent token (~100 years validity) for admin operations
- **One-Time Initialization** - Can only be called once, prevents duplicate admins
- **Secure Storage** - Keys stored in `configs/keys/` directory

**Usage:**

```bash
# Initialize admin (only works once)
curl -X POST http://localhost:8080/api/v1/bootstrap/init

# Returns admin token (save securely!)
{
  "admin_token": "eyJhbGc...",
  "public_key_path": "./configs/keys/admin.pub"
}
```

### 2. JWT Token System

Four token types with different purposes:

| Token Type | Lifetime | Purpose | Algorithm |
|------------|----------|---------|-----------|
| **Admin Token** | ~100 years | Bootstrap operations | RS256 (RSA) |
| **Access Token** | 1 hour | API authentication | HS256 (HMAC) |
| **Refresh Token** | 30 days | Obtain new access tokens | HS256 (HMAC) |
| **Invitation Token** | 7 days | User registration | HS256 (HMAC) |

**Configuration:**

```toml
[jwt]
secret = "your-jwt-secret-min-32-chars"
access_token_expiry = "1h"
refresh_token_expiry = "720h"  # 30 days
invitation_expiry = "168h"      # 7 days
admin_key_path = "./configs/keys/admin.key"
admin_pub_key_path = "./configs/keys/admin.pub"
```

### 3. Invitation System

Secure registration flow managed by administrators:

**Admin creates invitation:**

```bash
POST /api/v1/admin/invite
Authorization: Bearer <admin_token>

{
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "role": "teacher"  # or "student"
}

# Returns invitation token
{
  "invitation_token": "eyJhbGc...",
  "expires_at": "2025-10-29T12:00:00Z"
}
```

**User completes registration:**

```bash
POST /api/v1/invite/{token}/complete

{
  "username": "john.doe",
  "password": "SecurePassword123!"
}

# Returns access and refresh tokens
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "user": { ... }
}
```

**Security Features:**
- Single-use tokens (invalidated after use)
- Time-limited (7 days by default)
- Bcrypt password hashing (cost factor 12)
- Username uniqueness validation

### 4. Login & Session Management

**Login Flow:**

```bash
POST /api/v1/auth/login

{
  "username": "john.doe",
  "password": "SecurePassword123!"
}

# Returns tokens
{
  "access_token": "eyJhbGc...",   # 1 hour lifetime
  "refresh_token": "eyJhbGc...",  # 30 days lifetime
  "user": {
    "id": "uuid",
    "username": "john.doe",
    "role": "teacher"
  }
}
```

**Token Refresh:**

```bash
POST /api/v1/auth/refresh

{
  "refresh_token": "eyJhbGc..."
}

# Returns new access token
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc..."  # New refresh token
}
```

**Logout:**

```bash
POST /api/v1/auth/logout
Authorization: Bearer <access_token>

# Invalidates all refresh tokens for the user
{
  "message": "Logged out successfully"
}
```

### 5. Middleware & Authorization

**JWT Middleware:**

Automatically validates tokens on protected routes:

```go
// Require authentication
app.Post("/api/v1/resource",
    jwtMiddleware.RequireAuth(),
    handler.Create,
)

// Require admin role
app.Post("/api/v1/admin/action",
    jwtMiddleware.RequireAuth(),
    middleware.RequireAdmin(),
    handler.AdminAction,
)
```

**RBAC (Role-Based Access Control):**

Three roles with hierarchical permissions:

```
Admin (Superuser)
  ├─ Full access to all resources
  ├─ Create and manage invitations
  ├─ Manage all users (students, teachers)
  ├─ System configuration
  └─ View all data

Teacher
  ├─ Sign reports for assigned students
  ├─ Manage assigned students' profiles
  ├─ View reports for assigned students
  └─ Update own profile

Student
  ├─ Upload training reports
  ├─ View own reports and signatures
  ├─ Update own profile
  └─ View assigned teacher
```

**Ownership Checks:**

Resources are protected by ownership validation:

```go
// Students can only access their own data
if user.Role == "student" && user.UserUUID != resourceOwner {
    return fiber.ErrForbidden
}

// Teachers can only access assigned students
if user.Role == "teacher" && !isAssigned(user, student) {
    return fiber.ErrForbidden
}
```

### 6. Database Integration

**PostgreSQL with GORM:**

- Connection pooling (25 max, 5 min connections)
- Automatic schema migration on startup
- Health checks for database connectivity
- Prepared statements (SQL injection prevention)

**Tables:**

- `users` - User accounts and credentials
- `refresh_tokens` - Active refresh tokens
- `invitations` - Pending invitations
- `students` - Student profiles
- `teachers` - Teacher profiles
- `teacher_keys` - RSA keys for signatures

## Quick Start

### 1. Start PostgreSQL

```bash
make db-up
```

### 2. Start Server

```bash
make run
```

Server starts and automatically:
- Creates database if not exists
- Runs all migrations
- Generates admin keypair
- Starts on http://localhost:8080

### 3. Bootstrap Admin

```bash
curl -X POST http://localhost:8080/api/v1/bootstrap/init

# Save the admin_token!
```

### 4. Create Teacher Invitation

```bash
curl -X POST http://localhost:8080/api/v1/admin/invite \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "teacher@school.com",
    "first_name": "Jane",
    "last_name": "Smith",
    "role": "teacher",
    "department": "Computer Science"
  }'

# Get invitation_token
```

### 5. Complete Registration

```bash
curl -X POST http://localhost:8080/api/v1/invite/{invitation_token}/complete \
  -H "Content-Type: application/json" \
  -d '{
    "username": "jane.smith",
    "password": "SecurePassword123!"
  }'

# Get access_token and refresh_token
```

### 6. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "jane.smith",
    "password": "SecurePassword123!"
  }'
```

## Authentication Flow Diagram

```
1. Admin Bootstrap
   └─> POST /api/v1/bootstrap/init
       └─> Returns admin_token (permanent)

2. Create Invitation
   └─> POST /api/v1/admin/invite
       └─> Returns invitation_token (7 days)

3. User Registration
   └─> POST /api/v1/invite/{token}/complete
       └─> Creates user account
       └─> Returns access_token + refresh_token

4. User Login
   └─> POST /api/v1/auth/login
       └─> Returns access_token (1h) + refresh_token (30d)

5. API Access
   └─> Request with Authorization: Bearer <access_token>
       └─> Middleware validates token
       └─> RBAC checks permissions
       └─> Ownership validated
       └─> Access granted/denied

6. Token Refresh
   └─> POST /api/v1/auth/refresh
       └─> Returns new access_token + refresh_token

7. Logout
   └─> POST /api/v1/auth/logout
       └─> Invalidates all refresh tokens
```

## Security Features

### Password Security

- **Bcrypt Hashing** - Industry-standard password hashing
- **Cost Factor 12** - Balanced security vs. performance
- **Password Validation** - Minimum requirements enforced
- **No Plain Text** - Passwords never stored unencrypted

### Token Security

- **HS256 Algorithm** - HMAC with SHA-256 for user tokens
- **RS256 Algorithm** - RSA with SHA-256 for admin tokens
- **Token Expiration** - All tokens have limited lifetime
- **Refresh Rotation** - New refresh token on each refresh
- **Logout Invalidation** - Tokens can be revoked

### API Security

- **RBAC** - Role-based access control on all endpoints
- **Ownership Checks** - Users can only access their data
- **Rate Limiting** - Protection against brute force
- **HTTPS Ready** - TLS support for production
- **Security Headers** - CORS, CSP, and best practices

### Database Security

- **Prepared Statements** - SQL injection prevention
- **Connection Pooling** - Limited connections per user
- **Password Protection** - Database credentials required
- **SSL Support** - Encrypted database connections

## Configuration

### JWT Settings

```toml
[jwt]
# Secret key (minimum 32 characters)
secret = "CHANGE-THIS-SECRET-KEY-IN-PRODUCTION"

# Token lifetimes
access_token_expiry = "1h"
refresh_token_expiry = "720h"  # 30 days
invitation_expiry = "168h"      # 7 days

# RSA key paths (auto-generated)
admin_key_path = "./configs/keys/admin.key"
admin_pub_key_path = "./configs/keys/admin.pub"
```

### Database Settings

```toml
[database]
host = "localhost"
port = 5432
user = "fits_user"
password = "fits_password"  # Change in production
database = "fits_db"
ssl_mode = "disable"        # Use "require" in production
max_conns = 25
min_conns = 5
```

## Testing

### Automated Tests

```bash
# All tests
make test

# Auth domain tests
go test -v ./internal/domain/auth/...

# Middleware tests
go test -v ./internal/middleware/...

# Crypto tests
go test -v ./pkg/crypto/...
```

### Manual Testing

Use the Swagger UI for interactive testing:

```
http://localhost:8080/docs
```

Or use the test scripts:

```bash
# Full authentication flow test
./scripts/test_auth_flow.sh

# Complete E2E test
./scripts/test_full_flow.sh
```

## Troubleshooting

### Common Issues

**Problem: Bootstrap fails**
```
Error: "Admin already initialized"
```
Solution: Bootstrap can only run once. Reset database if needed:
```bash
make db-reset
```

**Problem: Token expired**
```
Error: "Token is expired"
```
Solution: Use refresh token to get new access token:
```bash
POST /api/v1/auth/refresh
```

**Problem: Invalid credentials**
```
Error: "Invalid username or password"
```
Solution: Check username (case-sensitive) and password. Ensure user is registered.

**Problem: Forbidden access**
```
Error: "Forbidden"
```
Solution: Check user role and permissions. Ensure proper RBAC middleware is applied.

## Related Documentation

- [Swagger UI Guide](swagger-ui.md) - Interactive API documentation
- [Security Guide](security.md) - Security best practices
- [Makefile Guide](makefile.md) - Development workflow
- [Main README](../../README.md) - Project overview and quick start

## License

This project is licensed under the GNU General Public License v3.0 - see [LICENSE](../../LICENSE) for details.

---

**Status:** ✅ Production Ready | **Last Updated:** 2025-10-22
