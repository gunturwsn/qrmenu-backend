package usecase

import (
	"qrmenu/internal/domain"
	"qrmenu/internal/platform/logging"
	"qrmenu/internal/repository"
)

type OrderUC struct {
	repo repository.OrderRepository
}

func NewOrderUC(r repository.OrderRepository) *OrderUC { return &OrderUC{repo: r} }

// CreateGuestOrder forwards the domain payload straight to the repository.
// Return values: orderID, status(string), error.
func (u *OrderUC) CreateGuestOrder(req domain.OrderCreateRequest) (string, string, error) {
	logging.UsecaseInfo("Order.CreateGuestOrder", "creating guest order", "order_create_requested", "tenant", req.Tenant, "table_token", req.TableToken, "items", len(req.Items))
	id, st, err := u.repo.CreateGuestOrder(req)
	if err != nil {
		logging.UsecaseError("Order.CreateGuestOrder", "repository error", "order_create_failed", err, "tenant", req.Tenant, "table_token", req.TableToken)
		return "", "", err
	}
	logging.UsecaseInfo("Order.CreateGuestOrder", "guest order created", "order_created", "tenant", req.Tenant, "order_id", id, "status", st)
	return id, string(st), nil
}
