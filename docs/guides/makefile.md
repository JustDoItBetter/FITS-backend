# FITS Backend - Makefile Guide

## Overview

The FITS Backend Makefile provides a **complete, production-ready development workflow** with full PostgreSQL context and automation. All commands are designed to work from a fresh environment without manual steps.

---

## 🚀 Quick Start

### First Time Setup

```bash
# Clone the repository
git clone <repository-url>
cd FITS-backend

# Complete setup (installs deps, tools, starts PostgreSQL)
make setup

# Start the server
make run

# In another terminal: seed the database
make db-seed

# Open Swagger UI
make docs-serve
```

**That's it!** You now have a fully functional development environment.

---

## 📚 Command Reference

### Quick Start Commands

| Command | Description |
|---------|-------------|
| `make setup` | Complete first-time setup (deps + tools + PostgreSQL) |
| `make dev` | Start complete development environment (DB + Server) |
| `make run` | Start server only (requires PostgreSQL running) |

**Example workflow:**
```bash
make setup    # Once per machine
make dev      # Starts DB and server together
```

---

### 🔨 Building

| Command | Description | Use Case |
|---------|-------------|----------|
| `make build` | Build development binary | Testing, debugging |
| `make build-prod` | Build optimized production binary | Deployment, releases |
| `make clean` | Remove all build artifacts and caches | Fresh rebuild needed |
| `make swagger` | Regenerate Swagger/OpenAPI documentation | After changing API annotations |

**Build Details:**
- **Development build**: Includes debug symbols, faster compilation
- **Production build**: Stripped symbols (`-w -s`), CGO disabled, smaller binary
- Binary output: `bin/fits-server`

**Example:**
```bash
make build              # Development: ~50MB
make build-prod         # Production: ~35MB (optimized)
```

---

### 🗄️ Database Management

#### Status & Control

| Command | Description |
|---------|-------------|
| `make db-status` | Check PostgreSQL container status and connection |
| `make db-up` | Start PostgreSQL container (with health check) |
| `make db-down` | Stop PostgreSQL container |
| `make db-logs` | View PostgreSQL container logs (live) |

#### Data Management

| Command | Description | Warning Level |
|---------|-------------|---------------|
| `make db-reset` | Truncate all tables (preserves structure) | ⚠️ Deletes all data |
| `make db-seed` | Populate with development test data | ✅ Safe, requires running server |
| `make db-destroy` | Destroy database and Docker volumes | 🔴 **PERMANENT DATA LOSS** |

#### Maintenance

| Command | Description |
|---------|-------------|
| `make db-migrate` | Run pending migrations (auto-runs on server start) |
| `make db-verify` | Verify database schema integrity |

**Database Architecture:**
- **Automatic migrations**: Migrations run automatically when server starts
- **Migration tracking**: Uses `schema_migrations` table
- **Idempotent**: Safe to run migrations multiple times

**Database Seed Data:**
The `db-seed` command creates:
- **1 Admin account** (via bootstrap)
- **3 Teachers**:
  - Dr. Anna Schmidt (Computer Science) - `anna.schmidt` / `SecurePassword123!`
  - Prof. Michael Weber (Mathematics) - `michael.weber` / `SecurePassword123!`
  - Dr. Sarah Müller (Physics) - `sarah.mueller` / `SecurePassword123!`
- **5 Students** (all assigned to Anna Schmidt):
  - Max Mustermann - `max.mustermann` / `StudentPass123!`
  - Lisa Schneider - `lisa.schneider` / `StudentPass123!`
  - Tom Fischer - `tom.fischer` / `StudentPass123!`
  - Emma Wagner - `emma.wagner` / `StudentPass123!`
  - Paul Becker - `paul.becker` / `StudentPass123!`

**Example Workflow:**
```bash
# Check database status
make db-status

# Reset and seed for fresh data
make db-reset
make db-seed

# Verify schema
make db-verify

# View logs if issues
make db-logs
```

---

### 🧪 Testing

| Command | Description | Coverage |
|---------|-------------|----------|
| `make test` | Run all unit tests | Standard output |
| `make test-verbose` | Run tests with verbose output | Detailed logs |
| `make test-coverage` | Generate HTML coverage report | Full metrics |
| `make test-watch` | Run tests on file changes (requires `entr`) | Development |
| `make test-integration` | Run integration tests | Requires running server |
| `make bench` | Run benchmarks | Performance testing |

**Coverage Reports:**
- Terminal summary: Shows function-level coverage
- HTML report: `coverage.html` (open with `xdg-open coverage.html`)

**Example:**
```bash
# Quick test
make test

# Detailed coverage analysis
make test-coverage
xdg-open coverage.html

# Integration tests (server must be running)
make run &
make test-integration
```

---

### 🔍 Code Quality

| Command | Description | Auto-fix |
|---------|-------------|----------|
| `make fmt` | Format Go code (gofmt) | ✅ Yes |
| `make vet` | Run go vet (potential bugs) | ❌ No (reports only) |
| `make lint` | Run golangci-lint (style, issues) | ❌ No (reports only) |
| `make check` | Run fmt + vet + lint + test | Partial (fmt only) |
| `make security` | Run gosec (security vulnerabilities) | ❌ No (reports only) |

**Recommended Workflow:**
```bash
# Before committing
make pre-commit    # Runs fmt + vet + lint + test

# For comprehensive check
make check

# Security audit
make security
```

---

### 🔧 Dependencies & Tools

#### Dependency Management

| Command | Description | Safe |
|---------|-------------|------|
| `make deps` | Download and tidy Go dependencies | ✅ Yes |
| `make deps-upgrade` | Upgrade all dependencies to latest | ⚠️ Test after |
| `make deps-verify` | Verify dependency integrity | ✅ Yes |

#### Development Tools

| Command | Description | Tools Installed |
|---------|-------------|-----------------|
| `make tools` | Install all development tools | swag, golangci-lint, gosec |

**Installed Tools:**
- **swag**: Swagger/OpenAPI generator
- **golangci-lint**: Meta-linter (runs multiple linters)
- **gosec**: Security vulnerability scanner

**Example:**
```bash
# Initial setup
make tools

# Regular updates
make deps

# Major upgrade (careful!)
make deps-upgrade
make test  # Verify everything still works
```

---

### 🎯 Workflow Commands

These are high-level commands that combine multiple operations:

| Command | What It Does | Use Case |
|---------|--------------|----------|
| `make fresh` | Destroy + setup + seed | Complete fresh start |
| `make ci` | deps + check + coverage | CI/CD pipeline |
| `make pre-commit` | fmt + vet + lint + test | Before git commit |
| `make verify` | Comprehensive system check | Health check |
| `make docker-clean` | Remove all Docker resources | Cleanup |

**Workflow Examples:**

**Fresh Start (New Day):**
```bash
make fresh
make run
```

**Before Committing:**
```bash
make pre-commit
git add .
git commit -m "Your message"
```

**CI/CD Pipeline:**
```bash
make ci  # This is what CI runs
```

**System Health Check:**
```bash
make verify
```

---

## 📚 Documentation Commands

| Command | Description |
|---------|-------------|
| `make docs` | Regenerate all documentation (swagger) |
| `make docs-serve` | Open Swagger UI in browser |

**Example:**
```bash
# After modifying API annotations
make swagger

# View documentation
make docs-serve
```

---

## 🎨 Colored Output

The Makefile uses colored output for better readability:

- **🔵 Cyan (→)**: Action in progress
- **🟢 Green (✓)**: Success
- **🟡 Yellow (⚠)**: Warning
- **🔴 Red (✗)**: Error
- **ℹ️ Blue**: Information

---

## 🔄 Common Workflows

### Daily Development

```bash
# Morning
make db-status          # Check if DB is running
make db-up              # Start if needed
make run                # Start server

# Development
make test-watch         # Auto-run tests on save

# Before commit
make pre-commit         # Quality checks
```

### Fresh Environment

```bash
# Complete reset
make fresh              # Destroys everything, rebuilds, seeds

# Or step by step
make db-destroy         # Confirm with 'yes'
make setup              # Setup from scratch
make build              # Build binary
make run                # Start server
make db-seed            # Add test data
```

### Testing Workflow

```bash
# Unit tests
make test               # Quick test
make test-coverage      # With coverage

# Integration tests
make run &              # Start server in background
make test-integration   # Run E2E tests
```

### Production Build

```bash
# Clean build
make clean
make build-prod

# Verify
./bin/fits-server --version

# Run with production config
./bin/fits-server
```

### Debugging Database

```bash
# Check status
make db-status

# View logs
make db-logs

# Reset and reseed
make db-reset
make db-seed

# Verify schema
make db-verify
```

---

## 🔧 Customization

### Environment Variables

You can override database configuration:

```bash
# Custom database settings
DB_HOST=localhost \
DB_PORT=5432 \
DB_NAME=fits_db \
DB_USER=fits_user \
DB_PASSWORD=fits_password \
make db-status
```

### Server URL

For remote servers:

```bash
SERVER_URL=http://production:8080 make db-seed
```

---

## 🐛 Troubleshooting

### "PostgreSQL container not running"

```bash
make db-up              # Start PostgreSQL
make db-status          # Verify
```

### "Server not running" (for db-seed)

```bash
make run &              # Start server in background
make db-seed            # Now seed
```

### "swag not found"

```bash
make tools              # Install all dev tools
```

### "Tests failing after deps-upgrade"

```bash
git checkout go.mod go.sum    # Revert
make deps                     # Reinstall stable versions
```

### "Database connection refused"

```bash
make db-logs            # Check PostgreSQL logs
make db-down
make db-up              # Restart
```

### Build Issues

```bash
make clean              # Remove all caches
make deps               # Reinstall dependencies
make build              # Try again
```

---

## 📝 Best Practices

### ✅ DO

- Run `make pre-commit` before every commit
- Use `make fresh` when switching branches
- Run `make verify` after major changes
- Keep dependencies updated regularly (`make deps-upgrade`)
- Use `make test-coverage` to maintain high coverage

### ❌ DON'T

- Don't run `make db-destroy` without backup in production
- Don't skip `make pre-commit` checks
- Don't commit without running tests
- Don't use development builds in production
- Don't upgrade dependencies without testing

---

## 🎓 Learning Path

### Beginner

1. `make help` - Learn available commands
2. `make setup` - First-time setup
3. `make run` - Start server
4. `make test` - Run tests
5. `make db-seed` - Add test data

### Intermediate

1. `make build` - Build binaries
2. `make test-coverage` - Coverage analysis
3. `make check` - Quality checks
4. `make swagger` - Update documentation
5. `make db-reset` - Database management

### Advanced

1. `make ci` - CI/CD pipeline
2. `make build-prod` - Production builds
3. `make security` - Security audits
4. `make verify` - System health checks
5. `make fresh` - Complete environment reset

---

## 📊 Performance Tips

### Faster Builds

```bash
# Use build cache
make build              # Subsequent builds are faster

# Parallel tests
make test               # Already uses -p (parallel)
```

### Faster Tests

```bash
# Skip slow tests in development
go test ./... -short    # Skips tests marked with testing.Short()

# Run specific package
make test ARGS="./internal/domain/auth"
```

---

## 🔐 Security

### Security Checks

```bash
make security           # Run gosec scanner
```

### Production Checklist

- [ ] Run `make build-prod` (not `make build`)
- [ ] Change all default passwords in `config.toml`
- [ ] Enable TLS (`tls_enabled = true`)
- [ ] Set strong JWT secret
- [ ] Configure CORS for specific origins (not `*`)
- [ ] Run `make security` and address findings
- [ ] Set appropriate rate limits

---

## 📖 Additional Resources

- **Swagger UI**: http://localhost:8080/docs
- **Health Check**: http://localhost:8080/health
- **PostgreSQL**: `make db-logs` for live logs
- **Coverage**: Open `coverage.html` after `make test-coverage`

---

## 💡 Tips & Tricks

### Alias for Convenience

Add to your `.bashrc` or `.zshrc`:

```bash
alias mr="make run"
alias mt="make test"
alias mf="make fresh"
alias mpc="make pre-commit"
```

### VS Code Integration

Add to `.vscode/tasks.json`:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "FITS: Run Server",
      "type": "shell",
      "command": "make run",
      "problemMatcher": []
    },
    {
      "label": "FITS: Run Tests",
      "type": "shell",
      "command": "make test",
      "problemMatcher": []
    }
  ]
}
```

---

## 🆘 Getting Help

1. **Command help**: `make help`
2. **Check status**: `make verify`
3. **View logs**: `make db-logs`
4. **Fresh start**: `make fresh`

Still having issues? Check:
- Go version: `go version` (should be 1.21+)
- Docker: `docker --version`
- Database: `make db-status`

---

**Made with ❤️ for FITS Backend Development**
