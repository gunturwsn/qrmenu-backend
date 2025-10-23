package handler

import (
	"qrmenu/internal/domain"

	"github.com/gofiber/fiber/v2"
)

// Pakai alias ke domain agar semua layer share satu tipe yang sama
type OrderCreateRequest = domain.OrderCreateRequest

// Kalau perlu, kamu juga bisa alias item-nya:
type OrderItemCreate = domain.OrderItemCreate

type OrderCreator interface {
	CreateGuestOrder(req OrderCreateRequest) (orderID string, status string, err error)
}

type OrderPublicHandler struct{ svc OrderCreator }

func NewOrderPublicHandler(s OrderCreator) *OrderPublicHandler { return &OrderPublicHandler{svc: s} }

func (h *OrderPublicHandler) Create(c *fiber.Ctx) error {
	var req OrderCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	// (opsional) validasi ringan
	if req.Tenant == "" || req.TableToken == "" || req.GuestSession == "" || len(req.Items) == 0 {
		return fiber.ErrBadRequest
	}

	id, status, err := h.svc.CreateGuestOrder(req)
	if err != nil {
		return fiber.ErrBadRequest
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"order_id": id,
		"status":   status,
	})
}
