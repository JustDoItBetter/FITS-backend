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
	"github.com/gofiber/fiber/v2/middleware/helmet"
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
	app.Use(helmet.New(helmet.Config{
		XSSProtection:      "1; mode=block", // Prevents XSS attacks in older browsers
		ContentTypeNosniff: "nosniff",       // Prevents MIME sniffing
		XFrameOptions:      "SAMEORIGIN",    // Prevents clickjacking
		ReferrerPolicy:     "strict-origin-when-cross-origin",
		// CSP allows Swagger UI to load fonts and styles from CDN
		ContentSecurityPolicy: "default-src 'self'; style-src 'self' 'unsafe-inline' fonts.googleapis.com; font-src 'self' fonts.gstatic.com; img-src 'self' data:; script-src 'self' 'unsafe-inline'",
	}))

	// Rate limiting prevents abuse and DoS attacks
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

	app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	// CORS configured from settings to support both development (*) and production (specific origins)
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.Server.AllowedOrigins,
		AllowMethods: "GET,POST,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Setup routes
	setupRoutes(app, db, cfg)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	logger.Info("FITS Backend ready",
		zap.String("address", addr),
		zap.String("docs", fmt.Sprintf("http://%s/docs", addr)),
		zap.String("bootstrap", fmt.Sprintf("POST http://%s/api/v1/bootstrap/init", addr)),
	)

	// Console output for convenience
	log.Printf("")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("FITS Backend v1.0.0 - Ready")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("")
	log.Printf("Server:        http://%s", addr)
	log.Printf("Documentation: http://%s/docs", addr)
	log.Printf("Health Check:  http://%s/health", addr)
	log.Printf("")
	log.Printf("Quick Start:")
	log.Printf("  1. Bootstrap Admin: POST http://%s/api/v1/bootstrap/init", addr)
	log.Printf("  2. Login:           POST http://%s/api/v1/auth/login", addr)
	log.Printf("  3. View API Docs:   http://%s/docs", addr)
	log.Printf("")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("")

	// Start server in goroutine for graceful shutdown support
	go func() {
		if err := app.Listen(addr); err != nil {
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

func setupRoutes(app *fiber.App, db *database.DB, cfg *config.Config) {
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

	// API v1 routes
	api := app.Group("/api/v1")

	// Register signing routes (protected)
	signingGroup := api.Group("/signing")
	signingGroup.Use(jwtMiddleware.RequireAuth())
	signingHandler.RegisterRoutes(signingGroup)

	// Register student routes
	studentGroup := api.Group("/student")
	// POST for creation follows REST conventions and allows server-generated UUIDs if client doesn't provide one
	studentGroup.Post("/",
		jwtMiddleware.RequireAuth(),
		middleware.RequireRole(crypto.RoleAdmin), // Prevents unauthorized student registration
		studentHandler.Create,
	)
	// Optional auth allows public profile viewing while enabling role-based data filtering for authenticated users
	studentGroup.Get("/:uuid",
		jwtMiddleware.OptionalAuth(),
		studentHandler.GetByUUID,
	)
	// PUT for full resource replacement follows REST semantics for idempotent updates
	studentGroup.Put("/:uuid",
		jwtMiddleware.RequireAuth(),
		middleware.RequireRole(crypto.RoleAdmin), // Student data modifications restricted to admins for data integrity
		studentHandler.Update,
	)
	// Soft delete preferred over hard delete for audit trail and data recovery
	studentGroup.Delete("/:uuid",
		jwtMiddleware.RequireAuth(),
		middleware.RequireRole(crypto.RoleAdmin),
		studentHandler.Delete,
	)
	// Optional auth enables public listing for discovery while allowing filtering based on user role
	studentGroup.Get("/",
		jwtMiddleware.OptionalAuth(),
		studentHandler.List,
	)

	// Register teacher routes
	teacherGroup := api.Group("/teacher")
	// POST for creation follows REST conventions
	teacherGroup.Post("/",
		jwtMiddleware.RequireAuth(),
		middleware.RequireRole(crypto.RoleAdmin), // Teacher accounts must be created by admins to ensure valid credentials
		teacherHandler.Create,
	)
	// Optional auth allows public teacher directory viewing
	teacherGroup.Get("/:uuid",
		jwtMiddleware.OptionalAuth(),
		teacherHandler.GetByUUID,
	)
	// PUT for full resource replacement follows REST semantics
	teacherGroup.Put("/:uuid",
		jwtMiddleware.RequireAuth(),
		middleware.RequireRole(crypto.RoleAdmin), // Teacher profile changes require admin approval for consistency
		teacherHandler.Update,
	)
	// Cascading deletes handled at database level to maintain referential integrity
	teacherGroup.Delete("/:uuid",
		jwtMiddleware.RequireAuth(),
		middleware.RequireRole(crypto.RoleAdmin),
		teacherHandler.Delete,
	)
	// Optional auth enables public teacher directory
	teacherGroup.Get("/",
		jwtMiddleware.OptionalAuth(),
		teacherHandler.List,
	)
}

// customErrorHandler provides centralized error handling across the entire application
// Ensures consistent error response format and proper HTTP status codes
func customErrorHandler(c *fiber.Ctx, err error) error {
	return response.Error(c, err)
}
