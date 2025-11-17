// internal/transport/http/routes.go
package http

import (
	"qrmenu/internal/handler"
	"qrmenu/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type Deps struct {
	Auth      *handler.AuthHandler
	Menu      *handler.MenuHandler
	Table     *handler.TableHandler
	OrderPub  *handler.OrderPublicHandler
	AdminMenu *handler.AdminMenuHandler
	AdminOrd  *handler.AdminOrdersHandler
	Setup     *handler.SetupHandler
	JWTSecret string
}

func Register(app *fiber.App, d Deps) {
	// ---- Health (public) ----
	// Ensure handler.Health() exists and returns a fiber.Handler.
	app.Get("/health", handler.Health())

	// ---- Setup (public) ----
	app.Get("/setup/status", d.Setup.Status)
	app.Post("/setup/admin", d.Setup.SetupTenant)

	// ---- Public / Customer ----
	app.Get("/api/v1/table/:token", d.Table.Resolve)
	app.Get("/api/v1/menu", d.Menu.Get)
	app.Post("/api/v1/orders", d.OrderPub.Create)

	// ---- Auth (cookie) ----
	app.Post("/auth/login", d.Auth.Login)
	app.Post("/auth/logout", d.Auth.Logout)

	// ---- Admin (cookie protected) ----
	admin := app.Group("/admin", middleware.AdminCookieOnly(d.JWTSecret))

	// Orders
	admin.Get("/orders", d.AdminOrd.List)
	admin.Patch("/orders/:id/status", d.AdminOrd.PatchStatus)

	// Categories
	admin.Get("/categories", d.AdminMenu.ListCategories)
	admin.Post("/categories", d.AdminMenu.CreateCategory)
	admin.Put("/categories/:id", d.AdminMenu.ReplaceCategory)
	admin.Patch("/categories/:id", d.AdminMenu.PatchCategory)
	admin.Delete("/categories/:id", d.AdminMenu.DeleteCategory)

	// Items
	admin.Get("/items", d.AdminMenu.ListItems)
	admin.Post("/items", d.AdminMenu.CreateItem)
	admin.Put("/items/:id", d.AdminMenu.ReplaceItem)
	admin.Patch("/items/:id", d.AdminMenu.PatchItem)
	admin.Delete("/items/:id", d.AdminMenu.DeleteItem)
	admin.Patch("/items/:id/oos", d.AdminMenu.ToggleOOS)

	// Options
	admin.Get("/items/:id/options", d.AdminMenu.ListItemOptions)
	admin.Post("/items/:id/options", d.AdminMenu.CreateItemOption)
	admin.Get("/options/:option_id/values", d.AdminMenu.ListOptionValues)
	admin.Post("/options/:option_id/values", d.AdminMenu.CreateOptionValue)

	// Tables
	admin.Post("/tables/:id/qr", d.AdminMenu.GenerateQR)
}
