package handler

import (
	"errors"
	"qrmenu/internal/domain"
	"qrmenu/internal/platform/logging"
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
		logging.HandlerError(c, "AdminOrders.List", "query failed", fiber.StatusBadRequest, "orders_query_failed", err, "tenant_id", tenantID, "status", status, "cursor", cursor)
		return fiber.ErrBadRequest
	}
	logging.HandlerInfo(c, "AdminOrders.List", "orders retrieved", fiber.StatusOK, "orders_listed", "tenant_id", tenantID, "status", status, "count", len(page.Data))
	return c.JSON(page)
}

// PATCH /admin/orders/:id/status
func (h *AdminOrdersHandler) PatchStatus(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	id := c.Params("id")

	var body struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		logging.HandlerError(c, "AdminOrders.PatchStatus", "failed to parse body", fiber.StatusBadRequest, "invalid_body", err, "tenant_id", tenantID, "order_id", id)
		return fiber.ErrBadRequest
	}
	if body.Status == "" {
		logging.HandlerError(c, "AdminOrders.PatchStatus", "status missing", fiber.StatusBadRequest, "status_missing", errors.New("status required"), "tenant_id", tenantID, "order_id", id)
		return fiber.ErrBadRequest
	}

	ord, err := h.q.UpdateStatus(tenantID, id, body.Status)
	if err != nil {
		logging.HandlerError(c, "AdminOrders.PatchStatus", "update failed", fiber.StatusBadRequest, "order_status_update_failed", err, "tenant_id", tenantID, "order_id", id, "status", body.Status)
		return fiber.ErrBadRequest
	}
	logging.HandlerInfo(c, "AdminOrders.PatchStatus", "status updated", fiber.StatusOK, "status_updated", "tenant_id", tenantID, "order_id", id, "status", body.Status)
	return c.JSON(ord)
}
