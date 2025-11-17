package handler

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"
)

func RegisterSwaggerUI(app *fiber.App) {
	// Serve the canonical OpenAPI document located at ./openapi/openapi.yaml.
	serveSpec := func(c *fiber.Ctx) error {
		// Disable client-side HTTP caching.
		c.Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		c.Set("Pragma", "no-cache")
		c.Set("Expires", "0")
		return c.SendFile("openapi/openapi.yaml")
	}
	app.Get("/openapi.yaml", serveSpec)
	app.Get("/docs/openapi.yaml", serveSpec)

	// Append a cache-busting query parameter on each startup.
	cb := fmt.Sprintf("%d", time.Now().Unix())
	app.Get("/swagger/*", fiberSwagger.New(fiberSwagger.Config{
		URL: "/openapi.yaml?v=" + cb, // Ensure Swagger UI always reloads the latest spec.
	}))
}
