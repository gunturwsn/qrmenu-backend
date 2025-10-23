package handler

import (
	"github.com/gofiber/fiber/v2"
)

type TableResolver interface {
	ResolveByToken(token string) (tenant any, table any, err error)
}

type TableHandler struct{ svc TableResolver }

func NewTableHandler(s TableResolver) *TableHandler { return &TableHandler{svc: s} }

func (h *TableHandler) Resolve(c *fiber.Ctx) error {
	token := c.Params("token")
	tenant, table, err := h.svc.ResolveByToken(token)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(fiber.Map{"tenant": tenant, "table": table})
}
