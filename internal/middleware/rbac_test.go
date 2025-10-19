package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JustDoItBetter/FITS-backend/pkg/crypto"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequireRole(t *testing.T) {
	t.Run("allows user with correct role", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/admin", func(c *fiber.Ctx) error {
			// Simulate authenticated user with admin role
			c.Locals("role", crypto.RoleAdmin)
			return c.Next()
		}, RequireRole(crypto.RoleAdmin), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("rejects user with wrong role", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/admin", func(c *fiber.Ctx) error {
			// Simulate authenticated user with student role
			c.Locals("role", crypto.RoleStudent)
			return c.Next()
		}, RequireRole(crypto.RoleAdmin), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("rejects unauthenticated user", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/admin", RequireRole(crypto.RoleAdmin), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("student cannot access teacher endpoint", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/teacher", func(c *fiber.Ctx) error {
			c.Locals("role", crypto.RoleStudent)
			return c.Next()
		}, RequireRole(crypto.RoleTeacher), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/teacher", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("teacher cannot access admin endpoint", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/admin", func(c *fiber.Ctx) error {
			c.Locals("role", crypto.RoleTeacher)
			return c.Next()
		}, RequireRole(crypto.RoleAdmin), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}

func TestRequireAdmin(t *testing.T) {
	t.Run("allows admin", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/admin", func(c *fiber.Ctx) error {
			c.Locals("role", crypto.RoleAdmin)
			return c.Next()
		}, RequireAdmin(), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("rejects teacher", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/admin", func(c *fiber.Ctx) error {
			c.Locals("role", crypto.RoleTeacher)
			return c.Next()
		}, RequireAdmin(), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("rejects student", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/admin", func(c *fiber.Ctx) error {
			c.Locals("role", crypto.RoleStudent)
			return c.Next()
		}, RequireAdmin(), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}

func TestRequireTeacher(t *testing.T) {
	t.Run("allows teacher", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/teacher", func(c *fiber.Ctx) error {
			c.Locals("role", crypto.RoleTeacher)
			return c.Next()
		}, RequireTeacher(), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/teacher", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("rejects student", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/teacher", func(c *fiber.Ctx) error {
			c.Locals("role", crypto.RoleStudent)
			return c.Next()
		}, RequireTeacher(), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/teacher", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("allows admin", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/teacher", func(c *fiber.Ctx) error {
			c.Locals("role", crypto.RoleAdmin)
			return c.Next()
		}, RequireTeacher(), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/teacher", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestRequireStudent(t *testing.T) {
	t.Run("allows student", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/student", func(c *fiber.Ctx) error {
			c.Locals("role", crypto.RoleStudent)
			return c.Next()
		}, RequireStudent(), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/student", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("allows teacher", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/student", func(c *fiber.Ctx) error {
			c.Locals("role", crypto.RoleTeacher)
			return c.Next()
		}, RequireStudent(), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/student", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("allows admin", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/student", func(c *fiber.Ctx) error {
			c.Locals("role", crypto.RoleAdmin)
			return c.Next()
		}, RequireStudent(), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/student", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestRequireOwnership(t *testing.T) {
	t.Run("allows user accessing own resource", func(t *testing.T) {
		app := setupTestApp()

		userID := "user-123"
		app.Get("/user/:id", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			c.Locals("role", crypto.RoleStudent)
			return c.Next()
		}, RequireOwnership("id"), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/user/"+userID, nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("rejects user accessing other resource", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/user/:id", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user-123")
			c.Locals("role", crypto.RoleStudent)
			return c.Next()
		}, RequireOwnership("id"), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/user/user-456", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("allows admin accessing any resource", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/user/:id", func(c *fiber.Ctx) error {
			c.Locals("user_id", "admin-001")
			c.Locals("role", crypto.RoleAdmin)
			return c.Next()
		}, RequireOwnership("id"), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/user/user-999", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("rejects unauthenticated user", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/user/:id", RequireOwnership("id"), func(c *fiber.Ctx) error {
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/user/user-123", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestSecretAuth(t *testing.T) {
	expectedSecret := "super-secret-key"

	t.Run("allows request with correct secret", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/metrics", SecretAuth(expectedSecret), func(c *fiber.Ctx) error {
			return c.SendString("metrics data")
		})

		req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		req.Header.Set("Authorization", "Bearer "+expectedSecret)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("rejects request with wrong secret", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/metrics", SecretAuth(expectedSecret), func(c *fiber.Ctx) error {
			return c.SendString("metrics data")
		})

		req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		req.Header.Set("Authorization", "Bearer wrong-secret")
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("rejects request without secret", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/metrics", SecretAuth(expectedSecret), func(c *fiber.Ctx) error {
			return c.SendString("metrics data")
		})

		req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("handles empty secret configuration", func(t *testing.T) {
		app := setupTestApp()

		app.Get("/metrics", SecretAuth(""), func(c *fiber.Ctx) error {
			return c.SendString("metrics data")
		})

		req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		// Should still reject if no secret provided
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestRBACIntegration(t *testing.T) {
	t.Run("full auth flow with role checking", func(t *testing.T) {
		app := setupTestApp()

		// Simulate full authentication + RBAC
		app.Get("/admin/users", func(c *fiber.Ctx) error {
			// Simulates JWT middleware setting user context
			c.Locals("user_id", "admin-001")
			c.Locals("role", crypto.RoleAdmin)
			return c.Next()
		}, RequireAdmin(), func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "admin access granted"})
		})

		req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("chained middleware - auth + ownership", func(t *testing.T) {
		app := setupTestApp()

		userID := "user-123"
		app.Get("/user/:id/profile", func(c *fiber.Ctx) error {
			c.Locals("user_id", userID)
			c.Locals("role", crypto.RoleStudent)
			return c.Next()
		}, RequireOwnership("id"), func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"profile": "data"})
		})

		req := httptest.NewRequest(http.MethodGet, "/user/"+userID+"/profile", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// Benchmark tests
func BenchmarkRequireRole(b *testing.B) {
	app := setupTestApp()

	app.Get("/test", func(c *fiber.Ctx) error {
		c.Locals("role", crypto.RoleAdmin)
		return c.Next()
	}, RequireRole(crypto.RoleAdmin), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = app.Test(req)
	}
}

func BenchmarkRequireOwnership(b *testing.B) {
	app := setupTestApp()

	userID := "user-123"
	app.Get("/user/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		c.Locals("role", crypto.RoleStudent)
		return c.Next()
	}, RequireOwnership("id"), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest(http.MethodGet, "/user/"+userID, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = app.Test(req)
	}
}
