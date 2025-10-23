# Swagger UI Diagnosis Report

**Date**: 2025-10-22
**Status**: ✅ DIAGNOSIS COMPLETE
**Severity**: 🔴 CRITICAL - Swagger UI completely non-functional

---

## Executive Summary

Swagger UI is **completely non-functional** due to incorrect route paths in Go Swagger annotations. All signing endpoints return **404 Not Found** and the student update endpoint uses the **wrong HTTP method**.

**Impact**: 100% of Swagger UI "Try it out" attempts fail for signing endpoints, student updates fail due to wrong HTTP method.

---

## Root Causes Identified

### 🔴 CRITICAL ISSUE #1: Wrong Signing Endpoint Paths

**Severity**: CRITICAL
**Impact**: All 3 signing endpoints unreachable from Swagger UI
**Root Cause**: Swagger annotations don't match actual router configuration

**Evidence:**
```bash
# Swagger spec says (WRONG):
POST /api/v1/upload                → 404 Not Found
GET  /api/v1/sign_requests         → 404 Not Found
POST /api/v1/sign_uploads          → 404 Not Found

# Server actually expects (CORRECT):
POST /api/v1/signing/upload        → 401 Unauthorized (route exists, needs auth)
GET  /api/v1/signing/sign_requests → 401 Unauthorized (route exists, needs auth)
POST /api/v1/signing/sign_uploads  → 401 Unauthorized (route exists, needs auth)
```

**Code Analysis:**

File: `cmd/server/main.go:258-260`
```go
signingGroup := api.Group("/signing")  // Creates /api/v1/signing
signingGroup.Use(jwtMiddleware.RequireAuth())
signingHandler.RegisterRoutes(signingGroup)
```

File: `internal/domain/signing/handler.go:21-24`
```go
func (h *Handler) RegisterRoutes(router fiber.Router) {
    router.Post("/upload", h.Upload)           // → /api/v1/signing/upload
    router.Get("/sign_requests", h.GetSignRequests)  // → /api/v1/signing/sign_requests
    router.Post("/sign_uploads", h.SignUploads)      // → /api/v1/signing/sign_uploads
}
```

But Swagger annotations say:
```go
// @Router /api/v1/upload [post]           ❌ WRONG - missing "/signing"
// @Router /api/v1/sign_requests [get]     ❌ WRONG - missing "/signing"
// @Router /api/v1/sign_uploads [post]     ❌ WRONG - missing "/signing"
```

---

### 🔴 CRITICAL ISSUE #2: Student Update Wrong HTTP Method

**Severity**: CRITICAL
**Impact**: Student update fails from Swagger UI
**Root Cause**: Swagger annotation says POST, code uses PUT

**Evidence:**

File: `internal/domain/student/handler.go:34-36`
```go
router.Put("/:uuid",  // ✅ CORRECT - uses PUT
    jwtMW.RequireAuth(),
    rbacMW.RequireAdmin(),
    h.Update,
)
```

But Swagger annotation says:
```go
// @Router /api/v1/student/{uuid} [post]  ❌ WRONG - should be [put]
```

**Test Results:**
```bash
# Swagger UI will try:
POST /api/v1/student/{uuid} → 404 Not Found (no POST handler)

# Server expects:
PUT /api/v1/student/{uuid} → Works correctly
```

---

### ⚠️ ISSUE #3: Missing Security Documentation

**Severity**: HIGH
**Impact**: Swagger UI shows signing endpoints as public when they require auth
**Root Cause**: @Security annotations missing

**Evidence:**
- Signing endpoints require JWT auth (code shows `jwtMiddleware.RequireAuth()`)
- Swagger annotations don't include `@Security BearerAuth`
- Users try endpoints without tokens → get confusing errors

---

### ⚠️ ISSUE #4: Swagger 2.0 vs OpenAPI 3.0 Confusion

**Severity**: MEDIUM
**Impact**: Two conflicting specs exist
**Root Cause**: Project has both auto-generated Swagger 2.0 and manual OpenAPI 3.0.3

**Files:**
- `docs/swagger.{json,yaml}` - Auto-generated Swagger 2.0 (served by app, has bugs)
- `docs/openapi.yaml` - Manual OpenAPI 3.0.3 (correct but not served)

**Current State:**
- Swagger UI uses `docs/swagger.json` (broken)
- `docs/openapi.yaml` is correct but ignored

---

## Test Results

### Endpoint Status

| Endpoint | Swagger Path | Actual Path | HTTP Method | Status |
|----------|--------------|-------------|-------------|--------|
| **Signing - Upload** | `/api/v1/upload` | `/api/v1/signing/upload` | POST | 🔴 404 |
| **Signing - Get Requests** | `/api/v1/sign_requests` | `/api/v1/signing/sign_requests` | GET | 🔴 404 |
| **Signing - Upload Signed** | `/api/v1/sign_uploads` | `/api/v1/signing/sign_uploads` | POST | 🔴 404 |
| **Student - Update** | `/api/v1/student/{uuid}` POST | `/api/v1/student/{uuid}` PUT | PUT | 🔴 Wrong Method |
| **Student - List** | `/api/v1/student` | `/api/v1/student` | GET | ✅ Works |
| **Auth - Login** | `/api/v1/auth/login` | `/api/v1/auth/login` | POST | ✅ Works |
| **Health Check** | `/health` | `/health` | GET | ✅ Works |

**Success Rate**: 3/7 endpoints work correctly (43%)
**Failure Rate**: 4/7 endpoints broken (57%)

---

## Verification Commands

```bash
# Start server
go run cmd/server/main.go

# Test wrong paths (what Swagger UI tries):
curl -X POST http://localhost:8080/api/v1/upload  # 404
curl -X GET http://localhost:8080/api/v1/sign_requests  # 404
curl -X POST http://localhost:8080/api/v1/sign_uploads  # 404

# Test correct paths (what server expects):
curl -X POST http://localhost:8080/api/v1/signing/upload \
  -H "Authorization: Bearer TOKEN"  # 401 (route exists, needs valid token)

# Test student update method mismatch:
curl -X POST http://localhost:8080/api/v1/student/UUID  # 404 (no POST handler)
curl -X PUT http://localhost:8080/api/v1/student/UUID \
  -H "Authorization: Bearer TOKEN"  # Works
```

---

## Priority Fixes

### Priority 1: Fix Signing Routes (CRITICAL - Quick Win)
**Time**: 2 minutes
**Risk**: None
**Impact**: Fixes 3/4 broken endpoints

**Fix**: Update Swagger annotations in `internal/domain/signing/handler.go`
```diff
- // @Router /api/v1/upload [post]
+ // @Router /api/v1/signing/upload [post]

- // @Router /api/v1/sign_requests [get]
+ // @Router /api/v1/signing/sign_requests [get]

- // @Router /api/v1/sign_uploads [post]
+ // @Router /api/v1/signing/sign_uploads [post]
```

### Priority 2: Fix Student Update Method (CRITICAL - Quick Win)
**Time**: 1 minute
**Risk**: None
**Impact**: Fixes student update endpoint

**Fix**: Update Swagger annotation in `internal/domain/student/handler.go:135`
```diff
- // @Router /api/v1/student/{uuid} [post]
+ // @Router /api/v1/student/{uuid} [put]
```

### Priority 3: Add Security Annotations (HIGH)
**Time**: 5 minutes
**Risk**: None
**Impact**: Documents auth requirements correctly

**Fix**: Add `@Security BearerAuth` to signing handler annotations

---

## Long-Term Improvements

1. **Add CI validation**: Run `swag validate` in CI to catch annotation errors
2. **Add integration tests**: Test Swagger spec matches actual routes
3. **Consider OpenAPI 3.0**: Migrate from Swagger 2.0 to OpenAPI 3.0
4. **Add spec validation**: Run `swagger-cli validate` on generated docs

---

## Next Steps

1. ✅ **Apply Priority 1 & 2 fixes** (3 minutes total)
2. ✅ **Regenerate Swagger**: Run `make swagger`
3. ✅ **Test in Swagger UI**: Verify all endpoints work
4. ✅ **Add tests**: Prevent regression
5. ✅ **Document**: Update API docs

---

**Estimated Total Fix Time**: 15 minutes
**Estimated Test Time**: 10 minutes
**Total Time to Resolution**: 25 minutes

