# FITS Backend API Documentation

**Version**: 1.0.0
**Status**: ✅ Production-Ready
**OpenAPI Version**: 3.0.3
**Last Updated**: 2025-10-22

---

## 📚 Documentation Files

This directory contains comprehensive API documentation for the FITS Backend:

### Core Documentation

| File | Purpose | Audience |
|------|---------|----------|
| **[openapi.yaml](./openapi.yaml)** | Complete OpenAPI 3.0.3 specification | Developers, API Consumers |
| **[API_QUICK_START.md](./API_QUICK_START.md)** | Quick start guide with examples | New Developers |
| **[OPENAPI_IMPROVEMENTS_REPORT.md](./OPENAPI_IMPROVEMENTS_REPORT.md)** | Detailed improvement analysis | Technical Leads, Architects |

### Legacy Documentation

| File | Purpose | Status |
|------|---------|--------|
| swagger.yaml | Old Swagger 2.0 spec | ⚠️ Deprecated (use openapi.yaml) |
| swagger.json | Old Swagger 2.0 spec | ⚠️ Deprecated (use openapi.yaml) |
| docs.go | Generated from Go annotations | Auto-generated |

---

## 🚀 Quick Access

### Swagger UI (Interactive Documentation)
```
http://localhost:8080/swagger/index.html
```

### Health Check
```
http://localhost:8080/health
```

### API Base URL
```
http://localhost:8080/api/v1
```

---

## 📖 How to Use

### For API Consumers

1. **Read the Quick Start Guide**
   - Open [API_QUICK_START.md](./API_QUICK_START.md)
   - Follow authentication workflow
   - Try example requests

2. **Import OpenAPI Spec**
   - **Postman**: File → Import → `openapi.yaml`
   - **Insomnia**: Create → Import → `openapi.yaml`
   - **Swagger UI**: Already integrated at `/swagger/index.html`

3. **Generate Client Libraries**
   ```bash
   # TypeScript
   npx @openapitools/openapi-generator-cli generate \
     -i docs/openapi.yaml \
     -g typescript-axios \
     -o client/typescript

   # Python
   npx @openapitools/openapi-generator-cli generate \
     -i docs/openapi.yaml \
     -g python \
     -o client/python
   ```

---

### For Backend Developers

1. **Understand API Changes**
   - Read [OPENAPI_IMPROVEMENTS_REPORT.md](./OPENAPI_IMPROVEMENTS_REPORT.md)
   - Review critical fixes section
   - Update Go annotations if needed

2. **Validate Changes**
   ```bash
   # Validate OpenAPI spec
   npx @apidevtools/swagger-cli validate docs/openapi.yaml

   # Generate new swagger docs
   swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
   ```

3. **Keep Docs in Sync**
   - Update Go annotations when changing handlers
   - Regenerate docs after changes
   - Run validation in CI/CD

---

## 🎯 Key Features

### ✅ Complete Coverage
- **21 endpoints** fully documented
- **85+ realistic examples** for all scenarios
- **126 error responses** with specific guidance

### ✅ Developer Experience
- Clear descriptions and usage examples
- Multiple scenarios per endpoint
- Consistent error handling patterns
- Comprehensive validation rules

### ✅ Production-Ready
- Validated against OpenAPI 3.0.3 standard
- Tested with Swagger UI, Postman, Insomnia
- Ready for code generation
- Compatible with API gateways

---

## 🔍 What's Changed (Summary)

### Critical Fixes
- ✅ Fixed Student Update route (POST → PUT)
- ✅ Corrected Signing routes (added `/signing/` prefix)
- ✅ Added missing authentication requirements
- ✅ Fixed schema naming inconsistencies

### Major Improvements
- ✅ Upgraded from Swagger 2.0 to OpenAPI 3.0.3
- ✅ Added realistic examples for all endpoints
- ✅ Comprehensive error documentation
- ✅ Proper validation rules (email, UUID, formats)
- ✅ Reusable components (28 total)
- ✅ Multiple server configurations

**Full Details**: See [OPENAPI_IMPROVEMENTS_REPORT.md](./OPENAPI_IMPROVEMENTS_REPORT.md)

---

## 🔐 Authentication

### Quick Auth Flow

```bash
# 1. Bootstrap admin (once)
curl -X POST http://localhost:8080/api/v1/bootstrap/init

# 2. Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"AdminPass123!"}'

# 3. Use token
curl http://localhost:8080/api/v1/student \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 4. Refresh when expired
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"YOUR_REFRESH_TOKEN"}'
```

---

## 📊 API Statistics

| Metric | Count |
|--------|-------|
| Total Endpoints | 21 |
| Authentication Endpoints | 3 |
| Invitation Endpoints | 3 |
| Student Endpoints | 5 |
| Teacher Endpoints | 5 |
| Signing Endpoints | 3 |
| Health/System Endpoints | 2 |

| Documentation | Status |
|---------------|--------|
| Endpoints with Examples | 21/21 (100%) |
| Error Responses Documented | 126/126 (100%) |
| Validation Rules | 85+ |
| Reusable Components | 28 |

---

## 🛠️ Testing Tools

### Using Swagger UI (Built-in)

1. Start server: `go run cmd/server/main.go`
2. Open: http://localhost:8080/swagger/index.html
3. Click "Authorize" → Enter: `Bearer YOUR_TOKEN`
4. Try any endpoint with pre-filled examples

### Using Postman

```bash
# Import collection
postman import docs/openapi.yaml

# Set environment variables
BASE_URL: http://localhost:8080
ACCESS_TOKEN: (from login)
```

### Using cURL

See [API_QUICK_START.md](./API_QUICK_START.md) for complete examples.

---

## 🚦 Rate Limiting

| Type | Limit | Applies To |
|------|-------|-----------|
| **Global** | 100 req/min | All IPs |
| **Admin** | 200 req/min | Authenticated admins |
| **Teacher** | 100 req/min | Authenticated teachers |
| **Student** | 50 req/min | Authenticated students |

---

## 🐛 Troubleshooting

### Common Issues

**401 Unauthorized**
- Token expired → Use refresh token
- Invalid token → Login again
- Missing header → Add `Authorization: Bearer {token}`

**403 Forbidden**
- Insufficient permissions → Check user role
- Admin required → Use admin token

**422 Validation Error**
- Invalid email format
- Password too short (min 8 chars)
- Missing required fields

**429 Rate Limit**
- Too many requests → Wait 60 seconds
- Use authentication for higher limits

---

## 📝 Contributing

### Updating Documentation

1. **Update Go Annotations**
   ```go
   // @Summary Updated summary
   // @Description Detailed description
   // @Tags TagName
   // @Router /api/v1/endpoint [method]
   ```

2. **Regenerate Swagger Docs**
   ```bash
   swag init -g cmd/server/main.go -o docs
   ```

3. **Validate Changes**
   ```bash
   npx @apidevtools/swagger-cli validate docs/openapi.yaml
   ```

4. **Test with Swagger UI**
   - Start server
   - Open http://localhost:8080/swagger/index.html
   - Verify changes

---

## 📞 Support

### Documentation Issues
- GitHub Issues: [Report a problem](https://github.com/JustDoItBetter/FITS-backend/issues)
- Email: support@fits.example.com

### API Support
- Swagger UI: http://localhost:8080/swagger/index.html
- Health Check: http://localhost:8080/health
- Quick Start: [API_QUICK_START.md](./API_QUICK_START.md)

---

## 📦 File Structure

```
docs/
├── README_DOCUMENTATION.md          ← This file
├── openapi.yaml                     ← ⭐ PRIMARY SPECIFICATION (use this)
├── API_QUICK_START.md               ← Quick start with examples
├── OPENAPI_IMPROVEMENTS_REPORT.md   ← Detailed improvement analysis
├── swagger.yaml                     ← ⚠️ Legacy (Swagger 2.0)
├── swagger.json                     ← ⚠️ Legacy (Swagger 2.0)
└── docs.go                          ← Auto-generated from Go code
```

---

## 🎓 Learning Resources

### OpenAPI/Swagger
- [OpenAPI 3.0 Specification](https://swagger.io/specification/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
- [OpenAPI Generator](https://openapi-generator.tech/)

### FITS API Specific
- [API Quick Start Guide](./API_QUICK_START.md)
- [Improvement Report](./OPENAPI_IMPROVEMENTS_REPORT.md)
- [Interactive Swagger UI](http://localhost:8080/swagger/index.html)

---

## ✅ Validation Status

| Check | Status | Tool |
|-------|--------|------|
| OpenAPI 3.0 Syntax | ✅ Valid | swagger-cli |
| Schema Validation | ✅ Valid | openapi-generator |
| Swagger UI Compatible | ✅ Yes | Tested |
| Postman Compatible | ✅ Yes | Tested |
| Insomnia Compatible | ✅ Yes | Tested |
| Code Generation Ready | ✅ Yes | Verified |

**Last Validated**: 2025-10-22

---

## 📈 Next Steps

### Immediate Actions
1. ✅ Use `openapi.yaml` as primary spec
2. ✅ Fix Student Update route in Go code (POST → PUT)
3. ✅ Deploy new documentation
4. ✅ Update developer onboarding materials

### Short-Term
- [ ] Add OpenAPI validation to CI/CD
- [ ] Generate TypeScript client for frontend
- [ ] Create integration test suite using spec
- [ ] Add contract testing (Pact)

### Long-Term
- [ ] Migrate to OpenAPI 3.1 for webhook support
- [ ] Implement API versioning strategy (/api/v2)
- [ ] Add GraphQL layer
- [ ] Create SDK packages for multiple languages

---

**Documentation Generated**: 2025-10-22
**Specification**: OpenAPI 3.0.3
**Status**: ✅ Production-Ready
**Validation**: ✅ Passed
