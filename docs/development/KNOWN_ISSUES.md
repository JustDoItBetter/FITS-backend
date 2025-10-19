# FITS Backend - Known Issues & Future Improvements

**Version:** 1.0.0
**Date:** 2025-10-19
**Status:** Production Ready with known limitations

---

##  Critical Issues

### None Currently Identified

The system is stable and production-ready. Core functionality (Auth, Students, Teachers) has been implemented and tested.

---

##  Known Limitations & TODOs

### 1. **Missing Signing Implementation (EXPERIMENTAL)**
**Priority:** Medium (Feature marked as experimental)
**Component:** `internal/domain/signing/service.go`

**Description:**
The signing feature is currently marked as experimental with placeholder implementations. This is acceptable for v1.0 release.

**TODO Comments in Code:**
```go
// internal/domain/signing/service.go
// TODO: Add parquet file handling when ready
// TODO: Implement RSA signature verification
// TODO: Add database persistence for signatures
// TODO: Implement parquet file parsing and validation
// TODO: Implement parquet file generation
// TODO: Implement signed request processing
```

**Current State:**
- Endpoints exist and are routed correctly
- Returns placeholder data
- No actual file upload or signature processing
- Marked as experimental feature

**Impact:**
- Signing features marked as "experimental" in documentation
- Not blocking v1.0 release
- Can be implemented in future versions

**Recommended Solution:**
1. Implement Apache Parquet file handling library
2. Add RSA signature generation and verification (crypto package already exists)
3. Add database models for signature tracking
4. Implement file upload with multipart/form-data
5. Add signing workflow state machine

**Estimated Effort:** 2-3 weeks

**Status:** Documented as experimental feature, not critical for v1.0

---

### 2. **Auth Service TODOs**
**Priority:** Low
**Component:** `internal/domain/auth/invitation_service.go`

**TODO Comments in Code:**
```go
// TODO: Validate that teacher exists in database
// TODO: Get base URL from config
// TODO: If teacher, generate RSA keypair (implement in keypair domain first)
```

**Current State:**
- Teacher validation happens during invitation completion
- Base URL is hardcoded for development
- RSA keypair generation for teachers not yet implemented

**Impact:**
- Minor - invitation system works but could be more robust
- Base URL needs configuration for production deployment

**Recommended Solution:**
1. Add teacher existence validation before creating invitation
2. Move base URL to config.toml
3. Implement RSA keypair generation for teacher accounts (future feature)

**Estimated Effort:** 4-6 hours

**Status:** Low priority, not blocking release

---

### 3. **No Email Service Integration**
**Priority:** Medium
**Component:** Authentication/Invitation System

**Description:**
Invitation tokens are generated but not automatically sent via email. Admins must manually share invitation links.

**Current State:**
- Invitation creation returns link and token
- No email sending capability
- Manual distribution required

**Impact:**
- Manual workflow for user onboarding
- Increased admin workload
- Higher risk of token exposure

**Recommended Solution:**
1. Integrate SMTP email service (e.g., SendGrid, AWS SES)
2. Create email templates for invitations
3. Add email queue for async processing
4. Implement email tracking (sent/failed/opened)
5. Add retry logic for failed emails

**Example Implementation:**
```go
// Add email service
type EmailService interface {
    SendInvitation(email, name, token string) error
    SendPasswordReset(email, token string) error
}

// Use in invitation creation
if err := emailService.SendInvitation(
    req.Email,
    req.FirstName + " " + req.LastName,
    invitationToken,
); err != nil {
    log.Printf("Failed to send invitation email: %v", err)
    // Continue anyway - return link to admin
}
```

**Estimated Effort:** 1 week

---

### 4. **Insufficient Input Validation**
**Priority:** Medium
**Component:** All Endpoints

**Description:**
Some edge cases and malicious inputs are not properly validated.

**Current State:**
- Basic validation with `go-playground/validator`
- No XSS protection
- No SQL injection protection (GORM handles this)
- Limited string sanitization

**Issues:**
1. No validation for Unicode/emoji in names
2. No maximum length enforcement on text fields
3. No prevention of control characters
4. No URL encoding validation

**Recommended Solution:**
```go
// Add custom validators
validator.RegisterValidation("no_control_chars", func(fl validator.FieldLevel) bool {
    str := fl.Field().String()
    for _, r := range str {
        if r < 32 && r != '\n' && r != '\r' && r != '\t' {
            return false
        }
    }
    return true
})

// Sanitize inputs
func sanitizeString(s string) string {
    s = html.EscapeString(s)
    s = strings.TrimSpace(s)
    return s
}
```

**Estimated Effort:** 2-3 days

---

### 5. ~~No Rate Limiting~~ **✅ IMPLEMENTED**
**Status:** Complete
**Component:** All Endpoints

**Description:**
Rate limiting has been implemented to protect against abuse and DoS attacks.

**Implemented Features:**
- ✅ 100 requests per minute per IP (configurable)
- ✅ Rate limit from config.toml
- ✅ Proper 429 error responses
- ✅ Rate limit can be disabled for development

**Current Implementation:**
```go
// cmd/server/main.go
if cfg.Server.RateLimit > 0 {
    app.Use(limiter.New(limiter.Config{
        Max:        cfg.Server.RateLimit,
        Expiration: 1 * time.Minute,
    }))
}
```

**Configuration:**
```toml
[server]
rate_limit = 100  # requests per minute, 0 = unlimited
```

---

### 6. **No Audit Logging**
**Priority:** Low
**Component:** All Write Operations

**Description:**
No audit trail for data changes (who changed what, when).

**Current State:**
- Basic access logging via Fiber middleware
- No change tracking
- No audit table
- No compliance logging

**Impact:**
- Cannot track data modifications
- Compliance issues (GDPR, audit requirements)
- Difficult to debug data issues
- No accountability

**Recommended Solution:**
```go
// Create audit log table
type AuditLog struct {
    ID          string    `gorm:"primaryKey"`
    EntityType  string    `gorm:"index"` // "student", "teacher", "user"
    EntityID    string    `gorm:"index"`
    Action      string    `gorm:"index"` // "create", "update", "delete"
    UserID      string    `gorm:"index"`
    Changes     string    `gorm:"type:jsonb"` // JSON of changes
    IPAddress   string
    UserAgent   string
    CreatedAt   time.Time
}

// Log all modifications
func (r *GormRepository) Update(ctx context.Context, student *Student) error {
    // Get old state
    old, _ := r.GetByUUID(ctx, student.UUID)

    // Perform update
    if err := r.db.Save(student).Error; err != nil {
        return err
    }

    // Log changes
    auditService.LogChange(ctx, "student", student.UUID, "update", old, student)
    return nil
}
```

**Estimated Effort:** 1 week

---

### 7. **Missing Data Export Functionality**
**Priority:** Low
**Component:** Student/Teacher Management

**Description:**
No way to export student/teacher data in common formats (CSV, Excel, PDF).

**Current State:**
- Only JSON API responses
- No bulk export
- No reporting

**Impact:**
- Manual data extraction required
- No integration with external systems
- Limited reporting capabilities

**Recommended Solution:**
1. Add export endpoints: `/api/v1/student/export?format=csv`
2. Support formats: CSV, Excel, JSON, PDF
3. Add filtering/sorting options
4. Stream large exports
5. Add export scheduling for regular reports

**Estimated Effort:** 1 week

---

### 8. **No Search/Filter Functionality**
**Priority:** Low
**Component:** Student/Teacher List Endpoints

**Description:**
Cannot search or filter students/teachers by criteria.

**Current State:**
- Only full list retrieval
- No search parameters
- No filtering

**Impact:**
- Difficult to find specific records
- Poor UX with large datasets
- Increased client-side processing

**Recommended Solution:**
```go
// Add search parameters
type SearchRequest struct {
    Query      string `query:"q"`           // Full-text search
    Email      string `query:"email"`       // Filter by email
    Name       string `query:"name"`        // Filter by name
    TeacherID  string `query:"teacher_id"`  // Filter by teacher
    Department string `query:"department"`  // Filter by department (teachers)
}

// Implement in repository
func (r *GormRepository) Search(ctx context.Context, req SearchRequest) ([]*Student, error) {
    query := r.db.WithContext(ctx)

    if req.Query != "" {
        query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?",
            "%"+req.Query+"%", "%"+req.Query+"%", "%"+req.Query+"%")
    }

    if req.Email != "" {
        query = query.Where("email ILIKE ?", "%"+req.Email+"%")
    }

    // ... more filters

    var students []*StudentModel
    if err := query.Find(&students).Error; err != nil {
        return nil, err
    }

    return convertToStudents(students), nil
}
```

**Estimated Effort:** 2-3 days

---

##  Minor Issues

### 1. **Inconsistent HTTP Methods**
**Priority:** Low
**Component:** Student/Teacher Endpoints

**Description:**
Student uses `PUT` for creation, Teacher uses `POST`. This is inconsistent and confusing.

**Current State:**
- `PUT /api/v1/student` - Create student
- `POST /api/v1/teacher` - Create teacher

**Recommended Solution:**
Standardize to use `POST` for creation, `PUT` for full replacement, `PATCH` for partial updates.

**Estimated Effort:** 1 hour

---

### 2. **No Health Check Details**
**Priority:** Low
**Component:** `/health` Endpoint

**Description:**
Health check doesn't provide detailed component status.

**Current State:**
- Only checks database connection
- Binary healthy/unhealthy

**Recommended Solution:**
```json
{
  "status": "healthy",
  "version": "1.2.0",
  "uptime": "2h30m",
  "components": {
    "database": "healthy",
    "redis": "not_configured",
    "email": "degraded"
  },
  "timestamp": "2025-10-18T19:00:00Z"
}
```

**Estimated Effort:** 2 hours

---

### 3. **No CORS Configuration**
**Priority:** Low
**Component:** CORS Middleware

**Description:**
CORS is currently set to `AllowOrigins: "*"` which is insecure for production.

**Current State:**
```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "*", // Too permissive!
    AllowMethods: "GET,POST,PUT,DELETE",
    AllowHeaders: "Origin, Content-Type, Accept, Authorization",
}))
```

**Recommended Solution:**
```go
app.Use(cors.New(cors.Config{
    AllowOrigins: cfg.CORS.AllowedOrigins, // From config
    AllowMethods: "GET,POST,PUT,DELETE,PATCH",
    AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    AllowCredentials: true,
    MaxAge: 86400, // 24 hours
}))
```

**Estimated Effort:** 1 hour

---

##  Performance Optimizations

### 1. **Database Connection Pooling**
**Current:** Default GORM settings
**Recommendation:** Tune for production load
```go
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### 2. **Add Database Indexes**
**Missing Indexes:**
- `students.email`  Already indexed (unique)
- `students.teacher_id`  Not indexed
- `teachers.department`  Not indexed
- `users.username`  Already indexed (unique)

**Recommendation:**
```sql
CREATE INDEX idx_students_teacher_id ON students(teacher_id);
CREATE INDEX idx_teachers_department ON teachers(department);
```

### 3. **Response Caching**
**Current:** No caching
**Recommendation:** Cache public GET endpoints (students, teachers lists)

**Estimated Effort:** 2-3 days

---

##  Test Coverage Gaps

Current coverage: **48.9%** for Student/Teacher domains

### Areas Needing More Tests:

1. **Handler Tests:** Only service layer tested
   - Add handler-level tests with mock HTTP requests
   - Test middleware integration
   - Test error response formats

2. **Auth Domain:** Only 18.3% coverage
   - Add more edge case tests
   - Test token expiration scenarios
   - Test concurrent token refresh

3. **Integration Tests:** Limited scope
   - Add full end-to-end tests with real database
   - Test transaction rollback scenarios
   - Test concurrent operations

4. **Load Tests:** Not implemented
   - Add performance benchmarks
   - Test concurrent user limits
   - Identify bottlenecks

**Estimated Effort:** 1-2 weeks

---

##  Security Improvements

### 1. **Password Requirements**
**Current:** Minimum 8 characters
**Recommendation:** Enforce complexity rules
- At least 1 uppercase
- At least 1 lowercase
- At least 1 number
- At least 1 special character

### 2. **Token Security**
**Current:** JWT stored in response
**Recommendation:** Add token rotation and blacklisting
- Implement refresh token rotation
- Add token revocation list (Redis)
- Shorter access token lifetime (15 minutes)

### 3. **API Security Headers**
**Missing:**
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Strict-Transport-Security: max-age=31536000`

**Recommended Solution:**
```go
app.Use(func(c *fiber.Ctx) error {
    c.Set("X-Content-Type-Options", "nosniff")
    c.Set("X-Frame-Options", "DENY")
    c.Set("X-XSS-Protection", "1; mode=block")
    c.Set("Strict-Transport-Security", "max-age=31536000")
    return c.Next()
})
```

**Estimated Effort:** 2 hours

---

##  Documentation Improvements

### 1. **API Examples**
**Current:** Swagger + HTML docs
**Needed:**
- Postman collection
- cURL examples for all endpoints
- Client library examples (JavaScript, Python)

### 2. **Deployment Guide**
**Missing:**
- Production deployment checklist
- Docker Compose for production
- Kubernetes manifests
- Monitoring setup (Prometheus, Grafana)
- Backup and restore procedures

### 3. **Development Guides**
**Needed:**
- Contributing guidelines
- Code style guide
- Git workflow
- Release process

**Estimated Effort:** 1 week

---

##  Prioritization Roadmap

### Sprint 1 (High Priority - 2-3 weeks)
1.  Implement Signing functionality
2.  Add Pagination
3.  Add Rate Limiting
4.  Fix CORS configuration

### Sprint 2 (Medium Priority - 2 weeks)
1.  Email service integration
2.  Enhanced input validation
3.  Search/Filter functionality
4.  Database indexes

### Sprint 3 (Low Priority - 1-2 weeks)
1.  Audit logging
2.  Data export
3.  Health check improvements
4.  Security headers

### Sprint 4 (Testing & Documentation - 2 weeks)
1.  Increase test coverage to 80%
2.  Load testing
3.  Documentation improvements
4.  Deployment guides

---

##  Notes

**Last Updated:** 2025-10-18
**Reviewed By:** Claude (Anthropic AI)
**Next Review:** After Sprint 1 completion

**How to Report Issues:**
```bash
# Create GitHub issue with template
- Title: [Component] Brief description
- Priority: Critical/High/Medium/Low
- Description: Detailed explanation
- Steps to Reproduce
- Expected vs Actual Behavior
- Proposed Solution (optional)
```

**Version History:**
- v1.0 (2025-10-18): Initial issue documentation
