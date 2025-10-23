# ============================================================================
# FITS Backend - Professional Development Makefile
# ============================================================================
# Provides a complete, modern development workflow with PostgreSQL context
# All commands are designed to work from a fresh environment
# ============================================================================

.PHONY: help

# Default target - show help
help:
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "FITS Backend - Development Makefile"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""
	@echo "🚀 Quick Start"
	@echo "  make setup           - Complete first-time setup (deps + Docker + wait)"
	@echo "  make dev             - Start development environment (DB + Server)"
	@echo "  make run             - Start server only (requires DB running)"
	@echo ""
	@echo "🔨 Building"
	@echo "  make build           - Build development binary"
	@echo "  make build-prod      - Build optimized production binary"
	@echo "  make clean           - Remove build artifacts and cache"
	@echo "  make swagger         - Regenerate Swagger/OpenAPI documentation"
	@echo ""
	@echo "🗄️  Database Management"
	@echo "  make db-status       - Check PostgreSQL status and connection"
	@echo "  make db-up           - Start PostgreSQL container"
	@echo "  make db-down         - Stop PostgreSQL container"
	@echo "  make db-logs         - View PostgreSQL container logs"
	@echo "  make db-reset        - Reset database (truncate all tables)"
	@echo "  make db-seed         - Populate database with development data"
	@echo "  make db-destroy      - Destroy database and volumes (complete wipe)"
	@echo "  make db-migrate      - Run pending migrations (auto-runs on startup)"
	@echo "  make db-verify       - Verify database schema integrity"
	@echo ""
	@echo "🧪 Testing"
	@echo "  make test            - Run all unit tests"
	@echo "  make test-verbose    - Run tests with verbose output"
	@echo "  make test-coverage   - Generate coverage report (HTML + terminal)"
	@echo "  make test-watch      - Run tests in watch mode"
	@echo "  make test-integration - Run integration tests (requires running server)"
	@echo "  make bench           - Run benchmarks"
	@echo ""
	@echo "🔍 Code Quality"
	@echo "  make fmt             - Format Go code (gofmt)"
	@echo "  make lint            - Run linter (golangci-lint)"
	@echo "  make vet             - Run go vet"
	@echo "  make check           - Run fmt + vet + lint + test"
	@echo "  make security        - Run security checks (gosec)"
	@echo ""
	@echo "🔧 Utilities"
	@echo "  make deps            - Download and tidy Go dependencies"
	@echo "  make deps-upgrade    - Upgrade all dependencies to latest"
	@echo "  make deps-verify     - Verify dependencies integrity"
	@echo "  make tools           - Install required development tools"
	@echo "  make verify          - Comprehensive system verification"
	@echo ""
	@echo "🎯 Workflows"
	@echo "  make fresh           - Complete fresh start (destroy + setup + build)"
	@echo "  make fresh-seed      - Fresh start + auto-seed (starts server in background)"
	@echo "  make ci              - Run CI pipeline (deps + check + coverage)"
	@echo "  make pre-commit      - Run before committing (fmt + check)"
	@echo "  make docker-clean    - Remove all Docker containers and volumes"
	@echo ""
	@echo "📚 Documentation"
	@echo "  make docs            - Generate all documentation"
	@echo "  make docs-serve      - Serve Swagger UI (opens browser)"
	@echo ""

# ============================================================================
# Configuration
# ============================================================================

# Go parameters
GO := go
GOCMD := $(GO)
GOBUILD := $(GO) build
GOTEST := $(GO) test
GOMOD := $(GO) mod
GOVET := $(GO) vet
GOFMT := $(GO) fmt

# Binary configuration
BINARY_NAME := fits-server
BINARY_PATH := bin/$(BINARY_NAME)
CMD_PATH := ./cmd/server

# Build flags
LDFLAGS := -w -s
BUILD_FLAGS := -trimpath
PROD_FLAGS := $(BUILD_FLAGS) -ldflags="$(LDFLAGS)"

# PostgreSQL configuration (from config.toml)
DB_HOST := localhost
DB_PORT := 5432
DB_NAME := fits_db
DB_USER := fits_user
DB_PASSWORD := fits_password
PGPASSWORD := $(DB_PASSWORD)

# Docker configuration
DOCKER_COMPOSE := docker-compose
DOCKER := docker
CONTAINER_NAME := fits_postgres

# Paths
SCRIPTS_DIR := scripts
MIGRATIONS_DIR := migrations
UPLOADS_DIR := uploads
KEYS_DIR := configs/keys
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

# Colors for output
COLOR_RESET := \033[0m
COLOR_BOLD := \033[1m
COLOR_RED := \033[31m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_BLUE := \033[34m
COLOR_CYAN := \033[36m

# ============================================================================
# Quick Start & Development
# ============================================================================

.PHONY: setup dev run

# Complete first-time setup
setup: deps tools db-up
	@echo "$(COLOR_BLUE)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(COLOR_RESET)"
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Setup complete!"
	@echo ""
	@echo "$(COLOR_CYAN)Next steps:$(COLOR_RESET)"
	@echo "  1. $(COLOR_YELLOW)make run$(COLOR_RESET)       - Start the server"
	@echo "  2. $(COLOR_YELLOW)make db-seed$(COLOR_RESET)   - Add development data (optional)"
	@echo "  3. Visit $(COLOR_YELLOW)http://localhost:8080/docs$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(COLOR_RESET)"

# Start complete development environment
dev: db-up
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Starting development environment..."
	@sleep 2
	@$(MAKE) run

# Run server (requires PostgreSQL to be running)
run:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Starting FITS Backend Server..."
	@$(GOCMD) run $(CMD_PATH)/main.go

# ============================================================================
# Building
# ============================================================================

.PHONY: build build-prod clean

# Build development binary
build:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Building development binary..."
	@mkdir -p bin
	@$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_PATH) $(CMD_PATH)/main.go
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Binary created: $(BINARY_PATH)"

# Build optimized production binary
build-prod:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Building production binary (optimized)..."
	@mkdir -p bin
	@CGO_ENABLED=0 $(GOBUILD) $(PROD_FLAGS) -o $(BINARY_PATH) $(CMD_PATH)/main.go
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Production binary created: $(BINARY_PATH)"
	@ls -lh $(BINARY_PATH)

# Clean build artifacts
clean:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@rm -f *.test
	@rm -f *.out
	@$(GO) clean -cache -testcache -modcache
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Cleanup complete"

# ============================================================================
# Database Management
# ============================================================================

.PHONY: db-status db-up db-down db-logs db-reset db-seed db-destroy db-migrate db-verify

# Check PostgreSQL status
db-status:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Checking PostgreSQL status..."
	@if $(DOCKER) ps | grep -q $(CONTAINER_NAME); then \
		echo "$(COLOR_GREEN)✓$(COLOR_RESET) PostgreSQL container is running"; \
		$(DOCKER) exec $(CONTAINER_NAME) pg_isready -U $(DB_USER) -d $(DB_NAME) && \
		echo "$(COLOR_GREEN)✓$(COLOR_RESET) Database is accepting connections" || \
		echo "$(COLOR_YELLOW)⚠$(COLOR_RESET) Database is not ready yet"; \
	else \
		echo "$(COLOR_RED)✗$(COLOR_RESET) PostgreSQL container is not running"; \
		echo "  Run: $(COLOR_YELLOW)make db-up$(COLOR_RESET)"; \
	fi

# Start PostgreSQL container
db-up:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Starting PostgreSQL container..."
	@$(DOCKER_COMPOSE) up -d
	@echo "$(COLOR_YELLOW)⏳$(COLOR_RESET) Waiting for PostgreSQL to be ready..."
	@for i in 1 2 3 4 5 6 7 8 9 10; do \
		if $(DOCKER) exec $(CONTAINER_NAME) pg_isready -U $(DB_USER) -d $(DB_NAME) >/dev/null 2>&1; then \
			echo "$(COLOR_GREEN)✓$(COLOR_RESET) PostgreSQL is ready on $(DB_HOST):$(DB_PORT)"; \
			break; \
		fi; \
		if [ $$i -eq 10 ]; then \
			echo "$(COLOR_RED)✗$(COLOR_RESET) Timeout waiting for PostgreSQL"; \
			exit 1; \
		fi; \
		sleep 1; \
	done

# Stop PostgreSQL container
db-down:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Stopping PostgreSQL container..."
	@$(DOCKER_COMPOSE) down
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) PostgreSQL stopped"

# View PostgreSQL logs
db-logs:
	@$(DOCKER_COMPOSE) logs -f postgres

# Reset database (truncate all tables)
db-reset:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Resetting database..."
	@$(SCRIPTS_DIR)/db_reset.sh
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Database reset complete"

# Populate database with development seed data
db-seed:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Seeding database with development data..."
	@$(SCRIPTS_DIR)/db_seed.sh
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Database seeded successfully"

# Destroy database completely (including Docker volumes)
db-destroy:
	@echo "$(COLOR_YELLOW)⚠$(COLOR_RESET)  This will destroy all data permanently!"
	@read -p "Continue? (yes/no): " confirm && [ "$$confirm" = "yes" ] || exit 1
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Destroying database and volumes..."
	@$(DOCKER_COMPOSE) down -v
	@rm -rf $(KEYS_DIR)/*.key $(KEYS_DIR)/*.pub 2>/dev/null || true
	@rm -rf $(UPLOADS_DIR)/* 2>/dev/null || true
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Database destroyed"

# Run database migrations (migrations auto-run on server start)
db-migrate:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Running database migrations..."
	@echo "$(COLOR_YELLOW)ℹ$(COLOR_RESET)  Note: Migrations run automatically when server starts"
	@echo "  Starting server to trigger migrations..."
	@timeout 5s $(MAKE) run 2>/dev/null || true
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Migrations complete"

# Verify database schema integrity
db-verify:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Verifying database schema..."
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -c "\
		SELECT table_name FROM information_schema.tables \
		WHERE table_schema = 'public' \
		ORDER BY table_name;" -t | grep -v '^$$' | while read table; do \
		echo "  ✓ $$table"; \
	done
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Schema verification complete"

# ============================================================================
# Swagger Documentation
# ============================================================================

.PHONY: swagger docs docs-serve

# Regenerate Swagger documentation
swagger:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Regenerating Swagger documentation..."
	@if [ -f ~/go/bin/swag ]; then \
		~/go/bin/swag init -g $(CMD_PATH)/main.go -o docs --parseDependency --parseInternal && \
		echo "$(COLOR_GREEN)✓$(COLOR_RESET) Swagger docs regenerated"; \
	else \
		echo "$(COLOR_RED)✗$(COLOR_RESET) swag not found"; \
		echo "  Install with: $(COLOR_YELLOW)make tools$(COLOR_RESET)"; \
		exit 1; \
	fi

# Generate all documentation
docs: swagger
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Documentation generated"
	@echo "  View at: $(COLOR_YELLOW)http://localhost:8080/docs$(COLOR_RESET)"

# Serve Swagger UI (requires running server)
docs-serve:
	@if curl -s http://localhost:8080/health >/dev/null 2>&1; then \
		echo "$(COLOR_GREEN)✓$(COLOR_RESET) Server is running"; \
		echo "  Swagger UI: $(COLOR_YELLOW)http://localhost:8080/docs$(COLOR_RESET)"; \
		command -v xdg-open >/dev/null && xdg-open http://localhost:8080/docs || true; \
	else \
		echo "$(COLOR_RED)✗$(COLOR_RESET) Server is not running"; \
		echo "  Start with: $(COLOR_YELLOW)make run$(COLOR_RESET)"; \
	fi

# ============================================================================
# Testing
# ============================================================================

.PHONY: test test-verbose test-coverage test-watch test-integration bench

# Run all unit tests
test:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Running unit tests..."
	@$(GOTEST) ./... -short
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Tests passed"

# Run tests with verbose output
test-verbose:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Running tests (verbose)..."
	@$(GOTEST) ./... -v -short

# Generate coverage report
test-coverage:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Running tests with coverage..."
	@$(GOTEST) ./... -coverprofile=$(COVERAGE_FILE) -covermode=atomic
	@echo ""
	@echo "$(COLOR_CYAN)Coverage Summary:$(COLOR_RESET)"
	@$(GO) tool cover -func=$(COVERAGE_FILE) | tail -20
	@echo ""
	@$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Coverage report generated"
	@echo "  HTML report: $(COLOR_YELLOW)$(COVERAGE_HTML)$(COLOR_RESET)"
	@echo "  Open with: $(COLOR_YELLOW)xdg-open $(COVERAGE_HTML)$(COLOR_RESET)"

# Run tests in watch mode (requires entr)
test-watch:
	@if command -v entr >/dev/null 2>&1; then \
		echo "$(COLOR_CYAN)→$(COLOR_RESET) Running tests in watch mode (Ctrl+C to stop)..."; \
		find . -name '*.go' | entr -c make test; \
	else \
		echo "$(COLOR_RED)✗$(COLOR_RESET) entr not found"; \
		echo "  Install with: $(COLOR_YELLOW)apt-get install entr$(COLOR_RESET) (or your package manager)"; \
	fi

# Run integration tests (requires running server)
test-integration:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Running integration tests..."
	@if ! curl -s http://localhost:8080/health >/dev/null 2>&1; then \
		echo "$(COLOR_RED)✗$(COLOR_RESET) Server not running"; \
		echo "  Start with: $(COLOR_YELLOW)make run$(COLOR_RESET)"; \
		exit 1; \
	fi
	@$(SCRIPTS_DIR)/test_full_flow.sh
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Integration tests passed"

# Run benchmarks
bench:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Running benchmarks..."
	@$(GOTEST) ./pkg/crypto -bench=. -benchmem -benchtime=3s
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Benchmarks complete"

# ============================================================================
# Code Quality
# ============================================================================

.PHONY: fmt lint vet check security

# Format Go code
fmt:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Formatting code..."
	@$(GOFMT) ./...
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Code formatted"

# Run linter
lint:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "$(COLOR_GREEN)✓$(COLOR_RESET) Linting passed"; \
	else \
		echo "$(COLOR_YELLOW)⚠$(COLOR_RESET) golangci-lint not installed"; \
		echo "  Install with: $(COLOR_YELLOW)make tools$(COLOR_RESET)"; \
		echo "  Or visit: https://golangci-lint.run/usage/install/"; \
	fi

# Run go vet
vet:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Running go vet..."
	@$(GOVET) ./...
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Vet passed"

# Run all checks
check: fmt vet lint test
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) All checks passed"

# Run security checks
security:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec -quiet ./...; \
		echo "$(COLOR_GREEN)✓$(COLOR_RESET) Security checks passed"; \
	else \
		echo "$(COLOR_YELLOW)⚠$(COLOR_RESET) gosec not installed"; \
		echo "  Install with: $(COLOR_YELLOW)go install github.com/securego/gosec/v2/cmd/gosec@latest$(COLOR_RESET)"; \
	fi

# ============================================================================
# Dependencies & Tools
# ============================================================================

.PHONY: deps deps-upgrade deps-verify tools

# Download and tidy dependencies
deps:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Installing Go dependencies..."
	@$(GOMOD) download
	@$(GOMOD) tidy
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Dependencies installed"

# Upgrade dependencies to latest
deps-upgrade:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Upgrading dependencies..."
	@$(GO) get -u ./...
	@$(GOMOD) tidy
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Dependencies upgraded"
	@echo "$(COLOR_YELLOW)ℹ$(COLOR_RESET)  Run $(COLOR_YELLOW)make test$(COLOR_RESET) to verify everything still works"

# Verify dependencies
deps-verify:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Verifying dependencies..."
	@$(GOMOD) verify
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Dependencies verified"

# Install development tools
tools:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Installing development tools..."
	@echo "  Installing swag (Swagger generator)..."
	@$(GO) install github.com/swaggo/swag/cmd/swag@latest
	@echo "  Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	@echo "  Installing gosec (security scanner)..."
	@$(GO) install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Development tools installed"

# ============================================================================
# Workflows
# ============================================================================

.PHONY: fresh fresh-seed ci pre-commit docker-clean verify

# Complete fresh start (no seeding - requires manual server start)
fresh: db-destroy setup build
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Triggering migrations..."
	@timeout 5s $(MAKE) run 2>/dev/null || true
	@sleep 1
	@echo ""
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Fresh environment ready!"
	@echo ""
	@echo "$(COLOR_CYAN)Next steps:$(COLOR_RESET)"
	@echo "  1. $(COLOR_YELLOW)make run$(COLOR_RESET) &      - Start server in background"
	@echo "  2. $(COLOR_YELLOW)make db-seed$(COLOR_RESET)    - Populate with test data"
	@echo ""
	@echo "Or use: $(COLOR_YELLOW)make fresh-seed$(COLOR_RESET) to automate everything"
	@echo ""

# Complete fresh start with automatic seeding (starts server in background)
fresh-seed: db-destroy setup build
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Starting server in background..."
	@$(MAKE) run > /tmp/fits-server.log 2>&1 & echo $$! > /tmp/fits-server.pid
	@echo "$(COLOR_YELLOW)⏳$(COLOR_RESET) Waiting for server to be ready..."
	@for i in 1 2 3 4 5 6 7 8 9 10; do \
		if curl -s http://localhost:8080/health >/dev/null 2>&1; then \
			echo "$(COLOR_GREEN)✓$(COLOR_RESET) Server is ready"; \
			break; \
		fi; \
		if [ $$i -eq 10 ]; then \
			echo "$(COLOR_RED)✗$(COLOR_RESET) Timeout waiting for server"; \
			exit 1; \
		fi; \
		sleep 1; \
	done
	@echo ""
	@$(MAKE) db-seed
	@echo ""
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Fresh environment ready and seeded!"
	@echo ""
	@echo "$(COLOR_CYAN)Server Info:$(COLOR_RESET)"
	@echo "  PID: $$(cat /tmp/fits-server.pid 2>/dev/null || echo 'unknown')"
	@echo "  Logs: tail -f /tmp/fits-server.log"
	@echo "  Stop: kill $$(cat /tmp/fits-server.pid 2>/dev/null)"
	@echo ""
	@echo "$(COLOR_CYAN)Ready to use:$(COLOR_RESET)"
	@echo "  API Docs: $(COLOR_YELLOW)http://localhost:8080/docs$(COLOR_RESET)"
	@echo "  Health: $(COLOR_YELLOW)http://localhost:8080/health$(COLOR_RESET)"
	@echo ""

# CI pipeline
ci: deps check test-coverage
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) CI pipeline complete"

# Pre-commit checks
pre-commit: fmt vet lint test
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Ready to commit"

# Clean Docker resources
docker-clean:
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Cleaning Docker resources..."
	@$(DOCKER_COMPOSE) down -v --remove-orphans
	@$(DOCKER) system prune -f
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) Docker resources cleaned"

# Comprehensive system verification
verify: deps-verify db-status
	@echo "$(COLOR_CYAN)→$(COLOR_RESET) Running comprehensive verification..."
	@echo ""
	@echo "$(COLOR_CYAN)Checking Go environment:$(COLOR_RESET)"
	@$(GO) version
	@echo ""
	@echo "$(COLOR_CYAN)Checking build:$(COLOR_RESET)"
	@$(MAKE) build >/dev/null 2>&1 && echo "$(COLOR_GREEN)✓$(COLOR_RESET) Build successful" || echo "$(COLOR_RED)✗$(COLOR_RESET) Build failed"
	@echo ""
	@echo "$(COLOR_CYAN)Checking database schema:$(COLOR_RESET)"
	@$(MAKE) db-verify
	@echo ""
	@echo "$(COLOR_GREEN)✓$(COLOR_RESET) System verification complete"

# ============================================================================
# End of Makefile
# ============================================================================
