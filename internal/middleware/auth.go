package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// SecretAuth creates a middleware that validates Bearer token authentication using a shared secret.
func SecretAuth(expectedSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization format, expected: Bearer <token>",
			})
		}

		token := parts[1]

		if token != expectedSecret {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid secret",
			})
		}

		return c.Next()
	}
}

// MetricsAuth creates a middleware that validates authentication for metrics endpoints.
func MetricsAuth(secret string) fiber.Handler {
	return SecretAuth(secret)
}

// RegistrationAuth creates a middleware that validates authentication for registration endpoints.
func RegistrationAuth(secret string) fiber.Handler {
	return SecretAuth(secret)
}

// DeletionAuth creates a middleware that validates authentication for deletion endpoints.
func DeletionAuth(secret string) fiber.Handler {
	return SecretAuth(secret)
}

// UpdateAuth creates a middleware that validates authentication for update endpoints.
func UpdateAuth(secret string) fiber.Handler {
	return SecretAuth(secret)
}
