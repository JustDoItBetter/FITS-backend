package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JustDoItBetter/FITS-backend/pkg/crypto"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestApp() *fiber.App {
	return fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
}

func TestNewJWTMiddleware(t *testing.T) {
	t.Run("creates middleware successfully", func(t *testing.T) {
		jwtService := crypto.NewJWTService("test-secret")
		middleware := NewJWTMiddleware(jwtService)

		assert.NotNil(t, middleware)
		assert.Equal(t, jwtService, middleware.jwtService)
	})
}

func TestJWTMiddleware_RequireAuth(t *testing.T) {
	jwtService := crypto.NewJWTService("test-secret")
	middleware := NewJWTMiddleware(jwtService)

	t.Run("allows request with valid token", func(t *testing.T) {
		app := setupTestApp()

		// Generate valid token
		token, err := jwtService.GenerateToken("user-123", crypto.RoleStudent, crypto.TokenTypeAccess, time.Hour)
		require.NoError(t, err)

		// Setup route with middleware
		app.Get("/protected", middleware.RequireAuth(), func(c *fiber.Ctx) error {
			// Check if user info was set in context
			userID := c.Locals("user_id")
			role := c.Locals("role")

			return c.JSON(fiber.Map{
				"user_id": userID,
				"role":    role,
			})
		})

		// Make request
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("rejects request without token", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/protected", middleware.RequireAuth(), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("rejects request with invalid token format", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/protected", middleware.RequireAuth(), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat token")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("rejects request with malformed token", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/protected", middleware.RequireAuth(), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid.jwt.token")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("rejects request with expired token", func(t *testing.T) {
		app := setupTestApp()

		// Generate expired token
		token, err := jwtService.GenerateToken("user-123", crypto.RoleStudent, crypto.TokenTypeAccess, -time.Hour)
		require.NoError(t, err)

		app.Get("/protected", middleware.RequireAuth(), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("rejects request with refresh token", func(t *testing.T) {
		app := setupTestApp()

		// Generate refresh token (not access token)
		token, err := jwtService.GenerateToken("user-123", crypto.RoleStudent, crypto.TokenTypeRefresh, time.Hour)
		require.NoError(t, err)

		app.Get("/protected", middleware.RequireAuth(), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("sets user_id and role in context", func(t *testing.T) {
		app := setupTestApp()

		userID := "test-user-456"
		role := crypto.RoleTeacher
		token, err := jwtService.GenerateToken(userID, role, crypto.TokenTypeAccess, time.Hour)
		require.NoError(t, err)

		var capturedUserID string
		var capturedRole crypto.Role

		app.Get("/protected", middleware.RequireAuth(), func(c *fiber.Ctx) error {
			capturedUserID = c.Locals("user_id").(string)
			capturedRole = c.Locals("role").(crypto.Role)
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, userID, capturedUserID)
		assert.Equal(t, role, capturedRole)
	})

	t.Run("allows admin token", func(t *testing.T) {
		app := setupTestApp()

		token, err := jwtService.GenerateToken("admin-001", crypto.RoleAdmin, crypto.TokenTypeAdmin, time.Hour)
		require.NoError(t, err)

		app.Get("/protected", middleware.RequireAuth(), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestJWTMiddleware_OptionalAuth(t *testing.T) {
	jwtService := crypto.NewJWTService("test-secret")
	middleware := NewJWTMiddleware(jwtService)

	t.Run("allows request with valid token", func(t *testing.T) {
		app := setupTestApp()

		token, err := jwtService.GenerateToken("user-123", crypto.RoleStudent, crypto.TokenTypeAccess, time.Hour)
		require.NoError(t, err)

		var hasUserID bool
		app.Get("/optional", middleware.OptionalAuth(), func(c *fiber.Ctx) error {
			hasUserID = c.Locals("user_id") != nil
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/optional", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, hasUserID, "user_id should be set with valid token")
	})

	t.Run("allows request without token", func(t *testing.T) {
		app := setupTestApp()

		var hasUserID bool
		app.Get("/optional", middleware.OptionalAuth(), func(c *fiber.Ctx) error {
			hasUserID = c.Locals("user_id") != nil
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/optional", nil)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.False(t, hasUserID, "user_id should not be set without token")
	})

	t.Run("ignores invalid token and proceeds", func(t *testing.T) {
		app := setupTestApp()

		var hasUserID bool
		app.Get("/optional", middleware.OptionalAuth(), func(c *fiber.Ctx) error {
			hasUserID = c.Locals("user_id") != nil
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/optional", nil)
		req.Header.Set("Authorization", "Bearer invalid.token")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.False(t, hasUserID, "user_id should not be set with invalid token")
	})

	t.Run("ignores expired token and proceeds", func(t *testing.T) {
		app := setupTestApp()

		token, err := jwtService.GenerateToken("user-123", crypto.RoleStudent, crypto.TokenTypeAccess, -time.Hour)
		require.NoError(t, err)

		var hasUserID bool
		app.Get("/optional", middleware.OptionalAuth(), func(c *fiber.Ctx) error {
			hasUserID = c.Locals("user_id") != nil
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/optional", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.False(t, hasUserID)
	})
}

// Note: extractToken is a private function, tested indirectly through RequireAuth and OptionalAuth

func TestJWTMiddleware_MultipleRoles(t *testing.T) {
	jwtService := crypto.NewJWTService("test-secret")
	middleware := NewJWTMiddleware(jwtService)

	roles := []crypto.Role{
		crypto.RoleStudent,
		crypto.RoleTeacher,
		crypto.RoleAdmin,
	}

	for _, role := range roles {
		t.Run("allows_"+string(role), func(t *testing.T) {
			app := setupTestApp()

			token, err := jwtService.GenerateToken("user-123", role, crypto.TokenTypeAccess, time.Hour)
			require.NoError(t, err)

			app.Get("/protected", middleware.RequireAuth(), func(c *fiber.Ctx) error {
				userRole := c.Locals("role").(crypto.Role)
				return c.JSON(fiber.Map{"role": userRole})
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

// Benchmark tests
func BenchmarkJWTMiddleware_RequireAuth(b *testing.B) {
	jwtService := crypto.NewJWTService("test-secret")
	middleware := NewJWTMiddleware(jwtService)
	app := setupTestApp()

	token, _ := jwtService.GenerateToken("user-123", crypto.RoleStudent, crypto.TokenTypeAccess, time.Hour)

	app.Get("/protected", middleware.RequireAuth(), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = app.Test(req)
	}
}

func BenchmarkJWTMiddleware_OptionalAuth(b *testing.B) {
	jwtService := crypto.NewJWTService("test-secret")
	middleware := NewJWTMiddleware(jwtService)
	app := setupTestApp()

	token, _ := jwtService.GenerateToken("user-123", crypto.RoleStudent, crypto.TokenTypeAccess, time.Hour)

	app.Get("/optional", middleware.OptionalAuth(), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest(http.MethodGet, "/optional", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = app.Test(req)
	}
}
