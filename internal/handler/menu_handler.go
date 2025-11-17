package handler

import (
	"qrmenu/internal/platform/logging"
	"qrmenu/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type MenuHandler struct{ uc usecase.MenuUC }

func NewMenuHandler(uc usecase.MenuUC) *MenuHandler { return &MenuHandler{uc: uc} }

func (h *MenuHandler) Get(c *fiber.Ctx) error {
	code := c.Query("tenant_code")
	if code == "" {
		logging.HandlerError(c, "Menu.Get", "tenant_code missing", fiber.StatusBadRequest, "tenant_code_missing", fiber.ErrBadRequest)
		return fiber.ErrBadRequest
	}

	res, err := h.uc.GetMenuByTenantCode(code)
	if err != nil {
		logging.HandlerError(c, "Menu.Get", "menu lookup failed", fiber.StatusNotFound, "menu_not_found", err, "tenant_code", code)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "tenant not found"})
	}
	logging.HandlerInfo(c, "Menu.Get", "menu served", fiber.StatusOK, "menu_success", "tenant_code", code)
	return c.JSON(res)
}
