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
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/JustDoItBetter/FITS-backend/internal/config"
	"github.com/JustDoItBetter/FITS-backend/internal/middleware"
)

func main() {
	cfg, err := config.Load("configs/config.toml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.GetReadTimeout(),
		WriteTimeout: cfg.GetWriteTimeout(),
		BodyLimit:    int(cfg.Storage.MaxFileSize),
		ErrorHandler: customErrorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	setupRoutes(app, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", addr)

	go func() {
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRoutes(app *fiber.App, cfg *config.Config) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	app.Get("/metrics",
		middleware.SecretAuth(cfg.Secrets.MetricsSecret),
		adaptor.HTTPHandler(promhttp.Handler()),
	)

	api := app.Group("/api/v1")

	api.Post("/upload", func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{"error": "not implemented yet"})
	})

	api.Get("/sign_requests", func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{"error": "not implemented yet"})
	})

	api.Post("/sign_uploads", func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{"error": "not implemented yet"})
	})

	student := api.Group("/student")
	student.Put("/", func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{"error": "not implemented yet"})
	})
	student.Post("/", func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{"error": "not implemented yet"})
	})
	student.Delete("/:uuid", func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{"error": "not implemented yet"})
	})

	teacher := api.Group("/teacher")
	teacher.Post("/", func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{"error": "not implemented yet"})
	})
	teacher.Post("/update", func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{"error": "not implemented yet"})
	})
	teacher.Delete("/:uuid", func(c *fiber.Ctx) error {
		return c.Status(501).JSON(fiber.Map{"error": "not implemented yet"})
	})
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
		"code":  code,
	})
}
