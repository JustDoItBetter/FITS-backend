# Security Guide

## Secrets Management

### Development

For local development, you can use the default secrets in `configs/config.toml`. These are safe for local testing only.

### Production

**CRITICAL: Change ALL secrets before deploying to production!**

#### 1. Generate Strong Secrets

```bash
# JWT Secret (64 characters)
openssl rand -base64 64

# API Secrets (32 characters each)
openssl rand -base64 32
openssl rand -base64 32
openssl rand -base64 32
openssl rand -base64 32
```

#### 2. Update Configuration

Copy `configs/config.example.toml` to `configs/config.toml` and update:

- `[jwt].secret` - JWT signing key (64+ characters)
- `[database].password` - Database password
- `[secrets].metrics_secret` - Metrics endpoint protection
- `[secrets].registration_secret` - Registration API protection
- `[secrets].deletion_secret` - Deletion API protection
- `[secrets].update_secret` - Update API protection

#### 3. Environment Variables (Recommended)

For production, use environment variables instead of config files:

```bash
export FITS_JWT_SECRET="your-generated-secret"
export FITS_DB_PASSWORD="your-db-password"
export FITS_METRICS_SECRET="your-metrics-secret"
```

Update `internal/config/config.go` to read from environment variables if available.

## RSA Keys

### Generate Admin Keys

```bash
# Create keys directory
mkdir -p configs/keys

# Generate RSA 4096-bit key pair
openssl genrsa -out configs/keys/admin.key 4096
openssl rsa -in configs/keys/admin.key -pubout -out configs/keys/admin.pub

# Set correct permissions
chmod 600 configs/keys/admin.key
chmod 644 configs/keys/admin.pub
```

**IMPORTANT:** Never commit private keys (.key files) to Git!

## CORS Configuration

### Development
```toml
allowed_origins = "*"
```

### Production
```toml
allowed_origins = "https://yourdomain.com,https://app.yourdomain.com"
```

## Database Security

### Production Settings

```toml
[database]
ssl_mode = "require"  # Enable SSL
password = "use-strong-password-64-chars"
```

### Connection Pool

- `max_conns = 25` - Maximum connections
- `min_conns = 5` - Minimum connections

Adjust based on your server capacity.

## Rate Limiting

```toml
rate_limit = 100  # requests per minute per IP
```

For production APIs:
- Public endpoints: 100 req/min
- Authenticated: 500 req/min
- Admin: 1000 req/min

Adjust in code if needed.

## Security Headers

Already configured in `cmd/server/main.go`:

- ✅ X-XSS-Protection
- ✅ Content-Type-Nosniff
- ✅ X-Frame-Options (SAMEORIGIN)
- ✅ Content-Security-Policy
- ✅ Referrer-Policy

## Password Requirements

Enforced in `pkg/crypto/password.go`:

- Minimum 12 characters
- At least 1 uppercase letter
- At least 1 lowercase letter
- At least 1 number
- At least 1 special character
- Maximum 100 characters

## JWT Token Security

- **Access Token:** 1 hour expiry
- **Refresh Token:** 30 days expiry
- **Invitation Token:** 7 days expiry

Tokens are signed with HS256 and validated on every protected request.

## Input Validation

All user inputs are validated:

- ✅ Email format validation
- ✅ UUID format validation
- ✅ SQL injection prevention (GORM parameterized queries)
- ✅ XSS prevention (HTML escaping)
- ✅ File upload size limits

## Security Checklist

Before deploying to production:

- [ ] Change JWT secret
- [ ] Change database password
- [ ] Change all API secrets
- [ ] Generate RSA keys
- [ ] Enable database SSL
- [ ] Set CORS to specific domains
- [ ] Review rate limits
- [ ] Enable HTTPS (reverse proxy)
- [ ] Set secure cookies (if using)
- [ ] Configure firewall rules
- [ ] Enable audit logging
- [ ] Set up monitoring alerts

## Vulnerability Reporting

If you discover a security vulnerability, please email: security@fits.example.com

Do not create public GitHub issues for security vulnerabilities.

## Security Updates

- Check dependencies regularly: `go list -m all`
- Update dependencies: `go get -u ./...`
- Monitor security advisories: https://github.com/advisories

## Resources

- OWASP Top 10: https://owasp.org/www-project-top-ten/
- Go Security: https://golang.org/doc/security
- NIST Guidelines: https://nvd.nist.gov/
