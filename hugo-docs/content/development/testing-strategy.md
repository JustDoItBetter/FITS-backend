---
title: "Testing Strategy"
weight: 2
---

# Testing Strategy

FITS Backend follows a comprehensive testing strategy to ensure code quality and reliability.

## Testing Pyramid

```
         ┌─────────────┐
         │  E2E Tests  │ (Few - High-level workflows)
         ├─────────────┤
         │Integration  │ (Some - Database operations)
         │    Tests    │
         ├─────────────┤
         │   Unit      │ (Many - Business logic)
         │   Tests     │
         └─────────────┘
```

## Test Types

### Unit Tests

Test individual functions and methods in isolation.

**Location**: `*_test.go` files alongside source code

**Example**: `internal/domain/auth/auth_service_test.go`

```go
func TestAuthService_Login_Success(t *testing.T) {
    // Arrange
    mockRepo := NewMockRepository()
    jwtService := crypto.NewJWTService("test-secret")
    config := &config.JWTConfig{
        AccessTokenExpiry:  15 * time.Minute,
        RefreshTokenExpiry: 7 * 24 * time.Hour,
    }
    service := auth.NewAuthService(mockRepo, jwtService, config)
    
    ctx := context.Background()
    email := "test@example.com"
    password := "TestPass123!"
    
    // Mock user
    hashedPassword, _ := crypto.HashPassword(password)
    mockUser := &auth.User{
        ID:           1,
        Email:        email,
        PasswordHash: hashedPassword,
        Role:         "student",
        IsActive:     true,
    }
    
    mockRepo.On("FindByEmail", ctx, email).Return(mockUser, nil)
    mockRepo.On("CreateSession", ctx, mock.Anything).Return(nil)
    
    // Act
    result, err := service.Login(ctx, email, password)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.NotEmpty(t, result.Token)
    assert.NotEmpty(t, result.RefreshToken)
    assert.Equal(t, email, result.User.Email)
    
    mockRepo.AssertExpectations(t)
}
```

**What to Test:**
- Business logic
- Input validation
- Error handling
- Edge cases
- Boundary conditions

**What NOT to Test:**
- Third-party libraries
- Framework code
- Simple getters/setters

### Integration Tests

Test interactions with external systems (database, APIs).

**Location**: `*_test.go` files, typically in repository implementations

**Example**: `internal/domain/student/repository_gorm_test.go`

```go
func TestStudentRepository_Create(t *testing.T) {
    // Setup test database
    db, err := setupTestDB()
    require.NoError(t, err)
    defer db.Close()
    
    repo := student.NewGormRepository(db.DB)
    ctx := context.Background()
    
    // Create student
    s := &student.Student{
        FirstName: "John",
        LastName:  "Doe",
        Email:     "john@example.com",
        StudentID: "S001",
    }
    
    err = repo.Create(ctx, s)
    assert.NoError(t, err)
    assert.NotZero(t, s.ID)
    
    // Verify in database
    found, err := repo.FindByID(ctx, s.ID)
    assert.NoError(t, err)
    assert.Equal(t, s.FirstName, found.FirstName)
    assert.Equal(t, s.Email, found.Email)
}

func setupTestDB() (*database.DB, error) {
    config := &config.DatabaseConfig{
        Driver:   "sqlite",
        Database: ":memory:",
    }
    return database.New(config)
}
```

**What to Test:**
- Database CRUD operations
- Transaction handling
- Foreign key constraints
- Data integrity

### Table-Driven Tests

Use table-driven tests for testing multiple scenarios:

```go
func TestPasswordValidation(t *testing.T) {
    tests := []struct {
        name     string
        password string
        wantErr  bool
        errType  string
    }{
        {
            name:     "valid password",
            password: "StrongPass123!",
            wantErr:  false,
        },
        {
            name:     "too short",
            password: "Short1!",
            wantErr:  true,
            errType:  "too_short",
        },
        {
            name:     "no uppercase",
            password: "weakpass123!",
            wantErr:  true,
            errType:  "no_uppercase",
        },
        {
            name:     "no lowercase",
            password: "WEAKPASS123!",
            wantErr:  true,
            errType:  "no_lowercase",
        },
        {
            name:     "no number",
            password: "WeakPassword!",
            wantErr:  true,
            errType:  "no_number",
        },
        {
            name:     "no special char",
            password: "WeakPassword123",
            wantErr:  true,
            errType:  "no_special",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validation.ValidatePassword(tt.password)
            
            if tt.wantErr {
                assert.Error(t, err)
                // Check error type if needed
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## Testing Tools

### Testing Frameworks

- **testing**: Standard Go testing package
- **testify**: Assertions and mocking
  ```go
  import (
      "github.com/stretchr/testify/assert"
      "github.com/stretchr/testify/mock"
      "github.com/stretchr/testify/require"
  )
  ```

### Assertions

```go
// Basic assertions
assert.Equal(t, expected, actual)
assert.NotEqual(t, unexpected, actual)
assert.Nil(t, obj)
assert.NotNil(t, obj)
assert.True(t, condition)
assert.False(t, condition)

// Error assertions
assert.NoError(t, err)
assert.Error(t, err)
assert.ErrorIs(t, err, targetErr)
assert.ErrorContains(t, err, "substring")

// Collection assertions
assert.Contains(t, slice, element)
assert.Len(t, collection, length)
assert.Empty(t, collection)

// Require (fails immediately if condition not met)
require.NoError(t, err) // Stops test execution
```

### Mocking

Use testify/mock for creating mocks:

```go
// Define mock
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

// Use in test
mockRepo := new(MockRepository)
mockRepo.On("FindByEmail", ctx, "test@example.com").Return(user, nil)

// Verify expectations
mockRepo.AssertExpectations(t)
mockRepo.AssertCalled(t, "FindByEmail", ctx, "test@example.com")
mockRepo.AssertNumberOfCalls(t, "FindByEmail", 1)
```

## Running Tests

### All Tests

```bash
go test ./...
```

### Specific Package

```bash
go test ./internal/domain/auth/...
go test ./pkg/crypto/...
```

### Specific Test

```bash
go test -run TestAuthService_Login ./internal/domain/auth/
```

### Verbose Output

```bash
go test -v ./...
```

### With Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage by package
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
```

### Race Detection

```bash
go test -race ./...
```

### Parallel Execution

```bash
# Run tests in parallel (default GOMAXPROCS)
go test -parallel 4 ./...
```

### Short Mode

Skip slow tests:

```go
func TestSlowOperation(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping slow test in short mode")
    }
    // Slow test code...
}
```

Run:
```bash
go test -short ./...
```

## Test Organization

### File Naming

- Test files: `*_test.go`
- Same directory as code being tested
- Example: `service.go` → `service_test.go`

### Test Naming

```go
// Pattern: Test{FunctionName}_{Scenario}_{ExpectedResult}

func TestAuthService_Login_ValidCredentials_ReturnsToken(t *testing.T) {}
func TestAuthService_Login_InvalidPassword_ReturnsError(t *testing.T) {}
func TestAuthService_Login_InactiveUser_ReturnsError(t *testing.T) {}
```

### Test Structure (AAA Pattern)

```go
func TestSomething(t *testing.T) {
    // Arrange - Set up test data and dependencies
    input := "test data"
    expected := "expected result"
    
    // Act - Execute the code under test
    result := SomeFunction(input)
    
    // Assert - Verify the results
    assert.Equal(t, expected, result)
}
```

## Coverage Targets

### Current Coverage

Check current coverage:
```bash
go test -cover ./...
```

### Coverage Goals

- **Overall**: 70%+
- **Business Logic** (services): 90%+
- **Repositories**: 80%+
- **Handlers**: 70%+
- **Utilities**: 85%+

### Coverage by Package

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
```

## Testing Best Practices

### 1. Write Tests First (TDD)

```go
// 1. Write failing test
func TestNewFeature(t *testing.T) {
    result := NewFeature()
    assert.Equal(t, "expected", result)
}

// 2. Implement feature to make test pass
func NewFeature() string {
    return "expected"
}

// 3. Refactor while keeping tests green
```

### 2. Test One Thing Per Test

```go
// Bad - tests multiple things
func TestUserService(t *testing.T) {
    user := service.Create(...)
    assert.NotNil(t, user)
    
    updated := service.Update(...)
    assert.NotNil(t, updated)
    
    err := service.Delete(...)
    assert.NoError(t, err)
}

// Good - separate tests
func TestUserService_Create(t *testing.T) {
    user := service.Create(...)
    assert.NotNil(t, user)
}

func TestUserService_Update(t *testing.T) {
    updated := service.Update(...)
    assert.NotNil(t, updated)
}

func TestUserService_Delete(t *testing.T) {
    err := service.Delete(...)
    assert.NoError(t, err)
}
```

### 3. Use Meaningful Test Names

```go
// Bad
func TestLogin(t *testing.T) {}

// Good
func TestAuthService_Login_ValidCredentials_ReturnsToken(t *testing.T) {}
func TestAuthService_Login_InvalidPassword_ReturnsError(t *testing.T) {}
```

### 4. Don't Test Implementation Details

```go
// Bad - tests internal implementation
func TestCache_Set_CallsRedis(t *testing.T) {
    cache.Set("key", "value")
    assert.True(t, redisClient.WasCalled())
}

// Good - tests behavior
func TestCache_Set_ValueIsRetrievable(t *testing.T) {
    cache.Set("key", "value")
    result, err := cache.Get("key")
    assert.NoError(t, err)
    assert.Equal(t, "value", result)
}
```

### 5. Use Test Fixtures

```go
// test_fixtures.go
func createTestUser() *User {
    return &User{
        Email:    "test@example.com",
        Role:     "student",
        IsActive: true,
    }
}

func createTestDB(t *testing.T) *database.DB {
    db, err := database.New(&config.DatabaseConfig{
        Driver:   "sqlite",
        Database: ":memory:",
    })
    require.NoError(t, err)
    return db
}
```

### 6. Clean Up After Tests

```go
func TestWithDatabase(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()  // Clean up
    
    // Test code...
}

func TestWithTempFile(t *testing.T) {
    file, err := os.CreateTemp("", "test-*")
    require.NoError(t, err)
    defer os.Remove(file.Name())  // Clean up
    
    // Test code...
}
```

### 7. Test Error Cases

```go
func TestService_Create_DuplicateEmail_ReturnsError(t *testing.T) {
    service := NewService(repo)
    
    // First create should succeed
    _, err := service.Create(ctx, user)
    assert.NoError(t, err)
    
    // Duplicate should fail
    _, err = service.Create(ctx, user)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "duplicate")
}
```

## Continuous Integration

Tests run automatically on:
- Pull requests
- Commits to main branch
- Scheduled nightly builds

### GitHub Actions Example

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_DB: fits_test
          POSTGRES_USER: fits
          POSTGRES_PASSWORD: test
        ports:
          - 5432:5432
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.25'
      
      - name: Install dependencies
        run: go mod download
      
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
```

## Next Steps

- [Contribution Guidelines](/development/contribution-guidelines/)
- [Local Development Setup](/development/local-setup/)
- [API Reference](/api/endpoints/)
