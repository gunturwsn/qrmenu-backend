package handler

import (
	"qrmenu/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type MenuHandler struct{ uc usecase.MenuUC }

func NewMenuHandler(uc usecase.MenuUC) *MenuHandler { return &MenuHandler{uc: uc} }

func (h *MenuHandler) Get(c *fiber.Ctx) error {
	code := c.Query("tenant_code")
	if code == "" {
		return fiber.ErrBadRequest
	}

	res, err := h.uc.GetMenuByTenantCode(code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "tenant not found"})
	}
	return c.JSON(res)
}
