package middleware

import (
	"sync"
	"time"

	"github.com/JustDoItBetter/FITS-backend/internal/common/response"
	"github.com/JustDoItBetter/FITS-backend/pkg/crypto"
	"github.com/gofiber/fiber/v2"
)

// UserRateLimiter implements per-user rate limiting with role-based limits
type UserRateLimiter struct {
	requests      map[string]*userRequests
	mu            sync.RWMutex
	window        time.Duration
	defaultLimit  int
	adminLimit    int
	teacherLimit  int
	studentLimit  int
	cleanupTicker *time.Ticker
}

// userRequests tracks requests for a single user within a time window
type userRequests struct {
	count     int
	resetTime time.Time
	mu        sync.Mutex
}

// UserRateLimitConfig configures per-user rate limiting
type UserRateLimitConfig struct {
	Window       time.Duration // Time window for rate limiting (default: 1 minute)
	DefaultLimit int           // Limit for unauthenticated requests (default: 20)
	AdminLimit   int           // Limit for admin users (default: 1000)
	TeacherLimit int           // Limit for teacher users (default: 200)
	StudentLimit int           // Limit for student users (default: 100)
}

// DefaultUserRateLimitConfig returns sensible defaults
func DefaultUserRateLimitConfig() UserRateLimitConfig {
	return UserRateLimitConfig{
		Window:       time.Minute,
		DefaultLimit: 20,
		AdminLimit:   1000,
		TeacherLimit: 200,
		StudentLimit: 100,
	}
}

// NewUserRateLimiter creates a new per-user rate limiter
func NewUserRateLimiter(config UserRateLimitConfig) *UserRateLimiter {
	limiter := &UserRateLimiter{
		requests:     make(map[string]*userRequests),
		window:       config.Window,
		defaultLimit: config.DefaultLimit,
		adminLimit:   config.AdminLimit,
		teacherLimit: config.TeacherLimit,
		studentLimit: config.StudentLimit,
	}

	// Start background cleanup goroutine to prevent memory leaks
	limiter.cleanupTicker = time.NewTicker(config.Window * 2)
	go limiter.cleanup()

	return limiter
}

// cleanup removes expired rate limit entries to prevent memory leaks
func (l *UserRateLimiter) cleanup() {
	for range l.cleanupTicker.C {
		l.mu.Lock()
		now := time.Now()
		for userID, req := range l.requests {
			req.mu.Lock()
			if now.After(req.resetTime) {
				delete(l.requests, userID)
			}
			req.mu.Unlock()
		}
		l.mu.Unlock()
	}
}

// Stop stops the cleanup ticker (call when shutting down)
func (l *UserRateLimiter) Stop() {
	if l.cleanupTicker != nil {
		l.cleanupTicker.Stop()
	}
}

// Middleware returns a Fiber middleware handler for per-user rate limiting
func (l *UserRateLimiter) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from context (set by JWT middleware)
		// If not authenticated, use IP address as identifier
		userID := l.getUserIdentifier(c)

		// Get user's role to determine rate limit
		limit := l.getRoleLimit(c)

		// Check and increment rate limit
		if !l.allow(userID, limit) {
			return response.Error(c, fiber.NewError(
				fiber.StatusTooManyRequests,
				"Rate limit exceeded. Please try again later.",
			))
		}

		return c.Next()
	}
}

// getUserIdentifier gets a unique identifier for the user
// Uses authenticated user ID if available, otherwise falls back to IP
func (l *UserRateLimiter) getUserIdentifier(c *fiber.Ctx) string {
	// Try to get authenticated user ID
	userIDRaw := c.Locals("user_id")
	if userIDRaw != nil {
		if userID, ok := userIDRaw.(string); ok && userID != "" {
			return "user:" + userID
		}
	}

	// Fall back to IP address for unauthenticated requests
	return "ip:" + c.IP()
}

// getRoleLimit returns the rate limit based on user's role
func (l *UserRateLimiter) getRoleLimit(c *fiber.Ctx) int {
	roleRaw := c.Locals("role")
	if roleRaw == nil {
		return l.defaultLimit
	}

	role, ok := roleRaw.(crypto.Role)
	if !ok {
		return l.defaultLimit
	}

	switch role {
	case crypto.RoleAdmin:
		return l.adminLimit
	case crypto.RoleTeacher:
		return l.teacherLimit
	case crypto.RoleStudent:
		return l.studentLimit
	default:
		return l.defaultLimit
	}
}

// allow checks if a request should be allowed based on rate limiting
func (l *UserRateLimiter) allow(userID string, limit int) bool {
	now := time.Now()

	l.mu.Lock()
	req, exists := l.requests[userID]
	if !exists {
		// First request from this user in this window
		req = &userRequests{
			count:     1,
			resetTime: now.Add(l.window),
		}
		l.requests[userID] = req
		l.mu.Unlock()
		return true
	}
	l.mu.Unlock()

	// Check if window has expired
	req.mu.Lock()
	defer req.mu.Unlock()

	if now.After(req.resetTime) {
		// Window expired, reset counter
		req.count = 1
		req.resetTime = now.Add(l.window)
		return true
	}

	// Check if limit exceeded
	if req.count >= limit {
		return false
	}

	// Increment counter
	req.count++
	return true
}

// GetStatus returns current rate limit status for a user (for debugging/monitoring)
func (l *UserRateLimiter) GetStatus(userID string) (count int, resetTime time.Time, exists bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	req, exists := l.requests[userID]
	if !exists {
		return 0, time.Time{}, false
	}

	req.mu.Lock()
	defer req.mu.Unlock()

	return req.count, req.resetTime, true
}
