package handler

import (
	"qrmenu/internal/domain"
	"qrmenu/internal/platform/logging"

	"github.com/gofiber/fiber/v2"
)

// Alias the domain type so every layer shares the same contract.
type OrderCreateRequest = domain.OrderCreateRequest

// Alias for individual order items if needed by callers.
type OrderItemCreate = domain.OrderItemCreate

type OrderCreator interface {
	CreateGuestOrder(req OrderCreateRequest) (orderID string, status string, err error)
}

type OrderPublicHandler struct{ svc OrderCreator }

func NewOrderPublicHandler(s OrderCreator) *OrderPublicHandler { return &OrderPublicHandler{svc: s} }

func (h *OrderPublicHandler) Create(c *fiber.Ctx) error {
	var req OrderCreateRequest
	if err := c.BodyParser(&req); err != nil {
		logging.HandlerError(c, "OrderPublic.Create", "failed to parse body", fiber.StatusBadRequest, "invalid_body", err)
		return fiber.ErrBadRequest
	}

	// Perform lightweight validation before invoking the use case.
	if req.Tenant == "" || req.TableToken == "" || req.GuestSession == "" || len(req.Items) == 0 {
		logging.HandlerError(c, "OrderPublic.Create", "invalid payload", fiber.StatusBadRequest, "invalid_payload", fiber.ErrBadRequest, "tenant", req.Tenant, "table_token", req.TableToken)
		return fiber.ErrBadRequest
	}

	id, status, err := h.svc.CreateGuestOrder(req)
	if err != nil {
		logging.HandlerError(c, "OrderPublic.Create", "failed to create order", fiber.StatusBadRequest, "order_create_failed", err, "tenant", req.Tenant, "table_token", req.TableToken)
		return fiber.ErrBadRequest
	}
	logging.HandlerInfo(c, "OrderPublic.Create", "guest order created", fiber.StatusCreated, "order_created", "order_id", id, "tenant", req.Tenant)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"order_id": id,
		"status":   status,
	})
}
