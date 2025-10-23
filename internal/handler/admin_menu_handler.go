package handler

import "github.com/gofiber/fiber/v2"

// Kontrak service yang dipakai handler Admin
type AdminMenuService interface {
	// Categories
	ListCategories(tenantID string) ([]map[string]any, error)
	CreateCategory(tenantID string, body map[string]any) (map[string]any, error)
	ReplaceCategory(tenantID, id string, body map[string]any) (map[string]any, error)
	PatchCategory(tenantID, id string, body map[string]any) (map[string]any, error)
	DeleteCategory(tenantID, id string) error

	// Items
	ListItems(tenantID, categoryID string) ([]map[string]any, error)
	CreateItem(tenantID string, body map[string]any) (map[string]any, error)
	ReplaceItem(tenantID, id string, body map[string]any) (map[string]any, error)
	PatchItem(tenantID, id string, body map[string]any) (map[string]any, error)
	DeleteItem(tenantID, id string) error
	ToggleOOS(tenantID, id string, isActive bool) (map[string]any, error)

	// Options
	ListItemOptions(itemID, tenantID string) ([]map[string]any, error)
	CreateItemOption(itemID, tenantID string, body map[string]any) (map[string]any, error)
	ListOptionValues(optionID, tenantID string) ([]map[string]any, error)
	CreateOptionValue(optionID, tenantID string, body map[string]any) (map[string]any, error)

	// Tables
	GenerateTableQR(tableID, tenantID string) (string, error)
}

// ====== Handler concrete (EXPORTed) ======
type AdminMenuHandler struct{ s AdminMenuService }

func NewAdminMenuHandler(s AdminMenuService) *AdminMenuHandler {
	return &AdminMenuHandler{s: s}
}

// === Categories
func (h *AdminMenuHandler) ListCategories(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	out, err := h.s.ListCategories(tenantID)
	if err != nil { return fiber.ErrBadRequest }
	return c.JSON(out)
}
func (h *AdminMenuHandler) CreateCategory(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	var body map[string]any
	if err := c.BodyParser(&body); err != nil { return fiber.ErrBadRequest }
	obj, err := h.s.CreateCategory(tenantID, body)
	if err != nil { return fiber.ErrBadRequest }
	return c.Status(201).JSON(obj)
}
func (h *AdminMenuHandler) ReplaceCategory(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	id := c.Params("id")
	var body map[string]any
	if err := c.BodyParser(&body); err != nil { return fiber.ErrBadRequest }
	obj, err := h.s.ReplaceCategory(tenantID, id, body)
	if err != nil { return fiber.ErrBadRequest }
	return c.JSON(obj)
}
func (h *AdminMenuHandler) PatchCategory(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	id := c.Params("id")
	var body map[string]any
	if err := c.BodyParser(&body); err != nil { return fiber.ErrBadRequest }
	obj, err := h.s.PatchCategory(tenantID, id, body)
	if err != nil { return fiber.ErrBadRequest }
	return c.JSON(obj)
}
func (h *AdminMenuHandler) DeleteCategory(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	id := c.Params("id")
	if err := h.s.DeleteCategory(tenantID, id); err != nil { return fiber.ErrBadRequest }
	return c.SendStatus(204)
}

// === Items
func (h *AdminMenuHandler) ListItems(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	cat := c.Query("category_id")
	out, err := h.s.ListItems(tenantID, cat)
	if err != nil { return fiber.ErrBadRequest }
	return c.JSON(out)
}
func (h *AdminMenuHandler) CreateItem(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	var body map[string]any
	if err := c.BodyParser(&body); err != nil { return fiber.ErrBadRequest }
	obj, err := h.s.CreateItem(tenantID, body)
	if err != nil { return fiber.ErrBadRequest }
	return c.Status(201).JSON(obj)
}
func (h *AdminMenuHandler) ReplaceItem(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	id := c.Params("id")
	var body map[string]any
	if err := c.BodyParser(&body); err != nil { return fiber.ErrBadRequest }
	obj, err := h.s.ReplaceItem(tenantID, id, body)
	if err != nil { return fiber.ErrBadRequest }
	return c.JSON(obj)
}
func (h *AdminMenuHandler) PatchItem(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	id := c.Params("id")
	var body map[string]any
	if err := c.BodyParser(&body); err != nil { return fiber.ErrBadRequest }
	obj, err := h.s.PatchItem(tenantID, id, body)
	if err != nil { return fiber.ErrBadRequest }
	return c.JSON(obj)
}
func (h *AdminMenuHandler) DeleteItem(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	id := c.Params("id")
	if err := h.s.DeleteItem(tenantID, id); err != nil { return fiber.ErrBadRequest }
	return c.SendStatus(204)
}
func (h *AdminMenuHandler) ToggleOOS(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	id := c.Params("id")
	var body struct{ IsActive bool `json:"is_active"` }
	if err := c.BodyParser(&body); err != nil { return fiber.ErrBadRequest }
	obj, err := h.s.ToggleOOS(tenantID, id, body.IsActive)
	if err != nil { return fiber.ErrBadRequest }
	return c.JSON(obj)
}

// === Options
func (h *AdminMenuHandler) ListItemOptions(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	itemID := c.Params("id")
	out, err := h.s.ListItemOptions(itemID, tenantID)
	if err != nil { return fiber.ErrBadRequest }
	return c.JSON(out)
}
func (h *AdminMenuHandler) CreateItemOption(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	itemID := c.Params("id")
	var body map[string]any
	if err := c.BodyParser(&body); err != nil { return fiber.ErrBadRequest }
	obj, err := h.s.CreateItemOption(itemID, tenantID, body)
	if err != nil { return fiber.ErrBadRequest }
	return c.Status(201).JSON(obj)
}
func (h *AdminMenuHandler) ListOptionValues(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	optionID := c.Params("option_id")
	out, err := h.s.ListOptionValues(optionID, tenantID)
	if err != nil { return fiber.ErrBadRequest }
	return c.JSON(out)
}
func (h *AdminMenuHandler) CreateOptionValue(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	optionID := c.Params("option_id")
	var body map[string]any
	if err := c.BodyParser(&body); err != nil { return fiber.ErrBadRequest }
	obj, err := h.s.CreateOptionValue(optionID, tenantID, body)
	if err != nil { return fiber.ErrBadRequest }
	return c.Status(201).JSON(obj)
}

// === Tables
func (h *AdminMenuHandler) GenerateQR(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	id := c.Params("id")
	url, err := h.s.GenerateTableQR(id, tenantID)
	if err != nil { return fiber.ErrBadRequest }
	return c.JSON(fiber.Map{"url": url})
}
