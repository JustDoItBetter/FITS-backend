# FITS Backend Documentation

> Central index for all FITS Backend documentation, guides, and API references.

## Quick Links

- [Main README](../README.md) - Project overview and quick start
- [Swagger UI](http://localhost:8080/docs) - Interactive API documentation (when server is running)
- [Health Check](http://localhost:8080/health) - Server status endpoint

## Getting Started

### Essential Guides

- [Makefile Guide](guides/makefile.md) - Complete development workflow reference
- [Authentication System](guides/authentication.md) - JWT, RBAC, and invitation-based registration
- [Swagger UI Guide](guides/swagger-ui.md) - Using the interactive API documentation

### Quick Start Commands

```bash
# First time setup
make setup        # Install dependencies, tools, and start PostgreSQL
make run          # Start the server
make db-seed      # Populate with test data

# Daily development
make dev          # Start database and server together
make test         # Run tests
make pre-commit   # Run quality checks before committing
```

## Documentation Structure

### API Documentation

**OpenAPI/Swagger Files:**
- `api/docs.go` - Go-generated Swagger documentation
- `api/swagger.json` - OpenAPI 3.0 specification (JSON)
- `api/swagger.yaml` - OpenAPI 3.0 specification (YAML)

**Guides:**
- [Swagger UI Guide](guides/swagger-ui.md) - How to use the interactive API interface
- [API Quick Start](API_QUICK_START.md) - Quick API usage examples

### Development Guides

**Core Workflows:**
- [Makefile Guide](guides/makefile.md) - Complete command reference
- [Authentication Guide](guides/authentication.md) - JWT, sessions, and RBAC

**Best Practices:**
- Testing workflow (see [Makefile Guide](guides/makefile.md#testing))
- Pre-commit checks (see [Makefile Guide](guides/makefile.md#workflow-commands))
- Database management (see [Makefile Guide](guides/makefile.md#database-management))

### Architecture & Design

**Current Documentation:**
- [Authentication System](guides/authentication.md) - Complete auth system design
- Database schema (auto-migrated, see `internal/database/migrations/`)

### Historical Reports (Archive)

Archived development reports and fix documentation:

**2025-10 Reports:**
- [CORS Fix](archive/2025-10-cors-fix.md) - Swagger UI CORS configuration resolution
- [Swagger Fix](archive/2025-10-swagger-fix.md) - Swagger UI browser issues resolution
- [Security Fixes](archive/2025-10-security-fixes.md) - Security improvements
- [Makefile Modernization](archive/2025-10-makefile-modernization.md) - Build system improvements
- [Makefile Before/After](archive/2025-10-makefile-before-after.md) - Build system comparison
- [Implementation Complete](archive/2025-10-implementation.md) - Feature implementation summary
- [Phase 2 Improvements](archive/2025-10-phase2.md) - Secondary enhancements
- [Diagnosis Report](archive/2025-10-diagnosis.md) - System diagnosis and fixes
- [Final Summary](archive/2025-10-final-summary.md) - Project completion summary

## Common Tasks

### Development Workflow

```bash
# Start development environment
make dev          # Starts PostgreSQL and server

# Run tests
make test         # Unit tests
make test-coverage # With HTML coverage report

# Code quality
make fmt          # Format code
make lint         # Run linters
make check        # Complete quality check
```

### Database Management

```bash
# Status and control
make db-status    # Check database status
make db-up        # Start PostgreSQL container
make db-logs      # View database logs

# Data management
make db-reset     # Clear all data (keeps schema)
make db-seed      # Populate test data
make db-destroy   # Complete database removal
```

### Build and Deploy

```bash
# Development build
make build        # Build with debug symbols

# Production build
make build-prod   # Optimized binary for deployment

# Clean build
make clean        # Remove all build artifacts
make build-prod   # Fresh production build
```

### Documentation

```bash
# Regenerate API docs
make swagger      # Update Swagger/OpenAPI files

# View documentation
make docs-serve   # Open Swagger UI in browser
```

## API Endpoints

### Authentication

- `POST /api/v1/bootstrap/init` - Initialize admin (one-time only)
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `POST /api/v1/auth/refresh` - Refresh access token

### Invitations

- `POST /api/v1/admin/invite` - Create user invitation (admin only)
- `GET /api/v1/invite/{token}` - Get invitation details
- `POST /api/v1/invite/{token}/complete` - Complete registration

### Students

- `POST /api/v1/student` - Create student (admin only)
- `GET /api/v1/student` - List students (paginated)
- `GET /api/v1/student/{uuid}` - Get student details
- `PUT /api/v1/student/{uuid}` - Update student (admin only)
- `DELETE /api/v1/student/{uuid}` - Delete student (admin only)

### Teachers

- `POST /api/v1/teacher` - Create teacher (admin only)
- `GET /api/v1/teacher` - List teachers (paginated)
- `GET /api/v1/teacher/{uuid}` - Get teacher details
- `PUT /api/v1/teacher/{uuid}` - Update teacher (admin only)
- `DELETE /api/v1/teacher/{uuid}` - Delete teacher (admin only)

### System

- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics (requires secret)

Full API documentation available at [http://localhost:8080/docs](http://localhost:8080/docs) when the server is running.

## Configuration

### Environment Configuration

Main configuration file: `configs/config.toml`

**Key sections:**
- `[server]` - HTTP server settings, CORS, TLS
- `[database]` - PostgreSQL connection settings
- `[jwt]` - Token lifetimes and key paths
- `[logging]` - Log levels and output
- `[security]` - Rate limiting, password requirements

### Database Configuration

Default PostgreSQL settings (development):

```toml
[database]
host = "localhost"
port = 5432
user = "fits_user"
password = "fits_password"
database = "fits_db"
ssl_mode = "disable"
max_conns = 25
min_conns = 5
```

**Production:** Use environment variables or update `config.toml` with secure values.

### JWT Configuration

```toml
[jwt]
secret = "your-jwt-secret-min-32-chars"
access_token_expiry = "1h"
refresh_token_expiry = "720h"  # 30 days
invitation_expiry = "168h"      # 7 days
```

**Production:** Generate a strong random secret with `openssl rand -base64 32`.

## Troubleshooting

### Server Issues

**Server won't start:**
```bash
make db-status    # Check PostgreSQL is running
make db-up        # Start if needed
make db-logs      # Check for errors
```

**Port already in use:**
```bash
lsof -i :8080     # Find process using port 8080
kill <PID>        # Kill the process
make run          # Try again
```

### Database Issues

**Connection refused:**
```bash
make db-status    # Check container status
make db-down      # Stop container
make db-up        # Restart with health check
```

**Migration errors:**
```bash
make db-destroy   # Remove database (WARNING: data loss)
make setup        # Recreate from scratch
```

### Authentication Issues

**Bootstrap already initialized:**
- This is normal, bootstrap only runs once
- Use `POST /api/v1/auth/login` instead
- Or reset database with `make db-reset` for fresh start

**Token expired:**
- Access tokens expire after 1 hour
- Use `POST /api/v1/auth/refresh` with refresh token
- Refresh tokens expire after 30 days

### Build Issues

**Build fails:**
```bash
make clean        # Remove all build artifacts
make deps         # Reinstall dependencies
make build        # Try again
```

**Swagger generation fails:**
```bash
make tools        # Reinstall swag
make swagger      # Regenerate docs
```

## Testing

### Unit Tests

```bash
make test         # Run all unit tests
make test-verbose # Detailed output
make test-watch   # Auto-run on file changes (requires entr)
```

### Integration Tests

```bash
make run &              # Start server in background
make test-integration   # Run E2E tests
```

### Coverage Analysis

```bash
make test-coverage      # Generate HTML coverage report
xdg-open coverage.html  # View in browser
```

### Benchmarks

```bash
make bench        # Run performance benchmarks
```

## Development Tools

### Required Tools

- **Go 1.21+** - Primary language
- **Docker** - PostgreSQL container
- **Make** - Build automation

### Optional Tools

- **golangci-lint** - Meta-linter (installed via `make tools`)
- **swag** - Swagger generator (installed via `make tools`)
- **gosec** - Security scanner (installed via `make tools`)
- **entr** - File watcher for `make test-watch`

Install all optional tools:
```bash
make tools
```

## Security

### Best Practices

**Development:**
- Use `localhost:8080`, not `0.0.0.0:8080` in browser
- Keep JWT secret at least 32 characters
- Run `make security` regularly for vulnerability scans

**Production:**
- Enable TLS (`tls_enabled = true`)
- Use strong JWT secret (generate with `openssl rand -base64 32`)
- Configure CORS for specific origins (not `*`)
- Enable SSL for database connections (`ssl_mode = "require"`)
- Change all default passwords in `config.toml`
- Set appropriate rate limits
- Run `make security` and address all findings

### Security Checklist

- [ ] Change default database password
- [ ] Generate strong JWT secret (32+ characters)
- [ ] Enable TLS for API server
- [ ] Configure SSL for database
- [ ] Set specific CORS origins (not `*`)
- [ ] Review and set rate limits
- [ ] Run security audit (`make security`)
- [ ] Update all dependencies (`make deps-upgrade`)
- [ ] Test authentication flows
- [ ] Verify RBAC permissions

## Contributing

### Before Committing

Always run pre-commit checks:

```bash
make pre-commit   # Runs fmt, vet, lint, test
```

This ensures:
- Code is properly formatted (`gofmt`)
- No potential bugs (`go vet`)
- Passes linting standards (`golangci-lint`)
- All tests pass

### Code Quality

```bash
make fmt          # Format all Go code
make vet          # Check for potential bugs
make lint         # Run comprehensive linting
make check        # Complete quality check
```

### Testing Requirements

- Maintain test coverage above 80%
- Add tests for new features
- Update tests when modifying code
- Run `make test-coverage` to verify

## Technology Stack

**Backend:**
- **Go 1.21+** - Primary language
- **Fiber v2** - HTTP framework
- **GORM** - ORM for PostgreSQL
- **JWT** - Authentication (RS256 for admin, HS256 for users)

**Database:**
- **PostgreSQL 15+** - Primary database
- **Docker** - Container runtime

**Tools:**
- **Swagger/OpenAPI 3.0** - API documentation
- **golangci-lint** - Code quality
- **gosec** - Security scanning

## Project Structure

```
FITS-backend/
├── cmd/
│   └── server/         # Main application entry point
├── configs/            # Configuration files
│   ├── config.toml     # Main configuration
│   └── keys/           # RSA keypairs (auto-generated)
├── docs/               # Documentation (you are here)
│   ├── api/            # OpenAPI/Swagger files
│   ├── guides/         # User and developer guides
│   └── archive/        # Historical reports
├── internal/           # Private application code
│   ├── config/         # Configuration loading
│   ├── domain/         # Business logic (auth, student, teacher)
│   ├── middleware/     # HTTP middleware
│   └── common/         # Shared utilities
├── pkg/                # Public libraries
│   ├── crypto/         # Cryptography utilities
│   └── database/       # Database utilities
├── scripts/            # Helper scripts
│   ├── db_reset.sh     # Database reset utility
│   └── db_seed.sh      # Database seeding utility
└── Makefile            # Build automation

```

## Support

### Documentation Issues

If you find documentation errors or have suggestions:

1. Check this index for the correct guide
2. Review the [Makefile Guide](guides/makefile.md) for common commands
3. Check [Troubleshooting](#troubleshooting) section above

### Development Issues

For development problems:

1. Run `make verify` for comprehensive health check
2. Check `make db-logs` for database errors
3. Review [Swagger UI Guide](guides/swagger-ui.md) for API testing

### Getting Help

```bash
make help         # Show all available commands
make verify       # Run system health checks
```

## License

This project is licensed under the GNU General Public License v3.0 - see [LICENSE](../LICENSE) for details.

---

**Last Updated:** 2025-10-22
**Status:** Production Ready
