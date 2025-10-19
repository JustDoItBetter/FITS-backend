package middleware

import (
	"github.com/JustDoItBetter/FITS-backend/internal/common/errors"
	"github.com/JustDoItBetter/FITS-backend/internal/common/response"
	"github.com/JustDoItBetter/FITS-backend/pkg/crypto"
	"github.com/gofiber/fiber/v2"
)

// RequireRole creates a middleware that checks if the user has one of the specified roles
// Must be used after RequireAuth middleware to ensure role context is set
func RequireRole(allowedRoles ...crypto.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleRaw := c.Locals("role")
		if roleRaw == nil {
			return response.Error(c, errors.Unauthorized("authentication required"))
		}

		// Safe type assertion prevents panic if context value has unexpected type
		userRole, ok := roleRaw.(crypto.Role)
		if !ok {
			return response.Error(c, errors.Internal("invalid role type in context"))
		}

		// Check if user has one of the allowed roles
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				return c.Next()
			}
		}

		// User doesn't have any of the required roles
		return response.Error(c, errors.NewAppError(
			fiber.StatusForbidden,
			"Forbidden",
			"this action requires one of the following roles: "+rolesToString(allowedRoles),
		))
	}
}

// RequireAdmin is a convenience middleware that requires admin role
func RequireAdmin() fiber.Handler {
	return RequireRole(crypto.RoleAdmin)
}

// RequireTeacher is a convenience middleware that requires teacher role
// (Admin can also access teacher endpoints)
func RequireTeacher() fiber.Handler {
	return RequireRole(crypto.RoleAdmin, crypto.RoleTeacher)
}

// RequireStudent is a convenience middleware that requires student role
// (Teachers and admins can also access student endpoints for supervision)
func RequireStudent() fiber.Handler {
	return RequireRole(crypto.RoleAdmin, crypto.RoleTeacher, crypto.RoleStudent)
}

// RequireOwnership creates a middleware that checks if the user owns the resource
// Must be used after RequireAuth middleware to ensure user context exists
// The resource UUID should be in the route parameter specified by paramName
// Admins bypass ownership checks to enable administrative access
func RequireOwnership(paramName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIDRaw := c.Locals("user_id")
		if userIDRaw == nil {
			return response.Error(c, errors.Unauthorized("authentication required"))
		}

		// Safe type assertion prevents runtime panic if context has wrong type
		userID, ok := userIDRaw.(string)
		if !ok {
			return response.Error(c, errors.Internal("invalid user ID type in context"))
		}

		roleRaw := c.Locals("role")
		if roleRaw == nil {
			return response.Error(c, errors.Unauthorized("authentication required"))
		}

		// Safe type assertion for role
		userRole, ok := roleRaw.(crypto.Role)
		if !ok {
			return response.Error(c, errors.Internal("invalid role type in context"))
		}

		// Admins bypass ownership checks for administrative operations
		if userRole == crypto.RoleAdmin {
			return c.Next()
		}

		resourceUUID := c.Params(paramName)
		if resourceUUID == "" {
			return response.Error(c, errors.BadRequest("resource UUID not provided"))
		}

		// Simple ownership check: user_uuid from JWT must match resource UUID
		// NOTE: For complex resources, query database to verify user_uuid field
		if userID != resourceUUID {
			return response.Error(c, errors.NewAppError(
				fiber.StatusForbidden,
				"Forbidden",
				"you don't have permission to access this resource",
			))
		}

		return c.Next()
	}
}

// Helper function to convert roles to string
func rolesToString(roles []crypto.Role) string {
	if len(roles) == 0 {
		return ""
	}

	result := string(roles[0])
	for i := 1; i < len(roles); i++ {
		result += ", " + string(roles[i])
	}
	return result
}
