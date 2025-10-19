package middleware

import (
	"strings"

	"github.com/JustDoItBetter/FITS-backend/pkg/crypto"
	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware creates a middleware that validates JWT tokens
type JWTMiddleware struct {
	jwtService *crypto.JWTService
}

// NewJWTMiddleware creates a new JWT middleware
func NewJWTMiddleware(jwtService *crypto.JWTService) *JWTMiddleware {
	return &JWTMiddleware{
		jwtService: jwtService,
	}
}

// RequireAuth is a middleware that requires valid JWT authentication
func (m *JWTMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "missing authorization header",
			})
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "invalid authorization format, expected: Bearer <token>",
			})
		}

		token := parts[1]

		// Validate token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "invalid or expired token",
				"details": err.Error(),
			})
		}

		// Check token type (should be access or admin token)
		if claims.TokenType != crypto.TokenTypeAccess && claims.TokenType != crypto.TokenTypeAdmin {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "invalid token type",
			})
		}

		// Store claims in context for later use
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)
		c.Locals("token_type", claims.TokenType)

		return c.Next()
	}
}

// OptionalAuth is a middleware that extracts JWT info if present, but doesn't require it
func (m *JWTMiddleware) OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next() // No auth provided, continue
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next() // Invalid format, but don't fail
		}

		token := parts[1]

		// Validate token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			return c.Next() // Invalid token, but don't fail
		}

		// Store claims in context
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)
		c.Locals("token_type", claims.TokenType)

		return c.Next()
	}
}
