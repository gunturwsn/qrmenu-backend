package handler

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"

	"qrmenu/internal/domain"
	"qrmenu/internal/platform/logging"
)

// AdminMenuUseCase models the business operations required by the admin menu HTTP adapter.
type AdminMenuUseCase interface {
	ListCategories(tenantID string) ([]domain.Category, error)
	CreateCategory(tenantID string, body map[string]any) (*domain.Category, error)
	ReplaceCategory(tenantID, id string, body map[string]any) (*domain.Category, error)
	PatchCategory(tenantID, id string, body map[string]any) (*domain.Category, error)
	DeleteCategory(tenantID, id string) error

	ListItems(tenantID, categoryID string) ([]domain.Item, error)
	CreateItem(tenantID string, body map[string]any) (*domain.Item, error)
	ReplaceItem(tenantID, id string, body map[string]any) (*domain.Item, error)
	PatchItem(tenantID, id string, body map[string]any) (*domain.Item, error)
	DeleteItem(tenantID, id string) error
	ToggleOOS(tenantID, id string, isActive bool) (*domain.Item, error)

	ListItemOptions(itemID, tenantID string) ([]domain.ItemOption, error)
	CreateItemOption(itemID, tenantID string, body map[string]any) (*domain.ItemOption, error)
	ListOptionValues(optionID, tenantID string) ([]domain.ItemOptionValue, error)
	CreateOptionValue(optionID, tenantID string, body map[string]any) (*domain.ItemOptionValue, error)

	GenerateTableQR(tableID, tenantID string) (string, error)
}

// AdminMenuHandler exposes HTTP handlers that orchestrate admin menu use cases.
type AdminMenuHandler struct {
	uc AdminMenuUseCase
}

// NewAdminMenuHandler wires the admin menu use case into a HTTP handler instance.
func NewAdminMenuHandler(uc AdminMenuUseCase) *AdminMenuHandler {
	return &AdminMenuHandler{uc: uc}
}

// categoryResponse describes the JSON payload returned for category endpoints.
type categoryResponse struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Name     string `json:"name"`
	Sort     int    `json:"sort"`
	IsActive bool   `json:"is_active"`
}

// itemResponse describes the JSON payload returned for item endpoints.
type itemResponse struct {
	ID          string            `json:"id"`
	TenantID    string            `json:"tenant_id"`
	CategoryID  string            `json:"category_id"`
	Name        string            `json:"name"`
	Description *string           `json:"description"`
	Price       int64             `json:"price"`
	PhotoURL    *string           `json:"photo_url"`
	Flags       datatypes.JSONMap `json:"flags,omitempty"`
	IsActive    bool              `json:"is_active"`
}

// optionResponse describes the JSON payload returned for item option endpoints.
type optionResponse struct {
	ID       string `json:"id"`
	ItemID   string `json:"item_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

// optionValueResponse describes the JSON payload returned for option value endpoints.
type optionValueResponse struct {
	ID         string `json:"id"`
	OptionID   string `json:"option_id"`
	Label      string `json:"label"`
	DeltaPrice int64  `json:"delta_price"`
}

// ListCategories returns all categories owned by the authenticated tenant.
func (h *AdminMenuHandler) ListCategories(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)

	cats, err := h.uc.ListCategories(tenantID)
	if err != nil {
		logging.HandlerError(c, "AdminMenu.ListCategories", "service error", fiber.StatusBadRequest, "categories_list_failed", err, "tenant_id", tenantID)
		return fiber.ErrBadRequest
	}

	resp := make([]categoryResponse, 0, len(cats))
	for _, cat := range cats {
		resp = append(resp, newCategoryResponse(cat))
	}

	logging.HandlerInfo(c, "AdminMenu.ListCategories", "categories listed", fiber.StatusOK, "categories_listed", "tenant_id", tenantID, "count", len(resp))
	return c.JSON(resp)
}

// CreateCategory persists a new category for the authenticated tenant.
func (h *AdminMenuHandler) CreateCategory(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)

	var payload map[string]any
	if err := c.BodyParser(&payload); err != nil {
		logging.HandlerError(c, "AdminMenu.CreateCategory", "failed to parse body", fiber.StatusBadRequest, "invalid_body", err, "tenant_id", tenantID)
		return fiber.ErrBadRequest
	}

	cat, err := h.uc.CreateCategory(tenantID, payload)
	if err != nil {
		logging.HandlerError(c, "AdminMenu.CreateCategory", "service error", fiber.StatusBadRequest, "category_create_failed", err, "tenant_id", tenantID)
		return fiber.ErrBadRequest
	}

	resp := newCategoryResponse(*cat)
	logging.HandlerInfo(c, "AdminMenu.CreateCategory", "category created", fiber.StatusCreated, "category_created", "tenant_id", tenantID, "category_id", resp.ID)
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// ReplaceCategory performs a full update on an existing category.
func (h *AdminMenuHandler) ReplaceCategory(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	categoryID := c.Params("id")

	var payload map[string]any
	if err := c.BodyParser(&payload); err != nil {
		logging.HandlerError(c, "AdminMenu.ReplaceCategory", "failed to parse body", fiber.StatusBadRequest, "invalid_body", err, "tenant_id", tenantID, "category_id", categoryID)
		return fiber.ErrBadRequest
	}

	cat, err := h.uc.ReplaceCategory(tenantID, categoryID, payload)
	if err != nil {
		logging.HandlerError(c, "AdminMenu.ReplaceCategory", "service error", fiber.StatusBadRequest, "category_replace_failed", err, "tenant_id", tenantID, "category_id", categoryID)
		return fiber.ErrBadRequest
	}

	resp := newCategoryResponse(*cat)
	logging.HandlerInfo(c, "AdminMenu.ReplaceCategory", "category replaced", fiber.StatusOK, "category_replaced", "tenant_id", tenantID, "category_id", categoryID)
	return c.JSON(resp)
}

// PatchCategory partially updates an existing category.
func (h *AdminMenuHandler) PatchCategory(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	categoryID := c.Params("id")

	var payload map[string]any
	if err := c.BodyParser(&payload); err != nil {
		logging.HandlerError(c, "AdminMenu.PatchCategory", "failed to parse body", fiber.StatusBadRequest, "invalid_body", err, "tenant_id", tenantID, "category_id", categoryID)
		return fiber.ErrBadRequest
	}

	cat, err := h.uc.PatchCategory(tenantID, categoryID, payload)
	if err != nil {
		logging.HandlerError(c, "AdminMenu.PatchCategory", "service error", fiber.StatusBadRequest, "category_patch_failed", err, "tenant_id", tenantID, "category_id", categoryID)
		return fiber.ErrBadRequest
	}

	resp := newCategoryResponse(*cat)
	logging.HandlerInfo(c, "AdminMenu.PatchCategory", "category patched", fiber.StatusOK, "category_patched", "tenant_id", tenantID, "category_id", categoryID)
	return c.JSON(resp)
}

// DeleteCategory removes a category owned by the authenticated tenant.
func (h *AdminMenuHandler) DeleteCategory(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	categoryID := c.Params("id")

	if err := h.uc.DeleteCategory(tenantID, categoryID); err != nil {
		logging.HandlerError(c, "AdminMenu.DeleteCategory", "service error", fiber.StatusBadRequest, "category_delete_failed", err, "tenant_id", tenantID, "category_id", categoryID)
		return fiber.ErrBadRequest
	}

	logging.HandlerInfo(c, "AdminMenu.DeleteCategory", "category deleted", fiber.StatusNoContent, "category_deleted", "tenant_id", tenantID, "category_id", categoryID)
	return c.SendStatus(fiber.StatusNoContent)
}

// ListItems returns items filtered by tenant and optional category.
func (h *AdminMenuHandler) ListItems(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	categoryID := c.Query("category_id")

	items, err := h.uc.ListItems(tenantID, categoryID)
	if err != nil {
		logging.HandlerError(c, "AdminMenu.ListItems", "service error", fiber.StatusBadRequest, "items_list_failed", err, "tenant_id", tenantID, "category_id", categoryID)
		return fiber.ErrBadRequest
	}

	resp := make([]itemResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, newItemResponse(item))
	}

	logging.HandlerInfo(c, "AdminMenu.ListItems", "items listed", fiber.StatusOK, "items_listed", "tenant_id", tenantID, "category_id", categoryID, "count", len(resp))
	return c.JSON(resp)
}

// CreateItem persists a new menu item for the tenant.
func (h *AdminMenuHandler) CreateItem(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)

	var payload map[string]any
	if err := c.BodyParser(&payload); err != nil {
		logging.HandlerError(c, "AdminMenu.CreateItem", "failed to parse body", fiber.StatusBadRequest, "invalid_body", err, "tenant_id", tenantID)
		return fiber.ErrBadRequest
	}

	item, err := h.uc.CreateItem(tenantID, payload)
	if err != nil {
		logging.HandlerError(c, "AdminMenu.CreateItem", "service error", fiber.StatusBadRequest, "item_create_failed", err, "tenant_id", tenantID)
		return fiber.ErrBadRequest
	}

	resp := newItemResponse(*item)
	logging.HandlerInfo(c, "AdminMenu.CreateItem", "item created", fiber.StatusCreated, "item_created", "tenant_id", tenantID, "item_id", resp.ID)
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// ReplaceItem performs a full update on an existing menu item.
func (h *AdminMenuHandler) ReplaceItem(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	itemID := c.Params("id")

	var payload map[string]any
	if err := c.BodyParser(&payload); err != nil {
		logging.HandlerError(c, "AdminMenu.ReplaceItem", "failed to parse body", fiber.StatusBadRequest, "invalid_body", err, "tenant_id", tenantID, "item_id", itemID)
		return fiber.ErrBadRequest
	}

	item, err := h.uc.ReplaceItem(tenantID, itemID, payload)
	if err != nil {
		logging.HandlerError(c, "AdminMenu.ReplaceItem", "service error", fiber.StatusBadRequest, "item_replace_failed", err, "tenant_id", tenantID, "item_id", itemID)
		return fiber.ErrBadRequest
	}

	resp := newItemResponse(*item)
	logging.HandlerInfo(c, "AdminMenu.ReplaceItem", "item replaced", fiber.StatusOK, "item_replaced", "tenant_id", tenantID, "item_id", itemID)
	return c.JSON(resp)
}

// PatchItem partially updates fields on an existing item.
func (h *AdminMenuHandler) PatchItem(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	itemID := c.Params("id")

	var payload map[string]any
	if err := c.BodyParser(&payload); err != nil {
		logging.HandlerError(c, "AdminMenu.PatchItem", "failed to parse body", fiber.StatusBadRequest, "invalid_body", err, "tenant_id", tenantID, "item_id", itemID)
		return fiber.ErrBadRequest
	}

	item, err := h.uc.PatchItem(tenantID, itemID, payload)
	if err != nil {
		logging.HandlerError(c, "AdminMenu.PatchItem", "service error", fiber.StatusBadRequest, "item_patch_failed", err, "tenant_id", tenantID, "item_id", itemID)
		return fiber.ErrBadRequest
	}

	resp := newItemResponse(*item)
	logging.HandlerInfo(c, "AdminMenu.PatchItem", "item patched", fiber.StatusOK, "item_patched", "tenant_id", tenantID, "item_id", itemID)
	return c.JSON(resp)
}

// DeleteItem removes an item from the tenant's catalogue.
func (h *AdminMenuHandler) DeleteItem(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	itemID := c.Params("id")

	if err := h.uc.DeleteItem(tenantID, itemID); err != nil {
		logging.HandlerError(c, "AdminMenu.DeleteItem", "service error", fiber.StatusBadRequest, "item_delete_failed", err, "tenant_id", tenantID, "item_id", itemID)
		return fiber.ErrBadRequest
	}

	logging.HandlerInfo(c, "AdminMenu.DeleteItem", "item deleted", fiber.StatusNoContent, "item_deleted", "tenant_id", tenantID, "item_id", itemID)
	return c.SendStatus(fiber.StatusNoContent)
}

// ToggleOOS flips the availability of an item through the out-of-stock flag.
func (h *AdminMenuHandler) ToggleOOS(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	itemID := c.Params("id")

	var payload struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.BodyParser(&payload); err != nil {
		logging.HandlerError(c, "AdminMenu.ToggleOOS", "failed to parse body", fiber.StatusBadRequest, "invalid_body", err, "tenant_id", tenantID, "item_id", itemID)
		return fiber.ErrBadRequest
	}

	item, err := h.uc.ToggleOOS(tenantID, itemID, payload.IsActive)
	if err != nil {
		logging.HandlerError(c, "AdminMenu.ToggleOOS", "service error", fiber.StatusBadRequest, "item_toggle_failed", err, "tenant_id", tenantID, "item_id", itemID, "is_active", payload.IsActive)
		return fiber.ErrBadRequest
	}

	resp := newItemResponse(*item)
	logging.HandlerInfo(c, "AdminMenu.ToggleOOS", "item toggled", fiber.StatusOK, "item_toggled", "tenant_id", tenantID, "item_id", itemID, "is_active", payload.IsActive)
	return c.JSON(resp)
}

// ListItemOptions returns all modifer options for a given item.
func (h *AdminMenuHandler) ListItemOptions(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	itemID := c.Params("id")

	opts, err := h.uc.ListItemOptions(itemID, tenantID)
	if err != nil {
		logging.HandlerError(c, "AdminMenu.ListItemOptions", "service error", fiber.StatusBadRequest, "options_list_failed", err, "tenant_id", tenantID, "item_id", itemID)
		return fiber.ErrBadRequest
	}

	resp := make([]optionResponse, 0, len(opts))
	for _, opt := range opts {
		resp = append(resp, newOptionResponse(opt))
	}

	logging.HandlerInfo(c, "AdminMenu.ListItemOptions", "item options listed", fiber.StatusOK, "options_listed", "tenant_id", tenantID, "item_id", itemID, "count", len(resp))
	return c.JSON(resp)
}

// CreateItemOption persists a new option for a given item.
func (h *AdminMenuHandler) CreateItemOption(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	itemID := c.Params("id")

	var payload map[string]any
	if err := c.BodyParser(&payload); err != nil {
		logging.HandlerError(c, "AdminMenu.CreateItemOption", "failed to parse body", fiber.StatusBadRequest, "invalid_body", err, "tenant_id", tenantID, "item_id", itemID)
		return fiber.ErrBadRequest
	}

	opt, err := h.uc.CreateItemOption(itemID, tenantID, payload)
	if err != nil {
		logging.HandlerError(c, "AdminMenu.CreateItemOption", "service error", fiber.StatusBadRequest, "option_create_failed", err, "tenant_id", tenantID, "item_id", itemID)
		return fiber.ErrBadRequest
	}

	resp := newOptionResponse(*opt)
	logging.HandlerInfo(c, "AdminMenu.CreateItemOption", "item option created", fiber.StatusCreated, "option_created", "tenant_id", tenantID, "item_id", itemID, "option_id", resp.ID)
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// ListOptionValues returns all values belonging to a specific option.
func (h *AdminMenuHandler) ListOptionValues(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	optionID := c.Params("option_id")

	values, err := h.uc.ListOptionValues(optionID, tenantID)
	if err != nil {
		logging.HandlerError(c, "AdminMenu.ListOptionValues", "service error", fiber.StatusBadRequest, "option_values_list_failed", err, "tenant_id", tenantID, "option_id", optionID)
		return fiber.ErrBadRequest
	}

	resp := make([]optionValueResponse, 0, len(values))
	for _, val := range values {
		resp = append(resp, newOptionValueResponse(val))
	}

	logging.HandlerInfo(c, "AdminMenu.ListOptionValues", "option values listed", fiber.StatusOK, "option_values_listed", "tenant_id", tenantID, "option_id", optionID, "count", len(resp))
	return c.JSON(resp)
}

// CreateOptionValue persists a new value for a given option.
func (h *AdminMenuHandler) CreateOptionValue(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	optionID := c.Params("option_id")

	var payload map[string]any
	if err := c.BodyParser(&payload); err != nil {
		logging.HandlerError(c, "AdminMenu.CreateOptionValue", "failed to parse body", fiber.StatusBadRequest, "invalid_body", err, "tenant_id", tenantID, "option_id", optionID)
		return fiber.ErrBadRequest
	}

	val, err := h.uc.CreateOptionValue(optionID, tenantID, payload)
	if err != nil {
		logging.HandlerError(c, "AdminMenu.CreateOptionValue", "service error", fiber.StatusBadRequest, "option_value_create_failed", err, "tenant_id", tenantID, "option_id", optionID)
		return fiber.ErrBadRequest
	}

	resp := newOptionValueResponse(*val)
	logging.HandlerInfo(c, "AdminMenu.CreateOptionValue", "option value created", fiber.StatusCreated, "option_value_created", "tenant_id", tenantID, "option_id", optionID, "value_id", resp.ID)
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// GenerateQR returns a QR code URL for a table belonging to the tenant.
func (h *AdminMenuHandler) GenerateQR(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	tableID := c.Params("id")

	url, err := h.uc.GenerateTableQR(tableID, tenantID)
	if err != nil {
		logging.HandlerError(c, "AdminMenu.GenerateQR", "service error", fiber.StatusBadRequest, "qr_generate_failed", err, "tenant_id", tenantID, "table_id", tableID)
		return fiber.ErrBadRequest
	}

	logging.HandlerInfo(c, "AdminMenu.GenerateQR", "qr generated", fiber.StatusOK, "qr_generated", "tenant_id", tenantID, "table_id", tableID)
	return c.JSON(fiber.Map{"url": url})
}

// newCategoryResponse converts a domain category into its JSON representation.
func newCategoryResponse(cat domain.Category) categoryResponse {
	return categoryResponse{
		ID:       cat.ID,
		TenantID: cat.TenantID,
		Name:     cat.Name,
		Sort:     cat.Sort,
		IsActive: cat.IsActive,
	}
}

// newItemResponse converts a domain item into its JSON representation.
func newItemResponse(item domain.Item) itemResponse {
	return itemResponse{
		ID:          item.ID,
		TenantID:    item.TenantID,
		CategoryID:  item.CategoryID,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
		PhotoURL:    item.PhotoURL,
		Flags:       item.Flags,
		IsActive:    item.IsActive,
	}
}

// newOptionResponse converts a domain option into its JSON representation.
func newOptionResponse(opt domain.ItemOption) optionResponse {
	return optionResponse{
		ID:       opt.ID,
		ItemID:   opt.ItemID,
		Name:     opt.Name,
		Type:     opt.Type,
		Required: opt.Required,
	}
}

// newOptionValueResponse converts a domain option value into its JSON representation.
func newOptionValueResponse(val domain.ItemOptionValue) optionValueResponse {
	return optionValueResponse{
		ID:         val.ID,
		OptionID:   val.OptionID,
		Label:      val.Label,
		DeltaPrice: val.DeltaPrice,
	}
}
