package handler

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"
)

// RegisterSwaggerUI men-setup Swagger UI agar memakai file YAML kita.
func RegisterSwaggerUI(app *fiber.App) {
	// Serve file YAML di /docs/openapi.yaml
	app.Static("/docs", "./openapi")

	// Arahkan Swagger UI ke YAML tsb
	app.Get("/swagger/*", fiberSwagger.New(fiberSwagger.Config{
		URL: "/docs/openapi.yaml", // <— penting: path ke YAML
	}))
}
