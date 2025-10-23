package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swagger "github.com/swaggo/fiber-swagger"

	"github.com/JustDoItBetter/FITS-backend/internal/common/response"
	"github.com/JustDoItBetter/FITS-backend/internal/config"
	"github.com/JustDoItBetter/FITS-backend/internal/domain/auth"
	"github.com/JustDoItBetter/FITS-backend/internal/domain/signing"
	"github.com/JustDoItBetter/FITS-backend/internal/domain/student"
	"github.com/JustDoItBetter/FITS-backend/internal/domain/teacher"
	"github.com/JustDoItBetter/FITS-backend/internal/middleware"
	"github.com/JustDoItBetter/FITS-backend/pkg/crypto"
	"github.com/JustDoItBetter/FITS-backend/pkg/database"
	"github.com/JustDoItBetter/FITS-backend/pkg/logger"

	_ "github.com/JustDoItBetter/FITS-backend/docs" // Swagger docs
	"go.uber.org/zap"
)

// @title FITS Backend API
// @version 1.0
// @description FITS (Flexible IT Training System) Backend API for managing students, teachers, and signing requests with full authentication
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@fits.example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token authentication. Format: "Bearer {token}"

func main() {
	// Load configuration
	cfg, err := config.Load("configs/config.toml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize structured logger early for consistent logging throughout application
	if err := logger.Init(cfg.Logging.Level, cfg.Logging.Format); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync() // Flush buffered logs before shutdown

	logger.Info("FITS Backend starting up",
		zap.String("version", "1.0.0"),
		zap.String("log_level", cfg.Logging.Level),
		zap.String("log_format", cfg.Logging.Format),
	)

	// Initialize database connection
	db, err := database.New(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	logger.Info("Database connected successfully",
		zap.String("host", cfg.Database.Host),
		zap.Int("port", cfg.Database.Port),
		zap.String("database", cfg.Database.Database),
	)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.GetReadTimeout(),
		WriteTimeout: cfg.GetWriteTimeout(),
		BodyLimit:    int(cfg.Storage.MaxFileSize),
		ErrorHandler: customErrorHandler,
	})

	// Global middleware stack - order matters for proper request processing

	// Panic recovery must be first to catch panics in other middleware
	app.Use(recover.New())

	// Security headers protect against common web vulnerabilities
	// Note: COOP/COEP/CORP headers are NOT set as they break Swagger UI fetch functionality
	// Using custom middleware instead of Helmet to have full control over headers
	app.Use(func(c *fiber.Ctx) error {
		// Essential security headers that don't interfere with Swagger UI
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// CSP allows Swagger UI to function properly with inline scripts/styles
		c.Set("Content-Security-Policy", "default-src 'self'; connect-src 'self' http://localhost:8080; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' data: https:; script-src 'self' 'unsafe-inline'")

		// IMPORTANT: DO NOT set Cross-Origin-Opener-Policy, Cross-Origin-Embedder-Policy,
		// or Cross-Origin-Resource-Policy - they break Swagger UI fetch functionality

		return c.Next()
	})

	// Global IP-based rate limiting (coarse-grained) prevents abuse and DoS attacks
	// This provides basic protection for unauthenticated endpoints
	if cfg.Server.RateLimit > 0 {
		app.Use(limiter.New(limiter.Config{
			Max:        cfg.Server.RateLimit,
			Expiration: 1 * time.Minute,
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"success": false,
					"error":   "Rate limit exceeded",
					"details": fmt.Sprintf("Maximum %d requests per minute allowed", cfg.Server.RateLimit),
					"code":    fiber.StatusTooManyRequests,
				})
			},
		}))
	}

	// Per-user rate limiting (fine-grained) provides role-based limits
	// This prevents authenticated users from being blocked by global limits
	// Must come AFTER JWT middleware to have access to user context
	userRateLimiter := middleware.NewUserRateLimiter(middleware.DefaultUserRateLimitConfig())
	defer userRateLimiter.Stop() // Clean up on shutdown

	app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	// CORS configured from settings to support both development (*) and production (specific origins)
	// IMPORTANT: Access the API via http://localhost:8080 (NOT http://0.0.0.0:8080)
	// Browsers treat 0.0.0.0 as untrustworthy and will block requests
	// Note: AllowCredentials must be false when using wildcard origins (security requirement)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.AllowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false, // Must be false with wildcard origins
		ExposeHeaders:    "Content-Length,Content-Type",
		MaxAge:           3600, // Cache preflight responses for 1 hour
	}))

	// Setup routes with per-user rate limiting
	setupRoutes(app, db, cfg, userRateLimiter)

	// Start server with optional TLS support and graceful shutdown handling
	startServer(app, cfg)
}

func setupRoutes(app *fiber.App, db *database.DB, cfg *config.Config, userRateLimiter *middleware.UserRateLimiter) {
	// Serve static files from web directory
	app.Static("/", "./web", fiber.Static{
		Index:         "login.html",
		Browse:        false,
		CacheDuration: 10 * time.Second,
	})

	// Health check endpoint
	// @Summary Health check
	// @Description Check if the API is running and database is accessible
	// @Tags health
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /health [get]
	app.Get("/health", func(c *fiber.Ctx) error {
		// Check database health
		if err := db.Health(c.Context()); err != nil {
			return c.Status(503).JSON(fiber.Map{
				"status":   "unhealthy",
				"database": "disconnected",
				"time":     time.Now().Format(time.RFC3339),
			})
		}

		return c.JSON(fiber.Map{
			"status":   "ok",
			"database": "connected",
			"time":     time.Now().Format(time.RFC3339),
		})
	})

	// Prometheus metrics endpoint
	app.Get("/metrics",
		middleware.SecretAuth(cfg.Secrets.MetricsSecret),
		adaptor.HTTPHandler(promhttp.Handler()),
	)

	// API Documentation - Auto-generated from code annotations with Swagger
	app.Get("/swagger/*", swagger.WrapHandler)

	// Redirect /docs and /api to Swagger UI
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html", fiber.StatusMovedPermanently)
	})
	app.Get("/api", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html", fiber.StatusMovedPermanently)
	})

	// Initialize JWT Service
	jwtService := crypto.NewJWTService(cfg.JWT.Secret)

	// Initialize JWT Middleware
	jwtMiddleware := middleware.NewJWTMiddleware(jwtService)

	// Initialize Middleware Adapter for clean dependency injection
	mwAdapter := middleware.NewMiddlewareAdapter(jwtMiddleware)

	// Initialize Auth Domain
	authRepo := auth.NewGormRepository(db.DB)
	bootstrapService := auth.NewBootstrapService(authRepo, &cfg.JWT)
	invitationService := auth.NewInvitationService(authRepo, jwtService, &cfg.JWT)
	authService := auth.NewAuthService(authRepo, jwtService, &cfg.JWT)
	authHandler := auth.NewHandler(bootstrapService, invitationService, authService)

	// Register auth routes (these don't require authentication)
	authHandler.RegisterRoutes(app)

	// Protected auth endpoints
	app.Post("/api/v1/auth/logout",
		jwtMiddleware.RequireAuth(),
		authHandler.Logout,
	)

	// Protected admin endpoints
	app.Post("/api/v1/admin/invite",
		jwtMiddleware.RequireAuth(),
		middleware.RequireAdmin(),
		authHandler.CreateInvitation,
	)

	// Initialize repositories with GORM (PostgreSQL persistence)
	studentRepo := student.NewGormRepository(db.DB)
	teacherRepo := teacher.NewGormRepository(db.DB)

	// Initialize services
	studentService := student.NewService(studentRepo)
	teacherService := teacher.NewService(teacherRepo)
	signingService := signing.NewService()

	// Initialize handlers
	studentHandler := student.NewHandler(studentService)
	teacherHandler := teacher.NewHandler(teacherService)
	signingHandler := signing.NewHandler(signingService)

	// API v1 routes - Single source of truth for all routes and security
	// Apply per-user rate limiting to all API routes (after JWT middleware extracts user info)
	api := app.Group("/api/v1")
	api.Use(userRateLimiter.Middleware())

	// Register signing routes (protected - requires authentication)
	signingGroup := api.Group("/signing")
	signingGroup.Use(jwtMiddleware.RequireAuth())
	signingHandler.RegisterRoutes(signingGroup)

	// Register student routes with their security requirements
	// Routes and middleware are now defined in one place within the handler
	studentGroup := api.Group("/student")
	studentHandler.RegisterRoutes(studentGroup, mwAdapter, mwAdapter)

	// Register teacher routes with their security requirements
	// Routes and middleware are now defined in one place within the handler
	teacherGroup := api.Group("/teacher")
	teacherHandler.RegisterRoutes(teacherGroup, mwAdapter, mwAdapter)
}

// startServer starts the HTTP/HTTPS server with optional TLS support
func startServer(app *fiber.App, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	scheme := "http"
	if cfg.Server.TLSEnabled {
		scheme = "https"
	}

	logger.Info("FITS Backend ready",
		zap.String("address", addr),
		zap.String("scheme", scheme),
		zap.Bool("tls_enabled", cfg.Server.TLSEnabled),
		zap.String("docs", fmt.Sprintf("%s://%s/docs", scheme, addr)),
		zap.String("bootstrap", fmt.Sprintf("POST %s://%s/api/v1/bootstrap/init", scheme, addr)),
	)

	// Console output for convenience
	// Use localhost for browser access, not 0.0.0.0 (browsers treat it as untrustworthy)
	displayAddr := addr
	if cfg.Server.Host == "0.0.0.0" {
		displayAddr = fmt.Sprintf("localhost:%d", cfg.Server.Port)
	}

	log.Printf("")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("FITS Backend v1.0.1 - Ready")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("")
	if cfg.Server.TLSEnabled {
		log.Printf("🔒 HTTPS Enabled")
		log.Printf("Server:        https://%s (listening on %s)", displayAddr, addr)
		log.Printf("Certificate:   %s", cfg.Server.TLSCertFile)
	} else {
		log.Printf("⚠️  HTTP Mode (Development Only)")
		log.Printf("Server:        http://%s (listening on %s)", displayAddr, addr)
		log.Printf("Browser:       http://%s 👈 USE THIS IN BROWSER", displayAddr)
	}
	log.Printf("Documentation: %s://%s/docs", scheme, displayAddr)
	log.Printf("Health Check:  %s://%s/health", scheme, displayAddr)
	log.Printf("")
	log.Printf("Quick Start:")
	log.Printf("  1. Bootstrap Admin: POST %s://%s/api/v1/bootstrap/init", scheme, displayAddr)
	log.Printf("  2. Login:           POST %s://%s/api/v1/auth/login", scheme, displayAddr)
	log.Printf("  3. View API Docs:   %s://%s/docs", scheme, displayAddr)
	log.Printf("")
	log.Printf("⚠️  IMPORTANT: Use http://localhost:8080 in your browser")
	log.Printf("   Do NOT use http://0.0.0.0:8080 - browsers will block it!")
	log.Printf("")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("")

	// Start server in goroutine for graceful shutdown support
	go func() {
		var err error
		if cfg.Server.TLSEnabled {
			logger.Info("Starting HTTPS server with TLS",
				zap.String("cert", cfg.Server.TLSCertFile),
				zap.String("key", cfg.Server.TLSKeyFile),
			)
			err = app.ListenTLS(addr, cfg.Server.TLSCertFile, cfg.Server.TLSKeyFile)
		} else {
			logger.Warn("Starting HTTP server without TLS - use HTTPS in production!")
			err = app.Listen(addr)
		}

		if err != nil {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Graceful shutdown - wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutdown signal received, gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Forced shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}

// customErrorHandler provides centralized error handling across the entire application
// Ensures consistent error response format and proper HTTP status codes
func customErrorHandler(c *fiber.Ctx, err error) error {
	return response.Error(c, err)
}
