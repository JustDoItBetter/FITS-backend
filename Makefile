# FITS Backend - Professional Makefile
# Simplified development and deployment tasks

.PHONY: help run test test-verbose test-coverage clean reset reset-db build install-deps docker-up docker-down e2e-test fmt lint bench quickstart bootstrap dev ci test-all swagger

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=fits-server
BINARY_PATH=bin/$(BINARY_NAME)

# Default target
help:
	@echo "FITS Backend - Available Commands"
	@echo ""
	@echo "Development:"
	@echo "  make run              - Start development server"
	@echo "  make dev              - Start complete dev environment (Docker + Server)"
	@echo "  make quickstart       - Initial setup (deps + Docker)"
	@echo "  make fmt              - Format Go code"
	@echo "  make lint             - Run code linting"
	@echo ""
	@echo "Building:"
	@echo "  make build            - Compile binary"
	@echo "  make clean            - Remove build artifacts"
	@echo "  make swagger          - Regenerate Swagger documentation"
	@echo ""
	@echo "Testing:"
	@echo "  make test             - Run all unit tests"
	@echo "  make test-verbose     - Run tests with verbose output"
	@echo "  make test-coverage    - Run tests with coverage report"
	@echo "  make test-all         - Run unit and E2E tests"
	@echo "  make e2e-test         - Run end-to-end tests (requires running server)"
	@echo "  make bench            - Run benchmarks"
	@echo ""
	@echo "Database:"
	@echo "  make reset-db         - Reset database only"
	@echo "  make reset            - Complete reset (DB + Keys + Uploads)"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up        - Start Docker containers"
	@echo "  make docker-down      - Stop Docker containers"
	@echo ""
	@echo "Utilities:"
	@echo "  make install-deps     - Install/update Go dependencies"
	@echo "  make bootstrap        - Initialize admin account (requires running server)"
	@echo "  make ci               - Run CI pipeline (deps + tests + coverage)"
	@echo ""

# Server operations
run:
	@echo "[INFO] Starting FITS Backend Server..."
	$(GOCMD) run cmd/server/main.go

dev: docker-up
	@echo "[INFO] Starting development environment..."
	@echo "[INFO] PostgreSQL is running. Starting server..."
	@$(MAKE) run

# Build operations
build:
	@echo "[BUILD] Compiling binary..."
	@mkdir -p bin
	$(GOBUILD) -o $(BINARY_PATH) cmd/server/main.go
	@echo "[SUCCESS] Binary created: $(BINARY_PATH)"

clean:
	@echo "[CLEAN] Removing build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f *.test
	@echo "[SUCCESS] Cleanup complete"

swagger:
	@echo "[BUILD] Regenerating Swagger documentation..."
	@if [ -f ~/go/bin/swag ]; then \
		~/go/bin/swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal; \
		echo "[SUCCESS] Swagger docs regenerated"; \
	else \
		echo "[ERROR] swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest"; \
		exit 1; \
	fi

# Testing
test:
	@echo "[TEST] Running unit tests..."
	$(GOTEST) ./...

test-verbose:
	@echo "[TEST] Running unit tests (verbose)..."
	$(GOTEST) ./... -v

test-coverage:
	@echo "[TEST] Running unit tests with coverage..."
	$(GOTEST) ./... -coverprofile=coverage.out -covermode=atomic
	@echo ""
	@echo "[REPORT] Coverage summary:"
	$(GOCMD) tool cover -func=coverage.out | tail -20
	@echo ""
	@echo "[INFO] Generate HTML report: go tool cover -html=coverage.out"

test-all: test
	@echo ""
	@echo "[TEST] Starting E2E tests..."
	@if curl -s http://localhost:8080/health > /dev/null 2>&1; then \
		$(MAKE) e2e-test; \
	else \
		echo "[WARN] Server not running - E2E tests skipped"; \
		echo "[INFO] Start server with 'make run' then run 'make e2e-test'"; \
	fi

e2e-test:
	@echo "[TEST] Running End-to-End tests..."
	@if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then \
		echo "[ERROR] Server not running. Start with 'make run' first."; \
		exit 1; \
	fi
	@./scripts/test_full_flow.sh

bench:
	@echo "[BENCH] Running benchmarks..."
	$(GOTEST) ./pkg/crypto -bench=. -benchmem

# Code quality
fmt:
	@echo "[FORMAT] Formatting code..."
	$(GOCMD) fmt ./...
	@echo "[SUCCESS] Code formatted"

lint:
	@echo "[LINT] Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "[WARN] golangci-lint not installed"; \
		echo "[INFO] Installation: https://golangci-lint.run/usage/install/"; \
	fi

# Database operations
reset-db:
	@echo "[DB] Resetting database..."
	@./scripts/reset_db_only.sh

reset:
	@echo "[RESET] Complete reset (DB, Keys, Uploads)..."
	@./scripts/reset.sh

# Docker operations
docker-up:
	@echo "[DOCKER] Starting containers..."
	docker-compose up -d
	@echo "[SUCCESS] Containers started"
	@echo "[INFO] PostgreSQL running on localhost:5432"

docker-down:
	@echo "[DOCKER] Stopping containers..."
	docker-compose down
	@echo "[SUCCESS] Containers stopped"

# Dependencies
install-deps:
	@echo "[DEPS] Installing Go dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "[SUCCESS] Dependencies installed"

# Quick setup
quickstart: install-deps docker-up
	@echo "[WAIT] Waiting 5 seconds for PostgreSQL..."
	@sleep 5
	@echo ""
	@echo "[SUCCESS] Setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. make run         - Start server"
	@echo "  2. make bootstrap   - Initialize admin account"
	@echo ""

# Bootstrap admin (requires running server)
bootstrap:
	@echo "[BOOTSTRAP] Initializing admin account..."
	@curl -s -X POST http://localhost:8080/api/v1/bootstrap/init \
		-H "Content-Type: application/json" \
		-d '{"username":"admin","password":"SecurePassword123!","email":"admin@fits.example.com"}' | jq .

# CI/CD pipeline
ci: install-deps test-coverage
	@echo "[SUCCESS] CI pipeline completed successfully"
