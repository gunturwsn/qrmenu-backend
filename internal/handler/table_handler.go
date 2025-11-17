package handler

import (
	"github.com/gofiber/fiber/v2"

	"qrmenu/internal/platform/logging"
)

type TableResolver interface {
	ResolveByToken(token string) (tenant any, table any, err error)
}

type TableHandler struct{ svc TableResolver }

func NewTableHandler(s TableResolver) *TableHandler { return &TableHandler{svc: s} }

func (h *TableHandler) Resolve(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		logging.HandlerError(c, "Table.Resolve", "token missing", fiber.StatusBadRequest, "token_missing", fiber.ErrBadRequest)
		return fiber.ErrBadRequest
	}
	tenant, table, err := h.svc.ResolveByToken(token)
	if err != nil {
		logging.HandlerError(c, "Table.Resolve", "resolve failed", fiber.StatusNotFound, "table_not_found", err, "token", token)
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	logging.HandlerInfo(c, "Table.Resolve", "table resolved", fiber.StatusOK, "table_resolved", "token", token)
	return c.JSON(fiber.Map{"tenant": tenant, "table": table})
}
