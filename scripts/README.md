# Scripts

This directory contains utility scripts for testing, database management, and development workflows.

## Available Scripts

### Database Management

#### `reset.sh`
Complete reset of the FITS backend system.

**What it does:**
- Stops running server processes
- Drops and recreates PostgreSQL database
- Removes generated RSA keys
- Cleans uploads directory
- Removes test artifacts

**Usage:**
```bash
./scripts/reset.sh
```

**Warning:** This will delete all data! You will be prompted for confirmation.

---

#### `reset_db_only.sh`
Quick database reset - truncates all tables without dropping the database.

**What it does:**
- Truncates all tables (preserves schema)
- Maintains database structure
- Faster than full reset

**Usage:**
```bash
./scripts/reset_db_only.sh
```

**When to use:**
- Need to clear data but keep schema
- Quick testing iterations
- Database migrations are up to date

---

### Testing Scripts

#### `test_auth_flow.sh`
End-to-end test for the complete authentication flow.

**Tests:**
1. Admin initialization (bootstrap)
2. Invitation creation
3. Invitation retrieval
4. User registration (invitation completion)
5. User login
6. Token authentication
7. Token refresh
8. User logout

**Prerequisites:**
- Server must be running: `go run cmd/server/main.go`
- `jq` must be installed for JSON parsing
- `curl` for HTTP requests

**Usage:**
```bash
./scripts/test_auth_flow.sh
```

---

#### `test_full_flow.sh`
Comprehensive end-to-end test covering authentication and user management.

**Tests:**
- Admin bootstrap
- Student invitation and registration
- Teacher invitation and registration
- Login for both roles
- Security (wrong password rejection)
- Token refresh
- Logout for both users

**Prerequisites:**
- Server running
- Clean database (use `reset_db_only.sh` if needed)
- `jq` and `curl` installed

**Usage:**
```bash
./scripts/test_full_flow.sh
```

---

#### `test_student_teacher_crud.sh`
Tests complete CRUD operations for Student and Teacher management.

**Tests:**
- Admin bootstrap
- Student: Create, Read, Update, Delete
- Teacher: Create, Read, Update, Delete
- List operations with pagination
- Student-Teacher assignment
- Duplicate email validation
- Data integrity checks

**Prerequisites:**
- Server running
- Clean database
- `jq` and `curl` installed

**Usage:**
```bash
./scripts/test_student_teacher_crud.sh
```

---

#### `comprehensive_api_test.sh`
Systematic test suite covering all API endpoints with detailed reporting.

**Test Suites:**
1. System Health & Monitoring
2. Authentication & Authorization
3. Invitation System
4. Student Management (CRUD)
5. Teacher Management (CRUD)
6. Digital Signing (experimental)
7. Input Validation & Security
8. Error Handling
9. Database Validation

**Features:**
- Detailed test reporting
- Success/failure statistics
- Known issue tracking
- Security testing (XSS, SQL injection, etc.)
- Automatic cleanup

**Prerequisites:**
- Server running
- `jq` and `curl` installed
- Python 3 (for long string test)

**Usage:**
```bash
./scripts/comprehensive_api_test.sh
```

**Output:**
- Colored test results
- Summary with success rate
- List of critical problems
- Warnings for known issues

---

## General Prerequisites

All scripts require:
- PostgreSQL 15+ running
- Go 1.21+ installed
- `jq` for JSON processing: `sudo apt install jq` (Ubuntu/Debian) or `brew install jq` (macOS)
- `curl` for HTTP requests (usually pre-installed)

## Environment Variables

Scripts respect standard environment variables:

```bash
# Database connection
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=fits_db
export DB_USER=fits_user
export DB_PASSWORD=fits_password

# API endpoint
export API_URL=http://localhost:8080
```

## Testing Workflow

### Recommended Testing Sequence

1. **Start fresh:**
   ```bash
   ./scripts/reset.sh
   ```

2. **Start server:**
   ```bash
   go run cmd/server/main.go
   ```

3. **Run tests (in order):**
   ```bash
   ./scripts/test_auth_flow.sh              # Auth basics
   ./scripts/test_student_teacher_crud.sh   # CRUD operations
   ./scripts/test_full_flow.sh              # Complete flow
   ./scripts/comprehensive_api_test.sh      # Full validation
   ```

### Quick Testing

For quick iterations without full reset:

```bash
./scripts/reset_db_only.sh
./scripts/test_auth_flow.sh
```

## Notes

- **Manual Testing**: These scripts are designed for manual testing and development
- **Automated Testing**: For CI/CD, use Go tests: `go test ./...`
- **Test Data**: Scripts create test data that is cleaned up automatically
- **Idempotency**: Most scripts can be run multiple times safely

## Troubleshooting

### Server Not Running
```
Error: Server is not running on http://localhost:8080
```
**Solution:** Start server with `go run cmd/server/main.go`

### Admin Already Initialized
```
Admin already initialized
```
**Solution:** Either use existing token or run `./scripts/reset_db_only.sh`

### Permission Denied
```bash
chmod +x scripts/*.sh
```

### jq Not Found
```bash
# Ubuntu/Debian
sudo apt install jq

# macOS
brew install jq

# Arch Linux
sudo pacman -S jq
```

## Historical Notes

- `setup_db.sh` - **DEPRECATED** - Database is now automatically initialized by Go application
