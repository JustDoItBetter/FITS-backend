---
title: "Contribution Guidelines"
weight: 3
---

# Contribution Guidelines

Thank you for considering contributing to FITS Backend! This document outlines the process and standards for contributing.

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on constructive feedback
- Assume good intentions

## Getting Started

### 1. Find an Issue

- Browse [open issues](https://github.com/JustDoItBetter/FITS-backend/issues)
- Look for issues labeled `good-first-issue` or `help-wanted`
- Ask questions if anything is unclear

### 2. Discuss Major Changes

For significant changes:
1. Open an issue first to discuss the approach
2. Wait for maintainer feedback
3. Get approval before starting implementation

### 3. Fork and Clone

```bash
# Fork on GitHub, then clone your fork
git clone https://github.com/YOUR_USERNAME/FITS-backend.git
cd FITS-backend

# Add upstream remote
git remote add upstream https://github.com/JustDoItBetter/FITS-backend.git
```

### 4. Create a Branch

```bash
git checkout -b feature/your-feature-name
```

Branch naming:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test additions/updates

## Development Process

### 1. Make Your Changes

- Follow the [code style guidelines](#code-style)
- Write tests for new code
- Update documentation as needed
- Keep commits focused and atomic

### 2. Write Tests

All new code requires tests:
- Unit tests for business logic
- Integration tests for database operations
- Table-driven tests for multiple scenarios

```go
func TestNewFeature(t *testing.T) {
    // Arrange
    input := "test"
    
    // Act
    result := NewFeature(input)
    
    // Assert
    assert.Equal(t, expected, result)
}
```

### 3. Run Quality Checks

Before committing:

```bash
# Format code
make fmt

# Run linters
make lint

# Run tests
make test

# Check coverage
make test-cover
```

### 4. Commit Changes

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `style:` Code style changes (formatting)
- `refactor:` Code refactoring
- `test:` Test additions/changes
- `chore:` Build process, dependencies

**Examples:**

```bash
# Feature
git commit -m "feat(auth): add password reset functionality"

# Bug fix
git commit -m "fix(student): correct validation for student ID format"

# Documentation
git commit -m "docs(api): update authentication flow diagram"

# Refactoring
git commit -m "refactor(repo): simplify database query methods"
```

**Good commit message:**
```
feat(auth): add OAuth2 authentication support

- Implement OAuth2 flow with Google and GitHub providers
- Add configuration for OAuth clients
- Update authentication middleware to support OAuth tokens
- Add tests for OAuth authentication flow

Closes #123
```

### 5. Push Changes

```bash
# Push to your fork
git push origin feature/your-feature-name
```

### 6. Create Pull Request

1. Go to GitHub and create a PR from your fork
2. Fill out the PR template completely
3. Link related issues
4. Wait for review

## Pull Request Guidelines

### PR Title

Follow conventional commit format:
```
feat(domain): add user profile management
fix(auth): resolve token expiration issue
docs: update installation guide
```

### PR Description

Include:
- **Summary**: What does this PR do?
- **Motivation**: Why is this change needed?
- **Changes**: List of main changes
- **Testing**: How was this tested?
- **Screenshots**: If UI changes (if applicable)
- **Checklist**: Completed items

**Template:**

```markdown
## Summary
Brief description of changes.

## Motivation
Why is this change necessary?

## Changes
- Change 1
- Change 2
- Change 3

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing performed

## Checklist
- [ ] Code follows style guidelines
- [ ] Tests pass locally
- [ ] Documentation updated
- [ ] No breaking changes (or documented if present)
```

### Review Process

1. **Automated Checks**: CI must pass (tests, linters)
2. **Code Review**: At least one maintainer review required
3. **Changes Requested**: Address feedback
4. **Approval**: Maintainer approves
5. **Merge**: Maintainer merges PR

## Code Style

### Go Style Guidelines

Follow [Effective Go](https://golang.org/doc/effective_go.html) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

#### Formatting

```bash
# Format all code
go fmt ./...

# Or use make
make fmt
```

#### Naming Conventions

```go
// Variables: camelCase
var userID uint
var studentEmail string

// Constants: CamelCase or UPPER_CASE
const MaxRetries = 3
const API_VERSION = "v1"

// Exported: Start with uppercase
func HandleRequest() {}
type User struct {}

// Unexported: Start with lowercase
func validateEmail() {}
type session struct {}

// Interfaces: -er suffix when possible
type Reader interface {}
type Validator interface {}
```

#### Package Organization

```go
// Standard library first
import (
    "context"
    "fmt"
    "time"
)

// Then third-party packages
import (
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
)

// Finally internal packages
import (
    "github.com/JustDoItBetter/FITS-backend/internal/domain/auth"
    "github.com/JustDoItBetter/FITS-backend/pkg/logger"
)
```

#### Error Handling

```go
// Return errors, don't panic
func DoSomething() error {
    if err != nil {
        return fmt.Errorf("failed to do something: %w", err)
    }
    return nil
}

// Check errors immediately
result, err := SomeFunction()
if err != nil {
    return err
}

// Use custom error types when appropriate
var ErrNotFound = errors.New("resource not found")
```

#### Comments

```go
// Exported functions need doc comments
// CreateUser creates a new user in the database.
// Returns the created user or an error if creation fails.
func CreateUser(ctx context.Context, user *User) (*User, error) {
    // Implementation...
}

// Explain why, not what
// Use pagination to avoid loading all records into memory
results := repo.FindWithPagination(page, limit)
```

#### Testing

```go
// Test naming: Test{Function}_{Scenario}_{Expected}
func TestAuthService_Login_ValidCredentials_ReturnsToken(t *testing.T) {
    // Arrange
    // Act
    // Assert
}

// Use testify assertions
assert.Equal(t, expected, actual)
assert.NoError(t, err)
require.NotNil(t, result)
```

### API Documentation

All API endpoints require Swagger annotations:

```go
// @Summary Create a new student
// @Description Creates a student record with the provided information
// @Tags students
// @Accept json
// @Produce json
// @Param student body CreateStudentRequest true "Student information"
// @Success 201 {object} StudentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/student [post]
// @Security BearerAuth
func (h *Handler) CreateStudent(c *fiber.Ctx) error {
    // Implementation...
}
```

Regenerate docs after changes:
```bash
make docs
```

## Testing Requirements

### Minimum Coverage

- Overall: 70%
- Business logic (services): 90%
- New code: 80%

### Required Tests

- Unit tests for all business logic
- Integration tests for database operations
- Table-driven tests for multiple scenarios
- Error case testing

### Running Tests

```bash
# All tests
make test

# With coverage
make test-cover

# Race detection
make test-race
```

## Documentation

### Code Documentation

- Add GoDoc comments for exported functions/types
- Explain complex logic with inline comments
- Update package documentation

### User Documentation

Update relevant documentation:
- README if setup changes
- API docs for endpoint changes
- Architecture docs for design changes
- This documentation for contribution process changes

## Common Pitfalls

### 1. Not Running Tests Locally

Always run tests before pushing:
```bash
make test
```

### 2. Not Formatting Code

Format code before committing:
```bash
make fmt
```

### 3. Large Pull Requests

Keep PRs focused:
- One feature/fix per PR
- Aim for < 400 lines changed
- Break large features into multiple PRs

### 4. Not Updating Documentation

Update docs alongside code changes.

### 5. Breaking Changes

Avoid breaking changes when possible. If necessary:
- Document clearly in PR
- Update CHANGELOG
- Provide migration guide

## Review Feedback

### Addressing Feedback

1. **Read carefully**: Understand the feedback
2. **Ask questions**: If unclear, ask for clarification
3. **Make changes**: Address all comments
4. **Respond**: Mark conversations as resolved
5. **Update**: Push changes and notify reviewers

### Disagreements

If you disagree with feedback:
1. Explain your reasoning politely
2. Be open to discussion
3. Maintainers have final say
4. Focus on what's best for the project

## Recognition

Contributors are recognized in:
- CONTRIBUTORS.md file
- GitHub contributors page
- Release notes

## Getting Help

- **Questions**: Open a [discussion](https://github.com/JustDoItBetter/FITS-backend/discussions)
- **Bug reports**: Open an [issue](https://github.com/JustDoItBetter/FITS-backend/issues)
- **Security issues**: Email security@fits.example.com (do not open public issues)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Thank You!

Your contributions make FITS Backend better for everyone. Thank you for taking the time to contribute!
