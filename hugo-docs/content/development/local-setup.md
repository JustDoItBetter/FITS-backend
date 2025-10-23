---
title: "Local Development Setup"
weight: 1
---

# Local Development Setup

## Development Environment

### Recommended Tools

#### IDE / Editor

- **VS Code** (recommended)
  - Extensions:
    - Go (official Go extension)
    - GitLens
    - REST Client
    - PostgreSQL Explorer
    - YAML
    - TOML Language Support
- **GoLand** (JetBrains)
- **Vim/Neovim** with Go plugins

#### Command Line Tools

```bash
# Go tools
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/air-verse/air@latest  # Hot reload

# Database tools
brew install postgresql-client  # macOS
apt install postgresql-client   # Ubuntu

# API testing
brew install httpie
brew install jq
```

### Project Setup

1. **Fork and Clone**

```bash
# Fork the repository on GitHub
git clone https://github.com/YOUR_USERNAME/FITS-backend.git
cd FITS-backend

# Add upstream remote
git remote add upstream https://github.com/JustDoItBetter/FITS-backend.git
```

2. **Install Dependencies**

```bash
go mod download
go mod verify
```

3. **Configure Environment**

```bash
cp .env.example .env
cp configs/config.toml.example configs/config.toml  # If available
```

Edit `configs/config.toml` for local development:

```toml
[server]
host = "0.0.0.0"
port = 8080
rate_limit = 1000  # Higher for development

[database]
host = "localhost"
port = 5432
user = "fits_dev"
password = "dev_password"
database = "fits_dev"

[logging]
level = "debug"     # More verbose logging
format = "console"  # Human-readable

[jwt]
secret = "dev-secret-change-in-production"
access_token_expiry = "60m"  # Longer for development
```

4. **Set Up Database**

```bash
# Create development database
createdb -U postgres fits_dev

# Or use Docker
docker run --name fits-dev-postgres \
  -e POSTGRES_DB=fits_dev \
  -e POSTGRES_USER=fits_dev \
  -e POSTGRES_PASSWORD=dev_password \
  -p 5432:5432 \
  -d postgres:15
```

5. **Run Migrations**

Migrations run automatically on server start, or manually:

```bash
go run cmd/server/main.go
# Migrations will run on startup
```

## Development Workflow

### Running the Server

#### Standard Run

```bash
go run cmd/server/main.go
```

#### With Makefile

```bash
make run
```

#### With Hot Reload (Air)

Install Air:
```bash
go install github.com/air-verse/air@latest
```

Create `.air.toml`:
```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/server cmd/server/main.go"
  bin = "tmp/server"
  include_ext = ["go", "toml", "yaml"]
  exclude_dir = ["tmp", "vendor", "docs"]
  delay = 1000
```

Run:
```bash
air
# Server will automatically reload on file changes
```

### Code Generation

#### Generate Swagger Docs

After modifying API endpoints or annotations:

```bash
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

# Or with Makefile
make docs
```

#### Generate Mocks (if using)

```bash
# Install mockgen
go install github.com/golang/mock/mockgen@latest

# Generate mocks
go generate ./...
```

### Testing

#### Run All Tests

```bash
go test ./...

# Or with Makefile
make test
```

#### Run Tests with Coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Or with Makefile
make test-cover
```

#### Run Specific Package Tests

```bash
go test ./internal/domain/auth/...
go test ./pkg/crypto/...
```

#### Run with Race Detector

```bash
go test -race ./...

# Or with Makefile
make test-race
```

#### Verbose Output

```bash
go test -v ./internal/domain/auth/...
```

#### Run Specific Test

```bash
go test -run TestAuthService_Login ./internal/domain/auth/
```

### Code Quality

#### Format Code

```bash
go fmt ./...

# Or with Makefile
make fmt
```

#### Run Linters

```bash
golangci-lint run

# Or with Makefile
make lint
```

#### Run go vet

```bash
go vet ./...

# Or with Makefile
make vet
```

#### Pre-commit Checks

Run all checks before committing:

```bash
make fmt && make vet && make lint && make test
```

### Debugging

#### VS Code Debug Configuration

Create `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Server",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/server",
      "env": {
        "FITS_LOG_LEVEL": "debug"
      },
      "args": []
    },
    {
      "name": "Attach to Process",
      "type": "go",
      "request": "attach",
      "mode": "local",
      "processId": "${command:pickProcess}"
    },
    {
      "name": "Test Current File",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${file}"
    }
  ]
}
```

#### Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug server
dlv debug cmd/server/main.go

# Debug tests
dlv test ./internal/domain/auth/
```

#### Debug Logging

Enable debug logging in `configs/config.toml`:

```toml
[logging]
level = "debug"
format = "console"
```

Or via environment:

```bash
FITS_LOG_LEVEL=debug go run cmd/server/main.go
```

### Database Management

#### View Database Schema

```bash
psql -U fits_dev fits_dev

\dt          # List tables
\d users     # Describe users table
\q           # Quit
```

#### Reset Database

```bash
# Drop and recreate
dropdb -U postgres fits_dev
createdb -U postgres fits_dev

# Restart server to run migrations
```

#### Seed Test Data

Create a seed script `scripts/seed.sh`:

```bash
#!/bin/bash

TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/bootstrap/init \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@dev.local",
    "password": "DevAdmin123!",
    "first_name": "Dev",
    "last_name": "Admin"
  }' | jq -r '.data.token')

# Create test students
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/student \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
      \"first_name\": \"Student\",
      \"last_name\": \"$i\",
      \"email\": \"student$i@dev.local\",
      \"student_id\": \"S$(printf '%03d' $i)\"
    }"
done
```

### API Testing

#### Using curl

```bash
# Bootstrap
curl -X POST http://localhost:8080/api/v1/bootstrap/init \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@dev.local","password":"DevPass123!","first_name":"Dev","last_name":"Admin"}'
```

#### Using HTTPie

```bash
# Bootstrap
http POST :8080/api/v1/bootstrap/init \
  email=admin@dev.local \
  password=DevPass123! \
  first_name=Dev \
  last_name=Admin

# Login
http POST :8080/api/v1/auth/login \
  email=admin@dev.local \
  password=DevPass123!

# Authenticated request
http GET :8080/api/v1/student \
  "Authorization: Bearer $TOKEN"
```

#### VS Code REST Client

Create `api-tests.http`:

```http
### Variables
@baseUrl = http://localhost:8080
@token = {{login.response.body.data.token}}

### Bootstrap Admin
POST {{baseUrl}}/api/v1/bootstrap/init
Content-Type: application/json

{
  "email": "admin@dev.local",
  "password": "DevPass123!",
  "first_name": "Dev",
  "last_name": "Admin"
}

### Login
# @name login
POST {{baseUrl}}/api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@dev.local",
  "password": "DevPass123!"
}

### List Students
GET {{baseUrl}}/api/v1/student
Authorization: Bearer {{token}}
```

### Project Structure

```
FITS-backend/
├── cmd/
│   └── server/
│       ├── main.go              # Entry point
│       └── swagger_handler.go   # Swagger configuration
├── configs/
│   └── config.toml              # Configuration file
├── internal/
│   ├── common/                  # Shared utilities
│   │   ├── errors/
│   │   ├── pagination/
│   │   ├── response/
│   │   └── validation/
│   ├── config/                  # Configuration logic
│   ├── domain/                  # Business domains
│   │   ├── auth/
│   │   ├── signing/
│   │   ├── student/
│   │   └── teacher/
│   └── middleware/              # HTTP middleware
├── pkg/                         # Public packages
│   ├── crypto/                  # Cryptography utilities
│   ├── database/                # Database connection
│   └── logger/                  # Logging
├── docs/                        # Generated Swagger docs
├── migrations/                  # Database migrations
├── scripts/                     # Utility scripts
├── web/                         # Static web files
├── .env.example                 # Environment template
├── Makefile                     # Build automation
└── go.mod                       # Go module definition
```

### Common Tasks

#### Add a New Endpoint

1. Define model in `internal/domain/{entity}/model.go`
2. Add repository method in `repository.go`
3. Implement in `repository_gorm.go`
4. Add business logic in `service.go`
5. Create handler in `handler.go`
6. Add Swagger annotations
7. Register route in `cmd/server/main.go`
8. Generate docs: `make docs`
9. Write tests

#### Add a New Domain

1. Create directory: `internal/domain/{entity}/`
2. Create files:
   - `model.go` - Domain models
   - `repository.go` - Repository interface
   - `repository_gorm.go` - GORM implementation
   - `service.go` - Business logic
   - `handler.go` - HTTP handlers
   - `*_test.go` - Tests
3. Register in `cmd/server/main.go`

#### Add Middleware

1. Create `internal/middleware/{name}.go`
2. Implement `fiber.Handler`
3. Apply in `cmd/server/main.go`

### Git Workflow

```bash
# Create feature branch
git checkout -b feature/add-new-endpoint

# Make changes, commit
git add .
git commit -m "feat: add new endpoint for user profile"

# Keep up to date with upstream
git fetch upstream
git rebase upstream/main

# Push to your fork
git push origin feature/add-new-endpoint

# Create pull request on GitHub
```

### Commit Message Convention

Follow conventional commits:

```
feat: Add new feature
fix: Fix bug
docs: Update documentation
style: Format code
refactor: Refactor code
test: Add tests
chore: Update dependencies
```

## Next Steps

- [Testing Strategy](/development/testing-strategy/)
- [Contribution Guidelines](/development/contribution-guidelines/)
- [API Reference](/api/endpoints/)
