package main

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"
)

// ServeOpenAPISpec serves the OpenAPI 3.0 specification
func ServeOpenAPISpec(c *fiber.Ctx) error {
	// Read the OpenAPI YAML file
	yamlFile, err := os.ReadFile("docs/openapi.yaml")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to read OpenAPI specification",
		})
	}

	// Parse YAML to JSON for Swagger UI
	var spec interface{}
	if err := yaml.Unmarshal(yamlFile, &spec); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to parse OpenAPI specification",
		})
	}

	// Return as JSON (Swagger UI prefers JSON)
	return c.JSON(spec)
}

// ServeSwaggerUI serves a custom Swagger UI page that loads our OpenAPI spec
func ServeSwaggerUI(c *fiber.Ctx) error {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FITS Backend API - Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui.css">
    <style>
        body { margin: 0; padding: 0; }
        .swagger-ui .topbar { display: none; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = () => {
            window.ui = SwaggerUIBundle({
                url: '/api/openapi.json',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                displayRequestDuration: true,
                filter: true,
                tryItOutEnabled: true,
                persistAuthorization: true
            });
        };
    </script>
</body>
</html>
`
	c.Set("Content-Type", "text/html")
	return c.SendString(html)
}

// GetOpenAPIYAML returns the raw YAML file
func GetOpenAPIYAML(c *fiber.Ctx) error {
	// Get absolute path to OpenAPI file
	yamlPath, _ := filepath.Abs("docs/openapi.yaml")
	return c.SendFile(yamlPath)
}
