# FITS Backend

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-GPL--3.0-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/Version-1.0.0-green.svg)]()
[![API Docs](https://img.shields.io/badge/API-Swagger-85EA2D.svg)](http://localhost:8080/docs)

##  Description

FITS (Flexible IT Training System) is a modern backend system for managing training reports and digital signatures in educational environments. Built with Go and Fiber, it provides a robust REST API with JWT authentication, role-based access control, and comprehensive student/teacher management capabilities.

##  Features

- **JWT Authentication** - Secure token-based authentication with access and refresh tokens
- **Role-Based Access Control (RBAC)** - Three-tier permission system (Admin, Teacher, Student)
- **User Management** - Complete student and teacher profile management
- **Invitation System** - Token-based secure registration flow
- **Digital Signatures** - RSA 4096-bit cryptographic signatures for reports (experimental)
- **Pagination Support** - Efficient data retrieval with configurable page sizes
- **Rate Limiting** - Protection against API abuse
- **Security Headers** - CORS, CSRF, and security best practices
- **Structured Logging** - Production-ready logging with Zap
- **Auto-Migration** - Automatic database schema initialization and versioning
- **API Documentation** - Interactive Swagger UI for API exploration
- **Health Checks** - Built-in health and metrics endpoints
- **Docker Support** - Quick setup with Docker Compose

##  Tech Stack

- **Language:** Go 1.21+
- **Web Framework:** [Fiber](https://gofiber.io/) - Express-inspired web framework
- **Database:** PostgreSQL 15+
- **ORM:** [GORM](https://gorm.io/) - Feature-rich ORM with auto-migration
- **Authentication:** JWT (golang-jwt) with HS256
- **Password Hashing:** Bcrypt (cost factor 12)
- **Cryptography:** RSA-PSS 4096-bit with SHA-256
- **Logging:** [Zap](https://github.com/uber-go/zap) - Structured, leveled logging
- **Testing:** Testify with 70%+ coverage for core components
- **Documentation:** Swagger/OpenAPI 3.0

##  Quick Start

### Prerequisites

```bash
# Required
- Go 1.21 or higher
- PostgreSQL 15+
- make
- jq (for E2E tests)

# Optional
- Docker & Docker Compose (recommended for quick setup)
```

### Installation (3 Commands!)

```bash
# 1. Start PostgreSQL
make docker-up

# 2. In a new terminal: Start the server
make run

# 3. In another terminal: Initialize admin user
make bootstrap
```

**That's it!** The server is running and the admin user is initialized.

### Available Make Commands

```bash
make help              # Show all available commands
make run               # Start the server
make test              # Run unit tests
make test-coverage     # Run tests with coverage report
make e2e-test          # Run end-to-end tests
make build             # Build binary
make clean             # Clean build artifacts
make reset-db          # Reset database
make docker-up         # Start PostgreSQL in Docker
make docker-down       # Stop PostgreSQL
```

### Manual Installation

```bash
# 1. Clone the repository
git clone <repo-url>
cd FITS-backend

# 2. Install dependencies
go mod download

# 3. Start PostgreSQL (Docker Compose)
docker-compose up -d

# 4. Wait for PostgreSQL to be ready (~5 seconds)
sleep 5

# 5. Configure environment (optional)
cp .env.example .env
# Edit .env if needed

# 6. Start the server
go run cmd/server/main.go
```

The database will be created automatically and all migrations will run on startup.

### Access Points

- **API:** http://localhost:8080
- **Swagger UI:** http://localhost:8080/swagger/index.html
- **Health Check:** http://localhost:8080/health
- **Metrics:** http://localhost:8080/metrics?secret=i-forgot-my-token=

##  API Documentation

Full interactive API documentation is available via Swagger UI:

**http://localhost:8080/swagger/index.html**

### Main Endpoints

**Authentication:**
```
POST   /api/v1/bootstrap/init              # Initialize admin user
POST   /api/v1/auth/login                  # User login
POST   /api/v1/auth/logout                 # User logout (auth required)
POST   /api/v1/auth/refresh                # Refresh access token
```

**Invitations:**
```
POST   /api/v1/admin/invite                # Create invitation (Admin only)
GET    /api/v1/invite/:token               # Get invitation details
POST   /api/v1/invite/:token/complete      # Complete registration
```

**Students:**
```
POST   /api/v1/student                     # Create student (Admin)
GET    /api/v1/student/:uuid               # Get student
PUT    /api/v1/student/:uuid               # Update student (Admin)
DELETE /api/v1/student/:uuid               # Delete student (Admin)
GET    /api/v1/student                     # List all students (paginated)
```

**Teachers:**
```
POST   /api/v1/teacher                     # Create teacher (Admin)
GET    /api/v1/teacher/:uuid               # Get teacher
PUT    /api/v1/teacher/:uuid               # Update teacher (Admin)
DELETE /api/v1/teacher/:uuid               # Delete teacher (Admin)
GET    /api/v1/teacher                     # List all teachers (paginated)
```

**System:**
```
GET    /health                             # Health check
GET    /metrics?secret=xxx                 # Prometheus metrics
GET    /swagger/*                          # API documentation
```

##  Project Structure

```
FITS-backend/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── configs/
│   ├── config.toml                 # Configuration file
│   └── keys/                       # RSA keys (auto-generated)
├── internal/
│   ├── common/                     # Shared code
│   │   ├── errors/                 # Custom error types
│   │   └── response/               # API response wrapper
│   ├── config/                     # Configuration loader
│   ├── domain/                     # Domain logic
│   │   ├── auth/                   # Authentication domain
│   │   ├── signing/                # Digital signatures (v1.1)
│   │   ├── student/                # Student management
│   │   └── teacher/                # Teacher management
│   └── middleware/                 # HTTP middleware
│       ├── jwt.go                  # JWT authentication
│       └── rbac.go                 # Role-based access control
├── pkg/                            # Reusable packages
│   ├── crypto/                     # Cryptography utilities
│   │   ├── jwt.go                  # JWT implementation
│   │   ├── password.go             # Bcrypt password hashing
│   │   └── rsa.go                  # RSA signatures
│   └── database/                   # Database layer
│       ├── database.go             # Connection & auto-setup
│       └── migrations.go           # Migration system
├── migrations/                     # SQL migrations (deprecated)
├── docs/                           # Swagger documentation
├── docker-compose.yml              # Docker Compose setup
└── Makefile                        # Build automation
```

##  Configuration

### Configuration File: `configs/config.toml`

```toml
[server]
port = 8080
host = "0.0.0.0"
read_timeout = "30s"
write_timeout = "30s"

[database]
host = "localhost"
port = 5432
user = "fits_user"
password = "fits_password"  #  CHANGE in production!
database = "fits_db"
ssl_mode = "disable"        #  Use "require" in production!
max_conns = 25
min_conns = 5

[jwt]
secret = "your-jwt-secret-change-in-production"  #  CHANGE!
access_token_expiry = "1h"
refresh_token_expiry = "720h"    # 30 days
invitation_expiry = "168h"       # 7 days
admin_key_path = "./configs/keys/admin.key"
admin_pub_key_path = "./configs/keys/admin.pub"

[storage]
upload_dir = "./uploads"
max_file_size = 104857600  # 100 MB

[secrets]
metrics_secret = "i-forgot-my-token="      #  CHANGE!
registration_secret = "i-forgot-my-token=" #  CHANGE!
deletion_secret = "i-forgot-my-token="     #  CHANGE!
update_secret = "i-forgot-my-token="       #  CHANGE!
```

### Environment Variables

You can override configuration values using environment variables:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=fits_user
export DB_PASSWORD=secure_password
export JWT_SECRET=your-super-secret-key
```

### Production Security Checklist

Before deploying to production:

- [ ] Change JWT secret (minimum 32 characters)
- [ ] Change database password
- [ ] Enable SSL/TLS for database (`ssl_mode = "require"`)
- [ ] Update all secrets in `configs/config.toml`
- [ ] Enable rate limiting
- [ ] Enforce HTTPS
- [ ] Configure CORS properly
- [ ] Set up audit logging
- [ ] Review and restrict file upload sizes
- [ ] Enable monitoring and alerting

##  Development

### Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details on:
- Code style guidelines
- Development workflow
- Pull request process
- Testing requirements

### Building

```bash
# Build binary
go build -o bin/fits-server cmd/server/main.go

# Build with optimizations
go build -ldflags="-s -w" -o bin/fits-server cmd/server/main.go

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o bin/fits-server-linux cmd/server/main.go
GOOS=windows GOARCH=amd64 go build -o bin/fits-server.exe cmd/server/main.go
```

### Updating Dependencies

```bash
go get -u ./...
go mod tidy
```

### Regenerate Swagger Documentation

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init -g cmd/server/main.go
```

### Code Quality

```bash
# Linting
go vet ./...
golangci-lint run

# Formatting
go fmt ./...
gofmt -s -w .

# Security check
gosec ./...
```

##  Testing

### Run Tests

```bash
# All tests
make test

# With coverage
make test-coverage

# View coverage in browser
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Specific package tests
go test -v ./pkg/crypto/...
go test -v ./internal/middleware/...
go test -v ./internal/domain/auth/...

# Run benchmarks
go test -bench=. -benchmem ./pkg/crypto/
```

### Test Coverage

| Component | Coverage | Status |
|-----------|----------|--------|
| pkg/crypto | 73.7% |  Excellent |
| internal/middleware | 81.2% |  Excellent |
| internal/domain/auth | 18.1% |  Basic coverage |
| **Overall** | **20.3%** |  **Growing** |

**Note:** Core security components (crypto, middleware) have >70% coverage.

### End-to-End Tests

```bash
# Run full E2E test suite
make e2e-test
```

##  Deployment

### Docker Deployment (Recommended)

```bash
# Build and run with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Systemd Service

Create `/etc/systemd/system/fits-backend.service`:

```ini
[Unit]
Description=FITS Backend Server
After=network.target postgresql.service

[Service]
Type=simple
User=fits
WorkingDirectory=/opt/fits-backend
ExecStart=/opt/fits-backend/fits-server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable fits-backend
sudo systemctl start fits-backend
sudo systemctl status fits-backend
```

### Deployment Documentation

For detailed deployment instructions, see:
- [Deployment Guide](./docs/guides/deployment.md) - Docker and production setup

##  Monitoring

### Health Check

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "ok",
  "database": "connected",
  "time": "2025-10-19T12:00:00Z"
}
```

### Metrics (Prometheus Format)

```bash
curl "http://localhost:8080/metrics?secret=your-metrics-secret"
```

The metrics endpoint provides Prometheus-compatible metrics for monitoring.

##  Security

### Authentication Flow

1. **Login:** User provides credentials → receives access token (1h) and refresh token (30 days)
2. **API Access:** Client sends access token in Authorization header
3. **Token Refresh:** When access token expires, use refresh token to get new access token
4. **Logout:** Refresh token is invalidated

### Token Types

- **Access Token** - Short-lived (1 hour), used for API requests
- **Refresh Token** - Long-lived (30 days), used to obtain new access tokens
- **Invitation Token** - Single-use (7 days), for user registration
- **Admin Token** - Permanent (~100 years), for bootstrap process

### Role-Based Access Control

```
Admin (Superuser)
  ├─ Full access to all resources
  ├─ Can create invitations
  └─ Can manage all users

Teacher
  ├─ Can sign reports
  ├─ Can manage assigned students
  └─ Can view reports

Student
  ├─ Can upload reports
  ├─ Can view own reports
  └─ Can manage own profile
```

### Security Features

- **Password Security:** Bcrypt hashing with cost factor 12
- **JWT Security:** HS256 signing with token expiration
- **RSA Security:** 4096-bit keys with RSA-PSS and SHA-256
- **API Security:** RBAC on all endpoints with resource ownership checks
- **Database Security:** Prepared statements to prevent SQL injection

## Additional Documentation

- [Documentation Index](./docs/README.md) - Complete documentation navigation
- [Implementation Summary](./docs/development/IMPLEMENTATION_SUMMARY.md) - Feature implementation status
- [Known Issues](./docs/development/KNOWN_ISSUES.md) - Current limitations and bugs
- [Quick Reference](./docs/guides/QUICK_REFERENCE.md) - Quick command reference

##  Roadmap

###  Completed (v1.0)

- [x] JWT authentication system
- [x] RBAC middleware
- [x] Automatic database migration
- [x] Invitation system
- [x] Student/Teacher management with pagination
- [x] Unit tests (>70% core coverage)
- [x] Swagger API documentation
- [x] Health check & metrics
- [x] Security headers and rate limiting
- [x] Structured logging with Zap

###  In Progress (v1.1)

- [ ] Digital signatures (RSA) - **Experimental**
- [ ] Report upload functionality
- [ ] Teacher keypair generation
- [ ] Email service for invitations
- [ ] Enhanced audit logging

###  Planned (v2.0)

- [ ] Frontend integration
- [ ] PDF report generation
- [ ] Batch operations
- [ ] Advanced search & filtering
- [ ] Backup & restore system
- [ ] Multi-tenancy support

##  Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for detailed guidelines.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

##  Acknowledgments

- [Fiber](https://gofiber.io/) - Express-inspired web framework for Go
- [GORM](https://gorm.io/) - The fantastic ORM library for Golang
- [Testify](https://github.com/stretchr/testify) - Testing toolkit
- [golang-jwt](https://github.com/golang-jwt/jwt) - JWT implementation
- [Zap](https://github.com/uber-go/zap) - Blazing fast structured logging

---

**Version:** 1.0.0
**Status:** Production Ready 
**Last Updated:** 2025-10-19
**Maintainer:** FITS Development Team

For questions or support, please open an issue on GitHub.
