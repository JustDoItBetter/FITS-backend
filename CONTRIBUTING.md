# Contributing to FITS Backend

Thank you for your interest in contributing to the FITS (Flexible IT Training System) Backend! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow
- Maintain professionalism in all interactions

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Git
- Docker (optional, for local development)

### Fork and Clone

```bash
# Fork the repository on GitHub first, then:
git clone https://github.com/YOUR_USERNAME/FITS-backend.git
cd FITS-backend
```

## Development Setup

1. **Install Dependencies**
   ```bash
   go mod download
   ```

2. **Configure Environment**
   ```bash
   cp configs/config.toml.example configs/config.toml
   # Edit config.toml with your local settings
   ```

3. **Start PostgreSQL**
   ```bash
   # Using Docker:
   docker-compose up -d postgres

   # Or use your local PostgreSQL instance
   ```

4. **Run Database Migrations**
   ```bash
   # Migrations run automatically on startup
   go run cmd/server/main.go
   ```

5. **Verify Setup**
   ```bash
   # Run tests
   go test ./...

   # Build application
   go build -o bin/fits-server cmd/server/main.go
   ```

## Project Structure

```
FITS-backend/
├── cmd/server/          # Application entry point
├── internal/
│   ├── common/          # Shared utilities (errors, responses, pagination)
│   ├── config/          # Configuration management
│   ├── domain/          # Business logic (DDD pattern)
│   │   ├── auth/        # Authentication domain
│   │   ├── student/     # Student management
│   │   ├── teacher/     # Teacher management
│   │   └── signing/     # Digital signing ( experimental)
│   └── middleware/      # HTTP middleware (JWT, RBAC)
├── pkg/                 # Reusable packages
│   ├── crypto/          # Cryptography utilities
│   ├── database/        # Database connection
│   └── logger/          # Structured logging
├── docs/                # Documentation
└── migrations/          # Database migration reference
```

## Development Workflow

### Creating a New Feature

1. **Create a Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Changes**
   - Follow the coding standards below
   - Write tests for new functionality
   - Update documentation as needed

3. **Test Your Changes**
   ```bash
   go test ./...
   go test -race ./...  # Check for race conditions
   go test -coverprofile=coverage.out ./...
   ```

4. **Commit Changes**
   ```bash
   git add .
   git commit -m "feat: add user profile update endpoint"
   ```

### Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>: <description>

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat: add pagination to student list endpoint
fix: prevent SQL injection in email validation
docs: update API documentation for teacher endpoints
refactor: extract common CRUD logic into generic service
test: add integration tests for auth flow
```

## Coding Standards

### Go Code Style

1. **Follow Standard Go Conventions**
   - Run `gofmt` and `goimports` before committing
   - Use `golangci-lint` for linting

2. **Comments Should Explain "Why", Not "What"**
   ```go
   //  Bad - explains what
   // Create student in database
   studentRepo.Create(ctx, student)

   //  Good - explains why
   // Students must be created before sending invitation to prevent orphaned invites
   studentRepo.Create(ctx, student)
   ```

3. **Error Handling**
   - Always handle errors explicitly
   - Use custom AppError types for API errors
   - Log errors with structured logging
   ```go
   student, err := repo.GetByUUID(ctx, uuid)
   if err != nil {
       logger.Error("Failed to retrieve student",
           zap.String("uuid", uuid),
           zap.Error(err),
       )
       return nil, errors.NotFound("student")
   }
   ```

4. **Input Validation and Sanitization**
   - Always sanitize user input to prevent XSS
   - Use validation package for struct validation
   ```go
   import "github.com/JustDoItBetter/FITS-backend/internal/common/validation"

   func (r *CreateRequest) ToEntity() *Entity {
       return &Entity{
           Name:  validation.SanitizeName(r.Name),
           Email: validation.SanitizeEmail(r.Email),
       }
   }
   ```

### Domain-Driven Design Pattern

Each domain follows a 4-layer architecture:

1. **Model** (`model.go`) - Entities, DTOs, value objects
2. **Repository** (`repository.go`, `repository_gorm.go`) - Data access
3. **Service** (`service.go`) - Business logic
4. **Handler** (`handler.go`) - HTTP endpoints

When adding a new domain:
```bash
internal/domain/newdomain/
├── model.go           # Domain entities and DTOs
├── repository.go      # Interface definition
├── repository_gorm.go # GORM implementation
├── service.go         # Business logic
└── handler.go         # HTTP handlers
```

### REST API Design

- Use standard HTTP methods correctly:
  - `POST` for creation
  - `GET` for retrieval
  - `PUT` for full updates
  - `PATCH` for partial updates
  - `DELETE` for deletion

- Return appropriate status codes:
  - `200 OK` - Success with data
  - `201 Created` - Resource created
  - `204 No Content` - Success without data
  - `400 Bad Request` - Invalid input
  - `401 Unauthorized` - Authentication required
  - `403 Forbidden` - Insufficient permissions
  - `404 Not Found` - Resource doesn't exist
  - `409 Conflict` - Duplicate resource
  - `422 Unprocessable Entity` - Validation errors
  - `500 Internal Server Error` - Server errors

- Use consistent response format:
  ```json
  {
    "success": true,
    "data": { ... }
  }
  ```

### Swagger Documentation

Add Swagger annotations to all handlers:

```go
// Create godoc
// @Summary Create a new student
// @Description Creates a new student record with the provided information
// @Tags Students
// @Accept json
// @Produce json
// @Param request body CreateStudentRequest true "Student creation request"
// @Success 201 {object} response.SuccessResponse{data=Student}
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /api/v1/student [post]
// @Security BearerAuth
func (h *Handler) Create(c *fiber.Ctx) error {
    // Implementation
}
```

## Testing

### Unit Tests

- Place tests in `*_test.go` files
- Use table-driven tests when appropriate
- Mock external dependencies

```go
func TestService_Create(t *testing.T) {
    tests := []struct {
        name    string
        request *CreateRequest
        wantErr bool
    }{
        {
            name: "valid student creation",
            request: &CreateRequest{
                FirstName: "John",
                LastName: "Doe",
                Email: "john@example.com",
            },
            wantErr: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Tests

- Use SQLite in-memory database for integration tests
- Clean up test data after each test
- Test complete workflows

### Test Coverage

Aim for:
- **>70% overall coverage**
- **>80% for critical business logic**
- **100% for security-related code**

Check coverage:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Submitting Changes

### Pull Request Process

1. **Update Documentation**
   - Update README.md if needed
   - Add/update API documentation
   - Update CHANGELOG.md

2. **Run Full Test Suite**
   ```bash
   go test ./...
   go test -race ./...
   golangci-lint run
   ```

3. **Create Pull Request**
   - Use a clear, descriptive title
   - Reference related issues
   - Describe changes and rationale
   - Include screenshots for UI changes
   - List breaking changes if any

4. **PR Description Template**
   ```markdown
   ## Description
   Brief description of changes

   ## Type of Change
   - [ ] Bug fix
   - [ ] New feature
   - [ ] Breaking change
   - [ ] Documentation update

   ## Testing
   - [ ] Unit tests pass
   - [ ] Integration tests pass
   - [ ] Manual testing completed

   ## Checklist
   - [ ] Code follows style guidelines
   - [ ] Self-review completed
   - [ ] Comments added for complex code
   - [ ] Documentation updated
   - [ ] No new warnings
   - [ ] Tests added/updated
   ```

### Review Process

- All PRs require at least one approval
- Address review feedback promptly
- Keep PRs focused and reasonably sized
- Squash commits before merging (if requested)

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Fiber Framework](https://docs.gofiber.io/)
- [GORM Documentation](https://gorm.io/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## Getting Help

- Open an issue for bugs or feature requests
- Ask questions in discussions
- Check existing documentation first

## License

By contributing, you agree that your contributions will be licensed under the project's MIT License.

---

Thank you for contributing to FITS Backend! 
