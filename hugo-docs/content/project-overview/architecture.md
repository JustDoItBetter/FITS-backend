---
title: "Architecture"
weight: 1
description: "FITS Backend system architecture using Clean Architecture pattern. Layer-by-layer breakdown of HTTP, middleware, domain, and infrastructure layers."
---

# System Architecture

## Overview

FITS Backend follows a **Clean Architecture** pattern with clear separation of concerns across multiple layers. The architecture promotes maintainability, testability, and scalability.

## Architecture Layers

```
┌─────────────────────────────────────────────────────┐
│                   HTTP Layer                        │
│            (Fiber Web Framework)                    │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────┴──────────────────────────────┐
│               Middleware Layer                      │
│  (Auth, RBAC, Rate Limiting, Logging, Security)     │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────┴──────────────────────────────┐
│                Domain Layer                         │
│         ┌────────────┬──────────┬────────┐          │
│         │   Auth     │ Student  │Teacher │          │
│         │  Domain    │ Domain   │ Domain │          │
│         └─────┬──────┴────┬─────┴────┬───┘          │
│               │           │          │              │
│         ┌─────┴───────────┴──────────┴──┐           │
│         │      Repository Interface     │           │
│         └─────────────────┬─────────────┘           │
└───────────────────────────┴─────────────────────────┘
                            │
┌───────────────────────────┴─────────────────────────┐
│           Infrastructure Layer                      │
│                                                     │
│  ┌──────────────────┐      ┌──────────────────┐     │
│  │  Database (GORM) │      │  Crypto/JWT      │     │
│  │  PostgreSQL      │      │  Services        │     │
│  └──────────────────┘      └──────────────────┘     │
└─────────────────────────────────────────────────────┘
```

## Core Components

### 1. HTTP Layer (`cmd/server/main.go`)

The entry point for the application using the Fiber web framework.

**Responsibilities:**
- Server initialization and configuration
- Route registration
- Middleware pipeline setup
- Graceful shutdown handling
- TLS/HTTPS support

**Key Features:**
- Fast HTTP server built on valyala/fasthttp
- Automatic request/response handling
- Built-in CORS support
- Configurable timeouts and limits

### 2. Middleware Layer (`internal/middleware/`)

Security and cross-cutting concerns implemented as middleware.

**Components:**

#### Authentication (`jwt.go`)
- JWT token validation
- Token extraction from Authorization header
- User context injection

#### Authorization (`rbac.go`)
- Role-based access control
- Permission checking
- Admin-only endpoint protection

#### Rate Limiting (`ratelimit.go`)
- Two-tier rate limiting:
  1. **Global**: IP-based rate limiting for all requests
  2. **Per-User**: Role-based limits (admin, teacher, student)
- Automatic cleanup of expired entries
- Configurable limits per role

#### Security Headers
- XSS Protection
- Content-Type Options
- Frame Options
- CSP (Content Security Policy)
- Referrer Policy

### 3. Domain Layer (`internal/domain/`)

Business logic organized by domain entities following Domain-Driven Design (DDD).

#### Structure per Domain

```
domain/{entity}/
├── model.go          # Domain entities and DTOs
├── repository.go     # Repository interface
├── repository_gorm.go # GORM implementation
├── service.go        # Business logic
├── handler.go        # HTTP handlers
└── *_test.go        # Unit tests
```

#### Domains

##### Auth Domain (`internal/domain/auth/`)
- User authentication (login/logout)
- Session management
- Token refresh
- Invitation system
- Bootstrap admin creation

**Key Files:**
- `auth/auth_service.go:1` - Core authentication logic
- `auth/invitation_service.go:1` - Invitation creation and validation
- `auth/bootstrap_service.go:1` - Initial admin setup

##### Student Domain (`internal/domain/student/`)
- Student CRUD operations
- Student data validation
- Student-specific business rules

##### Teacher Domain (`internal/domain/teacher/`)
- Teacher CRUD operations
- Teacher data validation
- Teacher-specific business rules

##### Signing Domain (`internal/domain/signing/`)
- Digital signing workflow
- Signing request management

### 4. Common Utilities (`internal/common/`)

Shared functionality across domains.

#### Error Handling (`errors/errors.go`)
- Centralized error definitions
- HTTP status code mapping
- Consistent error responses

#### Response Formatting (`response/response.go`)
- Standard API response structure
- Success/error response helpers

#### Pagination (`pagination/pagination.go`)
- Query parameter parsing
- Offset/limit calculation
- Metadata generation

#### Validation (`validation/`)
- Password strength validation
- Input sanitization
- Custom validators

### 5. Infrastructure Layer (`pkg/`)

Low-level technical concerns and external dependencies.

#### Database (`pkg/database/`)
- GORM connection management
- Migration handling
- Transaction support
- Health checks

**Key Features:**
- PostgreSQL for production
- SQLite for testing
- Automatic migrations
- Connection pooling

#### Crypto (`pkg/crypto/`)
- Password hashing (bcrypt)
- JWT token generation/validation
- RSA key management

**Key Files:**
- `crypto/jwt.go:1` - JWT operations
- `crypto/password.go:1` - Password hashing
- `crypto/rsa.go:1` - RSA encryption

#### Logger (`pkg/logger/`)
- Structured logging with Uber Zap
- Configurable log levels
- JSON and console output formats

## Design Principles

### 1. Dependency Injection
All dependencies are injected through constructors, making testing and mocking straightforward.

```go
// Example from cmd/server/main.go
authRepo := auth.NewGormRepository(db.DB)
authService := auth.NewAuthService(authRepo, jwtService, &cfg.JWT)
authHandler := auth.NewHandler(bootstrapService, invitationService, authService)
```

### 2. Interface-Based Design
Repository interfaces allow swapping implementations without changing business logic.

```go
type Repository interface {
    Create(ctx context.Context, student *Student) error
    FindByID(ctx context.Context, id uint) (*Student, error)
    // ...
}
```

### 3. Clean Separation of Concerns
- **Handlers**: HTTP request/response handling
- **Services**: Business logic
- **Repositories**: Data persistence
- **Models**: Data structures

### 4. Configuration Management
All configuration is externalized in `configs/config.toml` and environment variables.

## Security Architecture

### Authentication Flow

```
┌─────────┐                ┌─────────┐                ┌──────────┐
│ Client  │                │  API    │                │ Database │
└────┬────┘                └────┬────┘                └─────┬────┘
     │                          │                           │
     │  POST /auth/login        │                           │
     ├─────────────────────────>│                           │
     │  {email, password}       │                           │
     │                          │  Validate credentials     │
     │                          ├──────────────────────────>│
     │                          │                           │
     │                          │<──────────────────────────┤
     │                          │  User data                │
     │                          │                           │
     │                          │  Generate JWT             │
     │                          │  Create session           │
     │                          ├──────────────────────────>│
     │                          │                           │
     │  200 OK                  │                           │
     │<─────────────────────────┤                           │
     │  {token, refresh_token}  │                           │
     │                          │                           │
     │  GET /api/v1/student     │                           │
     │  Authorization: Bearer   │                           │
     ├─────────────────────────>│                           │
     │                          │  Validate JWT             │
     │                          │  Check permissions        │
     │                          │  Apply rate limit         │
     │                          │                           │
     │                          │  Query data               │
     │                          ├──────────────────────────>│
     │                          │                           │
     │  200 OK                  │                           │
     │<─────────────────────────┤                           │
     │  {data}                  │                           │
```

### Authorization Layers

1. **Authentication Middleware** (`middleware/jwt.go:1`): Validates JWT token
2. **RBAC Middleware** (`middleware/rbac.go:1`): Checks user permissions
3. **Rate Limiting**: Prevents abuse

### Role Hierarchy

- **Admin**: Full system access, can invite users
- **Teacher**: Manage students, view reports
- **Student**: View own data, submit requests

## Request Flow

1. **Request arrives** → Global rate limiter checks IP
2. **Security headers** applied
3. **CORS validation**
4. **Route matching**
5. **JWT authentication** (if required)
6. **RBAC authorization** (if required)
7. **Per-user rate limiting**
8. **Handler execution**
9. **Response formatting**
10. **Logging**

## Performance Considerations

### Rate Limiting Strategy

**Global Rate Limit** (IP-based):
- Prevents basic DoS attacks
- Protects unauthenticated endpoints
- Default: 100 requests/minute per IP

**Per-User Rate Limit** (Role-based):
- Prevents authenticated user abuse
- Different limits per role:
  - Admin: 1000 req/min
  - Teacher: 300 req/min
  - Student: 100 req/min

### Database Optimization

- Connection pooling via GORM
- Prepared statements
- Indexed queries
- Lazy loading

### Caching Strategy

Currently no caching layer; future considerations:
- Redis for session storage
- In-memory cache for frequently accessed data
- CDN for static assets

## Testing Strategy

- **Unit Tests**: Domain logic and utilities
- **Integration Tests**: Database operations
- **Test Coverage**: Focus on critical paths
- **Mocking**: Interface-based design enables easy mocking

Example test locations:
- `internal/domain/auth/auth_service_test.go:1`
- `internal/middleware/jwt_test.go:1`
- `pkg/crypto/password_test.go:1`

## Deployment Architecture

```
┌─────────────────────────────────────┐
│                Load Balancer        │
└──────────────────┬──────────────────┘
                   │
       ┌───────────┴───────────┐
       │                       │
┌──────┴──────┐        ┌───────┴─────┐
│  FITS API   │        │  FITS API   │
│  Instance 1 │        │  Instance 2 │
└──────┬──────┘        └──────┬──────┘
       │                       │
       └───────────┬───────────┘
                   │
         ┌─────────┴─────────┐
         │   PostgreSQL      │
         │   (Primary)       │
         └───────────────────┘
```

## Configuration

Configuration is loaded from `configs/config.toml` with support for:
- Database connection
- JWT secrets
- Server settings
- TLS certificates
- Rate limits
- CORS origins
- Logging levels

See `internal/config/config.go:1` for full configuration structure.

## Monitoring & Observability

- **Structured Logging**: Uber Zap with JSON output
- **Metrics**: Prometheus endpoint at `/metrics`
- **Health Checks**: `/health` endpoint
- **Request Logging**: HTTP request/response logging

## Future Enhancements

1. **Caching Layer**: Redis integration
2. **Message Queue**: Async task processing
3. **Distributed Tracing**: OpenTelemetry
4. **API Versioning**: Multiple API versions
5. **WebSockets**: Real-time notifications
6. **File Storage**: S3-compatible storage
