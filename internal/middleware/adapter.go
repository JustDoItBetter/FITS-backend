package middleware

import (
	"github.com/JustDoItBetter/FITS-backend/pkg/crypto"
	"github.com/gofiber/fiber/v2"
)

// MiddlewareAdapter wraps JWT and RBAC middleware for dependency injection into handlers
// This allows handlers to register routes with their security requirements in one place
type MiddlewareAdapter struct {
	jwtMiddleware *JWTMiddleware
}

// NewMiddlewareAdapter creates a new middleware adapter
func NewMiddlewareAdapter(jwtMiddleware *JWTMiddleware) *MiddlewareAdapter {
	return &MiddlewareAdapter{
		jwtMiddleware: jwtMiddleware,
	}
}

// RequireAuth returns the JWT authentication middleware
func (m *MiddlewareAdapter) RequireAuth() fiber.Handler {
	return m.jwtMiddleware.RequireAuth()
}

// OptionalAuth returns the optional JWT authentication middleware
func (m *MiddlewareAdapter) OptionalAuth() fiber.Handler {
	return m.jwtMiddleware.OptionalAuth()
}

// RequireAdmin returns middleware that requires admin role
func (m *MiddlewareAdapter) RequireAdmin() fiber.Handler {
	return RequireAdmin()
}

// RequireRole returns middleware that requires one of the specified roles
func (m *MiddlewareAdapter) RequireRole(roles ...interface{}) fiber.Handler {
	// Convert interface{} to crypto.Role
	cryptoRoles := make([]crypto.Role, len(roles))
	for i, role := range roles {
		switch v := role.(type) {
		case crypto.Role:
			cryptoRoles[i] = v
		case string:
			cryptoRoles[i] = crypto.Role(v)
		default:
			// Fallback to admin if type is unexpected
			cryptoRoles[i] = crypto.RoleAdmin
		}
	}
	return RequireRole(cryptoRoles...)
}
