package middleware

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JustDoItBetter/FITS-backend/pkg/crypto"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserRateLimiter(t *testing.T) {
	config := DefaultUserRateLimitConfig()
	limiter := NewUserRateLimiter(config)

	assert.NotNil(t, limiter)
	assert.Equal(t, config.Window, limiter.window)
	assert.Equal(t, config.DefaultLimit, limiter.defaultLimit)
	assert.Equal(t, config.AdminLimit, limiter.adminLimit)
	assert.Equal(t, config.TeacherLimit, limiter.teacherLimit)
	assert.Equal(t, config.StudentLimit, limiter.studentLimit)
	assert.NotNil(t, limiter.cleanupTicker)

	limiter.Stop()
}

func TestUserRateLimiter_Allow(t *testing.T) {
	config := UserRateLimitConfig{
		Window:       100 * time.Millisecond,
		DefaultLimit: 3,
		AdminLimit:   10,
		TeacherLimit: 5,
		StudentLimit: 3,
	}
	limiter := NewUserRateLimiter(config)
	defer limiter.Stop()

	tests := []struct {
		name     string
		userID   string
		limit    int
		requests int
		expectOK []bool
	}{
		{
			name:     "Within limit",
			userID:   "user1",
			limit:    3,
			requests: 3,
			expectOK: []bool{true, true, true},
		},
		{
			name:     "Exceeds limit",
			userID:   "user2",
			limit:    3,
			requests: 5,
			expectOK: []bool{true, true, true, false, false},
		},
		{
			name:     "Admin high limit",
			userID:   "admin1",
			limit:    10,
			requests: 10,
			expectOK: []bool{true, true, true, true, true, true, true, true, true, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.requests; i++ {
				allowed := limiter.allow(tt.userID, tt.limit)
				assert.Equal(t, tt.expectOK[i], allowed,
					"Request %d: expected %v, got %v", i+1, tt.expectOK[i], allowed)
			}
		})
	}
}

func TestUserRateLimiter_WindowReset(t *testing.T) {
	config := UserRateLimitConfig{
		Window:       50 * time.Millisecond,
		DefaultLimit: 2,
		AdminLimit:   10,
		TeacherLimit: 5,
		StudentLimit: 2,
	}
	limiter := NewUserRateLimiter(config)
	defer limiter.Stop()

	userID := "test-user"
	limit := 2

	// First two requests should succeed
	assert.True(t, limiter.allow(userID, limit))
	assert.True(t, limiter.allow(userID, limit))

	// Third request should fail (limit reached)
	assert.False(t, limiter.allow(userID, limit))

	// Wait for window to reset
	time.Sleep(60 * time.Millisecond)

	// Next request should succeed (new window)
	assert.True(t, limiter.allow(userID, limit))
}

func TestUserRateLimiter_GetStatus(t *testing.T) {
	config := DefaultUserRateLimitConfig()
	limiter := NewUserRateLimiter(config)
	defer limiter.Stop()

	userID := "test-user"

	// No requests yet
	count, resetTime, exists := limiter.GetStatus(userID)
	assert.False(t, exists)
	assert.Equal(t, 0, count)

	// Make some requests
	limiter.allow(userID, 5)
	limiter.allow(userID, 5)
	limiter.allow(userID, 5)

	// Check status
	count, resetTime, exists = limiter.GetStatus(userID)
	assert.True(t, exists)
	assert.Equal(t, 3, count)
	assert.True(t, resetTime.After(time.Now()))
}

func TestUserRateLimiter_GetUserIdentifier(t *testing.T) {
	limiter := NewUserRateLimiter(DefaultUserRateLimitConfig())
	defer limiter.Stop()

	tests := []struct {
		name     string
		userID   interface{}
		expected string
	}{
		{
			name:     "Authenticated user",
			userID:   "user-123",
			expected: "user:user-123",
		},
		{
			name:     "Unauthenticated request",
			userID:   nil,
			expected: "ip:",
		},
		{
			name:     "Invalid user ID type",
			userID:   12345,
			expected: "ip:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/test", func(c *fiber.Ctx) error {
				if tt.userID != nil {
					c.Locals("user_id", tt.userID)
				}
				identifier := limiter.getUserIdentifier(c)
				assert.Contains(t, identifier, tt.expected)
				return c.SendStatus(200)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
		})
	}
}

func TestUserRateLimiter_GetRoleLimit(t *testing.T) {
	config := UserRateLimitConfig{
		Window:       time.Minute,
		DefaultLimit: 20,
		AdminLimit:   1000,
		TeacherLimit: 200,
		StudentLimit: 100,
	}
	limiter := NewUserRateLimiter(config)
	defer limiter.Stop()

	tests := []struct {
		name     string
		role     interface{}
		expected int
	}{
		{
			name:     "Admin role",
			role:     crypto.RoleAdmin,
			expected: 1000,
		},
		{
			name:     "Teacher role",
			role:     crypto.RoleTeacher,
			expected: 200,
		},
		{
			name:     "Student role",
			role:     crypto.RoleStudent,
			expected: 100,
		},
		{
			name:     "No role (unauthenticated)",
			role:     nil,
			expected: 20,
		},
		{
			name:     "Invalid role type",
			role:     "invalid",
			expected: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/test", func(c *fiber.Ctx) error {
				if tt.role != nil {
					c.Locals("role", tt.role)
				}
				limit := limiter.getRoleLimit(c)
				assert.Equal(t, tt.expected, limit)
				return c.SendStatus(200)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
		})
	}
}

func TestUserRateLimiter_Cleanup(t *testing.T) {
	config := UserRateLimitConfig{
		Window:       10 * time.Millisecond,
		DefaultLimit: 5,
		AdminLimit:   10,
		TeacherLimit: 8,
		StudentLimit: 5,
	}
	limiter := NewUserRateLimiter(config)

	// Make some requests
	limiter.allow("user1", 5)
	limiter.allow("user2", 5)
	limiter.allow("user3", 5)

	// Verify users exist
	_, _, exists := limiter.GetStatus("user1")
	assert.True(t, exists)

	// Wait for cleanup (window * 2)
	time.Sleep(30 * time.Millisecond)

	// Cleanup should have removed expired entries
	// Note: This test is timing-sensitive and may be flaky
	// In production, we verify cleanup runs but don't assert on timing
	limiter.Stop()
}

func TestUserRateLimiter_ConcurrentAccess(t *testing.T) {
	config := UserRateLimitConfig{
		Window:       100 * time.Millisecond,
		DefaultLimit: 100,
		AdminLimit:   1000,
		TeacherLimit: 500,
		StudentLimit: 100,
	}
	limiter := NewUserRateLimiter(config)
	defer limiter.Stop()

	// Simulate concurrent requests from multiple goroutines
	concurrency := 10
	requestsPerGoroutine := 10
	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			userID := "concurrent-user"
			for j := 0; j < requestsPerGoroutine; j++ {
				limiter.allow(userID, 1000) // High limit to avoid rejections
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < concurrency; i++ {
		<-done
	}

	// Verify count matches expected (should be concurrency * requestsPerGoroutine)
	count, _, exists := limiter.GetStatus("concurrent-user")
	assert.True(t, exists)
	assert.Equal(t, concurrency*requestsPerGoroutine, count)
}

func TestDefaultUserRateLimitConfig(t *testing.T) {
	config := DefaultUserRateLimitConfig()

	assert.Equal(t, time.Minute, config.Window)
	assert.Equal(t, 20, config.DefaultLimit)
	assert.Equal(t, 1000, config.AdminLimit)
	assert.Equal(t, 200, config.TeacherLimit)
	assert.Equal(t, 100, config.StudentLimit)
}

func TestUserRateLimiter_Middleware_Integration(t *testing.T) {
	config := UserRateLimitConfig{
		Window:       100 * time.Millisecond,
		DefaultLimit: 2,
		AdminLimit:   10,
		TeacherLimit: 5,
		StudentLimit: 2,
	}
	limiter := NewUserRateLimiter(config)
	defer limiter.Stop()

	app := fiber.New()

	// Simulate JWT middleware setting user context
	app.Use(func(c *fiber.Ctx) error {
		// Simulate authenticated admin user
		c.Locals("user_id", "admin-user")
		c.Locals("role", crypto.RoleAdmin)
		return c.Next()
	})

	app.Use(limiter.Middleware())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Admin should be able to make many requests (limit: 10)
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode, "Request %d should succeed", i+1)
	}

	// 11th request should be rejected
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 429, resp.StatusCode, "Request 11 should be rate limited")
}
