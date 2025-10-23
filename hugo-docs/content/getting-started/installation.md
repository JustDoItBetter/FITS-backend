---
title: "Installation"
weight: 1
description: "Complete installation guide for FITS Backend. Prerequisites, PostgreSQL setup, environment configuration, and first run instructions."
---

# Installation Guide

## Prerequisites

### Required

- **Go**: Version 1.25.1 or higher
- **PostgreSQL**: Version 12 or higher (production)
- **Git**: For cloning the repository

### Optional

- **Docker**: For containerized deployment
- **Make**: For using Makefile commands
- **Swag**: For regenerating API documentation

## Installation Steps

### 1. Clone the Repository

```bash
git clone https://github.com/JustDoItBetter/FITS-backend.git
cd FITS-backend
```

### 2. Install Dependencies

```bash
go mod download
```

This downloads all required Go modules defined in `go.mod`.

### 3. Install Development Tools

#### Swag (API Documentation Generator)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

#### golangci-lint (Linter)

```bash
# macOS/Linux
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Or via Go
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 4. Set Up PostgreSQL

#### Option A: Local PostgreSQL

Install PostgreSQL for your platform:

**macOS (Homebrew):**
```bash
brew install postgresql@15
brew services start postgresql@15
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
```

**Create Database:**
```bash
# Login as postgres user
sudo -u postgres psql

# Create database and user
CREATE DATABASE fits_db;
CREATE USER fits WITH ENCRYPTED PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE fits_db TO fits;
\q
```

#### Option B: Docker PostgreSQL

```bash
docker run --name fits-postgres \
  -e POSTGRES_DB=fits_db \
  -e POSTGRES_USER=fits \
  -e POSTGRES_PASSWORD=your_secure_password \
  -p 5432:5432 \
  -d postgres:15
```

#### Option C: Docker Compose

The repository includes a `docker-compose.yml`:

```bash
docker-compose up -d postgres
```

### 5. Configure the Application

Copy the example configuration:

```bash
cp .env.example .env
```

Edit `configs/config.toml`:

```toml
[server]
host = "0.0.0.0"
port = 8080
read_timeout = "30s"
write_timeout = "30s"
rate_limit = 100
allowed_origins = "*"  # Use specific origins in production!

[database]
host = "localhost"
port = 5432
user = "fits"
password = "your_secure_password"
database = "fits_db"
sslmode = "disable"  # Use "require" in production!

[jwt]
secret = "change-this-to-a-secure-random-string"
access_token_expiry = "15m"
refresh_token_expiry = "168h"  # 7 days

[logging]
level = "info"  # debug, info, warn, error
format = "console"  # console or json

[storage]
upload_dir = "./uploads"
max_file_size = 10485760  # 10MB

[secrets]
metrics_secret = "change-this-metrics-secret"

[server.tls]
enabled = false
cert_file = ""
key_file = ""
```

**Important Security Notes:**
- Change the JWT secret to a secure random string (min 32 characters)
- Use strong database passwords
- Enable TLS in production
- Restrict CORS origins in production
- Use SSL/TLS for database connections in production

### 6. Generate API Documentation

```bash
# Using swag directly
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

# Or using Makefile
make docs
```

This generates Swagger documentation in the `docs/` directory.

### 7. Build the Application

```bash
# Build binary
go build -o server cmd/server/main.go

# Or using Makefile
make build
```

### 8. Run Database Migrations

Migrations run automatically on first startup. The application will create all necessary tables.

### 9. Start the Server

```bash
# Run directly
./server

# Or using go run
go run cmd/server/main.go

# Or using Makefile
make run
```

You should see output like:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
FITS Backend v1.0.1 - Ready
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

⚠️  HTTP Mode (Development Only)
Server:        http://localhost:8080 (listening on 0.0.0.0:8080)
Browser:       http://localhost:8080 👈 USE THIS IN BROWSER
Documentation: http://localhost:8080/docs
Health Check:  http://localhost:8080/health

Quick Start:
  1. Bootstrap Admin: POST http://localhost:8080/api/v1/bootstrap/init
  2. Login:           POST http://localhost:8080/api/v1/auth/login
  3. View API Docs:   http://localhost:8080/docs

⚠️  IMPORTANT: Use http://localhost:8080 in your browser
   Do NOT use http://0.0.0.0:8080 - browsers will block it!

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 10. Verify Installation

Check the health endpoint:

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "database": "connected",
  "time": "2025-10-23T18:00:00Z"
}
```

View API documentation:
- Open browser: http://localhost:8080/docs

## Makefile Commands

The project includes a comprehensive Makefile:

```bash
# Development
make run          # Run the server
make build        # Build binary
make dev          # Run with hot reload (requires air)

# Testing
make test         # Run all tests
make test-cover   # Run tests with coverage report
make test-race    # Run tests with race detector

# Code Quality
make fmt          # Format code
make lint         # Run linters
make vet          # Run go vet

# Documentation
make docs         # Generate Swagger docs

# Database
make migrate-up   # Run migrations
make migrate-down # Rollback migrations

# Docker
make docker-build # Build Docker image
make docker-run   # Run Docker container

# Cleanup
make clean        # Remove build artifacts
```

See `Makefile` for all available commands.

## Environment Variables

The application supports environment variables as an alternative to TOML configuration:

```bash
# Server
export FITS_SERVER_HOST="0.0.0.0"
export FITS_SERVER_PORT="8080"

# Database
export FITS_DB_HOST="localhost"
export FITS_DB_PORT="5432"
export FITS_DB_USER="fits"
export FITS_DB_PASSWORD="your_password"
export FITS_DB_NAME="fits_db"

# JWT
export FITS_JWT_SECRET="your-jwt-secret"
export FITS_JWT_ACCESS_EXPIRY="15m"
export FITS_JWT_REFRESH_EXPIRY="168h"

# Logging
export FITS_LOG_LEVEL="info"
export FITS_LOG_FORMAT="json"
```

## Docker Installation

### Using Docker

Build and run with Docker:

```bash
# Build image
docker build -t fits-backend:latest .

# Run container
docker run -d \
  --name fits-backend \
  -p 8080:8080 \
  -e FITS_DB_HOST=host.docker.internal \
  -e FITS_DB_PASSWORD=your_password \
  -e FITS_JWT_SECRET=your_secret \
  fits-backend:latest
```

### Using Docker Compose

The repository includes a complete Docker Compose setup:

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

This starts:
- PostgreSQL database
- FITS Backend API
- (Optional) Prometheus
- (Optional) Grafana

## Troubleshooting

### Common Issues

#### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

#### Database Connection Failed

1. Verify PostgreSQL is running:
   ```bash
   # macOS/Linux
   pg_isready -h localhost -p 5432
   
   # Docker
   docker ps | grep postgres
   ```

2. Check credentials in `configs/config.toml`

3. Verify network connectivity:
   ```bash
   telnet localhost 5432
   ```

#### Permission Denied on Binary

```bash
chmod +x server
```

#### Missing Dependencies

```bash
# Clean and reinstall
go clean -modcache
go mod download
go mod verify
```

#### Swagger Documentation Not Found

```bash
# Regenerate docs
make docs

# Or manually
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
```

### Debug Mode

Run with debug logging:

```bash
# Edit configs/config.toml
[logging]
level = "debug"

# Or use environment variable
FITS_LOG_LEVEL=debug ./server
```

### Database Debugging

Enable GORM debug logging in `pkg/database/database.go`:

```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

## Next Steps

- [Quick Start Guide](/getting-started/quick-start/) - Make your first API requests
- [Development Setup](/development/local-setup/) - Set up your development environment
- [API Documentation](http://localhost:8080/docs) - Explore the API endpoints

## Production Deployment

For production deployment, see:
- [Deployment Guide](/infrastructure/deployment/)
- [Security Considerations](/infrastructure/security-considerations/)
