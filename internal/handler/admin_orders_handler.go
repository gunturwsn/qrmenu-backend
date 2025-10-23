package handler

import (
	"qrmenu/internal/domain"
	"qrmenu/internal/repository"

	"github.com/gofiber/fiber/v2"
)

// <- ALIAS: samakan tipe OrdersPage di handler dengan repository
type OrdersPage = repository.OrdersPage

// AdminOrdersQuery dikonsumsi handler; diimplementasikan oleh usecase.AdminOrdersUC
type AdminOrdersQuery interface {
	List(tenantID, status, cursor string) (OrdersPage, error)
	UpdateStatus(tenantID, id, status string) (*domain.Order, error)
}

type AdminOrdersHandler struct{ q AdminOrdersQuery }

func NewAdminOrdersHandler(q AdminOrdersQuery) *AdminOrdersHandler { return &AdminOrdersHandler{q: q} }

// GET /admin/orders?status=&cursor=
func (h *AdminOrdersHandler) List(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	status := c.Query("status")
	cursor := c.Query("cursor")

	page, err := h.q.List(tenantID, status, cursor)
	if err != nil {
		return fiber.ErrBadRequest
	}
	return c.JSON(page)
}

// PATCH /admin/orders/:id/status
func (h *AdminOrdersHandler) PatchStatus(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	id := c.Params("id")

	var body struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil || body.Status == "" {
		return fiber.ErrBadRequest
	}

	ord, err := h.q.UpdateStatus(tenantID, id, body.Status)
	if err != nil {
		return fiber.ErrBadRequest
	}
	return c.JSON(ord)
}
