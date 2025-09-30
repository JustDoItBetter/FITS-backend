package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

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

func MetricsAuth(secret string) fiber.Handler {
	return SecretAuth(secret)
}

func RegistrationAuth(secret string) fiber.Handler {
	return SecretAuth(secret)
}

func DeletionAuth(secret string) fiber.Handler {
	return SecretAuth(secret)
}

func UpdateAuth(secret string) fiber.Handler {
	return SecretAuth(secret)
}
